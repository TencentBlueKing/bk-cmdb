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

// SelectObjectTopoGraphics TODO
func (t *object) SelectObjectTopoGraphics(ctx context.Context, scopeType string, scopeID string,
	h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/find/objecttopo/scope_type/%s/scope_id/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, scopeType, scopeID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateObjectTopoGraphics TODO
func (t *object) UpdateObjectTopoGraphics(ctx context.Context, scopeType string, scopeID string, h http.Header,
	data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/objecttopo/scope_type/%s/scope_id/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, scopeType, scopeID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
