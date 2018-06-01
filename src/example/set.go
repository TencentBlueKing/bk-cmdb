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

// CreateSet create a set
func CreateSet(appID int) (int, error) {

	// construct a set
	data := map[string]interface{}{
		common.BKSetNameField:  "example_set",
		common.BKOwnerIDField:  "sys_cc",
		common.BKInstParentStr: appID,
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// create the set
	rst, rstErr := client.POST(fmt.Sprintf("/set/%d", appID), nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to create the set, error info is ", rstErr.Error())
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
		switch id := t["id"].(type) {
		case float64:
			return int(id), nil
		default:
			fmt.Println("kind:", reflect.TypeOf(t["id"]).Kind())
		}
	default:
		fmt.Println("kind:", reflect.TypeOf(rstObj.Data))
	}

	return 0, nil
}

// DeleteSet delete the set by id
func DeleteSet(appID, id int) error {

	// delete the set by id
	rst, rstErr := client.DELETE(fmt.Sprintf("/set/%d/%d", appID, id), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to delete the set, the error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	// print the result
	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// UpdateSet update the set by id
func UpdateSet(appID, id int) error {

	// construct the data
	data := map[string]interface{}{
		common.BKSetNameField: "example_set_new",
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// update the set by id
	rst, rstErr := client.PUT(fmt.Sprintf("/set/%d/%d", appID, id), nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to update the set, error info is ", rstErr.Error())
	}

	// print the result
	fmt.Printf("the result is %+v\n", rst)
	return nil
}

// SearchSet search the set by condition
func SearchSet(ownerID string, appID int) error {

	// construct the condition
	condition := map[string]interface{}{
		"fields": []string{
			common.BKSetIDField,
			common.BKSetNameField,
		},
		"page": map[string]interface{}{
			"start": 0,
			"limit": 100,
			"sort":  common.BKSetNameField,
		},
		"condition": map[string]interface{}{
			common.BKSetNameField: "example_set",
		},
	}

	// marshal the condition
	conditionStr, _ := json.Marshal(condition)

	// search set
	rst, rstErr := client.POST(fmt.Sprintf("/set/search/%s/%d", ownerID, appID), nil, conditionStr)
	if nil != rstErr {
		fmt.Println("failed to search sets, error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal result, error info is ", jsErr.Error())
		return jsErr
	}

	// print the result
	fmt.Printf("the result is %+v\n", rstObj)

	return nil
}
