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
	"fmt"
	"strconv"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"

	"github.com/rs/xid"
)

// ObjectOperationInterface object operation methods
type ObjectOperationInterface interface {
	CreateObjectBatch(params types.ContextParams, data frtypes.MapStr) (frtypes.MapStr, error)
	FindObjectBatch(params types.ContextParams, data frtypes.MapStr) (frtypes.MapStr, error)
	CreateObject(params types.ContextParams, isMainline bool, data frtypes.MapStr) (model.Object, error)
	CanDelete(params types.ContextParams, targetObj model.Object) error
	DeleteObject(params types.ContextParams, id int64, cond condition.Condition, needCheckInst bool) error
	FindObject(params types.ContextParams, cond condition.Condition) ([]model.Object, error)
	FindObjectTopo(params types.ContextParams, cond condition.Condition) ([]metadata.ObjectTopo, error)
	FindSingleObject(params types.ContextParams, objectID string) (model.Object, error)
	UpdateObject(params types.ContextParams, data frtypes.MapStr, id int64, cond condition.Condition) error

	SetProxy(modelFactory model.Factory, instFactory inst.Factory, cls ClassificationOperationInterface, asst AssociationOperationInterface, inst InstOperationInterface, attr AttributeOperationInterface, grp GroupOperationInterface, unique UniqueOperationInterface)
	IsValidObject(params types.ContextParams, objID string) error
}

// NewObjectOperation create a new object operation instance
func NewObjectOperation(client apimachinery.ClientSetInterface) ObjectOperationInterface {
	return &object{
		clientSet: client,
	}
}

type object struct {
	clientSet    apimachinery.ClientSetInterface
	modelFactory model.Factory
	instFactory  inst.Factory
	cls          ClassificationOperationInterface
	grp          GroupOperationInterface
	unique       UniqueOperationInterface
	asst         AssociationOperationInterface
	inst         InstOperationInterface
	attr         AttributeOperationInterface
}

func (o *object) SetProxy(modelFactory model.Factory, instFactory inst.Factory, cls ClassificationOperationInterface, asst AssociationOperationInterface, inst InstOperationInterface, attr AttributeOperationInterface, grp GroupOperationInterface, unique UniqueOperationInterface) {
	o.modelFactory = modelFactory
	o.instFactory = instFactory
	o.asst = asst
	o.inst = inst
	o.attr = attr
	o.grp = grp
	o.unique = unique
}

func (o *object) IsValidObject(params types.ContextParams, objID string) error {

	checkObjCond := condition.CreateCondition()
	checkObjCond.Field(metadata.AttributeFieldObjectID).Eq(objID)
	checkObjCond.Field(metadata.AttributeFieldSupplierAccount).Eq(params.SupplierAccount)

	objItems, err := o.FindObject(params, checkObjCond)
	if nil != err {
		blog.Errorf("[opeartion-attr] failed to check the object repeated, err: %s", err.Error())
		return params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	if 0 == len(objItems) {
		return params.Err.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("the object id  '%s' is invalid", objID))
	}

	return nil
}

