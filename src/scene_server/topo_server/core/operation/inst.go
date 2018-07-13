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
	"strings"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	metatype "configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// InstOperationInterface inst operation methods
type InstOperationInterface interface {
	CreateInst(params types.ContextParams, obj model.Object, data frtypes.MapStr) (inst.Inst, error)
	DeleteInst(params types.ContextParams, obj model.Object, cond condition.Condition) error
	FindInst(params types.ContextParams, obj model.Object, cond *metatype.QueryInput) (count int, results []inst.Inst, err error)
	UpdateInst(params types.ContextParams, data frtypes.MapStr, obj model.Object, cond condition.Condition) error

	SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface)
}

// NewInstOperation create a new inst operation instance
func NewInstOperation(client apimachinery.ClientSetInterface) InstOperationInterface {
	return &commonInst{
		clientSet: client,
	}
}

type FieldName string
type AssociationObjectID string
type RowIndex int
type InputKey string
type InstID int64

type asstObjectAttribute struct {
	obj   model.Object
	attrs []model.Attribute
}

type batchResult struct {
	Errors       []string `json:"error"`
	Success      []string `json:"success"`
	UpdateErrors []string `json:"update_error"`
}

type commonInst struct {
	clientSet    apimachinery.ClientSetInterface
	modelFactory model.Factory
	instFactory  inst.Factory
	asst         AssociationOperationInterface
}

func (c *commonInst) SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface) {
	c.modelFactory = modelFactory
	c.instFactory = instFactory
	c.asst = asst
}

func (c *commonInst) getAsstObjectPrimaryFields(obj model.Object) ([]model.Attribute, map[FieldName]asstObjectAttribute, error) {

	fields, err := obj.GetAttributes()
	if nil != err {
		blog.Errorf("[operation-inst] failed to get the object(%s)'s fields, error info is %s", obj.GetID(), err.Error())
		return nil, nil, err
	}

	asstPrimaryFields := make(map[FieldName]asstObjectAttribute, 0)
	for _, fieldValue := range fields {

		if !fieldValue.IsAssociationType() {
			continue
		}

		asstObjects, err := obj.GetParentObjectByFieldID(fieldValue.GetID())
		if nil != err {
			blog.Errorf("[operation-inst] failed to get the object(%s)'s objects, error info is %s", obj.GetID(), err.Error())
			return nil, nil, err
		}

		for _, asstObj := range asstObjects {

			asstFields, err := asstObj.GetAttributes()
			if nil != err {
				blog.Errorf("[operation-inst] failed to get the object(%s)'s fields, error info is %s", obj.GetID(), err.Error())
				return nil, nil, err
			}

			asstPrimaryField := asstObjectAttribute{
				attrs: make([]model.Attribute, 0),
			}
			asstPrimaryField.obj = asstObj
			for _, field := range asstFields {
				if field.GetIsOnly() && field.IsAssociationType() {
					asstPrimaryField.attrs = append(asstPrimaryField.attrs, field)
				}
			}

			asstPrimaryFields[FieldName(fieldValue.GetID())] = asstPrimaryField
		}

	}
	return fields, asstPrimaryFields, nil
}

