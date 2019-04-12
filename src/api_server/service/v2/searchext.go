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
	"github.com/emicklei/go-restful"

	"configcenter/src/api_server/logics/v2/common/converter"
	"configcenter/src/api_server/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/blog"
)

func (s *Service) getHostListByOwner(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getHostListByOwner error:%v,rid:%s", err, srvData.rid)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	res, msg := utils.ValidateFormData(formData, []string{"OwnerQQ"})
	if !res {
		blog.Errorf("getHostListByOwner error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	//  reset http header owner
	srvData.header.Set(common.BKHTTPOwnerID, formData["OwnerQQ"][0])
	srvData.ownerID = formData["OwnerQQ"][0]
	dataArr, dataErr := srvData.lgc.GetAllHostAndModuleRelation(srvData.ctx)
	if nil != dataErr {
		blog.Errorf("getHostListByOwner error: %s,input:%#v,rid:%s", dataErr.Error(), formData, srvData.rid)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, dataErr.Error(), resp)
		return
	}

	converter.RespSuccessV2(converter.GeneralV2Data(dataArr), resp)

}
