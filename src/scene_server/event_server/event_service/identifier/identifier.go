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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/types"
	"configcenter/src/source_controller/api/metadata"
	"configcenter/src/source_controller/common/commondata"
	"configcenter/src/source_controller/common/instdata"
	"encoding/json"
	"fmt"
	redis "gopkg.in/redis.v5"
	"strconv"
	"time"
)

var hostIndentDiffFiels = map[string][]string{
	common.BKInnerObjIDApp:    {common.BKAppNameField},
	common.BKInnerObjIDSet:    {common.BKSetNameField, "bk_service_status", "bk_set_env"},
	common.BKInnerObjIDModule: {common.BKModuleNameField},
	common.BKInnerObjIDPlat:   {common.BKCloudNameField},
	common.BKInnerObjIDHost: {common.BKHostNameField,
		common.BKCloudIDField, common.BKHostInnerIPField, common.BKHostOuterIPField,
		common.BKOSTypeField, common.BKOSNameField,
		"bk_mem", "bk_cpu", "bk_disk"},
}

func handleInst(e *types.EventInst) {
	redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
	hostIdentify := *e
	hostIdentify.Data = nil
	hostIdentify.EventType = types.EventTypeRelation
	hostIdentify.ObjType = "hostidentifier"
	hostIdentify.Action = types.EventActionUpdate

	// add new dist if event belong to hostidentifier
	if diffFields, ok := hostIndentDiffFiels[e.ObjType]; ok && e.Action == types.EventActionUpdate && e.EventType == types.EventTypeInstData {
		blog.InfoJSON("identifier: handle inst %s", e)
		for dataIndex := range e.Data {
			curdata := e.Data[dataIndex].CurData.(map[string]interface{})
			predata := e.Data[dataIndex].PreData.(map[string]interface{})
			if checkDifferent(curdata, predata, diffFields...) {

				instIDField := util.GetObjIDByType(e.ObjType)

				instID := getInt(curdata, instIDField)
				if 0 == instID {
					// this should wound happen -_-
					blog.Errorf("identifier: conver instID faile the raw is %+v", curdata[instIDField])
					continue
				}

				inst, err := getCache(e.ObjType, instID, false)
				if err != nil {
					blog.Errorf("identifier: getCache error %+v", err)
					continue
				}
				if nil == inst {
					blog.Errorf("identifier: inst == nil, continue")
					// the inst may be deleted, just ignore
					continue
				}
				for _, field := range diffFields {
					inst.set(field, curdata[field])
				}
				err = inst.saveCache()
				if err != nil {
					blog.Errorf("identifier: SaveCache error %+v", err)
					continue
				}

				if common.BKInnerObjIDHost == e.ObjType {
					hostIdentify.ID = redisCli.Incr(types.EventCacheEventIDKey).Val()
					d := types.EventData{CurData: inst.ident.fillIden()}
					hostIdentify.Data = append(hostIdentify.Data, d)
					// TODO handle error
					redisCli.LPush(types.EventCacheEventQueueKey, &hostIdentify)
					blog.InfoJSON("identifier: pushed event inst %s", hostIdentify)
				} else {
					hosIDs := findHost(e.ObjType, instID)
					blog.Infof("identifier: hostIDs: %v", hosIDs)
					total := len(hosIDs)
					index := 0
					leftIndex := 0
					// pack identifiers into 1 distribution to prevent send too many messages
					for leftIndex < total {
						leftIndex = index + 256
						if leftIndex > total {
							leftIndex = total
						}
						hostIdentify.Data = nil
						idens := redisCli.MGet(hosIDs[index:leftIndex]...).Val()
						index += 256
						for identIndex := range idens {
							iden := HostIdentifier{}
							if err = json.Unmarshal([]byte(fmt.Sprint(idens[identIndex])), &iden); err != nil {
								blog.Errorf("identifier: unmarshal error %s", err.Error())
								continue
							}
							d := types.EventData{CurData: iden.fillIden()}
							hostIdentify.Data = append(hostIdentify.Data, d)
						}

						// handle error
						hostIdentify.ID = redisCli.Incr(types.EventCacheEventIDKey).Val()
						redisCli.LPush(types.EventCacheEventQueueKey, &hostIdentify)
						blog.InfoJSON("identifier: pushed event inst %s", hostIdentify)
					}
				}
			}
		}
	} else if types.EventTypeRelation == e.EventType && "moduletransfer" == e.ObjType {
		blog.Infof("identifier: handle inst %+v", e)
		go func() {
			time.Sleep(time.Second * 60) // delay to ensure moduletransfer ended
			for index := range e.Data {
				var curdata map[string]interface{}

				if types.EventActionDelete == e.Action {
					curdata, ok = e.Data[index].PreData.(map[string]interface{})
				} else {
					curdata, ok = e.Data[index].CurData.(map[string]interface{})
				}
				if !ok {
					continue
				}

				instID := getInt(curdata, common.BKHostIDField)
				if 0 == instID {
					// this should wound happen -_-
					blog.Errorf("identifier: conver instID faile the raw is %+v", curdata[common.BKHostIDField])
					continue
				}

				inst, err := getCache(common.BKInnerObjIDHost, instID, true)
				if err != nil {
					blog.Errorf("identifier: getCache error %+v", err)
					continue
				}
				if nil == inst {
					// the inst may be deleted, just ignore
					continue
				}

				// belong, ok := inst.data["associations"].(map[string]interface{})

				// // TODO 处理数据类型
				// moduleID := fmt.Sprint(curdata[common.BKModuleIDField])
				// switch e.Action {
				// case types.EventActionCreate:
				// 	if ok {
				// 		belong[moduleID] = curdata
				// 	}
				// 	inst.ident.Module[moduleID] = NewModule(curdata)
				// case types.EventActionDelete:
				// 	if ok {
				// 		delete(belong, moduleID)
				// 	}
				// 	delete(inst.ident.Module, moduleID)
				// }
				inst.saveCache()
				d := types.EventData{CurData: inst.ident.fillIden()}
				hostIdentify.Data = append(hostIdentify.Data, d)
			}
			hostIdentify.ID = redisCli.Incr(types.EventCacheEventIDKey).Val()
			redisCli.LPush(types.EventCacheEventQueueKey, &hostIdentify)
			blog.InfoJSON("identifier: pushed event inst %s", hostIdentify)
		}()
	}
}

