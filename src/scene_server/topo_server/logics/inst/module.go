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
	"bytes"
	"reflect"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/version"
)

// ModuleOperationInterface module operation methods
type ModuleOperationInterface interface {
	CreateModule(kit *rest.Kit, bizID, setID int64, data mapstr.MapStr) (mapstr.MapStr, error)
	DeleteModule(kit *rest.Kit, bizID int64, setID, moduleIDS []int64) error
	UpdateModule(kit *rest.Kit, data mapstr.MapStr, bizID, setID, moduleID int64) error
	GetInternalModule(kit *rest.Kit, bizID int64) (count int, result *metadata.InnterAppTopo, err errors.CCErrorCoder)
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

// GetInternalModule 获取内置模型
func (m *module) GetInternalModule(kit *rest.Kit, bizID int64) (count int, result *metadata.InnterAppTopo,
	err errors.CCErrorCoder) {
	// get default set model
	querySet := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKAppIDField:   bizID,
			common.BKDefaultField: common.DefaultResSetFlag,
		},
		Fields: []string{common.BKSetIDField, common.BKSetNameField},
	}
	querySet.Page.Limit = 1

	setRsp := &metadata.ResponseSetInstance{}
	// 返回数据不包含自定义字段
	if err = m.clientSet.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDSet,
		querySet, setRsp); err != nil {
		return 0, nil, err
	}
	if err := setRsp.CCError(); err != nil {
		blog.Errorf("query set failed, err: %v, rid: %s", err, kit.Rid)
		return 0, nil, err
	}

	// search modules
	queryModule := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKDefaultField: map[string]interface{}{
				common.BKDBNE: 0,
			},
		},
		Fields: []string{common.BKModuleIDField, common.BKModuleNameField, common.BKDefaultField,
			common.HostApplyEnabledField},
	}
	queryModule.Page.Limit = common.BKNoLimit

	moduleResp := &metadata.ResponseModuleInstance{}
	// 返回数据不包含自定义字段
	if err = m.clientSet.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header,
		common.BKInnerObjIDModule, queryModule, moduleResp); err != nil {
		return 0, nil, err
	}
	if err := moduleResp.CCError(); err != nil {
		blog.Errorf("query module failed, err: %v, rid: %s", err, kit.Rid)
		return 0, nil, err
	}

	// construct result
	result = &metadata.InnterAppTopo{}
	for _, set := range setRsp.Data.Info {
		result.SetID = set.SetID
		result.SetName = set.SetName
		break // should be only one set
	}

	for _, module := range moduleResp.Data.Info {
		result.Module = append(result.Module, metadata.InnerModule{
			ModuleID:         module.ModuleID,
			ModuleName:       module.ModuleName,
			Default:          module.Default,
			HostApplyEnabled: module.HostApplyEnabled,
		})
	}

	return 0, result, nil
}

func (m *module) validBizSetID(kit *rest.Kit, bizID int64, setID int64) error {
	query := &metadata.Condition{
		Condition: mapstr.MapStr{
			common.BKSetIDField: setID,
			common.BKAppIDField: bizID,
		},
	}

	rsp, err := m.clientSet.CoreService().Instance().CountInstances(kit.Ctx, kit.Header, common.BKInnerObjIDSet, query)
	if err != nil {
		blog.Errorf("get module instance failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if int(rsp.Count) > 0 {
		return nil
	}

	return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField+"/"+common.BKSetIDField)
}

// CreateModule create a new module
func (m *module) CreateModule(kit *rest.Kit, bizID, setID int64, data mapstr.MapStr) (mapstr.MapStr, error) {
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
	if err := m.validBizSetID(kit, bizID, setID); err != nil {
		return nil, err
	}

	serviceTemplateID, serviceCategoryID, err := m.checkModuleServiceTemplate(kit, bizID, defaultVal, data)
	if err != nil {
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
			blog.Errorf("get module parent id failed, err: %v, rid: %s", err, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKParentIDField)
		}
		if parentID != setID {
			blog.Errorf("module parent id not equal set id, rid: %s", kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKParentIDField)
		}
		data[common.BKParentIDField] = parentID
	}
	data.Remove(common.MetadataField)

	// if service template has attributes, initialize module using these attributes
	data, err = m.initModuleWithSvcTemp(kit, bizID, serviceTemplateID, data)
	if err != nil {
		return nil, err
	}

	inst, createErr := m.inst.CreateInst(kit, common.BKInnerObjIDModule, data)
	if createErr != nil {
		blog.Errorf("create module failed, err: %v, rid: %s", createErr, kit.Rid)
		return inst, createErr
	}

	return inst, nil
}

