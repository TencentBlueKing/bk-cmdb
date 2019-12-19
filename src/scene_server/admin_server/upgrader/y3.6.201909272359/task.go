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

package y3_6_201909272359

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func taskMigrate(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	existTable, err := db.HasTable(ctx, common.BKTableNameAPITask)
	if err != nil {
		blog.Errorf("has table %s error. err:%s", common.BKTableNameAPITask, err.Error())
		return err
	}

	if !existTable {
		err := db.CreateTable(ctx, common.BKTableNameAPITask)
		if err != nil {
			blog.Errorf("create table %s error. err:%s", common.BKTableNameAPITask, err.Error())
			return err
		}
	}

	indexArr := []dal.Index{
		dal.Index{
			Keys:       map[string]int32{"task_id": 1},
			Name:       "idx_taskID",
			Unique:     true,
			Background: true,
		},
		dal.Index{
			Keys:       map[string]int32{"name": 1, "status": 1, "create_time": 1},
			Name:       "idx_name_status_createTime",
			Background: true,
		},
		dal.Index{
			Keys:       map[string]int32{"status": 1, "last_time": 1},
			Name:       "idx_status_lastTime",
			Background: true,
		},
		dal.Index{
			Keys:       map[string]int32{"name": 1, "flag": 1, "create_time": 1},
			Name:       "idx_name_flag_createTime",
			Background: true,
		},
	}

	for _, index := range indexArr {
		err := db.Table(common.BKTableNameAPITask).CreateIndex(ctx, index)
		if err != nil {
			blog.ErrorJSON("create table %s  index %s error. err:%s", common.BKTableNameAPITask, index, err.Error())
			return err
		}
	}

	return nil
}
