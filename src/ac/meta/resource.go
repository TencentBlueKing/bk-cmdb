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

package meta

// ResourceType TODO
type ResourceType string

// String 用于打印
func (r ResourceType) String() string {
	return string(r)
}

// ResourceType 表示 CMDB 这一侧的资源类型， 对应的有 ResourceTypeID 表示 IAM 一侧的资源类型
// 两者之间有映射关系，详情见 ConvertResourceType
const (
	Business                 ResourceType = "business"
	BizSet                   ResourceType = "bizSet"
	Model                    ResourceType = "model"
	ModelModule              ResourceType = "modelModule"
	ModelSet                 ResourceType = "modelSet"
	MainlineModel            ResourceType = "mainlineObject"
	MainlineModelTopology    ResourceType = "mainlineObjectTopology"
	MainlineInstanceTopology ResourceType = "mainlineInstanceTopology"
	MainlineInstance         ResourceType = "mainlineInstance"
	AssociationType          ResourceType = "associationType"
	ModelAssociation         ResourceType = "modelAssociation"
	ModelInstanceTopology    ResourceType = "modelInstanceTopology"
	ModelTopology            ResourceType = "modelTopology"
	ModelClassification      ResourceType = "modelClassification"
	ModelAttributeGroup      ResourceType = "modelAttributeGroup"
	ModelAttribute           ResourceType = "modelAttribute"
	ModelUnique              ResourceType = "modelUnique"
	HostFavorite             ResourceType = "hostFavorite"
	Process                  ResourceType = "process"
	ProcessServiceCategory   ResourceType = "processServiceCategory"
	ProcessServiceTemplate   ResourceType = "processServiceTemplate"
	ProcessTemplate          ResourceType = "processTemplate"
	ProcessServiceInstance   ResourceType = "processServiceInstance"
	BizTopology              ResourceType = "bizTopology"
	HostInstance             ResourceType = "hostInstance"
	NetDataCollector         ResourceType = "netDataCollector"
	DynamicGrouping          ResourceType = "dynamicGrouping" // 动态分组
	EventWatch               ResourceType = "eventWatch"
	CloudAreaInstance        ResourceType = "plat"
	AuditLog                 ResourceType = "auditlog"   // 操作审计
	UserCustom               ResourceType = "usercustom" // 用户自定义
	SystemBase               ResourceType = "systemBase"
	InstallBK                ResourceType = "installBK"
	SystemConfig             ResourceType = "systemConfig"
	SetTemplate              ResourceType = "setTemplate"
	OperationStatistic       ResourceType = "operationStatistic" // 运营统计
	HostApply                ResourceType = "hostApply"
	ResourcePoolDirectory    ResourceType = "resourcePoolDirectory"
	CloudAccount             ResourceType = "cloudAccount"
	CloudResourceTask        ResourceType = "cloudResourceTask"
	ConfigAdmin              ResourceType = "configAdmin"
)

// kube related auth resource in CMDB
const (
	// KubeCluster auth resource type in CMDB
	KubeCluster ResourceType = "kube_cluster"

	// KubeNode auth resource type in CMDB
	KubeNode ResourceType = "kube_node"

	// KubeNamespace auth resource type in CMDB
	KubeNamespace ResourceType = "kube_namespace"

	// KubeWorkload auth resource type in CMDB, including deployment, statefulSet, daemonSet ...
	KubeWorkload ResourceType = "kube_workload"

	// KubePod auth resource type in CMDB
	KubePod ResourceType = "kube_pod"

	// KubeContainer auth resource type in CMDB
	KubeContainer ResourceType = "kube_container"

	// below are specific workload auth resource types in CMDB, reserved for later use

	// KubeDeployment auth resource type in CMDB
	KubeDeployment ResourceType = "kube_deployment"

	// KubeStatefulSet auth resource type in CMDB
	KubeStatefulSet ResourceType = "kube_statefulSet"

	// KubeDaemonSet auth resource type in CMDB
	KubeDaemonSet ResourceType = "kube_daemonSet"

	// KubeGameStatefulSet auth resource type in CMDB
	KubeGameStatefulSet ResourceType = "kube_gameStatefulSet"

	// KubeGameDeployment auth resource type in CMDB
	KubeGameDeployment ResourceType = "kube_gameDeployment"

	// KubeCronJob auth resource type in CMDB
	KubeCronJob ResourceType = "kube_cronJob"

	// KubeJob auth resource type in CMDB
	KubeJob ResourceType = "kube_job"

	// KubePodWorkload pod workload auth resource type in CMDB
	KubePodWorkload ResourceType = "kube_pods"
)

const (
	// CMDBSysInstTypePrefix TODO
	// CMDB侧资源的通用模型实例前缀标识
	CMDBSysInstTypePrefix = "comobj_"
)

const (
	// NetCollector TODO
	NetCollector = "netCollector"
	// NetDevice TODO
	NetDevice = "netDevice"
	// NetProperty TODO
	NetProperty = "netProperty"
	// NetReport TODO
	NetReport = "netReport"
)

// ResourceDescribe TODO
type ResourceDescribe struct {
	Type    ResourceType
	Actions []Action
}
