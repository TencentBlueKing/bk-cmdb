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

type ModuleInterface interface {
	CreateModule(ctx *types.CreateModuleCtx) (int64, error)
	DeleteModule(ctx *types.DeleteModuleCtx) error
	UpdateModule(ctx *types.UpdateModuleCtx) error
	ListModules(ctx *types.ListModulesCtx) (*types.ListInfo, error)
}

var _ ModuleInterface = &module{}

type module struct {
	client rest.ClientInterface
}

func (m *module) CreateModule(ctx *types.CreateModuleCtx) (int64, error) {
	resp := new(types.CreateModuleResult)
	subPath := fmt.Sprintf("/module/%d/%d", ctx.BusinessID, ctx.SetID)
	err := m.client.Post().
		WithContext(ctx.Ctx).
		Body(ctx.Module).
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

func (m *module) DeleteModule(ctx *types.DeleteModuleCtx) error {
	resp := new(types.Response)
	subPath := fmt.Sprintf("/module/%d/%d/%d", ctx.BusinessID, ctx.SetID, ctx.ModuleID)
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

func (m *module) UpdateModule(ctx *types.UpdateModuleCtx) error {
	resp := new(types.Response)
	subPath := fmt.Sprintf("/module/%d/%d/%d", ctx.BusinessID, ctx.SetID, ctx.ModuleID)
	err := m.client.Put().
		WithContext(ctx.Ctx).
		Body(ctx.Module).
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

func (m *module) ListModules(ctx *types.ListModulesCtx) (*types.ListInfo, error) {
	resp := new(types.ListModulesResult)
	subPath := fmt.Sprintf("/module/search/%s/%d/%d", ctx.Tenancy, ctx.BusinessID, ctx.SetID)
	err := m.client.Post().
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
