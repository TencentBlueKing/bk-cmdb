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
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// BusinessOperationInterface business operation methods
type BusinessOperationInterface interface {
	CreateBusiness(params types.ContextParams, obj model.Object, data mapstr.MapStr) (inst.Inst, error)
	DeleteBusiness(params types.ContextParams, obj model.Object, bizID int64) error
	// FindBusiness(params types.ContextParams, fields []string, cond condition.Condition) (count int, results []mapstr.MapStr, err error)
	FindBusiness(params types.ContextParams, cond *metadata.QueryBusinessRequest) (count int, results []mapstr.MapStr, err error)
	GetInternalModule(params types.ContextParams, obj model.Object, bizID int64) (count int, result *metadata.InnterAppTopo, err error)
	UpdateBusiness(params types.ContextParams, data mapstr.MapStr, obj model.Object, bizID int64) error
	HasHosts(params types.ContextParams, bizID int64) (bool, error)
	SetProxy(set SetOperationInterface, module ModuleOperationInterface, inst InstOperationInterface, obj ObjectOperationInterface)
}

// NewBusinessOperation create a business instance
func NewBusinessOperation(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) BusinessOperationInterface {
	return &business{
		clientSet:   client,
		authManager: authManager,
	}
}

type business struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
	inst        InstOperationInterface
	set         SetOperationInterface
	module      ModuleOperationInterface
	obj         ObjectOperationInterface
}

func (b *business) SetProxy(set SetOperationInterface, module ModuleOperationInterface, inst InstOperationInterface, obj ObjectOperationInterface) {
	b.inst = inst
	b.set = set
	b.module = module
	b.obj = obj
}

func (b *business) HasHosts(params types.ContextParams, bizID int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
	}
	rsp, err := b.clientSet.CoreService().Host().GetHostModuleRelation(context.Background(), params.Header, option)
	if nil != err {
		blog.Errorf("[operation-set] failed to request the object controller, error info is %s", err.Error())
		return false, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-set]  failed to search the host set configures, error info is %s", rsp.ErrMsg)
		return false, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != len(rsp.Data.Info), nil
}

