/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"configcenter/src/ac/iam/types"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/auth_server/sdk/operator"
	"configcenter/src/thirdparty/apigw/apigwutil"
	"configcenter/src/thirdparty/apigw/apigwutil/user"
)

// handleIamResp handle IAM V4 response.
func handleIamResp[T any](result *rest.Result) (T, error) {
	var data T
	if result.Err != nil {
		return data, result.Err
	}

	if result.StatusCode >= http.StatusOK && result.StatusCode < http.StatusMultipleChoices {
		if len(result.Body) == 0 {
			return data, nil
		}

		resp := new(struct {
			Data T `json:"data"`
		})
		if err := json.Unmarshal(result.Body, resp); err != nil {
			return data, &AuthError{
				RequestID:  result.Header.Get(IamRequestHeader),
				StatusCode: result.StatusCode,
				Reason:     fmt.Errorf("unmarshal iam response failed: %v, body: %s", err, result.Body),
			}
		}
		return resp.Data, nil
	}

	authErr := &AuthError{
		RequestID:  result.Header.Get(IamRequestHeader),
		StatusCode: result.StatusCode,
		Reason:     fmt.Errorf("body: %s", result.Body),
	}

	if len(result.Body) == 0 {
		return data, authErr
	}

	errResp := new(IamErrorResp)
	if err := json.Unmarshal(result.Body, errResp); err == nil && errResp.Error != nil {
		authErr.Reason = fmt.Errorf("code: %s, msg: %s", errResp.Error.Code, errResp.Error.Message)
	}
	return data, authErr
}

// GetNoAuthSkipUrl returns the url which can helps to launch the bk-iam when user do not have the authority to
// access resource(s).
func (i *iam) GetNoAuthSkipUrl(ctx context.Context, header http.Header, p metadata.IamPermission) (string, error) {
	resp := new(iamPermissionURLResp)
	subPath := "/api/v1/open/application/"

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return "", err
	}

	params := &apiGWIamPermissionParams{
		IamPermission: p,
	}
	err = i.service.Client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return "", err
	}

	if resp.Code != 0 {
		return "", fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp.Data.Url, nil
}

// RegisterResourceCreatorAction register iam resource instance with creator, returns related actions with policy id
// that the creator gained
func (i *iam) RegisterResourceCreatorAction(ctx context.Context, header http.Header,
	instance metadata.IamInstanceWithCreator) ([]metadata.IamCreatorActionPolicy, error) {

	resp := new(iamCreatorActionResp)
	subPath := "/api/v1/open/authorization/resource_creator_action/"
	params := &iamInstanceParams{
		IamInstanceWithCreator: instance,
	}

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}
	err = i.service.Client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp.Data, nil
}

// BatchRegisterResourceCreatorAction batch register iam resource instances with creator, returns related actions with
// policy id that the creator gained
func (i *iam) BatchRegisterResourceCreatorAction(ctx context.Context, header http.Header,
	instances metadata.IamInstancesWithCreator) ([]metadata.IamCreatorActionPolicy, error) {

	resp := new(iamCreatorActionResp)
	url := "/api/v1/open/authorization/batch_resource_creator_action/"
	params := &iamInstancesParams{
		IamInstancesWithCreator: instances,
	}

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}
	err = i.service.Client.Post().
		SubResourcef(url).
		WithContext(ctx).
		WithHeaders(h).
		Body(params).
		Do().
		Into(&resp)

	if err != nil {
		return nil, err
	}

	if !resp.Result || resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp.Data, nil
}

// BatchOperateInstanceAuth batch grant or revoke iam resource instances' authorization
func (i *iam) BatchOperateInstanceAuth(ctx context.Context, header http.Header,
	req *metadata.IamBatchOperateInstanceAuthReq) ([]metadata.IamBatchOperateInstanceAuthRes, error) {

	resp := new(iamBatchOperateInstanceAuthResp)
	url := "/api/v1/open/authorization/batch_instance/"
	params := &iamBatchOperateInstanceAuthParams{
		IamBatchOperateInstanceAuthReq: req,
	}

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}
	err = i.service.Client.Post().
		SubResourcef(url).
		WithContext(ctx).
		WithHeaders(h).
		Body(params).
		Do().
		Into(&resp)

	if err != nil {
		return nil, err
	}

	if !resp.Result || resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp.Data, nil
}

// RegisterSystem register a system in IAM, returns the registered system id
func (i *iam) RegisterSystem(ctx context.Context, header http.Header, sys *System) (string, error) {
	subPath := "/api/v1/open/rbac/model/systems/"
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return "", err
	}

	data, err := handleIamResp[RegisterSystemData](i.service.Client.Post().
		SubResourcef(subPath).
		WithContext(ctx).
		WithHeaders(h).
		Body(sys).
		Do())
	if err != nil {
		return "", err
	}
	return data.ID, nil
}

// GetSystem get system info from IAM
func (i *iam) GetSystem(ctx context.Context, header http.Header) (*System, error) {
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/"
	return handleIamResp[*System](i.service.Client.Get().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(nil).
		Do())
}

