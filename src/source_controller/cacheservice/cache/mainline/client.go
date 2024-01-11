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

package mainline

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal"
	dalrds "configcenter/src/storage/dal/redis"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream"
)

var client *Client
var clientOnce sync.Once
var cache *cacheCollection

// NewMainlineCache is to initialize a mainline cache handle instance.
// It will start to cache the business's mainline topology instance's
// cache when a instance's event is occurred.
// It's a event triggered and ttl refresh combined cache mechanism.
// which help us to refresh the cache in time with event, and refresh
// cache with ttl without event triggered.
// Note: it can only be called for once.
func NewMainlineCache(event stream.LoopInterface) error {

	if cache != nil {
		return nil
	}

	// cache has not been initialized.
	biz := &business{
		key:   bizKey,
		event: event,
		rds:   redis.Client(),
		db:    mongodb.Client(),
	}

	if err := biz.Run(); err != nil {
		return fmt.Errorf("run biz cache failed, err: %v", err)
	}

	module := &module{
		key:   moduleKey,
		event: event,
		rds:   redis.Client(),
		db:    mongodb.Client(),
	}
	if err := module.Run(); err != nil {
		return fmt.Errorf("run module cache failed, err: %v", err)
	}

	set := &set{
		key:   setKey,
		event: event,
		rds:   redis.Client(),
		db:    mongodb.Client(),
	}
	if err := set.Run(); err != nil {
		return fmt.Errorf("run set cache failed, err: %v", err)
	}

	custom := &customLevel{
		event: event,
		rds:   redis.Client(),
		db:    mongodb.Client(),
	}

	if err := custom.Run(); err != nil {
		return fmt.Errorf("run biz custom level cache failed, err: %v", err)
	}

	cache = &cacheCollection{
		business: biz,
		set:      set,
		module:   module,
		custom:   custom,
	}
	return nil
}

// NewMainlineClient new a mainline cache client, which is used to get the business's
// mainline topology's instance cache.
// this client can only be initialized for once.
func NewMainlineClient() *Client {

	if client != nil {
		return client
	}

	// initialize for once.
	clientOnce.Do(func() {
		client = &Client{
			rds: redis.Client(),
			db:  mongodb.Client(),
		}
	})

	return client
}

// Client is a business's topology cache client instance,
// which is used to get cache from redis and refresh cache
// with ttl policy.
type Client struct {
	rds dalrds.Client
	db  dal.DB
}

// GetBusiness get a business's all info with business id
func (c *Client) GetBusiness(ctx context.Context, bizID int64) (string, error) {
	rid := ctx.Value(common.BKHTTPCCRequestID)

	key := bizKey.detailKey(bizID)
	biz, err := c.rds.Get(ctx, key).Result()
	if err == nil {
		return biz, nil
	}

	blog.Errorf("get business %d info from cache failed, will get from db, err: %v, rid: %v", bizID, err, rid)

	// error occurs, get from db directly.
	biz, err = c.getBusinessFromMongo(bizID)
	if err != nil {
		blog.Errorf("get biz detail from db failed, err: %v, rid: %v", err, rid)
		return "", err
	}

	// refresh biz cache with the latest info.
	err = c.rds.Set(ctx, bizKey.detailKey(bizID), biz, bizKey.detailExpireDuration).Err()
	if err != nil {
		blog.Errorf("update biz cache failed, err: %v, rid: %s", err, rid)
		// do not return, cache will be refresh with next round.
	}

	// get from db
	return biz, nil
}

