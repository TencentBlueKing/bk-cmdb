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
	"fmt"
	"strconv"
	"time"

	fullsynccond "configcenter/pkg/cache/full-sync-cond"
	"configcenter/pkg/cache/general"
	"configcenter/pkg/filter"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/general/types"
	ccredis "configcenter/src/storage/dal/redis"
	"configcenter/src/storage/driver/redis"

	rawRedis "github.com/go-redis/redis/v7"
)

// AddData add data to general resource cache
func (c *Cache) AddData(ctx context.Context, dataArr []types.WatchEventData, rid string) error {
	if !c.NeedCache() {
		return nil
	}

	if len(dataArr) == 0 {
		return nil
	}

	pip := redis.Client().Pipeline()
	defer pip.Close()

	idKeyMap := make(map[string]*addIDToListOpt)

	for _, data := range dataArr {
		pip, idKeyMap = c.parseAddData(data, pip, idKeyMap, rid)
	}

	var err error
	for _, addOpt := range idKeyMap {
		addOpt.pip = pip
		pip, err = c.addIDToListWithRefresh(ctx, addOpt, rid)
		if err != nil {
			return err
		}
	}

	if _, err = pip.Exec(); err != nil {
		blog.Errorf("add data to %s cache failed, err: %v, data: %+v, rid: %s", c.key.Resource(), err, dataArr, rid)
		return err
	}

	return nil
}

// parseAddData parse added event watch data
func (c *Cache) parseAddData(data types.WatchEventData, pip ccredis.Pipeliner, idKeyMap map[string]*addIDToListOpt,
	rid string) (ccredis.Pipeliner, map[string]*addIDToListOpt) {

	info, err := c.parseData(data)
	if err != nil {
		blog.Errorf("parse %s data: %+v failed, err: %v, rid: %s", c.key.Resource(), data, err, rid)
		return pip, idKeyMap
	}

	id, score := c.key.IDKey(info.id, info.oid)

	// add id to system id lists
	if c.needCacheAll {
		if len(info.subRes) == 0 {
			idKeyMap = c.recordIDToAddForSystem(idKeyMap, id, score, "", info.supplier)
		} else {
			for _, subRes := range info.subRes {
				idKeyMap = c.recordIDToAddForSystem(idKeyMap, id, score, subRes, info.supplier)
			}
		}
	}

	// generate full sync cond id list to id & score map
	c.fullSyncCondMap.Range(func(idListKey string, cond *types.FullSyncCondInfo) bool {
		// remove id from the id list if not matches the full sync cond
		if !c.isFullSyncCondMatched(data.Data, info, cond, rid) {
			pip.ZRem(idListKey, id)
			return false
		}

		_, exists := idKeyMap[idListKey]
		if !exists {
			idKeyMap[idListKey] = &addIDToListOpt{
				idMap:            make(map[string]float64),
				refreshIDListOpt: genFullSyncCondRefreshIDListOpt(idListKey, cond),
			}
		}
		idKeyMap[idListKey].idMap[id] = score
		return false
	})

	// set related id to unique key cache
	for typ, lgc := range c.uniqueKeyLogics {
		redisKeys, err := lgc.genKey(data, info)
		if err != nil {
			blog.Errorf("generate %s %s redis key from data: %+v failed, err: %v, rid: %s", c.key.Resource(), typ,
				data, err, rid)
			continue
		}

		for _, redisKey := range redisKeys {
			pip.Set(c.key.UniqueKey(string(typ), redisKey), id, c.withRandomExpireSeconds(c.expireSeconds))
		}
	}

	return pip, idKeyMap
}

