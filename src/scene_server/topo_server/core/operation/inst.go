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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	gparams "configcenter/src/common/paraparse"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// InstOperationInterface inst operation methods
type InstOperationInterface interface {
	CreateInst(params types.ContextParams, obj model.Object, data mapstr.MapStr) (inst.Inst, error)
	CreateInstBatch(params types.ContextParams, obj model.Object, batchInfo *InstBatchInfo) (*BatchResult, error)
	DeleteInst(params types.ContextParams, obj model.Object, cond condition.Condition, needCheckHost bool) error
	DeleteMainlineInstWithID(params types.ContextParams, obj model.Object, instID int64) error
	DeleteInstByInstID(params types.ContextParams, obj model.Object, instID []int64, needCheckHost bool) error
	FindOriginInst(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (*metadata.InstResult, error)
	FindInst(params types.ContextParams, obj model.Object, cond *metadata.QueryInput, needAsstDetail bool) (count int, results []inst.Inst, err error)
	FindInstByAssociationInst(params types.ContextParams, obj model.Object, data mapstr.MapStr) (cont int, results []inst.Inst, err error)
	FindInstChildTopo(params types.ContextParams, obj model.Object, instID int64, query *metadata.QueryInput) (count int, results []*CommonInstTopo, err error)
	FindInstParentTopo(params types.ContextParams, obj model.Object, instID int64, query *metadata.QueryInput) (count int, results []*CommonInstTopo, err error)
	FindInstTopo(params types.ContextParams, obj model.Object, instID int64, query *metadata.QueryInput) (count int, results []CommonInstTopoV2, err error)
	UpdateInst(params types.ContextParams, data mapstr.MapStr, obj model.Object, cond condition.Condition, instID int64) error

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

type BatchResult struct {
	Errors         []string `json:"error"`
	Success        []string `json:"success"`
	SuccessCreated []int64  `json:"success_created"`
	SuccessUpdated []int64  `json:"success_updated"`
	UpdateErrors   []string `json:"update_error"`
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

// CreateInstBatch
func (c *commonInst) CreateInstBatch(params types.ContextParams, obj model.Object, batchInfo *InstBatchInfo) (*BatchResult, error) {
	var err error
	var bizID int64
	if params.MetaData != nil {
		bizID, err = metadata.BizIDFromMetadata(*params.MetaData)
		if err != nil {
			return nil, fmt.Errorf("parse business id from metadata failed, err: %+v", err)
		}
	}

	object := obj.Object()

	// forbidden create inner model instance with common api
	if common.IsInnerModel(object.ObjectID) == true {
		blog.V(5).Infof("CreateInstBatch failed, create %s instance with common create api forbidden, rid: %s", object.ObjectID, params.ReqID)
		return nil, params.Err.Error(common.CCErrTopoImportMainlineForbidden)
	}

	isMainlin, err := obj.IsMainlineObject()
	if err != nil {
		blog.Errorf("[operation-inst] failed to get if the object(%s) is mainline object, err: %s, rid: %s", object.ObjectID, err.Error(), params.ReqID)
		return nil, err
	}
	if isMainlin {
		blog.V(5).Infof("CreateInstBatch failed, create %s instance with common create api forbidden, rid: %s", object.ObjectID, params.ReqID)
		return nil, params.Err.Error(common.CCErrTopoImportMainlineForbidden)

	}

	results := &BatchResult{}
	if batchInfo.InputType != common.InputTypeExcel {
		return results, fmt.Errorf("unexpected input_type: %s", batchInfo.InputType)
	}
	if len(batchInfo.BatchInfo) == 0 {
		return results, fmt.Errorf("BatchInfo empty")
	}

	// 1. 检查实例与URL参数指定的模型一致
	for line, inst := range batchInfo.BatchInfo {
		objID, exist := inst[common.BKObjIDField]
		if exist == true && objID != object.ObjectID {
			blog.Errorf("create object[%s] instance batch failed, because bk_obj_id field conflict with url field, rid: %s", object.ObjectID, params.ReqID)
			return nil, params.Err.Errorf(common.CCErrorTopoObjectInstanceObjIDFieldConflictWithURL, line)
		}
	}

	// 2. 检查批量数据中实例名称是否重复
	instNameMap := make(map[string]bool)
	for line, inst := range batchInfo.BatchInfo {
		iName, exist := inst[common.BKInstNameField]
		if !exist {
			blog.Errorf("create object[%s] instance batch failed, because missing bk_inst_name field., rid: %s", object.ObjectID, params.ReqID)
			return nil, params.Err.Errorf(common.CCErrorTopoObjectInstanceMissingInstanceNameField, line)
		}

		name, can := iName.(string)
		if !can {
			blog.Errorf("create object[%s] instance batch failed, because  bk_inst_name value type is not string., rid: %s", object.ObjectID, params.ReqID)
			return nil, params.Err.Errorf(common.CCErrorTopoInvalidObjectInstanceNameFieldValue, line)
		}

		// check if this instance name is already exist.
		if _, ok := instNameMap[name]; ok {
			blog.Errorf("create object[%s] instance batch, but bk_inst_name %s is duplicated., rid: %s", object.ObjectID, name, params.ReqID)
			return nil, params.Err.Errorf(common.CCErrorTopoMultipleObjectInstanceName, name)
		}

		instNameMap[name] = true
	}

	nonInnerAttributes, err := obj.GetNonInnerAttributes()
	if err != nil {
		blog.Errorf("[audit]failed to get the object(%s)' attribute, err: %s, rid: %s", obj.Object().ObjectID, err.Error(), params.ReqID)
		return nil, err
	}

	updatedInstanceIDs := make([]int64, 0)
	createdInstanceIDs := make([]int64, 0)
	idFieldname := metadata.GetInstIDFieldByObjID(obj.GetObjectID())
	for colIdx, colInput := range batchInfo.BatchInfo {
		if colInput == nil {
			// ignore empty excel line
			continue
		}

		delete(colInput, "import_from")
		// create memory object
		item := c.instFactory.CreateInst(params, obj)

		item.SetValues(colInput)

		// 实例id 为空，表示要新建实例
		// 实例ID已经赋值，更新数据.  (已经赋值, value not equal 0 or nil)

		// 是否存在实例ID字段
		instID, existInstID := colInput[idFieldname]
		// 实例ID字段是否设置值
		if existInstID && (instID == "" || instID == nil) {
			existInstID = false
		}
		if existInstID {
			delete(colInput, idFieldname)
			filter := condition.CreateCondition()
			filter = filter.Field(idFieldname).Eq(instID)

			preAuditLog := NewSupplementary().Audit(params, c.clientSet, obj, c).CreateSnapshot(-1, filter.ToMapStr())
			err = item.UpdateInstance(filter, colInput, nonInnerAttributes)
			if nil != err {
				blog.Errorf("[operation-inst] failed to update the object(%s) inst data (%#v), err: %s, rid: %s", object.ObjectID, colInput, err.Error(), params.ReqID)
				results.Errors = append(results.Errors, params.Lang.Languagef("import_row_int_error_str", colIdx, err.Error()))
				continue
			}
			instID, err := item.GetInstID()
			if err != nil {
				blog.ErrorJSON("update inst success, but get id field failed, inst: %s, err: %s, rid: %s", item.GetValues(), err.Error(), params.ReqID)
				results.Errors = append(results.Errors, params.Lang.Languagef("import_row_int_error_str", colIdx, err.Error()))
				continue
			}
			updatedInstanceIDs = append(updatedInstanceIDs, instID)
			results.Success = append(results.Success, strconv.FormatInt(colIdx, 10))
			currAuditLog := NewSupplementary().Audit(params, c.clientSet, obj, c).CreateSnapshot(-1, filter.ToMapStr())
			NewSupplementary().Audit(params, c.clientSet, item.GetObject(), c).CommitUpdateLog(preAuditLog, currAuditLog, nil, nonInnerAttributes)
			continue
		}

		// create with metadata
		if bizID != 0 {
			colInput[metadata.BKMetadata] = metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))
		}
		// set data
		// call CoreService.CreateInstance
		err = item.Create()
		if nil != err {
			blog.Errorf("[operation-inst] failed to save the object(%s) inst data (%#v), err: %s, rid: %s", object.ObjectID, colInput, err.Error(), params.ReqID)
			results.Errors = append(results.Errors, params.Lang.Languagef("import_row_int_error_str", colIdx, err.Error()))
			continue
		}
		results.Success = append(results.Success, strconv.FormatInt(colIdx, 10))
		NewSupplementary().Audit(params, c.clientSet, item.GetObject(), c).CommitCreateLog(nil, nil, item, nonInnerAttributes)

		instanceID, err := item.GetInstID()
		if err != nil {
			blog.Errorf("unexpected error, instances created success, but get id failed, err: %+v, rid: %s", err, params.ReqID)
			continue
		}
		createdInstanceIDs = append(createdInstanceIDs, instanceID)
	}

	results.SuccessCreated = createdInstanceIDs
	results.SuccessUpdated = updatedInstanceIDs

	return results, nil
}

func (c *commonInst) isValidBizInstID(params types.ContextParams, obj metadata.Object, instID int64, bizID int64) error {

	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)
	if bizID != 0 {
		or := cond.NewOR()
		or.Item(mapstr.MapStr{common.BKAppIDField: bizID})
		meta := metadata.Metadata{
			Label: map[string]string{
				common.BKAppIDField: strconv.FormatInt(bizID, 10),
			},
		}
		or.Item(mapstr.MapStr{metadata.BKMetadata: meta})
	}
	if obj.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(obj.ObjectID)
	}

	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit

	rsp, err := c.clientSet.CoreService().Instance().ReadInstance(context.Background(), params.Header, obj.GetObjectID(), &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-inst] faild to read the object(%s) inst by the condition(%#v), err: %s, rid: %s", obj.ObjectID, cond, rsp.ErrMsg, params.ReqID)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	if rsp.Data.Count > 0 {
		return nil
	}

	return params.Err.Error(common.CCErrTopoInstSelectFailed)
}

