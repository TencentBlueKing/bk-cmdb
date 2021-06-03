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
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addPlatData(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := common.BKTableNameBasePlat
	blog.Infof("add data for  %s table ", tableName)
	row := map[string]interface{}{
		common.BKCloudNameField: "default area",
		common.BKOwnerIDField:   common.BKDefaultOwnerID,
		common.BKCloudIDField:   common.BKDefaultDirSubArea,
		common.CreateTimeField:  time.Now(),
		common.LastTimeField:    time.Now(),
	}

	// ensure id > 1, 1 is reserved for direct connecting area
	_, err := db.NextSequence(ctx, tableName)
	if err != nil {
		return err
	}

	_, _, err = upgrader.Upsert(ctx, db, tableName, row, "", []string{common.BKCloudNameField, common.BKOwnerIDField}, []string{common.BKCloudIDField})
	if nil != err {
		blog.Errorf("add data for  %s table error  %s", tableName, err)
		return err
	}

	return nil
}
