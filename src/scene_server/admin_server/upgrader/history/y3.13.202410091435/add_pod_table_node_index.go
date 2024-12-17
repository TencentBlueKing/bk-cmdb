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

package y3_13_202410091435

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/kube/types"
	"configcenter/src/storage/dal"
	dbtypes "configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func addPodTableNodeIDIndex(ctx context.Context, db dal.RDB) error {
	nodeIDIndex := dbtypes.Index{
		Name: common.CCLogicIndexNamePrefix + "node_id",
		Keys: bson.D{
			{types.BKNodeIDField, 1},
		},
		Background: true,
	}

	existIndexes, err := db.Table(types.BKTableNameBasePod).Indexes(ctx)
	if err != nil {
		blog.Errorf("get pod index failed, err: %v", err)
		return err
	}

	for _, index := range existIndexes {
		if index.Name == nodeIDIndex.Name || (len(index.Keys) == 1 && index.Keys[0].Key == types.BKNodeIDField) {
			return nil
		}
	}

	err = db.Table(types.BKTableNameBasePod).CreateIndex(ctx, nodeIDIndex)
	if err != nil && !db.IsDuplicatedError(err) {
		blog.Errorf("create pod index %+v failed, err: %s", nodeIDIndex, err)
		return err
	}

	return nil
}
