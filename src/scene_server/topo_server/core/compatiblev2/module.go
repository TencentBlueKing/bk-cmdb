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
	"configcenter/src/scene_server/topo_server/core/types"
)

// ModuleInterface module interface
type ModuleInterface interface {
	UpdateMultiModule(bizID int64, moduleIDS interface{}, innerData mapstr.MapStr) error
	SearchModuleByApp(query *metadata.QueryCondition) (*metadata.InstResult, error)
	SearchModuleBySetProperty(bizID int64, cond condition.Condition) (*metadata.InstResult, error)
	AddMultiModule(bizID, setID int64, moduleNames []string, data mapstr.MapStr) error
	DeleteMultiModule(bizID int64, moduleIDS []int64) error
}

// NewModule create a module instance
func NewModule(params types.ContextParams, client apimachinery.ClientSetInterface) ModuleInterface {
	return &module{
		params: params,
		client: client,
	}
}

type module struct {
	params types.ContextParams
	client apimachinery.ClientSetInterface
}

func (m *module) hasHost(bizID int64, moduleIDS []int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		ModuleIDArr:   moduleIDS,
	}
	rsp, err := m.client.CoreService().Host().GetHostModuleRelation(context.Background(), m.params.Header, option)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request the object controller, err: %s, rid: %s", err.Error(), m.params.ReqID)
		return false, m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module]  failed to search the host module configures, err: %s, rid: %s", rsp.ErrMsg, m.params.ReqID)
		return false, m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != len(rsp.Data.Info), nil
}

func (m *module) isValidSet(bizID int64, setID int64) (bool, error) {
	inputParam := &metadata.QueryCondition{
		Limit: metadata.SearchLimit{
			Limit: common.BKNoLimit,
		},
		SortArr: nil,
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKSetIDField: setID,
		},
	}
	rsp, err := m.client.CoreService().Instance().ReadInstance(m.params.Context, m.params.Header, common.BKInnerObjIDSet, inputParam)
	if nil != err {
		blog.Errorf("[compatiblev2-module]failed to request object controller, err: %s, rid: %s", err.Error(), m.params.ReqID)
		return false, m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module] failed to search the sets, err: %s, rid: %s", rsp.ErrMsg, m.params.ReqID)
		return false, m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}
	if rsp.Data.Count == 0 {
		blog.Errorf("[compatiblev2-module] failed to search the sets, err: not found, rid: %s, rid: %s", m.params.ReqID, m.params.ReqID)
	}

	return true, nil
}

func (m *module) isRepeated(moduleName string, excludeModuleIDS interface{}) (bool, error) {
	inputParam := &metadata.QueryCondition{
		Limit: metadata.SearchLimit{
			Limit: common.BKNoLimit,
		},
		Condition: map[string]interface{}{
			common.BKModuleNameField: moduleName,
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBNIN: excludeModuleIDS,
			},
		},
	}
	rsp, err := m.client.CoreService().Instance().ReadInstance(m.params.Context, m.params.Header, common.BKInnerObjIDModule, inputParam)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request object controller, err: %s, rid: %s", err.Error(), m.params.ReqID)
		return false, m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module] failed to search the module, err: %s, rid: %s", rsp.ErrMsg, m.params.ReqID)
		return false, m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != rsp.Data.Count, nil
}

func (m *module) UpdateMultiModule(bizID int64, moduleIDS interface{}, innerData mapstr.MapStr) error {

	// check the module name
	moduleName, err := innerData.String(common.BKModuleNameField)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to parse the module name , err: %s, rid: %s", err.Error(), m.params.ReqID)
		return m.params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	repeated, err := m.isRepeated(moduleName, moduleIDS)
	if nil != err {
		return err
	}
	if repeated {
		blog.Error("[compatiblev2-module] the module name is repeated, rid: %s", m.params.ReqID)
		return m.params.Err.Errorf(common.CCErrCommDuplicateItem, moduleName)
	}

	// update module

	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	cond.Field(common.BKModuleIDField).In(moduleIDS)

	updateParam := &metadata.UpdateOption{
		Data:      innerData,
		Condition: cond.ToMapStr(),
	}
	rsp, err := m.client.CoreService().Instance().UpdateInstance(m.params.Context, m.params.Header, common.BKInnerObjIDModule, updateParam)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request object controller, err: %s, rid: %s", err.Error(), m.params.ReqID)
		return m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module] failed to search the module, err: %s, rid: %s", rsp.ErrMsg, m.params.ReqID)
		return m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (m *module) SearchModuleByApp(query *metadata.QueryCondition) (*metadata.InstResult, error) {
	rsp, err := m.client.CoreService().Instance().ReadInstance(m.params.Context, m.params.Header, common.BKInnerObjIDModule, query)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request object controller, err: %s, rid: %s", err.Error(), m.params.ReqID)
		return nil, m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module] failed to search the module, err: %s, rid: %s", rsp.ErrMsg, m.params.ReqID)
		return nil, m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	result := &metadata.InstResult{
		Count: rsp.Data.Count,
		Info:  rsp.Data.Info,
	}
	return result, nil
}

