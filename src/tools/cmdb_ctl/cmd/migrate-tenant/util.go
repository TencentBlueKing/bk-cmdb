/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package migratetenant

import (
	"fmt"

	"configcenter/pkg/tenant/types"
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/rs/xid"
)

// removeSupplierAccount remove supplier account field
func (s *migrateTenantService) removeSupplierAccount(data mapstr.MapStr) (mapstr.MapStr, error) {
	delete(data, "bk_supplier_account")
	return data, nil
}

// migrateTenantIDField remove supplier account field and add default tenant id field
func (s *migrateTenantService) migrateTenantIDField(data mapstr.MapStr) (mapstr.MapStr, error) {
	delete(data, "bk_supplier_account")
	data[common.TenantID] = common.BKSingleTenantID
	return data, nil
}

// migrateCloudIDField convert unassigned cloud id to -1
func (s *migrateTenantService) migrateCloudIDField(data mapstr.MapStr) (mapstr.MapStr, error) {
	delete(data, "bk_supplier_account")

	rawCloudID, exists := data[common.BKCloudIDField]
	if !exists {
		return data, nil
	}

	cloudID, err := util.GetInt64ByInterface(rawCloudID)
	if err != nil {
		return nil, fmt.Errorf("invalid cloud id: %v, err: %v, data: %+v", rawCloudID, err, data)
	}

	if cloudID == 90000001 {
		data[common.BKCloudIDField] = -1
	}
	return data, nil
}

// migrateObject add uuid field to object
func (s *migrateTenantService) migrateObject(data mapstr.MapStr) (mapstr.MapStr, error) {
	delete(data, "bk_supplier_account")

	objID := util.GetStrByInterface(data[common.BKObjIDField])

	uuid, exists := s.objUUIDMap[objID]
	if exists {
		data[metadata.ModelFieldObjUUID] = uuid
		return data, nil
	}

	uuid = util.GetStrByInterface(data[metadata.ModelFieldObjUUID])
	if uuid == "" {
		uuid = xid.New().String()
	}
	s.objUUIDMap[objID] = uuid

	data[metadata.ModelFieldObjUUID] = uuid
	return data, nil
}

// insertTenant insert default tenant info
func (s *migrateTenantService) insertTenant(kit *rest.Kit) error {
	fmt.Println("=================================")
	printInfo("start insert default tenant info\n")

	filter := map[string]interface{}{
		common.TenantID: common.BKSingleTenantID,
	}
	err := s.sysDB.Table(common.BKTableNameTenant).Upsert(kit.Ctx, filter, types.Tenant{
		TenantID: common.BKSingleTenantID,
		Status:   types.EnabledStatus,
		Database: s.dbUUID,
	})
	if err != nil {
		return fmt.Errorf("insert default tenant info failed, err: %v", err)
	}
	return nil
}

func printInfo(format string, a ...interface{}) {
	fmt.Printf("INFO: %s", fmt.Sprintf(format, a...))
}
