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
	"strconv"
	"strings"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
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
	CreateInstBatch(params types.ContextParams, obj model.Object, batchInfo *InstBatchInfo) (*BatchResult, error)
	DeleteInst(params types.ContextParams, obj model.Object, cond condition.Condition, needCheckHost bool) error
	DeleteInstByInstID(params types.ContextParams, obj model.Object, instID []int64, needCheckHost bool) error
	FindOriginInst(params types.ContextParams, obj model.Object, cond *metatype.QueryInput) (*metatype.InstResult, error)
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

type BatchResult struct {
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

func (c *commonInst) CreateInstBatch(params types.ContextParams, obj model.Object, batchInfo *InstBatchInfo) (*BatchResult, error) {

	var rowErr map[int64]error
	results := &BatchResult{}
	if common.InputTypeExcel != batchInfo.InputType || nil == batchInfo.BatchInfo {
		return results, nil
	}

	assObjectInt := NewAsstObjectInst(params.Header, params.Engin, params.SupplierAccount, nil)
	assObjectInt.SetMapFields(obj.GetID())
	err := assObjectInt.GetObjAsstObjectPrimaryKey()
	if nil != err {
		blog.Error("failed to read the object att, error is %s ", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommSearchPropertyFailed, err.Error())
		//return fmt.Errorf("get host assocate object  property failure, error:%s", err.Error())
	}
	rowErr, err = assObjectInt.InitInstFromData(*batchInfo.BatchInfo)
	if nil != err {
		blog.Error("failed to read the object att, error is %s ", err.Error())
		return nil, params.Err.Error(common.CCErrTopoInstSelectFailed)
		//return fmt.Errorf("get host assocate object instance data failure, error:%s", err.Error()), nil, nil, nil
	}

	for errIdx, err := range rowErr {
		results.Errors = append(results.Errors, params.Lang.Languagef("import_row_int_error_str", errIdx, err.Error()))
	}

	for colIdx, colInput := range *batchInfo.BatchInfo {
		delete(colInput, "import_from")

		if err := assObjectInt.SetObjAsstPropertyVal(colInput); nil != err {
			results.Errors = append(results.Errors, params.Lang.Languagef("import_row_int_error_str", colIdx, err.Error()))
			continue
		}

		//fmt.Println("input:", colInput)

		item := c.instFactory.CreateInst(params, obj)
		item.SetValues(colInput)

		if item.GetValues().Exists(obj.GetInstIDFieldName()) {
			// check update
			targetInstID, err := item.GetInstID()
			if nil != err {
				blog.Errorf("[operation-inst] failed to get inst id, error info is %s", err.Error())
				results.Errors = append(results.Errors, params.Lang.Languagef("import_row_int_error_str", colIdx, err.Error()))
				continue
			}
			if err = NewSupplementary().Validator(c).ValidatorUpdate(params, obj, item.ToMapStr(), targetInstID, nil); nil != err {
				blog.Errorf("[operation-inst] failed to valid, error info is %s", err.Error())
				results.Errors = append(results.Errors, params.Lang.Languagef("import_row_int_error_str", colIdx, err.Error()))
				continue
			}

		} else {
			// check create
			if err = NewSupplementary().Validator(c).ValidatorCreate(params, obj, item.ToMapStr()); nil != err {
				switch tmpErr := err.(type) {
				case errors.CCErrorCoder:
					if tmpErr.GetCode() != common.CCErrCommDuplicateItem {
						blog.Errorf("[operation-inst] failed to valid, input value(%#v) the instname is %s, error info is %s", item.GetValues(), obj.GetInstNameFieldName(), err.Error())
						results.Errors = append(results.Errors, params.Lang.Languagef("import_row_int_error_str", colIdx, err.Error()))
						continue
					}
				}

			}
		}

		// set data
		err = item.Save(colInput)
		if nil != err {
			blog.Errorf("[operation-inst] failed to save the object(%s) inst data (%#v), error info is %s", obj.GetID(), colInput, err.Error())
			results.Errors = append(results.Errors, params.Lang.Languagef("import_row_int_error_str", colIdx, err.Error()))
			continue
		}
		NewSupplementary().Audit(params, c.clientSet, item.GetObject(), c).CommitCreateLog(nil, nil, item)
		if err := c.setInstAsst(params, obj, item); nil != err {
			blog.Errorf("[operation-inst] failed to set the inst asst, error info is %s", err.Error())
			return nil, err
		}
	} // end foreach batchinfo

	return results, nil
}

func (c *commonInst) isValidInstID(params types.ContextParams, obj metatype.Object, instID int64) error {

	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)
	if obj.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(obj.ObjectID)
	}

	query := &metatype.QueryInput{}
	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit

	rsp, err := c.clientSet.ObjectController().Instance().SearchObjects(context.Background(), obj.GetObjectType(), params.Header, query)
	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-inst] faild to delete the object(%s) inst by the condition(%#v), error info is %s", obj.ObjectID, cond, rsp.ErrMsg)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	if rsp.Data.Count > 0 {
		return nil
	}

	return params.Err.Error(common.CCErrTopoInstSelectFailed)
}

