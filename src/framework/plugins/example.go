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
	"configcenter/src/framework/core/types"
	"fmt"
	//"time"
)

func init() {

	_, sender, _ := api.CreateCustomOutputer("example_output", func(data types.MapStr) error {
		fmt.Println("outputer:", data)
		return nil
	})

	api.RegisterInputer(target, sender, nil)
	//api.RegisterTimingInputer(target, sender, time.Second*5, nil)
}

var target = &myInputer{}

type myInputer struct {
}

// Description the Inputer description.
// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
func (cli *myInputer) Name() string {
	return "name_myinputer"
}

// Run the input should not be blocked
func (cli *myInputer) Run() interface{} {

	set, err := api.CreateSet("0")
	if nil != err {
		fmt.Println("err:", err.Error())
	}
	set.SetValue("bk_set_desc", "setdesc")
	set.SetValue("description", "descrip")
	set.SetValue("bk_capacity", 15)
	set.SetValue("bk_parent_id", 2)
	set.SetValue("bk_biz_id", 2)
	set.SetValue("bk_set_name", "test_demo")
	err = set.Save()
	if nil != err {
		fmt.Println("err:", err)
	}
	return nil

}

func (cli *myInputer) Stop() error {
	return nil
}
