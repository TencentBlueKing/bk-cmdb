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
 
package userapi

import (
	"configcenter/src/common"
	"configcenter/src/source_controller/common/commondata"
	"fmt"
	"testing"
)

var (
	c *Client = NewClient("127.0.0.1:50002")
)

func TestCreateUserAPI(t *testing.T) {
	input := common.KvMap{"Name": "test111", "ApplicationID": 12, "User": "admin_default", "Info": "{\"application\":{\"id\":12},\"exact_search\":false,\"inner_ip\":true,\"outer_ip\":true,\"ip_list\":\"\",\"Set\":\"\",\"Module\":\"\"}"}
	_, result, err := c.Create(input)

	if nil != err {
		t.Error(err.Error())
		return
	}

	fmt.Println(result)
}

func TestUpdate(t *testing.T) {
	//b986nctmjrcfcrr7n0v0
	input := common.KvMap{"Name": "test1", "ApplicationID": 12, "ModifyUser": "admin_edit", "Info": "{\"application\":{\"id\":12},\"exact_search\":false,\"inner_ip\":true,\"outer_ip\":true,\"ip_list\":\"\",\"Set\":\"\",\"Module\":\"\"}"}
	_, result, err := c.Update(input, "12", "b986nctmjrcfcrr7n0v0")

	if nil != err {
		t.Error(err.Error())
		return
	}

	fmt.Println(result)
}

func TestSearch(t *testing.T) {
	var input commondata.ObjQueryInput
	input.Start = 0
	input.Limit = 0
	input.Condition = common.KvMap{"ApplicationID": 12}
	_, result, err := c.GetUserAPI(input)
	if nil != err {
		t.Error(err.Error())
		return
	}

	fmt.Println(result)

}

func TestDetail(t *testing.T) {
	_, result, err := c.Detail("12", "b986nctmjrcfcrr7n0v0")
	if nil != err {
		t.Error(err.Error())
		return
	}

	fmt.Println(result)
}

func TestDelete(t *testing.T) {
	_, result, err := c.Delete("12", "b986phtmjrcfcrr7n0vg")
	if nil != err {
		t.Error(err.Error())
		return
	}

	fmt.Println(result)
}
