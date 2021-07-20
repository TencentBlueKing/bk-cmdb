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
	"strconv"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func (s *service) AuthVerify(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	ownerID := util.GetOwnerID(pheader)
	rid := util.GetHTTPCCRequestID(pheader)

	if auth.EnableAuthorize() == false {
		blog.Errorf("inappropriate calling, auth is disabled, rid: %s", rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommInappropriateVisitToIAM)})
		return
	}

	body := metadata.AuthBathVerifyRequest{}
	if err := json.NewDecoder(req.Request.Body).Decode(&body); err != nil {
		blog.Errorf("get user's resource auth verify status, but decode body failed, err: %v, rid: %s", err, rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	user := meta.UserInfo{
		UserName:        util.GetUser(pheader),
		SupplierAccount: ownerID,
	}

	resources := make([]metadata.AuthBathVerifyResult, len(body.Resources), len(body.Resources))

	attrs := make([]meta.ResourceAttribute, 0)
	needExactAuthAttrs := make([]meta.ResourceAttribute, 0)
	needExactAuthMap := make(map[int]bool)

	for i, res := range body.Resources {
		resources[i].AuthResource = res
		attr := meta.ResourceAttribute{
			Basic: meta.Basic{
				Type:         meta.ResourceType(res.ResourceType),
				Action:       meta.Action(res.Action),
				InstanceID:   res.ResourceID,
				InstanceIDEx: res.ResourceIDEx,
			},
			SupplierAccount: ownerID,
			BusinessID:      res.BizID,
		}
		for _, item := range res.ParentLayers {
			attr.Layers = append(attr.Layers, meta.Item{Type: meta.ResourceType(item.ResourceType), InstanceID: item.ResourceID, InstanceIDEx: item.ResourceIDEx})
		}
		// contains exact resource info, need exact authorize
		if res.ResourceID > 0 || res.ResourceIDEx != "" || res.BizID > 0 || len(res.ParentLayers) > 0 {
			needExactAuthMap[i] = true
			needExactAuthAttrs = append(needExactAuthAttrs, attr)
		} else {
			attrs = append(attrs, attr)
		}
	}

	ctx := context.WithValue(context.Background(), common.ContextRequestIDField, rid)

	if len(needExactAuthAttrs) > 0 {
		verifyResults, err := s.authorizer.AuthorizeBatch(ctx, pheader, user, needExactAuthAttrs...)
		if err != nil {
			blog.ErrorJSON("get user's resource auth verify status, but authorize batch failed, err: %s, attrs: %s, rid: %s", err, needExactAuthAttrs, rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrAPIGetUserResourceAuthStatusFailed)})
			return
		}
		index := 0
		resourceLen := len(body.Resources)
		for i := 0; i < resourceLen; i++ {
			if needExactAuthMap[i] {
				resources[i].Passed = verifyResults[index].Authorized
				index++
			}
		}
	}

	if len(attrs) > 0 {
		verifyResults, err := s.authorizer.AuthorizeAnyBatch(ctx, pheader, user, attrs...)
		if err != nil {
			blog.ErrorJSON("get user's resource auth verify status, but authorize any batch failed, err: %s, attrs: %s, rid: %s", err, attrs, rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrAPIGetUserResourceAuthStatusFailed)})
			return
		}
		index := 0
		resourceLen := len(body.Resources)
		for i := 0; i < resourceLen; i++ {
			if !needExactAuthMap[i] {
				resources[i].Passed = verifyResults[index].Authorized
				index++
			}
		}
	}

	resp.WriteEntity(metadata.NewSuccessResp(resources))
}

func (s *service) GetAnyAuthorizedAppList(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	rid := util.GetHTTPCCRequestID(pheader)

	if auth.EnableAuthorize() == false {
		blog.Errorf("inappropriate calling, auth is disabled, rid: %s", rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommInappropriateVisitToIAM)})
		return
	}

	userInfo := meta.UserInfo{
		UserName:        util.GetUser(pheader),
		SupplierAccount: util.GetOwnerID(pheader),
	}

	authInput := meta.ListAuthorizedResourcesParam{
		UserName:     util.GetUser(pheader),
		ResourceType: meta.Business,
		Action:       meta.ViewBusinessResource,
	}
	authorizedResources, err := s.authorizer.ListAuthorizedResources(req.Request.Context(), pheader, authInput)
	if err != nil {
		blog.Errorf("get user: %s authorized business list failed, err: %v, rid: %s", userInfo.UserName, err, rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrAPIGetAuthorizedAppListFromAuthFailed)})
		return
	}
	appIDList := make([]int64, 0)
	for _, resourceID := range authorizedResources {
		bizID, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			blog.Errorf("parse bizID(%s) failed, err: %v, rid: %s", bizID, err, rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
			return
		}
		appIDList = append(appIDList, bizID)
	}

	if len(appIDList) == 0 {
		resp.WriteEntity(metadata.NewSuccessResp(metadata.InstResult{Info: make([]mapstr.MapStr, 0)}))
		return
	}

	input := params.SearchParams{
		Condition: mapstr.MapStr{common.BKAppIDField: mapstr.MapStr{"$in": appIDList}},
	}

	result, err := s.engine.CoreAPI.TopoServer().Instance().SearchApp(req.Request.Context(), userInfo.SupplierAccount, req.Request.Header, &input)
	if err != nil {
		blog.Errorf("get authorized business list, but get apps[%v] failed, err: %v, rid: %s", appIDList, err, rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrAPIGetAuthorizedAppListFromAuthFailed)})
		return
	}

	if !result.Result {
		blog.Errorf("get authorized business list, but get apps[%v] failed, err: %v, rid: %s", appIDList, result.ErrMsg, rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrAPIGetAuthorizedAppListFromAuthFailed)})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(result.Data))
}

// GetUserNoAuthSkipURL returns the url which can helps to launch the bk-auth-center when a user do not
// have the authorize to access resource(s).
func (s *service) GetUserNoAuthSkipURL(req *restful.Request, resp *restful.Response) {
	reqHeader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(reqHeader))
	rid := util.GetHTTPCCRequestID(reqHeader)

	p := new(metadata.IamPermission)
	err := json.NewDecoder(req.Request.Body).Decode(p)
	if err != nil {
		blog.Errorf("get user's skip url when no auth, but decode request failed, err: %v, rid: %s", err, rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	url, err := s.authorizer.GetNoAuthSkipUrl(req.Request.Context(), reqHeader, p)
	if err != nil {
		blog.Errorf("get user's skip url when no auth, but request to auth center failed, err: %v, rid: %s", err, rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrGetNoAuthSkipURLFailed)})
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(url))
}
