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

package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"

	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (s *Service) updateHostStatus(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("update host status error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("update host status data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"appId", "platId", "ip"})
	if !res {
		blog.Error("update host status error: %s", msg)
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

	if err != nil {
		blog.Error("update host status error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}
	result, err := s.CoreAPI.HostServer().UpdateHost(context.Background(), appID, pheader, param)
	if err != nil {
		blog.Error("update host status  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)

}

func (s *Service) updateHostByAppID(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	user := util.GetUser(pheader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("update host by app id error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("update host by app id data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"appId", "platId", "proxyList"})
	if !res {
		blog.Error("update host by app id error: %s", msg)
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
		proxyNew, inputErr := s.Logics.AutoInputV3Field(proxyNew, common.BKInnerObjIDHost, user, pheader)

		if inputErr != nil {
			blog.Error("AutoInputV3Field error:%v", inputErr)
			converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, inputErr.Error()).Error(), resp)
			return
		}
		proxyListArrV3 = append(proxyListArrV3, proxyNew)
	}

	blog.Debug("proxyListArrV3:%v", proxyListArrV3)
	param := map[string]interface{}{

		common.BKCloudIDField:   platID,
		common.BKProxyListField: proxyListArrV3,
	}

	result, err := s.CoreAPI.HostServer().UpdateHostByAppID(context.Background(), appID, pheader, param)
	if err != nil {
		blog.Error("update host by app id   error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)

}

func (s *Service) getCompanyIDByIps(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("get company id by ips error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form
	blog.Debug("get company id by ips data: %v", formData)

	if len(formData["Ips"]) == 0 {
		blog.Error("get company id by ips error: param ips is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "Ips").Error(), resp)
		return
	}

	ipArr := strings.Split(formData["Ips"][0], ",")
	//build v3 params
	param := map[string]interface{}{
		common.BKIPListField: ipArr,
	}

	result, err := s.CoreAPI.HostServer().HostSearchByIP(context.Background(), pheader, param)
	if err != nil {
		blog.Error("get company id by ips  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForCpyHost(result.Result, result.ErrMsg, result.Data)
	if err != nil {
		blog.Error("convert host res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) getHostListByAppIDAndField(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("get host list by appid and field error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	res, msg := utils.ValidateFormData(formData, []string{"appId", "field"})
	if !res {
		blog.Error("get host list by appid and field error: %s", msg)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	appID := formData["appId"][0]
	if nil != err {
		blog.Error("get host list by appid and field error: %v", err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "appId").Error(), resp)
		return
	}

	field := formData["field"][0]
	result, err := s.CoreAPI.HostServer().GetHostListByAppidAndField(context.Background(), appID, converter.ConverterV2FieldsToV3(field, common.BKInnerObjIDHost), pheader)
	if err != nil {
		blog.Error("get host list by appid and field  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForHostGroup(result.Result, result.ErrMsg, result.Data)
	if err != nil {
		blog.Error("convert host res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) updateHostModule(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	const (
		NormalTrans = "nomal"
		EmptyTrans  = "empty"
		FaultTrans  = "fault"
	)
	var hostTransType string
	var hostModuleParam metadata.DefaultModuleHostConfigParams
	var result *metadata.Response
	var err error

	err = req.Request.ParseForm()
	if err != nil {
		blog.Error("update host module error %v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form

	blog.Infof("update host module data:%v", formData)

	appID := formData.Get("ApplicationID")
	platID := formData.Get("platId")
	moduleID := formData.Get("dstModuleID")
	ips := formData.Get("ip")

	if "" == appID {
		blog.Errorf("update host module error ApplicationID empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}
	if "" == platID {
		blog.Errorf("update host module error platID empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "platID").Error(), resp)
		return
	}
	if "" == moduleID {
		blog.Errorf("update host module error moduleID empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "dstModuleID").Error(), resp)
		return
	}
	if "" == ips {
		blog.Errorf("update host module error ips empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ips").Error(), resp)
		return
	}
	ipArr := strings.Split(ips, ",")
	moduleIDArr, err := utils.SliceStrToInt(strings.Split(moduleID, ","))
	if nil != err {
		blog.Error("update host module error: %v", err)
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

	result, err = s.CoreAPI.HostServer().HostSearchByIP(context.Background(), pheader, param)

	hostsMap, ok := result.Data.([]interface{})
	if false == ok {
		blog.Errorf("update host module error js.Map error, data:%v", result.Data)
		converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}
	HostIDArr := make([]int, 0)
	HostID64Arr := make([]int64, 0)
	for _, item := range hostsMap {

		hostID := item.(map[string]interface{})[common.BKHostIDField]
		blog.Infof("hostIDInt:%d ", hostID)
		hostIDInt, _ := util.GetIntByInterface(hostID)
		HostIDArr = append(HostIDArr, int(hostIDInt))
		HostID64Arr = append(HostID64Arr, int64(hostIDInt))

	}

	blog.Infof("HostIDArr:%v", HostIDArr)

	// host translate module
	inputparam := map[string]interface{}{
		common.BKAppIDField: appID}
	moduleMap, err := s.getModuleInfoByApp(inputparam, pheader)

	if nil != err {
		blog.Error("GetModulesByApp error: %v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Errorf(common.CCErrCommReplyDataFormatError, err.Error()).Error(), resp)
		return
	}

	input := make(common.KvMap)
	input[common.BKAppIDField] = appIDInt
	input[common.BKHostIDField] = HostIDArr
	hostModuleParam.ApplicationID = int64(appIDInt)
	hostModuleParam.HostID = HostID64Arr
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
		hostTransType = NormalTrans
		input[common.BKModuleIDField] = moduleIDArr
		input[common.BKIsIncrementField] = false
	} else {
		if moduleMap[moduleIDArr[0]].(map[string]interface{})[common.BKModuleNameField].(string) == common.DefaultFaultModuleName {
			hostTransType = FaultTrans
		}
		if moduleMap[moduleIDArr[0]].(map[string]interface{})[common.BKModuleNameField].(string) == common.DefaultResModuleName {
			hostTransType = EmptyTrans
		} else {
			hostTransType = NormalTrans
			input[common.BKModuleIDField] = moduleIDArr
			input[common.BKIsIncrementField] = false
		}
	}
	switch hostTransType {
	case NormalTrans:
		result, err = s.CoreAPI.HostServer().HostModuleRelation(context.Background(), pheader, input)
	case EmptyTrans:
		result, err = s.CoreAPI.HostServer().MoveHost2EmptyModule(context.Background(), pheader, &hostModuleParam)
	case FaultTrans:
		result, err = s.CoreAPI.HostServer().MoveHost2FaultModule(context.Background(), pheader, &hostModuleParam)
	}

	if err != nil {
		blog.Error("update host module  error:%v ", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)
}

func (s *Service) updateCustomProperty(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if nil != err {
		blog.Error("update custom property Error %v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Error("formData:%v", formData)
	appId := formData.Get("ApplicationID")
	hostId := formData.Get("HostID")
	property := formData.Get("Property")
	if "" == appId {
		blog.Errorf("update custom property error platId empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}
	if "" == hostId {
		blog.Errorf("update custom property error host empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "HostID").Error(), resp)
		return
	}
	param := make(common.KvMap)
	param[common.BKAppIDField] = appId
	param[common.BKHostIDField] = hostId
	param["property"] = property

	result, err := s.CoreAPI.HostServer().UpdateCustomProperty(context.Background(), pheader, param)
	if err != nil {
		blog.Error("update custom property error:%v ", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	if false == result.Result {
		converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}
	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)
}

func (s *Service) cloneHostProperty(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()

	if nil != err {
		blog.Error("clone host property Error %v", err)
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
		blog.Error("clone host property  error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	appId := formData.Get("ApplicationID")
	orgIp := formData.Get("orgIp")
	dstIp := formData.Get("dstIp")
	platId := formData.Get("platId")

	appIdI, err := strconv.ParseInt(appId, 10, 64)
	if nil != err {
		blog.Error("clone host property  error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	platIdI, err := strconv.ParseInt(platId, 10, 64)

	if nil != err {
		blog.Error("clone host property  error: %s", msg)
		platIdI = 0
	}

	var param metadata.HostCloneInputParams
	param.AppID = appIdI
	param.DstIP = dstIp
	param.OrgIP = orgIp
	param.PlatID = platIdI

	result, err := s.CoreAPI.HostServer().CloneHostProperty(context.Background(), pheader, &param)
	if err != nil {
		blog.Error("clone host property error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)
}

func (s *Service) delHostInApp(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if nil != err {
		blog.Error("delete host in app Error %v", err)
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

	result, err := s.CoreAPI.HostServer().DelHostInApp(context.Background(), pheader, param)
	if err != nil {
		blog.Error("delete host in app error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)
}

func (s *Service) getGitServerIp(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	appName := common.WhiteListAppName
	setName := common.WhiteListSetName
	moduleName := common.WhiteListModuleName
	param := make(common.KvMap)
	param[common.BKAppNameField] = appName
	param[common.BKSetNameField] = setName
	param[common.BKModuleNameField] = moduleName

	result, err := s.CoreAPI.HostServer().GetGitServerIp(context.Background(), pheader, param)
	if err != nil {
		blog.Error("get git server ip, error:%v ", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForHostList(result.Result, result.ErrMsg, result.Data)
	if err != nil {
		blog.Error("convert host res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}
	converter.RespSuccessV2(resDataV2, resp)
}
