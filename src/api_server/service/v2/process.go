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
	"strings"

	"configcenter/src/api_server/logics/v2/common/converter"
	"configcenter/src/api_server/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccError "configcenter/src/common/errors"
	"configcenter/src/common/mapstr"

	"github.com/emicklei/go-restful"
)

func (s *Service) getProcessPortByApplicationID(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr
	defLang := srvData.ccLang

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getProcessPortByApplicationID error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("getProcessPortByApplicationID  data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID"})
	if !res {
		blog.Errorf("getProcessPortByApplicationID error: %s", msg)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}

	appID := formData.Get("ApplicationID")
	if nil != err {

		blog.Error("getProcessPortByApplicationID error: ApplicationID is not number")
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	modulesMap, err := s.getModulesByAppId(srvData.ctx, appID, srvData)
	blog.V(5).Infof("modules data:%v,input:%+v,rid:%s", modulesMap, formData, srvData.rid)
	if nil != err {
		blog.Errorf("getProcessPortByApplicationID error:%v,input:%+v,rid:%s", err, formData, srvData.rid)
		converter.RespFailV2Error(err, resp)
		return
	}

	result, err := s.CoreAPI.ProcServer().OpenAPI().GetProcessPortByApplicationID(srvData.ctx, appID, srvData.header, modulesMap)
	if err != nil {
		blog.Errorf("getProcessPortByApplicationID http do error.err:%v,input:%#v,rid:%s", err, formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	if !result.Result {
		blog.Errorf("getProcessPortByApplicationID http reply error.reply:%#v, input:%+v,rid:%s", result, formData, srvData.rid)
		converter.RespFailV2(result.Code, result.ErrMsg, resp)
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
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr
	defLang := srvData.ccLang

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getProcessPortByIP error:%v,rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Infof("getProcessPortByIP data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ips"})
	if !res {
		blog.Errorf("getProcessPortByIP error: %v,input:%#v,rid:%s", msg, formData, srvData.rid)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	ips := formData.Get("ips")
	ipArr := strings.Split(ips, ",")
	if len(ipArr) == 0 {
		blog.Errorf("getProcessPortByIP error: ips is required,input:%#v,rid:%s", formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ips").Error(), resp)
		return
	}
	param := make(common.KvMap)
	param["ipArr"] = ipArr
	result, err := s.CoreAPI.ProcServer().OpenAPI().GetProcessPortByIP(srvData.ctx, srvData.header, param)
	if err != nil {
		blog.Errorf("getProcessPortByIP http do error.err:%v,input:%#v,rid:%s", err, formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	if !result.Result {
		blog.Errorf("getProcessPortByIP http reply error.reply:%#v,input:%#v,rid:%s", result, formData, srvData.rid)
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

func (s *Service) getModulesByAppId(ctx context.Context, appID string, srvData *srvComm) ([]mapstr.MapStr, ccError.CCError) {

	searchParams := mapstr.MapStr{
		"condition": mapstr.MapStr{},
		"fields":    []string{common.BKModuleIDField, common.BKModuleNameField},
		"page": mapstr.MapStr{
			"start": 0,
			"limit": 0,
			"sort":  "",
		},
	}
	result, err := s.CoreAPI.TopoServer().OpenAPI().SearchModuleByApp(ctx, appID, srvData.header, searchParams)
	if nil != err {
        blog.Errorf("getModulesByAppId http do error,err:%v,appID:%v,rid:%s", err.Error(), appID, srvData.rid)
		return nil, srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("getModulesByAppId http do error,err code:%d, err msg:%v,appID:%v,rid:%s", result.Code, result.ErrMsg, appID, srvData.rid)
		return nil, srvData.ccErr.New(result.Code, result.ErrMsg)
	}

	resData := make([]mapstr.MapStr, 0)
	for _, module := range result.Data.Info {
		resData = append(resData, module)
	}
	return resData, nil

}
