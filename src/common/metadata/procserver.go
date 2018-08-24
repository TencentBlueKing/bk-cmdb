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
)

type ProcModuleConfig struct {
	ApplicationID int    `json:"bk_biz_id"`
	ModuleName    string `json:"bk_module_name"`
	ProcessID     int    `json:"bk_process_id"`
}

type ProcInstanceModel struct {
	ApplicationID uint64 `json: "bk_biz_id" bson:"bk_biz_id"`
	SetID         uint64 `json: "bk_set_id" bson:"bk_set_id,omitempty"`
	SetName       string `json: "bk_set_name" bson:"bk_set_name"`
	ModuleID      uint64 `json: "bk_module_id" bson:"bk_module_id,omitempty"`
	ModuleName    string `json: "bk_module_name" bson:"bk_module_name"`
	ProcID        uint64 `json: "bk_process_id" bson:"bk_process_id"`
	FuncID        uint64 `json: "bk_func_id" bson:"bk_func_id"`
	InstanceID    uint64 `json: "bk_instance_id" bson:"bk_instance_id"`
	HostId        uint64 `json: "bk_host_id" bson:"bk_host_id"`
}

type ProcessOperate struct {
	ApplicationID string `json: "bk_biz_id"`
	SetName       string `json: "bk_set_name"`
	ModuleName    string `json: "bk_module_name"`
	FuncID        string `json: "bk_func_id"`
	InstanceID    string `json: "bk_instance_id"`
	OpType        int    `json: "bk_proc_optype"`
}

type ProcModuleResult struct {
	BaseResp `json:",inline"`
	Data     []ProcModuleConfig `json:"data"`
}

type ProcInstModelResult struct {
	BaseResp `json:",inline"`
	Data     []ProcInstanceModel `json:"data"`
}

type GseProcRespone struct {
	BaseResp `json:",inline"`
	Data     map[string]interface{} `json:"data"`
}

type GseHost struct {
	Ip           string `json:"ip,omitempty"`
	BkCloudId    int    `json:"bk_cloud_id,omitempty"`
	BkSupplierId int    `json:"bk_supplier_ed,omitempty"`
}

type GseProcMeta struct {
	Namespace string            `json:"namespace,omitempty"`
	Name      string            `json:"name,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

type ProcInfoArrResult struct {
	BaseResp `json:",inline"`
	Data     []mapstr.MapStr `json:"data"`
}
type GseProcRequest struct {
	Meta   GseProcMeta `json:"meta,omitempty"`
	Hosts  []GseHost   `json:"hosts,omitempty"`
	OpType int         `json:"op_type,omitempty"`
	Spec   GseProcSpec `json:"spec,omitempty"`
}

type GseProcSpec struct {
	Identity         GseProcIdentity         `json:"identity,omitempty"`
	Control          GseProcControl          `json:"control,omitempty"`
	Resource         GseProcResource         `json:"resource,omitempty"`
	MonitorPolicy    GseProcMonitorPlolicy   `json:"monitor_policy,omitempty"`
	WarnReportPolicy GseProcWarnReportPolicy `json:"warn_report_policy,omitempty"`
	Configmap        GseProcConfigmap        `json:"configmap,omitempty"`
}

type GseProcIdentity struct {
	IndexKey   string `json:"index_key,omitempty"`
	ProcName   string `json:"proc_name,omitempty"`
	SetupPath  string `json:"setup_path,omitempty"`
	PidPath    string `json:"pid_path,omitempty"`
	ConfigPath string `json:"config_path,omitempty"`
	LogPath    string `json:"log_path,omitempty"`
}

type GseProcControl struct {
	StartCmd   string `json:"start_cmd,omitempty"`
	StopCmd    string `json:"stop_cmd,omitempty"`
	RestartCmd string `json:"restart_cmd,omitempty"`
	ReloadCmd  string `json:"reload_cmd,omitempty"`
	KillCmd    string `json:"kill_cmd,omitempty"`
	VersionCmd string `json:"version_cmd,omitempty"`
	HealthCmd  string `json:"health_cmd,omitempty"`
}

type GseProcResource struct {
	Cpu  float64 `json:"cpu,omitempty"`
	Mem  float64 `json:"mem,omitempty"`
	Fd   int     `json:"fd,omitempty"`
	Disk int     `json:"disk,omitempty"`
	Net  int     `json:"net,omitempty"`
}

type GseProcMonitorPlolicy struct {
	AutoType       int    `json:"auto_type,omitempty"`
	StartCheckSecs int    `json:"start_check_secs,omitempty"`
	StopCheckSecs  int    `json:"stop_check_secs,omitempty"`
	StartRetries   int    `json:"start_retries,omitempty"`
	StartInterval  int    `json:"start_interval,omitempty"`
	CrontabRule    string `json:"crontab_rule,omitempty"`
}

type GseProcWarnReportPolicy struct {
	ReportId int `json:"report_id,omitempty"`
}

type GseProcConfigmap struct {
	Path string `json:"path,omitempty"`
	Md5  string `json:"md5,omitempty"`
}

type ProcessModule struct {
	AppID      int64  `json:"bk_biz_id" bson:"bk_biz_id"`
	ModuleName string `json:"bk_module_name" bson:"bk_module_name"`
	ProcessID  int64  `json:"bk_process_id" bson:"bk_process_id"`
	OwnerID    string `json:"bk_supplier_account" bson:"bk_supplier_account"`
}
