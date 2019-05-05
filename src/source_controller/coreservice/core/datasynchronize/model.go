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
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type model struct {
	base         *synchronizeAdapter
	dataType     metadata.SynchronizeOperateDataType
	dbProxy      dal.RDB
	DataClassify string
}

func NewSynchronizeModelAdapter(s *metadata.SynchronizeParameter, dbProxy dal.RDB) dataTypeInterface {

	return &model{
		base:         newSynchronizeAdapter(s, dbProxy),
		dataType:     s.OperateDataType,
		DataClassify: s.DataClassify,
		dbProxy:      dbProxy,
	}
}

func (m *model) PreSynchronizeFilter(ctx core.ContextParams) errors.CCError {
	err := m.preSynchronizeFilterBefore(ctx)
	if err != nil {
		return err
	}
	return m.base.PreSynchronizeFilter(ctx)
}

func (m *model) GetErrorStringArr(ctx core.ContextParams) ([]metadata.ExceptionResult, errors.CCError) {

	if len(m.base.errorArray) == 0 {
		return nil, nil
	}

	return m.base.GetErrorStringArr(ctx)

}
func (m *model) SaveSynchronize(ctx core.ContextParams) errors.CCError {
	// Each model is written separately for subsequent expansion,
	// each type may be processed differently.
	switch m.DataClassify {
	case common.SynchronizeModelTypeClassification:
		return m.saveSynchronizeModelClassification(ctx)
	case common.SynchronizeModelTypeAttribute:
		return m.saveSynchronizeModelAttribute(ctx)
	case common.SynchronizeModelTypeAttributeGroup:
		return m.saveSynchronizeModelAttributeGroup(ctx)
	case common.SynchronizeModelTypeBase:
		return m.saveSynchronizeModelBase(ctx)
	default:
		return ctx.Error.Errorf(common.CCErrCoreServiceSyncDataClassifyNotExistError, m.dataType, m.DataClassify)
	}
}

func (m *model) saveSynchronizeModelClassification(ctx core.ContextParams) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	// "cc_ObjClassification"
	dbParam.tableName = common.BKTableNameObjClassifiction
	dbParam.InstIDField = common.BKFieldID
	m.base.saveSynchronize(ctx, dbParam)
	return nil
}

func (m *model) saveSynchronizeModelAttribute(ctx core.ContextParams) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	// "cc_ObjAttDes"
	dbParam.tableName = common.BKTableNameObjAttDes
	dbParam.InstIDField = common.BKFieldID
	m.base.saveSynchronize(ctx, dbParam)
	return nil
}

func (m *model) saveSynchronizeModelAttributeGroup(ctx core.ContextParams) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	// cc_PropertyGroup
	dbParam.tableName = common.BKTableNamePropertyGroup
	dbParam.InstIDField = common.BKFieldID
	m.base.saveSynchronize(ctx, dbParam)
	return nil
}

func (m *model) saveSynchronizeModelBase(ctx core.ContextParams) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	// cc_ObjDes
	dbParam.tableName = common.BKTableNameObjDes
	dbParam.InstIDField = common.BKFieldID
	m.base.saveSynchronize(ctx, dbParam)
	return nil
}

func (m *model) preSynchronizeFilterBefore(ctx core.ContextParams) errors.CCError {
	return nil
}
func (m *model) preSynchronizeFilterEnd(ctx core.ContextParams) errors.CCError {
	return nil
}
