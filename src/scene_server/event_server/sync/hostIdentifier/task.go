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

package hostIdentifier

import (
	"encoding/json"
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
	TaskID      string      `json:"task_id"`      // 任务id
	HostInfos   []*HostInfo `json:"host_infos"`   // 该任务包含的主机信息
	ExpiredTime int64       `json:"expired_time"` // 当超过该时间后，视该任务已经超时，将此任务中还没有查到推送状态的主机视为推送失败
}

// MarshalBinary marshal Task struct
func (t Task) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

// HostInfo 存到redis的失败主机list的主机信息结构体
type HostInfo struct {
	HostInnerIP string `json:"bk_host_innerip"` // 主机ip
	CloudID     int64  `json:"bk_cloud_id"`     // 云区域id
	Times       int64  `json:"times"`           // 重试了的次数
	HasResult   bool   `json:"has_result"`      // true或者false，为true时，表示已经拿到该主机的推送结果，false时未拿到
}

// MarshalBinary marshal HostInfo struct
func (h HostInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(h)
}

// GetTaskExecutionStatus get task execution status
func (h *HostIdentifier) GetTaskExecutionStatus() {
	for {
		if !h.engine.Discovery().IsMaster() {
			return
		}

		taskMessage, err := h.getTaskFromTaskList()
		if err != nil {
			blog.Errorf("get task from redis task list error, err: %v", err)
			continue
		}
		task := new(Task)
		if err := json.Unmarshal([]byte(taskMessage), task); err != nil {
			blog.Errorf("Unmarshal task error, task: %s, err: %v", taskMessage, err)
			continue
		}
		taskStatusResp := new(push_file_forsyncdata.API_MapRsp)
		failCount := 0
		for failCount < retryTimes {
			taskStatusResp, err = h.gseTaskServerClient.GetPushFileRst(h.ctx, task.TaskID)
			if err != nil || taskStatusResp.MErrcode != common.CCSuccess {
				blog.Errorf("get task status from gse error, taskMessage: %s, err: %v", taskMessage, err)
				failCount++
				sleepForFail(failCount)
				continue
			}
			break
		}
		if failCount >= retryTimes {
			failCount = 0
			continue
		}

		taskResultMap := buildTaskResultMap(taskStatusResp.MRsp)
		for _, hostInfo := range task.HostInfos {
			if hostInfo.HasResult {
				continue
			}
			key := strconv.FormatInt(hostInfo.CloudID, 10) + ":" + hostInfo.HostInnerIP
			if taskResultMap[key] == "" {
				blog.Errorf("can not get host identifier push message from task, hostInfo: %v, taskID: %s",
					hostInfo, task.TaskID)
				// 超过规定时间还没有拿到结果，并且没超过最大重试次数时，将没拿到结果的主机信息放到失败主机队列中
				if time.Now().Unix() >= task.ExpiredTime && hostInfo.Times < retryTimes {
					hostInfo.Times++
					h.addHostToFailHostList(hostInfo)
				}
				continue
			}

			// 把推送失败且没超过最大重试次数的主机放到失败主机队列中
			if gjson.Get(taskResultMap[key], "error_code").Int() != common.CCSuccess {
				blog.Errorf("push host identifier error, hostInfo: %v, taskID: %s", hostInfo, task.TaskID)
				if hostInfo.Times < retryTimes {
					hostInfo.Times++
					h.addHostToFailHostList(hostInfo)
				}
			}
			hostInfo.HasResult = true
		}

		// 该任务包含的主机还没有拿到全部的结果，并且还没超过规定时间时，把任务重新放入任务队列中
		if len(task.HostInfos) != len(taskResultMap) && time.Now().Unix() < task.ExpiredTime {
			h.addTaskToTaskList(task)
		}
	}
}

