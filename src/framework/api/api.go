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
	"configcenter/src/framework/core/httpserver"
	"configcenter/src/framework/core/input"
	"configcenter/src/framework/core/manager"
	"configcenter/src/framework/core/output"
	"context"
)

//  mgr the global variable for the manager
var mgr *manager.Manager

// Init init the framework
func Init() {

	ctx, cancel := context.WithCancel(context.Background())

	// create Framework
	mgr = manager.New()

	/** initialize the default configuration */

	// set outputer manager
	mgr.OutputerMgr = output.New()

	// set inputer manager
	mgr.InputerMgr = input.New()

	// init inputer
	for _, inputer := range inputers {
		mgr.InputerMgr.AddInputer(inputer)
	}

	// register events
	for _, eve := range events {
		mgr.RegisterEvent(eve.key, eve.eveType, eve.eveCallbackFunc)
	}

	/** start the main business loop */
	common.GoRun(func() {
		mgr.Run(ctx, cancel)
	}, nil)

}

// Actions return the framework actions
func Actions() []httpserver.Action {
	return mgr.Actions()
}

// UnInit destory the framework
func UnInit() error {
	defer func() {
		mgr = nil
	}()
	if nil == mgr {
		return nil
	}
	return manager.Delete(mgr)
}
