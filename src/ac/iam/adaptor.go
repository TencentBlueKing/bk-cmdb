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
	"net/http"
	"strconv"

	"configcenter/src/ac/meta"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

var NotEnoughLayer = fmt.Errorf("not enough layer")

func AdaptAuthOptions(a *meta.ResourceAttribute) (ActionID, []types.Resource, error) {

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

	resource, err := genIamResource(action, *rscType, a)
	if err != nil {
		return "", nil, err
	}

	return action, resource, nil
}

func ConvertResourceType(resourceType meta.ResourceType, businessID int64) (*TypeID, error) {
	var iamResourceType TypeID
	switch resourceType {
	case meta.Business:
		iamResourceType = Business
	case meta.Model,
		meta.ModelUnique,
		meta.ModelAttributeGroup:
		iamResourceType = SysModel
	case meta.ModelAttribute:
		if businessID > 0 {
			iamResourceType = BizCustomField
		} else {
			iamResourceType = SysModel
		}
	case meta.ModelModule, meta.ModelSet, meta.MainlineInstance, meta.MainlineInstanceTopology:
		iamResourceType = BizTopology
	case meta.MainlineModel, meta.ModelTopology:
		iamResourceType = SysSystemBase
	case meta.ModelClassification:
		iamResourceType = SysModelGroup
	case meta.AssociationType:
		iamResourceType = SysAssociationType
	case meta.ModelAssociation:
		iamResourceType = SysModel
	case meta.ModelInstanceAssociation:
		iamResourceType = SysInstance
	case meta.MainlineModelTopology:
		iamResourceType = SysSystemBase
	case meta.ModelInstance:
		iamResourceType = SysInstance
	case meta.ModelInstanceTopology:
		iamResourceType = SkipType
	case meta.CloudAreaInstance:
		iamResourceType = SysCloudArea
	case meta.HostInstance:
		iamResourceType = Host
	case meta.HostFavorite:
		iamResourceType = Host
	case meta.Process:
		iamResourceType = BizProcessServiceInstance
	case meta.EventPushing:
		iamResourceType = SysEventPushing
	case meta.DynamicGrouping:
		iamResourceType = BizCustomQuery
	case meta.AuditLog:
		iamResourceType = SysAuditLog
	case meta.SystemBase:
		iamResourceType = SysSystemBase
	case meta.UserCustom:
		iamResourceType = UserCustom
	case meta.NetDataCollector:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	case meta.ProcessServiceTemplate:
		iamResourceType = BizProcessServiceTemplate
	case meta.ProcessServiceCategory:
		iamResourceType = BizProcessServiceCategory
	case meta.ProcessServiceInstance:
		iamResourceType = BizProcessServiceInstance
	case meta.BizTopology:
		iamResourceType = BizTopology
	case meta.SetTemplate:
		iamResourceType = BizSetTemplate
	case meta.OperationStatistic:
		iamResourceType = SysOperationStatistic
	case meta.HostApply:
		iamResourceType = BizHostApply
	case meta.ResourcePoolDirectory:
		iamResourceType = SysResourcePoolDirectory
	case meta.CloudAccount:
		iamResourceType = SysCloudAccount
	case meta.CloudResourceTask:
		iamResourceType = SysCloudResourceTask
	case meta.EventWatch:
		iamResourceType = SysEventWatch
	case meta.ConfigAdmin:
		iamResourceType = SysSystemBase
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	return &iamResourceType, nil
}

func ConvertResourceAction(resourceType meta.ResourceType, action meta.Action, businessID int64) (ActionID, error) {
	if action == meta.SkipAction {
		return Skip, nil
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

	if resourceType == meta.ModelAttribute && businessID > 0 {
		if convertAction == meta.Delete || convertAction == meta.Update || convertAction == meta.Create {
			if businessID > 0 {
				return EditBusinessCustomField, nil
			} else {
				return EditSysModel, nil
			}
		}
	}

	if resourceType == meta.HostInstance && convertAction == meta.Update {
		if businessID > 0 {
			return EditBusinessHost, nil
		} else {
			return EditResourcePoolHost, nil
		}
	}

	if resourceType == meta.ModelAttributeGroup && businessID > 0 {
		return EditBusinessCustomField, nil
	}

	if _, exist := resourceActionMap[resourceType]; exist {
		actionID, ok := resourceActionMap[resourceType][convertAction]
		if ok && actionID != Unsupported {
			return actionID, nil
		}
	}
	return Unsupported, fmt.Errorf("unsupported type %s action: %s", resourceType, action)
}

var resourceActionMap = map[meta.ResourceType]map[meta.Action]ActionID{
	meta.ModelInstance: {
		meta.Delete:   DeleteSysInstance,
		meta.Update:   EditSysInstance,
		meta.Create:   CreateSysInstance,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.ModelAttributeGroup: {
		meta.Delete:   EditSysModel,
		meta.Update:   EditSysModel,
		meta.Create:   EditSysModel,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.ModelUnique: {
		meta.Delete:   EditSysModel,
		meta.Update:   EditSysModel,
		meta.Create:   EditSysModel,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.Business: {
		meta.Archive:              ArchiveBusiness,
		meta.Create:               CreateBusiness,
		meta.Update:               EditBusiness,
		meta.Find:                 FindBusiness,
		meta.ViewBusinessResource: ViewBusinessResource,
	},
	meta.DynamicGrouping: {
		meta.Delete:   DeleteBusinessCustomQuery,
		meta.Update:   EditBusinessCustomQuery,
		meta.Create:   CreateBusinessCustomQuery,
		meta.Find:     FindBusinessCustomQuery,
		meta.FindMany: FindBusinessCustomQuery,
		meta.Execute:  FindBusinessCustomQuery,
	},
	meta.MainlineModel: {
		meta.Find:   Skip,
		meta.Create: EditBusinessLayer,
		meta.Delete: EditBusinessLayer,
	},
	meta.ModelTopology: {
		meta.Find:   EditModelTopologyView,
		meta.Update: EditModelTopologyView,
	},
	meta.MainlineModelTopology: {
		meta.Find: Skip,
	},
	meta.Process: {
		meta.Find:                   Skip,
		meta.Create:                 EditBusinessServiceInstance,
		meta.Delete:                 EditBusinessServiceInstance,
		meta.Update:                 EditBusinessServiceInstance,
	},
	meta.HostInstance: {
		meta.MoveResPoolHostToBizIdleModule: ResourcePoolHostTransferToBusiness,
		meta.MoveResPoolHostToDirectory:     ResourcePoolHostTransferToDirectory,
		meta.MoveBizHostFromModuleToResPool: BusinessHostTransferToResourcePool,
		meta.AddHostToResourcePool:          CreateResourcePoolHost,
		meta.Create:                         CreateResourcePoolHost,
		meta.Delete:                         DeleteResourcePoolHost,
		meta.MoveHostToBizFaultModule:       EditBusinessHost,
		meta.MoveHostToBizIdleModule:        EditBusinessHost,
		meta.MoveHostToBizRecycleModule:     EditBusinessHost,
		meta.MoveHostToAnotherBizModule:     EditBusinessHost,
		meta.CleanHostInSetOrModule:         EditBusinessHost,
		meta.MoveHostWithinBusiness:         EditBusinessHost,
		meta.Find:                           Skip,
		meta.FindMany:                       Skip,
	},
	meta.ProcessServiceCategory: {
		meta.Delete: DeleteBusinessServiceCategory,
		meta.Update: EditBusinessServiceCategory,
		meta.Create: CreateBusinessServiceCategory,
		meta.Find:   Skip,
	},
	meta.ProcessServiceInstance: {
		meta.Delete:   DeleteBusinessServiceInstance,
		meta.Update:   EditBusinessServiceInstance,
		meta.Create:   CreateBusinessServiceInstance,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.ProcessServiceTemplate: {
		meta.Delete:   DeleteBusinessServiceTemplate,
		meta.Update:   EditBusinessServiceTemplate,
		meta.Create:   CreateBusinessServiceTemplate,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.SetTemplate: {
		meta.Delete:   DeleteBusinessSetTemplate,
		meta.Update:   EditBusinessSetTemplate,
		meta.Create:   CreateBusinessSetTemplate,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.ModelModule: {
		meta.Delete:   DeleteBusinessTopology,
		meta.Update:   EditBusinessTopology,
		meta.Create:   CreateBusinessTopology,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.ModelSet: {
		meta.Delete:   DeleteBusinessTopology,
		meta.Update:   EditBusinessTopology,
		meta.Create:   CreateBusinessTopology,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.MainlineInstance: {
		meta.Delete:   DeleteBusinessTopology,
		meta.Update:   EditBusinessTopology,
		meta.Create:   CreateBusinessTopology,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.MainlineInstanceTopology: {
		meta.Delete: Skip,
		meta.Update: Skip,
		meta.Create: Skip,
		meta.Find:   Skip,
	},
	meta.HostApply: {
		meta.Create: EditBusinessHostApply,
		meta.Update: EditBusinessHostApply,
		meta.Delete: EditBusinessHostApply,
		meta.Find:   Skip,
	},
	meta.ResourcePoolDirectory: {
		meta.Delete:                DeleteResourcePoolDirectory,
		meta.Update:                EditResourcePoolDirectory,
		meta.Create:                CreateResourcePoolDirectory,
		meta.AddHostToResourcePool: CreateResourcePoolHost,
		meta.Find:                  Skip,
	},
	meta.CloudAreaInstance: {
		meta.Delete:   DeleteCloudArea,
		meta.Update:   EditCloudArea,
		meta.Create:   CreateCloudArea,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.EventPushing: {
		meta.Delete:   DeleteEventPushing,
		meta.Update:   EditEventPushing,
		meta.Create:   CreateEventPushing,
		meta.Find:     FindEventPushing,
		meta.FindMany: FindEventPushing,
	},
	meta.CloudAccount: {
		meta.Delete:   DeleteCloudAccount,
		meta.Update:   EditCloudAccount,
		meta.Create:   CreateCloudAccount,
		meta.Find:     FindCloudAccount,
		meta.FindMany: FindCloudAccount,
	},
	meta.CloudResourceTask: {
		meta.Delete: DeleteCloudResourceTask,
		meta.Update: EditCloudResourceTask,
		meta.Create: CreateCloudResourceTask,
		meta.Find:   FindCloudResourceTask,
	},
	meta.Model: {
		meta.Delete:   DeleteSysModel,
		meta.Update:   EditSysModel,
		meta.Create:   CreateSysModel,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.AssociationType: {
		meta.Delete:   DeleteAssociationType,
		meta.Update:   EditAssociationType,
		meta.Create:   CreateAssociationType,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.ModelClassification: {
		meta.Delete:   DeleteModelGroup,
		meta.Update:   EditModelGroup,
		meta.Create:   CreateModelGroup,
		meta.Find:     Skip,
		meta.FindMany: Skip,
	},
	meta.OperationStatistic: {
		meta.Create:   EditOperationStatistic,
		meta.Delete:   EditOperationStatistic,
		meta.Update:   EditOperationStatistic,
		meta.Find:     FindOperationStatistic,
		meta.FindMany: FindOperationStatistic,
	},
	meta.AuditLog: {
		meta.FindMany: FindAuditLog,
	},
	meta.SystemBase: {
		meta.ModelTopologyView:      EditModelTopologyView,
		meta.ModelTopologyOperation: EditBusinessTopology,
	},
	meta.EventWatch: {
		meta.WatchHost:         WatchHostEvent,
		meta.WatchHostRelation: WatchHostRelationEvent,
		meta.WatchBiz:          WatchBizEvent,
		meta.WatchSet:          WatchSetEvent,
		meta.WatchModule:       WatchModuleEvent,
	},
	meta.UserCustom: {
		meta.Find:   Skip,
		meta.Update: Skip,
		meta.Delete: Skip,
		meta.Create: Skip,
	},
	meta.ModelInstanceAssociation: {
		meta.Find:     Skip,
		meta.FindMany: Skip,
		meta.Update:   EditSysInstance,
		meta.Delete:   EditSysInstance,
		meta.Create:   EditSysInstance,
	},
	meta.ModelAssociation: {
		meta.Find:     Skip,
		meta.FindMany: Skip,
		meta.Update:   EditSysModel,
		meta.Delete:   EditSysModel,
		meta.Create:   EditSysModel,
	},
	meta.ModelInstanceTopology: {
		meta.Find:   Skip,
		meta.Update: Skip,
		meta.Delete: Skip,
		meta.Create: Skip,
	},
	meta.ModelAttribute: {
		meta.Find:   Skip,
		meta.Update: EditSysModel,
		meta.Delete: DeleteSysModel,
		meta.Create: CreateSysModel,
	},
	meta.HostFavorite: {
		meta.Find:   Skip,
		meta.Update: EditBusinessHost,
		meta.Delete: DeleteResourcePoolHost,
		meta.Create: CreateResourcePoolHost,
	},

	meta.ProcessTemplate: {
		meta.Find:   Skip,
		meta.Delete: DeleteBusinessServiceTemplate,
		meta.Update: EditBusinessServiceTemplate,
		meta.Create: CreateBusinessServiceTemplate,
	},
	meta.BizTopology: {
		meta.Find:   Skip,
		meta.Update: EditBusinessTopology,
		meta.Delete: DeleteBusinessTopology,
		meta.Create: CreateBusinessTopology,
	},
	// unsupported resource actions for now
	meta.NetDataCollector: {
		meta.Find:   Unsupported,
		meta.Update: Unsupported,
		meta.Delete: Unsupported,
		meta.Create: Unsupported,
	},
	meta.InstallBK: {
		meta.Update: Skip,
	},
	// TODO: confirm this
	meta.SystemConfig: {
		meta.Find:   Skip,
		meta.Update: Skip,
		meta.Delete: Skip,
		meta.Create: Skip,
	},
	meta.ConfigAdmin: {
		meta.Find:   Skip,
		meta.Update: GlobalSettings,
		// unsupported action
		meta.Delete: Unsupported,
		// unsupported action
		meta.Create: Unsupported,
	},
}

func GetBizNameByID(clientSet apimachinery.ClientSetInterface, header http.Header, bizID int64) (string, error) {
	ctx := util.NewContextFromHTTPHeader(header)

	result, err := clientSet.TopoServer().Instance().GetAppBasicInfo(ctx, header, bizID)
	if err != nil {
		return "", err
	}
	if result.Code == common.CCNoPermission {
		return "", nil
	}
	if !result.Result {
		return "", result.Error()
	}
	bizName := result.Data.BizName
	return bizName, nil
}

// AdoptPermissions 用于鉴权没有通过时，根据鉴权的资源信息生成需要申请的权限信息
func AdoptPermissions(h http.Header, api apimachinery.ClientSetInterface, rs []meta.ResourceAttribute) (*metadata.IamPermission, error) {
	permission := new(metadata.IamPermission)
	permission.SystemID = SystemIDCMDB
	bizIDMap := make(map[int64]string)
	// permissionMap maps ResourceActionID and ResourceTypeID to ResourceInstances
	permissionMap := make(map[string]map[string][]metadata.IamResourceInstance, 0)
	for _, r := range rs {
		actionID, err := ConvertResourceAction(r.Type, r.Action, r.BusinessID)
		if err != nil {
			return nil, err
		}
		rscType, err := ConvertResourceType(r.Basic.Type, r.BusinessID)
		if err != nil {
			return nil, err
		}
		// TODO: rearrange this.
		rscIDs := make([]RscTypeAndID, 0)
		// rscIDs, err := GenerateResourceID(*rscType, &r)
		// if err != nil {
		// 	return nil, err
		// }

		// generate iam resource instances by biz id and instance itself
		instances := make([]metadata.IamResourceInstance, 0)
		if r.BusinessID > 0 {
			bizName, exist := bizIDMap[r.BusinessID]
			if !exist {
				var err error
				bizName, err = GetBizNameByID(api, h, r.BusinessID)
				if err != nil {
					return nil, err
				} else {
					bizIDMap[r.BusinessID] = bizName
				}
			}
			instances = append(instances, metadata.IamResourceInstance{
				Type: string(Business),
				ID:   strconv.FormatInt(r.BusinessID, 10),
				Name: bizName,
			})
		}
		instance := metadata.IamResourceInstance{
			Type: string(*rscType),
			Name: r.Basic.Name,
		}
		if len(rscIDs) != 0 {
			instance.ID = rscIDs[len(rscIDs)-1].ResourceID
		}
		instances = append(instances, instance)

		if _, ok := permissionMap[string(actionID)]; !ok {
			permissionMap[string(actionID)] = make(map[string][]metadata.IamResourceInstance, 0)
		}
		permissionMap[string(actionID)][string(*rscType)] = instances
	}

	for actionID, permissionTypeMap := range permissionMap {
		action := metadata.IamAction{ID: actionID}
		for rscType, instances := range permissionTypeMap {
			action.RelatedResourceTypes = append(action.RelatedResourceTypes, metadata.IamResourceType{
				SystemID:  SystemIDCMDB,
				Type:      rscType,
				Instances: instances,
			})
		}
		permission.Actions = append(permission.Actions, action)
	}
	return permission, nil
}
