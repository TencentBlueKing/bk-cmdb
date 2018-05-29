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

	"configcenter/src/framework/core/input"
	"configcenter/src/framework/core/output"
	"context"
)

// Manager contains the basic framework data and the publisher client used to publis events.
type Manager struct {
	cancel      context.CancelFunc
	eventMgr    *eventSubscription
	OutputerMgr output.Manager
	InputerMgr  input.Manager
}

// CreateFrameworkContext create a new framework context instance
func (cli *Manager) CreateFrameworkContext() FrameworkContext {
	return cli
}

// RegisterEvent register cmdb 3.0 event
func (cli *Manager) RegisterEvent(eventType types.EventType, eventFunc types.EventCallbackFunc) types.EventKey {
	return cli.eventMgr.register(eventType, eventFunc)
}

// UnRegisterEvent unregister cmdb 3.0 event
func (cli *Manager) UnRegisterEvent(eventKey types.EventKey) {
	cli.eventMgr.unregister(eventKey)
}

// stop used to stop the business cycles.
func (cli *Manager) stop() error {

	if nil != cli.cancel {
		cli.cancel()
	}

	return cli.InputerMgr.Stop()
}

// Run start the business cycle until the stop method is called.
func (cli *Manager) Run(ctx context.Context, cancel context.CancelFunc) {

	cli.cancel = cancel
	cli.InputerMgr.Run(ctx, cli)
}
