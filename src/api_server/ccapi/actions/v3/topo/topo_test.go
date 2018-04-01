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
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func createInst(objID string, parentID int) {
	param := map[string]interface{}{
		"InstName": "app_test_inst",
		"ParentID": parentID,
	}

	paramsJs, _ := json.Marshal(param)
	rsp, rspErr := hclient.POST(Address+"/inst/test_owner/"+objID, nil, paramsJs)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}
	fmt.Printf("rsp:%s", rsp)
}

func createMainObject(objID string) {

	param := map[string]interface{}{
		"Creator":          "app_test",
		"Modifier":         "app_test",
		"Description":      "app_test, main line",
		"ClassificationID": "test",
		"OwnerID":          "test_owner",
		"ObjID":            objID,
		"ObjName":          fmt.Sprintf("%s_%d", objID, time.Now().UnixNano()),
	}

	paramsJs, _ := json.Marshal(param)

	rsp, rspErr := hclient.POST(Address+"/object", nil, paramsJs)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Printf("rsp:%s", rsp)
}

func createMainModule(objID, associationID string) {

	param := map[string]interface{}{
		"Creator":          "app_test",
		"Modifier":         "app_test",
		"Description":      "app_test, main line",
		"ClassificationID": "test",
		"OwnerID":          "test_owner",
		"ObjID":            objID,
		"ObjName":          fmt.Sprintf("%s_%d", objID, time.Now().UnixNano()),
		"AssociationID":    associationID, // only association with app
	}

	paramsJs, _ := json.Marshal(param)
	rsp, rspErr := hclient.POST(Address+"/topo/model/mainline", nil, paramsJs)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Printf("rsp:%s", rsp)
}

func TestCreateInst(t *testing.T) {

	// init main module line
	//createMainObject("app_test")
	//createMainModule("set_test", "app_test")
	//createMainModule("module_test", "set_test")

	// create parent
	//createInst("app", 0)
	//createInst("set", 136)
	//createInst("module", 137)
}
func TestSearchMainTopo(t *testing.T) {

	rsp, rspErr := hclient.GET(Address+"/topo/model/test_owner", nil, nil)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Printf("rsp:%s", rsp)
}
func TestCreateMainModule(t *testing.T) {

	param := map[string]interface{}{
		"Creator":          "app_test",
		"Modifier":         "app_test",
		"Description":      "app_test, main line",
		"ClassificationID": "test",
		"OwnerID":          "0",
		"ObjID":            fmt.Sprintf("app_test_%d", time.Now().UnixNano()),
		"ObjName":          fmt.Sprintf("app_test_%d", time.Now().UnixNano()),
		"AssociationID":    "app", // only association with app
	}

	paramsJs, _ := json.Marshal(param)
	rsp, rspErr := hclient.POST(Address+"/topo/model/mainline", nil, paramsJs)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Printf("rsp:%s", rsp)
}

func TestDeleteMainLineModule(t *testing.T) {

	objid := "app_test_1515571586120245000"
	rsp, rspErr := hclient.DELETE(Address+"/topo/model/mainline/owners/0/objectids/"+objid, nil, nil)
	if nil != rspErr {
		fmt.Printf("err: %s\n", rspErr.Error())
		return
	}

	fmt.Printf("rsp:%s", rsp)
}
