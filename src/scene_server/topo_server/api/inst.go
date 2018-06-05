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
	apiInst.initFuncs = append(apiInst.initFuncs, apiInst.initInst)
}

func (cli *topoAPI) initInst() {
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/{owner_id}/{obj_id}", HandlerFunc: cli.CreateInst})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/inst/{owner_id}/{obj_id}/{inst_id}", HandlerFunc: cli.DeleteInst})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/inst/{owner_id}/{obj_id}/{inst_id}", HandlerFunc: cli.UpdateInst})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/search/{owner_id}/{obj_id}", HandlerFunc: cli.SearchInst})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/search/owner/{owner_id}/object/{obj_id}/detail", HandlerFunc: cli.SearchInstAndAssociationDetail})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/search/owner/{owner_id}/object/{obj_id}", HandlerFunc: cli.SearchInstByObject})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/association/search/owner/{owner_id}/object/{obj_id}", HandlerFunc: cli.SearchInstByAssociation})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/search/{owner_id}/{obj_id}/{inst_id}", HandlerFunc: cli.SearchInstByInstID})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/search/topo/owner/{owner_id}/object/{object_id}/inst/{inst_id}", HandlerFunc: cli.SearchInstChildTopo})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/inst/association/topo/search/owner/{owner_id}/object/{object_id}/inst/{inst_id}", HandlerFunc: cli.SearchInstTopo})
}

// CreateInst create a new inst
func (cli *topoAPI) CreateInst(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	// /inst/{owner_id}/{obj_id}

	objID := pathParams("obj_id")

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(objID)

	objItems, err := cli.core.FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the %s, %s", objID, err.Error())
		return nil, err
	}

	data.Set(common.BKAppIDField, params.OwnerID)
	data.Set(common.BKObjIDField, objID)

	for _, item := range objItems {
		setInst, err := cli.core.CreateInst(params, item, data)
		if nil != err {
			blog.Errorf("failed to create a new %s, %s", objID, err.Error())
			return nil, err
		}

		err = setInst.Save()
		if nil != err {
			blog.Errorf("failed to create a new %s, %s", objID, err.Error())
			return nil, err
		}

		return setInst.ToMapStr() // only one item
	}

	return nil, nil
}

// DeleteInst delete the inst
func (cli *topoAPI) DeleteInst(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(pathParams("obj_id")).
		Field(common.BKInstIDField).Eq(pathParams("inst_id"))

	err := cli.core.DeleteInst(params, cond)
	return nil, err
}

// UpdateInst update the inst
func (cli *topoAPI) UpdateInst(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	// /inst/{owner_id}/{obj_id}/{inst_id}

	objID := pathParams("obj_id")

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(objID).
		Field(common.BKInstIDField).Eq(pathParams("inst_id"))

	err := cli.core.UpdateInst(params, data, cond)
	return nil, err
}

// SearchInst search the inst
func (cli *topoAPI) SearchInst(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	// /inst/search/{owner_id}/{obj_id}

	objID := pathParams("obj_id")

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(objID)

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

// SearchInstAndAssociationDetail search the inst with association details
func (cli *topoAPI) SearchInstAndAssociationDetail(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	// /inst/search/owner/{owner_id}/object/{obj_id}/detail

	objID := pathParams("obj_id")

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(objID)

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

// SearchInstByObject search the inst of the object
func (cli *topoAPI) SearchInstByObject(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	// /inst/search/owner/{owner_id}/object/{obj_id}

	objID := pathParams("obj_id")

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(objID)

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

// SearchInstByAssociation search inst by the association inst
func (cli *topoAPI) SearchInstByAssociation(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	// /inst/association/search/owner/{owner_id}/object/{obj_id}

	objID := pathParams("obj_id")

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(objID)

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

// SearchInstByInstID search the inst by inst ID
func (cli *topoAPI) SearchInstByInstID(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {
	// /inst/search/{owner_id}/{obj_id}/{inst_id}

	objID := pathParams("obj_id")

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(objID).
		Field(common.BKInstIDField).Eq(pathParams("inst_id"))

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

// SearchInstChildTopo search the child inst topo for a inst
func (cli *topoAPI) SearchInstChildTopo(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	// /inst/search/topo/owner/{owner_id}/object/{object_id}/inst/{inst_id}

	objID := pathParams("object_id")

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(objID).
		Field(common.BKInstIDField).Eq("inst_id")

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

// SearchInstTopo search the inst topo
func (cli *topoAPI) SearchInstTopo(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (frtypes.MapStr, error) {

	// /inst/association/topo/search/owner/{owner_id}/object/{object_id}/inst/{inst_id}

	objID := pathParams("object_id")

	cond := frcommon.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.OwnerID).
		Field(common.BKObjIDField).Eq(objID).
		Field(common.BKInstIDField).Eq(pathParams("inst_id"))

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
