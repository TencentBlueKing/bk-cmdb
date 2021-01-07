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

package x08_09_17_01

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

func createTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	for tablename, indexs := range tables {
		exists, err := db.HasTable(ctx, tablename)
		if err != nil {
			return err
		}
		if !exists {
			if err = db.CreateTable(ctx, tablename); err != nil && !db.IsDuplicatedError(err) {
				return err
			}
		}
		for index := range indexs {
			if err = db.Table(tablename).CreateIndex(ctx, indexs[index]); err != nil && !db.IsDuplicatedError(err) {
				return err
			}
		}
	}
	return nil
}

var tables = map[string][]types.Index{
	common.BKTableNameNetcollectDevice: []types.Index{
		{Keys: map[string]int32{"device_id": 1}, Background: true},
		{Keys: map[string]int32{"device_name": 1}, Background: true},
		{Keys: map[string]int32{"bk_supplier_account": 1}, Background: true},
	},

	common.BKTableNameNetcollectProperty: []types.Index{
		{Keys: map[string]int32{"netcollect_property_id": 1}, Background: true},
		{Keys: map[string]int32{"bk_supplier_account": 1}, Background: true},
	},
}
