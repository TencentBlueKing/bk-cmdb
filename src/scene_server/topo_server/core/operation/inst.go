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
	"fmt"
	"sort"
	"strconv"
	"strings"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	gparams "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
)

// InstOperationInterface inst operation methods
type InstOperationInterface interface {
	CreateInst(kit *rest.Kit, obj model.Object, data mapstr.MapStr) (inst.Inst, error)
	CreateInstBatch(kit *rest.Kit, obj model.Object, batchInfo *InstBatchInfo) (*BatchResult, error)
	DeleteInst(kit *rest.Kit, objectID string, cond mapstr.MapStr, needCheckHost bool) error
	DeleteMainlineInstWithID(kit *rest.Kit, obj model.Object, instID int64) error
	DeleteInstByInstID(kit *rest.Kit, objectID string, instID []int64, needCheckHost bool) error
	FindOriginInst(kit *rest.Kit, objID string, cond *metadata.QueryInput) (*metadata.InstResult, errors.CCError)
	FindInst(kit *rest.Kit, obj model.Object, cond *metadata.QueryInput, needAsstDetail bool) (count int, results []inst.Inst, err error)
	FindInstByAssociationInst(kit *rest.Kit, objID string, asstParamCond *AssociationParams) (*metadata.InstResult, error)
	FindInstChildTopo(kit *rest.Kit, obj model.Object, instID int64, query *metadata.QueryInput) (count int, results []*CommonInstTopo, err error)
	FindInstParentTopo(kit *rest.Kit, obj model.Object, instID int64, query *metadata.QueryInput) (count int, results []*CommonInstTopo, err error)
	FindInstTopo(kit *rest.Kit, obj model.Object, instID int64, query *metadata.QueryInput) (count int, results []CommonInstTopoV2, err error)
	UpdateInst(kit *rest.Kit, data mapstr.MapStr, obj model.Object, cond condition.Condition, instID int64) error

	SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface, obj ObjectOperationInterface)
}

// NewInstOperation create a new inst operation instance
func NewInstOperation(client apimachinery.ClientSetInterface, languageIf language.CCLanguageIf, authManager *extensions.AuthManager) InstOperationInterface {
	return &commonInst{
		clientSet:   client,
		language:    languageIf,
		authManager: authManager,
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
	language     language.CCLanguageIf
	authManager  *extensions.AuthManager
}

func (c *commonInst) SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface, obj ObjectOperationInterface) {
	c.modelFactory = modelFactory
	c.instFactory = instFactory
	c.asst = asst
	c.obj = obj
}

