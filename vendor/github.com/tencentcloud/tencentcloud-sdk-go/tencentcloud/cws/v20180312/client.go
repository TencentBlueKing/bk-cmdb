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

package v20180312

import (
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2018-03-12"

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


func NewCreateMonitorsRequest() (request *CreateMonitorsRequest) {
    request = &CreateMonitorsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "CreateMonitors")
    return
}

func NewCreateMonitorsResponse() (response *CreateMonitorsResponse) {
    response = &CreateMonitorsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口（CreateMonitors）用于新增一个或多个站点的监测任务。
func (c *Client) CreateMonitors(request *CreateMonitorsRequest) (response *CreateMonitorsResponse, err error) {
    if request == nil {
        request = NewCreateMonitorsRequest()
    }
    response = NewCreateMonitorsResponse()
    err = c.Send(request, response)
    return
}

func NewCreateSitesRequest() (request *CreateSitesRequest) {
    request = &CreateSitesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "CreateSites")
    return
}

func NewCreateSitesResponse() (response *CreateSitesResponse) {
    response = &CreateSitesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口（CreateSites）用于新增一个或多个站点。
func (c *Client) CreateSites(request *CreateSitesRequest) (response *CreateSitesResponse, err error) {
    if request == nil {
        request = NewCreateSitesRequest()
    }
    response = NewCreateSitesResponse()
    err = c.Send(request, response)
    return
}

func NewCreateSitesScansRequest() (request *CreateSitesScansRequest) {
    request = &CreateSitesScansRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "CreateSitesScans")
    return
}

func NewCreateSitesScansResponse() (response *CreateSitesScansResponse) {
    response = &CreateSitesScansResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口（CreateSitesScans）用于新增一个或多个站点的单次扫描任务。
func (c *Client) CreateSitesScans(request *CreateSitesScansRequest) (response *CreateSitesScansResponse, err error) {
    if request == nil {
        request = NewCreateSitesScansRequest()
    }
    response = NewCreateSitesScansResponse()
    err = c.Send(request, response)
    return
}

func NewCreateVulsMisinformationRequest() (request *CreateVulsMisinformationRequest) {
    request = &CreateVulsMisinformationRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "CreateVulsMisinformation")
    return
}

func NewCreateVulsMisinformationResponse() (response *CreateVulsMisinformationResponse) {
    response = &CreateVulsMisinformationResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口（CreateVulsMisinformation）用于新增一个或多个漏洞误报信息。
func (c *Client) CreateVulsMisinformation(request *CreateVulsMisinformationRequest) (response *CreateVulsMisinformationResponse, err error) {
    if request == nil {
        request = NewCreateVulsMisinformationRequest()
    }
    response = NewCreateVulsMisinformationResponse()
    err = c.Send(request, response)
    return
}

func NewCreateVulsReportRequest() (request *CreateVulsReportRequest) {
    request = &CreateVulsReportRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "CreateVulsReport")
    return
}

