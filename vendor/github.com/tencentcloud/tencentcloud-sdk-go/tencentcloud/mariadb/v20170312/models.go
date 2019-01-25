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

type CloneAccountRequest struct {
	*tchttp.BaseRequest

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 源用户账户名
	SrcUser *string `json:"SrcUser" name:"SrcUser"`

	// 源用户HOST
	SrcHost *string `json:"SrcHost" name:"SrcHost"`

	// 目的用户账户名
	DstUser *string `json:"DstUser" name:"DstUser"`

	// 目的用户HOST
	DstHost *string `json:"DstHost" name:"DstHost"`

	// 目的用户账户描述
	DstDesc *string `json:"DstDesc" name:"DstDesc"`
}

func (r *CloneAccountRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CloneAccountRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CloneAccountResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务流程ID。
		FlowId *uint64 `json:"FlowId" name:"FlowId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CloneAccountResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CloneAccountResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CloseDBExtranetAccessRequest struct {
	*tchttp.BaseRequest

	// 待关闭外网访问的实例ID。形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
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

		// 异步任务Id，可通过 DescribeFlow 查询任务状态。
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

type ConstraintRange struct {

	// 约束类型为section时的最小值
	Min *string `json:"Min" name:"Min"`

	// 约束类型为section时的最大值
	Max *string `json:"Max" name:"Max"`
}

type CopyAccountPrivilegesRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 源用户名
	SrcUserName *string `json:"SrcUserName" name:"SrcUserName"`

	// 源用户允许的访问 host
	SrcHost *string `json:"SrcHost" name:"SrcHost"`

	// 目的用户名
	DstUserName *string `json:"DstUserName" name:"DstUserName"`

	// 目的用户允许的访问 host
	DstHost *string `json:"DstHost" name:"DstHost"`

	// 源账号的 ReadOnly 属性
	SrcReadOnly *string `json:"SrcReadOnly" name:"SrcReadOnly"`

	// 目的账号的 ReadOnly 属性
	DstReadOnly *string `json:"DstReadOnly" name:"DstReadOnly"`
}

func (r *CopyAccountPrivilegesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CopyAccountPrivilegesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CopyAccountPrivilegesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CopyAccountPrivilegesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CopyAccountPrivilegesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateAccountRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 登录用户名，由字幕、数字、下划线和连字符组成，长度为1~32位。
	UserName *string `json:"UserName" name:"UserName"`

	// 可以登录的主机，与mysql 账号的 host 格式一致，可以支持通配符，例如 %，10.%，10.20.%。
	Host *string `json:"Host" name:"Host"`

	// 账号密码，由字母、数字或常见符号组成，不能包含分号、单引号和双引号，长度为6~32位。
	Password *string `json:"Password" name:"Password"`

	// 是否创建为只读账号，0：否， 1：该账号的sql请求优先选择备机执行，备机不可用时选择主机执行，2：优先选择备机执行，备机不可用时操作失败。
	ReadOnly *int64 `json:"ReadOnly" name:"ReadOnly"`

	// 账号备注，可以包含中文、英文字符、常见符号和数字，长度为0~256字符
	Description *string `json:"Description" name:"Description"`
}

func (r *CreateAccountRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateAccountRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateAccountResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例Id，透传入参。
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 用户名，透传入参。
		UserName *string `json:"UserName" name:"UserName"`

		// 允许访问的 host，透传入参。
		Host *string `json:"Host" name:"Host"`

		// 透传入参。
		ReadOnly *int64 `json:"ReadOnly" name:"ReadOnly"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateAccountResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateAccountResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 实例节点可用区分布，最多可填两个可用区。当分片规格为一主两从时，其中两个节点在第一个可用区。
	Zones []*string `json:"Zones" name:"Zones" list`

	// 节点个数大小，可以通过 DescribeDBInstanceSpecs
	//  查询实例规格获得。
	NodeCount *int64 `json:"NodeCount" name:"NodeCount"`

	// 内存大小，单位：GB，可以通过 DescribeDBInstanceSpecs
	//  查询实例规格获得。
	Memory *int64 `json:"Memory" name:"Memory"`

	// 存储空间大小，单位：GB，可以通过 DescribeDBInstanceSpecs
	//  查询实例规格获得不同内存大小对应的磁盘规格下限和上限。
	Storage *int64 `json:"Storage" name:"Storage"`

	// 欲购买的时长，单位：月。
	Period *int64 `json:"Period" name:"Period"`

	// 欲购买的数量，默认查询购买1个实例的价格。
	Count *int64 `json:"Count" name:"Count"`

	// 是否自动使用代金券进行支付，默认不使用。
	AutoVoucher *bool `json:"AutoVoucher" name:"AutoVoucher"`

	// 代金券ID列表，目前仅支持指定一张代金券。
	VoucherIds []*string `json:"VoucherIds" name:"VoucherIds" list`

	// 虚拟私有网络 ID，不传或传 0 表示创建为基础网络
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 虚拟私有网络子网 ID，VpcId 不为0时必填
	SubnetId *string `json:"SubnetId" name:"SubnetId"`

	// 项目 ID，可以通过查看项目列表获取，不传则关联到默认项目
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 数据库引擎版本，当前可选：10.0.10，10.1.9，5.7.17
	DbVersionId *string `json:"DbVersionId" name:"DbVersionId"`
}

