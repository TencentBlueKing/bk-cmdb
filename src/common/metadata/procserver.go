// Package metadata TODO
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
	"net/http"
	"time"

	"configcenter/src/common/mapstr"
)

// ProcModuleConfig TODO
type ProcModuleConfig struct {
	ApplicationID int64  `json:"bk_biz_id"`
	ModuleName    string `json:"bk_module_name"`
	ProcessID     int64  `json:"bk_process_id"`
}

// ProcInstanceModel TODO
type ProcInstanceModel struct {
	ApplicationID  int64  `json:"bk_biz_id" bson:"bk_biz_id"`
	SetID          int64  `json:"bk_set_id" bson:"bk_set_id,omitempty"`
	ModuleID       int64  `json:"bk_module_id" bson:"bk_module_id,omitempty"`
	ProcID         int64  `json:"bk_process_id" bson:"bk_process_id"`
	ProcInstanceID uint64 `json:"proc_instance_id" bson:"proc_instance_id"`
	HostID         int64  `json:"bk_host_id" bson:"bk_host_id"`
	HostInstanID   uint64 `json:"bk_host_instance_id" bson:"bk_host_instance_id"`
	HostProcID     uint64 `json:"host_proc_id" bson:"host_proc_id"`
	OwnerID        string `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// MatchProcInstParam TODO
type MatchProcInstParam struct {
	ApplicationID  int64  `json:"bk_biz_id" bson:"bk_biz_id"`
	SetName        string `json:"bk_set_name" bson:"bk_set_name"`
	ModuleName     string `json:"bk_module_name" bson:"bk_module_name"`
	HostInstanceID string `json:"bk_host_instance_id" bson:"bk_host_instance_id"`
}

// ProcessOperate TODO
type ProcessOperate struct {
	MatchProcInstParam `json:",inline"`
	OpType             int `json:"bk_proc_optype"`
}

// ProcModuleResult TODO
type ProcModuleResult struct {
	BaseResp `json:",inline"`
	Data     []ProcModuleConfig `json:"data"`
}

// ProcInstModelResult TODO
type ProcInstModelResult struct {
	BaseResp `json:",inline"`
	Data     struct {
		Count int                 `json:"count"`
		Info  []ProcInstanceModel `json:"info"`
	} `json:"data"`
}

// GseHost TODO
type GseHost struct {
	HostID    int64  `json:"bk_host_id,omitempty"`
	Ip        string `json:"ip,omitempty"`
	BkCloudId int64  `json:"bk_cloud_id"`
}

// GseProcMeta TODO
type GseProcMeta struct {
	Namespace string            `json:"namespace,omitempty"`
	Name      string            `json:"name,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// ProcInfoArrResult TODO
type ProcInfoArrResult struct {
	BaseResp `json:",inline"`
	Data     []mapstr.MapStr `json:"data"`
}

// GseProcRequest TODO
type GseProcRequest struct {
	AppID    int64       `json:"bk_biz_id"  bson:"bk_biz_id"`
	ModuleID int64       `json:"bk_module_id" bson:"bk_module_id"`
	ProcID   int64       `json:"bk_process_id" bson:"bk_process_id"`
	Meta     GseProcMeta `json:"meta,omitempty" bson:"meta"`
	Hosts    []GseHost   `json:"hosts,omitempty" bson:"hosts"`
	OpType   int         `json:"op_type,omitempty" bson:"-"`
	Spec     GseProcSpec `json:"spec,omitempty" bson:"spec"`
}

// ProcInstanceDetail TODO
type ProcInstanceDetail struct {
	GseProcRequest `json:",inline" bson:",inline"`
	OwnerID        string                   `json:"bk_supplier_account" bson:"bk_supplier_account"`
	HostID         int64                    `json:"bk_host_id" bson:"bk_host_id"`
	Status         ProcInstanceDetailStatus `json:"status" bson:"status"` // 1 register gse sucess, 2 register error need retry 3 unregister error need retry
}

// ProcInstanceDetailResult TODO
type ProcInstanceDetailResult struct {
	BaseResp `json:",inline"`

	Data struct {
		Count int                  `json:"count"`
		Info  []ProcInstanceDetail `json:"info"`
	} `json:"data"`
}

// ProcInstanceDetailStatus TODO
type ProcInstanceDetailStatus int64

var (
	// ProcInstanceDetailStatusRegisterSucc TODO
	ProcInstanceDetailStatusRegisterSucc ProcInstanceDetailStatus = 1
	// ProcInstanceDetailStatusRegisterFailed TODO
	ProcInstanceDetailStatusRegisterFailed ProcInstanceDetailStatus = 2
	// ProcInstanceDetailStatusUnRegisterFailed TODO
	ProcInstanceDetailStatusUnRegisterFailed ProcInstanceDetailStatus = 10
)

// ModifyProcInstanceDetail TODO
type ModifyProcInstanceDetail struct {
	Conds map[string]interface{} `json:"condition"`
	Data  map[string]interface{} `json:"data"`
}

// GseProcSpec TODO
type GseProcSpec struct {
	Identity         GseProcIdentity         `json:"identity,omitempty"`
	Control          GseProcControl          `json:"control,omitempty"`
	Resource         GseProcResource         `json:"resource,omitempty"`
	MonitorPolicy    GseProcMonitorPlolicy   `json:"monitor_policy,omitempty"`
	WarnReportPolicy GseProcWarnReportPolicy `json:"warn_report_policy,omitempty"`
	Configmap        []GseProcConfigmap      `json:"configmap,omitempty"`
}

// GseProcIdentity TODO
type GseProcIdentity struct {
	IndexKey   string `json:"index_key,omitempty"`
	ProcName   string `json:"proc_name,omitempty"`
	SetupPath  string `json:"setup_path,omitempty"`
	PidPath    string `json:"pid_path,omitempty"`
	ConfigPath string `json:"config_path,omitempty"`
	LogPath    string `json:"log_path,omitempty"`
}

// GseProcControl TODO
type GseProcControl struct {
	StartCmd   string `json:"start_cmd,omitempty"`
	StopCmd    string `json:"stop_cmd,omitempty"`
	RestartCmd string `json:"restart_cmd,omitempty"`
	ReloadCmd  string `json:"reload_cmd,omitempty"`
	KillCmd    string `json:"kill_cmd,omitempty"`
	VersionCmd string `json:"version_cmd,omitempty"`
	HealthCmd  string `json:"health_cmd,omitempty"`
}

// GseProcResource TODO
type GseProcResource struct {
	Cpu  float64 `json:"cpu,omitempty"`
	Mem  float64 `json:"mem,omitempty"`
	Fd   int     `json:"fd,omitempty"`
	Disk int     `json:"disk,omitempty"`
	Net  int     `json:"net,omitempty"`
}

// GseProcMonitorPlolicy TODO
type GseProcMonitorPlolicy struct {
	AutoType       int    `json:"auto_type,omitempty"`
	StartCheckSecs int    `json:"start_check_secs,omitempty"`
	StopCheckSecs  int    `json:"stop_check_secs,omitempty"`
	StartRetries   int    `json:"start_retries,omitempty"`
	StartInterval  int    `json:"start_interval,omitempty"`
	CrontabRule    string `json:"crontab_rule,omitempty"`
}

// GseProcWarnReportPolicy TODO
type GseProcWarnReportPolicy struct {
	ReportId int `json:"report_id,omitempty"`
}

// GseProcConfigmap TODO
type GseProcConfigmap struct {
	Path string `json:"path,omitempty"`
	Md5  string `json:"md5,omitempty"`
}

// FilePriviewMap TODO
type FilePriviewMap struct {
	Content string `json:"content"`
	Inst    string `json:"inst"`
}

// InlineProcInfo process info convert gse proc info
type InlineProcInfo struct {
	// Meta    GseProcMeta
	// Spec    GseProcSpec
	ProcInfo map[string]interface{}
	ProcNum  int64
	AppID    int64 // use gse proc namespace
	FunID    int64
	ProcID   int64
}

// ProcessOperateTask TODO
type ProcessOperateTask struct {
	OperateInfo *ProcessOperate                     `json:"operate_info" bson:"operate_info"`
	TaskID      string                              `json:"task_id" bson:"task_id"`
	GseTaskID   string                              `json:"gse_task_id" bson:"gse_task_id"`
	Namespace   string                              `json:"namespace" bson:"namespace"`
	Status      ProcOpTaskStatus                    `json:"status" bson:"status"`
	CreateTime  time.Time                           `json:"create_time" bson:"create_time"`
	OwnerID     string                              `json:"bk_supplier_account" bson:"bk_supplier_account"`
	User        string                              `json:"user,omitempty" bson:"user,omitempty"`
	Detail      map[string]ProcessOperateTaskDetail `json:"detail" bson:"detail"`
	Host        []GseHost                           `json:"host_info" bson:"host_info"`
	ProcName    string                              `json:"bk_process_name" bson:"bk_process_name"`
	HTTPHeader  http.Header                         `json:"http_header" bson:"http_header"`
}

// ProcOpTaskStatus TODO
type ProcOpTaskStatus int64

var (
	// ProcOpTaskStatusWaitOP TODO
	ProcOpTaskStatusWaitOP ProcOpTaskStatus = 1
	// ProcOpTaskStatusExecuteing TODO
	ProcOpTaskStatusExecuteing ProcOpTaskStatus = 115
	// ProcOpTaskStatusErr TODO
	ProcOpTaskStatusErr ProcOpTaskStatus = 2
	// ProcOpTaskStatusSucc TODO
	ProcOpTaskStatusSucc ProcOpTaskStatus = 3
	// ProcOpTaskStatusHTTPErr TODO
	ProcOpTaskStatusHTTPErr ProcOpTaskStatus = 1101
	// ProcOpTaskStatusNotTaskIDErr TODO
	ProcOpTaskStatusNotTaskIDErr ProcOpTaskStatus = 1112
)

// ProcessOperateTaskResult TODO
type ProcessOperateTaskResult struct {
	BaseResp `json:",inline"`
	Data     struct {
		Count int                  `json:"count"`
		Info  []ProcessOperateTask `json:"info"`
	} `json:"data"`
}

// ProcessOperateTaskDetail TODO
type ProcessOperateTaskDetail struct {
	Errcode int    `json:"errcode" bson:"error_code"`
	ErrMsg  string `json:"errmsg" bson:"error_msg"`
}

// GseProcessOperateTaskResult TODO
type GseProcessOperateTaskResult struct {
	Data            map[string]ProcessOperateTaskDetail `json:"data"`
	EsbBaseResponse `json:",inline"`
}

// EsbResponse TODO
type EsbResponse struct {
	EsbBaseResponse `json:",inline"`
	Data            mapstr.MapStr `json:"data"`
}

// UserInfo TODO
type UserInfo struct {
	Qq          string `json:"qq"`
	Status      string `json:"status"`
	WxUserid    string `json:"wx_userid"`
	DisplayName string `json:"display_name"`
	Language    string `json:"language"`
	StaffStatus string `json:"staff_status"`
	BkUsername  string `json:"bk_username"`
	Telephone   string `json:"telephone"`
	BkRole      int    `json:"bk_role"`
	TimeZone    string `json:"time_zone"`
	Email       string `json:"email"`
}

// EsbUserListResponse TODO
type EsbUserListResponse struct {
	EsbBaseResponse `json:",inline"`
	Data            []UserInfo `json:"data"`
}

// EsbListUserResponse TODO
type EsbListUserResponse struct {
	EsbBaseResponse `json:",inline"`
	Data            []ListUserItem `json:"data"`
}

// ListUserItem TODO
type ListUserItem struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

// EsbDepartmentResponse TODO
type EsbDepartmentResponse struct {
	EsbBaseResponse `json:",inline"`
	Data            DepartmentData `json:"data"`
}

// DepartmentData TODO
type DepartmentData struct {
	Count   int64            `json:"count"`
	Results []DepartmentItem `json:"results"`
}

// DepartmentItem TODO
type DepartmentItem struct {
	ID          int64      `json:"id"`
	Parent      int64      `json:"parent"`
	Name        string     `json:"name"`
	FullName    string     `json:"full_name"`
	Level       int        `json:"level"`
	HasChildren bool       `json:"has_children"`
	Ancestors   []Ancestor `json:"ancestors"`
}

// Ancestor TODO
type Ancestor struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// EsbDepartmentProfileResponse TODO
type EsbDepartmentProfileResponse struct {
	EsbBaseResponse `json:",inline"`
	Data            DepartmentProfileData `json:"data"`
}

// DepartmentProfileData TODO
type DepartmentProfileData struct {
	Count   int64                   `json:"count"`
	Results []DepartmentProfileItem `json:"results"`
}

// DepartmentProfileItem TODO
type DepartmentProfileItem struct {
	ID   int64  `json:"id"`
	Name string `json:"username"`
}

// EsbBaseResponse TODO
type EsbBaseResponse struct {
	Result       bool   `json:"result"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
	EsbRequestID string `json:"request_id"`
}

// ProcessModule TODO
type ProcessModule struct {
	AppID      int64  `json:"bk_biz_id" bson:"bk_biz_id"`
	ModuleName string `json:"bk_module_name" bson:"bk_module_name"`
	ProcessID  int64  `json:"bk_process_id" bson:"bk_process_id"`
	OwnerID    string `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// TemplateVersion TODO
type TemplateVersion struct {
	Content     string `json:"content" field:"content"`
	Status      string `json:"status" field:"status"`
	Description string `json:"description" field:"description"`
}

// ListProcessRelatedInfoResponse TODO
type ListProcessRelatedInfoResponse struct {
	BaseResp `json:",inline"`
	Data     struct {
		Count int                            `json:"count"`
		Info  []ListProcessRelatedInfoResult `json:"info"`
	} `json:"data"`
}
