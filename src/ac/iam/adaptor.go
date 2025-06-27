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

package iam

import (
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/ac/iam/types"
	"configcenter/src/ac/meta"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/apigw/iam"
)

// NotEnoughLayer TODO
var NotEnoughLayer = fmt.Errorf("not enough layer")

// AdaptAuthOptions TODO
func AdaptAuthOptions(a *meta.ResourceAttribute) (types.ActionID, []iam.Resource, error) {

	var action types.ActionID

	action, err := ConvertResourceAction(a.Type, a.Action, a.BusinessID)
	if err != nil {
		return "", nil, err
	}

	// convert different cmdb resource's to resource's registered to iam
	rscType, err := ConvertResourceType(a.Type, a.BusinessID)
	if err != nil {
		return "", nil, err
	}

	resource, err := GenIamResource(action, *rscType, a)
	if err != nil {
		return "", nil, err
	}

	return action, resource, nil
}

var ccIamResTypeMap = map[meta.ResourceType]types.TypeID{
	meta.Business:                 types.Business,
	meta.BizSet:                   types.BizSet,
	meta.Project:                  types.Project,
	meta.Model:                    types.SysModel,
	meta.ModelUnique:              types.SysModel,
	meta.ModelAttributeGroup:      types.SysModel,
	meta.ModelModule:              types.BizTopology,
	meta.ModelSet:                 types.BizTopology,
	meta.MainlineInstance:         types.BizTopology,
	meta.MainlineInstanceTopology: types.BizTopology,
	meta.MainlineModel:            types.TypeID(""),
	meta.ModelTopology:            types.TypeID(""),
	meta.ModelClassification:      types.SysModelGroup,
	meta.AssociationType:          types.SysAssociationType,
	meta.ModelAssociation:         types.SysModel,
	meta.MainlineModelTopology:    types.TypeID(""),
	meta.ModelInstanceTopology:    types.SkipType,
	meta.CloudAreaInstance:        types.SysCloudArea,
	meta.HostInstance:             types.Host,
	meta.HostFavorite:             types.SkipType,
	meta.Process:                  types.BizProcessServiceInstance,
	meta.DynamicGrouping:          types.BizCustomQuery,
	meta.AuditLog:                 types.SysAuditLog,
	meta.SystemBase:               types.TypeID(""),
	meta.UserCustom:               types.UserCustom,
	meta.ProcessServiceTemplate:   types.BizProcessServiceTemplate,
	meta.ProcessServiceCategory:   types.BizProcessServiceCategory,
	meta.ProcessServiceInstance:   types.BizProcessServiceInstance,
	meta.BizTopology:              types.BizTopology,
	meta.SetTemplate:              types.BizSetTemplate,
	meta.HostApply:                types.BizHostApply,
	meta.ResourcePoolDirectory:    types.SysResourcePoolDirectory,
	meta.EventWatch:               types.SysEventWatch,
	meta.ConfigAdmin:              types.TypeID(""),
	meta.SystemConfig:             types.TypeID(""),
	meta.KubeCluster:              types.TypeID(""),
	meta.KubeNode:                 types.TypeID(""),
	meta.KubeNamespace:            types.TypeID(""),
	meta.KubeWorkload:             types.TypeID(""),
	meta.KubeDeployment:           types.TypeID(""),
	meta.KubeStatefulSet:          types.TypeID(""),
	meta.KubeDaemonSet:            types.TypeID(""),
	meta.KubeGameStatefulSet:      types.TypeID(""),
	meta.KubeGameDeployment:       types.TypeID(""),
	meta.KubeCronJob:              types.TypeID(""),
	meta.KubeJob:                  types.TypeID(""),
	meta.KubePodWorkload:          types.TypeID(""),
	meta.KubePod:                  types.TypeID(""),
	meta.KubeContainer:            types.TypeID(""),
	meta.FieldTemplate:            types.FieldGroupingTemplate,
	meta.FulltextSearch:           types.TypeID(""),
	meta.IDRuleIncrID:             types.TypeID(""),
	meta.FullSyncCond:             types.TypeID(""),
	meta.GeneralCache:             types.GeneralCache,
	meta.TenantSet:                types.TenantSet,
}

