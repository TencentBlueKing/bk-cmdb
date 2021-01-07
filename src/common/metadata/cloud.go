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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
)

// 云账户
type CloudAccount struct {
	AccountName string    `json:"bk_account_name" bson:"bk_account_name"`
	CloudVendor string    `json:"bk_cloud_vendor" bson:"bk_cloud_vendor"`
	AccountID   int64     `json:"bk_account_id" bson:"bk_account_id"`
	SecretID    string    `json:"bk_secret_id" bson:"bk_secret_id"`
	SecretKey   string    `json:"bk_secret_key" bson:"bk_secret_key"`
	Description string    `json:"bk_description" bson:"bk_description"`
	OwnerID     string    `json:"bk_supplier_account" bson:"bk_supplier_account"`
	Creator     string    `json:"bk_creator" bson:"bk_creator"`
	LastEditor  string    `json:"bk_last_editor" bson:"bk_last_editor"`
	CreateTime  time.Time `json:"create_time" bson:"create_time"`
	LastTime    time.Time `json:"last_time" bson:"last_time"`
}

func (c *CloudAccount) ToMapStr() mapstr.MapStr {
	return mapstr.SetValueToMapStrByTags(c)
}

func (c *CloudAccount) Validate() (rawError errors.RawErrorInfo) {
	if c.AccountName == "" {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"bk_account_name"},
		}
	}

	if c.CloudVendor == "" {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"bk_cloud_vendor"},
		}
	}

	if !util.InStrArr(SupportedCloudVendors, c.CloudVendor) {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCloudVendorNotSupport,
		}
	}

	if c.SecretID == "" {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"bk_secret_id"},
		}
	}

	if c.SecretKey == "" {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"bk_secret_key"},
		}
	}

	return errors.RawErrorInfo{}
}

// 带有额外信息的云账户
type CloudAccountWithExtraInfo struct {
	CloudAccount `json:",inline"`
	// 是否能删除账户，只有该账户下不存在同步任务了，才能删除，此时才能为true，否则为false
	CanDeleteAccount bool `json:"bk_can_delete_account" bson:"bk_can_delete_account"`
}

// 云厂商
// 和属性表中的bk_cloud_vendor值相对应
const (
	AWS          string = "1"
	TencentCloud string = "2"
)

// 支持的云厂商
// 实现了相应的云厂商插件
var SupportedCloudVendors = []string{AWS, TencentCloud}

// 云同步任务同步状态
const (
	CloudSyncSuccess    string = "cloud_sync_success"
	CloudSyncFail       string = "cloud_sync_fail"
	CloudSyncInProgress string = "cloud_sync_in_progress"
)

// 云厂商账户配置
type CloudAccountConf struct {
	AccountID  int64  `json:"bk_account_id" bson:"bk_account_id"`
	VendorName string `json:"bk_cloud_vendor" bson:"bk_cloud_vendor"`
	SecretID   string `json:"bk_secret_id" bson:"bk_secret_id"`
	SecretKey  string `json:"bk_secret_key" bson:"bk_secret_key"`
}

type SearchCloudOption struct {
	Condition mapstr.MapStr `json:"condition" bson:"condition" field:"condition"`
	Page      BasePage      `json:"page" bson:"page" field:"page"`
	Fields    []string      `json:"fields,omitempty" bson:"fields,omitempty"`
	// 对于condition里的属性值是否模糊匹配，默认为false，即不采用模糊匹配，而使用精确匹配
	IsFuzzy bool `json:"is_fuzzy" bson:"is_fuzzy"`
}

type SearchSyncTaskOption struct {
	SearchCloudOption `json:",inline"`
	// 是否实时获取云厂商vpc下最新的主机数
	LastestHostCount bool `json:"latest_hostcount" bson:"latest_host_count"`
}

type SearchVpcHostCntOption struct {
	RegionVpcs []RegionVpc
}

