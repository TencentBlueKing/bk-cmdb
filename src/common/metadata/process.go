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
	"errors"
	"fmt"
	"time"
)

type SocketBindType string

const (
	BindLocalHost SocketBindType = "1"
	BindAll       SocketBindType = "2"
	BindInnerIP   SocketBindType = "3"
	BindOtterIP   SocketBindType = "4"
)

func (p SocketBindType) String() string {
	// TODO: how to support internationalization?
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
	ProtocolTypeTCP ProtocolType = "1"
	ProtocolTypeUDP ProtocolType = "2"
)

func (p ProtocolType) String() string {
	switch p {
	case ProtocolTypeTCP:
		return "TCP"
	case ProtocolTypeUDP:
		return "UDP"
	default:
		return ""
	}
}

type Process struct {
	ProcNum         int64          `field:"proc_num" json:"proc_num,omitempty" bson:"proc_num"`
	StopCmd         string         `field:"stop_cmd" json:"stop_cmd,omitempty" bson:"stop_cmd"`
	RestartCmd      string         `field:"restart_cmd" json:"restart_cmd,omitempty" bson:"restart_cmd"`
	ForceStopCmd    string         `field:"face_stop_cmd" json:"face_stop_cmd,omitempty" bson:"face_stop_cmd"`
	ProcessID       int64          `field:"bk_process_id" json:"bk_process_id,omitempty" bson:"bk_process_id"`
	FuncName        string         `field:"bk_func_name" json:"bk_func_name,omitempty" bson:"bk_func_name"`
	WorkPath        string         `field:"work_path" json:"work_path,omitempty" bson:"work_path"`
	BindIP          SocketBindType `field:"bind_ip" json:"bind_ip,omitempty" bson:"bind_ip"`
	Priority        int64          `field:"priority" json:"priority,omitempty" bson:"priority"`
	ReloadCmd       string         `field:"reload_cmd" json:"reload_cmd,omitempty" bson:"reload_cmd"`
	ProcessName     string         `field:"bk_process_name" json:"bk_process_name,omitempty" bson:"bk_process_name"`
	Port            string         `field:"port" json:"port,omitempty" bson:"port"`
	PidFile         string         `field:"pid_file" json:"pid_file,omitempty" bson:"pid_file"`
	AutoStart       bool           `field:"auto_start" json:"auto_start,omitempty" bson:"auto_start"`
	AutoTimeGap     int64          `field:"auto_time_gap" json:"auto_time_gap,omitempty" bson:"auto_time_gap"`
	LastTime        time.Time      `field:"last_time" json:"last_time,omitempty" bson:"last_time"`
	CreateTime      time.Time      `field:"create_time" json:"create_time,omitempty" bson:"create_time"`
	BusinessID      int64          `field:"bk_biz_id" json:"bk_biz_id,omitempty" bson:"bk_biz_id"`
	StartCmd        string         `field:"start_cmd" json:"start_cmd,omitempty" bson:"start_cmd"`
	FuncID          string         `field:"bk_func_id" json:"bk_func_id,omitempty" bson:"bk_func_id"`
	User            string         `field:"user" json:"user,omitempty" bson:"user"`
	TimeoutSeconds  int64          `field:"timeout" json:"timeout,omitempty" bson:"timeout"`
	Protocol        ProtocolType   `field:"protocol" json:"protocol,omitempty" bson:"protocol"`
	Description     string         `field:"description" json:"description,omitempty" bson:"description"`
	SupplierAccount string         `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
}

type ServiceCategory struct {
	Metadata `field:"metadata" json:"metadata" bson:"metadata"`

	ID   int64  `field:"id" json:"id,omitempty" bson:"id"`
	Name string `field:"name" json:"name,omitempty" bson:"name"`

	RootID          int64  `field:"root_id" json:"root_id,omitempty" bson:"root_id"`
	ParentID        int64  `field:"parent_id" json:"parent_id,omitempty" bson:"parent_id"`
	SupplierAccount string `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
}

func (sc *ServiceCategory) Validate() (field string, err error) {
	MaxLen := 128
	if len(sc.Name) == 0 {
		return "name", errors.New("name can't be empty")
	}

	if len(sc.Name) > MaxLen {
		return "name", fmt.Errorf("name too long, input: %d > max: %d", len(sc.Name), MaxLen)
	}
	return "", nil
}

type ServiceCategoryWithStatistics struct {
	ServiceCategory ServiceCategory `field:"category" json:"category" bson:"category"`
	UsageAmount     int64           `field:"usageAmount" json:"usageAmount" bson:"usageAmount"`
}

