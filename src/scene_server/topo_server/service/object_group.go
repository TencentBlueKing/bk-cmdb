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
	"encoding/json"
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateObjectGroup create a new object group

func (s *Service) CreateObjectGroup(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	rsp, err := s.Core.GroupOperation().CreateObjectGroup(params, data)
	if nil != err {
		return nil, err
	}

	// auth: register attribute group
	if err := s.AuthManager.RegisterModelAttributeGroup(params.Context, params.Header, rsp.Group()); err != nil {
		return nil, fmt.Errorf("register attribute group to iam failed, err: %+v", err)
	}
	return rsp.ToMapStr()
}

// UpdateObjectGroup update the object group information
func (s *Service) UpdateObjectGroup(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	cond := &metadata.UpdateGroupCondition{}

	err := data.MarshalJSONInto(cond)
	if nil != err {
		return nil, err
	}
	data.Remove(metadata.BKMetadata)

	err = s.Core.GroupOperation().UpdateObjectGroup(params, cond)
	if nil != err {
		return nil, err
	}

	// query attribute groups with given condition, so that update them to iam after updated
	searchCondition := condition.CreateCondition()
	if cond.Condition.ID != 0 {
		searchCondition.Field(common.BKFieldID).Eq(cond.Condition.ID)
	}
	if cond.Condition.GroupID != "" {
		searchCondition.Field(common.BKPropertyGroupIDField).Eq(cond.Condition.GroupID)
	}
	if cond.Condition.ObjID != "" {
		searchCondition.Field(common.BKObjIDField).Eq(cond.Condition.ObjID)
	}
	queryCond := metadata.QueryCondition{
		Condition: searchCondition.ToMapStr(),
	}
	result, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(params.Context, params.Header, common.BKTableNamePropertyGroup, &queryCond)
	if err != nil {
		return nil, fmt.Errorf("search attribute group by condition failed, err: %+v", err)
	}
	attributeGroups := make([]metadata.Group, 0)
	for _, item := range result.Data.Info {
		ag := metadata.Group{}
		ag.Parse(item)
		attributeGroups = append(attributeGroups, ag)
	}

	// auth: register attribute group
	if err := s.AuthManager.UpdateRegisteredModelAttributeGroup(params.Context, params.Header, attributeGroups...); err != nil {
		return nil, fmt.Errorf("update attribute group to iam failed, err: %+v", err)
	}
	return nil, nil
}

// DeleteObjectGroup delete the object group
func (s *Service) DeleteObjectGroup(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	gid, err := strconv.ParseInt(pathParams("id"), 10, 64)
	if nil != err {
		return nil, err
	}

	data.Remove(metadata.BKMetadata)

	err = s.Core.GroupOperation().DeleteObjectGroup(params, gid)
	if nil != err {
		return nil, err
	}

	// auth: deregister attribute group
	if err := s.AuthManager.DeregisterModelAttributeGroupByID(params.Context, params.Header, gid); err != nil {
		return nil, fmt.Errorf("deregister attribute group to iam failed, err: %+v", err)
	}
	return nil, nil
}

func (s *Service) ParseUpdateObjectAttributeGroupPropertyInput(data []byte) (mapstr.MapStr, error) {
	requestBody := struct{
		data []metadata.PropertyGroupObjectAtt `json:"data",field:"json"`
		metadata metadata.Metadata `json:"metadata",field:"metadata"`
	}{}
	err := json.Unmarshal(data, &requestBody)
	if nil != err {
		return nil, err
	}
	result := mapstr.MapStr{}
	result.Set("origin", requestBody.data)
	return result, nil
}

// UpdateObjectAttributeGroupProperty update the object attribute belongs to group information
func (s *Service) UpdateObjectAttributeGroupProperty(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	datas := make([]metadata.PropertyGroupObjectAtt, 0)
	val, exists := data.Get("origin")
	if !exists {
		return nil, params.Err.New(common.CCErrCommParamsIsInvalid, "not set anything")
	}

	datas, _ = val.([]metadata.PropertyGroupObjectAtt)

	err := s.Core.GroupOperation().UpdateObjectAttributeGroup(params, datas)
	if nil != err {
		return nil, err
	}

	return nil, nil
}

// DeleteObjectAttributeGroup delete the object attribute belongs to group information

func (s *Service) DeleteObjectAttributeGroup(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	data.Remove(metadata.BKMetadata)
	err := s.Core.GroupOperation().DeleteObjectAttributeGroup(params, pathParams("bk_object_id"), pathParams("property_id"), pathParams("group_id"))
	if nil != err {
		return nil, err
	}

	// query attribute groups with given condition, so that update them to iam after updated
	searchCondition := condition.CreateCondition()
	searchCondition.Field(common.BKObjIDField).Eq(pathParams("bk_object_id"))
	searchCondition.Field(common.BKPropertyIDField).Eq(pathParams("property_id"))
	searchCondition.Field(common.BKPropertyGroupIDField).Eq(pathParams("group_id"))
	queryCondition := metadata.QueryCondition{
		Condition: searchCondition.ToMapStr(),
	}
	result, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(params.Context, params.Header, common.BKTableNamePropertyGroup, &queryCondition)
	if err != nil {
		return nil, fmt.Errorf("search attribute group by condition failed, err: %+v", err)
	}
	attributeGroups := make([]metadata.Group, 0)
	for _, item := range result.Data.Info {
		ag := metadata.Group{}
		ag.Parse(item)
		attributeGroups = append(attributeGroups, ag)
	}

	// auth: deregister attribute group
	if err := s.AuthManager.DeregisterModelAttributeGroup(params.Context, params.Header, attributeGroups...); err != nil {
		return nil, fmt.Errorf("update attribute group to iam failed, err: %+v", err)
	}
	return nil, nil
}

// SearchGroupByObject search the groups by the object
func (s *Service) SearchGroupByObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()

	return s.Core.GroupOperation().FindGroupByObject(params, pathParams("bk_obj_id"), cond)
}
