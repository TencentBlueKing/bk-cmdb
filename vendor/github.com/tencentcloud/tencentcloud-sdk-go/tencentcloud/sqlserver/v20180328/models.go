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

package v20180328

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type AccountCreateInfo struct {

	// 实例用户名
	UserName *string `json:"UserName" name:"UserName"`

	// 实例密码
	Password *string `json:"Password" name:"Password"`

	// DB权限列表
	DBPrivileges []*DBPrivilege `json:"DBPrivileges" name:"DBPrivileges" list`

	// 账号备注信息
	Remark *string `json:"Remark" name:"Remark"`
}

type AccountDetail struct {

	// 账户名
	Name *string `json:"Name" name:"Name"`

	// 账户备注
	Remark *string `json:"Remark" name:"Remark"`

	// 账户创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 账户状态，1-创建中，2-正常，3-修改中，4-密码重置中，-1-删除中
	Status *int64 `json:"Status" name:"Status"`

	// 账户更新时间
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`

	// 密码更新时间
	PassTime *string `json:"PassTime" name:"PassTime"`

	// 账户内部状态，正常为enable
	InternalStatus *string `json:"InternalStatus" name:"InternalStatus"`

	// 该账户对相关db的读写权限信息
	Dbs []*DBPrivilege `json:"Dbs" name:"Dbs" list`
}

type AccountPassword struct {

	// 用户名
	UserName *string `json:"UserName" name:"UserName"`

	// 密码
	Password *string `json:"Password" name:"Password"`
}

type AccountPrivilege struct {

	// 数据库用户名
	UserName *string `json:"UserName" name:"UserName"`

	// 数据库权限。ReadWrite表示可读写，ReadOnly表示只读
	Privilege *string `json:"Privilege" name:"Privilege"`
}

type AccountPrivilegeModifyInfo struct {

	// 数据库用户名
	UserName *string `json:"UserName" name:"UserName"`

	// 账号权限变更信息
	DBPrivileges []*DBPrivilegeModifyInfo `json:"DBPrivileges" name:"DBPrivileges" list`
}

type AccountRemark struct {

	// 账户名
	UserName *string `json:"UserName" name:"UserName"`

	// 对应账户新的备注信息
	Remark *string `json:"Remark" name:"Remark"`
}

type Backup struct {

	// 文件名
	FileName *string `json:"FileName" name:"FileName"`

	// 文件大小，单位 KB
	Size *int64 `json:"Size" name:"Size"`

	// 备份开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 备份结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 内网下载地址
	InternalAddr *string `json:"InternalAddr" name:"InternalAddr"`

	// 外网下载地址
	ExternalAddr *string `json:"ExternalAddr" name:"ExternalAddr"`

	// 备份文件唯一标识，RestoreInstance接口会用到该字段
	Id *uint64 `json:"Id" name:"Id"`

	// 备份文件状态（0-创建中；1-成功；2-失败）
	Status *uint64 `json:"Status" name:"Status"`

	// 多库备份时的DB列表
	DBs []*string `json:"DBs" name:"DBs" list`

	// 备份策略（0-实例备份；1-多库备份）
	Strategy *int64 `json:"Strategy" name:"Strategy"`

	// 备份方式，0-定时备份；1-手动临时备份
	BackupWay *int64 `json:"BackupWay" name:"BackupWay"`
}

type CreateAccountRequest struct {
	*tchttp.BaseRequest

	// 数据库实例ID，形如mssql-njj2mtpl
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 数据库实例账户信息
	Accounts []*AccountCreateInfo `json:"Accounts" name:"Accounts" list`
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

		// 任务流id
		FlowId *int64 `json:"FlowId" name:"FlowId"`

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

type CreateBackupRequest struct {
	*tchttp.BaseRequest

	// 备份策略(0-实例备份 1-多库备份)
	Strategy *int64 `json:"Strategy" name:"Strategy"`

	// 需要备份库名的列表(多库备份才填写)
	DBNames []*string `json:"DBNames" name:"DBNames" list`

	// 实例ID，形如mssql-i1z41iwd
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
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

		// 异步任务ID
		FlowId *int64 `json:"FlowId" name:"FlowId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
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

type CreateDBInstancesRequest struct {
	*tchttp.BaseRequest

	// 实例可用区，类似ap-guangzhou-1（广州一区）；实例可售卖区域可以通过接口DescribeZones获取
	Zone *string `json:"Zone" name:"Zone"`

	// 实例内存大小，单位GB
	Memory *int64 `json:"Memory" name:"Memory"`

	// 实例磁盘大小，单位GB
	Storage *int64 `json:"Storage" name:"Storage"`

	// 付费模式，目前只支持预付费，其值为PREPAID。可不填，默认值为PREPAID
	InstanceChargeType *string `json:"InstanceChargeType" name:"InstanceChargeType"`

	// 项目ID
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 本次购买几个实例，默认值为1。取值不超过10
	GoodsNum *int64 `json:"GoodsNum" name:"GoodsNum"`

	// VPC子网ID，形如subnet-bdoe83fa；SubnetId和VpcId需同时设置或者同时不设置
	SubnetId *string `json:"SubnetId" name:"SubnetId"`

	// VPC网络ID，形如vpc-dsp338hz；SubnetId和VpcId需同时设置或者同时不设置
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 购买实例周期，默认取值为1，表示一个月。取值不超过48
	Period *int64 `json:"Period" name:"Period"`

	// 是否自动使用代金券；1 - 是，0 - 否，默认不使用
	AutoVoucher *int64 `json:"AutoVoucher" name:"AutoVoucher"`

	// 代金券ID数组，目前单个订单只能使用一张
	VoucherIds []*string `json:"VoucherIds" name:"VoucherIds" list`

	// 数据库版本号，目前取值有2012SP3，表示SQL Server 2012；2008R2，表示SQL Server 2008 R2；2016SP1，表示SQL Server 2016 SP1。每个地域支持售卖的版本可能不一样，可以通过DescribeZones接口来拉取每个地域可售卖的版本信息。不填的话，默认为版本2008R2
	DBVersion *string `json:"DBVersion" name:"DBVersion"`
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

		// 订单名称
		DealName *string `json:"DealName" name:"DealName"`

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

type CreateDBRequest struct {
	*tchttp.BaseRequest

	// 实例id
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 数据库创建信息
	DBs []*DBCreateInfo `json:"DBs" name:"DBs" list`
}

func (r *CreateDBRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDBRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDBResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务流id
		FlowId *int64 `json:"FlowId" name:"FlowId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateDBResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDBResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateMigrationRequest struct {
	*tchttp.BaseRequest

	// 迁移任务的名称
	MigrateName *string `json:"MigrateName" name:"MigrateName"`

	// 迁移类型（1:结构迁移 2:数据迁移 3:增量同步）
	MigrateType *uint64 `json:"MigrateType" name:"MigrateType"`

	// 迁移源的类型 1:CDB for SQLServer 2:云服务器自建SQLServer数据库 4:SQLServer备份还原 5:SQLServer备份还原（COS方式）
	SourceType *uint64 `json:"SourceType" name:"SourceType"`

	// 迁移源
	Source *MigrateSource `json:"Source" name:"Source"`

	// 迁移目标
	Target *MigrateTarget `json:"Target" name:"Target"`

	// 迁移DB对象 ，离线迁移不使用（SourceType=4或SourceType=5）。
	MigrateDBSet []*MigrateDB `json:"MigrateDBSet" name:"MigrateDBSet" list`
}

func (r *CreateMigrationRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateMigrationRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateMigrationResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 迁移任务ID
		MigrateId *int64 `json:"MigrateId" name:"MigrateId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateMigrationResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateMigrationResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DBCreateInfo struct {

	// 数据库名
	DBName *string `json:"DBName" name:"DBName"`

	// 字符集。可选值包括：Chinese_PRC_CI_AS, Chinese_PRC_CS_AS, Chinese_PRC_BIN, Chinese_Taiwan_Stroke_CI_AS, SQL_Latin1_General_CP1_CI_AS, SQL_Latin1_General_CP1_CS_AS。不填默认为Chinese_PRC_CI_AS
	Charset *string `json:"Charset" name:"Charset"`

	// 数据库账号权限信息
	Accounts []*AccountPrivilege `json:"Accounts" name:"Accounts" list`

	// 备注
	Remark *string `json:"Remark" name:"Remark"`
}

type DBDetail struct {

	// 实例id
	Name *string `json:"Name" name:"Name"`

	// 字符集
	Charset *string `json:"Charset" name:"Charset"`

	// 备注
	Remark *string `json:"Remark" name:"Remark"`

	// 数据库创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 数据库状态。1--创建中， 2--运行中， 3--修改中，-1--删除中
	Status *int64 `json:"Status" name:"Status"`

	// 数据库账号权限信息
	Accounts []*AccountPrivilege `json:"Accounts" name:"Accounts" list`

	// 内部状态。ONLINE表示运行中
	InternalStatus *string `json:"InternalStatus" name:"InternalStatus"`
}

type DBInstance struct {

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 实例名称
	Name *string `json:"Name" name:"Name"`

	// 实例所在项目ID
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 实例所在地域ID
	RegionId *int64 `json:"RegionId" name:"RegionId"`

	// 实例所在可用区ID
	ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

	// 实例所在私有网络ID，基础网络时为 0
	VpcId *int64 `json:"VpcId" name:"VpcId"`

	// 实例所在私有网络子网ID，基础网络时为 0
	SubnetId *int64 `json:"SubnetId" name:"SubnetId"`

	// 实例状态。取值范围： <li>1：申请中</li> <li>2：运行中</li> <li>3：受限运行中 (主备切换中)</li> <li>4：已隔离</li> <li>5：回收中</li> <li>6：已回收</li> <li>7：任务执行中 (实例做备份、回档等操作)</li> <li>8：已下线</li> <li>9：实例扩容中</li> <li>10：实例迁移中</li> <li>11：只读</li> <li>12：重启中</li>
	Status *int64 `json:"Status" name:"Status"`

	// 实例访问IP
	Vip *string `json:"Vip" name:"Vip"`

	// 实例访问端口
	Vport *int64 `json:"Vport" name:"Vport"`

	// 实例创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 实例更新时间
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`

	// 实例计费开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 实例计费结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 实例隔离时间
	IsolateTime *string `json:"IsolateTime" name:"IsolateTime"`

	// 实例内存大小，单位G
	Memory *int64 `json:"Memory" name:"Memory"`

	// 实例已经使用存储空间大小，单位G
	UsedStorage *int64 `json:"UsedStorage" name:"UsedStorage"`

	// 实例存储空间大小，单位G
	Storage *int64 `json:"Storage" name:"Storage"`

	// 实例版本
	VersionName *string `json:"VersionName" name:"VersionName"`

	// 实例续费标记，0-正常续费，1-自动续费，2-到期不续费
	RenewFlag *int64 `json:"RenewFlag" name:"RenewFlag"`

	// 实例高可用， 1-双机高可用，2-单机
	Model *int64 `json:"Model" name:"Model"`

	// 实例所在地域名称，如 ap-guangzhou
	Region *string `json:"Region" name:"Region"`

	// 实例所在可用区名称，如 ap-guangzhou-1
	Zone *string `json:"Zone" name:"Zone"`

	// 备份时间点
	BackupTime *string `json:"BackupTime" name:"BackupTime"`
}

type DBPrivilege struct {

	// 数据库名
	DBName *string `json:"DBName" name:"DBName"`

	// 数据库权限，ReadWrite表示可读写，ReadOnly表示只读
	Privilege *string `json:"Privilege" name:"Privilege"`
}

type DBPrivilegeModifyInfo struct {

	// 数据库名
	DBName *string `json:"DBName" name:"DBName"`

	// 权限变更信息。ReadWrite表示可读写，ReadOnly表示只读，Delete表示删除账号对该DB的权限
	Privilege *string `json:"Privilege" name:"Privilege"`
}

type DBRemark struct {

	// 据库名
	Name *string `json:"Name" name:"Name"`

	// 备注信息
	Remark *string `json:"Remark" name:"Remark"`
}

type DbRollbackTimeInfo struct {

	// 数据库名称
	DBName *string `json:"DBName" name:"DBName"`

	// 可回档开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 可回档结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`
}

type DealInfo struct {

	// 订单名
	DealName *string `json:"DealName" name:"DealName"`

	// 商品数量
	Count *uint64 `json:"Count" name:"Count"`

	// 关联的流程 Id，可用于查询流程执行状态
	FlowId *int64 `json:"FlowId" name:"FlowId"`

	// 只有创建实例的订单会填充该字段，表示该订单创建的实例的 ID。
	InstanceIdSet []*string `json:"InstanceIdSet" name:"InstanceIdSet" list`

	// 所属账号
	OwnerUin *string `json:"OwnerUin" name:"OwnerUin"`
}

type DeleteAccountRequest struct {
	*tchttp.BaseRequest

	// 数据库实例ID，形如mssql-njj2mtpl
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 实例用户名数组
	UserNames []*string `json:"UserNames" name:"UserNames" list`
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

		// 任务流id
		FlowId *int64 `json:"FlowId" name:"FlowId"`

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

type DeleteDBRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如mssql-rljoi3bf
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 数据库名数组
	Names []*string `json:"Names" name:"Names" list`
}

func (r *DeleteDBRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteDBRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteDBResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务流id
		FlowId *int64 `json:"FlowId" name:"FlowId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteDBResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteDBResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteMigrationRequest struct {
	*tchttp.BaseRequest

	// 迁移任务ID
	MigrateId *uint64 `json:"MigrateId" name:"MigrateId"`
}

func (r *DeleteMigrationRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteMigrationRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteMigrationResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteMigrationResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteMigrationResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAccountsRequest struct {
	*tchttp.BaseRequest

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 分页返回，每页返回的数目，取值为1-100，默认值为20
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 分页返回，从第几页开始返回。从第0页开始，默认第0页
	Offset *uint64 `json:"Offset" name:"Offset"`
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

		// 实例ID
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 账户信息列表
		Accounts []*AccountDetail `json:"Accounts" name:"Accounts" list`

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

type DescribeBackupsRequest struct {
	*tchttp.BaseRequest

	// 开始时间(yyyy-MM-dd HH:mm:ss)
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 结束时间(yyyy-MM-dd HH:mm:ss)
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 实例ID，形如mssql-njj2mtpl
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 分页返回，每页返回数量，默认为20，最大值为 100
	Limit *int64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为 0
	Offset *int64 `json:"Offset" name:"Offset"`
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

		// 备份总数量
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 备份列表详情
		Backups []*Backup `json:"Backups" name:"Backups" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
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

type DescribeDBInstancesRequest struct {
	*tchttp.BaseRequest

	// 项目ID
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`

	// 实例状态。取值范围：
	// <li>1：申请中</li>
	// <li>2：运行中</li>
	// <li>3：受限运行中 (主备切换中)</li>
	// <li>4：已隔离</li>
	// <li>5：回收中</li>
	// <li>6：已回收</li>
	// <li>7：任务执行中 (实例做备份、回档等操作)</li>
	// <li>8：已下线</li>
	// <li>9：实例扩容中</li>
	// <li>10：实例迁移中</li>
	// <li>11：只读</li>
	// <li>12：重启中</li>
	Status *int64 `json:"Status" name:"Status"`

	// 偏移量，默认为 0
	Offset *int64 `json:"Offset" name:"Offset"`

	// 返回数量，默认为50
	Limit *int64 `json:"Limit" name:"Limit"`

	// 一个或者多个实例ID。实例ID，格式如：mssql-si2823jyl
	InstanceIdSet []*string `json:"InstanceIdSet" name:"InstanceIdSet" list`
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

		// 符合条件的实例总数。分页返回的话，这个值指的是所有符合条件的实例的个数，而非当前根据Limit和Offset值返回的实例个数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 实例列表
		DBInstances []*DBInstance `json:"DBInstances" name:"DBInstances" list`

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

type DescribeDBsRequest struct {
	*tchttp.BaseRequest

	// 实例ID
	InstanceIdSet []*string `json:"InstanceIdSet" name:"InstanceIdSet" list`

	// 每页记录数，最大为100，默认20
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 页编号，从第0页开始
	Offset *uint64 `json:"Offset" name:"Offset"`
}

func (r *DescribeDBsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDBsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 数据库数量
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 实例数据库列表
		DBInstances []*InstanceDBDetail `json:"DBInstances" name:"DBInstances" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDBsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDBsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeFlowStatusRequest struct {
	*tchttp.BaseRequest

	// 流程ID
	FlowId *int64 `json:"FlowId" name:"FlowId"`
}

func (r *DescribeFlowStatusRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeFlowStatusRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeFlowStatusResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 流程状态，0：成功，1：失败，2：运行中
		Status *int64 `json:"Status" name:"Status"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeFlowStatusResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeFlowStatusResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMigrationDetailRequest struct {
	*tchttp.BaseRequest

	// 迁移任务ID
	MigrateId *uint64 `json:"MigrateId" name:"MigrateId"`
}

func (r *DescribeMigrationDetailRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMigrationDetailRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMigrationDetailResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 迁移任务ID
		MigrateId *uint64 `json:"MigrateId" name:"MigrateId"`

		// 迁移任务名称
		MigrateName *string `json:"MigrateName" name:"MigrateName"`

		// 迁移任务所属的用户ID
		AppId *uint64 `json:"AppId" name:"AppId"`

		// 迁移任务所属的地域
		Region *string `json:"Region" name:"Region"`

		// 迁移源的类型 1:CDB for SQLServer 2:云服务器自建SQLServer数据库 4:SQLServer备份还原 5:SQLServer备份还原（COS方式）
		SourceType *int64 `json:"SourceType" name:"SourceType"`

		// 迁移任务的创建时间
		CreateTime *string `json:"CreateTime" name:"CreateTime"`

		// 迁移任务的开始时间
		StartTime *string `json:"StartTime" name:"StartTime"`

		// 迁移任务的结束时间
		EndTime *string `json:"EndTime" name:"EndTime"`

		// 迁移任务的状态（1:初始化,4:迁移中,5.迁移失败,6.迁移成功）
		Status *uint64 `json:"Status" name:"Status"`

		// 迁移任务当前进度
		Progress *int64 `json:"Progress" name:"Progress"`

		// 迁移类型（1:结构迁移 2:数据迁移 3:增量同步）
		MigrateType *int64 `json:"MigrateType" name:"MigrateType"`

		// 迁移源
		Source *MigrateSource `json:"Source" name:"Source"`

		// 迁移目标
		Target *MigrateTarget `json:"Target" name:"Target"`

		// 迁移DB对象 ，离线迁移（SourceType=4或SourceType=5）不使用。
		MigrateDBSet []*MigrateDB `json:"MigrateDBSet" name:"MigrateDBSet" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeMigrationDetailResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMigrationDetailResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMigrationsRequest struct {
	*tchttp.BaseRequest

	// 状态集合。只要符合集合中某一状态的迁移任务，就会查出来
	StatusSet []*int64 `json:"StatusSet" name:"StatusSet" list`

	// 迁移任务的名称，模糊匹配
	MigrateName *string `json:"MigrateName" name:"MigrateName"`

	// 每页的记录数
	Limit *int64 `json:"Limit" name:"Limit"`

	// 查询第几页的记录
	Offset *int64 `json:"Offset" name:"Offset"`

	// 查询结果按照关键字排序，可选值为name、createTime、startTime，endTime，status
	OrderBy *string `json:"OrderBy" name:"OrderBy"`

	// 排序方式，可选值为desc、asc
	OrderByType *string `json:"OrderByType" name:"OrderByType"`
}

func (r *DescribeMigrationsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMigrationsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMigrationsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 查询结果的总数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 查询结果的列表
		MigrateTaskSet []*MigrateTask `json:"MigrateTaskSet" name:"MigrateTaskSet" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeMigrationsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMigrationsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeOrdersRequest struct {
	*tchttp.BaseRequest

	// 订单数组。发货时会返回订单名字，利用该订单名字调用DescribeOrders接口查询发货情况
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

		// 订单信息数组
		Deals []*DealInfo `json:"Deals" name:"Deals" list`

		// 返回多少个订单的信息
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

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

	// 可用区英文ID，形如ap-guangzhou-1
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

		// 规格信息数组
		SpecInfoList []*SpecInfo `json:"SpecInfoList" name:"SpecInfoList" list`

		// 返回总共多少条数据
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

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

		// 返回地域信息总的条目
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 地域信息数组
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

type DescribeRollbackTimeRequest struct {
	*tchttp.BaseRequest

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 需要查询的数据库列表
	DBs []*string `json:"DBs" name:"DBs" list`
}

func (r *DescribeRollbackTimeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeRollbackTimeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeRollbackTimeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 数据库可回档实例信息
		Details []*DbRollbackTimeInfo `json:"Details" name:"Details" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeRollbackTimeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeRollbackTimeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSlowlogsRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如mssql-k8voqdlz
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 查询开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 查询结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 分页返回结果，分页大小，默认20，不超过100
	Limit *int64 `json:"Limit" name:"Limit"`

	// 从第几页开始返回，起始页，从0开始，默认为0
	Offset *uint64 `json:"Offset" name:"Offset"`
}

func (r *DescribeSlowlogsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSlowlogsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSlowlogsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 查询总数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 慢查询日志信息列表
		Slowlogs []*SlowlogInfo `json:"Slowlogs" name:"Slowlogs" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeSlowlogsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSlowlogsResponse) FromJsonString(s string) error {
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

		// 返回多少个可用区信息
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 可用区数组
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

type InquiryPriceCreateDBInstancesRequest struct {
	*tchttp.BaseRequest

	// 可用区ID。该参数可以通过调用 DescribeZones 接口的返回值中的Zone字段来获取。
	Zone *string `json:"Zone" name:"Zone"`

	// 内存大小，单位：GB
	Memory *int64 `json:"Memory" name:"Memory"`

	// 实例容量大小，单位：GB。
	Storage *int64 `json:"Storage" name:"Storage"`

	// 计费类型，当前只支持预付费，即包年包月，取值为PREPAID。默认值为PREPAID
	InstanceChargeType *string `json:"InstanceChargeType" name:"InstanceChargeType"`

	// 购买时长，单位：月。取值为1到48，默认为1
	Period *int64 `json:"Period" name:"Period"`

	// 一次性购买的实例数量。取值1-100，默认取值为1
	GoodsNum *int64 `json:"GoodsNum" name:"GoodsNum"`

	// sqlserver版本，目前只支持：2008R2（SQL Server 2008 R2），2012SP3（SQL Server 2012），2016SP1（SQL Server 2016 SP1）两种版本。默认为2008R2版本
	DBVersion *string `json:"DBVersion" name:"DBVersion"`
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

		// 未打折前价格，其值除以100表示多少钱。比如10010表示100.10元
		OriginalPrice *int64 `json:"OriginalPrice" name:"OriginalPrice"`

		// 实际需要支付的价格，其值除以100表示多少钱。比如10010表示100.10元
		Price *int64 `json:"Price" name:"Price"`

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
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 续费周期。按月续费最多48个月。默认查询续费一个月的价格
	Period *uint64 `json:"Period" name:"Period"`

	// 续费周期单位。month表示按月续费，当前只支持按月付费查询价格
	TimeUnit *string `json:"TimeUnit" name:"TimeUnit"`
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

		// 未打折的原价，其值除以100表示最终的价格。比如10094表示100.94元
		OriginalPrice *uint64 `json:"OriginalPrice" name:"OriginalPrice"`

		// 实际需要支付价格，其值除以100表示最终的价格。比如10094表示100.94元
		Price *uint64 `json:"Price" name:"Price"`

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

	// 实例ID，形如mssql-njj2mtpl
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 实例升级后的内存大小，单位GB，其值不能比当前实例内存小
	Memory *int64 `json:"Memory" name:"Memory"`

	// 实例升级后的磁盘大小，单位GB，其值不能比当前实例磁盘小
	Storage *int64 `json:"Storage" name:"Storage"`
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

		// 未打折的原价，其值除以100表示最终的价格。比如10094表示100.94元
		OriginalPrice *int64 `json:"OriginalPrice" name:"OriginalPrice"`

		// 实际需要支付价格，其值除以100表示最终的价格。比如10094表示100.94元
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

type InstanceDBDetail struct {

	// 实例id
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 数据库信息列表
	DBDetails []*DBDetail `json:"DBDetails" name:"DBDetails" list`
}

type InstanceRenewInfo struct {

	// 实例ID，形如mssql-j8kv137v
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 实例续费标记。0：正常续费，1：自动续费，2：到期不续
	RenewFlag *int64 `json:"RenewFlag" name:"RenewFlag"`
}

type MigrateDB struct {

	// 迁移数据库的名称
	DBName *string `json:"DBName" name:"DBName"`
}

type MigrateDetail struct {

	// 当前环节的名称
	StepName *string `json:"StepName" name:"StepName"`

	// 当前环节的进度（单位是%）
	Progress *int64 `json:"Progress" name:"Progress"`
}

type MigrateSource struct {

	// 迁移源实例的ID，MigrateType=1(CDB for SQLServers)时使用，格式如：mssql-si2823jyl
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 迁移源Cvm的ID，MigrateType=2(云服务器自建SQLServer数据库)时使用
	CvmId *string `json:"CvmId" name:"CvmId"`

	// 迁移源Cvm的Vpc网络标识，MigrateType=2(云服务器自建SQLServer数据库)时使用，格式如：vpc-6ys9ont9
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 迁移源Cvm的Vpc下的子网标识，MigrateType=2(云服务器自建SQLServer数据库)时使用，格式如：subnet-h9extioi
	SubnetId *string `json:"SubnetId" name:"SubnetId"`

	// 用户名，MigrateType=1或MigrateType=2使用
	UserName *string `json:"UserName" name:"UserName"`

	// 密码，MigrateType=1或MigrateType=2使用
	Password *string `json:"Password" name:"Password"`

	// 迁移源Cvm自建库的内网IP，MigrateType=2(云服务器自建SQLServer数据库)时使用
	Ip *string `json:"Ip" name:"Ip"`

	// 迁移源Cvm自建库的端口号，MigrateType=2(云服务器自建SQLServer数据库)时使用
	Port *uint64 `json:"Port" name:"Port"`

	// 离线迁移的源备份地址，MigrateType=4或MigrateType=5使用
	Url []*string `json:"Url" name:"Url" list`

	// 离线迁移的源备份密码，MigrateType=4或MigrateType=5使用
	UrlPassword *string `json:"UrlPassword" name:"UrlPassword"`
}

type MigrateTarget struct {

	// 迁移目标实例的ID，格式如：mssql-si2823jyl
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 迁移目标实例的用户名
	UserName *string `json:"UserName" name:"UserName"`

	// 迁移目标实例的密码
	Password *string `json:"Password" name:"Password"`
}

type MigrateTask struct {

	// 迁移任务ID
	MigrateId *uint64 `json:"MigrateId" name:"MigrateId"`

	// 迁移任务名称
	MigrateName *string `json:"MigrateName" name:"MigrateName"`

	// 迁移任务所属的用户ID
	AppId *uint64 `json:"AppId" name:"AppId"`

	// 迁移任务所属的地域
	Region *string `json:"Region" name:"Region"`

	// 迁移源的类型 1:CDB for SQLServer 2:云服务器自建SQLServer数据库 4:SQLServer备份还原 5:SQLServer备份还原（COS方式）
	SourceType *int64 `json:"SourceType" name:"SourceType"`

	// 迁移任务的创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 迁移任务的开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 迁移任务的结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 迁移任务的状态（1:初始化,4:迁移中,5.迁移失败,6.迁移成功）
	Status *uint64 `json:"Status" name:"Status"`

	// 信息
	Message *string `json:"Message" name:"Message"`

	// 是否迁移任务经过检查（0:未校验,1:校验成功,2:校验失败,3:校验中）
	CheckFlag *uint64 `json:"CheckFlag" name:"CheckFlag"`

	// 迁移任务当前进度（单位%）
	Progress *int64 `json:"Progress" name:"Progress"`

	// 迁移任务进度细节
	MigrateDetail *MigrateDetail `json:"MigrateDetail" name:"MigrateDetail"`
}

type ModifyAccountPrivilegeRequest struct {
	*tchttp.BaseRequest

	// 数据库实例ID，形如mssql-njj2mtpl
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 账号权限变更信息
	Accounts []*AccountPrivilegeModifyInfo `json:"Accounts" name:"Accounts" list`
}

func (r *ModifyAccountPrivilegeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAccountPrivilegeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAccountPrivilegeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务流程ID
		FlowId *uint64 `json:"FlowId" name:"FlowId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyAccountPrivilegeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAccountPrivilegeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAccountRemarkRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如mssql-j8kv137v
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 修改备注的账户信息
	Accounts []*AccountRemark `json:"Accounts" name:"Accounts" list`
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

	// 数据库实例ID，形如mssql-njj2mtpl
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

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

type ModifyDBInstanceProjectRequest struct {
	*tchttp.BaseRequest

	// 实例ID数组，形如mssql-j8kv137v
	InstanceIdSet []*string `json:"InstanceIdSet" name:"InstanceIdSet" list`

	// 项目ID，为0的话表示默认项目
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`
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

		// 修改成功的实例个数
		Count *int64 `json:"Count" name:"Count"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
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

type ModifyDBInstanceRenewFlagRequest struct {
	*tchttp.BaseRequest

	// 实例续费状态标记信息
	RenewFlags []*InstanceRenewInfo `json:"RenewFlags" name:"RenewFlags" list`
}

func (r *ModifyDBInstanceRenewFlagRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBInstanceRenewFlagRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBInstanceRenewFlagResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 修改成功的个数
		Count *int64 `json:"Count" name:"Count"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDBInstanceRenewFlagResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBInstanceRenewFlagResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBNameRequest struct {
	*tchttp.BaseRequest

	// 实例id
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 旧数据库名
	OldDBName *string `json:"OldDBName" name:"OldDBName"`

	// 新数据库名
	NewDBName *string `json:"NewDBName" name:"NewDBName"`
}

func (r *ModifyDBNameRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBNameRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBNameResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务流id
		FlowId *int64 `json:"FlowId" name:"FlowId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDBNameResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBNameResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBRemarkRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如mssql-rljoi3bf
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 数据库名称及备注数组，每个元素包含数据库名和对应的备注
	DBRemarks []*DBRemark `json:"DBRemarks" name:"DBRemarks" list`
}

func (r *ModifyDBRemarkRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBRemarkRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDBRemarkResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDBRemarkResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDBRemarkResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyMigrationRequest struct {
	*tchttp.BaseRequest

	// 迁移任务ID
	MigrateId *uint64 `json:"MigrateId" name:"MigrateId"`

	// 新的迁移任务的名称，若不填则不修改
	MigrateName *string `json:"MigrateName" name:"MigrateName"`

	// 新的迁移类型（1:结构迁移 2:数据迁移 3:增量同步），若不填则不修改
	MigrateType *uint64 `json:"MigrateType" name:"MigrateType"`

	// 迁移源的类型 1:CDB for SQLServer 2:云服务器自建SQLServer数据库 4:SQLServer备份还原 5:SQLServer备份还原（COS方式），若不填则不修改
	SourceType *uint64 `json:"SourceType" name:"SourceType"`

	// 迁移源，若不填则不修改
	Source *MigrateSource `json:"Source" name:"Source"`

	// 迁移目标，若不填则不修改
	Target *MigrateTarget `json:"Target" name:"Target"`

	// 迁移DB对象 ，离线迁移（SourceType=4或SourceType=5）不使用，若不填则不修改
	MigrateDBSet []*MigrateDB `json:"MigrateDBSet" name:"MigrateDBSet" list`
}

func (r *ModifyMigrationRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyMigrationRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyMigrationResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 迁移任务ID
		MigrateId *uint64 `json:"MigrateId" name:"MigrateId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyMigrationResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyMigrationResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RegionInfo struct {

	// 地域英文ID，类似ap-guanghou
	Region *string `json:"Region" name:"Region"`

	// 地域中文名称
	RegionName *string `json:"RegionName" name:"RegionName"`

	// 地域数字ID
	RegionId *int64 `json:"RegionId" name:"RegionId"`

	// 该地域目前是否可以售卖，UNAVAILABLE-不可售卖；AVAILABLE-可售卖
	RegionState *string `json:"RegionState" name:"RegionState"`
}

type RenewDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如mssql-j8kv137v
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 续费多少个月，取值范围为1-48，默认为1
	Period *uint64 `json:"Period" name:"Period"`

	// 是否自动使用代金券，0-不使用；1-使用；默认不实用
	AutoVoucher *int64 `json:"AutoVoucher" name:"AutoVoucher"`

	// 代金券ID数组，目前只支持使用1张代金券
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

		// 订单名称
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

	// 数据库实例ID，形如mssql-njj2mtpl
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 更新后的账户密码信息数组
	Accounts []*AccountPassword `json:"Accounts" name:"Accounts" list`
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

		// 修改帐号密码的异步任务流程ID
		FlowId *int64 `json:"FlowId" name:"FlowId"`

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

	// 数据库实例ID，形如mssql-njj2mtpl
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
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

		// 异步任务流程ID
		FlowId *uint64 `json:"FlowId" name:"FlowId"`

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

type RestoreInstanceRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如mssql-j8kv137v
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 备份文件ID，该ID可以通过DescribeBackups接口返回数据中的Id字段获得
	BackupId *int64 `json:"BackupId" name:"BackupId"`
}

func (r *RestoreInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RestoreInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RestoreInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步流程任务ID，使用FlowId调用DescribeFlowStatus接口获取任务执行状态
		FlowId *int64 `json:"FlowId" name:"FlowId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RestoreInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RestoreInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RollbackInstanceRequest struct {
	*tchttp.BaseRequest

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 回档类型，0-回档的数据库覆盖原库；1-回档的数据库以重命名的形式生成，不覆盖原库
	Type *uint64 `json:"Type" name:"Type"`

	// 需要回档的数据库
	DBs []*string `json:"DBs" name:"DBs" list`

	// 回档目标时间点
	Time *string `json:"Time" name:"Time"`
}

func (r *RollbackInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RollbackInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RollbackInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 异步任务ID
		FlowId *uint64 `json:"FlowId" name:"FlowId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RollbackInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RollbackInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RunMigrationRequest struct {
	*tchttp.BaseRequest

	// 迁移任务ID
	MigrateId *uint64 `json:"MigrateId" name:"MigrateId"`
}

func (r *RunMigrationRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RunMigrationRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RunMigrationResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 迁移流程启动后，返回流程ID
		FlowId *int64 `json:"FlowId" name:"FlowId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RunMigrationResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RunMigrationResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SlowlogInfo struct {

	// 慢查询日志文件唯一标识
	Id *int64 `json:"Id" name:"Id"`

	// 文件生成的开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 文件生成的结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 文件大小（KB）
	Size *int64 `json:"Size" name:"Size"`

	// 文件中log条数
	Count *int64 `json:"Count" name:"Count"`

	// 内网下载地址
	InternalAddr *string `json:"InternalAddr" name:"InternalAddr"`

	// 外网下载地址
	ExternalAddr *string `json:"ExternalAddr" name:"ExternalAddr"`
}

type SpecInfo struct {

	// 实例规格ID，利用DescribeZones返回的SpecId，结合DescribeProductConfig返回的可售卖规格信息，可获悉某个可用区下可购买什么规格的实例
	SpecId *int64 `json:"SpecId" name:"SpecId"`

	// 机型ID
	MachineType *string `json:"MachineType" name:"MachineType"`

	// 机型中文名称
	MachineTypeName *string `json:"MachineTypeName" name:"MachineTypeName"`

	// 数据库版本信息。取值为2008R2（表示SQL Server 2008 R2），2012SP3（表示SQL Server 2012），2016SP1（表示SQL Server 2016 SP1）
	Version *string `json:"Version" name:"Version"`

	// Version字段对应的版本名称
	VersionName *string `json:"VersionName" name:"VersionName"`

	// 内存大小，单位GB
	Memory *int64 `json:"Memory" name:"Memory"`

	// CPU核数
	CPU *int64 `json:"CPU" name:"CPU"`

	// 此规格下最小的磁盘大小，单位GB
	MinStorage *int64 `json:"MinStorage" name:"MinStorage"`

	// 此规格下最大的磁盘大小，单位GB
	MaxStorage *int64 `json:"MaxStorage" name:"MaxStorage"`

	// 此规格对应的QPS大小
	QPS *int64 `json:"QPS" name:"QPS"`

	// 此规格的中文描述信息
	SuitInfo *string `json:"SuitInfo" name:"SuitInfo"`

	// 此规格对应的Pid
	Pid *int64 `json:"Pid" name:"Pid"`
}

type UpgradeDBInstanceRequest struct {
	*tchttp.BaseRequest

	// 实例ID，形如mssql-j8kv137v
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 实例升级后内存大小，单位GB，其值不能小于当前实例内存大小
	Memory *int64 `json:"Memory" name:"Memory"`

	// 实例升级后磁盘大小，单位GB，其值不能小于当前实例磁盘大小
	Storage *int64 `json:"Storage" name:"Storage"`

	// 是否自动使用代金券，0 - 不使用；1 - 默认使用。取值默认为0
	AutoVoucher *int64 `json:"AutoVoucher" name:"AutoVoucher"`

	// 代金券ID，目前单个订单只能使用一张代金券
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

		// 订单名称
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

type ZoneInfo struct {

	// 可用区英文ID，形如ap-guangzhou-1，表示广州一区
	Zone *string `json:"Zone" name:"Zone"`

	// 可用区中文名称
	ZoneName *string `json:"ZoneName" name:"ZoneName"`

	// 可用区数字ID
	ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

	// 该可用区目前可售卖的规格ID，利用SpecId，结合接口DescribeProductConfig返回的数据，可获悉该可用区目前可售卖的规格大小
	SpecId *int64 `json:"SpecId" name:"SpecId"`

	// 当前可用区与规格下，可售卖的数据库版本，形如2008R2（表示SQL Server 2008 R2）。其可选值有2008R2（表示SQL Server 2008 R2），2012SP3（表示SQL Server 2012），2016SP1（表示SQL Server 2016 SP1）
	Version *string `json:"Version" name:"Version"`
}
