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

type associationInstance struct {
	dbProxy dal.RDB
	*associationKind
	*associationModel
	dependent OperationDependences
}

func (m *associationInstance) isExists(ctx core.ContextParams, instID, asstInstID int64, objAsstID string, meta metadata.Metadata) (origin *metadata.InstAsst, exists bool, err error) {
	cond := mongo.NewCondition()
	origin = &metadata.InstAsst{}
	cond.Element(
		&mongo.Eq{Key: common.BKInstIDField, Val: instID},
		&mongo.Eq{Key: common.BKAsstInstIDField, Val: asstInstID},
		&mongo.Eq{Key: common.AssociationObjAsstIDField, Val: objAsstID})
	isExsit, bizID := meta.Label.Get(common.BKAppIDField)
	if isExsit {
		_, metaCond := cond.Embed(metadata.BKMetadata)
		_, lableCond := metaCond.Embed(metadata.BKLabel)
		lableCond.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bizID})
	}

	err = m.dbProxy.Table(common.BKTableNameInstAsst).Find(cond.ToMapStr()).One(ctx, origin)
	if m.dbProxy.IsNotFoundError(err) {
		return origin, !m.dbProxy.IsNotFoundError(err), nil
	}
	return origin, !m.dbProxy.IsNotFoundError(err), err
}

func (m *associationInstance) instCount(ctx core.ContextParams, cond mapstr.MapStr) (cnt uint64, err error) {
	innerCnt, err := m.dbProxy.Table(common.BKTableNameInstAsst).Find(cond).Count(ctx)
	return innerCnt, err
}

func (m *associationInstance) searchInstanceAssociation(ctx core.ContextParams, inputParam metadata.QueryCondition) (results []metadata.InstAsst, err error) {
	results = []metadata.InstAsst{}
	instHandler := m.dbProxy.Table(common.BKTableNameInstAsst).Find(inputParam.Condition).Fields(inputParam.Fields...)
	for _, sort := range inputParam.SortArr {
		fileld := sort.Field
		if sort.IsDsc {
			fileld = "-" + fileld
		}
		instHandler = instHandler.Sort(fileld)
	}
	err = instHandler.Start(uint64(inputParam.Limit.Offset)).Limit(uint64(inputParam.Limit.Limit)).All(ctx, &results)
	return results, err
}

func (m *associationInstance) countInstanceAssociation(ctx core.ContextParams, cond mapstr.MapStr) (count uint64, err error) {
	count, err = m.dbProxy.Table(common.BKTableNameInstAsst).Find(cond).Count(ctx)

	return count, err
}

func (m *associationInstance) save(ctx core.ContextParams, asstInst metadata.InstAsst) (id uint64, err error) {

	id, err = m.dbProxy.NextSequence(ctx, common.BKTableNameInstAsst)
	if err != nil {
		return id, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	asstInst.ID = int64(id)
	asstInst.OwnerID = ctx.SupplierAccount

	err = m.dbProxy.Table(common.BKTableNameInstAsst).Insert(ctx, asstInst)
	return id, err
}

func (m *associationInstance) CreateOneInstanceAssociation(ctx core.ContextParams, inputParam metadata.CreateOneInstanceAssociation) (*metadata.CreateOneDataResult, error) {
	inputParam.Data.OwnerID = ctx.SupplierAccount
	_, exists, err := m.isExists(ctx, inputParam.Data.InstID, inputParam.Data.AsstInstID, inputParam.Data.ObjectAsstID, inputParam.Data.Metadata)
	if nil != err {
		blog.Errorf("check instance (%#v)is duplicated error, rid: %s", inputParam.Data, ctx.ReqID)
		return nil, err
	}
	if exists {
		blog.Errorf("association instance (%#v)is duplicated, rid: %s", inputParam.Data, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommDuplicateItem, "association")
	}
	//check association kind
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.AssociationObjAsstIDField, Val: inputParam.Data.ObjectAsstID})
	_, exists, err = m.associationModel.isExists(ctx, cond)
	if nil != err {
		blog.Errorf("check asst kind(%#v)is not exist, rid: %s", inputParam.Data.ObjectAsstID, ctx.ReqID)
		return nil, err
	}
	if !exists {
		blog.Errorf("association asst kind(%#v)is not exist, rid: %s", inputParam.Data.ObjectAsstID, ctx.ReqID)
		return nil, ctx.Error.Error(common.CCErrorTopoAsstKindIsNotExist)
	}
	//check association inst
	exists, err = m.dependent.IsInstanceExist(ctx, inputParam.Data.ObjectID, uint64(inputParam.Data.InstID))
	if nil != err {
		return nil, err
	}

	if !exists {
		blog.Errorf("inst to asst is not exist objid(%#v), instid(%#v), rid: %s", inputParam.Data.ObjectID, inputParam.Data.InstID, ctx.ReqID)
		return nil, ctx.Error.Error(common.CCErrorAsstInstIsNotExist)
	}
	//check inst to asst
	exists, err = m.dependent.IsInstanceExist(ctx, inputParam.Data.AsstObjectID, uint64(inputParam.Data.AsstInstID))
	if nil != err {
		return nil, err
	}

	if !exists {
		blog.Errorf("asst inst is not exist objid(%#v), instid(%#v), rid: %s", inputParam.Data.ObjectID, inputParam.Data.InstID, ctx.ReqID)
		return nil, ctx.Error.Error(common.CCErrorInstToAsstIsNotExist)
	}
	id, err := m.save(ctx, inputParam.Data)
	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, err
}

