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

// Package topo defines business topology related common logics
package topo

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/key"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/topo"
	"configcenter/src/storage/driver/redis"
)

// RefreshBizTopo get biz topo info from db and update it to cache
func RefreshBizTopo(kit *rest.Kit, topoKey key.Key, bizID int64, byCache bool) error {
	ctx := context.WithValue(context.Background(), common.ContextRequestIDField, kit.Rid)
	ctx = util.SetDBReadPreference(ctx, common.SecondaryPreferredMode)
	kit = kit.WithCtx(ctx)
	topoType := topoKey.Type()

	blog.Infof("start refreshing biz %d %s topology, by cache: %v, rid: %s", bizID, topoType, byCache, kit.Rid)

	bizTopo, err := topo.GenBizTopo(kit, bizID, topoType, byCache)
	if err != nil {
		blog.Errorf("get biz %d %s topology to refresh failed, err: %v, rid: %s", bizID, topoType, err, kit.Rid)
		return err
	}

	// update it to cache directly.
	_, err = topoKey.UpdateTopology(kit, bizID, bizTopo)
	if err != nil {
		blog.Errorf("refresh biz %d %s topology cache failed, err: %v, rid: %s", bizID, topoType, err, kit.Rid)
		return err
	}

	queue, exists := bizRefreshQueuePool[topoType]
	if exists {
		queue.Remove(kit.TenantID, bizID)
	}

	blog.Infof("refresh biz %d %s topology success, by cache: %v, rid: %s", bizID, topoType, byCache, kit.Rid)
	return nil
}

// TryRefreshBizTopoByCache try refresh biz topo cache by separate node cache, refresh from db data for the first time
func TryRefreshBizTopoByCache(kit *rest.Kit, topoKey key.Key, bizID int64) error {
	ctx := context.WithValue(context.Background(), common.ContextRequestIDField, kit.Rid)
	ctx = util.SetDBReadPreference(ctx, common.SecondaryPreferredMode)
	kit = kit.WithCtx(ctx)

	// check if biz topo cache exists, if not, refresh from db data
	bizTopoKey := topoKey.BizTopoKey(kit.TenantID, bizID)
	existRes, err := redis.Client().Exists(ctx, bizTopoKey).Result()
	if err != nil {
		blog.Errorf("check if biz %d topo cache exists failed, key: %s, err: %v, rid: %s", bizID, bizTopoKey, err,
			kit.Rid)
		return err
	}

	if existRes != 1 {
		return RefreshBizTopo(kit, topoKey, bizID, false)
	}

	// refresh biz topo from cache
	return RefreshBizTopo(kit, topoKey, bizID, true)
}
