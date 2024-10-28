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

package hostidentifier

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/event_server/types"
	"configcenter/src/thirdparty/apigw/gse"
	pushfile "configcenter/src/thirdparty/gse/push_file_forsyncdata"
)

const (
	// makeNewTaskFromFailHostCount the count about fail host to make a new task
	makeNewTaskFromFailHostCount = 50
	// redisTaskListName task list name in redis
	redisTaskListName = "host_identifier:task_list"
	// RedisFailHostListName fail host list name in redis
	RedisFailHostListName = "host_identifier:fail_host_list"
	// Handling the task is being processed
	Handling = 115
)

// Task 存到redis任务队列中的任务结构体
type Task struct {
	// 任务id
	TaskID string `json:"task_id"`
	// 该任务包含的主机信息
	HostInfos []*HostInfo `json:"host_infos"`
	// 当超过该时间后，视该任务已经超时，将此任务中还没有查到推送状态的主机视为推送失败
	ExpiredTime int64 `json:"expired_time"`
}

// MarshalBinary marshal Task struct
func (t Task) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

// HostInfo 存到redis的失败主机list的主机信息结构体
type HostInfo struct {
	// 主机id
	HostID int64 `json:"bk_host_id"`
	// 主机ip
	HostInnerIP string `json:"bk_host_innerip"`
	// 云区域id
	CloudID int64 `json:"bk_cloud_id"`
	// 重试了的次数
	Times int64 `json:"times"`
	// true或者false，为true时，表示已经拿到该主机的推送结果，false时未拿到
	HasResult bool `json:"has_result"`
	// 对于1.0的agent为cloudID:innerIP, 对于2.0的agent，为bk_agent_id
	AgentID      string `json:"bk_agent_id"`
	BKAddressing string `json:"bk_addressing"`
}

// MarshalBinary marshal HostInfo struct
func (h HostInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(h)
}

// GetTaskExecutionStatus get task execution status
func (h *HostIdentifier) GetTaskExecutionStatus() {
	for {
		if !h.engine.Discovery().IsMaster() {
			time.Sleep(time.Minute)
			continue
		}

		header, rid := newHeaderWithRid()
		// 1. 拿到任务
		task, err := h.getFromTaskList()
		if err != nil {
			blog.Errorf("get task from redis task list error, err: %v, rid: %s", err, rid)
			continue
		}

		// 2. 判断任务是否过期, 过期不处理
		if task.ExpiredTime < time.Now().Unix() {
			blog.Errorf("the task is expired, skip it, taskID: %s, rid: %s", task.TaskID, rid)
			continue
		}

		// 3. 拿任务执行结果
		taskResultMap, err := h.GetTaskExecutionResultMap(task, header, rid)
		if err != nil {
			blog.Errorf("get task result error, taskID: %s, err: %v, rid: %s", task.TaskID, err, rid)
			continue
		}

		// 4.遍历任务里的主机信息，与查到的任务结果进行对比，判断任务中的主机身份下发操作是否成功, 是否需要重新查任务状态
		failHosts, retry := h.compareTaskResult(task, taskResultMap)
		if len(failHosts) != 0 {
			h.addToFailHostList(failHosts)
		}

		// 5.该任务包含的主机还没有拿到全部的结果，并且还没超过规定时间时，把任务重新放入任务队列中
		if retry && time.Now().Unix() < task.ExpiredTime {
			if err := h.addToTaskList(task); err != nil {
				blog.Errorf("add task to redis list error, task: %v, err: %v, rid: %s", task, err, rid)
			}
		}
	}
}

func (h *HostIdentifier) compareTaskResult(task *Task, taskResultMap map[string]int64) ([]*HostInfo, bool) {

	// 此变量用于表示是否需要把task重新放回任务队列重新查询任务结果
	retry := false
	failHosts := make([]*HostInfo, 0)

	for _, hostInfo := range task.HostInfos {
		if hostInfo.HasResult {
			continue
		}

		key := hostInfo.AgentID
		if key == "" {
			key = HostKey(strconv.FormatInt(hostInfo.CloudID, 10), hostInfo.HostInnerIP)
		}
		code, exist := taskResultMap[key]
		// 如果拿不到主机下发的结果或者处于还在执行中的状态，那么表示这个task需要重新放入任务队列进行查询结果
		if !exist || code == Handling {
			retry = true
			blog.V(5).Infof("can not get push host identifier result, hostInfo: %v, taskID: %s", hostInfo, task.TaskID)
			continue
		}

		hostInfo.HasResult = true

		// 记录推送的失败主机
		if code != common.CCSuccess {
			blog.Errorf("push host identifier error, hostInfo: %v, taskID: %s", hostInfo, task.TaskID)
			failHosts = append(failHosts, hostInfo)
			h.metric.hostResultTotal.WithLabelValues("failed").Inc()
			continue
		}

		h.metric.hostResultTotal.WithLabelValues("success").Inc()
		blog.V(5).Infof("push identifier to host success, host: %v, taskID: %s", hostInfo, task.TaskID)
	}

	return failHosts, retry
}

