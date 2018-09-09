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
	"fmt"
	"io"

	"configcenter/src/common/metadata"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"

	"configcenter/src/scene_server/topo_server/core/types"
)

// Object model operation interface declaration
type Object interface {
	Operation

	Parse(data frtypes.MapStr) (*meta.Object, error)

	Origin() meta.Object
	IsMainlineObject() (bool, error)
	IsCommon() bool

	SetRecordID(id int64)
	GetRecordID() int64
	GetMainlineParentObject() (Object, error)
	GetMainlineChildObject() (Object, error)

	GetChildObjectByFieldID(fieldID string) ([]Object, error)
	GetParentObject() ([]Object, error)
	GetChildObject() ([]Object, error)

	SetMainlineParentObject(objID string) error
	SetMainlineChildObject(objID string) error

	CreateGroup() Group
	CreateAttribute() Attribute

	GetGroups() ([]Group, error)
	GetAttributes() ([]Attribute, error)
	GetAttributesExceptInnerFields() ([]Attribute, error)

	SetClassification(class Classification)
	GetClassification() (Classification, error)

	SetIcon(objectIcon string)
	GetIcon() string

	SetID(objectID string)
	GetID() string

	SetName(objectName string)
	GetName() string

	SetIsPre(isPre bool)
	GetIsPre() bool

	SetIsPaused(isPaused bool)
	GetIsPaused() bool

	SetPosition(position string)
	GetPosition() string

	SetSupplierAccount(supplierAccount string)
	GetSupplierAccount() string

	SetDescription(description string)
	GetDescription() string

	SetCreator(creator string)
	GetCreator() string

	SetModifier(modifier string)
	GetModifier() string

	ToMapStr() (frtypes.MapStr, error)

	GetInstIDFieldName() string
	GetInstNameFieldName() string
	GetDefaultInstPropertyName() string
	GetObjectType() string
}

var _ Object = (*object)(nil)

type object struct {
	FieldValid
	obj       meta.Object
	isNew     bool
	params    types.ContextParams
	clientSet apimachinery.ClientSetInterface
}

func (o *object) Origin() meta.Object {
	return o.obj
}

func (o *object) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.obj)
}

func (o *object) GetDefaultInstPropertyName() string {
	return o.obj.GetDefaultInstPropertyName()
}
func (o *object) GetInstIDFieldName() string {
	return o.obj.GetInstIDFieldName()
}

func (o *object) GetInstNameFieldName() string {
	return o.obj.GetInstNameFieldName()
}

func (o *object) GetObjectType() string {
	return o.obj.GetObjectType()

}
func (o *object) IsCommon() bool {
	return o.obj.IsCommon()
}

func (o *object) IsMainlineObject() (bool, error) {
	attrs, err := o.GetAttributes()
	if nil != err {
		return false, err
	}

	for _, att := range attrs {
		if att.IsMainlineField() {
			return true, nil
		}
	}

	return false, nil
}

func (o *object) searchAttributes(cond condition.Condition) ([]Attribute, error) {

	rsp, err := o.clientSet.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), o.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", o.obj.ObjectID, rsp.ErrMsg)
		return nil, o.params.Err.Error(rsp.Code)
	}

	rstItems := make([]Attribute, 0)
	for _, item := range rsp.Data {

		attr := &attribute{
			attr:      item,
			params:    o.params,
			clientSet: o.clientSet,
		}

		// reset the group name
		grp, err := attr.GetGroup()
		if nil != err {
			blog.Errorf("[model-obj] failed to get the attribute group info , error info is %s", err.Error())
			return nil, err
		}
		attr.SetGroup(grp)

		rstItems = append(rstItems, attr)

	}

	return rstItems, nil
}

func (o *object) search(cond condition.Condition) ([]meta.Object, error) {

	rsp, err := o.clientSet.ObjectController().Meta().SelectObjects(context.Background(), o.params.Header, cond.ToMapStr())

	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", o.obj.ObjectID, rsp.ErrMsg)
		return nil, o.params.Err.Error(rsp.Code)
	}

	return rsp.Data, nil

}

