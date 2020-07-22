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

// pull model instance related resource
func (s *AuthService) PullInstanceResource(ctx *rest.Contexts) {
	query := types.PullResourceReq{}
	err := ctx.DecodeInto(&query)
	if err != nil {
		ctx.RespHTTPBody(types.BaseResp{
			Code:    types.InternalServerErrorCode,
			Message: err.Error(),
		})
		return
	}

	switch query.Type {
	case iam.Host, iam.Business, iam.SysCloudArea, iam.SysInstance:
	default:
		ctx.RespHTTPBody(types.BaseResp{
			Code:    types.NotFoundErrorCode,
			Message: fmt.Sprintf("resource type %s not found", query.Type),
		})
		return
	}

	switch query.Method {
	case types.ListAttrMethod:
		res, err := s.lgc.ListAttr(ctx.Kit, query.Type)
		if err != nil {
			ctx.RespHTTPBody(types.BaseResp{
				Code:    types.InternalServerErrorCode,
				Message: err.Error(),
			})
			return
		}
		ctx.RespHTTPBody(types.ListAttrResourceResp{
			BaseResp: types.SuccessBaseResp,
			Data:     res,
		})
		return
	case types.ListAttrValueMethod:
		res, err := s.lgc.ListAttrValue(ctx.Kit, query)
		if err != nil {
			ctx.RespHTTPBody(types.BaseResp{
				Code:    types.InternalServerErrorCode,
				Message: err.Error(),
			})
			return
		}
		ctx.RespHTTPBody(types.ListAttrValueResourceResp{
			BaseResp: types.SuccessBaseResp,
			Data:     *res,
		})
		return
	case types.ListInstanceMethod:
		var res *types.ListInstanceResult
		var err error
		switch query.Type {
		case iam.Host:
			res, err = s.lgc.ListHostInstance(ctx.Kit, query)
		case iam.Business, iam.SysCloudArea, iam.SysResourcePoolDirectory, iam.SysHostRscPoolDirectory:
			res, err = s.lgc.ListSystemInstance(ctx.Kit, query)
		case iam.SysInstance:
			res, err = s.lgc.ListModelInstance(ctx.Kit, query)
		}
		if err != nil {
			ctx.RespHTTPBody(types.BaseResp{
				Code:    types.InternalServerErrorCode,
				Message: err.Error(),
			})
			return
		}
		ctx.RespHTTPBody(types.ListInstanceResourceResp{
			BaseResp: types.SuccessBaseResp,
			Data:     *res,
		})
		return
	case types.FetchInstanceInfoMethod:
		res, err := s.lgc.FetchInstanceInfo(ctx.Kit, query)
		if err != nil {
			ctx.RespHTTPBody(types.BaseResp{
				Code:    types.InternalServerErrorCode,
				Message: err.Error(),
			})
			return
		}
		ctx.RespHTTPBody(types.FetchInstanceInfoResp{
			BaseResp: types.SuccessBaseResp,
			Data:     res,
		})
		return
	case types.ListInstanceByPolicyMethod:
		res, err := s.lgc.ListInstanceByPolicy(ctx.Kit, query)
		if err != nil {
			ctx.RespHTTPBody(types.BaseResp{
				Code:    types.InternalServerErrorCode,
				Message: err.Error(),
			})
			return
		}
		ctx.RespHTTPBody(types.ListInstanceResourceResp{
			BaseResp: types.SuccessBaseResp,
			Data:     *res,
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
