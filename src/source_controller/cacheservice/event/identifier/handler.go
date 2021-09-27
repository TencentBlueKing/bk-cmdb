/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package identifier

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/stream/types"
)

var _ = types.TokenHandler(&identityHandler{})

type identityHandler struct {
	key     event.Key
	watchDB dal.DB
	metrics *event.EventMetrics
}

func newIdentityTokenHandler(key event.Key, watchDB dal.DB, metrics *event.EventMetrics) *identityHandler {
	return &identityHandler{
		key:     key,
		watchDB: watchDB,
		metrics: metrics,
	}
}

/* SetLastWatchToken do not use this function in this host identity events(set after events are successfully inserted)
   when there are several masters watching db event, we use db transaction to avoid inserting duplicate data by setting
   the last token after the insertion of db chain nodes in one transaction, since we have a unique index on the cursor
   field, the later one will encounters an error when inserting nodes and roll back without setting the token and watch
   another round from the last token of the last inserted node, thus ensures the sequence of db chain nodes.
*/
func (f *identityHandler) SetLastWatchToken(ctx context.Context, token string) error {
	return nil
}

// setLastWatchToken set last watch token(used after events are successfully inserted)
func (f *identityHandler) setLastWatchToken(ctx context.Context, data map[string]interface{}) error {
	filter := map[string]interface{}{
		"_id": event.HostIdentityKey.Collection(),
	}

	tokenInfo := mapstr.MapStr{
		f.key.Collection(): data,
	}

	// update id and cursor field if set, to compensate for the scenario of searching with an outdated but latest cursor
	if id, exists := data[common.BKFieldID]; exists {
		tokenInfo[common.BKFieldID] = id
	}

	if cursor, exists := data[common.BKCursorField]; exists {
		tokenInfo[common.BKCursorField] = cursor
	}

	if err := f.watchDB.Table(common.BKTableNameWatchToken).Update(ctx, filter, tokenInfo); err != nil {
		blog.Errorf("set host identity %s last watch token failed, err: %v, data: %+v", f.key.Collection(),
			err, tokenInfo)
		return err
	}
	return nil
}

// GetStartWatchToken get start watch token from watch token db first, if an error occurred, get from chain db
func (f *identityHandler) GetStartWatchToken(ctx context.Context) (token string, err error) {
	filter := map[string]interface{}{
		"_id": event.HostIdentityKey.Collection(),
	}

	data := make(map[string]watch.LastChainNodeData)
	if err := f.watchDB.Table(common.BKTableNameWatchToken).Find(filter).Fields(f.key.Collection()).
		One(ctx, &data); err != nil {
		if !f.watchDB.IsNotFoundError(err) {
			f.metrics.CollectMongoError()
			blog.ErrorJSON("get host identity start watch token, will get the last event's time and start watch, "+
				"err: %s, filter: %s", err, filter)
		}

		tailNode := new(watch.ChainNode)
		if err := f.watchDB.Table(event.HostIdentityKey.ChainCollection()).Find(nil).Fields(common.BKTokenField).
			Sort(common.BKFieldID+":-1").One(context.Background(), tailNode); err != nil {

			if !f.watchDB.IsNotFoundError(err) {
				f.metrics.CollectMongoError()
				blog.Errorf("get host identity last watch token from mongo failed, err: %v", err)
				return "", err
			}
			// the tail node is not exist.
			return "", nil
		}
		return tailNode.Token, nil
	}

	// check whether this field is exists or not
	node, exists := data[f.key.Collection()]
	if !exists {
		// watch from now on.
		return "", nil
	}

	return node.Token, nil
}

// resetWatchToken set watch token to empty and set the start watch time to the given one for next watch
func (f *identityHandler) resetWatchToken(startAtTime types.TimeStamp) error {
	data := mapstr.MapStr{
		f.key.Collection(): mapstr.MapStr{
			common.BKTokenField:       "",
			common.BKStartAtTimeField: startAtTime,
		},
	}

	filter := map[string]interface{}{
		"_id": event.HostIdentityKey.Collection(),
	}

	if err := f.watchDB.Table(common.BKTableNameWatchToken).Update(context.Background(), filter, data); err != nil {
		blog.ErrorJSON("clear watch token failed, err: %s, collection: %s, data: %s", err, f.key.Collection(), data)
		return err
	}
	return nil
}

func (f *identityHandler) getStartWatchTime(ctx context.Context) (*types.TimeStamp, error) {
	filter := map[string]interface{}{
		"_id": event.HostIdentityKey.Collection(),
	}

	data := make(map[string]watch.LastChainNodeData)
	err := f.watchDB.Table(common.BKTableNameWatchToken).Find(filter).One(ctx, &data)
	if err != nil {
		blog.Errorf("get host identity %s start watch time, but find in watch token failed, err: %v",
			f.key.Collection(), err)

		if !f.watchDB.IsNotFoundError(err) {
			f.metrics.CollectMongoError()
			blog.ErrorJSON("run flow, but get start watch time failed, err: %v, filter: %+v", err, filter)
			return nil, err
		}

		blog.Infof("get host identity %s start watch time, but not find in watch token, "+
			"start watch from a minute ago", f.key.Collection())
		return new(types.TimeStamp), nil
	}

	node, exist := data[f.key.Collection()]
	if !exist {
		// can not find, start watch from one minute ago.
		blog.Infof("get host identity %s start watch time, but not find in watch token, "+
			"start watch from a minute ago", f.key.Collection())
		return new(types.TimeStamp), nil
	}

	return &node.StartAtTime, nil
}
