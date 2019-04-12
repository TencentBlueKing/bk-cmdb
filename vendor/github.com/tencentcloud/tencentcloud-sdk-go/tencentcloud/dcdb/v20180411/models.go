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

package v20180411

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type AddShardConfig struct {

	// 新增分片的数量
	ShardCount *int64 `json:"ShardCount" name:"ShardCount"`

	// 分片内存大小，单位 GB
	ShardMemory *int64 `json:"ShardMemory" name:"ShardMemory"`

	// 分片存储大小，单位 GB
	ShardStorage *int64 `json:"ShardStorage" name:"ShardStorage"`
}

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

		// 异步任务流程ID
		FlowId *int64 `json:"FlowId" name:"FlowId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

	// 待关闭外网访问的实例ID。形如：dcdbt-ow728lmc，可以通过 DescribeDCDBInstances 查询实例详情获得。
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

	// 实例 ID，形如：dcdbt-ow728lmc。
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

	// 实例 ID，形如：dcdbt-ow728lmc，可以通过 DescribeDCDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// AccountName
	UserName *string `json:"UserName" name:"UserName"`

	// 可以登录的主机，与mysql 账号的 host 格式一致，可以支持通配符，例如 %，10.%，10.20.%。
	Host *string `json:"Host" name:"Host"`

	// 账号密码，由字母、数字或常见符号组成，不能包含分号、单引号和双引号，长度为6~32位。
	Password *string `json:"Password" name:"Password"`

	// 是否创建为只读账号，0：否， 1：该账号的sql请求优先选择备机执行，备机不可用时选择主机执行，2：优先选择备机执行，备机不可用时操作失败，3：只从备机读取。
	ReadOnly *int64 `json:"ReadOnly" name:"ReadOnly"`

	// 账号备注，可以包含中文、英文字符、常见符号和数字，长度为0~256字符
	Description *string `json:"Description" name:"Description"`

	// 如果备机延迟超过本参数设置值，系统将认为备机发生故障
	// 建议该参数值大于10。当ReadOnly选择1、2时该参数生效。
	DelayThresh *int64 `json:"DelayThresh" name:"DelayThresh"`
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type CreateDCDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 分片节点可用区分布，最多可填两个可用区。当分片规格为一主两从时，其中两个节点在第一个可用区。
	Zones []*string `json:"Zones" name:"Zones" list`

	// 欲购买的时长，单位：月。
	Period *int64 `json:"Period" name:"Period"`

	// 分片内存大小，单位：GB，可以通过 DescribeShardSpec
	//  查询实例规格获得。
	ShardMemory *int64 `json:"ShardMemory" name:"ShardMemory"`

	// 分片存储空间大小，单位：GB，可以通过 DescribeShardSpec
	//  查询实例规格获得。
	ShardStorage *int64 `json:"ShardStorage" name:"ShardStorage"`

	// 单个分片节点个数，可以通过 DescribeShardSpec
	//  查询实例规格获得。
	ShardNodeCount *int64 `json:"ShardNodeCount" name:"ShardNodeCount"`

	// 实例分片个数，可选范围2-8，可以通过升级实例进行新增分片到最多64个分片。
	ShardCount *int64 `json:"ShardCount" name:"ShardCount"`

	// 欲购买实例的数量，目前只支持购买1个实例
	Count *int64 `json:"Count" name:"Count"`

	// 项目 ID，可以通过查看项目列表获取，不传则关联到默认项目
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 虚拟私有网络 ID，不传或传空表示创建为基础网络
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 虚拟私有网络子网 ID，VpcId不为空时必填
	SubnetId *string `json:"SubnetId" name:"SubnetId"`

	// 数据库引擎版本，当前可选：10.0.10，10.1.9，5.7.17
	DbVersionId *string `json:"DbVersionId" name:"DbVersionId"`

	// 是否自动使用代金券进行支付，默认不使用。
	AutoVoucher *bool `json:"AutoVoucher" name:"AutoVoucher"`

	// 代金券ID列表，目前仅支持指定一张代金券。
	VoucherIds []*string `json:"VoucherIds" name:"VoucherIds" list`
}

func (r *CreateDCDBInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDCDBInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDCDBInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 长订单号。可以据此调用 DescribeOrders
	//  查询订单详细信息，或在支付失败时调用用户账号相关接口进行支付。
		DealName *string `json:"DealName" name:"DealName"`

		// 订单对应的实例 ID 列表，如果此处没有返回实例 ID，可以通过订单查询接口获取。还可通过实例查询接口查询实例是否创建完成。
		InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateDCDBInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDCDBInstanceResponse) FromJsonString(s string) error {
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

	// 如果备机延迟超过本参数设置值，系统将认为备机发生故障
	// 建议该参数值大于10。当ReadOnly选择1、2时该参数生效。
	DelayThresh *int64 `json:"DelayThresh" name:"DelayThresh"`
}

type DBParamValue struct {

	// 参数名称
	Param *string `json:"Param" name:"Param"`

	// 参数值
	Value *string `json:"Value" name:"Value"`
}

type DCDBInstanceInfo struct {

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// APPID
	AppId *int64 `json:"AppId" name:"AppId"`

	// 项目ID
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 地域
	Region *string `json:"Region" name:"Region"`

	// 可用区
	Zone *string `json:"Zone" name:"Zone"`

	// VPC数字ID
	VpcId *int64 `json:"VpcId" name:"VpcId"`

	// Subnet数字ID
	SubnetId *int64 `json:"SubnetId" name:"SubnetId"`

	// 状态中文描述
	StatusDesc *string `json:"StatusDesc" name:"StatusDesc"`

	// 状态
	Status *int64 `json:"Status" name:"Status"`

	// 内网IP
	Vip *string `json:"Vip" name:"Vip"`

	// 内网端口
	Vport *int64 `json:"Vport" name:"Vport"`

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 自动续费标志
	AutoRenewFlag *int64 `json:"AutoRenewFlag" name:"AutoRenewFlag"`

	// 内存大小，单位 GB
	Memory *int64 `json:"Memory" name:"Memory"`

	// 存储大小，单位 GB
	Storage *int64 `json:"Storage" name:"Storage"`

	// 分片个数
	ShardCount *int64 `json:"ShardCount" name:"ShardCount"`

	// 到期时间
	PeriodEndTime *string `json:"PeriodEndTime" name:"PeriodEndTime"`

	// 隔离时间
	IsolatedTimestamp *string `json:"IsolatedTimestamp" name:"IsolatedTimestamp"`

	// UIN
	Uin *string `json:"Uin" name:"Uin"`

	// 分片详情
	ShardDetail []*ShardInfo `json:"ShardDetail" name:"ShardDetail" list`

	// 节点数，2 为一主一从， 3 为一主二从
	NodeCount *int64 `json:"NodeCount" name:"NodeCount"`

	// 临时实例标记，0 为非临时实例
	IsTmp *int64 `json:"IsTmp" name:"IsTmp"`

	// 独享集群Id，为空表示非独享集群实例
	ExclusterId *string `json:"ExclusterId" name:"ExclusterId"`

	// 字符串型的私有网络Id
	UniqueVpcId *string `json:"UniqueVpcId" name:"UniqueVpcId"`

	// 字符串型的私有网络子网Id
	UniqueSubnetId *string `json:"UniqueSubnetId" name:"UniqueSubnetId"`

	// 数字实例Id（过时字段，请勿依赖该值）
	Id *uint64 `json:"Id" name:"Id"`

	// 外网访问的域名，公网可解析
	WanDomain *string `json:"WanDomain" name:"WanDomain"`

	// 外网 IP 地址，公网可访问
	WanVip *string `json:"WanVip" name:"WanVip"`

	// 外网端口
	WanPort *int64 `json:"WanPort" name:"WanPort"`

	// 产品类型 Id（过时字段，请勿依赖该值）
	Pid *int64 `json:"Pid" name:"Pid"`

	// 实例最后更新时间，格式为 2006-01-02 15:04:05
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`

	// 数据库引擎
	DbEngine *string `json:"DbEngine" name:"DbEngine"`

	// 数据库引擎版本
	DbVersion *string `json:"DbVersion" name:"DbVersion"`

	// 付费模式
	Paymode *string `json:"Paymode" name:"Paymode"`

	// 实例处于异步任务状态时，表示异步任务流程ID
	Locker *int64 `json:"Locker" name:"Locker"`

	// 外网状态，0-未开通；1-已开通；2-关闭；3-开通中
	WanStatus *int64 `json:"WanStatus" name:"WanStatus"`
}

