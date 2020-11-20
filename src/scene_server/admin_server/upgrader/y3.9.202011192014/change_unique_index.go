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

package y3_9_202011192014

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

func changeUniqueIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	idxUniqueObjIDGroupName := types.Index{
		Keys: map[string]int32{common.BKObjIDField: sortFlag, common.BKAppIDField: sortFlag,
			common.BKPropertyGroupNameField: sortFlag},
		Unique:     true,
		Background: true,
		Name:       "idx_unique_objID_groupName",
	}

	tableName := common.BKTableNamePropertyGroup
	dbIndexes, err := db.Table(tableName).Indexes(ctx)
	if err != nil {
		blog.ErrorJSON("find table(%s) index error. err: %s", tableName, err.Error())
		return err
	}
	idxUniqueGroupName := "idx_unique_objID_groupName"

	isCreateUniqueGroupName := true
	for _, index := range dbIndexes {

		if index.Keys == nil {
			continue
		}
		// 已经存在业务ID不需要再操作了
		if _, ok := index.Keys[common.BKAppIDField]; ok {
			if index.Name == idxUniqueGroupName {
				isCreateUniqueGroupName = false
			}
			continue
		}

		switch index.Name {
		case idxUniqueGroupName:
			if err := db.Table(tableName).DropIndex(ctx, idxUniqueGroupName); err != nil {
				blog.ErrorJSON("drop table(%s) index error. idx name: %s, err: %s",
					tableName, idxUniqueGroupName, err.Error())
				return err
			}
			if err := db.Table(tableName).CreateIndex(ctx, idxUniqueObjIDGroupName); err != nil {
				blog.ErrorJSON("create table(%s) index error. idx name: %s, index: %s, err: %s",
					tableName, idxUniqueGroupName, idxUniqueObjIDGroupName, err.Error())
				return err
			}
		case "idx_unique_objID_groupIdx":
			if err := db.Table(tableName).DropIndex(ctx, "idx_unique_objID_groupIdx"); err != nil {
				blog.ErrorJSON("drop table(%s) index error. idx name: %s, err: %s",
					tableName, "idx_unique_objID_groupIdx", err.Error())
				return err
			}
		}

	}
	if isCreateUniqueGroupName {
		if err := db.Table(tableName).CreateIndex(ctx, idxUniqueObjIDGroupName); err != nil {
			blog.ErrorJSON("create table(%s) index error. idx name: %s, index: %s, err: %s",
				tableName, idxUniqueGroupName, idxUniqueObjIDGroupName, err.Error())
			if db.IsDuplicatedError(err) {
				return nil
			}
			return err
		}
	}

	return nil
}
