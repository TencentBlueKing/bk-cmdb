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

func (t *meta) SelectObjects(ctx context.Context, h http.Header, data interface{}) (resp *metatype.QueryObjectResult, err error) {
	subPath := "/meta/objects"
	resp = new(metatype.QueryObjectResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) DeleteObject(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metatype.DeleteResult, err error) {
	subPath := fmt.Sprintf("/meta/object/%d", objID)
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

func (t *meta) CreateObject(ctx context.Context, h http.Header, dat *metatype.Object) (resp *metatype.CreateObjectResult, err error) {
	subPath := "/meta/object"
	resp = new(metatype.CreateObjectResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) UpdateObject(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metatype.UpdateResult, err error) {
	subPath := fmt.Sprintf("/meta/object/%d", objID)
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

func (t *meta) SelectObjectAssociations(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metatype.QueryObjectAssociationResult, err error) {
	subPath := "/meta/objectassts"
	resp = new(metatype.QueryObjectAssociationResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) DeleteObjectAssociation(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metatype.DeleteResult, err error) {
	subPath := fmt.Sprintf("/meta/objectasst/%d", objID)
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

func (t *meta) CreateObjectAssociation(ctx context.Context, h http.Header, dat *metatype.Association) (resp *metatype.CreateResult, err error) {
	subPath := "/meta/objectasst"
	resp = new(metatype.CreateResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) UpdateObjectAssociation(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metatype.UpdateResult, err error) {
	subPath := fmt.Sprintf("/meta/objectasst/%d", objID)
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

func (t *meta) SelectObjectAttByID(ctx context.Context, objID int64, h http.Header) (resp *metatype.QueryObjectAttributeResult, err error) {
	resp = new(metatype.QueryObjectAttributeResult)

	subPath := fmt.Sprintf("/meta/objectatt/%d", objID)
	err = t.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) SelectObjectAttWithParams(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metatype.QueryObjectAttributeResult, err error) {
	resp = new(metatype.QueryObjectAttributeResult)
	subPath := "/meta/objectatts"

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) DeleteObjectAttByID(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metatype.DeleteResult, err error) {
	subPath := fmt.Sprintf("/meta/objectatt/%d", objID)
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

func (t *meta) CreateObjectAtt(ctx context.Context, h http.Header, dat *metatype.Attribute) (resp *metatype.CreateObjectAttributeResult, err error) {
	subPath := "/meta/objectatt"
	resp = new(metatype.CreateObjectAttributeResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) UpdateObjectAttByID(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metatype.UpdateResult, err error) {
	subPath := fmt.Sprintf("/meta/objectatt/%d", objID)
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
