/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */
package metadata

import (
	"configcenter/src/common/mapstr"
)

// 云账户
type CloudAccount struct {
	AccountName string `json:"bk_account_name" bson:"bk_account_name"`
	CloudVendor string `json:"bk_cloud_vendor" bson:"bk_cloud_vendor"`
	AccountID   int64  `json:"bk_account_id" bson:"bk_account_id"`
	SecretID    string `json:"bk_secret_id" bson:"bk_secret_id"`
	SecretKey   string `json:"bk_secret_key" bson:"bk_secret_key"`
	Description string `json:"bk_description" bson:"bk_description"`
	// 是否能删除账户，只有该账户下不存在同步任务了，才能删除，此时才能为true，否则为false
	CanDeleteAccount bool   `json:"bk_can_delete_account" bson:"bk_can_delete_account"`
	OwnerID          string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	Creator          string `json:"bk_creator" bson:"bk_creator"`
	LastEditor       string `json:"bk_last_editor" bson:"bk_last_editor"`
	CreateTime       Time   `json:"create_time" bson:"create_time"`
	LastTime         Time   `json:"last_time" bson:"last_time"`
}

// 云厂商
const (
	AWS          string = "aws"
	TencentCloud string = "tencent_cloud"
)

// 同步状态
const (
	CloudSyncSuccess    string = "cloud_sync_success"
	CloudSyncFail       string = "cloud_sync_fail"
	CloudSyncInProgress string = "cloud_sync_in_progress"
)

var SupportedCloudVendors = []string{"aws", "tencent_cloud"}

type SearchCloudOption struct {
	Condition mapstr.MapStr `json:"condition" bson:"condition" field:"condition"`
	Page      BasePage      `json:"page" bson:"page" field:"page"`
	Fields    []string      `json:"fields,omitempty" bson:"fields,omitempty"`
	// 对于condition里的属性值是否精确匹配，默认为false，即使用模糊匹配和忽略大小写
	Exact bool `json:"exact" bson:"exact"`
}

type SearchVpcOption struct {
	Region string `json:"bk_region"`
}

type SearchSyncRegionOption struct {
	AccountID int64 `json:"bk_account_id" bson:"bk_account_id"`
	// 是否返回地域下的主机数，返回主机数会导致请求更耗时，默认为false
	WithHostCount bool `json:"with_host_count" bson:"with_host_count"`
}

type SearchSyncHistoryOption struct {
	SearchCloudOption `json:",inline"`
	TaskID            int64  `json:"bk_task_id" bson:"bk_task_id"`
	StarTime          string `json:"start_time" bson:"start_time"`
	EndTime           string `json:"end_time" bson:"end_time"`
}

type MultipleSyncHistory struct {
	Count int64         `json:"count"`
	Info  []SyncHistory `json:"info"`
}

type MultipleSyncRegion struct {
	Count int64        `json:"count"`
	Info  []SyncRegion `json:"info"`
}

type MultipleCloudAccount struct {
	Count int64          `json:"count"`
	Info  []CloudAccount `json:"info"`
}

type CloudAccountVerify struct {
	SecretID    string `json:"bk_secret_id"`
	SecretKey   string `json:"bk_secret_key"`
	CloudVendor string `json:"bk_cloud_vendor"`
}

type VpcInfo struct {
	VpcName string `json:"bk_vpc_name"`
	VpcID   string `json:"bk_vpc_id"`
	Region  string `json:"bk_region"`
}

type SearchVpcResult struct {
	Count string            `json:"count"`
	Info  []VpcInstancesCnt `json:"info"`
}

type VpcInstancesCnt struct {
	VpcId     string `json:"bk_vpc_id"`
	VpcName   string `json:"bk_vpc_name"`
	Region    string `json:"bk_region"`
	HostCount int64  `json:"bk_host_count"`
}

