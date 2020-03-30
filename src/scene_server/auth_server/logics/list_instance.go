/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logics

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/types"

	"gopkg.in/redis.v5"
)

// TODO confirm 422 need to be used in which case

// list instances by condition
func (lgc *Logics) ListInstance(kit *rest.Kit, cond map[string]interface{}, resourceType iam.ResourceTypeID, page types.Page) (*types.ListInstanceResult, error) {
	idField := types.GetResourceIDField(resourceType)
	nameField := types.GetResourceNameField(resourceType)
	if idField == "" || nameField == "" {
		blog.Errorf("request type %s is invalid, rid: %s", resourceType, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "type")
	}

	param := metadata.PullResourceParam{
		Condition: cond,
		Fields:    []string{idField, nameField},
		Limit:     page.Limit,
		Offset:    page.Offset,
	}
	data, err := lgc.searchAuthResource(kit, param, resourceType)
	if err != nil {
		blog.ErrorJSON("search auth resource failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
		return nil, err
	}

	instances := make([]types.InstanceResource, 0)
	for _, instance := range data.Info {
		instances = append(instances, types.InstanceResource{
			ID:          util.GetStrByInterface(instance[idField]),
			DisplayName: util.GetStrByInterface(instance[nameField]),
		})
	}
	return &types.ListInstanceResult{
		Count:   data.Count,
		Results: instances,
	}, nil
}

// search auth resource instances from database
func (lgc *Logics) searchAuthResource(kit *rest.Kit, param metadata.PullResourceParam, resourceType iam.ResourceTypeID) (*metadata.PullResourceResult, error) {
	param.Collection = getResourceTableName(resourceType)
	if param.Collection == "" {
		blog.Errorf("request type %s is invalid, rid: %s", resourceType, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "type")
	}

	var err error
	param.Condition, err = lgc.generateSpecialCondition(kit, resourceType, param.Condition)
	if err != nil {
		blog.ErrorJSON("generate special condition failed, error: %s, cond: %s, rid: %s", err.Error(), param.Condition, kit.Rid)
		return nil, err
	}

	res, err := lgc.CoreAPI.CoreService().Auth().SearchAuthResource(kit.Ctx, kit.Header, param)
	if err != nil {
		blog.ErrorJSON("search auth resource failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
		return nil, err
	}
	if !res.Result {
		blog.ErrorJSON("search auth resource failed, error code: %s, error message: %s, param: %s, rid: %s", res.Code, res.ErrMsg, param, kit.Rid)
		return nil, res.Error()
	}
	return &res.Data, nil
}

// list system scope instances that have no parent
func (lgc *Logics) ListSystemInstance(kit *rest.Kit, req types.PullResourceReq) (*types.ListInstanceResult, error) {
	if req.Page.IsIllegal() {
		blog.Errorf("request page limit %d exceeds max page size, rid: %s", req.Page.Limit, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommPageLimitIsExceeded)
	}
	cond := make(map[string]interface{})
	if req.Filter == nil {
		return lgc.ListInstance(kit, cond, req.Type, req.Page)
	}
	filter, ok := req.Filter.(types.ListInstanceFilter)
	if !ok {
		blog.ErrorJSON("request filter %s is not the right type for list_instance method, rid: %s", req.Filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "filter")
	}
	// system resource has no parent
	if filter.Parent != nil || len(filter.ResourceTypeChain) != 0 {
		return &types.ListInstanceResult{
			Count:   0,
			Results: make([]types.InstanceResource, 0),
		}, nil
	}
	for resourceType, keywords := range filter.Search {
		if resourceType != req.Type {
			return &types.ListInstanceResult{
				Count:   0,
				Results: make([]types.InstanceResource, 0),
			}, nil
		}
		keywordCond := map[string]interface{}{
			common.BKDBLIKE: strings.Join(keywords, "|"),
		}
		cond[common.BKDBOR] = []map[string]interface{}{
			{types.GetResourceIDField(req.Type): keywordCond},
			{types.GetResourceNameField(req.Type): keywordCond},
		}
	}
	return lgc.ListInstance(kit, cond, req.Type, req.Page)
}

// list business scope instances TODO confirm if these resource only has biz as parent
func (lgc *Logics) ListBusinessInstance(kit *rest.Kit, req types.PullResourceReq) (*types.ListInstanceResult, error) {
	if req.Page.IsIllegal() {
		blog.Errorf("request page limit %d exceeds max page size, rid: %s", req.Page.Limit, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommPageLimitIsExceeded)
	}
	cond := make(map[string]interface{})
	if req.Filter == nil {
		return lgc.ListInstance(kit, cond, req.Type, req.Page)
	}
	filter, ok := req.Filter.(types.ListInstanceFilter)
	if !ok {
		blog.ErrorJSON("request filter %s is not the right type for list_instance method, rid: %s", req.Filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "filter")
	}

	if filter.Parent != nil {
		parentIDField := types.GetResourceIDField(filter.Parent.Type)
		if parentIDField == "" {
			blog.Errorf("filter parent type %s is invalid, rid: %s", filter.Parent.Type, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "filter.parent.type")
		}
		if isResourceIDStringType(req.Type) {
			cond[parentIDField] = filter.Parent.ID
		} else {
			// if parent id type is not string, filter parent id is not int64, return empty result
			id, err := strconv.ParseInt(filter.Parent.ID, 10, 64)
			if err != nil {
				blog.Errorf("filter.parent.id %s parse int failed, error: %s, rid: %s", filter.Parent.ID, err.Error(), kit.Rid)
				return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
			}
			cond[parentIDField] = id
		}
	}

	// stores ancestor mapping, ancestorMap[ancestor type][ancestor id] = ancestor display_name
	ancestorMap := make(map[iam.ResourceTypeID]map[string]string)
	idField := types.GetResourceIDField(req.Type)
	nameField := types.GetResourceNameField(req.Type)
	// all ancestor conditions are 'and' relation
	ancestorCond := make([]map[string]interface{}, 0)
	for resourceType, keywords := range filter.Search {
		keywordCond := map[string]interface{}{
			common.BKDBLIKE: strings.Join(keywords, "|"),
		}
		if resourceType == req.Type {
			cond[common.BKDBOR] = []map[string]interface{}{
				{idField: keywordCond},
				{nameField: keywordCond},
			}
			continue
		}

		// filter ancestor topology, try to get from cache first
		objID := GetInstanceResourceObjID(resourceType)
		ancestorIDField := types.GetResourceIDField(resourceType)
		ancestorNameField := types.GetResourceNameField(resourceType)
		if objID == "" {
			// ancestor type not exist
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
		// get ancestor from cache
		keys, err := lgc.cache.Keys(common.BKCacheKeyV3Prefix + objID + "*").Result()
		if err != nil {
			blog.ErrorJSON("get cache keys for object %s failed, error: %s, rid: %s", objID, err.Error(), kit.Rid)
			return nil, err
		}
		pipeline := lgc.cache.Pipeline()
		for _, key := range keys {
			pipeline.HGetAll(key)
		}
		results, err := pipeline.Exec()
		if err != nil {
			blog.ErrorJSON("get cached instances for object %s using keys: %v failed, error: %s, rid: %s", objID, keys, err.Error(), kit.Rid)
			return nil, err
		}
		// filter instance using ancestorIDs
		ancestorIDs := make([]interface{}, 0)
		for _, result := range results {
			cmd := result.(*redis.StringStringMapCmd)
			instance, err := cmd.Result()
			if err != nil {
				blog.ErrorJSON("get cached instance result failed, error: %s, keys: %v, rid: %s", objID, keys, err.Error(), kit.Rid)
				return nil, err
			}
			regex := regexp.MustCompile(fmt.Sprintf(".*%s.*", strings.Join(keywords, "|")))
			ancestorName := util.GetStrByInterface(instance[ancestorNameField])
			ancestorID := util.GetStrByInterface(instance[ancestorIDField])
			if regex.MatchString(ancestorName) || regex.MatchString(ancestorID) {
				if ancestorMap[resourceType] == nil {
					ancestorMap[resourceType] = map[string]string{
						ancestorID: ancestorName,
					}
				} else {
					ancestorMap[resourceType][ancestorID] = ancestorName
				}
				if isResourceIDStringType(resourceType) {
					ancestorIDs = append(ancestorIDs, instance[ancestorIDField])
				} else {
					ancestorID, err := strconv.ParseInt(instance[ancestorIDField], 10, 64)
					if err != nil {
						blog.ErrorJSON("parse ancestorID %s to int64 failed, error: %s, rid: %s", ancestorID, err.Error(), kit.Rid)
						return nil, err
					}
					ancestorIDs = append(ancestorIDs, ancestorID)
				}
			}
		}
		if len(ancestorIDs) > 0 {
			// business instances stores relation in field like bk_biz_id which is the ancestor's id field
			ancestorCond = append(ancestorCond, map[string]interface{}{
				ancestorIDField: map[string]interface{}{common.BKDBOR: ancestorIDs},
			})
		}
	}
	if len(ancestorCond) > 0 {
		cond[common.BKDBAND] = ancestorCond
	}

	if len(filter.ResourceTypeChain) == 0 {
		return lgc.ListInstance(kit, cond, req.Type, req.Page)
	}

	// get instances with path info
	param := metadata.PullResourceParam{
		Condition: cond,
		Limit:     req.Page.Limit,
		Offset:    req.Page.Offset,
	}
	data, err := lgc.searchAuthResource(kit, param, req.Type)
	if err != nil {
		blog.ErrorJSON("search auth resource failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
		return nil, err
	}
	instances := make([]types.InstanceResource, 0)
	for _, inst := range data.Info {
		instance := types.InstanceResource{
			ID:          util.GetStrByInterface(inst[idField]),
			DisplayName: util.GetStrByInterface(inst[nameField]),
		}
		for _, ancestor := range filter.ResourceTypeChain {
			ancestorID := util.GetStrByInterface(inst[types.GetResourceIDField(ancestor.ID)])
			instance.Path = append(instance.Path, types.InstancePath{
				Type:        ancestor.ID,
				ID:          ancestorID,
				DisplayName: ancestorMap[ancestor.ID][ancestorID],
			})
		}
		instances = append(instances, instance)
	}
	return &types.ListInstanceResult{
		Count:   data.Count,
		Results: instances,
	}, nil
}

// list model instances, parent is model
func (lgc *Logics) ListModelInstance(kit *rest.Kit, req types.PullResourceReq) (*types.ListInstanceResult, error) {
	if req.Page.IsIllegal() {
		blog.Errorf("request page limit %d exceeds max page size, rid: %s", req.Page.Limit, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommPageLimitIsExceeded)
	}
	cond := make(map[string]interface{})
	if req.Filter == nil {
		return lgc.ListInstance(kit, cond, req.Type, req.Page)
	}
	filter, ok := req.Filter.(types.ListInstanceFilter)
	if !ok {
		blog.ErrorJSON("request filter %s is not the right type for list_instance method, rid: %s", req.Filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "filter")
	}

	var parentObjID string
	if filter.Parent != nil {
		if filter.Parent.Type != iam.SysModel {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
		parentObjID = filter.Parent.ID
		cond[common.BKObjIDField] = filter.Parent.ID
	}

	// stores model mapping, modelMap[model id] = model display_name
	modelMap := make(map[string]string)
	idField := types.GetResourceIDField(req.Type)
	nameField := types.GetResourceNameField(req.Type)
	for resourceType, keywords := range filter.Search {
		keywordCond := map[string]interface{}{
			common.BKDBLIKE: strings.Join(keywords, "|"),
		}
		if resourceType == req.Type {
			cond[common.BKDBOR] = []map[string]interface{}{
				{idField: keywordCond},
				{nameField: keywordCond},
			}
			continue
		}

		if resourceType != iam.SysModel {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
		regex := regexp.MustCompile(fmt.Sprintf(".*%s.*", strings.Join(keywords, "|")))
		// check if parent object matches the search filter
		if parentObjID != "" {
			if regex.MatchString(parentObjID) {
				continue
			}
			modelRes, err := lgc.CoreAPI.CoreService().Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{
				Fields:    []string{common.BKObjNameField},
				Page:      metadata.BasePage{Limit: 1},
				Condition: map[string]interface{}{common.BKObjIDField: parentObjID},
			})
			if err != nil {
				blog.Errorf("get parent model %s failed, error: %s, rid: %s", parentObjID, err.Error(), kit.Rid)
				return nil, err
			}
			if len(modelRes.Data.Info) == 0 {
				return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
			}
			if regex.MatchString(modelRes.Data.Info[0].Spec.ObjectName) {
				modelMap[parentObjID] = modelRes.Data.Info[0].Spec.ObjectName
				continue
			}
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}

		modelRes, err := lgc.CoreAPI.CoreService().Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{
			Fields: []string{common.BKObjIDField, common.BKObjNameField},
			Page:   metadata.BasePage{Limit: 1},
			Condition: map[string]interface{}{common.BKDBOR: []map[string]interface{}{
				{common.BKObjIDField: keywordCond},
				{common.BKObjNameField: keywordCond},
			}},
		})
		if err != nil {
			blog.Errorf("get model failed, error: %s, rid: %s", err.Error(), kit.Rid)
			return nil, err
		}
		objIDs := make([]string, 0)
		for _, model := range modelRes.Data.Info {
			modelMap[model.Spec.ObjectID] = model.Spec.ObjectName
			objIDs = append(objIDs, model.Spec.ObjectID)
		}
		if len(objIDs) == 0 {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
		cond[common.BKObjIDField] = map[string]interface{}{
			common.BKDBIN: objIDs,
		}
	}

	if len(filter.ResourceTypeChain) == 0 {
		return lgc.ListInstance(kit, cond, req.Type, req.Page)
	}

	// get all model if no model is filtered
	if len(modelMap) == 0 {
		modelRes, err := lgc.CoreAPI.CoreService().Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{
			Fields:    []string{common.BKObjIDField, common.BKObjNameField},
			Page:      metadata.BasePage{Limit: 1},
			Condition: map[string]interface{}{},
		})
		if err != nil {
			blog.Errorf("get model failed, error: %s, rid: %s", err.Error(), kit.Rid)
			return nil, err
		}
		for _, model := range modelRes.Data.Info {
			modelMap[model.Spec.ObjectID] = model.Spec.ObjectName
		}
	}

	// get instances with path info
	param := metadata.PullResourceParam{
		Condition: cond,
		Limit:     req.Page.Limit,
		Offset:    req.Page.Offset,
	}
	data, err := lgc.searchAuthResource(kit, param, req.Type)
	if err != nil {
		blog.ErrorJSON("search auth resource failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
		return nil, err
	}
	instances := make([]types.InstanceResource, 0)
	for _, inst := range data.Info {
		instance := types.InstanceResource{
			ID:          util.GetStrByInterface(inst[idField]),
			DisplayName: util.GetStrByInterface(inst[nameField]),
		}
		for _, ancestor := range filter.ResourceTypeChain {
			if ancestor.ID == iam.SysModel {
				ancestorID := util.GetStrByInterface(inst[types.GetResourceIDField(ancestor.ID)])
				instance.Path = append(instance.Path, types.InstancePath{
					Type:        ancestor.ID,
					ID:          ancestorID,
					DisplayName: modelMap[ancestorID],
				})
			}
		}
		instances = append(instances, instance)
	}
	return &types.ListInstanceResult{
		Count:   data.Count,
		Results: instances,
	}, nil
}

// list host instances
func (lgc *Logics) ListHostInstance(kit *rest.Kit, req types.PullResourceReq) (*types.ListInstanceResult, error) {
	if req.Page.IsIllegal() {
		blog.Errorf("request page limit %d exceeds max page size, rid: %s", req.Page.Limit, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommPageLimitIsExceeded)
	}
	cond := make(map[string]interface{})
	if req.Filter == nil {
		return lgc.ListInstance(kit, cond, req.Type, req.Page)
	}
	filter, ok := req.Filter.(types.ListInstanceFilter)
	if !ok {
		blog.ErrorJSON("request filter %s is not the right type for list_instance method, rid: %s", req.Filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "filter")
	}

	hostIDs := make([]int64, 0)
	// stores host relation in the for hostRelations[objID][instance ID] = host IDs
	hostRelations := map[string]map[string][]int64{
		common.BKInnerObjIDApp:    make(map[string][]int64),
		common.BKInnerObjIDSet:    make(map[string][]int64),
		common.BKInnerObjIDModule: make(map[string][]int64),
	}
	// maps hostID to host relation
	hostIDRelationMap := make(map[int64]map[string]string)
	if filter.Parent != nil || filter.Search != nil {
		// get host topology relation from cache
		keys, err := lgc.cache.Keys(common.BKCacheKeyV3Prefix + "host_module:*").Result()
		if err != nil {
			blog.ErrorJSON("get cache keys for host module relation failed, error: %s, rid: %s", err.Error(), kit.Rid)
			return nil, err
		}
		pipeline := lgc.cache.Pipeline()
		for _, key := range keys {
			pipeline.HGetAll(key)
		}
		results, err := pipeline.Exec()
		if err != nil {
			blog.ErrorJSON("get cached host module relations using keys: %v failed, error: %s, rid: %s", keys, err.Error(), kit.Rid)
			return nil, err
		}
		for _, result := range results {
			cmd := result.(*redis.StringStringMapCmd)
			relation, err := cmd.Result()
			if err != nil {
				blog.ErrorJSON("get cached host module relation result failed, error: %s, rid: %s", err.Error(), kit.Rid)
				return nil, err
			}
			hostID, err := util.GetInt64ByInterface(relation[common.BKHostIDField])
			if err != nil {
				blog.ErrorJSON("get host module relation hostID failed, error: %s, rid: %s", err.Error(), kit.Rid)
				return nil, err
			}
			hostIDRelationMap[hostID] = relation
			hostRelations[common.BKInnerObjIDApp][relation[common.BKAppIDField]] = append(hostRelations[common.BKInnerObjIDApp][relation[common.BKAppIDField]], hostID)
			hostRelations[common.BKInnerObjIDSet][relation[common.BKSetIDField]] = append(hostRelations[common.BKInnerObjIDSet][relation[common.BKSetIDField]], hostID)
			hostRelations[common.BKInnerObjIDModule][relation[common.BKModuleIDField]] = append(hostRelations[common.BKInnerObjIDModule][relation[common.BKModuleIDField]], hostID)
		}
	}
	if filter.Parent != nil {
		objID := GetInstanceResourceObjID(filter.Parent.Type)
		// compatible for case like module that has no specified auth resource TODO confirm this
		if objID == "" {
			objID = string(filter.Parent.Type)
		}
		hostIDs = hostRelations[objID][filter.Parent.ID]
		if len(hostIDs) == 0 {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
	}

	// stores ancestor mapping, ancestorMap[ancestor type][ancestor id] = ancestor display_name
	ancestorMap := make(map[iam.ResourceTypeID]map[string]string)
	idField := types.GetResourceIDField(req.Type)
	nameField := types.GetResourceNameField(req.Type)
	for resourceType, keywords := range filter.Search {
		keywordCond := map[string]interface{}{
			common.BKDBLIKE: strings.Join(keywords, "|"),
		}
		if resourceType == req.Type {
			cond[common.BKDBOR] = []map[string]interface{}{
				{idField: keywordCond},
				{nameField: keywordCond},
			}
			continue
		}

		// filter ancestor topology, try to get from cache first
		objID := GetInstanceResourceObjID(resourceType)
		ancestorIDField := types.GetResourceIDField(resourceType)
		ancestorNameField := types.GetResourceNameField(resourceType)
		// compatible for case like module that has no specified auth resource TODO confirm this
		if objID == "" {
			objID = string(resourceType)
			ancestorIDField = common.GetInstIDField(objID)
			ancestorNameField = common.GetInstNameField(objID)
		}
		// get ancestor from cache
		keys, err := lgc.cache.Keys(common.BKCacheKeyV3Prefix + objID + "*").Result()
		if err != nil {
			blog.ErrorJSON("get cache keys for object %s failed, error: %s, rid: %s", objID, err.Error(), kit.Rid)
			return nil, err
		}
		pipeline := lgc.cache.Pipeline()
		for _, key := range keys {
			pipeline.HGetAll(key)
		}
		results, err := pipeline.Exec()
		if err != nil {
			blog.ErrorJSON("get cached instances for object %s using keys: %v failed, error: %s, rid: %s", objID, keys, err.Error(), kit.Rid)
			return nil, err
		}
		// filter instance using ancestorIDs
		ancestorIDs := make([]string, 0)
		for _, result := range results {
			cmd := result.(*redis.StringStringMapCmd)
			instance, err := cmd.Result()
			if err != nil {
				blog.ErrorJSON("get cached instance result failed, error: %s, keys: %v, rid: %s", objID, keys, err.Error(), kit.Rid)
				return nil, err
			}
			regex := regexp.MustCompile(fmt.Sprintf(".*%s.*", strings.Join(keywords, "|")))
			ancestorName := util.GetStrByInterface(instance[ancestorNameField])
			ancestorID := util.GetStrByInterface(instance[ancestorIDField])
			if regex.MatchString(ancestorName) || regex.MatchString(ancestorID) {
				if ancestorMap[resourceType] == nil {
					ancestorMap[resourceType] = map[string]string{
						ancestorID: ancestorName,
					}
				} else {
					ancestorMap[resourceType][ancestorID] = ancestorName
				}
				ancestorIDs = append(ancestorIDs, instance[ancestorIDField])
			}
		}
		filterHostIDs := make([]int64, 0)
		for _, ancestorID := range ancestorIDs {
			filterHostIDs = append(filterHostIDs, hostRelations[objID][ancestorID]...)
		}
		hostIDs = util.IntArrIntersection(hostIDs, filterHostIDs)
		if len(hostIDs) == 0 {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
	}
	cond[common.BKHostIDField] = map[string]interface{}{
		common.BKDBIN: hostIDs,
	}

	if len(filter.ResourceTypeChain) == 0 {
		return lgc.ListInstance(kit, cond, req.Type, req.Page)
	}

	// get ancestor needed for ResourceTypeChain from cache
	relatedIDs := make(map[string][]string)
	for _, hostID := range hostIDs {
		for idField, id := range hostIDRelationMap[hostID] {
			relatedIDs[idField] = append(relatedIDs[idField], id)
		}
	}
	for _, ancestor := range filter.ResourceTypeChain {
		resourceType := ancestor.ID
		if ancestorMap[resourceType] != nil {
			continue
		}
		objID := GetInstanceResourceObjID(resourceType)
		ancestorIDField := types.GetResourceIDField(resourceType)
		ancestorNameField := types.GetResourceNameField(resourceType)
		if objID == "" {
			objID = string(resourceType)
			ancestorIDField = common.GetInstIDField(objID)
			ancestorNameField = common.GetInstNameField(objID)
		}
		pipeline := lgc.cache.Pipeline()
		for _, id := range relatedIDs[idField] {
			pipeline.HGetAll(common.BKCacheKeyV3Prefix + objID + "*id:" + id)
		}
		results, err := pipeline.Exec()
		if err != nil {
			blog.ErrorJSON("get cached instances for object %s with ids: %v failed, error: %s, rid: %s", objID, relatedIDs[idField], err.Error(), kit.Rid)
			return nil, err
		}
		for _, result := range results {
			cmd := result.(*redis.StringStringMapCmd)
			instance, err := cmd.Result()
			if err != nil {
				blog.ErrorJSON("get cached instance result failed, error: %s, rid: %s", objID, err.Error(), kit.Rid)
				return nil, err
			}
			ancestorName := util.GetStrByInterface(instance[ancestorNameField])
			ancestorID := util.GetStrByInterface(instance[ancestorIDField])
			if ancestorMap[resourceType] == nil {
				ancestorMap[resourceType] = map[string]string{
					ancestorID: ancestorName,
				}
			} else {
				ancestorMap[resourceType][ancestorID] = ancestorName
			}
		}
	}

	// get instances with path info
	param := metadata.PullResourceParam{
		Condition: cond,
		Limit:     req.Page.Limit,
		Offset:    req.Page.Offset,
	}
	data, err := lgc.searchAuthResource(kit, param, req.Type)
	if err != nil {
		blog.ErrorJSON("search auth resource failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
		return nil, err
	}
	instances := make([]types.InstanceResource, 0)
	for _, inst := range data.Info {
		instance := types.InstanceResource{
			ID:          util.GetStrByInterface(inst[idField]),
			DisplayName: util.GetStrByInterface(inst[nameField]),
		}
		for _, ancestor := range filter.ResourceTypeChain {
			ancestorIDField := types.GetResourceIDField(ancestor.ID)
			// compatible for case like module that has no specified auth resource TODO confirm this
			if ancestorIDField == "" {
				ancestorIDField = common.GetInstIDField(string(ancestor.ID))
			}
			hostID, _ := util.GetInt64ByInterface(inst[idField])
			ancestorID := hostIDRelationMap[hostID][ancestorIDField]
			instance.Path = append(instance.Path, types.InstancePath{
				Type:        ancestor.ID,
				ID:          ancestorID,
				DisplayName: ancestorMap[ancestor.ID][ancestorID],
			})
		}
		instances = append(instances, instance)
	}
	return &types.ListInstanceResult{
		Count:   data.Count,
		Results: instances,
	}, nil
}
