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
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/topo_server/logics/model"
)

// SetOperationInterface set operation methods
type SetOperationInterface interface {
	CreateSet(kit *rest.Kit, bizID int64, data mapstr.MapStr) (*mapstr.MapStr, error)
	DeleteSet(kit *rest.Kit, bizID int64, setIDS []int64) error
	UpdateSet(kit *rest.Kit, data mapstr.MapStr, bizID, setID int64) error
	SetProxy(obj model.ObjectOperationInterface, inst InstOperationInterface)
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
	inst      InstOperationInterface
	obj       model.ObjectOperationInterface
	language  language.CCLanguageIf
}

func (s *set) SetProxy(obj model.ObjectOperationInterface, inst InstOperationInterface) {
	s.inst = inst
	s.obj = obj
}

// isSetDuplicateError check set exist
func (s *set) isSetDuplicateError(inputErr error) bool {
	ccErr, ok := inputErr.(errors.CCErrorCoder)
	if !ok {
		return false
	}

	if ccErr.GetCode() == common.CCErrCommDuplicateItem {
		return true
	}

	return false
}

// 依赖需要删除
func (s *set) validBizSetID(kit *rest.Kit, bizID int64, setID int64) error {
	cond := mapstr.MapStr{
		common.BKDBOR: []mapstr.MapStr{
			{common.BKSetIDField: setID},
			{common.BKAppIDField: bizID},
		},
	}

	query := &metadata.QueryCondition{
		Condition: cond,
	}

	rsp, err := s.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, query)
	if err != nil {
		blog.Errorf("get module instance failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if rsp.Count > 0 {
		return nil
	}

	return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField+"/"+common.BKSetIDField)
}

