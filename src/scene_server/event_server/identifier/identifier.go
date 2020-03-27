/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package identifier

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/types"
	"configcenter/src/storage/dal"

	"gopkg.in/redis.v5"
)

var delayTime = time.Second * 30

// TODO: add event for layer
var hostIndentDiffFields = map[string][]string{
	common.BKInnerObjIDApp: {
		common.BKAppNameField,
	},
	common.BKInnerObjIDSet: {
		common.BKSetNameField,
		common.BKSetStatusField,
		common.BKSetEnvField,
	},
	common.BKInnerObjIDModule: {
		common.BKModuleNameField,
	},
	common.BKInnerObjIDPlat: {
		common.BKCloudNameField,
	},
	common.BKInnerObjIDProc: {
		common.BKProcessNameField,
		common.BKFuncIDField,
		common.BKFuncName,
		common.BKBindIP,
		common.BKProtocol,
		common.BKPort,
		common.BKStartParamRegex,
	},
	common.BKInnerObjIDHost: {
		common.BKHostNameField,
		common.BKCloudIDField,
		common.BKHostInnerIPField,
		common.BKHostOuterIPField,
		common.BKOSTypeField,
		common.BKOSNameField,
		"bk_mem",
		"bk_cpu",
		"bk_disk",
	},
}

func (ih *IdentifierHandler) handleEvent(event *metadata.EventInstCtx) {
	if diffFields, ok := hostIndentDiffFields[event.ObjType]; ok &&
		event.Action == metadata.EventActionUpdate &&
		event.EventType == metadata.EventTypeInstData {
		ih.handleInstFieldChange(event, diffFields)
	} else if event.EventType == metadata.EventTypeRelation &&
		(event.ObjType == metadata.EventObjTypeModuleTransfer ||
			event.ObjType == metadata.EventObjTypeProcModule) {
		ih.handleHostRelationChange(event)
	}
}

func (ih *IdentifierHandler) handleInstFieldChange(event *metadata.EventInstCtx, diffFields []string) {
	rid := util.ExtractRequestIDFromContext(ih.ctx)
	blog.InfoJSON("identifier: handle inst %s, rid: %s", event, rid)

	hostIdentify := event.EventInst
	hostIdentify.Data = nil
	hostIdentify.EventType = metadata.EventTypeRelation
	hostIdentify.ObjType = objectTypeIdentifier
	hostIdentify.Action = metadata.EventActionUpdate

	for dataIndex := range event.Data {
		curData := event.Data[dataIndex].CurData.(map[string]interface{})
		preData := event.Data[dataIndex].PreData.(map[string]interface{})
		if hasChanged(curData, preData, diffFields...) == false {
			continue
		}

		instIDField := common.GetInstIDField(event.ObjType)

		instID, err := getInt(curData, instIDField)
		if err != nil || 0 == instID {
			blog.Errorf("identifier: convert instID failed the raw is %+v, rid: %s", curData[instIDField], rid)
			continue
		}

		inst, err := getCache(ih.ctx, ih.cache, ih.db, event.ObjType, instID, false)
		if err != nil {
			blog.Errorf("identifier: getCache error %+v, rid: %s", err, rid)
			continue
		}
		if nil == inst {
			blog.Errorf("identifier: inst == nil, continue, rid:%s", rid)
			continue
		}
		// save previous host identifier data for later usage.
		preIdentifier := *inst.ident
		for _, field := range diffFields {
			err = inst.set(field, curData[field])
			if err != nil {
				blog.Errorf("key %s, value: %s, convert error %s", field, curData[field], err.Error())
				break
			}
		}
		if err != nil {
			continue
		}

		err = inst.saveCache(ih.cache)
		if err != nil {
			blog.ErrorJSON("saveCache inst data %s failed, err: %s, rid: %s", inst.data, err, rid)
			continue
		}

		if common.BKInnerObjIDHost == event.ObjType {
			hostIdentify.ID = ih.cache.Incr(types.EventCacheEventIDKey).Val()
			d := metadata.EventData{CurData: inst.ident, PreData: preIdentifier}
			hostIdentify.Data = append(hostIdentify.Data, d)

			ih.cache.LPush(types.EventCacheEventQueueKey, &hostIdentify)
			blog.InfoJSON("identifier: pushed event inst %s, rid: %s", hostIdentify, rid)
		} else {
			if err := ih.handleRelatedInst(hostIdentify, event.ObjType, instID); err != nil {
				blog.Warnf("handleRelatedInst failed objType: %s, inst: %d, error: %v, rid: %s", event.ObjType, instID, err, rid)
			}
		}
	}
}

