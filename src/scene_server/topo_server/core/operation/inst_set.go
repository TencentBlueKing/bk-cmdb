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
	cond := map[string][]int64{
		common.BKAppIDField: []int64{bizID},
		common.BKSetIDField: setIDS,
	}

	rsp, err := s.clientSet.HostController().Module().GetModulesHostConfig(context.Background(), params.Header, cond)
	if nil != err {
		blog.Errorf("[operation-set] failed to request the object controller, error info is %s", err.Error())
		return false, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-set]  failed to search the host set configures, error info is %s", err.Error())
		return false, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != len(rsp.Data), nil
}

func (s *set) CreateSet(params types.ContextParams, obj model.Object, bizID int64, data mapstr.MapStr) (inst.Inst, error) {

	data.Set(common.BKAppIDField, bizID)

	if !data.Exists(common.BKDefaultField) {
		data.Set(common.BKDefaultField, 0)
	}

	//data.Set(common.CreateTimeField, util.GetCurrentTimeStr())
	return s.inst.CreateInst(params, obj, data)
}

func (s *set) DeleteSet(params types.ContextParams, obj model.Object, bizID int64, setIDS []int64) error {

	setCond := condition.CreateCondition()

	// clear the sets

	setCond.Field(common.BKAppIDField).Eq(bizID)
	if nil != setIDS {
		setCond.Field(common.BKSetIDField).In(setIDS)
	}

	exists, err := s.hasHost(params, bizID, setIDS)
	if nil != err {
		blog.Errorf("[operation-set] failed to check the host, error info is %s", err.Error())
		return err
	}

	if exists {
		blog.Errorf("[operation-set] the sets(%#v) has some hosts", setIDS)
		return params.Err.Error(common.CCErrTopoHasHostCheckFailed)
	}

	// clear the moudle belong to deleted sets
	moduleObj, err := s.obj.FindSingleObject(params, common.BKInnerObjIDModule)
	if nil != err {
		blog.Errorf("[operation-set] failed to find the object , error info is %s", err.Error())
		return err
	}

	if err = s.module.DeleteModule(params, moduleObj, bizID, setIDS, nil); nil != err {
		blog.Errorf("[operation-set] failed to delete the modules, error info is %s", err.Error())
		return params.Err.New(common.CCErrTopoSetDeleteFailed, err.Error())
	}

	// clear the sets
	return s.inst.DeleteInst(params, obj, setCond, false)
}

func (s *set) FindSet(params types.ContextParams, obj model.Object, cond *metadata.QueryInput) (count int, results []inst.Inst, err error) {
	return s.inst.FindInst(params, obj, cond, false)
}

func (s *set) UpdateSet(params types.ContextParams, data mapstr.MapStr, obj model.Object, bizID, setID int64) error {

	innerCond := condition.CreateCondition()

	innerCond.Field(common.BKAppIDField).Eq(bizID)
	innerCond.Field(common.BKSetIDField).Eq(setID)

	return s.inst.UpdateInst(params, data, obj, innerCond, setID)
}