func (c *commonInst) isValidInstID(params types.ContextParams, obj metadata.Object, instID int64) error {
	return c.isValidBizInstID(params, obj, instID, 0)
}

func (c *commonInst) validMainLineParentID(params types.ContextParams, obj model.Object, data mapstr.MapStr) error {
	if obj.Object().ObjectID == common.BKInnerObjIDApp {
		return nil
	}
	def, exist := data.Get(common.BKDefaultField)
	if exist && def.(int) != common.DefaultFlagDefaultValue {
		return nil
	}
	parent, err := obj.GetMainlineParentObject()
	if err != nil {
		blog.Errorf("[operation-inst] failed to get the object(%s) mainline parent, err: %s, rid: %s", obj.Object().ObjectID, err.Error(), params.ReqID)
		return err
	}
	bizID, err := data.Int64(common.BKAppIDField)
	if err != nil {
		bizID, err = metadata.ParseBizIDFromData(data)
		if err != nil {
			blog.Errorf("[operation-inst]failed to parse the biz id, err: %s, rid: %s", err.Error(), params.ReqID)
			return params.Err.Errorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField)
		}
	}
	parentID, err := data.Int64(common.BKParentIDField)
	if err != nil {
		blog.Errorf("[operation-inst]failed to parse the parent id, err: %s, rid: %s", err.Error(), params.ReqID)
		return params.Err.Errorf(common.CCErrCommParamsIsInvalid, common.BKParentIDField)
	}
	if err = c.isValidBizInstID(params, parent.Object(), parentID, bizID); err != nil {
		blog.Errorf("[operation-inst]parent id %d is invalid, err: %s, rid: %s", parentID, err.Error(), params.ReqID)
		return params.Err.Errorf(common.CCErrCommParamsIsInvalid, common.BKParentIDField)
	}
	return nil
}