// ListBusiness list business's cache with options.
func (c *Client) ListBusiness(ctx context.Context, opt *metadata.ListWithIDOption) ([]string, error) {
	rid := ctx.Value(common.ContextRequestIDField)

	if len(opt.IDs) == 0 {
		return make([]string, 0), nil
	}

	keys := make([]string, len(opt.IDs))
	for idx, bizID := range opt.IDs {
		keys[idx] = bizKey.detailKey(bizID)
	}

	bizList, err := c.rds.MGet(context.Background(), keys...).Result()
	if err != nil {
		blog.Errorf("get business %d info from cache failed, get from db directly, err: %v, rid: %v", opt.IDs, err, rid)
		return c.listBusinessWithRefreshCache(ctx, opt.IDs, opt.Fields)
	}

	all := make([]string, 0)
	toAdd := make([]int64, 0)
	for idx, biz := range bizList {
		if biz == nil {
			// can not find in cache
			toAdd = append(toAdd, opt.IDs[idx])
			continue
		}

		detail, ok := biz.(string)
		if !ok {
			blog.Errorf("got invalid biz cache %v, rid: %v", biz, rid)
			return nil, fmt.Errorf("got invalid biz cache %v", biz)
		}

		if len(opt.Fields) != 0 {
			all = append(all, *json.CutJsonDataWithFields(&detail, opt.Fields))
		} else {
			all = append(all, detail)
		}

	}

	if len(toAdd) != 0 {
		// several business caches is not hit, try get from db and refresh them to cache.
		details, err := c.listBusinessWithRefreshCache(ctx, toAdd, opt.Fields)
		if err != nil {
			blog.Errorf("get business list from db failed, err: %v, rid: %v", err, rid)
			return nil, err
		}

		all = append(all, details...)
	}

	return all, nil
}

// ListModules list modules cache with options from redis.
func (c *Client) ListModules(ctx context.Context, opt *metadata.ListWithIDOption) ([]string, error) {
	rid := ctx.Value(common.ContextRequestIDField)

	if len(opt.IDs) == 0 {
		return make([]string, 0), nil
	}

	keys := make([]string, len(opt.IDs))
	list, err := c.rds.MGet(context.Background(), keys...).Result()
	if err != nil {
		blog.Errorf("list module %d info from cache failed, get from db directly, err: %v, rid: %v", opt.IDs, err, rid)
		return c.listModuleWithRefreshCache(ctx, opt.IDs, opt.Fields)
	}

	all := make([]string, 0)
	toAdd := make([]int64, 0)
	for idx, module := range list {
		if module == nil {
			// can not find in cache
			toAdd = append(toAdd, opt.IDs[idx])
			continue
		}

		detail, ok := module.(string)
		if !ok {
			blog.Errorf("got invalid module cache %v, rid: %v", module, rid)
			return nil, fmt.Errorf("got invalid module cache %v", module)
		}

		if len(opt.Fields) != 0 {
			all = append(all, *json.CutJsonDataWithFields(&detail, opt.Fields))
		} else {
			all = append(all, detail)
		}
	}

	if len(toAdd) != 0 {
		// several module caches is not hit, try get from db and refresh them to cache.
		details, err := c.listModuleWithRefreshCache(ctx, toAdd, opt.Fields)
		if err != nil {
			blog.Errorf("get module list from db failed, err: %v, rid: %v", err, rid)
			return nil, err
		}

		all = append(all, details...)
	}

	return all, nil
}

// ListSets list sets from cache with options.
func (c *Client) ListSets(ctx context.Context, opt *metadata.ListWithIDOption) ([]string, error) {
	rid := ctx.Value(common.ContextRequestIDField)

	if len(opt.IDs) == 0 {
		return make([]string, 0), nil
	}

	keys := make([]string, len(opt.IDs))
	for idx, id := range opt.IDs {
		keys[idx] = setKey.detailKey(id)
	}

	list, err := c.rds.MGet(ctx, keys...).Result()
	if err != nil {
		blog.Errorf("list set %d info from cache failed, get from db directly, err: %v, rid: %v", opt.IDs, err, rid)
		return c.listSetWithRefreshCache(ctx, opt.IDs, opt.Fields)
	}

	all := make([]string, 0)
	toAdd := make([]int64, 0)
	for idx, set := range list {
		if set == nil {
			// can not find in cache
			toAdd = append(toAdd, opt.IDs[idx])
			continue
		}

		detail, ok := set.(string)
		if !ok {
			blog.Errorf("got invalid set cache %v, rid: %v", set, rid)
			return nil, fmt.Errorf("got invalid set cache %v", set)
		}

		if len(opt.Fields) != 0 {
			all = append(all, *json.CutJsonDataWithFields(&detail, opt.Fields))
		} else {
			all = append(all, detail)
		}
	}

	if len(toAdd) != 0 {
		// several set caches is not hit, try get from db and refresh them to cache.
		details, err := c.listSetWithRefreshCache(ctx, toAdd, opt.Fields)
		if err != nil {
			blog.Errorf("get set list from db failed, err: %v, rid: %v", err, rid)
			return nil, err
		}

		all = append(all, details...)
	}

	return all, nil
}

