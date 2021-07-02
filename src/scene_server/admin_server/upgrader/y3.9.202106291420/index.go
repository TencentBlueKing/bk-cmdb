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
package y3_9_202106291420

import (
	"context"
	"fmt"
	"reflect"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

func addIndexex(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {

	indexes := map[string]types.Index{
		"cc_ServiceInstance": {
			Name:       "bkcc_idx_bkBizID_ID",
			Background: true,
			Keys:       map[string]int32{"bk_biz_id": 1, "id": 1},
		},

		"cc_ModuleHostConfig": {
			Name:       "bkcc_idx_bkBizID_bkHostID",
			Background: true,
			Keys:       map[string]int32{"bk_biz_id": 1, "bk_host_id": 1},
		},
	}

	for tableName, index := range indexes {

		dbIndexes, err := db.Table(tableName).Indexes(ctx)
		if err != nil {
			blog.ErrorJSON("find collection(%s) index error. err: %s", tableName, err.Error())
			return err
		}
		indexExist := false
		for _, dbIndex := range dbIndexes {
			if dbIndex.Name == index.Name {
				blog.InfoJSON("start  collection(%s) equal db index(%s) logic index(%s)", tableName, dbIndex, index)
				if reflect.DeepEqual(dbIndex.Keys, index.Keys) {
					indexExist = true
					continue
				} else {
					err := fmt.Errorf("collection(%s) index(%s) has exists, but keys not equal", tableName, index.Name)
					blog.ErrorJSON("%s, db index: %s, logic index: %s", err.Error(), dbIndex, index)
					return err
				}
			}
		}
		if !indexExist {
			if err := db.Table(tableName).CreateIndex(ctx, index); err != nil {
				blog.ErrorJSON("collection(%s)  create index(%s) error. index: %s, err: %s", tableName, index.Name, index, err.Error())
				return err
			}
		}

	}

	return
}