func (c *commonInst) CreateInst(params types.ContextParams, obj model.Object, data mapstr.MapStr) (inst.Inst, error) {

	// create new insts
	item := c.instFactory.CreateInst(params, obj)
	item.SetValues(data)

	iData := item.ToMapStr()
	if obj.Object().ObjectID == common.BKInnerObjIDPlat {
		iData["bk_supplier_account"] = params.SupplierAccount
	}

	isMainline, err := obj.IsMainlineObject()
	if err != nil {
		blog.Errorf("[operation-inst] failed to get if the object(%s) is mainline object, err: %s, rid: %s", obj.Object().ObjectID, err.Error(), params.ReqID)
		return nil, err
	}
	if isMainline {
		if err := c.validMainLineParentID(params, obj, data); nil != err {
			blog.Errorf("[operation-inst] the mainline object(%s) parent id invalid, err: %s, rid: %s", obj.Object().ObjectID, err.Error(), params.ReqID)
			return nil, err
		}
	}

	if err := item.Create(); nil != err {
		blog.Errorf("[operation-inst] failed to save the object(%s) inst data (%#v), err: %s, rid: %s", obj.Object().ObjectID, data, err.Error(), params.ReqID)
		return nil, err
	}

	NewSupplementary().Audit(params, c.clientSet, item.GetObject(), c).CommitCreateLog(nil, nil, item, nil)

	instID, err := item.GetInstID()
	if err != nil {
		return nil, params.Err.Error(common.CCErrTopoInstCreateFailed)
	}
	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)
	_, insts, err := c.FindInst(params, obj, &metadata.QueryInput{Condition: cond.ToMapStr()}, false)
	if err != nil {
		return nil, err
	}

	for _, inst := range insts {
		return inst, nil
	}

	return item, nil
}

