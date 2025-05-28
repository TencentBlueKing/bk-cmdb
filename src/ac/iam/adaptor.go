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

	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/ac/meta"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/apigw/iam"
)

// NotEnoughLayer TODO
var NotEnoughLayer = fmt.Errorf("not enough layer")

// AdaptAuthOptions TODO
func AdaptAuthOptions(a *meta.ResourceAttribute) (iamtypes.ActionID, []iam.Resource, error) {

	var action iamtypes.ActionID

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

var ccIamResTypeMap = map[meta.ResourceType]iamtypes.TypeID{
	meta.Business:                 iamtypes.Business,
	meta.BizSet:                   iamtypes.BizSet,
	meta.Project:                  iamtypes.Project,
	meta.Model:                    iamtypes.SysModel,
	meta.ModelUnique:              iamtypes.SysModel,
	meta.ModelAttributeGroup:      iamtypes.SysModel,
	meta.ModelModule:              iamtypes.BizTopology,
	meta.ModelSet:                 iamtypes.BizTopology,
	meta.MainlineInstance:         iamtypes.BizTopology,
	meta.MainlineInstanceTopology: iamtypes.BizTopology,
	meta.MainlineModel:            iamtypes.TypeID(""),
	meta.ModelTopology:            iamtypes.TypeID(""),
	meta.ModelClassification:      iamtypes.SysModelGroup,
	meta.AssociationType:          iamtypes.SysAssociationType,
	meta.ModelAssociation:         iamtypes.SysModel,
	meta.MainlineModelTopology:    iamtypes.TypeID(""),
	meta.ModelInstanceTopology:    iamtypes.SkipType,
	meta.CloudAreaInstance:        iamtypes.SysCloudArea,
	meta.HostInstance:             iamtypes.Host,
	meta.HostFavorite:             iamtypes.SkipType,
	meta.Process:                  iamtypes.BizProcessServiceInstance,
	meta.DynamicGrouping:          iamtypes.BizCustomQuery,
	meta.AuditLog:                 iamtypes.SysAuditLog,
	meta.SystemBase:               iamtypes.TypeID(""),
	meta.UserCustom:               iamtypes.UserCustom,
	meta.ProcessServiceTemplate:   iamtypes.BizProcessServiceTemplate,
	meta.ProcessServiceCategory:   iamtypes.BizProcessServiceCategory,
	meta.ProcessServiceInstance:   iamtypes.BizProcessServiceInstance,
	meta.BizTopology:              iamtypes.BizTopology,
	meta.SetTemplate:              iamtypes.BizSetTemplate,
	meta.HostApply:                iamtypes.BizHostApply,
	meta.ResourcePoolDirectory:    iamtypes.SysResourcePoolDirectory,
	meta.EventWatch:               iamtypes.SysEventWatch,
	meta.ConfigAdmin:              iamtypes.TypeID(""),
	meta.SystemConfig:             iamtypes.TypeID(""),
	meta.KubeCluster:              iamtypes.TypeID(""),
	meta.KubeNode:                 iamtypes.TypeID(""),
	meta.KubeNamespace:            iamtypes.TypeID(""),
	meta.KubeWorkload:             iamtypes.TypeID(""),
	meta.KubeDeployment:           iamtypes.TypeID(""),
	meta.KubeStatefulSet:          iamtypes.TypeID(""),
	meta.KubeDaemonSet:            iamtypes.TypeID(""),
	meta.KubeGameStatefulSet:      iamtypes.TypeID(""),
	meta.KubeGameDeployment:       iamtypes.TypeID(""),
	meta.KubeCronJob:              iamtypes.TypeID(""),
	meta.KubeJob:                  iamtypes.TypeID(""),
	meta.KubePodWorkload:          iamtypes.TypeID(""),
	meta.KubePod:                  iamtypes.TypeID(""),
	meta.KubeContainer:            iamtypes.TypeID(""),
	meta.FieldTemplate:            iamtypes.FieldGroupingTemplate,
	meta.FulltextSearch:           iamtypes.TypeID(""),
	meta.IDRuleIncrID:             iamtypes.TypeID(""),
	meta.FullSyncCond:             iamtypes.TypeID(""),
	meta.GeneralCache:             iamtypes.GeneralCache,
	meta.TenantSet:                iamtypes.TenantSet,
}

// ConvertResourceType convert resource type from CMDB to IAM
func ConvertResourceType(resourceType meta.ResourceType, businessID int64) (*iamtypes.TypeID, error) {
	var iamResourceType iamtypes.TypeID

	switch resourceType {
	case meta.ModelAttribute:
		if businessID > 0 {
			iamResourceType = iamtypes.BizCustomField
		} else {
			iamResourceType = iamtypes.SysModel
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
		iamResourceType = iamtypes.TypeID(resourceType)
		return &iamResourceType, nil
	}

	return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
}

// ConvertResourceAction convert resource action from CMDB to IAM
func ConvertResourceAction(resourceType meta.ResourceType, action meta.Action, businessID int64) (iamtypes.ActionID,
	error) {
	if action == meta.SkipAction {
		return iamtypes.Skip, nil
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
				return iamtypes.EditBusinessCustomField, nil
			} else {
				return iamtypes.EditSysModel, nil
			}
		}
	case meta.HostInstance:
		switch convertAction {
		case meta.Update:
			if businessID > 0 {
				return iamtypes.EditBusinessHost, nil
			} else {
				return iamtypes.EditResourcePoolHost, nil
			}
		case meta.Find:
			if businessID > 0 {
				return iamtypes.ViewBusinessResource, nil
			} else {
				return iamtypes.ViewResourcePoolHost, nil
			}
		}
	}

	if IsCMDBSysInstance(resourceType) {
		return ConvertSysInstanceActionID(resourceType, convertAction)
	}

	if _, exist := resourceActionMap[resourceType]; exist {
		actionID, ok := resourceActionMap[resourceType][convertAction]
		if ok && actionID != iamtypes.Unsupported {
			return actionID, nil
		}
	}

	return iamtypes.Unsupported, fmt.Errorf("unsupported type %s action: %s", resourceType, action)
}

