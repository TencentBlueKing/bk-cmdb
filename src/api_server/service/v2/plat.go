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
	"strconv"
	"strings"

	"configcenter/src/api_server/logics/v2/common/converter"
	"configcenter/src/api_server/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func (s *Service) updateHost(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("updateHost error:%v,rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.V(5).Infof("updateHost data:%#v,rid:%s", formData, srvData.rid)

	res, msg := utils.ValidateFormData(formData, []string{"appId", "orgPlatId", "ip", "dstPlatId"})
	if !res {
		blog.Errorf("ValidateFormData error:%s,input:%#v,rid:%s", msg, formData, srvData.rid)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	appID := formData["appId"][0]
	orgPlatID, _ := util.GetIntByInterface(formData["orgPlatId"][0])
	ip := strings.Trim(formData["ip"][0], " ")
	dstPlatID, _ := util.GetIntByInterface(formData["dstPlatId"][0])

	param := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKHostInnerIPField: ip,
			common.BKSubAreaField:     orgPlatID,
		},
		"data": map[string]interface{}{
			common.BKSubAreaField: dstPlatID,
		},
	}

	result, err := s.CoreAPI.HostServer().UpdateHost(srvData.ctx, appID, srvData.header, param)
	if err != nil {
		blog.Errorf("updateHost http do error.err:%v,input:%#v,rid:%s", err, formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("updateHost http reply error.reply:%#v,input:%#v,rid:%s", result, formData, srvData.rid)
		converter.RespFailV2(result.Code, result.ErrMsg, resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)

}

func (s *Service) getPlats(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	result, err := s.CoreAPI.HostServer().GetPlat(srvData.ctx, srvData.header)
	if err != nil {
		blog.Errorf("getPlats error:%v, rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("getPlats http response error, err code:%d, err msg:%s,rid:%s", result.Code, result.ErrMsg, srvData.rid)
		converter.RespFailV2(result.Code, result.ErrMsg, resp)
		return
	}

	resDataV2, err := converter.ResToV2ForPlatList(result.Data)
	if err != nil {
		blog.Errorf("convert plat res to v2 error:%v, rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) deletePlats(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("deletePlats error:%v,rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.V(5).Infof("deletePlats data:%#v,rid:%s", formData, srvData.rid)

	res, msg := utils.ValidateFormData(formData, []string{"platId"})
	if !res {
		blog.Errorf("ValidateFormData error:%s,input:%#v,rid:%s", msg, formData, srvData.rid)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	platID := formData["platId"][0]

	result, err := s.CoreAPI.HostServer().DelPlat(srvData.ctx, platID, srvData.header)
	if err != nil {
		blog.Errorf("deletePlats http do error.err:%v,input:%#v,rid:%s", err, formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("deletePlats http reply error.err:%#v,input:%#v,rid:%s", result, formData, srvData.rid)
		converter.RespFailV2(result.Code, result.ErrMsg, resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)
}

func (s *Service) createPlats(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("createPlats error:%v,rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.V(5).Infof("createPlats data:%#v,rid:%s", formData, srvData.rid)

	res, msg := utils.ValidateFormData(formData, []string{"platName"})
	if !res {
		blog.Errorf("ValidateFormData error:%s,input:%#v,rid:%s", msg, formData, srvData.rid)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	platName := formData["platName"][0]
	param := mapstr.MapStr{
		common.BKCloudNameField: platName,
	}

	param, err = srvData.lgc.AutoInputV3Field(srvData.ctx, param, common.BKInnerObjIDPlat)
	if err != nil {
		blog.Errorf("AutoInputV3Field err,err:%s,input:%#v,rid:%s", err.Error(), formData, srvData.rid)
		converter.RespFailV2Error(err, resp)
		return
	}
	result, err := s.CoreAPI.HostServer().CreatePlat(srvData.ctx, srvData.header, param)
	if nil != err {
		blog.Errorf("createPlats http do error.err:%v,input:%#v,rid:%s", err, formData, srvData.rid)
		if strings.Contains(err.Error(), strconv.Itoa(common.CCErrCommDuplicateItem)) {
			converter.RespFailV2(common.CCErrCommDuplicateItem, defErr.Error(common.CCErrCommDuplicateItem).Error(), resp)
			return
		}
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("createPlats http reply error.reply:%#v,input:%#v,rid:%s", result, formData, srvData.rid)
		converter.RespFailV2(result.Code, result.ErrMsg, resp)
		return
	}

	rspDataV3Map, ok := result.Data.(map[string]interface{})
	if false == ok {
		blog.Errorf("createPlats http reply error.data not ma[string]interface.reply:%#v,input:%#v,rid:%s", result.Data, formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(rspDataV3Map[common.BKCloudIDField], resp)
}
