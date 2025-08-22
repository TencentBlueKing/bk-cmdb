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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/types"

	"github.com/rs/xid"
)

// inPlaceUpgrade in place upgrade data to multi-tenant version
func (s *migrateTenantService) inPlaceUpgrade(kit *rest.Kit, skipRemoveSupplierAccount bool) error {
	// upgrade object related table data
	if err := s.inPlaceUpgradeObject(kit, skipRemoveSupplierAccount); err != nil {
		return err
	}

	// upgrade table name and remove supplier account field for all tables
	for table := range s.tableHandlers {
		switch table {
		case common.BKTableNameBasePlat, common.BKTableNameBaseHost:
			// convert unassigned cloud id to -1
			fmt.Println("=================================")
			printInfo("start upgrade unassigned cloud id\n")
			cond := mapstr.MapStr{common.BKCloudIDField: 90000001}
			updateData := types.ModeUpdate{
				Op:  "set",
				Doc: mapstr.MapStr{common.BKCloudIDField: -1},
			}
			if err := upgradeTableData(kit, s.oldDB, table, cond, updateData); err != nil {
				return err
			}
		case common.BKTableNameObjDes:
			continue
		}

		if err := s.inPlaceUpgradeTable(kit, table, table, skipRemoveSupplierAccount); err != nil {
			return err
		}
	}

	// delete supplier account object attribute
	objAttrCond := mapstr.MapStr{common.BKPropertyIDField: "bk_supplier_account"}
	if err := s.newDB.Table(common.BKTableNameObjDes).Delete(kit.Ctx, objAttrCond); err != nil {
		return fmt.Errorf("delete supplier account object attribute failed, err: %v", err)
	}

	// copy id generator data and remove old id generator table
	if err := s.copyIDGenerator(kit); err != nil {
		return err
	}

	// upgrade watch info
	if err := s.inPlaceUpgradeWatchInfo(kit, skipRemoveSupplierAccount); err != nil {
		return err
	}

	// copy system data and remove old system table
	if err := s.copySystemData(kit); err != nil {
		return err
	}

	// drop deprecated tables
	fmt.Println("=================================")
	tables, err := s.oldDB.ListTables(kit.Ctx)
	if err != nil {
		return fmt.Errorf("list tables failed, err: %v", err)
	}

	for _, table := range tables {
		if !strings.HasPrefix(table, "cc_") {
			continue
		}
		printInfo("drop %s table\n", table)
		if err = s.oldDB.DropTable(kit.Ctx, table); err != nil {
			return fmt.Errorf("drop deprecated table %s failed, err: %v", table, err)
		}
	}

	return nil
}

// upgradeTableData upgrade table data
func upgradeTableData(kit *rest.Kit, db local.DB, table string, cond mapstr.MapStr,
	updateData ...types.ModeUpdate) error {

	for {
		data := make([]mapstr.MapStr, 0)
		err := db.Table(table).Find(cond).Fields("_id").Limit(common.BKMaxPageSize).All(kit.Ctx, &data)
		if err != nil {
			return fmt.Errorf("get table %s data oids by cond(%+v) failed, err: %v", table, cond, err)
		}

		if len(data) == 0 {
			break
		}

		oids := make([]interface{}, len(data))
		for i, oidData := range data {
			oids[i] = oidData["_id"]
		}

		oidCond := mapstr.MapStr{"_id": mapstr.MapStr{common.BKDBIN: oids}}
		if err = db.Table(table).UpdateMultiModel(kit.Ctx, oidCond, updateData...); err != nil {
			return fmt.Errorf("update table %s data(%+v) failed, err: %v", table, updateData, err)
		}
	}
	return nil
}

