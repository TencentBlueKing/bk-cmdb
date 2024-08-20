/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package y3_13_202407191507

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/kube/types"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	dbtypes "configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func upgradeContainer(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	if err := addContainerIndex(ctx, db); err != nil {
		blog.Errorf("add container table index failed, err: %v", err)
		return err
	}

	if err := upgradeContainerData(ctx, db); err != nil {
		blog.Errorf("upgrade container data failed, err: %v", err)
		return err
	}

	blog.Infof("upgrade container data successfully")

	return nil
}

func addContainerIndex(ctx context.Context, db dal.RDB) error {
	indexes := []dbtypes.Index{
		{
			Name: common.CCLogicIndexNamePrefix + "biz_id",
			Keys: bson.D{
				{types.BKBizIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "cluster_id",
			Keys: bson.D{
				{types.BKClusterIDFiled, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "namespace_id",
			Keys: bson.D{
				{types.BKNamespaceIDField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "reference_id_reference_kind",
			Keys: bson.D{
				{types.RefIDField, 1},
				{types.RefKindField, 1},
				{common.BkSupplierAccount, 1},
			},
			Background: true,
		},
	}

	return addIndexes(ctx, db, types.BKTableNameBaseContainer, indexes)
}

func upgradeContainerData(ctx context.Context, db dal.RDB) error {
	filter := mapstr.MapStr{
		common.BKAppIDField: mapstr.MapStr{common.BKDBExists: false},
	}

	for {
		containers := make([]types.Container, 0)
		err := db.Table(types.BKTableNameBaseContainer).Find(filter).Fields(types.BKPodIDField).
			Limit(common.BKMaxPageSize).All(ctx, &containers)
		if err != nil {
			blog.Errorf("get expired del archive data failed, err: %v", err)
			return err
		}

		if len(containers) == 0 {
			break
		}

		podIDs := make([]int64, len(containers))
		podIDMap := make(map[int64]struct{})
		for idx, container := range containers {
			podIDs[idx] = container.PodID
			podIDMap[container.PodID] = struct{}{}
		}

		podCond := mapstr.MapStr{
			types.BKIDField: mapstr.MapStr{common.BKDBIN: util.IntArrayUnique(podIDs)},
		}

		pods := make([]types.Pod, 0)
		err = db.Table(types.BKTableNameBasePod).Find(podCond).Fields(common.BKFieldID, common.BKAppIDField,
			types.BKClusterIDFiled, types.BKNamespaceIDField, types.RefIDField, types.RefKindField).
			Limit(common.BKMaxPageSize).All(ctx, &pods)
		if err != nil {
			blog.Errorf("get expired del archive data failed, err: %v", err)
			return err
		}

		wlPodIDMap := make(map[types.WorkloadType]map[int64][]int64)
		wlUpdateDataMap := make(map[types.WorkloadType]map[int64]mapstr.MapStr)

		// update pod related resource id info to container
		for _, pod := range pods {
			delete(podIDMap, pod.ID)

			updateData := mapstr.MapStr{
				common.BKAppIDField:      pod.BizID,
				types.BKClusterIDFiled:   pod.ClusterID,
				types.BKNamespaceIDField: pod.NamespaceID,
			}

			if pod.Ref != nil {
				updateData[types.RefKindField] = pod.Ref.Kind
				updateData[types.RefIDField] = pod.Ref.ID

				_, exists := wlPodIDMap[pod.Ref.Kind]
				if !exists {
					wlPodIDMap[pod.Ref.Kind] = make(map[int64][]int64)
					wlUpdateDataMap[pod.Ref.Kind] = make(map[int64]mapstr.MapStr)
				}

				wlPodIDMap[pod.Ref.Kind][pod.Ref.ID] = append(wlPodIDMap[pod.Ref.Kind][pod.Ref.ID], pod.ID)
				wlUpdateDataMap[pod.Ref.Kind][pod.Ref.ID] = updateData
				continue
			}

			updateCond := mapstr.MapStr{types.BKPodIDField: pod.ID}
			if err = db.Table(types.BKTableNameBaseContainer).Update(ctx, updateCond, updateData); err != nil {
				blog.Errorf("update container failed, err: %v, cond: %+v, data: %+v", err, updateCond, updateData)
				return err
			}
		}

		for wlKind, wlIDPodIDMap := range wlPodIDMap {
			for wlID, wlPodIDs := range wlIDPodIDMap {
				updateCond := mapstr.MapStr{
					types.BKPodIDField: mapstr.MapStr{common.BKDBIN: wlPodIDs},
				}
				updateData := wlUpdateDataMap[wlKind][wlID]

				if err = db.Table(types.BKTableNameBaseContainer).Update(ctx, updateCond, updateData); err != nil {
					blog.Errorf("update container failed, err: %v, cond: %+v, data: %+v", err, updateCond, updateData)
					return err
				}
			}
		}

		// delete containers that has no pod
		if len(podIDMap) > 0 {
			delPodIDs := make([]int64, 0)
			for podID := range podIDMap {
				delPodIDs = append(delPodIDs, podID)
			}

			delCond := mapstr.MapStr{
				types.BKPodIDField: mapstr.MapStr{common.BKDBIN: delPodIDs},
			}
			if err = db.Table(types.BKTableNameBaseContainer).Delete(ctx, delCond); err != nil {
				blog.Errorf("delete no pod container failed, err: %v, cond: %+v", err, delCond)
				return err
			}
		}

		if len(containers) < common.BKMaxPageSize {
			break
		}

		time.Sleep(time.Millisecond * 5)
	}

	return nil
}
