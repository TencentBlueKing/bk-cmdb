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
	"strings"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	metatype "configcenter/src/common/metadata"
	gparams "configcenter/src/common/paraparse"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// InstOperationInterface inst operation methods
type InstOperationInterface interface {
	CreateInst(params types.ContextParams, obj model.Object, data frtypes.MapStr) (inst.Inst, error)
	DeleteInst(params types.ContextParams, obj model.Object, cond condition.Condition) error
	DeleteInstByInstID(params types.ContextParams, obj model.Object, instID []int64) error
	FindInst(params types.ContextParams, obj model.Object, cond *metatype.QueryInput, needAsstDetail bool) (count int, results []inst.Inst, err error)
	FindInstByAssociationInst(params types.ContextParams, obj model.Object, data frtypes.MapStr) (cont int, results []inst.Inst, err error)
	FindInstChildTopo(params types.ContextParams, obj model.Object, instID int64, query *metatype.QueryInput) (count int, results []interface{}, err error)
	FindInstParentTopo(params types.ContextParams, obj model.Object, instID int64, query *metatype.QueryInput) (count int, results []interface{}, err error)
	FindInstTopo(params types.ContextParams, obj model.Object, instID int64, query *metatype.QueryInput) (count int, results []commonInstTopoV2, err error)
	UpdateInst(params types.ContextParams, data frtypes.MapStr, obj model.Object, cond condition.Condition, instID int64) error

	SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface, obj ObjectOperationInterface)
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
	obj          ObjectOperationInterface
}

func (c *commonInst) SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface, obj ObjectOperationInterface) {
	c.modelFactory = modelFactory
	c.instFactory = instFactory
	c.asst = asst
	c.obj = obj
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
		_, insts, err := c.FindInst(params, target.obj, query, false)
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
			if err = NewSupplementary().Validator(c).ValidatorCreate(params, obj, item.ToMapStr()); nil != err {
				return nil, err
			}
			err = item.Create()
			if nil != err {
				blog.Errorf("[operation-inst] failed to save the object(%s) inst data (%#v), error info is %s", obj.GetID(), data, err.Error())
				return nil, err
			}
			NewSupplementary().Audit(params, c.clientSet, item.GetObject(), c).CommitCreateLog(nil, nil, item)
			return item, nil
		} // end foreach batchinfo
	}

	// create new insts
	blog.Infof("the data inst:%#v", data)
	item := c.instFactory.CreateInst(params, obj)
	item.SetValues(data)
	if err = NewSupplementary().Validator(c).ValidatorCreate(params, obj, item.ToMapStr()); nil != err {
		blog.Errorf("[operation-inst] valid is bad, the data is (%#v)  error info is %s", item.ToMapStr(), err.Error())
		return nil, err
	}
	err = item.Create()
	if nil != err {
		blog.Errorf("[operation-inst] failed to save the object(%s) inst data (%#v), error info is %s", obj.GetID(), data, err.Error())
		return nil, err
	}

	NewSupplementary().Audit(params, c.clientSet, item.GetObject(), c).CommitCreateLog(nil, nil, item)

	return item, nil
}

