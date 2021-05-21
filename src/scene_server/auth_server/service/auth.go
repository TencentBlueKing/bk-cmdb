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

var (
	staticActionList            []iam.ResourceAction
	staticInstanceSelectionList []iam.InstanceSelection
	staticResourceTypeList      []iam.ResourceType
	staticActionGroupList       []iam.ActionGroup
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

// RegisterModelResourceTypes add new iam resourceType to IAM.
func (s *AuthService) RegisterModelResourceTypes(ctx *rest.Contexts) {
	models := make([]metadata.Object, 0)
	err := ctx.DecodeInto(&models)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	resourceActions := iam.GenDynamicResourceTypeWithModel(models)
	// Direct call IAM.
	if err := s.acIam.Client.RegisterResourcesTypes(ctx.Kit.Ctx, resourceActions); err != nil {
		blog.ErrorJSON("register resource actions failed, error: %s, resource actions: %s, rid: %s", err.Error(), resourceActions, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// UnregisterModelResourceTypes add new iam resourceType to IAM.
func (s *AuthService) UnregisterModelResourceTypes(ctx *rest.Contexts) {
	models := make([]metadata.Object, 0)
	err := ctx.DecodeInto(&models)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	idList := []iam.TypeID{}
	for _, obj := range models {
		idList = append(idList, iam.GenDynamicResourceType(obj).ID)
	}
	// Direct call IAM.
	if err := s.acIam.Client.DeleteResourcesTypes(ctx.Kit.Ctx, idList); err != nil {
		blog.ErrorJSON("register resource actions failed, error: %s, models: %s, rid: %s", err.Error(), models, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// CreateModelInstanceActions create iam resource instance actions.
func (s *AuthService) RegisterModelInstanceSelections(ctx *rest.Contexts) {
	models := make([]metadata.Object, 0)
	err := ctx.DecodeInto(&models)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	instanceSelections := iam.GenDynamicInstanceSelectionWithModel(models)
	// Direct call IAM.
	if err := s.acIam.Client.CreateInstanceSelection(ctx.Kit.Ctx, instanceSelections); err != nil {
		blog.ErrorJSON("register instance selections failed, error: %s, instance selections: %s, rid: %s", err.Error(), instanceSelections, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// DeleteModelInstanceSelections create iam resource instance actions.
func (s *AuthService) UnregisterModelInstanceSelections(ctx *rest.Contexts) {
	models := make([]metadata.Object, 0)
	err := ctx.DecodeInto(&models)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	instanceSelectionIDList := []iam.InstanceSelectionID{}
	instanceSelections := iam.GenDynamicInstanceSelectionWithModel(models)
	for _, instanceSelection := range instanceSelections {
		instanceSelectionIDList = append(instanceSelectionIDList, instanceSelection.ID)
	}
	// Direct call IAM.
	if err := s.acIam.Client.DeleteInstanceSelection(ctx.Kit.Ctx, instanceSelectionIDList); err != nil {
		blog.ErrorJSON("Unregister instance selections failed, error: %s, instance selection ID list: %s, rid: %s", err.Error(), instanceSelectionIDList, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// CreateModelInstanceActions create iam resource instance actions.
func (s *AuthService) CreateModelInstanceActions(ctx *rest.Contexts) {
	models := make([]metadata.Object, 0)
	err := ctx.DecodeInto(&models)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	resourceActions := iam.GenModelInstanceActions(models)
	// Direct call IAM.
	if err := s.acIam.Client.CreateAction(ctx.Kit.Ctx, resourceActions); err != nil {
		blog.ErrorJSON("register resource actions failed, error: %s, resource actions: %s, rid: %s", err.Error(), resourceActions, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// DeleteModelInstanceActions delete iam resource instance actions.
func (s *AuthService) DeleteModelInstanceActions(ctx *rest.Contexts) {
	models := make([]metadata.Object, 0)
	err := ctx.DecodeInto(&models)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	actionIDList := []iam.ActionID{}
	for _, obj := range models {
		actionIDList = append(actionIDList, iam.GenDynamicActionIDListWithModel(obj)...)
	}

	// Direct call IAM.
	if err := s.acIam.Client.DeleteAction(ctx.Kit.Ctx, actionIDList); err != nil {
		blog.ErrorJSON("unregister resource actions failed, error: %s, resource actions: %s, rid: %s", err.Error(), actionIDList, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// CreateModelInstanceActionGroup create iam resource instance action group.
func (s *AuthService) UpdateModelInstanceActionGroups(ctx *rest.Contexts) {
	// 入参没有用, 由于IAM仅提供了全量更新的接口, 所以只能重新全量拉取models列表
	models, err := s.lgc.GetCustomObjects(ctx.Kit)
	if err != nil {
		blog.Errorf("Synchronize actions with IAM failed, collect notPre-models failed, err: %s, rid:%s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if staticActionGroupList == nil {
		staticActionGroupList = iam.GenerateStaticActionGroups()
	}
	actionGroups := staticActionGroupList
	// generate model instance manage action groups
	actionGroups = append(actionGroups, iam.GenModelInstanceManageActionGroups(models)...)

	// Direct call IAM.
	if err := s.acIam.Client.UpdateActionGroups(ctx.Kit.Ctx, actionGroups); err != nil {
		blog.ErrorJSON("register resource action groups failed, error: %s, resource action groups: %s, rid: %s", err.Error(), actionGroups, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// SyncIAMModelResourcesCall will call SyncIAMModelResources
func (s *AuthService) SyncIAMModelResourcesCall(ctx *rest.Contexts) {

	err := s.SyncIAMModelResources(*ctx.Kit)
	if err != nil {
		blog.Errorf("sync action with IAM failed, err:%v", err)
		ctx.RespAutoError(err)
	}
	ctx.RespEntity(nil)
}
