/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

import (
	"configcenter/src/common/mapstr"
	"net/http"
	"time"
)

type DeleteHostBatchOpt struct {
	HostID string `json:"bk_host_id"`
}

type HostInstanceProperties struct {
	PropertyID    string      `json:"bk_property_id"`
	PropertyName  string      `json:"bk_property_name"`
	PropertyValue interface{} `json:"bk_property_value"`
}

type HostInstancePropertiesResult struct {
	BaseResp `json:",inline"`
	Data     []HostInstanceProperties `json:"data"`
}

type HostSnapResult struct {
	BaseResp `json:",inline"`
	Data     map[string]interface{} `json:"data"`
}

type UserCustomQueryDetailResult struct {
	BaseResp `json:",inline"`
	Data     map[string]interface{} `json:"data"`
}

type HostInputType string

const (
	ExecelType  HostInputType = "excel"
	CollectType HostInputType = "collect"
)

type HostList struct {
	ApplicationID int64                            `json:"bk_biz_id"`
	HostInfo      map[int64]map[string]interface{} `json:"host_info"`
	SupplierID    int64                            `json:"bk_supplier_id"`
	InputType     HostInputType                    `json:"input_type"`
}

type AddHostFromAgentHostList struct {
	HostInfo map[string]interface{} `json:"host_info"`
}

type HostSyncList struct {
	ApplicationID int64                            `json:"bk_biz_id"`
	HostInfo      map[int64]map[string]interface{} `json:"host_info"`
	SupplierID    int64                            `json:"bk_supplier_id"`
	ModuleID      []int64                          `json:"bk_module_id"`
	InputType     HostInputType                    `json:"input_type"`
}

type HostsModuleRelation struct {
	ApplicationID int64   `json:"bk_biz_id"`
	HostID        []int64 `json:"bk_host_id"`
	ModuleID      []int64 `json:"bk_module_id"`
	IsIncrement   bool    `json:"is_increment"`
}

type HostToAppModule struct {
	Ips         []string `json:"ips"`
	HostName    []string `json:"bk_host_name"`
	ModuleName  string   `json:"bk_module_name"`
	SetName     string   `json:"bk_set_name"`
	AppName     string   `json:"bk_biz_name"`
	OsType      string   `json:"bk_os_type"`
	OwnerID     string   `json:"bk_supplier_account"`
	PlatID      int64    `json:"bk_cloud_id"`
	IsIncrement bool     `json:"is_increment"`
}

type HostCommonSearch struct {
	AppID     int64             `json:"bk_biz_id,omitempty"`
	Ip        IPInfo            `json:"ip"`
	Condition []SearchCondition `json:"condition"`
	Page      BasePage          `json:"page"`
	Pattern   string            `json:"pattern,omitempty"`
}

type HostModuleFind struct {
	ModuleIDS []int64  `json:"bk_module_ids"`
	Metadata  Metadata `json:"metadata"`
	Page      BasePage `json:"page"`
}

//ip search info
type IPInfo struct {
	Data  []string `json:"data"`
	Exact int64    `json:"exact"`
	Flag  string   `json:"flag"`
}

//search condition
type SearchCondition struct {
	Fields    []string        `json:"fields"`
	Condition []ConditionItem `json:"condition"`
	ObjectID  string          `json:"bk_obj_id"`
}

type SearchHost struct {
	Count int             `json:"count"`
	Info  []mapstr.MapStr `json:"info"`
}

type SearchHostResult struct {
	BaseResp `json:",inline"`
	Data     SearchHost `json:"data"`
}

type HostCloneInputParams struct {
	OrgIP  string `json:"bk_org_ip"`
	DstIP  string `json:"bk_dst_ip"`
	AppID  int64  `json:"bk_biz_id"`
	PlatID int64  `json:"bk_cloud_id"`
}

type SetHostConfigParams struct {
	ApplicationID int64 `json:"bk_biz_id"`
	SetID         int64 `json:"bk_set_id"`
	ModuleID      int64 `json:"bk_module_id"`
}

type CloneHostPropertyParams struct {
	AppID   int64  `json:"bk_biz_id"`
	OrgIP   string `json:"bk_org_ip"`
	DstIP   string `json:"bk_dst_ip"`
	CloudID int64  `json:"bk_cloud_id"`
}