func (o *object) GetMainlineParentObject() (Object, error) {
	cond := condition.CreateCondition()
	cond.Field(meta.AssociationFieldSupplierAccount).Eq(o.params.SupplierAccount)
	cond.Field(meta.AssociationFieldObjectID).Eq(o.obj.ObjectID)
	cond.Field(meta.AssociationFieldObjectAttributeID).Eq(common.BKChildStr)

	rsp, err := o.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), o.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[model-obj] failed to request the object controller, error info is %s", err.Error())
		return nil, err
	}

	for _, asst := range rsp.Data {

		cond := condition.CreateCondition()
		cond.Field(common.BKOwnerIDField).Eq(o.params.SupplierAccount)
		cond.Field(metadata.ModelFieldObjectID).Eq(asst.AsstObjID)

		rspRst, err := o.search(cond)
		if nil != err {
			blog.Errorf("[model-obj] failed to search the object(%s)'s parent, error info is %s", asst.ObjectID, err.Error())
			return nil, err
		}

		objItems := CreateObject(o.params, o.clientSet, rspRst)
		for _, item := range objItems { // only one parent in the main-line
			return item, nil
		}

	}

	return nil, io.EOF
}

func (o *object) GetMainlineChildObject() (Object, error) {

	cond := condition.CreateCondition()
	cond.Field(meta.AssociationFieldSupplierAccount).Eq(o.params.SupplierAccount)
	cond.Field(meta.AssociationFieldAssociationObjectID).Eq(o.obj.ObjectID)
	cond.Field(meta.AssociationFieldObjectAttributeID).Eq(common.BKChildStr)

	rsp, err := o.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), o.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[model-obj] failed to request the object controller, error info is %s", err.Error())
		return nil, err
	}

	for _, asst := range rsp.Data {
		cond := condition.CreateCondition()
		cond.Field(common.BKOwnerIDField).Eq(o.params.SupplierAccount)
		cond.Field(metadata.ModelFieldObjectID).Eq(asst.ObjectID)
		rspRst, err := o.search(cond)
		if nil != err {
			blog.Errorf("[model-obj] failed to search the object(%s)'s child, error info is %s", asst.ObjectID, err.Error())
			return nil, err
		}

		objItems := CreateObject(o.params, o.clientSet, rspRst)
		for _, item := range objItems { // only one child in the main-line
			return item, nil
		}
	}

	return nil, io.EOF
}

func (o *object) searchObjects(isNeedChild bool, cond condition.Condition) ([]Object, error) {
	rsp, err := o.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), o.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[model-obj] failed to request the object controller, error info is %s", err.Error())
		return nil, err
	}

	objItems := make([]Object, 0)
	for _, asst := range rsp.Data {
		cond := condition.CreateCondition()
		cond.Field(common.BKOwnerIDField).Eq(o.params.SupplierAccount)
		if isNeedChild {
			cond.Field(metadata.ModelFieldObjectID).Eq(asst.AsstObjID)
		} else {
			cond.Field(metadata.ModelFieldObjectID).Eq(asst.ObjectID)
		}
		rspRst, err := o.search(cond)
		if nil != err {
			blog.Errorf("[model-obj] failed to search the object(%s)'s parent, error info is %s", asst.ObjectID, err.Error())
			return nil, err
		}

		objItems = append(objItems, CreateObject(o.params, o.clientSet, rspRst)...)

	}

	return objItems, nil
}
func (o *object) GetChildObjectByFieldID(fieldID string) ([]Object, error) {
	cond := condition.CreateCondition()
	cond.Field(meta.AssociationFieldSupplierAccount).Eq(o.params.SupplierAccount)
	cond.Field(meta.AssociationFieldObjectID).Eq(o.obj.ObjectID)
	cond.Field(meta.AssociationFieldObjectAttributeID).Eq(fieldID)

	return o.searchObjects(true, cond)
}
func (o *object) GetParentObject() ([]Object, error) {

	cond := condition.CreateCondition()
	cond.Field(meta.AssociationFieldSupplierAccount).Eq(o.params.SupplierAccount)
	cond.Field(meta.AssociationFieldAssociationObjectID).Eq(o.obj.ObjectID)
	cond.Field(meta.AssociationFieldObjectAttributeID).NotEq(common.BKChildStr)

	return o.searchObjects(false, cond)
}

