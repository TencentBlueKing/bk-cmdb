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

	"configcenter/src/common"
	httpheader "configcenter/src/common/http/header"
)

// GetTenants get all tenants from bk user
func (u *user) GetTenants(ctx context.Context, h http.Header) ([]*Tenant, error) {
	resp := new(BkUserResponse[[]*Tenant])

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
