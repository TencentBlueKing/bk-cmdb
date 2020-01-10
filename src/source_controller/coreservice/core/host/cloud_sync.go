/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package host

import (
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (hm *hostManager) CreateCloudSyncTask(kit *rest.Kit, input *metadata.CloudTaskList) (uint64, error) {
	id, err := hm.DbProxy.NextSequence(kit.Ctx, common.BKTableNameCloudTask)
	if err != nil {
		return 0, err
	}

	input.TaskID = int64(id)
	if err := hm.DbProxy.Table(common.BKTableNameCloudTask).Insert(kit.Ctx, input); err != nil {
		return 0, err
	}

	return id, nil
}

func (hm *hostManager) CreateResourceConfirm(kit *rest.Kit, input *metadata.ResourceConfirm) (uint64, error) {
	id, err := hm.DbProxy.NextSequence(kit.Ctx, common.BKTableNameCloudResourceConfirm)
	if err != nil {
		return 0, err
	}

	input.ResourceID = int64(id)
	if err := hm.DbProxy.Table(common.BKTableNameCloudResourceConfirm).Insert(kit.Ctx, input); err != nil {
		return 0, err
	}

	return id, nil
}

func (hm *hostManager) CreateCloudSyncHistory(kit *rest.Kit, input *metadata.CloudHistory) (uint64, error) {
	id, err := hm.DbProxy.NextSequence(kit.Ctx, common.BKTableNameCloudSyncHistory)
	if err != nil {
		return 0, err
	}

	input.HistoryID = int64(id)
	if err := hm.DbProxy.Table(common.BKTableNameCloudSyncHistory).Insert(kit.Ctx, input); err != nil {
		return 0, err
	}

	return id, nil
}

func (hm *hostManager) CreateConfirmHistory(kit *rest.Kit, input mapstr.MapStr) (uint64, error) {
	id, err := hm.DbProxy.NextSequence(kit.Ctx, common.BKTableNameResourceConfirmHistory)
	if err != nil {
		return 0, err
	}

	input[common.CloudSyncConfirmHistoryID] = id
	if err := hm.DbProxy.Table(common.BKTableNameResourceConfirmHistory).Insert(kit.Ctx, input); err != nil {
		return 0, err
	}

	return id, nil
}