func (c *commonInst) innerHasHost(params types.ContextParams, moduleIDS []int64) (bool, error) {
	cond := map[string][]int64{
		common.BKModuleIDField: moduleIDS,
	}

	rsp, err := c.clientSet.HostController().Module().GetModulesHostConfig(context.Background(), params.Header, cond)
	if nil != err {
		blog.Errorf("[operation-module] failed to request the object controller, error info is %s", err.Error())
		return false, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-module]  failed to search the host module configures, error info is %s", err.Error())
		return false, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != len(rsp.Data), nil
}
func (c *commonInst) hasHost(params types.ContextParams, targetInst inst.Inst) ([]deletedInst, bool, error) {

	id, err := targetInst.GetInstID()
	if nil != err {
		return nil, false, err
	}

	targetObj := targetInst.GetObject()
	if !targetObj.IsCommon() {
		if targetObj.GetObjectType() == common.BKInnerObjIDModule {
			exists, err := c.innerHasHost(params, []int64{id})
			if nil != err {
				return nil, false, err
			}

			if exists {
				return nil, true, nil
			}
		}
	}

	instIDS := []deletedInst{}
	instIDS = append(instIDS, deletedInst{instID: id, obj: targetObj})
	childInsts, err := targetInst.GetChildInst()
	if nil != err {
		return nil, false, err
	}

	for _, childInst := range childInsts {

		ids, exists, err := c.hasHost(params, childInst)
		if nil != err {
			return nil, false, err
		}
		if exists {
			return instIDS, true, nil
		}
		instIDS = append(instIDS, ids...)
	}

	return instIDS, false, nil
}

func (c *commonInst) DeleteInstByInstID(params types.ContextParams, obj model.Object, instID []int64) error {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	cond.Field(obj.GetInstIDFieldName()).In(instID)
	if obj.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(obj.GetID())
	}

	query := &metatype.QueryInput{}
	query.Condition = cond.ToMapStr()

	_, insts, err := c.FindInst(params, obj, query, false)
	if nil != err {
		return err
	}

	deleteIDS := []deletedInst{}
	for _, inst := range insts {
		ids, exists, err := c.hasHost(params, inst)
		if nil != err {
			return params.Err.Error(common.CCErrTopoHasHostCheckFailed)
		}

		if exists {
			return params.Err.Error(common.CCErrTopoHasHostCheckFailed)
		}

		deleteIDS = append(deleteIDS, ids...)
	}

	for _, delInst := range deleteIDS {

		cond = condition.CreateCondition()
		cond.Field(common.BKObjIDField).Eq(delInst.obj.GetID())
		cond.Field(common.BKInstIDField).Eq(delInst.instID)
		if err = c.asst.DeleteInstAssociation(params, cond); nil != err {
			return err
		}

		cond = condition.CreateCondition()
		cond.Field(common.BKAsstInstIDField).Eq(delInst.instID)
		cond.Field(common.BKAsstObjIDField).Eq(delInst.obj.GetID())
		if err = c.asst.DeleteAssociation(params, cond); nil != err {
			return err
		}

		cond = condition.CreateCondition()
		cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
		cond.Field(obj.GetInstIDFieldName()).Eq(delInst.instID)

		if err = c.DeleteInst(params, delInst.obj, cond); nil != err {
			return err
		}

	}

	return nil

}
func (c *commonInst) DeleteInst(params types.ContextParams, obj model.Object, cond condition.Condition) error {

	preAudit := NewSupplementary().Audit(params, c.clientSet, obj, c).CreateSnapshot(-1, cond.ToMapStr())

	rsp, err := c.clientSet.ObjectController().Instance().DelObject(context.Background(), obj.GetObjectType(), params.Header, cond.ToMapStr())

	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-inst] faild to delete the object(%s) inst by the condition(%#v), error info is %s", obj.GetID(), cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}

	NewSupplementary().Audit(params, c.clientSet, obj, c).CommitDeleteLog(preAudit, nil, nil)

	return nil
}
func (c *commonInst) convertInstIDIntoStruct(params types.ContextParams, asstObj metatype.Association, instIDS []string, needAsstDetail bool) ([]instNameAsst, error) {

	obj, err := c.obj.FindSingleObject(params, asstObj.AsstObjID)
	if nil != err {
		return nil, err
	}

	ids := []int64{}
	for _, id := range instIDS {
		idbit, err := strconv.ParseInt(id, 10, 64)
		if nil != err {
			return nil, err
		}

		ids = append(ids, idbit)
	}

	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).In(ids)

	query := &metatype.QueryInput{}
	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit
	rsp, err := c.clientSet.ObjectController().Instance().SearchObjects(context.Background(), obj.GetObjectType(), params.Header, query)

	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, error info is %s", err.Error())
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-inst] faild to delete the object(%s) inst by the condition(%#v), error info is %s", obj.GetID(), cond, rsp.ErrMsg)
		return nil, params.Err.Error(rsp.Code)
	}

	instAsstNames := []instNameAsst{}
	for _, instInfo := range rsp.Data.Info {
		instName, err := instInfo.String(obj.GetInstNameFieldName())
		if nil != err {
			return nil, err
		}
		instID, err := instInfo.Int64(obj.GetInstIDFieldName())
		if nil != err {
			return nil, err
		}

		if needAsstDetail {
			instAsstNames = append(instAsstNames, instNameAsst{
				ObjID:      obj.GetID(),
				ObjectName: obj.GetName(),
				ObjIcon:    obj.GetIcon(),
				InstID:     instID,
				InstName:   instName,
				InstInfo:   instInfo,
			})
			continue
		}

		instAsstNames = append(instAsstNames, instNameAsst{
			ObjID:      obj.GetID(),
			ObjectName: obj.GetName(),
			ObjIcon:    obj.GetIcon(),
			InstID:     instID,
			InstName:   instName,
		})

	}

	return instAsstNames, nil
}

