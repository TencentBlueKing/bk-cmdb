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

package service

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	taskUtil "configcenter/src/apimachinery/taskserver/util"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/task_server/logics"
	"configcenter/src/scene_server/task_server/taskconfig"
)

var (
	dbMaxRetry = 3
)

// TaskInfo TODO
type TaskInfo struct {
	Name    string
	Addr    func() ([]string, error)
	Path    string
	Retry   int64
	LockTTL int64
}

// TaskQueue TODO
type TaskQueue struct {
	task  []TaskInfo
	close bool
	sync.WaitGroup
	service *Service
}

// NewQueue TODO
func (s *Service) NewQueue(taskMap map[string]TaskInfo) *TaskQueue {
	var taskArr []TaskInfo
	codeTaskInfoMap := s.initCodeTaskConfig()
	if taskMap == nil {
		taskMap = make(map[string]TaskInfo)
	}
	for name, taskInfo := range codeTaskInfoMap {
		taskMap[name] = taskInfo
	}

	for _, taskItem := range taskMap {
		taskUtil.UpdateTaskServerConfigServ(taskItem.Name, taskItem.Addr)
		taskArr = append(taskArr, taskItem)
	}

	return &TaskQueue{
		task:    taskArr,
		service: s,
	}
}

// Stop TODO
func (tq *TaskQueue) Stop() {
	tq.close = true
	tq.Wait()
	return
}

// Start TODO
func (tq *TaskQueue) Start() {
	go tq.compensate(context.Background())

	for _, taskInfo := range tq.task {
		go func(taskInfo TaskInfo) {
			tq.Add(1)
			defer tq.Done()
			tq.execute(context.Background(), taskInfo)
		}(taskInfo)
	}
}

func (tq *TaskQueue) execute(ctx context.Context, task TaskInfo) {
	blog.Infof("start execute task queue: %s", task.Name)

	defer func() {
		if fetalErr := recover(); fetalErr != nil {
			blog.Errorf("err:%s, panic:%s", fetalErr, debug.Stack())
		}
	}()
	if task.Retry < 1 {
		task.Retry = 1
	}

	if task.LockTTL < 1 {
		task.LockTTL = 1
	}

	for {
		if tq.close {
			return
		}

		isMaster := tq.service.Engine.ServiceManageInterface.IsMaster()
		if !isMaster {
			blog.Infof("execute task, but is not master, skip")
			time.Sleep(time.Minute)
			continue
		}

		canSleep := true
		taskQueueInfoArr, err := tq.getWaitExecute(ctx, task.Name)
		if err != nil {
			blog.Errorf("get wait execute task failed, err: %v, task type: %s", err, task.Name)
			// select db failed, sleep 10s
			time.Sleep(time.Second * 10)
			continue
		}

		if len(taskQueueInfoArr) == 0 {
			// no task, sleep 5s
			time.Sleep(time.Second * 5)
			continue
		}

		for _, taskQueueInfo := range taskQueueInfoArr {
			if tq.close {
				return
			}

			isMaster := tq.service.Engine.ServiceManageInterface.IsMaster()
			if !isMaster {
				blog.Infof("execute task %s, but is not master, skip", taskQueueInfo.TaskID)
				time.Sleep(time.Minute)
				continue
			}

			execute := tq.executeTaskQueueItem(ctx, task, taskQueueInfo)
			if execute {
				canSleep = false
			}
		}

		if canSleep {
			time.Sleep(time.Second * 10)
		}
	}
}

