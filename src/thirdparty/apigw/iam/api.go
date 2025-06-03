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
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"configcenter/src/ac/iam/types"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/auth_server/sdk/operator"
	"configcenter/src/thirdparty/apigw/apigwutil"
	"configcenter/src/thirdparty/apigw/apigwutil/user"
)

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

// RegisterSystem register a system in IAM
func (i *iam) RegisterSystem(ctx context.Context, header http.Header, sys System) error {
	resp := new(apigwutil.ApiGWBaseResponse)

	subPath := "/api/v1/model/systems"
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	blog.Errorf("register system, url: %s, body: %+v, config: %+v", subPath, sys, *sys.ProviderConfig)
	result := i.service.Client.Post().
		SubResourcef(subPath).
		WithContext(ctx).
		WithHeaders(h).
		Body(sys).
		Do()

	err = result.Into(resp)
	if err != nil {
		blog.Errorf("err: %s", err)
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

// GetSystemInfo get a system info from IAM, if fields is empty, find all system info
func (i *iam) GetSystemInfo(ctx context.Context, header http.Header, fields []types.SystemQueryField) (
	*RegisteredSystemInfo, error) {

	resp := new(SystemResp)
	fieldsStr := ""
	if len(fields) > 0 {
		fieldArr := make([]string, len(fields))
		for idx, field := range fields {
			fieldArr[idx] = string(field)
		}
		fieldsStr = strings.Join(fieldArr, ",")
	}

	subPath := "/api/v1/model/systems/%s/query"
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return nil, err
	}

	result := i.service.Client.Get().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		WithParam("fields", fieldsStr).
		Body(nil).
		Do()

	err = result.Into(resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		if resp.Code == codeNotFound {
			return nil, ErrNotFound
		}

		return nil, &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("code: %d, msg:%s", resp.Code, resp.Message),
		}
	}

	return &resp.Data, nil
}

