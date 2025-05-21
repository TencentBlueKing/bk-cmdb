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

// Package tenant defines tenant related logics
package tenant

import (
	"sync"

	"configcenter/pkg/tenant/types"
)

var (
	allTenants = make([]types.Tenant, 0)
	tenantMap  = make(map[string]*types.Tenant)
	lock       sync.RWMutex
)

// GetTenant get tenant
func GetTenant(tenantID string) (*types.Tenant, bool) {
	lock.RLock()
	defer lock.RUnlock()
	if data, ok := tenantMap[tenantID]; ok {
		return data, ok
	}
	return nil, false
}

// SetTenant set tenant
func SetTenant(tenant []types.Tenant) {
	lock.Lock()
	allTenants = tenant
	tenantMap = make(map[string]*types.Tenant)
	for _, t := range allTenants {
		tenantMap[t.TenantID] = &t
	}
	lock.Unlock()
	generateAndPushTenantEvent(allTenants)
}

// GetAllTenants get all tenants
func GetAllTenants() []types.Tenant {
	lock.RLock()
	defer lock.RUnlock()
	return allTenants
}

// ExecForAllTenants execute handler for all tenants
func ExecForAllTenants(handler func(tenantID string) error) error {
	for _, tenant := range GetAllTenants() {
		if err := handler(tenant.TenantID); err != nil {
			return err
		}
	}
	return nil
}
