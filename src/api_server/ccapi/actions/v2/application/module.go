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
	"configcenter/src/api_server/ccapi/actions/v2"
	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"strconv"
	"strings"
)

var module *moduleAction = &moduleAction{}

type moduleAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Module/getmodules", Params: nil, Handler: module.GetModulesByApp, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "module/editmodule", Params: nil, Handler: module.UpdateModule, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "module/addModule", Params: nil, Handler: module.AddModule, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "module/delModule", Params: nil, Handler: module.DeleteModule, FilterHandler: nil, Version: v2.APIVersion})

	// set cc api interface
	module.CreateAction()
}

// GetModulesByApp 查询业务下的所有模块
func (cli *moduleAction) GetModulesByApp(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetModulesByApp start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("GetModulesByApp error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	if len(formData["ApplicationID"]) == 0 || formData["ApplicationID"][0] == "" {
		blog.Error("GetModulesByApp error: ApplicationID is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}

	appID := formData["ApplicationID"][0]

	//set empty to get all fields
	param := map[string]interface{}{
		"fields":    []string{},
		"condition": make(map[string]interface{}),
		"page": map[string]interface{}{
			"start": 0,
			"limit": 0,
		},
	}
	paramJson, _ := json.Marshal(param)

	url := fmt.Sprintf("%s/topo/v1/openapi/module/searchByApp/%s", cli.CC.TopoAPI(), appID)
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))

	if err != nil {
		blog.Error("GetModulesByApp url:%s, params:%s, error:%v", url, string(paramJson), err)
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

// UpdateModule: 更新一个或多个Module，更新多个Module时候名称无效
func (cli *moduleAction) UpdateModule(req *restful.Request, resp *restful.Response) {
	blog.Debug("updateModule start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("updateModule error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("updateModule http body data: %s", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "ModuleID"})
	if !res {
		blog.Error("updateModule error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	reqParam := make(map[string]interface{})
	reqData := make(map[string]interface{})

	moduleIDStrArr := strings.Split(formData["ModuleID"][0], ",")
	moduleIDArr, err := utils.SliceStrToInt(moduleIDStrArr)
	if nil != err {
		blog.Error("updateModule error:%v", err)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, err.Error()).Error(), resp)
		return
	}

	appID := formData["ApplicationID"][0]

	reqParam[common.BKModuleIDField] = moduleIDArr
	moduleName := formData.Get("ModuleName")
	Operator := formData.Get("Operator")
	BakOperator := formData.Get("BakOperator")
	ModuleType := formData.Get("ModuleType")
	if "" != Operator {
		reqData[common.BKOperatorField] = Operator
	}
	if "" != Operator {
		reqData[common.BKBakOperatorField] = BakOperator
	}
	if "" != Operator {
		if ModuleType == "1" {
			ModuleType = "普通"
		} else if ModuleType == "2" {
			ModuleType = "数据库"
		} else {
			msg := "模块类型不正确"
			blog.Error("updateModule error:%v", msg)
			converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
			return
		}

		reqData[common.BKModuleTypeField] = ModuleType
	}
	//仅当ID只有一个时候，才能更新ModuleName
	if len(moduleIDArr) == 1 && len(formData["ModuleName"]) > 0 {
		reqData[common.BKModuleNameField] = moduleName
	} else {
		msg := "一次只能更新一个模块"
		blog.Debug("updateModule error:%v", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	if len(moduleName) > 24 {
		msg := "模块名长度不能大于24个字节"
		blog.Debug("updateModule error:%v", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	reqParam["data"] = reqData

	reqParamJson, err := json.Marshal(reqParam)
	if nil != err {
		blog.Error("updateModule error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}

	url := fmt.Sprintf("%s/topo/v1/openapi/module/multi/%s", cli.CC.TopoAPI(), appID)
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, []byte(reqParamJson))
	blog.Debug("rspV3:%v", rspV3)
	if err != nil {
		blog.Error("updateModule url:%s, params:%s, error:%v", url, string(reqParamJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2([]byte(rspV3), resp)
}

// AddModule 新增模块
func (cli *moduleAction) AddModule(req *restful.Request, resp *restful.Response) {
	blog.Debug("addModule start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("addModule error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("addModule http body data: %s", formData)

	res, msg := utils.ValidateFormData(formData, []string{
		"ApplicationID",
		"SetID",
		"ModuleName",
		"Operator",
		"BakOperator",
		"ModuleType",
	})
	if !res {
		blog.Error("addModule error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	reqParam := make(map[string]interface{})
	moduleName := formData.Get("ModuleName")
	appId, _ := strconv.Atoi(formData.Get("ApplicationID"))
	setId, _ := strconv.Atoi(formData.Get("SetID"))
	Operator := formData.Get("Operator")
	BakOperator := formData.Get("BakOperator")
	moduleType := formData.Get("ModuleType")

	moduleTypeMap := map[string]string{"1": "普通", "2": "数据库"}
	if "1" != moduleType && "2" != moduleType {
		if moduleType == "普通" || moduleType == "数据库" {
			reqParam[common.BKModuleTypeField] = moduleType
		} else {
			msg = "bk_module_type 不正确"
			blog.Error("addModule error: %s", msg)
			converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
			return
		}
	} else {
		reqParam[common.BKModuleTypeField] = moduleTypeMap[moduleType]
	}

	if "1" != moduleType && "2" != moduleType {
		msg = "bk_module_type 不正确"
		blog.Error("addModule error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	reqParam[common.BKModuleNameField] = moduleName
	reqParam[common.BKAppIDField] = appId
	reqParam[common.BKSetIDField] = setId
	reqParam[common.BKOperatorField] = Operator
	reqParam[common.BKBakOperatorField] = BakOperator

	reqParamJson, err := json.Marshal(reqParam)
	if nil != err {
		blog.Error("addModule error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}

	url := fmt.Sprintf("%s/topo/v1/openapi/module/multi", cli.CC.TopoAPI())
	rsp_v3, err := httpcli.ReqHttp(req, url, common.HTTPCreate, []byte(reqParamJson))
	if err != nil {
		blog.Error("addModule url:%s, params:%s, error:%v", url, string(reqParamJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	rspV3Map := make(map[string]interface{})

	err = json.Unmarshal([]byte(rsp_v3), &rspV3Map)
	if !rspV3Map["result"].(bool) {
		msg = fmt.Sprintf("%s", rspV3Map[common.HTTPBKAPIErrorMessage])
		blog.Error("CreatePlats error:%s", msg)
		converter.RespFailV2(common.CCErrTopoModuleCreateFailed, msg, resp)
		return
	}

	rspDataV3Map := rspV3Map["data"].(map[string]interface{})
	resData := make(map[string]interface{})
	resData["ModuleID"] = rspDataV3Map[common.BKModuleIDField]
	resData["success"] = true
	converter.RespSuccessV2(resData, resp)
}

// DeleteModule 删除模块
func (cli *moduleAction) DeleteModule(req *restful.Request, resp *restful.Response) {
	blog.Debug("DeleteModule start!")

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("DeleteModule error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("DeleteModule http body data: %s", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "ModuleID"})
	if !res {
		blog.Error("DeleteModule error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	reqParam := make(map[string]interface{})
	moduleIdStrArr := strings.Split(formData["ModuleID"][0], ",")
	moduleIdArr, err := utils.SliceStrToInt(moduleIdStrArr)
	if nil != err {
		blog.Error("DeleteModule error:%v", err)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, err.Error()).Error(), resp)
		return
	}

	appId := formData["ApplicationID"][0]

	reqParam[common.BKModuleIDField] = moduleIdArr

	reqParamJson, err := json.Marshal(reqParam)
	if nil != err {
		blog.Error("DeleteModule error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}

	url := fmt.Sprintf("%s/topo/v1/openapi/module/multi/%s", cli.CC.TopoAPI(), appId)
	rsp_v3, err := httpcli.ReqHttp(req, url, common.HTTPDelete, []byte(reqParamJson))
	blog.Debug("rsp_v3:%v", rsp_v3)
	if err != nil {
		blog.Error("DeleteModule url:%s, params:%s, error:%v", url, string(reqParamJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	converter.RespCommonResV2([]byte(rsp_v3), resp)
}
