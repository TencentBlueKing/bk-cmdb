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

package inst

import (
	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/version"
	"strings"
)

// ModuleOperationInterface module operation methods
type ModuleOperationInterface interface {
	CreateModule(kit *rest.Kit, obj metadata.Object, bizID, setID int64,
		data mapstr.MapStr) (*metadata.CreateOneDataResult, error)
	DeleteModule(kit *rest.Kit, bizID int64, setID, moduleIDS []int64) error
	FindModule(kit *rest.Kit, objID string, cond *metadata.QueryInput) (count int, results []mapstr.MapStr, err error)
	UpdateModule(kit *rest.Kit, data mapstr.MapStr, obj *metadata.Object, bizID, setID, moduleID int64) error
}

// NewModuleOperation create a new module
func NewModuleOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) ModuleOperationInterface {
	return &module{
		clientSet:   client,
		authManager: authManager,
	}
}

type module struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
}

// CreateInst create instance by object and create message
func (m *module) CreateInst(kit *rest.Kit, obj metadata.Object,
	data mapstr.MapStr) (*metadata.CreateOneDataResult, error) {

	if obj.ObjectID == common.BKInnerObjIDPlat {
		data.Set(common.BkSupplierAccount, kit.SupplierAccount)
	}

	assoc, err := m.validObject(kit, obj)
	if err != nil {
		blog.Errorf("valid object (%s) failed, err: %v, rid: %s", obj.ObjectID, err, kit.Rid)
		return nil, err
	}

	if assoc != nil {
		if err := m.validMainLineParentID(kit, assoc, data); err != nil {
			blog.Errorf("the mainline object(%s) parent id invalid, err: %v, rid: %s", obj.ObjectID, err, kit.Rid)
			return nil, err
		}
	}

	data.Set(common.BKObjIDField, obj.ObjectID)

	instCond := &metadata.CreateModelInstance{Data: data}
	rsp, err := m.clientSet.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, obj.ObjectID, instCond)
	if err != nil {
		blog.Errorf("failed to create object instance, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to create object instance ,err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if rsp.Data.Created.ID == 0 {
		blog.Errorf("failed to create object instance, return nothing, rid: %s", kit.Rid)
		return nil, kit.CCError.Error(common.CCErrTopoInstCreateFailed)
	}

	data.Set(obj.GetInstIDFieldName(), rsp.Data.Created.ID)
	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	audit := auditlog.NewInstanceAudit(m.clientSet.CoreService())
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, obj.GetObjectID(), []mapstr.MapStr{data})
	if err != nil {
		blog.Errorf(" creat inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return &rsp.Data, nil
}

// validObject check object valid
func (m *module) validObject(kit *rest.Kit, obj metadata.Object) (*metadata.Association, error) {

	if !metadata.IsCommon(obj.ObjectID) {
		blog.Errorf("object (%s) isn't common object, rid: %s", obj.ID, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI)
	}

	// 暂停使用的model不允许创建实例
	if obj.IsPaused {
		blog.Errorf("object (%s) is paused, rid: %s", obj.ObjectID, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrorTopoModelStopped)
	}

	cond := mapstr.MapStr{
		common.BKObjIDField:           obj.ObjectID,
		common.AssociationKindIDField: common.AssociationKindMainline,
	}
	asst, err := m.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("search object association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if err = asst.CCError(); err != nil {
		blog.Errorf("search object association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if len(asst.Data.Info) > 1 {
		return nil, kit.CCError.CCErrorf(common.CCErrTopoGotMultipleAssociationInstance)
	}

	if len(asst.Data.Info) == 0 {
		return nil, nil
	}

	return &asst.Data.Info[0], nil
}

// validMainLineParentID check parent id is or not mainline
func (m *module) validMainLineParentID(kit *rest.Kit, assoc *metadata.Association, data mapstr.MapStr) error {
	if assoc.ObjectID == common.BKInnerObjIDApp {
		return nil
	}

	def, exist := data.Get(common.BKDefaultField)
	if exist && def.(int) != common.DefaultFlagDefaultValue {
		return nil
	}

	bizID, err := data.Int64(common.BKAppIDField)
	if err != nil {
		blog.Errorf("[operation-inst]failed to parse the biz id, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField)
	}

	parentID, err := data.Int64(common.BKParentIDField)
	if err != nil {
		blog.Errorf("[operation-inst]failed to parse the parent id, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKParentIDField)
	}

	if err = m.isValidBizInstID(kit, assoc.AsstObjID, parentID, bizID); err != nil {
		blog.Errorf("[operation-inst]parent id %d is invalid, err: %s, rid: %s", parentID, err.Error(), kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKParentIDField)
	}
	return nil
}

// isValidBizInstID check biz instance id is or not valid
func (m *module) isValidBizInstID(kit *rest.Kit, objID string, instID int64, bizID int64) error {

	cond := mapstr.MapStr{
		metadata.GetInstIDFieldByObjID(objID): instID,
	}

	if bizID != 0 {
		cond.Set(common.BKAppIDField, bizID)
	}

	if metadata.IsCommon(objID) {
		cond.Set(common.BKObjIDField, objID)
	}

	rsp, err := m.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("failed to request object controller, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to read the object(%s) inst by the condition(%#v), err: %v, rid: %s", objID, cond,
			err, kit.Rid)
		return err
	}

	if rsp.Data.Count > 0 {
		return nil
	}

	return kit.CCError.Error(common.CCErrTopoInstSelectFailed)
}

// DeleteInst delete instance by objectid and condition
func (m *module) DeleteInst(kit *rest.Kit, objectID string, cond mapstr.MapStr, needCheckHost bool) error {
	return m.deleteInstByCond(kit, objectID, cond, needCheckHost)
}

// deleteInstByCond delete instance by condition
func (m *module) deleteInstByCond(kit *rest.Kit, objectID string, cond mapstr.MapStr, needCheckHost bool) error {
	query := &metadata.QueryInput{
		Condition: cond,
		Limit:     common.BKNoLimit,
	}

	instRsp, err := m.FindInst(kit, objectID, query)
	if err != nil {
		return err
	}

	if len(instRsp.Info) == 0 {
		return nil
	}

	delObjInstsMap, exists, err := m.hasHosts(kit, instRsp.Info, objectID, needCheckHost)
	if err != nil {
		return err
	}
	if exists {
		return kit.CCError.Error(common.CCErrTopoHasHostCheckFailed)
	}

	bizSetMap := make(map[int64][]int64)
	audit := auditlog.NewInstanceAudit(m.clientSet.CoreService())
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
		input := &metadata.Condition{
			Condition: mapstr.MapStr{common.BKDBOR: []mapstr.MapStr{
				{common.BKObjIDField: objID, common.BKInstIDField: mapstr.MapStr{common.BKDBIN: delInstIDs}},
				{common.BKAsstObjIDField: objID, common.BKAsstInstIDField: mapstr.MapStr{common.BKDBIN: delInstIDs}},
			}}}
		cnt, err := m.clientSet.CoreService().Association().CountInstanceAssociations(kit.Ctx, kit.Header, objID, input)
		if err != nil {
			blog.Errorf("count instance association failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		if err = cnt.CCError(); err != nil {
			blog.Errorf("count instance association failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		if cnt.Data.Count != 0 {
			return kit.CCError.CCError(common.CCErrorInstHasAsst)
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
		rsp, err := m.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, objID, dc)
		if err != nil {
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
			if ccErr := m.clientSet.CoreService().SetTemplate().DeleteSetTemplateSyncStatus(kit.Ctx, kit.Header,
				bizID, setIDs); ccErr != nil {
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

// FindInst search instance by condition
func (m *module) FindInst(kit *rest.Kit, objID string, cond *metadata.QueryInput) (*metadata.InstResult, error) {

	result := new(metadata.InstResult)
	switch objID {
	case common.BKInnerObjIDHost:
		rsp, err := m.clientSet.CoreService().Host().GetHosts(kit.Ctx, kit.Header, cond)
		if err != nil {
			blog.Errorf("get host failed, err: %v, rid: %s", err, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if err = rsp.CCError(); err != nil {
			blog.Errorf("search object(%s) inst by the condition(%#v) failed, err: %v, rid: %s",
				objID, cond, err, kit.Rid)
			return nil, err
		}

		result.Count = rsp.Data.Count
		result.Info = rsp.Data.Info
		return result, nil

	default:
		input := &metadata.QueryCondition{Condition: cond.Condition, TimeCondition: cond.TimeCondition}
		input.Page.Start = cond.Start
		input.Page.Limit = cond.Limit
		input.Page.Sort = cond.Sort
		input.Fields = strings.Split(cond.Fields, ",")
		rsp, err := m.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, input)
		if err != nil {
			blog.Errorf("search instance failed, err: %v, rid: %s", err, kit.Rid)
			return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if err = rsp.CCError(); err != nil {
			blog.Errorf("search object(%s) inst by the condition(%#v) failed, err: %v, rid: %s",
				objID, cond, err, kit.Rid)
			return nil, err
		}

		result.Count = rsp.Data.Count
		result.Info = rsp.Data.Info
		return result, nil
	}
}

// hasHost get objID and instances map for mainline instances with its children topology, and check if they have hosts
func (m *module) hasHosts(kit *rest.Kit, instances []mapstr.MapStr, objID string, checkHost bool) (
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

		moduleRsp, err := m.FindInst(kit, common.BKInnerObjIDModule, query)
		if err != nil {
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
				blog.Errorf("can not convert ID to int64, err: %v, module: %s, rid: %s", err, module, kit.Rid)
				return nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKModuleIDField)
			}
			moduleIDs[index] = moduleID
		}
	} else {
		// get mainline object relation(excluding hosts) by mainline associations
		mainlineCond := &metadata.QueryCondition{
			Condition: map[string]interface{}{
				common.AssociationKindIDField: common.AssociationKindMainline,
				common.BKObjIDField: mapstr.MapStr{
					common.BKDBNIN: []string{common.BKInnerObjIDSet, common.BKInnerObjIDModule},
				}}}
		asstRsp, err := m.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, mainlineCond)
		if err != nil {
			blog.Errorf("search mainline association failed, error: %v, rid: %s", err, kit.Rid)
			return nil, false, err
		}

		if err = asstRsp.CCError(); err != nil {
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

			childRsp, err := m.FindInst(kit, childObjID, query)
			if err != nil {
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
					blog.Errorf("can not convert ID to int64, err: %v, inst: %s, rid: %s", err, instance, kit.Rid)
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
		exists, err := m.innerHasHost(kit, moduleIDs)
		if err != nil {
			return nil, false, err
		}

		if exists {
			return nil, true, nil
		}
	}

	return objInstMap, false, nil
}

// innerHasHost check host is or not inner ip
func (m *module) innerHasHost(kit *rest.Kit, moduleIDS []int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ModuleIDArr: moduleIDS,
		Fields:      []string{common.BKHostIDField},
		Page:        metadata.BasePage{Limit: 1},
	}
	rsp, err := m.clientSet.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, option)
	if nil != err {
		blog.Errorf("searh host object relation failed, err: %v, rid: %s", err, kit.Rid)
		return false, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to search the host module configures, err: %v, rid: %s", err, kit.Rid)
		return false, err
	}

	return 0 != len(rsp.Data.Info), nil
}

// UpdateInst update instance by condition
func (m *module) UpdateInst(kit *rest.Kit, cond, data mapstr.MapStr, objID string) error {
	// not allowed to update these fields, need to use specialized function
	data.Remove(common.BKParentIDField)
	data.Remove(common.BKAppIDField)

	inputParams := metadata.UpdateOption{
		Data:      data,
		Condition: cond,
	}

	// generate audit log of instance.
	audit := auditlog.NewInstanceAudit(m.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(data)
	auditLog, ccErr := audit.GenerateAuditLogByCondGetData(generateAuditParameter, objID, cond)
	if ccErr != nil {
		blog.Errorf(" update inst, generate audit log failed, err: %v, rid: %s", ccErr, kit.Rid)
		return ccErr
	}

	// to update.
	rsp, err := m.clientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, objID, &inputParams)
	if err != nil {
		blog.Errorf("update instance failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if err = rsp.CCError(); err != nil {
		blog.Errorf("update the object(%s) inst by the condition(%#v) failed, err: %v, rid: %s",
			objID, cond, err, kit.Rid)
		return err
	}

	// save audit log.
	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}
	return nil
}

// validBizSetID check biz, set id is or not valid
func (m *module) validBizSetID(kit *rest.Kit, bizID int64, setID int64) error {
	cond := condition.CreateCondition()
	cond.Field(common.BKSetIDField).Eq(setID)
	or := cond.NewOR()
	or.Item(mapstr.MapStr{common.BKAppIDField: bizID})

	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit

	rsp, err := m.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet,
		&metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, err: %s, rid: %s", err.Error(),
			kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !rsp.Result {
		blog.Errorf("[operation-inst] failed to read the object(%s) inst by the condition(%#v), err: %s, rid: %s",
			common.BKInnerObjIDSet, cond, rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}
	if rsp.Data.Count > 0 {
		return nil
	}

	return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField+"/"+common.BKSetIDField)
}

// IsModuleNameDuplicateError check module name is or not exist
func (m *module) IsModuleNameDuplicateError(kit *rest.Kit, bizID, setID int64, moduleName string,
	inputErr error) (bool, error) {

	ccErr, ok := inputErr.(errors.CCErrorCoder)
	if ok == false {
		return false, nil
	}
	if ccErr.GetCode() != common.CCErrCommDuplicateItem {
		return false, nil
	}

	// 检测模块名重复并返回定制提示信息
	nameDuplicateFilter := &metadata.QueryCondition{
		Page: metadata.BasePage{
			Limit: 1,
		},
		Condition: map[string]interface{}{
			common.BKParentIDField:   setID,
			common.BKAppIDField:      bizID,
			common.BKModuleNameField: moduleName,
		},
	}
	result, err := m.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule, nameDuplicateFilter)
	if err != nil {
		blog.ErrorJSON("IsModuleNameDuplicateError failed, filter: %s, err: %s, rid: %s", nameDuplicateFilter, err.Error(), kit.Rid)
		return false, err
	}
	if ccErr := result.CCError(); ccErr != nil {
		blog.ErrorJSON("IsModuleNameDuplicateError failed, result false, filter: %s, result: %s, err: %s, rid: %s",
			nameDuplicateFilter, result, ccErr, kit.Rid)
		return false, ccErr
	}
	if result.Data.Count > 0 {
		return true, nil
	}
	return false, nil
}

// CreateModule create a new module
func (m *module) CreateModule(kit *rest.Kit, obj metadata.Object, bizID, setID int64,
	data mapstr.MapStr) (*metadata.CreateOneDataResult, error) {

	data.Set(common.BKSetIDField, setID)
	data.Set(common.BKAppIDField, bizID)
	if !data.Exists(common.BKDefaultField) {
		data.Set(common.BKDefaultField, common.DefaultFlagDefaultValue)
	}
	defaultVal, err := data.Int64(common.BKDefaultField)
	if err != nil {
		blog.Errorf("parse default field into int failed, data: %+v, rid: %s", data, kit.Rid)
		err := kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKDefaultField)
		return nil, err
	}

	if err := m.validBizSetID(kit, bizID, setID); err != nil {
		return nil, err
	}

	// validate service category id and service template id
	// 如果服务分类没有设置，则从服务模版中获取，如果服务模版也没有设置，则参数错误
	// 有效参数参数形式:
	// 1. serviceCategoryID > 0  && serviceTemplateID == 0
	// 2. serviceCategoryID unset && serviceTemplateID > 0
	// 3. serviceCategoryID > 0 && serviceTemplateID > 0 && serviceTemplate.ServiceCategoryID == serviceCategoryID
	// 4. serviceCategoryID unset && serviceTemplateID unset, then module create with default category
	var serviceCategoryID int64
	serviceCategoryIDIf, serviceCategoryExist := data.Get(common.BKServiceCategoryIDField)
	if serviceCategoryExist == true {
		scID, err := util.GetInt64ByInterface(serviceCategoryIDIf)
		if err != nil {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
		}
		serviceCategoryID = scID
	}

	var serviceTemplateID int64
	serviceTemplateIDIf, serviceTemplateFieldExist := data.Get(common.BKServiceTemplateIDField)
	if serviceTemplateFieldExist == true {
		serviceTemplateID, err = util.GetInt64ByInterface(serviceTemplateIDIf)
		if err != nil {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
	}

	// if need create module using service template
	if serviceTemplateID == 0 && !version.CanCreateSetModuleWithoutTemplate && defaultVal == 0 {
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, "service_template_id can not be 0")
	}

	if serviceCategoryID == 0 && serviceTemplateID == 0 {
		// set default service template id
		defaultServiceCategory, err := m.clientSet.CoreService().Process().GetDefaultServiceCategory(kit.Ctx, kit.Header)
		if err != nil {
			blog.Errorf("create module failed, GetDefaultServiceCategory failed, err: %s, rid: %s", err.Error(),
				kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrProcGetDefaultServiceCategoryFailed)
		}
		serviceCategoryID = defaultServiceCategory.ID
	} else if serviceTemplateID != common.ServiceTemplateIDNotSet {
		// 校验 serviceCategoryID 与 serviceTemplateID 对应
		templateIDs := []int64{serviceTemplateID}
		option := metadata.ListServiceTemplateOption{
			BusinessID:         bizID,
			ServiceTemplateIDs: templateIDs,
		}
		stResult, err := m.clientSet.CoreService().Process().ListServiceTemplates(kit.Ctx, kit.Header, &option)
		if err != nil {
			return nil, err
		}
		if len(stResult.Info) == 0 {
			blog.ErrorJSON("create module failed, service template not found, filter: %s, rid: %s", option,
				kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
		if serviceCategoryExist == true && serviceCategoryID != stResult.Info[0].ServiceCategoryID {
			return nil, kit.CCError.Error(common.CCErrProcServiceTemplateAndCategoryNotCoincide)
		}
		serviceCategoryID = stResult.Info[0].ServiceCategoryID
	} else {
		// 检查 service category id 是否有效
		serviceCategory, err := m.clientSet.CoreService().Process().GetServiceCategory(kit.Ctx, kit.Header,
			serviceCategoryID)
		if err != nil {
			return nil, err
		}
		if serviceCategory.BizID != 0 && serviceCategory.BizID != bizID {
			blog.V(3).Info("create module failed, service category and module belong to two business, "+
				"categoryBizID: %d, bizID: %d, rid: %s", serviceCategory.BizID, bizID, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
		}
	}
	data.Set(common.BKServiceCategoryIDField, serviceCategoryID)
	data.Set(common.BKServiceTemplateIDField, serviceTemplateID)
	data.Set(common.HostApplyEnabledField, false)

	// set default set template
	_, exist := data[common.BKSetTemplateIDField]
	if exist == false {
		data[common.BKSetTemplateIDField] = common.SetTemplateIDNotSet
	}

	// convert bk_parent_id to int
	parentIDIf, ok := data[common.BKParentIDField]
	if ok == true {
		parentID, err := util.GetInt64ByInterface(parentIDIf)
		if err != nil {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKParentIDField)
		}
		data[common.BKParentIDField] = parentID
	}

	data.Remove(common.MetadataField)
	inst, createErr := m.CreateInst(kit, obj, data)
	if createErr != nil {
		moduleNameStr, exist := data[common.BKModuleNameField]
		if exist == false {
			return inst, err
		}
		moduleName := util.GetStrByInterface(moduleNameStr)
		isDuplicate, err := m.IsModuleNameDuplicateError(kit, bizID, setID, moduleName, createErr)
		if err != nil {
			blog.Infof("create module failed and check whether is name duplicated err failed, bizID: %d,"+
				"setID: %d , moduleName: %s, err: %+v, rid: %s", bizID, setID, moduleName, err, kit.Rid)
			return inst, err
		}
		if isDuplicate {
			return inst, kit.CCError.CCError(common.CCErrorTopoModuleNameDuplicated)
		}
		return inst, createErr
	}

	return inst, nil
}

// hasHost check module is or not has host
func (m *module) hasHost(kit *rest.Kit, bizID int64, setIDs, moduleIDS []int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		ModuleIDArr:   moduleIDS,
		Fields:        []string{common.BKHostIDField},
		Page:          metadata.BasePage{Limit: 1},
	}
	if len(setIDs) > 0 {
		option.SetIDArr = setIDs
	}
	if len(moduleIDS) > 0 {
		option.ModuleIDArr = moduleIDS
	}
	rsp, err := m.clientSet.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, option)
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

// DeleteModule delete module
func (m *module) DeleteModule(kit *rest.Kit, bizID int64, setIDs, moduleIDs []int64) error {

	exists, err := m.hasHost(kit, bizID, setIDs, moduleIDs)
	if nil != err {
		blog.Errorf("[operation-module] failed to delete the modules, err: %s, rid: %s", err.Error(), kit.Rid)
		return err
	}

	if exists {
		blog.Errorf("[operation-module]the module has some hosts, can not be deleted, rid: %s", kit.Rid)
		return kit.CCError.Error(common.CCErrTopoHasHost)
	}

	innerCond := map[string]interface{}{common.BKAppIDField: bizID}
	if nil != setIDs {
		innerCond[common.BKSetIDField] = map[string]interface{}{common.BKDBIN: setIDs}
	}

	if nil != moduleIDs {
		innerCond[common.BKModuleIDField] = map[string]interface{}{common.BKDBIN: moduleIDs}
	}

	// module table doesn't have metadata field
	err = m.DeleteInst(kit, common.BKInnerObjIDModule, innerCond, false)
	if err != nil {
		blog.Errorf("delete module failed, DeleteInst failed, err: %+v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

// FindModule find module by obj id and other condition
func (m *module) FindModule(kit *rest.Kit, objID string, cond *metadata.QueryInput) (count int, results []mapstr.MapStr,
	err error) {
	resultData, err := m.FindInst(kit, objID, cond)
	if err != nil {
		return 0, nil, err
	}
	moduleInstances := make([]mapstr.MapStr, 0)
	for _, item := range resultData.Info {
		moduleInstance := make(map[string]interface{})
		if err := mapstr.DecodeFromMapStr(&moduleInstance, item); err != nil {
			blog.Errorf("unmarshal module into struct failed, module: %+v, rid: %s", item, kit.Rid)
			return 0, nil, err
		}
		moduleInstances = append(moduleInstances, moduleInstance)
	}
	return count, moduleInstances, err
}

// UpdateModule update module
func (m *module) UpdateModule(kit *rest.Kit, data mapstr.MapStr, obj *metadata.Object, bizID, setID,
	moduleID int64) error {
	innerCond := condition.CreateCondition()

	innerCond.Field(common.BKAppIDField).Eq(bizID)
	innerCond.Field(common.BKSetIDField).Eq(setID)
	innerCond.Field(common.BKModuleIDField).Eq(moduleID)

	findCond := &metadata.QueryInput{
		Condition: innerCond.ToMapStr(),
	}
	var err error
	count, moduleInstances, err := m.FindModule(kit, obj.ObjectID, findCond)
	if err != nil {
		blog.Errorf("update module failed, find module failed, filter: %+v, err: %s, rid: %s", findCond,
			err.Error(), kit.Rid)
		return err
	}
	if count == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommNotFound)
	}
	if count > 1 {
		return kit.CCError.CCErrorf(common.CCErrCommGetMultipleObject)
	}
	if len(moduleInstances) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommNotFound)
	}

	moduleMapStr := moduleInstances[0]
	moduleInstance := metadata.ModuleInst{}
	if err := mapstruct.Decode2Struct(moduleMapStr, &moduleInstance); err != nil {
		blog.ErrorJSON("unmarshal db data into module failed, module: %s, err: %s, rid: %s", moduleMapStr,
			err.Error(), kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParseDBFailed)
	}

	// 检查并提示禁止修改集群模板ID字段
	if val, ok := data[common.BKSetTemplateIDField]; ok == true {
		setTemplateID, err := util.GetInt64ByInterface(val)
		if err != nil {
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField)
		}
		if setTemplateID != moduleInstance.SetTemplateID {
			return kit.CCError.CCErrorf(common.CCErrCommModifyFieldForbidden, common.BKSetTemplateIDField)
		}
	}

	// 检查并提示禁止修改集服务模板ID字段
	if val, ok := data[common.BKServiceTemplateIDField]; ok == true {
		serviceTemplateID, err := util.GetInt64ByInterface(val)
		if err != nil {
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
		if serviceTemplateID != moduleInstance.ServiceTemplateID {
			return kit.CCError.CCErrorf(common.CCErrCommModifyFieldForbidden, common.BKServiceTemplateIDField)
		}
	}

	if moduleInstance.ServiceTemplateID != common.ServiceTemplateIDNotSet {
		// 检查并提示禁止修改服务分类
		if val, ok := data[common.BKServiceCategoryIDField]; ok == true {
			serviceCategoryID, err := util.GetInt64ByInterface(val)
			if err != nil {
				return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
			}
			if serviceCategoryID != moduleInstance.ServiceCategoryID {
				return kit.CCError.CCError(common.CCErrorTopoUpdateModuleFromTplServiceCategoryForbidden)
			}
		}

		// 检查并提示禁止修改通过模板创建的模块名称
		if val, ok := data[common.BKModuleNameField]; ok == true {
			name := util.GetStrByInterface(val)
			if len(name) == 0 {
				delete(data, common.BKModuleNameField)
			} else if name != moduleInstance.ModuleName {
				return kit.CCError.CCError(common.CCErrorTopoUpdateModuleFromTplNameForbidden)
			}
		}
	}

	data.Remove(common.BKAppIDField)
	data.Remove(common.BKSetIDField)
	data.Remove(common.BKModuleIDField)
	data.Remove(common.BKParentIDField)
	data.Remove(common.MetadataField)
	updateErr := m.UpdateInst(kit, innerCond.ToMapStr(), data, obj.ObjectID)
	if updateErr != nil {
		moduleNameStr, exist := data[common.BKModuleNameField]
		if exist == false {
			return updateErr
		}
		moduleName := util.GetStrByInterface(moduleNameStr)
		isDuplicate, err := m.IsModuleNameDuplicateError(kit, bizID, setID, moduleName, updateErr)
		if err != nil {
			blog.Infof("update module failed and check whether is name duplicated err failed, bizID: %d, "+
				"setID: %d,  moduleName: %s, err: %+v, rid: %s", bizID, setID, moduleName, err, kit.Rid)
			return err
		}
		if isDuplicate {
			return kit.CCError.CCError(common.CCErrorTopoModuleNameDuplicated)
		}
		return updateErr
	}

	return nil
}
