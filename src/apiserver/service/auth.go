/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
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
	"net/http"

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func (s *service) AuthVerify(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)

	body := metadata.AuthBathVerifyRequest{}
	if err := json.NewDecoder(req.Request.Body).Decode(&body); err != nil {
		blog.Errorf("add subscription, but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	user := meta.UserInfo{
		UserName:        util.GetUser(pheader),
		SupplierAccount: ownerID,
	}

	resources := make([]metadata.AuthBathVerifyResult, len(body.Resources), len(body.Resources))

	attrs := make([]meta.ResourceAttribute, len(body.Resources), len(body.Resources))
	for i, res := range body.Resources {
		resources[i].AuthResource = res
		attrs[i].BusinessID = res.BizID
		attrs[i].SupplierAccount = ownerID
		attrs[i].Type = meta.ResourceType(res.ResourceType)
		attrs[i].InstanceID = res.ResourceID
		attrs[i].Action = meta.Action(res.Action)
		attrs[i].SupplierAccount = ownerID
		for _, item := range res.ParentLayers {
			attrs[i].Layers = append(attrs[i].Layers, meta.Item{Type: meta.ResourceType(item.ResourceType), InstanceID: item.ResourceID})
		}
	}

	verifyResults, err := s.authorizer.AuthorizeBatch(context.Background(), user, attrs...)
	if err != nil {
		blog.Errorf("add subscription, but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	for i, verifyResult := range verifyResults {
		resources[i].Passed = verifyResult.Authorized
		resources[i].Reason = verifyResult.Reason
	}

	resp.WriteAsJson(metadata.NewSuccessResp(resources))

}