type DCDBShardInfo struct {

	// 所属实例Id
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 分片SQL透传Id，用于将sql透传到指定分片执行
	ShardSerialId *string `json:"ShardSerialId" name:"ShardSerialId"`

	// 全局唯一的分片Id
	ShardInstanceId *string `json:"ShardInstanceId" name:"ShardInstanceId"`

	// 状态：0 创建中，1 流程处理中， 2 运行中，3 分片未初始化
	Status *int64 `json:"Status" name:"Status"`

	// 状态中文描述
	StatusDesc *string `json:"StatusDesc" name:"StatusDesc"`

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 字符串格式的私有网络Id
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 字符串格式的私有网络子网Id
	SubnetId *string `json:"SubnetId" name:"SubnetId"`

	// 项目ID
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 地域
	Region *string `json:"Region" name:"Region"`

	// 可用区
	Zone *string `json:"Zone" name:"Zone"`

	// 内存大小，单位 GB
	Memory *int64 `json:"Memory" name:"Memory"`

	// 存储大小，单位 GB
	Storage *int64 `json:"Storage" name:"Storage"`

	// 到期时间
	PeriodEndTime *string `json:"PeriodEndTime" name:"PeriodEndTime"`

	// 节点数，2 为一主一从， 3 为一主二从
	NodeCount *int64 `json:"NodeCount" name:"NodeCount"`

	// 存储使用率，单位为 %
	StorageUsage *float64 `json:"StorageUsage" name:"StorageUsage"`

	// 内存使用率，单位为 %
	MemoryUsage *float64 `json:"MemoryUsage" name:"MemoryUsage"`

	// 数字分片Id（过时字段，请勿依赖该值）
	ShardId *int64 `json:"ShardId" name:"ShardId"`
}

