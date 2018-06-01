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

// CreateModule create a module
func CreateModule(setID, appID, parentID int, ownerID string) (int, error) {

	// construct a module
	data := map[string]interface{}{
		common.BKSetIDField:      setID,
		common.BKAppIDField:      appID,
		common.BKOwnerIDField:    ownerID,
		common.BKInstParentStr:   parentID,
		common.BKModuleNameField: "sys_cc_modulename",
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// create module
	rst, rstErr := client.POST(fmt.Sprintf("/module/%s/%d", ownerID, appID), nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to create module, error info is ", rstErr.Error())
		return 0, rstErr
	}

	// unmarshal the result
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
		switch id := t[common.BKModuleIDField].(type) {
		case float64:
			return int(id), nil
		default:
			fmt.Println("kind:", reflect.TypeOf(t[common.BKModuleIDField]).Kind())
		}
	default:
		fmt.Println("kind:", reflect.TypeOf(rstObj.Data))
	}

	return 0, nil
}

// DeleteModule delete the module by condition
func DeleteModule(appID, setID, moduleID int) error {

	// delete the module by condition
	rst, rstErr := client.DELETE(fmt.Sprintf("/module/%d/%d/%d", appID, setID, moduleID), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to delete the module, error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the module, error info is ", jsErr.Error())
		return jsErr
	}

	// print the result
	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// UpdateModule update the module by condition
func UpdateModule(appID, setID, moduleID int) error {

	// construct the data
	data := map[string]interface{}{
		common.BKModuleNameField: "sys_cc_modulename_new",
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// update the module
	rst, rstErr := client.PUT(fmt.Sprintf("/module/%d/%d/%d", appID, setID, moduleID), nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to update the module, error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the module
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	// print the result
	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// SearchModule search the module by condition
func SearchModule(ownerID string, appID, setID int) error {

	// construct the condition
	condition := map[string]interface{}{
		"fields": []string{
			common.BKModuleIDField,
			common.BKModuleNameField,
		},
		"page": map[string]interface{}{
			"start": 0,
			"limit": 100,
			"sort":  common.BKModuleNameField,
		},
		"condition": map[string]interface{}{
			common.BKModuleNameField: "sys_cc_modulename",
		},
	}

	// marshal the condition
	conditionStr, _ := json.Marshal(condition)

	// search the module
	rst, rstErr := client.POST(fmt.Sprintf("/module/search/%s/%d/%d", ownerID, appID, setID), nil, conditionStr)
	if nil != rstErr {
		fmt.Println("failed to search the module, error info is ", rstErr.Error())
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
