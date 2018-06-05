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
	"configcenter/src/common/blog"
	frcommon "configcenter/src/framework/common"
	frtypes "configcenter/src/framework/core/types"
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
func (cli *topoAPI) CreateBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(common.BKInnerObjIDApp)

	objItems, err := cli.core.FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	data.Set(common.BKOwnerIDField, params.OwnerID)

	for _, item := range objItems {
		setInst, err := cli.core.CreateInst(params, item, data)
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
func (cli *topoAPI) DeleteBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(common.BKInnerObjIDApp).
		Field(common.BKAppIDField).Eq(pathParams("app_id"))

	err := cli.core.DeleteInst(params, cond)

	return nil, err
}

// UpdateBusiness update the business
func (cli *topoAPI) UpdateBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(common.BKInnerObjIDModule).
		Field(common.BKAppIDField).Eq(pathParams("app_id"))

	data.Set(common.BKAppIDField, pathParams("app_id"))
	err := cli.core.UpdateInst(params, data, cond)

	return nil, err
}

// UpdateBusinessStatus update the business status
func (cli *topoAPI) UpdateBusinessStatus(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	// /app/status/{flag}/{owner_id}/{app_id}

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(common.BKInnerObjIDModule).
		Field(common.BKAppIDField).Eq(pathParams("app_id"))

	data.Set("flag", pathParams("flag"))
	err := cli.core.UpdateInst(params, data, cond)

	return nil, err
}

// SearchBusiness search the business by condition
func (cli *topoAPI) SearchBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	// "/app/search/{owner_id}

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(common.BKInnerObjIDApp)

	data.Set(common.BKOwnerIDField, params.OwnerID)

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

// SearchDefaultBusiness search the business by condition
func (cli *topoAPI) SearchDefaultBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(common.BKInnerObjIDApp)

	data.Set(common.BKOwnerIDField, params.OwnerID)

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

// CreateDefaultBusiness create the default business
func (cli *topoAPI) CreateDefaultBusiness(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(common.BKInnerObjIDApp)

	objItems, err := cli.core.FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the business, %s", err.Error())
		return nil, err
	}

	data.Set(common.BKOwnerIDField, params.OwnerID)

	for _, item := range objItems {
		setInst, err := cli.core.CreateInst(params, item, data)
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
