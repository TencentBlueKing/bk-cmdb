/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package metadata defines cmdb data syncer's metadata info
package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	commonutil "configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/transfer-service/app/options"
	"configcenter/src/storage/driver/mongodb"

	"github.com/tidwall/gjson"
)

// Metadata is cmdb data syncer's metadata info
type Metadata struct {
	role options.SyncRole
	// InnerIDInfo is inner data id info of this environment
	InnerIDInfo *options.InnerDataIDConf
	// blueking is the blueking biz info, is used to skip the resource in blueking biz for source environment
	blueking *bluekingBizInfo
}

// BluekingBizID is the blueking biz info, resource in blueking biz should not be synced
type bluekingBizInfo struct {
	bizID         int64
	lock          sync.RWMutex
	hostModuleMap map[int64]map[int64]struct{}
}

func (m *Metadata) getBluekingHostIDs() []int64 {
	if m.blueking == nil {
		return make([]int64, 0)
	}

	m.blueking.lock.RLock()
	defer m.blueking.lock.RUnlock()

	hostIDs := make([]int64, 0)
	for hostID := range m.blueking.hostModuleMap {
		hostIDs = append(hostIDs, hostID)
	}
	return hostIDs
}

func (m *Metadata) isHostInBluekingBiz(hostID int64) bool {
	if m.blueking == nil {
		return false
	}

	m.blueking.lock.RLock()
	defer m.blueking.lock.RUnlock()

	_, exists := m.blueking.hostModuleMap[hostID]
	return exists
}

// NewMetadata new cmdb data syncer's metadata info
func NewMetadata(role options.SyncRole) (*Metadata, error) {
	meta := &Metadata{
		role: role,
		InnerIDInfo: &options.InnerDataIDConf{
			HostPool: new(options.HostPoolInfo),
		},
	}

	err := meta.init()
	if err != nil {
		return nil, fmt.Errorf("init metadata failed, err: %v", err)
	}

	return meta, nil
}

func (m *Metadata) init() error {
	ctx := commonutil.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)

	// get inner data id info
	hostPoolBiz := new(metadata.BizBasicInfo)
	hostPoolBizCond := mapstr.MapStr{common.BKDefaultFiled: common.DefaultAppFlag}
	if err := mongodb.Client().Table(common.BKTableNameBaseApp).Find(hostPoolBizCond).Fields(common.BKAppIDField).
		One(ctx, &hostPoolBiz); err != nil {
		blog.Errorf("get host pool biz id by cond(%+v) failed, err: %v", hostPoolBizCond, err)
		return err
	}
	m.InnerIDInfo.HostPool.Biz = hostPoolBiz.BizID

	hostPoolSet := new(metadata.SetInst)
	hostPoolSetCond := mapstr.MapStr{
		common.BKAppIDField:   hostPoolBiz.BizID,
		common.BKDefaultFiled: common.DefaultResSetFlag,
	}
	if err := mongodb.Client().Table(common.BKTableNameBaseSet).Find(hostPoolSetCond).Fields(common.BKSetIDField).
		One(ctx, &hostPoolSet); err != nil {
		blog.Errorf("get host pool set id by cond(%+v) failed, err: %v", hostPoolSetCond, err)
		return err
	}
	m.InnerIDInfo.HostPool.Set = hostPoolSet.SetID

	hostPoolModule := new(metadata.ModuleInst)
	hostPoolModuleCond := mapstr.MapStr{
		common.BKSetIDField:   hostPoolSet.SetID,
		common.BKDefaultFiled: common.DefaultResModuleFlag,
	}
	if err := mongodb.Client().Table(common.BKTableNameBaseModule).Find(hostPoolModuleCond).
		Fields(common.BKModuleIDField).One(ctx, &hostPoolModule); err != nil {
		blog.Errorf("get host pool module id by cond(%+v) failed, err: %v", hostPoolModuleCond, err)
		return err
	}
	m.InnerIDInfo.HostPool.Module = hostPoolModule.ModuleID

	// get blueking biz id and host ids for source environment
	if m.role != options.SyncRoleSrc {
		return nil
	}
	m.blueking = new(bluekingBizInfo)

	bluekingBiz := new(metadata.BizBasicInfo)
	bluekingBizCond := mapstr.MapStr{common.BKAppNameField: common.BKAppName}
	if err := mongodb.Client().Table(common.BKTableNameBaseApp).Find(bluekingBizCond).Fields(common.BKAppIDField).
		One(ctx, &bluekingBiz); err != nil {
		blog.Errorf("get blueking biz id by cond(%+v) failed, err: %v", hostPoolBizCond, err)
		return err
	}
	m.blueking.bizID = bluekingBiz.BizID

	bkHostRel := make([]metadata.ModuleHost, 0)
	bkHostRelCond := mapstr.MapStr{common.BKAppIDField: m.blueking.bizID}
	if err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(bkHostRelCond).
		Fields(common.BKHostIDField, common.BKModuleIDField).All(ctx, &bkHostRel); err != nil {
		blog.Errorf("get blueking host relations by cond(%+v) failed, err: %v", bkHostRelCond, err)
		return err
	}

	m.blueking.hostModuleMap = make(map[int64]map[int64]struct{})
	for _, relation := range bkHostRel {
		_, exists := m.blueking.hostModuleMap[relation.HostID]
		if !exists {
			m.blueking.hostModuleMap[relation.HostID] = make(map[int64]struct{})
		}
		m.blueking.hostModuleMap[relation.HostID][relation.ModuleID] = struct{}{}
	}

	return nil
}

