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

package authcenter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/auth/meta"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// clients contains all the client api which is used to
// interact with blueking auth center.

const (
	AuthSupplierAccountHeaderKey = "HTTP_BK_SUPPLIER_ACCOUNT"
)

const (
	codeDuplicated = 1901409
	codeNotFound   = 1901404
)

// Error define
var (
	ErrDuplicated = errors.New("Duplicated item")
	ErrNotFound   = errors.New("Not Found")
)

type authClient struct {
	Config AuthConfig
	// http client instance
	client rest.ClientInterface
	// http header info
	basicHeader http.Header
}

func (a *authClient) verifyExactResourceBatch(ctx context.Context, header http.Header, batch *AuthBatch) ([]BatchStatus, error) {
	util.CopyHeader(a.basicHeader, header)
	resp := new(BatchResult)
	err := a.client.Post().
		SubResourcef("/bkiam/api/v1/perm/systems/%s/resources-perms/batch-verify", a.Config.SystemID).
		WithContext(ctx).
		WithHeaders(header).
		Body(batch).
		Do().Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, &AuthError{
			RequestID: resp.RequestID,
			Reason:    fmt.Errorf("register resource failed, error code: %d, message: %s", resp.Code, resp.Message),
		}
	}

	if len(batch.ResourceActions) != len(resp.Data) {
		return nil, fmt.Errorf("expect %d result, IAM returns %d result", len(batch.ResourceActions), len(resp.Data))
	}

	return resp.Data, nil
}

func (a *authClient) verifyAnyResourceBatch(ctx context.Context, header http.Header, batch *AuthBatch) ([]BatchStatus, error) {
	util.CopyHeader(a.basicHeader, header)
	resp := new(BatchResult)
	err := a.client.Post().
		SubResourcef("/bkiam/api/v1/perm/systems/%s/any-resources-perms/batch-verify", a.Config.SystemID).
		WithContext(ctx).
		WithHeaders(header).
		Body(batch).
		Do().Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, &AuthError{
			RequestID: resp.RequestID,
			Reason:    fmt.Errorf("register resource failed, error code: %d, message: %s", resp.Code, resp.Message),
		}
	}

	if len(batch.ResourceActions) != len(resp.Data) {
		return nil, fmt.Errorf("expect %d result, IAM returns %d result", len(batch.ResourceActions), len(resp.Data))
	}

	return resp.Data, nil
}

func (a *authClient) registerResource(ctx context.Context, header http.Header, info *RegisterInfo) error {
	// register resource with empty id will make crash
	for _, resource := range info.Resources {
		if resource.ResourceID == nil || len(resource.ResourceID) == 0 {
			return fmt.Errorf("resource id can't be empty, resource: %+v", resource)
		}
	}

	util.CopyHeader(a.basicHeader, header)
	resp := new(ResourceResult)
	err := a.client.Post().
		SubResourcef("/bkiam/api/v1/perm/systems/%s/resources/batch-register", a.Config.SystemID).
		WithContext(ctx).
		WithHeaders(header).
		Body(info).
		Do().Into(resp)

	if err != nil {
		return err
	}

	if resp.Code != 0 {
		// 1901409 is for: resource already exist, can not created repeatedly
		if resp.Code == codeDuplicated {
			return ErrDuplicated
		}
		return &AuthError{RequestID: resp.RequestID, Reason: fmt.Errorf("register resource failed, error code: %d, message: %s", resp.Code, resp.Message)}
	}

	if !resp.Data.IsCreated {
		return &AuthError{resp.RequestID, fmt.Errorf("register resource failed, error code: %d", resp.Code)}
	}

	return nil
}

func (a *authClient) deregisterResource(ctx context.Context, header http.Header, info *DeregisterInfo) error {
	util.CopyHeader(a.basicHeader, header)
	resp := new(ResourceResult)
	err := a.client.Delete().
		SubResourcef("/bkiam/api/v1/perm/systems/%s/resources/batch-delete", a.Config.SystemID).
		WithContext(ctx).
		WithHeaders(header).
		Body(info).
		Do().Into(resp)

	if err != nil {
		return err
	}

	// 1901404: resource not exists
	if resp.Code == 1901404 {
		return nil
	}

	if resp.Code != 0 {
		return &AuthError{resp.RequestID, fmt.Errorf("deregister resource failed, error code: %d, message: %s", resp.Code, resp.Message)}
	}

	if !resp.Data.IsDeleted {
		return &AuthError{resp.RequestID, fmt.Errorf("deregister resource failed, error code: %d", resp.Code)}
	}

	return nil
}

