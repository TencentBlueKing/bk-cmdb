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

package cache

import (
	"context"
	"errors"

	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/custom/types"
	"configcenter/src/storage/driver/redis"
)

// StrCache caches redis key to data
type StrCache struct {
	key Key
}

// NewStrCache new string type cache
func NewStrCache(key Key) *StrCache {
	return &StrCache{
		key: key,
	}
}

// List get data list by cache keys
func (c *StrCache) List(ctx context.Context, keys []string, rid string) (map[string]string, error) {
	if len(keys) == 0 {
		return make(map[string]string), nil
	}
	keys = util.StrArrayUnique(keys)

	redisKeys := make([]string, len(keys))
	for i, key := range keys {
		redisKeys[i] = c.key.Key(key)
	}

	result, err := redis.Client().MGet(ctx, redisKeys...).Result()
	if err != nil {
		blog.Errorf("list %s data from redis failed, err: %v, keys: %+v, rid: %s", c.key.Type(), err, keys, rid)
		return nil, err
	}

	if len(result) != len(keys) {
		blog.Errorf("%s redis result(%+v) length is invalid, keys: %+v, rid: %s", c.key.Type(), result, keys, rid)
		return nil, errors.New("redis result length is invalid")
	}

	dataMap := make(map[string]string)
	for idx, res := range result {
		if res == nil {
			continue
		}
		detail, ok := res.(string)
		if !ok {
			blog.Errorf("%s redis result type %T is invalid, result: %+v, rid: %s", keys[idx], res, res, rid)
			return nil, errors.New("redis result type is invalid")
		}
		dataMap[keys[idx]] = detail
	}

	return dataMap, nil
}

// BatchUpdate batch update cache by map[key]data
func (c *StrCache) BatchUpdate(ctx context.Context, dataMap map[string]interface{}, rid string) error {
	if len(dataMap) == 0 {
		return nil
	}

	pip := redis.Client().Pipeline()
	defer pip.Close()

	for key, data := range dataMap {
		pip.Set(c.key.Key(key), data, c.key.ttl)
	}

	_, err := pip.Exec()
	if err != nil {
		blog.Errorf("update %s cache failed, err: %v, dataMap: %+v, rid: %s", c.key.Type(), err, dataMap, rid)
		return err
	}

	return nil
}

// BatchDelete batch delete cache keys
func (c *StrCache) BatchDelete(ctx context.Context, keys []string, rid string) error {
	if len(keys) == 0 {
		return nil
	}
	keys = util.StrArrayUnique(keys)

	for i, key := range keys {
		keys[i] = c.key.Key(key)
	}

	err := redis.Client().Del(ctx, keys...).Err()
	if err != nil {
		blog.Errorf("delete %s cache failed, err: %v, keys: %+v, rid: %s", c.key.Type(), err, keys, rid)
		return err
	}

	return nil
}

// Refresh replace the cache info to map[data]count, returns the deleted data list
func (c *StrCache) Refresh(ctx context.Context, match string, dataMap map[string]interface{}, rid string) error {
	pip := redis.Client().Pipeline()
	defer pip.Close()

	keyDataMap := make(map[string]interface{})
	for key, data := range dataMap {
		redisKey := c.key.Key(key)
		keyDataMap[redisKey] = data
		pip.Set(redisKey, data, c.key.ttl)
	}

	match = c.key.Key(match)
	cursor := uint64(0)

	for {
		list, nextCursor, err := redis.Client().Scan(ctx, cursor, match, types.RedisPage).Result()
		if err != nil {
			blog.Errorf("scan %s cache matching %s by cursor %d failed, err: %v, rid: %s", c.key.Type(), match, cursor,
				err, rid)
			return err
		}

		for _, key := range list {
			_, exists := keyDataMap[key]
			if !exists {
				pip.Del(key)
			}
		}

		if nextCursor == uint64(0) {
			break
		}
		cursor = nextCursor
	}

	_, err := pip.Exec()
	if err != nil {
		blog.Errorf("refresh %s cache matching %s failed, err: %v, dataMap: %+v, rid: %s", c.key.Type(), match, err,
			dataMap, rid)
		return err
	}

	return nil
}
