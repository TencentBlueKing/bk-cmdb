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

package datasynchronize

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
)

type setIdentifierFlag struct {
	params *metadata.SetIdenifierFlag

	// db collection name
	tableName string
	// instance id field name
	instIDField string
}

// NewSetIdentifierFlag net set identifier handler
func NewSetIdentifierFlag(params *metadata.SetIdenifierFlag) *setIdentifierFlag {
	return &setIdentifierFlag{
		params: params,
	}
}

// Run execute set identifier flag logic
func (s *setIdentifierFlag) Run(kit *rest.Kit) errors.CCErrorCoder {
	blog.V(5).Infof("SetIdentifierFlag api Run. input:%#v, rid:%s", s.params, kit.Rid)
	if s.params.DataType != metadata.SynchronizeOperateDataTypeInstance {
		// not support DataType
		blog.Errorf(" not support data_type. input:%#v, rid:%s", s.params, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsValueInvalidError, "data_type", s.params.DataType)
	}
	err := s.dbCollectionInfo(kit)
	if err != nil {
		blog.Errorf("dbCollectionInfo handle db identifier id field error. err:%s, input:%#v, rid:%s", err.Error(), s.params, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsValueInvalidError, "data_type", s.params.DataType)

	}
	switch s.params.OperateType {
	case metadata.SynchronizeOperateTypeAdd:
		return s.addFlag(kit)
	case metadata.SynchronizeOperateTypeRepalce:
		return s.replaceFlag(kit)
	case metadata.SynchronizeOperateTypeDelete:
		return s.deleteFlag(kit)
	default:
		// not upport OperateType
		blog.Errorf(" not support op_type. input:%#v, rid:%s", s.params, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsValueInvalidError, "op_type", s.params.OperateType)
	}
}

// dbCollectionInfo get data in db table name and data identifier id field
func (s *setIdentifierFlag) dbCollectionInfo(kit *rest.Kit) errors.CCErrorCoder {
	switch s.params.DataType {
	case metadata.SynchronizeOperateDataTypeInstance:
		s.tableName = common.GetInstTableName(s.params.DataClassify)
		s.instIDField = common.GetInstIDField(s.params.DataClassify)

	default:
		// not suuport DataType
		blog.Errorf(" not support data_type. input:%#v, rid:%s", s.params, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsValueInvalidError, "data_type", s.params.DataType)

	}

	return nil
}

// addFlag 添加新的需要同步数据cmdb身份标识
func (s *setIdentifierFlag) addFlag(kit *rest.Kit) errors.CCErrorCoder {
	conds := condition.CreateCondition()
	conds.Field(s.instIDField).In(s.params.IdentifierID)
	condMap := conds.ToMapStr()

	// 如果metadata = null。初始化为{}, 避免数据库报出错
	fixMetadataCondMap := condMap.Clone()
	fixMetadataCondMap[common.MetadataField] = nil

	err := mongodb.Client().Table(s.tableName).Update(kit.Ctx, fixMetadataCondMap, mapstr.MapStr{common.MetadataField: mapstr.New()})
	if err != nil {
		blog.ErrorJSON("addFlag db update error. err:%s,  cond:%s, rid:%s", err.Error(), fixMetadataCondMap, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}
	// 如果同步数据cmdb身份标识已经存在，不做任何操作。否则新加同步标志
	data := mapstr.MapStr{
		util.BuildMongoSyncItemField(common.MetaDataSynchronizeIdentifierField): s.params.Flag,
	}
	err = mongodb.Client().Table(s.tableName).UpdateMultiModel(kit.Ctx, condMap, types.ModeUpdate{Op: types.UpdateOpAddToSet, Doc: data})
	if err != nil {
		blog.ErrorJSON("addFlag db update error. err:%s,  cond:%s, data:%s, rid:%s", err.Error(), condMap, data, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

// repalceFlag 使用当前cmdb系统同步身份标识。该操作会删除之前的所有的cmdb系统同步身份标识
func (s *setIdentifierFlag) replaceFlag(kit *rest.Kit) errors.CCErrorCoder {
	conds := condition.CreateCondition()
	conds.Field(s.instIDField).In(s.params.IdentifierID)
	condMap := conds.ToMapStr()
	data := mapstr.MapStr{
		util.BuildMongoSyncItemField(common.MetaDataSynchronizeIdentifierField): []string{s.params.Flag},
	}
	err := mongodb.Client().Table(s.tableName).Update(kit.Ctx, condMap, data)
	if err != nil {
		blog.ErrorJSON("addFlag db update error. err:%s,  cond:%s, data:%s, rid:%s", err.Error(), condMap, data, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

// deleteFlag
func (s *setIdentifierFlag) deleteFlag(kit *rest.Kit) errors.CCErrorCoder {
	conds := condition.CreateCondition()
	conds.Field(s.instIDField).In(s.params.IdentifierID)
	conds.Field(common.MetadataField).NotEq(nil)
	condMap := conds.ToMapStr()
	data := mapstr.MapStr{
		util.BuildMongoSyncItemField(common.MetaDataSynchronizeIdentifierField): s.params.Flag,
	}
	err := mongodb.Client().Table(s.tableName).UpdateMultiModel(kit.Ctx, condMap, types.ModeUpdate{Op: types.UpdateOpPull, Doc: data})
	if err != nil {
		blog.ErrorJSON("addFlag db update error. err:%s,  cond:%s, data:%s, rid:%s", err.Error(), condMap, data, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}
	return nil
}
