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

package key

import (
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
)

var (
	defaultBriefRefreshInterval = 15 * time.Minute
	briefRefreshIntervalConfig  = "cacheService.briefTopologySyncIntervalMinutes"
)

func init() {
	TopoKeyMap[types.BriefType] = Key{
		topoType:  types.BriefType,
		namespace: fmt.Sprintf("%stopology:%s", common.BKCacheKeyV3Prefix, types.BriefType),
		ttl:       24 * time.Hour,
		GetRefreshInterval: func() time.Duration {
			if !configcenter.IsExist(briefRefreshIntervalConfig) {
				return defaultBriefRefreshInterval
			}

			duration, err := configcenter.Int(briefRefreshIntervalConfig)
			if err != nil {
				blog.Errorf("get brief biz topology cache refresh interval failed, err: %v, use default value", err)
				return defaultBriefRefreshInterval
			}

			if duration < 2 {
				blog.Warnf("brief biz topology cache refresh interval %d is invalid, use default value", duration)
				return defaultBriefRefreshInterval
			}

			return time.Duration(duration) * time.Minute
		},
	}
}
