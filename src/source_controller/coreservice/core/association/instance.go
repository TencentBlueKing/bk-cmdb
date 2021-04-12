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

func (m *associationInstance) isExists(kit *rest.Kit, instID, asstInstID int64, objAsstID, objID string, bizID int64) (
	bool, error) {

	cond := map[string]interface{}{
		common.BKInstIDField:             instID,
		common.BKAsstInstIDField:         asstInstID,
		common.AssociationObjAsstIDField: objAsstID,
	}

	if bizID > 0 {
		cond[common.BKAppIDField] = bizID
	}

	count, err := m.countInstanceAssociation(kit, objID, cond)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *associationInstance) searchInstanceAssociation(kit *rest.Kit, objID string, param metadata.QueryCondition) (
	[]metadata.InstAsst, error) {

	results := make([]metadata.InstAsst, 0)
	asstTableName := common.GetObjectInstAsstTableName(objID, kit.SupplierAccount)
	instHandler := mongodb.Client().Table(asstTableName).Find(param.Condition).Fields(param.Fields...)
	err := instHandler.Start(uint64(param.Page.Start)).Limit(uint64(param.Page.Limit)).
		Sort(param.Page.Sort).All(kit.Ctx, &results)
	return results, err
}

func (m *associationInstance) countInstanceAssociation(kit *rest.Kit, objID string, cond mapstr.MapStr) (uint64, error) {
	asstTableName := common.GetObjectInstAsstTableName(objID, kit.SupplierAccount)
	return mongodb.Client().Table(asstTableName).Find(cond).Count(kit.Ctx)
}

func (m *associationInstance) save(kit *rest.Kit, asstInst metadata.InstAsst) (id uint64, err error) {
	id, err = mongodb.Client().NextSequence(kit.Ctx, common.BKTableNameInstAsst)
	if err != nil {
		return id, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	asstInst.ID = int64(id)
	asstInst.OwnerID = kit.SupplierAccount

	objInstAsstTableName := common.GetObjectInstAsstTableName(asstInst.ObjectID, kit.SupplierAccount)
	err = mongodb.Client().Table(objInstAsstTableName).Insert(kit.Ctx, asstInst)
	if err != nil {
		return id, err
	}

	// do not insert twice for self related association
	if asstInst.ObjectID == asstInst.AsstObjectID {
		return id, nil
	}

	asstObjInstAsstTableName := common.GetObjectInstAsstTableName(asstInst.AsstObjectID, kit.SupplierAccount)
	err = mongodb.Client().Table(asstObjInstAsstTableName).Insert(kit.Ctx, asstInst)
	return id, err
}

func (m *associationInstance) deleteInstanceAssociation(kit *rest.Kit, objID string,
	cond mapstr.MapStr) (uint64, error) {
	asstInstTableName := common.GetObjectInstAsstTableName(objID, kit.SupplierAccount)
	associations := make([]metadata.InstAsst, 0)
	if err := mongodb.Client().Table(asstInstTableName).Find(cond).Fields(common.BKObjIDField, common.BKAsstObjIDField).
		All(kit.Ctx, &associations); err != nil {
		blog.ErrorJSON("delete instance association error. objID: %s, cond: %s, err: %s, rid: %s",
			objID, cond, err.Error(), kit.Rid)
		return 0, err
	}

	objIDMap := make(map[string]struct{})
	for _, asst := range associations {
		var asstObjID string
		if asst.ObjectID != objID {
			asstObjID = asst.ObjectID
		} else if asst.AsstObjectID != objID {
			asstObjID = asst.AsstObjectID
		} else {
			// 自关联再循环外处理
			continue
		}

		if _, exists := objIDMap[asstObjID]; exists {
			continue
		}
		objIDMap[asstObjID] = struct{}{}

		asstTableName := common.GetObjectInstAsstTableName(asstObjID, kit.SupplierAccount)
		err := mongodb.Client().Table(asstTableName).Delete(kit.Ctx, cond)
		if err != nil {
			blog.ErrorJSON("delete instance association error. objID: %s, cond: %s, err: %s, rid: %s",
				asstObjID, cond, err.Error(), kit.Rid)
			return 0, err
		}
	}

	cnt, err := mongodb.Client().Table(asstInstTableName).DeleteMany(kit.Ctx, cond)
	if err != nil {
		blog.ErrorJSON("delete instance association error. objID: %s, cond: %s, err: %s, rid: %s",
			objID, cond, err.Error(), kit.Rid)
		return 0, err
	}
	return cnt, nil
}

func (m *associationInstance) CreateOneInstanceAssociation(kit *rest.Kit, inputParam metadata.CreateOneInstanceAssociation) (*metadata.CreateOneDataResult, error) {
	inputParam.Data.OwnerID = kit.SupplierAccount
	exists, err := m.isExists(kit, inputParam.Data.InstID, inputParam.Data.AsstInstID, inputParam.Data.ObjectAsstID,
		inputParam.Data.ObjectID, inputParam.Data.BizID)
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
		exists, err := m.isExists(kit, item.InstID, item.AsstInstID, item.ObjectAsstID, item.ObjectID, item.BizID)
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

func (m *associationInstance) SearchInstanceAssociation(kit *rest.Kit, objID string, param metadata.QueryCondition) (
	*metadata.QueryResult, error) {

	if len(objID) == 0 {
		return &metadata.QueryResult{}, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIDField)
	}

	param.Condition = util.SetQueryOwner(param.Condition, kit.SupplierAccount)
	instAsstItems, err := m.searchInstanceAssociation(kit, objID, param)
	if nil != err {
		blog.ErrorJSON("search inst association err: %s, objID: %s, param: %s, rid: %s", err, objID, param, kit.Rid)
		return &metadata.QueryResult{}, err
	}

	dataResult := &metadata.QueryResult{}

	// the InstAsst number will be counted by default.
	if !param.DisableCounter {
		count, err := m.countInstanceAssociation(kit, objID, param.Condition)
		if nil != err {
			blog.Errorf("search inst association count err [%#v], rid: %s", err, kit.Rid)
			return &metadata.QueryResult{}, err
		}
		dataResult.Count = count
	}
	dataResult.Info = make([]mapstr.MapStr, 0)
	for _, item := range instAsstItems {
		dataResult.Info = append(dataResult.Info, mapstr.NewFromStruct(item, "field"))
	}

	return dataResult, nil
}

func (m *associationInstance) DeleteInstanceAssociation(kit *rest.Kit, objID string, inputParam metadata.DeleteOption) (
	*metadata.DeletedCount, error) {

	if len(objID) == 0 {
		return &metadata.DeletedCount{}, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIDField)
	}

	inputParam.Condition = util.SetModOwner(inputParam.Condition, kit.SupplierAccount)

	cnt, err := m.deleteInstanceAssociation(kit, objID, inputParam.Condition)
	if nil != err {
		blog.Errorf("delete inst association [%#v] err [%#v], rid: %s", inputParam.Condition, err, kit.Rid)
		return &metadata.DeletedCount{}, err
	}
	return &metadata.DeletedCount{Count: cnt}, nil
}
