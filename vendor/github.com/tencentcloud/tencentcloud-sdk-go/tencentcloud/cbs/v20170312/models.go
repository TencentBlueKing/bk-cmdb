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

package v20170312

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type ApplySnapshotRequest struct {
	*tchttp.BaseRequest

	// 快照ID, 可通过[DescribeSnapshots](/document/product/362/15647)查询。
	SnapshotId *string `json:"SnapshotId" name:"SnapshotId"`

	// 快照原云硬盘ID，可通过[DescribeDisks](/document/product/362/16315)接口查询。
	DiskId *string `json:"DiskId" name:"DiskId"`
}

func (r *ApplySnapshotRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ApplySnapshotRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ApplySnapshotResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ApplySnapshotResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ApplySnapshotResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AttachDetail struct {

	// 实例ID。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 实例已挂载数据盘的数量。
	AttachedDiskCount *uint64 `json:"AttachedDiskCount" name:"AttachedDiskCount"`

	// 实例最大可挂载数据盘的数量。
	MaxAttachCount *uint64 `json:"MaxAttachCount" name:"MaxAttachCount"`
}

type AttachDisksRequest struct {
	*tchttp.BaseRequest

	// 将要被挂载的弹性云盘ID。通过[DescribeDisks](/document/product/362/16315)接口查询。单次最多可挂载10块弹性云盘。
	DiskIds []*string `json:"DiskIds" name:"DiskIds" list`

	// 云服务器实例ID。云盘将被挂载到此云服务器上，通过[DescribeInstances](/document/product/213/15728)接口查询。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 可选参数，不传该参数则仅执行挂载操作。传入`True`时，会在挂载成功后将云硬盘设置为随云主机销毁模式，仅对按量计费云硬盘有效。
	DeleteWithInstance *bool `json:"DeleteWithInstance" name:"DeleteWithInstance"`
}

func (r *AttachDisksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AttachDisksRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AttachDisksResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AttachDisksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AttachDisksResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDisksRequest struct {
	*tchttp.BaseRequest

	// 硬盘介质类型。取值范围：<br><li>CLOUD_BASIC：表示普通云硬盘<br><li>CLOUD_PREMIUM：表示高性能云硬盘<br><li>CLOUD_SSD：表示SSD云硬盘。
	DiskType *string `json:"DiskType" name:"DiskType"`

	// 云硬盘计费类型。<br><li>PREPAID：预付费，即包年包月<br><li>POSTPAID_BY_HOUR：按小时后付费<br>各类型价格请参考云硬盘[价格总览](/document/product/362/2413)。
	DiskChargeType *string `json:"DiskChargeType" name:"DiskChargeType"`

	// 实例所在的位置。通过该参数可以指定实例所属可用区，所属项目。若不指定项目，将在默认项目下进行创建。
	Placement *Placement `json:"Placement" name:"Placement"`

	// 云盘显示名称。不传则默认为“未命名”。最大长度不能超60个字节。
	DiskName *string `json:"DiskName" name:"DiskName"`

	// 创建云硬盘数量，不传则默认为1。单次请求最多可创建的云盘数有限制，具体参见[云硬盘使用限制](https://cloud.tencent.com/doc/product/362/5145)。
	DiskCount *uint64 `json:"DiskCount" name:"DiskCount"`

	// 预付费模式，即包年包月相关参数设置。通过该参数指定包年包月云盘的购买时长、是否设置自动续费等属性。<br>创建预付费云盘该参数必传，创建按小时后付费云盘无需传该参数。
	DiskChargePrepaid *DiskChargePrepaid `json:"DiskChargePrepaid" name:"DiskChargePrepaid"`

	// 云硬盘大小，单位为GB。<br><li>如果传入`SnapshotId`则可不传`DiskSize`，此时新建云盘的大小为快照大小<br><li>如果传入`SnapshotId`同时传入`DiskSize`，则云盘大小必须大于或等于快照大小<br><li>云盘大小取值范围参见云硬盘[产品分类](/document/product/362/2353)的说明。
	DiskSize *uint64 `json:"DiskSize" name:"DiskSize"`

	// 快照ID，如果传入则根据此快照创建云硬盘，快照类型必须为数据盘快照，可通过[DescribeSnapshots](/document/product/362/15647)接口查询快照，见输出参数DiskUsage解释。
	SnapshotId *string `json:"SnapshotId" name:"SnapshotId"`

	// 用于保证请求幂等性的字符串。该字符串由客户生成，需保证不同请求之间唯一，最大值不超过64个ASCII字符。若不指定该参数，则无法保证请求的幂等性。
	ClientToken *string `json:"ClientToken" name:"ClientToken"`

	// 传入该参数用于创建加密云盘，取值固定为ENCRYPT。
	Encrypt *string `json:"Encrypt" name:"Encrypt"`

	// 云盘绑定的标签。
	Tags []*Tag `json:"Tags" name:"Tags" list`

	// 可选参数，不传该参数则仅执行挂载操作。传入True时，新创建的云盘将设置为随云主机销毁模式，仅对按量计费云硬盘有效。
	DeleteWithInstance *bool `json:"DeleteWithInstance" name:"DeleteWithInstance"`
}

func (r *CreateDisksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDisksRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateDisksResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 创建的云硬盘ID列表。
		DiskIdSet []*string `json:"DiskIdSet" name:"DiskIdSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateDisksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateDisksResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateSnapshotRequest struct {
	*tchttp.BaseRequest

	// 需要创建快照的云硬盘ID，可通过[DescribeDisks](/document/product/362/16315)接口查询。
	DiskId *string `json:"DiskId" name:"DiskId"`

	// 快照名称，不传则新快照名称默认为“未命名”。
	SnapshotName *string `json:"SnapshotName" name:"SnapshotName"`
}

func (r *CreateSnapshotRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateSnapshotRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateSnapshotResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 新创建的快照ID。
		SnapshotId *string `json:"SnapshotId" name:"SnapshotId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateSnapshotResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateSnapshotResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteSnapshotsRequest struct {
	*tchttp.BaseRequest

	// 要删除的快照ID列表，可通过[DescribeSnapshots](/document/product/362/15647)查询。
	SnapshotIds []*string `json:"SnapshotIds" name:"SnapshotIds" list`
}

func (r *DeleteSnapshotsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteSnapshotsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteSnapshotsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteSnapshotsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteSnapshotsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDiskConfigQuotaRequest struct {
	*tchttp.BaseRequest

	// 查询类别，取值范围。<br><li>INQUIRY_CBS_CONFIG：查询云盘配置列表<br><li>INQUIRY_CVM_CONFIG：查询云盘与实例搭配的配置列表。
	InquiryType *string `json:"InquiryType" name:"InquiryType"`

	// 查询一个或多个[可用区](/document/api/213/9452#zone)下的配置。
	Zones []*string `json:"Zones" name:"Zones" list`

	// 付费模式。取值范围：<br><li>PREPAID：预付费<br><li>POSTPAID_BY_HOUR：后付费。
	DiskChargeType *string `json:"DiskChargeType" name:"DiskChargeType"`

	// 硬盘介质类型。取值范围：<br><li>CLOUD_BASIC：表示普通云硬盘<br><li>CLOUD_PREMIUM：表示高性能云硬盘<br><li>CLOUD_SSD：表示SSD云硬盘。
	DiskTypes []*string `json:"DiskTypes" name:"DiskTypes" list`

	// 系统盘或数据盘。取值范围：<br><li>SYSTEM_DISK：表示系统盘<br><li>DATA_DISK：表示数据盘。
	DiskUsage *string `json:"DiskUsage" name:"DiskUsage"`

	// 按照实例机型系列过滤。实例机型系列形如：S1、I1、M1等。详见[实例类型](https://cloud.tencent.com/document/product/213/11518)
	InstanceFamilies []*string `json:"InstanceFamilies" name:"InstanceFamilies" list`

	// 实例CPU核数。
	CPU *uint64 `json:"CPU" name:"CPU"`

	// 实例内存大小。
	Memory *uint64 `json:"Memory" name:"Memory"`
}

func (r *DescribeDiskConfigQuotaRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDiskConfigQuotaRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDiskConfigQuotaResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 云盘配置列表。
		DiskConfigSet []*DiskConfig `json:"DiskConfigSet" name:"DiskConfigSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDiskConfigQuotaResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDiskConfigQuotaResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDiskOperationLogsRequest struct {
	*tchttp.BaseRequest

	// 过滤条件。支持以下条件：
	// <li>disk-id - Array of String - 是否必填：是 - 按云盘ID过滤，每个请求最多可指定10个云盘ID。
	Filters []*Filter `json:"Filters" name:"Filters" list`
}

func (r *DescribeDiskOperationLogsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDiskOperationLogsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDiskOperationLogsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 云盘的操作日志列表。
		DiskOperationLogSet []*DiskOperationLog `json:"DiskOperationLogSet" name:"DiskOperationLogSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDiskOperationLogsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDiskOperationLogsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDisksRequest struct {
	*tchttp.BaseRequest

	// 按照一个或者多个云硬盘ID查询。云硬盘ID形如：`disk-11112222`，此参数的具体格式可参考API[简介](/document/product/362/15633)的ids.N一节）。参数不支持同时指定`DiskIds`和`Filters`。
	DiskIds []*string `json:"DiskIds" name:"DiskIds" list`

	// 过滤条件。参数不支持同时指定`DiskIds`和`Filters`。<br><li>disk-usage - Array of String - 是否必填：否 -（过滤条件）按云盘类型过滤。 (SYSTEM_DISK：表示系统盘 | DATA_DISK：表示数据盘)<br><li>disk-charge-type - Array of String - 是否必填：否 -（过滤条件）按照云硬盘计费模式过滤。 (PREPAID：表示预付费，即包年包月 | POSTPAID_BY_HOUR：表示后付费，即按量计费。)<br><li>portable - Array of String - 是否必填：否 -（过滤条件）按是否为弹性云盘过滤。 (TRUE：表示弹性云盘 | FALSE：表示非弹性云盘。)<br><li>project-id - Array of Integer - 是否必填：否 -（过滤条件）按云硬盘所属项目ID过滤。<br><li>disk-id - Array of String - 是否必填：否 -（过滤条件）按照云硬盘ID过滤。云盘ID形如：`disk-11112222`。<br><li>disk-name - Array of String - 是否必填：否 -（过滤条件）按照云盘名称过滤。<br><li>disk-type - Array of String - 是否必填：否 -（过滤条件）按照云盘介质类型过滤。(CLOUD_BASIC：表示普通云硬盘 | CLOUD_PREMIUM：表示高性能云硬盘。| CLOUD_SSD：SSD表示SSD云硬盘。)<br><li>disk-state - Array of String - 是否必填：否 -（过滤条件）按照云盘状态过滤。(UNATTACHED：未挂载 | ATTACHING：挂载中 | ATTACHED：已挂载 | DETACHING：解挂中 | EXPANDING：扩容中 | ROLLBACKING：回滚中 | TORECYCLE：待回收。)<br><li>instance-id - Array of String - 是否必填：否 -（过滤条件）按照云盘挂载的云主机实例ID过滤。可根据此参数查询挂载在指定云主机下的云硬盘。<br><li>zone - Array of String - 是否必填：否 -（过滤条件）按照[可用区](/document/api/213/9452#zone)过滤。<br><li>instance-ip-address - Array of String - 是否必填：否 -（过滤条件）按云盘所挂载云主机的内网或外网IP过滤。<br><li>instance-name - Array of String - 是否必填：否 -（过滤条件）按云盘所挂载的实例名称过滤。<br><li>tag-key - Array of String - 是否必填：否 -（过滤条件）按照标签键进行过滤。<br><li>tag-value - Array of String - 是否必填：否 -（过滤条件）照标签值进行过滤。<br><li>tag:tag-key - Array of String - 是否必填：否 -（过滤条件）按照标签键值对进行过滤。 tag-key使用具体的标签键进行替换。
	Filters []*Filter `json:"Filters" name:"Filters" list`

	// 偏移量，默认为0。关于`Offset`的更进一步介绍请参考API[简介](/document/product/362/15633)中的相关小节。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 返回数量，默认为20，最大值为100。关于`Limit`的更进一步介绍请参考 API [简介](/document/product/362/15633)中的相关小节。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 输出云盘列表的排列顺序。取值范围：<br><li>ASC：升序排列<br><li>DESC：降序排列。
	Order *string `json:"Order" name:"Order"`

	// 云盘列表排序的依据字段。取值范围：<br><li>CREATE_TIME：依据云盘的创建时间排序<br><li>DEADLINE：依据云盘的到期时间排序<br>默认按云盘创建时间排序。
	OrderField *string `json:"OrderField" name:"OrderField"`

	// 云盘详情中是否需要返回云盘绑定的定期快照策略ID，TRUE表示需要返回，FALSE表示不返回。
	ReturnBindAutoSnapshotPolicy *bool `json:"ReturnBindAutoSnapshotPolicy" name:"ReturnBindAutoSnapshotPolicy"`
}

func (r *DescribeDisksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDisksRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDisksResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合条件的云硬盘数量。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 云硬盘的详细信息列表。
		DiskSet []*Disk `json:"DiskSet" name:"DiskSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDisksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDisksResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeInstancesDiskNumRequest struct {
	*tchttp.BaseRequest

	// 云服务器实例ID，通过[DescribeInstances](/document/product/213/15728)接口查询。
	InstanceIds []*string `json:"InstanceIds" name:"InstanceIds" list`
}

func (r *DescribeInstancesDiskNumRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeInstancesDiskNumRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeInstancesDiskNumResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 各个云服务器已挂载和可挂载弹性云盘的数量。
		AttachDetail []*AttachDetail `json:"AttachDetail" name:"AttachDetail" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeInstancesDiskNumResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeInstancesDiskNumResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSnapshotsRequest struct {
	*tchttp.BaseRequest

	// 要查询快照的ID列表。参数不支持同时指定`SnapshotIds`和`Filters`。
	SnapshotIds []*string `json:"SnapshotIds" name:"SnapshotIds" list`

	// 过滤条件。参数不支持同时指定`SnapshotIds`和`Filters`。<br><li>snapshot-id - Array of String - 是否必填：否 -（过滤条件）按照快照的ID过滤。快照ID形如：`snap-11112222`。<br><li>snapshot-name - Array of String - 是否必填：否 -（过滤条件）按照快照名称过滤。<br><li>snapshot-state - Array of String - 是否必填：否 -（过滤条件）按照快照状态过滤。 (NORMAL：正常 | CREATING：创建中 | ROLLBACKING：回滚中。)<br><li>disk-usage - Array of String - 是否必填：否 -（过滤条件）按创建快照的云盘类型过滤。 (SYSTEM_DISK：代表系统盘 | DATA_DISK：代表数据盘。)<br><li>project-id  - Array of String - 是否必填：否 -（过滤条件）按云硬盘所属项目ID过滤。<br><li>disk-id  - Array of String - 是否必填：否 -（过滤条件）按照创建快照的云硬盘ID过滤。<br><li>zone - Array of String - 是否必填：否 -（过滤条件）按照[可用区](/document/api/213/9452#zone)过滤。<br><li>encrypt - Array of String - 是否必填：否 -（过滤条件）按是否加密盘快照过滤。 (TRUE：表示加密盘快照 | FALSE：表示非加密盘快照。)
	Filters []*Filter `json:"Filters" name:"Filters" list`

	// 偏移量，默认为0。关于`Offset`的更进一步介绍请参考API[简介](/document/product/362/15633)中的相关小节。
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 返回数量，默认为20，最大值为100。关于`Limit`的更进一步介绍请参考 API [简介](/document/product/362/15633)中的相关小节。
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 输出云盘列表的排列顺序。取值范围：<br><li>ASC：升序排列<br><li>DESC：降序排列。
	Order *string `json:"Order" name:"Order"`

	// 快照列表排序的依据字段。取值范围：<br><li>CREATE_TIME：依据快照的创建时间排序<br>默认按创建时间排序。
	OrderField *string `json:"OrderField" name:"OrderField"`
}

func (r *DescribeSnapshotsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSnapshotsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeSnapshotsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 快照的数量。
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 快照的详情列表。
		SnapshotSet []*Snapshot `json:"SnapshotSet" name:"SnapshotSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeSnapshotsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeSnapshotsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DetachDisksRequest struct {
	*tchttp.BaseRequest

	// 将要解挂的云硬盘ID， 通过[DescribeDisks](/document/product/362/16315)接口查询，单次请求最多可解挂10块弹性云盘。
	DiskIds []*string `json:"DiskIds" name:"DiskIds" list`
}

func (r *DetachDisksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DetachDisksRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DetachDisksResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DetachDisksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DetachDisksResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Disk struct {

	// 云硬盘ID。
	DiskId *string `json:"DiskId" name:"DiskId"`

	// 云硬盘类型。取值范围：<br><li>SYSTEM_DISK：系统盘<br><li>DATA_DISK：数据盘。
	DiskUsage *string `json:"DiskUsage" name:"DiskUsage"`

	// 付费模式。取值范围：<br><li>PREPAID：预付费，即包年包月<br><li>POSTPAID_BY_HOUR：后付费，即按量计费。
	DiskChargeType *string `json:"DiskChargeType" name:"DiskChargeType"`

	// 是否为弹性云盘，false表示非弹性云盘，true表示弹性云盘。
	Portable *bool `json:"Portable" name:"Portable"`

	// 云硬盘所在的位置。
	Placement *Placement `json:"Placement" name:"Placement"`

	// 云盘是否具备创建快照的能力。取值范围：<br><li>false表示不具备<br><li>true表示具备。
	SnapshotAbility *bool `json:"SnapshotAbility" name:"SnapshotAbility"`

	// 云硬盘名称。
	DiskName *string `json:"DiskName" name:"DiskName"`

	// 云硬盘大小，单位GB。
	DiskSize *uint64 `json:"DiskSize" name:"DiskSize"`

	// 云盘状态。取值范围：<br><li>UNATTACHED：未挂载<br><li>ATTACHING：挂载中<br><li>ATTACHED：已挂载<br><li>DETACHING：解挂中<br><li>EXPANDING：扩容中<br><li>ROLLBACKING：回滚中。
	DiskState *string `json:"DiskState" name:"DiskState"`

	// 云盘介质类型。取值范围：<br><li>CLOUD_BASIC：表示普通云硬<br><li>CLOUD_PREMIUM：表示高性能云硬盘<br><li>CLOUD_SSD：SSD表示SSD云硬盘。
	DiskType *string `json:"DiskType" name:"DiskType"`

	// 云盘是否挂载到云主机上。取值范围：<br><li>false:表示未挂载<br><li>true:表示已挂载。
	Attached *bool `json:"Attached" name:"Attached"`

	// 云硬盘挂载的云主机ID。
	InstanceId *string `json:"InstanceId" name:"InstanceId"`

	// 云硬盘的创建时间。
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 云硬盘的到期时间。
	DeadlineTime *string `json:"DeadlineTime" name:"DeadlineTime"`

	// 云盘是否处于快照回滚状态。取值范围：<br><li>false:表示不处于快照回滚状态<br><li>true:表示处于快照回滚状态。
	Rollbacking *bool `json:"Rollbacking" name:"Rollbacking"`

	// 云盘快照回滚的进度。
	RollbackPercent *uint64 `json:"RollbackPercent" name:"RollbackPercent"`

	// 云盘是否为加密盘。取值范围：<br><li>false:表示非加密盘<br><li>true:表示加密盘。
	Encrypt *bool `json:"Encrypt" name:"Encrypt"`

	// 云盘已挂载到子机，且子机与云盘都是包年包月。<br><li>true：子机设置了自动续费标识，但云盘未设置<br><li>false：云盘自动续费标识正常。
	AutoRenewFlagError *bool `json:"AutoRenewFlagError" name:"AutoRenewFlagError"`

	// 自动续费标识。取值范围：<br><li>NOTIFY_AND_AUTO_RENEW：通知过期且自动续费<br><li>NOTIFY_AND_MANUAL_RENEW：通知过期不自动续费<br><li>DISABLE_NOTIFY_AND_MANUAL_RENEW：不通知过期不自动续费。
	RenewFlag *string `json:"RenewFlag" name:"RenewFlag"`

	// 在云盘已挂载到实例，且实例与云盘都是包年包月的条件下，此字段才有意义。<br><li>true:云盘到期时间早于实例。<br><li>false：云盘到期时间晚于实例。
	DeadlineError *bool `json:"DeadlineError" name:"DeadlineError"`

	// 判断预付费的云盘是否支持主动退还。<br><li>true:支持主动退还<br><li>false:不支持主动退还。
	IsReturnable *bool `json:"IsReturnable" name:"IsReturnable"`

	// 预付费云盘在不支持主动退还的情况下，该参数表明不支持主动退还的具体原因。取值范围：<br><li>1：云硬盘已经退还<br><li>2：云硬盘已过期<br><li>3：云盘不支持退还<br><li>8：超过可退还数量的限制。
	ReturnFailCode *int64 `json:"ReturnFailCode" name:"ReturnFailCode"`

	// 云盘关联的定期快照ID。只有在调用DescribeDisks接口时，入参ReturnBindAutoSnapshotPolicy取值为TRUE才会返回该参数。
	AutoSnapshotPolicyIds []*string `json:"AutoSnapshotPolicyIds" name:"AutoSnapshotPolicyIds" list`

	// 与云盘绑定的标签，云盘未绑定标签则取值为空。
	Tags []*Tag `json:"Tags" name:"Tags" list`

	// 云盘是否与挂载的实例一起销毁。<br><li>true:销毁实例时会同时销毁云盘，只支持按小时后付费云盘。<br><li>false：销毁实例时不销毁云盘。
	DeleteWithInstance *bool `json:"DeleteWithInstance" name:"DeleteWithInstance"`

	// 当前时间距离盘到期的天数（仅对预付费盘有意义）。
	DifferDaysOfDeadline *int64 `json:"DifferDaysOfDeadline" name:"DifferDaysOfDeadline"`
}

type DiskChargePrepaid struct {

	// 购买云盘的时长，默认单位为月，取值范围：1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36。
	Period *uint64 `json:"Period" name:"Period"`

	// 自动续费标识。取值范围：<br><li>NOTIFY_AND_AUTO_RENEW：通知过期且自动续费<br><li>NOTIFY_AND_MANUAL_RENEW：通知过期不自动续费<br><li>DISABLE_NOTIFY_AND_MANUAL_RENEW：不通知过期不自动续费<br><br>默认取值：NOTIFY_AND_MANUAL_RENEW：通知过期不自动续费。
	RenewFlag *string `json:"RenewFlag" name:"RenewFlag"`

	// 需要将云盘的到期时间与挂载的子机对齐时，可传入该参数。该参数表示子机当前的到期时间，此时Period如果传入，则表示子机需要续费的时长，云盘会自动按对齐到子机续费后的到期时间续费。
	CurInstanceDeadline *string `json:"CurInstanceDeadline" name:"CurInstanceDeadline"`
}

type DiskConfig struct {

	// 配置是否可用。
	Available *bool `json:"Available" name:"Available"`

	// 云盘介质类型。取值范围：<br><li>CLOUD_BASIC：表示普通云硬盘<br><li>CLOUD_PREMIUM：表示高性能云硬盘<br><li>CLOUD_SSD：SSD表示SSD云硬盘。
	DiskType *string `json:"DiskType" name:"DiskType"`

	// 云盘类型。取值范围：<br><li>SYSTEM_DISK：表示系统盘<br><li>DATA_DISK：表示数据盘。
	DiskUsage *string `json:"DiskUsage" name:"DiskUsage"`

	// 付费模式。取值范围：<br><li>PREPAID：表示预付费，即包年包月<br><li>POSTPAID_BY_HOUR：表示后付费，即按量计费。
	DiskChargeType *string `json:"DiskChargeType" name:"DiskChargeType"`

	// 最大可配置云盘大小，单位GB。
	MaxDiskSize *uint64 `json:"MaxDiskSize" name:"MaxDiskSize"`

	// 最小可配置云盘大小，单位GB。
	MinDiskSize *uint64 `json:"MinDiskSize" name:"MinDiskSize"`

	// 所在[可用区](/document/api/213/9452#zone)。
	Zone *string `json:"Zone" name:"Zone"`

	// 实例机型。
	DeviceClass *string `json:"DeviceClass" name:"DeviceClass"`

	// 实例机型系列。详见[实例类型](https://cloud.tencent.com/document/product/213/11518)
	InstanceFamily *string `json:"InstanceFamily" name:"InstanceFamily"`
}

type DiskOperationLog struct {

	// 操作者的UIN。
	Operator *string `json:"Operator" name:"Operator"`

	// 操作类型。取值范围：
	// CBS_OPERATION_ATTACH：挂载云硬盘
	// CBS_OPERATION_DETACH：解挂云硬盘
	// CBS_OPERATION_RENEW：续费
	// CBS_OPERATION_EXPAND：扩容
	// CBS_OPERATION_CREATE：创建
	// CBS_OPERATION_ISOLATE：隔离
	// CBS_OPERATION_MODIFY：修改云硬盘属性
	// ASP_OPERATION_BIND：关联定期快照策略
	// ASP_OPERATION_UNBIND：取消关联定期快照策略
	Operation *string `json:"Operation" name:"Operation"`

	// 操作的云盘ID。
	DiskId *string `json:"DiskId" name:"DiskId"`

	// 操作的状态。取值范围：
	// SUCCESS :表示操作成功 
	// FAILED :表示操作失败 
	// PROCESSING :表示操作中。
	OperationState *string `json:"OperationState" name:"OperationState"`

	// 开始时间。
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 结束时间。
	EndTime *string `json:"EndTime" name:"EndTime"`
}

type Filter struct {

	// 过滤键的名称。
	Name *string `json:"Name" name:"Name"`

	// 一个或者多个过滤值。
	Values []*string `json:"Values" name:"Values" list`
}

type InquiryPriceCreateDisksRequest struct {
	*tchttp.BaseRequest

	// 云硬盘类型。取值范围：<br><li>普通云硬盘：CLOUD_BASIC<br><li>高性能云硬盘：CLOUD_PREMIUM<br><li>SSD云硬盘：CLOUD_SSD。
	DiskType *string `json:"DiskType" name:"DiskType"`

	// 云硬盘大小，单位为GB。云盘大小取值范围参见云硬盘[产品分类](/document/product/362/2353)的说明。
	DiskSize *uint64 `json:"DiskSize" name:"DiskSize"`

	// 云硬盘计费类型。<br><li>PREPAID：预付费，即包年包月<br><li>POSTPAID_BY_HOUR：按小时后付费
	DiskChargeType *string `json:"DiskChargeType" name:"DiskChargeType"`

	// 预付费模式，即包年包月相关参数设置。通过该参数指定包年包月云盘的购买时长、是否设置自动续费等属性。<br>创建预付费云盘该参数必传，创建按小时后付费云盘无需传该参数。
	DiskChargePrepaid *DiskChargePrepaid `json:"DiskChargePrepaid" name:"DiskChargePrepaid"`

	// 购买云盘的数量。不填则默认为1。
	DiskCount *uint64 `json:"DiskCount" name:"DiskCount"`

	// 云盘所属项目ID。
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`
}

func (r *InquiryPriceCreateDisksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceCreateDisksRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InquiryPriceCreateDisksResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 描述了新购云盘的价格。
		DiskPrice *Price `json:"DiskPrice" name:"DiskPrice"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *InquiryPriceCreateDisksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceCreateDisksResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InquiryPriceRenewDisksRequest struct {
	*tchttp.BaseRequest

	// 云硬盘ID， 通过[DescribeDisks](/document/product/362/16315)接口查询。
	DiskIds []*string `json:"DiskIds" name:"DiskIds" list`

	// 预付费模式，即包年包月相关参数设置。通过该参数可以指定包年包月云盘的购买时长。如果在该参数中指定CurInstanceDeadline，则会按对齐到子机到期时间来续费。如果是批量续费询价，该参数与Disks参数一一对应，元素数量需保持一致。
	DiskChargePrepaids []*DiskChargePrepaid `json:"DiskChargePrepaids" name:"DiskChargePrepaids" list`

	// 指定云盘新的到期时间，形式如：2017-12-17 00:00:00。参数`NewDeadline`和`DiskChargePrepaids`是两种指定询价时长的方式，两者必传一个。
	NewDeadline *string `json:"NewDeadline" name:"NewDeadline"`

	// 云盘所属项目ID。 如传入则仅用于鉴权。
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`
}

func (r *InquiryPriceRenewDisksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceRenewDisksRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InquiryPriceRenewDisksResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 描述了续费云盘的价格。
		DiskPrice *PrepayPrice `json:"DiskPrice" name:"DiskPrice"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *InquiryPriceRenewDisksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceRenewDisksResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InquiryPriceResizeDiskRequest struct {
	*tchttp.BaseRequest

	// 云硬盘ID， 通过[DescribeDisks](/document/product/362/16315)接口查询。
	DiskId *string `json:"DiskId" name:"DiskId"`

	// 云硬盘扩容后的大小，单位为GB，不得小于当前云硬盘大小。云盘大小取值范围参见云硬盘[产品分类](/document/product/362/2353)的说明。
	DiskSize *uint64 `json:"DiskSize" name:"DiskSize"`

	// 云盘所属项目ID。 如传入则仅用于鉴权。
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`
}

func (r *InquiryPriceResizeDiskRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceResizeDiskRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InquiryPriceResizeDiskResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 描述了扩容云盘的价格。
		DiskPrice *PrepayPrice `json:"DiskPrice" name:"DiskPrice"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *InquiryPriceResizeDiskResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InquiryPriceResizeDiskResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDiskAttributesRequest struct {
	*tchttp.BaseRequest

	// 一个或多个待操作的云硬盘ID。如果传入多个云盘ID，仅支持所有云盘修改为同一属性。
	DiskIds []*string `json:"DiskIds" name:"DiskIds" list`

	// 新的云硬盘项目ID，只支持修改弹性云盘的项目ID。通过[DescribeProject](/document/api/378/4400)接口查询可用项目及其ID。
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`

	// 新的云硬盘名称。
	DiskName *string `json:"DiskName" name:"DiskName"`

	// 是否为弹性云盘，FALSE表示非弹性云盘，TRUE表示弹性云盘。仅支持非弹性云盘修改为弹性云盘。
	Portable *bool `json:"Portable" name:"Portable"`

	// 成功挂载到云主机后该云硬盘是否随云主机销毁，TRUE表示随云主机销毁，FALSE表示不随云主机销毁。仅支持按量计费云硬盘数据盘。
	DeleteWithInstance *bool `json:"DeleteWithInstance" name:"DeleteWithInstance"`
}

func (r *ModifyDiskAttributesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDiskAttributesRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDiskAttributesResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDiskAttributesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDiskAttributesResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDisksRenewFlagRequest struct {
	*tchttp.BaseRequest

	// 一个或多个待操作的云硬盘ID。
	DiskIds []*string `json:"DiskIds" name:"DiskIds" list`

	// 云盘的续费标识。取值范围：<br><li>NOTIFY_AND_AUTO_RENEW：通知过期且自动续费<br><li>NOTIFY_AND_MANUAL_RENEW：通知过期不自动续费<br><li>DISABLE_NOTIFY_AND_MANUAL_RENEW：不通知过期不自动续费。
	RenewFlag *string `json:"RenewFlag" name:"RenewFlag"`
}

func (r *ModifyDisksRenewFlagRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDisksRenewFlagRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyDisksRenewFlagResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyDisksRenewFlagResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyDisksRenewFlagResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifySnapshotAttributeRequest struct {
	*tchttp.BaseRequest

	// 快照ID, 可通过[DescribeSnapshots](/document/product/362/15647)查询。
	SnapshotId *string `json:"SnapshotId" name:"SnapshotId"`

	// 新的快照名称。最长为60个字符。
	SnapshotName *string `json:"SnapshotName" name:"SnapshotName"`

	// 快照的保留时间，FALSE表示非永久保留，TRUE表示永久保留。仅支持将非永久快照修改为永久快照。
	IsPermanent *bool `json:"IsPermanent" name:"IsPermanent"`
}

func (r *ModifySnapshotAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifySnapshotAttributeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifySnapshotAttributeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifySnapshotAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifySnapshotAttributeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Placement struct {

	// 实例所属的[可用区](/document/api/213/9452#zone)ID。该参数也可以通过调用  [DescribeZones](/document/product/213/15707) 的返回值中的Zone字段来获取。
	Zone *string `json:"Zone" name:"Zone"`

	// 实例所属项目ID。该参数可以通过调用 [DescribeProject](/document/api/378/4400) 的返回值中的 projectId 字段来获取。不填为默认项目。
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`
}

type PrepayPrice struct {

	// 预付费云盘或快照预支费用的原价，单位：元。
	OriginalPrice *float64 `json:"OriginalPrice" name:"OriginalPrice"`

	// 预付费云盘或快照预支费用的折扣价，单位：元。
	DiscountPrice *float64 `json:"DiscountPrice" name:"DiscountPrice"`
}

type Price struct {

	// 预付费云盘预支费用的原价，单位：元。
	OriginalPrice *float64 `json:"OriginalPrice" name:"OriginalPrice"`

	// 预付费云盘预支费用的折扣价，单位：元。
	DiscountPrice *float64 `json:"DiscountPrice" name:"DiscountPrice"`

	// 后付费云盘的单价，单位：元。
	UnitPrice *float64 `json:"UnitPrice" name:"UnitPrice"`

	// 后付费云盘的计价单元，取值范围：<br><li>HOUR：表示后付费云盘的计价单元是按小时计算。
	ChargeUnit *string `json:"ChargeUnit" name:"ChargeUnit"`
}

type RenewDiskRequest struct {
	*tchttp.BaseRequest

	// 预付费模式，即包年包月相关参数设置。通过该参数可以指定包年包月云盘的续费时长。<br>在云盘与挂载的实例一起续费的场景下，可以指定参数CurInstanceDeadline，此时云盘会按对齐到实例续费后的到期时间来续费。
	DiskChargePrepaid *DiskChargePrepaid `json:"DiskChargePrepaid" name:"DiskChargePrepaid"`

	// 云硬盘ID， 通过[DescribeDisks](/document/product/362/16315)接口查询。
	DiskId *string `json:"DiskId" name:"DiskId"`
}

func (r *RenewDiskRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RenewDiskRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RenewDiskResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RenewDiskResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RenewDiskResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ResizeDiskRequest struct {
	*tchttp.BaseRequest

	// 云硬盘ID， 通过[DescribeDisks](/document/product/362/16315)接口查询。
	DiskId *string `json:"DiskId" name:"DiskId"`

	// 云硬盘扩容后的大小，单位为GB，必须大于当前云硬盘大小。云盘大小取值范围参见云硬盘[产品分类](/document/product/362/2353)的说明。
	DiskSize *uint64 `json:"DiskSize" name:"DiskSize"`
}

func (r *ResizeDiskRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ResizeDiskRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ResizeDiskResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ResizeDiskResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ResizeDiskResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Snapshot struct {

	// 快照ID。
	SnapshotId *string `json:"SnapshotId" name:"SnapshotId"`

	// 快照所在的位置。
	Placement *Placement `json:"Placement" name:"Placement"`

	// 创建此快照的云硬盘类型。取值范围：<br><li>SYSTEM_DISK：系统盘<br><li>DATA_DISK：数据盘。
	DiskUsage *string `json:"DiskUsage" name:"DiskUsage"`

	// 创建此快照的云硬盘ID。
	DiskId *string `json:"DiskId" name:"DiskId"`

	// 创建此快照的云硬盘大小，单位GB。
	DiskSize *uint64 `json:"DiskSize" name:"DiskSize"`

	// 快照的状态。取值范围：<br><li>NORMAL：正常<br><li>CREATING：创建中<br><li>ROLLBACKING：回滚中<br><li>COPYING_FROM_REMOTE：跨地域复制快照拷贝中。
	SnapshotState *string `json:"SnapshotState" name:"SnapshotState"`

	// 快照名称，用户自定义的快照别名。调用[ModifySnapshotAttribute](/document/product/362/15650)可修改此字段。
	SnapshotName *string `json:"SnapshotName" name:"SnapshotName"`

	// 快照创建进度百分比，快照创建成功后此字段恒为100。
	Percent *uint64 `json:"Percent" name:"Percent"`

	// 快照的创建时间。
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 快照到期时间。如果快照为永久保留，此字段为空。
	DeadlineTime *string `json:"DeadlineTime" name:"DeadlineTime"`

	// 是否为加密盘创建的快照。取值范围：<br><li>true：该快照为加密盘创建的<br><li>false:非加密盘创建的快照。
	Encrypt *bool `json:"Encrypt" name:"Encrypt"`

	// 是否为永久快照。取值范围：<br><li>true：永久快照<br><li>false：非永久快照。
	IsPermanent *bool `json:"IsPermanent" name:"IsPermanent"`

	// 快照正在跨地域复制的目的地域，默认取值为[]。
	CopyingToRegions []*string `json:"CopyingToRegions" name:"CopyingToRegions" list`

	// 是否为跨地域复制的快照。取值范围：<br><li>true：表示为跨地域复制的快照。<br><li>false:本地域的快照。
	CopyFromRemote *bool `json:"CopyFromRemote" name:"CopyFromRemote"`
}

type Tag struct {

	// 标签健。
	Key *string `json:"Key" name:"Key"`

	// 标签值。
	Value *string `json:"Value" name:"Value"`
}

type TerminateDisksRequest struct {
	*tchttp.BaseRequest

	// 需退还的云盘ID列表。
	DiskIds []*string `json:"DiskIds" name:"DiskIds" list`
}

func (r *TerminateDisksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *TerminateDisksRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type TerminateDisksResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *TerminateDisksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *TerminateDisksResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}
