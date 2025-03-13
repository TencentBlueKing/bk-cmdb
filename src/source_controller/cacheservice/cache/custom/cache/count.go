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
	"strconv"

	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/source_controller/cacheservice/cache/custom/types"
	"configcenter/src/storage/driver/redis"
)

// CountCache caches the mapping of data to its count
type CountCache struct {
	key Key
}

// NewCountCache new count cache
func NewCountCache(key Key) *CountCache {
	return &CountCache{
		key: key,
	}
}

// GetDataList get data list by cache key
func (c *CountCache) GetDataList(kit *rest.Kit, key string) ([]string, error) {
	cacheKey := c.key.Key(kit.TenantID, key)
	cursor := uint64(0)

	all := make([]string, 0)
	for {
		list, nextCursor, err := redis.Client().HScan(kit.Ctx, cacheKey, cursor, "", types.RedisPage).Result()
		if err != nil {
			blog.Errorf("scan %s data list cache by cursor %d failed, err: %v, rid: %s", cacheKey, cursor, err, kit.Rid)
			return nil, err
		}

		for i := 0; i < len(list)-1; i += 2 {
			all = append(all, list[i])
		}

		if nextCursor == uint64(0) {
			return all, nil
		}
		cursor = nextCursor
	}
}

// UpdateCount batch update cache by map[key]map[data]count
func (c *CountCache) UpdateCount(kit *rest.Kit, cntMap map[string]map[string]int64) error {
	for key, dataCnt := range cntMap {
		var args []interface{}
		for k, v := range dataCnt {
			args = append(args, k, v)
		}
		args = append(args, int64(c.key.TTL().Seconds()))

		err := redis.Client().Eval(kit.Ctx, updateCountScript, []string{c.key.Key(kit.TenantID, key)}, args...).Err()
		if err != nil {
			blog.Errorf("update type: %s key: %s count cache failed, err: %v, data: %+v, rid: %s", c.key.Type(), key,
				err, dataCnt, kit.Rid)
			return err
		}
	}

	return nil
}

const updateCountScript = `
local cacheKey = KEYS[1]

for i = 1, #ARGV-1, 2 do
	local res = redis.call('HINCRBY', cacheKey, ARGV[i], ARGV[i+1])
	if tonumber(res) < 1 then
		redis.call('HDEL', cacheKey, ARGV[i])
	end
end

redis.call('EXPIRE', cacheKey, ARGV[#ARGV]) 

return 1
`

// RefreshCount replace the cache info to map[data]count, returns the deleted data list
func (c *CountCache) RefreshCount(kit *rest.Kit, key string, cntMap map[string]int64) ([]string,
	error) {

	cacheKey := c.key.Key(kit.TenantID, key)

	pip := redis.Client().Pipeline()
	defer pip.Close()

	cntKeyValues := make([]string, 0)
	for data, count := range cntMap {
		cntKeyValues = append(cntKeyValues, data, strconv.FormatInt(count, 10))
		if len(cntKeyValues) >= 2*types.RedisPage {
			pip.HMSet(cacheKey, cntKeyValues)
			cntKeyValues = make([]string, 0)
		}
	}

	if len(cntKeyValues) > 0 {
		pip.HMSet(cacheKey, cntKeyValues)
	}

	list, err := c.GetDataList(kit, key)
	if err != nil {
		return nil, err
	}

	delList := make([]string, 0)
	for _, data := range list {
		_, exists := cntMap[data]
		if !exists {
			delList = append(delList, data)
		}

		if len(delList) >= types.RedisPage {
			pip.HDel(cacheKey, delList...)
			delList = make([]string, 0)
		}
	}

	if len(delList) > 0 {
		pip.HDel(cacheKey, delList...)
	}

	pip.Expire(cacheKey, c.key.TTL())

	_, err = pip.Exec()
	if err != nil {
		blog.Errorf("refresh type: %s key: %s count cache failed, err: %v, data: %+v, rid: %s", c.key.Type(), key, err,
			cntMap, kit.Rid)
		return nil, err
	}

	return delList, nil
}

// Delete delete cache key
func (c *CountCache) Delete(kit *rest.Kit, key string) error {
	cacheKey := c.key.Key(kit.TenantID, key)

	pip := redis.Client().Pipeline()
	defer pip.Close()

	list, err := c.GetDataList(kit, key)
	if err != nil {
		return err
	}

	delList := make([]string, 0)
	for _, data := range list {
		delList = append(delList, data)

		if len(delList) >= types.RedisPage {
			pip.HDel(cacheKey, delList...)
			delList = make([]string, 0)
		}
	}

	if len(delList) > 0 {
		pip.HDel(cacheKey, delList...)
	}

	_, err = pip.Exec()
	if err != nil {
		blog.Errorf("delete %s count cache key failed, err: %v, rid: %s", cacheKey, err, kit.Rid)
		return err
	}

	return nil
}
