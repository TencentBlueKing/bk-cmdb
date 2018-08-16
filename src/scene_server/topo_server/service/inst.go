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

package service

import (
	"strconv"
	"strings"

	"configcenter/src/scene_server/topo_server/core/operation"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	paraparse "configcenter/src/common/paraparse"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateInst create a new inst
func (s *topoService) CreateInst(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	// /inst/{owner_id}/{obj_id}

	objID := pathParams("obj_id")

	obj, err := s.core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("failed to search the inst, %s", err.Error())
		return nil, err
	}

	if data.Exists("BatchInfo") {
		batchInfo := new(operation.InstBatchInfo)
		data.MarshalJSONInto(batchInfo)
		setInst, err := s.core.InstOperation().CreateInstBatch(params, obj, batchInfo)
		if nil != err {
			blog.Errorf("failed to create a new %s, %s", objID, err.Error())
			return nil, err
		}
		return setInst, nil
	}

	setInst, err := s.core.InstOperation().CreateInst(params, obj, data)
	if nil != err {
		blog.Errorf("failed to create a new %s, %s", objID, err.Error())
		return nil, err
	}

	return setInst.ToMapStr(), nil
}
func (s *topoService) DeleteInsts(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	obj, err := s.core.ObjectOperation().FindSingleObject(params, pathParams("obj_id"))
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	deleteCondition := &operation.OpCondition{}
	if err := data.MarshalJSONInto(deleteCondition); nil != err {
		return nil, err
	}

	return nil, s.core.InstOperation().DeleteInstByInstID(params, obj, deleteCondition.Delete.InstID)
}

// DeleteInst delete the inst
func (s *topoService) DeleteInst(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	if "batch" == pathParams("inst_id") {
		return s.DeleteInsts(params, pathParams, queryParams, data)
	}

	instID, err := strconv.ParseInt(pathParams("inst_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-inst]failed to parse the inst id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "inst id")
	}

	obj, err := s.core.ObjectOperation().FindSingleObject(params, pathParams("obj_id"))
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	err = s.core.InstOperation().DeleteInstByInstID(params, obj, []int64{instID})
	return nil, err
}
func (s *topoService) UpdateInsts(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	objID := pathParams("obj_id")

	updateCondition := &operation.OpCondition{}
	if err := data.MarshalJSONInto(updateCondition); nil != err {
		blog.Errorf("[api-inst] failed to parse the input data(%v), error info is %s", data, err.Error())
		return nil, err
	}

	obj, err := s.core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	for _, item := range updateCondition.Update {

		cond := condition.CreateCondition()
		cond.Field(obj.GetInstIDFieldName()).Eq(item.InstID)
		err = s.core.InstOperation().UpdateInst(params, item.InstInfo, obj, cond, item.InstID)
		if nil != err {
			blog.Errorf("[api-inst] failed to update the object(%s) inst (%d),the data (%#v), error info is %s", obj.GetID(), item.InstID, data, err.Error())
			return nil, err
		}
	}

	return nil, nil
}

// UpdateInst update the inst
func (s *topoService) UpdateInst(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	// /inst/{owner_id}/{obj_id}/{inst_id}

	if "batch" == pathParams("inst_id") {
		return s.UpdateInsts(params, pathParams, queryParams, data)
	}

	objID := pathParams("obj_id")

	obj, err := s.core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	instID, err := strconv.ParseInt(pathParams("inst_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-inst]failed to parse the inst id, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "inst id")
	}

	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)
	err = s.core.InstOperation().UpdateInst(params, data, obj, cond, instID)
	if nil != err {
		blog.Errorf("[api-inst] failed to update the object(%s) inst (%s),the data (%#v), error info is %s", obj.GetID(), pathParams("inst_id"), data, err.Error())
		return nil, err
	}

	return nil, err
}

// SearchInst search the inst
func (s *topoService) SearchInsts(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	objID := pathParams("obj_id")

	obj, err := s.core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	// construct the query inst condition
	queryCond := &paraparse.SearchParams{}
	if err := data.MarshalJSONInto(queryCond); nil != err {
		blog.Errorf("[api-inst] failed to parse the data and the condition, the input (%#v), error info is %s", data, err.Error())
		return nil, err
	}
	page := metadata.ParsePage(queryCond.Page)
	query := &metadata.QueryInput{}
	query.Condition = queryCond.Condition
	query.Fields = strings.Join(queryCond.Fields, ",")
	query.Limit = page.Limit
	query.Sort = page.Sort
	query.Start = page.Start

	cnt, instItems, err := s.core.InstOperation().FindInst(params, obj, query, false)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	result := frtypes.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
	return result, nil
}

// SearchInstAndAssociationDetail search the inst with association details
func (s *topoService) SearchInstAndAssociationDetail(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	//fmt.Println("SearchInstAndAssociationDetail")
	// /inst/search/owner/{owner_id}/object/{obj_id}/detail

	objID := pathParams("obj_id")

	obj, err := s.core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	// construct the query inst condition

	queryCond := &paraparse.SearchParams{}
	if err := data.MarshalJSONInto(queryCond); nil != err {
		blog.Errorf("[api-inst] failed to parse the data and the condition, the input (%#v), error info is %s", data, err.Error())
		return nil, err
	}
	page := metadata.ParsePage(queryCond.Page)
	query := &metadata.QueryInput{}
	query.Condition = queryCond.Condition
	query.Fields = strings.Join(queryCond.Fields, ",")
	query.Limit = page.Limit
	query.Sort = page.Sort
	query.Start = page.Start

	cnt, instItems, err := s.core.InstOperation().FindInst(params, obj, query, true)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	result := frtypes.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
	return result, nil
}

// SearchInstByObject search the inst of the object
func (s *topoService) SearchInstByObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	// /inst/search/owner/{owner_id}/object/{obj_id}

	objID := pathParams("obj_id")
	obj, err := s.core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	queryCond := &paraparse.SearchParams{}
	if err := data.MarshalJSONInto(queryCond); nil != err {
		blog.Errorf("[api-inst] failed to parse the data and the condition, the input (%#v), error info is %s", data, err.Error())
		return nil, err
	}
	page := metadata.ParsePage(queryCond.Page)
	query := &metadata.QueryInput{}
	query.Condition = queryCond.Condition
	query.Fields = strings.Join(queryCond.Fields, ",")
	query.Limit = page.Limit
	query.Sort = page.Sort
	query.Start = page.Start
	cnt, instItems, err := s.core.InstOperation().FindInst(params, obj, query, false)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	result := frtypes.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)

	return result, nil

}

