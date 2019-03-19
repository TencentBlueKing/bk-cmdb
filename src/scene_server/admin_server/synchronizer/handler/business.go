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

	"configcenter/src/auth/authcenter"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/synchronizer/meta"
	"configcenter/src/scene_server/admin_server/synchronizer/utils"
)

// HandleBusinessSync do sync all business to iam
func (ih *IAMHandler) HandleBusinessSync(task *meta.WorkRequest) error {
	header := utils.NewListBusinessAPIHeader()
	condition := metadata.QueryCondition{}
	result, err := ih.CoreAPI.CoreService().Instance().ReadInstance(context.TODO(), *header, common.BKInnerObjIDApp, &condition)
	if err != nil {
		blog.Errorf("list business failed, err: %v, job: %+v", err, task)
		return err
	}

	// step1 get busineses from core service
	businessList := make([]meta.BusinessSimplify, 0)
	for _, business := range result.Data.Info {
		businessSimplify := meta.BusinessSimplify{
			BKAppIDField:      int64(business[common.BKAppIDField].(float64)),
			BKSupplierIDField: int64(business[common.BKSupplierIDField].(float64)),
			BKOwnerIDField:    business[common.BKOwnerIDField].(string),
			BKAppNameField:    business[common.BKAppNameField].(string),
		}
		// businessID := business[common.BKAppIDField].(int64)
		// businessIDArr = append(businessIDArr, businessID)
		businessList = append(businessList, businessSimplify)
	}
	blog.Info("list business businessList: %+v", businessList)

	// step2 get businesses from iam
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.Business,
		},
		SupplierAccount: "",
		// BusinessID:      businessID,
		// iam don't support host layers yet.
		// Layers:          []authmeta.Item{{Type: authmeta.Business, InstanceID: businessID,}},
	}
	realResources, err := ih.Authorizer.ListResources(context.Background(), rs)
	if err != nil {
		blog.Errorf("synchronize host instance failed, ListResources failed, err: %+v", err)
		return err
	}
	blog.InfoJSON("realResources is: %s", realResources)

	// init key:hit map for
	iamResourceKeyMap := map[string]int{}
	for _, iamResource := range realResources {
		key := generateIAMResourceKey(iamResource)
		// init hit count 0
		iamResourceKeyMap[key] = 0
	}

	// step3 register host not exist in iam
	scope := authcenter.ScopeInfo{}
	needRegister := make([]authmeta.ResourceAttribute, 0)
	for _, business := range businessList {
		resource := authmeta.ResourceAttribute{
			Basic: authmeta.Basic{
				Type:       authmeta.Business,
				Name:       business.BKAppNameField,
				InstanceID: business.BKAppIDField,
			},
			SupplierAccount: "",
			BusinessID:      business.BKAppIDField,
			// Layers:          layer[0:1],
		}
		targetResource, err := ih.Authorizer.DryRunRegisterResource(context.Background(), resource)
		if err != nil {
			blog.Errorf("synchronize host instance failed, dry run register resource failed, err: %+v", err)
			return err
		}
		if len(targetResource.Resources) != 1 {
			blog.Errorf("synchronize instance:%+v failed, dry run register result is: %+v", resource, targetResource)
			continue
		}
		scope.ScopeID = targetResource.Resources[0].ScopeID
		scope.ScopeType = targetResource.Resources[0].ScopeType
		resourceKey := generateCMDBResourceKey(&targetResource.Resources[0])
		_, exist := iamResourceKeyMap[resourceKey]
		if exist {
			iamResourceKeyMap[resourceKey]++
		} else {
			needRegister = append(needRegister, resource)
		}
	}
	blog.V(5).Infof("iamResourceKeyMap: %+v", iamResourceKeyMap)
	blog.V(5).Infof("needRegister: %+v", needRegister)
	if len(needRegister) > 0 {
		blog.V(2).Infof("sychronizer register resource that only in cmdb, resources: %+v", needRegister)
		err = ih.Authorizer.RegisterResource(context.Background(), needRegister...)
		if err != nil {
			blog.ErrorJSON("sychronizer register resource that only in cmdb failed, resources: %s, err: %+v", needRegister, err)
		}
	}

	// step4 deregister resource id that hasn't been hit
	needDeregister := make([]authmeta.BackendResource, 0)
	for _, iamResource := range realResources {
		resourceKey := generateIAMResourceKey(iamResource)
		if iamResourceKeyMap[resourceKey] == 0 {
			needDeregister = append(needDeregister, iamResource)
		}
	}
	if len(needDeregister) != 0 {
		blog.V(2).Infof("sychronizer deregister resource that only in iam, resources: %+v", needDeregister)
		err = ih.Authorizer.RawDeregisterResource(context.Background(), scope, needDeregister...)
		if err != nil {
			blog.ErrorJSON("sychronizer deregister resource that only in iam failed, resources: %s, err: %+v", needDeregister, err)
		}
	}

	return nil
}
