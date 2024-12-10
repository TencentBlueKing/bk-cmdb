/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package data

import (
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal"
)

var moduleMap = map[string]map[string]interface{}{
	common.DefaultResModuleName: {
		common.BKModuleNameField: common.DefaultResModuleName,
		common.BKDefaultField:    common.DefaultResModuleFlag,
	},
	common.DefaultFaultModuleName: {
		common.BKModuleNameField: common.DefaultFaultModuleName,
		common.BKDefaultField:    common.DefaultFaultModuleFlag,
	},
	common.DefaultRecycleModuleName: {
		common.BKModuleNameField: common.DefaultRecycleModuleName,
		common.BKDefaultField:    common.DefaultRecycleModuleFlag,
	},
}

func addModuleData(kit *rest.Kit, db dal.Dal, bizID int64, moduleNames []string, setID uint64) error {

	categoryID := defaultID
	var moduleAdd []interface{}
	var moduleAudit []tools.AuditType
	for _, moduleName := range moduleNames {
		moduleData := moduleMap[moduleName]
		moduleData[common.BKModuleNameField] = moduleName
		moduleData[common.HostApplyEnabledField] = false
		moduleData[common.BKSetTemplateIDField] = common.SetTemplateIDNotSet
		moduleData[common.BKModuleTypeField] = common.DefaultModuleType
		moduleData[common.BKServiceCategoryIDField] = categoryID
		moduleData[common.BKOperatorField] = ""
		moduleData[common.BKBakOperatorField] = ""
		moduleData[common.BKServiceTemplateIDField] = common.ServiceTemplateIDNotSet
		moduleData[common.BKAppIDField] = bizID
		moduleData[common.BKSetIDField] = setID
		moduleData[common.CreateTimeField] = time.Now()
		moduleData[common.LastTimeField] = time.Now()
		moduleAdd = append(moduleAdd, moduleData)
		moduleAudit = append(moduleAudit, tools.AuditType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModuleRes,
		})
	}

	cmpField := &tools.CmpFiled{
		UniqueFields: []string{common.BKAppIDField, common.BKSetIDField, common.BKModuleNameField},
		IgnoreKeys:   []string{common.BKModuleIDField},
		IDField:      common.BKModuleIDField,
	}
	auditDataField := &tools.AuditDataField{
		BusinessID:   common.BKAppIDField,
		ResourceID:   common.BKModuleIDField,
		ResourceName: common.BKModuleNameField,
	}
	_, err := tools.InsertData(kit, db.Shard(kit.ShardOpts()), common.BKTableNameBaseModule, moduleAdd, cmpField,
		moduleAudit, auditDataField)
	if err != nil {
		blog.Errorf("insert data for table %s failed,common.BKTableNameAsstDes, err: %v", err)
		return err
	}
	return nil
}

func getDefaultServiceCategoryID(kit *rest.Kit, db dal.Dal) (int64, error) {
	serviceCategory := metadata.ServiceCategory{}
	filter := map[string]interface{}{
		"is_built_in": true,
		"name":        "Default",
		"bk_parent_id": map[string]interface{}{
			common.BKDBGT: 0,
		},
	}
	if err := db.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceCategory).Find(filter).Fields("id").One(kit.Ctx,
		&serviceCategory); err != nil {
		blog.Errorf("find service category failed, filter: %v, err: %v", filter, err)
		return 0, fmt.Errorf("get default service category failed, err: %v", err)
	}
	return serviceCategory.ID, nil
}
