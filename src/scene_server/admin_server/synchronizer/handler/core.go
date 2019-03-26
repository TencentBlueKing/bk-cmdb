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
	"strings"
	
	"configcenter/src/auth/authcenter"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common/blog"
)

func (ih *IAMHandler) diffAndSync(ra *authmeta.ResourceAttribute, resources []authmeta.ResourceAttribute) error {
	iamResources, err := ih.authManager.Authorize.ListResources(context.Background(), ra)
	if err != nil {
		blog.Errorf("synchronize set instance failed, ListResources failed, err: %+v", err)
		return err
	}

	realResources := make([]authmeta.BackendResource, 0)
	for _, iamResources := range iamResources {
		if strings.Contains(iamResources[len(iamResources)-1].ResourceID, "set") {
			realResources = append(realResources, iamResources)
		}
	}
	blog.InfoJSON("realResources is: %s", realResources)

	scope := authcenter.ScopeInfo{}
	needRegister := make([]authmeta.ResourceAttribute, 0)
	// init key:hit map for
	iamResourceKeyMap := map[string]int{}
	for _, iamResource := range realResources {
		key := generateIAMResourceKey(iamResource)
		// init hit count 0
		iamResourceKeyMap[key] = 0
	}
	
	for _, resource := range resources {
		targetResource, err := ih.authManager.Authorize.DryRunRegisterResource(context.Background(), resource)
		if err != nil {
			blog.Errorf("synchronize set instance failed, dry run register resource failed, err: %+v", err)
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
		err = ih.authManager.Authorize.RegisterResource(context.Background(), needRegister...)
		if err != nil {
			blog.ErrorJSON("sychronizer register resource that only in cmdb failed, resources: %s, err: %+v", needRegister, err)
		}
	}

	// step7 deregister resource id that hasn't been hit
	needDeregister := make([]authmeta.BackendResource, 0)
	for _, iamResource := range realResources {
		resourceKey := generateIAMResourceKey(iamResource)
		if iamResourceKeyMap[resourceKey] == 0 {
			needDeregister = append(needDeregister, iamResource)
		}
	}
	blog.V(5).Infof("needDeregister: %+v", needDeregister)
	if len(needDeregister) != 0 {
		blog.V(2).Infof("sychronize deregister resource that only in iam, resources: %+v", needDeregister)
		err = ih.authManager.Authorize.RawDeregisterResource(context.Background(), scope, needDeregister...)
		if err != nil {
			blog.ErrorJSON("sychronize deregister resource that only in iam failed, resources: %s, err: %+v", needDeregister, err)
		}
	}
	return nil
}
