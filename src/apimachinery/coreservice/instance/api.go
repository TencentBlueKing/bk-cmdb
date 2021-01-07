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

package instance

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

func (inst *instance) CreateInstance(ctx context.Context, h http.Header, objID string, input *metadata.CreateModelInstance) (resp *metadata.CreatedOneOptionResult, err error) {
	resp = new(metadata.CreatedOneOptionResult)
	subPath := "/create/model/%s/instance"

	err = inst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (inst *instance) CreateManyInstance(ctx context.Context, h http.Header, objID string, input *metadata.CreateManyModelInstance) (resp *metadata.CreatedManyOptionResult, err error) {
	resp = new(metadata.CreatedManyOptionResult)
	subPath := "/createmany/model/%s/instance"

	err = inst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (inst *instance) SetManyInstance(ctx context.Context, h http.Header, objID string, input *metadata.SetManyModelInstance) (resp *metadata.SetOptionResult, err error) {
	resp = new(metadata.SetOptionResult)
	subPath := "/setmany/model/%s/instances"

	err = inst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (inst *instance) UpdateInstance(ctx context.Context, h http.Header, objID string, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error) {
	resp = new(metadata.UpdatedOptionResult)
	subPath := "/update/model/%s/instance"

	err = inst.client.Put().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (inst *instance) ReadInstance(ctx context.Context, h http.Header, objID string, input *metadata.QueryCondition) (resp *metadata.QueryConditionResult, err error) {
	resp = new(metadata.QueryConditionResult)
	subPath := "/read/model/%s/instances"

	err = inst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (inst *instance) DeleteInstance(ctx context.Context, h http.Header, objID string, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error) {
	resp = new(metadata.DeletedOptionResult)
	subPath := "/delete/model/%s/instance"

	err = inst.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (inst *instance) DeleteInstanceCascade(ctx context.Context, h http.Header, objID string, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error) {
	resp = new(metadata.DeletedOptionResult)
	subPath := "/delete/model/%s/instance/cascade"

	err = inst.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
