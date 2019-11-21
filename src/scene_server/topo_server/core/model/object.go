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

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	meta "configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// Object model operation interface declaration
type Object interface {
	Operation

	Parse(data mapstr.MapStr) error

	Object() meta.Object
	IsMainlineObject() (bool, error)
	IsCommon() bool

	SetRecordID(id int64)
	GetMainlineParentObject() (Object, error)
	GetMainlineChildObject() (Object, error)

	GetParentObject() ([]ObjectAssoPair, error)
	GetChildObject() ([]ObjectAssoPair, error)

	SetMainlineParentObject(objID string) error

	CreateMainlineObjectAssociation(relateToObjID string) error

	CreateGroup() GroupInterface
	CreateAttribute() AttributeInterface

	GetGroups() ([]GroupInterface, error)
	GetAttributes() ([]AttributeInterface, error)
	GetNonInnerAttributes() ([]AttributeInterface, error)

	CreateUnique() Unique
	GetUniques() ([]Unique, error)

	SetClassification(class Classification)
	GetClassification() (Classification, error)

	SetSupplierAccount(supplierAccount string)
	GetSupplierAccount() string

	ToMapStr() (mapstr.MapStr, error)

	GetInstIDFieldName() string
	GetInstNameFieldName() string
	GetDefaultInstPropertyName() string
	GetObjectType() string
	GetObjectID() string
}

var _ Object = (*object)(nil)

type object struct {
	FieldValid
	obj       meta.Object
	isNew     bool
	params    types.ContextParams
	clientSet apimachinery.ClientSetInterface
}

func (o *object) Object() meta.Object {
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

func (o *object) GetObjectID() string {
	return o.obj.GetObjectID()
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

func (o *object) searchAttributes(cond condition.Condition) ([]AttributeInterface, error) {
	rsp, err := o.clientSet.CoreService().Model().ReadModelAttr(context.Background(), o.params.Header, o.obj.ObjectID, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s, rid: %s", err.Error(), o.params.ReqID)
		return nil, o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("failed to search the object(%s), error info is %s, rid: %s", o.obj.ObjectID, rsp.ErrMsg, o.params.ReqID)
		return nil, o.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	rstItems := make([]AttributeInterface, 0)
	for _, item := range rsp.Data.Info {

		attr := &attribute{
			attr:      item,
			params:    o.params,
			clientSet: o.clientSet,
		}

		rstItems = append(rstItems, attr)

	}

	return rstItems, nil
}

func (o *object) search(cond condition.Condition) ([]meta.Object, error) {
	rsp, err := o.clientSet.CoreService().Model().ReadModel(context.Background(), o.params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s, rid: %s", err.Error(), o.params.ReqID)
		return nil, o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("failed to search the object(%s), error info is %s, rid: %s", o.obj.ObjectID, rsp.ErrMsg, o.params.ReqID)
		return nil, o.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	models := make([]meta.Object, 0)
	for _, info := range rsp.Data.Info {
		models = append(models, info.Spec)
	}

	return models, nil
}

// GetMainlineParentObject get mainline relationship model
// the parent not exactly mean parent in a tree case
func (o *object) GetMainlineParentObject() (Object, error) {
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(o.obj.ObjectID)
	cond.Field(common.AssociationKindIDField).Eq(common.AssociationKindMainline)

	rsp, err := o.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), o.params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[model-obj] failed to request the object controller, error info is %s, rid: %s", err.Error(), o.params.ReqID)
		return nil, err
	}

	for _, asst := range rsp.Data.Info {
		cond := condition.CreateCondition()
		cond.Field(common.BKObjIDField).Eq(asst.AsstObjID)

		rspRst, err := o.search(cond)
		if nil != err {
			blog.Errorf("[model-obj] failed to search the object(%s)'s parent, error info is %s, rid: %s", asst.ObjectID, err.Error(), o.params.ReqID)
			return nil, err
		}

		objItems := CreateObject(o.params, o.clientSet, rspRst)
		for _, item := range objItems {
			// only one parent in the main-line
			return item, nil
		}

	}

	return nil, io.EOF
}

func (o *object) GetMainlineChildObject() (Object, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKAsstObjIDField).Eq(o.obj.ObjectID)
	cond.Field(common.AssociationKindIDField).Eq(common.AssociationKindMainline)

	rsp, err := o.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), o.params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[model-obj] failed to request the object controller, error info is %s, rid: %s", err.Error(), o.params.ReqID)
		return nil, err
	}

	for _, asst := range rsp.Data.Info {
		cond := condition.CreateCondition()
		cond.Field(common.BKObjIDField).Eq(asst.ObjectID)
		rspRst, err := o.search(cond)
		if nil != err {
			blog.Errorf("[model-obj] failed to search the object(%s)'s child, error info is %s, rid: %s", asst.ObjectID, err.Error(), o.params.ReqID)
			return nil, err
		}

		objItems := CreateObject(o.params, o.clientSet, rspRst)
		if len(objItems) > 1 {
			blog.Errorf("[model-obj] get multiple(%d) children for object(%s), rid: %s", len(objItems), asst.ObjectID, o.params.ReqID)
		}
		for _, item := range objItems {
			// only one child in the main-line
			return item, nil
		}
	}

	return nil, io.EOF
}

