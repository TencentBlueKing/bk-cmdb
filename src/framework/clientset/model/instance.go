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

package model

import (
	"fmt"

	"configcenter/src/framework/clientset/types"
	"configcenter/src/framework/common/rest"
)

type InstanceInterface interface {
}

type instClient struct {
	client rest.ClientInterface
}

func (s instClient) CreateObjectInstance(ctx *types.CreateSetCtx) (int64, error) {
	resp := new(types.CreateInstanceResult)
	subPath := fmt.Sprintf("/set/%d", ctx.SetID)
	err := s.client.Post().
		WithContext(ctx.Ctx).
		Body(ctx.Set).
		SubResource(subPath).
		WithHeaders(ctx.Header).
		Do().
		Into(resp)

	if err != nil {
		return 0, &types.ErrorDetail{Code: types.HttpRequestFailed, Message: err.Error()}
	}

	if !resp.BaseResp.Result {
		return 0, &types.ErrorDetail{Code: resp.Code, Message: resp.ErrMsg}
	}
	return resp.Data.ID, nil
}

func (s instClient) DeleteObjectInstance(ctx *types.DeleteObjectCtx) error {
	resp := new(types.Response)
	subPath := fmt.Sprintf("/delete/instance/object/%s/inst/%d", ctx.ObjectID, ctx.InstanceID)
	err := s.client.Delete().
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

func (s instClient) UpdateObjectInstance(ctx *types.UpdateObjectCtx) error {
	resp := new(types.Response)
	subPath := fmt.Sprintf("/update/instance/object/%s/inst/%d", ctx.ObjectID, ctx.InstanceID)
	err := s.client.Put().
		WithContext(ctx.Ctx).
		Body(ctx.Object).
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

func (s instClient) ListObjectInstance(ctx *types.ListInstanceCtx) (*types.ListInfo, error) {
	resp := new(types.ListInstanceResult)
	subPath := fmt.Sprintf("/inst/search/owener/%s/object/%s", ctx.Tenancy, ctx.ObjectID)
	err := s.client.Post().
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
