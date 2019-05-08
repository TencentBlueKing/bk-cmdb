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

	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"gopkg.in/mgo.v2"
)

func createTable(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	for tablename, indexs := range tables {
		exists, err := db.HasTable(tablename)
		if err != nil {
			return err
		}
		if !exists {
			if err = db.CreateTable(tablename); err != nil && !mgo.IsDup(err) {
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

var tables = map[string][]dal.Index{
	"cc_ApplicationBase": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_biz_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_biz_name": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"default": 1}, Background: true},
	},

	"cc_HostBase": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_host_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_host_name": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_host_innerip": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_host_outerip": 1}, Background: true},
	},
	"cc_ModuleBase": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_module_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_module_name": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"default": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_biz_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_supplier_account": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_set_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_parent_id": 1}, Background: true},
	},
	"cc_ModuleHostConfig": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_biz_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_host_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_module_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_set_id": 1}, Background: true},
	},
	"cc_ObjAsst": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_obj_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_asst_obj_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_supplier_account": 1}, Background: true},
	},
	"cc_ObjAttDes": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_obj_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_supplier_account": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"id": 1}, Background: true},
	},
	"cc_ObjClassification": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_classification_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_classification_name": 1}, Background: true},
	},
	"cc_ObjDes": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_obj_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_classification_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_obj_name": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_supplier_account": 1}, Background: true},
	},
	"cc_ObjectBase": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_obj_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_supplier_account": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_inst_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_parent_id": 1}, Background: true},
	},
	"cc_OperationLog": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"op_target": 1, "inst_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_supplier_account": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_biz_id": 1, "bk_supplier_account": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"ext_key": 1, "bk_supplier_account": 1}, Background: true},
	},
	"cc_PlatBase": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_supplier_account": 1}, Background: true},
	},
	"cc_Proc2Module": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_biz_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_process_id": 1}, Background: true},
	},
	"cc_Process": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_process_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_biz_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_supplier_account": 1}, Background: true},
	},
	"cc_PropertyGroup": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_obj_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_supplier_account": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_group_id": 1}, Background: true},
	},
	"cc_SetBase": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_set_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_parent_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_biz_id": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_supplier_account": 1}, Background: true},
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_set_name": 1}, Background: true},
	},
	"cc_Subscription": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"subscription_id": 1}, Background: true},
	},
	"cc_TopoGraphics": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"scope_type": 1, "scope_id": 1, "node_type": 1, "bk_obj_id": 1, "bk_inst_id": 1}, Background: true, Unique: true},
	},
	"cc_InstAsst": []dal.Index{
		dal.Index{Name: "", Keys: map[string]interface{}{"bk_obj_id": 1, "bk_inst_id": 1}, Background: true},
	},

	"cc_Privilege":          []dal.Index{},
	"cc_History":            []dal.Index{},
	"cc_HostFavourite":      []dal.Index{},
	"cc_UserAPI":            []dal.Index{},
	"cc_UserCustom":         []dal.Index{},
	"cc_UserGroup":          []dal.Index{},
	"cc_UserGroupPrivilege": []dal.Index{}, "cc_idgenerator": []dal.Index{}, "cc_System": []dal.Index{}}
