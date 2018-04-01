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
 
package models

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	dbStorage "configcenter/src/storage"
	"time"
)

func AddPlatData(tableName string, insCli dbStorage.DI, metaCli dbStorage.DI) error {
	blog.Errorf("add data for  %s table ", tableName)
	rows := []map[string]interface{}{
		map[string]interface{}{
			common.BKCloudNameField: "default area",
			common.BKOwnerIDField:   "",
			common.BKCloudIDField:   common.BKDefaultDirSubArea,
			common.CreateTimeField:  time.Now(),
			common.LastTimeField:    time.Now(),
		},
	}
	for _, row := range rows {

		selector := map[string]interface{}{
			common.BKCloudNameField: row[common.BKCloudNameField],
			common.BKOwnerIDField:   row[common.BKOwnerIDField],
		}
		isExist, err := insCli.GetCntByCondition(tableName, selector)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tableName, err)
			return err
		}
		if isExist > 0 {
			return nil
		}

		// ensure id plug > 1, 1Reserved
		platID, _ := getIncID(tableName, metaCli)
		// Direct connecting area id = 1
		if common.BKDefaultDirSubArea == row[common.BKCloudIDField].(int) {
			platID = common.BKDefaultDirSubArea
		}

		row[common.BKCloudIDField] = platID
		_, err = insCli.Insert(tableName, row)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tableName, err)
			return err
		}

		return nil

	}

	blog.Errorf("add data for  %s table  ", tableName)
	return nil
}

func getIncID(tableName string, DataH dbStorage.DI) (int, error) {
	id, err := DataH.GetIncID(tableName)
	if nil != err {
		return 0, err
	}
	return int(id), nil
}