func (c *Cache) recordIDToAddForSystem(idKeyMap map[string]*addIDToListOpt, id string, score float64,
	subRes, supplier string) map[string]*addIDToListOpt {

	idListKey := c.Key().IDListKey()
	if subRes != "" {
		idListKey = c.Key().IDListKey(subRes)
	}

	_, exists := idKeyMap[idListKey]
	if !exists {
		idKeyMap[idListKey] = &addIDToListOpt{
			idMap: make(map[string]float64),
			refreshIDListOpt: &refreshIDListOpt{
				ttl: c.expireSeconds,
				filterOpt: &types.IDListFilterOpt{
					IDListKey: idListKey,
					BasicFilter: &types.BasicFilter{
						SubRes:          subRes,
						SupplierAccount: supplier,
						IsSystem:        true,
					},
					IsAll: true,
				},
			},
		}
	}

	idKeyMap[idListKey].idMap[id] = score
	return idKeyMap
}

// isFullSyncCondMatched check if data matches full sync cond
func (c *Cache) isFullSyncCondMatched(data filter.MatchedData, info *basicInfo, cond *types.FullSyncCondInfo,
	rid string) bool {

	if info.supplier != common.BKDefaultOwnerID && info.supplier != cond.SupplierAccount &&
		cond.SupplierAccount != common.BKSuperOwnerID {
		blog.V(4).Infof("%s data(%+v) supplier not matches cond(%+v), rid: %s", c.key.Resource(), data, cond, rid)
		return false
	}

	subResMatched := true
	if cond.SubResource != "" {
		subResMatched = false
		for _, sub := range info.subRes {
			if cond.SubResource == sub {
				subResMatched = true
				break
			}
		}
	}

	if !subResMatched {
		blog.V(4).Infof("%s data(%+v) sub resource(%+v) not matches cond(%+v), rid: %s", c.key.Resource(), data,
			info.subRes, cond, rid)
		return false
	}

	if cond.IsAll {
		return true
	}

	matched, err := cond.Condition.Match(data)
	if err != nil {
		blog.Errorf("check if %s data: %+v match cond %v failed, err: %v, rid: %s", c.key.Resource(), data,
			cond.Condition, err, rid)
		return false
	}

	if !matched {
		blog.V(4).Infof("%s data(%+v) not matches cond(%+v) condition, rid: %s", c.key.Resource(), data, cond, rid)
		return false
	}

	return true
}

type addIDToListOpt struct {
	pip   ccredis.Pipeliner
	idMap map[string]float64
	*refreshIDListOpt
}

// addIDToListWithRefresh add id to id list cache, refresh the id list if needed
func (c *Cache) addIDToListWithRefresh(ctx context.Context, opt *addIDToListOpt, rid string) (ccredis.Pipeliner,
	error) {

	idListKey := opt.filterOpt.IDListKey

	// try refresh id list if it's not exist or is expired
	notExists, expired, err := c.tryRefreshIDListIfNeeded(ctx, opt.refreshIDListOpt, rid)
	if err != nil {
		blog.Errorf("try refresh id list %s failed, err: %v, opt: %+v, rid: %s", idListKey, err, opt.filterOpt, rid)
		return nil, err
	}

	// id list is refreshing not exist or expired, add id to temp id list key
	if notExists || expired {
		tempKey, err := redis.Client().Get(ctx, c.key.IDListTempKey(idListKey)).Result()
		if err != nil {
			if !redis.IsNilErr(err) {
				blog.Errorf("get id list %s temp key failed, err: %v, rid: %s", idListKey, err, rid)
				return nil, err
			}
			tempKey = c.key.IDListTempKey(idListKey, rid)
		}

		for id, score := range opt.idMap {
			opt.pip.ZAdd(tempKey, &rawRedis.Z{
				Score:  score,
				Member: id,
			})
		}

		opt.pip.Expire(tempKey, c.withRandomExpireSeconds(opt.ttl))
		return opt.pip, nil
	}

	// add id to id list key
	opt.pip = c.addIDToList(opt)
	return opt.pip, nil
}

