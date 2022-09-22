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
	"configcenter/cmd/scene_server/auth_server/sdk/types"
	meta2 "configcenter/pkg/ac/meta"
	"fmt"
	"strconv"
	"strings"

	"configcenter/pkg/blog"
	"configcenter/pkg/metadata"
)

// NotEnoughLayer TODO
var NotEnoughLayer = fmt.Errorf("not enough layer")

// AdaptAuthOptions TODO
func AdaptAuthOptions(a *meta2.ResourceAttribute) (ActionID, []types.Resource, error) {

	var action ActionID

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

// ConvertResourceType convert resource type from CMDB to IAM
func ConvertResourceType(resourceType meta2.ResourceType, businessID int64) (*TypeID, error) {
	var iamResourceType TypeID
	switch resourceType {
	case meta2.Business:
		iamResourceType = Business
	case meta2.BizSet:
		iamResourceType = BizSet
	case meta2.Model,
		meta2.ModelUnique,
		meta2.ModelAttributeGroup:
		iamResourceType = SysModel
	case meta2.ModelAttribute:
		if businessID > 0 {
			iamResourceType = BizCustomField
		} else {
			iamResourceType = SysModel
		}
	case meta2.ModelModule, meta2.ModelSet, meta2.MainlineInstance, meta2.MainlineInstanceTopology:
		iamResourceType = BizTopology
	case meta2.MainlineModel, meta2.ModelTopology:
	case meta2.ModelClassification:
		iamResourceType = SysModelGroup
	case meta2.AssociationType:
		iamResourceType = SysAssociationType
	case meta2.ModelAssociation:
		iamResourceType = SysModel
	case meta2.MainlineModelTopology:
	case meta2.ModelInstanceTopology:
		iamResourceType = SkipType
	case meta2.CloudAreaInstance:
		iamResourceType = SysCloudArea
	case meta2.HostInstance:
		iamResourceType = Host
	case meta2.HostFavorite:
		iamResourceType = SkipType
	case meta2.Process:
		iamResourceType = BizProcessServiceInstance
	case meta2.DynamicGrouping:
		iamResourceType = BizCustomQuery
	case meta2.AuditLog:
		iamResourceType = SysAuditLog
	case meta2.SystemBase:
	case meta2.UserCustom:
		iamResourceType = UserCustom
	case meta2.NetDataCollector:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	case meta2.ProcessServiceTemplate:
		iamResourceType = BizProcessServiceTemplate
	case meta2.ProcessServiceCategory:
		iamResourceType = BizProcessServiceCategory
	case meta2.ProcessServiceInstance:
		iamResourceType = BizProcessServiceInstance
	case meta2.BizTopology:
		iamResourceType = BizTopology
	case meta2.SetTemplate:
		iamResourceType = BizSetTemplate
	case meta2.OperationStatistic:
		iamResourceType = SysOperationStatistic
	case meta2.HostApply:
		iamResourceType = BizHostApply
	case meta2.ResourcePoolDirectory:
		iamResourceType = SysResourcePoolDirectory
	case meta2.CloudAccount:
		iamResourceType = SysCloudAccount
	case meta2.CloudResourceTask:
		iamResourceType = SysCloudResourceTask
	case meta2.EventWatch:
		iamResourceType = SysEventWatch
	case meta2.ConfigAdmin:
	case meta2.SystemConfig:
	case meta2.KubeCluster, meta2.KubeNode, meta2.KubeNamespace, meta2.KubeWorkload, meta2.KubeDeployment,
		meta2.KubeStatefulSet, meta2.KubeDaemonSet, meta2.KubeGameStatefulSet, meta2.KubeGameDeployment, meta2.KubeCronJob,
		meta2.KubeJob, meta2.KubePodWorkload, meta2.KubePod, meta2.KubeContainer:
	default:
		if IsCMDBSysInstance(resourceType) {
			iamResourceType = TypeID(resourceType)
			return &iamResourceType, nil
		}

		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	return &iamResourceType, nil
}

// ConvertResourceAction convert resource action from CMDB to IAM
func ConvertResourceAction(resourceType meta2.ResourceType, action meta2.Action, businessID int64) (ActionID, error) {
	if action == meta2.SkipAction {
		return Skip, nil
	}

	convertAction := action
	switch action {
	case meta2.CreateMany:
		convertAction = meta2.Create
	case meta2.FindMany:
		convertAction = meta2.Find
	case meta2.DeleteMany:
		convertAction = meta2.Delete
	case meta2.UpdateMany:
		convertAction = meta2.Update
	}

	if resourceType == meta2.ModelAttribute || resourceType == meta2.ModelAttributeGroup {
		if convertAction == meta2.Delete || convertAction == meta2.Update || convertAction == meta2.Create {
			if businessID > 0 {
				return EditBusinessCustomField, nil
			} else {
				return EditSysModel, nil
			}
		}
	}

	if resourceType == meta2.HostInstance && convertAction == meta2.Update {
		if businessID > 0 {
			return EditBusinessHost, nil
		} else {
			return EditResourcePoolHost, nil
		}
	}

	if IsCMDBSysInstance(resourceType) {
		return ConvertSysInstanceActionID(resourceType, convertAction)
	}

	if _, exist := resourceActionMap[resourceType]; exist {
		actionID, ok := resourceActionMap[resourceType][convertAction]
		if ok && actionID != Unsupported {
			return actionID, nil
		}
	}

	return Unsupported, fmt.Errorf("unsupported type %s action: %s", resourceType, action)
}

// ConvertSysInstanceActionID convert system instances action from CMDB to IAM
func ConvertSysInstanceActionID(resourceType meta2.ResourceType, action meta2.Action) (ActionID, error) {
	var actionType ActionType
	switch action {
	case meta2.Create:
		actionType = Create
	case meta2.Update:
		actionType = Edit
	case meta2.Delete:
		actionType = Delete
	case meta2.Find:
		return Skip, nil
	default:
		return Unsupported, fmt.Errorf("unsupported action: %s", action)
	}
	id := strings.TrimPrefix(string(resourceType), meta2.CMDBSysInstTypePrefix)
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		return Unsupported, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
	return ActionID(fmt.Sprintf("%s_%s%s", actionType, IAMSysInstTypePrefix, id)), nil
}

var resourceActionMap = map[meta2.ResourceType]map[meta2.Action]ActionID{
	meta2.ModelAttributeGroup: {
		meta2.Delete:   EditSysModel,
		meta2.Update:   EditSysModel,
		meta2.Create:   EditSysModel,
		meta2.Find:     Skip,
		meta2.FindMany: Skip,
	},
	meta2.ModelUnique: {
		meta2.Delete:   EditSysModel,
		meta2.Update:   EditSysModel,
		meta2.Create:   EditSysModel,
		meta2.Find:     Skip,
		meta2.FindMany: Skip,
	},
	meta2.Business: {
		meta2.Archive:              ArchiveBusiness,
		meta2.Create:               CreateBusiness,
		meta2.Update:               EditBusiness,
		meta2.Find:                 FindBusiness,
		meta2.ViewBusinessResource: ViewBusinessResource,
	},
	meta2.BizSet: {
		meta2.Create:       CreateBizSet,
		meta2.Update:       EditBizSet,
		meta2.Delete:       DeleteBizSet,
		meta2.Find:         ViewBizSet,
		meta2.AccessBizSet: AccessBizSet,
	},
	meta2.DynamicGrouping: {
		meta2.Delete:   DeleteBusinessCustomQuery,
		meta2.Update:   EditBusinessCustomQuery,
		meta2.Create:   CreateBusinessCustomQuery,
		meta2.Find:     ViewBusinessResource,
		meta2.FindMany: ViewBusinessResource,
		meta2.Execute:  ViewBusinessResource,
	},
	meta2.MainlineModel: {
		meta2.Find:   Skip,
		meta2.Create: EditBusinessLayer,
		meta2.Delete: EditBusinessLayer,
	},
	meta2.ModelTopology: {
		meta2.Find:   EditModelTopologyView,
		meta2.Update: EditModelTopologyView,
	},
	meta2.MainlineModelTopology: {
		meta2.Find: Skip,
	},
	meta2.Process: {
		meta2.Find:   Skip,
		meta2.Create: EditBusinessServiceInstance,
		meta2.Delete: EditBusinessServiceInstance,
		meta2.Update: EditBusinessServiceInstance,
	},
	meta2.HostInstance: {
		meta2.MoveResPoolHostToBizIdleModule: ResourcePoolHostTransferToBusiness,
		meta2.MoveResPoolHostToDirectory:     ResourcePoolHostTransferToDirectory,
		meta2.MoveBizHostFromModuleToResPool: BusinessHostTransferToResourcePool,
		meta2.AddHostToResourcePool:          CreateResourcePoolHost,
		meta2.Create:                         CreateResourcePoolHost,
		meta2.Delete:                         DeleteResourcePoolHost,
		meta2.MoveHostToAnotherBizModule:     HostTransferAcrossBusiness,
		meta2.Find:                           Skip,
		meta2.FindMany:                       Skip,
	},
	meta2.ProcessServiceCategory: {
		meta2.Delete: DeleteBusinessServiceCategory,
		meta2.Update: EditBusinessServiceCategory,
		meta2.Create: CreateBusinessServiceCategory,
		meta2.Find:   Skip,
	},
	meta2.ProcessServiceInstance: {
		meta2.Delete:   DeleteBusinessServiceInstance,
		meta2.Update:   EditBusinessServiceInstance,
		meta2.Create:   CreateBusinessServiceInstance,
		meta2.Find:     Skip,
		meta2.FindMany: Skip,
	},
	meta2.ProcessServiceTemplate: {
		meta2.Delete:   DeleteBusinessServiceTemplate,
		meta2.Update:   EditBusinessServiceTemplate,
		meta2.Create:   CreateBusinessServiceTemplate,
		meta2.Find:     Skip,
		meta2.FindMany: Skip,
	},
	meta2.SetTemplate: {
		meta2.Delete:   DeleteBusinessSetTemplate,
		meta2.Update:   EditBusinessSetTemplate,
		meta2.Create:   CreateBusinessSetTemplate,
		meta2.Find:     Skip,
		meta2.FindMany: Skip,
	},
	meta2.ModelModule: {
		meta2.Delete:   DeleteBusinessTopology,
		meta2.Update:   EditBusinessTopology,
		meta2.Create:   CreateBusinessTopology,
		meta2.Find:     Skip,
		meta2.FindMany: Skip,
	},
	meta2.ModelSet: {
		meta2.Delete:   DeleteBusinessTopology,
		meta2.Update:   EditBusinessTopology,
		meta2.Create:   CreateBusinessTopology,
		meta2.Find:     Skip,
		meta2.FindMany: Skip,
	},
	meta2.MainlineInstance: {
		meta2.Delete:   DeleteBusinessTopology,
		meta2.Update:   EditBusinessTopology,
		meta2.Create:   CreateBusinessTopology,
		meta2.Find:     Skip,
		meta2.FindMany: Skip,
	},
	meta2.MainlineInstanceTopology: {
		meta2.Delete: Skip,
		meta2.Update: Skip,
		meta2.Create: Skip,
		meta2.Find:   Skip,
	},
	meta2.HostApply: {
		meta2.Create:           EditBusinessHostApply,
		meta2.Update:           EditBusinessHostApply,
		meta2.Delete:           EditBusinessHostApply,
		meta2.Find:             Skip,
		meta2.DefaultHostApply: ViewBusinessResource,
	},
	meta2.ResourcePoolDirectory: {
		meta2.Delete:                DeleteResourcePoolDirectory,
		meta2.Update:                EditResourcePoolDirectory,
		meta2.Create:                CreateResourcePoolDirectory,
		meta2.AddHostToResourcePool: CreateResourcePoolHost,
		meta2.Find:                  Skip,
	},
	meta2.CloudAreaInstance: {
		meta2.Delete:   DeleteCloudArea,
		meta2.Update:   EditCloudArea,
		meta2.Create:   CreateCloudArea,
		meta2.Find:     Skip,
		meta2.FindMany: Skip,
	},
	meta2.CloudAccount: {
		meta2.Delete:   DeleteCloudAccount,
		meta2.Update:   EditCloudAccount,
		meta2.Create:   CreateCloudAccount,
		meta2.Find:     FindCloudAccount,
		meta2.FindMany: FindCloudAccount,
	},
	meta2.CloudResourceTask: {
		meta2.Delete: DeleteCloudResourceTask,
		meta2.Update: EditCloudResourceTask,
		meta2.Create: CreateCloudResourceTask,
		meta2.Find:   FindCloudResourceTask,
	},
	meta2.Model: {
		meta2.Delete:   DeleteSysModel,
		meta2.Update:   EditSysModel,
		meta2.Create:   CreateSysModel,
		meta2.Find:     Skip,
		meta2.FindMany: Skip,
	},
	meta2.AssociationType: {
		meta2.Delete:   DeleteAssociationType,
		meta2.Update:   EditAssociationType,
		meta2.Create:   CreateAssociationType,
		meta2.Find:     Skip,
		meta2.FindMany: Skip,
	},
	meta2.ModelClassification: {
		meta2.Delete:   DeleteModelGroup,
		meta2.Update:   EditModelGroup,
		meta2.Create:   CreateModelGroup,
		meta2.Find:     Skip,
		meta2.FindMany: Skip,
	},
	meta2.OperationStatistic: {
		meta2.Create:   EditOperationStatistic,
		meta2.Delete:   EditOperationStatistic,
		meta2.Update:   EditOperationStatistic,
		meta2.Find:     FindOperationStatistic,
		meta2.FindMany: FindOperationStatistic,
	},
	meta2.AuditLog: {
		meta2.Find:     FindAuditLog,
		meta2.FindMany: FindAuditLog,
	},
	meta2.SystemBase: {
		meta2.ModelTopologyView:      EditModelTopologyView,
		meta2.ModelTopologyOperation: EditBusinessLayer,
	},
	meta2.EventWatch: {
		meta2.WatchHost:             WatchHostEvent,
		meta2.WatchHostRelation:     WatchHostRelationEvent,
		meta2.WatchBiz:              WatchBizEvent,
		meta2.WatchSet:              WatchSetEvent,
		meta2.WatchModule:           WatchModuleEvent,
		meta2.WatchProcess:          WatchProcessEvent,
		meta2.WatchCommonInstance:   WatchCommonInstanceEvent,
		meta2.WatchMainlineInstance: WatchMainlineInstanceEvent,
		meta2.WatchInstAsst:         WatchInstAsstEvent,
		meta2.WatchBizSet:           WatchBizSetEvent,
		meta2.WatchKubeCluster:      WatchKubeClusterEvent,
		meta2.WatchKubeNode:         WatchKubeNodeEvent,
		meta2.WatchKubeNamespace:    WatchKubeNamespaceEvent,
		meta2.WatchKubeWorkload:     WatchKubeWorkloadEvent,
		meta2.WatchKubePod:          WatchKubePodEvent,
	},
	meta2.UserCustom: {
		meta2.Find:   Skip,
		meta2.Update: Skip,
		meta2.Delete: Skip,
		meta2.Create: Skip,
	},
	meta2.ModelAssociation: {
		meta2.Find:     Skip,
		meta2.FindMany: Skip,
		meta2.Update:   EditSysModel,
		meta2.Delete:   EditSysModel,
		meta2.Create:   EditSysModel,
	},
	meta2.ModelInstanceTopology: {
		meta2.Find:   Skip,
		meta2.Update: Skip,
		meta2.Delete: Skip,
		meta2.Create: Skip,
	},
	meta2.ModelAttribute: {
		meta2.Find:   Skip,
		meta2.Update: EditSysModel,
		meta2.Delete: DeleteSysModel,
		meta2.Create: CreateSysModel,
	},
	meta2.HostFavorite: {
		meta2.Find:   Skip,
		meta2.Update: Skip,
		meta2.Delete: Skip,
		meta2.Create: Skip,
	},

	meta2.ProcessTemplate: {
		meta2.Find:   Skip,
		meta2.Delete: DeleteBusinessServiceTemplate,
		meta2.Update: EditBusinessServiceTemplate,
		meta2.Create: CreateBusinessServiceTemplate,
	},
	meta2.BizTopology: {
		meta2.Find:   Skip,
		meta2.Update: EditBusinessTopology,
		meta2.Delete: DeleteBusinessTopology,
		meta2.Create: CreateBusinessTopology,
	},
	// unsupported resource actions for now
	meta2.NetDataCollector: {
		meta2.Find:   Unsupported,
		meta2.Update: Unsupported,
		meta2.Delete: Unsupported,
		meta2.Create: Unsupported,
	},
	meta2.InstallBK: {
		meta2.Update: Skip,
	},
	// TODO: confirm this
	meta2.SystemConfig: {
		meta2.FindMany: Skip,
		meta2.Find:     Skip,
		meta2.Update:   Skip,
		meta2.Delete:   Skip,
		meta2.Create:   Skip,
	},
	meta2.ConfigAdmin: {

		// reuse GlobalSettings permissions
		meta2.Find:   Skip,
		meta2.Update: GlobalSettings,
		meta2.Delete: Unsupported,
		meta2.Create: Unsupported,
	},
	meta2.KubeCluster: {
		meta2.Find:   Skip,
		meta2.Update: EditContainerCluster,
		meta2.Delete: DeleteContainerCluster,
		meta2.Create: CreateContainerCluster,
	},
	meta2.KubeNode: {
		meta2.Find:   Skip,
		meta2.Update: EditContainerNode,
		meta2.Delete: DeleteContainerNode,
		meta2.Create: CreateContainerNode,
	},
	meta2.KubeNamespace: {
		meta2.Find:   Skip,
		meta2.Update: EditContainerNamespace,
		meta2.Delete: DeleteContainerNamespace,
		meta2.Create: CreateContainerNamespace,
	},
	meta2.KubeWorkload: {
		meta2.Find:   Skip,
		meta2.Update: EditContainerWorkload,
		meta2.Delete: DeleteContainerWorkload,
		meta2.Create: CreateContainerWorkload,
	},
	meta2.KubeDeployment: {
		meta2.Find:   Skip,
		meta2.Update: EditContainerWorkload,
		meta2.Delete: DeleteContainerWorkload,
		meta2.Create: CreateContainerWorkload,
	},
	meta2.KubeStatefulSet: {
		meta2.Find:   Skip,
		meta2.Update: EditContainerWorkload,
		meta2.Delete: DeleteContainerWorkload,
		meta2.Create: CreateContainerWorkload,
	},
	meta2.KubeDaemonSet: {
		meta2.Find:   Skip,
		meta2.Update: EditContainerWorkload,
		meta2.Delete: DeleteContainerWorkload,
		meta2.Create: CreateContainerWorkload,
	},
	meta2.KubeGameStatefulSet: {
		meta2.Find:   Skip,
		meta2.Update: EditContainerWorkload,
		meta2.Delete: DeleteContainerWorkload,
		meta2.Create: CreateContainerWorkload,
	},
	meta2.KubeGameDeployment: {
		meta2.Find:   Skip,
		meta2.Update: EditContainerWorkload,
		meta2.Delete: DeleteContainerWorkload,
		meta2.Create: CreateContainerWorkload,
	},
	meta2.KubeCronJob: {
		meta2.Find:   Skip,
		meta2.Update: EditContainerWorkload,
		meta2.Delete: DeleteContainerWorkload,
		meta2.Create: CreateContainerWorkload,
	},
	meta2.KubeJob: {
		meta2.Find:   Skip,
		meta2.Update: EditContainerWorkload,
		meta2.Delete: DeleteContainerWorkload,
		meta2.Create: CreateContainerWorkload,
	},
	meta2.KubePodWorkload: {
		meta2.Find:   Skip,
		meta2.Update: EditContainerWorkload,
		meta2.Delete: DeleteContainerWorkload,
		meta2.Create: CreateContainerWorkload,
	},
	meta2.KubePod: {
		meta2.Find:   Skip,
		meta2.Delete: DeleteContainerPod,
		meta2.Create: CreateContainerPod,
	},
	meta2.KubeContainer: {
		meta2.Find: Skip,
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
				TypeName: ResourceTypeIDMap[TypeID(typeAndID[0])],
				ID:       id,
			})
		}
	}
	return instances, nil
}

