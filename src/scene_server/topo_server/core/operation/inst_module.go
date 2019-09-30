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

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// ModuleOperationInterface module operation methods
type ModuleOperationInterface interface {
	CreateModule(params types.ContextParams, obj model.Object, bizID, setID int64, data mapstr.MapStr) (inst.Inst, error)
	DeleteModule(params types.ContextParams, obj model.Object, bizID int64, setID, moduleIDS []int64) error
	FindModule(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []mapstr.MapStr, err error)
	UpdateModule(params types.ContextParams, data mapstr.MapStr, obj model.Object, bizID, setID, moduleID int64) error

	SetProxy(inst InstOperationInterface)
}

// NewModuleOperation create a new module
func NewModuleOperation(client apimachinery.ClientSetInterface) ModuleOperationInterface {
	return &module{
		clientSet: client,
	}
}

type module struct {
	clientSet apimachinery.ClientSetInterface
	inst      InstOperationInterface
}

func (m *module) SetProxy(inst InstOperationInterface) {
	m.inst = inst
}

func (m *module) hasHost(params types.ContextParams, bizID int64, setIDs, moduleIDS []int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		ModuleIDArr:   moduleIDS,
	}
	if len(setIDs) > 0 {
		option.SetIDArr = setIDs
	}
	if len(moduleIDS) > 0 {
		option.ModuleIDArr = moduleIDS
	}
	rsp, err := m.clientSet.CoreService().Host().GetHostModuleRelation(context.Background(), params.Header, option)
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

