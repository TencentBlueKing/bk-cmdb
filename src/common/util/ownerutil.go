/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package util

import (
	"configcenter/src/common"
)

// SetQueryOwner returns condition that in default ownerID and request ownerID
func SetQueryOwner(condition map[string]interface{}, ownerID string) map[string]interface{} {
	if condition == nil {
		condition = make(map[string]interface{})
	}
	if ownerID == common.BKSuperOwnerID {
		return condition
	}
	if ownerID == common.BKDefaultOwnerID {
		condition[common.BKOwnerIDField] = common.BKDefaultOwnerID
		return condition
	}
	condition[common.BKOwnerIDField] = map[string]interface{}{common.BKDBIN: []string{common.BKDefaultOwnerID, ownerID}}
	return condition
}

// SetModOwner set condition equal owner id, the condition must be a map or struct
func SetModOwner(condition map[string]interface{}, ownerID string) map[string]interface{} {
	if nil == condition {
		condition = make(map[string]interface{})
	}
	if ownerID == common.BKSuperOwnerID {
		return condition
	}
	condition[common.BKOwnerIDField] = ownerID
	return condition
}
