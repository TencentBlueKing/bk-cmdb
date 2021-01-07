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

package v20180419

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type AttachInstancesRequest struct {
	*tchttp.BaseRequest

	// 伸缩组ID
	AutoScalingGroupId *string `json:"AutoScalingGroupId" name:"AutoScalingGroupId"`

	// CVM实例ID列表
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`
}

func (r *AttachInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AttachInstancesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AttachInstancesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AttachInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AttachInstancesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AutoScalingGroup struct {

	// 伸缩组ID
	AutoScalingGroupId *string `json:"AutoScalingGroupId" name:"AutoScalingGroupId"`

	// 伸缩组名称
	AutoScalingGroupName *string `json:"AutoScalingGroupName" name:"AutoScalingGroupName"`

	// 伸缩组状态
	AutoScalingGroupStatus *string `json:"AutoScalingGroupStatus" name:"AutoScalingGroupStatus"`

	// 创建时间，采用UTC标准计时
	CreatedTime *string `json:"CreatedTime" name:"CreatedTime"`

	// 默认冷却时间，单位秒
	DefaultCooldown *uint64 `json:"DefaultCooldown" name:"DefaultCooldown"`

	// 期望实例数
	DesiredCapacity *uint64 `json:"DesiredCapacity" name:"DesiredCapacity"`

	// 启用状态，取值包括`ENABLED`和`DISABLED`
	EnabledStatus *string `json:"EnabledStatus" name:"EnabledStatus"`

	// 应用型负载均衡器列表
	ForwardLoadBalancerSet []*ForwardLoadBalancer `json:"ForwardLoadBalancerSet" name:"ForwardLoadBalancerSet" list`

	// 实例数量
	InstanceCount *uint64 `json:"InstanceCount" name:"InstanceCount"`

	// 状态为`IN_SERVICE`实例的数量
	InServiceInstanceCount *uint64 `json:"InServiceInstanceCount" name:"InServiceInstanceCount"`

	// 启动配置ID
	LaunchConfigurationId *string `json:"LaunchConfigurationId" name:"LaunchConfigurationId"`

	// 启动配置名称
	LaunchConfigurationName *string `json:"LaunchConfigurationName" name:"LaunchConfigurationName"`

	// 传统型负载均衡器ID列表
	LoadBalancerIdSet []*string `json:"LoadBalancerIdSet" name:"LoadBalancerIdSet" list`

	// 最大实例数
	MaxSize []*uint64 `json:"MaxSize" name:"MaxSize" list`

	// 最小实例数
	MinSize []*uint64 `json:"MinSize" name:"MinSize" list`

	// 项目ID
	ProjectId []*uint64 `json:"ProjectId" name:"ProjectId" list`

	// 子网ID列表
	SubnetIdSet []*string `json:"SubnetIdSet" name:"SubnetIdSet" list`

	// 销毁策略
	TerminationPolicySet []*string `json:"TerminationPolicySet" name:"TerminationPolicySet" list`

	// VPC标识
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 可用区列表
	ZoneSet []*string `json:"ZoneSet" name:"ZoneSet" list`

	// 重试策略
	RetryPolicy *string `json:"RetryPolicy" name:"RetryPolicy"`
}

type AutoScalingGroupAbstract struct {

	// 伸缩组ID。
	AutoScalingGroupId *string `json:"AutoScalingGroupId" name:"AutoScalingGroupId"`

	// 伸缩组名称。
	AutoScalingGroupName *string `json:"AutoScalingGroupName" name:"AutoScalingGroupName"`
}

type CreateAutoScalingGroupRequest struct {
	*tchttp.BaseRequest

	// 伸缩组名称，在您账号中必须唯一。名称仅支持中文、英文、数字、下划线、分隔符"-"、小数点，最大长度不能超55个字节。
	AutoScalingGroupName *string `json:"AutoScalingGroupName" name:"AutoScalingGroupName"`

	// 启动配置ID
	LaunchConfigurationId *string `json:"LaunchConfigurationId" name:"LaunchConfigurationId"`

	// 最大实例数，取值范围为0-2000。
	MaxSize *uint64 `json:"MaxSize" name:"MaxSize"`

	// 最小实例数，取值范围为0-2000。
	MinSize *uint64 `json:"MinSize" name:"MinSize"`

	// VPC ID，基础网络则填空字符串
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 默认冷却时间，单位秒，默认值为300
	DefaultCooldown *uint64 `json:"DefaultCooldown" name:"DefaultCooldown"`

	// 期望实例数，大小介于最小实例数和最大实例数之间
	DesiredCapacity *uint64 `json:"DesiredCapacity" name:"DesiredCapacity"`

	// 传统负载均衡器ID列表，目前长度上限为1，LoadBalancerIds 和 ForwardLoadBalancers 二者同时最多只能指定一个
	LoadBalancerIds []*string `json:"LoadBalancerIds" name:"LoadBalancerIds" list`

	// 项目ID
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`

	// 应用型负载均衡器列表，目前长度上限为1，LoadBalancerIds 和 ForwardLoadBalancers 二者同时最多只能指定一个
	ForwardLoadBalancers []*ForwardLoadBalancer `json:"ForwardLoadBalancers" name:"ForwardLoadBalancers" list`

	// 子网ID列表，VPC场景下必须指定子网
	SubnetIds []*string `json:"SubnetIds" name:"SubnetIds" list`

	// 销毁策略，目前长度上限为1。取值包括 OLDEST_INSTANCE 和 NEWEST_INSTANCE，默认取值为 OLDEST_INSTANCE。
	// <br><li> OLDEST_INSTANCE 优先销毁伸缩组中最旧的实例。
	// <br><li> NEWEST_INSTANCE，优先销毁伸缩组中最新的实例。
	TerminationPolicies []*string `json:"TerminationPolicies" name:"TerminationPolicies" list`

	// 可用区列表，基础网络场景下必须指定可用区
	Zones []*string `json:"Zones" name:"Zones" list`

	// 重试策略，取值包括 IMMEDIATE_RETRY 和 INCREMENTAL_INTERVALS，默认取值为 IMMEDIATE_RETRY。
	// <br><li> IMMEDIATE_RETRY，立即重试，在较短时间内快速重试，连续失败超过一定次数（5次）后不再重试。
	// <br><li> INCREMENTAL_INTERVALS，间隔递增重试，随着连续失败次数的增加，重试间隔逐渐增大，重试间隔从秒级到1天不等。在连续失败超过一定次数（25次）后不再重试。
	RetryPolicy *string `json:"RetryPolicy" name:"RetryPolicy"`

	// 可用区校验策略，取值包括 ALL 和 ANY，默认取值为ANY。
	// <br><li> ALL，所有可用区（Zone）或子网（SubnetId）都可用则通过校验，否则校验报错。
	// <br><li> ANY，存在任何一个可用区（Zone）或子网（SubnetId）可用则通过校验，否则校验报错。
	// 
	// 可用区或子网不可用的常见原因包括该可用区CVM实例类型售罄、该可用区CBS云盘售罄、该可用区配额不足、该子网IP不足等。
	// 如果 Zones/SubnetIds 中可用区或者子网不存在，则无论 ZonesCheckPolicy 采用何种取值，都会校验报错。
	ZonesCheckPolicy *string `json:"ZonesCheckPolicy" name:"ZonesCheckPolicy"`
}

