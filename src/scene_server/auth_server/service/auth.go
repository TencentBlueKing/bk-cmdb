/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"configcenter/src/ac/iam"
	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/ac/meta"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/resource/apigw"
	apigwiam "configcenter/src/thirdparty/apigw/iam"
)

// AuthorizeBatch works to check if a user has the authority to operate resources.
func (s *AuthService) AuthorizeBatch(ctx *rest.Contexts) {
	opts := new(apigwiam.AuthBatchOptions)
	err := ctx.DecodeInto(opts)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	decisions, err := s.authorizer.AuthorizeBatch(ctx.Kit.Ctx, ctx.Kit.Header, opts)
	if err != nil {
		blog.ErrorJSON("authorize batch failed, err: %s, ops: %s, rid: %s", err, opts, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(decisions)
}

// AuthorizeAnyBatch works to check if a user has any authority for actions.
func (s *AuthService) AuthorizeAnyBatch(ctx *rest.Contexts) {
	opts := new(apigwiam.AuthBatchOptions)
	err := ctx.DecodeInto(opts)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	blog.InfoJSON("-> authorize any request: %s, rid: %s", opts, ctx.Kit.Rid)

	decisions, err := s.authorizer.AuthorizeAnyBatch(ctx.Kit.Ctx, ctx.Kit.Header, opts)
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

	iamResourceType, err := iam.ConvertResourceType(input.ResourceType, 0)
	if err != nil {
		blog.Errorf("ConvertResourceType failed, err: %+v, resourceType: %s, rid: %s", err, input.ResourceType,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	iamActionID, err := iam.ConvertResourceAction(input.ResourceType, input.Action, 0)
	if err != nil {
		blog.ErrorJSON("ConvertResourceAction failed, err: %s, input: %s, rid: %s", err, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ops := &apigwiam.AuthOptions{
		System: iamtypes.SystemIDCMDB,
		Subject: apigwiam.Subject{
			Type: "user",
			ID:   input.UserName,
		},
		Action: apigwiam.Action{
			ID: string(iamActionID),
		},
	}
	authorizeList, err := s.authorizer.ListAuthorizedInstances(ctx.Kit.Ctx, ctx.Kit.Header, ops,
		apigwiam.IamResourceType(*iamResourceType))
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

	req := &apigwiam.PermApplyURLReq{
		SystemID:    input.SystemID,
		Permissions: make([]apigwiam.PermApplyItem, 0, len(input.Actions)),
	}

	for _, action := range input.Actions {
		item := apigwiam.PermApplyItem{
			ActionID: action.ID,
		}
		for _, resType := range action.RelatedResourceTypes {
			for _, path := range resType.Instances {
				if len(path) == 0 {
					continue
				}

				self := path[len(path)-1]
				resource := apigwiam.PermApplyResource{
					ID:   self.ID,
					Type: self.Type,
				}
				if resource.Type == "" {
					resource.Type = resType.Type
				}
				if len(path) > 1 {
					resource.Ancestors = make([]apigwiam.PermApplyAncestor, 0, len(path)-1)
					for _, ancestor := range path[:len(path)-1] {
						resource.Ancestors = append(resource.Ancestors, apigwiam.PermApplyAncestor{
							ID:   ancestor.ID,
							Type: ancestor.Type,
						})
					}
				}
				item.Resources = append(item.Resources, resource)
			}
		}
		req.Permissions = append(req.Permissions, item)
	}

	url, err := apigw.Client().Iam().GetNoAuthSkipUrl(ctx.Kit.Ctx, ctx.Kit.Header, req)
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
	input.System = iamtypes.SystemIDCMDB
	policies, err := apigw.Client().Iam().RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, *input)
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
	input.System = iamtypes.SystemIDCMDB

	policies, err := apigw.Client().Iam().BatchRegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, *input)
	if err != nil {
		blog.ErrorJSON("register resource creator action failed, err: %s, input: %s, rid: %s", err, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(policies)
}