func (o *object) GetChildObject() ([]Object, error) {
	cond := condition.CreateCondition()
	cond.Field(meta.AssociationFieldSupplierAccount).Eq(o.params.SupplierAccount)
	cond.Field(meta.AssociationFieldObjectID).Eq(o.obj.ObjectID)
	cond.Field(meta.AssociationFieldObjectAttributeID).NotEq(common.BKChildStr)

	return o.searchObjects(true, cond)
}

func (o *object) SetMainlineParentObject(objID string) error {

	cond := condition.CreateCondition()

	cond.Field(meta.AssociationFieldSupplierAccount).Eq(o.params.SupplierAccount)
	cond.Field(meta.AssociationFieldObjectID).Eq(o.obj.ObjectID)
	cond.Field(meta.AssociationFieldObjectAttributeID).Eq(common.BKChildStr)

	rsp, err := o.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), o.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[model-obj] failed to request the object controller, error info is %s", err.Error())
		return o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[model-obj] failed to search the main line association, error info is %s", rsp.ErrMsg)
		return o.params.Err.Error(rsp.Code)
	}

	// create
	if 0 == len(rsp.Data) {

		asst := &meta.Association{}
		asst.OwnerID = o.params.SupplierAccount
		asst.ObjectAttID = common.BKChildStr
		asst.ObjectID = o.obj.ObjectID
		asst.AsstObjID = objID

		rsp, err := o.clientSet.ObjectController().Meta().CreateObjectAssociation(context.Background(), o.params.Header, asst)

		if nil != err {
			blog.Errorf("[model-obj] failed to request the object controller, error info is %s", err.Error())
			return o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if common.CCSuccess != rsp.Code {
			blog.Errorf("[model-obj] failed to set the main line association parent, error info is %s", rsp.ErrMsg)
			return o.params.Err.Error(rsp.Code)
		}

		return nil
	}

	// update
	for _, asst := range rsp.Data {

		asst.AsstObjID = objID
		asst.ObjectAttID = common.BKChildStr

		rsp, err := o.clientSet.ObjectController().Meta().UpdateObjectAssociation(context.Background(), asst.ID, o.params.Header, asst.ToMapStr())
		if nil != err {
			blog.Errorf("[model-obj] failed to request object controller, error info is %s", err.Error())
			return err
		}

		if common.CCSuccess != rsp.Code {
			blog.Errorf("[model-obj] failed to update the parent association, error info is %s", rsp.ErrMsg)
			return o.params.Err.Error(rsp.Code)
		}
	}

	return nil
}
func (o *object) SetMainlineChildObject(objID string) error {

	cond := condition.CreateCondition()

	cond.Field(meta.AssociationFieldSupplierAccount).Eq(o.params.SupplierAccount)
	cond.Field(meta.AssociationFieldObjectAttributeID).Eq(common.BKChildStr)
	cond.Field(meta.AssociationFieldAssociationObjectID).Eq(o.obj.ObjectID)

	rsp, err := o.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), o.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[model-obj] failed to request the object controller, error info is %s", err.Error())
		return err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[model-obj] failed to set the main line association, error info is %s", rsp.ErrMsg)
		return o.params.Err.Error(rsp.Code)
	}

	// create
	if 0 == len(rsp.Data) {

		asst := &meta.Association{}
		asst.OwnerID = o.params.SupplierAccount
		asst.ObjectAttID = common.BKChildStr
		asst.ObjectID = objID
		asst.AsstObjID = o.obj.ObjectID

		rsp, err := o.clientSet.ObjectController().Meta().CreateObjectAssociation(context.Background(), o.params.Header, asst)

		if nil != err {
			blog.Errorf("[model-obj] failed to request the object controller, error info is %s", err.Error())
			return o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if common.CCSuccess != rsp.Code {
			blog.Errorf("[model-obj] failed to set the main line association parent, error info is %s", rsp.ErrMsg)
			return o.params.Err.Error(rsp.Code)
		}

		return nil
	}

	// update
	for _, asst := range rsp.Data { // should be only one item

		asst.ObjectID = objID
		asst.ObjectAttID = common.BKChildStr

		rsp, err := o.clientSet.ObjectController().Meta().UpdateObjectAssociation(context.Background(), asst.ID, o.params.Header, asst.ToMapStr())
		if nil != err {
			blog.Errorf("[model-obj] failed to request object controller, error info is %s", err.Error())
			return err
		}

		if common.CCSuccess != rsp.Code {
			blog.Errorf("[model-obj] failed to update the child association, error info is %s", rsp.ErrMsg)
			return o.params.Err.Error(rsp.Code)
		}
	}

	return nil
}

