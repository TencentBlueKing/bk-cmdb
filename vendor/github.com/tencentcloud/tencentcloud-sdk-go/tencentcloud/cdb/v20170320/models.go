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

package v20170320

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type Account struct {

	// 新账户的名称
	User *string `json:"User" name:"User"`

	// 新账户的域名
	Host *string `json:"Host" name:"Host"`
}

type AccountInfo struct {

	// 账号备注信息
	Notes *string `json:"Notes" name:"Notes"`

	// 账号的域名
	Host *string `json:"Host" name:"Host"`

	// 账号的名称
	User *string `json:"User" name:"User"`

	// 账号信息修改时间
	ModifyTime *string `json:"ModifyTime" name:"ModifyTime"`

	// 修改密码的时间
	ModifyPasswordTime *string `json:"ModifyPasswordTime" name:"ModifyPasswordTime"`

	// 账号的创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`
}

type AssociateSecurityGroupsRequest struct {
	*tchttp.BaseRequest

	// 安全组Id。
	SecurityGroupId *string `json:"SecurityGroupId" name:"SecurityGroupId"`

	// 实例ID列表，一个或者多个实例Id组成的数组。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`
}

func (r *AssociateSecurityGroupsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AssociateSecurityGroupsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AssociateSecurityGroupsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AssociateSecurityGroupsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AssociateSecurityGroupsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type BackupConfig struct {

	// 第二个从库复制方式，可能的返回值：async-异步，semisync-半同步
	ReplicationMode *string `json:"ReplicationMode" name:"ReplicationMode"`

	// 第二个从库可用区的正式名称，如ap-shanghai-1
	Zone *string `json:"Zone" name:"Zone"`

	// 第二个从库内网IP地址
	Vip *string `json:"Vip" name:"Vip"`

	// 第二个从库访问端口
	Vport *string `json:"Vport" name:"Vport"`
}

type BackupInfo struct {

	// 备份文件名
	Name *string `json:"Name" name:"Name"`

	// 备份文件大小，单位：Byte
	Size *int64 `json:"Size" name:"Size"`

	// 备份快照时间，时间格式：2016-03-17 02:10:37
	Date *string `json:"Date" name:"Date"`

	// 内网下载地址
	IntranetUrl *string `json:"IntranetUrl" name:"IntranetUrl"`

	// 外网下载地址
	InternetUrl *string `json:"InternetUrl" name:"InternetUrl"`

	// 日志具体类型，可能的值有：logic - 逻辑冷备，physical - 物理冷备
	Type *string `json:"Type" name:"Type"`

	// 备份子任务的ID，删除备份文件时使用
	BackupId *int64 `json:"BackupId" name:"BackupId"`

	// 备份任务状态
	Status *string `json:"Status" name:"Status"`

	// 备份任务的完成时间
	FinishTime *string `json:"FinishTime" name:"FinishTime"`

	// 备份的创建者，可能的值：SYSTEM - 系统创建，Uin - 发起者Uin值
	Creator *string `json:"Creator" name:"Creator"`
}

type BinlogInfo struct {

	// 备份文件名
	Name *string `json:"Name" name:"Name"`

	// 备份文件大小，单位：Byte
	Size *int64 `json:"Size" name:"Size"`

	// 备份快照时间，时间格式：2016-03-17 02:10:37
	Date *string `json:"Date" name:"Date"`

	// 内网下载地址
	IntranetUrl *string `json:"IntranetUrl" name:"IntranetUrl"`

	// 外网下载地址
	InternetUrl *string `json:"InternetUrl" name:"InternetUrl"`

	// 日志具体类型，可能的值有：binlog - 二进制日志
	Type *string `json:"Type" name:"Type"`
}

type CloseWanServiceRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同，可使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *CloseWanServiceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CloseWanServiceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CloseWanServiceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务的请求ID，可使用此ID查询异步任务的执行结果。
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CloseWanServiceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CloseWanServiceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ColumnPrivilege struct {

	// 数据库名
	Database *string `json:"Database" name:"Database"`

	// 数据库表名
	Table *string `json:"Table" name:"Table"`

	// 数据库列名
	Column *string `json:"Column" name:"Column"`

	// 权限信息
	Privileges []*string `json:"Privileges" name:"Privileges" list`
}

type CreateAccountsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 云数据库账号。
	Accounts []*Account `json:"Accounts" name:"Accounts" list`

	// 新账户的密码。
	Password *string `json:"Password" name:"Password"`

	// 备注信息。
	Description *string `json:"Description" name:"Description"`
}

func (r *CreateAccountsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateAccountsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateAccountsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务的请求ID，可使用此ID查询异步任务的执行结果。
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateAccountsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateAccountsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateBackupRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 目标备份方法，可选的值：logical - 逻辑冷备，physical - 物理冷备。
	BackupMethod *string `json:"BackupMethod" name:"BackupMethod"`
}

func (r *CreateBackupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateBackupRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateBackupResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 备份任务ID。
		BackupId *uint64 `json:"BackupId" name:"BackupId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateBackupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateBackupResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDBImportJobRequest struct {
	*tchttp.BaseRequest

	// 实例的ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 文件名称。该文件是指用户已上传到腾讯云的文件。
	FileName *string `json:"FileName" name:"FileName"`

	// 云数据库的用户名。
	User *string `json:"User" name:"User"`

	// 云数据库实例User账号的密码。
	Password *string `json:"Password" name:"Password"`

	// 导入的目标数据库名，不传表示不指定数据库。
	DbName *string `json:"DbName" name:"DbName"`
}

func (r *CreateDBImportJobRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDBImportJobRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDBImportJobResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务的请求ID，可使用此ID查询异步任务的执行结果
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateDBImportJobResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDBImportJobResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDBInstanceHourRequest struct {
	*tchttp.BaseRequest

	// 实例数量，默认值为1, 最小值1，最大值为100
	GoodsNum *int64 `json:"GoodsNum" name:"GoodsNum"`

	// 实例内存大小，单位：MB，请使用[获取云数据库可售卖规格](https://cloud.tencent.com/document/api/236/17229)接口获取可创建的内存规格
	Memory *int64 `json:"Memory" name:"Memory"`

	// 实例硬盘大小，单位：GB，请使用[获取云数据库可售卖规格](https://cloud.tencent.com/document/api/236/17229)接口获取可创建的硬盘范围
	Volume *int64 `json:"Volume" name:"Volume"`

	// MySQL版本，值包括：5.5、5.6和5.7，请使用[获取云数据库可售卖规格](https://cloud.tencent.com/document/api/236/17229)接口获取可创建的实例版本
	EngineVersion *string `json:"EngineVersion" name:"EngineVersion"`

	// 私有网络ID，如果不传则默认选择基础网络，请使用[查询私有网络列表](/document/api/215/15778)
	UniqVpcId *string `json:"UniqVpcId" name:"UniqVpcId"`

	// 私有网络下的子网ID，如果设置了 UniqVpcId，则 UniqSubnetId 必填，请使用[查询子网列表](/document/api/215/15784)
	UniqSubnetId *string `json:"UniqSubnetId" name:"UniqSubnetId"`

	// 项目ID，不填为默认项目。请使用[查询项目列表](https://cloud.tencent.com/document/product/378/4400)接口获取项目ID
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 可用区信息，该参数缺省时，系统会自动选择一个可用区，请使用[获取云数据库可售卖规格](https://cloud.tencent.com/document/api/236/17229)接口获取可创建的可用区
	Zone *string `json:"Zone" name:"Zone"`

	// 实例ID，购买只读实例或者灾备实例时必填，该字段表示只读实例或者灾备实例的主实例ID，请使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872)接口查询云数据库实例ID
	MasterInstanceId *string `json:"MasterInstanceId" name:"MasterInstanceId"`

	// 实例类型，默认为 master，支持值包括：master-表示主实例，dr-表示灾备实例，ro-表示只读实例
	InstanceRole *string `json:"InstanceRole" name:"InstanceRole"`

	// 主实例的可用区信息，购买灾备实例时必填
	MasterRegion *string `json:"MasterRegion" name:"MasterRegion"`

	// 自定义端口，端口支持范围：[ 1024-65535 ]
	Port *int64 `json:"Port" name:"Port"`

	// 设置root帐号密码，密码规则：8-64个字符，至少包含字母、数字、字符（支持的字符：_+-&=!@#$%^*()）中的两种，购买主实例时可指定该参数，购买只读实例或者灾备实例时指定该参数无意义
	Password *string `json:"Password" name:"Password"`

	// 参数列表，参数格式如ParamList.0.Name=auto_increment_increment&ParamList.0.Value=1。可通过[查询参数列表](/document/product/236/6369)查询支持设置的参数
	ParamList []*ParamInfo `json:"ParamList" name:"ParamList" list`

	// 数据复制方式，默认为0，支持值包括：0-表示异步复制，1-表示半同步复制，2-表示强同步复制，购买主实例时可指定该参数，购买只读实例或者灾备实例时指定该参数无意义
	ProtectMode *int64 `json:"ProtectMode" name:"ProtectMode"`

	// 多可用区域，默认为0，支持值包括：0-表示单可用区，1-表示多可用区，购买主实例时可指定该参数，购买只读实例或者灾备实例时指定该参数无意义
	DeployMode *int64 `json:"DeployMode" name:"DeployMode"`

	// 备库1的可用区ID，默认为zoneId的值，购买主实例时可指定该参数，购买只读实例或者灾备实例时指定该参数无意义
	SlaveZone *string `json:"SlaveZone" name:"SlaveZone"`

	// 备库2的可用区ID，默认为0，购买主实例时可指定该参数，购买只读实例或者灾备实例时指定该参数无意义
	BackupZone *string `json:"BackupZone" name:"BackupZone"`

	// 安全组参数，可使用[查询项目安全组信息](https://cloud.tencent.com/document/api/236/15850)接口查询某个项目的安全组详情
	SecurityGroup []*string `json:"SecurityGroup" name:"SecurityGroup" list`

	// 只读实例信息
	RoGroup *RoGroup `json:"RoGroup" name:"RoGroup"`

	// 自动续费标记，值为0或1。购买按量计费实例该字段无意义
	AutoRenewFlag *int64 `json:"AutoRenewFlag" name:"AutoRenewFlag"`

	// 实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`
}

