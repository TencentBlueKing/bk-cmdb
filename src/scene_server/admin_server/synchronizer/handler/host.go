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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/synchronizer/meta"
	"configcenter/src/scene_server/admin_server/synchronizer/utils"
)

// HandleHostSync do sync host of one business
func (ih *IAMHandler) HandleHostSync(task *meta.WorkRequest) error {
	businessSimplify := task.Data.(meta.BusinessSimplify)
	header := utils.NewAPIHeaderByBusiness(&businessSimplify)
	coreService := ih.CoreAPI.CoreService()

	// step1 get host by business from core service
	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(businessSimplify.BKAppIDField)
	query := &metadata.QueryCondition{
		Fields:    []string{common.BKHostInnerIPField, common.BKHostIDField},
		Condition: cond.ToMapStr(),
		Limit:     metadata.SearchLimit{Limit: common.BKNoLimit},
	}
	hosts, err := coreService.Instance().ReadInstance(context.Background(), *header, common.BKInnerObjIDHost, query)
	if err != nil {
		blog.Errorf("get host:%+v by businessID:%d failed, err: %+v", businessSimplify.BKAppIDField, err)
		return fmt.Errorf("get host:%+v by businessID:%d failed, err: %+v", businessSimplify.BKAppIDField, err)
	}

	if len(hosts.Data.Info) == 0 {
		blog.V(2).Infof("business: %d has no hosts, skip synchronize hosts.", businessSimplify.BKAppIDField)
		return nil
	}

	// extract hostID
	hostIDArr := make([]int64, 0)
	for _, host := range hosts.Data.Info {
		hostIDVal, exist := host[common.BKHostIDField]
		if exist == false {
			continue
		}
		hostID, err := util.GetInt64ByInterface(hostIDVal)
		if err != nil {
			blog.V(2).Infof("synchronize task skip host:%+v, as parse hostID field failed, err: %+v", host, err)
			continue
		}
		hostIDArr = append(hostIDArr, hostID)
	}

	blog.V(4).Infof("list hosts by business:%d result: %+v", businessSimplify.BKAppIDField, hostIDArr)

	// step2 generate host layers
	businessID, batchLayers, err := extensions.GetHostLayers(coreService, header, &hostIDArr)

	// step3 generate host resource id
	resources := make([]authmeta.ResourceAttribute, 0)
	for _, layer := range batchLayers {
		lasteItem := layer[len(layer)-1]
		resource := authmeta.ResourceAttribute{
			Basic: authmeta.Basic{
				Type:       lasteItem.Type,
				Name:       lasteItem.Name,
				InstanceID: lasteItem.InstanceID,
			},
			SupplierAccount: "",
			BusinessID:      businessID,
			Layers:          layer[0:1],
		}
		resources = append(resources, resource)
	}
	desiredResources, err := ih.Authorizer.DryRunRegisterResource(context.Background(), resources...)
	if err != nil {
		blog.Errorf("synchronize host instance failed, dry run register resource faileld, err: %+v", err)
		return err
	}

	// step4 get host by business from iam
	// ListResources
	item := authmeta.Item{
		Type:       authmeta.Business,
		InstanceID: businessID,
	}
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.HostInstance,
		},
		SupplierAccount: "",
		BusinessID:      businessID,
		Layers:          []authmeta.Item{item},
	}
	realResources, err := ih.Authorizer.ListResources(context.Background(), rs)
	if err != nil {
		blog.Errorf("synchronize host instance failed, DryRunRegisterResource faileld, err: %+v", err)
		return err
	}

	// step5 diff step2 and step4 result

	// step6 register host not exist in iam

	// step7 deregister and register hosts that layers has changed

	// step8 deregister resource id that not in cmdb
	return nil
}
