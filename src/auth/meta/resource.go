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

type ResourceType string

func (r ResourceType) String() string {
	return string(r)
}

const (
	Business                  ResourceType = "business"
	Object                    ResourceType = "object"
	ObjectModule              ResourceType = "objectModule"
	ObjectSet                 ResourceType = "objectSet"
	MainlineObject            ResourceType = "mainlineObject"
	MainlineObjectTopology    ResourceType = "mainlineObjectTopology"
	MainlineInstanceTopology  ResourceType = "mainlineInstanceTopology"
	AssociationType           ResourceType = "associationType"
	ObjectAssociation         ResourceType = "objectAssociation"
	ObjectInstanceAssociation ResourceType = "objectInstanceAssociation"
	ObjectInstance            ResourceType = "objectInstance"
	ObjectInstanceTopology    ResourceType = "objectInstanceTopology"
	ObjectTopology            ResourceType = "objectTopology"
	ObjectClassification      ResourceType = "objectClassification"
	ObjectAttributeGroup      ResourceType = "objectAttributeGroup"
	ObjectAttribute           ResourceType = "objectAttribute"
	ObjectUnique              ResourceType = "objectUnique"

	HostUserCustom        ResourceType = "hostUserCustom"
	HostFavorite          ResourceType = "hostFavorite"
	Host                  ResourceType = "host"
	AddHostToResourcePool ResourceType = "addHostToResourcePool"
	MoveHostToModule      ResourceType = "moveHostToModule"
	// move resource pool hosts to a business idle module
	MoveResPoolHostToBizIdleModule ResourceType = "moveResPoolHostToBizIdleModule"
	MoveHostToBizFaultModule       ResourceType = "moveHostToBizFaultModule"
	MoveHostToBizIdleModule        ResourceType = "moveHostToBizIdleModule"
	MoveHostFromModuleToResPool    ResourceType = "moveHostFromModuleToResPool"
	MoveHostToAnotherBizModule     ResourceType = "moveHostToAnotherBizModule"
	CleanHostInSetOrModule         ResourceType = "cleanHostInSetOrModule"
	MoveHostsToOrBusinessModule    ResourceType = "moveHostsToBusinessOrModule"

	Process                      ResourceType = "process"
	ProcessConfigTemplate        ResourceType = "processConfigTemplate"
	ProcessConfigTemplateVersion ResourceType = "processConfigTemplateVersion"
	ProcessBoundConfig           ResourceType = "processBoundConfig"

	NetDataCollector ResourceType = "netDataCollector"
)

const (
	NetCollector = "netCollector"
	NetDevice    = "netDevice"
	NetProperty  = "netProperty"
	NetReport    = "netReport"
)
