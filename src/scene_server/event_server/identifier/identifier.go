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
	"runtime/debug"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/types"
	"configcenter/src/storage/dal"

	"gopkg.in/redis.v5"
)

var delayTime = time.Second * 30

var hostIndentDiffFields = map[string][]string{
	common.BKInnerObjIDApp: {
		common.BKAppNameField,
	},
	common.BKInnerObjIDSet: {
		common.BKSetNameField,
		"bk_service_status",
		"bk_set_env",
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
		"bk_start_param_regex",
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
		event.ObjType == metadata.EventObjTypeModuleTransfer {
		ih.handleModuleTransfer(event)
	} else if event.EventType == metadata.EventTypeRelation &&
		event.ObjType == metadata.EventObjTypeProcModule {
		ih.handleBindProcess(event)
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

		instID := getInt(curData, instIDField)
		if 0 == instID {
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
			inst.set(field, curData[field])
		}
		err = inst.saveCache(ih.cache)
		if err != nil {
			blog.Errorf("identifier: SaveCache error %+v, rid: %s", err, rid)
			continue
		}

		if common.BKInnerObjIDHost == event.ObjType {
			hostIdentify.ID = ih.cache.Incr(types.EventCacheEventIDKey).Val()
			d := metadata.EventData{CurData: inst.ident.fillIdentifier(ih.ctx, ih.cache, ih.db), PreData: preIdentifier}
			hostIdentify.Data = append(hostIdentify.Data, d)

			ih.cache.LPush(types.EventCacheEventQueueKey, &hostIdentify)
			blog.InfoJSON("identifier: pushed event inst %s, rid: %s", hostIdentify, rid)
		} else {
			if err := ih.handleRelatedInst(hostIdentify, event.ObjType, instID, false); err != nil {
				blog.Warnf("handleRelatedInst failed objType: %s, inst: %d, error: %v, rid: %s", event.ObjType, instID, err, rid)
			}
		}
	}
}

func (ih *IdentifierHandler) handleModuleTransfer(e *metadata.EventInstCtx) {
	rid := util.ExtractRequestIDFromContext(ih.ctx)
	blog.InfoJSON("identifier: handle inst %s, rid: %s", e, rid)

	hostIdentify := e.EventInst
	hostIdentify.Data = nil
	hostIdentify.EventType = metadata.EventTypeRelation
	hostIdentify.ObjType = objectTypeIdentifier
	hostIdentify.Action = metadata.EventActionUpdate

	go func() {
		time.Sleep(delayTime)
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

			instID := getInt(curData, common.BKHostIDField)
			if 0 == instID {
				blog.Errorf("identifier: convert instID failed the raw is %+v, rid: %s", curData[common.BKHostIDField], rid)
				continue
			}

			inst, err := getCache(ih.ctx, ih.cache, ih.db, common.BKInnerObjIDHost, instID, true)
			if err != nil {
				blog.Errorf("identifier: getCache error %+v, rid: %s", err, rid)
				continue
			}
			if nil == inst {
				continue
			}

			if err := inst.saveCache(ih.cache); err != nil {
				blog.Errorf("saveCache failed, err: %+v", err)
			}
			d := metadata.EventData{CurData: inst.ident.fillIdentifier(ih.ctx, ih.cache, ih.db)}
			hostIdentify.Data = append(hostIdentify.Data, d)
		}
		hostIdentify.ID = ih.cache.Incr(types.EventCacheEventIDKey).Val()
		ih.cache.LPush(types.EventCacheEventQueueKey, &hostIdentify)
		blog.InfoJSON("identifier: pushed event inst %s, rid: %s", hostIdentify, rid)
	}()
}

const objectTypeIdentifier = "hostidentifier"

