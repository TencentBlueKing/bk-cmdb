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

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreateInstance create instance
func (inst *instance) CreateInstance(ctx context.Context, h http.Header, objID string,
	input *metadata.CreateModelInstance) (*metadata.CreateOneDataResult, error) {

	resp := new(metadata.CreatedOneOptionResult)
	subPath := "/create/model/%s/instance"

	err := inst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// CreateManyInstance batch create instances
func (inst *instance) CreateManyInstance(ctx context.Context, h http.Header, objID string,
	input *metadata.CreateManyModelInstance) (*metadata.CreateManyDataResult, errors.CCErrorCoder) {

	resp := new(metadata.CreatedManyOptionResult)
	subPath := "/createmany/model/%s/instance"

	err := inst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// SetManyInstance TODO
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

// UpdateInstance update instance
func (inst *instance) UpdateInstance(ctx context.Context, h http.Header, objID string, input *metadata.UpdateOption) (
	*metadata.UpdatedCount, error) {

	resp := new(metadata.UpdatedOptionResult)
	subPath := "/update/model/%s/instance"

	err := inst.client.Put().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// ReadInstance search instance
func (inst *instance) ReadInstance(ctx context.Context, h http.Header, objID string, input *metadata.QueryCondition) (
	*metadata.InstDataInfo, error) {

	resp := new(metadata.QueryConditionResult)
	subPath := "/read/model/%s/instances"

	err := inst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// DeleteInstance delete instance
func (inst *instance) DeleteInstance(ctx context.Context, h http.Header, objID string, input *metadata.DeleteOption) (
	*metadata.DeletedCount, error) {

	resp := new(metadata.DeletedOptionResult)
	subPath := "/delete/model/%s/instance"

	err := inst.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// DeleteInstanceCascade TODO
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

// ReadInstanceStruct TODO
//  ReadInstanceStruct 按照结构体返回实例数据
func (inst *instance) ReadInstanceStruct(ctx context.Context, h http.Header, objID string,
	input *metadata.QueryCondition, result interface{}) errors.CCErrorCoder {

	rid := util.GetHTTPCCRequestID(h)
	subPath := "/read/model/%s/instances"

	err := inst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(result)

	if err != nil {
		blog.Errorf("search instance failed, err: %v, filter: %#v, rid: %s", err, input, rid)
		return errors.CCHttpError
	}

	return nil
}

// CountInstances counts target model instances num.
func (inst *instance) CountInstances(ctx context.Context, header http.Header, objID string, input *metadata.Condition) (
	*metadata.CountResponseContent, error) {

	resp := new(metadata.CountResponse)
	subPath := "/count/model/%s/instances"

	err := inst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// GetInstanceObjectMapping get instance to bk_obj_id mapping by instance ids
func (inst *instance) GetInstanceObjectMapping(ctx context.Context, header http.Header, ids []int64) (
	[]metadata.ObjectMapping, errors.CCErrorCoder) {

	resp := new(metadata.InstanceObjectMappingsResult)
	subPath := "/get/instance/object/mapping"
	input := metadata.GetInstanceObjectMappingsOption{IDs: ids}

	err := inst.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
