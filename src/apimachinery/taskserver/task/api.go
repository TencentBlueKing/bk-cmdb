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

package task

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// TaskClientInterface TODO
type TaskClientInterface interface {
	// Create  新加任务， name 任务名，flag:任务标识，留给业务方做识别任务, instID:任务的执行源实例id, data 每一项任务需要的参数
	Create(ctx context.Context, header http.Header, flag string, instID int64, data []interface{}) (
		metadata.APITaskDetail, errors.CCErrorCoder)

	CreateBatch(c context.Context, h http.Header, tasks []metadata.CreateTaskRequest) ([]metadata.APITaskDetail, error)
	CreateFieldTemplateBatch(c context.Context, h http.Header, tasks []metadata.CreateTaskRequest) (
		[]metadata.APITaskDetail, error)

	ListTask(ctx context.Context, header http.Header, name string, data *metadata.ListAPITaskRequest) (
		resp *metadata.ListAPITaskResponse, err error)

	ListLatestTask(ctx context.Context, header http.Header, name string, data *metadata.ListAPITaskLatestRequest) (
		[]metadata.APITaskDetail, errors.CCErrorCoder)

	TaskDetail(ctx context.Context, header http.Header, taskID string) (resp *metadata.TaskDetailResponse, err error)

	DeleteTask(ctx context.Context, header http.Header, taskCond *metadata.DeleteOption) error

	ListLatestSyncStatus(ctx context.Context, header http.Header, option *metadata.ListLatestSyncStatusRequest) (
		[]metadata.APITaskSyncStatus, errors.CCErrorCoder)

	ListSyncStatusHistory(ctx context.Context, header http.Header, option *metadata.QueryCondition) (
		*metadata.ListAPITaskSyncStatusResult, errors.CCErrorCoder)

	ListFieldTemplateTaskSyncResult(ctx context.Context, header http.Header,
		data *metadata.ListFieldTmplSyncTaskStatusOption) ([]metadata.ListFieldTmplTaskSyncResult, errors.CCErrorCoder)
}

// NewTaskClientInterface TODO
func NewTaskClientInterface(client rest.ClientInterface) TaskClientInterface {
	return &task{client: client}
}

type task struct {
	client rest.ClientInterface
}