// GenIAMDynamicResTypeID 生成IAM侧资源的的dynamic resource typeID
func GenIAMDynamicResTypeID(modelID int64) TypeID {
	return TypeID(fmt.Sprintf("%s%d", IAMSysInstTypePrefix, modelID))
}

// GenCMDBDynamicResType 生成CMDB侧资源的的dynamic resourceType
func GenCMDBDynamicResType(modelID int64) meta2.ResourceType {
	return meta2.ResourceType(fmt.Sprintf("%s%d", meta2.CMDBSysInstTypePrefix, modelID))
}

// genDynamicResourceType generate dynamic resourceType
func genDynamicResourceType(obj metadata.Object) ResourceType {
	return ResourceType{
		ID:      GenIAMDynamicResTypeID(obj.ID),
		Name:    obj.ObjectName,
		NameEn:  obj.ObjectID,
		Parents: nil,
		ProviderConfig: ResourceConfig{
			Path: "/auth/v3/find/resource",
		},
		Version: 1,
	}
}

// genDynamicResourceTypes generate dynamic resourceTypes
func genDynamicResourceTypes(objects []metadata.Object) []ResourceType {
	resourceTypes := make([]ResourceType, 0)
	for _, obj := range objects {
		resourceTypes = append(resourceTypes, genDynamicResourceType(obj))
	}
	return resourceTypes
}

