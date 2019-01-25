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

package v20180416

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type CreateInstanceRequest struct {
	*tchttp.BaseRequest

	// 可用区
	Zone *string `json:"Zone" name:"Zone"`

	// 节点数量
	NodeNum *uint64 `json:"NodeNum" name:"NodeNum"`

	// 实例版本,当前只支持5.6.4
	EsVersion *string `json:"EsVersion" name:"EsVersion"`

	// 节点规格： 
	// ES.S1.SMALL2: 1核2G
	// ES.S1.MEDIUM4: 2核4G
	// ES.S1.MEDIUM8: 2核8G
	// ES.S1.LARGE16: 4核16G
	// ES.S1.2XLARGE32: 8核32G
	// ES.S1.4XLARGE64: 16核64G
	NodeType *string `json:"NodeType" name:"NodeType"`

	// 节点存储容量，单位GB
	DiskSize *uint64 `json:"DiskSize" name:"DiskSize"`

	// 私有网络ID
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 子网ID
	SubnetId *string `json:"SubnetId" name:"SubnetId"`

	// 访问密码，密码需8到16位，至少包括两项（[a-z,A-Z],[0-9]和[()`~!@#$%^&*-+=_|{}:;' <>,.?/]的特殊符号
	Password *string `json:"Password" name:"Password"`

	// 实例名称，1-50 个英文、汉字、数字、连接线-或下划线_
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// 计费类型: 
	// PREPAID：预付费，即包年包月 
	// POSTPAID_BY_HOUR：按小时后付费，默认值
	ChargeType *string `json:"ChargeType" name:"ChargeType"`

	// 包年包月购买时长，单位由TimeUint决定，默认为月
	ChargePeriod *uint64 `json:"ChargePeriod" name:"ChargePeriod"`

	// 自动续费标识，取值范围： 
	// RENEW_FLAG_AUTO：自动续费
	// RENEW_FLAG_MANUAL：不自动续费，用户手动续费
	// 如不传递该参数，普通用于默认不自动续费，SVIP用户自动续费
	RenewFlag *string `json:"RenewFlag" name:"RenewFlag"`

	// 节点存储类型,取值范围:  
	// LOCAL_BASIC: 本地硬盘  
	// LOCAL_SSD: 本地SSD硬盘，默认值  
	// CLOUD_BASIC: 普通云硬盘  
	// CLOUD_PREMIUM: 高硬能云硬盘  
	// CLOUD_SSD: SSD云硬盘
	DiskType *string `json:"DiskType" name:"DiskType"`

	// 计费时长单位，当前只支持“m”，表示月
	TimeUnit *string `json:"TimeUnit" name:"TimeUnit"`

	// 是否自动使用代金券，1是，0否，默认不使用
	AutoVoucher *int64 `json:"AutoVoucher" name:"AutoVoucher"`

	// 代金券ID列表，目前仅支持指定一张代金券
	VoucherIds []*string `json:"VoucherIds" name:"VoucherIds" list`
}

