/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package containerserver

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

func (cs *containerServer) CreatePod(ctx context.Context, h http.Header, bizID int64, data interface{}) (resp *metadata.CreatedOneOptionResult, err error) {
	resp = new(metadata.CreatedOneOptionResult)
	subPath := "/create/container/biz/%d/pod"

	err = cs.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (cs *containerServer) CreateManyPod(ctx context.Context, h http.Header, bizID int64, data interface{}) (resp *metadata.CreatedManyOptionResult, err error) {
	resp = new(metadata.CreatedManyOptionResult)
	subPath := "/createmany/container/biz/%d/pod"

	err = cs.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (cs *containerServer) UpdatePod(ctx context.Context, h http.Header, bizID int64, data interface{}) (resp *metadata.UpdatedOptionResult, err error) {
	resp = new(metadata.UpdatedOptionResult)
	subPath := "/update/container/biz/%d/pod"

	err = cs.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (cs *containerServer) DeletePod(ctx context.Context, h http.Header, bizID int64, data interface{}) (resp *metadata.DeletedOptionResult, err error) {
	resp = new(metadata.DeletedOptionResult)
	subPath := "/delete/container/biz/%d/pod"

	err = cs.client.Delete().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (cs *containerServer) ListPods(ctx context.Context, h http.Header, bizID int64, data interface{}) (resp *metadata.ListPodsResult, err error) {
	resp = new(metadata.ListPodsResult)
	subPath := "/list/container/pod"

	err = cs.client.Delete().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
