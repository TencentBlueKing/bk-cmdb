/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package object

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

// CreateClassification TODO
func (t *object) CreateClassification(ctx context.Context, h http.Header,
	obj *metadata.Classification) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/create/objectclassification"

	err = t.client.Post().
		WithContext(ctx).
		Body(obj).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SelectClassificationWithObjects TODO
func (t *object) SelectClassificationWithObjects(ctx context.Context, h http.Header,
	data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/find/classificationobject"

	err = t.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SelectClassificationWithParams TODO
func (t *object) SelectClassificationWithParams(ctx context.Context, h http.Header,
	data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/find/objectclassification"

	err = t.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateClassification TODO
func (t *object) UpdateClassification(ctx context.Context, classID string, h http.Header,
	data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/objectclassification/%s"

	err = t.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, classID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// DeleteClassification TODO
func (t *object) DeleteClassification(ctx context.Context, classID string, h http.Header,
	data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/delete/objectclassification/%s"

	err = t.client.Delete().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, classID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
