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

package instances

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

var _ core.InstanceOperation = (*instanceManager)(nil)

type instanceManager struct {
	dbProxy   dal.RDB
	dependent OperationDependences
	validator validator
}

// New create a new instance manager instance
func New(dbProxy dal.RDB, dependent OperationDependences) core.InstanceOperation {
	return &instanceManager{
		dbProxy:   dbProxy,
		dependent: dependent,
	}
}

func (m *instanceManager) instCnt(ctx core.ContextParams, objID string, cond mapstr.MapStr) (cnt uint64, exists bool, err error) {
	tableName := common.GetInstTableName(objID)
	cnt, err = m.dbProxy.Table(tableName).Find(cond).Count(ctx)
	exists = 0 != cnt
	return cnt, exists, err
}

func (m *instanceManager) CreateModelInstance(ctx core.ContextParams, objID string, inputParam metadata.CreateModelInstance) (*metadata.CreateOneDataResult, error) {
	err := m.validCreateInstanceData(ctx, objID, inputParam.Data)
	if nil != err {
		blog.Errorf("create inst valid error: %v", err)
		return nil, err
	}
	id, err := m.save(ctx, objID, inputParam.Data)
	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, err
}

func (m *instanceManager) CreateManyModelInstance(ctx core.ContextParams, objID string, inputParam metadata.CreateManyModelInstance) (*metadata.CreateManyDataResult, error) {
	dataResult := &metadata.CreateManyDataResult{}
	for itemIdx, item := range inputParam.Datas {
		item.Set(common.BKOwnerIDField, ctx.SupplierAccount)
		err := m.validCreateInstanceData(ctx, objID, item)
		if nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}
		item.Set(common.BKOwnerIDField, ctx.SupplierAccount)
		id, err := m.save(ctx, objID, item)
		if nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}

		dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{
			ID: id,
		})

	}

	return dataResult, nil
}

func (m *instanceManager) UpdateModelInstance(ctx core.ContextParams, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {
	instIDFieldName := common.GetInstIDField(objID)
	inputParam.Condition.Set(common.BKOwnerIDField, ctx.SupplierAccount)
	origins, _, err := m.getInsts(ctx, objID, inputParam.Condition)
	if nil != err {
		blog.Errorf("update module instance get inst error :%v ", err)
		return nil, err
	}

	var instMedataData metadata.Metadata
	instMedataData.Label = make(metadata.Label)
	for key, val := range inputParam.Condition {
		if metadata.BKMetadata == key {
			bizID := metadata.GetBusinessIDFromMeta(val)
			if "" != bizID {
				instMedataData.Label.Set(metadata.LabelBusinessID, metadata.GetBusinessIDFromMeta(val))
			}
			continue
		}
	}

	for _, origin := range origins {
		instIDI := origin[instIDFieldName]
		instID, _ := util.GetInt64ByInterface(instIDI)
		err := m.validUpdateInstanceData(ctx, objID, inputParam.Data, instMedataData, uint64(instID))
		if nil != err {
			blog.Errorf("update module instance validate error :%v ", err)
			return nil, err
		}
	}

	if nil != err {
		blog.Errorf("update module instance validate error :%v ", err)
		return &metadata.UpdatedCount{}, err
	}
	cnt, err := m.update(ctx, objID, inputParam.Data, inputParam.Condition)
	return &metadata.UpdatedCount{Count: cnt}, err
}

func (m *instanceManager) SearchModelInstance(ctx core.ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	condition, err := mongo.NewConditionFromMapStr(inputParam.Condition)
	if nil != err {
		blog.Errorf("parse conditon  error %v, [%v]", err)
		return &metadata.QueryResult{}, err
	}
	ownerIDArr := []string{ctx.SupplierAccount, common.BKDefaultOwnerID}
	condition.Element(&mongo.In{Key: common.BKOwnerIDField, Val: ownerIDArr})
	inputParam.Condition = condition.ToMapStr()

	instItems, err := m.searchInstance(ctx, objID, inputParam)
	if nil != err {
		blog.Errorf("search instance error [%v]", err)
		return &metadata.QueryResult{}, err
	}

	dataResult := &metadata.QueryResult{}
	dataResult.Count, err = m.countInstance(ctx, objID, inputParam.Condition)
	if nil != err {
		blog.Errorf("count instance error [%v]", err)
		return &metadata.QueryResult{}, err
	}
	dataResult.Info = instItems

	return dataResult, nil
}

func (m *instanceManager) DeleteModelInstance(ctx core.ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	tableName := common.GetInstTableName(objID)
	instIDFieldName := common.GetInstIDField(objID)
	inputParam.Condition.Set(common.BKOwnerIDField, ctx.SupplierAccount)
	origins, _, err := m.getInsts(ctx, objID, inputParam.Condition)
	if nil != err {
		return &metadata.DeletedCount{}, err
	}

	for _, origin := range origins {
		instID, err := util.GetInt64ByInterface(origin[instIDFieldName])
		if nil != err {
			return nil, err
		}
		exists, err := m.dependent.IsInstAsstExist(ctx, objID, uint64(instID))
		if nil != err {
			return nil, err
		}
		if exists {
			return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrorInstHasAsst)
		}
	}
	err = m.dbProxy.Table(tableName).Delete(ctx, inputParam.Condition)
	if nil != err {
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: uint64(len(origins))}, nil
}

func (m *instanceManager) CascadeDeleteModelInstance(ctx core.ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	tableName := common.GetInstTableName(objID)
	instIDFieldName := common.GetInstIDField(objID)
	origins, _, err := m.getInsts(ctx, objID, inputParam.Condition)
	blog.Errorf("cascade delete model instance get inst error:%v", origins)
	if nil != err {
		blog.Errorf("cascade delete model instance get inst error:%v", err)
		return &metadata.DeletedCount{}, err
	}

	for _, origin := range origins {
		instID, err := util.GetInt64ByInterface(origin[instIDFieldName])
		if nil != err {
			return &metadata.DeletedCount{}, err
		}
		err = m.dependent.DeleteInstAsst(ctx, objID, uint64(instID))
		if nil != err {
			return &metadata.DeletedCount{}, err
		}
	}
	inputParam.Condition.Set(common.BKOwnerIDField, ctx.SupplierAccount)
	err = m.dbProxy.Table(tableName).Delete(ctx, inputParam.Condition)
	if nil != err {
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: uint64(len(origins))}, nil
}
