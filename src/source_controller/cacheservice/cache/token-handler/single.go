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
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"
)

// SingleHandler is a token handler for single event that watches db only once
// token data: {"_id": $key, "token": $token, "start_at_time": $start_at_time}
type SingleHandler struct {
	key string
}

// NewSingleTokenHandler generate a new event token handler
func NewSingleTokenHandler(key string) *SingleHandler {
	return &SingleHandler{
		key: key,
	}
}

// SetLastWatchToken set last event watch token
func (m *SingleHandler) SetLastWatchToken(ctx context.Context, uuid string, watchDB local.DB,
	token *types.TokenInfo) error {

	filter := map[string]interface{}{
		"_id": m.genWatchTokenID(uuid),
	}

	if err := watchDB.Table(common.BKTableNameCacheWatchToken).Upsert(ctx, filter, token); err != nil {
		blog.Errorf("set event %s-%s last watch token failed, data: %+v, err: %v", m.key, uuid, *token, err)
		return err
	}
	return nil
}

// GetStartWatchToken get event start watch token
func (m *SingleHandler) GetStartWatchToken(ctx context.Context, uuid string, watchDB local.DB) (*types.TokenInfo,
	error) {

	filter := map[string]interface{}{
		"_id": m.genWatchTokenID(uuid),
	}

	tokenInfo := new(types.TokenInfo)
	err := watchDB.Table(common.BKTableNameCacheWatchToken).Find(filter).One(ctx, &tokenInfo)
	if err != nil {
		if !mongodb.IsNotFoundError(err) {
			blog.Errorf("get event start watch token by filter: %+v failed, err: %v", filter, err)
			return nil, err
		}

		return &types.TokenInfo{Token: "", StartAtTime: &types.TimeStamp{Sec: uint32(time.Now().Unix())}}, nil
	}

	return tokenInfo, nil
}

func (m *SingleHandler) genWatchTokenID(uuid string) string {
	return m.key + ":" + uuid
}

// IsTokenExists check if any event token exists for all watch dbs
func (m *SingleHandler) IsTokenExists(ctx context.Context, watchDal dal.Dal) (bool, error) {
	filter := map[string]interface{}{
		"_id": map[string]interface{}{
			common.BKDBLIKE: fmt.Sprintf("^%s:", m.key),
		},
	}

	exists := false
	err := watchDal.ExecForAllDB(func(db local.DB) error {
		if exists {
			return nil
		}

		cnt, err := db.Table(common.BKTableNameCacheWatchToken).Find(filter).Count(ctx)
		if err != nil {
			blog.Errorf("check if event token exists failed, filter: %+v, err: %v", filter, err)
			return err
		}

		exists = cnt > 0
		return nil
	})
	if err != nil {
		return false, err
	}

	return exists, nil
}
