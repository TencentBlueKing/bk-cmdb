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
	"configcenter/src/common"
	"encoding/json"
	"fmt"
	"reflect"
)

// CreateObject create a common object
func CreateObject() (int, error) {

	// construct object data
	data := map[string]interface{}{
		common.CreatorField:              "sys_cc",
		common.ModifierField:             "sys_cc",
		common.BKDescriptionField:       "example data",
		common.BKClassificationIDField: "sys_cc_cls",
		common.BKObjIDField:            "sys_cc_objid",
		common.BKObjNameField:          "sys_cc_objname",
		common.BKOwnerIDField:          "sys_cc",
		common.BKObjIconField:          "obj_icon",
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// save the data
	rst, rstErr := client.POST("/object", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to create object, error info is ", rstErr.Error())
		return 0, rstErr
	}
	// unmarshal the data
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return 0, jsErr
	}

	// print the result
	fmt.Printf("the result is %+v\n", rstObj)

	// parse the data and return the id
	switch t := rstObj.Data.(type) {
	case map[string]interface{}:
		switch id := t[common.BKInstIDField].(type) {
		case float64:
			return int(id), nil
		default:
			fmt.Println("kind:", reflect.TypeOf(t[common.BKInstIDField]).Kind())
		}
	default:
		fmt.Println("kind:", reflect.TypeOf(rstObj.Data))
	}
	return 0, nil
}

// DeleteObject delete the object by id
func DeleteObject(id int) error {

	// delete the object
	rst, rstErr := client.DELETE(fmt.Sprintf("/object/%d", id), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to delete the object, error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	// print the data
	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// UpdateObject update the object by id
func UpdateObject(id int) error {

	// update the object
	data := map[string]interface{}{
		common.CreatorField:        "sys_cc_new",
		common.ModifierField:       "sys_cc_new",
		common.BKDescriptionField: "update the object",
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// updte the object
	rst, rstErr := client.PUT(fmt.Sprintf("/object/%d", id), nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to update the object , error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	// print the data
	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// SearchObjects search the objects
func SearchObjects() error {

	// construct the search condition
	condition := map[string]interface{}{
		common.BKObjIDField:   "sys_cc_objid",
		common.BKOwnerIDField: "sys_cc",
	}

	// marshal the condition
	conditionStr, _ := json.Marshal(condition)

	// search the objects by condition
	rst, rstErr := client.POST("/objects", nil, conditionStr)
	if nil != rstErr {
		fmt.Println("failed to search the objects, error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	// print the result
	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// SearchObjectTopo search the object topo
func SearchObjectTopo() error {

	// construct the search condition
	condition := map[string]interface{}{
		common.ModifierField:    "sys_cc",
		common.BKOwnerIDField: "sys_cc",
	}

	// marshal the condition
	conditionStr, _ := json.Marshal(condition)

	// search the object topo by condition
	rst, rstErr := client.POST("/objects/topo", nil, conditionStr)
	if nil != rstErr {
		fmt.Println("failed to search the objects, error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	// print the result
	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}