// checkModuleServiceTemplate TODO
// validate service category id and service template id
// 如果服务分类没有设置，则从服务模版中获取，如果服务模版也没有设置，则参数错误
// 有效参数参数形式:
// 1. serviceCategoryID > 0  && serviceTemplateID == 0
// 2. serviceCategoryID unset && serviceTemplateID > 0
// 3. serviceCategoryID > 0 && serviceTemplateID > 0 && serviceTemplate.ServiceCategoryID == serviceCategoryID
// 4. serviceCategoryID unset && serviceTemplateID unset, then module create with default category
func (m *module) checkModuleServiceTemplate(kit *rest.Kit, bizID, defaultVal int64, data mapstr.MapStr) (int64,
	int64, error) {

	var err error
	var serviceCategoryID int64
	serviceCategoryIDIf, serviceCategoryExist := data.Get(common.BKServiceCategoryIDField)
	if serviceCategoryExist {
		serviceCategoryID, err = util.GetInt64ByInterface(serviceCategoryIDIf)
		if err != nil {
			blog.Errorf("get service category id failed, err: %v, rid: %s", err, kit.Rid)
			return 0, 0, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
		}
	}

	var serviceTemplateID int64
	serviceTemplateIDIf, serviceTemplateFieldExist := data.Get(common.BKServiceTemplateIDField)
	if serviceTemplateFieldExist {
		serviceTemplateID, err = util.GetInt64ByInterface(serviceTemplateIDIf)
		if err != nil {
			blog.Errorf("get service template id failed, err: %v, rid: %s", err, kit.Rid)
			return 0, 0, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
	}

	// if need create module using service template
	if serviceTemplateID == 0 && !version.CanCreateSetModuleWithoutTemplate && defaultVal == 0 {
		blog.Errorf("not use  service template create set module, rid: %s", kit.Rid)
		return 0, 0, kit.CCError.Errorf(common.CCErrCommParamsInvalid, "service_template_id can not be 0")
	}

	serviceCategoryID, err = m.checkServiceTemplateParam(kit, serviceCategoryID, serviceTemplateID, bizID,
		serviceCategoryExist)
	if err != nil {
		return 0, 0, err
	}

	return serviceTemplateID, serviceCategoryID, nil
}

func (m *module) checkServiceTemplateParam(kit *rest.Kit, serviceCategoryID, serviceTemplateID, bizID int64,
	serviceCategoryExist bool) (int64, error) {
	if serviceCategoryID == 0 && serviceTemplateID == 0 {
		// set default service template id
		defaultServiceCategory, err := m.clientSet.CoreService().Process().GetDefaultServiceCategory(kit.Ctx,
			kit.Header)
		if err != nil {
			blog.Errorf("get default service category failed, err: %v, rid: %s", err, kit.Rid)
			return serviceCategoryID, kit.CCError.Errorf(common.CCErrProcGetDefaultServiceCategoryFailed)
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
			return serviceCategoryID, err
		}
		if len(stResult.Info) == 0 {
			blog.Errorf("get service template not found, filter: %#v, rid: %s", option, kit.Rid)
			return serviceCategoryID, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
		if serviceCategoryExist && serviceCategoryID != stResult.Info[0].ServiceCategoryID {
			return serviceCategoryID, kit.CCError.Error(common.CCErrProcServiceTemplateAndCategoryNotCoincide)
		}
		serviceCategoryID = stResult.Info[0].ServiceCategoryID
	} else {
		// 检查 service category id 是否有效
		serviceCategory, err := m.clientSet.CoreService().Process().GetServiceCategory(kit.Ctx, kit.Header,
			serviceCategoryID)
		if err != nil {
			return serviceCategoryID, err
		}
		if serviceCategory.BizID != 0 && serviceCategory.BizID != bizID {
			blog.Errorf("service category and module belong to two business, categoryBizID: %d, bizID: %d, "+
				"rid: %s", serviceCategory.BizID, bizID, kit.Rid)
			return serviceCategoryID, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
		}
	}
	return serviceCategoryID, nil
}

// initModuleWithSvcTemp initialize module using the service template attributes
func (m *module) initModuleWithSvcTemp(kit *rest.Kit, bizID, svcTempID int64, module mapstr.MapStr) (mapstr.MapStr,
	error) {

	if svcTempID == common.ServiceTemplateIDNotSet {
		return module, nil
	}

	// get service template attributes
	tempAttrOpt := &metadata.ListServTempAttrOption{
		BizID: bizID,
		ID:    svcTempID,
	}

	tempAttrs, err := m.clientSet.CoreService().Process().ListServiceTemplateAttribute(kit.Ctx, kit.Header, tempAttrOpt)
	if err != nil {
		blog.Errorf("get service template attributes failed, opt: %+v, err: %v, rid: %s", tempAttrOpt, err, kit.Rid)
		return nil, err
	}

	if len(tempAttrs.Attributes) == 0 {
		return module, nil
	}

	// get corresponding module attributes
	attrIDs := make([]int64, len(tempAttrs.Attributes))
	for idx, tempAttr := range tempAttrs.Attributes {
		attrIDs[idx] = tempAttr.AttributeID
	}

	attrOpt := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKFieldID: mapstr.MapStr{common.BKDBIN: attrIDs},
		},
		Fields:         []string{common.BKFieldID, common.BKPropertyIDField},
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		DisableCounter: true,
	}

	attrs, e := m.clientSet.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, common.BKInnerObjIDModule, attrOpt)
	if e != nil {
		blog.Errorf("get module attributes failed, opt: %+v, err: %v, rid: %s", attrOpt, err, kit.Rid)
		return nil, e
	}

	// use service template attributes to initialize module data
	attrIDMap := make(map[int64]string)
	for _, attr := range attrs.Info {
		attrIDMap[attr.ID] = attr.PropertyID
	}

	for _, tempAttr := range tempAttrs.Attributes {
		propertyID, exists := attrIDMap[tempAttr.AttributeID]
		if !exists {
			blog.Errorf("service template %d attr %d is not exist, rid: %s", svcTempID, tempAttr.AttributeID, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField)
		}
		module[propertyID] = tempAttr.PropertyValue
	}

	return module, nil
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
		blog.Errorf("delete module failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// UpdateModule update module
func (m *module) UpdateModule(kit *rest.Kit, data mapstr.MapStr, bizID, setID, moduleID int64) error {

	innerCond := mapstr.MapStr{common.BKAppIDField: bizID, common.BKSetIDField: setID, common.BKModuleIDField: moduleID}
	fields := []string{common.BKSetTemplateIDField, common.BKServiceTemplateIDField, common.BKServiceCategoryIDField,
		common.BKModuleNameField}
	for field := range data {
		fields = append(fields, field)
	}

	findCond := &metadata.QueryCondition{
		Fields:         fields,
		Condition:      innerCond,
		DisableCounter: true,
	}

	moduleInstance, err := m.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header,
		common.BKInnerObjIDModule, findCond)
	if err != nil {
		blog.Errorf("get list modules failed, findCond: %#v, err: %v, rid: %s", findCond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParseDBFailed)
	}
	if len(moduleInstance.Info) > 1 {
		return kit.CCError.CCErrorf(common.CCErrCommGetMultipleObject)
	}
	if len(moduleInstance.Info) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommNotFound)
	}

	if err := m.validUpdateModuleData(kit, bizID, data, moduleInstance.Info[0]); err != nil {
		blog.Errorf("valid input data by module instance failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	data.Remove(common.BKAppIDField)
	data.Remove(common.BKSetIDField)
	data.Remove(common.BKModuleIDField)
	data.Remove(common.BKParentIDField)
	data.Remove(common.MetadataField)

	updateErr := m.inst.UpdateInst(kit, innerCond, data, common.BKInnerObjIDModule)
	if updateErr != nil {
		blog.Errorf("update module failed,  err: %v, rid: %s", updateErr, kit.Rid)
		return updateErr
	}

	return nil
}

func (m *module) validUpdateModuleData(kit *rest.Kit, bizID int64, data, moduleData mapstr.MapStr) error {
	module := new(metadata.ModuleInst)
	if err := mapstruct.Decode2Struct(moduleData, module); err != nil {
		blog.Errorf("decode original module(%+v) failed, err: %v, rid: %s", moduleData, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommJSONUnmarshalFailed)
	}

	// 检查并提示禁止修改集群模板ID字段
	if val, ok := data[common.BKSetTemplateIDField]; ok {
		setTemplateID, err := util.GetInt64ByInterface(val)
		if err != nil {
			blog.Errorf("get set template id failed, err: %v, rid: %s", err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField)
		}
		if setTemplateID != module.SetTemplateID {
			blog.Errorf("forbidden to modify set template id, rid: %s", kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommModifyFieldForbidden, common.BKSetTemplateIDField)
		}
	}

	// 检查并提示禁止修改服务模板ID字段
	if val, ok := data[common.BKServiceTemplateIDField]; ok {
		serviceTemplateID, err := util.GetInt64ByInterface(val)
		if err != nil {
			blog.Errorf("get set service template id failed, err: %v, rid: %s", err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
		if serviceTemplateID != module.ServiceTemplateID {
			blog.Errorf("forbidden to modify set service template id, rid: %s", kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommModifyFieldForbidden, common.BKServiceTemplateIDField)
		}
	}

	if module.ServiceTemplateID != common.ServiceTemplateIDNotSet {
		// 检查并提示禁止修改服务分类
		if val, ok := data[common.BKServiceCategoryIDField]; ok {
			serviceCategoryID, err := util.GetInt64ByInterface(val)
			if err != nil {
				blog.Errorf("get service category id failed, err: %v, rid: %s", err, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
			}
			if serviceCategoryID != module.ServiceCategoryID {
				blog.Errorf("forbidden to modify service category id, rid: %s", kit.Rid)
				return kit.CCError.CCError(common.CCErrorTopoUpdateModuleFromTplServiceCategoryForbidden)
			}
		}

		// 检查并提示禁止修改通过服务模板创建的模块名称
		if val, ok := data[common.BKModuleNameField]; ok {
			name := util.GetStrByInterface(val)
			if name == "" {
				delete(data, common.BKModuleNameField)
			} else if name != module.ModuleName {
				blog.Errorf("forbidden to modify module name by service template create, rid: %s", kit.Rid)
				return kit.CCError.CCError(common.CCErrorTopoUpdateModuleFromTplNameForbidden)
			}
		}

		// 检查并提示禁止修改服务模板配置的属性字段
		if err := m.validUpdateModuleDataWithTemplateAttr(kit, bizID, data, moduleData, module); err != nil {
			blog.Errorf("valid update module data(%+v) with template attr failed, err: %v, rid: %s", data, err, kit.Rid)
			return err
		}

	}

	return nil
}

// validUpdateModuleDataWithTemplateAttr validate if module update data contains service template attributes
func (m *module) validUpdateModuleDataWithTemplateAttr(kit *rest.Kit, bizID int64, data, moduleData mapstr.MapStr,
	module *metadata.ModuleInst) error {

	// get service template attributes
	tempAttrOpt := &metadata.ListServTempAttrOption{
		BizID: bizID,
		ID:    module.ServiceTemplateID,
	}

	tempAttrs, err := m.clientSet.CoreService().Process().ListServiceTemplateAttribute(kit.Ctx, kit.Header, tempAttrOpt)
	if err != nil {
		blog.Errorf("get service template attributes failed, opt: %+v, err: %v, rid: %s", tempAttrOpt, err, kit.Rid)
		return err
	}

	if len(tempAttrs.Attributes) == 0 {
		return nil
	}

	// check if update module data contains service template attributes, these attributes are forbidden to update
	attrIDs := make([]int64, len(tempAttrs.Attributes))
	for idx, tempAttr := range tempAttrs.Attributes {
		attrIDs[idx] = tempAttr.AttributeID
	}

	propertyIDs := make([]string, 0)
	for key := range data {
		propertyIDs = append(propertyIDs, key)
	}

	attrOpt := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKFieldID:         mapstr.MapStr{common.BKDBIN: attrIDs},
			common.BKPropertyIDField: mapstr.MapStr{common.BKDBIN: propertyIDs},
		},
		Fields:         []string{common.BKFieldID, common.BKPropertyIDField},
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		DisableCounter: true,
	}

	attrs, e := m.clientSet.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, common.BKInnerObjIDModule, attrOpt)
	if e != nil {
		blog.Errorf("get module attributes failed, opt: %+v, err: %v, rid: %s", attrOpt, err, kit.Rid)
		return e
	}

	fields := bytes.Buffer{}
	for _, attr := range attrs.Info {
		switch attr.PropertyType {
		case common.FieldTypeTime:
			// convert property value to time type for comparison
			updateVal, err := util.ConvToTime(data[attr.PropertyID])
			if err != nil {
				blog.Errorf("parse updated value(%+v) failed, err: %v, rid: %s", data[attr.PropertyID], err, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, attr.PropertyID)
			}

			prevVal, err := util.ConvToTime(moduleData[attr.PropertyID])
			if err != nil {
				blog.Errorf("parse prev value(%+v) failed, err: %v, rid: %s", moduleData[attr.PropertyID], err, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, attr.PropertyID)
			}

			if reflect.DeepEqual(prevVal, updateVal) {
				continue
			}
		default:
			if reflect.DeepEqual(data[attr.PropertyID], moduleData[attr.PropertyID]) {
				continue
			}
		}

		fields.WriteString(attr.PropertyID)
		fields.WriteByte(',')
	}

	if fields.Len() > 0 {
		forbiddenFields := fields.String()
		return kit.CCError.CCErrorf(common.CCErrCommModifyFieldForbidden, forbiddenFields[:len(forbiddenFields)-1])
	}

	return nil
}
