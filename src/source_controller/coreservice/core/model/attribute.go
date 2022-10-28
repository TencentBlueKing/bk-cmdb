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
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/lock"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
)

type modelAttribute struct {
	model    *modelManager
	language language.CCLanguageIf
}

// CreateModelAttributes TODO
func (m *modelAttribute) CreateModelAttributes(kit *rest.Kit, objID string, inputParam metadata.CreateModelAttributes) (dataResult *metadata.CreateManyDataResult, err error) {

	dataResult = &metadata.CreateManyDataResult{
		CreateManyInfoResult: metadata.CreateManyInfoResult{
			Created:    []metadata.CreatedDataResult{},
			Repeated:   []metadata.RepeatedDataResult{},
			Exceptions: []metadata.ExceptionResult{},
		},
	}

	if err := m.model.isValid(kit, objID); nil != err {
		blog.Errorf("CreateModelAttributes failed, validate model(%s) failed, err: %s, rid: %s", objID, err.Error(), kit.Rid)
		return dataResult, err
	}

	addExceptionFunc := func(idx int64, err errors.CCErrorCoder, attr *metadata.Attribute) {
		dataResult.CreateManyInfoResult.Exceptions = append(dataResult.CreateManyInfoResult.Exceptions, metadata.ExceptionResult{
			OriginIndex: idx,
			Message:     err.Error(),
			Code:        int64(err.GetCode()),
			Data:        attr,
		})
	}

	for attrIdx, attr := range inputParam.Attributes {
		// fmt.Sprintf("coreservice:create:model:%s:attr:%s", objID, attr.PropertyID)
		redisKey := lock.GetLockKey(lock.CreateModuleAttrFormat, objID, attr.PropertyID)

		locker := lock.NewLocker(redis.Client())
		locked, err := locker.Lock(redisKey, time.Second*35)
		defer locker.Unlock()
		if err != nil {
			blog.ErrorJSON("create model error. get create look error. err:%s, input:%s, rid:%s", err.Error(), inputParam, kit.Rid)
			addExceptionFunc(int64(attrIdx), kit.CCError.CCErrorf(common.CCErrCommRedisOPErr), &attr)
			continue
		}
		if !locked {
			blog.ErrorJSON("create model have same task in progress. input:%s, rid:%s", inputParam, kit.Rid)
			addExceptionFunc(int64(attrIdx), kit.CCError.CCErrorf(common.CCErrCommOPInProgressErr, fmt.Sprintf("create object(%s) attribute(%s)", attr.ObjectID, attr.PropertyName)), &attr)
			continue
		}
		if attr.IsPre {
			if attr.PropertyID == common.BKInstNameField {
				language := util.GetLanguage(kit.Header)
				lang := m.language.CreateDefaultCCLanguageIf(language)
				attr.PropertyName = util.FirstNotEmptyString(lang.Language("common_property_"+attr.PropertyID), attr.PropertyName, attr.PropertyID)
			}
		}

		attr.OwnerID = kit.SupplierAccount
		_, exists, err := m.isExists(kit, attr.ObjectID, attr.PropertyID, attr.BizID)
		blog.V(5).Infof("CreateModelAttributes isExists info. property id:%s, bizID:%#v, exit:%v, rid:%s", attr.PropertyID, attr.BizID, exists, kit.Rid)
		if nil != err {
			blog.Errorf("CreateModelAttributes failed, attribute field propertyID(%s) exists, err: %s, rid: %s", attr.PropertyID, err.Error(), kit.Rid)
			addExceptionFunc(int64(attrIdx), err.(errors.CCErrorCoder), &attr)
			continue
		}

		if exists {
			dataResult.CreateManyInfoResult.Repeated = append(dataResult.CreateManyInfoResult.Repeated, metadata.RepeatedDataResult{
				OriginIndex: int64(attrIdx),
				Data:        mapstr.NewFromStruct(attr, "field"),
			})
			continue
		}
		id, err := m.save(kit, attr)
		if nil != err {
			blog.Errorf("CreateModelAttributes failed, failed to save the attribute(%#v), err: %s, rid: %s", attr, err.Error(), kit.Rid)
			addExceptionFunc(int64(attrIdx), err.(errors.CCErrorCoder), &attr)
			continue
		}

		dataResult.CreateManyInfoResult.Created = append(dataResult.CreateManyInfoResult.Created, metadata.CreatedDataResult{
			OriginIndex: int64(attrIdx),
			ID:          id,
		})
	}

	return dataResult, nil
}

