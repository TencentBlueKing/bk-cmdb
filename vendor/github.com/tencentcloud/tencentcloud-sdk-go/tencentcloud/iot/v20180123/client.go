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

package v20180123

import (
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2018-01-23"

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


func NewActivateRuleRequest() (request *ActivateRuleRequest) {
    request = &ActivateRuleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "ActivateRule")
    return
}

func NewActivateRuleResponse() (response *ActivateRuleResponse) {
    response = &ActivateRuleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 启用规则
func (c *Client) ActivateRule(request *ActivateRuleRequest) (response *ActivateRuleResponse, err error) {
    if request == nil {
        request = NewActivateRuleRequest()
    }
    response = NewActivateRuleResponse()
    err = c.Send(request, response)
    return
}

func NewAddDeviceRequest() (request *AddDeviceRequest) {
    request = &AddDeviceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AddDevice")
    return
}

func NewAddDeviceResponse() (response *AddDeviceResponse) {
    response = &AddDeviceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 提供在指定的产品Id下创建一个设备的能力，生成设备名称与设备秘钥。
func (c *Client) AddDevice(request *AddDeviceRequest) (response *AddDeviceResponse, err error) {
    if request == nil {
        request = NewAddDeviceRequest()
    }
    response = NewAddDeviceResponse()
    err = c.Send(request, response)
    return
}

func NewAddProductRequest() (request *AddProductRequest) {
    request = &AddProductRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AddProduct")
    return
}

func NewAddProductResponse() (response *AddProductResponse) {
    response = &AddProductResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口(AddProduct)用于创建、定义某款硬件产品。
func (c *Client) AddProduct(request *AddProductRequest) (response *AddProductResponse, err error) {
    if request == nil {
        request = NewAddProductRequest()
    }
    response = NewAddProductResponse()
    err = c.Send(request, response)
    return
}

func NewAddRuleRequest() (request *AddRuleRequest) {
    request = &AddRuleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AddRule")
    return
}

func NewAddRuleResponse() (response *AddRuleResponse) {
    response = &AddRuleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 新增规则
func (c *Client) AddRule(request *AddRuleRequest) (response *AddRuleResponse, err error) {
    if request == nil {
        request = NewAddRuleRequest()
    }
    response = NewAddRuleResponse()
    err = c.Send(request, response)
    return
}

func NewAddTopicRequest() (request *AddTopicRequest) {
    request = &AddTopicRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AddTopic")
    return
}

func NewAddTopicResponse() (response *AddTopicResponse) {
    response = &AddTopicResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 新增Topic，用于设备或应用发布消息至该Topic或订阅该Topic的消息。
func (c *Client) AddTopic(request *AddTopicRequest) (response *AddTopicResponse, err error) {
    if request == nil {
        request = NewAddTopicRequest()
    }
    response = NewAddTopicResponse()
    err = c.Send(request, response)
    return
}

func NewAppAddUserRequest() (request *AppAddUserRequest) {
    request = &AppAddUserRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AppAddUser")
    return
}

func NewAppAddUserResponse() (response *AppAddUserResponse) {
    response = &AppAddUserResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 为APP提供用户注册功能
func (c *Client) AppAddUser(request *AppAddUserRequest) (response *AppAddUserResponse, err error) {
    if request == nil {
        request = NewAppAddUserRequest()
    }
    response = NewAppAddUserResponse()
    err = c.Send(request, response)
    return
}

func NewAppDeleteDeviceRequest() (request *AppDeleteDeviceRequest) {
    request = &AppDeleteDeviceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AppDeleteDevice")
    return
}

func NewAppDeleteDeviceResponse() (response *AppDeleteDeviceResponse) {
    response = &AppDeleteDeviceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用户解除与设备的关联关系，解除后APP用户无法控制设备，获取设备数据
func (c *Client) AppDeleteDevice(request *AppDeleteDeviceRequest) (response *AppDeleteDeviceResponse, err error) {
    if request == nil {
        request = NewAppDeleteDeviceRequest()
    }
    response = NewAppDeleteDeviceResponse()
    err = c.Send(request, response)
    return
}

func NewAppGetDeviceRequest() (request *AppGetDeviceRequest) {
    request = &AppGetDeviceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AppGetDevice")
    return
}

func NewAppGetDeviceResponse() (response *AppGetDeviceResponse) {
    response = &AppGetDeviceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取绑定设备的基本信息与数据模板定义
func (c *Client) AppGetDevice(request *AppGetDeviceRequest) (response *AppGetDeviceResponse, err error) {
    if request == nil {
        request = NewAppGetDeviceRequest()
    }
    response = NewAppGetDeviceResponse()
    err = c.Send(request, response)
    return
}

func NewAppGetDeviceDataRequest() (request *AppGetDeviceDataRequest) {
    request = &AppGetDeviceDataRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AppGetDeviceData")
    return
}

func NewAppGetDeviceDataResponse() (response *AppGetDeviceDataResponse) {
    response = &AppGetDeviceDataResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取绑定设备数据，用于实时展示设备的最新数据
func (c *Client) AppGetDeviceData(request *AppGetDeviceDataRequest) (response *AppGetDeviceDataResponse, err error) {
    if request == nil {
        request = NewAppGetDeviceDataRequest()
    }
    response = NewAppGetDeviceDataResponse()
    err = c.Send(request, response)
    return
}

func NewAppGetDeviceStatusesRequest() (request *AppGetDeviceStatusesRequest) {
    request = &AppGetDeviceStatusesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AppGetDeviceStatuses")
    return
}

func NewAppGetDeviceStatusesResponse() (response *AppGetDeviceStatusesResponse) {
    response = &AppGetDeviceStatusesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取绑定设备的上下线状态
func (c *Client) AppGetDeviceStatuses(request *AppGetDeviceStatusesRequest) (response *AppGetDeviceStatusesResponse, err error) {
    if request == nil {
        request = NewAppGetDeviceStatusesRequest()
    }
    response = NewAppGetDeviceStatusesResponse()
    err = c.Send(request, response)
    return
}

func NewAppGetDevicesRequest() (request *AppGetDevicesRequest) {
    request = &AppGetDevicesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AppGetDevices")
    return
}

func NewAppGetDevicesResponse() (response *AppGetDevicesResponse) {
    response = &AppGetDevicesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取用户的绑定设备列表
func (c *Client) AppGetDevices(request *AppGetDevicesRequest) (response *AppGetDevicesResponse, err error) {
    if request == nil {
        request = NewAppGetDevicesRequest()
    }
    response = NewAppGetDevicesResponse()
    err = c.Send(request, response)
    return
}

func NewAppGetTokenRequest() (request *AppGetTokenRequest) {
    request = &AppGetTokenRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AppGetToken")
    return
}

func NewAppGetTokenResponse() (response *AppGetTokenResponse) {
    response = &AppGetTokenResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取用户token
func (c *Client) AppGetToken(request *AppGetTokenRequest) (response *AppGetTokenResponse, err error) {
    if request == nil {
        request = NewAppGetTokenRequest()
    }
    response = NewAppGetTokenResponse()
    err = c.Send(request, response)
    return
}

func NewAppGetUserRequest() (request *AppGetUserRequest) {
    request = &AppGetUserRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AppGetUser")
    return
}

func NewAppGetUserResponse() (response *AppGetUserResponse) {
    response = &AppGetUserResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取用户信息
func (c *Client) AppGetUser(request *AppGetUserRequest) (response *AppGetUserResponse, err error) {
    if request == nil {
        request = NewAppGetUserRequest()
    }
    response = NewAppGetUserResponse()
    err = c.Send(request, response)
    return
}

func NewAppIssueDeviceControlRequest() (request *AppIssueDeviceControlRequest) {
    request = &AppIssueDeviceControlRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AppIssueDeviceControl")
    return
}

func NewAppIssueDeviceControlResponse() (response *AppIssueDeviceControlResponse) {
    response = &AppIssueDeviceControlResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用户通过APP控制设备
func (c *Client) AppIssueDeviceControl(request *AppIssueDeviceControlRequest) (response *AppIssueDeviceControlResponse, err error) {
    if request == nil {
        request = NewAppIssueDeviceControlRequest()
    }
    response = NewAppIssueDeviceControlResponse()
    err = c.Send(request, response)
    return
}

func NewAppResetPasswordRequest() (request *AppResetPasswordRequest) {
    request = &AppResetPasswordRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AppResetPassword")
    return
}

func NewAppResetPasswordResponse() (response *AppResetPasswordResponse) {
    response = &AppResetPasswordResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 重置APP用户密码
func (c *Client) AppResetPassword(request *AppResetPasswordRequest) (response *AppResetPasswordResponse, err error) {
    if request == nil {
        request = NewAppResetPasswordRequest()
    }
    response = NewAppResetPasswordResponse()
    err = c.Send(request, response)
    return
}

func NewAppSecureAddDeviceRequest() (request *AppSecureAddDeviceRequest) {
    request = &AppSecureAddDeviceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AppSecureAddDevice")
    return
}

func NewAppSecureAddDeviceResponse() (response *AppSecureAddDeviceResponse) {
    response = &AppSecureAddDeviceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 用户绑定设备，绑定后可以在APP端进行控制。绑定设备前需调用“获取设备绑定签名”接口
func (c *Client) AppSecureAddDevice(request *AppSecureAddDeviceRequest) (response *AppSecureAddDeviceResponse, err error) {
    if request == nil {
        request = NewAppSecureAddDeviceRequest()
    }
    response = NewAppSecureAddDeviceResponse()
    err = c.Send(request, response)
    return
}

func NewAppUpdateDeviceRequest() (request *AppUpdateDeviceRequest) {
    request = &AppUpdateDeviceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AppUpdateDevice")
    return
}

func NewAppUpdateDeviceResponse() (response *AppUpdateDeviceResponse) {
    response = &AppUpdateDeviceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 修改设备别名，便于用户个性化定义设备的名称
func (c *Client) AppUpdateDevice(request *AppUpdateDeviceRequest) (response *AppUpdateDeviceResponse, err error) {
    if request == nil {
        request = NewAppUpdateDeviceRequest()
    }
    response = NewAppUpdateDeviceResponse()
    err = c.Send(request, response)
    return
}

func NewAppUpdateUserRequest() (request *AppUpdateUserRequest) {
    request = &AppUpdateUserRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AppUpdateUser")
    return
}

func NewAppUpdateUserResponse() (response *AppUpdateUserResponse) {
    response = &AppUpdateUserResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 修改用户信息
func (c *Client) AppUpdateUser(request *AppUpdateUserRequest) (response *AppUpdateUserResponse, err error) {
    if request == nil {
        request = NewAppUpdateUserRequest()
    }
    response = NewAppUpdateUserResponse()
    err = c.Send(request, response)
    return
}

func NewAssociateSubDeviceToGatewayProductRequest() (request *AssociateSubDeviceToGatewayProductRequest) {
    request = &AssociateSubDeviceToGatewayProductRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "AssociateSubDeviceToGatewayProduct")
    return
}

func NewAssociateSubDeviceToGatewayProductResponse() (response *AssociateSubDeviceToGatewayProductResponse) {
    response = &AssociateSubDeviceToGatewayProductResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 关联子设备产品和网关产品
func (c *Client) AssociateSubDeviceToGatewayProduct(request *AssociateSubDeviceToGatewayProductRequest) (response *AssociateSubDeviceToGatewayProductResponse, err error) {
    if request == nil {
        request = NewAssociateSubDeviceToGatewayProductRequest()
    }
    response = NewAssociateSubDeviceToGatewayProductResponse()
    err = c.Send(request, response)
    return
}

func NewDeactivateRuleRequest() (request *DeactivateRuleRequest) {
    request = &DeactivateRuleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "DeactivateRule")
    return
}

func NewDeactivateRuleResponse() (response *DeactivateRuleResponse) {
    response = &DeactivateRuleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 禁用规则
func (c *Client) DeactivateRule(request *DeactivateRuleRequest) (response *DeactivateRuleResponse, err error) {
    if request == nil {
        request = NewDeactivateRuleRequest()
    }
    response = NewDeactivateRuleResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteDeviceRequest() (request *DeleteDeviceRequest) {
    request = &DeleteDeviceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "DeleteDevice")
    return
}

func NewDeleteDeviceResponse() (response *DeleteDeviceResponse) {
    response = &DeleteDeviceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 提供在指定的产品Id下删除一个设备的能力。
func (c *Client) DeleteDevice(request *DeleteDeviceRequest) (response *DeleteDeviceResponse, err error) {
    if request == nil {
        request = NewDeleteDeviceRequest()
    }
    response = NewDeleteDeviceResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteProductRequest() (request *DeleteProductRequest) {
    request = &DeleteProductRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "DeleteProduct")
    return
}

func NewDeleteProductResponse() (response *DeleteProductResponse) {
    response = &DeleteProductResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 删除用户指定的产品Id对应的信息。
func (c *Client) DeleteProduct(request *DeleteProductRequest) (response *DeleteProductResponse, err error) {
    if request == nil {
        request = NewDeleteProductRequest()
    }
    response = NewDeleteProductResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteRuleRequest() (request *DeleteRuleRequest) {
    request = &DeleteRuleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "DeleteRule")
    return
}

func NewDeleteRuleResponse() (response *DeleteRuleResponse) {
    response = &DeleteRuleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 删除规则
func (c *Client) DeleteRule(request *DeleteRuleRequest) (response *DeleteRuleResponse, err error) {
    if request == nil {
        request = NewDeleteRuleRequest()
    }
    response = NewDeleteRuleResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteTopicRequest() (request *DeleteTopicRequest) {
    request = &DeleteTopicRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "DeleteTopic")
    return
}

func NewDeleteTopicResponse() (response *DeleteTopicResponse) {
    response = &DeleteTopicResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 删除Topic
func (c *Client) DeleteTopic(request *DeleteTopicRequest) (response *DeleteTopicResponse, err error) {
    if request == nil {
        request = NewDeleteTopicRequest()
    }
    response = NewDeleteTopicResponse()
    err = c.Send(request, response)
    return
}

func NewGetDataHistoryRequest() (request *GetDataHistoryRequest) {
    request = &GetDataHistoryRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetDataHistory")
    return
}

func NewGetDataHistoryResponse() (response *GetDataHistoryResponse) {
    response = &GetDataHistoryResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 批量获取设备某一段时间范围的设备上报数据。该接口适用于使用高级版类型的产品
func (c *Client) GetDataHistory(request *GetDataHistoryRequest) (response *GetDataHistoryResponse, err error) {
    if request == nil {
        request = NewGetDataHistoryRequest()
    }
    response = NewGetDataHistoryResponse()
    err = c.Send(request, response)
    return
}

func NewGetDebugLogRequest() (request *GetDebugLogRequest) {
    request = &GetDebugLogRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetDebugLog")
    return
}

func NewGetDebugLogResponse() (response *GetDebugLogResponse) {
    response = &GetDebugLogResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取设备的调试日志，用于定位问题
func (c *Client) GetDebugLog(request *GetDebugLogRequest) (response *GetDebugLogResponse, err error) {
    if request == nil {
        request = NewGetDebugLogRequest()
    }
    response = NewGetDebugLogResponse()
    err = c.Send(request, response)
    return
}

func NewGetDeviceRequest() (request *GetDeviceRequest) {
    request = &GetDeviceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetDevice")
    return
}

func NewGetDeviceResponse() (response *GetDeviceResponse) {
    response = &GetDeviceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 提供查询某个设备详细信息的能力。
func (c *Client) GetDevice(request *GetDeviceRequest) (response *GetDeviceResponse, err error) {
    if request == nil {
        request = NewGetDeviceRequest()
    }
    response = NewGetDeviceResponse()
    err = c.Send(request, response)
    return
}

func NewGetDeviceDataRequest() (request *GetDeviceDataRequest) {
    request = &GetDeviceDataRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetDeviceData")
    return
}

func NewGetDeviceDataResponse() (response *GetDeviceDataResponse) {
    response = &GetDeviceDataResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取某个设备当前上报到云端的数据，该接口适用于使用数据模板协议的产品。
func (c *Client) GetDeviceData(request *GetDeviceDataRequest) (response *GetDeviceDataResponse, err error) {
    if request == nil {
        request = NewGetDeviceDataRequest()
    }
    response = NewGetDeviceDataResponse()
    err = c.Send(request, response)
    return
}

func NewGetDeviceLogRequest() (request *GetDeviceLogRequest) {
    request = &GetDeviceLogRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetDeviceLog")
    return
}

func NewGetDeviceLogResponse() (response *GetDeviceLogResponse) {
    response = &GetDeviceLogResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 批量获取设备与云端的详细通信日志，该接口适用于使用高级版类型的产品。
func (c *Client) GetDeviceLog(request *GetDeviceLogRequest) (response *GetDeviceLogResponse, err error) {
    if request == nil {
        request = NewGetDeviceLogRequest()
    }
    response = NewGetDeviceLogResponse()
    err = c.Send(request, response)
    return
}

func NewGetDeviceSignaturesRequest() (request *GetDeviceSignaturesRequest) {
    request = &GetDeviceSignaturesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetDeviceSignatures")
    return
}

func NewGetDeviceSignaturesResponse() (response *GetDeviceSignaturesResponse) {
    response = &GetDeviceSignaturesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取设备绑定签名，用于用户绑定某个设备的应用场景
func (c *Client) GetDeviceSignatures(request *GetDeviceSignaturesRequest) (response *GetDeviceSignaturesResponse, err error) {
    if request == nil {
        request = NewGetDeviceSignaturesRequest()
    }
    response = NewGetDeviceSignaturesResponse()
    err = c.Send(request, response)
    return
}

func NewGetDeviceStatisticsRequest() (request *GetDeviceStatisticsRequest) {
    request = &GetDeviceStatisticsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetDeviceStatistics")
    return
}

func NewGetDeviceStatisticsResponse() (response *GetDeviceStatisticsResponse) {
    response = &GetDeviceStatisticsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 查询某段时间范围内产品的在线、激活设备数
func (c *Client) GetDeviceStatistics(request *GetDeviceStatisticsRequest) (response *GetDeviceStatisticsResponse, err error) {
    if request == nil {
        request = NewGetDeviceStatisticsRequest()
    }
    response = NewGetDeviceStatisticsResponse()
    err = c.Send(request, response)
    return
}

func NewGetDeviceStatusesRequest() (request *GetDeviceStatusesRequest) {
    request = &GetDeviceStatusesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetDeviceStatuses")
    return
}

func NewGetDeviceStatusesResponse() (response *GetDeviceStatusesResponse) {
    response = &GetDeviceStatusesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 批量获取设备的当前状态，状态包括在线、离线或未激活状态。
func (c *Client) GetDeviceStatuses(request *GetDeviceStatusesRequest) (response *GetDeviceStatusesResponse, err error) {
    if request == nil {
        request = NewGetDeviceStatusesRequest()
    }
    response = NewGetDeviceStatusesResponse()
    err = c.Send(request, response)
    return
}

func NewGetDevicesRequest() (request *GetDevicesRequest) {
    request = &GetDevicesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetDevices")
    return
}

func NewGetDevicesResponse() (response *GetDevicesResponse) {
    response = &GetDevicesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 提供分页查询某个产品Id下设备信息的能力。
func (c *Client) GetDevices(request *GetDevicesRequest) (response *GetDevicesResponse, err error) {
    if request == nil {
        request = NewGetDevicesRequest()
    }
    response = NewGetDevicesResponse()
    err = c.Send(request, response)
    return
}

func NewGetProductRequest() (request *GetProductRequest) {
    request = &GetProductRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetProduct")
    return
}

func NewGetProductResponse() (response *GetProductResponse) {
    response = &GetProductResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取产品定义的详细信息，包括产品名称、产品描述，鉴权模式等信息。
func (c *Client) GetProduct(request *GetProductRequest) (response *GetProductResponse, err error) {
    if request == nil {
        request = NewGetProductRequest()
    }
    response = NewGetProductResponse()
    err = c.Send(request, response)
    return
}

func NewGetProductsRequest() (request *GetProductsRequest) {
    request = &GetProductsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetProducts")
    return
}

func NewGetProductsResponse() (response *GetProductsResponse) {
    response = &GetProductsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取用户在物联网套件所创建的所有产品信息。
func (c *Client) GetProducts(request *GetProductsRequest) (response *GetProductsResponse, err error) {
    if request == nil {
        request = NewGetProductsRequest()
    }
    response = NewGetProductsResponse()
    err = c.Send(request, response)
    return
}

func NewGetRuleRequest() (request *GetRuleRequest) {
    request = &GetRuleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetRule")
    return
}

func NewGetRuleResponse() (response *GetRuleResponse) {
    response = &GetRuleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取转发规则信息
func (c *Client) GetRule(request *GetRuleRequest) (response *GetRuleResponse, err error) {
    if request == nil {
        request = NewGetRuleRequest()
    }
    response = NewGetRuleResponse()
    err = c.Send(request, response)
    return
}

func NewGetRulesRequest() (request *GetRulesRequest) {
    request = &GetRulesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetRules")
    return
}

func NewGetRulesResponse() (response *GetRulesResponse) {
    response = &GetRulesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取转发规则列表
func (c *Client) GetRules(request *GetRulesRequest) (response *GetRulesResponse, err error) {
    if request == nil {
        request = NewGetRulesRequest()
    }
    response = NewGetRulesResponse()
    err = c.Send(request, response)
    return
}

func NewGetTopicRequest() (request *GetTopicRequest) {
    request = &GetTopicRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetTopic")
    return
}

func NewGetTopicResponse() (response *GetTopicResponse) {
    response = &GetTopicResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取Topic信息
func (c *Client) GetTopic(request *GetTopicRequest) (response *GetTopicResponse, err error) {
    if request == nil {
        request = NewGetTopicRequest()
    }
    response = NewGetTopicResponse()
    err = c.Send(request, response)
    return
}

func NewGetTopicsRequest() (request *GetTopicsRequest) {
    request = &GetTopicsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "GetTopics")
    return
}

func NewGetTopicsResponse() (response *GetTopicsResponse) {
    response = &GetTopicsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取Topic列表
func (c *Client) GetTopics(request *GetTopicsRequest) (response *GetTopicsResponse, err error) {
    if request == nil {
        request = NewGetTopicsRequest()
    }
    response = NewGetTopicsResponse()
    err = c.Send(request, response)
    return
}

func NewIssueDeviceControlRequest() (request *IssueDeviceControlRequest) {
    request = &IssueDeviceControlRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "IssueDeviceControl")
    return
}

func NewIssueDeviceControlResponse() (response *IssueDeviceControlResponse) {
    response = &IssueDeviceControlResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 提供下发控制指令到指定设备的能力，该接口适用于使用高级版类型的产品。
func (c *Client) IssueDeviceControl(request *IssueDeviceControlRequest) (response *IssueDeviceControlResponse, err error) {
    if request == nil {
        request = NewIssueDeviceControlRequest()
    }
    response = NewIssueDeviceControlResponse()
    err = c.Send(request, response)
    return
}

func NewPublishMsgRequest() (request *PublishMsgRequest) {
    request = &PublishMsgRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "PublishMsg")
    return
}

func NewPublishMsgResponse() (response *PublishMsgResponse) {
    response = &PublishMsgResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 提供向指定的Topic发布消息的能力，常用于向设备下发控制指令。该接口只适用于产品版本为“基础版”类型的产品，使用高级版的产品需使用“下发设备控制指令”接口
func (c *Client) PublishMsg(request *PublishMsgRequest) (response *PublishMsgResponse, err error) {
    if request == nil {
        request = NewPublishMsgRequest()
    }
    response = NewPublishMsgResponse()
    err = c.Send(request, response)
    return
}

func NewResetDeviceRequest() (request *ResetDeviceRequest) {
    request = &ResetDeviceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "ResetDevice")
    return
}

func NewResetDeviceResponse() (response *ResetDeviceResponse) {
    response = &ResetDeviceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 重置设备操作，将会为设备生成新的证书及清空最新数据，需谨慎操作。
func (c *Client) ResetDevice(request *ResetDeviceRequest) (response *ResetDeviceResponse, err error) {
    if request == nil {
        request = NewResetDeviceRequest()
    }
    response = NewResetDeviceResponse()
    err = c.Send(request, response)
    return
}

func NewUnassociateSubDeviceFromGatewayProductRequest() (request *UnassociateSubDeviceFromGatewayProductRequest) {
    request = &UnassociateSubDeviceFromGatewayProductRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "UnassociateSubDeviceFromGatewayProduct")
    return
}

func NewUnassociateSubDeviceFromGatewayProductResponse() (response *UnassociateSubDeviceFromGatewayProductResponse) {
    response = &UnassociateSubDeviceFromGatewayProductResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 取消子设备产品与网关设备产品的关联
func (c *Client) UnassociateSubDeviceFromGatewayProduct(request *UnassociateSubDeviceFromGatewayProductRequest) (response *UnassociateSubDeviceFromGatewayProductResponse, err error) {
    if request == nil {
        request = NewUnassociateSubDeviceFromGatewayProductRequest()
    }
    response = NewUnassociateSubDeviceFromGatewayProductResponse()
    err = c.Send(request, response)
    return
}

func NewUpdateProductRequest() (request *UpdateProductRequest) {
    request = &UpdateProductRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "UpdateProduct")
    return
}

func NewUpdateProductResponse() (response *UpdateProductResponse) {
    response = &UpdateProductResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 提供修改产品信息及数据模板的能力。
func (c *Client) UpdateProduct(request *UpdateProductRequest) (response *UpdateProductResponse, err error) {
    if request == nil {
        request = NewUpdateProductRequest()
    }
    response = NewUpdateProductResponse()
    err = c.Send(request, response)
    return
}

func NewUpdateRuleRequest() (request *UpdateRuleRequest) {
    request = &UpdateRuleRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("iot", APIVersion, "UpdateRule")
    return
}

func NewUpdateRuleResponse() (response *UpdateRuleResponse) {
    response = &UpdateRuleResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 更新规则
func (c *Client) UpdateRule(request *UpdateRuleRequest) (response *UpdateRuleResponse, err error) {
    if request == nil {
        request = NewUpdateRuleRequest()
    }
    response = NewUpdateRuleResponse()
    err = c.Send(request, response)
    return
}