func (o *object) searchAssoObjects(isNeedChild bool, cond condition.Condition) ([]ObjectAssoPair, error) {
	rsp, err := o.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), o.params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[model-obj] failed to request the object controller, error info is %s, rid: %s", err.Error(), o.params.ReqID)
		return nil, err
	}

	pair := make([]ObjectAssoPair, 0)
	for _, asst := range rsp.Data.Info {
		cond := condition.CreateCondition()
		if isNeedChild {
			cond.Field(metadata.ModelFieldObjectID).Eq(asst.AsstObjID)
		} else {
			cond.Field(metadata.ModelFieldObjectID).Eq(asst.ObjectID)
		}
		rspRst, err := o.search(cond)
		if nil != err {
			blog.Errorf("[model-obj] failed to search the object(%s)'s parent, error info is %s, rid: %s", asst.ObjectID, err.Error(), o.params.ReqID)
			return nil, err
		}

		if len(rspRst) == 0 {
			blog.Errorf("search asso object, but can not found object with cond: %v, rid: %s", cond.ToMapStr(), o.params.ReqID)
			return nil, fmt.Errorf("can not found object %v", cond.ToMapStr())
		}

		pair = append(pair, ObjectAssoPair{
			Object:      CreateObject(o.params, o.clientSet, rspRst)[0],
			Association: asst,
		})

	}

	return pair, nil
}

func (o *object) GetParentObject() ([]ObjectAssoPair, error) {

	cond := condition.CreateCondition()
	cond.Field(meta.AssociationFieldSupplierAccount).Eq(o.params.SupplierAccount)
	cond.Field(meta.AssociationFieldAssociationObjectID).Eq(o.obj.ObjectID)

	return o.searchAssoObjects(false, cond)
}

func (o *object) GetChildObject() ([]ObjectAssoPair, error) {
	cond := condition.CreateCondition()
	cond.Field(meta.AssociationFieldSupplierAccount).Eq(o.params.SupplierAccount)
	cond.Field(meta.AssociationFieldObjectID).Eq(o.obj.ObjectID)

	return o.searchAssoObjects(true, cond)
}

