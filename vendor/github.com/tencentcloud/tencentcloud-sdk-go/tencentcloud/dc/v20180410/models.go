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

package v20180410

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type AcceptDirectConnectTunnelRequest struct {
	*tchttp.BaseRequest

	// 物理专线拥有者接受共享专用通道申请
	DirectConnectTunnelId *string `json:"DirectConnectTunnelId" name:"DirectConnectTunnelId"`
}

func (r *AcceptDirectConnectTunnelRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AcceptDirectConnectTunnelRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AcceptDirectConnectTunnelResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AcceptDirectConnectTunnelResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AcceptDirectConnectTunnelResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type BgpPeer struct {

	// 用户侧，BGP Asn
	Asn *int64 `json:"Asn" name:"Asn"`

	// 用户侧BGP密钥
	AuthKey *string `json:"AuthKey" name:"AuthKey"`
}

type CreateDirectConnectTunnelRequest struct {
	*tchttp.BaseRequest

	// 专线 ID，例如：dc-kd7d06of
	DirectConnectId *string `json:"DirectConnectId" name:"DirectConnectId"`

	// 专用通道名称
	DirectConnectTunnelName *string `json:"DirectConnectTunnelName" name:"DirectConnectTunnelName"`

	// 物理专线 owner，缺省为当前客户（物理专线 owner）
	// 共享专线时这里需要填写共享专线的开发商账号 ID
	DirectConnectOwnerAccount *string `json:"DirectConnectOwnerAccount" name:"DirectConnectOwnerAccount"`

	// 网络类型，分别为VPC、BMVPC，CCN，默认是VPC
	// VPC：私有网络
	// BMVPC：黑石网络
	// CCN：云联网
	NetworkType *string `json:"NetworkType" name:"NetworkType"`

	// 网络地域
	NetworkRegion *string `json:"NetworkRegion" name:"NetworkRegion"`

	// 私有网络统一 ID 或者黑石网络统一 ID
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 专线网关 ID，例如 dcg-d545ddf
	DirectConnectGatewayId *string `json:"DirectConnectGatewayId" name:"DirectConnectGatewayId"`

	// 专线带宽，单位：Mbps
	// 默认是物理专线带宽值
	Bandwidth *int64 `json:"Bandwidth" name:"Bandwidth"`

	// BGP ：BGP路由
	// STATIC：静态
	// 默认为 BGP 路由
	RouteType *string `json:"RouteType" name:"RouteType"`

	// BgpPeer，用户侧bgp信息，包括Asn和AuthKey
	BgpPeer *BgpPeer `json:"BgpPeer" name:"BgpPeer"`

	// 静态路由，用户IDC的网段地址
	RouteFilterPrefixes []*RouteFilterPrefix `json:"RouteFilterPrefixes" name:"RouteFilterPrefixes" list`

	// vlan，范围：0 ~ 3000
	// 0：不开启子接口
	// 默认值是非0
	Vlan *int64 `json:"Vlan" name:"Vlan"`

	// TencentAddress，腾讯侧互联 IP
	TencentAddress *string `json:"TencentAddress" name:"TencentAddress"`

	// CustomerAddress，用户侧互联 IP
	CustomerAddress *string `json:"CustomerAddress" name:"CustomerAddress"`
}

func (r *CreateDirectConnectTunnelRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDirectConnectTunnelRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDirectConnectTunnelResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 专用通道ID
		DirectConnectTunnelIdSet []*string `json:"DirectConnectTunnelIdSet" name:"DirectConnectTunnelIdSet" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateDirectConnectTunnelResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDirectConnectTunnelResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteDirectConnectTunnelRequest struct {
	*tchttp.BaseRequest

	// 专用通道ID
	DirectConnectTunnelId *string `json:"DirectConnectTunnelId" name:"DirectConnectTunnelId"`
}

func (r *DeleteDirectConnectTunnelRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteDirectConnectTunnelRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteDirectConnectTunnelResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteDirectConnectTunnelResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteDirectConnectTunnelResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDirectConnectTunnelsRequest struct {
	*tchttp.BaseRequest

	// 过滤条件:
	// 参数不支持同时指定DirectConnectTunnelIds和Filters。
	// <li> direct-connect-tunnel-name, 专用通道名称。</li>
	// <li> direct-connect-tunnel-id, 专用通道实例ID，如dcx-abcdefgh。</li>
	// <li>direct-connect-id, 物理专线实例ID，如，dc-abcdefgh。</li>
	Filters []*Filter `json:"Filters" name:"Filters" list`

	// 专用通道 ID数组
	DirectConnectTunnelIds []*string `json:"DirectConnectTunnelIds" name:"DirectConnectTunnelIds" list`

	// 偏移量，默认为0
	Offset *int64 `json:"Offset" name:"Offset"`

	// 返回数量，默认为20，最大值为100
	Limit *int64 `json:"Limit" name:"Limit"`
}

func (r *DescribeDirectConnectTunnelsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDirectConnectTunnelsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDirectConnectTunnelsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 专用通道列表
		DirectConnectTunnelSet []*DirectConnectTunnel `json:"DirectConnectTunnelSet" name:"DirectConnectTunnelSet" list`

		// 符合专用通道数量。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDirectConnectTunnelsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDirectConnectTunnelsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DirectConnectTunnel struct {

	// 专线通道ID
	DirectConnectTunnelId *string `json:"DirectConnectTunnelId" name:"DirectConnectTunnelId"`

	// 物理专线ID
	DirectConnectId *string `json:"DirectConnectId" name:"DirectConnectId"`

	// 专线通道状态
	// AVAILABLE:就绪或者已连接
	// PENDING:申请中
	// ALLOCATING:配置中
	// ALLOCATED:配置完成
	// ALTERING:修改中
	// DELETING:删除中
	// DELETED:删除完成
	// COMFIRMING:待接受
	// REJECTED:拒绝
	State *string `json:"State" name:"State"`

	// 物理专线的拥有者，开发商账号 ID
	DirectConnectOwnerAccount *string `json:"DirectConnectOwnerAccount" name:"DirectConnectOwnerAccount"`

	// 专线通道的拥有者，开发商账号 ID
	OwnerAccount *string `json:"OwnerAccount" name:"OwnerAccount"`

	// 网络类型，分别为VPC、BMVPC、CCN
	//  VPC：私有网络 ，BMVPC：黑石网络，CCN：云联网
	NetworkType *string `json:"NetworkType" name:"NetworkType"`

	// VPC地域
	NetworkRegion *string `json:"NetworkRegion" name:"NetworkRegion"`

	// 私有网络统一 ID 或者黑石网络统一 ID
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 专线网关 ID
	DirectConnectGatewayId *string `json:"DirectConnectGatewayId" name:"DirectConnectGatewayId"`

	// BGP ：BGP路由 STATIC：静态 默认为 BGP 路由
	RouteType *string `json:"RouteType" name:"RouteType"`

	// 用户侧BGP，Asn，AuthKey
	BgpPeer *BgpPeer `json:"BgpPeer" name:"BgpPeer"`

	// 用户侧网段地址
	RouteFilterPrefixes []*RouteFilterPrefix `json:"RouteFilterPrefixes" name:"RouteFilterPrefixes" list`

	// 专线通道的Vlan
	Vlan *int64 `json:"Vlan" name:"Vlan"`

	// TencentAddress，腾讯侧互联 IP
	TencentAddress *string `json:"TencentAddress" name:"TencentAddress"`

	// CustomerAddress，用户侧互联 IP
	CustomerAddress *string `json:"CustomerAddress" name:"CustomerAddress"`

	// 专线通道名称
	DirectConnectTunnelName *string `json:"DirectConnectTunnelName" name:"DirectConnectTunnelName"`

	// 专线通道创建时间
	CreatedTime *string `json:"CreatedTime" name:"CreatedTime"`

	// 专线通道带宽值
	Bandwidth *int64 `json:"Bandwidth" name:"Bandwidth"`
}

type Filter struct {

	// 需要过滤的字段。
	Name *string `json:"Name" name:"Name"`

	// 字段的过滤值。
	Values []*string `json:"Values" name:"Values" list`
}

type ModifyDirectConnectTunnelAttributeRequest struct {
	*tchttp.BaseRequest

	// 专用通道ID
	DirectConnectTunnelId *string `json:"DirectConnectTunnelId" name:"DirectConnectTunnelId"`

	// 专用通道名称
	DirectConnectTunnelName *string `json:"DirectConnectTunnelName" name:"DirectConnectTunnelName"`

	// 用户侧BGP，包括Asn，AuthKey
	BgpPeer *BgpPeer `json:"BgpPeer" name:"BgpPeer"`

	// 用户侧网段地址
	RouteFilterPrefixes []*RouteFilterPrefix `json:"RouteFilterPrefixes" name:"RouteFilterPrefixes" list`

	// 腾讯侧互联IP
	TencentAddress *string `json:"TencentAddress" name:"TencentAddress"`

	// 用户侧互联IP
	CustomerAddress *string `json:"CustomerAddress" name:"CustomerAddress"`

	// 专用通道带宽值，单位为M。
	Bandwidth *int64 `json:"Bandwidth" name:"Bandwidth"`
}

func (r *ModifyDirectConnectTunnelAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDirectConnectTunnelAttributeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDirectConnectTunnelAttributeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDirectConnectTunnelAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDirectConnectTunnelAttributeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RejectDirectConnectTunnelRequest struct {
	*tchttp.BaseRequest

	// 无
	DirectConnectTunnelId *string `json:"DirectConnectTunnelId" name:"DirectConnectTunnelId"`
}

func (r *RejectDirectConnectTunnelRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RejectDirectConnectTunnelRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RejectDirectConnectTunnelResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RejectDirectConnectTunnelResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RejectDirectConnectTunnelResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RouteFilterPrefix struct {

	// 用户侧网段地址
	Cidr *string `json:"Cidr" name:"Cidr"`
}
