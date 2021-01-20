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
	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
)

// ModuleOperationInterface module operation methods
type ModuleOperationInterface interface {
	CreateModule(kit *rest.Kit, obj model.Object, bizID, setID int64, data mapstr.MapStr) (inst.Inst, error)
	DeleteModule(kit *rest.Kit, bizID int64, setID, moduleIDS []int64) error
	FindModule(kit *rest.Kit, obj model.Object, cond *metadata.QueryInput) (count int, results []mapstr.MapStr, err error)
	UpdateModule(kit *rest.Kit, data mapstr.MapStr, obj model.Object, bizID, setID, moduleID int64) error

	SetProxy(inst InstOperationInterface)
}

// NewModuleOperation create a new module
func NewModuleOperation(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) ModuleOperationInterface {
	return &module{
		clientSet:   client,
		authManager: authManager,
	}
}

type module struct {
	clientSet   apimachinery.ClientSetInterface
	inst        InstOperationInterface
	authManager *extensions.AuthManager
}

func (m *module) SetProxy(inst InstOperationInterface) {
	m.inst = inst
}

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

func (m *module) validBizSetID(kit *rest.Kit, bizID int64, setID int64) error {
	cond := condition.CreateCondition()
	cond.Field(common.BKSetIDField).Eq(setID)
	or := cond.NewOR()
	or.Item(mapstr.MapStr{common.BKAppIDField: bizID})

	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit

	rsp, err := m.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, &metadata.QueryCondition{Condition: cond.ToMapStr()})
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

func (m *module) CreateModule(kit *rest.Kit, obj model.Object, bizID, setID int64, data mapstr.MapStr) (inst.Inst, error) {

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
			blog.Errorf("create module failed, GetDefaultServiceCategory failed, err: %s, rid: %s", err.Error(), kit.Rid)
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
			blog.ErrorJSON("create module failed, service template not found, filter: %s, rid: %s", option, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
		if serviceCategoryExist == true && serviceCategoryID != stResult.Info[0].ServiceCategoryID {
			return nil, kit.CCError.Error(common.CCErrProcServiceTemplateAndCategoryNotCoincide)
		}
		serviceCategoryID = stResult.Info[0].ServiceCategoryID
	} else {
		// 检查 service category id 是否有效
		serviceCategory, err := m.clientSet.CoreService().Process().GetServiceCategory(kit.Ctx, kit.Header, serviceCategoryID)
		if err != nil {
			return nil, err
		}
		if serviceCategory.BizID != 0 && serviceCategory.BizID != bizID {
			blog.V(3).Info("create module failed, service category and module belong to two business, categoryBizID: %d, bizID: %d, rid: %s", serviceCategory.BizID, bizID, kit.Rid)
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
	inst, createErr := m.inst.CreateInst(kit, obj, data)
	if createErr != nil {
		moduleNameStr, exist := data[common.BKModuleNameField]
		if exist == false {
			return inst, err
		}
		moduleName := util.GetStrByInterface(moduleNameStr)
		isDuplicate, err := m.IsModuleNameDuplicateError(kit, bizID, setID, moduleName, createErr)
		if err != nil {
			blog.Infof("create module failed and check whether is name duplicated err failed, bizID: %d, setID: %d, moduleName: %s, err: %+v, rid: %s", bizID, setID, moduleName, err, kit.Rid)
			return inst, err
		}
		if isDuplicate {
			return inst, kit.CCError.CCError(common.CCErrorTopoModuleNameDuplicated)
		}
		return inst, createErr
	}

	return inst, nil
}

func (m *module) IsModuleNameDuplicateError(kit *rest.Kit, bizID, setID int64, moduleName string, inputErr error) (bool, error) {

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
	if result.Result == false || result.Code != 0 {
		blog.ErrorJSON("IsModuleNameDuplicateError failed, result false, filter: %s, result: %s, err: %s, rid: %s", nameDuplicateFilter, result, err.Error(), kit.Rid)
		return false, errors.New(result.Code, result.ErrMsg)
	}
	if result.Data.Count > 0 {
		return true, nil
	}
	return false, nil
}

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
	err = m.inst.DeleteInst(kit, common.BKInnerObjIDModule, innerCond, false)
	if err != nil {
		blog.Errorf("delete module failed, DeleteInst failed, err: %+v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func (m *module) FindModule(kit *rest.Kit, obj model.Object, cond *metadata.QueryInput) (count int, results []mapstr.MapStr, err error) {
	count, resultData, err := m.inst.FindInst(kit, obj, cond, false)
	if err != nil {
		return 0, nil, err
	}
	moduleInstances := make([]mapstr.MapStr, 0)
	for _, item := range resultData {
		moduleInstance := make(map[string]interface{})
		if err := mapstr.DecodeFromMapStr(&moduleInstance, item.ToMapStr()); err != nil {
			blog.Errorf("unmarshal module into struct failed, module: %+v, rid: %s", item, kit.Rid)
			return 0, nil, err
		}
		moduleInstances = append(moduleInstances, moduleInstance)
	}
	return count, moduleInstances, err
}

func (m *module) UpdateModule(kit *rest.Kit, data mapstr.MapStr, obj model.Object, bizID, setID, moduleID int64) error {
	innerCond := condition.CreateCondition()

	innerCond.Field(common.BKAppIDField).Eq(bizID)
	innerCond.Field(common.BKSetIDField).Eq(setID)
	innerCond.Field(common.BKModuleIDField).Eq(moduleID)

	findCond := &metadata.QueryInput{
		Condition: innerCond.ToMapStr(),
	}
	var err error
	count, moduleInstances, err := m.FindModule(kit, obj, findCond)
	if err != nil {
		blog.Errorf("update module failed, find module failed, filter: %+v, err: %s, rid: %s", findCond, err.Error(), kit.Rid)
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
		blog.ErrorJSON("unmarshal db data into module failed, module: %s, err: %s, rid: %s", moduleMapStr, err.Error(), kit.Rid)
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
	updateErr := m.inst.UpdateInst(kit, data, obj, innerCond, -1)
	if updateErr != nil {
		moduleNameStr, exist := data[common.BKModuleNameField]
		if exist == false {
			return updateErr
		}
		moduleName := util.GetStrByInterface(moduleNameStr)
		isDuplicate, err := m.IsModuleNameDuplicateError(kit, bizID, setID, moduleName, updateErr)
		if err != nil {
			blog.Infof("update module failed and check whether is name duplicated err failed, bizID: %d, setID: %d, moduleName: %s, err: %+v, rid: %s", bizID, setID, moduleName, err, kit.Rid)
			return err
		}
		if isDuplicate {
			return kit.CCError.CCError(common.CCErrorTopoModuleNameDuplicated)
		}
		return updateErr
	}

	return nil
}
