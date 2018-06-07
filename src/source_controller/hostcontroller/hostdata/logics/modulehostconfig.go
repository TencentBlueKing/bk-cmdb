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

package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/util"
	eventtypes "configcenter/src/scene_server/event_server/types"
	metadataTable "configcenter/src/source_controller/api/metadata"
	"configcenter/src/source_controller/common/eventdata"
	"configcenter/src/source_controller/common/instdata"
	"errors"

	"gopkg.in/mgo.v2/bson"
)

var (
	moduleBaseTaleName = "cc_ModuleBase"
)

type moduleHostConfigParams struct {
	ApplicationID int   `json:"bk_biz_id"`
	HostID        int   `json:"bk_host_id"`
	ModuleID      []int `json:"bk_module_id"`
}

//DelSingleHostModuleRelation delete single host module relation
func DelSingleHostModuleRelation(ec *eventdata.EventContext, cc *api.APIResource, hostID, moduleID, appID int, ownerID string) (bool, error) {

	//get host info
	hostFieldArr := []string{common.BKHostInnerIPField}
	hostResult := make(map[string]interface{}, 0)
	errHost := instdata.GetObjectByID(common.BKInnerObjIDHost, hostFieldArr, hostID, &hostResult, common.BKHostIDField)
	blog.Infof("DelSingleHostModuleRelation hostID:%d, hostinfo:%v", hostID, hostResult)
	if errHost != nil {
		blog.Error("delSingleHostModuleRelation get host error:%s, host:%v", errHost.Error(), hostID)
		return false, errHost
	}

	moduleFieldArr := []string{common.BKModuleNameField}
	var moduleResult interface{}
	errModule := instdata.GetObjectByID(common.BKInnerObjIDModule, moduleFieldArr, moduleID, &moduleResult, common.BKModuleNameField)
	blog.Infof("DelSingleHostModuleRelation module:%d, module info:%v", moduleID, moduleResult)
	if errModule != nil {
		blog.Error("delSingleHostModuleRelation get module moduleID:%d, error:%s,", moduleID, errModule.Error())
		return false, errModule
	}

	tableName := metadataTable.ModuleHostConfig{}

	delCondition := make(map[string]interface{})
	delCondition[common.BKAppIDField] = appID
	delCondition[common.BKHostIDField] = hostID
	delCondition[common.BKModuleIDField] = moduleID
	delCondition = util.SetModOwner(delCondition, ownerID)
	num, numError := cc.InstCli.GetCntByCondition(tableName.TableName(), delCondition)
	blog.Infof("DelSingleHostModuleRelation  get module host relation condition:%v", delCondition)
	if numError != nil {
		blog.Error("delSingleHostModuleRelation get module host relation error:", numError.Error())
		return false, numError
	}
	//no config, return
	if num == 0 {
		return true, nil
	}

	// retrieve original datas
	origindatas := make([]map[string]interface{}, 0)
	getErr := cc.InstCli.GetMutilByCondition(tableName.TableName(), nil, delCondition, &origindatas, "", 0, 0)
	if getErr != nil {
		blog.Error("retrieve original datas error:%v", getErr)
		return false, getErr
	}

	delErr := cc.InstCli.DelByCondition(tableName.TableName(), delCondition)
	blog.Infof("DelSingleHostModuleRelation delCondition:%v", delCondition)
	if delErr != nil {
		blog.Error("delSingleHostModuleRelation del module host relation error:", delErr.Error())
		return false, delErr
	}

	// send events
	for _, origindata := range origindatas {
		err := ec.InsertEvent(eventtypes.EventTypeRelation, "moduletransfer", eventtypes.EventActionDelete, nil, origindata)
		if err != nil {
			blog.Error("create event error:%v", err)
		}
	}

	return true, nil
}

