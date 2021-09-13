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

package y3_9_202109132155

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// removeNoProcSvcInst remove service instances with no process
func removeNoProcSvcInst(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	svcInstFilter := make(map[string]interface{})
	for {
		// get service instance ids in one page by last id
		serviceInstances := make([]metadata.ServiceInstance, 0)
		err := db.Table(common.BKTableNameServiceInstance).Find(svcInstFilter).Fields(common.BKFieldID).Start(0).
			Limit(common.BKMaxPageSize).Sort(common.BKFieldID).All(ctx, &serviceInstances)
		if err != nil {
			blog.Errorf("get service instances failed, err: %v", err)
			return err
		}

		if len(serviceInstances) == 0 {
			return nil
		}

		serviceInstanceIDs := make([]int64, len(serviceInstances))
		for index, serviceInstance := range serviceInstances {
			serviceInstanceIDs[index] = serviceInstance.ID
		}

		// get service instances that has processe in it
		relationFilter := map[string]interface{}{
			common.BKServiceInstanceIDField: map[string]interface{}{common.BKDBIN: serviceInstanceIDs},
		}
		hasRelationIDs, err := db.Table(common.BKTableNameProcessInstanceRelation).Distinct(ctx,
			common.BKServiceInstanceIDField, relationFilter)
		if err != nil {
			blog.Errorf("get service instance ids that has processes failed, err: %v", err)
			return err
		}

		// delete those service instances that has no process
		hasRelationIDMap := make(map[int64]struct{})
		for _, rawID := range hasRelationIDs {
			id, err := util.GetInt64ByInterface(rawID)
			if err != nil {
				blog.Errorf("service instance id %v is invalid, err: %v", rawID, err)
				return err
			}
			hasRelationIDMap[id] = struct{}{}
		}

		deleteIDs := make([]int64, 0)
		for _, id := range serviceInstanceIDs {
			if _, exists := hasRelationIDMap[id]; !exists {
				deleteIDs = append(deleteIDs, id)
			}
		}

		if len(deleteIDs) > 0 {
			deleteFilter := map[string]interface{}{
				common.BKFieldID: map[string]interface{}{common.BKDBIN: deleteIDs},
			}
			if err := db.Table(common.BKTableNameServiceInstance).Delete(ctx, deleteFilter); err != nil {
				blog.Errorf("delete service instances(%+v) failed, err: %v", deleteIDs, err)
				return err
			}
		}

		if len(serviceInstances) < common.BKMaxPageSize {
			return nil
		}

		svcInstFilter = map[string]interface{}{
			common.BKFieldID: map[string]interface{}{common.BKDBGT: serviceInstanceIDs[len(serviceInstanceIDs)-1]},
		}
	}
}
