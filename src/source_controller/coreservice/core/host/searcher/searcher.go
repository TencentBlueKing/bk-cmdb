/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

// Package searcher TODO
package searcher

import (
	"context"
	"fmt"
	"strings"
	"time"

	"configcenter/src/apimachinery/cacheservice/cache/host"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
)

// Searcher TODO
type Searcher struct {
	CacheHost host.Interface
}

// New TODO
func New(cacheHost host.Interface) Searcher {
	return Searcher{
		CacheHost: cacheHost,
	}
}

func (s *Searcher) buildListHostsFilter(kit *rest.Kit, option metadata.ListHosts) ([]map[string]interface{},
	[]interface{}, bool, bool, error) {

	needHostIDFilter, hostIDs, err := s.listHostIDsByRelation(kit, &option)
	if err != nil {
		return nil, nil, false, false, err
	}

	filters := make([]map[string]interface{}, 0)
	if needHostIDFilter {
		if len(hostIDs) == 0 {
			return nil, nil, true, false, nil
		}
		filters = append(filters, map[string]interface{}{
			common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDs},
		})
	}

	propertyFilter, err := option.GetHostPropertyFilter(kit.Ctx)
	if err != nil {
		blog.Errorf("get host property filter failed, filter: %+v, err: %v, rid: %s", option, err, kit.Rid)
		return nil, nil, false, false, err
	}
	if len(propertyFilter) > 0 {
		filters = append(filters, propertyFilter)
	}

	return filters, hostIDs, false, needHostIDFilter, nil
}

// ListHosts search host with topo node info and host property
func (s *Searcher) ListHosts(kit *rest.Kit, option metadata.ListHosts) (*metadata.ListHostResult, error) {
	if key, err := option.Validate(); err != nil {
		blog.Errorf("list hosts option %+v is invalid, key: %s, err: %v, rid: %s", option, key, err, kit.Rid)
		return nil, err
	}

	filters, hostIDs, noHostsMatched, needHostIDFilter, err := s.buildListHostsFilter(kit, option)
	if err != nil {
		return nil, err
	}
	if noHostsMatched {
		return &metadata.ListHostResult{Count: 0, Info: make([]map[string]interface{}, 0)}, nil
	}

	finalFilter := map[string]interface{}{}
	finalFilter = util.SetQueryOwner(finalFilter, util.ExtractOwnerFromContext(kit.Ctx))
	if len(filters) > 0 {
		finalFilter[common.BKDBAND] = filters
	}

	if needHostIDFilter && len(filters) == 1 && option.BizID != 0 {
		sort := strings.TrimLeft(option.Page.Sort, "+-")
		if len(option.Page.Sort) == 0 || sort == common.BKHostIDField || strings.Contains(sort, ",") == false &&
			strings.HasPrefix(sort, common.BKHostIDField+":") {
			searchResult, err := s.listAllBizHostsPage(kit.Ctx, option.Fields, option.Page, hostIDs, kit.Rid)
			if err != nil {
				return nil, err
			}
			return searchResult, nil
		}

	}

	if len(filters) == 0 {
		// return info use cache
		// fix: has question when multi-supplier
		sort := strings.TrimLeft(option.Page.Sort, "+-")
		if len(option.Page.Sort) == 0 || sort == common.BKHostIDField || strings.Contains(sort, ",") == false &&
			strings.HasPrefix(sort, common.BKHostIDField+":") {
			searchResult, skip, err := s.ListHostsWithCache(kit, option.Fields, option.Page)
			if err != nil {
				return nil, err
			}
			if skip == false {
				return searchResult, nil
			}
		}
	}

	return s.listHostFromDB(kit, finalFilter, &option)
}

func (s *Searcher) listHostIDsByRelation(kit *rest.Kit, option *metadata.ListHosts) (bool, []interface{}, error) {
	relationFilter := map[string]interface{}{}
	var cacheKey string
	if option.BizID != 0 {
		relationFilter[common.BKAppIDField] = option.BizID
		cacheKey = fmt.Sprintf("%sbiz_host_ids:%d", common.BKCacheKeyV3Prefix, option.BizID)
	}

	if len(option.ModuleIDs) != 0 {
		relationFilter[common.BKModuleIDField] = map[string]interface{}{
			common.BKDBIN: option.ModuleIDs,
		}
		cacheKey = ""
	}

	if len(option.SetIDs) != 0 {
		relationFilter[common.BKSetIDField] = map[string]interface{}{
			common.BKDBIN: option.SetIDs,
		}
		cacheKey = ""
	}

	if len(relationFilter) == 0 {
		return false, make([]interface{}, 0), nil
	}

	// get biz host ids from cache
	if cacheKey != "" {
		hostIDsCache, err := redis.Client().Get(context.Background(), cacheKey).Result()
		if err == nil {
			hostIDs := make([]interface{}, 0)
			err = json.UnmarshalFromString(hostIDsCache, &hostIDs)
			if err == nil {
				return true, hostIDs, nil
			}
			blog.Errorf("unmarshal host ids from cache(%s) failed, err: %v, rid: %s", hostIDsCache, err, kit.Rid)
		}

		if !redis.IsNilErr(err) {
			blog.Errorf("get biz host ids from redis failed, key: %s, err: %v, rid: %s", cacheKey, err, kit.Rid)
		}
	}

	hostIDs, err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Distinct(kit.Ctx,
		common.BKHostIDField, relationFilter)
	if err != nil {
		blog.Errorf("get host ids by relation failed, filter: %+v, err: %v, rid: %s", relationFilter, err, kit.Rid)
		return false, nil, err
	}

	// set biz host ids to cache
	if cacheKey != "" {
		// default cache expiration is 2s, we set the expiration to 10s for big bizs with more than 5w hosts
		expiration := 2 * time.Second
		if len(hostIDs) > 50000 {
			expiration = 10 * time.Second
		}

		hostIDCache, err := json.MarshalToString(hostIDs)
		if err != nil {
			blog.Errorf("marshal host ids(%+v) failed, err: %v, rid: %s", hostIDs, err, kit.Rid)
			return true, hostIDs, nil
		}
		if err = redis.Client().Set(context.Background(), cacheKey, hostIDCache, expiration).Err(); err != nil {
			blog.Errorf("set biz host ids to redis failed, key: %s, err: %v, rid: %s", cacheKey, err, kit.Rid)
		}
	}

	return true, hostIDs, nil
}

