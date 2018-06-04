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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"

	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/defs"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
)

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
		appIDArr, sliceErr := utils.SliceStrToInt(appIDStrArr)
		if nil != sliceErr {
			blog.Error("GetHostListByIP error: %v", sliceErr)
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
		setIDArr, sliceErr := utils.SliceStrToInt(strings.Split(formData["SetID"][0], ","))
		if nil != sliceErr {
			blog.Error("GetHostsByProperty error: %v", sliceErr)
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
	resData := js.Get("data").Interface()
	converter.RespSuccessV2(resData, resp)
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
