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
	"context"
	"fmt"

	"configcenter/pkg/tenant"
	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/refresh"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/header/util"
	"configcenter/src/common/http/rest"
	"configcenter/src/storage/dal/mongo/local"
)

// NewTenantInterface get new tenant cli interface
type NewTenantInterface interface {
	NewTenantCli(tenant string) (local.DB, string, error)
}

// GetNewTenantCli get new tenant db
func GetNewTenantCli(kit *rest.Kit, cli interface{}) (local.DB, string, error) {
	newTenantCli, ok := cli.(NewTenantInterface)
	if !ok {
		blog.Errorf("get new tenant cli failed, rid: %s", kit.Rid)
		return nil, "", fmt.Errorf("get new tenant cli failed")
	}

	dbCli, dbUUID, err := newTenantCli.NewTenantCli(kit.TenantID)
	if err != nil || dbCli == nil {
		blog.Errorf("get new tenant cli failed, err: %v, tenant: %s, rid: %s", err, kit.TenantID, kit.Rid)
		return nil, "", fmt.Errorf("get new tenant cli failed, err: %v", err)
	}

	return dbCli, dbUUID, nil
}

// RefreshTenants refresh tenant info, skip tenant verify for apiserver
func RefreshTenants(coreAPI apimachinery.ClientSetInterface) error {
	clientSet, isOK := coreAPI.(*apimachinery.ClientSet)
	if !isOK {
		blog.Errorf("get client set from coreAPI failed")
		return fmt.Errorf("get client set from coreAPI failed")
	}

	refreshApiCli := refresh.NewRefreshClientInterface(apimachinery.GetRefreshCapability(clientSet))
	tenants, err := refreshApiCli.RefreshTenant(context.Background(), util.GenDefaultHeader())
	if err != nil {
		blog.Errorf("refresh tenant info failed, err: %v", err)
		return err
	}

	tenant.SetTenant(tenants)
	return nil
}
