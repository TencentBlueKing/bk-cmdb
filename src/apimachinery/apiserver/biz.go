/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by bizlicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package apiserver

import (
	"context"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
)

func (a *apiServer) CreateBiz(ctx context.Context, ownerID string, h http.Header, params map[string]interface{}) (resp *metadata.CreateInstResult, err error) {
	resp = new(metadata.CreateInstResult)
	subPath := "/biz/%s"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, ownerID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) UpdateBiz(ctx context.Context, ownerID string, bizID string, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/biz/%s/%s"
	err = a.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, ownerID, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) UpdateBizDataStatus(ctx context.Context, ownerID string, flag common.DataStatusFlag, bizID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/biz/status/%s/%s/%s"
	err = a.client.Put().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, flag, ownerID, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) SearchBiz(ctx context.Context, ownerID string, h http.Header, s *params.SearchParams) (resp *metadata.SearchInstResult, err error) {
	resp = new(metadata.SearchInstResult)
	subPath := "/biz/search/%s"
	err = a.client.Post().
		WithContext(ctx).
		Body(s).
		SubResourcef(subPath, ownerID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
