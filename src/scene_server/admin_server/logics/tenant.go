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

package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/storage/dal/mongo/local"
)

// NewTenantInterface get new tenant cli interface
type NewTenantInterface interface {
	NewTenantCli(tenant string) local.DB
	NewTenantDB() string
}

// GetNewTenantCli get new tenant db
func GetNewTenantCli(kit *rest.Kit, cli interface{}) (local.DB, string) {
	newTenantCli := cli.(NewTenantInterface)
	return newTenantCli.NewTenantCli(kit.TenantID), newTenantCli.NewTenantDB()
}

// GetSystemTenant get system tenant # TODO get the default tenant when multi-tenancy is not enabled
func GetSystemTenant() string {
	return common.BKDefaultTenantID
}
