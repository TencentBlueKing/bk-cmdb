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

package host

import (
    "fmt"
    
    "configcenter/src/framework/clientset/types"
    "configcenter/src/framework/common/rest"
    types2 "configcenter/src/framework/core/types"
)

type hostClient struct {
	client rest.ClientInterface
}

func (h *hostClient) ListHosts(ctx *types.ListHostsCtx) (*types.HostsInfo, error) {
	resp := new(types.ListHostResult)
	subPath := "/hosts/search"
	err := h.client.Post().
		WithContext(ctx.Ctx).
		Body(ctx.Filter).
		SubResource(subPath).
		WithHeaders(ctx.Header).
		Do().
		Into(resp)

	if err != nil {
		return nil, &types.ErrorDetail{Code: types.HttpRequestFailed, Message: err.Error()}
	}

	if !resp.BaseResp.Result {
		return nil, &types.ErrorDetail{Code: resp.Code, Message: resp.ErrMsg}
	}
	return &resp.Data, nil
}

func (h *hostClient) GetHostDetails(ctx *types.GetHostCtx) ([]types.HostAttribute, error) {
	resp := new(types.GetHostResult)
	subPath := fmt.Sprintf("/hosts/%s/%d", ctx.Tenancy, ctx.HostID)
	err := h.client.Get().
		WithContext(ctx.Ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(ctx.Header).
		Do().
		Into(resp)

	if err != nil {
		return nil, &types.ErrorDetail{Code: types.HttpRequestFailed, Message: err.Error()}
	}

	if !resp.BaseResp.Result {
		return nil, &types.ErrorDetail{Code: resp.Code, Message: resp.ErrMsg}
	}
	return resp.Data, nil
}

func (h *hostClient) GetHostSnapshot(ctx *types.GetHostSnapshotCtx) (types2.MapStr, error) {
	resp := new(types.GetHostSnapshotResult)
	subPath := fmt.Sprintf("/hosts/snapshot/%d", ctx.HostID)
	err := h.client.Get().
		WithContext(ctx.Ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(ctx.Header).
		Do().
		Into(resp)

	if err != nil {
		return nil, &types.ErrorDetail{Code: types.HttpRequestFailed, Message: err.Error()}
	}

	if !resp.BaseResp.Result {
		return nil, &types.ErrorDetail{Code: resp.Code, Message: resp.ErrMsg}
	}
	return resp.Data, nil
}

func (h *hostClient) UpdateHostsAttributes(ctx *types.UpdateHostsAttributesCtx) error {
	resp := new(types.GetHostSnapshotResult)
	subPath := "/hosts/batch"
	err := h.client.Put().
		WithContext(ctx.Ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(ctx.Header).
		Do().
		Into(resp)

	if err != nil {
		return &types.ErrorDetail{Code: types.HttpRequestFailed, Message: err.Error()}
	}

	if !resp.BaseResp.Result {
		return &types.ErrorDetail{Code: resp.Code, Message: resp.ErrMsg}
	}
	return nil
}

func (h *hostClient) DeleteHosts(ctx *types.DeleteHostsCtx) error {
	resp := new(types.GetHostSnapshotResult)
	subPath := "/hosts/batch"
	err := h.client.Delete().
		WithContext(ctx.Ctx).
		Body(ctx.Hosts).
		SubResource(subPath).
		WithHeaders(ctx.Header).
		Do().
		Into(resp)

	if err != nil {
		return &types.ErrorDetail{Code: types.HttpRequestFailed, Message: err.Error()}
	}

	if !resp.BaseResp.Result {
		return &types.ErrorDetail{Code: resp.Code, Message: resp.ErrMsg}
	}
	return nil
}
