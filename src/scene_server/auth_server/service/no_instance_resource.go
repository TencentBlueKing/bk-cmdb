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

	"configcenter/src/ac/iam"
	"configcenter/src/common/http/rest"
	"configcenter/src/scene_server/auth_server/types"
)

// pull resource whose actions don't need to be related to instance
func (s *AuthService) PullNoRelatedInstanceResource(ctx *rest.Contexts) {
	query := new(types.PullResourceReq)
	err := ctx.DecodeInto(query)
	if err != nil {
		ctx.RespHTTPBody(types.BaseResp{
			Code:    types.InternalServerErrorCode,
			Message: err.Error(),
		})
		return
	}

	switch query.Type {
	case iam.SysSystemBase, iam.SysOperationStatistic, iam.SysAuditLog, iam.BizCustomField, iam.BizHostApply, iam.BizTopology:
	default:
		ctx.RespHTTPBody(types.BaseResp{
			Code:    types.NotFoundErrorCode,
			Message: fmt.Sprintf("resource type %s not found", query.Type),
		})
		return
	}

	switch query.Method {
	case types.ListAttrMethod:
		ctx.RespHTTPBody(types.ListAttrResourceResp{
			BaseResp: types.BaseResp{
				Code:    types.SuccessCode,
				Message: "success",
			},
			Data: []types.AttrResource{},
		})
		return
	case types.ListAttrValueMethod:
		ctx.RespHTTPBody(types.ListAttrValueResourceResp{
			BaseResp: types.BaseResp{
				Code:    types.SuccessCode,
				Message: "success",
			},
			Data: types.ListAttrValueResult{
				Count:   0,
				Results: []types.AttrValueResource{},
			},
		})
		return
	case types.ListInstanceMethod, types.ListInstanceByPolicyMethod:
		ctx.RespHTTPBody(types.ListInstanceResourceResp{
			BaseResp: types.BaseResp{
				Code:    types.SuccessCode,
				Message: "success",
			},
			Data: types.ListInstanceResult{
				Count:   0,
				Results: []types.InstanceResource{},
			},
		})
		return
	case types.FetchInstanceInfoMethod:
		ctx.RespHTTPBody(types.FetchInstanceInfoResp{
			BaseResp: types.BaseResp{
				Code:    types.SuccessCode,
				Message: "success",
			},
			Data: []map[string]interface{}{},
		})
		return
	default:
		ctx.RespHTTPBody(types.BaseResp{
			Code:    types.NotFoundErrorCode,
			Message: fmt.Sprintf("method %s not found", query.Method),
		})
		return
	}
}