func (c *commonInst) constructAssociationInstSearchCondition(params types.ContextParams, fields []model.Attribute, asstPrimaryFields map[FieldName]asstObjectAttribute, batch *instBatchInfo) (map[AssociationObjectID][]frtypes.MapStr, map[RowIndex]error, error) {

	results := make(map[AssociationObjectID][]frtypes.MapStr)
	errs := make(map[RowIndex]error, 0)
	for rowIdx, dataVal := range *batch.BatchInfo {

		batchData, err := frtypes.NewFromInterface(dataVal)
		if nil != err {
			blog.Errorf("[operation-inst] failed to parse the data(%#v), error info is %s", dataVal, err.Error())
			return nil, nil, err
		}

		batchData.ForEach(func(key string, val interface{}) {

			for _, fieldValue := range fields {

				if fieldValue.GetID() != key {
					continue
				}

				if !fieldValue.IsAssociationType() {
					continue
				}

				asstFields, exists := asstPrimaryFields[FieldName(fieldValue.GetID())]
				if !exists {
					errs[RowIndex(rowIdx)] = params.Err.New(common.CCErrTopoInstCreateFailed, params.Lang.Languagef("import_asst_property_str_not_found", key))
					return
				}

				valStr, ok := val.(string)
				if !ok {
					errs[RowIndex(rowIdx)] = params.Err.New(common.CCErrTopoInstCreateFailed, params.Lang.Languagef("import_asst_property_str_not_found", key))
					return
				}

				if common.ExcelDelAsstObjectRelation == strings.TrimSpace(valStr) {
					continue
				}

				valStrItems := strings.Split(valStr, common.ExcelAsstPrimaryKeyRowChar)
				asstConds := make([]frtypes.MapStr, 0)
				for _, valStrItem := range valStrItems {

					if 0 == len(valStrItem) {
						continue
					}

					primaryKeys := strings.Split(valStrItem, common.ExcelAsstPrimaryKeySplitChar)

					if len(primaryKeys) != len(asstFields.attrs) {
						errs[RowIndex(rowIdx)] = params.Err.New(common.CCErrTopoInstCreateFailed, params.Lang.Languagef("import_asst_property_str_primary_count_len", key))
						continue
					}

					conds := frtypes.New()
					if asstFields.obj.IsCommon() {
						conds.Set(common.BKObjIDField, asstFields.obj.GetID())
					}

					for _, inputVal := range primaryKeys {

						for _, attr := range asstFields.attrs {

							if attr.GetID() != inputVal {
								continue
							}

							var err error
							conds[attr.GetID()], err = ConvByPropertytype(attr, inputVal)
							if nil != err {
								errs[RowIndex(rowIdx)] = params.Err.New(common.CCErrTopoInstCreateFailed, params.Lang.Languagef("import_asst_property_str_primary_count_len", key))
								continue
							}
						}

					} // end foreach primaryKeys
					asstConds = append(asstConds, conds)
				} // end foreach valStrItems

				if _, exists := results[AssociationObjectID(asstFields.obj.GetID())]; exists {
					results[AssociationObjectID(asstFields.obj.GetID())] = append(results[AssociationObjectID(asstFields.obj.GetID())], asstConds...)
					continue
				}

				results[AssociationObjectID(asstFields.obj.GetID())] = asstConds

			} // end for fields

		}) // end for dataVal

	}

	return results, errs, nil
}

func (c *commonInst) searchAssociationInstByConditions(params types.ContextParams, fields []model.Attribute, asstPrimaryFields map[FieldName]asstObjectAttribute, conds map[AssociationObjectID][]frtypes.MapStr) (map[AssociationObjectID]map[InputKey]InstID, error) {

	results := make(map[AssociationObjectID]map[InputKey]InstID)
	for asstObjID, asstCond := range conds {

		var target *asstObjectAttribute
	labelFor:
		for _, obj := range asstPrimaryFields {
			if obj.obj.GetID() != string(asstObjID) {
				target = &obj
				break labelFor
			}

		} // end foreach asst primary fields

		if nil != target {
			continue
		}

		query := &metatype.QueryInput{}
		query.Condition = asstCond
		_, insts, err := c.FindInst(params, target.obj, query)
		if nil != err {
			blog.Errorf("[operation-inst] failed to find the object(%s)'s insts, error info is %s", target.obj.GetID(), err.Error())
			return nil, err
		}

		for _, asstInst := range insts {
			inputKeys := []string{}
			instsMapStr := asstInst.GetValues()
			for _, attr := range target.attrs {
				attrVal, exists := instsMapStr.Get(attr.GetID())
				if !exists {
					return nil, params.Err.New(common.CCErrTopoInstCreateFailed, params.Lang.Languagef("import_str_asst_str_query_data_format_error", target.obj.GetID(), attr.GetID()))
				}
				inputKeys = append(inputKeys, fmt.Sprintf("%v", attrVal))
			} // end foreach target.attrs

			if _, ok := results[AssociationObjectID(target.obj.GetID())]; !ok {
				results[AssociationObjectID(target.obj.GetID())] = make(map[InputKey]InstID)
			}

			instID, err := asstInst.GetInstID()
			if nil != err {
				blog.Errorf("[operation-inst] failed to get the inst (%#v) id", instsMapStr)
				return nil, params.Err.New(common.CCErrTopoInstCreateFailed, err.Error())
			}

			results[AssociationObjectID(target.obj.GetID())][InputKey(strings.Join(inputKeys, common.ExcelAsstPrimaryKeySplitChar))] = InstID(instID)

		} // end foreach insts

	}

	return results, nil
}

