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
	"strings"

	"configcenter/src/auth/authcenter"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common/blog"
)

func (ih *IAMHandler) getIamResources(taskName string, ra *authmeta.ResourceAttribute, iamIDPrefix string) ([]authmeta.BackendResource, error) {
	iamResources, err := ih.authManager.Authorize.ListResources(context.Background(), ra)
	if err != nil {
		blog.Errorf("synchronize failed, ListResources from iam failed, task: %s, err: %+v", taskName, err)
		return nil, err
	}

	blog.V(5).Infof("ih.authManager.Authorize.ListResources result: %+v", iamResources)
	realResources := make([]authmeta.BackendResource, 0)
	for _, iamResource := range iamResources {
		if len(iamResource) == 0 {
			continue
		}
		if strings.HasPrefix(iamResource[len(iamResource)-1].ResourceID, iamIDPrefix) {
			realResources = append(realResources, iamResource)
		}
	}
	blog.InfoJSON("task: %s, realResources is: %s", taskName, realResources)
	return realResources, nil
}

func (ih *IAMHandler) diffAndSync(taskName string, ra *authmeta.ResourceAttribute, iamIDPrefix string, resources []authmeta.ResourceAttribute, skipDeregister bool) error {
	iamResources, err := ih.getIamResources(taskName, ra, iamIDPrefix)
	if err != nil {
		blog.Errorf("task: %s, get iam resources failed, err: %+v", taskName, err)
		return fmt.Errorf("get iam resources failed, err: %+v", err)
	}

	scope := authcenter.ScopeInfo{}
	needRegister := make([]authmeta.ResourceAttribute, 0)
	// init key:hit map for
	iamResourceKeyMap := map[string]int{}
	iamResourceMap := map[string]authmeta.BackendResource{}
	for _, iamResource := range iamResources {
		key := generateIAMResourceKey(iamResource)
		// init hit count 0
		iamResourceKeyMap[key] = 0
		iamResourceMap[key] = iamResource
	}

	for _, resource := range resources {
		targetResource, err := ih.authManager.Authorize.DryRunRegisterResource(context.Background(), resource)
		if err != nil {
			blog.Errorf("task: %s, synchronize set instance failed, dry run register resource failed, err: %+v", taskName, err)
			return err
		}
		if len(targetResource.Resources) != 1 {
			blog.Errorf("task: %s, synchronize instance:%+v failed, dry run register result is: %+v", taskName, resource, targetResource)
			continue
		}
		scope.ScopeID = targetResource.Resources[0].ScopeID
		scope.ScopeType = targetResource.Resources[0].ScopeType
		resourceKey := generateCMDBResourceKey(&targetResource.Resources[0])
		_, exist := iamResourceKeyMap[resourceKey]
		if exist {
			iamResourceKeyMap[resourceKey]++
			// TODO compare name and decide whether need update
			// iamResource := iamResourceMap[resourceKey]
			// resource.Name != iamResource[len(iamResource) - 1].ResourceName
		} else {
			needRegister = append(needRegister, resource)
		}
	}
	blog.V(5).Infof("task: %s, iamResourceKeyMap: %+v, needRegister: %+v", taskName, iamResourceKeyMap, needRegister)

	if len(needRegister) > 0 {
		blog.InfoJSON("synchronize register resource that only in cmdb, resources: %s", needRegister)
		err = ih.authManager.Authorize.RegisterResource(context.Background(), needRegister...)
		if err != nil {
			blog.ErrorJSON("synchronize register resource that only in cmdb failed, resources: %s, err: %+v", needRegister, err)
		}
	}

	if skipDeregister == true {
		return nil
	}

	// deregister resource id that hasn't been hit
	if len(resources) == 0 {
		blog.Info("cmdb resource not found of current category, skip deregister resource for safety.")
		return nil
	}
	needDeregister := make([]authmeta.BackendResource, 0)
	for _, iamResource := range iamResources {
		resourceKey := generateIAMResourceKey(iamResource)
		if iamResourceKeyMap[resourceKey] == 0 {
			needDeregister = append(needDeregister, iamResource)
		}
	}

	if len(needDeregister) != 0 {
		blog.V(5).Infof("task: %s, synchronize deregister resource that only in iam, resources: %+v", taskName, needDeregister)
		err = ih.authManager.Authorize.RawDeregisterResource(context.Background(), scope, needDeregister...)
		if err != nil {
			blog.ErrorJSON("task: %s, synchronize deregister resource that only in iam failed, resources: %s, err: %+v", taskName, needDeregister, err)
		}
	}
	return nil
}

func generateCMDBResourceKey(resource *authcenter.ResourceEntity) string {
	resourcesIDs := make([]string, 0)
	for _, resourceID := range resource.ResourceID {
		resourcesIDs = append(resourcesIDs, fmt.Sprintf("%s:%s", resourceID.ResourceType, resourceID.ResourceID))
	}
	key := strings.Join(resourcesIDs, "-")
	return key
}

func generateIAMResourceKey(iamResource authmeta.BackendResource) string {
	resourcesIDs := make([]string, 0)
	for _, iamLayer := range iamResource {
		resourcesIDs = append(resourcesIDs, fmt.Sprintf("%s:%s", iamLayer.ResourceType, iamLayer.ResourceID))
	}
	key := strings.Join(resourcesIDs, "-")
	return key
}
