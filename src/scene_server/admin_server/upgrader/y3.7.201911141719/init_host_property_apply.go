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

package y3_7_201911141719

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func InitHostPropertyApplyDataModel(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// check attribute exist
	moduleAttributeFilter := map[string]interface{}{
		"bk_obj_id":      common.BKInnerObjIDModule,
		"bk_property_id": common.HostApplyEnabledField,
	}
	count, err := db.Table(common.BKTableNameObjAttDes).Find(moduleAttributeFilter).Count(ctx)
	if err != nil {
		blog.Errorf("InitHostPropertyApplyDataModel failed, count module attribute failed, filter: %+v, err: %s", moduleAttributeFilter, err.Error())
		return fmt.Errorf("count module attribute failed, err: %s", err.Error())
	}
	if count > 0 {
		return nil
	}

	// add module attribute field
	newAttributeID, err := db.NextSequence(ctx, common.BKTableNameObjDes)
	if err != nil {
		blog.Errorf("InitHostPropertyApplyDataModel failed, NextSequence failed, err: %s", err.Error())
		return fmt.Errorf("NextSequence failed, err: %s", err.Error())
	}

	now := time.Now()
	moduleAttribute := map[string]interface{}{
		"id":                  newAttributeID,
		"bk_obj_id":           common.BKInnerObjIDModule,
		"editable":            true,
		"bk_supplier_account": conf.OwnerID,
		"ispre":               true,
		"isreadonly":          false,
		"bk_issystem":         false,
		"bk_property_index":   0,
		"unit":                "",
		"isrequired":          false,
		"bk_property_type":    common.FieldTypeBool,
		"option":              make(map[string]interface{}),
		"bk_property_id":      common.HostApplyEnabledField,
		"bk_property_name":    "主机属性自动应用",
		"bk_property_group":   "default",
		"placeholder":         "是否开启主机属性自动应用",
		"bk_isapi":            true,
		"creator":             conf.User,
		"create_time":         now,
		"last_time":           now,
	}
	if err := db.Table(common.BKTableNameObjAttDes).Insert(ctx, moduleAttribute); err != nil {
		blog.ErrorJSON("InitHostPropertyApplyDataModel failed, insert failed, attribute: %s, err: %s", moduleAttribute, err.Error())
		return fmt.Errorf("db insert failed, err: %s", err.Error())
	}

	// add module flag, default value false
	filter := map[string]interface{}{
		common.HostApplyEnabledField: map[string]interface{}{
			common.BKDBExists: false,
		},
	}
	doc := map[string]interface{}{
		common.HostApplyEnabledField: false,
	}
	if err := db.Table(common.BKTableNameBaseModule).Update(ctx, filter, doc); err != nil {
		blog.Errorf("InitHostPropertyApplyDataModel failed, init module flag failed, doc: %+v, err: %s", doc, err.Error())
		return fmt.Errorf("init module flag failed, err: %s", err.Error())
	}

	// check table exist
	hasTable, err := db.HasTable(ctx, common.BKTableNameHostApplyRule)
	if err != nil {
		blog.Errorf("InitHostPropertyApplyDataModel failed, check table exist failed, err: %s", err.Error())
		return fmt.Errorf("check table exist failed, table: %s, err: %s", common.BKTableNameHostApplyRule, err.Error())
	}
	if hasTable {
		return nil
	}

	// add table
	/*
		- id(pk)
		- bk_biz_id
		- bk_module_id
		- bk_property_id
		- value
	*/
	if err := db.CreateTable(ctx, common.BKTableNameHostApplyRule); err != nil {
		blog.Errorf("InitHostPropertyApplyDataModel failed, create tabled failed, table: %s, err: %s", common.BKTableNameHostApplyRule, err.Error())
		return fmt.Errorf("create table failed, table: %s, err: %s", common.BKTableNameHostApplyRule, err.Error())
	}

	// add index
	indexes := []dal.Index{
		{
			Keys: map[string]int32{
				common.BKFieldID: 1,
			},
			Name:       common.BKFieldID,
			Unique:     true,
			Background: false,
		}, {
			Keys: map[string]int32{
				common.BKModuleIDField: 1,
			},
			Name:       common.BKModuleIDField,
			Unique:     false,
			Background: false,
		}, {
			Keys: map[string]int32{
				common.BKModuleIDField:    1,
				common.BKAttributeIDField: 1,
			},
			Name:       "host_property_under_module",
			Unique:     true,
			Background: false,
		},
	}
	for _, index := range indexes {
		if err := db.Table(common.BKTableNameHostApplyRule).CreateIndex(ctx, index); err != nil {
			blog.Errorf("InitHostPropertyApplyDataModel failed, add index failed, table: %s, index: %+v, err: %s", common.BKTableNameHostApplyRule, index, err.Error())
			return fmt.Errorf("add index failed, table: %s, index: %s, err: %s", common.BKTableNameHostApplyRule, index.Name, err.Error())
		}
	}
	return nil
}
