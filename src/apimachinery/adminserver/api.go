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

package adminserver

import (
	"context"
	"fmt"
	"net/http"

    "configcenter/src/common/metadata"
)

func (a *adminServer) ClearDatabase(ctx context.Context, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/clear"

	err = a.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *adminServer) Set(ctx context.Context, ownerID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/migrate/system/hostcrossbiz/%s", ownerID)

	err = a.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *adminServer) Migrate(ctx context.Context, ownerID string, distribution string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/migrate/%s/%s", distribution, ownerID)

	err = a.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
