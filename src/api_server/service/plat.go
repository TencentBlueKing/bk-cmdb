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
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func (s *Service) updateHost(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("updateHost error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("updateHost data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"appId", "orgPlatId", "ip", "dstPlatId"})
	if !res {
		blog.Errorf("ValidateFormData error:%s", msg)
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

	result, err := s.CoreAPI.HostServer().UpdateHost(context.Background(), appID, pheader, param)
	if err != nil {
		blog.Errorf("updateHost u error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)

}

func (s *Service) getPlats(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	result, err := s.CoreAPI.HostServer().GetPlat(context.Background(), pheader)
	if err != nil {
		blog.Errorf("getPlats error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForPlatList(result.Result, result.ErrMsg, result.Data)
	if err != nil {
		blog.Error("convert plat res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) deletePlats(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("deletePlats error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("deletePlats data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"platId"})
	if !res {
		blog.Error("ValidateFormData error:%s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	platID := formData["platId"][0]

	result, err := s.CoreAPI.HostServer().DelPlat(context.Background(), platID, pheader)
	if err != nil {
		blog.Errorf("deletePlats error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)
}

func (s *Service) createPlats(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	user := util.GetUser(pheader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("createPlats error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("createPlats data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"platName"})
	if !res {
		blog.Error("ValidateFormData error:%s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	platName := formData["platName"][0]
	param := map[string]interface{}{
		common.BKCloudNameField: platName,
	}

	param, _ = s.Logics.AutoInputV3Field(param, common.BKInnerObjIDPlat, user, pheader)

	result, err := s.CoreAPI.HostServer().CreatePlat(context.Background(), pheader, param)
	if nil != err {
		blog.Errorf("createPlats  error:%v", err)
		if strings.Contains(err.Error(), strconv.Itoa(common.CCErrCommDuplicateItem)) {
			converter.RespFailV2(common.CCErrCommDuplicateItem, defErr.Error(common.CCErrCommDuplicateItem).Error(), resp)
			return
		}
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	if !result.Result {
		blog.Errorf("createPlats error:%s", result.ErrMsg)
		converter.RespFailV2(result.Code, result.ErrMsg, resp)
		return
	}

	rspDataV3Map, ok := result.Data.(map[string]interface{})
	if false == ok {
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(rspDataV3Map[common.BKCloudIDField], resp)
}
