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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
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
	idField := types.GetResourceIDField(req.Type)
	nameField := types.GetResourceNameField(req.Type)
	if collection == "" || idField == "" || nameField == "" {
		blog.Errorf("request type %s is invalid", req.Type)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "type")
	}
	cond, err := types.ParseFilterToMongo(filter.Expression, req.Type)
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
