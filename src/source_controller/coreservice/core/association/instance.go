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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

type associationInstance struct {
	*associationKind
	*associationModel
	dependent OperationDependencies
}

func (m *associationInstance) isExists(kit *rest.Kit, instID, asstInstID int64, objAsstID string, bizID int64) (origin *metadata.InstAsst, exists bool, err error) {
	cond := mongo.NewCondition()
	origin = &metadata.InstAsst{}
	cond.Element(
		&mongo.Eq{Key: common.BKInstIDField, Val: instID},
		&mongo.Eq{Key: common.BKAsstInstIDField, Val: asstInstID},
		&mongo.Eq{Key: common.AssociationObjAsstIDField, Val: objAsstID})
	if bizID > 0 {
		cond.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bizID})
	}

	err = mongodb.Client().Table(common.BKTableNameInstAsst).Find(cond.ToMapStr()).One(kit.Ctx, origin)
	if mongodb.Client().IsNotFoundError(err) {
		return origin, !mongodb.Client().IsNotFoundError(err), nil
	}
	return origin, !mongodb.Client().IsNotFoundError(err), err
}

func (m *associationInstance) instCount(kit *rest.Kit, cond mapstr.MapStr) (cnt uint64, err error) {
	innerCnt, err := mongodb.Client().Table(common.BKTableNameInstAsst).Find(cond).Count(kit.Ctx)
	return innerCnt, err
}

func (m *associationInstance) searchInstanceAssociation(kit *rest.Kit, inputParam metadata.QueryCondition) (results []metadata.InstAsst, err error) {
	results = []metadata.InstAsst{}
	instHandler := mongodb.Client().Table(common.BKTableNameInstAsst).Find(inputParam.Condition).Fields(inputParam.Fields...)
	err = instHandler.Start(uint64(inputParam.Page.Start)).Limit(uint64(inputParam.Page.Limit)).Sort(inputParam.Page.Sort).All(kit.Ctx, &results)
	return results, err
}

func (m *associationInstance) countInstanceAssociation(kit *rest.Kit, cond mapstr.MapStr) (count uint64, err error) {
	count, err = mongodb.Client().Table(common.BKTableNameInstAsst).Find(cond).Count(kit.Ctx)

	return count, err
}

func (m *associationInstance) save(kit *rest.Kit, asstInst metadata.InstAsst) (id uint64, err error) {

	id, err = mongodb.Client().NextSequence(kit.Ctx, common.BKTableNameInstAsst)
	if err != nil {
		return id, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	asstInst.ID = int64(id)
	asstInst.OwnerID = kit.SupplierAccount

	err = mongodb.Client().Table(common.BKTableNameInstAsst).Insert(kit.Ctx, asstInst)
	return id, err
}

func (m *associationInstance) CreateOneInstanceAssociation(kit *rest.Kit, inputParam metadata.CreateOneInstanceAssociation) (*metadata.CreateOneDataResult, error) {
	inputParam.Data.OwnerID = kit.SupplierAccount
	_, exists, err := m.isExists(kit, inputParam.Data.InstID, inputParam.Data.AsstInstID, inputParam.Data.ObjectAsstID, inputParam.Data.BizID)
	if nil != err {
		blog.Errorf("check instance (%#v)is duplicated error, rid: %s", inputParam.Data, kit.Rid)
		return nil, err
	}
	if exists {
		blog.Errorf("association instance (%#v)is duplicated, rid: %s", inputParam.Data, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommDuplicateItem, "association")
	}
	//check association kind
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.AssociationObjAsstIDField, Val: inputParam.Data.ObjectAsstID})
	_, exists, err = m.associationModel.isExists(kit, cond)
	if nil != err {
		blog.Errorf("check asst kind(%#v)is not exist, rid: %s", inputParam.Data.ObjectAsstID, kit.Rid)
		return nil, err
	}
	if !exists {
		blog.Errorf("association asst kind(%#v)is not exist, rid: %s", inputParam.Data.ObjectAsstID, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrorTopoAsstKindIsNotExist)
	}
	//check association inst
	exists, err = m.dependent.IsInstanceExist(kit, inputParam.Data.ObjectID, uint64(inputParam.Data.InstID))
	if nil != err {
		return nil, err
	}

	if !exists {
		blog.Errorf("inst to asst is not exist objid(%#v), instid(%#v), rid: %s", inputParam.Data.ObjectID, inputParam.Data.InstID, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrorAsstInstIsNotExist)
	}
	//check inst to asst
	exists, err = m.dependent.IsInstanceExist(kit, inputParam.Data.AsstObjectID, uint64(inputParam.Data.AsstInstID))
	if nil != err {
		return nil, err
	}

	if !exists {
		blog.Errorf("asst inst is not exist objid(%#v), instid(%#v), rid: %s", inputParam.Data.ObjectID, inputParam.Data.InstID, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrorInstToAsstIsNotExist)
	}
	id, err := m.save(kit, inputParam.Data)
	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, err
}

