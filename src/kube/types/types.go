/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package types

import (
	"errors"
	"time"
)

// identification of k8s in cc
const (
	// KubeCluster k8s cluster type
	KubeCluster = "cluster"

	// KubeNode k8s node type
	KubeNode = "node"

	// KubeNamespace k8s namespace type
	KubeNamespace = "namespace"

	// KubeWorkload k8s workload type
	KubeWorkload = "workload"

	// KubePod k8s pod type
	KubePod = "pod"

	// KubeContainer k8s container type
	KubeContainer = "container"
)

// WorkloadType workload type enum
type WorkloadType string

const (
	// KubeDeployment k8s deployment type
	KubeDeployment WorkloadType = "deployment"

	// KubeStatefulSet k8s statefulSet type
	KubeStatefulSet WorkloadType = "statefulSet"

	// KubeDaemonSet k8s daemonSet type
	KubeDaemonSet WorkloadType = "daemonSet"

	// KubeGameStatefulSet k8s gameStatefulSet type
	KubeGameStatefulSet WorkloadType = "gameStatefulSet"

	// KubeGameDeployment k8s gameDeployment type
	KubeGameDeployment WorkloadType = "gameDeployment"

	// KubeCronJob k8s cronJob type
	KubeCronJob WorkloadType = "cronJob"

	// KubeJob k8s job type
	KubeJob WorkloadType = "job"

	// KubePodWorkload k8s pod workload type
	KubePodWorkload WorkloadType = "pods"
)

// table names
const (
	// BKTableNameBaseCluster the table name of the Cluster
	BKTableNameBaseCluster = "cc_ClusterBase"

	// BKTableNameBaseNode the table name of the Node
	BKTableNameBaseNode = "cc_NodeBase"

	// BKTableNameBaseNamespace the table name of the Namespace
	BKTableNameBaseNamespace = "cc_NamespaceBase"

	// BKTableNameBaseDeployment the table name of the Deployment
	BKTableNameBaseDeployment = "cc_DeploymentBase"

	// BKTableNameBaseStatefulSet the table name of the StatefulSet
	BKTableNameBaseStatefulSet = "cc_StatefulSetBase"

	// BKTableNameBaseDaemonSet the table name of the DaemonSet
	BKTableNameBaseDaemonSet = "cc_DaemonSetBase"

	// BKTableNameGameDeployment the table name of the GameDeployment
	BKTableNameGameDeployment = "cc_GameDeploymentBase"

	// BKTableNameGameStatefulSet the table name of the GameStatefulSet
	BKTableNameGameStatefulSet = "cc_GameStatefulSetBase"

	// BKTableNameBaseCronJob the table name of the CronJob
	BKTableNameBaseCronJob = "cc_CronJobBase"

	// BKTableNameBaseJob the table name of the Job
	BKTableNameBaseJob = "cc_JobBase"

	// BKTableNameBasePodWorkload the table name of the Pod Workload
	BKTableNameBasePodWorkload = "cc_PodWorkloadBase"

	// BKTableNameBaseCustom the table name of the Custom Workload
	BKTableNameBaseCustom = "cc_CustomBase"

	// BKTableNameBasePod the table name of the Pod
	BKTableNameBasePod = "cc_PodBase"

	// BKTableNameBaseContainer the table name of the Container
	BKTableNameBaseContainer = "cc_ContainerBase"
)

// common field names
const (
	// UidField unique id field in third party platform
	UidField = "uid"

	// LabelsField object labels field
	LabelsField = "labels"

	// KindField object kind field
	KindField = "kind"
)

