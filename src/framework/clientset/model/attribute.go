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

type AttributeInterface interface {
	CreateAttribute(ctx *types.CreateAttributeCtx) (int64, error)
	DeleteAttribute(ctx *types.DeleteAttributeCtx) error
	UpdateAttribute(ctx *types.UpdateAttributeCtx) error
	GetAttribute(ctx *types.GetAttributeCtx) ([]types.Attribute, error)
}

var _ AttributeInterface = &attribute{}

type attribute struct {
	client rest.ClientInterface
}

func (a *attribute) CreateAttribute(ctx *types.CreateAttributeCtx) (int64, error) {
	resp := new(types.CreateAttributeResult)
	subPath := "/create/objectattr"
	err := a.client.Post().
		WithContext(ctx.Ctx).
		Body(ctx.Attribute).
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

func (a *attribute) DeleteAttribute(ctx *types.DeleteAttributeCtx) error {
	resp := new(types.Response)
	subPath := fmt.Sprintf("/delete/objectattr/%d", ctx.AttributeID)
	err := a.client.Delete().
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

func (a *attribute) UpdateAttribute(ctx *types.UpdateAttributeCtx) error {
	resp := new(types.Response)
	subPath := fmt.Sprintf("/update/objectattr/%d", ctx.AttributeID)
	err := a.client.Put().
		WithContext(ctx.Ctx).
		Body(ctx.Attribute).
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

func (a *attribute) GetAttribute(ctx *types.GetAttributeCtx) ([]types.Attribute, error) {
	resp := new(types.GetAttributeResult)
	subPath := "/find/objectattr"
	err := a.client.Post().
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
	return resp.Data, nil
}
