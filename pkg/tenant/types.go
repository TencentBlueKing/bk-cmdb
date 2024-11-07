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

package tenant

// Tenant is the tenant info
type Tenant struct {
	TenantID string `bson:"tenant_id"`
	Status   Status `bson:"status"`
	Database string `bson:"database"`
}

// Status is the tenant status
type Status string

const (
	// DisabledStatus is the disabled status for tenant
	DisabledStatus Status = "disabled"
	// EnabledStatus is the enabled status for tenant
	EnabledStatus Status = "enabled"
)