// MakeNewTaskFromFailHost make new task from redis fail host list
func (h *HostIdentifier) MakeNewTaskFromFailHost() {
	for {
		if !h.engine.Discovery().IsMaster() {
			return
		}

		hostInfoArray, agentStatusRequest := h.collectFailHost()
		if len(hostInfoArray) == 0 {
			continue
		}

		// 查询主机状态
		agentStatus := new(get_agent_state_forsyncdata.AgentStatusResponse)
		var err error
		failCount := 0
		for failCount < retryTimes {
			agentStatus, err = h.gseApiServerClient.GetAgentStatus(h.ctx, agentStatusRequest)
			if err != nil || agentStatus.BkErrorCode != common.CCSuccess {
				blog.Errorf("get host agent status error: err: %v, errCode: %d, errMessage: %s",
					err, agentStatus.BkErrorCode, agentStatus.BkErrorMsg)
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
		var fileList []*push_file_forsyncdata.API_FileInfoV2
		var hostInfos []*HostInfo
		for _, hostInfo := range hostInfoArray {
			key := strconv.FormatInt(hostInfo.CloudID, 10) + ":" + hostInfo.HostInnerIP
			if gjson.Get(agentStatus.Result_[key], "bk_agent_alive").Int() == agentOnStatus {
				header, rid := newHeaderWithRid()
				identifier, err := h.findHostIdentifier(hostInfo.CloudID, hostInfo.HostInnerIP, header)
				if err != nil {
					blog.Errorf("search identifier error, cloudID: %d, innerIP: %s, err: %v, rid: %s",
						hostInfo.CloudID, hostInfo.HostInnerIP, err, header, rid)
					continue
				}
				fileList = append(fileList, h.buildPushFile(identifier, hostInfo.HostInnerIP, hostInfo.CloudID))
				hostInfos = append(hostInfos, hostInfo)
				continue
			}
			blog.Infof("host %v agent status is off", key)
		}

		if len(fileList) == 0 {
			continue
		}

		// 推送主机身份信息
		failCount = 0
		pushFileResp := new(push_file_forsyncdata.API_CommRsp)
		for failCount < retryTimes {
			pushFileResp, err = h.gseTaskServerClient.PushFileV2(h.ctx, fileList)
			if err != nil || pushFileResp.MErrcode != common.CCSuccess {
				blog.Errorf("push host identifier to gse error: err: %v, errCode: %d, errMessage: %s",
					err, pushFileResp.MErrcode, pushFileResp.MErrmsg)
				failCount++
				sleepForFail(failCount)
				continue
			}
			break
		}
		if failCount >= retryTimes {
			continue
		}
		blog.Infof("push host identifier success: file: %v", fileList)

		// 将任务放到redis维护的任务list中
		if err := h.createNewTaskToTaskList(pushFileResp.MContent, hostInfos); err != nil {
			blog.Errorf("add task to redis task list error, taskID: %s, hostInfos: %v, err: %v",
				pushFileResp.MContent, hostInfos, err)
		}
	}
}

// collectFailHost collect fail host
func (h *HostIdentifier) collectFailHost() ([]*HostInfo, *get_agent_state_forsyncdata.AgentStatusRequest) {
	start := time.Now()
	var hostInfoArray []*HostInfo
	agentStatusRequest := new(get_agent_state_forsyncdata.AgentStatusRequest)
	uniqueMap := make(map[string]bool)
	// 从redis的保存失败的主机的list中拿出一定数量主机，并进行去重
	for i := 0; i < makeNewTaskFromFailHostCount; i++ {
		if time.Now().Sub(start) > 5*time.Minute {
			break
		}
		hostInfoMessage, err := h.getFailHostFromFailHostList()
		if err != nil {
			blog.Errorf("get host from redis fail_host_list error, err: %v", err)
			continue
		}
		hostInfo := new(HostInfo)
		if err := json.Unmarshal([]byte(hostInfoMessage), hostInfo); err != nil {
			blog.Errorf("Unmarshal hostInfo error, hostInfo: %s, err: %v", hostInfoMessage, err)
			continue
		}
		cloudIDStr := strconv.FormatInt(hostInfo.CloudID, 10)
		uniqueKey := cloudIDStr + ":" + hostInfo.HostInnerIP
		if uniqueMap[uniqueKey] {
			continue
		}
		uniqueMap[uniqueKey] = true
		agentStatusRequest.Hosts = append(agentStatusRequest.Hosts, &get_agent_state_forsyncdata.CacheIPInfo{
			GseCompositeID: cloudIDStr,
			IP:             hostInfo.HostInnerIP,
		})
		hostInfoArray = append(hostInfoArray, hostInfo)
	}
	return hostInfoArray, agentStatusRequest
}

// CreateNewTaskToTaskList create a new task to redis task list
func (h *HostIdentifier) createNewTaskToTaskList(taskID string, hostInfos []*HostInfo) error {
	var err error
	task := &Task{
		TaskID:      taskID,
		HostInfos:   hostInfos,
		ExpiredTime: time.Now().Add(24 * time.Hour).Unix(),
	}
	failCount := 0
	for failCount < retryTimes {
		if err := h.addTaskToTaskList(task); err != nil {
			failCount++
			sleepForFail(failCount)
			continue
		}
		break
	}
	return err
}

// addTaskToTaskList add task to redis task list
func (h *HostIdentifier) addTaskToTaskList(task *Task) error {
	return h.redisCli.RPush(h.ctx, redisTaskListName, task).Err()
}

// GetTaskFromTaskList get task from redis task list
func (h *HostIdentifier) getTaskFromTaskList() (string, error) {
	result, err := h.redisCli.BLPop(h.ctx, 0, redisTaskListName).Result()
	if err != nil {
		return "", err
	}
	return result[1], nil
}

// addHostToFailHostList add host to redis fail host list
func (h *HostIdentifier) addHostToFailHostList(host *HostInfo) error {
	return h.redisCli.RPush(h.ctx, RedisFailHostListName, host).Err()
}

// getFailHostFromFailHostList get fail host from redis fail host list
func (h *HostIdentifier) getFailHostFromFailHostList() (string, error) {
	result, err := h.redisCli.BLPop(h.ctx, 5*time.Minute, RedisFailHostListName).Result()
	if err != nil {
		return "", err
	}
	return result[1], nil
}
