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
	"fmt"

	"configcenter/src/storage/dal/table"
)

// identification of k8s in cc
const (
	// KubeBusiness k8s business type
	KubeBusiness = "biz"

	// KubeCluster k8s cluster type
	KubeCluster = "cluster"

	// KubeNode k8s node type
	KubeNode = "node"

	// KubeNamespace k8s namespace type
	KubeNamespace = "namespace"

	// KubeFolder k8s folder type
	KubeFolder = "folder"

	// KubeWorkload k8s workload type
	KubeWorkload = "workload"

	// KubePod k8s pod type
	KubePod = "pod"

	// KubeContainer k8s container type
	KubeContainer = "container"
)

// WorkloadType workload type enum
type WorkloadType string

// Validate validate WorkloadType
func (t WorkloadType) Validate() error {
	switch t {
	case KubeDeployment, KubeStatefulSet, KubeDaemonSet,
		KubeGameStatefulSet, KubeGameDeployment, KubeCronJob,
		KubeJob, KubePodWorkload:
		return nil
	default:
		return fmt.Errorf("can not support this type of workload, kind: %s", t)
	}
}

// Table get the table name based on the workload type
func (t WorkloadType) Table() (string, error) {
	switch t {
	case KubeDeployment:
		return BKTableNameBaseDeployment, nil

	case KubeStatefulSet:
		return BKTableNameBaseStatefulSet, nil

	case KubeDaemonSet:
		return BKTableNameBaseDaemonSet, nil

	case KubeGameStatefulSet:
		return BKTableNameGameStatefulSet, nil

	case KubeGameDeployment:
		return BKTableNameGameDeployment, nil

	case KubeCronJob:
		return BKTableNameBaseCronJob, nil

	case KubeJob:
		return BKTableNameBaseJob, nil

	case KubePodWorkload:
		return BKTableNameBasePodWorkload, nil

	default:
		return "", fmt.Errorf("can not find table name, kind: %s", t)
	}
}

// Fields get the workload type related table fields
func (t WorkloadType) Fields() (*table.Fields, error) {
	switch t {
	case KubeDeployment:
		return DeploymentFields, nil

	case KubeStatefulSet:
		return StatefulSetFields, nil

	case KubeDaemonSet:
		return DaemonSetFields, nil

	case KubeGameStatefulSet:
		return GameStatefulSetFields, nil

	case KubeGameDeployment:
		return GameDeploymentFields, nil

	case KubeCronJob:
		return CronJobFields, nil

	case KubeJob:
		return JobFields, nil

	case KubePodWorkload:
		return PodsWorkloadFields, nil

	default:
		return nil, fmt.Errorf("workload type %s is not supported", t)
	}
}

// NewInst new a workload instance according to workload type
func (t WorkloadType) NewInst() (WorkloadInterface, error) {
	switch t {
	case KubeDeployment:
		return new(Deployment), nil

	case KubeStatefulSet:
		return new(StatefulSet), nil

	case KubeDaemonSet:
		return new(DaemonSet), nil

	case KubeGameDeployment:
		return new(GameDeployment), nil

	case KubeGameStatefulSet:
		return new(GameStatefulSet), nil

	case KubeCronJob:
		return new(CronJob), nil

	case KubeJob:
		return new(Job), nil

	case KubePodWorkload:
		return new(PodsWorkload), nil

	default:
		return nil, fmt.Errorf("workload type %s is not supported", t)
	}
}

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

	// BKTableNameBaseWorkload virtual table name of Workload, specific workload data are stored in separate tables
	BKTableNameBaseWorkload = "cc_WorkloadBase"

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

const (
	// KubeHostKind host kind
	KubeHostKind = "host"
	// KubePodKind pod kind
	KubePodKind = "pod"
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

	// ClusterEnvironmentField cluster environment field
	ClusterEnvironmentField = "environment"

	// NetworkField cluster network field
	NetworkField = "network"

	// TypeField cluster type field
	TypeField = "type"
	// SchedulingEngineField scheduling engine
	SchedulingEngineField = "scheduling_engine"

	// todo: 后续项目的pr合入之后，需要统一定义

	// ProjectNameField project name field
	ProjectNameField = "bk_project_name"

	// ProjectIDField project id field
	ProjectIDField = "bk_project_id"

	// ProjectCodeField project code field
	ProjectCodeField = "bk_project_code"
)

// node field names
const (
	// RolesField node role field
	RolesField = "roles"

	// TaintsField node taints field
	TaintsField = "taints"

	// HasPodField node taints field
	HasPodField = "has_pod"

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
	NodeField = "node_name"
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

	// RefField pod relate workload field
	RefField = "ref"

	// RefIDField pod relate workload id field
	RefIDField = "ref.id"

	// RefNameField pod relate workload name field
	RefNameField = "ref.name"

	// RefKindField pod relate workload kind field
	RefKindField = "ref.kind"

	// NodeNameFiled pod relate node name field
	NodeNameFiled = "node_name"
)

const (
	// KubeFolderID 每个cluster只有唯一的 folder，此节点没有表，统一用 999 表示 folder的ID，如果需要确认具体的 folder
	// 需要与clusterID结合起来使用
	KubeFolderID   = 999
	KubeFolderName = "空Pod节点"
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

	// RequestsField requests field
	RequestsField = "requests"

	// LivenessField container liveness field
	LivenessField = "liveness"

	// EnvironmentField container environment field
	EnvironmentField = "environment"

	// MountsField container mounts field
	MountsField = "mounts"
)
