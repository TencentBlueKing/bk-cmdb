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

package synchronize

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
)

type bizSyncer struct {
}

// ParseDataArr parse data array to actual type
func (b *bizSyncer) ParseDataArr(kit *rest.Kit, data any) (any, error) {
	return parseDataArr[mapstr.MapStr](kit, data)
}

// Validate biz sync data
func (b *bizSyncer) Validate(kit *rest.Kit, subRes string, data any) error {
	bizs, ok := data.([]mapstr.MapStr)
	if !ok {
		return kit.CCError.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("data type %T is invalid", data))
	}

	for _, biz := range bizs {
		bizID, err := util.GetInt64ByInterface(biz[common.BKAppIDField])
		if err != nil {
			blog.Errorf("parse biz id(%v) failed, err: %v, rid: %s", biz[common.BKAppIDField], err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField)
		}
		if bizID <= 0 {
			blog.Errorf("biz id(%d) is invalid,  rid: %s", bizID, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField)
		}
	}

	return nil
}

// TableName returns table name for biz syncer
func (b *bizSyncer) TableName(subRes, supplierAccount string) string {
	return common.BKTableNameBaseApp
}

type setSyncer struct {
}

// ParseDataArr parse data array to actual type
func (s *setSyncer) ParseDataArr(kit *rest.Kit, data any) (any, error) {
	return parseDataArr[mapstr.MapStr](kit, data)
}

// Validate set sync data
func (s *setSyncer) Validate(kit *rest.Kit, subRes string, data any) error {
	sets, ok := data.([]mapstr.MapStr)
	if !ok {
		return kit.CCError.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("data type %T is invalid", data))
	}

	bizIDs, tempIDs, parentIDs := make([]int64, 0), make([]int64, 0), make([]int64, 0)

	for _, set := range sets {
		setID, err := util.GetInt64ByInterface(set[common.BKSetIDField])
		if err != nil {
			blog.Errorf("parse set id(%v) failed, err: %v, rid: %s", set[common.BKSetIDField], err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKSetIDField)
		}
		if setID <= 0 {
			blog.Errorf("set id(%d) is invalid,  rid: %s", setID, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKSetIDField)
		}

		bizID, err := util.GetInt64ByInterface(set[common.BKAppIDField])
		if err != nil {
			blog.Errorf("parse biz id(%v) failed, err: %v, rid: %s", set[common.BKAppIDField], err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField)
		}
		bizIDs = append(bizIDs, bizID)

		parentID, err := util.GetInt64ByInterface(set[common.BKParentIDField])
		if err != nil {
			blog.Errorf("parse parent id(%v) failed, err: %v, rid: %s", set[common.BKParentIDField], err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKParentIDField)
		}

		defaultVal, err := util.GetIntByInterface(set[common.BKDefaultField])
		if err != nil {
			blog.Errorf("parse default(%v) failed, err: %v, rid: %s", set[common.BKDefaultField], err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKDefaultField)
		}

		if defaultVal != common.DefaultResSetFlag {
			parentIDs = append(parentIDs, parentID)
		}

		tempIDVal, exists := set[common.BKSetTemplateIDField]
		if !exists {
			continue
		}

		tempID, err := util.GetInt64ByInterface(tempIDVal)
		if err != nil {
			blog.Errorf("parse temp id(%v) failed, err: %v, rid: %s", set[common.BKSetTemplateIDField], err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKSetTemplateIDField)
		}
		if tempID != 0 {
			tempIDs = append(tempIDs, tempID)
		}
	}

	if err := validateDependency(kit, common.BKTableNameBaseApp, common.BKAppIDField, bizIDs); err != nil {
		return err
	}

	if err := validateDependency(kit, common.BKTableNameSetTemplate, common.BKFieldID, tempIDs); err != nil {
		return err
	}

	// validate parent mainline instance id
	parentObj, _, err := getMainlineParentObj(kit, common.BKInnerObjIDSet)
	if err != nil {
		return err
	}

	if parentObj == common.BKInnerObjIDApp {
		return nil
	}

	table := common.GetObjectInstTableName(parentObj, kit.SupplierAccount)
	if err = validateDependency(kit, table, common.BKInstIDField, parentIDs); err != nil {
		return err
	}

	return nil
}

// TableName returns table name for set syncer
func (s *setSyncer) TableName(subRes, supplierAccount string) string {
	return common.BKTableNameBaseSet
}

type moduleSyncer struct {
}

// ParseDataArr parse data array to actual type
func (m *moduleSyncer) ParseDataArr(kit *rest.Kit, data any) (any, error) {
	return parseDataArr[mapstr.MapStr](kit, data)
}

// Validate module sync data
func (m *moduleSyncer) Validate(kit *rest.Kit, subRes string, data any) error {
	modules, ok := data.([]mapstr.MapStr)
	if !ok {
		return kit.CCError.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("data type %T is invalid", data))
	}

	setIDs, categoryIDs, tempIDs := make([]int64, 0), make([]int64, 0), make([]int64, 0)

	for _, module := range modules {
		moduleID, err := util.GetInt64ByInterface(module[common.BKModuleIDField])
		if err != nil {
			blog.Errorf("parse module id(%v) failed, err: %v, rid: %s", module[common.BKModuleIDField], err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKModuleIDField)
		}
		if moduleID <= 0 {
			blog.Errorf("module id(%d) is invalid,  rid: %s", moduleID, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKModuleIDField)
		}

		setID, err := util.GetInt64ByInterface(module[common.BKSetIDField])
		if err != nil {
			blog.Errorf("parse set id(%v) failed, err: %v, rid: %s", module[common.BKSetIDField], err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKSetIDField)
		}
		setIDs = append(setIDs, setID)

		categoryID, err := util.GetInt64ByInterface(module[common.BKServiceCategoryIDField])
		if err != nil {
			blog.Errorf("parse category id(%v) failed, err: %v, rid: %s", module[common.BKServiceCategoryIDField], err,
				kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKServiceCategoryIDField)
		}
		categoryIDs = append(categoryIDs, categoryID)

		tempIDVal, exists := module[common.BKServiceTemplateIDField]
		if !exists {
			continue
		}

		tempID, err := util.GetInt64ByInterface(tempIDVal)
		if err != nil {
			blog.Errorf("parse temp id(%v) failed, err: %v, rid: %s", module[common.BKServiceTemplateIDField], err,
				kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKServiceTemplateIDField)
		}
		if tempID != 0 {
			tempIDs = append(tempIDs, tempID)
		}
	}

	if err := validateDependency(kit, common.BKTableNameBaseSet, common.BKSetIDField, setIDs); err != nil {
		return err
	}

	if err := validateDependency(kit, common.BKTableNameServiceCategory, common.BKFieldID, categoryIDs); err != nil {
		return err
	}

	if err := validateDependency(kit, common.BKTableNameServiceTemplate, common.BKFieldID, tempIDs); err != nil {
		return err
	}

	return nil
}

// TableName returns table name for module syncer
func (m *moduleSyncer) TableName(subRes, supplierAccount string) string {
	return common.BKTableNameBaseModule
}
