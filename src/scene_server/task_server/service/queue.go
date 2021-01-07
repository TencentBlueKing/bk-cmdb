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

type TaskInfo struct {
	Name  string
	Addr  func() ([]string, error)
	Path  string
	Retry int64
}

type TaskQueue struct {
	task  []TaskInfo
	close bool
	sync.WaitGroup
	service *Service
}

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

func (tq *TaskQueue) Stop() {
	tq.close = true
	tq.Wait()
	return
}

func (tq *TaskQueue) Start() {

	go tq.compensate(context.Background())
	for _, taskInfo := range tq.task {

		go func(taskInfo TaskInfo) {
			tq.Add(1)
			defer tq.Done()
			tq.executeWrap(context.Background(), taskInfo)
		}(taskInfo)
	}
}

func (tq *TaskQueue) executeWrap(ctx context.Context, taskInfo TaskInfo) {
	for {
		if tq.close {
			return
		}
		tq.execute(ctx, taskInfo)
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

	for {
		if tq.close {
			return
		}
		canSleep := true
		taskQueueInfoArr, err := tq.getWaitExectue(ctx, task.Name)
		if err != nil {
			blog.Errorf("execute get wait execute task error. task name:%s, err:%s", task.Name, err.Error())
			// select db error. sleep 10s
			time.Sleep(time.Second * 10)
			continue
		}
		if len(taskQueueInfoArr) == 0 {
			// not task. sleep 5s
			time.Sleep(time.Second * 5)
			continue
		}
		for _, taskQueueInfo := range taskQueueInfoArr {
			if tq.close {
				return
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

// executeTaskQueueItem  返回是否执行任务
func (tq *TaskQueue) executeTaskQueueItem(ctx context.Context, taskInfo TaskInfo, taskQueueInfo metadata.APITaskDetail) (execute bool) {

	locked, err := tq.lockTask(ctx, taskQueueInfo.TaskID)
	blog.Infof("start task %s", taskQueueInfo.TaskID)
	if err != nil {
		blog.Errorf("exceute task. lock error. task name:%s, taskID:%s, err:%s", taskInfo.Name, taskQueueInfo.TaskID, err.Error())
		time.Sleep(time.Second)
		return
	}
	if !locked {
		return
	}
	canExecute, err := tq.changeTaskToExecuting(ctx, taskQueueInfo.TaskID)
	blog.Infof("change task %s to executing, can execute %v", taskQueueInfo.TaskID, canExecute)
	if err != nil {
		tq.unLockTask(ctx, taskQueueInfo.TaskID)
		time.Sleep(time.Second)
		return
	}
	if !canExecute {
		if err := tq.unLockTask(ctx, taskQueueInfo.TaskID); err != nil {
			blog.Errorf("exceute task. task cann't execute. unlock error. task name:%s, taskID:%s, err:%s", taskInfo.Name, taskQueueInfo.TaskID, err.Error())
		}
		return
	}
	tq.executePush(ctx, taskInfo, &taskQueueInfo)
	return true
}

func (tq *TaskQueue) executePush(ctx context.Context, taskInfo TaskInfo, taskQueue *metadata.APITaskDetail) {
	var resp *metadata.Response
	var err error
	blog.InfoJSON("task execute task id:%s", taskQueue.TaskID)

	header := logics.GetDBHTTPHeader(taskQueue.Header)

	allSucc := true

	for _, subTask := range taskQueue.Detail {

		if subTask.Status == metadata.APITaskStatusSuccess {
			continue
		}

		if subTask.Status != metadata.APITaskStatusNew && subTask.Status != metadata.APITaskStatusWaitExecute {
			blog.ErrorJSON("task execute http do error. taskID:%s, taskqueue:%s, queue info: status not wait execute ", taskQueue.TaskID, taskQueue)
			allSucc = false
			break
		}
		for retry := int64(0); retry < taskInfo.Retry; retry++ {
			resp, err = tq.service.CoreAPI.TaskServer().Queue(taskInfo.Name).Post(ctx, header, taskInfo.Path, subTask.Data)
			if err != nil {
				time.Sleep(time.Millisecond * 100)
				blog.ErrorJSON("task execute http do error. taskID:%s, path:%s, taskName:%s, header: %s, err:%s", taskQueue.TaskID, taskInfo.Path, taskInfo.Name, header, err.Error())
				continue
			}
			break
		}

		updateConditon := mapstr.New()
		updateConditon.Set("task_id", taskQueue.TaskID)
		updateConditon.Set("detail.sub_task_id", subTask.SubTaskID)
		updateData := mapstr.New()
		errResponse := &metadata.Response{}

		if err != nil {
			ccErr := tq.service.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(taskQueue.Header)).CCError(common.CCErrCommHTTPDoRequestFailed)
			errResponse.Result = false
			errResponse.Code = ccErr.GetCode()
			errResponse.ErrMsg = ccErr.Error()
		} else {
			errResponse = resp
		}

		if err != nil || !resp.Result {
			allSucc = false
			updateData.Set("detail.$.status", metadata.APITAskStatusFail)
			updateData.Set("status", metadata.APITAskStatusFail)
		} else {
			updateData.Set("detail.$.status", metadata.APITaskStatusSuccess)
		}
		updateData.Set("detail.$.response", errResponse)
		updateData.Set(common.LastTimeField, time.Now())

		for dbRetry := 0; dbRetry < dbMaxRetry; dbRetry++ {
			dbErr := tq.service.DB.Table(common.BKTableNameAPITask).Update(ctx, updateConditon, updateData)
			if dbErr != nil {
				blog.ErrorJSON("task execute http do error. taskID:%s, err:%s", taskQueue.TaskID, dbErr)
				time.Sleep(time.Second * 3)
				continue
			}
			break
		}

	}

	// 所有任务执行完成，修改整个任务状态
	updateConditon := mapstr.New()
	updateConditon.Set("task_id", taskQueue.TaskID)
	updateData := mapstr.New()
	if allSucc {
		updateData.Set("status", metadata.APITaskStatusSuccess)
	} else {
		updateData.Set("status", metadata.APITAskStatusFail)
	}
	updateData.Set(common.LastTimeField, time.Now())

	for dbRetry := 0; dbRetry < dbMaxRetry; dbRetry++ {
		dbErr := tq.service.DB.Table(common.BKTableNameAPITask).Update(ctx, updateConditon, updateData)
		if dbErr != nil {
			blog.ErrorJSON("task execute http do error. taskID:%s, err:%s", taskQueue.TaskID, dbErr)
			time.Sleep(time.Second * 3)
			continue
		}
		break
	}

	return
}

func (tq *TaskQueue) lockTask(ctx context.Context, taskID string) (locked bool, err error) {

	key := fmt.Sprintf("%s:apiTask:%s", common.BKCacheKeyV3Prefix, taskID)
	locked, err = tq.service.CacheDB.SetNX(ctx, key, time.Now(), time.Minute*2).Result()
	if err != nil {
		blog.Errorf("lock task error. err:%s, taskID:%s", err.Error(), taskID)
		return false, tq.service.CCErr.Error("zh-cn", common.CCErrTaskLockedTaskFail)
	}
	return locked, nil
}

func (tq *TaskQueue) unLockTask(ctx context.Context, taskID string) (err error) {

	key := fmt.Sprintf("%s:apiTask:%s", common.BKCacheKeyV3Prefix, taskID)
	_, err = tq.service.CacheDB.Del(ctx, key).Result()
	if err != nil {
		blog.Errorf("unlock task error. err:%s, taskID:%s", err.Error(), taskID)
		return tq.service.CCErr.Error("zh-cn", common.CCErrTaskUnLockedTaskFail)
	}
	return nil
}

func (tq *TaskQueue) getWaitExectue(ctx context.Context, name string) ([]metadata.APITaskDetail, error) {

	cond := condition.CreateCondition()
	cond.Field("name").Eq(name)
	cond.Field("status").In([]metadata.APITaskStatus{metadata.APITaskStatusNew, metadata.APITaskStatusWaitExecute})

	rows := make([]metadata.APITaskDetail, 0)
	err := tq.service.DB.Table(common.BKTableNameAPITask).Find(cond.ToMapStr()).Sort("create_time").Limit(20).All(ctx, &rows)
	if err != nil {
		blog.ErrorJSON("query wait execute error:%s, task queue task:%s, cond:%s", err.Error(), name, cond.ToMapStr())
		return nil, tq.service.CCErr.Error("zh-cn", common.CCErrCommDBSelectFailed)
	}

	return rows, nil
}

func (tq *TaskQueue) changeTaskToExecuting(ctx context.Context, taskID string) (bool, error) {
	cond := condition.CreateCondition()
	cond.Field("task_id").Eq(taskID)
	cond.Field("status").In([]metadata.APITaskStatus{metadata.APITaskStatusNew, metadata.APITaskStatusWaitExecute})

	cnt, err := tq.service.DB.Table(common.BKTableNameAPITask).Find(cond.ToMapStr()).Count(ctx)
	if err != nil {
		blog.ErrorJSON("query wait execute error, taskID:%s, err:%s,, cond:%s", taskID, err.Error(), cond.ToMapStr())
		return false, tq.service.CCErr.Error("zh-cn", common.CCErrCommDBSelectFailed)
	}
	if cnt != 1 {
		blog.ErrorJSON("query wait execute , task not equal 1, taskID:%s, cond:%s, cnt:%s", taskID, cond.ToMapStr(), cnt)
		return false, nil
	}
	data := mapstr.MapStr{
		"status":             metadata.APITaskStatuExecute,
		common.LastTimeField: time.Now(),
	}
	err = tq.service.DB.Table(common.BKTableNameAPITask).Update(ctx, cond.ToMapStr(), data)
	if err != nil {
		blog.ErrorJSON("update task to execute error:%s, task queue task:%s, cond:%s", err.Error(), taskID, cond.ToMapStr())
		return false, tq.service.CCErr.Error("zh-cn", common.CCErrCommDBUpdateFailed)
	}
	return true, nil
}

func (tq *TaskQueue) compensate(ctx context.Context) {
	go func() {
		tq.compensateDBExecute(ctx)
		timer := time.NewTicker(time.Minute * 10)
		for range timer.C {
			if tq.close {
				return
			}
			tq.compensateDBExecute(ctx)

		}
	}()
}

func (tq *TaskQueue) compensateDBExecute(ctx context.Context) {
	cond := condition.CreateCondition()
	cond.Field("status").In([]metadata.APITaskStatus{metadata.APITaskStatuExecute})
	cond.Field(common.LastTimeField).Lt(time.Now().Add(-time.Minute * 20))
	data := mapstr.MapStr{
		"status":             metadata.APITaskStatusWaitExecute,
		common.LastTimeField: time.Now(),
	}
	err := tq.service.DB.Table(common.BKTableNameAPITask).Update(ctx, cond.ToMapStr(), data)
	if err != nil {
		blog.ErrorJSON("update task to wait execute error:%s, cond:%s", err.Error(), cond.ToMapStr())
	}
}

func (s *Service) initCodeTaskConfig() map[string]TaskInfo {
	taskInfoMap := make(map[string]TaskInfo, 0)
	codeTaskConfigArr := taskconfig.GetCodeTaskConfig()

	for _, codeTaskConfig := range codeTaskConfigArr {
		ti := TaskInfo{
			Name:  codeTaskConfig.Name,
			Retry: codeTaskConfig.Retry,
			Path:  codeTaskConfig.Path,
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
			panicErr := fmt.Sprintf("task code init. task:%s, svrType:%s, not exist", ti.Name, codeTaskConfig.SvrType)
			panic(panicErr)
		}
		taskInfoMap[ti.Name] = ti
	}

	return taskInfoMap
}
