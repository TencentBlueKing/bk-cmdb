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

	err := s.core.AssociationOperation().CreateCommonAssociation(params, asso)
	if nil != err {
		return nil, err
	}

	return asso.ToMapStr(), nil

}

// SearchObjectAssociation search  object association by object id
func (s *topoService) SearchObjectAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	objId, err := data.String(metadata.AssociationFieldObjectID)
	if err != nil {
		blog.Errorf("search object association, but get object id failed from: %v, err: %v", data, err)
		return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
	}

	return s.core.AssociationOperation().SearchObjectAssociation(params, objId)
}

// DeleteObjectAssociation delete object association
func (s *topoService) DeleteObjectAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	id, err := strconv.ParseInt(pathParams("id"), 10, 64)
	if err != nil {
		blog.Errorf("delete object association failed, got a invalid objasst id[%v], err: %v", pathParams("id"), err)
		return nil, params.Err.Error(common.CCErrTopoInvalidObjectAssociaitonID)
	}

	if id <= 0 {
		blog.Errorf("delete object association failed, got a invalid objasst id[%d]", id)
		return nil, params.Err.Error(common.CCErrTopoInvalidObjectAssociaitonID)
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
