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

package mixevent

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

var _ = types.TokenHandler(&mixEventHandler{})

// mixEventHandler a token handler for mix event consisting of multiple types of events, stores the events' token in the
// form of {"_id": $mix_key_collection, {$event_key_collection: {"token": $token, "start_at_time": $start_at_time}}... }
type mixEventHandler struct {
	mixKey  event.Key
	key     event.Key
	watchDB dal.DB
	metrics *event.EventMetrics
}

// newMixEventTokenHandler generate a new mix event token handler
func newMixEventTokenHandler(mixKey event.Key, key event.Key, watchDB dal.DB,
	metrics *event.EventMetrics) *mixEventHandler {

	return &mixEventHandler{
		mixKey:  mixKey,
		key:     key,
		watchDB: watchDB,
		metrics: metrics,
	}
}

// SetLastWatchToken TODO
/* SetLastWatchToken do not use this function in the mix events(set after events are successfully inserted)
   when there are several masters watching db event, we use db transaction to avoid inserting duplicate data by setting
   the last token after the insertion of db chain nodes in one transaction, since we have a unique index on the cursor
   field, the later one will encounters an error when inserting nodes and roll back without setting the token and watch
   another round from the last token of the last inserted node, thus ensures the sequence of db chain nodes.
*/
func (m *mixEventHandler) SetLastWatchToken(ctx context.Context, token string) error {
	return nil
}

// setLastWatchToken set last watch token(used after events are successfully inserted)
func (m *mixEventHandler) setLastWatchToken(ctx context.Context, data map[string]interface{}) error {
	filter := map[string]interface{}{
		"_id": m.mixKey.Collection(),
	}

	// only update the needed fields to avoid erasing the previous exist fields
	tokenInfo := make(mapstr.MapStr)
	for key, value := range data {
		tokenInfo[m.key.Collection()+"."+key] = value
	}

	// update id and cursor field if set, to compensate for the scenario of searching with an outdated but latest cursor
	if id, exists := data[common.BKFieldID]; exists {
		tokenInfo[common.BKFieldID] = id
	}

	if cursor, exists := data[common.BKCursorField]; exists {
		tokenInfo[common.BKCursorField] = cursor
	}

	if err := m.watchDB.Table(common.BKTableNameWatchToken).Update(ctx, filter, tokenInfo); err != nil {
		blog.Errorf("set mix event %s last watch token failed, data: %+v, err: %v", m.key.Collection(), tokenInfo, err)
		return err
	}
	return nil
}

// GetStartWatchToken get start watch token from watch token db first, if an error occurred, get from chain db
func (m *mixEventHandler) GetStartWatchToken(ctx context.Context) (token string, err error) {
	filter := map[string]interface{}{
		"_id": m.mixKey.Collection(),
	}

	data := make(map[string]watch.LastChainNodeData)
	if err := m.watchDB.Table(common.BKTableNameWatchToken).Find(filter).Fields(m.key.Collection()).
		One(ctx, &data); err != nil {
		if !m.watchDB.IsNotFoundError(err) {
			m.metrics.CollectMongoError()
			blog.Errorf("get mix event start watch token, will get the last event's time and start watch, "+
				"filter: %+v, err: %v", filter, err)
		}

		tailNode := new(watch.ChainNode)
		if err := m.watchDB.Table(m.mixKey.ChainCollection()).Find(nil).Fields(common.BKTokenField).
			Sort(common.BKFieldID+":-1").One(context.Background(), tailNode); err != nil {

			if !m.watchDB.IsNotFoundError(err) {
				m.metrics.CollectMongoError()
				blog.Errorf("get mix event last watch token from mongo failed, err: %v", err)
				return "", err
			}
			// the tail node is not exist.
			return "", nil
		}
		return tailNode.Token, nil
	}

	// check whether this field is exists or not
	node, exists := data[m.key.Collection()]
	if !exists {
		// watch from now on.
		return "", nil
	}

	return node.Token, nil
}

// resetWatchToken set watch token to empty and set the start watch time to the given one for next watch
func (m *mixEventHandler) resetWatchToken(startAtTime types.TimeStamp) error {
	data := mapstr.MapStr{
		m.key.Collection(): mapstr.MapStr{
			common.BKTokenField:       "",
			common.BKStartAtTimeField: startAtTime,
		},
	}

	filter := map[string]interface{}{
		"_id": m.mixKey.Collection(),
	}

	if err := m.watchDB.Table(common.BKTableNameWatchToken).Update(context.Background(), filter, data); err != nil {
		blog.Errorf("clear watch token failed, collection: %s, data: %+v, err: %v", m.key.Collection(), data, err)
		return err
	}
	return nil
}

// getStartWatchTime get start watch tim of the event key in mix event token
func (m *mixEventHandler) getStartWatchTime(ctx context.Context) (*types.TimeStamp, error) {
	filter := map[string]interface{}{
		"_id": m.mixKey.Collection(),
	}

	data := make(map[string]watch.LastChainNodeData)
	err := m.watchDB.Table(common.BKTableNameWatchToken).Find(filter).Fields(m.key.Collection()).One(ctx, &data)
	if err != nil {
		blog.Errorf("get mix event %s start watch time, but find in watch token failed, err: %v", m.key.Collection(),
			err)

		if !m.watchDB.IsNotFoundError(err) {
			m.metrics.CollectMongoError()
			blog.Errorf("run flow, but get start watch time failed, filter: %+v, err: %v", filter, err)
			return nil, err
		}

		blog.Infof("get mix event %s start watch time, but not find in watch token, start watch from a minute ago",
			m.key.Collection())
		return new(types.TimeStamp), nil
	}

	node, exist := data[m.key.Collection()]
	if !exist {
		// can not find, start watch from one minute ago.
		blog.Infof("get mix event %s start watch time, but not find in watch token, start watch from a minute ago",
			m.key.Collection())
		return new(types.TimeStamp), nil
	}

	return &node.StartAtTime, nil
}