func (c *commonInst) importExcelData(params types.ContextParams, obj model.Object, batch *instBatchInfo) (map[AssociationObjectID]map[InputKey]InstID, map[RowIndex]error, error) {

	fields, asstPrimaryFields, err := c.getAsstObjectPrimaryFields(obj)
	if nil != err {
		return nil, nil, err
	}

	searchAsstInstConds, errs, err := c.constructAssociationInstSearchCondition(params, fields, asstPrimaryFields, batch)
	if nil != err {
		return nil, errs, err
	}
	if 0 != len(errs) {
		return nil, errs, params.Err.New(common.CCErrTopoInstCreateFailed, "some exceptions")
	}
	insts, err := c.searchAssociationInstByConditions(params, fields, asstPrimaryFields, searchAsstInstConds)
	if nil != err {
		return nil, nil, err
	}
	return insts, nil, nil

}

func (c *commonInst) dealBatchImportInsts(params types.ContextParams, rowErrs map[RowIndex]error, obj model.Object, importAsstInsts map[AssociationObjectID]map[InputKey]InstID, batch frtypes.MapStr) error {

	asstObjs, err := obj.GetParentObject()
	if nil != err {
		blog.Errorf("[operation-inst] faild to find the association object, error info is %s", err.Error())
		return nil
	}

	for _, asstObj := range asstObjs {

		asstAttr, err := asstObj.GetAttributes()
		if nil != err {
			blog.Errorf("[operation-inst] not found the association attributes, error info is %s", err.Error())
			return err
		}

		asstKeyIDS, ok := importAsstInsts[AssociationObjectID(asstObj.GetID())]
		if !ok {
			blog.Errorf("[operation-inst] not found the association object(%s)", asstObj.GetID())
			return params.Err.New(common.CCErrCommParamsIsInvalid, asstObj.GetID())
		}

		for _, asstAttr := range asstAttr {

			if !asstAttr.IsAssociationType() {
				continue
			}

			var strIDS []string
			strInst, err := batch.String(asstAttr.GetID())
			if nil != err {
				blog.Errorf("[operation-inst] the asst key(%s) is invalid, error info is %s", asstAttr.GetID(), err.Error())
				return params.Err.New(common.CCErrCommParamsInvalid, asstAttr.GetID())
			}

			if common.ExcelDelAsstObjectRelation == strings.TrimSpace(strInst) {
				batch.Set(asstAttr.GetID(), "")
				continue
			}

			strAsstInstKeys := strings.Split(strInst, common.ExcelAsstPrimaryKeyRowChar)

			for _, asstKey := range strAsstInstKeys {
				asstKey = strings.TrimSpace(asstKey)
				if 0 == len(asstKey) {
					continue
				}
				if id, ok := asstKeyIDS[InputKey(asstKey)]; ok {
					strIDS = append(strIDS, fmt.Sprintf("%d", id))
				}
			}

			batch.Set(asstAttr.GetID(), strings.Join(strIDS, common.InstAsstIDSplit))

		}

	} // end for

	return nil
}

