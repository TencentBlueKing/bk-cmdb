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

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/resource/esb"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

// AuthorizeBath works to check if a user has the authority to operate resources.
func (s *AuthService) AuthorizeBatch(ctx *rest.Contexts) {
	opts := new(types.AuthBatchOptions)
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
	opts := new(types.AuthBatchOptions)
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

	iamResourceType, err := iam.ConvertResourceType(input.ResourceType, 0)
	if err != nil {
		blog.Errorf("ConvertResourceType failed, err: %+v, resourceType: %s, rid: %s", err, input.ResourceType, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	iamActionID, err := iam.ConvertResourceAction(input.ResourceType, input.Action, input.BizID)
	if err != nil {
		blog.ErrorJSON("ConvertResourceAction failed, err: %s, input: %s, rid: %s", err, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	resources := make([]types.Resource, 0)
	if input.BizID > 0 {
		businessPath := "/" + string(iam.Business) + "," + strconv.FormatInt(input.BizID, 10) + "/"
		resource := types.Resource{
			System: iam.SystemIDCMDB,
			Type:   types.ResourceType(*iamResourceType),
			Attribute: map[string]interface{}{
				types.IamPathKey: []string{businessPath},
			},
		}
		resources = append(resources, resource)
	}

	ops := &types.AuthOptions{
		System: iam.SystemIDCMDB,
		Subject: types.Subject{
			Type: "user",
			ID:   input.UserName,
		},
		Action: types.Action{
			ID: string(iamActionID),
		},
		Resources: resources,
	}
	resourceIDs, err := s.authorizer.ListAuthorizedInstances(ctx.Kit.Ctx, ops, types.ResourceType(*iamResourceType))
	if err != nil {
		blog.ErrorJSON("ListAuthorizedInstances failed, err: %+v, input: %s, ops: %s, input: %s, rid: %s", err, input, ops, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resourceIDs)
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

// RegisterResourceCreatorAction registers iam resource instance so that creator will be authorized on related actions
func (s *AuthService) RegisterResourceCreatorAction(ctx *rest.Contexts) {
	input := new(metadata.IamInstanceWithCreator)
	err := ctx.DecodeInto(input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	input.System = iam.SystemIDCMDB

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
	input.System = iam.SystemIDCMDB

	policies, err := esb.EsbClient().IamSrv().BatchRegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, *input)
	if err != nil {
		blog.ErrorJSON("register resource creator action failed, err: %s, input: %s, rid: %s", err, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(policies)
}
