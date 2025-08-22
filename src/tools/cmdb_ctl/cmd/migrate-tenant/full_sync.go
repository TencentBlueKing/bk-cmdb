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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/index"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/logics"
	dbtypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
)

// copyFullSyncData copy full sync data from old db to new db
func (s *migrateTenantService) copyFullSyncData(kit *rest.Kit) error {
	// record full sync start time for incremental sync
	startTimeCond := mapstr.MapStr{"_id": "full_sync_start_time"}
	count, err := s.sysDB.Table(common.BKTableNameSystem).Find(startTimeCond).Count(kit.Ctx)
	if err != nil {
		return fmt.Errorf("check if full sync start time exists failed, err: %v", err)
	}
	if count == 0 {
		startTime := mapstr.MapStr{
			"_id":        "full_sync_start_time",
			"start_from": time.Now().Unix(),
		}
		if err = s.sysDB.Table(common.BKTableNameSystem).Insert(kit.Ctx, startTime); err != nil {
			return fmt.Errorf("insert full sync start time info failed, err: %v", err)
		}
	}

	// copy data from old db to new db for tables in tableHandlers
	for table, handler := range s.tableHandlers {
		if err := s.copyTableData(kit, table, table, handler); err != nil {
			return err
		}
	}

	// copy data from old db to new db for object instance and instance association tables
	for objID, uuid := range s.objUUIDMap {
		if !common.IsInnerModel(objID) {
			oldTable := fmt.Sprintf("%s0_pub_%s", common.BKObjectInstShardingTablePrefix, objID)
			newTable := common.GetObjInstTableName(uuid)
			if err := s.copyTableData(kit, oldTable, newTable, s.removeSupplierAccount); err != nil {
				return err
			}
		}

		oldAsstTable := fmt.Sprintf("%s0_pub_%s", common.BKObjectInstAsstShardingTablePrefix, objID)
		newAsstTable := common.GetObjInstAsstTableName(uuid)
		if err := s.copyTableData(kit, oldAsstTable, newAsstTable, s.removeSupplierAccount); err != nil {
			return err
		}
	}

	// copy id generator data from old db to new db
	if err = s.copyIDGenerator(kit); err != nil {
		return err
	}

	// copy system data from old db to new db
	if err = s.copySystemData(kit); err != nil {
		return err
	}

	// copy watch data from old db to new db
	if err = s.copyWatchInfo(kit); err != nil {
		return err
	}

	return nil
}

// copyTableData copy data from old db to new db for specified table
func (s *migrateTenantDBInfo) copyTableData(kit *rest.Kit, oldTable, newTable string,
	handler func(data mapstr.MapStr) (mapstr.MapStr, error)) error {

	db := s.newDB
	if common.IsPlatformTable(newTable) {
		db = s.sysDB
	}

	opts := dbtypes.NewFindOpts().SetWithObjectID(true)

	fmt.Println("=================================")
	printInfo("start creating table %s for new db\n", newTable)

	// create table and indexes
	if err := logics.CreateTable(kit, db, newTable); err != nil {
		return fmt.Errorf("create table %s failed, err: %v", newTable, err)
	}

	indexes := index.TableIndexes()[newTable]
	if len(indexes) > 0 {
		if err := logics.CreateIndexes(kit, db, newTable, indexes); err != nil {
			return fmt.Errorf("create table %s indexes failed, err: %v", newTable, err)
		}
	}

	// get start oid from new db, ignores the synced data
	startOidInfo := make(mapstr.MapStr)
	err := db.Table(newTable).Find(nil).Fields("_id").Sort("-_id").One(kit.Ctx, &startOidInfo)
	if err != nil && !mongodb.IsNotFoundError(err) {
		return fmt.Errorf("get table %s start oid failed, err: %v", newTable, err)
	}

	cond := make(mapstr.MapStr)
	startOid, exists := startOidInfo[common.MongoMetaID]
	if exists {
		cond[common.MongoMetaID] = mapstr.MapStr{common.BKDBGT: startOid}
	}

	for {
		printInfo("start sync table %s data, cond: %+v\n", newTable, cond)

		dataArr := make([]mapstr.MapStr, 0)
		err = s.oldDB.Table(oldTable).Find(cond, opts).Sort("_id").Limit(common.BKMaxPageSize).All(kit.Ctx, &dataArr)
		if err != nil {
			return fmt.Errorf("get %s table data failed, cond: %+v, err: %v", oldTable, cond, err)
		}

		if len(dataArr) == 0 {
			break
		}

		insertData := make([]mapstr.MapStr, 0)
		for _, data := range dataArr {
			if newTable == common.BKTableNameObjAttDes && data[common.BKPropertyIDField] == "bk_supplier_account" {
				continue
			}

			newData, err := handler(data)
			if err != nil {
				return fmt.Errorf("convert %s data(%+v) to new version failed, err: %v", newTable, data, err)
			}
			insertData = append(insertData, newData)
		}

		if len(insertData) > 0 {
			if err = db.Table(newTable).Insert(kit.Ctx, insertData); err != nil {
				return fmt.Errorf("insert %s table data(%+v) failed, err: %v", newTable, insertData, err)
			}
		}

		cond = mapstr.MapStr{
			common.MongoMetaID: mapstr.MapStr{common.BKDBGT: dataArr[len(dataArr)-1][common.MongoMetaID]},
		}
	}
	return nil
}