func (a *authClient) updateResource(ctx context.Context, header http.Header, info *UpdateInfo) error {
	util.CopyHeader(a.basicHeader, header)
	resp := new(ResourceResult)
	err := a.client.Put().
		SubResourcef("/bkiam/api/v1/perm/systems/%s/resources", a.Config.SystemID).
		WithContext(ctx).
		WithHeaders(header).
		Body(info).
		Do().Into(resp)

	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{resp.RequestID, fmt.Errorf("update resource failed, error code: %d, message: %s", resp.Code, resp.Message)}
	}

	if !resp.Data.IsUpdated {
		return &AuthError{resp.RequestID, fmt.Errorf("update resource failed, error code: %d", resp.Code)}
	}

	return nil
}

func (a *authClient) QuerySystemInfo(ctx context.Context, header http.Header, systemID string, detail bool) (*SystemDetail, error) {
	util.CopyHeader(a.basicHeader, header)

	resp := struct {
		BaseResponse
		Data SystemDetail `json:"data"`
	}{}

	isDetail := "0"
	if detail {
		isDetail = "1"
	}

	err := a.client.Get().
		SubResourcef("/bkiam/api/v1/perm-model/systems/%s", systemID).
		WithParam("is_detail", isDetail).
		WithContext(ctx).
		WithHeaders(header).
		Do().Into(&resp)
	if err != nil {
		return nil, err
	}

	if !resp.Result {
		if resp.Code == codeNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query system info for [%s] failed, err: %v", systemID, resp.ErrorString())
	}

	return &resp.Data, nil
}

func (a *authClient) RegistSystem(ctx context.Context, header http.Header, system System) error {
	util.CopyHeader(a.basicHeader, header)
	const url = "/bkiam/api/v1/perm-model/systems"
	resp := struct {
		BaseResponse
		Data System `json:"data"`
	}{}

	err := a.client.Post().
		SubResource(url).
		WithContext(ctx).
		WithHeaders(header).
		Body(system).
		Do().Into(&resp)
	if err != nil {
		return err
	}

	if !resp.Result {
		if resp.Code == codeDuplicated {
			return ErrDuplicated
		}
		return fmt.Errorf("regist system info for [%s] failed, err: %v", system.SystemID, resp.ErrorString())
	}

	return nil
}

func (a *authClient) UpdateSystem(ctx context.Context, header http.Header, system System) error {
	util.CopyHeader(a.basicHeader, header)
	resp := struct {
		BaseResponse
		Data System `json:"data"`
	}{}

	err := a.client.Put().
		SubResourcef("/bkiam/api/v1/perm-model/systems/%s", system.SystemID).
		WithContext(ctx).
		WithHeaders(header).
		Body(system).
		Do().Into(&resp)
	if err != nil {
		return err
	}

	if !resp.Result {
		return fmt.Errorf("regist system info for [%s] failed, err: %v", system.SystemID, resp.ErrorString())
	}

	return nil
}

func (a *authClient) InitSystemBatch(ctx context.Context, header http.Header, detail SystemDetail) error {
	util.CopyHeader(a.basicHeader, header)
	const url = "/bkiam/api/v1/perm-model/systems/init"
	resp := BaseResponse{}

	err := a.client.Put().
		SubResource(url).
		WithContext(ctx).
		WithHeaders(header).
		Body(detail).
		Do().Into(&resp)
	if err != nil {
		return fmt.Errorf("init system resource failed, error: %v", err)
	}
	if !resp.Result {
		return fmt.Errorf("init system resource failed, err: %v", resp.ErrorString())
	}

	return nil
}

func (a *authClient) RegistResourceTypeBatch(ctx context.Context, header http.Header, systemID, scopeType string, resources []ResourceType) error {
	util.CopyHeader(a.basicHeader, header)
	resp := BaseResponse{}

	err := a.client.Put().
		SubResourcef("/bkiam/api/v1/perm-model/systems/%s/scope-types/%s/resource-types/batch-register", systemID, scopeType).
		WithContext(ctx).
		WithHeaders(header).
		Body(struct {
			ResourceTypes []ResourceType `json:"resource_types"`
		}{resources}).
		Do().Into(&resp)
	if err != nil {
		return fmt.Errorf("regist resource %+v for [%s] failed, error: %v", resources, systemID, err)
	}
	if !resp.Result {
		return fmt.Errorf("regist resource %+v for [%s] failed, err: %v", resources, systemID, resp.ErrorString())
	}

	return nil
}

