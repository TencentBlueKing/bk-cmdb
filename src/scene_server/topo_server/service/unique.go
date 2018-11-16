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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateObjectUnique create a new object unique
func (s *topoService) CreateObjectUnique(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	request := &metadata.CreateUniqueRequest{}

	if err := data.MarshalJSONInto(request); err != nil {
		blog.Errorf("format", "args")
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}

	objectID := pathParams(common.BKObjIDField)

	id, err := s.core.UniqueOperation().Create(params, objectID, request)
	if err != nil {
		return nil, err
	}
	return metadata.CreateUniqueResult{BaseResp: metadata.SuccessBaseResp, Data: *id}, nil
}

// UpdateObjectUnique update a object unique
func (s *topoService) UpdateObjectUnique(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	request := &metadata.UpdateUniqueRequest{}

	if err := data.MarshalJSONInto(request); err != nil {
		blog.Errorf("format", "args")
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}

	objectID := pathParams(common.BKObjIDField)
	id, err := strconv.ParseUint(pathParams("id"), 10, 64)
	if err != nil {
		return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, "id")
	}

	err = s.core.UniqueOperation().Update(params, objectID, id, request)
	if err != nil {
		return nil, err
	}
	return metadata.UpdateUniqueResult{BaseResp: metadata.SuccessBaseResp}, nil
}

// DeleteObjectUnique delete a object unique
func (s *topoService) DeleteObjectUnique(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	objectID := pathParams(common.BKObjIDField)
	id, err := strconv.ParseUint(pathParams("id"), 10, 64)
	if err != nil {
		return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, "id")
	}

	err = s.core.UniqueOperation().Delete(params, objectID, id)
	if err != nil {
		return nil, err
	}
	return metadata.DeleteUniqueResult{BaseResp: metadata.SuccessBaseResp}, nil
}

// SearchObjectUnique search object uniques
func (s *topoService) SearchObjectUnique(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	objectID := pathParams(common.BKObjIDField)
	uniques, err := s.core.UniqueOperation().Search(params, objectID)
	if err != nil {
		return nil, err
	}
	return metadata.SearchUniqueResult{BaseResp: metadata.SuccessBaseResp, Data: uniques}, nil
}
