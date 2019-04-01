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
	"context"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	redis "gopkg.in/redis.v5"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
)

type Logics struct {
	Instance dal.RDB
	Cache    *redis.Client
	*backbone.Engine
}

const (
	ModuleBaseCollectioin     = "cc_ModuleBase"
	ModuleHostCollection      = "cc_ModuleHostConfig"
	ApplicationBaseCollection = "cc_ApplicationBase"
)

//DelSingleHostModuleRelation delete single host module relation
func (lgc *Logics) DelSingleHostModuleRelation(ctx context.Context, ec *eventclient.EventContext, hostID, moduleID, appID int64, ownerID string) (bool, error) {

	hostFieldArr := []string{common.BKHostInnerIPField}
	hostResult := make(map[string]interface{}, 0)
	errHost := lgc.GetObjectByID(ctx, common.BKInnerObjIDHost, hostFieldArr, hostID, &hostResult, common.BKHostIDField)
	if errHost != nil {
		blog.Errorf("delete single host relation failed, host: %v, err: %v", hostID, errHost)
		return false, errHost
	}

	moduleFieldArr := []string{common.BKModuleNameField}
	var moduleResult interface{}
	errModule := lgc.GetObjectByID(ctx, common.BKInnerObjIDModule, moduleFieldArr, moduleID, &moduleResult, common.BKModuleNameField)
	if errModule != nil {
		blog.Errorf("delete single host relation, but get module failed,  moduleID:%d, error:%s,", moduleID, errModule.Error())
		return false, errModule
	}

	delCondition := common.KvMap{common.BKAppIDField: appID, common.BKHostIDField: hostID, common.BKModuleIDField: moduleID}
	delCondition = util.SetModOwner(delCondition, ownerID)
	num, numError := lgc.Instance.Table(common.BKTableNameModuleHostConfig).Find(delCondition).Count(ctx) //.GetCntByCondition(ModuleHostCollection, delCondition)
	if numError != nil {
		blog.Errorf("delete single host relation, but get module host relation failed, err: %v", numError)
		return false, numError
	}

	if num == 0 {
		return true, nil
	}

	// retrieve original datas
	origindatas := make([]map[string]interface{}, 0)
	getErr := lgc.Instance.Table(common.BKTableNameModuleHostConfig).Find(delCondition).All(ctx, &origindatas)
	if getErr != nil {
		blog.Errorf("delete single host relation, retrieve original data error:%v", getErr)
		return false, getErr
	}

	delErr := lgc.Instance.Table(common.BKTableNameModuleHostConfig).Delete(ctx, delCondition) //.DelByCondition(ModuleHostCollection, delCondition)
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

// AddSingleHostModuleRelation add single host module relation
func (lgc *Logics) AddSingleHostModuleRelation(ctx context.Context, ec *eventclient.EventContext, hostID, moduleID, appID int64, ownerID string) (bool, error) {
	hostFieldArr := []string{common.BKHostInnerIPField}
	hostResult := make(map[string]interface{})
	errHost := lgc.GetObjectByID(ctx, common.BKInnerObjIDHost, hostFieldArr, hostID, &hostResult, common.BKHostIDField)
	if errHost != nil {
		blog.Errorf("add single host module relation, but get host error:%s", errHost.Error())
		return false, errHost
	}

	moduleFieldArr := []string{common.BKModuleNameField, common.BKSetIDField}
	moduleResult := make(map[string]interface{})
	errModule := lgc.GetObjectByID(ctx, common.BKInnerObjIDModule, moduleFieldArr, moduleID, &moduleResult, common.BKModuleIDField)
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
		blog.Errorf("add single host module relation, get module error:not find module width ModuleID: %d", moduleID)
		return false, errors.New("can not find it's module")
	}

	moduleHostConfig := common.KvMap{common.BKAppIDField: appID, common.BKHostIDField: hostID, common.BKModuleIDField: moduleID}
	num, numError := lgc.Instance.Table(common.BKTableNameModuleHostConfig).Find(moduleHostConfig).Count(ctx)
	if numError != nil {
		blog.Errorf("add single host module relation, get module host relation error: %v", numError)
		return false, numError
	}

	if num > 0 {
		return true, nil
	}

	moduleHostConfig[common.BKSetIDField] = setID
	moduleHostConfig = util.SetModOwner(moduleHostConfig, ownerID)
	err = lgc.Instance.Table(common.BKTableNameModuleHostConfig).Insert(ctx, moduleHostConfig) //.Insert(common.BKTableNameModuleHostConfig, moduleHostConfig)
	if err != nil {
		blog.Errorf("add single host module relation, add module host relation error: %v", err)
		return false, err
	}

	err = ec.InsertEvent(metadata.EventTypeRelation, "moduletransfer", metadata.EventActionCreate, moduleHostConfig, nil)
	if err != nil {
		blog.Errorf("add single host module relation, create event error:%v", err)
	}

	return true, nil
}

