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

// HandleModelSync do sync model of one business
func (ih *IAMHandler) HandleModelSync(task *meta.WorkRequest) error {
	businessSimplify := task.Data.(extensions.BusinessSimplify)
	header := utils.NewAPIHeaderByBusiness(&businessSimplify)

	// step1 get instances by business from core service
	businessID := businessSimplify.BKAppIDField
	objects, err := ih.authManager.CollectObjectsByBusinessID(context.Background(), *header, businessID)

	if err != nil {
		blog.Errorf("get models by business %d failed, err: %+v", businessSimplify.BKAppIDField, err)
		return fmt.Errorf("get models by business %d failed, err: %+v", businessSimplify.BKAppIDField, err)
	}
	blog.V(4).Infof("list model by business %d result: %+v", businessID, objects)

	resources, err := ih.authManager.MakeResourcesByObjects(context.Background(), *header, authmeta.EmptyAction, objects...)
	if err != nil {
		blog.Errorf("make iam resource from models failed, err: %+v", err)
		return nil
	}

	// step2 get models from iam
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.Model,
		},
		SupplierAccount: "",
		BusinessID:      businessSimplify.BKAppIDField,
		Layers:          make([]authmeta.Item, 0),
	}

	taskName := fmt.Sprintf("sync model for business: %d", businessSimplify.BKAppIDField)
	iamIDPrefix := ""
	skipDeregister := false
	return ih.diffAndSync(taskName, rs, iamIDPrefix, resources, skipDeregister)
}
