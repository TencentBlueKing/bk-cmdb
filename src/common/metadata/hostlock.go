/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

import (
	"time"

	"configcenter/src/common/mapstr"
)

type HostLockRequest struct {
	IPS     []string `json:"ip_list"`
	CloudID int64    `json:"bk_cloud_id"`
}

type QueryHostLockRequest struct {
	IPS     []string `json:"ip_list"`
	CloudID int64    `json:"bk_cloud_id"`
}

type HostLockResultResponse struct {
	BaseResp `json:",inline"`
	Data     map[string]bool `json:"data"`
}

type HostLockData struct {
	User       string    `json:"bk_user" bson:"bk_user"`
	IP         string    `json:"bk_host_innerip" bson:"bk_host_innerip"`
	CloudID    int64     `json:"bk_cloud_id" bson:"bk_cloud_id"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
	OwnerID    string    `json:"-" bson:"bk_supplier_account"`
}

type HostLockQueryResponse struct {
	BaseResp `json:",inline"`
	Data     struct {
		Info  []HostLockData `json:"info"`
		Count int64          `json:"count"`
	}
}

type HostLockResponse struct {
	BaseResp `json:",inline"`
	Data     mapstr.MapStr `json:"data"`
}
