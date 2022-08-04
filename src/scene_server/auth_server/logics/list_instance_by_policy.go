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
	"context"
	"fmt"
	"strconv"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/sdk/operator"
	sdktypes "configcenter/src/scene_server/auth_server/sdk/types"
	"configcenter/src/scene_server/auth_server/types"
)

// ListInstanceByPolicy TODO
// list resource instances that user is privileged to access by policy
func (lgc *Logics) ListInstanceByPolicy(kit *rest.Kit, resourceType iam.TypeID,
	filter *types.ListInstanceByPolicyFilter, page types.Page, extraCond map[string]interface{}) (
	*types.ListInstanceResult, error) {

	if resourceType == iam.Host {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKResourceTypeField)
	}

	collection := ""
	if iam.IsIAMSysInstance(resourceType) {
		objID, err := lgc.GetObjIDFromResourceType(kit.Ctx, kit.Header, resourceType)
		if err != nil {
			blog.ErrorJSON("get object id from resource type failed, error: %s, resource type: %s, rid: %s",
				err, resourceType, kit.Rid)
			return nil, err
		}
		collection = common.GetObjectInstTableName(objID, kit.SupplierAccount)
	} else {
		collection = getResourceTableName(resourceType)
	}

	idField := GetResourceIDField(resourceType)
	nameField := GetResourceNameField(resourceType)
	if collection == "" || idField == "" || nameField == "" {
		blog.Errorf("request type %s is invalid", resourceType)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "type")
	}

	cond, err := lgc.parseFilterToMongo(kit.Ctx, kit.Header, filter.Expression, resourceType)
	if err != nil {
		blog.ErrorJSON("parse request filter expression %s failed, error: %s, rid: %s",
			filter.Expression, err.Error(), kit.Rid)
		return nil, err
	}

	if cond == nil {
		return &types.ListInstanceResult{
			Count:   0,
			Results: make([]types.InstanceResource, 0),
		}, nil
	}

	if len(extraCond) > 0 {
		cond = map[string]interface{}{
			common.BKDBAND: []map[string]interface{}{cond, extraCond},
		}
	}

	return lgc.listInstance(kit, cond, resourceType, page)
}

// ListHostByPolicy TODO
// list host instances that user is privileged to access by policy
func (lgc *Logics) ListHostByPolicy(kit *rest.Kit, resourceType iam.TypeID, filter *types.ListInstanceByPolicyFilter,
	page types.Page) (*types.ListInstanceResult, error) {

	if resourceType != iam.Host {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKResourceTypeField)
	}

	cond, err := lgc.parseFilterToMongo(kit.Ctx, kit.Header, filter.Expression, resourceType)
	if err != nil {
		blog.ErrorJSON("parse request filter expression %s failed, error: %s, rid: %s",
			filter.Expression, err.Error(), kit.Rid)
		return nil, err
	}

	if cond == nil {
		return &types.ListInstanceResult{Count: 0, Results: make([]types.InstanceResource, 0)}, nil
	}

	param := metadata.PullResourceParam{
		Condition: cond,
		Fields:    []string{common.BKHostIDField, common.BKHostInnerIPField, common.BKCloudIDField},
		Limit:     page.Limit,
		Offset:    page.Offset,
	}

	hostRes, err := lgc.searchAuthResource(kit, param, resourceType)
	if err != nil {
		blog.ErrorJSON("search auth resource failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
		return nil, err
	}
	hosts := hostRes.Info

	cloudIDs := make([]int64, 0)
	cloudExistMap := make(map[int64]bool)
	for _, host := range hosts {
		cloudID, err := util.GetInt64ByInterface(host[common.BKCloudIDField])
		if err != nil {
			blog.Errorf("parse cloud area id failed, err: %v, host: %+v", err, host)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudIDField)
		}

		if !cloudExistMap[cloudID] {
			cloudIDs = append(cloudIDs, cloudID)
		}
	}

	cloudMap, err := lgc.getCloudNameMapByIDs(kit, cloudIDs)
	if err != nil {
		return nil, err
	}

	// covert id and display_name field
	instances := make([]types.InstanceResource, len(hosts))
	for index, host := range hosts {
		cloudID, err := util.GetInt64ByInterface(host[common.BKCloudIDField])
		if err != nil {
			blog.Errorf("parse cloud area id failed, err: %v, host: %+v", err, host)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudIDField)
		}

		instances[index] = types.InstanceResource{
			ID:          util.GetStrByInterface(host[common.BKHostIDField]),
			DisplayName: getHostDisplayName(util.GetStrByInterface(host[common.BKHostInnerIPField]), cloudMap[cloudID]),
		}
	}
	return &types.ListInstanceResult{
		Count:   hostRes.Count,
		Results: instances,
	}, nil
}

