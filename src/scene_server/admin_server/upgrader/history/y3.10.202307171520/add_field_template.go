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

package y3_10_202307171520

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func addFieldTemplateCollection(ctx context.Context, db dal.RDB) error {
	collections := []string{
		common.BKTableNameFieldTemplate, common.BKTableNameObjAttDesTemplate,
		common.BKTableNameObjectUniqueTemplate, common.BKTableNameObjFieldTemplateRelation,
	}

	for _, collection := range collections {
		exists, err := db.HasTable(ctx, collection)
		if err != nil {
			blog.Errorf("check if %s table exists failed, err: %v", collection, err)
			return err
		}

		if exists {
			continue
		}

		if err := db.CreateTable(ctx, collection); err != nil {
			blog.Errorf("create %s table failed, err: %v", collection, err)
			return err
		}

	}
	return nil
}

func addFieldTemplateIndex(ctx context.Context, db dal.RDB) error {
	if err := addFieldTemplateIndexes(ctx, db); err != nil {
		return err
	}

	if err := addObjAttDesTemplateIndexes(ctx, db); err != nil {
		return err
	}

	if err := addObjectUniqueTemplateIndexes(ctx, db); err != nil {
		return err
	}

	if err := addObjFieldTemplateRelationIndexes(ctx, db); err != nil {
		return err
	}

	return nil
}

func addFieldTemplateIndexes(ctx context.Context, db dal.RDB) error {
	indexes := []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
			Keys: bson.D{
				{
					common.BKFieldID, 1,
				},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + common.BKFieldName,
			Keys: bson.D{
				{
					common.BKFieldName, 1,
				},
			},
			Background: true,
			Unique:     true,
		},
	}

	if err := addIndexIfNotExist(ctx, db, common.BKTableNameFieldTemplate, indexes); err != nil {
		return err
	}

	return nil
}

func addObjAttDesTemplateIndexes(ctx context.Context, db dal.RDB) error {
	indexes := []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
			Keys: bson.D{
				{
					common.BKFieldID, 1,
				},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkTemplateID_bkPropertyID",
			Keys: bson.D{
				{
					common.BKTemplateID, 1,
				},
				{
					common.BKPropertyIDField, 1,
				},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkTemplateID_bkPropertyName",
			Keys: bson.D{
				{
					common.BKTemplateID, 1,
				},
				{
					common.BKPropertyNameField, 1,
				},
			},
			Background: true,
			Unique:     true,
		},
	}

	if err := addIndexIfNotExist(ctx, db, common.BKTableNameObjAttDesTemplate, indexes); err != nil {
		return err
	}

	return nil
}

func addObjectUniqueTemplateIndexes(ctx context.Context, db dal.RDB) error {
	indexes := []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
			Keys: bson.D{
				{
					common.BKFieldID, 1,
				},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bkTemplateID_bkSupplierAccount",
			Keys: bson.D{
				{
					common.BKTemplateID, 1,
				},
				{
					"bk_supplier_account", 1,
				},
			},
			Background: true,
		},
	}

	if err := addIndexIfNotExist(ctx, db, common.BKTableNameObjectUniqueTemplate, indexes); err != nil {
		return err
	}

	return nil
}

func addObjFieldTemplateRelationIndexes(ctx context.Context, db dal.RDB) error {
	indexes := []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkTemplateID_objectID",
			Keys: bson.D{
				{
					common.BKTemplateID, 1,
				},
				{
					common.ObjectIDField, 1,
				},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "objectID_bkSupplierAccount",
			Keys: bson.D{
				{
					common.ObjectIDField, 1,
				},
				{
					"bk_supplier_account", 1,
				},
			},
			Background: true,
		},
	}

	if err := addIndexIfNotExist(ctx, db, common.BKTableNameObjFieldTemplateRelation, indexes); err != nil {
		return err
	}

	return nil
}

func addIndexIfNotExist(ctx context.Context, db dal.RDB, collection string, indexes []types.Index) error {
	existIndexArr, err := db.Table(collection).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for %s table failed, err: %v", collection, err)
		return err
	}

	existIdxMap := make(map[string]struct{})
	for _, index := range existIndexArr {
		// skip the default "_id" index for the database
		if index.Name == "_id_" {
			continue
		}
		existIdxMap[index.Name] = struct{}{}
	}

	needAddIndexes := make([]types.Index, 0)
	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}
		needAddIndexes = append(needAddIndexes, index)
	}

	if len(needAddIndexes) == 0 {
		return nil
	}

	err = db.Table(collection).BatchCreateIndexes(ctx, needAddIndexes)
	if err != nil && !db.IsDuplicatedError(err) {
		blog.Errorf("create index failed, table: %s, index: %+v, err: %v", collection, needAddIndexes, err)
		return err
	}

	return nil
}
