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
	"gopkg.in/mgo.v2"

	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage"
)

func createTable(db storage.DI, conf *upgrader.Config) (err error) {
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
			if err = db.Index(tablename, &indexs[index]); err != nil && !mgo.IsDup(err) {
				return err
			}
		}
	}
	return nil
}

var tables = map[string][]storage.Index{
	"cc_ApplicationBase": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_biz_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_biz_name"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"default"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},

	"cc_HostBase": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_host_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_host_name"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_host_innerip"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_host_outerip"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_ModuleBase": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_module_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_module_name"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"default"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_biz_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_set_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_parent_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_ModuleHostConfig": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_biz_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_host_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_module_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_set_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_ObjAsst": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_obj_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_asst_obj_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_ObjAttDes": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_obj_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_ObjClassification": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_classification_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_classification_name"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_ObjDes": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_obj_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_classification_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_obj_name"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_ObjectBase": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_obj_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_inst_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_OperationLog": []storage.Index{
		storage.Index{Name: "", Columns: []string{"op_target", "inst_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_biz_id", "bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"ext_key", "bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_PlatBase": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_Proc2Module": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_biz_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_process_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_Process": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_process_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_biz_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_PropertyGroup": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_obj_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_group_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_SetBase": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_set_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_parent_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_biz_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_set_name"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_Subscription": []storage.Index{
		storage.Index{Name: "", Columns: []string{"subscription_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},
	"cc_TopoGraphics": []storage.Index{
		storage.Index{Name: "", Columns: []string{"scope_type", "scope_id", "node_type", "bk_obj_id", "bk_inst_id"}, Type: storage.INDEX_TYPE_BACKGROUP_UNIQUE},
	},
	"cc_InstAsst": []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_obj_id", "bk_inst_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	},

	"cc_Privilege":          []storage.Index{},
	"cc_History":            []storage.Index{},
	"cc_HostFavourite":      []storage.Index{},
	"cc_UserAPI":            []storage.Index{},
	"cc_UserCustom":         []storage.Index{},
	"cc_UserGroup":          []storage.Index{},
	"cc_UserGroupPrivilege": []storage.Index{},
	"cc_idgenerator":        []storage.Index{}, "cc_System": []storage.Index{}}
