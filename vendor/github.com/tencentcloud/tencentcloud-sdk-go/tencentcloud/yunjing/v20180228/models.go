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

package v20180228

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type Account struct {

	// 唯一ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 云镜客户端唯一Uuid
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 主机内网IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 主机名称。
	MachineName *string `json:"MachineName" name:"MachineName"`

	// 帐号名。
	Username *string `json:"Username" name:"Username"`

	// 帐号所属组。
	Groups *string `json:"Groups" name:"Groups"`

	// 帐号类型。
	// <li>ORDINARY：普通帐号</li>
	// <li>SUPPER：超级管理员帐号</li>
	Privilege *string `json:"Privilege" name:"Privilege"`

	// 帐号创建时间。
	AccountCreateTime *string `json:"AccountCreateTime" name:"AccountCreateTime"`

	// 帐号最后登录时间。
	LastLoginTime *string `json:"LastLoginTime" name:"LastLoginTime"`
}

type AccountStatistics struct {

	// 用户名。
	Username *string `json:"Username" name:"Username"`

	// 主机数量。
	MachineNum *uint64 `json:"MachineNum" name:"MachineNum"`
}

type AgentVul struct {

	// 漏洞ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 主机IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 漏洞名称。
	VulName *string `json:"VulName" name:"VulName"`

	// 漏洞危害等级。
	// <li>HIGH：高危</li>
	// <li>MIDDLE：中危</li>
	// <li>LOW：低危</li>
	// <li>NOTICE：提示</li>
	VulLevel *string `json:"VulLevel" name:"VulLevel"`

	// 最后扫描时间。
	LastScanTime *string `json:"LastScanTime" name:"LastScanTime"`

	// 漏洞描述。
	Description *string `json:"Description" name:"Description"`

	// 漏洞种类ID。
	VulId *uint64 `json:"VulId" name:"VulId"`

	// 漏洞状态。
	// <li>UN_OPERATED : 待处理</li>
	// <li>FIXED : 已修复</li>
	VulStatus *string `json:"VulStatus" name:"VulStatus"`
}

type BruteAttack struct {

	// 事件ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 主机IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 破解事件状态
	// <li>BRUTEATTACK_FAIL_ACCOUNT： 暴力破解事件-失败(存在帐号)  </li>
	// <li>BRUTEATTACK_FAIL_NOACCOUNT：暴力破解事件-失败(帐号不存在)</li>
	// <li>BRUTEATTACK_SUCCESS：暴力破解事件-成功</li>
	Status *string `json:"Status" name:"Status"`

	// 用户名称。
	UserName *string `json:"UserName" name:"UserName"`

	// 城市ID。
	City *uint64 `json:"City" name:"City"`

	// 国家ID。
	Country *uint64 `json:"Country" name:"Country"`

	// 省份ID。
	Province *uint64 `json:"Province" name:"Province"`

	// 来源IP。
	SrcIp *string `json:"SrcIp" name:"SrcIp"`

	// 尝试破解次数。
	Count *uint64 `json:"Count" name:"Count"`

	// 发生时间。
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 主机名称。
	MachineName *string `json:"MachineName" name:"MachineName"`

	// 云镜客户端唯一标识UUID。
	Uuid *string `json:"Uuid" name:"Uuid"`
}

type ChargePrepaid struct {

	// 购买实例的时长，单位：月。取值范围：1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36。
	Period *uint64 `json:"Period" name:"Period"`

	// 自动续费标识。取值范围：
	// <li>NOTIFY_AND_AUTO_RENEW：通知过期且自动续费</li>
	// <li>NOTIFY_AND_MANUAL_RENEW：通知过期不自动续费</li>
	// <li>DISABLE_NOTIFY_AND_MANUAL_RENEW：不通知过期不自动续费</li>
	// 
	// 默认取值：NOTIFY_AND_MANUAL_RENEW。若该参数指定为NOTIFY_AND_AUTO_RENEW，在账户余额充足的情况下，实例到期后将按月自动续费。
	RenewFlag *string `json:"RenewFlag" name:"RenewFlag"`
}

type CloseProVersionRequest struct {
	*tchttp.BaseRequest

	// 主机唯一标识Uuid。
	// 黑石的InstanceId，CVM的Uuid
	Quuid *string `json:"Quuid" name:"Quuid"`
}

