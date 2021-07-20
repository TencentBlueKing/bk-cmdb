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

package resourcedir

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

func (r *ResourceDirectory) CreateResourceDirectory(ctx context.Context, header http.Header, data map[string]interface{}) (resp *metadata.CreatedOneOptionResult, err error) {
	resp = new(metadata.CreatedOneOptionResult)
	subPath := "/create/resource/directory"

	err = r.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (r *ResourceDirectory) UpdateResourceDirectory(ctx context.Context, header http.Header, moduleID int64, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/resource/directory/%d"

	err = r.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, moduleID).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (r *ResourceDirectory) SearchResourceDirectory(ctx context.Context, header http.Header, data map[string]interface{}) (resp *metadata.SearchResp, err error) {
	resp = new(metadata.SearchResp)
	subPath := "/findmany/resource/directory"

	err = r.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (r *ResourceDirectory) DeleteResourceDirectory(ctx context.Context, header http.Header, moduleID int64) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/delete/resource/directory/%d"

	err = r.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, moduleID).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}