func (o *object) SetMainlineParentObject(relateToObjID string) error {
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(o.obj.ObjectID)
	cond.Field(common.AssociationKindIDField).Eq(common.AssociationKindMainline)

	resp, err := o.clientSet.CoreService().Association().DeleteModelAssociation(context.Background(), o.params.Header, &metadata.DeleteOption{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("update mainline object[%s] association to %s, search object association failed, err: %v, rid: %s", o.params.ReqID,
			o.obj.ObjectID, relateToObjID, err)
		return o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !resp.Result {
		blog.Errorf("update mainline object[%s] association to %s, search object association failed, err: %v, rid: %s", o.params.ReqID,
			o.obj.ObjectID, relateToObjID, resp.ErrMsg)
		return o.params.Err.Errorf(resp.Code, resp.ErrMsg)
	}
	return o.CreateMainlineObjectAssociation(relateToObjID)
}

func (o *object) generateObjectAssociatioinID(srcObjID, asstID, destObjID string) string {
	return fmt.Sprintf("%s_%s_%s", srcObjID, asstID, destObjID)
}

func (o *object) CreateMainlineObjectAssociation(relateToObjID string) error {
	objAsstID := o.generateObjectAssociatioinID(o.obj.ObjectID, common.AssociationKindMainline, relateToObjID)
	defined := false
	association := meta.Association{
		OwnerID:              o.params.SupplierAccount,
		AssociationName:      objAsstID,
		AssociationAliasName: objAsstID,
		ObjectID:             o.obj.ObjectID,
		// related to it's parent object id
		AsstObjID:  relateToObjID,
		AsstKindID: common.AssociationKindMainline,
		Mapping:    metadata.OneToOneMapping,
		OnDelete:   metadata.NoAction,
		IsPre:      &defined,
	}

	result, err := o.clientSet.CoreService().Association().CreateMainlineModelAssociation(context.Background(), o.params.Header, &metadata.CreateModelAssociation{Spec: association})
	if err != nil {
		blog.Errorf("[model-obj] create mainline object association failed, err: %v, rid: %s", err, o.params.ReqID)
		return err
	}

	if result.Code != common.CCSuccess {
		blog.Errorf("[model-obj] create mainline object association failed, err: %s, rid: %s", result.ErrMsg, o.params.ReqID)
		return o.params.Err.Error(result.Code)
	}

	return nil
}

func (o *object) IsExists() (bool, error) {

	// check id
	cond := condition.CreateCondition()
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
	cond.Field(common.BKObjIDField).Eq(o.obj.ObjectID)
	cond.Field(o.GetInstNameFieldName()).Eq(o.obj.ObjectName)
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
func (o *object) IsValid(isUpdate bool, data mapstr.MapStr) error {

	if !isUpdate || data.Exists(metadata.ModelFieldObjectID) {
		val, err := o.FieldValid.Valid(o.params, data, metadata.ModelFieldObjectID)
		if nil != err {
			blog.Errorf("[model-obj] failed to valid the object id(%s), rid: %s", metadata.ModelFieldObjectID, o.params.ReqID)
			return o.params.Err.New(common.CCErrCommParamsIsInvalid, metadata.ModelFieldObjectID+" "+err.Error())
		}

		if err = o.FieldValid.ValidID(o.params, val); nil != err {
			blog.Errorf("[model-obj] failed to valid the object id(%s), rid: %s", metadata.ModelFieldObjectID, o.params.ReqID)
			return o.params.Err.New(common.CCErrCommParamsIsInvalid, metadata.ModelFieldObjectID+" "+err.Error())
		}
	}

	if !isUpdate || data.Exists(metadata.ModelFieldObjectName) {
		val, err := o.FieldValid.Valid(o.params, data, metadata.ModelFieldObjectName)
		if nil != err {
			blog.Errorf("[model-obj] failed to valid the object name(%s), rid: %s", metadata.ModelFieldObjectName, o.params.ReqID)
			return o.params.Err.New(common.CCErrCommParamsIsInvalid, metadata.ModelFieldObjectName+" "+err.Error())
		}
		if err = o.FieldValid.ValidName(o.params, val); nil != err {
			blog.Errorf("[model-obj] failed to valid the object name(%s), rid: %s", metadata.ModelFieldObjectName, o.params.ReqID)
			return o.params.Err.New(common.CCErrCommParamsIsInvalid, metadata.ModelFieldObjectName+" "+err.Error())
		}
	}

	if !isUpdate || data.Exists(metadata.ModelFieldObjCls) {
		if _, err := o.FieldValid.Valid(o.params, data, metadata.ModelFieldObjCls); nil != err {
			return err
		}
	}

	if !isUpdate && !o.IsCommon() {
		return o.params.Err.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("'%s' the built-in object id, please use a new one", o.obj.ObjectID))
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
		return o.params.Err.Errorf(common.CCErrCommDuplicateItem, o.obj.ObjectID+"/"+o.obj.ObjectName)
	}

	if o.obj.ObjIcon == "" {
		return o.params.Err.Errorf(common.CCErrCommParamsNeedSet, common.BKObjIconField)
	}

	rsp, err := o.clientSet.CoreService().Model().CreateModel(context.Background(), o.params.Header, &metadata.CreateModel{Spec: o.obj})
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s, rid: %s", err.Error(), o.params.ReqID)
		return o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("failed to search the object(%s), error info is %s, rid: %s", o.obj.ObjectID, rsp.ErrMsg, o.params.ReqID)
		return o.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	o.obj.ID = int64(rsp.Data.Created.ID)

	return nil
}

func (o *object) Update(data mapstr.MapStr) error {

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
		return o.params.Err.Errorf(common.CCErrCommDuplicateItem, o.obj.ObjectName)
	}

	// update action
	cond := condition.CreateCondition()
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
		input := metadata.UpdateOption{
			Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(item.ID).ToMapStr(),
			Data:      data,
		}
		rsp, err := o.clientSet.CoreService().Model().UpdateModel(context.Background(), o.params.Header, &input)
		if nil != err {
			blog.Errorf("failed to request the object controller, error info is %s, rid: %s", err.Error(), o.params.ReqID)
			return o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("failed to search the object(%s), error info is %s, rid: %s", o.obj.ObjectID, rsp.ErrMsg, o.params.ReqID)
			return o.params.Err.New(rsp.Code, rsp.ErrMsg)
		}
	}
	return nil
}

func (o *object) Parse(data mapstr.MapStr) error {
	tmp, err := data.ToJSON()
	if err != nil {
		return err
	}

	if err = json.Unmarshal(tmp, &o.obj); err != nil {
		return err
	}

	// err = mapstr.SetValueToStructByTags(&o.obj, data)
	// if nil != err {
	// 	return err
	// }

	return nil
}

