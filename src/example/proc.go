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
)

// AddProc add process info
func AddProc(ownerID, appID string) error {

	// construct the data
	data := map[string]interface{}{
		common.BKProcessNameField: "nginx",
		"Port":     80,
		"BindIP":   "127.0.0.1",
		"Protocol": "tcp",
		"FuncName": "nginx",
		"WorkPath": "/data/cc/ruinning",
		"User":     "cc",
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// add operation
	rst, rstErr := client.POST("/proc/"+ownerID+"/"+appID, nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to create process, error info is ", rstErr.Error())
		return rstErr
	}
	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}
	return nil
}

// SearchProc search process info
func SearchProc(ownerID, appID string) error {

	// construct the data
	data := map[string]interface{}{
		"page": map[string]interface{}{
			"start": 0,
			"limit": 10,
			"sort":  common.BKProcessNameField,
		},
		"fields": []string{"ProcessID", common.BKProcessNameField},
		"condition": map[string]interface{}{
			common.BKAppIDField:       "123",
			common.BKProcessNameField: "nginx",
		},
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// add operation
	rst, rstErr := client.POST("/proc/search/"+ownerID+"/"+appID, nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to create process, error info is ", rstErr.Error())
		return rstErr
	}
	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}
	return nil
}

// GetProcDetailInfo get a process detail info by id
func GetProcDetailInfo(ownerID, appID, procID string) error {

	// get proc detail info
	rst, rstErr := client.GET("/proc/"+ownerID+"/"+appID+"/"+procID, nil, nil)
	if nil != rstErr {
		fmt.Println("failed to create process, error info is ", rstErr.Error())
		return rstErr
	}
	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}
	return nil
}

// DeleteProcDetailInfo delete a process detail info by id
func DeleteProcDetailInfo(ownerID, appID, procID string) error {

	// get proc detail info
	rst, rstErr := client.DELETE("/proc/"+ownerID+"/"+appID+"/"+procID, nil, nil)
	if nil != rstErr {
		fmt.Println("failed to create process, error info is ", rstErr.Error())
		return rstErr
	}
	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}
	return nil
}

// UpdateProcDetailInfo update a process detail info by id
func UpdateProcDetailInfo(ownerID, appID, procID string) error {

	// construct the data
	data := map[string]interface{}{
		common.BKProcessNameField: "nginx",
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// get proc detail info
	rst, rstErr := client.PUT("/proc/"+ownerID+"/"+appID+"/"+procID, nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to create process, error info is ", rstErr.Error())
		return rstErr
	}
	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}
	return nil
}

// SearchProcBindModuleInfo update a process detail info by id
func SearchProcBindModuleInfo(ownerID, appID, procID string) error {

	// get proc detail info
	rst, rstErr := client.GET("/proc/module/"+ownerID+"/"+appID+"/"+procID, nil, nil)
	if nil != rstErr {
		fmt.Println("failed to create process, error info is ", rstErr.Error())
		return rstErr
	}
	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}
	return nil
}

// BindModuleInfo bind a process detail info by id
func BindModuleInfo(ownerID, appID, procID, moduleName string) error {

	// get proc detail info
	rst, rstErr := client.PUT("/proc/module/"+ownerID+"/"+appID+"/"+procID+"/"+moduleName, nil, nil)
	if nil != rstErr {
		fmt.Println("failed to create process, error info is ", rstErr.Error())
		return rstErr
	}
	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}
	return nil
}

// UnBindModuleInfo bind a process detail info by id
func UnBindModuleInfo(ownerID, appID, procID, moduleName string) error {

	// get proc detail info
	rst, rstErr := client.DELETE("/proc/module/"+ownerID+"/"+appID+"/"+procID+"/"+moduleName, nil, nil)
	if nil != rstErr {
		fmt.Println("failed to create process, error info is ", rstErr.Error())
		return rstErr
	}
	// unmarshal the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}
	return nil
}
