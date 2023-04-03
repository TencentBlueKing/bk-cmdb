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
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func addModelQuoteRelationCollection(ctx context.Context, db dal.RDB) error {

	exists, err := db.HasTable(ctx, common.BKTableNameModelQuoteRelation)
	if err != nil {
		blog.Errorf("check if table %s exists failed, err: %v, rid: %s", common.BKTableNameModelQuoteRelation, err)
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

func addModelQuoteRelationIndex(ctx context.Context, db dal.RDB) error {
	indexes := []types.Index{
		{
			Name: common.CCLogicIndexNamePrefix + "dest_model_src_model",
			Keys: bson.D{
				{common.BKDestModelField, 1},
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
