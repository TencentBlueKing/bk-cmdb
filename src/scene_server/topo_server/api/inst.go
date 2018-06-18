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
	"configcenter/src/common/metadata"
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/inst"
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
func (cli *topoAPI) CreateInst(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	// /inst/{owner_id}/{obj_id}

	objID := pathParams("obj_id")

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(objID)

	objItems, err := cli.core.ObjectOperation().FindObject(params, cond)

	if nil != err {
		blog.Errorf("failed to search the %s, %s", objID, err.Error())
		return nil, err
	}

	data.Set(common.BKAppIDField, params.Header.OwnerID)
	data.Set(common.BKObjIDField, objID)

	for _, item := range objItems {

		setInst, err := cli.core.InstOperation().CreateInst(params, item, data)
		if nil != err {
			blog.Errorf("failed to create a new %s, %s", objID, err.Error())
			return nil, err
		}

		return setInst.ToMapStr(), nil // only one item
	}

	return nil, nil
}

// DeleteInst delete the inst
func (cli *topoAPI) DeleteInst(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(pathParams("obj_id"))

	objs, err := cli.core.ObjectOperation().FindObject(params, cond)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	innerCond := condition.CreateCondition()
	paramPath := frtypes.MapStr{}
	paramPath.Set("inst_id", pathParams("inst_id"))
	id, err := paramPath.Int64("inst_id")
	if nil != err {
		blog.Errorf("[api-inst] failed to parse the path params id(%s), error info is %s ", pathParams("inst_id"), err.Error())
		return nil, err
	}
	innerCond.Field(common.BKInstIDField).Eq(id)
	for _, objItem := range objs {
		err = cli.core.InstOperation().DeleteInst(params, objItem, innerCond)
		if nil != err {
			blog.Errorf("[api-inst] failed to delete the object(%s) inst (%s), error info is %s", objItem.GetID(), pathParams("inst_id"), err.Error())
			return nil, err
		}
	}

	return nil, err
}

// UpdateInst update the inst
func (cli *topoAPI) UpdateInst(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	// /inst/{owner_id}/{obj_id}/{inst_id}

	objID := pathParams("obj_id")

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(objID)

	objs, err := cli.core.ObjectOperation().FindObject(params, cond)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	innerCond := condition.CreateCondition()
	paramPath := frtypes.MapStr{}
	paramPath.Set("inst_id", pathParams("inst_id"))
	id, err := paramPath.Int64("inst_id")
	if nil != err {
		blog.Errorf("[api-inst] failed to parse the path params id(%s), error info is %s ", pathParams("inst_id"), err.Error())
		return nil, err
	}
	innerCond.Field(common.BKInstIDField).Eq(id)
	for _, objItem := range objs {
		err = cli.core.InstOperation().UpdateInst(params, data, objItem, innerCond)
		if nil != err {
			blog.Errorf("[api-inst] failed to update the object(%s) inst (%s),the data (%#v), error info is %s", objItem.GetID(), pathParams("inst_id"), data, err.Error())
			return nil, err
		}
	}

	return nil, err
}