func (r *CreateDBInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDBInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDBInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 长订单号。可以据此调用 DescribeOrders
	//  查询订单详细信息，或在支付失败时调用用户账号相关接口进行支付。
		DealName *string `json:"DealName" name:"DealName"`

		// 订单对应的实例 ID 列表，如果此处没有返回实例 ID，可以通过订单查询接口获取。还可通过实例查询接口查询实例是否创建完成。
		InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateDBInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDBInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DBAccount struct {

	// 用户名
	UserName *string `json:"UserName" name:"UserName"`

	// 用户可以从哪台主机登录（对应 MySQL 用户的 host 字段，UserName + Host 唯一标识一个用户，IP形式，IP段以%结尾；支持填入%；为空默认等于%）
	Host *string `json:"Host" name:"Host"`

	// 用户备注信息
	Description *string `json:"Description" name:"Description"`

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 最后更新时间
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`

	// 只读标记，0：否， 1：该账号的sql请求优先选择备机执行，备机不可用时选择主机执行，2：优先选择备机执行，备机不可用时操作失败。
	ReadOnly *int64 `json:"ReadOnly" name:"ReadOnly"`

	// 该字段对只读帐号有意义，表示选择主备延迟小于该值的备机
	DelayThresh *int64 `json:"DelayThresh" name:"DelayThresh"`
}

type DBBackupTimeConfig struct {

	// 实例 Id
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 每天备份执行的区间的开始时间，格式 mm:ss，形如 22:00
	StartBackupTime *string `json:"StartBackupTime" name:"StartBackupTime"`

	// 每天备份执行的区间的结束时间，格式 mm:ss，形如 23:00
	EndBackupTime *string `json:"EndBackupTime" name:"EndBackupTime"`
}

type DBInstance struct {

	// 实例 Id，唯一标识一个 TDSQL 实例
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 实例名称，用户可修改
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// 实例所属应用 Id
	AppId *int64 `json:"AppId" name:"AppId"`

	// 实例所属项目 Id
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 实例所在地域名称，如 ap-shanghai
	Region *string `json:"Region" name:"Region"`

	// 实例所在可用区名称，如 ap-shanghai-1
	Zone *string `json:"Zone" name:"Zone"`

	// 私有网络 Id，基础网络时为 0
	VpcId *int64 `json:"VpcId" name:"VpcId"`

	// 子网 Id，基础网络时为 0
	SubnetId *int64 `json:"SubnetId" name:"SubnetId"`

	// 实例状态：0 创建中，1 流程处理中， 2 运行中，3 实例未初始化，-1 实例已隔离，-2 实例已删除
	Status *int64 `json:"Status" name:"Status"`

	// 内网 IP 地址
	Vip *string `json:"Vip" name:"Vip"`

	// 内网端口
	Vport *int64 `json:"Vport" name:"Vport"`

	// 外网访问的域名，公网可解析
	WanDomain *string `json:"WanDomain" name:"WanDomain"`

	// 外网 IP 地址，公网可访问
	WanVip *string `json:"WanVip" name:"WanVip"`

	// 外网端口
	WanPort *int64 `json:"WanPort" name:"WanPort"`

	// 实例创建时间，格式为 2006-01-02 15:04:05
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 实例最后更新时间，格式为 2006-01-02 15:04:05
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`

	// 自动续费标志：0 否，1 是
	AutoRenewFlag *int64 `json:"AutoRenewFlag" name:"AutoRenewFlag"`

	// 实例到期时间，格式为 2006-01-02 15:04:05
	PeriodEndTime *string `json:"PeriodEndTime" name:"PeriodEndTime"`

	// 实例所属账号
	Uin *string `json:"Uin" name:"Uin"`

	// TDSQL 版本信息
	TdsqlVersion *string `json:"TdsqlVersion" name:"TdsqlVersion"`

	// 实例内存大小，单位 GB
	Memory *int64 `json:"Memory" name:"Memory"`

	// 实例存储大小，单位 GB
	Storage *int64 `json:"Storage" name:"Storage"`

	// 字符串型的私有网络Id
	UniqueVpcId *string `json:"UniqueVpcId" name:"UniqueVpcId"`

	// 字符串型的私有网络子网Id
	UniqueSubnetId *string `json:"UniqueSubnetId" name:"UniqueSubnetId"`

	// 原始实例ID（过时字段，请勿依赖该值）
	OriginSerialId *string `json:"OriginSerialId" name:"OriginSerialId"`

	// 节点数，2为一主一从，3为一主二从
	NodeCount *uint64 `json:"NodeCount" name:"NodeCount"`

	// 是否临时实例，0为否，非0为是
	IsTmp *uint64 `json:"IsTmp" name:"IsTmp"`

	// 独享集群Id，为空表示为普通实例
	ExclusterId *string `json:"ExclusterId" name:"ExclusterId"`

	// 数字实例Id（过时字段，请勿依赖该值）
	Id *uint64 `json:"Id" name:"Id"`

	// 产品类型 Id
	Pid *int64 `json:"Pid" name:"Pid"`

	// 最大 Qps 值
	Qps *int64 `json:"Qps" name:"Qps"`

	// 付费模式
	Paymode *string `json:"Paymode" name:"Paymode"`

	// 实例处于异步任务时的异步任务流程ID
	Locker *int64 `json:"Locker" name:"Locker"`

	// 实例目前运行状态描述
	StatusDesc *string `json:"StatusDesc" name:"StatusDesc"`
}

type DBParamValue struct {

	// 参数名称
	Param *string `json:"Param" name:"Param"`

	// 参数值
	Value *string `json:"Value" name:"Value"`
}

type Deal struct {

	// 订单号
	DealName *string `json:"DealName" name:"DealName"`

	// 所属账号
	OwnerUin *string `json:"OwnerUin" name:"OwnerUin"`

	// 商品数量
	Count *int64 `json:"Count" name:"Count"`

	// 关联的流程 Id，可用于查询流程执行状态
	FlowId *int64 `json:"FlowId" name:"FlowId"`

	// 只有创建实例的订单会填充该字段，表示该订单创建的实例的 ID。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

	// 付费模式，0后付费/1预付费
	PayMode *int64 `json:"PayMode" name:"PayMode"`
}

type DeleteAccountRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 用户名
	UserName *string `json:"UserName" name:"UserName"`

	// 用户允许的访问 host
	Host *string `json:"Host" name:"Host"`
}

