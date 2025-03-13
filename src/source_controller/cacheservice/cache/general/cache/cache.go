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

// Package cache defines the general resource cache
package cache

import (
	"math/rand"
	"sync"
	"time"

	"configcenter/pkg/cache/general"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/source_controller/cacheservice/cache/general/types"
	"configcenter/src/source_controller/cacheservice/cache/tools"
)

// cacheSet is the general resource cache key -> cache logics map
var cacheSet = make(map[general.ResType]*Cache)

// addCache add general resource cache to cache set
func addCache(cache *Cache) {
	cacheSet[cache.key.Resource()] = cache
	go cache.handleFullSyncCondEvent()
}

// GetAllCache get all general resource cache type -> cache logics map
func GetAllCache() map[general.ResType]*Cache {
	return cacheSet
}

// Cache is the general resource cache key
type Cache struct {
	key *general.Key
	// expireSeconds is the ttl of all the general resource built-in cache keys.
	expireSeconds time.Duration
	// expireRangeSeconds defines the additional min range and max range of all the built-in general res cache keys ttl
	expireRangeSeconds [2]int
	// refreshingLock is used to lock the general resource cache key when refreshing
	refreshingLock tools.RefreshingLock

	// needCacheAll is used to indicate whether to cache all the general resource id list by default
	needCacheAll bool

	// parseData parse general resource data
	parseData dataParser
	// getDataByID get general resource data from db by id keys, used to get detail when detail cache is not available
	getDataByID dataGetterByKeys
	// listData list general resource data from db, used to get data when id list cache is not available
	listData dataLister

	// uniqueKeyLogics is the unique key type to related logics map
	uniqueKeyLogics map[general.UniqueKeyType]uniqueKeyLogics

	// fullSyncCondMap stores the mapping of id list key to the parsed full sync cond info
	fullSyncCondMap *fullSyncCondInfoMap
	// fullSyncCondCh is the channel to receive full sync cond event
	fullSyncCondCh chan types.FullSyncCondEvent
	// cacheChangeCh is the channel to notify cache change event
	cacheChangeCh chan struct{}
}

// NewCache new empty Cache
func NewCache() *Cache {
	return &Cache{
		refreshingLock:  tools.NewRefreshingLock(),
		uniqueKeyLogics: make(map[general.UniqueKeyType]uniqueKeyLogics),
		fullSyncCondMap: &fullSyncCondInfoMap{
			info: make(map[string]*types.FullSyncCondInfo),
			lock: sync.RWMutex{},
		},
		fullSyncCondCh: make(chan types.FullSyncCondEvent, 1),
		cacheChangeCh:  make(chan struct{}, 1),
	}
}

// basicInfo is the general resource data's basic info to generate redis key
type basicInfo struct {
	id     int64
	oid    string
	subRes []string
}

// dataParser parse the general resource data to basic info
type dataParser func(data any) (*basicInfo, error)

// getDataByKeysOpt is get general resource data from db by keys option
type getDataByKeysOpt struct {
	*types.BasicFilter
	Keys []string
}

// dataGetterByKeys get mongodb data by redis keys
type dataGetterByKeys func(kit *rest.Kit, opt *getDataByKeysOpt) ([]any, error)

// listDataOpt is list general resource data from db option
type listDataOpt struct {
	*types.BasicFilter
	Fields []string
	// OnlyListID only list the fields that are used to generated general resource cache id key
	OnlyListID bool
	Cond       mapstr.MapStr
	Page       *general.PagingOption
}

// listDataRes is list general resource data from db result
type listDataRes struct {
	Count uint64
	Data  []any
}

// dataLister list mongodb data
type dataLister func(kit *rest.Kit, opt *listDataOpt) (*listDataRes, error)

type uniqueKeyLogics struct {
	// genKey generator redis key of the unique key type
	// the unique key cache stores the related id key to get detail indirectly
	genKey redisKeyGenerator
	// getData get general resource data from db by unique keys, used to get detail when cache is not available
	getData dataGetterByKeys
}

// redisKeyGenerator generates redis keys from data
type redisKeyGenerator func(data any, info *basicInfo) ([]string, error)

// Key returns the general resource cache key
func (c *Cache) Key() *general.Key {
	return c.key
}

// generateID generates general resource id key and score from data
func (c *Cache) generateID(data any) (string, float64, error) {
	info, err := c.parseData(data)
	if err != nil {
		return "", 0, err
	}

	idKey, score := c.key.IDKey(info.id, info.oid)
	return idKey, score, nil
}

// withRandomExpireSeconds generate random redis key expire in seconds
func (c *Cache) withRandomExpireSeconds(expireSeconds time.Duration) time.Duration {
	seconds := c.expireRangeSeconds[0] + rand.New(rand.NewSource(time.Now().UnixNano())).
		Intn(c.expireRangeSeconds[1]-c.expireRangeSeconds[0])
	return expireSeconds + time.Duration(seconds)*time.Second
}

// FullSyncCondCh returns the full sync cond channel
func (c *Cache) FullSyncCondCh() chan types.FullSyncCondEvent {
	return c.fullSyncCondCh
}

// CacheChangeCh returns the general resource cache change channel
func (c *Cache) CacheChangeCh() chan struct{} {
	return c.cacheChangeCh
}

// NeedWatchRes returns whether all resource data needs to be watched, and the specified sub-resources to be watched
func (c *Cache) NeedWatchRes() (bool, map[string][]string) {
	if c.needCacheAll || len(c.uniqueKeyLogics) != 0 {
		return true, nil
	}

	needWatchAllTenants := make([]string, 0)
	tenantSubResMap := make(map[string][]string)
	c.fullSyncCondMap.Range(func(idListKey string, cond *types.FullSyncCondInfo) bool {
		if cond.SubResource == "" {
			needWatchAllTenants = append(needWatchAllTenants, cond.TenantID)
			return false
		}

		tenantSubResMap[cond.TenantID] = append(tenantSubResMap[cond.TenantID], cond.SubResource)
		return false
	})

	for _, tenantID := range needWatchAllTenants {
		tenantSubResMap[tenantID] = make([]string, 0)
	}

	return false, tenantSubResMap
}

// NeedCache returns if the general resource needs to be cached
func (c *Cache) NeedCache() bool {
	return c.needCacheAll || c.fullSyncCondMap.Len() != 0 || len(c.uniqueKeyLogics) != 0
}
