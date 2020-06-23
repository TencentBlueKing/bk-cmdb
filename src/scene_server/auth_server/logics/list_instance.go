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
	"strconv"
	"strings"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/types"
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

// list business scope instances which has parent id field in its data, like biz parent instances has bk_biz_id field
func (lgc *Logics) ListBusinessInstance(kit *rest.Kit, req types.PullResourceReq, parentType iam.ResourceTypeID) (*types.ListInstanceResult, error) {
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

	parentIDField := types.GetResourceIDField(parentType)
	parentNameField := types.GetResourceNameField(parentType)
	parentObjID := GetInstanceResourceObjID(parentType)
	if filter.Parent != nil {
		if filter.Parent.Type != parentType {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
		// if filter parent id is not int64 type, return empty result
		id, err := strconv.ParseInt(filter.Parent.ID, 10, 64)
		if err != nil {
			blog.Errorf("filter.parent.id %s parse int failed, error: %s, rid: %s", filter.Parent.ID, err.Error(), kit.Rid)
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
		cond[parentIDField] = id
		return lgc.ListInstance(kit, cond, req.Type, req.Page)
	}

	// stores ancestor mapping, ancestorMap[ancestor id] = ancestor display_name
	ancestorMap := make(map[string]string)
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

		if resourceType != parentType {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
		// filter ancestor topology, try to get from cache first
		// TODO get ancestor from cache
		//keys, err := lgc.cache.Keys(common.BKCacheKeyV3Prefix + parentObjID + "*").Result()
		//if err != nil {
		//	blog.ErrorJSON("get cache keys for object %s failed, error: %s, rid: %s", objID, err.Error(), kit.Rid)
		//	return nil, err
		//}
		//pipeline := lgc.cache.Pipeline()
		//for _, key := range keys {
		//	pipeline.HGetAll(key)
		//}
		//results, err := pipeline.Exec()
		//if err != nil {
		//	blog.ErrorJSON("get cached instances for object %s using keys: %v failed, error: %s, rid: %s", objID, keys, err.Error(), kit.Rid)
		//	return nil, err
		//}
		//// filter instance using ancestorIDs
		ancestorIDs := make([]int64, 0)
		//for _, result := range results {
		//	cmd := result.(*redis.StringStringMapCmd)
		//	instance, err := cmd.Result()
		//	if err != nil {
		//		blog.ErrorJSON("get cached instance result failed, error: %s, keys: %v, rid: %s", err.Error(), keys, kit.Rid)
		//		return nil, err
		//	}
		//	regex := regexp.MustCompile(fmt.Sprintf(".*%s.*", strings.Join(keywords, "|")))
		//	ancestorName := util.GetStrByInterface(instance[parentNameField])
		//	if regex.MatchString(ancestorName) {
		//		ancestorID, err := util.GetInt64ByInterface(instance[parentIDField])
		//		if err != nil {
		//			blog.ErrorJSON("parse ancestorID %s to int64 failed, error: %s, rid: %s", ancestorID, err.Error(), kit.Rid)
		//			return nil, err
		//		}
		//		ancestorMap[strconv.FormatInt(ancestorID, 10)] = ancestorName
		//		ancestorIDs = append(ancestorIDs, ancestorID)
		//	}
		//}
		ancestorQuery := &metadata.QueryCondition{
			Condition: map[string]interface{}{
				parentNameField: keywordCond,
			},
			Page:   metadata.BasePage{Start: 0, Limit: common.BKNoLimit},
			Fields: []string{parentIDField, parentNameField},
		}
		result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, parentObjID, ancestorQuery)
		if err != nil {
			blog.ErrorJSON("get ancestor failed, err: %s, ancestorQuery: %s,rid:%s", err.Error(), ancestorQuery, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.ErrorJSON("get ancestor failed, error message: %s, ancestorQuery: %s,rid:%s", result.ErrMsg, ancestorQuery, kit.Rid)
			return nil, result.CCError()
		}
		for _, ancestor := range result.Data.Info {
			ancestorID, err := util.GetInt64ByInterface(ancestor[parentIDField])
			if err != nil {
				blog.ErrorJSON("parse ancestorID %s to int64 failed, error: %s, rid: %s", ancestor[parentIDField], err.Error(), kit.Rid)
				return nil, err
			}
			ancestorMap[strconv.FormatInt(ancestorID, 10)] = util.GetStrByInterface(ancestor[parentNameField])
			ancestorIDs = append(ancestorIDs, ancestorID)
		}
		if len(ancestorIDs) > 0 {
			// business instances stores relation in field like bk_biz_id which is the ancestor's id field
			cond[parentIDField] = map[string]interface{}{common.BKDBIN: ancestorIDs}
		}
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

	// get all instances' ancestor if none is filtered
	for _, ancestor := range filter.ResourceTypeChain {
		if ancestor.ID != parentType {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
	}
	ancestorIDs := make([]int64, 0)
	for _, inst := range data.Info {
		ancestorID, err := util.GetInt64ByInterface(inst[parentIDField])
		if err != nil {
			blog.ErrorJSON("parse ancestorID %s to int64 failed, error: %s, rid: %s", inst[parentIDField], err.Error(), kit.Rid)
			return nil, err
		}
		ancestorIDs = append(ancestorIDs, ancestorID)
	}
	if len(ancestorMap) == 0 {
		ancestorQuery := &metadata.QueryCondition{
			Condition: map[string]interface{}{
				parentIDField: map[string]interface{}{
					common.BKDBIN: ancestorIDs,
				},
			},
			Page:   metadata.BasePage{Start: 0, Limit: common.BKNoLimit},
			Fields: []string{parentIDField, parentNameField},
		}
		result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, parentObjID, ancestorQuery)
		if err != nil {
			blog.ErrorJSON("get ancestor failed, err: %s, ancestorQuery: %s,rid:%s", err.Error(), ancestorQuery, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.ErrorJSON("get ancestor failed, error message: %s, ancestorQuery: %s,rid:%s", result.ErrMsg, ancestorQuery, kit.Rid)
			return nil, result.CCError()
		}
		for _, ancestor := range result.Data.Info {
			ancestorMap[util.GetStrByInterface(ancestor[parentIDField])] = util.GetStrByInterface(ancestor[parentNameField])
		}
	}
	instances := make([]types.InstanceResource, 0)
	for _, inst := range data.Info {
		instance := types.InstanceResource{
			ID:          util.GetStrByInterface(inst[idField]),
			DisplayName: util.GetStrByInterface(inst[nameField]),
		}
		ancestorID := util.GetStrByInterface(inst[parentIDField])
		instance.Path = append(instance.Path, types.InstancePath{
			Type:        parentType,
			ID:          ancestorID,
			DisplayName: ancestorMap[ancestorID],
		})
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

	expectParentType := iam.SysInstanceModel
	expectParentIDField := common.BKObjIDField
	expectParentNameField := common.BKObjNameField

	if filter.Parent != nil {
		if filter.Parent.Type != expectParentType {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
		cond[expectParentIDField] = filter.Parent.ID
		return lgc.ListInstance(kit, cond, req.Type, req.Page)
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

		if resourceType != expectParentType {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}

		queryCond := &metadata.QueryCondition{
			Fields: []string{expectParentIDField, expectParentNameField},
			Page:   metadata.BasePage{Limit: common.BKNoLimit},
			Condition: map[string]interface{}{common.BKDBOR: []map[string]interface{}{
				{expectParentIDField: keywordCond},
				{expectParentNameField: keywordCond},
			}},
		}
		modelRes, err := lgc.CoreAPI.CoreService().Model().ReadModel(kit.Ctx, kit.Header, queryCond)
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
		cond[expectParentIDField] = map[string]interface{}{
			common.BKDBIN: objIDs,
		}
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

	// get all instances' model if no model is filtered
	for _, ancestor := range filter.ResourceTypeChain {
		if ancestor.ID != expectParentType {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
	}
	modelIDs := make([]string, 0)
	for _, inst := range data.Info {
		modelID := util.GetStrByInterface(inst[expectParentIDField])
		modelIDs = append(modelIDs, modelID)
	}
	if len(modelMap) == 0 {
		modelRes, err := lgc.CoreAPI.CoreService().Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{
			Fields:    []string{expectParentIDField, expectParentNameField},
			Page:      metadata.BasePage{Limit: common.BKNoLimit},
			Condition: map[string]interface{}{expectParentIDField: map[string]interface{}{common.BKDBIN: modelIDs}},
		})
		if err != nil {
			blog.Errorf("get model failed, error: %s, rid: %s", err.Error(), kit.Rid)
			return nil, err
		}
		for _, model := range modelRes.Data.Info {
			modelMap[model.Spec.ObjectID] = model.Spec.ObjectName
		}
	}

	instances := make([]types.InstanceResource, 0)
	for _, inst := range data.Info {
		instance := types.InstanceResource{
			ID:          util.GetStrByInterface(inst[idField]),
			DisplayName: util.GetStrByInterface(inst[nameField]),
		}
		ancestorID := util.GetStrByInterface(inst[expectParentIDField])
		instance.Path = append(instance.Path, types.InstancePath{
			Type:        expectParentType,
			ID:          ancestorID,
			DisplayName: modelMap[ancestorID],
		})
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

	// TODO use cache
	// stores host relation in the for hostRelations[objID][instance ID] = host IDs
	//hostRelations := map[string]map[string][]int64{
	//	common.BKInnerObjIDApp:    make(map[string][]int64),
	//	common.BKInnerObjIDSet:    make(map[string][]int64),
	//	common.BKInnerObjIDModule: make(map[string][]int64),
	//}
	//// maps hostID to host relation
	//hostIDRelationMap := make(map[int64]map[string]string)
	//if filter.Parent != nil || filter.Search != nil {
	//	// get host topology relation from cache
	//	keys, err := lgc.cache.Keys(common.BKCacheKeyV3Prefix + "host_module:*").Result()
	//	if err != nil {
	//		blog.ErrorJSON("get cache keys for host module relation failed, error: %s, rid: %s", err.Error(), kit.Rid)
	//		return nil, err
	//	}
	//	pipeline := lgc.cache.Pipeline()
	//	for _, key := range keys {
	//		pipeline.HGetAll(key)
	//	}
	//	results, err := pipeline.Exec()
	//	if err != nil {
	//		blog.ErrorJSON("get cached host module relations using keys: %v failed, error: %s, rid: %s", keys, err.Error(), kit.Rid)
	//		return nil, err
	//	}
	//	for _, result := range results {
	//		cmd := result.(*redis.StringStringMapCmd)
	//		relation, err := cmd.Result()
	//		if err != nil {
	//			blog.ErrorJSON("get cached host module relation result failed, error: %s, rid: %s", err.Error(), kit.Rid)
	//			return nil, err
	//		}
	//		hostID, err := util.GetInt64ByInterface(relation[common.BKHostIDField])
	//		if err != nil {
	//			blog.ErrorJSON("get host module relation hostID failed, error: %s, rid: %s", err.Error(), kit.Rid)
	//			return nil, err
	//		}
	//		hostIDRelationMap[hostID] = relation
	//		hostRelations[common.BKInnerObjIDApp][relation[common.BKAppIDField]] = append(hostRelations[common.BKInnerObjIDApp][relation[common.BKAppIDField]], hostID)
	//		hostRelations[common.BKInnerObjIDSet][relation[common.BKSetIDField]] = append(hostRelations[common.BKInnerObjIDSet][relation[common.BKSetIDField]], hostID)
	//		hostRelations[common.BKInnerObjIDModule][relation[common.BKModuleIDField]] = append(hostRelations[common.BKInnerObjIDModule][relation[common.BKModuleIDField]], hostID)
	//	}
	//}
	if filter.Parent != nil {
		var relationReq *metadata.HostModuleRelationRequest
		parentID, err := strconv.ParseInt(filter.Parent.ID, 10, 64)
		if err != nil {
			blog.ErrorJSON("parse parent id %s to int64 failed, error: %s, rid: %s", filter.Parent.ID, err.Error(), kit.Rid)
			return nil, err
		}
		if req.Type != iam.Host {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
		if filter.Parent.Type != iam.SysResourcePoolDirectory && filter.Parent.Type != iam.Business /* iam.Module */ {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
		relationReq = &metadata.HostModuleRelationRequest{ModuleIDArr: []int64{parentID}}
		hostIDs := make([]int64, 0)
		hostRsp, err := lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, relationReq)
		if err != nil {
			blog.Errorf("get host ids by parent failed, err: %s, rid: %s", err.Error(), kit.Rid)
			return nil, err
		}
		if !hostRsp.Result {
			blog.Errorf("get host ids by parent failed, err code: %d, err msg: %s, rid: %s", hostRsp.Code, hostRsp.ErrMsg, kit.Rid)
			return nil, hostRsp.Error()
		}
		for _, relation := range hostRsp.Data.Info {
			hostIDs = append(hostIDs, relation.HostID)
		}
		cond[common.BKHostIDField] = map[string]interface{}{common.BKDBIN: hostIDs}
		return lgc.ListInstance(kit, cond, req.Type, req.Page)
	}
	return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil

	// TODO implement host search with topo
	//hostIDs := make([]int64, 0)
	// stores ancestor mapping, ancestorMap[ancestor type][ancestor id] = ancestor display_name
	//ancestorMap := make(map[iam.ResourceTypeID]map[string]string)
	//idField := types.GetResourceIDField(req.Type)
	//nameField := types.GetResourceNameField(req.Type)
	//for resourceType, keywords := range filter.Search {
	//	keywordCond := map[string]interface{}{
	//		common.BKDBLIKE: strings.Join(keywords, "|"),
	//	}
	//	if resourceType == req.Type {
	//		cond[common.BKDBOR] = []map[string]interface{}{
	//			{idField: keywordCond},
	//			{nameField: keywordCond},
	//		}
	//		continue
	//	}
	//
	//	// filter ancestor topology, try to get from cache first
	//	objID := GetInstanceResourceObjID(resourceType)
	//	ancestorIDField := types.GetResourceIDField(resourceType)
	//	ancestorNameField := types.GetResourceNameField(resourceType)
	//	// compatible for case like module that has no specified auth resource TODO confirm this
	//	if objID == "" {
	//		objID = string(resourceType)
	//		ancestorIDField = common.GetInstIDField(objID)
	//		ancestorNameField = common.GetInstNameField(objID)
	//	}
	//	// get ancestor from cache
	//	keys, err := lgc.cache.Keys(common.BKCacheKeyV3Prefix + objID + "*").Result()
	//	if err != nil {
	//		blog.ErrorJSON("get cache keys for object %s failed, error: %s, rid: %s", objID, err.Error(), kit.Rid)
	//		return nil, err
	//	}
	//	pipeline := lgc.cache.Pipeline()
	//	for _, key := range keys {
	//		pipeline.HGetAll(key)
	//	}
	//	results, err := pipeline.Exec()
	//	if err != nil {
	//		blog.ErrorJSON("get cached instances for object %s using keys: %v failed, error: %s, rid: %s", objID, keys, err.Error(), kit.Rid)
	//		return nil, err
	//	}
	//	// filter instance using ancestorIDs
	//	ancestorIDs := make([]string, 0)
	//	for _, result := range results {
	//		cmd := result.(*redis.StringStringMapCmd)
	//		instance, err := cmd.Result()
	//		if err != nil {
	//			blog.ErrorJSON("get cached instance result failed, error: %s, keys: %v, rid: %s", objID, keys, err.Error(), kit.Rid)
	//			return nil, err
	//		}
	//		regex := regexp.MustCompile(fmt.Sprintf(".*%s.*", strings.Join(keywords, "|")))
	//		ancestorName := util.GetStrByInterface(instance[ancestorNameField])
	//		ancestorID := util.GetStrByInterface(instance[ancestorIDField])
	//		if regex.MatchString(ancestorName) || regex.MatchString(ancestorID) {
	//			if ancestorMap[resourceType] == nil {
	//				ancestorMap[resourceType] = map[string]string{
	//					ancestorID: ancestorName,
	//				}
	//			} else {
	//				ancestorMap[resourceType][ancestorID] = ancestorName
	//			}
	//			ancestorIDs = append(ancestorIDs, instance[ancestorIDField])
	//		}
	//	}
	//	filterHostIDs := make([]int64, 0)
	//	for _, ancestorID := range ancestorIDs {
	//		filterHostIDs = append(filterHostIDs, hostRelations[objID][ancestorID]...)
	//	}
	//	hostIDs = util.IntArrIntersection(hostIDs, filterHostIDs)
	//	if len(hostIDs) == 0 {
	//		return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
	//	}
	//}
	//cond[common.BKHostIDField] = map[string]interface{}{
	//	common.BKDBIN: hostIDs,
	//}
	//
	//if len(filter.ResourceTypeChain) == 0 {
	//	return lgc.ListInstance(kit, cond, req.Type, req.Page)
	//}
	//
	//// get ancestor needed for ResourceTypeChain from cache
	//relatedIDs := make(map[string][]string)
	//for _, hostID := range hostIDs {
	//	for idField, id := range hostIDRelationMap[hostID] {
	//		relatedIDs[idField] = append(relatedIDs[idField], id)
	//	}
	//}
	//for _, ancestor := range filter.ResourceTypeChain {
	//	resourceType := ancestor.ID
	//	if ancestorMap[resourceType] != nil {
	//		continue
	//	}
	//	objID := GetInstanceResourceObjID(resourceType)
	//	ancestorIDField := types.GetResourceIDField(resourceType)
	//	ancestorNameField := types.GetResourceNameField(resourceType)
	//	if objID == "" {
	//		objID = string(resourceType)
	//		ancestorIDField = common.GetInstIDField(objID)
	//		ancestorNameField = common.GetInstNameField(objID)
	//	}
	//	pipeline := lgc.cache.Pipeline()
	//	for _, id := range relatedIDs[idField] {
	//		pipeline.HGetAll(common.BKCacheKeyV3Prefix + objID + "*id:" + id)
	//	}
	//	results, err := pipeline.Exec()
	//	if err != nil {
	//		blog.ErrorJSON("get cached instances for object %s with ids: %v failed, error: %s, rid: %s", objID, relatedIDs[idField], err.Error(), kit.Rid)
	//		return nil, err
	//	}
	//	for _, result := range results {
	//		cmd := result.(*redis.StringStringMapCmd)
	//		instance, err := cmd.Result()
	//		if err != nil {
	//			blog.ErrorJSON("get cached instance result failed, error: %s, rid: %s", objID, err.Error(), kit.Rid)
	//			return nil, err
	//		}
	//		ancestorName := util.GetStrByInterface(instance[ancestorNameField])
	//		ancestorID := util.GetStrByInterface(instance[ancestorIDField])
	//		if ancestorMap[resourceType] == nil {
	//			ancestorMap[resourceType] = map[string]string{
	//				ancestorID: ancestorName,
	//			}
	//		} else {
	//			ancestorMap[resourceType][ancestorID] = ancestorName
	//		}
	//	}
	//}
	//
	//// get instances with path info
	//param := metadata.PullResourceParam{
	//	Condition: cond,
	//	Limit:     req.Page.Limit,
	//	Offset:    req.Page.Offset,
	//}
	//data, err := lgc.searchAuthResource(kit, param, req.Type)
	//if err != nil {
	//	blog.ErrorJSON("search auth resource failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
	//	return nil, err
	//}
	//instances := make([]types.InstanceResource, 0)
	//for _, inst := range data.Info {
	//	instance := types.InstanceResource{
	//		ID:          util.GetStrByInterface(inst[idField]),
	//		DisplayName: util.GetStrByInterface(inst[nameField]),
	//	}
	//	for _, ancestor := range filter.ResourceTypeChain {
	//		ancestorIDField := types.GetResourceIDField(ancestor.ID)
	//		// compatible for case like module that has no specified auth resource TODO confirm this
	//		if ancestorIDField == "" {
	//			ancestorIDField = common.GetInstIDField(string(ancestor.ID))
	//		}
	//		hostID, _ := util.GetInt64ByInterface(inst[idField])
	//		ancestorID := hostIDRelationMap[hostID][ancestorIDField]
	//		instance.Path = append(instance.Path, types.InstancePath{
	//			Type:        ancestor.ID,
	//			ID:          ancestorID,
	//			DisplayName: ancestorMap[ancestor.ID][ancestorID],
	//		})
	//	}
	//	instances = append(instances, instance)
	//}
	//return &types.ListInstanceResult{
	//	Count:   data.Count,
	//	Results: instances,
	//}, nil
}
