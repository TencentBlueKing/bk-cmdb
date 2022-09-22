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

package client

import (
	"configcenter/cmd/scene_server/auth_server/sdk/operator"
	types2 "configcenter/cmd/scene_server/auth_server/sdk/types"
	"context"
)

// GetUserPolicy get a user's policy with a action and resources
func (ac *authClient) GetUserPolicy(ctx context.Context, opt *types2.GetPolicyOption) (*operator.Policy, error) {
	resp := new(types2.GetPolicyResp)

	// iam requires resources to be set
	if opt.Resources == nil {
		opt.Resources = make([]types2.Resource, 0)
	}

	result := ac.client.Post().
		SubResourcef("/api/v1/policy/query").
		WithContext(ctx).
		WithHeaders(ac.cloneHeader(ctx)).
		Body(opt).
		Do()

	err := result.Into(resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, &types2.AuthError{
			Rid:     result.Header.Get(types2.RequestIDHeaderKey),
			Code:    resp.Code,
			Message: resp.Message,
		}
	}

	return resp.Data, nil
}

// ListUserPolicies get a user's policy with multiple actions and resources
func (ac *authClient) ListUserPolicies(ctx context.Context, opts *types2.ListPolicyOptions) (
	[]*types2.ActionPolicy, error) {

	resp := new(types2.ListPolicyResp)

	// iam requires resources to be set
	if opts.Resources == nil {
		opts.Resources = make([]types2.Resource, 0)
	}

	result := ac.client.Post().
		SubResourcef("/api/v1/policy/query_by_actions").
		WithContext(ctx).
		WithHeaders(ac.cloneHeader(ctx)).
		Body(opts).
		Do()

	err := result.Into(resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, &types2.AuthError{
			Rid:     result.Header.Get(types2.RequestIDHeaderKey),
			Code:    resp.Code,
			Message: resp.Message,
		}
	}

	return resp.Data, nil
}

// GetSystemToken get system token from iam, used to validate if request is from iam
func (ac *authClient) GetSystemToken(ctx context.Context) (string, error) {
	resp := new(struct {
		types2.BaseResp
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	})
	result := ac.client.Get().
		SubResourcef("/api/v1/model/systems/%s/token", ac.config.SystemID).
		WithContext(ctx).
		WithHeaders(ac.basicHeader).
		Body(nil).Do()
	err := result.Into(resp)
	if err != nil {
		return "", err
	}

	if resp.Code != 0 {
		return "", &types2.AuthError{
			Rid:     result.Header.Get(types2.RequestIDHeaderKey),
			Code:    resp.Code,
			Message: resp.Message,
		}
	}

	return resp.Data.Token, nil
}
