/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package searcher TODO
package searcher

import (
	"strings"

	"configcenter/src/apimachinery/cacheservice/cache/host"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
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

// ListHosts search host with topo node info and host property
func (s *Searcher) ListHosts(kit *rest.Kit, option metadata.ListHosts) (*metadata.ListHostResult, error) {
	if key, err := option.Validate(); err != nil {
		blog.Errorf("list hosts option %+v is invalid, key: %s, err: %v, rid: %s", option, key, err, kit.Rid)
		return nil, err
	}

	needHostIDFilter, hostIDs, err := s.listHostIDsByRelation(kit, &option)
	if err != nil {
		return nil, err
	}

	filters := make([]map[string]interface{}, 0)
	if needHostIDFilter {
		if len(hostIDs) == 0 {
			return &metadata.ListHostResult{Count: 0, Info: make([]map[string]interface{}, 0)}, nil
		}
		filters = append(filters, map[string]interface{}{
			common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDs},
		})
	}

	propertyFilter, err := option.GetHostPropertyFilter(kit.Ctx)
	if err != nil {
		blog.Errorf("get host property filter failed, filter: %+v, err: %v, rid: %s", option, err, kit.Rid)
		return nil, err
	}
	if len(propertyFilter) > 0 {
		filters = append(filters, propertyFilter)
	}

	finalFilter := map[string]interface{}{}
	if len(filters) > 0 {
		finalFilter[common.BKDBAND] = filters
	}

	if needHostIDFilter && len(filters) == 1 && option.BizID != 0 {
		sort := strings.TrimLeft(option.Page.Sort, "+-")
		if len(option.Page.Sort) == 0 || sort == common.BKHostIDField || strings.Contains(sort, ",") == false &&
			strings.HasPrefix(sort, common.BKHostIDField+":") {
			searchResult, err := s.listAllBizHostsPage(kit, option.Fields, option.Page, hostIDs)
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
	if option.BizID != 0 {
		relationFilter[common.BKAppIDField] = option.BizID
	}

	if len(option.ModuleIDs) != 0 {
		relationFilter[common.BKModuleIDField] = map[string]interface{}{
			common.BKDBIN: option.ModuleIDs,
		}
	}

	if len(option.SetIDs) != 0 {
		relationFilter[common.BKSetIDField] = map[string]interface{}{
			common.BKDBIN: option.SetIDs,
		}
	}

	if len(relationFilter) == 0 {
		return false, make([]interface{}, 0), nil
	}

	hostIDs, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameModuleHostConfig).Distinct(kit.Ctx,
		common.BKHostIDField, relationFilter)
	if err != nil {
		blog.Errorf("get host ids by relation failed, filter: %+v, err: %v, rid: %s", relationFilter, err, kit.Rid)
		return false, nil, err
	}

	return true, hostIDs, nil
}

func (s *Searcher) listHostFromDB(kit *rest.Kit, finalFilter map[string]interface{}, option *metadata.ListHosts) (
	*metadata.ListHostResult, error) {

	total, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(finalFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count hosts failed, filter: %+v, err: %v, rid: %s", finalFilter, err, kit.Rid)
		return nil, err
	}

	searchResult := &metadata.ListHostResult{Count: int(total)}

	limit := uint64(option.Page.Limit)
	start := uint64(option.Page.Start)
	query := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(finalFilter).Limit(limit).
		Start(start).Fields(option.Fields...)
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
func (s *Searcher) listAllBizHostsPage(kit *rest.Kit, fields []string, page metadata.BasePage,
	allBizHostIDs []interface{}) (searchResult *metadata.ListHostResult, err error) {
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
	blog.V(10).Infof("list biz host start: %d, end: %d, rid: %s", start, end, kit.Rid)

	finalFilter := make(map[string]interface{}, 0)
	// 针对获取业务下所有的主机的场景优
	finalFilter[common.BKHostIDField] = map[string]interface{}{
		common.BKDBIN: allBizHostIDs[start:end],
	}

	hosts := make([]metadata.HostMapStr, 0)
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(finalFilter).Fields(fields...).
		Sort(page.Sort).All(kit.Ctx, &hosts); err != nil {
		blog.Errorf("list hosts failed, filter: %+v, err: %v, rid: %s", finalFilter, err, kit.Rid)
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
