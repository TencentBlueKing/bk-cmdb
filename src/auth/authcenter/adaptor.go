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
	"errors"
	"fmt"

	"configcenter/src/auth/meta"
	"configcenter/src/common/metadata"
)

var NotEnoughLayer = fmt.Errorf("not enough layer")

// ResourceTypeID is resource's type in auth center.
func adaptor(attribute *meta.ResourceAttribute) (*ResourceInfo, error) {
	var err error
	info := new(ResourceInfo)
	info.ResourceName = attribute.Basic.Name

	resourceTypeID, err := convertResourceType(attribute.Type, attribute.BusinessID)
	if err != nil {
		return nil, err
	}
	info.ResourceType = *resourceTypeID

	info.ResourceID, err = GenerateResourceID(info.ResourceType, attribute)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// Adaptor is a middleware wrapper which works for converting concepts
// between bk-cmdb and blueking auth center. Especially the policies
// in auth center.
func convertResourceType(resourceType meta.ResourceType, businessID int64) (*ResourceTypeID, error) {
	var iamResourceType ResourceTypeID
	switch resourceType {
	case meta.Business:
		iamResourceType = SysBusinessInstance

	case meta.Model,
		meta.ModelUnique,
		meta.ModelAttribute,
		meta.ModelAttributeGroup:
		if businessID > 0 {
			iamResourceType = BizModel
		} else {
			iamResourceType = SysModel
		}

	case meta.ModelModule, meta.ModelSet, meta.MainlineInstance, meta.MainlineInstanceTopology:
		iamResourceType = BizTopology

	case meta.MainlineModel, meta.ModelTopology:
		iamResourceType = SysSystemBase

	case meta.ModelClassification:
		if businessID > 0 {
			iamResourceType = BizModelGroup
		} else {
			iamResourceType = SysModelGroup
		}

	case meta.AssociationType:
		iamResourceType = SysAssociationType

	case meta.ModelAssociation:
		return nil, errors.New("model association does not support auth now")

	case meta.ModelInstanceAssociation:
		return nil, errors.New("model instance association does not support  auth now")
	case meta.MainlineModelTopology:
		iamResourceType = SysSystemBase

	case meta.ModelInstance:
		if businessID <= 0 {
			iamResourceType = SysInstance
		} else {
			iamResourceType = BizInstance
		}

	case meta.Plat:
		iamResourceType = SysInstance
	case meta.HostInstance:
		if businessID <= 0 {
			iamResourceType = SysHostInstance
		} else {
			iamResourceType = BizHostInstance
		}

	case meta.HostFavorite:
		return nil, errors.New("host favorite does not support auth now")

	case meta.Process:
		iamResourceType = BizProcessInstance
	case meta.EventPushing:
		iamResourceType = SysEventPushing
	case meta.DynamicGrouping:
		iamResourceType = BizCustomQuery
	case meta.AuditLog:
		if businessID <= 0 {
			iamResourceType = SysAuditLog
		} else {
			iamResourceType = BizAuditLog
		}
	case meta.SystemBase:
		iamResourceType = SysSystemBase
	case meta.UserCustom:
		iamResourceType = UserCustom
	case meta.NetDataCollector:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	return &iamResourceType, nil
}

// ResourceTypeID is resource's type in auth center.
type ResourceTypeID string

// System Resource
const (
	SysSystemBase       ResourceTypeID = "sys_system_base"
	SysBusinessInstance ResourceTypeID = "sys_business_instance"
	SysHostInstance     ResourceTypeID = "sys_host_instance"
	SysEventPushing     ResourceTypeID = "sys_event_pushing"
	SysModelGroup       ResourceTypeID = "sys_model_group"
	SysModel            ResourceTypeID = "sys_model"
	SysInstance         ResourceTypeID = "sys_instance"
	SysAssociationType  ResourceTypeID = "sys_association_type"
	SysAuditLog         ResourceTypeID = "sys_audit_log"
)

// Business Resource
const (
	// the alias name maybe "dynamic classification"
	BizCustomQuery     ResourceTypeID = "biz_custom_query"
	BizHostInstance    ResourceTypeID = "biz_host_instance"
	BizProcessInstance ResourceTypeID = "biz_process_instance"
	BizTopology        ResourceTypeID = "biz_topology"
	BizModelGroup      ResourceTypeID = "biz_model_group"
	BizModel           ResourceTypeID = "biz_model"
	BizInstance        ResourceTypeID = "biz_instance"
	BizAuditLog        ResourceTypeID = "biz_audit_log"
)

const (
	UserCustom ResourceTypeID = "userCustom"
)

var ResourceTypeIDMap = map[ResourceTypeID]string{
	SysSystemBase:       "系统基础",
	SysBusinessInstance: "业务",
	SysHostInstance:     "主机",
	SysEventPushing:     "事件推送",
	SysModelGroup:       "模型分级",
	SysModel:            "模型",
	SysInstance:         "实例",
	SysAssociationType:  "关联类型",
	SysAuditLog:         "操作审计",
	BizCustomQuery:      "动态分组",
	BizHostInstance:     "业务主机",
	BizProcessInstance:  "进程",
	// TODO: delete this when upgrade to v3.5.x
	BizTopology:   "拓扑",
	BizModelGroup: "模型分组",
	BizModel:      "模型",
	BizInstance:   "实例",
	BizAuditLog:   "操作审计",
	UserCustom:    "",
}

type ActionID string

// ActionID define
const (
	// Unknown action is a action that can not be recognized by the auth center.
	Unknown ActionID = "unknown"
	Edit    ActionID = "edit"
	Create  ActionID = "create"
	Get     ActionID = "get"
	Delete  ActionID = "delete"

	// Archive for business
	Archive ActionID = "archive"
	// host action
	ModuleTransfer ActionID = "module_transfer"
	// business topology action
	HostTransfer ActionID = "host_transfer"
	// system base action, related to model topology
	ModelTopologyView ActionID = "model_topology_view"
	// business model topology operation.
	ModelTopologyOperation ActionID = "model_topology_operation"
	// assign host(s) to a business
	// located system/host/assignHostsToBusiness in auth center.
	AssignHostsToBusiness ActionID = "assign_hosts_to_business"
	BindModule            ActionID = "bind_module"
	AdminEntrance         ActionID = "admin_entrance"
)

var ActionIDNameMap = map[ActionID]string{
	Unknown:                "未知操作",
	Edit:                   "编辑",
	Create:                 "新建",
	Get:                    "查询",
	Delete:                 "删除",
	Archive:                "归档",
	ModelTopologyOperation: "拓扑层级管理",
	// TODO: delete this when upgrade to v3.5.x
	BindModule: "绑定到模块",
}

func adaptorAction(r *meta.ResourceAttribute) (ActionID, error) {
	if r.Basic.Type == meta.ModelAttributeGroup ||
		r.Basic.Type == meta.ModelUnique ||
		r.Basic.Type == meta.ModelAttribute {
		if r.Action == meta.Delete || r.Action == meta.Update || r.Action == meta.Create {
			return Edit, nil
		}
	}

	if r.Basic.Type == meta.Business {
		if r.Action == meta.Archive {
			return Archive, nil
		}

		// edit a business.
		if r.Action == meta.Create {
			return Create, nil
		}

		if r.Action == meta.Update {
			return Edit, nil
		}
	}

	if r.Action == meta.Find || r.Action == meta.Delete || r.Action == meta.Create {
		if r.Basic.Type == meta.MainlineModel {
			return ModelTopologyOperation, nil
		}
	}

	// if r.Basic.Type == meta.ModelModule || r.Basic.Type == meta.ModelSet || r.Basic.Type == meta.MainlineInstance {
	// 	return ModelTopologyOperation, nil
	// }

	if r.Action == meta.Find || r.Action == meta.Update {
		if r.Basic.Type == meta.ModelTopology {
			return ModelTopologyView, nil
		}
		if r.Basic.Type == meta.MainlineModelTopology {
			return ModelTopologyOperation, nil
		}

	}

	if r.Basic.Type == meta.Process {
		if r.Action == meta.BoundModuleToProcess || r.Action == meta.UnboundModuleToProcess {
			return Edit, nil
		}
	}

	if r.Basic.Type == meta.HostInstance {
		if r.Action == meta.MoveResPoolHostToBizIdleModule {
			return Edit, nil
		}

		if r.Action == meta.AddHostToResourcePool {
			return Create, nil
		}

		if r.Action == meta.MoveResPoolHostToBizIdleModule {
			return Edit, nil
		}
	}

	switch r.Action {
	case meta.Create, meta.CreateMany:
		return Create, nil

	case meta.Find, meta.FindMany:
		return Get, nil

	case meta.Delete, meta.DeleteMany:
		return Delete, nil

	case meta.Update, meta.UpdateMany:
		return Edit, nil

	case meta.MoveResPoolHostToBizIdleModule:
		if r.Basic.Type == meta.ModelInstance && r.Basic.Name == meta.Host {
			return Edit, nil
		}

	case meta.MoveHostToBizFaultModule,
		meta.MoveHostToBizIdleModule,
		meta.MoveHostToAnotherBizModule,
		meta.CleanHostInSetOrModule,
		meta.TransferHost,
		meta.MoveBizHostToModule:
		return Edit, nil

	case meta.MoveHostFromModuleToResPool:
		return Delete, nil

	case meta.MoveHostsToBusinessOrModule:
		return Edit, nil
	case meta.ModelTopologyView:
		return ModelTopologyView, nil
	case meta.ModelTopologyOperation:
		return ModelTopologyOperation, nil
	case meta.AdminEntrance:
		return AdminEntrance, nil
	}

	return Unknown, fmt.Errorf("unsupported action: %s", r.Action)
}

// TODO: add multiple language support
func AdoptPermissions(rs []meta.ResourceAttribute) ([]metadata.Permission, error) {

	ps := make([]metadata.Permission, 0)
	for _, r := range rs {
		var p metadata.Permission
		p.SystemID = SystemIDCMDB
		p.SystemName = SystemNameCMDB

		if r.BusinessID > 0 {
			p.ScopeType = ScopeTypeIDBiz
			p.ScopeTypeName = ScopeTypeIDBizName
		} else {
			p.ScopeType = ScopeTypeIDSystem
			p.ScopeTypeName = ScopeTypeIDSystemName
		}

		actID, err := adaptorAction(&r)
		if err != nil {
			return nil, err
		}
		p.ActionID = string(actID)
		p.ActionName = ActionIDNameMap[actID]

		rscType, err := convertResourceType(r.Basic.Type, r.BusinessID)
		if err != nil {
			return nil, err
		}

		rscIDs, err := GenerateResourceID(*rscType, &r)
		if err != nil {
			return nil, err
		}

		var rsc metadata.Resource
		rsc.ResourceType = string(*rscType)
		rsc.ResourceTypeName = ResourceTypeIDMap[*rscType]
		if len(rscIDs) != 0 {
			rsc.ResourceID = rscIDs[0].ResourceID
		}
		rsc.ResourceName = r.Basic.Name
		p.Resources = [][]metadata.Resource{{rsc}}
		ps = append(ps, p)
	}
	return ps, nil
}

type ResourceDetail struct {
	// the resource type in auth center.
	Type ResourceTypeID
	// all the actions that this resource supported.
	Actions []ActionID
}

var (
	CustomQueryDescribe = ResourceDetail{
		Type:    BizCustomQuery,
		Actions: []ActionID{Get, Delete, Edit, Create},
	}

	AppModelDescribe = ResourceDetail{
		Type:    BizModel,
		Actions: []ActionID{Get, Delete, Edit, Create},
	}

	HostDescribe = ResourceDetail{
		Type:    BizHostInstance,
		Actions: []ActionID{Get, Delete, Edit, Create, ModuleTransfer},
	}

	ProcessDescribe = ResourceDetail{
		Type:    BizProcessInstance,
		Actions: []ActionID{Get, Delete, Edit, Create},
	}

	TopologyDescribe = ResourceDetail{
		Type:    BizTopology,
		Actions: []ActionID{Get, Delete, Edit, Create, HostTransfer},
	}

	AppInstanceDescribe = ResourceDetail{
		Type:    BizInstance,
		Actions: []ActionID{Get, Delete, Edit, Create},
	}

	InstanceManagementDescribe = ResourceDetail{
		Type:    SysInstance,
		Actions: []ActionID{Get, Delete, Edit, Create},
	}

	ModelManagementDescribe = ResourceDetail{
		Type:    SysModel,
		Actions: []ActionID{Get, Delete, Edit, Create},
	}

	AssociationTypeDescribe = ResourceDetail{
		Type:    SysAssociationType,
		Actions: []ActionID{Get, Delete, Edit, Create},
	}

	ModelGroupDescribe = ResourceDetail{
		Type:    SysModelGroup,
		Actions: []ActionID{Get, Delete, Edit, Create},
	}

	EventDescribe = ResourceDetail{
		Type:    SysEventPushing,
		Actions: []ActionID{Get, Delete, Edit, Create},
	}

	SystemBaseDescribe = ResourceDetail{
		Type:    SysSystemBase,
		Actions: []ActionID{Get, Delete, Edit, Create},
	}
)
