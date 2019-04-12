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

package v20180123

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type Action struct {

	// 转发至topic
	Topic *TopicAction `json:"Topic" name:"Topic"`

	// 转发至第三发
	Service *ServiceAction `json:"Service" name:"Service"`

	// 转发至第三发Ckafka
	Ckafka *CkafkaAction `json:"Ckafka" name:"Ckafka"`
}

type ActivateRuleRequest struct {
	*tchttp.BaseRequest

	// 规则Id
	RuleId *string `json:"RuleId" name:"RuleId"`
}

func (r *ActivateRuleRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ActivateRuleRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ActivateRuleResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ActivateRuleResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ActivateRuleResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AddDeviceRequest struct {
	*tchttp.BaseRequest

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称，唯一标识某产品下的一个设备
	DeviceName *string `json:"DeviceName" name:"DeviceName"`
}

func (r *AddDeviceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AddDeviceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AddDeviceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备信息
		Device *Device `json:"Device" name:"Device"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AddDeviceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AddDeviceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AddProductRequest struct {
	*tchttp.BaseRequest

	// 产品名称，同一区域产品名称需唯一，支持中文、英文字母、中划线和下划线，长度不超过31个字符，中文占两个字符
	Name *string `json:"Name" name:"Name"`

	// 产品描述
	Description *string `json:"Description" name:"Description"`

	// 数据模版
	DataTemplate []*DataTemplate `json:"DataTemplate" name:"DataTemplate" list`

	// 产品版本（native表示基础版，template表示高级版，默认值为template）
	DataProtocol *string `json:"DataProtocol" name:"DataProtocol"`

	// 设备认证方式（1：动态令牌，2：签名直连鉴权）
	AuthType *uint64 `json:"AuthType" name:"AuthType"`

	// 通信方式（other/wifi/cellular/nb-iot）
	CommProtocol *string `json:"CommProtocol" name:"CommProtocol"`

	// 产品的设备类型（device: 直连设备；sub_device：子设备；gateway：网关设备）
	DeviceType *string `json:"DeviceType" name:"DeviceType"`
}

func (r *AddProductRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AddProductRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AddProductResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 产品信息
		Product *Product `json:"Product" name:"Product"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AddProductResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AddProductResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AddRuleRequest struct {
	*tchttp.BaseRequest

	// 名称
	Name *string `json:"Name" name:"Name"`

	// 描述
	Description *string `json:"Description" name:"Description"`

	// 查询
	Query *RuleQuery `json:"Query" name:"Query"`

	// 转发动作列表
	Actions []*Action `json:"Actions" name:"Actions" list`

	// 数据类型（0：文本，1：二进制）
	DataType *uint64 `json:"DataType" name:"DataType"`
}

func (r *AddRuleRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AddRuleRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AddRuleResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 规则
		Rule *Rule `json:"Rule" name:"Rule"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AddRuleResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AddRuleResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AddTopicRequest struct {
	*tchttp.BaseRequest

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// Topic名称
	TopicName *string `json:"TopicName" name:"TopicName"`
}

func (r *AddTopicRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AddTopicRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AddTopicResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// Topic信息
		Topic *Topic `json:"Topic" name:"Topic"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AddTopicResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AddTopicResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppAddUserRequest struct {
	*tchttp.BaseRequest

	// 用户名
	UserName *string `json:"UserName" name:"UserName"`

	// 密码
	Password *string `json:"Password" name:"Password"`
}

func (r *AppAddUserRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppAddUserRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppAddUserResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 应用用户
		AppUser *AppUser `json:"AppUser" name:"AppUser"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AppAddUserResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppAddUserResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppDeleteDeviceRequest struct {
	*tchttp.BaseRequest

	// 访问Token
	AccessToken *string `json:"AccessToken" name:"AccessToken"`

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`
}

func (r *AppDeleteDeviceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppDeleteDeviceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppDeleteDeviceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AppDeleteDeviceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppDeleteDeviceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppDevice struct {

	// 设备Id
	DeviceId *string `json:"DeviceId" name:"DeviceId"`

	// 所属产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 别名
	AliasName *string `json:"AliasName" name:"AliasName"`

	// 地区
	Region *string `json:"Region" name:"Region"`

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 更新时间
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`
}

type AppDeviceDetail struct {

	// 设备Id
	DeviceId *string `json:"DeviceId" name:"DeviceId"`

	// 所属产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 别名
	AliasName *string `json:"AliasName" name:"AliasName"`

	// 地区
	Region *string `json:"Region" name:"Region"`

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 更新时间
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`

	// 设备信息（json）
	DeviceInfo *string `json:"DeviceInfo" name:"DeviceInfo"`

	// 数据模板
	DataTemplate []*DataTemplate `json:"DataTemplate" name:"DataTemplate" list`
}

type AppGetDeviceDataRequest struct {
	*tchttp.BaseRequest

	// 访问Token
	AccessToken *string `json:"AccessToken" name:"AccessToken"`

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`
}

func (r *AppGetDeviceDataRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppGetDeviceDataRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppGetDeviceDataResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备数据。
		DeviceData *string `json:"DeviceData" name:"DeviceData"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AppGetDeviceDataResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppGetDeviceDataResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppGetDeviceRequest struct {
	*tchttp.BaseRequest

	// 访问Token
	AccessToken *string `json:"AccessToken" name:"AccessToken"`

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`
}

func (r *AppGetDeviceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppGetDeviceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppGetDeviceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 绑定设备详情
		AppDevice *AppDeviceDetail `json:"AppDevice" name:"AppDevice"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AppGetDeviceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppGetDeviceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppGetDeviceStatusesRequest struct {
	*tchttp.BaseRequest

	// 访问Token
	AccessToken *string `json:"AccessToken" name:"AccessToken"`

	// 设备Id列表（单次限制1000个设备）
	DeviceIds []*string `json:"DeviceIds" name:"DeviceIds" list`
}

func (r *AppGetDeviceStatusesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppGetDeviceStatusesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppGetDeviceStatusesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备状态
		DeviceStatuses []*DeviceStatus `json:"DeviceStatuses" name:"DeviceStatuses" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AppGetDeviceStatusesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppGetDeviceStatusesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppGetDevicesRequest struct {
	*tchttp.BaseRequest

	// 访问Token
	AccessToken *string `json:"AccessToken" name:"AccessToken"`
}

func (r *AppGetDevicesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppGetDevicesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppGetDevicesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 绑定设备列表
		Devices []*AppDevice `json:"Devices" name:"Devices" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AppGetDevicesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppGetDevicesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppGetTokenRequest struct {
	*tchttp.BaseRequest

	// 用户名
	UserName *string `json:"UserName" name:"UserName"`

	// 密码
	Password *string `json:"Password" name:"Password"`

	// TTL
	Expire *uint64 `json:"Expire" name:"Expire"`
}

func (r *AppGetTokenRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppGetTokenRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppGetTokenResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 访问Token
		AccessToken *string `json:"AccessToken" name:"AccessToken"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AppGetTokenResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppGetTokenResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppGetUserRequest struct {
	*tchttp.BaseRequest

	// 访问Token
	AccessToken *string `json:"AccessToken" name:"AccessToken"`
}

func (r *AppGetUserRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppGetUserRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppGetUserResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 用户信息
		AppUser *AppUser `json:"AppUser" name:"AppUser"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AppGetUserResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppGetUserResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppIssueDeviceControlRequest struct {
	*tchttp.BaseRequest

	// 访问Token
	AccessToken *string `json:"AccessToken" name:"AccessToken"`

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 控制数据（json）
	ControlData *string `json:"ControlData" name:"ControlData"`

	// 是否发送metadata字段
	Metadata *bool `json:"Metadata" name:"Metadata"`
}

func (r *AppIssueDeviceControlRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppIssueDeviceControlRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppIssueDeviceControlResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AppIssueDeviceControlResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppIssueDeviceControlResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppResetPasswordRequest struct {
	*tchttp.BaseRequest

	// 访问Token
	AccessToken *string `json:"AccessToken" name:"AccessToken"`

	// 旧密码
	OldPassword *string `json:"OldPassword" name:"OldPassword"`

	// 新密码
	NewPassword *string `json:"NewPassword" name:"NewPassword"`
}

func (r *AppResetPasswordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppResetPasswordRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppResetPasswordResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AppResetPasswordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppResetPasswordResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppSecureAddDeviceRequest struct {
	*tchttp.BaseRequest

	// 访问Token
	AccessToken *string `json:"AccessToken" name:"AccessToken"`

	// 设备签名
	DeviceSignature *string `json:"DeviceSignature" name:"DeviceSignature"`
}

func (r *AppSecureAddDeviceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppSecureAddDeviceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppSecureAddDeviceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 绑定设备信息
		AppDevice *AppDevice `json:"AppDevice" name:"AppDevice"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AppSecureAddDeviceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppSecureAddDeviceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppUpdateDeviceRequest struct {
	*tchttp.BaseRequest

	// 访问Token
	AccessToken *string `json:"AccessToken" name:"AccessToken"`

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 设备别名
	AliasName *string `json:"AliasName" name:"AliasName"`
}

func (r *AppUpdateDeviceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppUpdateDeviceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppUpdateDeviceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备信息
		AppDevice *AppDevice `json:"AppDevice" name:"AppDevice"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AppUpdateDeviceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppUpdateDeviceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppUpdateUserRequest struct {
	*tchttp.BaseRequest

	// 访问Token
	AccessToken *string `json:"AccessToken" name:"AccessToken"`

	// 昵称
	NickName *string `json:"NickName" name:"NickName"`
}

func (r *AppUpdateUserRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppUpdateUserRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppUpdateUserResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 应用用户
		AppUser *AppUser `json:"AppUser" name:"AppUser"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AppUpdateUserResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AppUpdateUserResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AppUser struct {

	// 应用Id
	ApplicationId *string `json:"ApplicationId" name:"ApplicationId"`

	// 用户名
	UserName *string `json:"UserName" name:"UserName"`

	// 昵称
	NickName *string `json:"NickName" name:"NickName"`

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 修改时间
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`
}

type AssociateSubDeviceToGatewayProductRequest struct {
	*tchttp.BaseRequest

	// 子设备产品Id
	SubDeviceProductId *string `json:"SubDeviceProductId" name:"SubDeviceProductId"`

	// 网关产品Id
	GatewayProductId *string `json:"GatewayProductId" name:"GatewayProductId"`
}

func (r *AssociateSubDeviceToGatewayProductRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AssociateSubDeviceToGatewayProductRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AssociateSubDeviceToGatewayProductResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AssociateSubDeviceToGatewayProductResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AssociateSubDeviceToGatewayProductResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type BoolData struct {

	// 名称
	Name *string `json:"Name" name:"Name"`

	// 描述
	Desc *string `json:"Desc" name:"Desc"`

	// 读写模式
	Mode *string `json:"Mode" name:"Mode"`

	// 取值列表
	Range []*bool `json:"Range" name:"Range" list`
}

type CkafkaAction struct {

	// 实例Id
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// topic名称
	TopicName *string `json:"TopicName" name:"TopicName"`

	// 地域
	Region *string `json:"Region" name:"Region"`
}

type DataHistoryEntry struct {

	// 日志id
	Id *string `json:"Id" name:"Id"`

	// 时间戳
	Timestamp *uint64 `json:"Timestamp" name:"Timestamp"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 数据
	Data *string `json:"Data" name:"Data"`
}

type DataTemplate struct {

	// 数字类型
	Number *NumberData `json:"Number" name:"Number"`

	// 字符串类型
	String *StringData `json:"String" name:"String"`

	// 枚举类型
	Enum *EnumData `json:"Enum" name:"Enum"`

	// 布尔类型
	Bool *BoolData `json:"Bool" name:"Bool"`
}

type DeactivateRuleRequest struct {
	*tchttp.BaseRequest

	// 规则Id
	RuleId *string `json:"RuleId" name:"RuleId"`
}

func (r *DeactivateRuleRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeactivateRuleRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeactivateRuleResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeactivateRuleResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeactivateRuleResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DebugLogEntry struct {

	// 日志id
	Id *string `json:"Id" name:"Id"`

	// 行为（事件）
	Event *string `json:"Event" name:"Event"`

	// shadow/action/mqtt, 分别表示：影子/规则引擎/上下线日志
	LogType *string `json:"LogType" name:"LogType"`

	// 时间戳
	Timestamp *uint64 `json:"Timestamp" name:"Timestamp"`

	// success/fail
	Result *string `json:"Result" name:"Result"`

	// 日志详细内容
	Data *string `json:"Data" name:"Data"`

	// 数据来源topic
	Topic *string `json:"Topic" name:"Topic"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`
}

type DeleteDeviceRequest struct {
	*tchttp.BaseRequest

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

	// 产品Id
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

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type DeleteRuleRequest struct {
	*tchttp.BaseRequest

	// 规则Id
	RuleId *string `json:"RuleId" name:"RuleId"`
}

func (r *DeleteRuleRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteRuleRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteRuleResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteRuleResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteRuleResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteTopicRequest struct {
	*tchttp.BaseRequest

	// TopicId
	TopicId *string `json:"TopicId" name:"TopicId"`

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`
}

func (r *DeleteTopicRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteTopicRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteTopicResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteTopicResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteTopicResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Device struct {

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 设备密钥
	DeviceSecret *string `json:"DeviceSecret" name:"DeviceSecret"`

	// 更新时间
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 设备信息（json）
	DeviceInfo *string `json:"DeviceInfo" name:"DeviceInfo"`
}

type DeviceEntry struct {

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 设备密钥
	DeviceSecret *string `json:"DeviceSecret" name:"DeviceSecret"`

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`
}

type DeviceLogEntry struct {

	// 日志id
	Id *string `json:"Id" name:"Id"`

	// 日志内容
	Msg *string `json:"Msg" name:"Msg"`

	// 状态码
	Code *string `json:"Code" name:"Code"`

	// 时间戳
	Timestamp *uint64 `json:"Timestamp" name:"Timestamp"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 设备动作
	Method *string `json:"Method" name:"Method"`
}

type DeviceSignature struct {

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 设备签名
	DeviceSignature *string `json:"DeviceSignature" name:"DeviceSignature"`
}

type DeviceStatData struct {

	// 时间点
	Datetime *string `json:"Datetime" name:"Datetime"`

	// 在线设备数
	DeviceOnline *uint64 `json:"DeviceOnline" name:"DeviceOnline"`

	// 激活设备数
	DeviceActive *uint64 `json:"DeviceActive" name:"DeviceActive"`

	// 设备总数
	DeviceTotal *uint64 `json:"DeviceTotal" name:"DeviceTotal"`
}

type DeviceStatus struct {

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 设备状态（inactive, online, offline）
	Status *string `json:"Status" name:"Status"`

	// 首次上线时间
	FirstOnline *string `json:"FirstOnline" name:"FirstOnline"`

	// 最后上线时间
	LastOnline *string `json:"LastOnline" name:"LastOnline"`

	// 上线次数
	OnlineTimes *uint64 `json:"OnlineTimes" name:"OnlineTimes"`
}

type EnumData struct {

	// 名称
	Name *string `json:"Name" name:"Name"`

	// 描述
	Desc *string `json:"Desc" name:"Desc"`

	// 读写模式
	Mode *string `json:"Mode" name:"Mode"`

	// 取值列表
	Range []*string `json:"Range" name:"Range" list`
}

type GetDataHistoryRequest struct {
	*tchttp.BaseRequest

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称列表，允许最多一次100台
	DeviceNames []*string `json:"DeviceNames" name:"DeviceNames" list`

	// 查询开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 查询结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 查询数据量
	Size []*uint64 `json:"Size" name:"Size" list`

	// 时间排序（desc/asc）
	Order *string `json:"Order" name:"Order"`

	// 查询游标
	ScrollId *string `json:"ScrollId" name:"ScrollId"`
}

func (r *GetDataHistoryRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDataHistoryRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDataHistoryResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 数据历史
		DataHistory []*DataHistoryEntry `json:"DataHistory" name:"DataHistory" list`

		// 查询游标
		ScrollId *string `json:"ScrollId" name:"ScrollId"`

		// 查询游标超时
		ScrollTimeout *uint64 `json:"ScrollTimeout" name:"ScrollTimeout"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetDataHistoryResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDataHistoryResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDebugLogRequest struct {
	*tchttp.BaseRequest

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称列表，最大支持100台
	DeviceNames []*string `json:"DeviceNames" name:"DeviceNames" list`

	// 查询开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 查询结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 查询数据量
	Size *uint64 `json:"Size" name:"Size"`

	// 时间排序（desc/asc）
	Order *string `json:"Order" name:"Order"`

	// 查询游标
	ScrollId *string `json:"ScrollId" name:"ScrollId"`

	// 日志类型（shadow/action/mqtt）
	Type *string `json:"Type" name:"Type"`
}

func (r *GetDebugLogRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDebugLogRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDebugLogResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 调试日志
		DebugLog []*DebugLogEntry `json:"DebugLog" name:"DebugLog" list`

		// 查询游标
		ScrollId *string `json:"ScrollId" name:"ScrollId"`

		// 游标超时
		ScrollTimeout *uint64 `json:"ScrollTimeout" name:"ScrollTimeout"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetDebugLogResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDebugLogResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDeviceDataRequest struct {
	*tchttp.BaseRequest

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`
}

func (r *GetDeviceDataRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDeviceDataRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDeviceDataResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备数据
		DeviceData *string `json:"DeviceData" name:"DeviceData"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetDeviceDataResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDeviceDataResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDeviceLogRequest struct {
	*tchttp.BaseRequest

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称列表，最大支持100台
	DeviceNames []*string `json:"DeviceNames" name:"DeviceNames" list`

	// 查询开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 查询结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 查询数据量
	Size *uint64 `json:"Size" name:"Size"`

	// 时间排序（desc/asc）
	Order *string `json:"Order" name:"Order"`

	// 查询游标
	ScrollId *string `json:"ScrollId" name:"ScrollId"`

	// 日志类型（comm/status）
	Type *string `json:"Type" name:"Type"`
}

func (r *GetDeviceLogRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDeviceLogRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDeviceLogResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备日志
		DeviceLog []*DeviceLogEntry `json:"DeviceLog" name:"DeviceLog" list`

		// 查询游标
		ScrollId *string `json:"ScrollId" name:"ScrollId"`

		// 游标超时
		ScrollTimeout *uint64 `json:"ScrollTimeout" name:"ScrollTimeout"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetDeviceLogResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDeviceLogResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDeviceRequest struct {
	*tchttp.BaseRequest

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`
}

func (r *GetDeviceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDeviceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDeviceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备信息
		Device *Device `json:"Device" name:"Device"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetDeviceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDeviceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDeviceSignaturesRequest struct {
	*tchttp.BaseRequest

	// 产品ID
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称列表（单次限制1000个设备）
	DeviceNames []*string `json:"DeviceNames" name:"DeviceNames" list`

	// 过期时间
	Expire *uint64 `json:"Expire" name:"Expire"`
}

func (r *GetDeviceSignaturesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDeviceSignaturesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDeviceSignaturesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备绑定签名列表
		DeviceSignatures []*DeviceSignature `json:"DeviceSignatures" name:"DeviceSignatures" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetDeviceSignaturesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDeviceSignaturesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDeviceStatisticsRequest struct {
	*tchttp.BaseRequest

	// 产品Id列表
	Products []*string `json:"Products" name:"Products" list`

	// 开始日期
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 结束日期
	EndDate *string `json:"EndDate" name:"EndDate"`
}

func (r *GetDeviceStatisticsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDeviceStatisticsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDeviceStatisticsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 统计数据
		DeviceStatistics []*DeviceStatData `json:"DeviceStatistics" name:"DeviceStatistics" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetDeviceStatisticsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDeviceStatisticsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDeviceStatusesRequest struct {
	*tchttp.BaseRequest

	// 产品ID
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称列表（单次限制1000个设备）
	DeviceNames []*string `json:"DeviceNames" name:"DeviceNames" list`
}

func (r *GetDeviceStatusesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDeviceStatusesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDeviceStatusesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备状态列表
		DeviceStatuses []*DeviceStatus `json:"DeviceStatuses" name:"DeviceStatuses" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetDeviceStatusesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDeviceStatusesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDevicesRequest struct {
	*tchttp.BaseRequest

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 偏移
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 长度
	Length *uint64 `json:"Length" name:"Length"`

	// 关键字查询
	Keyword *string `json:"Keyword" name:"Keyword"`
}

func (r *GetDevicesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDevicesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetDevicesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备列表
		Devices []*DeviceEntry `json:"Devices" name:"Devices" list`

		// 设备总数
		Total *uint64 `json:"Total" name:"Total"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetDevicesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetDevicesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetProductRequest struct {
	*tchttp.BaseRequest

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`
}

func (r *GetProductRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetProductRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetProductResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 产品信息
		Product *Product `json:"Product" name:"Product"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetProductResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetProductResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetProductsRequest struct {
	*tchttp.BaseRequest

	// 偏移
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 长度
	Length *uint64 `json:"Length" name:"Length"`
}

func (r *GetProductsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetProductsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetProductsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// Product列表
		Products []*ProductEntry `json:"Products" name:"Products" list`

		// Product总数
		Total *uint64 `json:"Total" name:"Total"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetProductsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetProductsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetRuleRequest struct {
	*tchttp.BaseRequest

	// 规则Id
	RuleId *string `json:"RuleId" name:"RuleId"`
}

func (r *GetRuleRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetRuleRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetRuleResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 规则
		Rule *Rule `json:"Rule" name:"Rule"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetRuleResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetRuleResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetRulesRequest struct {
	*tchttp.BaseRequest

	// 偏移
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 长度
	Length *uint64 `json:"Length" name:"Length"`
}

func (r *GetRulesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetRulesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetRulesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 规则列表
		Rules []*Rule `json:"Rules" name:"Rules" list`

		// 规则总数
		Total *uint64 `json:"Total" name:"Total"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetRulesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetRulesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetTopicRequest struct {
	*tchttp.BaseRequest

	// TopicId
	TopicId *string `json:"TopicId" name:"TopicId"`

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`
}

func (r *GetTopicRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetTopicRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetTopicResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// Topic信息
		Topic *Topic `json:"Topic" name:"Topic"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetTopicResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetTopicResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetTopicsRequest struct {
	*tchttp.BaseRequest

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 偏移
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 长度
	Length *uint64 `json:"Length" name:"Length"`
}

func (r *GetTopicsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetTopicsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetTopicsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// Topic列表
		Topics []*Topic `json:"Topics" name:"Topics" list`

		// Topic总数
		Total *uint64 `json:"Total" name:"Total"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetTopicsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetTopicsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type IssueDeviceControlRequest struct {
	*tchttp.BaseRequest

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`

	// 控制数据（json）
	ControlData *string `json:"ControlData" name:"ControlData"`

	// 是否发送metadata字段
	Metadata *bool `json:"Metadata" name:"Metadata"`
}

func (r *IssueDeviceControlRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *IssueDeviceControlRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type IssueDeviceControlResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *IssueDeviceControlResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *IssueDeviceControlResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type NumberData struct {

	// 名称
	Name *string `json:"Name" name:"Name"`

	// 描述
	Desc *string `json:"Desc" name:"Desc"`

	// 读写模式
	Mode *string `json:"Mode" name:"Mode"`

	// 取值范围
	Range []*float64 `json:"Range" name:"Range" list`
}

type Product struct {

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 产品Key
	ProductKey *string `json:"ProductKey" name:"ProductKey"`

	// AppId
	AppId *uint64 `json:"AppId" name:"AppId"`

	// 产品名称
	Name *string `json:"Name" name:"Name"`

	// 产品描述
	Description *string `json:"Description" name:"Description"`

	// 连接域名
	Domain *string `json:"Domain" name:"Domain"`

	// 产品规格
	Standard *uint64 `json:"Standard" name:"Standard"`

	// 鉴权类型（0：直连，1：Token）
	AuthType *uint64 `json:"AuthType" name:"AuthType"`

	// 删除（0未删除）
	Deleted *uint64 `json:"Deleted" name:"Deleted"`

	// 备注
	Message *string `json:"Message" name:"Message"`

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 更新时间
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`

	// 数据模版
	DataTemplate []*DataTemplate `json:"DataTemplate" name:"DataTemplate" list`

	// 数据协议（native/template）
	DataProtocol *string `json:"DataProtocol" name:"DataProtocol"`

	// 直连用户名
	Username *string `json:"Username" name:"Username"`

	// 直连密码
	Password *string `json:"Password" name:"Password"`

	// 通信方式
	CommProtocol *string `json:"CommProtocol" name:"CommProtocol"`

	// qps
	Qps *uint64 `json:"Qps" name:"Qps"`

	// 地域
	Region *string `json:"Region" name:"Region"`

	// 产品的设备类型
	DeviceType *string `json:"DeviceType" name:"DeviceType"`

	// 关联的产品列表
	AssociatedProducts []*string `json:"AssociatedProducts" name:"AssociatedProducts" list`
}

type ProductEntry struct {

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 产品Key
	ProductKey *string `json:"ProductKey" name:"ProductKey"`

	// AppId
	AppId *uint64 `json:"AppId" name:"AppId"`

	// 产品名称
	Name *string `json:"Name" name:"Name"`

	// 产品描述
	Description *string `json:"Description" name:"Description"`

	// 连接域名
	Domain *string `json:"Domain" name:"Domain"`

	// 鉴权类型（0：直连，1：Token）
	AuthType *uint64 `json:"AuthType" name:"AuthType"`

	// 数据协议（native/template）
	DataProtocol *string `json:"DataProtocol" name:"DataProtocol"`

	// 删除（0未删除）
	Deleted *uint64 `json:"Deleted" name:"Deleted"`

	// 备注
	Message *string `json:"Message" name:"Message"`

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 通信方式
	CommProtocol *string `json:"CommProtocol" name:"CommProtocol"`

	// 地域
	Region *string `json:"Region" name:"Region"`

	// 设备类型
	DeviceType *string `json:"DeviceType" name:"DeviceType"`
}

type PublishMsgRequest struct {
	*tchttp.BaseRequest

	// Topic
	Topic *string `json:"Topic" name:"Topic"`

	// 消息内容
	Message *string `json:"Message" name:"Message"`

	// Qos(目前QoS支持0与1)
	Qos *int64 `json:"Qos" name:"Qos"`
}

func (r *PublishMsgRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *PublishMsgRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type PublishMsgResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *PublishMsgResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *PublishMsgResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ResetDeviceRequest struct {
	*tchttp.BaseRequest

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 设备名称
	DeviceName *string `json:"DeviceName" name:"DeviceName"`
}

func (r *ResetDeviceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ResetDeviceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ResetDeviceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 设备信息
		Device *Device `json:"Device" name:"Device"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ResetDeviceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ResetDeviceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Rule struct {

	// 规则Id
	RuleId *string `json:"RuleId" name:"RuleId"`

	// AppId
	AppId *uint64 `json:"AppId" name:"AppId"`

	// 名称
	Name *string `json:"Name" name:"Name"`

	// 描述
	Description *string `json:"Description" name:"Description"`

	// 查询
	Query *RuleQuery `json:"Query" name:"Query"`

	// 转发
	Actions []*Action `json:"Actions" name:"Actions" list`

	// 已启动
	Active *uint64 `json:"Active" name:"Active"`

	// 已删除
	Deleted *uint64 `json:"Deleted" name:"Deleted"`

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 更新时间
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`

	// 消息顺序
	MsgOrder *uint64 `json:"MsgOrder" name:"MsgOrder"`

	// 数据类型（0：文本，1：二进制）
	DataType *uint64 `json:"DataType" name:"DataType"`
}

type RuleQuery struct {

	// 字段
	Field *string `json:"Field" name:"Field"`

	// 过滤规则
	Condition *string `json:"Condition" name:"Condition"`

	// Topic
	Topic *string `json:"Topic" name:"Topic"`

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`
}

type ServiceAction struct {

	// 服务url地址
	Url *string `json:"Url" name:"Url"`
}

type StringData struct {

	// 名称
	Name *string `json:"Name" name:"Name"`

	// 描述
	Desc *string `json:"Desc" name:"Desc"`

	// 读写模式
	Mode *string `json:"Mode" name:"Mode"`

	// 长度范围
	Range []*uint64 `json:"Range" name:"Range" list`
}

type Topic struct {

	// TopicId
	TopicId *string `json:"TopicId" name:"TopicId"`

	// Topic名称
	TopicName *string `json:"TopicName" name:"TopicName"`

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 消息最大生命周期
	MsgLife *uint64 `json:"MsgLife" name:"MsgLife"`

	// 消息最大大小
	MsgSize *uint64 `json:"MsgSize" name:"MsgSize"`

	// 消息最大数量
	MsgCount *uint64 `json:"MsgCount" name:"MsgCount"`

	// 已删除
	Deleted *uint64 `json:"Deleted" name:"Deleted"`

	// Topic完整路径
	Path *string `json:"Path" name:"Path"`

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 更新时间
	UpdateTime *string `json:"UpdateTime" name:"UpdateTime"`
}

type TopicAction struct {

	// 目标topic
	Topic *string `json:"Topic" name:"Topic"`
}

type UnassociateSubDeviceFromGatewayProductRequest struct {
	*tchttp.BaseRequest

	// 子设备产品Id
	SubDeviceProductId *string `json:"SubDeviceProductId" name:"SubDeviceProductId"`

	// 网关设备产品Id
	GatewayProductId *string `json:"GatewayProductId" name:"GatewayProductId"`
}

func (r *UnassociateSubDeviceFromGatewayProductRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UnassociateSubDeviceFromGatewayProductRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UnassociateSubDeviceFromGatewayProductResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UnassociateSubDeviceFromGatewayProductResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UnassociateSubDeviceFromGatewayProductResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateProductRequest struct {
	*tchttp.BaseRequest

	// 产品Id
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 产品名称
	Name *string `json:"Name" name:"Name"`

	// 产品描述
	Description *string `json:"Description" name:"Description"`

	// 数据模版
	DataTemplate []*DataTemplate `json:"DataTemplate" name:"DataTemplate" list`
}

func (r *UpdateProductRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateProductRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateProductResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 更新后的产品信息
		Product *Product `json:"Product" name:"Product"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UpdateProductResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateProductResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateRuleRequest struct {
	*tchttp.BaseRequest

	// 规则Id
	RuleId *string `json:"RuleId" name:"RuleId"`

	// 名称
	Name *string `json:"Name" name:"Name"`

	// 描述
	Description *string `json:"Description" name:"Description"`

	// 查询
	Query *RuleQuery `json:"Query" name:"Query"`

	// 转发动作列表
	Actions []*Action `json:"Actions" name:"Actions" list`

	// 数据类型（0：文本，1：二进制）
	DataType *uint64 `json:"DataType" name:"DataType"`
}

func (r *UpdateRuleRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateRuleRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateRuleResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 规则
		Rule *Rule `json:"Rule" name:"Rule"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UpdateRuleResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateRuleResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}
