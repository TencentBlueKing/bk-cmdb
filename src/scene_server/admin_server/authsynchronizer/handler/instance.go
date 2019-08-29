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

	ignoreObjectIDs := []string{
		common.BKInnerObjIDHost,
		common.BKInnerObjIDSet,
		common.BKInnerObjIDModule,
		common.BKInnerObjIDApp,
		common.BKInnerObjIDProc,
	}
	objectID := object.ObjectID
	if util.InStrArr(ignoreObjectIDs, objectID) {
		blog.V(5).Infof("ignore instance sync task: %s", objectID)
		return nil
	}
	// step1 construct instances resource query parameter for iam
	bizIDMap, err := ih.authManager.ExtractBusinessIDFromObjects(object)
	if err != nil {
		blog.Errorf("extract business id from model failed, model: %+v, err: %+v", object, err)
		return err
	}
	bizID := bizIDMap[object.ID]
	objectResource, err := ih.authManager.MakeResourcesByObjects(context.Background(), header, authmeta.EmptyAction, object)
	if err != nil {
		blog.Errorf("make auth resource from model failed, model: %+v, err: %+v", object, err)
		return err
	}

	layers := objectResource[0].Layers
	layers = append(layers, authmeta.Item{
		Type:       authmeta.Model,
		Name:       object.ObjectID,
		InstanceID: object.ID,
	})
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.Model,
		},
		BusinessID: bizID,
		Layers:     layers,
	}

	// step2. collect instances by model, and convert to iam interface format
	instances, err := ih.authManager.CollectInstancesByModelID(context.Background(), header, object.ObjectID)
	if err != nil {
		blog.Errorf("CollectInstancesByModelID failed, err: %+v", err)
		return err
	}
	resources, err := ih.authManager.MakeResourcesByInstances(context.Background(), header, authmeta.EmptyAction, instances...)
	if err != nil {
		blog.Errorf("diff and sync resource between iam and cmdb failed, err: %+v", err)
		return nil
	}

	taskName := fmt.Sprintf("sync instance for business: %d model: %s", bizID, object.ObjectID)
	iamIDPrefix := ""
	skipDeregister := false
	return ih.diffAndSync(taskName, rs, iamIDPrefix, resources, skipDeregister)
}
