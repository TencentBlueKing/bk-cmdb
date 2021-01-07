/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */
package datasynchronize

import (
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

type model struct {
	base         *synchronizeAdapter
	dataType     metadata.SynchronizeOperateDataType
	DataClassify string
}

func NewSynchronizeModelAdapter(s *metadata.SynchronizeParameter) dataTypeInterface {

	return &model{
		base:         newSynchronizeAdapter(s),
		dataType:     s.OperateDataType,
		DataClassify: s.DataClassify,
	}
}

func (m *model) PreSynchronizeFilter(kit *rest.Kit) errors.CCError {
	err := m.preSynchronizeFilterBefore(kit)
	if err != nil {
		return err
	}
	return m.base.PreSynchronizeFilter(kit)
}

func (m *model) GetErrorStringArr(kit *rest.Kit) ([]metadata.ExceptionResult, errors.CCError) {

	if len(m.base.errorArray) == 0 {
		return nil, nil
	}

	return m.base.GetErrorStringArr(kit)

}
func (m *model) SaveSynchronize(kit *rest.Kit) errors.CCError {
	// Each model is written separately for subsequent expansion,
	// each type may be processed differently.
	switch m.DataClassify {
	case common.SynchronizeModelTypeClassification:
		return m.saveSynchronizeModelClassification(kit)
	case common.SynchronizeModelTypeAttribute:
		return m.saveSynchronizeModelAttribute(kit)
	case common.SynchronizeModelTypeAttributeGroup:
		return m.saveSynchronizeModelAttributeGroup(kit)
	case common.SynchronizeModelTypeBase:
		return m.saveSynchronizeModelBase(kit)
	default:
		return kit.CCError.Errorf(common.CCErrCoreServiceSyncDataClassifyNotExistError, m.dataType, m.DataClassify)
	}
}

func (m *model) saveSynchronizeModelClassification(kit *rest.Kit) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	// "cc_ObjClassification"
	dbParam.tableName = common.BKTableNameObjClassification
	dbParam.InstIDField = common.BKFieldID
	m.base.saveSynchronize(kit, dbParam)
	return nil
}

func (m *model) saveSynchronizeModelAttribute(kit *rest.Kit) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	// "cc_ObjAttDes"
	dbParam.tableName = common.BKTableNameObjAttDes
	dbParam.InstIDField = common.BKFieldID
	m.base.saveSynchronize(kit, dbParam)
	return nil
}

func (m *model) saveSynchronizeModelAttributeGroup(kit *rest.Kit) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	// cc_PropertyGroup
	dbParam.tableName = common.BKTableNamePropertyGroup
	dbParam.InstIDField = common.BKFieldID
	m.base.saveSynchronize(kit, dbParam)
	return nil
}

func (m *model) saveSynchronizeModelBase(kit *rest.Kit) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	// cc_ObjDes
	dbParam.tableName = common.BKTableNameObjDes
	dbParam.InstIDField = common.BKFieldID
	m.base.saveSynchronize(kit, dbParam)
	return nil
}

func (m *model) preSynchronizeFilterBefore(kit *rest.Kit) errors.CCError {
	return nil
}
func (m *model) preSynchronizeFilterEnd(kit *rest.Kit) errors.CCError {
	return nil
}
