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
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"

	v2 "configcenter/src/api_server/ccapi/actions/v2/phpapi"
	logics "configcenter/src/api_server/ccapi/logics/v2"
	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/defs"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
)

var app *appAction = &appAction{}

type appAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "App/getapplist", Params: nil, Handler: app.GetAppList, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "app/getapplist", Params: nil, Handler: app.GetAppList, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "App/getAppByID", Params: nil, Handler: app.GetAppByID, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "app/getAppByID", Params: nil, Handler: app.GetAppByID, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "app/getappbyid", Params: nil, Handler: app.GetAppByID, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "App/getappbyuin", Params: nil, Handler: app.GetAppByUIN, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "app/getappbyuin", Params: nil, Handler: app.GetAppByUIN, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "User/getUserRoleApp", Params: nil, Handler: app.GetUserRoleApp, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "user/getUserRoleApp", Params: nil, Handler: app.GetUserRoleApp, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "TopSetModule/getappsetmoduletreebyappid", Params: nil, Handler: app.GetAppSetModuleTreeByAppId, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "app/addApp", Params: nil, Handler: app.AddApp, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "app/deleteApp", Params: nil, Handler: app.DeleteApp, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "app/editApp", Params: nil, Handler: app.EditApp, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "App/getHostAppByCompanyId", Params: nil, Handler: app.GetHostAppByCompanyId, FilterHandler: nil, Version: v2.APIVersion})

	// set cc api interface
	app.CreateAction()
}

