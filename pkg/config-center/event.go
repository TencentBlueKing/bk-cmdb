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

import (
	"encoding/json/v2"
	"fmt"

	"github.com/spf13/cast"
)

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
type Event[T any] struct {
	// Type is the event type.
	Type EventType
	// Pre is the previous config item data.
	Pre T
	// Data is the config item data.
	Data T
}

// EventHandler is the config item change event handler.
type EventHandler[T any] func(*Event[T]) error

// eventHandlers stores config item change event handlers.
var eventHandlers = make(map[string]EventHandler[any])

// registerEventHandler registers event handler.
// Note: must be called before new config reader or writer, otherwise the event handler will not be used.
func registerEventHandler[T any](key string, handler EventHandler[T], conv func(data any) (T, error)) error {
	_, exists := eventHandlers[key]
	if exists {
		return fmt.Errorf("config event handler for key %s already exists", key)
	}

	eventHandlers[key] = func(e *Event[any]) error {
		event := &Event[T]{Type: e.Type}
		if e.Pre != nil {
			pre, err := conv(e.Pre)
			if err != nil {
				return fmt.Errorf("convert pre failed: %w", err)
			}
			event.Pre = pre
		}

		if e.Data != nil {
			data, err := conv(e.Data)
			if err != nil {
				return fmt.Errorf("convert data failed: %w", err)
			}
			event.Data = data
		}

		return handler(event)
	}
	return nil
}

// RegisterBasicEventHandler registers event handler for basic types.
func RegisterBasicEventHandler[T cast.Basic](key string, handler EventHandler[T]) error {
	return registerEventHandler(key, handler, convertBasic[T])
}

// RegisterPtrEventHandler registers event handler for pointer types.
func RegisterPtrEventHandler[T any](key string, handler EventHandler[*T]) error {
	return registerEventHandler(key, handler, convert[T])
}

// convertBasic converts config value to specified basic type value.
func convertBasic[T cast.Basic](data any) (T, error) {
	return cast.ToE[T](data)
}

// convert converts config value to pointer of specified type.
func convert[T any](data any) (*T, error) {
	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %w", err)
	}

	result := new(T)
	if err = json.Unmarshal(marshal, result); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	return result, nil
}
