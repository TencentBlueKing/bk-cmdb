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
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/util"
)

type DeleteCategoryInput struct {
	Metadata Metadata `json:"metadata"`
	ID       int64    `json:"id"`
}

type CreateProcessTemplateBatchInput struct {
	Metadata          Metadata        `json:"metadata"`
	ServiceTemplateID int64           `json:"serviceTemplateID"`
	Processes         []ProcessDetail `json:"processes"`
}

type DeleteProcessTemplateBatchInput struct {
	Metadata          Metadata `json:"metadata"`
	ServiceTemplateID int64    `json:"serviceTemplateID"`
	ProcessTemplates  []int64  `json:"processTemplates"`
}

type ProcessDetail struct {
	Spec *ProcessProperty `json:"spec"`
}

type ListServiceTemplateInput struct {
	Metadata Metadata `json:"metadata"`
	// this field can be empty, it a optional condition.
	ServiceCategoryID int64 `json:"serviceCategoryID,omitempty"`
}

type DeleteServiceTemplatesInput struct {
	Metadata          Metadata `json:"metadata"`
	ServiceTemplateID int64    `json:"serviceTemplateID"`
}

type CreateServiceInstanceForServiceTemplateInput struct {
	Metadata   Metadata                `json:"metadata"`
	Name       string                  `json:"name"`
	TemplateID int64                   `json:"template_id"`
	ModuleID   int64                   `json:"module_id"`
	Instances  []ServiceInstanceDetail `json:"instances"`
}

type DeleteProcessInstanceInServiceInstanceInput struct {
	Metadata           Metadata `json:"metadata"`
	ProcessInstanceIDs []int64  `json:"process_instance_ids"`
}

type GetServiceInstanceInModuleInput struct {
	Metadata Metadata `json:"metadata"`
	ModuleID int64    `json:"module_id"`
	Page     BasePage `json:"page"`
}

type FindServiceTemplateAndInstanceDifferenceOption struct {
	Metadata          Metadata `json:"metadata"`
	ModuleID          int64    `json:"module_id"`
	ServiceTemplateID int64    `json:"service_template_id"`
}

type DeleteServiceInstanceOption struct {
	Metadata          Metadata `json:"metadata"`
	ServiceInstanceID int64    `json:"id"`
}

type FindServiceAndProcessInstanceOption struct {
	Metadata          Metadata `json:"metadata"`
	ModuleID          int64    `json:"module_id"`
	ServiceTemplateID int64    `json:"service_template_id"`
}

// to describe the differences between service instance and it's service template's
// process template's attribute.
type ServiceProcessInstanceDifference struct {
	ServiceInstanceID   int64             `json:"service_instance_id"`
	ServiceInstanceName string            `json:"service_instance_name"`
	HostID              int64             `json:"host_id"`
	Differences         *DifferenceDetail `json:"differences"`
}

type DifferenceDetail struct {
	Unchanged []ProcessDifferenceDetail `json:"unchanged"`
	Changed   []ProcessDifferenceDetail `json:"changed"`
	Added     []ProcessDifferenceDetail `json:"added"`
	Removed   []ProcessDifferenceDetail `json:"removed"`
}

type ProcessDifferenceDetail struct {
	ProcessTemplateID int64                     `json:"process_template_id"`
	ProcessInstance   Process                   `json:"process_instance"`
	ChangedAttributes []ProcessChangedAttribute `json:"changed_attributes"`
}

type ProcessChangedAttribute struct {
	ID                    int64       `json:"id"`
	PropertyID            string      `json:"property_id"`
	PropertyName          string      `json:"property_name"`
	PropertyValue         interface{} `json:"property_value"`
	TemplatePropertyValue interface{} `json:"template_property_value"`
}

type ServiceInstanceDetail struct {
	HostID    int64                   `json:"host_id"`
	Processes []ProcessInstanceDetail `json:"processes"`
}

type ProcessInstanceDetail struct {
	ProcessTemplateID int64   `json:"process_template_id"`
	ProcessInfo       Process `json:"process_info"`
}

