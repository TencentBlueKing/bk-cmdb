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
	"fmt"

	"configcenter/src/common"
)

var (
	// NilStr is special NIL string.
	NilStr = "nil"
)

const (
	// MetricsNamespacePrefix is prefix of metrics namespace.
	MetricsNamespacePrefix = "cmdb_eventserver"

	// EventCacheEventIDKey the event instance id key in cache
	EventCacheEventIDKey = common.BKCacheKeyV3Prefix + "event:inst_id"

	// EventCacheDistIDPrefix is prefix of event dist id key in cache.
	EventCacheDistIDPrefix = common.BKCacheKeyV3Prefix + "event:dist_id_"

	// EventCacheIdentInstPrefix is prefix of host identifier instance key in cache.
	EventCacheIdentInstPrefix = common.BKCacheKeyV3Prefix + "ident:inst_"

	// EventCacheEventQueueKey is main event queue key in cache.
	EventCacheEventQueueKey = common.BKCacheKeyV3Prefix + "event:queue"

	// EventCacheEventQueueDuplicateKey is duplicate event queue key in cache.
	EventCacheEventQueueDuplicateKey = common.BKCacheKeyV3Prefix + "event:duplicate_queue"

	// EventCacheSubscriberEventQueueKeyPrefix is prefix of subscriber event queue key in cache.
	EventCacheSubscriberEventQueueKeyPrefix = common.BKCacheKeyV3Prefix + "event:subscriber_queue_"

	// EventCacheSubscriberCursorPrefixis prefix for subscriber on target resource event type.
	// e.g: cc:v3:event:subscriber_cursor_hostcreate:1 -> MarshalChainNodeStr
	EventCacheSubscriberCursorPrefix = common.BKCacheKeyV3Prefix + "event:subscriber_cursor"

	// EventCacheDistCallBackCountPrefix is prefix of event callback stats key in cache.
	EventCacheDistCallBackCountPrefix = common.BKCacheKeyV3Prefix + "event:dist_callback_"
)

// EventCacheSubscriberCursorKey returns redis key for subscriber cursor cache.
func EventCacheSubscriberCursorKey(eventType string, subid int64) string {
	return fmt.Sprintf("%s_%s_%d", EventCacheSubscriberCursorPrefix, eventType, subid)
}
