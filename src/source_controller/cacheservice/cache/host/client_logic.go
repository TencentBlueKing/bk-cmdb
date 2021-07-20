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

package host

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"

	rawRedis "github.com/go-redis/redis/v7"
)

func (c *Client) tryRefreshHostDetail(hostID int64, ips string, cloudID int64, detail []byte) {
	hostKey := hostKey.HostDetailKey(hostID)
	if !c.lock.CanRefresh(hostKey) {
		return
	}
	// set refreshing status
	c.lock.SetRefreshing(hostKey)

	// now, we check whether we can refresh the host detail cache
	go func() {
		defer func() {
			c.lock.SetUnRefreshing(hostKey)
		}()

		refreshHostDetailCache(hostID, ips, cloudID, detail)
	}()
}

func (c *Client) tryRefreshHostIDList(rid string) {
	ctx := context.WithValue(context.Background(), common.ContextRequestIDField, rid)

	// get local lock
	if !c.lock.CanRefresh(hostKey.HostIDListLockKey()) {
		return
	}

	// set refreshing status
	c.lock.SetRefreshing(hostKey.HostIDListLockKey())
	defer c.lock.SetUnRefreshing(hostKey.HostIDListLockKey())

	forceRefresh := false
	expire, err := redis.Client().Get(ctx, hostKey.HostIDListExpireKey()).Result()
	if err != nil {
		if !redis.IsNilErr(err) {
			blog.Errorf("get host id list expire key failed, err: %v, rid :%v", err, rid)
			return
		} else {
			// force refresh
			forceRefresh = true
		}
	}

	if !forceRefresh {
		expireAt, err := strconv.ParseInt(expire, 10, 64)
		if err != nil {
			blog.Errorf("parse host id list expire time failed, err: %v, rid: %v", err, rid)
			return
		}

		if time.Now().Unix()-expireAt > 0 && (time.Now().Unix()-expireAt <= hostKey.HostIDListKeyExpireSeconds()) {
			// not expired
			return
		}
	}

	// set expire key with a value which can will enforce the host id list key expire in a minute, which will
	// help us block the refresh request in the following one minute. this policy is used to avoid refresh key
	// when the redis in high pressure or not well performed.
	redis.Client().Set(ctx, hostKey.HostIDListExpireKey(), time.Now().Unix()+60,
		time.Duration(hostKey.HostIDListKeyExpireSeconds())*time.Second)

	// expired, we refresh it now.
	blog.Infof("host id list key: %s is expired, refresh it now. rid: %v", hostKey.HostIDListKey(), rid)

	// then get distribute lock.
	locked, err := redis.Client().SetNX(ctx, hostKey.HostIDListLockKey(), true,
		time.Duration(hostKey.HostIDListKeyExpireSeconds())*time.Second).Result()

	if err != nil {
		blog.Errorf("get host id list key lock failed, err: %v, rid: %v", err, rid)
		return
	}

	if !locked {
		// locked by others, skip refresh operation
		return
	}

	go func() {
		// already get lock, force refresh host id list now.
		c.refreshHostIDListCache(rid)

		if err := redis.Client().Del(ctx, hostKey.HostIDListLockKey()).Err(); err != nil {
			blog.Errorf("delete host id list lock key: %s failed, err: %v, rid: %v", hostKey.HostIDListLockKey(),
				err, rid)
		}
	}()

}

func (c *Client) forceRefreshHostIDList(ctx context.Context) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// get local lock
	if !c.lock.CanRefresh(hostKey.HostIDListLockKey()) {
		return
	}

	// set refreshing status
	c.lock.SetRefreshing(hostKey.HostIDListLockKey())
	defer c.lock.SetUnRefreshing(hostKey.HostIDListLockKey())

	// then get distribute lock.
	locked, err := redis.Client().SetNX(context.Background(), hostKey.HostIDListLockKey(), true,
		time.Duration(hostKey.HostIDListKeyExpireSeconds())*time.Second).Result()

	if err != nil {
		blog.Errorf("get host id list key lock failed, err: %v, rid: %v", err, rid)
		return
	}

	if !locked {
		return
	}

	blog.Infof("start force fresh host id list key: %s. rid: %v", hostKey.HostIDListKey(), rid)

	go func() {
		// already get lock, force refresh host id list now.
		c.refreshHostIDListCache(rid)

		if err := redis.Client().Del(context.Background(), hostKey.HostIDListLockKey()).Err(); err != nil {
			blog.Errorf("delete host id list lock key: %s failed, err: %v, rid: %v", hostKey.HostIDListLockKey(),
				err, rid)
		}
	}()

}

