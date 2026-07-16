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

package y3_15_202607151050

import (
	"configcenter/pkg/tenant"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/scene_server/admin_server/logics"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func addCustomResourceCollection(kit *rest.Kit, db dal.Dal) error {
	table := kubetypes.BKTableNameBaseCustom

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
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "cluster_id",
			Keys: bson.D{
				{kubetypes.BKClusterIDFiled, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "name",
			Keys: bson.D{
				{common.BKFieldName, 1},
			},
			Background: true,
		},
	}

	err := tenant.ExecForAllTenants(func(tenantID string) error {
		tenantKit := kit.NewKit().WithTenant(tenantID)
		tenantDB := db.Shard(tenantKit.ShardOpts())

		if err := logics.CreateTable(tenantKit, tenantDB, table); err != nil {
			blog.Errorf("create table %s failed, err: %v", table, err)
			return err
		}
		if err := logics.CreateIndexes(tenantKit, tenantDB, table, indexes); err != nil {
			blog.Errorf("create table %s indexes failed, err: %v", table, err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