// addIDToList add id to id list cache
func (c *Cache) addIDToList(opt *addIDToListOpt) ccredis.Pipeliner {
	// add id to id list key
	for id, score := range opt.idMap {
		opt.pip.ZAdd(opt.filterOpt.IDListKey, &rawRedis.Z{
			Score:  score,
			Member: id,
		})
	}
	opt.pip.Expire(opt.filterOpt.IDListKey, c.withRandomExpireSeconds(opt.ttl*2))
	opt.pip.Set(c.key.IDListExpireKey(opt.filterOpt.IDListKey), time.Now().Unix(), c.withRandomExpireSeconds(opt.ttl))
	return opt.pip
}

// RemoveData remove data from general resource cache
func (c *Cache) RemoveData(ctx context.Context, dataArr []types.WatchEventData, rid string) error {
	if !c.NeedCache() {
		return nil
	}

	if len(dataArr) == 0 {
		return nil
	}

	pip := redis.Client().Pipeline()
	defer pip.Close()

	idKeyMap := make(map[string]*removeIDFromListOpt)

	for _, data := range dataArr {
		pip, idKeyMap = c.parseRemoveData(data, pip, idKeyMap, rid)
	}

	var err error
	for _, delOpt := range idKeyMap {
		delOpt.pip = pip
		pip, err = c.removeIDFromListWithRefresh(ctx, delOpt, rid)
		if err != nil {
			return err
		}
	}

	if _, err = pip.Exec(); err != nil {
		blog.Errorf("del data from %s cache failed, err: %v, data: %+v, rid: %s", c.key.Resource(), err, dataArr, rid)
		return err
	}

	return nil
}

func (c *Cache) parseRemoveData(data types.WatchEventData, pip ccredis.Pipeliner,
	idKeyMap map[string]*removeIDFromListOpt, rid string) (ccredis.Pipeliner, map[string]*removeIDFromListOpt) {

	info, err := c.parseData(data)
	if err != nil {
		blog.Errorf("parse %s data: %+v failed, err: %v, rid: %s", c.key.Resource(), data, err, rid)
		return pip, idKeyMap
	}

	id, _ := c.key.IDKey(info.id, info.oid)

	// remove id from system id lists
	if c.needCacheAll {
		if len(info.subRes) == 0 {
			idKeyMap = c.recordIDToRemoveForSystem(idKeyMap, id, "", info.supplier)
		} else {
			for _, subRes := range info.subRes {
				idKeyMap = c.recordIDToRemoveForSystem(idKeyMap, id, subRes, info.supplier)
			}
		}
	}

	// generate full sync cond id list to id & score map
	c.fullSyncCondMap.Range(func(idListKey string, cond *types.FullSyncCondInfo) bool {
		if !c.isFullSyncCondMatched(data.Data, info, cond, rid) {
			return false
		}

		_, exists := idKeyMap[idListKey]
		if !exists {
			idKeyMap[idListKey] = &removeIDFromListOpt{
				ids:              make([]string, 0),
				refreshIDListOpt: genFullSyncCondRefreshIDListOpt(idListKey, cond),
			}
		}
		idKeyMap[idListKey].ids = append(idKeyMap[idListKey].ids, id)
		return false
	})

	// remove related id from unique key cache
	for typ, lgc := range c.uniqueKeyLogics {
		redisKeys, err := lgc.genKey(data, info)
		if err != nil {
			blog.Errorf("generate %s %s redis key from data: %+v failed, err: %v, rid: %s", c.key.Resource(), typ,
				data, err, rid)
			continue
		}

		for _, redisKey := range redisKeys {
			pip.Del(c.key.UniqueKey(string(typ), redisKey))
		}
	}
	return pip, idKeyMap
}

