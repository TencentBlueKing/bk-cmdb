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
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/types"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
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
		common.BKFuncName,
		common.BKBindIP,
		common.BKProtocol,
		common.BKPort,
		common.BKStartParamRegex,
		common.BKProcBindInfo,
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
	/*
	   NOTE: Identifier event handle rules here,
	       For EventTypeInstData:
	           only handle events App: update, Set: update, Module: update, Process: update, Host: create, update
	       For EventTypeRelation:
	           only handle events ModuleHostRelation: create/update/delete ProcessInstanceRelation: create/update/delete
	*/

	if event.EventType == metadata.EventTypeInstData {
		diffFields, ok := hostIndentDiffFields[event.ObjType]
		if !ok {
			return
		}

		// host create and all object type update event would be handled.
		if event.Action == metadata.EventActionUpdate ||
			(event.ObjType == common.BKInnerObjIDHost && event.Action == metadata.EventActionCreate) {
			ih.handleInstFieldChange(event, diffFields)
		}

	} else if event.EventType == metadata.EventTypeRelation {
		// handle module and process relation events.
		if event.ObjType == metadata.EventObjTypeModuleTransfer ||
			event.ObjType == metadata.EventObjTypeProcModule {
			ih.handleHostRelationChange(event)
		}
	}
}

