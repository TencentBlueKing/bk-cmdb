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

// CreateApp create a app
func CreateApp(ownerID string) (int, error) {

	// construct the data
	data := map[string]interface{}{
		common.BKAppIDField:      "cmdb",
		common.BKMaintainersField: "sys_cc;sys_cc_1",
		common.BKDeveloperField:   "sys_developer;sys_developer_1",
		common.BKTesterField:      "sys_tester;sys_tester_1",
		common.BKOperatorField:    "sys_owner;sys_owner_1",
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// create a app
	rst, rstErr := client.POST(fmt.Sprintf("/biz/%s", ownerID), nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to create a app, error info is ", rstErr.Error())
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
		switch id := t[common.BKAppIDField].(type) {
		case float64:
			return int(id), nil
		default:
			fmt.Println("kind:", reflect.TypeOf(t[common.BKAppIDField]).Kind())
		}
	default:
		fmt.Println("kind:", reflect.TypeOf(rstObj.Data))
	}

	return 0, nil
}

// DeleteApp delete a app
func DeleteApp(ownerID string, appID int) error {

	// delete a app
	rst, rstErr := client.DELETE(fmt.Sprintf("/biz/%s/%d", ownerID, appID), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to delete a app, error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal a app, error info is ", jsErr.Error())
		return jsErr
	}

	// print the result
	fmt.Printf("the result %+v\n", rstObj)

	return nil
}

// UpdateApp update a app
func UpdateApp(ownerID string, appID int) error {

	//  construct the data
	data := map[string]interface{}{
		common.BKAppName: "app_new",
	}

	// marshal data
	dataStr, _ := json.Marshal(data)

	// update a app
	rst, rstErr := client.PUT(fmt.Sprintf("/app/%s/%d", ownerID, appID), nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to update a app, error info is ", rstErr.Error())
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

// SearchApp search a app
func SearchApp(ownerID string) error {

	// construct the data
	data := map[string]interface{}{
		"page": map[string]interface{}{
			"start": 0,
			"limit": 100,
			"sort":  common.BKAppNameField,
		},
		"fields": []string{common.BKAppIDField},
		"condition": map[string]interface{}{
			common.BKAppNameField: "cmdb",
		},
		"native": 1,
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// search the app
	rst, rstErr := client.POST(fmt.Sprintf("/app/search/%s", ownerID), nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to search the app, error info is ", rstErr.Error())
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