func (ih *IdentifierHandler) handleBindProcess(e *metadata.EventInstCtx) {
	rid := util.ExtractRequestIDFromContext(ih.ctx)
	blog.InfoJSON("identifier: handle inst %s, rid: %s", e, rid)

	hostIdentify := e.EventInst
	hostIdentify.Data = nil
	hostIdentify.EventType = metadata.EventTypeRelation
	hostIdentify.ObjType = objectTypeIdentifier
	hostIdentify.Action = metadata.EventActionUpdate

	go func() {
		time.Sleep(delayTime)
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

			instID := getInt(curData, common.BKProcIDField)
			if 0 == instID {
				blog.Errorf("identifier: convert instID failed the raw is %+v, rid: %s", curData[common.BKProcIDField], rid)
				continue
			}

			modules := make([]metadata.ModuleInst, 0)
			filter := map[string]interface{}{
				common.BKSupplierIDField: curData[common.BKSupplierIDField],
				common.BKAppIDField:      curData[common.BKAppIDField],
				common.BKModuleNameField: curData[common.BKModuleNameField],
			}
			if err := ih.db.Table(common.BKTableNameBaseModule).Find(filter).All(ih.ctx, &modules); err != nil {
				continue
			}

			for _, module := range modules {
				if err := ih.handleRelatedInst(hostIdentify, common.BKInnerObjIDModule, module.ModuleID, true); err != nil {
					blog.Warnf("handleRelatedInst failed objtype: %s, inst: %d, error: %v, rid: %s", e.ObjType, instID, err, rid)
				}
			}

		}
	}()
}

func (ih *IdentifierHandler) handleRelatedInst(hostIdentify metadata.EventInst, objType string, instID int64, formdb bool) error {
	rid := util.ExtractRequestIDFromContext(ih.ctx)
	hosIDs, err := ih.findHost(objType, instID)
	if err != nil {
		blog.Warnf("identifier: find host failure: %v, rid: %s", err, rid)
		return err
	}
	blog.V(3).Infof("identifier: handleRelatedInst by objType %s, instID %d,  hostIDs: %v, fromdb: %v, rid: %s", objType, instID, hosIDs, formdb, rid)
	total := len(hosIDs)
	index := 0
	leftIndex := 0

	for leftIndex < total {
		leftIndex = index + 256
		if leftIndex > total {
			leftIndex = total
		}
		hostIdentify.Data = nil

		if formdb {
			for _, hostID := range hosIDs[index:leftIndex] {
				inst, getCacheErr := getCache(ih.ctx, ih.cache, ih.db, common.BKInnerObjIDHost, hostID, true)
				if getCacheErr != nil {
					blog.Errorf("identifier: getCache error %+v, rid: %s", getCacheErr, rid)
					continue
				}
				if nil == inst {
					continue
				}
				if err := inst.saveCache(ih.cache); err != nil {
					blog.Errorf("saveCache failed, err: %+v", err)
				}
				d := metadata.EventData{CurData: inst.ident.fillIdentifier(ih.ctx, ih.cache, ih.db)}
				hostIdentify.Data = append(hostIdentify.Data, d)
			}
		} else {
			hostIDKeys := make([]string, 0)
			for _, hostID := range hosIDs[index:leftIndex] {
				hostIDKeys = append(hostIDKeys, getInstCacheKey(common.BKInnerObjIDHost, hostID))
			}
			idens, err := ih.cache.MGet(hostIDKeys...).Result()
			if err != nil {
				blog.Errorf("identifier: ih.cache.MGet by %v,%v. we will try to fetch it from db instead, rid: %s", hostIDKeys, err, rid)
				idens = make([]interface{}, len(hostIDKeys))
				for index := range idens {
					// simulate that redis returns all nil
					idens[index] = nilStr
				}
			}
			for identIndex := range idens {
				iden := HostIdentifier{}
				if err = json.Unmarshal([]byte(getString(idens[identIndex])), &iden); err != nil {
					blog.Warnf("identifier: unmarshal error %s, rid: %s", err.Error(), rid)
					hostID := hosIDs[index:leftIndex][identIndex]
					inst, err := getCache(ih.ctx, ih.cache, ih.db, common.BKInnerObjIDHost, hostID, true)
					if err != nil {
						blog.Errorf("identifier: getCache error %+v, rid: %s", err, rid)
						continue
					}
					if nil == inst {
						continue
					}
					if err := inst.saveCache(ih.cache); err != nil {
						blog.Errorf("saveCache failed, err: %+v", err)
					}
					d := metadata.EventData{CurData: inst.ident.fillIdentifier(ih.ctx, ih.cache, ih.db)}
					hostIdentify.Data = append(hostIdentify.Data, d)
					continue
				}
				d := metadata.EventData{CurData: iden.fillIdentifier(ih.ctx, ih.cache, ih.db)}
				hostIdentify.Data = append(hostIdentify.Data, d)
			}
		}
		index += 256

		hostIdentify.ID = ih.cache.Incr(types.EventCacheEventIDKey).Val()
		if err = ih.cache.LPush(types.EventCacheEventQueueKey, &hostIdentify).Err(); err != nil {
			blog.Warnf("identifier: push event inst %v failure %v, rid: %s", hostIdentify, err, rid)
		} else {
			blog.InfoJSON("identifier: pushed event inst %s, rid: %s", hostIdentify, rid)
		}

	}
	return nil
}

