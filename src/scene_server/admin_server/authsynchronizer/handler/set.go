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

// HandleSetSync do sync set of one business
func (ih *IAMHandler) HandleSetSync(task *meta.WorkRequest) error {
	businessSimplify := task.Data.(extensions.BusinessSimplify)
	header := utils.NewAPIHeaderByBusiness(&businessSimplify)

	// step1 get instances by business from logics service
	businessID := businessSimplify.BKAppIDField
	sets, err := ih.authManager.CollectSetByBusinessID(context.Background(), *header, businessID)
	if err != nil {
		blog.Errorf("get set by business id:%d failed, err: %+v", businessID, err)
		return err
	}
	if len(sets) == 0 {
		blog.Infof("no set found for business: %d", businessID)
		return nil
	}
	resources := ih.authManager.MakeResourcesBySet(*header, authmeta.EmptyAction, businessID, sets...)
	if len(resources) == 0 && len(sets) > 0 {
		blog.Errorf("make iam resource for sets %+v return empty", sets)
		return nil
	}

	// step2 get set by business from iam
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.ModelSet,
		},
		BusinessID: businessSimplify.BKAppIDField,
		Layers:     make([]authmeta.Item, 0),
	}

	taskName := fmt.Sprintf("sync set for business: %d", businessSimplify.BKAppIDField)
	iamIDPrefix := "set:"
	skipDeregister := false
	return ih.diffAndSync(taskName, rs, iamIDPrefix, resources, skipDeregister)
}
