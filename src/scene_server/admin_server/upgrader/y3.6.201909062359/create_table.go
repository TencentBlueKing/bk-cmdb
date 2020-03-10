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

package y3_6_201909062359

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

func createSetTemplateTables(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tables := []string{
		common.BKTableNameSetTemplate,
		common.BKTableNameSetServiceTemplateRelation,
	}

	for _, tableName := range tables {
		exists, err := db.HasTable(ctx, tableName)
		if err != nil {
			return err
		}
		if exists == true {
			continue
		}

		// add table
		if err = db.CreateTable(ctx, tableName); err != nil && !db.IsDuplicatedError(err) {
			return err
		}
	}

	// list indexes
	indices, err := db.Table(common.BKTableNameSetTemplate).Indexes(ctx)
	if err != nil {
		blog.ErrorJSON("list index for table: %s failed, err:%s", common.BKTableNameSetTemplate, err.Error())
		return err
	}
	idxMap := make(map[string]bool)
	for _, idx := range indices {
		idxMap[idx.Name] = true
	}

	// check index exist and add
	idxName := "idx_id"
	if _, ok := idxMap[idxName]; ok == false {
		index := types.Index{
			Keys:       map[string]int32{common.BKFieldID: 1},
			Name:       idxName,
			Unique:     true,
			Background: true,
		}
		if err := db.Table(common.BKTableNameSetTemplate).CreateIndex(ctx, index); err != nil {
			blog.ErrorJSON("add index for table: %s failed, index: %s, err:%s", common.BKTableNameSetTemplate, index, err.Error())
			return err
		}
	}
	return nil
}
