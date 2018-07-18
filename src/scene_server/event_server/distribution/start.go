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
	redis "gopkg.in/redis.v5"

	"configcenter/src/scene_server/event_server/identifier"
	"configcenter/src/storage"
)

func Start(cache *redis.Client, db storage.DI) error {
	chErr := make(chan error)

	eh := &EventHandler{cache: cache}
	go func() {
		chErr <- eh.StartHandleInsts()
	}()

	dh := &DistHandler{cache: cache, db: db}
	go func() {
		chErr <- dh.StartDistribute()
	}()

	ih := identifier.NewIdentifierHandler(cache, db)
	go func() {
		chErr <- ih.StartHandleInsts()
	}()

	return <-chErr
}

type EventHandler struct{ cache *redis.Client }
type DistHandler struct {
	cache *redis.Client
	db    storage.DI
}
