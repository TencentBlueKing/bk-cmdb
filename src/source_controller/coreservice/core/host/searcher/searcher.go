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

func (s *Searcher) ListHostByTopoNode(ctx context.Context, option metadata.ListHostByTopoNodeOption) (searchResult *metadata.ListHostResult, err error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	if key, err := option.Validate(); err != nil {
		blog.Errorf("ListHostByTopoNode failed, invalid parameters, key: %s, err: %s, rid: %s", key, err.Error(), rid)
		return nil, err
	}

	filter := map[string]interface{}{
		common.BKAppIDField: option.BizID,
	}
	if option.ModuleIDs != nil {
		filter[common.BKModuleIDField] = map[string]interface{}{
			common.BKDBIN: option.ModuleIDs,
		}
	}
	if option.SetIDs != nil {
		filter[common.BKSetIDField] = map[string]interface{}{
			common.BKDBIN: option.SetIDs,
		}
	}
	total, err := s.DbProxy.Table(common.BKTableNameModuleHostConfig).Find(filter).Count(ctx)
	if err != nil {
		blog.Errorf("ListHostByTopoNode failed, db select failed, filter: %+v, err: %+v, rid: %s", filter, err, rid)
		return nil, err
	}

	searchResult = &metadata.ListHostResult{}
	searchResult.Count = int(total)

	relations := make([]metadata.ModuleHost, 0)
	if err := s.DbProxy.Table(common.BKTableNameModuleHostConfig).Find(filter).All(ctx, &relations); err != nil {
		blog.Errorf("ListHostByTopoNode failed, db select failed, filter: %+v, err: %+v, rid: %s", filter, err, rid)
		return nil, err
	}

	hostIDs := make([]int64, 0)
	for _, relation := range relations {
		hostIDs = append(hostIDs, relation.HostID)
	}
	hostIDs = util.IntArrayUnique(hostIDs)

	hostFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDs,
		},
	}
	limit := uint64(option.Page.Limit)
	start := uint64(option.Page.Start)
	query := s.DbProxy.Table(common.BKTableNameBaseHost).Find(hostFilter).Limit(limit).Start(start)
	if len(option.Page.Sort) > 0 {
		query = query.Sort(option.Page.Sort)
	}

	hosts := make([]map[string]interface{}, 0)
	if err := query.All(ctx, &hosts); err != nil {
		blog.Errorf("ListHostByTopoNode failed, db select hosts failed, filter: %+v, err: %+v, rid: %s", hostFilter, err, rid)
		return nil, err
	}
	searchResult.Info = hosts
	return searchResult, nil
}