func (ih *IdentifierHandler) handleHostRelationChange(e *metadata.EventInstCtx) {
	rid := util.ExtractRequestIDFromContext(ih.ctx)
	blog.InfoJSON("identifier: handle inst %s, rid: %s", e, rid)

	hostIdentify := e.EventInst
	hostIdentify.Data = nil
	hostIdentify.EventType = metadata.EventTypeRelation
	hostIdentify.ObjType = objectTypeIdentifier
	hostIdentify.Action = metadata.EventActionUpdate

	go func() {
		time.Sleep(delayTime)

		hostIDs := make([]int64, 0)
		for index := range e.Data {
			var curData map[string]interface{}
			var ok bool
			if metadata.EventActionDelete == e.Action {
				curData, ok = e.Data[index].PreData.(map[string]interface{})
			} else {
				curData, ok = e.Data[index].CurData.(map[string]interface{})
			}
			if !ok {
				continue
			}

			hostID, err := getInt(curData, common.BKHostIDField)
			if err != nil || 0 == hostID {
				blog.Errorf("identifier: convert instID failed the raw is %+v, rid: %s", curData[common.BKHostIDField], rid)
				continue
			}

			hostIDs = append(hostIDs, hostID)
		}

		hostIDs = util.IntArrayUnique(hostIDs)

		if err := ih.handleHostBatch(hostIdentify, hostIDs); err != nil {
			blog.Warnf("handleHostBatch failed , hostIDs: %v, error: %v, rid: %s", hostIDs, err, rid)
		}
	}()
}

const objectTypeIdentifier = "hostidentifier"

func (ih *IdentifierHandler) handleRelatedInst(hostIdentify metadata.EventInst, objType string, instID int64) error {
	rid := util.ExtractRequestIDFromContext(ih.ctx)
	hostIDs, err := ih.findHost(objType, instID)
	if err != nil {
		blog.Warnf("identifier: find host failure: %v, rid: %s", err, rid)
		return err
	}
	blog.V(3).Infof("identifier: handleRelatedInst by objType %s, instID %d,  hostIDs: %v, rid: %s", objType, instID, hostIDs, rid)
	return ih.handleHostBatch(hostIdentify, hostIDs)
}