func (c *commonInst) searchAssociationInst(params types.ContextParams, objID string, searchCond frtypes.MapStr) ([]int64, error) {

	obj, err := c.obj.FindSingleObject(params, objID)
	if nil != err {
		return nil, err
	}

	query := &metatype.QueryInput{}
	query.Condition = searchCond

	_, insts, err := c.FindInst(params, obj, query, false)
	if nil != err {
		return nil, err
	}
	//fmt.Println("search cond:", searchCond, obj.GetID(), insts)

	instIDS := make([]int64, 0)
	for _, inst := range insts {
		id, err := inst.GetInstID()
		if nil != err {
			return nil, err
		}
		instIDS = append(instIDS, id)
	}

	return instIDS, nil
}

func (c *commonInst) FindInstChildTopo(params types.ContextParams, obj model.Object, instID int64, query *metatype.QueryInput) (count int, results []interface{}, err error) {
	results = []interface{}{}
	if nil == query {
		query = &metatype.QueryInput{}
		cond := condition.CreateCondition()
		cond.Field(obj.GetInstIDFieldName()).Eq(instID)
		cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
		query.Condition = cond.ToMapStr()
	}

	_, insts, err := c.FindInst(params, obj, query, false)
	if nil != err {
		return 0, nil, err
	}

	tmpResults := map[string]*commonInstTopo{}
	for _, inst := range insts {

		childInsts, err := inst.GetChildInst()
		if nil != err {
			return 0, nil, err
		}

		for _, childInst := range childInsts {

			childObj := childInst.GetObject()
			commonInst, exists := tmpResults[childObj.GetID()]
			if !exists {
				commonInst = &commonInstTopo{}
				commonInst.ObjectName = childObj.GetName()
				commonInst.ObjIcon = childObj.GetIcon()
				commonInst.ObjID = childObj.GetID()
				tmpResults[childObj.GetID()] = commonInst
			}

			instAsst := instNameAsst{}
			id, err := childInst.GetInstID()
			if nil != err {
				return 0, nil, err
			}

			name, err := childInst.GetInstName()
			if nil != err {
				return 0, nil, err
			}

			instAsst.InstID = id
			instAsst.InstName = name
			instAsst.ObjectName = childInst.GetObject().GetName()
			instAsst.ObjIcon = childInst.GetObject().GetIcon()
			instAsst.ObjID = childInst.GetObject().GetID()

			commonInst.Children = append(commonInst.Children, instAsst)
		}
	}

	for _, subResult := range tmpResults {
		results = append(results, subResult)
	}

	return 0, results, nil
}