// GetTaskExecutionResultMap get task execution result map from gse
func (h *HostIdentifier) GetTaskExecutionResultMap(task *Task, header http.Header, rid string) (map[string]int64,
	error) {
	switch h.apiVersion {
	case types.V2:
		return h.getV2TaskExecutionResultMap(task, header, rid)

	case types.V1:
		return h.getV1TaskExecutionResultMap(task, rid)
	}

	return nil, fmt.Errorf("can not find the version about gse client, version: %s", h.apiVersion)
}

func (h *HostIdentifier) getV2TaskExecutionResultMap(task *Task, header http.Header, rid string) (map[string]int64,
	error) {
	var err error
	resp := new(gse.GetTransferFileResultResp)
	failCount := 0
	agentIDList := make([]string, 0)
	for _, host := range task.HostInfos {
		agentIDList = append(agentIDList, host.AgentID)
	}
	req := &gse.GetTransferFileResultRequest{
		TaskID:      task.TaskID,
		AgentIDList: agentIDList,
	}
	for failCount < retryTimes {
		resp, err = h.gseApiGWClient.GetTransferFileResult(h.ctx, header, req)
		if err != nil {
			blog.Errorf("get task status from gse error, task: %v, err: %v, rid: %s", task, err, rid)
			h.metric.getResultTotal.WithLabelValues("failed").Inc()
			failCount++
			sleepForFail(failCount)
			continue
		}
		break
	}
	if failCount >= retryTimes {
		return nil, errors.New("get task push result error")
	}

	h.metric.getResultTotal.WithLabelValues("success").Inc()
	return buildV2TaskResultMap(resp.Data.Result), nil
}

func (h *HostIdentifier) getV1TaskExecutionResultMap(task *Task, rid string) (map[string]int64, error) {
	var err error
	resp := new(pushfile.API_MapRsp)
	failCount := 0
	for failCount < retryTimes {
		resp, err = h.gseTaskServerClient.GetPushFileRst(h.ctx, task.TaskID)
		if err != nil {
			blog.Errorf("get task status from gse error, task: %v, err: %v, rid: %s", task, err, rid)
			h.metric.getResultTotal.WithLabelValues("failed").Inc()
			failCount++
			sleepForFail(failCount)
			continue
		}

		if resp.MErrcode != common.CCSuccess {
			blog.Errorf("get task status from gse fail, task: %v, code: %d, msg: %s, rid: %s", task, resp.MErrcode,
				resp.MErrmsg, rid)
			h.metric.getResultTotal.WithLabelValues("failed").Inc()
			failCount++
			sleepForFail(failCount)
			continue
		}
		break
	}
	if failCount >= retryTimes {
		return nil, errors.New("get task push result error")
	}

	h.metric.getResultTotal.WithLabelValues("success").Inc()
	return buildV1TaskResultMap(resp.MRsp), nil
}

// LaunchTaskForFailedHost launch task for failed host
func (h *HostIdentifier) LaunchTaskForFailedHost() {
	for {
		if !h.engine.Discovery().IsMaster() {
			time.Sleep(time.Minute)
			continue
		}

		// 1、收集失败的主机
		header, rid := newHeaderWithRid()
		hostInfoArray, statusReq, success := h.collectFailHost(rid)
		if !success {
			continue
		}

		// 2、查询主机的agent状态
		resp, err := h.getAgentStatus(statusReq, header, rid)
		if err != nil {
			blog.Errorf("get agent status error, hostInfo: %v, err: %v, rid: %s", hostInfoArray, err, rid)
			continue
		}

		// 3、将处于on状态的主机拿出来构造推送信息
		hostIDs, hostInfos, hostMap := h.getOnStatusAgentFromHostInfo(hostInfoArray, resp, rid)
		if len(hostIDs) == 0 {
			blog.Warnf("get fail host success, but agent status is off, hostInfos: %v, rid: %s", hostInfoArray, rid)
			continue
		}

		// 4、查询主机身份并推送
		if _, err := h.getHostIdentifierAndPush(hostIDs, hostMap, hostInfos, rid, header); err != nil {
			blog.Errorf("launch task for failed host error, err: %v, rid: %s", err, rid)
		}
	}
}

func (h *HostIdentifier) getOnStatusAgentFromHostInfo(hosts []*HostInfo, statusMap map[string]string, rid string) (
	[]int64, []*HostInfo, map[int64]string) {

	switch h.apiVersion {
	case types.V2:
		return h.getV2OnStatusAgentFromHostInfo(hosts, statusMap, rid)

	case types.V1:
		return h.getV1OnStatusAgentFromHostInfo(hosts, statusMap, rid)
	}

	return nil, nil, nil
}

