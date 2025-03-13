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

import (
	"sync"

	"configcenter/pkg/tenant/types"
)

var (
	prevTenantInfo      = make(map[string]types.Tenant)
	tenantEventChannels = make(map[string]chan TenantEvent)
	tenantEventChLock   sync.RWMutex
)

// TenantEvent is the tenant event info
type TenantEvent struct {
	EventType EventType
	TenantID  string
}

// EventType is the tenant event type
type EventType string

const (
	// Create is the create or enable tenant event type
	Create EventType = "create"
	// Delete is the delete or disable tenant event type
	Delete EventType = "delete"
)

// NewTenantEventChan generate a new tenant event chan
func NewTenantEventChan(name string) <-chan TenantEvent {
	tenantEventChLock.Lock()
	defer tenantEventChLock.Unlock()

	if ch, exists := tenantEventChannels[name]; exists {
		return ch
	}

	eventChan := make(chan TenantEvent)
	tenantEventChannels[name] = eventChan
	go func() {
		for _, tenant := range allTenants {
			if tenant.Status == types.EnabledStatus {
				eventChan <- TenantEvent{
					EventType: Create,
					TenantID:  tenant.TenantID,
				}
			}
		}
	}()
	return eventChan
}

// RemoveTenantEventChan remove tenant event chan
func RemoveTenantEventChan(name string) {
	tenantEventChLock.Lock()
	defer tenantEventChLock.Unlock()

	ch, exists := tenantEventChannels[name]
	if !exists {
		return
	}

	close(ch)
	delete(tenantEventChannels, name)
}

// generateAndPushTenantEvent compare the tenant with the previous tenant info to generate and push event
func generateAndPushTenantEvent(tenants []types.Tenant) {
	tenantEventChLock.RLock()
	defer tenantEventChLock.RUnlock()

	prevTenantMap := make(map[string]types.Tenant)

	for _, tenant := range tenants {
		tenantID := tenant.TenantID
		prevTenantMap[tenantID] = tenant

		prevTenant, exists := prevTenantInfo[tenantID]
		if !exists && tenant.Status == types.EnabledStatus {
			for _, eventChan := range tenantEventChannels {
				eventChan <- TenantEvent{
					EventType: Create,
					TenantID:  tenantID,
				}
			}
			continue
		}

		if prevTenant.Status != tenant.Status {
			eventType := Create
			if tenant.Status == types.DisabledStatus {
				eventType = Delete
			}
			for _, eventChan := range tenantEventChannels {
				eventChan <- TenantEvent{
					EventType: eventType,
					TenantID:  tenantID,
				}
			}
		}

		delete(prevTenantInfo, tenantID)
	}

	for tenantID := range prevTenantInfo {
		for _, eventChan := range tenantEventChannels {
			eventChan <- TenantEvent{
				EventType: Delete,
				TenantID:  tenantID,
			}
		}
	}

	prevTenantInfo = prevTenantMap
}
