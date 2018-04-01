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
 
package instdata

import (
	"configcenter/src/storage"
)

var DataH storage.DI

// GetHostCntByCondition query host count by condition
func GetHostCntByCondition(condition interface{}) (int, error) {
	cnt, err := DataH.GetCntByCondition("cc_HostBase", condition)
	if nil != err {
		return 0, err
	}
	return cnt, nil
}

// DelHostByCondition delete host by condition
func DelHostByCondition(condition interface{}) error {
	err := DataH.DelByCondition("cc_HostBase", condition)
	if nil != err {
		return err
	}
	return nil
}

// UpdateHostByCondition update host by condition
func UpdateHostByCondition(data interface{}, condition interface{}) error {
	err := DataH.UpdateByCondition("cc_HostBase", data, condition)
	if nil != err {
		return err
	}
	return nil
}

// GetHostByCondition query
func GetHostByCondition(fields []string, condition, result interface{}, sort string, skip, limit int) error {
	return DataH.GetMutilByCondition("cc_HostBase", fields, condition, result, sort, skip, limit)
}

func GetOneHostByCondition(fields []string, condition, result interface{}) error {
	return DataH.GetOneByCondition("cc_HostBase", fields, condition, result)
}

// CreateHost create host
func CreateHost(input interface{}, idName *string) (int, error) {
	hostID, err := DataH.GetIncID("cc_HostBase")
	if err != nil {
		return 0, err
	}
	inputc := input.(map[string]interface{})
	inputc["ObjectID"] = hostID
	*idName = "ObjectID"
	DataH.Insert("cc_HostBase", inputc)
	return int(hostID), nil
}
