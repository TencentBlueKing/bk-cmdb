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
	"sync"
	"time"

	fullsynccond "configcenter/pkg/cache/full-sync-cond"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/general/types"
	"configcenter/src/storage/driver/redis"
)

// fullSyncCondInfoMap stores the mapping of id list key to the parsed full sync cond info
type fullSyncCondInfoMap struct {
	// info is the mapping of id list key to the parsed full sync cond info
	info map[string]*types.FullSyncCondInfo
	// lock is the rw lock for fullSyncCondInfo
	lock sync.RWMutex
}

func (f *fullSyncCondInfoMap) Get(idListKey string) (*types.FullSyncCondInfo, bool) {
	f.lock.RLock()
	cond, exist := f.info[idListKey]
	f.lock.RUnlock()
	return cond, exist
}

func (f *fullSyncCondInfoMap) Set(idListKey string, cond *types.FullSyncCondInfo) {
	f.lock.Lock()
	f.info[idListKey] = cond
	f.lock.Unlock()
}

func (f *fullSyncCondInfoMap) Remove(idListKey string) {
	f.lock.Lock()
	delete(f.info, idListKey)
	f.lock.Unlock()
}

func (f *fullSyncCondInfoMap) Clear() {
	f.lock.Lock()
	f.info = make(map[string]*types.FullSyncCondInfo)
	f.lock.Unlock()
}

func (f *fullSyncCondInfoMap) Len() int {
	f.lock.RLock()
	length := len(f.info)
	f.lock.RUnlock()
	return length
}

func (f *fullSyncCondInfoMap) Range(handler func(idListKey string, cond *types.FullSyncCondInfo) bool) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	for idListKey, condInfo := range f.info {
		needBreak := handler(idListKey, condInfo)
		if needBreak {
			return
		}
	}
}

// handleFullSyncCondEvent handle full sync cond event, update fullSyncCondInfo and notify the resource watch
func (c *Cache) handleFullSyncCondEvent() {
	for {
		select {
		case e := <-c.fullSyncCondCh:
			rid := util.GenerateRID()
			blog.V(4).Infof("received %s full sync cond event: %+v, rid: %s", c.key.Resource(), e, rid)

			for eventType, conds := range e.EventMap {
				switch eventType {
				case types.Init:
					c.fullSyncCondMap.Clear()
					fallthrough
				case types.Upsert:
					for _, cond := range conds {
						idListKey, err := c.GenFullSyncCondIDListKey(cond)
						if err != nil {
							blog.Errorf("gen full sync cond(%+v) id list key failed, err: %v, rid: %s", cond, err, rid)
							continue
						}

						ttl := time.Duration(cond.Interval) * time.Hour

						// add the non-exist cond info and update the exist id list ttl if changed
						condInfo, exists := c.fullSyncCondMap.Get(idListKey)
						if !exists {
							if !cond.IsAll && cond.Condition == nil {
								blog.Errorf("full sync cond %d is invalid, rid: %s", cond.ID, rid)
								continue
							}
							c.fullSyncCondMap.Set(idListKey, &types.FullSyncCondInfo{
								SubResource:     cond.SubResource,
								IsAll:           cond.IsAll,
								Interval:        ttl,
								Condition:       cond.Condition,
								SupplierAccount: cond.SupplierAccount,
							})
							continue
						}

						if condInfo.Interval != ttl {
							// update full sync cond info ttl
							condInfo.Interval = ttl
							c.fullSyncCondMap.Set(idListKey, condInfo)

							for retry := 0; retry < 3; retry++ {
								if err = c.updateFullSyncCondTTL(idListKey, ttl, rid); err == nil {
									break
								}
								time.Sleep(100 * time.Millisecond * time.Duration(retry))
							}
						}
					}
				case types.Delete:
					for _, cond := range conds {
						idListKey, err := c.GenFullSyncCondIDListKey(cond)
						if err != nil {
							blog.Errorf("gen full sync cond(%+v) id list key failed, err: %v, rid: %s", cond, err, rid)
							continue
						}

						// remove id list key
						c.fullSyncCondMap.Remove(idListKey)

						for retry := 0; retry < 3; retry++ {
							if err = c.deleteFullSyncCondIDList(idListKey, rid); err == nil {
								break
							}
							time.Sleep(100 * time.Millisecond * time.Duration(retry))
						}
					}
				}
			}
			c.cacheChangeCh <- struct{}{}
		}
	}
}

