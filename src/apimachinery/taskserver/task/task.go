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
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
)

// Create  新加任务， name 任务名，flag:任务标识，留给业务方做识别任务, data 每一项任务需要的参数
func (t *task) Create(ctx context.Context, header http.Header, name, flag string, data []interface{}) (resp *metadata.CreateTaskResponse, err error) {
	resp = new(metadata.CreateTaskResponse)
	subPath := "/task/create"
	body := metadata.CreateTaskRequest{
		Name: name,
		Flag: flag,
		Data: data,
	}

	err = t.client.Post().
		WithContext(ctx).
		Body(body).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (t *task) ListTask(ctx context.Context, header http.Header, name string, data *metadata.ListAPITaskRequest) (resp *metadata.ListAPITaskResponse, err error) {
	resp = new(metadata.ListAPITaskResponse)
	subPath := "/task/findmany/list/" + name

	err = t.client.Post().
		WithContext(ctx).
		Body(data).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (t *task) TaskDetail(ctx context.Context, header http.Header, taskID string) (resp *metadata.TaskDetailResponse, err error) {
	resp = new(metadata.TaskDetailResponse)
	subPath := "/task/findone/detail/" + taskID

	err = t.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (t *task) TaskStatusToSuccess(ctx context.Context, header http.Header, taskID, subTaskID string) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/task/set/status/sucess/id/%s/sub_id/%s", taskID, subTaskID)

	err = t.client.Put().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (t *task) TaskStatusToFailure(ctx context.Context, header http.Header, taskID, subTaskID string, errResponse *metadata.Response) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/task/set/status/failure/id/%s/sub_id/%s", taskID, subTaskID)

	err = t.client.Put().
		WithContext(ctx).
		Body(errResponse).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

/*


 http.MethodPost, Path: "/task/create", Handler: s.CreateTask})
 http.MethodPost, Path: "/task/findmany/list/{name}", Handler: s.ListTask})
 http.MethodPost, Path: "/task/findone/detail/{task_id}", Handler: s.DetailTask})
 http.MethodPut, Path: "/task/set/status/sucess/id/{task_id}/sub_id/{sub_task_id}", Handler: s.StatusToSuccess})
 http.MethodPut, Path: "/task/set/status/failure/id/{task_id}/sub_id/{sub_task_id}", Handler: s.StatusToFailure})

*/
