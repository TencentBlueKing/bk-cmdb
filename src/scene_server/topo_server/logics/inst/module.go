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
	"fmt"
	"strings"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/version"
)

// ModuleOperationInterface module operation methods
type ModuleOperationInterface interface {
	CreateModule(kit *rest.Kit, bizID, setID int64, data mapstr.MapStr) (*metadata.CreateOneDataResult, error)
	DeleteModule(kit *rest.Kit, bizID int64, setID, moduleIDS []int64) error
	UpdateModule(kit *rest.Kit, data mapstr.MapStr, bizID, setID, moduleID int64) error
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
func (m *module) CreateInst(kit *rest.Kit, obj metadata.Object, data mapstr.MapStr) (*metadata.CreateOneDataResult,
	error) {

	if obj.ObjectID == common.BKInnerObjIDPlat {
		data.Set(common.BkSupplierAccount, kit.SupplierAccount)
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

// IsModuleNameDuplicateError check module name is or not exist
func (m *module) IsModuleNameDuplicateError(kit *rest.Kit, bizID, setID int64, moduleName string) error {
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
	result, err := m.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		nameDuplicateFilter)
	if err != nil {
		blog.Errorf("module name duplicate err, filter: %s, err: %s, rid: %s", nameDuplicateFilter, err, kit.Rid)
		return err
	}

	if result.Data.Count > 0 {
		return kit.CCError.CCError(common.CCErrorTopoModuleNameDuplicated)
	}
	return nil
}

// CreateModule create a new module
func (m *module) CreateModule(kit *rest.Kit, bizID, setID int64, data mapstr.MapStr) (*metadata.CreateOneDataResult,
	error) {

	data.Set(common.BKSetIDField, setID)
	data.Set(common.BKAppIDField, bizID)
	if !data.Exists(common.BKDefaultField) {
		data.Set(common.BKDefaultField, common.DefaultFlagDefaultValue)
	}

	defaultVal, err := data.Int64(common.BKDefaultField)
	if err != nil {
		blog.Errorf("parse default field into int failed, data: %+v, rid: %s", data, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKDefaultField)
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
			blog.Errorf("get default service category failed, err: %s, rid: %s", err, kit.Rid)
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
			blog.Errorf("get service template not found, filter: %s, rid: %s", option, kit.Rid)
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
			blog.V(3).Info("get service category and module belong to two business, categoryBizID: %d, "+
				"bizID: %d, rid: %s", serviceCategory.BizID, bizID, kit.Rid)
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
	if ok {
		parentID, err := util.GetInt64ByInterface(parentIDIf)
		if err != nil {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKParentIDField)
		}
		if parentID != setID {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKParentIDField)
		}
		data[common.BKParentIDField] = parentID
	}
	data.Remove(common.MetadataField)

	// TODO 替换依赖CreateInst,obj参数从何而来？如果CreateModule()去掉obj入参？？？
	obj := metadata.Object{}
	inst, createErr := m.CreateInst(kit, obj, data)
	if createErr != nil {
		moduleNameStr, exist := data[common.BKModuleNameField]
		if exist == false {
			return inst, fmt.Errorf("get module name failed, moduleNameStr: %s", moduleNameStr)
		}
		moduleName := util.GetStrByInterface(moduleNameStr)
		if err := m.IsModuleNameDuplicateError(kit, bizID, setID, moduleName); err != nil {
			blog.Infof("create module failed and check whether is name duplicated err failed, bizID: %d,"+
				"setID: %d , moduleName: %s, err: %+v, rid: %s", bizID, setID, moduleName, err, kit.Rid)
			return inst, err
		}

		return inst, createErr
	}

	return inst, nil
}

// DeleteModule delete module
func (m *module) DeleteModule(kit *rest.Kit, bizID int64, setIDs, moduleIDs []int64) error {
	innerCond := map[string]interface{}{common.BKAppIDField: bizID}

	if len(setIDs) > 0 {
		innerCond[common.BKSetIDField] = map[string]interface{}{common.BKDBIN: setIDs}
	}

	if len(moduleIDs) > 0 {
		innerCond[common.BKModuleIDField] = map[string]interface{}{common.BKDBIN: moduleIDs}
	}

	// TODO 替换依赖DeleteInst
	// module table doesn't have metadata field
	err := m.DeleteInst(kit, common.BKInnerObjIDModule, innerCond, false)
	if err != nil {
		blog.Errorf("delete module failed, DeleteInst failed, err: %+v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// UpdateModule update module
func (m *module) UpdateModule(kit *rest.Kit, data mapstr.MapStr, bizID, setID, moduleID int64) error {
	objIDStr, ok := data.Get(common.BKObjIDField)
	if !ok {
		return kit.CCError.CCErrorf(common.CCErrCommNotFound)
	}
	objID := util.GetStrByInterface(objIDStr)

	innerCond := mapstr.MapStr{
		common.BKAppIDField:    bizID,
		common.BKSetIDField:    setID,
		common.BKModuleIDField: moduleID,
	}

	findCond := &metadata.QueryCondition{
		Fields: []string{common.BKSetTemplateIDField, common.BKServiceTemplateIDField, common.BKServiceCategoryIDField,
			common.BKModuleNameField},
		Condition:      innerCond,
		DisableCounter: true,
	}

	moduleInstance := new(metadata.ModuleInst)
	if err := m.clientSet.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, objID, findCond,
		moduleInstance); err != nil {
		blog.Errorf("list modules failed, bizID: %s, setID: %s, moduleID: %s, err: %s, rid: %s", bizID, setID,
			moduleID, err, kit.Rid)
		return err
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
	// TODO 后续替换依赖UpdateInst
	updateErr := m.UpdateInst(kit, innerCond, data, objID)
	if updateErr != nil {
		moduleNameStr, exist := data[common.BKModuleNameField]
		if exist == false {
			return updateErr
		}
		moduleName := util.GetStrByInterface(moduleNameStr)
		if err := m.IsModuleNameDuplicateError(kit, bizID, setID, moduleName); err != nil {
			blog.Infof("update module failed and check whether is name duplicated err failed, bizID: %d, "+
				"setID: %d,  moduleName: %s, err: %+v, rid: %s", bizID, setID, moduleName, err, kit.Rid)
			return err
		}

		return updateErr
	}

	return nil
}
