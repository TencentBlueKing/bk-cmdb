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

package association

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type associationKind struct {
	dbProxy dal.RDB
	*associationModel
}

func (m *associationKind) isExists(ctx core.ContextParams, associationKindID string) (origin *metadata.AssociationKind, exists bool, err error) {
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: associationKindID})
	err = m.dbProxy.Table(common.BKTableNameAsstDes).Find(cond.ToMapStr()).One(ctx, origin)
	return origin, !m.dbProxy.IsNotFoundError(err), err
}

func (m *associationKind) update(ctx core.ContextParams, data mapstr.MapStr, cond mapstr.MapStr) error {

	return m.dbProxy.Table(common.BKTableNameAsstDes).Update(ctx, cond, data)
}

func (m *associationKind) searchAssociationKind(ctx core.ContextParams, cond mapstr.MapStr) ([]metadata.AssociationKind, error) {

	results := []metadata.AssociationKind{}
	err := m.dbProxy.Table(common.BKTableNameObjClassifiction).Find(cond).All(ctx, &results)

	return results, err
}

func (m *associationKind) isPrPreAssociationKind(ctx core.ContextParams, cond metadata.DeleteOption) (exists bool, err error) {

	innerCnt, err := m.dbProxy.Table(common.BKTableNameAsstDes).Find(cond).Count(ctx)
	exists = 0 != innerCnt
	return exists, err
}

func (m *associationKind) isApplyToObject(ctx core.ContextParams, cond metadata.DeleteOption) (cnt int64, exists bool, err error) {

	innerCnt, err := m.dbProxy.Table(common.BKTableNameAsstDes).Find(cond).Count(ctx)
	cnt = int64(innerCnt)
	exists = 0 != cnt
	return cnt, exists, err
}

func (m *associationKind) save(ctx core.ContextParams, associationKind metadata.AssociationKind) (id uint64, err error) {

	id, err = m.dbProxy.NextSequence(ctx, common.BKTableNameAsstDes)
	if err != nil {
		return id, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	associationKind.ID = int64(id)

	err = m.dbProxy.Table(common.BKTableNameAsstDes).Insert(ctx, associationKind)
	return id, err
}

func (m *associationKind) CreateAssociationKind(ctx core.ContextParams, inputParam metadata.CreateAssociationKind) (*metadata.CreateOneDataResult, error) {
	_, exists, err := m.isExists(ctx, inputParam.Data.AssociationKindID)
	if nil != err {
		return nil, err
	}
	if exists {
		blog.Errorf("association kind (%v)is duplicated", inputParam.Data)
		return nil, ctx.Error.Error(common.CCErrCommDuplicateItem)
	}

	id, err := m.save(ctx, inputParam.Data)
	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, err
}

func (m *associationKind) CreateManyAssociationKind(ctx core.ContextParams, inputParam metadata.CreateManyAssociationKind) (*metadata.CreateManyDataResult, error) {
	dataResult := &metadata.CreateManyDataResult{}
	for itemIdx, item := range inputParam.Datas {

		_, exists, err := m.isExists(ctx, item.AssociationKindID)
		if nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}

		if exists {
			dataResult.Repeated = append(dataResult.Repeated, metadata.RepeatedDataResult{OriginIndex: int64(itemIdx), Data: mapstr.NewFromStruct(item, "field")})
			continue
		}

		id, err := m.save(ctx, item)
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

func (m *associationKind) SetAssociationKind(ctx core.ContextParams, inputParam metadata.SetAssociationKind) (*metadata.SetDataResult, error) {
	origin, exists, err := m.isExists(ctx, inputParam.Data.AssociationKindID)
	if nil != err {
		return nil, err
	}

	dataResult := &metadata.SetDataResult{}

	if exists {

		cond := mongo.NewCondition()
		data := mapstr.NewFromStruct(inputParam.Data, "field")
		data.Remove(common.BKIsPre)
		data.Remove(common.AssociationKindIDField)
		if err := m.update(ctx, data, cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: origin.ID}).ToMapStr()); nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        inputParam.Data,
				OriginIndex: 0,
			})
			return dataResult, nil
		}
		dataResult.Updated = append(dataResult.Updated, metadata.UpdatedDataResult{ID: uint64(origin.ID)})
		return dataResult, err
	}

	id, err := m.save(ctx, inputParam.Data)
	if nil != err {
		dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
			Message:     err.Error(),
			Code:        int64(err.(errors.CCErrorCoder).GetCode()),
			Data:        origin,
			OriginIndex: 0,
		})
	}
	dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{ID: id})
	return dataResult, err
}