// SetModelAttributes TODO
func (m *modelAttribute) SetModelAttributes(kit *rest.Kit, objID string, inputParam metadata.SetModelAttributes) (dataResult *metadata.SetDataResult, err error) {

	dataResult = &metadata.SetDataResult{
		Created:    []metadata.CreatedDataResult{},
		Updated:    []metadata.UpdatedDataResult{},
		Exceptions: []metadata.ExceptionResult{},
	}

	if err := m.model.isValid(kit, objID); nil != err {
		return dataResult, err
	}

	addExceptionFunc := func(idx int64, err errors.CCErrorCoder, attr *metadata.Attribute) {
		dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
			OriginIndex: idx,
			Message:     err.Error(),
			Code:        int64(err.GetCode()),
			Data:        attr,
		})
	}

	for attrIdx, attr := range inputParam.Attributes {

		existsAttr, exists, err := m.isExists(kit, attr.ObjectID, attr.PropertyID, attr.BizID)
		if nil != err {
			addExceptionFunc(int64(attrIdx), err.(errors.CCErrorCoder), &attr)
			continue
		}
		attr.OwnerID = kit.SupplierAccount
		if exists {
			cond := mongo.NewCondition()
			cond.Element(&mongo.Eq{Key: metadata.AttributeFieldSupplierAccount, Val: kit.SupplierAccount})
			cond.Element(&mongo.Eq{Key: metadata.AttributeFieldID, Val: existsAttr.ID})

			_, err := m.update(kit, mapstr.NewFromStruct(attr, "field"), cond)
			if nil != err {
				blog.Errorf("SetModelAttributes failed, failed to update the attribute(%#v) by the condition(%#v), err: %s, rid: %s", attr, cond.ToMapStr(), err.Error(), kit.Rid)
				addExceptionFunc(int64(attrIdx), err.(errors.CCErrorCoder), &attr)
				continue
			}
			dataResult.Updated = append(dataResult.Updated, metadata.UpdatedDataResult{
				OriginIndex: int64(attrIdx),
				ID:          uint64(existsAttr.ID),
			})
			continue
		}
		id, err := m.save(kit, attr)
		if nil != err {
			blog.Errorf("SetModelAttributes failed, failed to save the attribute(%#v), err: %s, rid: %s", attr, err.Error(), kit.Rid)
			addExceptionFunc(int64(attrIdx), err.(errors.CCErrorCoder), &attr)
			continue
		}

		dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{
			OriginIndex: int64(attrIdx),
			ID:          id,
		})

	}

	return dataResult, nil
}

// UpdateModelAttributes TODO
func (m *modelAttribute) UpdateModelAttributes(kit *rest.Kit, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	if err := m.model.isValid(kit, objID); nil != err {
		blog.Errorf("UpdateModelAttributes failed, validate model(%s) failed, err: %s, rid: %s", objID, err.Error(), kit.Rid)
		return &metadata.UpdatedCount{}, err
	}

	cond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), kit.SupplierAccount))
	if nil != err {
		blog.Errorf("UpdateModelAttributes failed, failed to convert mapstr(%#v) into a condition object, err: %s, rid: %s", inputParam.Condition, err.Error(), kit.Rid)
		return &metadata.UpdatedCount{}, err
	}

	cnt, err := m.update(kit, inputParam.Data, cond)
	if nil != err {
		blog.ErrorJSON("UpdateModelAttributes failed, update attributes failed, model:%s, attributes:%s, condition: %s, err: %s, rid: %s", inputParam.Data, objID, cond, err.Error(), kit.Rid)
		return &metadata.UpdatedCount{}, err
	}

	return &metadata.UpdatedCount{Count: cnt}, nil
}