// CreateInstBatch
func (c *commonInst) CreateInstBatch(kit *rest.Kit, obj model.Object, batchInfo *InstBatchInfo) (*BatchResult, error) {
	object := obj.Object()

	// forbidden create inner model instance with common api
	if common.IsInnerModel(object.ObjectID) == true {
		blog.V(5).Infof("CreateInstBatch failed, create %s instance with common create api forbidden, rid: %s", object.ObjectID, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrTopoImportMainlineForbidden)
	}

	isMainlin, err := obj.IsMainlineObject()
	if err != nil {
		blog.Errorf("[operation-inst] failed to get if the object(%s) is mainline object, err: %s, rid: %s", object.ObjectID, err.Error(), kit.Rid)
		return nil, err
	}
	if isMainlin {
		blog.V(5).Infof("CreateInstBatch failed, create %s instance with common create api forbidden, rid: %s", object.ObjectID, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrTopoImportMainlineForbidden)

	}

	results := &BatchResult{}
	colIdxErrMap := map[int]string{}
	colIdxList := []int{}
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
			blog.Errorf("create object[%s] instance batch failed, because bk_obj_id field conflict with url field, rid: %s", object.ObjectID, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrorTopoObjectInstanceObjIDFieldConflictWithURL, line)
		}
	}

	nonInnerAttributes, err := obj.GetNonInnerAttributes()
	if err != nil {
		blog.Errorf("[audit]failed to get the object(%s)' attribute, err: %s, rid: %s", obj.Object().ObjectID, err.Error(), kit.Rid)
		return nil, err
	}

	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	updateAuditLogs := make([]metadata.AuditLog, 0)
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
		item := c.instFactory.CreateInst(kit, obj)

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
			filter := condition.CreateCondition()
			filter = filter.Field(idFieldname).Eq(instID)

			// remove unchangeable fields.
			delete(colInput, idFieldname)
			delete(colInput, common.BKParentIDField)
			delete(colInput, common.BKAppIDField)

			// generate audit log of instance.
			generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(colInput)
			auditLog, ccErr := audit.GenerateAuditLogByCondGetData(generateAuditParameter, obj.GetObjectID(), filter.ToMapStr())
			if ccErr != nil {
				blog.Errorf(" update inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
				return nil, ccErr
			}
			updateAuditLogs = append(updateAuditLogs, auditLog...)

			// to update.
			err = item.UpdateInstance(filter, colInput, nonInnerAttributes)
			if nil != err {
				blog.Errorf("[operation-inst] failed to update the object(%s) inst data (%#v), err: %s, rid: %s", object.ObjectID, colInput, err.Error(), kit.Rid)
				errStr := c.language.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)).Languagef("import_row_int_error_str", colIdx, err.Error())
				colIdxList = append(colIdxList, int(colIdx))
				colIdxErrMap[int(colIdx)] = errStr
				continue
			}
			instID, err := item.GetInstID()
			if err != nil {
				blog.ErrorJSON("update inst success, but get id field failed, inst: %s, err: %s, rid: %s", item.GetValues(), err.Error(), kit.Rid)
				errStr := c.language.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)).Languagef("import_row_int_error_str", colIdx, err.Error())
				colIdxList = append(colIdxList, int(colIdx))
				colIdxErrMap[int(colIdx)] = errStr
				continue
			}
			updatedInstanceIDs = append(updatedInstanceIDs, instID)
			results.Success = append(results.Success, strconv.FormatInt(colIdx, 10))
			continue
		}

		// set data
		// call CoreService.CreateInstance
		err = item.Create()
		if nil != err {
			blog.Errorf("[operation-inst] failed to save the object(%s) inst data (%#v), err: %s, rid: %s", object.ObjectID, colInput, err.Error(), kit.Rid)
			errStr := c.language.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)).Languagef("import_row_int_error_str", colIdx, err.Error())
			colIdxList = append(colIdxList, int(colIdx))
			colIdxErrMap[int(colIdx)] = errStr
			continue
		}
		results.Success = append(results.Success, strconv.FormatInt(colIdx, 10))

		instanceID, err := item.GetInstID()
		if err != nil {
			blog.Errorf("unexpected error, instances created success, but get id failed, err: %+v, rid: %s", err, kit.Rid)
			continue
		}
		createdInstanceIDs = append(createdInstanceIDs, instanceID)
	}

	// generate audit log of instance.
	cond := map[string]interface{}{
		obj.GetInstIDFieldName(): map[string]interface{}{
			common.BKDBIN: createdInstanceIDs,
		},
	}
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLogByCondGetData(generateAuditParameter, obj.GetObjectID(), cond)
	if err != nil {
		blog.Errorf(" creat inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	// save audit log.
	err = audit.SaveAuditLog(kit, append(updateAuditLogs, auditLog...)...)
	if err != nil {
		blog.Errorf("creat inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	results.SuccessCreated = createdInstanceIDs
	results.SuccessUpdated = updatedInstanceIDs
	sort.Strings(results.Success)

	//sort error
	sort.Ints(colIdxList)
	for colIdx := range colIdxList {
		results.Errors = append(results.Errors, colIdxErrMap[colIdx])
	}

	return results, nil
}

func (c *commonInst) isValidBizInstID(kit *rest.Kit, obj metadata.Object, instID int64, bizID int64) error {

	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)

	if bizID != 0 {
		cond.Field(common.BKAppIDField).Eq(bizID)
	}

	if obj.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(obj.ObjectID)
	}

	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit

	rsp, err := c.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, obj.GetObjectID(), &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-inst] faild to read the object(%s) inst by the condition(%#v), err: %s, rid: %s", obj.ObjectID, cond, rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	if rsp.Data.Count > 0 {
		return nil
	}

	return kit.CCError.Error(common.CCErrTopoInstSelectFailed)
}