func (c *commonInst) FindInstParentTopo(params types.ContextParams, obj model.Object, instID int64, query *metatype.QueryInput) (count int, results []interface{}, err error) {

	results = []interface{}{}
	if nil == query {
		query = &metatype.QueryInput{}
		cond := condition.CreateCondition()
		cond.Field(obj.GetInstIDFieldName()).Eq(instID)
		cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
		query.Condition = cond.ToMapStr()
	}

	_, insts, err := c.FindInst(params, obj, query, false)
	if nil != err {
		return 0, nil, err
	}

	tmpResults := map[string]*commonInstTopo{}
	for _, inst := range insts {

		parentInsts, err := inst.GetParentInst()
		if nil != err {
			return 0, nil, err
		}

		for _, parentInst := range parentInsts {

			childObj := parentInst.GetObject()
			commonInst, exists := tmpResults[childObj.GetID()]
			if !exists {
				commonInst = &commonInstTopo{}
				commonInst.ObjectName = childObj.GetName()
				commonInst.ObjIcon = childObj.GetIcon()
				commonInst.ObjID = childObj.GetID()
				tmpResults[childObj.GetID()] = commonInst
			}

			instAsst := instNameAsst{}
			id, err := parentInst.GetInstID()
			if nil != err {
				return 0, nil, err
			}

			name, err := parentInst.GetInstName()
			if nil != err {
				return 0, nil, err
			}

			instAsst.InstID = id
			instAsst.InstName = name
			instAsst.ObjectName = parentInst.GetObject().GetName()
			instAsst.ObjIcon = parentInst.GetObject().GetIcon()
			instAsst.ObjID = parentInst.GetObject().GetID()

			commonInst.Children = append(commonInst.Children, instAsst)
		}
	}

	for _, subResult := range tmpResults {
		results = append(results, subResult)
	}

	return 0, results, nil
}

func (c *commonInst) FindInstTopo(params types.ContextParams, obj model.Object, instID int64, query *metatype.QueryInput) (count int, results []commonInstTopoV2, err error) {

	if nil == query {
		query = &metatype.QueryInput{}
		cond := condition.CreateCondition()
		cond.Field(obj.GetInstIDFieldName()).Eq(instID)
		cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
		query.Condition = cond.ToMapStr()
	}

	_, insts, err := c.FindInst(params, obj, query, false)
	if nil != err {
		blog.Errorf("[operation-inst] failed to find the inst, error info is %s", err.Error())
		return 0, nil, err
	}

	for _, inst := range insts {

		//fmt.Println("the insts:", inst.GetValues(), query)
		id, err := inst.GetInstID()
		if nil != err {
			blog.Errorf("[operation-inst] failed to find the inst, error info is %s", err.Error())
			return 0, nil, err
		}

		name, err := inst.GetInstName()
		if nil != err {
			blog.Errorf("[operation-inst] failed to find the inst, error info is %s", err.Error())
			return 0, nil, err
		}

		commonInst := commonInstTopo{Children: []instNameAsst{}}
		commonInst.ObjectName = inst.GetObject().GetName()
		commonInst.ObjID = inst.GetObject().GetID()
		commonInst.ObjIcon = inst.GetObject().GetIcon()
		commonInst.InstID = id
		commonInst.InstName = name

		_, parentInsts, err := c.FindInstParentTopo(params, inst.GetObject(), id, nil)
		if nil != err {
			blog.Errorf("[operation-inst] failed to find the inst, error info is %s", err.Error())
			return 0, nil, err
		}

		_, childInsts, err := c.FindInstChildTopo(params, inst.GetObject(), id, nil)
		if nil != err {
			blog.Errorf("[operation-inst] failed to find the inst, error info is %s", err.Error())
			return 0, nil, err
		}

		results = append(results, commonInstTopoV2{
			Prev: parentInsts,
			Next: childInsts,
			Curr: commonInst,
		})

	}

	return 0, results, nil
}

