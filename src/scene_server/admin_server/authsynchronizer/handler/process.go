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

// HandleProcessSync do sync process of one business
func (ih *IAMHandler) HandleProcessSync(task *meta.WorkRequest) error {
	businessSimplify := task.Data.(extensions.BusinessSimplify)
	header := utils.NewAPIHeaderByBusiness(&businessSimplify)

	// step1 get instances by business from logics service
	bizID := businessSimplify.BKAppIDField
	processes, err := ih.authManager.CollectProcessesByBusinessID(context.Background(), *header, bizID)
	if err != nil {
		blog.Errorf("get processes by business %d failed, err: %+v", businessSimplify.BKAppIDField, err)
		return err
	}
	resources := ih.authManager.MakeResourcesByProcesses(*header, authmeta.EmptyAction, bizID, processes...)

	if len(resources) == 0 {
		return nil
	}

	// step2 get host by business from iam
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.Process,
		},
		BusinessID: bizID,
	}

	taskName := fmt.Sprintf("sync processes for business: %d", businessSimplify.BKAppIDField)
	iamIDPrefix := ""
	skipDeregister := false
	return ih.diffAndSync(taskName, rs, iamIDPrefix, resources, skipDeregister)
}
