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

package openapi

import (
	"context"
	"fmt"
	"net/http"

    "configcenter/src/common/metadata"
)

func (t *openapi) SearchAllApp(ctx context.Context, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/app/searchAll"

	err = t.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *openapi) UpdateMultiModule(ctx context.Context, appID string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/openapi/module/multi/%s", appID)

	err = t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *openapi) SearchModuleByApp(ctx context.Context, appID string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/openapi/module/searchByApp/%s", appID)

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *openapi) SearchModuleByProperty(ctx context.Context, appID string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/openapi/module/searchByProperty/%s", appID)

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *openapi) AddMultiModule(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/openapi/module/multi"

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *openapi) DeleteMultiModule(ctx context.Context, appID string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/openapi/module/multi/%s", appID)

	err = t.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *openapi) UpdateMultiSet(ctx context.Context, appID string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/openapi/set/multi/%s", appID)

	err = t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *openapi) DeleteMultiSet(ctx context.Context, appID string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/openapi/set/multi/%s", appID)

	err = t.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *openapi) DeleteSetHost(ctx context.Context, appID string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/openapi/set/setHost/%s", appID)

	err = t.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