func (r *CreateAutoScalingGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateAutoScalingGroupRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateAutoScalingGroupResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 伸缩组ID
		AutoScalingGroupId *string `json:"AutoScalingGroupId" name:"AutoScalingGroupId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateAutoScalingGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateAutoScalingGroupResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateLaunchConfigurationRequest struct {
	*tchttp.BaseRequest

	// 启动配置显示名称。名称仅支持中文、英文、数字、下划线、分隔符"-"、小数点，最大长度不能超60个字节。
	LaunchConfigurationName *string `json:"LaunchConfigurationName" name:"LaunchConfigurationName"`

	// 指定有效的[镜像](https://cloud.tencent.com/document/product/213/4940)ID，格式形如`img-8toqc6s3`。镜像类型分为四种：<br/><li>公共镜像</li><li>自定义镜像</li><li>共享镜像</li><li>服务市场镜像</li><br/>可通过以下方式获取可用的镜像ID：<br/><li>`公共镜像`、`自定义镜像`、`共享镜像`的镜像ID可通过登录[控制台](https://console.cloud.tencent.com/cvm/image?rid=1&imageType=PUBLIC_IMAGE)查询；`服务镜像市场`的镜像ID可通过[云市场](https://market.cloud.tencent.com/list)查询。</li><li>通过调用接口 [DescribeImages](https://cloud.tencent.com/document/api/213/15715) ，取返回信息中的`ImageId`字段。</li>
	ImageId *string `json:"ImageId" name:"ImageId"`

	// 实例所属项目ID。该参数可以通过调用 [DescribeProject](https://cloud.tencent.com/document/api/378/4400) 的返回值中的`projectId`字段来获取。不填为默认项目。
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`

	// 实例机型。不同实例机型指定了不同的资源规格，具体取值可通过调用接口 [DescribeInstanceTypeConfigs](https://cloud.tencent.com/document/api/213/15749) 来获得最新的规格表或参见[实例类型](https://cloud.tencent.com/document/product/213/11518)描述。
	// `InstanceType`和`InstanceTypes`参数互斥，二者必填一个且只能填写一个。
	InstanceType *string `json:"InstanceType" name:"InstanceType"`

	// 实例系统盘配置信息。若不指定该参数，则按照系统默认值进行分配。
	SystemDisk *SystemDisk `json:"SystemDisk" name:"SystemDisk"`

	// 实例数据盘配置信息。若不指定该参数，则默认不购买数据盘，最多支持指定11块数据盘。
	DataDisks []*DataDisk `json:"DataDisks" name:"DataDisks" list`

	// 公网带宽相关信息设置。若不指定该参数，则默认公网带宽为0Mbps。
	InternetAccessible *InternetAccessible `json:"InternetAccessible" name:"InternetAccessible"`

	// 实例登录设置。通过该参数可以设置实例的登录方式密码、密钥或保持镜像的原始登录设置。默认情况下会随机生成密码，并以站内信方式知会到用户。
	LoginSettings *LoginSettings `json:"LoginSettings" name:"LoginSettings"`

	// 实例所属安全组。该参数可以通过调用 [DescribeSecurityGroups](https://cloud.tencent.com/document/api/215/15808) 的返回值中的`SecurityGroupId`字段来获取。若不指定该参数，则默认不绑定安全组。
	SecurityGroupIds []*string `json:"SecurityGroupIds" name:"SecurityGroupIds" list`

	// 增强服务。通过该参数可以指定是否开启云安全、云监控等服务。若不指定该参数，则默认开启云监控、云安全服务。
	EnhancedService *EnhancedService `json:"EnhancedService" name:"EnhancedService"`

	// 经过 Base64 编码后的自定义数据，最大长度不超过16KB。
	UserData *string `json:"UserData" name:"UserData"`

	// 实例计费类型，CVM默认值按照POSTPAID_BY_HOUR处理。
	// <br><li>POSTPAID_BY_HOUR：按小时后付费
	// <br><li>SPOTPAID：竞价付费
	InstanceChargeType *string `json:"InstanceChargeType" name:"InstanceChargeType"`

	// 实例的市场相关选项，如竞价实例相关参数，若指定实例的付费模式为竞价付费则该参数必传。
	InstanceMarketOptions *InstanceMarketOptionsRequest `json:"InstanceMarketOptions" name:"InstanceMarketOptions"`

	// 实例机型列表，不同实例机型指定了不同的资源规格，最多支持5中实例机型。
	// `InstanceType`和`InstanceTypes`参数互斥，二者必填一个且只能填写一个。
	InstanceTypes []*string `json:"InstanceTypes" name:"InstanceTypes" list`

	// 实例类型校验策略，取值包括 ALL 和 ANY，默认取值为ANY。
	// <br><li> ALL，所有实例类型（InstanceType）都可用则通过校验，否则校验报错。
	// <br><li> ANY，存在任何一个实例类型（InstanceType）可用则通过校验，否则校验报错。
	// 
	// 实例类型不可用的常见原因包括该实例类型售罄、对应云盘售罄等。
	// 如果 InstanceTypes 中一款机型不存在或者已下线，则无论 InstanceTypesCheckPolicy 采用何种取值，都会校验报错。
	InstanceTypesCheckPolicy *string `json:"InstanceTypesCheckPolicy" name:"InstanceTypesCheckPolicy"`
}

func (r *CreateLaunchConfigurationRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateLaunchConfigurationRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateLaunchConfigurationResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 当通过本接口来创建启动配置时会返回该参数，表示启动配置ID。
		LaunchConfigurationId *string `json:"LaunchConfigurationId" name:"LaunchConfigurationId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateLaunchConfigurationResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateLaunchConfigurationResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateScheduledActionRequest struct {
	*tchttp.BaseRequest

	// 伸缩组ID
	AutoScalingGroupId *string `json:"AutoScalingGroupId" name:"AutoScalingGroupId"`

	// 定时任务名称。名称仅支持中文、英文、数字、下划线、分隔符"-"、小数点，最大长度不能超60个字节。同一伸缩组下必须唯一。
	ScheduledActionName *string `json:"ScheduledActionName" name:"ScheduledActionName"`

	// 当定时任务触发时，设置的伸缩组最大实例数。
	MaxSize *uint64 `json:"MaxSize" name:"MaxSize"`

	// 当定时任务触发时，设置的伸缩组最小实例数。
	MinSize *uint64 `json:"MinSize" name:"MinSize"`

	// 当定时任务触发时，设置的伸缩组期望实例数。
	DesiredCapacity *uint64 `json:"DesiredCapacity" name:"DesiredCapacity"`

	// 定时任务的首次触发时间，取值为`北京时间`（UTC+8），按照`ISO8601`标准，格式：`YYYY-MM-DDThh:mm:ss+08:00`。
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 定时任务的结束时间，取值为`北京时间`（UTC+8），按照`ISO8601`标准，格式：`YYYY-MM-DDThh:mm:ss+08:00`。<br><br>此参数与`Recurrence`需要同时指定，到达结束时间之后，定时任务将不再生效。
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 定时任务的重复方式。为标准[Cron](https://zh.wikipedia.org/wiki/Cron)格式<br><br>此参数与`EndTime`需要同时指定。
	Recurrence *string `json:"Recurrence" name:"Recurrence"`
}

func (r *CreateScheduledActionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateScheduledActionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateScheduledActionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 定时任务ID
		ScheduledActionId *string `json:"ScheduledActionId" name:"ScheduledActionId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateScheduledActionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateScheduledActionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DataDisk struct {

	// 数据盘类型。数据盘类型限制详见[CVM实例配置](https://cloud.tencent.com/document/product/213/2177)。取值范围：<br><li>LOCAL_BASIC：本地硬盘<br><li>LOCAL_SSD：本地SSD硬盘<br><li>CLOUD_BASIC：普通云硬盘<br><li>CLOUD_PREMIUM：高性能云硬盘<br><li>CLOUD_SSD：SSD云硬盘<br><br>默认取值：LOCAL_BASIC。
	DiskType *string `json:"DiskType" name:"DiskType"`

	// 数据盘大小，单位：GB。最小调整步长为10G，不同数据盘类型取值范围不同，具体限制详见：[CVM实例配置](https://cloud.tencent.com/document/product/213/2177)。默认值为0，表示不购买数据盘。更多限制详见产品文档。
	DiskSize *uint64 `json:"DiskSize" name:"DiskSize"`
}

type DeleteAutoScalingGroupRequest struct {
	*tchttp.BaseRequest

	// 伸缩组ID
	AutoScalingGroupId *string `json:"AutoScalingGroupId" name:"AutoScalingGroupId"`
}