// ConvertSysInstanceActionID convert system instances action from CMDB to IAM
func ConvertSysInstanceActionID(resourceType meta.ResourceType, action meta.Action) (iamtypes.ActionID, error) {
	var actionType iamtypes.ActionType
	switch action {
	case meta.Create:
		actionType = iamtypes.Create
	case meta.Update:
		actionType = iamtypes.Edit
	case meta.Delete:
		actionType = iamtypes.Delete
	case meta.Find:
		actionType = iamtypes.View
	default:
		return iamtypes.Unsupported, fmt.Errorf("unsupported action: %s", action)
	}
	id := strings.TrimPrefix(string(resourceType), meta.CMDBSysInstTypePrefix)
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		return iamtypes.Unsupported, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
	return iamtypes.ActionID(fmt.Sprintf("%s_%s%s", actionType, iamtypes.IAMSysInstTypePrefix, id)), nil
}

var resourceActionMap = map[meta.ResourceType]map[meta.Action]iamtypes.ActionID{
	meta.ModelAttributeGroup: {
		meta.Delete:   iamtypes.EditSysModel,
		meta.Update:   iamtypes.EditSysModel,
		meta.Create:   iamtypes.EditSysModel,
		meta.Find:     iamtypes.ViewSysModel,
		meta.FindMany: iamtypes.ViewSysModel,
	},
	meta.ModelUnique: {
		meta.Delete:   iamtypes.EditSysModel,
		meta.Update:   iamtypes.EditSysModel,
		meta.Create:   iamtypes.EditSysModel,
		meta.Find:     iamtypes.ViewSysModel,
		meta.FindMany: iamtypes.ViewSysModel,
	},
	meta.Business: {
		meta.Archive:              iamtypes.ArchiveBusiness,
		meta.Create:               iamtypes.CreateBusiness,
		meta.Update:               iamtypes.EditBusiness,
		meta.Find:                 iamtypes.FindBusiness,
		meta.ViewBusinessResource: iamtypes.ViewBusinessResource,
	},
	meta.BizSet: {
		meta.Create:       iamtypes.CreateBizSet,
		meta.Update:       iamtypes.EditBizSet,
		meta.Delete:       iamtypes.DeleteBizSet,
		meta.Find:         iamtypes.ViewBizSet,
		meta.AccessBizSet: iamtypes.AccessBizSet,
	},
	meta.DynamicGrouping: {
		meta.Delete:   iamtypes.DeleteBusinessCustomQuery,
		meta.Update:   iamtypes.EditBusinessCustomQuery,
		meta.Create:   iamtypes.CreateBusinessCustomQuery,
		meta.Find:     iamtypes.ViewBusinessResource,
		meta.FindMany: iamtypes.ViewBusinessResource,
		meta.Execute:  iamtypes.ViewBusinessResource,
	},
	meta.MainlineModel: {
		meta.Find:   iamtypes.Skip,
		meta.Create: iamtypes.EditBusinessLayer,
		meta.Delete: iamtypes.EditBusinessLayer,
	},
	meta.ModelTopology: {
		meta.Find:              iamtypes.EditModelTopologyView,
		meta.Update:            iamtypes.EditModelTopologyView,
		meta.ModelTopologyView: iamtypes.ViewModelTopo,
	},
	meta.MainlineModelTopology: {
		meta.Find: iamtypes.Skip,
	},
	meta.Process: {
		meta.Find:   iamtypes.Skip,
		meta.Create: iamtypes.EditBusinessServiceInstance,
		meta.Delete: iamtypes.EditBusinessServiceInstance,
		meta.Update: iamtypes.EditBusinessServiceInstance,
	},
	meta.HostInstance: {
		meta.MoveResPoolHostToBizIdleModule: iamtypes.ResourcePoolHostTransferToBusiness,
		meta.MoveResPoolHostToDirectory:     iamtypes.ResourcePoolHostTransferToDirectory,
		meta.MoveBizHostFromModuleToResPool: iamtypes.BusinessHostTransferToResourcePool,
		meta.AddHostToResourcePool:          iamtypes.CreateResourcePoolHost,
		meta.Create:                         iamtypes.CreateResourcePoolHost,
		meta.Delete:                         iamtypes.DeleteResourcePoolHost,
		meta.MoveHostToAnotherBizModule:     iamtypes.HostTransferAcrossBusiness,
		meta.Find:                           iamtypes.ViewResourcePoolHost,
		meta.FindMany:                       iamtypes.ViewResourcePoolHost,
		meta.ManageHostAgentID:              iamtypes.ManageHostAgentID,
	},
	meta.ProcessServiceCategory: {
		meta.Delete: iamtypes.DeleteBusinessServiceCategory,
		meta.Update: iamtypes.EditBusinessServiceCategory,
		meta.Create: iamtypes.CreateBusinessServiceCategory,
		meta.Find:   iamtypes.Skip,
	},
	meta.ProcessServiceInstance: {
		meta.Delete:   iamtypes.DeleteBusinessServiceInstance,
		meta.Update:   iamtypes.EditBusinessServiceInstance,
		meta.Create:   iamtypes.CreateBusinessServiceInstance,
		meta.Find:     iamtypes.Skip,
		meta.FindMany: iamtypes.Skip,
	},
	meta.ProcessServiceTemplate: {
		meta.Delete:   iamtypes.DeleteBusinessServiceTemplate,
		meta.Update:   iamtypes.EditBusinessServiceTemplate,
		meta.Create:   iamtypes.CreateBusinessServiceTemplate,
		meta.Find:     iamtypes.Skip,
		meta.FindMany: iamtypes.Skip,
	},
	meta.SetTemplate: {
		meta.Delete:   iamtypes.DeleteBusinessSetTemplate,
		meta.Update:   iamtypes.EditBusinessSetTemplate,
		meta.Create:   iamtypes.CreateBusinessSetTemplate,
		meta.Find:     iamtypes.Skip,
		meta.FindMany: iamtypes.Skip,
	},
	meta.ModelModule: {
		meta.Delete:   iamtypes.DeleteBusinessTopology,
		meta.Update:   iamtypes.EditBusinessTopology,
		meta.Create:   iamtypes.CreateBusinessTopology,
		meta.Find:     iamtypes.Skip,
		meta.FindMany: iamtypes.Skip,
	},
	meta.ModelSet: {
		meta.Delete:   iamtypes.DeleteBusinessTopology,
		meta.Update:   iamtypes.EditBusinessTopology,
		meta.Create:   iamtypes.CreateBusinessTopology,
		meta.Find:     iamtypes.Skip,
		meta.FindMany: iamtypes.Skip,
	},
	meta.MainlineInstance: {
		meta.Delete:   iamtypes.DeleteBusinessTopology,
		meta.Update:   iamtypes.EditBusinessTopology,
		meta.Create:   iamtypes.CreateBusinessTopology,
		meta.Find:     iamtypes.Skip,
		meta.FindMany: iamtypes.Skip,
	},
	meta.MainlineInstanceTopology: {
		meta.Delete: iamtypes.Skip,
		meta.Update: iamtypes.Skip,
		meta.Create: iamtypes.Skip,
		meta.Find:   iamtypes.Skip,
	},
	meta.HostApply: {
		meta.Create:           iamtypes.EditBusinessHostApply,
		meta.Update:           iamtypes.EditBusinessHostApply,
		meta.Delete:           iamtypes.EditBusinessHostApply,
		meta.Find:             iamtypes.Skip,
		meta.DefaultHostApply: iamtypes.ViewBusinessResource,
	},
	meta.ResourcePoolDirectory: {
		meta.Delete:                iamtypes.DeleteResourcePoolDirectory,
		meta.Update:                iamtypes.EditResourcePoolDirectory,
		meta.Create:                iamtypes.CreateResourcePoolDirectory,
		meta.AddHostToResourcePool: iamtypes.CreateResourcePoolHost,
		meta.Find:                  iamtypes.Skip,
	},
	meta.CloudAreaInstance: {
		meta.Delete:   iamtypes.DeleteCloudArea,
		meta.Update:   iamtypes.EditCloudArea,
		meta.Create:   iamtypes.CreateCloudArea,
		meta.Find:     iamtypes.ViewCloudArea,
		meta.FindMany: iamtypes.ViewCloudArea,
	},
	meta.Model: {
		meta.Delete:   iamtypes.DeleteSysModel,
		meta.Update:   iamtypes.EditSysModel,
		meta.Create:   iamtypes.CreateSysModel,
		meta.Find:     iamtypes.ViewSysModel,
		meta.FindMany: iamtypes.ViewSysModel,
	},
	meta.AssociationType: {
		meta.Delete:   iamtypes.DeleteAssociationType,
		meta.Update:   iamtypes.EditAssociationType,
		meta.Create:   iamtypes.CreateAssociationType,
		meta.Find:     iamtypes.Skip,
		meta.FindMany: iamtypes.Skip,
	},
	meta.ModelClassification: {
		meta.Delete:   iamtypes.DeleteModelGroup,
		meta.Update:   iamtypes.EditModelGroup,
		meta.Create:   iamtypes.CreateModelGroup,
		meta.Find:     iamtypes.Skip,
		meta.FindMany: iamtypes.Skip,
	},
	meta.AuditLog: {
		meta.Find:     iamtypes.FindAuditLog,
		meta.FindMany: iamtypes.FindAuditLog,
	},
	meta.SystemBase: {
		meta.ModelTopologyView:      iamtypes.EditModelTopologyView,
		meta.ModelTopologyOperation: iamtypes.EditBusinessLayer,
	},
	meta.EventWatch: {
		meta.WatchHost:             iamtypes.WatchHostEvent,
		meta.WatchHostRelation:     iamtypes.WatchHostRelationEvent,
		meta.WatchBiz:              iamtypes.WatchBizEvent,
		meta.WatchSet:              iamtypes.WatchSetEvent,
		meta.WatchModule:           iamtypes.WatchModuleEvent,
		meta.WatchProcess:          iamtypes.WatchProcessEvent,
		meta.WatchCommonInstance:   iamtypes.WatchCommonInstanceEvent,
		meta.WatchMainlineInstance: iamtypes.WatchMainlineInstanceEvent,
		meta.WatchInstAsst:         iamtypes.WatchInstAsstEvent,
		meta.WatchBizSet:           iamtypes.WatchBizSetEvent,
		meta.WatchPlat:             iamtypes.WatchPlatEvent,
		meta.WatchKubeCluster:      iamtypes.WatchKubeClusterEvent,
		meta.WatchKubeNode:         iamtypes.WatchKubeNodeEvent,
		meta.WatchKubeNamespace:    iamtypes.WatchKubeNamespaceEvent,
		meta.WatchKubeWorkload:     iamtypes.WatchKubeWorkloadEvent,
		meta.WatchKubePod:          iamtypes.WatchKubePodEvent,
		meta.WatchProject:          iamtypes.WatchProjectEvent,
	},
	meta.UserCustom: {
		meta.Find:   iamtypes.Skip,
		meta.Update: iamtypes.Skip,
		meta.Delete: iamtypes.Skip,
		meta.Create: iamtypes.Skip,
	},
	meta.ModelAssociation: {
		meta.Find:     iamtypes.ViewSysModel,
		meta.FindMany: iamtypes.ViewSysModel,
		meta.Update:   iamtypes.EditSysModel,
		meta.Delete:   iamtypes.EditSysModel,
		meta.Create:   iamtypes.EditSysModel,
	},
	meta.ModelInstanceTopology: {
		meta.Find:   iamtypes.Skip,
		meta.Update: iamtypes.Skip,
		meta.Delete: iamtypes.Skip,
		meta.Create: iamtypes.Skip,
	},
	meta.ModelAttribute: {
		meta.Find:   iamtypes.ViewSysModel,
		meta.Update: iamtypes.EditSysModel,
		meta.Delete: iamtypes.DeleteSysModel,
		meta.Create: iamtypes.CreateSysModel,
	},
	meta.HostFavorite: {
		meta.Find:   iamtypes.Skip,
		meta.Update: iamtypes.Skip,
		meta.Delete: iamtypes.Skip,
		meta.Create: iamtypes.Skip,
	},

	meta.ProcessTemplate: {
		meta.Find:   iamtypes.Skip,
		meta.Delete: iamtypes.DeleteBusinessServiceTemplate,
		meta.Update: iamtypes.EditBusinessServiceTemplate,
		meta.Create: iamtypes.CreateBusinessServiceTemplate,
	},
	meta.BizTopology: {
		meta.Find:   iamtypes.Skip,
		meta.Update: iamtypes.EditBusinessTopology,
		meta.Delete: iamtypes.DeleteBusinessTopology,
		meta.Create: iamtypes.CreateBusinessTopology,
	},
	// unsupported resource actions for now
	meta.NetDataCollector: {
		meta.Find:   iamtypes.Unsupported,
		meta.Update: iamtypes.Unsupported,
		meta.Delete: iamtypes.Unsupported,
		meta.Create: iamtypes.Unsupported,
	},
	meta.InstallBK: {
		meta.Update: iamtypes.Skip,
	},
	// TODO: confirm this
	meta.SystemConfig: {
		meta.FindMany: iamtypes.Skip,
		meta.Find:     iamtypes.Skip,
		meta.Update:   iamtypes.Skip,
		meta.Delete:   iamtypes.Skip,
		meta.Create:   iamtypes.Skip,
	},
	meta.ConfigAdmin: {

		// reuse GlobalSettings permissions
		meta.Find:   iamtypes.Skip,
		meta.Update: iamtypes.GlobalSettings,
		meta.Delete: iamtypes.Unsupported,
		meta.Create: iamtypes.Unsupported,
	},
	meta.KubeCluster: {
		meta.Find:   iamtypes.ViewBusinessResource,
		meta.Update: iamtypes.EditContainerCluster,
		meta.Delete: iamtypes.DeleteContainerCluster,
		meta.Create: iamtypes.CreateContainerCluster,
	},
	meta.KubeNode: {
		meta.Find:   iamtypes.ViewBusinessResource,
		meta.Update: iamtypes.EditContainerNode,
		meta.Delete: iamtypes.DeleteContainerNode,
		meta.Create: iamtypes.CreateContainerNode,
	},
	meta.KubeNamespace: {
		meta.Find:   iamtypes.ViewBusinessResource,
		meta.Update: iamtypes.EditContainerNamespace,
		meta.Delete: iamtypes.DeleteContainerNamespace,
		meta.Create: iamtypes.CreateContainerNamespace,
	},
	meta.KubeWorkload: {
		meta.Find:   iamtypes.ViewBusinessResource,
		meta.Update: iamtypes.EditContainerWorkload,
		meta.Delete: iamtypes.DeleteContainerWorkload,
		meta.Create: iamtypes.CreateContainerWorkload,
	},
	meta.KubeDeployment: {
		meta.Find:   iamtypes.ViewBusinessResource,
		meta.Update: iamtypes.EditContainerWorkload,
		meta.Delete: iamtypes.DeleteContainerWorkload,
		meta.Create: iamtypes.CreateContainerWorkload,
	},
	meta.KubeStatefulSet: {
		meta.Find:   iamtypes.ViewBusinessResource,
		meta.Update: iamtypes.EditContainerWorkload,
		meta.Delete: iamtypes.DeleteContainerWorkload,
		meta.Create: iamtypes.CreateContainerWorkload,
	},
	meta.KubeDaemonSet: {
		meta.Find:   iamtypes.ViewBusinessResource,
		meta.Update: iamtypes.EditContainerWorkload,
		meta.Delete: iamtypes.DeleteContainerWorkload,
		meta.Create: iamtypes.CreateContainerWorkload,
	},
	meta.KubeGameStatefulSet: {
		meta.Find:   iamtypes.ViewBusinessResource,
		meta.Update: iamtypes.EditContainerWorkload,
		meta.Delete: iamtypes.DeleteContainerWorkload,
		meta.Create: iamtypes.CreateContainerWorkload,
	},
	meta.KubeGameDeployment: {
		meta.Find:   iamtypes.ViewBusinessResource,
		meta.Update: iamtypes.EditContainerWorkload,
		meta.Delete: iamtypes.DeleteContainerWorkload,
		meta.Create: iamtypes.CreateContainerWorkload,
	},
	meta.KubeCronJob: {
		meta.Find:   iamtypes.ViewBusinessResource,
		meta.Update: iamtypes.EditContainerWorkload,
		meta.Delete: iamtypes.DeleteContainerWorkload,
		meta.Create: iamtypes.CreateContainerWorkload,
	},
	meta.KubeJob: {
		meta.Find:   iamtypes.ViewBusinessResource,
		meta.Update: iamtypes.EditContainerWorkload,
		meta.Delete: iamtypes.DeleteContainerWorkload,
		meta.Create: iamtypes.CreateContainerWorkload,
	},
	meta.KubePodWorkload: {
		meta.Find:   iamtypes.ViewBusinessResource,
		meta.Update: iamtypes.EditContainerWorkload,
		meta.Delete: iamtypes.DeleteContainerWorkload,
		meta.Create: iamtypes.CreateContainerWorkload,
	},
	meta.KubePod: {
		meta.Find:   iamtypes.ViewBusinessResource,
		meta.Delete: iamtypes.DeleteContainerPod,
		meta.Create: iamtypes.CreateContainerPod,
	},
	meta.KubeContainer: {
		meta.Find: iamtypes.ViewBusinessResource,
	},
	meta.Project: {
		meta.Find:   iamtypes.ViewProject,
		meta.Update: iamtypes.EditProject,
		meta.Delete: iamtypes.DeleteProject,
		meta.Create: iamtypes.CreateProject,
	},
	meta.FulltextSearch: {
		meta.Find: iamtypes.UseFulltextSearch,
	},
	meta.FieldTemplate: {
		meta.Create: iamtypes.CreateFieldGroupingTemplate,
		meta.Find:   iamtypes.ViewFieldGroupingTemplate,
		meta.Update: iamtypes.EditFieldGroupingTemplate,
		meta.Delete: iamtypes.DeleteFieldGroupingTemplate,
	},
	meta.IDRuleIncrID: {
		meta.Update: iamtypes.EditIDRuleIncrID,
	},
	meta.FullSyncCond: {
		meta.Create: iamtypes.CreateFullSyncCond,
		meta.Find:   iamtypes.ViewFullSyncCond,
		meta.Update: iamtypes.EditFullSyncCond,
		meta.Delete: iamtypes.DeleteFullSyncCond,
	},
	meta.GeneralCache: {
		meta.Find: iamtypes.ViewGeneralCache,
	},
	meta.TenantSet: {
		meta.Find:            iamtypes.ViewTenantSet,
		meta.AccessTenantSet: iamtypes.AccessTenantSet,
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
				TypeName: ResourceTypeIDMap[iamtypes.TypeID(typeAndID[0])],
				ID:       id,
			})
		}
	}
	return instances, nil
}

