/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import (
	"encoding/json"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// Attribute attribute opeartion interface declaration
type AttributeInterface interface {
	Operation
	Parse(data mapstr.MapStr) error

	Attribute() *metadata.Attribute
	SetAttribute(attr metadata.Attribute)
	IsMainlineField() bool
	ToMapStr() (mapstr.MapStr, error)
}

var _ AttributeInterface = (*attribute)(nil)

// attribute the metadata structure definition of the model attribute
type attribute struct {
	FieldValid
	attr      metadata.Attribute
	isNew     bool
	kit       *rest.Kit
	clientSet apimachinery.ClientSetInterface
}

func (a *attribute) Attribute() *metadata.Attribute {
	return &a.attr
}
func (a *attribute) SetAttribute(attr metadata.Attribute) {
	a.attr = attr
}

func (a *attribute) IsMainlineField() bool {
	return a.attr.PropertyID == common.BKInstParentStr
}

func (a *attribute) searchObjects(objID string) ([]metadata.Object, error) {
	cond := condition.CreateCondition()

	input := metadata.QueryCondition{
		Condition: cond.ToMapStr(),
	}
	rsp, err := a.clientSet.CoreService().Model().ReadModel(a.kit.Ctx, a.kit.Header, &input)
	if nil != err {
		blog.Errorf("failed to request the object controller, err: %s, rid: %s", err.Error(), a.kit.Rid)
		return nil, a.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("failed to search the object(%s), err: %s, rid: %s", objID, rsp.ErrMsg, a.kit.Rid)
		return nil, a.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	models := []metadata.Object{}
	for index := range rsp.Data.Info {
		models = append(models, rsp.Data.Info[index].Spec)
	}
	return models, nil

}

func (a *attribute) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.attr)
}

func (a *attribute) Parse(data mapstr.MapStr) error {
	attr, err := a.attr.Parse(data)
	if nil != err {
		return err
	}

	a.attr = *attr
	if a.attr.IsOnly {
		a.attr.IsRequired = true
	}

	if 0 == len(a.attr.PropertyGroup) {
		a.attr.PropertyGroup = "default"
	}

	return err
}

func (a *attribute) ToMapStr() (mapstr.MapStr, error) {
	rst := mapstr.SetValueToMapStrByTags(&a.attr)
	return rst, nil

}

func (a *attribute) IsValid(isUpdate bool, data mapstr.MapStr) error {

	if a.attr.PropertyID == common.BKInstParentStr {
		return nil
	}

	// check if property type for creation is valid, can't update property type
	if !isUpdate {
		if _, err := a.FieldValid.Valid(a.kit, data, metadata.AttributeFieldPropertyType); nil != err {
			return err
		}
	}

	if !isUpdate || data.Exists(metadata.AttributeFieldPropertyID) {
		val, err := a.FieldValid.Valid(a.kit, data, metadata.AttributeFieldPropertyID)
		if nil != err {
			return err
		}
		if err = a.FieldValid.ValidID(a.kit, val); nil != err {
			return err
		}
	}

	if !isUpdate || data.Exists(metadata.AttributeFieldPropertyName) {
		val, err := a.FieldValid.Valid(a.kit, data, metadata.AttributeFieldPropertyName)
		if nil != err {
			return err
		}
		if err = a.FieldValid.ValidNameWithRegex(a.kit, val); nil != err {
			return err
		}
	}

	// check option validity for creation, update validation is in coreservice cause property type need to be obtained from db
	if !isUpdate {
		propertyType, err := data.String(metadata.AttributeFieldPropertyType)
		if nil != err {
			return a.kit.CCError.New(common.CCErrCommParamsIsInvalid, err.Error())
		}

		option, exists := data.Get(metadata.AttributeFieldOption)
		if exists && a.isPropertyTypeIntEnumList(propertyType) {
			if err := util.ValidPropertyOption(propertyType, option, a.kit.CCError); nil != err {
				return err
			}
		}
	}

	if val, ok := data[metadata.AttributeFieldPlaceHolder]; ok && val != "" {
		if placeholder, ok := val.(string); ok {
			if err := a.FieldValid.ValidPlaceHolder(a.kit, placeholder); nil != err {
				return err
			}
		}
	}

	return nil
}

func (a *attribute) Create() error {

	if err := a.IsValid(false, a.attr.ToMapStr()); nil != err {
		return err
	}

	// check the property id repeated
	a.attr.OwnerID = a.kit.SupplierAccount

	// create a new record
	input := metadata.CreateModelAttributes{Attributes: []metadata.Attribute{a.attr}}
	rsp, err := a.clientSet.CoreService().Model().CreateModelAttrs(a.kit.Ctx, a.kit.Header, a.attr.ObjectID, &input)
	if nil != err {
		blog.ErrorJSON("failed to request coreService to create model attrs, the err: %s, ObjectID: %s, input: %s, rid: %s", err.Error(), a.attr.ObjectID, input, a.kit.Rid)
		return a.kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.ErrorJSON("create model attrs failed, ObjectID: %s, input: %s, rid: %s", a.attr.ObjectID, input, a.kit.Rid)
		return rsp.CCError()
	}

	for _, exception := range rsp.Data.Exceptions {
		return a.kit.CCError.New(int(exception.Code), exception.Message)
	}

	if len(rsp.Data.Repeated) > 0 {
		blog.ErrorJSON("create model attrs failed, the attr is duplicated, ObjectID: %s, input: %s, rid: %s", a.attr.ObjectID, input, a.kit.Rid)
		return a.kit.CCError.CCError(common.CCErrorAttributeNameDuplicated)
	}

	if len(rsp.Data.Created) != 1 {
		blog.ErrorJSON("create model attrs created amount error, ObjectID: %s, input: %s, rid: %s", a.attr.ObjectID, input, a.kit.Rid)
		return a.kit.CCError.CCError(common.CCErrTopoObjectAttributeCreateFailed)
	}
	a.attr.ID = int64(rsp.Data.Created[0].ID)

	return nil
}