// cluster field names
const (

	// BKIDField the id definition
	BKIDField = "id"

	// KubeNameField the name definition
	KubeNameField = "name"

	// BKBizIDField business id field
	BKBizIDField = "bk_biz_id"

	// BKSupplierAccountField supplier account
	BKSupplierAccountField = "bk_supplier_account"

	// CreatorField the creator field
	CreatorField = "creator"

	// ModifierField the modifier field
	ModifierField = "modifier"

	// CreateTimeField the create time field
	CreateTimeField = "create_time"

	// LastTimeField the last time field
	LastTimeField = "last_time"

	// BKClusterIDFiled cluster unique id field in cc
	BKClusterIDFiled = "bk_cluster_id"

	// ClusterUIDField cluster unique id field in third party platform
	ClusterUIDField = "cluster_uid"

	// XidField base cluster id field
	XidField = "xid"

	// VersionField cluster version field
	VersionField = "version"

	// NetworkTypeField cluster network type field
	NetworkTypeField = "network_type"

	// RegionField cluster region field
	RegionField = "region"

	// VpcField cluster vpc field
	VpcField = "vpc"

	// NetworkField cluster network field
	NetworkField = "network"

	// TypeField cluster type field
	TypeField = "type"

	// SchedulingEngineField scheduling engine
	SchedulingEngineField = "scheduling_engine"
)

// node field names
const (
	// RolesField node role field
	RolesField = "roles"

	// TaintsField node taints field
	TaintsField = "taints"

	// UnschedulableField node unschedulable field
	UnschedulableField = "unschedulable"

	// InternalIPField node internal ip field
	InternalIPField = "internal_ip"

	// ExternalIPField node external ip field
	ExternalIPField = "external_ip"

	// HostnameField node hostname field
	HostnameField = "hostname"

	// RuntimeComponentField node runtime component field
	RuntimeComponentField = "runtime_component"

	// KubeProxyModeField node proxy mode field
	KubeProxyModeField = "kube_proxy_mode"

	// PodCidrField pod address allocation range
	PodCidrField = "pod_cidr"

	// BKNodeIDField cluster unique id field in cc
	BKNodeIDField = "bk_node_id"

	// NodeField node name field in third party platform
	NodeField = "node"
)

// namespace field names
const (
	// ResourceQuotasField namespace resource quotas field
	ResourceQuotasField = "resource_quotas"

	// BKNamespaceIDField namespace unique id field in cc
	BKNamespaceIDField = "bk_namespace_id"

	// NamespaceField namespace name field in third party platform
	NamespaceField = "namespace"
)

// workload fields names
const (
	// SelectorField workload selector field
	SelectorField = "selector"

	// ReplicasField workload replicas field
	ReplicasField = "replicas"

	// StrategyTypeField workload strategy type field
	StrategyTypeField = "strategy_type"

	// MinReadySecondsField workload minimum ready seconds field
	MinReadySecondsField = "min_ready_seconds"

	// RollingUpdateStrategyField workload rolling update strategy field
	RollingUpdateStrategyField = "rolling_update_strategy"
)

// pod field names
const (
	// PriorityField pod priority field
	PriorityField = "priority"

	// IPField pod ip field
	IPField = "ip"

	// IPsField pod ips field
	IPsField = "ips"

	// ControlledBy owning replica controller
	ControlledBy = "controlled_by"

	// VolumesField pod volumes field
	VolumesField = "volumes"

	// QOSClassField pod qos class field
	QOSClassField = "qos_class"

	// NodeSelectorsField pod node selectors field
	NodeSelectorsField = "node_selectors"

	// TolerationsField pod tolerations field
	TolerationsField = "tolerations"

	// BKPodIDField pod unique id field in cc
	BKPodIDField = "bk_pod_id"

	// PodUIDField pod unique id field in third party platform
	PodUIDField = "pod_uid"
)

// container field names
const (
	// ContainerUIDField container unique id field in third party platform
	ContainerUIDField = "container_uid"

	// BKContainerIDField container unique id field in cc
	BKContainerIDField = "bk_container_id"

	// ImageField container image field
	ImageField = "image"

	// PortsField container ports field
	PortsField = "ports"

	// HostPortsField container host ports field
	HostPortsField = "host_ports"

	// ArgsField container args field
	ArgsField = "args"

	// StartedField container started Field
	StartedField = "started"

	// LimitsField container limits field
	LimitsField = "limits"

	// RequestsField container requests field
	RequestsField = "requests"

	// LivenessField container liveness field
	LivenessField = "liveness"

	// EnvironmentField container environment field
	EnvironmentField = "environment"

	// MountsField container mounts field
	MountsField = "mounts"
)

