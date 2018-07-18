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

// ObjectOperationInterface object operation methods
type ObjectOperationInterface interface {
	CreateObjectBatch(params types.ContextParams, data frtypes.MapStr) error
	FindObjectBatch(params types.ContextParams, data frtypes.MapStr) error
	CreateObject(params types.ContextParams, data frtypes.MapStr) (model.Object, error)
	DeleteObject(params types.ContextParams, id int64, cond condition.Condition) error
	FindObject(params types.ContextParams, cond condition.Condition) ([]model.Object, error)
	FindObjectTopo(params types.ContextParams, cond condition.Condition) ([]metadata.ObjectTopo, error)
	FindSingleObject(params types.ContextParams, objectID string) (model.Object, error)
	UpdateObject(params types.ContextParams, data frtypes.MapStr, id int64, cond condition.Condition) error

	SetProxy(modelFactory model.Factory, instFactory inst.Factory, cls ClassificationOperationInterface, asst AssociationOperationInterface, inst InstOperationInterface, attr AttributeOperationInterface)
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
	asst         AssociationOperationInterface
	inst         InstOperationInterface
	attr         AttributeOperationInterface
}

func (o *object) SetProxy(modelFactory model.Factory, instFactory inst.Factory, cls ClassificationOperationInterface, asst AssociationOperationInterface, inst InstOperationInterface, attr AttributeOperationInterface) {
	o.modelFactory = modelFactory
	o.instFactory = instFactory
	o.asst = asst
	o.inst = inst
	o.attr = attr
}

func (o *object) IsValidObject(params types.ContextParams, objID string) error {

	checkObjCond := condition.CreateCondition()
	checkObjCond.Field(metadata.AttributeFieldObjectID).Eq(objID)
	checkObjCond.Field(metadata.AttributeFieldSupplierAccount).Eq(params.SupplierAccount)

	objItems, err := o.FindObject(params, checkObjCond)
	if nil != err {
		blog.Errorf("[opeartion-attr] failed to check the object repeated, error info is %s", err.Error())
		return params.Err.New(common.CCErrTopoObjectAttributeCreateFailed, err.Error())
	}

	if 0 == len(objItems) {
		return params.Err.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("the object id  '%s' is invalid", objID))
	}

	return nil
}