func (b *business) CreateBusiness(params types.ContextParams, obj model.Object, data mapstr.MapStr) (inst.Inst, error) {

	defaultFieldVal, err := data.Int64(common.BKDefaultField)
	if nil != err {
		blog.Errorf("[operation-biz] failed to create business, error info is did not set the default field, %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}
	if defaultFieldVal == int64(common.DefaultAppFlag) && params.SupplierAccount != common.BKDefaultOwnerID {
		// this is a new supplier owner and prepare to create a new business.
		asstQuery := map[string]interface{}{
			common.BKOwnerIDField: common.BKDefaultOwnerID,
		}
		defaultOwnerHeader := util.CloneHeader(params.Header)
		defaultOwnerHeader.Set(common.BKHTTPOwnerID, common.BKDefaultOwnerID)

		asstRsp, err := b.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), defaultOwnerHeader, &metadata.QueryCondition{Condition: asstQuery})
		if nil != err {
			blog.Errorf("create business failed to get default assoc, error info is %s, rid: %s", err.Error(), params.ReqID)
			return nil, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
		}
		if !asstRsp.Result {
			return nil, params.Err.Error(asstRsp.Code)
		}
		expectAssts := asstRsp.Data.Info
		blog.Infof("copy asst for %s, %+v, rid: %s", params.SupplierAccount, expectAssts, params.ReqID)

		existAsstRsp, err := b.clientSet.CoreService().Association().ReadModelAssociation(context.Background(), params.Header, &metadata.QueryCondition{Condition: asstQuery})
		if nil != err {
			blog.Errorf("create business failed to get default assoc, error info is %s, rid: %s", err.Error(), params.ReqID)
			return nil, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
		}
		if !existAsstRsp.Result {
			return nil, params.Err.Error(existAsstRsp.Code)
		}
		existAssts := existAsstRsp.Data.Info

	expectLoop:
		for _, asst := range expectAssts {
			asst.OwnerID = params.SupplierAccount
			for _, existAsst := range existAssts {
				if existAsst.ObjectID == asst.ObjectID &&
					existAsst.AsstObjID == asst.AsstObjID &&
					existAsst.AsstKindID == asst.AsstKindID {
					continue expectLoop
				}
			}

			var createAsstRsp *metadata.CreatedOneOptionResult
			var err error
			if asst.AsstKindID == common.AssociationKindMainline {
				// bk_mainline is a inner association type that can only create in special case,
				// so we separate bk_mainline association type creation with a independent method,
				createAsstRsp, err = b.clientSet.CoreService().Association().CreateMainlineModelAssociation(context.Background(), params.Header, &metadata.CreateModelAssociation{Spec: asst})
			} else {
				createAsstRsp, err = b.clientSet.CoreService().Association().CreateModelAssociation(context.Background(), params.Header, &metadata.CreateModelAssociation{Spec: asst})
			}
			if nil != err {
				blog.Errorf("create business failed to copy default assoc, error info is %s, rid: %s", err.Error(), params.ReqID)
				return nil, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
			}
			if !createAsstRsp.Result {
				return nil, params.Err.Error(createAsstRsp.Code)
			}

		}
	}

	if util.IsExistSupplierID(params.Header) {
		supplierID, err := util.GetSupplierID(params.Header)
		if err != nil {
			return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, common.BKSupplierIDField)
		}
		data[common.BKSupplierIDField] = supplierID
	}

	bizInst, err := b.inst.CreateInst(params, obj, data)
	if nil != err {
		blog.Errorf("[operation-biz] failed to create business, error info is %s, rid: %s", err.Error(), params.ReqID)
		return bizInst, err
	}

	bizID, err := bizInst.GetInstID()
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), params.ReqID)
		return bizInst, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// register business to auth
	bizName, err := data.String(common.BKAppNameField)
	if err != nil {
		blog.Errorf("create business, but got invalid business name. err: %v, rid: %s", err, params.ReqID)
		return bizInst, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	if err := b.authManager.RegisterBusinessesByID(params.Context, params.Header, bizID); err != nil {
		blog.Errorf("create business: %s, but register business resource failed, err: %v, rid: %s", bizName, err, params.ReqID)
		return bizInst, params.Err.New(common.CCErrCommRegistResourceToIAMFailed, err.Error())
	}

	// create set
	objSet, err := b.obj.FindSingleObject(params, common.BKInnerObjIDSet)
	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	setData := mapstr.New()
	setData.Set(common.BKAppIDField, bizID)
	setData.Set(common.BKInstParentStr, bizID)
	setData.Set(common.BKSetNameField, common.DefaultResSetName)
	setData.Set(common.BKDefaultField, common.DefaultResSetFlag)

	setInst, err := b.set.CreateSet(params, objSet, bizID, setData)
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), params.ReqID)
		return bizInst, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	setID, err := setInst.GetInstID()
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), params.ReqID)
		return bizInst, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// create module
	objModule, err := b.obj.FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	defaultCategory, err := b.clientSet.CoreService().Process().GetDefaultServiceCategory(params.Context, params.Header)
	if err != nil {
		blog.Errorf("failed to search default category, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.New(common.CCErrProcGetDefaultServiceCategoryFailed, err.Error())
	}

	idleModuleData := mapstr.New()
	idleModuleData.Set(common.BKSetIDField, setID)
	idleModuleData.Set(common.BKInstParentStr, setID)
	idleModuleData.Set(common.BKAppIDField, bizID)
	idleModuleData.Set(common.BKModuleNameField, common.DefaultResModuleName)
	idleModuleData.Set(common.BKDefaultField, common.DefaultResModuleFlag)
	idleModuleData.Set(common.BKServiceTemplateIDField, common.ServiceTemplateIDNotSet)
	idleModuleData.Set(common.BKSetTemplateIDField, common.SetTemplateIDNotSet)
	idleModuleData.Set(common.BKServiceCategoryIDField, defaultCategory.ID)

	_, err = b.module.CreateModule(params, objModule, bizID, setID, idleModuleData)
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), params.ReqID)
		return bizInst, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// create fault module
	faultModuleData := mapstr.New()
	faultModuleData.Set(common.BKSetIDField, setID)
	faultModuleData.Set(common.BKInstParentStr, setID)
	faultModuleData.Set(common.BKAppIDField, bizID)
	faultModuleData.Set(common.BKModuleNameField, common.DefaultFaultModuleName)
	faultModuleData.Set(common.BKDefaultField, common.DefaultFaultModuleFlag)
	faultModuleData.Set(common.BKServiceTemplateIDField, common.ServiceTemplateIDNotSet)
	faultModuleData.Set(common.BKSetTemplateIDField, common.SetTemplateIDNotSet)
	faultModuleData.Set(common.BKServiceCategoryIDField, defaultCategory.ID)

	_, err = b.module.CreateModule(params, objModule, bizID, setID, faultModuleData)
	if nil != err {
		blog.Errorf("create business failed to create business, error info is %s, rid: %s", err.Error(), params.ReqID)
		return bizInst, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	// create recycle module
	recycleModuleData := mapstr.New()
	recycleModuleData.Set(common.BKSetIDField, setID)
	recycleModuleData.Set(common.BKInstParentStr, setID)
	recycleModuleData.Set(common.BKAppIDField, bizID)
	recycleModuleData.Set(common.BKModuleNameField, common.DefaultRecycleModuleName)
	recycleModuleData.Set(common.BKDefaultField, common.DefaultRecycleModuleFlag)
	recycleModuleData.Set(common.BKServiceTemplateIDField, common.ServiceTemplateIDNotSet)
	recycleModuleData.Set(common.BKSetTemplateIDField, common.SetTemplateIDNotSet)
	recycleModuleData.Set(common.BKServiceCategoryIDField, defaultCategory.ID)

	_, err = b.module.CreateModule(params, objModule, bizID, setID, recycleModuleData)
	if nil != err {
		blog.Errorf("create business failed, create recycle module failed, err: %s, rid: %s", err.Error(), params.ReqID)
		return bizInst, params.Err.New(common.CCErrTopoAppCreateFailed, err.Error())
	}

	return bizInst, nil
}