func (c *commonInst) isValidInstID(kit *rest.Kit, obj metadata.Object, instID int64) error {
	return c.isValidBizInstID(kit, obj, instID, 0)
}

func (c *commonInst) validMainLineParentID(kit *rest.Kit, obj model.Object, data mapstr.MapStr) error {
	if obj.Object().ObjectID == common.BKInnerObjIDApp {
		return nil
	}
	def, exist := data.Get(common.BKDefaultField)
	if exist && def.(int) != common.DefaultFlagDefaultValue {
		return nil
	}
	parent, err := obj.GetMainlineParentObject()
	if err != nil {
		blog.Errorf("[operation-inst] failed to get the object(%s) mainline parent, err: %s, rid: %s", obj.Object().ObjectID, err.Error(), kit.Rid)
		return err
	}
	bizID, err := data.Int64(common.BKAppIDField)
	if err != nil {
		bizID, err = metadata.ParseBizIDFromData(data)
		if err != nil {
			blog.Errorf("[operation-inst]failed to parse the biz id, err: %s, rid: %s", err.Error(), kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField)
		}
	}
	parentID, err := data.Int64(common.BKParentIDField)
	if err != nil {
		blog.Errorf("[operation-inst]failed to parse the parent id, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKParentIDField)
	}
	if err = c.isValidBizInstID(kit, parent.Object(), parentID, bizID); err != nil {
		blog.Errorf("[operation-inst]parent id %d is invalid, err: %s, rid: %s", parentID, err.Error(), kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKParentIDField)
	}
	return nil
}

func (c *commonInst) CreateInst(kit *rest.Kit, obj model.Object, data mapstr.MapStr) (inst.Inst, error) {

	// create new insts
	item := c.instFactory.CreateInst(kit, obj)
	item.SetValues(data)

	iData := item.ToMapStr()
	if obj.Object().ObjectID == common.BKInnerObjIDPlat {
		iData["bk_supplier_account"] = kit.SupplierAccount
	}

	isMainline, err := obj.IsMainlineObject()
	if err != nil {
		blog.Errorf("[operation-inst] failed to get if the object(%s) is mainline object, err: %s, rid: %s", obj.Object().ObjectID, err.Error(), kit.Rid)
		return nil, err
	}
	if isMainline {
		if err := c.validMainLineParentID(kit, obj, data); nil != err {
			blog.Errorf("[operation-inst] the mainline object(%s) parent id invalid, err: %s, rid: %s", obj.Object().ObjectID, err.Error(), kit.Rid)
			return nil, err
		}
	}

	if err := item.Create(); nil != err {
		blog.Errorf("[operation-inst] failed to save the object(%s) inst data (%#v), err: %s, rid: %s", obj.Object().ObjectID, data, err.Error(), kit.Rid)
		return nil, err
	}

	instID, err := item.GetInstID()
	if err != nil {
		return nil, kit.CCError.Error(common.CCErrTopoInstCreateFailed)
	}
	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)
	_, insts, err := c.FindInst(kit, obj, &metadata.QueryInput{Condition: cond.ToMapStr()}, false)
	if err != nil {
		return nil, err
	}

	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	for _, inst := range insts {
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, obj.GetObjectID(), []mapstr.MapStr{inst.GetValues()})
		if err != nil {
			blog.Errorf(" creat inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		err = audit.SaveAuditLog(kit, auditLog...)
		if err != nil {
			blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
		}
		return inst, nil
	}

	return item, nil
}

func (c *commonInst) innerHasHost(kit *rest.Kit, moduleIDS []int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ModuleIDArr: moduleIDS,
		Fields:      []string{common.BKHostIDField},
		Page:        metadata.BasePage{Limit: 1},
	}
	rsp, err := c.clientSet.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, option)
	if nil != err {
		blog.Errorf("[operation-module] failed to request the object controller, err: %s, rid: %s", err.Error(), kit.Rid)
		return false, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-module]  failed to search the host module configures, err: %s, rid: %s", rsp.ErrMsg, kit.Rid)
		return false, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != len(rsp.Data.Info), nil
}