func NewModule(m map[string]interface{}) *Module {
	belong := Module{}
	belong.BizID = getInt(m, common.BKAppIDField)
	belong.SetID = getInt(m, common.BKSetIDField)
	belong.ModuleID = getInt(m, common.BKModuleIDField)
	return &belong
}

func getInt(data map[string]interface{}, key string) int {
	i, err := strconv.Atoi(fmt.Sprint(data[key]))
	if err != nil {
		blog.Errorf("identifier: getInt error: %+v", err)
	}
	return i
}

func findHost(objType string, instID int) (hostIDs []string) {
	relations := []metadata.ModuleHostConfig{}
	condiction := map[string]interface{}{
		util.GetObjIDByType(objType): instID,
	}
	if objType == common.BKInnerObjIDPlat {
		// TODO handle error
		api.GetAPIResource().InstCli.GetMutilByCondition(common.BKTableNameBaseHost, []string{common.BKHostIDField}, condiction, &relations, "", -1, -1)
	} else {
		api.GetAPIResource().InstCli.GetMutilByCondition(common.BKTableNameModuleHostConfig, []string{common.BKHostIDField}, condiction, &relations, "", -1, -1)
	}

	for index := range relations {
		// TODO 抽象拼key
		hostIDs = append(hostIDs, types.EventCacheIdentInstPrefix+"host_"+strconv.Itoa(relations[index].HostID))
	}
	return hostIDs
}

type Inst struct {
	objType string
	instID  int
	data    map[string]interface{}
	ident   *HostIdentifier
}

