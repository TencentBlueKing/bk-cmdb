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

package y3_10_202203111313

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

// addServiceTemplateTableColumn Added host_apply_enabled field to service template table.
func addServiceTemplateTableColumn(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	err := db.Table(common.BKTableNameServiceTemplate).AddColumn(ctx, common.HostApplyEnabledField, false)
	if err != nil {
		blog.Errorf("add host_apply_enabled column to service template failed, err: %v", err)
		return err
	}
	return nil
}

// addHostApplyRuleTableColumn The host automatic application table adds the service_template_id field.
func addHostApplyRuleTableColumn(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	err := db.Table(common.BKTableNameHostApplyRule).AddColumn(ctx, common.BKServiceTemplateIDField, 0)
	if err != nil {
		blog.Errorf("add service_template_id column to host apply rule table failed, err: %v", err)
		return err
	}

	// add index
	indexes := []types.Index{
		{
			Keys: map[string]int32{
				common.BKServiceTemplateIDField: 1,
				common.BKAttributeIDField:       1,
			},
			Name:       common.CCLogicIndexNamePrefix + "host_property_under_service_template",
			Background: true,
		},
		{
			Keys: map[string]int32{
				common.BKAppIDField:             1,
				common.BKServiceTemplateIDField: 1,
				common.BKAttributeIDField:       1,
			},
			Background: true,
			Name:       common.CCLogicIndexNamePrefix + "bizID_serviceTemplateID_attrID",
		},
		{
			Keys: map[string]int32{
				common.BKServiceTemplateIDField: 1,
			},
			Name:       common.CCLogicIndexNamePrefix + common.BKServiceTemplateIDField,
			Background: true,
		},
		{
			Keys: map[string]int32{
				common.BKAppIDField:             1,
				common.BKServiceTemplateIDField: 1,
			},
			Name:       common.CCLogicIndexNamePrefix + "bizID_serviceTemplateID",
			Background: true,
		},
		// complement the composite index of BizID and moduleID.
		{
			Keys: map[string]int32{
				common.BKAppIDField:    1,
				common.BKModuleIDField: 1,
			},
			Name:       common.CCLogicIndexNamePrefix + "bizID_ModuleID",
			Background: true,
		},
	}

	idxArr, err := db.Table(common.BKTableNameHostApplyRule).Indexes(ctx)
	if err != nil {
		blog.Errorf("get table %s index error, err: %v", common.BKTableNameHostApplyRule, err)
		return err
	}

	for _, index := range indexes {
		exist := false
		for _, existIdx := range idxArr {
			if existIdx.Name == index.Name {
				exist = true
				break
			}
		}
		// index already exist, skip create
		if exist {
			continue
		}

		if err := db.Table(common.BKTableNameHostApplyRule).CreateIndex(ctx, index); err != nil {
			blog.Errorf("add host property apply index failed, table: %s, index: %+v, err: %v",
				common.BKTableNameHostApplyRule, index, err)
			return fmt.Errorf("add index failed, table: %s, index: %s, err: %v",
				common.BKTableNameHostApplyRule, index.Name, err)
		}
	}
	return nil
}
