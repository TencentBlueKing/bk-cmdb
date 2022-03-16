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

package model

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
)

type modelClassification struct {
	model *modelManager
}

func (m *modelClassification) CreateOneModelClassification(kit *rest.Kit, inputParam metadata.CreateOneModelClassification) (*metadata.CreateOneDataResult, error) {

	if 0 == len(inputParam.Data.ClassificationID) {
		blog.Errorf("request(%s): it is failed to create the model classification, because of the classificationID (%#v) is not set", kit.Rid, inputParam.Data)
		return &metadata.CreateOneDataResult{}, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.ClassFieldClassificationID)
	}

	_, exists, err := m.isExists(kit, inputParam.Data.ClassificationID)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check if the classification ID (%s)is exists, error info is %s", kit.Rid, inputParam.Data.ClassificationID, err.Error())
		return nil, err
	}
	if exists {
		blog.Errorf("classification (%#v)is duplicated, rid: %s", inputParam.Data, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommDuplicateItem, inputParam.Data.ClassificationID)
	}

	inputParam.Data.OwnerID = kit.SupplierAccount

	id, err := m.save(kit, inputParam.Data)
	if nil != err {
		blog.Errorf("request(%s): it is failed to save the classification(%#v), error info is %s", kit.Rid, inputParam.Data, err.Error())
		return &metadata.CreateOneDataResult{}, err
	}
	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, err
}

func (m *modelClassification) CreateManyModelClassification(kit *rest.Kit, inputParam metadata.CreateManyModelClassifiaction) (*metadata.CreateManyDataResult, error) {

	dataResult := &metadata.CreateManyDataResult{
		CreateManyInfoResult: metadata.CreateManyInfoResult{
			Created:    []metadata.CreatedDataResult{},
			Repeated:   []metadata.RepeatedDataResult{},
			Exceptions: []metadata.ExceptionResult{},
		},
	}

	addExceptionFunc := func(idx int64, err errors.CCErrorCoder, classification *metadata.Classification) {
		dataResult.CreateManyInfoResult.Exceptions = append(dataResult.CreateManyInfoResult.Exceptions, metadata.ExceptionResult{
			OriginIndex: idx,
			Message:     err.Error(),
			Code:        int64(err.GetCode()),
			Data:        classification,
		})
	}

	for itemIdx, item := range inputParam.Data {

		if 0 == len(item.ClassificationID) {
			blog.Errorf("request(%s): it is failed to create the model classification, because of the classificationID (%#v) is not set", kit.Rid, item.ClassificationID)
			addExceptionFunc(int64(itemIdx), kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.ClassFieldClassificationID).(errors.CCErrorCoder), &item)
			continue
		}

		_, exists, err := m.isExists(kit, item.ClassificationID)
		if nil != err {
			blog.Errorf("request(%s): it is failed to check the classification ID (%s) is exists, error info is %s", kit.Rid, item.ClassificationID, err.Error())
			addExceptionFunc(int64(itemIdx), err.(errors.CCErrorCoder), &item)
			continue
		}

		if exists {
			dataResult.Repeated = append(dataResult.Repeated, metadata.RepeatedDataResult{OriginIndex: int64(itemIdx), Data: mapstr.NewFromStruct(item, "field")})
			continue
		}

		item.OwnerID = kit.SupplierAccount
		id, err := m.save(kit, item)
		if nil != err {
			blog.Errorf("request(%s): it is failed to save the classification(%#v), error info is %s", kit.Rid, item, err.Error())
			addExceptionFunc(int64(itemIdx), err.(errors.CCErrorCoder), &item)
			continue
		}

		dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{
			ID: id,
		})

	}

	return dataResult, nil
}
func (m *modelClassification) SetManyModelClassification(kit *rest.Kit, inputParam metadata.SetManyModelClassification) (*metadata.SetDataResult, error) {

	dataResult := &metadata.SetDataResult{
		Created:    []metadata.CreatedDataResult{},
		Updated:    []metadata.UpdatedDataResult{},
		Exceptions: []metadata.ExceptionResult{},
	}

	addExceptionFunc := func(idx int64, err errors.CCErrorCoder, classification *metadata.Classification) {
		dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
			OriginIndex: idx,
			Message:     err.Error(),
			Code:        int64(err.GetCode()),
			Data:        classification,
		})
	}

	for itemIdx, item := range inputParam.Data {

		if 0 == len(item.ClassificationID) {
			blog.Errorf("request(%s): it is failed to create the model classification, because of the classificationID (%#v) is not set", kit.Rid, item.ClassificationID)
			addExceptionFunc(int64(itemIdx), kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.ClassFieldClassificationID).(errors.CCErrorCoder), &item)
			continue
		}

		origin, exists, err := m.isExists(kit, item.ClassificationID)
		if nil != err {
			blog.Errorf("request(%s): it is failed to check the classification ID (%s) is exists, error info is %s", kit.Rid, item.ClassificationID, err.Error())
			addExceptionFunc(int64(itemIdx), err.(errors.CCErrorCoder), &item)
			continue
		}

		if exists {

			cond := mongo.NewCondition()
			cond.Element(&mongo.Eq{Key: metadata.ClassificationFieldID, Val: origin.ID})
			if _, err := m.update(kit, mapstr.NewFromStruct(item, "field"), cond); nil != err {
				blog.Errorf("request(%s): it is failed to update some fields(%#v) of the classification by the condition(%#v), error info is %s", kit.Rid, item, cond.ToMapStr(), err.Error())
				addExceptionFunc(int64(itemIdx), err.(errors.CCErrorCoder), &item)
				continue
			}

			dataResult.UpdatedCount.Count++
			dataResult.Updated = append(dataResult.Updated, metadata.UpdatedDataResult{
				OriginIndex: int64(itemIdx),
				ID:          uint64(origin.ID),
			})
			continue
		}

		item.OwnerID = kit.SupplierAccount

		id, err := m.save(kit, item)
		if nil != err {
			blog.Errorf("request(%s): it is failed to save the classification(%#v), error info is %s", kit.Rid, item, err.Error())
			addExceptionFunc(int64(itemIdx), err.(errors.CCErrorCoder), &item)
			continue
		}

		dataResult.CreatedCount.Count++
		dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{
			ID: id,
		})

	}

	return dataResult, nil
}

