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
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
)

// TransferToInnerModule  transfer host to inner module  eg:idle module and fault module
func (h *host) TransferToInnerModule(ctx context.Context, header http.Header, input *metadata.TransferHostToInnerModule) (resp *metadata.OperaterException, err error) {
	resp = new(metadata.OperaterException)
	subPath := "/set/module/host/relation/inner/module"

	err = h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// TransferHostModule  transfer host to  module
func (h *host) TransferToNormalModule(ctx context.Context, header http.Header, input *metadata.HostsModuleRelation) (resp *metadata.OperaterException, err error) {
	resp = new(metadata.OperaterException)
	subPath := "/set/module/host/relation/module"

	err = h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// RemoveFromModule 将主机从模块中移出
// 如果主机属于n+1个模块（n>0），操作之后，主机属于n个模块
// 如果主机属于1个模块, 且非空闲机模块，操作之后，主机属于空闲机模块
// 如果主机属于空闲机模块，操作失败
// 如果主机属于故障机模块，操作失败
// 如果主机不在参数指定的模块中，操作失败
func (h *host) RemoveFromModule(ctx context.Context, header http.Header, input *metadata.RemoveHostsFromModuleOption) (resp *metadata.OperaterException, err error) {
	resp = new(metadata.OperaterException)
	subPath := "/update/host/host_module_relations"

	err = h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// TransferHostCrossBusiness  transfer host to other bussiness module
func (h *host) TransferToAnotherBusiness(ctx context.Context, header http.Header, input *metadata.TransferHostsCrossBusinessRequest) (resp *metadata.OperaterException, err error) {
	resp = new(metadata.OperaterException)
	subPath := "/set/module/host/relation/cross/business"

	err = h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// GetHostModuleRelation get host module relation
func (h *host) GetHostModuleRelation(ctx context.Context, header http.Header, input *metadata.HostModuleRelationRequest) (resp *metadata.HostConfig, err error) {
	resp = new(metadata.HostConfig)
	subPath := "/read/module/host/relation"

	err = h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// DeleteHost delete host
func (h *host) DeleteHost(ctx context.Context, header http.Header, input *metadata.DeleteHostRequest) (resp *metadata.OperaterException, err error) {
	resp = new(metadata.OperaterException)
	subPath := "/delete/host"

	err = h.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// FindIdentifier  query host identifier
func (h *host) FindIdentifier(ctx context.Context, header http.Header, input *metadata.SearchHostIdentifierParam) (resp *metadata.SearchHostIdentifierResult, err error) {
	resp = new(metadata.SearchHostIdentifierResult)
	subPath := "/read/host/indentifier"

	err = h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) GetHostByID(ctx context.Context, header http.Header, hostID string) (resp *metadata.HostInstanceResult, err error) {
	resp = new(metadata.HostInstanceResult)
	subPath := fmt.Sprintf("/find/host/%s", hostID)

	err = h.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

func (h *host) GetHosts(ctx context.Context, header http.Header, opt *metadata.QueryInput) (resp *metadata.GetHostsResult, err error) {
	resp = new(metadata.GetHostsResult)
	subPath := "/findmany/hosts/search"

	err = h.client.Post().
		Body(opt).
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

func (h *host) GetHostSnap(ctx context.Context, header http.Header, hostID string) (resp *metadata.GetHostSnapResult, err error) {
	resp = new(metadata.GetHostSnapResult)
	subPath := fmt.Sprintf("/find/host/snapshot/%s", hostID)

	err = h.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

func (h *host) LockHost(ctx context.Context, header http.Header, input *metadata.HostLockRequest) (resp *metadata.HostLockResponse, err error) {
	resp = new(metadata.HostLockResponse)
	subPath := "/find/host/lock"

	err = h.client.Post().
		Body(input).
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

func (h *host) UnlockHost(ctx context.Context, header http.Header, input *metadata.HostLockRequest) (resp *metadata.HostLockResponse, err error) {
	resp = new(metadata.HostLockResponse)
	subPath := "/delete/host/lock"

	err = h.client.Delete().
		Body(input).
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

func (h *host) QueryHostLock(ctx context.Context, header http.Header, input *metadata.QueryHostLockRequest) (resp *metadata.HostLockQueryResponse, err error) {
	resp = new(metadata.HostLockQueryResponse)
	subPath := "/findmany/host/lock/search"

	err = h.client.Post().
		Body(input).
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}
