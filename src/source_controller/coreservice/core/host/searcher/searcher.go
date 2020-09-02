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
	"configcenter/src/source_controller/coreservice/cache"
	"configcenter/src/storage/dal"

	"gopkg.in/redis.v5"
)

type Searcher struct {
	DbProxy  dal.RDB
	Cache    *redis.Client
	CacheSet *cache.ClientSet
}

func New(db dal.RDB, cache *redis.Client, CacheSet *cache.ClientSet) Searcher {
	return Searcher{
		DbProxy:  db,
		Cache:    cache,
		CacheSet: CacheSet,
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
	if len(relationFilter) > 0 {

		hostIDs, err := s.DbProxy.Table(common.BKTableNameModuleHostConfig).Distinct(ctx, common.BKHostIDField, relationFilter)
		if err != nil {
			blog.Errorf("ListHosts failed, db select failed, filter: %+v, err: %+v, rid: %s", relationFilter, err, rid)
			return nil, err
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

	if len(filters) == 0 {

		// return info use cache
		// fix: has question when multi-supplier
		sort := strings.TrimLeft(option.Page.Sort, "+-")
		if len(option.Page.Sort) == 0 || sort == common.BKHostIDField || strings.Contains(sort, ",") == false && strings.HasPrefix(sort, common.BKHostIDField+":") {
			opt := &metadata.ListHostWithPage{
				Fields: option.Fields,
				Page:   option.Page,
			}
			total, infos, err := s.CacheSet.Host.ListHostsWithPage(ctx, opt)
			if err != nil {
				blog.ErrorJSON("list host from redis error, filter: %s, err: %s, rid: %s", opt, err.Error(), rid)
				// HOOK: 缓存出现错误从db中获取数据
			} else {
				searchResult = &metadata.ListHostResult{}
				searchResult.Count = int(total)

				searchResult.Info = make([]map[string]interface{}, 0)
				if err := json.UnmarshalArray(infos, &searchResult.Info); err != nil {
					blog.ErrorJSON("list host from redis error, filter: %s, item host: %s, err: %s,rid: %s", opt, infos, err, rid)
					// TODO： use cc error. keep the same as before code
					return nil, err
				}

				return searchResult, nil
			}

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
