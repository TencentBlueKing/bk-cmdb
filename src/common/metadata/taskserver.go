/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */
package metadata

import (
	"net/http"
	"time"

	"configcenter/src/common/mapstr"
)

// CreateTaskRequest create task request parameters
type CreateTaskRequest struct {
	// task name
	Name string `json:"name"`

	Header http.Header `json:"header"`

	Data []interface{} `json:"data"`
}

// APITaskDetail task info detaill
type APITaskDetail struct {
	// task id
	TaskID string `json:"task_id"`
	// task name. 表示所在的任务队列
	Name string `json:"name"`
	// task create user
	User string `json:"user"`
	//  http header
	Header http.Header `json:"header"`
	// task status
	Status APITaskStatus `json:"status"`
	// sub task detail
	Detail []APISubTaskDetail `json:"detail"`

	CreateTime time.Time `json:"create_time"`
	LastTime   time.Time `json:"last_time"`
}

// APISubTaskDetail task data and execute detail
type APISubTaskDetail struct {
	SubTaskID string        `json:"sub_task_id"`
	Data      interface{}   `json:"data"`
	Status    APITaskStatus `json:"status"`
	Response  *Response     `json:"response"`
}

// APITaskStatus task status type
type APITaskStatus int64

const (
	// APITaskStatusNew new task ,waiting execute
	APITaskStatusNew APITaskStatus = 0

	// APITaskStatuExecute task executing
	APITaskStatuExecute APITaskStatus = 100

	// APITaskStatusSuccess task execute success
	APITaskStatusSuccess APITaskStatus = 200

	// APITAskStatusFail task execute failure
	APITAskStatusFail APITaskStatus = 500
)

type ListAPITaskRequest struct {
	Condition mapstr.MapStr `json:"condition"`
	Page      BasePage      `json:"page"`
}

type ListAPITaskData struct {
	Info  []APITaskDetail `json:"data"`
	Count int64           `json:"count"`
	Page  BasePage        `json:"page"`
}

type ListAPITaskResponse struct {
	BaseResp
	Data ListAPITaskData `json:"data"`
}

type CreateTaskResponse struct {
	BaseResp
	Data struct {
		Info APITaskDetail `json:"info"`
	} `json:"data"`
}

type TaskDetailResponse struct {
	BaseResp
	Data struct {
		Info APITaskDetail `json:"info"`
	} `json:"data"`
}
