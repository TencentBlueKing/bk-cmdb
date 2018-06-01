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
 
package topo_test

import (
	"configcenter/src/common/http/httpclient"
	"encoding/json"
	"fmt"
	"testing"
)

const Address = "http://127.0.0.1:50001/api/v1"

var hclient = httpclient.NewHttpClient()

func TestCreateGroup(t *testing.T) {

	input := map[string]interface{}{
		"groupID":    "group_id_test",
		"groupName":  "group_name_test",
		"groupIndex": 0,
		"ownerID":    "test_owner",
		"objectID":   "app_test",
		"isDefault":  false,
	}

	paramsJs, _ := json.Marshal(input)

	rsp, rspErr := hclient.POST(Address+"/objectatt/group/new", nil, paramsJs)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Println("rsp: ", string(rsp))

}

func TestUpdateGroupSelf(t *testing.T) {
	input := map[string]interface{}{
		"condition": map[string]interface{}{
			"ID": 2,
		},
		"data": map[string]interface{}{
			"groupID":    "group_id_test_new",
			"groupName":  "group_name_test_new",
			"groupIndex": 2,
			"isDefault":  false,
		},
	}

	paramsJs, _ := json.Marshal(input)
	rsp, rspErr := hclient.PUT(Address+"/objectatt/group/update", nil, paramsJs)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Println("rsp: ", string(rsp))
}

func TestDeleteGroupSelf(t *testing.T) {

	id := "2"
	rsp, rspErr := hclient.DELETE(Address+"/objectatt/group/groupid/"+id, nil, nil)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Println("rsp: ", string(rsp))
}

func TestUpdateGroupProperty(t *testing.T) {
	input := map[string]interface{}{
		"condition": map[string]interface{}{
			"propertyID": "childid",
			"ownerID":    "test_owner",
			"objectID":   "app_test_third",
		},
		"data": map[string]interface{}{
			"groupID":       "group_id_test",
			"propertyIndex": 1,
		},
	}

	paramsJs, _ := json.Marshal(input)

	rsp, rspErr := hclient.PUT(Address+"/objectatt/group/property", nil, paramsJs)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Println("rsp: ", string(rsp))
}

func TestDeleteGroupProperty(t *testing.T) {

	ownerid := "test_owner"
	objectid := "app_test_third"
	propertyid := "childid"
	groupid := "group_id_test"
	rsp, rspErr := hclient.DELETE(Address+"/objectatt/group/owner/"+ownerid+"/object/"+objectid+"/propertyids/"+propertyid+"/groupids/"+groupid, nil, nil)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Println("rsp: ", string(rsp))
}
