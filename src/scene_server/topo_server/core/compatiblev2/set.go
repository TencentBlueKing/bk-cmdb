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
	"configcenter/src/scene_server/topo_server/core/types"
)

// SetInterface set interface
type SetInterface interface {
	UpdateMultiSet(bizID int64, data mapstr.MapStr, cond condition.Condition) error
	DeleteMultiSet(bizID int64, setIDS []int64) error
	DeleteSetHost(bizID int64, cond condition.Condition) error
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
	cond := map[string][]int64{
		"ApplicationID": []int64{bizID},
		"SetID":         setIDS,
	}

	rsp, err := s.client.HostController().Module().GetModulesHostConfig(context.Background(), s.params.Header, cond)
	if nil != err {
		blog.Errorf("[compatiblev2-set] failed to request the object controller, error info is %s", err.Error())
		return false, s.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-set]  failed to search the host module configures, error info is %s", err.Error())
		return false, s.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != len(rsp.Data), nil
}

func (s *set) deleteModules(bizID int64, setIDS []int64) error {

	cond := condition.CreateCondition()
	cond.Field(common.BKSetIDField).In(setIDS)
	cond.Field(common.BKAppIDField).Eq(bizID)

	rsp, err := s.client.ObjectController().Instance().DelObject(context.Background(), common.BKInnerObjIDModule, s.params.Header, cond.ToMapStr())

	if nil != err {
		blog.Errorf("[compatiblev2-set] failed to request the object controller, error info is %s", err.Error())
		return s.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-set]  failed to delete module by the condition(%#v), error info is %s", cond.ToMapStr(), rsp.ErrMsg)
		return s.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (s *set) UpdateMultiSet(bizID int64, data mapstr.MapStr, cond condition.Condition) error {

	input := mapstr.New()
	input.Set("data", data)
	input.Set("condition", cond.ToMapStr())

	rsp, err := s.client.ObjectController().Instance().UpdateObject(context.Background(), common.BKInnerObjIDSet, s.params.Header, input)
	if nil != err {
		blog.Errorf("[compatiblev2-set] failed to request the object controller, error info is %s", err.Error())
		return s.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-set]  failed to update the set, error info is %s", err.Error())
		return s.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (s *set) DeleteMultiSet(bizID int64, setIDS []int64) error {

	// checkout host
	existsHost, err := s.hasHost(bizID, setIDS)
	if nil != err {
		blog.Errorf("[compatiblev2-set] faield to check host, error info is %s", err.Error())
		return err
	}

	if existsHost {
		blog.Error("[compatiblev2-set] the setids has host, disable be deleted")
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

	rsp, err := s.client.ObjectController().Instance().DelObject(context.Background(), common.BKInnerObjIDSet, s.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[compatiblev2-set] faield to check host, error info is %s", err.Error())
		return err
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-set]  failed to delete the sets by the condition(%#v), error info is %s", cond.ToMapStr(), err.Error())
		return s.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}
func (s *set) DeleteSetHost(bizID int64, cond condition.Condition) error {

	rsp, err := s.client.ObjectController().OpenAPI().DeleteSetHost(context.Background(), s.params.Header, cond.ToMapStr())

	if nil != err {
		blog.Errorf("[compatiblev2-set] failed to delete the set hosts, error info is %s", err.Error())
		return s.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-set] failed to delete the set hosts, error info is %s", rsp.ErrMsg)
		return s.params.Err.New(rsp.Code, rsp.ErrMsg)
	}
	return nil
}
