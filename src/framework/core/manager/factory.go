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

package manager

import (
	"configcenter/src/framework/core/types"
	"net/http"
)

// New return a new  Manager instance
func New() *Manager {
	evn := &eventSubscription{
		datas:             make(chan types.MapStr, 4096),
		registers:         make(map[types.EventType][]*eventRegister),
		hostMgr:           &eventHost{},
		businessMgr:       &eventBusiness{},
		setMgr:            &eventSet{},
		moduleMgr:         &eventModule{},
		moduleTransferMgr: &eventModuleTransfer{},
		hostIdentifierMgr: &eventHostIdentifier{},
	}
	return &Manager{
		eventMgr: evn,
		ms: []Action{
			Action{
				Method:      http.MethodPost,
				Path:        "/api/v1/event/puts",
				HandlerFunc: evn.puts,
			},
		},
	}
}

// Delete delete the framework instance
func Delete(mgr *Manager) error {

	if nil != mgr {
		return mgr.stop()
	}

	return nil
}