// genIAMDynamicInstanceSelection generate IAM dynamic instanceSelection
func genIAMDynamicInstanceSelection(modelID int64) InstanceSelectionID {
	return InstanceSelectionID(fmt.Sprintf("%s%d", IAMSysInstTypePrefix, modelID))
}

// genDynamicInstanceSelection generate dynamic instanceSelection
func genDynamicInstanceSelection(obj metadata.Object) InstanceSelection {
	return InstanceSelection{
		ID:     genIAMDynamicInstanceSelection(obj.ID),
		Name:   obj.ObjectName,
		NameEn: obj.ObjectID,
		ResourceTypeChain: []ResourceChain{{
			SystemID: SystemIDCMDB,
			ID:       GenIAMDynamicResTypeID(obj.ID),
		}},
	}
}

// genDynamicInstanceSelections generate dynamic instanceSelections
func genDynamicInstanceSelections(objects []metadata.Object) []InstanceSelection {
	instanceSelections := make([]InstanceSelection, 0)
	for _, obj := range objects {
		instanceSelections = append(instanceSelections, genDynamicInstanceSelection(obj))
	}
	return instanceSelections
}

// genDynamicAction generate dynamic action
func genDynamicAction(obj metadata.Object) []DynamicAction {
	return []DynamicAction{
		genDynamicCreateAction(obj),
		genDynamicEditAction(obj),
		genDynamicDeleteAction(obj),
	}
}

