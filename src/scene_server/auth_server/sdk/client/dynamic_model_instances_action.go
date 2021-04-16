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
	"context"
	"fmt"

	"configcenter/src/ac/iam"
)

func (ac *authClient) UpdateAction(ctx context.Context, action iam.ResourceAction) error {

	resp := new(iam.BaseResponse)
	result := ac.client.Put().
		SubResourcef("/api/v1/model/systems/%s/actions/%s", ac.config.SystemID, action.ID).
		WithContext(ctx).
		WithHeaders(ac.basicHeader).
		Body(action).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &iam.AuthError{
			RequestID: result.Header.Get(iam.IamRequestHeader),
			Reason:    fmt.Errorf("udpate resource action %v failed, code: %d, msg:%s", action, resp.Code, resp.Message),
		}
	}

	return nil
}

func (ac *authClient) CreateActions(ctx context.Context, actions []iam.ResourceAction) error {

	resp := new(iam.BaseResponse)
	result := ac.client.Post().
		SubResourcef("/api/v1/model/systems/%s/actions", ac.config.SystemID).
		WithContext(ctx).
		WithHeaders(ac.basicHeader).
		Body(actions).Do()
	err := result.Into(resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &iam.AuthError{
			RequestID: result.Header.Get(iam.IamRequestHeader),
			Reason:    fmt.Errorf("add resource actions %v failed, code: %d, msg:%s", actions, resp.Code, resp.Message),
		}
	}

	return nil
}
