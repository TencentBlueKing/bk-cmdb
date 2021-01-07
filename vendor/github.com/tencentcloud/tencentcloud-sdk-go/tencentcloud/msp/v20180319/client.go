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

package v20180319

import (
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2018-03-19"

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


func NewDeregisterMigrationTaskRequest() (request *DeregisterMigrationTaskRequest) {
    request = &DeregisterMigrationTaskRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("msp", APIVersion, "DeregisterMigrationTask")
    return
}

func NewDeregisterMigrationTaskResponse() (response *DeregisterMigrationTaskResponse) {
    response = &DeregisterMigrationTaskResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 取消注册迁移任务
func (c *Client) DeregisterMigrationTask(request *DeregisterMigrationTaskRequest) (response *DeregisterMigrationTaskResponse, err error) {
    if request == nil {
        request = NewDeregisterMigrationTaskRequest()
    }
    response = NewDeregisterMigrationTaskResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeMigrationTaskRequest() (request *DescribeMigrationTaskRequest) {
    request = &DescribeMigrationTaskRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("msp", APIVersion, "DescribeMigrationTask")
    return
}

func NewDescribeMigrationTaskResponse() (response *DescribeMigrationTaskResponse) {
    response = &DescribeMigrationTaskResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取指定迁移任务详情
func (c *Client) DescribeMigrationTask(request *DescribeMigrationTaskRequest) (response *DescribeMigrationTaskResponse, err error) {
    if request == nil {
        request = NewDescribeMigrationTaskRequest()
    }
    response = NewDescribeMigrationTaskResponse()
    err = c.Send(request, response)
    return
}

func NewListMigrationProjectRequest() (request *ListMigrationProjectRequest) {
    request = &ListMigrationProjectRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("msp", APIVersion, "ListMigrationProject")
    return
}

func NewListMigrationProjectResponse() (response *ListMigrationProjectResponse) {
    response = &ListMigrationProjectResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取迁移项目名称列表
func (c *Client) ListMigrationProject(request *ListMigrationProjectRequest) (response *ListMigrationProjectResponse, err error) {
    if request == nil {
        request = NewListMigrationProjectRequest()
    }
    response = NewListMigrationProjectResponse()
    err = c.Send(request, response)
    return
}

func NewListMigrationTaskRequest() (request *ListMigrationTaskRequest) {
    request = &ListMigrationTaskRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("msp", APIVersion, "ListMigrationTask")
    return
}

func NewListMigrationTaskResponse() (response *ListMigrationTaskResponse) {
    response = &ListMigrationTaskResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取迁移任务列表
func (c *Client) ListMigrationTask(request *ListMigrationTaskRequest) (response *ListMigrationTaskResponse, err error) {
    if request == nil {
        request = NewListMigrationTaskRequest()
    }
    response = NewListMigrationTaskResponse()
    err = c.Send(request, response)
    return
}

func NewModifyMigrationTaskBelongToProjectRequest() (request *ModifyMigrationTaskBelongToProjectRequest) {
    request = &ModifyMigrationTaskBelongToProjectRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("msp", APIVersion, "ModifyMigrationTaskBelongToProject")
    return
}

func NewModifyMigrationTaskBelongToProjectResponse() (response *ModifyMigrationTaskBelongToProjectResponse) {
    response = &ModifyMigrationTaskBelongToProjectResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 更改迁移任务所属项目
func (c *Client) ModifyMigrationTaskBelongToProject(request *ModifyMigrationTaskBelongToProjectRequest) (response *ModifyMigrationTaskBelongToProjectResponse, err error) {
    if request == nil {
        request = NewModifyMigrationTaskBelongToProjectRequest()
    }
    response = NewModifyMigrationTaskBelongToProjectResponse()
    err = c.Send(request, response)
    return
}

func NewModifyMigrationTaskStatusRequest() (request *ModifyMigrationTaskStatusRequest) {
    request = &ModifyMigrationTaskStatusRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("msp", APIVersion, "ModifyMigrationTaskStatus")
    return
}

func NewModifyMigrationTaskStatusResponse() (response *ModifyMigrationTaskStatusResponse) {
    response = &ModifyMigrationTaskStatusResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 更新迁移任务状态
func (c *Client) ModifyMigrationTaskStatus(request *ModifyMigrationTaskStatusRequest) (response *ModifyMigrationTaskStatusResponse, err error) {
    if request == nil {
        request = NewModifyMigrationTaskStatusRequest()
    }
    response = NewModifyMigrationTaskStatusResponse()
    err = c.Send(request, response)
    return
}

func NewRegisterMigrationTaskRequest() (request *RegisterMigrationTaskRequest) {
    request = &RegisterMigrationTaskRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("msp", APIVersion, "RegisterMigrationTask")
    return
}

func NewRegisterMigrationTaskResponse() (response *RegisterMigrationTaskResponse) {
    response = &RegisterMigrationTaskResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 注册迁移任务
func (c *Client) RegisterMigrationTask(request *RegisterMigrationTaskRequest) (response *RegisterMigrationTaskResponse, err error) {
    if request == nil {
        request = NewRegisterMigrationTaskRequest()
    }
    response = NewRegisterMigrationTaskResponse()
    err = c.Send(request, response)
    return
}
