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

package cache

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/cache/general/types"
	"configcenter/src/storage/driver/redis"
)

// convertToAnyArr convert array into array of interface{} type
func convertToAnyArr[T any](arr []T) []any {
	result := make([]any, len(arr))
	for i, data := range arr {
		result[i] = data
	}

	return result
}

// parseWatchChainNode parse watch chain node into basic info
func parseWatchChainNode(node *watch.ChainNode) (*basicInfo, error) {
	if node == nil {
		return nil, fmt.Errorf("chain node is nil")
	}

	return &basicInfo{
		id:       node.InstanceID,
		oid:      node.Oid,
		subRes:   node.SubResource,
		supplier: node.SupplierAccount,
	}, nil
}

// isIDListExists check if id list exists
func isIDListExists(ctx context.Context, key string, rid string) (bool, error) {
	existRes, err := redis.Client().Exists(ctx, key).Result()
	if err != nil {
		blog.Errorf("check if id list %s exists failed, err: %v, opt: %+v, rid: %s", key, err, rid)
		return false, err
	}

	if existRes != 1 {
		blog.V(4).Infof("id list %s key not exists. rid: %s", key, rid)
		return false, nil
	}

	return true, nil
}

// setQueryOwner set supplier account for general resource query, system resource do not need to set supplier account
func setQueryOwner(cond map[string]interface{}, opt *types.BasicFilter) map[string]interface{} {
	if !opt.IsSystem {
		return util.SetQueryOwner(cond, opt.SupplierAccount)
	}
	return cond
}

// validateIDList check if id list is valid, returns the id list's ttl
func (c *Cache) validateIDList(opt *types.IDListFilterOpt) (time.Duration, error) {
	fullSyncCondInfo, exists := c.fullSyncCondMap.Get(opt.IDListKey)
	if exists {
		return fullSyncCondInfo.Interval, nil
	}

	if !c.needCacheAll {
		return 0, fmt.Errorf("resource type %s do not support cache all id list", c.key.Resource())
	}

	return c.expireSeconds, nil
}