func (m *associationInstance) CreateManyInstanceAssociation(kit *rest.Kit, inputParam metadata.CreateManyInstanceAssociation) (*metadata.CreateManyDataResult, error) {
	dataResult := &metadata.CreateManyDataResult{}
	for itemIdx, item := range inputParam.Datas {
		item.OwnerID = kit.SupplierAccount
		//check is exist
		_, exists, err := m.isExists(kit, item.InstID, item.AsstInstID, item.ObjectAsstID, item.BizID)
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
		_, exists, err = m.associationKind.isExists(kit, item.ObjectAsstID)
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
			blog.InfoJSON("CreateManyInstanceAssociation error. obj:%s,rid:%s", item.ObjectAsstID, kit.Rid)
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     kit.CCError.Error(common.CCErrorAsstInstIsNotExist).Error(),
				Code:        int64(common.CCErrorAsstInstIsNotExist),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}
		//check asst inst exist
		exists, err = m.dependent.IsInstanceExist(kit, item.ObjectID, uint64(item.InstID))
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
				Message:     kit.CCError.Error(common.CCErrorAsstInstIsNotExist).Error(),
				Code:        int64(common.CCErrorAsstInstIsNotExist),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}
		//check  inst to asst exist
		exists, err = m.dependent.IsInstanceExist(kit, item.AsstObjectID, uint64(item.AsstInstID))
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
				Message:     kit.CCError.Error(common.CCErrorInstToAsstIsNotExist).Error(),
				Code:        int64(common.CCErrorInstToAsstIsNotExist),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}
		//save asst inst
		id, err := m.save(kit, item)
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

func (m *associationInstance) SearchInstanceAssociation(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	inputParam.Condition = util.SetQueryOwner(inputParam.Condition, kit.SupplierAccount)
	instAsstItems, err := m.searchInstanceAssociation(kit, inputParam)
	if nil != err {
		blog.Errorf("search inst association array err [%#v], rid: %s", err, kit.Rid)
		return &metadata.QueryResult{}, err
	}

	dataResult := &metadata.QueryResult{}
	dataResult.Count, err = m.countInstanceAssociation(kit, inputParam.Condition)
	dataResult.Info = make([]mapstr.MapStr, 0)
	if nil != err {
		blog.Errorf("search inst association count err [%#v], rid: %s", err, kit.Rid)
		return &metadata.QueryResult{}, err
	}
	for _, item := range instAsstItems {
		dataResult.Info = append(dataResult.Info, mapstr.NewFromStruct(item, "field"))
	}

	return dataResult, nil
}

func (m *associationInstance) DeleteInstanceAssociation(kit *rest.Kit, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	inputParam.Condition = util.SetModOwner(inputParam.Condition, kit.SupplierAccount)
	cnt, err := m.instCount(kit, inputParam.Condition)
	if nil != err {
		blog.Errorf("delete inst association get inst [%#v] count err [%#v], rid: %s", inputParam.Condition, err, kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	err = mongodb.Client().Table(common.BKTableNameInstAsst).Delete(kit.Ctx, inputParam.Condition)
	if nil != err {
		blog.Errorf("delete inst association [%#v] err [%#v], rid: %s", inputParam.Condition, err, kit.Rid)
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: cnt}, nil
}