func (c *commonInst) FindInstByAssociationInst(params types.ContextParams, obj model.Object, data frtypes.MapStr) (cont int, results []inst.Inst, err error) {

	asstParamCond := &AssociationParams{}

	if err := data.MarshalJSONInto(asstParamCond); nil != err {
		blog.Errorf("[operation-inst] find inst by association inst , error info is %s", err.Error())
		return 0, nil, params.Err.New(common.CCErrTopoInstSelectFailed, err.Error())
	}

	instCond := map[string]interface{}{}
	instCond[common.BKOwnerIDField] = params.SupplierAccount
	if obj.IsCommon() {
		instCond[common.BKObjIDField] = obj.GetID()
	}
	targetInstIDS := []int64{}

	for keyObjID, objs := range asstParamCond.Condition {
		// Extract the ID of the instance according to the associated object.
		cond := map[string]interface{}{}
		if common.GetObjByType(keyObjID) == common.BKINnerObjIDObject {
			cond[common.BKObjIDField] = keyObjID
			cond[common.BKOwnerIDField] = params.SupplierAccount
		}

		for _, objCondition := range objs {

			if objCondition.Operator != common.BKDBEQ {

				if obj.GetID() == keyObjID {
					// deal self condition
					instCond[objCondition.Field] = map[string]interface{}{
						objCondition.Operator: objCondition.Value,
					}
				} else {
					// deal association condition
					cond[objCondition.Field] = map[string]interface{}{
						objCondition.Operator: objCondition.Value,
					}
				}
			} else {
				if obj.GetID() == keyObjID {
					// deal self condition
					switch t := objCondition.Value.(type) {
					case string:
						instCond[objCondition.Field] = map[string]interface{}{
							common.BKDBEQ: gparams.SpeceialCharChange(t),
						}
					default:
						instCond[objCondition.Field] = objCondition.Value
					}

				} else {
					// deal association condition
					cond[objCondition.Field] = objCondition.Value
				}
			}

		}

		if obj.GetID() == keyObjID {
			// no need to search the association objects
			continue
		}

		innerCond := frtypes.New()
		if fields, ok := asstParamCond.Fields[keyObjID]; ok {
			innerCond.Set("fields", strings.Join(fields, ","))

		} else {
			innerCond.Set("fields", "")
		}

		innerCond.Set("condition", cond)
		innerCond.Set("start", 0)
		innerCond.Set("limit", common.BKNoLimit)
		innerCond.Set("sort", "")
		//fmt.Println("input cond:", innerCond)
		asstInstIDS, err := c.searchAssociationInst(params, keyObjID, innerCond)
		if nil != err {
			blog.Errorf("[operation-inst]failed to search the association inst, error info is %s", err.Error())
			return 0, nil, err
		}

		input := map[string]interface{}{
			"page": map[string]interface{}{
				"start": 0,
				"sort":  "",
				"limit": common.BKNoLimit,
			},
			"condition": map[string]interface{}{
				"bk_asst_inst_id": map[string]interface{}{
					common.BKDBIN: asstInstIDS,
				},
				"bk_asst_obj_id": keyObjID,
				"bk_obj_id":      obj.GetID(),
			},
			"fields": "",
		}
		query := &metatype.QueryInput{}
		query.Condition = input

		asstInst, err := c.asst.SearchInstAssociation(params, obj.GetID(), keyObjID, query)
		if nil != err {
			blog.Errorf("[operation-inst] failed to search the association inst, error info is %s", err.Error())
			return 0, nil, err
		}

		for _, asst := range asstInst {
			targetInstIDS = append(targetInstIDS, asst.InstID)
		}

	} // end foreach conditions

	if 0 != len(targetInstIDS) {
		instCond[obj.GetInstIDFieldName()] = map[string]interface{}{
			common.BKDBIN: targetInstIDS,
		}
	} else if 0 != len(asstParamCond.Condition) {
		if _, ok := asstParamCond.Condition[obj.GetID()]; !ok {
			instCond[obj.GetInstIDFieldName()] = map[string]interface{}{
				common.BKDBIN: targetInstIDS,
			}
		}
	}

	//fmt.Println("the targetids:", targetInstIDS, instCond)

	query := &metatype.QueryInput{}
	query.Condition = instCond
	if fields, ok := asstParamCond.Fields[obj.GetID()]; ok {
		query.Fields = strings.Join(fields, ",")
	}
	query.Limit = asstParamCond.Page.Limit
	query.Sort = asstParamCond.Page.Sort
	query.Start = asstParamCond.Page.Start
	return c.FindInst(params, obj, query, false)
}

