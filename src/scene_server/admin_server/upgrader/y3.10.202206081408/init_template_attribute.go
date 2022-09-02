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

package y3_10_202206081408

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func initTemplateAttribute(ctx context.Context, db dal.RDB, tableName, templateIDField string) error {
	// create table if it isn't exist
	hasTable, err := db.HasTable(ctx, tableName)
	if err != nil {
		blog.Errorf("check if %s table exist failed, err: %v", tableName, err)
		return err
	}
	if hasTable {
		return nil
	}

	if err := db.CreateTable(ctx, tableName); err != nil {
		blog.Errorf("create %s table failed, err: %v", tableName, err)
		return err
	}

	// add index if it isn't exist
	indexes := []types.Index{
		{
			Keys: bson.D{
				{
					common.BKFieldID, 1,
				},
			},
			Name:       common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
			Unique:     true,
			Background: true,
		},
		{
			Keys: bson.D{
				{
					common.BKAppIDField, 1,
				},
				{
					templateIDField, 1,
				},
				{
					common.BKAttributeIDField, 1,
				},
			},
			Name: common.CCLogicUniqueIdxNamePrefix + common.BKAppIDField + "_" + templateIDField + "_" +
				common.BKAttributeIDField,
			Background: true,
			Unique:     true,
		},
	}

	existIndexArr, err := db.Table(tableName).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for %s table failed, err: %v", tableName, err)
		return err
	}

	existIdxMap := make(map[string]struct{})
	for _, index := range existIndexArr {
		existIdxMap[index.Name] = struct{}{}
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}

		err = db.Table(tableName).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create %s table index(%+v) failed, err: %v", tableName, index, err)
			return err
		}
	}
	return nil
}
