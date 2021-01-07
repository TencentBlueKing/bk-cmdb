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
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/log"
	"context"
	"sync"
	"time"
)

// manager implements the Manager interface
type manager struct {
	ctx         InputerContext
	inputerLock sync.RWMutex
	inputers    MapInputer
	inputerChan chan *wrapInputer
}

func (cli *manager) AddInputer(params InputerParams) InputerKey {

	key := makeInputerKey()

	target := &wrapInputer{
		frequency: params.Frequency,
		isTiming:  params.IsTiming,
		inputer:   params.Target,
		status:    NormalStatus,
		kind:      params.Kind,
		putter:    params.Putter,
	}

	cli.inputerLock.Lock()
	cli.inputers[key] = target
	cli.inputerLock.Unlock()

	select {
	case cli.inputerChan <- target:
	default:
		log.Fatal("failed to puth the inputer")
	}

	return key
}

// RemoveInputer remove the Inputer by a InputerKey
func (cli *manager) RemoveInputer(key InputerKey) {

	cli.inputerLock.Lock()
	defer cli.inputerLock.Unlock()

	deleteInputer(cli.inputers, key)
}

// Stop used to stop the business cycles.
func (cli *manager) Stop() error {

	// stop the all Inputers
	cli.inputerLock.Lock()
	for _, inputer := range cli.inputers {
		inputer.Stop()
	}
	cli.inputerLock.Unlock()

	return nil
}

// Run start the business cycle until the stop method is called.
func (cli *manager) Run(ctx context.Context, inputerCtx InputerContext) {

	// catch the framework context
	cli.ctx = inputerCtx

	// check the stat of the Inputer regularly, and start it if there is any new
	for {
		select {
		case <-ctx.Done():
			log.Info("will exit from inputer main business cycle")
			goto end

		case target := <-cli.inputerChan:
			common.GoRun(func() {
				cli.executeInputer(ctx, target)
			}, func() {
				target.SetStatus(ExceptionExitStatus)
			})

		case <-time.After(time.Second * 10):

			cli.inputerLock.RLock()

			// scan the all Inputers and restart the stoped Inputer
			for _, inputer := range cli.inputers {
				switch inputer.GetStatus() {
				case NormalStatus:
					common.GoRun(func() {
						cli.executeInputer(ctx, inputer)
					}, func() {
						inputer.SetStatus(ExceptionExitStatus)
					})

				case WaitingToRunStatus:
					common.GoRun(func() {
						cli.executeInputer(ctx, inputer)
					}, func() {
						inputer.SetStatus(ExceptionExitStatus)
					})

				case RunningStatus:
					// pass
				case StoppingStatus:
					// pass
				case StoppedStatus:
					// pass
				case ExceptionExitStatus:
					common.GoRun(func() {
						cli.executeInputer(ctx, inputer)
					}, func() {
						inputer.SetStatus(ExceptionExitStatus)
					})

				default:
					log.Fatalf("unknown the Inputer status (%d)", inputer.GetStatus())
				}
			}

			cli.inputerLock.RUnlock()
		}
	}

end:
	log.Info("finish the inputer main business cycle")
}
