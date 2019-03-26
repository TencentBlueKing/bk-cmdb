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

package input

import (
    "sync"
    "time"

	"configcenter/src/framework/core/output"
)

// wrapInputer the Inputer wrapper
type wrapInputer struct {
    sync.Mutex
	isTiming  bool
	frequency time.Duration
	kind      InputerType
	status    InputerStatus
	inputer   Inputer
	putter    output.Puter
	exception ExceptionFunc
}

func (cli *wrapInputer) SetStatus(status InputerStatus) {
    cli.Lock()
    defer cli.Unlock()
	cli.status = status
}

func (cli *wrapInputer) GetStatus() InputerStatus {
    cli.Lock()
    defer cli.Unlock()
	return cli.status
}

func (cli *wrapInputer) GetFrequency() time.Duration {
    cli.Lock()
    defer cli.Unlock()
	return cli.frequency
}

func (cli *wrapInputer) Name() string {
    cli.Lock()
    defer cli.Unlock()
	return cli.inputer.Name()
}

func (cli *wrapInputer) Run(ctx InputerContext) *InputerResult {
	return cli.inputer.Run(ctx)
}

func (cli *wrapInputer) Stop() {
	cli.inputer.Stop()
}
