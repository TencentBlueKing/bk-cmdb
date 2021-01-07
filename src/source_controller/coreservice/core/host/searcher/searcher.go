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

package searcher

import (
	"context"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/host"
	"configcenter/src/storage/driver/mongodb"
)

type Searcher struct {
	CacheHost *host.Client
}

func New() Searcher {
	return Searcher{
		CacheHost: host.NewClient(),
	}
}

// ListHosts search host with topo node info and host property
func (s *Searcher) ListHosts(ctx context.Context, option metadata.ListHosts) (searchResult *metadata.ListHostResult, err error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	if key, err := option.Validate(); err != nil {
		blog.Errorf("ListHosts failed, invalid parameters, key: %s, err: %s, rid: %s", key, err.Error(), rid)
		return nil, err
	}

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
	hostIDFilter := map[string]interface{}{}
	needHostIDFilter := false
	var hostIDs []interface{}
	if len(relationFilter) > 0 {
		needHostIDFilter = true
		var err error
		hostIDs, err = mongodb.Client().Table(common.BKTableNameModuleHostConfig).Distinct(ctx, common.BKHostIDField, relationFilter)
		if err != nil {
			blog.Errorf("ListHosts failed, db select failed, filter: %+v, err: %+v, rid: %s", relationFilter, err, rid)
			return nil, err
		}
		if len(hostIDs) == 0 {
			return &metadata.ListHostResult{
				Count: 0,
				Info:  []map[string]interface{}{},
			}, nil
		}

		hostIDFilter = map[string]interface{}{
			common.BKHostIDField: map[string]interface{}{
				common.BKDBIN: hostIDs,
			},
		}
	}
	filters := make([]map[string]interface{}, 0)
	if len(hostIDFilter) > 0 {
		filters = append(filters, hostIDFilter)
	}

	propertyFilter, err := option.GetHostPropertyFilter(ctx)
	if err != nil {
		blog.Errorf("ListHosts failed, db select failed, filter: %+v, err: %+v, rid: %s", hostIDFilter, err, rid)
		return nil, err
	}
	if len(propertyFilter) > 0 {
		filters = append(filters, propertyFilter)
	}

	finalFilter := map[string]interface{}{}
	finalFilter = util.SetQueryOwner(finalFilter, util.ExtractOwnerFromContext(ctx))
	if len(filters) > 0 {
		finalFilter[common.BKDBAND] = filters
	}

	if needHostIDFilter && len(filters) == 1 && option.BizID != 0 {
		sort := strings.TrimLeft(option.Page.Sort, "+-")
		if len(option.Page.Sort) == 0 || sort == common.BKHostIDField || strings.Contains(sort, ",") == false &&
			strings.HasPrefix(sort, common.BKHostIDField+":") {
			searchResult, err = s.listAllBizHostsPage(ctx, option.Fields, option.Page, hostIDs, rid)
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
			var skip bool
			searchResult, skip, err = s.ListHostsWithCache(ctx, option.Fields, option.Page, rid)
			if err != nil {
				return nil, err
			}
			if skip == false {
				return searchResult, nil
			}
		}
	}

	total, err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(finalFilter).Count(ctx)
	if err != nil {
		blog.Errorf("ListHosts failed, db select failed, filter: %+v, err: %+v, rid: %s", finalFilter, err, rid)
		return nil, err
	}
	searchResult = &metadata.ListHostResult{}
	searchResult.Count = int(total)

	limit := uint64(option.Page.Limit)
	start := uint64(option.Page.Start)
	query := mongodb.Client().Table(common.BKTableNameBaseHost).Find(finalFilter).Limit(limit).Start(start).Fields(option.Fields...)
	if len(option.Page.Sort) > 0 {
		query = query.Sort(option.Page.Sort)
	} else {
		query = query.Sort(common.BKHostIDField)
	}

	hosts := make([]metadata.HostMapStr, 0)
	if err := query.All(ctx, &hosts); err != nil {
		blog.Errorf("ListHosts failed, db select hosts failed, filter: %+v, err: %+v, rid: %s", finalFilter, err, rid)
		return nil, err
	}
	searchResult.Info = make([]map[string]interface{}, len(hosts))
	for index, host := range hosts {
		searchResult.Info[index] = host
	}
	return searchResult, nil
}

func (s *Searcher) ListHostsWithCache(ctx context.Context, fields []string, page metadata.BasePage, rid string) (
	searchResult *metadata.ListHostResult, skip bool, err error) {
	opt := &metadata.ListHostWithPage{
		Fields: fields,
		Page:   page,
	}
	total, infos, err := s.CacheHost.ListHostsWithPage(ctx, opt)
	if err != nil {
		blog.ErrorJSON("list host from redis error, filter: %s, err: %s, rid: %s", opt, err.Error(), rid)
		// HOOK: 缓存出现错误从db中获取数据
		return nil, true, nil
	} else {
		searchResult = &metadata.ListHostResult{}
		searchResult.Count = int(total)

		searchResult.Info = make([]map[string]interface{}, 0)
		if err := json.UnmarshalArray(infos, &searchResult.Info); err != nil {
			blog.ErrorJSON("list host from redis error, filter: %s, item host: %s, err: %s,rid: %s", opt, infos, err, rid)
			// TODO： use cc error. keep the same as before code
			return nil, false, err
		}

		return searchResult, false, nil
	}
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
