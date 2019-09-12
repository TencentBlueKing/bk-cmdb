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

	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/authsynchronizer/meta"
	"configcenter/src/scene_server/admin_server/authsynchronizer/utils"
)

// HandleBusinessSync do sync all business to iam
func (ih *IAMHandler) HandleBusinessSync(task *meta.WorkRequest) error {
	header := utils.NewListBusinessAPIHeader()
	businesses, err := ih.authManager.CollectAllBusiness(context.Background(), *header)
	if err != nil {
		blog.Errorf("collect business failed, err: %+v", err)
		return err
	}
	blog.Infof("start sync business, count: %d", len(businesses))
	if len(businesses) == 0 {
		blog.Info("no business found")
	}

	resources := ih.authManager.MakeResourcesByBusiness(*header, authmeta.EmptyAction, businesses...)

	// step2 get businesses from iam
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.Business,
		},
	}

	taskName := fmt.Sprintf("sync all business")
	iamIDPrefix := ""
	skipDeregister := false
	return ih.diffAndSync(taskName, rs, iamIDPrefix, resources, skipDeregister)
}
