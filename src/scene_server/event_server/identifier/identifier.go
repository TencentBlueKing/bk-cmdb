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

	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/types"
	"configcenter/src/storage/dal"
)

var delayTime = time.Second * 30

var hostIndentDiffFiels = map[string][]string{
	common.BKInnerObjIDApp:    {common.BKAppNameField},
	common.BKInnerObjIDSet:    {common.BKSetNameField, "bk_service_status", "bk_set_env"},
	common.BKInnerObjIDModule: {common.BKModuleNameField},
	common.BKInnerObjIDPlat:   {common.BKCloudNameField},
	common.BKInnerObjIDProc: {common.BKProcessNameField, common.BKFuncIDField, common.BKFuncName,
		common.BKBindIP, common.BKProtocol, common.BKPort, "bk_start_param_regex"},
	common.BKInnerObjIDHost: {common.BKHostNameField,
		common.BKCloudIDField, common.BKHostInnerIPField, common.BKHostOuterIPField,
		common.BKOSTypeField, common.BKOSNameField,
		"bk_mem", "bk_cpu", "bk_disk"},
}

func (ih *IdentifierHandler) handleInst(e *metadata.EventInstCtx) {
	if diffFields, ok := hostIndentDiffFiels[e.ObjType]; ok && e.Action == metadata.EventActionUpdate && e.EventType == metadata.EventTypeInstData {
		ih.handleInstFieldChange(e, diffFields)
	} else if metadata.EventTypeRelation == e.EventType && "moduletransfer" == e.ObjType {
		ih.handleModuleTransfer(e)
	} else if metadata.EventTypeRelation == e.EventType && "processmodule" == e.ObjType {
		ih.handleBindProcess(e)
	}
}

func (ih *IdentifierHandler) handleInstFieldChange(e *metadata.EventInstCtx, diffFields []string) {
	blog.InfoJSON("identifier: handle inst %s", e)

	hostIdentify := e.EventInst
	hostIdentify.Data = nil
	hostIdentify.EventType = metadata.EventTypeRelation
	hostIdentify.ObjType = "hostidentifier"
	hostIdentify.Action = metadata.EventActionUpdate

	for dataIndex := range e.Data {
		curdata := e.Data[dataIndex].CurData.(map[string]interface{})
		predata := e.Data[dataIndex].PreData.(map[string]interface{})
		if checkDifferent(curdata, predata, diffFields...) {

			instIDField := common.GetInstIDField(e.ObjType)

			instID := getInt(curdata, instIDField)
			if 0 == instID {

				blog.Errorf("identifier: conver instID faile the raw is %+v", curdata[instIDField])
				continue
			}

			inst, err := getCache(ih.ctx, ih.cache, ih.db, e.ObjType, instID, false)
			if err != nil {
				blog.Errorf("identifier: getCache error %+v", err)
				continue
			}
			if nil == inst {
				blog.Errorf("identifier: inst == nil, continue")
				continue
			}
			for _, field := range diffFields {
				inst.set(field, curdata[field])
			}
			err = inst.saveCache(ih.cache)
			if err != nil {
				blog.Errorf("identifier: SaveCache error %+v", err)
				continue
			}

			if common.BKInnerObjIDHost == e.ObjType {
				hostIdentify.ID = ih.cache.Incr(types.EventCacheEventIDKey).Val()
				d := metadata.EventData{CurData: inst.ident.fillIden(ih.ctx, ih.cache, ih.db)}
				hostIdentify.Data = append(hostIdentify.Data, d)

				ih.cache.LPush(types.EventCacheEventQueueKey, &hostIdentify)
				blog.InfoJSON("identifier: pushed event inst %s", hostIdentify)
			} else {
				if err := ih.handleRelatedInst(hostIdentify, e.ObjType, instID, false); err != nil {
					blog.Warnf("handleRelatedInst faile objtype: %s, inst: %d, error: %v", e.ObjType, instID, err)
				}
			}
		}
	}
}