// UpdateModelAttributeIndex update model attribute index
func (m *modelAttribute) UpdateModelAttributeIndex(kit *rest.Kit, objID string, id int64,
	input *metadata.UpdateAttrIndexInput) error {

	// check if attribute exists
	attrCond := mapstr.MapStr{
		common.BKFieldID:    id,
		common.BKObjIDField: objID,
		common.BKAppIDField: input.BizID,
	}
	attrCond = util.SetQueryOwner(attrCond, kit.SupplierAccount)
	cnt, err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(attrCond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("check if attribute exists failed, err: %v, cond: %+v, rid: %s", err, attrCond, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	if cnt == 0 {
		blog.Errorf("attributes is not exist, condition: %+v, rid: %s", attrCond, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKFieldID)
	}

	// check if index is used by attribute in the same group of the same biz, if not, use it directly
	indexCond := mapstr.MapStr{
		common.BKObjIDField:         objID,
		common.BKAppIDField:         input.BizID,
		common.BKPropertyGroupField: input.PropertyGroup,
		common.BKPropertyIndexField: input.PropertyIndex,
	}
	indexCond = util.SetQueryOwner(indexCond, kit.SupplierAccount)
	count, err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(indexCond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("check if index is used failed, err: %v, cond: %+v, rid: %s", err, indexCond, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	// increase attributes index whose index is greater than the current index and is conflict with it
	if count > 0 {
		incCond := mapstr.MapStr{
			common.BKObjIDField:         objID,
			common.BKAppIDField:         input.BizID,
			common.BKPropertyGroupField: input.PropertyGroup,
			common.BKPropertyIndexField: mapstr.MapStr{common.BKDBGTE: input.PropertyIndex},
			common.BKFieldID:            mapstr.MapStr{common.BKDBNE: id},
		}

		incData := mapstr.MapStr{common.BKPropertyIndexField: int64(1)}
		err = mongodb.Client().Table(common.BKTableNameObjAttDes).UpdateMultiModel(kit.Ctx, incCond,
			types.ModeUpdate{Op: "inc", Doc: incData})
		if err != nil {
			blog.Errorf("increase attributes index failed, err: %v, cond: %+v, rid: %s", err, incCond, kit.Rid)
			return kit.CCError.Error(common.CCErrCommDBSelectFailed)
		}
	}

	// update attribute index now
	data := mapstr.MapStr{
		common.BKPropertyIndexField: input.PropertyIndex,
		common.BKPropertyGroupField: input.PropertyGroup,
	}
	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Update(kit.Ctx, attrCond, data)
	if err != nil {
		blog.Errorf("update attribute index failed, err: %v, cond: %+v, rid: %s", err, attrCond, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	return nil
}

// UpdateModelAttributesByCondition TODO
func (m *modelAttribute) UpdateModelAttributesByCondition(kit *rest.Kit, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	cond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), kit.SupplierAccount))
	if nil != err {
		blog.Errorf("UpdateModelAttributesByCondition failed, failed to convert mapstr(%#v) into a condition object, err: %s, rid: %s", inputParam.Condition, err.Error(), kit.Rid)
		return &metadata.UpdatedCount{}, err
	}

	cnt, err := m.update(kit, inputParam.Data, cond)
	if nil != err {
		blog.Errorf("UpdateModelAttributesByCondition failed, failed to update fields (%#v) by condition(%#v), err: %s, rid: %s", inputParam.Data, cond.ToMapStr(), err.Error(), kit.Rid)
		return &metadata.UpdatedCount{}, err
	}

	return &metadata.UpdatedCount{Count: cnt}, nil
}

// DeleteModelAttributes TODO
func (m *modelAttribute) DeleteModelAttributes(kit *rest.Kit, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	if err := m.model.isValid(kit, objID); nil != err {
		blog.Errorf("request(%s): it is failed to check if the model(%s) is valid, error info is %s", kit.Rid, objID, err.Error())
		return &metadata.DeletedCount{}, err
	}

	cond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), kit.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert from mapstr(%#v) into a condition object, error info is %s", kit.Rid, inputParam.Condition, err.Error())
		return &metadata.DeletedCount{}, err
	}

	cond.Element(&mongo.Eq{Key: metadata.AttributeFieldSupplierAccount, Val: kit.SupplierAccount})
	cnt, err := m.delete(kit, cond)
	return &metadata.DeletedCount{Count: cnt}, err
}

// SearchModelAttributes search model's attributes
func (m *modelAttribute) SearchModelAttributes(kit *rest.Kit, objID string, inputParam metadata.QueryCondition) (
	*metadata.QueryModelAttributeDataResult, error) {

	if err := m.model.isValid(kit, objID); err != nil {
		blog.Errorf("failed to check if the model(%s) is valid, err: %v, rid: %s", objID, err, kit.Rid)
		return nil, err
	}

	inputParam.Condition = util.SetQueryOwner(inputParam.Condition, kit.SupplierAccount)
	inputParam.Condition[common.BKObjIDField] = objID

	attrResult, err := m.newSearch(kit, inputParam.Condition)
	if err != nil {
		blog.Errorf("failed to search the attributes of the model(%s), err: %v, rid: %s", objID, err, kit.Rid)
		return nil, err
	}

	return &metadata.QueryModelAttributeDataResult{
		Count: int64(len(attrResult)),
		Info:  attrResult,
	}, nil
}

// SearchModelAttributesByCondition TODO
func (m *modelAttribute) SearchModelAttributesByCondition(kit *rest.Kit, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeDataResult, error) {

	dataResult := &metadata.QueryModelAttributeDataResult{
		Info: []metadata.Attribute{},
	}

	inputParam.Condition = util.SetQueryOwner(inputParam.Condition, kit.SupplierAccount)

	attrResult, err := m.searchWithSort(kit, inputParam)
	if nil != err {
		blog.Errorf("request(%s): it is failed to search the attributes of the model(%+v), error info is %s", kit.Rid, inputParam, err.Error())
		return &metadata.QueryModelAttributeDataResult{}, err
	}

	dataResult.Count = int64(len(attrResult))
	dataResult.Info = attrResult
	return dataResult, nil
}