func (c *commonInst) innerHasHost(params types.ContextParams, moduleIDS []int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ModuleIDArr: moduleIDS,
	}
	rsp, err := c.clientSet.CoreService().Host().GetHostModuleRelation(context.Background(), params.Header, option)
	if nil != err {
		blog.Errorf("[operation-module] failed to request the object controller, err: %s, rid: %s", err.Error(), params.ReqID)
		return false, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-module]  failed to search the host module configures, err: %s, rid: %s", rsp.ErrMsg, params.ReqID)
		return false, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != len(rsp.Data.Info), nil
}
func (c *commonInst) hasHost(params types.ContextParams, targetInst inst.Inst, checkhost bool) ([]deletedInst, bool, error) {

	id, err := targetInst.GetInstID()
	if nil != err {
		return nil, false, err
	}

	targetObj := targetInst.GetObject()
	// if this is a module object and need to check host, then check.
	if !targetObj.IsCommon() && targetObj.GetObjectType() == common.BKInnerObjIDModule && checkhost {
		exists, err := c.innerHasHost(params, []int64{id})
		if nil != err {
			return nil, false, err
		}

		if exists {
			return nil, true, nil
		}
	}

	instIDS := make([]deletedInst, 0)
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
	object := obj.Object()
	objID := object.ID
	objectID := object.ObjectID

	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).In(instID)
	if obj.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(objectID)
	}

	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()

	_, insts, err := c.FindInst(params, obj, query, false)
	if nil != err {
		return err
	}

	deleteIDS := make([]deletedInst, 0)
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
		auditFilter := condition.CreateCondition().ToMapStr()
		preAudit := NewSupplementary().Audit(params, c.clientSet, delInst.obj, c).CreateSnapshot(delInst.instID, auditFilter)

		// if this instance has been bind to a instance by the association, then this instance should not be deleted.
		innerCond := condition.CreateCondition()
		innerCond.Field(common.BKAsstObjIDField).Eq(objID)
		innerCond.Field(common.BKAsstInstIDField).Eq(delInst.instID)
		err := c.asst.CheckBeAssociation(params, obj, innerCond)
		if nil != err {
			return err
		}

		// this instance has not be bind to another instance, we can delete all the associations it created
		// by the association with other instances.
		innerCond = condition.CreateCondition()
		innerCond.Field(common.BKObjIDField).Eq(objID)
		innerCond.Field(common.BKInstIDField).Eq(delInst.instID)
		if err := c.asst.DeleteInstAssociation(params, innerCond); nil != err {
			blog.Errorf("[operation-inst] failed to delete the inst asst, err: %s, rid: %s", err.Error(), params.ReqID)
			return err
		}

		// delete this instance now.
		delCond := condition.CreateCondition()
		delCond.Field(delInst.obj.GetInstIDFieldName()).In(delInst.instID)
		if delInst.obj.IsCommon() {
			delCond.Field(common.BKObjIDField).Eq(objID)
		}
		// clear association
		dc := &metadata.DeleteOption{Condition: delCond.ToMapStr()}
		instObjID := delInst.obj.GetObjectID()
		rsp, err := c.clientSet.CoreService().Instance().DeleteInstance(params.Context, params.Header, instObjID, dc)
		if nil != err {
			blog.Errorf("[operation-inst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
			return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("[operation-inst] failed to delete the object(%s) inst by the condition(%#v), err: %s, rid: %s", objectID, delCond.ToMapStr(), rsp.ErrMsg, params.ReqID)
			return params.Err.New(rsp.Code, rsp.ErrMsg)
		}

		NewSupplementary().Audit(params, c.clientSet, delInst.obj, c).CommitDeleteLog(preAudit, nil, nil)
	}
	return nil
}

