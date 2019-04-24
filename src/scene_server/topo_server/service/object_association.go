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
	"context"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateObjectAssociation create a new object association
func (s *topoService) CreateObjectAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	asso := &metadata.Association{}
	if err := data.MarshalJSONInto(asso); err != nil {
		return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
	}

	association, err := s.core.AssociationOperation().CreateCommonAssociation(params, asso)
	if nil != err {
		return nil, err
	}

	return association, nil

}

// SearchObjectAssociation search  object association by object id
func (s *topoService) SearchObjectAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	if data.Exists("condition") {
		// ATTENTION:
		// compatible with new query structures
		// the new condition format:
		// { "condition":{}}

		cond, err := data.MapStr("condition")
		if nil != err {
			blog.Errorf("search object association, failed to get the condition, error info is %s", err.Error())
			return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
		}

		if len(cond) == 0 {
			return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
		}

		resp, err := s.core.AssociationOperation().SearchObject(params, &metadata.SearchAssociationObjectRequest{Condition: cond})
		if err != nil {
			blog.Errorf("search object association with cond[%v] failed, err: %v", cond, err)
			return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !resp.Result {
			blog.Errorf("search object association with cond[%v] failed, err: %s", cond, resp.ErrMsg)
			return nil, params.Err.Error(resp.Code)
		}

		return resp.Data, err
	}

	objID, err := data.String(metadata.AssociationFieldObjectID)
	if err != nil {
		blog.Errorf("search object association, but get object id failed from: %v, err: %v", data, err)
		return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
	}

	if len(objID) == 0 {
		return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
	}

	return s.core.AssociationOperation().SearchObjectAssociation(params, objID)
}

// DeleteObjectAssociation delete object association
func (s *topoService) DeleteObjectAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	id, err := strconv.ParseInt(pathParams("id"), 10, 64)
	if err != nil {
		blog.Errorf("delete object association failed, got a invalid object association id[%v], err: %v", pathParams("id"), err)
		return nil, params.Err.Error(common.CCErrTopoInvalidObjectAssociationID)
	}

	if id <= 0 {
		blog.Errorf("delete object association failed, got a invalid objasst id[%d]", id)
		return nil, params.Err.Error(common.CCErrTopoInvalidObjectAssociationID)
	}

	return nil, s.core.AssociationOperation().DeleteAssociationWithPreCheck(params, id)
}

// UpdateObjectAssociation update object association
func (s *topoService) UpdateObjectAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	id, err := strconv.ParseInt(pathParams("id"), 10, 64)
	if err != nil {
		blog.Errorf("update object association, but got invalid id[%v], err: %v", pathParams("id"), err)
		return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
	}
	err = s.core.AssociationOperation().UpdateAssociation(params, data, id)
	return nil, err

}

// ImportInstanceAssociation import instance  association
func (s *topoService) ImportInstanceAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	objID := pathParams("obj_id")
	request := new(metadata.RequestImportAssociation)
	if err := data.MarshalJSONInto(request); err != nil {
		blog.Errorf("ImportInstanceAssociation, json unmarshal error, objID:%S, err: %v", objID, err)
		return nil, params.Err.New(common.CCErrCommJSONUnmarshalFailed, err.Error())
	}

	resp, err := s.core.AssociationOperation().ImportInstAssociation(context.Background(), params, objID, request.AssociationInfoMap)
	return resp, err

}
