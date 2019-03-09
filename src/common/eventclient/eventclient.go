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

package eventclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"

	"gopkg.in/redis.v5"
)

type Client interface {
	Push(context.Context, ...*metadata.EventInst) error
}

func NewEventWithHeader(header http.Header) *metadata.EventInst {
	return &metadata.EventInst{
		OwnerID:     util.GetOwnerID(header),
		TxnID:       util.GetHTTPCCTransaction(header),
		RequestID:   util.GetHTTPCCRequestID(header),
		RequestTime: metadata.Now(),
		ActionTime:  metadata.Now(),
	}
}

type ClientViaRedis struct {
	rdb       dal.RDB
	cache     *redis.Client
	queue     chan *eventtmp
	pending   *eventtmp
	queueLock sync.Mutex
}

func NewClientViaRedis(cache *redis.Client, rdb dal.RDB) *ClientViaRedis {
	// we limit the queue size to 4k*2500=10M， assume that 4k per event
	const queuesize = 2500

	ec := &ClientViaRedis{
		rdb:   rdb,
		cache: cache,
		queue: make(chan *eventtmp, queuesize),
	}
	go ec.runPusher()
	return ec
}

func (c *ClientViaRedis) Push(ctx context.Context, events ...*metadata.EventInst) error {
	c.queueLock.Lock()
	for i := range events {
		if events[i] != nil {
			eventID, err := c.rdb.NextSequence(ctx, common.EventCacheEventIDKey)
			if err != nil {
				c.queueLock.Unlock()
				return fmt.Errorf("[event] generate eventID failed: %v", err)
			}
			events[i].ID = int64(eventID)

			value, err := json.Marshal(events[i])
			if err != nil {
				c.queueLock.Unlock()
				return fmt.Errorf("[event] marshal json error: %v, raw: %#v", err, events[i])
			}
			et := &eventtmp{EventInst: events[i], data: value}
			select {
			case c.queue <- et:
			default:
				// channel fulled, so we drop 200 oldest events from queue
				// TODO save to disk if possible
				c.pending = nil
				var ok bool
				for i := 0; i < 200; i-- {
					_, ok = <-c.queue
					if !ok {
						break
					}
				}
				c.queue <- et
			}

		} else {
			c.queueLock.Unlock()
			return fmt.Errorf("[event] event could not be nil")
		}
	}
	c.queueLock.Unlock()
	return nil
}

type eventtmp struct {
	*metadata.EventInst
	data []byte
}

func (c *ClientViaRedis) runPusher() {
	var err error
	for {
		// 1. get el
		var event *eventtmp
		if c.pending != nil {
			c.queueLock.Lock()
			event = c.pending
			c.queueLock.Unlock()
		} else {
			event = <-c.queue
		}

		// 2. ignore if el is nil
		if event == nil {
			continue
		}

		// 3. push event to reids
		if err = c.pushToRedis(event); err != nil {
			c.queueLock.Lock()
			c.pending = event
			c.queueLock.Unlock()
			blog.Errorf("[event] push event to redis failed: %v, we will retry 3s later", err)
			time.Sleep(time.Second * 3)
			continue
		}

		// 4. clear pushed el
		c.queueLock.Lock()
		c.pending = nil
		c.queueLock.Unlock()
	}
}

func (c *ClientViaRedis) pushToRedis(event *eventtmp) error {
	if event.TxnID != "" {
		z := redis.Z{
			Score:  float64(time.Now().UTC().Unix()),
			Member: event.TxnID,
		}
		if err := c.cache.ZAddNX(common.EventCacheEventTxnSet, z).Err(); err != nil {
			return err
		}
		return c.cache.LPush(common.EventCacheEventTxnQueuePrefix+event.TxnID, event.data).Err()
	}
	return c.cache.LPush(common.EventCacheEventQueueKey, event.data).Err()
}
