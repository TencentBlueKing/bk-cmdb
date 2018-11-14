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

package apiserver

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (a *apiServer) AddDefaultApp(ctx context.Context, h http.Header, ownerID string, params mapstr.MapStr) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("biz/default/%s", ownerID)

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) SearchDefaultApp(ctx context.Context, h http.Header, ownerID string, params mapstr.MapStr) (resp *metadata.QueryInstResult, err error) {
	resp = new(metadata.QueryInstResult)
	subPath := fmt.Sprintf("biz/default/%s/search", ownerID)

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) GetRolePrivilege(ctx context.Context, h http.Header, ownerID, objID, role string) (resp *metadata.RolePriResult, err error) {
	resp = new(metadata.RolePriResult)
	subPath := fmt.Sprintf("/topo/privilege/%s/%s/%s", ownerID, objID, role)

	err = a.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) GetAppRole(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.RoleAppResult, err error) {

	resp = new(metadata.RoleAppResult)
	subPath := fmt.Sprintf("/object/attr/search")

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) GetUserPrivilegeApp(ctx context.Context, h http.Header, ownerID, userName string, params mapstr.MapStr) (resp *metadata.AppQueryResult, err error) {

	resp = new(metadata.AppQueryResult)
	subPath := fmt.Sprintf("biz/search/%s")

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) GetUserPrivilegeConfig(ctx context.Context, h http.Header, ownerID, userName string) (resp *metadata.UserPriviResult, err error) {
	resp = new(metadata.UserPriviResult)
	subPath := fmt.Sprintf("topo/privilege/user/detail/%s/%s", ownerID, userName)

	err = a.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) GetAllMainLineObject(ctx context.Context, h http.Header, ownerID, userName string) (resp *metadata.MainLineResult, err error) {
	resp = new(metadata.MainLineResult)
	subPath := fmt.Sprintf("topo/model/%s", ownerID)

	err = a.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) GetObjectData(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ObjectAttrBatchResult, err error) {
	resp = new(metadata.ObjectAttrBatchResult)
	subPath := fmt.Sprintf("object/search/batch")

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) GetInstDetail(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (resp *metadata.QueryInstResult, err error) {

	resp = new(metadata.QueryInstResult)
	subPath := fmt.Sprintf("inst/search/owner/%s/object/%s/detail", ownerID, objID)

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) GetObjectAttr(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ObjectAttrResult, err error) {

	resp = new(metadata.ObjectAttrResult)
	subPath := fmt.Sprintf("object/attr/search")

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) GetHostData(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.QueryInstResult, err error) {

	resp = new(metadata.QueryInstResult)
	subPath := fmt.Sprintf("hosts/search/asstdetail")

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) GetObjectGroup(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (resp *metadata.ObjectAttrGroupResult, err error) {

	resp = new(metadata.ObjectAttrGroupResult)
	subPath := fmt.Sprintf("objectatt/group/property/owner/%s/object/%s", ownerID, objID)
	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) AddHost(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.Response, err error) {

	resp = new(metadata.Response)
	subPath := fmt.Sprintf("hosts/add")

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) AddInst(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (resp *metadata.Response, err error) {

	resp = new(metadata.Response)
	subPath := fmt.Sprintf("inst/%s/%s", ownerID, objID)

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (a *apiServer) AddObjectBatch(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (resp *metadata.Response, err error) {

	resp = new(metadata.Response)
	subPath := "object/batch"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
