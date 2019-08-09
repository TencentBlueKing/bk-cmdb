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
	if len(relationFilter) > 0 {
		relations := make([]metadata.ModuleHost, 0)
		if err := s.DbProxy.Table(common.BKTableNameModuleHostConfig).Find(relationFilter).All(ctx, &relations); err != nil {
			blog.Errorf("ListHosts failed, db select failed, filter: %+v, err: %+v, rid: %s", relationFilter, err, rid)
			return nil, err
		}

		hostIDs := make([]int64, 0)
		for _, relation := range relations {
			hostIDs = append(hostIDs, relation.HostID)
		}
		hostIDs = util.IntArrayUnique(hostIDs)

		hostIDFilter = map[string]interface{}{
			common.BKHostIDField: map[string]interface{}{
				common.BKDBIN: hostIDs,
			},
		}
	}
	filters := make([]map[string]interface{}, 0)
	if hostIDFilter != nil && len(hostIDFilter) > 0 {
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

	total, err := s.DbProxy.Table(common.BKTableNameBaseHost).Find(finalFilter).Count(ctx)
	if err != nil {
		blog.Errorf("ListHosts failed, db select failed, filter: %+v, err: %+v, rid: %s", finalFilter, err, rid)
		return nil, err
	}
	searchResult = &metadata.ListHostResult{}
	searchResult.Count = int(total)

	limit := uint64(option.Page.Limit)
	start := uint64(option.Page.Start)
	query := s.DbProxy.Table(common.BKTableNameBaseHost).Find(finalFilter).Limit(limit).Start(start)
	if len(option.Page.Sort) > 0 {
		query = query.Sort(option.Page.Sort)
	}

	hosts := make([]map[string]interface{}, 0)
	if err := query.All(ctx, &hosts); err != nil {
		blog.Errorf("ListHosts failed, db select hosts failed, filter: %+v, err: %+v, rid: %s", finalFilter, err, rid)
		return nil, err
	}
	searchResult.Info = hosts
	return searchResult, nil
}