func (c *Cache) recordIDToRemoveForSystem(idKeyMap map[string]*removeIDFromListOpt, id string,
	subRes, supplier string) map[string]*removeIDFromListOpt {

	idListKey := c.Key().IDListKey()
	if subRes != "" {
		idListKey = c.Key().IDListKey(subRes)
	}

	_, exists := idKeyMap[idListKey]
	if !exists {
		idKeyMap[idListKey] = &removeIDFromListOpt{
			ids: make([]string, 0),
			refreshIDListOpt: &refreshIDListOpt{
				ttl: c.expireSeconds,
				filterOpt: &types.IDListFilterOpt{
					IDListKey: idListKey,
					BasicFilter: &types.BasicFilter{
						SubRes:          subRes,
						SupplierAccount: supplier,
						IsSystem:        true,
					},
					IsAll: true,
				},
			},
		}
	}

	idKeyMap[idListKey].ids = append(idKeyMap[idListKey].ids, id)
	return idKeyMap
}

type removeIDFromListOpt struct {
	pip ccredis.Pipeliner
	ids []string
	*refreshIDListOpt
}

// removeIDFromList remove id from id list cache, refresh the id list if needed
func (c *Cache) removeIDFromListWithRefresh(ctx context.Context, opt *removeIDFromListOpt, rid string) (
	ccredis.Pipeliner, error) {

	idListKey := opt.filterOpt.IDListKey

	// try refresh id list if it's not exist or is expired
	notExists, expired, err := c.tryRefreshIDListIfNeeded(ctx, opt.refreshIDListOpt, rid)
	if err != nil {
		blog.Errorf("try refresh id list %s failed, err: %v, opt: %+v, rid: %s", idListKey, err, opt.filterOpt, rid)
		return nil, err
	}

	// id list is refreshing not exist or expired, remove id from temp id list key
	if notExists || expired {
		tempKey, err := redis.Client().Get(ctx, c.key.IDListTempKey(idListKey)).Result()
		if err != nil {
			blog.Errorf("get  id list %s temp key failed, err: %v, rid: %s", idListKey, err, rid)
			return nil, err
		}

		for _, id := range opt.ids {
			opt.pip.ZRem(tempKey, id)
		}
		opt.pip.Expire(tempKey, c.withRandomExpireSeconds(opt.ttl))
		return opt.pip, nil
	}

	// remove id from id list key
	opt.pip = c.removeIDFromList(opt)
	return opt.pip, nil
}

// removeIDFromList remove id from id list cache
func (c *Cache) removeIDFromList(opt *removeIDFromListOpt) ccredis.Pipeliner {
	for _, id := range opt.ids {
		opt.pip.ZRem(opt.filterOpt.IDListKey, id)
	}
	opt.pip.Expire(opt.filterOpt.IDListKey, c.withRandomExpireSeconds(opt.ttl*2))
	opt.pip.Set(c.key.IDListExpireKey(opt.filterOpt.IDListKey), time.Now().Unix(), c.withRandomExpireSeconds(opt.ttl))
	return opt.pip
}

type refreshIDListOpt struct {
	filterOpt *types.IDListFilterOpt
	ttl       time.Duration
}

// tryRefreshIDListIfNeeded try refresh id list cache if it's not exist or expired
// return params: notExists: returns if the id list is not exist, expired: returns if the id list is expired
func (c *Cache) tryRefreshIDListIfNeeded(ctx context.Context, opt *refreshIDListOpt, rid string) (notExists bool,
	expired bool, err error) {

	idListKey := opt.filterOpt.IDListKey
	exists, err := isIDListExists(ctx, idListKey, rid)
	if err != nil {
		return false, false, err
	}

	if !exists {
		c.tryRefreshIDList(ctx, opt, rid)
		return true, false, nil
	}

	expire, err := redis.Client().Get(ctx, c.key.IDListExpireKey(idListKey)).Result()
	if err != nil {
		if redis.IsNilErr(err) {
			blog.V(4).Infof("id list %s expire key not exists, refresh it now. rid: %s", idListKey, rid)
			c.tryRefreshIDList(ctx, opt, rid)
			return false, true, nil
		}

		blog.Errorf("get host id list expire key failed, err: %v, rid :%v", err, rid)
		return false, false, err
	}

	expireAt, err := strconv.ParseInt(expire, 10, 64)
	if err != nil {
		blog.Errorf("parse id list %s expire time %s failed, err: %v, rid: %s", idListKey, expire, err, rid)
		return false, false, err
	}

	expireSeconds := int64(opt.ttl.Seconds())
	if time.Now().Unix()-expireAt <= expireSeconds {
		// not expired
		return false, false, nil
	}

	// set expire key with a value which will enforce the id list key to expire within one minute
	// which will block the refresh request for the next minute. This policy is used to avoid refreshing keys
	// when redis is under high pressure or not well performed.
	redis.Client().Set(ctx, c.key.IDListExpireKey(idListKey), time.Now().Unix()-expireSeconds+60, time.Minute)

	// expired, we refresh it now.
	blog.V(4).Infof("id list %s is expired, refresh it now. rid: %s", idListKey, rid)
	c.tryRefreshIDList(ctx, opt, rid)
	return false, true, nil
}

