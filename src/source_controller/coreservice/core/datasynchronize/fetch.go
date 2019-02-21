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

type associationFetchDataInterface interface {
	Fetch(ctx core.ContextParams) ([]mapstr.MapStr, uint64, errors.CCError)
}

type associationFetchData struct {
	dataClassify string
	dataType     metadata.SynchronizeOperateDataType
	dbProxy      dal.RDB
	start        uint64
	limit        uint64
}

func NewSynchronizeFetchAdapter(fetch *metadata.SynchronizeFetchInfoParameter, dbProxy dal.RDB) associationFetchDataInterface {

	return &associationFetchData{
		dataClassify: fetch.DataClassify,
		dataType:     fetch.DataType,
		dbProxy:      dbProxy,
		start:        fetch.Start,
		limit:        fetch.Limit,
	}
}

func (a *associationFetchData) Fetch(ctx core.ContextParams) ([]mapstr.MapStr, uint64, errors.CCError) {
	switch a.dataType {
	case metadata.SynchronizeOperateDataTypeAssociation:
		return a.fetchAssociation(ctx)
	case metadata.SynchronizeOperateDataTypeModel:
		return a.fetchModel(ctx)
	}
	return nil, 0, nil
}

func (a *associationFetchData) fetchModel(ctx core.ContextParams) ([]mapstr.MapStr, uint64, errors.CCError) {
	switch a.dataClassify {
	case common.SynchronizeModelTypeBase:
		return a.findModel(ctx, common.BKTableNameObjDes)
	case common.SynchronizeModelTypeAttribute:
		return a.findModel(ctx, common.BKTableNameObjAttDes)
	case common.SynchronizeModelTypeAttributeGroup:
		return a.findModel(ctx, common.BKTableNamePropertyGroup)
	case common.SynchronizeModelTypeClassification:
		return a.findModel(ctx, common.BKTableNameObjClassifiction)
	}
	return nil, 0, nil
}

func (a *associationFetchData) findModel(ctx core.ContextParams, tableName string) ([]mapstr.MapStr, uint64, errors.CCError) {
	info := make([]mapstr.MapStr, 0)
	err := a.dbProxy.Table(tableName).Find(nil).Start(a.start).Limit(a.limit).All(ctx, info)
	if err != nil {
		blog.Errorf("findModel info error. error:%s,rid:%s", err.Error(), ctx.ReqID)
		return nil, 0, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}
	cnt, err := a.dbProxy.Table(tableName).Find(nil).Count(ctx)
	if err != nil {
		blog.Errorf("findModel count error. error:%s,rid:%s", err.Error(), ctx.ReqID)
		return nil, 0, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}
	return info, cnt, nil
}

func (a *associationFetchData) fetchAssociation(ctx core.ContextParams) ([]mapstr.MapStr, uint64, errors.CCError) {
	switch a.dataClassify {
	case common.SynchronizeAssociationTypeModelHost:
		return a.associationFetchHostModelData(ctx)
	}
	return nil, 0, nil
}

func (a *associationFetchData) associationFetchHostModelData(ctx core.ContextParams) ([]mapstr.MapStr, uint64, errors.CCError) {
	info := make([]mapstr.MapStr, 0)
	err := a.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(nil).Start(a.start).Limit(a.limit).All(ctx, info)
	if err != nil {
		blog.Errorf("associationFetchHostModelData info error. error:%s,rid:%s", err.Error(), ctx.ReqID)
		return nil, 0, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}
	cnt, err := a.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(nil).Count(ctx)
	if err != nil {
		blog.Errorf("associationFetchHostModelData count error. error:%s,rid:%s", err.Error(), ctx.ReqID)
		return nil, 0, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}

	return info, cnt, nil

}
