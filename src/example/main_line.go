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

// CreateMainlineObject create a main line object
func CreateMainlineObject(ownerID, objID, associationObjID, classificationID string) (int, error) {

	// construct a data
	data := map[string]interface{}{
		common.CreatorField:              "sys_cc",
		common.BKOperatorField:          "sys_cc",
		common.BKDescriptionField:       "introduce the main object",
		common.BKClassificationIDField: classificationID,
		common.BKObjIDField:            objID,
		common.BKObjNameField:          "main_line_obj_name",
		common.BKOwnerIDField:          ownerID,
		common.BKAsstObjIDField:       associationObjID,
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// create main line object
	rst, rstErr := client.POST("/topo/model/mainline", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to create main model object, error info is ", rstErr.Error())
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

// DeleteMainlineObject delete a main line object
func DeleteMainlineObject(ownerID, objID string) error {

	// delete the main module object
	rst, rstErr := client.DELETE(fmt.Sprintf("/topo/model/mainline/owners/%s/objectids/%s", ownerID, objID), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to delete the main line object, error info is ", rstErr.Error())
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

// SearchTopoObject search the topo objects by ownerID
func SearchTopoObject(ownerID string) error {

	// search the topo main object
	rst, rstErr := client.POST(fmt.Sprintf("/topo/model/%s", ownerID), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to search the main object topo, error info is ", rstErr.Error())
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

// SearchTopoInst search the main object topo inst
func SearchTopoInst(ownerID string, appID int) error {

	// search the main line topo  inst
	rst, rstErr := client.GET(fmt.Sprintf("/topo/inst/%s/%d", ownerID, appID), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to search the topo inst, error info is ", rstErr.Error())
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

// SearchTopoInstChild search the topo inst child
func SearchTopoInstChild(ownerID, objID string, appID, instID int) error {

	// search the main line topo inst
	rst, rstErr := client.GET(fmt.Sprintf("/topo/inst/child/%s/%s/%d/%d", ownerID, objID, appID, instID), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to search the result , error info is ", rstErr.Error())
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
