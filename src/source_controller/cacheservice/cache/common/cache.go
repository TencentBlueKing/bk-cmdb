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

// Package common defines the common resource cache
package common

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/source_controller/cacheservice/cache/common/key"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/reflector"
	"configcenter/src/storage/stream/types"

	rawRedis "github.com/go-redis/redis/v7"
)

var cacheMap = make(map[key.KeyType]*commonCache)
var cacheOnce = make(map[key.KeyType]*sync.Once)

// InitCache initialize common resource cache
func InitCache(event reflector.Interface) error {
	allKeyGen := key.GetAllKeyGenerator()

	for typ, keyGen := range allKeyGen {
		_, exists := cacheMap[typ]
		if exists {
			continue
		}

		_, exists = cacheOnce[typ]
		if !exists {
			cacheOnce[typ] = new(sync.Once)
		}

		var err error
		cacheOnce[typ].Do(func() {
			cacheMap[typ] = &commonCache{
				key:   keyGen,
				event: event,
			}

			err = cacheMap[typ].Run()
		})

		if err != nil {
			return fmt.Errorf("run %s cache failed, err: %v", typ, err)
		}
	}

	return nil
}

type commonCache struct {
	key   key.KeyGenerator
	event reflector.Interface
}

// Run common resource watch
func (c *commonCache) Run() error {
	_, err := redis.Client().Get(context.Background(), c.key.ListDoneKey()).Result()
	if err != nil {
		if !redis.IsNilErr(err) {
			blog.Errorf("get %s list done redis key %s failed, err: %v", c.key.Type(), c.key.ListDoneKey(), err)
			return fmt.Errorf("get %s list done redis key failed, err: %v", err)
		}

		// do with list watcher.
		page := 500
		listOpts := &types.ListWatchOptions{
			Options:  c.key.GetWatchOpt(),
			PageSize: &page,
		}

		listCap := &reflector.Capable{
			OnChange: reflector.OnChangeEvent{
				OnLister:     c.onUpsert,
				OnAdd:        c.onUpsert,
				OnUpdate:     c.onUpsert,
				OnListerDone: c.onListDone,
				OnDelete:     c.onDelete,
			},
		}

		blog.Info("do %s cache with list watcher.", c.key.Type())
		return c.event.ListWatcher(context.Background(), listOpts, listCap)
	}

	// do with watcher only.
	watchOpts := &types.WatchOptions{
		Options: c.key.GetWatchOpt(),
	}

	watchCap := &reflector.Capable{
		OnChange: reflector.OnChangeEvent{
			OnAdd:    c.onUpsert,
			OnUpdate: c.onUpsert,
			OnDelete: c.onDelete,
		},
	}

	blog.Info("do %s cache with only watcher")
	return c.event.Watcher(context.Background(), watchOpts, watchCap)
}

func (c *commonCache) onUpsert(e *types.Event) {
	if blog.V(4) {
		blog.Infof("received %s upsert event, oid: %s, doc: %s", c.key.Type(), e.Oid, e.DocBytes)
	}

	idKey, _, err := c.key.GetIDKey(e.Document)
	if err != nil {
		blog.Errorf("generate %s id key from upsert event failed, oid: %s, doc: %s", c.key.Type(), e.Oid, e.DocBytes)
		return
	}

	// get resource details from db again to avoid dirty data.
	mongoData, err := c.key.GetMongoData(key.IDKind, mongodb.Client(), idKey)
	if err != nil {
		blog.Errorf("get %s mongo data from id key %s failed, oid: %s, doc: %s", c.key.Type(), idKey, e.Oid, e.DocBytes)
		return
	}

	for _, data := range mongoData {
		refreshCache(c.key, data, e.Oid)
	}
}

func (c *commonCache) onListDone() {
	if err := redis.Client().Set(context.Background(), c.key.ListDoneKey(), "done", 0).Err(); err != nil {
		blog.Errorf("list %s data to cache and list done, but set list done key failed, err: %v", c.key.Type(), err)
		return
	}
	blog.Info("list %s data to cache and list done")
}

