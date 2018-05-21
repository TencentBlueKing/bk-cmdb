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
	"configcenter/src/source_controller/common/instdata"
	"encoding/json"
	redis "gopkg.in/redis.v5"
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
	hostIdentify := *e
	var ds []types.EventInst

	// add new dist if event belong to hostidentifier
	if diffFields, ok := hostIndentDiffFiels[e.ObjType]; ok && e.Action == types.EventActionUpdate && e.EventType == types.EventTypeInstData {
		for dataIndex := range e.Data {
			curdata := e.Data[dataIndex].CurData.(map[string]interface{})
			predata := e.Data[dataIndex].PreData.(map[string]interface{})
			if checkDifferent(curdata, predata, diffFields...) {
				hostIdentify.Data = nil
				hostIdentify.EventType = types.EventTypeRelation
				hostIdentify.ObjType = "hostidentifier"

				instID, _ := curdata[common.GetInstIDField(e.ObjType)].(int)
				if instID == 0 {
					// this should wound happen -_-
					blog.Errorf("conver instID faile the raw is %v", curdata[common.GetInstIDField(e.ObjType)])
					continue
				}

				count := 0
				inst, err := getCache(e.ObjType, instID)
				if err != nil {
					blog.Errorf("getCache error %v", err)
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
					redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
					hostIdentify.ID = redisCli.Incr(types.EventCacheEventIDKey).Val()
					ds = append(ds, hostIdentify)
				} else {

					total := len(identifiers)
					// pack identifiers into 1 distribution to prevent send too many messages
					for ident := range identifiers {
						count++
						d := types.EventData{PreData: *ident}
						d.CurData = *ident
						hostIdentify.Data = append(hostIdentify.Data, d)
						// each group is divided into 1000 units in order to limit the message size
						if count%1000 == 0 || count == total {
							ds = append(ds, hostIdentify)
							hostIdentify.Data = nil
						}
					}
				}
			}
		}
	} else if e.EventType == types.EventTypeRelation && hostIdentify.ObjType == "moduletransfer" {

	}

	return ds
}

func findHost(objType string, instID int) {
	api.GetAPIResource().InstCli.GetMutilByCondition(common.BKTableNameModuleHostConfig, fields, condiction, result, sort, start, limit)
}

type Inst map[string]interface{}

func (i Inst) set(key string, value interface{}) {
	i[key] = value
}

func (i Inst) saveCache() error {
	redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
	out, err := json.Marshal(i)
	if err != nil {
		return err
	}
	err = redisCli.Set("key", string(out), 0).Err()
	if err != nil {
		return err
	}
}

func getCache(objType string, instID int) (Inst, error) {
	redisCli := api.GetAPIResource().CacheCli.GetSession().(*redis.Client)
	ret := redisCli.Get("").String()
	inst := Inst{}
	if ret == "" || ret == "nil" {
		err := instdata.GetObjectByID(objType, nil, instID, inst, "")
		if err != nil {
			return nil, err
		}
	} else {
		err := json.Unmarshal([]byte(ret), inst)
		if err != nil {
			return nil, err
		}
	}

	if len(inst) <= 0 {
		return nil, nil
	}

	return inst, nil
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
