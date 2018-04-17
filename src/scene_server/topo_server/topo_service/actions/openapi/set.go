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
	"configcenter/src/common/util"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

var set = &setAction{}

type setAction struct {
	base.BaseAction
}

func init() {
	// init actionDeleteSetHost  所有的调用都要改
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/openapi/set/multi/{appid}", Params: nil, Handler: set.UpdateMultiSet})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/openapi/set/multi/{appid}", Params: nil, Handler: set.DeleteMultiSet})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/openapi/set/setHost/{appid}", Params: nil, Handler: set.DeleteSetHost})

	// set cc interface
	set.CreateAction()
}

// UpdateMultiSet 更新一个或多个Set，更新多个Set时候名称无效
func (cli *setAction) UpdateMultiSet(req *restful.Request, resp *restful.Response) {
	defErr := m.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	m.CallResponseEx(func() (int, interface{}, error) {
		blog.Debug("UpdateMultiSet start!")

		appID, err := strconv.Atoi(req.PathParameter("appid"))
		if nil != err {
			blog.Error("convert appid to int error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}

		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)

		}

		blog.Debug("UpdateMultiSet http body data: %s", value)

		input := make(map[string]interface{})
		err = json.Unmarshal(value, &input)
		if nil != err {
			blog.Error("unmarshal json error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}

		setIDArr := input[common.BKSetIDField]

		condition := make(map[string]interface{})
		condition[common.BKAppIDField] = appID
		condition[common.BKSetIDField] = map[string]interface{}{
			"$in": setIDArr,
		}

		param := make(map[string]interface{})
		param["condition"] = condition
		param["data"] = input["data"]

		uURL := cli.CC.ObjCtrl() + "/object/v1/insts/set"

		paramJson, err := json.Marshal(param)
		if nil != err {
			blog.Error("Marshal json error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)

		}

		res, err := httpcli.ReqHttp(req, uURL, common.HTTPUpdate, []byte(paramJson))
		blog.Error("uUrl:%v", uURL)
		blog.Error("res:%v", res)

		if nil != err {
			blog.Error("request ctrl error:%v", err)
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoSetUpdateFailed)

		}

		// deal result
		var rst api.BKAPIRsp
		if jserr := json.Unmarshal([]byte(res), &rst); nil == jserr {
			// cli.Response(&rst, resp)
			return http.StatusOK, rst.Data, nil
		} else {
			blog.Error("unmarshal the json failed, error: %v", jserr)
		}
		return http.StatusOK, nil, nil
	}, resp)
}

// DeleteMultiSet 删除一个或多个Set
func (cli *setAction) DeleteMultiSet(req *restful.Request, resp *restful.Response) {
	blog.Debug("DeleteMultiSet start !")
	appID, err := strconv.Atoi(req.PathParameter("appid"))
	if nil != err {
		blog.Error("convert appid to int error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	blog.Debug("DeleteMultiSet appid: %s", appID)

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}

	blog.Debug("DeleteMultiSet http body data: %s", value)

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)

	if nil != err {
		blog.Error("unmarshal json error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	setIDStr := input[common.BKSetIDField]
	setIDStrArr := strings.Split(setIDStr.(string), ",")
	hostIDArr := make([]int, 0)
	moduleIDArr := make([]int, 0)
	blog.Error("setIdArr:%v", setIDStrArr)
	setIDArr, err := util.SliceStrToInt(setIDStrArr)
	//判断集群下是否有主机
	rstOk, rstErr := hasHost(req, cli.CC.HostCtrl(), map[string][]int{
		"ApplicationID": []int{appID},
		"SetID":         setIDArr,
	})

	if nil != rstErr {
		blog.Error("failed to check set wether it has hosts, error info is %s", rstErr.Error())
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}
	if !rstOk {
		msg := "failed to delete set, because of it has some hosts"
		blog.Error(msg)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, fmt.Sprintf("%s", msg), resp)
		return
	}

	configs, err := getConfigByCond(req, cli.CC.HostCtrl(), map[string]interface{}{
		common.BKSetIDField: setIDArr,
	})
	condition := make(map[string]interface{})
	condition[common.BKSetIDField] = map[string]interface{}{
		"$in": setIDArr,
	}
	moduleIDArr, err = getModuleByCond(req, cli.CC.ObjCtrl(), condition)
	blog.Debug("configs:%v", configs)
	for _, config := range configs {
		hostIDArr = append(hostIDArr, config[common.BKHostIDField])
	}

	////删除主机
	//
	//paramHost := make(map[string]interface{})
	//paramHost["HostID"] = map[string]interface{}{
	//	"$in": hostIdArr,
	//}
	//err = deleteObj(req, paramHost, "host", cli.CC.ObjCtrl())
	//if nil != err {
	//	blog.Error("delete host error %v", err)
	//	return
	//}
	//删除模块

	paramModule := make(map[string]interface{})
	paramModule[common.BKModuleIDField] = map[string]interface{}{
		"$in": moduleIDArr,
	}
	err = deleteObj(req, paramModule, "module", cli.CC.ObjCtrl())
	if nil != err {
		blog.Error("delete module error %v", err)
		return
	}

	// 删除Set
	paramSet := make(map[string]interface{})
	paramSet[common.BKAppIDField] = appID
	paramSet[common.BKSetIDField] = map[string]interface{}{
		"$in": setIDArr,
	}

	err = deleteObj(req, paramSet, "set", cli.CC.ObjCtrl())
	if nil != err {
		blog.Error("delete set error %v", err)
		return
	}
	// deal result
	var rst api.BKAPIRsp
	cli.Response(&rst, resp)
	return
}

func deleteObj(req *restful.Request, param map[string]interface{}, objType, objCtrl string) error {

	uURL := objCtrl + "/object/v1/insts/" + objType
	blog.Debug("uURL%v", uURL)
	inputJson, err := json.Marshal(param)
	blog.Debug("inputJson%v", string(inputJson))

	if nil != err {
		blog.Error("Marshal json error:%v", err)
		return err
	}

	res, err := httpcli.ReqHttp(req, uURL, common.HTTPDelete, []byte(inputJson))
	blog.Debug("del res:%v", res)

	if nil != err {
		blog.Error("request ctrl error:%v", err)
		return err
	}
	if "" != res {
		// deal result
		var rst api.APIRsp
		if jserr := json.Unmarshal([]byte(res), &rst); nil == jserr {
			if rst.Result {
				return nil
			} else {
				return errors.New(rst.Message.(string))
			}

		} else {
			blog.Error("unmarshal the json failed, error information is %v", jserr)
			return jserr
		}
	} else {
		return nil
	}

}

// DeleteSetHost 删除set下所有主机，只是在cc_HostModuleConfig删除对应的关系，其他的操作没有
func (cli *setAction) DeleteSetHost(req *restful.Request, resp *restful.Response) {
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	cli.CallResponseEx(func() (int, interface{}, error) {
		blog.Debug("DeleteSetHost start !")
		appID, err := strconv.Atoi(req.PathParameter("appid"))
		if nil != err {
			blog.Error("convert appid to int error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}

		blog.Debug("DeleteSetHost appid: %s", appID)

		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		blog.Debug("DeleteSetHost http body data: %s", value)

		input := make(map[string]interface{})
		err = json.Unmarshal(value, &input)
		if nil != err {
			blog.Error("unmarshal json error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}

		setIDArr := input[common.BKSetIDField]

		param := make(map[string]interface{})
		param[common.BKAppIDField] = appID
		param[common.BKSetIDField] = map[string]interface{}{
			"$in": setIDArr,
		}

		uUrl := cli.CC.ObjCtrl() + "/object/v1/openapi/set/delhost"
		blog.Debug("uUrl%v", uUrl)
		inputJson, err := json.Marshal(param)
		blog.Debug("inputJson%v", string(inputJson))

		if nil != err {
			blog.Error("Marshal json error:%v", err)
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		res, err := httpcli.ReqHttp(req, uUrl, common.HTTPDelete, []byte(inputJson))
		blog.Debug("del res:%v", res)
		if nil != err {
			blog.Error("request ctrl error:%v", err)
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommHTTPDoRequestFailed)

		}
		blog.Debug("res:%v", res)
		//err = delSetConfigHost(param)
		var rst api.BKAPIRsp
		if "not found" == fmt.Sprintf("%v", err) {
			// cli.Response(&rst, resp)
			return http.StatusOK, rst.Data, nil
		}
		if nil != err {
			blog.Error("delSetConfigHost error:%v", err)
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoSetDeleteFailed)
		}

		// deal result

		// cli.Response(&rst, resp)
		return http.StatusOK, rst.Data, nil
	}, resp)
}

//helper
func delSetConfigHost(params map[string]interface{}) error {
	blog.Debug("params:%v", params)
	err := set.CC.InstCli.DelByCondition("cc_ModuleHostConfig", params)
	if err != nil {
		blog.Error("fail to delSetConfigHost: %v", err)
		return err
	}
	return nil
}

func getConfigByCond(req *restful.Request, hostURL string, cond interface{}) ([]map[string]int, error) {
	configArr := make([]map[string]int, 0)
	bodyContent, _ := json.Marshal(cond)
	url := hostURL + "/host/v1/meta/hosts/module/config/search"
	blog.Info("Get ModuleHostConfig url :%s", url)
	blog.Info("Get ModuleHostConfig content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("Get ModuleHostConfig content :%s", string(reply))
	if err != nil {
		return configArr, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, err := js.Map()
	configData := output["data"]
	configInfo, _ := configData.([]interface{})
	for _, mh := range configInfo {
		celh := mh.(map[string]interface{})

		hostID, _ := celh["HostID"].(json.Number).Int64()
		setID, _ := celh["SetID"].(json.Number).Int64()
		moduleID, _ := celh["ModuleID"].(json.Number).Int64()
		appID, _ := celh["ApplicationID"].(json.Number).Int64()
		data := make(map[string]int)
		data["ApplicationID"] = int(appID)
		data["SetID"] = int(setID)
		data["ModuleID"] = int(moduleID)
		data["HostID"] = int(hostID)
		configArr = append(configArr, data)
	}
	return configArr, nil
}

func getModuleByCond(req *restful.Request, objCtrl string, cond interface{}) ([]int, error) {
	searchParams := make(map[string]interface{})
	searchParams["condition"] = cond
	searchParams["fields"] = ""

	sURL := objCtrl + "/object/v1/insts/module/search"
	inputJson, err := json.Marshal(searchParams)
	if nil != err {
		blog.Error("inputJson :%v", err)

	}
	moduleRes, err := httpcli.ReqHttp(req, sURL, "POST", []byte(inputJson))
	blog.Debug("moduleRes:%v", moduleRes)
	if nil != err {
		blog.Error("get module :%v", err)
	}
	js, _ := simplejson.NewJson([]byte(moduleRes))
	moduleMap, _ := js.Map()
	blog.Debug("moduleMap:%v", moduleMap)
	moduleArr := make([]int, 0)
	for _, module := range moduleMap["data"].(map[string]interface{})["info"].([]interface{}) {
		moduleInt64, _ := module.(map[string]interface{})["ModuleID"].(json.Number).Int64()
		moduleArr = append(moduleArr, int(moduleInt64))
	}
	return moduleArr, nil
}

func hasHost(req *restful.Request, hostCtr string, condition map[string][]int) (bool, error) {

	url := fmt.Sprintf("%s/host/v1/meta/hosts/module/config/search", hostCtr)
	inputJSON, _ := json.Marshal(condition)
	rst, rstErr := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(inputJSON))
	if nil != rstErr {
		return false, rstErr
	}
	blog.Debug("get host(%s) config:%s input:%s", url, rst, inputJSON)
	js, err := simplejson.NewJson([]byte(rst))
	if nil != err {
		return false, err
	}
	rstData, _ := js.Map()
	if subData, ok := rstData["data"]; ok {
		if info, infoOk := subData.([]interface{}); infoOk {
			if len(info) > 0 {
				return false, nil
			}
		} else if nil == subData {
			return true, nil
		} else {
			return false, fmt.Errorf("the data is not array, the kind is %s", reflect.TypeOf(subData))
		}
	} else {
		return false, fmt.Errorf("not found the data in result")
	}
	return true, nil
}