// GenDynamicActionID generate dynamic ActionID
func GenDynamicActionID(actionType ActionType, modelID int64) ActionID {
	return ActionID(fmt.Sprintf("%s_%s%d", actionType, IAMSysInstTypePrefix, modelID))
}

// genDynamicCreateAction generate dynamic create action
func genDynamicCreateAction(obj metadata.Object) DynamicAction {
	return DynamicAction{
		ActionID:     GenDynamicActionID(Create, obj.ID),
		ActionType:   Create,
		ActionNameCN: fmt.Sprintf("%s%s%s", obj.ObjectName, "实例", "新建"),
		ActionNameEN: fmt.Sprintf("%s %s %s", "create", obj.ObjectID, "instance"),
	}
}

// genDynamicEditAction generate dynamic edit action
func genDynamicEditAction(obj metadata.Object) DynamicAction {
	return DynamicAction{
		ActionID:     GenDynamicActionID(Edit, obj.ID),
		ActionType:   Edit,
		ActionNameCN: fmt.Sprintf("%s%s%s", obj.ObjectName, "实例", "编辑"),
		ActionNameEN: fmt.Sprintf("%s %s %s", "edit", obj.ObjectID, "instance"),
	}
}

// genDynamicDeleteAction generate dynamic delete action
func genDynamicDeleteAction(obj metadata.Object) DynamicAction {
	return DynamicAction{
		ActionID:     GenDynamicActionID(Delete, obj.ID),
		ActionType:   Delete,
		ActionNameCN: fmt.Sprintf("%s%s%s", obj.ObjectName, "实例", "删除"),
		ActionNameEN: fmt.Sprintf("%s %s %s", "delete", obj.ObjectID, "instance"),
	}
}

