/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"configcenter/src/common"
	httpheader "configcenter/src/common/http/header"
)

// GetTenants get all tenants from bk user
func (u *user) GetTenants(ctx context.Context, h http.Header) ([]Tenant, error) {
	resp := new(BkUserResponse[[]Tenant])

	httpheader.SetTenantID(h, common.BKDefaultTenantID)
	err := u.service.Client.Get().
		WithContext(ctx).
		SubResourcef("/api/v3/open/tenants/").
		WithHeaders(httpheader.SetBkAuth(h, u.service.Auth)).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("code: %s, message: %s", resp.Error.Code, resp.Error.Message)
	}

	return resp.Data, nil
}

// ListUsers list users
func (u *user) ListUsers(ctx context.Context, h http.Header, page *PageOptions) (*ListUserResult, error) {
	resp := new(BkUserResponse[*ListUserResult])

	params := make(map[string]string)
	if page != nil {
		if page.Page != 0 {
			params["page"] = strconv.Itoa(page.Page)
		}
		if page.PageSize != 0 {
			params["page_size"] = strconv.Itoa(page.PageSize)
		}
	}

	err := u.service.Client.Get().
		WithContext(ctx).
		WithParams(params).
		SubResourcef("/api/v3/open/tenant/users/").
		WithHeaders(httpheader.SetBkAuth(h, u.service.Auth)).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("code: %s, message: %s", resp.Error.Code, resp.Error.Message)
	}

	return resp.Data, nil
}

// BatchQueryUserDisplayInfo batch query user display name info
func (u *user) BatchQueryUserDisplayInfo(ctx context.Context, h http.Header, opts *QueryUserDisplayInfoOpts) (
	[]UserDisplayInfo, error) {

	resp := new(BkUserResponse[[]UserDisplayInfo])

	err := u.service.Client.Get().
		WithContext(ctx).
		WithParam("bk_usernames", strings.Join(opts.BkUsernames, ",")).
		SubResourcef("/api/v3/open/tenant/users/-/display_info/").
		WithHeaders(httpheader.SetBkAuth(h, u.service.Auth)).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("code: %s, message: %s", resp.Error.Code, resp.Error.Message)
	}

	return resp.Data, nil
}

// BatchLookupDept batch lookup department
func (u *user) BatchLookupDept(ctx context.Context, h http.Header, opts *BatchLookupDeptOpts) ([]DepartmentItem,
	error) {

	strDeptIDs := make([]string, len(opts.DeptIDs))
	for i, deptID := range opts.DeptIDs {
		strDeptIDs[i] = strconv.FormatInt(deptID, 10)
	}

	params := map[string]string{
		"department_ids": strings.Join(strDeptIDs, ","),
	}
	if opts.WithOrgPath {
		params["with_organization_path"] = "true"
	}

	resp := new(BkUserResponse[[]DepartmentItem])

	err := u.service.Client.Get().
		WithContext(ctx).
		WithParams(params).
		SubResourcef("/api/v3/open/tenant/departments/-/lookup/").
		WithHeaders(httpheader.SetBkAuth(h, u.service.Auth)).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("code: %s, message: %s", resp.Error.Code, resp.Error.Message)
	}

	return resp.Data, nil
}