func (h *HostIdentifier) getV2OnStatusAgentFromHostInfo(hosts []*HostInfo, statusMap map[string]string, rid string) (
	[]int64, []*HostInfo, map[int64]string) {

	hostIDs := make([]int64, 0)
	hostInfos := make([]*HostInfo, 0)
	// 此map保存hostID和该host处于on的agent的agentID的对应关系
	hostMap := make(map[int64]string)
	for _, hostInfo := range hosts {
		if statusMap[hostInfo.AgentID] != v2ApiAgentOnStatus {
			blog.Infof("agent status is off, agentID: %s, hostID: %d, rid: %s", hostInfo.AgentID, hostInfo.HostID, rid)
			h.metric.agentStatusTotal.WithLabelValues("off").Inc()
			continue
		}

		h.metric.agentStatusTotal.WithLabelValues("on").Inc()
		hostIDs = append(hostIDs, hostInfo.HostID)
		hostMap[hostInfo.HostID] = hostInfo.AgentID
		hostInfos = append(hostInfos, hostInfo)
	}

	return hostIDs, hostInfos, hostMap
}

func (h *HostIdentifier) getV1OnStatusAgentFromHostInfo(hosts []*HostInfo, statusMap map[string]string, rid string) (
	[]int64, []*HostInfo, map[int64]string) {

	hostIDs := make([]int64, 0)
	// 此map保存hostID和该host处于on的agent的ip的对应关系
	hostMap := make(map[int64]string)
	hostInfos := make([]*HostInfo, 0)
	for _, hostInfo := range hosts {
		cloudID := strconv.FormatInt(hostInfo.CloudID, 10)
		isOn, hostIP := getStatusOnAgentIP(cloudID, hostInfo.HostInnerIP, statusMap)
		if !isOn {
			blog.Infof("host %v agent status is off, rid: %s", hostInfo, rid)
			continue
		}

		h.metric.agentStatusTotal.WithLabelValues("on").Inc()
		hostIDs = append(hostIDs, hostInfo.HostID)
		hostMap[hostInfo.HostID] = hostIP
		hostInfo.HostInnerIP = hostIP
		hostInfos = append(hostInfos, hostInfo)
	}

	return hostIDs, hostInfos, hostMap
}

// collectFailHost collect fail host
func (h *HostIdentifier) collectFailHost(rid string) ([]*HostInfo, []StatusReq, bool) {
	start := time.Now()
	hostInfoArray := make([]*HostInfo, 0)
	statusReqList := make([]StatusReq, 0)
	uniqueMap := make(map[int64]struct{})

	// 从redis的保存失败的主机的list中拿出一定数量主机，并进行去重
	for time.Now().Sub(start) < time.Minute {
		if len(hostInfoArray) >= makeNewTaskFromFailHostCount {
			break
		}

		val, err := h.redisCli.LLen(context.Background(), RedisFailHostListName).Result()
		if err != nil {
			blog.Errorf("get fail_host_list list length error, err: %v, rid: %s", err, rid)
			continue
		}

		if val == 0 {
			time.Sleep(time.Second)
			continue
		}

		hostInfo, err := h.getFromFailHostList()
		if err != nil {
			blog.Errorf("get host from redis fail_host_list error, err: %v, rid: %s", err, rid)
			continue
		}

		// 去重
		if _, ok := uniqueMap[hostInfo.HostID]; ok {
			continue
		}
		uniqueMap[hostInfo.HostID] = struct{}{}

		statusReqList = append(statusReqList, StatusReq{
			CloudID:      strconv.FormatInt(hostInfo.CloudID, 10),
			InnerIP:      hostInfo.HostInnerIP,
			BKAddressing: hostInfo.BKAddressing,
			AgentID:      hostInfo.AgentID,
		})
	}

	return hostInfoArray, statusReqList, len(hostInfoArray) > 0
}

// pushFile push host identifier file to gse and create a new task to redis task_list
func (h *HostIdentifier) pushFile(always bool, hostInfos []*HostInfo, taskInfo *TaskInfo, header http.Header,
	rid string) (*Task, error) {

	if taskInfo == nil || (taskInfo.V1Task == nil && taskInfo.V2Task == nil) {
		return nil, errors.New("the task info about host is empty")
	}

	// 1、调用gse taskServer接口，推送主机身份
	var taskID string
	var err error
	switch h.apiVersion {
	case types.V1:
		taskID, err = h.pushFileByV1Api(taskInfo.V1Task, rid)
		if err != nil {
			return nil, err
		}

	case types.V2:
		taskID, err = h.pushFileByV2Api(taskInfo.V2Task, header, rid)
		if err != nil {
			return nil, err
		}
	}

	h.metric.pushFileTotal.WithLabelValues("success").Inc()
	blog.V(5).Infof("push host identifier to gse success: taskID: %s, rid: %s", taskID, rid)

	// 2、构建task放到redis维护的任务队列中
	task := &Task{
		TaskID:      taskID,
		HostInfos:   hostInfos,
		ExpiredTime: time.Now().Add(50 * time.Minute).Unix(),
	}
	failCount := 0
	for failCount < retryTimes {
		if err = h.addToTaskList(task); err != nil {
			failCount++
			sleepForFail(failCount)
			continue
		}
		break
	}
	if failCount >= retryTimes {
		blog.Errorf("add task to redis error, taskID: %s, err: %v, rid: %s", taskID, err, rid)
		return nil, err
	}
	return task, nil
}

