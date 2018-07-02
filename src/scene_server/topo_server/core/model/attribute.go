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
	"io"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	metadata "configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// Attribute attribute opeartion interface declaration
type Attribute interface {
	Operation
	Parse(data frtypes.MapStr) (*metadata.Attribute, error)

	SetSupplierAccount(supplierAccount string)
	GetSupplierAccount() string

	GetParentObject() (Object, error)
	GetChildObject() (Object, error)

	SetParentObject(objID string) error
	SetChildObject(objID string) error

	SetObjectID(objectID string)
	GetObjectID() string

	SetID(attributeID string)
	GetID() string

	SetName(attributeName string)
	GetName() string

	SetGroup(grp Group)
	GetGroup() (Group, error)

	SetGroupIndex(attGroupIndex int64)
	GetGroupIndex() int64

	SetUnint(unit string)
	GetUnint() string

	SetPlaceholder(placeHolder string)
	GetPlaceholder() string

	SetIsEditable(isEditable bool)
	GetIsEditable() bool

	SetIsPre(isPre bool)
	GetIsPre() bool

	SetIsReadOnly(isReadOnly bool)
	GetIsReadOnly() bool

	SetIsOnly(isOnly bool)
	GetIsOnly() bool

	SetIsSystem(isSystem bool)
	GetIsSystem() bool

	SetIsAPI(isAPI bool)
	GetIsAPI() bool

	SetType(attributeType string)
	GetType() string

	SetOption(attributeOption interface{})
	GetOption() interface{}

	SetDescription(attributeDescription string)
	GetDescription() string

	SetCreator(attributeCreator string)
	GetCreator() string

	ToMapStr() (frtypes.MapStr, error)
}

var _ Attribute = (*attribute)(nil)

// attribute the metadata structure definition of the model attribute
type attribute struct {
	attr      metadata.Attribute
	isNew     bool
	params    types.LogicParams
	clientSet apimachinery.ClientSetInterface
}

func (cli *attribute) searchObjects(objID string) ([]metadata.Object, error) {
	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(cli.params.Header.OwnerID).Field(common.BKObjIDField).Eq(objID)

	condStr, err := cond.ToMapStr().ToJSON()
	if nil != err {
		return nil, err
	}
	rsp, err := cli.clientSet.ObjectController().Meta().SelectObjects(context.Background(), cli.params.Header.ToHeader(), condStr)

	if nil != err {
		blog.Errorf("failed to request the object controller, error info is %s", err.Error())
		return nil, cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to search the object(%s), error info is %s", objID, rsp.ErrMsg)
		return nil, cli.params.Err.Error(rsp.Code)
	}

	return rsp.Data, nil

}

func (cli *attribute) GetParentObject() (Object, error) {

	cond := condition.CreateCondition()
	cond.Field(metadata.AssociationFieldSupplierAccount).Eq(cli.params.Header.OwnerID)
	cond.Field(metadata.AssociationFieldObjectID).Eq(cli.attr.ObjectID)
	cond.Field(metadata.AssociationFieldObjectAttributeID).Eq(cli.attr.PropertyID)

	rsp, err := cli.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), cli.params.Header.ToHeader(), cond.ToMapStr())
	if nil != err {
		blog.Errorf("[model-obj] failed to request the object controller, error info is %s", err.Error())
		return nil, err
	}

	for _, asst := range rsp.Data {

		rspRst, err := cli.searchObjects(asst.ObjectID)
		if nil != err {
			blog.Errorf("[model-obj] failed to search the object(%s)'s parent, error info is %s", asst.ObjectID, err.Error())
			return nil, err
		}

		objItems := CreateObject(cli.params, cli.clientSet, rspRst)
		for _, item := range objItems { // only one object
			return item, nil
		}

	}

	return nil, io.EOF
}
func (cli *attribute) GetChildObject() (Object, error) {

	cond := condition.CreateCondition()
	cond.Field(metadata.AssociationFieldSupplierAccount).Eq(cli.params.Header.OwnerID)
	cond.Field(metadata.AssociationFieldAssociationObjectID).Eq(cli.attr.ObjectID)
	cond.Field(metadata.AssociationFieldObjectAttributeID).Eq(cli.attr.PropertyID)

	rsp, err := cli.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), cli.params.Header.ToHeader(), cond.ToMapStr())
	if nil != err {
		blog.Errorf("[model-obj] failed to request the object controller, error info is %s", err.Error())
		return nil, err
	}

	for _, asst := range rsp.Data {

		rspRst, err := cli.searchObjects(asst.ObjectID)
		if nil != err {
			blog.Errorf("[model-obj] failed to search the object(%s)'s child, error info is %s", asst.ObjectID, err.Error())
			return nil, err
		}

		objItems := CreateObject(cli.params, cli.clientSet, rspRst)
		for _, item := range objItems { // only one object
			return item, nil
		}

	}

	return nil, io.EOF
}