func NewCreateVulsReportResponse() (response *CreateVulsReportResponse) {
    response = &CreateVulsReportResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (CreateVulsReport) 用于生成漏洞报告并返回下载链接。
func (c *Client) CreateVulsReport(request *CreateVulsReportRequest) (response *CreateVulsReportResponse, err error) {
    if request == nil {
        request = NewCreateVulsReportRequest()
    }
    response = NewCreateVulsReportResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteMonitorsRequest() (request *DeleteMonitorsRequest) {
    request = &DeleteMonitorsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "DeleteMonitors")
    return
}

func NewDeleteMonitorsResponse() (response *DeleteMonitorsResponse) {
    response = &DeleteMonitorsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (DeleteMonitors) 用于删除监控任务。
func (c *Client) DeleteMonitors(request *DeleteMonitorsRequest) (response *DeleteMonitorsResponse, err error) {
    if request == nil {
        request = NewDeleteMonitorsRequest()
    }
    response = NewDeleteMonitorsResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteSitesRequest() (request *DeleteSitesRequest) {
    request = &DeleteSitesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "DeleteSites")
    return
}

func NewDeleteSitesResponse() (response *DeleteSitesResponse) {
    response = &DeleteSitesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (DeleteSites) 用于删除站点。
func (c *Client) DeleteSites(request *DeleteSitesRequest) (response *DeleteSitesResponse, err error) {
    if request == nil {
        request = NewDeleteSitesRequest()
    }
    response = NewDeleteSitesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeConfigRequest() (request *DescribeConfigRequest) {
    request = &DescribeConfigRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "DescribeConfig")
    return
}

func NewDescribeConfigResponse() (response *DescribeConfigResponse) {
    response = &DescribeConfigResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (DescribeConfig) 用于查询用户配置的详细信息。
func (c *Client) DescribeConfig(request *DescribeConfigRequest) (response *DescribeConfigResponse, err error) {
    if request == nil {
        request = NewDescribeConfigRequest()
    }
    response = NewDescribeConfigResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeMonitorsRequest() (request *DescribeMonitorsRequest) {
    request = &DescribeMonitorsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "DescribeMonitors")
    return
}

func NewDescribeMonitorsResponse() (response *DescribeMonitorsResponse) {
    response = &DescribeMonitorsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (DescribeMonitors) 用于查询一个或多个监控任务的详细信息。
func (c *Client) DescribeMonitors(request *DescribeMonitorsRequest) (response *DescribeMonitorsResponse, err error) {
    if request == nil {
        request = NewDescribeMonitorsRequest()
    }
    response = NewDescribeMonitorsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeSiteQuotaRequest() (request *DescribeSiteQuotaRequest) {
    request = &DescribeSiteQuotaRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "DescribeSiteQuota")
    return
}

func NewDescribeSiteQuotaResponse() (response *DescribeSiteQuotaResponse) {
    response = &DescribeSiteQuotaResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (DescribeSiteQuota) 用于查询用户购买的扫描次数总数和已使用数。
func (c *Client) DescribeSiteQuota(request *DescribeSiteQuotaRequest) (response *DescribeSiteQuotaResponse, err error) {
    if request == nil {
        request = NewDescribeSiteQuotaRequest()
    }
    response = NewDescribeSiteQuotaResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeSitesRequest() (request *DescribeSitesRequest) {
    request = &DescribeSitesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "DescribeSites")
    return
}

func NewDescribeSitesResponse() (response *DescribeSitesResponse) {
    response = &DescribeSitesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (DescribeSites) 用于查询一个或多个站点的详细信息。
func (c *Client) DescribeSites(request *DescribeSitesRequest) (response *DescribeSitesResponse, err error) {
    if request == nil {
        request = NewDescribeSitesRequest()
    }
    response = NewDescribeSitesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeSitesVerificationRequest() (request *DescribeSitesVerificationRequest) {
    request = &DescribeSitesVerificationRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "DescribeSitesVerification")
    return
}

func NewDescribeSitesVerificationResponse() (response *DescribeSitesVerificationResponse) {
    response = &DescribeSitesVerificationResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (DescribeSitesVerification) 用于查询一个或多个待验证站点的验证信息。
func (c *Client) DescribeSitesVerification(request *DescribeSitesVerificationRequest) (response *DescribeSitesVerificationResponse, err error) {
    if request == nil {
        request = NewDescribeSitesVerificationRequest()
    }
    response = NewDescribeSitesVerificationResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeVulsRequest() (request *DescribeVulsRequest) {
    request = &DescribeVulsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "DescribeVuls")
    return
}

func NewDescribeVulsResponse() (response *DescribeVulsResponse) {
    response = &DescribeVulsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (DescribeVuls) 用于查询一个或多个漏洞的详细信息。
func (c *Client) DescribeVuls(request *DescribeVulsRequest) (response *DescribeVulsResponse, err error) {
    if request == nil {
        request = NewDescribeVulsRequest()
    }
    response = NewDescribeVulsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeVulsNumberRequest() (request *DescribeVulsNumberRequest) {
    request = &DescribeVulsNumberRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "DescribeVulsNumber")
    return
}

func NewDescribeVulsNumberResponse() (response *DescribeVulsNumberResponse) {
    response = &DescribeVulsNumberResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (DescribeVulsNumber) 用于查询用户网站的漏洞总计数量。
func (c *Client) DescribeVulsNumber(request *DescribeVulsNumberRequest) (response *DescribeVulsNumberResponse, err error) {
    if request == nil {
        request = NewDescribeVulsNumberRequest()
    }
    response = NewDescribeVulsNumberResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeVulsNumberTimelineRequest() (request *DescribeVulsNumberTimelineRequest) {
    request = &DescribeVulsNumberTimelineRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "DescribeVulsNumberTimeline")
    return
}

func NewDescribeVulsNumberTimelineResponse() (response *DescribeVulsNumberTimelineResponse) {
    response = &DescribeVulsNumberTimelineResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (DescribeVulsNumberTimeline) 用于查询漏洞数随时间变化统计信息。
func (c *Client) DescribeVulsNumberTimeline(request *DescribeVulsNumberTimelineRequest) (response *DescribeVulsNumberTimelineResponse, err error) {
    if request == nil {
        request = NewDescribeVulsNumberTimelineRequest()
    }
    response = NewDescribeVulsNumberTimelineResponse()
    err = c.Send(request, response)
    return
}

func NewModifyConfigAttributeRequest() (request *ModifyConfigAttributeRequest) {
    request = &ModifyConfigAttributeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "ModifyConfigAttribute")
    return
}

func NewModifyConfigAttributeResponse() (response *ModifyConfigAttributeResponse) {
    response = &ModifyConfigAttributeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (ModifyConfigAttribute) 用于修改用户配置的属性。
func (c *Client) ModifyConfigAttribute(request *ModifyConfigAttributeRequest) (response *ModifyConfigAttributeResponse, err error) {
    if request == nil {
        request = NewModifyConfigAttributeRequest()
    }
    response = NewModifyConfigAttributeResponse()
    err = c.Send(request, response)
    return
}

func NewModifyMonitorAttributeRequest() (request *ModifyMonitorAttributeRequest) {
    request = &ModifyMonitorAttributeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "ModifyMonitorAttribute")
    return
}

func NewModifyMonitorAttributeResponse() (response *ModifyMonitorAttributeResponse) {
    response = &ModifyMonitorAttributeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (ModifyMonitorAttribute) 用于修改监测任务的属性。
func (c *Client) ModifyMonitorAttribute(request *ModifyMonitorAttributeRequest) (response *ModifyMonitorAttributeResponse, err error) {
    if request == nil {
        request = NewModifyMonitorAttributeRequest()
    }
    response = NewModifyMonitorAttributeResponse()
    err = c.Send(request, response)
    return
}

func NewModifySiteAttributeRequest() (request *ModifySiteAttributeRequest) {
    request = &ModifySiteAttributeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "ModifySiteAttribute")
    return
}

func NewModifySiteAttributeResponse() (response *ModifySiteAttributeResponse) {
    response = &ModifySiteAttributeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (ModifySiteAttribute) 用于修改站点的属性。
func (c *Client) ModifySiteAttribute(request *ModifySiteAttributeRequest) (response *ModifySiteAttributeResponse, err error) {
    if request == nil {
        request = NewModifySiteAttributeRequest()
    }
    response = NewModifySiteAttributeResponse()
    err = c.Send(request, response)
    return
}

func NewVerifySitesRequest() (request *VerifySitesRequest) {
    request = &VerifySitesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("cws", APIVersion, "VerifySites")
    return
}

func NewVerifySitesResponse() (response *VerifySitesResponse) {
    response = &VerifySitesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口 (VerifySites) 用于验证一个或多个待验证站点。
func (c *Client) VerifySites(request *VerifySitesRequest) (response *VerifySitesResponse, err error) {
    if request == nil {
        request = NewVerifySitesRequest()
    }
    response = NewVerifySitesResponse()
    err = c.Send(request, response)
    return
}
