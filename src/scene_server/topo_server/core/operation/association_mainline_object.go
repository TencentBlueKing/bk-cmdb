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
	"io"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (a *association) DeleteMainlineAssociaton(params types.ContextParams, objID string) error {

	targetObj, err := a.obj.FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[opeartion-asst] failed to find the target object(%s), error info is %s", objID, err.Error())
		return err
	}

	parentObj, err := targetObj.GetMainlineParentObject()
	if nil != err {
		blog.Errorf("[operation-asst] failed to find the object(%s)'s parent, error info is %s", objID, err.Error())
		return err
	}

	// update associations
	childObj, err := targetObj.GetMainlineChildObject()
	if nil != err && io.EOF != err {
		blog.Errorf("[operation-asst] failed to find the object(%s)'s child, error info is %s", objID, err.Error())
		return err
	}

	if err = a.ResetMainlineInstAssociatoin(params, targetObj); nil != err && io.EOF != err {
		blog.Errorf("[operation-asst] failed to delete the object(%s)'s insts, error info %s", objID, err.Error())
		return err
	}

	if nil != childObj { // FIX: 正常情况下 childObj 不可以能为 nil，只有在拓扑异常的时候才会出现

		if err = childObj.SetMainlineParentObject(parentObj.GetID()); nil != err && io.EOF != err {
			blog.Errorf("[operation-asst] failed to update the association, error info is %s", err.Error())
			return err
		}

	}
	// delete objects
	if err = a.obj.DeleteObject(params, targetObj.GetRecordID(), nil, false); nil != err && io.EOF != err {
		blog.Errorf("[operation-asst] failed to delete the object(%s), error info is %s", targetObj.GetID(), err.Error())
		return err
	}

	// delete the object associations
	cond := condition.CreateCondition()
	cond.Field(metadata.AssociationFieldObjectID).Eq(targetObj.GetID())
	cond.Field(common.BKOwnerIDField).Eq(targetObj.GetSupplierAccount())
	if err = a.DeleteAssociation(params, cond); nil != err {
		blog.Errorf("[operation-asst] failed to delete the association, error info is %s", err.Error())
		return err
	}

	cond = condition.CreateCondition()
	cond.Field(metadata.AssociationFieldAssociationObjectID).Eq(targetObj.GetID())
	cond.Field(common.BKOwnerIDField).Eq(targetObj.GetSupplierAccount())
	if err = a.DeleteAssociation(params, cond); nil != err {
		blog.Errorf("[operation-asst] failed to delete the association, error info is %s", err.Error())
		return err
	}

	return nil
}

func (a *association) SearchMainlineAssociationTopo(params types.ContextParams, targetObj model.Object) ([]*metadata.MainlineObjectTopo, error) {

	results := make([]*metadata.MainlineObjectTopo, 0)

	for {

		tmpRst := &metadata.MainlineObjectTopo{}
		tmpRst.ObjID = targetObj.GetID()
		tmpRst.ObjName = targetObj.GetName()
		tmpRst.OwnerID = params.SupplierAccount

		parentObj, err := targetObj.GetMainlineParentObject()
		if nil == err {
			tmpRst.PreObjID = parentObj.GetID()
			tmpRst.PreObjName = parentObj.GetName()
		} else if nil != err && io.EOF != err {
			return nil, err
		}

		childObj, err := targetObj.GetMainlineChildObject()
		if nil == err {
			tmpRst.NextObj = childObj.GetID()
			tmpRst.NextName = childObj.GetName()
		} else if nil != err {
			if io.EOF == err {
				results = append(results, tmpRst)
				return results, nil
			}
			return nil, err
		}

		results = append(results, tmpRst)
		targetObj = childObj
	}

}