func (cli *attribute) SetParentObject(objID string) error {

	cond := condition.CreateCondition()
	cond.Field(metadata.AssociationFieldSupplierAccount).Eq(cli.params.Header.OwnerID)
	cond.Field(metadata.AssociationFieldObjectAttributeID).Eq(cli.attr.PropertyID)
	cond.Field(metadata.AssociationFieldObjectID).Eq(cli.attr.ObjectID)

	rsp, err := cli.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), cli.params.Header.ToHeader(), cond.ToMapStr())
	if nil != err {
		blog.Errorf("[model-attr] failed to request the object controller, error info is %s", err.Error())
		return cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[model-attr] failed to search the parent association, error info is %s", rsp.ErrMsg)
		return cli.params.Err.Error(rsp.Code)
	}

	// create
	if 0 == len(rsp.Data) {

		asst := &metadata.Association{}
		asst.OwnerID = cli.params.Header.OwnerID
		asst.ObjectAttID = cli.attr.PropertyID
		asst.AsstObjID = objID
		asst.ObjectID = cli.attr.ObjectID

		rsp, err := cli.clientSet.ObjectController().Meta().CreateObjectAssociation(context.Background(), cli.params.Header.ToHeader(), asst)

		if nil != err {
			blog.Errorf("[model-obj] failed to request the object controller, error info is %s", err.Error())
			return cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if common.CCSuccess != rsp.Code {
			blog.Errorf("[model-obj] failed to set the main line association parent, error info is %s", rsp.ErrMsg)
			return cli.params.Err.Error(rsp.Code)
		}

		return nil
	}

	// update
	for _, asst := range rsp.Data {

		asst.AsstObjID = objID

		rsp, err := cli.clientSet.ObjectController().Meta().UpdateObjectAssociation(context.Background(), asst.ID, cli.params.Header.ToHeader(), nil)
		if nil != err {
			blog.Errorf("[model-obj] failed to request object controller, error info is %s", err.Error())
			return err
		}

		if common.CCSuccess != rsp.Code {
			blog.Errorf("[model-obj] failed to update the child association, error info is %s", rsp.ErrMsg)
			return cli.params.Err.Error(rsp.Code)
		}
	}

	return nil
}
func (cli *attribute) SetChildObject(objID string) error {

	cond := condition.CreateCondition()
	cond.Field(metadata.AssociationFieldSupplierAccount).Eq(cli.params.Header.OwnerID)
	cond.Field(metadata.AssociationFieldObjectAttributeID).Eq(cli.attr.PropertyID)
	cond.Field(metadata.AssociationFieldAssociationObjectID).Eq(cli.attr.ObjectID)

	rsp, err := cli.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), cli.params.Header.ToHeader(), cond.ToMapStr())
	if nil != err {
		blog.Errorf("[model-attr] failed to request the object controller, error info is %s", err.Error())
		return cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[model-attr] failed to search the child association, error info is %s", rsp.ErrMsg)
		return cli.params.Err.Error(rsp.Code)
	}

	// create
	if 0 == len(rsp.Data) {

		asst := &metadata.Association{}
		asst.OwnerID = cli.params.Header.OwnerID
		asst.ObjectAttID = cli.attr.PropertyID
		asst.AsstObjID = cli.attr.ObjectID
		asst.ObjectID = objID

		rsp, err := cli.clientSet.ObjectController().Meta().CreateObjectAssociation(context.Background(), cli.params.Header.ToHeader(), asst)

		if nil != err {
			blog.Errorf("[model-obj] failed to request the object controller, error info is %s", err.Error())
			return cli.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if common.CCSuccess != rsp.Code {
			blog.Errorf("[model-obj] failed to set the main line association parent, error info is %s", rsp.ErrMsg)
			return cli.params.Err.Error(rsp.Code)
		}

		return nil
	}

	// update
	for _, asst := range rsp.Data {

		asst.ObjectID = objID

		rsp, err := cli.clientSet.ObjectController().Meta().UpdateObjectAssociation(context.Background(), asst.ID, cli.params.Header.ToHeader(), nil)
		if nil != err {
			blog.Errorf("[model-obj] failed to request object controller, error info is %s", err.Error())
			return err
		}

		if common.CCSuccess != rsp.Code {
			blog.Errorf("[model-obj] failed to update the child association, error info is %s", rsp.ErrMsg)
			return cli.params.Err.Error(rsp.Code)
		}
	}

	return nil
}

