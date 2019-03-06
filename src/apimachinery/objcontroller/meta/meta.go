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

	"configcenter/src/common/metadata"
)

func (t *meta) SelectObjects(ctx context.Context, h http.Header, data interface{}) (resp *metadata.QueryObjectResult, err error) {
	subPath := "/meta/objects"
	resp = new(metadata.QueryObjectResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) DeleteObject(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metadata.DeleteResult, err error) {
	subPath := fmt.Sprintf("/meta/object/%d", objID)
	resp = new(metadata.DeleteResult)
	err = t.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) CreateObject(ctx context.Context, h http.Header, dat *metadata.Object) (resp *metadata.CreateObjectResult, err error) {
	subPath := "/meta/object"
	resp = new(metadata.CreateObjectResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) UpdateObject(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metadata.UpdateResult, err error) {
	subPath := fmt.Sprintf("/meta/object/%d", objID)
	resp = new(metadata.UpdateResult)
	err = t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) SelectObjectAssociations(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.QueryObjectAssociationResult, err error) {
	subPath := "/meta/objectassts"
	resp = new(metadata.QueryObjectAssociationResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) DeleteObjectAssociation(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metadata.DeleteResult, err error) {
	subPath := fmt.Sprintf("/meta/objectasst/%d", objID)
	resp = new(metadata.DeleteResult)
	err = t.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) CreateObjectAssociation(ctx context.Context, h http.Header, dat *metadata.Association) (resp *metadata.CreateResult, err error) {
	subPath := "/meta/objectasst"
	resp = new(metadata.CreateResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) CreateMainlineObjectAssociation(ctx context.Context, h http.Header, dat *metadata.Association) (resp *metadata.CreateResult, err error) {
	subPath := "/meta/mainlineobjectasst"
	resp = new(metadata.CreateResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) UpdateObjectAssociation(ctx context.Context, asstID int64, h http.Header, dat map[string]interface{}) (resp *metadata.UpdateResult, err error) {
	subPath := fmt.Sprintf("/meta/objectasst/%d", asstID)
	resp = new(metadata.UpdateResult)
	err = t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) SelectObjectAttByID(ctx context.Context, objID int64, h http.Header) (resp *metadata.QueryObjectAttributeResult, err error) {
	resp = new(metadata.QueryObjectAttributeResult)

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

func (t *meta) SelectObjectAttWithParams(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.QueryObjectAttributeResult, err error) {
	resp = new(metadata.QueryObjectAttributeResult)
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

func (t *meta) DeleteObjectAttByID(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metadata.DeleteResult, err error) {
	subPath := fmt.Sprintf("/meta/objectatt/%d", objID)
	resp = new(metadata.DeleteResult)
	err = t.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) CreateObjectAtt(ctx context.Context, h http.Header, dat *metadata.Attribute) (resp *metadata.CreateObjectAttributeResult, err error) {
	subPath := "/meta/objectatt"
	resp = new(metadata.CreateObjectAttributeResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) UpdateObjectAttByID(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metadata.UpdateResult, err error) {
	subPath := fmt.Sprintf("/meta/objectatt/%d", objID)
	resp = new(metadata.UpdateResult)
	err = t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