// FieldType define the table's field data type.
type FieldType string

const (
	// Numeric means this field is numeric data type.
	Numeric FieldType = "numeric"
	// Boolean means this field is boolean data type.
	Boolean FieldType = "bool"
	// String means this field is string data type.
	String FieldType = "string"
	// MapString means this field is map string type.
	MapString FieldType = "mapString"
	// Array means this field is array data type.
	Array FieldType = "array"
	// Object means this field is object data type.
	Object FieldType = "object"
	// Enum means this field is object enum type.
	Enum FieldType = "enum"
	// Note: subsequent support for other types can be added here.
	// after adding a type, pay attention to adding a verification
	// function for this type synchronously. special attention is
	// paid to whether the array elements also need to synchronize support for this type.
)

// Fields table's fields details.
type Fields struct {
	// descriptors specific description of the field.
	descriptors []FieldDescriptor
	// fields defines all the table's fields.
	fields []string
	// fieldType the type corresponding to the field.
	fieldType map[string]FieldType
}

// FieldsDescriptors table of field descriptor.
type FieldsDescriptors []FieldDescriptor

func mergeFields(all ...FieldsDescriptors) *Fields {
	result := &Fields{
		descriptors: make([]FieldDescriptor, 0),
		fields:      make([]string, 0),
		fieldType:   make(map[string]FieldType),
	}

	if len(all) == 0 {
		return result
	}

	for _, col := range all {
		for _, f := range col {
			result.descriptors = append(result.descriptors, f)
			result.fieldType[f.Field] = f.Type
			result.fields = append(result.fields, f.Field)
		}
	}
	return result
}

// FieldsType returns the corresponding type of all fields.
func (f Fields) FieldsType() map[string]FieldType {
	copied := make(map[string]FieldType)
	for k, v := range f.fieldType {
		copied[k] = v
	}

	return copied
}

// OneFieldType returns the type corresponding to the specified field.
func (f Fields) OneFieldType(field string) FieldType {

	var fieldType FieldType
	if field == "" {
		return fieldType
	}

	for k, v := range f.fieldType {
		if k == field {
			fieldType = v
			break
		}
	}
	return fieldType
}

// FieldsDescriptor returns table's all fields descriptor.
func (f Fields) FieldsDescriptor() []FieldDescriptor {
	return f.descriptors
}

// OneFieldDescriptor returns one field's descriptor.
func (f Fields) OneFieldDescriptor(field string) FieldDescriptor {
	if field == "" {
		return FieldDescriptor{}
	}

	for idx := range f.descriptors {
		if f.descriptors[idx].Field == field {
			return f.descriptors[idx]
		}
	}
	return FieldDescriptor{}
}

// Fields returns all the table's fields.
func (f Fields) Fields() []string {
	copied := make([]string, len(f.fields))
	for idx := range f.fields {
		copied[idx] = f.fields[idx]
	}
	return copied
}

// mergeFieldDescriptors merge all fields of a table together.
func mergeFieldDescriptors(resources ...FieldsDescriptors) FieldsDescriptors {
	if len(resources) == 0 {
		return make([]FieldDescriptor, 0)
	}

	merged := make([]FieldDescriptor, 0)
	for _, one := range resources {
		merged = append(merged, one...)
	}

	return merged
}

// FieldDescriptor defines a table's field related information.
type FieldDescriptor struct {
	// Field is field's name.
	Field string
	// Type is this field's data type.
	Type FieldType
	// Required is it required.
	Required bool
	// IsEditable is it editable.
	IsEditable bool
	// Option additional information for the field.
	// the content corresponding to different fields may be different.
	Option interface{}
	_      struct{}
}

// Revision resource version information.
type Revision struct {
	Creator    string `json:"creator" bson:"creator"`
	Modifier   string `json:"modifier" bson:"modifier"`
	CreateTime int64  `json:"create_time" bson:"create_time"`
	LastTime   int64  `json:"last_time" bson:"last_time"`
}

