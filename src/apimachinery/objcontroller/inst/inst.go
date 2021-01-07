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

package inst

import (
	"context"
	"fmt"
	"net/http"

	metatype "configcenter/src/common/metadata"
)

func (t *instance) SearchObjects(ctx context.Context, objType string, h http.Header, dat *metatype.QueryInput) (resp *metatype.QueryInstResult, err error) {
	subPath := fmt.Sprintf("/insts/%s/search", objType)
	resp = new(metatype.QueryInstResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *instance) CreateObject(ctx context.Context, objType string, h http.Header, dat interface{}) (resp *metatype.CreateInstResult, err error) {
	resp = new(metatype.CreateInstResult)
	subPath := fmt.Sprintf("/insts/%s", objType)

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *instance) DelObject(ctx context.Context, objType string, h http.Header, dat map[string]interface{}) (resp *metatype.DeleteResult, err error) {
	resp = new(metatype.DeleteResult)
	subPath := fmt.Sprintf("/insts/%s", objType)

	err = t.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *instance) UpdateObject(ctx context.Context, objType string, h http.Header, dat map[string]interface{}) (resp *metatype.UpdateResult, err error) {
	resp = new(metatype.UpdateResult)
	subPath := fmt.Sprintf("/insts/%s", objType)

	err = t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