func (o *object) CreateObjectBatch(params types.ContextParams, data frtypes.MapStr) error {
	/*
		* new fuction
			result := map[string]interface{}{}

			// parse the json get the object id
			for objID := range data {

				subResult := map[string]interface{}{}

				// check the object
				cond := condition.CreateCondition()
				cond.Field(common.BKOwnerIDField).Eq(ownerID)
				cond.Field(common.BKObjIDField).Eq(objID)

				items, err := o.FindObject(params, cond)
				if nil != err {
					blog.Error("failed to search, the error info is :%s", err.Error())
					subResult["errors"] = fmt.Sprintf("the object(%s) is invalid", objID)
					result[objID] = subResult
					continue
				}
				if 0 == len(items) {
					// TODO: may be need to create the object in the future version
					blog.Error("not found the  objid: %s", objID)
					subResult["errors"] = fmt.Sprintf("the object(%s) is invalid", objID)
					result[objID] = subResult
					continue
				}

				// update the object attribute
				conditionAtt := map[string]interface{}{}
				attr, err := data.MapStr(objID)
				if nil != err {
					blog.Error("can not convert to map, error info is %s", mapErr.Error())
					subResult["errors"] = defErr.Errorf(common.CCErrCommParamsLostField, "attr")
					result[objID] = subResult
					continue
				}

				for keyIdx := range attr {

					colIdx, err := strconv.ParseInt(keyIdx, 10, 64)
					if nil != err {
						blog.Errorf("the attribute index(%d) is invalid, error info is %s", colIdx, err.Error())
						continue
					}

					propertyGroupName := ""
					var err error
					grpName, _ := attr.MapStr(keyIdx)
					propertyGroupName, err := grpName.String("bk_property_group_name")
					if nil == err {
						// check group name
						grpName.Remove("bk_property_group_name")
						if nil != err {
							blog.Error("failed to parse the bk_property_group_name, error info is %s", err.Error())
							errStr := defLang.Languagef("import_row_int_error_str", colIdx, defErr.Errorf(common.CCErrCommParamsNeedString, "bk_property_group_name"))
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

					// check group name
					if 0 == len(propertyGroupName) {
						jsObjAttr.Get(keyIdx).Set("bk_property_group", "default") // set default, if set nothing
					} else {
						data := map[string]interface{}{
							common.BKOwnerIDField: ownerID,
							common.BKObjIDField:   objID,
							"bk_group_name":       propertyGroupName,
						}
						dataStr, _ := json.Marshal(data)

						grps, err := cli.mgr.SelectPropertyGroupByObjectID(forward, ownerID, objID, dataStr, defErr)
						if nil != err {
							blog.Error("failed to search the group, error info is %s", err.Error())
							errStr := defLang.Languagef("import_row_int_error_str", colIdx, defErr.Errorf(common.CCErrCommParamsNeedString, "bk_property_group_name"))
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

						if 0 != len(grps) {
							jsObjAttr.Get(keyIdx).Set("bk_property_group", grps[0].GroupID) // only one group, not any more
						} else {
							grp := api.ObjAttGroupDes{}
							grp.ObjectID = objID
							grp.OwnerID = ownerID
							grp.GroupID = xid.New().String()
							grp.GroupName = propertyGroupName
							grpStr, _ := json.Marshal(grp)
							if _, err := cli.mgr.CreateObjectGroup(forward, grpStr, defErr); nil != err {
								blog.Error("failed to create the group, error info is %s", err.Error())
								errStr := defLang.Languagef("import_row_int_error_str", colIdx, defErr.Error(common.CCErrTopoObjectGroupCreateFailed))
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

							jsObjAttr.Get(keyIdx).Set("bk_property_group", grp.GroupID) // only one group, not any more
						}

					}

					// check base attribute
					propertyID, err := jsObjAttr.Get(keyIdx).Get("bk_property_id").String()
					if 0 == len(propertyID) {
						blog.Error("not set the bk_property_id")
						errStr := defLang.Languagef("import_row_int_error_str", colIdx, defErr.Errorf(common.CCErrCommParamsNeedSet, "bk_property_id"))
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
					if nil != err {
						blog.Error("failed to parse the bk_property_id, error info is %s", err.Error())
						errStr := defLang.Languagef("import_row_int_error_str", colIdx, defErr.Errorf(common.CCErrCommParamsNeedString, "bk_property_id"))
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

					// check the property id
					conditionAtt[common.BKOwnerIDField] = ownerID
					conditionAtt[common.BKObjIDField] = objID
					conditionAtt["bk_property_id"] = propertyID

					conditionAttVal, _ := json.Marshal(conditionAtt)
					if items, err := cli.mgr.SelectObjectAtt(forward, conditionAttVal, defErr); nil != err {
						blog.Error("failed to search the object attribute, the condition is %+v, error info is %s", conditionAtt, err.Error())
						errStr := defLang.Languagef("import_row_int_error_str", colIdx, err.Error())
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

					} else if 0 != len(items) {

						// need to update
						for _, tmpItem := range items {

							item, itemErr := cli.updateObjectAttribute(&tmpItem, jsObjAttr.Get(keyIdx), defErr)
							if nil != itemErr {
								blog.Error("failed to reset the object attribute, error info is %s ", itemErr.Error())
								errStr := defLang.Languagef("import_row_int_error_str", colIdx, itemErr.Error())
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

							itemVal, _ := json.Marshal(item)
							blog.Debug("the new attribute:%s", string(itemVal))
							if updateErr := cli.mgr.UpdateObjectAtt(forward, item.ID, itemVal, defErr); nil != updateErr {
								blog.Error("failed to update the object attribute, error info is %s", updateErr.Error())
								errStr := defLang.Languagef("import_row_int_error_str", colIdx, updateErr.Error())
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

					} else {
						// need to create
						tmpItem := &api.ObjAttDes{}
						tmpItem.ObjectID = objID
						tmpItem.OwnerID = ownerID
						item, itemErr := cli.updateObjectAttribute(tmpItem, jsObjAttr.Get(keyIdx), defErr)
						if nil != itemErr {
							blog.Error("failed to reset the object attribute, error info is %s ", itemErr.Error())
							errStr := defLang.Languagef("import_row_int_error_str", colIdx, itemErr.Error())
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

						if _, insertErr := cli.mgr.CreateObjectAtt(forward, *item, defErr); nil != insertErr {
							blog.Error("failed to create the object attribute, error info is %s", insertErr.Error())
							errStr := defLang.Languagef("import_row_int_error_str", colIdx, insertErr.Error())
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

					} // end else  create attribute

					if failed, ok := subResult["success"]; ok {
						failedArr := failed.([]string)
						failedArr = append(failedArr, keyIdx)
						subResult["success"] = failedArr
					} else {
						subResult["success"] = []string{
							keyIdx,
						}
					}

					result[objID] = subResult

				} // end foreach objid
			}
	*/
	return nil
}
func (o *object) FindObjectBatch(params types.ContextParams, data frtypes.MapStr) error {
	return nil
}