// 云同步任务
type CloudSyncTask struct {
	TaskID            int64         `json:"bk_task_id" bson:"bk_task_id"`
	TaskName          string        `json:"bk_task_name" bson:"bk_task_name"`
	ResourceType      string        `json:"bk_resource_type" bson:"bk_resource_type"`
	AccountID         int64         `json:"bk_account_id" bson:"bk_account_id"`
	CloudVendor       string        `json:"bk_cloud_vendor" bson:"bk_cloud_vendor"`
	SyncStatus        string        `json:"bk_sync_status" bson:"bk_sync_status"`
	OwnerID           string        `json:"bk_supplier_account" bson:"bk_supplier_account"`
	StatusDescription string        `json:"bk_status_description" bson:"bk_status_description"`
	LastSyncTime      string        `json:"bk_last_sync_time" bson:"bk_last_sync_time"`
	SyncAll           bool          `json:"bk_sync_all" bson:"bk_sync_all"`
	SyncAllDir        int64         `json:"bk_sync_all_dir" bson:"bk_sync_all_dir"`
	SyncVpcs          []VpcSyncInfo `json:"bk_sync_vpcs" bson:"bk_sync_vpcs"`
	Creator           string        `json:"bk_creator" bson:"bk_creator"`
	CreateTime        Time          `json:"create_time" bson:"create_time"`
	LastEditor        string        `json:"bk_last_editor" bson:"bk_last_editor"`
	LastTime          Time          `json:"last_time" bson:"last_time"`
}

type VpcSyncInfo struct {
	VpcID        string `json:"bk_vpc_id" bson:"bk_vpc_id"`
	VpcName      string `json:"bk_vpc_name" bson:"bk_vpc_name"`
	Region       string `json:"bk_region" bson:"bk_region"`
	VpcHostCount int64  `json:"bk_host_count" bson:"bk_host_count"`
	SyncDir      int64  `json:"bk_sync_dir,omitempty" bson:"bk_sync_dir,omitempty"`
}

type MultipleCloudSyncTask struct {
	Count int64           `json:"count"`
	Info  []CloudSyncTask `json:"info"`
}

type VpcHostCntResult struct {
	Count int64         `json:"count"`
	Info  []VpcSyncInfo `json:"info"`
}

type RegionsInfo struct {
	Count     int64     `json:"count" bson:"count"`
	RegionSet []*Region `json:"region_set" bson:"region_set"`
}

type Region struct {
	RegionId    string `json:"bk_region" bson:"bk_region"`
	RegionName  string `json:"bk_region_name" bson:"bk_region_name"`
	RegionState string `json:"bk_region_state" bson:"bk_region_state"`
}

type VpcsInfo struct {
	Count  int64  `json:"count" bson:"count"`
	VpcSet []*Vpc `json:"vpc_set" bson:"vpc_set"`
}

type Vpc struct {
	VpcId   string `json:"bk_vpc_id" bson:"bk_vpc_id"`
	VpcName string `json:"bk_vpc_name" bson:"bk_vpc_name"`
}

type InstancesInfo struct {
	Count       int64       `json:"count" bson:"count"`
	InstanceSet []*Instance `json:"instance_set" bson:"instance_set"`
}

type Instance struct {
	InstanceId    string `json:"bk_cloud_inst_id" bson:"bk_cloud_inst_id"`
	InstanceName  string `json:"bk_host_name" bson:"bk_host_name"`
	PrivateIp     string `json:"bk_host_innerip" bson:"bk_host_innerip"`
	PublicIp      string `json:"bk_host_outerip" bson:"bk_host_outerip"`
	InstanceState string `json:"bk_cloud_host_status" bson:"bk_cloud_host_status"`
	VpcId         string `json:"bk_vpc_id" bson:"bk_vpc_id"`
	OsName        string `json:"bk_os_name" bson:"bk_os_name"`
}

// 云主机资源
type CloudHostResource struct {
	HostResource []*VpcInstances
	TaskID       int64 `json:"bk_task_id" bson:"bk_task_id"`
}

type VpcInstances struct {
	Vpc       *VpcSyncInfo
	CloudID   int64 `json:"bk_cloud_id" bson:"bk_cloud_id"`
	Instances []*Instance
}