func (ih *IdentifierHandler) handleHostBatch(hostIdentify metadata.EventInst, hosIDs []int64) error {
	rid := util.ExtractRequestIDFromContext(ih.ctx)
	blog.V(3).Infof("identifier: handleHostBatch hostIDs: %v, rid: %s", hosIDs, rid)
	total := len(hosIDs)
	bufSize := 256
	index := 0
	leftIndex := 0

	for leftIndex < total {
		leftIndex = index + bufSize
		if leftIndex > total {
			leftIndex = total
		}
		hostIdentify.Data = nil

		hostIDKeys := make([]string, 0)
		for _, hostID := range hosIDs[index:leftIndex] {
			hostIDKeys = append(hostIDKeys, getInstCacheKey(common.BKInnerObjIDHost, hostID))
		}
		idens, err := ih.cache.MGet(hostIDKeys...).Result()
		if err != nil {
			blog.Warnf("identifier: ih.cache.MGet by %v,%v. we will try to fetch it from db instead, rid: %s", hostIDKeys, err, rid)
			idens = make([]interface{}, len(hostIDKeys))
			for index := range idens {
				// simulate that redis returns all nil
				idens[index] = nilStr
			}
		}
		for identIndex, idenVal := range idens {
			hostID := hosIDs[index:leftIndex][identIndex]
			iden := metadata.HostIdentifier{}
			var preIdentifier metadata.HostIdentifier

			err = json.Unmarshal([]byte(getString(idenVal)), &iden)
			if err != nil {
				blog.Warnf("hostID %d cache value %v unmarshal error: %s, fetch it from db, rid: %s", hostID, idenVal, err.Error(), rid)
				preInst, err := getCache(ih.ctx, ih.cache, ih.db, common.BKInnerObjIDHost, hostID, false)
				if err != nil {
					blog.Errorf("identifier: getCache error %+v, rid: %s", err, rid)
					continue
				}
				preIdentifier = *preInst.ident
			} else {
				preIdentifier = iden
			}

			inst, err := getCache(ih.ctx, ih.cache, ih.db, common.BKInnerObjIDHost, hostID, true)
			if err != nil {
				blog.Errorf("identifier: getCache error %+v, rid: %s", err, rid)
				continue
			}
			if nil == inst {
				continue
			}
			if err := inst.saveCache(ih.cache); err != nil {
				blog.ErrorJSON("saveCache inst data %s failed, err: %s, rid: %s", inst.data, err, rid)
			}

			if !reflect.DeepEqual(preIdentifier, *inst) {
				d := metadata.EventData{CurData: inst.ident, PreData: preIdentifier}
				hostIdentify.Data = append(hostIdentify.Data, d)
			}
		}
		index += bufSize

		hostIdentify.ID = ih.cache.Incr(types.EventCacheEventIDKey).Val()
		if err = ih.cache.LPush(types.EventCacheEventQueueKey, &hostIdentify).Err(); err != nil {
			blog.Warnf("identifier: push event inst %v failure %v, rid: %s", hostIdentify, err, rid)
		} else {
			blog.InfoJSON("identifier: pushed event inst %s, rid: %s", hostIdentify, rid)
		}

	}
	return nil
}

func getInt(data map[string]interface{}, key string) (int64, error) {
	i, err := util.GetInt64ByInterface(data[key])
	if err != nil {
		blog.ErrorJSON("identifier: getInt error: %s, key: %s, value: %s", err, key, data[key])
	}
	return i, err
}

func getString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%s", value)
}

func (ih *IdentifierHandler) findHost(objType string, instID int64) (hostIDs []int64, err error) {
	rid := util.ExtractRequestIDFromContext(ih.ctx)
	relations := make([]metadata.ModuleHost, 0)
	filter := map[string]interface{}{
		common.GetInstIDField(objType): instID,
	}

	if objType == common.BKInnerObjIDPlat {
		if err = ih.db.Table(common.BKTableNameBaseHost).Find(filter).Fields(common.BKHostIDField).All(ih.ctx, &relations); err != nil {
			return nil, err
		}
	} else if objType == common.BKInnerObjIDProc {
		// 根据进程获取主机信息
		serviceInstRelationArr := make([]metadata.ProcessInstanceRelation, 0)
		if err = ih.db.Table(common.BKTableNameProcessInstanceRelation).Find(filter).All(context.Background(), &serviceInstRelationArr); err != nil {
			blog.ErrorJSON("find table(%s) data error. err:%s, filter:%s, rid: %s", err.Error(), filter, rid)
			return nil, err
		}
		for _, item := range serviceInstRelationArr {
			hostIDs = append(hostIDs, item.HostID)
		}
		return hostIDs, nil

	} else {
		if err = ih.db.Table(common.BKTableNameModuleHostConfig).Find(filter).Fields(common.BKHostIDField).All(ih.ctx, &relations); err != nil {
			return nil, err
		}
	}

	for index := range relations {
		hostIDs = append(hostIDs, relations[index].HostID)
	}
	return hostIDs, nil
}

type Inst struct {
	objType string
	instID  int64
	data    map[string]interface{}
	ident   *metadata.HostIdentifier
}