// genDynamicActionSubGroup 动态的按模型生成动作分组作为‘模型实例管理’分组的subGroup
func genDynamicActionSubGroup(obj metadata.Object) ActionGroup {
	actions := genDynamicAction(obj)
	actionWithIDs := make([]ActionWithID, len(actions))
	for idx, action := range actions {
		actionWithIDs[idx] = ActionWithID{ID: action.ActionID}
	}
	return ActionGroup{
		Name:    obj.ObjectName,
		NameEn:  obj.ObjectID,
		Actions: actionWithIDs,
	}
}

// genDynamicActionIDs generate dynamic model actionIDs
func genDynamicActionIDs(object metadata.Object) []ActionID {
	actions := genDynamicAction(object)
	actionIDs := make([]ActionID, len(actions))
	for idx, action := range actions {
		actionIDs[idx] = action.ActionID
	}
	return actionIDs
}

// genDynamicActions generate dynamic model actions
func genDynamicActions(objects []metadata.Object) []ResourceAction {
	resActions := make([]ResourceAction, 0)
	for _, obj := range objects {
		relatedResource := []RelateResourceType{
			{
				SystemID:    SystemIDCMDB,
				ID:          GenIAMDynamicResTypeID(obj.ID),
				NameAlias:   "",
				NameAliasEn: "",
				Scope:       nil,
				// 配置权限时可选择实例和配置属性, 后者用于属性鉴权
				SelectionMode: modeAll,
				InstanceSelections: []RelatedInstanceSelection{{
					SystemID: SystemIDCMDB,
					ID:       genIAMDynamicInstanceSelection(obj.ID),
				}},
			},
		}

		actions := genDynamicAction(obj)
		for _, action := range actions {
			switch action.ActionType {
			case Create:
				resActions = append(resActions, ResourceAction{
					ID:                   action.ActionID,
					Name:                 action.ActionNameCN,
					NameEn:               action.ActionNameEN,
					Type:                 Create,
					RelatedResourceTypes: nil,
					RelatedActions:       nil,
					Version:              1,
				})
			case Edit:
				resActions = append(resActions, ResourceAction{
					ID:                   action.ActionID,
					Name:                 action.ActionNameCN,
					NameEn:               action.ActionNameEN,
					Type:                 Edit,
					RelatedActions:       nil,
					Version:              1,
					RelatedResourceTypes: relatedResource,
				})

			case Delete:
				resActions = append(resActions, ResourceAction{
					ID:                   action.ActionID,
					Name:                 action.ActionNameCN,
					NameEn:               action.ActionNameEN,
					Type:                 Delete,
					RelatedResourceTypes: relatedResource,
					RelatedActions:       nil,
					Version:              1,
				})
			default:
				return nil
			}
		}
	}

	return resActions
}

