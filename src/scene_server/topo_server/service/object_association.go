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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateObjectAssociation create a new object association
func (s *Service) CreateObjectAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	assoc := &metadata.Association{}
	if err := data.MarshalJSONInto(assoc); err != nil {
		return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
	}
	params.MetaData = &assoc.Metadata
	association, err := s.Core.AssociationOperation().CreateCommonAssociation(params, assoc)
	if nil != err {
		return nil, err
	}

	return association, nil

}

// SearchObjectAssociation search  object association by object id
func (s *Service) SearchObjectAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	if data.Exists("condition") {
		// ATTENTION:
		// compatible with new query structures
		// the new condition format:
		// { "condition":{}}

		cond, err := data.MapStr("condition")
		if nil != err {
			blog.Errorf("search object association, failed to get the condition, error info is %s, rid: %s", err.Error(), params.ReqID)
			return nil, params.Err.New(common.CCErrCommParamsIsInvalid, err.Error())
		}

		if len(cond) == 0 {
			return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
		}

		if nil != params.MetaData {
			cond.Merge(metadata.PublicAndBizCondition(*params.MetaData))
			cond.Remove(metadata.BKMetadata)
		} else {
			cond.Merge(metadata.BizLabelNotExist)
		}

		resp, err := s.Core.AssociationOperation().SearchObject(params, &metadata.SearchAssociationObjectRequest{Condition: cond})
		if err != nil {
			blog.Errorf("search object association with cond[%v] failed, err: %v, rid: %s", cond, err, params.ReqID)
			return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !resp.Result {
			blog.Errorf("search object association with cond[%v] failed, err: %s, rid: %s", cond, resp.ErrMsg, params.ReqID)
			return nil, params.Err.New(resp.Code, resp.ErrMsg)
		}

		return resp.Data, err
	}

	objID, err := data.String(metadata.AssociationFieldObjectID)
	if err != nil {
		blog.Errorf("search object association, but get object id failed from: %v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
	}

	if len(objID) == 0 {
		return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
	}

	return s.Core.AssociationOperation().SearchObjectAssociation(params, objID)
}

// DeleteObjectAssociation delete object association
func (s *Service) DeleteObjectAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	id, err := strconv.ParseInt(pathParams("id"), 10, 64)
	if err != nil {
		blog.Errorf("delete object association failed, got a invalid object association id[%v], err: %v, rid: %s", pathParams("id"), err, params.ReqID)
		return nil, params.Err.Error(common.CCErrTopoInvalidObjectAssociationID)
	}

	if id <= 0 {
		blog.Errorf("delete object association failed, got a invalid objAsst id[%d], rid: %s", id, params.ReqID)
		return nil, params.Err.Error(common.CCErrTopoInvalidObjectAssociationID)
	}

	data.Remove(metadata.BKMetadata)
	return nil, s.Core.AssociationOperation().DeleteAssociationWithPreCheck(params, id)
}

// UpdateObjectAssociation update object association
func (s *Service) UpdateObjectAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	id, err := strconv.ParseInt(pathParams("id"), 10, 64)
	if err != nil {
		blog.Errorf("update object association, but got invalid id[%v], err: %v, rid: %s", pathParams("id"), err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
	}

	err = s.Core.AssociationOperation().UpdateAssociation(params, data, id)
	return nil, err

}

// ImportInstanceAssociation import instance  association
func (s *Service) ImportInstanceAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	objID := pathParams("bk_obj_id")
	request := new(metadata.RequestImportAssociation)
	if err := data.MarshalJSONInto(request); err != nil {
		blog.Errorf("ImportInstanceAssociation, json unmarshal error, objID:%S, err: %v, rid:%s", objID, err, params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}

	resp, err := s.Core.AssociationOperation().ImportInstAssociation(context.Background(), params, objID, request.AssociationInfoMap)
	return resp, err

}
