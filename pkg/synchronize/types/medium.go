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

package types

import (
	"encoding/json"
)

// PushSyncDataOpt is the push full sync data option
type PushSyncDataOpt struct {
	ResType     ResType `json:"resource_type"`
	SubRes      string  `json:"sub_resource"`
	IsIncrement bool    `json:"is_increment"`
	Data        any     `json:"data"`
}

// PullSyncDataOpt is the pull sync data option
type PullSyncDataOpt struct {
	ResType     ResType `json:"resource_type"`
	SubRes      string  `json:"sub_resource"`
	IsIncrement bool    `json:"is_increment"`
	Ack         bool    `json:"ack"`
}

// TransferMediumResp is the transfer medium response
type TransferMediumResp[T any] struct {
	Result  bool   `json:"result"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// PullSyncDataRes is the pull sync data result data
type PullSyncDataRes struct {
	Total int64           `json:"total"`
	Info  json.RawMessage `json:"info"`
}

// FullSyncTransData is the full sync transfer data
type FullSyncTransData struct {
	Name  string           `json:"name"`
	Start map[string]int64 `json:"start"`
	End   map[string]int64 `json:"end"`
	Data  any              `json:"data"`
}

// IncrSyncTransData is the incremental sync transfer data
type IncrSyncTransData struct {
	Name       string                       `json:"name"`
	UpsertInfo map[string][]json.RawMessage `json:"upsert_info"`
	DeleteInfo map[string][]json.RawMessage `json:"delete_info"`
}
