/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package v3v0v9beta1

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func fixesSupplierAccount(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	for _, tablename := range shouldAddSupplierAccountFieldTables {
		condition := map[string]interface{}{
			common.BKOwnerIDField: map[string]interface{}{
				"$in": []interface{}{nil, ""},
			},
		}
		data := map[string]interface{}{
			common.BKOwnerIDField: common.BKDefaultOwnerID,
		}
		err := db.Table(tablename).Update(ctx, condition, data)
		if nil != err {
			return err
		}
	}
	return nil
}

var shouldAddSupplierAccountFieldTables = []string{
	"cc_ApplicationBase",
	"cc_HostBase",
	"cc_ModuleBase",
	"cc_ModuleHostConfig",
	"cc_ObjAsst",
	"cc_ObjAttDes",
	"cc_ObjClassification",
	"cc_ObjDes",
	"cc_ObjectBase",
	"cc_PlatBase",
	"cc_Proc2Module",
	"cc_Process",
	"cc_PropertyGroup",
	"cc_SetBase",
	"cc_Subscription",
	"cc_TopoGraphics",
	"cc_InstAsst",
	"cc_History",
	"cc_HostFavourite",
	"cc_UserAPI",
	"cc_UserCustom",
}
