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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (h *hostManager) CreateCloudSyncTask(ctx core.ContextParams, input *metadata.CloudTaskList) (uint64, error) {
	id, err := h.DbProxy.NextSequence(ctx, common.BKTableNameCloudTask)
	if err != nil {
		return 0, nil
	}

	input.TaskID = int64(id)
	if err := h.DbProxy.Table(common.BKTableNameCloudTask).Insert(ctx, input); err != nil {
		return 0, nil
	}

	return 0, nil
}

func (h *hostManager) CreateResourceConfirm(ctx core.ContextParams, input *metadata.ResourceConfirm) (uint64, error) {
	id, err := h.DbProxy.NextSequence(ctx, common.BKTableNameCloudResourceConfirm)
	if err != nil {
		return 0, err
	}

	input.ResourceID = int64(id)
	if err := h.DbProxy.Table(common.BKTableNameCloudResourceConfirm).Insert(ctx, input); err != nil {
		return 0, nil
	}

	return 0, nil
}

func (h *hostManager) CreateCloudSyncHistory(ctx core.ContextParams, input *metadata.CloudHistory) (uint64, error) {
	id, err := h.DbProxy.NextSequence(ctx, common.BKTableNameCloudSyncHistory)
	if err != nil {
		return 0, nil
	}

	input.HistoryID = int64(id)
	if err := h.DbProxy.Table(common.BKTableNameCloudSyncHistory).Insert(ctx, input); err != nil {
		return 0, nil
	}

	return id, nil
}

func (h *hostManager) CreateConfirmHistory(ctx core.ContextParams, input mapstr.MapStr) (uint64, error) {
	id, err := h.DbProxy.NextSequence(ctx, common.BKTableNameResourceConfirmHistory)
	if err != nil {
		return 0, nil
	}

	input[common.CloudSyncConfirmHistoryID] = id
	if err := h.DbProxy.Table(common.BKTableNameResourceConfirmHistory).Insert(ctx, input); err != nil {
		return 0, nil
	}

	return id, nil
}
