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

package types

import (
	"configcenter/src/common"
)

// Event Cache Keys
const (

	// EventCacheEventIDKey the event instance id key in cache
	EventCacheEventIDKey             = common.BKCacheKeyV3Prefix + "event:inst_id"
	EventCacheEventQueueKey          = common.BKCacheKeyV3Prefix + "event:inst_queue"
	EventCacheEventQueueDuplicateKey = common.BKCacheKeyV3Prefix + "event:inst_queue_duplicate"
	EventCacheEventPendingKey        = common.BKCacheKeyV3Prefix + "event:inst_pending"
	EventCacheEventRunningPrefix     = common.BKCacheKeyV3Prefix + "event:inst_running_"
	EventCacheEventTimeoutKey        = common.BKCacheKeyV3Prefix + "event:inst_timeout"
	EventCacheEventDoneKey           = common.BKCacheKeyV3Prefix + "event:inst_done"

	EventCacheDistIDPrefix      = common.BKCacheKeyV3Prefix + "event:dist_id_"
	EventCacheDistQueuePrefix   = common.BKCacheKeyV3Prefix + "event:dist_queue_"
	EventCacheDistPendingPrefix = common.BKCacheKeyV3Prefix + "event:dist_pending_"
	EventCacheDistRunningPrefix = common.BKCacheKeyV3Prefix + "event:dist_running_"
	EventCacheDistTimeoutPrefix = common.BKCacheKeyV3Prefix + "event:dist_timeout_"
	EventCacheDistDonePrefix    = common.BKCacheKeyV3Prefix + "event:dist_done_"

	EventCacheDistCallBackCountPrefix = common.BKCacheKeyV3Prefix + "event:dist_callback_"

	// EventCacheSubscribeformKey the key prefix in cache
	EventCacheSubscribeformKey = common.BKCacheKeyV3Prefix + "event:subscribeform:"
	EventCacheSubscribesKey    = common.BKCacheKeyV3Prefix + "event:subscribers"
	EventCacheProcessChannel   = common.BKCacheKeyV3Prefix + "event_process_channel"

	EventCacheIdentInstPrefix = common.BKCacheKeyV3Prefix + "ident:inst_"
)

// EventSubscriberCacheKey returns EventSubscriberCacheKey
func EventSubscriberCacheKey(ownerID, eventtype string) string {
	return EventCacheSubscribeformKey + ownerID + ":" + eventtype
}
