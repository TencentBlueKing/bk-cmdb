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
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/version"
)

// ModuleOperationInterface module operation methods
type ModuleOperationInterface interface {
	CreateModule(kit *rest.Kit, bizID, setID int64, data mapstr.MapStr) (*mapstr.MapStr, error)
	DeleteModule(kit *rest.Kit, bizID int64, setID, moduleIDS []int64) error
	UpdateModule(kit *rest.Kit, data mapstr.MapStr, bizID, setID, moduleID int64) error
	SetProxy(inst InstOperationInterface)
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
	inst        InstOperationInterface
}

// SetProxy 初始化依赖
func (m *module) SetProxy(inst InstOperationInterface) {
	m.inst = inst
}

// CreateModule create a new module
func (m *module) CreateModule(kit *rest.Kit, bizID, setID int64, data mapstr.MapStr) (*mapstr.MapStr, error) {
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

	if err := m.checkServiceTemplateParam(kit, serviceCategoryID, serviceTemplateID, bizID,
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

	inst, createErr := m.inst.CreateInst(kit, common.BKInnerObjIDModule, data)
	if createErr != nil {
		blog.Errorf("create module failed, err: %s, rid: %s", createErr, kit.Rid)
		return inst, createErr
	}

	return inst, nil
}

func (m *module) checkServiceTemplateParam(kit *rest.Kit, serviceCategoryID, serviceTemplateID, bizID int64,
	serviceCategoryExist bool) error{
	if serviceCategoryID == 0 && serviceTemplateID == 0 {
		// set default service template id
		defaultServiceCategory, err := m.clientSet.CoreService().Process().GetDefaultServiceCategory(kit.Ctx, kit.Header)
		if err != nil {
			blog.Errorf("get default service category failed, err: %s, rid: %s", err, kit.Rid)
			return kit.CCError.Errorf(common.CCErrProcGetDefaultServiceCategoryFailed)
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
		serviceCategory, err := m.clientSet.CoreService().Process().GetServiceCategory(kit.Ctx, kit.Header,
			serviceCategoryID)
		if err != nil {
			return err
		}
		if serviceCategory.BizID != 0 && serviceCategory.BizID != bizID {
			blog.V(3).Info("get service category and module belong to two business, categoryBizID: %d, "+
				"bizID: %d, rid: %s", serviceCategory.BizID, bizID, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
		}
	}
	return nil
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

	// module table doesn't have metadata field
	err := m.inst.DeleteInst(kit, common.BKInnerObjIDModule, innerCond, true)
	if err != nil {
		blog.Errorf("delete module failed, err: %s, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// UpdateModule update module
func (m *module) UpdateModule(kit *rest.Kit, data mapstr.MapStr, bizID, setID, moduleID int64) error {
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

	moduleInstance := new(metadata.ResponseModuleInstance)
	if err := m.clientSet.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		findCond, moduleInstance); err != nil {
		blog.Errorf("get list modules failed, bizID: %d, setID: %d, moduleID: %d, err: %s, rid: %s", bizID, setID,
			moduleID, err, kit.Rid)
		return err
	}

	// 检查并提示禁止修改集群模板ID字段
	if val, ok := data[common.BKSetTemplateIDField]; ok {
		setTemplateID, err := util.GetInt64ByInterface(val)
		if err != nil {
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField)
		}
		if setTemplateID != moduleInstance.Data.Info[0].SetTemplateID {
			return kit.CCError.CCErrorf(common.CCErrCommModifyFieldForbidden, common.BKSetTemplateIDField)
		}
	}

	// 检查并提示禁止修改服务模板ID字段
	if val, ok := data[common.BKServiceTemplateIDField]; ok {
		serviceTemplateID, err := util.GetInt64ByInterface(val)
		if err != nil {
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
		if serviceTemplateID != moduleInstance.Data.Info[0].ServiceTemplateID {
			return kit.CCError.CCErrorf(common.CCErrCommModifyFieldForbidden, common.BKServiceTemplateIDField)
		}
	}

	if moduleInstance.Data.Info[0].ServiceTemplateID != common.ServiceTemplateIDNotSet {
		// 检查并提示禁止修改服务分类
		if val, ok := data[common.BKServiceCategoryIDField]; ok {
			serviceCategoryID, err := util.GetInt64ByInterface(val)
			if err != nil {
				return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
			}
			if serviceCategoryID != moduleInstance.Data.Info[0].ServiceCategoryID {
				return kit.CCError.CCError(common.CCErrorTopoUpdateModuleFromTplServiceCategoryForbidden)
			}
		}

		// 检查并提示禁止修改通过服务模板创建的模块名称
		if val, ok := data[common.BKModuleNameField]; ok {
			name := util.GetStrByInterface(val)
			if name == "" {
				delete(data, common.BKModuleNameField)
			} else if name != moduleInstance.Data.Info[0].ModuleName {
				return kit.CCError.CCError(common.CCErrorTopoUpdateModuleFromTplNameForbidden)
			}
		}
	}

	data.Remove(common.BKAppIDField)
	data.Remove(common.BKSetIDField)
	data.Remove(common.BKModuleIDField)
	data.Remove(common.BKParentIDField)
	data.Remove(common.MetadataField)

	updateErr := m.inst.UpdateInst(kit, innerCond, data, common.BKInnerObjIDModule)
	if updateErr != nil {
		blog.Errorf("update module failed,  err: %s, rid: %s", updateErr, kit.Rid)
		return updateErr
	}

	return nil
}
