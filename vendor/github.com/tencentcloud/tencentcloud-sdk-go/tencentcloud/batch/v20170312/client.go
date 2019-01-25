// Copyright (c) 2017-2018 THL A29 Limited, a Tencent company. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v20170312

import (
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2017-03-12"

type Client struct {
    common.Client
}

// Deprecated
func NewClientWithSecretId(secretId, secretKey, region string) (client *Client, err error) {
    cpf := profile.NewClientProfile()
    client = &Client{}
    client.Init(region).WithSecretId(secretId, secretKey).WithProfile(cpf)
    return
}

func NewClient(credential *common.Credential, region string, clientProfile *profile.ClientProfile) (client *Client, err error) {
    client = &Client{}
    client.Init(region).
        WithCredential(credential).
        WithProfile(clientProfile)
    return
}


func NewCreateComputeEnvRequest() (request *CreateComputeEnvRequest) {
    request = &CreateComputeEnvRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "CreateComputeEnv")
    return
}

func NewCreateComputeEnvResponse() (response *CreateComputeEnvResponse) {
    response = &CreateComputeEnvResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于创建计算环境
func (c *Client) CreateComputeEnv(request *CreateComputeEnvRequest) (response *CreateComputeEnvResponse, err error) {
    if request == nil {
        request = NewCreateComputeEnvRequest()
    }
    response = NewCreateComputeEnvResponse()
    err = c.Send(request, response)
    return
}

func NewCreateTaskTemplateRequest() (request *CreateTaskTemplateRequest) {
    request = &CreateTaskTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "CreateTaskTemplate")
    return
}

