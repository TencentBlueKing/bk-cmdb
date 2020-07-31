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
	EventCacheEventDoneKey           = common.BKCacheKeyV3Prefix + "event:inst_done"

	EventCacheDistIDPrefix      = common.BKCacheKeyV3Prefix + "event:dist_id_"
	EventCacheDistQueuePrefix   = common.BKCacheKeyV3Prefix + "event:dist_queue_"
	EventCacheDistRunningPrefix = common.BKCacheKeyV3Prefix + "event:dist_running_"
	EventCacheDistDonePrefix    = common.BKCacheKeyV3Prefix + "event:dist_done_"

	EventCacheDistCallBackCountPrefix = common.BKCacheKeyV3Prefix + "event:dist_callback_"

	// EventCacheSubscribeFormKey the key prefix in cache
	EventCacheSubscribeFormKey = common.BKCacheKeyV3Prefix + "event:subscribeform:"
	EventCacheProcessChannel   = common.BKCacheKeyV3Prefix + "event_process_channel"

	EventCacheIdentInstPrefix = common.BKCacheKeyV3Prefix + "ident:inst_"
)

// EventSubscriberCacheKey returns EventSubscriberCacheKey
func EventSubscriberCacheKey(ownerID, eventType string) string {
	return EventCacheSubscribeFormKey + ownerID + ":" + eventType
}
