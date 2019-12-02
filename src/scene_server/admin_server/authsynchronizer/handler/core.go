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
	"net/http"
	"strings"
	"time"

	"configcenter/src/auth/authcenter"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/authsynchronizer/meta"
)

func (ih *IAMHandler) getIamResources(taskName string, ra *authmeta.ResourceAttribute, iamIDPrefix string) ([]authmeta.BackendResource, error) {
	offset := int64(0)
	iamResources := make([]authmeta.BackendResource, 0)
	for {
		iamResult, err := ih.authManager.Authorize.ListPageResources(context.Background(), ra, meta.IamPageLimit, offset)
		if err != nil {
			blog.Errorf("synchronize failed, ListResources from iam failed, task: %s, err: %+v", taskName, err)
			return nil, err
		}
		iamResources = append(iamResources, iamResult.Results...)
		offset += meta.IamPageLimit
		if int64(offset) > iamResult.Count {
			break
		}
		time.Sleep(meta.IamRequestIntervalMillisecond * time.Millisecond)
	}

	realResources := make([]authmeta.BackendResource, 0)
	for _, iamResource := range iamResources {
		if len(iamResource) == 0 {
			continue
		}
		if strings.HasPrefix(iamResource[len(iamResource)-1].ResourceID, iamIDPrefix) {
			realResources = append(realResources, iamResource)
		}
	}
	return realResources, nil
}

// diffAndSyncInstances only for instances
// which can use a more extract filter condition without list all common instances from iam
func (ih *IAMHandler) diffAndSyncInstances(header http.Header, taskName string, searchCondition authcenter.SearchCondition, iamIDPrefix string, resources []authmeta.ResourceAttribute, skipDeregister bool) error {
	offset := int64(0)
	iamResources := make([]authmeta.BackendResource, 0)
	for {
		iamResult, err := ih.authManager.Authorize.RawPageListResources(context.Background(), header, searchCondition, meta.IamPageLimit, offset)
		if err != nil {
			blog.Errorf("synchronize failed, RawPageListResources from iam failed, task: %s, err: %+v", taskName, err)
			return err
		}
		iamResources = append(iamResources, iamResult.Results...)
		offset += meta.IamPageLimit
		if int64(offset) > iamResult.Count {
			break
		}
		time.Sleep(meta.IamRequestIntervalMillisecond * time.Millisecond)
	}
	if blog.V(5) {
		blog.InfoJSON("ih.authManager.Authorize.RawPageListResources, count: %d,  result: %v", len(iamResources), iamResources)
	}
	realResources := make([]authmeta.BackendResource, 0)
	for _, iamResource := range iamResources {
		if len(iamResource) == 0 {
			continue
		}
		if strings.HasPrefix(iamResource[len(iamResource)-1].ResourceID, iamIDPrefix) {
			realResources = append(realResources, iamResource)
		}
	}
	if blog.V(5) {
		blog.InfoJSON("task: %s, count: %d, iam realResources is: %v", taskName, len(realResources), realResources)
	}
	return ih.diffAndSyncCore(taskName, realResources, resources, skipDeregister)
}

func (ih *IAMHandler) diffAndSync(taskName string, ra *authmeta.ResourceAttribute, iamIDPrefix string, resources []authmeta.ResourceAttribute, skipDeregister bool) error {
	iamResources, err := ih.getIamResources(taskName, ra, iamIDPrefix)
	if err != nil {
		blog.Errorf("task: %s, get iam resources failed, err: %+v", taskName, err)
		return fmt.Errorf("get iam resources failed, err: %+v", err)
	}
	if blog.V(5) {
		blog.InfoJSON("getIamResources by %s result is: %s", ra, iamResources)
	}
	return ih.diffAndSyncCore(taskName, iamResources, resources, skipDeregister)
}

func (ih *IAMHandler) diffAndSyncCore(taskName string, iamResources []authmeta.BackendResource, resources []authmeta.ResourceAttribute, skipDeregister bool) error {
	// check final resource type related with resourceID
	dryRunResources, err := ih.authManager.Authorize.DryRunRegisterResource(context.Background(), resources...)
	if err != nil {
		blog.ErrorJSON("diffAndSyncCore failed, DryRunRegisterResource failed, %s, resources: %s, err: %s", taskName, resources, err)
		return nil
	}
	if len(dryRunResources.Resources) == 0 {
		if blog.V(5) {
			blog.InfoJSON("no cmdb resource found, skip sync for safe, %s", resources)
		}
		return nil
	}
	resourceType := dryRunResources.Resources[0].ResourceType
	if !authcenter.IsRelatedToResourceID(resourceType) {
		blog.V(5).Infof("skip-sync for resourceType: %s, as it doesn't related to resourceID", resourceType)
		return nil
	}

	scope := authcenter.ScopeInfo{}
	needRegister := make([]authmeta.ResourceAttribute, 0)
	needUpdate := make([]authmeta.ResourceAttribute, 0)
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
			iamResource, ok := iamResourceMap[resourceKey]
			if ok == false {
				continue
			}
			if len(iamResource) == 0 {
				continue
			}
			if iamResource[0].ResourceName != resource.Name {
				needUpdate = append(needUpdate, resource)
				blog.Infof("need update resource, type: %s, name: %s, id: %d", resource.Type, resource.Name, resource.InstanceID)
			}
		} else {
			needRegister = append(needRegister, resource)
			blog.Infof("need register resource, type: %s, name: %s, id: %d", resource.Type, resource.Name, resource.InstanceID)
		}
	}

	if len(needRegister) > 0 {
		blog.InfoJSON("synchronize register %d resource that only in cmdb, resources: %s", len(needRegister), needRegister)
		err := ih.authManager.Authorize.RegisterResource(context.Background(), needRegister...)
		if err != nil {
			blog.ErrorJSON("synchronize register %d resource that only in cmdb failed, resources: %s, err: %+v", len(needRegister), needRegister, err)
		}
	}

	if len(needUpdate) > 0 {
		blog.InfoJSON("synchronize update %d resource that only in cmdb, resources: %s", len(needUpdate), needUpdate)
		for _, resource := range resources {
			err := ih.authManager.Authorize.UpdateResource(context.Background(), &resource)
			if err != nil {
				blog.ErrorJSON("synchronize update resource failed, resource: %s, err: %+v", resource, err)
			}
		}
	}

	if skipDeregister == true {
		return nil
	}

	// deregister resource id that hasn't been hit
	// if len(resources) == 0 {
	// 	blog.Info("cmdb resource not found of current category, skip deregister resource for safety.")
	// 	return nil
	// }
	needDeregister := make([]authmeta.BackendResource, 0)
	for _, iamResource := range iamResources {
		resourceKey := generateIAMResourceKey(iamResource)
		if iamResourceKeyMap[resourceKey] == 0 {
			needDeregister = append(needDeregister, iamResource)
			if len(iamResource) != 0 {
				blog.Infof("need deregister, type: %s, name: %s, id: %d", iamResource[0].ResourceType,
					iamResource[0].ResourceName, iamResource[0].ResourceID)
			}
		}
	}

	if len(needDeregister) != 0 {
		blog.InfoJSON("task: %s, synchronize deregister %d resource that only in iam, resources: %s", taskName, len(needDeregister), needDeregister)
		err := ih.authManager.Authorize.RawDeregisterResource(context.Background(), scope, needDeregister...)
		if err != nil {
			blog.ErrorJSON("task: %s, synchronize deregister resource that only in iam failed, resources: %s, err: %+v", taskName, needDeregister, err)
		}
	}
	blog.Infof("%s finished.", taskName)
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