// GenIAMDynamicResTypeID 生成IAM侧资源的的dynamic resource typeID
func GenIAMDynamicResTypeID(modelID int64) iamtypes.TypeID {
	return iamtypes.TypeID(fmt.Sprintf("%s%d", iamtypes.IAMSysInstTypePrefix, modelID))
}

// GenCMDBDynamicResType 生成CMDB侧资源的的dynamic resourceType
func GenCMDBDynamicResType(modelID int64) meta.ResourceType {
	return meta.ResourceType(fmt.Sprintf("%s%d", meta.CMDBSysInstTypePrefix, modelID))
}

// genDynamicResourceType generate dynamic resourceType
func genDynamicResourceType(obj metadata.Object) iam.ResourceType {
	return iam.ResourceType{
		ID:      GenIAMDynamicResTypeID(obj.ID),
		Name:    obj.ObjectName,
		NameEn:  obj.ObjectID,
		Parents: nil,
		ProviderConfig: iam.ResourceConfig{
			Path: "/auth/v3/find/resource",
		},
		Version: 1,
	}
}

// genDynamicResourceTypes generate dynamic resourceTypes
func genDynamicResourceTypes(objects []metadata.Object) []iam.ResourceType {
	resourceTypes := make([]iam.ResourceType, 0)
	for _, obj := range objects {
		resourceTypes = append(resourceTypes, genDynamicResourceType(obj))
	}
	return resourceTypes
}

