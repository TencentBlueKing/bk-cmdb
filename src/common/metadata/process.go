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

type Process struct {
	ProcNum         int64     `json:"proc_num"`
	StopCmd         string    `json:"stop_cmd"`
	RestartCmd      string    `json:"restart_cmd"`
	FaceStopCmd     string    `json:"face_stop_cmd"`
	ProcessID       int       `json:"bk_process_id"`
	FuncName        string    `json:"bk_func_name"`
	SupplierAccount string    `json:"bk_supplier_account"`
	WorkPath        string    `json:"work_path"`
	BindIP          string    `json:"bind_ip"`
	Priority        int64     `json:"priority"`
	ReloadCmd       string    `json:"reload_cmd"`
	LastTime        time.Time `json:"last_time"`
	ProcessName     string    `json:"bk_process_name"`
	Port            string    `json:"port"`
	PidFile         string    `json:"pid_file"`
	AutoStart       bool      `json:"auto_start"`
	AutoTimeGap     int64     `json:"auto_time_gap"`
	CreateTime      time.Time `json:"create_time"`
	BizID           int64     `json:"bk_biz_id"`
	StartCmd        string    `json:"start_cmd"`
	FuncID          string    `json:"bk_func_id"`
	User            string    `json:"user"`
	Timeout         int64     `json:"timeout"`
	Protocol        string    `json:"protocol"`
	Description     string    `json:"description"`
}

type ServiceCategory struct {
	Metadata `json:"metadata"`

	ID   int64  `json:"id"`
	Name string `json:"name"`

	RootID   int64 `json:"root_id"`
	ParentID int64 `json:"parent_id"`
}

type ServiceCategoryWithStatistics struct {
	ServiceCategory ServiceCategory `json:"category"`
	UsageAmount     int64           `json:"usageAmount"`
}

type ServiceTemplate struct {
	Metadata

	ID int64
	// name of this service, can not be empty
	Name string

	// the class of this service, each field means a class label.
	// now, the class must have two labels.
	ServiceCategoryID int64

	Creator    string
	Modifier   string
	CreateTime time.Time
	LastTime   time.Time
}

// this works for the process instance which is used for a template.
type ProcessTemplate struct {
	ID int64
	// the service template's, which this process template belongs to.
	ServiceTemplateID int64

	// stores a process instance's data includes all the process's
	// properties's value.
	Spec Process

	// it records the relations between process instance's property and
	// whether it's used as a default value, the empty value can also be
	// a default value.
	// key is property's name,
	// value is true if this property's value is used as a default value.
	Status map[string]bool
}

// ServiceInstance is a service, which created when a host binding with a service template.
type ServiceInstance struct {
	Metadata

	ID   int64
	Name string

	// the template id can not be updated, once the service is created.
	// it can be 0 when the service is not created with a service template.
	ServiceTemplateID int64
	HostID            int64

	// the module that this service belongs to.
	ModuleID int64
}

// ServiceInstanceRelations record which service instance and process template are current process binding, process identified by ProcessID
type ServiceInstanceRelations struct {
	// unique field, 1:1 mapping with ProcessInstance.
	ProcessID int64

	ServiceInstanceID int64
	// ProcessTemplateID indicate which template are current process instantiate from.
	ProcessTemplateID int64

	// redundant field for accelerating processes by HostID
	HostID int64
}
