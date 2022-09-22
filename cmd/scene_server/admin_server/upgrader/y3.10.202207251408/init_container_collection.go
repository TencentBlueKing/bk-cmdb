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

package y3_10_202207251408

import (
	kubeTypes "configcenter/pkg/kube/types"
	"configcenter/pkg/storage/dal"
	"configcenter/pkg/storage/dal/types"
	"context"
	"errors"

	"configcenter/pkg/blog"
	"configcenter/pkg/common"
)

func addContainerCollection(ctx context.Context, db dal.RDB) error {

	collections := []string{
		kubeTypes.BKTableNameBaseCluster, kubeTypes.BKTableNameBaseNode,
		kubeTypes.BKTableNameBaseNamespace, kubeTypes.BKTableNameBasePod,
		kubeTypes.BKTableNameBaseContainer, kubeTypes.BKTableNameBaseDeployment,
		kubeTypes.BKTableNameBaseDaemonSet, kubeTypes.BKTableNameBaseStatefulSet,
		kubeTypes.BKTableNameGameStatefulSet, kubeTypes.BKTableNameGameDeployment,
		kubeTypes.BKTableNameBaseCronJob, kubeTypes.BKTableNameBaseJob,
		kubeTypes.BKTableNameBasePodWorkload,
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

func addContainerCollectionIndex(ctx context.Context, db dal.RDB) error {

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

	workLoadTables := []string{
		kubeTypes.BKTableNameBaseDeployment, kubeTypes.BKTableNameBaseDaemonSet,
		kubeTypes.BKTableNameBaseStatefulSet, kubeTypes.BKTableNameGameStatefulSet,
		kubeTypes.BKTableNameGameDeployment, kubeTypes.BKTableNameBaseCronJob,
		kubeTypes.BKTableNameBaseJob, kubeTypes.BKTableNameBasePodWorkload,
	}

	for _, table := range workLoadTables {
		if err := addWorkLoadTableIndexes(ctx, db, table); err != nil {
			return err
		}
	}
	return nil
}

func addContainerTableIndexes(ctx context.Context, db dal.RDB) error {
	indexes := []types.Index{
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
			Keys:       map[string]int32{common.BKFieldID: 1},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_pod_id_container_uid",
			Keys: map[string]int32{
				kubeTypes.BKPodIDField:      1,
				kubeTypes.ContainerUIDField: 1,
			},
			Background: true,
			Unique:     true,
		},
	}

	existIndexArr, err := db.Table(kubeTypes.BKTableNameBaseContainer).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for container table failed, err: %v", err)
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
		blog.Errorf("the number of indexes is not as expected, existId: %+v, indexes: %v", existIdxMap, indexes)
		return errors.New("the number of indexes is not as expected")
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}
		err = db.Table(kubeTypes.BKTableNameBaseContainer).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create index for container table failed, index: %+v, err: %v", index, err)
			return err
		}
	}
	return nil
}

var podIndexes = []types.Index{
	{
		Name: common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
		Keys: map[string]int32{
			common.BKFieldID: 1,
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bk_reference_id_reference_kind_name",
		Keys: map[string]int32{
			kubeTypes.RefIDField:   1,
			kubeTypes.RefKindField: 1,
			common.BKFieldName:     1,
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix +
			"bk_biz_id_cluster_uid_namespace_reference_kind_reference_name_name",
		Keys: map[string]int32{
			common.BKAppIDField:       1,
			kubeTypes.ClusterUIDField: 1,
			kubeTypes.KubeNamespace:   1,
			kubeTypes.RefKindField:    1,
			kubeTypes.RefNameField:    1,
			common.BKFieldName:        1,
		},
		Unique:     true,
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bk_biz_id_reference_name_reference_kind",
		Keys: map[string]int32{
			common.BKAppIDField:    1,
			kubeTypes.RefNameField: 1,
			kubeTypes.RefKindField: 1,
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bk_reference_id_reference_kind",
		Keys: map[string]int32{
			kubeTypes.RefIDField:   1,
			kubeTypes.RefKindField: 1,
		},
		Background: true,
	},
	{
		Name:       common.CCLogicIndexNamePrefix + common.BKHostIDField,
		Keys:       map[string]int32{common.BKHostIDField: 1},
		Background: true,
	},
}

func addPodTableIndexes(ctx context.Context, db dal.RDB) error {

	existIndexArr, err := db.Table(kubeTypes.BKTableNameBasePod).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for pod table failed, err: %v", err)
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
	if len(existIdxMap) != 0 && (len(existIdxMap) < len(podIndexes)) {
		blog.Errorf("the number of indexes is not as expected, existId: %+v, indexes: %v", existIdxMap, podIndexes)
		return errors.New("the number of indexes is not as expected")
	}

	for _, index := range podIndexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}
		err = db.Table(kubeTypes.BKTableNameBasePod).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create index for pod table failed, err: %v, index: %+v", err, index)
			return err
		}
	}
	return nil
}

