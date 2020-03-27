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

package x18_12_13_02

import (
	"context"
	"gopkg.in/mgo.v2"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addCloudTaskTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := common.BKTableNameCloudTask
	exists, err := db.HasTable(tableName)

	if err != nil {
		return err
	}
	if !exists {
		if err = db.CreateTable(tableName); err != nil && !mgo.IsDup(err) {
			return err
		}
	}

	indexs := []dal.Index{
		dal.Index{Name: "", Keys: map[string]int32{"bk_task_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]int32{"bk_task_name": 1}, Background: true},
	}

	for _, index := range indexs {
		if err = db.Table(tableName).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			return err
		}
	}
	return nil
}

func addCloudResourceConfirmTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := common.BKTableNameCloudResourceConfirm
	exists, err := db.HasTable(tableName)

	if err != nil {
		return err
	}
	if !exists {
		if err = db.CreateTable(tableName); err != nil && !mgo.IsDup(err) {
			return err
		}
	}

	indexs := []dal.Index{
		dal.Index{Name: "", Keys: map[string]int32{"bk_resource_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]int32{"bk_obj_id": 1}, Background: true},
	}

	for _, index := range indexs {
		if err = db.Table(tableName).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			return err
		}
	}
	return nil

}

func addCloudSyncHistoryTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := common.BKTableNameCloudSyncHistory
	exists, err := db.HasTable(tableName)

	if err != nil {
		return err
	}
	if !exists {
		if err = db.CreateTable(tableName); err != nil && !mgo.IsDup(err) {
			return err
		}
	}

	indexs := []dal.Index{
		dal.Index{Name: "", Keys: map[string]int32{"bk_task_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]int32{"bk_history_id": 1}, Background: true},
	}

	for _, index := range indexs {
		if err = db.Table(tableName).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			return err
		}
	}
	return nil

}

func addCloudConfirmHistoryTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableName := common.BKTableNameResourceConfirmHistory
	exists, err := db.HasTable(tableName)

	if err != nil {
		return err
	}
	if !exists {
		if err = db.CreateTable(tableName); err != nil && !mgo.IsDup(err) {
			return err
		}
	}

	indexs := []dal.Index{
		dal.Index{Name: "", Keys: map[string]int32{"bk_resource_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]int32{"confirm_history_id": 1}, Background: true},
	}

	for _, index := range indexs {
		if err = db.Table(tableName).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			return err
		}
	}
	return nil
}
