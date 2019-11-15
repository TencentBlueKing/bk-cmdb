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

// SetOperationInterface set operation methods
type SetOperationInterface interface {
	CreateSet(params types.ContextParams, obj model.Object, bizID int64, data mapstr.MapStr) (inst.Inst, error)
	DeleteSet(params types.ContextParams, obj model.Object, bizID int64, setIDS []int64) error
	FindSet(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []inst.Inst, err error)
	UpdateSet(params types.ContextParams, data mapstr.MapStr, obj model.Object, bizID, setID int64) error

	SetProxy(obj ObjectOperationInterface, inst InstOperationInterface, module ModuleOperationInterface)
}

// NewSetOperation create a set instance
func NewSetOperation(client apimachinery.ClientSetInterface) SetOperationInterface {
	return &set{
		clientSet: client,
	}
}

type set struct {
	clientSet apimachinery.ClientSetInterface
	inst      InstOperationInterface
	obj       ObjectOperationInterface
	module    ModuleOperationInterface
}

func (s *set) SetProxy(obj ObjectOperationInterface, inst InstOperationInterface, module ModuleOperationInterface) {
	s.inst = inst
	s.obj = obj
	s.module = module
}

func (s *set) hasHost(params types.ContextParams, bizID int64, setIDS []int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		SetIDArr:      setIDS,
	}
	rsp, err := s.clientSet.CoreService().Host().GetHostModuleRelation(context.Background(), params.Header, option)
	if nil != err {
		blog.Errorf("[operation-set] failed to request the object controller, error info is %s, rid: %s", err.Error(), params.ReqID)
		return false, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-set]  failed to search the host set configures, error info is %s, rid: %s", rsp.ErrMsg, params.ReqID)
		return false, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != len(rsp.Data.Info), nil
}

func (s *set) CreateSet(params types.ContextParams, obj model.Object, bizID int64, data mapstr.MapStr) (inst.Inst, error) {

	data.Set(common.BKAppIDField, bizID)

	if !data.Exists(common.BKDefaultField) {
		data.Set(common.BKDefaultField, common.DefaultFlagDefaultValue)
	}

	setTemplate := metadata.SetTemplate{}
	// validate foreign key
	if setTemplateIDIf, ok := data[common.BKSetTemplateIDField]; ok == true {
		setTemplateID, err := util.GetInt64ByInterface(setTemplateIDIf)
		if err != nil {
			blog.Errorf("parse set_template_id field into int failed, id: %+v, rid: %s", setTemplateIDIf, params.ReqID)
			err := params.Err.CCErrorf(common.CCErrCommParamsInvalid, params.Lang.Language("set_property_set_template_id"))
			return nil, err
		}
		if setTemplateID != common.SetTemplateIDNotSet {
			st, err := s.clientSet.CoreService().SetTemplate().GetSetTemplate(params.Context, params.Header, bizID, setTemplateID)
			if err != nil {
				err := params.Err.CCErrorf(common.CCErrCommParamsInvalid, params.Lang.Language("set_property_set_template_id"))
				return nil, err
			}
			setTemplate = st
		}
	} else {
		data[common.BKSetTemplateIDField] = common.SetTemplateIDNotSet
	}

	// TODO: run in transaction
	data.Set(common.BKSetTemplateIDField, setTemplate.ID)
	data.Remove(common.MetadataField)
	setInstance, err := s.inst.CreateInst(params, obj, data)
	if err != nil {
		blog.Errorf("create set instance failed, object: %+v, data: %+v, err: %s, rid: %s", obj, data, err.Error(), params.ReqID)
		return setInstance, err
	}
	if setTemplate.ID == 0 {
		return setInstance, nil
	}

	setID, err := setInstance.GetInstID()
	if err != nil {
		blog.Errorf("create set instance success, but read instance id field failed, bizID: %d, setInstance: %+v, err: %s, rid: %s", bizID, setInstance, err.Error(), params.ReqID)
		return setInstance, err
	}

	// set create by template should create module at the same time
	serviceTemplates, err := s.clientSet.CoreService().SetTemplate().ListSetTplRelatedSvcTpl(params.Context, params.Header, bizID, setTemplate.ID)
	if err != nil {
		blog.Errorf("create set failed, ListSetTplRelatedSvcTpl failed, bizID: %d, setTemplateID: %d, err: %s, rid: %s", bizID, setTemplate.ID, err.Error(), params.ReqID)
		return setInstance, err
	}

	moduleObj, err := s.obj.FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("[operation-set] failed to find module object, error info is %s, rid: %s", err.Error(), params.ReqID)
		return setInstance, err
	}
	for _, serviceTemplate := range serviceTemplates {
		createModuleParam := map[string]interface{}{
			common.BKModuleNameField:        serviceTemplate.Name,
			common.BKServiceTemplateIDField: serviceTemplate.ID,
			common.BKSetTemplateIDField:     setTemplate.ID,
			common.BKParentIDField:          setID,
			common.BKServiceCategoryIDField: serviceTemplate.ServiceCategoryID,
			common.BKAppIDField:             bizID,
		}
		_, err := s.module.CreateModule(params, moduleObj, bizID, setID, createModuleParam)
		if err != nil {
			blog.Errorf("create set instance failed, create module instance failed, bizID: %d, setID: %d, param: %+v, err: %s, rid: %s", bizID, setID, createModuleParam, err.Error(), params.ReqID)
			return setInstance, err
		}
	}

	return setInstance, nil
}

