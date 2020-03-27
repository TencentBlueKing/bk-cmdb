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

package object

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

func (t *object) CreateObjectUnique(ctx context.Context, objID string, h http.Header, data *metadata.CreateUniqueRequest) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/create/objectunique/object/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *object) SearchObjectUnique(ctx context.Context, objID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/find/objectunique/object/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *object) UpdateObjectUnique(ctx context.Context, objID string, h http.Header, uniqueID uint64, data *metadata.UpdateUniqueRequest) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/objectunique/object/%s/unique/%d"

	err = t.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, objID, uniqueID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *object) DeleteObjectUnique(ctx context.Context, objID string, h http.Header, uniqueID uint64) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/delete/objectunique/object/%s/unique/%d"

	err = t.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, objID, uniqueID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
