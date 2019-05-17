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
	"strings"
)

// MapStr the common event data definition
type MapStr map[string]interface{}

// EventType  CMDB event definition
type EventType string

// EventKey the event key type
type EventKey string

// Compare compare the event key
func (cli EventKey) Compare(target EventKey) int {
	return strings.Compare(string(cli), string(target))
}

// EventCallbackFunc the event deal function
type EventCallbackFunc func(evn []*Event) error

// ContextKey the context key type
type ContextKey string

const (
	// FrameworkKey the framework identifier
	FrameworkKey ContextKey = "cmdb_v3_framework"

	// EventHostType the host event
	EventHostType EventType = "cmdb_v3_event_host_type"

	// EventBusinessType the business event
	EventBusinessType EventType = "cmdb_v3_event_business_type"

	// EventModuleType the module event
	EventModuleType EventType = "cmdb_v3_event_module_type"

	// EventModuleTransferType the module transfer event
	EventModuleTransferType EventType = "cmdb_v3_event_moduletransfer_type"

	// EventHostIdentifierType the host identifier
	EventHostIdentifierType EventType = "cmdb_v3_event_hostidentifier_type"

	// EventSetType the set event
	EventSetType EventType = "cmdb_v3_event_set_type"

	// EventInstType the custom inst event
	EventInstType EventType = "cmdb_v3_event_inst_type"
)