// genIAMDynamicInstanceSelection generate IAM dynamic instanceSelection
func genIAMDynamicInstanceSelection(modelID int64) iamtypes.InstanceSelectionID {
	return iamtypes.InstanceSelectionID(fmt.Sprintf("%s%d", iamtypes.IAMSysInstTypePrefix, modelID))
}

// genDynamicInstanceSelection generate dynamic instanceSelection
func genDynamicInstanceSelection(obj metadata.Object) iam.InstanceSelection {
	return iam.InstanceSelection{
		ID:     genIAMDynamicInstanceSelection(obj.ID),
		Name:   obj.ObjectName,
		NameEn: obj.ObjectID,
		ResourceTypeChain: []iam.ResourceChain{{
			SystemID: iamtypes.SystemIDCMDB,
			ID:       GenIAMDynamicResTypeID(obj.ID),
		}},
	}
}

// genDynamicInstanceSelections generate dynamic instanceSelections
func genDynamicInstanceSelections(objects []metadata.Object) []iam.InstanceSelection {
	instanceSelections := make([]iam.InstanceSelection, 0)
	for _, obj := range objects {
		instanceSelections = append(instanceSelections, genDynamicInstanceSelection(obj))
	}
	return instanceSelections
}

// genDynamicAction generate dynamic action
// Note: view action must be in the first place
func genDynamicAction(obj metadata.Object) []iamtypes.DynamicAction {
	return []iamtypes.DynamicAction{
		genDynamicViewAction(obj),
		genDynamicCreateAction(obj),
		genDynamicEditAction(obj),
		genDynamicDeleteAction(obj),
	}
}