func (ih *IdentifierHandler) handleInstFieldChange(event *metadata.EventInstCtx, diffFields []string) error {
	rid := util.ExtractRequestIDFromContext(ih.ctx)
	blog.InfoJSON("identifier: handle inst %s, rid: %s", event, rid)

	hostIdentify := event.EventInst
	hostIdentify.Data = nil
	hostIdentify.EventType = metadata.EventTypeRelation
	hostIdentify.ObjType = objectTypeIdentifier
	hostIdentify.Action = metadata.EventActionUpdate

	if len(event.Data) == 0 {
		blog.Errorf("empty identifier event data, %+v", event)
		return errors.New("empty event data")
	}
	metaEventData := event.Data[0]
	curData, ok := metaEventData.CurData.(map[string]interface{})
	if !ok {
		return errors.New("invalid event current data")
	}

	updateFields := event.UpdateFields
	deletedFields := event.DeletedFields

	if !hasChanged(updateFields, deletedFields, diffFields...) {
		return nil
	}

	instIDField := common.GetInstIDField(event.ObjType)
	instID, err := getInt(curData, instIDField)
	if err != nil || 0 == instID {
		blog.Errorf("identifier: convert instID failed the raw is %+v, rid: %s", curData[instIDField], rid)
		return err
	}

	inst, err := getCache(ih.ctx, ih.cache, ih.clientSet, ih.db, event.ObjType, instID)
	if err != nil {
		blog.Errorf("identifier: getCache error %+v, rid: %s", err, rid)
		return err
	}
	if inst == nil {
		blog.Errorf("identifier: inst == nil, continue, rid:%s", rid)
		return err
	}

	if common.BKInnerObjIDHost == event.ObjType {
		// save previous host identifier data for later usage.
		preInst := inst.copy()

		for _, field := range diffFields {
			if err := inst.set(field, curData[field]); err != nil {
				blog.Errorf("key %s, value: %s, convert error %s", field, curData[field], err.Error())
				return err
			}
		}

		if err := inst.saveCache(ih.cache); err != nil {
			blog.ErrorJSON("saveCache inst data %s failed, err: %s, rid: %s", inst.data, err, rid)
			return err
		}

		hostIdentify.ID = ih.cache.Incr(ih.ctx, types.EventCacheEventIDKey).Val()
		d := metadata.EventData{CurData: inst.ident, PreData: preInst.ident}
		hostIdentify.Data = append(hostIdentify.Data, d)

		ih.cache.LPush(ih.ctx, types.EventCacheEventQueueKey, &hostIdentify)
		blog.InfoJSON("identifier: pushed event inst %s, rid: %s", hostIdentify, rid)
	} else {
		if err := ih.handleRelatedInst(hostIdentify, event.ObjType, instID, curData, diffFields); err != nil {
			blog.Warnf("handleRelatedInst failed objType: %s, inst: %d, error: %v, rid: %s", event.ObjType, instID, err, rid)
			return err
		}
	}

	return nil
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
		// NOTE: why 30 seconds?
		time.Sleep(delayTime)

		hostIDs := make([]int64, 0)
		var preData, curData map[string]interface{}

		for index := range e.Data {
			var hostID int64
			var ok bool
			var err error

			if e.Action == metadata.EventActionDelete {
				preData, ok = e.Data[index].PreData.(map[string]interface{})
				if !ok {
					continue
				}
				hostID, err = getInt(preData, common.BKHostIDField)
			}

			if e.Action == metadata.EventActionCreate || e.Action == metadata.EventActionUpdate {
				curData, ok = e.Data[index].CurData.(map[string]interface{})
				if !ok {
					continue
				}
				hostID, err = getInt(curData, common.BKHostIDField)
			}

			if err != nil || 0 == hostID {
				blog.Errorf("identifier: convert instID failed the raw is %+v, rid: %s", curData[common.BKHostIDField], rid)
				continue
			}
			hostIDs = append(hostIDs, hostID)
		}

		hostIDs = util.IntArrayUnique(hostIDs)

		handler := func(hostID int64, inst *Inst) error {
			switch e.ObjType {
			case metadata.EventObjTypeModuleTransfer:
				if len(preData) != 0 {
					deletedModuleID, err := getInt(preData, common.BKModuleIDField)
					if err != nil || 0 == deletedModuleID {
						blog.Errorf("identifier: convert moduleID failed the raw is %+v, rid: %s", preData[common.BKModuleIDField], rid)
						return err
					}
					delete(inst.ident.HostIdentModule, strconv.FormatInt(deletedModuleID, 10))
				}

				if len(curData) != 0 {
					createdModuleID, err := getInt(curData, common.BKModuleIDField)
					if err != nil || 0 == createdModuleID {
						blog.Errorf("identifier: convert moduleID failed the raw is %+v, rid: %s", curData[common.BKModuleIDField], rid)
						return err
					}
					_, exist := inst.ident.HostIdentModule[strconv.FormatInt(createdModuleID, 10)]
					if exist {
						return nil
					}

					bizID, err := getInt(curData, common.BKAppIDField)
					if err != nil || 0 == bizID {
						blog.Errorf("identifier: convert bizID failed the raw is %+v, rid: %s", curData[common.BKAppIDField], rid)
						return err
					}

					setID, err := getInt(curData, common.BKSetIDField)
					if err != nil || 0 == setID {
						blog.Errorf("identifier: convert setID failed the raw is %+v, rid: %s", curData[common.BKSetIDField], rid)
						return err
					}

					hostIdentModule := &metadata.HostIdentModule{
						BizID:    bizID,
						SetID:    setID,
						ModuleID: createdModuleID,
					}

					err = fillModule(inst.ident, hostIdentModule, nil, ih.ctx, ih.cache, ih.clientSet, ih.db)
					if err != nil {
						blog.ErrorJSON("identifier: fillModule error %s, hostIdentModule: %s", err, hostIdentModule)
						return err
					}
					inst.ident.HostIdentModule[strconv.FormatInt(createdModuleID, 10)] = hostIdentModule
				}

			case metadata.EventObjTypeProcModule:
				var deletedProcessID, createdProcessID int64 = -1, -1
				var err error
				needInsert := false

				if len(preData) != 0 {
					deletedProcessID, err = getInt(preData, common.BKProcessIDField)
					if err != nil || 0 == deletedProcessID {
						blog.Errorf("identifier: convert processID failed the raw is %+v, rid: %s", preData[common.BKProcessIDField], rid)
						return err
					}
				}

				if len(curData) != 0 {
					createdProcessID, err = getInt(curData, common.BKProcessIDField)
					if err != nil || 0 == createdProcessID {
						blog.Errorf("identifier: convert processID failed the raw is %+v, rid: %s", curData[common.BKProcessIDField], rid)
						return err
					}
					needInsert = true
				}

				delIndex := -1
				for index, process := range inst.ident.Process {
					if process.ProcessID == deletedProcessID {
						delIndex = index
					}
					if process.ProcessID == createdProcessID {
						needInsert = false
					}
				}
				if delIndex != -1 {
					inst.ident.Process = append(inst.ident.Process[:delIndex], inst.ident.Process[delIndex+1:]...)
				}

				if needInsert {
					serviceInstanceID, err := getInt(curData, common.BKServiceInstanceIDField)
					if err != nil || 0 == createdProcessID {
						blog.Errorf("identifier: convert serviceInstanceID failed the raw is %+v, rid: %s", curData[common.BKServiceInstanceIDField], rid)
						return err
					}

					procInfo, err := genProcInfo(ih.ctx, ih.db, createdProcessID, serviceInstanceID)
					if err != nil {
						blog.Errorf("genProcInfo by process %d and service instance %d failed, err: %s, rid: %s", createdProcessID, serviceInstanceID, rid)
						return err
					}

					err = fillProcess(&procInfo, ih.ctx, ih.cache, ih.clientSet, ih.db)
					if err != nil {
						blog.ErrorJSON("identifier: fillProcess error %s, process: %s", err, procInfo)
						return err
					}
					inst.ident.Process = append(inst.ident.Process, procInfo)
				}
			}

			return nil
		}

		if err := ih.handleHostBatch(hostIdentify, hostIDs, handler); err != nil {
			blog.Warnf("handleHostBatch failed , hostIDs: %v, error: %v, rid: %s", hostIDs, err, rid)
		}
	}()
}

