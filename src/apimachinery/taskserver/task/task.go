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

// Package task TODO
package task

import (
	"context"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// Create 新加任务，taskType: 任务标识，留给业务方做识别任务，instID: 任务的执行源实例id，data: 每一项任务需要的参数
func (t *task) Create(ctx context.Context, header http.Header, taskType string, instID int64, data []interface{}) (
	metadata.APITaskDetail, errors.CCErrorCoder) {

	resp := new(metadata.CreateTaskResponse)
	subPath := "/task/create"
	body := metadata.CreateTaskRequest{
		TaskType: taskType,
		InstID:   instID,
		Data:     data,
	}

	err := t.client.Post().
		WithContext(ctx).
		Body(body).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	if err != nil {
		return metadata.APITaskDetail{}, errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return metadata.APITaskDetail{}, resp.CCError()
	}
	return resp.Data, nil
}

// CreateBatch create task batch, returns the created task details
func (t *task) CreateBatch(ctx context.Context, header http.Header, tasks []metadata.CreateTaskRequest) (
	[]metadata.APITaskDetail, error) {

	resp := new(metadata.CreateTaskBatchResponse)
	subPath := "/createmany/task"

	err := t.client.Post().
		WithContext(ctx).
		Body(tasks).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return nil, resp.CCError()
	}
	return resp.Data, nil
}

// ListTask TODO
func (t *task) ListTask(ctx context.Context, header http.Header, name string, data *metadata.ListAPITaskRequest) (resp *metadata.ListAPITaskResponse, err error) {
	resp = new(metadata.ListAPITaskResponse)
	subPath := "/task/findmany/list/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, name).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// ListLatestTask list the latest task for each inst by bk_inst_id
func (t *task) ListLatestTask(ctx context.Context, header http.Header, name string,
	data *metadata.ListAPITaskLatestRequest) ([]metadata.APITaskDetail, errors.CCErrorCoder) {

	resp := new(metadata.ListAPITaskLatestResponse)
	subPath := "/task/findmany/list/latest/%s"

	err := t.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, name).
		WithHeaders(header).
		Do().
		Into(resp)
	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return nil, resp.CCError()
	}
	return resp.Data, nil
}

// TaskDetail TODO
func (t *task) TaskDetail(ctx context.Context, header http.Header, taskID string) (resp *metadata.TaskDetailResponse, err error) {
	resp = new(metadata.TaskDetailResponse)
	subPath := "/task/findone/detail/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, taskID).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// DeleteTask delete task
func (t *task) DeleteTask(ctx context.Context, header http.Header, taskCond *metadata.DeleteOption) error {
	resp := new(metadata.Response)
	subPath := "/task/deletemany"

	err := t.client.Post().
		WithContext(ctx).
		Body(taskCond).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		blog.Errorf("delete task failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if resp.CCError() != nil {
		return resp.CCError()
	}

	return nil
}

// ListLatestSyncStatus list latest sync status by condition
func (t *task) ListLatestSyncStatus(ctx context.Context, header http.Header,
	option *metadata.ListLatestSyncStatusRequest) ([]metadata.APITaskSyncStatus, errors.CCErrorCoder) {

	resp := new(metadata.ListLatestSyncStatusResponse)
	subPath := "/findmany/latest/sync_status"

	err := t.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(&resp)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return nil, resp.CCError()
	}

	return resp.Data, nil
}

// ListSyncStatusHistory list sync status history by condition
func (t *task) ListSyncStatusHistory(ctx context.Context, header http.Header,
	option *metadata.QueryCondition) (*metadata.ListAPITaskSyncStatusResult, errors.CCErrorCoder) {

	resp := new(metadata.ListSyncStatusHistoryResponse)
	subPath := "/findmany/sync_status_history"

	err := t.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(&resp)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return nil, resp.CCError()
	}

	return resp.Data, nil
}
