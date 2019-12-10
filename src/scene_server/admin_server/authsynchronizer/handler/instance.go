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

	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/extensions"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/authsynchronizer/meta"
)

// HandleSetSync do sync set of one business
func (ih *IAMHandler) HandleInstanceSync(task *meta.WorkRequest) error {
	object := task.Data.(metadata.Object)
	header := task.Header.(http.Header)
	rid := util.GetHTTPCCRequestID(header)
	ctx := util.NewContextFromHTTPHeader(header)

	ignoreObjectIDs := []string{
		common.BKInnerObjIDHost,
		common.BKInnerObjIDSet,
		common.BKInnerObjIDModule,
		common.BKInnerObjIDApp,
		common.BKInnerObjIDProc,
		common.BKInnerObjIDPlat,
	}
	objectID := object.ObjectID
	if util.InStrArr(ignoreObjectIDs, objectID) {
		blog.V(5).Infof("HandleInstanceSync, ignore instance sync task for model: %s, rid: %s", objectID, rid)
		return nil
	}

	parentResources, err := ih.authManager.MakeResourcesByObjects(ctx, header, authmeta.EmptyAction, object)
	if err != nil {
		blog.Errorf("HandleInstanceSync failed, MakeResourcesByObjects failed, make parent auth resource by objects failed, object: %+v, err: %+v, rid: %s", object, err, rid)
		return fmt.Errorf("make parent auth resource by objects failed, err: %+v", err)
	}
	if len(parentResources) != 1 {
		blog.Errorf("HandleInstanceSync failed, MakeResourcesByInstances failed, make parent auth resource by objects failed, get %d with object %s, rid: %s", len(parentResources), object.ObjectID, rid)
		return fmt.Errorf("make parent auth resource by objects failed, get %d with object %d", len(parentResources), object.ID)
	}

	parentResource := parentResources[0]
	layers := parentResource.Layers
	layers = append(layers, authmeta.Item{
		Type:       parentResource.Type,
		Action:     parentResource.Action,
		Name:       parentResource.Name,
		InstanceID: parentResource.InstanceID,
	})

	// step2. collect instances by model, and convert to iam interface format
	instances, err := ih.authManager.CollectInstancesByModelID(context.Background(), header, object.ObjectID)
	if err != nil {
		blog.Errorf("HandleInstanceSync failed, CollectInstancesByModelID failed, objectID: %s, err: %+v, rid: %s", object.ObjectID, err, rid)
		return err
	}

	bizInstanceMap := make(map[int64][]extensions.InstanceSimplify, 0)
	for _, inst := range instances {
		if _, ok := bizInstanceMap[inst.BizID]; ok == false {
			bizInstanceMap[inst.BizID] = make([]extensions.InstanceSimplify, 0)
		}
		bizInstanceMap[inst.BizID] = append(bizInstanceMap[inst.BizID], inst)
	}

	for bizID, instances := range bizInstanceMap {
		resources, err := ih.authManager.MakeResourcesByInstances(context.Background(), header, authmeta.EmptyAction, instances...)
		if err != nil {
			blog.Errorf("HandleInstanceSync failed, MakeResourcesByInstances failed, object: %s, instances: %+v, err: %+v", objectID, instances, err)
			return err
		}
		iamResources, err := ih.authManager.Authorize.DryRunRegisterResource(ctx, resources...)
		if err != nil {
			blog.Errorf("HandleInstanceSync failed, DryRunRegisterResource failed, object: %s, instances: %+v, err: %+v", objectID, instances, err)
			return nil
		}
		if len(iamResources.Resources) == 0 {
			if blog.V(5) {
				blog.InfoJSON("HandleInstanceSync failed, no cmdb resource found, skip sync for safe, %s", resources)
			}
			return nil
		}
		first := iamResources.Resources[0]
		searchCondition := authcenter.SearchCondition{
			ScopeInfo: authcenter.ScopeInfo{
				ScopeType: first.ScopeType,
				ScopeID:   first.ScopeID,
			},
			ResourceType:    first.ResourceType,
			ParentResources: first.ResourceID[0 : len(first.ResourceID)-1],
		}

		taskName := fmt.Sprintf("sync instance for business: %d model: %s", bizID, object.ObjectID)
		iamIDPrefix := ""
		skipDeregister := false
		if err := ih.diffAndSyncInstances(header, taskName, searchCondition, iamIDPrefix, resources, skipDeregister); err != nil {
			blog.Errorf("HandleInstanceSync failed, diffAndSyncInstances failed, err: %+v, rid: %s", err, rid)
			return fmt.Errorf("diffAndSyncInstances failed, err: %+v, rid: %s", err, rid)
		}
	}
	return nil
}
