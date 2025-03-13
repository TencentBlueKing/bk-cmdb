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

package flow

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"
)

type flowTokenHandler struct {
	key     event.Key
	metrics *event.EventMetrics
}

// NewFlowTokenHandler new flow token handler
func NewFlowTokenHandler(key event.Key, metrics *event.EventMetrics) *flowTokenHandler {
	return &flowTokenHandler{
		key:     key,
		metrics: metrics,
	}
}

/*
SetLastWatchToken do not set last watch token in the do batch action(set it after events are successfully inserted)
when there are several masters watching db event, we use db transaction to avoid inserting duplicate data by setting
the last token after the insertion of db chain nodes in one transaction, since we have a unique index on the cursor
field, the later one will encounter an error when inserting nodes and roll back without setting the token and watch
another round from the last token of the last inserted node, thus ensures the sequence of db chain nodes.
*/
func (f *flowTokenHandler) SetLastWatchToken(_ context.Context, _ string, _ local.DB, _ *types.TokenInfo) error {
	return nil
}

// setLastWatchToken set last watch token(used after events are successfully inserted)
func (f *flowTokenHandler) setLastWatchToken(ctx context.Context, uuid string, watchDB local.DB,
	data map[string]any) error {

	filter := map[string]interface{}{
		"_id": watch.GenDBWatchTokenID(uuid, f.key.Collection()),
	}
	if err := watchDB.Table(common.BKTableNameWatchToken).Update(ctx, filter, data); err != nil {
		blog.Errorf("set last watch token failed, err: %v, data: %+v", err, data)
		return err
	}
	return nil
}

// GetStartWatchToken get start watch token from watch token db first, if an error occurred, get from chain db
func (f *flowTokenHandler) GetStartWatchToken(ctx context.Context, uuid string, watchDB local.DB) (*types.TokenInfo,
	error) {

	filter := map[string]interface{}{
		"_id": watch.GenDBWatchTokenID(uuid, f.key.Collection()),
	}

	data := new(types.TokenInfo)
	err := watchDB.Table(common.BKTableNameWatchToken).Find(filter).Fields(common.BKTokenField,
		common.BKStartAtTimeField).One(ctx, data)
	if err != nil {
		if !mongodb.IsNotFoundError(err) {
			f.metrics.CollectMongoError()
			blog.Errorf("run flow, but get start watch token failed, err: %v, filter: %+v", err, filter)
		}
		// the tail node is not exist.
		return &types.TokenInfo{Token: "", StartAtTime: &types.TimeStamp{Sec: uint32(time.Now().Unix())}}, nil
	}

	return data, nil
}
