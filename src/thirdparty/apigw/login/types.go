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

import "configcenter/src/thirdparty/apigw/apigwutil"

// BkLoginResponse is bk login api gateway response
type BkLoginResponse[T any] struct {
	apigwutil.ApiGWBaseResponse `json:",inline"`
	Error                       *BkLoginError `json:"error"`
	Data                        T             `json:"data"`
}

// BkLoginError is bk login api gateway error
type BkLoginError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// VerifyTokenRes is the result of the verify token api
type VerifyTokenRes struct {
	Username string `json:"bk_username"`
	TenantID string `json:"tenant_id"`
}

// UserInfo is the user info
type UserInfo struct {
	Username    string `json:"bk_username"`
	TenantID    string `json:"tenant_id"`
	DisplayName string `json:"display_name"`
	Language    string `json:"language"`
	TimeZone    string `json:"time_zone"`
}
