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
	"fmt"
)

func init() {

	api.RegisterInputer(platMgr)
	//api.RegisterFrequencyInputer(platMgr, time.Minute*5)
}

var platMgr = &platInputer{}

type platInputer struct {
}

// Name the Inputer name.
// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
func (cli *platInputer) Name() string {
	return "plat_inputer"
}

// Init initialization method
func (cli *platInputer) Init(ctx input.InputerContext) error {

	return nil
}

// Run the input should not be blocked
func (cli *platInputer) Run(ctx input.InputerContext) *input.InputerResult {

	plat, err := api.CreatePlat("0")

	if nil != err {
		fmt.Println("err:", err)
	}
	plat.SetName("plat(test)")

	err = plat.Save()

	if nil != err {
		fmt.Println("err save:", err)
	}
	return nil
}

func (cli *platInputer) Stop() error {
	return nil
}
