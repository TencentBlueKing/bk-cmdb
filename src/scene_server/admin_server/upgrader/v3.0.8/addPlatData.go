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

package v3v0v8

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage"
)

func addPlatData(db storage.DI, conf *upgrader.Config) error {
	tablename := "cc_PlatBase"
	blog.Errorf("add data for  %s table ", tablename)
	rows := []map[string]interface{}{
		map[string]interface{}{
			common.BKCloudNameField: "default area",
			common.BKOwnerIDField:   common.BKDefaultOwnerID,
			common.BKCloudIDField:   common.BKDefaultDirSubArea,
			common.CreateTimeField:  time.Now(),
			common.LastTimeField:    time.Now(),
		},
	}
	for _, row := range rows {
		// ensure id plug > 1, 1Reserved
		platID, err := db.GetIncID(tablename)
		if err != nil {
			return err
		}
		// Direct connecting area id = 1
		if common.BKDefaultDirSubArea == row[common.BKCloudIDField] {
			platID = common.BKDefaultDirSubArea
		}

		row[common.BKCloudIDField] = platID
		_, _, err = upgrader.Upsert(db, tablename, row, "", []string{common.BKCloudNameField, common.BKOwnerIDField}, []string{common.BKCloudIDField})
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tablename, err)
			return err
		}

		return nil

	}

	blog.Errorf("add data for  %s table  ", tablename)
	return nil
}