// executeTaskQueueItem 执行异步任务，返回是否执行了任务
func (tq *TaskQueue) executeTaskQueueItem(ctx context.Context, taskInfo TaskInfo,
	taskQueueInfo metadata.APITaskDetail) bool {

	// set timeout and execute the task
	ctx, cancel := context.WithTimeout(ctx, time.Minute*time.Duration(taskInfo.LockTTL))
	defer cancel()

	blog.Infof("start task %s", taskQueueInfo.TaskID)

	// lock task to avoid conflict
	locked, err := tq.lockTask(ctx, taskQueueInfo.TaskID, taskInfo.LockTTL)
	if err != nil {
		blog.Errorf("lock task failed, task name: %s, taskID: %s, err: %v", taskInfo.Name, taskQueueInfo.TaskID, err)
		time.Sleep(time.Second)
		return false
	}
	if !locked {
		blog.Errorf("task(type %s, id: %s) is locked, return and retry later", taskInfo.Name, taskQueueInfo.TaskID)
		return false
	}

	defer func() {
		if err := tq.unLockTask(ctx, taskQueueInfo.TaskID); err != nil {
			blog.Errorf("unlock failed, task type: %s, taskID: %s, err: %v", taskInfo.Name, taskQueueInfo.TaskID, err)
		}
	}()

	canExecute, err := tq.changeTaskToExecuting(ctx, taskQueueInfo.TaskID)
	if err != nil {
		time.Sleep(time.Second)
		return false
	}

	blog.Infof("change task %s to executing, can execute %v", taskQueueInfo.TaskID, canExecute)
	if !canExecute {
		return false
	}

	tq.executePush(ctx, taskInfo, &taskQueueInfo)
	return true
}

func (tq *TaskQueue) executePush(ctx context.Context, taskInfo TaskInfo, taskQueue *metadata.APITaskDetail) {
	header := logics.GetDBHTTPHeader(taskQueue.Header)
	kit := rest.NewKitFromHeader(header, tq.service.CCErr)
	kit.Ctx = ctx

	blog.InfoJSON("start execute task, id: %s, rid: %s", taskQueue.TaskID, kit.Rid)

	allSucc := true

	for _, subTask := range taskQueue.Detail {
		success, needReturn := tq.executeSubTask(kit, taskInfo, taskQueue.TaskID, &subTask)
		if needReturn {
			return
		}

		if !success {
			allSucc = false
			break
		}
	}

	// 所有任务执行完成，修改整个任务状态
	blog.Infof("execute task %s done, all subtask success: %v, rid: %s", taskQueue.TaskID, allSucc, kit.Rid)

	updateCond := mapstr.MapStr{common.BKTaskIDField: taskQueue.TaskID}
	var updateStatus metadata.APITaskStatus
	if allSucc {
		updateStatus = metadata.APITaskStatusSuccess
	} else {
		updateStatus = metadata.APITAskStatusFail
	}

	needReturn := retryWrapper(kit, dbMaxRetry, func() error {
		if err := tq.service.Logics.UpdateTaskStatus(kit.Ctx, updateCond, updateStatus, kit.Rid); err != nil {
			time.Sleep(time.Second * 3)
			return err
		}
		return nil
	})

	if needReturn {
		return
	}

	blog.Infof("successfully updated task %s status to %s, rid: %s", taskQueue.TaskID, updateStatus, kit.Rid)
}

