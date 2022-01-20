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
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/thirdparty/gse/get_agent_state_forsyncdata"
	"configcenter/src/thirdparty/gse/push_file_forsyncdata"

	"github.com/tidwall/gjson"
)

const (
	// makeNewTaskFromFailHostCount the count about fail host to make a new task
	makeNewTaskFromFailHostCount = 50
	// redisTaskListName task list name in redis
	redisTaskListName = "host_identifier:task_list"
	// redisFailHostListName fail host list name in redis
	RedisFailHostListName = "host_identifier:fail_host_list"
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

		task, err := h.getFromTaskList()
		if err != nil {
			blog.Errorf("get task from redis task list error, err: %v", err)
			continue
		}

		taskResultMap, err := h.GetTaskExecutionResultMap(task)
		if err != nil {
			blog.Errorf("get task result error, taskID: %s, err: %v", task.TaskID, err)
			continue
		}

		// 遍历任务里的主机信息，与查到的任务结果进行对比，判断任务中的主机身份下发操作是否成功
		for _, hostInfo := range task.HostInfos {
			if hostInfo.HasResult {
				continue
			}
			key := HostKey(strconv.FormatInt(hostInfo.CloudID, 10), hostInfo.HostInnerIP)
			if taskResultMap[key] == "" {
				blog.Errorf("can not get host identifier push message from task, hostInfo: %v, taskID: %s",
					hostInfo, task.TaskID)
				// 超过规定时间还没有拿到结果，并且没超过最大重试次数时，将没拿到结果的主机信息放到失败主机队列中
				if time.Now().Unix() >= task.ExpiredTime && hostInfo.Times < retryTimes {
					hostInfo.Times++
					if err := h.addToFailHostList(hostInfo); err != nil {
						blog.Errorf("add fail host to redis list error, hostInfo; %v, err: %v", hostInfo, err)
					}
				}
				continue
			}

			hostInfo.HasResult = true
			// 把推送失败且没超过最大重试次数的主机放到失败主机队列中
			if gjson.Get(taskResultMap[key], "error_code").Int() != common.CCSuccess {
				blog.Errorf("push host identifier error, hostInfo: %v, taskID: %s", hostInfo, task.TaskID)
				if hostInfo.Times < retryTimes {
					hostInfo.Times++
					if err := h.addToFailHostList(hostInfo); err != nil {
						blog.Errorf("add fail host to redis list error, hostInfo; %v, err: %v", hostInfo, err)
					}
				}
				continue
			}

			blog.V(5).Infof("push identifier to host success, host: %v, taskID: %s", hostInfo, task.TaskID)
		}

		// 该任务包含的主机还没有拿到全部的结果，并且还没超过规定时间时，把任务重新放入任务队列中
		if len(task.HostInfos) != len(taskResultMap) && time.Now().Unix() < task.ExpiredTime {
			if err := h.addToTaskList(task); err != nil {
				blog.Errorf("add task to redis list error, task: %v, err: %v", task, err)
			}
		}
	}
}

// GetTaskExecutionResultMap get task execution result map from gse
func (h *HostIdentifier) GetTaskExecutionResultMap(task *Task) (map[string]string, error) {
	var err error
	taskStatusResp := new(push_file_forsyncdata.API_MapRsp)
	failCount := 0
	for failCount < retryTimes {
		taskStatusResp, err = h.gseTaskServerClient.GetPushFileRst(h.ctx, task.TaskID)
		if err != nil {
			blog.Errorf("get task status from gse error, task: %v, err: %v", task, err)
			failCount++
			sleepForFail(failCount)
			continue
		}

		if taskStatusResp.MErrcode != common.CCSuccess {
			blog.Errorf("get task status from gse fail, task: %v, errCode: %d, errMessage: %s",
				task, taskStatusResp.MErrcode, taskStatusResp.MErrmsg)
			failCount++
			sleepForFail(failCount)
			continue
		}
		break
	}
	if failCount >= retryTimes {
		return nil, err
	}

	return buildTaskResultMap(taskStatusResp.MRsp), nil
}

// LaunchTaskForFailedHost launch task for failed host
func (h *HostIdentifier) LaunchTaskForFailedHost() {
	for {
		if !h.engine.Discovery().IsMaster() {
			time.Sleep(time.Minute)
			continue
		}

		hostInfoArray, agentStatusRequest, success := h.collectFailHost()
		if !success {
			continue
		}

		// 查询主机状态
		agentStatus := new(get_agent_state_forsyncdata.AgentStatusResponse)
		var err error
		failCount := 0
		for failCount < retryTimes {
			agentStatus, err = h.gseApiServerClient.GetAgentStatus(h.ctx, agentStatusRequest)
			if err != nil {
				blog.Errorf("get host agent status error: err: %v", agentStatus.BkErrorCode, agentStatus.BkErrorMsg)
				failCount++
				sleepForFail(failCount)
				continue
			}

			if agentStatus.BkErrorCode != common.CCSuccess {
				blog.Errorf("get host agent status fail, errCode: %d, errMessage: %s",
					agentStatus.BkErrorCode, agentStatus.BkErrorMsg)
				failCount++
				sleepForFail(failCount)
				continue
			}
			break
		}
		if failCount >= retryTimes {
			continue
		}

		// 将处于on状态的主机拿出来构造推送信息
		hostIDs := make([]int64, 0)
		hostInfos := make([]*HostInfo, 0)
		// 此map保存hostID和该host处于on的agent的ip的对应关系
		hostMap := make(map[int64]string)
		for _, hostInfo := range hostInfoArray {
			isOn, hostIP := getStatusOnAgentIP(strconv.FormatInt(hostInfo.CloudID, 10),
				hostInfo.HostInnerIP, agentStatus.Result_)
			if !isOn {
				continue
			}
			hostIDs = append(hostIDs, hostInfo.HostID)
			hostMap[hostInfo.HostID] = hostIP
			hostInfo.HostInnerIP = hostIP
			hostInfos = append(hostInfos, hostInfo)
		}

		if len(hostIDs) == 0 {
			blog.Warnf("get fail host from redis success, but host agent status is off, hostInfos: %v", hostInfoArray)
			continue
		}

		if _, err := h.getHostIdentifierAndPush(hostIDs, hostMap, hostInfos); err != nil {
			blog.Errorf("launch task for failed host error, err: %v", err)
		}
	}
}

