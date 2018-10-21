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
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
	"strconv"
)

// CreateObjectAssociation create a new object association
func (s *topoService) CreateObjectAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	asso := &metadata.Association{}
	if err := data.MarshalJSONInto(asso); err != nil {
		return nil, err
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
		blog.Errorf("failed to get object id, error info is %s", err.Error())
		return nil, err
	}

	return s.core.AssociationOperation().SearchObjectAssociation(params, objId)
}

// DeleteObjectAssociation delete object association
func (s *topoService) DeleteObjectAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	cond.Field(metadata.AssociationFieldAssociationId).Eq(pathParams("id"))
	err := s.core.AssociationOperation().DeleteAssociation(params, cond)
	return nil, err

}

// UpdateObjectAssociation update object association
func (s *topoService) UpdateObjectAssociation(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	id, _ := strconv.ParseInt(pathParams("id"), 10, 64)
	cond := condition.CreateCondition()
	cond.Field(metadata.AssociationFieldAssociationId).Eq(id)
	err := s.core.AssociationOperation().UpdateAssociation(params, data, cond)
	return nil, err

}
