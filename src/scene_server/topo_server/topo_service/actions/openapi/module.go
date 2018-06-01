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
 
package openapi

import (
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

var m *moduleAction = &moduleAction{}

type moduleAction struct {
	base.BaseAction
}

func init() {
	// init action
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/openapi/module/multi/{" + common.BKAppIDField + "}", Params: nil, Handler: m.UpdateMultiModule})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/openapi/module/searchByApp/{" + common.BKAppIDField + "}", Params: nil, Handler: m.SearchModuleByApp})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/openapi/module/searchByProperty/{" + common.BKAppIDField + "}", Params: nil, Handler: m.SearchModuleByProperty})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/openapi/module/multi", Params: nil, Handler: m.AddMultiModule})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/openapi/module/multi/{" + common.BKAppIDField + "}", Params: nil, Handler: m.DeleteMultiModule})

	m.CreateAction()
}

// SearchModuleByApp: 获取业务下的所有模块
func (cli *moduleAction) SearchModuleByApp(req *restful.Request, resp *restful.Response) {
	blog.Debug("SearchModuleByApp")

	appID, err := strconv.Atoi(req.PathParameter(common.BKAppIDField))
	if nil != err {
		blog.Error("convert appid to int error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}

	blog.Debug("SearchModuleByApp http body data: %s", value)

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Error("unmarshal json error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	condition := input["condition"].(map[string]interface{})
	condition[common.BKAppIDField] = appID

	page := input["page"].(map[string]interface{})
	fields := input["fields"].([]interface{})
	strFields := make([]string, 0)
	for _, i := range fields {
		strI := i.(string)
		strFields = append(strFields, strI)
	}

	searchParams := make(map[string]interface{})
	searchParams["condition"] = condition
	searchParams["fields"] = strings.Join(strFields, ",")
	searchParams["start"] = page["start"]
	searchParams["limit"] = page["limit"]
	searchParams["sort"] = page["sort"]

	//search
	sURL := cli.CC.ObjCtrl() + "/object/v1/insts/module/search"
	inputJson, _ := json.Marshal(searchParams)
	moduleRes, err := httpcli.ReqHttp(req, sURL, "POST", []byte(inputJson))
	blog.Debug("search module params: %s", string(inputJson))
	if nil != err {
		blog.Error("search module error: %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Module_QUERY_FAIL, common.CC_Err_Comm_Module_QUERY_FAIL_STR, resp)
		return
	}

	// deal result
	var rst api.BKAPIRsp
	if jserr := json.Unmarshal([]byte(moduleRes), &rst); nil == jserr {
		cli.Response(&rst, resp)
		return
	} else {
		blog.Error("unmarshal the json failed, error information is %v", jserr)
	}
}

// SearchModuleByProperty: 根据属性获取模块列表
func (cli *moduleAction) SearchModuleByProperty(req *restful.Request, resp *restful.Response) {
	blog.Debug("SearchModuleByProperty start!")

	appID, err := strconv.Atoi(req.PathParameter(common.BKAppIDField))
	if nil != err {
		blog.Error("convert appid to int error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}

	blog.Debug("SearchModuleByProperty http body data: %s", value)

	properties := make(map[string]interface{})
	err = json.Unmarshal(value, &properties)
	if nil != err {
		blog.Error("unmarshal json error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	setCondition := make(map[string]interface{})
	for proName, proValue := range properties {
		setCondition[proName] = map[string]interface{}{
			"$in": proValue,
		}
	}

	blog.Debug("getSetIDByCond condition: %v", setCondition)

	setIDArr, err := getSetIDByCond(req, cli.CC.ObjCtrl(), setCondition)
	if nil != err {
		blog.Error("SearchModuleByProperty error: %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Module_QUERY_FAIL, common.CC_Err_Comm_Module_QUERY_FAIL_STR, resp)
		return
	}

	blog.Debug("getSetIDByCond res: %v", setIDArr)

	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appID
	condition[common.BKSetIDField] = map[string]interface{}{
		"$in": setIDArr,
	}

	searchParams := make(map[string]interface{})
	searchParams["condition"] = condition
	searchParams["fields"] = ""
	sURL := cli.CC.ObjCtrl() + "/object/v1/insts/module/search"
	inputJson, _ := json.Marshal(searchParams)
	moduleRes, err := httpcli.ReqHttp(req, sURL, "POST", []byte(inputJson))

	if nil != err {
		blog.Error("SearchModuleByProperty error: %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Module_QUERY_FAIL, common.CC_Err_Comm_Module_QUERY_FAIL_STR, resp)
		return
	}

	// deal result
	var rst api.BKAPIRsp
	if jserr := json.Unmarshal([]byte(moduleRes), &rst); nil == jserr {
		cli.Response(&rst, resp)
		return
	} else {
		blog.Error("unmarshal the json failed, error information is %v", jserr)
	}
	return
}

// UpdateMultiModule 更新一个或多个Module，更新多个更新一个或多个Module时候名称无效
func (cli *moduleAction) UpdateMultiModule(req *restful.Request, resp *restful.Response) {
	blog.Debug("UpdateMultiModule start!")

	appID, err := strconv.Atoi(req.PathParameter(common.BKAppIDField))
	if nil != err {
		blog.Error("convert appid to int error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}

	blog.Debug("UpdateMultiModule http body data: %s", value)

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Error("unmarshal json error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	moduleIDArr := input[common.BKModuleIDField]

	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appID
	condition[common.BKModuleIDField] = map[string]interface{}{
		"$in": moduleIDArr,
	}

	param := make(map[string]interface{})
	param["condition"] = condition
	param["data"] = input["data"]

	moduleName := input["data"].(map[string]interface{})[common.BKModuleNameField]
	moduleCon := map[string]interface{}{
		common.BKModuleNameField: moduleName,
		common.BKModuleIDField: map[string]interface{}{
			"$nin": moduleIDArr,
		},
	}
	blog.Error("edit module moduleCon:%v", moduleCon)

	moduleData, err := GetModuleMapByCond(req, "", cli.CC.ObjCtrl(), moduleCon)
	blog.Error("edit module moduleData:%v", moduleData)
	if err != nil {
		blog.Error("Marshal json error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}
	if len(moduleData) > 0 {
		msg := "模块名已存在"
		blog.Debug("addModule error: %s", msg)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, fmt.Sprintf("%s%s", common.CC_Err_Comm_http_Input_Params_STR, msg), resp)
		return
	}
	uURL := cli.CC.ObjCtrl() + "/object/v1/insts/module"

	paramJson, err := json.Marshal(param)
	if nil != err {
		blog.Error("Marshal json error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	res, err := httpcli.ReqHttp(req, uURL, common.HTTPUpdate, []byte(paramJson))

	if nil != err {
		blog.Error("request ctrl error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Module_Update_FAIL, common.CC_Err_Comm_Module_Update_FAIL_STR, resp)
		return
	}
	blog.Debug("res:%v", res)
	if "" == res {
		msg := "没有找到模块"
		blog.Error("request ctrl error:%v", msg)
		cli.ResponseFailed(common.CC_Err_Comm_Module_Update_FAIL, fmt.Sprintf("%s:%s", common.CC_Err_Comm_Module_Update_FAIL_STR, msg), resp)
		return
	}
	// deal result
	var rst api.BKAPIRsp
	if jserr := json.Unmarshal([]byte(res), &rst); nil == jserr {
		cli.Response(&rst, resp)
		return
	} else {
		blog.Error("unmarshal the json failed, error: %v", jserr)
		cli.ResponseFailed(common.CC_Err_Comm_Module_Update_FAIL, fmt.Sprintf("%s:%v", common.CC_Err_Comm_Module_Update_FAIL_STR, jserr), resp)
		return
	}
}

// AddMultiModule 添加模块
func (cli *moduleAction) AddMultiModule(req *restful.Request, resp *restful.Response) {
	blog.Debug("AddMultiModule start!")

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}

	blog.Debug("AddMultiModule http body data: %s", value)

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Error("unmarshal json error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	appId := input[common.BKAppIDField]
	setId := input[common.BKSetIDField]
	condition := map[string]interface{}{
		common.BKAppIDField: appId,
		common.BKSetIDField: setId,
	}
	setIdArr, err := getSetIDByCond(req, cli.CC.ObjCtrl(), condition)
	if nil != err {
		blog.Debug("add module error getSetIdByCond err:%v", err)
		cli.ResponseFailed(common.CCErrTopoModuleCreateFailed, "create module error", resp)
		return
	}
	if len(setIdArr) == 0 {
		msg := "SetID is not find "
		blog.Debug("add module error getSetIdByCond err:%v", msg)
		cli.ResponseFailed(common.CCErrTopoModuleCreateFailed, msg, resp)
		return
	}
	moduleNameArr := strings.Split(input[common.BKModuleNameField].(string), ",")
	param := make(map[string]interface{})

	for _, moduleName := range moduleNameArr {
		if len(moduleName) > 24 {
			msg := "模块名不能大于24个字节"
			blog.Debug("addModule error: %s", msg)
			cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, fmt.Sprintf("%s%s", common.CC_Err_Comm_http_Input_Params_STR, msg), resp)
			return
		}
		moduleCon := map[string]interface{}{
			common.BKModuleNameField: moduleName,
		}
		moduleData, getErr := GetModuleMapByCond(req, "", cli.CC.ObjCtrl(), moduleCon)
		if getErr != nil {
			blog.Error("Marshal json error:%v", getErr)
			cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
			return
		}
		if len(moduleData) > 0 {
			msg := "模块已存在"
			blog.Debug("addModule error: %s", msg)
			cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, fmt.Sprintf("%s%s", common.CC_Err_Comm_http_Input_Params_STR, msg), resp)
			return
		}
		param[common.BKModuleNameField] = moduleName
		param[common.BKAppIDField] = appId
		param[common.BKSetIDField] = input[common.BKSetIDField]
		param[common.BKOperatorField] = input[common.BKOperatorField]
		param[common.BKBakOperatorField] = input[common.BKBakOperatorField]
		param[common.BKModuleTypeField] = input[common.BKModuleTypeField]
		param[common.BKDefaultField] = 0
		param[common.BKInstParentStr] = input[common.BKSetIDField]
		param[common.BKOwnerIDField] = common.BKDefaultOwnerID
	}

	uUrl := cli.CC.ObjCtrl() + "/object/v1/insts/module"

	paramJson, err := json.Marshal(param)
	if nil != err {
		blog.Error("Marshal json error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	res, err := httpcli.ReqHttp(req, uUrl, common.HTTPCreate, []byte(paramJson))
	if nil != err {
		blog.Error("request ctrl error:%v", err)
		cli.ResponseFailed(common.CCErrTopoModuleCreateFailed, "create module error", resp)
		return
	}

	// deal result
	var rst api.BKAPIRsp
	if jserr := json.Unmarshal([]byte(res), &rst); nil == jserr {
		cli.Response(&rst, resp)
		return
	} else {
		blog.Error("unmarshal the json failed, error: %v", jserr)
	}
}

// DeleteMultiModule 删除模块，模块下存在主机则不删除
func (cli *moduleAction) DeleteMultiModule(req *restful.Request, resp *restful.Response) {
	blog.Debug("DeleteMultiModule start!")

	appId, err := strconv.Atoi(req.PathParameter(common.BKAppIDField))
	if nil != err {
		blog.Error("convert appid to int error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}

	blog.Debug("DeleteMultiModule http body data: %s", value)

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Error("unmarshal json error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	moduleIdArr := input[common.BKModuleIDField]
	configCond := map[string]interface{}{
		common.BKModuleIDField: moduleIdArr,
	}
	configData, err := GetConfigByCond(req, cli.CC.HostCtrl(), configCond)
	blog.Debug("configData:%v", configData)
	if nil != err {
		blog.Error("delete module GetConfigByCond error:%v", err)
		cli.ResponseFailed(common.CCErrTopoModuleDeleteFailed, "delete module error", resp)
		return
	}
	if len(configData) > 0 {
		msg := "module has host ,not delete"
		blog.Error("delete module GetConfigByCond error:%v", msg)
		cli.ResponseFailed(common.CCErrTopoModuleDeleteFailed, msg, resp)
		return
	}

	param := make(map[string]interface{})
	param[common.BKAppIDField] = appId
	param[common.BKModuleIDField] = map[string]interface{}{
		"$in": moduleIdArr,
	}
	blog.Debug(param)
	uUrl := cli.CC.ObjCtrl() + "/object/v1/insts/module"
	paramJson, err := json.Marshal(param)
	if nil != err {
		blog.Error("Marshal json error:%v", err)
		cli.ResponseFailed(common.CCErrTopoModuleDeleteFailed, "delete module error", resp)
		return
	}

	res, err := httpcli.ReqHttp(req, uUrl, common.HTTPDelete, []byte(paramJson))

	if "" == res {
		msg := "not find module"
		blog.Error("request ctrl error:%v", msg)
		cli.ResponseFailed(common.CCErrTopoModuleDeleteFailed, msg, resp)
		return
	}
	if nil != err {
		blog.Error("request ctrl error:%v", err)
		cli.ResponseFailed(common.CCErrTopoModuleDeleteFailed, "delete module error", resp)
		return
	}
	blog.Debug("delete module res%v", res)
	// deal result
	var rst api.BKAPIRsp
	if jserr := json.Unmarshal([]byte(res), &rst); nil == jserr {
		cli.Response(&rst, resp)
		return
	} else {
		blog.Error("unmarshal the json failed, error: %v", jserr)
	}
}

// helpers
func getSetIDByCond(req *restful.Request, objURL string, cond map[string]interface{}) ([]int, error) {
	setIDArr := make([]int, 0)
	condition := make(map[string]interface{})
	condition["fields"] = common.BKSetIDField
	condition["sort"] = common.BKSetIDField
	condition["start"] = 0
	condition["limit"] = 0
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/set/search"
	reply, err := httpcli.ReqHttp(req, url, "POST", []byte(bodyContent))
	if err != nil {
		return setIDArr, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	setData := output["data"]
	setResult := setData.(map[string]interface{})
	setInfo := setResult["info"].([]interface{})
	for _, i := range setInfo {
		set := i.(map[string]interface{})
		setID, _ := set[common.BKSetIDField].(json.Number).Int64()
		setIDArr = append(setIDArr, int(setID))
	}
	return setIDArr, nil
}

//get config by condition
func GetConfigByCond(req *restful.Request, hostUrl string, cond interface{}) ([]map[string]int, error) {
	configArr := make([]map[string]int, 0)
	bodyContent, _ := json.Marshal(cond)
	fmt.Println(string(bodyContent))
	url := hostUrl + "/host/v1/meta/hosts/module/config/search"
	blog.Info("Get ModuleHostConfig url :%s", url)
	blog.Info("Get ModuleHostConfig content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("Get ModuleHostConfig content :%s", string(reply))
	fmt.Println("get moduleHostConfig", url, reply)
	if err != nil {
		return configArr, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, err := js.Map()
	configData := output["data"]
	configInfo, _ := configData.([]interface{})
	for _, mh := range configInfo {
		celh := mh.(map[string]interface{})

		hostId, _ := celh[common.BKHostIDField].(json.Number).Int64()
		setId, _ := celh[common.BKSetIDField].(json.Number).Int64()
		moduleId, _ := celh[common.BKModuleIDField].(json.Number).Int64()
		appId, _ := celh[common.BKAppIDField].(json.Number).Int64()
		data := make(map[string]int)
		data[common.BKAppIDField] = int(appId)
		data[common.BKSetIDField] = int(setId)
		data[common.BKModuleIDField] = int(moduleId)
		data[common.BKHostIDField] = int(hostId)
		configArr = append(configArr, data)
	}
	return configArr, nil
}

//get modulemap by cond
func GetModuleMapByCond(req *restful.Request, fields string, objUrl string, cond interface{}) (map[int]interface{}, error) {
	moduleMap := make(map[int]interface{})
	condition := make(map[string]interface{})
	condition["fields"] = fields
	condition["sort"] = common.BKModuleIDField
	condition["start"] = 0
	condition["limit"] = 0
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := objUrl + "/object/v1/insts/module/search"
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
		moduleId, _ := module[common.BKModuleIDField].(json.Number).Int64()
		moduleMap[int(moduleId)] = i
	}
	return moduleMap, nil
}
