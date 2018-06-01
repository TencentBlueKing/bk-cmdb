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

// CreateFieldGroup create a field group
func CreateFieldGroup(ownerID, objID string) (int, error) {

	// construct the data
	data := map[string]interface{}{
		common.BKPropertyGroupIDField:    "group_id_cc",
		common.BKPropertyGroupNameField:  "group_name_cc",
		common.BKPropertyGroupIndexField: -1,
		common.BKObjIDField:               objID,
		common.BKOwnerIDField:             ownerID,
		"bk_isdefault":                       false,
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// create group
	rst, rstErr := client.POST("/objectatt/group/new", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to create group, error info is ", rstErr.Error())
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
		switch id := t[common.BKObjIDField].(type) {
		case float64:
			return int(id), nil
		default:
			fmt.Println("kind:", reflect.TypeOf(t[common.BKObjIDField]).Kind())
		}
	default:
		fmt.Println("kind:", reflect.TypeOf(rstObj.Data))
	}
	return 0, nil
}

// SearchFieldGroup search the field group by condition
func SearchFieldGroup(ownerID, objID string) error {

	// construct the condition
	condition := map[string]interface{}{
		"isDefault": false,
		"page": map[string]interface{}{
			"start": 0,
			"limit": 100,
			"sort":  common.BKPropertyGroupNameField,
		},
	}

	// marshal the condition
	conditionStr, _ := json.Marshal(condition)

	// search field group by condition
	rst, rstErr := client.POST(fmt.Sprintf("/objectatt/group/property/owner/%s/object/%s", ownerID, objID), nil, conditionStr)
	if nil != rstErr {
		fmt.Println("failed to search the field group , error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result , error info is ", jsErr.Error())
		return jsErr
	}

	// print the result
	fmt.Printf("the result is %+v\n", rstObj)

	return nil
}

// UpdateFieldGroup update a field group
func UpdateFieldGroup() error {

	// construct the data
	data := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKPropertyGroupIDField: "group_id_cc",
		},
		"data": map[string]interface{}{
			common.BKPropertyGroupNameField: "group_name_cc_new",
		},
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// update the field group
	rst, rstErr := client.PUT("/objectatt/group/update", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to update the field group , error info is ", rstErr.Error())
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

// DeleteFieldGroup delete the field group by id
func DeleteFieldGroup(id int) error {

	// delete the field group
	rst, rstErr := client.DELETE(fmt.Sprintf("/objectatt/group/groupid/%d", id), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to delete the field group, error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, rstObj); nil != jsErr {
		fmt.Println("failed unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	// print the result
	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// UpdatePropertyGroup update property group
func UpdatePropertyGroup(ownerID, objectID string) error {

	// construct the data
	data := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKPropertyGroupNameField: "sys_cc_field",
			common.BKOwnerIDField:            ownerID,
			common.BKObjIDField:              objectID,
		},
		"data": map[string]interface{}{
			common.BKPropertyGroupIDField:    "group_id_cc",
			common.BKPropertyGroupIndexField: -1,
		},
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// update the data
	rst, rstErr := client.PUT("/objectatt/group/property", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to update the field group, error info is ", rstErr.Error())
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

// DeletePropertyGroup delete the property group
func DeletePropertyGroup(ownerID, objID, propertyID, groupID string) error {

	// delete the property group
	rst, rstErr := client.DELETE(fmt.Sprintf("/objectatt/group/owner/%s/object/%s/propertyids/%s/groupids/%s", ownerID, objID, propertyID, groupID), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to delete the property group, error info is ", rstErr.Error())
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
