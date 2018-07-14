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
	"configcenter/src/framework/core/log"
	"context"
	"time"
)

func (cli *manager) subExecuteInputer(inputer *wrapInputer) error {

	inputObj := inputer.Run(cli.ctx)

	if nil == inputObj {
		return nil
	}

	return inputObj.Err
}

// executeInputer start the Inputer
func (cli *manager) executeInputer(ctx context.Context, inputer *wrapInputer) {

	inputer.SetStatus(RunningStatus)

	log.Infof("the Inputer(%s) will to run", inputer.Name())
	// non timing inputer
	if !inputer.isTiming {
		if err := cli.subExecuteInputer(inputer); nil != err {
			log.Fatalf("the inputer(%s) return some error and exit , %s", inputer.Name(), err.Error())
			inputer.SetStatus(ExceptionExitStatus)
			return
		}
		inputer.SetStatus(StoppedStatus)
		log.Infof("the Inputer(%s) normal exit", inputer.Name())
		return
	}

	log.Infof("the Inputer(%s) is timing runing", inputer.Name())

	cli.subExecuteInputer(inputer) // execute onece
	tick := time.NewTicker(inputer.frequency)

	for {
		//fmt.Println("tick:", tick)
		select {
		case <-ctx.Done():
			inputer.SetStatus(StoppedStatus)
			log.Infof("the Inputer(%s) normal exit", inputer.Name())
			return
		case <-tick.C:
			tick.Stop()
			log.Infof("timing frequency(%s)", inputer.Name())
			if err := cli.subExecuteInputer(inputer); nil != err {
				log.Fatalf("the inputer(%s) return some error and exit , %s", inputer.Name(), err.Error())
				inputer.SetStatus(ExceptionExitStatus)
				return
			}
			tick = time.NewTicker(inputer.frequency)
		}
	}

}