// executeSubTask execute subtask, returns if it is successful and if the task needs return
func (tq *TaskQueue) executeSubTask(kit *rest.Kit, taskInfo TaskInfo, taskID string,
	subTask *metadata.APISubTaskDetail) (bool, bool) {

	blog.Infof("start execute task(id: %s) subtask(id: %s)", taskID, subTask.SubTaskID)

	if subTask.Status == metadata.APITaskStatusSuccess {
		return true, false
	}

	if subTask.Status != metadata.APITaskStatusNew && subTask.Status != metadata.APITaskStatusWaitExecute {
		blog.Errorf("task(id: %s) subtask(id: %s) status is wrong", taskID, subTask.SubTaskID)
		return false, false
	}

	var resp *metadata.Response
	var err error
	needReturn := retryWrapper(kit, int(taskInfo.Retry), func() error {
		if resp, err = tq.service.CoreAPI.TaskServer().Queue(taskInfo.Name).Post(kit.Ctx, kit.Header, taskInfo.Path,
			subTask.Data); err != nil {
			time.Sleep(time.Millisecond * 100)
			blog.Errorf("execute task http request failed, err: %v, taskID: %s, path: %s, data: %#v, rid: %s",
				err, taskID, taskInfo.Path, subTask.Data, kit.Rid)
			return err
		}
		return nil
	})

	if needReturn {
		return false, true
	}

	if err != nil {
		resp.Result = false
		resp.Code = common.CCErrCommHTTPDoRequestFailed
		resp.ErrMsg = kit.CCError.CCErrorf(common.CCErrCommHTTPDoRequestFailed).Error()
	}

	updateCond := mapstr.MapStr{"task_id": taskID, "detail.sub_task_id": subTask.SubTaskID}
	updateData := mapstr.New()

	if err != nil || !resp.Result {
		updateData.Set("detail.$.status", metadata.APITAskStatusFail)
		updateData.Set("status", metadata.APITAskStatusFail)
	} else {
		updateData.Set("detail.$.status", metadata.APITaskStatusSuccess)
	}
	updateData.Set("detail.$.response", resp)
	updateData.Set(common.LastTimeField, time.Now())

	needReturn = retryWrapper(kit, dbMaxRetry, func() error {
		err := tq.service.DB.Table(common.BKTableNameAPITask).Update(kit.Ctx, updateCond, updateData)
		if err != nil {
			blog.Errorf("update sub task resp failed, err: %v, cond: %#v, data: %#v", err, updateCond, updateData)
			time.Sleep(time.Second * 3)
			return err
		}
		return nil
	})

	if needReturn {
		return false, true
	}

	blog.Infof("finished executing task(id: %s) subtask(id: %s)", taskID, subTask.SubTaskID)

	// the subtask is not successful, returns the status after updating its response to end the task
	if err != nil || !resp.Result {
		return false, false
	}

	return true, false
}

// retryWrapper retry task execute step wrapper, returns if task is terminated.
func retryWrapper(kit *rest.Kit, maxRetry int, handler func() error) bool {
	for retry := 0; retry < maxRetry; retry++ {
		select {
		case <-kit.Ctx.Done():
			blog.Errorf("context canceled when executing task, rid: %s", kit.Rid)
			return true
		default:
			err := handler()
			if err == nil {
				return false
			}
		}
	}

	return false
}

func (tq *TaskQueue) taskLockKey(taskID string) string {
	return fmt.Sprintf("%s:apiTask:%s", common.BKCacheKeyV3Prefix, taskID)
}

func (tq *TaskQueue) lockTask(ctx context.Context, taskID string, ttl int64) (bool, error) {

	key := tq.taskLockKey(taskID)
	locked, err := tq.service.CacheDB.SetNX(ctx, key, time.Now(), time.Minute*time.Duration(ttl)).Result()
	if err != nil {
		blog.Errorf("lock task failed, err: %v, taskID: %s", err, taskID)
		return false, tq.service.CCErr.Error("zh-cn", common.CCErrTaskLockedTaskFail)
	}
	return locked, nil
}

func (tq *TaskQueue) unLockTask(ctx context.Context, taskID string) (err error) {

	key := tq.taskLockKey(taskID)
	_, err = tq.service.CacheDB.Del(ctx, key).Result()
	if err != nil {
		blog.Errorf("unlock task failed, err: %v, taskID: %s", err, taskID)
		return tq.service.CCErr.Error("zh-cn", common.CCErrTaskUnLockedTaskFail)
	}
	return nil
}

// isTaskLocked returns if task is locked.
func (tq *TaskQueue) isTaskLocked(ctx context.Context, taskID string) (bool, error) {
	key := tq.taskLockKey(taskID)
	result, err := tq.service.CacheDB.Exists(ctx, key).Result()
	if err != nil {
		blog.Errorf("check if task %s is locked failed, err: %v", taskID, err)
		return false, err
	}
	return result == 1, nil
}