// copyIDGenerator copy id generator data with new oid and remove the old ones
func (s *migrateTenantDBInfo) copyIDGenerator(kit *rest.Kit) error {
	fmt.Println("=================================")
	printInfo("start copy id generator data\n")

	idGenerators := make([]mapstr.MapStr, 0)
	if err := s.oldDB.Table(common.BKTableNameIDgenerator).Find(nil, dbtypes.NewFindOpts().SetWithObjectID(true)).
		All(kit.Ctx, &idGenerators); err != nil {
		return fmt.Errorf("get id generators failed, err: %v", err)
	}

	for _, idGenerator := range idGenerators {
		id := strings.TrimPrefix(util.GetStrByInterface(idGenerator["_id"]), "cc_")
		if strings.HasPrefix(id, "id_rule:incr_id:") {
			id = "id_rule:incr_id:default:" + strings.TrimPrefix(id, "id_rule:incr_id:")
		}

		cond := mapstr.MapStr{"_id": id}
		idGenerator["_id"] = id
		if err := s.sysDB.Table(common.BKTableNameIDgenerator).Upsert(kit.Ctx, cond, idGenerator); err != nil {
			return fmt.Errorf("insert id generators failed, err: %v", err)
		}
	}

	return nil
}

// copySystemData copy system data to new db table
func (s *migrateTenantService) copySystemData(kit *rest.Kit) error {
	fmt.Println("=================================")
	printInfo("start copying system data\n")

	// copy config admin info
	if err := s.copyConfigAdmin(kit); err != nil {
		return err
	}

	// copy version info
	versionCond := mapstr.MapStr{
		"type": "version",
	}
	versionInfo := make(mapstr.MapStr)
	err := s.oldDB.Table(common.BKTableNameSystem).Find(versionCond).One(kit.Ctx, &versionInfo)
	if err != nil {
		if mongodb.IsNotFoundError(err) {
			return nil
		}
		return fmt.Errorf("get system version info failed, err: %v", err)
	}

	versionInfo["current_version"] = "y3.15.202411071530"
	if err = s.sysDB.Table(common.BKTableNameSystem).Upsert(kit.Ctx, versionCond, versionInfo); err != nil {
		return fmt.Errorf("upsert system version info(%+v) failed, err: %v", versionInfo, err)
	}

	return nil
}

// copyConfigAdmin copy config admin info to new db
func (s *migrateTenantService) copyConfigAdmin(kit *rest.Kit) error {
	cond := mapstr.MapStr{common.BKFieldDBID: common.ConfigAdminID}
	info := make(map[string]string)
	err := s.oldDB.Table(common.BKTableNameSystem).Find(cond).Fields(common.ConfigAdminValueField).One(kit.Ctx, &info)
	if err != nil {
		if mongodb.IsNotFoundError(err) {
			return nil
		}
		return fmt.Errorf("get config admin info failed, err: %v", err)
	}

	configAdmin := new(metadata.OldPlatformSettingConfig)
	if err = json.Unmarshal([]byte(info[common.ConfigAdminValueField]), configAdmin); err != nil {
		return fmt.Errorf("parse config admin info %s failed, err: %v", info[common.ConfigAdminValueField], err)
	}

	platConfCond := map[string]interface{}{
		common.BKFieldDBID: common.PlatformConfig,
	}
	platConf := mapstr.MapStr{
		common.BKFieldDBID:         common.PlatformConfig,
		metadata.IDGeneratorConfig: configAdmin.IDGenerator,
	}
	if err = s.sysDB.Table(common.BKTableNameSystem).Upsert(kit.Ctx, platConfCond, platConf); err != nil {
		return fmt.Errorf("upsert platform config(%+v) failed, err: %v", platConf, err)
	}

	globalConfCond := mapstr.MapStr{common.TenantID: common.BKSingleTenantID}
	globalConf := metadata.GlobalSettingConfig{
		TenantID:            common.BKSingleTenantID,
		Backend:             metadata.AdminBackendCfg{MaxBizTopoLevel: configAdmin.Backend.MaxBizTopoLevel},
		ValidationRules:     configAdmin.ValidationRules,
		BuiltInSetName:      configAdmin.BuiltInSetName,
		BuiltInModuleConfig: configAdmin.BuiltInModuleConfig,
		CreateTime:          metadata.Now(),
		LastTime:            metadata.Now(),
	}
	if err = s.sysDB.Table(common.BKTableNameGlobalConfig).Upsert(kit.Ctx, globalConfCond, globalConf); err != nil {
		return fmt.Errorf("upsert global setting config(%+v) failed, err: %v", globalConf, err)
	}
	return nil
}