// IsIAMSysInstance judge whether the resource type is a system instance in iam resource
func IsIAMSysInstance(resourceType TypeID) bool {
	return strings.HasPrefix(string(resourceType), IAMSysInstTypePrefix)
}

// IsCMDBSysInstance judge whether the resource type is a system instance in cmdb resource
func IsCMDBSysInstance(resourceType meta2.ResourceType) bool {
	return strings.HasPrefix(string(resourceType), meta2.CMDBSysInstTypePrefix)
}

// isIAMSysInstanceSelection judge whether the instance selection is a system instance selection in iam resource
func isIAMSysInstanceSelection(instanceSelectionID InstanceSelectionID) bool {
	return strings.Contains(string(instanceSelectionID), IAMSysInstTypePrefix)
}

// isIAMSysInstanceAction judge whether the action is a system instance action in iam resource
func isIAMSysInstanceAction(actionID ActionID) bool {
	return strings.Contains(string(actionID), IAMSysInstTypePrefix)
}

// GetModelIDFromIamSysInstance get model id from iam system instance
func GetModelIDFromIamSysInstance(resourceType TypeID) (int64, error) {
	if !IsIAMSysInstance(resourceType) {
		return 0, fmt.Errorf("resourceType %s is not an iam system instance, it must start with prefix %s",
			resourceType, IAMSysInstTypePrefix)
	}
	modelIDStr := strings.TrimPrefix(string(resourceType), IAMSysInstTypePrefix)
	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		blog.ErrorJSON("modelID convert to int64 failed, err:%s, input:%s", err, modelID)
		return 0, fmt.Errorf("get model id failed, parse to int err:%s, the format of resourceType:%s is wrong",
			err.Error(), resourceType)
	}

	return modelID, nil
}

// GetActionTypeFromIAMSysInstance get action type from iam system instance
func GetActionTypeFromIAMSysInstance(actionID ActionID) ActionType {
	actionIDStr := string(actionID)
	return ActionType(actionIDStr[:strings.Index(actionIDStr, "_")])
}
