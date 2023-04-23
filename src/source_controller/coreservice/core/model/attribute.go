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
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
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

var forbiddenCreateAttrObjList = []string{
	common.BKInnerObjIDProject,
}

// CreateTableModelAttributes create model table attributes
func (m *modelAttribute) CreateTableModelAttributes(kit *rest.Kit, objID string,
	inputParam metadata.CreateModelAttributes) (*metadata.CreateManyDataResult, error) {

	dataResult := &metadata.CreateManyDataResult{
		CreateManyInfoResult: metadata.CreateManyInfoResult{
			Created:    []metadata.CreatedDataResult{},
			Repeated:   []metadata.RepeatedDataResult{},
			Exceptions: []metadata.ExceptionResult{},
		},
	}
	if err := m.model.isValid(kit, objID); err != nil {
		blog.Errorf("validate model(%s) failed, err: %v, rid: %s", objID, err, kit.Rid)
		return dataResult, err
	}

	addExceptionFunc := func(idx int64, err errors.CCErrorCoder, attr *metadata.Attribute) {
		dataResult.CreateManyInfoResult.Exceptions = append(dataResult.CreateManyInfoResult.Exceptions,
			metadata.ExceptionResult{OriginIndex: idx,
				Message: err.Error(),
				Code:    int64(err.GetCode()),
				Data:    attr,
			})
	}

	for attrIdx, attr := range inputParam.Attributes {
		redisKey := lock.GetLockKey(lock.CreateModuleAttrFormat, objID, attr.PropertyID)

		locker := lock.NewLocker(redis.Client())
		locked, err := locker.Lock(redisKey, time.Second*35)
		defer locker.Unlock()
		if err != nil {
			blog.Errorf("get create lock error, input: %+v, err: %v, rid: %s", inputParam, err, kit.Rid)
			addExceptionFunc(int64(attrIdx), kit.CCError.CCErrorf(common.CCErrCommRedisOPErr), &attr)
			continue
		}

		if !locked {
			blog.Errorf("create model have same task in progress. input: %v, rid: %s", inputParam, kit.Rid)
			addExceptionFunc(int64(attrIdx), kit.CCError.CCErrorf(common.CCErrCommOPInProgressErr,
				fmt.Sprintf("create table object(%s) attribute(%s)", attr.ObjectID, attr.PropertyName)), &attr)
			continue
		}

		if attr.IsPre {
			if attr.PropertyID == common.BKInstNameField {
				lang := m.language.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header))
				attr.PropertyName = util.FirstNotEmptyString(lang.Language("common_property_"+attr.PropertyID),
					attr.PropertyName, attr.PropertyID)
			}
		}

		attr.ObjectID = objID
		attr.OwnerID = kit.SupplierAccount
		_, exists, err := m.isExists(kit, attr.ObjectID, attr.PropertyID, attr.BizID)
		blog.V(5).Infof("table model attributes, property id: %s, bizID: %d, exists: %v, rid: %s", attr.PropertyID,
			attr.BizID, exists, kit.Rid)
		if err != nil {
			blog.Errorf("create model attrs failed, property id(%s), err: %s, rid: %s", attr.PropertyID, err, kit.Rid)
			addExceptionFunc(int64(attrIdx), err.(errors.CCErrorCoder), &attr)
			continue
		}

		if exists {
			dataResult.CreateManyInfoResult.Repeated = append(dataResult.CreateManyInfoResult.Repeated,
				metadata.RepeatedDataResult{
					OriginIndex: int64(attrIdx),
					Data:        mapstr.NewFromStruct(attr, "field"),
				})
			continue
		}

		id, err := m.saveTableAttr(kit, attr)
		if err != nil {
			blog.Errorf("failed to save the table attribute(%#v), err: %v, rid: %s", attr, err, kit.Rid)
			addExceptionFunc(int64(attrIdx), err.(errors.CCErrorCoder), &attr)
			continue
		}
		dataResult.CreateManyInfoResult.Created = append(dataResult.CreateManyInfoResult.Created,
			metadata.CreatedDataResult{
				OriginIndex: int64(attrIdx),
				ID:          id,
			})
	}
	return dataResult, nil
}