func (ih *IdentifierHandler) handleModuleTransfer(e *metadata.EventInstCtx) {
	blog.InfoJSON("identifier: handle inst %s", e)

	hostIdentify := e.EventInst
	hostIdentify.Data = nil
	hostIdentify.EventType = metadata.EventTypeRelation
	hostIdentify.ObjType = "hostidentifier"
	hostIdentify.Action = metadata.EventActionUpdate

	go func() {
		time.Sleep(delayTime)
		for index := range e.Data {
			var curdata map[string]interface{}
			var ok bool
			if metadata.EventActionDelete == e.Action {
				curdata, ok = e.Data[index].PreData.(map[string]interface{})
			} else {
				curdata, ok = e.Data[index].CurData.(map[string]interface{})
			}
			if !ok {
				continue
			}

			instID := getInt(curdata, common.BKHostIDField)
			if 0 == instID {
				blog.Errorf("identifier: conver instID faile the raw is %+v", curdata[common.BKHostIDField])
				continue
			}

			inst, err := getCache(ih.ctx, ih.cache, ih.db, common.BKInnerObjIDHost, instID, true)
			if err != nil {
				blog.Errorf("identifier: getCache error %+v", err)
				continue
			}
			if nil == inst {
				continue
			}

			inst.saveCache(ih.cache)
			d := metadata.EventData{CurData: inst.ident.fillIden(ih.ctx, ih.cache, ih.db)}
			hostIdentify.Data = append(hostIdentify.Data, d)
		}
		hostIdentify.ID = ih.cache.Incr(types.EventCacheEventIDKey).Val()
		ih.cache.LPush(types.EventCacheEventQueueKey, &hostIdentify)
		blog.InfoJSON("identifier: pushed event inst %s", hostIdentify)
	}()
}

func (ih *IdentifierHandler) handleBindProcess(e *metadata.EventInstCtx) {
	blog.InfoJSON("identifier: handle inst %s", e)

	hostIdentify := e.EventInst
	hostIdentify.Data = nil
	hostIdentify.EventType = metadata.EventTypeRelation
	hostIdentify.ObjType = "hostidentifier"
	hostIdentify.Action = metadata.EventActionUpdate

	go func() {
		time.Sleep(delayTime)
		for index := range e.Data {
			var curdata map[string]interface{}
			var ok bool
			if metadata.EventActionDelete == e.Action {
				curdata, ok = e.Data[index].PreData.(map[string]interface{})
			} else {
				curdata, ok = e.Data[index].CurData.(map[string]interface{})
			}
			if !ok {
				continue
			}

			instID := getInt(curdata, common.BKProcIDField)
			if 0 == instID {
				blog.Errorf("identifier: conver instID faile the raw is %+v", curdata[common.BKProcIDField])
				continue
			}

			modules := []metadata.ModuleInst{}
			cond := condition.CreateCondition().Field(common.BKSupplierIDField).Eq(curdata[common.BKSupplierIDField]).
				Field(common.BKAppIDField).Eq(curdata[common.BKAppIDField]).
				Field(common.BKModuleNameField).Eq(curdata[common.BKModuleNameField])
			if err := ih.db.Table(common.BKTableNameBaseModule).Find(cond.ToMapStr()).All(ih.ctx, &modules); err != nil {
				continue
			}

			for _, module := range modules {
				if err := ih.handleRelatedInst(hostIdentify, common.BKInnerObjIDModule, module.ModuleID, true); err != nil {
					blog.Warnf("handleRelatedInst faile objtype: %s, inst: %d, error: %v", e.ObjType, instID, err)
				}
			}

		}
	}()
}

