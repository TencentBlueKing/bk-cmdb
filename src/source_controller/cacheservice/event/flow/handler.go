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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/stream/types"
)

var _ = types.TokenHandler(&flowTokenHandler{})

type flowTokenHandler struct {
	key     event.Key
	watchDB dal.DB
	metrics *event.EventMetrics
}

// NewFlowTokenHandler new flow token handler
func NewFlowTokenHandler(key event.Key, watchDB dal.DB, metrics *event.EventMetrics) *flowTokenHandler {
	return &flowTokenHandler{
		key:     key,
		watchDB: watchDB,
		metrics: metrics,
	}
}

/* SetLastWatchToken do not set last watch token in the do batch action(set it after events are successfully inserted)
   when there are several masters watching db event, we use db transaction to avoid inserting duplicate data by setting
   the last token after the insertion of db chain nodes in one transaction, since we have a unique index on the cursor
   field, the later one will encounters an error when inserting nodes and roll back without setting the token and watch
   another round from the last token of the last inserted node, thus ensures the sequence of db chain nodes.
*/
func (f *flowTokenHandler) SetLastWatchToken(ctx context.Context, token string) error {
	return nil
}

// setLastWatchToken set last watch token(used after events are successfully inserted)
func (f *flowTokenHandler) setLastWatchToken(ctx context.Context, data map[string]interface{}) error {
	filter := map[string]interface{}{
		"_id": f.key.Collection(),
	}
	if err := f.watchDB.Table(common.BKTableNameWatchToken).Update(ctx, filter, data); err != nil {
		blog.Errorf("set last watch token failed, err: %v, data: %+v", err, data)
		return err
	}
	return nil
}

// GetStartWatchToken get start watch token from watch token db first, if an error occurred, get from chain db
func (f *flowTokenHandler) GetStartWatchToken(ctx context.Context) (token string, err error) {
	filter := map[string]interface{}{
		"_id": f.key.Collection(),
	}

	data := new(watch.LastChainNodeData)
	err = f.watchDB.Table(common.BKTableNameWatchToken).Find(filter).Fields(common.BKTokenField).One(ctx, data)
	if err != nil {
		if !f.watchDB.IsNotFoundError(err) {
			f.metrics.CollectMongoError()
			blog.ErrorJSON("run flow, but get start watch token failed, err: %v, filter: %+v", err, filter)
		}

		tailNode := new(watch.ChainNode)
		if err := f.watchDB.Table(f.key.ChainCollection()).Find(map[string]interface{}{}).Fields(common.BKTokenField).
			Sort(common.BKFieldID+":-1").One(context.Background(), tailNode); err != nil {

			if !f.watchDB.IsNotFoundError(err) {
				f.metrics.CollectMongoError()
				blog.Errorf("get last watch token from mongo failed, err: %v", err)
				return "", err
			}
			// the tail node is not exist.
			return "", nil
		}
		return tailNode.Token, nil
	}

	return data.Token, nil
}

// resetWatchToken set watch token to empty and set the start watch time to the given one for next watch
func (f *flowTokenHandler) resetWatchToken(startAtTime types.TimeStamp) error {
	data := map[string]interface{}{
		common.BKTokenField:       "",
		common.BKStartAtTimeField: startAtTime,
	}

	filter := map[string]interface{}{
		"_id": f.key.Collection(),
	}

	if err := f.watchDB.Table(common.BKTableNameWatchToken).Update(context.Background(), filter, data); err != nil {
		blog.ErrorJSON("clear watch token failed, err: %s, collection: %s, data: %s", err, f.key.Collection(), data)
		return err
	}
	return nil
}

func (f *flowTokenHandler) getStartWatchTime(ctx context.Context) (*types.TimeStamp, error) {
	filter := map[string]interface{}{
		"_id": f.key.Collection(),
	}

	data := new(watch.LastChainNodeData)
	err := f.watchDB.Table(common.BKTableNameWatchToken).Find(filter).Fields(common.BKStartAtTimeField).One(ctx, data)
	if err != nil {
		if !f.watchDB.IsNotFoundError(err) {
			f.metrics.CollectMongoError()
			blog.ErrorJSON("run flow, but get start watch time failed, err: %v, filter: %+v", err, filter)
			return nil, err
		}
		return new(types.TimeStamp), nil
	}
	return &data.StartAtTime, nil
}
