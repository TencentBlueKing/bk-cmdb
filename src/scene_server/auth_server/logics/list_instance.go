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
	"strconv"
	"strings"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/types"
)

// TODO: confirm 422 need to be used in which case

// listInstance list instances by condition
func (lgc *Logics) listInstance(kit *rest.Kit, cond map[string]interface{}, resourceType iam.TypeID, page types.Page,
	extraCond map[string]interface{}) (*types.ListInstanceResult, error) {

	idField := GetResourceIDField(resourceType)
	nameField := GetResourceNameField(resourceType)
	if idField == "" || nameField == "" {
		blog.Errorf("request type %s is invalid, rid: %s", resourceType, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "type")
	}

	if len(extraCond) > 0 {
		cond = map[string]interface{}{
			common.BKDBAND: []map[string]interface{}{cond, extraCond},
		}
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

// searchAuthResource search auth resource instances from database
func (lgc *Logics) searchAuthResource(kit *rest.Kit, param metadata.PullResourceParam, resourceType iam.TypeID) (*metadata.PullResourceResult, error) {
	param.Collection = getResourceTableName(resourceType)
	if param.Collection == "" {
		blog.Errorf("request type %s is invalid, rid: %s", resourceType, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "type")
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

// ListSystemInstance list system scope instances that have no parent
func (lgc *Logics) ListSystemInstance(kit *rest.Kit, resourceType iam.TypeID, filter *types.ListInstanceFilter,
	page types.Page, extraCond map[string]interface{}) (*types.ListInstanceResult, error) {

	cond := make(map[string]interface{})
	if filter == nil {
		return lgc.listInstance(kit, cond, resourceType, page, extraCond)
	}

	// system resource has no parent
	if filter.Parent != nil || len(filter.ResourceTypeChain) != 0 {
		return &types.ListInstanceResult{
			Count:   0,
			Results: make([]types.InstanceResource, 0),
		}, nil
	}

	for resourceType, keywords := range filter.Search {
		if resourceType != resourceType {
			return &types.ListInstanceResult{
				Count:   0,
				Results: make([]types.InstanceResource, 0),
			}, nil
		}
		keywordCond := map[string]interface{}{
			common.BKDBLIKE: strings.Join(keywords, "|"),
		}
		cond[common.BKDBOR] = []map[string]interface{}{
			{GetResourceIDField(resourceType): keywordCond},
			{GetResourceNameField(resourceType): keywordCond},
		}
	}
	return lgc.listInstance(kit, cond, resourceType, page, extraCond)
}

// ListBusinessInstance list business scope instances whose parent is biz, and has parent id field bk_biz_id in its data
func (lgc *Logics) ListBusinessInstance(kit *rest.Kit, resourceType iam.TypeID, filter *types.ListInstanceFilter,
	page types.Page) (*types.ListInstanceResult, error) {

	cond := make(map[string]interface{})
	if filter == nil {
		return lgc.listInstance(kit, cond, resourceType, page, nil)
	}

	parentType := iam.Business
	parentIDField := common.BKAppIDField
	parentNameField := common.BKAppNameField

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
		return lgc.listInstance(kit, cond, resourceType, page, nil)
	}

	// stores ancestor mapping, ancestorMap[ancestor id] = ancestor display_name
	ancestorMap := make(map[string]string)
	idField := GetResourceIDField(resourceType)
	nameField := GetResourceNameField(resourceType)

	for rscType, keywords := range filter.Search {
		keywordCond := map[string]interface{}{
			common.BKDBLIKE: strings.Join(keywords, "|"),
		}

		if rscType == resourceType {
			cond[common.BKDBOR] = []map[string]interface{}{
				{idField: keywordCond},
				{nameField: keywordCond},
			}
			continue
		}

		if rscType != parentType {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}

		// filter instance using ancestorIDs
		ancestorQuery := &metadata.QueryCondition{
			Condition: map[string]interface{}{
				parentNameField: keywordCond,
			},
			Page:   metadata.BasePage{Start: 0, Limit: common.BKNoLimit},
			Fields: []string{parentIDField, parentNameField},
		}

		result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, ancestorQuery)
		if err != nil {
			blog.ErrorJSON("get ancestor failed, err: %s, ancestorQuery: %s,rid:%s", err.Error(), ancestorQuery, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.ErrorJSON("get ancestor failed, error message: %s, ancestorQuery: %s,rid:%s", result.ErrMsg, ancestorQuery, kit.Rid)
			return nil, result.CCError()
		}

		ancestorIDs := make([]int64, len(result.Data.Info))
		for index, ancestor := range result.Data.Info {
			ancestorID, err := util.GetInt64ByInterface(ancestor[parentIDField])
			if err != nil {
				blog.ErrorJSON("parse ancestorID %s to int64 failed, error: %s, rid: %s", ancestor[parentIDField], err.Error(), kit.Rid)
				return nil, err
			}
			ancestorMap[strconv.FormatInt(ancestorID, 10)] = util.GetStrByInterface(ancestor[parentNameField])
			ancestorIDs[index] = ancestorID
		}

		if len(ancestorIDs) == 0 {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}

		cond[parentIDField] = map[string]interface{}{common.BKDBIN: ancestorIDs}
	}

	if len(filter.ResourceTypeChain) == 0 {
		return lgc.listInstance(kit, cond, resourceType, page, nil)
	}

	// get instances with path info
	param := metadata.PullResourceParam{
		Condition: cond,
		Limit:     page.Limit,
		Offset:    page.Offset,
	}
	data, err := lgc.searchAuthResource(kit, param, resourceType)
	if err != nil {
		blog.ErrorJSON("search auth resource failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
		return nil, err
	}
	if len(data.Info) == 0 {
		return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
	}

	for _, ancestor := range filter.ResourceTypeChain {
		if ancestor.ID != parentType {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}
	}

	// get all instances' ancestors if none is filtered
	if len(ancestorMap) == 0 {
		ancestorIDs := make([]int64, 0)
		for _, inst := range data.Info {
			ancestorID, err := util.GetInt64ByInterface(inst[parentIDField])
			if err != nil {
				blog.ErrorJSON("parse ancestorID %s to int64 failed, error: %s, rid: %s", inst[parentIDField], err.Error(), kit.Rid)
				return nil, err
			}
			ancestorIDs = append(ancestorIDs, ancestorID)
		}

		listBizParam := &metadata.ListWithIDOption{
			IDs:    ancestorIDs,
			Fields: []string{common.BKAppIDField, common.BKAppNameField},
		}
		bizArrStr, err := lgc.CoreAPI.CoreService().Cache().ListBusiness(kit.Ctx, kit.Header, listBizParam)
		if err != nil {
			blog.Errorf("get biz from cache failed, err: %v, bizIDs: %+v", err, ancestorIDs)
			return nil, err
		}

		if len(bizArrStr) == 0 {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}

		bizArr := make([]metadata.BizInst, 0)
		err = json.Unmarshal([]byte(bizArrStr), &bizArr)
		if err != nil {
			blog.Errorf("unmarshal biz array %s failed, err: %v", bizArrStr, err)
			return nil, err
		}

		for _, biz := range bizArr {
			ancestorMap[strconv.FormatInt(biz.BizID, 10)] = biz.BizName
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

// ListModelInstance list model instances, parent is model
func (lgc *Logics) ListModelInstance(kit *rest.Kit, resourceType iam.TypeID, filter *types.ListInstanceFilter,
	page types.Page) (*types.ListInstanceResult, error) {

	cond := make(map[string]interface{})
	if filter == nil {
		return lgc.listInstance(kit, cond, resourceType, page, nil)
	}

	expectParentType := iam.SysInstanceModel
	expectParentIDField := common.BKObjIDField
	expectParentNameField := common.BKObjNameField

	if filter.Parent != nil {
		if filter.Parent.Type != expectParentType {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}

		objID, err := lgc.getModelObjectIDWithIamParentID(kit, filter.Parent.ID)
		if err != nil {
			return nil, err
		}

		cond[expectParentIDField] = objID

		return lgc.listInstance(kit, cond, resourceType, page, nil)
	}

	// stores model mapping, modelMap[model id] = model display_name
	modelMap := make(map[string]string)
	idField := GetResourceIDField(resourceType)
	nameField := GetResourceNameField(resourceType)
	for resourceType, keywords := range filter.Search {
		keywordCond := map[string]interface{}{
			common.BKDBLIKE: strings.Join(keywords, "|"),
		}
		if resourceType == resourceType {
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
		return lgc.listInstance(kit, cond, resourceType, page, nil)
	}

	// get instances with path info
	param := metadata.PullResourceParam{
		Condition: cond,
		Limit:     page.Limit,
		Offset:    page.Offset,
	}
	data, err := lgc.searchAuthResource(kit, param, resourceType)
	if err != nil {
		blog.ErrorJSON("search auth resource failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
		return nil, err
	}
	if len(data.Info) == 0 {
		return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
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

func (lgc *Logics) getModelObjectIDWithIamParentID(kit *rest.Kit, parentID string) (string, error) {
	id, err := strconv.ParseInt(parentID, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid parent id %s type, should be a int string", parentID)
	}

	opt := &metadata.QueryCondition{
		Fields:    []string{common.BKObjIDField},
		Condition: mapstr.MapStr{common.BKFieldID: id},
	}

	result, err := lgc.CoreAPI.CoreService().Model().ReadModel(kit.Ctx, kit.Header, opt)
	if err != nil {
		return "", err
	}

	if len(result.Data.Info) != 1 {
		return "", fmt.Errorf("got multiple model with id: %s", parentID)
	}
	return result.Data.Info[0].Spec.ObjectID, nil
}

// ListHostInstance list host instances
func (lgc *Logics) ListHostInstance(kit *rest.Kit, resourceType iam.TypeID, filter *types.ListInstanceFilter,
	page types.Page) (*types.ListInstanceResult, error) {

	if resourceType != iam.Host {
		return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
	}

	if filter == nil {
		return lgc.listHostInstanceFromCache(kit, nil, page)
	}

	if filter.Parent != nil {
		if filter.Parent.Type != iam.SysHostRscPoolDirectory && filter.Parent.Type != iam.Business /* iam.Module */ {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}

		parentID, err := strconv.ParseInt(filter.Parent.ID, 10, 64)
		if err != nil {
			blog.ErrorJSON("parse parent id %s to int64 failed, error: %s, rid: %s", filter.Parent.ID, err.Error(), kit.Rid)
			return nil, err
		}

		var relationReq *metadata.HostModuleRelationRequest
		if filter.Parent.Type == iam.Business {
			relationReq = &metadata.HostModuleRelationRequest{ApplicationID: parentID, Fields: []string{common.BKHostIDField}}
		} else {
			relationReq = &metadata.HostModuleRelationRequest{ModuleIDArr: []int64{parentID}, Fields: []string{common.BKHostIDField}}
		}

		hostRsp, err := lgc.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, relationReq)
		if err != nil {
			blog.Errorf("get host ids by parent failed, err: %s, rid: %s", err.Error(), kit.Rid)
			return nil, err
		}
		if !hostRsp.Result {
			blog.Errorf("get host ids by parent failed, err code: %d, err msg: %s, rid: %s", hostRsp.Code, hostRsp.ErrMsg, kit.Rid)
			return nil, hostRsp.Error()
		}

		if len(hostRsp.Data.Info) == 0 {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}

		hostIDs := make([]int64, 0)
		for _, relation := range hostRsp.Data.Info {
			hostIDs = append(hostIDs, relation.HostID)
		}

		return lgc.listHostInstanceFromCache(kit, hostIDs, page)
	}
	return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil

	// TODO: implement host search with topo
}

type hostInstance struct {
	ID      int64  `json:"bk_host_id"`
	InnerIP string `json:"bk_host_innerip"`
}

func (lgc *Logics) listHostInstanceFromCache(kit *rest.Kit, hostIDs []int64, page types.Page) (*types.ListInstanceResult, error) {
	listHostParam := &metadata.ListHostWithPage{
		HostIDs: hostIDs,
		Fields:  []string{common.BKHostIDField, common.BKHostInnerIPField},
		Page: metadata.BasePage{
			Start: int(page.Offset),
			Limit: int(page.Limit),
		},
	}
	count, hostArrStr, err := lgc.CoreAPI.CoreService().Cache().ListHostWithPage(kit.Ctx, kit.Header, listHostParam)
	if err != nil {
		blog.Errorf("get hosts from cache failed, err: %v, hostIDs: %+v", err, hostIDs)
		return nil, err
	}

	if len(hostArrStr) == 0 {
		return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
	}

	hosts := make([]hostInstance, 0)
	err = json.Unmarshal([]byte(hostArrStr), &hosts)
	if err != nil {
		blog.Errorf("unmarshal hosts %s failed, err: %v", hostArrStr, err)
		return nil, err
	}

	instances := make([]types.InstanceResource, 0)
	for _, host := range hosts {
		instances = append(instances, types.InstanceResource{
			ID:          strconv.FormatInt(host.ID, 10),
			DisplayName: host.InnerIP,
		})
	}

	return &types.ListInstanceResult{
		Count:   count,
		Results: instances,
	}, nil
}

func (lgc *Logics) ValidateListInstanceRequest(kit *rest.Kit, req *types.PullResourceReq) (*types.ListInstanceFilter, error) {
	if req.Page.IsIllegal() {
		blog.Errorf("request page limit %d exceeds max page size, rid: %s", req.Page.Limit, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommPageLimitIsExceeded)
	}
	if req.Filter == nil {
		return nil, nil
	}
	filter, ok := req.Filter.(types.ListInstanceFilter)
	if !ok {
		blog.ErrorJSON("request filter %s is not the right type for list_instance method, rid: %s", filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "filter")
	}
	return &filter, nil
}
