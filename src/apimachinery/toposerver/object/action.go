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

func (t *object) CreateModel(ctx context.Context, h http.Header, model *metadata.MainLineObject) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/create/topomodelmainline"

	err = t.client.Post().
		WithContext(ctx).
		Body(model).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *object) DeleteModel(ctx context.Context, objID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/delete/topomodelmainline/object/%s"

	err = t.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *object) SelectModel(ctx context.Context, h http.Header) (resp *metadata.MainlineObjectTopoResult, err error) {
	resp = new(metadata.MainlineObjectTopoResult)
	subPath := "/find/topomodelmainline"

	err = t.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *object) SelectModelByClsID(ctx context.Context, ownerID string, clsID string, objID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/topo/model/%s/%s/%s"

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, ownerID, clsID, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *object) SelectInst(ctx context.Context, bizID int64, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/find/topoinst/biz/%d"

	err = t.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
