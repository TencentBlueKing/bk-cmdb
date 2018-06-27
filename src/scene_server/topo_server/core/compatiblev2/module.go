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
	"configcenter/src/common/condition"
	"context"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// ModuleInterface module interface
type ModuleInterface interface {
	UpdateMultiModule()
	SearchModuleByApp(query *metadata.QueryInput) (*metadata.InstResult, error)
	SearchModuleByProperty(bizID int64, cond condition.Condition) (*metadata.InstResult, error)
	AddMultiModule()
	DeleteMultiModule()
}

// NewModule create a module instance
func NewModule() ModuleInterface {
	return &module{}
}

type module struct {
	params types.LogicParams
	client apimachinery.ClientSetInterface
}

func (m *module) UpdateMultiModule() {

}
func (m *module) SearchModuleByApp(query *metadata.QueryInput) (*metadata.InstResult, error) {

	rsp, err := m.client.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, m.params.Header.ToHeader(), query)
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
func (m *module) SearchModuleByProperty(bizID int64, cond condition.Condition) {

	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()

	// search sets
	rsp, err := m.client.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDSet, m.params.Header.ToHeader(), query)
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

	rsp, err := m.client.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, m.params.Header.ToHeader(), query)
	if nil != err {
		blog.Errorf("[compatiblev2-module] failed to request object controller, error info is %s", err.Error())
		return nil, m.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[compatiblev2-module] failed to search the sets, error info is %s", rsp.ErrMsg)
		return nil, m.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return &rsp.Data, nil
}
func (m *module) AddMultiModule() {

}
func (m *module) DeleteMultiModule() {

}
