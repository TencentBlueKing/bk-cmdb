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

package user

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
)

func (u *user) AddUserConfig(ctx context.Context, h http.Header, dat *metadata.UserConfig) (resp *metadata.IDResult, err error) {
	resp = new(metadata.IDResult)
	subPath := "/userapi"

	err = u.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (u *user) UpdateUserConfig(ctx context.Context, businessID string, id string, h http.Header, dat map[string]interface{}) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := fmt.Sprintf("/userapi/%s/%s", businessID, id)

	err = u.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (u *user) DeleteUserConfig(ctx context.Context, businessID string, id string, h http.Header) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := fmt.Sprintf("/userapi/%s/%s", businessID, id)

	err = u.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (u *user) GetUserConfig(ctx context.Context, h http.Header, opt *metadata.QueryInput) (resp *metadata.GetUserConfigResult, err error) {
	resp = new(metadata.GetUserConfigResult)
	subPath := "/userapi/search"

	err = u.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (u *user) GetUserConfigDetail(ctx context.Context, businessID string, id string, h http.Header) (resp *metadata.GetUserConfigDetailResult, err error) {
	resp = new(metadata.GetUserConfigDetailResult)
	subPath := fmt.Sprintf("/userapi/detail/%s/%s", businessID, id)

	err = u.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (u *user) AddUserCustom(ctx context.Context, user string, h http.Header, dat map[string]interface{}) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := fmt.Sprintf("/usercustom/%s", user)

	err = u.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (u *user) UpdateUserCustomByID(ctx context.Context, user string, id string, h http.Header, dat map[string]interface{}) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := fmt.Sprintf("/usercustom/%s/%s", user, id)

	err = u.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (u *user) GetUserCustomByUser(ctx context.Context, user string, h http.Header) (resp *metadata.GetUserCustomResult, err error) {
	resp = new(metadata.GetUserCustomResult)
	subPath := fmt.Sprintf("/usercustom/user/search/%s", user)

	err = u.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (u *user) GetDefaultUserCustom(ctx context.Context, user string, h http.Header) (resp *metadata.GetUserCustomResult, err error) {
	resp = new(metadata.GetUserCustomResult)
	subPath := fmt.Sprintf("/usercustom/default/search/%s", user)

	err = u.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