func (a *association) CreateMainlineAssociation(params types.ContextParams, data *metadata.Association) (model.Object, error) {

	// check the level
	bizObj, err := a.obj.FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("[operation-asst] failed to check the mainline topo level, error info is %s", err.Error())
		return nil, err
	}

	items, err := a.SearchMainlineAssociationTopo(params, bizObj)
	if nil != err {
		blog.Errorf("[operation-asst] failed to check the mainline topo level, error info is %s", err.Error())
		return nil, err
	}

	if len(items) >= params.MaxTopoLevel {
		blog.Errorf("[operation-asst] the mainline topo leve is %d, the max limit is %d", len(items), params.MaxTopoLevel)
		return nil, params.Err.Error(common.CCErrTopoBizTopoLevelOverLimit)
	}

	// check and fetch the association object's classification
	objCls, err := a.cls.FindSingleClassification(params, data.ClassificationID)
	if nil != err {
		blog.Errorf("[opration-asst] failed to find the single classification, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, data.ClassificationID)
	}

	// find the mainline parent object
	parentObj, err := a.obj.FindSingleObject(params, data.AsstObjID)
	switch t := err.(type) {
	case nil:
	default:
		blog.Errorf("[operation-asst] failed to find the single object(%s), error info is %s", data.ObjectID, t.Error())
		return nil, t
	case errors.CCErrorCoder:
		if t.GetCode() == common.CCErrTopoObjectSelectFailed {
			blog.Errorf("[operation-asst] failed to find the single object(%s), error info is %s", data.ObjectID, t.Error())
			return nil, t
		}
	}

	// find the mainline child object for the parent
	childObj, err := parentObj.GetMainlineChildObject()
	if nil != err {
		blog.Errorf("[operation-asst] failed to find the child object for the object(%s), error info is %s", parentObj.GetID(), err.Error())
		return nil, err
	}

	// check and create the association mainline object
	if err = a.obj.IsValidObject(params, data.ObjectID); nil == err {
		blog.Errorf("[operation-asst] the object(%s) is duplicate", data.ObjectID)
		return nil, params.Err.Errorf(common.CCErrCommDuplicateItem, data.ObjectID)
	}

	currentObj := a.modelFactory.CreaetObject(params)
	currentObj.SetID(data.ObjectID)
	currentObj.SetName(data.ObjectName)
	currentObj.SetIcon(data.ObjectIcon)
	currentObj.SetClassification(objCls)

	if err = currentObj.Save(nil); nil != err {
		blog.Errorf("[operation-asst] failed to create the object(%s), error info is %s", currentObj.GetID(), err.Error())
		return nil, err
	}

	attr := currentObj.CreateAttribute()
	attr.SetIsSystem(true)
	attr.SetID(common.BKChildStr)
	attr.SetType(common.FieldTypeLongChar)
	attr.SetName(common.BKChildStr)
	attr.SetOption(nil)

	if err = attr.Save(nil); nil != err {
		blog.Errorf("[operation-asst] failed to create the object(%s) attribute(%s), error info is %s", currentObj.GetID(), common.BKChildStr, err.Error())
		return nil, err
	}

	// create the default group
	grp := currentObj.CreateGroup()
	grp.SetDefault(true)
	grp.SetIndex(-1)
	grp.SetName("Default")
	grp.SetID("default")
	if err = grp.Save(nil); nil != err {
		blog.Errorf("[operation-obj] failed to create the default group, error info is %s", err.Error())
		return nil, err
	}

	defaultInstNameAttr := currentObj.CreateAttribute()
	defaultInstNameAttr.SetIsSystem(false)
	defaultInstNameAttr.SetIsOnly(true)
	defaultInstNameAttr.SetIsPre(true)
	defaultInstNameAttr.SetIsEditable(true)
	defaultInstNameAttr.SetType(common.FieldTypeLongChar)
	defaultInstNameAttr.SetIsRequired(true)
	defaultInstNameAttr.SetID(currentObj.GetInstNameFieldName())
	defaultInstNameAttr.SetName(currentObj.GetDefaultInstPropertyName())
	defaultInstNameAttr.SetGroupIndex(-1)
	defaultInstNameAttr.SetGroup(grp)

	if err = defaultInstNameAttr.Save(nil); nil != err {
		blog.Errorf("[operation-asst] failed to create the object(%s) attribute(%s), error info is %s", currentObj.GetID(), currentObj.GetDefaultInstPropertyName(), err.Error())
		return nil, err
	}

	defaultInstParentAttr := currentObj.CreateAttribute()
	defaultInstParentAttr.SetIsSystem(true)
	defaultInstParentAttr.SetIsOnly(true)
	defaultInstParentAttr.SetIsEditable(false)
	defaultInstParentAttr.SetType(common.FieldTypeInt)
	defaultInstParentAttr.SetIsRequired(true)
	defaultInstParentAttr.SetID(common.BKInstParentStr)
	defaultInstParentAttr.SetName(common.BKInstParentStr)

	if err = defaultInstParentAttr.Save(nil); nil != err {
		blog.Errorf("[operation-asst] failed to create the object(%s) attribute(%s), error info is %s", currentObj.GetID(), common.BKInstParentStr, err.Error())
		return nil, err
	}

	// update the mainline topo inst association
	if err = a.SetMainlineInstAssociation(params, parentObj, currentObj, childObj); nil != err {
		blog.Errorf("[operation-asst] failed set the mainline inst association, error info is %s", err.Error())
		return nil, err
	}

	// reset the parent's child object
	if err = parentObj.SetMainlineChildObject(currentObj.GetID()); nil != err {
		blog.Errorf("[operation-asst] failed to set the mainline object, error info is %s", err.Error())
		return nil, err
	}

	// reset the current's child object
	if err = currentObj.SetMainlineChildObject(childObj.GetID()); nil != err {
		blog.Errorf("[operation-asst] failed to set the mainline object, error info is %s ", err.Error())
		return nil, err
	}

	return currentObj, nil
}
