/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
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
func New(dbProxy dal.RDB) core.InstanceOperation {
	return &instanceManager{
		dbProxy: dbProxy,
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
	origins, _, err := m.getInsts(ctx, objID, inputParam.Condition)
	if nil != err {
		return nil, err
	}

	for _, origin := range origins {
		instIDI := origin[instIDFieldName]
		instID, _ := util.GetInt64ByInterface(instIDI)
		err := m.validUpdateInstanceData(ctx, objID, inputParam.Data, uint64(instID))
		if nil != err {
			return nil, err
		}
	}

	updateCond, err := mongo.NewConditionFromMapStr(inputParam.Condition)
	if nil != err {
		return &metadata.UpdatedCount{}, err
	}

	cnt, err := m.update(ctx, objID, inputParam.Data, updateCond)
	return &metadata.UpdatedCount{Count: cnt}, err
}

func (m *instanceManager) SearchModelInstance(ctx core.ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	instItems, err := m.searchInstance(ctx, objID, inputParam)
	if nil != err {
		return &metadata.QueryResult{}, err
	}

	dataResult := &metadata.QueryResult{}
	dataResult.Count, err = m.countInstance(ctx, objID, inputParam.Condition)
	if nil != err {
		return &metadata.QueryResult{}, err
	}
	dataResult.Info = instItems

	return dataResult, nil
}

func (m *instanceManager) DeleteModelInstance(ctx core.ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	tableName := common.GetInstTableName(objID)
	instIDFieldName := common.GetInstIDField(objID)
	origins, _, err := m.getInsts(ctx, objID, inputParam.Condition)
	if nil != err {
		return nil, err
	}

	for _, origin := range origins {
		instID, err := util.GetInt64ByInterface(origin[instIDFieldName])
		if nil != err {
			return &metadata.DeletedCount{}, err
		}
		exists, err := m.dependent.IsInstAsstExist(ctx, objID, uint64(instID))
		if nil != err {
			return &metadata.DeletedCount{}, err
		}
		if exists {
			return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrorInstHasAsst)
		}
	}

	m.dbProxy.Table(tableName).Delete(ctx, inputParam.Condition)
	return &metadata.DeletedCount{Count: uint64(len(origins))}, nil
}

func (m *instanceManager) CascadeDeleteModelInstance(ctx core.ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	tableName := common.GetInstTableName(objID)
	instIDFieldName := common.GetInstIDField(objID)
	origins, _, err := m.getInsts(ctx, objID, inputParam.Condition)
	if nil != err {
		return nil, err
	}

	for _, origin := range origins {
		instID, err := util.GetInt64ByInterface(origin[instIDFieldName])
		if nil != err {
			return &metadata.DeletedCount{}, err
		}
		_, err = m.dependent.DeleteInstAsst(ctx, objID, uint64(instID))
		if nil != err {
			return &metadata.DeletedCount{}, err
		}
	}

	m.dbProxy.Table(tableName).Delete(ctx, inputParam.Condition)
	return &metadata.DeletedCount{Count: uint64(len(origins))}, nil
}