func (r *DeleteAccountRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteAccountRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteAccountResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteAccountResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteAccountResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAccountPrivilegesRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 登录用户名。
	UserName *string `json:"UserName" name:"UserName"`

	// 用户允许的访问 host，用户名+host唯一确定一个账号。
	Host *string `json:"Host" name:"Host"`

	// 数据库名。如果为 \*，表示查询全局权限（即 \*.\*），此时忽略 Type 和 Object 参数
	DbName *string `json:"DbName" name:"DbName"`

	// 类型,可以填入 table 、 view 、 proc 、 func 和 \*。当 DbName 为具体数据库名，Type为 \* 时，表示查询该数据库权限（即db.\*），此时忽略 Object 参数
	Type *string `json:"Type" name:"Type"`

	// 具体的 Type 的名称，比如 Type 为 table 时就是具体的表名。DbName 和 Type 都为具体名称，则 Object 表示具体对象名，不能为 \* 或者为空
	Object *string `json:"Object" name:"Object"`

	// 当 Type=table 时，ColName 为 \* 表示查询表的权限，如果为具体字段名，表示查询对应字段的权限
	ColName *string `json:"ColName" name:"ColName"`
}

func (r *DescribeAccountPrivilegesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAccountPrivilegesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAccountPrivilegesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例Id
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 权限列表。
		Privileges []*string `json:"Privileges" name:"Privileges" list`

		// 数据库账号用户名
		UserName *string `json:"UserName" name:"UserName"`

		// 数据库账号Host
		Host *string `json:"Host" name:"Host"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeAccountPrivilegesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAccountPrivilegesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAccountsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
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

		// 实例ID，透传入参。
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 实例用户列表。
		Users []*DBAccount `json:"Users" name:"Users" list`

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

type DescribeBackupTimeRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`
}