func (c *commonInst) setInstAsst(params types.ContextParams, obj model.Object, inst inst.Inst) error {

	currInstID, err := inst.GetInstID()
	if nil != err {
		return err
	}
	attrs, err := obj.GetAttributesExceptInnerFields()
	if nil != err {
		return err
	}

	assts, err := c.asst.SearchObjectAssociation(params, obj.GetID())
	if nil != err {
		return err
	}

	data := inst.GetValues()
	for _, attr := range attrs {
		if !attr.IsAssociationType() {
			continue
		}

		// extract the assocaition insts ids
		asstVal, exists := data.Get(attr.GetID())
		if !exists {
			continue
		}

		asstInstIDS := []int64{}
		switch targetAssts := asstVal.(type) {
		default:
		case string:
			tmpIDS := strings.Split(targetAssts, common.InstAsstIDSplit)
			for _, asstID := range tmpIDS {
				if 0 == len(strings.TrimSpace(asstID)) {
					continue
				}

				id, err := strconv.ParseInt(asstID, 10, 64)
				if nil != err {
					blog.Errorf("[operation-inst] failed to parse  the value(%s), error info is %s", asstID, err.Error())
					return err
				}
				asstInstIDS = append(asstInstIDS, id)
			}

		case []metatype.InstNameAsst:
			for _, item := range targetAssts {
				asstInstIDS = append(asstInstIDS, item.InstID)
			}
		}

		// update the inst association
		for _, asst := range assts {
			if asst.ObjectAttID != attr.GetID() {
				continue
			}

			// delete the inst asst
			innerCond := condition.CreateCondition()
			innerCond.Field(common.BKObjIDField).Eq(obj.GetID())
			innerCond.Field(common.BKAsstObjIDField).Eq(asst.AsstObjID)
			innerCond.Field(common.BKInstIDField).Eq(currInstID)
			if err = c.asst.DeleteInstAssociation(params, innerCond); nil != err {
				blog.Errorf("[operation-inst] failed to delete inst association,by the condition(%#v), error info is %s", innerCond.ToMapStr(), err.Error())
				return err
			}

			// create a new new asst
			validInstIDS := []string{}
			for _, asstInstID := range asstInstIDS {

				// check the inst id
				if err := c.isValidInstID(params, metatype.Object{ObjectID: asst.AsstObjID}, asstInstID); nil != err {
					blog.Warnf("[operation-int] the association object(%s)' instid(%d) is invalid", asst.AsstObjID, asstInstID)
					continue
				}
				validInstIDS = append(validInstIDS, strconv.Itoa(int(asstInstID)))
				// create a new inst in inst asst table
				if err = c.asst.CreateCommonInstAssociation(params, &metatype.InstAsst{InstID: currInstID, ObjectID: obj.GetID(), AsstInstID: asstInstID, AsstObjectID: asst.AsstObjID}); nil != err {
					blog.Errorf("[operation-inst] failed to create inst association, error info is %s", err.Error())
					return err
				}
			}

			// update asst value
			inst.SetValue(attr.GetID(), strings.Join(validInstIDS, common.InstAsstIDSplit))
			break
		} // end foreach assts

	}

	return nil
}
func (c *commonInst) CreateInst(params types.ContextParams, obj model.Object, data frtypes.MapStr) (inst.Inst, error) {

	// create new insts
	blog.V(3).Infof("the data inst:%#v", data)
	item := c.instFactory.CreateInst(params, obj)
	item.SetValues(data)

	if err := NewSupplementary().Validator(c).ValidatorCreate(params, obj, item.ToMapStr()); nil != err {
		blog.Errorf("[operation-inst] valid is bad, the data is (%#v)  error info is %s", item.ToMapStr(), err.Error())
		return nil, err
	}

	if err := item.Create(); nil != err {
		blog.Errorf("[operation-inst] failed to save the object(%s) inst data (%#v), error info is %s", obj.GetID(), data, err.Error())
		return nil, err
	}

	NewSupplementary().Audit(params, c.clientSet, item.GetObject(), c).CommitCreateLog(nil, nil, item)

	if err := c.setInstAsst(params, obj, item); nil != err {
		blog.Errorf("[operation-inst] failed to set the inst asst, error info is %s", err.Error())
		return nil, err
	}

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
func (c *commonInst) hasHost(params types.ContextParams, targetInst inst.Inst, checkhost bool) ([]deletedInst, bool, error) {

	id, err := targetInst.GetInstID()
	if nil != err {
		return nil, false, err
	}

	targetObj := targetInst.GetObject()
	if !targetObj.IsCommon() &&
		targetObj.GetObjectType() == common.BKInnerObjIDModule &&
		checkhost {

		exists, err := c.innerHasHost(params, []int64{id})
		if nil != err {
			return nil, false, err
		}

		if exists {

			return nil, true, nil
		}
	}

	instIDS := []deletedInst{}
	instIDS = append(instIDS, deletedInst{instID: id, obj: targetObj})
	childInsts, err := targetInst.GetMainlineChildInst()
	if nil != err {
		return nil, false, err
	}

	for _, childInst := range childInsts {

		ids, exists, err := c.hasHost(params, childInst, checkhost)
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

func (c *commonInst) DeleteInstByInstID(params types.ContextParams, obj model.Object, instID []int64, needCheckHost bool) error {

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
		ids, exists, err := c.hasHost(params, inst, needCheckHost)
		if nil != err {
			return params.Err.Error(common.CCErrTopoHasHostCheckFailed)
		}

		if exists {
			return params.Err.Error(common.CCErrTopoHasHostCheckFailed)
		}

		deleteIDS = append(deleteIDS, ids...)
	}

	for _, delInst := range deleteIDS {
		preAudit := NewSupplementary().Audit(params, c.clientSet, obj, c).CreateSnapshot(delInst.instID, condition.CreateCondition().ToMapStr())

		innerCond := condition.CreateCondition()
		innerCond.Field(common.BKAsstObjIDField).Eq(obj.GetID())
		innerCond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
		innerCond.Field(common.BKAsstInstIDField).Eq(delInst.instID)
		err := c.asst.CheckBeAssociation(params, obj, innerCond)
		if nil != err {
			return err
		}

		if err := c.asst.DeleteInstAssociation(params, innerCond); nil != err {
			blog.Errorf("[operation-inst] failed to set the inst asst, error info is %s", err.Error())
			return err
		}

		innerCond = condition.CreateCondition()
		innerCond.Field(common.BKObjIDField).Eq(obj.GetID())
		innerCond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
		innerCond.Field(common.BKInstIDField).Eq(delInst)
		if err := c.asst.DeleteInstAssociation(params, innerCond); nil != err {
			blog.Errorf("[operation-inst] failed to set the inst asst, error info is %s", err.Error())
			return err
		}

		// clear association
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

	}
	return nil
}

func (c *commonInst) DeleteInst(params types.ContextParams, obj model.Object, cond condition.Condition, needCheckHost bool) error {

	// clear inst associations
	query := &metatype.QueryInput{}
	query.Limit = common.BKNoLimit
	query.Condition = cond.ToMapStr()

	_, insts, err := c.FindInst(params, obj, query, false)
	instIDs := []int64{}
	for _, inst := range insts {
		instID, _ := inst.GetInstID()
		instIDs = append(instIDs, instID)
	}
	blog.V(4).Infof("[DeleteInst] find inst by %+v, returns %+v", query, instIDs)
	if nil != err {
		blog.Errorf("[operation-inst] failed to search insts by the condition(%#v), error info is %s", cond.ToMapStr(), err.Error())
		return err
	}
	for _, inst := range insts {
		targetInstID, err := inst.GetInstID()
		if nil != err {
			return err
		}
		err = c.DeleteInstByInstID(params, obj, []int64{targetInstID}, needCheckHost)
		if nil != err {
			return err
		}

	}

	return nil
}
func (c *commonInst) convertInstIDIntoStruct(params types.ContextParams, asstObj metatype.Association, instIDS []string, needAsstDetail bool) ([]metatype.InstNameAsst, error) {

	obj, err := c.obj.FindSingleObject(params, asstObj.AsstObjID)
	if nil != err {
		return nil, err
	}

	ids := []int64{}
	for _, id := range instIDS {
		if 0 == len(strings.TrimSpace(id)) {
			continue
		}
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

	instAsstNames := []metatype.InstNameAsst{}
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
			instAsstNames = append(instAsstNames, metatype.InstNameAsst{
				ObjID:      obj.GetID(),
				ObjectName: obj.GetName(),
				ObjIcon:    obj.GetIcon(),
				InstID:     instID,
				InstName:   instName,
				InstInfo:   instInfo,
			})
			continue
		}

		instAsstNames = append(instAsstNames, metatype.InstNameAsst{
			ObjID:      obj.GetID(),
			ObjectName: obj.GetName(),
			ObjIcon:    obj.GetIcon(),
			InstID:     instID,
			InstName:   instName,
		})

	}

	return instAsstNames, nil
}

func (c *commonInst) searchAssociationInst(params types.ContextParams, objID string, query *metatype.QueryInput) ([]int64, error) {

	obj, err := c.obj.FindSingleObject(params, objID)
	if nil != err {
		return nil, err
	}

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
	//fmt.Println("cond:", obj.GetID(), query, len(insts))
	tmpResults := map[string]*commonInstTopo{}
	for _, inst := range insts {

		childs, err := inst.GetChildObjectWithInsts()
		if nil != err {
			return 0, nil, err
		}

		for _, child := range childs {

			commonInst, exists := tmpResults[child.Object.GetID()]
			if !exists {
				commonInst = &commonInstTopo{}
				commonInst.ObjectName = child.Object.GetName()
				commonInst.ObjIcon = child.Object.GetIcon()
				commonInst.ObjID = child.Object.GetID()
				commonInst.Children = []metatype.InstNameAsst{}
				tmpResults[child.Object.GetID()] = commonInst
			}

			commonInst.Count = commonInst.Count + len(child.Insts)

			for _, childInst := range child.Insts {

				instAsst := metatype.InstNameAsst{}
				id, err := childInst.GetInstID()
				if nil != err {
					return 0, nil, err
				}

				name, err := childInst.GetInstName()
				if nil != err {
					return 0, nil, err
				}

				instAsst.ID = strconv.Itoa(int(id))
				instAsst.InstID = id
				instAsst.InstName = name
				instAsst.ObjectName = child.Object.GetName()
				instAsst.ObjIcon = child.Object.GetIcon()
				instAsst.ObjID = child.Object.GetID()

				tmpResults[child.Object.GetID()].Children = append(tmpResults[child.Object.GetID()].Children, instAsst)
			}
		}
	}

	for _, subResult := range tmpResults {
		results = append(results, subResult)
	}

	return len(results), results, nil
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

		parents, err := inst.GetParentObjectWithInsts()
		if nil != err {
			return 0, nil, err
		}

		for _, parent := range parents {

			commonInst, exists := tmpResults[parent.Object.GetID()]
			if !exists {
				commonInst = &commonInstTopo{}
				commonInst.ObjectName = parent.Object.GetName()
				commonInst.ObjIcon = parent.Object.GetIcon()
				commonInst.ObjID = parent.Object.GetID()
				commonInst.Children = []metatype.InstNameAsst{}
				tmpResults[parent.Object.GetID()] = commonInst
			}

			commonInst.Count = commonInst.Count + len(parent.Insts)

			for _, parentInst := range parent.Insts {
				instAsst := metatype.InstNameAsst{}
				id, err := parentInst.GetInstID()
				if nil != err {
					return 0, nil, err
				}

				name, err := parentInst.GetInstName()
				if nil != err {
					return 0, nil, err
				}
				instAsst.ID = strconv.Itoa(int(id))
				instAsst.InstID = id
				instAsst.InstName = name
				instAsst.ObjectName = parent.Object.GetName()
				instAsst.ObjIcon = parent.Object.GetIcon()
				instAsst.ObjID = parent.Object.GetID()

				tmpResults[parent.Object.GetID()].Children = append(tmpResults[parent.Object.GetID()].Children, instAsst)
			}
		}
	}

	for _, subResult := range tmpResults {
		results = append(results, subResult)
	}

	return len(results), results, nil
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

		commonInst := commonInstTopo{Children: []metatype.InstNameAsst{}}
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

	return len(results), results, nil
}

