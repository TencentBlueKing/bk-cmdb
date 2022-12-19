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

package y3_10_202209281408

import (
	"context"
	"errors"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func addKubeCollection(ctx context.Context, db dal.RDB) error {

	collections := []string{
		kubetypes.BKTableNameBaseCluster, kubetypes.BKTableNameBaseNode,
		kubetypes.BKTableNameBaseNamespace, kubetypes.BKTableNameBasePod,
		kubetypes.BKTableNameBaseContainer, kubetypes.BKTableNameBaseDeployment,
		kubetypes.BKTableNameBaseDaemonSet, kubetypes.BKTableNameBaseStatefulSet,
		kubetypes.BKTableNameGameStatefulSet, kubetypes.BKTableNameGameDeployment,
		kubetypes.BKTableNameBaseCronJob, kubetypes.BKTableNameBaseJob,
		kubetypes.BKTableNameBasePodWorkload, kubetypes.BKTableNsClusterRelation,
		kubetypes.BKTableNodeClusterRelation,
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

func addKubeCollectionIndex(ctx context.Context, db dal.RDB) error {

	if err := addClusterTableIndexes(ctx, db); err != nil {
		return err
	}

	if err := addNodeTableIndexes(ctx, db); err != nil {
		return err
	}

	if err := addNamespaceTableIndexes(ctx, db); err != nil {
		return err
	}

	if err := addPodTableIndexes(ctx, db); err != nil {
		return err
	}

	if err := addContainerTableIndexes(ctx, db); err != nil {
		return err
	}

	if err := addKubeNsRelationTableIndexes(ctx, db); err != nil {
		return err
	}

	if err := addKubeNodeRelationTableIndexes(ctx, db); err != nil {
		return err
	}

	workLoadTables := []string{
		kubetypes.BKTableNameBaseDeployment, kubetypes.BKTableNameBaseDaemonSet,
		kubetypes.BKTableNameBaseStatefulSet, kubetypes.BKTableNameGameStatefulSet,
		kubetypes.BKTableNameGameDeployment, kubetypes.BKTableNameBaseCronJob,
		kubetypes.BKTableNameBaseJob, kubetypes.BKTableNameBasePodWorkload,
	}

	for _, table := range workLoadTables {
		if err := addWorkLoadTableIndexes(ctx, db, table); err != nil {
			return err
		}
	}
	return nil
}

func setIndex(ctx context.Context, db dal.RDB, collection string, indexes []types.Index) error {
	existIndexArr, err := db.Table(collection).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for %s table failed, err: %v", collection, err)
		return err
	}

	existIdxMap := make(map[string]bool)
	for _, index := range existIndexArr {
		// skip the default "_id" index for the database
		if index.Name == "_id_" {
			continue
		}
		existIdxMap[index.Name] = true
	}

	// the number of indexes is not as expected.
	if len(existIdxMap) != 0 && (len(existIdxMap) < len(indexes)) {
		blog.Errorf("the number of indexes is not as expected, collection: %s, existId: %+v, indexes: %v",
			collection, existIdxMap, indexes)
		return errors.New("the number of indexes is not as expected")
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}
		err = db.Table(collection).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create index for %s table failed, index: %+v, err: %v", collection, index, err)
			return err
		}
	}
	return nil
}

