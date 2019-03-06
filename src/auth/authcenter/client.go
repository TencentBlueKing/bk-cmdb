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
	"fmt"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/auth/meta"
)

// clients contains all the client api which is used to
// interact with blueking auth center.

const (
	AuthSupplierAccountHeaderKey = "HTTP_BK_SUPPLIER_ACCOUNT"
)

type authClient struct {
	Config AuthConfig
	// http client instance
	client rest.ClientInterface
	// http header info
	basicHeader http.Header
}

func (a *authClient) verifyInList(ctx context.Context, header http.Header, batch *AuthBatch) (meta.Decision, error) {
	for k, v := range a.basicHeader {
		header[k] = v
	}
	resp := new(BatchResult)
	url := fmt.Sprintf("/bkiam/api/v1/perm/systems/%s/resources-perms/verify", a.Config.SystemID)
	err := a.client.Post().
		SubResource(url).
		WithContext(ctx).
		WithHeaders(header).
		Body(batch).
		Do().Into(resp)

	if err != nil {
		return meta.Decision{}, err
	}

	if resp.Code != 0 {
		return meta.Decision{}, &AuthError{
			RequestID: resp.RequestID,
			Reason:    fmt.Errorf("register resource failed, error code: %d, message: %s", resp.Code, resp.ErrMsg),
		}
	}

	noAuth := make([]ResourceType, 0)
	for _, item := range resp.Data {
		if !item.IsPass {
			noAuth = append(noAuth, item.ResourceType)
		}
	}

	if len(noAuth) != 0 {
		return meta.Decision{
			Authorized: false,
			Reason:     fmt.Sprintf("resource [%v] do not have permission", noAuth),
		}, nil
	}

	return meta.Decision{Authorized: true}, nil
}

func (a *authClient) registerResource(ctx context.Context, header http.Header, info *RegisterInfo) error {
	for k, v := range a.basicHeader {
		header[k] = v
	}
	resp := new(ResourceResult)
	url := fmt.Sprintf("/bkiam/api/v1/perm/systems/%s/resources", a.Config.SystemID)
	err := a.client.Post().
		SubResource(url).
		WithContext(ctx).
		WithHeaders(header).
		Body(info).
		Do().Into(resp)

	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{RequestID: resp.RequestID, Reason: fmt.Errorf("register resource failed, error code: %d, message: %s", resp.Code, resp.ErrMsg)}
	}

	if !resp.Data.IsCreated {
		return &AuthError{resp.RequestID, fmt.Errorf("register resource failed, error code: %d", resp.Code)}
	}

	return nil
}

func (a *authClient) deregisterResource(ctx context.Context, header http.Header, info *DeregisterInfo) error {
	for k, v := range a.basicHeader {
		header[k] = v
	}
	resp := new(ResourceResult)
	url := fmt.Sprintf("/bkiam/api/v1/perm/systems/%s/resources", a.Config.SystemID)
	err := a.client.Delete().
		SubResource(url).
		WithContext(ctx).
		WithHeaders(header).
		Body(info).
		Do().Into(resp)

	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{resp.RequestID, fmt.Errorf("deregister resource failed, error code: %d, message: %s", resp.Code, resp.ErrMsg)}
	}

	if !resp.Data.IsDeleted {
		return &AuthError{resp.RequestID, fmt.Errorf("deregister resource failed, error code: %d", resp.Code)}
	}

	return nil
}

func (a *authClient) updateResource(ctx context.Context, header http.Header, info *UpdateInfo) error {
	for k, v := range a.basicHeader {
		header[k] = v
	}
	resp := new(ResourceResult)
	url := fmt.Sprintf("/bkiam/api/v1/perm/systems/%s/resources", a.Config.SystemID)
	err := a.client.Put().
		SubResource(url).
		WithContext(ctx).
		WithHeaders(header).
		Body(info).
		Do().Into(resp)

	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return &AuthError{resp.RequestID, fmt.Errorf("update resource failed, error code: %d, message: %s", resp.Code, resp.ErrMsg)}
	}

	if !resp.Data.IsUpdated {
		return &AuthError{resp.RequestID, fmt.Errorf("update resource failed, error code: %d", resp.Code)}
	}

	return nil
}