func (r *DeleteAutoScalingGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteAutoScalingGroupRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteAutoScalingGroupResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteAutoScalingGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteAutoScalingGroupResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteLaunchConfigurationRequest struct {
	*tchttp.BaseRequest

	// 需要删除的启动配置ID。
	LaunchConfigurationId *string `json:"LaunchConfigurationId" name:"LaunchConfigurationId"`
}

func (r *DeleteLaunchConfigurationRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteLaunchConfigurationRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteLaunchConfigurationResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteLaunchConfigurationResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteLaunchConfigurationResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteScheduledActionRequest struct {
	*tchttp.BaseRequest

	// 待删除的定时任务ID。
	ScheduledActionId *string `json:"ScheduledActionId" name:"ScheduledActionId"`
}

func (r *DeleteScheduledActionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteScheduledActionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteScheduledActionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteScheduledActionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteScheduledActionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAccountLimitsRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeAccountLimitsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAccountLimitsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAccountLimitsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 用户账户被允许创建的启动配置最大数量
		MaxNumberOfLaunchConfigurations *int64 `json:"MaxNumberOfLaunchConfigurations" name:"MaxNumberOfLaunchConfigurations"`

		// 用户账户启动配置的当前数量
		NumberOfLaunchConfigurations *int64 `json:"NumberOfLaunchConfigurations" name:"NumberOfLaunchConfigurations"`

		// 用户账户被允许创建的伸缩组最大数量
		MaxNumberOfAutoScalingGroups *int64 `json:"MaxNumberOfAutoScalingGroups" name:"MaxNumberOfAutoScalingGroups"`

		// 用户账户伸缩组的当前数量
		NumberOfAutoScalingGroups *int64 `json:"NumberOfAutoScalingGroups" name:"NumberOfAutoScalingGroups"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeAccountLimitsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAccountLimitsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAutoScalingGroupsRequest struct {
	*tchttp.BaseRequest

	// 按照一个或者多个伸缩组ID查询。伸缩组ID形如：`asg-nkdwoui0`。每次请求的上限为100。参数不支持同时指定`AutoScalingGroups`和`Filters`。
	AutoScalingGroupIds []*string `json:"AutoScalingGroupIds" name:"AutoScalingGroupIds" list`

	// 过滤条件。
	// <li> auto-scaling-group-id - String - 是否必填：否 -（过滤条件）按照伸缩组ID过滤。</li>
	// <li> auto-scaling-group-name - String - 是否必填：否 -（过滤条件）按照伸缩组名称过滤。</li>
	// <li> launch-configuration-id - String - 是否必填：否 -（过滤条件）按照启动配置ID过滤。</li>
	// 每次请求的`Filters`的上限为10，`Filter.Values`的上限为5。参数不支持同时指定`AutoScalingGroupIds`和`Filters`。
	Filters []*Filter `json:"Filters" name:"Filters" list`

	// 返回数量，默认为20，最大值为100。关于`Limit`的更进一步介绍请参考 API [简介](https://cloud.tencent.com/document/api/213/15688)中的相关小节。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。关于`Offset`的更进一步介绍请参考 API [简介](https://cloud.tencent.com/document/api/213/15688)中的相关小节。
	Offset *uint64 `json:"Offset" name:"Offset"`
}

func (r *DescribeAutoScalingGroupsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAutoScalingGroupsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAutoScalingGroupsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 伸缩组详细信息列表。
		AutoScalingGroupSet []*AutoScalingGroup `json:"AutoScalingGroupSet" name:"AutoScalingGroupSet" list`

		// 符合条件的伸缩组数量。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeAutoScalingGroupsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAutoScalingGroupsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAutoScalingInstancesRequest struct {
	*tchttp.BaseRequest

	// 待查询的云主机（CVM）实例ID。参数不支持同时指定InstanceIds和Filters。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`

	// 过滤条件。
	// <li> instance-id - String - 是否必填：否 -（过滤条件）按照实例ID过滤。</li>
	// <li> auto-scaling-group-id - String - 是否必填：否 -（过滤条件）按照伸缩组ID过滤。</li>
	// 每次请求的`Filters`的上限为10，`Filter.Values`的上限为5。参数不支持同时指定`InstanceIds`和`Filters`。
	Filters []*Filter `json:"Filters" name:"Filters" list`

	// 偏移量，默认为0。关于`Offset`的更进一步介绍请参考 API [简介](https://cloud.tencent.com/document/api/213/15688)中的相关小节。
	Offset *int64 `json:"Offset" name:"Offset"`

	// 返回数量，默认为20，最大值为100。关于`Limit`的更进一步介绍请参考 API [简介](https://cloud.tencent.com/document/api/213/15688)中的相关小节。
	Limit *int64 `json:"Limit" name:"Limit"`
}

func (r *DescribeAutoScalingInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAutoScalingInstancesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAutoScalingInstancesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 实例详细信息列表。
		AutoScalingInstanceSet []*Instance `json:"AutoScalingInstanceSet" name:"AutoScalingInstanceSet" list`

		// 符合条件的实例数量。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeAutoScalingInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAutoScalingInstancesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLaunchConfigurationsRequest struct {
	*tchttp.BaseRequest

	// 按照一个或者多个启动配置ID查询。启动配置ID形如：`asc-ouy1ax38`。每次请求的上限为100。参数不支持同时指定`LaunchConfigurationIds`和`Filters`。
	LaunchConfigurationIds []*string `json:"LaunchConfigurationIds" name:"LaunchConfigurationIds" list`

	// 过滤条件。
	// <li> launch-configuration-id - String - 是否必填：否 -（过滤条件）按照启动配置ID过滤。</li>
	// <li> launch-configuration-name - String - 是否必填：否 -（过滤条件）按照启动配置名称过滤。</li>
	// 每次请求的`Filters`的上限为10，`Filter.Values`的上限为5。参数不支持同时指定`LaunchConfigurationIds`和`Filters`。
	Filters []*Filter `json:"Filters" name:"Filters" list`

	// 返回数量，默认为20，最大值为100。关于`Limit`的更进一步介绍请参考 API [简介](https://cloud.tencent.com/document/api/213/15688)中的相关小节。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 偏移量，默认为0。关于`Offset`的更进一步介绍请参考 API [简介](https://cloud.tencent.com/document/api/213/15688)中的相关小节。
	Offset *uint64 `json:"Offset" name:"Offset"`
}

func (r *DescribeLaunchConfigurationsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLaunchConfigurationsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLaunchConfigurationsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合条件的启动配置数量。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 启动配置详细信息列表。
		LaunchConfigurationSet []*LaunchConfiguration `json:"LaunchConfigurationSet" name:"LaunchConfigurationSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeLaunchConfigurationsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLaunchConfigurationsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeScheduledActionsRequest struct {
	*tchttp.BaseRequest

	// 按照一个或者多个定时任务ID查询。实例ID形如：asst-am691zxo。每次请求的实例的上限为100。参数不支持同时指定ScheduledActionIds和Filters。
	ScheduledActionIds []*string `json:"ScheduledActionIds" name:"ScheduledActionIds" list`

	// 过滤条件。
	// <li> scheduled-action-id - String - 是否必填：否 -（过滤条件）按照定时任务ID过滤。</li>
	// <li> scheduled-action-name - String - 是否必填：否 - （过滤条件） 按照定时任务名称过滤。</li>
	// <li> auto-scaling-group-id - String - 是否必填：否 - （过滤条件） 按照伸缩组ID过滤。</li>
	Filters []*Filter `json:"Filters" name:"Filters" list`

	// 偏移量，默认为0。关于Offset的更进一步介绍请参考 API [简介](https://cloud.tencent.com/document/api/213/15688)中的相关小节。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 返回数量，默认为20，最大值为100。关于Limit的更进一步介绍请参考 API [简介](https://cloud.tencent.com/document/api/213/15688)中的相关小节。
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribeScheduledActionsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeScheduledActionsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeScheduledActionsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合条件的定时任务数量。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 定时任务详细信息列表。
		ScheduledActionSet []*ScheduledAction `json:"ScheduledActionSet" name:"ScheduledActionSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeScheduledActionsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeScheduledActionsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DetachInstancesRequest struct {
	*tchttp.BaseRequest

	// 伸缩组ID
	AutoScalingGroupId *string `json:"AutoScalingGroupId" name:"AutoScalingGroupId"`

	// CVM实例ID列表
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`
}

func (r *DetachInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DetachInstancesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DetachInstancesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DetachInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DetachInstancesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DisableAutoScalingGroupRequest struct {
	*tchttp.BaseRequest

	// 伸缩组ID
	AutoScalingGroupId *string `json:"AutoScalingGroupId" name:"AutoScalingGroupId"`
}

func (r *DisableAutoScalingGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DisableAutoScalingGroupRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DisableAutoScalingGroupResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DisableAutoScalingGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DisableAutoScalingGroupResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type EnableAutoScalingGroupRequest struct {
	*tchttp.BaseRequest

	// 伸缩组ID
	AutoScalingGroupId *string `json:"AutoScalingGroupId" name:"AutoScalingGroupId"`
}

func (r *EnableAutoScalingGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *EnableAutoScalingGroupRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type EnableAutoScalingGroupResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *EnableAutoScalingGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *EnableAutoScalingGroupResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type EnhancedService struct {

	// 开启云安全服务。若不指定该参数，则默认开启云安全服务。
	SecurityService *RunSecurityServiceEnabled `json:"SecurityService" name:"SecurityService"`

	// 开启云监控服务。若不指定该参数，则默认开启云监控服务。
	MonitorService *RunMonitorServiceEnabled `json:"MonitorService" name:"MonitorService"`
}

type Filter struct {

	// 需要过滤的字段。
	Name *string `json:"Name" name:"Name"`

	// 字段的过滤值。
	Values []*string `json:"Values" name:"Values" list`
}

type ForwardLoadBalancer struct {

	// 负载均衡器ID
	LoadBalancerId *string `json:"LoadBalancerId" name:"LoadBalancerId"`

	// 应用型负载均衡监听器 ID
	ListenerId *string `json:"ListenerId" name:"ListenerId"`

	// 目标规则属性列表
	TargetAttributes []*TargetAttribute `json:"TargetAttributes" name:"TargetAttributes" list`

	// 转发规则ID
	LocationId *string `json:"LocationId" name:"LocationId"`
}

type Instance struct {

	// 实例ID
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 伸缩组ID
	AutoScalingGroupId *string `json:"AutoScalingGroupId" name:"AutoScalingGroupId"`

	// 启动配置ID
	LaunchConfigurationId *string `json:"LaunchConfigurationId" name:"LaunchConfigurationId"`

	// 启动配置名称
	LaunchConfigurationName *string `json:"LaunchConfigurationName" name:"LaunchConfigurationName"`

	// 生命周期状态，取值包括IN_SERVICE, CREATING, TERMINATING, ATTACHING, DETACHING, ATTACHING_LB, DETACHING_LB等
	LifeCycleState *string `json:"LifeCycleState" name:"LifeCycleState"`

	// 健康状态，取值包括HEALTHY和UNHEALTHY
	HealthStatus *string `json:"HealthStatus" name:"HealthStatus"`

	// 是否加入缩容保护
	ProtectedFromScaleIn *bool `json:"ProtectedFromScaleIn" name:"ProtectedFromScaleIn"`

	// 可用区
	Zone *string `json:"Zone" name:"Zone"`

	// 创建类型，取值包括AUTO_CREATION, MANUAL_ATTACHING。
	CreationType *string `json:"CreationType" name:"CreationType"`

	// 实例加入时间
	AddTime *string `json:"AddTime" name:"AddTime"`

	// 实例类型
	InstanceType *string `json:"InstanceType" name:"InstanceType"`
}

type InstanceMarketOptionsRequest struct {
	*tchttp.BaseRequest

	// 竞价相关选项
	SpotOptions *SpotMarketOptions `json:"SpotOptions" name:"SpotOptions"`

	// 市场选项类型，当前只支持取值：spot
	MarketType *string `json:"MarketType" name:"MarketType"`
}

func (r *InstanceMarketOptionsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InstanceMarketOptionsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InternetAccessible struct {

	// 网络计费类型。取值范围：<br><li>BANDWIDTH_PREPAID：预付费按带宽结算<br><li>TRAFFIC_POSTPAID_BY_HOUR：流量按小时后付费<br><li>BANDWIDTH_POSTPAID_BY_HOUR：带宽按小时后付费<br><li>BANDWIDTH_PACKAGE：带宽包用户<br>默认取值：TRAFFIC_POSTPAID_BY_HOUR。
	InternetChargeType *string `json:"InternetChargeType" name:"InternetChargeType"`

	// 公网出带宽上限，单位：Mbps。默认值：0Mbps。不同机型带宽上限范围不一致，具体限制详见[购买网络带宽](https://cloud.tencent.com/document/product/213/509)。
	InternetMaxBandwidthOut *uint64 `json:"InternetMaxBandwidthOut" name:"InternetMaxBandwidthOut"`

	// 是否分配公网IP。取值范围：<br><li>TRUE：表示分配公网IP<br><li>FALSE：表示不分配公网IP<br><br>当公网带宽大于0Mbps时，可自由选择开通与否，默认开通公网IP；当公网带宽为0，则不允许分配公网IP。
	PublicIpAssigned *bool `json:"PublicIpAssigned" name:"PublicIpAssigned"`
}

type LaunchConfiguration struct {

	// 实例所属项目ID。
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 启动配置ID。
	LaunchConfigurationId *string `json:"LaunchConfigurationId" name:"LaunchConfigurationId"`

	// 启动配置名称。
	LaunchConfigurationName *string `json:"LaunchConfigurationName" name:"LaunchConfigurationName"`

	// 实例机型。
	InstanceType *string `json:"InstanceType" name:"InstanceType"`

	// 实例系统盘配置信息。
	SystemDisk *SystemDisk `json:"SystemDisk" name:"SystemDisk"`

	// 实例数据盘配置信息。
	DataDisks []*DataDisk `json:"DataDisks" name:"DataDisks" list`

	// 实例登录设置。
	LoginSettings *LimitedLoginSettings `json:"LoginSettings" name:"LoginSettings"`

	// 公网带宽相关信息设置。
	InternetAccessible *InternetAccessible `json:"InternetAccessible" name:"InternetAccessible"`

	// 实例所属安全组。
	SecurityGroupIds []*string `json:"SecurityGroupIds" name:"SecurityGroupIds" list`

	// 启动配置关联的伸缩组。
	AutoScalingGroupAbstractSet []*AutoScalingGroupAbstract `json:"AutoScalingGroupAbstractSet" name:"AutoScalingGroupAbstractSet" list`

	// 自定义数据。
	UserData *string `json:"UserData" name:"UserData"`

	// 启动配置创建时间。
	CreatedTime *string `json:"CreatedTime" name:"CreatedTime"`

	// 实例的增强服务启用情况与其设置。
	EnhancedService *EnhancedService `json:"EnhancedService" name:"EnhancedService"`

	// 镜像ID。
	ImageId *string `json:"ImageId" name:"ImageId"`

	// 启动配置当前状态。取值范围：<br><li>NORMAL：正常<br><li>IMAGE_ABNORMAL：启动配置镜像异常<br><li>CBS_SNAP_ABNORMAL：启动配置数据盘快照异常<br><li>SECURITY_GROUP_ABNORMAL：启动配置安全组异常<br>
	LaunchConfigurationStatus *string `json:"LaunchConfigurationStatus" name:"LaunchConfigurationStatus"`

	// 实例计费类型，CVM默认值按照POSTPAID_BY_HOUR处理。
	// <br><li>POSTPAID_BY_HOUR：按小时后付费
	// <br><li>SPOTPAID：竞价付费
	InstanceChargeType *string `json:"InstanceChargeType" name:"InstanceChargeType"`

	// 实例的市场相关选项，如竞价实例相关参数，若指定实例的付费模式为竞价付费则该参数必传。
	InstanceMarketOptions *InstanceMarketOptionsRequest `json:"InstanceMarketOptions" name:"InstanceMarketOptions"`

	// 实例机型列表。
	InstanceTypes []*string `json:"InstanceTypes" name:"InstanceTypes" list`
}

type LimitedLoginSettings struct {

	// 密钥ID列表。
	KeyIds []*string `json:"KeyIds" name:"KeyIds" list`
}

type LoginSettings struct {

	// 实例登录密码。不同操作系统类型密码复杂度限制不一样，具体如下：<br><li>Linux实例密码必须8到16位，至少包括两项[a-z，A-Z]、[0-9] 和 [( ) ` ~ ! @ # $ % ^ & * - + = | { } [ ] : ; ' , . ? / ]中的特殊符号。<br><li>Windows实例密码必须12到16位，至少包括三项[a-z]，[A-Z]，[0-9] 和 [( ) ` ~ ! @ # $ % ^ & * - + = { } [ ] : ; ' , . ? /]中的特殊符号。<br><br>若不指定该参数，则由系统随机生成密码，并通过站内信方式通知到用户。
	Password *string `json:"Password" name:"Password"`

	// 密钥ID列表。关联密钥后，就可以通过对应的私钥来访问实例；KeyId可通过接口DescribeKeyPairs获取，密钥与密码不能同时指定，同时Windows操作系统不支持指定密钥。当前仅支持购买的时候指定一个密钥。
	KeyIds []*string `json:"KeyIds" name:"KeyIds" list`

	// 保持镜像的原始设置。该参数与Password或KeyIds.N不能同时指定。只有使用自定义镜像、共享镜像或外部导入镜像创建实例时才能指定该参数为TRUE。取值范围：<br><li>TRUE：表示保持镜像的登录设置<br><li>FALSE：表示不保持镜像的登录设置<br><br>默认取值：FALSE。
	KeepImageLogin *bool `json:"KeepImageLogin" name:"KeepImageLogin"`
}

type ModifyAutoScalingGroupRequest struct {
	*tchttp.BaseRequest

	// 伸缩组ID
	AutoScalingGroupId *string `json:"AutoScalingGroupId" name:"AutoScalingGroupId"`

	// 伸缩组名称，在您账号中必须唯一。名称仅支持中文、英文、数字、下划线、分隔符"-"、小数点，最大长度不能超55个字节。
	AutoScalingGroupName *string `json:"AutoScalingGroupName" name:"AutoScalingGroupName"`

	// 默认冷却时间，单位秒，默认值为300
	DefaultCooldown *uint64 `json:"DefaultCooldown" name:"DefaultCooldown"`

	// 期望实例数，大小介于最小实例数和最大实例数之间
	DesiredCapacity *uint64 `json:"DesiredCapacity" name:"DesiredCapacity"`

	// 启动配置ID
	LaunchConfigurationId *string `json:"LaunchConfigurationId" name:"LaunchConfigurationId"`

	// 最大实例数，取值范围为0-2000。
	MaxSize *uint64 `json:"MaxSize" name:"MaxSize"`

	// 最小实例数，取值范围为0-2000。
	MinSize *uint64 `json:"MinSize" name:"MinSize"`

	// 项目ID
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`

	// 子网ID列表
	SubnetIds []*string `json:"SubnetIds" name:"SubnetIds" list`

	// 销毁策略，目前长度上限为1。取值包括 OLDEST_INSTANCE 和 NEWEST_INSTANCE。
	// <br><li> OLDEST_INSTANCE 优先销毁伸缩组中最旧的实例。
	// <br><li> NEWEST_INSTANCE，优先销毁伸缩组中最新的实例。
	TerminationPolicies []*string `json:"TerminationPolicies" name:"TerminationPolicies" list`

	// VPC ID，基础网络则填空字符串。修改为具体VPC ID时，需指定相应的SubnetIds；修改为基础网络时，需指定相应的Zones。
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 可用区列表
	Zones []*string `json:"Zones" name:"Zones" list`

	// 重试策略，取值包括 IMMEDIATE_RETRY 和 INCREMENTAL_INTERVALS，默认取值为 IMMEDIATE_RETRY。
	// <br><li> IMMEDIATE_RETRY，立即重试，在较短时间内快速重试，连续失败超过一定次数（5次）后不再重试。
	// <br><li> INCREMENTAL_INTERVALS，间隔递增重试，随着连续失败次数的增加，重试间隔逐渐增大，重试间隔从秒级到1天不等。在连续失败超过一定次数（25次）后不再重试。
	RetryPolicy *string `json:"RetryPolicy" name:"RetryPolicy"`

	// 可用区校验策略，取值包括 ALL 和 ANY，默认取值为ANY。在伸缩组实际变更资源相关字段时（启动配置、可用区、子网）发挥作用。
	// <br><li> ALL，所有可用区（Zone）或子网（SubnetId）都可用则通过校验，否则校验报错。
	// <br><li> ANY，存在任何一个可用区（Zone）或子网（SubnetId）可用则通过校验，否则校验报错。
	// 
	// 可用区或子网不可用的常见原因包括该可用区CVM实例类型售罄、该可用区CBS云盘售罄、该可用区配额不足、该子网IP不足等。
	// 如果 Zones/SubnetIds 中可用区或者子网不存在，则无论 ZonesCheckPolicy 采用何种取值，都会校验报错。
	ZonesCheckPolicy *string `json:"ZonesCheckPolicy" name:"ZonesCheckPolicy"`
}

func (r *ModifyAutoScalingGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAutoScalingGroupRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyAutoScalingGroupResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyAutoScalingGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyAutoScalingGroupResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDesiredCapacityRequest struct {
	*tchttp.BaseRequest

	// 伸缩组ID
	AutoScalingGroupId *string `json:"AutoScalingGroupId" name:"AutoScalingGroupId"`

	// 期望实例数
	DesiredCapacity *uint64 `json:"DesiredCapacity" name:"DesiredCapacity"`
}

func (r *ModifyDesiredCapacityRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDesiredCapacityRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDesiredCapacityResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDesiredCapacityResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDesiredCapacityResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyLaunchConfigurationAttributesRequest struct {
	*tchttp.BaseRequest

	// 启动配置ID
	LaunchConfigurationId *string `json:"LaunchConfigurationId" name:"LaunchConfigurationId"`

	// 指定有效的[镜像](https://cloud.tencent.com/document/product/213/4940)ID，格式形如`img-8toqc6s3`。镜像类型分为四种：<br/><li>公共镜像</li><li>自定义镜像</li><li>共享镜像</li><li>服务市场镜像</li><br/>可通过以下方式获取可用的镜像ID：<br/><li>`公共镜像`、`自定义镜像`、`共享镜像`的镜像ID可通过登录[控制台](https://console.cloud.tencent.com/cvm/image?rid=1&imageType=PUBLIC_IMAGE)查询；`服务镜像市场`的镜像ID可通过[云市场](https://market.cloud.tencent.com/list)查询。</li><li>通过调用接口 [DescribeImages](https://cloud.tencent.com/document/api/213/15715) ，取返回信息中的`ImageId`字段。</li>
	ImageId *string `json:"ImageId" name:"ImageId"`

	// 实例类型列表，不同实例机型指定了不同的资源规格，最多支持5中实例机型。
	// 启动配置，通过 InstanceType 表示单一实例类型，通过 InstanceTypes 表示多实例类型。指定 InstanceTypes 成功启动配置后，原有的 InstanceType 自动失效。
	InstanceTypes []*string `json:"InstanceTypes" name:"InstanceTypes" list`

	// 实例类型校验策略，在实际修改 InstanceTypes 时发挥作用，取值包括 ALL 和 ANY，默认取值为ANY。
	// <br><li> ALL，所有实例类型（InstanceType）都可用则通过校验，否则校验报错。
	// <br><li> ANY，存在任何一个实例类型（InstanceType）可用则通过校验，否则校验报错。
	// 
	// 实例类型不可用的常见原因包括该实例类型售罄、对应云盘售罄等。
	// 如果 InstanceTypes 中一款机型不存在或者已下线，则无论 InstanceTypesCheckPolicy 采用何种取值，都会校验报错。
	InstanceTypesCheckPolicy *string `json:"InstanceTypesCheckPolicy" name:"InstanceTypesCheckPolicy"`

	// 启动配置显示名称。名称仅支持中文、英文、数字、下划线、分隔符"-"、小数点，最大长度不能超60个字节。
	LaunchConfigurationName *string `json:"LaunchConfigurationName" name:"LaunchConfigurationName"`
}

func (r *ModifyLaunchConfigurationAttributesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyLaunchConfigurationAttributesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyLaunchConfigurationAttributesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyLaunchConfigurationAttributesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyLaunchConfigurationAttributesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyScheduledActionRequest struct {
	*tchttp.BaseRequest

	// 待修改的定时任务ID
	ScheduledActionId *string `json:"ScheduledActionId" name:"ScheduledActionId"`

	// 定时任务名称。名称仅支持中文、英文、数字、下划线、分隔符"-"、小数点，最大长度不能超60个字节。同一伸缩组下必须唯一。
	ScheduledActionName *string `json:"ScheduledActionName" name:"ScheduledActionName"`

	// 当定时任务触发时，设置的伸缩组最大实例数。
	MaxSize *uint64 `json:"MaxSize" name:"MaxSize"`

	// 当定时任务触发时，设置的伸缩组最小实例数。
	MinSize *uint64 `json:"MinSize" name:"MinSize"`

	// 当定时任务触发时，设置的伸缩组期望实例数。
	DesiredCapacity *uint64 `json:"DesiredCapacity" name:"DesiredCapacity"`

	// 定时任务的首次触发时间，取值为`北京时间`（UTC+8），按照`ISO8601`标准，格式：`YYYY-MM-DDThh:mm:ss+08:00`。
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 定时任务的结束时间，取值为`北京时间`（UTC+8），按照`ISO8601`标准，格式：`YYYY-MM-DDThh:mm:ss+08:00`。<br>此参数与`Recurrence`需要同时指定，到达结束时间之后，定时任务将不再生效。
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 定时任务的重复方式。为标准[Cron](https://zh.wikipedia.org/wiki/Cron)格式<br>此参数与`EndTime`需要同时指定。
	Recurrence *string `json:"Recurrence" name:"Recurrence"`
}

func (r *ModifyScheduledActionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyScheduledActionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyScheduledActionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyScheduledActionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyScheduledActionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RemoveInstancesRequest struct {
	*tchttp.BaseRequest

	// 伸缩组ID
	AutoScalingGroupId *string `json:"AutoScalingGroupId" name:"AutoScalingGroupId"`

	// CVM实例ID列表
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`
}

func (r *RemoveInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RemoveInstancesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RemoveInstancesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RemoveInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RemoveInstancesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RunMonitorServiceEnabled struct {

	// 是否开启[云监控](https://cloud.tencent.com/document/product/248)服务。取值范围：<br><li>TRUE：表示开启云监控服务<br><li>FALSE：表示不开启云监控服务<br><br>默认取值：TRUE。
	Enabled *bool `json:"Enabled" name:"Enabled"`
}

type RunSecurityServiceEnabled struct {

	// 是否开启[云安全](https://cloud.tencent.com/document/product/296)服务。取值范围：<br><li>TRUE：表示开启云安全服务<br><li>FALSE：表示不开启云安全服务<br><br>默认取值：TRUE。
	Enabled *bool `json:"Enabled" name:"Enabled"`
}

type ScheduledAction struct {

	// 定时任务ID。
	ScheduledActionId *string `json:"ScheduledActionId" name:"ScheduledActionId"`

	// 定时任务名称。
	ScheduledActionName *string `json:"ScheduledActionName" name:"ScheduledActionName"`

	// 定时任务所在伸缩组ID。
	AutoScalingGroupId *string `json:"AutoScalingGroupId" name:"AutoScalingGroupId"`

	// 定时任务的开始时间。取值为`北京时间`（UTC+8），按照`ISO8601`标准，格式：`YYYY-MM-DDThh:mm:ss+08:00`。
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 定时任务的重复方式。
	Recurrence *string `json:"Recurrence" name:"Recurrence"`

	// 定时任务的结束时间。取值为`北京时间`（UTC+8），按照`ISO8601`标准，格式：`YYYY-MM-DDThh:mm:ss+08:00`。
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 定时任务设置的最大实例数。
	MaxSize *uint64 `json:"MaxSize" name:"MaxSize"`

	// 定时任务设置的期望实例数。
	DesiredCapacity *uint64 `json:"DesiredCapacity" name:"DesiredCapacity"`

	// 定时任务设置的最小实例数。
	MinSize *uint64 `json:"MinSize" name:"MinSize"`

	// 定时任务的创建时间。取值为`UTC`时间，按照`ISO8601`标准，格式：`YYYY-MM-DDThh:mm:ssZ`。
	CreatedTime *string `json:"CreatedTime" name:"CreatedTime"`
}

type SpotMarketOptions struct {

	// 竞价出价，例如“1.05”
	MaxPrice *string `json:"MaxPrice" name:"MaxPrice"`

	// 竞价请求类型，当前仅支持类型：one-time，默认值为one-time
	SpotInstanceType *string `json:"SpotInstanceType" name:"SpotInstanceType"`
}

type SystemDisk struct {

	// 系统盘类型。系统盘类型限制详见[CVM实例配置](https://cloud.tencent.com/document/product/213/2177)。取值范围：<br><li>LOCAL_BASIC：本地硬盘<br><li>LOCAL_SSD：本地SSD硬盘<br><li>CLOUD_BASIC：普通云硬盘<br><li>CLOUD_PREMIUM：高性能云硬盘<br><li>CLOUD_SSD：SSD云硬盘<br><br>默认取值：LOCAL_BASIC。
	DiskType *string `json:"DiskType" name:"DiskType"`

	// 系统盘大小，单位：GB。默认值为 50
	DiskSize *uint64 `json:"DiskSize" name:"DiskSize"`
}

type TargetAttribute struct {

	// 端口
	Port *uint64 `json:"Port" name:"Port"`

	// 权重
	Weight *uint64 `json:"Weight" name:"Weight"`
}
