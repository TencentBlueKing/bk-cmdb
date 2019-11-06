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
	"context"
	"fmt"

	"configcenter/src/auth/extensions"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/authsynchronizer/meta"
	"configcenter/src/scene_server/admin_server/authsynchronizer/utils"
)

// HandleModuleSync do sync module of one business
func (ih *IAMHandler) HandleModuleSync(task *meta.WorkRequest) error {
	businessSimplify := task.Data.(extensions.BusinessSimplify)
	header := utils.NewAPIHeaderByBusiness(&businessSimplify)

	// step1 get instances by business from logics service
	bizID := businessSimplify.BKAppIDField
	modules, err := ih.authManager.CollectModuleByBusinessIDs(context.Background(), *header, bizID)
	if err != nil {
		blog.Errorf("collect module by business id failed, err: %+v", err)
		return err
	}
	if len(modules) == 0 {
		blog.Infof("no modules found for business: %d", bizID)
		return nil
	}
	resources := ih.authManager.MakeResourcesByModule(*header, authmeta.EmptyAction, bizID, modules...)
	if len(resources) == 0 && len(modules) > 0 {
		blog.Errorf("make iam resource for modules %+v return empty", modules)
		return nil
	}

	// step2 get modules by business from iam
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.ModelModule,
		},
		SupplierAccount: "",
		BusinessID:      businessSimplify.BKAppIDField,
		Layers:          make([]authmeta.Item, 0),
	}

	taskName := fmt.Sprintf("sync module for business: %d", businessSimplify.BKAppIDField)
	iamIDPrefix := "module:"
	skipDeregister := false
	return ih.diffAndSync(taskName, rs, iamIDPrefix, resources, skipDeregister)
}