// UpdateSystem update system info in IAM, the system id can not be updated
func (i *iam) UpdateSystem(ctx context.Context, header http.Header, sys *System) error {
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/"
	_, err = handleIamResp[struct{}](i.service.Client.Put().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(sys).
		Do())

	return err
}

// ListResourceTypes list resource types by page, the max page size is MaxListPageSize
func (i *iam) ListResourceTypes(ctx context.Context, header http.Header, page, pageSize int64) (
	*ListResourceTypesData, error) {

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/resource-types/"
	return handleIamResp[*ListResourceTypesData](i.service.Client.Get().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		WithParam(pageParam, strconv.FormatInt(page, 10)).
		WithParam(pageSizeParam, strconv.FormatInt(pageSize, 10)).
		Body(nil).
		Do())
}

// RegisterResourcesTypes register resource types in IAM, returns the registered resource type ids
func (i *iam) RegisterResourcesTypes(ctx context.Context, header http.Header, resTypes []ResourceType) (
	[]string, error) {

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/resource-types/"
	return handleIamResp[[]string](i.service.Client.Post().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(resTypes).
		Do())
}

// UpdateResourcesType update resource type in IAM, only name and ancestors can be updated
func (i *iam) UpdateResourcesType(ctx context.Context, header http.Header, resTypeID types.TypeID,
	req *UpdateResourceTypeReq) error {

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/resource-types/%s/"
	_, err = handleIamResp[struct{}](i.service.Client.Put().
		SubResourcef(subPath, types.SystemIDCMDB, resTypeID).
		WithContext(ctx).
		WithHeaders(h).
		Body(req).
		Do())

	return err
}

// DeleteResourcesType delete resource type in IAM
func (i *iam) DeleteResourcesType(ctx context.Context, header http.Header, resTypeID types.TypeID) error {
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/resource-types/%s/"
	_, err = handleIamResp[struct{}](i.service.Client.Delete().
		SubResourcef(subPath, types.SystemIDCMDB, resTypeID).
		WithContext(ctx).
		WithHeaders(h).
		Do())

	return err
}

// ListActions list actions by page, the max page size is MaxListPageSize
func (i *iam) ListActions(ctx context.Context, header http.Header, page, pageSize int64) (*ListActionsData, error) {
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/actions/"
	return handleIamResp[*ListActionsData](i.service.Client.Get().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		WithParam(pageParam, strconv.FormatInt(page, 10)).
		WithParam(pageSizeParam, strconv.FormatInt(pageSize, 10)).
		Body(nil).
		Do())
}

// RegisterActions register actions in IAM, returns the registered action ids
func (i *iam) RegisterActions(ctx context.Context, header http.Header, actions []ResourceAction) ([]string, error) {
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/actions/"
	return handleIamResp[[]string](i.service.Client.Post().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(actions).
		Do())
}

// UpdateAction update action in IAM, only the action name can be updated, changing the related resource type
// needs to delete the action and register it again
func (i *iam) UpdateAction(ctx context.Context, header http.Header, actionID types.ActionID,
	req *UpdateActionReq) error {

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/actions/%s/"
	_, err = handleIamResp[struct{}](i.service.Client.Put().
		SubResourcef(subPath, types.SystemIDCMDB, actionID).
		WithContext(ctx).
		WithHeaders(h).
		Body(req).
		Do())

	return err
}

// DeleteAction delete action in IAM
func (i *iam) DeleteAction(ctx context.Context, header http.Header, actionID types.ActionID) error {
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/actions/%s/"
	_, err = handleIamResp[struct{}](i.service.Client.Delete().
		SubResourcef(subPath, types.SystemIDCMDB, actionID).
		WithContext(ctx).
		WithHeaders(h).
		Do())

	return err
}

// ListRoles list roles by page, the max page size is MaxListPageSize
func (i *iam) ListRoles(ctx context.Context, header http.Header, page, pageSize int64) (*ListRolesData, error) {
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/roles/"
	return handleIamResp[*ListRolesData](i.service.Client.Get().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		WithParam(pageParam, strconv.FormatInt(page, 10)).
		WithParam(pageSizeParam, strconv.FormatInt(pageSize, 10)).
		Body(nil).
		Do())
}

// RegisterRoles register roles in IAM, returns the registered role ids
func (i *iam) RegisterRoles(ctx context.Context, header http.Header, roles []Role) ([]string, error) {
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/roles/"
	return handleIamResp[[]string](i.service.Client.Post().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(roles).
		Do())
}

// UpdateRole update role in IAM, only name and description can be updated, the bound actions are updated by
// AddRoleActions and DeleteRoleActions
func (i *iam) UpdateRole(ctx context.Context, header http.Header, roleID types.RoleID, req *UpdateRoleReq) error {
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/roles/%s/"
	_, err = handleIamResp[struct{}](i.service.Client.Put().
		SubResourcef(subPath, types.SystemIDCMDB, roleID).
		WithContext(ctx).
		WithHeaders(h).
		Body(req).
		Do())

	return err
}

