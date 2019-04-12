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

package v20180317

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type Backend struct {

	// 转发目标的类型，目前仅可取值为 CVM
	Type *string `json:"Type" name:"Type"`

	// 云服务器的唯一 ID，可通过 DescribeInstances 接口返回字段中的 unInstanceId 字段获取
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 后端云服务器监听端口
	Port *int64 `json:"Port" name:"Port"`

	// 后端云服务器的转发权重，取值范围：0~100，默认为 10。
	Weight *int64 `json:"Weight" name:"Weight"`

	// 云服务器的外网 IP
	PublicIpAddresses []*string `json:"PublicIpAddresses" name:"PublicIpAddresses" list`

	// 云服务器的内网 IP
	PrivateIpAddresses []*string `json:"PrivateIpAddresses" name:"PrivateIpAddresses" list`

	// 云服务器实例名称
	InstanceName *string `json:"InstanceName" name:"InstanceName"`

	// 云服务器被绑定到监听器的时间
	RegisteredTime *string `json:"RegisteredTime" name:"RegisteredTime"`
}

type CertificateInput struct {

	// 认证类型，UNIDIRECTIONAL：单向认证，MUTUAL：双向认证
	SSLMode *string `json:"SSLMode" name:"SSLMode"`

	// 服务端证书的 ID，如果不填写此项则必须上传证书，包括 CertContent，CertKey，CertName。
	CertId *string `json:"CertId" name:"CertId"`

	// 客户端证书的 ID，如果 SSLMode=mutual，监听器如果不填写此项则必须上传客户端证书，包括 CertCaContent，CertCaName。
	CertCaId *string `json:"CertCaId" name:"CertCaId"`

	// 上传服务端证书的名称，如果没有 CertId，则此项必传。
	CertName *string `json:"CertName" name:"CertName"`

	// 上传服务端证书的 key，如果没有 CertId，则此项必传。
	CertKey *string `json:"CertKey" name:"CertKey"`

	// 上传服务端证书的内容，如果没有 CertId，则此项必传。
	CertContent *string `json:"CertContent" name:"CertContent"`

	// 上传客户端 CA 证书的名称，如果 SSLMode=mutual，如果没有 CertCaId，则此项必传。
	CertCaName *string `json:"CertCaName" name:"CertCaName"`

	// 上传客户端证书的内容，如果 SSLMode=mutual，如果没有 CertCaId，则此项必传。
	CertCaContent *string `json:"CertCaContent" name:"CertCaContent"`
}

type CertificateOutput struct {

	// 认证类型，unidirectional：单向认证，mutual：双向认证
	SSLMode *string `json:"SSLMode" name:"SSLMode"`

	// 服务端证书的 ID。
	CertId *string `json:"CertId" name:"CertId"`

	// 客户端证书的 ID。
	CertCaId *string `json:"CertCaId" name:"CertCaId"`
}

type CreateListenerRequest struct {
	*tchttp.BaseRequest

	// 负载均衡实例 ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 要将监听器创建到哪些端口，每个端口对应一个新的监听器
	Ports []*int64 `json:"Ports" name:"Ports" list`

	// 监听器协议：HTTP | HTTPS | TCP | TCP_SSL
	Protocol *string `json:"Protocol" name:"Protocol"`

	// 要创建的监听器名称列表，名称与Ports数组按序一一对应，如不需立即命名，则无需提供此参数
	ListenerNames []*string `json:"ListenerNames" name:"ListenerNames" list`

	// 健康检查相关参数，此参数仅适用于TCP/UDP/TCP_SSL监听器
	HealthCheck *HealthCheck `json:"HealthCheck" name:"HealthCheck"`

	// 证书相关信息，此参数仅适用于HTTPS/TCP_SSL监听器
	Certificate *CertificateInput `json:"Certificate" name:"Certificate"`

	// 会话保持时间，单位：秒。可选值：30~3600，默认 0，表示不开启。此参数仅适用于TCP/UDP监听器。
	SessionExpireTime *int64 `json:"SessionExpireTime" name:"SessionExpireTime"`

	// 监听器转发的方式。可选值：WRR、LEAST_CONN
	// 分别表示按权重轮询、最小连接数， 默认为 WRR。此参数仅适用于TCP/UDP/TCP_SSL监听器。
	Scheduler *string `json:"Scheduler" name:"Scheduler"`

	// 是否开启SNI特性，此参数仅适用于HTTPS监听器。
	SniSwitch *int64 `json:"SniSwitch" name:"SniSwitch"`
}

