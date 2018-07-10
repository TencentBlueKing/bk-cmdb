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
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateObjectGroup create a new object group
func (s *topoService) CreateObjectGroup(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	rsp, err := s.core.GroupOperation().CreateObjectGroup(params, data)
	if nil != err {
		return nil, err
	}

	return rsp.ToMapStr()
}

// UpdateObjectGroup update the object group information
func (s *topoService) UpdateObjectGroup(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	fmt.Println("UpdateObjectGroup")

	cond := &metadata.UpdateGroupCondition{}

	err := data.MarshalJSONInto(cond)
	if nil != err {
		return nil, err
	}

	err = s.core.GroupOperation().UpdateObjectGroup(params, cond)
	if nil != err {
		return nil, err
	}

	return nil, nil
}

// DeleteObjectGroup delete the object group
func (s *topoService) DeleteObjectGroup(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	err := s.core.GroupOperation().DeleteObjectGroup(params, pathParams("id"))
	if nil != err {
		return nil, err
	}

	return nil, nil
}

// UpdateObjectAttributeGroup update the object attribute belongs to group information
func (s *topoService) UpdateObjectAttributeGroup(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("UpdateObjectAttributeGroup")
	cond := condition.CreateCondition()

	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	val, exists := data.Get(metadata.ModelFieldObjectID)
	if !exists {
		return nil, params.Err.Errorf(common.CCErrCommParamsLostField, metadata.ModelFieldObjectID)
	}
	cond.Field(metadata.ModelFieldObjectID).Eq(val)
	data.Remove(metadata.ModelFieldObjectID)

	val, exists = data.Get(metadata.AttributeFieldPropertyID)
	if !exists {
		return nil, params.Err.Errorf(common.CCErrCommParamsLostField, metadata.AttributeFieldPropertyID)
	}
	cond.Field(metadata.AttributeFieldPropertyID).Eq(val)
	data.Remove(metadata.AttributeFieldPropertyID)

	err := s.core.AttributeOperation().UpdateObjectAttribute(params, data, -1, cond)
	if nil != err {
		return nil, err
	}

	return nil, nil
}

// DeleteObjectAttributeGroup delete the object attribute belongs to group information
func (s *topoService) DeleteObjectAttributeGroup(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("DeleteObjectAttributeGroup")
	cond := condition.CreateCondition()

	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	cond.Field(metadata.ModelFieldObjectID).Eq(pathParams("object_id"))
	cond.Field(metadata.AttributeFieldPropertyID).Eq(pathParams("property_id"))
	cond.Field(metadata.AttributeFieldPropertyGroup).Eq("group_id")

	innerData := frtypes.MapStr{}
	innerData.Set(metadata.AttributeFieldPropertyGroup, "none")
	innerData.Set(metadata.AttributeFieldPropertyIndex, -1)

	err := s.core.AttributeOperation().UpdateObjectAttribute(params, innerData, -1, cond)
	if nil != err {
		return nil, err
	}

	return nil, nil
}

// SearchGroupByObject search the groups by the object
func (s *topoService) SearchGroupByObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	return s.core.GroupOperation().FindGroupByObject(params, pathParams("object_id"), cond)
}
