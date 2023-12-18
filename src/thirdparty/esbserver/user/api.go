// Package user TODO
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
	"net/http"
	"net/url"

	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/esbserver/esbutil"
)

// GetAllUsers get all users from user manager
func (p *user) GetAllUsers(ctx context.Context, h http.Header) (resp *metadata.EsbUserListResponse, err error) {
	resp = &metadata.EsbUserListResponse{}
	subPath := "/v2/usermanage/get_all_users/"
	h.Set("Accept", "application/json")

	err = p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Do().
		Into(resp)

	return
}

// ListUsers get user list from user manager
func (p *user) ListUsers(ctx context.Context, h http.Header,
	params map[string]string) (resp *metadata.EsbListUserResponse, err error) {
	// response demo
	/*
		{
		  "status": "NORMAL",
		  "domain": "default.local",
		  "code": null,
		  "telephone": "11111111111",
		  "create_time": "2019-12-11T19:45:21.518283Z",
		  "country_code": "86",
		  "iso_code": "CN",
		  "logo": null,
		  "id": 22222,
		  "display_name": "foo",
		  "role": 0,
		  "type": "",
		  "leader": [],
		  "username": "foo",
		  "update_time": "2019-12-11T19:50:07.903422Z",
		  "wx_userid": "",
		  "staff_status": "IN",
		  "position": 0,
		  "qq": "",
		  "language": "zh-cn",
		  "enabled": true,
		  "time_zone": "Asia/Shanghai",
		  "departments": [
			{
			  "order": 1,
			  "id": 22222,
			  "full_name": "总公司",
			  "name": "总公司"
			}
		  ],
		  "email": "foo@qq.com",
		  "extras": [],
		  "wx_openid": "",
		  "password_valid_days": 180,
		  "category_id": 1
		}
	*/
	resp = &metadata.EsbListUserResponse{}
	h.Set("Accept", "application/json")

	if params == nil {
		params = make(map[string]string)
	}
	if _, ok := params["fields"]; ok == false {
		params["fields"] = "username,id,display_name"
	}
	params["no_page"] = "true"
	err = p.client.Get().
		WithContext(ctx).
		WithParams(params).
		SubResourcef("/v2/usermanage/list_users/").
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Do().
		Into(resp)

	return
}

// GetDepartment This function is used for front-end non paged fuzzy query of organizational data in user management.
// The parameters are included in the url
func (p *user) GetDepartment(ctx context.Context, h http.Header, u *url.URL) (
	resp *metadata.EsbDepartmentResponse, err error) {

	resp = &metadata.EsbDepartmentResponse{}
	subPath := "/v2/usermanage/list_departments/"
	h.Set("Accept", "application/json")

	request := p.client.Get().
		WithContext(ctx).
		// 通过WithParamsFromURL中的url参数指定page_size可以覆盖此默认的page_size
		WithParam("page_size", "500")

	if u != nil {
		request.WithParamsFromURL(u)
	}

	err = request.
		// 确保只使用分页查找，指定一个大page_size同样能达到不分页的效果
		WithParam("no_page", "false").
		SubResourcef(subPath).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Do().
		Into(resp)

	return
}

// GetAllDepartment This function is used for backend paging to accurately pull organizational data from user manager,
// without url parameters
func (p *user) GetAllDepartment(ctx context.Context, h http.Header, params map[string]string) (
	resp *metadata.EsbDepartmentResponse, err error) {

	resp = &metadata.EsbDepartmentResponse{}
	subPath := "/v2/usermanage/list_departments/"
	h.Set("Accept", "application/json")

	err = p.client.Get().
		WithContext(ctx).
		WithParams(params).
		SubResourcef(subPath).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Do().
		Into(resp)

	return
}

// GetDepartmentProfile get department profile from user manager
func (p *user) GetDepartmentProfile(ctx context.Context, q *http.Request) (resp *metadata.EsbDepartmentProfileResponse,
	err error) {
	resp = &metadata.EsbDepartmentProfileResponse{}
	subPath := "/v2/usermanage/list_department_profiles/"
	h := q.Header
	h.Set("Accept", "application/json")

	err = p.client.Get().
		WithContext(ctx).
		// 通过WithParamsFromURL中的url参数指定page_size可以覆盖此默认的page_size
		WithParam("page_size", "500").
		WithParamsFromURL(q.URL).
		// 确保只使用分页查找，指定一个大page_size同样能达到不分页的效果
		WithParam("no_page", "false").
		SubResourcef(subPath).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Do().
		Into(resp)

	return
}