// collectFailHost collect fail host
func (h *HostIdentifier) collectFailHost() ([]*HostInfo, *get_agent_state_forsyncdata.AgentStatusRequest, bool) {
	start := time.Now()
	hostInfoArray := make([]*HostInfo, 0)
	agentStatusRequest := new(get_agent_state_forsyncdata.AgentStatusRequest)
	uniqueMap := make(map[int64]bool)

	// 从redis的保存失败的主机的list中拿出一定数量主机，并进行去重
	for i := 0; i < makeNewTaskFromFailHostCount; i++ {
		if time.Now().Sub(start) > 5*time.Minute {
			break
		}

		hostInfoMessage, err := h.getFromFailHostList()
		if err != nil {
			blog.Errorf("get host from redis fail_host_list error, err: %v", err)
			continue
		}
		hostInfo := new(HostInfo)
		if err := json.Unmarshal([]byte(hostInfoMessage), hostInfo); err != nil {
			blog.Errorf("Unmarshal hostInfo error, hostInfo: %s, err: %v", hostInfoMessage, err)
			continue
		}

		// 去重
		if uniqueMap[hostInfo.HostID] {
			continue
		}
		uniqueMap[hostInfo.HostID] = true

		agentStatusRequest.Hosts = append(agentStatusRequest.Hosts, &get_agent_state_forsyncdata.CacheIPInfo{
			GseCompositeID: strconv.FormatInt(hostInfo.CloudID, 10),
			IP:             hostInfo.HostInnerIP,
		})
		hostInfoArray = append(hostInfoArray, hostInfo)
	}

	return hostInfoArray, agentStatusRequest, len(hostInfoArray) > 0
}

// pushFile push host identifier file to gse and create a new task to redis task_list
func (h *HostIdentifier) pushFile(always bool, hostInfos []*HostInfo,
	fileList []*push_file_forsyncdata.API_FileInfoV2) (*Task, error) {

	var err error
	pushFileResp := new(push_file_forsyncdata.API_CommRsp)
	failCount := 0
	// 调用gse taskServer接口，推送主机身份
	for always || failCount < retryTimes {
		pushFileResp, err = h.gseTaskServerClient.PushFileV2(h.ctx, fileList)
		if err != nil {
			blog.Errorf("push host identifier to gse error, file: %v, err: %v", fileList, err)
			failCount++
			sleepForFail(failCount)
			continue
		}

		if pushFileResp.MErrcode != common.CCSuccess {
			blog.Errorf("push host identifier to gse fail, file: %v, errCode: %d, errMessage: %s",
				fileList, pushFileResp.MErrcode, pushFileResp.MErrmsg)
			failCount++
			sleepForFail(failCount)
			continue
		}
		break
	}

	if !always && failCount >= retryTimes {
		return nil, err
	}
	blog.V(5).Infof("push host identifier to gse success: file: %v, taskID: %s", fileList, pushFileResp.MContent)

	// 构建task放到redis维护的任务队列中
	task := &Task{
		TaskID:      pushFileResp.MContent,
		HostInfos:   hostInfos,
		ExpiredTime: time.Now().Add(3 * time.Hour).Unix(),
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
		blog.Errorf("add task to redis task list error, taskID: %s, hostInfos: %v, err: %v",
			pushFileResp.MContent, hostInfos, err)
		return nil, err
	}
	return task, nil
}

// addToTaskList add task to redis task list
func (h *HostIdentifier) addToTaskList(task *Task) error {
	return h.redisCli.RPush(h.ctx, redisTaskListName, task).Err()
}

// GetFromTaskList get task from redis task list
func (h *HostIdentifier) getFromTaskList() (*Task, error) {
	result, err := h.redisCli.BLPop(h.ctx, 0, redisTaskListName).Result()
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
func (h *HostIdentifier) addToFailHostList(host *HostInfo) error {
	return h.redisCli.RPush(h.ctx, RedisFailHostListName, host).Err()
}

// getFailHostList get fail host from redis fail host list
func (h *HostIdentifier) getFromFailHostList() (string, error) {
	result, err := h.redisCli.BLPop(h.ctx, 5*time.Minute, RedisFailHostListName).Result()
	if err != nil {
		return "", err
	}
	return result[1], nil
}