// CreateModule create a new module 依赖需要删除
func (s *set) CreateModule(kit *rest.Kit, bizID, setID int64, data mapstr.MapStr) (*mapstr.MapStr, error) {
	data.Set(common.BKSetIDField, setID)
	data.Set(common.BKAppIDField, bizID)
	if !data.Exists(common.BKDefaultField) {
		data.Set(common.BKDefaultField, common.DefaultFlagDefaultValue)
	}

	defaultVal, err := data.Int64(common.BKDefaultField)
	if err != nil {
		blog.Errorf("parse default field into int failed, data: %#v, rid: %s", data, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKDefaultField)
	}
	if err := s.validBizSetID(kit, bizID, setID); err != nil {
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
	if serviceCategoryExist {
		serviceCategoryID, err = util.GetInt64ByInterface(serviceCategoryIDIf)
		if err != nil {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
		}
	}

	var serviceTemplateID int64
	serviceTemplateIDIf, serviceTemplateFieldExist := data.Get(common.BKServiceTemplateIDField)
	if serviceTemplateFieldExist {
		serviceTemplateID, err = util.GetInt64ByInterface(serviceTemplateIDIf)
		if err != nil {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
	}

	// if need create module using service template
	if serviceTemplateID == 0 && !version.CanCreateSetModuleWithoutTemplate && defaultVal == 0 {
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, "service_template_id can not be 0")
	}

	if err := s.checkServiceTemplateParam(kit, serviceCategoryID, serviceTemplateID, bizID,
		serviceCategoryExist); err != nil {
		return nil, err
	}
	data.Set(common.BKServiceCategoryIDField, serviceCategoryID)
	data.Set(common.BKServiceTemplateIDField, serviceTemplateID)
	data.Set(common.HostApplyEnabledField, false)

	// set default set template
	_, exist := data[common.BKSetTemplateIDField]
	if !exist {
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

	inst, createErr := s.inst.CreateInst(kit, common.BKInnerObjIDModule, data)
	if createErr != nil {
		blog.Errorf("create module failed, err: %v, rid: %s", createErr, kit.Rid)
		return inst, createErr
	}

	return inst, nil
}

// 依赖需要删除
func (s *set) checkServiceTemplateParam(kit *rest.Kit, serviceCategoryID, serviceTemplateID, bizID int64,
	serviceCategoryExist bool) error {
	if serviceCategoryID == 0 && serviceTemplateID == 0 {
		// set default service template id
		defaultServiceCategory, err := s.clientSet.CoreService().Process().GetDefaultServiceCategory(kit.Ctx, kit.Header)
		if err != nil {
			blog.Errorf("get default service category failed, err: %v, rid: %s", err, kit.Rid)
			return err
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
			return err
		}
		if len(stResult.Info) == 0 {
			blog.Errorf("get service template not found, filter: %s, rid: %s", option, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
		if serviceCategoryExist == true && serviceCategoryID != stResult.Info[0].ServiceCategoryID {
			return kit.CCError.Error(common.CCErrProcServiceTemplateAndCategoryNotCoincide)
		}
		serviceCategoryID = stResult.Info[0].ServiceCategoryID
	} else {
		// 检查 service category id 是否有效
		serviceCategory, err := s.clientSet.CoreService().Process().GetServiceCategory(kit.Ctx, kit.Header,
			serviceCategoryID)
		if err != nil {
			return err
		}
		if serviceCategory.BizID != 0 && serviceCategory.BizID != bizID {
			blog.V(3).Info("service category and module belong to two business, categoryBizID: %d, bizID: %d, "+
				"rid: %s", serviceCategory.BizID, bizID, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
		}
	}
	return nil
}

// CreateSet create a new set
func (s *set) CreateSet(kit *rest.Kit, bizID int64, data mapstr.MapStr) (*mapstr.MapStr, error) {
	data.Set(common.BKAppIDField, bizID)

	if !data.Exists(common.BKDefaultField) {
		data.Set(common.BKDefaultField, common.DefaultFlagDefaultValue)
	}
	defaultVal, err := data.Int64(common.BKDefaultField)
	if err != nil {
		blog.Errorf("parse default field into int failed, data: %#v, rid: %s", data, kit.Rid)
		return nil, err
	}

	setTemplate := metadata.SetTemplate{}
	// validate foreign key
	if setTemplateIDIf, ok := data[common.BKSetTemplateIDField]; ok {
		setTemplateID, err := util.GetInt64ByInterface(setTemplateIDIf)
		if err != nil {
			blog.Errorf("parse set_template_id field into int failed, err: %v, rid: %s", err, kit.Rid)
			err := kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, s.language.CreateDefaultCCLanguageIf(util.
				GetLanguage(kit.Header)).Language("set_property_set_template_id"))
			return nil, err
		}
		if setTemplateID != common.SetTemplateIDNotSet {
			st, err := s.clientSet.CoreService().SetTemplate().GetSetTemplate(kit.Ctx, kit.Header, bizID, setTemplateID)
			if err != nil {
				return nil, err
			}
			setTemplate = st
		}
	}

	// if need create set using set template
	if setTemplate.ID == common.SetTemplateIDNotSet && !version.CanCreateSetModuleWithoutTemplate && defaultVal == 0 {
		blog.Errorf("service template not exist, can not create set, rid: %s", kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, "set_template_id can not be 0")
	}

	data.Set(common.BKSetTemplateIDField, setTemplate.ID)
	data.Remove(common.MetadataField)

	setInstance, err := s.inst.CreateInst(kit, common.BKInnerObjIDSet, data)
	if err != nil {
		blog.Errorf("create set instance failed, data: %#v, err: %v, rid: %s", data, err, kit.Rid)
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
		blog.Errorf("get list set tpl related svc tpl failed, bizID: %d, setTemplateID: %d, err: %v, rid: %s",
			bizID, setTemplate.ID, err, kit.Rid)
		return setInstance, err
	}

	setID, err := metadata.GetInstID(common.BKInnerObjIDSet, *setInstance)
	if err != nil {
		blog.Errorf("get inst id failed, err: %v, rid: %s", err, kit.Rid)
		return setInstance, err
	}

	for _, serviceTemplate := range serviceTemplates {
		createModuleParam := mapstr.MapStr{
			common.BKModuleNameField:        serviceTemplate.Name,
			common.BKServiceTemplateIDField: serviceTemplate.ID,
			common.BKSetTemplateIDField:     setTemplate.ID,
			common.BKParentIDField:          setID,
			common.BKServiceCategoryIDField: serviceTemplate.ServiceCategoryID,
			common.BKAppIDField:             bizID,
		}

		// TODO 替换依赖CreateModule
		if _, err := s.CreateModule(kit, bizID, setID, createModuleParam); err != nil {
			blog.Errorf("create module instance failed, bizID: %d, setID: %d, "+
				"param: %#v, err: %v, rid: %s", bizID, setID, createModuleParam, err, kit.Rid)
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

	// clear the module belong to deleted sets
	err := s.inst.DeleteInst(kit, common.BKInnerObjIDModule, setCond, true)
	if err != nil {
		blog.Errorf("delete module failed, err: %v, cond: %#v, rid: %s", err, setCond, kit.Rid)
		return err
	}

	// clear set template sync status
	if ccErr := s.clientSet.CoreService().SetTemplate().DeleteSetTemplateSyncStatus(kit.Ctx, kit.Header, bizID,
		setIDs); ccErr != nil {
		blog.Errorf("failed to delete set template sync status failed, bizID: %d, setIDs: %#v, "+
			"err: %v, rid: %s", bizID, setIDs, ccErr.Error(), kit.Rid)
		return ccErr
	}

	taskCond := &metadata.DeleteOption{
		Condition: map[string]interface{}{
			common.BKInstIDField: map[string]interface{}{
				common.BKDBIN: setIDs,
			},
		},
	}
	if err = s.clientSet.TaskServer().Task().DeleteTask(kit.Ctx, kit.Header, taskCond); err != nil {
		blog.Errorf("failed to delete set sync task message failed, bizID: %d, setIDs: %#v, err: %v, rid: %s",
			bizID, setIDs, err, kit.Rid)
		return err
	}

	// clear the sets
	return s.inst.DeleteInst(kit, common.BKInnerObjIDSet, setCond, true)
}

// UpdateSet update set
func (s *set) UpdateSet(kit *rest.Kit, data mapstr.MapStr, bizID, setID int64) error {
	innerCond := mapstr.MapStr{
		common.BKAppIDField: bizID,
		common.BKSetIDField: setID,
	}

	data.Remove(common.MetadataField)
	data.Remove(common.BKSetIDField)
	data.Remove(common.BKSetTemplateIDField)

	err := s.inst.UpdateInst(kit, innerCond, data, common.BKInnerObjIDSet)
	if err != nil {
		blog.Errorf("update set instance failed, data: %#v, innerCond:%#v, err: %v, rid: %s", data, innerCond, err,
			kit.Rid)
		// return this duplicate error for unique validation failed
		if s.isSetDuplicateError(err) {
			blog.Errorf("update set instance failed, set name duplicated, rid: %s", kit.Rid)
			return kit.CCError.CCError(common.CCErrorSetNameDuplicated)
		}
		return err
	}

	return nil
}
