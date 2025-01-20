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

package login

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common"
	httpheader "configcenter/src/common/http/header"
)

// VerifyToken verify user token
func (l *login) VerifyToken(ctx context.Context, h http.Header, token string) (*VerifyTokenRes, error) {
	resp := new(BkLoginResponse[*VerifyTokenRes])

	httpheader.SetTenantID(h, common.BKDefaultTenantID)

	err := l.service.Client.Get().
		WithContext(ctx).
		WithParam("bk_token", token).
		SubResourcef("/login/api/v3/open/bk-tokens/verify/").
		WithHeaders(httpheader.SetBkAuth(h, l.service.Auth)).
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

// GetUserByToken get user info by token
func (l *login) GetUserByToken(ctx context.Context, h http.Header, token string) (*UserInfo, error) {
	resp := new(BkLoginResponse[*UserInfo])

	httpheader.SetTenantID(h, common.BKDefaultTenantID)

	err := l.service.Client.Get().
		WithContext(ctx).
		WithParam("bk_token", token).
		SubResourcef("/login/api/v3/open/bk-tokens/userinfo/").
		WithHeaders(httpheader.SetBkAuth(h, l.service.Auth)).
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