func genProcInfo(ctx context.Context, db dal.RDB, procID int64, serviceInstanceID int64) (metadata.HostIdentProcess, error) {
	serviceInstanceInfo := metadata.ServiceInstance{}
	serviceInstFilter := map[string]interface{}{
		common.BKFieldID: serviceInstanceID,
	}
	err := db.Table(common.BKTableNameServiceInstance).Find(serviceInstFilter).Fields(common.BKFieldID, common.BKModuleIDField).One(ctx, &serviceInstanceInfo)
	if err != nil {
		blog.ErrorJSON("getHostIdentifierProcInfo query table %s err. cond:%s", common.BKTableNameServiceInstance, serviceInstFilter)
		return metadata.HostIdentProcess{}, err
	}
	blog.V(5).Infof("getHostIdentifierProcInfo query service instance info. service instance id:%d, info:%#v", serviceInstanceID, serviceInstanceInfo)
	procInfo := metadata.HostIdentProcess{
		ProcessID:   procID,
		BindModules: []int64{serviceInstanceInfo.ModuleID},
	}
	return procInfo, nil
}

const objectTypeIdentifier = "hostidentifier"

func (ih *IdentifierHandler) handleRelatedInst(hostIdentify metadata.EventInst, objType string, instID int64, curData map[string]interface{}, diffFields []string) error {
	rid := util.ExtractRequestIDFromContext(ih.ctx)
	hostIDs, err := ih.findHost(objType, instID)
	if err != nil {
		blog.Warnf("identifier: find host failure: %v, rid: %s", err, rid)
		return err
	}
	blog.V(3).Infof("identifier: handleRelatedInst by objType %s, instID %d,  hostIDs: %v, rid: %s", objType, instID, hostIDs, rid)
	handler := func(hostID int64, inst *Inst) error {
		switch objType {
		case common.BKInnerObjIDPlat:
			for _, field := range diffFields {
				if field == common.BKCloudNameField {
					inst.ident.CloudName = getString(curData[common.BKCloudNameField])
				}
			}
		case common.BKInnerObjIDApp:
			for _, hostIdentModule := range inst.ident.HostIdentModule {
				if hostIdentModule.BizID != instID {
					continue
				}
				for _, field := range diffFields {
					if field == common.BKAppNameField {
						hostIdentModule.BizName = getString(curData[field])
					}
				}
			}
		case common.BKInnerObjIDSet:
			for _, hostIdentModule := range inst.ident.HostIdentModule {
				if hostIdentModule.SetID != instID {
					continue
				}
				for _, field := range diffFields {
					switch field {
					case common.BKSetNameField:
						hostIdentModule.SetName = getString(curData[field])
					case common.BKSetStatusField:
						hostIdentModule.SetStatus = getString(curData[field])
					case common.BKSetEnvField:
						hostIdentModule.SetEnv = getString(curData[field])
					}
				}
			}
		case common.BKInnerObjIDModule:
			for _, hostIdentModule := range inst.ident.HostIdentModule {
				if hostIdentModule.ModuleID != instID {
					continue
				}
				for _, field := range diffFields {
					if field == common.BKModuleNameField {
						hostIdentModule.ModuleName = getString(curData[field])
					}
				}
			}
		case common.BKInnerObjIDProc:
			for index, process := range inst.ident.Process {
				if process.ProcessID != instID {
					continue
				}

				for _, field := range diffFields {
					switch field {
					case common.BKProcessNameField:
						process.ProcessName = getString(curData[field])
					case common.BKFuncName:
						process.FuncName = getString(curData[field])
					case common.BKStartParamRegex:
						process.StartParamRegex = getString(curData[field])
					}
				}

				ip, port, protocol, enable, bindInfoArr := getBindInfo(curData[common.BKProcBindInfo])
				process.BindIP = ip
				process.Port = port
				process.Protocol = protocol
				process.PortEnable = enable
				process.BindInfo = bindInfoArr
				inst.ident.Process[index] = process
			}
		}
		return nil
	}
	return ih.handleHostBatch(hostIdentify, hostIDs, handler)
}

