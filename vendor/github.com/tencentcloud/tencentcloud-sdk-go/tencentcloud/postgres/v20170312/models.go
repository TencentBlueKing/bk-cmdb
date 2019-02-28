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

package v20170312

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type AccountInfo struct {

	// 实例ID，形如postgres-lnp6j617
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`

	// 帐号
	UserName *string `json:"UserName" name:"UserName"`

	// 帐号备注
	Remark *string `json:"Remark" name:"Remark"`

	// 帐号状态。 1-创建中，2-正常，3-修改中，4-密码重置中，-1-删除中
	Status *int64 `json:"Status" name:"Status"`

	// 帐号创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 帐号最后一次更新时间
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`
}

type CloseDBExtranetAccessRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如postgres-6r233v55
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`
}

func (r *CloseDBExtranetAccessRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CloseDBExtranetAccessRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CloseDBExtranetAccessResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务流程ID
		FlowId *int64 `json:"FlowId" name:"FlowId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CloseDBExtranetAccessResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CloseDBExtranetAccessResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDBInstancesRequest struct {
	*tchttp.BaseRequest

	// 售卖规格ID。该参数可以通过调用DescribeProductConfig的返回值中的SpecCode字段来获取。
	SpecCode *string `json:"SpecCode" name:"SpecCode"`

	// PostgreSQL内核版本，目前只支持：9.3.5、9.5.4两种版本。
	DBVersion *string `json:"DBVersion" name:"DBVersion"`

	// 实例容量大小，单位：GB。
	Storage *uint64 `json:"Storage" name:"Storage"`

	// 一次性购买的实例数量。取值1-100
	InstanceCount *uint64 `json:"InstanceCount" name:"InstanceCount"`

	// 购买时长，单位：月。目前只支持1,2,3,4,5,6,7,8,9,10,11,12,24,36这些值。
	Period *uint64 `json:"Period" name:"Period"`

	// 可用区ID。该参数可以通过调用 DescribeZones 接口的返回值中的Zone字段来获取。
	Zone *string `json:"Zone" name:"Zone"`

	// 项目ID。
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 实例计费类型。目前只支持：PREPAID（预付费，即包年包月）。
	InstanceChargeType *string `json:"InstanceChargeType" name:"InstanceChargeType"`

	// 是否自动使用代金券。1（是），0（否），默认不使用。
	AutoVoucher *uint64 `json:"AutoVoucher" name:"AutoVoucher"`

	// 代金券ID列表，目前仅支持指定一张代金券。
	VoucherIds []*string `json:"VoucherIds" name:"VoucherIds" list`

	// 私有网络ID。
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 私有网络子网ID。
	SubnetId *string `json:"SubnetId" name:"SubnetId"`
}

