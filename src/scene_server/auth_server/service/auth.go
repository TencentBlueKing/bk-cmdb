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
	"configcenter/src/ac/iam/permit"
	"configcenter/src/ac/meta"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

// Authorize works to check if a user has the authority to operate resources
func (s *AuthService) Authorize(ctx *rest.Contexts) {
	if !auth.IsAuthed() {
		ctx.RespEntity(meta.Decision{Authorized: true})
		return
	}

	authAttribute := new(meta.AuthAttribute)
	err := ctx.DecodeInto(authAttribute)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// filter out SkipAction, which set by api server to skip authorization
	noSkipResources := make([]meta.ResourceAttribute, 0)
	for _, resource := range authAttribute.Resources {
		if resource.Action == meta.SkipAction {
			continue
		}
		noSkipResources = append(noSkipResources, resource)
	}
	if len(noSkipResources) == 0 {
		blog.V(5).Infof("Authorize skip. auth attribute: %+v, rid: %s", authAttribute, ctx.Kit.Rid)
		ctx.RespEntity(meta.Decision{Authorized: true})
		return
	}

	resource := noSkipResources[0]
	actionID, err := iam.ConvertResourceAction(resource.Type, resource.Action, resource.BusinessID)
	if err != nil {
		blog.ErrorJSON("ConvertResourceAction failed, err: %s, resource: %s, rid: %s", err, resource, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	resources, err := iam.Adaptor(noSkipResources)
	if err != nil {
		blog.ErrorJSON("Adaptor failed, err: %s, noSkipResources: %s, rid: %s", err, noSkipResources, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ops := &types.AuthOptions{
		System: iam.SystemIDCMDB,
		Subject: types.Subject{
			Type: "user",
			ID:   authAttribute.User.UserName,
		},
		Action: types.Action{
			ID: string(actionID),
		},
		Resources: resources,
	}
	decision, err := s.authorizer.Authorize(ctx.Kit.Ctx, ops)
	if err != nil {
		blog.ErrorJSON("Authorize failed, err: %s, ops: %s, rid: %s", err, ops, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(decision)
}

// AuthorizeBath works to check if a user has the authority to operate resources.
func (s *AuthService) AuthorizeBatch(ctx *rest.Contexts) {
	authAttribute := new(meta.AuthAttribute)
	err := ctx.DecodeInto(authAttribute)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if !auth.IsAuthed() {
		decisions := make([]meta.Decision, len(authAttribute.Resources))
		for i := range decisions {
			decisions[i].Authorized = true
		}
		ctx.RespEntity(decisions)
		return
	}

	authBatchArr := make([]*types.AuthBatch, 0)
	decisions := make([]meta.Decision, len(authAttribute.Resources))
	for index, resource := range authAttribute.Resources {
		// pick out skip resource at first.
		if permit.ShouldSkipAuthorize(&resource) {
			// this resource should be skipped, do not need to verify in auth center.
			decisions[index].Authorized = true
			blog.V(5).Infof("skip authorization for resource: %+v, rid: %s", resource, ctx.Kit.Rid)
			continue
		}

		actionID, err := iam.ConvertResourceAction(resource.Type, resource.Action, resource.BusinessID)
		if err != nil {
			blog.ErrorJSON("ConvertResourceAction failed, err: %s, resource: %s, rid: %s", err, resource, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		resources, err := iam.Adaptor([]meta.ResourceAttribute{resource})
		if err != nil {
			blog.ErrorJSON("Adaptor failed, err: %s, resource: %s, rid: %s", err, resource, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		authBatchArr = append(authBatchArr, &types.AuthBatch{
			Action: types.Action{
				ID: string(actionID),
			},
			Resources: resources,
		})
	}

	if len(authBatchArr) == 0 {
		ctx.RespEntity(decisions)
		return
	}

	ops := &types.AuthBatchOptions{
		System: iam.SystemIDCMDB,
		Subject: types.Subject{
			Type: "user",
			ID:   authAttribute.User.UserName,
		},
		Batch: authBatchArr,
	}
	authDecisions, err := s.authorizer.AuthorizeBatch(ctx.Kit.Ctx, ops)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	index := 0
	for _, decision := range authDecisions {
		// skip resources' decisions are already set as authorized
		for decisions[index].Authorized {
			index++
		}
		decisions[index].Authorized = decision.Authorized
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
	resource := types.Resource{
		System: iam.SystemIDCMDB,
		Type:   types.ResourceType(*iamResourceType),
	}
	if input.BizID > 0 {
		businessPath := "/" + string(iam.Business) + "," + strconv.FormatInt(input.BizID, 10) + "/"
		pathArr := []string{businessPath}
		resource.Attribute = map[string]interface{}{
			types.IamPathKey: pathArr,
		}
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
		Resources: []types.Resource{resource},
	}
	resources, err := s.authorizer.ListAuthorizedInstances(ctx.Kit.Ctx, ops)
	if err != nil {
		blog.ErrorJSON("ListAuthorizedInstances failed, err: %+v, input: %s, ops: %s, rid: %s", err, input, ops, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resources)
}