func addWorkLoadTableIndexes(ctx context.Context, db dal.RDB, workLoadKind string) error {
	indexes := []types.Index{
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
			Keys:       map[string]int32{common.BKFieldID: 1},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_namespace_id_name",
			Keys: map[string]int32{
				kubeTypes.BKNamespaceIDField: 1,
				common.BKFieldName:           1,
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_biz_id_bk_cluster_uid_namespace_name",
			Keys: map[string]int32{
				common.BKAppIDField:       1,
				kubeTypes.ClusterUIDField: 1,
				kubeTypes.NamespaceField:  1,
				common.BKFieldName:        1,
			},
			Unique:     true,
			Background: true,
		},
		{
			Name:       common.CCLogicIndexNamePrefix + kubeTypes.ClusterUIDField,
			Keys:       map[string]int32{kubeTypes.ClusterUIDField: 1},
			Background: true,
		},
		{
			Name:       common.CCLogicIndexNamePrefix + kubeTypes.BKClusterIDFiled,
			Keys:       map[string]int32{kubeTypes.BKClusterIDFiled: 1},
			Background: true,
		},
		{
			Name:       common.CCLogicIndexNamePrefix + common.BKFieldName,
			Keys:       map[string]int32{common.BKFieldName: 1},
			Background: true,
		},
	}

	existIndexArr, err := db.Table(workLoadKind).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for %s table failed, err: %v", workLoadKind, err)
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
		blog.Errorf("the number of indexes is not as expected, existId: %+v, indexes: %v", existIdxMap, indexes)
		return errors.New("the number of indexes is not as expected")
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}
		err = db.Table(workLoadKind).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create index for %s table failed, err: %v, index: %+v", workLoadKind, err, index)
			return err
		}
	}
	return nil
}

func addNamespaceTableIndexes(ctx context.Context, db dal.RDB) error {
	indexes := []types.Index{
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
			Keys:       map[string]int32{common.BKFieldID: 1},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_cluster_id_name",
			Keys: map[string]int32{
				kubeTypes.BKClusterIDFiled: 1,
				common.BKFieldName:         1,
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_biz_id_cluster_uid_name",
			Keys: map[string]int32{
				common.BKAppIDField:       1,
				kubeTypes.ClusterUIDField: 1,
				common.BKFieldName:        1,
			},
			Unique:     true,
			Background: true,
		},
		{
			Name:       common.CCLogicIndexNamePrefix + kubeTypes.ClusterUIDField,
			Keys:       map[string]int32{kubeTypes.ClusterUIDField: 1},
			Background: true,
		},
		{
			Name:       common.CCLogicIndexNamePrefix + kubeTypes.BKClusterIDFiled,
			Keys:       map[string]int32{kubeTypes.BKClusterIDFiled: 1},
			Background: true,
		},
	}

	existIndexArr, err := db.Table(kubeTypes.BKTableNameBaseNamespace).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for namespace table failed, err: %v", err)
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
		blog.Errorf("the number of indexes is not as expected, existId: %+v, indexes: %v", existIdxMap, indexes)
		return errors.New("the number of indexes is not as expected")
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}
		err = db.Table(kubeTypes.BKTableNameBaseNamespace).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create index for namespace table failed, err: %v, index: %+v", err, index)
			return err
		}
	}
	return nil
}

func addClusterTableIndexes(ctx context.Context, db dal.RDB) error {
	indexes := []types.Index{
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
			Keys:       map[string]int32{common.BKFieldID: 1},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_biz_id_uid",
			Keys: map[string]int32{
				common.BKAppIDField: 1,
				kubeTypes.UidField:  1,
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_biz_id_name",
			Keys: map[string]int32{
				common.BKAppIDField: 1,
				common.BKFieldName:  1,
			},
			Unique:     true,
			Background: true,
		},
		{
			Name:       common.CCLogicIndexNamePrefix + common.BKAppIDField,
			Keys:       map[string]int32{common.BKAppIDField: 1},
			Background: true,
		},
		{
			Name:       common.CCLogicIndexNamePrefix + kubeTypes.XidField,
			Keys:       map[string]int32{kubeTypes.XidField: 1},
			Background: true,
		},
	}

	existIndexArr, err := db.Table(kubeTypes.BKTableNameBaseCluster).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for cluster table failed, err: %v", err)
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
		blog.Errorf("the number of indexes is not as expected, existId: %+v, indexes: %v", existIdxMap, indexes)
		return errors.New("the number of indexes is not as expected")
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}
		err = db.Table(kubeTypes.BKTableNameBaseCluster).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create index for cluster table failed, err: %v, index: %+v", err, index)
			return err
		}
	}
	return nil
}

func addNodeTableIndexes(ctx context.Context, db dal.RDB) error {
	indexes := []types.Index{
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
			Keys:       map[string]int32{common.BKFieldID: 1},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_biz_id_cluster_uid_name",
			Keys: map[string]int32{
				common.BKAppIDField:       1,
				kubeTypes.ClusterUIDField: 1,
				common.BKFieldName:        1,
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bk_cluster_id_id",
			Keys: map[string]int32{
				kubeTypes.BKClusterIDFiled: 1,
				common.BKFieldID:           1,
			},
			Unique:     true,
			Background: true,
		},
		{
			Name:       common.CCLogicIndexNamePrefix + kubeTypes.ClusterUIDField,
			Keys:       map[string]int32{kubeTypes.ClusterUIDField: 1},
			Background: true,
		},
		{
			Name:       common.CCLogicIndexNamePrefix + kubeTypes.BKClusterIDFiled,
			Keys:       map[string]int32{kubeTypes.BKClusterIDFiled: 1},
			Background: true,
		},
		{
			Name:       common.CCLogicIndexNamePrefix + common.BKHostIDField,
			Keys:       map[string]int32{common.BKHostIDField: 1},
			Background: true,
		},
	}

	existIndexArr, err := db.Table(kubeTypes.BKTableNameBaseNode).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for node table failed, err: %v", err)
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
		blog.Errorf("the number of indexes is not as expected, existId: %+v, indexes: %v", existIdxMap, indexes)
		return errors.New("the number of indexes is not as expected")
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}
		err = db.Table(kubeTypes.BKTableNameBaseNode).CreateIndex(ctx, index)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create index for node table failed, err: %v, index: %+v", err, index)
			return err
		}
	}
	return nil
}
