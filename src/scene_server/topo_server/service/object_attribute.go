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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateObjectAttribute create a new object attribute
func (s *Service) CreateObjectAttribute(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	attr, err := s.Core.AttributeOperation().CreateObjectAttribute(params, data)
	if nil != err {
		return nil, err
	}

	// auth: register resource
	attribute := attr.Attribute()
	if err := s.AuthManager.RegisterModelAttribute(params.Context, params.Header, *attribute); err != nil {
		blog.Errorf("create object attribute success, but register model attribute to auth failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}

	cond := condition.CreateCondition()
	cond.Field("id").Eq(attribute.ID)
	attrInfo, err := s.Core.AttributeOperation().FindObjectAttributeWithDetail(params, cond)
	if err != nil {
		blog.Errorf("create object attribute success, but get attributes detail failed, err: %v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrorTopoSearchModelAttriFailedPleaseRefresh)
	}
	if len(attrInfo) <= 0 {
		blog.Errorf("create object attribute success, but get attributes detail failed, err: %v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrorTopoSearchModelAttriFailedPleaseRefresh)
	}

	return attrInfo[0], nil
}

// SearchObjectAttribute search the object attributes
func (s *Service) SearchObjectAttribute(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	cond := condition.CreateCondition()
	if data.Exists(metadata.PageName) {
		page, err := data.MapStr(metadata.PageName)
		if err != nil {
			blog.Errorf("SearchObjectAttribute failed, page info convert to mapstr failed, page: %v, err: %v, rid: %s", data[metadata.PageName], err, params.ReqID)
			return nil, err
		}
		if err := cond.SetPage(page); err != nil {
			blog.Errorf("SearchObjectAttribute, cond set page failed, page: %v, err: %v, rid: %v", page, err, params.ReqID)
			return nil, err
		}
		data.Remove(metadata.PageName)
	}

	if err := cond.Parse(data); nil != err {
		blog.Errorf("search object attribute, but failed to parse the data into condition, err: %v, rid: %s", err, params.ReqID)
		return nil, err
	}

	cond.Field(metadata.AttributeFieldIsSystem).NotEq(true)
	cond.Field(metadata.AttributeFieldIsAPI).NotEq(true)

	return s.Core.AttributeOperation().FindObjectAttributeWithDetail(params, cond)
}

// UpdateObjectAttribute update the object attribute
func (s *Service) UpdateObjectAttribute(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	paramPath := mapstr.MapStr{}
	paramPath.Set("id", pathParams("id"))
	id, err := paramPath.Int64("id")
	if nil != err {
		blog.Errorf("[api-att] failed to parse the path params id(%s), error info is %s, rid: %s", pathParams("id"), err.Error(), params.ReqID)
		return nil, err
	}
	// TODO: why does remove this????
	data.Remove(metadata.BKMetadata)

	// UpdateObjectAttribute should not update bk_property_index、bk_property_group
	data.Remove(common.BKPropertyIndexField)
	data.Remove(common.BKPropertyGroupField)

	err = s.Core.AttributeOperation().UpdateObjectAttribute(params, data, id)

	// auth: update registered resource
	if err := s.AuthManager.UpdateRegisteredModelAttributeByID(params.Context, params.Header, id); err != nil {
		blog.Errorf("update object attribute success , but update registered model attribute to auth failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}

	return nil, err
}

// DeleteObjectAttribute delete the object attribute
func (s *Service) DeleteObjectAttribute(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	paramPath := mapstr.MapStr{}
	paramPath.Set("id", pathParams("id"))
	id, err := paramPath.Int64("id")
	if nil != err {
		blog.Errorf("[api-att] failed to parse the path params id(%s), error info is %s , rid: %s", pathParams("id"), err.Error(), params.ReqID)
		return nil, err
	}

	cond := condition.CreateCondition()
	cond.Field(metadata.AttributeFieldID).Eq(id)

	data.Remove(metadata.BKMetadata)

	listRuleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: []int64{id},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	ruleResult, err := s.Engine.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(params.Context, params.Header, 0, listRuleOption)
	if err != nil {
		blog.Errorf("delete object attribute failed, ListHostApplyRule failed, listRuleOption: %+v, err: %+v, rid: %s", listRuleOption, err, params.ReqID)
		return nil, err
	}
	ruleIDs := make([]int64, 0)
	for _, item := range ruleResult.Info {
		ruleIDs = append(ruleIDs, item.ID)
	}

	// auth: update registered resource
	if err := s.AuthManager.DeregisterModelAttributeByID(params.Context, params.Header, id); err != nil {
		blog.Errorf("delete object attribute failed, deregister model attribute to auth failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommUnRegistResourceToIAMFailed)
	}

	err = s.Core.AttributeOperation().DeleteObjectAttribute(params, cond)
	if err != nil {
		blog.Errorf("delete object attribute failed, DeleteObjectAttribute failed, params: %+v, err: %+v, rid: %s", params, err, params.ReqID)
		return nil, err
	}

	if len(ruleIDs) > 0 {
		deleteRuleOption := metadata.DeleteHostApplyRuleOption{
			RuleIDs: ruleIDs,
		}
		if err := s.Engine.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(params.Context, params.Header, 0, deleteRuleOption); err != nil {
			blog.Errorf("delete object attribute success, but DeleteHostApplyRule failed, params: %+v, err: %+v, rid: %s", deleteRuleOption, err, params.ReqID)
			return nil, err
		}
	}

	return nil, err
}

func (s *Service) UpdateObjectAttributeIndex(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	paramPath := mapstr.MapStr{}
	paramPath.Set("id", pathParams("id"))
	paramPath.Set(common.BKObjIDField, pathParams(common.BKObjIDField))
	objID, err := paramPath.String(common.BKObjIDField)
	if nil != err {
		blog.Errorf("[api-att] failed to parse the path params bk_obj_id(%s), error info is %s , rid: %s", objID, err.Error(), params.ReqID)
		return nil, params.Err.Error(common.CCErrorTopoPathParamPaserFailed)
	}

	id, err := paramPath.Int64("id")
	if nil != err {
		blog.Errorf("[api-att] failed to parse the params id(%s), error info is %s , rid: %s", data["id"], err.Error(), params.ReqID)
		return nil, params.Err.Error(common.CCErrorTopoPathParamPaserFailed)
	}
	result, err := s.Core.AttributeOperation().UpdateObjectAttributeIndex(params, objID, data, id)
	if err != nil {
		blog.Errorf("UpdateObjectAttributeIndex failed, error info is %s , rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	return result, nil
}

// ListHostModelAttribute list host model's attributes
func (s *Service) ListHostModelAttribute(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	cond := condition.CreateCondition()
	data.Remove(metadata.PageName)
	if err := cond.Parse(data); nil != err {
		blog.Errorf("search object attribute, but failed to parse the data into condition, err: %v, rid: %s", err, params.ReqID)
		return nil, err
	}
	cond.Field(metadata.AttributeFieldIsSystem).NotEq(true)
	cond.Field(metadata.AttributeFieldIsAPI).NotEq(true)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDHost)
	attributes, err := s.Core.AttributeOperation().FindObjectAttributeWithDetail(params, cond)
	if err != nil {
		return nil, err
	}
	hostAttributes := make([]metadata.HostObjAttDes, 0)
	for _, item := range attributes {
		if item == nil {
			continue
		}
		hostApplyEnabled := metadata.CheckAllowHostApplyOnField(item.PropertyID)
		hostAttribute := metadata.HostObjAttDes{
			ObjAttDes:        *item,
			HostApplyEnabled: hostApplyEnabled,
		}
		hostAttributes = append(hostAttributes, hostAttribute)
	}
	return hostAttributes, nil
}
