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
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/defs"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

const (
	AppStatusTest   = "测试中"
	AppStatusOnline = "已上线"
	AppStatusStop   = "停运"
)

const (
	AppStatusTestI   = "1"
	AppStatusOnlineI = "2"
	AppStatusStopI   = "3"
)

func (s *Service) getAppList(req *restful.Request, resp *restful.Response) {

	param := &params.SearchParams{Condition: nil}

	pheader := req.Request.Header
	user := util.GetUser(pheader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	result, err := s.CoreAPI.TopoServer().Instance().SearchApp(context.Background(), user, pheader, param)

	if err != nil {
		blog.Errorf("getAppList failed, err: %v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForAppList(result.Data)

	if err != nil {
		blog.Errorf("convert v3 applist to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}
	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) getAppByID(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	user := util.GetUser(pheader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getAppByID error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("getAppByID data:%v", formData)

	if len(formData["ApplicationID"]) == 0 {
		blog.Errorf("getAppByID error: ApplicationID is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}

	appIDsArr := strings.Split(formData["ApplicationID"][0], ",")

	appIDArr := make([]int64, 0)
	for _, appID := range appIDsArr {
		id, err := strconv.ParseInt(appID, 10, 64)
		if err != nil {
			blog.Error("convert str to int error")
			converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
			return
		}
		appIDArr = append(appIDArr, id)
	}
	param := &params.SearchParams{Condition: mapstr.MapStr{common.BKAppIDField: mapstr.MapStr{common.BKDBIN: appIDArr}}}

	result, err := s.CoreAPI.TopoServer().Instance().SearchApp(context.Background(), user, pheader, param)

	if err != nil {
		blog.Errorf("getAppByID  params:%v, error:%v", param, err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForAppList(result.Data)
	if err != nil {
		blog.Errorf("convert v3 app res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) getAppByUin(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	user := util.GetUser(pheader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getAppByUin error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("getAppByUin data:%v", formData)

	if len(formData["userName"]) == 0 || formData["userName"][0] == "" {
		blog.Error("get app by uin error: userName is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedSet, "userName").Error(), resp)
		return
	}

	userName := formData["userName"][0]
	v3Username := strings.Trim(converter.DecorateUserName(userName), ",")
	orCondition := []mapstr.MapStr{
		mapstr.MapStr{
			common.BKMaintainersField: mapstr.MapStr{
				common.BKDBLIKE: fmt.Sprintf("^%s,|,%s,|,%s$|^%s$", v3Username, v3Username, v3Username, v3Username),
			},
		},
	}

	//is application maintainer
	filterOnly := false
	if len(formData["filterOnly"]) > 0 && formData["filterOnly"][0] != "" {
		if "true" == formData["filterOnly"][0] {
			filterOnly = true
		}
	}

	if filterOnly == false {
		orCondition = append(orCondition, mapstr.MapStr{
			common.BKProductPMField: mapstr.MapStr{
				common.BKDBLIKE: fmt.Sprintf("^%s,|,%s,|,%s$|^%s$", v3Username, v3Username, v3Username, v3Username),
			},
		})
	}

	param := &params.SearchParams{Condition: mapstr.MapStr{common.BKDBOR: orCondition}}

	result, err := s.CoreAPI.TopoServer().Instance().SearchApp(context.Background(), user, pheader, param)

	if err != nil {
		blog.Errorf("getAppByUin  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForAppList(result.Data)
	if err != nil {
		blog.Errorf("convert res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) getUserRoleApp(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	user := util.GetUser(pheader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("GetUserRoleApp error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	res, msg := utils.ValidateFormData(formData, []string{"uin", "roleList"})

	if !res {
		blog.Errorf("GetUserRoleApp error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	userName := formData["uin"][0]
	roleArr := strings.Split(formData["roleList"][0], ",")

	v3Username := converter.DecorateUserName(userName)

	roleOrCondition := make([]mapstr.MapStr, 0)
	for _, roleStr := range roleArr {
		roleStrV3, ok := defs.RoleMap[roleStr]
		if !ok {
			continue
		}

		roleOrCondition = append(roleOrCondition, mapstr.MapStr{
			roleStrV3: mapstr.MapStr{
				common.BKDBLIKE: fmt.Sprintf("^%s,|,%s,|,%s$|^%s$", v3Username, v3Username, v3Username, v3Username),
			},
		})
	}

	param := &params.SearchParams{Condition: mapstr.MapStr{common.BKDBOR: roleOrCondition}}

	result, err := s.CoreAPI.TopoServer().Instance().SearchApp(context.Background(), user, pheader, param)
	if err != nil {
		blog.Errorf("GetUserRoleApp  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForRoleApp(result.Data, converter.DecorateUserName(userName), roleArr)
	if err != nil {
		blog.Error("convert res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) getAppSetModuleTreeByAppId(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	user := util.GetUser(pheader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getAppSetModuleTreeByAppId error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	if len(formData["ApplicationID"]) == 0 {
		blog.Errorf("getAppSetModuleTreeByAppId error error: ApplicationID is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	intAppID, err := util.GetInt64ByInterface(formData["ApplicationID"][0])

	conds := mapstr.MapStr{}
	if 0 < len(formData["ModuleType"]) {
		conds[common.BKModuleTypeField] = formData["ModuleType"][0]
	}

	if nil != err {
		blog.Errorf("getAppSetModuleTreeByAppId   appID:%v, error:%v", formData["ApplicationID"][0], err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	topo, errCode := s.Logics.GetAppTopo(user, pheader, intAppID, conds)
	if 0 != errCode {
		converter.RespFailV2(errCode, defErr.Error(errCode).Error(), resp)
		return
	}

	if nil != topo {
		s.Logics.SetModuleHostCount([]map[string]interface{}{topo}, user, pheader)
	}
	converter.RespSuccessV2(topo, resp)
}

func (s *Service) addApp(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	user := util.GetUser(pheader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("AddApp error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form

	blog.Infof("AddApp formData:%v", formData)
	res, msg := utils.ValidateFormData(formData, []string{
		"ApplicationName",
		"Maintainers",
		"Creator",
		"LifeCycle",
	})
	if !res {
		blog.Errorf("AddApp error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	if len(formData.Get("ApplicationName")) > 32 {
		converter.RespFailV2(common.CCErrAPIServerV2APPNameLenErr, defErr.Error(common.CCErrAPIServerV2APPNameLenErr).Error(), resp)
		return
	}
	lifeCycle := formData.Get("LifeCycle")

	if AppStatusTestI == lifeCycle || AppStatusTest == lifeCycle {
		lifeCycle = AppStatusTestI
	} else if AppStatusOnlineI == lifeCycle || AppStatusOnline == lifeCycle {
		lifeCycle = AppStatusOnlineI
	} else if AppStatusStopI == lifeCycle || AppStatusStop == lifeCycle {
		lifeCycle = AppStatusStopI
	} else {

		converter.RespFailV2(common.CCErrCommParamsIsInvalid, defErr.Errorf(common.CCErrCommParamsIsInvalid, "LifeCycle").Error(), resp)
		return
	}

	param := make(mapstr.MapStr)
	param[common.BKAppNameField] = formData.Get("ApplicationName")
	param[common.BKMaintainersField] = formData.Get("Maintainers")
	param[common.BKLanguageField] = "1"

	timeZone := formData.Get("TimeZone")
	if "" != timeZone {
		param[common.BKTimeZoneField] = timeZone
	} else {
		param[common.BKTimeZoneField] = "Asia/Shanghai"
	}

	//param[common.CreatorField] = formData.Get("Creator")
	param[common.BKLifeCycleField] = lifeCycle
	param[common.BKProductPMField] = formData.Get("ProductPm")
	param[common.BKDeveloperField] = formData.Get("Developer")
	param[common.BKTesterField] = formData.Get("Tester")
	param[common.BKOperatorField] = formData.Get("Operator")

	blog.Infof("AddApp v3 param data1: %v", param)

	param, err = s.Logics.AutoInputV3Field(param, common.BKInnerObjIDApp, user, pheader)

	if nil != err {
		blog.Errorf("AddApp params:%s, error:%v", param, err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	result, err := s.CoreAPI.TopoServer().Instance().CreateApp(context.Background(), user, pheader, param)

	blog.Infof("AddApp  params:%v", param)
	if nil != err {
		blog.Errorf("AddApp  params:%v, error:%s", param, err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("AddApp  error:%s", result.ErrMsg)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error()+" :"+result.ErrMsg, resp)
		return
	}

	rspDataV3Map := result.Data
	resData := make(map[string]interface{})
	resData["appId"] = rspDataV3Map[common.BKAppIDField]
	resData["success"] = true
	converter.RespSuccessV2(resData, resp)
}

func (s *Service) deleteApp(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	user := util.GetUser(pheader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("deleteApp error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	appID := formData.Get("ApplicationID")
	if "" == appID {
		blog.Errorf("deleteApp error ApplicationID empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}
	result, err := s.CoreAPI.TopoServer().Instance().DeleteApp(context.Background(), appID, user, pheader)

	if nil != err {
		blog.Errorf("deleteApp  error:%s", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)

}

func (s *Service) editApp(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	user := util.GetUser(pheader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("editApp error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Infof("editApp formData:%v", formData)
	res, msg := utils.ValidateFormData(formData, []string{
		"ApplicationID",
	})
	if !res {
		blog.Errorf("editApp error: %s", msg)
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
		lifeMap := map[string]bool{AppStatusTest: true, AppStatusOnline: true, AppStatusStop: true}
		if !lifeMap[LifeCycle] {
			if LifeCycle == AppStatusTestI {
				LifeCycle = AppStatusTest
			} else if LifeCycle == AppStatusOnlineI {
				LifeCycle = AppStatusOnline
			} else if LifeCycle == AppStatusStopI {
				LifeCycle = AppStatusStop
			} else {
				converter.RespFailV2(common.CCErrCommParamsIsInvalid, defErr.Errorf(common.CCErrCommParamsIsInvalid, "LifeCycle").Error(), resp)
				return
			}
		}
	}

	param := make(mapstr.MapStr)
	appID := formData.Get("ApplicationID")

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

	blog.Errorf("editApp param:%v", param)

	result, err := s.CoreAPI.TopoServer().Instance().UpdateApp(context.Background(), user, appID, pheader, param)
	if nil != err {
		blog.Errorf("editApp  params:%s, error:%v", param, err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)
}

func (s *Service) getHostAppByCompanyId(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getHostAppByCompanyId error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.V(3).Infof("GetHostAppByCompanyId formData:%v", formData)
	res, msg := utils.ValidateFormData(formData, []string{
		"companyId",
		"ip",
		"platId",
	})
	if !res {
		blog.Errorf("getHostAppByCompanyId error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	param := make(mapstr.MapStr)
	param[common.BKOwnerIDField] = formData.Get("companyId")
	param["ip"] = formData.Get("ip")
	param[common.BKCloudIDField] = formData.Get("platId")

	result, err := s.CoreAPI.HostServer().GetHostAppByCompanyId(context.Background(), pheader, param)
	if nil != err {
		blog.Errorf("getHostAppByCompanyId  params:%v, error:%v", param, err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForHostDataList(result.Result, result.ErrMsg, result.Data)
	if err != nil {
		blog.Error("convert res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}
	blog.Infof("GetHostAppByCompanyId success, data length: %d", len(resDataV2))
	converter.RespSuccessV2(resDataV2, resp)
}
