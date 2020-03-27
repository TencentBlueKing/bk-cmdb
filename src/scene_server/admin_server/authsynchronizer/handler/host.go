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
	"errors"
	"fmt"

	"configcenter/src/auth/extensions"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/authsynchronizer/meta"
	"configcenter/src/scene_server/admin_server/authsynchronizer/utils"
)

// HandleHostSync do sync host of one business
func (ih *IAMHandler) HandleHostSync(task *meta.WorkRequest) error {
	businessSimplify := task.Data.(extensions.BusinessSimplify)
	header := utils.NewAPIHeaderByBusiness(&businessSimplify)
	blog.Infof("sync host with biz: %d", businessSimplify.BKAppIDField)

	// step1 get instances by business from logics service
	bizID := businessSimplify.BKAppIDField
	hosts, err := ih.authManager.CollectHostByBusinessID(context.Background(), *header, bizID)
	if err != nil {
		blog.Errorf("get host by business %d failed, err: %+v", businessSimplify.BKAppIDField, err)
		return err
	}
	resources, err := ih.authManager.MakeResourcesByHosts(context.Background(), *header, authmeta.EmptyAction, hosts...)
	if err != nil {
		blog.Errorf("make host resources failed, bizID: %d, err: %+v", businessSimplify.BKAppIDField, err)
		return err
	}

	// step2 get host by business from iam
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.HostInstance,
		},
		BusinessID: bizID,
	}

	taskName := fmt.Sprintf("sync host for business: %d", businessSimplify.BKAppIDField)
	iamIDPrefix := ""
	skipDeregister := false
	return ih.diffAndSync(taskName, rs, iamIDPrefix, resources, skipDeregister)
}

func (ih *IAMHandler) HandleHostResourcePoolSync(task *meta.WorkRequest) error {
	blog.Info("sync system host instance with iam.")
	header := utils.NewListBusinessAPIHeader()
	// find only resource pool business.
	condition := metadata.QueryCondition{Condition: mapstr.MapStr{"default": 1}}
	result, err := ih.clientSet.CoreService().Instance().ReadInstance(context.TODO(), *header, common.BKInnerObjIDApp, &condition)
	if err != nil {
		return fmt.Errorf("list resource pool business failed, err: %v", err)
	}
	if len(result.Data.Info) != 1 {
		return errors.New("sync resource pool host, but can not find resource pool business")
	}
	biz, ok := result.Data.Info[0].Get(common.BKAppIDField)
	if !ok {
		return errors.New("sync resource pool host, but can not find resource pool business id")
	}

	bizID, err := util.GetInt64ByInterface(biz)
	if err != nil {
		return fmt.Errorf("sync resource pool host, but got invalid biz id: %v", biz)
	}

	// step1 get instances by business from core service
	hosts, err := ih.authManager.CollectHostByBusinessID(context.Background(), *header, bizID)
	if err != nil {
		blog.Errorf("get host by business %d failed, err: %+v", bizID, err)
		return err
	}
	blog.Infof("resource pool host in cmdb: %d", len(hosts))
	resources, err := ih.authManager.MakeResourcesByHosts(context.Background(), *header, authmeta.EmptyAction, hosts...)
	if err != nil {
		blog.Errorf("make host resources failed, bizID: %d, err: %+v", bizID, err)
		return err
	}

	// step2 get host by business from iam
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.HostInstance,
		},
		// resource pool is system resource in iam, so business id must be 0.
		BusinessID: 0,
	}

	taskName := fmt.Sprintf("sync resource pool host for business: %d", bizID)
	iamIDPrefix := ""
	skipDeregister := false
	return ih.diffAndSync(taskName, rs, iamIDPrefix, resources, skipDeregister)
}