// ListModuleDetails list module's all details from cache with module ids.
func (c *Client) ListModuleDetails(ctx context.Context, moduleIDs []int64) ([]string, error) {
	if len(moduleIDs) == 0 {
		return make([]string, 0), nil
	}

	rid := ctx.Value(common.ContextRequestIDField)

	keys := make([]string, len(moduleIDs))
	for idx, id := range moduleIDs {
		keys[idx] = moduleKey.detailKey(id)
	}

	modules, err := c.rds.MGet(context.Background(), keys...).Result()
	if err == nil {
		list := make([]string, 0)
		for idx, m := range modules {
			if m == nil {
				detail, isNotFound, err := c.getModuleDetailCheckNotFoundWithRefreshCache(ctx, moduleIDs[idx])
				// 跳过不存在的模块，因为作为批量查询的API，调用方希望查询到存在的资源，并自动过滤掉不存在的资源
				if isNotFound {
					blog.Errorf("module %d not exist, err: %v, rid: %v", moduleIDs[idx], err, rid)
					continue
				}

				if err != nil {
					blog.Errorf("get module %d detail from db failed, err: %v, rid: %v", moduleIDs[idx], err, rid)
					return nil, err
				}

				list = append(list, detail)
				continue
			}
			list = append(list, m.(string))
		}
		return list, nil
	}
	blog.Errorf("get modules details from redis failed, err: %v, rid: %v", err, rid)

	// can not get from redis, get from db directly and refresh cache.
	return c.listModuleWithRefreshCache(ctx, moduleIDs, nil)
}

// GetModuleDetail get a module's details with id from cache.
func (c *Client) GetModuleDetail(ctx context.Context, moduleID int64) (string, error) {
	rid := ctx.Value(common.ContextRequestIDField)

	mod, err := c.rds.Get(ctx, moduleKey.detailKey(moduleID)).Result()
	if err == nil {
		return mod, nil
	}

	blog.Errorf("get module: %d failed from redis, err: %v, rid: %v", moduleID, err, rid)
	// get from db directly and refresh the cache.
	detail, _, err := c.getModuleDetailCheckNotFoundWithRefreshCache(ctx, moduleID)
	return detail, err
}

// GetSet get a set's details from cache with id.
func (c *Client) GetSet(ctx context.Context, setID int64) (string, error) {
	rid := ctx.Value(common.ContextRequestIDField)

	set, err := c.rds.Get(context.Background(), setKey.detailKey(setID)).Result()
	if err == nil {
		return set, nil
	}

	blog.Errorf("get set: %d failed from redis failed, err: %v, rid: %v", setID, err, rid)

	// can not get set from cache, get from db directly and refresh cache.
	detail, _, err := c.getSetDetailCheckNotFoundWithRefreshCache(ctx, setID)
	return detail, err
}