func (r *CreateDBInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDBInstancesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDBInstancesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 订单号列表。每个实例对应一个订单号。
		DealNames []*string `json:"DealNames" name:"DealNames" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateDBInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDBInstancesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DBBackup struct {

	// 备份文件唯一标识
	Id *int64 `json:"Id" name:"Id"`

	// 文件生成的开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 文件生成的结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 文件大小(K)
	Size *int64 `json:"Size" name:"Size"`

	// 策略（0-实例备份；1-多库备份）
	Strategy *int64 `json:"Strategy" name:"Strategy"`

	// 类型（0-定时；1-临时）
	Way *int64 `json:"Way" name:"Way"`

	// 备份方式（1-完整；2-日志；3-差异）
	Type *int64 `json:"Type" name:"Type"`

	// 状态（0-创建中；1-成功；2-失败）
	Status *int64 `json:"Status" name:"Status"`

	// DB列表
	DbList []*string `json:"DbList" name:"DbList" list`

	// 内网下载地址
	InternalAddr *string `json:"InternalAddr" name:"InternalAddr"`

	// 外网下载地址
	ExternalAddr *string `json:"ExternalAddr" name:"ExternalAddr"`
}

type DBInstance struct {

	// 实例所属地域，如: ap-guangzhou，对应RegionSet的Region字段
	Region *string `json:"Region" name:"Region"`

	// 实例所属可用区， 如：ap-guangzhou-3，对应ZoneSet的Zone字段
	Zone *string `json:"Zone" name:"Zone"`

	// 项目ID
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`

	// 私有网络ID
	VpcId *string `json:"VpcId" name:"VpcId"`

	// SubnetId
	SubnetId *string `json:"SubnetId" name:"SubnetId"`

	// 实例ID
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`

	// 实例名称
	DBInstanceName *string `json:"DBInstanceName" name:"DBInstanceName"`

	// 实例状态
	DBInstanceStatus *string `json:"DBInstanceStatus" name:"DBInstanceStatus"`

	// 实例分配的内存大小，单位：GB
	DBInstanceMemory *uint64 `json:"DBInstanceMemory" name:"DBInstanceMemory"`

	// 实例分配的存储空间大小，单位：GB
	DBInstanceStorage *uint64 `json:"DBInstanceStorage" name:"DBInstanceStorage"`

	// 实例分配的CPU数量，单位：个
	DBInstanceCpu *uint64 `json:"DBInstanceCpu" name:"DBInstanceCpu"`

	// 售卖规格ID
	DBInstanceClass *string `json:"DBInstanceClass" name:"DBInstanceClass"`

	// 实例类型，类型有：1、primary（主实例）；2、readonly（只读实例）；3、guard（灾备实例）；4、temp（临时实例）
	DBInstanceType *string `json:"DBInstanceType" name:"DBInstanceType"`

	// 实例版本，目前只支持standard（双机高可用版, 一主一从）
	DBInstanceVersion *string `json:"DBInstanceVersion" name:"DBInstanceVersion"`

	// 实例DB字符集
	DBCharset *string `json:"DBCharset" name:"DBCharset"`

	// PostgreSQL内核版本
	DBVersion *string `json:"DBVersion" name:"DBVersion"`

	// 实例创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 实例执行最后一次更新的时间
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`

	// 实例到期时间
	ExpireTime *string `json:"ExpireTime" name:"ExpireTime"`

	// 实例隔离时间
	IsolatedTime *string `json:"IsolatedTime" name:"IsolatedTime"`

	// 计费模式，1、prepaid（包年包月,预付费）；2、postpaid（按量计费，后付费）
	PayType *string `json:"PayType" name:"PayType"`

	// 是否自动续费，1：自动续费，0：不自动续费
	AutoRenew *uint64 `json:"AutoRenew" name:"AutoRenew"`

	// 实例网络连接信息
	DBInstanceNetInfo []*DBInstanceNetInfo `json:"DBInstanceNetInfo" name:"DBInstanceNetInfo" list`
}

type DBInstanceNetInfo struct {

	// DNS域名
	Address *string `json:"Address" name:"Address"`

	// Ip
	Ip *string `json:"Ip" name:"Ip"`

	// 连接Port地址
	Port *uint64 `json:"Port" name:"Port"`

	// 网络类型，1、inner（内网地址）；2、public（外网地址）
	NetType *string `json:"NetType" name:"NetType"`

	// 网络连接状态
	Status *string `json:"Status" name:"Status"`
}

type DescribeAccountsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如postgres-6fego161
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`

	// 分页返回，每页最大返回数目，默认20，取值范围为1-100
	Limit *int64 `json:"Limit" name:"Limit"`

	// 分页返回，返回第几页的用户数据。页码从0开始计数
	Offset *int64 `json:"Offset" name:"Offset"`

	// 返回数据按照创建时间或者用户名排序。取值只能为createTime或者name。createTime-按照创建时间排序；name-按照用户名排序
	OrderBy *string `json:"OrderBy" name:"OrderBy"`

	// 返回结果是升序还是降序。取值只能为desc或者asc。desc-降序；asc-升序
	OrderByType *string `json:"OrderByType" name:"OrderByType"`
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

		// 本次调用接口共返回了多少条数据。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 帐号列表详细信息。
		Details []*AccountInfo `json:"Details" name:"Details" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
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

type DescribeDBBackupsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如postgres-4wdeb0zv。
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`

	// 备份方式（1-全量）。目前只支持全量，取值为1。
	Type *int64 `json:"Type" name:"Type"`

	// 查询开始时间，形如2018-06-10 17:06:38，起始时间不得小于7天以前
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 查询结束时间，形如2018-06-10 17:06:38
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 备份列表分页返回，每页返回数量，默认为 20，最小为1，最大值为 100。
	Limit *int64 `json:"Limit" name:"Limit"`

	// 返回结果中的第几页，从第0页开始。默认为0。
	Offset *int64 `json:"Offset" name:"Offset"`
}

func (r *DescribeDBBackupsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBBackupsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBBackupsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 返回备份列表中备份文件的个数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 备份列表
		BackupList []*DBBackup `json:"BackupList" name:"BackupList" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBBackupsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBBackupsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBErrlogsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如postgres-5bq3wfjd
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`

	// 查询起始时间，形如2018-01-01 00:00:00，起始时间不得小于7天以前
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 查询结束时间，形如2018-01-01 00:00:00
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 数据库名字
	DatabaseName *string `json:"DatabaseName" name:"DatabaseName"`

	// 搜索关键字
	SearchKeys []*string `json:"SearchKeys" name:"SearchKeys" list`

	// 分页返回，每页返回的最大数量。取值为1-100
	Limit *int64 `json:"Limit" name:"Limit"`

	// 分页返回，返回第几页的数据，从第0页开始计数
	Offset *int64 `json:"Offset" name:"Offset"`
}

func (r *DescribeDBErrlogsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBErrlogsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBErrlogsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 本次调用返回了多少条数据
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 错误日志列表
		Details []*ErrLogDetail `json:"Details" name:"Details" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBErrlogsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBErrlogsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstanceAttributeRequest struct {
	*tchttp.BaseRequest

	// 实例ID
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`
}

func (r *DescribeDBInstanceAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBInstanceAttributeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstanceAttributeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例详细信息。
		DBInstance *DBInstance `json:"DBInstance" name:"DBInstance"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBInstanceAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBInstanceAttributeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstancesRequest struct {
	*tchttp.BaseRequest

	// 过滤条件，目前支持：db-instance-id、db-instance-name两种。
	Filters []*Filter `json:"Filters" name:"Filters" list`

	// 每页显示数量，默认返回10条。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 分页序号，从0开始。
	Offset *uint64 `json:"Offset" name:"Offset"`
}

func (r *DescribeDBInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBInstancesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstancesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 查询到的实例数量。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 实例详细信息集合。
		DBInstanceSet []*DBInstance `json:"DBInstanceSet" name:"DBInstanceSet" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBInstancesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBSlowlogsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如postgres-lnp6j617
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`

	// 查询起始时间，形如2018-06-10 17:06:38，起始时间不得小于7天以前
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 查询结束时间，形如2018-06-10 17:06:38
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 数据库名字
	DatabaseName *string `json:"DatabaseName" name:"DatabaseName"`

	// 按照何种指标排序，取值为sum_calls或者sum_cost_time。sum_calls-总调用次数；sum_cost_time-总的花费时间
	OrderBy *string `json:"OrderBy" name:"OrderBy"`

	// 排序规则。desc-降序；asc-升序
	OrderByType *string `json:"OrderByType" name:"OrderByType"`

	// 分页返回结果，每页最大返回数量，取值为1-100，默认20
	Limit *int64 `json:"Limit" name:"Limit"`

	// 分页返回结果，返回结果的第几页，从0开始计数
	Offset *int64 `json:"Offset" name:"Offset"`
}

func (r *DescribeDBSlowlogsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBSlowlogsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBSlowlogsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 本次返回多少条数据
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 慢查询日志详情
		Detail *SlowlogDetail `json:"Detail" name:"Detail"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBSlowlogsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBSlowlogsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBXlogsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如postgres-4wdeb0zv。
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`

	// 查询开始时间，形如2018-06-10 17:06:38，起始时间不得小于7天以前
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 查询结束时间，形如2018-06-10 17:06:38
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 分页返回，表示返回第几页的条目。从第0页开始计数。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 分页返回，表示每页有多少条目。取值为1-100。
	Limit *int64 `json:"Limit" name:"Limit"`
}

func (r *DescribeDBXlogsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBXlogsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBXlogsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 表示此次返回结果有多少条数据。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// Xlog列表
		XlogList []*Xlog `json:"XlogList" name:"XlogList" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBXlogsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBXlogsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeOrdersRequest struct {
	*tchttp.BaseRequest

	// 订单名集合
	DealNames []*string `json:"DealNames" name:"DealNames" list`
}

func (r *DescribeOrdersRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeOrdersRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeOrdersResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 订单数量
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 订单数组
		Deals []*PgDeal `json:"Deals" name:"Deals" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeOrdersResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeOrdersResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProductConfigRequest struct {
	*tchttp.BaseRequest

	// 可用区名称
	Zone *string `json:"Zone" name:"Zone"`
}

func (r *DescribeProductConfigRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProductConfigRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProductConfigResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 售卖规格列表。
		SpecInfoList []*SpecInfo `json:"SpecInfoList" name:"SpecInfoList" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeProductConfigResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProductConfigResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeRegionsRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeRegionsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeRegionsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeRegionsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 返回的结果数量。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 地域信息集合。
		RegionSet []*RegionInfo `json:"RegionSet" name:"RegionSet" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeRegionsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeRegionsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZonesRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeZonesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZonesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZonesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 返回的结果数量。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 可用区信息集合。
		ZoneSet []*ZoneInfo `json:"ZoneSet" name:"ZoneSet" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeZonesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZonesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ErrLogDetail struct {

	// 用户名
	UserName *string `json:"UserName" name:"UserName"`

	// 数据库名字
	Database *string `json:"Database" name:"Database"`

	// 错误发生时间
	ErrTime *string `json:"ErrTime" name:"ErrTime"`

	// 错误消息
	ErrMsg *string `json:"ErrMsg" name:"ErrMsg"`
}

type Filter struct {

	// 过滤键的名称。
	Name *string `json:"Name" name:"Name"`

	// 一个或者多个过滤值。
	Values []*string `json:"Values" name:"Values" list`
}

type InitDBInstancesRequest struct {
	*tchttp.BaseRequest

	// 实例ID集合。
	DBInstanceIdSet []*string `json:"DBInstanceIdSet" name:"DBInstanceIdSet" list`

	// 实例根账号用户名。
	AdminName *string `json:"AdminName" name:"AdminName"`

	// 实例根账号用户名对应的密码。
	AdminPassword *string `json:"AdminPassword" name:"AdminPassword"`

	// 实例字符集，目前只支持：UTF8、LATIN1。
	Charset *string `json:"Charset" name:"Charset"`
}

func (r *InitDBInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InitDBInstancesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InitDBInstancesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例ID集合。
		DBInstanceIdSet []*string `json:"DBInstanceIdSet" name:"DBInstanceIdSet" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *InitDBInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InitDBInstancesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InquiryPriceCreateDBInstancesRequest struct {
	*tchttp.BaseRequest

	// 可用区ID。该参数可以通过调用 DescribeZones 接口的返回值中的Zone字段来获取。
	Zone *string `json:"Zone" name:"Zone"`

	// 规格ID。该参数可以通过调用DescribeProductConfig接口的返回值中的SpecCode字段来获取。
	SpecCode *string `json:"SpecCode" name:"SpecCode"`

	// 存储容量大小，单位：GB。
	Storage *uint64 `json:"Storage" name:"Storage"`

	// 实例数量。目前最大数量不超过100，如需一次性创建更多实例，请联系客服支持。
	InstanceCount *uint64 `json:"InstanceCount" name:"InstanceCount"`

	// 购买时长，单位：月。目前只支持1,2,3,4,5,6,7,8,9,10,11,12,24,36这些值。
	Period *uint64 `json:"Period" name:"Period"`

	// 计费ID。该参数可以通过调用DescribeProductConfig接口的返回值中的Pid字段来获取。
	Pid *uint64 `json:"Pid" name:"Pid"`

	// 实例计费类型。目前只支持：PREPAID（预付费，即包年包月）。
	InstanceChargeType *string `json:"InstanceChargeType" name:"InstanceChargeType"`
}

func (r *InquiryPriceCreateDBInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceCreateDBInstancesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InquiryPriceCreateDBInstancesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 原始价格，单位：分
		OriginalPrice *uint64 `json:"OriginalPrice" name:"OriginalPrice"`

		// 折后价格，单位：分
		Price *uint64 `json:"Price" name:"Price"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *InquiryPriceCreateDBInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceCreateDBInstancesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InquiryPriceRenewDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 实例ID
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`

	// 续费周期，按月计算，最大不超过48
	Period *int64 `json:"Period" name:"Period"`
}

func (r *InquiryPriceRenewDBInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceRenewDBInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InquiryPriceRenewDBInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 总费用，打折前的。比如24650表示246.5元
		OriginalPrice *int64 `json:"OriginalPrice" name:"OriginalPrice"`

		// 实际需要付款金额。比如24650表示246.5元
		Price *int64 `json:"Price" name:"Price"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *InquiryPriceRenewDBInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceRenewDBInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InquiryPriceUpgradeDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 实例的磁盘大小，单位GB
	Storage *int64 `json:"Storage" name:"Storage"`

	// 实例的内存大小，单位GB
	Memory *int64 `json:"Memory" name:"Memory"`

	// 实例ID，形如postgres-hez4fh0v
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`

	// 实例计费类型，预付费或者后付费。PREPAID-预付费。目前只支持预付费。
	InstanceChargeType *string `json:"InstanceChargeType" name:"InstanceChargeType"`
}

func (r *InquiryPriceUpgradeDBInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceUpgradeDBInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InquiryPriceUpgradeDBInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 总费用，打折前的
		OriginalPrice *int64 `json:"OriginalPrice" name:"OriginalPrice"`

		// 实际需要付款金额
		Price *int64 `json:"Price" name:"Price"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *InquiryPriceUpgradeDBInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceUpgradeDBInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAccountRemarkRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如postgres-4wdeb0zv
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`

	// 实例用户名
	UserName *string `json:"UserName" name:"UserName"`

	// 用户UserName对应的新备注
	Remark *string `json:"Remark" name:"Remark"`
}

func (r *ModifyAccountRemarkRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAccountRemarkRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAccountRemarkResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyAccountRemarkResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAccountRemarkResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBInstanceNameRequest struct {
	*tchttp.BaseRequest

	// 数据库实例ID，形如postgres-6fego161
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`

	// 新的数据库实例名字
	InstanceName *string `json:"InstanceName" name:"InstanceName"`
}

func (r *ModifyDBInstanceNameRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBInstanceNameRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBInstanceNameResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDBInstanceNameResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBInstanceNameResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBInstancesProjectRequest struct {
	*tchttp.BaseRequest

	// postgresql实例ID数组
	DBInstanceIdSet []*string `json:"DBInstanceIdSet" name:"DBInstanceIdSet" list`

	// postgresql实例所属新项目的ID
	ProjectId *string `json:"ProjectId" name:"ProjectId"`
}

func (r *ModifyDBInstancesProjectRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBInstancesProjectRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBInstancesProjectResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 转移项目成功的实例个数
		Count *int64 `json:"Count" name:"Count"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDBInstancesProjectResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBInstancesProjectResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type NormalQueryItem struct {

	// 用户名
	UserName *string `json:"UserName" name:"UserName"`

	// 调用次数
	Calls *int64 `json:"Calls" name:"Calls"`

	// 粒度点
	CallsGrids []*int64 `json:"CallsGrids" name:"CallsGrids" list`

	// 花费总时间
	CostTime *float64 `json:"CostTime" name:"CostTime"`

	// 影响的行数
	Rows *int64 `json:"Rows" name:"Rows"`

	// 花费最小时间
	MinCostTime *float64 `json:"MinCostTime" name:"MinCostTime"`

	// 花费最大时间
	MaxCostTime *float64 `json:"MaxCostTime" name:"MaxCostTime"`

	// 最早一条慢SQL时间
	FirstTime *string `json:"FirstTime" name:"FirstTime"`

	// 最晚一条慢SQL时间
	LastTime *string `json:"LastTime" name:"LastTime"`

	// 读共享内存块数
	SharedReadBlks *int64 `json:"SharedReadBlks" name:"SharedReadBlks"`

	// 写共享内存块数
	SharedWriteBlks *int64 `json:"SharedWriteBlks" name:"SharedWriteBlks"`

	// 读io总耗时
	ReadCostTime *int64 `json:"ReadCostTime" name:"ReadCostTime"`

	// 写io总耗时
	WriteCostTime *int64 `json:"WriteCostTime" name:"WriteCostTime"`

	// 数据库名字
	DatabaseName *string `json:"DatabaseName" name:"DatabaseName"`

	// 脱敏后的慢SQL
	NormalQuery *string `json:"NormalQuery" name:"NormalQuery"`
}

type OpenDBExtranetAccessRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如postgres-hez4fh0v
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`
}

func (r *OpenDBExtranetAccessRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *OpenDBExtranetAccessRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type OpenDBExtranetAccessResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务流程ID
		FlowId *int64 `json:"FlowId" name:"FlowId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *OpenDBExtranetAccessResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *OpenDBExtranetAccessResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type PgDeal struct {

	// 订单名
	DealName *string `json:"DealName" name:"DealName"`

	// 所属用户
	OwnerUin *string `json:"OwnerUin" name:"OwnerUin"`

	// 订单涉及多少个实例
	Count *int64 `json:"Count" name:"Count"`

	// 付费模式。1-预付费；0-后付费
	PayMode *int64 `json:"PayMode" name:"PayMode"`

	// 异步任务流程ID
	FlowId *int64 `json:"FlowId" name:"FlowId"`

	// 实例ID数组
	DBInstanceIdSet []*string `json:"DBInstanceIdSet" name:"DBInstanceIdSet" list`
}

type RegionInfo struct {

	// 该地域对应的英文名称
	Region *string `json:"Region" name:"Region"`

	// 该地域对应的中文名称
	RegionName *string `json:"RegionName" name:"RegionName"`

	// 该地域对应的数字编号
	RegionId *uint64 `json:"RegionId" name:"RegionId"`

	// 可用状态，UNAVAILABLE表示不可用，AVAILABLE表示可用
	RegionState *string `json:"RegionState" name:"RegionState"`
}

type RenewInstanceRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如postgres-6fego161
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`

	// 续费多少个月
	Period *int64 `json:"Period" name:"Period"`

	// 是否自动使用代金券,1是,0否，默认不使用
	AutoVoucher *int64 `json:"AutoVoucher" name:"AutoVoucher"`

	// 代金券ID列表，目前仅支持指定一张代金券
	VoucherIds []*string `json:"VoucherIds" name:"VoucherIds" list`
}

func (r *RenewInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RenewInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RenewInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 订单名
		DealName *string `json:"DealName" name:"DealName"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RenewInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RenewInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ResetAccountPasswordRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如postgres-4wdeb0zv
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`

	// 实例账户名
	UserName *string `json:"UserName" name:"UserName"`

	// UserName账户对应的新密码
	Password *string `json:"Password" name:"Password"`
}

func (r *ResetAccountPasswordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ResetAccountPasswordRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ResetAccountPasswordResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ResetAccountPasswordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ResetAccountPasswordResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RestartDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如postgres-6r233v55
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`
}

func (r *RestartDBInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RestartDBInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RestartDBInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步流程ID
		FlowId *int64 `json:"FlowId" name:"FlowId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RestartDBInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RestartDBInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SetAutoRenewFlagRequest struct {
	*tchttp.BaseRequest

	// 实例ID数组
	DBInstanceIdSet []*string `json:"DBInstanceIdSet" name:"DBInstanceIdSet" list`

	// 续费标记。0-正常续费；1-自动续费；2-到期不续费
	AutoRenewFlag *int64 `json:"AutoRenewFlag" name:"AutoRenewFlag"`
}

func (r *SetAutoRenewFlagRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SetAutoRenewFlagRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SetAutoRenewFlagResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设置成功的实例个数
		Count *int64 `json:"Count" name:"Count"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *SetAutoRenewFlagResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SetAutoRenewFlagResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SlowlogDetail struct {

	// 花费总时间
	TotalTime *float64 `json:"TotalTime" name:"TotalTime"`

	// 调用总次数
	TotalCalls *int64 `json:"TotalCalls" name:"TotalCalls"`

	// 脱敏后的慢SQL列表
	NormalQueries []*NormalQueryItem `json:"NormalQueries" name:"NormalQueries" list`
}

type SpecInfo struct {

	// 地域英文编码，对应RegionSet的Region字段
	Region *string `json:"Region" name:"Region"`

	// 区域英文编码，对应ZoneSet的Zone字段
	Zone *string `json:"Zone" name:"Zone"`

	// 规格详细信息列表
	SpecItemInfoList []*SpecItemInfo `json:"SpecItemInfoList" name:"SpecItemInfoList" list`
}

type SpecItemInfo struct {

	// 规格ID
	SpecCode *string `json:"SpecCode" name:"SpecCode"`

	// PostgreSQL的内核版本编号
	Version *string `json:"Version" name:"Version"`

	// 内核编号对应的完整版本名称
	VersionName *string `json:"VersionName" name:"VersionName"`

	// CPU核数
	Cpu *uint64 `json:"Cpu" name:"Cpu"`

	// 内存大小，单位：MB
	Memory *uint64 `json:"Memory" name:"Memory"`

	// 该规格所支持最大存储容量，单位：GB
	MaxStorage *uint64 `json:"MaxStorage" name:"MaxStorage"`

	// 该规格所支持最小存储容量，单位：GB
	MinStorage *uint64 `json:"MinStorage" name:"MinStorage"`

	// 该规格的预估QPS
	Qps *uint64 `json:"Qps" name:"Qps"`

	// 该规格对应的计费ID
	Pid *uint64 `json:"Pid" name:"Pid"`
}

type UpgradeDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 升级后的实例内存大小，单位GB
	Memory *int64 `json:"Memory" name:"Memory"`

	// 升级后的实例磁盘大小，单位GB
	Storage *int64 `json:"Storage" name:"Storage"`

	// 实例ID，形如postgres-lnp6j617
	DBInstanceId *string `json:"DBInstanceId" name:"DBInstanceId"`

	// 是否自动使用代金券,1是,0否，默认不使用
	AutoVoucher *int64 `json:"AutoVoucher" name:"AutoVoucher"`

	// 代金券ID列表，目前仅支持指定一张代金券
	VoucherIds []*string `json:"VoucherIds" name:"VoucherIds" list`
}

func (r *UpgradeDBInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpgradeDBInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpgradeDBInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 交易名字。
		DealName *string `json:"DealName" name:"DealName"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UpgradeDBInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpgradeDBInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Xlog struct {

	// 备份文件唯一标识
	Id *int64 `json:"Id" name:"Id"`

	// 文件生成的开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 文件生成的结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 内网下载地址
	InternalAddr *string `json:"InternalAddr" name:"InternalAddr"`

	// 外网下载地址
	ExternalAddr *string `json:"ExternalAddr" name:"ExternalAddr"`
}

type ZoneInfo struct {

	// 该可用区的英文名称
	Zone *string `json:"Zone" name:"Zone"`

	// 该可用区的中文名称
	ZoneName *string `json:"ZoneName" name:"ZoneName"`

	// 该可用区对应的数字编号
	ZoneId *uint64 `json:"ZoneId" name:"ZoneId"`

	// 可用状态，UNAVAILABLE表示不可用，AVAILABLE表示可用
	ZoneState *string `json:"ZoneState" name:"ZoneState"`
}
