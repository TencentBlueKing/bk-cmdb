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
func (lgc *Logics) listInstance(kit *rest.Kit, cond map[string]interface{}, resourceType iam.TypeID, page types.Page) (
	*types.ListInstanceResult, error) {

	idField := GetResourceIDField(resourceType)
	nameField := GetResourceNameField(resourceType)
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

	instances := make([]types.InstanceResource, len(data.Info))
	for index, instance := range data.Info {
		instances[index] = types.InstanceResource{
			ID:          util.GetStrByInterface(instance[idField]),
			DisplayName: util.GetStrByInterface(instance[nameField]),
		}
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

	if err := res.Error(); err != nil {
		blog.ErrorJSON("search auth resource failed, error: %v, param: %s, rid: %s", err, param, kit.Rid)
		return nil, res.Error()
	}
	return &res.Data, nil
}

// ListSystemInstance list system scope instances that have no parent
func (lgc *Logics) ListSystemInstance(kit *rest.Kit, resourceType iam.TypeID, filter *types.ListInstanceFilter,
	page types.Page, extraCond map[string]interface{}) (*types.ListInstanceResult, error) {

	if filter == nil {
		return lgc.listInstance(kit, extraCond, resourceType, page)
	}

	// system resource has no parent
	if filter.Parent != nil {
		return &types.ListInstanceResult{
			Count:   0,
			Results: make([]types.InstanceResource, 0),
		}, nil
	}

	cond := make(map[string]interface{})
	if len(filter.Keyword) != 0 {
		cond[GetResourceNameField(resourceType)] = map[string]interface{}{common.BKDBLIKE: filter.Keyword, common.BKDBOPTIONS: "i"}
	}

	if len(extraCond) > 0 {
		cond = map[string]interface{}{
			common.BKDBAND: []map[string]interface{}{cond, extraCond},
		}
	}

	return lgc.listInstance(kit, cond, resourceType, page)
}

// ListBusinessInstance list business scope instances whose parent is biz, and has parent id field bk_biz_id in its data
func (lgc *Logics) ListBusinessInstance(kit *rest.Kit, resourceType iam.TypeID, filter *types.ListInstanceFilter,
	page types.Page) (*types.ListInstanceResult, error) {

	cond := make(map[string]interface{})
	if filter == nil {
		return lgc.listInstance(kit, cond, resourceType, page)
	}

	if len(filter.Keyword) != 0 {
		cond[GetResourceNameField(resourceType)] = map[string]interface{}{common.BKDBLIKE: filter.Keyword, common.BKDBOPTIONS: "i"}
	}

	if filter.Parent == nil {
		return lgc.listInstance(kit, cond, resourceType, page)
	}

	if filter.Parent.Type != iam.Business {
		return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
	}

	// if filter parent id is not int64 type, return empty result
	id, err := strconv.ParseInt(filter.Parent.ID, 10, 64)
	if err != nil {
		blog.Errorf("filter.parent.id %s parse int failed, error: %s, rid: %s", filter.Parent.ID, err.Error(), kit.Rid)
		return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
	}

	cond[common.BKAppIDField] = id
	return lgc.listInstance(kit, cond, resourceType, page)
}

// ListModelInstance list model instances, parent is model
func (lgc *Logics) ListModelInstance(kit *rest.Kit, resourceType iam.TypeID, filter *types.ListInstanceFilter,
	page types.Page) (*types.ListInstanceResult, error) {

	cond := make(map[string]interface{})
	if filter == nil {
		return lgc.listInstance(kit, cond, resourceType, page)
	}

	if len(filter.Keyword) != 0 {
		cond[common.BKInstNameField] = map[string]interface{}{common.BKDBLIKE: filter.Keyword, common.BKDBOPTIONS: "i"}
	}

	if filter.Parent == nil {
		return lgc.listInstance(kit, cond, resourceType, page)
	}

	if filter.Parent.Type != iam.SysInstanceModel {
		return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
	}

	objID, err := lgc.getModelObjectIDWithIamParentID(kit, filter.Parent.ID)
	if err != nil {
		return nil, err
	}

	cond[common.BKObjIDField] = objID

	return lgc.listInstance(kit, cond, resourceType, page)
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

	if filter.Parent == nil {
		if filter.Keyword != "" {
			return lgc.listHostInstanceFromDB(kit, nil, page, filter.Keyword)
		}
		return lgc.listHostInstanceFromCache(kit, nil, page)

	}

	if filter.Parent.Type != iam.SysHostRscPoolDirectory && filter.Parent.Type != iam.Business /* iam.Module */ {
		return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
	}

	parentID, err := strconv.ParseInt(filter.Parent.ID, 10, 64)
	if err != nil {
		blog.ErrorJSON("parse parent id %s to int64 failed, error: %s, rid: %s", filter.Parent.ID, err.Error(), kit.Rid)
		return nil, err
	}

	var relationReq *metadata.DistinctHostIDByTopoRelationRequest
	if filter.Parent.Type == iam.Business {
		relationReq = &metadata.DistinctHostIDByTopoRelationRequest{ApplicationIDArr: []int64{parentID}}
	} else {
		relationReq = &metadata.DistinctHostIDByTopoRelationRequest{ModuleIDArr: []int64{parentID}}
	}

	hostRsp, err := lgc.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(kit.Ctx, kit.Header, relationReq)
	if err != nil {
		blog.Errorf("get host ids by parent failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}
	if err := hostRsp.Error(); err != nil {
		blog.Errorf("get host ids by parent failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, hostRsp.Error()
	}

	if len(hostRsp.Data.IDArr) == 0 || int64(len(hostRsp.Data.IDArr)) <= page.Offset {
		return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
	}

	if filter.Keyword != "" {
		return lgc.listHostInstanceFromDB(kit, hostRsp.Data.IDArr, page, filter.Keyword)
	}

	return lgc.listHostInstanceFromCache(kit, hostRsp.Data.IDArr, page)
}

func (lgc *Logics) listHostInstanceFromDB(kit *rest.Kit, hostIDs []int64, page types.Page, keyword string) (*types.ListInstanceResult, error) {
	condition := make(map[string]interface{})

	if len(hostIDs) != 0 {
		condition[common.BKHostIDField] = map[string]interface{}{common.BKDBIN: hostIDs}
	}
	if keyword != "" {
		condition[common.BKHostInnerIPField] = map[string]interface{}{common.BKDBLIKE: keyword}
	}

	input := &metadata.QueryInput{
		Condition: condition,
		Fields:    common.BKHostIDField + "," + common.BKHostInnerIPField + "," + common.BKCloudIDField,
		Start:     int(page.Offset),
		Limit:     int(page.Limit),
	}

	hostResp, err := lgc.CoreAPI.CoreService().Host().GetHosts(kit.Ctx, kit.Header, input)
	if err != nil {
		return nil, err
	}

	// get cloud area to generate host display name
	cloudIDs := make([]int64, len(hostResp.Data.Info))
	for index, host := range hostResp.Data.Info {
		cloudID, err := util.GetInt64ByInterface(host[common.BKCloudIDField])
		if err != nil {
			return nil, err
		}

		cloudIDs[index] = cloudID
	}

	cloudMap, err := lgc.getCloudNameMapByIDs(kit, cloudIDs)
	if err != nil {
		return nil, err
	}

	instances := make([]types.InstanceResource, 0)
	for _, host := range hostResp.Data.Info {
		cloudID, _ := util.GetInt64ByInterface(host[common.BKCloudIDField])
		instances = append(instances, types.InstanceResource{
			ID:          util.GetStrByInterface(host[common.BKHostIDField]),
			DisplayName: getHostDisplayName(util.GetStrByInterface(host[common.BKHostInnerIPField]), cloudMap[cloudID]),
		})
	}

	return &types.ListInstanceResult{
		Count:   int64(hostResp.Data.Count),
		Results: instances,
	}, nil
}

type hostInstance struct {
	ID      int64  `json:"bk_host_id"`
	InnerIP string `json:"bk_host_innerip"`
	CloudID int64  `json:"bk_cloud_id"`
}

func (lgc *Logics) listHostInstanceFromCache(kit *rest.Kit, hostIDs []int64, page types.Page) (*types.ListInstanceResult, error) {

	// if hostIDs are set, get hosts from cache returns hosts using ids directly without paging, we need to do it here
	hosts := make([]hostInstance, 0)
	var count int64

	hostLen := int64(len(hostIDs))
	if hostLen > 0 {
		count = hostLen

		hostIDLen := page.Offset + page.Limit
		if hostIDLen > hostLen {
			hostIDLen = hostLen
		}

		for offset := page.Offset; offset < hostIDLen; offset += 500 {
			limit := offset + 500
			if limit > hostIDLen {
				limit = hostIDLen
			}
			listHostParam := &metadata.ListWithIDOption{
				IDs:    hostIDs[offset:limit],
				Fields: []string{common.BKHostIDField, common.BKHostInnerIPField},
			}
			hostArrStr, err := lgc.CoreAPI.CacheService().Cache().Host().ListHostWithHostID(kit.Ctx, kit.Header, listHostParam)
			if err != nil {
				blog.Errorf("get hosts from cache failed, err: %v, hostIDs: %+v", err, hostIDs)
				return nil, err
			}

			hostArr := make([]hostInstance, 0)
			err = json.Unmarshal([]byte(hostArrStr), &hostArr)
			if err != nil {
				blog.Errorf("unmarshal hosts %s failed, err: %v", hostArrStr, err)
				return nil, err
			}

			hosts = append(hosts, hostArr...)
		}
	} else {
		listHostParam := &metadata.ListHostWithPage{
			Fields: []string{common.BKHostIDField, common.BKHostInnerIPField},
			Page: metadata.BasePage{
				Start: int(page.Offset),
				Limit: int(page.Limit),
			},
		}

		cnt, hostArrStr, err := lgc.CoreAPI.CacheService().Cache().Host().ListHostWithPage(kit.Ctx, kit.Header, listHostParam)
		if err != nil {
			blog.Errorf("get hosts from cache failed, err: %v, hostIDs: %+v", err, hostIDs)
			return nil, err
		}

		if len(hostArrStr) == 0 {
			return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
		}

		err = json.Unmarshal([]byte(hostArrStr), &hosts)
		if err != nil {
			blog.Errorf("unmarshal hosts %s failed, err: %v", hostArrStr, err)
			return nil, err
		}

		count = cnt
	}

	// get cloud area to generate host display name
	cloudIDs := make([]int64, len(hosts))
	for index, host := range hosts {
		cloudIDs[index] = host.CloudID
	}

	cloudMap, err := lgc.getCloudNameMapByIDs(kit, cloudIDs)
	if err != nil {
		return nil, err
	}

	instances := make([]types.InstanceResource, 0)
	for _, host := range hosts {
		instances = append(instances, types.InstanceResource{
			ID:          strconv.FormatInt(host.ID, 10),
			DisplayName: getHostDisplayName(host.InnerIP, cloudMap[host.CloudID]),
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
