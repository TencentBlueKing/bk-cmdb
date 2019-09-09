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

package compatiblev2

import (
	"context"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/types"
)

// SetInterface set interface
type SetInterface interface {
	UpdateMultiSet(bizID int64, data mapstr.MapStr, cond condition.Condition) error
	DeleteMultiSet(bizID int64, setIDS []int64) error
	// DeleteSetHost(bizID int64, option interface{}) error
	DeleteSetHost(bizID int64, setIDs []int64) error
}

// NewSet ceate a new set instance
func NewSet(params types.ContextParams, client apimachinery.ClientSetInterface) SetInterface {
	return &set{
		params: params,
		client: client,
	}
}

type set struct {
	params types.ContextParams
	client apimachinery.ClientSetInterface
}

func (s *set) hasHost(bizID int64, setIDS []int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		SetIDArr:      setIDS,
	}
	rsp, err := s.client.CoreService().Host().GetHostModuleRelation(context.Background(), s.params.Header, option)
	if nil != err {
		blog.Errorf("[compatiblev2-set] failed to request the object controller, err: %s, rid: %s", err.Error(), s.params.ReqID)
		return false, s.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-set]  failed to search the host module configures, err: %s, rid: %s", rsp.ErrMsg, s.params.ReqID)
		return false, s.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != len(rsp.Data.Info), nil
}

func (s *set) deleteModules(bizID int64, setIDS []int64) error {

	cond := condition.CreateCondition()
	cond.Field(common.BKSetIDField).In(setIDS)
	cond.Field(common.BKAppIDField).Eq(bizID)

	deleteParam := &metadata.DeleteOption{
		Condition: cond.ToMapStr(),
	}
	rsp, err := s.client.CoreService().Instance().DeleteInstance(s.params.Context, s.params.Header, common.BKInnerObjIDModule, deleteParam)

	if nil != err {
		blog.Errorf("[compatiblev2-set] failed to request the object controller, err: %s, rid: %s", err.Error(), s.params.ReqID)
		return s.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-set]  failed to delete module by the condition(%#v), err: %s, rid: %s", cond.ToMapStr(), rsp.ErrMsg, s.params.ReqID)
		return s.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (s *set) UpdateMultiSet(bizID int64, data mapstr.MapStr, cond condition.Condition) error {
	updateParam := &metadata.UpdateOption{
		Data:      data,
		Condition: cond.ToMapStr(),
	}
	rsp, err := s.client.CoreService().Instance().UpdateInstance(s.params.Context, s.params.Header, common.BKInnerObjIDSet, updateParam)
	if nil != err {
		blog.Errorf("[compatiblev2-set] failed to request the object controller, err: %s, rid: %s", rsp.ErrMsg, s.params.ReqID)
		return s.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-set]  failed to update the set, err: %s, rid: %s", err.Error(), s.params.ReqID)
		return s.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (s *set) DeleteMultiSet(bizID int64, setIDS []int64) error {

	// checkout host
	existsHost, err := s.hasHost(bizID, setIDS)
	if nil != err {
		blog.Errorf("[compatiblev2-set] faield to check host, err: %s, rid: %s", err.Error(), s.params.ReqID)
		return err
	}

	if existsHost {
		blog.Error("[compatiblev2-set] the setids has host, disable be deleted, rid: %s", s.params.ReqID)
		return s.params.Err.Error(common.CCErrTopoHasHost)
	}

	// delete module
	if err = s.deleteModules(bizID, setIDS); nil != err {
		return err
	}

	// delete sets
	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	cond.Field(common.BKSetIDField).In(setIDS)

	deleteParam := &metadata.DeleteOption{
		Condition: cond.ToMapStr(),
	}
	rsp, err := s.client.CoreService().Instance().DeleteInstance(s.params.Context, s.params.Header, common.BKInnerObjIDSet, deleteParam)
	if nil != err {
		blog.Errorf("[compatiblev2-set] faield to check host, err: %s, rid: %s", err.Error(), s.params.ReqID)
		return err
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-set]  failed to delete the sets by the condition(%#v), err: %s, rid: %s", cond.ToMapStr(), rsp.ErrMsg, s.params.ReqID)
		return s.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (s *set) DeleteSetHost(bizID int64, setIDs []int64) error {
	defaultModule, err := s.client.CoreService().Process().GetBusinessDefaultSetModuleInfo(s.params.Context, s.params.Header, bizID)
	if err != nil {
		blog.Errorf("[compatiblev2-set] failed to get biz default modules, err: %s, rid: %s, rid: %s", err.Error(), s.params.ReqID, s.params.ReqID)
		return err
	}

	relationParam := metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		SetIDArr:      setIDs,
	}
	relations, ccerr := s.client.CoreService().Host().GetHostModuleRelation(s.params.Context, s.params.Header, &relationParam)
	if ccerr != nil {
		blog.Errorf("[compatiblev2-set] failed find hosts, err: %s, rid: %s, rid: %s", err.Error(), s.params.ReqID, s.params.ReqID)
		return ccerr
	}
	hostIDs := make([]int64, 0)
	for _, relation := range relations.Data.Info {
		hostIDs = append(hostIDs, relation.HostID)
	}
	hostIDs = util.IntArrayUnique(hostIDs)
	transferOption := &metadata.TransferHostToInnerModule{
		ApplicationID: bizID,
		ModuleID:      defaultModule.IdleModuleID,
		HostID:        hostIDs,
	}
	rsp, ccerr := s.client.CoreService().Host().TransferToInnerModule(s.params.Context, s.params.Header, transferOption)
	if nil != ccerr {
		blog.Errorf("[compatiblev2-set] failed to delete the set hosts, err: %s, rid: %s, rid: %s", err.Error(), s.params.ReqID, s.params.ReqID)
		return ccerr
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-set] failed to delete the set hosts, err: %s, rid: %s", rsp.ErrMsg, s.params.ReqID)
		return s.params.Err.New(rsp.Code, rsp.ErrMsg)
	}
	return nil
}
