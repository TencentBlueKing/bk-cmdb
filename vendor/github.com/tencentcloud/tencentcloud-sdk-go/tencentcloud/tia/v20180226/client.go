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

package v20180226

import (
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2018-02-26"

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


func NewCreateJobRequest() (request *CreateJobRequest) {
    request = &CreateJobRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tia", APIVersion, "CreateJob")
    return
}

func NewCreateJobResponse() (response *CreateJobResponse) {
    response = &CreateJobResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 创建训练任务
func (c *Client) CreateJob(request *CreateJobRequest) (response *CreateJobResponse, err error) {
    if request == nil {
        request = NewCreateJobRequest()
    }
    response = NewCreateJobResponse()
    err = c.Send(request, response)
    return
}

func NewCreateModelRequest() (request *CreateModelRequest) {
    request = &CreateModelRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tia", APIVersion, "CreateModel")
    return
}

func NewCreateModelResponse() (response *CreateModelResponse) {
    response = &CreateModelResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 在指定的集群上部署一个模型，用以提供服务。
func (c *Client) CreateModel(request *CreateModelRequest) (response *CreateModelResponse, err error) {
    if request == nil {
        request = NewCreateModelRequest()
    }
    response = NewCreateModelResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteJobRequest() (request *DeleteJobRequest) {
    request = &DeleteJobRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tia", APIVersion, "DeleteJob")
    return
}

func NewDeleteJobResponse() (response *DeleteJobResponse) {
    response = &DeleteJobResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 删除训练任务
func (c *Client) DeleteJob(request *DeleteJobRequest) (response *DeleteJobResponse, err error) {
    if request == nil {
        request = NewDeleteJobRequest()
    }
    response = NewDeleteJobResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteModelRequest() (request *DeleteModelRequest) {
    request = &DeleteModelRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tia", APIVersion, "DeleteModel")
    return
}

func NewDeleteModelResponse() (response *DeleteModelResponse) {
    response = &DeleteModelResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 删除一个指定的Model
func (c *Client) DeleteModel(request *DeleteModelRequest) (response *DeleteModelResponse, err error) {
    if request == nil {
        request = NewDeleteModelRequest()
    }
    response = NewDeleteModelResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeJobRequest() (request *DescribeJobRequest) {
    request = &DescribeJobRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tia", APIVersion, "DescribeJob")
    return
}

func NewDescribeJobResponse() (response *DescribeJobResponse) {
    response = &DescribeJobResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取训练任务详情
func (c *Client) DescribeJob(request *DescribeJobRequest) (response *DescribeJobResponse, err error) {
    if request == nil {
        request = NewDescribeJobRequest()
    }
    response = NewDescribeJobResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeModelRequest() (request *DescribeModelRequest) {
    request = &DescribeModelRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tia", APIVersion, "DescribeModel")
    return
}

func NewDescribeModelResponse() (response *DescribeModelResponse) {
    response = &DescribeModelResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 描述Model
func (c *Client) DescribeModel(request *DescribeModelRequest) (response *DescribeModelResponse, err error) {
    if request == nil {
        request = NewDescribeModelRequest()
    }
    response = NewDescribeModelResponse()
    err = c.Send(request, response)
    return
}

func NewInstallAgentRequest() (request *InstallAgentRequest) {
    request = &InstallAgentRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tia", APIVersion, "InstallAgent")
    return
}

func NewInstallAgentResponse() (response *InstallAgentResponse) {
    response = &InstallAgentResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 安装agent
func (c *Client) InstallAgent(request *InstallAgentRequest) (response *InstallAgentResponse, err error) {
    if request == nil {
        request = NewInstallAgentRequest()
    }
    response = NewInstallAgentResponse()
    err = c.Send(request, response)
    return
}

func NewListJobsRequest() (request *ListJobsRequest) {
    request = &ListJobsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tia", APIVersion, "ListJobs")
    return
}

func NewListJobsResponse() (response *ListJobsResponse) {
    response = &ListJobsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 列举训练任务
func (c *Client) ListJobs(request *ListJobsRequest) (response *ListJobsResponse, err error) {
    if request == nil {
        request = NewListJobsRequest()
    }
    response = NewListJobsResponse()
    err = c.Send(request, response)
    return
}

func NewListModelsRequest() (request *ListModelsRequest) {
    request = &ListModelsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tia", APIVersion, "ListModels")
    return
}

func NewListModelsResponse() (response *ListModelsResponse) {
    response = &ListModelsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 列举某个指定集群上运行的模型。
func (c *Client) ListModels(request *ListModelsRequest) (response *ListModelsResponse, err error) {
    if request == nil {
        request = NewListModelsRequest()
    }
    response = NewListModelsResponse()
    err = c.Send(request, response)
    return
}

func NewQueryLogsRequest() (request *QueryLogsRequest) {
    request = &QueryLogsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tia", APIVersion, "QueryLogs")
    return
}

func NewQueryLogsResponse() (response *QueryLogsResponse) {
    response = &QueryLogsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 查询日志
func (c *Client) QueryLogs(request *QueryLogsRequest) (response *QueryLogsResponse, err error) {
    if request == nil {
        request = NewQueryLogsRequest()
    }
    response = NewQueryLogsResponse()
    err = c.Send(request, response)
    return
}
