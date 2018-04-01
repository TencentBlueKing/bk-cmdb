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
	"encoding/json"

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
	blog.Info("GetModuleIDByCond url :%s", url)
	blog.Info("GetModuleIDByCond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetModuleIDByCond return :%s", string(reply))
	if err != nil {
		return moduleIDArr, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	moduleData := output["data"]
	moduleResult := moduleData.(map[string]interface{})
	moduleInfo := moduleResult["info"].([]interface{})
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