func (ih *IdentifierHandler) handleHostBatch(hostIdentify metadata.EventInst, hostIDs []int64,
	handler func(hostID int64, inst *Inst) error) error {
	rid := util.ExtractRequestIDFromContext(ih.ctx)
	blog.V(3).Infof("identifier: handleHostBatch hostIDs: %v, rid: %s", hostIDs, rid)
	total := len(hostIDs)
	bufSize := 256
	index := 0
	leftIndex := 0

	for leftIndex < total {
		leftIndex = index + bufSize
		if leftIndex > total {
			leftIndex = total
		}
		hostIdentify.Data = nil

		for _, hostID := range hostIDs[index:leftIndex] {
			inst, err := getCache(ih.ctx, ih.cache, ih.clientSet, ih.db, common.BKInnerObjIDHost, hostID)
			if err != nil {
				blog.Errorf("identifier: getCache error %+v, rid: %s", err, rid)
				continue
			}
			if nil == inst {
				continue
			}

			preIdentifier := inst.copy()
			err = handler(hostID, inst)
			if err != nil {
				blog.Errorf("handleHostBatch handler error %+v, rid: %s", err, rid)
				continue
			}
			if err := inst.saveCache(ih.cache); err != nil {
				blog.ErrorJSON("saveCache inst data %s failed, err: %s, rid: %s", inst.data, err, rid)
			}
			d := metadata.EventData{CurData: inst.ident, PreData: preIdentifier.ident}
			hostIdentify.Data = append(hostIdentify.Data, d)
		}
		index += bufSize

		hostIdentify.ID = ih.cache.Incr(ih.ctx, types.EventCacheEventIDKey).Val()
		if err := ih.cache.LPush(ih.ctx, types.EventCacheEventQueueKey, &hostIdentify).Err(); err != nil {
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
			i.ident.InnerIP = metadata.StringArrayToString(getString(value))
		case common.BKHostOuterIPField:
			i.ident.OuterIP = metadata.StringArrayToString(getString(value))
		case common.BKOSTypeField:
			i.ident.OSType = getString(value)
		case common.BKOSNameField:
			i.ident.OSName = getString(value)
		case "bk_mem":
			if value != nil {
				i.ident.Memory, err = util.GetInt64ByInterface(value)
			}
		case "bk_cpu":
			if value != nil {
				i.ident.CPU, err = util.GetInt64ByInterface(value)
			}
		case "bk_disk":
			if value != nil {
				i.ident.Disk, err = util.GetInt64ByInterface(value)
			}
		}
		if nil != err {
			blog.Errorf("key %s	convert error %s", key, err.Error())
		}
	}
	return err
}

func (i *Inst) saveCache(cache redis.Client) error {
	out, err := json.Marshal(i.ident)
	if err != nil {
		return err
	}
	return cache.Set(context.Background(), getInstCacheKey(i.objType, i.instID), string(out), time.Minute*20).Err()
}