func (c *commonInst) DeleteMainlineInstWithID(params types.ContextParams, obj model.Object, instID int64) error {
	object := obj.Object()
	preAudit := NewSupplementary().Audit(params, c.clientSet, obj, c).CreateSnapshot(instID, condition.CreateCondition().ToMapStr())
	// if this instance has been bind to a instance by the association, then this instance should not be deleted.
	innerCond := condition.CreateCondition()
	innerCond.Field(common.BKAsstObjIDField).Eq(object.ObjectID)
	innerCond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	innerCond.Field(common.BKAsstInstIDField).Eq(instID)
	err := c.asst.CheckBeAssociation(params, obj, innerCond)
	if nil != err {
		return err
	}

	// this instance has not be bind to another instance, we can delete all the associations it created
	// by the association with other instances.
	innerCond = condition.CreateCondition()
	innerCond.Field(common.BKObjIDField).Eq(object.ObjectID)
	innerCond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	innerCond.Field(common.BKInstIDField).Eq(instID)
	if err = c.asst.DeleteInstAssociation(params, innerCond); nil != err {
		blog.Errorf("[operation-inst] failed to delete the inst asst, err: %s", err.Error())
		return err
	}

	// delete this instance now.
	delCond := condition.CreateCondition()
	delCond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	delCond.Field(obj.GetInstIDFieldName()).Eq(instID)
	if obj.IsCommon() {
		delCond.Field(common.BKObjIDField).Eq(object.ObjectID)
	}

	ops := metadata.DeleteOption{
		Condition: delCond.ToMapStr(),
	}
	rsp, err := c.clientSet.CoreService().Instance().DeleteInstance(params.Context, params.Header, object.ObjectID, &ops)
	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, err: %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-inst] failed to delete the object(%s) inst by the condition(%#v), err: %s", object.ObjectID, delCond.ToMapStr(), rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}

	NewSupplementary().Audit(params, c.clientSet, obj, c).CommitDeleteLog(preAudit, nil, nil)

	return nil
}

func (c *commonInst) DeleteInst(params types.ContextParams, obj model.Object, cond condition.Condition, needCheckHost bool) error {

	// clear inst associations
	query := &metadata.QueryInput{}
	query.Limit = common.BKNoLimit
	query.Condition = cond.ToMapStr()

	_, insts, err := c.FindInst(params, obj, query, false)
	instIDs := make([]int64, 0)
	for _, inst := range insts {
		instID, _ := inst.GetInstID()
		instIDs = append(instIDs, instID)
	}
	blog.V(4).Infof("[DeleteInst] find inst by %+v, returns %+v, rid: %s", query, instIDs, params.ReqID)
	if nil != err {
		blog.Errorf("[operation-inst] failed to search insts by the condition(%#v), err: %s, rid: %s", cond.ToMapStr(), err.Error(), params.ReqID)
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
func (c *commonInst) convertInstIDIntoStruct(params types.ContextParams, asstObj metadata.Association, instIDS []string, needAsstDetail bool) ([]metadata.InstNameAsst, error) {

	obj, err := c.obj.FindSingleObject(params, asstObj.AsstObjID)
	if nil != err {
		return nil, err
	}
	object := obj.Object()

	ids := make([]int64, 0)
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

	query := &metadata.QueryCondition{}
	query.Condition = cond.ToMapStr()
	rsp, err := c.clientSet.CoreService().Instance().ReadInstance(context.Background(), params.Header, obj.GetObjectID(), query)

	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-inst] faild to delete the object(%s) inst by the condition(%#v), err: %s, rid: %s", object.ObjectID, cond, rsp.ErrMsg, params.ReqID)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	instAsstNames := []metadata.InstNameAsst{}
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
			instAsstNames = append(instAsstNames, metadata.InstNameAsst{
				ID:         strconv.Itoa(int(instID)),
				ObjID:      object.ObjectID,
				ObjectName: object.ObjectName,
				ObjIcon:    object.ObjIcon,
				InstID:     instID,
				InstName:   instName,
				InstInfo:   instInfo,
			})
			continue
		}

		instAsstNames = append(instAsstNames, metadata.InstNameAsst{
			ID:         strconv.Itoa(int(instID)),
			ObjID:      object.ObjectID,
			ObjectName: object.ObjectName,
			ObjIcon:    object.ObjIcon,
			InstID:     instID,
			InstName:   instName,
		})

	}

	return instAsstNames, nil
}

