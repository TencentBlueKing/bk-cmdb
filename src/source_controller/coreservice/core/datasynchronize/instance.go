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

type instance struct {
	base         *synchronizeAdapter
	dataType     metadata.SynchronizeOperateDataType
	DataClassify string
}

func NewSynchronizeInstanceAdapter(s *metadata.SynchronizeParameter) dataTypeInterface {

	return &instance{
		base:     newSynchronizeAdapter(s),
		dataType: s.OperateDataType,
		// instance data classify
		DataClassify: s.DataClassify,
	}
}

func (inst *instance) PreSynchronizeFilter(kit *rest.Kit) errors.CCError {
	err := inst.preSynchronizeFilterBefore(kit)
	if err != nil {
		return err
	}
	return inst.base.PreSynchronizeFilter(kit)
}

func (inst *instance) GetErrorStringArr(kit *rest.Kit) ([]metadata.ExceptionResult, errors.CCError) {

	if len(inst.base.errorArray) == 0 {
		return nil, nil
	}

	return inst.base.GetErrorStringArr(kit)

}

func (inst *instance) SaveSynchronize(kit *rest.Kit) errors.CCError {
	// Each model is written separately for subsequent expansion,
	// each model may be processed differently.

	switch inst.DataClassify {
	case common.BKInnerObjIDApp:
		return inst.saveSynchronizeAppInstance(kit)
	case common.BKInnerObjIDSet:
		return inst.saveSynchronizeSetInstance(kit)
	case common.BKInnerObjIDModule:
		return inst.saveSynchronizeModuleInstance(kit)
	case common.BKInnerObjIDProc:
		return inst.saveSynchronizeModuleInstance(kit)
	case common.BKInnerObjIDPlat:
		return inst.saveSynchronizePlatInstance(kit)
	case common.BKInnerObjIDHost:
		return inst.saveSynchronizeHostInstance(kit)
	default:
		return inst.saveSynchronizeObjectInstance(kit)
	}
}

func (inst *instance) saveSynchronizeAppInstance(kit *rest.Kit) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	dbParam.tableName = common.BKTableNameBaseApp
	dbParam.InstIDField = common.BKAppIDField
	inst.base.saveSynchronize(kit, dbParam)
	return nil
}

func (inst *instance) saveSynchronizeSetInstance(kit *rest.Kit) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	dbParam.tableName = common.BKTableNameBaseSet
	dbParam.InstIDField = common.BKSetIDField
	inst.base.saveSynchronize(kit, dbParam)
	return nil
}

func (inst *instance) saveSynchronizeModuleInstance(kit *rest.Kit) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	dbParam.tableName = common.BKTableNameBaseModule
	dbParam.InstIDField = common.BKModuleIDField
	inst.base.saveSynchronize(kit, dbParam)
	return nil
}

func (inst *instance) saveSynchronizeProcessInstance(kit *rest.Kit) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	dbParam.tableName = common.BKTableNameBaseProcess
	dbParam.InstIDField = common.BKProcIDField
	inst.base.saveSynchronize(kit, dbParam)
	return nil
}

func (inst *instance) saveSynchronizePlatInstance(kit *rest.Kit) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	dbParam.tableName = common.BKTableNameBasePlat
	dbParam.InstIDField = common.BKCloudIDField
	inst.base.saveSynchronize(kit, dbParam)
	return nil
}

func (inst *instance) saveSynchronizeHostInstance(kit *rest.Kit) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	dbParam.tableName = common.BKTableNameBaseHost
	dbParam.InstIDField = common.BKHostIDField
	for _, info := range inst.base.syncData.InfoArray {
		info.Info = metadata.ConvertHostSpecialStringToArray(info.Info)
	}
	inst.base.saveSynchronize(kit, dbParam)
	return nil
}

func (inst *instance) saveSynchronizeObjectInstance(kit *rest.Kit) errors.CCError {
	var dbParam synchronizeAdapterDBParameter
	dbParam.tableName = common.BKTableNameBaseInst
	dbParam.InstIDField = common.BKInstIDField
	inst.base.saveSynchronize(kit, dbParam)
	return nil
}

func (inst *instance) getErrorStringArr(kit *rest.Kit) ([]metadata.ExceptionResult, errors.CCError) {

	return inst.base.GetErrorStringArr(kit)

}

func (inst *instance) preSynchronizeFilterBefore(kit *rest.Kit) errors.CCError {
	return nil
}
func (inst *instance) preSynchronizeFilterEnd(kit *rest.Kit) errors.CCError {
	return nil
}