type ServiceTemplate struct {
	Metadata `field:"metadata" json:"metadata" bson:"metadata"`

	ID int64 `field:"id" json:"id,omitempty" bson:"id"`
	// name of this service, can not be empty
	Name string `field:"name" json:"name,omitempty" bson:"name"`

	// the class of this service, each field means a class label.
	// now, the class must have two labels.
	ServiceCategoryID int64 `field:"service_category_id" json:"service_category_id,omitempty" bson:"service_category_id"`

	Creator         string    `field:"creator" json:"creator,omitempty" bson:"creator"`
	Modifier        string    `field:"modifier" json:"modifier,omitempty" bson:"modifier"`
	CreateTime      time.Time `field:"create_time" json:"create_time,omitempty" bson:"create_time"`
	LastTime        time.Time `field:"last_time" json:"last_time,omitempty" bson:"last_time"`
	SupplierAccount string    `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
}

// this works for the process instance which is used for a template.
type ProcessTemplate struct {
	ID int64 `json:"id"`
	// the service template's, which this process template belongs to.
	ServiceTemplateID int64 `json:"serviceTemplateID"`

	// stores a process instance's data includes all the process's
	// properties's value.
	Template ProcessProperty `json:"template"`
}

type ProcessProperty struct {
	ProcNum        PropertyDetail `field:"proc_num" json:"proc_num,omitempty" bson:"proc_num,omitempty"`
	StopCmd        PropertyDetail `field:"stop_cmd" json:"stop_cmd,omitempty" bson:"stop_cmd,omitempty"`
	RestartCmd     PropertyDetail `field:"restart_cmd" json:"restart_cmd,omitempty" bson:"restart_cmd,omitempty"`
	ForceStopCmd   PropertyDetail `field:"face_stop_cmd" json:"face_stop_cmd,omitempty" bson:"face_stop_cmd,omitempty"`
	ProcessID      PropertyDetail `field:"bk_process_id" json:"bk_process_id,omitempty" bson:"bk_process_id,omitempty"`
	FuncName       PropertyDetail `field:"bk_func_name" json:"bk_func_name,omitempty" bson:"bk_func_name,omitempty"`
	WorkPath       PropertyDetail `field:"work_path" json:"work_path,omitempty" bson:"work_path,omitempty"`
	BindIP         PropertyDetail `field:"bind_ip" json:"bind_ip,omitempty" bson:"bind_ip,omitempty"`
	Priority       PropertyDetail `field:"priority" json:"priority,omitempty" bson:"priority,omitempty"`
	ReloadCmd      PropertyDetail `field:"reload_cmd" json:"reload_cmd,omitempty" bson:"reload_cmd,omitempty"`
	ProcessName    PropertyDetail `field:"bk_process_name" json:"bk_process_name,omitempty" bson:"bk_process_name,omitempty"`
	Port           PropertyDetail `field:"port" json:"port,omitempty" bson:"port,omitempty"`
	PidFile        PropertyDetail `field:"pid_file" json:"pid_file,omitempty" bson:"pid_file,omitempty"`
	AutoStart      PropertyDetail `field:"auto_start" json:"auto_start,omitempty" bson:"auto_start,omitempty"`
	AutoTimeGap    PropertyDetail `field:"auto_time_gap" json:"auto_time_gap,omitempty" bson:"auto_time_gap,omitempty"`
	StartCmd       PropertyDetail `field:"start_cmd" json:"start_cmd,omitempty" bson:"start_cmd,omitempty"`
	FuncID         PropertyDetail `field:"bk_func_id" json:"bk_func_id,omitempty" bson:"bk_func_id,omitempty"`
	User           PropertyDetail `field:"user" json:"user,omitempty" bson:"user,omitempty"`
	TimeoutSeconds PropertyDetail `field:"timeout" json:"timeout,omitempty" bson:"timeout,omitempty"`
	Protocol       PropertyDetail `field:"protocol" json:"protocol,omitempty" bson:"protocol,omitempty"`
	Description    PropertyDetail `field:"description" json:"description,omitempty" bson:"description,omitempty"`
}

type PropertyDetail struct {
	Value interface{} `field:"value" json:"value" bson:"value"`
	// it records the relations between process instance's property and
	// whether it's used as a default value, the empty value can also be a default value.
	// If the property's value is used as a default value, then this property
	// can not be changed in all the process instance's created by this process
	// template. or, it can only be changed to this default value.
	AsDefaultValue bool `field:"asDefaultValue" json:"asDefaultValue" bson:"asDefaultValue"`
}

// ServiceInstance is a service, which created when a host binding with a service template.
type ServiceInstance struct {
	Metadata `field:"metadata" json:"metadata" bson:"metadata"`
	ID       int64  `field:"id" json:"id" bson:"id"`
	Name     string `field:"name" json:"name,omitempty" bson:"name"`

	// the template id can not be updated, once the service is created.
	// it can be 0 when the service is not created with a service template.
	ServiceTemplateID int64 `field:"serviceTemplateID" json:"serviceTemplateID,omitempty" bson:"serviceTemplateID"`
	HostID            int64 `field:"hostID" json:"hostID,omitempty" bson:"hostID"`

	// the module that this service belongs to.
	ModuleID        int64  `field:"moduleID" json:"moduleID,omitempty" bson:"moduleID"`
	SupplierAccount string `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
}

// ServiceInstanceRelations record which service instance and process template are current process binding, process identified by ProcessID
type ServiceInstanceRelations struct {
	// unique field, 1:1 mapping with ProcessInstance.
	ProcessID         int64 `field:"processID" json:"processID" bson:"processID"`
	ServiceInstanceID int64 `field:"serviceInstanceID" json:"serviceInstanceID" bson:"serviceInstanceID"`

	// ProcessTemplateID indicate which template are current process instantiate from.
	ProcessTemplateID int64 `field:"processTemplateID" json:"processTemplateID" bson:"processTemplateID"`

	// redundant field for accelerating processes by HostID
	HostID          int64  `field:"hostID" json:"hostID" bson:"hostID"`
	SupplierAccount string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type ProcessInstance struct {
}