func (r *CreateListenerRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateListenerRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateListenerResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 创建的监听器的唯一标识数组
		ListenerIds []*string `json:"ListenerIds" name:"ListenerIds" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateListenerResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateListenerResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateLoadBalancerRequest struct {
	*tchttp.BaseRequest

	// 负载均衡实例的网络类型：
	// OPEN：公网属性， INTERNAL：内网属性。
	LoadBalancerType *string `json:"LoadBalancerType" name:"LoadBalancerType"`

	// 负载均衡实例。1：应用型，0：传统型，默认为应用型负载均衡实例。
	Forward *int64 `json:"Forward" name:"Forward"`

	// 负载均衡实例的名称，只用来创建一个的时候生效。规则：1-50 个英文、汉字、数字、连接线“-”或下划线“_”。
	// 注意：如果名称与系统中已有负载均衡实例的名称重复的话，则系统将会自动生成此次创建的负载均衡实例的名称。
	LoadBalancerName *string `json:"LoadBalancerName" name:"LoadBalancerName"`

	// 负载均衡后端实例所属网络 ID，可以通过 DescribeVpcEx 接口获取。 不填则默认为基础网络。
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 在私有网络内购买内网负载均衡实例的时候需要指定子网 ID，内网负载均衡实例的 VIP 将从这个子网中产生。其他情况不用填写该字段。
	SubnetId *string `json:"SubnetId" name:"SubnetId"`

	// 负载均衡实例所属的项目 ID，可以通过 DescribeProject 接口获取。不填则属于默认项目。
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`
}