func (tq *TaskQueue) getWaitExecute(ctx context.Context, name string) ([]metadata.APITaskDetail, error) {
	cond := mapstr.MapStr{
		common.BKTaskTypeField: name,
		common.BKStatusField: mapstr.MapStr{
			common.BKDBIN: []metadata.APITaskStatus{metadata.APITaskStatusNew, metadata.APITaskStatusWaitExecute},
		},
	}

	rows := make([]metadata.APITaskDetail, 0)
	err := tq.service.DB.Table(common.BKTableNameAPITask).Find(cond).Sort("create_time").Limit(20).All(ctx, &rows)
	if err != nil {
		blog.ErrorJSON("query wait execute failed, err: %v, task type: %s, cond: %#v", err, name, cond)
		return nil, tq.service.CCErr.Error("zh-cn", common.CCErrCommDBSelectFailed)
	}

	return rows, nil
}

func (tq *TaskQueue) changeTaskToExecuting(ctx context.Context, taskID string) (bool, error) {
	cond := mapstr.MapStr{
		common.BKTaskIDField: taskID,
		common.BKStatusField: mapstr.MapStr{
			common.BKDBIN: []metadata.APITaskStatus{metadata.APITaskStatusNew, metadata.APITaskStatusWaitExecute},
		},
	}

	cnt, err := tq.service.DB.Table(common.BKTableNameAPITask).Find(cond).Count(ctx)
	if err != nil {
		blog.Errorf("query wait execute error, taskID: %s, err: %v, cond: %#v", taskID, err, cond)
		return false, tq.service.CCErr.Error("zh-cn", common.CCErrCommDBSelectFailed)
	}
	if cnt != 1 {
		blog.Errorf("query wait execute, task not equal 1, taskID: %s, cond: %#v, cnt: %d", taskID, cond, cnt)
		return false, nil
	}

	data := mapstr.MapStr{
		common.BKStatusField: metadata.APITaskStatusExecute,
		common.LastTimeField: time.Now(),
	}
	err = tq.service.DB.Table(common.BKTableNameAPITask).Update(ctx, cond, data)
	if err != nil {
		blog.Errorf("update task to execute failed, err: %v, task: %s, cond: %#v", err, taskID, cond)
		return false, tq.service.CCErr.Error("zh-cn", common.CCErrCommDBUpdateFailed)
	}

	err = tq.service.DB.Table(common.BKTableNameAPITaskSyncHistory).Update(ctx, cond, data)
	if err != nil {
		blog.Errorf("update task sync history to execute failed, err: %v, task: %s, cond: %#v", err, taskID, cond)
		return false, tq.service.CCErr.Error("zh-cn", common.CCErrCommDBUpdateFailed)
	}
	return true, nil
}

func (tq *TaskQueue) compensate(ctx context.Context) {
	go func() {
		timer := time.NewTicker(time.Minute * 10)
		for range timer.C {
			if tq.close {
				return
			}

			isMaster := tq.service.Engine.ServiceManageInterface.IsMaster()
			if !isMaster {
				blog.Infof("compensate task, but is not master, skip")
				continue
			}
			tq.compensateDBExecute(ctx)
		}
	}()
}