func (o *object) CreateObjectBatch(params types.ContextParams, data frtypes.MapStr) (frtypes.MapStr, error) {

	inputData := map[string]ImportObjectData{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	result := frtypes.New()
	for objID, inputData := range inputData {
		subResult := frtypes.New()
		if err := o.IsValidObject(params, objID); nil != err {
			blog.Errorf("not found the  objid: %s", objID)
			subResult["errors"] = fmt.Sprintf("the object(%s) is invalid", objID)
			result[objID] = subResult
			continue
		}

		// update the object's attribute
		for idx, attr := range inputData.Attr {

			metaAttr := metadata.Attribute{}
			targetAttr, err := metaAttr.Parse(attr)
			targetAttr.OwnerID = params.SupplierAccount
			targetAttr.ObjectID = objID
			if nil != err {
				blog.Errorf("not found the  objid: %s", objID)
				subResult["errors"] = err.Error()
				result[objID] = subResult
				continue
			}

			if targetAttr.PropertyID == common.BKChildStr || targetAttr.PropertyID == common.BKInstParentStr {
				continue
			}

			if 0 == len(targetAttr.PropertyGroupName) {
				targetAttr.PropertyGroup = "Default"
			}

			// find group
			grpCond := condition.CreateCondition()
			grpCond.Field(metadata.GroupFieldObjectID).Eq(objID)
			grpCond.Field(metadata.GroupFieldGroupName).Eq(targetAttr.PropertyGroupName)
			grps, err := o.grp.FindObjectGroup(params, grpCond)
			if nil != err {
				blog.Errorf("not found the  objid: %s", objID)
				errStr := params.Lang.Languagef("import_row_int_error_str", idx, err)
				subResult["errors"] = errStr
				result[objID] = subResult
				continue
			}

			if 0 != len(grps) {
				targetAttr.PropertyGroup = grps[0].GetID() // should be only one group
			} else {

				newGrp := o.modelFactory.CreateGroup(params)
				newGrp.SetName(targetAttr.PropertyGroupName)
				newGrp.SetID(xid.New().String())
				newGrp.SetSupplierAccount(params.SupplierAccount)
				newGrp.SetObjectID(objID)
				err := newGrp.Save(nil)
				if nil != err {
					errStr := params.Lang.Languagef("import_row_int_error_str", idx, params.Err.Error(common.CCErrTopoObjectGroupCreateFailed))
					if failed, ok := subResult["insert_failed"]; ok {
						failedArr := failed.([]string)
						failedArr = append(failedArr, errStr)
						subResult["insert_failed"] = failedArr
					} else {
						subResult["insert_failed"] = []string{
							errStr,
						}
					}
					result[objID] = subResult
					continue
				}

				targetAttr.PropertyGroup = newGrp.GetID()
			}

			// create or update the attribute
			attrID, err := attr.String(metadata.AttributeFieldPropertyID)
			if nil != err {
				errStr := params.Lang.Languagef("import_row_int_error_str", idx, err.Error())
				if failed, ok := subResult["insert_failed"]; ok {
					failedArr := failed.([]string)
					failedArr = append(failedArr, errStr)
					subResult["insert_failed"] = failedArr
				} else {
					subResult["insert_failed"] = []string{
						errStr,
					}
				}
				result[objID] = subResult
				continue
			}
			attrCond := condition.CreateCondition()
			attrCond.Field(metadata.AttributeFieldSupplierAccount).Eq(params.SupplierAccount)
			attrCond.Field(metadata.AttributeFieldObjectID).Eq(objID)
			attrCond.Field(metadata.AttributeFieldPropertyID).Eq(attrID)
			attrs, err := o.attr.FindObjectAttribute(params, attrCond)
			if nil != err {
				errStr := params.Lang.Languagef("import_row_int_error_str", idx, err.Error())
				if failed, ok := subResult["insert_failed"]; ok {
					failedArr := failed.([]string)
					failedArr = append(failedArr, errStr)
					subResult["insert_failed"] = failedArr
				} else {
					subResult["insert_failed"] = []string{
						errStr,
					}
				}
				result[objID] = subResult
				continue
			}

			if 0 == len(attrs) {

				newAttr := o.modelFactory.CreateAttribute(params)
				if err = newAttr.Save(targetAttr.ToMapStr()); nil != err {
					errStr := params.Lang.Languagef("import_row_int_error_str", idx, err.Error())
					if failed, ok := subResult["insert_failed"]; ok {
						failedArr := failed.([]string)
						failedArr = append(failedArr, errStr)
						subResult["insert_failed"] = failedArr
					} else {
						subResult["insert_failed"] = []string{
							errStr,
						}
					}
					result[objID] = subResult
					continue
				}

			}

			for _, newAttr := range attrs {
				if err := newAttr.Update(targetAttr.ToMapStr()); nil != err {
					errStr := params.Lang.Languagef("import_row_int_error_str", idx, err.Error())
					if failed, ok := subResult["update_failed"]; ok {
						failedArr := failed.([]string)
						failedArr = append(failedArr, errStr)
						subResult["update_failed"] = failedArr
					} else {
						subResult["update_failed"] = []string{
							errStr,
						}
					}
					result[objID] = subResult
					continue
				}

			}

			if failed, ok := subResult["success"]; ok {
				failedArr := failed.([]string)
				failedArr = append(failedArr, strconv.FormatInt(idx, 10))
				subResult["success"] = failedArr
			} else {
				subResult["success"] = []string{
					strconv.FormatInt(idx, 10),
				}
			}
			result[objID] = subResult
		}

	}

	return result, nil
}
func (o *object) FindObjectBatch(params types.ContextParams, data frtypes.MapStr) (frtypes.MapStr, error) {

	cond := &ExportObjectCondition{}
	if err := data.MarshalJSONInto(cond); nil != err {
		return nil, err
	}

	result := frtypes.New()

	for _, objID := range cond.ObjIDS {
		obj, err := o.FindSingleObject(params, objID)
		if nil != err {
			return nil, err
		}

		attrs, err := obj.GetAttributesExceptInnerFields()
		if nil != err {
			return nil, err
		}

		result.Set(objID, frtypes.MapStr{
			"attr": attrs,
		})
	}

	return result, nil
}

func (o *object) FindSingleObject(params types.ContextParams, objectID string) (model.Object, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(objectID)

	objs, err := o.FindObject(params, cond)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the supplier account(%s) objects(%s), err: %s", params.SupplierAccount, objectID, err.Error())
		return nil, err
	}
	for _, item := range objs {
		return item, nil
	}
	return nil, params.Err.New(common.CCErrTopoObjectSelectFailed, params.Err.Errorf(common.CCErrCommParamsIsInvalid, objectID).Error())
}
func (o *object) CreateObject(params types.ContextParams, isMainline bool, data frtypes.MapStr) (model.Object, error) {
	obj := o.modelFactory.CreateObject(params)

	_, err := obj.Parse(data)
	if nil != err {
		blog.Errorf("[operation-obj] failed to parse the data(%#v), err: %s", data, err.Error())
		return nil, err
	}

	// check the classification
	_, err = obj.GetClassification()
	if nil != err {
		blog.Errorf("[operation-obj] failed to create the object, err: %s", err.Error())
		return nil, params.Err.New(common.CCErrTopoObjectCreateFailed, err.Error())
	}

	// check repeated
	exists, err := obj.IsExists()
	if nil != err {
		blog.Errorf("[operation-obj] failed to create the object(%#v), err: %s", data, err.Error())
		return nil, params.Err.New(common.CCErrTopoObjectCreateFailed, err.Error())
	}

	if exists {
		blog.Errorf("[operation-obj] the object(%#v) is repeated", data)
		return nil, params.Err.Error(common.CCErrCommDuplicateItem)
	}

	err = obj.Create()
	if nil != err {
		blog.Errorf("[operation-obj] failed to save the data(%#v), err: %s", data, err.Error())
		return nil, err
	}

	// create the default group
	grp := obj.CreateGroup()
	grp.SetDefault(true)
	grp.SetIndex(-1)
	grp.SetName("Default")
	grp.SetID("default")
	if err = grp.Save(nil); nil != err {
		blog.Errorf("[operation-obj] failed to create the default group, err: %s", err.Error())
		return nil, params.Err.Error(common.CCErrTopoObjectGroupCreateFailed)
	}

	keys := make([]metadata.UinqueKey, 0)
	// create the default inst name
	attr := obj.CreateAttribute()
	attr.SetIsOnly(true)
	attr.SetIsPre(true)
	attr.SetCreator("user")
	attr.SetIsEditable(true)
	attr.SetGroupIndex(-1)
	attr.SetGroup(grp)
	attr.SetIsRequired(true)
	attr.SetType(common.FieldTypeSingleChar)
	attr.SetID(obj.GetInstNameFieldName())
	attr.SetName(obj.GetDefaultInstPropertyName())
	if err = attr.Create(); nil != err {
		blog.Errorf("[operation-obj] failed to create the default inst name field, err: %s", err.Error())
		return nil, params.Err.Error(common.CCErrTopoObjectAttributeCreateFailed)
	}
	keys = append(keys, metadata.UinqueKey{Kind: metadata.UinqueKeyKindProperty, ID: uint64(attr.GetRecordID())})

	if isMainline {
		pAttr := obj.CreateAttribute()
		pAttr.SetIsOnly(true)
		pAttr.SetIsPre(true)
		pAttr.SetCreator("user")
		pAttr.SetIsEditable(true)
		pAttr.SetGroupIndex(-1)
		pAttr.SetGroup(grp)
		pAttr.SetIsRequired(true)
		pAttr.SetIsSystem(true)
		pAttr.SetType(common.FieldTypeInt)
		pAttr.SetID(common.BKInstParentStr)
		pAttr.SetName(obj.GetDefaultInstPropertyName())
		if err = pAttr.Create(); nil != err {
			blog.Errorf("[operation-obj] failed to create the default inst name field, err: %s", err.Error())
			return nil, params.Err.Error(common.CCErrTopoObjectAttributeCreateFailed)
		}
		keys = append(keys, metadata.UinqueKey{Kind: metadata.UinqueKeyKindProperty, ID: uint64(pAttr.GetRecordID())})
	}

	uni := obj.CreateUnique()
	uni.SetKeys(keys)
	uni.SetIsPre(false)
	uni.SetMustCheck(true)
	if err = uni.Save(nil); nil != err {
		blog.Errorf("[operation-obj] failed to create the default inst name field, err: %s", err.Error())
		return nil, err
	}

	return obj, nil
}