// IsCreateEmpty insert data case validator and creator.
func (r Revision) IsCreateEmpty() bool {
	if len(r.Creator) != 0 {
		return false
	}

	if r.CreateTime == 0 {
		return false
	}

	return true
}

// lagSeconds fault tolerance for ntp errors of different devices.
const lagSeconds = 5 * 60

// ValidateCreate insert data case validator and creator.
func (r Revision) ValidateCreate() error {

	if len(r.Creator) == 0 {
		return errors.New("creator can not be empty")
	}

	now := time.Now().Unix()
	if (r.CreateTime <= (now - lagSeconds)) || (r.CreateTime >= (now + lagSeconds)) {
		return errors.New("invalid create time")
	}

	return nil
}

// IsModifyEmpty the update data scene verifies the revisioner and modification time of the updated data.
func (r Revision) IsModifyEmpty() bool {
	if len(r.Modifier) != 0 {
		return false
	}

	if r.LastTime == 0 {
		return false
	}

	return true
}

// ValidateUpdate validate revision when updated.
func (r Revision) ValidateUpdate() error {
	if len(r.Modifier) == 0 {
		return errors.New("reviser can not be empty")
	}

	if len(r.Creator) != 0 {
		return errors.New("creator can not be updated")
	}

	now := time.Now().Unix()
	if (r.LastTime <= (now - lagSeconds)) || (r.LastTime >= (now + lagSeconds)) {
		return errors.New("invalid update time")
	}

	if r.LastTime < r.CreateTime-lagSeconds {
		return errors.New("update time must be later than create time")
	}
	return nil
}

// ClusterFields merge the fields of the cluster and the details corresponding to the fields together.
var ClusterFields = mergeFields(ClusterFieldsDescriptor)

// ClusterFieldsDescriptor cluster's fields descriptors.
var ClusterFieldsDescriptor = mergeFieldDescriptors(
	FieldsDescriptors{
		{Field: BKIDField, Type: Numeric, Required: true, IsEditable: false},
		{Field: BKBizIDField, Type: Numeric, Required: true, IsEditable: false},
		{Field: BKSupplierAccountField, Type: String, Required: true, IsEditable: false},
		{Field: CreatorField, Type: String, Required: true, IsEditable: false},
		{Field: ModifierField, Type: String, Required: true, IsEditable: true},
		{Field: CreateTimeField, Type: Numeric, Required: true, IsEditable: false},
		{Field: LastTimeField, Type: Numeric, Required: true, IsEditable: true},
	},
	mergeFieldDescriptors(ClusterSpecFieldsDescriptor),
)

// ClusterSpecFieldsDescriptor cluster spec's fields descriptors.
var ClusterSpecFieldsDescriptor = FieldsDescriptors{
	{Field: KubeNameField, Type: String, Required: true, IsEditable: false},
	{Field: SchedulingEngineField, Type: String, Required: false, IsEditable: false},
	{Field: UidField, Type: String, Required: true, IsEditable: false},
	{Field: XidField, Type: String, Required: false, IsEditable: false},
	{Field: VersionField, Type: String, Required: false, IsEditable: true},
	{Field: NetworkTypeField, Type: Enum, Required: false, IsEditable: true, Option: map[string]string{}},
	{Field: RegionField, Type: String, Required: false, IsEditable: true},
	{Field: VpcField, Type: String, Required: false, IsEditable: false},
	{Field: NetworkField, Type: String, Required: false, IsEditable: false},
	{Field: TypeField, Type: String, Required: false, IsEditable: true},
}

// NamespaceSpecFieldsDescriptor namespace spec's fields descriptors.
var NamespaceSpecFieldsDescriptor = FieldsDescriptors{
	{Field: KubeNameField, Type: String, Required: true, IsEditable: false},
	{Field: LabelsField, Type: MapString, Required: false, IsEditable: true},
	{Field: ClusterUIDField, Type: String, Required: true, IsEditable: false},
	{Field: ResourceQuotasField, Type: Array, Required: false, IsEditable: true},
}

