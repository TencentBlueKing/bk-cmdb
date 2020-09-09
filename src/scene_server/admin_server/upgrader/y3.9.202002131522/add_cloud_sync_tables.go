/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package y3_9_202002131522

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

var CreateTableOptions = []struct {
	TableName  string
	TableIndex []types.Index
}{
	{
		TableName: common.BKTableNameCloudSyncTask,
		TableIndex: []types.Index{
			{Name: "bk_task_id", Keys: map[string]int32{"bk_task_id": 1}, Unique: true, Background: true},
		},
	},
	{
		TableName: common.BKTableNameCloudAccount,
		TableIndex: []types.Index{
			{Name: "bk_account_id", Keys: map[string]int32{"bk_account_id": 1}, Background: true},
		},
	},
	{
		TableName: common.BKTableNameCloudSyncHistory,
		TableIndex: []types.Index{
			{Name: "bk_history_id", Keys: map[string]int32{"bk_history_id": 1}, Background: true},
		},
	},
}

func upsertTable(ctx context.Context, db dal.RDB, conf *upgrader.Config, tableName string, indices []types.Index) error {
	exists, err := db.HasTable(ctx, tableName)
	if err != nil {
		return fmt.Errorf("check HasTable failed, tableName: %s, err: %+v", tableName, err)
	}
	if exists == false {
		if err = db.CreateTable(ctx, tableName); err != nil && !db.IsDuplicatedError(err) {
			return fmt.Errorf("CreateTable failed, tableName: %s, err: %+v", tableName, err)
		}
	}

	existIndices, err := db.Table(tableName).Indexes(ctx)
	if err != nil {
		return fmt.Errorf("upsertTable failed, Indexes failed, tableName: %s, err:%+v", tableName, err)
	}
	existIdxMap := make(map[string]bool)
	for _, idx := range existIndices {
		existIdxMap[idx.Name] = true
	}
	for _, index := range indices {
		if _, ok := existIdxMap[index.Name]; ok == true {
			continue
		}
		if err = db.Table(tableName).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			return fmt.Errorf("CreateIndex failed, tableName: %s, err:%+v", tableName, err)
		}
	}
	return nil
}

func CreateTables(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	for _, item := range CreateTableOptions {
		err := upsertTable(ctx, db, conf, item.TableName, item.TableIndex)
		if err != nil {
			return fmt.Errorf("upsertTable failed, err: %s", err.Error())
		}
	}
	return nil
}
