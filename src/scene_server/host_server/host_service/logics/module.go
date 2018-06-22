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
	httpcli "configcenter/src/common/http/httpclient"
	parse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/source_controller/common/commondata"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

//get moduleid by cond
func GetModuleIDByCond(req *restful.Request, objURL string, cond []interface{}) ([]int, error) {
	moduleIDArr := make([]int, 0)
	condition := make(map[string]interface{})
	condition["fields"] = common.BKModuleIDField
	condition["sort"] = common.BKModuleIDField
	condition["start"] = 0
	condition["limit"] = 0
	condc := make(map[string]interface{})
	parse.ParseCommonParams(cond, condc)
	condition["condition"] = condc
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/module/search"
	blog.Infof("GetModuleIDByCond url :%s content:%s", url, string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetModuleIDByCond return :%s", string(reply))
	if err != nil {
		return moduleIDArr, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	moduleData := output["data"]
	moduleResult, ok := moduleData.(map[string]interface{})
	if !ok {
		return moduleIDArr, nil
	}
	moduleInfo, ok := moduleResult["info"].([]interface{})
	if !ok {
		return moduleIDArr, nil
	}
	for _, i := range moduleInfo {
		module := i.(map[string]interface{})
		moduleID, _ := module[common.BKModuleIDField].(json.Number).Int64()
		moduleIDArr = append(moduleIDArr, int(moduleID))
	}
	return moduleIDArr, nil
}

//get modulemap by cond
func GetModuleMapByCond(req *restful.Request, fields string, objURL string, cond interface{}) (map[int]interface{}, error) {
	moduleMap := make(map[int]interface{})
	condition := make(map[string]interface{})
	condition["fields"] = fields
	condition["sort"] = common.BKModuleIDField
	condition["start"] = 0
	condition["limit"] = 0
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/module/search"
	blog.Info("GetModuleMapByCond url :%s", url)
	blog.Info("GetModuleMapByCond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetModuleMapByCond return :%s", string(reply))
	if err != nil {
		return moduleMap, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	moduleData := output["data"]
	moduleResult := moduleData.(map[string]interface{})
	moduleInfo := moduleResult["info"].([]interface{})
	for _, i := range moduleInfo {
		module := i.(map[string]interface{})
		moduleID, _ := module[common.BKModuleIDField].(json.Number).Int64()
		moduleMap[int(moduleID)] = i
	}
	return moduleMap, nil
}

//GetModuleByModuleID  get module by module id
func GetModuleByModuleID(req *restful.Request, appID int, moduleID int, hostAddr string) ([]interface{}, error) {
	URL := hostAddr + "/object/v1/insts/module/search"
	params := make(map[string]interface{})

	conditon := make(map[string]interface{})
	conditon[common.BKAppIDField] = appID
	conditon[common.BKModuleIDField] = moduleID
	params["condition"] = conditon
	params["sort"] = common.BKModuleIDField
	params["start"] = 0
	params["limit"] = 1
	params["fields"] = common.BKModuleIDField
	isSuccess, errMsg, data := GetHttpResult(req, URL, common.HTTPSelectPost, params)
	if !isSuccess {
		blog.Error("get idle module error, params:%v, error:%s", params, errMsg)
		return nil, errors.New(errMsg)
	}
	dataStrArry := data.(map[string]interface{})
	dataInfo, ok := dataStrArry["info"].([]interface{})
	if !ok {
		blog.Error("get idle module error, params:%v, error:%s", params, errMsg)
		return nil, errors.New(errMsg)
	}

	return dataInfo, nil
}

//GetSingleModuleID  get single module id
func GetSingleModuleID(req *restful.Request, conds interface{}, hostAddr string) (int, error) {
	//moduleURL := "http://" + cc.ObjCtrl + "/object/v1/insts/module/search"
	url := hostAddr + "/object/v1/insts/module/search"
	params := make(map[string]interface{})

	params["condition"] = conds
	params["sort"] = common.BKModuleIDField
	params["start"] = 0
	params["limit"] = 1
	params["fields"] = common.BKModuleIDField
	isSuccess, errMsg, data := GetHttpResult(req, url, common.HTTPSelectPost, params)
	if !isSuccess {
		blog.Error("get idle module error, params:%v, error:%s", params, errMsg)
		return 0, errors.New(errMsg)
	}
	dataInterface := data.(map[string]interface{})
	info := dataInterface["info"].([]interface{})
	if 1 != len(info) {
		blog.Error("not find module error, params:%v, error:%s", params, errMsg)
		return 0, errors.New("获取集群，返回数据格式错误")
	}
	row := info[0].(map[string]interface{})
	moduleID, _ := util.GetIntByInterface(row[common.BKModuleIDField])

	if 0 == moduleID {
		blog.Error("not find module error, params:%v, error:%s", params, errMsg)
		return 0, errors.New("获取集群信息失败")
	}

	return moduleID, nil
}

// NewHostSyncValidModule
// 1. check module is exist,
// 2. multiple moduleID  Check whether  all the module default 0
func NewHostSyncValidModule(req *restful.Request, appID int, moduleID []int, objAddr string) ([]int, error) {
	if 0 == len(moduleID) {
		return nil, fmt.Errorf("module id number must be > 1")
	}

	conds := common.KvMap{
		common.BKAppIDField:    appID,
		common.BKModuleIDField: common.KvMap{common.BKDBIN: moduleID},
	}

	condition := new(commondata.ObjQueryInput)
	condition.Sort = common.BKModuleIDField
	condition.Start = 0
	condition.Limit = 0
	condition.Condition = conds
	bodyContent, err := json.Marshal(condition)
	if nil != err {
		return nil, fmt.Errorf("query module parameters not json")
	}
	url := objAddr + "/object/v1/insts/module/search"
	blog.Info("NewHostSyncValidModule url :%s", url)
	blog.Info("NewHostSyncValidModule content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("NewHostSyncValidModule return :%s", string(reply))

	js, err := simplejson.NewJson([]byte(reply))
	if nil != err {
		blog.Errorf("NewHostSyncValidModule get module reply not json,  url:%s, params:%s, reply:%s", string(bodyContent), url, reply)
		return nil, fmt.Errorf("get moduel reply not json, reply:%s", reply)
	}

	moduleInfos, err := js.Get("data").Get("info").Array()
	if nil != err {
		blog.Errorf("NewHostSyncValidModule get module reply not foound data.info,  url:%s, params:%s, reply:%s", string(bodyContent), url, reply)
		return nil, fmt.Errorf("get moduel reply not found info from data, reply:%s", reply)
	}

	// only module  and module exist return true
	if 1 == len(moduleID) && 1 == len(moduleInfos) {
		return moduleID, nil
	}

	// use module id is exist
	moduleIDMap := make(map[int64]bool)
	for _, id := range moduleID {
		moduleIDMap[int64(id)] = false
	}

	// multiple module all module default = 0
	for _, module := range moduleInfos {
		moduleMap, ok := module.(map[string]interface{})
		if !ok {
			blog.Errorf("NewHostSyncValidModule item not map[string]interface{},  module:%v", module)
			return nil, fmt.Errorf("item not map[string]interface{},  module:%v", module)
		}
		moduelDefault, err := util.GetInt64ByInterface(moduleMap[common.BKDefaultField])
		if nil != err {
			blog.Errorf("NewHostSyncValidModule item not found default,  module:%v", module)
			return nil, fmt.Errorf("module information not found default,  module:%v", module)
		}
		if 0 != moduelDefault {
			return nil, fmt.Errorf("multiple module cannot appear system module")
		}
		moduleID, err := util.GetInt64ByInterface(moduleMap[common.BKModuleIDField])
		if nil != err {
			blog.Errorf("NewHostSyncValidModule item not found module id,  module:%v", module)
			return nil, fmt.Errorf("module information not found module id,  module:%v", module)
		}
		moduleIDMap[moduleID] = true //module id  exist db

	}

	var dbModuleID []int
	var notExistModuleId []string
	for id, exist := range moduleIDMap {
		if exist {
			dbModuleID = append(dbModuleID, int(id))
		} else {
			notExistModuleId = append(notExistModuleId, strconv.FormatInt(id, 10))
		}
	}
	if 0 < len(notExistModuleId) {
		return nil, fmt.Errorf("module id %s not found", strings.Join(notExistModuleId, ","))
	}

	return dbModuleID, nil
}
