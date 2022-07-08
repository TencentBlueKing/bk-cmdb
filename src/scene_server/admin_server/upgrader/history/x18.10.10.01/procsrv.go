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
package x18_10_10_01

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"gopkg.in/mgo.v2"
)

func addProcOpTaskTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := "cc_ProcOpTask"
	exists, err := db.HasTable(ctx, tableName)
	if err != nil {
		return err
	}
	if !exists {
		if err = db.CreateTable(ctx, tableName); err != nil && !mgo.IsDup(err) {
			return err
		}
	}
	indexs := []types.Index{
		types.Index{Name: "idx_taskID_gseTaskID", Keys: map[string]int32{common.BKTaskIDField: 1, common.BKGseOpTaskIDField: 1}, Background: true},
	}
	for _, index := range indexs {

		if err = db.Table(tableName).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			return err
		}

	}
	return nil
}
func addProcInstanceModelTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := "cc_ProcInstanceModel"
	exists, err := db.HasTable(ctx, tableName)
	if err != nil {
		return err
	}
	if !exists {
		if err = db.CreateTable(ctx, tableName); err != nil && !mgo.IsDup(err) {
			return err
		}
	}
	indexs := []types.Index{
		types.Index{Name: "idx_bkBizID_bkSetID_bkModuleID_bkHostInstanceID", Keys: map[string]int32{common.BKAppIDField: 1, common.BKSetIDField: 1, common.BKModuleIDField: 1, "bk_host_instance_id": 1}, Background: true},
		types.Index{Name: "idx_bkBizID_bkHostID", Keys: map[string]int32{common.BKAppIDField: 1, common.BKHostIDField: 1}, Background: true},
		types.Index{Name: "idx_bkBizID_bkProcessID", Keys: map[string]int32{common.BKAppIDField: 1, common.BKProcessIDField: 1}, Background: true},
	}
	for _, index := range indexs {
		if err = db.Table(tableName).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			return err
		}
	}
	return nil
}
func addProcInstanceDetailTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := "cc_ProcInstanceDetail"
	exists, err := db.HasTable(ctx, tableName)
	if err != nil {
		return err
	}
	if !exists {
		if err = db.CreateTable(ctx, tableName); err != nil && !mgo.IsDup(err) {
			return err
		}
	}
	indexs := []types.Index{
		types.Index{Name: "idx_bkBizID_bkModuleID_bkProcessID", Keys: map[string]int32{common.BKAppIDField: 1, common.BKModuleIDField: 1, common.BKProcessIDField: 1}, Background: true},
		types.Index{Name: "idx_bkBizID_status", Keys: map[string]int32{common.BKAppIDField: 1, common.BKStatusField: 1}, Background: true},
		types.Index{Name: "idx_bkBizID_bkHostID", Keys: map[string]int32{common.BKAppIDField: 1, common.BKHostIDField: 1}, Background: true},
	}
	for _, index := range indexs {
		if err = db.Table(tableName).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			return err
		}
	}
	return nil
}