func (o *object) IsExists() (bool, error) {

	// check id
	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(o.params.SupplierAccount)
	cond.Field(common.BKObjIDField).Eq(o.obj.ObjectID)
	cond.Field(metadata.ModelFieldID).NotIn([]int64{o.obj.ID})

	items, err := o.search(cond)
	if nil != err {
		return false, err
	}

	if 0 != len(items) {
		return true, nil
	}

	// check name
	cond = condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(o.params.SupplierAccount)
	cond.Field(common.BKObjIDField).Eq(o.obj.ObjectName)
	cond.Field(metadata.ModelFieldID).NotIn([]int64{o.obj.ID})

	items, err = o.search(cond)
	if nil != err {
		return false, err
	}

	if 0 != len(items) {
		return true, nil
	}

	return false, nil
}
func (o *object) IsValid(isUpdate bool, data frtypes.MapStr) error {

	if !isUpdate || data.Exists(metadata.ModelFieldObjectID) {
		val, err := o.FieldValid.Valid(o.params, data, metadata.ModelFieldObjectID)
		if nil != err {
			blog.Errorf("[model-obj] failed to valid the object id(%s)", metadata.ModelFieldObjectID)
			return o.params.Err.New(common.CCErrCommParamsIsInvalid, metadata.ModelFieldObjectID+" "+err.Error())
		}

		if err = o.FieldValid.ValidID(o.params, val); nil != err {
			blog.Errorf("[model-obj] failed to valid the object id(%s)", metadata.ModelFieldObjectID)
			return o.params.Err.New(common.CCErrCommParamsIsInvalid, metadata.ModelFieldObjectID+" "+err.Error())
		}
	}

	if !isUpdate || data.Exists(metadata.ModelFieldObjectName) {
		val, err := o.FieldValid.Valid(o.params, data, metadata.ModelFieldObjectName)
		if nil != err {
			blog.Errorf("[model-obj] failed to valid the object name(%s)", metadata.ModelFieldObjectName)
			return o.params.Err.New(common.CCErrCommParamsIsInvalid, metadata.ModelFieldObjectName+" "+err.Error())
		}
		if err = o.FieldValid.ValidName(o.params, val); nil != err {
			blog.Errorf("[model-obj] failed to valid the object name(%s)", metadata.ModelFieldObjectName)
			return o.params.Err.New(common.CCErrCommParamsIsInvalid, metadata.ModelFieldObjectName+" "+err.Error())
		}
	}

	if !isUpdate || data.Exists(metadata.ModelFieldObjCls) {
		if _, err := o.FieldValid.Valid(o.params, data, metadata.ModelFieldObjCls); nil != err {
			return err
		}
	}

	if !isUpdate && !o.IsCommon() {
		return o.params.Err.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("'%s' the built-in object id, please use a new one", o.GetID()))
	}

	return nil
}
func (o *object) Create() error {

	if err := o.IsValid(false, o.obj.ToMapStr()); nil != err {
		return err
	}

	o.obj.OwnerID = o.params.SupplierAccount
	exists, err := o.IsExists()
	if nil != err {
		return err
	}

	if exists {
		return o.params.Err.Error(common.CCErrCommDuplicateItem)
	}

	rsp, err := o.clientSet.ObjectController().Meta().CreateObject(context.Background(), o.params.Header, &o.obj)

	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", o.obj.ObjectID, rsp.ErrMsg)
		return o.params.Err.Error(rsp.Code)
	}

	o.obj.ID = rsp.Data.ID

	return nil
}