type RegionVpc struct {
	Region string `json:"bk_region"`
	VpcID  string `json:"bk_vpc_id"`
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

type DeleteDestroyedHostRelatedOption struct {
	HostIDs []int64 `json:"host_ids" bson:"host_ids"`
}

type MultipleSyncHistory struct {
	Count int64         `json:"count"`
	Info  []SyncHistory `json:"info"`
}

type MultipleCloudAccount struct {
	Count int64                       `json:"count"`
	Info  []CloudAccountWithExtraInfo `json:"info"`
}

type MultipleCloudAccountConf struct {
	Count int64              `json:"count"`
	Info  []CloudAccountConf `json:"info"`
}

type CloudAccountVerify struct {
	SecretID    string `json:"bk_secret_id"`
	SecretKey   string `json:"bk_secret_key"`
	CloudVendor string `json:"bk_cloud_vendor"`
}

type SearchAccountValidityOption struct {
	AccountIDs []int64 `json:"account_ids" bson:"account_ids"`
}

func (s *SearchAccountValidityOption) Validate() (rawError errors.RawErrorInfo) {
	if len(s.AccountIDs) == 0 || len(s.AccountIDs) > common.BKMaxInstanceLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrArrayLengthWrong,
			Args:    []interface{}{"account_ids", common.BKMaxInstanceLimit},
		}
	}

	return errors.RawErrorInfo{}
}

type AccountValidityInfo struct {
	AccountID int64  `json:"bk_account_id" bson:"bk_account_id"`
	ErrMsg    string `json:"err_msg" bson:"err_msg"`
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
	TaskID            int64          `json:"bk_task_id" bson:"bk_task_id"`
	TaskName          string         `json:"bk_task_name" bson:"bk_task_name"`
	ResourceType      string         `json:"bk_resource_type" bson:"bk_resource_type"`
	AccountID         int64          `json:"bk_account_id" bson:"bk_account_id"`
	CloudVendor       string         `json:"bk_cloud_vendor" bson:"bk_cloud_vendor"`
	SyncStatus        string         `json:"bk_sync_status" bson:"bk_sync_status"`
	OwnerID           string         `json:"bk_supplier_account" bson:"bk_supplier_account"`
	StatusDescription SyncStatusDesc `json:"bk_status_description" bson:"bk_status_description"`
	LastSyncTime      *time.Time     `json:"bk_last_sync_time" bson:"bk_last_sync_time"`
	SyncAll           bool           `json:"bk_sync_all" bson:"bk_sync_all"`
	SyncAllDir        int64          `json:"bk_sync_all_dir" bson:"bk_sync_all_dir"`
	SyncVpcs          []VpcSyncInfo  `json:"bk_sync_vpcs" bson:"bk_sync_vpcs"`
	Creator           string         `json:"bk_creator" bson:"bk_creator"`
	LastEditor        string         `json:"bk_last_editor" bson:"bk_last_editor"`
	CreateTime        time.Time      `json:"create_time" bson:"create_time"`
	LastTime          time.Time      `json:"last_time" bson:"last_time"`
}

// ToMapStr to mapstr
func (c *CloudSyncTask) ToMapStr() mapstr.MapStr {
	return mapstr.SetValueToMapStrByTags(c)
}

type VpcSyncInfo struct {
	VpcID        string `json:"bk_vpc_id" bson:"bk_vpc_id"`
	VpcName      string `json:"bk_vpc_name" bson:"bk_vpc_name"`
	Region       string `json:"bk_region" bson:"bk_region"`
	VpcHostCount int64  `json:"bk_host_count" bson:"bk_host_count"`
	SyncDir      int64  `json:"bk_sync_dir,omitempty" bson:"bk_sync_dir,omitempty"`
	CloudID      int64  `json:"bk_cloud_id" bson:"bk_cloud_id"`
	// 该vpc在云端是否被销毁
	Destroyed bool `json:"destroyed" bson:"destroyed"`
}

type MultipleCloudSyncTask struct {
	Count int64           `json:"count"`
	Info  []CloudSyncTask `json:"info"`
}

type VpcHostCntResult struct {
	Count int64         `json:"count"`
	Info  []VpcSyncInfo `json:"info"`
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
	PrivateIp     string `json:"bk_host_innerip" bson:"bk_host_innerip"`
	PublicIp      string `json:"bk_host_outerip" bson:"bk_host_outerip"`
	InstanceState string `json:"bk_cloud_host_status" bson:"bk_cloud_host_status"`
	VpcId         string `json:"bk_vpc_id" bson:"bk_vpc_id"`
}

// 云主机同步时的资源数据
type CloudHostResource struct {
	HostResource  []*VpcInstances
	DestroyedVpcs []*VpcSyncInfo
	TaskID        int64             `json:"bk_task_id" bson:"bk_task_id"`
	AccountConf   *CloudAccountConf `json:"account_conf" bson:"account_conf"`
}

type VpcInstances struct {
	Vpc       *VpcSyncInfo
	CloudID   int64 `json:"bk_cloud_id" bson:"bk_cloud_id"`
	Instances []*Instance
}