type CloudTaskList struct {
	User            string `json:"bk_user"`
	TaskName        string `json:"bk_task_name"`
	TaskID          int64  `json:"bk_task_id"`
	AccountType     string `json:"bk_account_type"`
	AccountAdmin    string `json:"bk_account_admin"`
	PeriodType      string `json:"bk_period_type"`
	Period          string `json:"bk_period"`
	LastSyncTime    string `json:"bk_last_sync_time"`
	ObjID           string `json:"bk_obj_id"`
	Status          bool   `json:"bk_status"`
	ResourceConfirm bool   `json:"bk_confirm"`
	AttrConfirm     bool   `json:"bk_attr_confirm"`
	SecretID        string `json:"bk_secret_id"`
	SecretKey       string `json:"bk_secret_key"`
}

type ResourceConfirm struct {
	ObjID        string          `json:"bk_obj_id"`
	ResourceName []mapstr.MapStr `json:"bk_resource_name"`
	SourceType   string          `json:"bk_source_type"`
	SourceName   string          `json:"bk_source_name"`
	CreateTime   string          `json:"bk_create_time"`
	TaskID       string          `json:"bk_task_id"`
	ResourceID   int64           `json:"bk_resource_id"`
	ConfirmType  string          `json:"bk_confirm_type`
	Incharge     string          `json:"bk_in_charge"`
}

type CloudHistory struct {
	ObjID       string `json:"bk_obj_id"`
	Status      string `json:"bk_status"`
	TimeConsume string `json:"bk_time_consume"`
	NewAdd      int    `json:"new_add"`
	AttrChanged int    `json:"attr_changed"`
	StartTime   string `json:"bk_start_time"`
	TaskID      int64  `json:"bk_task_id"`
	HistoryID   int64  `json:"bk_history_id"`
	FailReason  string `json:"fail_reason"`
}

type DeleteCloudTask struct {
	TaskID int64 `json:"bk_task_id"`
}

type RegionResponse struct {
	Response RegionSet `json:"Response"`
}

type RegionSet struct {
	Data []Region `json:"RegionSet"`
}

type Region struct {
	Region string `json:"Region"`
}

type HostResponse struct {
	HostResponse InstanceSet `json:"Response"`
}

type InstanceSet struct {
	InstanceSet []CloudHostInfo `json:"InstanceSet"`
}

type CloudHostInfo struct {
	PrivateIpAddresses []string `json:"PrivateIpAddresses"`
	PublicIpAddresses  []string `json:"PublicIpAddresses"`
	OsName             string   `json:"OsName"`
}

type TaskInfo struct {
	Args        CloudTaskInfo
	Method      string
	NextTrigger int64
	ManagerChn  chan bool
}

type CloudSyncRedisPendingStart struct {
	NewHeader    http.Header `json:"new_header"`
	TaskID       int64       `json:"bk_task_id"`
	TaskItemInfo TaskInfo    `json:"task_item_info"`
	OwnerID      string      `json:"bk_supplier_account"`
}

type CloudSyncRedisAlreadyStarted struct {
	LastSyncTime time.Time   `json:"last_sync_time"`
	NewHeader    http.Header `json:"new_header"`
	TaskID       int64       `json:"bk_task_id"`
	TaskItemInfo TaskInfo    `json:"task_item_info"`
	OwnerID      string      `json:"bk_supplier_account"`
}

type CloudSyncRedisPendingStop struct {
	TaskID  int64  `json:"bk_task_id"`
	OwnerID string `json:"bk_supplier_account"`
}

// TransferHostAcrossBusinessParameter Transfer host across business request parameter
type TransferHostAcrossBusinessParameter struct {
	SrcAppID       int64   `json:"src_bk_biz_id"`
	DstAppID       int64   `json:"dst_bk_biz_id"`
	HostID         int64   `json:"bk_host_id"`
	DstModuleIDArr []int64 `json:"bk_module_ids"`
}

// HostModuleRelationParameter host and module  relation parameter
type HostModuleRelationParameter struct {
	AppID  int64   `json:"bk_biz_id"`
	HostID []int64 `json:"bk_host_id"`
}

// DeleteHostFromBizParameter delete host from business
type DeleteHostFromBizParameter struct {
	AppID     int64   `json:"bk_biz_id"`
	HostIDArr []int64 `json:"bk_host_ids"`
}

// OperaterException synchronize result
type OperaterException struct {
	BaseResp `json:",inline"`
	Data     []ExceptionResult `json:"data"`
}
