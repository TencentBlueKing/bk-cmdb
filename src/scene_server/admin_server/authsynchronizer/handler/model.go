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

// HandleSetSync do sync set of one business
func (ih *IAMHandler) HandleModelSync(task *meta.WorkRequest) error {
	businessSimplify := task.Data.(extensions.BusinessSimplify)
	header := utils.NewAPIHeaderByBusiness(&businessSimplify)
	ctx := util.NewContextFromHTTPHeader(*header)

	// step1 get instances by business from logics service
	businessID := businessSimplify.BKAppIDField
	objects, err := ih.authManager.CollectObjectsByBusinessID(ctx, *header, businessID)
	if err != nil {
		blog.Errorf("HandleModelSync failed, get models by business %d failed, err: %+v", businessSimplify.BKAppIDField, err)
		return fmt.Errorf("get models by business %d failed, err: %+v", businessSimplify.BKAppIDField, err)
	}
	blog.V(4).Infof("HandleModelSync, list model by business %d result: %+v", businessID, objects)

	resources, err := ih.authManager.MakeResourcesByObjects(ctx, *header, authmeta.EmptyAction, objects...)
	if err != nil {
		blog.Errorf("HandleModelSync failed, make iam resource from models failed, err: %+v", err)
		return fmt.Errorf("make iam resource from models failed, err: %+v", err)
	}

	// append global model in business scope
	if businessID != 0 {
		globalModels, err := ih.authManager.CollectObjectsByBusinessID(ctx, *header, 0)
		if err != nil {
			blog.Errorf("HandleModelSync failed, get global models failed, err: %+v", err)
			return fmt.Errorf("get global models failed, err: %+v", err)
		}
		blog.V(4).Infof("HandleModelSync, list global model result: %+v", globalModels)

		globalResources, err := ih.authManager.MakeGlobalModelAsBizResources(ctx, *header, businessID, authmeta.EmptyAction, globalModels...)
		if err != nil {
			blog.Errorf("HandleModelSync failed, make global resource in biz scope failed, err: %+v", err)
			return fmt.Errorf("make global resource in biz scope failed, err: %+v", err)
		}
		resources = append(resources, globalResources...)
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
