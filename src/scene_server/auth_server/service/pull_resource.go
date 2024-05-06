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
	"regexp"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/scene_server/auth_server/types"
)

// PullResource iam pull resource callback function, returns resource attributes or instances based on query condition
func (s *AuthService) PullResource(ctx *rest.Contexts) {
	start := time.Now()

	query := new(types.PullResourceReq)
	err := ctx.DecodeInto(query)
	if err != nil {
		ctx.RespBkError(types.InternalServerErrorCode, err.Error())
		return
	}

	defer func() {
		if time.Since(start) > time.Second {
			blog.V(4).Infof("[iam pull resource] request exceeded max latency time, cost: %dms, user: %s, body: %#v, "+
				"rid: %s", time.Since(start)/time.Millisecond, ctx.Kit.User, query, ctx.Kit.Rid)
		}
	}()

	method, err := s.genResourcePullMethod(ctx.Kit, query.Type)
	if err != nil {
		ctx.RespBkError(types.NotFoundErrorCode, err.Error())
		return
	}

	// get response data for each iam query method, if callback method is not set, returns empty data
	var res interface{}
	switch query.Method {
	case types.ListAttrMethod:
		res, err = s.listAttr(ctx.Kit, method, query)
	case types.ListAttrValueMethod:
		res, err = s.listAttrValue(ctx.Kit, method, query)
	case types.ListInstanceMethod, types.SearchInstanceMethod:
		res, err = s.listInstance(ctx.Kit, method, query)
	case types.FetchInstanceInfoMethod:
		res, err = s.fetchInstanceInfo(ctx.Kit, method, query)
	case types.ListInstanceByPolicyMethod:
		res, err = s.listInstanceByPolicy(ctx.Kit, method, query)
	default:
		ctx.RespBkError(types.NotFoundErrorCode, fmt.Sprintf("method %s not found", query.Method))
		return
	}

	if err != nil {
		ctx.RespBkError(types.InternalServerErrorCode, err.Error())
		return
	}

	ctx.RespBkEntity(res)
}

func (s *AuthService) listAttr(kit *rest.Kit, method types.ResourcePullMethod, query *types.PullResourceReq) (
	[]types.AttrResource, error) {

	if method.ListAttr == nil {
		return make([]types.AttrResource, 0), nil
	}

	res, err := method.ListAttr(kit, query.Type)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *AuthService) listAttrValue(kit *rest.Kit, method types.ResourcePullMethod, query *types.PullResourceReq) (
	*types.ListAttrValueResult, error) {

	if method.ListAttrValue == nil {
		return &types.ListAttrValueResult{Count: 0, Results: []types.AttrValueResource{}}, nil
	}

	filter, err := s.lgc.ValidateListAttrValueRequest(kit, query)
	if err != nil {
		return nil, err
	}

	res, err := method.ListAttrValue(kit, query.Type, filter, query.Page)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *AuthService) listInstance(kit *rest.Kit, method types.ResourcePullMethod, query *types.PullResourceReq) (
	*types.ListInstanceResult, error) {

	if method.ListInstance == nil {
		return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
	}

	filter, err := s.lgc.ValidateListInstanceRequest(kit, query)
	if err != nil {
		return nil, err
	}

	if filter.Keyword != "" {
		filter.Keyword = regexp.QuoteMeta(filter.Keyword)
	}

	res, err := method.ListInstance(kit, query.Type, filter, query.Page)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *AuthService) fetchInstanceInfo(kit *rest.Kit, method types.ResourcePullMethod, query *types.PullResourceReq) (
	[]map[string]interface{}, error) {

	if method.FetchInstanceInfo == nil {
		return make([]map[string]interface{}, 0), nil
	}

	filter, err := s.lgc.ValidateFetchInstanceInfoRequest(kit, query)
	if err != nil {
		return nil, err
	}

	if len(filter.IDs) > common.BKMaxPageSize {
		return nil, fmt.Errorf("filter.ids length exceeds maximum limit %d", common.BKMaxPageSize)
	}

	res, err := method.FetchInstanceInfo(kit, query.Type, filter)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *AuthService) listInstanceByPolicy(kit *rest.Kit, method types.ResourcePullMethod,
	query *types.PullResourceReq) (*types.ListInstanceResult, error) {

	if method.ListInstanceByPolicy == nil {
		return &types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}}, nil
	}

	filter, err := s.lgc.ValidateListInstanceByPolicyRequest(kit, query)
	if err != nil {
		return nil, err
	}

	res, err := method.ListInstanceByPolicy(kit, query.Type, filter, query.Page)
	if err != nil {
		return nil, err
	}
	return res, nil
}
