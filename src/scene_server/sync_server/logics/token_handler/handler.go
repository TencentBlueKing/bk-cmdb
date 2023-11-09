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

// Package tokenhandler defines the common token handler for incremental sync logics using watch
package tokenhandler

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/stream/types"
)

var _ types.TokenHandler = new(TokenHandler)

// TokenHandler is the token handler that manges watch token data
type TokenHandler struct {
	key     string
	watchDB dal.DB
	metrics *event.EventMetrics
}

// New creates a TokenHandler
func New(key string, watchDB dal.DB, metrics *event.EventMetrics) (*TokenHandler, error) {
	// create watch token table if not exists
	exists, err := watchDB.HasTable(context.Background(), WatchTokenTable)
	if err != nil {
		blog.Errorf("check if table %s exists failed, err: %v", WatchTokenTable, err)
		return nil, err
	}

	if !exists {
		err = watchDB.CreateTable(context.Background(), WatchTokenTable)
		if err != nil && !watchDB.IsDuplicatedError(err) {
			blog.Errorf("create table %s failed, err: %v", WatchTokenTable, err)
			return nil, err
		}
	}

	return &TokenHandler{
		key:     key,
		watchDB: watchDB,
		metrics: metrics,
	}, nil
}

// InitWatchToken initialize watch token data
func (t TokenHandler) InitWatchToken(ctx context.Context) error {
	data := map[string]interface{}{
		common.MongoMetaID:        t.key,
		common.BKTokenField:       "",
		common.BKStartAtTimeField: time.Now(),
	}

	if err := t.watchDB.Table(WatchTokenTable).Insert(ctx, data); err != nil {
		blog.Errorf("set %s last watch token failed, err: %v, data: %+v", t.key, err, data)
		return err
	}

	return t.SetLastWatchTokenData(context.Background(), data)
}

// SetLastWatchToken set last watch token
func (t TokenHandler) SetLastWatchToken(ctx context.Context, token string) error {
	data := map[string]interface{}{
		common.BKTokenField: token,
	}
	return t.SetLastWatchTokenData(ctx, data)
}

// SetLastWatchTokenData set last watch token info
func (t TokenHandler) SetLastWatchTokenData(ctx context.Context, data map[string]interface{}) error {
	filter := map[string]interface{}{
		common.MongoMetaID: t.key,
	}
	if err := t.watchDB.Table(WatchTokenTable).Update(ctx, filter, data); err != nil {
		blog.Errorf("set %s last watch token failed, err: %v, data: %+v", t.key, err, data)
		return err
	}
	return nil
}

// GetStartWatchToken get start watch token from watch token db
func (t TokenHandler) GetStartWatchToken(ctx context.Context) (string, error) {
	filter := map[string]interface{}{
		common.MongoMetaID: t.key,
	}

	data := new(WatchToken)
	err := t.watchDB.Table(WatchTokenTable).Find(filter).Fields(common.BKTokenField).One(ctx, data)
	if err != nil {
		t.metrics.CollectMongoError()
		blog.Errorf("get %s start watch token failed, err: %v, filter: %+v", t.key, err, filter)
		return "", err
	}

	return data.Token, nil
}

// ResetWatchToken set watch token to empty and set the start watch time to the given one for next watch
func (t TokenHandler) ResetWatchToken(startAtTime types.TimeStamp) error {
	data := map[string]interface{}{
		common.BKTokenField:       "",
		common.BKStartAtTimeField: startAtTime,
	}

	return t.SetLastWatchTokenData(context.Background(), data)
}

// GetStartWatchTime get start watch token data from watch token db
func (t TokenHandler) GetStartWatchTime(ctx context.Context) (bool, types.TimeStamp, error) {
	filter := map[string]interface{}{
		common.MongoMetaID: t.key,
	}

	data := new(WatchToken)
	err := t.watchDB.Table(WatchTokenTable).Find(filter).Fields(common.BKStartAtTimeField).One(ctx, data)
	if err != nil {
		if t.watchDB.IsNotFoundError(err) {
			return false, types.TimeStamp{}, nil
		}

		t.metrics.CollectMongoError()
		blog.Errorf("get %s start watch token data failed, err: %v, filter: %+v", t.key, err, filter)
		return false, types.TimeStamp{}, err
	}

	return true, data.StartAtTime, nil
}
