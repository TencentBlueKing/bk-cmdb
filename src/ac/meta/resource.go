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
