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


func NewCreateBindInstanceRequest() (request *CreateBindInstanceRequest) {
    request = &CreateBindInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ms", APIVersion, "CreateBindInstance")
    return
}

func NewCreateBindInstanceResponse() (response *CreateBindInstanceResponse) {
    response = &CreateBindInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 将应用和资源进行绑定
func (c *Client) CreateBindInstance(request *CreateBindInstanceRequest) (response *CreateBindInstanceResponse, err error) {
    if request == nil {
        request = NewCreateBindInstanceRequest()
    }
    response = NewCreateBindInstanceResponse()
    err = c.Send(request, response)
    return
}

func NewCreateCosSecKeyInstanceRequest() (request *CreateCosSecKeyInstanceRequest) {
    request = &CreateCosSecKeyInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ms", APIVersion, "CreateCosSecKeyInstance")
    return
}

func NewCreateCosSecKeyInstanceResponse() (response *CreateCosSecKeyInstanceResponse) {
    response = &CreateCosSecKeyInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取云COS文件存储临时密钥，密钥仅限于临时上传文件，有访问限制和时效性。
func (c *Client) CreateCosSecKeyInstance(request *CreateCosSecKeyInstanceRequest) (response *CreateCosSecKeyInstanceResponse, err error) {
    if request == nil {
        request = NewCreateCosSecKeyInstanceRequest()
    }
    response = NewCreateCosSecKeyInstanceResponse()
    err = c.Send(request, response)
    return
}

func NewCreateResourceInstancesRequest() (request *CreateResourceInstancesRequest) {
    request = &CreateResourceInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ms", APIVersion, "CreateResourceInstances")
    return
}

func NewCreateResourceInstancesResponse() (response *CreateResourceInstancesResponse) {
    response = &CreateResourceInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用户可以使用该接口自建资源，只支持白名单用户
func (c *Client) CreateResourceInstances(request *CreateResourceInstancesRequest) (response *CreateResourceInstancesResponse, err error) {
    if request == nil {
        request = NewCreateResourceInstancesRequest()
    }
    response = NewCreateResourceInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewCreateScanInstancesRequest() (request *CreateScanInstancesRequest) {
    request = &CreateScanInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ms", APIVersion, "CreateScanInstances")
    return
}

func NewCreateScanInstancesResponse() (response *CreateScanInstancesResponse) {
    response = &CreateScanInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用户通过该接口批量提交应用进行应用扫描，扫描后需通过DescribeScanResults接口查询扫描结果
func (c *Client) CreateScanInstances(request *CreateScanInstancesRequest) (response *CreateScanInstancesResponse, err error) {
    if request == nil {
        request = NewCreateScanInstancesRequest()
    }
    response = NewCreateScanInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewCreateShieldInstanceRequest() (request *CreateShieldInstanceRequest) {
    request = &CreateShieldInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ms", APIVersion, "CreateShieldInstance")
    return
}

func NewCreateShieldInstanceResponse() (response *CreateShieldInstanceResponse) {
    response = &CreateShieldInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用户通过该接口提交应用进行应用加固，加固后需通过DescribeShieldResult接口查询加固结果
func (c *Client) CreateShieldInstance(request *CreateShieldInstanceRequest) (response *CreateShieldInstanceResponse, err error) {
    if request == nil {
        request = NewCreateShieldInstanceRequest()
    }
    response = NewCreateShieldInstanceResponse()
    err = c.Send(request, response)
    return
}

func NewCreateShieldPlanInstanceRequest() (request *CreateShieldPlanInstanceRequest) {
    request = &CreateShieldPlanInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ms", APIVersion, "CreateShieldPlanInstance")
    return
}

func NewCreateShieldPlanInstanceResponse() (response *CreateShieldPlanInstanceResponse) {
    response = &CreateShieldPlanInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 对资源进行策略新增
func (c *Client) CreateShieldPlanInstance(request *CreateShieldPlanInstanceRequest) (response *CreateShieldPlanInstanceResponse, err error) {
    if request == nil {
        request = NewCreateShieldPlanInstanceRequest()
    }
    response = NewCreateShieldPlanInstanceResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteScanInstancesRequest() (request *DeleteScanInstancesRequest) {
    request = &DeleteScanInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ms", APIVersion, "DeleteScanInstances")
    return
}

func NewDeleteScanInstancesResponse() (response *DeleteScanInstancesResponse) {
    response = &DeleteScanInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 删除一个或者多个app扫描信息
func (c *Client) DeleteScanInstances(request *DeleteScanInstancesRequest) (response *DeleteScanInstancesResponse, err error) {
    if request == nil {
        request = NewDeleteScanInstancesRequest()
    }
    response = NewDeleteScanInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteShieldInstancesRequest() (request *DeleteShieldInstancesRequest) {
    request = &DeleteShieldInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ms", APIVersion, "DeleteShieldInstances")
    return
}

func NewDeleteShieldInstancesResponse() (response *DeleteShieldInstancesResponse) {
    response = &DeleteShieldInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 删除一个或者多个app加固信息
func (c *Client) DeleteShieldInstances(request *DeleteShieldInstancesRequest) (response *DeleteShieldInstancesResponse, err error) {
    if request == nil {
        request = NewDeleteShieldInstancesRequest()
    }
    response = NewDeleteShieldInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeResourceInstancesRequest() (request *DescribeResourceInstancesRequest) {
    request = &DescribeResourceInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ms", APIVersion, "DescribeResourceInstances")
    return
}

func NewDescribeResourceInstancesResponse() (response *DescribeResourceInstancesResponse) {
    response = &DescribeResourceInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取某个用户的所有资源信息
func (c *Client) DescribeResourceInstances(request *DescribeResourceInstancesRequest) (response *DescribeResourceInstancesResponse, err error) {
    if request == nil {
        request = NewDescribeResourceInstancesRequest()
    }
    response = NewDescribeResourceInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeScanInstancesRequest() (request *DescribeScanInstancesRequest) {
    request = &DescribeScanInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ms", APIVersion, "DescribeScanInstances")
    return
}

func NewDescribeScanInstancesResponse() (response *DescribeScanInstancesResponse) {
    response = &DescribeScanInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口用于查看app列表。
// 可以通过指定任务唯一标识ItemId来查询指定app的详细信息，或通过设定过滤器来查询满足过滤条件的app的详细信息。 指定偏移(Offset)和限制(Limit)来选择结果中的一部分，默认返回满足条件的前20个app信息。
func (c *Client) DescribeScanInstances(request *DescribeScanInstancesRequest) (response *DescribeScanInstancesResponse, err error) {
    if request == nil {
        request = NewDescribeScanInstancesRequest()
    }
    response = NewDescribeScanInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeScanResultsRequest() (request *DescribeScanResultsRequest) {
    request = &DescribeScanResultsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ms", APIVersion, "DescribeScanResults")
    return
}

func NewDescribeScanResultsResponse() (response *DescribeScanResultsResponse) {
    response = &DescribeScanResultsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用户通过CreateScanInstances接口提交应用进行风险批量扫描后，用此接口批量获取风险详细信息,包含漏洞信息，广告信息，插件信息和病毒信息
func (c *Client) DescribeScanResults(request *DescribeScanResultsRequest) (response *DescribeScanResultsResponse, err error) {
    if request == nil {
        request = NewDescribeScanResultsRequest()
    }
    response = NewDescribeScanResultsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeShieldInstancesRequest() (request *DescribeShieldInstancesRequest) {
    request = &DescribeShieldInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ms", APIVersion, "DescribeShieldInstances")
    return
}

func NewDescribeShieldInstancesResponse() (response *DescribeShieldInstancesResponse) {
    response = &DescribeShieldInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口用于查看app列表。
// 可以通过指定任务唯一标识ItemId来查询指定app的详细信息，或通过设定过滤器来查询满足过滤条件的app的详细信息。 指定偏移(Offset)和限制(Limit)来选择结果中的一部分，默认返回满足条件的前20个app信息。
func (c *Client) DescribeShieldInstances(request *DescribeShieldInstancesRequest) (response *DescribeShieldInstancesResponse, err error) {
    if request == nil {
        request = NewDescribeShieldInstancesRequest()
    }
    response = NewDescribeShieldInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeShieldPlanInstanceRequest() (request *DescribeShieldPlanInstanceRequest) {
    request = &DescribeShieldPlanInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ms", APIVersion, "DescribeShieldPlanInstance")
    return
}

func NewDescribeShieldPlanInstanceResponse() (response *DescribeShieldPlanInstanceResponse) {
    response = &DescribeShieldPlanInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 查询加固策略
func (c *Client) DescribeShieldPlanInstance(request *DescribeShieldPlanInstanceRequest) (response *DescribeShieldPlanInstanceResponse, err error) {
    if request == nil {
        request = NewDescribeShieldPlanInstanceRequest()
    }
    response = NewDescribeShieldPlanInstanceResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeShieldResultRequest() (request *DescribeShieldResultRequest) {
    request = &DescribeShieldResultRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ms", APIVersion, "DescribeShieldResult")
    return
}

func NewDescribeShieldResultResponse() (response *DescribeShieldResultResponse) {
    response = &DescribeShieldResultResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 通过唯一标识获取加固的结果
func (c *Client) DescribeShieldResult(request *DescribeShieldResultRequest) (response *DescribeShieldResultResponse, err error) {
    if request == nil {
        request = NewDescribeShieldResultRequest()
    }
    response = NewDescribeShieldResultResponse()
    err = c.Send(request, response)
    return
}