func (c *commonInst) FindInst(params types.ContextParams, obj model.Object, cond *metatype.QueryInput, needAsstDetail bool) (count int, results []inst.Inst, err error) {

	rsp, err := c.clientSet.ObjectController().Instance().SearchObjects(context.Background(), obj.GetObjectType(), params.Header, cond)

	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, error info is %s", err.Error())
		return 0, nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {

		blog.Errorf("[operation-inst] faild to delete the object(%s) inst by the condition(%#v), error info is %s", obj.GetID(), cond, rsp.ErrMsg)
		return 0, nil, params.Err.Error(rsp.Code)
	}

	asstObjAttrs, err := c.asst.SearchObjectAssociation(params, obj.GetID())
	if nil != err {
		blog.Errorf("[operation-inst] failed to search object associations, error info is %s", err.Error())
		return 0, nil, err
	}

	for idx, instInfo := range rsp.Data.Info {

		for _, attrAsst := range asstObjAttrs {
			if attrAsst.ObjectAttID == common.BKChildStr || attrAsst.ObjectAttID == common.BKInstParentStr {
				continue
			}

			if !instInfo.Exists(attrAsst.ObjectAttID) { // the inst data is old, but the attribute is new.
				continue
			}

			asstFieldValue, err := instInfo.String(attrAsst.ObjectAttID)
			if nil != err {
				blog.Errorf("[operation-inst] failed to get the inst'attr(%s) value int the data(%#v), error info is %s", attrAsst.ObjectAttID, instInfo, err.Error())
				return 0, nil, err
			}
			instVals, err := c.convertInstIDIntoStruct(params, attrAsst, strings.Split(asstFieldValue, ","), needAsstDetail)
			if nil != err {
				blog.Errorf("[operation-inst] failed to convert association asst(%#v) origin value(%#v) value(%s), error info is %s", attrAsst, instInfo, asstFieldValue, err.Error())
				return 0, nil, err
			}
			rsp.Data.Info[idx].Set(attrAsst.ObjectAttID, instVals)

		}
	}
	return rsp.Data.Count, inst.CreateInst(params, c.clientSet, obj, rsp.Data.Info), nil
}

func (c *commonInst) UpdateInst(params types.ContextParams, data frtypes.MapStr, obj model.Object, cond condition.Condition, instID int64) error {

	if err := NewSupplementary().Validator(c).ValidatorUpdate(params, obj, data, instID, cond); nil != err {
		return err
	}

	inputParams := frtypes.New()
	inputParams.Set("data", data)
	inputParams.Set("condition", cond.ToMapStr())
	preAuditLog := NewSupplementary().Audit(params, c.clientSet, obj, c).CreateSnapshot(-1, cond.ToMapStr())
	rsp, err := c.clientSet.ObjectController().Instance().UpdateObject(context.Background(), obj.GetObjectType(), params.Header, inputParams)

	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-inst] faild to set the object(%s) inst by the condition(%#v), error info is %s", obj.GetID(), cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}
	currAuditLog := NewSupplementary().Audit(params, c.clientSet, obj, c).CreateSnapshot(-1, cond.ToMapStr())
	NewSupplementary().Audit(params, c.clientSet, obj, c).CommitUpdateLog(preAuditLog, currAuditLog, nil)
	return nil
}