// ConvertResourceType convert resource type from CMDB to IAM
func ConvertResourceType(resourceType meta.ResourceType, businessID int64) (*types.TypeID, error) {
	var iamResourceType types.TypeID

	switch resourceType {
	case meta.ModelAttribute:
		if businessID > 0 {
			iamResourceType = types.BizCustomField
		} else {
			iamResourceType = types.SysModel
		}
		return &iamResourceType, nil
	case meta.NetDataCollector:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	iamResourceType, exists := ccIamResTypeMap[resourceType]
	if exists {
		return &iamResourceType, nil
	}

	if IsCMDBSysInstance(resourceType) {
		iamResourceType = types.TypeID(resourceType)
		return &iamResourceType, nil
	}

	return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
}

// ConvertResourceAction convert resource action from CMDB to IAM
func ConvertResourceAction(resourceType meta.ResourceType, action meta.Action, businessID int64) (types.ActionID,
	error) {
	if action == meta.SkipAction {
		return types.Skip, nil
	}

	convertAction := action
	switch action {
	case meta.CreateMany:
		convertAction = meta.Create
	case meta.FindMany:
		convertAction = meta.Find
	case meta.DeleteMany:
		convertAction = meta.Delete
	case meta.UpdateMany:
		convertAction = meta.Update
	}

	switch resourceType {
	case meta.ModelAttribute, meta.ModelAttributeGroup:
		switch convertAction {
		case meta.Delete, meta.Update, meta.Create:
			if businessID > 0 {
				return types.EditBusinessCustomField, nil
			} else {
				return types.EditSysModel, nil
			}
		}
	case meta.HostInstance:
		switch convertAction {
		case meta.Update:
			if businessID > 0 {
				return types.EditBusinessHost, nil
			} else {
				return types.EditResourcePoolHost, nil
			}
		case meta.Find:
			if businessID > 0 {
				return types.ViewBusinessResource, nil
			} else {
				return types.ViewResourcePoolHost, nil
			}
		}
	}

	if IsCMDBSysInstance(resourceType) {
		return ConvertSysInstanceActionID(resourceType, convertAction)
	}

	if _, exist := resourceActionMap[resourceType]; exist {
		actionID, ok := resourceActionMap[resourceType][convertAction]
		if ok && actionID != types.Unsupported {
			return actionID, nil
		}
	}

	return types.Unsupported, fmt.Errorf("unsupported type %s action: %s", resourceType, action)
}

// ConvertSysInstanceActionID convert system instances action from CMDB to IAM
func ConvertSysInstanceActionID(resourceType meta.ResourceType, action meta.Action) (types.ActionID, error) {
	var actionType types.ActionType
	switch action {
	case meta.Create:
		actionType = types.Create
	case meta.Update:
		actionType = types.Edit
	case meta.Delete:
		actionType = types.Delete
	case meta.Find:
		actionType = types.View
	default:
		return types.Unsupported, fmt.Errorf("unsupported action: %s", action)
	}
	id := strings.TrimPrefix(string(resourceType), meta.CMDBSysInstTypePrefix)
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		return types.Unsupported, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
	return types.ActionID(fmt.Sprintf("%s_%s%s", actionType, types.IAMSysInstTypePrefix, id)), nil
}

var resourceActionMap = map[meta.ResourceType]map[meta.Action]types.ActionID{
	meta.ModelAttributeGroup: {
		meta.Delete:   types.EditSysModel,
		meta.Update:   types.EditSysModel,
		meta.Create:   types.EditSysModel,
		meta.Find:     types.ViewSysModel,
		meta.FindMany: types.ViewSysModel,
	},
	meta.ModelUnique: {
		meta.Delete:   types.EditSysModel,
		meta.Update:   types.EditSysModel,
		meta.Create:   types.EditSysModel,
		meta.Find:     types.ViewSysModel,
		meta.FindMany: types.ViewSysModel,
	},
	meta.Business: {
		meta.Archive:              types.ArchiveBusiness,
		meta.Create:               types.CreateBusiness,
		meta.Update:               types.EditBusiness,
		meta.Find:                 types.FindBusiness,
		meta.ViewBusinessResource: types.ViewBusinessResource,
	},
	meta.BizSet: {
		meta.Create:       types.CreateBizSet,
		meta.Update:       types.EditBizSet,
		meta.Delete:       types.DeleteBizSet,
		meta.Find:         types.ViewBizSet,
		meta.AccessBizSet: types.AccessBizSet,
	},
	meta.DynamicGrouping: {
		meta.Delete:   types.DeleteBusinessCustomQuery,
		meta.Update:   types.EditBusinessCustomQuery,
		meta.Create:   types.CreateBusinessCustomQuery,
		meta.Find:     types.ViewBusinessResource,
		meta.FindMany: types.ViewBusinessResource,
		meta.Execute:  types.ViewBusinessResource,
	},
	meta.MainlineModel: {
		meta.Find:   types.Skip,
		meta.Create: types.EditBusinessLayer,
		meta.Delete: types.EditBusinessLayer,
	},
	meta.ModelTopology: {
		meta.Find:              types.EditModelTopologyView,
		meta.Update:            types.EditModelTopologyView,
		meta.ModelTopologyView: types.ViewModelTopo,
	},
	meta.MainlineModelTopology: {
		meta.Find: types.Skip,
	},
	meta.Process: {
		meta.Find:   types.Skip,
		meta.Create: types.EditBusinessServiceInstance,
		meta.Delete: types.EditBusinessServiceInstance,
		meta.Update: types.EditBusinessServiceInstance,
	},
	meta.HostInstance: {
		meta.MoveResPoolHostToBizIdleModule: types.ResourcePoolHostTransferToBusiness,
		meta.MoveResPoolHostToDirectory:     types.ResourcePoolHostTransferToDirectory,
		meta.MoveBizHostFromModuleToResPool: types.BusinessHostTransferToResourcePool,
		meta.AddHostToResourcePool:          types.CreateResourcePoolHost,
		meta.Create:                         types.CreateResourcePoolHost,
		meta.Delete:                         types.DeleteResourcePoolHost,
		meta.MoveHostToAnotherBizModule:     types.HostTransferAcrossBusiness,
		meta.Find:                           types.ViewResourcePoolHost,
		meta.FindMany:                       types.ViewResourcePoolHost,
		meta.ManageHostAgentID:              types.ManageHostAgentID,
	},
	meta.ProcessServiceCategory: {
		meta.Delete: types.DeleteBusinessServiceCategory,
		meta.Update: types.EditBusinessServiceCategory,
		meta.Create: types.CreateBusinessServiceCategory,
		meta.Find:   types.Skip,
	},
	meta.ProcessServiceInstance: {
		meta.Delete:   types.DeleteBusinessServiceInstance,
		meta.Update:   types.EditBusinessServiceInstance,
		meta.Create:   types.CreateBusinessServiceInstance,
		meta.Find:     types.Skip,
		meta.FindMany: types.Skip,
	},
	meta.ProcessServiceTemplate: {
		meta.Delete:   types.DeleteBusinessServiceTemplate,
		meta.Update:   types.EditBusinessServiceTemplate,
		meta.Create:   types.CreateBusinessServiceTemplate,
		meta.Find:     types.Skip,
		meta.FindMany: types.Skip,
	},
	meta.SetTemplate: {
		meta.Delete:   types.DeleteBusinessSetTemplate,
		meta.Update:   types.EditBusinessSetTemplate,
		meta.Create:   types.CreateBusinessSetTemplate,
		meta.Find:     types.Skip,
		meta.FindMany: types.Skip,
	},
	meta.ModelModule: {
		meta.Delete:   types.DeleteBusinessTopology,
		meta.Update:   types.EditBusinessTopology,
		meta.Create:   types.CreateBusinessTopology,
		meta.Find:     types.Skip,
		meta.FindMany: types.Skip,
	},
	meta.ModelSet: {
		meta.Delete:   types.DeleteBusinessTopology,
		meta.Update:   types.EditBusinessTopology,
		meta.Create:   types.CreateBusinessTopology,
		meta.Find:     types.Skip,
		meta.FindMany: types.Skip,
	},
	meta.MainlineInstance: {
		meta.Delete:   types.DeleteBusinessTopology,
		meta.Update:   types.EditBusinessTopology,
		meta.Create:   types.CreateBusinessTopology,
		meta.Find:     types.Skip,
		meta.FindMany: types.Skip,
	},
	meta.MainlineInstanceTopology: {
		meta.Delete: types.Skip,
		meta.Update: types.Skip,
		meta.Create: types.Skip,
		meta.Find:   types.Skip,
	},
	meta.HostApply: {
		meta.Create:           types.EditBusinessHostApply,
		meta.Update:           types.EditBusinessHostApply,
		meta.Delete:           types.EditBusinessHostApply,
		meta.Find:             types.Skip,
		meta.DefaultHostApply: types.ViewBusinessResource,
	},
	meta.ResourcePoolDirectory: {
		meta.Delete:                types.DeleteResourcePoolDirectory,
		meta.Update:                types.EditResourcePoolDirectory,
		meta.Create:                types.CreateResourcePoolDirectory,
		meta.AddHostToResourcePool: types.CreateResourcePoolHost,
		meta.Find:                  types.Skip,
	},
	meta.CloudAreaInstance: {
		meta.Delete:   types.DeleteCloudArea,
		meta.Update:   types.EditCloudArea,
		meta.Create:   types.CreateCloudArea,
		meta.Find:     types.ViewCloudArea,
		meta.FindMany: types.ViewCloudArea,
	},
	meta.Model: {
		meta.Delete:   types.DeleteSysModel,
		meta.Update:   types.EditSysModel,
		meta.Create:   types.CreateSysModel,
		meta.Find:     types.ViewSysModel,
		meta.FindMany: types.ViewSysModel,
	},
	meta.AssociationType: {
		meta.Delete:   types.DeleteAssociationType,
		meta.Update:   types.EditAssociationType,
		meta.Create:   types.CreateAssociationType,
		meta.Find:     types.Skip,
		meta.FindMany: types.Skip,
	},
	meta.ModelClassification: {
		meta.Delete:   types.DeleteModelGroup,
		meta.Update:   types.EditModelGroup,
		meta.Create:   types.CreateModelGroup,
		meta.Find:     types.Skip,
		meta.FindMany: types.Skip,
	},
	meta.AuditLog: {
		meta.Find:     types.FindAuditLog,
		meta.FindMany: types.FindAuditLog,
	},
	meta.SystemBase: {
		meta.ModelTopologyView:      types.EditModelTopologyView,
		meta.ModelTopologyOperation: types.EditBusinessLayer,
	},
	meta.EventWatch: {
		meta.WatchHost:             types.WatchHostEvent,
		meta.WatchHostRelation:     types.WatchHostRelationEvent,
		meta.WatchBiz:              types.WatchBizEvent,
		meta.WatchSet:              types.WatchSetEvent,
		meta.WatchModule:           types.WatchModuleEvent,
		meta.WatchProcess:          types.WatchProcessEvent,
		meta.WatchCommonInstance:   types.WatchCommonInstanceEvent,
		meta.WatchMainlineInstance: types.WatchMainlineInstanceEvent,
		meta.WatchInstAsst:         types.WatchInstAsstEvent,
		meta.WatchBizSet:           types.WatchBizSetEvent,
		meta.WatchPlat:             types.WatchPlatEvent,
		meta.WatchKubeCluster:      types.WatchKubeClusterEvent,
		meta.WatchKubeNode:         types.WatchKubeNodeEvent,
		meta.WatchKubeNamespace:    types.WatchKubeNamespaceEvent,
		meta.WatchKubeWorkload:     types.WatchKubeWorkloadEvent,
		meta.WatchKubePod:          types.WatchKubePodEvent,
		meta.WatchProject:          types.WatchProjectEvent,
	},
	meta.UserCustom: {
		meta.Find:   types.Skip,
		meta.Update: types.Skip,
		meta.Delete: types.Skip,
		meta.Create: types.Skip,
	},
	meta.ModelAssociation: {
		meta.Find:     types.ViewSysModel,
		meta.FindMany: types.ViewSysModel,
		meta.Update:   types.EditSysModel,
		meta.Delete:   types.EditSysModel,
		meta.Create:   types.EditSysModel,
	},
	meta.ModelInstanceTopology: {
		meta.Find:   types.Skip,
		meta.Update: types.Skip,
		meta.Delete: types.Skip,
		meta.Create: types.Skip,
	},
	meta.ModelAttribute: {
		meta.Find:   types.ViewSysModel,
		meta.Update: types.EditSysModel,
		meta.Delete: types.DeleteSysModel,
		meta.Create: types.CreateSysModel,
	},
	meta.HostFavorite: {
		meta.Find:   types.Skip,
		meta.Update: types.Skip,
		meta.Delete: types.Skip,
		meta.Create: types.Skip,
	},

	meta.ProcessTemplate: {
		meta.Find:   types.Skip,
		meta.Delete: types.DeleteBusinessServiceTemplate,
		meta.Update: types.EditBusinessServiceTemplate,
		meta.Create: types.CreateBusinessServiceTemplate,
	},
	meta.BizTopology: {
		meta.Find:   types.Skip,
		meta.Update: types.EditBusinessTopology,
		meta.Delete: types.DeleteBusinessTopology,
		meta.Create: types.CreateBusinessTopology,
	},
	// unsupported resource actions for now
	meta.NetDataCollector: {
		meta.Find:   types.Unsupported,
		meta.Update: types.Unsupported,
		meta.Delete: types.Unsupported,
		meta.Create: types.Unsupported,
	},
	meta.InstallBK: {
		meta.Update: types.Skip,
	},
	// TODO: confirm this
	meta.SystemConfig: {
		meta.FindMany: types.Skip,
		meta.Find:     types.Skip,
		meta.Update:   types.Skip,
		meta.Delete:   types.Skip,
		meta.Create:   types.Skip,
	},
	meta.ConfigAdmin: {

		// reuse GlobalSettings permissions
		meta.Find:   types.Skip,
		meta.Update: types.GlobalSettings,
		meta.Delete: types.Unsupported,
		meta.Create: types.Unsupported,
	},
	meta.KubeCluster: {
		meta.Find:   types.ViewBusinessResource,
		meta.Update: types.EditContainerCluster,
		meta.Delete: types.DeleteContainerCluster,
		meta.Create: types.CreateContainerCluster,
	},
	meta.KubeNode: {
		meta.Find:   types.ViewBusinessResource,
		meta.Update: types.EditContainerNode,
		meta.Delete: types.DeleteContainerNode,
		meta.Create: types.CreateContainerNode,
	},
	meta.KubeNamespace: {
		meta.Find:   types.ViewBusinessResource,
		meta.Update: types.EditContainerNamespace,
		meta.Delete: types.DeleteContainerNamespace,
		meta.Create: types.CreateContainerNamespace,
	},
	meta.KubeWorkload: {
		meta.Find:   types.ViewBusinessResource,
		meta.Update: types.EditContainerWorkload,
		meta.Delete: types.DeleteContainerWorkload,
		meta.Create: types.CreateContainerWorkload,
	},
	meta.KubeDeployment: {
		meta.Find:   types.ViewBusinessResource,
		meta.Update: types.EditContainerWorkload,
		meta.Delete: types.DeleteContainerWorkload,
		meta.Create: types.CreateContainerWorkload,
	},
	meta.KubeStatefulSet: {
		meta.Find:   types.ViewBusinessResource,
		meta.Update: types.EditContainerWorkload,
		meta.Delete: types.DeleteContainerWorkload,
		meta.Create: types.CreateContainerWorkload,
	},
	meta.KubeDaemonSet: {
		meta.Find:   types.ViewBusinessResource,
		meta.Update: types.EditContainerWorkload,
		meta.Delete: types.DeleteContainerWorkload,
		meta.Create: types.CreateContainerWorkload,
	},
	meta.KubeGameStatefulSet: {
		meta.Find:   types.ViewBusinessResource,
		meta.Update: types.EditContainerWorkload,
		meta.Delete: types.DeleteContainerWorkload,
		meta.Create: types.CreateContainerWorkload,
	},
	meta.KubeGameDeployment: {
		meta.Find:   types.ViewBusinessResource,
		meta.Update: types.EditContainerWorkload,
		meta.Delete: types.DeleteContainerWorkload,
		meta.Create: types.CreateContainerWorkload,
	},
	meta.KubeCronJob: {
		meta.Find:   types.ViewBusinessResource,
		meta.Update: types.EditContainerWorkload,
		meta.Delete: types.DeleteContainerWorkload,
		meta.Create: types.CreateContainerWorkload,
	},
	meta.KubeJob: {
		meta.Find:   types.ViewBusinessResource,
		meta.Update: types.EditContainerWorkload,
		meta.Delete: types.DeleteContainerWorkload,
		meta.Create: types.CreateContainerWorkload,
	},
	meta.KubePodWorkload: {
		meta.Find:   types.ViewBusinessResource,
		meta.Update: types.EditContainerWorkload,
		meta.Delete: types.DeleteContainerWorkload,
		meta.Create: types.CreateContainerWorkload,
	},
	meta.KubePod: {
		meta.Find:   types.ViewBusinessResource,
		meta.Delete: types.DeleteContainerPod,
		meta.Create: types.CreateContainerPod,
	},
	meta.KubeContainer: {
		meta.Find: types.ViewBusinessResource,
	},
	meta.Project: {
		meta.Find:   types.ViewProject,
		meta.Update: types.EditProject,
		meta.Delete: types.DeleteProject,
		meta.Create: types.CreateProject,
	},
	meta.FulltextSearch: {
		meta.Find: types.UseFulltextSearch,
	},
	meta.FieldTemplate: {
		meta.Create: types.CreateFieldGroupingTemplate,
		meta.Find:   types.ViewFieldGroupingTemplate,
		meta.Update: types.EditFieldGroupingTemplate,
		meta.Delete: types.DeleteFieldGroupingTemplate,
	},
	meta.IDRuleIncrID: {
		meta.Update: types.EditIDRuleIncrID,
	},
	meta.FullSyncCond: {
		meta.Create: types.CreateFullSyncCond,
		meta.Find:   types.ViewFullSyncCond,
		meta.Update: types.EditFullSyncCond,
		meta.Delete: types.DeleteFullSyncCond,
	},
	meta.GeneralCache: {
		meta.Find: types.ViewGeneralCache,
	},
	meta.TenantSet: {
		meta.Find:            types.ViewTenantSet,
		meta.AccessTenantSet: types.AccessTenantSet,
	},
}

// ParseIamPathToAncestors TODO
func ParseIamPathToAncestors(iamPath []string) ([]metadata.IamResourceInstance, error) {
	instances := make([]metadata.IamResourceInstance, 0)
	for _, path := range iamPath {
		pathItemArr := strings.Split(strings.Trim(path, "/"), "/")
		for _, pathItem := range pathItemArr {
			typeAndID := strings.Split(pathItem, ",")
			if len(typeAndID) != 2 {
				return nil, fmt.Errorf("pathItem %s invalid", pathItem)
			}
			id := typeAndID[1]
			if id == "*" {
				continue
			}
			instances = append(instances, metadata.IamResourceInstance{
				Type:     typeAndID[0],
				TypeName: ResourceTypeIDMap[types.TypeID(typeAndID[0])],
				ID:       id,
			})
		}
	}
	return instances, nil
}

// GenIAMDynamicResTypeID 生成IAM侧资源的的dynamic resource typeID
func GenIAMDynamicResTypeID(modelID int64) types.TypeID {
	return types.TypeID(fmt.Sprintf("%s%d", types.IAMSysInstTypePrefix, modelID))
}

// GenCMDBDynamicResType 生成CMDB侧资源的的dynamic resourceType
func GenCMDBDynamicResType(modelID int64) meta.ResourceType {
	return meta.ResourceType(fmt.Sprintf("%s%d", meta.CMDBSysInstTypePrefix, modelID))
}

// genDynamicResourceType generate dynamic resourceType
func genDynamicResourceType(tenantID string, obj metadata.Object) iam.ResourceType {
	return iam.ResourceType{
		ID:      GenIAMDynamicResTypeID(obj.ID),
		Name:    obj.ObjectName,
		NameEn:  obj.ObjectID,
		Parents: nil,
		ProviderConfig: iam.ResourceConfig{
			Path: "/auth/v3/find/resource",
		},
		Version:  1,
		TenantID: tenantID,
	}
}

// genDynamicResourceTypes generate dynamic resourceTypes
func genDynamicResourceTypes(tenantObjects map[string][]metadata.Object) []iam.ResourceType {
	resourceTypes := make([]iam.ResourceType, 0)

	for tenantID, objects := range tenantObjects {
		for _, obj := range objects {
			resourceTypes = append(resourceTypes, genDynamicResourceType(tenantID, obj))
		}
	}

	return resourceTypes
}

// genIAMDynamicInstanceSelection generate IAM dynamic instanceSelection
func genIAMDynamicInstanceSelection(modelID int64) types.InstanceSelectionID {
	return types.InstanceSelectionID(fmt.Sprintf("%s%d", types.IAMSysInstTypePrefix, modelID))
}

// genDynamicInstanceSelection generate dynamic instanceSelection
func genDynamicInstanceSelection(tenantID string, obj metadata.Object) iam.InstanceSelection {
	return iam.InstanceSelection{
		ID:     genIAMDynamicInstanceSelection(obj.ID),
		Name:   obj.ObjectName,
		NameEn: obj.ObjectID,
		ResourceTypeChain: []iam.ResourceChain{{
			SystemID: types.SystemIDCMDB,
			ID:       GenIAMDynamicResTypeID(obj.ID),
		}},
		TenantID: tenantID,
	}
}

// genDynamicInstanceSelections generate dynamic instanceSelections
func genDynamicInstanceSelections(tenantObjects map[string][]metadata.Object) []iam.InstanceSelection {
	instanceSelections := make([]iam.InstanceSelection, 0)

	for tenantID, objects := range tenantObjects {
		for _, obj := range objects {
			instanceSelections = append(instanceSelections, genDynamicInstanceSelection(tenantID, obj))
		}
	}

	return instanceSelections
}

// genDynamicAction generate dynamic action
// Note: view action must be in the first place
func genDynamicAction(obj metadata.Object) []types.DynamicAction {
	return []types.DynamicAction{
		genDynamicViewAction(obj),
		genDynamicCreateAction(obj),
		genDynamicEditAction(obj),
		genDynamicDeleteAction(obj),
	}
}

// GenDynamicActionID generate dynamic ActionID
func GenDynamicActionID(actionType types.ActionType, modelID int64) types.ActionID {
	return types.ActionID(fmt.Sprintf("%s_%s%d", actionType, types.IAMSysInstTypePrefix, modelID))
}

// genDynamicViewAction generate dynamic view action
func genDynamicViewAction(obj metadata.Object) types.DynamicAction {
	return types.DynamicAction{
		ActionID:     GenDynamicActionID(types.View, obj.ID),
		ActionType:   types.View,
		ActionNameCN: fmt.Sprintf("%s%s%s", obj.ObjectName, "实例", "查看"),
		ActionNameEN: fmt.Sprintf("%s %s %s", "view", obj.ObjectID, "instance"),
	}
}

// genDynamicCreateAction generate dynamic create action
func genDynamicCreateAction(obj metadata.Object) types.DynamicAction {
	return types.DynamicAction{
		ActionID:     GenDynamicActionID(types.Create, obj.ID),
		ActionType:   types.Create,
		ActionNameCN: fmt.Sprintf("%s%s%s", obj.ObjectName, "实例", "新建"),
		ActionNameEN: fmt.Sprintf("%s %s %s", "create", obj.ObjectID, "instance"),
	}
}

// genDynamicEditAction generate dynamic edit action
func genDynamicEditAction(obj metadata.Object) types.DynamicAction {
	return types.DynamicAction{
		ActionID:     GenDynamicActionID(types.Edit, obj.ID),
		ActionType:   types.Edit,
		ActionNameCN: fmt.Sprintf("%s%s%s", obj.ObjectName, "实例", "编辑"),
		ActionNameEN: fmt.Sprintf("%s %s %s", "edit", obj.ObjectID, "instance"),
	}
}

// genDynamicDeleteAction generate dynamic delete action
func genDynamicDeleteAction(obj metadata.Object) types.DynamicAction {
	return types.DynamicAction{
		ActionID:     GenDynamicActionID(types.Delete, obj.ID),
		ActionType:   types.Delete,
		ActionNameCN: fmt.Sprintf("%s%s%s", obj.ObjectName, "实例", "删除"),
		ActionNameEN: fmt.Sprintf("%s %s %s", "delete", obj.ObjectID, "instance"),
	}
}

// genDynamicActionSubGroup 动态的按模型生成动作分组作为‘模型实例管理’分组的subGroup
func genDynamicActionSubGroup(obj metadata.Object) iam.ActionGroup {
	actions := genDynamicAction(obj)
	actionWithIDs := make([]iam.ActionWithID, len(actions))
	for idx, action := range actions {
		actionWithIDs[idx] = iam.ActionWithID{ID: action.ActionID}
	}
	return iam.ActionGroup{
		Name:    obj.ObjectName,
		NameEn:  obj.ObjectID,
		Actions: actionWithIDs,
	}
}

// genDynamicActionIDs generate dynamic model actionIDs
func genDynamicActionIDs(object metadata.Object) []types.ActionID {
	actions := genDynamicAction(object)
	actionIDs := make([]types.ActionID, len(actions))
	for idx, action := range actions {
		actionIDs[idx] = action.ActionID
	}
	return actionIDs
}

// genDynamicActions generate dynamic model actions
func genDynamicActions(tenantObjects map[string][]metadata.Object) []iam.ResourceAction {
	resActions := make([]iam.ResourceAction, 0)
	for tenantID, objects := range tenantObjects {
		for _, obj := range objects {
			relatedResource := []iam.RelateResourceType{
				{
					SystemID: types.SystemIDCMDB,
					ID:       GenIAMDynamicResTypeID(obj.ID),
					// 配置权限时可选择实例和配置属性, 后者用于属性鉴权
					SelectionMode: types.ModeAll,
					InstanceSelections: []iam.RelatedInstanceSelection{{
						SystemID: types.SystemIDCMDB,
						ID:       genIAMDynamicInstanceSelection(obj.ID),
					}},
				},
			}

			actions := genDynamicAction(obj)
			var relatedActions []types.ActionID
			for _, action := range actions {
				switch action.ActionType {
				case types.View:
					resActions = append(resActions, iam.ResourceAction{
						ID:       action.ActionID,
						Name:     action.ActionNameCN,
						NameEn:   action.ActionNameEN,
						Type:     types.View,
						Version:  1,
						TenantID: tenantID,
					})
					relatedActions = []types.ActionID{action.ActionID}

				case types.Create:
					resActions = append(resActions, iam.ResourceAction{
						ID:       action.ActionID,
						Name:     action.ActionNameCN,
						NameEn:   action.ActionNameEN,
						Type:     types.Create,
						Version:  1,
						TenantID: tenantID,
					})
				case types.Edit:
					resActions = append(resActions, iam.ResourceAction{
						ID:                   action.ActionID,
						Name:                 action.ActionNameCN,
						NameEn:               action.ActionNameEN,
						Type:                 types.Edit,
						RelatedActions:       relatedActions,
						Version:              1,
						RelatedResourceTypes: relatedResource,
						TenantID:             tenantID,
					})

				case types.Delete:
					resActions = append(resActions, iam.ResourceAction{
						ID:                   action.ActionID,
						Name:                 action.ActionNameCN,
						NameEn:               action.ActionNameEN,
						Type:                 types.Delete,
						RelatedResourceTypes: relatedResource,
						RelatedActions:       relatedActions,
						Version:              1,
						TenantID:             tenantID,
					})
				default:
					return nil
				}
			}
		}
	}

	return resActions
}

// IsIAMSysInstance judge whether the resource type is a system instance in iam resource
func IsIAMSysInstance(resourceType types.TypeID) bool {
	return strings.HasPrefix(string(resourceType), types.IAMSysInstTypePrefix)
}

// IsCMDBSysInstance judge whether the resource type is a system instance in cmdb resource
func IsCMDBSysInstance(resourceType meta.ResourceType) bool {
	return strings.HasPrefix(string(resourceType), meta.CMDBSysInstTypePrefix)
}

// isIAMSysInstanceSelection judge whether the instance selection is a system instance selection in iam resource
func isIAMSysInstanceSelection(instanceSelectionID types.InstanceSelectionID) bool {
	return strings.Contains(string(instanceSelectionID), types.IAMSysInstTypePrefix)
}

// isIAMSysInstanceAction judge whether the action is a system instance action in iam resource
func isIAMSysInstanceAction(actionID types.ActionID) bool {
	return strings.Contains(string(actionID), types.IAMSysInstTypePrefix)
}

// GetModelIDFromIamSysInstance get model id from iam system instance
func GetModelIDFromIamSysInstance(resourceType types.TypeID) (int64, error) {
	if !IsIAMSysInstance(resourceType) {
		return 0, fmt.Errorf("resourceType %s is not an iam system instance, it must start with prefix %s",
			resourceType, types.IAMSysInstTypePrefix)
	}
	modelIDStr := strings.TrimPrefix(string(resourceType), types.IAMSysInstTypePrefix)
	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		blog.ErrorJSON("modelID convert to int64 failed, err:%s, input:%s", err, modelID)
		return 0, fmt.Errorf("get model id failed, parse to int err:%s, the format of resourceType:%s is wrong",
			err.Error(), resourceType)
	}

	return modelID, nil
}

// GetActionTypeFromIAMSysInstance get action type from iam system instance
func GetActionTypeFromIAMSysInstance(actionID types.ActionID) types.ActionType {
	actionIDStr := string(actionID)
	return types.ActionType(actionIDStr[:strings.Index(actionIDStr, "_")])
}
