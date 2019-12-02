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

	"configcenter/src/auth/extensions"
	authmeta "configcenter/src/auth/meta"
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
		blog.Errorf("collect plat by business id failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("CollectAllPlats failed, err: %s", err.Error())
	}
	if len(plats) == 0 {
		blog.Info("no plat found")
		return nil
	}
	resources, err := ih.authManager.MakeResourcesByPlat(*header, authmeta.EmptyAction, plats...)
	if err != nil {
		blog.Errorf("MakeResourcesByPlat failed, err: %+v, rid: %s", err, rid)
		return fmt.Errorf("MakeResourcesByPlat failed, err: %s", err.Error())
	}
	if len(resources) == 0 && len(plats) > 0 {
		blog.Errorf("make iam resource for plat %+v return empty, rid: %s", plats, rid)
		return nil
	}

	// step2 get all plat from iam
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.Plat,
		},
		SupplierAccount: util.GetOwnerID(*header),
	}

	taskName := "sync all plat"
	iamIDPrefix := "plat:"
	skipDeregister := false
	return ih.diffAndSync(taskName, rs, iamIDPrefix, resources, skipDeregister)
}