// GetCommonObjIDs get all objIDs and quoted objIDs for object instance resource full sync, do not include inner objects
func (m *Metadata) GetCommonObjIDs() ([]string, []string, error) {
	ctx := commonutil.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)

	objects := make([]metadata.Object, 0)
	objCond := mapstr.MapStr{common.BKIsPre: false}
	err := mongodb.Client().Table(common.BKTableNameObjDes).Find(objCond).Fields(common.BKObjIDField).All(ctx, &objects)
	if err != nil {
		blog.Errorf("get all object ids failed, err: %v, cond: %+v", err, objCond)
		return nil, nil, err
	}

	quoteRelations := make([]metadata.ModelQuoteRelation, 0)
	err = mongodb.Client().Table(common.BKTableNameModelQuoteRelation).Find(nil).All(ctx, &quoteRelations)
	if err != nil {
		blog.Errorf("get model quote relations failed, err: %v", err)
		return nil, nil, err
	}

	quotedObjMap := make(map[string]struct{})
	for _, relation := range quoteRelations {
		quotedObjMap[relation.DestModel] = struct{}{}
	}

	objIDs, quotedObjIDs := make([]string, 0), make([]string, 0)
	for _, object := range objects {
		_, exists := quotedObjMap[object.ObjectID]
		if exists {
			quotedObjIDs = append(quotedObjIDs, object.ObjectID)
			continue
		}
		objIDs = append(objIDs, object.ObjectID)
	}

	return objIDs, quotedObjIDs, nil
}

var bizRelatedResTypeMap = map[types.ResType]struct{}{types.Biz: {}, types.ObjectInstance: {}, types.Set: {},
	types.Module: {}, types.HostRelation: {}, types.ServiceInstance: {}, types.Process: {}, types.ProcessRelation: {}}

// AddListCond add list condition for resource full sync list data logics
func (m *Metadata) AddListCond(resType types.ResType, cond mapstr.MapStr) mapstr.MapStr {
	extraCond := make(mapstr.MapStr)
	switch resType {
	case types.Biz, types.ObjectInstance:
		// do not sync host pool and blueking biz
		if m.blueking == nil || m.blueking.bizID == 0 {
			extraCond[common.BKAppIDField] = mapstr.MapStr{common.BKDBNE: m.InnerIDInfo.HostPool.Biz}
			return mergeCond(cond, extraCond)
		}
		extraCond[common.BKAppIDField] = mapstr.MapStr{common.BKDBNIN: []int64{m.InnerIDInfo.HostPool.Biz,
			m.blueking.bizID}}
		return mergeCond(cond, extraCond)
	case types.Set:
		// do not sync host pool set
		extraCond[common.BKSetIDField] = mapstr.MapStr{common.BKDBNE: m.InnerIDInfo.HostPool.Set}
	case types.Module:
		// do not sync host pool module
		extraCond[common.BKModuleIDField] = mapstr.MapStr{common.BKDBNE: m.InnerIDInfo.HostPool.Module}
	case types.Host:
		// do not sync blueking hosts
		hostIDs := m.getBluekingHostIDs()
		if len(hostIDs) == 0 {
			return cond
		}
		extraCond[common.BKHostIDField] = mapstr.MapStr{common.BKDBNIN: hostIDs}
		return mergeCond(cond, extraCond)
	case types.InstAsst:
		// do not sync blueking host associations
		hostIDs := m.getBluekingHostIDs()
		if len(hostIDs) == 0 {
			return cond
		}
		extraCond = mapstr.MapStr{
			common.BKDBNOR: []mapstr.MapStr{{
				common.BKObjIDField:  common.BKInnerObjIDHost,
				common.BKInstIDField: mapstr.MapStr{common.BKDBIN: hostIDs},
			}, {
				common.BKAsstObjIDField:  common.BKInnerObjIDHost,
				common.BKAsstInstIDField: mapstr.MapStr{common.BKDBIN: hostIDs},
			}},
		}
	}

	// do not sync resource in blueking biz
	if m.blueking == nil || m.blueking.bizID == 0 {
		return mergeCond(cond, extraCond)
	}

	_, exists := bizRelatedResTypeMap[resType]
	if exists {
		extraCond[common.BKAppIDField] = mapstr.MapStr{common.BKDBNE: m.blueking.bizID}
	}

	return mergeCond(cond, extraCond)
}

func mergeCond(cond mapstr.MapStr, extraCond mapstr.MapStr) mapstr.MapStr {
	if len(extraCond) == 0 {
		return cond
	}

	if len(cond) == 0 {
		return extraCond
	}

	return mapstr.MapStr{common.BKDBAND: []mapstr.MapStr{cond, extraCond}}
}