// GetAppList search application, return all application
func (cli *appAction) GetAppList(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetAppList start !")

	//set empty to get all fields
	param := map[string]interface{}{
		"condition": nil,
	}

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	paramJson, _ := json.Marshal(param)
	url := fmt.Sprintf("%s/topo/v1/app/search/"+common.BKDefaultOwnerID, cli.CC.TopoAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	if err != nil {
		blog.Error("GetAppList url:%s, params:%s error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForAppList(rspV3)
	if err != nil {
		blog.Error("convert app res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

// GetAppByID  get application by application id
func (cli *appAction) GetAppByID(req *restful.Request, resp *restful.Response) {
	blog.Info("GetAppByID start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("GetAppByID error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("GetAppByID http body data:%v", formData)

	if len(formData["ApplicationID"]) == 0 {
		blog.Error("GetAppByID error: ApplicationID is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}

	appIDsStr := strings.Split(formData["ApplicationID"][0], ",")

	appIDs := make([]int, 0)
	for _, idStr := range appIDsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			msg := fmt.Sprintf("convert str to int error:%v", err)
			blog.Error("convert str to int error:%s", msg)
			converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
			return
		}
		appIDs = append(appIDs, id)
	}

	// build v3 parameters
	param := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKAppIDField: map[string]interface{}{
				"$in": appIDs,
			},
		},
	}

	paramJson, _ := json.Marshal(param)
	url := fmt.Sprintf("%s/topo/v1/app/search/"+common.BKDefaultOwnerID, cli.CC.TopoAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	if err != nil {
		blog.Error("GetAppByID  url:%s, params:%s error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForAppList(rspV3)
	if err != nil {
		blog.Error("convert app res to v2 error:%v, reply:%s", err, string(rspV3))
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

// GetAppByUIN  get application with username permission
func (cli *appAction) GetAppByUIN(req *restful.Request, resp *restful.Response) {
	blog.Debug("getAppByUIN start!")

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("getAppByUIN error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("getAppByUIN http body data:%v", formData)

	if len(formData["userName"]) == 0 || formData["userName"][0] == "" {
		blog.Error("getAppByUIN error: userName is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedSet, "userName").Error(), resp)
		return
	}

	userName := formData["userName"][0]
	v3Username := strings.Trim(converter.DecorateUserName(userName), ",")
	orCondition := []map[string]interface{}{
		map[string]interface{}{
			"bk_biz_maintainer": map[string]interface{}{
				common.BKDBLIKE: fmt.Sprintf("^%s,|,%s,|,%s$|^%s$", v3Username, v3Username, v3Username, v3Username),
			},
		},
	}
	filterOnly := false //is application maintainer
	if len(formData["filterOnly"]) > 0 && formData["filterOnly"][0] != "" {
		if "true" == formData["filterOnly"][0] {
			filterOnly = true
		}
	}

	if filterOnly == false {
		orCondition = append(orCondition, map[string]interface{}{
			"bk_biz_productor": map[string]interface{}{
				common.BKDBLIKE: fmt.Sprintf("^%s,|,%s,|,%s$|^%s$", v3Username, v3Username, v3Username, v3Username),
			},
		})
	}

	//build v3 parameters
	param := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKDBOR: orCondition,
		},
	}

	paramJson, err := json.Marshal(param)

	if err != nil {
		blog.Error("GetAppByUIN error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}

	url := fmt.Sprintf("%s/topo/v1/app/search/"+common.BKDefaultOwnerID, cli.CC.TopoAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	if err != nil {
		blog.Error("GetAppByUIN url:%s, params:%s, error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForAppList(rspV3)
	if err != nil {
		blog.Error("convert res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	blog.Debug("GetAppByUIN success, data length: %d", len(resDataV2.([]map[string]interface{})))
	converter.RespSuccessV2(resDataV2, resp)
}

// GetUserRoleApp  search application with user role, multiple user role split comma
func (cli *appAction) GetUserRoleApp(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetUserRoleApp start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("GetUserRoleApp error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	res, msg := utils.ValidateFormData(formData, []string{"uin", "roleList"})

	if !res {
		blog.Error("GetUserRoleApp error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	userName := formData["uin"][0]
	roleArr := strings.Split(formData["roleList"][0], ",")

	v3Username := converter.DecorateUserName(userName)
	roleOrCondition := make([]map[string]interface{}, 0)
	for _, roleStr := range roleArr {
		roleStrV3, ok := defs.RoleMap[roleStr]
		if !ok {
			continue
		}

		roleOrCondition = append(roleOrCondition, map[string]interface{}{
			roleStrV3: map[string]interface{}{
				common.BKDBLIKE: fmt.Sprintf("^%s,|,%s,|,%s$|^%s$", v3Username, v3Username, v3Username, v3Username),
			},
		})
	}

	//build v3 parameters
	param := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKDBOR: roleOrCondition,
		},
	}

	blog.Debug("GetUserRoleApp v3 param: %v", param)

	paramJson, err := json.Marshal(param)
	if err != nil {
		blog.Error("GetUserRoleApp error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}
	url := fmt.Sprintf("%s/topo/v1/app/search/"+common.BKDefaultOwnerID, cli.CC.TopoAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	if err != nil {
		blog.Error("GetUserRoleApp url:%, params:%s, error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForRoleApp(rspV3, converter.DecorateUserName(userName), roleArr)
	if err != nil {
		blog.Error("convert res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

// GetAppSetModuleTreeByAppId    get topology  by application id, default ownerid = 0
func (cli *appAction) GetAppSetModuleTreeByAppId(req *restful.Request, resp *restful.Response) {
	blog.Debug("getAppSetModuleTreeByAppId start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("getAppSetModuleTreeByAppID error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form
	blog.Debug("getAppSetModuleTreeByAppID http body data: %v", formData)

	if len(formData["ApplicationID"]) == 0 {
		blog.Error("getAppSetModuleTreeByAppID error: ApplicationID is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	intAppID, err := util.GetInt64ByInterface(formData["ApplicationID"][0])

	conds := common.KvMap{}
	if 0 < len(formData["ModuleType"]) {
		// 1 normal module, 2 database module
		conds[common.BKModuleTypeField] = formData["ModuleType"][0]
	}

	if nil != err {
		blog.Error("getAppSetModuleTreeByAppID   appID:%v, error:%v", formData["ApplicationID"][0], err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	topo, errCode := logics.GetAppTopo(req, cli.CC, common.BKDefaultOwnerID, intAppID, conds)
	if 0 != errCode {
		converter.RespFailV2(errCode, defErr.Error(errCode).Error(), resp)
		return
	}

	if nil != topo {
		logics.SetModuleHostCount([]map[string]interface{}{topo}, req, cli.CC)
	}
	converter.RespSuccessV2(topo, resp)
}

// addApp 新增业务
func (cli *appAction) AddApp(req *restful.Request, resp *restful.Response) {
	blog.Debug("AddApp start")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("AddApp error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Debug("add app formData:%v", formData)
	res, msg := utils.ValidateFormData(formData, []string{
		"ApplicationName",
		"Maintainers",
		//"CompanyName",
		//"Level",
		"Creator",
		"LifeCycle",
	})
	if !res {
		blog.Error("add app error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	if len(formData.Get("ApplicationName")) > 32 {
		converter.RespFailV2(common.CCErrAPIServerV2APPNameLenErr, defErr.Error(common.CCErrAPIServerV2APPNameLenErr).Error(), resp)
		return
	}
	lifeCycle := formData.Get("LifeCycle")

	if lifeCycle == "1" || "测试中" == lifeCycle {
		lifeCycle = "1"
	} else if lifeCycle == "2" || "已上线" == lifeCycle {
		lifeCycle = "2"
	} else if lifeCycle == "3" || "停运" == lifeCycle {
		lifeCycle = "3"
	} else {
		//msg := "生成周期字段值不合法"
		converter.RespFailV2(common.CCErrCommParamsIsInvalid, defErr.Errorf(common.CCErrCommParamsIsInvalid, "LifeCycle").Error(), resp)
		return
	}

	param := make(common.KvMap)
	param[common.BKAppNameField] = formData.Get("ApplicationName")
	param[common.BKMaintainersField] = formData.Get("Maintainers")
	param[common.BKLanguageField] = "1"

	timeZone := formData.Get("TimeZone")
	if "" != timeZone {
		param[common.BKTimeZoneField] = timeZone
	} else {
		param[common.BKTimeZoneField] = "Asia/Shanghai"
	}

	//param[common.CreatorField] = formData.GetUserConfig("Creator")
	param[common.BKLifeCycleField] = lifeCycle
	param[common.BKProductPMField] = formData.Get("ProductPm")
	param[common.BKDeveloperField] = formData.Get("Developer")
	param[common.BKTesterField] = formData.Get("Tester")
	param[common.BKOperatorField] = formData.Get("Operator")

	blog.Debug("AddApp v3 param data1: %v", param)
	//填充v3版本需要的参数
	param, err = logics.AutoInputV3Field(param, common.BKInnerObjIDApp, app.CC.TopoAPI(), req.Request.Header)

	paramJson, err := json.Marshal(param)
	blog.Debug("AddApp v3 param data: %v", param)
	url := fmt.Sprintf("%s/topo/v1/app/%s", app.CC.TopoAPI(), common.BKDefaultOwnerID)
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPCreate, []byte(paramJson))

	blog.Infof("create app url:%s, params:%s", url, paramJson)
	if nil != err {
		blog.Errorf("create app url:%s, params:%s, error:%s", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	rspV3Map := make(map[string]interface{})
	err = json.Unmarshal([]byte(rspV3), &rspV3Map)
	if !rspV3Map["result"].(bool) {
		blog.Error("Create App url:%s error:%s", url, rspV3Map[common.HTTPBKAPIErrorMessage])
		errMsg, _ := rspV3Map[common.HTTPBKAPIErrorMessage].(string)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error()+" :"+errMsg, resp)
		return
	}

	rspDataV3Map := rspV3Map["data"].(map[string]interface{})
	resData := make(map[string]interface{})
	resData["appId"] = rspDataV3Map[common.BKAppIDField]
	resData["success"] = true
	converter.RespSuccessV2(resData, resp)
}

//deleteApp: 删除业务
func (cli *appAction) DeleteApp(req *restful.Request, resp *restful.Response) {
	blog.Debug("DeletedApp start")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("DeletedApp error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	appId := formData.Get("ApplicationID")
	if "" == appId {
		blog.Errorf("DeletedApp error ApplicationID empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}
	//param := make(common.KvMap)
	//param["appId"] = appId
	//paramJson, err := json.Marshal(param)
	//blog.Debug("DeletedApp http body data: %v", formData)
	//url := fmt.Sprintf("%s/topo/v1/openapi/app/deleteApp", app.CC.TopoAPI())
	url := fmt.Sprintf("%s/topo/v1/app/%s/%s", app.CC.TopoAPI(), common.BKDefaultOwnerID, appId)
	blog.Debug("url:%v", url)
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPDelete, nil)
	if nil != err {
		blog.Errorf("DeleteApp url:%s, error:%s", url, err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	blog.Debug("rspV3:%v", rspV3)
	rspV3Map := make(map[string]interface{})
	err = json.Unmarshal([]byte(rspV3), &rspV3Map)
	if nil != err {
		blog.Error("DeletedApp Unmarshal json error:%v, rspV3:%s", err, rspV3)
		return
	}

	converter.RespCommonResV2([]byte(rspV3), resp)
}

//editApp: edit app
func (cli *appAction) EditApp(req *restful.Request, resp *restful.Response) {
	blog.Debug("EditApp start")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("EditApp error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Debug("edit app formData:%v", formData)
	res, msg := utils.ValidateFormData(formData, []string{
		"ApplicationID",
	})
	if !res {
		blog.Error("add app error: %s", msg)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}
	appName := formData.Get("ApplicationName")

	if "" != appName && len(appName) > 32 {
		converter.RespFailV2(common.CCErrAPIServerV2APPNameLenErr, defErr.Error(common.CCErrAPIServerV2APPNameLenErr).Error(), resp)
		return
	}
	LifeCycle := formData.Get("LifeCycle")
	if "" != LifeCycle {
		lifeMap := map[string]bool{"测试中": true, "已上线": true, "停运": true}
		if !lifeMap[LifeCycle] {
			if LifeCycle == "1" {
				LifeCycle = "测试中"
			} else if LifeCycle == "2" {
				LifeCycle = "已上线"
			} else if LifeCycle == "3" {
				LifeCycle = "停运"
			} else {
				converter.RespFailV2(common.CCErrCommParamsIsInvalid, defErr.Errorf(common.CCErrCommParamsIsInvalid, "LifeCycle").Error(), resp)
				return
			}
		}
	}

	param := make(common.KvMap)
	appId := formData.Get("ApplicationID")
	//param["ApplicationID"] = appId
	param[common.BKAppNameField] = formData.Get("ApplicationName")
	if formData.Get("LifeCycle") != "" {
		param[common.BKLifeCycleField] = LifeCycle
	}
	if formData.Get("Maintainers") != "" {
		param[common.BKMaintainersField] = formData.Get("Maintainers")
	}
	if formData.Get("Creator") != "" {
		param[common.CreatorField] = formData.Get("Creator")
	}
	if formData.Get("ProductPm") != "" {
		param[common.BKProductPMField] = formData.Get("ProductPm")
	}
	if formData.Get("Developer") != "" {
		param[common.BKDeveloperField] = formData.Get("Developer")
	}
	if formData.Get("Tester") != "" {
		param[common.BKTesterField] = formData.Get("Tester")
	}
	if formData.Get("Operator") != "" {
		param[common.BKOperatorField] = formData.Get("Operator")
	}
	blog.Error("edit_app param:%v", param)

	//param, err = logics.AutoInputV3Field(param, common.BKInnerObjIDApp, app.CC.TopoAPI(), req.Request.Header)

	paramJson, err := json.Marshal(param)
	blog.Debug("edit app http body data: %v", formData)
	//url := fmt.Sprintf("%s/topo/v1/openapi/app/editApp", app.CC.TopoAPI())
	url := fmt.Sprintf("%s/topo/v1/app/%s/%s", app.CC.TopoAPI(), common.BKDefaultOwnerID, appId)
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, []byte(paramJson))
	if nil != err {
		blog.Errorf("EditApp url:%s, params:%s, error:%s", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	blog.Debug("rspV3:%v", rspV3)
	rspV3Map := make(map[string]interface{})
	err = json.Unmarshal([]byte(rspV3), &rspV3Map)
	if nil != err {
		blog.Error("edit app Unmarshal json error:%v, rspV3:%s", err, rspV3)
		return
	}
	converter.RespCommonResV2([]byte(rspV3), resp)
}

//GetHostAppByCompanyId:    get application by companyid and host ip
func (cli *appAction) GetHostAppByCompanyId(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetHostAppByCompanyId start")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("GetHostAppByCompanyId error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Debug("GetHostAppByCompanyId formData:%v", formData)
	res, msg := utils.ValidateFormData(formData, []string{
		"companyId",
		"ip",
		"platId",
	})
	if !res {
		blog.Error("GetHostAppByCompanyId error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	param := make(common.KvMap)
	param[common.BKOwnerIDField] = formData.Get("companyId")
	param["ip"] = formData.Get("ip")
	param[common.BKCloudIDField] = formData.Get("platId")

	paramJson, err := json.Marshal(param)
	blog.Debug("GetHostAppByCompanyId http body data: %v", formData)
	url := fmt.Sprintf("%s/host/v1/openapi/host/getHostAppByCompanyId", app.CC.HostAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	if nil != err {
		blog.Errorf("GetHostAppByCommpanyId url:%s, params:%s, error:%s", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	blog.Debug("rspV3:%v", rspV3)
	resDataV2, err := converter.ResToV2ForHostDataList(rspV3)
	if err != nil {
		blog.Error("convert res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}
	blog.Debug("GetHostAppByCompanyId success, data length: %d", len(resDataV2))
	converter.RespSuccessV2(resDataV2, resp)
}