func (s *Searcher) listHostFromDB(kit *rest.Kit, finalFilter map[string]interface{}, option *metadata.ListHosts) (
	*metadata.ListHostResult, error) {

	total, err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(finalFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count hosts failed, filter: %+v, err: %v, rid: %s", finalFilter, err, kit.Rid)
		return nil, err
	}

	searchResult := &metadata.ListHostResult{Count: int(total)}

	limit := uint64(option.Page.Limit)
	start := uint64(option.Page.Start)
	query := mongodb.Client().Table(common.BKTableNameBaseHost).Find(finalFilter).Limit(limit).Start(start).
		Fields(option.Fields...)
	if len(option.Page.Sort) > 0 {
		query = query.Sort(option.Page.Sort)
	} else {
		query = query.Sort(common.BKHostIDField)
	}

	hosts := make([]metadata.HostMapStr, 0)
	if err = query.All(kit.Ctx, &hosts); err != nil {
		blog.Errorf("list hosts from db failed, filter: %+v, err: %v, rid: %s", finalFilter, err, kit.Rid)
		return nil, err
	}

	searchResult.Info = make([]map[string]interface{}, len(hosts))
	for index, host := range hosts {
		searchResult.Info[index] = host
	}
	return searchResult, nil
}

// ListHostsWithCache TODO
func (s *Searcher) ListHostsWithCache(kit *rest.Kit, fields []string, page metadata.BasePage) (
	searchResult *metadata.ListHostResult, skip bool, err error) {
	opt := &metadata.ListHostWithPage{
		Fields: fields,
		Page:   page,
	}
	total, infos, err := s.CacheHost.ListHostWithPage(kit.Ctx, kit.Header, opt)
	if err != nil {
		blog.Errorf("list host from redis failed, filter: %+v, err: %v, rid: %s", opt, err, kit.Rid)
		// HOOK: 缓存出现错误从db中获取数据
		return nil, true, nil
	}

	searchResult = &metadata.ListHostResult{}
	searchResult.Count = int(total)

	searchResult.Info = make([]map[string]interface{}, 0)
	if err := json.UnmarshalFromString(infos, &searchResult.Info); err != nil {
		blog.Errorf("list host from redis failed, filter: %+v, item host: %s, err: %v, rid: %s", opt, infos, err,
			kit.Rid)
		// TODO： use cc error. keep the same as before code
		return nil, false, err
	}

	return searchResult, false, nil
}

// listAllBizHostsPage 专有流程
func (s *Searcher) listAllBizHostsPage(ctx context.Context, fields []string, page metadata.BasePage,
	allBizHostIDs []interface{}, rid string) (searchResult *metadata.ListHostResult, err error) {
	cnt := len(allBizHostIDs)
	start := page.Start
	if start > cnt {
		return &metadata.ListHostResult{
			Count: 0,
			Info:  []map[string]interface{}{},
		}, nil
	}
	if start < 0 {
		start = 0
	}
	end := start + page.Limit
	if end > cnt {
		end = cnt
	}
	blog.V(10).Infof("list biz host start: %d, end: %d, rid: %s", start, end, rid)

	finalFilter := make(map[string]interface{}, 0)
	// 针对获取业务下所有的主机的场景优
	finalFilter[common.BKHostIDField] = map[string]interface{}{
		common.BKDBIN: allBizHostIDs[start:end],
	}
	finalFilter = util.SetQueryOwner(finalFilter, util.ExtractOwnerFromContext(ctx))

	hosts := make([]metadata.HostMapStr, 0)
	if err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(finalFilter).Fields(fields...).
		Sort(page.Sort).All(ctx, &hosts); err != nil {
		blog.Errorf("ListHosts failed, db select hosts failed, filter: %+v, err: %+v, rid: %s", finalFilter, err, rid)
		return nil, err
	}
	searchResult = &metadata.ListHostResult{
		Count: cnt,
	}
	searchResult.Info = make([]map[string]interface{}, len(hosts))
	for index, host := range hosts {
		searchResult.Info[index] = host
	}
	return searchResult, nil
}