func (a *authClient) UpdateResourceTypeBatch(ctx context.Context, header http.Header, systemID, scopeType string, resources []ResourceType) error {
	util.CopyHeader(a.basicHeader, header)
	resp := BaseResponse{}

	err := a.client.Put().
		SubResourcef("/bkiam/api/v1/perm-model/systems/%s/scope-types/%s/resource-types/batch-update", systemID, scopeType).
		WithContext(ctx).
		WithHeaders(header).
		Body(struct {
			ResourceTypes []ResourceType `json:"resource_types"`
		}{resources}).
		Do().Into(&resp)
	if err != nil {
		return fmt.Errorf("regist resource %+v for [%s] failed, error: %v", resources, systemID, err)
	}
	if !resp.Result {
		return fmt.Errorf("regist resource %+v for [%s] failed, err: %v", resources, systemID, resp.ErrorString())
	}

	return nil
}

func (a *authClient) UpdateResourceTypeActionBatch(ctx context.Context, header http.Header, systemID, scopeType string, resources []ResourceType) error {
	util.CopyHeader(a.basicHeader, header)
	resp := BaseResponse{}

	err := a.client.Put().
		SubResourcef("/bkiam/api/v1/perm-model/systems/%s/scope-types/%s/resource-type-actions/batch-update", systemID, scopeType).
		WithContext(ctx).
		WithHeaders(header).
		Body(struct {
			ResourceTypes []ResourceType `json:"resource_types"`
		}{resources}).
		Do().Into(&resp)
	if err != nil {
		return fmt.Errorf("regist resource %+v for [%s] failed, error: %v", resources, systemID, err)
	}
	if !resp.Result {
		return fmt.Errorf("regist resource %+v for [%s] failed, err: %v", resources, systemID, resp.ErrorString())
	}

	return nil
}

func (a *authClient) UpsertResourceTypeBatch(ctx context.Context, header http.Header, systemID, scopeType string, resources []ResourceType) error {
	util.CopyHeader(a.basicHeader, header)
	resp := BaseResponse{}

	err := a.client.Post().
		SubResourcef("/bkiam/api/v1/perm-model/systems/%s/scope-types/%s/resource-types/batch-upsert", systemID, scopeType).
		WithContext(ctx).
		WithHeaders(header).
		Body(struct {
			ResourceTypes []ResourceType `json:"resource_types"`
		}{resources}).
		Do().Into(&resp)
	if err != nil {
		return fmt.Errorf("regist resource %+v for [%s] failed, error: %v", resources, systemID, err)
	}
	if !resp.Result {
		return fmt.Errorf("regist resource %+v for [%s] failed, message: %s, code: %v", resources, systemID, resp.Message, resp.Code)
	}

	return nil
}

func (a *authClient) DeleteResourceType(ctx context.Context, header http.Header, systemID, scopeType, resourceType string) error {
	util.CopyHeader(a.basicHeader, header)
	resp := BaseResponse{}

	err := a.client.Delete().
		SubResourcef("/bkiam/api/v1/perm-model/systems/%s/scope-types/%s/resource-types/%s", systemID, scopeType, resourceType).
		WithContext(ctx).
		WithHeaders(header).
		Do().Into(&resp)
	if err != nil {
		return fmt.Errorf("delete resource type %+v for [%s] failed, error: %v", resourceType, systemID, err)
	}
	if !resp.Result {
		return fmt.Errorf("regist resource %+v for [%s] failed, err: %v", resourceType, systemID, resp.ErrorString())
	}

	return nil
}