// UpdateSystemConfig update system config in IAM, can only update provider_config.host field.
func (i *iam) UpdateSystemConfig(ctx context.Context, header http.Header, config *SysConfig) error {
	sys := new(System)
	config.Auth = "basic"
	sys.ProviderConfig = config
	resp := new(apigwutil.ApiGWBaseResponse)
	subPath := "/api/v1/model/systems/%s"

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	result := i.service.Client.Put().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(sys).
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

// RegisterResourcesTypes register resource types in IAM
func (i *iam) RegisterResourcesTypes(ctx context.Context, header http.Header, resTypes []ResourceType) error {
	if len(resTypes) == 0 {
		return nil
	}

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	subPath := "/api/v1/model/systems/%s/resource-types"
	resp := new(apigwutil.ApiGWBaseResponse)
	result := i.service.Client.Post().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(resTypes).
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

// UpdateResourcesType update resource type in IAM
func (i *iam) UpdateResourcesType(ctx context.Context, header http.Header, resType ResourceType) error {
	resp := new(apigwutil.ApiGWBaseResponse)

	subPath := "/api/v1/model/systems/%s/resource-types/%s"
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	result := i.service.Client.Put().
		SubResourcef(subPath, types.SystemIDCMDB, resType.ID).
		WithContext(ctx).
		WithHeaders(h).
		Body(resType).
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

// DeleteResourcesTypes delete resource types in IAM
func (i *iam) DeleteResourcesTypes(ctx context.Context, header http.Header, resTypeIDs []types.TypeID) error {
	if len(resTypeIDs) == 0 {
		return nil
	}

	ids := make([]struct {
		ID types.TypeID `json:"id"`
	}, len(resTypeIDs))
	for idx := range resTypeIDs {
		ids[idx].ID = resTypeIDs[idx]
	}

	subPath := "/api/v1/model/systems/%s/resource-types"
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	resp := new(apigwutil.ApiGWBaseResponse)
	result := i.service.Client.Delete().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(ids).
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

// RegisterActions register actions in IAM
func (i *iam) RegisterActions(ctx context.Context, header http.Header, actions []ResourceAction) error {
	if len(actions) == 0 {
		return nil
	}

	subPath := "/api/v1/model/systems/%s/actions"
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	resp := new(apigwutil.ApiGWBaseResponse)
	result := i.service.Client.Post().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(actions).
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

// UpdateAction update action in IAM
func (i *iam) UpdateAction(ctx context.Context, header http.Header, action ResourceAction) error {
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	resp := new(apigwutil.ApiGWBaseResponse)
	subPath := "/api/v1/model/systems/%s/actions/%s"
	result := i.service.Client.Put().
		SubResourcef(subPath, types.SystemIDCMDB, action.ID).
		WithContext(ctx).
		WithHeaders(h).
		Body(action).
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

// DeleteActions delete actions in IAM
func (i *iam) DeleteActions(ctx context.Context, header http.Header, actionIDs []types.ActionID) error {

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	ids := make([]struct {
		ID types.ActionID `json:"id"`
	}, len(actionIDs))
	for idx := range actionIDs {
		ids[idx].ID = actionIDs[idx]
	}

	resp := new(apigwutil.ApiGWBaseResponse)
	subPath := "/api/v1/model/systems/%s/actions"
	result := i.service.Client.Delete().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(ids).
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

// RegisterActionGroups register action groups in IAM
func (i *iam) RegisterActionGroups(ctx context.Context, header http.Header, actionGroups []ActionGroup) error {
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	resp := new(apigwutil.ApiGWBaseResponse)
	subPath := "/api/v1/model/systems/%s/configs/action_groups"
	result := i.service.Client.Post().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(actionGroups).
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

// UpdateActionGroups update action groups in IAM
func (i *iam) UpdateActionGroups(ctx context.Context, header http.Header, actionGroups []ActionGroup) error {
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	resp := new(apigwutil.ApiGWBaseResponse)
	subPath := "/api/v1/model/systems/%s/configs/action_groups"
	result := i.service.Client.Put().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(actionGroups).
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

// RegisterInstanceSelections register instance selections in IAM
func (i *iam) RegisterInstanceSelections(ctx context.Context, header http.Header,
	instanceSelections []InstanceSelection) error {

	if len(instanceSelections) == 0 {
		return nil
	}

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	subPath := "/api/v1/model/systems/%s/instance-selections"

	resp := new(apigwutil.ApiGWBaseResponse)
	result := i.service.Client.Post().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(instanceSelections).
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

// UpdateInstanceSelection update instance selection in IAM
func (i *iam) UpdateInstanceSelection(ctx context.Context, header http.Header,
	instanceSelection InstanceSelection) error {

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}
	subPath := "/api/v1/model/systems/%s/instance-selections/%s"

	resp := new(apigwutil.ApiGWBaseResponse)
	result := i.service.Client.Put().
		SubResourcef(subPath, types.SystemIDCMDB, instanceSelection.ID).
		WithContext(ctx).
		WithHeaders(h).
		Body(instanceSelection).
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

// DeleteInstanceSelections delete instance selections in IAM
func (i *iam) DeleteInstanceSelections(ctx context.Context, header http.Header,
	instanceSelectionIDs []types.InstanceSelectionID) error {

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}
	subPath := "/api/v1/model/systems/%s/instance-selections"

	if len(instanceSelectionIDs) == 0 {
		return nil
	}

	ids := make([]struct {
		ID types.InstanceSelectionID `json:"id"`
	}, len(instanceSelectionIDs))
	for idx := range instanceSelectionIDs {
		ids[idx].ID = instanceSelectionIDs[idx]
	}

	resp := new(apigwutil.ApiGWBaseResponse)
	result := i.service.Client.Delete().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(ids).
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

// RegisterResourceCreatorActions register resource creator actions in IAM
func (i *iam) RegisterResourceCreatorActions(ctx context.Context, header http.Header,
	resourceCreatorActions ResourceCreatorActions) error {

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}
	subPath := "/api/v1/model/systems/%s/configs/resource_creator_actions"

	resp := new(apigwutil.ApiGWBaseResponse)
	result := i.service.Client.Post().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(resourceCreatorActions).
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

// UpdateResourceCreatorActions update resource creator actions in IAM
func (i *iam) UpdateResourceCreatorActions(ctx context.Context, header http.Header,
	resourceCreatorActions ResourceCreatorActions) error {

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}
	subPath := "/api/v1/model/systems/%s/configs/resource_creator_actions"

	resp := new(apigwutil.ApiGWBaseResponse)
	result := i.service.Client.Put().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(resourceCreatorActions).
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

// RegisterCommonActions register common actions in IAM
func (i *iam) RegisterCommonActions(ctx context.Context, header http.Header, commonActions []CommonAction) error {

	resp := new(apigwutil.ApiGWBaseResponse)
	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}
	subPath := "/api/v1/model/systems/%s/configs/common_actions"

	result := i.service.Client.Post().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(commonActions).
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

// UpdateCommonActions update common actions in IAM
func (i *iam) UpdateCommonActions(ctx context.Context, header http.Header, commonActions []CommonAction) error {

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return err
	}

	resp := new(apigwutil.ApiGWBaseResponse)
	subPath := "/api/v1/model/systems/%s/configs/common_actions"
	result := i.service.Client.Put().
		SubResourcef(subPath, types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(commonActions).
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
	resp := new(struct {
		apigwutil.ApiGWBaseResponse
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	})

	h, err := user.SetBKAuthHeader(ctx, i.service.Config, header, i.userCli)
	if err != nil {
		return "", err
	}

	result := i.service.Client.Get().
		SubResourcef("/api/v1/model/systems/%s/token", types.SystemIDCMDB).
		WithContext(ctx).
		WithHeaders(h).
		Body(nil).
		Do()

	err = result.Into(resp)
	if err != nil {
		return "", err
	}

	if resp.Code != 0 {
		return "", &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("code: %d, msg:%s", resp.Code, resp.Message),
		}
	}

	return resp.Data.Token, nil
}
