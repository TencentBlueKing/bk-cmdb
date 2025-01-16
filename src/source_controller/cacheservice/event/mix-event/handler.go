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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"
)

// mixEventHandler a token handler for mix event consisting of multiple types of events, stores the events' token in the
// form of {"_id": $mix_key_collection, {$event_key_collection: {"token": $token, "start_at_time": $start_at_time}}... }
type mixEventHandler struct {
	mixKey  event.Key
	key     event.Key
	metrics *event.EventMetrics
}

// newMixEventTokenHandler generate a new mix event token handler
func newMixEventTokenHandler(mixKey event.Key, key event.Key, metrics *event.EventMetrics) *mixEventHandler {
	return &mixEventHandler{
		mixKey:  mixKey,
		key:     key,
		metrics: metrics,
	}
}

/*
SetLastWatchToken do not use this function in the mix events(set after events are successfully inserted)
when there are several masters watching db event, we use db transaction to avoid inserting duplicate data by setting
the last token after the insertion of db chain nodes in one transaction, since we have a unique index on the cursor
field, the later one will encounter an error when inserting nodes and roll back without setting the token and watch
another round from the last token of the last inserted node, thus ensures the sequence of db chain nodes.
*/
func (m *mixEventHandler) SetLastWatchToken(_ context.Context, _ string, _ local.DB, _ *types.TokenInfo) error {
	return nil
}

// setLastWatchToken set last watch token(used after events are successfully inserted)
func (m *mixEventHandler) setLastWatchToken(ctx context.Context, uuid string, watchDB local.DB,
	data map[string]any) error {

	filter := map[string]interface{}{
		"_id": watch.GenDBWatchTokenID(uuid, m.mixKey.Collection()),
	}

	// only update the needed fields to avoid erasing the previous exist fields
	tokenInfo := make(mapstr.MapStr)
	for key, value := range data {
		tokenInfo[m.key.Collection()+"."+key] = value
	}

	if err := watchDB.Table(common.BKTableNameWatchToken).Update(ctx, filter, tokenInfo); err != nil {
		blog.Errorf("set mix event %s coll %s last watch token failed, data: %+v, err: %v", m.mixKey.Namespace(),
			m.key.Collection(), tokenInfo, err)
		return err
	}
	return nil
}

// GetStartWatchToken get start watch token from watch token db first, if an error occurred, get from chain db
func (m *mixEventHandler) GetStartWatchToken(ctx context.Context, uuid string, watchDB local.DB) (*types.TokenInfo,
	error) {

	filter := map[string]interface{}{
		"_id": watch.GenDBWatchTokenID(uuid, m.mixKey.Collection()),
	}

	data := make(map[string]types.TokenInfo)
	if err := watchDB.Table(common.BKTableNameWatchToken).Find(filter).Fields(m.key.Collection()).
		One(ctx, &data); err != nil {
		if !mongodb.IsNotFoundError(err) {
			m.metrics.CollectMongoError()
			blog.Errorf("get mix event start watch token, will get the last event's time and start watch, "+
				"filter: %+v, err: %v", filter, err)
		}

		// the tail node is not exist.
		return &types.TokenInfo{Token: "", StartAtTime: &types.TimeStamp{Sec: uint32(time.Now().Unix())}}, nil
	}

	// check whether this field is exists or not
	node, exists := data[m.key.Collection()]
	if !exists {
		// watch from now on.
		return &types.TokenInfo{Token: "", StartAtTime: &types.TimeStamp{Sec: uint32(time.Now().Unix())}}, nil
	}

	return &node, nil
}