// inPlaceUpgradeObject in place upgrade object related tables
func (s *migrateTenantService) inPlaceUpgradeObject(kit *rest.Kit, skipRemoveSupplierAccount bool) error {
	fmt.Println("=================================")
	printInfo("start upgrade object table\n")

	// rename object table and upgrade object data
	exists, err := s.oldDB.HasTable(kit.Ctx, common.BKTableNameObjDes)
	if err != nil {
		return fmt.Errorf("check if table %s exists failed, err: %v", common.BKTableNameObjDes, err)
	}
	if exists {
		err := s.oldDB.RenameTable(kit.Ctx, "cc_"+common.BKTableNameObjDes, "default_"+common.BKTableNameObjDes)
		if err != nil {
			return fmt.Errorf("rename table %s failed, err: %v", common.BKTableNameObjDes, err)
		}
	}

	objects := make([]metadata.Object, 0)
	if err = s.newDB.Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField, metadata.ModelFieldObjUUID).
		All(kit.Ctx, &objects); err != nil {
		return fmt.Errorf("get object id and uuid info failed, err: %v", err)
	}

	for _, object := range objects {
		if object.UUID == "" {
			object.UUID = xid.New().String()

			cond := mapstr.MapStr{common.BKObjIDField: object.ObjectID}
			data := []types.ModeUpdate{{Op: "set", Doc: mapstr.MapStr{metadata.ModelFieldObjUUID: object.UUID}},
				{Op: "unset", Doc: mapstr.MapStr{"bk_supplier_account": ""}}}
			if err = s.newDB.Table(common.BKTableNameObjDes).UpdateMultiModel(kit.Ctx, cond, data...); err != nil {
				return fmt.Errorf("update object %s uuid failed, err: %v", object.ObjectID, err)
			}
		}
		s.objUUIDMap[object.ObjectID] = object.UUID
	}

	// upgrade object instance and instance association data
	for objID, uuid := range s.objUUIDMap {
		if !common.IsInnerModel(objID) {
			oldTable := fmt.Sprintf("%s0_pub_%s", common.BKObjectInstShardingTablePrefix, objID)
			newTable := common.GetObjInstTableName(uuid)
			if err = s.inPlaceUpgradeTable(kit, oldTable, newTable, skipRemoveSupplierAccount); err != nil {
				return err
			}
		}

		oldAsstTable := fmt.Sprintf("%s0_pub_%s", common.BKObjectInstAsstShardingTablePrefix, objID)
		newAsstTable := common.GetObjInstAsstTableName(uuid)
		if err = s.inPlaceUpgradeTable(kit, oldAsstTable, newAsstTable, skipRemoveSupplierAccount); err != nil {
			return err
		}
	}

	return nil
}

// inPlaceUpgradeTable upgrade table name and remove supplier account field
func (s *migrateTenantDBInfo) inPlaceUpgradeTable(kit *rest.Kit, oldTable, newTable string,
	skipRemoveSupplierAccount bool) error {

	fmt.Println("=================================")
	printInfo("start upgrade table %s\n", newTable)

	renamedTable := newTable
	db := s.sysDB
	if !common.IsPlatformTable(newTable) && !common.IsPlatformTableWithTenant(newTable) {
		renamedTable = "default_" + newTable
		db = s.newDB
	}

	exists, err := s.oldDB.HasTable(kit.Ctx, oldTable)
	if err != nil {
		return fmt.Errorf("check if table %s exists failed, err: %v", oldTable, err)
	}
	if exists {
		if err = s.oldDB.RenameTable(kit.Ctx, "cc_"+oldTable, renamedTable); err != nil {
			return fmt.Errorf("rename table %s failed, err: %v", newTable, err)
		}
	}

	if skipRemoveSupplierAccount && !common.IsPlatformTableWithTenant(newTable) {
		return nil
	}

	updateData := []types.ModeUpdate{{
		Op:  "unset",
		Doc: mapstr.MapStr{"bk_supplier_account": ""},
	}}
	if common.IsPlatformTableWithTenant(newTable) {
		updateData = append(updateData, types.ModeUpdate{
			Op:  "set",
			Doc: mapstr.MapStr{common.TenantID: common.BKSingleTenantID},
		})
	}
	updateCond := mapstr.MapStr{"bk_supplier_account": mapstr.MapStr{common.BKDBExists: true}}

	return upgradeTableData(kit, db, newTable, updateCond, updateData...)
}