// tryRefreshIDList try refresh the general resource id list cache if it's not locked
func (c *Cache) tryRefreshIDList(ctx context.Context, opt *refreshIDListOpt, rid string) {
	idListKey := opt.filterOpt.IDListKey
	if idListKey == "" {
		blog.Errorf("id list key is not set, opt: %+v, rid: %s", opt, rid)
		return
	}

	lockKey := c.key.IDListLockKey(idListKey)

	// get local lock
	if !c.refreshingLock.CanRefresh(lockKey) {
		blog.V(4).Infof("%s id list lock %s is locked, skip refresh, rid: %s", c.key.Resource(), lockKey, rid)
		return
	}

	// set refreshing status
	c.refreshingLock.SetRefreshing(lockKey)

	// then get distribute lock
	locked, err := redis.Client().SetNX(ctx, lockKey, rid, 5*time.Minute).Result()
	if err != nil {
		blog.Errorf("get id list %s lock failed, err: %v, rid: %s", idListKey, err, rid)
		c.refreshingLock.SetUnRefreshing(lockKey)
		return
	}

	if !locked {
		blog.V(4).Infof("%s id list key redis lock %s is locked, skip refresh, rid: %s", c.key.Resource(), lockKey, rid)
		c.refreshingLock.SetUnRefreshing(lockKey)
		return
	}

	go func() {
		ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)
		blog.V(4).Infof("start refresh %s id list cache %s, rid: %s", c.key.Resource(), idListKey, rid)

		defer c.refreshingLock.SetUnRefreshing(lockKey)
		defer redis.Client().Del(ctx, lockKey)

		// already get lock, refresh the id list cache now
		err = c.refreshIDList(ctx, opt, rid)
		if err != nil {
			blog.Errorf("refresh %s id list cache %s failed, err: %v, rid: %s", c.key.Resource(), idListKey, err, rid)
			return
		}

		blog.V(4).Infof("refresh %s id list cache %s success, rid: %s", c.key.Resource(), idListKey, rid)
	}()
}

