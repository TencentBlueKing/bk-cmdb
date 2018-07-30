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
	SearchModuleByApp(query *metadata.QueryInput) (*metadata.InstResult, error)
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
	cond := map[string][]int64{
		"ApplicationID": []int64{bizID},
		"ModuleID":      moduleIDS,
	}

	rsp, err := m.client.HostController().Module().GetModulesHostConfig(context.Background(), m.params.Header, cond)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request the object controller, error info is %s", err.Error())
		return false, m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module]  failed to search the host module configures, error info is %s", err.Error())
		return false, m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != len(rsp.Data), nil
}
func (m *module) isValidSet(bizID, setID int64) (bool, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	cond.Field(common.BKSetIDField).Eq(setID)

	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()

	rsp, err := m.client.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDSet, m.params.Header, query)
	if nil != err {
		blog.Errorf("[compatiblev2-module]failed to request object controller, error info is %s", err.Error())
		return false, m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module] failed to search the sets, error info is %s", rsp.ErrMsg)
		return false, m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return true, nil
}

func (m *module) isRepeated(moduleName string, excludeModuleIDS interface{}) (bool, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKModuleNameField).Eq(moduleName)
	cond.Field(common.BKModuleIDField).NotIn(excludeModuleIDS)

	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()

	rsp, err := m.client.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, m.params.Header, query)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request object controller, error info is %s", err.Error())
		return false, m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module] failed to search the module, error info is %s", rsp.ErrMsg)
		return false, m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != rsp.Data.Count, nil
}

func (m *module) UpdateMultiModule(bizID int64, moduleIDS interface{}, innerData mapstr.MapStr) error {

	// check the module name
	moduleName, err := innerData.String(common.BKModuleNameField)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to parse the module name , error info is %s", err.Error())
		return m.params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	repeated, err := m.isRepeated(moduleName, moduleIDS)
	if nil != err {
		return err
	}
	if repeated {
		blog.Error("[compatiblev2-module] the module name is repeated")
		return m.params.Err.Errorf(common.CCErrCommDuplicateItem, moduleName)
	}

	// update module

	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	cond.Field(common.BKModuleIDField).In(moduleIDS)

	updateData := mapstr.New()
	updateData.Set("condition", cond.ToMapStr())
	updateData.Set("data", innerData)
	rsp, err := m.client.ObjectController().Instance().UpdateObject(context.Background(), common.BKInnerObjIDModule, m.params.Header, updateData)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request object controller, error info is %s", err.Error())
		return m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module] failed to search the module, error info is %s", rsp.ErrMsg)
		return m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}
func (m *module) SearchModuleByApp(query *metadata.QueryInput) (*metadata.InstResult, error) {

	rsp, err := m.client.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, m.params.Header, query)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request object controller, error info is %s", err.Error())
		return nil, m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module] failed to search the module, error info is %s", rsp.ErrMsg)
		return nil, m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return &rsp.Data, nil
}
func (m *module) SearchModuleBySetProperty(bizID int64, cond condition.Condition) (*metadata.InstResult, error) {

	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit
	//fmt.Println("cond:", cond.ToMapStr())
	// search sets
	rsp, err := m.client.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDSet, m.params.Header, query)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request object controller, error info is %s", err.Error())
		return nil, m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module] failed to search the sets, error info is %s", rsp.ErrMsg)
		return nil, m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	// construct the sets search condition
	setIDS := make([]int64, 0)
	for _, set := range rsp.Data.Info {

		id, err := set.Int64(common.BKSetIDField)
		if nil != err {
			blog.Errorf("[compatiblev2-module] failed to search sets, error info is %s", err.Error())
			return nil, m.params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
		}
		setIDS = append(setIDS, id)
	}

	// search modules
	cond = condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	cond.Field(common.BKSetIDField).In(setIDS)
	query.Condition = cond.ToMapStr()

	rspModule, err := m.client.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, m.params.Header, query)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request object controller, error info is %s", err.Error())
		return nil, m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rspModule.Result {
		blog.Errorf("[compatiblev2-module] failed to search the sets, error info is %s", rspModule.ErrMsg)
		return nil, m.params.Err.New(rspModule.Code, rspModule.ErrMsg)
	}

	return &rspModule.Data, nil
}
func (m *module) AddMultiModule(bizID, setID int64, moduleNames []string, data mapstr.MapStr) error {

	ok, err := m.isValidSet(bizID, setID)
	if nil != err {
		return err
	}

	if !ok {
		blog.Errorf("[compatiblev2-module] the set id (%d) is invalid", setID)
		return m.params.Err.New(common.CCErrCommParamsIsInvalid, "the set id is invalid")
	}

	for _, moduleName := range moduleNames {

		if len(moduleName) > 24 { // the module name lengh limit max
			blog.Errorf("[compatiblev2-module] the module name is over the max length limit")
			return m.params.Err.New(common.CCErrCommParamsIsInvalid, "the module name is over max limit")
		}

		exists, err := m.isRepeated(moduleName, []int64{})
		if nil != err {
			blog.Errorf("[compatiblev2-module] failed to check the module name(%s), error info is %s ", moduleNames, err.Error())
			return err
		}

		if exists {
			blog.Errorf("[compatiblev2-module] the module name (%s) is duplicated ", moduleName)
			return m.params.Err.Errorf(common.CCErrCommDuplicateItem, moduleName)
		}

		data.Set(common.BKModuleNameField, moduleName)
		data.Set(common.BKAppIDField, bizID)
		data.Set(common.BKDefaultField, 0)
		data.Set(common.BKOwnerIDField, common.BKDefaultOwnerID)

		rsp, err := m.client.ObjectController().Instance().CreateObject(context.Background(), common.BKInnerObjIDModule, m.params.Header, data)
		if nil != err {
			blog.Errorf("[compatiblev2-module] failed to request object controller, error info is %s", err.Error())
			return m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("[compatiblev2-module] failed to add the module name, error info is %s", rsp.ErrMsg)
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
		blog.Errorf("[compatiblev2-module] failed to delete  the modues(%#v) for the business(%d), which has some hosts ", moduleIDS, bizID)
		return m.params.Err.Error(common.CCErrTopoHasHost)
	}

	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	cond.Field(common.BKModuleIDField).In(moduleIDS)
	//fmt.Println("cond:", cond.ToMapStr())
	rsp, err := m.client.ObjectController().Instance().DelObject(context.Background(), common.BKInnerObjIDModule, m.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request object controller, error info is %s", err.Error())
		return m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module] failed to delete the modules(%#v), error info is %s", moduleIDS, rsp.ErrMsg)
		return m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}
