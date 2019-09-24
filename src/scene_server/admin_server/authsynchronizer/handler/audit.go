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

// HandleAuditSync do sync all audit category to iam
func (ih *IAMHandler) HandleAuditSync(task *meta.WorkRequest) error {
	businessSimplify := task.Data.(extensions.BusinessSimplify)
	header := utils.NewAPIHeaderByBusiness(&businessSimplify)

	// step1 get instances by business from logics service
	categories := make([]extensions.AuditCategorySimplify, 0)
	businessID := businessSimplify.BKAppIDField
	objects, err := ih.authManager.CollectObjectsByBusinessID(context.Background(), *header, businessID)
	if err != nil {
		blog.Errorf("get categories by business id:%d failed, err: %+v", businessID, err)
		return err
	}
	for _, object := range objects {
		categories = append(categories, extensions.AuditCategorySimplify{
			BKAppIDField:    businessID,
			BKOpTargetField: object.ObjectID,
		})
	}
	if businessID != 0 {
		objects, err := ih.authManager.CollectObjectsByBusinessID(context.Background(), *header, 0)
		if err != nil {
			blog.Errorf("get objects by business id:%d failed, err: %+v", 0, err)
			return err
		}
		for _, object := range objects {
			categories = append(categories, extensions.AuditCategorySimplify{
				BKAppIDField:    businessID,
				BKOpTargetField: object.ObjectID,
			})
		}
	}
	if len(categories) == 0 {
		blog.Infof("no categories found for business: %d", businessID)
		return nil
	}
	resources, err := ih.authManager.MakeResourcesByAuditCategories(context.Background(), *header, authmeta.EmptyAction, businessID, categories...)
	if err != nil {
		blog.Errorf("make iam resource for audit categories %+v failed, err: %+v", categories, err)
		return err
	}
	if len(resources) == 0 && len(categories) > 0 {
		blog.Errorf("make iam resource for categories %+v return empty", categories)
		return nil
	}

	// step2 get set by business from iam
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.AuditLog,
		},
		BusinessID: businessSimplify.BKAppIDField,
		Layers:     make([]authmeta.Item, 0),
	}

	taskName := fmt.Sprintf("sync audit categories for business: %d", businessSimplify.BKAppIDField)
	iamIDPrefix := ""
	skipDeregister := true
	return ih.diffAndSync(taskName, rs, iamIDPrefix, resources, skipDeregister)
}