func (i *Inst) set(key string, value interface{}) error {
	i.data[key] = value
	var err error
	if i.objType == common.BKInnerObjIDHost {
		switch key {
		case common.BKHostNameField:
			i.ident.HostName = getString(value)
		case common.BKCloudIDField:
			i.ident.CloudID, err = util.GetInt64ByInterface(value)
		case common.BKHostInnerIPField:
			i.ident.InnerIP = getString(value)
		case common.BKHostOuterIPField:
			i.ident.OuterIP = getString(value)
		case common.BKOSTypeField:
			i.ident.OSType = getString(value)
		case common.BKOSNameField:
			i.ident.OSName = getString(value)
		case "bk_mem":
			i.ident.Memory, err = util.GetInt64ByInterface(value)
		case "bk_cpu":
			i.ident.CPU, err = util.GetInt64ByInterface(value)
		case "bk_disk":
			i.ident.Disk, err = util.GetInt64ByInterface(value)
		}
		if nil != err {
			blog.Errorf("key %s	convert error %s", key, err.Error())
		}
	}
	return err
}

func (i *Inst) saveCache(cache *redis.Client) error {
	out, err := json.Marshal(i.data)
	if err != nil {
		return err
	}
	return cache.Set(getInstCacheKey(i.objType, i.instID), string(out), time.Minute*20).Err()
}

func NewHostIdentifier(m map[string]interface{}) (*metadata.HostIdentifier, error) {
	var err error
	ident := metadata.HostIdentifier{}
	ident.HostName = getString(m[common.BKHostNameField])
	ident.CloudID, err = util.GetInt64ByInterface(m[common.BKCloudIDField])
	if nil != err {
		blog.Errorf("%s is not integer, %+v", common.BKCloudIDField, m)
		return nil, err
	}
	ident.InnerIP = getString(m[common.BKHostInnerIPField])
	ident.OuterIP = getString(m[common.BKHostOuterIPField])
	ident.OSType = getString(m[common.BKOSTypeField])
	ident.OSName = getString(m[common.BKOSNameField])
	ident.HostID, err = util.GetInt64ByInterface(m[common.BKHostIDField])
	if nil != err {
		blog.Errorf("%s is not integer, %+v", common.BKHostIDField, m)
		return nil, err
	}
	if m["bk_mem"] != nil {
		ident.Memory, err = util.GetInt64ByInterface(m["bk_mem"])
		if nil != err {
			blog.Errorf("bk_mem is not integer, %+v", m)
			return nil, err
		}
	}
	if m["bk_cpu"] != nil {
		ident.CPU, err = util.GetInt64ByInterface(m["bk_cpu"])
		if nil != err {
			blog.Errorf("bk_cpu is not integer, %+v", m)
			return nil, err
		}
	}
	if m["bk_disk"] != nil {
		ident.Disk, err = util.GetInt64ByInterface(m["bk_disk"])
		if nil != err {
			blog.Errorf("bk_disk is not integer, %+v", m)
			return nil, err
		}
	}
	ident.HostIdentModule = map[string]*metadata.HostIdentModule{}
	return &ident, nil
}