type Database struct {

	// 数据库名称
	DbName *string `json:"DbName" name:"DbName"`
}

type DatabaseFunction struct {

	// 函数名称
	Func *string `json:"Func" name:"Func"`
}

type DatabaseProcedure struct {

	// 存储过程名称
	Proc *string `json:"Proc" name:"Proc"`
}

type DatabaseTable struct {

	// 表名
	Table *string `json:"Table" name:"Table"`
}

type DatabaseView struct {

	// 视图名称
	View *string `json:"View" name:"View"`
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

	// 实例ID，形如：dcdbt-ow728lmc，可以通过 DescribeDCDBInstances 查询实例详情获得。
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

	// 实例 ID，形如：dcdbt-ow7t8lmc。
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

	// 实例ID，形如：dcdbt-ow728lmc。
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

type DescribeDBLogFilesRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：dcdbt-ow7t8lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 分片 ID，形如：shard-7noic7tv
	ShardId *string `json:"ShardId" name:"ShardId"`

	// 请求日志类型，取值只能为1、2、3或者4。1-binlog，2-冷备，3-errlog，4-slowlog。
	Type *int64 `json:"Type" name:"Type"`
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

		// 实例 ID，形如：dcdbt-ow728lmc。
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 请求日志类型。1-binlog，2-冷备，3-errlog，4-slowlog。
		Type *uint64 `json:"Type" name:"Type"`

		// 请求日志总数
		Total *uint64 `json:"Total" name:"Total"`

		// 日志文件列表
		Files []*LogFileInfo `json:"Files" name:"Files" list`

		// 如果是VPC网络的实例，做用本前缀加上URI为下载地址
		VpcPrefix *string `json:"VpcPrefix" name:"VpcPrefix"`

		// 如果是普通网络的实例，做用本前缀加上URI为下载地址
		NormalPrefix *string `json:"NormalPrefix" name:"NormalPrefix"`

		// 分片 ID，形如：shard-7noic7tv
		ShardId *string `json:"ShardId" name:"ShardId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

	// 实例 ID，形如：dcdbt-ow7t8lmc。
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

		// 实例 ID，形如：dcdbt-ow7t8lmc。
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 请求DB的当前参数值
		Params []*ParamDesc `json:"Params" name:"Params" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type DescribeDBSyncModeRequest struct {
	*tchttp.BaseRequest

	// 待修改同步模式的实例ID。形如：dcdbt-ow728lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *DescribeDBSyncModeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBSyncModeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBSyncModeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 同步模式：0 异步，1 强同步， 2 强同步可退化
		SyncMode *int64 `json:"SyncMode" name:"SyncMode"`

		// 是否有修改流程在执行中：1 是， 0 否。
		IsModifying *int64 `json:"IsModifying" name:"IsModifying"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBSyncModeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBSyncModeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDCDBInstancesRequest struct {
	*tchttp.BaseRequest

	// 按照一个或者多个实例 ID 查询。实例 ID 形如：dcdbt-2t4cf98d
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

	// 返回数量，默认为 10，最大值为 100。
	Limit *int64 `json:"Limit" name:"Limit"`

	// 1非独享集群，2独享集群， 0全部
	ExclusterType *int64 `json:"ExclusterType" name:"ExclusterType"`

	// 标识是否使用ExclusterType字段, false不使用，true使用
	IsFilterExcluster *bool `json:"IsFilterExcluster" name:"IsFilterExcluster"`
}

func (r *DescribeDCDBInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDCDBInstancesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDCDBInstancesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合条件的实例数量
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 实例详细信息列表
		Instances []*DCDBInstanceInfo `json:"Instances" name:"Instances" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDCDBInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDCDBInstancesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDCDBPriceRequest struct {
	*tchttp.BaseRequest

	// 欲新购实例的可用区ID。
	Zone *string `json:"Zone" name:"Zone"`

	// 欲购买实例的数量，目前只支持购买1个实例
	Count *int64 `json:"Count" name:"Count"`

	// 欲购买的时长，单位：月。
	Period *int64 `json:"Period" name:"Period"`

	// 单个分片节点个数大小，可以通过 DescribeShardSpec
	//  查询实例规格获得。
	ShardNodeCount *int64 `json:"ShardNodeCount" name:"ShardNodeCount"`

	// 分片内存大小，单位：GB，可以通过 DescribeShardSpec
	//  查询实例规格获得。
	ShardMemory *int64 `json:"ShardMemory" name:"ShardMemory"`

	// 分片存储空间大小，单位：GB，可以通过 DescribeShardSpec
	//  查询实例规格获得。
	ShardStorage *int64 `json:"ShardStorage" name:"ShardStorage"`

	// 实例分片个数，可选范围2-8，可以通过升级实例进行新增分片到最多64个分片。
	ShardCount *int64 `json:"ShardCount" name:"ShardCount"`
}

func (r *DescribeDCDBPriceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDCDBPriceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDCDBPriceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 原价，单位：分
		OriginalPrice *int64 `json:"OriginalPrice" name:"OriginalPrice"`

		// 实际价格，单位：分。受折扣等影响，可能和原价不同。
		Price *int64 `json:"Price" name:"Price"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDCDBPriceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDCDBPriceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDCDBRenewalPriceRequest struct {
	*tchttp.BaseRequest

	// 待续费的实例ID。形如：dcdbt-ow728lmc，可以通过 DescribeDCDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 续费时长，单位：月。不传则默认为1个月。
	Period *int64 `json:"Period" name:"Period"`
}

func (r *DescribeDCDBRenewalPriceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDCDBRenewalPriceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDCDBRenewalPriceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 原价，单位：分
		OriginalPrice *int64 `json:"OriginalPrice" name:"OriginalPrice"`

		// 实际价格，单位：分。受折扣等影响，可能和原价不同。
		Price *int64 `json:"Price" name:"Price"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDCDBRenewalPriceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDCDBRenewalPriceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDCDBSaleInfoRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeDCDBSaleInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDCDBSaleInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDCDBSaleInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 可售卖地域信息列表
		RegionList []*RegionInfo `json:"RegionList" name:"RegionList" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDCDBSaleInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDCDBSaleInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDCDBShardsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如：dcdbt-ow728lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 分片Id列表。
	ShardInstanceIds []*string `json:"ShardInstanceIds" name:"ShardInstanceIds" list`

	// 偏移量，默认为 0
	Offset *int64 `json:"Offset" name:"Offset"`

	// 返回数量，默认为 20，最大值为 100。
	Limit *int64 `json:"Limit" name:"Limit"`

	// 排序字段， 目前仅支持 createtime
	OrderBy *string `json:"OrderBy" name:"OrderBy"`

	// 排序类型， desc 或者 asc
	OrderByType *string `json:"OrderByType" name:"OrderByType"`
}

func (r *DescribeDCDBShardsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDCDBShardsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDCDBShardsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合条件的分片数量
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 分片信息列表
		Shards []*DCDBShardInfo `json:"Shards" name:"Shards" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDCDBShardsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDCDBShardsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDCDBUpgradePriceRequest struct {
	*tchttp.BaseRequest

	// 待升级的实例ID。形如：dcdbt-ow728lmc，可以通过 DescribeDCDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 升级类型，取值范围: 
	// <li> ADD: 新增分片 </li> 
	//  <li> EXPAND: 升级实例中的已有分片 </li> 
	//  <li> SPLIT: 将已有分片中的数据切分到新增分片上</li>
	UpgradeType *string `json:"UpgradeType" name:"UpgradeType"`

	// 新增分片配置，当UpgradeType为ADD时生效。
	AddShardConfig *AddShardConfig `json:"AddShardConfig" name:"AddShardConfig"`

	// 扩容分片配置，当UpgradeType为EXPAND时生效。
	ExpandShardConfig *ExpandShardConfig `json:"ExpandShardConfig" name:"ExpandShardConfig"`

	// 切分分片配置，当UpgradeType为SPLIT时生效。
	SplitShardConfig *SplitShardConfig `json:"SplitShardConfig" name:"SplitShardConfig"`
}

func (r *DescribeDCDBUpgradePriceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDCDBUpgradePriceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDCDBUpgradePriceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 原价，单位：分
		OriginalPrice *int64 `json:"OriginalPrice" name:"OriginalPrice"`

		// 实际价格，单位：分。受折扣等影响，可能和原价不同。
		Price *int64 `json:"Price" name:"Price"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDCDBUpgradePriceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDCDBUpgradePriceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDatabaseObjectsRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：dcdbt-ow7t8lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 数据库名称，通过 DescribeDatabases 接口获取。
	DbName *string `json:"DbName" name:"DbName"`
}

func (r *DescribeDatabaseObjectsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDatabaseObjectsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDatabaseObjectsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 透传入参。
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 数据库名称。
		DbName *string `json:"DbName" name:"DbName"`

		// 表列表。
		Tables []*DatabaseTable `json:"Tables" name:"Tables" list`

		// 视图列表。
		Views []*DatabaseView `json:"Views" name:"Views" list`

		// 存储过程列表。
		Procs []*DatabaseProcedure `json:"Procs" name:"Procs" list`

		// 函数列表。
		Funcs []*DatabaseFunction `json:"Funcs" name:"Funcs" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDatabaseObjectsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDatabaseObjectsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDatabaseTableRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：dcdbt-ow7t8lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 数据库名称，通过 DescribeDatabases 接口获取。
	DbName *string `json:"DbName" name:"DbName"`

	// 表名称，通过 DescribeDatabaseObjects 接口获取。
	Table *string `json:"Table" name:"Table"`
}

func (r *DescribeDatabaseTableRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDatabaseTableRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDatabaseTableResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例名称。
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 数据库名称。
		DbName *string `json:"DbName" name:"DbName"`

		// 表名称。
		Table *string `json:"Table" name:"Table"`

		// 列信息。
		Cols []*TableColumn `json:"Cols" name:"Cols" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDatabaseTableResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDatabaseTableResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDatabasesRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：dcdbt-ow7t8lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *DescribeDatabasesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDatabasesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDatabasesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 该实例上的数据库列表。
		Databases []*Database `json:"Databases" name:"Databases" list`

		// 透传入参。
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDatabasesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDatabasesResponse) FromJsonString(s string) error {
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type DescribeShardSpecRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeShardSpecRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeShardSpecRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeShardSpecResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 按机型分类的可售卖规格列表
		SpecConfig []*SpecConfig `json:"SpecConfig" name:"SpecConfig" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeShardSpecResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeShardSpecResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSqlLogsRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：dcdbt-ow728lmc，可以通过 DescribeDCDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// SQL日志偏移。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 拉取数量（0-1000，为0时拉取总数信息）。
	Limit *uint64 `json:"Limit" name:"Limit"`
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type ExpandShardConfig struct {

	// 分片ID数组
	ShardInstanceIds []*string `json:"ShardInstanceIds" name:"ShardInstanceIds" list`

	// 分片内存大小，单位 GB
	ShardMemory *int64 `json:"ShardMemory" name:"ShardMemory"`

	// 分片存储大小，单位 GB
	ShardStorage *int64 `json:"ShardStorage" name:"ShardStorage"`
}

type GrantAccountPrivilegesRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：dcdbt-ow728lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 登录用户名。
	UserName *string `json:"UserName" name:"UserName"`

	// 用户允许的访问 host，用户名+host唯一确定一个账号。
	Host *string `json:"Host" name:"Host"`

	// 数据库名。如果为 \*，表示查询全局权限（即 \*.\*），此时忽略 Type 和 Object 参数
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type InitDCDBInstancesRequest struct {
	*tchttp.BaseRequest

	// 待初始化的实例Id列表，形如：dcdbt-ow728lmc，可以通过 DescribeDCDBInstances 查询实例详情获得。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

	// 参数列表。本接口的可选值为：character_set_server（字符集，必传），lower_case_table_names（表名大小写敏感，必传），innodb_page_size（innodb数据页，默认16K），sync_mode（同步模式：0 - 异步； 1 - 强同步；2 - 强同步可退化。默认为强同步）。
	Params []*DBParamValue `json:"Params" name:"Params" list`
}

func (r *InitDCDBInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InitDCDBInstancesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InitDCDBInstancesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务Id，可通过 DescribeFlow 查询任务状态。
		FlowIds []*uint64 `json:"FlowIds" name:"FlowIds" list`

		// 透传入参。
		InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *InitDCDBInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InitDCDBInstancesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
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

	// 实例 ID，形如：dcdbt-ow728lmc。
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type ModifyDBInstancesProjectRequest struct {
	*tchttp.BaseRequest

	// 待修改的实例ID列表。实例 ID 形如：dcdbt-ow728lmc。
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

	// 实例 ID，形如：dcdbt-ow728lmc。
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

		// 实例 ID，形如：dcdbt-ow728lmc。
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 各参数修改结果
		Result []*ParamModifyResult `json:"Result" name:"Result" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type ModifyDBSyncModeRequest struct {
	*tchttp.BaseRequest

	// 待修改同步模式的实例ID。形如：dcdbt-ow728lmc。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 同步模式：0 异步，1 强同步， 2 强同步可退化
	SyncMode *int64 `json:"SyncMode" name:"SyncMode"`
}

func (r *ModifyDBSyncModeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBSyncModeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBSyncModeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务Id，可通过 DescribeFlow 查询任务状态。
		FlowId *int64 `json:"FlowId" name:"FlowId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDBSyncModeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBSyncModeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type OpenDBExtranetAccessRequest struct {
	*tchttp.BaseRequest

	// 待开放外网访问的实例ID。形如：dcdbt-ow728lmc。
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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
	AvailableChoice []*ShardZoneChooseInfo `json:"AvailableChoice" name:"AvailableChoice" list`
}

type RenewDCDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 待续费的实例ID。形如：dcdbt-ow728lmc，可以通过 DescribeDCDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 续费时长，单位：月。
	Period *int64 `json:"Period" name:"Period"`

	// 是否自动使用代金券进行支付，默认不使用。
	AutoVoucher *bool `json:"AutoVoucher" name:"AutoVoucher"`

	// 代金券ID列表，目前仅支持指定一张代金券。
	VoucherIds []*string `json:"VoucherIds" name:"VoucherIds" list`
}

func (r *RenewDCDBInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RenewDCDBInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RenewDCDBInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 长订单号。可以据此调用 DescribeOrders
	//  查询订单详细信息，或在支付失败时调用用户账号相关接口进行支付。
		DealName *string `json:"DealName" name:"DealName"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RenewDCDBInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RenewDCDBInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ResetAccountPasswordRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，形如：dcdbt-ow728lmc。
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type ShardInfo struct {

	// 分片ID
	ShardInstanceId *string `json:"ShardInstanceId" name:"ShardInstanceId"`

	// 分片Set ID
	ShardSerialId *string `json:"ShardSerialId" name:"ShardSerialId"`

	// 状态：0 创建中，1 流程处理中， 2 运行中，3 分片未初始化，-2 分片已删除
	Status *int64 `json:"Status" name:"Status"`

	// 创建时间
	Createtime *string `json:"Createtime" name:"Createtime"`

	// 内存大小，单位 GB
	Memory *int64 `json:"Memory" name:"Memory"`

	// 存储大小，单位 GB
	Storage *int64 `json:"Storage" name:"Storage"`

	// 分片数字ID
	ShardId *int64 `json:"ShardId" name:"ShardId"`

	// 节点数，2 为一主一从， 3 为一主二从
	NodeCount *int64 `json:"NodeCount" name:"NodeCount"`

	// 产品类型 Id（过时字段，请勿依赖该值）
	Pid *int64 `json:"Pid" name:"Pid"`
}

type ShardZoneChooseInfo struct {

	// 主可用区
	MasterZone *ZonesInfo `json:"MasterZone" name:"MasterZone"`

	// 可选的从可用区
	SlaveZones []*ZonesInfo `json:"SlaveZones" name:"SlaveZones" list`
}

type SpecConfig struct {

	// 规格机型
	Machine *string `json:"Machine" name:"Machine"`

	// 规格列表
	SpecConfigInfos []*SpecConfigInfo `json:"SpecConfigInfos" name:"SpecConfigInfos" list`
}

type SpecConfigInfo struct {

	// 节点个数，2 表示一主一从，3 表示一主二从
	NodeCount *uint64 `json:"NodeCount" name:"NodeCount"`

	// 内存大小，单位 GB
	Memory *int64 `json:"Memory" name:"Memory"`

	// 数据盘规格最小值，单位 GB
	MinStorage *int64 `json:"MinStorage" name:"MinStorage"`

	// 数据盘规格最大值，单位 GB
	MaxStorage *int64 `json:"MaxStorage" name:"MaxStorage"`

	// 推荐的使用场景
	SuitInfo *string `json:"SuitInfo" name:"SuitInfo"`

	// 产品类型 Id
	Pid *int64 `json:"Pid" name:"Pid"`

	// 最大 Qps 值
	Qps *int64 `json:"Qps" name:"Qps"`
}

type SplitShardConfig struct {

	// 分片ID数组
	ShardInstanceIds []*string `json:"ShardInstanceIds" name:"ShardInstanceIds" list`

	// 数据切分比例
	SplitRate *int64 `json:"SplitRate" name:"SplitRate"`

	// 分片内存大小，单位 GB
	ShardMemory *int64 `json:"ShardMemory" name:"ShardMemory"`

	// 分片存储大小，单位 GB
	ShardStorage *int64 `json:"ShardStorage" name:"ShardStorage"`
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

type TableColumn struct {

	// 列名称
	Col *string `json:"Col" name:"Col"`

	// 列类型
	Type *string `json:"Type" name:"Type"`
}

type UpgradeDCDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 待升级的实例ID。形如：dcdbt-ow728lmc，可以通过 DescribeDCDBInstances 查询实例详情获得。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 升级类型，取值范围: 
	// <li> ADD: 新增分片 </li> 
	//  <li> EXPAND: 升级实例中的已有分片 </li> 
	//  <li> SPLIT: 将已有分片中的数据切分到新增分片上</li>
	UpgradeType *string `json:"UpgradeType" name:"UpgradeType"`

	// 新增分片配置，当UpgradeType为ADD时生效。
	AddShardConfig *AddShardConfig `json:"AddShardConfig" name:"AddShardConfig"`

	// 扩容分片配置，当UpgradeType为EXPAND时生效。
	ExpandShardConfig *ExpandShardConfig `json:"ExpandShardConfig" name:"ExpandShardConfig"`

	// 切分分片配置，当UpgradeType为SPLIT时生效。
	SplitShardConfig *SplitShardConfig `json:"SplitShardConfig" name:"SplitShardConfig"`

	// 是否自动使用代金券进行支付，默认不使用。
	AutoVoucher *bool `json:"AutoVoucher" name:"AutoVoucher"`

	// 代金券ID列表，目前仅支持指定一张代金券。
	VoucherIds []*string `json:"VoucherIds" name:"VoucherIds" list`
}

func (r *UpgradeDCDBInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpgradeDCDBInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpgradeDCDBInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 长订单号。可以据此调用 DescribeOrders
	//  查询订单详细信息，或在支付失败时调用用户账号相关接口进行支付。
		DealName *string `json:"DealName" name:"DealName"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UpgradeDCDBInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpgradeDCDBInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ZonesInfo struct {

	// 可用区英文ID
	Zone *string `json:"Zone" name:"Zone"`

	// 可用区数字ID
	ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

	// 可用区中文名
	ZoneName *string `json:"ZoneName" name:"ZoneName"`
}