// SearchInstByAssociation search inst by the association inst
func (s *topoService) SearchInstByAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	// fmt.Println("SearchInstByAssociation")
	// /inst/association/search/owner/{owner_id}/object/{obj_id}

	objID := pathParams("obj_id")
	obj, err := s.core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	cnt, instItems, err := s.core.InstOperation().FindInstByAssociationInst(params, obj, data)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	result := frtypes.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
	return result, nil
}

// SearchInstByInstID search the inst by inst ID
func (s *topoService) SearchInstByInstID(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	// /inst/search/{owner_id}/{obj_id}/{inst_id}

	objID := pathParams("obj_id")

	instID, err := strconv.ParseInt(pathParams("inst_id"), 10, 64)
	if nil != err {
		return nil, params.Err.New(common.CCErrTopoInstSelectFailed, err.Error())
	}

	obj, err := s.core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)
	queryCond := &metadata.QueryInput{}
	queryCond.Condition = cond.ToMapStr()

	cnt, instItems, err := s.core.InstOperation().FindInst(params, obj, queryCond, false)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	result := frtypes.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)

	return result, nil
}

// SearchInstChildTopo search the child inst topo for a inst
func (s *topoService) SearchInstChildTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	//fmt.Println("SearchInstChildTopo")
	// /inst/search/topo/owner/{owner_id}/object/{object_id}/inst/{inst_id}

	objID := pathParams("object_id")

	obj, err := s.core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	instID, err := strconv.ParseInt(pathParams("inst_id"), 10, 64)
	if nil != err {
		return nil, err
	}

	query := &metadata.QueryInput{}

	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)

	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit

	_, instItems, err := s.core.InstOperation().FindInstChildTopo(params, obj, instID, query)
	return instItems, err

}

// SearchInstTopo search the inst topo
func (s *topoService) SearchInstTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	//fmt.Println("SearchInstTopo")
	// /inst/association/topo/search/owner/{owner_id}/object/{object_id}/inst/{inst_id}

	objID := pathParams("object_id")

	obj, err := s.core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the objects(%s), error info is %s", pathParams("obj_id"), err.Error())
		return nil, err
	}

	instID, err := strconv.ParseInt(pathParams("inst_id"), 10, 64)
	if nil != err {
		return nil, err
	}

	query := &metadata.QueryInput{}

	cond := condition.CreateCondition()
	cond.Field(obj.GetInstIDFieldName()).Eq(instID)
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)

	query.Condition = cond.ToMapStr()
	query.Limit = common.BKNoLimit

	_, instItems, err := s.core.InstOperation().FindInstTopo(params, obj, instID, query)
	return instItems, err
}
