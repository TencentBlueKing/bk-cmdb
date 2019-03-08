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

package privilege

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
)

func (t *privilege) CreateUserGroup(ctx context.Context, ownerID string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/privilege/group/%s", ownerID)

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *privilege) UpdateUserGroup(ctx context.Context, ownerID string, groupID string, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/privilege/group/%s/%s", ownerID, groupID)

	err = t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *privilege) DeleteUserGroup(ctx context.Context, ownerID string, groupID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/privilege/group/%s/%s", ownerID, groupID)

	err = t.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *privilege) SearchUserGroup(ctx context.Context, ownerID string, h http.Header, dat interface{}) (resp *metadata.PermissionGroupListResult, err error) {
	resp = new(metadata.PermissionGroupListResult)
	subPath := fmt.Sprintf("/privilege/group/%s/search", ownerID)

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *privilege) CreateUserGroupPrivi(ctx context.Context, ownerID string, groupID string, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/privilege/group/detail/%s/%s", ownerID, groupID)

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *privilege) UpdateUserGroupPrivi(ctx context.Context, ownerID string, groupID string, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/privilege/group/detail/%s/%s", ownerID, groupID)

	err = t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *privilege) GetUserGroupPrivi(ctx context.Context, ownerID string, groupID string, h http.Header) (resp *metadata.GroupPriviResult, err error) {
	resp = new(metadata.GroupPriviResult)
	subPath := fmt.Sprintf("/privilege/group/detail/%s/%s", ownerID, groupID)

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *privilege) CreateRolePri(ctx context.Context, ownerID string, objID string, propertyID string, h http.Header, role []string) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/role/%s/%s/%s", ownerID, objID, propertyID)

	err = t.client.Post().
		WithContext(ctx).
		Body(role).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *privilege) GetRolePri(ctx context.Context, ownerID string, objID string, propertyID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	resp.Data = []interface{}{}
	subPath := fmt.Sprintf("/role/%s/%s/%s", ownerID, objID, propertyID)

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *privilege) UpdateRolePri(ctx context.Context, ownerID string, objID string, propertyID string, h http.Header, role []string) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/role/%s/%s/%s", ownerID, objID, propertyID)

	err = t.client.Put().
		WithContext(ctx).
		Body(role).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *privilege) GetSystemFlag(ctx context.Context, ownerID string, flag string, h http.Header) (resp *metadata.PermissionSystemResponse, err error) {
	resp = new(metadata.PermissionSystemResponse)
	subPath := fmt.Sprintf("/system/%s/%s", flag, ownerID)

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
