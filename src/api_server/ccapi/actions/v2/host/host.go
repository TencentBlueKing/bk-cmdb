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

package host

import (
	"configcenter/src/common"

	logics "configcenter/src/api_server/ccapi/logics/v2"
	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/api_server/ccapi/logics/v2/common/defs"

	"configcenter/src/api_server/ccapi/actions/v2"
	"errors"

	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

var host *hostAction = &hostAction{}

type hostAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Host/gethostlistbyip", Params: nil, Handler: host.GetHostListByIP, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/gethostlistbyip", Params: nil, Handler: host.GetHostListByIP, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Host/getsethostlist", Params: nil, Handler: host.GetSetHostList, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/getmodulehostlist", Params: nil, Handler: host.GetModuleHostList, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Host/getmodulehostlist", Params: nil, Handler: host.GetModuleHostList, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/getapphostlist", Params: nil, Handler: host.GetAppHostList, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Host/getapphostlist", Params: nil, Handler: host.GetAppHostList, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Set/gethostsbyproperty", Params: nil, Handler: host.GetHostsByProperty, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "set/gethostsbyproperty", Params: nil, Handler: host.GetHostsByProperty, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "Host/updateHostStatus", Params: nil, Handler: host.UpdateHostStatus, FilterHandler: nil, Version: v2.APIVersion})

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "Host/updateHostByAppId", Params: nil, Handler: host.UpdateHostByAppID, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Host/getCompanyIdByIps", Params: nil, Handler: host.GetCompanyIDByIps, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/getCompanyIdByIps", Params: nil, Handler: host.GetCompanyIDByIps, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Host/getHostListByAppidAndField", Params: nil, Handler: host.GetHostListByAppIDAndField, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/getHostListByAppidAndField", Params: nil, Handler: host.GetHostListByAppIDAndField, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/enterIp", Params: nil, Handler: host.EnterIP, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/enterip", Params: nil, Handler: host.EnterIP, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Host/getIPAndProxyByCompany", Params: nil, Handler: host.GetIPAndProxyByCompany, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "Host/updatehostmodule", Params: nil, Handler: host.UpdateHostModule, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "host/updatehostmodule", Params: nil, Handler: host.UpdateHostModule, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "host/updateCustomProperty", Params: nil, Handler: host.UpdateCustomProperty, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "host/cloneHostProperty", Params: nil, Handler: host.CloneHostProperty, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "host/delHostInApp", Params: nil, Handler: host.DelHostInApp, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/getgitServerIp", Params: nil, Handler: host.GetGitServerIp, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/addhost", Params: nil, Handler: host.AddHost, FilterHandler: nil, Version: v2.APIVersion})

	// set cc api interface
	host.CreateAction()
}

