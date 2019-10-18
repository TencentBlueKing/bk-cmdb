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

// HandleSetTemplateSync do sync set template of one business
func (ih *IAMHandler) HandleSetTemplateSync(task *meta.WorkRequest) error {
	businessSimplify := task.Data.(extensions.BusinessSimplify)
	header := utils.NewAPIHeaderByBusiness(&businessSimplify)

	// step1 get instances by business from core service
	bizID := businessSimplify.BKAppIDField
	setTemplates, err := ih.authManager.CollectSetTemplatesByBusinessIDs(context.Background(), *header, bizID)
	if err != nil {
		blog.Errorf("collect setTemplates by business id failed, err: %+v", err)
		return err
	}
	if len(setTemplates) == 0 {
		blog.Infof("no setTemplates found for business: %d", bizID)
		return nil
	}
	resources := ih.authManager.MakeResourcesBySetTemplate(*header, authmeta.EmptyAction, bizID, setTemplates...)
	if len(resources) == 0 && len(setTemplates) > 0 {
		blog.Errorf("make iam resource for set template %+v return empty", setTemplates)
		return nil
	}

	// step2 get set template by business from iam
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.SetTemplate,
		},
		SupplierAccount: "",
		BusinessID:      businessSimplify.BKAppIDField,
		Layers:          make([]authmeta.Item, 0),
	}

	taskName := fmt.Sprintf("sync set template for business: %d", businessSimplify.BKAppIDField)
	iamIDPrefix := ""
	skipDeregister := false
	return ih.diffAndSync(taskName, rs, iamIDPrefix, resources, skipDeregister)
}
