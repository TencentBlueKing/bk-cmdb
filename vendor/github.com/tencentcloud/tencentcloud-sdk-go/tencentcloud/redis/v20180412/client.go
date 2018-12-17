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

package v20180412

import (
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2018-04-12"

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


func NewClearInstanceRequest() (request *ClearInstanceRequest) {
    request = &ClearInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("redis", APIVersion, "ClearInstance")
    return
}

func NewClearInstanceResponse() (response *ClearInstanceResponse) {
    response = &ClearInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 清空Redis实例的实例数据。
func (c *Client) ClearInstance(request *ClearInstanceRequest) (response *ClearInstanceResponse, err error) {
    if request == nil {
        request = NewClearInstanceRequest()
    }
    response = NewClearInstanceResponse()
    err = c.Send(request, response)
    return
}

func NewCreateInstancesRequest() (request *CreateInstancesRequest) {
    request = &CreateInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("redis", APIVersion, "CreateInstances")
    return
}

func NewCreateInstancesResponse() (response *CreateInstancesResponse) {
    response = &CreateInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 创建redis实例
func (c *Client) CreateInstances(request *CreateInstancesRequest) (response *CreateInstancesResponse, err error) {
    if request == nil {
        request = NewCreateInstancesRequest()
    }
    response = NewCreateInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeAutoBackupConfigRequest() (request *DescribeAutoBackupConfigRequest) {
    request = &DescribeAutoBackupConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("redis", APIVersion, "DescribeAutoBackupConfig")
    return
}

func NewDescribeAutoBackupConfigResponse() (response *DescribeAutoBackupConfigResponse) {
    response = &DescribeAutoBackupConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取备份配置
func (c *Client) DescribeAutoBackupConfig(request *DescribeAutoBackupConfigRequest) (response *DescribeAutoBackupConfigResponse, err error) {
    if request == nil {
        request = NewDescribeAutoBackupConfigRequest()
    }
    response = NewDescribeAutoBackupConfigResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeInstanceBackupsRequest() (request *DescribeInstanceBackupsRequest) {
    request = &DescribeInstanceBackupsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("redis", APIVersion, "DescribeInstanceBackups")
    return
}

func NewDescribeInstanceBackupsResponse() (response *DescribeInstanceBackupsResponse) {
    response = &DescribeInstanceBackupsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 查询 CRS 实例备份列表
func (c *Client) DescribeInstanceBackups(request *DescribeInstanceBackupsRequest) (response *DescribeInstanceBackupsResponse, err error) {
    if request == nil {
        request = NewDescribeInstanceBackupsRequest()
    }
    response = NewDescribeInstanceBackupsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeInstanceDealDetailRequest() (request *DescribeInstanceDealDetailRequest) {
    request = &DescribeInstanceDealDetailRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("redis", APIVersion, "DescribeInstanceDealDetail")
    return
}

func NewDescribeInstanceDealDetailResponse() (response *DescribeInstanceDealDetailResponse) {
    response = &DescribeInstanceDealDetailResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 查询订单信息
func (c *Client) DescribeInstanceDealDetail(request *DescribeInstanceDealDetailRequest) (response *DescribeInstanceDealDetailResponse, err error) {
    if request == nil {
        request = NewDescribeInstanceDealDetailRequest()
    }
    response = NewDescribeInstanceDealDetailResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeInstancesRequest() (request *DescribeInstancesRequest) {
    request = &DescribeInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("redis", APIVersion, "DescribeInstances")
    return
}

func NewDescribeInstancesResponse() (response *DescribeInstancesResponse) {
    response = &DescribeInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 查询Redis实例列表
func (c *Client) DescribeInstances(request *DescribeInstancesRequest) (response *DescribeInstancesResponse, err error) {
    if request == nil {
        request = NewDescribeInstancesRequest()
    }
    response = NewDescribeInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeProductInfoRequest() (request *DescribeProductInfoRequest) {
    request = &DescribeProductInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("redis", APIVersion, "DescribeProductInfo")
    return
}

func NewDescribeProductInfoResponse() (response *DescribeProductInfoResponse) {
    response = &DescribeProductInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口查询指定可用区和实例类型下 Redis 的售卖规格， 如果用户不在购买白名单中，将不能查询该可用区或该类型的售卖规格详情。申请购买某地域白名单可以提交工单
func (c *Client) DescribeProductInfo(request *DescribeProductInfoRequest) (response *DescribeProductInfoResponse, err error) {
    if request == nil {
        request = NewDescribeProductInfoRequest()
    }
    response = NewDescribeProductInfoResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTaskInfoRequest() (request *DescribeTaskInfoRequest) {
    request = &DescribeTaskInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("redis", APIVersion, "DescribeTaskInfo")
    return
}

func NewDescribeTaskInfoResponse() (response *DescribeTaskInfoResponse) {
    response = &DescribeTaskInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用于查询任务结果
func (c *Client) DescribeTaskInfo(request *DescribeTaskInfoRequest) (response *DescribeTaskInfoResponse, err error) {
    if request == nil {
        request = NewDescribeTaskInfoRequest()
    }
    response = NewDescribeTaskInfoResponse()
    err = c.Send(request, response)
    return
}

func NewManualBackupInstanceRequest() (request *ManualBackupInstanceRequest) {
    request = &ManualBackupInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("redis", APIVersion, "ManualBackupInstance")
    return
}

func NewManualBackupInstanceResponse() (response *ManualBackupInstanceResponse) {
    response = &ManualBackupInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 手动备份Redis实例
func (c *Client) ManualBackupInstance(request *ManualBackupInstanceRequest) (response *ManualBackupInstanceResponse, err error) {
    if request == nil {
        request = NewManualBackupInstanceRequest()
    }
    response = NewManualBackupInstanceResponse()
    err = c.Send(request, response)
    return
}

func NewModfiyInstancePasswordRequest() (request *ModfiyInstancePasswordRequest) {
    request = &ModfiyInstancePasswordRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("redis", APIVersion, "ModfiyInstancePassword")
    return
}

func NewModfiyInstancePasswordResponse() (response *ModfiyInstancePasswordResponse) {
    response = &ModfiyInstancePasswordResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 修改redis密码
func (c *Client) ModfiyInstancePassword(request *ModfiyInstancePasswordRequest) (response *ModfiyInstancePasswordResponse, err error) {
    if request == nil {
        request = NewModfiyInstancePasswordRequest()
    }
    response = NewModfiyInstancePasswordResponse()
    err = c.Send(request, response)
    return
}

func NewModifyAutoBackupConfigRequest() (request *ModifyAutoBackupConfigRequest) {
    request = &ModifyAutoBackupConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("redis", APIVersion, "ModifyAutoBackupConfig")
    return
}

func NewModifyAutoBackupConfigResponse() (response *ModifyAutoBackupConfigResponse) {
    response = &ModifyAutoBackupConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 设置自动备份时间
func (c *Client) ModifyAutoBackupConfig(request *ModifyAutoBackupConfigRequest) (response *ModifyAutoBackupConfigResponse, err error) {
    if request == nil {
        request = NewModifyAutoBackupConfigRequest()
    }
    response = NewModifyAutoBackupConfigResponse()
    err = c.Send(request, response)
    return
}

func NewRenewInstanceRequest() (request *RenewInstanceRequest) {
    request = &RenewInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("redis", APIVersion, "RenewInstance")
    return
}

func NewRenewInstanceResponse() (response *RenewInstanceResponse) {
    response = &RenewInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 续费实例
func (c *Client) RenewInstance(request *RenewInstanceRequest) (response *RenewInstanceResponse, err error) {
    if request == nil {
        request = NewRenewInstanceRequest()
    }
    response = NewRenewInstanceResponse()
    err = c.Send(request, response)
    return
}

func NewResetPasswordRequest() (request *ResetPasswordRequest) {
    request = &ResetPasswordRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("redis", APIVersion, "ResetPassword")
    return
}

func NewResetPasswordResponse() (response *ResetPasswordResponse) {
    response = &ResetPasswordResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 重置密码
func (c *Client) ResetPassword(request *ResetPasswordRequest) (response *ResetPasswordResponse, err error) {
    if request == nil {
        request = NewResetPasswordRequest()
    }
    response = NewResetPasswordResponse()
    err = c.Send(request, response)
    return
}

func NewUpgradeInstanceRequest() (request *UpgradeInstanceRequest) {
    request = &UpgradeInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("redis", APIVersion, "UpgradeInstance")
    return
}

func NewUpgradeInstanceResponse() (response *UpgradeInstanceResponse) {
    response = &UpgradeInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 升级实例
func (c *Client) UpgradeInstance(request *UpgradeInstanceRequest) (response *UpgradeInstanceResponse, err error) {
    if request == nil {
        request = NewUpgradeInstanceRequest()
    }
    response = NewUpgradeInstanceResponse()
    err = c.Send(request, response)
    return
}
