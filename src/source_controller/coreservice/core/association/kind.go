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

func (m *associationKind) CreateAssociationKind(ctx core.ContextParams, inputParam metadata.CreateAssociationKind) (*metadata.CreateOneDataResult, error) {
	_, exists, err := m.isExists(ctx, inputParam.Data.AssociationKindID)
	if nil != err {
		blog.Errorf("check association kind is exist error (%#v), rid: %s", err, ctx.ReqID)
		return nil, err
	}
	if exists {
		blog.Errorf("association kind (%#v)is duplicated, rid: %s", inputParam.Data, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommDuplicateItem, inputParam.Data.AssociationKindID)
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
			OriginIndex: int64(itemIdx),
			ID:          id,
		})

	}

	return dataResult, nil
}

func (m *associationKind) SetAssociationKind(ctx core.ContextParams, inputParam metadata.SetAssociationKind) (*metadata.SetDataResult, error) {
	origin, exists, err := m.isExists(ctx, inputParam.Data.AssociationKindID)
	if nil != err {
		blog.Errorf("check association kind is exist error (%#v), rid: %s", err, ctx.ReqID)
		return nil, err
	}
	dataResult := &metadata.SetDataResult{}
	if exists {
		cond := mongo.NewCondition()
		data := mapstr.NewFromStruct(inputParam.Data, "field")
		data.Remove(common.BKIsPre)
		data.Remove(common.AssociationKindIDField)
		data.Remove(common.BKFieldID)
		if err := m.update(ctx, data, cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: origin.AssociationKindID}).ToMapStr()); nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        inputParam.Data,
				OriginIndex: 0,
			})
			return dataResult, nil
		}
		dataResult.Updated = append(dataResult.Updated, metadata.UpdatedDataResult{ID: uint64(origin.ID)})
		dataResult.UpdatedCount.Count++
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
	dataResult.CreatedCount.Count++
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
			data.Remove(common.BKFieldID)
			if err := m.update(ctx, data, cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: origin.AssociationKindID}).ToMapStr()); nil != err {
				dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
					Message:     err.Error(),
					Code:        int64(err.(errors.CCErrorCoder).GetCode()),
					Data:        item,
					OriginIndex: int64(itemIdx),
				})
				continue
			}
			dataResult.Updated = append(dataResult.Updated, metadata.UpdatedDataResult{ID: uint64(origin.ID), OriginIndex: int64(itemIdx)})
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
			ID:          id,
			OriginIndex: int64(itemIdx),
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
	return &metadata.UpdatedCount{Count: cnt}, nil
}

func (m *associationKind) DeleteAssociationKind(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	queryCond := metadata.QueryCondition{Condition: inputParam.Condition}
	origins, err := m.searchAssociationKind(ctx, queryCond)
	if nil != err {
		blog.Errorf("search association kind by condition error:%s, rid: %s", err.Error(), ctx.ReqID)
		return &metadata.DeletedCount{}, err
	}

	for _, origin := range origins {
		cond := mongo.NewCondition()
		cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: origin.AssociationKindID})
		origin, exist, err := m.associationModel.isExists(ctx, cond)
		if nil != err {
			blog.Errorf("get association kind apply error:%s, rid: %s", err.Error(), ctx.ReqID)
			//return &metadata.DeletedCount{}, err
		}
		if exist {
			blog.Errorf("the association kind [%#v] has been apply to model [%#v], rid: %s", inputParam.Condition, origin, ctx.ReqID)
			return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrorTopoAssKindHasApplyToObject)
		}
	}

	exist, err := m.isPreAssociationKind(ctx, inputParam)
	if nil != err {
		blog.Errorf("search pre association kind by condition error:%s, rid: %s", err.Error(), ctx.ReqID)
		return &metadata.DeletedCount{}, err
	}
	if exist {
		blog.Errorf(" pre association can not be delete [%#v], rid: %s", inputParam.Condition, ctx.ReqID)
		return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrorTopoPreAssKindCanNotBeDelete)
	}
	err = m.dbProxy.Table(common.BKTableNameAsstDes).Delete(ctx, inputParam.Condition)
	if nil != err {
		blog.Errorf("delete association kind by condition [%#v],error:%s, rid: %s", inputParam.Condition, err.Error(), ctx.ReqID)
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: uint64(len(origins))}, nil
}

func (m *associationKind) CascadeDeleteAssociationKind(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	condition := metadata.QueryCondition{Condition: inputParam.Condition}
	associationKindItems, err := m.searchAssociationKind(ctx, condition)
	if nil != err {
		blog.Errorf("search association kind by condition [%#v],error:%s, rid: %s", inputParam.Condition, err.Error(), ctx.ReqID)
		return &metadata.DeletedCount{}, err
	}

	for _, item := range associationKindItems {
		cond := mongo.NewCondition()
		cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: item.AssociationKindID})
		deleteModelAsstParam := metadata.DeleteOption{Condition: cond.ToMapStr()}
		if _, err := m.associationModel.CascadeDeleteModelAssociation(ctx, deleteModelAsstParam); nil != err {
			blog.Errorf("cascade delete association kind by condition [%#v],error:%s, rid: %s", deleteModelAsstParam, err.Error(), ctx.ReqID)
			return &metadata.DeletedCount{}, err
		}
	}

	err = m.dbProxy.Table(common.BKTableNameAsstDes).Delete(ctx, inputParam.Condition)
	if nil != err {
		blog.Errorf("delete association kind by condition [%#v],error:%s, rid: %s", inputParam.Condition, err.Error(), ctx.ReqID)
		return &metadata.DeletedCount{}, err
	}

	return &metadata.DeletedCount{Count: uint64(len(associationKindItems))}, nil
}

func (m *associationKind) SearchAssociationKind(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	associationKindItems, err := m.searchAssociationKind(ctx, inputParam)
	if nil != err {
		return &metadata.QueryResult{}, err
	}

	dataResult := &metadata.QueryResult{}
	dataResult.Count, err = m.countInstanceAssociation(ctx, inputParam.Condition)
	dataResult.Info = make([]mapstr.MapStr, 0)
	if nil != err {
		return &metadata.QueryResult{}, err
	}
	for _, item := range associationKindItems {
		dataResult.Info = append(dataResult.Info, mapstr.NewFromStruct(item, "field"))
	}

	return dataResult, nil
}
