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

package v20180408

import (
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2018-04-08"

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


func NewCreateContainerInstanceRequest() (request *CreateContainerInstanceRequest) {
    request = &CreateContainerInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cis", APIVersion, "CreateContainerInstance")
    return
}

func NewCreateContainerInstanceResponse() (response *CreateContainerInstanceResponse) {
    response = &CreateContainerInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 此接口（CreateContainerInstance）用于创建容器实例
func (c *Client) CreateContainerInstance(request *CreateContainerInstanceRequest) (response *CreateContainerInstanceResponse, err error) {
    if request == nil {
        request = NewCreateContainerInstanceRequest()
    }
    response = NewCreateContainerInstanceResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteContainerInstanceRequest() (request *DeleteContainerInstanceRequest) {
    request = &DeleteContainerInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cis", APIVersion, "DeleteContainerInstance")
    return
}

func NewDeleteContainerInstanceResponse() (response *DeleteContainerInstanceResponse) {
    response = &DeleteContainerInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 此接口（DeleteContainerInstance）用于删除容器实例
func (c *Client) DeleteContainerInstance(request *DeleteContainerInstanceRequest) (response *DeleteContainerInstanceResponse, err error) {
    if request == nil {
        request = NewDeleteContainerInstanceRequest()
    }
    response = NewDeleteContainerInstanceResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeContainerInstanceRequest() (request *DescribeContainerInstanceRequest) {
    request = &DescribeContainerInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cis", APIVersion, "DescribeContainerInstance")
    return
}

func NewDescribeContainerInstanceResponse() (response *DescribeContainerInstanceResponse) {
    response = &DescribeContainerInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 此接口（DescribeContainerInstance）用于获取容器实例详情
func (c *Client) DescribeContainerInstance(request *DescribeContainerInstanceRequest) (response *DescribeContainerInstanceResponse, err error) {
    if request == nil {
        request = NewDescribeContainerInstanceRequest()
    }
    response = NewDescribeContainerInstanceResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeContainerInstanceEventsRequest() (request *DescribeContainerInstanceEventsRequest) {
    request = &DescribeContainerInstanceEventsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cis", APIVersion, "DescribeContainerInstanceEvents")
    return
}

func NewDescribeContainerInstanceEventsResponse() (response *DescribeContainerInstanceEventsResponse) {
    response = &DescribeContainerInstanceEventsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 此接口（DescribeContainerInstanceEvents）用于查询容器实例事件列表
func (c *Client) DescribeContainerInstanceEvents(request *DescribeContainerInstanceEventsRequest) (response *DescribeContainerInstanceEventsResponse, err error) {
    if request == nil {
        request = NewDescribeContainerInstanceEventsRequest()
    }
    response = NewDescribeContainerInstanceEventsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeContainerInstancesRequest() (request *DescribeContainerInstancesRequest) {
    request = &DescribeContainerInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cis", APIVersion, "DescribeContainerInstances")
    return
}

func NewDescribeContainerInstancesResponse() (response *DescribeContainerInstancesResponse) {
    response = &DescribeContainerInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 此接口（DescribeContainerInstances）查询容器实例列表
func (c *Client) DescribeContainerInstances(request *DescribeContainerInstancesRequest) (response *DescribeContainerInstancesResponse, err error) {
    if request == nil {
        request = NewDescribeContainerInstancesRequest()
    }
    response = NewDescribeContainerInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeContainerLogRequest() (request *DescribeContainerLogRequest) {
    request = &DescribeContainerLogRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cis", APIVersion, "DescribeContainerLog")
    return
}

func NewDescribeContainerLogResponse() (response *DescribeContainerLogResponse) {
    response = &DescribeContainerLogResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 此接口（DescribeContainerLog）用于获取容器日志信息
func (c *Client) DescribeContainerLog(request *DescribeContainerLogRequest) (response *DescribeContainerLogResponse, err error) {
    if request == nil {
        request = NewDescribeContainerLogRequest()
    }
    response = NewDescribeContainerLogResponse()
    err = c.Send(request, response)
    return
}

func NewInquiryPriceCreateCisRequest() (request *InquiryPriceCreateCisRequest) {
    request = &InquiryPriceCreateCisRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cis", APIVersion, "InquiryPriceCreateCis")
    return
}

func NewInquiryPriceCreateCisResponse() (response *InquiryPriceCreateCisResponse) {
    response = &InquiryPriceCreateCisResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 此接口（InquiryPriceCreateCis）用于查询容器实例价格
func (c *Client) InquiryPriceCreateCis(request *InquiryPriceCreateCisRequest) (response *InquiryPriceCreateCisResponse, err error) {
    if request == nil {
        request = NewInquiryPriceCreateCisRequest()
    }
    response = NewInquiryPriceCreateCisResponse()
    err = c.Send(request, response)
    return
}
