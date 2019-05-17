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

package v20180319

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type DeregisterMigrationTaskRequest struct {
	*tchttp.BaseRequest

	// 任务ID
	TaskId *string `json:"TaskId" name:"TaskId"`
}

func (r *DeregisterMigrationTaskRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeregisterMigrationTaskRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeregisterMigrationTaskResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeregisterMigrationTaskResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeregisterMigrationTaskResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMigrationTaskRequest struct {
	*tchttp.BaseRequest

	// 任务ID
	TaskId *string `json:"TaskId" name:"TaskId"`
}

func (r *DescribeMigrationTaskRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMigrationTaskRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMigrationTaskResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 迁移详情列表
		TaskStatus []*TaskStatus `json:"TaskStatus" name:"TaskStatus" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeMigrationTaskResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMigrationTaskResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DstInfo struct {

	// 迁移目的地域
	Region *string `json:"Region" name:"Region"`

	// 迁移目的Ip
	Ip *string `json:"Ip" name:"Ip"`

	// 迁移目的端口
	Port *string `json:"Port" name:"Port"`

	// 迁移目的实例Id
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

type ListMigrationProjectRequest struct {
	*tchttp.BaseRequest

	// 记录起始数，默认值为0
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 返回条数，默认值为500
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *ListMigrationProjectRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ListMigrationProjectRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ListMigrationProjectResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 项目列表
		Projects []*Project `json:"Projects" name:"Projects" list`

		// 项目总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ListMigrationProjectResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ListMigrationProjectResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ListMigrationTaskRequest struct {
	*tchttp.BaseRequest

	// 记录起始数，默认值为0
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 记录条数，默认值为10
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 项目ID，默认值为空
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`
}

func (r *ListMigrationTaskRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ListMigrationTaskRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ListMigrationTaskResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 记录总条数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 迁移任务列表
		Tasks []*Task `json:"Tasks" name:"Tasks" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ListMigrationTaskResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ListMigrationTaskResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyMigrationTaskBelongToProjectRequest struct {
	*tchttp.BaseRequest

	// 任务ID
	TaskId *string `json:"TaskId" name:"TaskId"`

	// 项目ID
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`
}

func (r *ModifyMigrationTaskBelongToProjectRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyMigrationTaskBelongToProjectRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyMigrationTaskBelongToProjectResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyMigrationTaskBelongToProjectResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyMigrationTaskBelongToProjectResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyMigrationTaskStatusRequest struct {
	*tchttp.BaseRequest

	// 任务状态
	Status *string `json:"Status" name:"Status"`

	// 任务ID
	TaskId *string `json:"TaskId" name:"TaskId"`
}

func (r *ModifyMigrationTaskStatusRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyMigrationTaskStatusRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyMigrationTaskStatusResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyMigrationTaskStatusResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyMigrationTaskStatusResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Project struct {

	// 项目ID
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`

	// 项目名称
	ProjectName *string `json:"ProjectName" name:"ProjectName"`
}

type RegisterMigrationTaskRequest struct {
	*tchttp.BaseRequest

	// 任务类型，取值database（数据库迁移）、file（文件迁移）、host（主机迁移）
	TaskType *string `json:"TaskType" name:"TaskType"`

	// 任务名称
	TaskName *string `json:"TaskName" name:"TaskName"`

	// 服务提供商名称
	ServiceSupplier *string `json:"ServiceSupplier" name:"ServiceSupplier"`

	// 迁移任务源信息
	SrcInfo *SrcInfo `json:"SrcInfo" name:"SrcInfo"`

	// 迁移任务目的信息
	DstInfo *DstInfo `json:"DstInfo" name:"DstInfo"`

	// 迁移任务创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 迁移任务更新时间
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`

	// 迁移类别，如数据库迁移中mysql:mysql代表从mysql迁移到mysql，文件迁移中oss:cos代表从阿里云oss迁移到腾讯云cos
	MigrateClass *string `json:"MigrateClass" name:"MigrateClass"`

	// 源实例接入类型
	SrcAccessType *string `json:"SrcAccessType" name:"SrcAccessType"`

	// 源实例数据库类型
	SrcDatabaseType *string `json:"SrcDatabaseType" name:"SrcDatabaseType"`

	// 目标实例接入类型
	DstAccessType *string `json:"DstAccessType" name:"DstAccessType"`

	// 目标实例数据库类型
	DstDatabaseType *string `json:"DstDatabaseType" name:"DstDatabaseType"`
}

func (r *RegisterMigrationTaskRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RegisterMigrationTaskRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RegisterMigrationTaskResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务ID
		TaskId *string `json:"TaskId" name:"TaskId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RegisterMigrationTaskResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RegisterMigrationTaskResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SrcInfo struct {

	// 迁移源地域
	Region *string `json:"Region" name:"Region"`

	// 迁移源Ip
	Ip *string `json:"Ip" name:"Ip"`

	// 迁移源端口
	Port *string `json:"Port" name:"Port"`

	// 迁移源实例Id
	InstanceId *string `json:"InstanceId" name:"InstanceId"`
}

type Task struct {

	// 任务Id
	TaskId *string `json:"TaskId" name:"TaskId"`

	// 任务名称
	TaskName *string `json:"TaskName" name:"TaskName"`

	// 迁移类型
	MigrationType *string `json:"MigrationType" name:"MigrationType"`

	// 迁移状态
	Status *string `json:"Status" name:"Status"`

	// 项目Id
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`

	// 项目名称
	ProjectName *string `json:"ProjectName" name:"ProjectName"`

	// 迁移源信息
	SrcInfo *SrcInfo `json:"SrcInfo" name:"SrcInfo"`

	// 迁移时间信息
	MigrationTimeLine *TimeObj `json:"MigrationTimeLine" name:"MigrationTimeLine"`

	// 状态更新时间
	Updated *string `json:"Updated" name:"Updated"`

	// 迁移目的信息
	DstInfo *DstInfo `json:"DstInfo" name:"DstInfo"`
}

type TaskStatus struct {

	// 迁移状态
	Status *string `json:"Status" name:"Status"`

	// 迁移进度
	Progress *string `json:"Progress" name:"Progress"`

	// 迁移日期
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`
}

type TimeObj struct {

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`
}
