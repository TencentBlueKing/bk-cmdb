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

// Package level defines the topology level related logics
package level

import (
	"configcenter/src/common/http/rest"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/key"
	nodelgc "configcenter/src/source_controller/cacheservice/cache/biz-topo/logics/node"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
)

// commonCacheLevel defines common level that implements GetNodesByCache
type commonCacheLevel struct {
	topoKey   key.Key
	nextLevel LevelI
	kinds     []string
}

func newCommonCacheLevel(topoType types.TopoType, nextLevel LevelI, kinds ...string) *commonCacheLevel {
	return &commonCacheLevel{
		topoKey:   key.TopoKeyMap[topoType],
		nextLevel: nextLevel,
		kinds:     kinds,
	}
}

// GetNodesByCache get topo nodes info by cache
func (l *commonCacheLevel) GetNodesByCache(kit *rest.Kit, bizID int64) ([]types.Node, error) {
	allNodes := make([]types.Node, 0)
	for _, kind := range l.kinds {
		nodes, err := nodelgc.GetNodeInfoCache(kit, l.topoKey, bizID, kind)
		if err != nil {
			return nil, err
		}
		allNodes = append(allNodes, nodes...)
	}

	if l.nextLevel == nil {
		return allNodes, nil
	}

	childNodes, err := l.nextLevel.GetNodesByCache(kit, bizID)
	if err != nil {
		return nil, err
	}

	allNodes = nodelgc.CombineChildNodes(allNodes, childNodes)
	return allNodes, nil
}
