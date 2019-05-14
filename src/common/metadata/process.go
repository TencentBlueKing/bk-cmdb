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

import "time"

type SocketBindType string

const (
	BindLocalHost SocketBindType = "1"
	BindAll       SocketBindType = "2"
	BindInnerIP   SocketBindType = "3"
	BindOtterIP   SocketBindType = "4"
)

func (p SocketBindType) String() string {
	switch p {
	case BindLocalHost:
		return "127.0.0.1"
	case BindAll:
		return "0.0.0.0"
	case BindInnerIP:
		return "第一内网IP"
	case BindOtterIP:
		return "第一外网IP"
	default:
		return ""
	}
}

type ProtocolType string

const (
	ProtocolType_TCP ProtocolType = "1"
	ProtocolType_UDP ProtocolType = "2"
)

func (p ProtocolType) String() string {
	switch p {
	case ProtocolType_TCP:
		return "TCP"
	case ProtocolType_UDP:
		return "UDP"
	default:
		return ""
	}
}

type Process struct {
	ProcNum         int64          `field:"proc_num" json:"proc_num" bson:"proc_num"`
	StopCmd         string         `field:"stop_cmd" json:"stop_cmd" bson:"stop_cmd"`
	RestartCmd      string         `field:"restart_cmd" json:"restart_cmd" bson:"restart_cmd"`
	ForceStopCmd    string         `field:"face_stop_cmd" json:"face_stop_cmd" bson:"face_stop_cmd"`
	ProcessID       int64          `field:"bk_process_id" json:"bk_process_id" bson:"bk_process_id"`
	FuncName        string         `field:"bk_func_name" json:"bk_func_name" bson:"bk_func_name"`
	WorkPath        string         `field:"work_path" json:"work_path" bson:"work_path"`
	BindIP          SocketBindType `field:"bind_ip" json:"bind_ip" bson:"bind_ip"`
	Priority        int64          `field:"priority" json:"priority" bson:"priority"`
	ReloadCmd       string         `field:"reload_cmd" json:"reload_cmd" bson:"reload_cmd"`
	ProcessName     string         `field:"bk_process_name" json:"bk_process_name" bson:"bk_process_name"`
	Port            string         `field:"port" json:"port" bson:"port"`
	PidFile         string         `field:"pid_file" json:"pid_file" bson:"pid_file"`
	AutoStart       bool           `field:"auto_start" json:"auto_start" bson:"auto_start"`
	AutoTimeGap     int64          `field:"auto_time_gap" json:"auto_time_gap" bson:"auto_time_gap"`
	LastTime        time.Time      `field:"last_time" json:"last_time" bson:"last_time"`
	CreateTime      time.Time      `field:"create_time" json:"create_time" bson:"create_time"`
	BusinessID      int64          `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	StartCmd        string         `field:"start_cmd" json:"start_cmd" bson:"start_cmd"`
	FuncID          string         `field:"bk_func_id" json:"bk_func_id" bson:"bk_func_id"`
	User            string         `field:"user" json:"user" bson:"user"`
	TimeoutSeconds  int64          `field:"timeout" json:"timeout" bson:"timeout"`
	Protocol        ProtocolType   `field:"protocol" json:"protocol" bson:"protocol"`
	Description     string         `field:"description" json:"description" bson:"description"`
	SupplierAccount string         `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type ServiceCategory struct {
	Metadata `field:"metadata" json:"metadata" bson:"metadata"`

	ID   int64  `field:"id" json:"id" bson:"id"`
	Name string `field:"name" json:"name" bson:"name"`

	RootID          int64  `field:"root_id" json:"root_id" bson:"root_id"`
	ParentID        int64  `field:"parent_id" json:"parent_id" bson:"parent_id"`
	SupplierAccount string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type ServiceCategoryWithStatistics struct {
	ServiceCategory ServiceCategory `field:"category" json:"category" bson:"category"`
	UsageAmount     int64           `field:"usageAmount" json:"usageAmount" bson:"usageAmount"`
}

type ServiceTemplate struct {
	Metadata `field:"metadata" json:"metadata" bson:"metadata"`

	ID int64 `field:"id" json:"id" bson:"id"`
	// name of this service, can not be empty
	Name string `field:"name" json:"name" bson:"name"`

	// the class of this service, each field means a class label.
	// now, the class must have two labels.
	ServiceCategoryID int64 `field:"service_category_id" json:"service_category_id" bson:"service_category_id"`

	Creator         string    `field:"creator" json:"creator" bson:"creator"`
	Modifier        string    `field:"modifier" json:"modifier" bson:"modifier"`
	CreateTime      time.Time `field:"create_time" json:"create_time" bson:"create_time"`
	LastTime        time.Time `field:"last_time" json:"last_time" bson:"last_time"`
	SupplierAccount string    `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// this works for the process instance which is used for a template.
type ProcessTemplate struct {
	ID int64 `field:"id" json:"id" bson:"id"`
	// the service template's, which this process template belongs to.
	ServiceTemplateID int64 `field:"service_template_id" json:"service_template_id" bson:"service_template_id"`

	// stores a process instance's data includes all the process's
	// properties's value.
	Spec Process `field:"spec" json:"spec" bson:"spec"`

	// it records the relations between process instance's property and
	// whether it's used as a default value, the empty value can also be
	// a default value.
	// key is property's name,
	// value is true if this property's value is used as a default value.
	Status          map[string]bool `field:"status" json:"status" bson:"status"`
	SupplierAccount string          `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// ServiceInstance is a service, which created when a host binding with a service template.
type ServiceInstance struct {
	Metadata `field:"metadata" json:"metadata" bson:"metadata"`
	ID       int64  `field:"id" json:"id" bson:"id"`
	Name     string `field:"name" json:"name" bson:"name"`

	// the template id can not be updated, once the service is created.
	// it can be 0 when the service is not created with a service template.
	ServiceTemplateID int64 `field:"service_template_id" json:"service_template_id" bson:"service_template_id"`
	HostID            int64 `field:"host_id" json:"host_id" bson:"host_id"`

	// the module that this service belongs to.
	ModuleID        int64  `field:"module_id" json:"module_id" bson:"module_id"`
	SupplierAccount string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// ServiceInstanceRelations record which service instance and process template are current process binding, process identified by ProcessID
type ServiceInstanceRelations struct {
	// unique field, 1:1 mapping with ProcessInstance.
	ProcessID         int64 `field:"process_id" json:"process_id" bson:"process_id"`
	ServiceInstanceID int64 `field:"service_instance_id" json:"service_instance_id" bson:"service_instance_id"`

	// ProcessTemplateID indicate which template are current process instantiate from.
	ProcessTemplateID int64 `field:"process_template_id" json:"process_template_id" bson:"process_template_id"`

	// redundant field for accelerating processes by HostID
	HostID          int64  `field:"host_id" json:"host_id" bson:"host_id"`
	SupplierAccount string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
}