type hostID struct {
	ID int64 `bson:"bk_host_id"`
}

const (
	// step to list host ids from mongodb
	listStep = 50000
)

func (c *Client) refreshHostIDListCache(rid string) error {
	ctx := context.Background()

	// get all host id list at first
	total, err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(nil).Count(ctx)
	if err != nil {
		blog.Errorf("refresh host id list, but count failed, err: %v, rid: %v", err, rid)
		return err
	}

	tempIDListKey := fmt.Sprintf("%s-%s", hostKey.HostIDListTempKey(), rid)
	blog.Infof("try to refresh host id list with temp key: %s, rid: %s", tempIDListKey, rid)

	for start := 0; start < int(total); start += listStep {
		stepID := make([]hostID, 0)
		if err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(nil).Start(uint64(start)).
			Limit(uint64(listStep)).Fields(common.BKHostIDField).All(ctx, &stepID); err != nil {
			blog.Errorf("refresh host id list, but get host id list failed, err: %v, rid: %v", err, rid)
			return err
		}

		pip := redis.Client().Pipeline()
		// because the temp key is a random key, so we set a expire time so that it can be gc,
		// but we will reset expire to unlimited when this key is renamed to a normal key.
		pip.Expire(tempIDListKey, time.Duration(hostKey.HostIDListKeyExpireSeconds())*time.Second)
		for _, h := range stepID {
			key := &rawRedis.Z{
				// set zset score with host id, so we can sort with host id
				Score: float64(h.ID),
				// set zset member with host id, so that we can get host id with score directly.
				Member: h.ID,
			}
			// write to the temp key
			pip.ZAdd(tempIDListKey, key)
		}

		// it cost about 600ms to zadd 100000 host id to redis from test case.
		_, err := pip.Exec()
		if err != nil {
			blog.Errorf("update host id list failed, err: %v, rid: %v", err, rid)
			return err
		}
	}

	pipe := redis.Client().Pipeline()
	// rename temp key to real key
	pipe.Rename(tempIDListKey, hostKey.HostIDListKey())
	// reset id_list key's expire time to a new one.
	pipe.Expire(hostKey.HostIDListKey(), 48*time.Hour)
	// set expire key with unix time seconds now value.
	pipe.Set(hostKey.HostIDListExpireKey(), time.Now().Unix(),
		time.Duration(hostKey.HostIDListKeyExpireSeconds())*time.Second)

	if _, err := pipe.Exec(); err != nil {
		blog.Errorf("rename host id list key form %s to %s failed, err :%v, rid: %v", hostKey.HostIDListTempKey(),
			hostKey.HostIDListKey(), err, rid)
		return err
	}

	blog.Infof("fresh host id list key: %s success, count: %d. rid: %v", hostKey.HostIDListKey(), total, rid)

	return nil
}

func (c *Client) getHostsWithPage(ctx context.Context, opt *metadata.ListHostWithPage) (int64, []string, error) {
	rid := ctx.Value(common.ContextRequestIDField)

	total, err := c.countHost(ctx, nil)
	if err != nil {
		blog.Errorf("get host with page, but count failed, err: %v, rid: %v", err, rid)
		return 0, nil, err
	}

	list := make([]metadata.HostMapStr, 0)
	if err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(nil).Start(uint64(opt.Page.Start)).
		Limit(uint64(opt.Page.Limit)).Sort(common.BKHostIDField).Fields(opt.Fields...).All(ctx, &list); err != nil {

		blog.Errorf("get host id list with page failed, err: %v, rid: %v", err, rid)
		return 0, nil, err
	}

	all := make([]string, len(list))
	for idx := range list {
		// the err can be ignore because it's unmarshal from bson upper, marshal it again is also available.
		js, _ := json.Marshal(list[idx])
		all[idx] = string(js)
	}

	return int64(total), all, nil
}

func (c *Client) countHost(ctx context.Context, filter map[string]interface{}) (uint64, error) {
	return mongodb.Client().Table(common.BKTableNameBaseHost).Find(filter).Count(ctx)
}
