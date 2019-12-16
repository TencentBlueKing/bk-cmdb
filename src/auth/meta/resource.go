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

import (
	"configcenter/src/common"
)

type ResourceType string

func (r ResourceType) String() string {
	return string(r)
}

// ResourceType 表示 CMDB 这一侧的资源类型， 对应的有 ResourceTypeID 表示 IAM 一侧的资源类型
// 两者之间有映射关系，详情见 ConvertResourceType
const (
	Business                 ResourceType = "business"
	Model                    ResourceType = "model"
	ModelModule              ResourceType = "modelModule"
	ModelSet                 ResourceType = "modelSet"
	MainlineModel            ResourceType = "mainlineObject"
	MainlineModelTopology    ResourceType = "mainlineObjectTopology"
	MainlineInstanceTopology ResourceType = "mainlineInstanceTopology"
	MainlineInstance         ResourceType = "mainlineInstance"
	AssociationType          ResourceType = "associationType"
	ModelAssociation         ResourceType = "modelAssociation"
	ModelInstanceAssociation ResourceType = "modelInstanceAssociation"
	ModelInstance            ResourceType = "modelInstance"
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
	EventPushing             ResourceType = "eventPushing"
	Plat                     ResourceType = "plat"
	AuditLog                 ResourceType = "auditlog"     // 操作审计
	ResourceSync             ResourceType = "resourceSync" // 云资源发现
	UserCustom               ResourceType = "usercustom"   // 用户自定义
	SystemBase               ResourceType = "systemBase"
	InstallBK                ResourceType = "installBK"
	SetTemplate              ResourceType = "setTemplate"
	OperationStatistic       ResourceType = "operationStatistic" // 运营统计
)

const (
	Host                         = "host"
	ProcessConfigTemplate        = "processConfigTemplate"
	ProcessConfigTemplateVersion = "processConfigTemplateVersion"
	ProcessBoundConfig           = "processBoundConfig"
	SystemFunctionality          = "systemFunctionality"

	NetCollector = "netCollector"
	NetDevice    = "netDevice"
	NetProperty  = "netProperty"
	NetReport    = "netReport"
)

type ResourceDescribe struct {
	Type    ResourceType
	Actions []Action
}

var (
	BusinessDescribe = ResourceDescribe{
		Type:    Business,
		Actions: []Action{Create, Update, Delete, FindMany},
	}

	ModelDescribe = ResourceDescribe{
		Type:    Model,
		Actions: []Action{Create, Update, Delete, FindMany},
	}

	ModelModuleDescribe = ResourceDescribe{
		Type:    ModelModule,
		Actions: []Action{Create, Update, Delete, FindMany},
	}

	ModelSetDescribe = ResourceDescribe{
		Type:    ModelSet,
		Actions: []Action{Create, Update, Delete, FindMany, DeleteMany},
	}

	MainlineModelDescribe = ResourceDescribe{
		Type:    MainlineModel,
		Actions: []Action{Create, Delete, Find},
	}

	MainlineModelTopologyDescribe = ResourceDescribe{
		Type:    MainlineModelTopology,
		Actions: []Action{Find},
	}

	MainlineInstanceTopologyDescribe = ResourceDescribe{
		Type:    MainlineInstanceTopology,
		Actions: []Action{Find},
	}

	AssociationTypeDescribe = ResourceDescribe{
		Type:    AssociationType,
		Actions: []Action{FindMany, Create, Update, Delete},
	}

	ModelAssociationDescribe = ResourceDescribe{
		Type:    ModelAssociation,
		Actions: []Action{FindMany, Create, Update, Delete},
	}

	ModelInstanceAssociationDescribe = ResourceDescribe{
		Type:    ModelInstanceAssociation,
		Actions: []Action{FindMany, Create, Delete},
	}

	ModelInstanceDescribe = ResourceDescribe{
		Type: ModelInstance,
		Actions: []Action{
			DeleteMany,
			FindMany,
			UpdateMany,
			Create,
			Find,
			Update,
			DeleteMany,
			Delete,
			// the following actions is the host actions for only.
			MoveResPoolHostToBizIdleModule,
			MoveHostToBizFaultModule,
			MoveHostToBizIdleModule,
			MoveHostToBizRecycleModule,
			MoveHostFromModuleToResPool,
			MoveHostToAnotherBizModule,
			CleanHostInSetOrModule,
			MoveHostsToBusinessOrModule,
			AddHostToResourcePool,
			MoveBizHostToModule,
		},
	}

	ModelInstanceTopologyDescribe = ResourceDescribe{
		Type:    ModelInstanceTopology,
		Actions: []Action{Find, FindMany},
	}

	ModelTopologyDescribe = ResourceDescribe{
		Type:    ModelTopology,
		Actions: []Action{Find, Update},
	}

	ModelClassificationDescribe = ResourceDescribe{
		Type:    ModelClassification,
		Actions: []Action{FindMany, Create, Update, Delete},
	}

	ModelAttributeGroupDescribe = ResourceDescribe{
		Type:    ModelAttributeGroup,
		Actions: []Action{Find, Create, Delete},
	}

	ModelAttributeDescribe = ResourceDescribe{
		Type:    ModelAttribute,
		Actions: []Action{Find, Create, Update, Delete},
	}

	ModelUniqueDescribe = ResourceDescribe{
		Type:    ModelUnique,
		Actions: []Action{FindMany, Create, Update, Delete},
	}

	HostUserCustomDescribe = ResourceDescribe{
		Type:    UserCustom,
		Actions: []Action{Find, FindMany, Create, Update, Delete},
	}

	HostFavoriteDescribe = ResourceDescribe{
		Type:    HostFavorite,
		Actions: []Action{FindMany, Create, Update, Delete, DeleteMany},
	}

	ProcessDescribe = ResourceDescribe{
		Type:    Process,
		Actions: []Action{Create, Find, FindMany, Delete, DeleteMany, Update, UpdateMany},
	}

	NetDataCollectorDescribe = ResourceDescribe{
		Type:    NetDataCollector,
		Actions: []Action{Find, FindMany, Update, UpdateMany, DeleteMany, Create, DeleteMany},
	}
)

func GetResourceTypeByObjectType(object string) ResourceType {
	switch object {
	case common.BKInnerObjIDApp:
		return Business
	case common.BKInnerObjIDSet:
		return ModelSet
	case common.BKInnerObjIDModule:
		return ModelModule
	default:
		return Model
	}
}
