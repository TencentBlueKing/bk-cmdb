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
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"

	"gopkg.in/redis.v5"
)

type Searcher struct {
	DbProxy dal.RDB
	Cache   *redis.Client
}

func New(db dal.RDB, cache *redis.Client) Searcher {
	return Searcher{
		DbProxy: db,
		Cache:   cache,
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
	if option.ModuleIDs != nil {
		relationFilter[common.BKModuleIDField] = map[string]interface{}{
			common.BKDBIN: option.ModuleIDs,
		}
	}
	if option.SetIDs != nil {
		relationFilter[common.BKSetIDField] = map[string]interface{}{
			common.BKDBIN: option.SetIDs,
		}
	}
	hostIDFilter := map[string]interface{}{}
	needHostIDFilter := false
	hostIDs := make([]int64, 0)
	if len(relationFilter) > 0 {
		needHostIDFilter = true
		if err := s.DbProxy.Table(common.BKTableNameModuleHostConfig).Distinct(ctx, common.BKHostIDField, relationFilter, &hostIDs); err != nil {
			blog.Errorf("ListHosts failed, db select failed, filter: %+v, err: %+v, rid: %s", relationFilter, err, rid)
			return nil, err
		}

		if len(hostIDs) == 0 {
			return new(metadata.ListHostResult), nil
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
	if propertyFilter != nil {
		filters = append(filters, propertyFilter)
	}
	finalFilter := map[string]interface{}{}
	util.SetQueryOwner(&finalFilter, util.ExtractOwnerFromContext(ctx))
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

	total, err := s.DbProxy.Table(common.BKTableNameBaseHost).Find(finalFilter).Count(ctx)
	if err != nil {
		blog.Errorf("ListHosts failed, db select failed, filter: %+v, err: %+v, rid: %s", finalFilter, err, rid)
		return nil, err
	}
	searchResult = &metadata.ListHostResult{}
	searchResult.Count = int(total)

	limit := uint64(option.Page.Limit)
	start := uint64(option.Page.Start)
	query := s.DbProxy.Table(common.BKTableNameBaseHost).Find(finalFilter).Limit(limit).Start(start).Fields(option.Fields...)
	if len(option.Page.Sort) > 0 {
		query = query.Sort(option.Page.Sort)
	} else {
		query = query.Sort(common.BKHostIDField)
	}

	hosts := make([]map[string]interface{}, 0)
	if err := query.All(ctx, &hosts); err != nil {
		blog.Errorf("ListHosts failed, db select hosts failed, filter: %+v, err: %+v, rid: %s", finalFilter, err, rid)
		return nil, err
	}
	searchResult.Info = hosts
	return searchResult, nil
}

// listAllBizHostsPage 专有流程
func (s *Searcher) listAllBizHostsPage(ctx context.Context, fields []string, page metadata.BasePage,
	allBizHostIDs []int64, rid string) (searchResult *metadata.ListHostResult, err error) {
	cnt := len(allBizHostIDs)
	start := page.Start
	if start > cnt {
		return new(metadata.ListHostResult), nil
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
	// 针对获取业务下所有的主机的场景优化
	finalFilter[common.BKHostIDField] = map[string]interface{}{
		common.BKDBIN: allBizHostIDs[start:end],
	}
	finalFilter = util.SetQueryOwner(finalFilter, util.ExtractOwnerFromContext(ctx))

	hosts := make([]map[string]interface{}, 0)
	if err := s.DbProxy.Table(common.BKTableNameBaseHost).Find(finalFilter).Fields(fields...).
		Sort(page.Sort).All(ctx, &hosts); err != nil {
		blog.Errorf("list hosts failed, db select hosts failed, filter: %+v, err: %+v, rid: %s", finalFilter, err, rid)
		return nil, err
	}
	searchResult = &metadata.ListHostResult{
		Count: cnt,
		Info:  hosts,
	}
	return searchResult, nil
}