func NewCreateTaskTemplateResponse() (response *CreateTaskTemplateResponse) {
    response = &CreateTaskTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于创建任务模板
func (c *Client) CreateTaskTemplate(request *CreateTaskTemplateRequest) (response *CreateTaskTemplateResponse, err error) {
    if request == nil {
        request = NewCreateTaskTemplateRequest()
    }
    response = NewCreateTaskTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteComputeEnvRequest() (request *DeleteComputeEnvRequest) {
    request = &DeleteComputeEnvRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DeleteComputeEnv")
    return
}

func NewDeleteComputeEnvResponse() (response *DeleteComputeEnvResponse) {
    response = &DeleteComputeEnvResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于删除计算环境
func (c *Client) DeleteComputeEnv(request *DeleteComputeEnvRequest) (response *DeleteComputeEnvResponse, err error) {
    if request == nil {
        request = NewDeleteComputeEnvRequest()
    }
    response = NewDeleteComputeEnvResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteJobRequest() (request *DeleteJobRequest) {
    request = &DeleteJobRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DeleteJob")
    return
}

func NewDeleteJobResponse() (response *DeleteJobResponse) {
    response = &DeleteJobResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于删除作业记录。
// 删除作业的效果相当于删除作业相关的所有信息。删除成功后，作业相关的所有信息都无法查询。
// 待删除的作业必须处于完结状态，且其内部包含的所有任务实例也必须处于完结状态，否则会禁止操作。完结状态，是指处于 SUCCEED 或 FAILED 状态。
func (c *Client) DeleteJob(request *DeleteJobRequest) (response *DeleteJobResponse, err error) {
    if request == nil {
        request = NewDeleteJobRequest()
    }
    response = NewDeleteJobResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteTaskTemplatesRequest() (request *DeleteTaskTemplatesRequest) {
    request = &DeleteTaskTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DeleteTaskTemplates")
    return
}

func NewDeleteTaskTemplatesResponse() (response *DeleteTaskTemplatesResponse) {
    response = &DeleteTaskTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于删除任务模板信息
func (c *Client) DeleteTaskTemplates(request *DeleteTaskTemplatesRequest) (response *DeleteTaskTemplatesResponse, err error) {
    if request == nil {
        request = NewDeleteTaskTemplatesRequest()
    }
    response = NewDeleteTaskTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeAvailableCvmInstanceTypesRequest() (request *DescribeAvailableCvmInstanceTypesRequest) {
    request = &DescribeAvailableCvmInstanceTypesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DescribeAvailableCvmInstanceTypes")
    return
}

func NewDescribeAvailableCvmInstanceTypesResponse() (response *DescribeAvailableCvmInstanceTypesResponse) {
    response = &DescribeAvailableCvmInstanceTypesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 查看可用的CVM机型配置信息
func (c *Client) DescribeAvailableCvmInstanceTypes(request *DescribeAvailableCvmInstanceTypesRequest) (response *DescribeAvailableCvmInstanceTypesResponse, err error) {
    if request == nil {
        request = NewDescribeAvailableCvmInstanceTypesRequest()
    }
    response = NewDescribeAvailableCvmInstanceTypesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeComputeEnvRequest() (request *DescribeComputeEnvRequest) {
    request = &DescribeComputeEnvRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DescribeComputeEnv")
    return
}

func NewDescribeComputeEnvResponse() (response *DescribeComputeEnvResponse) {
    response = &DescribeComputeEnvResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于查询计算环境的详细信息
func (c *Client) DescribeComputeEnv(request *DescribeComputeEnvRequest) (response *DescribeComputeEnvResponse, err error) {
    if request == nil {
        request = NewDescribeComputeEnvRequest()
    }
    response = NewDescribeComputeEnvResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeComputeEnvActivitiesRequest() (request *DescribeComputeEnvActivitiesRequest) {
    request = &DescribeComputeEnvActivitiesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DescribeComputeEnvActivities")
    return
}

func NewDescribeComputeEnvActivitiesResponse() (response *DescribeComputeEnvActivitiesResponse) {
    response = &DescribeComputeEnvActivitiesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于查询计算环境的活动信息
func (c *Client) DescribeComputeEnvActivities(request *DescribeComputeEnvActivitiesRequest) (response *DescribeComputeEnvActivitiesResponse, err error) {
    if request == nil {
        request = NewDescribeComputeEnvActivitiesRequest()
    }
    response = NewDescribeComputeEnvActivitiesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeComputeEnvCreateInfoRequest() (request *DescribeComputeEnvCreateInfoRequest) {
    request = &DescribeComputeEnvCreateInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DescribeComputeEnvCreateInfo")
    return
}

func NewDescribeComputeEnvCreateInfoResponse() (response *DescribeComputeEnvCreateInfoResponse) {
    response = &DescribeComputeEnvCreateInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 查看计算环境的创建信息。
func (c *Client) DescribeComputeEnvCreateInfo(request *DescribeComputeEnvCreateInfoRequest) (response *DescribeComputeEnvCreateInfoResponse, err error) {
    if request == nil {
        request = NewDescribeComputeEnvCreateInfoRequest()
    }
    response = NewDescribeComputeEnvCreateInfoResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeComputeEnvCreateInfosRequest() (request *DescribeComputeEnvCreateInfosRequest) {
    request = &DescribeComputeEnvCreateInfosRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DescribeComputeEnvCreateInfos")
    return
}

func NewDescribeComputeEnvCreateInfosResponse() (response *DescribeComputeEnvCreateInfosResponse) {
    response = &DescribeComputeEnvCreateInfosResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于查看计算环境创建信息列表，包括名称、描述、类型、环境参数、通知及期望节点数等。
func (c *Client) DescribeComputeEnvCreateInfos(request *DescribeComputeEnvCreateInfosRequest) (response *DescribeComputeEnvCreateInfosResponse, err error) {
    if request == nil {
        request = NewDescribeComputeEnvCreateInfosRequest()
    }
    response = NewDescribeComputeEnvCreateInfosResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeComputeEnvsRequest() (request *DescribeComputeEnvsRequest) {
    request = &DescribeComputeEnvsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DescribeComputeEnvs")
    return
}

func NewDescribeComputeEnvsResponse() (response *DescribeComputeEnvsResponse) {
    response = &DescribeComputeEnvsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于查看计算环境列表
func (c *Client) DescribeComputeEnvs(request *DescribeComputeEnvsRequest) (response *DescribeComputeEnvsResponse, err error) {
    if request == nil {
        request = NewDescribeComputeEnvsRequest()
    }
    response = NewDescribeComputeEnvsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeCvmZoneInstanceConfigInfosRequest() (request *DescribeCvmZoneInstanceConfigInfosRequest) {
    request = &DescribeCvmZoneInstanceConfigInfosRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DescribeCvmZoneInstanceConfigInfos")
    return
}

func NewDescribeCvmZoneInstanceConfigInfosResponse() (response *DescribeCvmZoneInstanceConfigInfosResponse) {
    response = &DescribeCvmZoneInstanceConfigInfosResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取批量计算可用区机型配置信息
func (c *Client) DescribeCvmZoneInstanceConfigInfos(request *DescribeCvmZoneInstanceConfigInfosRequest) (response *DescribeCvmZoneInstanceConfigInfosResponse, err error) {
    if request == nil {
        request = NewDescribeCvmZoneInstanceConfigInfosRequest()
    }
    response = NewDescribeCvmZoneInstanceConfigInfosResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeJobRequest() (request *DescribeJobRequest) {
    request = &DescribeJobRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DescribeJob")
    return
}

func NewDescribeJobResponse() (response *DescribeJobResponse) {
    response = &DescribeJobResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于查看一个作业的详细信息，包括内部任务（Task）和依赖（Dependence）信息。
func (c *Client) DescribeJob(request *DescribeJobRequest) (response *DescribeJobResponse, err error) {
    if request == nil {
        request = NewDescribeJobRequest()
    }
    response = NewDescribeJobResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeJobSubmitInfoRequest() (request *DescribeJobSubmitInfoRequest) {
    request = &DescribeJobSubmitInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DescribeJobSubmitInfo")
    return
}

func NewDescribeJobSubmitInfoResponse() (response *DescribeJobSubmitInfoResponse) {
    response = &DescribeJobSubmitInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于查询指定作业的提交信息，其返回内容包括 JobId 和 SubmitJob 接口中作为输入参数的作业提交信息
func (c *Client) DescribeJobSubmitInfo(request *DescribeJobSubmitInfoRequest) (response *DescribeJobSubmitInfoResponse, err error) {
    if request == nil {
        request = NewDescribeJobSubmitInfoRequest()
    }
    response = NewDescribeJobSubmitInfoResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeJobsRequest() (request *DescribeJobsRequest) {
    request = &DescribeJobsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DescribeJobs")
    return
}

func NewDescribeJobsResponse() (response *DescribeJobsResponse) {
    response = &DescribeJobsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于查询若干个作业的概览信息
func (c *Client) DescribeJobs(request *DescribeJobsRequest) (response *DescribeJobsResponse, err error) {
    if request == nil {
        request = NewDescribeJobsRequest()
    }
    response = NewDescribeJobsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTaskRequest() (request *DescribeTaskRequest) {
    request = &DescribeTaskRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DescribeTask")
    return
}

func NewDescribeTaskResponse() (response *DescribeTaskResponse) {
    response = &DescribeTaskResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于查询指定任务的详细信息，包括任务内部的任务实例信息。
func (c *Client) DescribeTask(request *DescribeTaskRequest) (response *DescribeTaskResponse, err error) {
    if request == nil {
        request = NewDescribeTaskRequest()
    }
    response = NewDescribeTaskResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTaskLogsRequest() (request *DescribeTaskLogsRequest) {
    request = &DescribeTaskLogsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DescribeTaskLogs")
    return
}

func NewDescribeTaskLogsResponse() (response *DescribeTaskLogsResponse) {
    response = &DescribeTaskLogsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于获取任务多个实例标准输出和标准错误日志。
func (c *Client) DescribeTaskLogs(request *DescribeTaskLogsRequest) (response *DescribeTaskLogsResponse, err error) {
    if request == nil {
        request = NewDescribeTaskLogsRequest()
    }
    response = NewDescribeTaskLogsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTaskTemplatesRequest() (request *DescribeTaskTemplatesRequest) {
    request = &DescribeTaskTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "DescribeTaskTemplates")
    return
}

func NewDescribeTaskTemplatesResponse() (response *DescribeTaskTemplatesResponse) {
    response = &DescribeTaskTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于查询任务模板信息
func (c *Client) DescribeTaskTemplates(request *DescribeTaskTemplatesRequest) (response *DescribeTaskTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeTaskTemplatesRequest()
    }
    response = NewDescribeTaskTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewModifyComputeEnvRequest() (request *ModifyComputeEnvRequest) {
    request = &ModifyComputeEnvRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "ModifyComputeEnv")
    return
}

func NewModifyComputeEnvResponse() (response *ModifyComputeEnvResponse) {
    response = &ModifyComputeEnvResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于修改计算环境属性
func (c *Client) ModifyComputeEnv(request *ModifyComputeEnvRequest) (response *ModifyComputeEnvResponse, err error) {
    if request == nil {
        request = NewModifyComputeEnvRequest()
    }
    response = NewModifyComputeEnvResponse()
    err = c.Send(request, response)
    return
}

func NewModifyTaskTemplateRequest() (request *ModifyTaskTemplateRequest) {
    request = &ModifyTaskTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "ModifyTaskTemplate")
    return
}

func NewModifyTaskTemplateResponse() (response *ModifyTaskTemplateResponse) {
    response = &ModifyTaskTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于修改任务模板
func (c *Client) ModifyTaskTemplate(request *ModifyTaskTemplateRequest) (response *ModifyTaskTemplateResponse, err error) {
    if request == nil {
        request = NewModifyTaskTemplateRequest()
    }
    response = NewModifyTaskTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewSubmitJobRequest() (request *SubmitJobRequest) {
    request = &SubmitJobRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "SubmitJob")
    return
}

func NewSubmitJobResponse() (response *SubmitJobResponse) {
    response = &SubmitJobResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于提交一个作业
func (c *Client) SubmitJob(request *SubmitJobRequest) (response *SubmitJobResponse, err error) {
    if request == nil {
        request = NewSubmitJobRequest()
    }
    response = NewSubmitJobResponse()
    err = c.Send(request, response)
    return
}

func NewTerminateComputeNodeRequest() (request *TerminateComputeNodeRequest) {
    request = &TerminateComputeNodeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "TerminateComputeNode")
    return
}

func NewTerminateComputeNodeResponse() (response *TerminateComputeNodeResponse) {
    response = &TerminateComputeNodeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于销毁计算节点。
// 对于状态为CREATED、CREATION_FAILED、RUNNING和ABNORMAL的节点，允许销毁处理。
func (c *Client) TerminateComputeNode(request *TerminateComputeNodeRequest) (response *TerminateComputeNodeResponse, err error) {
    if request == nil {
        request = NewTerminateComputeNodeRequest()
    }
    response = NewTerminateComputeNodeResponse()
    err = c.Send(request, response)
    return
}

func NewTerminateComputeNodesRequest() (request *TerminateComputeNodesRequest) {
    request = &TerminateComputeNodesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "TerminateComputeNodes")
    return
}

func NewTerminateComputeNodesResponse() (response *TerminateComputeNodesResponse) {
    response = &TerminateComputeNodesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于批量销毁计算节点，不允许重复销毁同一个节点。
func (c *Client) TerminateComputeNodes(request *TerminateComputeNodesRequest) (response *TerminateComputeNodesResponse, err error) {
    if request == nil {
        request = NewTerminateComputeNodesRequest()
    }
    response = NewTerminateComputeNodesResponse()
    err = c.Send(request, response)
    return
}

func NewTerminateJobRequest() (request *TerminateJobRequest) {
    request = &TerminateJobRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "TerminateJob")
    return
}

func NewTerminateJobResponse() (response *TerminateJobResponse) {
    response = &TerminateJobResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于终止作业。
// 当作业处于“SUBMITTED”状态时，禁止终止操作；当作业处于“SUCCEED”状态时，终止操作不会生效。
// 终止作业是一个异步过程。整个终止过程的耗时和任务总数成正比。终止的效果相当于所含的所有任务实例进行TerminateTaskInstance操作。具体效果和用法可参考TerminateTaskInstance。
func (c *Client) TerminateJob(request *TerminateJobRequest) (response *TerminateJobResponse, err error) {
    if request == nil {
        request = NewTerminateJobRequest()
    }
    response = NewTerminateJobResponse()
    err = c.Send(request, response)
    return
}

func NewTerminateTaskInstanceRequest() (request *TerminateTaskInstanceRequest) {
    request = &TerminateTaskInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("batch", APIVersion, "TerminateTaskInstance")
    return
}

func NewTerminateTaskInstanceResponse() (response *TerminateTaskInstanceResponse) {
    response = &TerminateTaskInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于终止任务实例。
// 对于状态已经为“SUCCEED”和“FAILED”的任务实例，不做处理。
// 对于状态为“SUBMITTED”、“PENDING”、“RUNNABLE”的任务实例，状态将置为“FAILED”状态。
// 对于状态为“STARTING”、“RUNNING”、“FAILED_INTERRUPTED”的任务实例，分区两种情况：如果未显示指定计算环境，会先销毁CVM服务器，然后将状态置为“FAILED”，具有一定耗时；如果指定了计算环境EnvId，任务实例状态置为“FAILED”，并重启执行该任务的CVM服务器，具有一定的耗时。
// 对于状态为“FAILED_INTERRUPTED”的任务实例，终止操作实际成功之后，相关资源和配额才会释放。
func (c *Client) TerminateTaskInstance(request *TerminateTaskInstanceRequest) (response *TerminateTaskInstanceResponse, err error) {
    if request == nil {
        request = NewTerminateTaskInstanceRequest()
    }
    response = NewTerminateTaskInstanceResponse()
    err = c.Send(request, response)
    return
}