func (c *commonCache) onDelete(e *types.Event) {
	blog.Infof("received %s delete event, oid: %s", e.Oid)

	filter := mapstr.MapStr{
		"oid":  e.Oid,
		"coll": e.Collection,
	}
	doc := make(map[string]mapstr.MapStr)
	err := mongodb.Client().Table(common.BKTableNameDelArchive).Find(filter).Fields("detail").One(context.Background(),
		&doc)
	if err != nil {
		blog.Errorf("get del archive failed, err: %v, oid: %s, coll: %s", err, e.Oid, e.Collection)
		return
	}

	pipe := redis.Client().Pipeline()

	// remove common resource detail cache
	idKey, _, err := c.key.GetIDKey(doc["detail"])
	if err != nil {
		blog.Errorf("generate %s id key from del archive failed, err: %v, doc: %+v", c.key.Type(), err, doc)
		return
	}
	pipe.Del(c.key.DetailKey(idKey))

	// delete common resource key kinds related redis keys
	for _, kind := range c.key.GetAllKeyKinds() {
		redisKey, err := c.key.GenerateRedisKey(kind, doc["detail"])
		if err != nil {
			blog.Errorf("generate %s %s redis key from del archive failed, err: %v, doc: %+v", c.key.Type(), kind, err,
				doc)
			return
		}

		pipe.SRem(redisKey, idKey)
	}

	// remove common resource id from common resource id list
	pipe.ZRem(c.key.IDListKey(), idKey)

	_, err = pipe.Exec()
	if err != nil {
		blog.Errorf("delete redis cache failed, err: %v, oid: %s, coll: %s", err, e.Oid, e.Collection)
		return
	}
	blog.Infof("received %s delete event, oid: %s, delete redis keys success", c.key.Type(), e.Oid)
}

// refreshCache refresh the common resource cache
func refreshCache(key key.KeyGenerator, data interface{}, rid string) {
	idKey, score, err := key.GetIDKey(data)
	if err != nil {
		blog.Errorf("generate %s refresh data id key failed, err: %v, data: %+v, rid: %s", key.Type(), err, data, rid)
		return
	}

	detailLockKey := key.DetailLockKey(idKey)

	// get refresh lock to avoid concurrency
	success, err := redis.Client().SetNX(context.Background(), detailLockKey, 1, 10*time.Second).Result()
	if err != nil {
		blog.Errorf("get %s detail lock %s failed, err: %v, rid: %s", key.Type(), detailLockKey, err, rid)
		return
	}

	if !success {
		blog.V(4).Infof("do not get %s detail lock %s, skip, rid: %s", key.Type(), detailLockKey, data, rid)
		return
	}

	defer func() {
		if err = redis.Client().Del(context.Background(), detailLockKey).Err(); err != nil {
			blog.Errorf("delete %s detail lock %s failed, err: %v, rid: %s", key.Type(), detailLockKey, err, rid)
		}
	}()

	// refresh all key kind cache after acquiring the lock
	pipeline := redis.Client().Pipeline()
	ttl := key.WithRandomExpireSeconds()

	// upsert all other key kinds' redis key to id relation cache
	for _, kind := range key.GetAllKeyKinds() {
		redisKey, err := key.GenerateRedisKey(kind, data)
		if err != nil {
			blog.Errorf("generate %s %s redis key from refresh data: %+v failed, err: %v, rid: %s", key.Type(), kind,
				data, err, rid)
			return
		}

		pipeline.Expire(redisKey, ttl)
		pipeline.SAdd(redisKey, idKey)
	}

	// update common resource details
	detail, err := json.Marshal(data)
	if err != nil {
		blog.Errorf("marshal %s data %+v failed, err: %v, rid: %s", key.Type(), detail, err, rid)
		return
	}
	pipeline.Set(key.DetailKey(idKey), string(detail), ttl)

	// add common resource id to id list.
	pipeline.ZAddNX(key.IDListKey(), &rawRedis.Z{
		// set common resource id as it's score number
		Score:  score,
		Member: idKey,
	})

	_, err = pipeline.Exec()
	if err != nil {
		blog.Errorf("refresh %s %s redis cache failed, err: %v, data: %+v, rid: %s", key.Type(), idKey, err, data, rid)
		return
	}

	blog.V(4).Infof("refresh %s cache success, id: %s, ttl: %ds, rid: %s", key.Type(), idKey, ttl/time.Second, rid)
}
