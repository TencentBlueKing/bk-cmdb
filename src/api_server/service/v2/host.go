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

package v2

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"configcenter/src/api_server/logics/v2/common/converter"
	"configcenter/src/api_server/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func (s *Service) updateHostStatus(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("updateHostStatus error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("updateHostStatus data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"appId", "platId", "ip"})
	if !res {
		blog.Errorf("updateHostStatus error: %s", msg)
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
		blog.Errorf("updateHostStatus error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}
	result, err := s.CoreAPI.HostServer().UpdateHost(context.Background(), appID, pheader, param)
	if err != nil {
		blog.Errorf("updateHostStatus  error:%v", err)
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
		blog.Errorf("updateHostByAppID error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("updateHostByAppID data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"appId", "platId", "proxyList"})
	if !res {
		blog.Errorf("updateHostByAppID error: %s", msg)
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
			blog.Errorf("AutoInputV3Field error:%v", inputErr)
			converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, inputErr.Error()).Error(), resp)
			return
		}
		proxyListArrV3 = append(proxyListArrV3, proxyNew)
	}

	blog.Infof("proxyListArrV3:%v", proxyListArrV3)
	param := map[string]interface{}{

		common.BKCloudIDField:   platID,
		common.BKProxyListField: proxyListArrV3,
	}

	result, err := s.CoreAPI.HostServer().UpdateHostByAppID(context.Background(), appID, pheader, param)
	if err != nil {
		blog.Errorf("updateHostByAppID   error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)

}

func (s *Service) getCompanyIDByIps(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	rid := util.GetHTTPCCRequestID(pheader)

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getCompanyIDByIps error:%v,rid:%s", err, rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form
	blog.V(5).Infof("getCompanyIDByIps data: %+v,rid:%s", formData, rid)

	if len(formData["Ips"]) == 0 {
		blog.Errorf("getCompanyIDByIps error: param ips is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "Ips").Error(), resp)
		return
	}

	ipArr := strings.Split(formData["Ips"][0], ",")
	//build v3 params
	param := &metadata.HostSearchByIPParams{
		IpList: ipArr,
	}
	// param := map[string]interface{}{
	// 	common.BKIPListField: ipArr,
	// }

	result, err := s.CoreAPI.HostServer().HostSearchByIP(context.Background(), pheader, param)
	if err != nil {
		blog.Errorf("getCompanyIDByIps  error:%v, input:%+v,rid:%s", err, formData, rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("getCompanyIDByIps  error:%v, input:%+v,rid:%s", err, formData, rid)
		converter.RespFailV2(result.Code, result.ErrMsg, resp)
		return
	}

	resDataV2, err := converter.ResToV2ForCpyHost(result.Result, result.ErrMsg, result.Data)
	if err != nil {
		blog.Errorf("convert host res to v2 error:%v, input:%+v,rid:%s", err, formData, rid)
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
		blog.Errorf("getHostListByAppIDAndField error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	res, msg := utils.ValidateFormData(formData, []string{"appId", "field"})
	if !res {
		blog.Errorf("getHostListByAppIDAndField error: %s", msg)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	appID := formData["appId"][0]
	if nil != err {
		blog.Errorf("getHostListByAppIDAndField error: %v", err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "appId").Error(), resp)
		return
	}

	field := formData["field"][0]
	result, err := s.CoreAPI.HostServer().GetHostListByAppidAndField(context.Background(), appID, converter.ConverterV2FieldsToV3(field, common.BKInnerObjIDHost), pheader)
	if err != nil {
		blog.Errorf("getHostListByAppIDAndField  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForHostGroup(result.Result, result.ErrMsg, result.Data)
	if err != nil {
		blog.Errorf("convert host res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) updateHostModule(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	rid := util.GetHTTPCCRequestID(pheader)

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
		blog.Errorf("updateHostModule error %v,rid:%s", err, rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form

	blog.V(5).Infof("updateHostModule data:%v,rid:%s", formData, rid)

	appID := formData.Get("ApplicationID")
	platID := formData.Get("platId")
	moduleID := formData.Get("dstModuleID")
	ips := formData.Get("ip")

	if "" == appID {
		blog.Errorf("updateHostModule error ApplicationID empty, input:%+v,rid:%s", formData, rid)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}
	if "" == platID {
		blog.Errorf("updateHostModule error platID empty, input:%+v,rid:%s", formData, rid)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "platID").Error(), resp)
		return
	}
	if "" == moduleID {
		blog.Errorf("updateHostModule error moduleID empty, input:%+v,rid:%s", formData, rid)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "dstModuleID").Error(), resp)
		return
	}
	if "" == ips {
		blog.Errorf("updateHostModule error ips empty, input:%+v,rid:%s", formData, rid)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ips").Error(), resp)
		return
	}
	ipArr := strings.Split(ips, ",")
	moduleIDArr, err := utils.SliceStrToInt(strings.Split(moduleID, ","))
	if nil != err {
		blog.Errorf("updateHostModule error: %v, input:%+v,rid:%s", err, formData, rid)
		converter.RespFailV2(common.CCErrAPIServerV2MultiModuleIDErr, defErr.Error(common.CCErrAPIServerV2MultiModuleIDErr).Error(), resp)
		return
	}
	appIDInt, err := util.GetInt64ByInterface(appID)
	if err != nil {
		blog.Errorf("updateHostModule error ApplicationID (%s) not integer. input:%+v,rid:%s", appID, formData, rid)
		converter.RespFailV2Error(defErr.Errorf(common.CCErrCommParamsNeedSet, "platID"), resp)
		return
	}
	platIDInt, err := util.GetInt64ByInterface(platID)
	if err != nil {
		blog.Errorf("updateHostModule error platID(%s)  not integer. input:%+v,rid:%s", platID, formData, rid)
		converter.RespFailV2Error(defErr.Errorf(common.CCErrCommParamsNeedSet, "platID"), resp)
		return
	}
	appIDArr := make([]int64, 0)
	appIDArr = append(appIDArr, appIDInt)
	param := &metadata.HostSearchByIPParams{
		IpList:  ipArr,
		AppID:   appIDArr,
		CloudID: &platIDInt,
	}
	// param := map[string]interface{}{
	// 	common.BKIPListField:  ipArr,
	// 	common.BKAppIDField:   appIDArr,
	// 	common.BKSubAreaField: platIDInt,
	// }

	result, err = s.CoreAPI.HostServer().HostSearchByIP(context.Background(), pheader, param)

	hostsMap, ok := result.Data.([]interface{})
	if false == ok {
		blog.Errorf("updateHostModule error js.Map error, data:%+v, input:%+v,rid:%s", result.Data, formData, rid)
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

	blog.V(5).Infof("HostIDArr:%+v,rid:%s", HostIDArr, rid)

	// host translate module
	moduleMap, err := s.getModuleInfoByApp(appIDInt, pheader)

	if nil != err {
		blog.Errorf("updateHostModule error: %v. appID:%v,input:%+v,rid:%s", err, appIDInt, formData, rid)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Errorf(common.CCErrCommReplyDataFormatError, err.Error()).Error(), resp)
		return
	}

	input := make(common.KvMap)
	input[common.BKAppIDField] = appIDInt
	input[common.BKHostIDField] = HostIDArr
	hostModuleParam.ApplicationID = int64(appIDInt)
	hostModuleParam.HostID = HostID64Arr
	if len(moduleIDArr) > 1 {
		for _, moduleID := range moduleIDArr {
			moduleInfo, ok := moduleMap[moduleID]
			if !ok {
				continue
			}
			moduleName, ok := moduleInfo[common.BKModuleNameField].(string)
			if !ok {
				continue
			}
			if moduleName == common.DefaultFaultModuleName || moduleName == common.DefaultResModuleName {
				msg := defErr.Error(common.CCErrAPIServerV2HostModuleContainDefaultModuleErr).Error()
				blog.Errorf("updateHostModule error: %v", msg)
				converter.RespFailV2(common.CCErrAPIServerV2HostModuleContainDefaultModuleErr, msg, resp)
				return
			}
		}
		hostTransType = NormalTrans
		input[common.BKModuleIDField] = moduleIDArr
		input[common.BKIsIncrementField] = false
	} else {
		moduleName, err := moduleMap[moduleIDArr[0]].String(common.BKModuleNameField)
		if err != nil {
			blog.Errorf("convert res to v2  key:%s, error:%v, moduleInfo:%+v,input:%+v,rid:%s", common.BKModuleNameField, err.Error(), moduleMap, formData, rid)
			converter.RespFailV2Error(defErr.Errorf(common.CCErrCommInstFieldConvFail, "module", "ModuleName", "int", err.Error()), resp)
			return
		}
		if moduleName == common.DefaultFaultModuleName {
			hostTransType = FaultTrans
		}
		if moduleName == common.DefaultResModuleName {
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
		blog.Errorf("updateHostModule  error:%v ", err)
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
		blog.Errorf("updateCustomProperty Error %v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form

	blog.Errorf("formData:%v", formData)

	appId := formData.Get("ApplicationID")
	hostId := formData.Get("HostID")
	property := formData.Get("Property")
	if "" == appId {
		blog.Error("updateCustomProperty error platId empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}
	if "" == hostId {
		blog.Error("updateCustomProperty error host empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "HostID").Error(), resp)
		return
	}
	param := make(common.KvMap)
	param[common.BKAppIDField] = appId
	param[common.BKHostIDField] = hostId
	param["property"] = property

	result, err := s.CoreAPI.HostServer().UpdateCustomProperty(context.Background(), pheader, param)
	if err != nil {
		blog.Errorf("updateCustomProperty error:%v ", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	if false == result.Result {
		converter.RespFailV2(result.Code, result.ErrMsg, resp)
		return
	}
	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)
}

func (s *Service) cloneHostProperty(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()

	if nil != err {
		blog.Errorf("cloneHostProperty Error %v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form

	blog.Infof("formData: %v", formData)
	res, msg := utils.ValidateFormData(formData, []string{
		"ApplicationID",
		"orgIp",
		"dstIp",
	})
	if !res {
		blog.Errorf("cloneHostProperty error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	appIDStr := formData.Get("ApplicationID")
	orgIP := formData.Get("orgIp")
	dstIP := formData.Get("dstIp")
	platIDStr := formData.Get("platId")

	appID, err := util.GetInt64ByInterface(appIDStr)
	if nil != err {
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, "ApplicationID not integer").Error(), resp)
		return
	}

	platID, err := util.GetInt64ByInterface(platIDStr)
	if nil != err {
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, "platId not integer").Error(), resp)
		return
	}

	var param metadata.CloneHostPropertyParams
	param.AppID = appID
	param.DstIP = dstIP
	param.OrgIP = orgIP
	param.CloudID = platID

	result, err := s.CoreAPI.HostServer().CloneHostProperty(context.Background(), pheader, &param)
	if err != nil {
		blog.Errorf("cloneHostProperty error:%v", err)
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
		blog.Errorf("delHostInApp Error %v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form

	blog.Infof("formData:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{
		"ApplicationID",
		"HostID",
	})
	if !res {
		blog.Errorf("delHostInApp error: %s", msg)
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
		blog.Errorf("delHostInApp error:%v", err)
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
		blog.Errorf("getGitServerIp, error:%v ", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForHostList(result.Result, result.ErrMsg, result.Data)
	if err != nil {
		blog.Errorf("convert host res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}
	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) GetHostHardInfo(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("GetHostExtInfo error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Infof("GetHostExtInfo data: %v", formData)
	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID"})
	if !res {
		blog.Errorf("GetHostExtInfo error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	appID, err := util.GetInt64ByInterface(formData["ApplicationID"][0])
	if nil != err {
		blog.Errorf("GetHostExtInfo error: %v", err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}
	param := map[string]interface{}{
		common.BKAppIDField: appID,
	}
	result, err := s.CoreAPI.HostServer().HostSearchByAppID(context.Background(), pheader, param)
	if err != nil {
		blog.Errorf("GetHostExtInfo  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("GetHostExtInfo  error, error code:%s, error message:%s", result.Code, result.ErrMsg)
		converter.RespFailV2(result.Code, result.ErrMsg, resp)
		return
	}
	dataArr := result.Data.([]interface{})
	var dataMapArr []mapstr.MapStr
	for _, item := range dataArr {
		mapItem, err := mapstr.NewFromInterface(item)
		if nil != err {
			blog.Errorf("GetHostExtInfo  error, error:%s, host info:%#v, request parammetes:%#v, request-id:%s", err.Error(), item, formData, util.GetHTTPCCRequestID(pheader))
			converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, "host info not map[string]interface").Error(), resp)
			return
		}
		dataMapArr = append(dataMapArr, mapItem)
	}
	data := converter.GetHostHardInfo(appID, dataMapArr)
	converter.RespSuccessV2(data, resp)
}
