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
	"errors"
	"fmt"

	"gopkg.in/mgo.v2/bson"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage"
)

type Logics struct {
	Instance storage.DI
}

const (
	ModuleBaseCollectioin     = "cc_ModuleBase"
	ModuleHostCollection      = "cc_ModuleHostConfig"
	ApplicationBaseCollection = "cc_ApplicationBase"
)

//DelSingleHostModuleRelation delete single host module relation
func (lgc *Logics) DelSingleHostModuleRelation(ec *eventclient.EventContext, hostID, moduleID, appID int64, ownerID string) (bool, error) {

	hostFieldArr := []string{common.BKHostInnerIPField}
	hostResult := make(map[string]interface{}, 0)
	errHost := lgc.GetObjectByID(common.BKInnerObjIDHost, hostFieldArr, hostID, &hostResult, common.BKHostIDField)
	blog.Infof("delete single host relation, hostID: %d, host info: %v", hostID, hostResult)
	if errHost != nil {
		blog.Errorf("delete single host relation failed, host: %v, err: %v", hostID, errHost)
		return false, errHost
	}

	moduleFieldArr := []string{common.BKModuleNameField}
	var moduleResult interface{}
	errModule := lgc.GetObjectByID(common.BKInnerObjIDModule, moduleFieldArr, moduleID, &moduleResult, common.BKModuleNameField)
	if errModule != nil {
		blog.Errorf("delete single host relation, but get module failed,  moduleID:%d, error:%s,", moduleID, errModule.Error())
		return false, errModule
	}

	delCondition := common.KvMap{common.BKAppIDField: appID, common.BKHostIDField: hostID, common.BKModuleIDField: moduleID}
	delCondition = util.SetModOwner(delCondition, ownerID)
	num, numError := lgc.Instance.GetCntByCondition(ModuleHostCollection, delCondition)
	if numError != nil {
		blog.Errorf("delete single host relation, but get module host relation failed, err: %v", numError)
		return false, numError
	}

	if num == 0 {
		return true, nil
	}

	// retrieve original datas
	origindatas := make([]map[string]interface{}, 0)
	getErr := lgc.Instance.GetMutilByCondition(ModuleHostCollection, nil, delCondition, &origindatas, "", 0, 0)
	if getErr != nil {
		blog.Errorf("delete single host relation, retrieve original data error:%v", getErr)
		return false, getErr
	}

	delErr := lgc.Instance.DelByCondition(ModuleHostCollection, delCondition)
	if delErr != nil {
		blog.Errorf("delete single host relation, but del module host relation failed, err: %v", delErr)
		return false, delErr
	}

	// send events
	for _, origindata := range origindatas {
		err := ec.InsertEvent(metadata.EventTypeRelation, "moduletransfer", metadata.EventActionDelete, nil, origindata)
		if err != nil {
			blog.Errorf("delete single host relation failed, but create event error:%v", err)
		}
	}

	return true, nil
}

//AddSingleHostModuleRelation add single host module relation
func (lgc *Logics) AddSingleHostModuleRelation(ec *eventclient.EventContext, hostID, moduleID, appID int64, ownerID string) (bool, error) {
	hostFieldArr := []string{common.BKHostInnerIPField}
	hostResult := make(map[string]interface{})
	errHost := lgc.GetObjectByID(common.BKInnerObjIDHost, hostFieldArr, hostID, &hostResult, common.BKHostIDField)
	if errHost != nil {
		blog.Errorf("add single host module relation, but get host error:%s", errHost.Error())
		return false, errHost
	}

	moduleFieldArr := []string{common.BKModuleNameField, common.BKSetIDField}
	moduleResult := make(map[string]interface{})
	errModule := lgc.GetObjectByID(common.BKInnerObjIDModule, moduleFieldArr, moduleID, &moduleResult, common.BKModuleIDField)
	if errModule != nil {
		blog.Errorf("add single host module relation, get module moduleid:%d, error:%s", moduleID, errModule.Error())
		return false, errModule
	}
	moduleName, ok := moduleResult[common.BKModuleNameField].(string)
	if !ok {
		return false, errors.New("invalid module name")
	}
	setID, err := util.GetInt64ByInterface(moduleResult[common.BKSetIDField])
	if err != nil {
		return false, fmt.Errorf("invalid set id, err: %v", err)
	}

	if "" == moduleName || 0 == setID {
		blog.Errorf("add single host module relation, get module error:not find module width ModuleID:%d", moduleID)
		return false, errors.New("can not find it's module")
	}

	moduleHostConfig := common.KvMap{common.BKAppIDField: appID, common.BKHostIDField: hostID, common.BKModuleIDField: moduleID}
	num, numError := lgc.Instance.GetCntByCondition(common.BKTableNameModuleHostConfig, moduleHostConfig)
	if numError != nil {
		blog.Error("add single host module relation, get module host relation error: %v", numError)
		return false, numError
	}

	if num > 0 {
		return true, nil
	}

	moduleHostConfig[common.BKSetIDField] = setID
	moduleHostConfig = util.SetModOwner(moduleHostConfig, ownerID)
	_, err = lgc.Instance.Insert(common.BKTableNameModuleHostConfig, moduleHostConfig)
	if err != nil {
		blog.Errorf("add single host module relation, add module host relation error:", err.Error())
		return false, err
	}

	err = ec.InsertEvent(metadata.EventTypeRelation, "moduletransfer", metadata.EventActionCreate, moduleHostConfig, nil)
	if err != nil {
		blog.Errorf("add single host module relation, create event error:%v", err)
	}

	return true, nil
}

