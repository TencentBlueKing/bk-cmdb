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

package metadata

import (
	"time"

	"configcenter/src/common/mapstr"
)

// HostLockRequest TODO
type HostLockRequest struct {
	IDS []int64 `json:"id_list"`
}

// QueryHostLockRequest TODO
type QueryHostLockRequest struct {
	IDS []int64 `json:"id_list"`
}

// HostLockResultResponse TODO
type HostLockResultResponse struct {
	BaseResp `json:",inline"`
	Data     map[int64]bool `json:"data"`
}

// HostLockData TODO
type HostLockData struct {
	User       string    `json:"bk_user" bson:"bk_user"`
	ID         int64     `json:"bk_host_id" bson:"bk_host_id"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
	OwnerID    string    `json:"-" bson:"bk_supplier_account"`
}

// HostLockQueryResponse TODO
type HostLockQueryResponse struct {
	BaseResp `json:",inline"`
	Data     struct {
		Info  []HostLockData `json:"info"`
		Count int64          `json:"count"`
	}
}

// HostLockResponse TODO
type HostLockResponse struct {
	BaseResp `json:",inline"`
	Data     mapstr.MapStr `json:"data"`
}