func getCache(ctx context.Context, cache *redis.Client, db dal.RDB, objType string, instID int64, fromdb bool) (*Inst, error) {
	var err error
	ret := cache.Get(getInstCacheKey(objType, instID)).Val()
	inst := Inst{objType: objType, instID: instID, ident: &metadata.HostIdentifier{}, data: map[string]interface{}{}}
	if "" == ret || nilStr == ret || fromdb {
		blog.Warnf("object %s inst %d not in cache, fetch it from db", objType, instID)
		getObjCondition := map[string]interface{}{
			common.GetInstIDField(objType): instID,
		}
		if err = db.Table(common.GetInstTableName(objType)).Find(getObjCondition).One(ctx, &inst.data); err != nil {
			blog.ErrorJSON("find object %s inst %s error: %s", objType, instID, err)
			return nil, err
		}
		if common.BKInnerObjIDHost == objType {
			inst.ident, err = NewHostIdentifier(inst.data)
			if err != nil {
				blog.ErrorJSON("NewHostIdentifier using inst data %s error: %s", inst.data, err)
				return nil, err
			}

			// 1. fill modules
			relations := make([]metadata.ModuleHost, 0)
			moduleHostCond := map[string]interface{}{
				common.GetInstIDField(objType): instID,
			}
			if err = db.Table(common.BKTableNameModuleHostConfig).Find(moduleHostCond).All(ctx, &relations); err != nil {
				blog.ErrorJSON("find module host relation hostID %s error: %s", instID, err)
				return nil, err
			}
			for _, relate := range relations {
				inst.ident.HostIdentModule[strconv.FormatInt(relate.ModuleID, 10)] = &metadata.HostIdentModule{
					SetID:    relate.SetID,
					ModuleID: relate.ModuleID,
					BizID:    relate.AppID,
				}
			}

			// 2. fill process
			hostProcessMap, err := getHostIdentifierProcInfo(ctx, db, []int64{instID})
			if err != nil {
				blog.ErrorJSON("find host process error: %s", err)
				return nil, err
			}
			inst.ident.Process = hostProcessMap[instID]

			// 3. fill other instances' detail
			inst.ident, err = fillIdentifier(inst.ident, ctx, cache, db)
			if err != nil {
				blog.ErrorJSON("fillIdentifier for ident: %s, error: %s", inst.ident, err)
				return nil, err
			}
			inst.data["associations"] = inst.ident.HostIdentModule
			inst.data["process"] = inst.ident.Process
		}
		if err := inst.saveCache(cache); err != nil {
			blog.ErrorJSON("saveCache inst data %s failed, err: %s", inst.data, err)
		}
	} else {
		err := json.Unmarshal([]byte(ret), &inst.data)
		if nil != err {
			blog.Errorf("unmarshal error %v, raw is %s", err, ret)
			return nil, err
		}
		if objType == common.BKInnerObjIDHost {
			err = json.Unmarshal([]byte(ret), inst.ident)
			if err != nil {
				blog.Errorf("unmarshal error %s, raw is %s", err.Error(), ret)
				return nil, err
			}
		}
	}

	if len(inst.data) <= 0 {
		return nil, nil
	}

	return &inst, nil
}

// getHostIdentifierProcInfo 根据主机ID生成主机身份
func getHostIdentifierProcInfo(ctx context.Context, db dal.RDB, hostIDs []int64) (map[int64][]metadata.HostIdentProcess, error) {
	relationFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDs,
		},
	}
	relations := make([]metadata.ProcessInstanceRelation, 0)

	// query process id with host id
	err := db.Table(common.BKTableNameProcessInstanceRelation).Find(relationFilter).All(ctx, &relations)
	if err != nil {
		blog.ErrorJSON("getHostIdentifierProcInfo query table %s err. cond:%s", common.BKTableNameProcessInstanceRelation, relationFilter)
		return nil, err
	}

	blog.V(5).Infof("getHostIdentifierProcInfo query host and process relation. hostID:%#v, relation:%#v", hostIDs, relations)

	procIDs := make([]int64, 0)
	serviceInstIDs := make([]int64, 0)
	// 进程与服务实例的关系
	procServiceInstMap := make(map[int64][]int64, 0)
	for _, relation := range relations {
		procIDs = append(procIDs, relation.ProcessID)
		serviceInstIDs = append(serviceInstIDs, relation.ServiceInstanceID)
		procServiceInstMap[relation.ProcessID] = append(procServiceInstMap[relation.ProcessID], relation.ServiceInstanceID)
	}

	serviceInstInfos := make([]metadata.ServiceInstance, 0)
	serviceInstFilter := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: serviceInstIDs,
		},
	}
	err = db.Table(common.BKTableNameServiceInstance).Find(serviceInstFilter).All(ctx, &serviceInstInfos)
	if err != nil {
		blog.ErrorJSON("getHostIdentifierProcInfo query table %s err. cond:%s", common.BKTableNameServiceInstance, serviceInstFilter)
		return nil, err
	}
	blog.V(5).Infof("getHostIdentifierProcInfo query service instance info. service instance id:%#v, info:%#v", serviceInstIDs, serviceInstInfos)
	// 服务实例与模块的关系
	serviceInstModuleRelation := make(map[int64][]int64, 0)
	for _, serviceInstInfo := range serviceInstInfos {
		serviceInstModuleRelation[serviceInstInfo.ID] = append(serviceInstModuleRelation[serviceInstInfo.ID], serviceInstInfo.ModuleID)
	}

	procModuleRelation := make(map[int64][]int64, 0)
	for procID, serviceInstIDs := range procServiceInstMap {
		for _, serviceInstID := range serviceInstIDs {
			procModuleRelation[procID] = append(procModuleRelation[procID], serviceInstModuleRelation[serviceInstID]...)
		}
	}

	procs := make(map[int64]metadata.HostIdentProcess, 0)
	for _, procID := range procIDs {
		procInfo := metadata.HostIdentProcess{
			ProcessID:   procID,
			BindModules: procModuleRelation[procID],
		}
		procs[procInfo.ProcessID] = procInfo
	}

	hostProcRelation := make(map[int64][]metadata.HostIdentProcess, 0)
	// 主机和进程之间的关系,生成主机与进程的关系
	for _, relation := range relations {
		if procInfo, ok := procs[relation.ProcessID]; ok {
			hostProcRelation[relation.HostID] = append(hostProcRelation[relation.HostID], procInfo)
		}
	}
	return hostProcRelation, nil
}