func (c *commonInst) CreateInst(params types.ContextParams, obj model.Object, data frtypes.MapStr) (inst.Inst, error) {

	// extract internal data
	batchInfo := &instBatchInfo{}
	err := data.MarshalJSONInto(batchInfo)
	if nil != err {
		blog.Errorf("[operation-inst] failed to unmarshal the data(%#v) into the inst batch info struct, error info is %s", data, err.Error())
		return nil, params.Err.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	// create association
	var importInsts map[AssociationObjectID]map[InputKey]InstID
	var rowsErrors map[RowIndex]error
	results := &batchResult{}
	if common.InputTypeExcel == batchInfo.InputType && nil != batchInfo.BatchInfo {
		importInsts, rowsErrors, err = c.importExcelData(params, obj, batchInfo)
		if nil != err {
			return nil, err
		}

		for colIdx, colInput := range *batchInfo.BatchInfo {
			delete(colInput, "import_from")
			err, ok := rowsErrors[RowIndex(colIdx)]
			if ok {
				results.Errors = append(results.Errors, params.Lang.Languagef("import_row_int_error_str", colIdx, err.Error()))
				continue
			}

			err = c.dealBatchImportInsts(params, rowsErrors, obj, importInsts, colInput)
			if nil != err {
				results.Errors = append(results.Errors, params.Lang.Languagef("import_row_int_error_str", colIdx, err.Error()))
				continue
			}

			item := c.instFactory.CreateInst(params, obj)

			item.SetValues(colInput)
			err = item.Create()
			if nil != err {
				blog.Errorf("[operation-inst] failed to save the object(%s) inst data (%#v), error info is %s", obj.GetID(), data, err.Error())
				return nil, err
			}

			return item, nil
		} // end foreach batchinfo
	}

	// create new insts
	blog.Infof("the data inst:%#v", data)
	item := c.instFactory.CreateInst(params, obj)

	item.SetValues(data)

	err = item.Create()
	if nil != err {
		blog.Errorf("[operation-inst] failed to save the object(%s) inst data (%#v), error info is %s", obj.GetID(), data, err.Error())
		return nil, err
	}

	return item, nil
}

func (c *commonInst) DeleteInst(params types.ContextParams, obj model.Object, cond condition.Condition) error {

	rsp, err := c.clientSet.ObjectController().Instance().DelObject(context.Background(), obj.GetObjectType(), params.Header, cond.ToMapStr())

	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-inst] faild to delete the object(%s) inst by the condition(%#v), error info is %s", obj.GetID(), cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}

	return nil
}

func (c *commonInst) FindInst(params types.ContextParams, obj model.Object, cond *metatype.QueryInput) (count int, results []inst.Inst, err error) {

	rsp, err := c.clientSet.ObjectController().Instance().SearchObjects(context.Background(), obj.GetObjectType(), params.Header, cond)

	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, error info is %s", err.Error())
		return 0, nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-inst] faild to delete the object(%s) inst by the condition(%#v), error info is %s", obj.GetID(), cond, rsp.ErrMsg)
		return 0, nil, params.Err.Error(rsp.Code)
	}

	return rsp.Data.Count, inst.CreateInst(params, c.clientSet, obj, rsp.Data.Info), nil
}

func (c *commonInst) UpdateInst(params types.ContextParams, data frtypes.MapStr, obj model.Object, cond condition.Condition) error {

	inputParams := frtypes.New()
	inputParams.Set("data", data)
	inputParams.Set("condition", cond.ToMapStr())
	blog.Infof("data condition:%#v", inputParams)
	rsp, err := c.clientSet.ObjectController().Instance().UpdateObject(context.Background(), obj.GetObjectType(), params.Header, inputParams)

	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-inst] faild to set the object(%s) inst by the condition(%#v), error info is %s", obj.GetID(), cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}
	return nil
}
