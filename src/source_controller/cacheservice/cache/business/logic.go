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

package business

import (
	"context"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/storage/driver/redis"
)

func upsertListCache(ms *forUpsertCache) {
	listKeyName := ms.listKey
	blog.V(3).Infof("received key %s %d/%s change event, and try refresh cache, data: %s", listKeyName, ms.instID, ms.name, string(ms.doc))

	listKeys, err := redis.Client().SMembers(context.Background(), listKeyName).Result()
	if err != nil {
		blog.Errorf("upsert cache, but get all cached keys %s failed. err: %v", listKeyName, err)
		return
	}

	hitKey := ""
	for _, key := range listKeys {
		id, parentID, _, err := ms.parseListKeyValue(key)
		if err != nil {
			blog.Errorf("got invalid cache %s key: %s, err: %v", listKeyName, key, err)
			// invalid key, delete immediately
			// normally, this can not be happen.
			// we try our best to correct the data
			if redis.Client().SRem(context.Background(), listKeyName, key).Err() != nil {
				blog.Errorf("delete invalid list key name %s: value: %s failed,", listKeyName, key)
				// we do not return and continue
			}

			if id != 0 {
				// try to find it from db and add again
				name, err := ms.getInstName(id)
				if err != nil {
					blog.Errorf("upsert cache, got invalid cache key %s, try to correct it, but get name failed, err: %v", listKeyName, key, err)
				} else {
					newKey := ms.genListKeyValue(id, parentID, name)
					if err := redis.Client().SAdd(context.Background(), listKeyName, newKey).Err(); err != nil {
						blog.Errorf("add new cache key failed, key: %s, err: %v", newKey, err)
						// we do not return and continue
					}
				}
			}
		}
		if ms.instID == id {
			hitKey = key
		}
	}

	pipeline := redis.Client().Pipeline()
	defer pipeline.Close()

	// check the key is change or a new one.
	newKey := ms.genListKeyValue(ms.instID, ms.parentID, ms.name)
	if len(hitKey) == 0 {
		// a new one, add it now
		pipeline.SAdd(listKeyName, newKey)
	} else {
		// we check the key is same or not, if not, means it's name has been changed
		if newKey != hitKey {
			// name changed, update now
			pipeline.SRem(listKeyName, hitKey)
			pipeline.SAdd(listKeyName, newKey)
		}
	}

	// set the new data
	// TODO: get the detailed info from db and set the key with the latest data
	pipeline.Set(ms.detailKey, string(ms.doc), 0)

	// set the expire key
	pipeline.Set(ms.listExpireKey, time.Now().Unix(), 0)
	pipeline.Set(ms.detailExpireKey, time.Now().Unix(), 0)

	_, err = pipeline.Exec()
	if err != nil {
		blog.Errorf("received key %s %d/%s change event, but upsert to redis failed, err: %v", listKeyName, ms.instID, ms.name, err)
		return
	}

	blog.V(4).Infof("received key %s %d/%s change event, and refresh cache success, data: %s", listKeyName, ms.instID, ms.name, string(ms.doc))
}
