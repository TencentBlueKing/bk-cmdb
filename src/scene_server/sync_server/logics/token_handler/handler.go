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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"
)

var _ types.TaskTokenHandler = new(TokenHandler)

// TokenHandler is the token handler that manges watch token data
type TokenHandler struct {
	key       string
	metrics   *event.EventMetrics
	startTime *types.TimeStamp
}

// InitWatchTokenTable initialize watch token table if not exists
func InitWatchTokenTable(ctx context.Context, watchDal dal.Dal) error {
	return watchDal.ExecForAllDB(func(watchDB local.DB) error {
		// create watch token table if not exists
		exists, err := watchDB.HasTable(ctx, common.BKTableNameSyncWatchToken)
		if err != nil {
			blog.Errorf("check if table %s exists failed, err: %v", common.BKTableNameSyncWatchToken, err)
			return err
		}

		if !exists {
			err = watchDB.CreateTable(ctx, common.BKTableNameSyncWatchToken)
			if err != nil && !watchDB.IsDuplicatedError(err) {
				blog.Errorf("create table %s failed, err: %v", common.BKTableNameSyncWatchToken, err)
				return err
			}
		}
		return nil
	})
}

// New creates a TokenHandler
func New(key string, metrics *event.EventMetrics) (*TokenHandler, error) {
	return &TokenHandler{
		key:       key,
		metrics:   metrics,
		startTime: &types.TimeStamp{Sec: uint32(time.Now().Unix())},
	}, nil
}

// SetLastWatchToken set last watch token
func (t TokenHandler) SetLastWatchToken(ctx context.Context, uuid string, watchDB local.DB,
	token *types.TokenInfo) error {

	filter := map[string]interface{}{
		common.MongoMetaID: watch.GenDBWatchTokenID(uuid, t.key),
	}

	tokenData := mapstr.MapStr{
		common.BKTokenField:       token.Token,
		common.BKStartAtTimeField: token.StartAtTime,
	}
	t.startTime = token.StartAtTime

	if err := watchDB.Table(common.BKTableNameSyncWatchToken).Update(ctx, filter, tokenData); err != nil {
		blog.Errorf("set %s db %s last watch token failed, err: %v, data: %+v", t.key, uuid, err, tokenData)
		return err
	}
	return nil
}

// GetStartWatchToken get start watch token from watch token db
func (t TokenHandler) GetStartWatchToken(ctx context.Context, uuid string, watchDB local.DB) (*types.TokenInfo, error) {
	filter := map[string]interface{}{
		common.MongoMetaID: watch.GenDBWatchTokenID(uuid, t.key),
	}

	data := new(types.TokenInfo)
	err := watchDB.Table(common.BKTableNameSyncWatchToken).Find(filter).Fields(common.BKTokenField,
		common.BKStartAtTimeField).One(ctx, data)
	if err != nil {
		// token info not exists, init token info with start time
		if !mongodb.IsNotFoundError(err) {
			blog.Errorf("get %s db %s start watch token failed, err: %v, filter: %+v", t.key, uuid, err, filter)
			return nil, err
		}

		tokenData := mapstr.MapStr{
			common.MongoMetaID:        watch.GenDBWatchTokenID(uuid, t.key),
			common.BKTokenField:       "",
			common.BKStartAtTimeField: t.startTime,
		}

		if err = watchDB.Table(common.BKTableNameSyncWatchToken).Insert(ctx, tokenData); err != nil {
			blog.Errorf("create %s db %s watch token failed, err: %v, data: %+v", t.key, uuid, err, tokenData)
			return nil, err
		}

		return &types.TokenInfo{StartAtTime: t.startTime}, nil
	}

	return data, nil
}
