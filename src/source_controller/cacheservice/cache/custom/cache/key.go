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
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/source_controller/cacheservice/cache/custom/types"
)

// Namespace is the custom resource cache key namespace
const Namespace = common.BKCacheKeyV3Prefix + "custom_res"

// Key is the custom resource cache key
type Key struct {
	// resType is the custom resource type
	resType types.ResType
	// ttl is the ttl of the custom resource cache
	ttl time.Duration
}

// Type returns the cache key's type
func (k Key) Type() types.ResType {
	return k.resType
}

// TTL the cache key's ttl
func (k Key) TTL() time.Duration {
	return k.ttl
}

// Key is the redis key to store the custom resource cache data
func (k Key) Key(tenantID, key string) string {
	return fmt.Sprintf("%s:%s:%s", Namespace, k.resType, key)
}
