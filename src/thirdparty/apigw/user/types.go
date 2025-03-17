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

import "configcenter/src/thirdparty/apigw/apigwutil"

// Tenant is the result from bk-user.
type Tenant struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status Status `json:"status"`
}

// Status is the tenant status
type Status string

const (
	// DisabledStatus is the disabled status for tenant
	DisabledStatus Status = "disabled"
	// EnabledStatus is the enabled status for tenant
	EnabledStatus Status = "enabled"
)

// BkUserResponse is bk user api gateway response
type BkUserResponse[T any] struct {
	apigwutil.ApiGWBaseResponse `json:",inline"`
	Error                       *BkUserError `json:"error"`
	Data                        T            `json:"data"`
}

// BkUserError is bk user api gateway error
type BkUserError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
