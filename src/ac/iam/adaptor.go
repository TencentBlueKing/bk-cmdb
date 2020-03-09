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
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"configcenter/src/ac/meta"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

var NotEnoughLayer = fmt.Errorf("not enough layer")

// ResourceTypeID is resource's type in auth center.
func adaptor(attribute *meta.ResourceAttribute) (*ResourceInfo, error) {
	var err error
	info := new(ResourceInfo)
	info.ResourceName = attribute.Basic.Name

	resourceTypeID, err := ConvertResourceType(attribute.Type, attribute.BusinessID)
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

func ConvertResourceType(resourceType meta.ResourceType, businessID int64) (*ResourceTypeID, error) {
	var iamResourceType ResourceTypeID
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
		iamResourceType = SysInstance
	case meta.ModelInstanceAssociation:
		iamResourceType = SysInstance
	case meta.MainlineModelTopology:
		iamResourceType = SysSystemBase
	case meta.ModelInstance:
		iamResourceType = SysInstance
	case meta.Plat:
		iamResourceType = SysCloudArea
	case meta.HostInstance:
		if businessID <= 0 {
			iamResourceType = SysHostInstance
		} else {
			iamResourceType = BizHostInstance
		}
	case meta.HostFavorite:
		iamResourceType = BizHostInstance
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
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
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
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	return &iamResourceType, nil
}

func ConvertResourceAction(resourceType meta.ResourceType, action meta.Action, businessID int64) (ResourceActionID, error) {
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
				return EditModel, nil
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

	resourceActionMap := map[meta.ResourceType]map[meta.Action]ResourceActionID{
		meta.ModelInstance: {
			meta.Delete: DeleteInstance,
			meta.Update: EditInstance,
			meta.Create: CreateInstance,
			meta.Find:   FindInstance,
		},
		meta.ModelAttributeGroup: {
			meta.Delete: EditModel,
			meta.Update: EditModel,
			meta.Create: EditModel,
		},
		meta.ModelUnique: {
			meta.Delete: EditModel,
			meta.Update: EditModel,
			meta.Create: EditModel,
		},
		meta.Business: {
			meta.Archive: ArchiveBusiness,
			meta.Create:  CreateBusiness,
			meta.Update:  EditBusiness,
			meta.Find:    FindBusiness,
		},
		meta.DynamicGrouping: {
			meta.Delete:  DeleteBusinessCustomQuery,
			meta.Update:  EditBusinessCustomQuery,
			meta.Create:  CreateBusinessCustomQuery,
			meta.Find:    FindBusinessCustomQuery,
			meta.Execute: FindBusinessCustomQuery,
		},
		meta.MainlineModel: {
			meta.Find:   EditBusinessLayer,
			meta.Create: EditBusinessLayer,
			meta.Delete: EditBusinessLayer,
		},
		meta.ModelTopology: {
			meta.Find:   EditModelTopologyView,
			meta.Update: EditModelTopologyView,
		},
		meta.MainlineModelTopology: {
			meta.Find:   EditBusinessLayer,
			meta.Update: EditBusinessLayer,
		},
		meta.Process: {
			meta.BoundModuleToProcess:   EditBusinessServiceInstance,
			meta.UnboundModuleToProcess: EditBusinessServiceInstance,
		},
		meta.HostInstance: {
			meta.MoveResPoolHostToBizIdleModule: ResourcePoolHostTransferToBusiness,
			meta.MoveResPoolHostToDirectory:     ResourcePoolHostTransferToDirectory,
			meta.MoveHostFromModuleToResPool:    BusinessHostTransferToResourcePool,
			meta.AddHostToResourcePool:          CreateResourcePoolHost,
			meta.Create:                         CreateResourcePoolHost,
			meta.Delete:                         DeleteResourcePoolHost,
			meta.MoveHostToBizFaultModule:       EditBusinessHost,
			meta.MoveHostToBizIdleModule:        EditBusinessHost,
			meta.MoveHostToBizRecycleModule:     EditBusinessHost,
			meta.MoveHostToAnotherBizModule:     EditBusinessHost,
			meta.CleanHostInSetOrModule:         EditBusinessHost,
			meta.TransferHost:                   EditBusinessHost,
			meta.MoveBizHostToModule:            EditBusinessHost,
		},
		meta.ProcessServiceCategory: {
			meta.Delete: DeleteBusinessServiceCategory,
			meta.Update: EditBusinessServiceCategory,
			meta.Create: CreateBusinessServiceCategory,
		},
		meta.ProcessServiceInstance: {
			meta.Delete: DeleteBusinessServiceInstance,
			meta.Update: EditBusinessServiceInstance,
			meta.Create: CreateBusinessServiceInstance,
		},
		meta.ProcessServiceTemplate: {
			meta.Delete: DeleteBusinessServiceTemplate,
			meta.Update: EditBusinessServiceTemplate,
			meta.Create: CreateBusinessServiceTemplate,
		},
		meta.SetTemplate: {
			meta.Delete: DeleteBusinessSetTemplate,
			meta.Update: EditBusinessSetTemplate,
			meta.Create: CreateBusinessSetTemplate,
		},
		meta.ModelModule: {
			meta.Delete: DeleteBusinessTopology,
			meta.Update: EditBusinessTopology,
			meta.Create: CreateBusinessTopology,
		},
		meta.ModelSet: {
			meta.Delete: DeleteBusinessTopology,
			meta.Update: EditBusinessTopology,
			meta.Create: CreateBusinessTopology,
		},
		meta.MainlineInstance: {
			meta.Delete: DeleteBusinessTopology,
			meta.Update: EditBusinessTopology,
			meta.Create: CreateBusinessTopology,
		},
		meta.MainlineInstanceTopology: {
			meta.Delete: DeleteBusinessTopology,
			meta.Update: EditBusinessTopology,
			meta.Create: CreateBusinessTopology,
		},
		meta.HostApply: {
			meta.Update: EditBusinessHostApply,
		},
		meta.ResourcePoolDirectory: {
			meta.Delete: DeleteResourcePoolDirectory,
			meta.Update: EditResourcePoolDirectory,
			meta.Create: CreateResourcePoolDirectory,
		},
		meta.Plat: {
			meta.Delete: DeleteCloudArea,
			meta.Update: EditCloudArea,
			meta.Create: CreateCloudArea,
		},
		meta.EventPushing: {
			meta.Delete: DeleteEventPushing,
			meta.Update: EditEventPushing,
			meta.Create: CreateEventPushing,
			meta.Find:   FindEventPushing,
		},
		meta.CloudAccount: {
			meta.Delete: DeleteCloudAccount,
			meta.Update: EditCloudAccount,
			meta.Create: CreateCloudAccount,
			meta.Find:   FindCloudAccount,
		},
		meta.CloudResourceTask: {
			meta.Delete: DeleteCloudResourceTask,
			meta.Update: EditCloudResourceTask,
			meta.Create: CreateCloudResourceTask,
			meta.Find:   FindCloudResourceTask,
		},
		meta.Model: {
			meta.Delete: DeleteModel,
			meta.Update: EditModel,
			meta.Create: CreateModel,
			meta.Find:   FindModel,
		},
		meta.AssociationType: {
			meta.Delete: DeleteAssociationType,
			meta.Update: EditAssociationType,
			meta.Create: CreateAssociationType,
		},
		meta.ModelClassification: {
			meta.Delete: DeleteModelGroup,
			meta.Update: EditModelGroup,
			meta.Create: CreateModelGroup,
		},
		meta.OperationStatistic: {
			meta.Update: EditOperationStatistic,
			meta.Find:   FindOperationStatistic,
		},
		meta.AuditLog: {
			meta.Find: FindAuditLog,
		},
	}
	if _, exist := resourceActionMap[resourceType]; exist {
		actionID, ok := resourceActionMap[resourceType][convertAction]
		if ok {
			return actionID, nil
		}
	}
	return Unknown, fmt.Errorf("unsupported action: %s", action)
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
		return "", errors.New(result.ErrMsg)
	}
	bizName := result.Data.BizName
	return bizName, nil
}

