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

// MapStr the common event data definition
type MapStr map[string]interface{}

// EventType  CMDB event definition
type EventType string

// EventKey the event key type
type EventKey string

// EventCallbackFunc the event deal function
type EventCallbackFunc func(evn Event) error

// ContextKey the context key type
type ContextKey string

const (
	// FrameworkKey the framework identifier
	FrameworkKey ContextKey = "cmdb_v3_framework"
)

// Event the cmdb event definition
type Event struct {
}

// Saver the save interface
type Saver interface {
	Save() error
}
