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
	defaultKubeRefreshInterval = 15 * time.Minute
	kubeRefreshIntervalConfig  = "cacheService.kubeTopoSyncIntervalMinutes"
)

func init() {
	TopoKeyMap[types.KubeType] = Key{
		topoType:  types.KubeType,
		namespace: fmt.Sprintf("%stopology:%s", common.BKCacheKeyV3Prefix, types.KubeType),
		ttl:       3 * time.Hour,
		GetRefreshInterval: func() time.Duration {
			if !configcenter.IsExist(kubeRefreshIntervalConfig) {
				return defaultKubeRefreshInterval
			}

			duration, err := configcenter.Int(kubeRefreshIntervalConfig)
			if err != nil {
				blog.Errorf("get kube biz topology cache refresh interval failed, err: %v, use default value", err)
				return defaultKubeRefreshInterval
			}

			interval := time.Duration(duration) * time.Minute

			if interval < defaultKubeRefreshInterval {
				blog.Warnf("kube biz topology cache refresh interval %d is invalid, use default value", duration)
				return defaultKubeRefreshInterval
			}

			return interval
		},
	}
}