func (ih *IdentifierHandler) handleRelatedInst(hostIdentify metadata.EventInst, objType string, instID int64, formdb bool) error {
	hosIDs, err := ih.findHost(objType, instID)
	if err != nil {
		blog.Warnf("identifier: find host faile: %v", err)
		return err
	}
	blog.V(3).Infof("identifier: handleRelatedInst by objType %s, instID %d,  hostIDs: %v, fromdb: %v", objType, instID, hosIDs, formdb)
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
				inst, err := getCache(ih.ctx, ih.cache, ih.db, common.BKInnerObjIDHost, hostID, true)
				if err != nil {
					blog.Errorf("identifier: getCache error %+v", err)
					continue
				}
				if nil == inst {
					continue
				}
				inst.saveCache(ih.cache)
				d := metadata.EventData{CurData: inst.ident.fillIden(ih.ctx, ih.cache, ih.db)}
				hostIdentify.Data = append(hostIdentify.Data, d)
			}
		} else {
			hostIDKeys := []string{}
			for _, hostID := range hosIDs[index:leftIndex] {
				hostIDKeys = append(hostIDKeys, getInstCacheKey(common.BKInnerObjIDHost, hostID))
			}
			idens, err := ih.cache.MGet(hostIDKeys...).Result()
			if err != nil {
				blog.Errorf("identifier: ih.cache.MGet by %v,%v. we will try to fetch it from db instead", hostIDKeys, err)
				idens = make([]interface{}, len(hostIDKeys))
				for index := range idens {
					// simulate that redis returns all nil
					idens[index] = "nil"
				}
			}
			for identIndex := range idens {
				iden := HostIdentifier{}
				if err = json.Unmarshal([]byte(getString(idens[identIndex])), &iden); err != nil {
					blog.Warnf("identifier: unmarshal error %s", err.Error())
					hostID := hosIDs[index:leftIndex][identIndex]
					inst, err := getCache(ih.ctx, ih.cache, ih.db, common.BKInnerObjIDHost, hostID, true)
					if err != nil {
						blog.Errorf("identifier: getCache error %+v", err)
						continue
					}
					if nil == inst {
						continue
					}
					inst.saveCache(ih.cache)
					d := metadata.EventData{CurData: inst.ident.fillIden(ih.ctx, ih.cache, ih.db)}
					hostIdentify.Data = append(hostIdentify.Data, d)
					continue
				}
				d := metadata.EventData{CurData: iden.fillIden(ih.ctx, ih.cache, ih.db)}
				hostIdentify.Data = append(hostIdentify.Data, d)
			}
		}
		index += 256

		hostIdentify.ID = ih.cache.Incr(types.EventCacheEventIDKey).Val()
		if err = ih.cache.LPush(types.EventCacheEventQueueKey, &hostIdentify).Err(); err != nil {
			blog.Warnf("identifier: push event inst %v faile %v", hostIdentify, err)
		} else {
			blog.InfoJSON("identifier: pushed event inst %s", hostIdentify)
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
	relations := []metadata.ModuleHost{}
	cond := condition.CreateCondition().Field(common.GetInstIDField(objType)).Eq(instID)

	if objType == common.BKInnerObjIDPlat {
		if err = ih.db.Table(common.BKTableNameBaseHost).Find(cond.ToMapStr()).Fields([]string{common.BKHostIDField}...).All(ih.ctx, &relations); err != nil {
			return nil, err
		}
	} else if objType == common.BKInnerObjIDProc {
		proc2module := []metadata.ProcessModule{}
		// get process to module
		if err = ih.db.Table(common.BKTableNameProcModule).Find(cond.ToMapStr()).All(ih.ctx, &proc2module); err != nil {
			return nil, err
		}
		if len(proc2module) > 0 {
			modulenames := make([]string, len(proc2module))
			for index := range proc2module {
				modulenames[index] = proc2module[index].ModuleName
			}
			// get module ids
			relations = []metadata.ModuleHost{}
			cond := condition.CreateCondition().Field(common.BKAppIDField).Eq(proc2module[0].AppID).Field(common.BKModuleNameField).In(modulenames)
			if err = ih.db.Table(common.BKTableNameBaseModule).Find(cond.ToMapStr()).Fields(common.BKModuleIDField).All(ih.ctx, &relations); err != nil {
				return nil, err
			}

			moduleids := make([]int64, len(relations))
			for index := range proc2module {
				moduleids[index] = relations[index].ModuleID
			}

			relations = []metadata.ModuleHost{}
			cond = condition.CreateCondition().Field(common.BKModuleIDField).In(moduleids)
			if err = ih.db.Table(common.BKTableNameModuleHostConfig).Find(cond.ToMapStr()).Fields(common.BKHostIDField).All(ih.ctx, &relations); err != nil {
				return nil, err
			}
		}
	} else {
		if err = ih.db.Table(common.BKTableNameModuleHostConfig).Find(cond.ToMapStr()).Fields(common.BKHostIDField).All(ih.ctx, &relations); err != nil {
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
	return cache.Set(getInstCacheKey(i.objType, i.instID), string(out), 0).Err()
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
	ident.Memory, err = util.GetInt64ByInterface(m["bk_mem"])
	if nil != err {
		blog.Errorf("%s is not integer, %+v ", "bk_mem", m)
	}
	ident.CPU, err = util.GetInt64ByInterface(m["bk_cpu"])
	if nil != err {
		blog.Errorf("%s is not integer, %+v ", "bk_cpu", m)
	}
	ident.Disk, err = util.GetInt64ByInterface(m["bk_disk"])
	if nil != err {
		blog.Errorf("%s is not integer, %+v ", "bk_disk", m)
	}
	ident.Module = map[string]*Module{}
	return &ident
}

func getCache(ctx context.Context, cache *redis.Client, db dal.RDB, objType string, instID int64, fromdb bool) (*Inst, error) {
	var err error
	ret := cache.Get(getInstCacheKey(objType, instID)).Val()
	inst := Inst{objType: objType, instID: instID, ident: &HostIdentifier{}, data: map[string]interface{}{}}
	if "" == ret || "nil" == ret || fromdb {
		blog.Infof("objType %s, instID %d not in cache, fetch it from db", objType, instID)
		getobjCondition := map[string]interface{}{
			common.GetInstIDField(objType): instID,
		}
		if err = db.Table(common.GetInstTableName(objType)).Find(getobjCondition).One(ctx, &inst.data); err != nil {
			return nil, err
		}
		if common.BKInnerObjIDHost == objType {
			inst.ident = NewHostIdentifier(inst.data)
			hostmoduleids := []int64{}

			// 1. fill modules
			relations := []metadata.ModuleHost{}
			moduleHostCond := map[string]interface{}{
				common.GetInstIDField(objType): instID,
			}
			if err = db.Table(common.BKTableNameModuleHostConfig).Find(moduleHostCond).All(ctx, &relations); err != nil {
				return nil, err
			}
			for _, rela := range relations {
				hostmoduleids = append(hostmoduleids, rela.ModuleID)
				inst.ident.Module[strconv.FormatInt(rela.ModuleID, 10)] = &Module{
					SetID:    rela.SetID,
					ModuleID: rela.ModuleID,
					BizID:    rela.AppID,
				}
			}
			inst.data["associations"] = inst.ident.Module

			// 2. fill process
			hostprocess := []Process{}

			// 2.1 find modules belongs to host
			modules := []metadata.ModuleInst{}
			cond := condition.CreateCondition().Field(common.BKModuleIDField).In(hostmoduleids)
			if err = db.Table(common.BKTableNameBaseModule).Find(cond.ToMapStr()).All(ctx, &modules); err != nil {
				return nil, err
			}

			// 2.2 find process belong to module within app
			appmodule := map[int64][]metadata.ModuleInst{}
			for _, module := range modules {
				appmodule[module.BizID] = append(appmodule[module.BizID], module)
			}
			for appid, modules := range appmodule {
				// 2.2.1 find process id belong to module within app
				moulename2ids := map[string][]int64{}
				procmoulenames := []string{}
				for _, module := range modules {
					moulename2ids[module.ModuleName] = append(moulename2ids[module.ModuleName], module.ModuleID)
					procmoulenames = append(procmoulenames, module.ModuleName)
				}
				proc2modules := []metadata.ProcessModule{}
				cond := condition.CreateCondition().Field(common.BKAppIDField).Eq(appid).Field(common.BKModuleNameField).In(procmoulenames)
				if err = db.Table(common.BKTableNameProcModule).Find(cond.ToMapStr()).All(ctx, &proc2modules); err != nil {
					return nil, err
				}

				// 2.2.2 find process by process id
				processids := []int64{}
				proc2moulenames := map[int64][]string{}
				for _, proc2module := range proc2modules {
					proc2moulenames[proc2module.ProcessID] = append(proc2moulenames[proc2module.ProcessID], proc2module.ModuleName)
					processids = append(processids, proc2module.ProcessID)
				}
				process := []Process{}
				cond = condition.CreateCondition().Field(common.BKProcIDField).In(processids)
				if err = db.Table(common.BKTableNameBaseProcess).Find(cond.ToMapStr()).All(ctx, &process); err != nil {
					return nil, err
				}

				// 2.3 bind module id
				for index := range process {
					for _, modulename := range proc2moulenames[process[index].ProcessID] {
						process[index].BindModules = moulename2ids[modulename]
					}
				}
				hostprocess = append(hostprocess, process...)
			}
			inst.ident.Process = hostprocess
			inst.data["process"] = hostprocess
		}
		inst.saveCache(cache)
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

func (ih *IdentifierHandler) StartHandleInsts() error {
	blog.Infof("identifier: handle identifiers started")
	go func() {
		ih.fetchHostCache()
		for range time.Tick(time.Minute * 10) {
			ih.fetchHostCache()
		}
	}()
	go ih.handleInstLoop()
	select {}
}

func (ih *IdentifierHandler) handleInstLoop() error {
	defer func() {
		procerr := recover()
		if procerr != nil {
			blog.Errorf("identifier: handleInstLoop panic: %v, stack:\n%s", procerr, debug.Stack())
		}
		go ih.handleInstLoop()
	}()
	for {
		event := ih.popEventInst()
		if nil == event {
			time.Sleep(time.Second * 2)
			continue
		}
		ih.handleInst(event)
	}
}

func (ih *IdentifierHandler) popEventInst() *metadata.EventInstCtx {

	eventstrs := ih.cache.BRPop(time.Second*60, types.EventCacheEventQueueDuplicateKey).Val()

	if 0 >= len(eventstrs) || "nil" == eventstrs[1] || "" == eventstrs[1] {
		return nil
	}

	eventstr := eventstrs[1]
	event := metadata.EventInst{}
	if err := json.Unmarshal([]byte(eventstr), &event); err != nil {
		blog.Errorf("identifier: event distribute fail, unmarshal error: %+v, date=[%s]", err, eventstr)
		return nil
	}

	return &metadata.EventInstCtx{EventInst: event, Raw: eventstr}
}

func (ih *IdentifierHandler) fetchHostCache() {

	relations := []metadata.ModuleHost{}
	hosts := []*HostIdentifier{}
	modules := []metadata.ModuleInst{}
	proc2modules := []metadata.ProcessModule{}

	err := ih.db.Table(common.BKTableNameModuleHostConfig).Find(map[string]interface{}{}).All(ih.ctx, &relations)
	if err != nil {
		blog.Errorf("[identifier][fetchHostCache] get cc_ModuleHostConfig error: %v", err)
		return
	}
	err = ih.db.Table(common.BKTableNameBaseHost).Find(map[string]interface{}{}).All(ih.ctx, &hosts)
	if err != nil {
		blog.Errorf("[identifier][fetchHostCache] get cc_HostBase error: %v", err)
		return
	}
	err = ih.db.Table(common.BKTableNameProcModule).Find(map[string]interface{}{}).All(ih.ctx, &proc2modules)
	if err != nil {
		blog.Errorf("[identifier][fetchHostCache] get cc_Proc2Module error: %v", err)
		return
	}
	err = ih.db.Table(common.BKTableNameBaseModule).Find(map[string]interface{}{}).All(ih.ctx, &modules)
	if err != nil {
		blog.Errorf("[identifier][fetchHostCache] get cc_ModuleBase error: %v", err)
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
		hostprocs := map[int64]bool{}
		for _, rela := range relationMap[ident.HostID] {
			ident.Module[strconv.FormatInt(rela.ModuleID, 10)] = &Module{
				SetID:    rela.SetID,
				ModuleID: rela.ModuleID,
				BizID:    rela.AppID,
			}
			if module, ok := modulesMap[rela.ModuleID]; ok {
				for _, procid := range proc2modulesMap[module.ModuleName] {
					hostprocs[procid] = true
				}
			}
		}

		ident.Process = []Process{}
		for procid := range hostprocs {
			bindModules := bindModulesMap[procid]
			if len(bindModules) == 0 {
				bindModules = []int64{}
			}
			ident.Process = append(ident.Process, Process{ProcessID: procid, BindModules: bindModules})
		}

		if err := ih.cache.Set(getInstCacheKey(common.BKInnerObjIDHost, ident.HostID), ident, 0).Err(); err != nil {
			blog.Errorf("set cache error %s", err.Error())
		}
	}
	blog.Infof("identifier: fetched %d hosts", len(hosts))

	objs := []string{common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule, common.BKInnerObjIDPlat, common.BKInnerObjIDProc}
	for _, objID := range objs {
		caches := []map[string]interface{}{}
		if err := ih.db.Table(common.GetInstTableName(objID)).Find(map[string]interface{}{}).All(ih.ctx, &caches); err != nil {
			blog.Errorf("set cache for objID %s error %v", objID, err)
		}

		for _, cache := range caches {
			out, _ := json.Marshal(cache)
			if err := ih.cache.Set(getInstCacheKey(objID, getInt(cache, common.GetInstIDField(objID))), string(out), 0).Err(); err != nil {
				blog.Errorf("set cache error %v", err)
			}
		}

		blog.Infof("identifier: fetched %d %s", len(caches), objID)
	}

}

func checkDifferent(curdata, predata map[string]interface{}, fields ...string) (isDifferent bool) {
	for _, field := range fields {
		if curdata[field] != predata[field] {
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