func (r *CreateLoadBalancerRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateLoadBalancerRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateLoadBalancerResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 由负载均衡实例统一 ID 组成的数组。
		LoadBalancerIds []*string `json:"LoadBalancerIds" name:"LoadBalancerIds" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateLoadBalancerResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateLoadBalancerResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateRuleRequest struct {
	*tchttp.BaseRequest

	// 负载均衡实例 ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 监听器 ID
	ListenerId *string `json:"ListenerId" name:"ListenerId"`

	// 新建转发规则的信息
	Rules []*RuleInput `json:"Rules" name:"Rules" list`
}

func (r *CreateRuleRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateRuleRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateRuleResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateRuleResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateRuleResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteListenerRequest struct {
	*tchttp.BaseRequest

	// 应用型负载均衡实例 ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 要删除的监听器 ID
	ListenerId *string `json:"ListenerId" name:"ListenerId"`
}

func (r *DeleteListenerRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteListenerRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteListenerResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteListenerResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteListenerResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteLoadBalancerRequest struct {
	*tchttp.BaseRequest

	// 要删除的负载均衡实例 ID数组
	LoadBalancerIds []*string `json:"LoadBalancerIds" name:"LoadBalancerIds" list`
}

func (r *DeleteLoadBalancerRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteLoadBalancerRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteLoadBalancerResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteLoadBalancerResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteLoadBalancerResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteRuleRequest struct {
	*tchttp.BaseRequest

	// 负载均衡实例 ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 应用型负载均衡监听器 ID
	ListenerId *string `json:"ListenerId" name:"ListenerId"`

	// 要删除的转发规则的ID组成的数组
	LocationIds []*string `json:"LocationIds" name:"LocationIds" list`

	// 要删除的转发规则的域名，已提供LocationIds参数时本参数不生效
	Domain *string `json:"Domain" name:"Domain"`

	// 要删除的转发规则的转发路径，已提供LocationIds参数时本参数不生效
	Url *string `json:"Url" name:"Url"`
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

type DeregisterTargetsRequest struct {
	*tchttp.BaseRequest

	// 负载均衡实例 ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 监听器 ID
	ListenerId *string `json:"ListenerId" name:"ListenerId"`

	// 要解绑的后端机器列表
	Targets []*Target `json:"Targets" name:"Targets" list`

	// 转发规则的ID，当从七层转发规则解绑机器时，必须提供此参数或Domain+Url两者之一
	LocationId *string `json:"LocationId" name:"LocationId"`

	// 目标规则的域名，提供LocationId参数时本参数不生效
	Domain *string `json:"Domain" name:"Domain"`

	// 目标规则的URL，提供LocationId参数时本参数不生效
	Url *string `json:"Url" name:"Url"`
}

func (r *DeregisterTargetsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeregisterTargetsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeregisterTargetsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeregisterTargetsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeregisterTargetsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeListenersRequest struct {
	*tchttp.BaseRequest

	// 负载均衡实例 ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 要查询的应用型负载均衡监听器 ID数组
	ListenerIds []*string `json:"ListenerIds" name:"ListenerIds" list`

	// 要查询的监听器协议类型，取值 TCP | UDP | HTTP | HTTPS | TCP_SSL
	Protocol *string `json:"Protocol" name:"Protocol"`

	// 要查询的监听器的端口
	Port *int64 `json:"Port" name:"Port"`
}

func (r *DescribeListenersRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeListenersRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeListenersResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 监听器列表
		Listeners []*Listener `json:"Listeners" name:"Listeners" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeListenersResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeListenersResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLoadBalancersRequest struct {
	*tchttp.BaseRequest

	// 负载均衡实例 ID。
	LoadBalancerIds []*string `json:"LoadBalancerIds" name:"LoadBalancerIds" list`

	// 负载均衡实例的网络类型：
	// OPEN：公网属性， INTERNAL：内网属性。
	LoadBalancerType *string `json:"LoadBalancerType" name:"LoadBalancerType"`

	// 1：应用型，0：传统型。
	Forward *int64 `json:"Forward" name:"Forward"`

	// 负载均衡实例名称。
	LoadBalancerName *string `json:"LoadBalancerName" name:"LoadBalancerName"`

	// 腾讯云为负载均衡实例分配的域名，应用型负载均衡该字段无意义。
	Domain *string `json:"Domain" name:"Domain"`

	// 负载均衡实例的 VIP 地址，支持多个。
	LoadBalancerVips []*string `json:"LoadBalancerVips" name:"LoadBalancerVips" list`

	// 后端云服务器的外网 IP。
	BackendPublicIps []*string `json:"BackendPublicIps" name:"BackendPublicIps" list`

	// 后端云服务器的内网 IP。
	BackendPrivateIps []*string `json:"BackendPrivateIps" name:"BackendPrivateIps" list`

	// 数据偏移量，默认为 0。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 返回负载均衡个数，默认为 20。
	Limit *int64 `json:"Limit" name:"Limit"`

	// 排序字段，支持以下字段：LoadBalancerName，CreateTime，Domain，LoadBalancerType。
	OrderBy *string `json:"OrderBy" name:"OrderBy"`

	// 1：倒序，0：顺序，默认按照创建时间倒序。
	OrderType *int64 `json:"OrderType" name:"OrderType"`

	// 搜索字段，模糊匹配名称、域名、VIP。
	SearchKey *string `json:"SearchKey" name:"SearchKey"`

	// 负载均衡实例所属的项目 ID，可以通过 DescribeProject 接口获取。
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 查询的负载均衡是否绑定后端服务器，0：没有绑定云服务器，1：绑定云服务器，-1：查询全部。
	WithRs *int64 `json:"WithRs" name:"WithRs"`
}

func (r *DescribeLoadBalancersRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLoadBalancersRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLoadBalancersResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 满足过滤条件的负载均衡实例总数。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 返回的负载均衡实例数组。
		LoadBalancerSet []*LoadBalancer `json:"LoadBalancerSet" name:"LoadBalancerSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeLoadBalancersResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLoadBalancersResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTargetsRequest struct {
	*tchttp.BaseRequest

	// 负载均衡实例 ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 监听器 ID列表
	ListenerIds []*string `json:"ListenerIds" name:"ListenerIds" list`

	// 监听器协议类型
	Protocol *string `json:"Protocol" name:"Protocol"`

	// 负载均衡监听器端口
	Port *int64 `json:"Port" name:"Port"`
}

func (r *DescribeTargetsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTargetsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTargetsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 监听器后端绑定的机器信息
		Listeners []*ListenerBackend `json:"Listeners" name:"Listeners" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeTargetsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTargetsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTaskStatusRequest struct {
	*tchttp.BaseRequest

	// 请求ID，即接口返回的RequestId
	TaskId *string `json:"TaskId" name:"TaskId"`
}

func (r *DescribeTaskStatusRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTaskStatusRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTaskStatusResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务的当前状态。 0：成功，1：失败，2：进行中。
		Status *int64 `json:"Status" name:"Status"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeTaskStatusResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTaskStatusResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type HealthCheck struct {

	// 是否开启健康检查：1（开启）、0（关闭）。
	HealthSwitch *int64 `json:"HealthSwitch" name:"HealthSwitch"`

	// 健康检查的响应超时时间，可选值：2~60，默认值：2，单位：秒。响应超时时间要小于检查间隔时间。
	TimeOut *int64 `json:"TimeOut" name:"TimeOut"`

	// 健康检查探测间隔时间，默认值：5，可选值：5~300，单位：秒。
	IntervalTime *int64 `json:"IntervalTime" name:"IntervalTime"`

	// 健康阈值，默认值：3，表示当连续探测三次健康则表示该转发正常，可选值：2~10，单位：次。
	HealthNum *int64 `json:"HealthNum" name:"HealthNum"`

	// 不健康阈值，默认值：3，表示当连续探测三次不健康则表示该转发异常，可选值：2~10，单位：次。
	UnHealthNum *int64 `json:"UnHealthNum" name:"UnHealthNum"`

	// 健康检查状态码（仅适用于HTTP/HTTPS转发规则）。可选值：1~31，默认 31。
	// 1 表示探测后返回值 1xx 表示健康，2 表示返回 2xx 表示健康，4 表示返回 3xx 表示健康，8 表示返回 4xx 表示健康，16 表示返回 5xx 表示健康。若希望多种码都表示健康，则将相应的值相加。
	HttpCode *int64 `json:"HttpCode" name:"HttpCode"`

	// 健康检查路径（仅适用于HTTP/HTTPS转发规则）。
	HttpCheckPath *string `json:"HttpCheckPath" name:"HttpCheckPath"`

	// 健康检查域名（仅适用于HTTP/HTTPS转发规则）。
	HttpCheckDomain *string `json:"HttpCheckDomain" name:"HttpCheckDomain"`

	// 健康检查方法（仅适用于HTTP/HTTPS转发规则），取值为HEAD或GET。
	HttpCheckMethod *string `json:"HttpCheckMethod" name:"HttpCheckMethod"`
}

type Listener struct {

	// 应用型负载均衡监听器 ID
	ListenerId *string `json:"ListenerId" name:"ListenerId"`

	// 监听器协议
	Protocol *string `json:"Protocol" name:"Protocol"`

	// 监听器端口
	Port *int64 `json:"Port" name:"Port"`

	// 监听器绑定的证书信息
	Certificate *CertificateOutput `json:"Certificate" name:"Certificate"`

	// 监听器的健康检查信息
	HealthCheck *HealthCheck `json:"HealthCheck" name:"HealthCheck"`

	// 请求调度方式
	Scheduler *string `json:"Scheduler" name:"Scheduler"`

	// 会话保持时间
	SessionExpireTime *int64 `json:"SessionExpireTime" name:"SessionExpireTime"`

	// 是否开启SNI特性（本参数仅对于HTTPS监听器有意义）
	SniSwitch *int64 `json:"SniSwitch" name:"SniSwitch"`

	// 监听器下的全部转发规则（本参数仅对于HTTP/HTTPS监听器有意义）
	Rules []*RuleOutput `json:"Rules" name:"Rules" list`

	// 监听器的名称
	ListenerName *string `json:"ListenerName" name:"ListenerName"`
}

type ListenerBackend struct {

	// 监听器 ID
	ListenerId *string `json:"ListenerId" name:"ListenerId"`

	// 监听器的协议
	Protocol *string `json:"Protocol" name:"Protocol"`

	// 监听器的端口
	Port *int64 `json:"Port" name:"Port"`

	// 监听器下的规则信息（仅适用于HTTP/HTTPS监听器）
	Rules []*RuleTargets `json:"Rules" name:"Rules" list`

	// 监听器上注册的机器列表（仅适用于TCP/UDP/TCP_SSL监听器）
	Targets []*Backend `json:"Targets" name:"Targets" list`
}

type LoadBalancer struct {

	// 负载均衡实例 ID。
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 负载均衡实例的名称。
	LoadBalancerName *string `json:"LoadBalancerName" name:"LoadBalancerName"`

	// 负载均衡实例的网络类型：
	// OPEN：公网属性， INTERNAL：内网属性。
	LoadBalancerType *string `json:"LoadBalancerType" name:"LoadBalancerType"`

	// 应用型负载均衡标识，1：应用型负载均衡，0：传统型的负载均衡。
	Forward *uint64 `json:"Forward" name:"Forward"`

	// 负载均衡实例的域名，内网类型负载均衡以及应用型负载均衡实例不提供该字段
	Domain *string `json:"Domain" name:"Domain"`

	// 负载均衡实例的 VIP 列表。
	LoadBalancerVips []*string `json:"LoadBalancerVips" name:"LoadBalancerVips" list`

	// 负载均衡实例的状态，包括
	// 0：创建中，1：正常运行。
	Status *uint64 `json:"Status" name:"Status"`

	// 负载均衡实例的创建时间。
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 负载均衡实例的上次状态转换时间。
	StatusTime *string `json:"StatusTime" name:"StatusTime"`

	// 负载均衡实例所属的项目 ID， 0 表示默认项目。
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`

	// 私有网络的 ID
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 高防 LB 的标识，1：高防负载均衡 0：非高防负载均衡。
	OpenBgp *uint64 `json:"OpenBgp" name:"OpenBgp"`

	// 在 2016 年 12 月份之前的传统型内网负载均衡都是开启了 snat 的。
	Snat *bool `json:"Snat" name:"Snat"`

	// 0：表示未被隔离，1：表示被隔离。
	Isolation *uint64 `json:"Isolation" name:"Isolation"`

	// 用户开启日志的信息，日志只有公网属性创建了 HTTP 、HTTPS 监听器的负载均衡才会有日志。
	Log *string `json:"Log" name:"Log"`

	// 负载均衡实例所在的子网（仅对内网VPC型LB有意义）
	SubnetId *string `json:"SubnetId" name:"SubnetId"`
}

type ModifyDomainRequest struct {
	*tchttp.BaseRequest

	// 负载均衡实例 ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 应用型负载均衡监听器 ID
	ListenerId *string `json:"ListenerId" name:"ListenerId"`

	// 监听器下的某个旧域名。
	Domain *string `json:"Domain" name:"Domain"`

	// 新域名，	长度限制为：1-80。有三种使用格式：非正则表达式格式，通配符格式，正则表达式格式。非正则表达式格式只能使用字母、数字、‘-’、‘.’。通配符格式的使用 ‘*’ 只能在开头或者结尾。正则表达式以'~'开头。
	NewDomain *string `json:"NewDomain" name:"NewDomain"`
}

func (r *ModifyDomainRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDomainRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDomainResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDomainResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDomainResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyListenerRequest struct {
	*tchttp.BaseRequest

	// 负载均衡实例 ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 负载均衡监听器 ID
	ListenerId *string `json:"ListenerId" name:"ListenerId"`

	// 新的监听器名称
	ListenerName *string `json:"ListenerName" name:"ListenerName"`

	// 会话保持时间，单位：秒。可选值：30~3600，默认 0，表示不开启。此参数仅适用于TCP/UDP监听器。
	SessionExpireTime *int64 `json:"SessionExpireTime" name:"SessionExpireTime"`

	// 健康检查相关参数，此参数仅适用于TCP/UDP/TCP_SSL监听器
	HealthCheck *HealthCheck `json:"HealthCheck" name:"HealthCheck"`

	// 证书相关信息，此参数仅适用于HTTPS/TCP_SSL监听器
	Certificate *CertificateInput `json:"Certificate" name:"Certificate"`

	// 监听器转发的方式。可选值：WRR、LEAST_CONN
	// 分别表示按权重轮询、最小连接数， 默认为 WRR。
	Scheduler *string `json:"Scheduler" name:"Scheduler"`
}

func (r *ModifyListenerRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyListenerRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyListenerResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyListenerResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyListenerResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyLoadBalancerAttributesRequest struct {
	*tchttp.BaseRequest

	// 负载均衡的唯一ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 负载均衡实例名称
	LoadBalancerName *string `json:"LoadBalancerName" name:"LoadBalancerName"`
}

func (r *ModifyLoadBalancerAttributesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyLoadBalancerAttributesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyLoadBalancerAttributesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyLoadBalancerAttributesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyLoadBalancerAttributesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyRuleRequest struct {
	*tchttp.BaseRequest

	// 负载均衡实例 ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 应用型负载均衡监听器 ID
	ListenerId *string `json:"ListenerId" name:"ListenerId"`

	// 要修改的转发规则的 ID。
	LocationId *string `json:"LocationId" name:"LocationId"`

	// 转发规则的新的转发路径，如不需修改Url，则不需提供此参数
	Url *string `json:"Url" name:"Url"`

	// 健康检查信息
	HealthCheck *HealthCheck `json:"HealthCheck" name:"HealthCheck"`

	// 规则的请求转发方式
	Scheduler *string `json:"Scheduler" name:"Scheduler"`

	// 会话保持时间
	SessionExpireTime *int64 `json:"SessionExpireTime" name:"SessionExpireTime"`
}

func (r *ModifyRuleRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyRuleRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyRuleResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyRuleResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyRuleResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyTargetPortRequest struct {
	*tchttp.BaseRequest

	// 负载均衡实例 ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 负载均衡监听器 ID
	ListenerId *string `json:"ListenerId" name:"ListenerId"`

	// 要修改端口的后端机器列表
	Targets []*Target `json:"Targets" name:"Targets" list`

	// 后端机器绑定到监听器的新端口
	NewPort *int64 `json:"NewPort" name:"NewPort"`

	// 转发规则的ID
	LocationId *string `json:"LocationId" name:"LocationId"`

	// 目标规则的域名，提供LocationId参数时本参数不生效
	Domain *string `json:"Domain" name:"Domain"`

	// 目标规则的URL，提供LocationId参数时本参数不生效
	Url *string `json:"Url" name:"Url"`
}

func (r *ModifyTargetPortRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyTargetPortRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyTargetPortResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyTargetPortResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyTargetPortResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyTargetWeightRequest struct {
	*tchttp.BaseRequest

	// 负载均衡实例 ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 负载均衡监听器 ID
	ListenerId *string `json:"ListenerId" name:"ListenerId"`

	// 后端云服务器新的转发权重，取值范围：0~100。
	Weight *int64 `json:"Weight" name:"Weight"`

	// 转发规则的ID，当绑定机器到七层转发规则时，必须提供此参数或Domain+Url两者之一
	LocationId *string `json:"LocationId" name:"LocationId"`

	// 目标规则的域名，提供LocationId参数时本参数不生效
	Domain *string `json:"Domain" name:"Domain"`

	// 目标规则的URL，提供LocationId参数时本参数不生效
	Url *string `json:"Url" name:"Url"`

	// 要修改权重的后端机器列表
	Targets []*Target `json:"Targets" name:"Targets" list`
}

func (r *ModifyTargetWeightRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyTargetWeightRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyTargetWeightResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyTargetWeightResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyTargetWeightResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RegisterTargetsRequest struct {
	*tchttp.BaseRequest

	// 负载均衡实例 ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 负载均衡监听器 ID
	ListenerId *string `json:"ListenerId" name:"ListenerId"`

	// 要注册的后端机器列表
	Targets []*Target `json:"Targets" name:"Targets" list`

	// 转发规则的ID，当注册机器到七层转发规则时，必须提供此参数或Domain+Url两者之一
	LocationId *string `json:"LocationId" name:"LocationId"`

	// 目标规则的域名，提供LocationId参数时本参数不生效
	Domain *string `json:"Domain" name:"Domain"`

	// 目标规则的URL，提供LocationId参数时本参数不生效
	Url *string `json:"Url" name:"Url"`
}

func (r *RegisterTargetsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RegisterTargetsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RegisterTargetsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RegisterTargetsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RegisterTargetsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RuleInput struct {

	// 转发规则的域名。
	Domain *string `json:"Domain" name:"Domain"`

	// 转发规则的路径。
	Url *string `json:"Url" name:"Url"`

	// 会话保持时间
	SessionExpireTime *int64 `json:"SessionExpireTime" name:"SessionExpireTime"`

	// 健康检查信息
	HealthCheck *HealthCheck `json:"HealthCheck" name:"HealthCheck"`

	// 证书信息
	Certificate *CertificateInput `json:"Certificate" name:"Certificate"`

	// 规则的请求转发方式
	Scheduler *string `json:"Scheduler" name:"Scheduler"`
}

type RuleOutput struct {

	// 转发规则的 ID，作为输入时无需此字段
	LocationId *string `json:"LocationId" name:"LocationId"`

	// 转发规则的域名。
	Domain *string `json:"Domain" name:"Domain"`

	// 转发规则的路径。
	Url *string `json:"Url" name:"Url"`

	// 会话保持时间
	SessionExpireTime *int64 `json:"SessionExpireTime" name:"SessionExpireTime"`

	// 健康检查信息
	HealthCheck *HealthCheck `json:"HealthCheck" name:"HealthCheck"`

	// 证书信息
	Certificate *CertificateOutput `json:"Certificate" name:"Certificate"`

	// 规则的请求转发方式
	Scheduler *string `json:"Scheduler" name:"Scheduler"`
}

type RuleTargets struct {

	// 转发规则的 ID
	LocationId *string `json:"LocationId" name:"LocationId"`

	// 转发规则的域名
	Domain *string `json:"Domain" name:"Domain"`

	// 转发规则的路径。
	Url *string `json:"Url" name:"Url"`

	// 后端机器的信息
	Targets []*Backend `json:"Targets" name:"Targets" list`
}

type Target struct {

	// 云服务器的唯一 ID，可通过 DescribeInstances 接口返回字段中的 unInstanceId 字段获取
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 后端云服务器监听端口
	Port *int64 `json:"Port" name:"Port"`

	// 转发目标的类型，目前仅可取值为 CVM
	Type *string `json:"Type" name:"Type"`

	// 后端云服务器的转发权重，取值范围：0~100，默认为 10。
	Weight *int64 `json:"Weight" name:"Weight"`
}