// GenDynamicActionID generate dynamic ActionID
func GenDynamicActionID(actionType iamtypes.ActionType, modelID int64) iamtypes.ActionID {
	return iamtypes.ActionID(fmt.Sprintf("%s_%s%d", actionType, iamtypes.IAMSysInstTypePrefix, modelID))
}

// genDynamicViewAction generate dynamic view action
func genDynamicViewAction(obj metadata.Object) iamtypes.DynamicAction {
	return iamtypes.DynamicAction{
		ActionID:     GenDynamicActionID(iamtypes.View, obj.ID),
		ActionType:   iamtypes.View,
		ActionNameCN: fmt.Sprintf("%s%s%s", obj.ObjectName, "实例", "查看"),
		ActionNameEN: fmt.Sprintf("%s %s %s", "view", obj.ObjectID, "instance"),
	}
}

// genDynamicCreateAction generate dynamic create action
func genDynamicCreateAction(obj metadata.Object) iamtypes.DynamicAction {
	return iamtypes.DynamicAction{
		ActionID:     GenDynamicActionID(iamtypes.Create, obj.ID),
		ActionType:   iamtypes.Create,
		ActionNameCN: fmt.Sprintf("%s%s%s", obj.ObjectName, "实例", "新建"),
		ActionNameEN: fmt.Sprintf("%s %s %s", "create", obj.ObjectID, "instance"),
	}
}

