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
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/api_server/logics/v2/common/converter"
	"configcenter/src/api_server/logics/v2/common/defs"
	"configcenter/src/api_server/logics/v2/common/utils"
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
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	param := &params.SearchParams{Condition: nil}
	result, err := s.CoreAPI.TopoServer().Instance().SearchApp(srvData.ctx, srvData.ownerID, srvData.header, param)
	if err != nil {
		blog.Errorf("getAppList http error, err: %v,rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("getAppList http reply error, err: %#v,rid:%s", result, srvData.rid)
		converter.RespFailV2Error(defErr.New(result.Code, result.ErrMsg), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForAppList(result.Data)

	if err != nil {
		blog.Errorf("convert v3 applist to v2 error:%v,rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}
	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) getAppByID(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getAppByID error:%v,rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.V(5).Infof("getAppByID data:%#v,rid:%s", formData, srvData.rid)

	if len(formData["ApplicationID"]) == 0 {
		blog.Errorf("getAppByID error: ApplicationID is empty!,input:%#v,rid:%s", formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}

	appIDsArr := strings.Split(formData["ApplicationID"][0], ",")

	appIDArr := make([]int64, 0)
	for _, appID := range appIDsArr {
		id, err := strconv.ParseInt(appID, 10, 64)
		if err != nil {
			blog.Errorf("convert str to int error,input:%#v,rid:%s", formData, srvData.rid)
			converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
			return
		}
		appIDArr = append(appIDArr, id)
	}
	param := &params.SearchParams{Condition: mapstr.MapStr{common.BKAppIDField: mapstr.MapStr{common.BKDBIN: appIDArr}}}

	result, err := s.CoreAPI.TopoServer().Instance().SearchApp(srvData.ctx, srvData.ownerID, srvData.header, param)
	if err != nil {
		blog.Errorf("getAppByID http do error, error:%v,input:%#v,params:%#v,rid:%s", err, formData, param, srvData.rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("getAppByID http reply error, reply:%#v,input:%#v,params:%#v,rid:%s", result, formData, param, srvData.rid)
		converter.RespFailV2Error(defErr.New(result.Code, result.ErrMsg), resp)
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
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getAppByUin error:%v,rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.V(5).Infof("getAppByUin data:%#v,rid:%s", formData, srvData.rid)

	if len(formData["userName"]) == 0 || formData["userName"][0] == "" {
		blog.Errorf("get app by uin error: userName is empty!,input:%#v,rid:%s", formData, srvData.rid)
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

	result, err := s.CoreAPI.TopoServer().Instance().SearchApp(srvData.ctx, srvData.ownerID, srvData.header, param)
	if err != nil {
		blog.Errorf("getAppByUin http do error, err:%v,input:%#v,http input:%#v,rid:%s", err, formData, param, srvData.rid)
		converter.RespFailV2Error(defErr.Error(common.CCErrCommHTTPDoRequestFailed), resp)
		return
	}
	if !result.Result {
		blog.Errorf("getAppByUin http reply error, reply:%v,input:%#v,http input:%#v,rid:%s", result, formData, param, srvData.rid)
		converter.RespFailV2Error(defErr.New(result.Code, result.ErrMsg), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForAppList(result.Data)
	if err != nil {
		blog.Errorf("convert res to v2 error:%v,input:%#v,http input:%#v,rid:%s", err, formData, param, srvData.rid)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) getUserRoleApp(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("GetUserRoleApp error:%v,rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	res, msg := utils.ValidateFormData(formData, []string{"uin", "roleList"})

	if !res {
		blog.Errorf("GetUserRoleApp error: %s,input:%#v,rid:%s", msg, formData, srvData.rid)
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

	result, err := s.CoreAPI.TopoServer().Instance().SearchApp(srvData.ctx, srvData.ownerID, srvData.header, param)
	if err != nil {
		blog.Errorf("GetUserRoleApp http do error. err:%v,input:%#v,http input:%#v,rid:%s", err, formData, param, srvData.rid)
		converter.RespFailV2Error(defErr.Error(common.CCErrCommHTTPDoRequestFailed), resp)
		return
	}
	if !result.Result {
		blog.Errorf("GetUserRoleApp http reply error. reply:%v,input:%#v,http input:%#v,rid:%s", result, formData, param, srvData.rid)
		converter.RespFailV2Error(defErr.New(result.Code, result.ErrMsg), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForRoleApp(result.Data, converter.DecorateUserName(userName), roleArr)
	if err != nil {
		blog.Errorf("convert res to v2 error:%v,input:%#v,http input:%#v,rid:%s", err, formData, param, srvData.rid)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) getAppSetModuleTreeByAppId(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getAppSetModuleTreeByAppId error:%v, rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	if len(formData["ApplicationID"]) == 0 {
		blog.Errorf("getAppSetModuleTreeByAppId error error: ApplicationID is empty!input:%#v,rid:%s", formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	intAppID, err := util.GetInt64ByInterface(formData["ApplicationID"][0])
	if nil != err {
		blog.Errorf("getAppSetModuleTreeByAppId   appID:%v, error:%v, input:%+v,rid:%s", formData["ApplicationID"][0], err, formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	conds := mapstr.MapStr{}
	if 0 < len(formData["ModuleType"]) {
		conds[common.BKModuleTypeField] = formData["ModuleType"][0]
	}

	topo, err := srvData.lgc.GetAppTopo(srvData.ctx, intAppID, conds)
	if err != nil {
		converter.RespFailV2Error(err, resp)
		return
	}

	if nil != topo {
		srvData.lgc.SetModuleHostCount(srvData.ctx, []mapstr.MapStr{topo})
	} else {
		converter.RespSuccessV2(make(map[string]interface{}), resp)
		return
	}
	converter.RespSuccessV2(topo, resp)
}

func (s *Service) addApp(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("AddApp error:%v,rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form

	blog.V(5).Infof("AddApp formData:%v,rid:%s", formData, srvData.rid)
	res, msg := utils.ValidateFormData(formData, []string{
		"ApplicationName",
		"Maintainers",
		"Creator",
		"LifeCycle",
	})
	if !res {
		blog.Errorf("AddApp error: %s,input:%#v,rid:%s", msg, formData, srvData.rid)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	if len(formData.Get("ApplicationName")) > 32 {
		blog.Errorf("AddApp ApplicationName > 32,input:%#v,rid:%s", formData, srvData.rid)
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

	blog.V(5).Infof("AddApp  params:%#v,input:%#v,rid:%s", param, formData, srvData.rid)

	param, err = srvData.lgc.AutoInputV3Field(srvData.ctx, param, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("AddApp AutoInputV3Field error.params:%#v, error:%v,rid:%s", param, err.Error(), srvData.rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	result, err := s.CoreAPI.TopoServer().Instance().CreateApp(srvData.ctx, srvData.ownerID, srvData.header, param)
	if nil != err {
		blog.Errorf("AddApp http do error. input:%#vparams:%#v, error:%s,rid:%s", formData, param, err.Error(), srvData.rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("AddApp http reply error, reply:%#v,input:%#v,param:%#v,rid:%s", result, formData, param, srvData.rid)
		converter.RespFailV2Error(defErr.New(result.Code, result.ErrMsg), resp)
		return
	}

	rspDataV3Map := result.Data
	resData := make(map[string]interface{})
	resData["appId"] = rspDataV3Map[common.BKAppIDField]
	resData["success"] = true
	converter.RespSuccessV2(resData, resp)
}

func (s *Service) deleteApp(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("deleteApp error:%v,rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	appID := formData.Get("ApplicationID")
	if "" == appID {
		blog.Errorf("deleteApp error ApplicationID empty.input:%#v,rid:%s", formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}
	result, err := s.CoreAPI.TopoServer().Instance().DeleteApp(srvData.ctx, appID, srvData.ownerID, srvData.header)
	if nil != err {
		blog.Errorf("deleteApp  error:%s,input:%#v,rid:%s", err, formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("deleteApp  error:%s,input:%#v,rid:%s", result, formData, srvData.rid)
		converter.RespFailV2Error(defErr.New(result.Code, result.ErrMsg), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)

}

func (s *Service) editApp(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("editApp error:%v,rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.V(5).Infof("editApp input:%#v,rid:%s", formData, srvData.rid)
	res, msg := utils.ValidateFormData(formData, []string{
		"ApplicationID",
	})
	if !res {
		blog.Errorf("editApp error: %s,input:%#v,rid:%s", msg, formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}
	appName := formData.Get("ApplicationName")

	if "" != appName && len(appName) > 32 {
		blog.Errorf("editApp ApplicationName > 32,input:%#v,rid:%s", formData, srvData.rid)
		converter.RespFailV2(common.CCErrAPIServerV2APPNameLenErr, defErr.Error(common.CCErrAPIServerV2APPNameLenErr).Error(), resp)
		return
	}

	LifeCycle := formData.Get("LifeCycle")

	if "" != LifeCycle {
		if LifeCycle == AppStatusTestI || LifeCycle == AppStatusTest {
			LifeCycle = AppStatusTestI
		} else if LifeCycle == AppStatusOnlineI || LifeCycle == AppStatusOnline {
			LifeCycle = AppStatusOnlineI
		} else if LifeCycle == AppStatusStopI || LifeCycle == AppStatusStop {
			LifeCycle = AppStatusStopI
		} else {
			blog.Errorf("editApp  lifeCycle value bad,input:%#v,rid:%s", formData, srvData.rid)
			converter.RespFailV2(common.CCErrCommParamsIsInvalid, defErr.Errorf(common.CCErrCommParamsIsInvalid, "LifeCycle").Error(), resp)
			return
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

	result, err := s.CoreAPI.TopoServer().Instance().UpdateApp(srvData.ctx, srvData.ownerID, appID, srvData.header, param)
	if nil != err {
		blog.Errorf("editApp http do error. error:%v,params:%#v,input:%#v,rid:%s", err, param, formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("editApp http reply error. reply:%s,input:%#v,rid:%s", result, formData, srvData.rid)
		converter.RespFailV2Error(defErr.New(result.Code, result.ErrMsg), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)
}

func (s *Service) getHostAppByCompanyId(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getHostAppByCompanyId error:%v,rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.V(5).Infof("GetHostAppByCompanyId formData:%v,rid:%s", formData, srvData.rid)
	res, msg := utils.ValidateFormData(formData, []string{
		"companyId",
		"ip",
		"platId",
	})
	if !res {
		blog.Errorf("getHostAppByCompanyId error: %s,input:%#v,rid:%s", msg, formData, srvData.rid)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	param := make(mapstr.MapStr)
	param[common.BKOwnerIDField] = formData.Get("companyId")
	param["ip"] = formData.Get("ip")
	param[common.BKCloudIDField] = formData.Get("platId")

	result, err := s.CoreAPI.HostServer().GetHostAppByCompanyId(srvData.ctx, srvData.header, param)
	if nil != err {
		blog.Errorf("getHostAppByCompanyId http do error. error:%v,input:%#v,param:%#v,rid:%s", err, formData, param, srvData.rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("getHostAppByCompanyId http do error. reply:%#v,input:%#v,param:%#v,rid:%s", result, formData, param, srvData.rid)
		converter.RespFailV2Error(defErr.New(result.Code, result.ErrMsg), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForHostDataList(result.Result, result.ErrMsg, result.Data)
	if err != nil {
		blog.Errorf("convert res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}
	blog.V(5).Infof("GetHostAppByCompanyId success, data length: %d,rid:%s", len(resDataV2), srvData.rid)
	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) getAppByOwnerAndUin(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getHostListByOwner error:%s, rid:%s", err.Error(), srvData.rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	res, msg := utils.ValidateFormData(formData, []string{"Owner_uin", "Uin"})
	if !res {
		blog.Errorf("getHostListByOwner error: %s,rid:%s", msg, srvData.rid)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	ownerID := formData["Owner_uin"][0]
	uin := formData["Uin"][0]
	// reset http header owner
	srvData.header.Set(common.BKHTTPOwnerID, ownerID)
	srvData.ownerID = ownerID
	appInfoArr, errCode, err := srvData.lgc.GetAppListByOwnerIDAndUin(srvData.ctx, uin)
	if nil != err {
		blog.Errorf("getHostListByOwner error:%s, input:%+v, rid:%s", err.Error(), formData, srvData.rid)
		converter.RespFailV2(errCode, err.Error(), resp)
		return
	}
	converter.RespSuccessV2(converter.GeneralV2Data(appInfoArr), resp)

}
