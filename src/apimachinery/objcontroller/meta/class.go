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

package meta

import (
	"context"
	"fmt"
	"net/http"

	metatype "configcenter/src/common/metadata"
)

func (t *meta) SelectClassificationWithObject(ctx context.Context, ownerID string, h http.Header, dat map[string]interface{}) (resp *metatype.QueryObjectClassificationWithObjectsResult, err error) {
	subPath := fmt.Sprintf("/meta/object/classification/%s/objects", ownerID)
	resp = new(metatype.QueryObjectClassificationWithObjectsResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) SelectClassifications(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metatype.QueryObjectClassificationResult, err error) {
	subPath := "/meta/object/classification/search"
	resp = new(metatype.QueryObjectClassificationResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) DeleteClassification(ctx context.Context, id int64, h http.Header, dat map[string]interface{}) (resp *metatype.DeleteResult, err error) {
	subPath := fmt.Sprintf("/meta/object/classification/%d", id)
	resp = new(metatype.DeleteResult)
	err = t.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) CreateClassification(ctx context.Context, h http.Header, dat *metatype.Classification) (resp *metatype.CreateObjectClassificationResult, err error) {
	subPath := "/meta/object/classification"
	resp = new(metatype.CreateObjectClassificationResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) UpdateClassification(ctx context.Context, id int64, h http.Header, dat map[string]interface{}) (resp *metatype.UpdateResult, err error) {
	subPath := fmt.Sprintf("/meta/object/classification/%d", id)
	resp = new(metatype.UpdateResult)
	err = t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
