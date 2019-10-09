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

package x19_10_09_01

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func SetTemplateSyncStatusMigrate(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableNames := []string{common.BKTableNameSetTemplateSyncStatus, common.BKTableNameSetTemplateSyncHistory}
	for _, tableName := range tableNames {
		existTable, err := db.HasTable(tableName)
		if err != nil {
			blog.Errorf("check table %s exist failed, err:%s", tableName, err.Error())
			return err
		}

		if !existTable {
			err := db.CreateTable(tableName)
			if err != nil {
				blog.Errorf("create table %s failed, err:%s", tableName, err.Error())
				return err
			}
		}

		setIDUnique := false
		if tableName == common.BKTableNameSetTemplateSyncStatus {
			setIDUnique = true
		}
		indexArr := []dal.Index{
			{
				Keys:       map[string]int32{"task_id": 1},
				Name:       "idx_taskID",
				Unique:     true,
				Background: true,
			},
			{
				Keys:       map[string]int32{"bk_set_id": 1},
				Name:       "idx_setID",
				Unique:     setIDUnique,
				Background: true,
			},
			{
				Keys:       map[string]int32{"last_time": 1, "create_time": 1},
				Name:       "idx_createLastTime",
				Background: true,
			},
		}

		for _, index := range indexArr {
			err := db.Table(tableName).CreateIndex(ctx, index)
			if err != nil {
				blog.ErrorJSON("add index %s for table %s failed, err:%s", index, tableName, err.Error())
				return err
			}
		}
	}

	return nil
}
