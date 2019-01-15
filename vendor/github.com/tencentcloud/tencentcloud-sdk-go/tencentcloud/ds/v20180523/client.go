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

package v20180523

import (
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2018-05-23"

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


func NewCheckVcodeRequest() (request *CheckVcodeRequest) {
    request = &CheckVcodeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ds", APIVersion, "CheckVcode")
    return
}

func NewCheckVcodeResponse() (response *CheckVcodeResponse) {
    response = &CheckVcodeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 检测验证码接口。此接口用于企业电子合同平台通过给用户发送短信验证码，以短信授权方式签署合同。此接口配合发送验证码接口使用。
// 
// 用户在企业电子合同平台输入收到的验证码后，由企业电子合同平台调用该接口向腾讯云提交确认受托签署合同验证码命令。验证码验证正确时，本次合同签署的授权成功。
func (c *Client) CheckVcode(request *CheckVcodeRequest) (response *CheckVcodeResponse, err error) {
    if request == nil {
        request = NewCheckVcodeRequest()
    }
    response = NewCheckVcodeResponse()
    err = c.Send(request, response)
    return
}

func NewCreateContractByUploadRequest() (request *CreateContractByUploadRequest) {
    request = &CreateContractByUploadRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ds", APIVersion, "CreateContractByUpload")
    return
}

func NewCreateContractByUploadResponse() (response *CreateContractByUploadResponse) {
    response = &CreateContractByUploadResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 此接口适用于：客户平台通过上传PDF文件作为合同，以备未来进行签署。接口返回任务号，可调用DescribeTaskStatus接口查看任务执行结果。
func (c *Client) CreateContractByUpload(request *CreateContractByUploadRequest) (response *CreateContractByUploadResponse, err error) {
    if request == nil {
        request = NewCreateContractByUploadRequest()
    }
    response = NewCreateContractByUploadResponse()
    err = c.Send(request, response)
    return
}

func NewCreateEnterpriseAccountRequest() (request *CreateEnterpriseAccountRequest) {
    request = &CreateEnterpriseAccountRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ds", APIVersion, "CreateEnterpriseAccount")
    return
}

func NewCreateEnterpriseAccountResponse() (response *CreateEnterpriseAccountResponse) {
    response = &CreateEnterpriseAccountResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 为企业电子合同平台的最终企业用户进行开户。在企业电子合同平台进行操作的企业用户，企业电子合同平台向腾讯云发送个人用户的信息，提交开户命令。腾讯云接到请求后，自动为企业电子合同平台的企业用户生成一张数字证书。
func (c *Client) CreateEnterpriseAccount(request *CreateEnterpriseAccountRequest) (response *CreateEnterpriseAccountResponse, err error) {
    if request == nil {
        request = NewCreateEnterpriseAccountRequest()
    }
    response = NewCreateEnterpriseAccountResponse()
    err = c.Send(request, response)
    return
}

func NewCreatePersonalAccountRequest() (request *CreatePersonalAccountRequest) {
    request = &CreatePersonalAccountRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ds", APIVersion, "CreatePersonalAccount")
    return
}

func NewCreatePersonalAccountResponse() (response *CreatePersonalAccountResponse) {
    response = &CreatePersonalAccountResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 为企业电子合同平台的最终个人用户进行开户。在企业电子合同平台进行操作的个人用户，企业电子合同平台向腾讯云发送个人用户的信息，提交开户命令。腾讯云接到请求后，自动为企业电子合同平台的个人用户生成一张数字证书。
func (c *Client) CreatePersonalAccount(request *CreatePersonalAccountRequest) (response *CreatePersonalAccountResponse, err error) {
    if request == nil {
        request = NewCreatePersonalAccountRequest()
    }
    response = NewCreatePersonalAccountResponse()
    err = c.Send(request, response)
    return
}

func NewCreateSealRequest() (request *CreateSealRequest) {
    request = &CreateSealRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ds", APIVersion, "CreateSeal")
    return
}

func NewCreateSealResponse() (response *CreateSealResponse) {
    response = &CreateSealResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 此接口用于客户电子合同平台增加某用户的印章图片。客户平台可以调用此接口增加某用户的印章图片。
func (c *Client) CreateSeal(request *CreateSealRequest) (response *CreateSealResponse, err error) {
    if request == nil {
        request = NewCreateSealRequest()
    }
    response = NewCreateSealResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteAccountRequest() (request *DeleteAccountRequest) {
    request = &DeleteAccountRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ds", APIVersion, "DeleteAccount")
    return
}

func NewDeleteAccountResponse() (response *DeleteAccountResponse) {
    response = &DeleteAccountResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 删除企业电子合同平台的最终用户。调用该接口后，腾讯云将删除该用户账号。删除账号后，已经签名的合同不受影响。
func (c *Client) DeleteAccount(request *DeleteAccountRequest) (response *DeleteAccountResponse, err error) {
    if request == nil {
        request = NewDeleteAccountRequest()
    }
    response = NewDeleteAccountResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteSealRequest() (request *DeleteSealRequest) {
    request = &DeleteSealRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ds", APIVersion, "DeleteSeal")
    return
}

func NewDeleteSealResponse() (response *DeleteSealResponse) {
    response = &DeleteSealResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 删除印章接口，删除指定账号的某个印章
func (c *Client) DeleteSeal(request *DeleteSealRequest) (response *DeleteSealResponse, err error) {
    if request == nil {
        request = NewDeleteSealRequest()
    }
    response = NewDeleteSealResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTaskStatusRequest() (request *DescribeTaskStatusRequest) {
    request = &DescribeTaskStatusRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ds", APIVersion, "DescribeTaskStatus")
    return
}

func NewDescribeTaskStatusResponse() (response *DescribeTaskStatusResponse) {
    response = &DescribeTaskStatusResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 接口使用于：客户平台可使用该接口查询任务执行状态或者执行结果
func (c *Client) DescribeTaskStatus(request *DescribeTaskStatusRequest) (response *DescribeTaskStatusResponse, err error) {
    if request == nil {
        request = NewDescribeTaskStatusRequest()
    }
    response = NewDescribeTaskStatusResponse()
    err = c.Send(request, response)
    return
}

func NewDownloadContractRequest() (request *DownloadContractRequest) {
    request = &DownloadContractRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ds", APIVersion, "DownloadContract")
    return
}

func NewDownloadContractResponse() (response *DownloadContractResponse) {
    response = &DownloadContractResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 下载合同接口。调用该接口可以下载签署中和签署完成的接口。接口返回任务号，可调用DescribeTaskStatus接口查看任务执行结果。
func (c *Client) DownloadContract(request *DownloadContractRequest) (response *DownloadContractResponse, err error) {
    if request == nil {
        request = NewDownloadContractRequest()
    }
    response = NewDownloadContractResponse()
    err = c.Send(request, response)
    return
}

func NewSendVcodeRequest() (request *SendVcodeRequest) {
    request = &SendVcodeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ds", APIVersion, "SendVcode")
    return
}

func NewSendVcodeResponse() (response *SendVcodeResponse) {
    response = &SendVcodeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 发送验证码接口。此接口用于：企业电子合同平台需要腾讯云发送验证码对其用户进行验证时调用，腾讯云将向其用户联系手机(企业电子合同平台为用户开户时通过接口传入)发送验证码，以验证码授权方式签署合同。企业电子合同平台可以选择签署合同时不校验验证码（需线下沟通）。用户验证工作由企业电子合同平台自身完成。
func (c *Client) SendVcode(request *SendVcodeRequest) (response *SendVcodeResponse, err error) {
    if request == nil {
        request = NewSendVcodeRequest()
    }
    response = NewSendVcodeResponse()
    err = c.Send(request, response)
    return
}

func NewSignContractByCoordinateRequest() (request *SignContractByCoordinateRequest) {
    request = &SignContractByCoordinateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ds", APIVersion, "SignContractByCoordinate")
    return
}

func NewSignContractByCoordinateResponse() (response *SignContractByCoordinateResponse) {
    response = &SignContractByCoordinateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 此接口适用于：客户平台在创建好合同后，由合同签署方对创建的合同内容进行确认，无误后再进行签署。客户平台使用该接口提供详细的PDF文档签名坐标进行签署。
func (c *Client) SignContractByCoordinate(request *SignContractByCoordinateRequest) (response *SignContractByCoordinateResponse, err error) {
    if request == nil {
        request = NewSignContractByCoordinateRequest()
    }
    response = NewSignContractByCoordinateResponse()
    err = c.Send(request, response)
    return
}

func NewSignContractByKeywordRequest() (request *SignContractByKeywordRequest) {
    request = &SignContractByKeywordRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("ds", APIVersion, "SignContractByKeyword")
    return
}

func NewSignContractByKeywordResponse() (response *SignContractByKeywordResponse) {
    response = &SignContractByKeywordResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 此接口适用于：客户平台在创建好合同后，由合同签署方对创建的合同内容进行确认，无误后再进行签署。客户平台使用该接口对PDF合同文档按照关键字和坐标进行签署。
func (c *Client) SignContractByKeyword(request *SignContractByKeywordRequest) (response *SignContractByKeywordResponse, err error) {
    if request == nil {
        request = NewSignContractByKeywordRequest()
    }
    response = NewSignContractByKeywordResponse()
    err = c.Send(request, response)
    return
}
