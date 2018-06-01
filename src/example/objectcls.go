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

// CreateObjectCls create a classification for a object
func CreateObjectCls() (int, error) {

	// construct object classification
	data := map[string]interface{}{
		common.BKClassificationIDField:   "sys_cc_cls",
		common.BKClassificationNameField: "sys_cc_cls_name",
		common.BKClassificationIconField: "sys_cc_cls_icon",
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// save the data
	rst, rstErr := client.POST("/object/classification", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to create the classification, error info is ", rstErr.Error())
		return 0, rstErr
	}

	// unmarshal the data
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return 0, jsErr
	}

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

// DeleteObjectCls delete the object classification by id
func DeleteObjectCls(id int) error {

	// delete the object classification
	rst, rstErr := client.DELETE(fmt.Sprintf("/objectcls/%d", id), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to create the classification, error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the data
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	return nil
}

// UpdateObjectCls update the object classification by id
func UpdateObjectCls(id int) error {

	// construct the data
	data := map[string]interface{}{
		common.BKClassificationNameField: "sys_cc_cls_name_new",
		common.BKClassificationIconField: "sys_cc_cls_icon_new",
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// update the data
	rst, rstErr := client.PUT(fmt.Sprintf("/object/classification/%d", id), nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to update the classification, error info is ", rstErr.Error())
		return rstErr
	}

	// unmarshal the data
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	// print the data
	fmt.Printf("result:%+v\n", rstObj)
	return nil

}

// SearchObjectCls search all objects by classificationid
func SearchObjectCls(id int) error {

	// construct the data
	data := map[string]interface{}{
		"ID": id,
		"page": map[string]interface{}{
			"sort":  "id",
			"limit": 5,
			"start": 0,
		},
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// search the data by condition
	rst, rstErr := client.POST("/object/classifications", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to search the data , error info is ", rstErr.Error())
		return rstErr
	}

	// print the data
	fmt.Printf("result:%+v", rst)
	return nil
}

// SearchObjectWithAssociationByCls select all object with association by classification
func SearchObjectWithAssociationByCls(ownerID string) error {

	// construct the data
	data := map[string]interface{}{
		common.BKClassificationNameField: "sys_cc_cls_name",
		common.BKClassificationIDField:   "sys_cc_cls",
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// search the data by condition
	rst, rstErr := client.POST("/object/classification/"+ownerID+"/objects", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to search the data, error info is ", rstErr.Error())
		return rstErr
	}

	// print the data
	fmt.Printf("result:%+v", rst)
	return nil
}