func (o *object) FindSingleObject(params types.ContextParams, objectID string) (model.Object, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	cond.Field(common.BKObjIDField).Eq(objectID)

	objs, err := o.FindObject(params, cond)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the supplier account(%s) objects(%s), error info is %s", params.SupplierAccount, objectID, err.Error())
		return nil, err
	}
	for _, item := range objs {
		return item, nil
	}
	return nil, params.Err.New(common.CCErrTopoObjectSelectFailed, params.Err.Errorf(common.CCErrCommParamsIsInvalid, objectID).Error())
}
func (o *object) CreateObject(params types.ContextParams, data frtypes.MapStr) (model.Object, error) {
	obj := o.modelFactory.CreaetObject(params)

	_, err := obj.Parse(data)
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

	// check repeated

	exists, err := obj.IsExists()
	if nil != err {
		blog.Errorf("[operation-obj] failed to create the object(%#v), error info is %s", data, err.Error())
		return nil, params.Err.New(common.CCErrTopoObjectCreateFailed, err.Error())
	}

	if exists {
		blog.Errorf("[operation-obj] the object(%#v) is repeated", data)
		return nil, params.Err.Error(common.CCErrCommDuplicateItem)
	}

	err = obj.Create()
	if nil != err {
		blog.Errorf("[operation-obj] failed to save the data(%#v), error info is %s", data, err.Error())
		return nil, err
	}

	// create the default group
	grp := obj.CreateGroup()
	grp.SetDefault(true)
	grp.SetIndex(-1)
	grp.SetName("Default")
	grp.SetID("default")

	if err = grp.Create(); nil != err {
		blog.Errorf("[operation-obj] failed to create the default group, error info is %s", err.Error())
	}

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
		blog.Errorf("[operation-obj] failed to create the default inst name field, error info is %s", err.Error())
	}

	return obj, nil
}

func (o *object) DeleteObject(params types.ContextParams, id int64, cond condition.Condition) error {

	if 0 < id {
		cond = condition.CreateCondition()
		cond.Field(metadata.ModelFieldID).Eq(id)
	}

	objs, err := o.FindObject(params, cond)
	if nil != err {
		blog.Errorf("[operation-obj] failed to find objects, the condition is (%v) error info is %s", cond, err.Error())
		return err
	}

	for _, obj := range objs {

		attrCond := condition.CreateCondition()
		attrCond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
		attrCond.Field(common.BKObjIDField).Eq(obj.GetID())

		if err := o.attr.DeleteObjectAttribute(params, -1, attrCond); nil != err {
			blog.Errorf("[operation-obj] failed to delete the object(%d)'s attribute, error info is %s", id, err.Error())
			return err
		}

		rsp, err := o.clientSet.ObjectController().Meta().DeleteObject(context.Background(), obj.GetRecordID(), params.Header, cond.ToMapStr())

		if nil != err {
			blog.Errorf("[operation-obj] failed to request the object controller, error info is %s", err.Error())
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
		blog.Errorf("[operation-obj] failed to find object, error info is %s", err.Error())
		return nil, err
	}

	results := []metadata.ObjectTopo{}
	for _, obj := range objs {

		asstItems, err := o.asst.SearchObjectAssociation(params, obj.GetID())
		if nil != err {
			return nil, err
		}

		for _, asst := range asstItems {

			if asst.ObjectAttID == common.BKChildStr {
				continue
			}

			cond = condition.CreateCondition()
			cond.Field(common.BKObjIDField).Eq(asst.AsstObjID)
			cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)

			asstObjs, err := o.FindObject(params, cond)
			if nil != err {
				blog.Errorf("[operation-obj] failed to find object, error info is %s", err.Error())
				return nil, err
			}

			for _, asstObj := range asstObjs {
				tmp := metadata.ObjectTopo{}
				tmp.Label = asst.ObjectAttID
				tmp.LabelName = asst.AsstName
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
				ok, err := o.isFrom(params, obj.GetID(), asstObj.GetID())
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
		blog.Errorf("[operation-obj] failed to request the object controller, error info is %s", err.Error())
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-obj] failed to search the objects by the condition(%#v) , error info is %s", cond.ToMapStr(), rsp.ErrMsg)
		return nil, params.Err.Error(rsp.Code)
	}

	return model.CreateObject(params, o.clientSet, rsp.Data), nil
}

func (o *object) UpdateObject(params types.ContextParams, data frtypes.MapStr, id int64, cond condition.Condition) error {

	obj := o.modelFactory.CreaetObject(params)
	obj.SetRecordID(id)
	_, err := obj.Parse(data)
	if nil != err {
		blog.Errorf("[operation-obj] failed to parse the data(%#v), error info is %s", data, err.Error())
		return err
	}

	if err = obj.Update(data); nil != err {
		blog.Errorf("[operation-obj] failed to update the object(%d), the new data(%#v), error info is %s", id, data, err.Error())
		return params.Err.New(common.CCErrTopoObjectUpdateFailed, err.Error())
	}

	return nil
}