type CloudHost struct {
	Instance `json:",inline"`
	CloudID  int64 `json:"bk_cloud_id" bson:"bk_cloud_id"`
	SyncDir  int64 `json:"bk_sync_dir,omitempty" bson:"bk_sync_dir,omitempty"`
	HostID   int64 `json:"bk_host_id" bson:"bk_host_id"`
}

type HostSyncInfo struct {
	HostID        int64  `json:"bk_host_id" bson:"bk_host_id"`
	CloudID       int64  `json:"bk_cloud_id" bson:"bk_cloud_id"`
	InstanceId    string `json:"bk_cloud_inst_id" bson:"bk_cloud_inst_id"`
	InstanceName  string `json:"bk_host_name" bson:"bk_host_name"`
	PrivateIp     string `json:"bk_host_innerip" bson:"bk_host_innerip"`
	PublicIp      string `json:"bk_host_outerip" bson:"bk_host_outerip"`
	InstanceState string `json:"bk_cloud_host_status" bson:"bk_cloud_host_status"`
	OsName        string `json:"bk_os_name" bson:"bk_os_name"`
	CreateTime    Time   `json:"create_time" bson:"create_time"`
	LastTime      Time   `json:"last_time" bson:"last_time"`
}

// 云区域
type CloudArea struct {
	CloudID     int64  `json:"bk_cloud_id" bson:"bk_cloud_id"`
	CloudName   string `json:"bk_cloud_name" bson:"bk_cloud_name"`
	Status      int    `json:"bk_status" bson:"bk_status"`
	CloudVendor string `json:"bk_cloud_vendor" bson:"bk_cloud_vendor"`
	OwnerID     string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	VpcID       string `json:"bk_vpc_id" bson:"bk_vpc_id"`
	VpcName     string `json:"bk_vpc_name" bson:"bk_vpc_name"`
	Region      string `json:"bk_region" bson:"bk_region"`
	AccountID   int64  `json:"bk_account_id" bson:"bk_account_id"`
	Creator     string `json:"bk_creator" bson:"bk_creator"`
	CreateTime  Time   `json:"create_time" bson:"create_time"`
	LastEditor  string `json:"bk_last_editor" bson:"bk_last_editor"`
	LastTime    Time   `json:"last_time" bson:"last_time"`
}

type SyncRegion struct {
	RegionId    string `json:"bk_region" bson:"bk_region"`
	RegionName  string `json:"bk_region_name" bson:"bk_region_name"`
	RegionState string `json:"bk_region_state" bson:"bk_region_state"`
	HostCount   int64  `json:"bk_host_count" bson:"bk_host_count"`
}

// 同步历史记录
type SyncHistory struct {
	HistoryID         int64      `json:"bk_history_id" bson:"bk_history_id"`
	TaskID            int64      `json:"bk_task_id" bson:"bk_task_id"`
	SyncStatus        string     `json:"bk_sync_status" bson:"bk_sync_status"`
	StatusDescription string     `json:"bk_status_description" bson:"bk_status_description"`
	OwnerID           string     `json:"bk_supplier_account" bson:"bk_supplier_account"`
	Detail            SyncDetail `json:"bk_detail" bson:"bk_detail"`
	CreateTime        Time       `json:"create_time" bson:"create_time"`
}

type SyncDetail struct {
	NewAdd SyncSuccessInfo `json:"new_add" bson:"new_add"`
	Update SyncSuccessInfo `json:"update" bson:"update"`
}

type SyncSuccessInfo struct {
	Count int64    `json:"count" bson:"count"`
	IPs   []string `json:"ips" bson:"ips"`
}

type SyncFailInfo struct {
	Count   int64             `json:"count" bson:"count"`
	IPError map[string]string `json:"ip_error" bson:"ip_error"`
}

type SyncResult struct {
	SuccessInfo SyncSuccessInfo `json:"success_info" bson:"success_info"`
	FailInfo    SyncFailInfo    `json:"fail_info" bson:"fail_info"`
	Detail      SyncDetail      `json:"detail" bson:"detail"`
}
