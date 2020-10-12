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

	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/scene_server/auth_server/types"
)

// PullResource iam pull resource callback function, returns resource attributes or instances based on query condition
func (s *AuthService) PullResource(ctx *rest.Contexts) {
	query := new(types.PullResourceReq)
	err := ctx.DecodeInto(query)
	if err != nil {
		ctx.RespBkError(types.InternalServerErrorCode, err.Error())
		return
	}

	method, err := s.genResourcePullMethod(ctx.Kit, query.Type)
	if err != nil {
		ctx.RespBkError(types.NotFoundErrorCode, err.Error())
		return
	}

	// get response data for each iam query method, if callback method is not set, returns empty data
	switch query.Method {
	case types.ListAttrMethod:
		if method.ListAttr == nil {
			ctx.RespBkEntity([]types.AttrResource{})
			return
		}

		res, err := method.ListAttr(ctx.Kit, query.Type)
		if err != nil {
			ctx.RespBkError(types.InternalServerErrorCode, err.Error())
			return
		}
		ctx.RespBkEntity(res)
		return

	case types.ListAttrValueMethod:
		if method.ListAttrValue == nil {
			ctx.RespBkEntity(types.ListAttrValueResult{Count: 0, Results: []types.AttrValueResource{}})
			return
		}

		filter, err := s.lgc.ValidateListAttrValueRequest(ctx.Kit, query)
		if err != nil {
			ctx.RespBkError(types.InternalServerErrorCode, err.Error())
			return
		}

		res, err := method.ListAttrValue(ctx.Kit, query.Type, filter, query.Page)
		if err != nil {
			ctx.RespBkError(types.InternalServerErrorCode, err.Error())
			return
		}
		ctx.RespBkEntity(res)
		return

	case types.ListInstanceMethod, types.SearchInstanceMethod:
		if method.ListInstance == nil {
			ctx.RespBkEntity(types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}})
			return
		}

		filter, err := s.lgc.ValidateListInstanceRequest(ctx.Kit, query)
		if err != nil {
			ctx.RespBkError(types.InternalServerErrorCode, err.Error())
			return
		}

		res, err := method.ListInstance(ctx.Kit, query.Type, filter, query.Page)
		if err != nil {
			ctx.RespBkError(types.InternalServerErrorCode, err.Error())
			return
		}
		ctx.RespBkEntity(res)

	case types.FetchInstanceInfoMethod:
		if method.FetchInstanceInfo == nil {
			ctx.RespBkEntity([]map[string]interface{}{})
			return
		}

		filter, err := s.lgc.ValidateFetchInstanceInfoRequest(ctx.Kit, query)
		if err != nil {
			ctx.RespBkError(types.InternalServerErrorCode, err.Error())
			return
		}

		if len(filter.IDs) > common.BKMaxPageSize {
			ctx.RespBkError(types.UnprocessableEntityErrorCode, fmt.Sprintf("filter.ids length exceeds maximum limit %d", common.BKMaxPageSize))
			return
		}

		res, err := method.FetchInstanceInfo(ctx.Kit, query.Type, filter)
		if err != nil {
			ctx.RespBkError(types.InternalServerErrorCode, err.Error())
			return
		}
		ctx.RespBkEntity(res)
		return

	case types.ListInstanceByPolicyMethod:
		if method.ListInstanceByPolicy == nil {
			ctx.RespBkEntity(types.ListInstanceResult{Count: 0, Results: []types.InstanceResource{}})
			return
		}

		filter, err := s.lgc.ValidateListInstanceByPolicyRequest(ctx.Kit, query)
		if err != nil {
			ctx.RespBkError(types.InternalServerErrorCode, err.Error())
			return
		}

		res, err := method.ListInstanceByPolicy(ctx.Kit, query.Type, filter, query.Page)
		if err != nil {
			ctx.RespBkError(types.InternalServerErrorCode, err.Error())
			return
		}
		ctx.RespBkEntity(res)
		return

	default:
		ctx.RespBkError(types.NotFoundErrorCode, fmt.Sprintf("method %s not found", query.Method))
		return
	}
}
