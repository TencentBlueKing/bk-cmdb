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
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	scmeta "configcenter/src/source_controller/api/metadata"
	"configcenter/src/storage"
	"configcenter/src/txn_server/types"
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

func (ih *IdentifierHandler) handleInst(e *metadata.EventInst) {
	hostIdentify := *e
	hostIdentify.Data = nil
	hostIdentify.EventType = metadata.EventTypeRelation
	hostIdentify.ObjType = "hostidentifier"
	hostIdentify.Action = metadata.EventActionUpdate

	if diffFields, ok := hostIndentDiffFiels[e.ObjType]; ok && e.Action == metadata.EventActionUpdate && e.EventType == metadata.EventTypeInstData {
		blog.InfoJSON("identifier: handle inst %s", e)
		for dataIndex := range e.Data {
			curdata := e.Data[dataIndex].CurData.(map[string]interface{})
			predata := e.Data[dataIndex].PreData.(map[string]interface{})
			if checkDifferent(curdata, predata, diffFields...) {

				instIDField := common.GetInstFieldByType(e.ObjType)

				instID := getInt(curdata, instIDField)
				if 0 == instID {

					blog.Errorf("identifier: conver instID faile the raw is %+v", curdata[instIDField])
					continue
				}

				inst, err := getCache(ih.cache, ih.db, e.ObjType, instID, false)
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
					d := metadata.EventData{CurData: inst.ident.fillIden(ih.cache, ih.db)}
					hostIdentify.Data = append(hostIdentify.Data, d)

					ih.cache.LPush(types.EventCacheEventQueueKey, &hostIdentify)
					blog.InfoJSON("identifier: pushed event inst %s", hostIdentify)
				} else {
					hosIDs := ih.findHost(e.ObjType, instID)
					blog.Infof("identifier: hostIDs: %v", hosIDs)
					total := len(hosIDs)
					index := 0
					leftIndex := 0

					for leftIndex < total {
						leftIndex = index + 256
						if leftIndex > total {
							leftIndex = total
						}
						hostIdentify.Data = nil
						idens := ih.cache.MGet(hosIDs[index:leftIndex]...).Val()
						index += 256
						for identIndex := range idens {
							iden := HostIdentifier{}
							if err = json.Unmarshal([]byte(fmt.Sprint(idens[identIndex])), &iden); err != nil {
								blog.Errorf("identifier: unmarshal error %s", err.Error())
								continue
							}
							d := metadata.EventData{CurData: iden.fillIden(ih.cache, ih.db)}
							hostIdentify.Data = append(hostIdentify.Data, d)
						}

						hostIdentify.ID = ih.cache.Incr(types.EventCacheEventIDKey).Val()
						ih.cache.LPush(types.EventCacheEventQueueKey, &hostIdentify)
						blog.InfoJSON("identifier: pushed event inst %s", hostIdentify)
					}
				}
			}
		}
	} else if metadata.EventTypeRelation == e.EventType && "moduletransfer" == e.ObjType {
		blog.Infof("identifier: handle inst %+v", e)
		go func() {
			time.Sleep(time.Second * 60)
			for index := range e.Data {
				var curdata map[string]interface{}

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

				inst, err := getCache(ih.cache, ih.db, common.BKInnerObjIDHost, instID, true)
				if err != nil {
					blog.Errorf("identifier: getCache error %+v", err)
					continue
				}
				if nil == inst {

					continue
				}

				inst.saveCache(ih.cache)
				d := metadata.EventData{CurData: inst.ident.fillIden(ih.cache, ih.db)}
				hostIdentify.Data = append(hostIdentify.Data, d)
			}
			hostIdentify.ID = ih.cache.Incr(types.EventCacheEventIDKey).Val()
			ih.cache.LPush(types.EventCacheEventQueueKey, &hostIdentify)
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

func (ih *IdentifierHandler) findHost(objType string, instID int) (hostIDs []string) {
	relations := []scmeta.ModuleHostConfig{}
	condiction := map[string]interface{}{
		common.GetInstFieldByType(objType): instID,
	}
	if objType == common.BKInnerObjIDPlat {

		ih.db.GetMutilByCondition(common.BKTableNameBaseHost, []string{common.BKHostIDField}, condiction, &relations, "", -1, -1)
	} else {
		ih.db.GetMutilByCondition(common.BKTableNameModuleHostConfig, []string{common.BKHostIDField}, condiction, &relations, "", -1, -1)
	}

	for index := range relations {

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

func (i *Inst) saveCache(cache *redis.Client) error {
	out, err := json.Marshal(i.data)
	if err != nil {
		return err
	}
	err = cache.Set(types.EventCacheIdentInstPrefix+i.objType+fmt.Sprint("_", i.instID), string(out), 0).Err()
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
	ident.HostID, err = util.GetIntByInterface(m[common.BKHostIDField])
	if nil != err {
		blog.Errorf("%s is not integer, %+v ", "bk_host_id", m)
	}
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
func getCache(cache *redis.Client, db storage.DI, objType string, instID int, fromdb bool) (*Inst, error) {
	ret := cache.Get(types.EventCacheIdentInstPrefix + objType + fmt.Sprint("_", instID)).Val()
	inst := Inst{objType: objType, instID: instID, ident: &HostIdentifier{}, data: map[string]interface{}{}}
	if "" == ret || "nil" == ret || fromdb {
		blog.Infof("objType %s, instID %d not in cache, fetch it from db", objType, instID)
		getobjCondition := map[string]interface{}{
			common.GetInstFieldByType(objType): instID,
		}
		err := db.GetOneByCondition(common.GetInstTableName(objType), nil, getobjCondition, &inst.data)
		if err != nil {
			return nil, err
		}
		if common.BKInnerObjIDHost == objType {
			inst.ident = NewHostIdentifier(inst.data)
			relations := []scmeta.ModuleHostConfig{}
			condiction := map[string]interface{}{
				common.GetInstFieldByType(objType): instID,
			}
			db.GetMutilByCondition(common.BKTableNameModuleHostConfig, nil, condiction, &relations, "", -1, -1)
			for _, rela := range relations {
				inst.ident.Module[fmt.Sprint(rela.ModuleID)] = &Module{
					SetID:    rela.SetID,
					ModuleID: rela.ModuleID,
					BizID:    rela.ApplicationID,
				}
			}
			inst.data["associations"] = inst.ident.Module
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

	for {
		event := ih.popEventInst()
		if nil == event {
			time.Sleep(time.Second * 2)
			continue
		}
		ih.handleInst(event)
	}
}

func (ih *IdentifierHandler) popEventInst() *metadata.EventInst {

	eventstr := ih.cache.BRPop(time.Second*60, types.EventCacheEventQueueDuplicateKey).Val()

	if 0 >= len(eventstr) || "nil" == eventstr[1] || "" == eventstr[1] {
		return nil
	}

	eventbytes := []byte(eventstr[1])
	event := metadata.EventInst{}
	if err := json.Unmarshal(eventbytes, &event); err != nil {
		blog.Errorf("identifier: event distribute fail, unmarshal error: %+v, date=[%s]", err, eventbytes)
		return nil
	}

	return &event
}

func (ih *IdentifierHandler) fetchHostCache() {

	relations := []scmeta.ModuleHostConfig{}
	hosts := []*HostIdentifier{}

	ih.db.GetMutilByCondition(common.BKTableNameModuleHostConfig, nil, map[string]interface{}{}, &relations, "", -1, -1)
	ih.db.GetMutilByCondition(common.BKTableNameBaseHost, nil, map[string]interface{}{}, &hosts, "", -1, -1)

	relationMap := map[int][]scmeta.ModuleHostConfig{}
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

		if err := ih.cache.Set(types.EventCacheIdentInstPrefix+common.BKInnerObjIDHost+fmt.Sprint("_", ident.HostID), ident, 0).Err(); err != nil {
			blog.Errorf("set cache error %s", err.Error())
		}
	}
	blog.Infof("identifier: fetched %d hosts", len(hosts))

	objs := []string{common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule, common.BKInnerObjIDPlat}
	for _, objID := range objs {
		caches := []map[string]interface{}{}
		ih.db.GetMutilByCondition(common.GetInstTableName(objID), nil, map[string]interface{}{}, &caches, "", -1, -1)

		for _, cache := range caches {
			out, _ := json.Marshal(cache)
			instID := fmt.Sprint(cache[common.GetInstFieldByType(objID)])
			if err := ih.cache.Set(types.EventCacheIdentInstPrefix+objID+fmt.Sprint("_", instID), string(out), 0).Err(); err != nil {
				blog.Errorf("set cache error %s", err.Error())
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
	db    storage.DI
}

func NewIdentifierHandler(cache *redis.Client, db storage.DI) *IdentifierHandler {
	return &IdentifierHandler{cache: cache, db: db}
}
