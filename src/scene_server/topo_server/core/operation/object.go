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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// ObjectOperationInterface object operation methods
type ObjectOperationInterface interface {
	CreateObjectBatch(params types.ContextParams, data mapstr.MapStr) (mapstr.MapStr, error)
	FindObjectBatch(params types.ContextParams, data mapstr.MapStr) (mapstr.MapStr, error)
	CreateObject(params types.ContextParams, isMainline bool, data mapstr.MapStr) (model.Object, error)
	CanDelete(params types.ContextParams, targetObj model.Object) error
	DeleteObject(params types.ContextParams, id int64, cond condition.Condition, needCheckInst bool) error
	FindObject(params types.ContextParams, cond condition.Condition) ([]model.Object, error)
	FindObjectTopo(params types.ContextParams, cond condition.Condition) ([]metadata.ObjectTopo, error)
	FindSingleObject(params types.ContextParams, objectID string) (model.Object, error)
	UpdateObject(params types.ContextParams, data mapstr.MapStr, id int64, cond condition.Condition) error

	SetProxy(modelFactory model.Factory, instFactory inst.Factory, cls ClassificationOperationInterface, asst AssociationOperationInterface, inst InstOperationInterface, attr AttributeOperationInterface, grp GroupOperationInterface, unique UniqueOperationInterface)
	IsValidObject(params types.ContextParams, objID string) error

	CreateOneObject(params types.ContextParams, data mapstr.MapStr) (model.Object, error)
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

func (o *object) CreateObjectBatch(params types.ContextParams, data mapstr.MapStr) (mapstr.MapStr, error) {

	inputData := map[string]ImportObjectData{}
	if err := data.MarshalJSONInto(&inputData); nil != err {
		return nil, err
	}

	result := mapstr.New()
	hasError := false
	for objID, inputData := range inputData {
		subResult := mapstr.New()
		if err := o.IsValidObject(params, objID); nil != err {
			blog.Errorf("not found the  objid: %s", objID)
			subResult["errors"] = fmt.Sprintf("the object(%s) is invalid", objID)
			result[objID] = subResult
			hasError = true
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
				hasError = true
				continue
			}

			if targetAttr.PropertyType == common.FieldTypeMultiAsst || targetAttr.PropertyType == common.FieldTypeSingleAsst {
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
				hasError = true
				continue
			}

			if 0 != len(grps) {
				targetAttr.PropertyGroup = grps[0].Group().GroupID // should be only one group
			} else {

				newGrp := o.modelFactory.CreateGroup(params)
				newGrp.SetGroup(metadata.Group{
					GroupName: targetAttr.PropertyGroupName,
					GroupID:   model.NewGroupID(false),
					ObjectID:  objID,
					OwnerID:   params.SupplierAccount,
				})
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
					hasError = true
					continue
				}

				targetAttr.PropertyGroup = newGrp.Group().GroupID
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
				hasError = true
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
				hasError = true
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
					hasError = true
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
					hasError = true
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

	if hasError {
		return result, params.Err.Error(common.CCErrCommNotAllSuccess)
	}
	return result, nil

}
func (o *object) FindObjectBatch(params types.ContextParams, data mapstr.MapStr) (mapstr.MapStr, error) {

	cond := &ExportObjectCondition{}
	if err := data.MarshalJSONInto(cond); nil != err {
		return nil, err
	}

	result := mapstr.New()

	for _, objID := range cond.ObjIDS {
		obj, err := o.FindSingleObject(params, objID)
		if nil != err {
			return nil, err
		}

		attrs, err := obj.GetAttributesExceptInnerFields()
		if nil != err {
			return nil, err
		}

		result.Set(objID, mapstr.MapStr{
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
func (o *object) CreateObject(params types.ContextParams, isMainline bool, data mapstr.MapStr) (model.Object, error) {
	obj := o.modelFactory.CreateObject(params)
	err := obj.Parse(data)
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
		return nil, params.Err.Errorf(common.CCErrCommDuplicateItem, "")
	}

	err = obj.Create()
	if nil != err {
		blog.Errorf("[operation-obj] failed to save the data(%#v), err: %s", data, err.Error())
		return nil, err
	}

	// create the default group
	grp := obj.CreateGroup()
	groupData := metadata.Group{
		IsDefault:  true,
		GroupIndex: -1,
		GroupName:  "Default",
		GroupID:    model.NewGroupID(true),
		ObjectID:   obj.Object().ObjectID,
		OwnerID:    obj.Object().OwnerID,
	}
	if nil != params.MetaData {
		groupData.Metadata = *params.MetaData
	}
	grp.SetGroup(groupData)

	if err = grp.Save(nil); nil != err {
		blog.Errorf("[operation-obj] failed to create the default group, err: %s", err.Error())
		return nil, params.Err.Error(common.CCErrTopoObjectGroupCreateFailed)
	}

	keys := make([]metadata.UniqueKey, 0)
	// create the default inst name
	group := grp.Group()
	attr := obj.CreateAttribute()
	attr.SetAttribute(metadata.Attribute{
		ObjectID:          obj.Object().ObjectID,
		IsOnly:            true,
		IsPre:             true,
		Creator:           "user",
		IsEditable:        true,
		PropertyIndex:     -1,
		PropertyGroup:     group.GroupID,
		PropertyGroupName: group.GroupName,
		IsRequired:        true,
		PropertyType:      common.FieldTypeSingleChar,
		PropertyID:        obj.GetInstNameFieldName(),
		PropertyName:      obj.GetDefaultInstPropertyName(),
	})
	if nil != params.MetaData {
		attr.Attribute().Metadata = *params.MetaData
	}
	if err = attr.Create(); nil != err {
		blog.Errorf("[operation-obj] failed to create the default inst name field, error info is %s", err.Error())
		return nil, err
	}

	keys = append(keys, metadata.UniqueKey{Kind: metadata.UniqueKeyKindProperty, ID: uint64(attr.Attribute().ID)})

	if isMainline {
		pAttr := obj.CreateAttribute()
		pAttr.SetAttribute(metadata.Attribute{
			ObjectID:          obj.Object().ObjectID,
			IsOnly:            true,
			IsPre:             true,
			Creator:           "user",
			IsEditable:        true,
			PropertyIndex:     -1,
			PropertyGroup:     group.GroupID,
			PropertyGroupName: group.GroupName,
			IsRequired:        true,
			PropertyType:      common.FieldTypeInt,
			PropertyID:        common.BKInstParentStr,
			PropertyName:      common.BKInstParentStr,
			IsSystem:          true,
		})

		if err = pAttr.Create(); nil != err {
			blog.Errorf("[operation-obj] failed to create the default inst name field, err: %s", err.Error())
			return nil, params.Err.Error(common.CCErrTopoObjectAttributeCreateFailed)
		}
		keys = append(keys, metadata.UniqueKey{Kind: metadata.UniqueKeyKindProperty, ID: uint64(pAttr.Attribute().ID)})
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

	tObject := targetObj.Object()
	cond := condition.CreateCondition()
	if targetObj.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(tObject.ObjectID)
	}

	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()
	findInstResponse, err := o.inst.FindOriginInst(params, targetObj, query)
	if nil != err {
		blog.Errorf("[operation-obj] failed to check if it (%s) has some insts, err: %s", tObject.ObjectID, err.Error())
		return err
	}
	if 0 != findInstResponse.Count {
		blog.Errorf("the object [%s] has been instantiated and cannot be deleted", tObject.ObjectID)
		return params.Err.Errorf(common.CCErrTopoObjectHasSomeInstsForbiddenToDelete, tObject.ObjectID)
	}

	or := make([]interface{}, 0)
	or = append(or, mapstr.MapStr{common.BKObjIDField: tObject.ObjectID})
	or = append(or, mapstr.MapStr{common.AssociatedObjectIDField: tObject.ObjectID})

	cond = condition.CreateCondition()
	cond.NewOR().Array(or)

	assoResult, err := o.asst.SearchObject(params, &metadata.SearchAssociationObjectRequest{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("check object[%s] can be deleted, but get object associate info failed, err: %v", tObject.ObjectID, err)
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !assoResult.Result {
		blog.Errorf("check if object[%s] can be deleted, but get object associate info failed, err: %v", tObject.ObjectID, err)
		return params.Err.Error(assoResult.Code)
	}

	if len(assoResult.Data) != 0 {
		blog.Errorf("check if object[%s] can be deleted, but object has already associate to another one.", tObject.ObjectID)
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
		object := obj.Object()
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
				if err = o.unique.Delete(params, object.ObjectID, uni.GetRecordID()); err != nil {
					blog.Errorf("[operation-asst] failed to delete the object's uniques, %s", err.Error())
					return err
				}
			}
		}

		attrCond := condition.CreateCondition()
		attrCond.Field(common.BKObjIDField).Eq(object.ObjectID)

		if err := o.attr.DeleteObjectAttribute(params, attrCond); nil != err {
			blog.Errorf("[operation-obj] failed to delete the object(%d)'s attribute, err: %s", id, err.Error())
			return err
		}

		if groups, err := obj.GetGroups(); err != nil {
			blog.Errorf("[operation-asst] failed to get the object's groups, err: %s", err.Error())
			return err
		} else {
			for _, group := range groups {
				if err = o.grp.DeleteObjectGroup(params, group.Group().ID); err != nil {
					blog.Errorf("[operation-asst] failed to delete the object's groups, err: %s", err.Error())
					return err
				}
			}
		}

		rsp, err := o.clientSet.CoreService().Model().DeleteModel(context.Background(), params.Header, &metadata.DeleteOption{Condition: cond.ToMapStr()})
		if nil != err {
			blog.Errorf("[operation-obj] failed to request the object controller, err: %s", err.Error())
			return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !rsp.Result {
			blog.Errorf("[opration-obj] failed to delete the object by the condition(%#v) or the id(%d)", cond.ToMapStr(), id)
			return params.Err.New(rsp.Code, rsp.ErrMsg)
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
		object := obj.Object()
		asstItems, err := o.asst.SearchObjectAssociation(params, object.ObjectID)
		if nil != err {
			return nil, err
		}

		for _, asst := range asstItems {

			// find association kind with association kind id.
			typeCond := condition.CreateCondition()
			typeCond.Field(common.AssociationKindIDField).Eq(asst.AsstKindID)
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

			asstObjs, err := o.FindObject(params, cond)
			if nil != err {
				blog.Errorf("[operation-obj] failed to find object, err: %s", err.Error())
				return nil, err
			}

			for _, asstObj := range asstObjs {
				assoObject := asstObj.Object()
				tmp := metadata.ObjectTopo{}
				tmp.Label = resp.Data.Info[0].AssociationKindName
				tmp.LabelName = resp.Data.Info[0].AssociationKindName
				tmp.From.ObjID = object.ObjectID
				cls, err := obj.GetClassification()
				if nil != err {
					return nil, err
				}
				tmp.From.ClassificationID = cls.Classify().ClassificationID
				tmp.From.Position = object.Position
				tmp.From.OwnerID = object.OwnerID
				tmp.From.ObjName = object.ObjectName
				tmp.To.OwnerID = assoObject.OwnerID
				tmp.To.ObjID = assoObject.ObjectID

				cls, err = asstObj.GetClassification()
				if nil != err {
					return nil, err
				}
				tmp.To.ClassificationID = cls.Classify().ClassificationID
				tmp.To.Position = assoObject.Position
				tmp.To.ObjName = assoObject.ObjectName
				ok, err := o.isFrom(params, assoObject.ObjectID, object.ObjectID)
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
	fCond := cond.ToMapStr()
	if nil != params.MetaData {
		fCond.Merge(metadata.PublicAndBizCondition(*params.MetaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}
	rsp, err := o.clientSet.CoreService().Model().ReadModel(context.Background(), params.Header, &metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-obj] failed to request the object controller, err: %s", err.Error())
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-obj] failed to search the objects by the condition(%#v) , error info is %s", fCond, rsp.ErrMsg)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	models := []metadata.Object{}
	for index := range rsp.Data.Info {
		models = append(models, rsp.Data.Info[index].Spec)
	}
	return model.CreateObject(params, o.clientSet, models), nil
}

func (o *object) UpdateObject(params types.ContextParams, data mapstr.MapStr, id int64, cond condition.Condition) error {

	obj := o.modelFactory.CreateObject(params)
	obj.SetRecordID(id)
	err := obj.Parse(data)
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
		return params.Err.Errorf(common.CCErrCommDuplicateItem, "")
	}
	if err = obj.Update(data); nil != err {
		blog.Errorf("[operation-obj] failed to update the object(%d), the new data(%#v), err: %s", id, data, err.Error())
		return params.Err.New(common.CCErrTopoObjectUpdateFailed, err.Error())
	}

	return nil
}

func (o *object) CreateOneObject(params types.ContextParams, data mapstr.MapStr) (model.Object, error) {
	obj := o.modelFactory.CreateObject(params)

	err := obj.Parse(data)
	if nil != err {
		blog.Errorf("[operation-obj] failed to parse the data(%#v), error info is %s", data, err.Error())
		return nil, err
	}

	// check the classification
	_, err = obj.GetClassification()
	if nil != err {
		blog.Errorf("[operation-obj] failed to create the object, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrTopoObjectCreateFailed, err.Error())
	}

	err = obj.Create()
	if nil != err {
		blog.Errorf("[operation-obj] failed to save the data(%#v), error info is %s", data, err.Error())
		return nil, err
	}

	// create the default group
	grp := obj.CreateGroup()
	grp.SetGroup(metadata.Group{
		IsDefault:  true,
		GroupIndex: -1,
		GroupName:  "Default",
		GroupID:    model.NewGroupID(true),
		ObjectID:   obj.Object().ObjectID,
		OwnerID:    obj.Object().OwnerID,
	})

	if err = grp.Save(nil); nil != err {
		blog.Errorf("[operation-obj] failed to create the default group, error info is %s", err.Error())
		return nil, err
	}

	group := grp.Group()
	// create the default inst name
	attr := obj.CreateAttribute()
	attr.SetAttribute(metadata.Attribute{
		IsOnly:            true,
		IsPre:             true,
		Creator:           "user",
		IsEditable:        true,
		PropertyIndex:     -1,
		PropertyGroup:     group.GroupID,
		PropertyGroupName: group.GroupName,
		IsRequired:        true,
		PropertyType:      common.FieldTypeSingleChar,
		PropertyID:        obj.GetInstNameFieldName(),
		PropertyName:      obj.GetDefaultInstPropertyName(),
	})
	if nil != params.MetaData {
		attr.Attribute().Metadata = *params.MetaData
	}

	if err = attr.Create(); nil != err {
		blog.Errorf("[operation-obj] failed to create the default inst name field, error info is %s", err.Error())
		return nil, err
	}

	uni := obj.CreateUnique()
	uni.SetKeys([]metadata.UniqueKey{{Kind: metadata.UniqueKeyKindProperty, ID: uint64(attr.Attribute().ID)}})
	uni.SetIsPre(false)
	uni.SetMustCheck(true)
	if err = uni.Save(nil); nil != err {
		blog.Errorf("[operation-obj] failed to create the default inst name field, error info is %s", err.Error())
		return nil, err
	}

	return obj, nil
}
