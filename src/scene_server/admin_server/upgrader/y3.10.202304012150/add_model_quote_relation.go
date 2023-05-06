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

package y3_10_202304012150

import (
	"context"
	"errors"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

func addModelQuoteRelationCollection(ctx context.Context, db dal.RDB) error {

	exists, err := db.HasTable(ctx, common.BKTableNameModelQuoteRelation)
	if err != nil {
		blog.Errorf("check if table %s exists failed, err: %v", common.BKTableNameModelQuoteRelation, err)
		return err
	}

	if exists {
		return nil
	}

	if err := db.CreateTable(ctx, common.BKTableNameModelQuoteRelation); err != nil {
		return err
	}

	return nil
}

func addHiddenClassification(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	result := make([]metadata.Classification, 0)
	filter := map[string]interface{}{
		common.BKDBOR: []map[string]interface{}{
			{
				metadata.ClassFieldClassificationID: "bk_table_classification",
			},
			{
				metadata.ClassFieldClassificationName: "表格分类",
			},
		},
	}

	err := db.Table(common.BKTableNameObjClassification).Find(filter).All(ctx, &result)
	if err != nil && !db.IsNotFoundError(err) {
		blog.Errorf("find obj classification failed, err: %v", err)
		return err
	}
	if len(result) > 0 {
		blog.Errorf("model category with the same id or name already exists in the system, "+
			"please upgrade after processing, result: %v", result)
		return errors.New("classification conflict")
	}

	id, err := mongodb.Client().NextSequence(ctx, common.BKTableNameObjClassification)
	if err != nil {
		blog.Errorf("it is failed to create a new sequence id on the table(%s) of the database, error %v",
			common.BKTableNameObjClassification, err)
		return err
	}

	classification := metadata.Classification{
		ID:                 int64(id),
		OwnerID:            conf.OwnerID,
		ClassificationType: metadata.HiddenType,
		ClassificationID:   "bk_table_classification",
		ClassificationName: "表格分类",
	}

	if err := db.Table(common.BKTableNameObjClassification).Insert(ctx, classification); err != nil {
		blog.Errorf("insert hidden classification failed err: %v", err)
	}
	return nil
}

func addModelQuoteRelationIndex(ctx context.Context, db dal.RDB) error {
	indexes := []types.Index{
		{
			Name: common.CCLogicIndexNamePrefix + "dest_model_src_model",
			Keys: bson.D{
				{common.BKDestModelField, 1},
				{common.BKOwnerIDField, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "dest_model_src_model",
			Keys: bson.D{
				{common.BKSrcModelField, 1},
				{common.BKPropertyIDField, 1},
				{common.BKOwnerIDField, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "dest_model_src_model",
			Keys: bson.D{
				{common.BKSrcModelField, 1},
				{common.BKOwnerIDField, 1},
			},
			Background: true,
		},
	}

	existIndexArr, err := db.Table(common.BKTableNameModelQuoteRelation).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for model quote relation table failed, err: %v", err)
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

	// the number of indexes is not as expected.
	if len(existIdxMap) != 0 && (len(existIdxMap) < len(indexes)) {
		blog.Errorf("the number of indexes is not as expected, existId: %+v, indexes: %v", existIdxMap, indexes)
		return errors.New("the number of indexes is not as expected")
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}
		err = db.Table(common.BKTableNameModelQuoteRelation).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create index for model quote relation table failed, err: %v, index: %+v", err, index)
			return err
		}
	}
	return nil
}
