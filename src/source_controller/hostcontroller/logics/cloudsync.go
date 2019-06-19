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

package logics

import (
	"context"

	"configcenter/src/common"
)

func (lgc *Logics) CreateCloudTask(ctx context.Context, input interface{}) error {
	objID, err := lgc.Instance.NextSequence(ctx, common.BKTableNameCloudTask)
	if err != nil {
		return err
	}

	inputc := input.(map[string]interface{})
	inputc["bk_task_id"] = objID
	if err := lgc.Instance.Table(common.BKTableNameCloudTask).Insert(ctx, inputc); err != nil {
		return err
	}

	return nil
}

func (lgc *Logics) CreateResourceConfirm(ctx context.Context, input interface{}) error {
	objID, err := lgc.Instance.NextSequence(ctx, common.BKTableNameCloudResourceConfirm)
	if err != nil {
		return err
	}

	inputc := input.(map[string]interface{})
	inputc["bk_resource_id"] = objID
	if err := lgc.Instance.Table(common.BKTableNameCloudResourceConfirm).Insert(ctx, inputc); err != nil {
		return err
	}

	return nil
}

func (lgc *Logics) CreateCloudHistory(ctx context.Context, input interface{}) error {
	objID, err := lgc.Instance.NextSequence(ctx, common.BKTableNameCloudSyncHistory)
	if err != nil {
		return err
	}

	inputc := input.(map[string]interface{})
	inputc["bk_history_id"] = objID
	if err := lgc.Instance.Table(common.BKTableNameCloudSyncHistory).Insert(ctx, inputc); err != nil {
		return err
	}

	return nil
}

func (lgc *Logics) CreateConfirmHistory(ctx context.Context, input interface{}) error {
	objID, err := lgc.Instance.NextSequence(ctx, common.BKTableNameResourceConfirmHistory)
	if err != nil {
		return err
	}

	inputc := input.(map[string]interface{})
	inputc["confirm_history_id"] = objID
	if err := lgc.Instance.Table(common.BKTableNameResourceConfirmHistory).Insert(ctx, inputc); err != nil {
		return err
	}

	return nil
}