func (b *business) DeleteBusiness(params types.ContextParams, obj model.Object, bizID int64) error {
	if err := b.authManager.DeregisterBusinessByRawID(params.Context, params.Header, bizID); err != nil {
		blog.Errorf("delete business: %d, but deregister business from auth failed, err: %v, rid: %s", bizID, err, params.ReqID)
		return params.Err.New(common.CCErrCommUnRegistResourceToIAMFailed, err.Error())
	}

	setObj, err := b.obj.FindSingleObject(params, common.BKInnerObjIDSet)
	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), params.ReqID)
		return err
	}

	bizModel, err := b.obj.FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the set, %s, rid: %s", err.Error(), params.ReqID)
		return err
	}

	if err = b.set.DeleteSet(params, setObj, bizID, nil); nil != err {
		blog.Errorf("[operation-biz] failed to delete the set, error info is %s, rid: %s", err.Error(), params.ReqID)
		return params.Err.New(common.CCErrTopoAppDeleteFailed, err.Error())
	}

	innerCond := condition.CreateCondition()
	innerCond.Field(common.BKAppIDField).Eq(bizID)

	return b.inst.DeleteInst(params, bizModel, innerCond, true)
}

func (b *business) FindBusiness(params types.ContextParams, cond *metadata.QueryBusinessRequest) (count int, results []mapstr.MapStr, err error) {

	cond.Condition[common.BKDefaultField] = 0
	query := &metadata.QueryCondition{
		Fields:    cond.Fields,
		Condition: cond.Condition,
		Limit: metadata.SearchLimit{
			Limit:  int64(cond.Page.Limit),
			Offset: int64(cond.Page.Start),
		},
		SortArr: metadata.NewSearchSortParse().String(cond.Page.Sort).ToSearchSortArr(),
	}

	result, err := b.clientSet.CoreService().Instance().ReadInstance(params.Context, params.Header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.ErrorJSON("failed to find business by query condition: %s, err: %s, rid: %s", query, err.Error(), params.ReqID)
		return 0, nil, err
	}

	if !result.Result {
		return 0, nil, params.Err.Errorf(result.Code, result.ErrMsg)
	}

	return result.Data.Count, result.Data.Info, err
}