func (a *attribute) Update(data mapstr.MapStr) error {

	data.Remove(metadata.AttributeFieldPropertyID)
	data.Remove(metadata.AttributeFieldObjectID)
	data.Remove(metadata.AttributeFieldID)

	if err := a.IsValid(true, data); nil != err {
		return err
	}

	a.attr.OwnerID = a.kit.SupplierAccount
	exists, err := a.IsExists()
	if nil != err {
		return err
	}

	if exists {
		return a.kit.CCError.Errorf(common.CCErrCommDuplicateItem, a.attr.PropertyName)
	}

	input := metadata.UpdateOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(a.attr.ID).ToMapStr(),
		Data:      data,
	}
	rsp, err := a.clientSet.CoreService().Model().UpdateModelAttrs(a.kit.Ctx, a.kit.Header, a.attr.ObjectID, &input)
	if nil != err {
		blog.Errorf("failed to request object controller, err: %s, rid: %s", err.Error(), a.kit.Rid)
		return err
	}

	if !rsp.Result {
		blog.Errorf("failed to update the object attribute(%s), err: %s, rid: %s", a.attr.PropertyID, rsp.ErrMsg, a.kit.Rid)
		return a.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}
	return nil
}
func (a *attribute) search(cond condition.Condition) ([]metadata.Attribute, error) {

	rsp, err := a.clientSet.CoreService().Model().ReadModelAttr(a.kit.Ctx, a.kit.Header, a.attr.ObjectID, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("failed to request to object controller, err: %s, rid: %s", err.Error(), a.kit.Rid)
		return nil, err
	}

	if !rsp.Result {
		blog.Errorf("failed to query the object controller, cond: %#v, err: %s, rid: %s", cond, rsp.ErrMsg, a.kit.Rid)
		return nil, a.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data.Info, nil
}
func (a *attribute) IsExists() (bool, error) {

	// check id
	cond := condition.CreateCondition()
	cond.Field(metadata.AttributeFieldObjectID).Eq(a.attr.ObjectID)
	cond.Field(metadata.AttributeFieldPropertyID).Eq(a.attr.PropertyID)
	cond.Field(metadata.AttributeFieldID).NotIn([]int64{a.attr.ID})

	items, err := a.search(cond)
	if nil != err {
		return false, err
	}

	if 0 != len(items) {
		return true, err
	}

	// ceck nam
	cond = condition.CreateCondition()
	cond.Field(metadata.AttributeFieldObjectID).Eq(a.attr.ObjectID)
	cond.Field(metadata.AttributeFieldPropertyName).Eq(a.attr.PropertyName)
	cond.Field(metadata.AttributeFieldID).NotIn([]int64{a.attr.ID})

	items, err = a.search(cond)
	if nil != err {
		return false, err
	}

	if 0 != len(items) {
		return true, err
	}

	return false, nil
}

func (a *attribute) Save(data mapstr.MapStr) error {

	if nil != data {
		if _, err := a.attr.Parse(data); nil != err {
			return err
		}
	}

	if exists, err := a.IsExists(); nil != err {
		return err
	} else if !exists {
		return a.Create()
	}

	return a.Update(data)
}

func (a *attribute) SetGroup(grp GroupInterface) {
	group := grp.Group()
	a.attr.PropertyGroup = group.GroupID
	a.attr.PropertyGroupName = group.GroupName
}

func (a *attribute) GetGroup() (GroupInterface, error) {

	cond := condition.CreateCondition()
	cond.Field(metadata.GroupFieldGroupID).Eq(a.attr.PropertyGroup)
	cond.Field(metadata.GroupFieldObjectID).Eq(a.attr.ObjectID)

	rsp, err := a.clientSet.CoreService().Model().ReadAttributeGroup(a.kit.Ctx, a.kit.Header, a.attr.ObjectID, metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[model-grp] failed to request the coreservice, err: %s, rid: %s", err.Error(), a.kit.Rid)
		return nil, a.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[model-grp] failed to search the group of the object(%s) by the condition (%#v), err: %s, rid: %s", a.attr.ObjectID, cond.ToMapStr(), rsp.ErrMsg, a.kit.Rid)
		return nil, a.kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	if 0 == len(rsp.Data.Info) {
		return CreateGroup(a.kit, a.clientSet, []metadata.Group{{GroupID: "default", GroupName: "Default", OwnerID: a.attr.OwnerID, ObjectID: a.attr.ObjectID}})[0], nil
	}

	return CreateGroup(a.kit, a.clientSet, rsp.Data.Info)[0], nil // should be one group
}

func (a *attribute) SetSupplierAccount(supplierAccount string) {
	a.attr.OwnerID = supplierAccount
}

func (a *attribute) isPropertyTypeIntEnumList(propertyType string) bool {
	switch propertyType {
	case common.FieldTypeInt, common.FieldTypeEnum, common.FieldTypeList:
		return true
	default:
		return false
	}
}
