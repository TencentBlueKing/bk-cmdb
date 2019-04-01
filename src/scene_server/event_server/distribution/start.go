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
	"sync"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/identifier"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/rpc"

	redis "gopkg.in/redis.v5"
)

func Start(ctx context.Context, cache *redis.Client, db dal.RDB, rc *rpc.Client) error {
	chErr := make(chan error)

	eh := &EventHandler{cache: cache}
	go func() {
		chErr <- eh.StartHandleInsts()
	}()

	dh := &DistHandler{cache: cache, db: db, ctx: ctx}
	go func() {
		chErr <- dh.StartDistribute()
	}()

	ih := identifier.NewIdentifierHandler(ctx, cache, db)
	go func() {
		chErr <- ih.StartHandleInsts()
	}()

	th := &TxnHandler{cache: cache, db: db, ctx: ctx, rc: rc, committed: make(chan string, 100), shouldClose: util.NewBool(false)}
	go func() {
		for {
			if err := th.Run(); err != nil {
				blog.Errorf("TxnHandler stoped with error: %v, we will try 1s later", err)
			}
			time.Sleep(time.Second)
		}
	}()

	return <-chErr
}

type EventHandler struct{ cache *redis.Client }
type DistHandler struct {
	cache *redis.Client
	db    dal.RDB
	ctx   context.Context
}

type TxnHandler struct {
	rc          *rpc.Client
	cache       *redis.Client
	db          dal.RDB
	ctx         context.Context
	committed   chan string
	shouldClose *util.AtomicBool
	wg          sync.WaitGroup
}