// GetDefaultModuleIDs get default module ids
func (lgc *Logics) GetDefaultModuleIDs(ctx context.Context, appID int64) ([]int64, error) {
	defaultModuleCond := make(map[string]interface{})
	defaultModuleCond[common.BKDefaultField] = common.KvMap{common.BKDBIN: []int64{int64(common.DefaultFaultModuleFlag), int64(common.DefaultResModuleFlag)}}
	defaultModuleCond[common.BKAppIDField] = appID
	result := make([]interface{}, 0)
	var ret []int64

	err := lgc.Instance.Table(common.BKTableNameBaseModule).Find(defaultModuleCond).Fields(common.BKModuleIDField, common.BKDefaultField).All(ctx, &result)
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

// GetModuleIDsByHostID get module id by hostid
func (lgc *Logics) GetModuleIDsByHostID(ctx context.Context, moduleCond interface{}) ([]int64, error) {
	result := make([]metadata.ModuleHost, 0)
	var ret []int64

	err := lgc.Instance.Table(common.BKTableNameModuleHostConfig).Find(moduleCond).Fields(common.BKModuleIDField).All(ctx, &result)
	if nil != err {
		blog.Errorf("get module id by host id failed, error: %s", err.Error())
		return ret, errors.New("can not find the module that host belongs to")
	}
	for _, r := range result {
		ret = append(ret, r.ModuleID)
	}
	return ret, err
}

//GetHostIDModuleIDMapsByHostID get hsot id and module id by hostid return map[hostid]moduleid
func (lgc *Logics) GetHostIDModuleIDMapsByHostID(ctx context.Context, moduleCond interface{}, header http.Header) (map[int64][]int64, error) {
	result := make([]metadata.ModuleHost, 0)
	ret := make(map[int64][]int64, 0)
	rid := util.GetHTTPCCRequestID(header)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	fileds := []string{common.BKModuleIDField, common.BKSetIDField, common.BKAppIDField, common.BKHostIDField}
	err := lgc.Instance.Table(common.BKTableNameModuleHostConfig).Find(moduleCond).Fields(fileds...).All(ctx, &result)
	if nil != err {
		blog.Errorf("get moudle id by host id failed, error: %s,rid:%s", err.Error(), rid)
		return ret, defErr.Error(common.CCErrCommDBSelectFailed)
	}

	for _, r := range result {
		ret[r.HostID] = append(ret[r.HostID], r.ModuleID)
	}
	return ret, nil
}

//GetResourcePoolApp get resource pool app
func (lgc *Logics) GetResourcePoolApp(ctx context.Context, ownerID int64) (int64, error) {
	params := common.KvMap{common.BKOwnerIDField: ownerID, common.BKDefaultField: 1}
	result := make(map[string]interface{}, 0)
	err := lgc.Instance.Table(common.BKTableNameBaseApp).Find(params).Fields(common.BKAppIDField).One(ctx, &result)
	if nil != err {
		blog.Errorf("get resource pool app failed,  error:%s", err.Error())
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
func (lgc *Logics) CheckHostInIDle(ctx context.Context, appID, emptyModuleID int64, hostIDs []int64) ([]int64, []int64, error) {

	conds := common.KvMap{common.BKHostIDField: bson.M{common.BKDBIN: hostIDs}}
	result := make([]metadata.ModuleHost, 0)

	err := lgc.Instance.Table(common.BKTableNameModuleHostConfig).Find(conds).Fields(common.BKHostIDField, common.BKModuleIDField, common.BKAppIDField).All(ctx, &result)
	if nil != err {
		return nil, nil, fmt.Errorf("get relation between host and module failed, err: %v", err)
	}
	var errHostIDs, faultHostIDs []int64

	mapHost := make(map[int64]int64, 0)
	for _, item := range result {
		//host not belong to this biz
		if item.AppID != appID {
			faultHostIDs = append(faultHostIDs, item.HostID)
		}
		//host belong to this biz, but not in idle module
		if item.ModuleID != emptyModuleID && item.AppID == appID {
			_, ok := mapHost[item.HostID]
			if !ok {
				errHostIDs = append(errHostIDs, item.HostID)
				mapHost[item.HostID] = item.HostID
			}
		}

	}

	return errHostIDs, faultHostIDs, err
}

func (lgc *Logics) GetIDleModuleID(ctx context.Context, appID int64) (int64, error) {
	cond := common.KvMap{common.BKDefaultField: common.DefaultResModuleFlag, common.BKAppIDField: appID}
	result := make(map[string]interface{}, 0)
	err := lgc.Instance.Table(common.BKTableNameBaseModule).Find(cond).Fields(common.BKModuleIDField).One(ctx, &result)
	if nil != err {
		return 0, fmt.Errorf("can not find module, err:　%v", err)
	}

	ID, err := util.GetInt64ByInterface(result[common.BKModuleIDField])
	if nil != err {
		return ID, fmt.Errorf("invalid module id, err:　%v", err)
	}

	return ID, nil
}

// GetSetAndModuleMapByModuleID return module id  and set id relation, return map[module id] set id
func (lgc *Logics) GetSetAndModuleMapByModuleID(ctx context.Context, appID int64, moduleID []int64, condMapStr mapstr.MapStr, header http.Header) (map[int64]int64, error) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(appID)
	cond.Field(common.BKModuleIDField).In(moduleID)
	if condMapStr == nil {
		condMapStr = mapstr.New()
	}
	condMapStr.Merge(cond.ToMapStr())
	rid := util.GetHTTPCCRequestID(header)
	dbResult := make([]struct {
		SetID    int64 `bson:"bk_set_id"`
		ModuleID int64 `bson:"bk_module_id"`
	}, 0)
	fields := []string{common.BKModuleIDField, common.BKSetIDField}
	err := lgc.Instance.Table(common.BKTableNameBaseModule).Find(condMapStr).Fields(fields...).All(ctx, &dbResult)
	if nil != err {
		blog.Errorf("GetSetAndModuleMapByModuleID query db error. condition:%#v, rid:%s", condMapStr, rid)
		return nil, defErr.Error(common.CCErrCommDBSelectFailed)
	}
	result := make(map[int64]int64, 0)
	for _, row := range dbResult {
		result[row.ModuleID] = row.SetID
	}

	return result, nil
}

// TransferHostToDefaultModuleConfig transfer host to default module config
func (lgc *Logics) TransferHostToDefaultModuleConfig(ctx context.Context, input *metadata.TransferHostToDefaultModuleConfig, header http.Header) error {
	rid := util.GetHTTPCCRequestID(header)
	ownerID := util.GetOwnerID(header)

	moduleCond := condition.CreateCondition()
	moduleCond.Field(common.BKAppIDField).Eq(input.ApplicationID)
	moduleCond.Field(common.BKHostIDField).In(input.HostID)
	hostIDModuleIDMap, err := lgc.GetHostIDModuleIDMapsByHostID(ctx, moduleCond.ToMapStr(), header)
	if nil != err {
		blog.Errorf("TransferHostToDefaultModuleConfig  GetHostIDModuleIDMapsByHostID , input:%#v, condition%#v, err: %v,rid:%s", input, moduleCond.ToMapStr(), err, rid)
		return err
	}
	ec := eventclient.NewEventContextByReq(header, lgc.Cache)
	for hostID, moduleIDArr := range hostIDModuleIDMap {
		for _, moduleID := range moduleIDArr {
			_, err := lgc.DelSingleHostModuleRelation(ctx, ec, hostID, moduleID, input.ApplicationID, ownerID)
			if nil != err {
				blog.Errorf("TransferHostToDefaultModuleConfig  DelSingleHostModuleRelation , input:%#v, hostID:%v,moduleID:%v, err: %v,rid:%s", input, hostID, moduleID, err, rid)
				return err
			}
		}
	}
	for _, hostID := range input.HostID {
		_, err := lgc.AddSingleHostModuleRelation(ctx, ec, hostID, input.ModuleID, input.ApplicationID, ownerID)
		if nil != err {
			blog.Errorf("TransferHostToDefaultModuleConfig  AddSingleHostModuleRelation , input:%#v, hostID:%v, err: %v,rid:%s", input, hostID, err, rid)
			return err
		}
	}
	return nil
}
