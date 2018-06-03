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

package example

import (
	"configcenter/src/framework/api"
	"configcenter/src/framework/core/input"
	"configcenter/src/framework/core/types"

	"fmt"
	//"time"
)

func init() {

	api.RegisterInputer(eve)
	//api.RegisterFrequencyInputer(business, time.Minute)
}

var eve = &eventInputer{}

type eventInputer struct {
	key types.EventKey
}

// Init initialization method
func (cli *eventInputer) Init(ctx input.InputerContext) error {

	cli.key = ctx.RegisterEvent(types.EventHostType, cli.eventCallback)

	return nil
}

func (cli *eventInputer) eventCallback(eves []*types.Event) error {

	for _, item := range eves {
		curData := item.GetCurrData()
		preData := item.GetPreData()
		//fmt.Println("cur:", string(curData.ToJSON()))
		//fmt.Println("pre:", string(preData.ToJSON()))
		more, less, changes := curData.Different(preData)
		fmt.Println("more:", string(more.ToJSON()))
		fmt.Println("less:", string(less.ToJSON()))
		fmt.Println("changes:", string(changes.ToJSON()))

	}
	return nil
}

// Name the Inputer name.
// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
func (cli *eventInputer) Name() string {
	return "eventer_inputer"
}

// Run the input should not be blocked
func (cli *eventInputer) Run(ctx input.InputerContext) *input.InputerResult {

	fmt.Println("hello event inputer")
	return nil
}

func (cli *eventInputer) Stop() error {

	return nil
}
