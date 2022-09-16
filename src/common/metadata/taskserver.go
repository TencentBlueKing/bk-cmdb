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
	// TaskType 任务标识，用于业务方识别任务，同时表示所在的任务队列
	TaskType string `json:"task_type"`

	// bk_inst_id 实例id，该任务关联的实例id
	InstID int64 `json:"bk_inst_id"`

	Data []interface{} `json:"data"`
}

// APITaskDetail task info detail
type APITaskDetail struct {
	// TaskID 任务ID，由taskserver生成的唯一ID
	TaskID string `json:"task_id,omitempty" bson:"task_id"`
	// TaskType 任务标识，用于业务方识别任务，同时表示所在的任务队列
	TaskType string `json:"task_type,omitempty" bson:"task_type"`
	// InstID 实例id，该任务关联的实例id
	InstID int64 `json:"bk_inst_id,omitempty" bson:"bk_inst_id"`
	// User 任务创建者
	User string `json:"user,omitempty" bson:"user"`
	// Header 请求的 http header
	Header http.Header `json:"header,omitempty" bson:"header"`
	// Status 任务执行状态
	Status APITaskStatus `json:"status,omitempty" bson:"status"`
	// Detail 子任务详情列表
	Detail []APISubTaskDetail `json:"detail,omitempty" bson:"detail"`

	// CreateTime 任务创建时间
	CreateTime time.Time `json:"create_time,omitempty" bson:"create_time"`
	// LastTime 任务最后更新时间
	LastTime time.Time `json:"last_time,omitempty" bson:"last_time"`
}

// APISubTaskDetail task data and execute detail
type APISubTaskDetail struct {
	SubTaskID string        `json:"sub_task_id,omitempty" bson:"sub_task_id"`
	Data      interface{}   `json:"data,omitempty" bson:"data"`
	Status    APITaskStatus `json:"status,omitempty" bson:"status"`
	Response  *Response     `json:"response,omitempty" bson:"response"`
}

// APITaskSyncStatus api task sync status
type APITaskSyncStatus struct {
	// TaskID 任务ID，对应APITaskDetail的TaskID
	TaskID string `json:"task_id,omitempty" bson:"task_id"`
	// TaskType 任务标识，用于业务方识别任务
	TaskType string `json:"task_type,omitempty" bson:"task_type"`
	// InstID 实例id，该任务关联的实例id
	InstID int64 `json:"bk_inst_id,omitempty" bson:"bk_inst_id"`
	// Status 任务执行状态
	Status APITaskStatus `json:"status,omitempty" bson:"status"`
	// Creator 任务创建者
	Creator string `json:"creator,omitempty" bson:"creator"`
	// CreateTime 任务创建时间
	CreateTime time.Time `json:"create_time,omitempty" bson:"create_time"`
	// LastTime 任务最后更新时间
	LastTime time.Time `json:"last_time,omitempty" bson:"last_time"`
	// SupplierAccount 开发商ID
	SupplierAccount string `json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
}

// APITaskStatus task status type
type APITaskStatus string

// IsFinished TODO
func (s APITaskStatus) IsFinished() bool {
	if s == APITaskStatusSuccess || s == APITAskStatusFail {
		return true
	}
	return false
}

// IsSuccessful TODO
func (s APITaskStatus) IsSuccessful() bool {
	if s == APITaskStatusSuccess {
		return true
	}
	return false
}

// IsFailure TODO
func (s APITaskStatus) IsFailure() bool {
	if s == APITAskStatusFail {
		return true
	}
	return false
}

const (
	// APITaskStatusNew new task ,waiting execute
	APITaskStatusNew APITaskStatus = "new"
	// APITaskStatusWaitExecute 正在执行的任务中断了。 补偿后。确定需要重新执行
	APITaskStatusWaitExecute APITaskStatus = "waiting"

	// APITaskStatusExecute task executing
	APITaskStatusExecute APITaskStatus = "executing"

	// APITaskStatusSuccess task execute success
	APITaskStatusSuccess APITaskStatus = "finished"

	// APITAskStatusFail task execute failure
	APITAskStatusFail APITaskStatus = "failure"

	// APITAskStatusNeedSync only used for instance with all tasks finished but actual status is not finished
	APITAskStatusNeedSync APITaskStatus = "need_sync"
)

// ListAPITaskRequest TODO
type ListAPITaskRequest struct {
	Condition mapstr.MapStr `json:"condition"`
	Page      BasePage      `json:"page"`
}

// ListAPITaskLatestRequest TODO
type ListAPITaskLatestRequest struct {
	Condition mapstr.MapStr `json:"condition"`
	Fields    []string      `json:"fields"`
}

// ListAPITaskLatestResponse TODO
type ListAPITaskLatestResponse struct {
	BaseResp
	Data []APITaskDetail `json:"data"`
}

// ListAPITaskData TODO
type ListAPITaskData struct {
	Info  []APITaskDetail `json:"info"`
	Count int64           `json:"count"`
	Page  BasePage        `json:"page"`
}

// ListAPITaskResponse TODO
type ListAPITaskResponse struct {
	BaseResp
	Data ListAPITaskData `json:"data"`
}

// CreateTaskResponse TODO
type CreateTaskResponse struct {
	BaseResp
	Data APITaskDetail `json:"data"`
}

// CreateTaskBatchResponse batch create task response
type CreateTaskBatchResponse struct {
	BaseResp
	Data []APITaskDetail `json:"data"`
}

// TaskDetailResponse api task detail response
type TaskDetailResponse struct {
	BaseResp
	Data struct {
		Info APITaskDetail `json:"info"`
	} `json:"data"`
}

// ListAPITaskDetail list api task detail condition
type ListAPITaskDetail struct {
	InstID []int64  `json:"bk_inst_id"`
	Fields []string `json:"fields"`
}

// ListLatestSyncStatusRequest list latest api task sync status request
type ListLatestSyncStatusRequest struct {
	Condition mapstr.MapStr `json:"condition"`
	Fields    []string      `json:"fields"`
	// 非必填，只能用来查时间，且与Condition是与关系
	TimeCondition *TimeCondition `json:"time_condition,omitempty"`
}

// ListLatestSyncStatusResponse list latest api task sync status response
type ListLatestSyncStatusResponse struct {
	BaseResp
	Data []APITaskSyncStatus `json:"data"`
}

// ListSyncStatusHistoryResponse list api task sync history response
type ListSyncStatusHistoryResponse struct {
	BaseResp
	Data *ListAPITaskSyncStatusResult `json:"data"`
}

// ListAPITaskSyncStatusResult list api task sync status paged result
type ListAPITaskSyncStatusResult struct {
	Count int64               `json:"count"`
	Info  []APITaskSyncStatus `json:"info"`
}