func (m *modelClassification) SetOneModelClassification(kit *rest.Kit, inputParam metadata.SetOneModelClassification) (*metadata.SetDataResult, error) {

	dataResult := &metadata.SetDataResult{
		Created:    []metadata.CreatedDataResult{},
		Updated:    []metadata.UpdatedDataResult{},
		Exceptions: []metadata.ExceptionResult{},
	}

	if 0 == len(inputParam.Data.ClassificationID) {
		blog.Errorf("request(%s): it is failed to set the model classification, because of the classificationID (%#v) is not set", kit.Rid, inputParam.Data)
		return dataResult, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.ClassFieldClassificationID)
	}

	origin, exists, err := m.isExists(kit, inputParam.Data.ClassificationID)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check the classification ID (%s) is exists, error info is %s", kit.Rid, inputParam.Data.ClassificationID, err.Error())
		return dataResult, err
	}

	addExceptionFunc := func(idx int64, err errors.CCErrorCoder, classification *metadata.Classification) {
		dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
			OriginIndex: idx,
			Message:     err.Error(),
			Code:        int64(err.GetCode()),
			Data:        classification,
		})
	}
	if exists {

		cond := mongo.NewCondition()
		cond.Element(&mongo.Eq{Key: metadata.ClassificationFieldID, Val: origin.ID})
		if _, err := m.update(kit, mapstr.NewFromStruct(inputParam.Data, "field"), cond); nil != err {
			blog.Errorf("request(%s): it is failed to update some fields(%#v) for a classification by the condition(%#v), error info is %s", kit.Rid, inputParam.Data, cond.ToMapStr(), err.Error())
			addExceptionFunc(0, err.(errors.CCErrorCoder), &inputParam.Data)
			return dataResult, nil
		}
		dataResult.UpdatedCount.Count++
		dataResult.Updated = append(dataResult.Updated, metadata.UpdatedDataResult{ID: uint64(origin.ID)})
		return dataResult, err
	}

	inputParam.Data.OwnerID = kit.SupplierAccount
	id, err := m.save(kit, inputParam.Data)
	if nil != err {
		blog.Errorf("request(%s): it is failed to save the classification(%#v), error info is %s", kit.Rid, inputParam.Data, err.Error())
		addExceptionFunc(0, err.(errors.CCErrorCoder), origin)
		return dataResult, err
	}
	dataResult.CreatedCount.Count++
	dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{ID: id})
	return dataResult, err
}

func (m *modelClassification) UpdateModelClassification(kit *rest.Kit, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	cond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), kit.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition(%#v) from mapstr into condition object, error info is %s", kit.Rid, inputParam.Condition, err.Error())
		return &metadata.UpdatedCount{}, err
	}

	cnt, err := m.update(kit, inputParam.Data, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to update some fields(%#v) for some classifications by the condition(%#v), error info is %s", kit.Rid, inputParam.Data, inputParam.Condition, err.Error())
		return &metadata.UpdatedCount{}, err
	}
	return &metadata.UpdatedCount{Count: cnt}, nil
}

func (m *modelClassification) DeleteModelClassification(kit *rest.Kit, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	deleteCond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), kit.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr into condition object, error info is %s", kit.Rid, inputParam.Condition, err.Error())
		return &metadata.DeletedCount{}, kit.CCError.New(common.CCErrCommHTTPInputInvalid, err.Error())
	}

	cnt, exists, err := m.hasModel(kit, deleteCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether the classifications which are marked by the condition (%#v) have some models, error info is %s", kit.Rid, deleteCond.ToMapStr(), err.Error())
		return &metadata.DeletedCount{}, err
	}
	if exists {
		return &metadata.DeletedCount{}, kit.CCError.Error(common.CCErrTopoObjectClassificationHasObject)
	}

	cnt, err = m.delete(kit, deleteCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to delete the classification whci are marked by the condition(%#v), error info is %s", kit.Rid, deleteCond.ToMapStr(), err.Error())
		return &metadata.DeletedCount{}, err
	}

	return &metadata.DeletedCount{Count: cnt}, nil
}

func (m *modelClassification) SearchModelClassification(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryModelClassificationDataResult, error) {

	dataResult := &metadata.QueryModelClassificationDataResult{
		Info: []metadata.Classification{},
	}
	searchCond, err := mongo.NewConditionFromMapStr(inputParam.Condition.ToMapInterface())
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr into condition object, error info is %s", kit.Rid, inputParam.Condition, err.Error())
		return dataResult, err
	}

	totalCount, err := m.count(kit, searchCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to get the count by the condition (%#v), error info is %s", kit.Rid, searchCond.ToMapStr(), err.Error())
		return dataResult, err
	}

	classificationItems, err := m.search(kit, searchCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to search some classifications by the condition (%#v), error info is %s", kit.Rid, searchCond.ToMapStr(), err.Error())
		return dataResult, err
	}

	dataResult.Count = int64(totalCount)
	dataResult.Info = classificationItems
	return dataResult, nil
}
