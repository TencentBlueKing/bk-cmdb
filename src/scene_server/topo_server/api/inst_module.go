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
	"net/http"

	"configcenter/src/common"
	frcommon "configcenter/src/common"
	"configcenter/src/common/blog"
	frtypes "configcenter/src/common/types"
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
func (cli *topoAPI) CreateModule(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID).Field(common.BKObjIDField).Eq(common.BKInnerObjIDModule)

	objItems, err := cli.core.FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the set, %s", err.Error())
		return nil, err
	}

	data.Set(common.BKAppIDField, pathParams("app_id"))
	data.Set(common.BKSetIDField, pathParams("set_id"))

	for _, item := range objItems {
		setInst, err := cli.core.CreateInst(params, item, data)
		if nil != err {
			blog.Errorf("failed to create a new set, %s", err.Error())
			return nil, err
		}

		err = setInst.Save()
		if nil != err {
			blog.Errorf("failed to create a new set, %s", err.Error())
			return nil, err
		}

		return setInst.ToMapStr() // only one item
	}

	return nil, nil
}

// DeleteModule delete the module
func (cli *topoAPI) DeleteModule(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID).
		Field(common.BKObjIDField).Eq(common.BKInnerObjIDModule).
		Field(common.BKAppIDField).Eq(pathParams("app_id")).
		Field(common.BKSetIDField).Eq(pathParams("set_id")).
		Field(common.BKModuleIDField).Eq(pathParams("module_id"))

	err := cli.core.DeleteInst(params, cond)

	return nil, err
}

// UpdateModule update the module
func (cli *topoAPI) UpdateModule(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID).
		Field(common.BKObjIDField).Eq(common.BKInnerObjIDModule).
		Field(common.BKAppIDField).Eq(pathParams("app_id")).
		Field(common.BKModuleIDField).Eq(pathParams("module_id"))

	data.Set(common.BKAppIDField, pathParams("app_id"))
	data.Set(common.BKSetIDField, pathParams("set_id"))
	data.Set(common.BKModuleIDField, pathParams("module_id"))

	err := cli.core.UpdateInst(params, data, cond)

	return nil, err
}

// SearchModule search the modules
func (cli *topoAPI) SearchModule(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	// {owner_id}/{app_id}/{set_id}

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID).
		Field(common.BKObjIDField).Eq(common.BKInnerObjIDModule).
		Field(common.BKAppIDField).Eq(pathParams("app_id")).
		Field(common.BKSetIDField).Eq(pathParams("set_id"))

	data.Set(common.BKAppIDField, pathParams("app_id"))
	data.Set(common.BKInnerObjIDSet, pathParams("set_id"))
	data.Set(common.BKOwnerIDField, pathParams("owner_id"))

	items, err := cli.core.FindInst(params, cond)
	if nil != err {
		return nil, err
	}

	results := make([]frtypes.MapStr, 0)
	for _, item := range items {
		toMapStr, err := item.ToMapStr()
		if nil != err {
			return nil, err
		}
		results = append(results, toMapStr)
	}

	resultData := frtypes.MapStr{}
	resultData.Set("data", results)
	return resultData, nil
}
