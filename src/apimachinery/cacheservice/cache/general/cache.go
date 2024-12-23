/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package general

import (
	"context"
	"net/http"

	fullsynccond "configcenter/pkg/cache/full-sync-cond"
	"configcenter/pkg/cache/general"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
)

// CreateFullSyncCond create full sync cache condition
func (c *cache) CreateFullSyncCond(ctx context.Context, h http.Header, opt *fullsynccond.CreateFullSyncCondOpt) (
	int64, errors.CCErrorCoder) {

	resp := new(metadata.CreateResult)

	err := c.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/create/full/sync/cond").
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return 0, errors.CCHttpError
	}

	if !resp.Result {
		return 0, resp.CCError()
	}

	return resp.Data.ID, nil
}

// UpdateFullSyncCond update full sync cache condition
func (c *cache) UpdateFullSyncCond(ctx context.Context, h http.Header,
	opt *fullsynccond.UpdateFullSyncCondOpt) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)

	err := c.client.Put().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/update/full/sync/cond").
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if !resp.Result {
		return resp.CCError()
	}

	return nil
}

// DeleteFullSyncCond delete full sync cache condition
func (c *cache) DeleteFullSyncCond(ctx context.Context, h http.Header,
	opt *fullsynccond.DeleteFullSyncCondOpt) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)

	err := c.client.Delete().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/delete/full/sync/cond").
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if !resp.Result {
		return resp.CCError()
	}

	return nil
}

// ListFullSyncCond list full sync cache condition
func (c *cache) ListFullSyncCond(ctx context.Context, h http.Header, opt *fullsynccond.ListFullSyncCondOpt) (
	*fullsynccond.ListFullSyncCondRes, errors.CCErrorCoder) {

	resp := new(fullsynccond.ListFullSyncCondResp)

	err := c.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/findmany/full/sync/cond").
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if !resp.Result {
		return nil, resp.CCError()
	}

	return resp.Data, nil
}

// ListCacheByFullSyncCond list resource cache by full sync condition
func (c *cache) ListCacheByFullSyncCond(ctx context.Context, h http.Header,
	opt *fullsynccond.ListCacheByFullSyncCondOpt) (*general.ListGeneralCacheRes, errors.CCErrorCoder) {

	resp := new(general.ListGeneralCacheResp)

	err := c.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/findmany/resource/by_full_sync_cond").
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if !resp.Result {
		return nil, resp.CCError()
	}

	return resp.Data, nil
}

// ListGeneralCacheByIDs list general resource cache by ids
func (c *cache) ListGeneralCacheByIDs(ctx context.Context, h http.Header, opt *general.ListDetailByIDsOpt) (
	*general.ListGeneralCacheRes, errors.CCErrorCoder) {

	httpheader.SetIsInnerReqHeader(h)

	resp := new(general.ListGeneralCacheResp)
	err := c.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/findmany/resource/by_ids").
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if !resp.Result {
		return nil, resp.CCError()
	}

	return resp.Data, nil
}

// ListGeneralCacheByUniqueKey list general resource cache by unique keys
func (c *cache) ListGeneralCacheByUniqueKey(ctx context.Context, h http.Header,
	opt *general.ListDetailByUniqueKeyOpt) (*general.ListGeneralCacheRes, errors.CCErrorCoder) {

	httpheader.SetIsInnerReqHeader(h)

	resp := new(general.ListGeneralCacheResp)
	err := c.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/findmany/resource/by_unique_keys").
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if !resp.Result {
		return nil, resp.CCError()
	}

	return resp.Data, nil
}
