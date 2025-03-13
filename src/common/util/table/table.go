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
	fullsynccond "configcenter/pkg/cache/full-sync-cond"
	"configcenter/src/common"
	kubetypes "configcenter/src/kube/types"
)

var delArchiveCollMap = map[string]struct{}{
	common.BKTableNameModuleHostConfig:        {},
	common.BKTableNameBaseHost:                {},
	common.BKTableNameBaseApp:                 {},
	common.BKTableNameBaseSet:                 {},
	common.BKTableNameBaseModule:              {},
	common.BKTableNameSetTemplate:             {},
	common.BKTableNameBaseProcess:             {},
	common.BKTableNameProcessInstanceRelation: {},
	common.BKTableNameBaseBizSet:              {},
	common.BKTableNameBasePlat:                {},
	common.BKTableNameBaseProject:             {},
	fullsynccond.BKTableNameFullSyncCond:      {},

	common.BKTableNameBaseInst:         {},
	common.BKTableNameMainlineInstance: {},
	common.BKTableNameInstAsst:         {},

	common.BKTableNameServiceInstance: {},

	kubetypes.BKTableNameBaseCluster:        {},
	kubetypes.BKTableNameBaseNode:           {},
	kubetypes.BKTableNameBaseNamespace:      {},
	kubetypes.BKTableNameBaseWorkload:       {},
	kubetypes.BKTableNameBaseDeployment:     {},
	kubetypes.BKTableNameBaseStatefulSet:    {},
	kubetypes.BKTableNameBaseDaemonSet:      {},
	kubetypes.BKTableNameGameDeployment:     {},
	kubetypes.BKTableNameGameStatefulSet:    {},
	kubetypes.BKTableNameBaseCronJob:        {},
	kubetypes.BKTableNameBaseJob:            {},
	kubetypes.BKTableNameBasePodWorkload:    {},
	kubetypes.BKTableNameBaseCustom:         {},
	kubetypes.BKTableNameBasePod:            {},
	kubetypes.BKTableNameBaseContainer:      {},
	kubetypes.BKTableNameNsSharedClusterRel: {},
}

// NeedPreImageTable check if table needs to enable change stream pre-image
func NeedPreImageTable(table string) bool {
	_, exists := delArchiveCollMap[table]
	if exists {
		return true
	}

	if !common.IsObjectShardingTable(table) {
		return false
	}

	return true
}
