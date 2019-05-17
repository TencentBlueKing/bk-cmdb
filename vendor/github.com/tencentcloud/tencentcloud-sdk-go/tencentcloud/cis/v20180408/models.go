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

package v20180408

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type Container struct {

	// 容器启动命令
	Command *string `json:"Command" name:"Command"`

	// 容器启动参数
	Args []*string `json:"Args" name:"Args" list`

	// 容器环境变量
	EnvironmentVars []*EnvironmentVar `json:"EnvironmentVars" name:"EnvironmentVars" list`

	// 镜像
	Image *string `json:"Image" name:"Image"`

	// 容器名，由小写字母、数字和 - 组成，由小写字母开头，小写字母或数字结尾，且长度不超过 63个字符
	Name *string `json:"Name" name:"Name"`

	// CPU，单位：核
	Cpu *float64 `json:"Cpu" name:"Cpu"`

	// 内存，单位：Gi
	Memory *float64 `json:"Memory" name:"Memory"`

	// 重启次数
	RestartCount *uint64 `json:"RestartCount" name:"RestartCount"`

	// 当前状态
	CurrentState *ContainerState `json:"CurrentState" name:"CurrentState"`

	// 上一次状态
	PreviousState *ContainerState `json:"PreviousState" name:"PreviousState"`

	// 容器工作目录
	WorkingDir *string `json:"WorkingDir" name:"WorkingDir"`

	// 容器ID
	ContainerId *string `json:"ContainerId" name:"ContainerId"`
}

type ContainerInstance struct {

	// 容器实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 容器实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// 容器实例所属VpcId
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 容器实例所属SubnetId
	SubnetId *string `json:"SubnetId" name:"SubnetId"`

	// 容器实例状态
	State *string `json:"State" name:"State"`

	// 容器列表
	Containers []*Container `json:"Containers" name:"Containers" list`

	// 重启策略
	RestartPolicy *string `json:"RestartPolicy" name:"RestartPolicy"`

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 启动时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 可用区
	Zone *string `json:"Zone" name:"Zone"`

	// Vpc名称
	VpcName *string `json:"VpcName" name:"VpcName"`

	// VpcCidr
	VpcCidr *string `json:"VpcCidr" name:"VpcCidr"`

	// SubnetName
	SubnetName *string `json:"SubnetName" name:"SubnetName"`

	// 子网Cidr
	SubnetCidr *string `json:"SubnetCidr" name:"SubnetCidr"`

	// 内网IP
	LanIp *string `json:"LanIp" name:"LanIp"`
}

type ContainerLog struct {

	// 容器名称
	Name *string `json:"Name" name:"Name"`

	// 日志
	Log *string `json:"Log" name:"Log"`

	// 日志记录时间
	Time *string `json:"Time" name:"Time"`
}

type ContainerState struct {

	// 容器运行开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 容器状态
	State *string `json:"State" name:"State"`

	// 状态详情
	Reason *string `json:"Reason" name:"Reason"`

	// 容器运行结束时间
	FinishTime *string `json:"FinishTime" name:"FinishTime"`

	// 容器运行退出码
	ExitCode *int64 `json:"ExitCode" name:"ExitCode"`
}

type CreateContainerInstanceRequest struct {
	*tchttp.BaseRequest

	// 可用区
	Zone *string `json:"Zone" name:"Zone"`

	// vpcId
	VpcId *string `json:"VpcId" name:"VpcId"`

	// subnetId
	SubnetId *string `json:"SubnetId" name:"SubnetId"`

	// 容器实例名称，由小写字母、数字和 - 组成，由小写字母开头，小写字母或数字结尾，且长度不超过 40个字符
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// 重启策略（Always,OnFailure,Never）
	RestartPolicy *string `json:"RestartPolicy" name:"RestartPolicy"`

	// 容器列表
	Containers []*Container `json:"Containers" name:"Containers" list`
}

func (r *CreateContainerInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateContainerInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateContainerInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 容器实例ID
		InstanceId *string `json:"InstanceId" name:"InstanceId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateContainerInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateContainerInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteContainerInstanceRequest struct {
	*tchttp.BaseRequest

	// 容器实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`
}