// NodeSpecFieldsDescriptor node spec's fields descriptors.
var NodeSpecFieldsDescriptor = FieldsDescriptors{
	{Field: KubeNameField, Type: String, Required: true, IsEditable: false},
	{Field: RolesField, Type: Enum, Required: false, IsEditable: true},
	{Field: LabelsField, Type: MapString, Required: false, IsEditable: true},
	{Field: TaintsField, Type: MapString, Required: false, IsEditable: true},
	{Field: UnschedulableField, Type: Boolean, Required: false, IsEditable: true},
	{Field: InternalIPField, Type: Array, Required: false, IsEditable: true},
	{Field: ExternalIPField, Type: Array, Required: false, IsEditable: true},
	{Field: HostnameField, Type: String, Required: false, IsEditable: true},
	{Field: RuntimeComponentField, Type: String, Required: false, IsEditable: true},
	{Field: KubeProxyModeField, Type: String, Required: false, IsEditable: true},
	{Field: PodCidrField, Type: String, Required: false, IsEditable: true},
}

// WorkLoadSpecFieldsDescriptor workLoad spec's fields descriptors.
var WorkLoadSpecFieldsDescriptor = FieldsDescriptors{
	{Field: KubeNameField, Type: String, Required: true, IsEditable: false},
	{Field: NamespaceField, Type: String, Required: true, IsEditable: false},
	{Field: LabelsField, Type: MapString, Required: false, IsEditable: true},
	{Field: SelectorField, Type: String, Required: false, IsEditable: true},
	{Field: ReplicasField, Type: Numeric, Required: true, IsEditable: true},
	{Field: StrategyTypeField, Type: String, Required: false, IsEditable: true},
	{Field: MinReadySecondsField, Type: Numeric, Required: false, IsEditable: true},
	{Field: RollingUpdateStrategyField, Type: Object, Required: false, IsEditable: true},
}

// PodSpecFieldsDescriptor pod spec's fields descriptors.
var PodSpecFieldsDescriptor = FieldsDescriptors{
	{Field: KubeNameField, Type: String, Required: true, IsEditable: false},
	{Field: NamespaceField, Type: String, Required: true, IsEditable: false},
	{Field: PriorityField, Type: Numeric, Required: false, IsEditable: true},
	{Field: LabelsField, Type: MapString, Required: false, IsEditable: true},
	{Field: IPField, Type: String, Required: false, IsEditable: true},
	{Field: IPsField, Type: String, Required: false, IsEditable: true},
	{Field: ControlledBy, Type: String, Required: false, IsEditable: true},
	{Field: ContainerUIDField, Type: Array, Required: false, IsEditable: true},
	{Field: VolumesField, Type: Object, Required: false, IsEditable: true},
	{Field: QOSClassField, Type: Enum, Required: false, IsEditable: true},
	{Field: NodeSelectorsField, Type: MapString, Required: false, IsEditable: true},
	{Field: TolerationsField, Type: Object, Required: false, IsEditable: true},
}

// ContainerSpecFieldsDescriptor container spec's fields descriptors.
var ContainerSpecFieldsDescriptor = FieldsDescriptors{
	{Field: KubeNameField, Type: String, Required: true, IsEditable: false},
	{Field: ContainerUIDField, Type: String, Required: true, IsEditable: false},
	{Field: ImageField, Type: String, Required: true, IsEditable: false},
	{Field: PortsField, Type: String, Required: false, IsEditable: true},
	{Field: HostPortsField, Type: String, Required: false, IsEditable: true},
	{Field: ArgsField, Type: String, Required: false, IsEditable: true},
	{Field: StartedField, Type: Numeric, Required: false, IsEditable: true},
	{Field: RequestsField, Type: Object, Required: false, IsEditable: true},
	{Field: LimitsField, Type: Object, Required: false, IsEditable: true},
	{Field: LivenessField, Type: Object, Required: false, IsEditable: true},
	{Field: EnvironmentField, Type: Object, Required: false, IsEditable: true},
	{Field: MountsField, Type: Object, Required: false, IsEditable: true},
}
