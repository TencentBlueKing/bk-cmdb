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

//GetSetIDByCond get setid by cond
func GetSetIDByCond(req *restful.Request, objURL string, cond []interface{}) ([]int, error) {
	setIDArr := make([]int, 0)
	condition := make(map[string]interface{})
	condition["fields"] = common.BKSetIDField
	condition["sort"] = common.BKSetIDField
	condition["start"] = 0
	condition["limit"] = 0
	condc := make(map[string]interface{})
	parse.ParseCommonParams(cond, condc)
	condition["condition"] = condc
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/set/search"
	blog.Infof("GetSetIDByCond url :%s content:%s", url, string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetsetIDByCond return :%s", string(reply))
	if err != nil {
		return setIDArr, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	setData := output["data"]
	setResult, ok := setData.(map[string]interface{})
	if !ok {
		return setIDArr, nil
	}
	setInfo, ok := setResult["info"].([]interface{})
	if !ok {
		return setIDArr, nil
	}
	for _, i := range setInfo {
		set := i.(map[string]interface{})
		setID, _ := set[common.BKSetIDField].(json.Number).Int64()
		setIDArr = append(setIDArr, int(setID))
	}
	return setIDArr, nil
}

//GetSetMapByCond get setmap by cond
func GetSetMapByCond(req *restful.Request, fields string, objURL string, cond interface{}) (map[int]interface{}, error) {
	setMap := make(map[int]interface{})
	condition := make(map[string]interface{})
	condition["fields"] = fields
	condition["sort"] = common.BKModuleIDField
	condition["start"] = 0
	condition["limit"] = 0
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/set/search"
	blog.Info("GetSetMapByCond url :%s", url)
	blog.Info("GetSetMapByCond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetSetMapByCond return :%s", string(reply))
	if err != nil {
		return setMap, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	setData := output["data"]
	setResult := setData.(map[string]interface{})
	setInfo := setResult["info"].([]interface{})
	for _, i := range setInfo {
		set := i.(map[string]interface{})
		setID, _ := set[common.BKSetIDField].(json.Number).Int64()
		setMap[int(setID)] = i
	}
	return setMap, nil
}
