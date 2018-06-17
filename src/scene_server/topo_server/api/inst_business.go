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
	apiInst.initFuncs = append(apiInst.initFuncs, apiInst.initBusiness)
}

func (cli *topoAPI) initBusiness() {
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/app/{owner_id}", HandlerFunc: cli.CreateBusiness})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/app/{owner_id}/{app_id}", HandlerFunc: cli.DeleteBusiness})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/app/{owner_id}/{app_id}", HandlerFunc: cli.UpdateBusiness})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/app/status/{flag}/{owner_id}/{app_id}", HandlerFunc: cli.UpdateBusinessStatus})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/app/search/{owner_id}", HandlerFunc: cli.SearchBusiness})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/app/default/{owner_id}/search", HandlerFunc: cli.SearchDefaultBusiness})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/app/default/{owner_id}", HandlerFunc: cli.CreateDefaultBusiness})
}

// CreateBusiness create a new business
func (cli *topoAPI) CreateBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	fmt.Println("CreateBusiness")

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDApp)

	objItems, err := cli.core.ObjectOperation().FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	data.Set(common.BKOwnerIDField, params.Header.OwnerID)

	for _, item := range objItems {
		setInst, err := cli.core.InstOperation().CreateInst(params, item, data)
		if nil != err {
			blog.Errorf("failed to create a new business, %s", err.Error())
			return nil, err
		}

		err = setInst.Save()
		if nil != err {
			blog.Errorf("failed to create a new business, %s", err.Error())
			return nil, err
		}

		return setInst.ToMapStr() // only one item
	}

	return nil, nil
}

// DeleteBusiness delete the business
func (cli *topoAPI) DeleteBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDApp)
	cond.Field(common.BKAppIDField).Eq(pathParams("app_id"))

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

// UpdateBusiness update the business
func (cli *topoAPI) UpdateBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDModule)
	cond.Field(common.BKAppIDField).Eq(pathParams("app_id"))

	objItems, err := cli.core.ObjectOperation().FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	for _, item := range objItems {
		if err = cli.core.InstOperation().UpdateInst(params, data, item, cond); nil != err {
			return nil, err
		}
	}

	return nil, err
}

// UpdateBusinessStatus update the business status
func (cli *topoAPI) UpdateBusinessStatus(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	// /app/status/{flag}/{owner_id}/{app_id}

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDModule)
	cond.Field(common.BKAppIDField).Eq(pathParams("app_id"))

	objItems, err := cli.core.ObjectOperation().FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	data.Set("flag", pathParams("flag"))
	for _, item := range objItems {
		if err = cli.core.InstOperation().UpdateInst(params, data, item, cond); nil != err {
			return nil, err
		}
	}

	return nil, err
}

// SearchBusiness search the business by condition
func (cli *topoAPI) SearchBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	// "/app/search/{owner_id}

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDApp)

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

// SearchDefaultBusiness search the business by condition
func (cli *topoAPI) SearchDefaultBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDApp)

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

// CreateDefaultBusiness create the default business
func (cli *topoAPI) CreateDefaultBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("CreateDefaultBusiness")
	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDApp)

	objItems, err := cli.core.ObjectOperation().FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	data.Set(common.BKOwnerIDField, params.Header.OwnerID)

	for _, item := range objItems {
		setInst, err := cli.core.InstOperation().CreateInst(params, item, data)
		if nil != err {
			blog.Errorf("failed to create a new business, %s", err.Error())
			return nil, err
		}

		err = setInst.Save()
		if nil != err {
			blog.Errorf("failed to create a new business, %s", err.Error())
			return nil, err
		}

		return setInst.ToMapStr() // only one item
	}

	return nil, nil
}
