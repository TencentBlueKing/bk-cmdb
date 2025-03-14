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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal/mongo/local"
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

func addModuleData(kit *rest.Kit, db local.DB, bizID int64, moduleNames []string, setID int64) error {

	categoryID := defaultServiceCategoryID
	moduleAdd := make([]mapstr.MapStr, 0)
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
		moduleData[common.BKParentIDField] = setID
		moduleData[common.BKSetIDField] = setID
		moduleData[common.CreateTimeField] = time.Now()
		moduleData[common.LastTimeField] = time.Now()
		moduleAdd = append(moduleAdd, moduleData)
	}
	moduleAuditType := &tools.AuditResType{
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModuleRes,
	}

	needField := &tools.InsertOptions{
		UniqueFields:   []string{common.BKAppIDField, common.BKSetIDField, common.BKModuleNameField},
		IgnoreKeys:     []string{common.BKModuleIDField, common.BKFieldDBID},
		IDField:        []string{common.BKModuleIDField},
		AuditTypeField: moduleAuditType,
		AuditDataField: &tools.AuditDataField{
			BizIDField:   common.BKAppIDField,
			ResIDField:   common.BKModuleIDField,
			ResNameField: common.BKModuleNameField,
		},
	}

	_, err := tools.InsertData(kit, db, common.BKTableNameBaseModule, moduleAdd, needField)
	if err != nil {
		blog.Errorf("insert data for table %s failed,common.BKTableNameAsstDes, err: %v", err)
		return err
	}
	return nil
}
