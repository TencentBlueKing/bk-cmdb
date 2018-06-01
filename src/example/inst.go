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

// CreateInst create a common inst
func CreateInst(appID, parentID int, ownerID, objID string) (int, error) {

	// construct the data
	data := map[string]interface{}{
		common.BKInstParentStr: parentID,
		common.BKInstNameField: "sys_cc_inst_name",
		common.BKAppIDField:    appID,
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// create a inst
	rst, rstErr := client.POST(fmt.Sprintf("/inst/%s/%s", ownerID, objID), nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to create a inst, error info is ", rstErr.Error())
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

// DeleteInst delete a inst
func DeleteInst(ownerID, objID string, instID int) error {
	// delete inst
	rst, rstErr := client.DELETE(fmt.Sprintf("/inst/%s/%s/%d", ownerID, objID, instID), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to delete the inst, error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the object
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the resul, error info is ", jsErr.Error())
		return jsErr
	}

	// print the result
	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// UpdateInst update the inst information
func UpdateInst(ownerID, objID string, parentID, appID, instID int) error {

	// construct the data
	data := map[string]interface{}{
		common.BKInstParentStr: parentID,
		common.BKInstNameField: "sys_cc_inst_name_new",
		common.BKAppIDField:    appID,
	}
	// marshal the data
	dataStr, _ := json.Marshal(data)

	// update the inst
	rst, rstErr := client.PUT(fmt.Sprintf("/inst/%s/%s/%d", ownerID, objID, instID), nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to update the inst ,error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the data
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	// print the result
	fmt.Printf("the result is %+v\n", rstObj)

	return nil
}

// SearchInst search the inst by condition
func SearchInst(ownerID, objID string, appID, instID int) error {
	// construct the condition
	condition := map[string]interface{}{
		"fields": []string{},
		"page": map[string]interface{}{
			"start": 0,
			"limit": 100,
			"sort":  common.BKInstNameField,
		},
		"condition": map[string]interface{}{
			common.BKAppIDField: appID,
		},
	}

	// marshal the condition
	conditionStr, _ := json.Marshal(condition)

	// search insts by condition
	rst, rstErr := client.POST(fmt.Sprintf("/inst/search/%s/%s/%d", ownerID, objID, instID), nil, conditionStr)
	if nil != rstErr {
		fmt.Println("failed to search the insts, error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, errror info is ", jsErr.Error())
		return jsErr
	}

	// print the result
	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// SearchAllInsts search all insts
func SearchAllInsts(ownerID, objID string, appID int) error {

	// construct the condition
	condition := map[string]interface{}{
		"fields": []string{},
		"page": map[string]interface{}{
			"start": 0,
			"limit": 100,
			"sort":  common.BKInstNameField,
		},
		"condition": map[string]interface{}{
			common.BKAppIDField: appID,
		},
	}

	// marshal the condition
	conditionStr, _ := json.Marshal(condition)

	// search all inst by objectID
	rst, rstErr := client.POST(fmt.Sprintf("/inst/search/%s/%s", ownerID, objID), nil, conditionStr)

	if nil != rstErr {
		fmt.Println("failed to search all insts, error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, the error info ", jsErr.Error())
		return jsErr
	}

	// print the result
	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}
