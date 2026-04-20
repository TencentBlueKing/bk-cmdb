/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package y3_14_202607151020

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	dbtypes "configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

// addObjLabelField add bk_labels field and index for object description table
func addObjLabelField(ctx context.Context, db dal.RDB, _ *upgrader.Config) error {
	if err := db.Table(common.BKTableNameObjDes).AddColumn(ctx, metadata.ModelFieldLabels,
		make([]string, 0)); err != nil {
		blog.Errorf("add bk_labels column failed, err: %v", err)
		return err
	}

	bkLabelsIndex := dbtypes.Index{
		Name: common.CCLogicIndexNamePrefix + "bkLabels",
		Keys: bson.D{
			{metadata.ModelFieldLabels, 1},
		},
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			metadata.ModelFieldLabels: map[string]string{common.BKDBType: "array"},
		},
	}

	existIndexes, err := db.Table(common.BKTableNameObjDes).Indexes(ctx)
	if err != nil {
		blog.Errorf("get obj des indexes failed, err: %v", err)
		return err
	}

	for _, index := range existIndexes {
		if index.Name == bkLabelsIndex.Name {
			return nil
		}
	}

	if err := db.Table(common.BKTableNameObjDes).CreateIndex(ctx, bkLabelsIndex); err != nil &&
		!db.IsDuplicatedError(err) {
		blog.Errorf("create bk_labels index %+v failed, err: %v", bkLabelsIndex, err)
		return err
	}

	return nil
}