// ListSetDetails list set's details from cache with ids.
func (c *Client) ListSetDetails(ctx context.Context, setIDs []int64) ([]string, error) {
	rid := ctx.Value(common.ContextRequestIDField)

	if len(setIDs) == 0 {
		return make([]string, 0), nil
	}

	keys := make([]string, len(setIDs))
	for idx, set := range setIDs {
		keys[idx] = setKey.detailKey(set)
	}

	sets, err := c.rds.MGet(context.Background(), keys...).Result()
	if err == nil && len(sets) != 0 {
		all := make([]string, 0)
		for idx, s := range sets {
			if s == nil {
				detail, isNotFound, err := c.getSetDetailCheckNotFoundWithRefreshCache(ctx, setIDs[idx])
				// 跳过不存在的集群，因为作为批量查询的API，调用方希望查询到存在的资源，并自动过滤掉不存在的资源
				if isNotFound {
					blog.Errorf("set %d not exist, err: %v, rid: %v", setIDs[idx], err, rid)
					continue
				}

				if err != nil {
					blog.Errorf("get set %d from mongodb failed, err: %v, rid: %v", setIDs[idx], err, rid)
					return nil, err
				}
				all = append(all, detail)
				continue
			}
			all = append(all, s.(string))
		}

		return all, nil
	}

	blog.Errorf("get sets: %v failed from redis failed, err: %v, rid: %v", setIDs, err, rid)

	// get from db directly and refresh the cache.
	return c.listSetWithRefreshCache(ctx, setIDs, nil)
}

// GetCustomLevelDetail get business's custom level object's instance detail information with instance id.
func (c *Client) GetCustomLevelDetail(ctx context.Context, objID, supplierAccount string, instID int64) (
	string, error) {

	rid := ctx.Value(common.ContextRequestIDField)
	key := newCustomKey(objID)
	custom, err := c.rds.Get(context.Background(), key.detailKey(instID)).Result()
	if err == nil {
		return custom, nil
	}

	blog.Errorf("get biz custom level, obj:%s, inst: %d failed from redis, err: %v, rid: %v", objID, instID, err, rid)

	detail, _, err := c.getCustomDetailCheckNotFoundWithRefreshCache(ctx, key, objID, supplierAccount, instID)
	return detail, err
}

// ListCustomLevelDetail business's custom level object's instance detail information with id list.
func (c *Client) ListCustomLevelDetail(ctx context.Context, objID, supplierAccount string, instIDs []int64) (
	[]string, error) {

	if len(instIDs) == 0 {
		return make([]string, 0), nil
	}

	rid := ctx.Value(common.ContextRequestIDField)

	customKey := newCustomKey(objID)
	keys := make([]string, len(instIDs))
	for idx, instID := range instIDs {
		keys[idx] = customKey.detailKey(instID)
	}

	customs, err := c.rds.MGet(context.Background(), keys...).Result()
	if err == nil && len(customs) != 0 {
		all := make([]string, 0)
		for idx, cu := range customs {
			if cu == nil {
				detail, isNotFound, err := c.getCustomDetailCheckNotFoundWithRefreshCache(ctx, customKey, objID,
					supplierAccount, instIDs[idx])
				// 跳过不存在的自定义节点，因为作为批量查询的API，调用方希望查询到存在的资源，并自动过滤掉不存在的资源
				if isNotFound {
					blog.Errorf("custom layer %s/%d not exist, err: %v, rid: %v", objID, instIDs[idx], err, rid)
					continue
				}

				if err != nil {
					blog.Errorf("get %s/%d detail from mongodb failed, err: %v, rid: %v", objID, instIDs[idx], err, rid)
					return nil, err
				}
				all = append(all, detail)
				continue
			}

			all = append(all, cu.(string))
		}
		return all, nil
	}

	blog.Errorf("get biz custom level, obj:%s, inst: %v failed from redis, err: %v, rid: %v", objID, instIDs, err, rid)
	// get from db directly and try refresh the cache.
	return c.listCustomLevelDetailWithRefreshCache(ctx, customKey, objID, supplierAccount, instIDs)
}

// GetTopology get business's mainline topology with rank from biz model to host model.
func (c *Client) GetTopology() ([]string, error) {

	rank, err := c.rds.Get(context.Background(), topologyKey).Result()
	if err != nil {
		blog.Errorf("get mainline topology from cache failed, get from db directly. err: %v", err)
		return c.refreshAndGetTopologyRank()
	}

	topo := strings.Split(rank, ",")
	if len(topo) < 4 {
		// invalid topology
		return c.refreshAndGetTopologyRank()
	}

	return topo, nil
}
