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

package v20180606

import (
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2018-06-06"

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


func NewDescribeCdnDataRequest() (request *DescribeCdnDataRequest) {
    request = &DescribeCdnDataRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cdn", APIVersion, "DescribeCdnData")
    return
}

func NewDescribeCdnDataResponse() (response *DescribeCdnDataResponse) {
    response = &DescribeCdnDataResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeCdnData 用于查询 CDN 实时访问监控数据，支持以下指标查询：
// 
// + 流量（单位为 byte）
// + 带宽（单位为 bps）
// + 请求数（单位为 次）
// + 流量命中率（单位为 %，小数点后保留两位）
// + 状态码 2XX 汇总及各 2 开头状态码明细（单位为 个）
// + 状态码 3XX 汇总及各 3 开头状态码明细（单位为 个）
// + 状态码 4XX 汇总及各 4 开头状态码明细（单位为 个）
// + 状态码 5XX 汇总及各 5 开头状态码明细（单位为 个）
func (c *Client) DescribeCdnData(request *DescribeCdnDataRequest) (response *DescribeCdnDataResponse, err error) {
    if request == nil {
        request = NewDescribeCdnDataRequest()
    }
    response = NewDescribeCdnDataResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeIpVisitRequest() (request *DescribeIpVisitRequest) {
    request = &DescribeIpVisitRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cdn", APIVersion, "DescribeIpVisit")
    return
}

func NewDescribeIpVisitResponse() (response *DescribeIpVisitResponse) {
    response = &DescribeIpVisitResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeIpVisit 用于查询 5 分钟活跃用户数，及日活跃用户数明细
// 
// + 5 分钟活跃用户数：根据日志中客户端 IP，5 分钟粒度去重统计
// + 日活跃用户数：根据日志中客户端 IP，按天粒度去重统计
func (c *Client) DescribeIpVisit(request *DescribeIpVisitRequest) (response *DescribeIpVisitResponse, err error) {
    if request == nil {
        request = NewDescribeIpVisitRequest()
    }
    response = NewDescribeIpVisitResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeMapInfoRequest() (request *DescribeMapInfoRequest) {
    request = &DescribeMapInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cdn", APIVersion, "DescribeMapInfo")
    return
}

func NewDescribeMapInfoResponse() (response *DescribeMapInfoResponse) {
    response = &DescribeMapInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeMapInfo 用于查询省份对应的 ID，运营商对应的 ID 信息。
func (c *Client) DescribeMapInfo(request *DescribeMapInfoRequest) (response *DescribeMapInfoResponse, err error) {
    if request == nil {
        request = NewDescribeMapInfoRequest()
    }
    response = NewDescribeMapInfoResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeOriginDataRequest() (request *DescribeOriginDataRequest) {
    request = &DescribeOriginDataRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cdn", APIVersion, "DescribeOriginData")
    return
}

func NewDescribeOriginDataResponse() (response *DescribeOriginDataResponse) {
    response = &DescribeOriginDataResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeOriginData 用于查询 CDN 实时回源监控数据，支持以下指标查询：
// 
// + 回源流量（单位为 byte）
// + 回源带宽（单位为 bps）
// + 回源请求数（单位为 次）
// + 回源失败请求数（单位为 次）
// + 回源失败率（单位为 %，小数点后保留两位）
// + 回源状态码 2XX 汇总及各 2 开头回源状态码明细（单位为 个）
// + 回源状态码 3XX 汇总及各 3 开头回源状态码明细（单位为 个）
// + 回源状态码 4XX 汇总及各 4 开头回源状态码明细（单位为 个）
// + 回源状态码 5XX 汇总及各 5 开头回源状态码明细（单位为 个）
func (c *Client) DescribeOriginData(request *DescribeOriginDataRequest) (response *DescribeOriginDataResponse, err error) {
    if request == nil {
        request = NewDescribeOriginDataRequest()
    }
    response = NewDescribeOriginDataResponse()
    err = c.Send(request, response)
    return
}

func NewDescribePayTypeRequest() (request *DescribePayTypeRequest) {
    request = &DescribePayTypeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cdn", APIVersion, "DescribePayType")
    return
}

func NewDescribePayTypeResponse() (response *DescribePayTypeResponse) {
    response = &DescribePayTypeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribePayType 用于查询用户的计费类型，计费周期等信息。
func (c *Client) DescribePayType(request *DescribePayTypeRequest) (response *DescribePayTypeResponse, err error) {
    if request == nil {
        request = NewDescribePayTypeRequest()
    }
    response = NewDescribePayTypeResponse()
    err = c.Send(request, response)
    return
}

func NewListTopDataRequest() (request *ListTopDataRequest) {
    request = &ListTopDataRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cdn", APIVersion, "ListTopData")
    return
}

func NewListTopDataResponse() (response *ListTopDataResponse) {
    response = &ListTopDataResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ListTopData 通过入参 Metric 和 Filter 组合不同，可以查询以下排序数据：
// 
// + 依据总流量、总请求数对访问 URL 排序，从大至小返回 TOP 1000 URL
// + 依据总流量、总请求数对客户端省份排序，从大至小返回省份列表
// + 依据总流量、总请求数对客户端运营商排序，从大至小返回运营商列表
// + 依据总流量、峰值带宽、总请求数、平均命中率、2XX/3XX/4XX/5XX 状态码对域名排序，从大至小返回域名列表
// + 依据总回源流量、回源峰值带宽、总回源请求数、平均回源失败率、2XX/3XX/4XX/5XX 回源状态码对域名排序，从大至小返回域名列表
func (c *Client) ListTopData(request *ListTopDataRequest) (response *ListTopDataResponse, err error) {
    if request == nil {
        request = NewListTopDataRequest()
    }
    response = NewListTopDataResponse()
    err = c.Send(request, response)
    return
}
