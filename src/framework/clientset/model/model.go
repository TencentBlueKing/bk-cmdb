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

type ModelInterface interface {
	CreateModel(ctx *types.CreateModelCtx) (int64, error)
	DeleteModel(ctx *types.DeleteModelCtx) error
	UpdateModel(ctx *types.UpdateModelCtx) error
	GetModels(ctx *types.GetModelsCtx) ([]types.ModelInfo, error)
}

type modelClient struct {
	client rest.ClientInterface
}

func (m *modelClient) CreateModel(ctx *types.CreateModelCtx) (int64, error) {
	resp := new(types.CreateModelResponse)
	subPath := "/create/object"
	err := m.client.Post().
		WithContext(ctx.Ctx).
		Body(ctx.ModelInfo).
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

func (m *modelClient) DeleteModel(ctx *types.DeleteModelCtx) error {
	resp := new(types.Response)
	subPath := fmt.Sprintf("/delete/object/%d", ctx.ModelID)
	err := m.client.Delete().
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

func (m *modelClient) UpdateModel(ctx *types.UpdateModelCtx) error {
	resp := new(types.Response)
	subPath := fmt.Sprintf("/update/object/%d", ctx.ModelID)
	err := m.client.Put().
		WithContext(ctx.Ctx).
		Body(ctx.ModelInfo).
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

func (m *modelClient) GetModels(ctx *types.GetModelsCtx) ([]types.ModelInfo, error) {
	resp := new(types.GetModelsResult)
	subPath := "/find/object"
	err := m.client.Post().
		WithContext(ctx.Ctx).
		Body(ctx.Filters).
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