func (o *object) Update(data frtypes.MapStr) error {

	data.Remove(metadata.ModelFieldObjectID)
	data.Remove(metadata.ModelFieldID)
	data.Remove(metadata.ModelFieldObjCls)

	if err := o.IsValid(true, data); nil != err {
		return err
	}

	exists, err := o.IsExists()
	if nil != err {
		return err
	}

	if exists {
		return o.params.Err.Error(common.CCErrCommDuplicateItem)
	}

	// update action
	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(o.params.SupplierAccount)
	if 0 != len(o.obj.ObjectID) {
		cond.Field(common.BKObjIDField).Eq(o.obj.ObjectID)
	} else {
		cond.Field(metadata.ModelFieldID).Eq(o.obj.ID)
	}

	items, err := o.search(cond)
	if nil != err {
		return err
	}

	for _, item := range items {

		rsp, err := o.clientSet.ObjectController().Meta().UpdateObject(context.Background(), item.ID, o.params.Header, data)

		if nil != err {
			blog.Errorf("failed to request the object controller, error info is %s", err.Error())
			return o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if common.CCSuccess != rsp.Code {
			blog.Errorf("failed to search the object(%s), error info is %s", o.obj.ObjectID, rsp.ErrMsg)
			return o.params.Err.Error(rsp.Code)
		}
	}
	return nil
}

func (o *object) Parse(data frtypes.MapStr) (*meta.Object, error) {

	err := meta.SetValueToStructByTags(&o.obj, data)
	if nil != err {
		return nil, err
	}

	/*
		if 0 == len(o.obj.ObjectID) {
			return nil, o.params.Err.Errorf(common.CCErrCommParamsNeedSet, meta.ModelFieldObjectID)
		}

		if 0 == len(o.obj.ObjCls) {
			return nil, o.params.Err.Errorf(common.CCErrCommParamsNeedSet, meta.ModelFieldObjCls)
		}
	*/
	return nil, err
}

func (o *object) ToMapStr() (frtypes.MapStr, error) {
	rst := meta.SetValueToMapStrByTags(&o.obj)
	return rst, nil
}

func (o *object) Save(data frtypes.MapStr) error {

	if nil != data {
		if _, err := o.obj.Parse(data); nil != err {
			return err
		}
	}

	if exists, err := o.IsExists(); nil != err {
		return err
	} else if exists {
		if nil != data {
			return o.Update(data)
		}
		return o.Update(o.obj.ToMapStr())
	}

	return o.Create()

}

func (o *object) CreateGroup() Group {
	return &group{
		params:    o.params,
		clientSet: o.clientSet,
		grp: meta.Group{

			OwnerID:  o.obj.OwnerID,
			ObjectID: o.obj.ObjectID,
		},
	}
}

func (o *object) CreateAttribute() Attribute {
	return &attribute{
		params:    o.params,
		clientSet: o.clientSet,
		attr: meta.Attribute{
			OwnerID:  o.obj.OwnerID,
			ObjectID: o.obj.ObjectID,
		},
	}
}

func (o *object) GetAttributesExceptInnerFields() ([]Attribute, error) {

	cond := condition.CreateCondition()
	cond.Field(meta.AttributeFieldObjectID).Eq(o.obj.ObjectID)
	cond.Field(meta.AttributeFieldSupplierAccount).Eq(o.params.SupplierAccount)
	cond.Field(meta.AttributeFieldIsSystem).NotEq(true)
	cond.Field(meta.AttributeFieldIsAPI).NotEq(true)
	return o.searchAttributes(cond)
}

func (o *object) GetAttributes() ([]Attribute, error) {

	cond := condition.CreateCondition()
	cond.Field(meta.AttributeFieldObjectID).Eq(o.obj.ObjectID).Field(meta.AttributeFieldSupplierAccount).Eq(o.params.SupplierAccount)
	return o.searchAttributes(cond)
}

func (o *object) GetGroups() ([]Group, error) {

	cond := condition.CreateCondition()

	cond.Field(meta.GroupFieldObjectID).Eq(o.obj.ObjectID).Field(meta.GroupFieldSupplierAccount).Eq(o.params.SupplierAccount)
	rsp, err := o.clientSet.ObjectController().Meta().SelectGroup(context.Background(), o.params.Header, cond.ToMapStr())

	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", o.obj.ObjectID, rsp.ErrMsg)
		return nil, o.params.Err.Error(rsp.Code)
	}

	rstItems := make([]Group, 0)
	for _, item := range rsp.Data {
		grp := &group{
			grp:       item,
			params:    o.params,
			clientSet: o.clientSet,
		}
		rstItems = append(rstItems, grp)
	}

	return rstItems, nil
}

