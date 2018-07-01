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
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/emicklei/go-restful"

	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
)

func (s *Service) getProcessPortByApplicationID(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	user := util.GetUser(pheader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	defLang := s.Language.CreateDefaultCCLanguageIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("GetProcessPortByApplicationID error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("GetProcessPortByApplicationID http body data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID"})
	if !res {
		blog.Error("GetProcessPortByApplicationID error: %s", msg)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}

	appID := formData.Get("ApplicationID")
	if nil != err {

		blog.Error("GetProcessPortByApplicationID error: ApplicationID is not number")
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	modules, err := s.getModulesByAppId(appID, user, pheader)
	blog.Debug("modules:%v", modules)
	if nil != err {
		blog.Error("getModulesByAppId error:%v", err)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, err.Error()).Error(), resp)
		return
	}
	modulesMap := modules.(map[string]interface{})

	result, err := s.CoreAPI.ProcServer().OpenAPI().GetProcessPortByApplicationID(context.Background(), appID, pheader, modulesMap)
	if err != nil {
		blog.Error("GetProcessPortByApplicationID  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	if !result.Result {
		code, _ := util.GetIntByInterface(result.Code)
		blog.Error("GetProcessPortByApplicationID error:%s", result.ErrMsg)
		converter.RespFailV2(code, result.ErrMsg, resp)
		return
	}
	if nil == result.Data {
		emptyData := make([]interface{}, 0)
		converter.RespSuccessV2(converter.ResV2ToForProcList(emptyData, defLang), resp)
		return
	}
	converter.RespSuccessV2(converter.GeneralV2Data(result.Data), resp)

}

func (s *Service) getProcessPortByIP(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	defLang := s.Language.CreateDefaultCCLanguageIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("GetProcessPortByIP error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Debug("GetProcessPortByIP http body data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ips"})
	if !res {
		blog.Error("GetProcessPortByIP error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	ips := formData.Get("ips")
	ipArr := strings.Split(ips, ",")
	if len(ipArr) == 0 {
		blog.Error("GetProcessPortByIP error: ips is required")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ips").Error(), resp)
		return
	}
	param := make(common.KvMap)
	param["ipArr"] = ipArr
	result, err := s.CoreAPI.ProcServer().OpenAPI().GetProcessPortByIP(context.Background(), pheader, param)
	if err != nil {
		blog.Error("GetProcessPortByIP url error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	if !result.Result {
		blog.Error("GetProcessPortByIP error:%s", result.ErrMsg)
		converter.RespFailV2(result.Code, result.ErrMsg, resp)
		return
	}
	if nil == result.Data {
		emptyData := make([]interface{}, 0)
		converter.RespSuccessV2(converter.ResV2ToForProcList(emptyData, defLang), resp)
		return
	}

	converter.RespSuccessV2(converter.ResV2ToForProcList(result.Data, defLang), resp)
}

func (s *Service) getModulesByAppId(appID string, user string, pheader http.Header) (interface{}, error) {

	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appID

	searchParams := make(map[string]interface{})
	searchParams["condition"] = condition
	searchParams["fields"] = fmt.Sprintf("%s,%s", common.BKModuleIDField, common.BKModuleNameField)
	searchParams["start"] = 0
	searchParams["limit"] = 0
	searchParams["sort"] = ""
	result, err := s.CoreAPI.TopoServer().OpenAPI().SearchModuleByApp(context.Background(), appID, pheader, searchParams)
	if nil != err {
		blog.Errorf("getModulesByAppId error:%v", err)
		return nil, err
	}

	resData := make([]mapstr.MapStr, 0)
	if result.Result {
		modules := (result.Data.(mapstr.MapStr))["info"]
		for _, module := range modules.([]interface{}) {
			resData = append(resData, module.(mapstr.MapStr))
		}
		return resData, nil
	} else {
		return nil, errors.New(result.ErrMsg)
	}

	return nil, nil
}