// hasHost get objID and instances map for mainline instances with its children topology, and check if they have hosts
func (c *commonInst) hasHost(kit *rest.Kit, instances []mapstr.MapStr, objID string, checkHost bool) (
	map[string][]mapstr.MapStr, bool, error) {

	if len(instances) == 0 {
		return nil, false, nil
	}

	objInstMap := map[string][]mapstr.MapStr{
		objID: instances,
	}

	instIDs := make([]int64, len(instances))
	for index, instance := range instances {
		instID, err := instance.Int64(common.GetInstIDField(objID))
		if err != nil {
			blog.ErrorJSON("can not convert ID to int64, err: %s, inst: %s, rid: %s", err, instance, kit.Rid)
			return nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.GetInstIDField(objID))
		}
		instIDs[index] = instID
	}

	var moduleIDs []int64
	if objID == common.BKInnerObjIDModule {
		moduleIDs = instIDs
	} else if objID == common.BKInnerObjIDSet {
		query := &metadata.QueryInput{
			Condition: map[string]interface{}{common.BKSetIDField: map[string]interface{}{common.BKDBIN: instIDs}},
			Limit:     common.BKNoLimit,
		}

		moduleRsp, err := c.FindOriginInst(kit, common.BKInnerObjIDModule, query)
		if nil != err {
			blog.Errorf("find modules for set failed, err: %v, set IDs: %+v, rid: %s", err, instIDs, kit.Rid)
			return nil, false, err
		}

		if len(moduleRsp.Info) == 0 {
			return objInstMap, false, nil
		}

		objInstMap[common.BKInnerObjIDModule] = moduleRsp.Info
		moduleIDs = make([]int64, len(moduleRsp.Info))
		for index, module := range moduleRsp.Info {
			moduleID, err := module.Int64(common.BKModuleIDField)
			if err != nil {
				blog.ErrorJSON("can not convert ID to int64, err: %s, module: %s, rid: %s", err, module, kit.Rid)
				return nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKModuleIDField)
			}
			moduleIDs[index] = moduleID
		}
	} else {
		// get mainline object relation(excluding hosts) by mainline associations
		mainlineCond := &metadata.QueryCondition{
			Condition: map[string]interface{}{common.AssociationKindIDField: common.AssociationKindMainline},
		}
		asstRsp, err := c.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, mainlineCond)
		if nil != err {
			blog.Errorf("search mainline association failed, error: %v, rid: %s", err, kit.Rid)
			return nil, false, err
		}

		objChildMap := make(map[string]string)
		isMainline := false
		for _, asst := range asstRsp.Data.Info {
			if asst.ObjectID == common.BKInnerObjIDHost {
				continue
			}
			objChildMap[asst.AsstObjID] = asst.ObjectID
			if asst.AsstObjID == objID || asst.ObjectID == objID {
				isMainline = true
			}
		}

		if !isMainline {
			return objInstMap, false, nil
		}

		// loop through the child topology level to get all instances
		parentIDs := instIDs
		for childObjID := objChildMap[objID]; len(childObjID) != 0; childObjID = objChildMap[childObjID] {
			cond := map[string]interface{}{common.BKParentIDField: map[string]interface{}{common.BKDBIN: parentIDs}}
			if metadata.IsCommon(childObjID) {
				cond[metadata.ModelFieldObjectID] = childObjID
			}

			if childObjID == common.BKInnerObjIDSet {
				cond[common.BKDefaultField] = common.DefaultFlagDefaultValue
			}

			query := &metadata.QueryInput{
				Condition: cond,
				Limit:     common.BKNoLimit,
			}

			childRsp, err := c.FindOriginInst(kit, childObjID, query)
			if nil != err {
				blog.Errorf("find children failed, err: %v, parent IDs: %+v, rid: %s", err, parentIDs, kit.Rid)
				return nil, false, err
			}

			if len(childRsp.Info) == 0 {
				return objInstMap, false, nil
			}

			parentIDs = make([]int64, len(childRsp.Info))
			for index, instance := range childRsp.Info {
				instID, err := instance.Int64(common.GetInstIDField(childObjID))
				if err != nil {
					blog.ErrorJSON("can not convert ID to int64, err: %s, inst: %s, rid: %s", err, instance, kit.Rid)
					return nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.GetInstIDField(childObjID))
				}
				parentIDs[index] = instID
			}

			if childObjID == common.BKInnerObjIDModule {
				moduleIDs = parentIDs
			}

			objInstMap[childObjID] = childRsp.Info
		}
	}

	// check if module contains hosts
	if checkHost && len(moduleIDs) > 0 {
		exists, err := c.innerHasHost(kit, moduleIDs)
		if nil != err {
			return nil, false, err
		}

		if exists {
			return nil, true, nil
		}
	}

	return objInstMap, false, nil
}