// SearchInst search the inst
func (cli *topoAPI) SearchInst(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("SearchInst")
	// /inst/search/{owner_id}/{obj_id}

	objID := pathParams("obj_id")

	// query the objects
	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(objID)

	objs, err := cli.core.ObjectOperation().FindObject(params, cond)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	// construct the query inst condition
	count := 0
	instRst := make([]inst.Inst, 0)
	queryCond := &metadata.QueryInput{}

	if err := data.MarshalJSONInto(queryCond); nil != err {
		blog.Errorf("[api-inst] failed to parse the data and the condition, the input (%#v), error info is %s", data, err.Error())
		return nil, err
	}

	innerQueryCond, err := frtypes.NewFromInterface(queryCond.Condition)
	if nil != err {
		blog.Errorf("[api-inst] failed to parse the condition, %s", err.Error())
		return nil, err
	}

	if err := cond.Parse(innerQueryCond); nil != err {
		blog.Errorf("[api-inst] failed to parse the condition(%#v)", innerQueryCond)
		return nil, err
	}
	queryCond.Condition = cond.ToMapStr()

	fmt.Println("the query condition:", queryCond)

	// query insts
	for _, objItem := range objs {

		cnt, instItems, err := cli.core.InstOperation().FindInst(params, objItem, queryCond)
		if nil != err {
			blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
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

// SearchInstAndAssociationDetail search the inst with association details
func (cli *topoAPI) SearchInstAndAssociationDetail(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("SearchInstAndAssociationDetail")
	// /inst/search/owner/{owner_id}/object/{obj_id}/detail

	objID := pathParams("obj_id")

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(objID)

	objs, err := cli.core.ObjectOperation().FindObject(params, cond)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	count := 0
	instRst := make([]inst.Inst, 0)
	queryCond := &metadata.QueryInput{}
	if err := data.MarshalJSONInto(queryCond); nil != err {
		blog.Errorf("[api-inst] failed to parse the data and the condition, the input (%#v), error info is %s", data, err.Error())
		return nil, err
	}

	for _, objItem := range objs {

		cnt, instItems, err := cli.core.InstOperation().FindInst(params, objItem, queryCond)
		if nil != err {
			blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
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

// SearchInstByObject search the inst of the object
func (cli *topoAPI) SearchInstByObject(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	// /inst/search/owner/{owner_id}/object/{obj_id}

	objID := pathParams("obj_id")

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(objID)

	objs, err := cli.core.ObjectOperation().FindObject(params, cond)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	count := 0
	instRst := make([]inst.Inst, 0)
	queryCond := &metadata.QueryInput{}
	if err := data.MarshalJSONInto(queryCond); nil != err {
		blog.Errorf("[api-inst] failed to parse the data and the condition, the input (%#v), error info is %s", data, err.Error())
		return nil, err
	}

	for _, objItem := range objs {

		cnt, instItems, err := cli.core.InstOperation().FindInst(params, objItem, queryCond)
		if nil != err {
			blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
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

// SearchInstByAssociation search inst by the association inst
func (cli *topoAPI) SearchInstByAssociation(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("SearchInstByAssociation")
	// /inst/association/search/owner/{owner_id}/object/{obj_id}

	objID := pathParams("obj_id")

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(objID)

	objs, err := cli.core.ObjectOperation().FindObject(params, cond)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	count := 0
	instRst := make([]inst.Inst, 0)
	queryCond := &metadata.QueryInput{}
	if err := data.MarshalJSONInto(queryCond); nil != err {
		blog.Errorf("[api-inst] failed to parse the data and the condition, the input (%#v), error info is %s", data, err.Error())
		return nil, err
	}

	for _, objItem := range objs {

		cnt, instItems, err := cli.core.InstOperation().FindInst(params, objItem, queryCond)
		if nil != err {
			blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
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

// SearchInstByInstID search the inst by inst ID
func (cli *topoAPI) SearchInstByInstID(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("SearchInstByInstID")
	// /inst/search/{owner_id}/{obj_id}/{inst_id}

	objID := pathParams("obj_id")

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(objID)
	cond.Field(common.BKInstIDField).Eq(pathParams("inst_id"))

	objs, err := cli.core.ObjectOperation().FindObject(params, cond)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	count := 0
	instRst := make([]inst.Inst, 0)
	queryCond := &metadata.QueryInput{}
	if err := data.MarshalJSONInto(queryCond); nil != err {
		blog.Errorf("[api-inst] failed to parse the data and the condition, the input (%#v), error info is %s", data, err.Error())
		return nil, err
	}

	for _, objItem := range objs {

		cnt, instItems, err := cli.core.InstOperation().FindInst(params, objItem, queryCond)
		if nil != err {
			blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
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

// SearchInstChildTopo search the child inst topo for a inst
func (cli *topoAPI) SearchInstChildTopo(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("SearchInstChildTopo")
	// /inst/search/topo/owner/{owner_id}/object/{object_id}/inst/{inst_id}

	objID := pathParams("object_id")

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(objID)

	objs, err := cli.core.ObjectOperation().FindObject(params, cond)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	data.Set(common.BKInstIDField, pathParams("inst_id"))

	count := 0
	instRst := make([]inst.Inst, 0)
	queryCond := &metadata.QueryInput{}

	if err := data.MarshalJSONInto(queryCond); nil != err {
		blog.Errorf("[api-inst] failed to parse the data and the condition, the input (%#v), error info is %s", data, err.Error())
		return nil, err
	}

	for _, objItem := range objs {

		cnt, instItems, err := cli.core.InstOperation().FindInst(params, objItem, queryCond)
		if nil != err {
			blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
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

// SearchInstTopo search the inst topo
func (cli *topoAPI) SearchInstTopo(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("SearchInstTopo")
	// /inst/association/topo/search/owner/{owner_id}/object/{object_id}/inst/{inst_id}

	objID := pathParams("object_id")

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(objID)

	return nil, nil
}
