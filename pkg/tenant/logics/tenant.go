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
	"sync"
	"time"

	"configcenter/pkg/tenant"
	"configcenter/pkg/tenant/types"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/header/util"
	commontypes "configcenter/src/common/types"
	"configcenter/src/storage/dal/mongo/local"
)

var (
	once            sync.Once
	db              local.DB
	apiMachineryCli apimachinery.ClientSetInterface
)

// InitTenant init tenant, refresh tenants info while server is starting
func InitTenant(apiMachineryCli apimachinery.ClientSetInterface) error {
	coreExist := false
	for retry := 0; retry < 10; retry++ {
		if _, err := apiMachineryCli.Healthz().HealthCheck(commontypes.CC_MODULE_CORESERVICE); err != nil {
			blog.Errorf("connect core server failed: %v", err)
			time.Sleep(time.Second * 2)
			continue
		}
		coreExist = true
		break
	}
	if !coreExist {
		blog.Errorf("core server not exist")
		return fmt.Errorf("core server not exist")
	}
	err := Init(&Options{ApiMachineryCli: apiMachineryCli})
	if err != nil {
		return err
	}
	return nil
}

// ValidatePlatformTenantMode validate platform multi-tenant mode
func ValidatePlatformTenantMode(tenantID string, enableTenantMode bool) bool {
	if enableTenantMode && tenantID == common.BKDefaultTenantID {
		return true
	}

	if !enableTenantMode && tenantID == common.BKSingleTenantID {
		return true
	}

	return false
}

// Options is tenant initialize options
type Options struct {
	DB              local.DB
	ApiMachineryCli apimachinery.ClientSetInterface
}

// Init initialize tenant info
func Init(opts *Options) error {
	if opts == nil {
		return fmt.Errorf("options is invalid")
	}
	if ((opts.DB != nil) && (opts.ApiMachineryCli != nil)) || ((opts.DB == nil) && (opts.ApiMachineryCli == nil)) {
		return fmt.Errorf("options is invalid, db: %v, api machinery: %v", opts.DB, opts.ApiMachineryCli)
	}
	db = opts.DB
	apiMachineryCli = opts.ApiMachineryCli
	var err error
	once.Do(func() {
		if err = refreshTenantInfo(); err != nil {
			return
		}

		// loop refresh tenant info
		go func() {
			for {
				time.Sleep(time.Minute)
				if err := refreshTenantInfo(); err != nil {
					blog.Errorf("refresh tenant info failed, err: %v", err)
					continue
				}
			}
		}()
	})
	if err != nil {
		return fmt.Errorf("init tenant info failed, err: %v", err)
	}

	return nil
}

func refreshTenantInfo() error {
	var tenants []types.Tenant
	var err error

	if db != nil {
		tenants, err = GetAllTenantsFromDB(context.Background(), db)
		if err != nil {
			blog.Errorf("get all tenants from db failed, err: %v", err)
			return err
		}
	}
	if apiMachineryCli != nil {
		tenants, err = apiMachineryCli.CoreService().Tenant().GetAllTenants(context.Background(),
			util.GenDefaultHeader())
		if err != nil {
			blog.Errorf("get all tenants from api machinery failed, err: %v", err)
			return err
		}
	}

	tenant.SetTenant(tenants)
	return nil
}

// GetAllTenantsFromDB get all tenants from db
func GetAllTenantsFromDB(ctx context.Context, db local.DB) ([]types.Tenant, error) {
	tenants := make([]types.Tenant, 0)
	err := db.Table(common.BKTableNameTenant).Find(nil).All(ctx, &tenants)
	if err != nil {
		return nil, fmt.Errorf("get all tenants from db failed, err: %v", err)
	}
	return tenants, nil
}
