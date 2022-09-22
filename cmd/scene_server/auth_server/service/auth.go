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
	"strconv"

	types2 "configcenter/cmd/scene_server/auth_server/sdk/types"
	iamtype "configcenter/pkg/ac/iam"
	"configcenter/pkg/ac/meta"
	"configcenter/pkg/blog"
	"configcenter/pkg/http/rest"
	"configcenter/pkg/metadata"
	"configcenter/pkg/resource/esb"
)

// AuthorizeBatch works to check if a user has the authority to operate resources.
func (s *AuthService) AuthorizeBatch(ctx *rest.Contexts) {
	opts := new(types2.AuthBatchOptions)
	err := ctx.DecodeInto(opts)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	decisions, err := s.authorizer.AuthorizeBatch(ctx.Kit.Ctx, opts)
	if err != nil {
		blog.ErrorJSON("authorize batch failed, err: %s, ops: %s, rid: %s", err, opts, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(decisions)
}

// AuthorizeAnyBatch works to check if a user has any authority for actions.
func (s *AuthService) AuthorizeAnyBatch(ctx *rest.Contexts) {
	opts := new(types2.AuthBatchOptions)
	err := ctx.DecodeInto(opts)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	blog.InfoJSON("-> authorize any request: %s, rid: %s", opts, ctx.Kit.Rid)

	decisions, err := s.authorizer.AuthorizeAnyBatch(ctx.Kit.Ctx, opts)
	if err != nil {
		blog.ErrorJSON("authorize any batch failed, err: %s, ops: %s, rid: %s", err, opts, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(decisions)
}

// ListAuthorizedResources returns all specified resources the user has the authority to operate.
func (s *AuthService) ListAuthorizedResources(ctx *rest.Contexts) {
	input := new(meta.ListAuthorizedResourcesParam)
	err := ctx.DecodeInto(input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	iamResourceType, err := iamtype.ConvertResourceType(input.ResourceType, 0)
	if err != nil {
		blog.Errorf("ConvertResourceType failed, err: %+v, resourceType: %s, rid: %s", err, input.ResourceType, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	iamActionID, err := iamtype.ConvertResourceAction(input.ResourceType, input.Action, input.BizID)
	if err != nil {
		blog.ErrorJSON("ConvertResourceAction failed, err: %s, input: %s, rid: %s", err, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	resources := make([]types2.Resource, 0)
	if input.BizID > 0 {
		businessPath := "/" + string(iamtype.Business) + "," + strconv.FormatInt(input.BizID, 10) + "/"
		resource := types2.Resource{
			System: iamtype.SystemIDCMDB,
			Type:   types2.ResourceType(*iamResourceType),
			Attribute: map[string]interface{}{
				types2.IamPathKey: []string{businessPath},
			},
		}
		resources = append(resources, resource)
	}

	ops := &types2.AuthOptions{
		System: iamtype.SystemIDCMDB,
		Subject: types2.Subject{
			Type: "user",
			ID:   input.UserName,
		},
		Action: types2.Action{
			ID: string(iamActionID),
		},
		Resources: resources,
	}
	authorizeList, err := s.authorizer.ListAuthorizedInstances(ctx.Kit.Ctx, ops, types2.ResourceType(*iamResourceType))
	if err != nil {
		blog.ErrorJSON("ListAuthorizedInstances failed, err: %+v,  ops: %s, input: %s, rid: %s", err, ops,
			input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(authorizeList)
}

// GetNoAuthSkipUrl returns the redirect url to iam for user to apply for specific authorizations
func (s *AuthService) GetNoAuthSkipUrl(ctx *rest.Contexts) {
	input := new(metadata.IamPermission)
	err := ctx.DecodeInto(input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	url, err := esb.EsbClient().IamSrv().GetNoAuthSkipUrl(ctx.Kit.Ctx, ctx.Kit.Header, *input)
	if err != nil {
		blog.ErrorJSON("GetNoAuthSkipUrl failed, err: %s, input: %s, rid: %s", err, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(url)
}

// GetPermissionToApply get the permissions to apply
// 用于鉴权没有通过时，根据鉴权的资源信息生成需要申请的权限信息
func (s *AuthService) GetPermissionToApply(ctx *rest.Contexts) {
	input := make([]meta.ResourceAttribute, 0)
	err := ctx.DecodeInto(&input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	permission, err := s.lgc.GetPermissionToApply(ctx.Kit, input)
	if err != nil {
		blog.ErrorJSON("GetPermissionToApply failed, err: %s, input: %s, rid: %s", err, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(permission)
}

// RegisterResourceCreatorAction registers iam resource instance so that creator will be authorized on related actions
// 创建者权限, 一个资源的创建者可以拥有这个资源的编辑和删除权限
func (s *AuthService) RegisterResourceCreatorAction(ctx *rest.Contexts) {
	input := new(metadata.IamInstanceWithCreator)
	err := ctx.DecodeInto(input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	input.System = iamtype.SystemIDCMDB

	policies, err := esb.EsbClient().IamSrv().RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, *input)
	if err != nil {
		blog.ErrorJSON("register resource creator action failed, err: %s, input: %s, rid: %s", err, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(policies)
}

// BatchRegisterResourceCreatorAction batch registers iam resource instance so that creator will be authorized on related actions
func (s *AuthService) BatchRegisterResourceCreatorAction(ctx *rest.Contexts) {
	input := new(metadata.IamInstancesWithCreator)
	err := ctx.DecodeInto(input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	input.System = iamtype.SystemIDCMDB

	policies, err := esb.EsbClient().IamSrv().BatchRegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, *input)
	if err != nil {
		blog.ErrorJSON("register resource creator action failed, err: %s, input: %s, rid: %s", err, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(policies)
}