func (h *HostIdentifier) pushFileByV2Api(task []*gse.Task, header http.Header, rid string) (string,
	error) {

	if len(task) == 0 {
		blog.Errorf("push file error, because the task is empty, rid: %s", rid)
		return "", errors.New("push file error, because the task is empty")
	}

	failCount := 0
	var err error
	var resp *gse.AsyncPushFileResp
	req := &gse.AsyncPushFileRequest{
		Tasks:          task,
		AutoMkdir:      true,
		TimeoutSeconds: 1000,
	}

	for failCount < retryTimes {
		resp, err = h.gseApiGWClient.AsyncPushFile(h.ctx, header, req)
		if err != nil {
			h.metric.pushFileTotal.WithLabelValues("failed").Inc()
			failCount++
			sleepForFail(failCount)
			continue
		}
		break
	}

	if failCount >= retryTimes {
		return "", err
	}

	return resp.Data.Result.TaskID, nil
}

func (h *HostIdentifier) pushFileByV1Api(task []*pushfile.API_FileInfoV2, rid string) (string, error) {
	if len(task) == 0 {
		blog.Errorf("push file error, because the task is empty, rid: %s", rid)
		return "", errors.New("push file error, because the task is empty")
	}

	var err error
	failCount := 0
	resp := new(pushfile.API_CommRsp)

	for failCount < retryTimes {
		resp, err = h.gseTaskServerClient.PushFileV2(context.Background(), task)
		if err != nil {
			blog.Errorf("push host identifier to gse error, err: %v, rid: %s", err, rid)
			h.metric.pushFileTotal.WithLabelValues("failed").Inc()
			failCount++
			sleepForFail(failCount)
			continue
		}

		if resp.MErrcode != common.CCSuccess {
			blog.Errorf("push host identifier fail, code: %d, msg: %s, rid: %s", resp.MErrcode, resp.MErrmsg, rid)
			h.metric.pushFileTotal.WithLabelValues("failed").Inc()
			failCount++
			sleepForFail(failCount)
			continue
		}
		break
	}

	if failCount >= retryTimes {
		return "", errors.New("push host identifier to gse taskServer error")
	}

	return resp.MContent, nil
}

// addToTaskList add task to redis task list
func (h *HostIdentifier) addToTaskList(task *Task) error {
	return h.redisCli.RPush(context.Background(), redisTaskListName, task).Err()
}

// getFromTaskList get task from redis task list
func (h *HostIdentifier) getFromTaskList() (*Task, error) {
	result, err := h.redisCli.BLPop(context.Background(), 0, redisTaskListName).Result()
	if err != nil {
		return nil, err
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("redis data is invalid, data: %v", result)
	}

	task := new(Task)
	if err := json.Unmarshal([]byte(result[1]), task); err != nil {
		blog.Errorf("Unmarshal task error, task: %s, err: %v", result[1], err)
		return nil, err
	}

	return task, nil
}

// addToFailHostList add host to redis fail host list
func (h *HostIdentifier) addToFailHostList(hosts []*HostInfo) {
	for _, host := range hosts {
		if host.Times >= retryTimes {
			blog.Errorf("host exceed the max retry times, hostInfo; %v", host)
			continue
		}

		host.Times++
		if err := h.redisCli.RPush(context.Background(), RedisFailHostListName, host).Err(); err != nil {
			blog.Errorf("add fail host to redis list error, hostInfo; %v, err: %v", host, err)
		}
	}
}

// getFromFailHostList get fail host from redis fail host list
func (h *HostIdentifier) getFromFailHostList() (*HostInfo, error) {
	result, err := h.redisCli.BLPop(context.Background(), 30*time.Second, RedisFailHostListName).Result()
	if err != nil {
		return nil, err
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("redis data is invalid, data: %v", result)
	}

	hostInfo := new(HostInfo)
	if err := json.Unmarshal([]byte(result[1]), hostInfo); err != nil {
		blog.Errorf("Unmarshal hostInfo error, hostInfo: %s, err: %v", result[1], err)
	}

	return hostInfo, nil
}
