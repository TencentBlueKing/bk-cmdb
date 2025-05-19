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

package tokenhandler

import (
	"context"

	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/stream/types"
)

var _ types.TaskTokenHandler = new(MemoryHandler)

// MemoryHandler is a token handler that stores the token in process memory
type MemoryHandler struct {
	dbTokenMap map[string]*types.TokenInfo
}

// NewMemoryTokenHandler generate a new memory event token handler
func NewMemoryTokenHandler() *MemoryHandler {
	return &MemoryHandler{
		dbTokenMap: make(map[string]*types.TokenInfo),
	}
}

// SetLastWatchToken set last event watch token
func (m *MemoryHandler) SetLastWatchToken(_ context.Context, uuid string, _ local.DB, token *types.TokenInfo) error {
	m.dbTokenMap[uuid] = token
	return nil
}

// GetStartWatchToken get event start watch token
func (m *MemoryHandler) GetStartWatchToken(_ context.Context, uuid string, _ local.DB) (*types.TokenInfo, error) {
	token, exists := m.dbTokenMap[uuid]
	if !exists {
		return new(types.TokenInfo), nil
	}
	return token, nil
}