type CloudHost struct {
	Instance   `json:",inline"`
	CloudID    int64  `json:"bk_cloud_id" bson:"bk_cloud_id"`
	SyncDir    int64  `json:"bk_sync_dir,omitempty" bson:"bk_sync_dir,omitempty"`
	HostID     int64  `json:"bk_host_id" bson:"bk_host_id"`
	VendorName string `json:"bk_cloud_vendor" bson:"bk_cloud_vendor"`
}

type HostSyncInfo struct {
	HostID        int64     `json:"bk_host_id" bson:"bk_host_id"`
	CloudID       int64     `json:"bk_cloud_id" bson:"bk_cloud_id"`
	InstanceId    string    `json:"bk_cloud_inst_id" bson:"bk_cloud_inst_id"`
	PrivateIp     string    `json:"bk_host_innerip" bson:"bk_host_innerip"`
	PublicIp      string    `json:"bk_host_outerip" bson:"bk_host_outerip"`
	InstanceState string    `json:"bk_cloud_host_status" bson:"bk_cloud_host_status"`
	OwnerID       string    `json:"bk_supplier_account" bson:"bk_supplier_account"`
	CreateTime    time.Time `json:"create_time" bson:"create_time"`
	LastTime      time.Time `json:"last_time" bson:"last_time"`
}

// 云区域
type CloudArea struct {
	CloudID     int64     `json:"bk_cloud_id" bson:"bk_cloud_id"`
	CloudName   string    `json:"bk_cloud_name" bson:"bk_cloud_name"`
	Status      int       `json:"bk_status" bson:"bk_status"`
	CloudVendor string    `json:"bk_cloud_vendor" bson:"bk_cloud_vendor"`
	OwnerID     string    `json:"bk_supplier_account" bson:"bk_supplier_account"`
	VpcID       string    `json:"bk_vpc_id" bson:"bk_vpc_id"`
	VpcName     string    `json:"bk_vpc_name" bson:"bk_vpc_name"`
	Region      string    `json:"bk_region" bson:"bk_region"`
	AccountID   int64     `json:"bk_account_id" bson:"bk_account_id"`
	Creator     string    `json:"bk_creator" bson:"bk_creator"`
	LastEditor  string    `json:"bk_last_editor" bson:"bk_last_editor"`
	CreateTime  time.Time `json:"create_time" bson:"create_time"`
	LastTime    time.Time `json:"last_time" bson:"last_time"`
}

type SyncRegion struct {
	RegionId    string `json:"bk_region" bson:"bk_region"`
	RegionName  string `json:"bk_region_name" bson:"bk_region_name"`
	RegionState string `json:"bk_region_state" bson:"bk_region_state"`
	HostCount   int64  `json:"bk_host_count" bson:"bk_host_count"`
}

// 同步历史记录
type SyncHistory struct {
	HistoryID         int64          `json:"bk_history_id" bson:"bk_history_id"`
	TaskID            int64          `json:"bk_task_id" bson:"bk_task_id"`
	SyncStatus        string         `json:"bk_sync_status" bson:"bk_sync_status"`
	OwnerID           string         `json:"bk_supplier_account" bson:"bk_supplier_account"`
	StatusDescription SyncStatusDesc `json:"bk_status_description" bson:"bk_status_description"`
	Detail            SyncDetail     `json:"bk_detail" bson:"bk_detail"`
	CreateTime        time.Time      `json:"create_time" bson:"create_time"`
}

type SyncStatusDesc struct {
	CostTime  float64 `json:"cost_time" bson:"cost_time"`
	ErrorInfo string  `json:"error_info" bson:"error_info"`
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
	SuccessInfo       SyncSuccessInfo `json:"success_info" bson:"success_info"`
	FailInfo          SyncFailInfo    `json:"fail_info" bson:"fail_info"`
	Detail            SyncDetail      `json:"detail" bson:"detail"`
	SyncStatus        string          `json:"bk_sync_status" bson:"bk_sync_status"`
	StatusDescription SyncStatusDesc  `json:"bk_status_description" bson:"bk_status_description"`
}

type SecretKeyResult struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Result  bool          `json:"result"`
	Data    SecretKeyInfo `json:"data"`
}

type SecretKeyInfo struct {
	Content SecretContent `json:"content"`
}

type SecretContent struct {
	SecretKey string `json:"secret_key"`
}
