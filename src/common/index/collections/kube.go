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
	registerIndexes(kubetypes.BKTableNameNsSharedClusterRel, nsSharedClusterRelIndexes)

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
		Name: common.CCLogicUniqueIdxNamePrefix + "ID",
		Keys: bson.D{
			{common.BKFieldID, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bkNamespaceID_name",
		Keys: bson.D{
			{kubetypes.BKNamespaceIDField, 1},
			{common.BKFieldName, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "clusterUID",
		Keys: bson.D{
			{kubetypes.ClusterUIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkClusterID",
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

var commContainerIndexes = []types.Index{
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "ID",
		Keys: bson.D{
			{common.BKFieldID, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bkPodID_containerUID",
		Keys: bson.D{
			{kubetypes.BKPodIDField, 1},
			{kubetypes.ContainerUIDField, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkPodID",
		Keys: bson.D{
			{kubetypes.BKPodIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkBizID",
		Keys: bson.D{
			{kubetypes.BKBizIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkClusterID",
		Keys: bson.D{
			{kubetypes.BKClusterIDFiled, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkNamespaceID",
		Keys: bson.D{
			{kubetypes.BKNamespaceIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "refID_refKind",
		Keys: bson.D{
			{kubetypes.RefIDField, 1},
			{kubetypes.RefKindField, 1},
		},
		Background: true,
	},
}

var commPodIndexes = []types.Index{
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "ID",
		Keys: bson.D{
			{common.BKFieldID, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "refID_refKind_name",
		Keys: bson.D{
			{kubetypes.RefIDField, 1},
			{kubetypes.RefKindField, 1},
			{common.BKFieldName, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkHostID",
		Keys: bson.D{
			{common.BKHostIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "refName_refID",
		Keys: bson.D{
			{kubetypes.RefNameField, 1},
			{kubetypes.RefIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "clusterID",
		Keys: bson.D{
			{kubetypes.BKClusterIDFiled, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "clusterUID",
		Keys: bson.D{
			{kubetypes.ClusterUIDField, 1},
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
	{
		Name: common.CCLogicIndexNamePrefix + "bkNodeID",
		Keys: bson.D{
			{kubetypes.BKNodeIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkNamespaceID",
		Keys: bson.D{
			{kubetypes.BKNamespaceIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "refID_refKind",
		Keys: bson.D{
			{kubetypes.RefIDField, 1},
			{kubetypes.RefKindField, 1},
		},
		Background: true,
	},
}

var commNamespaceIndexes = []types.Index{
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "ID",
		Keys: bson.D{
			{common.BKFieldID, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bkClusterID_name",
		Keys: bson.D{
			{kubetypes.BKClusterIDFiled, 1},
			{common.BKFieldName, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "clusterUID",
		Keys: bson.D{
			{kubetypes.ClusterUIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "clusterID",
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

var commNodeIndexes = []types.Index{
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "ID",
		Keys: bson.D{
			{common.BKFieldID, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bkClusterID_name",
		Keys: bson.D{
			{kubetypes.BKClusterIDFiled, 1},
			{common.BKFieldName, 1},
		},
		Unique:     true,
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkBizID_clusterUID",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{kubetypes.ClusterUIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkBizID_bkClusterID",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{kubetypes.BKClusterIDFiled, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkBizID_bkHostID",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{common.BKHostIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkBizID_name",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{common.BKFieldName, 1},
		},
		Background: true,
	},
}

var commClusterIndexes = []types.Index{
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "ID",
		Keys: bson.D{
			{common.BKFieldID, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "UID",
		Keys: bson.D{
			{kubetypes.UidField, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkBizID_name",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{common.BKFieldName, 1}},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkBizID",
		Keys: bson.D{
			{common.BKAppIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "xid",
		Keys: bson.D{
			{kubetypes.XidField, 1},
		},
		Background: true,
	},
}

var nsSharedClusterRelIndexes = []types.Index{
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bkNamespaceID",
		Keys: bson.D{
			{kubetypes.BKNamespaceIDField, 1},
		},
		Background: true,
		Unique:     true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkBizID",
		Keys: bson.D{
			{kubetypes.BKBizIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bkAsstBizID",
		Keys: bson.D{
			{kubetypes.BKAsstBizIDField, 1},
		},
		Background: true,
	},
}
