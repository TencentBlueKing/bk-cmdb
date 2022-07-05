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

package y3_10_202204271725

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

// addServiceTemplateTableColumn add host_apply_enabled field to service template table.
func addServiceTemplateTableColumn(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	err := db.Table(common.BKTableNameServiceTemplate).AddColumn(ctx, common.HostApplyEnabledField, false)
	if err != nil {
		blog.Errorf("add host_apply_enabled column to service template failed, err: %v", err)
		return err
	}
	return nil
}

// addHostApplyRuleTableColumn the host automatic application table adds the service_template_id field.
func addHostApplyRuleTableColumn(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	err := db.Table(common.BKTableNameHostApplyRule).AddColumn(ctx, common.BKServiceTemplateIDField, 0)
	if err != nil {
		blog.Errorf("add service_template_id column to host apply rule table failed, err: %v", err)
		return err
	}

	// add index
	indexes := []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bizID_ModuleID_serviceTemplateID_attrID",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{common.BKModuleIDField, 1},
				{common.BKServiceTemplateIDField, 1},
				{common.BKAttributeIDField, 1},
			},
			Unique:     true,
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "host_property_under_service_template",
			Keys: bson.D{
				{common.BKServiceTemplateIDField, 1},
				{common.BKAttributeIDField, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bizID_serviceTemplateID_attrID",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{common.BKServiceTemplateIDField, 1},
				{common.BKAttributeIDField, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bizID_moduleID_attrID",
			Keys: bson.D{
				{common.BKAppIDField, 1},
				{common.BKModuleIDField, 1},
				{common.BKAttributeIDField, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "moduleID_attrID",
			Keys: bson.D{
				{common.BKModuleIDField, 1},
				{common.BKAttributeIDField, 1},
			},
			Background: true,
		},
	}

	idxArr, err := db.Table(common.BKTableNameHostApplyRule).Indexes(ctx)
	if err != nil {
		blog.Errorf("get table %s index error, err: %v", common.BKTableNameHostApplyRule, err)
		return err
	}

	idArrMap := make(map[string]types.Index)

	for _, idx := range idxArr {
		idArrMap[idx.Name] = idx
	}
	// It is necessary to delete the previous unique indexes that do not meet the requirements of the new scenario.
	// According to the current factory, there are only two indexes that need to be deleted and rebuilt. They are
	// "idx_unique_bizID_moduleID_attrID" and "idx_unique_bizID_moduleID_attrID". After "host_property_under_module"
	// is deleted, the reconstruction needs to be performed according to the latest naming rules. The corresponding
	// new indexes are "bkcc_idx_bizID_moduleID_attrID" and "bkcc_idx_moduleID_attrID".
	for name, index := range idArrMap {
		if name == "idx_unique_bizID_moduleID_attrID" && index.Unique {
			if err := db.Table(common.BKTableNameHostApplyRule).DropIndex(ctx, index.Name); err != nil &&
				!db.IsNotFoundError(err) {
				blog.Errorf("remove table: %s index: %s error, err: %v, rid: %s", common.BKTableNameHostApplyRule,
					name, err)
				return err
			}
			delete(idArrMap, name)
		}
		if index.Name == "host_property_under_module" && index.Unique {
			if err := db.Table(common.BKTableNameHostApplyRule).DropIndex(ctx, index.Name); err != nil &&
				!db.IsNotFoundError(err) {
				blog.Errorf("remove table: %s index: %s error, err: %v, rid: %s", common.BKTableNameHostApplyRule,
					name, err)
				return err
			}
			delete(idArrMap, name)
		}
	}

	for _, index := range indexes {
		exist := false
		if _, ok := idArrMap[index.Name]; ok {
			exist = true
			break
		}

		// index already exist, skip create
		if exist {
			continue
		}

		if err := db.Table(common.BKTableNameHostApplyRule).CreateIndex(ctx, index); err != nil &&
			!db.IsDuplicatedError(err) {
			blog.Errorf("add host property apply index failed, table: %s, index: %+v, err: %v",
				common.BKTableNameHostApplyRule, index, err)
			return fmt.Errorf("add index failed, table: %s, index: %s, err: %v",
				common.BKTableNameHostApplyRule, index.Name, err)
		}
	}
	return nil
}