func (m *associationInstance) CreateManyInstanceAssociation(ctx core.ContextParams, inputParam metadata.CreateManyInstanceAssociation) (*metadata.CreateManyDataResult, error) {
	dataResult := &metadata.CreateManyDataResult{}
	for itemIdx, item := range inputParam.Datas {
		item.OwnerID = ctx.SupplierAccount
		//check is exist
		_, exists, err := m.isExists(ctx, item.InstID, item.AsstInstID, item.ObjectAsstID, item.Metadata)
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
		//check asst kind
		_, exists, err = m.associationKind.isExists(ctx, item.ObjectAsstID)
		if nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}
		if !exists {
			blog.InfoJSON("CreateManyInstanceAssociation error. obj:%s,rid:%s", item.ObjectAsstID, ctx.ReqID)
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     ctx.Error.Error(common.CCErrorAsstInstIsNotExist).Error(),
				Code:        int64(common.CCErrorAsstInstIsNotExist),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}
		//check asst inst exist
		exists, err = m.dependent.IsInstanceExist(ctx, item.ObjectID, uint64(item.InstID))
		if nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}

		if !exists {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     ctx.Error.Error(common.CCErrorAsstInstIsNotExist).Error(),
				Code:        int64(common.CCErrorAsstInstIsNotExist),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}
		//check  inst to asst exist
		exists, err = m.dependent.IsInstanceExist(ctx, item.AsstObjectID, uint64(item.AsstInstID))
		if nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}

		if !exists {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     ctx.Error.Error(common.CCErrorInstToAsstIsNotExist).Error(),
				Code:        int64(common.CCErrorInstToAsstIsNotExist),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}
		//save asst inst
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

func (m *associationInstance) SearchInstanceAssociation(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	condition, err := mongo.NewConditionFromMapStr(inputParam.Condition)
	if nil != err {
		blog.Errorf("parse conditon  error [%#v], rid: %s", err, ctx.ReqID)
		return &metadata.QueryResult{}, err
	}
	ownerIDArr := []string{ctx.SupplierAccount, common.BKDefaultOwnerID}
	condition.Element(&mongo.In{Key: common.BKOwnerIDField, Val: ownerIDArr})
	inputParam.Condition = condition.ToMapStr()
	instAsstItems, err := m.searchInstanceAssociation(ctx, inputParam)
	if nil != err {
		blog.Errorf("search inst association array err [%#v], rid: %s", err, ctx.ReqID)
		return &metadata.QueryResult{}, err
	}

	dataResult := &metadata.QueryResult{}
	dataResult.Count, err = m.countInstanceAssociation(ctx, inputParam.Condition)
	dataResult.Info = make([]mapstr.MapStr, 0)
	if nil != err {
		blog.Errorf("search inst association count err [%#v], rid: %s", err, ctx.ReqID)
		return &metadata.QueryResult{}, err
	}
	for _, item := range instAsstItems {
		dataResult.Info = append(dataResult.Info, mapstr.NewFromStruct(item, "field"))
	}

	return dataResult, nil
}

func (m *associationInstance) DeleteInstanceAssociation(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	inputParam.Condition.Set(common.BKOwnerIDField, ctx.SupplierAccount)
	cnt, err := m.instCount(ctx, inputParam.Condition)
	if nil != err {
		blog.Errorf("delete inst association get inst [%#v] count err [%#v], rid: %s", inputParam.Condition, err, ctx.ReqID)
		return &metadata.DeletedCount{}, err
	}

	err = m.dbProxy.Table(common.BKTableNameInstAsst).Delete(ctx, inputParam.Condition)
	if nil != err {
		blog.Errorf("delete inst association [%#v] err [%#v], rid: %s", inputParam.Condition, err, ctx.ReqID)
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: cnt}, nil
}
