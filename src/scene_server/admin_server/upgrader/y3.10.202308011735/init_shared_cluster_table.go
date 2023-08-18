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

package y3_10_202308011735

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func initSharedClusterTable(ctx context.Context, db dal.RDB) error {
	table := kubetypes.BKTableNameNsSharedClusterRel

	exists, err := db.HasTable(ctx, table)
	if err != nil {
		blog.Errorf("check if ns shared cluster relation table exists failed, err: %v", err)
		return err
	}

	if exists {
		return nil
	}

	if err = db.CreateTable(ctx, table); err != nil {
		blog.Errorf("create ns shared cluster relation table failed, err: %v", err)
		return err
	}

	return nil
}

func initSharedClusterIndex(ctx context.Context, db dal.RDB) error {
	table := kubetypes.BKTableNameNsSharedClusterRel

	indexes := []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "namespace_id",
			Keys: bson.D{
				{kubetypes.BKNamespaceIDField, 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "biz_id",
			Keys: bson.D{
				{kubetypes.BKBizIDField, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "asst_biz_id",
			Keys: bson.D{
				{kubetypes.BKAsstBizIDField, 1},
			},
			Background: true,
		},
	}

	existIndexes, err := db.Table(table).Indexes(ctx)
	if err != nil {
		blog.Errorf("get ns shared cluster relation table index failed, err: %v", err)
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
			blog.Errorf("create ns shared cluster relation table index %+v failed, err: %v", index, err)
			return err
		}
	}

	return nil
}
