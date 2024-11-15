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
	"context"
	"fmt"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal/mongo/local"
)

var (
	allTenants = make([]Tenant, 0)
	db         local.DB
	once       sync.Once
)

// Options is tenant initialize options
type Options struct {
	DB local.DB
	// TODO add redis cache and api machinery client
}

// Init initialize tenant info
func Init(opts *Options) error {
	if opts == nil {
		return fmt.Errorf("options is nil")
	}

	if opts.DB == nil {
		return fmt.Errorf("db is nil")
	}
	db = opts.DB

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
	var err error
	if db != nil {
		allTenants, err = GetAllTenantsFromDB(context.Background(), db)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAllTenants get all tenants
func GetAllTenants() []Tenant {
	// TODO right now only support default tenant for compatible, use actual tenants later
	return []Tenant{{TenantID: common.BKDefaultTenantID, Status: EnabledStatus}}
}

// GetAllTenantsFromDB get all tenants from db
func GetAllTenantsFromDB(ctx context.Context, db local.DB) ([]Tenant, error) {
	tenants := make([]Tenant, 0)
	err := db.Table(common.BKTableNameTenant).Find(nil).All(ctx, &tenants)
	if err != nil {
		return nil, fmt.Errorf("get all tenants from db failed, err: %v", err)
	}
	return tenants, nil
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
