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
	"context"
	"encoding/json"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/types"
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
	OwnerID, ObjectID string
	attr              metadata.Attribute
	isNew             bool
	params            types.ContextParams
	clientSet         apimachinery.ClientSetInterface
}

func (a *attribute) Attribute() *metadata.Attribute {
	return &a.attr
}
func (a *attribute) SetAttribute(attr metadata.Attribute) {
	a.attr = attr
	a.attr.OwnerID = a.OwnerID
	a.attr.ObjectID = a.ObjectID
}

func (a *attribute) IsMainlineField() bool {
	return a.attr.PropertyID == common.BKChildStr
}

func (a *attribute) searchObjects(objID string) ([]metadata.Object, error) {
	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(a.params.SupplierAccount).Field(common.BKObjIDField).Eq(objID)

	input := metadata.QueryCondition{
		Condition: cond.ToMapStr(),
	}
	rsp, err := a.clientSet.CoreService().Model().ReadModel(context.Background(), a.params.Header, &input)
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, a.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", objID, rsp.ErrMsg)
		return nil, a.params.Err.Error(rsp.Code)
	}

	models := []metadata.Object{}
	for index := range rsp.Data.Info {
		model := metadata.Object{}
		if err := rsp.Data.Info[index].Spec.ToStructByTag(&model, "field"); err != nil {
			return nil, a.params.Err.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		models = append(models, model)
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

	if a.attr.PropertyID == common.BKChildStr || a.attr.PropertyID == common.BKInstParentStr {
		return nil
	}

	if !isUpdate || data.Exists(metadata.AttributeFieldPropertyType) {
		if _, err := a.FieldValid.Valid(a.params, data, metadata.AttributeFieldPropertyType); nil != err {
			return err
		}
	}

	if !isUpdate || data.Exists(metadata.AttributeFieldPropertyID) {
		val, err := a.FieldValid.Valid(a.params, data, metadata.AttributeFieldPropertyID)
		if nil != err {
			return err
		}
		if err = a.FieldValid.ValidID(a.params, val); nil != err {
			return err
		}
	}

	if !isUpdate || data.Exists(metadata.AttributeFieldPropertyName) {
		val, err := a.FieldValid.Valid(a.params, data, metadata.AttributeFieldPropertyName)
		if nil != err {
			return err
		}
		if err = a.FieldValid.ValidNameWithRegex(a.params, val); nil != err {
			return err
		}
	}

	if !isUpdate || data.Exists(metadata.AttributeFieldOption) {
		propertyType, err := data.String(metadata.AttributeFieldPropertyType)
		if nil != err {
			return a.params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
		}

		if option, exists := data.Get(metadata.AttributeFieldOption); exists && (propertyType == common.FieldTypeInt || propertyType == common.FieldTypeEnum) {
			if err := util.ValidPropertyOption(propertyType, option, a.params.Err); nil != err {
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
	a.attr.OwnerID = a.params.SupplierAccount
	exists, err := a.IsExists()
	if nil != err {
		return err
	}

	if exists {
		return a.params.Err.Error(common.CCErrCommDuplicateItem)
	}

	// create a new record
	input := metadata.CreateModelAttributes{Attributes: []metadata.Attribute{a.attr}}
	rsp, err := a.clientSet.CoreService().Model().CreateModelAttrs(context.Background(), a.params.Header, a.ObjectID, &input)
	if nil != err {
		blog.Errorf("faield to request the object controller, the error info is %s", err.Error())
		return err
	}

	if common.CCSuccess != rsp.Code {
		return err
	}

	for _, id := range rsp.Data.Created {
		a.attr.ID = int64(id.ID)
	}

	return nil
}

func (a *attribute) Update(data mapstr.MapStr) error {

	data.Remove(metadata.AttributeFieldPropertyID)
	data.Remove(metadata.AttributeFieldObjectID)
	data.Remove(metadata.AttributeFieldID)

	if err := a.IsValid(true, data); nil != err {
		return err
	}

	a.attr.OwnerID = a.params.SupplierAccount
	exists, err := a.IsExists()
	if nil != err {
		return err
	}

	if exists {
		return a.params.Err.Error(common.CCErrCommDuplicateItem)
	}

	input := metadata.UpdateOption{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(a.attr.ID).ToMapStr(),
		Data:      data,
	}
	rsp, err := a.clientSet.CoreService().Model().UpdateModelAttrs(context.Background(), a.params.Header, a.ObjectID, &input)
	if nil != err {
		blog.Errorf("failed to request object controller, error info is %s", err.Error())
		return err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to update the object attribute(%s), error info is %s", a.attr.PropertyID, rsp.ErrMsg)
		return a.params.Err.Error(common.CCErrTopoObjectAttributeUpdateFailed)
	}

	return nil
}
func (a *attribute) search(cond condition.Condition) ([]metadata.Attribute, error) {

	rsp, err := a.clientSet.CoreService().Model().ReadModelAttr(context.Background(), a.params.Header, a.ObjectID, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("failed to request to object controller, error info is %s", err.Error())
		return nil, err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to query the object controller, error info is %s", err.Error())
		return nil, a.params.Err.Error(common.CCErrTopoObjectAttributeSelectFailed)
	}

	return rsp.Data.Info, nil
}
func (a *attribute) IsExists() (bool, error) {

	// check id
	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(a.params.SupplierAccount)
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
	cond.Field(common.BKOwnerIDField).Eq(a.params.SupplierAccount)
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

	rsp, err := a.clientSet.CoreService().Model().ReadAttributeGroup(context.Background(), a.params.Header, a.attr.ObjectID, metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[model-grp] failed to request the object controller, error info is %s", err.Error())
		return nil, a.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[model-grp] failed to search the group of the object(%s) by the condition (%#v), error info is %s", a.attr.ObjectID, cond.ToMapStr(), rsp.ErrMsg)
		return nil, a.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	if 0 == len(rsp.Data.Info) {
		return CreateGroup(a.params, a.clientSet, []metadata.Group{metadata.Group{GroupID: "default", GroupName: "Default", OwnerID: a.attr.OwnerID, ObjectID: a.attr.ObjectID}})[0], nil
	}

	return CreateGroup(a.params, a.clientSet, rsp.Data.Info)[0], nil // should be one group
}

func (a *attribute) SetSupplierAccount(supplierAccount string) {
	a.attr.OwnerID = supplierAccount
}