func (c *commonInst) FindInstByAssociationInst(params types.ContextParams, obj model.Object, data frtypes.MapStr) (cont int, results []inst.Inst, err error) {

	asstParamCond := &AssociationParams{}

	if err := data.MarshalJSONInto(asstParamCond); nil != err {
		blog.Errorf("[operation-inst] find inst by association inst , error info is %s", err.Error())
		return 0, nil, params.Err.Errorf(common.CCErrTopoInstSelectFailed, err.Error())
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

		innerCond := new(metatype.QueryInput)
		if fields, ok := asstParamCond.Fields[keyObjID]; ok {
			innerCond.Fields = strings.Join(fields, ",")
		}
		innerCond.Condition = cond

		asstInstIDS, err := c.searchAssociationInst(params, keyObjID, innerCond)
		if nil != err {
			blog.Errorf("[operation-inst]failed to search the association inst, error info is %s", err.Error())
			return 0, nil, err
		}
		blog.V(4).Infof("[FindInstByAssociationInst] search association insts, keyObjID %s, condition: %v, results: %v", keyObjID, innerCond, asstInstIDS)

		query := &metatype.QueryInput{}
		query.Condition = map[string]interface{}{
			"bk_asst_inst_id": map[string]interface{}{
				common.BKDBIN: asstInstIDS,
			},
			"bk_asst_obj_id": keyObjID,
			"bk_obj_id":      obj.GetID(),
		}

		asstInst, err := c.asst.SearchInstAssociation(params, query)
		if nil != err {
			blog.Errorf("[operation-inst] failed to search the association inst, error info is %s", err.Error())
			return 0, nil, err
		}

		for _, asst := range asstInst {
			targetInstIDS = append(targetInstIDS, asst.InstID)
		}
		blog.V(4).Infof("[FindInstByAssociationInst] search association, objectID=%s, keyObjID=%s, condition: %v, results: %v", obj.GetID(), keyObjID, query, targetInstIDS)

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

	query := &metatype.QueryInput{}
	query.Condition = instCond
	if fields, ok := asstParamCond.Fields[obj.GetID()]; ok {
		query.Fields = strings.Join(fields, ",")
	}
	query.Limit = asstParamCond.Page.Limit
	query.Sort = asstParamCond.Page.Sort
	query.Start = asstParamCond.Page.Start
	blog.V(4).Infof("[FindInstByAssociationInst] search inst condition: %v", instCond)
	return c.FindInst(params, obj, query, false)
}
func (c *commonInst) FindOriginInst(params types.ContextParams, obj model.Object, cond *metatype.QueryInput) (*metatype.InstResult, error) {

	switch obj.GetID() {
	case common.BKInnerObjIDHost:
		rsp, err := c.clientSet.HostController().Host().GetHosts(context.Background(), params.Header, cond)
		if nil != err {
			blog.Errorf("[operation-inst] failed to request object controller, error info is %s", err.Error())
			return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if common.CCSuccess != rsp.Code {

			blog.Errorf("[operation-inst] faild to delete the object(%s) inst by the condition(%#v), error info is %s", obj.GetID(), cond, rsp.ErrMsg)
			return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
		}

		return &metatype.InstResult{Count: rsp.Data.Count, Info: frtypes.NewArrayFromInterface(rsp.Data.Info)}, nil

	default:

		rsp, err := c.clientSet.ObjectController().Instance().SearchObjects(context.Background(), obj.GetObjectType(), params.Header, cond)

		if nil != err {
			blog.Errorf("[operation-inst] failed to request object controller, error info is %s", err.Error())
			return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if common.CCSuccess != rsp.Code {

			blog.Errorf("[operation-inst] faild to delete the object(%s) inst by the condition(%#v), error info is %s", obj.GetID(), cond, rsp.ErrMsg)
			return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
		}
		return &rsp.Data, nil
	}

}
func (c *commonInst) FindInst(params types.ContextParams, obj model.Object, cond *metatype.QueryInput, needAsstDetail bool) (count int, results []inst.Inst, err error) {

	rsp, err := c.FindOriginInst(params, obj, cond)
	if nil != err {
		blog.Errorf("[operation-inst] failed to find origin inst , error info is %s", err.Error())
		return 0, nil, err
	}

	asstObjAttrs, err := c.asst.SearchObjectAssociation(params, obj.GetID())
	if nil != err {
		blog.Errorf("[operation-inst] failed to search object associations, error info is %s", err.Error())
		return 0, nil, err
	}

	for idx, instInfo := range rsp.Info {

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
			rsp.Info[idx].Set(attrAsst.ObjectAttID, instVals)

		}
	}
	return rsp.Count, inst.CreateInst(params, c.clientSet, obj, rsp.Info), nil
}

func (c *commonInst) UpdateInst(params types.ContextParams, data frtypes.MapStr, obj model.Object, cond condition.Condition, instID int64) error {

	if err := NewSupplementary().Validator(c).ValidatorUpdate(params, obj, data, instID, cond); nil != err {
		return err
	}

	// update association
	query := &metatype.QueryInput{}
	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit
	if 0 < instID {
		innerCond := condition.CreateCondition()
		innerCond.Field(obj.GetInstIDFieldName()).Eq(instID)
		query.Condition = innerCond.ToMapStr()
	}
	_, insts, err := c.FindInst(params, obj, query, false)
	if nil != err {
		blog.Errorf("[operation-inst] failed to search insts by the condition(%#v), error info is %s", cond.ToMapStr(), err.Error())
		return err
	}
	for _, inst := range insts {

		data.ForEach(func(key string, val interface{}) {
			inst.SetValue(key, val)
		})

		if err := c.setInstAsst(params, obj, inst); nil != err {
			blog.Errorf("[operation-inst] failed to set the inst asst, error info is %s", err.Error())
			return err
		}
	}

	// update insts
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
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}
	currAuditLog := NewSupplementary().Audit(params, c.clientSet, obj, c).CreateSnapshot(-1, cond.ToMapStr())
	NewSupplementary().Audit(params, c.clientSet, obj, c).CommitUpdateLog(preAuditLog, currAuditLog, nil)
	return nil
}
