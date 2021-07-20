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

package y3_9_202104211151

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

var (
	sortFlag      = int32(1)
	idUniqueIndex = types.Index{
		Keys:       map[string]int32{common.BKFieldID: sortFlag},
		Unique:     true,
		Background: true,
		Name:       "idx_unique_id",
	}
)

func changeSetUniqueIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	idxUniqueParentIDSetName := types.Index{
		Keys:       map[string]int32{common.BKParentIDField: sortFlag, common.BKSetNameField: sortFlag},
		Unique:     true,
		Background: true,
		Name:       "idx_unique_parentID_setName",
	}
	tableName := common.BKTableNameBaseSet
	dbIndexes, err := db.Table(tableName).Indexes(ctx)
	if err != nil {
		blog.ErrorJSON("find table(%s) index error. err: %s", tableName, err.Error())
		return err
	}
	isCreateUniqueSetName := true
	for _, index := range dbIndexes {

		if index.Keys == nil {
			continue
		}

		if len(index.Keys) == 2 {
			_, ok1 := index.Keys[common.BKParentIDField]
			_, ok2 := index.Keys[common.BKSetNameField]
			if ok1 && ok2 {
				if index.Unique == true {
					isCreateUniqueSetName = false
				} else {
					if err := db.Table(tableName).DropIndex(ctx, index.Name); err != nil {
						blog.ErrorJSON("drop table(%s) index error. idx name: %s, err: %s",
							tableName, index.Name, err.Error())
						return err
					}
				}
				continue
			}
		}

		if index.Name == "idx_unique_bizID_setName" {
			if err := db.Table(tableName).DropIndex(ctx, "idx_unique_bizID_setName"); err != nil {
				blog.ErrorJSON("drop table(%s) index error. idx name: %s, err: %s",
					tableName, "idx_unique_bizID_setName", err.Error())
				return err
			}
		}
	}
	if isCreateUniqueSetName {
		if err := db.Table(tableName).CreateIndex(ctx, idxUniqueParentIDSetName); err != nil {
			blog.ErrorJSON("create table(%s) index error. idx name: %s, index: %s, err: %s",
				tableName, idxUniqueParentIDSetName.Name, idxUniqueParentIDSetName, err.Error())
			if db.IsDuplicatedError(err) {
				return nil
			}
			return err
		}
	}
	return nil
}
