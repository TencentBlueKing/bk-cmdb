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

// Package tenantset is the tenant set types
package tenantset

const (
	// DefaultTenantSetID is the default tenant set id
	DefaultTenantSetID = 1
	// DefaultTenantSetName is the default tenant set name
	DefaultTenantSetName = "All tenants"
	// DefaultTenantSetStr is the json string of the default tenant set that matches all tenants
	// right now we only support one default tenant set, so we mock it in code
	DefaultTenantSetStr = `{"id":1,"name":"All tenants","maintainer":"cc_system","description":"全租户","default":1,"bk_scope":{"match_all":true},"bk_created_at":"2025-03-03 10:00:00","bk_created_by":"cc_system","bk_updated_at":"2025-03-03 10:00:00","bk_updated_by":"cc_system"}`
)

// ListTenantResult is the result of list tenant in tenant set
type ListTenantResult struct {
	Data []Tenant `json:"data"`
}

// Tenant is the tenant info for list tenant api
type Tenant struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
