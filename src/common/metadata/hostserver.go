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
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"

	"github.com/imdario/mergo"
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

type RemoveHostsFromModuleOption struct {
	ApplicationID int64 `json:"bk_biz_id"`
	HostID        int64 `json:"bk_host_id"`
	ModuleID      int64 `json:"bk_module_id"`
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

type ListHostByTopoNodeParameter struct {
	Metadata           Metadata                 `json:"metadata"`
	SetIDs             []int64                  `json:"bk_set_ids"`
	ModuleIDs          []int64                  `json:"bk_module_ids"`
	StdProperty        HostStdPropertyCondition `json:"std_property"`
	ExcludeStdProperty HostStdPropertyCondition `json:"exclude_std_property"`
	Page               BasePage                 `json:"page"`
}

type TimeRange struct {
	Start *time.Time
	End   *time.Time
}
type HostStdPropertyCondition struct {
	Memory          []int64   `json:"bk_mem"`
	CPUMhz          []int64   `json:"bk_cpu_mhz"`
	InnerIP         []string  `json:"bk_host_innerip"`
	HostID          []int64   `json:"bk_host_id"`
	OsVersion       []string  `json:"bk_os_version"`
	CPUModule       []string  `json:"bk_cpu_module" `
	Sn              []string  `json:"bk_sn" `
	Operator        []string  `json:"operator"`
	Disk            []int64   `json:"bk_disk"`
	Sla             []string  `json:"bk_sla"`
	HostName        []string  `json:"bk_host_name"`
	AssetID         []int64   `json:"bk_asset_id"`
	OuterIP         []string  `json:"bk_host_outerip" `
	ImportFrom      []string  `json:"import_from"`
	CPU             []int     `json:"bk_cpu"`
	ProvinceName    []string  `json:"bk_province_name"`
	OSBit           []int     `json:"bk_os_bit"`
	OSType          []string  `json:"bk_os_type"`
	Comment         []string  `json:"bk_comment"`
	Mac             []string  `json:"bk_mac"`
	CloudID         []int64   `json:"bk_cloud_id"`
	OSName          []string  `json:"bk_os_name"`
	StateName       []string  `json:"bk_state_name"`
	OuterMac        []string  `json:"bk_outer_mac"`
	ServiceTerm     []string  `json:"bk_service_term"`
	BakOperator     []string  `json:"bk_bak_operator"`
	ISPName         []string  `json:"bk_isp_name"`
	SupplierAccount []string  `json:"bk_supplier_account"`
	CreateTime      TimeRange `json:"create_time"`
	LastTime        TimeRange `json:"last_time"`
}

func (hspc HostStdPropertyCondition) Validate(operator string) (string, error) {
	switch operator {
	case common.BKDBIN:
	case common.BKDBNIN:
	default:
		return "operator", fmt.Errorf("unexpected operator %s", operator)
	}
	return "", nil
}

func (hspc HostStdPropertyCondition) ToHostFilter(operator string) map[string]interface{} {
	filter := make(map[string]interface{})
	if hspc.Memory != nil {
		filter["bk_mem"] = map[string]interface{}{
			operator: hspc.Memory,
		}
	}
	if hspc.CPUMhz != nil {
		filter["bk_cpu_mhz"] = map[string]interface{}{
			operator: hspc.CPUMhz,
		}
	}
	if hspc.InnerIP != nil {
		filter["bk_host_innerip"] = map[string]interface{}{
			operator: hspc.InnerIP,
		}
	}
	if hspc.HostID != nil {
		filter["bk_host_id"] = map[string]interface{}{
			operator: hspc.HostID,
		}
	}
	if hspc.OsVersion != nil {
		filter["bk_os_version"] = map[string]interface{}{
			operator: hspc.OsVersion,
		}
	}
	if hspc.CPUModule != nil {
		filter["bk_cpu_module"] = map[string]interface{}{
			operator: hspc.CPUModule,
		}
	}
	if hspc.Sn != nil {
		filter["bk_sn"] = map[string]interface{}{
			operator: hspc.Sn,
		}
	}
	if hspc.Operator != nil {
		filter["operator"] = map[string]interface{}{
			operator: hspc.Operator,
		}
	}
	if hspc.Disk != nil {
		filter["bk_disk"] = map[string]interface{}{
			operator: hspc.Disk,
		}
	}
	if hspc.Sla != nil {
		filter["bk_sla"] = map[string]interface{}{
			operator: hspc.Sla,
		}
	}
	if hspc.HostName != nil {
		filter["bk_host_name"] = map[string]interface{}{
			operator: hspc.HostName,
		}
	}
	if hspc.AssetID != nil {
		filter["bk_asset_id"] = map[string]interface{}{
			operator: hspc.AssetID,
		}
	}
	if hspc.OuterIP != nil {
		filter["bk_host_outerip"] = map[string]interface{}{
			operator: hspc.OuterIP,
		}
	}
	if hspc.ImportFrom != nil {
		filter["import_from"] = map[string]interface{}{
			operator: hspc.ImportFrom,
		}
	}
	if hspc.CPU != nil {
		filter["bk_cpu"] = map[string]interface{}{
			operator: hspc.CPU,
		}
	}
	if hspc.ProvinceName != nil {
		filter["bk_province_name"] = map[string]interface{}{
			operator: hspc.ProvinceName,
		}
	}
	if hspc.OSBit != nil {
		filter["bk_os_bit"] = map[string]interface{}{
			operator: hspc.OSBit,
		}
	}
	if hspc.OSType != nil {
		filter["bk_os_type"] = map[string]interface{}{
			operator: hspc.OSType,
		}
	}
	if hspc.Comment != nil {
		filter["bk_comment"] = map[string]interface{}{
			operator: hspc.Comment,
		}
	}
	if hspc.Mac != nil {
		filter["bk_mac"] = map[string]interface{}{
			operator: hspc.Mac,
		}
	}
	if hspc.CloudID != nil {
		filter["bk_cloud_id"] = map[string]interface{}{
			operator: hspc.CloudID,
		}
	}
	if hspc.OSName != nil {
		filter["bk_os_name"] = map[string]interface{}{
			operator: hspc.OSName,
		}
	}
	if hspc.StateName != nil {
		filter["bk_state_name"] = map[string]interface{}{
			operator: hspc.StateName,
		}
	}
	if hspc.OuterMac != nil {
		filter["bk_outer_mac"] = map[string]interface{}{
			operator: hspc.OuterMac,
		}
	}
	if hspc.ServiceTerm != nil {
		filter["bk_service_term"] = map[string]interface{}{
			operator: hspc.ServiceTerm,
		}
	}
	if hspc.BakOperator != nil {
		filter["bk_bak_operator"] = map[string]interface{}{
			operator: hspc.BakOperator,
		}
	}
	if hspc.ISPName != nil {
		filter["bk_isp_name"] = map[string]interface{}{
			operator: hspc.ISPName,
		}
	}
	if hspc.SupplierAccount != nil {
		filter["bk_supplier_account"] = map[string]interface{}{
			operator: hspc.SupplierAccount,
		}
	}

	switch operator {
	case common.BKDBIN:
		if hspc.LastTime.Start != nil {
			filter["last_time"] = map[string]interface{}{
				common.BKDBGTE: hspc.LastTime.Start,
			}
		}
		if hspc.LastTime.End != nil {
			filter["last_time"] = map[string]interface{}{
				common.BKDBLTE: hspc.LastTime.End,
			}
		}
		if hspc.CreateTime.Start != nil {
			filter["create_time"] = map[string]interface{}{
				common.BKDBGTE: hspc.CreateTime.Start,
			}
		}
		if hspc.CreateTime.End != nil {
			filter["create_time"] = map[string]interface{}{
				common.BKDBLTE: hspc.CreateTime.End,
			}
		}
	case common.BKDBNIN:
		if hspc.LastTime.Start != nil {
			filter["last_time"] = map[string]interface{}{
				common.BKDBLT: hspc.LastTime.Start,
			}
		}
		if hspc.LastTime.End != nil {
			filter["last_time"] = map[string]interface{}{
				common.BKDBGT: hspc.LastTime.End,
			}
		}
		if hspc.CreateTime.Start != nil {
			filter["create_time"] = map[string]interface{}{
				common.BKDBLT: hspc.CreateTime.Start,
			}
		}
		if hspc.CreateTime.End != nil {
			filter["create_time"] = map[string]interface{}{
				common.BKDBGT: hspc.CreateTime.End,
			}
		}
	default:
		panic(fmt.Errorf("unexpected operator: %s", operator))
	}
	return filter
}

type ListHostByTopoNodeOption struct {
	BizID              int64                    `json:"bk_biz_id,omitempty"`
	SetIDs             []int64                  `json:"bk_set_ids"`
	ModuleIDs          []int64                  `json:"bk_module_ids"`
	StdProperty        HostStdPropertyCondition `json:"std_property"`
	ExcludeStdProperty HostStdPropertyCondition `json:"exclude_std_property"`
	Page               BasePage                 `json:"page"`
}

func (option ListHostByTopoNodeOption) Validate() (string, error) {
	if option.BizID == 0 {
		return "bk_biz_id", errors.New("bk_biz_id field shouldn't be empty")
	}

	if key, err := option.Page.Validate(); err != nil {
		return fmt.Sprintf("page.%s", key), err
	}

	if key, err := option.StdProperty.Validate(common.BKDBIN); err != nil {
		return fmt.Sprintf("std_property.%s", key), err
	}

	if key, err := option.ExcludeStdProperty.Validate(common.BKDBNIN); err != nil {
		return fmt.Sprintf("std_property.%s", key), err
	}
	return "", nil
}

func (option ListHostByTopoNodeOption) GetHostPropertyFilter(ctx context.Context) (map[string]interface{}, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	includePropertyFilter := option.StdProperty.ToHostFilter(common.BKDBIN)
	excludePropertyFilter := option.ExcludeStdProperty.ToHostFilter(common.BKDBNIN)
	if err := mergo.Merge(&includePropertyFilter, excludePropertyFilter); err != nil {
		blog.Errorf("GetHostPropertyFilter failed, err: %+v, rid: %s", err, rid)
		return nil, err
	}
	return includePropertyFilter, nil
}

// ip search info
type IPInfo struct {
	Data  []string `json:"data"`
	Exact int64    `json:"exact"`
	Flag  string   `json:"flag"`
}

// search condition
type SearchCondition struct {
	Fields    []string        `json:"fields"`
	Condition []ConditionItem `json:"condition"`
	ObjectID  string          `json:"bk_obj_id"`
}

type SearchHost struct {
	Count int             `json:"count"`
	Info  []mapstr.MapStr `json:"info"`
}

type ListHostResult struct {
	Count int                      `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}

func (sh SearchHost) ExtractHostIDs() *[]int64 {
	hostIDArray := make([]int64, 0)
	for _, h := range sh.Info {
		if _, exist := h["host"]; exist == false {
			blog.ErrorJSON("unexpected error, host: %s don't have host field.", h)
			continue
		}
		hostID, exist := h["host"].(mapstr.MapStr)[common.BKHostIDField]
		if exist == false {
			blog.ErrorJSON("unexpected error, host: %s don't have host.bk_host_id field.", h)
			continue
		}
		id, err := util.GetInt64ByInterface(hostID)
		if err != nil {
			blog.ErrorJSON("unexpected error, host: %s host.bk_host_id field is not integer.", h)
			continue
		}
		hostIDArray = append(hostIDArray, id)
	}
	return &hostIDArray
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
	User            string `json:"bk_user" bson:"bk_user"`
	TaskName        string `json:"bk_task_name" bson:"bk_task_name"`
	TaskID          int64  `json:"bk_task_id" bson:"bk_task_id"`
	AccountType     string `json:"bk_account_type" bson:"bk_account_type"`
	AccountAdmin    string `json:"bk_account_admin" bson:"bk_account_admin"`
	PeriodType      string `json:"bk_period_type" bson:"bk_period_type"`
	Period          string `json:"bk_period" bson:"bk_period"`
	LastSyncTime    string `json:"bk_last_sync_time" bson:"bk_last_sync_time"`
	ObjID           string `json:"bk_obj_id" bson:"bk_obj_id"`
	Status          bool   `json:"bk_status" bson:"bk_status"`
	ResourceConfirm bool   `json:"bk_confirm" bson:"bk_confirm"`
	AttrConfirm     bool   `json:"bk_attr_confirm" bson:"bk_attr_confirm"`
	SecretID        string `json:"bk_secret_id" bson:"bk_secret_id"`
	SecretKey       string `json:"bk_secret_key" bson:"bk_secret_key"`
	OwnerID         string `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type ResourceConfirm struct {
	ObjID        string          `json:"bk_obj_id"`
	ResourceName []mapstr.MapStr `json:"bk_resource_name"`
	SourceType   string          `json:"bk_source_type"`
	SourceName   string          `json:"bk_source_name"`
	CreateTime   time.Time       `json:"create_time"`
	TaskID       string          `json:"bk_task_id"`
	ResourceID   int64           `json:"bk_resource_id"`
	ConfirmType  string          `json:"bk_confirm_type"`
	InCharge     string          `json:"bk_in_charge"`
	OwnerID      string          `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type CloudHistory struct {
	ObjID       string `json:"bk_obj_id" bson:"bk_obj_id"`
	Status      string `json:"bk_status" bson:"bk_status"`
	TimeConsume string `json:"bk_time_consume" bson:"bk_time_consume"`
	NewAdd      int    `json:"new_add" bson:"new_add"`
	AttrChanged int    `json:"attr_changed" bson:"attr_changed"`
	StartTime   string `json:"bk_start_time" bson:"bk_start_time"`
	TaskID      int64  `json:"bk_task_id" bson:"bk_task_id"`
	HistoryID   int64  `json:"bk_history_id" bson:"bk_history_id"`
	FailReason  string `json:"fail_reason" bson:"fail_reason"`
	OwnerID     string `json:"bk_supplier_account" bson:"bk_supplier_account"`
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
}

type CloudSyncRedisPendingStart struct {
	NewHeader    http.Header `json:"new_header"`
	TaskID       int64       `json:"bk_task_id"`
	TaskItemInfo TaskInfo    `json:"task_item_info"`
	OwnerID      string      `json:"bk_supplier_account"`
	Update       bool        `json:"update"`
}

type CloudSyncRedisAlreadyStarted struct {
	LastSyncTime time.Time   `json:"last_sync_time"`
	NewHeader    http.Header `json:"new_header"`
	TaskID       int64       `json:"bk_task_id"`
	TaskItemInfo TaskInfo    `json:"task_item_info"`
	OwnerID      string      `json:"bk_supplier_account"`
}

// TransferHostAcrossBusinessParameter Transfer host across business request parameter
type TransferHostAcrossBusinessParameter struct {
	SrcAppID       int64   `json:"src_bk_biz_id"`
	DstAppID       int64   `json:"dst_bk_biz_id"`
	HostID         int64   `json:"bk_host_id"`
	DstModuleIDArr []int64 `json:"bk_module_ids"`
}

// HostModuleRelationParameter get host and module  relation parameter
type HostModuleRelationParameter struct {
	AppID  int64   `json:"bk_biz_id"`
	HostID []int64 `json:"bk_host_id"`
}

// DeleteHostFromBizParameter delete host from business
type DeleteHostFromBizParameter struct {
	AppID     int64   `json:"bk_biz_id"`
	HostIDArr []int64 `json:"bk_host_ids"`
}