func (c *commonInst) searchAssociationInst(params types.ContextParams, objID string, query *metadata.QueryInput) ([]int64, error) {

	obj, err := c.obj.FindSingleObject(params, objID)
	if nil != err {
		return nil, err
	}

	_, insts, err := c.FindInst(params, obj, query, false)
	if nil != err {
		return nil, err
	}

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

func (c *commonInst) FindInstChildTopo(params types.ContextParams, obj model.Object, instID int64, query *metadata.QueryInput) (count int, results []*CommonInstTopo, err error) {
	results = make([]*CommonInstTopo, 0)
	if nil == query {
		query = &metadata.QueryInput{}
		cond := condition.CreateCondition()
		cond.Field(obj.GetInstIDFieldName()).Eq(instID)
		query.Condition = cond.ToMapStr()
	}

	_, insts, err := c.FindInst(params, obj, query, false)
	if nil != err {
		return 0, nil, err
	}

	tmpResults := map[string]*CommonInstTopo{}
	for _, inst := range insts {

		children, err := inst.GetChildObjectWithInsts()
		if nil != err {
			return 0, nil, err
		}

		for _, child := range children {
			object := child.Object.Object()
			commonInst, exists := tmpResults[object.ObjectID]
			if !exists {
				commonInst = &CommonInstTopo{}
				commonInst.ObjectName = object.ObjectName
				commonInst.ObjIcon = object.ObjIcon
				commonInst.ObjID = object.ObjectID
				commonInst.Children = []metadata.InstNameAsst{}
				tmpResults[object.ObjectID] = commonInst
			}

			commonInst.Count = commonInst.Count + len(child.Insts)

			for _, childInst := range child.Insts {

				instAsst := metadata.InstNameAsst{}
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
				instAsst.ObjectName = object.ObjectName
				instAsst.ObjIcon = object.ObjIcon
				instAsst.ObjID = object.ObjectID
				instAsst.AssoID = childInst.GetAssoID()

				tmpResults[object.ObjectID].Children = append(tmpResults[object.ObjectID].Children, instAsst)
			}
		}
	}

	for _, subResult := range tmpResults {
		results = append(results, subResult)
	}

	return len(results), results, nil
}

func (c *commonInst) FindInstParentTopo(params types.ContextParams, obj model.Object, instID int64, query *metadata.QueryInput) (count int, results []*CommonInstTopo, err error) {

	results = make([]*CommonInstTopo, 0)
	if nil == query {
		query = &metadata.QueryInput{}
		cond := condition.CreateCondition()
		cond.Field(obj.GetInstIDFieldName()).Eq(instID)
		query.Condition = cond.ToMapStr()
	}

	_, insts, err := c.FindInst(params, obj, query, false)
	if nil != err {
		return 0, nil, err
	}

	tmpResults := map[string]*CommonInstTopo{}
	for _, inst := range insts {

		parents, err := inst.GetParentObjectWithInsts()
		if nil != err {
			return 0, nil, err
		}

		for _, parent := range parents {
			object := parent.Object.Object()
			commonInst, exists := tmpResults[object.ObjectID]
			if !exists {
				commonInst = &CommonInstTopo{}
				commonInst.ObjectName = object.ObjectName
				commonInst.ObjIcon = object.ObjIcon
				commonInst.ObjID = object.ObjectID
				commonInst.Children = []metadata.InstNameAsst{}
				tmpResults[object.ObjectID] = commonInst
			}

			commonInst.Count = commonInst.Count + len(parent.Insts)

			for _, parentInst := range parent.Insts {
				instAsst := metadata.InstNameAsst{}
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
				instAsst.ObjectName = object.ObjectName
				instAsst.ObjIcon = object.ObjIcon
				instAsst.ObjID = object.ObjectID
				instAsst.AssoID = parentInst.GetAssoID()

				tmpResults[object.ObjectID].Children = append(tmpResults[object.ObjectID].Children, instAsst)
			}
		}
	}

	for _, subResult := range tmpResults {
		results = append(results, subResult)
	}

	return len(results), results, nil
}

func (c *commonInst) FindInstTopo(params types.ContextParams, obj model.Object, instID int64, query *metadata.QueryInput) (count int, results []CommonInstTopoV2, err error) {

	if nil == query {
		query = &metadata.QueryInput{}
		cond := condition.CreateCondition()
		cond.Field(obj.GetInstIDFieldName()).Eq(instID)
		query.Condition = cond.ToMapStr()
	}

	_, insts, err := c.FindInst(params, obj, query, false)
	if nil != err {
		blog.Errorf("[operation-inst] failed to find the inst, err: %s, rid: %s", err.Error(), params.ReqID)
		return 0, nil, err
	}

	for _, inst := range insts {
		id, err := inst.GetInstID()
		if nil != err {
			blog.Errorf("[operation-inst] failed to find the inst, err: %s, rid: %s", err.Error(), params.ReqID)
			return 0, nil, err
		}

		name, err := inst.GetInstName()
		if nil != err {
			blog.Errorf("[operation-inst] failed to find the inst, err: %s, rid: %s", err.Error(), params.ReqID)
			return 0, nil, err
		}

		object := inst.GetObject().Object()

		commonInst := metadata.InstNameAsst{}
		commonInst.ObjectName = object.ObjectName
		commonInst.ObjID = object.ObjectID
		commonInst.ObjIcon = object.ObjIcon
		commonInst.InstID = id
		commonInst.ID = strconv.Itoa(int(id))
		commonInst.InstName = name

		_, parentInsts, err := c.FindInstParentTopo(params, inst.GetObject(), id, nil)
		if nil != err {
			blog.Errorf("[operation-inst] failed to find the inst, err: %s, rid: %s", err.Error(), params.ReqID)
			return 0, nil, err
		}

		_, childInsts, err := c.FindInstChildTopo(params, inst.GetObject(), id, nil)
		if nil != err {
			blog.Errorf("[operation-inst] failed to find the inst, err: %s, rid: %s", err.Error(), params.ReqID)
			return 0, nil, err
		}

		results = append(results, CommonInstTopoV2{
			Prev: parentInsts,
			Next: childInsts,
			Curr: commonInst,
		})

	}

	return len(results), results, nil
}

func (c *commonInst) FindInstByAssociationInst(params types.ContextParams, obj model.Object, data mapstr.MapStr) (cont int, results []inst.Inst, err error) {

	asstParamCond := &AssociationParams{}
	if err := data.MarshalJSONInto(asstParamCond); nil != err {
		blog.Errorf("[operation-inst] find inst by association inst , err: %s, rid: %s", err.Error(), params.ReqID)
		return 0, nil, params.Err.Errorf(common.CCErrTopoInstSelectFailed, err.Error())
	}

	object := obj.Object()

	instCond := map[string]interface{}{}
	if obj.IsCommon() {
		instCond[common.BKObjIDField] = object.ObjectID
	}
	targetInstIDS := make([]int64, 0)

	for keyObjID, objs := range asstParamCond.Condition {
		// Extract the ID of the instance according to the associated object.
		cond := map[string]interface{}{}
		if common.GetObjByType(keyObjID) == common.BKInnerObjIDObject {
			cond[common.BKObjIDField] = keyObjID
		}

		for _, objCondition := range objs {
			if objCondition.Operator != common.BKDBEQ {
				if object.ObjectID == keyObjID {
					if objCondition.Operator == common.BKDBLIKE ||
						objCondition.Operator == common.BKDBMULTIPLELike {
						switch t := objCondition.Value.(type) {
						case string:
							instCond[objCondition.Field] = map[string]interface{}{
								objCondition.Operator: gparams.SpecialCharChange(t),
							}
						default:
							// deal self condition
							instCond[objCondition.Field] = map[string]interface{}{
								objCondition.Operator: objCondition.Value,
							}
						}
					} else {
						// deal self condition
						instCond[objCondition.Field] = map[string]interface{}{
							objCondition.Operator: objCondition.Value,
						}
					}
				} else {
					// deal association condition
					cond[objCondition.Field] = map[string]interface{}{
						objCondition.Operator: objCondition.Value,
					}
				}
			} else {
				if object.ObjectID == keyObjID {
					// deal self condition
					switch t := objCondition.Value.(type) {
					case string:
						instCond[objCondition.Field] = map[string]interface{}{
							common.BKDBEQ: t,
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

		if object.ObjectID == keyObjID {
			// no need to search the association objects
			continue
		}

		innerCond := new(metadata.QueryInput)
		if fields, ok := asstParamCond.Fields[keyObjID]; ok {
			innerCond.Fields = strings.Join(fields, ",")
		}
		innerCond.Condition = cond

		asstInstIDS, err := c.searchAssociationInst(params, keyObjID, innerCond)
		if nil != err {
			blog.Errorf("[operation-inst]failed to search the association inst, err: %s, rid: %s", err.Error(), params.ReqID)
			return 0, nil, err
		}
		blog.V(4).Infof("[FindInstByAssociationInst] search association insts, keyObjID %s, condition: %v, results: %v, rid: %s", keyObjID, innerCond, asstInstIDS, params.ReqID)

		query := &metadata.QueryInput{}
		query.Condition = map[string]interface{}{
			"bk_asst_inst_id": map[string]interface{}{
				common.BKDBIN: asstInstIDS,
			},
			"bk_asst_obj_id": keyObjID,
			"bk_obj_id":      object.ObjectID,
		}

		asstInst, err := c.asst.SearchInstAssociation(params, query)
		if nil != err {
			blog.Errorf("[operation-inst] failed to search the association inst, err: %s, rid: %s", err.Error(), params.ReqID)
			return 0, nil, err
		}

		for _, asst := range asstInst {
			targetInstIDS = append(targetInstIDS, asst.InstID)
		}
		blog.V(4).Infof("[FindInstByAssociationInst] search association, objectID=%s, keyObjID=%s, condition: %v, results: %v, rid: %s", object.ObjectID, keyObjID, query, targetInstIDS, params.ReqID)

	}

	if 0 != len(targetInstIDS) {
		instCond[obj.GetInstIDFieldName()] = map[string]interface{}{
			common.BKDBIN: targetInstIDS,
		}
	} else if 0 != len(asstParamCond.Condition) {
		if _, ok := asstParamCond.Condition[object.ObjectID]; !ok {
			instCond[obj.GetInstIDFieldName()] = map[string]interface{}{
				common.BKDBIN: targetInstIDS,
			}
		}
	}

	query := &metadata.QueryInput{}
	query.Condition = instCond
	if fields, ok := asstParamCond.Fields[object.ObjectID]; ok {
		query.Fields = strings.Join(fields, ",")
	}
	query.Limit = asstParamCond.Page.Limit
	query.Sort = asstParamCond.Page.Sort
	query.Start = asstParamCond.Page.Start
	blog.V(4).Infof("[FindInstByAssociationInst] search object[%s] with inst condition: %v, rid: %s", object.ObjectID, instCond, params.ReqID)
	return c.FindInst(params, obj, query, false)
}

func (c *commonInst) FindOriginInst(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (*metadata.InstResult, error) {
	switch obj.Object().ObjectID {
	case common.BKInnerObjIDHost:
		rsp, err := c.clientSet.CoreService().Host().GetHosts(context.Background(), params.Header, cond)
		if nil != err {
			blog.Errorf("[operation-inst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
			return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {

			blog.Errorf("[operation-inst] failed to delete the object(%s) inst by the condition(%#v), err: %s, rid: %s", obj.Object().ObjectID, cond, rsp.ErrMsg, params.ReqID)
			return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
		}

		return &metadata.InstResult{Count: rsp.Data.Count, Info: mapstr.NewArrayFromMapStr(rsp.Data.Info)}, nil

	default:
		queryCond, err := mapstr.NewFromInterface(cond.Condition)
		input := &metadata.QueryCondition{Condition: queryCond}
		input.Limit.Offset = int64(cond.Start)
		input.Limit.Limit = int64(cond.Limit)
		input.Fields = strings.Split(cond.Fields, ",")
		input.SortArr = metadata.NewSearchSortParse().String(cond.Sort).ToSearchSortArr()
		rsp, err := c.clientSet.CoreService().Instance().ReadInstance(context.Background(), params.Header, obj.GetObjectID(), input)
		if nil != err {
			blog.Errorf("[operation-inst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
			return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("[operation-inst] failed to delete the object(%s) inst by the condition(%#v), err: %s, rid: %s", obj.Object().ObjectID, cond, rsp.ErrMsg, params.ReqID)
			return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
		}
		return &metadata.InstResult{Info: rsp.Data.Info, Count: rsp.Data.Count}, nil
	}
}

func (c *commonInst) FindInst(params types.ContextParams, obj model.Object, cond *metadata.QueryInput, needAsstDetail bool) (count int, results []inst.Inst, err error) {
	rsp, err := c.FindOriginInst(params, obj, cond)
	if nil != err {
		blog.Errorf("[operation-inst] failed to find origin inst , err: %s, rid: %s", err.Error(), params.ReqID)
		return 0, nil, err
	}

	return rsp.Count, inst.CreateInst(params, c.clientSet, obj, rsp.Info), nil
}

func (c *commonInst) UpdateInst(params types.ContextParams, data mapstr.MapStr, obj model.Object, cond condition.Condition, instID int64) error {
	// not allowed to update these fields, need to use specialized function
	data.Remove(common.BKParentIDField)
	data.Remove(common.BKAppIDField)
	data.Remove(metadata.BKMetadata)
	// update association
	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit
	if 0 < instID {
		innerCond := condition.CreateCondition()
		innerCond.Field(obj.GetInstIDFieldName()).Eq(instID)
		query.Condition = innerCond.ToMapStr()
	}

	// update insts
	fCond := cond.ToMapStr()
	if nil != params.MetaData {
		fCond.Set(metadata.BKMetadata, *params.MetaData)
	}
	inputParams := metadata.UpdateOption{
		Data:      data,
		Condition: fCond,
	}

	preAuditLog := NewSupplementary().Audit(params, c.clientSet, obj, c).CreateSnapshot(-1, fCond)
	rsp, err := c.clientSet.CoreService().Instance().UpdateInstance(params.Context, params.Header, obj.GetObjectID(), &inputParams)
	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, err: %s, rid: %s", err.Error(), params.ReqID)
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-inst] failed to set the object(%s) inst by the condition(%#v), err: %s, rid: %s", obj.Object().ObjectID, fCond, rsp.ErrMsg, params.ReqID)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}
	currAuditLog := NewSupplementary().Audit(params, c.clientSet, obj, c).CreateSnapshot(-1, cond.ToMapStr())
	NewSupplementary().Audit(params, c.clientSet, obj, c).CommitUpdateLog(preAuditLog, currAuditLog, nil, nil)
	return nil
}
