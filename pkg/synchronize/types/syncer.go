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

	"configcenter/src/common/watch"
)

// ListDataOpt is the list data option
type ListDataOpt struct {
	SubRes string
	Start  map[string]int64
	End    map[string]int64
}

// ListDataRes is the list data result
type ListDataRes struct {
	IsAll     bool
	Data      any
	NextStart map[string]int64
}

// CompDataRes is the compare data result
type CompDataRes struct {
	Insert       any
	Update       any
	Delete       any
	RemainingSrc any
}

// FullSyncLockKey is the full sync lock key
const FullSyncLockKey = "cmdb_syncer:full_sync_lock"

// EventInfo is the incremental sync event info
type EventInfo struct {
	EventType watch.EventType
	ResType   ResType
	Oid       string
	SubRes    []string
	Detail    json.RawMessage
}