func (i *Inst) copy() *Inst {
	inst := &Inst{
		objType: i.objType,
		instID:  i.instID,
		data:    make(map[string]interface{}),
		ident:   new(metadata.HostIdentifier),
	}
	for key, value := range i.data {
		inst.data[key] = value
	}
	if i.objType == common.BKInnerObjIDHost {
		ident := new(metadata.HostIdentifier)
		marshaledIdent, _ := json.Marshal(i.ident)
		_ = json.Unmarshal(marshaledIdent, ident)
		inst.ident = ident
	}
	return inst
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
	ident.InnerIP = metadata.StringArrayToString(getString(m[common.BKHostInnerIPField]))
	ident.OuterIP = metadata.StringArrayToString(getString(m[common.BKHostOuterIPField]))
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
	ident.SupplierAccount = getString(m[common.BKOwnerIDField])
	ident.HostIdentModule = map[string]*metadata.HostIdentModule{}
	return &ident, nil
}

func getCache(ctx context.Context, cache redis.Client, clientSet apimachinery.ClientSetInterface, db dal.RDB,
	objType string, instID int64) (*Inst, error) {
	var err error
	header := getHeader()
	var instDataStr string
	inst := Inst{objType: objType, instID: instID, ident: &metadata.HostIdentifier{}, data: map[string]interface{}{}}
	switch objType {
	case common.BKInnerObjIDHost:
		ret := cache.Get(ctx, getInstCacheKey(objType, instID)).Val()
		if "" != ret && types.NilStr != ret {
			err = json.Unmarshal([]byte(ret), inst.ident)
			if err == nil {
				return &inst, nil
			}
			blog.Errorf("unmarshal error %s, raw is %s", err.Error(), ret)
		}

		// 0. get host data
		host, err := clientSet.CacheService().Cache().Host().SearchHostWithHostID(ctx, header, &metadata.SearchHostWithIDOption{
			HostID: instID, Fields: append(hostIndentDiffFields[common.BKInnerObjIDHost], common.BKHostIDField, common.BKOwnerIDField)})
		if err != nil {
			blog.Errorf("search host with id: %d failed, err: %s", instID, err.Error())
			return nil, err
		}

		if len(host) == 0 {
			return nil, nil
		}
		err = json.Unmarshal([]byte(host), &inst.data)
		if err != nil {
			blog.Errorf("unmarshal host %s failed, err: %s", host, err.Error())
			return nil, err
		}
		inst.ident, err = NewHostIdentifier(inst.data)
		if err != nil {
			blog.ErrorJSON("NewHostIdentifier using inst data %s error: %s", inst.data, err)
			return nil, err
		}

		// 1. fill modules TODO use cache to get host module relation when it is supported
		relations := make([]metadata.ModuleHost, 0)
		moduleHostCond := map[string]interface{}{
			common.BKHostIDField: instID,
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
		inst.ident, err = fillIdentifier(inst.ident, ctx, cache, clientSet, db)
		if err != nil {
			blog.ErrorJSON("fillIdentifier for ident: %s, error: %s", inst.ident, err)
			return nil, err
		}
		inst.data["associations"] = inst.ident.HostIdentModule
		inst.data["process"] = inst.ident.Process

		if err := inst.saveCache(cache); err != nil {
			blog.ErrorJSON("saveCache inst data %s failed, err: %s", inst.data, err)
		}
		return &inst, nil
	case common.BKInnerObjIDApp:
		instDataStr, err = clientSet.CacheService().Cache().Topology().SearchBusiness(ctx, header, instID)
		if err != nil {
			blog.Errorf("search biz with id: %d failed, err: %s", instID, err.Error())
			return nil, err
		}
	case common.BKInnerObjIDSet:
		instDataStr, err = clientSet.CacheService().Cache().Topology().SearchSet(ctx, header, instID)
		if err != nil {
			blog.Errorf("search set with id: %d failed, err: %s", instID, err.Error())
			return nil, err
		}
	case common.BKInnerObjIDModule:
		instDataStr, err = clientSet.CacheService().Cache().Topology().SearchModule(ctx, header, instID)
		if err != nil {
			blog.Errorf("search module with id: %d failed, err: %s", instID, err.Error())
			return nil, err
		}
	case common.BKInnerObjIDPlat, common.BKInnerObjIDProc:
		getInstCondition := map[string]interface{}{
			common.GetInstIDField(objType): instID,
		}
		if err = db.Table(common.GetInstTableName(objType)).Find(getInstCondition).One(ctx, &inst.data); err != nil {
			blog.ErrorJSON("find object %s inst %s error: %s", objType, instID, err)
			return nil, err
		}
		if len(inst.data) <= 0 {
			return nil, nil
		}
		return &inst, nil
	default:
		instDataStr, err = clientSet.CacheService().Cache().Topology().SearchCustomLayer(ctx, header, objType, instID)
		if err != nil {
			blog.Errorf("search custom layer %s with id: %d failed, err: %s", objType, instID, err.Error())
			return nil, err
		}
	}
	// handle inst data from cache
	if len(instDataStr) == 0 {
		return nil, nil
	}
	err = json.Unmarshal([]byte(instDataStr), &inst.data)
	if err != nil {
		blog.Errorf("unmarshal object %s inst %s failed, err: %s", objType, instDataStr, err.Error())
		return nil, err
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
	err = db.Table(common.BKTableNameServiceInstance).Find(serviceInstFilter).Fields(common.BKFieldID, common.BKModuleIDField).All(ctx, &serviceInstInfos)
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

func (ih *IdentifierHandler) popEvent() *metadata.EventInstCtx {
	rid := util.ExtractRequestIDFromContext(ih.ctx)

	eventStrs := ih.cache.BRPop(ih.ctx, time.Second*60, types.EventCacheEventQueueDuplicateKey).Val()

	if len(eventStrs) == 0 || eventStrs[1] == types.NilStr || len(eventStrs[1]) == 0 {
		return nil
	}

	// eventStrs format is []string{key, event}
	eventStr := eventStrs[1]
	event := metadata.EventInst{}
	if err := json.Unmarshal([]byte(eventStr), &event); err != nil {
		blog.Errorf("identifier: event distribute fail, unmarshal error: %+v, data=[%s], rid: %s", err, eventStr, rid)
		return nil
	}
	blog.V(3).Infof("pop new event for identifier, %+v", event)

	return &metadata.EventInstCtx{EventInst: event, Raw: eventStr}
}

func hasChanged(updateFields, deletedFields []string, fields ...string) bool {
	updateFieldMap := make(map[string]bool)
	for _, updateField := range updateFields {
		updateFieldMap[updateField] = true
	}

	deletedFieldMap := make(map[string]bool)
	for _, deletedField := range deletedFields {
		deletedFieldMap[deletedField] = true
	}

	for _, field := range fields {
		if updateFieldMap[field] {
			return true
		}
		if deletedFieldMap[field] {
			return true
		}

	}
	return false
}

func getHeader() http.Header {
	header := http.Header{}
	header.Add(common.BKHTTPOwnerID, common.BKSuperOwnerID)
	header.Add(common.BKHTTPHeaderUser, common.CCSystemOperatorUserName)
	return header
}

type IdentifierHandler struct {
	cache     redis.Client
	db        dal.RDB
	ctx       context.Context
	clientSet apimachinery.ClientSetInterface
}

func getInstCacheKey(objType string, instID int64) string {
	return types.EventCacheIdentInstPrefix + objType + "_" + strconv.FormatInt(instID, 10)
}

func NewIdentifierHandler(ctx context.Context, cache redis.Client, db dal.RDB, clientSet apimachinery.ClientSetInterface) *IdentifierHandler {
	return &IdentifierHandler{ctx: ctx, cache: cache, db: db, clientSet: clientSet}
}

func getBindInfo(value interface{}) (ip, port, protocol string, enable bool, bindInfoArr []metadata.ProcBindInfo) {
	if value == nil {
		return
	}
	bindInfoByteArr, err := json.Marshal(value)
	if err != nil {
		blog.Errorf("%v marshal error.", value, err.Error())
		return
	}
	bindInfoArr = make([]metadata.ProcBindInfo, 0)
	if err := json.Unmarshal(bindInfoByteArr, &bindInfoArr); err != nil {
		blog.Errorf("%s ummarshal error.", string(bindInfoByteArr), err.Error())
		return
	}
	for _, row := range bindInfoArr {
		if row.Std == nil {
			continue
		}
		if ip == "" && row.Std.IP != nil {
			ip = *row.Std.IP
		}
		if port == "" && row.Std.Port != nil {
			port = *row.Std.Port
		}
		if protocol == "" && row.Std.Protocol != nil {
			protocol = *row.Std.Protocol
		}
		if row.Std.Enable != nil && *row.Std.Enable == true {
			enable = true
		}
	}
	return
}
