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
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	registerIndexes(kubetypes.BKTableNameBaseCluster, commClusterIndexes)
	registerIndexes(kubetypes.BKTableNameBaseNode, commNodeIndexes)
	registerIndexes(kubetypes.BKTableNameBaseNamespace, commNamespaceIndexes)
	registerIndexes(kubetypes.BKTableNameBasePod, commPodIndexes)
	registerIndexes(kubetypes.BKTableNameBaseContainer, commContainerIndexes)

	workLoadTables := []string{
		kubetypes.BKTableNameBaseDeployment, kubetypes.BKTableNameBaseDaemonSet,
		kubetypes.BKTableNameBaseStatefulSet, kubetypes.BKTableNameGameStatefulSet,
		kubetypes.BKTableNameGameDeployment, kubetypes.BKTableNameBaseCronJob,
		kubetypes.BKTableNameBaseJob, kubetypes.BKTableNameBasePodWorkload,
	}
	for _, table := range workLoadTables {
		registerIndexes(table, commWorkLoadIndexes)
	}
}

var commWorkLoadIndexes = []types.Index{
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
		Name: common.CCLogicIndexNamePrefix + "biz_id_name",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{common.BKFieldName, 1},
			{common.BkSupplierAccount, 1},
		},
		Background: true,
	},
}

var commContainerIndexes = []types.Index{
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
			{kubetypes.RefIDField, 1},
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
		Name: common.CCLogicIndexNamePrefix + "biz_id_reference_id_reference_kind",
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

var commNamespaceIndexes = []types.Index{
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

var commNodeIndexes = []types.Index{
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

var commClusterIndexes = []types.Index{

	{
		Name: common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
		Keys: bson.D{
			{common.BKFieldID, 1},
		},
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
		Name: common.CCLogicUniqueIdxNamePrefix + "bk_biz_id_name",
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
