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

// Package table is table related util package
package table

import (
	"configcenter/src/common"
	kubetypes "configcenter/src/kube/types"
)

var delArchiveCollMap = map[string]string{
	common.BKTableNameModuleHostConfig:        common.BKTableNameDelArchive,
	common.BKTableNameBaseHost:                common.BKTableNameDelArchive,
	common.BKTableNameBaseApp:                 common.BKTableNameDelArchive,
	common.BKTableNameBaseSet:                 common.BKTableNameDelArchive,
	common.BKTableNameBaseModule:              common.BKTableNameDelArchive,
	common.BKTableNameSetTemplate:             common.BKTableNameDelArchive,
	common.BKTableNameBaseProcess:             common.BKTableNameDelArchive,
	common.BKTableNameProcessInstanceRelation: common.BKTableNameDelArchive,
	common.BKTableNameBaseBizSet:              common.BKTableNameDelArchive,
	common.BKTableNameBasePlat:                common.BKTableNameDelArchive,
	common.BKTableNameBaseProject:             common.BKTableNameDelArchive,

	common.BKTableNameBaseInst:         common.BKTableNameDelArchive,
	common.BKTableNameMainlineInstance: common.BKTableNameDelArchive,
	common.BKTableNameInstAsst:         common.BKTableNameDelArchive,

	common.BKTableNameServiceInstance: common.BKTableNameDelArchive,

	kubetypes.BKTableNameBaseCluster:        common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameBaseNode:           common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameBaseNamespace:      common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameBaseWorkload:       common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameBaseDeployment:     common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameBaseStatefulSet:    common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameBaseDaemonSet:      common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameGameDeployment:     common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameGameStatefulSet:    common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameBaseCronJob:        common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameBaseJob:            common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameBasePodWorkload:    common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameBaseCustom:         common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameBasePod:            common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameBaseContainer:      common.BKTableNameKubeDelArchive,
	kubetypes.BKTableNameNsSharedClusterRel: common.BKTableNameKubeDelArchive,
}

// GetDelArchiveTable get delete archive table
func GetDelArchiveTable(table string) (string, bool) {
	delArchiveTable, exists := delArchiveCollMap[table]
	if exists {
		return delArchiveTable, true
	}

	if !common.IsObjectShardingTable(table) {
		return "", false
	}

	return common.BKTableNameDelArchive, true
}

// GetDelArchiveFields get delete archive fields by table
func GetDelArchiveFields(table string) []string {
	switch table {
	case common.BKTableNameServiceInstance:
		return []string{common.BKFieldID}
	}

	return make([]string, 0)
}
