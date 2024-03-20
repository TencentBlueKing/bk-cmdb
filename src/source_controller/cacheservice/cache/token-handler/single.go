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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/stream/types"
)

var _ types.TokenHandler = new(SingleHandler)

// SingleHandler is a token handler for single event that watches db only once
// token data: {"_id": $key, "token": $token, "start_at_time": $start_at_time}
type SingleHandler struct {
	key string
	db  dal.DB
}

// NewSingleTokenHandler generate a new event token handler
func NewSingleTokenHandler(key string, db dal.DB) *SingleHandler {
	return &SingleHandler{
		key: key,
		db:  db,
	}
}

// SetLastWatchToken set last event watch token
func (m *SingleHandler) SetLastWatchToken(ctx context.Context, token string) error {
	filter := map[string]interface{}{
		"_id": m.key,
	}

	tokenInfo := mapstr.MapStr{
		"token": token,
	}

	if err := m.db.Table(common.BKTableNameSystem).Upsert(ctx, filter, tokenInfo); err != nil {
		blog.Errorf("set event %s last watch token failed, data: %+v, err: %v", m.key, tokenInfo, err)
		return err
	}
	return nil
}

// GetStartWatchToken get event start watch token
func (m *SingleHandler) GetStartWatchToken(ctx context.Context) (string, error) {
	filter := map[string]interface{}{
		"_id": m.key,
	}

	tokenInfo := make(map[string]string)
	err := m.db.Table(common.BKTableNameSystem).Find(filter).Fields("token").One(ctx, &tokenInfo)
	if err != nil {
		if !m.db.IsNotFoundError(err) {
			blog.Errorf("get event start watch token by filter: %+v failed, err: %v", filter, err)
			return "", err
		}

		return "", nil
	}

	return tokenInfo["token"], nil
}

// ResetWatchToken reset watch token and start watch time
func (m *SingleHandler) ResetWatchToken(startAtTime types.TimeStamp) error {
	data := mapstr.MapStr{
		common.BKTokenField:       "",
		common.BKStartAtTimeField: startAtTime,
	}

	filter := map[string]interface{}{
		"_id": m.key,
	}

	if err := m.db.Table(common.BKTableNameSystem).Upsert(context.Background(), filter, data); err != nil {
		blog.Errorf("reset single watch token %s failed, data: %+v, err: %v", m.key, data, err)
		return err
	}
	return nil
}

// GetStartWatchTime get event start watch time
func (m *SingleHandler) GetStartWatchTime(ctx context.Context) (*types.TimeStamp, error) {
	filter := map[string]interface{}{
		"_id": m.key,
	}

	tokenInfo := make(map[string]metadata.Time)
	if err := m.db.Table(common.BKTableNameSystem).Find(filter).Fields("start_at_time").
		One(ctx, &tokenInfo); err != nil {

		if !m.db.IsNotFoundError(err) {
			blog.Errorf("get event start watch time by filter: %+v failed, err: %v", filter, err)
			return nil, err
		}

		blog.Infof("event %s start watch time is not found", m.key)
		return new(types.TimeStamp), nil
	}

	time := tokenInfo["start_at_time"].Time

	return &types.TimeStamp{
		Sec:  uint32(time.Unix()),
		Nano: uint32(time.Nanosecond()),
	}, nil
}

// IsTokenExists check if event token exists
func (m *SingleHandler) IsTokenExists(ctx context.Context) (bool, error) {
	filter := map[string]interface{}{
		"_id": m.key,
	}

	cnt, err := m.db.Table(common.BKTableNameSystem).Find(filter).Fields("token").Count(ctx)
	if err != nil {
		blog.Errorf("check if event token exists failed, filter: %+v, err: %v", filter, err)
		return false, err
	}

	return cnt > 0, nil
}
