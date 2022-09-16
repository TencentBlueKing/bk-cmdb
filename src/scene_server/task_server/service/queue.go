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
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/task_server/logics"
	"configcenter/src/scene_server/task_server/taskconfig"
)

var (
	dbMaxRetry = 20
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

	locked, err := tq.lockTask(ctx, taskQueueInfo.TaskID, taskInfo.LockTTL)
	blog.Infof("start task %s", taskQueueInfo.TaskID)
	if err != nil {
		blog.Errorf("lock task failed, task name: %s, taskID: %s, err: %v", taskInfo.Name, taskQueueInfo.TaskID, err)
		time.Sleep(time.Second)
		return false
	}
	if !locked {
		return false
	}

	canExecute, err := tq.changeTaskToExecuting(ctx, taskQueueInfo.TaskID)
	blog.Infof("change task %s to executing, can execute %v", taskQueueInfo.TaskID, canExecute)
	if err != nil {
		if err := tq.unLockTask(ctx, taskQueueInfo.TaskID); err != nil {
			blog.Errorf("unlock failed, task type: %s, taskID: %s, err: %s", taskInfo.Name, taskQueueInfo.TaskID, err)
		}
		time.Sleep(time.Second)
		return false
	}
	if !canExecute {
		if err := tq.unLockTask(ctx, taskQueueInfo.TaskID); err != nil {
			blog.Errorf("unlock failed, task type: %s, taskID: %s, err: %s", taskInfo.Name, taskQueueInfo.TaskID, err)
		}
		return false
	}

	tq.executePush(ctx, taskInfo, &taskQueueInfo)
	return true
}

func (tq *TaskQueue) executePush(ctx context.Context, taskInfo TaskInfo, taskQueue *metadata.APITaskDetail) {
	blog.InfoJSON("start execute task, id: %s", taskQueue.TaskID)

	header := logics.GetDBHTTPHeader(taskQueue.Header)

	allSucc := true

	for _, subTask := range taskQueue.Detail {

		if subTask.Status == metadata.APITaskStatusSuccess {
			continue
		}

		if subTask.Status != metadata.APITaskStatusNew && subTask.Status != metadata.APITaskStatusWaitExecute {
			blog.Errorf("task status is not wait execute, task: %#v", taskQueue)
			allSucc = false
			break
		}

		var resp *metadata.Response
		var err error
		for retry := int64(0); retry < taskInfo.Retry; retry++ {
			if resp, err = tq.service.CoreAPI.TaskServer().Queue(taskInfo.Name).Post(ctx, header, taskInfo.Path,
				subTask.Data); err != nil {
				time.Sleep(time.Millisecond * 100)
				blog.Errorf("execute task http request failed, err: %v, taskID: %s, path: %s, header: %#v, data: %#v",
					err, taskQueue.TaskID, taskInfo.Path, header, subTask.Data)
				continue
			}
			break
		}

		if err != nil {
			ccErr := tq.service.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(taskQueue.Header)).
				CCError(common.CCErrCommHTTPDoRequestFailed)
			resp.Result = false
			resp.Code = ccErr.GetCode()
			resp.ErrMsg = ccErr.Error()
		}

		updateCond := mapstr.MapStr{"task_id": taskQueue.TaskID, "detail.sub_task_id": subTask.SubTaskID}
		updateData := mapstr.New()

		if err != nil || !resp.Result {
			allSucc = false
			updateData.Set("detail.$.status", metadata.APITAskStatusFail)
			updateData.Set("status", metadata.APITAskStatusFail)
		} else {
			updateData.Set("detail.$.status", metadata.APITaskStatusSuccess)
		}
		updateData.Set("detail.$.response", resp)
		updateData.Set(common.LastTimeField, time.Now())

		for dbRetry := 0; dbRetry < dbMaxRetry; dbRetry++ {
			dbErr := tq.service.DB.Table(common.BKTableNameAPITask).Update(ctx, updateCond, updateData)
			if dbErr != nil {
				blog.Errorf("update sub task resp failed, err: %v, cond: %#v, data: %#v", dbErr, updateCond, updateData)
				time.Sleep(time.Second * 3)
				continue
			}
			break
		}

	}

	// 所有任务执行完成，修改整个任务状态
	updateCond := mapstr.MapStr{common.BKTaskIDField: taskQueue.TaskID}
	updateData := mapstr.MapStr{common.LastTimeField: time.Now()}
	if allSucc {
		updateData.Set("status", metadata.APITaskStatusSuccess)
	} else {
		updateData.Set("status", metadata.APITAskStatusFail)
	}

	for dbRetry := 0; dbRetry < dbMaxRetry; dbRetry++ {
		dbErr := tq.service.DB.Table(common.BKTableNameAPITask).Update(ctx, updateCond, updateData)
		if dbErr != nil {
			blog.Errorf("update task status failed, err: %v, cond: %#v, data: %#v", dbErr, updateCond, updateData)
			time.Sleep(time.Second * 3)
			continue
		}

		dbErr = tq.service.DB.Table(common.BKTableNameAPITaskSyncHistory).Update(ctx, updateCond, updateData)
		if dbErr != nil {
			blog.Errorf("update history status failed, err: %v, cond: %#v, data: %#v", dbErr, updateCond, updateData)
			time.Sleep(time.Second * 3)
			continue
		}
		break
	}
}

func (tq *TaskQueue) lockTask(ctx context.Context, taskID string, ttl int64) (bool, error) {

	key := fmt.Sprintf("%s:apiTask:%s", common.BKCacheKeyV3Prefix, taskID)
	locked, err := tq.service.CacheDB.SetNX(ctx, key, time.Now(), time.Minute*time.Duration(ttl)).Result()
	if err != nil {
		blog.Errorf("lock task failed, err: %v, taskID: %s", err, taskID)
		return false, tq.service.CCErr.Error("zh-cn", common.CCErrTaskLockedTaskFail)
	}
	return locked, nil
}

func (tq *TaskQueue) unLockTask(ctx context.Context, taskID string) (err error) {

	key := fmt.Sprintf("%s:apiTask:%s", common.BKCacheKeyV3Prefix, taskID)
	_, err = tq.service.CacheDB.Del(ctx, key).Result()
	if err != nil {
		blog.Errorf("unlock task failed, err: %v, taskID: %s", err, taskID)
		return tq.service.CCErr.Error("zh-cn", common.CCErrTaskUnLockedTaskFail)
	}
	return nil
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
				continue
			}
			tq.compensateDBExecute(ctx)
		}
	}()
}

func (tq *TaskQueue) compensateDBExecute(ctx context.Context) {
	cond := condition.CreateCondition()
	cond.Field(common.BKStatusField).In([]metadata.APITaskStatus{metadata.APITaskStatusExecute})
	cond.Field(common.LastTimeField).Lt(time.Now().Add(-time.Minute * 20))
	data := mapstr.MapStr{
		"status":             metadata.APITaskStatusWaitExecute,
		common.LastTimeField: time.Now(),
	}
	err := tq.service.DB.Table(common.BKTableNameAPITask).Update(ctx, cond, data)
	if err != nil {
		blog.ErrorJSON("update task to wait execute error:%s, cond:%s", err.Error(), cond)
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