func (r *DeleteContainerInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteContainerInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteContainerInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 操作信息
		Msg *string `json:"Msg" name:"Msg"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteContainerInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteContainerInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeContainerInstanceEventsRequest struct {
	*tchttp.BaseRequest

	// 容器实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`
}

func (r *DescribeContainerInstanceEventsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeContainerInstanceEventsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeContainerInstanceEventsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 容器实例事件列表
		EventList []*Event `json:"EventList" name:"EventList" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeContainerInstanceEventsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeContainerInstanceEventsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeContainerInstanceRequest struct {
	*tchttp.BaseRequest

	// 容器实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`
}

func (r *DescribeContainerInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeContainerInstanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeContainerInstanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 容器实例详细信息
		ContainerInstance *ContainerInstance `json:"ContainerInstance" name:"ContainerInstance"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeContainerInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeContainerInstanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeContainerInstancesRequest struct {
	*tchttp.BaseRequest

	// 偏移量，默认为0
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 返回数量，默认为10
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 过滤条件。
	// - Zone - String - 是否必填：否 -（过滤条件）按照可用区过滤。
	// - VpcId - String - 是否必填：否 -（过滤条件）按照VpcId过滤。
	// - InstanceName - String - 是否必填：否 -（过滤条件）按照容器实例名称做模糊查询。
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeContainerInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeContainerInstancesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeContainerInstancesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 容器实例列表
		ContainerInstanceList []*ContainerInstance `json:"ContainerInstanceList" name:"ContainerInstanceList" list`

		// 容器实例总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeContainerInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeContainerInstancesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeContainerLogRequest struct {
	*tchttp.BaseRequest

	// 容器实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// 容器名称
	ContainerName *string `json:"ContainerName" name:"ContainerName"`

	// 日志显示尾部行数
	Tail *uint64 `json:"Tail" name:"Tail"`

	// 日志起始时间
	SinceTime *string `json:"SinceTime" name:"SinceTime"`
}

func (r *DescribeContainerLogRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeContainerLogRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeContainerLogResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 容器日志数组
		ContainerLogList []*ContainerLog `json:"ContainerLogList" name:"ContainerLogList" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeContainerLogResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeContainerLogResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type EnvironmentVar struct {

	// 环境变量名
	Name *string `json:"Name" name:"Name"`

	// 环境变量值
	Value *string `json:"Value" name:"Value"`
}

type Event struct {

	// 事件首次出现时间
	FirstSeen *string `json:"FirstSeen" name:"FirstSeen"`

	// 事件上次出现时间
	LastSeen *string `json:"LastSeen" name:"LastSeen"`

	// 事件等级
	Level *string `json:"Level" name:"Level"`

	// 事件出现次数
	Count *string `json:"Count" name:"Count"`

	// 事件出现原因
	Reason *string `json:"Reason" name:"Reason"`

	// 事件消息
	Message *string `json:"Message" name:"Message"`
}

type Filter struct {

	// 过滤字段，可选值 - Zone，VpcId，InstanceName
	Name *string `json:"Name" name:"Name"`

	// 过滤值列表
	ValueList []*string `json:"ValueList" name:"ValueList" list`
}

type InquiryPriceCreateCisRequest struct {
	*tchttp.BaseRequest

	// 可用区
	Zone *string `json:"Zone" name:"Zone"`

	// CPU，单位：核
	Cpu *float64 `json:"Cpu" name:"Cpu"`

	// 内存，单位：Gi
	Memory *float64 `json:"Memory" name:"Memory"`
}

func (r *InquiryPriceCreateCisRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceCreateCisRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InquiryPriceCreateCisResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 价格
		Price *Price `json:"Price" name:"Price"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *InquiryPriceCreateCisResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceCreateCisResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Price struct {

	// 原价，单位：元
	DiscountPrice *float64 `json:"DiscountPrice" name:"DiscountPrice"`

	// 折扣价，单位：元
	OriginalPrice *float64 `json:"OriginalPrice" name:"OriginalPrice"`
}
