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
	"configcenter/src/common/util"
	"configcenter/src/common/version"
)

// SetOperationInterface set operation methods
type SetOperationInterface interface {
	CreateSet(kit *rest.Kit, bizID int64, data mapstr.MapStr) (*metadata.CreateOneDataResult, error)
	DeleteSet(kit *rest.Kit, bizID int64, setIDS []int64) error
	UpdateSet(kit *rest.Kit, data mapstr.MapStr, bizID, setID int64) error
}

// NewSetOperation create a set instance
func NewSetOperation(client apimachinery.ClientSetInterface, languageIf language.CCLanguageIf) SetOperationInterface {
	return &set{
		clientSet: client,
		language:  languageIf,
	}
}

type set struct {
	clientSet apimachinery.ClientSetInterface
	language  language.CCLanguageIf
}

// validObject check object vaild
func (s *set) validObject(kit *rest.Kit, obj metadata.Object) (*metadata.Association, error) {

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
	asst, err := s.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
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

// isValidBizInstID check biz id and set id valid
func (s *set) isValidBizInstID(kit *rest.Kit, objID string, instID int64, bizID int64) error {

	cond := mapstr.MapStr{
		metadata.GetInstIDFieldByObjID(objID): instID,
	}

	if bizID != 0 {
		cond.Set(common.BKAppIDField, bizID)
	}

	if metadata.IsCommon(objID) {
		cond.Set(common.BKObjIDField, objID)
	}

	rsp, err := s.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID,
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

// validMainLineParentID check parent id mainline valid
func (s *set) validMainLineParentID(kit *rest.Kit, assoc *metadata.Association, data mapstr.MapStr) error {
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

	if err = s.isValidBizInstID(kit, assoc.AsstObjID, parentID, bizID); err != nil {
		blog.Errorf("[operation-inst]parent id %d is invalid, err: %s, rid: %s", parentID, err.Error(), kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKParentIDField)
	}
	return nil
}

// CreateInst create instance by object and create message
func (s *set) CreateInst(kit *rest.Kit, obj metadata.Object,
	data mapstr.MapStr) (*metadata.CreateOneDataResult, error) {

	if obj.ObjectID == common.BKInnerObjIDPlat {
		data.Set(common.BkSupplierAccount, kit.SupplierAccount)
	}

	assoc, err := s.validObject(kit, obj)
	if err != nil {
		blog.Errorf("valid object (%s) failed, err: %v, rid: %s", obj.ObjectID, err, kit.Rid)
		return nil, err
	}

	if assoc != nil {
		if err := s.validMainLineParentID(kit, assoc, data); err != nil {
			blog.Errorf("the mainline object(%s) parent id invalid, err: %v, rid: %s", obj.ObjectID, err, kit.Rid)
			return nil, err
		}
	}

	data.Set(common.BKObjIDField, obj.ObjectID)

	instCond := &metadata.CreateModelInstance{Data: data}
	rsp, err := s.clientSet.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, obj.ObjectID, instCond)
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
	audit := auditlog.NewInstanceAudit(s.clientSet.CoreService())
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

// FindSingleObject find single object
func (s *set) FindSingleObject(kit *rest.Kit, objectID string) (*metadata.Object, error) {

	cond := mapstr.MapStr{
		common.BKObjIDField: objectID,
	}

	objs, err := s.FindObject(kit, cond)
	if err != nil {
		blog.Errorf("get model failed, failed to get model by supplier account(%s) objects(%s), err: %s, rid: %s",
			kit.SupplierAccount, objectID, err.Error(), kit.Rid)
		return nil, err
	}

	if len(objs) == 0 {
		blog.Errorf("get model failed, get model by supplier account(%s) objects(%s) not found, result: %+v, "+
			"rid: %s", kit.SupplierAccount, objectID, objs, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed,
			kit.CCError.Error(common.CCErrCommNotFound).Error())
	}

	if len(objs) > 1 {
		blog.Errorf("get model failed, get model by supplier account(%s) objects(%s) get multiple, result: %+v, "+
			"rid: %s", kit.SupplierAccount, objectID, objs, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed,
			kit.CCError.Error(common.CCErrCommGetMultipleObject).Error())
	}

	return &objs[0], nil
}

// FindObject find object
func (s *set) FindObject(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.Object, error) {
	rsp, err := s.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond},
	)
	if err != nil {
		blog.Errorf("find object failed, cond: %+v, err: %s, rid: %s", cond, err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to search the objects by the condition(%#v) , err: %v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}

	return rsp.Data.Info, nil
}

// CreateModule create a new module
func (s *set) CreateModule(kit *rest.Kit, bizID, setID int64, data mapstr.MapStr) (*metadata.CreateOneDataResult,
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
		defaultServiceCategory, err := s.clientSet.CoreService().Process().GetDefaultServiceCategory(kit.Ctx,
			kit.Header)
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
		stResult, err := s.clientSet.CoreService().Process().ListServiceTemplates(kit.Ctx, kit.Header, &option)
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
		serviceCategory, err := s.clientSet.CoreService().Process().GetServiceCategory(kit.Ctx, kit.Header,
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
	inst, createErr := s.CreateInst(kit, obj, data)
	if createErr != nil {
		moduleNameStr, exist := data[common.BKModuleNameField]
		if exist == false {
			return inst, fmt.Errorf("get module name failed, moduleNameStr: %s", moduleNameStr)
		}
		moduleName := util.GetStrByInterface(moduleNameStr)
		if _, err := s.IsModuleNameDuplicateError(kit, bizID, setID, moduleName, createErr); err != nil {
			blog.Infof("create module failed and check whether is name duplicated err failed, bizID: %d,"+
				"setID: %d , moduleName: %s, err: %+v, rid: %s", bizID, setID, moduleName, err, kit.Rid)
			return inst, err
		}

		return inst, createErr
	}

	return inst, nil
}

// validBizSetID check biz and set id valid
func (s *set) validBizSetID(kit *rest.Kit, bizID int64, setID int64) error {
	cond := condition.CreateCondition()
	cond.Field(common.BKSetIDField).Eq(setID)
	or := cond.NewOR()
	or.Item(mapstr.MapStr{common.BKAppIDField: bizID})

	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit

	rsp, err := s.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet,
		&metadata.QueryCondition{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !rsp.Result {
		blog.Errorf("[operation-inst] failed to read the object(%s) inst by the condition(%#v), err: %s, rid: %s", common.BKInnerObjIDSet, cond, rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}
	if rsp.Data.Count > 0 {
		return nil
	}

	return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField+"/"+common.BKSetIDField)
}

// IsModuleNameDuplicateError check moduel name duplicate
func (s *set) IsModuleNameDuplicateError(kit *rest.Kit, bizID, setID int64, moduleName string, inputErr error) (bool,
	error) {

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
	result, err := s.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		nameDuplicateFilter)
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

// DeleteInst delete instance by objectid and condition
func (s *set) DeleteInst(kit *rest.Kit, objectID string, cond mapstr.MapStr, needCheckHost bool) error {
	return s.deleteInstByCond(kit, objectID, cond, needCheckHost)
}

// deleteInstByCond delete instance by condition
func (s *set) deleteInstByCond(kit *rest.Kit, objectID string, cond mapstr.MapStr, needCheckHost bool) error {
	query := &metadata.QueryInput{
		Condition: cond,
		Limit:     common.BKNoLimit,
	}

	instRsp, err := s.FindInst(kit, objectID, query)
	if err != nil {
		return err
	}

	if len(instRsp.Info) == 0 {
		return nil
	}

	delObjInstsMap, exists, err := s.hasHosts(kit, instRsp.Info, objectID, needCheckHost)
	if err != nil {
		return err
	}
	if exists {
		return kit.CCError.Error(common.CCErrTopoHasHostCheckFailed)
	}

	bizSetMap := make(map[int64][]int64)
	audit := auditlog.NewInstanceAudit(s.clientSet.CoreService())
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
		cnt, err := s.clientSet.CoreService().Association().CountInstanceAssociations(kit.Ctx, kit.Header, objID, input)
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
		rsp, err := s.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, objID, dc)
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
			if ccErr := s.clientSet.CoreService().SetTemplate().DeleteSetTemplateSyncStatus(kit.Ctx, kit.Header,
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
func (s *set) FindInst(kit *rest.Kit, objID string, cond *metadata.QueryInput) (*metadata.InstResult, error) {

	result := new(metadata.InstResult)
	switch objID {
	case common.BKInnerObjIDHost:
		rsp, err := s.clientSet.CoreService().Host().GetHosts(kit.Ctx, kit.Header, cond)
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
		rsp, err := s.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, input)
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
func (s *set) hasHosts(kit *rest.Kit, instances []mapstr.MapStr, objID string, checkHost bool) (
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

		moduleRsp, err := s.FindInst(kit, common.BKInnerObjIDModule, query)
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
		asstRsp, err := s.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, mainlineCond)
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

			childRsp, err := s.FindInst(kit, childObjID, query)
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
		exists, err := s.innerHasHost(kit, moduleIDs)
		if err != nil {
			return nil, false, err
		}

		if exists {
			return nil, true, nil
		}
	}

	return objInstMap, false, nil
}

// innerHasHost check host ip is or not inner ip
func (s *set) innerHasHost(kit *rest.Kit, moduleIDS []int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ModuleIDArr: moduleIDS,
		Fields:      []string{common.BKHostIDField},
		Page:        metadata.BasePage{Limit: 1},
	}
	rsp, err := s.clientSet.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, option)
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
func (s *set) UpdateInst(kit *rest.Kit, cond, data mapstr.MapStr, objID string) error {
	// not allowed to update these fields, need to use specialized function
	data.Remove(common.BKParentIDField)
	data.Remove(common.BKAppIDField)

	inputParams := metadata.UpdateOption{
		Data:      data,
		Condition: cond,
	}

	// generate audit log of instance.
	audit := auditlog.NewInstanceAudit(s.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(data)
	auditLog, ccErr := audit.GenerateAuditLogByCondGetData(generateAuditParameter, objID, cond)
	if ccErr != nil {
		blog.Errorf(" update inst, generate audit log failed, err: %v, rid: %s", ccErr, kit.Rid)
		return ccErr
	}

	// to update.
	rsp, err := s.clientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, objID, &inputParams)
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

// isSetDuplicateError check set exist
func (s *set) isSetDuplicateError(inputErr error) bool {
	ccErr, ok := inputErr.(errors.CCErrorCoder)
	if ok == false {
		return false
	}

	if ccErr.GetCode() == common.CCErrCommDuplicateItem {
		return true
	}

	return false
}

// hasHost check module has host
func (s *set) hasHost(kit *rest.Kit, bizID int64, setIDS []int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		SetIDArr:      setIDS,
	}
	rsp, err := s.clientSet.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, option)
	if nil != err {
		blog.Errorf("[operation-set] failed to request the object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return false, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-set]  failed to search the host set configures, error info is %s, rid: %s", rsp.ErrMsg, kit.Rid)
		return false, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != len(rsp.Data.Info), nil
}

// CreateSet create a new set
func (s *set) CreateSet(kit *rest.Kit, bizID int64, data mapstr.MapStr) (*metadata.CreateOneDataResult, error) {
	data.Set(common.BKAppIDField, bizID)

	if !data.Exists(common.BKDefaultField) {
		data.Set(common.BKDefaultField, common.DefaultFlagDefaultValue)
	}
	defaultVal, err := data.Int64(common.BKDefaultField)
	if err != nil {
		blog.Errorf("parse default field into int failed, data: %+v, rid: %s", data, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKDefaultField)
	}

	setTemplate := new(metadata.SetTemplate)
	// validate foreign key
	if setTemplateIDIf, ok := data[common.BKSetTemplateIDField]; ok == true {
		setTemplateID, err := util.GetInt64ByInterface(setTemplateIDIf)
		if err != nil {
			blog.Errorf("parse set_template_id field into int failed, id: %+v, rid: %s", setTemplateIDIf, kit.Rid)
			err := kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, s.language.CreateDefaultCCLanguageIf(util.
				GetLanguage(kit.Header)).Language("set_property_set_template_id"))
			return nil, err
		}
		if setTemplateID != common.SetTemplateIDNotSet {
			st, err := s.clientSet.CoreService().SetTemplate().GetSetTemplate(kit.Ctx, kit.Header, bizID, setTemplateID)
			if err != nil {
				return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, s.language.CreateDefaultCCLanguageIf(util.
					GetLanguage(kit.Header)).Language("set_property_set_template_id"))
			}
			setTemplate = &st
		}
	}

	// if need create set using set template
	if setTemplate.ID == common.SetTemplateIDNotSet && !version.CanCreateSetModuleWithoutTemplate && defaultVal == 0 {
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, "set_template_id can not be 0")
	}
	data.Set(common.BKSetTemplateIDField, setTemplate.ID)
	data.Remove(common.MetadataField)

	// TODO 替换CreateInst依赖 obj数据从哪来？？？
	obj := metadata.Object{}
	setInstance, err := s.CreateInst(kit, obj, data)

	if err != nil {
		blog.Errorf("create set instance failed, object: %s, data: %s, err: %s, rid: %s", obj, data, err, kit.Rid)
		// return this duplicate error for unique validation failed
		if s.isSetDuplicateError(err) {
			return setInstance, kit.CCError.CCError(common.CCErrorSetNameDuplicated)
		}
		return setInstance, err
	}
	if setTemplate.ID == common.SetTemplateIDNotSet {
		return setInstance, nil
	}

	// set create by template should create module at the same time
	serviceTemplates, err := s.clientSet.CoreService().SetTemplate().ListSetTplRelatedSvcTpl(kit.Ctx, kit.Header,
		bizID, setTemplate.ID)
	if err != nil {
		blog.Errorf("create set failed, list set tpl related svc tpl failed, bizID: %d, setTemplateID: %d, " +
			"err: %s, rid: %s", bizID, setTemplate.ID, err, kit.Rid)
		return setInstance, err
	}

	setID := setInstance.Created.ID
	for _, serviceTemplate := range serviceTemplates {
		createModuleParam := mapstr.MapStr{
			common.BKModuleNameField:        serviceTemplate.Name,
			common.BKServiceTemplateIDField: serviceTemplate.ID,
			common.BKSetTemplateIDField:     setTemplate.ID,
			common.BKParentIDField:          setID,
			common.BKServiceCategoryIDField: serviceTemplate.ServiceCategoryID,
			common.BKAppIDField:             bizID,
		}
		// TODO 替换依赖CreateModule,目前创建模块只能一个一个创建，在CreateModule改成批量创建也是通过遍历循环创建？有区别？
		if _, err := s.CreateModule(kit, bizID, int64(setID), createModuleParam); err != nil {
			blog.Errorf("create set instance failed, create module instance failed, bizID: %d, setID: %d, "+
				"param: %+v, err: %s, rid: %s", bizID, setID, createModuleParam, err.Error(), kit.Rid)
			return setInstance, err
		}
	}

	return setInstance, nil
}

// DeleteSet delete set
func (s *set) DeleteSet(kit *rest.Kit, bizID int64, setIDs []int64) error {
	setCond := map[string]interface{}{common.BKAppIDField: bizID}
	if len(setIDs) > 0 {
		setCond[common.BKSetIDField] = map[string]interface{}{common.BKDBIN: setIDs}
	}

	// TODO 替换依赖DeleteInst
	// clear the module belong to deleted sets
	err := s.DeleteInst(kit, common.BKInnerObjIDModule, setCond, false)
	if err != nil {
		blog.Errorf("delete module failed, err: %v, cond: %#v, rid: %s", err, setCond, kit.Rid)
		return err
	}

	// clear set template sync status
	if ccErr := s.clientSet.CoreService().SetTemplate().DeleteSetTemplateSyncStatus(kit.Ctx, kit.Header, bizID,
		setIDs); ccErr != nil {
		blog.Errorf("failed to delete set template sync status failed, bizID: %d, setIDs: %+v, "+
			"err: %s, rid: %s", bizID, setIDs, ccErr.Error(), kit.Rid)
		return ccErr
	}

	taskCond := &metadata.DeleteOption{
		Condition: map[string]interface{}{
			common.BKInstIDField: map[string]interface{}{
				common.BKDBIN: setIDs,
			}}}
	if err = s.clientSet.TaskServer().Task().DeleteTask(kit.Ctx, kit.Header, taskCond); err != nil {
		blog.Errorf("failed to delete set sync task message failed, bizID: %d, setIDs: %+v, err: %s, rid: %s",
			bizID, setIDs, err.Error(), kit.Rid)
		return err
	}

	// clear the sets
	return s.DeleteInst(kit, common.BKInnerObjIDSet, setCond, false)
}

// UpdateSet update set
func (s *set) UpdateSet(kit *rest.Kit, data mapstr.MapStr, bizID, setID int64) error {
	innerCond := mapstr.MapStr{
		common.BKAppIDField: bizID,
		common.BKSetIDField: setID,
	}

	data.Remove(common.MetadataField)
	data.Remove(common.BKAppIDField)
	data.Remove(common.BKSetIDField)
	data.Remove(common.BKSetTemplateIDField)
	// TODO 不知前端是否会传递objID参数
	err := s.UpdateInst(kit, innerCond, data, "")
	if err != nil {
		blog.Errorf("update set instance failed, object: %s, data: %s, innerCond:%s, err: %s, rid: %s", "",
			data, innerCond, err, kit.Rid)
		// return this duplicate error for unique validation failed
		if s.isSetDuplicateError(err) {
			return kit.CCError.CCError(common.CCErrorSetNameDuplicated)
		}
		return err
	}

	return nil
}
