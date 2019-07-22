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

package distribution

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/identifier"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/rpc"
)

func Start(ctx context.Context, cache *redis.Client, db dal.RDB, rc rpc.Client) error {
	chErr := make(chan error, 1)
	err := migrateIDToMongo(ctx, cache, db)
	if err != nil {
		return fmt.Errorf("migrateIDToMongo failed: %v", err)
	}

	eh := &EventHandler{cache: cache}
	go func() {
		chErr <- eh.Run()
	}()

	dh := &DistHandler{cache: cache, db: db, ctx: ctx}
	go func() {
		chErr <- dh.StartDistribute()
	}()

	ih := identifier.NewIdentifierHandler(ctx, cache, db)
	go func() {
		chErr <- ih.Run()
	}()

	go cleanExpiredEvents(cache)

	if rc != nil {
		th := &TxnHandler{cache: cache, db: db, ctx: ctx, rc: rc, committed: make(chan string, 100), shouldClose: util.NewBool(false)}
		go func() {
			for {
				if err := th.Run(); err != nil {
					blog.Errorf("TxnHandler stopped with error: %v, we will try 1s later", err)
				}
				time.Sleep(time.Second)
			}
		}()
	}

	return <-chErr
}

func migrateIDToMongo(ctx context.Context, cache *redis.Client, db dal.RDB) error {
	sid, err := cache.Get(common.EventCacheEventIDKey).Result()
	if redis.Nil == err {
		return nil
	}
	if err != nil {
		return err
	}
	if sid == "" {
		return nil
	}

	id, err := strconv.ParseUint(sid, 10, 64)
	if err != nil {
		return err
	}

	docs := map[string]interface{}{
		"_id":        common.EventCacheEventIDKey,
		"SequenceID": id,
	}

	err = db.Table(common.BKTableNameIDgenerator).Insert(ctx, docs)
	if err != nil && !db.IsDuplicatedError(err) {
		return err
	}

	return cache.Del(common.EventCacheEventIDKey).Err()
}

type EventHandler struct{ cache *redis.Client }
type DistHandler struct {
	cache *redis.Client
	db    dal.RDB
	ctx   context.Context
}

type TxnHandler struct {
	rc          rpc.Client
	cache       *redis.Client
	db          dal.RDB
	ctx         context.Context
	committed   chan string
	shouldClose *util.AtomicBool
	wg          sync.WaitGroup
}
