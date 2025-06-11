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

package tools

import (
	"fmt"

	"configcenter/src/common"
)

var defaultTenant = common.BKSingleTenantID

// InitDefaultTenant init default tenant
func InitDefaultTenant(enableMultiTenantMode bool) {
	if enableMultiTenantMode {
		defaultTenant = common.BKDefaultTenantID
	}
}

// GetDefaultTenant get default tenant
func GetDefaultTenant() string {
	return defaultTenant
}

// ValidateDisableTenantMode validate disable multi-tenant mode
func ValidateDisableTenantMode(tenantID string, enableTenantMode bool) (string, error) {
	if !enableTenantMode {
		if tenantID == "" || tenantID == common.BKSingleTenantID {
			return common.BKSingleTenantID, nil
		}

		return "", fmt.Errorf("tenant mode is disable, but tenant id %s is set", tenantID)
	}

	return tenantID, nil
}