type ListProcessTemplateWithServiceTemplateInput struct {
	Metadata            Metadata `json:"metadata"`
	ProcessTemplatesIDs []int64  `json:"process_template_ids"`
	ServiceTemplateID   int64    `json:"service_template_id"`
}

type ForceSyncServiceInstanceWithTemplateInput struct {
	Metadata          Metadata `json:"metadata"`
	ServiceTemplateID int64    `json:"service_template_id"`
	ModuleID          int64    `json:"module_id"`
	ServiceInstances  []int64  `json:"service_instances"`
}

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

func (p SocketBindType) Validate() error {
	validValues := []SocketBindType{BindLocalHost, BindAll, BindInnerIP, BindOtterIP}
	if util.InArray(p, validValues) == false {
		return fmt.Errorf("invalid socket bind type, value: %s, available values: %+v", p, validValues)
	}
	return nil
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

func (p ProtocolType) Validate() error {
	validValues := []ProtocolType{ProtocolTypeTCP, ProtocolTypeUDP}
	if util.InArray(p, validValues) == false {
		return fmt.Errorf("invalid protocol type, value: %s, available values: %+v", p, validValues)
	}
	return nil
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
	Metadata Metadata `field:"metadata" json:"metadata" bson:"metadata"`

	ID   int64  `field:"id" json:"id,omitempty" bson:"id"`
	Name string `field:"name" json:"name,omitempty" bson:"name"`

	RootID          int64  `field:"root_id" json:"root_id,omitempty" bson:"root_id"`
	ParentID        int64  `field:"parent_id" json:"parent_id,omitempty" bson:"parent_id"`
	SupplierAccount string `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`

	// IsBuiltIn indicates internal system service category, which shouldn't be modified.
	IsBuiltIn bool `field:"is_built_in" json:"is_built_in" bson:"is_built_in"`
}

func (sc *ServiceCategory) Validate() (field string, err error) {
	if len(sc.Name) == 0 {
		return "name", errors.New("name can't be empty")
	}

	if len(sc.Name) > common.NameFieldMaxLength {
		return "name", fmt.Errorf("name too long, input: %d > max: %d", len(sc.Name), common.NameFieldMaxLength)
	}
	return "", nil
}

type ServiceCategoryWithStatistics struct {
	ServiceCategory ServiceCategory `field:"category" json:"category" bson:"category"`
	UsageAmount     int64           `field:"usage_amount" json:"usage_amount" bson:"usage_amount"`
}

type ServiceTemplate struct {
	Metadata Metadata `field:"metadata" json:"metadata" bson:"metadata"`

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

func (st *ServiceTemplate) Validate() (field string, err error) {
	if len(st.Name) == 0 {
		return "name", errors.New("name can't be empty")
	}

	if len(st.Name) > common.NameFieldMaxLength {
		return "name", fmt.Errorf("name too long, input: %d > max: %d", len(st.Name), common.NameFieldMaxLength)
	}
	return "", nil
}

// this works for the process instance which is used for a template.
type ProcessTemplate struct {
	ID       int64    `field:"id" json:"id,omitempty" bson:"id"`
	Metadata Metadata `field:"metadata" json:"metadata" bson:"metadata"`
	// the service template's, which this process template belongs to.
	ServiceTemplateID int64 `field:"service_template_id" json:"service_template_id" bson:"service_template_id"`

	// stores a process instance's data includes all the process's
	// properties's value.
	Property *ProcessProperty `field:"property" json:"property,omitempty" bson:"property"`

	Creator         string    `field:"creator" json:"creator,omitempty" bson:"creator"`
	Modifier        string    `field:"modifier" json:"modifier,omitempty" bson:"modifier"`
	CreateTime      time.Time `field:"create_time" json:"create_time,omitempty" bson:"create_time"`
	LastTime        time.Time `field:"last_time" json:"last_time,omitempty" bson:"last_time"`
	SupplierAccount string    `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
}

func (pt *ProcessTemplate) Validate() (field string, err error) {
	if pt.Property != nil {
		field, err = pt.Property.Validate()
		if err != nil {
			return field, err
		}
	}
	return "", nil
}

type ProcessProperty struct {
	ProcNum            PropertyInt64       `field:"proc_num" json:"proc_num,omitempty" bson:"proc_num,omitempty"`
	StopCmd            PropertyString      `field:"stop_cmd" json:"stop_cmd,omitempty" bson:"stop_cmd,omitempty"`
	RestartCmd         PropertyString      `field:"restart_cmd" json:"restart_cmd,omitempty" bson:"restart_cmd,omitempty"`
	ForceStopCmd       PropertyString      `field:"face_stop_cmd" json:"face_stop_cmd,omitempty" bson:"face_stop_cmd,omitempty"`
	FuncName           PropertyString      `field:"bk_func_name" json:"bk_func_name,omitempty" bson:"bk_func_name,omitempty"`
	WorkPath           PropertyString      `field:"work_path" json:"work_path,omitempty" bson:"work_path,omitempty"`
	BindIP             PropertyBindIP      `field:"bind_ip" json:"bind_ip,omitempty" bson:"bind_ip,omitempty"`
	Priority           PropertyInt64       `field:"priority" json:"priority,omitempty" bson:"priority,omitempty"`
	ReloadCmd          PropertyString      `field:"reload_cmd" json:"reload_cmd,omitempty" bson:"reload_cmd,omitempty"`
	ProcessName        PropertyString      `field:"bk_process_name" json:"bk_process_name,omitempty" bson:"bk_process_name,omitempty"`
	Port               PropertyPort        `field:"port" json:"port,omitempty" bson:"port,omitempty"`
	PidFile            PropertyString      `field:"pid_file" json:"pid_file,omitempty" bson:"pid_file,omitempty"`
	AutoStart          PropertyBool        `field:"auto_start" json:"auto_start,omitempty" bson:"auto_start,omitempty"`
	AutoTimeGapSeconds PropertyInt64       `field:"auto_time_gap" json:"auto_time_gap,omitempty" bson:"auto_time_gap,omitempty"`
	StartCmd           PropertyString      `field:"start_cmd" json:"start_cmd,omitempty" bson:"start_cmd,omitempty"`
	FuncID             PropertyInt64String `field:"bk_func_id" json:"bk_func_id,omitempty" bson:"bk_func_id,omitempty"`
	User               PropertyString      `field:"user" json:"user,omitempty" bson:"user,omitempty"`
	TimeoutSeconds     PropertyInt64       `field:"timeout" json:"timeout,omitempty" bson:"timeout,omitempty"`
	Protocol           PropertyProtocol    `field:"protocol" json:"protocol,omitempty" bson:"protocol,omitempty"`
	Description        PropertyString      `field:"description" json:"description,omitempty" bson:"description,omitempty"`
}

func (pt *ProcessProperty) Validate() (field string, err error) {
	// call all field's Validate method one by one
	propertyInterfaceType := reflect.TypeOf((*ProcessPropertyInterface)(nil)).Elem()
	selfVal := reflect.ValueOf(pt).Elem()
	selfType := reflect.TypeOf(pt).Elem()
	fieldCount := selfVal.NumField()
	for fieldIdx := 0; fieldIdx < fieldCount; fieldIdx++ {
		field := selfType.Field(fieldIdx)
		fieldVal := selfVal.Field(fieldIdx)

		// check implements interface
		fieldValType := fieldVal.Addr().Type()
		if !fieldValType.Implements(propertyInterfaceType) {
			msg := fmt.Sprintf("field %s of type: %s should implements %s", field.Name, fieldVal.Type().Elem().Name(), propertyInterfaceType.Name())
			panic(msg)
		}

		// call validate method by interface
		checkResult := fieldVal.Addr().MethodByName("Validate").Call([]reflect.Value{})
		out := checkResult[0]
		if !out.IsNil() {
			err := out.Interface().(error)
			tag := field.Tag.Get("json")
			fieldName := strings.Split(tag, ",")[0]
			return fieldName, err
		}
	}
	return "", nil
}

func (pt *ProcessProperty) Update(input ProcessProperty) {
	selfVal := reflect.ValueOf(pt).Elem()
	inputVal := reflect.ValueOf(input)
	fieldCount := selfVal.NumField()
	for fieldIdx := 0; fieldIdx < fieldCount; fieldIdx++ {
		inputField := inputVal.Field(fieldIdx)
		selfField := selfVal.Field(fieldIdx)
		subFieldCount := inputField.NumField()
		// subFields: Value, AsDefaultValue
		for subFieldIdx := 0; subFieldIdx < subFieldCount; subFieldIdx++ {
			inputFieldPtr := inputField.Field(subFieldIdx)
			if inputFieldPtr.IsNil() {
				continue
			}
			inputFieldValue := inputFieldPtr.Elem()

			selfFieldValuePtr := selfField.Field(subFieldIdx)
			if selfFieldValuePtr.Kind() == reflect.Ptr {
				if selfFieldValuePtr.IsNil() && selfFieldValuePtr.CanSet() {
					selfFieldValuePtr.Set(reflect.New(selfFieldValuePtr.Type().Elem()))
				}
			}

			selfFieldValue := selfFieldValuePtr.Elem()
			selfFieldValue.Set(inputFieldValue)
		}
	}
	return
}

type ProcessPropertyInterface interface {
	Validate() error
}

type PropertyInt64 struct {
	Value *int64 `field:"value" json:"value" bson:"value"`

	// AsDefaultValue records the relations between process instance's property and
	// whether it's used as a default value, the empty value can also be a default value.
	// If the property's value is used as a default value, then this property
	// can not be changed in all the process instance's created by this process
	// template. or, it can only be changed to this default value.
	AsDefaultValue *bool `field:"as_default_value" json:"as_default_value" bson:"as_default_value"`
}

func (ti *PropertyInt64) Validate() error {
	return nil
}

// PropertyInt64String is a string field that parse into int64
type PropertyInt64String struct {
	Value *string `field:"value" json:"value" bson:"value"`

	// AsDefaultValue records the relations between process instance's property and
	// whether it's used as a default value, the empty value can also be a default value.
	// If the property's value is used as a default value, then this property
	// can not be changed in all the process instance's created by this process
	// template. or, it can only be changed to this default value.
	AsDefaultValue *bool `field:"as_default_value" json:"as_default_value" bson:"as_default_value"`
}

func (ti *PropertyInt64String) Validate() error {
	if ti.Value != nil {
		if _, err := strconv.ParseInt(*ti.Value, 10, 64); err != nil {
			return err
		}
	}
	return nil
}

type PropertyBool struct {
	Value          *bool `field:"value" json:"value" bson:"value"`
	AsDefaultValue *bool `field:"as_default_value" json:"as_default_value" bson:"as_default_value"`
}

func (ti *PropertyBool) Validate() error {
	return nil
}

type PropertyString struct {
	Value          *string `field:"value" json:"value" bson:"value"`
	AsDefaultValue *bool   `field:"as_default_value" json:"as_default_value" bson:"as_default_value"`
}

func (ti *PropertyString) Validate() error {
	return nil
}

var (
	ProcessPortFormat = regexp.MustCompile(`^(((([1-9][0-9]{0,3})|([1-5][0-9]{4})|(6[0-4][0-9]{3})|(65[0-4][0-9]{2})|(655[0-2][0-9])|(6553[0-5]))-(([1-9][0-9]{0,3})|([1-5][0-9]{4})|(6[0-4][0-9]{3})|(65[0-4][0-9]{2})|(655[0-2][0-9])|(6553[0-5])))|((([1-9][0-9]{0,3})|([1-5][0-9]{4})|(6[0-4][0-9]{3})|(65[0-4][0-9]{2})|(655[0-2][0-9])|(6553[0-5]))))(,(((([1-9][0-9]{0,3})|([1-5][0-9]{4})|(6[0-4][0-9]{3})|(65[0-4][0-9]{2})|(655[0-2][0-9])|(6553[0-5])))|((([1-9][0-9]{0,3})|([1-5][0-9]{4})|(6[0-4][0-9]{3})|(65[0-4][0-9]{2})|(655[0-2][0-9])|(6553[0-5]))-(([1-9][0-9]{0,3})|([1-5][0-9]{4})|(6[0-4][0-9]{3})|(65[0-4][0-9]{2})|(655[0-2][0-9])|(6553[0-5])))))*$`)
)

type PropertyPort struct {
	Value          *string `field:"value" json:"value" bson:"value"`
	AsDefaultValue *bool   `field:"as_default_value" json:"as_default_value" bson:"as_default_value"`
}

func (ti *PropertyPort) Validate() error {
	if ti.Value != nil {
		if matched := ProcessPortFormat.MatchString(*ti.Value); matched == false {
			return fmt.Errorf("port format invalid")
		}
	}
	return nil
}

type PropertyBindIP struct {
	Value          *SocketBindType `field:"value" json:"value" bson:"value"`
	AsDefaultValue *bool           `field:"as_default_value" json:"as_default_value" bson:"as_default_value"`
}

func (ti *PropertyBindIP) Validate() error {
	if ti.Value != nil {
		if err := ti.Value.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type PropertyProtocol struct {
	Value          *ProtocolType `field:"value" json:"value" bson:"value"`
	AsDefaultValue *bool         `field:"as_default_value" json:"as_default_value" bson:"as_default_value"`
}

func (ti *PropertyProtocol) Validate() error {
	if ti.Value != nil {
		if err := ti.Value.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// ServiceInstance is a service, which created when a host binding with a service template.
type ServiceInstance struct {
	Metadata Metadata `field:"metadata" json:"metadata" bson:"metadata"`
	ID       int64    `field:"id" json:"id,omitempty" bson:"id"`
	Name     string   `field:"name" json:"name,omitempty" bson:"name"`

	// the template id can not be updated, once the service is created.
	// it can be 0 when the service is not created with a service template.
	ServiceTemplateID int64 `field:"service_template_id" json:"service_template_id,omitempty" bson:"service_template_id"`
	HostID            int64 `field:"host_id" json:"host_id,omitempty" bson:"host_id"`

	// the module that this service belongs to.
	ModuleID int64 `field:"module_id" json:"module_id,omitempty" bson:"module_id"`

	Creator         string    `field:"creator" json:"creator,omitempty" bson:"creator"`
	Modifier        string    `field:"modifier" json:"modifier,omitempty" bson:"modifier"`
	CreateTime      time.Time `field:"create_time" json:"create_time,omitempty" bson:"create_time"`
	LastTime        time.Time `field:"last_time" json:"last_time,omitempty" bson:"last_time"`
	SupplierAccount string    `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
}

func (si *ServiceInstance) Validate() (field string, err error) {
	if len(si.Name) == 0 {
		return "name", errors.New("name can't be empty")
	}

	if len(si.Name) > common.NameFieldMaxLength {
		return "name", fmt.Errorf("name too long, input: %d > max: %d", len(si.Name), common.NameFieldMaxLength)
	}
	return "", nil
}

// ServiceInstanceRelations record which service instance and process template are current process binding, process identified by ProcessID
type ProcessInstanceRelation struct {
	Metadata Metadata `field:"metadata" json:"metadata" bson:"metadata"`

	// unique field, 1:1 mapping with ProcessInstance.
	ProcessID         int64 `field:"process_id" json:"process_id" bson:"process_id"`
	ServiceInstanceID int64 `field:"service_instance_id" json:"service_instance_id" bson:"service_instance_id"`

	// ProcessTemplateID indicate which template are current process instantiate from.
	ProcessTemplateID int64 `field:"process_template_id" json:"process_template_id" bson:"process_template_id"`

	// redundant field for accelerating processes by HostID
	HostID          int64  `field:"host_id" json:"host_id" bson:"host_id"`
	SupplierAccount string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
}

func (pir *ProcessInstanceRelation) Validate() (field string, err error) {
	return "", nil
}