func (c *Cache) updateFullSyncCondTTL(idListKey string, ttl time.Duration, rid string) error {
	// update id list ttl
	err := redis.Client().Expire(context.Background(), idListKey, c.withRandomExpireSeconds(ttl*2)).Err()
	if err != nil {
		blog.Errorf("update id list key: %s ttl to %s failed, err: %v, rid: %s", idListKey, ttl, err, rid)
		return err
	}

	// update id list expire key ttl
	expireKey := c.key.IDListExpireKey(idListKey)
	err = redis.Client().Expire(context.Background(), expireKey, c.withRandomExpireSeconds(ttl)).Err()
	if err != nil {
		blog.Errorf("update id list expire key: %s ttl to %s failed, err: %v, rid: %s", expireKey, ttl, err, rid)
		return err
	}

	return nil
}

func (c *Cache) deleteFullSyncCondIDList(idListKey string, rid string) error {
	// remove id list expire key, the id list will be treated as expired
	expireKey := c.key.IDListExpireKey(idListKey)
	if err := redis.Client().Del(context.Background(), expireKey).Err(); err != nil {
		blog.Errorf("delete expire key: %s failed, err: %v, rid: %s", expireKey, err, rid)
		return err
	}

	exists, err := isIDListExists(context.Background(), idListKey, rid)
	if err != nil {
		blog.Errorf("check if id list key %s exists failed, err: %v, rid: %s", idListKey, err, rid)
		return err
	}

	if !exists {
		return nil
	}

	// rename the id list to avoid reusing the out-dated id list if same id list is watched again
	oldIDListKey := fmt.Sprintf("%s-old", c.key.IDListTempKey(idListKey, rid))
	err = redis.Client().Rename(context.Background(), idListKey, oldIDListKey).Err()
	if err != nil {
		return err
	}

	// delete old id list in background
	go c.deleteIDList(context.Background(), oldIDListKey, rid)
	return nil
}

// GenFullSyncCondIDListKey generate id list key by full sync cond, returns the sync all flag and the id list key
func (c *Cache) GenFullSyncCondIDListKey(cond *fullsynccond.FullSyncCond) (string, error) {
	// generate sync all id list key
	if cond.IsAll && cond.SubResource == "" {
		if c.key.HasSubRes() {
			return "", fmt.Errorf("do not allow sync all cond %d for %s", cond.ID, c.key.Resource())
		}
		return c.key.IDListKey(cond.SupplierAccount), nil
	}

	// generate id list key by sub resource and full sync cond id
	keys := []string{cond.SupplierAccount}

	if cond.SubResource != "" {
		keys = append(keys, cond.SubResource)
	}

	if !cond.IsAll {
		keys = append(keys, strconv.FormatInt(cond.ID, 10))
	}

	return c.key.IDListKey(keys...), nil
}

// genFullSyncCondRefreshIDListOpt generate refresh id list option by full sync cond
func genFullSyncCondRefreshIDListOpt(idListKey string, condInfo *types.FullSyncCondInfo) *refreshIDListOpt {
	return &refreshIDListOpt{
		ttl: condInfo.Interval,
		filterOpt: &types.IDListFilterOpt{
			IDListKey: idListKey,
			BasicFilter: &types.BasicFilter{
				SubRes:          condInfo.SubResource,
				SupplierAccount: condInfo.SupplierAccount,
				IsSystem:        false,
			},
			IsAll: condInfo.IsAll,
			Cond:  condInfo.Condition,
		},
	}
}
