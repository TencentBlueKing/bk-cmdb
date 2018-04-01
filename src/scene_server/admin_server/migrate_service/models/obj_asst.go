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
	"configcenter/src/source_controller/api/metadata"
	dbStorage "configcenter/src/storage"
)

func AddAsstData(tableName, ownerID string, metaCli dbStorage.DI) error {
	blog.Errorf("add data for  %s table ", tableName)
	rows := getAddAsstData(ownerID)
	for _, row := range rows {
		selector := map[string]interface{}{
			common.BKObjIDField:     row.ObjectID,
			common.BKObjAttIDField: row.ObjectAttID,
			common.BKOwnerIDField:   row.OwnerID,
		}
		isExist, err := metaCli.GetCntByCondition(tableName, selector)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tableName, err)
			return err
		}
		if isExist > 0 {
			continue
		}
		id, err := metaCli.GetIncID(tableName)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tableName, err)
			return err
		}
		row.ID = int(id)
		_, err = metaCli.Insert(tableName, row)
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tableName, err)
			return err
		}
	}

	blog.Errorf("add data for  %s table  ", tableName)
	return nil
}

func getAddAsstData(ownerID string) []*metadata.ObjectAsst {

	dataRows := []*metadata.ObjectAsst{
		&metadata.ObjectAsst{ObjectID: common.BKInnerObjIDSet, ObjectAttID: common.BKChildStr, AsstObjID: common.BKInnerObjIDApp},
		&metadata.ObjectAsst{ObjectID: common.BKInnerObjIDModule, ObjectAttID: common.BKChildStr, AsstObjID: common.BKInnerObjIDSet},
		&metadata.ObjectAsst{ObjectID: common.BKInnerObjIDHost, ObjectAttID: common.BKChildStr, AsstObjID: common.BKInnerObjIDModule},
		&metadata.ObjectAsst{ObjectID: common.BKInnerObjIDHost, ObjectAttID: common.BKCloudIDField, AsstObjID: common.BKInnerObjIDPlat},
	}
	for _, r := range dataRows {
		r.OwnerID = ownerID
		//r.AsstObjType = 0
	}

	return dataRows

}