func (m *associationKind) SetManyAssociationKind(ctx core.ContextParams, inputParam metadata.SetManyAssociationKind) (*metadata.SetDataResult, error) {
	dataResult := &metadata.SetDataResult{}
	for itemIdx, item := range inputParam.Datas {

		origin, exists, err := m.isExists(ctx, item.AssociationKindID)
		if nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}

		if exists {

			cond := mongo.NewCondition()
			data := mapstr.NewFromStruct(item, "field")
			data.Remove(common.BKIsPre)
			data.Remove(common.AssociationKindIDField)
			if err := m.update(ctx, data, cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: origin.ID}).ToMapStr()); nil != err {
				dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
					Message:     err.Error(),
					Code:        int64(err.(errors.CCErrorCoder).GetCode()),
					Data:        item,
					OriginIndex: int64(itemIdx),
				})
				continue
			}

			dataResult.UpdatedCount.Count++
			continue
		}

		id, err := m.save(ctx, item)
		if nil != err {

			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}

		dataResult.CreatedCount.Count++
		dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{
			ID: id,
		})

	}

	return dataResult, nil
}
func (m *associationKind) UpdateAssociationKind(ctx core.ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {
	cnt, err := m.dbProxy.Table(common.BKTableNameAsstDes).Find(inputParam.Condition).Count(ctx)
	if nil != err {
		return &metadata.UpdatedCount{}, err
	}
	if err := m.update(ctx, inputParam.Data, inputParam.Condition); nil != err {
		return &metadata.UpdatedCount{}, err
	}
	return &metadata.UpdatedCount{Count: int64(cnt)}, nil
}

func (m *associationKind) DeleteAssociationKind(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	cnt, exists, err := m.isApplyToObject(ctx, inputParam)
	if nil != err {
		return &metadata.DeletedCount{}, err
	}
	if exists {
		return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrorTopoAssKindHasApplyToObject)
	}

	exists, err = m.isPrPreAssociationKind(ctx, inputParam)
	if nil != err {
		return &metadata.DeletedCount{}, err
	}
	if exists {
		return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrorTopoPreAssKindCanNotBeDelete)
	}

	m.dbProxy.Table(common.BKTableNameAsstDes).Delete(ctx, inputParam.Condition)
	return &metadata.DeletedCount{Count: cnt}, nil
}

func (m *associationKind) CascadeDeleteAssociationKind(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	associationKindItems, err := m.searchAssociationKind(ctx, inputParam.Condition)
	if nil != err {
		return &metadata.DeletedCount{}, err
	}

	for _, item := range associationKindItems {
		cond := mongo.NewCondition()
		cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: item.AssociationKindID})
		if _, err := m.associationModel.CascadeDeleteModelAssociation(ctx, inputParam); nil != err {
			return &metadata.DeletedCount{}, err
		}
	}

	return &metadata.DeletedCount{Count: int64(len(associationKindItems))}, nil
}

func (m *associationKind) SearchAssociationKind(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	associationKindItems, err := m.searchAssociationKind(ctx, inputParam.Condition)
	if nil != err {
		return &metadata.QueryResult{}, err
	}

	dataResult := &metadata.QueryResult{}
	dataResult.Count = int64(len(associationKindItems))
	for item := range associationKindItems {
		dataResult.Info = append(dataResult.Info, mapstr.NewFromStruct(item, "field"))
	}

	return dataResult, nil
}
