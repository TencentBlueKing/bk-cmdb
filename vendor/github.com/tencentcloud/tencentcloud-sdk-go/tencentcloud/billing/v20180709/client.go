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

package v20180709

import (
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2018-07-09"

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


func NewDescribeAccountBalanceRequest() (request *DescribeAccountBalanceRequest) {
    request = &DescribeAccountBalanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("billing", APIVersion, "DescribeAccountBalance")
    return
}

func NewDescribeAccountBalanceResponse() (response *DescribeAccountBalanceResponse) {
    response = &DescribeAccountBalanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取云账户余额信息。
func (c *Client) DescribeAccountBalance(request *DescribeAccountBalanceRequest) (response *DescribeAccountBalanceResponse, err error) {
    if request == nil {
        request = NewDescribeAccountBalanceRequest()
    }
    response = NewDescribeAccountBalanceResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeBillDetailRequest() (request *DescribeBillDetailRequest) {
    request = &DescribeBillDetailRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("billing", APIVersion, "DescribeBillDetail")
    return
}

func NewDescribeBillDetailResponse() (response *DescribeBillDetailResponse) {
    response = &DescribeBillDetailResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 查询账单明细数据
func (c *Client) DescribeBillDetail(request *DescribeBillDetailRequest) (response *DescribeBillDetailResponse, err error) {
    if request == nil {
        request = NewDescribeBillDetailRequest()
    }
    response = NewDescribeBillDetailResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeBillResourceSummaryRequest() (request *DescribeBillResourceSummaryRequest) {
    request = &DescribeBillResourceSummaryRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("billing", APIVersion, "DescribeBillResourceSummary")
    return
}

func NewDescribeBillResourceSummaryResponse() (response *DescribeBillResourceSummaryResponse) {
    response = &DescribeBillResourceSummaryResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 查询账单资源汇总数据
func (c *Client) DescribeBillResourceSummary(request *DescribeBillResourceSummaryRequest) (response *DescribeBillResourceSummaryResponse, err error) {
    if request == nil {
        request = NewDescribeBillResourceSummaryRequest()
    }
    response = NewDescribeBillResourceSummaryResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDealsByCondRequest() (request *DescribeDealsByCondRequest) {
    request = &DescribeDealsByCondRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("billing", APIVersion, "DescribeDealsByCond")
    return
}

func NewDescribeDealsByCondResponse() (response *DescribeDealsByCondResponse) {
    response = &DescribeDealsByCondResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 查询订单
func (c *Client) DescribeDealsByCond(request *DescribeDealsByCondRequest) (response *DescribeDealsByCondResponse, err error) {
    if request == nil {
        request = NewDescribeDealsByCondRequest()
    }
    response = NewDescribeDealsByCondResponse()
    err = c.Send(request, response)
    return
}

func NewPayDealsRequest() (request *PayDealsRequest) {
    request = &PayDealsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("billing", APIVersion, "PayDeals")
    return
}

func NewPayDealsResponse() (response *PayDealsResponse) {
    response = &PayDealsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 支付订单
func (c *Client) PayDeals(request *PayDealsRequest) (response *PayDealsResponse, err error) {
    if request == nil {
        request = NewPayDealsRequest()
    }
    response = NewPayDealsResponse()
    err = c.Send(request, response)
    return
}