func (m *module) CreateModule(params types.ContextParams, obj model.Object, bizID, setID int64, data mapstr.MapStr) (inst.Inst, error) {

	data.Set(common.BKSetIDField, setID)
	data.Set(common.BKAppIDField, bizID)
	if !data.Exists(common.BKDefaultField) {
		data.Set(common.BKDefaultField, 0)
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
			return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
		}
		serviceCategoryID = scID
	}
	var serviceTemplateID int64
	var err error
	serviceTemplateIDIf, serviceTemplateFieldExist := data.Get(common.BKServiceTemplateIDField)
	if serviceTemplateFieldExist == true {
		serviceTemplateID, err = util.GetInt64ByInterface(serviceTemplateIDIf)
		if err != nil {
			return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
	}
	data.Set(common.BKServiceTemplateIDField, serviceTemplateID)
	if serviceCategoryExist == false && (serviceTemplateFieldExist == false || serviceTemplateID == common.ServiceTemplateIDNotSet) {
		// set default service template id
		defaultServiceCategory, err := m.clientSet.CoreService().Process().GetDefaultServiceCategory(params.Context, params.Header)
		if err != nil {
			blog.Errorf("create module failed, GetDefaultServiceCategory failed, err: %s, rid: %s", err.Error(), params.ReqID)
			return nil, params.Err.Errorf(common.CCErrProcGetDefaultServiceCategoryFailed)
		}
		serviceCategoryID = defaultServiceCategory.ID
	}
	if serviceTemplateID != common.ServiceTemplateIDNotSet {
		// 校验 serviceCategoryID 与 serviceTemplateID 对应
		templateIDs := []int64{serviceTemplateID}
		option := metadata.ListServiceTemplateOption{
			BusinessID:         bizID,
			ServiceTemplateIDs: templateIDs,
		}
		stResult, err := m.clientSet.CoreService().Process().ListServiceTemplates(params.Context, params.Header, &option)
		if err != nil {
			return nil, err
		}
		if len(stResult.Info) == 0 {
			blog.Errorf("create module failed, service template not found, filter: %+v, rid: %s", option, params.ReqID)
			return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
		if serviceCategoryExist == true && serviceCategoryID != stResult.Info[0].ServiceCategoryID {
			return nil, params.Err.Error(common.CCErrProcServiceTemplateAndCategoryNotCoincide)
		}
	} else {
		// 检查 service category id 是否有效
		serviceCategory, err := m.clientSet.CoreService().Process().GetServiceCategory(params.Context, params.Header, serviceCategoryID)
		if err != nil {
			return nil, err
		}
		categoryBizID, parseErr := serviceCategory.Metadata.ParseBizID()
		if parseErr != nil {
			blog.ErrorJSON("create module failed, parse biz id from db data failed, data: %s, rid: %s", categoryBizID, params.ReqID)
			return nil, params.Err.Errorf(common.CCErrCommParseDataFailed)
		}
		if categoryBizID != 0 && categoryBizID != bizID {
			blog.V(3).Info("create module failed, service category and module belong to two business, categoryBizID: %d, bizID: %d, rid: %s", categoryBizID, bizID, params.ReqID)
			return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
		}
	}
	data.Set(common.BKServiceCategoryIDField, serviceCategoryID)

	return m.inst.CreateInst(params, obj, data)
}

func (m *module) DeleteModule(params types.ContextParams, obj model.Object, bizID int64, setIDs, moduleIDS []int64) error {

	exists, err := m.hasHost(params, bizID, setIDs, moduleIDS)
	if nil != err {
		blog.Errorf("[operation-module] failed to delete the modules, err: %s, rid: %s", err.Error(), params.ReqID)
		return err
	}

	if exists {
		blog.Errorf("[operation-module]the module has some hosts, can not be deleted, rid: %s", params.ReqID)
		return params.Err.Error(common.CCErrTopoHasHost)
	}

	innerCond := condition.CreateCondition()
	innerCond.Field(common.BKAppIDField).Eq(bizID)
	if nil != setIDs {
		innerCond.Field(common.BKSetIDField).In(setIDs)
	}

	if nil != moduleIDS {
		innerCond.Field(common.BKModuleIDField).In(moduleIDS)
	}

	// module table doesn't have metadata field
	params.MetaData = nil
	return m.inst.DeleteInst(params, obj, innerCond, false)
}

func (m *module) FindModule(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []mapstr.MapStr, err error) {
	// module table doesn't have metadata field
	params.MetaData = nil
	count, resultData, err := m.inst.FindInst(params, obj, cond, false)
	if err != nil {
		return 0, nil, err
	}
	moduleInstances := make([]mapstr.MapStr, 0)
	for _, item := range resultData {
		moduleInstance := make(map[string]interface{})
		if err := mapstr.DecodeFromMapStr(&moduleInstance, item.ToMapStr()); err != nil {
			blog.Errorf("unmarshal module into struct failed, module: %+v, rid: %s", item, params.ReqID)
			return 0, nil, err
		}
		moduleInstances = append(moduleInstances, moduleInstance)
	}
	return count, moduleInstances, err
}

func (m *module) UpdateModule(params types.ContextParams, data mapstr.MapStr, obj model.Object, bizID, setID, moduleID int64) error {
	innerCond := condition.CreateCondition()

	innerCond.Field(common.BKAppIDField).Eq(bizID)
	innerCond.Field(common.BKSetIDField).Eq(setID)
	innerCond.Field(common.BKModuleIDField).Eq(moduleID)

	findCond := &metadata.QueryInput{
		Condition: innerCond.ToMapStr(),
	}
	var err error
	count, moduleInstances, err := m.FindModule(params, obj, findCond)
	if err != nil {
		blog.Errorf("update module failed, find module failed, filter: %+v, err: %s, rid: %s", findCond, err.Error(), params.ReqID)
		return err
	}
	if count == 0 {
		return params.Err.CCErrorf(common.CCErrCommNotFound)
	}
	if count > 1 {
		return params.Err.CCErrorf(common.CCErrCommGetMultipleObject)
	}
	if len(moduleInstances) == 0 {
		return params.Err.CCErrorf(common.CCErrCommNotFound)
	}

	moduleMapStr := moduleInstances[0]
	moduleInstance := metadata.ModuleInst{}
	if err := moduleMapStr.MarshalJSONInto(&moduleInstance); err != nil {
		blog.ErrorJSON("unmarshal db data into module failed, module: %s, err: %s, rid: %s", moduleMapStr, err.Error(), params.ReqID)
		return params.Err.CCError(common.CCErrCommParseDBFailed)
	}

	if moduleInstance.ServiceTemplateID != common.ServiceTemplateIDNotSet {
		// 检查并提示禁止修改服务分类
		if val, ok := data[common.BKServiceCategoryIDField]; ok == true {
			serviceCategoryID, err := util.GetInt64ByInterface(val)
			if err != nil {
				return params.Err.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
			}
			if serviceCategoryID != moduleInstance.ServiceCategoryID {
				return params.Err.CCError(common.CCErrorTopoUpdateModuleFromTplServiceCategoryForbidden)
			}
		}

		// 检查并提示禁止修改通过模板创建的模块名称
		if val, ok := data[common.BKModuleNameField]; ok == true {
			name := util.GetStrByInterface(val)
			if len(name) == 0 {
				delete(data, common.BKModuleNameField)
			} else if name != moduleInstance.ModuleName {
				return params.Err.CCError(common.CCErrorTopoUpdateModuleFromTplNameForbidden)
			}
		}
	}

	// module table don't have metadata field
	params.MetaData = nil
	return m.inst.UpdateInst(params, data, obj, innerCond, -1)
}
