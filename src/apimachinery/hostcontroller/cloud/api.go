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

package cloud

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
)

func (c *cloud) AddCloudTask(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/cloud/add"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) ResourceConfirm(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/cloud/confirm"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) TaskNameCheck(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/cloud/nameCheck"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) DeleteCloudTask(ctx context.Context, h http.Header, taskID string) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/hosts/cloud/delete/%v", taskID)

	err = c.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) SearchCloudTask(ctx context.Context, h http.Header, data interface{}) (resp *metadata.CloudTaskSearch, err error) {
	resp = new(metadata.CloudTaskSearch)
	subPath := "/hosts/cloud/search"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) UpdateCloudTask(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/cloud/update"

	err = c.client.Put().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) DeleteConfirm(ctx context.Context, h http.Header, ResourceID int64) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/hosts/cloud/confirm/delete/%v", ResourceID)

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
	subPath := "/hosts/cloud/confirm/search"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) AddSyncHistory(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/cloud/syncHistory/add"

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
	subPath := "/hosts/cloud/syncHistory/search"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (c *cloud) AddConfirmHistory(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/cloud/confirmHistory/add"

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
	subPath := "/hosts/cloud/confirmHistory/search"

	err = c.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