// AdoptPermissions 用于鉴权没有通过时，根据鉴权的资源信息生成需要申请的权限信息
var actionMap map[ResourceActionID]ResourceAction
var resourceTypeMap map[ResourceTypeID]ResourceType

func AdoptPermissions(h http.Header, api apimachinery.ClientSetInterface, rs []meta.ResourceAttribute) ([]metadata.Permission, error) {
	rid := util.GetHTTPCCRequestID(h)
	language := util.GetLanguage(h)

	ps := make([]metadata.Permission, 0)
	bizIDMap := make(map[int64]string)
	if actionMap == nil {
		actionMap = make(map[ResourceActionID]ResourceAction)
		for _, action := range GenerateActions() {
			actionMap[action.ID] = action
		}
	}
	if resourceTypeMap == nil {
		resourceTypeMap = make(map[ResourceTypeID]ResourceType)
		for _, resource := range GenerateResourceTypes() {
			resourceTypeMap[resource.ID] = resource
		}
	}
	for _, r := range rs {
		var p metadata.Permission
		p.SystemID = SystemIDCMDB
		if language == string(common.English) {
			p.SystemName = SystemNameCMDBEn
		} else {
			p.ScopeName = SystemNameCMDB
		}

		if r.BusinessID > 0 {
			p.ScopeType = ScopeTypeIDBiz
			if language == string(common.English) {
				p.ScopeTypeName = ScopeTypeIDBizNameEn
			} else {
				p.ScopeTypeName = ScopeTypeIDBizName
			}
			p.ScopeID = strconv.FormatInt(r.BusinessID, 10)
			scopeName, exist := bizIDMap[r.BusinessID]
			if !exist {
				var err error
				scopeName, err = GetBizNameByID(api, h, r.BusinessID)
				if err != nil {
					blog.Errorf("AdoptPermissions failed, GetBizNameByID failed, bizID: %d, err: %s, rid: %s", r.BusinessID, err.Error(), rid)
				} else {
					bizIDMap[r.BusinessID] = scopeName
				}
			}
			p.ScopeName = scopeName
		} else {
			p.ScopeType = ScopeTypeIDSystem
			p.ScopeID = SystemIDCMDB
			if language == string(common.English) {
				p.ScopeTypeName = ScopeTypeIDSystemNameEn
				p.ScopeName = SystemNameCMDBEn
			} else {
				p.ScopeTypeName = ScopeTypeIDSystemName
				p.ScopeName = SystemNameCMDB
			}
		}

		actID, err := ConvertResourceAction(r.Type, r.Action, r.BusinessID)
		if err != nil {
			return nil, err
		}
		p.ActionID = string(actID)
		if language == string(common.English) {
			p.ActionName = actionMap[actID].NameEn
		} else {
			p.ActionName = actionMap[actID].Name
		}

		rscType, err := ConvertResourceType(r.Basic.Type, r.BusinessID)
		if err != nil {
			return nil, err
		}

		rscIDs, err := GenerateResourceID(*rscType, &r)
		if err != nil {
			return nil, err
		}

		var rsc metadata.Resource
		rsc.ResourceType = string(*rscType)
		if language == string(common.English) {
			p.ResourceTypeName = resourceTypeMap[*rscType].NameEn
		} else {
			p.ResourceTypeName = resourceTypeMap[*rscType].Name
		}
		if len(rscIDs) != 0 {
			rsc.ResourceID = rscIDs[len(rscIDs)-1].ResourceID
		}
		rsc.ResourceName = r.Basic.Name
		p.Resources = [][]metadata.Resource{{rsc}}
		ps = append(ps, p)
	}
	return ps, nil
}
