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
	"net/http"

	"configcenter/src/common/metadata"
)

// TransferHostToInnerModule  transfer host to inner module  eg:idle module and fault module
func (h *host) TransferHostToInnerModule(ctx context.Context, header http.Header, input *metadata.TransferHostToInnerModule) (resp *metadata.OperaterException, err error) {
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

// TransferHostModule  transfer host to inner module  eg:idle module and fault module
func (h *host) TransferHostModule(ctx context.Context, header http.Header, input *metadata.HostsModuleRelation) (resp *metadata.OperaterException, err error) {
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

// TransferHostCrossBusiness  transfer host to inner module  eg:idle module and fault module
func (h *host) TransferHostCrossBusiness(ctx context.Context, header http.Header, input *metadata.TransferHostsCrossBusinessRequest) (resp *metadata.OperaterException, err error) {
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