func (m *module) SearchModuleBySetProperty(bizID int64, cond condition.Condition) (*metadata.InstResult, error) {
	inputParam := &metadata.QueryCondition{
		Limit: metadata.SearchLimit{
			Limit: common.BKNoLimit,
		},
		Condition: cond.ToMapStr(),
	}
	rsp, err := m.client.CoreService().Instance().ReadInstance(m.params.Context, m.params.Header, common.BKInnerObjIDSet, inputParam)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request object controller, err: %s, rid: %s", err.Error(), m.params.ReqID)
		return nil, m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module] failed to search the sets, err: %s, rid: %s", rsp.ErrMsg, m.params.ReqID)
		return nil, m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	// construct the sets search condition
	setIDS := make([]int64, 0)
	for _, set := range rsp.Data.Info {

		id, err := set.Int64(common.BKSetIDField)
		if nil != err {
			blog.Errorf("[compatiblev2-module] failed to search sets, err: %s, rid: %s", err.Error(), m.params.ReqID)
			return nil, m.params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
		}
		setIDS = append(setIDS, id)
	}

	// search modules
	inputParam = &metadata.QueryCondition{
		Limit: metadata.SearchLimit{
			Limit: common.BKNoLimit,
		},
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKSetIDField: map[string]interface{}{
				common.BKDBIN: setIDS,
			},
		},
	}
	rspModule, err := m.client.CoreService().Instance().ReadInstance(m.params.Context, m.params.Header, common.BKInnerObjIDModule, inputParam)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request object controller, err: %s, rid: %s", err.Error(), m.params.ReqID)
		return nil, m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rspModule.Result {
		blog.Errorf("[compatiblev2-module] failed to search the sets, err: %s, rid: %s", rspModule.ErrMsg, m.params.ReqID)
		return nil, m.params.Err.New(rspModule.Code, rspModule.ErrMsg)
	}

	result := &metadata.InstResult{
		Count: rspModule.Data.Count,
		Info:  rspModule.Data.Info,
	}
	return result, nil
}

func (m *module) AddMultiModule(bizID, setID int64, moduleNames []string, data mapstr.MapStr) error {

	ok, err := m.isValidSet(bizID, setID)
	if nil != err {
		return err
	}

	if !ok {
		blog.Errorf("[compatiblev2-module] the set id (%d) is invalid, rid: %s", setID, m.params.ReqID)
		return m.params.Err.New(common.CCErrCommParamsIsInvalid, "the set id is invalid")
	}

	for _, moduleName := range moduleNames {

		if len(moduleName) > 24 { // the module name lengh limit max
			blog.Errorf("[compatiblev2-module] the module name is over the max length limit, rid: %s", m.params.ReqID)
			return m.params.Err.New(common.CCErrCommParamsIsInvalid, "the module name is over max limit")
		}

		exists, err := m.isRepeated(moduleName, []int64{})
		if nil != err {
			blog.Errorf("[compatiblev2-module] failed to check the module name(%s), err: %s , rid: %s", moduleNames, err.Error(), m.params.ReqID)
			return err
		}

		if exists {
			blog.Errorf("[compatiblev2-module] the module name (%s) is duplicated , rid: %s", moduleName, m.params.ReqID)
			return m.params.Err.Errorf(common.CCErrCommDuplicateItem, moduleName)
		}

		data.Set(common.BKModuleNameField, moduleName)
		data.Set(common.BKAppIDField, bizID)
		data.Set(common.BKDefaultField, 0)
		data.Set(common.BKInstParentStr, setID)

		createParam := &metadata.CreateModelInstance{
			Data: data,
		}
		rsp, err := m.client.CoreService().Instance().CreateInstance(m.params.Context, m.params.Header, common.BKInnerObjIDModule, createParam)
		if nil != err {
			blog.Errorf("[compatiblev2-module] failed to request object controller, err: %s, rid: %s", err.Error(), m.params.ReqID)
			return m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("[compatiblev2-module] failed to add the module name, err: %s, rid: %s", rsp.ErrMsg, m.params.ReqID)
			return m.params.Err.New(rsp.Code, rsp.ErrMsg)
		}
	}

	return nil
}

func (m *module) DeleteMultiModule(bizID int64, moduleIDS []int64) error {

	exists, err := m.hasHost(bizID, moduleIDS)
	if nil != err {
		return err
	}

	if exists {
		blog.Errorf("[compatiblev2-module] failed to delete  the modues(%#v) for the business(%d), which has some hosts , rid: %s", moduleIDS, bizID, m.params.ReqID)
		return m.params.Err.Error(common.CCErrTopoHasHost)
	}

	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	cond.Field(common.BKModuleIDField).In(moduleIDS)

	deleteParam := &metadata.DeleteOption{
		Condition: cond.ToMapStr(),
	}
	rsp, err := m.client.CoreService().Instance().DeleteInstance(m.params.Context, m.params.Header, common.BKInnerObjIDModule, deleteParam)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request object controller, err: %s, rid: %s", err.Error(), m.params.ReqID)
		return m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module] failed to delete the modules(%#v), err: %s, rid: %s", moduleIDS, rsp.ErrMsg, m.params.ReqID)
		return m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}