func (ih *IdentifierHandler) Run() error {
	blog.Infof("identifier: handle identifiers started")
	go func() {
		ih.fetchHostCache()
	}()
	go func() {
		if err := ih.handleEventLoop(); err != nil {
			blog.Errorf("handleInstLoop failed, err: %+v", err)
		}
	}()
	select {}
}

func (ih *IdentifierHandler) handleEventLoop() error {
	rid := util.ExtractRequestIDFromContext(ih.ctx)
	defer func() {
		procErr := recover()
		if procErr != nil {
			blog.Errorf("identifier: handleEventLoop panic: %v, stack:\n%s, rid: %s", procErr, debug.Stack(), rid)
		}
		// keep handleInstLoop run forever
		go func() {
			if err := ih.handleEventLoop(); err != nil {
				blog.Errorf("handleEventLoop failed, err: %+v", err)
			}
		}()
	}()
	for {
		event := ih.popEvent()
		if nil == event {
			time.Sleep(time.Second * 2)
			continue
		}
		ih.handleEvent(event)
	}
}

const nilStr = "nil"

func (ih *IdentifierHandler) popEvent() *metadata.EventInstCtx {
	rid := util.ExtractRequestIDFromContext(ih.ctx)

	eventStrs := ih.cache.BRPop(time.Second*60, types.EventCacheEventQueueDuplicateKey).Val()

	if len(eventStrs) == 0 || eventStrs[1] == nilStr || len(eventStrs[1]) == 0 {
		return nil
	}

	// eventStrs format is []string{key, event}
	eventStr := eventStrs[1]
	event := metadata.EventInst{}
	if err := json.Unmarshal([]byte(eventStr), &event); err != nil {
		blog.Errorf("identifier: event distribute fail, unmarshal error: %+v, data=[%s], rid: %s", err, eventStr, rid)
		return nil
	}

	return &metadata.EventInstCtx{EventInst: event, Raw: eventStr}
}