// genDynamicEditAction generate dynamic edit action
func genDynamicEditAction(obj metadata.Object) iamtypes.DynamicAction {
	return iamtypes.DynamicAction{
		ActionID:     GenDynamicActionID(iamtypes.Edit, obj.ID),
		ActionType:   iamtypes.Edit,
		ActionNameCN: fmt.Sprintf("%s%s%s", obj.ObjectName, "实例", "编辑"),
		ActionNameEN: fmt.Sprintf("%s %s %s", "edit", obj.ObjectID, "instance"),
	}
}

// genDynamicDeleteAction generate dynamic delete action
func genDynamicDeleteAction(obj metadata.Object) iamtypes.DynamicAction {
	return iamtypes.DynamicAction{
		ActionID:     GenDynamicActionID(iamtypes.Delete, obj.ID),
		ActionType:   iamtypes.Delete,
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
func genDynamicActionIDs(object metadata.Object) []iamtypes.ActionID {
	actions := genDynamicAction(object)
	actionIDs := make([]iamtypes.ActionID, len(actions))
	for idx, action := range actions {
		actionIDs[idx] = action.ActionID
	}
	return actionIDs
}

// genDynamicActions generate dynamic model actions
func genDynamicActions(objects []metadata.Object) []iam.ResourceAction {
	resActions := make([]iam.ResourceAction, 0)
	for _, obj := range objects {
		relatedResource := []iam.RelateResourceType{
			{
				SystemID:    iamtypes.SystemIDCMDB,
				ID:          GenIAMDynamicResTypeID(obj.ID),
				NameAlias:   "",
				NameAliasEn: "",
				Scope:       nil,
				// 配置权限时可选择实例和配置属性, 后者用于属性鉴权
				SelectionMode: iamtypes.ModeAll,
				InstanceSelections: []iam.RelatedInstanceSelection{{
					SystemID: iamtypes.SystemIDCMDB,
					ID:       genIAMDynamicInstanceSelection(obj.ID),
				}},
			},
		}

		actions := genDynamicAction(obj)
		var relatedActions []iamtypes.ActionID
		for _, action := range actions {
			switch action.ActionType {
			case iamtypes.View:
				resActions = append(resActions, iam.ResourceAction{
					ID:                   action.ActionID,
					Name:                 action.ActionNameCN,
					NameEn:               action.ActionNameEN,
					Type:                 iamtypes.View,
					RelatedActions:       nil,
					RelatedResourceTypes: nil,
					Version:              1,
				})
				relatedActions = []iamtypes.ActionID{action.ActionID}

			case iamtypes.Create:
				resActions = append(resActions, iam.ResourceAction{
					ID:                   action.ActionID,
					Name:                 action.ActionNameCN,
					NameEn:               action.ActionNameEN,
					Type:                 iamtypes.Create,
					RelatedResourceTypes: nil,
					RelatedActions:       nil,
					Version:              1,
				})
			case iamtypes.Edit:
				resActions = append(resActions, iam.ResourceAction{
					ID:                   action.ActionID,
					Name:                 action.ActionNameCN,
					NameEn:               action.ActionNameEN,
					Type:                 iamtypes.Edit,
					RelatedActions:       relatedActions,
					Version:              1,
					RelatedResourceTypes: relatedResource,
				})

			case iamtypes.Delete:
				resActions = append(resActions, iam.ResourceAction{
					ID:                   action.ActionID,
					Name:                 action.ActionNameCN,
					NameEn:               action.ActionNameEN,
					Type:                 iamtypes.Delete,
					RelatedResourceTypes: relatedResource,
					RelatedActions:       relatedActions,
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
func IsIAMSysInstance(resourceType iamtypes.TypeID) bool {
	return strings.HasPrefix(string(resourceType), iamtypes.IAMSysInstTypePrefix)
}

// IsCMDBSysInstance judge whether the resource type is a system instance in cmdb resource
func IsCMDBSysInstance(resourceType meta.ResourceType) bool {
	return strings.HasPrefix(string(resourceType), meta.CMDBSysInstTypePrefix)
}

// isIAMSysInstanceSelection judge whether the instance selection is a system instance selection in iam resource
func isIAMSysInstanceSelection(instanceSelectionID iamtypes.InstanceSelectionID) bool {
	return strings.Contains(string(instanceSelectionID), iamtypes.IAMSysInstTypePrefix)
}

// isIAMSysInstanceAction judge whether the action is a system instance action in iam resource
func isIAMSysInstanceAction(actionID iamtypes.ActionID) bool {
	return strings.Contains(string(actionID), iamtypes.IAMSysInstTypePrefix)
}

// GetModelIDFromIamSysInstance get model id from iam system instance
func GetModelIDFromIamSysInstance(resourceType iamtypes.TypeID) (int64, error) {
	if !IsIAMSysInstance(resourceType) {
		return 0, fmt.Errorf("resourceType %s is not an iam system instance, it must start with prefix %s",
			resourceType, iamtypes.IAMSysInstTypePrefix)
	}
	modelIDStr := strings.TrimPrefix(string(resourceType), iamtypes.IAMSysInstTypePrefix)
	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		blog.ErrorJSON("modelID convert to int64 failed, err:%s, input:%s", err, modelID)
		return 0, fmt.Errorf("get model id failed, parse to int err:%s, the format of resourceType:%s is wrong",
			err.Error(), resourceType)
	}

	return modelID, nil
}

// GetActionTypeFromIAMSysInstance get action type from iam system instance
func GetActionTypeFromIAMSysInstance(actionID iamtypes.ActionID) iamtypes.ActionType {
	actionIDStr := string(actionID)
	return iamtypes.ActionType(actionIDStr[:strings.Index(actionIDStr, "_")])
}