func (cli *attribute) MarshalJSON() ([]byte, error) {
	return json.Marshal(cli.attr)
}

func (cli *attribute) Parse(data frtypes.MapStr) (*metadata.Attribute, error) {
	return cli.attr.Parse(data)
}

func (cli *attribute) ToMapStr() (frtypes.MapStr, error) {

	rst := metadata.SetValueToMapStrByTags(&cli.attr)
	return rst, nil

}

func (cli *attribute) Create() error {

	rsp, err := cli.clientSet.ObjectController().Meta().CreateObjectAtt(context.Background(), cli.params.Header.ToHeader(), &cli.attr)

	if nil != err {
		blog.Errorf("faield to request the object controller, the error info is %s", err.Error())
		return err
	}

	if common.CCSuccess != rsp.Code {
		return err
	}

	cli.attr.ID = rsp.Data.ID

	return nil
}

func (cli *attribute) Update() error {

	rsp, err := cli.clientSet.ObjectController().Meta().UpdateObjectAttByID(context.Background(), cli.attr.ID, cli.params.Header.ToHeader(), cli.attr.ToMapStr())

	if nil != err {
		blog.Errorf("failed to request object controller, error info is %s", err.Error())
		return err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to update the object attribute(%s), error info is %s", cli.attr.PropertyID, rsp.ErrMsg)
		return cli.params.Err.Error(common.CCErrTopoObjectAttributeUpdateFailed)
	}

	return nil
}
func (cli *attribute) search() ([]metadata.Attribute, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(cli.params.Header.OwnerID).
		Field(metadata.AttributeFieldObjectID).Eq(cli.attr.ObjectID).
		Field(metadata.AttributeFieldPropertyID).Eq(cli.attr.PropertyName)

	rsp, err := cli.clientSet.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), cli.params.Header.ToHeader(), cond.ToMapStr())

	if nil != err {
		blog.Errorf("failed to request to object controller, error info is %s", err.Error())
		return nil, err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to query the object controller, error info is %s", err.Error())
		return nil, cli.params.Err.Error(common.CCErrTopoObjectAttributeSelectFailed)
	}

	return rsp.Data, nil
}
func (cli *attribute) IsExists() (bool, error) {

	items, err := cli.search()
	if nil != err {
		return false, err
	}

	return 0 != len(items), nil
}

func (cli *attribute) Delete() error {

	cond := condition.CreateCondition()
	cond.Field(metadata.AttributeFieldObjectID).Eq(cli.attr.ObjectID).
		Field(metadata.AttributeFieldSupplierAccount).Eq(cli.params.Header.OwnerID).
		Field(metadata.AttributeFieldPropertyID).Eq(cli.attr.PropertyID)

	rsp, err := cli.clientSet.ObjectController().Meta().DeleteObjectAttByID(context.Background(), cli.attr.ID, cli.params.Header.ToHeader(), cond.ToMapStr())

	if nil != err {
		blog.Errorf("failed to request object, error info is %s", err.Error())
		return err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to delete attribute,error info is is %s", rsp.ErrMsg)
		return cli.params.Err.Error(common.CCErrTopoObjectAttributeDeleteFailed)
	}

	return nil
}