func (r *CloseProVersionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CloseProVersionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CloseProVersionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CloseProVersionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CloseProVersionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Component struct {

	// 唯一ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 云镜客户端唯一Uuid。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 主机内网IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 主机名。
	MachineName *string `json:"MachineName" name:"MachineName"`

	// 组件版本号。
	ComponentVersion *string `json:"ComponentVersion" name:"ComponentVersion"`

	// 组件类型。
	// <li>SYSTEM：系统组件</li>
	// <li>WEB：WEB组件</li>
	ComponentType *string `json:"ComponentType" name:"ComponentType"`

	// 组件名称。
	ComponentName *string `json:"ComponentName" name:"ComponentName"`

	// 组件检测更新时间。
	ModifyTime *string `json:"ModifyTime" name:"ModifyTime"`
}

type ComponentStatistics struct {

	// 组件ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 主机数量。
	MachineNum *uint64 `json:"MachineNum" name:"MachineNum"`

	// 组件名称。
	ComponentName *string `json:"ComponentName" name:"ComponentName"`

	// 组件类型。
	// <li>WEB：web组件</li>
	// <li>SYSTEM：系统组件</li>
	ComponentType *string `json:"ComponentType" name:"ComponentType"`

	// 组件描述。
	Description *string `json:"Description" name:"Description"`
}

type CreateProcessTaskRequest struct {
	*tchttp.BaseRequest

	// 云镜客户端唯一Uuid。
	Uuid *string `json:"Uuid" name:"Uuid"`
}

func (r *CreateProcessTaskRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateProcessTaskRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateProcessTaskResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateProcessTaskResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateProcessTaskResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateUsualLoginPlacesRequest struct {
	*tchttp.BaseRequest

	// 云镜客户端UUID数组。
	Uuids []*string `json:"Uuids" name:"Uuids" list`

	// 登录地域信息数组。
	Places []*Place `json:"Places" name:"Places" list`
}

func (r *CreateUsualLoginPlacesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateUsualLoginPlacesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateUsualLoginPlacesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateUsualLoginPlacesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateUsualLoginPlacesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteBruteAttacksRequest struct {
	*tchttp.BaseRequest

	// 暴力破解事件Id数组。
	Ids []*uint64 `json:"Ids" name:"Ids" list`
}

func (r *DeleteBruteAttacksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteBruteAttacksRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteBruteAttacksResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteBruteAttacksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteBruteAttacksResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteMachineRequest struct {
	*tchttp.BaseRequest

	// 云镜客户端Uuid。
	Uuid *string `json:"Uuid" name:"Uuid"`
}

func (r *DeleteMachineRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteMachineRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteMachineResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteMachineResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteMachineResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteMaliciousRequestsRequest struct {
	*tchttp.BaseRequest

	// 恶意请求记录ID数组，最大100条。
	Ids []*uint64 `json:"Ids" name:"Ids" list`
}

func (r *DeleteMaliciousRequestsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteMaliciousRequestsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteMaliciousRequestsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteMaliciousRequestsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteMaliciousRequestsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteMalwaresRequest struct {
	*tchttp.BaseRequest

	// 木马记录ID数组
	Ids []*uint64 `json:"Ids" name:"Ids" list`
}

func (r *DeleteMalwaresRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteMalwaresRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteMalwaresResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteMalwaresResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteMalwaresResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteNonlocalLoginPlacesRequest struct {
	*tchttp.BaseRequest

	// 异地登录事件Id数组。
	Ids []*uint64 `json:"Ids" name:"Ids" list`
}

func (r *DeleteNonlocalLoginPlacesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteNonlocalLoginPlacesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteNonlocalLoginPlacesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteNonlocalLoginPlacesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteNonlocalLoginPlacesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteUsualLoginPlacesRequest struct {
	*tchttp.BaseRequest

	// 云镜客户端Uuid
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 已添加常用登录地城市ID数组
	CityIds []*uint64 `json:"CityIds" name:"CityIds" list`
}

func (r *DeleteUsualLoginPlacesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteUsualLoginPlacesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteUsualLoginPlacesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteUsualLoginPlacesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteUsualLoginPlacesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAccountStatisticsRequest struct {
	*tchttp.BaseRequest

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>Username - String - 是否必填：否 - 帐号用户名</li>
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeAccountStatisticsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAccountStatisticsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAccountStatisticsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 帐号统计列表记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 帐号统计列表。
		AccountStatistics []*AccountStatistics `json:"AccountStatistics" name:"AccountStatistics" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeAccountStatisticsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAccountStatisticsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAccountsRequest struct {
	*tchttp.BaseRequest

	// 云镜客户端唯一Uuid。Username和Uuid必填其一，使用Uuid表示，查询该主机下列表信息。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 云镜客户端唯一Uuid。Username和Uuid必填其一，使用Username表示，查询该用户名下列表信息。
	Username *string `json:"Username" name:"Username"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>Username - String - 是否必填：否 - 帐号名</li>
	// <li>Privilege - String - 是否必填：否 - 帐号类型（ORDINARY: 普通帐号 | SUPPER: 超级管理员帐号）</li>
	// <li>MachineIp - String - 是否必填：否 - 主机内网IP</li>
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeAccountsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAccountsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAccountsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 帐号列表记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 帐号数据列表。
		Accounts []*Account `json:"Accounts" name:"Accounts" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeAccountsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAccountsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAgentVulsRequest struct {
	*tchttp.BaseRequest

	// 漏洞类型。
	// <li>WEB: Web应用漏洞</li>
	// <li>SYSTEM：系统组件漏洞</li>
	// <li>BASELINE：安全基线</li>
	VulType *string `json:"VulType" name:"VulType"`

	// 客户端UUID。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>Status - String - 是否必填：否 - 状态筛选（UN_OPERATED: 待处理 | FIXED：已修复）
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeAgentVulsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAgentVulsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAgentVulsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 记录总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 主机漏洞信息
		AgentVuls []*AgentVul `json:"AgentVuls" name:"AgentVuls" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeAgentVulsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAgentVulsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAlarmAttributeRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeAlarmAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAlarmAttributeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAlarmAttributeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 防护软件离线告警状态：
	// <li>OPEN：告警已开启</li>
	// <li>CLOSE： 告警已关闭</li>
		Offline *string `json:"Offline" name:"Offline"`

		// 发现木马告警状态：
	// <li>OPEN：告警已开启</li>
	// <li>CLOSE： 告警已关闭</li>
		Malware *string `json:"Malware" name:"Malware"`

		// 发现异地登录告警状态：
	// <li>OPEN：告警已开启</li>
	// <li>CLOSE： 告警已关闭</li>
		NonlocalLogin *string `json:"NonlocalLogin" name:"NonlocalLogin"`

		// 被暴力破解成功告警状态：
	// <li>OPEN：告警已开启</li>
	// <li>CLOSE： 告警已关闭</li>
		CrackSuccess *string `json:"CrackSuccess" name:"CrackSuccess"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeAlarmAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAlarmAttributeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBruteAttacksRequest struct {
	*tchttp.BaseRequest

	// 客户端唯一Uuid。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>Keywords - String - 是否必填：否 -  查询关键字</li>
	// <li>Status - String - 是否必填：否 -  查询状态（FAILED：破解失败 |SUCCESS：破解成功）</li>
	Filters []*Filter `json:"Filters" name:"Filters" list`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribeBruteAttacksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBruteAttacksRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBruteAttacksResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 事件数量
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 暴力破解事件列表
		BruteAttacks []*BruteAttack `json:"BruteAttacks" name:"BruteAttacks" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBruteAttacksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBruteAttacksResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeComponentInfoRequest struct {
	*tchttp.BaseRequest

	// 组件ID。
	ComponentId *uint64 `json:"ComponentId" name:"ComponentId"`
}

func (r *DescribeComponentInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeComponentInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeComponentInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 组件ID。
		Id *uint64 `json:"Id" name:"Id"`

		// 组件名称。
		ComponentName *string `json:"ComponentName" name:"ComponentName"`

		// 组件类型。
	// <li>WEB：web组件</li>
	// <li>SYSTEM：系统组件</li>
		ComponentType *string `json:"ComponentType" name:"ComponentType"`

		// 组件官网。
		Homepage *string `json:"Homepage" name:"Homepage"`

		// 组件描述。
		Description *string `json:"Description" name:"Description"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeComponentInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeComponentInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeComponentStatisticsRequest struct {
	*tchttp.BaseRequest

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// ComponentName - String - 是否必填：否 - 组件名称
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeComponentStatisticsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeComponentStatisticsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeComponentStatisticsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 组件统计列表记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 组件统计列表数据数组。
		ComponentStatistics []*ComponentStatistics `json:"ComponentStatistics" name:"ComponentStatistics" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeComponentStatisticsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeComponentStatisticsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeComponentsRequest struct {
	*tchttp.BaseRequest

	// 云镜客户端唯一Uuid。Uuid和ComponentId必填其一，使用Uuid表示，查询该主机列表信息。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 组件ID。Uuid和ComponentId必填其一，使用ComponentId表示，查询该组件列表信息。
	ComponentId *uint64 `json:"ComponentId" name:"ComponentId"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>ComponentVersion - String - 是否必填：否 - 组件版本号</li>
	// <li>MachineIp - String - 是否必填：否 - 主机内网IP</li>
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeComponentsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeComponentsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeComponentsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 组件列表记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 组件列表数据。
		Components []*Component `json:"Components" name:"Components" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeComponentsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeComponentsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeHistoryAccountsRequest struct {
	*tchttp.BaseRequest

	// 云镜客户端唯一Uuid。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>Username - String - 是否必填：否 - 帐号名</li>
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeHistoryAccountsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeHistoryAccountsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeHistoryAccountsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 帐号变更历史列表记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 帐号变更历史数据数组。
		HistoryAccounts []*HistoryAccount `json:"HistoryAccounts" name:"HistoryAccounts" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeHistoryAccountsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeHistoryAccountsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeImpactedHostsRequest struct {
	*tchttp.BaseRequest

	// 漏洞种类ID。
	VulId *uint64 `json:"VulId" name:"VulId"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>Status - String - 是否必填：否 - 状态筛选（UN_OPERATED：待处理 | FIXED：已修复）</li>
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeImpactedHostsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeImpactedHostsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeImpactedHostsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 记录总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 漏洞影响机器列表数组
		ImpactedHosts []*ImpactedHost `json:"ImpactedHosts" name:"ImpactedHosts" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeImpactedHostsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeImpactedHostsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMachineInfoRequest struct {
	*tchttp.BaseRequest

	// 云镜客户端唯一Uuid。
	Uuid *string `json:"Uuid" name:"Uuid"`
}

func (r *DescribeMachineInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMachineInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMachineInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 机器ip。
		MachineIp *string `json:"MachineIp" name:"MachineIp"`

		// 受云镜保护天数。
		ProtectDays *uint64 `json:"ProtectDays" name:"ProtectDays"`

		// 操作系统。
		MachineOs *string `json:"MachineOs" name:"MachineOs"`

		// 主机名称。
		MachineName *string `json:"MachineName" name:"MachineName"`

		// 在线状态。
	// <li>ONLINE： 在线</li>
	// <li>OFFLINE：离线</li>
		MachineStatus *string `json:"MachineStatus" name:"MachineStatus"`

		// CVM或BM主机唯一标识。
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 主机外网IP。
		MachineWanIp *string `json:"MachineWanIp" name:"MachineWanIp"`

		// CVM或BM主机唯一Uuid。
		Quuid *string `json:"Quuid" name:"Quuid"`

		// 云镜客户端唯一Uuid。
		Uuid *string `json:"Uuid" name:"Uuid"`

		// 是否开通专业版。
	// <li>true：是</li>
	// <li>false：否</li>
		IsProVersion *bool `json:"IsProVersion" name:"IsProVersion"`

		// 专业版开通时间。
		ProVersionOpenDate *string `json:"ProVersionOpenDate" name:"ProVersionOpenDate"`

		// 云主机类型。
	// <li>CVM: 虚拟主机</li>
	// <li>BM: 黑石物理机</li>
		MachineType *string `json:"MachineType" name:"MachineType"`

		// 机器所属地域。如：ap-guangzhou，ap-shanghai
		MachineRegion *string `json:"MachineRegion" name:"MachineRegion"`

		// 主机状态。
	// <li>POSTPAY: 表示后付费，即按量计费  </li>
	// <li>PREPAY: 表示预付费，即包年包月</li>
		PayMode *string `json:"PayMode" name:"PayMode"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeMachineInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMachineInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMachinesRequest struct {
	*tchttp.BaseRequest

	// 云主机类型。
	// <li>CVM：表示虚拟主机</li>
	// <li>BM:  表示黑石物理机</li>
	MachineType *string `json:"MachineType" name:"MachineType"`

	// 机器所属地域。如：ap-guangzhou，ap-shanghai
	MachineRegion *string `json:"MachineRegion" name:"MachineRegion"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>Keywords - String - 是否必填：否 - 查询关键字 </li>
	// <li>Status - String - 是否必填：否 - 客户端在线状态（OFFLINE: 离线 | ONLINE: 在线）</li>
	// <li>Version - String  是否必填：否 - 当前防护版本（ PRO_VERSION：专业版 | BASIC_VERSION：基础版）</li>
	// 每个过滤条件只支持一个值，暂不支持多个值“或”关系查询
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeMachinesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMachinesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMachinesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 主机列表
		Machines []*Machine `json:"Machines" name:"Machines" list`

		// 主机数量
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeMachinesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMachinesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMaliciousRequestsRequest struct {
	*tchttp.BaseRequest

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>Status - String - 是否必填：否 - 状态筛选（UN_OPERATED: 待处理 | TRUSTED：已信任 | UN_TRUSTED：已取消信任）</li>
	// <li>Domain - String - 是否必填：否 - 恶意请求的域名</li>
	// <li>MachineIp - String - 是否必填：否 - 主机内网IP</li>
	Filters []*Filter `json:"Filters" name:"Filters" list`

	// 云镜客户端唯一UUID。
	Uuid *string `json:"Uuid" name:"Uuid"`
}

func (r *DescribeMaliciousRequestsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMaliciousRequestsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMaliciousRequestsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 恶意请求记录数组。
		MaliciousRequests []*MaliciousRequest `json:"MaliciousRequests" name:"MaliciousRequests" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeMaliciousRequestsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMaliciousRequestsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMalwaresRequest struct {
	*tchttp.BaseRequest

	// 客户端唯一Uuid。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>Keywords - String - 是否必填：否 - 查询关键字 </li>
	// <li>Status - String - 是否必填：否 - 木马状态（UN_OPERATED: 未处理 | SEGREGATED: 已隔离|TRUSTED：信任）</li>
	// 每个过滤条件只支持一个值，暂不支持多个值“或”关系查询。
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeMalwaresRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMalwaresRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMalwaresResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 木马总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// Malware数组。
		Malwares []*Malware `json:"Malwares" name:"Malwares" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeMalwaresResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMalwaresResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeNonlocalLoginPlacesRequest struct {
	*tchttp.BaseRequest

	// 客户端唯一Uuid。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>Keywords - String - 是否必填：否 -  查询关键字</li>
	// <li>Status - String - 是否必填：否 -  登录状态（NON_LOCAL_LOGIN: 异地登录 | NORMAL_LOGIN : 正常登录）</li>
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeNonlocalLoginPlacesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeNonlocalLoginPlacesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeNonlocalLoginPlacesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 异地登录信息数组。
		NonLocalLoginPlaces []*NonLocalLoginPlace `json:"NonLocalLoginPlaces" name:"NonLocalLoginPlaces" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeNonlocalLoginPlacesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeNonlocalLoginPlacesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeOpenPortStatisticsRequest struct {
	*tchttp.BaseRequest

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>Port - Uint64 - 是否必填：否 - 端口号</li>
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeOpenPortStatisticsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeOpenPortStatisticsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeOpenPortStatisticsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 端口统计列表总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 端口统计数据列表
		OpenPortStatistics []*OpenPortStatistics `json:"OpenPortStatistics" name:"OpenPortStatistics" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeOpenPortStatisticsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeOpenPortStatisticsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeOpenPortsRequest struct {
	*tchttp.BaseRequest

	// 云镜客户端唯一Uuid。Port和Uuid必填其一，使用Uuid表示，查询该主机列表信息。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 开放端口号。Port和Uuid必填其一，使用Port表示查询该端口的列表信息。
	Port *uint64 `json:"Port" name:"Port"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>Port - Uint64 - 是否必填：否 - 端口号</li>
	// <li>ProcessName - String - 是否必填：否 - 进程名</li>
	// <li>MachineIp - String - 是否必填：否 - 主机内网IP</li>
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeOpenPortsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeOpenPortsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeOpenPortsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 端口列表记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 端口列表。
		OpenPorts []*OpenPort `json:"OpenPorts" name:"OpenPorts" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeOpenPortsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeOpenPortsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeOverviewStatisticsRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeOverviewStatisticsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeOverviewStatisticsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeOverviewStatisticsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 服务器在线数。
		OnlineMachineNum *uint64 `json:"OnlineMachineNum" name:"OnlineMachineNum"`

		// 专业服务器数。
		ProVersionMachineNum *uint64 `json:"ProVersionMachineNum" name:"ProVersionMachineNum"`

		// 木马文件数。
		MalwareNum *uint64 `json:"MalwareNum" name:"MalwareNum"`

		// 异地登录数。
		NonlocalLoginNum *uint64 `json:"NonlocalLoginNum" name:"NonlocalLoginNum"`

		// 暴力破解成功数。
		BruteAttackSuccessNum *uint64 `json:"BruteAttackSuccessNum" name:"BruteAttackSuccessNum"`

		// 漏洞数。
		VulNum *uint64 `json:"VulNum" name:"VulNum"`

		// 安全基线数。
		BaseLineNum *uint64 `json:"BaseLineNum" name:"BaseLineNum"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeOverviewStatisticsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeOverviewStatisticsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProVersionInfoRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeProVersionInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProVersionInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProVersionInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 后付费昨日扣费
		PostPayCost *uint64 `json:"PostPayCost" name:"PostPayCost"`

		// 新增主机是否自动开通专业版
		IsAutoOpenProVersion *bool `json:"IsAutoOpenProVersion" name:"IsAutoOpenProVersion"`

		// 开通专业版主机数
		ProVersionNum *uint64 `json:"ProVersionNum" name:"ProVersionNum"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeProVersionInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProVersionInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProcessStatisticsRequest struct {
	*tchttp.BaseRequest

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>ProcessName - String - 是否必填：否 - 进程名</li>
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeProcessStatisticsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProcessStatisticsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProcessStatisticsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 进程统计列表记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 进程统计列表数据数组。
		ProcessStatistics []*ProcessStatistics `json:"ProcessStatistics" name:"ProcessStatistics" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeProcessStatisticsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProcessStatisticsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProcessTaskStatusRequest struct {
	*tchttp.BaseRequest

	// 云镜客户端唯一Uuid。
	Uuid *string `json:"Uuid" name:"Uuid"`
}

func (r *DescribeProcessTaskStatusRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProcessTaskStatusRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProcessTaskStatusResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务状态。
	// <li>COMPLETE：完成（此时可以调用DescribeProcesses接口获取实时进程列表）</li>
	// <li>AGENT_OFFLINE：云镜客户端离线</li>
	// <li>COLLECTING：进程获取中</li>
	// <li>FAILED：进程获取失败</li>
		Status *string `json:"Status" name:"Status"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeProcessTaskStatusResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProcessTaskStatusResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProcessesRequest struct {
	*tchttp.BaseRequest

	// 云镜客户端唯一Uuid。Uuid和ProcessName必填其一，使用Uuid表示，查询该主机列表信息。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 进程名。Uuid和ProcessName必填其一，使用ProcessName表示，查询该进程列表信息。
	ProcessName *string `json:"ProcessName" name:"ProcessName"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>ProcessName - String - 是否必填：否 - 进程名</li>
	// <li>MachineIp - String - 是否必填：否 - 主机内网IP</li>
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeProcessesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProcessesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProcessesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 进程列表记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 进程列表数据数组。
		Processes []*Process `json:"Processes" name:"Processes" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeProcessesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProcessesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSecurityDynamicsRequest struct {
	*tchttp.BaseRequest

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`
}

func (r *DescribeSecurityDynamicsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSecurityDynamicsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSecurityDynamicsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 安全事件消息数组。
		SecurityDynamics []*SecurityDynamic `json:"SecurityDynamics" name:"SecurityDynamics" list`

		// 记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeSecurityDynamicsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSecurityDynamicsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSecurityTrendsRequest struct {
	*tchttp.BaseRequest

	// 开始时间。
	BeginDate *string `json:"BeginDate" name:"BeginDate"`

	// 结束时间。
	EndDate *string `json:"EndDate" name:"EndDate"`
}

func (r *DescribeSecurityTrendsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSecurityTrendsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSecurityTrendsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 木马事件统计数据数组。
		Malwares []*SecurityTrend `json:"Malwares" name:"Malwares" list`

		// 异地登录事件统计数据数组。
		NonLocalLoginPlaces []*SecurityTrend `json:"NonLocalLoginPlaces" name:"NonLocalLoginPlaces" list`

		// 密码破解事件统计数据数组。
		BruteAttacks []*SecurityTrend `json:"BruteAttacks" name:"BruteAttacks" list`

		// 漏洞统计数据数组。
		Vuls []*SecurityTrend `json:"Vuls" name:"Vuls" list`

		// 基线统计数据数组。
		BaseLines []*SecurityTrend `json:"BaseLines" name:"BaseLines" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeSecurityTrendsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSecurityTrendsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeUsualLoginPlacesRequest struct {
	*tchttp.BaseRequest

	// 云镜客户端UUID
	Uuid *string `json:"Uuid" name:"Uuid"`
}

func (r *DescribeUsualLoginPlacesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeUsualLoginPlacesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeUsualLoginPlacesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 常用登录地数组
		UsualLoginPlaces []*UsualPlace `json:"UsualLoginPlaces" name:"UsualLoginPlaces" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeUsualLoginPlacesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeUsualLoginPlacesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeVulInfoRequest struct {
	*tchttp.BaseRequest

	// 漏洞种类ID。
	VulId *uint64 `json:"VulId" name:"VulId"`
}

func (r *DescribeVulInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeVulInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeVulInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 漏洞种类ID。
		VulId *uint64 `json:"VulId" name:"VulId"`

		// 漏洞名称。
		VulName *string `json:"VulName" name:"VulName"`

		// 漏洞等级。
		VulLevel *string `json:"VulLevel" name:"VulLevel"`

		// 漏洞类型。
		VulType *string `json:"VulType" name:"VulType"`

		// 漏洞描述。
		Description *string `json:"Description" name:"Description"`

		// 修复方案。
		RepairPlan *string `json:"RepairPlan" name:"RepairPlan"`

		// 漏洞CVE。
		CveId *string `json:"CveId" name:"CveId"`

		// 参考链接。
		Reference *string `json:"Reference" name:"Reference"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeVulInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeVulInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeVulScanResultRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeVulScanResultRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeVulScanResultRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeVulScanResultResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 漏洞数量。
		VulNum *uint64 `json:"VulNum" name:"VulNum"`

		// 专业版机器数。
		ProVersionNum *uint64 `json:"ProVersionNum" name:"ProVersionNum"`

		// 受影响的专业版主机数。
		ImpactedHostNum *uint64 `json:"ImpactedHostNum" name:"ImpactedHostNum"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeVulScanResultResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeVulScanResultResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeVulsRequest struct {
	*tchttp.BaseRequest

	// 漏洞类型。
	// <li>WEB：Web应用漏洞</li>
	// <li>SYSTEM：系统组件漏洞</li>
	// <li>BASELINE：安全基线</li>
	VulType *string `json:"VulType" name:"VulType"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 过滤条件。
	// <li>Status - String - 是否必填：否 - 状态筛选（UN_OPERATED: 待处理 | FIXED：已修复）
	// 
	// Status过滤条件值只能取其一，不能是“或”逻辑。
	Filters []*Filter `json:"Filters" name:"Filters" list`
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

		// 漏洞列表数组。
		Vuls []*Vul `json:"Vuls" name:"Vuls" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type DescribeWeeklyReportBruteAttacksRequest struct {
	*tchttp.BaseRequest

	// 专业周报开始时间。
	BeginDate *string `json:"BeginDate" name:"BeginDate"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`
}

func (r *DescribeWeeklyReportBruteAttacksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeWeeklyReportBruteAttacksRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeWeeklyReportBruteAttacksResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 专业周报密码破解数组。
		WeeklyReportBruteAttacks []*WeeklyReportBruteAttack `json:"WeeklyReportBruteAttacks" name:"WeeklyReportBruteAttacks" list`

		// 记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeWeeklyReportBruteAttacksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeWeeklyReportBruteAttacksResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeWeeklyReportInfoRequest struct {
	*tchttp.BaseRequest

	// 专业周报开始时间。
	BeginDate *string `json:"BeginDate" name:"BeginDate"`
}

func (r *DescribeWeeklyReportInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeWeeklyReportInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeWeeklyReportInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 账号所属公司或个人名称。
		CompanyName *string `json:"CompanyName" name:"CompanyName"`

		// 机器总数。
		MachineNum *uint64 `json:"MachineNum" name:"MachineNum"`

		// 云镜客户端在线数。
		OnlineMachineNum *uint64 `json:"OnlineMachineNum" name:"OnlineMachineNum"`

		// 云镜客户端离线数。
		OfflineMachineNum *uint64 `json:"OfflineMachineNum" name:"OfflineMachineNum"`

		// 开通云镜专业版数量。
		ProVersionMachineNum *uint64 `json:"ProVersionMachineNum" name:"ProVersionMachineNum"`

		// 周报开始时间。
		BeginDate *string `json:"BeginDate" name:"BeginDate"`

		// 周报结束时间。
		EndDate *string `json:"EndDate" name:"EndDate"`

		// 安全等级。
	// <li>HIGH：高</li>
	// <li>MIDDLE：中</li>
	// <li>LOW：低</li>
		Level *string `json:"Level" name:"Level"`

		// 木马记录数。
		MalwareNum *uint64 `json:"MalwareNum" name:"MalwareNum"`

		// 异地登录数。
		NonlocalLoginNum *uint64 `json:"NonlocalLoginNum" name:"NonlocalLoginNum"`

		// 密码破解成功数。
		BruteAttackSuccessNum *uint64 `json:"BruteAttackSuccessNum" name:"BruteAttackSuccessNum"`

		// 漏洞数。
		VulNum *uint64 `json:"VulNum" name:"VulNum"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeWeeklyReportInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeWeeklyReportInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeWeeklyReportMalwaresRequest struct {
	*tchttp.BaseRequest

	// 专业周报开始时间。
	BeginDate *string `json:"BeginDate" name:"BeginDate"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`
}

func (r *DescribeWeeklyReportMalwaresRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeWeeklyReportMalwaresRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeWeeklyReportMalwaresResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 专业周报木马数据。
		WeeklyReportMalwares []*WeeklyReportMalware `json:"WeeklyReportMalwares" name:"WeeklyReportMalwares" list`

		// 记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeWeeklyReportMalwaresResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeWeeklyReportMalwaresResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeWeeklyReportNonlocalLoginPlacesRequest struct {
	*tchttp.BaseRequest

	// 专业周报开始时间。
	BeginDate *string `json:"BeginDate" name:"BeginDate"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`
}

func (r *DescribeWeeklyReportNonlocalLoginPlacesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeWeeklyReportNonlocalLoginPlacesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeWeeklyReportNonlocalLoginPlacesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 专业周报异地登录数据。
		WeeklyReportNonlocalLoginPlaces []*WeeklyReportNonlocalLoginPlace `json:"WeeklyReportNonlocalLoginPlaces" name:"WeeklyReportNonlocalLoginPlaces" list`

		// 记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeWeeklyReportNonlocalLoginPlacesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeWeeklyReportNonlocalLoginPlacesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeWeeklyReportVulsRequest struct {
	*tchttp.BaseRequest

	// 专业版周报开始时间。
	BeginDate *string `json:"BeginDate" name:"BeginDate"`

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`
}

func (r *DescribeWeeklyReportVulsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeWeeklyReportVulsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeWeeklyReportVulsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 专业周报漏洞数据数组。
		WeeklyReportVuls []*WeeklyReportVul `json:"WeeklyReportVuls" name:"WeeklyReportVuls" list`

		// 记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeWeeklyReportVulsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeWeeklyReportVulsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeWeeklyReportsRequest struct {
	*tchttp.BaseRequest

	// 返回数量，默认为10，最大值为100。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。
	Offset *uint64 `json:"Offset" name:"Offset"`
}

func (r *DescribeWeeklyReportsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeWeeklyReportsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeWeeklyReportsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 专业周报列表数组。
		WeeklyReports []*WeeklyReport `json:"WeeklyReports" name:"WeeklyReports" list`

		// 记录总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeWeeklyReportsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeWeeklyReportsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ExportMaliciousRequestsRequest struct {
	*tchttp.BaseRequest
}

func (r *ExportMaliciousRequestsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ExportMaliciousRequestsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ExportMaliciousRequestsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 导出文件下载链接地址。
		DownloadUrl *string `json:"DownloadUrl" name:"DownloadUrl"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ExportMaliciousRequestsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ExportMaliciousRequestsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Filter struct {

	// 过滤键的名称。
	Name *string `json:"Name" name:"Name"`

	// 一个或者多个过滤值。
	Values []*string `json:"Values" name:"Values" list`
}

type HistoryAccount struct {

	// 唯一ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 云镜客户端唯一Uuid。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 主机内网IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 主机名。
	MachineName *string `json:"MachineName" name:"MachineName"`

	// 帐号名。
	Username *string `json:"Username" name:"Username"`

	// 帐号变更类型。
	// <li>CREATE：表示新增帐号</li>
	// <li>MODIFY：表示修改帐号</li>
	// <li>DELETE：表示删除帐号</li>
	ModifyType *string `json:"ModifyType" name:"ModifyType"`

	// 变更时间。
	ModifyTime *string `json:"ModifyTime" name:"ModifyTime"`
}

type IgnoreImpactedHostsRequest struct {
	*tchttp.BaseRequest

	// 漏洞ID数组。
	Ids []*uint64 `json:"Ids" name:"Ids" list`
}

func (r *IgnoreImpactedHostsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *IgnoreImpactedHostsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type IgnoreImpactedHostsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *IgnoreImpactedHostsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *IgnoreImpactedHostsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ImpactedHost struct {

	// 漏洞ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 主机IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 主机名称。
	MachineName *string `json:"MachineName" name:"MachineName"`

	// 最后检测时间。
	LastScanTime *string `json:"LastScanTime" name:"LastScanTime"`

	// 漏洞状态。
	// <li>UN_OPERATED ：待处理</li>
	// <li>SCANING : 扫描中</li>
	// <li>FIXED : 已修复</li>
	VulStatus *string `json:"VulStatus" name:"VulStatus"`

	// 云镜客户端唯一标识UUID。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 漏洞描述。
	Description *string `json:"Description" name:"Description"`

	// 漏洞种类ID。
	VulId *uint64 `json:"VulId" name:"VulId"`
}

type InquiryPriceOpenProVersionPrepaidRequest struct {
	*tchttp.BaseRequest

	// 预付费模式(包年包月)参数设置。
	ChargePrepaid *ChargePrepaid `json:"ChargePrepaid" name:"ChargePrepaid"`

	// 需要开通专业版机器列表数组。
	Machines []*ProVersionMachine `json:"Machines" name:"Machines" list`
}

func (r *InquiryPriceOpenProVersionPrepaidRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceOpenProVersionPrepaidRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InquiryPriceOpenProVersionPrepaidResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 预支费用的原价，单位：元。
		DiscountPrice *float64 `json:"DiscountPrice" name:"DiscountPrice"`

		// 预支费用的折扣价，单位：元。
		DiscountPrice *float64 `json:"DiscountPrice" name:"DiscountPrice"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *InquiryPriceOpenProVersionPrepaidResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceOpenProVersionPrepaidResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Machine struct {

	// 主机名称。
	MachineName *string `json:"MachineName" name:"MachineName"`

	// 主机系统。
	MachineOs *string `json:"MachineOs" name:"MachineOs"`

	// 主机状态。
	// <li>OFFLINE: 离线  </li>
	// <li>ONLINE: 在线</li>
	MachineStatus *string `json:"MachineStatus" name:"MachineStatus"`

	// 云镜客户端唯一Uuid，若客户端长时间不在线将返回空字符。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// CVM或BM机器唯一Uuid。
	Quuid *string `json:"Quuid" name:"Quuid"`

	// 漏洞数，非专业版将返回：0。
	VulNum *uint64 `json:"VulNum" name:"VulNum"`

	// 主机IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 是否是专业版。
	// <li>true： 是</li>
	// <li>false：否</li>
	IsProVersion *bool `json:"IsProVersion" name:"IsProVersion"`

	// 主机外网IP。
	MachineWanIp *string `json:"MachineWanIp" name:"MachineWanIp"`

	// 主机状态。
	// <li>POSTPAY: 表示后付费，即按量计费  </li>
	// <li>PREPAY: 表示预付费，即包年包月</li>
	PayMode *string `json:"PayMode" name:"PayMode"`
}

type MaliciousRequest struct {
	*tchttp.BaseRequest

	// 记录ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 云镜客户端UUID。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 主机内网IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 主机名。
	MachineName *string `json:"MachineName" name:"MachineName"`

	// 恶意请求域名。
	Domain *string `json:"Domain" name:"Domain"`

	// 恶意请求数。
	Count *uint64 `json:"Count" name:"Count"`

	// 进程名。
	ProcessName *string `json:"ProcessName" name:"ProcessName"`

	// 记录状态。
	// <li>UN_OPERATED：待处理</li>
	// <li>TRUSTED：已信任</li>
	// <li>UN_TRUSTED：已取消信任</li>
	Status *string `json:"Status" name:"Status"`

	// 恶意请求域名描述。
	Description *string `json:"Description" name:"Description"`

	// 参考地址。
	Reference *string `json:"Reference" name:"Reference"`

	// 发现时间。
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 记录合并时间。
	MergeTime *string `json:"MergeTime" name:"MergeTime"`

	// 进程MD5
	// 值。
	ProcessMd5 *string `json:"ProcessMd5" name:"ProcessMd5"`

	// 执行命令行。
	CmdLine *string `json:"CmdLine" name:"CmdLine"`

	// 进程PID。
	Pid *uint64 `json:"Pid" name:"Pid"`
}

func (r *MaliciousRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *MaliciousRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Malware struct {

	// 事件ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 主机IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 当前木马状态。
	// <li>UN_OPERATED：未处理</li><li>SEGREGATED：已隔离</li><li>TRUSTED：已信任</li>
	// <li>SEPARATING：隔离中</li><li>RECOVERING：恢复中</li>
	Status *string `json:"Status" name:"Status"`

	// 木马所在的路径。
	FilePath *string `json:"FilePath" name:"FilePath"`

	// 木马描述。
	Description *string `json:"Description" name:"Description"`

	// 主机名称。
	MachineName *string `json:"MachineName" name:"MachineName"`

	// 木马文件创建时间。
	FileCreateTime *string `json:"FileCreateTime" name:"FileCreateTime"`

	// 木马文件修改时间。
	ModifyTime *string `json:"ModifyTime" name:"ModifyTime"`

	// 云镜客户端唯一标识UUID。
	Uuid *string `json:"Uuid" name:"Uuid"`
}

type MisAlarmNonlocalLoginPlacesRequest struct {
	*tchttp.BaseRequest

	// 异地登录事件Id数组。
	Ids []*uint64 `json:"Ids" name:"Ids" list`
}

func (r *MisAlarmNonlocalLoginPlacesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *MisAlarmNonlocalLoginPlacesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type MisAlarmNonlocalLoginPlacesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *MisAlarmNonlocalLoginPlacesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *MisAlarmNonlocalLoginPlacesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAlarmAttributeRequest struct {
	*tchttp.BaseRequest

	// 告警项目。
	// <li>Offline：防护软件离线</li>
	// <li>Malware：发现木马文件</li>
	// <li>NonlocalLogin：发现异地登录行为</li>
	// <li>CrackSuccess：被暴力破解成功</li>
	Attribute *string `json:"Attribute" name:"Attribute"`

	// 告警项目属性。
	// <li>CLOSE：关闭</li>
	// <li>OPEN：打开</li>
	Value *string `json:"Value" name:"Value"`
}

func (r *ModifyAlarmAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAlarmAttributeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAlarmAttributeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyAlarmAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAlarmAttributeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAutoOpenProVersionConfigRequest struct {
	*tchttp.BaseRequest

	// 设置自动开通状态。
	// <li>CLOSE：关闭</li>
	// <li>OPEN：打开</li>
	Status *string `json:"Status" name:"Status"`
}

func (r *ModifyAutoOpenProVersionConfigRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAutoOpenProVersionConfigRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAutoOpenProVersionConfigResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyAutoOpenProVersionConfigResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAutoOpenProVersionConfigResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyProVersionRenewFlagRequest struct {
	*tchttp.BaseRequest

	// 自动续费标识。取值范围：
	// <li>NOTIFY_AND_AUTO_RENEW：通知过期且自动续费</li>
	// <li>NOTIFY_AND_MANUAL_RENEW：通知过期不自动续费</li>
	// <li>DISABLE_NOTIFY_AND_MANUAL_RENEW：不通知过期不自动续费</li>
	RenewFlag *string `json:"RenewFlag" name:"RenewFlag"`

	// 主机唯一ID，对应CVM的uuid、BM的instanceId。
	Quuid *string `json:"Quuid" name:"Quuid"`
}

func (r *ModifyProVersionRenewFlagRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyProVersionRenewFlagRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyProVersionRenewFlagResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyProVersionRenewFlagResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyProVersionRenewFlagResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type NonLocalLoginPlace struct {

	// 事件ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 主机IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 登录状态
	// <li>NON_LOCAL_LOGIN：异地登录</li>
	// <li>NORMAL_LOGIN：正常登录</li>
	Status *string `json:"Status" name:"Status"`

	// 用户名。
	UserName *string `json:"UserName" name:"UserName"`

	// 城市ID。
	City *uint64 `json:"City" name:"City"`

	// 国家ID。
	Country *uint64 `json:"Country" name:"Country"`

	// 省份ID。
	Province *uint64 `json:"Province" name:"Province"`

	// 登录IP。
	SrcIp *string `json:"SrcIp" name:"SrcIp"`

	// 机器名称。
	MachineName *string `json:"MachineName" name:"MachineName"`

	// 登录时间。
	LoginTime *string `json:"LoginTime" name:"LoginTime"`

	// 云镜客户端唯一标识Uuid。
	Uuid *string `json:"Uuid" name:"Uuid"`
}

type OpenPort struct {

	// 唯一ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 云镜客户端唯一UUID。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 开放端口号。
	Port *uint64 `json:"Port" name:"Port"`

	// 主机IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 主机名。
	MachineName *string `json:"MachineName" name:"MachineName"`

	// 端口对应进程名。
	ProcessName *string `json:"ProcessName" name:"ProcessName"`

	// 端口对应进程Pid。
	Pid *uint64 `json:"Pid" name:"Pid"`

	// 记录创建时间。
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 记录更新时间。
	ModifyTime *string `json:"ModifyTime" name:"ModifyTime"`
}

type OpenPortStatistics struct {

	// 端口号
	Port *uint64 `json:"Port" name:"Port"`

	// 主机数量
	MachineNum *uint64 `json:"MachineNum" name:"MachineNum"`
}

type OpenProVersionPrepaidRequest struct {
	*tchttp.BaseRequest

	// 购买相关参数。
	ChargePrepaid *ChargePrepaid `json:"ChargePrepaid" name:"ChargePrepaid"`

	// 需要开通专业版主机信息数组。
	Machines []*ProVersionMachine `json:"Machines" name:"Machines" list`
}

func (r *OpenProVersionPrepaidRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *OpenProVersionPrepaidRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type OpenProVersionPrepaidResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 订单ID列表。
		DealIds []*string `json:"DealIds" name:"DealIds" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *OpenProVersionPrepaidResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *OpenProVersionPrepaidResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Place struct {

	// 城市 ID。
	CityId *uint64 `json:"CityId" name:"CityId"`

	// 省份 ID。
	ProvinceId *uint64 `json:"ProvinceId" name:"ProvinceId"`

	// 国家ID，暂只支持国内：1。
	CountryId *uint64 `json:"CountryId" name:"CountryId"`
}

type ProVersionMachine struct {

	// 主机类型。
	// <li>CVM: 虚拟主机</li>
	// <li>BM: 黑石物理机</li>
	MachineType *string `json:"MachineType" name:"MachineType"`

	// 主机所在地域。
	// 如：ap-guangzhou、ap-beijing
	MachineRegion *string `json:"MachineRegion" name:"MachineRegion"`

	// 主机唯一标识Uuid。
	// 黑石的InstanceId，CVM的Uuid
	Quuid *string `json:"Quuid" name:"Quuid"`
}

type Process struct {

	// 唯一ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 云镜客户端唯一UUID。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 主机内网IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 主机名。
	MachineName *string `json:"MachineName" name:"MachineName"`

	// 进程Pid。
	Pid *uint64 `json:"Pid" name:"Pid"`

	// 进程Ppid。
	Ppid *uint64 `json:"Ppid" name:"Ppid"`

	// 进程名。
	ProcessName *string `json:"ProcessName" name:"ProcessName"`

	// 进程用户名。
	Username *string `json:"Username" name:"Username"`

	// 所属平台。
	// <li>WIN32：windows32位</li>
	// <li>WIN64：windows64位</li>
	// <li>LINUX32：Linux32位</li>
	// <li>LINUX64：Linux64位</li>
	Platform *string `json:"Platform" name:"Platform"`

	// 进程路径。
	FullPath *string `json:"FullPath" name:"FullPath"`

	// 创建时间。
	CreateTime *string `json:"CreateTime" name:"CreateTime"`
}

type ProcessStatistics struct {

	// 进程名。
	ProcessName *string `json:"ProcessName" name:"ProcessName"`

	// 主机数量。
	MachineNum *uint64 `json:"MachineNum" name:"MachineNum"`
}

type RecoverMalwaresRequest struct {
	*tchttp.BaseRequest

	// 木马Id数组,单次最大删除不能超过200条
	Ids []*uint64 `json:"Ids" name:"Ids" list`
}

func (r *RecoverMalwaresRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RecoverMalwaresRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RecoverMalwaresResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 恢复成功id数组
		SuccessIds []*uint64 `json:"SuccessIds" name:"SuccessIds" list`

		// 恢复失败id数组
		FailedIds []*uint64 `json:"FailedIds" name:"FailedIds" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RecoverMalwaresResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RecoverMalwaresResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RenewProVersionRequest struct {
	*tchttp.BaseRequest

	// 购买相关参数。
	ChargePrepaid *ChargePrepaid `json:"ChargePrepaid" name:"ChargePrepaid"`

	// 主机唯一ID，对应CVM的uuid、BM的InstanceId。
	Quuid *string `json:"Quuid" name:"Quuid"`
}

func (r *RenewProVersionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RenewProVersionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RenewProVersionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RenewProVersionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RenewProVersionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RescanImpactedHostRequest struct {
	*tchttp.BaseRequest

	// 漏洞ID。
	Id *uint64 `json:"Id" name:"Id"`
}

func (r *RescanImpactedHostRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RescanImpactedHostRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RescanImpactedHostResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RescanImpactedHostResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RescanImpactedHostResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SecurityDynamic struct {

	// 云镜客户端UUID。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 安全事件发生事件。
	EventTime *string `json:"EventTime" name:"EventTime"`

	// 安全事件类型。
	// <li>MALWARE：木马事件</li>
	// <li>NON_LOCAL_LOGIN：异地登录</li>
	// <li>BRUTEATTACK_SUCCESS：密码破解成功</li>
	// <li>VUL：漏洞</li>
	// <li>BASELINE：安全基线</li>
	EventType *string `json:"EventType" name:"EventType"`

	// 安全事件消息。
	Message *string `json:"Message" name:"Message"`
}

type SecurityTrend struct {

	// 事件时间。
	Date *string `json:"Date" name:"Date"`

	// 事件数量。
	EventNum *uint64 `json:"EventNum" name:"EventNum"`
}

type SeparateMalwaresRequest struct {
	*tchttp.BaseRequest

	// 木马事件Id数组。
	Ids []*uint64 `json:"Ids" name:"Ids" list`
}

func (r *SeparateMalwaresRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SeparateMalwaresRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SeparateMalwaresResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 隔离成功的id数组。
		SuccessIds []*uint64 `json:"SuccessIds" name:"SuccessIds" list`

		// 隔离失败的id数组。
		FailedIds []*uint64 `json:"FailedIds" name:"FailedIds" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *SeparateMalwaresResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SeparateMalwaresResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type TrustMaliciousRequestRequest struct {
	*tchttp.BaseRequest

	// 恶意请求记录ID。
	Id *uint64 `json:"Id" name:"Id"`
}

func (r *TrustMaliciousRequestRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *TrustMaliciousRequestRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type TrustMaliciousRequestResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *TrustMaliciousRequestResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *TrustMaliciousRequestResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type TrustMalwaresRequest struct {
	*tchttp.BaseRequest

	// 木马ID数组。
	Ids []*uint64 `json:"Ids" name:"Ids" list`
}

func (r *TrustMalwaresRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *TrustMalwaresRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type TrustMalwaresResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *TrustMalwaresResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *TrustMalwaresResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UntrustMaliciousRequestRequest struct {
	*tchttp.BaseRequest

	// 受信任记录ID。
	Id *uint64 `json:"Id" name:"Id"`
}

func (r *UntrustMaliciousRequestRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UntrustMaliciousRequestRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UntrustMaliciousRequestResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UntrustMaliciousRequestResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UntrustMaliciousRequestResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UntrustMalwaresRequest struct {
	*tchttp.BaseRequest

	// 木马Id数组，单次最大处理不能超过200条。
	Ids []*uint64 `json:"Ids" name:"Ids" list`
}

func (r *UntrustMalwaresRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UntrustMalwaresRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UntrustMalwaresResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UntrustMalwaresResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UntrustMalwaresResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UsualPlace struct {

	// ID。
	Id *uint64 `json:"Id" name:"Id"`

	// 云镜客户端唯一标识UUID。
	Uuid *string `json:"Uuid" name:"Uuid"`

	// 国家 ID。
	CountryId *uint64 `json:"CountryId" name:"CountryId"`

	// 省份 ID。
	ProvinceId *uint64 `json:"ProvinceId" name:"ProvinceId"`

	// 城市 ID。
	CityId *uint64 `json:"CityId" name:"CityId"`
}

type Vul struct {

	// 漏洞种类ID
	VulId *uint64 `json:"VulId" name:"VulId"`

	// 漏洞名称
	VulName *string `json:"VulName" name:"VulName"`

	// 漏洞危害等级:
	// HIGH：高危
	// MIDDLE：中危
	// LOW：低危
	// NOTICE：提示
	VulLevel *string `json:"VulLevel" name:"VulLevel"`

	// 最后扫描时间
	LastScanTime *string `json:"LastScanTime" name:"LastScanTime"`

	// 受影响机器数量
	ImpactedHostNum *uint64 `json:"ImpactedHostNum" name:"ImpactedHostNum"`

	// 漏洞状态
	// * UN_OPERATED : 待处理
	// * FIXED : 已修复
	VulStatus *string `json:"VulStatus" name:"VulStatus"`
}

type WeeklyReport struct {

	// 周报开始时间。
	BeginDate *string `json:"BeginDate" name:"BeginDate"`

	// 周报结束时间。
	EndDate *string `json:"EndDate" name:"EndDate"`
}

type WeeklyReportBruteAttack struct {

	// 主机IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 被破解用户名。
	Username *string `json:"Username" name:"Username"`

	// 源IP。
	SrcIp *string `json:"SrcIp" name:"SrcIp"`

	// 尝试次数。
	Count *uint64 `json:"Count" name:"Count"`

	// 攻击时间。
	AttackTime *string `json:"AttackTime" name:"AttackTime"`
}

type WeeklyReportMalware struct {

	// 主机IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 木马文件路径。
	FilePath *string `json:"FilePath" name:"FilePath"`

	// 木马文件MD5值。
	Md5 *string `json:"Md5" name:"Md5"`

	// 木马发现时间。
	FindTime *string `json:"FindTime" name:"FindTime"`

	// 当前木马状态。
	// <li>UN_OPERATED：未处理</li>
	// <li>SEGREGATED：已隔离</li>
	// <li>TRUSTED：已信任</li>
	// <li>SEPARATING：隔离中</li>
	// <li>RECOVERING：恢复中</li>
	Status *string `json:"Status" name:"Status"`
}

type WeeklyReportNonlocalLoginPlace struct {

	// 主机IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 用户名。
	Username *string `json:"Username" name:"Username"`

	// 源IP。
	SrcIp *string `json:"SrcIp" name:"SrcIp"`

	// 国家ID。
	Country *uint64 `json:"Country" name:"Country"`

	// 省份ID。
	Province *uint64 `json:"Province" name:"Province"`

	// 城市ID。
	City *uint64 `json:"City" name:"City"`

	// 登录时间。
	LoginTime *string `json:"LoginTime" name:"LoginTime"`
}

type WeeklyReportVul struct {

	// 主机内网IP。
	MachineIp *string `json:"MachineIp" name:"MachineIp"`

	// 漏洞名称。
	VulName *string `json:"VulName" name:"VulName"`

	// 漏洞类型。
	// <li> WEB : WEB漏洞</li>
	// <li> SYSTEM :系统组件漏洞</li>
	// <li> BASELINE : 安全基线</li>
	VulType *string `json:"VulType" name:"VulType"`

	// 漏洞描述。
	Description *string `json:"Description" name:"Description"`

	// 漏洞状态。
	// <li> UN_OPERATED : 待处理</li>
	// <li> SCANING : 扫描中</li>
	// <li> FIXED : 已修复</li>
	VulStatus *string `json:"VulStatus" name:"VulStatus"`

	// 最后扫描时间。
	LastScanTime *string `json:"LastScanTime" name:"LastScanTime"`
}
