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

package v20180614

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type Attribute struct {

	// 属性列表
	Tags []*DeviceTag `json:"Tags" name:"Tags" list`
}

type BatchPublishMessage struct {

	// 消息发往的主题。为 Topic 权限中去除 ProductID 和 DeviceName 的部分，如 “event”
	Topic *string `json:"Topic" name:"Topic"`

	// 消息内容
	Payload *string `json:"Payload" name:"Payload"`
}

type BatchUpdateShadow struct {

	// 设备影子的期望状态，格式为 Json 对象序列化之后的字符串
	Desired *string `json:"Desired" name:"Desired"`
}

type CancelTaskRequest struct {
	*tchttp.BaseRequest

	// 任务 ID
	Id *string `json:"Id" name:"Id"`
}

func (r *CancelTaskRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CancelTaskRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CancelTaskResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CancelTaskResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CancelTaskResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDeviceRequest struct {
	*tchttp.BaseRequest

	// 产品 ID 。创建产品时腾讯云为用户分配全局唯一的 ID
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称。命名规则：[a-zA-Z0-9:_-]{1,48}。
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 设备属性
	Attribute *Attribute `json:"Attribute" name:"Attribute"`

	// 是否使用自定义PSK，默认不使用
	DefinedPsk *string `json:"DefinedPsk" name:"DefinedPsk"`

	// 运营商类型，当产品是NB-IoT产品时，此字段必填。1表示中国电信，2表示中国移动，3表示中国联通
	Isp *uint64 `json:"Isp" name:"Isp"`

	// IMEI，当产品是NB-IoT产品时，此字段必填
	Imei *string `json:"Imei" name:"Imei"`

	// LoRa设备的DevEui，当创建LoRa时，此字段必填
	LoraDevEui *string `json:"LoraDevEui" name:"LoraDevEui"`
}

func (r *CreateDeviceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDeviceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDeviceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备名称
		DeviceName *string `json:"DeviceName" name:"DeviceName"`

		// 对称加密密钥，base64编码。采用对称加密时返回该参数
		DevicePsk *string `json:"DevicePsk" name:"DevicePsk"`

		// 设备证书，用于 TLS 建立链接时校验客户端身份。采用非对称加密时返回该参数
		DeviceCert *string `json:"DeviceCert" name:"DeviceCert"`

		// 设备私钥，用于 TLS 建立链接时校验客户端身份，腾讯云后台不保存，请妥善保管。采用非对称加密时返回该参数
		DevicePrivateKey *string `json:"DevicePrivateKey" name:"DevicePrivateKey"`

		// LoRa设备的DevEui，当设备是LoRa设备时，会返回该字段
		LoraDevEui *string `json:"LoraDevEui" name:"LoraDevEui"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateDeviceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDeviceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateMultiDeviceRequest struct {
	*tchttp.BaseRequest

	// 产品 ID。创建产品时腾讯云为用户分配全局唯一的 ID
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 批量创建的设备名数组，单次最多创建 100 个设备。命名规则：[a-zA-Z0-9:_-]{1,48}
	DeviceNames []*string `json:"DeviceNames" name:"DeviceNames" list`
}

func (r *CreateMultiDeviceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateMultiDeviceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateMultiDeviceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务ID，腾讯云生成全局唯一的任务 ID，有效期一个月，一个月之后任务失效。可以调用获取创建多设备任务状态接口获取该任务的执行状态，当状态为成功时，可以调用获取创建多设备任务结果接口获取该任务的结果
		TaskId *string `json:"TaskId" name:"TaskId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateMultiDeviceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateMultiDeviceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateProductRequest struct {
	*tchttp.BaseRequest

	// 产品名称，名称不能和已经存在的产品名称重复。命名规则：[a-zA-Z0-9:_-]{1,32}
	ProductName *string `json:"ProductName" name:"ProductName"`

	// 产品属性
	ProductProperties *ProductProperties `json:"ProductProperties" name:"ProductProperties"`
}

func (r *CreateProductRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateProductRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateProductResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 产品名称
		ProductName *string `json:"ProductName" name:"ProductName"`

		// 产品 ID，腾讯云生成全局唯一 ID
		ProductId *string `json:"ProductId" name:"ProductId"`

		// 产品属性
		ProductProperties *ProductProperties `json:"ProductProperties" name:"ProductProperties"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateProductResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateProductResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateTaskRequest struct {
	*tchttp.BaseRequest

	// 任务类型，取值为 “UpdateShadow” 或者 “PublishMessage”
	TaskType *string `json:"TaskType" name:"TaskType"`

	// 执行任务的产品ID
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 执行任务的设备名的正则表达式
	DeviceNameFilter *string `json:"DeviceNameFilter" name:"DeviceNameFilter"`

	// 任务开始执行的时间。 取值为 Unix 时间戳，单位秒，且需大于等于当前时间时间戳，0为系统当前时间时间戳，即立即执行，最大为当前时间86400秒后，超过则取值为当前时间86400秒后
	ScheduleTimeInSeconds *uint64 `json:"ScheduleTimeInSeconds" name:"ScheduleTimeInSeconds"`

	// 任务描述细节，描述见下 Task
	Tasks *Task `json:"Tasks" name:"Tasks"`

	// 最长执行时间，单位秒，被调度后超过此时间仍未有结果则视为任务失败。取值为0-86400，默认为86400
	MaxExecutionTimeInSeconds *uint64 `json:"MaxExecutionTimeInSeconds" name:"MaxExecutionTimeInSeconds"`
}

func (r *CreateTaskRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateTaskRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateTaskResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 创建的任务ID
		TaskId *string `json:"TaskId" name:"TaskId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateTaskResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateTaskResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteDeviceRequest struct {
	*tchttp.BaseRequest

	// 设备所属的产品 ID
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 需要删除的设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`
}

func (r *DeleteDeviceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteDeviceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteDeviceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteDeviceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteDeviceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteProductRequest struct {
	*tchttp.BaseRequest

	// 需要删除的产品 ID
	ProductId *string `json:"ProductId" name:"ProductId"`
}

func (r *DeleteProductRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteProductRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteProductResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteProductResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteProductResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDeviceShadowRequest struct {
	*tchttp.BaseRequest

	// 产品 ID
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称。命名规则：[a-zA-Z0-9:_-]{1,48}
	DeviceName *string `json:"DeviceName" name:"DeviceName"`
}

func (r *DescribeDeviceShadowRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDeviceShadowRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDeviceShadowResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备影子数据
		Data *string `json:"Data" name:"Data"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDeviceShadowResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDeviceShadowResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDevicesRequest struct {
	*tchttp.BaseRequest

	// 需要查看设备列表的产品 ID
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 分页偏移
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 分页的大小，数值范围 10-100
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 设备固件版本号，若不带此参数会返回所有固件版本的设备
	FirmwareVersion *string `json:"FirmwareVersion" name:"FirmwareVersion"`
}

func (r *DescribeDevicesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDevicesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDevicesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 设备详细信息列表
		Devices []*DeviceInfo `json:"Devices" name:"Devices" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDevicesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDevicesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMultiDevTaskRequest struct {
	*tchttp.BaseRequest

	// 任务 ID，由批量创建设备接口返回
	TaskId *string `json:"TaskId" name:"TaskId"`

	// 产品 ID，创建产品时腾讯云为用户分配全局唯一的 ID
	ProductId *string `json:"ProductId" name:"ProductId"`
}

func (r *DescribeMultiDevTaskRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMultiDevTaskRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMultiDevTaskResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务 ID
		TaskId *string `json:"TaskId" name:"TaskId"`

		// 任务是否完成。0 代表任务未开始，1 代表任务正在执行，2 代表任务已完成
		TaskStatus *uint64 `json:"TaskStatus" name:"TaskStatus"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeMultiDevTaskResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMultiDevTaskResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMultiDevicesRequest struct {
	*tchttp.BaseRequest

	// 产品 ID，创建产品时腾讯云为用户分配全局唯一的 ID
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 任务 ID，由批量创建设备接口返回
	TaskId *string `json:"TaskId" name:"TaskId"`

	// 分页偏移
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 分页大小，每页返回的设备个数
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribeMultiDevicesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMultiDevicesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeMultiDevicesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务 ID，由批量创建设备接口返回
		TaskId *string `json:"TaskId" name:"TaskId"`

		// 设备详细信息列表
		DevicesInfo []*MultiDevicesInfo `json:"DevicesInfo" name:"DevicesInfo" list`

		// 该任务创建设备的总数
		TotalDevNum *uint64 `json:"TotalDevNum" name:"TotalDevNum"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeMultiDevicesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeMultiDevicesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProductsRequest struct {
	*tchttp.BaseRequest

	// 分页偏移，Offset从0开始
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 分页大小，当前页面中显示的最大数量，值范围 10-250。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 过滤条件
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeProductsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProductsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeProductsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 产品总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 产品详细信息列表
		Products []*ProductInfo `json:"Products" name:"Products" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeProductsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeProductsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTaskRequest struct {
	*tchttp.BaseRequest

	// 任务ID
	Id *string `json:"Id" name:"Id"`
}

func (r *DescribeTaskRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTaskRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTaskResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务类型，目前取值为 “UpdateShadow” 或者 “PublishMessage”
		Type *string `json:"Type" name:"Type"`

		// 任务 ID
		Id *string `json:"Id" name:"Id"`

		// 产品 ID
		ProductId *string `json:"ProductId" name:"ProductId"`

		// 状态。1表示等待处理，2表示调度处理中，3表示已完成，4表示失败，5表示已取消
		Status *uint64 `json:"Status" name:"Status"`

		// 任务创建时间，Unix 时间戳
		CreateTime *uint64 `json:"CreateTime" name:"CreateTime"`

		// 最后任务更新时间，Unix 时间戳
		UpdateTime *uint64 `json:"UpdateTime" name:"UpdateTime"`

		// 任务完成时间，Unix 时间戳
		DoneTime *uint64 `json:"DoneTime" name:"DoneTime"`

		// 被调度时间，Unix 时间戳
		ScheduleTime *uint64 `json:"ScheduleTime" name:"ScheduleTime"`

		// 返回的错误码
		RetCode *uint64 `json:"RetCode" name:"RetCode"`

		// 返回的错误信息
		ErrMsg *string `json:"ErrMsg" name:"ErrMsg"`

		// 完成任务的设备比例
		Percent *uint64 `json:"Percent" name:"Percent"`

		// 匹配到的需执行任务的设备数目
		AllDeviceCnt *uint64 `json:"AllDeviceCnt" name:"AllDeviceCnt"`

		// 已完成任务的设备数目
		DoneDeviceCnt *uint64 `json:"DoneDeviceCnt" name:"DoneDeviceCnt"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeTaskResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTaskResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTasksRequest struct {
	*tchttp.BaseRequest

	// 分页偏移，从0开始
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 分页的大小，数值范围 1-250
	Limit *uint64 `json:"Limit" name:"Limit"`
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

		// 用户一个月内创建的任务总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 此页任务对象的数组，按创建时间排序
		Tasks []*TaskInfo `json:"Tasks" name:"Tasks" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
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

type DeviceInfo struct {

	// 设备名
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 设备是否在线，0不在线，1在线
	Online *uint64 `json:"Online" name:"Online"`

	// 设备登陆时间
	LoginTime *uint64 `json:"LoginTime" name:"LoginTime"`

	// 设备版本
	Version *string `json:"Version" name:"Version"`

	// 设备证书，证书加密的设备返回
	DeviceCert *string `json:"DeviceCert" name:"DeviceCert"`

	// 设备密钥，密钥加密的设备返回
	DevicePsk *string `json:"DevicePsk" name:"DevicePsk"`

	// 设备属性
	Tags []*DeviceTag `json:"Tags" name:"Tags" list`

	// 设备类型
	DeviceType *uint64 `json:"DeviceType" name:"DeviceType"`

	// IMEI
	Imei *string `json:"Imei" name:"Imei"`

	// 运营商类型
	Isp *uint64 `json:"Isp" name:"Isp"`

	// NB IOT运营商处的DeviceID
	NbiotDeviceID *string `json:"NbiotDeviceID" name:"NbiotDeviceID"`

	// IP地址
	ConnIP *uint64 `json:"ConnIP" name:"ConnIP"`

	// 设备最后更新时间
	LastUpdateTime *uint64 `json:"LastUpdateTime" name:"LastUpdateTime"`

	// LoRa设备的dev eui
	LoraDevEui *string `json:"LoraDevEui" name:"LoraDevEui"`

	// LoRa设备的Mote type
	LoraMoteType *uint64 `json:"LoraMoteType" name:"LoraMoteType"`
}

type DeviceTag struct {

	// 属性名称
	Tag *string `json:"Tag" name:"Tag"`

	// 属性值的类型，1 int，2 string
	Type *uint64 `json:"Type" name:"Type"`

	// 属性的值
	Value *string `json:"Value" name:"Value"`
}

type Filter struct {

	// 过滤键的名称
	Name *string `json:"Name" name:"Name"`

	// 一个或者多个过滤值
	Values []*string `json:"Values" name:"Values" list`
}

type MultiDevicesInfo struct {

	// 设备名
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 对称加密密钥，base64 编码，采用对称加密时返回该参数
	DevicePsk *string `json:"DevicePsk" name:"DevicePsk"`

	// 设备证书，采用非对称加密时返回该参数
	DeviceCert *string `json:"DeviceCert" name:"DeviceCert"`

	// 设备私钥，采用非对称加密时返回该参数，腾讯云为用户缓存起来，其生命周期与任务生命周期一致
	DevicePrivateKey *string `json:"DevicePrivateKey" name:"DevicePrivateKey"`

	// 错误码
	Result *uint64 `json:"Result" name:"Result"`

	// 错误信息
	ErrMsg *string `json:"ErrMsg" name:"ErrMsg"`
}

type ProductInfo struct {

	// 产品ID
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 产品名
	ProductName *string `json:"ProductName" name:"ProductName"`

	// 产品元数据
	ProductMetadata *ProductMetadata `json:"ProductMetadata" name:"ProductMetadata"`

	// 产品属性
	ProductProperties *ProductProperties `json:"ProductProperties" name:"ProductProperties"`
}

type ProductMetadata struct {

	// 产品创建时间
	CreationDate *uint64 `json:"CreationDate" name:"CreationDate"`
}

type ProductProperties struct {

	// 产品描述
	ProductDescription *string `json:"ProductDescription" name:"ProductDescription"`

	// 加密类型，1表示非对称加密，2表示对称加密。如不填写，默认值是1
	EncryptionType *string `json:"EncryptionType" name:"EncryptionType"`

	// 产品所属区域，目前只支持广州（gz）
	Region *string `json:"Region" name:"Region"`

	// 产品类型，0表示正常设备，2表示NB-IoT设备，默认值是0
	ProductType *uint64 `json:"ProductType" name:"ProductType"`

	// 数据格式，取值为json或者custom，默认值是json
	Format *string `json:"Format" name:"Format"`

	// 产品所属平台，默认值是0
	Platform *string `json:"Platform" name:"Platform"`

	// LoRa产品运营侧APPEUI，只有LoRa产品需要填写
	Appeui *string `json:"Appeui" name:"Appeui"`
}

type PublishMessageRequest struct {
	*tchttp.BaseRequest

	// 消息发往的主题。命名规则：${ProductId}/${DeviceName}/[a-zA-Z0-9:_-]{1,128}
	Topic *string `json:"Topic" name:"Topic"`

	// 消息内容
	Payload *string `json:"Payload" name:"Payload"`

	// 产品ID
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 服务质量等级，取值为0， 1
	Qos *uint64 `json:"Qos" name:"Qos"`
}

func (r *PublishMessageRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *PublishMessageRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type PublishMessageResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *PublishMessageResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *PublishMessageResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Task struct {

	// 批量更新影子任务的描述细节，当 taskType 取值为 “UpdateShadow” 时，此字段必填。描述见下 BatchUpdateShadow
	UpdateShadowTask *BatchUpdateShadow `json:"UpdateShadowTask" name:"UpdateShadowTask"`

	// 批量下发消息任务的描述细节，当 taskType 取值为 “PublishMessage” 时，此字段必填。描述见下 BatchPublishMessage
	PublishMessageTask *BatchPublishMessage `json:"PublishMessageTask" name:"PublishMessageTask"`
}

type TaskInfo struct {

	// 任务类型，目前取值为 “UpdateShadow” 或者 “PublishMessage”
	Type *string `json:"Type" name:"Type"`

	// 任务 ID
	Id *string `json:"Id" name:"Id"`

	// 产品 ID
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 状态。1表示等待处理，2表示调度处理中，3表示已完成，4表示失败，5表示已取消
	Status *uint64 `json:"Status" name:"Status"`

	// 任务创建时间，Unix 时间戳
	CreateTime *uint64 `json:"CreateTime" name:"CreateTime"`

	// 最后任务更新时间，Unix 时间戳
	UpdateTime *uint64 `json:"UpdateTime" name:"UpdateTime"`

	// 返回的错误码
	RetCode *uint64 `json:"RetCode" name:"RetCode"`

	// 返回的错误信息
	ErrMsg *string `json:"ErrMsg" name:"ErrMsg"`
}

type UpdateDeviceShadowRequest struct {
	*tchttp.BaseRequest

	// 产品ID
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 虚拟设备的状态，JSON字符串格式，由desired结构组成
	State *string `json:"State" name:"State"`

	// 当前版本号，需要和后台的version保持一致，才能更新成功
	ShadowVersion *uint64 `json:"ShadowVersion" name:"ShadowVersion"`
}

func (r *UpdateDeviceShadowRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateDeviceShadowRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateDeviceShadowResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备影子数据，JSON字符串格式
		Data *string `json:"Data" name:"Data"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UpdateDeviceShadowResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateDeviceShadowResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}