// ParseEventDetail parse event detail for event watch logics, returns if the event needs sync
func (m *Metadata) ParseEventDetail(eventType watch.EventType, resType types.ResType, oid string,
	detail json.RawMessage) (*types.EventInfo, bool) {

	event := &types.EventInfo{
		EventType: eventType,
		ResType:   resType,
		Oid:       oid,
		Detail:    detail,
	}

	if resType == types.HostRelation {
		return m.parseHostRelEvent(event)
	}

	// do not sync resource in blueking biz
	_, exists := bizRelatedResTypeMap[resType]
	if exists {
		if m.blueking != nil && gjson.GetBytes(detail, common.BKAppIDField).Int() == m.blueking.bizID {
			return nil, false
		}
	}

	switch resType {
	case types.Biz:
		// do not sync host pool biz
		return event, gjson.GetBytes(detail, common.BKAppIDField).Int() != m.InnerIDInfo.HostPool.Biz
	case types.Set:
		// do not sync host pool set
		return event, gjson.GetBytes(detail, common.BKSetIDField).Int() != m.InnerIDInfo.HostPool.Set
	case types.Module:
		// do not sync host pool module
		return event, gjson.GetBytes(detail, common.BKModuleIDField).Int() != m.InnerIDInfo.HostPool.Module
	case types.Host:
		// do not sync blueking hosts
		return event, !m.isHostInBluekingBiz(gjson.GetBytes(detail, common.BKHostIDField).Int())
	case types.ObjectInstance:
		event.SubRes = []string{gjson.GetBytes(detail, common.BKObjIDField).String()}
		// do not sync host pool biz
		return event, gjson.GetBytes(detail, common.BKAppIDField).Int() != m.InnerIDInfo.HostPool.Biz
	case types.InstAsst:
		event.SubRes = []string{gjson.GetBytes(detail, common.BKObjIDField).String(),
			gjson.GetBytes(detail, common.BKAsstObjIDField).String()}
		asstInstIDs := []int64{gjson.GetBytes(detail, common.BKInstIDField).Int(),
			gjson.GetBytes(detail, common.BKAsstInstIDField).Int()}

		// do not sync blueking host associations
		for i, subRes := range event.SubRes {
			if subRes == common.BKInnerObjIDHost {
				if m.isHostInBluekingBiz(asstInstIDs[i]) {
					return nil, false
				}
			}
		}
		return event, true
	}

	return event, true
}

// parseHostRelEvent parse host relation event
func (m *Metadata) parseHostRelEvent(event *types.EventInfo) (*types.EventInfo, bool) {
	// if host relation event is not in blueking biz, sync this event
	if m.blueking == nil || m.blueking.bizID == 0 {
		return event, true
	}
	if gjson.GetBytes(event.Detail, common.BKAppIDField).Int() != m.blueking.bizID {
		return event, true
	}

	// update metadata blueking host id info by host relation event
	hostID := gjson.GetBytes(event.Detail, common.BKHostIDField).Int()
	moduleID := gjson.GetBytes(event.Detail, common.BKModuleIDField).Int()

	m.blueking.lock.Lock()
	defer m.blueking.lock.Unlock()

	if event.EventType == watch.Delete {
		_, exists := m.blueking.hostModuleMap[hostID]
		if exists {
			delete(m.blueking.hostModuleMap[hostID], moduleID)
		}
		if len(m.blueking.hostModuleMap[hostID]) == 0 {
			delete(m.blueking.hostModuleMap, hostID)
			// host is not in blueking biz, change this event to create host event
			host := new(metadata.HostMapStr)
			hostCond := mapstr.MapStr{common.BKHostIDField: hostID}
			err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(hostCond).One(context.Background(), host)
			if err != nil {
				if mongodb.Client().IsNotFoundError(err) {
					return nil, false
				}
				blog.Errorf("get not in blueking host by id %d failed, err: %v", hostID, err)
				return nil, false
			}

			hostJson, err := json.Marshal(host)
			if err != nil {
				blog.Errorf("marshal host(%+v) failed, err: %v", hostID, err)
				return nil, false
			}
			event = &types.EventInfo{
				EventType: watch.Create,
				ResType:   types.Host,
				Oid:       strconv.FormatInt(hostID, 10),
				Detail:    hostJson,
			}
			return event, true
		}
		return nil, false
	}

	_, exists := m.blueking.hostModuleMap[hostID]
	if !exists {
		m.blueking.hostModuleMap[hostID] = make(map[int64]struct{})
		// host is transferred to blueking biz, change this event to delete host event
		event = &types.EventInfo{
			EventType: watch.Delete,
			ResType:   types.Host,
			Oid:       strconv.FormatInt(hostID, 10),
			Detail:    json.RawMessage(fmt.Sprintf(`{"%s":%d}`, common.BKHostIDField, hostID)),
		}
		return event, true
	}
	m.blueking.hostModuleMap[hostID][moduleID] = struct{}{}
	return nil, false
}