// refreshIDList refresh the general resource id list cache
func (c *Cache) refreshIDList(ctx context.Context, opt *refreshIDListOpt, rid string) error {
	ctx = util.SetDBReadPreference(ctx, common.SecondaryPreferredMode)

	idListKey := opt.filterOpt.IDListKey
	tempKey := c.key.IDListTempKey(idListKey, rid)

	// set the temp id list key in redis for the event watch to judge which temp id list to write to
	err := redis.Client().Set(ctx, c.key.IDListTempKey(idListKey), tempKey, c.withRandomExpireSeconds(opt.ttl)).Err()
	if err != nil {
		blog.Errorf("set temp id list key %s failed, err: %v, rid: %s", tempKey, err, rid)
		return err
	}
	defer func() {
		if err := redis.Client().Del(context.Background(), c.key.IDListTempKey(idListKey)).Err(); err != nil {
			blog.Errorf("delete temp id list key %s failed, err: %v, rid: %s", tempKey, err, rid)
		}
	}()

	blog.V(4).Infof("try to refresh id list %s with temp key: %s, rid: %s", idListKey, tempKey, rid)

	listOpt := &types.ListDetailOpt{
		OnlyListID:   true,
		IDListFilter: opt.filterOpt,
		Page:         &general.PagingOption{Limit: types.PageSize},
	}
	total := 0
	for {
		dbRes, err := c.listDataFromDB(ctx, listOpt, rid)
		if err != nil {
			return err
		}
		dbData := dbRes.Data

		stepLen := len(dbData)
		if stepLen == 0 {
			break
		}
		total += stepLen

		pip := redis.Client().Pipeline()
		// because the temp key is a random key, so we set an expiry time so that it can be gc,
		// but we will reset the expiry time when this key is renamed to a normal key.
		pip.Expire(tempKey, c.withRandomExpireSeconds(opt.ttl))
		for _, data := range dbData {
			id, score, err := c.generateID(data)
			if err != nil {
				blog.Errorf("generate %s id from data: %+v failed, err: %v, rid: %s", c.key.Resource(), data, err, rid)
				continue
			}

			// write to the temp key
			pip.ZAdd(tempKey, &rawRedis.Z{
				Score:  score,
				Member: id,
			})
		}

		if _, err = pip.Exec(); err != nil {
			blog.Errorf("update temp id list %s failed, err: %v, data: %+v, rid: %s", tempKey, err, dbData, rid)
			return err
		}

		if stepLen < types.PageSize {
			break
		}

		info, err := c.parseData(dbData[stepLen-1])
		if err != nil {
			blog.Errorf("parse %s data(%+v) failed, err: %v, rid: %s", c.key.Resource(), dbData[stepLen-1], err, rid)
			return err
		}

		listOpt.Page.StartID = info.id
		listOpt.Page.StartOid = info.oid
		time.Sleep(50 * time.Millisecond)
	}

	if total == 0 {
		return nil
	}

	// if id list exists, we need to delete it
	exists, err := isIDListExists(ctx, idListKey, rid)
	if err != nil {
		return err
	}

	pipe := redis.Client().Pipeline()
	tempOldKey := fmt.Sprintf("%s-old", tempKey)
	if exists {
		// rename id list key to temp old id list key so that we can delete later, avoiding the implicit del in rename
		// which will block all the following redis operation
		pipe.Rename(idListKey, tempOldKey)
	}
	// rename temp key to real key
	pipe.RenameNX(tempKey, idListKey)
	// reset id_list key's expire time to a new one.
	pipe.Expire(idListKey, c.withRandomExpireSeconds(opt.ttl*2))
	// set expire key with unix time seconds now value.
	pipe.Set(c.key.IDListExpireKey(idListKey), time.Now().Unix(), c.withRandomExpireSeconds(opt.ttl))

	if _, err = pipe.Exec(); err != nil {
		blog.Errorf("refresh id list %s with temp key: %s failed, err :%v, rid: %s", idListKey, tempKey, err, rid)
		return err
	}

	if exists {
		// remove the old id list key in background
		go c.deleteIDList(context.Background(), tempOldKey, rid)
	}

	blog.V(4).Infof("refresh id list key: %s success, count: %d. rid: %s", idListKey, total, rid)
	return nil
}

// deleteIDList delete the general resource id list cache
func (c *Cache) deleteIDList(ctx context.Context, key string, rid string) error {
	for {
		cnt, err := redis.Client().ZRemRangeByRank(key, 0, types.PageSize).Result()
		if err != nil {
			blog.Errorf("delete id list: %s failed, err: %v, rid: %s", key, err, rid)
			return err
		}

		if cnt < types.PageSize {
			return nil
		}

		time.Sleep(200 * time.Millisecond)
	}
}