func (r *DescribeBackupTimeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBackupTimeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBackupTimeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 返回的配置数量
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 实例备份时间配置信息
		Items *DBBackupTimeConfig `json:"Items" name:"Items"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBackupTimeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBackupTimeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstanceSpecsRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeDBInstanceSpecsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBInstanceSpecsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstanceSpecsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 按机型分类的可售卖规格列表
		Specs []*InstanceSpec `json:"Specs" name:"Specs" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBInstanceSpecsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBInstanceSpecsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstancesRequest struct {
	*tchttp.BaseRequest

	// 按照一个或者多个实例 ID 查询。实例 ID 形如：tdsql-ow728lmc。每次请求的实例的上限为100。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

	// 搜索的字段名，当前支持的值有：instancename、vip、all。传 instancename 表示按实例名进行搜索；传 vip 表示按内网IP进行搜索；传 all 将会按实例ID、实例名和内网IP进行搜索。
	SearchName *string `json:"SearchName" name:"SearchName"`

	// 搜索的关键字，支持模糊搜索。多个关键字使用换行符（'\n'）分割。
	SearchKey *string `json:"SearchKey" name:"SearchKey"`

	// 按项目 ID 查询
	ProjectIds []*int64 `json:"ProjectIds" name:"ProjectIds" list`

	// 是否根据 VPC 网络来搜索
	IsFilterVpc *bool `json:"IsFilterVpc" name:"IsFilterVpc"`

	// 私有网络 ID， IsFilterVpc 为 1 时有效
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 私有网络的子网 ID， IsFilterVpc 为 1 时有效
	SubnetId *string `json:"SubnetId" name:"SubnetId"`

	// 排序字段， projectId， createtime， instancename 三者之一
	OrderBy *string `json:"OrderBy" name:"OrderBy"`

	// 排序类型， desc 或者 asc
	OrderByType *string `json:"OrderByType" name:"OrderByType"`

	// 偏移量，默认为 0
	Offset *int64 `json:"Offset" name:"Offset"`

	// 返回数量，默认为 20，最大值为 100。
	Limit *int64 `json:"Limit" name:"Limit"`

	// 按 OriginSerialId 查询
	OriginSerialIds []*string `json:"OriginSerialIds" name:"OriginSerialIds" list`

	// 标识是否使用ExclusterType字段, false不使用，true使用
	IsFilterExcluster *bool `json:"IsFilterExcluster" name:"IsFilterExcluster"`

	// 1非独享集群，2独享集群， 0全部
	ExclusterType *int64 `json:"ExclusterType" name:"ExclusterType"`
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

		// 符合条件的实例数量
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 实例详细信息列表
		Instances []*DBInstance `json:"Instances" name:"Instances" list`

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

type DescribeDBLogFilesRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 请求日志类型，取值只能为1、2、3或者4。1-binlog，2-冷备，3-errlog，4-slowlog。
	Type *uint64 `json:"Type" name:"Type"`
}

func (r *DescribeDBLogFilesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBLogFilesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBLogFilesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例 ID，形如：tdsql-ow728lmc。
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 请求日志类型，取值只能为1、2、3或者4。1-binlog，2-冷备，3-errlog，4-slowlog。
		Type *uint64 `json:"Type" name:"Type"`

		// 请求日志总数
		Total *uint64 `json:"Total" name:"Total"`

		// 包含uri、length、mtime（修改时间）等信息
		Files []*LogFileInfo `json:"Files" name:"Files" list`

		// 如果是VPC网络的实例，做用本前缀加上URI为下载地址
		VpcPrefix *string `json:"VpcPrefix" name:"VpcPrefix"`

		// 如果是普通网络的实例，做用本前缀加上URI为下载地址
		NormalPrefix *string `json:"NormalPrefix" name:"NormalPrefix"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBLogFilesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBLogFilesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBParametersRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *DescribeDBParametersRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBParametersRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBParametersResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例 ID，形如：tdsql-ow728lmc。
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 请求DB的当前参数值
		Params []*ParamDesc `json:"Params" name:"Params" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBParametersResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBParametersResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBPerformanceDetailsRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 开始日期，格式yyyy-mm-dd
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 结束日期，格式yyyy-mm-dd
	EndTime *string `json:"EndTime" name:"EndTime"`
}

func (r *DescribeDBPerformanceDetailsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBPerformanceDetailsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBPerformanceDetailsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 主节点性能监控数据
		Master *PerformanceMonitorSet `json:"Master" name:"Master"`

		// 备机1性能监控数据
		Slave1 *PerformanceMonitorSet `json:"Slave1" name:"Slave1"`

		// 备机2性能监控数据，如果实例是一主一从，则没有该字段
		Slave2 *PerformanceMonitorSet `json:"Slave2" name:"Slave2"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBPerformanceDetailsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBPerformanceDetailsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBPerformanceRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 开始日期，格式yyyy-mm-dd
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 结束日期，格式yyyy-mm-dd
	EndTime *string `json:"EndTime" name:"EndTime"`
}

func (r *DescribeDBPerformanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBPerformanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBPerformanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 慢查询数
		LongQuery *MonitorData `json:"LongQuery" name:"LongQuery"`

		// 查询操作数SELECT
		SelectTotal *MonitorData `json:"SelectTotal" name:"SelectTotal"`

		// 更新操作数UPDATE
		UpdateTotal *MonitorData `json:"UpdateTotal" name:"UpdateTotal"`

		// 插入操作数INSERT
		InsertTotal *MonitorData `json:"InsertTotal" name:"InsertTotal"`

		// 删除操作数DELETE
		DeleteTotal *MonitorData `json:"DeleteTotal" name:"DeleteTotal"`

		// 缓存命中率
		MemHitRate *MonitorData `json:"MemHitRate" name:"MemHitRate"`

		// 磁盘每秒IO次数
		DiskIops *MonitorData `json:"DiskIops" name:"DiskIops"`

		// 活跃连接数
		ConnActive *MonitorData `json:"ConnActive" name:"ConnActive"`

		// 是否发生主备切换，1为发生，0否
		IsMasterSwitched *MonitorData `json:"IsMasterSwitched" name:"IsMasterSwitched"`

		// 主备延迟
		SlaveDelay *MonitorData `json:"SlaveDelay" name:"SlaveDelay"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBPerformanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBPerformanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBResourceUsageDetailsRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 开始日期，格式yyyy-mm-dd
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 结束日期，格式yyyy-mm-dd
	EndTime *string `json:"EndTime" name:"EndTime"`
}

func (r *DescribeDBResourceUsageDetailsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBResourceUsageDetailsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBResourceUsageDetailsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 主节点资源使用情况监控数据
		Master *ResourceUsageMonitorSet `json:"Master" name:"Master"`

		// 备机1资源使用情况监控数据
		Slave1 *ResourceUsageMonitorSet `json:"Slave1" name:"Slave1"`

		// 备机2资源使用情况监控数据
		Slave2 *ResourceUsageMonitorSet `json:"Slave2" name:"Slave2"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBResourceUsageDetailsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBResourceUsageDetailsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBResourceUsageRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 开始日期，格式yyyy-mm-dd
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 结束日期，格式yyyy-mm-dd
	EndTime *string `json:"EndTime" name:"EndTime"`
}

func (r *DescribeDBResourceUsageRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBResourceUsageRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBResourceUsageResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// binlog日志磁盘可用空间,单位GB
		BinlogDiskAvailable *MonitorData `json:"BinlogDiskAvailable" name:"BinlogDiskAvailable"`

		// 磁盘可用空间,单位GB
		DataDiskAvailable *MonitorData `json:"DataDiskAvailable" name:"DataDiskAvailable"`

		// CPU利用率
		CpuUsageRate *MonitorData `json:"CpuUsageRate" name:"CpuUsageRate"`

		// 内存可用空间,单位GB
		MemAvailable *MonitorData `json:"MemAvailable" name:"MemAvailable"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBResourceUsageResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBResourceUsageResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBSlowLogsRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 从结果的第几条数据开始返回
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 返回的结果条数
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 查询的起始时间，形如2016-07-23 14:55:20
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 查询的结束时间，形如2016-08-22 14:55:20
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 要查询的具体数据库名称
	Db *string `json:"Db" name:"Db"`

	// 排序指标，取值为query_time_sum或者query_count
	OrderBy *string `json:"OrderBy" name:"OrderBy"`

	// 排序类型，desc或者asc
	OrderByType *string `json:"OrderByType" name:"OrderByType"`
}

func (r *DescribeDBSlowLogsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBSlowLogsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBSlowLogsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 慢查询日志数据
		Data []*SlowLogData `json:"Data" name:"Data" list`

		// 所有语句锁时间总和
		LockTimeSum *float64 `json:"LockTimeSum" name:"LockTimeSum"`

		// 所有语句查询总次数
		QueryCount *int64 `json:"QueryCount" name:"QueryCount"`

		// 总记录数
		Total *int64 `json:"Total" name:"Total"`

		// 所有语句查询时间总和
		QueryTimeSum *float64 `json:"QueryTimeSum" name:"QueryTimeSum"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBSlowLogsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBSlowLogsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeFlowRequest struct {
	*tchttp.BaseRequest

	// 异步请求接口返回的任务流程号。
	FlowId *int64 `json:"FlowId" name:"FlowId"`
}

func (r *DescribeFlowRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeFlowRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeFlowResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 流程状态，0：成功，1：失败，2：运行中
		Status *int64 `json:"Status" name:"Status"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeFlowResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeFlowResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLogFileRetentionPeriodRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *DescribeLogFileRetentionPeriodRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLogFileRetentionPeriodRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLogFileRetentionPeriodResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例 ID，形如：tdsql-ow728lmc。
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 日志备份天数
		Days *uint64 `json:"Days" name:"Days"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeLogFileRetentionPeriodResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLogFileRetentionPeriodResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeOrdersRequest struct {
	*tchttp.BaseRequest

	// 待查询的长订单号列表，创建实例、续费实例、扩容实例接口返回。
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

		// 返回的订单数量。
		TotalCount []*int64 `json:"TotalCount" name:"TotalCount" list`

		// 订单信息列表。
		Deals []*Deal `json:"Deals" name:"Deals" list`

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

type DescribePriceRequest struct {
	*tchttp.BaseRequest

	// 欲新购实例的可用区ID。
	Zone *string `json:"Zone" name:"Zone"`

	// 实例节点个数，可以通过 DescribeDBInstanceSpecs
	//  查询实例规格获得。
	NodeCount *int64 `json:"NodeCount" name:"NodeCount"`

	// 内存大小，单位：GB，可以通过 DescribeDBInstanceSpecs
	//  查询实例规格获得。
	Memory *int64 `json:"Memory" name:"Memory"`

	// 存储空间大小，单位：GB，可以通过 DescribeDBInstanceSpecs
	//  查询实例规格获得不同内存大小对应的磁盘规格下限和上限。
	Storage *int64 `json:"Storage" name:"Storage"`

	// 欲购买的时长，单位：月。
	Period *int64 `json:"Period" name:"Period"`

	// 欲购买的数量，默认查询购买1个实例的价格。
	Count *int64 `json:"Count" name:"Count"`
}

func (r *DescribePriceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePriceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePriceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 原价，单位：分
		OriginalPrice *int64 `json:"OriginalPrice" name:"OriginalPrice"`

		// 实际价格，单位：分。受折扣等影响，可能和原价不同。
		Price *int64 `json:"Price" name:"Price"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribePriceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePriceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeRenewalPriceRequest struct {
	*tchttp.BaseRequest

	// 待续费的实例ID。形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 续费时长，单位：月。不传则默认为1个月。
	Period *int64 `json:"Period" name:"Period"`
}

func (r *DescribeRenewalPriceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeRenewalPriceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeRenewalPriceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 原价，单位：分
		OriginalPrice *int64 `json:"OriginalPrice" name:"OriginalPrice"`

		// 实际价格，单位：分。受折扣等影响，可能和原价不同。
		Price *int64 `json:"Price" name:"Price"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeRenewalPriceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeRenewalPriceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSaleInfoRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeSaleInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSaleInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSaleInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 可售卖地域信息列表
		RegionList []*RegionInfo `json:"RegionList" name:"RegionList" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeSaleInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSaleInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSqlLogsRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// SQL日志偏移。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 拉取数量（0-1000，为0时拉取总数信息）。
	Limit *int64 `json:"Limit" name:"Limit"`
}

func (r *DescribeSqlLogsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSqlLogsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSqlLogsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 当前消息队列中的sql日志条目数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 消息队列中的sql日志起始偏移。
		StartOffset *uint64 `json:"StartOffset" name:"StartOffset"`

		// 消息队列中的sql日志结束偏移。
		EndOffset *uint64 `json:"EndOffset" name:"EndOffset"`

		// 返回的第一条sql日志的偏移。
		Offset *uint64 `json:"Offset" name:"Offset"`

		// 返回的sql日志数量。
		Count *uint64 `json:"Count" name:"Count"`

		// Sql日志列表。
		SqlItems []*SqlLogItem `json:"SqlItems" name:"SqlItems" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeSqlLogsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSqlLogsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeUpgradePriceRequest struct {
	*tchttp.BaseRequest

	// 待升级的实例ID。形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 内存大小，单位：GB，可以通过 DescribeDBInstanceSpecs
	//  查询实例规格获得。
	Memory *int64 `json:"Memory" name:"Memory"`

	// 存储空间大小，单位：GB，可以通过 DescribeDBInstanceSpecs
	//  查询实例规格获得不同内存大小对应的磁盘规格下限和上限。
	Storage *int64 `json:"Storage" name:"Storage"`
}

func (r *DescribeUpgradePriceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeUpgradePriceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeUpgradePriceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 原价，单位：分
		OriginalPrice *int64 `json:"OriginalPrice" name:"OriginalPrice"`

		// 实际价格，单位：分。受折扣等影响，可能和原价不同。
		Price *int64 `json:"Price" name:"Price"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeUpgradePriceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeUpgradePriceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GrantAccountPrivilegesRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 登录用户名。
	UserName *string `json:"UserName" name:"UserName"`

	// 用户允许的访问 host，用户名+host唯一确定一个账号。
	Host *string `json:"Host" name:"Host"`

	// 数据库名。如果为 \*，表示设置全局权限（即 \*.\*），此时忽略 Type 和 Object 参数
	DbName *string `json:"DbName" name:"DbName"`

	// 全局权限： SELECT，INSERT，UPDATE，DELETE，CREATE，DROP，REFERENCES，INDEX，ALTER，CREATE TEMPORARY TABLES，LOCK TABLES，EXECUTE，CREATE VIEW，SHOW VIEW，CREATE ROUTINE，ALTER ROUTINE，EVENT，TRIGGER，SHOW DATABASES 
	// 库权限： SELECT，INSERT，UPDATE，DELETE，CREATE，DROP，REFERENCES，INDEX，ALTER，CREATE TEMPORARY TABLES，LOCK TABLES，EXECUTE，CREATE VIEW，SHOW VIEW，CREATE ROUTINE，ALTER ROUTINE，EVENT，TRIGGER 
	// 表/视图权限： SELECT，INSERT，UPDATE，DELETE，CREATE，DROP，REFERENCES，INDEX，ALTER，CREATE VIEW，SHOW VIEW，TRIGGER 
	// 存储过程/函数权限： ALTER ROUTINE，EXECUTE 
	// 字段权限： INSERT，REFERENCES，SELECT，UPDATE
	Privileges []*string `json:"Privileges" name:"Privileges" list`

	// 类型,可以填入 table 、 view 、 proc 、 func 和 \*。当 DbName 为具体数据库名，Type为 \* 时，表示设置该数据库权限（即db.\*），此时忽略 Object 参数
	Type *string `json:"Type" name:"Type"`

	// 具体的 Type 的名称，比如 Type 为 table 时就是具体的表名。DbName 和 Type 都为具体名称，则 Object 表示具体对象名，不能为 \* 或者为空
	Object *string `json:"Object" name:"Object"`

	// 当 Type=table 时，ColName 为 \* 表示对表授权，如果为具体字段名，表示对字段授权
	ColName *string `json:"ColName" name:"ColName"`
}

func (r *GrantAccountPrivilegesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GrantAccountPrivilegesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GrantAccountPrivilegesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GrantAccountPrivilegesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GrantAccountPrivilegesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InitDBInstancesRequest struct {
	*tchttp.BaseRequest

	// 待初始化的实例Id列表，形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

	// 参数列表。本接口的可选值为：character_set_server（字符集，必传），lower_case_table_names（表名大小写敏感，必传），innodb_page_size（innodb数据页，默认16K），sync_mode（同步模式：0 - 异步； 1 - 强同步；2 - 强同步可退化。默认为强同步）。
	Params []*DBParamValue `json:"Params" name:"Params" list`
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

		// 异步任务Id，可通过 DescribeFlow 查询任务状态。
		FlowId *int64 `json:"FlowId" name:"FlowId"`

		// 透传入参。
		InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

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

type InstanceSpec struct {

	// 设备型号
	Machine *string `json:"Machine" name:"Machine"`

	// 该机型对应的可售卖规格列表
	SpecInfos []*SpecConfigInfo `json:"SpecInfos" name:"SpecInfos" list`
}

type LogFileInfo struct {

	// Log最后修改时间
	Mtime *uint64 `json:"Mtime" name:"Mtime"`

	// 文件长度
	Length *uint64 `json:"Length" name:"Length"`

	// 下载Log时用到的统一资源标识符
	Uri *string `json:"Uri" name:"Uri"`
}

type ModifyAccountDescriptionRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 登录用户名。
	UserName *string `json:"UserName" name:"UserName"`

	// 用户允许的访问 host，用户名+host唯一确定一个账号。
	Host *string `json:"Host" name:"Host"`

	// 新的账号备注，长度 0~256。
	Description *string `json:"Description" name:"Description"`
}

func (r *ModifyAccountDescriptionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAccountDescriptionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAccountDescriptionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyAccountDescriptionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAccountDescriptionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyBackupTimeRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 每天备份执行的区间的开始时间，格式 mm:ss，形如 22:00
	StartBackupTime *string `json:"StartBackupTime" name:"StartBackupTime"`

	// 每天备份执行的区间的结束时间，格式 mm:ss，形如 23:59
	EndBackupTime *string `json:"EndBackupTime" name:"EndBackupTime"`
}

func (r *ModifyBackupTimeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyBackupTimeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyBackupTimeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设置的状态，0 表示成功
		Status *int64 `json:"Status" name:"Status"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyBackupTimeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyBackupTimeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBInstanceNameRequest struct {
	*tchttp.BaseRequest

	// 待修改的实例 ID。形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 新的实例名称。允许的字符为字母、数字、下划线、连字符和中文。
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

	// 待修改的实例ID列表。实例 ID 形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

	// 要分配的项目 ID，可以通过 DescribeProjects 查询项目列表接口获取。
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`
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

type ModifyDBParametersRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 参数列表，每一个元素是Param和Value的组合
	Params []*DBParamValue `json:"Params" name:"Params" list`
}

func (r *ModifyDBParametersRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBParametersRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBParametersResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例 ID，形如：tdsql-ow728lmc。
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 参数修改结果
		Result []*ParamModifyResult `json:"Result" name:"Result" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDBParametersResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBParametersResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyLogFileRetentionPeriodRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 保存的天数,不能超过30
	Days *uint64 `json:"Days" name:"Days"`
}

func (r *ModifyLogFileRetentionPeriodRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyLogFileRetentionPeriodRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyLogFileRetentionPeriodResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例 ID，形如：tdsql-ow728lmc。
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyLogFileRetentionPeriodResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyLogFileRetentionPeriodResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type MonitorData struct {

	// 起始时间，形如 2018-03-24 23:59:59
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 结束时间，形如 2018-03-24 23:59:59
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 监控数据
	Data []*float64 `json:"Data" name:"Data" list`
}

type MonitorIntData struct {

	// 起始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 监控数据
	Data *int64 `json:"Data" name:"Data"`
}

type OpenDBExtranetAccessRequest struct {
	*tchttp.BaseRequest

	// 待开放外网访问的实例ID。形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
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

		// 异步任务Id，可通过 DescribeFlow 查询任务状态。
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

type ParamConstraint struct {

	// 约束类型,如枚举enum，区间section
	Type *string `json:"Type" name:"Type"`

	// 约束类型为enum时的可选值列表
	Enum *string `json:"Enum" name:"Enum"`

	// 约束类型为section时的范围
	Range *ConstraintRange `json:"Range" name:"Range"`

	// 约束类型为string时的可选值列表
	String *string `json:"String" name:"String"`
}

type ParamDesc struct {

	// 参数名字
	Param *string `json:"Param" name:"Param"`

	// 当前参数值
	Value *string `json:"Value" name:"Value"`

	// 设置过的值，参数生效后，该值和value一样。未设置过就不返回该字段。
	SetValue *string `json:"SetValue" name:"SetValue"`

	// 系统默认值
	Default *string `json:"Default" name:"Default"`

	// 参数限制
	Constraint *ParamConstraint `json:"Constraint" name:"Constraint"`
}

type ParamModifyResult struct {

	// 修改参数名字
	Param *string `json:"Param" name:"Param"`

	// 参数修改结果。0表示修改成功；-1表示修改失败；-2表示该参数值非法
	Code *int64 `json:"Code" name:"Code"`
}

type PerformanceMonitorSet struct {

	// 更新操作数UPDATE
	UpdateTotal *MonitorData `json:"UpdateTotal" name:"UpdateTotal"`

	// 磁盘每秒IO次数
	DiskIops *MonitorData `json:"DiskIops" name:"DiskIops"`

	// 活跃连接数
	ConnActive *MonitorData `json:"ConnActive" name:"ConnActive"`

	// 缓存命中率
	MemHitRate *MonitorData `json:"MemHitRate" name:"MemHitRate"`

	// 主备延迟
	SlaveDelay *MonitorData `json:"SlaveDelay" name:"SlaveDelay"`

	// 查询操作数SELECT
	SelectTotal *MonitorData `json:"SelectTotal" name:"SelectTotal"`

	// 慢查询数
	LongQuery *MonitorData `json:"LongQuery" name:"LongQuery"`

	// 删除操作数DELETE
	DeleteTotal *MonitorData `json:"DeleteTotal" name:"DeleteTotal"`

	// 插入操作数INSERT
	InsertTotal *MonitorData `json:"InsertTotal" name:"InsertTotal"`

	// 是否发生主备切换，1为发生，0否
	IsMasterSwitched *MonitorData `json:"IsMasterSwitched" name:"IsMasterSwitched"`
}

type RegionInfo struct {

	// 地域英文ID
	Region *string `json:"Region" name:"Region"`

	// 地域数字ID
	RegionId *int64 `json:"RegionId" name:"RegionId"`

	// 地域中文名
	RegionName *string `json:"RegionName" name:"RegionName"`

	// 可用区列表
	ZoneList []*ZonesInfo `json:"ZoneList" name:"ZoneList" list`

	// 可选择的主可用区和从可用区
	AvailableChoice []*ZoneChooseInfo `json:"AvailableChoice" name:"AvailableChoice" list`
}

type RenewDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 待续费的实例ID。形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 续费时长，单位：月。
	Period *int64 `json:"Period" name:"Period"`

	// 是否自动使用代金券进行支付，默认不使用。
	AutoVoucher *bool `json:"AutoVoucher" name:"AutoVoucher"`

	// 代金券ID列表，目前仅支持指定一张代金券。
	VoucherIds []*string `json:"VoucherIds" name:"VoucherIds" list`
}

func (r *RenewDBInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RenewDBInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RenewDBInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 长订单号。可以据此调用 DescribeOrders
	//  查询订单详细信息，或在支付失败时调用用户账号相关接口进行支付。
		DealName *string `json:"DealName" name:"DealName"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RenewDBInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RenewDBInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ResetAccountPasswordRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 登录用户名。
	UserName *string `json:"UserName" name:"UserName"`

	// 用户允许的访问 host，用户名+host唯一确定一个账号。
	Host *string `json:"Host" name:"Host"`

	// 新密码，由字母、数字或常见符号组成，不能包含分号、单引号和双引号，长度为6~32位。
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

type ResourceUsageMonitorSet struct {

	// binlog日志磁盘可用空间,单位GB
	BinlogDiskAvailable *MonitorData `json:"BinlogDiskAvailable" name:"BinlogDiskAvailable"`

	// CPU利用率
	CpuUsageRate *MonitorData `json:"CpuUsageRate" name:"CpuUsageRate"`

	// 内存可用空间,单位GB
	MemAvailable *MonitorData `json:"MemAvailable" name:"MemAvailable"`

	// 磁盘可用空间,单位GB
	DataDiskAvailable *MonitorIntData `json:"DataDiskAvailable" name:"DataDiskAvailable"`
}

type SlowLogData struct {

	// 语句校验和，用于查询详情
	CheckSum *string `json:"CheckSum" name:"CheckSum"`

	// 数据库名称
	Db *string `json:"Db" name:"Db"`

	// 抽象的SQL语句
	FingerPrint *string `json:"FingerPrint" name:"FingerPrint"`

	// 平均的锁时间
	LockTimeAvg *float64 `json:"LockTimeAvg" name:"LockTimeAvg"`

	// 最大锁时间
	LockTimeMax *float64 `json:"LockTimeMax" name:"LockTimeMax"`

	// 最小锁时间
	LockTimeMin *float64 `json:"LockTimeMin" name:"LockTimeMin"`

	// 锁时间总和
	LockTimeSum *float64 `json:"LockTimeSum" name:"LockTimeSum"`

	// 查询次数
	QueryCount *int64 `json:"QueryCount" name:"QueryCount"`

	// 平均查询时间
	QueryTimeAvg *float64 `json:"QueryTimeAvg" name:"QueryTimeAvg"`

	// 最大查询时间
	QueryTimeMax *float64 `json:"QueryTimeMax" name:"QueryTimeMax"`

	// 最小查询时间
	QueryTimeMin *float64 `json:"QueryTimeMin" name:"QueryTimeMin"`

	// 查询时间总和
	QueryTimeSum *float64 `json:"QueryTimeSum" name:"QueryTimeSum"`

	// 扫描行数
	RowsExaminedSum *uint64 `json:"RowsExaminedSum" name:"RowsExaminedSum"`

	// 发送行数
	RowsSentSum *uint64 `json:"RowsSentSum" name:"RowsSentSum"`

	// 首次执行时间
	TsMax *string `json:"TsMax" name:"TsMax"`

	// 最后执行时间
	TsMin *string `json:"TsMin" name:"TsMin"`

	// 帐号
	User *string `json:"User" name:"User"`
}

type SpecConfigInfo struct {

	// 设备型号
	Machine *string `json:"Machine" name:"Machine"`

	// 内存大小，单位 GB
	Memory *int64 `json:"Memory" name:"Memory"`

	// 数据盘规格最小值，单位 GB
	MinStorage *int64 `json:"MinStorage" name:"MinStorage"`

	// 数据盘规格最大值，单位 GB
	MaxStorage *int64 `json:"MaxStorage" name:"MaxStorage"`

	// 推荐的使用场景
	SuitInfo *string `json:"SuitInfo" name:"SuitInfo"`

	// 最大 Qps 值
	Qps *int64 `json:"Qps" name:"Qps"`

	// 产品类型 Id
	Pid *int64 `json:"Pid" name:"Pid"`

	// 节点个数，2 表示一主一从，3 表示一主二从
	NodeCount *int64 `json:"NodeCount" name:"NodeCount"`
}

type SqlLogItem struct {

	// 本条日志在消息队列中的偏移量。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 执行本条sql的用户。
	User *string `json:"User" name:"User"`

	// 执行本条sql的客户端IP+端口。
	Client *string `json:"Client" name:"Client"`

	// 数据库名称。
	DbName *string `json:"DbName" name:"DbName"`

	// 执行的sql语句。
	Sql *string `json:"Sql" name:"Sql"`

	// 返回的数据行数。
	SelectRowNum *uint64 `json:"SelectRowNum" name:"SelectRowNum"`

	// 影响行数。
	AffectRowNum *uint64 `json:"AffectRowNum" name:"AffectRowNum"`

	// Sql执行时间戳。
	Timestamp *uint64 `json:"Timestamp" name:"Timestamp"`

	// Sql耗时，单位为毫秒。
	TimeCostMs *uint64 `json:"TimeCostMs" name:"TimeCostMs"`

	// Sql返回码，0为成功。
	ResultCode *uint64 `json:"ResultCode" name:"ResultCode"`
}

type UpgradeDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 待升级的实例ID。形如：tdsql-ow728lmc，可以通过 DescribeDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 内存大小，单位：GB，可以通过 DescribeDBInstanceSpecs
	//  查询实例规格获得。
	Memory *int64 `json:"Memory" name:"Memory"`

	// 存储空间大小，单位：GB，可以通过 DescribeDBInstanceSpecs
	//  查询实例规格获得不同内存大小对应的磁盘规格下限和上限。
	Storage *int64 `json:"Storage" name:"Storage"`

	// 是否自动使用代金券进行支付，默认不使用。
	AutoVoucher *bool `json:"AutoVoucher" name:"AutoVoucher"`

	// 代金券ID列表，目前仅支持指定一张代金券。
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

		// 长订单号。可以据此调用 DescribeOrders
	//  查询订单详细信息，或在支付失败时调用用户账号相关接口进行支付。
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

type ZoneChooseInfo struct {

	// 主可用区
	MasterZone *ZonesInfo `json:"MasterZone" name:"MasterZone"`

	// 可选的从可用区
	SlaveZones []*ZonesInfo `json:"SlaveZones" name:"SlaveZones" list`
}

type ZonesInfo struct {

	// 可用区英文ID
	Zone *string `json:"Zone" name:"Zone"`

	// 可用区数字ID
	ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

	// 可用区中文名
	ZoneName *string `json:"ZoneName" name:"ZoneName"`
}
