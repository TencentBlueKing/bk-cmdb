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
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type objInstSyncer struct {
}

// ParseDataArr parse data array to actual type
func (o *objInstSyncer) ParseDataArr(kit *rest.Kit, data any) (any, error) {
	return parseDataArr[mapstr.MapStr](kit, data)
}

// Validate object instance sync data
func (o *objInstSyncer) Validate(kit *rest.Kit, subRes string, data any) error {
	instances, ok := data.([]mapstr.MapStr)
	if !ok {
		return kit.CCError.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("data type %T is invalid", data))
	}

	for _, inst := range instances {
		instID, err := util.GetInt64ByInterface(inst[common.BKInstIDField])
		if err != nil {
			blog.Errorf("parse obj inst id(%v) failed, err: %v, rid: %s", inst[common.BKInstIDField], err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKInstIDField)
		}
		if instID <= 0 {
			blog.Errorf("obj inst id(%d) is invalid,  rid: %s", instID, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKInstIDField)
		}
	}

	// validate parent mainline instance id
	parentObj, isMainline, err := getMainlineParentObj(kit, subRes)
	if err != nil {
		return err
	}

	if !isMainline {
		return nil
	}

	parentIDs := make([]int64, 0)
	for _, inst := range instances {
		parentID, err := util.GetInt64ByInterface(inst[common.BKParentIDField])
		if err != nil {
			blog.Errorf("parse parent id(%v) failed, err: %v, rid: %s", inst[common.BKParentIDField], err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKParentIDField)
		}
		parentIDs = append(parentIDs, parentID)
	}

	table := common.GetInstTableName(parentObj, kit.SupplierAccount)
	if err = validateDependency(kit, table, common.GetInstIDField(parentObj), parentIDs); err != nil {
		return err
	}

	return nil
}

// TableName returns table name for object instance syncer
func (o *objInstSyncer) TableName(subRes, supplierAccount string) string {
	return common.GetObjectInstTableName(subRes, supplierAccount)
}

type instAsstSyncer struct {
}

// ParseDataArr parse data array to actual type
func (i *instAsstSyncer) ParseDataArr(kit *rest.Kit, data any) (any, error) {
	return parseDataArr[metadata.InstAsst](kit, data)
}

// Validate inst asst sync data
func (i *instAsstSyncer) Validate(kit *rest.Kit, subRes string, data any) error {
	instAssts, ok := data.([]metadata.InstAsst)
	if !ok {
		return kit.CCError.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("data type %T is invalid", data))
	}

	objAsstIDs := make([]string, 0)
	objInstIDMap := make(map[string][]int64)

	for _, instAsst := range instAssts {
		if instAsst.ID <= 0 {
			blog.Errorf("inst asst id(%d) is invalid,  rid: %s", instAsst.ID, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKFieldID)
		}

		objAsstIDs = append(objAsstIDs, instAsst.ObjectAsstID)
		objInstIDMap[instAsst.ObjectID] = append(objInstIDMap[instAsst.ObjectID], instAsst.InstID)
		objInstIDMap[instAsst.AsstObjectID] = append(objInstIDMap[instAsst.AsstObjectID], instAsst.AsstInstID)
	}

	err := validateDependency(kit, common.BKTableNameObjAsst, common.AssociationObjAsstIDField, objAsstIDs)
	if err != nil {
		return err
	}

	for objID, instIDs := range objInstIDMap {
		table := common.GetInstTableName(objID, kit.SupplierAccount)
		if err = validateDependency(kit, table, common.GetInstIDField(objID), instIDs); err != nil {
			return err
		}
	}

	return nil
}

// TableName returns table name for instance association syncer
func (i *instAsstSyncer) TableName(subRes, supplierAccount string) string {
	return common.GetObjectInstAsstTableName(subRes, supplierAccount)
}
