/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package general

import (
	fullsynccond "configcenter/pkg/cache/full-sync-cond"
	"configcenter/pkg/cache/general"
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	fullsynccondcli "configcenter/src/source_controller/cacheservice/cache/general/full-sync-cond"
	"configcenter/src/source_controller/cacheservice/cache/general/types"
)

// FullSyncCond returns full sync cond client
func (c *Cache) FullSyncCond() *fullsynccondcli.FullSyncCond {
	return c.fullSyncCond
}

// ListDetailByIDs list general resource detail cache by ids
// NOTE: since event flow and cache are reused, this method may return deleted data
// because event ttl is long and event detail cache will not be deleted
func (c *Cache) ListDetailByIDs(kit *rest.Kit, opt *general.ListDetailByIDsOpt) ([]string, error) {
	cache, exists := c.cacheSet[opt.Resource]
	if !exists {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, general.ResourceField)
	}

	idKeys := make([]string, len(opt.IDs))
	for i, id := range opt.IDs {
		idKeys[i], _ = cache.Key().IDKey(id, "")
	}

	listOpt := &types.ListDetailByIDsOpt{
		SubRes: opt.SubResource,
		IDKeys: idKeys,
		Fields: opt.Fields,
	}

	return cache.ListDetailByIDs(kit, listOpt)
}

// ListDetailByUniqueKey list general resource detail cache by unique keys
func (c *Cache) ListDetailByUniqueKey(kit *rest.Kit, opt *general.ListDetailByUniqueKeyOpt) ([]string, error) {
	cache, exists := c.cacheSet[opt.Resource]
	if !exists {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, general.ResourceField)
	}

	listOpt := &types.ListDetailByUniqueKeyOpt{
		SubRes: opt.SubResource,
		Type:   opt.Type,
		Keys:   opt.Keys,
		Fields: opt.Fields,
	}

	return cache.ListDetailByUniqueKey(kit, listOpt)
}

// ListCacheByFullSyncCond list general resource detail cache by full sync cond
func (c *Cache) ListCacheByFullSyncCond(kit *rest.Kit, opt *fullsynccond.ListCacheByFullSyncCondOpt,
	cond *fullsynccond.FullSyncCond) ([]string, error) {

	cache, exists := c.cacheSet[cond.Resource]
	if !exists {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, general.ResourceField)
	}

	idListKey, err := cache.GenFullSyncCondIDListKey(cond)
	if err != nil {
		return nil, err
	}

	listOpt := &types.ListDetailOpt{
		Fields: opt.Fields,
		IDListFilter: &types.IDListFilterOpt{
			IDListKey: idListKey,
			BasicFilter: &types.BasicFilter{
				SubRes: cond.SubResource,
			},
			IsAll: cond.IsAll,
			Cond:  cond.Condition,
		},
		Page: &general.PagingOption{
			StartID: opt.Cursor,
			Limit:   opt.Limit,
		},
	}
	return cache.ListDetail(kit, listOpt)
}

// ListData list general resource count or detail from cache for system request
func (c *Cache) ListData(kit *rest.Kit, opt *general.ListDetailOpt) (int64, []string, error) {
	cache, exists := c.cacheSet[opt.Resource]
	if !exists {
		return 0, nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, general.ResourceField)
	}

	idListKey := cache.Key().IDListKey(kit.TenantID)
	if opt.SubResource != "" {
		idListKey = cache.Key().IDListKey(kit.TenantID, opt.SubResource)
	}

	listOpt := &types.ListDetailOpt{
		Fields: opt.Fields,
		IDListFilter: &types.IDListFilterOpt{
			IDListKey: idListKey,
			BasicFilter: &types.BasicFilter{
				SubRes: opt.SubResource,
			},
			IsAll: true,
		},
		Page: opt.Page,
	}

	if opt.Page.EnableCount {
		cnt, err := cache.CountData(kit, listOpt)
		if err != nil {
			return 0, nil, err
		}
		return cnt, make([]string, 0), nil
	}

	details, err := cache.ListDetail(kit, listOpt)
	if err != nil {
		return 0, nil, err
	}
	return 0, details, nil
}

// RefreshIDList refresh general resource id list cache
func (c *Cache) RefreshIDList(kit *rest.Kit, opt *general.RefreshIDListOpt,
	cond *fullsynccond.FullSyncCond) error {
	cache, exists := c.cacheSet[opt.Resource]
	if !exists {
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, general.ResourceField)
	}

	return cache.RefreshIDList(kit, opt, cond)
}

// RefreshDetailByIDs refresh general resource detail cache by ids
func (c *Cache) RefreshDetailByIDs(kit *rest.Kit, opt *general.RefreshDetailByIDsOpt) error {
	cache, exists := c.cacheSet[opt.Resource]
	if !exists {
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, general.ResourceField)
	}

	idKeys := make([]string, len(opt.IDs))
	for i, id := range opt.IDs {
		idKeys[i], _ = cache.Key().IDKey(id, "")
	}

	refreshOpt := &types.RefreshDetailByIDsOpt{
		SubResource: opt.SubResource,
		IDKeys:      idKeys,
	}

	return cache.RefreshDetailByIDs(kit, refreshOpt)
}
