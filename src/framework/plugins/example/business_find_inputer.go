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
	"fmt"
	//"time"
)

func init() {

	api.RegisterInputer(businessFind, nil)
	//api.RegisterTimingInputer(business, time.Second*5, nil)
}

var businessFind = &businessFindInputer{}

type businessFindInputer struct {
}

// Name the Inputer name.
// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
func (cli *businessFindInputer) Name() string {
	return "business_inputer"
}

// Run the input should not be blocked
func (cli *businessFindInputer) Run() interface{} {

	busIterator, err := api.FindBusinessLikeName("0", "蓝鲸")
	if nil != err {
		fmt.Println("failed to find:", err)
		return nil
	}

	busIterator.ForEach(func(item *api.BusinessWrapper) error {
		name, _ := item.GetName()
		dev, _ := item.GetMaintainer()
		fmt.Println("inst name:", name, dev)
		return nil
	})

	return nil

}

func (cli *businessFindInputer) Stop() error {
	return nil
}
