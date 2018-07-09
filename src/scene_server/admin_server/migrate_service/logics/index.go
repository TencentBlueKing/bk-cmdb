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

package logics

import (
	"configcenter/src/common/core/cc/api"
	"configcenter/src/storage"
)

func CreateIndex() error {
	a := api.GetAPIResource()

	tableindexs := getIndex()
	for tablename, indexs := range tableindexs {
		for _, index := range indexs {
			ii := index
			if err := a.InstCli.Index(tablename, &ii); err != nil {
				return err
			}
		}
	}
	return nil
}

func getIndex() map[string][]storage.Index {
	index := map[string][]storage.Index{}

	index["cc_ApplicationBase"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_biz_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_biz_name"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"default"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}

	index["cc_HostBase"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_host_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_host_name"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_host_innerip"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_host_outerip"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_ModuleBase"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_module_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_module_name"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"default"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_biz_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_set_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_parent_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_ModuleHostConfig"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_biz_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_host_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_module_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_set_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_ObjAsst"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_obj_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_asst_obj_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_ObjAttDes"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_obj_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_ObjClassification"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_classification_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_classification_name"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_ObjDes"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_obj_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_classification_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_obj_name"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_ObjectBase"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_obj_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_inst_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_OperationLog"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_obj_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_biz_id", "bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"ext_key", "bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_PlatBase"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_Proc2Module"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_biz_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_process_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_ApplicationBase"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_process_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_biz_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_PropertyGroup"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_obj_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_group_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_SetBase"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"bk_set_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_parent_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_biz_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_supplier_account"}, Type: storage.INDEX_TYPE_BACKGROUP},
		storage.Index{Name: "", Columns: []string{"bk_set_name"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_Subscription"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"subscription_id"}, Type: storage.INDEX_TYPE_BACKGROUP},
	}
	index["cc_TopoGraphics"] = []storage.Index{
		storage.Index{Name: "", Columns: []string{"scope_type", "scope_id", "node_type", "bk_obj_id", "bk_inst_id"}, Type: storage.INDEX_TYPE_BACKGROUP_UNIQUE},
	}

	return index
}