//AddSingleHostModuleRelation add single host module relation
func AddSingleHostModuleRelation(ec *eventdata.EventContext, cc *api.APIResource, hostID, moduleID, appID int, ownerID string) (bool, error) {
	//get host info
	hostFieldArr := []string{common.BKHostInnerIPField}
	hostResult := make(map[string]interface{})

	errHost := instdata.GetObjectByID(common.BKInnerObjIDHost, hostFieldArr, hostID, &hostResult, common.BKHostIDField)
	if errHost != nil {
		blog.Error("addSingleHostModuleRelation get host error:%s", errHost.Error())
		return false, errHost
	}

	moduleFieldArr := []string{common.BKModuleNameField, common.BKSetIDField}
	moduleResult := make(map[string]interface{})
	errModule := instdata.GetObjectByID(common.BKInnerObjIDModule, moduleFieldArr, moduleID, &moduleResult, common.BKModuleIDField)
	if errModule != nil {
		blog.Error("addSingleHostModuleRelation get module moduleid:%d, error:%s", moduleID, errModule.Error())
		return false, errModule
	}
	moduleName, _ := moduleResult[common.BKModuleNameField].(string)
	setID, _ := util.GetIntByInterface(moduleResult[common.BKSetIDField])

	if "" == moduleName || 0 == setID {
		blog.Error("addSingleHostModuleRelation get module error:not find module width ModuleID:%d", moduleID)
		return false, errors.New("未找到对应的模块")
	}

	tableName := metadataTable.ModuleHostConfig{}
	moduleHostConfig := make(map[string]interface{})

	moduleHostConfig[common.BKAppIDField] = appID
	moduleHostConfig[common.BKHostIDField] = hostID
	moduleHostConfig[common.BKModuleIDField] = moduleID
	moduleHostConfig = util.SetModOwner(moduleHostConfig, ownerID)

	num, numError := cc.InstCli.GetCntByCondition(tableName.TableName(), moduleHostConfig)
	if numError != nil {
		blog.Error("addSingleHostModuleRelation get module host relation error:", numError.Error())
		return false, numError
	}
	//config exsit, return
	if num > 0 {
		return true, nil
	}

	moduleHostConfig[common.BKSetIDField] = setID
	_, err := cc.InstCli.Insert(tableName.TableName(), moduleHostConfig)
	if err != nil {
		blog.Error("addSingleHostModuleRelation add module host relation error:", err.Error())
		return false, err
	}

	err = ec.InsertEvent(eventtypes.EventTypeRelation, "moduletransfer", eventtypes.EventActionCreate, moduleHostConfig, nil)
	if err != nil {
		blog.Error("create event error:%v", err)
	}

	return true, nil
}

//GetDefaultModuleIDs get default module ids
func GetDefaultModuleIDs(cc *api.APIResource, appID int, ownerID string) ([]int, error) {
	defaultModuleCond := make(map[string]interface{}, 2)
	defaultModuleCond[common.BKDefaultField] = common.KvMap{common.BKDBIN: []int{common.DefaultFaultModuleFlag, common.DefaultResModuleFlag}}
	defaultModuleCond[common.BKAppIDField] = appID
	defaultModuleCond = util.SetModOwner(defaultModuleCond, ownerID)
	result := make([]interface{}, 0)
	var ret []int

	err := cc.InstCli.GetMutilByCondition(moduleBaseTaleName, []string{common.BKModuleIDField, common.BKDefaultField}, defaultModuleCond, &result, "ID", 0, 100)
	blog.Infof("defaultModuleCond:%v", defaultModuleCond)
	if nil != err {
		blog.Errorf("getDefaultModuleIds error:%s, params:%v, %v", err.Error(), defaultModuleCond, result)
		return ret, errors.New("未找到模块")
	}

	for _, r := range result {
		item := r.(bson.M)
		ID, err := util.GetIntByInterface(item[common.BKModuleIDField])
		if nil != err {
			return ret, errors.New("未找到模块")
		}
		ret = append(ret, ID)
	}
	if 0 == len(ret) {
		return ret, errors.New("未找到模块")
	}

	return ret, nil
}

