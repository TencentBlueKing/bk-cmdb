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
	//"time"
)

func init() {

	api.RegisterInputer(business)
	//api.RegisterFrequencyInputer(business, time.Minute)
}

var business = &businessInputer{}

type businessInputer struct {
}

// Name the Inputer name.
// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
func (cli *businessInputer) Name() string {
	return "business_inputer"
}

// Run the input should not be blocked
func (cli *businessInputer) Run(ctx input.InputerContext) *input.InputerResult {

	business, err := api.CreateBusiness("0")
	if nil != err {
		fmt.Println("failed to create business:", err)
		return nil
	}

	business.SetName("demo_business")
	business.SetMaintainer("test_user")
	business.SetDeveloper("test_developer")
	business.SetOperator("test_operator")
	business.SetProductor("test_productor")
	business.SetLifeCycle(api.BusinessLifeCycleOnLine)
	business.SetTester("test_tester")
	business.SetValue("time_zone", "Asia/Shanghai")

	err = business.Save()
	if nil != err {
		fmt.Println("failed to save the business:", err)
	}
	bizID, _ := business.GetBusinessID()
	fmt.Println("the business id:", bizID)
	return nil

}

func (cli *businessInputer) Stop() error {
	return nil
}
