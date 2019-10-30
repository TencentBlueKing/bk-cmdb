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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type modelAttribute struct {
	model   *modelManager
	dbProxy dal.RDB
}

func (m *modelAttribute) CreateModelAttributes(ctx core.ContextParams, objID string, inputParam metadata.CreateModelAttributes) (dataResult *metadata.CreateManyDataResult, err error) {

	dataResult = &metadata.CreateManyDataResult{
		CreateManyInfoResult: metadata.CreateManyInfoResult{
			Created:    []metadata.CreatedDataResult{},
			Repeated:   []metadata.RepeatedDataResult{},
			Exceptions: []metadata.ExceptionResult{},
		},
	}

	if err := m.model.isValid(ctx, objID); nil != err {
		blog.Errorf("CreateModelAttributes failed, validate model(%s) failed, err: %s, rid: %s", objID, err.Error(), ctx.ReqID)
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
		if attr.IsPre {
			if attr.PropertyID == common.BKInstNameField {
				attr.PropertyName = util.FirstNotEmptyString(ctx.Lang.Language("common_property_"+attr.PropertyID), attr.PropertyName, attr.PropertyID)
			}
		}

		attr.OwnerID = ctx.SupplierAccount
		_, exists, err := m.isExists(ctx, attr.ObjectID, attr.PropertyID, attr.Metadata)
		blog.V(5).Infof("CreateModelAttributes isExists info. property id:%s, metadata:%#v, exit:%v, rid:%s", attr.PropertyID, attr.Metadata, exists, ctx.ReqID)
		if nil != err {
			blog.Errorf("CreateModelAttributes failed, attribute field propertyID(%s) exists, err: %s, rid: %s", attr.PropertyID, err.Error(), ctx.ReqID)
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
		id, err := m.save(ctx, attr)
		if nil != err {
			blog.Errorf("CreateModelAttributes failed, failed to save the attribute(%#v), err: %s, rid: %s", attr, err.Error(), ctx.ReqID)
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

func (m *modelAttribute) SetModelAttributes(ctx core.ContextParams, objID string, inputParam metadata.SetModelAttributes) (dataResult *metadata.SetDataResult, err error) {

	dataResult = &metadata.SetDataResult{
		Created:    []metadata.CreatedDataResult{},
		Updated:    []metadata.UpdatedDataResult{},
		Exceptions: []metadata.ExceptionResult{},
	}

	if err := m.model.isValid(ctx, objID); nil != err {
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

		existsAttr, exists, err := m.isExists(ctx, attr.ObjectID, attr.PropertyID, attr.Metadata)
		if nil != err {
			addExceptionFunc(int64(attrIdx), err.(errors.CCErrorCoder), &attr)
			continue
		}
		attr.OwnerID = ctx.SupplierAccount
		if exists {
			cond := mongo.NewCondition()
			cond.Element(&mongo.Eq{Key: metadata.AttributeFieldSupplierAccount, Val: ctx.SupplierAccount})
			cond.Element(&mongo.Eq{Key: metadata.AttributeFieldID, Val: existsAttr.ID})

			_, err := m.update(ctx, mapstr.NewFromStruct(attr, "field"), cond)
			if nil != err {
				blog.Errorf("SetModelAttributes failed, failed to update the attribute(%#v) by the condition(%#v), err: %s, rid: %s", attr, cond.ToMapStr(), err.Error(), ctx.ReqID)
				addExceptionFunc(int64(attrIdx), err.(errors.CCErrorCoder), &attr)
				continue
			}
			dataResult.Updated = append(dataResult.Updated, metadata.UpdatedDataResult{
				OriginIndex: int64(attrIdx),
				ID:          uint64(existsAttr.ID),
			})
			continue
		}
		id, err := m.save(ctx, attr)
		if nil != err {
			blog.Errorf("SetModelAttributes failed, failed to save the attribute(%#v), err: %s, rid: %s", attr, err.Error(), ctx.ReqID)
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
func (m *modelAttribute) UpdateModelAttributes(ctx core.ContextParams, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	if err := m.model.isValid(ctx, objID); nil != err {
		blog.Errorf("UpdateModelAttributes failed, validate model(%s) failed, err: %s, rid: %s", objID, err.Error(), ctx.ReqID)
		return &metadata.UpdatedCount{}, err
	}

	cond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("UpdateModelAttributes failed, failed to convert mapstr(%#v) into a condition object, err: %s, rid: %s", inputParam.Condition, err.Error(), ctx.ReqID)
		return &metadata.UpdatedCount{}, err
	}

	cnt, err := m.update(ctx, inputParam.Data, cond)
	if nil != err {
		blog.ErrorJSON("UpdateModelAttributes failed, update attributes failed, model:%s, attributes:%s, condition: %s, err: %s, rid: %s", inputParam.Data, objID, cond, err.Error(), ctx.ReqID)
		return &metadata.UpdatedCount{}, err
	}

	return &metadata.UpdatedCount{Count: cnt}, nil
}

func (m *modelAttribute) UpdateModelAttributesByCondition(ctx core.ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	cond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("UpdateModelAttributesByCondition failed, failed to convert mapstr(%#v) into a condition object, err: %s, rid: %s", inputParam.Condition, err.Error(), ctx.ReqID)
		return &metadata.UpdatedCount{}, err
	}

	cnt, err := m.update(ctx, inputParam.Data, cond)
	if nil != err {
		blog.Errorf("UpdateModelAttributesByCondition failed, failed to update fields (%#v) by condition(%#v), err: %s, rid: %s", inputParam.Data, cond.ToMapStr(), err.Error(), ctx.ReqID)
		return &metadata.UpdatedCount{}, err
	}

	return &metadata.UpdatedCount{Count: cnt}, nil
}

func (m *modelAttribute) DeleteModelAttributes(ctx core.ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	if err := m.model.isValid(ctx, objID); nil != err {
		blog.Errorf("request(%s): it is failed to check if the model(%s) is valid, error info is %s", ctx.ReqID, objID, err.Error())
		return &metadata.DeletedCount{}, err
	}

	cond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert from mapstr(%#v) into a condition object, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.DeletedCount{}, err
	}

	cond.Element(&mongo.Eq{Key: metadata.AttributeFieldSupplierAccount, Val: ctx.SupplierAccount})
	cnt, err := m.delete(ctx, cond)
	return &metadata.DeletedCount{Count: cnt}, err
}

func (m *modelAttribute) SearchModelAttributes(ctx core.ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeDataResult, error) {

	dataResult := &metadata.QueryModelAttributeDataResult{
		Info: []metadata.Attribute{},
	}

	if err := m.model.isValid(ctx, objID); nil != err {
		blog.Errorf("request(%s): it is failed to check if the model(%s) is valid, error info is %s", ctx.ReqID, objID, err.Error())
		return &metadata.QueryModelAttributeDataResult{}, err
	}

	cond, err := mongo.NewConditionFromMapStr(util.SetQueryOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert from mapstr(%#v) into a condition object, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.QueryModelAttributeDataResult{}, err
	}
	attrArr := []string{ctx.SupplierAccount, common.BKDefaultOwnerID}
	cond.Element(&mongo.In{Key: metadata.AttributeFieldSupplierAccount, Val: attrArr})
	cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: objID})
	attrResult, err := m.search(ctx, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to search the attributes of the model(%s), error info is %s", ctx.ReqID, objID, err.Error())
		return &metadata.QueryModelAttributeDataResult{}, err
	}

	dataResult.Count = int64(len(attrResult))
	dataResult.Info = attrResult
	return dataResult, nil
}

func (m *modelAttribute) SearchModelAttributesByCondition(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeDataResult, error) {

	dataResult := &metadata.QueryModelAttributeDataResult{
		Info: []metadata.Attribute{},
	}

	cond, err := mongo.NewConditionFromMapStr(util.SetQueryOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert from mapstr(%#v) into a condition object, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.QueryModelAttributeDataResult{}, err
	}
	attrArr := []string{ctx.SupplierAccount, common.BKDefaultOwnerID}
	cond.Element(&mongo.In{Key: metadata.AttributeFieldSupplierAccount, Val: attrArr})
	attrResult, err := m.search(ctx, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to search the attributes of the model(%+v), error info is %s", ctx.ReqID, cond, err.Error())
		return &metadata.QueryModelAttributeDataResult{}, err
	}

	dataResult.Count = int64(len(attrResult))
	dataResult.Info = attrResult
	return dataResult, nil
}
