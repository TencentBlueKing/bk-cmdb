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
	"configcenter/src/common/blog"
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/types"
	"context"
	"fmt"
	"reflect"
	"time"
)

func (cli *manager) subExecuteInputer(inputer *wrapInputer) {

	inputObj := inputer.Run()

	// inputer 分：事物、定时、常规实现
	switch t := inputObj.(type) {
	case error:
		log.Errorf("return some errors from the inputer, error info is %s", t.Error())
	case nil:
		//log.Info("return the data is nil")
	case types.MapStr:
		if nil != inputer.putter {
			blog.Errorf("the puter is not set")
			return
		}
		if err := inputer.putter.Put(t); nil != err {
			log.Errorf("puter return error, error info is %s", err.Error())
			if nil != inputer.exception {
				inputer.exception(t, err)
			}
		}
	case types.Saver:
		if err := t.Save(); nil != err {
			blog.Errorf("failed to execute saver, error info is %s", err.Error())
			if nil != inputer.exception {
				inputer.exception(t, err)
			}
		}
	default:
		unknown := reflect.TypeOf(t)
		log.Infof("unknown the type:%s", unknown.Kind())
		if nil != inputer.exception {
			inputer.exception(t, fmt.Errorf("unkown the input data type:%s", unknown.Kind()))
		}
	}
}

// executeInputer start the Inputer
func (cli *manager) executeInputer(ctx context.Context, inputer *wrapInputer) {

	inputer.SetStatus(RunningStatus)

	log.Infof("the Inputer(%s) will to run", inputer.Name())
	// non timing inputer
	if !inputer.isTiming {
		cli.subExecuteInputer(inputer)
		inputer.SetStatus(StoppedStatus)
		log.Infof("the Inputer(%s) normal exit", inputer.Name())
		return
	}
	log.Infof("the Inputer(%s) is timing runing", inputer.Name())
	tick := time.NewTicker(inputer.frequency)

	for {
		//fmt.Println("tick:", tick)
		select {
		case <-ctx.Done():
			inputer.SetStatus(StoppedStatus)
			log.Infof("the Inputer(%s) normal exit", inputer.Name())
			return
		case <-tick.C:
			cli.subExecuteInputer(inputer)
		}
	}

}
