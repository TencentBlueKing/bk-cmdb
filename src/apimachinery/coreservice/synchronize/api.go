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

package synchronize

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

func (sync *synchronize) SynchronizeInstance(ctx context.Context, h http.Header, input *metadata.SynchronizeParameter) (resp *metadata.SynchronizeResult, err error) {
	resp = new(metadata.SynchronizeResult)
	subPath := "/set/synchronize/instance"

	err = sync.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (sync *synchronize) SynchronizeModel(ctx context.Context, h http.Header, input *metadata.SynchronizeParameter) (resp *metadata.SynchronizeResult, err error) {
	resp = new(metadata.SynchronizeResult)
	subPath := "/set/synchronize/model"

	err = sync.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (sync *synchronize) SynchronizeAssociation(ctx context.Context, h http.Header, input *metadata.SynchronizeParameter) (resp *metadata.SynchronizeResult, err error) {
	resp = new(metadata.SynchronizeResult)
	subPath := "/set/synchronize/association"

	err = sync.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (sync *synchronize) SynchronizeFind(ctx context.Context, h http.Header, input *metadata.SynchronizeFindInfoParameter) (resp *metadata.ResponseInstData, err error) {
	resp = new(metadata.ResponseInstData)
	subPath := "/read/synchronize"

	err = sync.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (sync *synchronize) SynchronizeClearData(ctx context.Context, h http.Header, input *metadata.SynchronizeClearDataParameter) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/clear/synchronize/data"

	err = sync.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (sync *synchronize) SetIdentifierFlag(ctx context.Context, h http.Header, input *metadata.SetIdenifierFlag) (resp *metadata.SynchronizeResult, err error) {
	resp = new(metadata.SynchronizeResult)
	subPath := "/set/synchronize/identifier/flag"

	err = sync.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
