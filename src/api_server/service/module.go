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
	"strconv"
	"strings"

	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

const (
	ModuleTypeCommon = "普通"
	ModuleTypeDB     = "数据库"
)

const (
	ModuleTypeCommonI = "1"
	ModuleTypeDBI     = "2"
)

func (s *Service) getModulesByApp(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getModulesByApp error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	if len(formData["ApplicationID"]) == 0 || formData["ApplicationID"][0] == "" {
		blog.Errorf("getModulesByApp  error: ApplicationID is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}

	appID := formData["ApplicationID"][0]

	params := mapstr.MapStr{
		"fields":    []string{},
		"condition": mapstr.MapStr{},
		"page": mapstr.MapStr{
			"start": 0,
			"limit": 0,
		},
	}

	result, err := s.CoreAPI.TopoServer().OpenAPI().SearchModuleByApp(context.Background(), appID, pheader, params)

	if err != nil {
		blog.Errorf("getModulesByApp   error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	resDataV2, err := converter.ResToV2ForModuleMapList(result.Result, result.ErrMsg, result.Data)
	if err != nil {
		blog.Error("convert module res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) updateModule(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	defLang := s.Language.CreateDefaultCCLanguageIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("updateModule error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("updateModule data: %s", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "ModuleID"})
	if !res {
		blog.Errorf("updateModule error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	reqParam := make(mapstr.MapStr)
	reqData := make(mapstr.MapStr)

	moduleIDStrArr := strings.Split(formData["ModuleID"][0], ",")
	moduleIDArr, err := utils.SliceStrToInt(moduleIDStrArr)
	if nil != err {
		blog.Errorf("updateModule error:%v", err)
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
	if "" != ModuleType {
		if ModuleType == "1" {
		} else if ModuleType == "2" {

		} else {
			msg := defLang.Language("apiv2_module_type_error")
			blog.Errorf("updateModule error:%v", msg)
			converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
			return
		}

		reqData[common.BKModuleTypeField] = ModuleType
	}

	if len(moduleIDArr) == 1 && len(formData["ModuleName"]) > 0 {
		reqData[common.BKModuleNameField] = moduleName
	} else {
		msg := defLang.Language("apiv2_module_edit_multi_module_name")
		blog.Infof("updateModule error:%v", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	if len(moduleName) > 24 {
		msg := defLang.Language("apiv2_module_name_lt_24")
		blog.Infof("updateModule error:%v", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	reqParam["data"] = reqData

	result, err := s.CoreAPI.TopoServer().OpenAPI().UpdateMultiModule(context.Background(), appID, pheader, reqParam)

	if err != nil {
		blog.Errorf("updateModule  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)
}

func (s *Service) addModule(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	defLang := s.Language.CreateDefaultCCLanguageIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("addModule error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("addModule data: %v", formData)

	res, msg := utils.ValidateFormData(formData, []string{
		"ApplicationID",
		"SetID",
		"ModuleName",
		"Operator",
		"BakOperator",
		"ModuleType",
	})
	if !res {
		blog.Errorf("addModule error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	reqParam := mapstr.MapStr{}
	moduleName := formData.Get("ModuleName")
	appId, _ := strconv.Atoi(formData.Get("ApplicationID"))
	setId, _ := strconv.Atoi(formData.Get("SetID"))
	Operator := formData.Get("Operator")
	BakOperator := formData.Get("BakOperator")
	moduleType := formData.Get("ModuleType")

	moduleTypeMap := map[string]string{ModuleTypeCommonI: ModuleTypeCommon, ModuleTypeDBI: ModuleTypeDB}
	if ModuleTypeCommonI == moduleType && ModuleTypeCommonI != moduleType {
		if moduleType == ModuleTypeCommon || moduleType == ModuleTypeDB {
			reqParam[common.BKModuleTypeField] = moduleType
		} else {
			msg = defLang.Language("apiv2_module_type_error")
			blog.Errorf("addModule error: %s", msg)
			converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
			return
		}
	} else {
		reqParam[common.BKModuleTypeField] = moduleTypeMap[moduleType]
	}

	if "1" != moduleType && "2" != moduleType {
		msg = defLang.Language("apiv2_module_type_error")
		blog.Errorf("addModule error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	reqParam[common.BKModuleNameField] = moduleName
	reqParam[common.BKAppIDField] = appId
	reqParam[common.BKSetIDField] = setId
	reqParam[common.BKOperatorField] = Operator
	reqParam[common.BKBakOperatorField] = BakOperator

	result, err := s.CoreAPI.TopoServer().OpenAPI().AddMultiModule(context.Background(), pheader, reqParam)
	if err != nil {
		blog.Errorf("addModule  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	rspDataV3Map := result.Data.(map[string]interface{})
	resData := make(map[string]interface{})
	resData["ModuleID"] = rspDataV3Map[common.BKModuleIDField]
	resData["success"] = true
	converter.RespSuccessV2(resData, resp)

}

func (s *Service) deleteModule(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("deleteModule error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("deleteModule data: %v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "ModuleID"})
	if !res {
		blog.Errorf("deleteModule error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	reqParam := make(map[string]interface{})
	moduleIdStrArr := strings.Split(formData["ModuleID"][0], ",")
	moduleIdArr, err := utils.SliceStrToInt(moduleIdStrArr)
	if nil != err {
		blog.Errorf("deleteModule error:%v", err)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, err.Error()).Error(), resp)
		return
	}

	appID := formData["ApplicationID"][0]

	reqParam[common.BKModuleIDField] = moduleIdArr

	result, err := s.CoreAPI.TopoServer().OpenAPI().DeleteMultiModule(context.Background(), appID, pheader, reqParam)

	if err != nil {
		blog.Errorf("deleteModule  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)
}
