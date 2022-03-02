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

package mainline

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/stream"
	"configcenter/src/storage/stream/types"
	"github.com/tidwall/gjson"
)

// business is a instance to watch business's change event and
// then try to refresh it to the cache.
// it based one the event loop watch mechanism which can ensure
// all the event can be watched safely, which also means the cache
// can be refreshed without lost and immediately.
type business struct {
	key   keyGenerator
	event stream.LoopInterface
	rds   redis.Client
	db    dal.DB
}

// Run start to watch and refresh the business's cache.
func (b *business) Run() error {

	// initialize business token handler key.
	handler := newTokenHandler(b.key)
	startTime, err := handler.getStartTimestamp(context.Background())
	if err != nil {
		blog.Errorf("get business cache event start at time failed, err: %v", err)
		return err
	}

	loopOpts := &types.LoopOneOptions{
		LoopOptions: types.LoopOptions{
			Name: "biz_cache",
			WatchOpt: &types.WatchOptions{
				Options: types.Options{
					EventStruct: new(map[string]interface{}),
					Collection:  common.BKTableNameBaseApp,
					// start token will be automatically set when it's running,
					// so we do not set here.
					StartAfterToken:         nil,
					StartAtTime:             startTime,
					WatchFatalErrorCallback: handler.resetWatchTokenWithTimestamp,
				},
			},
			TokenHandler: handler,
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 4,
				RetryDuration: retryDuration,
			},
		},
		EventHandler: &types.OneHandler{
			DoAdd:    b.onUpsert,
			DoUpdate: b.onUpsert,
			DoDelete: b.onDelete,
		},
	}

	return b.event.WithOne(loopOpts)
}

// onUpsert set or update business cache when a add/update/upsert
// event is triggered.
func (b *business) onUpsert(e *types.Event) bool {
	if blog.V(4) {
		blog.Infof("received biz cache event, op: %s, doc: %s, rid: %s", e.OperationType, e.DocBytes, e.ID())
	}

	bizID := gjson.GetBytes(e.DocBytes, common.BKAppIDField).Int()
	if bizID <= 0 {
		blog.Errorf("received invalid biz event, skip, op: %s, doc: %s, rid: %s", e.OperationType, e.DocBytes, e.ID())
		return false
	}

	// update the cache.
	err := b.rds.Set(context.Background(), b.key.detailKey(bizID), e.DocBytes, b.key.detailExpireDuration).Err()
	if err != nil {
		blog.Errorf("set biz cache failed, op: %s, doc: %s, err: %v, rid: %s", e.OperationType, e.DocBytes, err, e.ID())
		return true
	}

	return false
}

// onDelete delete business cache when a business s delete.
func (b *business) onDelete(e *types.Event) bool {

	filter := mapstr.MapStr{
		"coll": common.BKTableNameBaseApp,
		"oid":  e.Oid,
	}

	biz := new(bizArchive)
	err := b.db.Table(common.BKTableNameDelArchive).Find(filter).Fields("detail").One(context.Background(), biz)
	if err != nil {
		blog.Errorf("get biz del archive detail failed, err: %v, rid: %s", err, e.ID())
		if b.db.IsNotFoundError(err) {
			return false
		}
		return true
	}

	blog.Infof("received delete biz %d/%s event, rid: %s", biz.Detail.BusinessID, biz.Detail.BusinessName, e.ID())

	// delete the cache.
	if err := b.rds.Del(context.Background(), b.key.detailKey(biz.Detail.BusinessID)).Err(); err != nil {
		blog.Errorf("delete biz cache failed, err: %v, rid: %s", err, e.ID())
		return true
	}

	return false
}
