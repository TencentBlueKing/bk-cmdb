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

package y3_15_202607151050

import (
	"configcenter/pkg/tenant"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
)

// updateObjAttDesBoolDefaultValue updates the default value for boolean attributes
// in the object attribute description table (cc_ObjAttDes) for each tenant.
//
//	bk_property_type:bool
//	- If option=false , set default=false and option=nil
//	- If option=true  , set default=true and option=nil
//	- If option=nil   , nothing
func updateObjAttDesBoolDefaultValue(kit *rest.Kit, db dal.Dal) error {
	return tenant.ExecForAllTenants(func(tenantID string) error {
		tenantKit := kit.NewKit().WithTenant(tenantID)
		tenantDB := db.Shard(tenantKit.ShardOpts())

		if err := updateBoolAttrDefaultValue(tenantKit, tenantDB, false); err != nil {
			return err
		}

		if err := updateBoolAttrDefaultValue(tenantKit, tenantDB, true); err != nil {
			return err
		}

		return nil
	})
}

func updateBoolAttrDefaultValue(tenantKit *rest.Kit, tenantDB local.DB, option bool) error {
	updateCond := map[string]interface{}{
		common.BKPropertyTypeField: common.FieldTypeBool,
		common.BKOptionField:       option,
	}
	updateData := map[string]interface{}{
		common.BKDefaultField: option,
		common.BKOptionField:  nil,
	}
	if err := tenantDB.Table(common.BKTableNameObjAttDes).Update(tenantKit.Ctx, updateCond, updateData); err != nil {
		blog.Errorf("update tenant %s bool attribute with option=%v failed, err: %v, cond: %v, updateData: %v", option,
			tenantKit.TenantID, err, updateCond, updateData)
		return err
	}
	return nil
}
