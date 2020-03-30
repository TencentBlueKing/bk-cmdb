/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cache

import (
	"context"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/reflector"
	"configcenter/src/storage/stream/types"

	"gopkg.in/redis.v5"
)

// sync database data to redis cache
func SyncDatabaseToRedis(ctx context.Context, conf local.MongoConf, cache *redis.Client, listDone chan<- bool, errChan chan<- error) {
	ref, err := reflector.NewReflector(conf)
	if err != nil {
		blog.Errorf("new reflector failed, error: %s", err.Error())
		errChan <- err
		return
	}

	var wg sync.WaitGroup

	// map database with redis key generate function
	collectionKeyMap := map[string]func(cacheMap map[string]string) string{
		common.BKTableNameObjAttDes: func(cacheMap map[string]string) string {
			return common.BKCacheKeyV3Prefix + "attribute:object" + cacheMap[common.BKObjIDField] + "id:" + cacheMap[common.BKPropertyIDField]
		},
		common.BKTableNameBaseApp: func(cacheMap map[string]string) string {
			return common.BKCacheKeyV3Prefix + common.BKInnerObjIDApp + ":id:" + cacheMap[common.BKAppIDField]
		},
		common.BKTableNameBaseSet: func(cacheMap map[string]string) string {
			return common.BKCacheKeyV3Prefix + common.BKInnerObjIDSet + ":biz:" + cacheMap[common.BKAppIDField] + ":id:" + cacheMap[common.BKSetIDField]
		},
		common.BKTableNameBaseModule: func(cacheMap map[string]string) string {
			return common.BKCacheKeyV3Prefix + common.BKInnerObjIDModule + ":biz:" + cacheMap[common.BKAppIDField] + ":set:" + cacheMap[common.BKSetIDField] + ":id:" + cacheMap[common.BKModuleIDField]
		},
		common.BKTableNameModuleHostConfig: func(cacheMap map[string]string) string {
			return common.BKCacheKeyV3Prefix + "host_module:host:" + cacheMap[common.BKHostIDField] + ":module:" + cacheMap[common.BKModuleIDField]
		},
	}

	// only sync once
	if cache.Exists(common.RedisMongoCacheSyncKey).Val() {
		listDone <- true
		return
	}
	cache.Set(common.RedisMongoCacheSyncKey, "", 0)
	defer cache.Del(common.RedisMongoCacheSyncKey)

	for collection, generateKeyFunc := range collectionKeyMap {
		isIgnoreField := getIgnoreFields(collection)
		isIgnoreFieldLen := len(isIgnoreField)

		// cache mongo data in redis
		cacheFunc := func(event *types.Event) {
			if doc, ok := event.Document.(map[string]interface{}); ok {
				cacheMap := make(map[string]string, len(doc)-isIgnoreFieldLen)
				for key, value := range doc {
					if !isIgnoreField[key] {
						cacheMap[key] = util.GetStrByInterface(value)
					}
				}
				key := generateKeyFunc(cacheMap)
				cache.HMSet(key, cacheMap)
			}
		}

		capable := &reflector.Capable{
			OnChange: reflector.OnChangeEvent{
				OnLister: cacheFunc,
				OnListerDone: func() {
					wg.Done()
				},
				OnAdd:    cacheFunc,
				OnUpdate: cacheFunc,
				OnDelete: func(event *types.Event) {
					if doc, ok := event.Document.(map[string]interface{}); ok {
						cacheMap := make(map[string]string, len(doc))
						for key, value := range doc {
							cacheMap[key] = util.GetStrByInterface(value)
						}
						key := generateKeyFunc(cacheMap)
						cache.Del(key)
					}
				},
			},
		}

		opt := &types.ListWatchOptions{
			Options: types.Options{
				EventStruct: new(map[string]interface{}),
				Collection:  collection,
			},
		}

		err = ref.ListWatcher(ctx, opt, capable)
		if err != nil {
			blog.Errorf("list watcher failed, error: %s, database: %s", err.Error(), collection)
			errChan <- err
			return
		}
		wg.Add(1)
	}

	// wait till all databases are cached and begin to watch changing event
	wg.Wait()
	listDone <- true

	<-ctx.Done()
	return
}

func getIgnoreFields(collection string) map[string]bool {
	switch collection {
	default:
		return map[string]bool{
			"_id": true,
		}
	}
}
