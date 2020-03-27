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

type instance struct {
	base         *synchronizeAdapter
	dataType     metadata.SynchronizeOperateDataType
	dbProxy      dal.RDB
	DataClassify string
}

func NewSynchronizeInstanceAdapter(s *metadata.SynchronizeParameter, dbProxy dal.RDB) dataTypeInterface {

	return &instance{
		base:     newSynchronizeAdapter(s, dbProxy),
		dataType: s.OperateDataType,
		// instance data classify
		DataClassify: s.DataClassify,
		dbProxy:      dbProxy,
	}
}

func (inst *instance) PreSynchronizeFilter(ctx core.ContextParams) errors.CCError {
	err := inst.preSynchronizeFilterBefore(ctx)
	if err != nil {
		return err
	}
	return inst.base.PreSynchronizeFilter(ctx)
}

func (inst *instance) GetErrorStringArr(ctx core.ContextParams) ([]metadata.ExceptionResult, errors.CCError) {

	if len(inst.base.errorArray) == 0 {
		return nil, nil
	}

	return inst.base.GetErrorStringArr(ctx)

}

func (inst *instance) SaveSynchronize(ctx core.ContextParams) errors.CCError {
	// Each model is written separately for subsequent expansion,
	// each model may be processed differently.

	switch inst.DataClassify {
	case common.BKInnerObjIDApp:
		return inst.saveSynchronizeAppInstance(ctx)
	case common.BKInnerObjIDSet:
		return inst.saveSynchronizeSetInstance(ctx)
	case common.BKInnerObjIDModule:
		return inst.saveSynchronizeModuleInstance(ctx)
	case common.BKInnerObjIDProc:
		return inst.saveSynchronizeModuleInstance(ctx)
	case common.BKInnerObjIDPlat:
		return inst.saveSynchronizePlatInstance(ctx)
	case common.BKInnerObjIDHost:
		return inst.saveSynchronizeHostInstance(ctx)
	default:
		return inst.saveSynchronizeObjectInstance(ctx)
	}
}

func (inst *instance) saveSynchronizeAppInstance(ctx core.ContextParams) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	dbParam.tableName = common.BKTableNameBaseApp
	dbParam.InstIDField = common.BKAppIDField
	inst.base.saveSynchronize(ctx, dbParam)
	return nil
}

func (inst *instance) saveSynchronizeSetInstance(ctx core.ContextParams) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	dbParam.tableName = common.BKTableNameBaseSet
	dbParam.InstIDField = common.BKSetIDField
	inst.base.saveSynchronize(ctx, dbParam)
	return nil
}

func (inst *instance) saveSynchronizeModuleInstance(ctx core.ContextParams) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	dbParam.tableName = common.BKTableNameBaseModule
	dbParam.InstIDField = common.BKModuleIDField
	inst.base.saveSynchronize(ctx, dbParam)
	return nil
}

func (inst *instance) saveSynchronizeProcessInstance(ctx core.ContextParams) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	dbParam.tableName = common.BKTableNameBaseProcess
	dbParam.InstIDField = common.BKProcIDField
	inst.base.saveSynchronize(ctx, dbParam)
	return nil
}

func (inst *instance) saveSynchronizePlatInstance(ctx core.ContextParams) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	dbParam.tableName = common.BKTableNameBasePlat
	dbParam.InstIDField = common.BKCloudIDField
	inst.base.saveSynchronize(ctx, dbParam)
	return nil
}

func (inst *instance) saveSynchronizeHostInstance(ctx core.ContextParams) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	dbParam.tableName = common.BKTableNameBaseHost
	dbParam.InstIDField = common.BKHostIDField
	inst.base.saveSynchronize(ctx, dbParam)
	return nil
}

func (inst *instance) saveSynchronizeObjectInstance(ctx core.ContextParams) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	dbParam.tableName = common.BKTableNameBaseInst
	dbParam.InstIDField = common.BKInstIDField
	inst.base.saveSynchronize(ctx, dbParam)
	return nil
}

func (inst *instance) getErrorStringArr(ctx core.ContextParams) ([]metadata.ExceptionResult, errors.CCError) {

	return inst.base.GetErrorStringArr(ctx)

}

func (inst *instance) preSynchronizeFilterBefore(ctx core.ContextParams) errors.CCError {
	return nil
}
func (inst *instance) preSynchronizeFilterEnd(ctx core.ContextParams) errors.CCError {
	return nil
}