func NewModule(m map[string]interface{}) *Module {
	belong := Module{}
	belong.BizID = getInt(m, common.BKAppIDField)
	belong.SetID = getInt(m, common.BKSetIDField)
	belong.ModuleID = getInt(m, common.BKModuleIDField)
	return &belong
}

func getInt(data map[string]interface{}, key string) int64 {
	i, err := util.GetInt64ByInterface(data[key])
	if err != nil {
		blog.Errorf("identifier: getInt error: %+v", err)
	}
	return i
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
	ident   *HostIdentifier
}

func (i *Inst) set(key string, value interface{}) {
	i.data[key] = value
	var err error
	if i.objType == common.BKInnerObjIDHost {
		switch key {
		case "bk_host_name":
			i.ident.HostName = getString(value)
		case "bk_cloud_id":
			i.ident.CloudID, err = util.GetInt64ByInterface(value)
		case "bk_host_innerip":
			i.ident.InnerIP = getString(value)
		case "bk_host_outerip":
			i.ident.OuterIP = getString(value)
		case "bk_os_type":
			i.ident.OSType = getString(value)
		case "bk_os_name":
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
}

func (i *Inst) saveCache(cache *redis.Client) error {
	out, err := json.Marshal(i.data)
	if err != nil {
		return err
	}
	return cache.Set(getInstCacheKey(i.objType, i.instID), string(out), time.Minute*20).Err()
}

func NewHostIdentifier(m map[string]interface{}) *HostIdentifier {
	var err error
	ident := HostIdentifier{}
	ident.HostName = getString(m["bk_host_name"])
	ident.CloudID, err = util.GetInt64ByInterface(m["bk_cloud_id"])
	if nil != err {
		blog.Errorf("%s is not integer, %+v", "bk_cloud_id", m)
	}
	ident.InnerIP = getString(m["bk_host_innerip"])
	ident.OuterIP = getString(m["bk_host_outerip"])
	ident.OSType = getString(m["bk_os_type"])
	ident.OSName = getString(m["bk_os_name"])
	ident.HostID, err = util.GetInt64ByInterface(m[common.BKHostIDField])
	if nil != err {
		blog.Errorf("%s is not integer, %+v ", "bk_host_id", m)
	}
	if m["bk_mem"] != nil {
		ident.Memory, err = util.GetInt64ByInterface(m["bk_mem"])
		if nil != err {
			blog.Warnf("%s is not integer, %+v ", "bk_mem", m)
		}
	}
	if m["bk_cpu"] != nil {
		ident.CPU, err = util.GetInt64ByInterface(m["bk_cpu"])
		if nil != err {
			blog.Warnf("%s is not integer, %+v ", "bk_cpu", m)
		}
	}
	if m["bk_disk"] != nil {
		ident.Disk, err = util.GetInt64ByInterface(m["bk_disk"])
		if nil != err {
			blog.Warnf("%s is not integer, %+v ", "bk_disk", m)
		}
	}
	ident.Module = map[string]*Module{}
	return &ident
}

func getCache(ctx context.Context, cache *redis.Client, db dal.RDB, objType string, instID int64, fromdb bool) (*Inst, error) {
	var err error
	ret := cache.Get(getInstCacheKey(objType, instID)).Val()
	inst := Inst{objType: objType, instID: instID, ident: &HostIdentifier{}, data: map[string]interface{}{}}
	if "" == ret || nilStr == ret || fromdb {
		blog.Infof("objType %s, instID %d not in cache, fetch it from db", objType, instID)
		getObjCondition := map[string]interface{}{
			common.GetInstIDField(objType): instID,
		}
		if err = db.Table(common.GetInstTableName(objType)).Find(getObjCondition).One(ctx, &inst.data); err != nil {
			return nil, err
		}
		if common.BKInnerObjIDHost == objType {
			inst.ident = NewHostIdentifier(inst.data)
			hostModuleIDs := make([]int64, 0)

			// 1. fill modules
			relations := make([]metadata.ModuleHost, 0)
			moduleHostCond := map[string]interface{}{
				common.GetInstIDField(objType): instID,
			}
			if err = db.Table(common.BKTableNameModuleHostConfig).Find(moduleHostCond).All(ctx, &relations); err != nil {
				return nil, err
			}
			for _, relate := range relations {
				hostModuleIDs = append(hostModuleIDs, relate.ModuleID)
				inst.ident.Module[strconv.FormatInt(relate.ModuleID, 10)] = &Module{
					SetID:    relate.SetID,
					ModuleID: relate.ModuleID,
					BizID:    relate.AppID,
				}
			}
			inst.data["associations"] = inst.ident.Module

			// 2. fill process
			hostProcessMap, err := getHostIdentifierProcInfo(ctx, db, []int64{instID})
			if err != nil {
				blog.InfoJSON("find host")
				return nil, err
			}
			inst.ident.Process = hostProcessMap[instID]
			inst.data["process"] = inst.ident.Process
		}
		if err := inst.saveCache(cache); err != nil {
			blog.Errorf("saveCache failed, err: %+v", err)
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
func getHostIdentifierProcInfo(ctx context.Context, db dal.RDB, hostIDs []int64) (map[int64][]Process, error) {
	relationFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDs,
		},
	}
	relations := make([]metadata.ProcessInstanceRelation, 0)

	// query process id with host id
	err := db.Table(common.BKTableNameProcessInstanceRelation).Find(relationFilter).All(ctx, &relations)
	if err != nil {
		blog.ErrorJSON("findHostServiceInst query table %s err. cond:%s", common.BKTableNameProcessInstanceRelation, relationFilter)
		return nil, err
	}

	blog.V(5).Infof("findHostServiceInst query host and process relation. hostID:%#v, relation:%#v", hostIDs, relations)

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
		blog.ErrorJSON("findHostServiceInst query table %s err. cond:%s", common.BKTableNameBaseProcess, serviceInstFilter)
		return nil, err
	}
	blog.V(5).Infof("findHostServiceInst query service instance info. service instance id:%#v, info:%#v", serviceInstIDs, serviceInstInfos)
	// 服务实例与模块的关系
	serviceInstModuleRelation := make(map[int64][]int64, 0)
	for _, serviceInstInfo := range serviceInstInfos {
		serviceInstModuleRelation[serviceInstInfo.ID] = append(serviceInstModuleRelation[serviceInstInfo.ID], serviceInstInfo.ModuleID)
	}

	procModuleRelation := make(map[int64][]int64, 0)
	for procID, serviceInstIDs := range procServiceInstMap {
		for _, serviceInstID := range serviceInstIDs {
			for _, moduleID := range serviceInstModuleRelation[serviceInstID] {
				procModuleRelation[procID] = append(procModuleRelation[procID], moduleID)
			}
		}
	}

	procInfos := make([]Process, 0)
	// query process info with process id
	processFilter := map[string]interface{}{
		common.BKProcIDField: map[string]interface{}{
			common.BKDBIN: procIDs,
		},
	}
	err = db.Table(common.BKTableNameBaseProcess).Find(processFilter).All(ctx, &procInfos)
	if err != nil {
		blog.ErrorJSON("findHostServiceInst query table %s err. cond:%s", common.BKTableNameBaseProcess, processFilter)
		return nil, err
	}

	blog.V(5).Infof("findHostServiceInst query process info. procIDs:%#v, info:%#v", procIDs, procInfos)

	procs := make(map[int64]Process, 0)
	for _, procInfo := range procInfos {
		if moduleIDs, ok := procModuleRelation[procInfo.ProcessID]; ok {
			procInfo.BindModules = moduleIDs
		}
		procs[procInfo.ProcessID] = procInfo
	}

	hostProcRelation := make(map[int64][]Process, 0)
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
		for range time.Tick(time.Minute * 10) {
			ih.fetchHostCache()
		}
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
		blog.Errorf("identifier: event distribute fail, unmarshal error: %+v, date=[%s], rid: %s", err, eventStr, rid)
		return nil
	}

	return &metadata.EventInstCtx{EventInst: event, Raw: eventStr}
}

func (ih *IdentifierHandler) fetchHostCache() {
	rid := util.ExtractRequestIDFromContext(ih.ctx)

	relations := make([]metadata.ModuleHost, 0)
	hosts := make([]*HostIdentifier, 0)
	modules := make([]metadata.ModuleInst, 0)
	proc2modules := make([]metadata.ProcessModule, 0)

	err := ih.db.Table(common.BKTableNameModuleHostConfig).Find(map[string]interface{}{}).All(ih.ctx, &relations)
	if err != nil {
		blog.Errorf("[identifier][fetchHostCache] get cc_ModuleHostConfig error: %v, rid: %s", err, rid)
		return
	}
	err = ih.db.Table(common.BKTableNameBaseHost).Find(map[string]interface{}{}).All(ih.ctx, &hosts)
	if err != nil {
		blog.Errorf("[identifier][fetchHostCache] get cc_HostBase error: %v, rid: %s", err, rid)
		return
	}
	err = ih.db.Table(common.BKTableNameProcModule).Find(map[string]interface{}{}).All(ih.ctx, &proc2modules)
	if err != nil {
		blog.Errorf("[identifier][fetchHostCache] get cc_Proc2Module error: %v, rid: %s", err, rid)
		return
	}
	err = ih.db.Table(common.BKTableNameBaseModule).Find(map[string]interface{}{}).All(ih.ctx, &modules)
	if err != nil {
		blog.Errorf("[identifier][fetchHostCache] get cc_ModuleBase error: %v, rid: %s", err, rid)
		return
	}

	relationMap := map[int64][]metadata.ModuleHost{}
	for _, relate := range relations {
		relationMap[relate.HostID] = append(relationMap[relate.HostID], relate)
	}

	proc2modulesMap := map[string][]int64{}
	bindModulesMap := map[int64][]int64{}
	for _, proc2module := range proc2modules {
		proc2modulesMap[proc2module.ModuleName] = append(proc2modulesMap[proc2module.ModuleName], proc2module.ProcessID)
		for _, module := range modules {
			if module.BizID == proc2module.AppID && module.ModuleName == proc2module.ModuleName {
				bindModulesMap[proc2module.ProcessID] = append(bindModulesMap[proc2module.ProcessID], module.ModuleID)
			}
		}
	}

	modulesMap := map[int64]metadata.ModuleInst{}
	for _, module := range modules {
		modulesMap[module.ModuleID] = module
	}

	for _, ident := range hosts {
		ident.Module = map[string]*Module{}
		hostProcs := map[int64]bool{}
		for _, relate := range relationMap[ident.HostID] {
			ident.Module[strconv.FormatInt(relate.ModuleID, 10)] = &Module{
				SetID:    relate.SetID,
				ModuleID: relate.ModuleID,
				BizID:    relate.AppID,
			}
			if module, ok := modulesMap[relate.ModuleID]; ok {
				for _, procID := range proc2modulesMap[module.ModuleName] {
					hostProcs[procID] = true
				}
			}
		}

		ident.Process = make([]Process, 0)
		for procID := range hostProcs {
			bindModules := bindModulesMap[procID]
			if len(bindModules) == 0 {
				bindModules = make([]int64, 0)
			}
			ident.Process = append(ident.Process, Process{ProcessID: procID, BindModules: bindModules})
		}

		if err := ih.cache.Set(getInstCacheKey(common.BKInnerObjIDHost, ident.HostID), ident, 0).Err(); err != nil {
			blog.Errorf("set cache error %s, rid: %s", err.Error(), rid)
		}
	}
	blog.Infof("identifier: fetched %d hosts", len(hosts))

	objs := []string{common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule, common.BKInnerObjIDPlat, common.BKInnerObjIDProc}
	for _, objID := range objs {
		caches := make([]map[string]interface{}, 0)
		if err := ih.db.Table(common.GetInstTableName(objID)).Find(map[string]interface{}{}).All(ih.ctx, &caches); err != nil {
			blog.Errorf("set cache for objID %s error %v, rid: %s", objID, err, rid)
		}

		for _, cache := range caches {
			out, _ := json.Marshal(cache)
			if err := ih.cache.Set(getInstCacheKey(objID, getInt(cache, common.GetInstIDField(objID))), string(out), 0).Err(); err != nil {
				blog.Errorf("set cache error %v, rid: %s", err, rid)
			}
		}

		blog.Infof("identifier: fetched %d %s, rid: %s", len(caches), objID, rid)
	}

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
