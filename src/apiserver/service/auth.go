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

	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/meta"
	"configcenter/src/common"
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

	if s.authorizer.Enabled() == false {
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

	attrs := make([]meta.ResourceAttribute, len(body.Resources))
	for i, res := range body.Resources {
		resources[i].AuthResource = res
		attrs[i].BusinessID = res.BizID
		attrs[i].SupplierAccount = ownerID
		attrs[i].Type = meta.ResourceType(res.ResourceType)
		attrs[i].InstanceID = res.ResourceID
		attrs[i].Action = meta.Action(res.Action)
		for _, item := range res.ParentLayers {
			attrs[i].Layers = append(attrs[i].Layers, meta.Item{Type: meta.ResourceType(item.ResourceType), InstanceID: item.ResourceID})
		}
	}

	ctx := context.WithValue(context.Background(), common.ContextRequestIDField, rid)
	verifyResults, err := s.authorizer.AuthorizeBatch(ctx, user, attrs...)
	if err != nil {
		blog.Errorf("get user's resource auth verify status, but authorize batch failed, err: %v, rid: %s", err, rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrAPIGetUserResourceAuthStatusFailed)})
		return
	}

	for i, verifyResult := range verifyResults {
		resources[i].Passed = verifyResult.Authorized
		resources[i].Reason = verifyResult.Reason
	}

	resp.WriteEntity(metadata.NewSuccessResp(resources))
}

func (s *service) GetAdminEntrance(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	rid := util.GetHTTPCCRequestID(pheader)

	if s.authorizer.Enabled() == false {
		blog.Errorf("inappropriate calling, auth is disabled, rid: %s", rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommInappropriateVisitToIAM)})
		return
	}

	userInfo := meta.UserInfo{
		UserName:        util.GetUser(pheader),
		SupplierAccount: util.GetOwnerID(pheader),
	}

	systemlist, err := s.authorizer.AdminEntrance(req.Request.Context(), userInfo)
	if err != nil {
		blog.Errorf("get user: %s authorized business list failed, err: %v, rid: %s", err, rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrAPIGetUserResourceAuthStatusFailed)})
		return
	}

	result := struct {
		Passed bool `json:"is_pass"`
	}{}
	if len(systemlist) > 0 {
		result.Passed = true
	}

	resp.WriteEntity(metadata.NewSuccessResp(result))
}

func (s *service) GetAnyAuthorizedAppList(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	rid := util.GetHTTPCCRequestID(pheader)

	if s.authorizer.Enabled() == false {
		blog.Errorf("inappropriate calling, auth is disabled, rid: %s", rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommInappropriateVisitToIAM)})
		return
	}

	userInfo := meta.UserInfo{
		UserName:        util.GetUser(pheader),
		SupplierAccount: util.GetOwnerID(pheader),
	}

	appIDList, err := s.authorizer.GetAnyAuthorizedBusinessList(req.Request.Context(), userInfo)
	if err != nil {
		blog.Errorf("get user: %s authorized business list failed, err: %v, rid: %s", userInfo.UserName, err, rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrAPIGetAuthorizedAppListFromAuthFailed)})
		return
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

	p := make([]metadata.Permission, 0)
	err := json.NewDecoder(req.Request.Body).Decode(&p)
	if err != nil {
		blog.Errorf("get user's skip url when no auth, but decode request failed, err: %v, rid: %s", err, rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	url, err := s.authorizer.GetNoAuthSkipUrl(req.Request.Context(), reqHeader, p)
	if err != nil {
		blog.Errorf("get user's skip url when no auth, but request to auth center failed, err: %v, rid: %s", err, rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrGetNoAuthSkipURLFailed)})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(url))
}

type ConvertedResource struct {
	Type   string `json:"type"`
	Action string `json:"action"`
}

type ScopeType string

const (
	Business ScopeType = "biz"
	System   ScopeType = "system"
)

type ResourceDetail struct {
	Attribute meta.ResourceAttribute `json:"attribute"`
	Scope     ScopeType              `json:"scope"`
}

type ConvertData struct {
	Data []ResourceDetail `json:"data"`
}

// used for web to get auth's resource with cmdb's resource. in a word, it's for converting.
func (s *service) GetCmdbConvertResources(req *restful.Request, resp *restful.Response) {
	reqHeader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(reqHeader))
	rid := util.GetHTTPCCRequestID(reqHeader)

	attributes := new(ConvertData)
	err := json.NewDecoder(req.Request.Body).Decode(attributes)
	if err != nil {
		blog.Errorf("convert cmdb resource with iam, but decode request failed, err: %v, rid: %s", err, rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	converts := make([]ConvertedResource, 0)
	for _, att := range attributes.Data {
		var bizID int64
		switch att.Scope {
		case Business:
			// set biz id = 1 means that this resource need to convert to a business resource.
			bizID = 1
		case System:
			// set biz id = 1 means that this resource need to convert to a global resource.
			bizID = 0
		default:
			blog.Errorf("convert cmdb resource with iam, but got invalid scope: %s, rid: %s", att.Scope, rid)
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsIsInvalid, att.Scope)})
			return
		}
		typ, err := authcenter.ConvertResourceType(att.Attribute.Type, bizID)
		if err != nil {
			blog.Errorf("convert attribute resource type: %+v failed, err: %v", att, err)
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, att.Attribute.Type)})
			return
		}

		action, err := authcenter.AdaptorAction(&att.Attribute)
		if err != nil {
			blog.Errorf("convert attribute resource action: %+v failed, err: %v", att, err)
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, att.Attribute.Type)})
			return
		}

		converts = append(converts, ConvertedResource{
			Type:   string(*typ),
			Action: string(action),
		})
	}

	resp.WriteEntity(metadata.NewSuccessResp(converts))
}
