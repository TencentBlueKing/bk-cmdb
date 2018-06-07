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

package api

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/types"
)

var (
	events = make([]eventWrapper, 0)
)

type eventWrapper struct {
	key             types.EventKey
	eveType         types.EventType
	eveCallbackFunc types.EventCallbackFunc
}

// registerEvent register the event
func registerEvent(eventType types.EventType, eventFunc types.EventCallbackFunc) types.EventKey {
	key := types.EventKey(common.UUID())
	events = append(events, eventWrapper{
		key:             key,
		eveType:         eventType,
		eveCallbackFunc: eventFunc,
	})
	return key
}

// UnRegisterEvent unregister event
func UnRegisterEvent(eventKey types.EventKey) {
	mgr.UnRegisterEvent(eventKey)
}

// RegisterEventHost register host event
func RegisterEventHost(eventFunc types.EventCallbackFunc) types.EventKey {
	return registerEvent(types.EventHostType, eventFunc)
}

// RegisterEventBusiness register business event
func RegisterEventBusiness(eventFunc types.EventCallbackFunc) types.EventKey {
	return registerEvent(types.EventBusinessType, eventFunc)
}

// RegisterEventModule register business event
func RegisterEventModule(eventFunc types.EventCallbackFunc) types.EventKey {
	return registerEvent(types.EventModuleType, eventFunc)
}

// RegisterEventHostIdentifier register host identifier event
func RegisterEventHostIdentifier(eventFunc types.EventCallbackFunc) types.EventKey {
	return registerEvent(types.EventHostIdentifierType, eventFunc)
}

// RegisterEventSet register host set event
func RegisterEventSet(eventFunc types.EventCallbackFunc) types.EventKey {
	return registerEvent(types.EventSetType, eventFunc)
}

// RegisterEventInst register host inst event
func RegisterEventInst(eventFunc types.EventCallbackFunc) types.EventKey {
	return registerEvent(types.EventInstType, eventFunc)
}

// RegisterEventModuleTransfer register module transfer event
func RegisterEventModuleTransfer(eventFunc types.EventCallbackFunc) types.EventKey {
	return registerEvent(types.EventModuleTransferType, eventFunc)
}