// GetHostListByIP get host list by ip , multiple ip use split comma
func (cli *hostAction) GetHostListByIP(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetHostListByIP start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("GetHostListByIP error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	if len(formData["IP"]) == 0 || formData["IP"][0] == "" {
		blog.Error("GetHostListByIP error: param IP is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "IP").Error(), resp)
		return
	}

	ipArr := strings.Split(formData["IP"][0], ",")

	//build v3 params
	param := map[string]interface{}{
		common.BKIPListField: ipArr,
	}
	if len(formData["ApplicationID"]) > 0 {
		appIDStrArr := strings.Split(formData["ApplicationID"][0], ",")
		appIDArr, err := utils.SliceStrToInt(appIDStrArr)
		if nil != err {
			blog.Error("GetHostListByIP error: %v", err)
			converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
			return
		}
		param[common.BKAppIDField] = appIDArr
	}

	if len(formData["platID"]) > 0 {
		param[common.BKCloudIDField] = formData["platID"][0]
	}

	paramJson, _ := json.Marshal(param)
	url := fmt.Sprintf("%s/host/v1/gethostlistbyip", cli.CC.HostAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	blog.Debug("rspV3:%v", rspV3)
	if err != nil {
		blog.Error("GetHostListByIP url:%s, params:%s, error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	resDataV2, err := converter.ResToV2ForHostList(rspV3)
	if err != nil {
		blog.Error("convert host res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}
	converter.RespSuccessV2(resDataV2, resp)
}

// getCompanyIDByIps  get company id by ips
func (cli *hostAction) GetCompanyIDByIps(req *restful.Request, resp *restful.Response) {
	blog.Debug("getCompanyIDByIps start!")

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("getCompanyIDByIps error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form
	blog.Debug("getCompanyIDByIps http body data: %v", formData)
	if len(formData["Ips"]) == 0 {
		blog.Error("getCompanyIDByIps error: param ips is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "Ips").Error(), resp)
		return
	}

	ipArr := strings.Split(formData["Ips"][0], ",")
	//build v3 params
	param := map[string]interface{}{
		common.BKIPListField: ipArr,
	}

	paramJson, _ := json.Marshal(param)
	url := fmt.Sprintf("%s/host/v1/gethostlistbyip", cli.CC.HostAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	if err != nil {
		blog.Error("getCompanyIDByIps url:%s, params:%s, error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForCpyHost(rspV3)
	if err != nil {
		blog.Error("convert host res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

// GetSetHostList get host list by set id
func (cli *hostAction) GetSetHostList(req *restful.Request, resp *restful.Response) {
	blog.Debug("getSetHostList start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("getSetHostList error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "SetID"})
	if !res {
		blog.Error("getSetHostList error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	appID, err := strconv.Atoi(formData["ApplicationID"][0])
	if nil != err {
		blog.Error("getSetHostList error: %v", err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	setIDStrArr := strings.Split(formData["SetID"][0], ",")
	setIDArr, err := utils.SliceStrToInt(setIDStrArr)
	if nil != err {
		blog.Error("getSetHostList error: %v", err)
		converter.RespFailV2(common.CCErrAPIServerV2MultiSetIDErr, defErr.Error(common.CCErrAPIServerV2MultiSetIDErr).Error(), resp)
		return
	}

	param := map[string]interface{}{
		common.BKAppIDField: appID,
		common.BKSetIDField: setIDArr,
	}
	paramJson, _ := json.Marshal(param)

	url := fmt.Sprintf("%s/host/v1/getsethostlist", cli.CC.HostAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	if err != nil {
		blog.Error("getSetHostList url:%s, params:%s error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForHostList(rspV3)
	if err != nil {
		blog.Error("convert host res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

// GetModuleHostList get host list by module
func (cli *hostAction) GetModuleHostList(req *restful.Request, resp *restful.Response) {
	blog.Debug("getModuleHostList start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("getModuleHostList error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "ModuleID"})
	if !res {
		blog.Error("getModuleHostList error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	appID, err := strconv.Atoi(formData["ApplicationID"][0])
	if nil != err {
		blog.Error("getModuleHostList error: %v", err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	moduleIDStrArr := strings.Split(formData["ModuleID"][0], ",")
	moduleIDArr, err := utils.SliceStrToInt(moduleIDStrArr)
	if nil != err {
		blog.Error("getModuleHostList error: %v", err)
		converter.RespFailV2(common.CCErrAPIServerV2MultiModuleIDErr, defErr.Error(common.CCErrAPIServerV2MultiModuleIDErr).Error(), resp)
		return
	}

	param := map[string]interface{}{
		common.BKAppIDField:    appID,
		common.BKModuleIDField: moduleIDArr,
	}
	paramJson, err := json.Marshal(param)
	if err != nil {
		blog.Error("getModuleHostList Marshal json error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}

	url := fmt.Sprintf("%s/host/v1/getmodulehostlist", cli.CC.HostAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	if err != nil {
		blog.Error("getModuleHostList url:%s, params:%s, error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForHostList(rspV3)
	if err != nil {
		blog.Error("convert host res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

// GetAppHostList get host list by ApplicationID
func (cli *hostAction) GetAppHostList(req *restful.Request, resp *restful.Response) {
	blog.Debug("getAppHostList start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("getAppHostList error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("getAppHostList http body data: %v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID"})
	if !res {
		blog.Error("getAppHostList error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	appID, err := strconv.Atoi(formData["ApplicationID"][0])
	if nil != err {
		blog.Error("getAppHostList error: %v", err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	param := map[string]interface{}{
		common.BKAppIDField: appID,
	}
	paramJson, _ := json.Marshal(param)

	url := fmt.Sprintf("%s/host/v1/getapphostlist", cli.CC.HostAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	if err != nil {
		blog.Error("getAppHostList url:%s, params:%s, error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForHostList(rspV3)
	if err != nil {
		blog.Error("convert host res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	blog.Debug("GetAppHostList success, data length: %d", len(resDataV2.([]map[string]interface{})))

	converter.RespSuccessV2(resDataV2, resp)
}

// GetHostsByProperty get host by set property
func (cli *hostAction) GetHostsByProperty(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetHostsByProperty start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("GetHostsByProperty error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID"})
	if !res {
		blog.Error("GetHostsByProperty error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	appID, err := strconv.Atoi(formData["ApplicationID"][0])
	if nil != err {
		blog.Error("GetHostsByProperty error: %v", err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	//build v3 params
	param := map[string]interface{}{
		common.BKAppIDField: appID,
	}

	if len(formData["SetID"]) > 0 && "" != formData["SetID"][0] {
		setIDArr, err := utils.SliceStrToInt(strings.Split(formData["SetID"][0], ","))
		if nil != err {
			blog.Error("GetHostsByProperty error: %v", err)
			converter.RespFailV2(common.CCErrAPIServerV2MultiSetIDErr, defErr.Errorf(common.CCErrAPIServerV2MultiSetIDErr).Error(), resp)
			return
		}
		param[common.BKSetIDField] = setIDArr
	}

	if len(formData["SetEnviType"]) > 0 && "" != formData["SetEnviType"][0] {

		//env 1：测试 2：体验 3：正式，默认为3
		setEnvArrTemp := strings.Split(formData["SetEnviType"][0], ",")
		setEnvArr := make([]string, 0)
		for _, setEnv := range setEnvArrTemp {
			setEnvV3, ok := defs.SetEnvMap[setEnv]
			if !ok {
				msg := fmt.Sprintf("SetEnviType not in 1,2,3, it is %s", setEnv)
				blog.Error(msg)
				converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
				return
			}
			setEnvArr = append(setEnvArr, setEnvV3)
		}
		param[common.BKSetEnvField] = setEnvArr
	}

	if len(formData["SetServiceStatus"]) > 0 && "" != formData["SetServiceStatus"][0] {

		//service status，包含0：关闭，1：开启，默认为1
		setStatusArrTemp := strings.Split(formData["SetServiceStatus"][0], ",")
		setStatusArr := make([]string, 0)
		for _, setStatus := range setStatusArrTemp {
			setStatusV3, ok := defs.SetStatusMap[setStatus]
			if !ok {
				msg := fmt.Sprintf("SetServiceStatus not in 0,1, it is %s", setStatus)
				blog.Error(msg)
				converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
				return
			}
			setStatusArr = append(setStatusArr, setStatusV3)
		}

		param[common.BKSetStatusField] = setStatusArr
	}

	paramJson, _ := json.Marshal(param)

	url := fmt.Sprintf("%s/host/v1/gethostsbyproperty", cli.CC.HostAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	if err != nil {
		blog.Error("GetHostsByProperty url:%s, params:%s error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForHostList(rspV3)
	if err != nil {
		blog.Error("convert host res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

// updateHostStatus update host proxy status , set GseProxy field equal 1
func (cli *hostAction) UpdateHostStatus(req *restful.Request, resp *restful.Response) {
	blog.Debug("updateHostStatus start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("updateHostStatus error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("updateHostStatus http body data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"appId", "platId", "ip"})
	if !res {
		blog.Error("updateHostStatus error: %s", msg)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	appID := formData["appId"][0]

	ip := formData["ip"][0]
	platID, _ := strconv.Atoi(formData["platId"][0])

	param := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKHostInnerIPField: ip,
			common.BKSubAreaField:     platID,
		},
		"data": map[string]interface{}{
			common.BKGseProxyField: "1",
			common.BKSubAreaField:  platID,
		},
	}
	paramJson, err := json.Marshal(param)
	if err != nil {
		blog.Error("updateHostStatus error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}

	url := fmt.Sprintf("%s/host/v1/openapi/host/%s", cli.CC.HostAPI(), appID)
	blog.Debug("url:%v", url)

	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, []byte(paramJson))
	blog.Debug("rspV3:%v", rspV3)
	if err != nil {
		blog.Error("updateHostStatus url:%s,  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2([]byte(rspV3), resp)
}

// updateHostStatus update proxy satus by ip
func (cli *hostAction) UpdateHostByAppID(req *restful.Request, resp *restful.Response) {
	blog.Debug("updateHostByAppID start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("updateHostByAppID error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("updateHostByAppID http body data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"appId", "platId", "proxyList"})
	if !res {
		blog.Error("updateHostByAppID error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	appID := formData["appId"][0]

	proxyList := formData["proxyList"][0]

	platID := formData.Get("platId")
	platIdInt, _ := strconv.Atoi(platID)
	proxyListArr := make([]map[string]interface{}, 0)
	json.Unmarshal([]byte(proxyList), &proxyListArr)
	proxyListArrV3 := make([]common.KvMap, 0)
	for _, proxy := range proxyListArr {
		proxyNew := make(map[string]interface{})
		proxyNew[common.BKCloudIDField] = platIdInt
		proxyNew[common.BKHostInnerIPField] = proxy["InnerIP"]
		proxyNew[common.BKHostOuterIPField] = proxy["OuterIP"]
		proxyNew, err := logics.AutoInputV3Field(proxyNew, common.BKInnerObjIDHost, host.CC.TopoAPI(), req.Request.Header)

		if err != nil {
			blog.Error("AutoInputV3Field error:%v", err)
			converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, err.Error()).Error(), resp)
			return
		}
		proxyListArrV3 = append(proxyListArrV3, proxyNew)
	}

	blog.Debug("proxyListArrV3:%v", proxyListArrV3)
	param := map[string]interface{}{

		common.BKCloudIDField:   platID,
		common.BKProxyListField: proxyListArrV3,
	}
	paramJson, _ := json.Marshal(param)

	url := fmt.Sprintf("%s/host/v1/host/updateHostByAppID/%s", cli.CC.HostAPI(), appID)
	blog.Debug(url)
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, []byte(paramJson))
	blog.Debug("rspV3:%v", rspV3)
	if err != nil {
		blog.Error("updateHostByAppID  url:%s, params:%s, error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2([]byte(rspV3), resp)
}

// getSetHostList get host list by application and fiels
func (cli *hostAction) GetHostListByAppIDAndField(req *restful.Request, resp *restful.Response) {
	blog.Debug("getHostListByAppidAndField start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("getHostListByAppidAndField error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	blog.Info("start info")
	formData := req.Request.Form

	res, msg := utils.ValidateFormData(formData, []string{"appId", "field"})
	if !res {
		blog.Error("getHostListByAppidAndField url error: %s", msg)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	appID, err := strconv.Atoi(formData["appId"][0])
	if nil != err {
		blog.Error("getHostListByAppidAndField error: %v", err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "appId").Error(), resp)
		return
	}

	field := formData["field"][0]

	url := fmt.Sprintf("%s/host/v1/host/getHostListByAppidAndField/%d/%s", cli.CC.HostAPI(), appID, converter.ConverterV2FieldsToV3(field, common.BKInnerObjIDHost))
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectGet, nil)
	if err != nil {
		blog.Error("getHostListByAppidAndField url:%s, params:%s, error:%v", url, err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	blog.Debug("rspV3: %v", rspV3)

	resDataV2, err := converter.ResToV2ForHostGroup(rspV3)
	if err != nil {
		blog.Error("convert host res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (cli *hostAction) EnterIP(req *restful.Request, resp *restful.Response) {
	blog.Debug("EnterIP start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("EnterIP error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form
	ips := formData.Get("ip")
	hostName := formData.Get("hostname")
	moduleName := formData.Get("moduleName")
	appName := formData.Get("appName")
	setName := formData.Get("setName")
	osType := formData.Get("osType")

	if "" == ips {
		blog.Errorf("EnterIP error ips empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ip").Error(), resp)
		return
	}
	ipArr := strings.Split(ips, ",")
	var hostNameArr []string
	if "" != hostName {
		hostNameArr = strings.Split(hostName, ",")
	}
	if osType == "window" {
		osType = "windows"
	}
	if "" != osType && osType != "windows" && osType != "linux" {
		blog.Errorf("osType mast be windows or linux; not %s", osType)
		converter.RespFailV2(common.CCErrAPIServerV2OSTypeErr, defErr.Error(common.CCErrAPIServerV2OSTypeErr).Error(), resp)
		return
	}
	osType = "1"
	if "windows" == osType {
		osType = "2"
	}
	input := make(common.KvMap)
	input["ips"] = ipArr
	input[common.BKHostNameField] = hostNameArr
	input[common.BKModuleNameField] = moduleName
	input[common.BKSetNameField] = setName
	input[common.BKAppNameField] = appName
	input[common.BKOSTypeField] = osType
	input[common.BKOwnerIDField] = common.BKDefaultOwnerID

	paramJson, _ := json.Marshal(input)
	url := fmt.Sprintf("%s/host/v1/host/add/module", cli.CC.HostAPI())
	blog.Infof("http request for add module url:%s, params:%s", url, string(paramJson))
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	blog.Infof("http request for add module url:%s, reply:%s", url, rspV3)
	if err != nil {
		blog.Error("EnterIP url:%s, params:%s, error:%s", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	err = converter.ResToV2ForEnterIP(rspV3)
	if err != nil {
		blog.Error("convert EnterIP result to v2 error:%s, reply:%s", err.Error(), rspV3)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}
	converter.RespSuccessV2("", resp)
}

//getIPAndProxxyByCompany   get proxy by  commpay
func (cli *hostAction) GetIPAndProxyByCompany(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetIPAndProxyByCompany start")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("GetIPAndProxyByCompany Error %v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	appID := formData.Get("appId")
	platID := formData.Get("platId")
	ips := formData.Get("ipList")
	if "" == appID {
		blog.Errorf("GetIPAndProxyByCompany error appID empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "appId").Error(), resp)
		return
	}
	if "" == platID {
		blog.Errorf("GetIPAndProxyByCompany error platID empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "platId").Error(), resp)
		return
	}
	if "" == ips {
		blog.Errorf("GetIPAndProxyByCompany error ipList empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ipList").Error(), resp)
		return
	}
	ipArr := strings.Split(ips, ",")
	input := make(common.KvMap)
	input["ips"] = ipArr
	input[common.BKAppIDField] = appID
	input[common.BKCloudIDField] = platID
	paramJson, _ := json.Marshal(input)
	blog.Debug("paramJson:%v", string(paramJson))
	url := fmt.Sprintf("%s/host/v1/getIPAndProxyByCompany", cli.CC.HostAPI())
	blog.Infof("http request for hosts search url:%s, params:%s", url, string(paramJson))
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	blog.Infof("http request for  hosts search url:%s, reply:%s", url, rspV3)
	if err != nil {
		blog.Error("GetIPAndProxyByCompany url:%s, params:%s, error:%s ", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	js, err := simplejson.NewJson([]byte(rspV3))
	if err != nil {
		blog.Error("simplejson error:%s , reply:%s", err.Error(), rspV3)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}
	resData, err := js.Map()
	converter.RespSuccessV2(resData, resp)
}

//updateHostModule update host module relation
func (cli *hostAction) UpdateHostModule(req *restful.Request, resp *restful.Response) {
	blog.Debug("UpdateHostModule start")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("UpdateHostModule Error %v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form

	blog.Infof("UpdateHostModule http body data:%v", formData)

	appID := formData.Get("ApplicationID")
	platID := formData.Get("platId")
	moduleID := formData.Get("dstModuleID")
	ips := formData.Get("ip")

	if "" == appID {
		blog.Errorf("UpdateHostModule error ApplicationID empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}
	if "" == platID {
		blog.Errorf("UpdateHostModule error platID empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "platID").Error(), resp)
		return
	}
	if "" == moduleID {
		blog.Errorf("UpdateHostModule error moduleID empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "dstModuleID").Error(), resp)
		return
	}
	if "" == ips {
		blog.Errorf("UpdateHostModule error ips empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ips").Error(), resp)
		return
	}
	ipArr := strings.Split(ips, ",")
	moduleIDArr, err := utils.SliceStrToInt(strings.Split(moduleID, ","))
	if nil != err {
		blog.Error("getHostListByAppidAndField error: %v", err)
		converter.RespFailV2(common.CCErrAPIServerV2MultiModuleIDErr, defErr.Error(common.CCErrAPIServerV2MultiModuleIDErr).Error(), resp)
		return
	}
	appIDArr := make([]int, 0)
	appIDInt, _ := strconv.Atoi(appID)
	appIDArr = append(appIDArr, appIDInt)
	platIDInt, _ := strconv.Atoi(platID)
	param := map[string]interface{}{
		common.BKIPListField:  ipArr,
		common.BKAppIDField:   appIDArr,
		common.BKSubAreaField: platIDInt,
	}
	paramJson, _ := json.Marshal(param)

	url := fmt.Sprintf("%s/host/v1/gethostlistbyip", cli.CC.HostAPI())
	hosts, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	if nil != err {
		blog.Error("UpdateHostModule url:%s, params:%s, error:%s ", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	js, err := simplejson.NewJson([]byte(hosts))
	if nil != err {
		converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		blog.Errorf("UpdateHostModule error simplejson.NewJson error, data:%s", string(hosts))
		return
	}
	hostsMap, err := js.Map()
	if nil != err {
		blog.Errorf("UpdateHostModule error js.Map error, data:%s", string(hosts))
		converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}
	HostIDArr := make([]int, 0)
	for _, item := range hostsMap["data"].([]interface{}) {
		blog.Debug("item:%v", item)
		hostID := item.(map[string]interface{})[common.BKHostIDField]
		blog.Debug("hostIDInt:%d", hostID)
		hostIDInt, _ := util.GetIntByInterface(hostID)
		HostIDArr = append(HostIDArr, int(hostIDInt))

	}
	blog.Debug("HostIDArr:%v", HostIDArr)
	urlInput := ""
	// host translate module
	moduleMap, err := GetModulesByApp(req, map[string]interface{}{
		common.BKAppIDField: appID}, cli.CC.TopoAPI())
	blog.Debug("moduleMap:%v", moduleMap)
	if nil != err {
		blog.Error("GetModulesByApp error: %v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Errorf(common.CCErrCommReplyDataFormatError, err.Error()).Error(), resp)
		return
	}

	input := make(common.KvMap)
	input[common.BKAppIDField] = appIDInt
	input[common.BKHostIDField] = HostIDArr
	if len(moduleIDArr) > 1 {
		for moduleId := range moduleIDArr {
			moduleName := moduleMap[moduleId].(map[string]interface{})[common.BKModuleNameField].(string)
			if moduleName == common.DefaultFaultModuleName || moduleName == common.DefaultResModuleName {
				msg := defErr.Error(common.CCErrAPIServerV2HostModuleContainDefaultModuleErr).Error()
				blog.Error("update host module error: %v", msg)
				converter.RespFailV2(common.CCErrAPIServerV2HostModuleContainDefaultModuleErr, msg, resp)
				return
			}
		}
		urlInput = fmt.Sprintf("%s/host/v1/hosts/modules", cli.CC.HostAPI())
		input[common.BKModuleIDField] = moduleIDArr
		input[common.BKIsIncrementField] = false
	} else {
		if moduleMap[moduleIDArr[0]].(map[string]interface{})[common.BKModuleNameField].(string) == common.DefaultFaultModuleName {
			urlInput = fmt.Sprintf("%s/host/v1/hosts/faultmodule", cli.CC.HostAPI())
		}
		if moduleMap[moduleIDArr[0]].(map[string]interface{})[common.BKModuleNameField].(string) == common.DefaultResModuleName {
			urlInput = fmt.Sprintf("%s/host/v1/hosts/emptymodule", cli.CC.HostAPI())
		} else {
			urlInput = fmt.Sprintf("%s/host/v1/hosts/modules", cli.CC.HostAPI())
			input[common.BKModuleIDField] = moduleIDArr
			input[common.BKIsIncrementField] = false
		}
	}

	inputJson, _ := json.Marshal(input)
	blog.Infof("update host module url:%v", url)
	blog.Infof("http request for hosts search url:%s, params:%s", urlInput, string(inputJson))
	rspV3, err := httpcli.ReqHttp(req, urlInput, common.HTTPCreate, []byte(inputJson))
	blog.Debug("http request for  hosts search url:%s, reply:%s", urlInput, rspV3)
	if err != nil {
		blog.Error("UpdateHostModule url:%s, params:%s, error:%s , reply:%s", url, string(inputJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	converter.RespCommonResV2([]byte(rspV3), resp)
}

// UpdateCustomProperty update host property
func (cli *hostAction) UpdateCustomProperty(req *restful.Request, resp *restful.Response) {
	blog.Debug("UpdateCustomProperty start")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if nil != err {
		blog.Error("UpdateHostModule Error %v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Error("formData:%v", formData)
	appId := formData.Get("ApplicationID")
	hostId := formData.Get("HostID")
	property := formData.Get("Property")
	if "" == appId {
		blog.Errorf("UpdateCustomProperty error platId empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}
	if "" == hostId {
		blog.Errorf("UpdateCustomProperty error host empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "HostID").Error(), resp)
		return
	}
	param := make(common.KvMap)
	param[common.BKAppIDField] = appId
	param[common.BKHostIDField] = hostId
	param["property"] = property
	paramJson, _ := json.Marshal(param)
	url := fmt.Sprintf("%s/host/v1/openapi/updatecustomproperty", cli.CC.HostAPI())
	blog.Infof("http request for UpdateCustomProperty url:%s, params:%s", url, string(paramJson))
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, []byte(paramJson))
	blog.Infof("http request for UpdateCustomProperty url:%s, reply:%s", url, rspV3)
	if err != nil {
		blog.Error("UpdateCustomProperty url:%s, params:%s, error:%s , reply:%s", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	rspV3Map := make(map[string]interface{})
	err = json.Unmarshal([]byte(rspV3), &rspV3Map)
	if nil != err {
		converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}
	converter.RespCommonResV2([]byte(rspV3), resp)
}

// CloneHostProperty clone host
func (cli *hostAction) CloneHostProperty(req *restful.Request, resp *restful.Response) {
	blog.Debug("CloneHostProperty start")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if nil != err {
		blog.Error("CloneHostProperty Error %v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Error("formData:%v", formData)
	res, msg := utils.ValidateFormData(formData, []string{
		"ApplicationID",
		"orgIp",
		"dstIp",
	})
	if !res {
		blog.Error("CloneHostProperty error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	appId := formData.Get("ApplicationID")
	orgIp := formData.Get("orgIp")
	dstIp := formData.Get("dstIp")
	platId := formData.Get("platId")

	param := make(common.KvMap)
	param[common.BKAppIDField] = appId
	param[common.BKOrgIPField] = orgIp
	param[common.BKDstIPField] = dstIp
	param[common.BKCloudIDField] = platId
	paramJson, _ := json.Marshal(param)

	url := fmt.Sprintf("%s/host/v1/openapi/host/clonehostproperty", cli.CC.HostAPI())
	blog.Infof("http request for CloneHostProperty url:%s, params:%s", url, string(paramJson))
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, []byte(paramJson))
	blog.Infof("http request for CloneHostProperty url:%s, reply:%s", url, rspV3)
	if err != nil {
		blog.Error("CloneHostProperty url:%s, params:%s, error:%s", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	rspV3Map := make(map[string]interface{})
	err = json.Unmarshal([]byte(rspV3), &rspV3Map)
	if nil != err {
		blog.Error("CloneHostProperty Unmarshal json error:%v, rspV3:%s", err, rspV3)
		return
	}
	converter.RespCommonResV2([]byte(rspV3), resp)
}

// DelHostInApp del host from application idle set
func (cli *hostAction) DelHostInApp(req *restful.Request, resp *restful.Response) {
	blog.Debug("DelHostInApp start")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if nil != err {
		blog.Error("DelHostInApp Error %v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Debug("formData:%v", formData)
	res, msg := utils.ValidateFormData(formData, []string{
		"ApplicationID",
		"HostID",
	})
	if !res {
		blog.Error("DelHostInApp error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	appId := formData.Get("ApplicationID")
	hostId, _ := util.GetInt64ByInterface(formData.Get("HostID"))
	param := make(common.KvMap)
	param[common.BKAppIDField], _ = util.GetInt64ByInterface(appId)
	param[common.BKHostIDField] = []int64{hostId}
	paramJson, _ := json.Marshal(param)
	//url := fmt.Sprintf("%s/host/v1/openapi/host/delhostinapp", cli.CC.HostAPI())
	url := fmt.Sprintf("%s/host/v1/hosts/resource", cli.CC.HostAPI())
	blog.Infof("http request for DelHostInApp url:%s, params:%s", url, string(paramJson))
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPCreate, []byte(paramJson))
	blog.Infof("http request for DelHostInApp url:%s, reply:%s", url, rspV3)
	if err != nil {
		blog.Error("DelHostInApp url:%s, params:%s, error:%s", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	rspV3Map := make(map[string]interface{})
	err = json.Unmarshal([]byte(rspV3), &rspV3Map)
	if nil != err {
		blog.Error("DelHostInApp Unmarshal json error:%v, rspV3:%s", err, rspV3)
		return
	}
	converter.RespCommonResV2([]byte(rspV3), resp)
}

// GetGitServerIp get white list ip
func (cli *hostAction) GetGitServerIp(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetGitServerIp start")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	appName := common.WhiteListAppName
	setName := common.WhiteListSetName
	moduleName := common.WhiteListModuleName
	param := make(common.KvMap)
	param[common.BKAppNameField] = appName
	param[common.BKSetNameField] = setName
	param[common.BKModuleNameField] = moduleName
	paramJson, _ := json.Marshal(param)

	url := fmt.Sprintf("%s/host/v1/openapi/host/getGitServerIp", cli.CC.HostAPI())
	blog.Infof("http request for GetGitServerIp url:%s, params:%s", url, string(paramJson))
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	blog.Infof("http request for GetGitServerIp url:%s, reply:%s", url, rspV3)
	if err != nil {
		blog.Error("GetGitServerIp url:%s, params:%s, error:%s ", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	blog.Debug("rspV3:%v", rspV3)

	resDataV2, err := converter.ResToV2ForHostList(rspV3)
	if err != nil {
		blog.Error("convert host res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}
	converter.RespSuccessV2(resDataV2, resp)
}

// GetModulesByApp get module by appliaciton
func GetModulesByApp(req *restful.Request, condition map[string]interface{}, topoApi string) (map[int]interface{}, error) {
	moduleMap := make(map[int]interface{})
	appId := condition[common.BKAppIDField]
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

	url := fmt.Sprintf("%s/topo/v1/openapi/module/searchByApp/%v", topoApi, appId)
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))

	if err != nil {
		blog.Error("convert module res to v2 url:%s, params:%s, error:%v", url, string(paramJson), err)
		return nil, err
	}
	blog.Info("GetModuleMapByCond return :%s", string(reply))
	if err != nil {
		blog.Errorf("GetModuleMapByCond params:%s  error:%v", string(paramJson), err)
		return moduleMap, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()

	result, err := js.Get("result").Bool()
	if false == result {
		msg, _ := js.Get(common.HTTPBKAPIErrorMessage).String()
		return nil, errors.New(msg)
	}
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
