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
	"context"

	"configcenter/src/common/mapstr"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
)

// LevelI is the interface for topology tree level
type LevelI interface {
	GetNodesByDB(ctx context.Context, bizID int64, cond []mapstr.MapStr, rid string) ([]types.Node, error)
	GetNodesByCache(ctx context.Context, bizID int64, rid string) ([]types.Node, error)
}
