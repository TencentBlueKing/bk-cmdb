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
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type associationFindDataInterface interface {
	Find(ctx core.ContextParams) ([]mapstr.MapStr, uint64, errors.CCError)
}

type associationFindData struct {
	dataClassify string
	dataType     metadata.SynchronizeOperateDataType
	dbProxy      dal.RDB
	start        uint64
	limit        uint64
	condition    mapstr.MapStr
}

func NewSynchronizeFindAdapter(input *metadata.SynchronizeFindInfoParameter, dbProxy dal.RDB) associationFindDataInterface {

	return &associationFindData{
		dataClassify: input.DataClassify,
		dataType:     input.DataType,
		dbProxy:      dbProxy,
		start:        input.Start,
		limit:        input.Limit,
		condition:    input.Condition,
	}
}

func (a *associationFindData) Find(ctx core.ContextParams) ([]mapstr.MapStr, uint64, errors.CCError) {
	switch a.dataType {
	case metadata.SynchronizeOperateDataTypeAssociation:
		return a.findAssociation(ctx)
	case metadata.SynchronizeOperateDataTypeModel:
		return a.findModel(ctx)
	}
	return nil, 0, nil
}

func (a *associationFindData) findModel(ctx core.ContextParams) ([]mapstr.MapStr, uint64, errors.CCError) {

	switch a.dataClassify {
	case common.SynchronizeModelTypeBase:
		return a.dbQueryModel(ctx, common.BKTableNameObjDes)
	case common.SynchronizeModelTypeAttribute:
		return a.dbQueryModel(ctx, common.BKTableNameObjAttDes)
	case common.SynchronizeModelTypeAttributeGroup:
		return a.dbQueryModel(ctx, common.BKTableNamePropertyGroup)
	case common.SynchronizeModelTypeClassification:
		return a.dbQueryModel(ctx, common.BKTableNameObjClassifiction)
	}
	return nil, 0, nil
}

func (a *associationFindData) dbQueryModel(ctx core.ContextParams, tableName string) ([]mapstr.MapStr, uint64, errors.CCError) {
	info := make([]mapstr.MapStr, 0)
	err := a.dbProxy.Table(tableName).Find(a.condition).Start(a.start).Limit(a.limit).All(ctx, &info)
	if err != nil {
		blog.Errorf("dbQueryModel info error. error:%s,rid:%s", err.Error(), ctx.ReqID)
		return nil, 0, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}
	cnt, err := a.dbProxy.Table(tableName).Find(nil).Count(ctx)
	if err != nil {
		blog.Errorf("dbQueryModel count error. error:%s,rid:%s", err.Error(), ctx.ReqID)
		return nil, 0, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}
	return info, cnt, nil
}

func (a *associationFindData) findAssociation(ctx core.ContextParams) ([]mapstr.MapStr, uint64, errors.CCError) {
	switch a.dataClassify {
	case common.SynchronizeAssociationTypeModelHost:
		return a.dbQueryAssociation(ctx)
	}
	return nil, 0, nil
}

func (a *associationFindData) dbQueryAssociation(ctx core.ContextParams) ([]mapstr.MapStr, uint64, errors.CCError) {
	info := make([]mapstr.MapStr, 0)
	err := a.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(a.condition).Start(a.start).Limit(a.limit).All(ctx, &info)
	if err != nil {
		blog.Errorf("dbQueryAssociation info error. error:%s,rid:%s", err.Error(), ctx.ReqID)
		return nil, 0, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}
	cnt, err := a.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(nil).Count(ctx)
	if err != nil {
		blog.Errorf("dbQueryAssociation count error. error:%s,rid:%s", err.Error(), ctx.ReqID)
		return nil, 0, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}

	return info, cnt, nil

}