// DeleteRole delete role in IAM
func (i *iam) DeleteRole(ctx context.Context, header http.Header, roleID types.RoleID) error {
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/roles/%s/"
	_, err = handleIamResp[struct{}](i.service.Client.Delete().
		SubResourcef(subPath, types.SystemIDCMDB, roleID).
		WithContext(ctx).
		WithHeaders(h).
		Do())

	return err
}

// AddRoleActions add actions to a role, returns the added action ids
func (i *iam) AddRoleActions(ctx context.Context, header http.Header, roleID types.RoleID, actions []RoleAction) (
	[]string, error) {

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/roles/%s/actions/"
	return handleIamResp[[]string](i.service.Client.Post().
		SubResourcef(subPath, types.SystemIDCMDB, roleID).
		WithContext(ctx).
		WithHeaders(h).
		Body(actions).
		Do())
}

// DeleteRoleActions delete actions from a role, the action ids are passed by the "ids" query parameter
func (i *iam) DeleteRoleActions(ctx context.Context, header http.Header, roleID types.RoleID,
	actionIDs []types.ActionID) error {

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	ids := make([]string, len(actionIDs))
	for idx, id := range actionIDs {
		ids[idx] = string(id)
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/roles/%s/actions/"
	_, err = handleIamResp[struct{}](i.service.Client.Delete().
		SubResourcef(subPath, types.SystemIDCMDB, roleID).
		WithContext(ctx).
		WithHeaders(h).
		WithParam(idsParam, strings.Join(ids, ",")).
		Body(nil).
		Do())

	return err
}

// DeleteActionPolicies delete action policies in IAM
func (i *iam) DeleteActionPolicies(ctx context.Context, header http.Header, actionID types.ActionID) error {

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}
	resp := new(apigwutil.ApiGWBaseResponse)
	subPath := "/api/v1/model/systems/%s/actions/%s/policies"

	result := i.service.Client.Delete().
		SubResourcef(subPath, types.SystemIDCMDB, actionID).
		WithContext(ctx).
		WithHeaders(h).
		Do()

	err = result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("code: %d, msg:%s", resp.Code, resp.Message),
		}
	}

	return nil
}

// ListPolicies list iam policies
func (i *iam) ListPolicies(ctx context.Context, header http.Header, params *ListPoliciesParams) (*ListPoliciesData,
	error) {

	parsedParams := map[string]string{"action_id": string(params.ActionID)}
	if params.Page != 0 {
		parsedParams["page"] = strconv.FormatInt(params.Page, 10)
	}
	if params.PageSize != 0 {
		parsedParams["page_size"] = strconv.FormatInt(params.PageSize, 10)
	}
	if params.Timestamp != 0 {
		parsedParams["timestamp"] = strconv.FormatInt(params.Timestamp, 10)
	}

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}
	subPath := "/api/v1/open/systems/%s/policies"

	resp := new(ListPoliciesResp)
	result := i.service.Client.Get().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		WithParams(parsedParams).
		Body(nil).
		Do()

	err = result.Into(resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("code: %d, msg:%s", resp.Code, resp.Message),
		}
	}
	return resp.Data, nil
}

// GetUserPolicy get a user's policy with a action and resources
func (i *iam) GetUserPolicy(ctx context.Context, header http.Header, opt *GetPolicyOption) (*operator.Policy, error) {
	resp := new(GetPolicyResp)

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}
	subPath := "/api/v1/policy/query"

	// iam requires resources to be set
	if opt.Resources == nil {
		opt.Resources = make([]Resource, 0)
	}

	result := i.service.Client.Post().
		SubResourcef(subPath).
		WithContext(ctx).
		WithHeaders(h).
		Body(opt).
		Do()

	err = result.Into(resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("code: %d, msg:%s", resp.Code, resp.Message),
		}
	}

	return resp.Data, nil
}

// ListUserPolicies get a user's policy with multiple actions and resources
func (i *iam) ListUserPolicies(ctx context.Context, header http.Header, opts *ListPolicyOptions) (
	[]*ActionPolicy, error) {

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}

	resp := new(ListPolicyResp)
	// iam requires resources to be set
	if opts.Resources == nil {
		opts.Resources = make([]Resource, 0)
	}

	result := i.service.Client.Post().
		SubResourcef("/api/v1/policy/query_by_actions").
		WithContext(ctx).
		WithHeaders(h).
		Body(opts).
		Do()

	err = result.Into(resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("code: %d, msg:%s", resp.Code, resp.Message),
		}
	}
	return resp.Data, nil
}

// GetSystemToken get system token from iam, used to validate if request is from iam
func (i *iam) GetSystemToken(ctx context.Context, header http.Header) (string, error) {
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return "", err
	}

	subPath := "/api/v1/open/rbac/model/systems/%s/auth-token/"
	data, err := handleIamResp[systemAuthToken](i.service.Client.Get().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(nil).
		Do())
	if err != nil {
		return "", err
	}
	return data.AuthToken, nil
}
