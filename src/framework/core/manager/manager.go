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
	"configcenter/src/framework/core/input"
	"configcenter/src/framework/core/output"
	"context"
)

// Manager contains the basic framework data and the publisher client used to publis events.
type Manager struct {
	OutputerMgr output.Manager
	InputerMgr  input.Manager
}

// stop used to stop the business cycles.
func (cli *Manager) stop() error {

	return cli.InputerMgr.Stop()
}

// Run start the business cycle until the stop method is called.
func (cli *Manager) Run(ctx context.Context, cancel context.CancelFunc) {

	cli.InputerMgr.Run(ctx, cancel)
}
