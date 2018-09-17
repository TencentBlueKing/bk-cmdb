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

package operation

import (
	"context"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// AttributeOperationInterface attribute operation methods
type AttributeOperationInterface interface {
	CreateObjectAttribute(params types.ContextParams, data frtypes.MapStr) (model.Attribute, error)
	DeleteObjectAttribute(params types.ContextParams, id int64, cond condition.Condition) error
	FindObjectAttributeWithDetail(params types.ContextParams, cond condition.Condition) ([]*metadata.ObjAttDes, error)
	FindObjectAttribute(params types.ContextParams, cond condition.Condition) ([]model.Attribute, error)
	UpdateObjectAttribute(params types.ContextParams, data frtypes.MapStr, attID int64, cond condition.Condition) error

	SetProxy(modelFactory model.Factory, instFactory inst.Factory, obj ObjectOperationInterface, asst AssociationOperationInterface, grp GroupOperationInterface)
}

// NewAttributeOperation create a new attribute operation instance
func NewAttributeOperation(client apimachinery.ClientSetInterface) AttributeOperationInterface {
	return &attribute{
		clientSet: client,
	}
}

type attribute struct {
	clientSet    apimachinery.ClientSetInterface
	modelFactory model.Factory
	instFactory  inst.Factory
	obj          ObjectOperationInterface
	asst         AssociationOperationInterface
	grp          GroupOperationInterface
}

func (a *attribute) SetProxy(modelFactory model.Factory, instFactory inst.Factory, obj ObjectOperationInterface, asst AssociationOperationInterface, grp GroupOperationInterface) {
	a.modelFactory = modelFactory
	a.instFactory = instFactory
	a.obj = obj
	a.asst = asst
	a.grp = grp
}

func (a *attribute) CreateObjectAttribute(params types.ContextParams, data frtypes.MapStr) (model.Attribute, error) {

	att := a.modelFactory.CreateAttribute(params)

	_, err := att.Parse(data)
	if nil != err {
		blog.Errorf("[operation-attr] failed to parse the attribute data (%#v), error info is %s", data, err.Error())
		return nil, err
	}

	if att.GetID() == common.BKChildStr || att.GetID() == common.BKInstParentStr {
		return nil, params.Err.New(common.CCErrTopoObjectAttributeCreateFailed, "could not create bk_childid or bk_parent_id")
	}

	// check the object id
	err = a.obj.IsValidObject(params, att.GetObjectID())
	if nil != err {
		return nil, params.Err.New(common.CCErrTopoObjectAttributeCreateFailed, err.Error())
	}

	// create a new one
	err = att.Create()
	if nil != err {
		blog.Errorf("[operation-attr] failed to save the attribute data (%#v), error info is %s", data, err.Error())
		return nil, err
	}

	// create association
	attrMeta := &metadata.Association{}
	if err = data.MarshalJSONInto(attrMeta); nil != err {
		blog.Errorf("[operation-attr] failed to parse the association data, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrTopoObjectAttributeCreateFailed, err.Error())
	}

	if 0 != len(attrMeta.AsstObjID) {

		// check the object id
		err = a.obj.IsValidObject(params, attrMeta.AsstObjID)
		if nil != err {
			return nil, params.Err.New(common.CCErrTopoObjectAttributeCreateFailed, err.Error())
		}

		attrMeta.ObjectAttID = att.GetID() // the structural difference
		if err := a.asst.CreateCommonAssociation(params, attrMeta); nil != err {
			blog.Errorf("[operation-attr] failed to create the association(%v), error info is %s", attrMeta, err.Error())
			return nil, params.Err.New(common.CCErrTopoObjectAttributeCreateFailed, err.Error())
		}
	}

	return att, nil
}

func (a *attribute) DeleteObjectAttribute(params types.ContextParams, id int64, cond condition.Condition) error {

	attrCond := condition.CreateCondition()
	if id < 0 {
		attrCond = cond
	} else {
		attrCond.Field(metadata.AttributeFieldSupplierAccount).Eq(params.SupplierAccount)
		attrCond.Field(metadata.AttributeFieldID).Eq(id)
	}

	attrItems, err := a.FindObjectAttribute(params, attrCond)
	if nil != err {
		blog.Errorf("[operation-attr] failed to find the attributes by the id(%d), error info is %s", id, err.Error())
		return params.Err.New(common.CCErrTopoObjectAttributeDeleteFailed, err.Error())
	}

	for _, attrItem := range attrItems {
		// delete the association
		//fmt.Println("attr:", attrItem)
		asstCond := condition.CreateCondition()
		asstCond.Field(metadata.AssociationFieldObjectID).Eq(attrItem.GetObjectID())
		asstCond.Field(metadata.AssociationFieldSupplierAccount).Eq(attrItem.GetSupplierAccount())
		asstCond.Field(metadata.AssociationFieldObjectAttributeID).Eq(attrItem.GetID())
		if err = a.asst.DeleteAssociation(params, asstCond); nil != err {
			blog.Errorf("[operation-attr] failed to delete the attribute association(%v), error info is %s", asstCond.ToMapStr(), err.Error())
			return params.Err.New(common.CCErrTopoObjectAttributeDeleteFailed, err.Error())
		}

		// delete the attribute
		rsp, err := a.clientSet.ObjectController().Meta().DeleteObjectAttByID(context.Background(), attrItem.Origin().ID, params.Header, cond.ToMapStr())

		if nil != err {
			blog.Errorf("[operation-attr] failed to request object controller, error info is %s", err.Error())
			return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if common.CCSuccess != rsp.Code {
			blog.Errorf("[operation-attr] failed to delete the attribute by the id(%d) or the condition(%#v), error info is %s", id, cond.ToMapStr(), rsp.ErrMsg)
			return params.Err.Error(rsp.Code)
		}
	}

	return nil
}
func (a *attribute) FindObjectAttributeWithDetail(params types.ContextParams, cond condition.Condition) ([]*metadata.ObjAttDes, error) {
	attrs, err := a.FindObjectAttribute(params, cond)
	if nil != err {
		return nil, err
	}
	results := []*metadata.ObjAttDes{}
	for _, attr := range attrs {
		result := &metadata.ObjAttDes{Attribute: attr.Origin()}

		grpCond := condition.CreateCondition()
		grpCond.Field(metadata.GroupFieldGroupID).Eq(attr.Origin().PropertyGroup)
		grpCond.Field(metadata.GroupFieldSupplierAccount).Eq(attr.GetSupplierAccount())
		grpCond.Field(metadata.GroupFieldObjectID).Eq(attr.GetObjectID())
		grps, err := a.grp.FindObjectGroup(params, grpCond)
		if nil != err {
			return nil, err
		}

		for _, grp := range grps { // should be only one
			result.PropertyGroupName = grp.GetName()
		}

		assts, err := a.asst.SearchObjectAssociation(params, attr.GetObjectID())
		if nil != err {
			return nil, err
		}

		for _, asst := range assts {
			if asst.ObjectAttID == attr.GetID() { // should be only one
				result.AssociationID = asst.AsstObjID
				result.AsstForward = asst.AsstForward
			}
		}

		results = append(results, result)
	}

	return results, nil
}
func (a *attribute) FindObjectAttribute(params types.ContextParams, cond condition.Condition) ([]model.Attribute, error) {

	rsp, err := a.clientSet.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), params.Header, cond.ToMapStr())

	if nil != err {
		blog.Errorf("[operation-attr] failed to request object controller, error info is %s", err.Error())
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-attr] failed to search attribute by the condition(%#v), error info is %s", cond.ToMapStr(), rsp.ErrMsg)
		return nil, params.Err.Error(rsp.Code)
	}

	return model.CreateAttribute(params, a.clientSet, rsp.Data), nil
}

func (a *attribute) UpdateObjectAttribute(params types.ContextParams, data frtypes.MapStr, attID int64, cond condition.Condition) error {

	rsp, err := a.clientSet.ObjectController().Meta().UpdateObjectAttByID(context.Background(), attID, params.Header, data)

	if nil != err {
		blog.Errorf("[operation-attr] failed to request object controller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-attr] failed to update the attribute by the condition(%#v) or the attr-id(%d), error info is %s", cond.ToMapStr(), attID, rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}

	return nil
}