func (o *object) ToMapStr() (mapstr.MapStr, error) {
	rst := mapstr.SetValueToMapStrByTags(&o.obj)
	return rst, nil
}

func (o *object) Save(data mapstr.MapStr) error {

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

	if o.obj.ObjIcon == "" {
		return o.params.Err.Errorf(common.CCErrCommParamsNeedSet, common.BKObjIconField)
	}

	return o.Create()

}

func (o *object) CreateGroup() GroupInterface {
	return NewGroup(o.params, o.clientSet)
}

func (o *object) CreateUnique() Unique {
	return &unique{
		params:    o.params,
		clientSet: o.clientSet,
		data: meta.ObjectUnique{
			OwnerID: o.obj.OwnerID,
			ObjID:   o.obj.ObjectID,
		},
	}
}

func (o *object) GetUniques() ([]Unique, error) {
	cond := condition.CreateCondition().Field(common.BKObjIDField).Eq(o.obj.ObjectID)
	rsp, err := o.clientSet.CoreService().Model().ReadModelAttrUnique(context.Background(), o.params.Header, metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s, rid: %s", err.Error(), o.params.ReqID)
		return nil, o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("failed to search the object(%s), error info is %s, rid: %s", o.obj.ObjectID, rsp.ErrMsg, o.params.ReqID)
		return nil, o.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	rstItems := make([]Unique, 0)
	for _, item := range rsp.Data.Info {
		grp := &unique{
			data:      item,
			params:    o.params,
			clientSet: o.clientSet,
		}
		rstItems = append(rstItems, grp)
	}

	return rstItems, nil
}

func (o *object) CreateAttribute() AttributeInterface {
	return &attribute{
		params:    o.params,
		clientSet: o.clientSet,
		attr:      meta.Attribute{},
	}
}

func (o *object) GetNonInnerAttributes() ([]AttributeInterface, error) {

	cond := condition.CreateCondition()
	cond.Field(meta.AttributeFieldObjectID).Eq(o.obj.ObjectID)
	cond.Field(meta.AttributeFieldSupplierAccount).Eq(o.params.SupplierAccount)
	cond.Field(meta.AttributeFieldIsSystem).NotEq(true)
	cond.Field(meta.AttributeFieldIsAPI).NotEq(true)
	return o.searchAttributes(cond)
}

func (o *object) GetAttributes() ([]AttributeInterface, error) {

	cond := condition.CreateCondition()
	cond.Field(meta.AttributeFieldObjectID).Eq(o.obj.ObjectID)
	cond.Field(meta.AttributeFieldSupplierAccount).Eq(o.params.SupplierAccount)
	return o.searchAttributes(cond)
}

func (o *object) GetGroups() ([]GroupInterface, error) {

	cond := condition.CreateCondition()
	cond.Field(meta.GroupFieldObjectID).Eq(o.obj.ObjectID)
	cond.Field(meta.GroupFieldSupplierAccount).Eq(o.params.SupplierAccount)

	rsp, err := o.clientSet.CoreService().Model().ReadAttributeGroup(context.Background(), o.params.Header, o.obj.ObjectID, metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s, rid: %s", err.Error(), o.params.ReqID)
		return nil, o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("failed to search the object(%s), error info is %s, rid: %s", o.obj.ObjectID, rsp.ErrMsg, o.params.ReqID)
		return nil, o.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	rstItems := make([]GroupInterface, 0)
	for _, item := range rsp.Data.Info {
		grp := NewGroup(o.params, o.clientSet)
		grp.SetGroup(item)
		rstItems = append(rstItems, grp)
	}

	return rstItems, nil
}

func (o *object) SetClassification(class Classification) {
	o.obj.ObjCls = class.Classify().ClassificationID
}

func (o *object) GetClassification() (Classification, error) {

	cond := condition.CreateCondition()
	cond.Field(meta.ClassFieldClassificationID).Eq(o.obj.ObjCls)

	rsp, err := o.clientSet.CoreService().Model().ReadModelClassification(context.Background(), o.params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s, rid: %s", err.Error(), o.params.ReqID)
		return nil, o.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("failed to search the object(%s), error info is %s, rid: %s", o.obj.ObjectID, rsp.ErrMsg, o.params.ReqID)
		return nil, o.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	for _, item := range rsp.Data.Info {

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

func (o *object) SetSupplierAccount(supplierAccount string) {
	o.obj.OwnerID = supplierAccount
}

func (o *object) GetSupplierAccount() string {
	return o.obj.OwnerID
}