func (o *object) SetClassification(class Classification) {
	o.obj.ObjCls = class.GetID()
}

func (o *object) GetClassification() (Classification, error) {

	cond := condition.CreateCondition()
	cond.Field(meta.ClassFieldClassificationID).Eq(o.obj.ObjCls)

	rsp, err := o.clientSet.ObjectController().Meta().SelectClassifications(context.Background(), o.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", o.obj.ObjectID, rsp.ErrMsg)
		return nil, o.params.Err.Error(rsp.Code)
	}

	for _, item := range rsp.Data {

		return &classification{
			cls:       item,
			params:    o.params,
			clientSet: o.clientSet,
		}, nil // only one classification
	}

	return nil, fmt.Errorf("invalid classification(%s) for the object(%s)", o.obj.ObjCls, o.obj.ObjectID)
}
func (o *object) SetRecordID(id int64) {
	o.obj.ID = id
}
func (o *object) GetRecordID() int64 {
	return o.obj.ID
}

func (o *object) SetIcon(objectIcon string) {
	o.obj.ObjIcon = objectIcon
}

func (o *object) GetIcon() string {
	return o.obj.ObjIcon
}

func (o *object) SetID(objectID string) {
	o.obj.ObjectID = objectID
}

func (o *object) GetID() string {
	return o.obj.ObjectID
}

func (o *object) SetName(objectName string) {
	o.obj.ObjectName = objectName
}

func (o *object) GetName() string {
	return o.obj.ObjectName
}

func (o *object) SetIsPre(isPre bool) {
	o.obj.IsPre = isPre
}

func (o *object) GetIsPre() bool {
	return o.obj.IsPre
}

func (o *object) SetIsPaused(isPaused bool) {
	o.obj.IsPaused = isPaused
}

func (o *object) GetIsPaused() bool {
	return o.obj.IsPaused
}

func (o *object) SetPosition(position string) {
	o.obj.Position = position
}

func (o *object) GetPosition() string {
	return o.obj.Position
}

func (o *object) SetSupplierAccount(supplierAccount string) {
	o.obj.OwnerID = supplierAccount
}

func (o *object) GetSupplierAccount() string {
	return o.obj.OwnerID
}

func (o *object) SetDescription(description string) {
	o.obj.Description = description
}

func (o *object) GetDescription() string {
	return o.obj.Description
}

func (o *object) SetCreator(creator string) {
	o.obj.Creator = creator
}

func (o *object) GetCreator() string {
	return o.obj.Creator
}

func (o *object) SetModifier(modifier string) {
	o.obj.Modifier = modifier
}

func (o *object) GetModifier() string {
	return o.obj.Modifier
}