func (cli *attribute) Save() error {

	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if !exists {
		return cli.Create()
	}

	return cli.Update()
}

func (cli *attribute) SetSupplierAccount(supplierAccount string) {

	cli.attr.OwnerID = supplierAccount
}

func (cli *attribute) GetSupplierAccount() string {
	return cli.attr.OwnerID
}

func (cli *attribute) SetObjectID(objectID string) {
	cli.attr.ObjectID = objectID
}

func (cli *attribute) GetObjectID() string {
	return cli.attr.ObjectID
}

func (cli *attribute) SetID(attributeID string) {
	cli.attr.PropertyID = attributeID
}

func (cli *attribute) GetID() string {
	return cli.attr.PropertyID
}

func (cli *attribute) SetName(attributeName string) {
	cli.attr.PropertyName = attributeName
}

func (cli *attribute) GetName() string {
	return cli.attr.PropertyName
}

func (cli *attribute) SetGroup(grp Group) {
	cli.attr.PropertyGroup = grp.GetID()
}

func (cli *attribute) GetGroup() (Group, error) {
	return nil, nil
}

func (cli *attribute) SetGroupIndex(attGroupIndex int64) {
	cli.attr.PropertyIndex = attGroupIndex
}

func (cli *attribute) GetGroupIndex() int64 {
	return cli.attr.PropertyIndex
}

func (cli *attribute) SetUnint(unit string) {
	cli.attr.Unit = unit
}

func (cli *attribute) GetUnint() string {
	return cli.attr.Unit
}

func (cli *attribute) SetPlaceholder(placeHolder string) {
	cli.attr.Placeholder = placeHolder
}

func (cli *attribute) GetPlaceholder() string {
	return cli.attr.Placeholder
}

func (cli *attribute) SetIsEditable(isEditable bool) {
	cli.attr.IsEditable = isEditable
}

func (cli *attribute) GetIsEditable() bool {
	return cli.attr.IsEditable
}

func (cli *attribute) SetIsPre(isPre bool) {
	cli.attr.IsPre = isPre
}

func (cli *attribute) GetIsPre() bool {
	return cli.attr.IsPre
}

func (cli *attribute) SetIsReadOnly(isReadOnly bool) {
	cli.attr.IsReadOnly = isReadOnly
}

func (cli *attribute) GetIsReadOnly() bool {
	return cli.attr.IsReadOnly
}

func (cli *attribute) SetIsOnly(isOnly bool) {
	cli.attr.IsOnly = isOnly
}

func (cli *attribute) GetIsOnly() bool {
	return cli.attr.IsOnly
}

func (cli *attribute) SetIsSystem(isSystem bool) {
	cli.attr.IsSystem = isSystem
}

func (cli *attribute) GetIsSystem() bool {
	return cli.attr.IsSystem
}

func (cli *attribute) SetIsAPI(isAPI bool) {
	cli.attr.IsAPI = isAPI
}

func (cli *attribute) GetIsAPI() bool {
	return cli.attr.IsAPI
}

func (cli *attribute) SetType(attributeType string) {
	cli.attr.PropertyType = attributeType
}

func (cli *attribute) GetType() string {
	return cli.attr.PropertyType
}

func (cli *attribute) SetOption(attributeOption interface{}) {
	cli.attr.Option = attributeOption
}

func (cli *attribute) GetOption() interface{} {
	return cli.attr.Option
}

func (cli *attribute) SetDescription(attributeDescription string) {
	cli.attr.Description = attributeDescription
}

func (cli *attribute) GetDescription() string {
	return cli.attr.Description
}

func (cli *attribute) SetCreator(attributeCreator string) {
	cli.attr.Creator = attributeCreator
}

func (cli *attribute) GetCreator() string {
	return cli.attr.Creator
}
