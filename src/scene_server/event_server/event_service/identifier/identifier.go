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
	"configcenter/src/scene_server/event_server/types"
	"configcenter/src/source_controller/api/metadata"
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

func GetDistInst(e *types.EventInst) []types.EventInst {
	redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
	hostIdentify := *e
	hostIdentify.Data = nil
	hostIdentify.EventType = types.EventTypeRelation
	hostIdentify.ObjType = "hostidentifier"
	hostIdentify.Action = types.EventActionUpdate
	var ds []types.EventInst

	// add new dist if event belong to hostidentifier
	if diffFields, ok := hostIndentDiffFiels[e.ObjType]; ok && e.Action == types.EventActionUpdate && e.EventType == types.EventTypeInstData {
		for dataIndex := range e.Data {
			curdata := e.Data[dataIndex].CurData.(map[string]interface{})
			predata := e.Data[dataIndex].PreData.(map[string]interface{})
			if checkDifferent(curdata, predata, diffFields...) {

				instIDField := common.GetInstIDField(e.ObjType)

				instID, _ := curdata[instIDField].(int)
				if instID == 0 {
					// this should wound happen -_-
					blog.Errorf("conver instID faile the raw is %v", curdata[instIDField])
					continue
				}

				inst, err := getCache(e.ObjType, instID)
				if err != nil {
					blog.Errorf("getCache error %v", err)
					continue
				}
				if inst == nil {
					// the inst may be deleted, just ignore
					continue
				}
				for _, field := range diffFields {
					inst.set(field, curdata[field])
				}
				err = inst.saveCache()
				if err != nil {
					blog.Errorf("SaveCache error %v", err)
					continue
				}

				if e.ObjType == common.BKInnerObjIDHost {
					hostIdentify.ID = redisCli.Incr(types.EventCacheEventIDKey).Val()
					d := types.EventData{CurData: inst.data}
					hostIdentify.Data = append(hostIdentify.Data, d)
					redisCli.LPush(types.EventCacheEventQueueKey, hostIdentify)
				} else {
					hosIDs := findHost(e.ObjType, instID)
					total := len(hosIDs)
					index := 0
					leftIndex := 0
					// pack identifiers into 1 distribution to prevent send too many messages
					for {
						leftIndex = index + 256
						if leftIndex > total {
							leftIndex = total - 1
						}
						hostIdentify.Data = nil
						idens := redisCli.MGet(hosIDs[index:leftIndex]...).Val()
						index += 256
						for identIndex := range idens {
							iden := HostIdentifier{}
							if err = json.Unmarshal([]byte(fmt.Sprint(idens[identIndex])), &iden); err != nil {
								continue
							}
							iden.fillIden()
							d := types.EventData{CurData: &iden}
							hostIdentify.Data = append(hostIdentify.Data, d)
						}

						hostIdentify.ID = redisCli.Incr(types.EventCacheEventIDKey).Val()
						redisCli.LPush(types.EventCacheEventQueueKey, hostIdentify)
					}
				}
			}
		}
	} else if e.EventType == types.EventTypeRelation && hostIdentify.ObjType == "moduletransfer" {
		for index := range e.Data {
			var curdata map[string]interface{}

			if e.Action == "delete" && len(e.Data) > 0 {
				curdata, ok = e.Data[index].PreData.(map[string]interface{})
			} else {
				curdata, ok = e.Data[index].CurData.(map[string]interface{})
			}
			if !ok {
				continue
			}

			instID := getInt(curdata, common.BKHostIDField)
			if instID == 0 {
				// this should wound happen -_-
				blog.Errorf("conver instID faile the raw is %v", curdata[common.BKHostIDField])
				continue
			}

			inst, err := getCache(common.BKInnerObjIDHost, instID)
			if err != nil {
				blog.Errorf("getCache error %v", err)
				continue
			}
			if inst == nil {
				// the inst may be deleted, just ignore
				continue
			}

			moduleID := fmt.Sprint(curdata[common.BKModuleIDField])
			switch e.Action {
			case types.EventActionCreate:
				inst.data[moduleID] = curdata
			case types.EventActionDelete:
				delete(inst.data, moduleID)
			}
		}
		hostIdentify.ID = redisCli.Incr(types.EventCacheEventIDKey).Val()
		redisCli.LPush(types.EventCacheEventQueueKey, hostIdentify)
	}

	return ds
}

func getInt(data map[string]interface{}, key string) int {
	i, err := strconv.Atoi(fmt.Sprint(data[key]))
	if err != nil {
		blog.Errorf("getInt error: %v", err)
	}
	return i
}

func findHost(objType string, instID int) (hostIDs []string) {
	// TODO cloud_name not handled
	if objType == common.BKInnerObjIDPlat {

	}
	relations := []metadata.ModuleHostConfig{}
	condiction := map[string]interface{}{
		common.GetInstIDField(objType): instID,
	}
	api.GetAPIResource().InstCli.GetMutilByCondition(common.BKTableNameModuleHostConfig, []string{common.BKHostIDField}, condiction, &relations, "", -1, -1)

	for index := range relations {
		hostIDs = append(hostIDs, types.EventCacheIdentInstPrefix+"_host_"+strconv.Itoa(relations[index].HostID))
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
}

func (i *Inst) saveCache() error {
	redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
	out, err := json.Marshal(i)
	if err != nil {
		return err
	}
	err = redisCli.Set(types.EventCacheIdentInstPrefix+i.objType+fmt.Sprint("_", i.instID), string(out), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func getCache(objType string, instID int) (*Inst, error) {
	redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
	ret := redisCli.Get(types.EventCacheIdentInstPrefix + objType + fmt.Sprint("_", instID)).String()
	inst := Inst{objType: objType, instID: instID}
	if ret == "" || ret == "nil" {
		err := instdata.GetObjectByID(objType, nil, instID, inst.data, "")
		if err != nil {
			return nil, err
		}
	} else {
		err := json.Unmarshal([]byte(ret), inst.data)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(ret), inst.ident)
		if err != nil {
			return nil, err
		}
	}

	if len(inst.data) <= 0 {
		return nil, nil
	}

	return &inst, nil
}

// StartHandleInsts handle the duplicate event queue
func StartHandleInsts() error {
	for {
		event := popEventInst()
		if event == nil {
			time.Sleep(time.Second * 2)
			continue
		}
		if err := handleInst(event); err != nil {
			blog.Errorf("error handle dist: %v, %v", err, event)
		}
	}
}

func handleInst(event *types.EventInstCtx) (err error) {
	blog.Info("handling event inst : %v", event.Raw)
	defer blog.Info("done event inst : %v", event.ID)

	origindists := GetDistInst(&event.EventInst)

	return
}

func popEventInst() *types.EventInstCtx {
	var eventstr string

	redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
	redisCli.BRPopLPush(types.EventCacheEventQueueKey, types.EventCacheEventQueueDuplicateKey, time.Second*60).Scan(&eventstr)

	if eventstr == "" {
		return nil
	}

	// Unmarshal event
	eventbytes := []byte(eventstr)
	event := types.EventInst{}
	if err := json.Unmarshal(eventbytes, &event); err != nil {
		blog.Errorf("event distribute fail, unmarshal error: %v, date=[%s]", err, eventbytes)
		return nil
	}

	return &types.EventInstCtx{EventInst: event, Raw: eventstr}
}

func checkDifferent(curdata, predata map[string]interface{}, fields ...string) (isDifferent bool) {
	for _, field := range fields {
		if curdata[field] != predata[field] {
			return true
		}
	}
	return false
}
