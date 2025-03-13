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

// Package mainline is the mainline instance cache
package mainline

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"configcenter/pkg/cache/general"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	generalcache "configcenter/src/source_controller/cacheservice/cache/general"
	"configcenter/src/storage/driver/redis"
)

var client *Client
var clientOnce sync.Once
var cache *mainlineCache

// NewMainlineCache is to initialize a mainline cache handle instance.
// Note: it can only be called for once.
func NewMainlineCache(isMaster discovery.ServiceManageInterface) error {
	if cache != nil {
		return nil
	}

	cache = &mainlineCache{
		isMaster: isMaster,
	}

	if err := cache.Run(); err != nil {
		return fmt.Errorf("run mainline cache failed, err: %v", err)
	}

	return nil
}

// NewMainlineClient new a mainline cache client, which is used to get the business's
// mainline topology's instance cache.
// this client can only be initialized for once.
func NewMainlineClient(cache *generalcache.Cache) *Client {
	if client != nil {
		return client
	}

	// initialize for once.
	clientOnce.Do(func() {
		client = &Client{
			cache: cache,
		}
	})

	return client
}

// Client is a business's topology cache client instance,
// which is used to get cache from redis and refresh cache
// with ttl policy.
type Client struct {
	cache *generalcache.Cache
}

// GetBusiness get a business's all info with business id
func (c *Client) GetBusiness(kit *rest.Kit, bizID int64) (string, error) {
	return c.getCacheDetailByID(kit, general.Biz, "", bizID)
}

// getCacheDetailByID get resource cache detail by id
func (c *Client) getCacheDetailByID(kit *rest.Kit, resource general.ResType, subRes string, id int64) (string, error) {
	listOpt := &general.ListDetailByIDsOpt{
		Resource:    resource,
		SubResource: subRes,
		IDs:         []int64{id},
	}
	details, err := c.cache.ListDetailByIDs(kit, listOpt)
	if err != nil {
		blog.Errorf("get %s:%s %d from cache failed, err: %v, rid: %v", resource, subRes, id, err, kit.Rid)
		return "", err
	}

	if len(details) != 1 || details[0] == "" {
		blog.Errorf("%s:%s %d cache detail %+v is invalid, rid: %s", resource, subRes, id, details, kit.Rid)
		return "", kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "cache detail")
	}
	return details[0], nil
}

// ListBusiness list business's cache with options.
func (c *Client) ListBusiness(kit *rest.Kit, opt *metadata.ListWithIDOption) ([]string, error) {
	return c.listCacheDetail(kit, general.Biz, "", opt)
}

// listCacheDetail list resource cache detail
func (c *Client) listCacheDetail(kit *rest.Kit, res general.ResType, subRes string, opt *metadata.ListWithIDOption) (
	[]string, error) {

	if len(opt.IDs) == 0 {
		return make([]string, 0), nil
	}

	listOpt := &general.ListDetailByIDsOpt{
		Resource:    res,
		SubResource: subRes,
		IDs:         opt.IDs,
		Fields:      opt.Fields,
	}
	details, err := c.cache.ListDetailByIDs(kit, listOpt)
	if err != nil {
		blog.Errorf("list %s:%s from cache failed, err: %v, opt: %+v, rid: %v", res, subRes, err, opt, kit.Rid)
		return nil, err
	}

	return details, nil
}

// ListModules list modules cache with options from redis.
func (c *Client) ListModules(kit *rest.Kit, opt *metadata.ListWithIDOption) ([]string, error) {
	return c.listCacheDetail(kit, general.Module, "", opt)
}

// ListSets list sets from cache with options.
func (c *Client) ListSets(kit *rest.Kit, opt *metadata.ListWithIDOption) ([]string, error) {
	return c.listCacheDetail(kit, general.Set, "", opt)
}

// ListModuleDetails list module's all details from cache with module ids.
func (c *Client) ListModuleDetails(kit *rest.Kit, moduleIDs []int64) ([]string, error) {
	return c.ListModules(kit, &metadata.ListWithIDOption{IDs: moduleIDs})
}

// GetModuleDetail get a module's details with id from cache.
func (c *Client) GetModuleDetail(kit *rest.Kit, moduleID int64) (string, error) {
	return c.getCacheDetailByID(kit, general.Module, "", moduleID)
}

// GetSet get a set's details from cache with id.
func (c *Client) GetSet(kit *rest.Kit, setID int64) (string, error) {
	return c.getCacheDetailByID(kit, general.Set, "", setID)
}

// ListSetDetails list set's details from cache with ids.
func (c *Client) ListSetDetails(kit *rest.Kit, setIDs []int64) ([]string, error) {
	return c.ListSets(kit, &metadata.ListWithIDOption{IDs: setIDs})
}

// GetCustomLevelDetail get business's custom level object's instance detail information with instance id.
func (c *Client) GetCustomLevelDetail(kit *rest.Kit, objID string, instID int64) (string, error) {
	return c.getCacheDetailByID(kit, general.MainlineInstance, objID, instID)
}

// ListCustomLevelDetail business's custom level object's instance detail information with id list.
func (c *Client) ListCustomLevelDetail(kit *rest.Kit, objID string, instIDs []int64) ([]string, error) {
	return c.listCacheDetail(kit, general.MainlineInstance, objID, &metadata.ListWithIDOption{IDs: instIDs})
}

// GetTopology get business's mainline topology with rank from biz model to host model.
func (c *Client) GetTopology(kit *rest.Kit) ([]string, error) {
	rank, err := redis.Client().Get(context.Background(), genTopologyKey(kit)).Result()
	if err != nil {
		blog.Errorf("get mainline topology from cache failed, get from db directly. err: %v, rid: %s", err, kit.Rid)
		return refreshAndGetTopologyRank(kit)
	}

	topo := strings.Split(rank, ",")
	if len(topo) < 4 {
		// invalid topology
		return refreshAndGetTopologyRank(kit)
	}

	return topo, nil
}
