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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	dbtypes "configcenter/src/storage/dal/types"
)

// copyWatchInfo copy watch related data from old db to new db
func (s *migrateTenantService) copyWatchInfo(kit *rest.Kit) error {
	// insert watch db relation info
	watchDBRel := mapstr.MapStr{
		"db":       s.dbUUID,
		"watch_db": s.watchDB.dbUUID,
	}
	err := s.watchDB.sysDB.Table(common.BKTableNameWatchDBRelation).Upsert(kit.Ctx, watchDBRel, watchDBRel)
	if err != nil {
		return fmt.Errorf("insert watch db relation info failed, err: %v", err)
	}

	// copy watch chain node data to new db
	tables, err := s.watchDB.oldDB.ListTables(kit.Ctx)
	if err != nil {
		return fmt.Errorf("list watch db tables failed, err: %v", err)
	}

	for _, table := range tables {
		if !strings.HasSuffix(table, "WatchChain") {
			continue
		}

		table = strings.TrimPrefix(table, "cc_")
		if err = s.watchDB.copyTableData(kit, table, table, s.removeSupplierAccount); err != nil {
			return err
		}
	}

	// copy id generator data and watch token data from old db to new db
	if err = s.watchDB.copyIDGenerator(kit); err != nil {
		return err
	}

	if err = s.copyWatchToken(kit); err != nil {
		return err
	}

	cond := mapstr.MapStr{"_id": s.dbUUID}
	token := mapstr.MapStr{
		"_id":                     s.dbUUID,
		common.BKTokenField:       "",
		common.BKStartAtTimeField: time.Now(),
	}
	if err = s.watchDB.sysDB.Table(common.BKTableNameWatchToken).Upsert(kit.Ctx, cond, token); err != nil {
		return fmt.Errorf("insert db watch token info failed, err: %v", err)
	}

	return nil
}

// inPlaceUpgradeWatchInfo in place upgrade watch related tables
func (s *migrateTenantService) inPlaceUpgradeWatchInfo(kit *rest.Kit, skipRemoveSupplierAccount bool) error {
	// insert watch db relation info
	watchDBRel := mapstr.MapStr{
		"db":       s.dbUUID,
		"watch_db": s.watchDB.dbUUID,
	}
	err := s.watchDB.sysDB.Table(common.BKTableNameWatchDBRelation).Upsert(kit.Ctx, watchDBRel, watchDBRel)
	if err != nil {
		return fmt.Errorf("insert watch db relation info failed, err: %v", err)
	}

	// in place upgrade watch chain node tables
	tables, err := s.watchDB.oldDB.ListTables(kit.Ctx)
	if err != nil {
		return fmt.Errorf("list watch db tables failed, err: %v", err)
	}

	for _, table := range tables {
		if !strings.HasSuffix(table, "WatchChain") {
			continue
		}

		table = strings.TrimPrefix(table, "cc_")
		if err = s.watchDB.inPlaceUpgradeTable(kit, table, table, skipRemoveSupplierAccount); err != nil {
			return err
		}
	}

	// copy id generator data from old db to new db
	if err = s.watchDB.copyIDGenerator(kit); err != nil {
		return err
	}

	if err = s.watchDB.oldDB.DropTable(kit.Ctx, common.BKTableNameIDgenerator); err != nil {
		return fmt.Errorf("remove old id generator table failed, err: %v", err)
	}

	// copy watch token data from old db to new db and remove old table
	if err = s.copyWatchToken(kit); err != nil {
		return err
	}

	cond := mapstr.MapStr{"_id": s.dbUUID}
	token := mapstr.MapStr{
		"_id":                     s.dbUUID,
		common.BKTokenField:       "",
		common.BKStartAtTimeField: time.Now(),
	}
	if err = s.watchDB.sysDB.Table(common.BKTableNameWatchToken).Upsert(kit.Ctx, cond, token); err != nil {
		return fmt.Errorf("insert db watch token info failed, err: %v", err)
	}

	if err = s.watchDB.oldDB.DropTable(kit.Ctx, common.BKTableNameWatchToken); err != nil {
		return fmt.Errorf("remove old watch token table failed, err: %v", err)
	}

	return nil
}

// inPlaceUpgradeWatchToken copy watch token data from old db to new db
func (s *migrateTenantService) copyWatchToken(kit *rest.Kit) error {
	fmt.Println("=================================")
	printInfo("start copy watch token data\n")

	tokens := make([]mapstr.MapStr, 0)
	if err := s.watchDB.oldDB.Table(common.BKTableNameWatchToken).Find(nil,
		dbtypes.NewFindOpts().SetWithObjectID(true)).All(kit.Ctx, &tokens); err != nil {
		return fmt.Errorf("get watch tokens failed, err: %v", err)
	}

	for _, token := range tokens {
		id := strings.TrimPrefix(util.GetStrByInterface(token["_id"]), "cc_")

		cond := mapstr.MapStr{"_id": id}
		lastEvent := mapstr.MapStr{
			"_id":                id,
			common.BKFieldID:     token[common.BKFieldID],
			common.BKCursorField: token[common.BKCursorField],
		}
		if err := s.watchDB.newDB.Table(common.BKTableNameLastWatchEvent).Upsert(kit.Ctx, cond, lastEvent); err != nil {
			return fmt.Errorf("insert id generators failed, err: %v", err)
		}

		delete(token, common.BKFieldID)
		delete(token, common.BKCursorField)
		id = s.dbUUID + ":" + id
		cond = mapstr.MapStr{"_id": id}
		token["_id"] = id
		for key, value := range token {
			if strings.HasPrefix(key, "cc_") {
				token[strings.TrimPrefix(key, "cc_")] = value
				delete(token, key)
			}
		}
		if err := s.watchDB.sysDB.Table(common.BKTableNameWatchToken).Upsert(kit.Ctx, cond, token); err != nil {
			return fmt.Errorf("insert id generators failed, err: %v", err)
		}
	}

	return nil
}
