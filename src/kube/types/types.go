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

	// container requests field
	RequestsField = "requests"

	// LivenessField container liveness field
	LivenessField = "liveness"

	// EnvironmentField container environment field
	EnvironmentField = "environment"

	// MountsField container mounts field
	MountsField = "mounts"
)
