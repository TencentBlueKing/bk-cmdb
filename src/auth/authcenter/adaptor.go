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

package authcenter

import (
	"fmt"

	"configcenter/src/auth/meta"
)

// Adaptor is a middleware wrapper which works for converting concepts
// between bk-cmdb and blueking auth center. Especially the policies
// in auth center.

func adaptor(attribute *meta.ResourceAttribute) (*ResourceInfo, error) {
	resourceType := attribute.Basic.Type
	info := new(ResourceInfo)
	info.ResourceName = attribute.Basic.Name

	var err error
	info.ResourceID, err = GenerateResourceID(attribute)
	if err != nil {
		return nil, err
	}

	switch resourceType {
	case meta.Business:
		return info, nil

	case meta.Model,
		meta.ModelModule,
		meta.ModelSet,
		meta.MainlineModel,
		meta.ModelUnique,
		meta.ModelClassification,
		meta.ModelAttribute,
		meta.ModelAttributeGroup:
		if attribute.BusinessID == 0 {
			info.ResourceType = AppModel
		} else {
			info.ResourceType = ModelManagement
		}

	case meta.AssociationType:
		info.ResourceType = AssociationType

	case meta.ModelAssociation:

	case meta.ModelInstanceAssociation:

	case meta.ModelInstance, meta.ModelInstanceTopology:
		info.ResourceType = InstanceManagement

	case meta.HostUserCustom:
		info.ResourceType = CustomQuery

	case meta.HostFavorite:

	case meta.Process:
		info.ResourceType = Process

	case meta.NetDataCollector:
		return nil, fmt.Errorf("unsupported resource type: %s", attribute.Basic.Type)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", attribute.Basic.Type)
	}

	return info, nil
}

// type is resource's type in auth center.
type Type string

const (
	// the alias name maybe "dynamic classification"
	CustomQuery        Type = "customQuery"
	AppModel           Type = "appModel"
	Host               Type = "host"
	Process            Type = "process"
	Topology           Type = "topology"
	AppInstance        Type = "appInstance"
	InstanceManagement Type = "instanceManagement"
	ModelManagement    Type = "modelManagement"
	AssociationType    Type = "associationType"
	ModelGroup         Type = "modelGroup"
	Event              Type = "event"
	SystemBase         Type = "systemBase"
)

type Action string

const (
	// unknown action is a action that can not be recognized by the
	// auth center.
	Unknown Action = "unknown"
	Edit    Action = "edit"
	Create  Action = "create"
	Get     Action = "get"
	Delete  Action = "delete"
	// host action
	ModuleTransfer Action = "moduleTransfer"
	// business topology action
	HostTransfer Action = "hostTransfer"
)

func adaptorAction(r *meta.ResourceAttribute) Action {
	switch r.Action {
	case meta.Create, meta.CreateMany:
		return Create

	case meta.Find, meta.FindMany:
		return Get

	case meta.Delete, meta.DeleteMany:
		return Delete

	case meta.Update, meta.UpdateMany:
		return Edit

	case meta.MoveResPoolHostToBizIdleModule,
		meta.MoveHostToBizFaultModule,
		meta.MoveHostToBizIdleModule,
		meta.MoveHostFromModuleToResPool,
		meta.MoveHostToAnotherBizModule,
		meta.CleanHostInSetOrModule,
		meta.MoveHostsToOrBusinessModule,
		meta.AddHostToResourcePool,
		meta.MoveHostToModule:

		return ModuleTransfer

	// TODO: add host transfer adaptor rule.

	default:
		return Unknown
	}
}

type ResourceDetail struct {
	// the resource type in auth center.
	Type Type
	// all the actions that this resource supported.
	Actions []Action
}

var (
	CustomQueryDescribe = ResourceDetail{
		Type:    CustomQuery,
		Actions: []Action{Get, Delete, Edit, Create},
	}

	AppModelDescribe = ResourceDetail{
		Type:    AppModel,
		Actions: []Action{Get, Delete, Edit, Create},
	}

	HostDescribe = ResourceDetail{
		Type:    Host,
		Actions: []Action{Get, Delete, Edit, Create, ModuleTransfer},
	}

	ProcessDescribe = ResourceDetail{
		Type:    Process,
		Actions: []Action{Get, Delete, Edit, Create},
	}

	TopologyDescribe = ResourceDetail{
		Type:    Topology,
		Actions: []Action{Get, Delete, Edit, Create, HostTransfer},
	}

	AppInstanceDescribe = ResourceDetail{
		Type:    AppInstance,
		Actions: []Action{Get, Delete, Edit, Create},
	}

	InstanceManagementDescribe = ResourceDetail{
		Type:    InstanceManagement,
		Actions: []Action{Get, Delete, Edit, Create},
	}

	ModelManagementDescribe = ResourceDetail{
		Type:    ModelManagement,
		Actions: []Action{Get, Delete, Edit, Create},
	}

	AssociationTypeDescribe = ResourceDetail{
		Type:    AssociationType,
		Actions: []Action{Get, Delete, Edit, Create},
	}

	ModelGroupDescribe = ResourceDetail{
		Type:    ModelGroup,
		Actions: []Action{Get, Delete, Edit, Create},
	}

	EventDescribe = ResourceDetail{
		Type:    Event,
		Actions: []Action{Get, Delete, Edit, Create},
	}

	SystemBaseDescribe = ResourceDetail{
		Type:    SystemBase,
		Actions: []Action{Get, Delete, Edit, Create},
	}
)
