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

package cloud

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"time"
)

func (c *cloudOperation) CreateSyncTask(kit *rest.Kit, task *metadata.CloudSyncTask) (*metadata.CloudSyncTask, errors.CCErrorCoder) {
	if err := c.validCreateSycTask(kit, task); nil != err {
		blog.Errorf("CreateAccount failed, valid error: %+v, rid: %s", err, kit.Rid)
		return nil, err
	}

	id, err := c.dbProxy.NextSequence(kit.Ctx, common.BKTableNameCloudSyncTask)
	if nil != err {
		blog.Errorf("CreateAccount failed, generate id failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}
	task.TaskID = int64(id)
	ts := time.Now()
	task.OwnerID = kit.SupplierAccount
	task.LastEditor = task.Creator
	task.CreateTime = ts
	task.LastTime = ts

	err = c.dbProxy.Table(common.BKTableNameCloudAccount).Insert(kit.Ctx, task)
	if err != nil {
		blog.ErrorJSON("CreateSyncTask failed, db insert failed, taskName: %s, err: %v, rid: %s", task.TaskName, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
	}

	return task, nil
}

func (c *cloudOperation) SearchSyncTask(kit *rest.Kit, option *metadata.SearchCloudOption) (*metadata.MultipleCloudSyncTask, errors.CCErrorCoder) {
	results := make([]metadata.CloudSyncTask, 0)
	err := c.dbProxy.Table(common.BKTableNameCloudSyncTask).Find(option.Condition).Fields(option.Fields...).
		Start(uint64(option.Page.Start)).Limit(uint64(option.Page.Limit)).Sort(option.Page.Sort).All(kit.Ctx, &results)
	if err != nil {
		blog.ErrorJSON("SearchSyncTask failed, db find failed, option: %v, err: %v, rid: %s", option, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return &metadata.MultipleCloudSyncTask{Count: int64(len(results)), Info: results}, nil
}

func (c *cloudOperation) UpdateSyncTask(kit *rest.Kit, taskID int64, option mapstr.MapStr) errors.CCErrorCoder {
	return nil
}

func (c *cloudOperation) DeleteSyncTask(kit *rest.Kit, taskID int64) errors.CCErrorCoder {
	cond := map[string]int64{common.BKCloudSyncTaskID: taskID}
	if e := c.dbProxy.Table(common.BKTableNameCloudSyncTask).Delete(kit.Ctx, cond); e != nil {
		blog.Errorf("DeleteSyncTask failed, mongodb operate failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameCloudAccount, cond, e, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBDeleteFailed)
	}
	return nil
}
