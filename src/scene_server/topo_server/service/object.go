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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateObjectBatch batch to create some objects
func (s *Service) CreateObjectBatch(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	data.Remove(metadata.BKMetadata)
	return s.Core.ObjectOperation().CreateObjectBatch(params, data)
}

// SearchObjectBatch batch to search some objects
func (s *Service) SearchObjectBatch(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	data.Remove(metadata.BKMetadata)
	return s.Core.ObjectOperation().FindObjectBatch(params, data)
}

// CreateObject create a new object
func (s *Service) CreateObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	rsp, err := s.Core.ObjectOperation().CreateObject(params, false, data)
	if nil != err {
		return nil, err
	}

	return rsp.ToMapStr()
}

// SearchObject search some objects by condition
func (s *Service) SearchObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	cond := condition.CreateCondition()
	if err := cond.Parse(data); nil != err {
		return nil, err
	}

	return s.Core.ObjectOperation().FindObject(params, cond)
}

// SearchObjectTopo search the object topo
func (s *Service) SearchObjectTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	cond := condition.CreateCondition()
	err := cond.Parse(data)
	if nil != err {
		return nil, params.Err.New(common.CCErrTopoObjectSelectFailed, err.Error())
	}

	return s.Core.ObjectOperation().FindObjectTopo(params, cond)
}

// UpdateObject update the object
func (s *Service) UpdateObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	idStr := pathParams(common.BKFieldID)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		blog.Errorf("[api-obj] failed to parse the path params id(%s), error info is %s , rid: %s", idStr, err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, common.BKFieldID)
	}
	err = s.Core.ObjectOperation().UpdateObject(params, data, id)
	return nil, err
}

// DeleteObject delete the object
func (s *Service) DeleteObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	idStr := pathParams(common.BKFieldID)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		blog.Errorf("[api-obj] failed to parse the path params id(%s), error info is %s , rid: %s", idStr, err.Error(), params.ReqID)
		return nil, params.Err.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID)
	}

	err = s.Core.ObjectOperation().DeleteObject(params, id, true)
	return nil, err
}

// GetModelStatistics 用于统计各个模型的实例数(Web页面展示需要)
func (s *Service) GetModelStatistics(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	result, err := s.Engine.CoreAPI.CoreService().Model().GetModelStatistics(params.Context, params.Header)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}
	return result.Data, err
}
