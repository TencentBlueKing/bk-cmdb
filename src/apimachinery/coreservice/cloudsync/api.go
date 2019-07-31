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

package cloudsync

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
)

func (c *cloud) CreateCloudSyncTask(ctx context.Context, header http.Header, input interface{}) (resp *metadata.Uint64DataResponse, err error) {
	resp = new(metadata.Uint64DataResponse)
	subPath := "/create/cloud/sync/task"

	err = c.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (c *cloud) DeleteCloudSyncTask(ctx context.Context, h http.Header, id int64) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/delete/cloud/sync/task/%v", id)

	err = c.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) UpdateCloudSyncTask(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/cloud/sync/task"

	err = c.client.Put().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) SearchCloudSyncTask(ctx context.Context, h http.Header, data interface{}) (resp *metadata.CloudTaskSearch, err error) {
	resp = new(metadata.CloudTaskSearch)
	subPath := "/search/cloud/sync/task"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) CreateConfirm(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Uint64DataResponse, err error) {
	resp = new(metadata.Uint64DataResponse)
	subPath := "/create/cloud/confirm"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) CheckTaskNameUnique(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Uint64Response, err error) {
	resp = new(metadata.Uint64Response)
	subPath := "/check/cloud/task/name"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) DeleteConfirm(ctx context.Context, h http.Header, id int64) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/delete/cloud/confirm/%v", id)

	err = c.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) SearchConfirm(ctx context.Context, h http.Header, data interface{}) (resp *metadata.FavoriteResult, err error) {
	resp = new(metadata.FavoriteResult)
	subPath := "/search/cloud/confirm"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) CreateSyncHistory(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Uint64Response, err error) {
	resp = new(metadata.Uint64Response)
	subPath := "/create/cloud/sync/history"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) SearchSyncHistory(ctx context.Context, h http.Header, data interface{}) (resp *metadata.FavoriteResult, err error) {
	resp = new(metadata.FavoriteResult)
	subPath := "/search/cloud/sync/history"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) CreateConfirmHistory(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/create/cloud/confirm/history"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) SearchConfirmHistory(ctx context.Context, h http.Header, data interface{}) (resp *metadata.FavoriteResult, err error) {
	resp = new(metadata.FavoriteResult)
	subPath := "/search/cloud/confirm/history"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
