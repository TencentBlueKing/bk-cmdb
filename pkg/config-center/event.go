/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

package cc

import "fmt"

// EventType is the config item change event type.
type EventType string

const (
	// UpsertEvent is the create and update event type.
	UpsertEvent EventType = "upsert"
	// DeleteEvent is the delete event type.
	DeleteEvent EventType = "delete"
)

// DiscoveryEvent is the config change event for discovery.
type DiscoveryEvent struct {
	// Type is the event type.
	Type EventType
	// Key is the config key.
	Key string
	// Data is the config data.
	Data []byte
}

// Event is the config item change event.
type Event struct {
	// Type is the event type.
	Type EventType
	// Pre is the previous config item data.
	Pre any
	// Data is the config item data.
	Data any
}

// EventHandler is the config item change event handler.
type EventHandler func(*Event) error

// eventHandlers stores config item change event handlers.
var eventHandlers = make(map[string]EventHandler)

// RegisterEventHandler registers event handler
// Note: must be called before new config reader or writer, otherwise the event handler will not be used.
func RegisterEventHandler(key string, handler EventHandler) error {
	_, exists := eventHandlers[key]
	if exists {
		return fmt.Errorf("config event handler for key %s already exists", key)
	}
	eventHandlers[key] = handler
	return nil
}
