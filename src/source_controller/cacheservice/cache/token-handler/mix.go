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

// Package tokenhandler defines the token handler for cache
package tokenhandler

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"
)

var _ types.TaskTokenHandler = new(MixHandler)

// MixHandler is a token handler for mix event composed of multiple types of events
// token data: {"_id": $mixKey, $collection: {"token": $token, "start_at_time": $start_at_time}}
type MixHandler struct {
	mixKey     string
	collection string
}

// NewMixTokenHandler generate a new mix event token handler
func NewMixTokenHandler(mixKey, collection string) *MixHandler {
	return &MixHandler{
		mixKey:     mixKey,
		collection: collection,
	}
}

// SetLastWatchToken set last mix event watch token
func (m *MixHandler) SetLastWatchToken(ctx context.Context, uuid string, watchDB local.DB,
	token *types.TokenInfo) error {

	filter := map[string]interface{}{
		"_id": m.genWatchTokenID(uuid),
	}

	tokenInfo := mapstr.MapStr{
		m.collection: token,
	}

	if err := watchDB.Table(common.BKTableNameCacheWatchToken).Upsert(ctx, filter, tokenInfo); err != nil {
		blog.Errorf("set mix event %s:%s last watch token failed, data: %+v, err: %v", m.mixKey, uuid, tokenInfo, err)
		return err
	}
	return nil
}

// GetStartWatchToken get mix event start watch token
func (m *MixHandler) GetStartWatchToken(ctx context.Context, uuid string, watchDB local.DB) (*types.TokenInfo, error) {
	filter := map[string]interface{}{
		"_id": m.genWatchTokenID(uuid),
	}

	data := make(map[string]*types.TokenInfo)
	err := watchDB.Table(common.BKTableNameCacheWatchToken).Find(filter).Fields(m.collection).One(ctx, &data)
	if err != nil {
		if !mongodb.IsNotFoundError(err) {
			blog.Errorf("get mix event start watch token by filter: %+v failed, err: %v", filter, err)
			return nil, err
		}

		return &types.TokenInfo{Token: "", StartAtTime: &types.TimeStamp{Sec: uint32(time.Now().Unix())}}, nil
	}

	tokenInfo, exist := data[m.collection]
	if !exist {
		blog.Infof("mix event %s:%s start watch token is not found", m.mixKey, uuid)
		return &types.TokenInfo{Token: "", StartAtTime: &types.TimeStamp{Sec: uint32(time.Now().Unix())}}, nil
	}

	return tokenInfo, nil
}

func (m *MixHandler) genWatchTokenID(uuid string) string {
	return m.mixKey + ":" + uuid
}