// ValidateListInstanceByPolicyRequest validate the request
func (lgc *Logics) ValidateListInstanceByPolicyRequest(kit *rest.Kit, req *types.PullResourceReq) (
	*types.ListInstanceByPolicyFilter, error) {
	filter, ok := req.Filter.(types.ListInstanceByPolicyFilter)
	if !ok {
		blog.ErrorJSON("request filter %s is not the right type for list_instance_by_policy method, rid: %s",
			req.Filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "filter")
	}

	if req.Page.IsIllegal() {
		blog.Errorf("request page limit %d exceeds max page size, rid: %s", req.Page.Limit, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommPageLimitIsExceeded)
	}
	return &filter, nil
}

// ListInstancesWithAttributes list resource instances that user is privileged to access by policy
func (lgc *Logics) ListInstancesWithAttributes(ctx context.Context, opts *sdktypes.ListWithAttributes) (
	[]string, error) {
	resourceType := iam.TypeID(opts.Type)
	rid := util.ExtractRequestIDFromContext(ctx)
	supplierAccount := util.ExtractOwnerFromContext(ctx)
	if supplierAccount == "" {
		supplierAccount = common.BKDefaultOwnerID
	}
	header := util.NewHeaderFromContext(ctx)
	collection := ""
	if iam.IsIAMSysInstance(resourceType) {
		objID, err := lgc.GetObjIDFromResourceType(ctx, header, resourceType)
		if err != nil {
			blog.Errorf("get object id from resource type(%s) failed, err: %v, rid: %s", resourceType, err, rid)
			return nil, err
		}
		collection = common.GetObjectInstTableName(objID, supplierAccount)
	} else {
		collection = getResourceTableName(resourceType)
	}
	idField := GetResourceIDField(resourceType)
	if collection == "" || idField == "" {
		return nil, fmt.Errorf("request type %s is invalid", opts.Type)
	}

	// get aggregated policy for all attribute policies
	policy, err := lgc.getAggrPolicy(opts.Operator, opts.AttrPolicies)
	if err != nil {
		blog.ErrorJSON("get aggregated policy failed, error: %s, operator:%s, attrPolicies: %s, rid: %s",
			err, opts.Operator, opts.AttrPolicies, rid)
		return nil, err
	}

	cond, err := lgc.parseFilterToMongo(ctx, header, policy, resourceType)
	if err != nil {
		blog.ErrorJSON("parse request filter expression %s failed, error: %s, rid: %s", policy, err.Error(), rid)
		return nil, err
	}
	if cond == nil {
		return make([]string, 0), nil
	}

	if len(opts.IDList) > 0 {
		idCond := make(map[string]interface{})
		if isResourceIDStringType(resourceType) {
			idCond[idField] = map[string]interface{}{common.BKDBIN: opts.IDList}
		} else {
			ids := make([]int64, len(opts.IDList))
			for idx, idStr := range opts.IDList {
				id, err := strconv.ParseInt(idStr, 10, 64)
				if err != nil {
					blog.Errorf("parse id %s to int failed, error: %v, rid: %s", idStr, err, rid)
					return nil, err
				}
				ids[idx] = id
			}
			idCond[idField] = map[string]interface{}{common.BKDBIN: ids}
		}

		cond = map[string]interface{}{common.BKDBAND: []map[string]interface{}{idCond, cond}}
	}

	param := metadata.PullResourceParam{
		Collection: collection,
		Condition:  cond,
		Fields:     []string{idField},
		Limit:      common.BKNoLimit,
	}
	res, err := lgc.CoreAPI.CoreService().Auth().SearchAuthResource(ctx, header, param)
	if err != nil {
		blog.Errorf("search auth resource failed, err: %v, param: %#v, rid: %s", err, param, rid)
		return nil, err
	}
	if err := res.CCError(); err != nil {
		blog.Errorf("search auth resource failed, err: %v, param: %#v, rid: %s", err, param, rid)
		return nil, err
	}

	idList := make([]string, 0)
	for _, instance := range res.Data.Info {
		id := util.GetStrByInterface(instance[idField])
		idList = append(idList, id)
	}
	return idList, nil
}

// getAggrPolicy get aggregated policy for all attribute policies
func (lgc *Logics) getAggrPolicy(opType operator.OperType, attrPolicies []*operator.Policy) (*operator.Policy, error) {
	if len(attrPolicies) == 0 {
		return nil, fmt.Errorf("attribute policies can't be empty")
	}

	var policy *operator.Policy
	// if operator is 'or' and 'and', the policy's element is used as type Content
	// else, the policy's element is used as type FieldValue
	if opType == operator.Or || opType == operator.And {
		policy = &operator.Policy{
			Operator: opType,
			Element: &operator.Content{
				Content: attrPolicies,
			},
		}
	} else {
		// in this condition, the length of attribute policies must be 1
		if len(attrPolicies) != 1 {
			return nil, fmt.Errorf("the length of attribute policies is not 1, policies:%#v", attrPolicies)
		}
		policy = attrPolicies[0]
	}

	return policy, nil
}
