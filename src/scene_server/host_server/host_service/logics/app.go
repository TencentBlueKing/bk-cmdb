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
	appParse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"encoding/json"
	"errors"
	"fmt"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

//GetAppIDByCond get appid by cond
func GetAppIDByCond(req *restful.Request, objURL string, cond []interface{}) ([]int, error) {
	appIDArr := make([]int, 0)
	condition := make(map[string]interface{})
	condition["fields"] = common.BKAppIDField
	condition["sort"] = common.BKAppIDField
	condition["start"] = 0
	condition["limit"] = 1000000
	condc := make(map[string]interface{})
	appParse.ParseCommonParams(cond, condc)
	condition["condition"] = condc
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/" + common.BKInnerObjIDApp + "/search"
	blog.Info("GetAppIDByCond url :%s", url)
	blog.Info("GetAppIDByCond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetAppIDByCond return :%s", string(reply))
	if err != nil {
		return appIDArr, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	appData := output["data"]
	appResult := appData.(map[string]interface{})
	appInfo := appResult["info"].([]interface{})
	for _, i := range appInfo {
		app := i.(map[string]interface{})
		appID, _ := app[common.BKAppIDField].(json.Number).Int64()
		appIDArr = append(appIDArr, int(appID))
	}
	return appIDArr, nil
}

//GetAppMapByCond get appmap by cond
func GetAppMapByCond(req *restful.Request, fields string, objURL string, cond interface{}) (map[int]interface{}, error) {
	appMap := make(map[int]interface{})
	condition := make(map[string]interface{})
	condition["fields"] = fields
	condition["sort"] = common.BKAppIDField
	condition["start"] = 0
	condition["limit"] = 0
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/" + common.BKInnerObjIDApp + "/search"
	blog.Info("GetAppMapByCond url :%s", url)
	blog.Info("GetAppMapByCond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetAppMapByCond return :%s", string(reply))
	if err != nil {
		blog.Errorf("GetAppMapByCond params:%s  error:%v", string(bodyContent), err)
		return appMap, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	appData := output["data"]
	appResult := appData.(map[string]interface{})
	appInfo := appResult["info"].([]interface{})
	for _, i := range appInfo {
		app := i.(map[string]interface{})
		appID, _ := app[common.BKAppIDField].(json.Number).Int64()
		appMap[int(appID)] = i
	}
	return appMap, nil
}

//GetSingleApp  get single app
func GetSingleApp(req *restful.Request, objURL string, cond interface{}) (map[string]interface{}, error) {
	condition := make(map[string]interface{})
	condition["sort"] = common.BKAppIDField
	condition["start"] = 0
	condition["limit"] = 1
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/" + common.BKInnerObjIDApp + "/search"
	fmt.Println("GetSingleApp", url, string(bodyContent))

	blog.Info("GetOneApp url :%s", url)
	blog.Info("GetOneApp content :%s", string(bodyContent))

	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	fmt.Println("GetSingleApp", url, string(reply))
	blog.Info("GetOneApp return :%s", string(reply))
	if err != nil {
		blog.Info("GetOneApp return http request error:%s", string(reply))
		return nil, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	appData := output["data"]
	appResult := appData.(map[string]interface{})
	appInfo := appResult["info"].([]interface{})
	for _, i := range appInfo {
		app, _ := i.(map[string]interface{})
		return app, nil
	}
	return nil, nil
}

//GetAppInfo get app info
func GetAppInfo(req *restful.Request, fields string, conditon map[string]interface{}, hostAddr string) (map[string]interface{}, error) {
	//moduleURL := "http://" + cc.ObjCtrl + "/object/v1/insts/module/search"
	URL := hostAddr + "/object/v1/insts/" + common.BKInnerObjIDApp + "/search"
	params := make(map[string]interface{})
	params["condition"] = conditon
	params["sort"] = common.BKAppIDField
	params["start"] = 0
	params["limit"] = 1
	params["fields"] = fields

	blog.Info("get application info  url:%s", URL)
	blog.Info("get application info  url:%v", params)
	isSuccess, errMsg, data := GetHttpResult(req, URL, common.HTTPSelectPost, params)
	if !isSuccess {
		blog.Error("get application info  error, params:%v, error:%s", params, errMsg)
		return nil, errors.New(errMsg)
	}
	dataInterface := data.(map[string]interface{})
	info := dataInterface["info"].([]interface{})
	if 1 != len(info) {
		blog.Error("not application info error, params:%v, error:%s", params, errMsg)
		return nil, errors.New("业务不存在")
	}
	row := info[0].(map[string]interface{})

	if 0 == len(row) {
		blog.Error("not application info error, params:%v, error:%s", params, errMsg)
		return nil, errors.New("业务存在")
	}

	return row, nil
}

//GetDefaultAppID get default biz id
func GetDefaultAppID(req *restful.Request, ownerID, fields, hostAddr string) (int, error) {
	conds := make(map[string]interface{})
	conds[common.BKOwnerIDField] = ownerID
	conds[common.BKDefaultField] = common.DefaultAppFlag
	appinfo, err := GetAppInfo(req, fields, conds, hostAddr)
	if nil != err {
		blog.Errorf("get default app info error:%v", err.Error())
		return 0, err
	}
	return util.GetIntByInterface(appinfo[common.BKAppIDField])
}

//GetDefaultAppID get supplier ID
func GetDefaultAppIDBySupplierID(req *restful.Request, supplierID int, fields, hostAddr string) (int, error) {
	conds := make(map[string]interface{})
	conds[common.BKSupplierIDField] = supplierID
	conds[common.BKDefaultField] = common.DefaultAppFlag
	appinfo, err := GetAppInfo(req, fields, conds, hostAddr)
	if nil != err {
		blog.Errorf("get default app info error:%v", err.Error())
		return 0, err
	}
	return util.GetIntByInterface(appinfo[common.BKAppIDField])
}
