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

package application

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"

	v2 "configcenter/src/api_server/ccapi/actions/v2/phpapi"
	logics "configcenter/src/api_server/ccapi/logics/v2"
	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
)

var set *setAction = &setAction{}

type setAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Set/getsetsbyproperty", Params: nil, Handler: set.GetSets, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Set/getsetproperty", Params: nil, Handler: set.Getsetproperty, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Set/getmodulesbyproperty", Params: nil, Handler: module.GetModulesByProperty, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "set/getmodulesbyproperty", Params: nil, Handler: module.GetModulesByProperty, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "set/addset", Params: nil, Handler: set.AddSet, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "set/updateset", Params: nil, Handler: set.UpdateSet, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "set/updateSetServiceStatus", Params: nil, Handler: set.UpdateSetServiceStatus, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "set/delset", Params: nil, Handler: set.DelSet, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "set/delSetHost", Params: nil, Handler: set.DelSetHost, FilterHandler: nil, Version: v2.APIVersion})

	// set cc api interface
	set.CreateAction()
}

// GetSets: 根据Set属性获取Set
func (cli *setAction) GetSets(req *restful.Request, resp *restful.Response) {
	blog.Debug("getSets start!")

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("getSets error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("getSets http body data:%v", formData)

	if len(formData["ApplicationID"]) == 0 {
		blog.Error("getSets error: ApplicationID is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}

	condition := make(map[string]interface{})

	appId := formData.Get("ApplicationID")
	condition[common.BKAppIDField] = appId

	blog.Error("len:%v", len(formData["SetServiceStatus"]))
	if "" != formData.Get("SetServiceStatus") {
		formStatus := formData.Get("SetServiceStatus")

		//服务状态，包含0：关闭，1：开启，默认为1
		strStatus := "1"
		if "0" == formStatus {
			strStatus = "0"
		}
		condition[common.BKSetStatusField] = strStatus
	}

	if "" != formData.Get("SetEnviType") {
		formEnv := formData.Get("SetEnviType")
		//env 1：测试 2：体验 3：正式，默认为3
		strEnv := "3"
		switch formEnv {
		case "1":
			strEnv = "1"
		case "2":
			strEnv = "2"
		}
		condition[common.BKSetEnvField] = strEnv
	}

	//set cc v3 request params
	reqParam := map[string]interface{}{
		"condition": condition,
		"fields":    []string{},
		"page": map[string]interface{}{
			"start": 0,
			"limit": 0,
			"sort":  "",
		},
	}

	reqParamJson, _ := json.Marshal(reqParam)
	blog.Debug("getSets reqParamJson:%v", string(reqParamJson))
	url := fmt.Sprintf("%s/topo/v1/set/search/0/%s", cli.CC.TopoAPI(), appId)

	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(reqParamJson))
	if err != nil {
		blog.Error("getSets url:%s, params:%s, error:%v", url, string(reqParamJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForSetList(rspV3)
	if err != nil {
		blog.Error("convert set res to v2 error:%v, rspV3:%v", err, rspV3)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

// addSet 新增单个Set
func (cli *setAction) AddSet(req *restful.Request, resp *restful.Response) {
	blog.Debug("addSet start!")

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("addSet error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("addSet http body data: %s", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "SetName"})
	if !res {
		blog.Error("addSet error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	reqParam := make(map[string]interface{})
	if len(formData["SetName"][0]) > 24 {
		msg = defLang.Language("apiv2_set_name_lt_24")
		blog.Error("add set error:%v", msg)
		converter.RespFailV2(common.CCErrAPIServerV2SetNameLenErr, defErr.Errorf(common.CCErrAPIServerV2SetNameLenErr, msg).Error(), resp)
		return
	}

	reqParam[common.BKSetNameField] = formData["SetName"][0]
	//reqParam[common.BKSetStatusField], reqParam[common.BKSetEnvField] = getSetServiceAndEnv(formData)

	//默认OwnerID：0
	reqParam[common.BKOwnerIDField] = common.BKDefaultOwnerID

	appID := formData["ApplicationID"][0]
	reqParam[common.BKInstParentStr], err = strconv.Atoi(appID)
	if nil != err {
		blog.Error("AddSet convert appid to int error:%v", err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	if len(formData["properties"]) > 0 {
		propertiesJson := formData["properties"][0]
		propertiesMap := make(map[string]interface{})
		err := json.Unmarshal([]byte(propertiesJson), &propertiesMap)
		if nil != err {
			blog.Error("addSet error:%v", err)
			converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
			return
		}
		for k, v := range propertiesMap {
			if k == "EnviType" {
				vInt, err := util.GetIntByInterface(v)
				if nil != err {
					blog.Error("addSet error:%v", err)
					converter.RespFailV2(common.CCErrCommParamsIsInvalid, defErr.Errorf(common.CCErrCommParamsIsInvalid, "EnviType").Error(), resp)
					return
				}
				vStr := strconv.Itoa(vInt)

				env := getAddSetEnv(vStr)
				reqParam[common.BKSetEnvField] = env
				continue
			}
			if k == "ServiceStatus" {
				vInt, err := util.GetIntByInterface(v)
				if nil != err {
					blog.Error("addSet error:%v", err)
					converter.RespFailV2(common.CCErrCommParamsIsInvalid, defErr.Errorf(common.CCErrCommParamsIsInvalid, "ServiceStatus").Error(), resp)
					return
				}
				vStr := strconv.Itoa(vInt)
				status := getAddSetService(vStr)
				reqParam[common.BKSetStatusField] = status
				continue
			}
			if k == "Description" {
				reqParam[common.BKSetDescField] = v
				continue
			}
			if k == "Capacity" {
				reqParam[common.BKSetCapacityField], err = util.GetIntByInterface(v)
				if nil != err {
					blog.Error("add set GetIntByInterface error:%v", err)
					converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsIsInvalid, "Capacity").Error(), resp)
					return
				}
				continue
			}
			reqParam[k] = v
		}
	}
	delete(reqParam, "ChnName")

	topoLevel, err := logics.CheckAppTopoIsThreeLevel(req, cli.CC)
	if err != nil {
		blog.Error("AddSet CheckAppTopoIsThreeLevel error:%v", err)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, err.Error()).Error(), resp)
		return
	}
	if false == topoLevel {
		blog.Error("AddSet CheckAppTopoIsThreeLevel  mainline topology level not three")
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, "business topology level is more than three, please use the v3 api instead").Error(), resp)
		return
	}
	reqParam, err = logics.AutoInputV3Field(reqParam, common.BKInnerObjIDSet, set.CC.TopoAPI(), req.Request.Header)
	blog.Debug("add set reqParam:%v", reqParam)

	if err != nil {
		blog.Error("AutoInputV3Field error:%v", err)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, err.Error()).Error(), resp)
		return
	}
	reqParamJson, err := json.Marshal(reqParam)
	if nil != err {
		blog.Error("addSet error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}

	blog.Error("add set reqParamJson:%v", string(reqParamJson))
	url := fmt.Sprintf("%s/topo/v1/set/%s", cli.CC.TopoAPI(), appID)
	rsp_v3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(reqParamJson))
	if err != nil {
		blog.Error("addSet url:%s, params:%s, error:%v", url, string(reqParamJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	rspV3Map := make(map[string]interface{})
	err = json.Unmarshal([]byte(rsp_v3), &rspV3Map)
	if err != nil {
		blog.Error("addSet not json url:%s, reply:%s", url, string(rsp_v3))
		converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}
	if !rspV3Map["result"].(bool) {
		msg = fmt.Sprintf("%s", rspV3Map[common.HTTPBKAPIErrorMessage])
		blog.Error("CreatePlats error:%s", msg)
		converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}
	rspDataV3Map, _ := rspV3Map["data"].(map[string]interface{})
	blog.Debug("rsp_v3:%v", rsp_v3)
	converter.RespSuccessV2(rspDataV3Map, resp)
}

// updateSet 更新一个或多个Set，更新多个Set时候名称无效
func (cli *setAction) UpdateSet(req *restful.Request, resp *restful.Response) {
	blog.Debug("updateSet start!")

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("updateSet error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("updateSet http body data: %s", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "SetID"})
	if !res {
		blog.Error("updateSet error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	reqParam := make(map[string]interface{})
	reqData := make(map[string]interface{})

	setIDStrArr := strings.Split(formData["SetID"][0], ",")
	setIDArr, err := utils.SliceStrToInt(setIDStrArr)
	if nil != err {
		blog.Error("updateSet error:%v", err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "SetID").Error(), resp)
		return
	}

	appID := formData["ApplicationID"][0]

	reqParam[common.BKSetIDField] = setIDArr
	setName := formData["SetName"][0]
	//仅当ID只有一个时候，才能更新SetName
	if len(setIDArr) == 1 && len(formData["SetName"]) > 0 {

		reqData[common.BKSetNameField] = setName
	} else {
		msg := "一次只能更新一个集群"
		blog.Error("updateSet error:%v", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	if len(formData["SetName"][0]) > 24 {
		msg = "集群名称长度不能超过24个字节"
		blog.Error("add set error:%v", msg)
		converter.RespFailV2(common.CCErrAPIServerV2SetNameLenErr, defErr.Error(common.CCErrAPIServerV2SetNameLenErr).Error(), resp)
		return
	}
	status, env := getSetServiceAndEnv(formData)
	if len(formData["SetEnviType"]) > 0 {
		reqData[common.BKSetEnvField] = env
	}
	if len(formData["Des"]) > 0 {
		reqData[common.BKSetDescField] = formData["Des"][0]
	}
	if len(formData["Capacity"]) > 0 {
		reqData[common.BKSetCapacityField], err = util.GetIntByInterface(formData["Capacity"][0])
		if nil != err {
			blog.Error("updateSet error:%v", err)
			converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "Capacity").Error(), resp)
			return
		}
	}

	if len(formData["ServiceStatus"]) > 0 {
		reqData[common.BKSetStatusField] = status
	}

	delete(reqData, "ChnName")

	reqParam["data"] = reqData

	reqParamJson, err := json.Marshal(reqParam)
	if nil != err {
		blog.Error("updateSet error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}

	url := fmt.Sprintf("%s/topo/v1/openapi/set/multi/%s", cli.CC.TopoAPI(), appID)
	rsp_v3, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, []byte(reqParamJson))
	if err != nil {
		blog.Error("updateSet url:%s,params:%s ,error:%v", url, string(reqParamJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2([]byte(rsp_v3), resp)
}

// updateSetServiceStatus 修改服务状态, Status 0:关闭，1:开启
func (cli *setAction) UpdateSetServiceStatus(req *restful.Request, resp *restful.Response) {
	blog.Debug("updateSetServiceStatus start!")

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("updateSetServiceStatus error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("updateSetServiceStatus http body data: %s", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "SetID", "Status"})
	if !res {
		blog.Error("updateSetServiceStatus error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	reqParam := make(map[string]interface{})
	reqData := make(map[string]interface{})

	setIDStrArr := strings.Split(formData["SetID"][0], ",")
	setIDArr, err := utils.SliceStrToInt(setIDStrArr)
	if nil != err {
		blog.Error("updateSetServiceStatus error:%v", err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "SetID").Error(), resp)
		return
	}

	appID := formData["ApplicationID"][0]

	reqData[common.BKSetStatusField] = formData["Status"][0]
	// service status  combin 0：关闭，1：开启，默认为1
	switch reqData[common.BKSetStatusField] {
	case "0":
		reqData[common.BKSetStatusField] = "2" //"关闭"
	case "1":
		reqData[common.BKSetStatusField] = "1" //"开放"

	}

	reqParam[common.BKSetIDField] = setIDArr
	reqParam["data"] = reqData

	reqParamJson, err := json.Marshal(reqParam)
	if nil != err {
		blog.Error("updateSetServiceStatus error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}

	url := fmt.Sprintf("%s/topo/v1/openapi/set/multi/%s", cli.CC.TopoAPI(), appID)
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, []byte(reqParamJson))
	if err != nil {
		blog.Error("updateSetServiceStatus url:%s, params:%s, error:%v", url, string(reqParamJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2([]byte(rspV3), resp)
}

// delSet 删除一个或多个Set
func (cli *setAction) DelSet(req *restful.Request, resp *restful.Response) {
	blog.Debug("delSet start!")

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("delSet error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("delSet http body data: %s", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "SetID"})
	if !res {
		blog.Error("delSet error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	appID := formData["ApplicationID"][0]

	reqParam := make(map[string]interface{})
	reqParam[common.BKSetIDField] = formData["SetID"][0]

	reqParamJson, err := json.Marshal(reqParam)
	if nil != err {
		blog.Error("delSet error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}

	url := fmt.Sprintf("%s/topo/v1/openapi/set/multi/%s", cli.CC.TopoAPI(), appID)

	rsp_v3, err := httpcli.ReqHttp(req, url, common.HTTPDelete, []byte(reqParamJson))
	blog.Debug("rsp_v3:%v", rsp_v3)
	if err != nil {
		blog.Error("delSet url:%s, params:%s , error:%v", url, string(reqParamJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2([]byte(rsp_v3), resp)
}

// delSetHost 删除set下所有主机，只是在cc_HostModuleConfig删除对应的关系，其他的操作没有
func (cli *setAction) DelSetHost(req *restful.Request, resp *restful.Response) {
	blog.Debug("delSetHost start!")

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("delSetHost error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("delSetHost http body data: %s", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "SetID"})
	if !res {
		blog.Error("delSetHost error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	appID := formData["ApplicationID"][0]
	setIDStrArr := strings.Split(formData["SetID"][0], ",")
	setIDArr, err := utils.SliceStrToInt(setIDStrArr)
	reqParam := make(map[string]interface{})
	reqParam[common.BKSetIDField] = setIDArr
	reqParamJson, err := json.Marshal(reqParam)
	if nil != err {
		blog.Error("delSetHost error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}

	url := fmt.Sprintf("%s/topo/v1/openapi/set/setHost/%s", cli.CC.TopoAPI(), appID)
	rsp_v3, err := httpcli.ReqHttp(req, url, common.HTTPDelete, []byte(reqParamJson))
	if err != nil {
		blog.Error("delSetHost url:%s, params:%s, error:%v", url, string(reqParamJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2([]byte(rsp_v3), resp)
}

// GetModulesByProperty: 通过属性获取模块
func (cli *moduleAction) GetModulesByProperty(req *restful.Request, resp *restful.Response) {
	blog.Debug("getModulesByProperty start!")

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("getModulesByProperty error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("GetModulesByProperty http body data:%v", formData)

	if len(formData["ApplicationID"]) == 0 {
		blog.Error("getModulesByProperty error: ApplicationID is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	appID := formData["ApplicationID"][0]

	param := make(map[string]interface{})

	if len(formData["SetID"]) > 0 && formData["SetID"][0] != "" {
		setIDArr, sliceErr := utils.SliceStrToInt(strings.Split(formData["SetID"][0], ","))
		if sliceErr != nil {
			blog.Error("SliceStrToInt error:%v", sliceErr)
			converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, sliceErr.Error()).Error(), resp)
			return
		}

		param[common.BKSetIDField] = setIDArr
	}

	if len(formData["SetEnviType"]) > 0 && "" != formData["SetEnviType"][0] && formData["SetID"][0] != "" {
		param[common.BKSetEnvField] = strings.Split(formData["SetEnviType"][0], ",")
	}

	if len(formData["SetServiceStatus"]) > 0 && "" != formData["SetServiceStatus"][0] {
		param[common.BKSetStatusField] = strings.Split(formData["SetServiceStatus"][0], ",")
	}

	paramJson, err := json.Marshal(param)
	if err != nil {
		blog.Error("getModulesByProperty error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}

	url := fmt.Sprintf("%s/topo/v1/openapi/module/searchByProperty/%s", cli.CC.TopoAPI(), appID)
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	if err != nil {
		blog.Error("getModulesByProperty url:%s, params:%s, error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	resDataV2, err := converter.ResToV2ForModuleMapList(rspV3)
	if err != nil {
		blog.Error("convert module res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

// GetAllSets 获取所有set徐行
func (cli *setAction) Getsetproperty(req *restful.Request, resp *restful.Response) {
	resDataV2 := common.KvMap{
		/*"option" : [
		    {
		        "id" : "2",
		        "name" : "体验",
		        "type" : "text",
		        "is_default" : false
		    },
		    {
		        "id" : "3",
		        "name" : "正式",
		        "type" : "text",
		        "is_default" : true
		    },
		    {
		        "id" : "1",
		        "name" : "测试",
		        "type" : "text",
		        "is_default" : false
		    }
		]*/
		"SetEnviType": []common.KvMap{
			common.KvMap{"Property": "1", "value": "1"},
			common.KvMap{"Property": "2", "value": "2"},
			common.KvMap{"Property": "3", "value": "3"},
		},
		/*"option" : [
		    {
		        "id" : "2",
		        "name" : "关闭",
		        "type" : "text",
		        "is_default" : false
		    },
		    {
		        "id" : "1",
		        "name" : "开放",
		        "type" : "text",
		        "is_default" : true
		    }
		]*/
		"SetServiceStatus": []common.KvMap{
			common.KvMap{"Property": "0", "value": "2"},
			common.KvMap{"Property": "1", "value": "1"},
		},
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func getSetService(formData url.Values) string {
	formStatus := ""
	if len(formData["ServiceStatus"]) > 0 {
		formStatus = formData["ServiceStatus"][0]
	}
	//服务状态，包含0：关闭，1：开启，默认为1
	strStatus := "1"
	if "0" == formStatus {
		strStatus = "2"
	}
	return strStatus
}

func getAddSetService(v string) string {
	formStatus := v

	//服务状态，包含0：关闭，1：开启，默认为1
	strStatus := "1"
	if "0" == formStatus {
		strStatus = "2"
	}
	return strStatus
}

func getAddSetEnv(v string) string {

	formEnv := v

	//env 1：测试 2：体验 3：正式，默认为3
	strEnv := "3"
	switch formEnv {
	case "1":
		strEnv = "1"
	case "2":
		strEnv = "2"
	}
	return strEnv
}

func getSetServiceAndEnv(formData url.Values) (string, string) {

	formEnv := ""

	if len(formData["SetEnviType"]) > 0 {
		formEnv = formData["SetEnviType"][0]
	}
	//env 1：测试 2：体验 3：正式，默认为3
	strEnv := "3"
	switch formEnv {
	case "1":
		strEnv = "1"
	case "2":
		strEnv = "1"
	}
	strStatus := getSetService(formData)
	return strStatus, strEnv
}
