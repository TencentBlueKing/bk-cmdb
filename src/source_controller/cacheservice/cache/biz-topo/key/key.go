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

// Package key defines the topology redis cache key logics
package key

import (
	"fmt"
	"time"

	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
	"configcenter/src/storage/driver/redis"
)

// TopoKeyMap is the topology type to cache key map
var TopoKeyMap = make(map[types.TopoType]Key)

// Key is the biz topology cache key
type Key struct {
	// topoType is the topology type
	topoType types.TopoType
	// namespace is the cache key's namespace
	namespace string
	// ttl is the ttl of the topology cache
	ttl time.Duration
	// GetRefreshInterval returns the topology cache force refresh time interval
	GetRefreshInterval func() time.Duration
}

// Type returns the cache key's type
func (k Key) Type() types.TopoType {
	return k.topoType
}

// Namespace returns the cache key's namespace
func (k Key) Namespace() string {
	return k.namespace
}

// TTL the cache key's ttl
func (k Key) TTL() time.Duration {
	return k.ttl
}

// BizTopoKey is the redis key to store the biz topology tree
func (k Key) BizTopoKey(tenantID string, biz int64) string {
	return fmt.Sprintf("%s:%s:%d", k.namespace, tenantID, biz)
}

// UpdateTopology update biz topology cache
func (k Key) UpdateTopology(kit *rest.Kit, bizID int64, topo any) (*string, error) {
	js, err := json.Marshal(topo)
	if err != nil {
		return nil, fmt.Errorf("marshal %s biz %d topology failed, err: %v", k.topoType, bizID, err)
	}
	value := string(js)

	return &value, redis.Client().Set(kit.Ctx, k.BizTopoKey(kit.TenantID, bizID), value, k.ttl).Err()
}

// GetTopology get biz Topology from cache
func (k Key) GetTopology(kit *rest.Kit, biz int64) (*string, error) {
	dat, err := redis.Client().Get(kit.Ctx, k.BizTopoKey(kit.TenantID, biz)).Result()
	if err != nil {
		if redis.IsNilErr(err) {
			empty := ""
			return &empty, nil
		}

		return nil, fmt.Errorf("get cache from redis failed, err: %v", err)
	}

	return &dat, nil
}