func (r *CreateDBInstanceHourRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDBInstanceHourRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDBInstanceHourResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 短订单ID
		DealIds []*string `json:"DealIds" name:"DealIds" list`

		// 实例ID列表
		InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateDBInstanceHourResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDBInstanceHourResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 实例内存大小，单位：MB，请使用[获取云数据库可售卖规格](https://cloud.tencent.com/document/api/236/17229)接口获取可创建的内存规格
	Memory *int64 `json:"Memory" name:"Memory"`

	// 实例硬盘大小，单位：GB，请使用[获取云数据库可售卖规格](https://cloud.tencent.com/document/api/236/17229)接口获取可创建的硬盘范围
	Volume *int64 `json:"Volume" name:"Volume"`

	// 实例时长，单位：月，可选值包括[1,2,3,4,5,6,7,8,9,10,11,12,24,36]
	Period *int64 `json:"Period" name:"Period"`

	// 实例数量，默认值为1, 最小值1，最大值为100
	GoodsNum *int64 `json:"GoodsNum" name:"GoodsNum"`

	// 可用区信息，该参数缺省时，系统会自动选择一个可用区，请使用[获取云数据库可售卖规格](https://cloud.tencent.com/document/api/236/17229)接口获取可创建的可用区
	Zone *string `json:"Zone" name:"Zone"`

	// 私有网络ID，如果不传则默认选择基础网络，请使用[查询私有网络列表](/document/api/215/15778)
	UniqVpcId *string `json:"UniqVpcId" name:"UniqVpcId"`

	// 私有网络下的子网ID，如果设置了 UniqVpcId，则 UniqSubnetId 必填，请使用[查询子网列表](/document/api/215/15784)
	UniqSubnetId *string `json:"UniqSubnetId" name:"UniqSubnetId"`

	// 项目ID，不填为默认项目。请使用[查询项目列表](https://cloud.tencent.com/document/product/378/4400)接口获取项目ID
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 自定义端口，端口支持范围：[ 1024-65535 ]
	Port *int64 `json:"Port" name:"Port"`

	// 实例类型，默认为 master，支持值包括：master-表示主实例，dr-表示灾备实例，ro-表示只读实例
	InstanceRole *string `json:"InstanceRole" name:"InstanceRole"`

	// 实例ID，购买只读实例时必填，该字段表示只读实例的主实例ID，请使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872)接口查询云数据库实例ID
	MasterInstanceId *string `json:"MasterInstanceId" name:"MasterInstanceId"`

	// MySQL版本，值包括：5.5、5.6和5.7，请使用[获取云数据库可售卖规格](https://cloud.tencent.com/document/api/236/17229)接口获取可创建的实例版本
	EngineVersion *string `json:"EngineVersion" name:"EngineVersion"`

	// 设置root帐号密码，密码规则：8-64个字符，至少包含字母、数字、字符（支持的字符：_+-&=!@#$%^*()）中的两种，购买主实例时可指定该参数，购买只读实例或者灾备实例时指定该参数无意义
	Password *string `json:"Password" name:"Password"`

	// 数据复制方式，默认为0，支持值包括：0-表示异步复制，1-表示半同步复制，2-表示强同步复制
	ProtectMode *int64 `json:"ProtectMode" name:"ProtectMode"`

	// 多可用区域，默认为0，支持值包括：0-表示单可用区，1-表示多可用区
	DeployMode *int64 `json:"DeployMode" name:"DeployMode"`

	// 备库1的可用区信息，默认为zone的值
	SlaveZone *string `json:"SlaveZone" name:"SlaveZone"`

	// 参数列表，参数格式如ParamList.0.Name=auto_increment&ParamList.0.Value=1。可通过[查询参数列表](/document/product/236/6369)查询支持设置的参数
	ParamList []*ParamInfo `json:"ParamList" name:"ParamList" list`

	// 备库2的可用区ID，默认为0，购买主实例时可指定该参数，购买只读实例或者灾备实例时指定该参数无意义
	BackupZone *string `json:"BackupZone" name:"BackupZone"`

	// 自动续费标记，可选值为：0-不自动续费；1-自动续费
	AutoRenewFlag *int64 `json:"AutoRenewFlag" name:"AutoRenewFlag"`

	// 主实例地域信息，购买灾备实例时，该字段必填
	MasterRegion *string `json:"MasterRegion" name:"MasterRegion"`

	// 安全组参数，可使用[查询项目安全组信息](https://cloud.tencent.com/document/api/236/15850)接口查询某个项目的安全组详情
	SecurityGroup []*string `json:"SecurityGroup" name:"SecurityGroup" list`

	// 只读实例参数
	RoGroup *RoGroup `json:"RoGroup" name:"RoGroup"`

	// 实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`
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

		// 短订单ID
		DealIds []*string `json:"DealIds" name:"DealIds" list`

		// 实例ID列表
		InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type DBSwitchInfo struct {

	// 切换时间，格式为：2017-09-03 01:34:31
	SwitchTime *string `json:"SwitchTime" name:"SwitchTime"`

	// 切换类型，可能的返回值为：TRANSFER - 数据迁移；MASTER2SLAVE - 主备切换；RECOVERY - 主从恢复
	SwitchType *string `json:"SwitchType" name:"SwitchType"`
}

type DatabaseName struct {

	// 数据库表名
	DatabaseName *string `json:"DatabaseName" name:"DatabaseName"`
}

type DatabasePrivilege struct {

	// 权限信息
	Privileges []*string `json:"Privileges" name:"Privileges" list`

	// 数据库名
	Database *string `json:"Database" name:"Database"`
}

type DeleteAccountsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 云数据库账号。
	Accounts []*Account `json:"Accounts" name:"Accounts" list`
}

func (r *DeleteAccountsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteAccountsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteAccountsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务的请求ID，可使用此ID查询异步任务的执行结果。
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteAccountsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteAccountsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteBackupRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 备份任务ID。该任务ID为[创建云数据库备份](https://cloud.tencent.com/document/api/236/15844)接口返回的任务ID。
	BackupId *int64 `json:"BackupId" name:"BackupId"`
}

func (r *DeleteBackupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteBackupRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteBackupResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteBackupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteBackupResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAccountPrivilegesRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 数据库的账号名称。
	User *string `json:"User" name:"User"`

	// 数据库的账号域名。
	Host *string `json:"Host" name:"Host"`
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

		// 全局权限数组。
		GlobalPrivileges []*string `json:"GlobalPrivileges" name:"GlobalPrivileges" list`

		// 数据库权限数组。
		DatabasePrivileges []*DatabasePrivilege `json:"DatabasePrivileges" name:"DatabasePrivileges" list`

		// 数据库中的表权限数组。
		TablePrivileges []*TablePrivilege `json:"TablePrivileges" name:"TablePrivileges" list`

		// 数据库表中的列权限数组。
		ColumnPrivileges []*ColumnPrivilege `json:"ColumnPrivileges" name:"ColumnPrivileges" list`

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

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 记录偏移量，默认值为0。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 单次请求返回的数量，默认值为20，最小值为1，最大值为100。
	Limit *int64 `json:"Limit" name:"Limit"`
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

		// 符合查询条件的账号数量。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 符合查询条件的账号详细信息。
		Items []*AccountInfo `json:"Items" name:"Items" list`

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

type DescribeAsyncRequestInfoRequest struct {
	*tchttp.BaseRequest

	// 异步任务的请求ID。
	AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`
}

func (r *DescribeAsyncRequestInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAsyncRequestInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAsyncRequestInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务执行结果。可能的取值：INITIAL - 初始化，RUNNING - 运行中，SUCCESS - 执行成功，FAILED - 执行失败，KILLED - 已终止，REMOVED - 已删除，PAUSED - 终止中。
		Status *string `json:"Status" name:"Status"`

		// 任务执行信息描述。
		Info *string `json:"Info" name:"Info"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeAsyncRequestInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAsyncRequestInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBackupConfigRequest struct {
	*tchttp.BaseRequest

	// 实例短实例ID，格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *DescribeBackupConfigRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBackupConfigRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBackupConfigResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 备份开始的最早时间点，单位为时刻。例如，2 - 凌晨2:00
		StartTimeMin *int64 `json:"StartTimeMin" name:"StartTimeMin"`

		// 备份开始的最晚时间点，单位为时刻。例如，6 - 凌晨6:00
		StartTimeMax *int64 `json:"StartTimeMax" name:"StartTimeMax"`

		// 备份过期时间，单位为天
		BackupExpireDays *int64 `json:"BackupExpireDays" name:"BackupExpireDays"`

		// 备份方式，可能的值为：physical - 物理备份，logical - 逻辑备份
		BackupMethod *string `json:"BackupMethod" name:"BackupMethod"`

		// Binlog过期时间，单位为天
		BinlogExpireDays *int64 `json:"BinlogExpireDays" name:"BinlogExpireDays"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBackupConfigResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBackupConfigResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBackupDatabasesRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 开始时间，格式为：2017-07-12 10:29:20。
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 要查询的数据库名前缀。
	SearchDatabase *string `json:"SearchDatabase" name:"SearchDatabase"`

	// 分页偏移量。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 分页大小，最小值为1，最大值为2000。
	Limit *int64 `json:"Limit" name:"Limit"`
}

func (r *DescribeBackupDatabasesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBackupDatabasesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBackupDatabasesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 返回的数据个数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 符合查询条件的数据库数组
		Items []*DatabaseName `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBackupDatabasesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBackupDatabasesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBackupTablesRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 开始时间，格式为：2017-07-12 10:29:20。
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 指定的数据库名。
	DatabaseName *string `json:"DatabaseName" name:"DatabaseName"`

	// 要查询的数据表名前缀。
	SearchTable *string `json:"SearchTable" name:"SearchTable"`

	// 分页偏移。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 分页大小，最小值为1，最大值为2000。
	Limit *int64 `json:"Limit" name:"Limit"`
}

func (r *DescribeBackupTablesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBackupTablesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBackupTablesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 返回的数据个数。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 符合条件的数据表数组。
		Items []*TableName `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBackupTablesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBackupTablesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBackupsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 偏移量，最小值为0。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 分页大小，默认值为20，最小值为1，最大值为100。
	Limit *int64 `json:"Limit" name:"Limit"`
}

func (r *DescribeBackupsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBackupsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBackupsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合查询条件的实例总数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 符合查询条件的备份信息详情
		Items []*BackupInfo `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBackupsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBackupsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBinlogsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 偏移量，最小值为0。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 分页大小，默认值为20，最小值为1，最大值为100。
	Limit *int64 `json:"Limit" name:"Limit"`
}

func (r *DescribeBinlogsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBinlogsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBinlogsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合查询条件的日志文件总数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 符合查询条件的二进制日志文件详情
		Items []*BinlogInfo `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBinlogsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBinlogsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBImportRecordsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 开始时间，时间格式如：2016-01-01 00:00:01。
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 结束时间，时间格式如：2016-01-01 23:59:59。
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 分页参数 , 偏移量 , 默认值为0。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 分页参数 , 单次请求返回的数量 , 默认值为20，最小值为1，最大值为100。
	Limit *int64 `json:"Limit" name:"Limit"`
}

func (r *DescribeDBImportRecordsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBImportRecordsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBImportRecordsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合查询条件的导入任务操作日志总数。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 返回的导入操作记录列表。
		Items []*ImportRecord `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBImportRecordsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBImportRecordsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstanceCharsetRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同，可使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *DescribeDBInstanceCharsetRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBInstanceCharsetRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstanceCharsetResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例的默认字符集，如"latin1", "utf8"等。
		Charset *string `json:"Charset" name:"Charset"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBInstanceCharsetResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBInstanceCharsetResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstanceConfigRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *DescribeDBInstanceConfigRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBInstanceConfigRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstanceConfigResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 主库数据保护方式，主实例属性，可能的返回值：0-异步复制方式，1-半同步复制方式，2-强同步复制方式。
		ProtectMode *int64 `json:"ProtectMode" name:"ProtectMode"`

		// 主库部署方式，主实例属性，可能的返回值：0-单可用部署，1-多可用区部署。
		DeployMode *int64 `json:"DeployMode" name:"DeployMode"`

		// 主库可用区的正式名称，如ap-shanghai-1。
		Zone *string `json:"Zone" name:"Zone"`

		// 从库的配置信息。
		SlaveConfig *SlaveConfig `json:"SlaveConfig" name:"SlaveConfig"`

		// ECDB第二个从库的配置信息，只有ECDB实例才有这个字段。
		BackupConfig *BackupConfig `json:"BackupConfig" name:"BackupConfig"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBInstanceConfigResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBInstanceConfigResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstanceGTIDRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同，可使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *DescribeDBInstanceGTIDRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBInstanceGTIDRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstanceGTIDResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// GTID是否开通的标记：0-未开通，1-已开通。
		IsGTIDOpen *int64 `json:"IsGTIDOpen" name:"IsGTIDOpen"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBInstanceGTIDResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBInstanceGTIDResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstanceRebootTimeRequest struct {
	*tchttp.BaseRequest

	// 实例的ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`
}

func (r *DescribeDBInstanceRebootTimeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBInstanceRebootTimeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstanceRebootTimeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合查询条件的实例总数。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 返回的参数信息。
		Items []*InstanceRebootTime `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBInstanceRebootTimeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBInstanceRebootTimeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBInstancesRequest struct {
	*tchttp.BaseRequest

	// 项目ID，可使用[查询项目列表](https://cloud.tencent.com/document/product/378/4400)接口查询项目ID
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 实例类型，可取值：1-主实例，2-灾备实例，3-只读实例
	InstanceTypes []*uint64 `json:"InstanceTypes" name:"InstanceTypes" list`

	// 实例的内网IP地址
	Vips []*string `json:"Vips" name:"Vips" list`

	// 实例状态，可取值：0-创建中，1-运行中，4-隔离中，5-已隔离
	Status []*uint64 `json:"Status" name:"Status" list`

	// 偏移量，默认值为0
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 单次请求返回的数量，默认值为20，最大值为2000
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 安全组ID
	SecurityGroupId *string `json:"SecurityGroupId" name:"SecurityGroupId"`

	// 付费类型，可取值：0-包年包月，1-小时计费
	PayTypes []*uint64 `json:"PayTypes" name:"PayTypes" list`

	// 实例名称
	InstanceNames []*string `json:"InstanceNames" name:"InstanceNames" list`

	// 实例任务状态，可能取值：<br>0-没有任务<br>1-升级中<br>2-数据导入中<br>3-开放Slave中<br>4-外网访问开通中<br>5-批量操作执行中<br>6-回档中<br>7-外网访问关闭中<br>8-密码修改中<br>9-实例名修改中<br>10-重启中<br>12-自建迁移中<br>13-删除库表中<br>14-灾备实例创建同步中
	TaskStatus []*uint64 `json:"TaskStatus" name:"TaskStatus" list`

	// 实例数据库引擎版本，可能取值：5.1、5.5、5.6和5.7
	EngineVersions []*string `json:"EngineVersions" name:"EngineVersions" list`

	// 私有网络的ID
	VpcIds []*uint64 `json:"VpcIds" name:"VpcIds" list`

	// 可用区的ID
	ZoneIds []*uint64 `json:"ZoneIds" name:"ZoneIds" list`

	// 子网ID
	SubnetIds []*uint64 `json:"SubnetIds" name:"SubnetIds" list`

	// 是否锁定标记
	CdbErrors []*int64 `json:"CdbErrors" name:"CdbErrors" list`

	// 返回结果集排序的字段，目前支持："InstanceId", "InstanceName", "CreateTime", "DeadlineTime"
	OrderBy *string `json:"OrderBy" name:"OrderBy"`

	// 返回结果集排序方式，目前支持："ASC"或者"DESC"
	OrderDirection *string `json:"OrderDirection" name:"OrderDirection"`

	// 是否包含安全组详细信息，可取值：0-不包含，1-包含
	WithSecurityGroup *int64 `json:"WithSecurityGroup" name:"WithSecurityGroup"`

	// 是否包含独享集群详细信息，可取值：0-不包含，1-包含
	WithExCluster *int64 `json:"WithExCluster" name:"WithExCluster"`

	// 独享集群ID
	ExClusterId *string `json:"ExClusterId" name:"ExClusterId"`

	// 实例ID
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

	// 初始化标记，可取值：0-未初始化，1-初始化
	InitFlag *int64 `json:"InitFlag" name:"InitFlag"`

	// 是否包含灾备实例，可取值：0-不包含，1-包含
	WithDr *int64 `json:"WithDr" name:"WithDr"`

	// 是否包含只读实例，可取值：0-不包含，1-包含
	WithRo *int64 `json:"WithRo" name:"WithRo"`

	// 是否包含主实例，可取值：0-不包含，1-包含
	WithMaster *int64 `json:"WithMaster" name:"WithMaster"`
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

		// 符合查询条件的实例总数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 实例详细信息
		Items []*InstanceInfo `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type DescribeDBPriceRequest struct {
	*tchttp.BaseRequest

	// 可用区信息，格式如"ap-guangzhou-1"
	Zone *string `json:"Zone" name:"Zone"`

	// 实例数量，默认值为1, 最小值1，最大值为100
	GoodsNum *int64 `json:"GoodsNum" name:"GoodsNum"`

	// 实例内存大小，单位：MB
	Memory *int64 `json:"Memory" name:"Memory"`

	// 实例硬盘大小，单位：GB
	Volume *int64 `json:"Volume" name:"Volume"`

	// 付费类型，支持值包括：PRE_PAID - 包年包月，HOUR_PAID - 按量计费
	PayType *string `json:"PayType" name:"PayType"`

	// 实例时长，单位：月，最小值1，最大值为36；查询按量计费价格时，该字段无效
	Period *int64 `json:"Period" name:"Period"`

	// 实例类型，默认为 master，支持值包括：master-表示主实例，ro-表示只读实例，dr-表示灾备实例
	InstanceRole *string `json:"InstanceRole" name:"InstanceRole"`

	// 数据复制方式，默认为0，支持值包括：0-表示异步复制，1-表示半同步复制，2-表示强同步复制
	ProtectMode *int64 `json:"ProtectMode" name:"ProtectMode"`
}

func (r *DescribeDBPriceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBPriceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBPriceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例价格，单位：分（人民币）
		Price *int64 `json:"Price" name:"Price"`

		// 实例原价，单位：分（人民币）
		OriginalPrice *int64 `json:"OriginalPrice" name:"OriginalPrice"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBPriceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBPriceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBSecurityGroupsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv或者cdbro-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *DescribeDBSecurityGroupsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBSecurityGroupsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBSecurityGroupsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 安全组详情。
		Groups []*SecurityGroup `json:"Groups" name:"Groups" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBSecurityGroupsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBSecurityGroupsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBSwitchRecordsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv或者cdbro-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 分页偏移量。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 分页大小，默认值为50，最小值为1，最大值为2000。
	Limit *int64 `json:"Limit" name:"Limit"`
}

func (r *DescribeDBSwitchRecordsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBSwitchRecordsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBSwitchRecordsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例切换记录的总数。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 实例切换记录详情。
		Items []*DBSwitchInfo `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBSwitchRecordsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBSwitchRecordsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBZoneConfigRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeDBZoneConfigRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBZoneConfigRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBZoneConfigResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 可售卖地域配置数量
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 可售卖地域配置详情
		Items []*RegionSellConf `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBZoneConfigResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBZoneConfigResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDatabasesRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 偏移量，最小值为0。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 单次请求数量，默认值为20，最小值为1，最大值为100。
	Limit *int64 `json:"Limit" name:"Limit"`

	// 匹配数据库库名的正则表达式，规则同MySQL官网
	DatabaseRegexp *string `json:"DatabaseRegexp" name:"DatabaseRegexp"`
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

		// 符合查询条件的实例总数。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 返回的实例信息。
		Items []*string `json:"Items" name:"Items" list`

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

type DescribeInstanceParamsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同，可使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *DescribeInstanceParamsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeInstanceParamsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeInstanceParamsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例的参数总数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 参数详情
		Items []*ParameterDetail `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeInstanceParamsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeInstanceParamsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProjectSecurityGroupsRequest struct {
	*tchttp.BaseRequest

	// 项目ID。
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`
}

func (r *DescribeProjectSecurityGroupsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProjectSecurityGroupsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProjectSecurityGroupsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 安全组详情。
		Groups []*SecurityGroup `json:"Groups" name:"Groups" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeProjectSecurityGroupsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProjectSecurityGroupsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeRollbackRangeTimeRequest struct {
	*tchttp.BaseRequest

	// 实例ID列表，单个实例Id的格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例ID相同。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`
}

func (r *DescribeRollbackRangeTimeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeRollbackRangeTimeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeRollbackRangeTimeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合查询条件的实例总数。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 返回的参数信息。
		Items []*InstanceRollbackRangeTime `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeRollbackRangeTimeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeRollbackRangeTimeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSlowLogsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 偏移量，最小值为0。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 分页大小，默认值为20，最小值为1，最大值为100。
	Limit *int64 `json:"Limit" name:"Limit"`
}

func (r *DescribeSlowLogsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSlowLogsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSlowLogsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合查询条件的慢查询日志总数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 符合查询条件的慢查询日志详情
		Items []*SlowLogInfo `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeSlowLogsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSlowLogsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTablesRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 数据库的名称。
	Database *string `json:"Database" name:"Database"`

	// 记录偏移量，默认值为0。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 单次请求返回的数量，默认值为20，最大值为2000。
	Limit *int64 `json:"Limit" name:"Limit"`

	// 匹配数据库表名的正则表达式，规则同MySQL官网
	TableRegexp *string `json:"TableRegexp" name:"TableRegexp"`
}

func (r *DescribeTablesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTablesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTablesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合查询条件的数据库表总数。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 返回的数据库表信息。
		Items []*string `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeTablesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTablesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTasksRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同，可使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 异步任务请求ID，执行 CDB 相关操作返回的 AsyncRequestId
	AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

	// 任务类型，不传值则查询所有任务类型，可能的值：1-数据库回档；2-SQL操作；3-数据导入；5-参数设置；6-初始化；7-重启；8-开启GTID；9-只读实例升级；10-数据库批量回档；11-主实例升级；12-删除库表；13-切换为主实例；
	TaskTypes []*int64 `json:"TaskTypes" name:"TaskTypes" list`

	// 任务状态，不传值则查询所有任务状态，可能的值：-1-未定义；0-初始化; 1-运行中；2-执行成功；3-执行失败；4-已终止；5-已删除；6-已暂停；
	TaskStatus []*int64 `json:"TaskStatus" name:"TaskStatus" list`

	// 第一个任务的开始时间，用于范围查询，时间格式如：2017-12-31 10:40:01
	StartTimeBegin *string `json:"StartTimeBegin" name:"StartTimeBegin"`

	// 最后一个任务的开始时间，用于范围查询，时间格式如：2017-12-31 10:40:01
	StartTimeEnd *string `json:"StartTimeEnd" name:"StartTimeEnd"`

	// 记录偏移量，默认值为0
	Offset *int64 `json:"Offset" name:"Offset"`

	// 单次请求返回的数量，默认值为20，最大值为100
	Limit *int64 `json:"Limit" name:"Limit"`
}

func (r *DescribeTasksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTasksRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTasksResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合查询条件的实例总数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 返回的实例任务信息
		Items []*string `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeTasksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTasksResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeUploadedFilesRequest struct {
	*tchttp.BaseRequest

	// 文件路径。该字段应填用户主账号的OwnerUin信息。
	Path *string `json:"Path" name:"Path"`

	// 记录偏移量，默认值为0。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 单次请求返回的数量，默认值为20。
	Limit *int64 `json:"Limit" name:"Limit"`
}

func (r *DescribeUploadedFilesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeUploadedFilesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeUploadedFilesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合查询条件的SQL文件总数。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 返回的SQL文件列表。
		Items []*SqlFileInfo `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeUploadedFilesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeUploadedFilesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DisassociateSecurityGroupsRequest struct {
	*tchttp.BaseRequest

	// 安全组Id。
	SecurityGroupId *string `json:"SecurityGroupId" name:"SecurityGroupId"`

	// 实例ID列表，一个或者多个实例Id组成的数组。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`
}

func (r *DisassociateSecurityGroupsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DisassociateSecurityGroupsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DisassociateSecurityGroupsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DisassociateSecurityGroupsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DisassociateSecurityGroupsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DrInfo struct {

	// 灾备实例状态
	Status *int64 `json:"Status" name:"Status"`

	// 可用区信息
	Zone *string `json:"Zone" name:"Zone"`

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 地域信息
	Region *string `json:"Region" name:"Region"`

	// 实例同步状态
	SyncStatus *int64 `json:"SyncStatus" name:"SyncStatus"`

	// 实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// 实例类型
	InstanceType *int64 `json:"InstanceType" name:"InstanceType"`
}

type ImportRecord struct {

	// 状态值
	Status *int64 `json:"Status" name:"Status"`

	// 状态值
	Code *int64 `json:"Code" name:"Code"`

	// 执行时间
	CostTime *int64 `json:"CostTime" name:"CostTime"`

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 后端任务ID
	WorkId *string `json:"WorkId" name:"WorkId"`

	// 导入文件名
	FileName *string `json:"FileName" name:"FileName"`

	// 执行进度
	Process *int64 `json:"Process" name:"Process"`

	// 任务创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 文件大小
	FileSize *string `json:"FileSize" name:"FileSize"`

	// 任务执行信息
	Message *string `json:"Message" name:"Message"`

	// 任务ID
	JobId *int64 `json:"JobId" name:"JobId"`

	// 导入库表名
	DbName *string `json:"DbName" name:"DbName"`

	// 异步任务的请求ID
	AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`
}

type Inbound struct {

	// 策略，ACCEPT或者DROP
	Action *string `json:"Action" name:"Action"`

	// 来源Ip或Ip段，例如192.168.0.0/16
	CidrIp *string `json:"CidrIp" name:"CidrIp"`

	// 端口
	PortRange *string `json:"PortRange" name:"PortRange"`

	// 网络协议，支持udp、tcp等
	IpProtocol *string `json:"IpProtocol" name:"IpProtocol"`

	// 规则限定的方向，进站规则为INPUT
	Dir *string `json:"Dir" name:"Dir"`
}

type InitDBInstancesRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同，可使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

	// 实例新的密码，密码规则：8-64个字符，至少包含字母、数字、字符（支持的字符：!@#$%^*()）中的两种。
	NewPassword *string `json:"NewPassword" name:"NewPassword"`

	// 实例的参数列表，目前支持设置“character_set_server”、“lower_case_table_names”参数。其中，“character_set_server”参数可选值为["utf8","latin1","gbk","utf8mb4"]；“lower_case_table_names”可选值为[“0”,“1”]。
	Parameters []*ParamInfo `json:"Parameters" name:"Parameters" list`

	// 实例的端口，取值范围为[1024, 65535]
	Vport *int64 `json:"Vport" name:"Vport"`
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

		// 异步任务的请求ID数组，可使用此ID查询异步任务的执行结果
		AsyncRequestIds []*string `json:"AsyncRequestIds" name:"AsyncRequestIds" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type InstanceInfo struct {

	// 外网状态，可能的返回值为：0-未开通外网；1-已开通外网；2-已关闭外网
	WanStatus *int64 `json:"WanStatus" name:"WanStatus"`

	// 可用区信息
	Zone *string `json:"Zone" name:"Zone"`

	// 初始化标志，可能的返回值为：0-未初始化；1-已初始化
	InitFlag *int64 `json:"InitFlag" name:"InitFlag"`

	// 只读vip信息。单独开通只读实例访问的只读实例才有该字段
	RoVipInfo *RoVipInfo `json:"RoVipInfo" name:"RoVipInfo"`

	// 内存容量，单位为MB
	Memory *int64 `json:"Memory" name:"Memory"`

	// 实例状态，可能的返回值：0-创建中；1-运行中；4-隔离中；5-已隔离
	Status *int64 `json:"Status" name:"Status"`

	// 私有网络ID，例如：51102
	VpcId *int64 `json:"VpcId" name:"VpcId"`

	// 备机信息
	SlaveInfo *SlaveInfo `json:"SlaveInfo" name:"SlaveInfo"`

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 硬盘容量，单位为GB
	Volume *int64 `json:"Volume" name:"Volume"`

	// 自动续费标志，可能的返回值：0-未开通自动续费；1-已开通自动续费；2-已关闭自动续费
	AutoRenew *int64 `json:"AutoRenew" name:"AutoRenew"`

	// 数据复制方式
	ProtectMode *int64 `json:"ProtectMode" name:"ProtectMode"`

	// 只读组详细信息
	RoGroups []*RoGroup `json:"RoGroups" name:"RoGroups" list`

	// 子网ID，例如：2333
	SubnetId *int64 `json:"SubnetId" name:"SubnetId"`

	// 实例类型，可能的返回值：1-主实例；2-灾备实例；3-只读实例
	InstanceType *int64 `json:"InstanceType" name:"InstanceType"`

	// 项目ID
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 地域信息
	Region *string `json:"Region" name:"Region"`

	// 实例到期时间
	DeadlineTime *string `json:"DeadlineTime" name:"DeadlineTime"`

	// 可用区部署方式
	DeployMode *int64 `json:"DeployMode" name:"DeployMode"`

	// 实例任务状态
	TaskStatus *int64 `json:"TaskStatus" name:"TaskStatus"`

	// 主实例详细信息
	MasterInfo *MasterInfo `json:"MasterInfo" name:"MasterInfo"`

	// 实例类型，可能的返回值：“HA”-高可用版；“BASIC”-基础版
	DeviceType *string `json:"DeviceType" name:"DeviceType"`

	// 内核版本
	EngineVersion *string `json:"EngineVersion" name:"EngineVersion"`

	// 实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// 灾备实例详细信息
	DrInfo []*DrInfo `json:"DrInfo" name:"DrInfo" list`

	// 外网域名
	WanDomain *string `json:"WanDomain" name:"WanDomain"`

	// 外网端口号
	WanPort *int64 `json:"WanPort" name:"WanPort"`

	// 付费类型，可能的返回值：0-包年包月；1-按量计费
	PayType *int64 `json:"PayType" name:"PayType"`

	// 实例创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 实例IP
	Vip *string `json:"Vip" name:"Vip"`

	// 端口号
	Vport *int64 `json:"Vport" name:"Vport"`

	// 是否锁定标记
	CdbError *int64 `json:"CdbError" name:"CdbError"`

	// 私有网络描述符，例如：“vpc-5v8wn9mg”
	UniqVpcId *string `json:"UniqVpcId" name:"UniqVpcId"`

	// 子网描述符，例如：“subnet-1typ0s7d”
	UniqSubnetId *string `json:"UniqSubnetId" name:"UniqSubnetId"`

	// 物理ID
	PhysicalId *string `json:"PhysicalId" name:"PhysicalId"`
}

type InstanceRebootTime struct {

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 预期重启时间
	TimeInSeconds *int64 `json:"TimeInSeconds" name:"TimeInSeconds"`
}

type InstanceRollbackRangeTime struct {

	// 查询数据库错误码
	Code *int64 `json:"Code" name:"Code"`

	// 查询数据库错误信息
	Message *string `json:"Message" name:"Message"`

	// 实例ID列表，单个实例Id的格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例ID相同
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 可回档时间范围
	Times []*RollbackTimeRange `json:"Times" name:"Times" list`
}

type IsolateDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同，可使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *IsolateDBInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *IsolateDBInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type IsolateDBInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务的请求ID，可使用此ID查询异步任务的执行结果
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *IsolateDBInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *IsolateDBInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type MasterInfo struct {

	// 地域信息
	Region *string `json:"Region" name:"Region"`

	// 地域ID
	RegionId *int64 `json:"RegionId" name:"RegionId"`

	// 可用区ID
	ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

	// 可用区信息
	Zone *string `json:"Zone" name:"Zone"`

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 实例长ID
	ResourceId *string `json:"ResourceId" name:"ResourceId"`

	// 实例状态
	Status *int64 `json:"Status" name:"Status"`

	// 实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// 实例类型
	InstanceType *int64 `json:"InstanceType" name:"InstanceType"`

	// 任务状态
	TaskStatus *int64 `json:"TaskStatus" name:"TaskStatus"`

	// 内存容量
	Memory *int64 `json:"Memory" name:"Memory"`

	// 硬盘容量
	Volume *int64 `json:"Volume" name:"Volume"`

	// 实例机型
	DeviceType *string `json:"DeviceType" name:"DeviceType"`

	// 每秒查询数
	Qps *int64 `json:"Qps" name:"Qps"`

	// 私有网络ID
	VpcId *int64 `json:"VpcId" name:"VpcId"`

	// 子网ID
	SubnetId *int64 `json:"SubnetId" name:"SubnetId"`

	// 独享集群ID
	ExClusterId *string `json:"ExClusterId" name:"ExClusterId"`

	// 独享集群名称
	ExClusterName *string `json:"ExClusterName" name:"ExClusterName"`
}

type ModifyAccountDescriptionRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 云数据库账号。
	Accounts []*Account `json:"Accounts" name:"Accounts" list`

	// 数据库账号的备注信息。
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

		// 异步任务的请求ID，可使用此ID查询异步任务的执行结果。
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

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

type ModifyAccountPasswordRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 数据库账号的新密码。
	NewPassword *string `json:"NewPassword" name:"NewPassword"`

	// 云数据库账号。
	Accounts []*Account `json:"Accounts" name:"Accounts" list`
}

func (r *ModifyAccountPasswordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAccountPasswordRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAccountPasswordResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务的请求ID，可使用此ID查询异步任务的执行结果。
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyAccountPasswordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAccountPasswordResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAccountPrivilegesRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 数据库的账号，包括用户名和域名。
	Accounts []*Account `json:"Accounts" name:"Accounts" list`

	// 全局权限。其中，GlobalPrivileges 中权限的可选值为："SELECT","INSERT","UPDATE","DELETE","CREATE",	"DROP","REFERENCES","INDEX","ALTER","SHOW DATABASES","CREATE TEMPORARY TABLES","LOCK TABLES","EXECUTE","CREATE VIEW","SHOW VIEW","CREATE ROUTINE","ALTER ROUTINE","EVENT","TRIGGER"。
	GlobalPrivileges []*string `json:"GlobalPrivileges" name:"GlobalPrivileges" list`

	// 数据库的权限。Privileges权限的可选值为："SELECT","INSERT","UPDATE","DELETE","CREATE",	"DROP","REFERENCES","INDEX","ALTER","CREATE TEMPORARY TABLES","LOCK TABLES","EXECUTE","CREATE VIEW","SHOW VIEW","CREATE ROUTINE","ALTER ROUTINE","EVENT","TRIGGER"。
	DatabasePrivileges []*DatabasePrivilege `json:"DatabasePrivileges" name:"DatabasePrivileges" list`

	// 数据库中表的权限。Privileges权限的可选值为：权限的可选值为："SELECT","INSERT","UPDATE","DELETE","CREATE",	"DROP","REFERENCES","INDEX","ALTER","CREATE VIEW","SHOW VIEW", "TRIGGER"。
	TablePrivileges []*TablePrivilege `json:"TablePrivileges" name:"TablePrivileges" list`

	// 数据库表中列的权限。Privileges权限的可选值为："SELECT","INSERT","UPDATE","REFERENCES"。
	ColumnPrivileges []*ColumnPrivilege `json:"ColumnPrivileges" name:"ColumnPrivileges" list`
}

func (r *ModifyAccountPrivilegesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAccountPrivilegesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAccountPrivilegesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务的请求ID，可使用此ID查询异步任务的执行结果。
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyAccountPrivilegesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAccountPrivilegesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAutoRenewFlagRequest struct {
	*tchttp.BaseRequest

	// 实例的ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

	// 自动续费标记，可取值的有：0-不自动续费，1-自动续费。
	AutoRenew *int64 `json:"AutoRenew" name:"AutoRenew"`
}

func (r *ModifyAutoRenewFlagRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAutoRenewFlagRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAutoRenewFlagResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyAutoRenewFlagResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAutoRenewFlagResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyBackupConfigRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 备份过期时间，单位为天，最小值为7天，最大值为732天。
	ExpireDays *int64 `json:"ExpireDays" name:"ExpireDays"`

	// 备份时间范围，格式为：02:00-06:00，起点和终点时间目前限制为整点，目前可以选择的范围为： 02:00-06:00，06：00-10：00，10:00-14:00，14:00-18:00，18:00-22:00，22:00-02:00。
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 目标备份方法，可选的值：logical - 逻辑冷备，physical - 物理冷备；默认备份方法为 逻辑冷备。
	BackupMethod *string `json:"BackupMethod" name:"BackupMethod"`
}

func (r *ModifyBackupConfigRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyBackupConfigRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyBackupConfigResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyBackupConfigResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyBackupConfigResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBInstanceNameRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同，可使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 实例名称。
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type ModifyDBInstanceProjectRequest struct {
	*tchttp.BaseRequest

	// 实例ID数组，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同，可使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

	// 项目的ID。
	NewProjectId *int64 `json:"NewProjectId" name:"NewProjectId"`
}

func (r *ModifyDBInstanceProjectRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBInstanceProjectRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBInstanceProjectResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDBInstanceProjectResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBInstanceProjectResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBInstanceSecurityGroupsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv或者cdbro-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 要修改的安全组ID列表，一个或者多个安全组Id组成的数组。
	SecurityGroupIds []*string `json:"SecurityGroupIds" name:"SecurityGroupIds" list`
}

func (r *ModifyDBInstanceSecurityGroupsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBInstanceSecurityGroupsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBInstanceSecurityGroupsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDBInstanceSecurityGroupsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBInstanceSecurityGroupsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBInstanceVipVportRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同，可使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 目标IP。该参数和DstPort参数，两者必传一个。
	DstIp *string `json:"DstIp" name:"DstIp"`

	// 目标端口，支持范围为：[1024-65535]。该参数和DstIp参数，两者必传一个。
	DstPort *int64 `json:"DstPort" name:"DstPort"`

	// 私有网络统一ID。
	UniqVpcId *string `json:"UniqVpcId" name:"UniqVpcId"`

	// 子网统一ID。
	UniqSubnetId *string `json:"UniqSubnetId" name:"UniqSubnetId"`
}

func (r *ModifyDBInstanceVipVportRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBInstanceVipVportRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBInstanceVipVportResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务ID，可使用[查询任务列表](https://cloud.tencent.com/document/api/236/8010)获取其执行情况。
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDBInstanceVipVportResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBInstanceVipVportResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyInstanceParamRequest struct {
	*tchttp.BaseRequest

	// 实例短Id列表。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

	// 要修改的参数列表。每一个元素是name和currentValue的组合。name是参数名，currentValue是要修改成的值。
	ParamList []*Parameter `json:"ParamList" name:"ParamList" list`
}

func (r *ModifyInstanceParamRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyInstanceParamRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyInstanceParamResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务Id，可用于查询任务进度。
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyInstanceParamResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyInstanceParamResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type OpenDBInstanceGTIDRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *OpenDBInstanceGTIDRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *OpenDBInstanceGTIDRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type OpenDBInstanceGTIDResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务的请求ID，可使用此ID查询异步任务的执行结果。
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *OpenDBInstanceGTIDResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *OpenDBInstanceGTIDResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type OpenWanServiceRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同，可使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *OpenWanServiceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *OpenWanServiceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type OpenWanServiceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务的请求ID，可使用此ID查询异步任务的执行结果
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *OpenWanServiceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *OpenWanServiceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Outbound struct {

	// 策略，ACCEPT或者DROP
	Action *string `json:"Action" name:"Action"`

	// 目的Ip或Ip段，例如172.16.0.0/12
	CidrIp *string `json:"CidrIp" name:"CidrIp"`

	// 端口或者端口范围
	PortRange *string `json:"PortRange" name:"PortRange"`

	// 网络协议，支持udp、tcp等
	IpProtocol *string `json:"IpProtocol" name:"IpProtocol"`

	// 规则限定的方向，进站规则为OUTPUT
	Dir *string `json:"Dir" name:"Dir"`
}

type ParamInfo struct {

	// 参数名
	Name *string `json:"Name" name:"Name"`

	// 参数值
	Value *string `json:"Value" name:"Value"`
}

type Parameter struct {

	// 参数名称
	Name *string `json:"Name" name:"Name"`

	// 参数值
	CurrentValue *string `json:"CurrentValue" name:"CurrentValue"`
}

type ParameterDetail struct {

	// 参数名称
	Name *string `json:"Name" name:"Name"`

	// 参数类型
	ParamType *string `json:"ParamType" name:"ParamType"`

	// 参数默认值
	Default *string `json:"Default" name:"Default"`

	// 参数描述
	Description *string `json:"Description" name:"Description"`

	// 参数当前值
	CurrentValue *string `json:"CurrentValue" name:"CurrentValue"`

	// 修改参数后，是否需要重启数据库以使参数生效。可能的值包括：0-不需要重启；1-需要重启
	NeedReboot *int64 `json:"NeedReboot" name:"NeedReboot"`

	// 参数允许的最大值
	Max *int64 `json:"Max" name:"Max"`

	// 参数允许的最小值
	Min *int64 `json:"Min" name:"Min"`

	// 参数的可选枚举值。如果为非枚举参数，则为空
	EnumValue []*string `json:"EnumValue" name:"EnumValue" list`
}

type RegionSellConf struct {

	// 地域中文名称
	RegionName *string `json:"RegionName" name:"RegionName"`

	// 所属大区
	Area *string `json:"Area" name:"Area"`

	// 是否为默认地域
	IsDefaultRegion *int64 `json:"IsDefaultRegion" name:"IsDefaultRegion"`

	// 地域名称
	Region *string `json:"Region" name:"Region"`

	// 可用区售卖配置
	ZonesConf []*ZoneSellConf `json:"ZonesConf" name:"ZonesConf" list`
}

type RenewDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 待续费的实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同，可使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872)
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 续费时长，单位：月，可选值包括[1,2,3,4,5,6,7,8,9,10,11,12,24,36]
	TimeSpan *int64 `json:"TimeSpan" name:"TimeSpan"`
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

		// 订单ID
		DealId *string `json:"DealId" name:"DealId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type RestartDBInstancesRequest struct {
	*tchttp.BaseRequest

	// 实例ID数组，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`
}

func (r *RestartDBInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RestartDBInstancesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RestartDBInstancesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务的请求ID，可使用此ID查询异步任务的执行结果。
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RestartDBInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RestartDBInstancesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RoGroup struct {

	// 只读组模式，可选值为：alone-系统自动分配只读组；allinone-新建只读组；join-使用现有只读组
	RoGroupMode *string `json:"RoGroupMode" name:"RoGroupMode"`

	// 只读组ID
	RoGroupId *string `json:"RoGroupId" name:"RoGroupId"`

	// 只读组名称
	RoGroupName *string `json:"RoGroupName" name:"RoGroupName"`

	// 是否启用延迟超限剔除功能，启用该功能后，只读实例与主实例的延迟超过延迟阈值值，只读实例将被隔离。可选值：1-启用；0-不启用
	RoOfflineDelay *int64 `json:"RoOfflineDelay" name:"RoOfflineDelay"`

	// 延迟阀值
	RoMaxDelayTime *int64 `json:"RoMaxDelayTime" name:"RoMaxDelayTime"`

	// 最少实例保留个数，若购买只读实例数量小于设置数量将不做剔除
	MinRoInGroup *int64 `json:"MinRoInGroup" name:"MinRoInGroup"`

	// 读写权重分配模式，可选值：system-系统自动分配；custom-自定义
	WeightMode *string `json:"WeightMode" name:"WeightMode"`

	// 权重值
	Weight *int64 `json:"Weight" name:"Weight"`

	// 只读组中的只读实例详情
	RoInstances []*RoInstanceInfo `json:"RoInstances" name:"RoInstances" list`

	// 只读组的内网IP
	Vip *string `json:"Vip" name:"Vip"`

	// 只读组的内网端口号
	Vport *int64 `json:"Vport" name:"Vport"`
}

type RoInstanceInfo struct {

	// RO组对应的主实例的ID
	MasterInstanceId *string `json:"MasterInstanceId" name:"MasterInstanceId"`

	// RO实例在RO组内的状态，可能的值：online-在线，offline-下线
	RoStatus *string `json:"RoStatus" name:"RoStatus"`

	// RO实例在RO组内上一次下线的时间
	OfflineTime *string `json:"OfflineTime" name:"OfflineTime"`

	// RO实例在RO组内的权重
	Weight *int64 `json:"Weight" name:"Weight"`

	// RO实例所在区域名称，如ap-shanghai
	Region *string `json:"Region" name:"Region"`

	// RO可用区的正式名称，如ap-shanghai-1
	Zone *string `json:"Zone" name:"Zone"`

	// RO实例ID，格式如：cdbro-c1nl9rpv
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// RO实例状态，可能返回值：0-创建中，1-运行中，4-删除中
	Status *int64 `json:"Status" name:"Status"`

	// 实例类型，可能返回值：1-主实例，2-灾备实例，3-只读实例
	InstanceType *int64 `json:"InstanceType" name:"InstanceType"`

	// RO实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// 按量计费状态，可能的取值：1-正常，2-欠费
	HourFeeStatus *int64 `json:"HourFeeStatus" name:"HourFeeStatus"`

	// RO实例任务状态，可能返回值：<br>0-没有任务<br>1-升级中<br>2-数据导入中<br>3-开放Slave中<br>4-外网访问开通中<br>5-批量操作执行中<br>6-回档中<br>7-外网访问关闭中<br>8-密码修改中<br>9-实例名修改中<br>10-重启中<br>12-自建迁移中<br>13-删除库表中<br>14-灾备实例创建同步中
	TaskStatus *int64 `json:"TaskStatus" name:"TaskStatus"`

	// RO实例内存大小，单位：MB
	Memory *int64 `json:"Memory" name:"Memory"`

	// RO实例硬盘大小，单位：GB
	Volume *int64 `json:"Volume" name:"Volume"`

	// 每次查询数量
	Qps *int64 `json:"Qps" name:"Qps"`

	// RO实例的内网IP地址
	Vip *string `json:"Vip" name:"Vip"`

	// RO实例访问端口
	Vport *int64 `json:"Vport" name:"Vport"`

	// RO实例所在私有网络ID
	VpcId *int64 `json:"VpcId" name:"VpcId"`

	// RO实例所在私有网络子网ID
	SubnetId *int64 `json:"SubnetId" name:"SubnetId"`

	// RO实例规格描述，目前可取值 CUSTOM
	DeviceType *string `json:"DeviceType" name:"DeviceType"`

	// RO实例数据库引擎版本，可能返回值：5.1、5.5、5.6和5.7
	EngineVersion *string `json:"EngineVersion" name:"EngineVersion"`

	// RO实例到期时间，时间格式：yyyy-mm-dd hh:mm:ss，如实例为按量计费模式，则此字段值为0000-00-00 00:00:00
	DeadlineTime *string `json:"DeadlineTime" name:"DeadlineTime"`

	// RO实例计费类型，可能返回值：0-包年包月，1-按量计费，2-后付费月结
	PayType *int64 `json:"PayType" name:"PayType"`
}

type RoVipInfo struct {

	// 只读vip状态
	RoVipStatus *int64 `json:"RoVipStatus" name:"RoVipStatus"`

	// 只读vip的子网
	RoSubnetId *int64 `json:"RoSubnetId" name:"RoSubnetId"`

	// 只读vip的私有网络
	RoVpcId *int64 `json:"RoVpcId" name:"RoVpcId"`

	// 只读vip的端口号
	RoVport *int64 `json:"RoVport" name:"RoVport"`

	// 只读vip
	RoVip *string `json:"RoVip" name:"RoVip"`
}

type RollbackDBName struct {

	// 回档前的原数据库名
	DatabaseName *string `json:"DatabaseName" name:"DatabaseName"`

	// 回档后的新数据库名
	NewDatabaseName *string `json:"NewDatabaseName" name:"NewDatabaseName"`
}

type RollbackInstancesInfo struct {

	// 云数据库实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 回档策略。可选值为：table、db、full；默认值为full。table - 急速回档模式，仅导入所选中表级别的备份和binlog，如有跨表操作，且关联表未被同时选中，将会导致回档失败，该模式下参数Databases必须为空；db - 快速模式，仅导入所选中库级别的备份和binlog，如有跨库操作，且关联库未被同时选中，将会导致回档失败；full - 普通回档模式，将导入整个实例的备份和binlog，速度较慢。
	Strategy *string `json:"Strategy" name:"Strategy"`

	// 数据库回档时间，时间格式为：yyyy-mm-dd hh:mm:ss
	RollbackTime *string `json:"RollbackTime" name:"RollbackTime"`

	// 待回档的数据库信息，表示整库回档
	Databases []*RollbackDBName `json:"Databases" name:"Databases" list`

	// 待回档的数据库表信息，表示按表回档
	Tables []*RollbackTables `json:"Tables" name:"Tables" list`
}

type RollbackTableName struct {

	// 回档前的原数据库表名
	TableName *string `json:"TableName" name:"TableName"`

	// 回档后的新数据库表名
	NewTableName *string `json:"NewTableName" name:"NewTableName"`
}

type RollbackTables struct {

	// 数据库名
	Database *string `json:"Database" name:"Database"`

	// 数据库表详情
	Table []*RollbackTableName `json:"Table" name:"Table" list`
}

type RollbackTimeRange struct {

	// 实例可回档开始时间，时间格式：2016-10-29 01:06:04
	Begin *string `json:"Begin" name:"Begin"`

	// 实例可回档结束时间，时间格式：2016-11-02 11:44:47
	End *string `json:"End" name:"End"`
}

type SecurityGroup struct {

	// 项目ID
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 创建时间，时间格式：yyyy-mm-dd hh:mm:ss
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 入站规则
	Inbound []*Inbound `json:"Inbound" name:"Inbound" list`

	// 出站规则
	Outbound []*Outbound `json:"Outbound" name:"Outbound" list`

	// 安全组ID
	SecurityGroupId *string `json:"SecurityGroupId" name:"SecurityGroupId"`

	// 安全组名称
	SecurityGroupName *string `json:"SecurityGroupName" name:"SecurityGroupName"`

	// 安全组备注
	SecurityGroupRemark *string `json:"SecurityGroupRemark" name:"SecurityGroupRemark"`
}

type SellConfig struct {

	// 设备类型
	Device *string `json:"Device" name:"Device"`

	// 售卖规格描述
	Type *string `json:"Type" name:"Type"`

	// 实例类型
	CdbType *string `json:"CdbType" name:"CdbType"`

	// 内存大小，单位为MB
	Memory *int64 `json:"Memory" name:"Memory"`

	// CPU核心数
	Cpu *int64 `json:"Cpu" name:"Cpu"`

	// 磁盘最小规格，单位为GB
	VolumeMin *int64 `json:"VolumeMin" name:"VolumeMin"`

	// 磁盘最大规格，单位为GB
	VolumeMax *int64 `json:"VolumeMax" name:"VolumeMax"`

	// 磁盘步长，单位为GB
	VolumeStep *int64 `json:"VolumeStep" name:"VolumeStep"`

	// 链接数
	Connection *int64 `json:"Connection" name:"Connection"`

	// 每秒查询数量
	Qps *int64 `json:"Qps" name:"Qps"`

	// 每秒IO数量
	Iops *int64 `json:"Iops" name:"Iops"`

	// 应用场景描述
	Info *string `json:"Info" name:"Info"`

	// 状态值
	Status *int64 `json:"Status" name:"Status"`
}

type SellType struct {

	// 售卖实例名称
	TypeName *string `json:"TypeName" name:"TypeName"`

	// 内核版本号
	EngineVersion []*string `json:"EngineVersion" name:"EngineVersion" list`

	// 售卖规格详细配置
	Configs []*SellConfig `json:"Configs" name:"Configs" list`
}

type SlaveConfig struct {

	// 从库复制方式，可能的返回值：aysnc-异步，semisync-半同步
	ReplicationMode *string `json:"ReplicationMode" name:"ReplicationMode"`

	// 从库可用区的正式名称，如ap-shanghai-1
	Zone *string `json:"Zone" name:"Zone"`
}

type SlaveInfo struct {

	// 第一备机信息
	First *SlaveInstanceInfo `json:"First" name:"First"`

	// 第二备机信息
	Second *SlaveInstanceInfo `json:"Second" name:"Second"`
}

type SlaveInstanceInfo struct {

	// 端口号
	Vport *int64 `json:"Vport" name:"Vport"`

	// 地域信息
	Region *string `json:"Region" name:"Region"`

	// 虚拟Ip信息
	Vip *string `json:"Vip" name:"Vip"`

	// 可用区信息
	Zone *string `json:"Zone" name:"Zone"`
}

type SlowLogInfo struct {

	// 备份文件名
	Name *string `json:"Name" name:"Name"`

	// 备份文件大小，单位：Byte
	Size *int64 `json:"Size" name:"Size"`

	// 备份快照时间，时间格式：2016-03-17 02:10:37
	Date *string `json:"Date" name:"Date"`

	// 内网下载地址
	IntranetUrl *string `json:"IntranetUrl" name:"IntranetUrl"`

	// 外网下载地址
	InternetUrl *string `json:"InternetUrl" name:"InternetUrl"`

	// 日志具体类型，可能的值：slowlog - 慢日志
	Type *string `json:"Type" name:"Type"`
}

type SqlFileInfo struct {

	// 上传时间
	UploadTime *string `json:"UploadTime" name:"UploadTime"`

	// 上传进度
	UploadInfo *UploadInfo `json:"UploadInfo" name:"UploadInfo"`

	// 文件名
	FileName *string `json:"FileName" name:"FileName"`

	// 文件大小，单位为Bytes
	FileSize *int64 `json:"FileSize" name:"FileSize"`

	// 上传是否完成标志，可选值：0 - 未完成，1 - 已完成
	IsUploadFinished *int64 `json:"IsUploadFinished" name:"IsUploadFinished"`

	// 文件ID
	FileId *string `json:"FileId" name:"FileId"`
}

type StartBatchRollbackRequest struct {
	*tchttp.BaseRequest

	// 用于回档的实例详情信息
	Instances []*RollbackInstancesInfo `json:"Instances" name:"Instances" list`
}

func (r *StartBatchRollbackRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *StartBatchRollbackRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type StartBatchRollbackResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务的请求ID，可使用此ID查询异步任务的执行结果
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *StartBatchRollbackResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *StartBatchRollbackResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type StopDBImportJobRequest struct {
	*tchttp.BaseRequest

	// 异步任务的请求ID。
	AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`
}

func (r *StopDBImportJobRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *StopDBImportJobRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type StopDBImportJobResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *StopDBImportJobResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *StopDBImportJobResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SwitchForUpgradeRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

func (r *SwitchForUpgradeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SwitchForUpgradeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SwitchForUpgradeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *SwitchForUpgradeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SwitchForUpgradeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type TableName struct {

	// 表名
	TableName *string `json:"TableName" name:"TableName"`
}

type TablePrivilege struct {

	// 数据库名
	Database *string `json:"Database" name:"Database"`

	// 数据库表名
	Table *string `json:"Table" name:"Table"`

	// 权限信息
	Privileges []*string `json:"Privileges" name:"Privileges" list`
}

type UpgradeDBInstanceEngineVersionRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv或者cdbro-c1nl9rpv。与云数据库控制台页面中显示的实例ID相同，可使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 主实例数据库引擎版本，支持值包括：5.6和5.7
	EngineVersion *string `json:"EngineVersion" name:"EngineVersion"`

	// 切换访问新实例的方式，默认为0，升级主实例时，可指定该参数，升级只读实例或者灾备实例时指定该参数无意义，支持值包括：0-立刻切换，1-时间窗切换；当该值为1时，升级中过程中，切换访问新实例的流程将会在时间窗内进行，或者用户主动调用接口[切换访问新实例](https://cloud.tencent.com/document/api/403/4392)触发该流程
	WaitSwitch *int64 `json:"WaitSwitch" name:"WaitSwitch"`
}

func (r *UpgradeDBInstanceEngineVersionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpgradeDBInstanceEngineVersionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpgradeDBInstanceEngineVersionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务ID，可使用[查询任务列表](https://cloud.tencent.com/document/api/236/8010)获取其执行情况
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UpgradeDBInstanceEngineVersionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpgradeDBInstanceEngineVersionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpgradeDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv或者cdbro-c1nl9rpv。与云数据库控制台页面中显示的实例ID相同，可使用[查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 升级后的内存大小，单位：MB，为保证传入 Memory 值有效，请使用[查询可创建规格（支持可用区、配置自定义）](https://cloud.tencent.com/document/api/253/6109)接口获取可升级的内存规格
	Memory *int64 `json:"Memory" name:"Memory"`

	// 升级后的硬盘大小，单位：GB，为保证传入 Volume 值有效，请使用[查询可创建规格（支持可用区、配置自定义）](https://cloud.tencent.com/document/api/253/6109)接口获取可升级的硬盘范围
	Volume *int64 `json:"Volume" name:"Volume"`

	// 数据复制方式，支持值包括：0-异步复制，1-半同步复制，2-强同步复制，升级主实例时可指定该参数，升级只读实例或者灾备实例时指定该参数无意义
	ProtectMode *int64 `json:"ProtectMode" name:"ProtectMode"`

	// 部署模式，默认为0，支持值包括：0-单可用区部署，1-多可用区部署，升级主实例时可指定该参数，升级只读实例或者灾备实例时指定该参数无意义
	DeployMode *int64 `json:"DeployMode" name:"DeployMode"`

	// 备库1的可用区信息，默认为实例的Zone，升级主实例为多可用区部署时可指定该参数，升级只读实例或者灾备实例时指定该参数无意义。可通过<a href='/document/product/236/6921' title='查询云数据库可售卖规格'>查询云数据库可售卖规格</a>查询支持的可用区
	SlaveZone *string `json:"SlaveZone" name:"SlaveZone"`

	// 主实例数据库引擎版本，支持值包括：5.5、5.6和5.7
	EngineVersion *string `json:"EngineVersion" name:"EngineVersion"`

	// 切换访问新实例的方式，默认为0，升级主实例时，可指定该参数，升级只读实例或者灾备实例时指定该参数无意义，支持值包括：0-立刻切换，1-时间窗切换；当该值为1时，升级中过程中，切换访问新实例的流程将会在时间窗内进行，或者用户主动调用接口[切换访问新实例](https://cloud.tencent.com/document/api/403/4392)触发该流程
	WaitSwitch *int64 `json:"WaitSwitch" name:"WaitSwitch"`

	// 备库2的可用区ID，默认为0，升级主实例时可指定该参数，升级只读实例或者灾备实例时指定该参数无意义
	BackupZone *string `json:"BackupZone" name:"BackupZone"`

	// 实例类型，默认为 master，支持值包括：master-表示主实例，dr-表示灾备实例，ro-表示只读实例
	InstanceRole *string `json:"InstanceRole" name:"InstanceRole"`
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

		// 订单ID，用于调用云API相关接口，如[获取订单信息](https://cloud.tencent.com/document/api/403/4392)
		DealIds []*string `json:"DealIds" name:"DealIds" list`

		// 异步任务的请求ID，可使用此ID查询异步任务的执行结果
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type UploadInfo struct {

	// 文件所有分片数
	AllSliceNum *int64 `json:"AllSliceNum" name:"AllSliceNum"`

	// 已完成分片数
	CompleteNum *int64 `json:"CompleteNum" name:"CompleteNum"`
}

type VerifyRootAccountRequest struct {
	*tchttp.BaseRequest

	// 实例ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例ID相同。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 实例ROOT账号的密码。
	Password *string `json:"Password" name:"Password"`
}

func (r *VerifyRootAccountRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *VerifyRootAccountRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type VerifyRootAccountResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务的请求ID，可使用此ID查询异步任务的执行结果
		AsyncRequestId *string `json:"AsyncRequestId" name:"AsyncRequestId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *VerifyRootAccountResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *VerifyRootAccountResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ZoneConf struct {

	// 可用区部署方式，可能的值为：0-单可用区；1-多可用区
	DeployMode []*int64 `json:"DeployMode" name:"DeployMode" list`

	// 主实例所在的可用区
	MasterZone []*string `json:"MasterZone" name:"MasterZone" list`

	// 实例为多可用区部署时，备库1所在的可用区
	SlaveZone []*string `json:"SlaveZone" name:"SlaveZone" list`

	// 实例为多可用区部署时，备库2所在的可用区
	BackupZone []*string `json:"BackupZone" name:"BackupZone" list`
}

type ZoneSellConf struct {

	// 可用区状态。可能的返回值为：0-未上线；1-上线；2-开放；3-停售；4-不展示
	Status *int64 `json:"Status" name:"Status"`

	// 可用区中文名称
	ZoneName *string `json:"ZoneName" name:"ZoneName"`

	// 实例类型是否为自定义类型
	IsCustom *bool `json:"IsCustom" name:"IsCustom"`

	// 是否支持灾备
	IsSupportDr *bool `json:"IsSupportDr" name:"IsSupportDr"`

	// 是否支持私有网络
	IsSupportVpc *bool `json:"IsSupportVpc" name:"IsSupportVpc"`

	// 小时计费实例最大售卖数量
	HourInstanceSaleMaxNum *int64 `json:"HourInstanceSaleMaxNum" name:"HourInstanceSaleMaxNum"`

	// 是否为默认可用区
	IsDefaultZone *bool `json:"IsDefaultZone" name:"IsDefaultZone"`

	// 是否为黑石区
	IsBm *bool `json:"IsBm" name:"IsBm"`

	// 支持的付费类型。可能的返回值为：0-包年包月；1-小时计费；2-后付费
	PayType []*string `json:"PayType" name:"PayType" list`

	// 数据复制类型。0-异步复制；1-半同步复制；2-强同步复制
	ProtectMode []*string `json:"ProtectMode" name:"ProtectMode" list`

	// 可用区名称
	Zone *string `json:"Zone" name:"Zone"`

	// 售卖实例类型数组
	SellType []*SellType `json:"SellType" name:"SellType" list`

	// 多可用区信息
	ZoneConf *ZoneConf `json:"ZoneConf" name:"ZoneConf"`
}
