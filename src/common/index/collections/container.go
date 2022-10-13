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

package collections

import (
	"configcenter/src/common"
	kubeTypes "configcenter/src/kube/types"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	registerIndexes(kubeTypes.BKTableNameBaseCluster, commClusterIndexes)
	registerIndexes(kubeTypes.BKTableNameBaseNode, commNodeIndexes)
	registerIndexes(kubeTypes.BKTableNameBaseNamespace, commNamespaceIndexes)
	registerIndexes(kubeTypes.BKTableNameBasePod, commPodIndexes)
	registerIndexes(kubeTypes.BKTableNameBaseContainer, commContainerIndexes)

	workLoadTables := []string{
		kubeTypes.BKTableNameBaseDeployment, kubeTypes.BKTableNameBaseDaemonSet,
		kubeTypes.BKTableNameBaseStatefulSet, kubeTypes.BKTableNameGameStatefulSet,
		kubeTypes.BKTableNameGameDeployment, kubeTypes.BKTableNameBaseCronJob,
		kubeTypes.BKTableNameBaseJob, kubeTypes.BKTableNameBasePodWorkload,
	}
	for _, table := range workLoadTables {
		registerIndexes(table, commWorkLoadIndexes)
	}
}

//Keys: bson.D{
//{common.BKAppIDField, 1},
//{kubeTypes.UidField, 1},
//},

var commWorkLoadIndexes = []types.Index{
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
		Keys:       bson.D{{common.BKFieldID, 1}},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bk_namespace_id_name",
		Keys: bson.D{
			{kubeTypes.BKNamespaceIDField, 1},
			{common.BKFieldName, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bk_biz_id_cluster_uid_namespace_name",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{kubeTypes.ClusterUIDField, 1},
			{kubeTypes.NamespaceField, 1},
			{common.BKFieldName, 1},
		},
		Unique:     true,
		Background: true,
	},
	{
		Name:       common.CCLogicIndexNamePrefix + kubeTypes.ClusterUIDField,
		Keys:       bson.D{{kubeTypes.ClusterUIDField, 1}},
		Background: true,
	},
	{
		Name:       common.CCLogicIndexNamePrefix + kubeTypes.BKClusterIDFiled,
		Keys:       bson.D{{kubeTypes.BKClusterIDFiled, 1}},
		Background: true,
	},
	{
		Name:       common.CCLogicIndexNamePrefix + common.BKFieldName,
		Keys:       bson.D{{common.BKFieldName, 1}},
		Background: true,
	},
}
var commContainerIndexes = []types.Index{
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
		Keys:       bson.D{{common.BKFieldID, 1}},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bk_pod_id_container_uid",
		Keys: bson.D{
			{kubeTypes.BKPodIDField, 1},
			{kubeTypes.ContainerUIDField, 1},
		},
		Background: true,
		Unique:     true,
	},
}
var commPodIndexes = []types.Index{
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
			{kubeTypes.RefIDField, 1},
			{kubeTypes.RefKindField, 1},
			{common.BKFieldName, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix +
			"bk_biz_id_cluster_uid_namespace_reference_kind_reference_name_name",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{kubeTypes.ClusterUIDField, 1},
			{kubeTypes.NamespaceField, 1},
			{kubeTypes.RefIDField, 1},
			{kubeTypes.RefNameField, 1},
			{common.BKFieldName, 1},
		},
		Unique:     true,
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bk_biz_id_reference_name_reference_kind",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{kubeTypes.RefNameField, 1},
			{kubeTypes.RefIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bk_reference_id_reference_kind",
		Keys: bson.D{
			{kubeTypes.RefIDField, 1},
			{kubeTypes.RefKindField, 1},
		},
		Background: true,
	},
	{
		Name:       common.CCLogicIndexNamePrefix + common.BKHostIDField,
		Keys:       bson.D{{common.BKHostIDField, 1}},
		Background: true,
	},
}
var commNamespaceIndexes = []types.Index{
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
		Keys:       bson.D{{common.BKFieldID, 1}},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bk_cluster_id_name",
		Keys: bson.D{
			{kubeTypes.BKClusterIDFiled, 1},
			{common.BKFieldName, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bk_biz_id_cluster_uid_name",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{kubeTypes.ClusterUIDField, 1},
			{common.BKFieldName, 1},
		},
		Unique:     true,
		Background: true,
	},
	{
		Name:       common.CCLogicIndexNamePrefix + kubeTypes.ClusterUIDField,
		Keys:       bson.D{{kubeTypes.ClusterUIDField, 1}},
		Background: true,
	},
	{
		Name:       common.CCLogicIndexNamePrefix + kubeTypes.BKClusterIDFiled,
		Keys:       bson.D{{kubeTypes.BKClusterIDFiled, 1}},
		Background: true,
	},
}
var commNodeIndexes = []types.Index{
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
		Keys:       bson.D{{common.BKFieldID, 1}},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bk_biz_id_cluster_uid_name",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{kubeTypes.ClusterUIDField, 1},
			{common.BKFieldID, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bk_cluster_id_id",
		Keys: bson.D{
			{kubeTypes.BKClusterIDFiled, 1},
			{common.BKFieldID, 1},
		},
		Unique:     true,
		Background: true,
	},
	{
		Name:       common.CCLogicIndexNamePrefix + kubeTypes.ClusterUIDField,
		Keys:       bson.D{{kubeTypes.ClusterUIDField, 1}},
		Background: true,
	},
	{
		Name:       common.CCLogicIndexNamePrefix + kubeTypes.BKClusterIDFiled,
		Keys:       bson.D{{kubeTypes.BKClusterIDFiled, 1}},
		Background: true,
	},
	{
		Name:       common.CCLogicIndexNamePrefix + common.BKHostIDField,
		Keys:       bson.D{{common.BKHostIDField, 1}},
		Background: true,
	},
}
var commClusterIndexes = []types.Index{

	{
		Name:       common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
		Keys:       bson.D{{common.BKFieldID, 1}},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bk_biz_id_uid",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{kubeTypes.UidField, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bk_biz_id_name",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{common.BKFieldName, 1},
		},
		Unique:     true,
		Background: true,
	},
	{
		Name:       common.CCLogicIndexNamePrefix + common.BKAppIDField,
		Keys:       bson.D{{common.BKAppIDField, 1}},
		Background: true,
	},
	{
		Name:       common.CCLogicIndexNamePrefix + kubeTypes.XidField,
		Keys:       bson.D{{kubeTypes.XidField, 1}},
		Background: true,
	},
}