func (o *object) CanDelete(params types.ContextParams, targetObj model.Object) error {
	if common.IsInnerModel(targetObj.GetID()) {
		return params.Err.Error(common.CCErrTopoForbiddenToDeleteModelFailed)
	}
	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	if targetObj.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(targetObj.GetID())
	}

	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()
	findInstResponse, err := o.inst.FindOriginInst(params, targetObj, query)
	if nil != err {
		blog.Errorf("[operation-obj] failed to check if it (%s) has some insts, err: %s", targetObj.GetID(), err.Error())
		return err
	}
	if 0 != findInstResponse.Count {
		blog.Errorf("the object [%s] has been instantiated and cannot be deleted", targetObj.GetID())
		return params.Err.Errorf(common.CCErrTopoObjectHasSomeInstsForbiddenToDelete, targetObj.GetID())
	}

	or := make([]interface{}, 0)
	or = append(or, frtypes.MapStr{common.BKObjIDField: targetObj.GetID()})
	or = append(or, frtypes.MapStr{common.AssociatedObjectIDField: targetObj.GetID()})

	cond = condition.CreateCondition()
	cond.NewOR().Array(or)
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)

	assoResult, err := o.asst.SearchObject(params, &metadata.SearchAssociationObjectRequest{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("check object[%s] can be deleted, but get object associate info failed, err: %v", targetObj.GetID(), err)
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !assoResult.Result {
		blog.Errorf("check if object[%s] can be deleted, but get object associate info failed, err: %v", targetObj.GetID(), err)
		return params.Err.Error(assoResult.Code)
	}

	if len(assoResult.Data) != 0 {
		blog.Errorf("check if object[%s] can be deleted, but object has already associate to another one.", targetObj.GetID())
		return params.Err.Error(common.CCErrorTopoObjectHasAlreadyAssociated)
	}

	return nil
}
func (o *object) DeleteObject(params types.ContextParams, id int64, cond condition.Condition, needCheckInst bool) error {

	if 0 < id {
		cond = condition.CreateCondition()
		cond.Field(metadata.ModelFieldID).Eq(id)
	}

	objs, err := o.FindObject(params, cond)
	if nil != err {
		blog.Errorf("[operation-obj] failed to find objects, the condition is (%v) err: %s", cond, err.Error())
		return err
	}

	for _, obj := range objs {

		// check if is can be deleted
		if needCheckInst {
			if err := o.CanDelete(params, obj); nil != err {
				return err
			}
		}

		// delete object
		if unis, err := obj.GetUniques(); err != nil {
			blog.Errorf("[operation-asst] failed to get the object's uniques, err: %s", err.Error())
			return err
		} else {
			for _, uni := range unis {
				if err = o.unique.Delete(params, obj.GetID(), uni.GetRecordID()); err != nil {
					blog.Errorf("[operation-asst] failed to delete the object's uniques, err: %s", err.Error())
					return err
				}
			}
		}

		attrCond := condition.CreateCondition()
		attrCond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
		attrCond.Field(common.BKObjIDField).Eq(obj.GetID())

		if err := o.attr.DeleteObjectAttribute(params, attrCond); nil != err {
			blog.Errorf("[operation-obj] failed to delete the object(%d)'s attribute, err: %s", id, err.Error())
			return err
		}

		if groups, err := obj.GetGroups(); err != nil {
			blog.Errorf("[operation-asst] failed to get the object's groups, err: %s", err.Error())
			return err
		} else {
			for _, group := range groups {
				if err = o.grp.DeleteObjectGroup(params, group.GetRecordID()); err != nil {
					blog.Errorf("[operation-asst] failed to delete the object's groups, err: %s", err.Error())
					return err
				}
			}
		}

		rsp, err := o.clientSet.ObjectController().Meta().DeleteObject(context.Background(), obj.GetRecordID(), params.Header, cond.ToMapStr())

		if nil != err {
			blog.Errorf("[operation-obj] failed to request the object controller, err: %s", err.Error())
			return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if common.CCSuccess != rsp.Code {
			blog.Errorf("[opration-obj] failed to delete the object by the condition(%#v) or the id(%d)", cond.ToMapStr(), id)
			return params.Err.Error(rsp.Code)
		}

	}
	return nil
}

func (o *object) isFrom(params types.ContextParams, fromObjID, toObjID string) (bool, error) {

	asstItems, err := o.asst.SearchObjectAssociation(params, fromObjID)
	if nil != err {
		return false, err
	}

	for _, asst := range asstItems {
		if asst.AsstObjID == toObjID {
			return true, nil
		}
	}

	return false, nil
}

func (o *object) FindObjectTopo(params types.ContextParams, cond condition.Condition) ([]metadata.ObjectTopo, error) {
	objs, err := o.FindObject(params, cond)
	if nil != err {
		blog.Errorf("[operation-obj] failed to find object, err: %s", err.Error())
		return nil, err
	}

	results := make([]metadata.ObjectTopo, 0)
	for _, obj := range objs {
		asstItems, err := o.asst.SearchObjectAssociation(params, obj.GetID())
		if nil != err {
			return nil, err
		}

		for _, asst := range asstItems {

			// find association kind with association kind id.
			typeCond := condition.CreateCondition()
			typeCond.Field(common.AssociationKindIDField).Eq(asst.AsstKindID)
			typeCond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
			request := &metadata.SearchAssociationTypeRequest{
				Condition: typeCond.ToMapStr(),
			}

			resp, err := o.asst.SearchType(params, request)
			if err != nil {
				blog.Errorf("find object topo failed, because get association kind[%s] failed, err: %v", asst.AsstKindID, err)
				return nil, params.Err.Errorf(common.CCErrTopoGetAssociationKindFailed, asst.AsstKindID)
			}
			if !resp.Result {
				blog.Errorf("find object topo failed, because get association kind[%s] failed, err: %v", asst.AsstKindID, resp.ErrMsg)
				return nil, params.Err.Errorf(common.CCErrTopoGetAssociationKindFailed, asst.AsstKindID)
			}

			// should only be one association kind.
			if len(resp.Data.Info) == 0 {
				blog.Errorf("find object topo failed, because get association kind[%s] failed, err: can not find this association kind.", asst.AsstKindID)
				return nil, params.Err.Errorf(common.CCErrTopoGetAssociationKindFailed, asst.AsstKindID)
			}

			cond = condition.CreateCondition()
			cond.Field(common.BKObjIDField).Eq(asst.AsstObjID)
			cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)

			asstObjs, err := o.FindObject(params, cond)
			if nil != err {
				blog.Errorf("[operation-obj] failed to find object, err: %s", err.Error())
				return nil, err
			}

			for _, asstObj := range asstObjs {
				tmp := metadata.ObjectTopo{}
				tmp.Label = resp.Data.Info[0].AssociationKindName
				tmp.LabelName = resp.Data.Info[0].AssociationKindName
				tmp.From.ObjID = obj.GetID()
				cls, err := obj.GetClassification()
				if nil != err {
					return nil, err
				}
				tmp.From.ClassificationID = cls.GetID()
				tmp.From.Position = obj.GetPosition()
				tmp.From.OwnerID = obj.GetSupplierAccount()
				tmp.From.ObjName = obj.GetName()
				tmp.To.OwnerID = asstObj.GetSupplierAccount()
				tmp.To.ObjID = asstObj.GetID()

				cls, err = asstObj.GetClassification()
				if nil != err {
					return nil, err
				}
				tmp.To.ClassificationID = cls.GetID()
				tmp.To.Position = asstObj.GetPosition()
				tmp.To.ObjName = asstObj.GetName()
				ok, err := o.isFrom(params, asstObj.GetID(), obj.GetID())
				if nil != err {
					return nil, err
				}

				if ok {
					tmp.Arrows = "to,from"
				} else {
					tmp.Arrows = "to"
				}

				results = append(results, tmp)
			}
		}

	}

	return results, nil
}

func (o *object) FindObject(params types.ContextParams, cond condition.Condition) ([]model.Object, error) {

	rsp, err := o.clientSet.ObjectController().Meta().SelectObjects(context.Background(), params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-obj] failed to request the object controller, err: %s", err.Error())
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-obj] failed to search the objects by the condition(%#v) , err: %s", cond.ToMapStr(), rsp.ErrMsg)
		return nil, params.Err.Error(rsp.Code)
	}

	return model.CreateObject(params, o.clientSet, rsp.Data), nil
}

func (o *object) UpdateObject(params types.ContextParams, data frtypes.MapStr, id int64, cond condition.Condition) error {

	obj := o.modelFactory.CreateObject(params)
	obj.SetRecordID(id)
	_, err := obj.Parse(data)
	if nil != err {
		blog.Errorf("[operation-obj] failed to parse the data(%#v), err: %s", data, err.Error())
		return err
	}

	// check repeated
	exists, err := obj.IsExists()
	if nil != err {
		blog.Errorf("[operation-obj] failed to update the object(%#v), err: %s", data, err.Error())
		return params.Err.New(common.CCErrTopoObjectCreateFailed, err.Error())
	}

	if exists {
		blog.Errorf("[operation-obj] the object(%#v) is repeated", data)
		return params.Err.Error(common.CCErrCommDuplicateItem)
	}
	if err = obj.Update(data); nil != err {
		blog.Errorf("[operation-obj] failed to update the object(%d), the new data(%#v), err: %s", id, data, err.Error())
		return params.Err.New(common.CCErrTopoObjectUpdateFailed, err.Error())
	}

	return nil
}