func (s *set) DeleteSet(params types.ContextParams, setModel model.Object, bizID int64, setIDS []int64) error {

	setCond := condition.CreateCondition()

	// clear the sets

	setCond.Field(common.BKAppIDField).Eq(bizID)
	if nil != setIDS {
		setCond.Field(common.BKSetIDField).In(setIDS)
	}

	exists, err := s.hasHost(params, bizID, setIDS)
	if nil != err {
		blog.Errorf("[operation-set] failed to check the host, error info is %s, rid: %s", err.Error(), params.ReqID)
		return err
	}

	if exists {
		blog.Errorf("[operation-set] the sets(%#v) has some hosts, rid: %s", setIDS, params.ReqID)
		return params.Err.Error(common.CCErrTopoHasHostCheckFailed)
	}

	// clear the module belong to deleted sets
	moduleObj, err := s.obj.FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("[operation-set] failed to find the object , error info is %s, rid: %s", err.Error(), params.ReqID)
		return err
	}

	if err = s.module.DeleteModule(params, moduleObj, bizID, setIDS, nil); nil != err {
		blog.Errorf("[operation-set] failed to delete the modules, error info is %s, rid: %s", err.Error(), params.ReqID)
		return params.Err.New(common.CCErrTopoSetDeleteFailed, err.Error())
	}

	// clear set template sync status
	if ccErr := s.clientSet.CoreService().SetTemplate().DeleteSetTemplateSyncStatus(params.Context, params.Header, bizID, setIDS); ccErr != nil {
		blog.Errorf("[operation-set] failed to delete set template sync status failed, bizID: %d, setIDs: %+v, err: %s, rid: %s", bizID, setIDS, ccErr.Error(), params.ReqID)
		return ccErr
	}

	// clear the sets
	return s.inst.DeleteInst(params, setModel, setCond, false)
}

func (s *set) FindSet(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []inst.Inst, err error) {
	return s.inst.FindInst(params, obj, cond, false)
}

func (s *set) UpdateSet(params types.ContextParams, data mapstr.MapStr, obj model.Object, bizID, setID int64) error {

	innerCond := condition.CreateCondition()

	innerCond.Field(common.BKAppIDField).Eq(bizID)
	innerCond.Field(common.BKSetIDField).Eq(setID)

	data.Remove(common.MetadataField)
	data.Remove(common.BKAppIDField)
	data.Remove(common.BKSetIDField)

	return s.inst.UpdateInst(params, data, obj, innerCond, setID)
}