func (b *business) GetInternalModule(params types.ContextParams, obj model.Object, bizID int64) (count int, result *metadata.InnterAppTopo, err error) {
	// get set model
	setObj, err := b.obj.FindSingleObject(params, common.BKInnerObjIDSet)
	if nil != err {
		return 0, nil, params.Err.New(common.CCErrTopoAppSearchFailed, err.Error())
	}

	// search internal sets
	querySet := &metadata.QueryInput{}
	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	cond.Field(common.BKDefaultField).Eq(common.DefaultResModuleFlag)
	querySet.Condition = cond.ToMapStr()
	_, sets, err := b.set.FindSet(params, setObj, querySet)
	if nil != err {
		return 0, nil, params.Err.New(common.CCErrTopoAppSearchFailed, err.Error())
	}

	// get module model
	moduleObj, err := b.obj.FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		return 0, nil, params.Err.New(common.CCErrTopoAppSearchFailed, err.Error())
	}

	// search internal modules
	queryModule := &metadata.QueryInput{}
	cond = condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	cond.Field(common.BKDefaultField).NotEq(common.DefaultFlagDefaultValue)
	queryModule.Condition = cond.ToMapStr()
	_, modules, err := b.module.FindModule(params, moduleObj, queryModule)
	if nil != err {
		return 0, nil, params.Err.New(common.CCErrTopoAppSearchFailed, err.Error())
	}

	// construct result
	result = &metadata.InnterAppTopo{}
	for _, set := range sets {
		id, err := set.GetInstID()
		if nil != err {
			return 0, nil, params.Err.New(common.CCErrTopoAppSearchFailed, err.Error())
		}
		name, err := set.GetInstName()
		if nil != err {
			return 0, nil, params.Err.New(common.CCErrTopoAppSearchFailed, err.Error())
		}

		result.SetID = id
		result.SetName = name
		break // should be only one set
	}

	for _, moduleMapStr := range modules {
		module := metadata.ModuleInst{}
		if err := mapstruct.Decode2Struct(moduleMapStr, &module); err != nil {
			blog.ErrorJSON("GetInternalModule failed, unmarshal module failed, module: %s, err: %s, rid: %s", moduleMapStr, err.Error(), params.ReqID)
			return 0, nil, params.Err.CCError(common.CCErrCommParseDBFailed)
		}
		result.Module = append(result.Module, metadata.InnerModule{
			ModuleID:   module.ModuleID,
			ModuleName: module.ModuleName,
			Default:    module.Default,
		})
	}

	return 0, result, nil
}

func (b *business) UpdateBusiness(params types.ContextParams, data mapstr.MapStr, obj model.Object, bizID int64) error {
	if biz, exist := data.Get(common.BKAppNameField); exist {
		bizName, err := data.String(common.BKAppNameField)
		if err != nil {
			blog.Errorf("update business, but got invalid business name: %v, id: %d, rid: %s", biz, bizID, params.ReqID)
			return params.Err.Error(common.CCErrCommParamsIsInvalid)
		}

		if err := b.authManager.UpdateRegisteredBusinessByID(params.Context, params.Header, bizID); err != nil {
			blog.Errorf("update business name: %s, but update resource to auth failed, err: %v, rid: %s", bizName, err, params.ReqID)
			return params.Err.New(common.CCErrCommRegistResourceToIAMFailed, err.Error())
		}
	}

	innerCond := condition.CreateCondition()
	innerCond.Field(common.BKAppIDField).Eq(bizID)

	return b.inst.UpdateInst(params, data, obj, innerCond, bizID)
}
