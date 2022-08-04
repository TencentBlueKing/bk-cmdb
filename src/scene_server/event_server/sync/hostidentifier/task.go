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
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	getstatus "configcenter/src/thirdparty/gse/get_agent_state_forsyncdata"
	pushfile "configcenter/src/thirdparty/gse/push_file_forsyncdata"

	"github.com/tidwall/gjson"
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

		// 1. 拿到任务
		task, err := h.getFromTaskList()
		if err != nil {
			blog.Errorf("get task from redis task list error, err: %v", err)
			continue
		}

		// 2. 判断任务是否过期, 过期不处理
		if task.ExpiredTime < time.Now().Unix() {
			blog.Errorf("the task is expired, skip it, taskID: %s", task.TaskID)
			continue
		}

		// 3. 拿任务执行结果
		taskResultMap, err := h.GetTaskExecutionResultMap(task)
		if err != nil {
			blog.Errorf("get task result error, taskID: %s, err: %v", task.TaskID, err)
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
				blog.Errorf("add task to redis list error, task: %v, err: %v", task, err)
			}
		}
	}
}

func (h *HostIdentifier) compareTaskResult(task *Task, taskResultMap map[string]string) ([]*HostInfo, bool) {

	// 此变量用于表示是否需要把task重新放回任务队列重新查询任务结果
	retry := false
	failHosts := make([]*HostInfo, 0)

	for _, hostInfo := range task.HostInfos {
		if hostInfo.HasResult {
			continue
		}

		key := HostKey(strconv.FormatInt(hostInfo.CloudID, 10), hostInfo.HostInnerIP)
		code := gjson.Get(taskResultMap[key], "error_code").Int()

		// 如果拿不到主机下发的结果或者处于还在执行中的状态，那么表示这个task需要重新放入任务队列进行查询结果
		if taskResultMap[key] == "" || code == Handling {
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
func (h *HostIdentifier) GetTaskExecutionResultMap(task *Task) (map[string]string, error) {
	var err error
	resp := new(pushfile.API_MapRsp)
	failCount := 0
	for failCount < retryTimes {
		resp, err = h.gseTaskServerClient.GetPushFileRst(h.ctx, task.TaskID)
		if err != nil {
			blog.Errorf("get task status from gse error, task: %v, err: %v", task, err)
			h.metric.getResultTotal.WithLabelValues("failed").Inc()
			failCount++
			sleepForFail(failCount)
			continue
		}

		if resp.MErrcode != common.CCSuccess {
			blog.Errorf("get task status from gse fail, task: %v, code: %d, msg: %s", task, resp.MErrcode, resp.MErrmsg)
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
	return buildTaskResultMap(resp.MRsp), nil
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
		resp, err := h.getAgentStatus(statusReq, false, rid)
		if err != nil {
			blog.Errorf("get agent status error, hostInfo: %v, err: %v, rid: %s", hostInfoArray, err, rid)
			continue
		}

		// 3、将处于on状态的主机拿出来构造推送信息
		hostIDs := make([]int64, 0)
		hostInfos := make([]*HostInfo, 0)
		// 此map保存hostID和该host处于on的agent的ip的对应关系
		hostMap := make(map[int64]string)
		for _, hostInfo := range hostInfoArray {
			cloudID := strconv.FormatInt(hostInfo.CloudID, 10)
			isOn, hostIP := getStatusOnAgentIP(cloudID, hostInfo.HostInnerIP, resp.Result_)
			if !isOn {
				blog.Infof("host %v agent status is off, rid: %s", hostInfo, rid)
				continue
			}

			blog.Infof("host %v agent status is on, ip: %s, rid: %s", hostInfo, hostIP, rid)

			hostIDs = append(hostIDs, hostInfo.HostID)
			hostMap[hostInfo.HostID] = hostIP
			hostInfo.HostInnerIP = hostIP
			hostInfos = append(hostInfos, hostInfo)
		}

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

// collectFailHost collect fail host
func (h *HostIdentifier) collectFailHost(rid string) ([]*HostInfo, *getstatus.AgentStatusRequest, bool) {
	start := time.Now()
	hostInfoArray := make([]*HostInfo, 0)
	agentStatusRequest := new(getstatus.AgentStatusRequest)
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

		agentStatusRequest.Hosts = append(agentStatusRequest.Hosts, &getstatus.CacheIPInfo{
			GseCompositeID: strconv.FormatInt(hostInfo.CloudID, 10),
			IP:             hostInfo.HostInnerIP,
		})
		hostInfoArray = append(hostInfoArray, hostInfo)
	}

	return hostInfoArray, agentStatusRequest, len(hostInfoArray) > 0
}

// pushFile push host identifier file to gse and create a new task to redis task_list
func (h *HostIdentifier) pushFile(always bool, hostInfos []*HostInfo, fileList []*pushfile.API_FileInfoV2,
	rid string) (*Task, error) {

	var err error
	failCount := 0
	resp := new(pushfile.API_CommRsp)

	// 1、调用gse taskServer接口，推送主机身份
	for always || failCount < retryTimes {
		resp, err = h.gseTaskServerClient.PushFileV2(context.Background(), fileList)
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

	if !always && failCount >= retryTimes {
		return nil, errors.New("push host identifier to gse taskServer error")
	}

	h.metric.pushFileTotal.WithLabelValues("success").Inc()
	blog.V(5).Infof("push host identifier to gse success: file: %v, taskID: %s, rid: %s", fileList, resp.MContent, rid)

	// 2、构建task放到redis维护的任务队列中
	task := &Task{
		TaskID:      resp.MContent,
		HostInfos:   hostInfos,
		ExpiredTime: time.Now().Add(50 * time.Minute).Unix(),
	}
	failCount = 0
	for failCount < retryTimes {
		if err = h.addToTaskList(task); err != nil {
			failCount++
			sleepForFail(failCount)
			continue
		}
		break
	}
	if failCount >= retryTimes {
		blog.Errorf("add task to redis error, taskID: %s, err: %v, rid: %s", resp.MContent, err, rid)
		return nil, err
	}
	return task, nil
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