// compensateDBExecute compensate old executing task status, if it is not finished, put back to wait execute queue.
func (tq *TaskQueue) compensateDBExecute(ctx context.Context) {
	cond := mapstr.MapStr{
		common.BKStatusField: mapstr.MapStr{
			common.BKDBIN: []metadata.APITaskStatus{metadata.APITaskStatusExecute},
		},
		common.LastTimeField: mapstr.MapStr{
			common.BKDBLT: time.Now().Add(-time.Minute * 20),
		},
	}

	for {
		rid := util.GenerateRID()
		blog.Infof("start compensate task, rid: %s", rid)

		rows := make([]metadata.APITaskDetail, 0)
		err := tq.service.DB.Table(common.BKTableNameAPITask).Find(cond).Sort("create_time").Limit(200).All(ctx, &rows)
		if err != nil {
			blog.ErrorJSON("get task to compensate failed, cond: %#v, err: %v, rid: %s", cond, err, rid)
			return
		}

		if len(rows) == 0 {
			return
		}

		// get all interrupted tasks that are not currently running, compensate their status
		interruptedTaskIDs := make([]metadata.APITaskDetail, 0)
		for _, row := range rows {
			locked, err := tq.isTaskLocked(ctx, row.TaskID)
			if err != nil {
				blog.Errorf("compensate duplicate tasks failed, tasks: %#v, err: %v, rid: %s", err, rows, rid)
				return
			}

			if locked {
				blog.Infof("task %s is locked, skip, rid: %s", row.TaskID, rid)
				continue
			}
			interruptedTaskIDs = append(interruptedTaskIDs, row)
		}

		blog.Infof("start compensate tasks(%+v) status, rid: %s", interruptedTaskIDs, rid)

		unfinishedTaskIDs, err := tq.service.Logics.CompensateStatus(ctx, interruptedTaskIDs, rid)
		if err != nil {
			blog.Errorf("compensate duplicate tasks failed, tasks: %#v, err: %v, rid: %s", err, rows, rid)
			return
		}

		if len(unfinishedTaskIDs) == 0 {
			continue
		}

		blog.Infof("compensate tasks(%+v) status success, remaining: %+v, rid: %s", interruptedTaskIDs,
			unfinishedTaskIDs, rid)

		// put these not finished task into wait execute queue
		updateCond := mapstr.MapStr{
			common.BKTaskIDField: mapstr.MapStr{
				common.BKDBIN: unfinishedTaskIDs,
			},
		}

		data := mapstr.MapStr{
			common.BKStatusField: metadata.APITaskStatusWaitExecute,
			common.LastTimeField: time.Now(),
		}
		err = tq.service.DB.Table(common.BKTableNameAPITask).Update(ctx, updateCond, data)
		if err != nil {
			blog.Errorf("update task to wait execute failed, err: %v, cond: %+v, rid: %s", err, cond, rid)
			return
		}

		blog.Infof("compensate tasks(%+v) to execute success, rid: %s", unfinishedTaskIDs, rid)
	}
}

func (s *Service) initCodeTaskConfig() map[string]TaskInfo {
	taskInfoMap := make(map[string]TaskInfo, 0)
	codeTaskConfigArr := taskconfig.GetCodeTaskConfig()

	for _, codeTaskConfig := range codeTaskConfigArr {
		ti := TaskInfo{
			Name:    codeTaskConfig.Name,
			Retry:   codeTaskConfig.Retry,
			Path:    codeTaskConfig.Path,
			LockTTL: codeTaskConfig.LockTTL,
		}
		switch codeTaskConfig.SvrType {
		case types.CC_MODULE_APISERVER:
			ti.Addr = s.Engine.Discovery().ApiServer().GetServers
		case types.CC_MODULE_HOST:
			ti.Addr = s.Engine.Discovery().HostServer().GetServers
		case types.CC_MODULE_PROC:
			ti.Addr = s.Engine.Discovery().ProcServer().GetServers
		case types.CC_MODULE_TOPO:
			ti.Addr = s.Engine.Discovery().TopoServer().GetServers
		case types.CC_MODULE_TASK:
			ti.Addr = s.Engine.Discovery().TaskServer().GetServers
		default:
			panicErr := fmt.Sprintf("init task queue(%s), but svrType(%s) is invalid", ti.Name, codeTaskConfig.SvrType)
			panic(panicErr)
		}
		taskInfoMap[ti.Name] = ti
	}

	return taskInfoMap
}
