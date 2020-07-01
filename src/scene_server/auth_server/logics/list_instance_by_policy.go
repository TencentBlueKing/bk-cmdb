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
	"net/http"

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

// list resource instances that user is privileged to access by policy
func (lgc *Logics) ListInstanceByPolicy(kit *rest.Kit, req types.PullResourceReq) (*types.ListInstanceResult, error) {
	filter, ok := req.Filter.(types.ListInstanceByPolicyFilter)
	if !ok {
		blog.ErrorJSON("request filter %s is not the right type for list_instance_by_policy method, rid: %s", req.Filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "filter")
	}
	if req.Page.IsIllegal() {
		blog.Errorf("request page limit %d exceeds max page size, rid: %s", req.Page.Limit, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommPageLimitIsExceeded)
	}
	collection := getResourceTableName(req.Type)
	idField := GetResourceIDField(req.Type)
	nameField := GetResourceNameField(req.Type)
	if collection == "" || idField == "" || nameField == "" {
		blog.Errorf("request type %s is invalid", req.Type)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "type")
	}
	cond, err := lgc.parseFilterToMongo(kit.Ctx, kit.Header, filter.Expression, req.Type)
	if err != nil {
		blog.ErrorJSON("parse request filter expression %s failed, error: %s, rid: %s", filter.Expression, err.Error(), kit.Rid)
		return nil, err
	}
	if cond == nil {
		return &types.ListInstanceResult{
			Count:   0,
			Results: make([]types.InstanceResource, 0),
		}, nil
	}
	return lgc.ListInstance(kit, cond, req.Type, req.Page)
}

// list resource instances that user is privileged to access by policy
func (lgc *Logics) ListInstancesWithAttributes(ctx context.Context, opts *sdktypes.ListWithAttributes) ([]string, error) {
	resourceType := iam.ResourceTypeID(opts.Type)
	collection := getResourceTableName(resourceType)
	idField := GetResourceIDField(resourceType)
	if collection == "" || idField == "" {
		return nil, fmt.Errorf("request type %s is invalid", opts.Type)
	}

	policyArr := make([]*operator.Policy, len(opts.Attributes))
	for index, element := range opts.Attributes {
		policyArr[index] = &operator.Policy{
			Operator: opts.Operator,
			Element:  element,
		}
	}
	policy := &operator.Policy{
		Operator: operator.And,
		Element: &operator.Content{
			Content: policyArr,
		},
	}
	header := make(http.Header)
	header.Add(common.BKHTTPOwnerID, "0")
	header.Add(common.BKSupplierIDField, "0")
	header.Add(common.BKHTTPHeaderUser, "admin")
	header.Add("Content-Type", "application/json")
	cond, err := lgc.parseFilterToMongo(ctx, header, policy, resourceType)
	if err != nil {
		blog.ErrorJSON("parse request filter expression %s failed, error: %s", policy, err.Error())
		return nil, err
	}
	if cond == nil {
		return []string{}, nil
	}

	param := metadata.PullResourceParam{
		Condition: cond,
		Fields:    []string{idField},
		Limit:     common.BKNoLimit,
	}
	res, err := lgc.CoreAPI.CoreService().Auth().SearchAuthResource(ctx, header, param)
	if err != nil {
		blog.ErrorJSON("search auth resource failed, error: %s, param: %s", err.Error(), param)
		return nil, err
	}
	if !res.Result {
		blog.ErrorJSON("search auth resource failed, error code: %s, error message: %s, param: %s", res.Code, res.ErrMsg, param)
		return nil, res.Error()
	}

	idMap := make(map[string]bool)
	needFilterID := false
	if len(opts.IDList) > 0 {
		needFilterID = true
		for _, id := range opts.IDList {
			idMap[id] = true
		}
	}
	idList := make([]string, 0)
	for _, instance := range res.Data.Info {
		id := util.GetStrByInterface(instance[idField])
		if needFilterID && !idMap[id] {
			continue
		}
		idList = append(idList, id)
	}
	return []string{}, nil
}
