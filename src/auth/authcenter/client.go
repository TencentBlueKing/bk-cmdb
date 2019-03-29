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

	"configcenter/src/apimachinery/rest"
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
	url := fmt.Sprintf("/bkiam/api/v1/perm/systems/%s/resources-perms/batch-verify", a.Config.SystemID)
	err := a.client.Post().
		SubResource(url).
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
	url := fmt.Sprintf("/bkiam/api/v1/perm/systems/%s/any-resources-perms/batch-verify", a.Config.SystemID)
	err := a.client.Post().
		SubResource(url).
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
	util.CopyHeader(a.basicHeader, header)
	resp := new(ResourceResult)
	url := fmt.Sprintf("/bkiam/api/v1/perm/systems/%s/resources/batch-register", a.Config.SystemID)
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
	url := fmt.Sprintf("/bkiam/api/v1/perm/systems/%s/resources/batch-delete", a.Config.SystemID)
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
		return &AuthError{resp.RequestID, fmt.Errorf("update resource failed, error code: %d, message: %s", resp.Code, resp.Message)}
	}

	if !resp.Data.IsUpdated {
		return &AuthError{resp.RequestID, fmt.Errorf("update resource failed, error code: %d", resp.Code)}
	}

	return nil
}

func (a *authClient) QuerySystemInfo(ctx context.Context, header http.Header, systemID string, detail bool) (*SystemDetail, error) {
	util.CopyHeader(a.basicHeader, header)
	url := fmt.Sprintf("/bkiam/api/v1/perm-model/systems/%s", systemID)

	resp := struct {
		BaseResponse
		Data SystemDetail `json:"data"`
	}{}

	isDetail := "0"
	if detail {
		isDetail = "1"
	}

	err := a.client.Get().
		SubResource(url).
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
	url := fmt.Sprintf("/bkiam/api/v1/perm-model/systems/%s", system.SystemID)
	resp := struct {
		BaseResponse
		Data System `json:"data"`
	}{}

	err := a.client.Put().
		SubResource(url).
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
	url := fmt.Sprintf("/bkiam/api/v1/perm-model/systems/%s/scope-types/%s/resource-types/batch-register", systemID, scopeType)
	resp := BaseResponse{}

	err := a.client.Put().
		SubResource(url).
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
	url := fmt.Sprintf("/bkiam/api/v1/perm-model/systems/%s/scope-types/%s/resource-types/batch-update", systemID, scopeType)
	resp := BaseResponse{}

	err := a.client.Put().
		SubResource(url).
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
	url := fmt.Sprintf("/bkiam/api/v1/perm-model/systems/%s/scope-types/%s/resource-type-actions/batch-update", systemID, scopeType)
	resp := BaseResponse{}

	err := a.client.Put().
		SubResource(url).
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
	url := fmt.Sprintf("/bkiam/api/v1/perm-model/systems/%s/scope-types/%s/resource-types/batch-upsert", systemID, scopeType)
	resp := BaseResponse{}

	err := a.client.Post().
		SubResource(url).
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
	url := fmt.Sprintf("/bkiam/api/v1/perm-model/systems/%s/scope-types/%s/resource-types/%s", systemID, scopeType, resourceType)
	resp := BaseResponse{}

	err := a.client.Delete().
		SubResource(url).
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
	url := fmt.Sprintf("/bkiam/api/v1/perm/systems/%s/authorized-resources/search", SystemIDCMDB)
	resp := ListAuthorizedResourcesResult{}

	err := a.client.Delete().
		SubResource(url).
		WithContext(ctx).
		WithContext(ctx).
		WithHeaders(a.basicHeader).
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