// CreateModelAttributes create model attributes
func (m *modelAttribute) CreateModelAttributes(kit *rest.Kit, objID string, inputParam metadata.CreateModelAttributes) (
	dataResult *metadata.CreateManyDataResult, err error) {

	if util.InStrArr(forbiddenCreateAttrObjList, objID) {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	dataResult = &metadata.CreateManyDataResult{
		CreateManyInfoResult: metadata.CreateManyInfoResult{
			Created:    []metadata.CreatedDataResult{},
			Repeated:   []metadata.RepeatedDataResult{},
			Exceptions: []metadata.ExceptionResult{},
		},
	}

	if err := m.model.isValid(kit, objID); err != nil {
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
		if err != nil {
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
		if err != nil {
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

	if err := m.model.isValid(kit, objID); err != nil {
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
		if err != nil {
			addExceptionFunc(int64(attrIdx), err.(errors.CCErrorCoder), &attr)
			continue
		}
		attr.OwnerID = kit.SupplierAccount
		if exists {
			cond := mongo.NewCondition()
			cond.Element(&mongo.Eq{Key: metadata.AttributeFieldSupplierAccount, Val: kit.SupplierAccount})
			cond.Element(&mongo.Eq{Key: metadata.AttributeFieldID, Val: existsAttr.ID})

			_, err := m.update(kit, mapstr.NewFromStruct(attr, "field"), cond)
			if err != nil {
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
		if err != nil {
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
	if err != nil {
		blog.Errorf("UpdateModelAttributes failed, failed to convert mapstr(%#v) into a condition object, err: %s, rid: %s", inputParam.Condition, err.Error(), kit.Rid)
		return &metadata.UpdatedCount{}, err
	}

	cnt, err := m.update(kit, inputParam.Data, cond)
	if err != nil {
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
	if err != nil {
		blog.Errorf("UpdateModelAttributesByCondition failed, failed to convert mapstr(%#v) into a condition object, err: %s, rid: %s", inputParam.Condition, err.Error(), kit.Rid)
		return &metadata.UpdatedCount{}, err
	}

	cnt, err := m.update(kit, inputParam.Data, cond)
	if err != nil {
		blog.Errorf("UpdateModelAttributesByCondition failed, failed to update fields (%#v) by condition(%#v), err: %s, rid: %s", inputParam.Data, cond.ToMapStr(), err.Error(), kit.Rid)
		return &metadata.UpdatedCount{}, err
	}

	return &metadata.UpdatedCount{Count: cnt}, nil
}

func assignmentUnchangeableFields(data mapstr.MapStr, dbAttr metadata.Attribute) mapstr.MapStr {

	data[common.BKAppIDField] = dbAttr.BizID
	data[common.BKPropertyIndexField] = dbAttr.PropertyIndex
	data[common.BKPropertyGroupField] = dbAttr.PropertyGroup
	data[common.BKFieldID] = dbAttr.ID
	data[common.BKObjIDField] = dbAttr.ObjectID
	return data
}

// UpdateTableModelAttributes update the attribute content of the form field
func (m *modelAttribute) UpdateTableModelAttributes(kit *rest.Kit, inputParam metadata.UpdateTableOption) error {
	inputParamCond := util.SetModOwner(inputParam.Condition.ToMapInterface(), kit.SupplierAccount)

	filter := metadata.QueryCondition{Condition: inputParamCond}
	attrs, err := m.searchWithSort(kit, filter)
	if err != nil {
		blog.Errorf("failed to search the attrs of the model(%+v), err: %s, rid: %s", inputParam, err, kit.Rid)
		return err
	}

	length := len(attrs)
	if length == 0 || length > 1 {
		blog.Errorf("attrs of the model length(%d) error, filter: %+v, err: %s, rid: %s", filter, length, err, kit.Rid)
		return err
	}

	if len(inputParam.CreateData.Data) > 0 {
		if err := m.model.isValid(kit, inputParam.CreateData.ObjID); err != nil {
			blog.Errorf("validate model(%s) failed, err: %v, rid: %s", inputParam.CreateData.ObjID, err, kit.Rid)
			return err
		}
		for _, attr := range inputParam.CreateData.Data {
			if attr.IsPre {
				if attr.PropertyID == common.BKInstNameField {
					lang := m.language.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header))
					attr.PropertyName = util.FirstNotEmptyString(lang.Language("common_property_"+attr.PropertyID),
						attr.PropertyName, attr.PropertyID)
				}
			}

			attr.OwnerID = kit.SupplierAccount
			_, exists, err := m.isExists(kit, attr.ObjectID, attr.PropertyID, attr.BizID)
			blog.V(5).Infof("table model attributes, property id: %s, bizID: %d, exists: %v, rid: %s", attr.PropertyID,
				attr.BizID, exists, kit.Rid)
			if err != nil {
				blog.Errorf("create model attr failed, propertyID(%s), err: %s, rid: %s", attr.PropertyID, err, kit.Rid)
				continue
			}

			if exists {
				continue
			}
			if err = m.saveTableAttrCheck(kit, attr); err != nil {
				return err
			}
		}

		// inputParam.UpdateData 中的option 转化成header default
		hOp, ok := inputParam.UpdateData["option"].(map[string]interface{})
		if !ok {
			return err
		}
		header := new(metadata.TableAttributesOption)
		if err := mapstruct.Decode2Struct(hOp, header); err != nil {
			return err
		}

		for _, data := range inputParam.CreateData.Data {
			d, ok := data.Option.(map[string]interface{})
			if !ok {
				return err
			}
			dataTbale := new(metadata.TableAttributesOption)
			if err := mapstruct.Decode2Struct(d, dataTbale); err != nil {
				return err
			}

			header.Header = append(header.Header, dataTbale.Header...)

		}
		inputParam.UpdateData[metadata.AttributeFieldOption] = header
	}

	inputParam.UpdateData = assignmentUnchangeableFields(inputParam.UpdateData, attrs[0])

	cond, err := mongo.NewConditionFromMapStr(inputParamCond)
	if err != nil {
		blog.Errorf("parse condition failed, err: %v, cond: %+v, rid: %s", err, inputParam.Condition, kit.Rid)
		return err
	}

	if err = m.unsetTableInstAttr(kit, inputParam.UpdateData, attrs[0]); err != nil {
		return err
	}

	if err := m.updateTableAttr(kit, inputParam.UpdateData, cond); err != nil {
		blog.Errorf("failed to update fields (%#v) by condition(%#v), err: %v, rid: %s",
			inputParam.UpdateData, cond.ToMapStr(), err, kit.Rid)
		return err
	}

	return nil
}

// unsetTableInstAttr unset instance attributes
func (m *modelAttribute) unsetTableInstAttr(kit *rest.Kit, data mapstr.MapStr, attr metadata.Attribute) error {
	// get deleted attributes
	attrOpt, err := metadata.ParseTableAttrOption(attr.Option)
	if err != nil {
		blog.Errorf("parse attribute option failed, err: %v, option: %+v, rid: %s", err, attr.Option, kit.Rid)
		return err
	}

	dataOpt, err := metadata.ParseTableAttrOption(data[common.BKOptionField])
	if err != nil {
		blog.Errorf("parse data option failed, err: %v， data: %+v, rid: %s", err, data, kit.Rid)
		return err
	}

	dataMap := make(map[string]struct{})
	for _, attr := range dataOpt.Header {
		dataMap[attr.PropertyID] = struct{}{}
	}

	deletedAttr := make([]string, 0)
	for _, attr := range attrOpt.Header {
		if _, exists := dataMap[attr.PropertyID]; !exists {
			deletedAttr = append(deletedAttr, attr.PropertyID)
		}
	}

	if len(deletedAttr) == 0 {
		return nil
	}

	// get quoted relation
	quoteCond := mapstr.MapStr{
		common.BKSrcModelField:   attr.ObjectID,
		common.BKPropertyIDField: attr.PropertyID,
	}
	quoteCond = util.SetQueryOwner(quoteCond, kit.SupplierAccount)

	quoteRel := new(metadata.ModelQuoteRelation)
	err = mongodb.Client().Table(common.BKTableNameModelQuoteRelation).Find(quoteCond).One(kit.Ctx, &quoteRel)
	if err != nil {
		blog.Errorf("get model quote relations failed, err: %v, filter: %+v, rid: %v", err, quoteCond, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	// drop instance columns
	instTable := common.GetInstTableName(quoteRel.DestModel, kit.SupplierAccount)

	existCond := make([]map[string]interface{}, len(deletedAttr))
	for index, field := range deletedAttr {
		existCond[index] = map[string]interface{}{
			field: map[string]interface{}{common.BKDBExists: true},
		}
	}
	instCond := util.SetModOwner(mapstr.MapStr{common.BKDBOR: existCond}, kit.SupplierAccount)

	if err = m.dropColumns(kit, quoteRel.DestModel, instTable, instCond, deletedAttr); err != nil {
		blog.Errorf("drop instance table attributes failed, err: %v, attr: %+v, rid: %s", err, deletedAttr, kit.Rid)
		return err
	}

	return nil
}

func (m *modelAttribute) updateTableAttr(kit *rest.Kit, data mapstr.MapStr, cond universalsql.Condition) error {

	if len(data) == 0 {
		return nil
	}

	err := m.checkTableAttrUpdate(kit, data, cond)
	if err != nil {
		blog.Errorf("checkUpdate error. data: %+v, cond: %+v, err: %v, rid:%s", data, cond, err, kit.Rid)
		return err
	}

	_, err = mongodb.Client().Table(common.BKTableNameObjAttDes).UpdateMany(kit.Ctx, cond.ToMapStr(), data)
	if err != nil {
		blog.Errorf("database operation is failed, error: %v, rid: %s", err, kit.Rid)
		return err
	}

	return err
}

// DeleteModelAttributes TODO
func (m *modelAttribute) DeleteModelAttributes(kit *rest.Kit, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	if err := m.model.isValid(kit, objID); nil != err {
		blog.Errorf("request(%s): it is failed to check if the model(%s) is valid, error info is %s", kit.Rid, objID, err.Error())
		return &metadata.DeletedCount{}, err
	}

	cond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), kit.SupplierAccount))
	if err != nil {
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

// SearchModelAttributesByCondition query for model attributes that do not contain table types
func (m *modelAttribute) SearchModelAttributesByCondition(kit *rest.Kit, inputParam metadata.QueryCondition) (
	*metadata.QueryModelAttributeDataResult, error) {

	dataResult := &metadata.QueryModelAttributeDataResult{
		Info: []metadata.Attribute{},
	}
	// NOTICE: exclude table attributes
	inputParam.Condition = map[string]interface{}{
		common.BKDBAND: []map[string]interface{}{
			inputParam.Condition,
			{
				common.BKPropertyTypeField: mapstr.MapStr{
					common.BKDBNE: common.FieldTypeInnerTable,
				},
			},
		},
	}
	inputParam.Condition = util.SetQueryOwner(inputParam.Condition, kit.SupplierAccount)

	attrResult, err := m.searchWithSort(kit, inputParam)
	if err != nil {
		blog.Errorf("failed to search the attributes of the model(%+v), err: %v, rid", inputParam, err, kit.Rid)
		return &metadata.QueryModelAttributeDataResult{}, err
	}

	dataResult.Count = int64(len(attrResult))
	dataResult.Info = attrResult
	return dataResult, nil
}

// SearchModelAttrsWithTableByCondition query includes table field model properties.
func (m *modelAttribute) SearchModelAttrsWithTableByCondition(kit *rest.Kit, inputParam metadata.QueryCondition) (
	*metadata.QueryModelAttributeDataResult, error) {

	inputParam.Condition = util.SetQueryOwner(inputParam.Condition, kit.SupplierAccount)

	dataResult := &metadata.QueryModelAttributeDataResult{
		Info: []metadata.Attribute{},
	}
	attrs, err := m.searchWithSort(kit, inputParam)
	if err != nil {
		blog.Errorf("failed to search the attrs with table of the model(%+v), err: %v, rid: %s",
			inputParam, err, kit.Rid)
		return dataResult, err
	}

	dataResult.Count = int64(len(attrs))
	dataResult.Info = attrs
	return dataResult, nil
}
