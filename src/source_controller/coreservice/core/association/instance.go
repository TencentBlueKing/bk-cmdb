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

type associationInstance struct {
	dbProxy dal.RDB
}

func (m *associationInstance) isExists(ctx core.ContextParams, instID, asstInstID int64, objAsstID string) (origin *metadata.InstAsst, exists bool, err error) {
	cond := mongo.NewCondition()
	cond.Element(
		&mongo.Eq{Key: common.BKInstIDField, Val: instID},
		&mongo.Eq{Key: common.BKAsstInstIDField, Val: asstInstID},
		&mongo.Eq{Key: common.AssociationObjAsstIDField, Val: objAsstID})
	err = m.dbProxy.Table(common.BKTableNameInstAsst).Find(cond.ToMapStr()).One(ctx, origin)
	return origin, !m.dbProxy.IsNotFoundError(err), err
}

func (m *associationInstance) instCount(ctx core.ContextParams, cond mapstr.MapStr) (cnt int64, err error) {
	innerCnt, err := m.dbProxy.Table(common.BKTableNameInstAsst).Find(cond).Count(ctx)
	cnt = int64(innerCnt)
	return cnt, err
}

func (m *associationInstance) searchInstanceAssociation(ctx core.ContextParams, cond mapstr.MapStr) ([]metadata.InstAsst, error) {

	results := []metadata.InstAsst{}
	err := m.dbProxy.Table(common.BKTableNameInstAsst).Find(cond).All(ctx, &results)

	return results, err
}

func (m *associationInstance) save(ctx core.ContextParams, asstInst metadata.InstAsst) (id uint64, err error) {

	id, err = m.dbProxy.NextSequence(ctx, common.BKTableNameInstAsst)
	if err != nil {
		return id, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	asstInst.ID = int64(id)

	err = m.dbProxy.Table(common.BKTableNameInstAsst).Insert(ctx, asstInst)
	return id, err
}

func (m *associationInstance) CreateOneInstanceAssociation(ctx core.ContextParams, inputParam metadata.CreateOneInstanceAssociation) (*metadata.CreateOneDataResult, error) {
	_, exists, err := m.isExists(ctx, inputParam.Data.InstID, inputParam.Data.AsstInstID, inputParam.Data.ObjectAsstID)
	if nil != err {
		return nil, err
	}
	if exists {
		blog.Errorf("association instance (%v)is duplicated", inputParam.Data)
		return nil, ctx.Error.Error(common.CCErrCommDuplicateItem)
	}

	id, err := m.save(ctx, inputParam.Data)
	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, err
}

func (m *associationInstance) SetOneInstanceAssociation(ctx core.ContextParams, inputParam metadata.SetOneInstanceAssociation) (*metadata.SetDataResult, error) {
	origin, exists, err := m.isExists(ctx, inputParam.Data.InstID, inputParam.Data.AsstInstID, inputParam.Data.ObjectAsstID)
	if nil != err {
		return nil, err
	}

	dataResult := &metadata.SetDataResult{}

	if exists {
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

func (m *associationInstance) CreateManyInstanceAssociation(ctx core.ContextParams, inputParam metadata.CreateManyInstanceAssociation) (*metadata.CreateManyDataResult, error) {
	dataResult := &metadata.CreateManyDataResult{}
	for itemIdx, item := range inputParam.Datas {

		_, exists, err := m.isExists(ctx, item.InstID, item.AsstInstID, item.ObjectAsstID)
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
func (m *associationInstance) SetManyInstanceAssociation(ctx core.ContextParams, inputParam metadata.SetManyInstanceAssociation) (*metadata.SetDataResult, error) {
	dataResult := &metadata.SetDataResult{}
	for itemIdx, item := range inputParam.Datas {

		_, exists, err := m.isExists(ctx, item.InstID, item.AsstInstID, item.ObjectAsstID)
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

func (m *associationInstance) SearchInstanceAssociation(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	instAsstItems, err := m.searchInstanceAssociation(ctx, inputParam.Condition)
	if nil != err {
		return &metadata.QueryResult{}, err
	}

	dataResult := &metadata.QueryResult{}
	dataResult.Count = int64(len(instAsstItems))
	for item := range instAsstItems {
		dataResult.Info = append(dataResult.Info, mapstr.NewFromStruct(item, "field"))
	}

	return dataResult, nil
}
func (m *associationInstance) DeleteInstanceAssociation(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	cnt, err := m.instCount(ctx, inputParam.Condition)
	if nil != err {
		return &metadata.DeletedCount{}, err
	}

	m.dbProxy.Table(common.BKTableNameInstAsst).Delete(ctx, inputParam.Condition)
	return &metadata.DeletedCount{Count: cnt}, nil
}