func (c *commonInst) DeleteInstByInstID(kit *rest.Kit, objectID string, instID []int64, needCheckHost bool) error {
	cond := map[string]interface{}{
		common.GetInstIDField(objectID): map[string]interface{}{common.BKDBIN: instID},
	}
	if metadata.IsCommon(objectID) {
		cond[common.BKObjIDField] = objectID
	}

	return c.deleteInstByCond(kit, objectID, cond, needCheckHost)
}

func (c *commonInst) deleteInstByCond(kit *rest.Kit, objectID string, cond mapstr.MapStr, needCheckHost bool) error {
	query := &metadata.QueryInput{
		Condition: cond,
		Limit:     common.BKNoLimit,
	}

	instRsp, err := c.FindOriginInst(kit, objectID, query)
	if nil != err {
		return err
	}

	if len(instRsp.Info) == 0 {
		return nil
	}

	delObjInstsMap, exists, err := c.hasHost(kit, instRsp.Info, objectID, needCheckHost)
	if nil != err {
		return err
	}
	if exists {
		return kit.CCError.Error(common.CCErrTopoHasHostCheckFailed)
	}

	bizSetMap := make(map[int64][]int64)
	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	auditLogs := make([]metadata.AuditLog, 0)

	for objID, delInsts := range delObjInstsMap {
		delInstIDs := make([]int64, len(delInsts))
		for index, instance := range delInsts {
			instID, err := instance.Int64(common.GetInstIDField(objID))
			if err != nil {
				blog.ErrorJSON("can not convert ID to int64, err: %s, inst: %s, rid: %s", err, instance, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.GetInstIDField(objID))
			}
			delInstIDs[index] = instID

			if objID == common.BKInnerObjIDSet {
				bizID, err := instance.Int64(common.BKAppIDField)
				if err != nil {
					blog.ErrorJSON("can not convert biz ID to int64, err: %s, set: %s, rid: %s", err, instance, kit.Rid)
					return kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)
				}
				bizSetMap[bizID] = append(bizSetMap[bizID], instID)
			}
		}

		// if any instance has been bind to a instance by the association, then these instances should not be deleted.
		err := c.asst.CheckAssociations(kit, objID, delInstIDs)
		if nil != err {
			return err
		}

		// generate audit log.
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, objID, delInsts)
		if err != nil {
			blog.Errorf(" delete inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
		auditLogs = append(auditLogs, auditLog...)

		// delete this instance now.
		delCond := map[string]interface{}{
			common.GetInstIDField(objID): map[string]interface{}{common.BKDBIN: delInstIDs},
		}
		if metadata.IsCommon(objID) {
			delCond[common.BKObjIDField] = objID
		}
		dc := &metadata.DeleteOption{Condition: delCond}
		rsp, err := c.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, objID, dc)
		if nil != err {
			blog.ErrorJSON("delete inst failed, err: %s, cond: %s rid: %s", err, delCond, kit.Rid)
			return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if err := rsp.CCError(); err != nil {
			blog.ErrorJSON("delete inst failed, err: %s, cond: %s rid: %s", err, delCond, kit.Rid)
			return err
		}
	}

	// clear set template sync status for set instances
	for bizID, setIDs := range bizSetMap {
		if len(setIDs) != 0 {
			if ccErr := c.clientSet.CoreService().SetTemplate().DeleteSetTemplateSyncStatus(kit.Ctx, kit.Header, bizID, setIDs); ccErr != nil {
				blog.Errorf("[operation-set] failed to delete set template sync status failed, bizID: %d, setIDs: %+v, err: %s, rid: %s", bizID, setIDs, ccErr.Error(), kit.Rid)
				return ccErr
			}
		}
	}

	err = audit.SaveAuditLog(kit, auditLogs...)
	if err != nil {
		blog.Errorf("delete inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

func (c *commonInst) DeleteMainlineInstWithID(kit *rest.Kit, obj model.Object, instID int64) error {
	object := obj.Object()
	// if this instance has been bind to a instance by the association, then this instance should not be deleted.
	err := c.asst.CheckAssociation(kit, object.ObjectID, instID)
	if nil != err {
		return err
	}

	// delete this instance now.
	delCond := condition.CreateCondition()
	delCond.Field(obj.GetInstIDFieldName()).Eq(instID)
	if obj.IsCommon() {
		delCond.Field(common.BKObjIDField).Eq(object.ObjectID)
	}

	// generate audit log.
	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	auditLog, err := audit.GenerateAuditLogByCondGetData(generateAuditParameter, obj.GetObjectID(), delCond.ToMapStr())
	if err != nil {
		blog.Errorf(" delete inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	// to delete.
	ops := metadata.DeleteOption{
		Condition: delCond.ToMapStr(),
	}
	rsp, err := c.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, object.ObjectID, &ops)
	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, err: %s", err.Error())
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-inst] failed to delete the object(%s) inst by the condition(%#v), err: %s", object.ObjectID, delCond.ToMapStr(), rsp.ErrMsg)
		return kit.CCError.Error(rsp.Code)
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, auditLog...); err != nil {
		blog.Errorf("delete inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

func (c *commonInst) DeleteInst(kit *rest.Kit, objectID string, cond mapstr.MapStr, needCheckHost bool) error {
	return c.deleteInstByCond(kit, objectID, cond, needCheckHost)
}

func (c *commonInst) convertInstIDIntoStruct(kit *rest.Kit, asstObj metadata.Association, instIDS []string, needAsstDetail bool, modelBizID int64) ([]metadata.InstNameAsst, error) {

	obj, err := c.obj.FindSingleObject(kit, asstObj.AsstObjID)
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
	rsp, err := c.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, obj.GetObjectID(), query)

	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-inst] faild to delete the object(%s) inst by the condition(%#v), err: %s, rid: %s", object.ObjectID, cond, rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
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

func (c *commonInst) searchAssociationInst(kit *rest.Kit, objID string, query *metadata.QueryInput) ([]int64, error) {

	obj, err := c.obj.FindSingleObject(kit, objID)
	if nil != err {
		return nil, err
	}

	_, insts, err := c.FindInst(kit, obj, query, false)
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

func (c *commonInst) FindInstChildTopo(kit *rest.Kit, obj model.Object, instID int64, query *metadata.QueryInput) (count int, results []*CommonInstTopo, err error) {
	results = make([]*CommonInstTopo, 0)
	if nil == query {
		query = &metadata.QueryInput{}
		cond := condition.CreateCondition()
		cond.Field(obj.GetInstIDFieldName()).Eq(instID)
		query.Condition = cond.ToMapStr()
	}

	_, insts, err := c.FindInst(kit, obj, query, false)
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

func (c *commonInst) FindInstParentTopo(kit *rest.Kit, obj model.Object, instID int64, query *metadata.QueryInput) (count int, results []*CommonInstTopo, err error) {

	results = make([]*CommonInstTopo, 0)
	if nil == query {
		query = &metadata.QueryInput{}
		cond := condition.CreateCondition()
		cond.Field(obj.GetInstIDFieldName()).Eq(instID)
		query.Condition = cond.ToMapStr()
	}

	_, insts, err := c.FindInst(kit, obj, query, false)
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

func (c *commonInst) FindInstTopo(kit *rest.Kit, obj model.Object, instID int64, query *metadata.QueryInput) (count int, results []CommonInstTopoV2, err error) {

	if nil == query {
		query = &metadata.QueryInput{}
		cond := condition.CreateCondition()
		cond.Field(obj.GetInstIDFieldName()).Eq(instID)
		query.Condition = cond.ToMapStr()
	}

	_, insts, err := c.FindInst(kit, obj, query, false)
	if nil != err {
		blog.Errorf("[operation-inst] failed to find the inst, err: %s, rid: %s", err.Error(), kit.Rid)
		return 0, nil, err
	}

	for _, inst := range insts {
		id, err := inst.GetInstID()
		if nil != err {
			blog.Errorf("[operation-inst] failed to find the inst, err: %s, rid: %s", err.Error(), kit.Rid)
			return 0, nil, err
		}

		name, err := inst.GetInstName()
		if nil != err {
			blog.Errorf("[operation-inst] failed to find the inst, err: %s, rid: %s", err.Error(), kit.Rid)
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

		_, parentInsts, err := c.FindInstParentTopo(kit, inst.GetObject(), id, nil)
		if nil != err {
			blog.Errorf("[operation-inst] failed to find the inst, err: %s, rid: %s", err.Error(), kit.Rid)
			return 0, nil, err
		}

		_, childInsts, err := c.FindInstChildTopo(kit, inst.GetObject(), id, nil)
		if nil != err {
			blog.Errorf("[operation-inst] failed to find the inst, err: %s, rid: %s", err.Error(), kit.Rid)
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

func (c *commonInst) FindInstByAssociationInst(kit *rest.Kit, objID string, asstParamCond *AssociationParams) (*metadata.InstResult, error) {

	instCond := map[string]interface{}{}
	if metadata.IsCommon(objID) {
		instCond[common.BKObjIDField] = objID
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
				if objID == keyObjID {
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
				if objID == keyObjID {
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

		if objID == keyObjID {
			// no need to search the association objects
			continue
		}

		innerCond := new(metadata.QueryInput)
		if fields, ok := asstParamCond.Fields[keyObjID]; ok {
			innerCond.Fields = strings.Join(fields, ",")
		}
		innerCond.Condition = cond

		asstInstIDS, err := c.searchAssociationInst(kit, keyObjID, innerCond)
		if nil != err {
			blog.Errorf("[operation-inst]failed to search the association inst, err: %s, rid: %s", err.Error(), kit.Rid)
			return nil, err
		}
		blog.V(4).Infof("[FindInstByAssociationInst] search association insts, keyObjID %s, condition: %v, results: %v, rid: %s", keyObjID, innerCond, asstInstIDS, kit.Rid)

		query := &metadata.QueryInput{}
		query.Condition = map[string]interface{}{
			"bk_asst_inst_id": map[string]interface{}{
				common.BKDBIN: asstInstIDS,
			},
			"bk_asst_obj_id": keyObjID,
			"bk_obj_id":      objID,
		}

		asstInst, err := c.asst.SearchInstAssociation(kit, query)
		if nil != err {
			blog.Errorf("[operation-inst] failed to search the association inst, err: %s, rid: %s", err.Error(), kit.Rid)
			return nil, err
		}

		for _, asst := range asstInst {
			targetInstIDS = append(targetInstIDS, asst.InstID)
		}
		blog.V(4).Infof("[FindInstByAssociationInst] search association, objectID=%s, keyObjID=%s, condition: %v, results: %v, rid: %s", objID, keyObjID, query, targetInstIDS, kit.Rid)
	}

	if 0 != len(targetInstIDS) {
		instCond[metadata.GetInstIDFieldByObjID(objID)] = map[string]interface{}{
			common.BKDBIN: targetInstIDS,
		}
	} else if 0 != len(asstParamCond.Condition) {
		if _, ok := asstParamCond.Condition[objID]; !ok {
			instCond[metadata.GetInstNameFieldName(objID)] = map[string]interface{}{
				common.BKDBIN: targetInstIDS,
			}
		}
	}

	query := &metadata.QueryInput{}
	query.Condition = instCond
	if fields, ok := asstParamCond.Fields[objID]; ok {
		query.Fields = strings.Join(fields, ",")
	}
	query.Limit = asstParamCond.Page.Limit
	query.Sort = asstParamCond.Page.Sort
	query.Start = asstParamCond.Page.Start
	blog.V(4).Infof("[FindInstByAssociationInst] search object[%s] with inst condition: %v, rid: %s", objID, instCond, kit.Rid)
	return c.FindOriginInst(kit, objID, query)
}

func (c *commonInst) FindOriginInst(kit *rest.Kit, objID string, cond *metadata.QueryInput) (*metadata.InstResult, errors.CCError) {
	switch objID {
	case common.BKInnerObjIDHost:
		rsp, err := c.clientSet.CoreService().Host().GetHosts(kit.Ctx, kit.Header, cond)
		if nil != err {
			blog.Errorf("[operation-inst] failed to request object controller, err: %s, rid: %s", err.Error(), kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("[operation-inst] failed to delete the object(%s) inst by the condition(%#v), err: %s, rid: %s", objID, cond, rsp.ErrMsg, kit.Rid)
			return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
		}

		return &metadata.InstResult{Count: rsp.Data.Count, Info: mapstr.NewArrayFromMapStr(rsp.Data.Info)}, nil

	default:
		queryCond, err := mapstr.NewFromInterface(cond.Condition)
		input := &metadata.QueryCondition{Condition: queryCond}
		input.Page.Start = cond.Start
		input.Page.Limit = cond.Limit
		input.Page.Sort = cond.Sort
		input.Fields = strings.Split(cond.Fields, ",")
		rsp, err := c.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, input)
		if nil != err {
			blog.Errorf("[operation-inst] failed to request object controller, err: %s, rid: %s", err.Error(), kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("[operation-inst] failed to delete the object(%s) inst by the condition(%#v), err: %s, rid: %s", objID, cond, rsp.ErrMsg, kit.Rid)
			return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
		}
		return &metadata.InstResult{Info: rsp.Data.Info, Count: rsp.Data.Count}, nil
	}
}

func (c *commonInst) FindInst(kit *rest.Kit, obj model.Object, cond *metadata.QueryInput, needAsstDetail bool) (count int, results []inst.Inst, err error) {
	rsp, err := c.FindOriginInst(kit, obj.GetObjectID(), cond)
	if nil != err {
		blog.Errorf("[operation-inst] failed to find origin inst , err: %s, rid: %s", err.Error(), kit.Rid)
		return 0, nil, err
	}

	return rsp.Count, inst.CreateInst(kit, c.clientSet, obj, rsp.Info), nil
}

func (c *commonInst) UpdateInst(kit *rest.Kit, data mapstr.MapStr, obj model.Object, cond condition.Condition, instID int64) error {
	// not allowed to update these fields, need to use specialized function
	data.Remove(common.BKParentIDField)
	data.Remove(common.BKAppIDField)

	// update association
	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit
	if 0 < instID {
		innerCond := condition.CreateCondition()
		innerCond.Field(obj.GetInstIDFieldName()).Eq(instID)
		query.Condition = innerCond.ToMapStr()
	}

	fCond := cond.ToMapStr()
	inputParams := metadata.UpdateOption{
		Data:      data,
		Condition: fCond,
	}

	// generate audit log of instance.
	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(data)
	auditLog, ccErr := audit.GenerateAuditLogByCondGetData(generateAuditParameter, obj.GetObjectID(), fCond)
	if ccErr != nil {
		blog.Errorf(" update inst, generate audit log failed, err: %v, rid: %s", ccErr, kit.Rid)
		return ccErr
	}

	// to update.
	rsp, err := c.clientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, obj.GetObjectID(), &inputParams)
	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !rsp.Result {
		blog.Errorf("[operation-inst] failed to set the object(%s) inst by the condition(%#v), err: %s, rid: %s", obj.Object().ObjectID, fCond, rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	// save audit log.
	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}
	return nil
}
