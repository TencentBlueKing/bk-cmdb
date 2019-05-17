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

package v20180412

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type ClearInstanceRequest struct {
	*tchttp.BaseRequest

	// 实例Id
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// redis的实例密码
	Password *string `json:"Password" name:"Password"`
}

func (r *ClearInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ClearInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ClearInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务Id
		TaskId *int64 `json:"TaskId" name:"TaskId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ClearInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ClearInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateInstancesRequest struct {
	*tchttp.BaseRequest

	// 实例所属的可用区id
	ZoneId *uint64 `json:"ZoneId" name:"ZoneId"`

	// 实例类型：2 – Redis2.8主从版，3 – Redis3.2主从版(CKV主从版)，4 – Redis3.2集群版(CKV集群版)，5-Redis2.8单机版，7 – Redis4.0集群版，
	TypeId *uint64 `json:"TypeId" name:"TypeId"`

	// 实例容量，单位MB， 取值大小以 查询售卖规格接口返回的规格为准
	MemSize *uint64 `json:"MemSize" name:"MemSize"`

	// 实例数量，单次购买实例数量以 查询售卖规格接口返回的规格为准
	GoodsNum *uint64 `json:"GoodsNum" name:"GoodsNum"`

	// 购买时长，单位：月，取值范围 [1,2,3,4,5,6,7,8,9,10,11,12,24,36]
	Period *uint64 `json:"Period" name:"Period"`

	// 实例密码，密码规则：1.长度为8-16个字符；2:至少包含字母、数字和字符!@^*()中的两种
	Password *string `json:"Password" name:"Password"`

	// 付费方式:0-按量计费，1-包年包月。
	BillingMode *int64 `json:"BillingMode" name:"BillingMode"`

	// 私有网络ID，如果不传则默认选择基础网络，请使用私有网络列表查询，如：vpc-sad23jfdfk
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 基础网络下， subnetId无效； vpc子网下，取值以查询子网列表，如：subnet-fdj24n34j2
	SubnetId *string `json:"SubnetId" name:"SubnetId"`

	// 项目id，取值以用户账户>用户账户相关接口查询>项目列表返回的projectId为准
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 自动续费表示。0 - 默认状态（手动续费）；1 - 自动续费；2 - 明确不自动续费
	AutoRenew *uint64 `json:"AutoRenew" name:"AutoRenew"`

	// 安全组id数组
	SecurityGroupIdList []*string `json:"SecurityGroupIdList" name:"SecurityGroupIdList" list`

	// 用户自定义的端口 不填则默认为6379
	VPort *uint64 `json:"VPort" name:"VPort"`

	// 实例分片数量，Redis2.8主从版、CKV主从版和Redis2.8单机版不需要填写
	RedisShardNum *int64 `json:"RedisShardNum" name:"RedisShardNum"`

	// 实例副本数量，Redis2.8主从版、CKV主从版和Redis2.8单机版不需要填写
	RedisReplicasNum *int64 `json:"RedisReplicasNum" name:"RedisReplicasNum"`

	// 是否支持副本只读，Redis2.8主从版、CKV主从版和Redis2.8单机版不需要填写
	ReplicasReadonly *bool `json:"ReplicasReadonly" name:"ReplicasReadonly"`
}

func (r *CreateInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateInstancesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateInstancesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 交易的Id
		DealId *string `json:"DealId" name:"DealId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateInstancesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAutoBackupConfigRequest struct {
	*tchttp.BaseRequest

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *DescribeAutoBackupConfigRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAutoBackupConfigRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAutoBackupConfigResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 备份类型。自动备份类型： 1 “定时回档”
		AutoBackupType *int64 `json:"AutoBackupType" name:"AutoBackupType"`

		// Monday，Tuesday，Wednesday，Thursday，Friday，Saturday，Sunday。
		WeekDays []*string `json:"WeekDays" name:"WeekDays" list`

		// 时间段。
		TimePeriod *string `json:"TimePeriod" name:"TimePeriod"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeAutoBackupConfigResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAutoBackupConfigResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeInstanceBackupsRequest struct {
	*tchttp.BaseRequest

	// 待操作的实例ID，可通过 DescribeInstance 接口返回值中的 InstanceId 获取。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 实例列表大小，默认大小20
	Limit *int64 `json:"Limit" name:"Limit"`

	// 偏移量，取Limit整数倍
	Offset *int64 `json:"Offset" name:"Offset"`

	// 开始时间，格式如：2017-02-08 16:46:34。查询实例在 [beginTime, endTime] 时间段内开始备份的备份列表。
	BeginTime *string `json:"BeginTime" name:"BeginTime"`

	// 结束时间，格式如：2017-02-08 19:09:26。查询实例在 [beginTime, endTime] 时间段内开始备份的备份列表。
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 1：备份在流程中，2：备份正常，3：备份转RDB文件处理中，4：已完成RDB转换，-1：备份已过期，-2：备份已删除。
	Status []*string `json:"Status" name:"Status" list`
}

func (r *DescribeInstanceBackupsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeInstanceBackupsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeInstanceBackupsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 备份总数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 实例的备份数组
		BackupSet []*RedisBackupSet `json:"BackupSet" name:"BackupSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeInstanceBackupsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeInstanceBackupsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeInstanceDealDetailRequest struct {
	*tchttp.BaseRequest

	// 订单ID数组
	DealIds []*string `json:"DealIds" name:"DealIds" list`
}

func (r *DescribeInstanceDealDetailRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeInstanceDealDetailRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeInstanceDealDetailResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 订单详细信息
		DealDetails []*TradeDealDetail `json:"DealDetails" name:"DealDetails" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeInstanceDealDetailResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeInstanceDealDetailResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeInstancesRequest struct {
	*tchttp.BaseRequest

	// 实例列表的大小，参数默认值20
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，取Limit整数倍
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 实例Id，如：crs-6ubhgouj
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 枚举范围： projectId,createtime,instancename,type,curDeadline
	OrderBy *string `json:"OrderBy" name:"OrderBy"`

	// 1倒序，0顺序，默认倒序
	OrderType *int64 `json:"OrderType" name:"OrderType"`

	// 私有网络ID数组，数组下标从0开始，如果不传则默认选择基础网络，如：47525
	VpcIds []*string `json:"VpcIds" name:"VpcIds" list`

	// 子网ID数组，数组下标从0开始，如：56854
	SubnetIds []*string `json:"SubnetIds" name:"SubnetIds" list`

	// 项目ID 组成的数组，数组下标从0开始
	ProjectIds []*int64 `json:"ProjectIds" name:"ProjectIds" list`

	// 查找实例的ID。
	SearchKey *string `json:"SearchKey" name:"SearchKey"`

	// 实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// 私有网络ID数组，数组下标从0开始，如果不传则默认选择基础网络，如：vpc-sad23jfdfk
	UniqVpcIds []*string `json:"UniqVpcIds" name:"UniqVpcIds" list`

	// 子网ID数组，数组下标从0开始，如：subnet-fdj24n34j2
	UniqSubnetIds []*string `json:"UniqSubnetIds" name:"UniqSubnetIds" list`

	// 地域ID，已经弃用，可通过公共参数Region查询对应地域
	RegionIds []*int64 `json:"RegionIds" name:"RegionIds" list`
}

func (r *DescribeInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeInstancesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeInstancesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 实例详细信息列表
		InstanceSet []*InstanceSet `json:"InstanceSet" name:"InstanceSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeInstancesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProductInfoRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeProductInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProductInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProductInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 地域售卖信息
		RegionSet []*RegionConf `json:"RegionSet" name:"RegionSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeProductInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProductInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTaskInfoRequest struct {
	*tchttp.BaseRequest

	// 任务ID
	TaskId *uint64 `json:"TaskId" name:"TaskId"`
}

func (r *DescribeTaskInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTaskInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTaskInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务状态preparing:待执行，running：执行中，succeed：成功，failed：失败，error 执行出错
		Status *string `json:"Status" name:"Status"`

		// 任务开始时间
		StartTime *string `json:"StartTime" name:"StartTime"`

		// 任务类型
		TaskType *string `json:"TaskType" name:"TaskType"`

		// 实例的ID
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 任务信息，错误时显示错误信息。执行中与成功则为空
		TaskMessage *string `json:"TaskMessage" name:"TaskMessage"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeTaskInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTaskInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InstanceSet struct {

	// 实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// 实例Id
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 用户的Appid
	Appid *int64 `json:"Appid" name:"Appid"`

	// 项目Id
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 地域id 1--广州 4--上海 5-- 香港 6--多伦多 7--上海金融 8--北京 9-- 新加坡 11--深圳金融 15--美西（硅谷）16--成都 17--德国 18--韩国 19--重庆 21--印度 22--美东（弗吉尼亚）23--泰国 24--俄罗斯 25--日本
	RegionId *int64 `json:"RegionId" name:"RegionId"`

	// 区域id
	ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

	// vpc网络id 如：75101
	VpcId *int64 `json:"VpcId" name:"VpcId"`

	// vpc网络下子网id 如：46315
	SubnetId *int64 `json:"SubnetId" name:"SubnetId"`

	// 实例当前状态，0：待初始化；1：实例在流程中；2：实例运行中；-2：实例已隔离；-3：实例待删除
	Status *int64 `json:"Status" name:"Status"`

	// 实例vip
	WanIp *string `json:"WanIp" name:"WanIp"`

	// 实例端口号
	Port *int64 `json:"Port" name:"Port"`

	// 实例创建时间
	Createtime *string `json:"Createtime" name:"Createtime"`

	// 实例容量大小，单位：MB
	Size *float64 `json:"Size" name:"Size"`

	// 实例当前已使用容量，单位：MB
	SizeUsed *float64 `json:"SizeUsed" name:"SizeUsed"`

	// 实例类型，1：Redis2.8集群版；2：Redis2.8主从版；3：CKV主从版（Redis3.2）；4：CKV集群版（Redis3.2）；5：Redis2.8单机版；7：Redis4.0集群版；
	Type *int64 `json:"Type" name:"Type"`

	// 实例是否设置自动续费标识，1：设置自动续费；0：未设置自动续费
	AutoRenewFlag *int64 `json:"AutoRenewFlag" name:"AutoRenewFlag"`

	// 实例到期时间
	DeadlineTime *string `json:"DeadlineTime" name:"DeadlineTime"`

	// 引擎：社区版Redis、腾讯云CKV
	Engine *string `json:"Engine" name:"Engine"`

	// 产品类型：Redis2.8集群版、Redis2.8主从版、Redis3.2主从版（CKV主从版）、Redis3.2集群版（CKV集群版）、Redis2.8单机版、Redis4.0集群版
	ProductType *string `json:"ProductType" name:"ProductType"`

	// vpc网络id 如：vpc-fk33jsf43kgv
	UniqVpcId *string `json:"UniqVpcId" name:"UniqVpcId"`

	// vpc网络下子网id 如：subnet-fd3j6l35mm0
	UniqSubnetId *string `json:"UniqSubnetId" name:"UniqSubnetId"`

	// 计费模式：0-按量计费，1-包年包月
	BillingMode *int64 `json:"BillingMode" name:"BillingMode"`
}

type ManualBackupInstanceRequest struct {
	*tchttp.BaseRequest

	// 待操作的实例ID，可通过 DescribeInstance接口返回值中的 InstanceId 获取。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 备份的备注信息
	Remark *string `json:"Remark" name:"Remark"`
}

func (r *ManualBackupInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ManualBackupInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ManualBackupInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务ID
		TaskId *int64 `json:"TaskId" name:"TaskId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ManualBackupInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ManualBackupInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModfiyInstancePasswordRequest struct {
	*tchttp.BaseRequest

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 实例旧密码
	OldPassword *string `json:"OldPassword" name:"OldPassword"`

	// 实例新密码
	Password *string `json:"Password" name:"Password"`
}

func (r *ModfiyInstancePasswordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModfiyInstancePasswordRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModfiyInstancePasswordResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务ID
		TaskId *int64 `json:"TaskId" name:"TaskId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModfiyInstancePasswordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModfiyInstancePasswordResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAutoBackupConfigRequest struct {
	*tchttp.BaseRequest

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 日期 Monday，Tuesday，Wednesday，Thursday，Friday，Saturday，Sunday
	WeekDays []*string `json:"WeekDays" name:"WeekDays" list`

	// 时间段 00:00-01:00, 01:00-02:00...... 23:00-00:00
	TimePeriod *string `json:"TimePeriod" name:"TimePeriod"`

	// 自动备份类型： 1 “定时回档”
	AutoBackupType *int64 `json:"AutoBackupType" name:"AutoBackupType"`
}

func (r *ModifyAutoBackupConfigRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAutoBackupConfigRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAutoBackupConfigResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 自动备份类型： 1 “定时回档”
		AutoBackupType *int64 `json:"AutoBackupType" name:"AutoBackupType"`

		// 日期Monday，Tuesday，Wednesday，Thursday，Friday，Saturday，Sunday。
		WeekDays []*string `json:"WeekDays" name:"WeekDays" list`

		// 时间段 00:00-01:00, 01:00-02:00...... 23:00-00:00
		TimePeriod *string `json:"TimePeriod" name:"TimePeriod"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyAutoBackupConfigResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAutoBackupConfigResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ProductConf struct {

	// 产品类型，2-Redis主从版，3-CKV主从版，4-CKV集群版，5-Redis单机版，7-Redis集群版
	Type *int64 `json:"Type" name:"Type"`

	// 产品名称，Redis主从版，CKV主从版，CKV集群版，Redis单机版，Redis集群版
	TypeName *string `json:"TypeName" name:"TypeName"`

	// 购买时的最小数量
	MinBuyNum *int64 `json:"MinBuyNum" name:"MinBuyNum"`

	// 购买时的最大数量
	MaxBuyNum *int64 `json:"MaxBuyNum" name:"MaxBuyNum"`

	// 产品是否售罄
	Saleout *bool `json:"Saleout" name:"Saleout"`

	// 产品引擎，腾讯云CKV或者社区版Redis
	Engine *string `json:"Engine" name:"Engine"`

	// 兼容版本，Redis-2.8，Redis-3.2，Redis-4.0
	Version *string `json:"Version" name:"Version"`

	// 规格总大小，单位G
	TotalSize []*string `json:"TotalSize" name:"TotalSize" list`

	// 每个分片大小，单位G
	ShardSize []*string `json:"ShardSize" name:"ShardSize" list`

	// 副本数量
	ReplicaNum []*string `json:"ReplicaNum" name:"ReplicaNum" list`

	// 分片数量
	ShardNum []*string `json:"ShardNum" name:"ShardNum" list`

	// 支持的计费模式，1-包年包月，0-按量计费
	PayMode *string `json:"PayMode" name:"PayMode"`

	// 是否支持副本只读
	EnableRepicaReadOnly *bool `json:"EnableRepicaReadOnly" name:"EnableRepicaReadOnly"`
}

type RedisBackupSet struct {

	// 开始备份的时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 备份ID
	BackupId *string `json:"BackupId" name:"BackupId"`

	// 备份类型。 manualBackupInstance：用户发起的手动备份； systemBackupInstance：凌晨系统发起的备份
	BackupType *string `json:"BackupType" name:"BackupType"`

	// 备份状态。  1:"备份被其它流程锁定";  2:"备份正常，没有被任何流程锁定";  -1:"备份已过期"； 3:"备份正在被导出";  4:"备份导出成功"
	Status *int64 `json:"Status" name:"Status"`

	// 备份的备注信息
	Remark *string `json:"Remark" name:"Remark"`

	// 备份是否被锁定，0：未被锁定；1：已被锁定
	Locked *int64 `json:"Locked" name:"Locked"`
}

type RegionConf struct {

	// 地域ID
	RegionId *string `json:"RegionId" name:"RegionId"`

	// 地域名称
	RegionName *string `json:"RegionName" name:"RegionName"`

	// 地域简称
	RegionShortName *string `json:"RegionShortName" name:"RegionShortName"`

	// 地域所在大区名称
	Area *string `json:"Area" name:"Area"`

	// 可用区信息
	ZoneSet []*ZoneCapacityConf `json:"ZoneSet" name:"ZoneSet" list`
}

type RenewInstanceRequest struct {
	*tchttp.BaseRequest

	// 购买时长，单位：月
	Period *uint64 `json:"Period" name:"Period"`

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
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

		// 交易Id
		DealId *string `json:"DealId" name:"DealId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type ResetPasswordRequest struct {
	*tchttp.BaseRequest

	// 重置的密码
	Password *string `json:"Password" name:"Password"`

	// Redis实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *ResetPasswordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ResetPasswordRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ResetPasswordResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务ID
		TaskId *int64 `json:"TaskId" name:"TaskId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ResetPasswordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ResetPasswordResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type TradeDealDetail struct {

	// 订单号ID，调用云API时使用此ID
	DealId *string `json:"DealId" name:"DealId"`

	// 长订单ID，反馈订单问题给官方客服使用此ID
	DealName *string `json:"DealName" name:"DealName"`

	// 可用区id
	ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

	// 订单关联的实例数
	GoodsNum *int64 `json:"GoodsNum" name:"GoodsNum"`

	// 创建用户uin
	Creater *string `json:"Creater" name:"Creater"`

	// 订单创建时间
	CreatTime *string `json:"CreatTime" name:"CreatTime"`

	// 订单超时时间
	OverdueTime *string `json:"OverdueTime" name:"OverdueTime"`

	// 订单完成时间
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 订单状态 1：未支付 2:已支付，未发货 3:发货中 4:发货成功 5:发货失败 6:已退款 7:已关闭订单 8:订单过期 9:订单已失效 10:产品已失效 11:代付拒绝 12:支付中
	Status *int64 `json:"Status" name:"Status"`

	// 订单状态描述
	Description *string `json:"Description" name:"Description"`

	// 订单实际总价，单位：分
	Price *int64 `json:"Price" name:"Price"`

	// 实例ID
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`
}

type UpgradeInstanceRequest struct {
	*tchttp.BaseRequest

	// 实例Id
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 分片大小 单位 MB
	MemSize *uint64 `json:"MemSize" name:"MemSize"`

	// 分片数量，Redis2.8主从版、CKV主从版和Redis2.8单机版不需要填写
	RedisShardNum *uint64 `json:"RedisShardNum" name:"RedisShardNum"`

	// 副本数量，Redis2.8主从版、CKV主从版和Redis2.8单机版不需要填写
	RedisReplicasNum *uint64 `json:"RedisReplicasNum" name:"RedisReplicasNum"`
}

func (r *UpgradeInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpgradeInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpgradeInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 订单ID
		DealId *string `json:"DealId" name:"DealId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UpgradeInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpgradeInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ZoneCapacityConf struct {

	// 可用区ID：如ap-guangzhou-3
	ZoneId *string `json:"ZoneId" name:"ZoneId"`

	// 可用区名称
	ZoneName *string `json:"ZoneName" name:"ZoneName"`

	// 可用区是否售罄
	IsSaleout *bool `json:"IsSaleout" name:"IsSaleout"`

	// 是否为默认可用区
	IsDefault *bool `json:"IsDefault" name:"IsDefault"`

	// 网络类型：basenet -- 基础网络；vpcnet -- VPC网络
	NetWorkType []*string `json:"NetWorkType" name:"NetWorkType" list`

	// 可用区内产品规格等信息
	ProductSet []*ProductConf `json:"ProductSet" name:"ProductSet" list`

	// 可用区ID：如100003
	OldZoneId *int64 `json:"OldZoneId" name:"OldZoneId"`
}