//GetModuleIDsByHostID get module id by hostid
func GetModuleIDsByHostID(cc *api.APIResource, moduleCond interface{}) ([]int, error) {
	result := make([]interface{}, 0)
	var ret []int

	tableName := metadataTable.ModuleHostConfig{}
	err := cc.InstCli.GetMutilByCondition(tableName.TableName(), []string{common.BKModuleIDField}, moduleCond, &result, "", 0, 100)
	blog.Infof("GetModuleIDsByHostID condition:%v", moduleCond)
	blog.Infof("result:%v", result)
	if nil != err {
		blog.Error("getModuleIDsByHostID error:%", err.Error())
		return ret, errors.New("未找到主机所属模块")
	}
	for _, r := range result {
		item := r.(bson.M)
		ID, getErr := util.GetIntByInterface(item[common.BKModuleIDField])
		if nil != getErr || ID == 0 {
			return ret, errors.New("未找到模块")
		}
		ret = append(ret, ID)
	}
	return ret, err
}

//GetResourcePoolApp get resource pool app
func GetResourcePoolApp(cc *api.APIResource, ownerID int) (int, error) {
	params := make(map[string]interface{})
	params[common.BKOwnerIDField] = ownerID
	params[common.BKDefaultField] = 1

	result := make(map[string]interface{})
	err := cc.InstCli.GetOneByCondition("cc_ApplicationBase", []string{common.BKAppIDField}, params, &result)
	if nil != err {
		blog.Error("getModuleIDsByHostID error:%", err.Error())
		return 0, errors.New("获取资源池业务失败")
	}
	appID, _ := util.GetIntByInterface(result[common.BKAppIDField])
	if 0 == appID {
		blog.Error("getModuleIDsByHostID error: 未找到默认业务")
		return 0, errors.New("未发现资源池业务")
	}

	return appID, nil
}

//check if host belong to empty module
func CheckHostInIDle(cc *api.APIResource, appID, emptyModuleID int, hostIDs []int) ([]int, []int, error) {

	moduleHostConfig := metadataTable.ModuleHostConfig{}
	conds := make(map[string]interface{}, 1)
	conds[common.BKHostIDField] = bson.M{common.BKDBIN: hostIDs}
	result := make([]interface{}, 0)

	err := cc.InstCli.GetMutilByCondition(moduleHostConfig.TableName(), []string{common.BKHostIDField, common.BKModuleIDField, common.BKAppIDField}, conds, &result, "", 0, common.BKNoLimit)
	if nil != err {
		blog.Error("get modulehostconfig error:%s", err.Error())
		return nil, nil, errors.New("获取主机与模块关系失败")
	}
	var errHostIDs []int
	var faultHostIDs []int
	mapHost := make(map[int]int, 0)
	for _, item := range result {
		row := item.(bson.M)
		moduleID, getErr := util.GetIntByInterface(row[common.BKModuleIDField])
		if nil != getErr {
			continue
		}
		hostID, getErr := util.GetIntByInterface(row[common.BKHostIDField])
		if nil != getErr {
			continue
		}
		rowAppID, getErr := util.GetIntByInterface(row[common.BKAppIDField])
		if nil != getErr {
			continue
		}
		//host not belong to this biz
		if rowAppID != appID {
			faultHostIDs = append(faultHostIDs, hostID)
		}
		//host belong to this biz, but not in idle module
		if moduleID != emptyModuleID && rowAppID == appID {
			_, ok := mapHost[hostID]
			if !ok {
				errHostIDs = append(errHostIDs, hostID)
				mapHost[hostID] = hostID
			}
		}

	}

	return errHostIDs, faultHostIDs, err
}

//获取业务下的默认模块
func GetIDleModuleID(cc *api.APIResource, appID int) (int, error) {
	defaultModuleCond := make(map[string]interface{}, 2)
	defaultModuleCond[common.BKDefaultField] = common.DefaultResModuleFlag
	defaultModuleCond[common.BKAppIDField] = appID
	result := make(map[string]interface{}, 0)
	err := cc.InstCli.GetOneByCondition(moduleBaseTaleName, []string{common.BKModuleIDField}, defaultModuleCond, &result)

	if nil != err {
		blog.Error("getDefaultModuleIDs error:%s", err.Error())
		return 0, errors.New("未找到模块")
	}

	ID, ok := util.GetIntByInterface(result[common.BKModuleIDField])
	if nil != ok {
		return ID, errors.New("未找到模块")
	}

	return ID, nil
}
