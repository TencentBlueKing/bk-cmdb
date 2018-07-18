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

func (t *meta) CreatePropertyGroup(ctx context.Context, h http.Header, dat *metatype.Group) (resp *metatype.CreateObjectGroupResult, err error) {
	subPath := "/meta/objectatt/group/new"
	resp = new(metatype.CreateObjectGroupResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) UpdatePropertyGroup(ctx context.Context, h http.Header, dat *metatype.UpdateGroupCondition) (resp *metatype.UpdateResult, err error) {
	subPath := "/meta/objectatt/group/update"
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

func (t *meta) DeletePropertyGroup(ctx context.Context, groupID string, h http.Header) (resp *metatype.DeleteResult, err error) {
	subPath := fmt.Sprintf("/meta/objectatt/group/groupid/%s", groupID)
	resp = new(metatype.DeleteResult)
	err = t.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) UpdatePropertyGroupObjectAtt(ctx context.Context, h http.Header, dat []metatype.PropertyGroupObjectAtt) (resp *metatype.UpdateResult, err error) {
	subPath := "/meta/objectatt/group/property"
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

func (t *meta) DeletePropertyGroupObjectAtt(ctx context.Context, ownerID, objID, propertyID, groupID string, h http.Header) (resp *metatype.DeleteResult, err error) {
	subPath := fmt.Sprintf("/meta/objectatt/group/owner/%s/object/%s/propertyids/%s/groupids/%s", ownerID, objID, propertyID, groupID)
	resp = new(metatype.DeleteResult)
	err = t.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) SelectPropertyGroupByObjectID(ctx context.Context, ownerID string, objID string, h http.Header, dat map[string]interface{}) (resp *metatype.QueryObjectGroupResult, err error) {
	subPath := fmt.Sprintf("/meta/objectatt/group/property/owner/%s/object/%s", ownerID, objID)
	resp = new(metatype.QueryObjectGroupResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *meta) SelectGroup(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metatype.QueryObjectGroupResult, err error) {
	subPath := "/meta/objectatt/group/search"
	resp = new(metatype.QueryObjectGroupResult)
	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