func (a *authClient) GetAuthorizedResources(ctx context.Context, body *ListAuthorizedResources) ([]AuthorizedResource, error) {
	header := util.CloneHeader(a.basicHeader)
	resp := ListAuthorizedResourcesResult{}

	err := a.client.Post().
		SubResourcef("/bkiam/api/v1/perm/systems/%s/authorized-resources/search", SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(header).
		Body(body).
		Do().Into(&resp)
	if err != nil {
		return nil, fmt.Errorf("get authorized resource failed, err: %v", err)
	}
	if !resp.Result {
		return nil, fmt.Errorf("get authorized resource failed, err: %v", resp.ErrorString())
	}

	return resp.Data, nil
}

// find resource list that a user got any authorized resources.
func (a *authClient) GetAnyAuthorizedScopes(ctx context.Context, scopeID string, body *Principal) ([]string, error) {
	header := util.CloneHeader(a.basicHeader)
	resp := ListAuthorizedScopeResult{}

	err := a.client.Post().
		SubResourcef("/bkiam/api/v1/perm/systems/%s/scope_type/%s/authorized-scopes", SystemIDCMDB, scopeID).
		WithContext(ctx).
		WithHeaders(header).
		Body(body).
		Do().Into(&resp)
	if err != nil {
		return nil, fmt.Errorf("get authorized resource failed, err: %v", err)
	}
	if !resp.Result {
		return nil, fmt.Errorf("get authorized resource failed, err: %v", resp.ErrorString())
	}

	return resp.Data, nil
}

func (a *authClient) ListPageResources(ctx context.Context, header http.Header, searchCondition SearchCondition, limit, offset int64) (result PageBackendResource, err error) {
	util.CopyHeader(a.basicHeader, header)

	resp := new(SearchPageResult)

	err = a.client.Post().
		SubResourcef("/bkiam/api/v1/perm/systems/%s/resources/search", a.Config.SystemID).
		WithContext(ctx).
		WithHeaders(header).
		Body(searchCondition).
		WithParam("page", "1").
		WithParam("limit", strconv.FormatInt(limit, 10)).
		WithParam("offset", strconv.FormatInt(offset, 10)).
		Do().Into(&resp)
	if err != nil {
		return resp.Data, fmt.Errorf("search resource with condition: %+v failed, error: %v", searchCondition, err)
	}
	if !resp.Result || resp.Code != 0 {
		return resp.Data, fmt.Errorf("search resource with condition: %+v failed, message: %s, code: %v", searchCondition, resp.Message, resp.Code)
	}

	return resp.Data, nil
}

func (a *authClient) ListResources(ctx context.Context, header http.Header, searchCondition SearchCondition) (result []meta.BackendResource, err error) {
	util.CopyHeader(a.basicHeader, header)

	resp := new(SearchResult)

	err = a.client.Post().
		SubResourcef("/bkiam/api/v1/perm/systems/%s/resources/search", a.Config.SystemID).
		WithContext(ctx).
		WithHeaders(header).
		Body(searchCondition).
		Do().Into(&resp)
	if err != nil {
		return nil, fmt.Errorf("search resource with condition: %+v failed, error: %v", searchCondition, err)
	}
	if !resp.Result || resp.Code != 0 {
		return nil, fmt.Errorf("search resource with condition: %+v failed, message: %s, code: %v", searchCondition, resp.Message, resp.Code)
	}

	return resp.Data, nil
}

func (a *authClient) RegisterUserRole(ctx context.Context, header http.Header, roles RoleWithAuthResources) (int64, error) {
	util.CopyHeader(a.basicHeader, header)
	resp := new(RegisterRoleResult)

	err := a.client.Post().
		SubResourcef("/bkiam/api/v1/perm-model/systems/%s/perm-templates", a.Config.SystemID).
		WithContext(ctx).
		WithHeaders(header).
		Body(roles).
		Do().Into(resp)
	if err != nil {
		return 0, err
	}
	if !resp.Result || resp.Code != 0 {
		return 0, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp.Data.TemplateID, nil
}

// returns the url which can helps to launch the bk-auth-center when a user do not
// have the authorize to access resource(s).
func (a *authClient) GetNoAuthSkipUrl(ctx context.Context, header http.Header, p []metadata.Permission) (skipUrl string, err error) {
	util.CopyHeader(a.basicHeader, header)
	url := "/o/bk_iam_app/api/v1/apply-permission/url/"
	req := map[string]interface{}{
		"permission": p,
	}
	resp := new(GetSkipUrlResult)
	err = a.client.Post().
		SubResource(url).
		WithContext(ctx).
		WithHeaders(header).
		Body(req).
		Do().Into(&resp)
	if err != nil {
		return "", err
	}
	if !resp.Result || resp.Code != 0 {
		return "", fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp.Data.Url, nil
}

// get user's group members from auth center
func (a *authClient) GetUserGroupMembers(ctx context.Context, header http.Header, bizID int64, groups []string) ([]UserGroupMembers, error) {
	util.CopyHeader(a.basicHeader, header)
	resp := new(UserGroupMembersResult)
	err := a.client.Get().
		SubResourcef("/bkiam/api/v1/perm/systems/%s/scope-types/%s/scopes/%d/group-users", SystemIDCMDB, "biz", bizID).
		WithContext(ctx).
		WithHeaders(header).
		WithParam("group_codes", strings.Join(groups, ",")).
		Do().Into(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.Result || resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp.Data, nil
}

// delete iam resource which has already registered from iam.
// scope type value can be enum of biz or system.
func (a *authClient) DeleteResources(ctx context.Context, header http.Header, scopeType string, resType ResourceTypeID) error {
	util.CopyHeader(a.basicHeader, header)
	resp := new(BaseResponse)
	err := a.client.Delete().
		SubResourcef("/bkiam/api/v1/perm-model/systems/%s/scope-types/%s/resource-types/%s", SystemIDCMDB, scopeType, resType).
		WithContext(ctx).
		WithHeaders(header).
		Do().Into(&resp)
	if err != nil {
		return err
	}
	if !resp.Result || resp.Code != 0 {
		// resource not exist error code
		if resp.Code == 1901002 {
			// delete a not exist resource, so it means success to us.
			return nil
		} else {
			return fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
		}
	}

	return nil
}
