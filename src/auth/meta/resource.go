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

import "fmt"

type ResourceType string

func (r ResourceType) String() string {
	return string(r)
}

func (r ResourceType) ResourceID(attribute ResourceAttribute) (string, error) {
	switch r {
	case Business:
		return BusinessDescribe.ResourceID(attribute), nil
	case Model:
		return ModelDescribe.ResourceID(attribute), nil
	case ModelModule:
		return ModelModuleDescribe.ResourceID(attribute), nil
	case ModelSet:
		return ModelSetDescribe.ResourceID(attribute), nil
	case MainlineModel:
		return MainlineModelDescribe.ResourceID(attribute), nil
	case MainlineModelTopology:
		return MainlineModelTopologyDescribe.ResourceID(attribute), nil
	case MainlineInstanceTopology:
		return MainlineInstanceTopologyDescribe.ResourceID(attribute), nil
	case AssociationType:
		return AssociationTypeDescribe.ResourceID(attribute), nil
	case ModelAssociation:
		return ModelAssociationDescribe.ResourceID(attribute), nil
	case ModelInstanceAssociation:
		return ModelInstanceAssociationDescribe.ResourceID(attribute), nil
	case ModelInstance:
		return ModelInstanceDescribe.ResourceID(attribute), nil
	case ModelInstanceTopology:
		return ModelInstanceTopologyDescribe.ResourceID(attribute), nil
	case ModelTopology:
		return ModelTopologyDescribe.ResourceID(attribute), nil
	case ModelClassification:
		return ModelClassificationDescribe.ResourceID(attribute), nil
	case ModelAttributeGroup:
		return ModelAttributeGroupDescribe.ResourceID(attribute), nil
	case ModelAttribute:
		return ModelAttributeDescribe.ResourceID(attribute), nil
	case ModelUnique:
		return ModelUniqueDescribe.ResourceID(attribute), nil
	case HostUserCustom:
		return HostUserCustomDescribe.ResourceID(attribute), nil
	case HostFavorite:
		return HostFavoriteDescribe.ResourceID(attribute), nil
	case Process:
		return ProcessDescribe.ResourceID(attribute), nil
	case NetDataCollector:
		return NetDataCollectorDescribe.ResourceID(attribute), nil
	default:
		return "", fmt.Errorf("unsupported resource type: %s", r)
	}

}

const (
	Business                 ResourceType = "business"
	Model                    ResourceType = "model"
	ModelModule              ResourceType = "modelModule"
	ModelSet                 ResourceType = "modelSet"
	MainlineModel            ResourceType = "mainlineObject"
	MainlineModelTopology    ResourceType = "mainlineObjectTopology"
	MainlineInstanceTopology ResourceType = "mainlineInstanceTopology"
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

	HostUserCustom   ResourceType = "hostUserCustom"
	HostFavorite     ResourceType = "hostFavorite"
	Process          ResourceType = "process"
	NetDataCollector ResourceType = "netDataCollector"
)

const (
	Host                         = "host"
	ProcessConfigTemplate        = "processConfigTemplate"
	ProcessConfigTemplateVersion = "processConfigTemplateVersion"
	ProcessBoundConfig           = "processBoundConfig"

	NetCollector = "netCollector"
	NetDevice    = "netDevice"
	NetProperty  = "netProperty"
	NetReport    = "netReport"
)

type ResourceDescribe struct {
	Type    ResourceType
	Actions []Action
	// the rule to generate the resource id to represent this resource.
	rule func(attribute ResourceAttribute) string
}

func (r ResourceDescribe) ResourceID(attribute ResourceAttribute) string {
	return r.rule(attribute)
}

var (
	BusinessDescribe = ResourceDescribe{
		Type:    Business,
		Actions: []Action{Create, Update, Delete, FindMany},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	ModelDescribe = ResourceDescribe{
		Type:    Model,
		Actions: []Action{Create, Update, Delete, FindMany},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	ModelModuleDescribe = ResourceDescribe{
		Type:    ModelModule,
		Actions: []Action{Create, Update, Delete, FindMany},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	ModelSetDescribe = ResourceDescribe{
		Type:    ModelSet,
		Actions: []Action{Create, Update, Delete, FindMany, DeleteMany},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	MainlineModelDescribe = ResourceDescribe{
		Type:    MainlineModel,
		Actions: []Action{Create, Delete, Find},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	MainlineModelTopologyDescribe = ResourceDescribe{
		Type:    MainlineModelTopology,
		Actions: []Action{Find},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	MainlineInstanceTopologyDescribe = ResourceDescribe{
		Type:    MainlineInstanceTopology,
		Actions: []Action{Find},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	AssociationTypeDescribe = ResourceDescribe{
		Type:    AssociationType,
		Actions: []Action{FindMany, Create, Update, Delete},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	ModelAssociationDescribe = ResourceDescribe{
		Type:    ModelAssociation,
		Actions: []Action{FindMany, Create, Update, Delete},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	ModelInstanceAssociationDescribe = ResourceDescribe{
		Type:    ModelInstanceAssociation,
		Actions: []Action{FindMany, Create, Delete},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
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
			MoveHostFromModuleToResPool,
			MoveHostToAnotherBizModule,
			CleanHostInSetOrModule,
			MoveHostsToOrBusinessModule,
			AddHostToResourcePool,
			MoveHostToModule,
		},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	ModelInstanceTopologyDescribe = ResourceDescribe{
		Type:    ModelInstanceTopology,
		Actions: []Action{Find, FindMany},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	ModelTopologyDescribe = ResourceDescribe{
		Type:    ModelTopology,
		Actions: []Action{Find, Update},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	ModelClassificationDescribe = ResourceDescribe{
		Type:    ModelClassification,
		Actions: []Action{FindMany, Create, Update, Delete},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	ModelAttributeGroupDescribe = ResourceDescribe{
		Type:    ModelAttributeGroup,
		Actions: []Action{Find, Create, Delete},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	ModelAttributeDescribe = ResourceDescribe{
		Type:    ModelAttribute,
		Actions: []Action{Find, Create, Update, Delete},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	ModelUniqueDescribe = ResourceDescribe{
		Type:    ModelUnique,
		Actions: []Action{FindMany, Create, Update, Delete},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	HostUserCustomDescribe = ResourceDescribe{
		Type:    HostUserCustom,
		Actions: []Action{Find, FindMany, Create, Update, Delete},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	HostFavoriteDescribe = ResourceDescribe{
		Type:    HostFavorite,
		Actions: []Action{FindMany, Create, Update, Delete, DeleteMany},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	ProcessDescribe = ResourceDescribe{
		Type:    Process,
		Actions: []Action{Create, Find, FindMany, Delete, DeleteMany, Update, UpdateMany, Create},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}

	NetDataCollectorDescribe = ResourceDescribe{
		Type:    NetDataCollector,
		Actions: []Action{Find, FindMany, Update, UpdateMany, DeleteMany, Create, DeleteMany},
		rule: func(attribute ResourceAttribute) string {
			return ""
		},
	}
)
