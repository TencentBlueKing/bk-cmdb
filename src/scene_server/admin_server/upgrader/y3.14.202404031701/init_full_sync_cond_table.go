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

package y3_14_202404031701

import (
	"context"

	fullsynccond "configcenter/pkg/cache/full-sync-cond"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func initFullSyncCondTable(ctx context.Context, db dal.RDB) error {
	table := fullsynccond.BKTableNameFullSyncCond

	exists, err := db.HasTable(ctx, table)
	if err != nil {
		blog.Errorf("check if full sync cond table exists failed, err: %v", err)
		return err
	}

	if exists {
		return nil
	}

	if err = db.CreateTable(ctx, table); err != nil {
		blog.Errorf("create full sync cond table failed, err: %v", err)
		return err
	}

	return nil
}

func initFullSyncCondIndex(ctx context.Context, db dal.RDB) error {
	table := fullsynccond.BKTableNameFullSyncCond

	indexes := []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + fullsynccond.IDField,
			Keys: bson.D{
				{fullsynccond.IDField, 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "resource_subResource",
			Keys: bson.D{
				{fullsynccond.ResourceField, 1},
				{fullsynccond.SubResField, 1},
				{common.BkSupplierAccount, 1},
			},
			PartialFilterExpression: map[string]interface{}{
				fullsynccond.IsAllField: true,
			},
			Background: true,
			Unique:     true,
		},
	}

	existIndexes, err := db.Table(table).Indexes(ctx)
	if err != nil {
		blog.Errorf("get full sync cond table index failed, err: %v", err)
		return err
	}

	existIndexMap := make(map[string]struct{})
	for _, index := range existIndexes {
		existIndexMap[index.Name] = struct{}{}
	}

	for _, index := range indexes {
		if _, exist := existIndexMap[index.Name]; exist {
			continue
		}

		err = db.Table(table).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create full sync cond table index %+v failed, err: %v", index, err)
			return err
		}
	}

	return nil
}
