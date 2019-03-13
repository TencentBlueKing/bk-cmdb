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

package business

import (
	"fmt"

	"configcenter/src/framework/clientset/types"
	"configcenter/src/framework/common/rest"
	"configcenter/src/framework/core/errors"
	types2 "configcenter/src/framework/core/types"
)

type biz struct {
	client rest.ClientInterface
}

func (b *biz) CreateBusiness(info *types.CreateBusinessCtx) (types2.MapStr, error) {
	resp := new(types.BusinessResponse)
	subPath := fmt.Sprintf("/biz/%s", info.Tenancy)
	err := b.client.Post().
		WithContext(info.Ctx).
		Body(info.BusinessInfo).
		SubResource(subPath).
		WithHeaders(info.Header).
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

func (b *biz) UpdateBusiness(info *types.UpdateBusinessCtx) error {
	resp := new(types.BusinessResponse)
	subPath := fmt.Sprintf("/biz/%s/%d", info.Tenancy, info.BusinessID)
	err := b.client.Put().
		WithContext(info.Ctx).
		Body(info.BusinessInfo).
		SubResource(subPath).
		WithHeaders(info.Header).
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

func (b *biz) DeleteBusiness(info *types.DeleteBusinessCtx) error {
	resp := new(types.BusinessResponse)
	subPath := fmt.Sprintf("/biz/%s/%d", info.Tenancy, info.BusinessID)
	err := b.client.Delete().
		WithContext(info.Ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(info.Header).
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

func (b *biz) ListBusiness(info *types.ListBusinessCtx) ([]types2.MapStr, error) {
	if len(info.Tenancy) == 0 {
		return nil, errors.New("business's tenancy can not be empty.")
	}

	resp := new(types.ListBusinessResult)
	subPath := fmt.Sprintf("/biz/search/%s", info.Tenancy)
	err := b.client.Post().
		WithContext(info.Ctx).
		Body(info.QueryInfo).
		SubResource(subPath).
		WithHeaders(info.Header).
		Do().
		Into(resp)

	if err != nil {
		return nil, &types.ErrorDetail{Code: types.HttpRequestFailed, Message: err.Error()}
	}

	if !resp.BaseResp.Result {
		return nil, &types.ErrorDetail{Code: resp.Code, Message: resp.ErrMsg}
	}
	return resp.Data.Info, nil
}
