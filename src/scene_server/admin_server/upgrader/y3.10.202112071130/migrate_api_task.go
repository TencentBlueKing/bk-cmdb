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

package y3_10_202112071130

import (
	"context"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

// apiTask api task info
type apiTaskInfo struct {
	TaskID interface{}   `bson:"task_id"`
	Status interface{}   `bson:"status"`
	Detail []subTaskInfo `bson:"detail"`
}

// subTask sub task data and execute detail
type subTaskInfo struct {
	Status interface{} `bson:"status"`
}

// migrateApiTask migrate api task table to common api task sync status table
func migrateApiTask(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	indexes := map[string]types.Index{
		"idx_taskID": {Name: "idx_taskID", Keys: bson.D{{"task_id", 1}}, Background: true},
		"idx_flag_status_createTime": {Name: "idx_flag_status_createTime", Keys: bson.D{{"flag", 1},
			{"status", 1}, {"create_time", 1}}, Background: true},
		"idx_lastTime_status": {Name: "idx_lastTime_status", Keys: bson.D{{"last_time", 1},
			{"status", 1}}, Background: true},
		"idx_lastTime": {Name: "idx_lastTime", Keys: bson.D{{"last_time", 1}}},
	}

	if err := reconcileIndexes(ctx, db, "cc_APITask", indexes); err != nil {
		blog.Errorf("reconcile cc_APITask table indexes failed, err: %v", err)
		return err
	}

	toRemoveFields := []string{"name"}
	taskFilter := map[string]interface{}{
		common.BKDBOR: []map[string]interface{}{
			{"name": map[string]interface{}{common.BKDBExists: true}},
			{"status": map[string]interface{}{common.BKDBType: "long"}},
			{"detail.status": map[string]interface{}{common.BKDBType: "long"}},
		},
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

		taskIDs := make([]interface{}, len(taskArr))
		for index, task := range taskArr {
			taskIDs[index] = task.TaskID

			if err := updateApiTaskStatus(ctx, db, task); err != nil {
				blog.Errorf("update task status failed, err: %v, task: %#v", err, task)
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

// taskStatusMapping previous int64 status to current enum status mapping
var taskStatusMapping = map[int64]string{
	0:   "new",
	1:   "waiting",
	100: "executing",
	200: "finished",
	500: "failure",
}

// updateApiTaskStatus update previous api task from int64 to current enum, skip those that are already string
func updateApiTaskStatus(ctx context.Context, db dal.RDB, task apiTaskInfo) error {
	updateStatusData := make(map[string]interface{})

	switch task.Status.(type) {
	case string:
	default:
		status, err := util.GetInt64ByInterface(task.Status)
		if err != nil {
			blog.Errorf("parse task status(%#v) failed, err: %v", task.Status, err)
			return err
		}

		updateStatusData["status"] = taskStatusMapping[status]
	}

	for idx, subTask := range task.Detail {
		switch subTask.Status.(type) {
		case string:
		default:
			status, err := util.GetInt64ByInterface(subTask.Status)
			if err != nil {
				blog.Errorf("parse sub task status(%#v) failed, err: %v", subTask.Status, err)
				return err
			}

			key := "detail." + strconv.FormatInt(int64(idx), 10) + ".status"
			updateStatusData[key] = taskStatusMapping[status]
		}
	}

	if len(updateStatusData) == 0 {
		return nil
	}

	updateStatusFilter := map[string]interface{}{
		"task_id": task.TaskID,
	}

	if err := db.Table("cc_APITask").Update(ctx, updateStatusFilter, updateStatusData); err != nil {
		blog.Errorf("update api task status type failed, filter: %#v, err: %v", updateStatusFilter, err)
		return err
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
		"idx_instID_flag_createTime": {Name: "idx_instID_flag_createTime", Keys: bson.D{{"bk_inst_id", 1},
			{"flag", 1}, {"create_time", 1}}, Background: true},
		"idx_lastTime": {Name: "idx_lastTime", Keys: bson.D{{"last_time", 1}},
			ExpireAfterSeconds: 3 * 30 * 24 * 60 * 60},
	}

	if err := reconcileIndexes(ctx, db, "cc_APITaskSyncHistory", indexes); err != nil {
		blog.Errorf("reconcile cc_APITaskSyncHistory table indexes failed, err: %v", err)
		return err
	}

	toRemoveFields := []string{"bk_biz_id", "bk_set_name", "set_template_id"}
	toAddField := map[string]interface{}{
		"task_type": "set_template_sync",
	}

	updateStatusField := map[string]interface{}{
		"status": "executing",
	}

	historyFilter := map[string]interface{}{
		"bk_set_id": map[string]interface{}{common.BKDBExists: true},
	}

	for {
		historyArr := make([]map[string]interface{}, 0)
		err := db.Table("cc_APITaskSyncHistory").Find(historyFilter).Start(0).Limit(common.BKMaxPageSize).
			Fields("task_id", "status").All(ctx, &historyArr)
		if err != nil {
			blog.Errorf("get task ids to update field failed, err: %v", err)
			return err
		}

		if len(historyArr) == 0 {
			break
		}

		taskIDs := make([]interface{}, len(historyArr))
		updateStatusTaskIDs := make([]interface{}, 0)
		for index, history := range historyArr {
			taskIDs[index] = history["task_id"]
			if history["status"] == "syncing" {
				updateStatusTaskIDs = append(updateStatusTaskIDs, history["task_id"])
			}
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

		if len(updateStatusTaskIDs) > 0 {
			filter := map[string]interface{}{
				"task_id": map[string]interface{}{
					common.BKDBIN: updateStatusTaskIDs,
				},
			}

			if err := db.Table("cc_APITaskSyncHistory").Update(ctx, filter, updateStatusField); err != nil {
				blog.Errorf("update status of cc_APITaskSyncHistory failed, filter: %#v, err: %v", filter, err)
				return err
			}
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
