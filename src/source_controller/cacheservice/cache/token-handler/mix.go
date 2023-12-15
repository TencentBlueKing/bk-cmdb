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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/stream/types"
)

var _ types.TokenHandler = new(MixHandler)

// MixHandler is a token handler for mix event composed of multiple types of events
// token data: {"_id": $mixKey, $collection: {"token": $token, "start_at_time": $start_at_time}}
type MixHandler struct {
	mixKey     string
	collection string
	db         dal.DB
}

// NewMixTokenHandler generate a new mix event token handler
func NewMixTokenHandler(mixKey, collection string, db dal.DB) *MixHandler {
	return &MixHandler{
		mixKey:     mixKey,
		collection: collection,
		db:         db,
	}
}

// SetLastWatchToken set last mix event watch token
func (m *MixHandler) SetLastWatchToken(ctx context.Context, token string) error {
	filter := map[string]interface{}{
		"_id": m.mixKey,
	}

	tokenInfo := mapstr.MapStr{
		m.collection + ".token": token,
	}

	if err := m.db.Table(common.BKTableNameSystem).Upsert(ctx, filter, tokenInfo); err != nil {
		blog.Errorf("set mix event %s last watch token failed, data: %+v, err: %v", m.mixKey, tokenInfo, err)
		return err
	}
	return nil
}

// GetStartWatchToken get mix event start watch token
func (m *MixHandler) GetStartWatchToken(ctx context.Context) (string, error) {
	filter := map[string]interface{}{
		"_id": m.mixKey,
	}

	data := make(map[string]map[string]string)
	err := m.db.Table(common.BKTableNameSystem).Find(filter).Fields(m.collection+".token").One(ctx, &data)
	if err != nil {
		if !m.db.IsNotFoundError(err) {
			blog.Errorf("get mix event start watch token by filter: %+v failed, err: %v", filter, err)
			return "", err
		}

		return "", nil
	}

	tokenInfo, exist := data[m.collection]
	if !exist {
		blog.Infof("mix event %s start watch token is not found", m.mixKey)
		return "", nil
	}

	return tokenInfo["token"], nil
}

// ResetWatchToken reset watch token and start watch time
func (m *MixHandler) ResetWatchToken(startAtTime types.TimeStamp) error {
	data := mapstr.MapStr{
		m.collection: mapstr.MapStr{
			common.BKTokenField:       "",
			common.BKStartAtTimeField: startAtTime,
		},
	}

	filter := map[string]interface{}{
		"_id": m.mixKey,
	}

	if err := m.db.Table(common.BKTableNameSystem).Upsert(context.Background(), filter, data); err != nil {
		blog.Errorf("reset mix watch token %s failed, data: %+v, err: %v", m.mixKey, data, err)
		return err
	}
	return nil
}

// GetStartWatchTime get mix event start watch time
func (m *MixHandler) GetStartWatchTime(ctx context.Context) (*types.TimeStamp, error) {
	filter := map[string]interface{}{
		"_id": m.mixKey,
	}

	data := make(map[string]map[string]metadata.Time)
	if err := m.db.Table(common.BKTableNameSystem).Find(filter).Fields(m.collection+".start_at_time").
		One(ctx, &data); err != nil {

		if !m.db.IsNotFoundError(err) {
			blog.Errorf("get mix event start watch time by filter: %+v failed, err: %v", filter, err)
			return nil, err
		}

		blog.Infof("mix event %s start watch time is not found", m.mixKey)
		return new(types.TimeStamp), nil
	}

	tokenInfo, exist := data[m.collection]
	if !exist {
		blog.Infof("mix event %s start watch time is not found", m.mixKey)
		return new(types.TimeStamp), nil
	}

	time := tokenInfo["start_at_time"].Time

	return &types.TimeStamp{
		Sec:  uint32(time.Unix()),
		Nano: uint32(time.Nanosecond()),
	}, nil
}