func addKubeNsRelationTableIndexes(ctx context.Context, db dal.RDB) error {

	indexes := []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "namespace_id_cluster_id",
			Keys: bson.D{
				{kubetypes.BKNamespaceIDField, 1},
				{kubetypes.BKClusterIDField, 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "namespace_id_cluster_uid",
			Keys: bson.D{
				{kubetypes.BKNamespaceIDField, 1},
				{kubetypes.ClusterUIDField, 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bk_biz_id",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
	}

	if err := setIndex(ctx, db, kubetypes.BKTableNsClusterRelation, indexes); err != nil {
		return err
	}
	return nil
}

func addKubeNodeRelationTableIndexes(ctx context.Context, db dal.RDB) error {

	indexes := []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "node_id_cluster_id",
			Keys: bson.D{
				{kubetypes.BKNodeIDField, 1},
				{kubetypes.BKClusterIDField, 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "node_id_cluster_uid",
			Keys: bson.D{
				{kubetypes.BKNodeIDField, 1},
				{kubetypes.ClusterUIDField, 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bk_biz_id",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
	}

	if err := setIndex(ctx, db, kubetypes.BKTableNodeClusterRelation, indexes); err != nil {
		return err
	}
	return nil
}

func addContainerTableIndexes(ctx context.Context, db dal.RDB) error {
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
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_pod_id_container_uid",
			Keys: bson.D{
				{kubetypes.BKPodIDField, 1},
				{kubetypes.ContainerUIDField, 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "pod_id",
			Keys: bson.D{
				{kubetypes.BKPodIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
	}

	if err := setIndex(ctx, db, kubetypes.BKTableNameBaseContainer, indexes); err != nil {
		return err
	}
	return nil
}

func addPodTableIndexes(ctx context.Context, db dal.RDB) error {

	var indexes = []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
			Keys: bson.D{
				{common.BKFieldID, 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_reference_id_reference_kind_name",
			Keys: bson.D{
				{kubetypes.RefIDField, 1},
				{kubetypes.RefKindField, 1},
				{common.BKFieldName, 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "biz_id_reference_name_reference_kind",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{kubetypes.RefNameField, 1},
				{kubetypes.RefKindField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "biz_id_cluster_id",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{kubetypes.BKClusterIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "biz_id_namespace_id",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{kubetypes.BKNamespaceIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "reference_id_reference_kind",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{kubetypes.RefIDField, 1},
				{kubetypes.RefKindField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "biz_id_name",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{common.BKFieldName, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bk_host_id",
			Keys: bson.D{
				{common.BKHostIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
	}

	if err := setIndex(ctx, db, kubetypes.BKTableNameBasePod, indexes); err != nil {
		return err
	}

	return nil
}

func addWorkLoadTableIndexes(ctx context.Context, db dal.RDB, workLoadKind string) error {
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
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_namespace_id_name",
			Keys: bson.D{
				{kubetypes.BKNamespaceIDField, 1},
				{common.BKFieldName, 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "biz_id_cluster_uid",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{kubetypes.ClusterUIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "biz_id_cluster_id",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{kubetypes.BKClusterIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "name",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{common.BKFieldName, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
	}

	if err := setIndex(ctx, db, workLoadKind, indexes); err != nil {
		return err
	}
	return nil
}

func addNamespaceTableIndexes(ctx context.Context, db dal.RDB) error {
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
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_cluster_id_name",
			Keys: bson.D{
				{kubetypes.BKClusterIDField, 1},
				{common.BKFieldName, 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "biz_id_cluster_uid",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{kubetypes.ClusterUIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "biz_id_cluster_id",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{kubetypes.BKClusterIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "name",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{common.BKFieldName, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
	}

	if err := setIndex(ctx, db, kubetypes.BKTableNameBaseNamespace, indexes); err != nil {
		return err
	}

	return nil
}

func addClusterTableIndexes(ctx context.Context, db dal.RDB) error {
	indexes := []types.Index{
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
			Keys:       bson.D{{common.BKFieldID, 1}},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "cluster_uid",
			Keys: bson.D{
				{kubetypes.UidField, 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bk_biz_id_name",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{common.BKFieldName, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + common.BKAppIDField,
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "xid",
			Keys: bson.D{
				{kubetypes.XidField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
	}

	if err := setIndex(ctx, db, kubetypes.BKTableNameBaseCluster, indexes); err != nil {
		return err
	}

	return nil
}

func addNodeTableIndexes(ctx context.Context, db dal.RDB) error {
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
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_cluster_id_name",
			Keys: bson.D{
				{kubetypes.BKClusterIDField, 1},
				{common.BKFieldName, 1},
			},
			Unique:     true,
			Background: true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "cluster_uid_name",
			Keys: bson.D{
				{kubetypes.ClusterUIDField, 1},
				{common.BKFieldName, 1},
			},
			Unique:     true,
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "biz_id_cluster_uid",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{kubetypes.ClusterUIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "biz_id_cluster_id",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{kubetypes.BKClusterIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "biz_id_host_id",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{common.BKHostIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "biz_id_name",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{common.BKFieldName, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
	}

	if err := setIndex(ctx, db, kubetypes.BKTableNameBaseNode, indexes); err != nil {
		return err
	}
	return nil
}
