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

package topology

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/driver/mongodb"
	drvRedis "configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream/types"
)

func newTokenHandler(object string) *tokenHandler {
	return &tokenHandler{
		doc: "brief_topology_cache_watch_token",
		key: object,
		db:  mongodb.Client(),
	}
}

type tokenHandler struct {
	doc string
	key string
	db  dal.DB
}

func (w *tokenHandler) SetLastWatchToken(ctx context.Context, token string) error {
	var err error
	// do with retry
	filter := map[string]interface{}{"_id": w.doc}
	tokenData := mapstr.MapStr{w.key: token}

	for try := 0; try < 5; try++ {
		err = w.db.Table(common.BKTableNameSystem).Upsert(ctx, filter, tokenData)
		if err != nil {
			time.Sleep(time.Duration(try/2+1) * time.Second)
			continue
		}
		return nil
	}

	return err
}

// get the former watched token.
// if Key is not exist, then token is "".
func (w *tokenHandler) GetStartWatchToken(ctx context.Context) (token string, err error) {
	// do with retry
	filter := map[string]interface{}{"_id": w.doc}
	for try := 0; try < 5; try++ {
		tokenData := make(map[string]string)
		err = w.db.Table(common.BKTableNameSystem).Find(filter).Fields(w.key).One(ctx, &tokenData)
		if err != nil {
			blog.Errorf("get %s start token failed, err: %v", w.key, err)
			if !w.db.IsNotFoundError(err) {
				time.Sleep(time.Duration(try/2+1) * time.Second)
				continue
			}
			return "", nil
		}
		return tokenData[w.key], nil
	}

	return "", err
}

// resetWatchToken set watch token to empty and set the start watch time to the given one for next watch
func (w *tokenHandler) resetWatchToken(startAtTime types.TimeStamp) error {
	filter := map[string]interface{}{"_id": w.doc}
	tokenData := mapstr.MapStr{
		w.key:                 "",
		w.key + "_start_time": startAtTime,
	}

	return w.db.Table(common.BKTableNameSystem).Upsert(context.Background(), filter, tokenData)
}

func (w *tokenHandler) getStartWatchTime(ctx context.Context) (*types.TimeStamp, error) {
	filter := map[string]interface{}{"_id": w.doc}

	data := make(map[string]types.TimeStamp)
	err := w.db.Table(common.BKTableNameSystem).Find(filter).Fields(w.key+"_start_time").One(ctx, &data)
	if err != nil {
		if !w.db.IsNotFoundError(err) {
			blog.Errorf("get %s start time failed, err: %v", w.key, err)
			return nil, err
		}
		return new(types.TimeStamp), nil
	}
	startTime := data[w.key+"_start_time"]
	return &startTime, nil
}

func newTopologyKey() *cacheKey {
	return &cacheKey{
		namespace: common.BKCacheKeyV3Prefix + "topology:brief",
		ttl:       24 * time.Hour,
		rds:       drvRedis.Client(),
	}
}

type cacheKey struct {
	name      string
	namespace string
	ttl       time.Duration
	rds       redis.Client
}

func (c *cacheKey) bizTopologyKey(biz int64) string {
	return fmt.Sprintf("%s:%d", c.namespace, biz)
}

// updateTopology update biz Topology cache
func (c *cacheKey) updateTopology(ctx context.Context, topo *BizBriefTopology) error {

	js, err := json.Marshal(topo)
	if err != nil {
		return fmt.Errorf("marshal topology failed, err: %v", err)
	}

	return c.rds.Set(ctx, c.bizTopologyKey(topo.Biz.ID), string(js), c.ttl).Err()
}

// getTopology get biz Topology from cache
func (c *cacheKey) getTopology(ctx context.Context, biz int64) (*string, error) {
	dat, err := c.rds.Get(ctx, c.bizTopologyKey(biz)).Result()
	if err != nil {
		if redis.IsNilErr(err) {
			empty := ""
			return &empty, nil
		}

		return nil, fmt.Errorf("get cache from redis failed, err: %v", err)
	}

	return &dat, nil
}
