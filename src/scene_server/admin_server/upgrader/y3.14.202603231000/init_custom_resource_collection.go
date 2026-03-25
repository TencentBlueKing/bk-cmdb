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

package y3_14_202603231000

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func addCustomResourceCollection(ctx context.Context, db dal.RDB) error {
	exists, err := db.HasTable(ctx, kubetypes.BKTableNameBaseCustom)
	if err != nil {
		blog.Errorf("check if %s table exists failed, err: %v", kubetypes.BKTableNameBaseCustom, err)
		return err
	}

	if exists {
		return nil
	}

	if err = db.CreateTable(ctx, kubetypes.BKTableNameBaseCustom); err != nil {
		blog.Errorf("create %s table failed, err: %v", kubetypes.BKTableNameBaseCustom, err)
		return err
	}

	return nil
}

func addCustomResourceCollectionIndex(ctx context.Context, db dal.RDB) error {
	indexes := []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
			Keys: bson.D{
				{common.BKFieldID, 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_namespace_id_cr_kind_cr_api_version_name",
			Keys: bson.D{
				{kubetypes.BKNamespaceIDField, 1},
				{kubetypes.CRKindField, 1},
				{kubetypes.CRApiVersionField, 1},
				{common.BKFieldName, 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "cluster_uid",
			Keys: bson.D{
				{kubetypes.ClusterUIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "cluster_id",
			Keys: bson.D{
				{kubetypes.BKClusterIDFiled, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "name",
			Keys: bson.D{
				{common.BKFieldName, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
	}

	existIndexArr, err := db.Table(kubetypes.BKTableNameBaseCustom).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for %s table failed, err: %v", kubetypes.BKTableNameBaseCustom, err)
		return err
	}

	existIdxMap := make(map[string]bool)
	for _, index := range existIndexArr {
		existIdxMap[index.Name] = true
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}

		err = db.Table(kubetypes.BKTableNameBaseCustom).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create custom resource index failed, index: %+v, err: %v", index, err)
			return err
		}
	}

	return nil
}
