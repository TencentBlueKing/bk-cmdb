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

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/scene_server/event_server/identifier"
	"configcenter/src/storage/dal"

	"gopkg.in/redis.v5"
)

func Start(ctx context.Context, cache *redis.Client, db dal.RDB, clientSet apimachinery.ClientSetInterface, disc discovery.ServiceManageInterface) error {
	chErr := make(chan error, 1)

	eh := &EventHandler{cache: cache}
	go func() {
		chErr <- eh.Run()
	}()

	dh := &DistHandler{cache: cache, db: db, ctx: ctx}
	go func() {
		chErr <- dh.StartDistribute()
	}()

	ih := identifier.NewIdentifierHandler(ctx, cache, db, clientSet)
	go func() {
		chErr <- ih.Run()
	}()

	go cleanExpiredEvents(cache)
	go cleanEventQueue(cache, disc)

	th := &TxnHandler{cache: cache}
	th.Run()

	return <-chErr
}

type EventHandler struct{ cache *redis.Client }

type DistHandler struct {
	cache *redis.Client
	db    dal.RDB
	ctx   context.Context
}

type TxnHandler struct {
	cache *redis.Client
}
