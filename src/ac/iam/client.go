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

package iam

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"configcenter/src/apimachinery/rest"
)

const (
	codeNotFound = 1901404
)

var (
	ErrNotFound = errors.New("Not Found")
)

type iamClient struct {
	Config AuthConfig
	// http client instance
	client rest.ClientInterface
	// http header info
	basicHeader http.Header
}

func (c *iamClient) RegisterSystem(ctx context.Context, sys System) error {
	resp := new(BaseResponse)
	result := c.client.Post().
		SubResourcef("/api/v1/model/systems").
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(sys).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("register system failed, code: %d, msg:%s", resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) GetSystemInfo(ctx context.Context) (*SystemResp, error) {
	resp := new(SystemResp)
	result := c.client.Get().
		SubResourcef("/api/v1/model/systems/%s/query", c.Config.SystemID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		WithParam("fields", "base_info,resource_types,actions,action_groups,instance_selections,resource_creator_actions").
		Body(nil).Do()
	err := result.Into(resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		if resp.Code == codeNotFound {
			return resp, ErrNotFound
		}
		return nil, &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("get system info failed, code: %d, msg:%s", resp.Code, resp.Message),
		}
	}

	return resp, nil
}

// Note: can only update provider_config.host field.
func (c *iamClient) UpdateSystemConfig(ctx context.Context, config *SysConfig) error {
	sys := new(System)
	config.Auth = "basic"
	sys.ProviderConfig = config
	resp := new(BaseResponse)
	result := c.client.Put().
		SubResourcef("/api/v1/model/systems/%s", c.Config.SystemID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(sys).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("update system config failed, code: %d, msg:%s", resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) RegisterResourcesTypes(ctx context.Context, resTypes []ResourceType) error {
	resp := new(BaseResponse)
	result := c.client.Post().
		SubResourcef("/api/v1/model/systems/%s/resource-types", c.Config.SystemID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(resTypes).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("register system failed, code: %d, msg:%s", resp.Code, resp.Message),
		}
	}

	return nil

}

func (c *iamClient) UpdateResourcesTypes(ctx context.Context, resType ResourceType) error {
	resp := new(BaseResponse)
	result := c.client.Put().
		SubResourcef("/api/v1/model/systems/%s/resource-types/%s", c.Config.SystemID, resType.ID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(resType).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("udpate resource type %s failed, code: %d, msg:%s", resType.ID, resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) DeleteResourcesTypes(ctx context.Context, resTypeIDs []TypeID) error {

	ids := make([]struct {
		ID TypeID `json:"id"`
	}, len(resTypeIDs))
	for idx := range resTypeIDs {
		ids[idx].ID = resTypeIDs[idx]
	}

	resp := new(BaseResponse)
	result := c.client.Delete().
		SubResourcef("/api/v1/model/systems/%s/resource-types", c.Config.SystemID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(ids).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("delete resource type %v failed, code: %d, msg:%s", resTypeIDs, resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) CreateAction(ctx context.Context, actions []ResourceAction) error {

	resp := new(BaseResponse)
	result := c.client.Post().
		SubResourcef("/api/v1/model/systems/%s/actions", c.Config.SystemID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(actions).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("add resource actions %v failed, code: %d, msg:%s", actions, resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) UpdateAction(ctx context.Context, action ResourceAction) error {

	resp := new(BaseResponse)
	result := c.client.Put().
		SubResourcef("/api/v1/model/systems/%s/actions/%s", c.Config.SystemID, action.ID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(action).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("udpate resource action %v failed, code: %d, msg:%s", action, resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) DeleteAction(ctx context.Context, actionIDs []ActionID) error {
	ids := make([]struct {
		ID ActionID `json:"id"`
	}, len(actionIDs))
	for idx := range actionIDs {
		ids[idx].ID = actionIDs[idx]
	}

	resp := new(BaseResponse)
	result := c.client.Delete().
		SubResourcef("/api/v1/model/systems/%s/actions", c.Config.SystemID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(ids).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("delete resource actions %v failed, code: %d, msg:%s", actionIDs, resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) RegisterActionGroups(ctx context.Context, actionGroups []ActionGroup) error {

	resp := new(BaseResponse)
	result := c.client.Post().
		SubResourcef("/api/v1/model/systems/%s/configs/action_groups", c.Config.SystemID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(actionGroups).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("register action groups %v failed, code: %d, msg:%s", actionGroups, resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) UpdateActionGroups(ctx context.Context, actionGroups []ActionGroup) error {

	resp := new(BaseResponse)
	result := c.client.Put().
		SubResourcef("/api/v1/model/systems/%s/configs/action_groups", c.Config.SystemID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(actionGroups).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("update action groups %v failed, code: %d, msg:%s", actionGroups, resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) CreateInstanceSelection(ctx context.Context, instanceSelections []InstanceSelection) error {

	resp := new(BaseResponse)
	result := c.client.Post().
		SubResourcef("/api/v1/model/systems/%s/instance-selections", c.Config.SystemID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(instanceSelections).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("add instance selections %v failed, code: %d, msg:%s", instanceSelections, resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) UpdateInstanceSelection(ctx context.Context, instanceSelection InstanceSelection) error {

	resp := new(BaseResponse)
	result := c.client.Put().
		SubResourcef("/api/v1/model/systems/%s/instance-selections/%s", c.Config.SystemID, instanceSelection.ID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(instanceSelection).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("udpate instance selections %v failed, code: %d, msg:%s", instanceSelection, resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) DeleteInstanceSelection(ctx context.Context, instanceSelectionIDs []InstanceSelectionID) error {
	ids := make([]struct {
		ID InstanceSelectionID `json:"id"`
	}, len(instanceSelectionIDs))
	for idx := range instanceSelectionIDs {
		ids[idx].ID = instanceSelectionIDs[idx]
	}

	resp := new(BaseResponse)
	result := c.client.Delete().
		SubResourcef("/api/v1/model/systems/%s/instance-selections", c.Config.SystemID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(ids).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("delete instance selections %v failed, code: %d, msg:%s", instanceSelectionIDs, resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) RegisterResourceCreatorActions(ctx context.Context, resourceCreatorActions ResourceCreatorActions) error {

	resp := new(BaseResponse)
	result := c.client.Post().
		SubResourcef("/api/v1/model/systems/%s/configs/resource_creator_actions", c.Config.SystemID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(resourceCreatorActions).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("register resource creator actions %v failed, code: %d, msg:%s", resourceCreatorActions, resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) UpdateResourceCreatorActions(ctx context.Context, resourceCreatorActions ResourceCreatorActions) error {

	resp := new(BaseResponse)
	result := c.client.Put().
		SubResourcef("/api/v1/model/systems/%s/configs/resource_creator_actions", c.Config.SystemID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		Body(resourceCreatorActions).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{
			RequestID: result.Header.Get(IamRequestHeader),
			Reason:    fmt.Errorf("update resource creator actions %v failed, code: %d, msg:%s", resourceCreatorActions, resp.Code, resp.Message),
		}
	}

	return nil
}
