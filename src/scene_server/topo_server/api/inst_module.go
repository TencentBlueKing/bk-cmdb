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

package api

import (
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/types"
)

func init() {
	apiInst.initFuncs = append(apiInst.initFuncs, apiInst.initModule)
}

func (cli *topoAPI) initModule() {
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/module/{app_id}/{set_id}", HandlerFunc: cli.CreateModule})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/module/{app_id}/{set_id}/{module_id}", HandlerFunc: cli.DeleteModule})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/module/{app_id}/{set_id}/{module_id}", HandlerFunc: cli.UpdateModule})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/module/search/{owner_id}/{app_id}/{set_id}", HandlerFunc: cli.SearchModule})

}

// CreateModule create a new module
func (cli *topoAPI) CreateModule(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("CreateModule")
	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID).Field(common.BKObjIDField).Eq(common.BKInnerObjIDModule)

	objItems, err := cli.core.ObjectOperation().FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the set, %s", err.Error())
		return nil, err
	}

	data.Set(common.BKAppIDField, pathParams("app_id"))
	data.Set(common.BKSetIDField, pathParams("set_id"))

	for _, item := range objItems {
		return cli.core.InstOperation().CreateInst(params, item, data) // should only one item
	}

	return nil, nil
}

// DeleteModule delete the module
func (cli *topoAPI) DeleteModule(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("DeleteModule")
	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDModule)
	cond.Field(common.BKAppIDField).Eq(pathParams("app_id"))
	cond.Field(common.BKSetIDField).Eq(pathParams("set_id"))
	cond.Field(common.BKModuleIDField).Eq(pathParams("module_id"))

	objItems, err := cli.core.ObjectOperation().FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	for _, item := range objItems {
		if err = cli.core.InstOperation().DeleteInst(params, item, cond); nil != err {
			return nil, err
		}
	}

	return nil, err
}

// UpdateModule update the module
func (cli *topoAPI) UpdateModule(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("UpdateModule")
	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDModule)
	cond.Field(common.BKAppIDField).Eq(pathParams("app_id"))
	cond.Field(common.BKModuleIDField).Eq(pathParams("module_id"))

	objItems, err := cli.core.ObjectOperation().FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the module, %s", err.Error())
		return nil, err
	}

	data.Set(common.BKAppIDField, pathParams("app_id"))
	data.Set(common.BKSetIDField, pathParams("set_id"))
	data.Set(common.BKModuleIDField, pathParams("module_id"))

	for _, item := range objItems {
		if err = cli.core.InstOperation().UpdateInst(params, data, item, cond); nil != err {
			return nil, err
		}
	}

	return nil, err
}

// SearchModule search the modules
func (cli *topoAPI) SearchModule(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("SearchModule")
	// {owner_id}/{app_id}/{set_id}

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDModule)
	cond.Field(common.BKAppIDField).Eq(pathParams("app_id"))
	cond.Field(common.BKSetIDField).Eq(pathParams("set_id"))

	objItems, err := cli.core.ObjectOperation().FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	count := 0
	instRst := make([]inst.Inst, 0)
	queryCond := &metadata.QueryInput{}
	for _, objItem := range objItems {

		cnt, instItems, err := cli.core.InstOperation().FindInst(params, objItem, queryCond)
		if nil != err {
			blog.Errorf("[api-business] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
			return nil, err
		}
		count = count + cnt
		instRst = append(instRst, instItems...)
	}

	result := frtypes.MapStr{}
	result.Set("count", count)
	result.Set("info", instRst)

	return result, nil
}
