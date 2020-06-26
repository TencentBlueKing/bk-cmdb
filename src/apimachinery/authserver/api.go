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

package authserver

import (
	"context"
	"net/http"

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

func (a *authServer) Authorize(ctx context.Context, h http.Header, authAttribute *meta.AuthAttribute) (meta.Decision, error) {
	response := new(struct {
		metadata.BaseResp `json:",inline"`
		Data              meta.Decision `json:"data"`
	})
	subPath := "/authorize"

	err := a.client.Post().
		WithContext(ctx).
		Body(authAttribute).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(response)

	if err != nil {
		return meta.Decision{}, errors.CCHttpError
	}
	if response.Code != 0 {
		return meta.Decision{}, response.CCError()
	}

	return response.Data, nil
}

func (a *authServer) AuthorizeBatch(ctx context.Context, h http.Header, user meta.UserInfo, resources ...meta.ResourceAttribute) ([]meta.Decision, error) {
	input := meta.AuthAttribute{
		User:      user,
		Resources: resources,
	}
	response := new(struct {
		metadata.BaseResp `json:",inline"`
		Data              []meta.Decision `json:"data"`
	})
	subPath := "/authorize/batch"

	err := a.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(response)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if response.Code != 0 {
		return nil, response.CCError()
	}

	return response.Data, nil
}

func (a *authServer) ListAuthorizedResources(ctx context.Context, h http.Header, username string, bizID int64,
	resourceType meta.ResourceType, action meta.Action) ([]iam.IamResource, error) {

	input := meta.ListAuthorizedResourcesParam{
		Username:     username,
		BizID:        bizID,
		ResourceType: resourceType,
		Action:       action,
	}
	response := new(struct {
		metadata.BaseResp `json:",inline"`
		Data              []iam.IamResource `json:"data"`
	})
	subPath := "/findmany/authorized_resource"

	err := a.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(response)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if response.Code != 0 {
		return nil, response.CCError()
	}

	return response.Data, nil
}
