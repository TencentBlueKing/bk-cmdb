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

package handler

import (
	"fmt"

	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/extensions"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/authsynchronizer/meta"
	"configcenter/src/scene_server/admin_server/authsynchronizer/utils"
)

// HandleModuleSync do sync all plat
func (ih *IAMHandler) HandlePlatSync(task *meta.WorkRequest) error {
	businessSimplify := task.Data.(extensions.BusinessSimplify)
	header := utils.NewAPIHeaderByBusiness(&businessSimplify)
	ctx := util.NewContextFromHTTPHeader(*header)
	rid := util.GetHTTPCCRequestID(*header)

	// step1 get instances by business from core service
	plats, err := ih.authManager.CollectAllPlats(ctx, *header)
	if err != nil {
		blog.Errorf("HandlePlatSync failed, collect plat by business id failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("CollectAllPlats failed, err: %s", err.Error())
	}
	if len(plats) == 0 {
		blog.Info("no plat found")
		return nil
	}

	resources, err := ih.authManager.MakeResourcesByPlat(*header, authmeta.EmptyAction, plats...)
	if err != nil {
		blog.Errorf("HandlePlatSync failed, MakeResourcesByPlat failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("MakeResourcesByPlat failed, err: %s", err.Error())
	}
	if len(resources) == 0 && len(plats) > 0 {
		blog.Errorf("make iam resource for plat %+v return empty, rid: %s", plats, rid)
		return nil
	}
	iamResources, err := ih.authManager.Authorize.DryRunRegisterResource(ctx, resources...)
	if err != nil {
		blog.Errorf("HandleInstanceSync failed, DryRunRegisterResource failed, object: %s, instances: %+v, err: %+v, rid: %s", common.BKInnerObjIDPlat, plats, err, rid)
		return nil
	}
	if len(iamResources.Resources) == 0 {
		if blog.V(5) {
			blog.InfoJSON("HandlePlatSync failed, no cmdb resource found, skip sync for safe, resource: %s, rid: %s", resources, rid)
		}
		return nil
	}
	first := iamResources.Resources[0]
	if len(first.ResourceID) < 2 {
		blog.ErrorJSON("HandlePlatSync failed, DryRunRegisterResource result unexpected, iamResources: %s, rid: %s", iamResources, rid)
		return fmt.Errorf("DryRunRegisterResource result unexpected, layer not enough, iamResources: %+v", iamResources)
	}
	searchCondition := authcenter.SearchCondition{
		ScopeInfo: authcenter.ScopeInfo{
			ScopeType: first.ScopeType,
			ScopeID:   first.ScopeID,
		},
		ResourceType:    first.ResourceType,
		ParentResources: first.ResourceID[0 : len(first.ResourceID)-1],
	}

	taskName := "sync all plat"
	iamIDPrefix := "plat:"
	skipDeregister := false
	if err := ih.diffAndSyncInstances(*header, taskName, searchCondition, iamIDPrefix, resources, skipDeregister); err != nil {
		blog.Errorf("HandlePlatSync failed, diffAndSyncInstances failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("diffAndSyncInstances failed, err: %+v, rid: %s", err, rid)
	}
	return nil
}
