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

package y3_9_202110211014

import (
	"context"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

// apiTask api task info
type apiTaskInfo struct {
	TaskID string        `bson:"task_id"`
	Status int64         `bson:"status"`
	Detail []subTaskInfo `bson:"detail"`
}

// subTask sub task data and execute detail
type subTaskInfo struct {
	Status int64 `bson:"status"`
}

// migrateApiTask migrate api task table to common api task sync status table
func migrateApiTask(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	indexes := map[string]types.Index{
		"idx_taskID": {Name: "idx_taskID", Keys: map[string]int32{"task_id": 1}, Background: true},
		"idx_flag_status_createTime": {Name: "idx_flag_status_createTime", Keys: map[string]int32{"flag": 1,
			"status": 1, "create_time": 1}, Background: true},
		"idx_lastTime_status": {Name: "idx_lastTime_status", Keys: map[string]int32{"last_time": 1, "status": 1},
			Background: true},
		"idx_taskType_instID": {Name: "idx_taskType_instID", Keys: map[string]int32{"task_type": 1, "bk_inst_id": 1},
			Background: true, Unique: true, PartialFilterExpression: map[string]interface{}{"status":
			map[string]interface{}{common.BKDBIN: []string{"0", "1", "2"}}}},
	}

	if err := reconcileIndexes(ctx, db, "cc_APITask", indexes); err != nil {
		blog.Errorf("reconcile cc_APITask table indexes failed, err: %v", err)
		return err
	}

	toRemoveFields := []string{"name"}
	taskFilter := map[string]interface{}{
		"name": map[string]interface{}{common.BKDBExists: true},
	}

	// previous int64 status to current enum status mapping
	taskStatusMapping := map[int64]string{
		0:   "new",
		1:   "waiting",
		100: "executing",
		200: "finished",
		500: "failure",
	}

	for {
		taskArr := make([]apiTaskInfo, 0)
		err := db.Table("cc_APITask").Find(taskFilter).Start(0).Limit(common.BKMaxPageSize).
			Fields("task_id", "status", "detail.status").All(ctx, &taskArr)
		if err != nil {
			blog.Errorf("get task ids to update field failed, err: %v", err)
			return err
		}

		if len(taskArr) == 0 {
			break
		}

		taskIDs := make([]string, len(taskArr))
		for index, task := range taskArr {
			taskIDs[index] = task.TaskID

			updateStatusFilter := map[string]interface{}{
				"task_id": task.TaskID,
			}

			updateStatusData := map[string]interface{}{
				"status": taskStatusMapping[task.Status],
			}
			for idx, subTask := range task.Detail {
				updateStatusData["detail."+strconv.FormatInt(int64(idx), 10)] = taskStatusMapping[subTask.Status]
			}
			if err := db.Table("cc_APITaskSyncHistory").Update(ctx, updateStatusFilter, updateStatusData); err != nil {
				blog.Errorf("update api task status type failed, filter: %#v, err: %v", updateStatusFilter, err)
				return err
			}
		}

		filter := map[string]interface{}{
			"task_id": map[string]interface{}{
				common.BKDBIN: taskIDs,
			},
		}

		if err := db.Table("cc_APITask").DropColumns(ctx, filter, toRemoveFields); err != nil {
			blog.Errorf("drop cc_APITask columns failed, filter: %#v, err: %v", filter, err)
			return err
		}

		if err := db.Table("cc_APITask").RenameColumn(ctx, filter, "flag", "task_type"); err != nil {
			blog.Errorf("rename cc_APITask bk_set_id column failed, filter: %#v, err: %v", filter, err)
			return err
		}

		if len(taskArr) < common.BKMaxPageSize {
			break
		}
	}
	return nil
}

// migrateAPITaskSyncStatus migrate set template sync status table to common api task sync status table
func migrateAPITaskSyncStatus(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	if err := db.DropTable(ctx, "cc_SetTemplateSyncStatus"); err != nil {
		blog.Errorf("drop cc_SetTemplateSyncStatus table failed, err: %v", err)
		return err
	}

	newTableExists, err := db.HasTable(ctx, "cc_APITaskSyncHistory")
	if err != nil {
		blog.Errorf("check if table cc_APITaskSyncHistory exists failed, err: %v", err)
		return err
	}

	oldTableExists, err := db.HasTable(ctx, "cc_SetTemplateSyncHistory")
	if err != nil {
		blog.Errorf("check if table cc_SetTemplateSyncHistory exists failed, err: %v", err)
		return err
	}

	if !newTableExists && oldTableExists {
		if err := db.RenameTable(ctx, "cc_SetTemplateSyncHistory", "cc_APITaskSyncHistory"); err != nil &&
			!strings.Contains(err.Error(), "NamespaceExists") {
			blog.Errorf("rename cc_SetTemplateSyncHistory table to cc_APITaskSyncHistory failed, err: %v", err)
			return err
		}
	}

	indexes := map[string]types.Index{
		"idx_instID_flag_createTime": {Name: "idx_instID_flag_createTime", Keys: map[string]int32{"bk_inst_id": 1,
			"flag": 1, "create_time": 1}, Background: true},
	}

	if err := reconcileIndexes(ctx, db, "cc_APITaskSyncHistory", indexes); err != nil {
		blog.Errorf("reconcile cc_APITaskSyncHistory table indexes failed, err: %v", err)
		return err
	}

	toRemoveFields := []string{"bk_biz_id", "bk_set_name", "set_template_id"}
	toAddField := map[string]interface{}{
		"task_type": "set_template_sync",
	}

	historyFilter := map[string]interface{}{
		"bk_set_id": map[string]interface{}{common.BKDBExists: true},
	}

	for {
		historyArr := make([]map[string]interface{}, 0)
		err := db.Table("cc_APITaskSyncHistory").Find(historyFilter).Start(0).Limit(common.BKMaxPageSize).
			Fields("task_id").All(ctx, &historyArr)
		if err != nil {
			blog.Errorf("get task ids to update field failed, err: %v", err)
			return err
		}

		if len(historyArr) == 0 {
			break
		}

		taskIDs := make([]interface{}, len(historyArr))
		for index, history := range historyArr {
			taskIDs[index] = history["task_id"]
		}

		filter := map[string]interface{}{
			"task_id": map[string]interface{}{
				common.BKDBIN: taskIDs,
			},
		}

		if err := db.Table("cc_APITaskSyncHistory").DropColumns(ctx, filter, toRemoveFields); err != nil {
			blog.Errorf("drop cc_APITaskSyncHistory columns failed, filter: %#v, err: %v", filter, err)
			return err
		}

		if err := db.Table("cc_APITaskSyncHistory").RenameColumn(ctx, filter, "bk_set_id", "bk_inst_id"); err != nil {
			blog.Errorf("rename cc_APITaskSyncHistory bk_set_id column failed, filter: %#v, err: %v", filter, err)
			return err
		}

		if err := db.Table("cc_APITaskSyncHistory").Update(ctx, filter, toAddField); err != nil {
			blog.Errorf("add task_type to cc_APITaskSyncHistory failed, filter: %#v, err: %v", filter, err)
			return err
		}

		if len(historyArr) < common.BKMaxPageSize {
			break
		}
	}
	return nil
}

// reconcileIndexes update table indexes
func reconcileIndexes(ctx context.Context, db dal.RDB, table string, indexes map[string]types.Index) error {
	existIndexArr, err := db.Table(table).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for audit table failed, err: %s", err.Error())
		return err
	}

	existIdxMap := make(map[string]struct{})
	removeIndexes := make([]string, 0)
	for _, index := range existIndexArr {
		if _, exists := indexes[index.Name]; exists || index.Name == "_id_" {
			existIdxMap[index.Name] = struct{}{}
			continue
		}
		removeIndexes = append(removeIndexes, index.Name)
	}

	for _, removeIndex := range removeIndexes {
		if err = db.Table(table).DropIndex(ctx, removeIndex); err != nil {
			blog.Errorf("remove %s index failed, err: %v, index: %s", table, err, removeIndex)
			return err
		}
	}

	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}

		if err = db.Table(table).CreateIndex(ctx, index); err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create %s index failed, err: %v, index: %#v", table, err, index)
			return err
		}
	}
	return nil
}