//GetDefaultModuleIDs get default module ids
func (lgc *Logics) GetDefaultModuleIDs(appID int64) ([]int64, error) {
	defaultModuleCond := make(map[string]interface{})
	defaultModuleCond[common.BKDefaultField] = common.KvMap{common.BKDBIN: []int64{int64(common.DefaultFaultModuleFlag), int64(common.DefaultResModuleFlag)}}
	defaultModuleCond[common.BKAppIDField] = appID
	result := make([]interface{}, 0)
	var ret []int64

	err := lgc.Instance.GetMutilByCondition(ModuleBaseCollectioin, []string{common.BKModuleIDField, common.BKDefaultField}, defaultModuleCond, &result, "ID", 0, 100)
	if nil != err {
		blog.Errorf("get default module ids failed,  error:%s, params:%v, %v", err.Error(), defaultModuleCond, result)
		return ret, errors.New("can not find the module")
	}

	for _, r := range result {
		item := r.(bson.M)
		ID, err := util.GetInt64ByInterface(item[common.BKModuleIDField])
		if nil != err {
			return ret, errors.New("can not find the module")
		}
		ret = append(ret, ID)
	}
	if 0 == len(ret) {
		return ret, errors.New("can not find the module")
	}

	return ret, nil
}

//GetModuleIDsByHostID get module id by hostid
func (lgc *Logics) GetModuleIDsByHostID(moduleCond interface{}) ([]int64, error) {
	result := make([]interface{}, 0)
	var ret []int64

	err := lgc.Instance.GetMutilByCondition(common.BKTableNameModuleHostConfig, []string{common.BKModuleIDField}, moduleCond, &result, "", 0, common.BKNoLimit)
	if nil != err {
		blog.Error("get moudle id by host id failed, error: %s", err.Error())
		return ret, errors.New("can not find the module that host belongs to")
	}
	for _, r := range result {
		item := r.(bson.M)
		ID, getErr := util.GetInt64ByInterface(item[common.BKModuleIDField])
		if nil != getErr || ID == 0 {
			return ret, errors.New("can not find the module")
		}
		ret = append(ret, ID)
	}
	return ret, err
}

//GetResourcePoolApp get resource pool app
func (lgc *Logics) GetResourcePoolApp(ownerID int64) (int64, error) {
	params := common.KvMap{common.BKOwnerIDField: ownerID, common.BKDefaultField: 1}
	result := make(map[string]interface{})
	err := lgc.Instance.GetOneByCondition(ApplicationBaseCollection, []string{common.BKAppIDField}, params, &result)
	if nil != err {
		blog.Errorf("get resource pool app failed,  error:%", err.Error())
		return 0, errors.New("get resource pool app failed")
	}
	appID, err := util.GetInt64ByInterface(result[common.BKAppIDField])
	if err != nil {
		return 0, err
	}
	if 0 == appID {
		blog.Error("get resource pool app failed, can not find the app")
		return 0, errors.New("can not find resource pool app")
	}

	return appID, nil
}

//check if host belong to empty module
func (lgc *Logics) CheckHostInIDle(appID, emptyModuleID int64, hostIDs []int64) ([]int64, []int64, error) {

	conds := common.KvMap{common.BKHostIDField: bson.M{common.BKDBIN: hostIDs}}
	result := make([]interface{}, 0)

	err := lgc.Instance.GetMutilByCondition(common.BKTableNameModuleHostConfig, []string{common.BKHostIDField, common.BKModuleIDField, common.BKAppIDField}, conds, &result, "", 0, common.BKNoLimit)
	if nil != err {
		return nil, nil, fmt.Errorf("get relation between host and module failed, err: %v", err)
	}
	var errHostIDs, faultHostIDs []int64

	mapHost := make(map[int64]int64, 0)
	for _, item := range result {
		row := item.(bson.M)
		moduleID, getErr := util.GetInt64ByInterface(row[common.BKModuleIDField])
		if nil != getErr {
			blog.Errorf("invalid module id: %v", row[common.BKModuleIDField])
			continue
		}
		hostID, getErr := util.GetInt64ByInterface(row[common.BKHostIDField])
		if nil != getErr {
			blog.Errorf("invalid host id :%v", row[common.BKHostIDField])
			continue
		}
		rowAppID, getErr := util.GetInt64ByInterface(row[common.BKAppIDField])
		if nil != getErr {
			blog.Errorf("invalid app id: %v", row[common.BKAppIDField])
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

func (lgc *Logics) GetIDleModuleID(appID int64) (int64, error) {
	cond := common.KvMap{common.BKDefaultField: common.DefaultResModuleFlag, common.BKAppIDField: appID}
	result := make(map[string]interface{}, 0)
	err := lgc.Instance.GetOneByCondition(ModuleBaseCollectioin, []string{common.BKModuleIDField}, cond, &result)
	if nil != err {
		return 0, fmt.Errorf("can not find module, err:　%v", err)
	}

	ID, err := util.GetInt64ByInterface(result[common.BKModuleIDField])
	if nil != err {
		return ID, fmt.Errorf("invalid module id, err:　%v", err)
	}

	return ID, nil
}
