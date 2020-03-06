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
	"fmt"
	"net/http"

	"configcenter/src/apimachinery/rest"
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
			RequestID: result.Header.Get(iamRequestHeader),
			Reason:    fmt.Errorf("register system failed, code: %d, msg:%s", resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) GetSystemInfo(ctx context.Context) (*SystemResp, error) {
	resp := new(SystemResp)
	result := c.client.Post().
		SubResourcef("/api/v1/model/systems/%s/query", c.Config.SystemID).
		WithContext(ctx).
		WithHeaders(c.basicHeader).
		WithParam("fields", "base_info").
		WithParam("fields", "resource_types").
		WithParam("fields", "actions").
		WithParam("fields", "action_topology").
		Body(nil).Do()
	err := result.Into(resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, &AuthError{
			RequestID: result.Header.Get(iamRequestHeader),
			Reason:    fmt.Errorf("get system info failed, code: %d, msg:%s", resp.Code, resp.Message),
		}
	}

	return resp, nil
}

// Note: can only update provider_config.host field.
func (c *iamClient) UpdateSystemConfig(ctx context.Context, config *SysConfig) error {
	sys := new(System)
	config.Auth = ""
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
			RequestID: result.Header.Get(iamRequestHeader),
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
			RequestID: result.Header.Get(iamRequestHeader),
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
			RequestID: result.Header.Get(iamRequestHeader),
			Reason:    fmt.Errorf("udpate resource type %s failed, code: %d, msg:%s", resType.ID, resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) DeleteResourcesTypes(ctx context.Context, resTypeIDs []ResourceTypeID) error {

	ids := make([]struct {
		ID ResourceTypeID `json:"id"`
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
			RequestID: result.Header.Get(iamRequestHeader),
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
			RequestID: result.Header.Get(iamRequestHeader),
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
			RequestID: result.Header.Get(iamRequestHeader),
			Reason:    fmt.Errorf("udpate resource action %v failed, code: %d, msg:%s", actions, resp.Code, resp.Message),
		}
	}

	return nil
}

func (c *iamClient) DeleteAction(ctx context.Context, actionIDs []ResourceActionID) error {
	ids := make([]struct {
		ID ResourceActionID `json:"id"`
	}, len(actionIDs))
	for idx := range actionIDs {
		ids[idx].ID = actionIDs[idx]
	}

	resp := new(BaseResponse)
	result := c.client.Put().
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
			RequestID: result.Header.Get(iamRequestHeader),
			Reason:    fmt.Errorf("delete resource actions %v failed, code: %d, msg:%s", actions, resp.Code, resp.Message),
		}
	}

	return nil
}
