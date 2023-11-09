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

package tokenhandler

import "configcenter/src/storage/stream/types"

// WatchTokenTable is the table to store the latest watch token for sync logics
const WatchTokenTable = "cc_SyncWatchToken"

// WatchToken is the watch token data for mongodb watch
type WatchToken struct {
	Coll        string          `json:"_id" bson:"_id"`
	Token       string          `json:"token" bson:"token"`
	StartAtTime types.TimeStamp `json:"start_at_time,omitempty" bson:"start_at_time,omitempty"`
}