func (i *Inst) set(key string, value interface{}) {
	i.data[key] = value
	var err error
	if i.objType == common.BKInnerObjIDHost {
		switch key {
		case "bk_host_name":
			i.ident.HostName = fmt.Sprint(value)
		case "bk_cloud_id":
			i.ident.CloudID, err = strconv.Atoi(fmt.Sprint(value))
		case "bk_host_innerip":
			i.ident.InnerIP = fmt.Sprint(value)
		case "bk_host_outerip":
			i.ident.OuterIP = fmt.Sprint(value)
		case "bk_os_type":
			i.ident.OSType = fmt.Sprint(value)
		case "bk_os_name":
			i.ident.OSName = fmt.Sprint(value)
		case "bk_mem":
			i.ident.Memory, err = strconv.ParseInt(fmt.Sprint(value), 10, 64)
		case "bk_cpu":
			i.ident.CPU, err = strconv.ParseInt(fmt.Sprint(value), 10, 64)
		case "bk_disk":
			i.ident.Disk, err = strconv.ParseInt(fmt.Sprint(value), 10, 64)
		}
		if nil != err {
			blog.Errorf("key %s	convert error %s", key, err.Error())
		}
	}
}

func (i *Inst) saveCache() error {
	redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
	out, err := json.Marshal(i.data)
	if err != nil {
		return err
	}
	err = redisCli.Set(types.EventCacheIdentInstPrefix+i.objType+fmt.Sprint("_", i.instID), string(out), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func NewHostIdentifier(m map[string]interface{}) *HostIdentifier {
	var err error
	ident := HostIdentifier{}
	ident.HostName = fmt.Sprint(m["bk_host_name"])
	ident.CloudID, err = strconv.Atoi(fmt.Sprint(m["bk_cloud_id"]))
	if nil != err {
		blog.Errorf("%s is not integer, %+v", "bk_cloud_id", m)
	}
	ident.InnerIP = fmt.Sprint(m["bk_host_innerip"])
	ident.OuterIP = fmt.Sprint(m["bk_host_outerip"])
	ident.OSType = fmt.Sprint(m["bk_os_type"])
	ident.OSName = fmt.Sprint(m["bk_os_name"])
	ident.Memory, err = strconv.ParseInt(fmt.Sprint(m["bk_mem"]), 10, 64)
	if nil != err {
		blog.Errorf("%s is not integer, %+v ", "bk_mem", m)
	}
	ident.CPU, err = strconv.ParseInt(fmt.Sprint(m["bk_cpu"]), 10, 64)
	if nil != err {
		blog.Errorf("%s is not integer, %+v ", "bk_cpu", m)
	}
	ident.Disk, err = strconv.ParseInt(fmt.Sprint(m["bk_disk"]), 10, 64)
	if nil != err {
		blog.Errorf("%s is not integer, %+v ", "bk_disk", m)
	}
	ident.Module = map[string]*Module{}
	return &ident
}
func getCache(objType string, instID int, fromdb bool) (*Inst, error) {
	redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
	ret := redisCli.Get(types.EventCacheIdentInstPrefix + objType + fmt.Sprint("_", instID)).Val()
	inst := Inst{objType: objType, instID: instID, ident: &HostIdentifier{}, data: map[string]interface{}{}}
	if "" == ret || "nil" == ret || fromdb {
		blog.Infof("objType %s, instID %d not in cache, fetch it from db", objType, instID)
		err := instdata.GetObjectByID(objType, nil, instID, &inst.data, "")
		if err != nil {
			return nil, err
		}
		if common.BKInnerObjIDHost == objType {
			inst.ident = NewHostIdentifier(inst.data)
			relations := []metadata.ModuleHostConfig{}
			condiction := map[string]interface{}{
				util.GetObjIDByType(objType): instID,
			}
			api.GetAPIResource().InstCli.GetMutilByCondition(common.BKTableNameModuleHostConfig, nil, condiction, &relations, "", -1, -1)
			for _, rela := range relations {
				inst.ident.Module[fmt.Sprint(rela.ModuleID)] = &Module{
					SetID:    rela.SetID,
					ModuleID: rela.ModuleID,
					BizID:    rela.ApplicationID,
				}
			}
			inst.data["associations"] = inst.ident.Module
		}
		inst.saveCache()
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

// StartHandleInsts handle the duplicate event queue
func StartHandleInsts() error {
	blog.Infof("identifier: handle identifiers started")
	go func() {
		fetchHostCache()
		for range time.Tick(time.Minute * 10) {
			fetchHostCache()
		}
	}()
	// TODO add
	for {
		event := popEventInst()
		if nil == event {
			time.Sleep(time.Second * 2)
			continue
		}
		handleInst(event)
	}
}

func popEventInst() *types.EventInst {

	redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)

	// TODO handle error
	eventstr := redisCli.BRPop(time.Second*60, types.EventCacheEventQueueDuplicateKey).Val()

	if 0 >= len(eventstr) || "nil" == eventstr[1] || "" == eventstr[1] {
		return nil
	}

	// Unmarshal event
	eventbytes := []byte(eventstr[1])
	event := types.EventInst{}
	if err := json.Unmarshal(eventbytes, &event); err != nil {
		blog.Errorf("identifier: event distribute fail, unmarshal error: %+v, date=[%s]", err, eventbytes)
		return nil
	}

	// blog.Infof("pop inst %s", eventbytes)
	return &event
}

func fetchHostCache() {
	redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)

	// fetch host cache
	relations := []metadata.ModuleHostConfig{}
	hosts := []*HostIdentifier{}

	// TODO handle db error, handle not found
	api.GetAPIResource().InstCli.GetMutilByCondition(common.BKTableNameModuleHostConfig, nil, map[string]interface{}{}, &relations, "", -1, -1)
	api.GetAPIResource().InstCli.GetMutilByCondition(common.BKTableNameBaseHost, nil, map[string]interface{}{}, &hosts, "", -1, -1)

	relationMap := map[int][]metadata.ModuleHostConfig{}
	for _, relate := range relations {
		relationMap[relate.HostID] = append(relationMap[relate.HostID], relate)
	}

	for _, ident := range hosts {
		ident.Module = map[string]*Module{}
		for _, rela := range relationMap[ident.HostID] {
			ident.Module[fmt.Sprint(rela.ModuleID)] = &Module{
				SetID:    rela.SetID,
				ModuleID: rela.ModuleID,
				BizID:    rela.ApplicationID,
			}
		}

		if err := redisCli.Set(types.EventCacheIdentInstPrefix+common.BKInnerObjIDHost+fmt.Sprint("_", ident.HostID), ident, 0).Err(); err != nil {
			blog.Errorf("set cache error %s", err.Error())
		}
	}
	blog.Infof("identifier: fetched %d hosts", len(hosts))

	// fetch others
	objs := []string{common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule, common.BKInnerObjIDPlat}
	for _, objID := range objs {
		caches := []map[string]interface{}{}
		api.GetAPIResource().InstCli.GetMutilByCondition(commondata.GetInstTableName(objID), nil, map[string]interface{}{}, &caches, "", -1, -1)

		for _, cache := range caches {
			out, _ := json.Marshal(cache)
			instID := fmt.Sprint(cache[util.GetObjIDByType(objID)])
			if err := redisCli.Set(types.EventCacheIdentInstPrefix+objID+fmt.Sprint("_", instID), string(out), 0).Err(); err != nil {
				blog.Errorf("set cache error %s", err.Error())
			}
		}

		blog.Infof("identifier: fetched %d %s", len(caches), objID)
	}

	// TODO compare data and build hostidentifier

}

func checkDifferent(curdata, predata map[string]interface{}, fields ...string) (isDifferent bool) {
	for _, field := range fields {
		if curdata[field] != predata[field] {
			return true
		}
	}
	return false
}
