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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateObjectBatch batch to create some objects
func (s *topoService) CreateObjectBatch(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	return s.core.ObjectOperation().CreateObjectBatch(params, data)
}

// SearchObjectBatch batch to search some objects
func (s *topoService) SearchObjectBatch(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	return s.core.ObjectOperation().FindObjectBatch(params, data)
}

// CreateObject create a new object
func (s *topoService) CreateObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	rsp, err := s.core.ObjectOperation().CreateObject(params, data)
	if nil != err {
		return nil, err
	}

	return rsp.ToMapStr()
}

// SearchObject search some objects by condition
func (s *topoService) SearchObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()

	if err := cond.Parse(data); nil != err {
		return nil, err
	}

	return s.core.ObjectOperation().FindObject(params, cond)
}

// SearchObjectTopo search the object topo
func (s *topoService) SearchObjectTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	cond := condition.CreateCondition()
	err := cond.Parse(data)
	if nil != err {
		return nil, params.Err.New(common.CCErrTopoObjectSelectFailed, err.Error())
	}

	return s.core.ObjectOperation().FindObjectTopo(params, cond)
}

// UpdateObject update the object
func (s *topoService) UpdateObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()

	id, err := strconv.ParseInt(pathParams("id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-obj] failed to parse the path params id(%s), error info is %s ", pathParams("id"), err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "object id")
	}

	err = s.core.ObjectOperation().UpdateObject(params, data, id, cond)
	return nil, err
}

// DeleteObject delete the object
func (s *topoService) DeleteObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()

	paramPath := frtypes.MapStr{}
	paramPath.Set("id", pathParams("id"))
	id, err := paramPath.Int64("id")
	if nil != err {
		blog.Errorf("[api-obj] failed to parse the path params id(%s), error info is %s ", pathParams("id"), err.Error())
		return nil, err
	}

	err = s.core.ObjectOperation().DeleteObject(params, id, cond)
	return nil, err
}