func (ih *IdentifierHandler) fetchHostCache() {
	rid := util.ExtractRequestIDFromContext(ih.ctx)

	objs := []string{common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule, common.BKInnerObjIDPlat, common.BKInnerObjIDProc}
	asstArr := make([]metadata.Association, 0)
	cond := condition.CreateCondition().Field(common.AssociationKindIDField).Eq(common.AssociationKindMainline)
	err := ih.db.Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).All(ih.ctx, &asstArr)
	if err != nil {
		blog.ErrorJSON("[identifier][fetchHostCache] get cc_ObjAsst error: %s, condition:%s, rid: %s", err, cond.ToMapStr(), rid)
		return
	}
	for _, asst := range asstArr {
		objs = append(objs, asst.ObjectID)
	}
	objs = util.StrArrayUnique(objs)
	for _, objID := range objs {
		caches := make([]map[string]interface{}, 0)
		if err := ih.db.Table(common.GetInstTableName(objID)).Find(map[string]interface{}{}).All(ih.ctx, &caches); err != nil {
			blog.Errorf("set cache for objID %s error %v, rid: %s", objID, err, rid)
			continue
		}

		for _, cache := range caches {
			out, _ := json.Marshal(cache)
			instID, err := getInt(cache, common.GetInstIDField(objID))
			if err != nil {
				blog.ErrorJSON("set cache key: %s, value: %s, error %s, rid: %s", getInstCacheKey(objID, instID), cache, err, rid)
				continue
			}
			if err := ih.cache.Set(getInstCacheKey(objID, instID), string(out), 0).Err(); err != nil {
				blog.ErrorJSON("set cache key: %s, value: %s, error %s, rid: %s", getInstCacheKey(objID, instID), cache, err, rid)
				continue
			}
		}

		blog.Infof("identifier: fetched %d %s, rid: %s", len(caches), objID, rid)
	}

	relations := make([]metadata.ModuleHost, 0)
	hosts := make([]*metadata.HostIdentifier, 0)
	hostIDs := make([]int64, 0)

	err = ih.db.Table(common.BKTableNameModuleHostConfig).Find(map[string]interface{}{}).All(ih.ctx, &relations)
	if err != nil {
		blog.Errorf("[identifier][fetchHostCache] get cc_ModuleHostConfig error: %v, rid: %s", err, rid)
		return
	}
	err = ih.db.Table(common.BKTableNameBaseHost).Find(map[string]interface{}{}).All(ih.ctx, &hosts)
	if err != nil {
		blog.Errorf("[identifier][fetchHostCache] get cc_HostBase error: %v, rid: %s", err, rid)
		return
	}

	for _, ident := range hosts {
		hostIDs = append(hostIDs, ident.HostID)
	}
	hostProcessMap, err := getHostIdentifierProcInfo(ih.ctx, ih.db, hostIDs)
	if err != nil {
		blog.ErrorJSON("find host process error: %s", err)
		return
	}

	relationMap := map[int64][]metadata.ModuleHost{}
	for _, relate := range relations {
		relationMap[relate.HostID] = append(relationMap[relate.HostID], relate)
	}

	for _, ident := range hosts {
		ident.HostIdentModule = map[string]*metadata.HostIdentModule{}
		for _, relate := range relationMap[ident.HostID] {
			ident.HostIdentModule[strconv.FormatInt(relate.ModuleID, 10)] = &metadata.HostIdentModule{
				SetID:    relate.SetID,
				ModuleID: relate.ModuleID,
				BizID:    relate.AppID,
			}
		}
		ident.Process = hostProcessMap[ident.HostID]

		ident, err = fillIdentifier(ident, ih.ctx, ih.cache, ih.db)
		if err != nil {
			blog.ErrorJSON("fillIdentifier for ident: %s, error: %s", ident, err)
			continue
		}

		if err := ih.cache.Set(getInstCacheKey(common.BKInnerObjIDHost, ident.HostID), ident, 0).Err(); err != nil {
			blog.Errorf("set cache error %s, rid: %s", err.Error(), rid)
			continue
		}
	}
	blog.Infof("identifier: fetched %d hosts", len(hosts))
}

func hasChanged(curData, preData map[string]interface{}, fields ...string) (isDifferent bool) {
	for _, field := range fields {
		if curData[field] != preData[field] {
			return true
		}
	}
	return false
}

type IdentifierHandler struct {
	cache *redis.Client
	db    dal.RDB
	ctx   context.Context
}

func NewIdentifierHandler(ctx context.Context, cache *redis.Client, db dal.RDB) *IdentifierHandler {
	return &IdentifierHandler{ctx: ctx, cache: cache, db: db}
}

func getInstCacheKey(objType string, instID int64) string {
	return types.EventCacheIdentInstPrefix + objType + "_" + strconv.FormatInt(instID, 10)
}