// listIDsFromRedis list general resource id list from redis
func (c *Cache) listIDsFromRedis(ctx context.Context, key string, opt *general.PagingOption, rid string) ([]string,
	error) {

	if opt.Limit == 0 {
		return make([]string, 0), nil
	}

	// list from start id
	if opt.StartID > 0 {
		redisOpt := &rawRedis.ZRangeBy{
			Min:   fmt.Sprintf("(%d", opt.StartID),
			Max:   "+inf",
			Count: opt.Limit,
		}
		ids, err := redis.Client().ZRangeByScore(ctx, key, redisOpt).Result()
		if err != nil {
			blog.Errorf("list %s ids from cache failed, err: %v, redis opt: %+v, rid: %s", key, err, redisOpt, rid)
			return nil, err
		}
		return ids, nil
	}

	// list from start oid
	if len(opt.StartOid) > 0 {
		redisOpt := &rawRedis.ZRangeBy{
			Min:   "(" + opt.StartOid,
			Max:   "+",
			Count: opt.Limit,
		}
		ids, err := redis.Client().ZRangeByLex(ctx, key, redisOpt).Result()
		if err != nil {
			blog.Errorf("list %s ids from cache failed, err: %v, redis opt: %+v, rid: %s", key, err, redisOpt, rid)
			return nil, err
		}
		return ids, nil
	}

	// list from start index
	ids, err := redis.Client().ZRange(ctx, key, opt.StartIndex, opt.StartIndex+opt.Limit-1).Result()
	if err != nil {
		blog.Errorf("list %s ids from cache failed, err: %v, opt: %+v, rid: %s", key, err, opt, rid)
		return nil, err
	}
	return ids, nil
}

// countIDsFromRedis count general resource id list from redis
func (c *Cache) countIDsFromRedis(ctx context.Context, key string, rid string) (int64, error) {
	cnt, err := redis.Client().ZCard(ctx, key).Result()
	if err != nil {
		blog.Errorf("count %s ids from cache failed, err: %v, rid: %s", key, err, rid)
		return 0, err
	}
	return cnt, nil
}

// RefreshIDList refresh general resource id list cache
func (c *Cache) RefreshIDList(kit *rest.Kit, opt *general.RefreshIDListOpt,
	cond *fullsynccond.FullSyncCond) error {
	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		blog.Errorf("refresh id list option is invalid, err: %v, opt: %+v, rid: %s", rawErr, opt, kit.Rid)
		return rawErr.ToCCError(kit.CCError)
	}

	refreshOpt := &refreshIDListOpt{
		filterOpt: &types.IDListFilterOpt{
			BasicFilter: &types.BasicFilter{
				SubRes:          opt.SubRes,
				SupplierAccount: kit.SupplierAccount,
			},
			IsAll: true,
		},
	}

	// refresh full sync cond id list cache
	if opt.CondID > 0 {
		refreshOpt.filterOpt.IsAll = cond.IsAll
		refreshOpt.filterOpt.Cond = cond.Condition

		var err error
		refreshOpt.filterOpt.IDListKey, err = c.GenFullSyncCondIDListKey(cond)
		if err != nil {
			return err
		}

		refreshOpt.ttl = time.Duration(cond.Interval) * time.Hour
		c.tryRefreshIDList(kit.Ctx, refreshOpt, kit.Rid)
		return nil
	}

	// refresh system id list cache
	refreshOpt.filterOpt.IsSystem = true
	idListTTL, err := c.validateIDList(refreshOpt.filterOpt)
	if err != nil {
		blog.Errorf("id list filter option is invalid, err: %v, opt: %+v, rid: %s", err, opt, kit.Rid)
		return err
	}
	refreshOpt.ttl = idListTTL

	if opt.SubRes != "" {
		refreshOpt.filterOpt.IDListKey = c.Key().IDListKey(opt.SubRes)
	} else {
		refreshOpt.filterOpt.IDListKey = c.Key().IDListKey()
	}

	c.tryRefreshIDList(kit.Ctx, refreshOpt, kit.Rid)
	return nil
}
