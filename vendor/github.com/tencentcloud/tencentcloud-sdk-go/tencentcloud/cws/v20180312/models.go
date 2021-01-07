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
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type CreateMonitorsRequest struct {
	*tchttp.BaseRequest

	// 站点的url列表
	Urls []*string `json:"Urls" name:"Urls" list`

	// 任务名称
	Name *string `json:"Name" name:"Name"`

	// 扫描模式，normal-正常扫描；deep-深度扫描
	ScannerType *string `json:"ScannerType" name:"ScannerType"`

	// 扫描周期，单位小时，每X小时执行一次
	Crontab *uint64 `json:"Crontab" name:"Crontab"`

	// 扫描速率限制，每秒发送X个HTTP请求
	RateLimit *uint64 `json:"RateLimit" name:"RateLimit"`

	// 首次扫描开始时间
	FirstScanStartTime *string `json:"FirstScanStartTime" name:"FirstScanStartTime"`
}

func (r *CreateMonitorsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateMonitorsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateMonitorsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateMonitorsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateMonitorsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateSitesRequest struct {
	*tchttp.BaseRequest

	// 站点的url列表
	Urls []*string `json:"Urls" name:"Urls" list`

	// 访问网站的客户端标识
	UserAgent *string `json:"UserAgent" name:"UserAgent"`
}

func (r *CreateSitesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateSitesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateSitesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 新增站点数。
		Number *uint64 `json:"Number" name:"Number"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateSitesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateSitesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateSitesScansRequest struct {
	*tchttp.BaseRequest

	// 站点的ID列表
	SiteIds []*uint64 `json:"SiteIds" name:"SiteIds" list`

	// 扫描模式，normal-正常扫描；deep-深度扫描
	ScannerType *string `json:"ScannerType" name:"ScannerType"`

	// 扫描速率限制，每秒发送X个HTTP请求
	RateLimit *uint64 `json:"RateLimit" name:"RateLimit"`
}

func (r *CreateSitesScansRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateSitesScansRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateSitesScansResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateSitesScansResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateSitesScansResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateVulsMisinformationRequest struct {
	*tchttp.BaseRequest

	// 漏洞ID列表
	VulIds []*uint64 `json:"VulIds" name:"VulIds" list`
}

func (r *CreateVulsMisinformationRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateVulsMisinformationRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateVulsMisinformationResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateVulsMisinformationResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateVulsMisinformationResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateVulsReportRequest struct {
	*tchttp.BaseRequest

	// 站点ID
	SiteId *uint64 `json:"SiteId" name:"SiteId"`

	// 监控任务ID
	MonitorId *uint64 `json:"MonitorId" name:"MonitorId"`
}

func (r *CreateVulsReportRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateVulsReportRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateVulsReportResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 报告下载地址
		ReportFileUrl *string `json:"ReportFileUrl" name:"ReportFileUrl"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateVulsReportResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateVulsReportResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteMonitorsRequest struct {
	*tchttp.BaseRequest

	// 监控任务ID列表
	MonitorIds []*uint64 `json:"MonitorIds" name:"MonitorIds" list`
}

func (r *DeleteMonitorsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteMonitorsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteMonitorsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteMonitorsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteMonitorsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteSitesRequest struct {
	*tchttp.BaseRequest

	// 站点ID列表
	SiteIds []*uint64 `json:"SiteIds" name:"SiteIds" list`
}

func (r *DeleteSitesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteSitesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteSitesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteSitesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteSitesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeConfigRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeConfigRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeConfigRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeConfigResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 漏洞告警通知等级，4位分别代表：高危、中危、低危、提示。
		NoticeLevel *string `json:"NoticeLevel" name:"NoticeLevel"`

		// 配置ID。
		Id *uint64 `json:"Id" name:"Id"`

		// 记录创建时间。
		CreatedAt *string `json:"CreatedAt" name:"CreatedAt"`

		// 记录更新新建。
		UpdatedAt *string `json:"UpdatedAt" name:"UpdatedAt"`

		// 云用户appid。
		Appid *uint64 `json:"Appid" name:"Appid"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeConfigResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeConfigResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMonitorsRequest struct {
	*tchttp.BaseRequest

	// 监控任务ID列表
	MonitorIds []*uint64 `json:"MonitorIds" name:"MonitorIds" list`

	// 过滤条件
	Filters []*Filter `json:"Filters" name:"Filters" list`

	// 偏移量，默认为0
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 返回数量，默认为10，最大值为100
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribeMonitorsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMonitorsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMonitorsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 监控任务列表。
		Monitors []*MonitorsDetail `json:"Monitors" name:"Monitors" list`

		// 监控任务数量。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeMonitorsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMonitorsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSiteQuotaRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeSiteQuotaRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSiteQuotaRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSiteQuotaResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 已购买的扫描次数。
		Total *uint64 `json:"Total" name:"Total"`

		// 已使用的扫描次数。
		Used *uint64 `json:"Used" name:"Used"`

		// 剩余可用的扫描次数。
		Available *uint64 `json:"Available" name:"Available"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeSiteQuotaResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSiteQuotaResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSitesRequest struct {
	*tchttp.BaseRequest

	// 站点ID列表
	SiteIds []*uint64 `json:"SiteIds" name:"SiteIds" list`

	// 过滤条件
	Filters []*Filter `json:"Filters" name:"Filters" list`

	// 偏移量，默认为0
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 返回数量，默认为10，最大值为100
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribeSitesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSitesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSitesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 站点数量。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 站点信息列表。
		Sites []*Site `json:"Sites" name:"Sites" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeSitesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSitesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSitesVerificationRequest struct {
	*tchttp.BaseRequest

	// 站点的url列表
	Urls []*string `json:"Urls" name:"Urls" list`
}

func (r *DescribeSitesVerificationRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSitesVerificationRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSitesVerificationResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 验证信息数量。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 验证信息列表。
		SitesVerification []*SitesVerification `json:"SitesVerification" name:"SitesVerification" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeSitesVerificationResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSitesVerificationResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeVulsNumberRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeVulsNumberRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeVulsNumberRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeVulsNumberResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 受影响的网站总数。
		ImpactSiteNumber *uint64 `json:"ImpactSiteNumber" name:"ImpactSiteNumber"`

		// 已验证的网站总数。
		SiteNumber *uint64 `json:"SiteNumber" name:"SiteNumber"`

		// 高风险漏洞总数。
		VulsHighNumber *uint64 `json:"VulsHighNumber" name:"VulsHighNumber"`

		// 中风险漏洞总数。
		VulsMiddleNumber *uint64 `json:"VulsMiddleNumber" name:"VulsMiddleNumber"`

		// 低高风险漏洞总数。
		VulsLowNumber *uint64 `json:"VulsLowNumber" name:"VulsLowNumber"`

		// 风险提示总数。
		VulsNoticeNumber *uint64 `json:"VulsNoticeNumber" name:"VulsNoticeNumber"`

		// 扫描页面总数。
		PageCount *uint64 `json:"PageCount" name:"PageCount"`

		// 已验证的网站列表。
		Sites []*MonitorMiniSite `json:"Sites" name:"Sites" list`

		// 受影响的网站列表。
		ImpactSites []*MonitorMiniSite `json:"ImpactSites" name:"ImpactSites" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeVulsNumberResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeVulsNumberResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeVulsNumberTimelineRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeVulsNumberTimelineRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeVulsNumberTimelineRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeVulsNumberTimelineResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 统计数据记录数量。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 用户漏洞数随时间变化统计数据。
		VulsTimeline []*VulsTimeline `json:"VulsTimeline" name:"VulsTimeline" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeVulsNumberTimelineResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeVulsNumberTimelineResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeVulsRequest struct {
	*tchttp.BaseRequest

	// 站点ID
	SiteId *uint64 `json:"SiteId" name:"SiteId"`

	// 监控任务ID
	MonitorId *uint64 `json:"MonitorId" name:"MonitorId"`

	// 过滤条件
	Filters []*Filter `json:"Filters" name:"Filters" list`

	// 偏移量，默认为0
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 返回数量，默认为10，最大值为100
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribeVulsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeVulsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeVulsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 漏洞数量。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 漏洞信息列表。
		Vuls []*Vul `json:"Vuls" name:"Vuls" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeVulsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeVulsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Filter struct {

	// 过滤键的名称。
	Name *string `json:"Name" name:"Name"`

	// 一个或者多个过滤值。
	Values []*string `json:"Values" name:"Values" list`
}

type ModifyConfigAttributeRequest struct {
	*tchttp.BaseRequest

	// 漏洞告警通知等级，4位分别代表：高危、中危、低危、提示
	NoticeLevel *string `json:"NoticeLevel" name:"NoticeLevel"`
}

func (r *ModifyConfigAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyConfigAttributeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyConfigAttributeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyConfigAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyConfigAttributeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyMonitorAttributeRequest struct {
	*tchttp.BaseRequest

	// 监测任务ID
	MonitorId *uint64 `json:"MonitorId" name:"MonitorId"`

	// 站点的url列表
	Urls []*string `json:"Urls" name:"Urls" list`

	// 任务名称
	Name *string `json:"Name" name:"Name"`

	// 扫描模式，normal-正常扫描；deep-深度扫描
	ScannerType *string `json:"ScannerType" name:"ScannerType"`

	// 扫描周期，单位小时，每X小时执行一次
	Crontab *uint64 `json:"Crontab" name:"Crontab"`

	// 扫描速率限制，每秒发送X个HTTP请求
	RateLimit *uint64 `json:"RateLimit" name:"RateLimit"`

	// 首次扫描开始时间
	FirstScanStartTime *string `json:"FirstScanStartTime" name:"FirstScanStartTime"`

	// 监测状态：1-监测中；2-暂停监测
	MonitorStatus *uint64 `json:"MonitorStatus" name:"MonitorStatus"`
}

func (r *ModifyMonitorAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyMonitorAttributeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyMonitorAttributeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyMonitorAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyMonitorAttributeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifySiteAttributeRequest struct {
	*tchttp.BaseRequest

	// 站点ID
	SiteId *uint64 `json:"SiteId" name:"SiteId"`

	// 站点名称
	Name *string `json:"Name" name:"Name"`

	// 网站是否需要登录扫描：0-未知；-1-不需要；1-需要
	NeedLogin *int64 `json:"NeedLogin" name:"NeedLogin"`

	// 登录后的cookie
	LoginCookie *string `json:"LoginCookie" name:"LoginCookie"`

	// 用于测试cookie是否有效的URL
	LoginCheckUrl *string `json:"LoginCheckUrl" name:"LoginCheckUrl"`

	// 用于测试cookie是否有效的关键字
	LoginCheckKw *string `json:"LoginCheckKw" name:"LoginCheckKw"`

	// 禁止扫描器扫描的目录关键字
	ScanDisallow *string `json:"ScanDisallow" name:"ScanDisallow"`
}

func (r *ModifySiteAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifySiteAttributeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifySiteAttributeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifySiteAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifySiteAttributeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Monitor struct {

	// 云用户appid。
	Appid *uint64 `json:"Appid" name:"Appid"`

	// 监控任务ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 监控名称。
	Name *string `json:"Name" name:"Name"`

	// 监测状态：1-监测中；2-暂停监测。
	MonitorStatus *uint64 `json:"MonitorStatus" name:"MonitorStatus"`

	// 监测模式，normal-正常扫描；deep-深度扫描。
	ScannerType *string `json:"ScannerType" name:"ScannerType"`

	// 扫描周期，单位小时，每X小时执行一次。
	Crontab *uint64 `json:"Crontab" name:"Crontab"`

	// 指定扫描类型，3位数每位依次表示：扫描Web漏洞、扫描系统漏洞、扫描系统端口。
	IncludedVulsTypes *string `json:"IncludedVulsTypes" name:"IncludedVulsTypes"`

	// 速率限制，每秒发送X个HTTP请求。
	RateLimit *uint64 `json:"RateLimit" name:"RateLimit"`

	// 首次扫描开始时间。
	FirstScanStartTime *string `json:"FirstScanStartTime" name:"FirstScanStartTime"`

	// 扫描状态：0-待扫描（无任何扫描结果）；1-扫描中（正在进行扫描）；2-已扫描（有扫描结果且不正在扫描）；3-扫描完成待同步结果。
	ScanStatus *uint64 `json:"ScanStatus" name:"ScanStatus"`

	// 上一次扫描完成时间。
	LastScanFinishTime *string `json:"LastScanFinishTime" name:"LastScanFinishTime"`

	// 当前扫描开始时间，如扫描完成则为上一次扫描的开始时间。
	CurrentScanStartTime *string `json:"CurrentScanStartTime" name:"CurrentScanStartTime"`

	// CreatedAt。
	CreatedAt *string `json:"CreatedAt" name:"CreatedAt"`

	// UpdatedAt。
	UpdatedAt *string `json:"UpdatedAt" name:"UpdatedAt"`
}

type MonitorMiniSite struct {

	// 站点ID。
	SiteId *uint64 `json:"SiteId" name:"SiteId"`

	// 站点Url。
	Url *string `json:"Url" name:"Url"`
}

type MonitorsDetail struct {

	// 监控任务包含的站点列表的平均扫描进度。
	Progress *uint64 `json:"Progress" name:"Progress"`

	// 扫描页面总数。
	PageCount *uint64 `json:"PageCount" name:"PageCount"`

	// 监控任务基础信息。
	Basic *Monitor `json:"Basic" name:"Basic"`

	// 监控任务包含的站点列表。
	Sites []*MonitorMiniSite `json:"Sites" name:"Sites" list`

	// 监控任务包含的站点列表数量。
	SiteNumber *uint64 `json:"SiteNumber" name:"SiteNumber"`

	// 监控任务包含的受漏洞威胁的站点列表。
	ImpactSites []*MonitorMiniSite `json:"ImpactSites" name:"ImpactSites" list`

	// 监控任务包含的受漏洞威胁的站点列表数量。
	ImpactSiteNumber *uint64 `json:"ImpactSiteNumber" name:"ImpactSiteNumber"`

	// 高风险漏洞数量。
	VulsHighNumber *uint64 `json:"VulsHighNumber" name:"VulsHighNumber"`

	// 中风险漏洞数量。
	VulsMiddleNumber *uint64 `json:"VulsMiddleNumber" name:"VulsMiddleNumber"`

	// 低风险漏洞数量。
	VulsLowNumber *uint64 `json:"VulsLowNumber" name:"VulsLowNumber"`

	// 提示数量。
	VulsNoticeNumber *uint64 `json:"VulsNoticeNumber" name:"VulsNoticeNumber"`
}

type Site struct {

	// 扫描进度，百分比整数
	Progress *uint64 `json:"Progress" name:"Progress"`

	// 云用户appid。
	Appid *uint64 `json:"Appid" name:"Appid"`

	// 云用户标识。
	Uin *string `json:"Uin" name:"Uin"`

	// 网站是否需要登录扫描：0-未知；-1-不需要；1-需要。
	NeedLogin *int64 `json:"NeedLogin" name:"NeedLogin"`

	// 登录后的cookie。
	LoginCookie *string `json:"LoginCookie" name:"LoginCookie"`

	// 登录后的cookie是否有效：0-无效；1-有效。
	LoginCookieValid *uint64 `json:"LoginCookieValid" name:"LoginCookieValid"`

	// 用于测试cookie是否有效的URL。
	LoginCheckUrl *string `json:"LoginCheckUrl" name:"LoginCheckUrl"`

	// 用于测试cookie是否有效的关键字。
	LoginCheckKw *string `json:"LoginCheckKw" name:"LoginCheckKw"`

	// 禁止扫描器扫描的目录关键字。
	ScanDisallow *string `json:"ScanDisallow" name:"ScanDisallow"`

	// 访问网站的客户端标识。
	UserAgent *string `json:"UserAgent" name:"UserAgent"`

	// 站点ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 监控任务ID，为0时表示未加入监控任务。
	MonitorId *uint64 `json:"MonitorId" name:"MonitorId"`

	// 站点url。
	Url *string `json:"Url" name:"Url"`

	// 站点名称。
	Name *string `json:"Name" name:"Name"`

	// 验证状态：0-未验证；1-已验证；2-验证失效，待重新验证。
	VerifyStatus *uint64 `json:"VerifyStatus" name:"VerifyStatus"`

	// 监测状态：0-未监测；1-监测中；2-暂停监测。
	MonitorStatus *uint64 `json:"MonitorStatus" name:"MonitorStatus"`

	// 扫描状态：0-待扫描（无任何扫描结果）；1-扫描中（正在进行扫描）；2-已扫描（有扫描结果且不正在扫描）；3-扫描完成待同步结果。
	ScanStatus *uint64 `json:"ScanStatus" name:"ScanStatus"`

	// 最近一次的AIScanner的扫描任务id，注意取消的情况。
	LastScanTaskId *uint64 `json:"LastScanTaskId" name:"LastScanTaskId"`

	// 最近一次扫描开始时间。
	LastScanStartTime *string `json:"LastScanStartTime" name:"LastScanStartTime"`

	// 最近一次扫描完成时间。
	LastScanFinishTime *string `json:"LastScanFinishTime" name:"LastScanFinishTime"`

	// 最近一次取消时间，取消即使用上一次扫描结果。
	LastScanCancelTime *string `json:"LastScanCancelTime" name:"LastScanCancelTime"`

	// 最近一次扫描扫描的页面数。
	LastScanPageCount *uint64 `json:"LastScanPageCount" name:"LastScanPageCount"`

	// normal-正常扫描；deep-深度扫描。
	LastScanScannerType *string `json:"LastScanScannerType" name:"LastScanScannerType"`

	// 最近一次扫描高风险漏洞数量。
	LastScanVulsHighNum *uint64 `json:"LastScanVulsHighNum" name:"LastScanVulsHighNum"`

	// 最近一次扫描中风险漏洞数量。
	LastScanVulsMiddleNum *uint64 `json:"LastScanVulsMiddleNum" name:"LastScanVulsMiddleNum"`

	// 最近一次扫描低风险漏洞数量。
	LastScanVulsLowNum *uint64 `json:"LastScanVulsLowNum" name:"LastScanVulsLowNum"`

	// 最近一次扫描提示信息数量。
	LastScanVulsNoticeNum *uint64 `json:"LastScanVulsNoticeNum" name:"LastScanVulsNoticeNum"`

	// 记录添加时间。
	CreatedAt *string `json:"CreatedAt" name:"CreatedAt"`

	// 记录最近修改时间。
	UpdatedAt *string `json:"UpdatedAt" name:"UpdatedAt"`

	// 速率限制，每秒发送X个HTTP请求。
	LastScanRateLimit *uint64 `json:"LastScanRateLimit" name:"LastScanRateLimit"`

	// 最近一次扫描漏洞总数量。
	LastScanVulsNum *uint64 `json:"LastScanVulsNum" name:"LastScanVulsNum"`

	// 最近一次扫描提示总数量
	LastScanNoticeNum *uint64 `json:"LastScanNoticeNum" name:"LastScanNoticeNum"`
}

type SitesVerification struct {

	// ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 云用户appid
	Appid *uint64 `json:"Appid" name:"Appid"`

	// 用于验证站点的url，即访问该url获取验证数据。
	VerifyUrl *string `json:"VerifyUrl" name:"VerifyUrl"`

	// 获取验证验证文件的url。
	VerifyFileUrl *string `json:"VerifyFileUrl" name:"VerifyFileUrl"`

	// 根域名。
	Domain *string `json:"Domain" name:"Domain"`

	// txt解析域名验证的name。
	TxtName *string `json:"TxtName" name:"TxtName"`

	// txt解析域名验证的text。
	TxtText *string `json:"TxtText" name:"TxtText"`

	// 验证有效期，在此之前有效。
	ValidTo *string `json:"ValidTo" name:"ValidTo"`

	// 验证状态：0-未验证；1-已验证；2-验证失效，待重新验证。
	VerifyStatus *uint64 `json:"VerifyStatus" name:"VerifyStatus"`

	// CreatedAt。
	CreatedAt *string `json:"CreatedAt" name:"CreatedAt"`

	// UpdatedAt。
	UpdatedAt *string `json:"UpdatedAt" name:"UpdatedAt"`
}

type VerifySitesRequest struct {
	*tchttp.BaseRequest

	// 站点的url列表
	Urls []*string `json:"Urls" name:"Urls" list`
}

func (r *VerifySitesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *VerifySitesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type VerifySitesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 验证成功的根域名数量。
		SuccessNumber *uint64 `json:"SuccessNumber" name:"SuccessNumber"`

		// 验证失败的根域名数量。
		FailNumber *uint64 `json:"FailNumber" name:"FailNumber"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *VerifySitesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *VerifySitesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Vul struct {

	// 是否已经添加误报，0-否，1-是。
	IsReported *uint64 `json:"IsReported" name:"IsReported"`

	// 云用户appid。
	Appid *uint64 `json:"Appid" name:"Appid"`

	// 云用户标识。
	Uin *string `json:"Uin" name:"Uin"`

	// 漏洞ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 站点ID。
	SiteId *uint64 `json:"SiteId" name:"SiteId"`

	// 扫描引擎的扫描任务ID。
	TaskId *uint64 `json:"TaskId" name:"TaskId"`

	// 漏洞级别：high、middle、low、notice。
	Level *string `json:"Level" name:"Level"`

	// 漏洞名称。
	Name *string `json:"Name" name:"Name"`

	// 出现漏洞的url。
	Url *string `json:"Url" name:"Url"`

	// 网址/细节。
	Html *string `json:"Html" name:"Html"`

	// 漏洞类型。
	Nickname *string `json:"Nickname" name:"Nickname"`

	// 危害说明。
	Harm *string `json:"Harm" name:"Harm"`

	// 漏洞描述。
	Describe *string `json:"Describe" name:"Describe"`

	// 解决方案。
	Solution *string `json:"Solution" name:"Solution"`

	// 漏洞参考。
	From *string `json:"From" name:"From"`

	// 漏洞通过该参数攻击。
	Parameter *string `json:"Parameter" name:"Parameter"`

	// CreatedAt。
	CreatedAt *string `json:"CreatedAt" name:"CreatedAt"`

	// UpdatedAt。
	UpdatedAt *string `json:"UpdatedAt" name:"UpdatedAt"`
}

type VulsTimeline struct {

	// ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 云用户appid。
	Appid *uint64 `json:"Appid" name:"Appid"`

	// 日期。
	Date *string `json:"Date" name:"Date"`

	// 扫描页面总数量。
	PageCount *uint64 `json:"PageCount" name:"PageCount"`

	// 已验证网站总数量。
	SiteNum *uint64 `json:"SiteNum" name:"SiteNum"`

	// 受影响的网站总数量。
	ImpactSiteNum *uint64 `json:"ImpactSiteNum" name:"ImpactSiteNum"`

	// 高危漏洞总数量。
	VulsHighNum *uint64 `json:"VulsHighNum" name:"VulsHighNum"`

	// 中危漏洞总数量。
	VulsMiddleNum *uint64 `json:"VulsMiddleNum" name:"VulsMiddleNum"`

	// 低危漏洞总数量。
	VulsLowNum *uint64 `json:"VulsLowNum" name:"VulsLowNum"`

	// 风险提示总数量
	VulsNoticeNum *uint64 `json:"VulsNoticeNum" name:"VulsNoticeNum"`

	// 记录添加时间。
	CreatedAt *string `json:"CreatedAt" name:"CreatedAt"`

	// 记录最近修改时间。
	UpdatedAt *string `json:"UpdatedAt" name:"UpdatedAt"`
}
