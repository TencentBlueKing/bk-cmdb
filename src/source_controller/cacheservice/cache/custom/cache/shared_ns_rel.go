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
	"fmt"
	"strconv"
	"time"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/lock"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/source_controller/cacheservice/cache/custom/types"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
)

// SharedNsRelCache is shared namespace relation cache
type SharedNsRelCache struct {
	isMaster       discovery.ServiceManageInterface
	nsAsstBizCache *StrCache
}

// NewSharedNsRelCache new shared namespace relation cache
func NewSharedNsRelCache(isMaster discovery.ServiceManageInterface) *SharedNsRelCache {
	return &SharedNsRelCache{
		isMaster:       isMaster,
		nsAsstBizCache: NewStrCache(Key{resType: types.SharedNsAsstBizType, ttl: 6 * time.Hour}),
	}
}

// genNsAsstBizRedisKey generate redis key for shared namespace to asst biz cache
func (c *SharedNsRelCache) genNsAsstBizRedisKey(nsID int64) string {
	return strconv.FormatInt(nsID, 10)
}

// GetAsstBiz get shared namespace associated biz info
func (c *SharedNsRelCache) GetAsstBiz(ctx context.Context, nsIDs []int64, rid string) (map[int64]int64, error) {
	redisKeys := make([]string, len(nsIDs))
	for i, nsID := range nsIDs {
		redisKeys[i] = c.genNsAsstBizRedisKey(nsID)
	}

	redisDataMap, err := c.nsAsstBizCache.List(ctx, redisKeys, rid)
	if err != nil {
		return nil, err
	}

	nsAsstBizMap := make(map[int64]int64)
	for redisKey, redisData := range redisDataMap {
		nsID, err := strconv.ParseInt(redisKey, 10, 64)
		if err != nil {
			blog.Errorf("parse shared ns id from redis key %s failed, err: %v", redisKey, err)
			return nil, err
		}

		asstID, err := strconv.ParseInt(redisData, 10, 64)
		if err != nil {
			blog.Errorf("parse asst biz id from redis data %s failed, err: %v", redisData, err)
			return nil, err
		}

		nsAsstBizMap[nsID] = asstID
	}

	return nsAsstBizMap, nil
}

// UpdateAsstBiz update shared namespace to associated biz cache by map[nsID]asstBizID
func (c *SharedNsRelCache) UpdateAsstBiz(ctx context.Context, nsAsstBizMap map[int64]int64, rid string) error {
	redisDataMap := make(map[string]interface{})
	for nsID, asstBizID := range nsAsstBizMap {
		redisDataMap[c.genNsAsstBizRedisKey(nsID)] = asstBizID
	}

	if err := c.nsAsstBizCache.BatchUpdate(ctx, redisDataMap, rid); err != nil {
		blog.Errorf("update shared ns asst biz cache failed, err: %v, data: %+v, rid: %s", err, nsAsstBizMap, rid)
		return err
	}
	return nil
}

// DeleteAsstBiz delete shared namespace to associated biz cache by nsIDs
func (c *SharedNsRelCache) DeleteAsstBiz(ctx context.Context, nsIDs []int64, rid string) error {
	redisKeys := make([]string, len(nsIDs))
	for i, nsID := range nsIDs {
		redisKeys[i] = c.genNsAsstBizRedisKey(nsID)
	}

	if err := c.nsAsstBizCache.BatchDelete(ctx, redisKeys, rid); err != nil {
		blog.Errorf("delete shared ns asst biz cache failed, err: %v, keys: %+v, rid: %s", err, redisKeys, rid)
		return err
	}
	return nil
}

// RefreshSharedNsRel refresh shared namespace relation key and value cache
func (c *SharedNsRelCache) RefreshSharedNsRel(ctx context.Context, rid string) error {
	// lock refresh shared namespace relation cache operation, returns error if it is already locked
	lockKey := fmt.Sprintf("%s:shared_ns_rel_refresh:lock", Namespace)

	locker := lock.NewLocker(redis.Client())
	locked, err := locker.Lock(lock.StrFormat(lockKey), 10*time.Minute)
	defer locker.Unlock()
	if err != nil {
		blog.Errorf("get %s lock failed, err: %v, rid: %s", lockKey, err, rid)
		return err
	}

	if !locked {
		blog.Errorf("%s task is already lock, rid: %s", lockKey, rid)
		return errors.New("there's a same refreshing task running, please retry later")
	}

	ctx = util.SetDBReadPreference(ctx, common.SecondaryPreferredMode)

	relations, err := c.getAllSharedNsRel(ctx, rid)
	if err != nil {
		return err
	}

	nsAsstBizMap := make(map[string]interface{})
	for _, rel := range relations {
		nsAsstBizMap[c.genNsAsstBizRedisKey(rel.NamespaceID)] = rel.AsstBizID
	}

	// refresh label key and value count cache
	err = c.nsAsstBizCache.Refresh(ctx, "*", nsAsstBizMap, rid)
	if err != nil {
		blog.Errorf("refresh shared ns asst biz cache failed, err: %v, data: %+v, rid: %s", err, nsAsstBizMap, rid)
		return err
	}

	return nil
}

// getAllSharedNsRel get all shared namespace relations
func (c *SharedNsRelCache) getAllSharedNsRel(ctx context.Context, rid string) ([]kubetypes.NsSharedClusterRel, error) {
	cond := make(mapstr.MapStr)

	all := make([]kubetypes.NsSharedClusterRel, 0)
	for {
		relations := make([]kubetypes.NsSharedClusterRel, 0)
		err := mongodb.Client().Table(kubetypes.BKTableNameNsSharedClusterRel).Find(cond).
			Sort(kubetypes.BKNamespaceIDField).Fields(kubetypes.BKNamespaceIDField, kubetypes.BKAsstBizIDField).
			All(ctx, &relations)
		if err != nil {
			blog.Errorf("list kube shared namespace rel failed, err: %v, cond: %+v, rid: %v", err, cond, rid)
			return nil, err
		}

		all = append(all, relations...)

		if len(relations) < types.DBPage {
			break
		}

		cond[kubetypes.BKNamespaceIDField] = mapstr.MapStr{common.BKDBGT: relations[len(relations)-1].NamespaceID}
	}
	return all, nil
}

// loopRefreshCache loop refresh shared namespace relation key and value cache every day at 3am
func (c *SharedNsRelCache) loopRefreshCache() {
	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)
	for {
		time.Sleep(2 * time.Hour)

		if !c.isMaster.IsMaster() {
			blog.V(4).Infof("loop refresh shared namespace relation cache, but not master, skip.")
			time.Sleep(time.Minute)
			continue
		}

		rid := util.GenerateRID()

		blog.Infof("start refresh shared namespace relation cache task, rid: %s", rid)
		err := c.RefreshSharedNsRel(ctx, rid)
		if err != nil {
			blog.Errorf("refresh shared namespace relation cache failed, err: %v, rid: %s", err, rid)
			continue
		}
		blog.Infof("refresh shared namespace relation cache successfully, rid: %s", rid)
	}
}
