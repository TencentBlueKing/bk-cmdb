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
	moduleResult := moduleData.(map[string]interface{})
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
