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

// Package util defines cmdb data syncer utility
package util

import (
	"math/rand"
	"time"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/source_controller/transfer-service/app/options"
)

// RetryWrapper retry the handler for maxRetry times if there is an error and need retry
func RetryWrapper(maxRetry int, handler func() (bool, error)) {
	retry := 0
	for {
		needRetry, err := handler()
		if err == nil {
			return
		}

		if !needRetry {
			return
		}

		retry++
		if retry > maxRetry {
			return
		}

		time.Sleep(time.Duration(rand.Intn(10)+1) * time.Second * time.Duration(retry))
	}
}

// MatchIDRule checks if any resource ids match id rules
func MatchIDRule(idRuleMap map[types.ResType]map[string][]options.IDRuleInfo, idMap map[types.ResType][]int64,
	env string) bool {

	for resType, ids := range idMap {
		ruleMap := idRuleMap[resType]
		if idsMatchIDRule(ruleMap[env], ids) {
			return true
		}
	}

	return false
}

// idsMatchIDRule checks if the ids match the id rules
func idsMatchIDRule(rules []options.IDRuleInfo, ids []int64) bool {
	for _, id := range ids {
		for _, rule := range rules {
			if rule.StartID > id {
				break
			}

			if (id <= rule.EndID || rule.EndID == types.InfiniteEndID) && (id-rule.StartID)%rule.Step == 0 {
				return true
			}
		}
	}

	return false
}
