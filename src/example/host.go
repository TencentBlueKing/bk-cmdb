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

// DeleteHostBatch delete a lot of hosts once
func DeleteHostBatch(ownerID string) error {

	// construct the condition
	condition := map[string]interface{}{
		common.BKHostIDField:  "1,2,3,4",
		common.BKOwnerIDField: ownerID,
	}

	// marshal the condition
	conditionStr, _ := json.Marshal(condition)

	// delete hosts
	rst, rstErr := client.DELETE("/hosts/batch", nil, conditionStr)
	if nil != rstErr {
		fmt.Println("failed to delete batch hosts, error info is ", rstErr.Error())
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

// MoveHostsToModule move some hosts to a module
func MoveHostsToModule() error {

	// construct the data
	data := map[string]interface{}{
		common.BKAppIDField:       1,
		common.BKHostIDField:      []int{10, 1},
		common.BKModuleIDField:    []int{1},
		common.BKIsIncrementField: "true",
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// action
	rst, rstErr := client.POST("/hosts/modules", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to move module, error info is ", rstErr.Error())
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

// AssginHostToSpecialModule assgin host to a special module
func AssginHostToSpecialModule() error {

	// construct a data
	data := map[string]interface{}{
		common.BKAppIDField:  11,
		common.BKHostIDField: []int{10, 9},
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// action
	rst, rstErr := client.POST("/hosts/assgin", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to assgin, error info is ", rstErr.Error())
		return rstErr
	}

	// print the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// MoveHostToFaultModule assgin host to a fault module
func MoveHostToFaultModule() error {

	// construct a data
	data := map[string]interface{}{
		common.BKAppIDField:  11,
		common.BKHostIDField: []int{10, 9},
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// action
	rst, rstErr := client.POST("/hosts/faultmodule", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to action, error info is ", rstErr.Error())
		return rstErr
	}

	// print the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// MoveHostToIDleModule assgin host to a idle module
func MoveHostToIDleModule() error {

	// construct a data
	data := map[string]interface{}{
		common.BKAppIDField:  11,
		common.BKHostIDField: []int{10, 9},
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// action
	rst, rstErr := client.POST("/hosts/emptymodule", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to action, error info is ", rstErr.Error())
		return rstErr
	}

	// print the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// MoveHostToResource assgin host to a resource module
func MoveHostToResource() error {

	// construct a data
	data := map[string]interface{}{
		common.BKAppIDField:  11,
		common.BKHostIDField: []int{10, 9},
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// action
	rst, rstErr := client.POST("/hosts/resource", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to action, error info is ", rstErr.Error())
		return rstErr
	}

	// print the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// SearchHosts search hosts by condition
func SearchHosts() error {

	// construct a data
	data := map[string]interface{}{
		common.BKAppIDField: -1,
		"ip": map[string]interface{}{
			"data":  nil,
			"exact": 1,
			"flag":  "innerIp",
		},
		"condition": []interface{}{
			map[string]interface{}{
				common.BKObjIDField: "Object",
				"fields":            []string{},
				"condition": []interface{}{
					map[string]interface{}{
						"field":    common.BKObjIDField,
						"operator": common.BKDBEQ,
						"value":    76,
					},
				},
			},
		},
		"page": map[string]interface{}{
			"start": 0,
			"limit": 100,
			"sort":  common.BKInstNameField,
		},
		"pattern": "",
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// action
	rst, rstErr := client.POST("/hosts/search", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to action, error info is ", rstErr.Error())
		return rstErr
	}

	// print the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// UpdateHosts update the hosts
func UpdateHosts() error {

	// construct a data
	data := map[string]interface{}{
		"HostName":           "host_name_test",
		common.BKHostIDField: []int{10, 9},
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// action
	rst, rstErr := client.POST("/hosts/batch", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to action, error info is ", rstErr.Error())
		return rstErr
	}

	// print the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// SearchHostDetail search host detail
func SearchHostDetail(ownerID, hostID string) error {

	// action
	rst, rstErr := client.POST(fmt.Sprintf("/hosts/%s/%s", ownerID, hostID), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to action, error info is ", rstErr.Error())
		return rstErr
	}

	// print the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// SearchBaseHostDetail search base host detail
func SearchBaseHostDetail(hostID string) error {

	// action
	rst, rstErr := client.GET(fmt.Sprintf("/hosts/baseinfo/search/%s", hostID), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to action, error info is ", rstErr.Error())
		return rstErr
	}

	// print the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// SearchHostSnapshot search host snapshot
func SearchHostSnapshot(hostID string) error {

	// action
	rst, rstErr := client.GET(fmt.Sprintf("/host/snapshot/%s", hostID), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to action, error info is ", rstErr.Error())
		return rstErr
	}

	// print the result
	var rstObj result
	if jsErr := json.Unmarshal(rst, &rstObj); nil != jsErr {
		fmt.Println("failed to unmarshal the result, error info is ", jsErr.Error())
		return jsErr
	}

	fmt.Printf("the result is %+v\n", rstObj)
	return nil
}

// CreateHostSearchHistory create host search history
func CreateHostSearchHistory() error {

	// construct the data
	data := map[string]interface{}{
		"Content": `{"` + common.BKHostIDField + `":"10"}`,
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)

	// create a new history
	rst, rstErr := client.POST("/host/history", nil, dataStr)
	if nil != rstErr {
		fmt.Println("failed to create a new history, error info is ", rstErr.Error())
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

// SearchHostHistory search host history
func SearchHostHistory() error {

	// serch hosts history
	skip := 0
	limit := 100
	rst, rstErr := client.GET(fmt.Sprintf("/host/history/%d/%d", skip, limit), nil, nil)
	if nil != rstErr {
		fmt.Println("failed to search the host history, error info is ", rstErr.Error())
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

// HostFavourite search the host favourite
func HostFavourite() error {

	// construct the data
	data := map[string]interface{}{
		"info": map[string]interface{}{
			"exact_search": true,
			"inner_ip":     true,
			"outer_ip":     true,
			"ip_list":      []string{"1.1.1.1", "2.2.2.2"},
		},
		"inner_ip": true,
		"outer_ip": true,
		"query_params": []map[string]interface{}{
			map[string]interface{}{"object_id": "host",
				"field":    "operator_system",
				"operator": "$in",
				"value":    "123",
			},
		},
		"operator":   "$in",
		"name":       "delete",
		"is_default": 1,
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)
	// save the favourite
	rst, rstErr := client.POST("/hosts/favorites", nil, dataStr)
	if rstErr != nil {
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

// HostUpdateFavourite update the host favourite
func HostUpdateFavourite(id string) error {

	// contruct the data
	data := map[string]interface{}{
		"info": map[string]interface{}{
			"exact_search": true,
			"inner_ip":     true,
			"outer_ip":     true,
			"ip_list":      []string{"1.1.1.1", "2.2.2.2"},
		},
		"inner_ip": true,
		"outer_ip": true,
		"query_params": []map[string]interface{}{
			map[string]interface{}{"object_id": "host",
				"field":    "operator_system",
				"operator": "$in",
				"value":    "123",
			},
		},
		"operator":   "$in",
		"name":       "delete",
		"is_default": 1,
	}

	// marshal the data
	dataStr, _ := json.Marshal(data)
	// save the favourite
	rst, rstErr := client.POST("/hosts/favorites/"+id, nil, dataStr)
	if rstErr != nil {
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

// HostSearchFavourite search the host favourite
func HostSearchFavourite() error {

	// contruct the data
	condition := map[string]interface{}{
		"condition": map[string]interface{}{
			"is_default": 1,
			"name":       "saved_name",
		},
		"limit": 10,
		"start": 0,
	}

	// marshal the data
	conditionStr, _ := json.Marshal(condition)
	// save the favourite
	rst, rstErr := client.POST("/hosts/favorites/search", nil, conditionStr)
	if rstErr != nil {
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

// DeleteHostFavourite delete the host favourite by id
func DeleteHostFavourite(id string) error {

	// delete operation
	rst, rstErr := client.DELETE("/hosts/favorites/"+id, nil, nil)
	if rstErr != nil {
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

// UpdateHostFavouriteIncr get favourite used count
func UpdateHostFavouriteIncr(id string) error {

	// search operation
	rst, rstErr := client.PUT("/hosts/favorites/"+id+"incr", nil, nil)
	if rstErr != nil {
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