func (r *CreateInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例ID
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteInstanceRequest struct {
	*tchttp.BaseRequest

	// 要销毁的实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *DeleteInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeInstancesRequest struct {
	*tchttp.BaseRequest

	// 集群实例所属可用区，不传则默认所有可用区
	Zone *string `json:"Zone" name:"Zone"`

	// 一个或多个集群实例ID
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

	// 一个或多个集群实例名称
	InstanceNames []*string `json:"InstanceNames" name:"InstanceNames" list`

	// 分页起始值, 默认值0
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 分页大小，默认值20
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 排序字段：1，实例ID；2，实例名称；3，可用区；4，创建时间，若orderKey未传递则按创建时间降序排序
	OrderByKey *uint64 `json:"OrderByKey" name:"OrderByKey"`

	// 排序方式：0，升序；1，降序；若传递了orderByKey未传递orderByType, 则默认升序
	OrderByType *uint64 `json:"OrderByType" name:"OrderByType"`
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

		// 返回的实例个数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 实例详细信息列表
		InstanceList []*InstanceInfo `json:"InstanceList" name:"InstanceList" list`

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

type DictInfo struct {

	// 词典键值
	Key *string `json:"Key" name:"Key"`

	// 词典名称
	Name *string `json:"Name" name:"Name"`

	// 词典大小，单位B
	Size *uint64 `json:"Size" name:"Size"`
}

type EsAcl struct {

	// kibana访问黑名单
	BlackIpList []*string `json:"BlackIpList" name:"BlackIpList" list`

	// kibana访问白名单
	WhiteIpList []*string `json:"WhiteIpList" name:"WhiteIpList" list`
}

type EsDictionaryInfo struct {

	// 启用词词典列表
	MainDict []*DictInfo `json:"MainDict" name:"MainDict" list`

	// 停用词词典列表
	Stopwords []*DictInfo `json:"Stopwords" name:"Stopwords" list`
}

type InstanceInfo struct {

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// 地域
	Region *string `json:"Region" name:"Region"`

	// 可用区
	Zone *string `json:"Zone" name:"Zone"`

	// 用户ID
	AppId *uint64 `json:"AppId" name:"AppId"`

	// 用户UIN
	Uin *string `json:"Uin" name:"Uin"`

	// 实例所属VPC的UID
	VpcUid *string `json:"VpcUid" name:"VpcUid"`

	// 实例所属子网的UID
	SubnetUid *string `json:"SubnetUid" name:"SubnetUid"`

	// 实例状态，0:处理中,1:正常,-1停止,-2:销毁中,-3:已销毁
	Status *int64 `json:"Status" name:"Status"`

	// 实例计费模式。取值范围：  PREPAID：表示预付费，即包年包月  POSTPAID_BY_HOUR：表示后付费，即按量计费  CDHPAID：CDH付费，即只对CDH计费，不对CDH上的实例计费。
	ChargeType *string `json:"ChargeType" name:"ChargeType"`

	// 包年包月购买时长,单位:月
	ChargePeriod *uint64 `json:"ChargePeriod" name:"ChargePeriod"`

	// 自动续费标识。取值范围：  NOTIFY_AND_AUTO_RENEW：通知过期且自动续费  NOTIFY_AND_MANUAL_RENEW：通知过期不自动续费  DISABLE_NOTIFY_AND_MANUAL_RENEW：不通知过期不自动续费  默认取值：NOTIFY_AND_AUTO_RENEW。若该参数指定为NOTIFY_AND_AUTO_RENEW，在账户余额充足的情况下，实例到期后将按月自动续费。
	RenewFlag *string `json:"RenewFlag" name:"RenewFlag"`

	// 节点规格:  ES.S1.SMALL2 : 1核2G  ES.S1.MEDIUM4 : 2核4G  ES.S1.MEDIUM8 : 2核8G  ES.S1.LARGE16 : 4核16G  ES.S1.2XLARGE32 : 8核32G  ES.S1.3XLARGE32 : 12核32G  ES.S1.6XLARGE32 : 24核32G
	NodeType *string `json:"NodeType" name:"NodeType"`

	// 节点个数
	NodeNum *uint64 `json:"NodeNum" name:"NodeNum"`

	// 节点CPU核数
	CpuNum *uint64 `json:"CpuNum" name:"CpuNum"`

	// 节点内存大小，单位GB
	MemSize *uint64 `json:"MemSize" name:"MemSize"`

	// 节点磁盘类型
	DiskType *string `json:"DiskType" name:"DiskType"`

	// 节点磁盘大小，单位GB
	DiskSize *uint64 `json:"DiskSize" name:"DiskSize"`

	// ES域名
	EsDomain *string `json:"EsDomain" name:"EsDomain"`

	// ES VIP
	EsVip *string `json:"EsVip" name:"EsVip"`

	// ES端口
	EsPort *uint64 `json:"EsPort" name:"EsPort"`

	// Kibana访问url
	KibanaUrl *string `json:"KibanaUrl" name:"KibanaUrl"`

	// ES版本号
	EsVersion *string `json:"EsVersion" name:"EsVersion"`

	// ES配置项
	EsConfig *string `json:"EsConfig" name:"EsConfig"`

	// ES访问控制配置
	EsAcl *EsAcl `json:"EsAcl" name:"EsAcl"`

	// 实例创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 实例最后修改操作时间
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`

	// 实例到期时间
	Deadline *string `json:"Deadline" name:"Deadline"`

	// 实例类型（实例类型标识，当前只有1,2两种）
	InstanceType *uint64 `json:"InstanceType" name:"InstanceType"`

	// Ik分词器配置
	IkConfig *EsDictionaryInfo `json:"IkConfig" name:"IkConfig"`

	// 专用主节点配置
	MasterNodeInfo *MasterNodeInfo `json:"MasterNodeInfo" name:"MasterNodeInfo"`
}

type MasterNodeInfo struct {

	// 是否启用了专用主节点
	EnableDedicatedMaster *bool `json:"EnableDedicatedMaster" name:"EnableDedicatedMaster"`

	// 专用主节点规格
	MasterNodeType *string `json:"MasterNodeType" name:"MasterNodeType"`

	// 专用主节点个数
	MasterNodeNum *uint64 `json:"MasterNodeNum" name:"MasterNodeNum"`

	// 专用主节点CPU核数
	MasterNodeCpuNum *uint64 `json:"MasterNodeCpuNum" name:"MasterNodeCpuNum"`

	// 专用主节点内存大小，单位GB
	MasterNodeMemSize *uint64 `json:"MasterNodeMemSize" name:"MasterNodeMemSize"`

	// 专用主节点磁盘大小，单位GB
	MasterNodeDiskSize *uint64 `json:"MasterNodeDiskSize" name:"MasterNodeDiskSize"`
}

type RestartInstanceRequest struct {
	*tchttp.BaseRequest

	// 要重启的实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *RestartInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RestartInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RestartInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RestartInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RestartInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateInstanceRequest struct {
	*tchttp.BaseRequest

	// 要操作的实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 修改后的实例名称, 1-50 个英文、汉字、数字、连接线-或下划线_
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// 横向扩缩容后的节点个数
	NodeNum *uint64 `json:"NodeNum" name:"NodeNum"`

	// 修改后的配置项, JSON格式字符串
	EsConfig *string `json:"EsConfig" name:"EsConfig"`

	// 重置后的Kibana密码, 8到16位，至少包括两项（[a-z,A-Z],[0-9]和[-!@#$%&^*+=_:;,.?]的特殊符号
	Password *string `json:"Password" name:"Password"`

	// 修改后的访问控制列表
	EsAcl *EsAcl `json:"EsAcl" name:"EsAcl"`

	// 磁盘大小,单位GB
	DiskSize *uint64 `json:"DiskSize" name:"DiskSize"`

	// 节点规格: 
	// ES.S1.SMALL2: 1 核 2G
	// ES.S1.MEDIUM4: 2 核 4G 
	// ES.S1.MEDIUM8: 2 核 8G 
	// ES.S1.LARGE16: 4 核 16G 
	// ES.S1.2XLARGE32: 8 核 32G 
	// ES.S1.4XLARGE64: 16 核 64G
	NodeType *string `json:"NodeType" name:"NodeType"`
}

func (r *UpdateInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UpdateInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}
