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

package service

import (
	"fmt"
	"strconv"
	"strings"

	"configcenter/pkg/cache/general"
	tenantset "configcenter/pkg/types/tenant-set"
	"configcenter/src/ac/iam"
	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/logics"
	"configcenter/src/scene_server/auth_server/types"
)

type resPullMethodGenerator func(*rest.Kit, *logics.Logics) (types.ResourcePullMethod, error)

var resPullMethodGeneratorMap = map[iamtypes.TypeID]resPullMethodGenerator{
	iamtypes.Host: func(kit *rest.Kit, lgc *logics.Logics) (types.ResourcePullMethod, error) {
		return types.ResourcePullMethod{
			ListAttr:             lgc.ListAttr,
			ListAttrValue:        lgc.ListAttrValue,
			ListInstance:         lgc.ListHostInstance,
			FetchInstanceInfo:    lgc.FetchHostInfo,
			ListInstanceByPolicy: lgc.ListHostByPolicy,
		}, nil
	},
	iamtypes.Business:                  getBusinessMethod,
	iamtypes.BusinessForHostTrans:      getBusinessMethod,
	iamtypes.SysCloudArea:              getSysCloudAreaMethod,
	iamtypes.BizCustomQuery:            getBizInstanceMethod,
	iamtypes.BizProcessServiceTemplate: getBizInstanceMethod,
	iamtypes.BizSetTemplate:            getBizInstanceMethod,
	iamtypes.SysModelGroup: getSystemInstanceMethod(mapstr.MapStr{
		common.BKClassificationTypeField: mapstr.MapStr{
			common.BKDBNE: metadata.HiddenType,
		},
	}),
	iamtypes.InstAsstEvent:             getSystemInstanceMethod(nil),
	iamtypes.BizSet:                    getSystemInstanceMethod(nil),
	iamtypes.Project:                   getSystemInstanceMethod(nil),
	iamtypes.FieldGroupingTemplate:     getSystemInstanceMethod(nil),
	iamtypes.SysModel:                  getModelMethod(iamtypes.SysModel),
	iamtypes.SysInstanceModel:          getModelMethod(iamtypes.SysInstanceModel),
	iamtypes.SysModelEvent:             getModelMethod(iamtypes.SysModelEvent),
	iamtypes.MainlineModelEvent:        getModelMethod(iamtypes.MainlineModelEvent),
	iamtypes.SysAssociationType:        getSysAssociationTypeMethod,
	iamtypes.SysResourcePoolDirectory:  getResourcePoolDirectoryMethod,
	iamtypes.SysHostRscPoolDirectory:   getResourcePoolDirectoryMethod,
	iamtypes.SysAuditLog:               getNoResourceMethod,
	iamtypes.BizCustomField:            getNoResourceMethod,
	iamtypes.BizHostApply:              getNoResourceMethod,
	iamtypes.BizTopology:               getNoResourceMethod,
	iamtypes.SysEventWatch:             getNoResourceMethod,
	iamtypes.BizProcessServiceCategory: getNoResourceMethod,
	iamtypes.BizProcessServiceInstance: getNoResourceMethod,
	iamtypes.KubeWorkloadEvent:         genKubeWorkloadEventMethod,
	iamtypes.GeneralCache:              genGeneralCacheMethod,
	iamtypes.Set: func(kit *rest.Kit, lgc *logics.Logics) (types.ResourcePullMethod, error) {
		return types.ResourcePullMethod{ListInstance: lgc.ListSetInstance,
			FetchInstanceInfo: lgc.FetchSetModuleNameInfo}, nil
	},
	iamtypes.Module: func(kit *rest.Kit, lgc *logics.Logics) (types.ResourcePullMethod, error) {
		return types.ResourcePullMethod{ListInstance: lgc.ListModuleInstance,
			FetchInstanceInfo: lgc.FetchSetModuleNameInfo}, nil
	},
	iamtypes.TenantSet: genTenantSetMethod,
}

// genResourcePullMethod generate iam callback methods for input resource type,
// method not set means not related to this kind of instances
func (s *AuthService) genResourcePullMethod(kit *rest.Kit, resourceType iamtypes.TypeID) (types.ResourcePullMethod,
	error) {

	generator, exists := resPullMethodGeneratorMap[resourceType]
	if exists {
		return generator(kit, s.lgc)
	}

	if iam.IsIAMSysInstance(resourceType) {
		return types.ResourcePullMethod{
			ListAttr:          s.lgc.ListAttr,
			ListAttrValue:     s.lgc.ListAttrValue,
			ListInstance:      s.lgc.ListModelInstance,
			FetchInstanceInfo: s.lgc.FetchObjInstInfo,
			ListInstanceByPolicy: func(kit *rest.Kit, resourceType iamtypes.TypeID,
				filter *types.ListInstanceByPolicyFilter, page types.Page) (result *types.ListInstanceResult,
				e error) {
				return s.lgc.ListInstanceByPolicy(kit, resourceType, filter, page, nil)
			},
		}, nil
	}
	return types.ResourcePullMethod{}, fmt.Errorf("gen method failed: unsupported resource type: %s", resourceType)
}

func getBusinessMethod(kit *rest.Kit, lgc *logics.Logics) (types.ResourcePullMethod, error) {
	// business instances should not include resource pool business
	extraCond := map[string]interface{}{
		common.BKDefaultField: map[string]interface{}{
			common.BKDBNE: common.DefaultAppFlag,
		},
	}

	return types.ResourcePullMethod{
		ListAttr:      lgc.ListAttr,
		ListAttrValue: lgc.ListAttrValue,
		ListInstance: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.ListInstanceFilter,
			page types.Page) (*types.ListInstanceResult, error) {
			return lgc.ListSystemInstance(kit, resourceType, filter, page, extraCond)
		},
		FetchInstanceInfo: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.FetchInstanceInfoFilter) (
			[]map[string]interface{}, error) {
			return lgc.FetchInstanceInfo(kit, resourceType, filter, extraCond)
		},
		ListInstanceByPolicy: func(kit *rest.Kit, resourceType iamtypes.TypeID,
			filter *types.ListInstanceByPolicyFilter,
			page types.Page) (result *types.ListInstanceResult, e error) {
			return lgc.ListInstanceByPolicy(kit, resourceType, filter, page, extraCond)
		},
	}, nil
}

func getSysCloudAreaMethod(kit *rest.Kit, lgc *logics.Logics) (types.ResourcePullMethod, error) {
	// cloud area instances should not include default cloud area, since it can't be operated
	extraCond := map[string]interface{}{
		common.BKCloudIDField: map[string]interface{}{
			common.BKDBNE: common.BKDefaultDirSubArea,
		},
	}

	return types.ResourcePullMethod{
		ListAttr:      lgc.ListAttr,
		ListAttrValue: lgc.ListAttrValue,
		ListInstance: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.ListInstanceFilter,
			page types.Page) (*types.ListInstanceResult, error) {
			return lgc.ListSystemInstance(kit, resourceType, filter, page, extraCond)
		},
		FetchInstanceInfo: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.FetchInstanceInfoFilter) (
			[]map[string]interface{}, error) {
			return lgc.FetchInstanceInfo(kit, resourceType, filter, extraCond)
		},
		ListInstanceByPolicy: func(kit *rest.Kit, resourceType iamtypes.TypeID,
			filter *types.ListInstanceByPolicyFilter,
			page types.Page) (result *types.ListInstanceResult, e error) {
			return lgc.ListInstanceByPolicy(kit, resourceType, filter, page, extraCond)
		},
	}, nil
}

func getBizInstanceMethod(kit *rest.Kit, lgc *logics.Logics) (types.ResourcePullMethod, error) {
	return types.ResourcePullMethod{
		ListInstance: lgc.ListBusinessInstance,
		FetchInstanceInfo: func(kit *rest.Kit, resourceType iamtypes.TypeID,
			filter *types.FetchInstanceInfoFilter) (
			[]map[string]interface{}, error) {
			return lgc.FetchInstanceInfo(kit, resourceType, filter, nil)
		},
		ListInstanceByPolicy: func(kit *rest.Kit, resourceType iamtypes.TypeID,
			filter *types.ListInstanceByPolicyFilter,
			page types.Page) (result *types.ListInstanceResult, e error) {
			return lgc.ListInstanceByPolicy(kit, resourceType, filter, page, nil)
		},
	}, nil
}

func getSystemInstanceMethod(extraCond map[string]interface{}) resPullMethodGenerator {
	return func(kit *rest.Kit, lgc *logics.Logics) (types.ResourcePullMethod, error) {
		return types.ResourcePullMethod{
			ListInstance: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.ListInstanceFilter,
				page types.Page) (*types.ListInstanceResult, error) {
				return lgc.ListSystemInstance(kit, resourceType, filter, page, extraCond)
			},
			FetchInstanceInfo: func(kit *rest.Kit, resourceType iamtypes.TypeID,
				filter *types.FetchInstanceInfoFilter) ([]map[string]interface{}, error) {
				return lgc.FetchInstanceInfo(kit, resourceType, filter, extraCond)
			},
			ListInstanceByPolicy: func(kit *rest.Kit, resourceType iamtypes.TypeID,
				filter *types.ListInstanceByPolicyFilter, page types.Page) (result *types.ListInstanceResult, e error) {
				return lgc.ListInstanceByPolicy(kit, resourceType, filter, page, extraCond)
			},
		}, nil
	}
}

func getModelMethod(resourceType iamtypes.TypeID) resPullMethodGenerator {
	return func(kit *rest.Kit, lgc *logics.Logics) (types.ResourcePullMethod, error) {
		// get mainline objects
		mainlineOpt := &metadata.QueryCondition{
			Condition: map[string]interface{}{common.AssociationKindIDField: common.AssociationKindMainline},
		}
		asstRes, err := lgc.CoreAPI.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, mainlineOpt)
		if err != nil {
			blog.Errorf("search mainline association failed, err: %v, rid: %s", err, kit.Rid)
			return types.ResourcePullMethod{}, err
		}

		mainlineObjIDs := make([]string, 0)
		for _, asst := range asstRes.Info {
			if metadata.IsCommon(asst.ObjectID) {
				mainlineObjIDs = append(mainlineObjIDs, asst.ObjectID)
			}
		}

		// process and cloud area are temporarily excluded TODO: remove this restriction when they are available for user
		// instance model is used as parent layer of instances, should exclude host model and mainline model as
		// they use separate operations
		excludedObjIDs := []string{common.BKInnerObjIDProc, common.BKInnerObjIDPlat}

		// get quoted objects
		quoteOpt := &metadata.CommonQueryOption{
			Fields: []string{common.BKDestModelField},
			Page:   metadata.BasePage{Limit: common.BKMaxPageSize},
		}
		relRes, err := lgc.CoreAPI.CoreService().ModelQuote().ListModelQuoteRelation(kit.Ctx, kit.Header, quoteOpt)
		if err != nil {
			blog.Errorf("list model quote relations failed, err: %v, rid: %s", err, kit.Rid)
			return types.ResourcePullMethod{}, err
		}

		for _, res := range relRes.Info {
			excludedObjIDs = append(excludedObjIDs, res.DestModel)
		}

		// generate extra condition for resource type
		var extraCond map[string]interface{}
		switch resourceType {
		case iamtypes.SysModelEvent, iamtypes.SysInstanceModel:
			excludedObjIDs = append(excludedObjIDs, common.BKInnerObjIDHost, common.BKInnerObjIDApp,
				common.BKInnerObjIDSet, common.BKInnerObjIDModule)
			excludedObjIDs = append(excludedObjIDs, mainlineObjIDs...)
			extraCond = map[string]interface{}{
				common.BKObjIDField: map[string]interface{}{common.BKDBNIN: excludedObjIDs},
			}
		case iamtypes.MainlineModelEvent:
			extraCond = map[string]interface{}{
				common.BKObjIDField: map[string]interface{}{common.BKDBIN: mainlineObjIDs},
			}
		default:
			extraCond = map[string]interface{}{
				common.BKObjIDField: map[string]interface{}{common.BKDBNIN: excludedObjIDs},
			}
		}

		return types.ResourcePullMethod{
			ListInstance: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.ListInstanceFilter,
				page types.Page) (*types.ListInstanceResult, error) {
				return lgc.ListSystemInstance(kit, resourceType, filter, page, extraCond)
			},
			FetchInstanceInfo: func(kit *rest.Kit, resourceType iamtypes.TypeID,
				filter *types.FetchInstanceInfoFilter) (
				[]map[string]interface{}, error) {
				return lgc.FetchInstanceInfo(kit, resourceType, filter, extraCond)
			},
			ListInstanceByPolicy: func(kit *rest.Kit, resourceType iamtypes.TypeID,
				filter *types.ListInstanceByPolicyFilter,
				page types.Page) (result *types.ListInstanceResult, e error) {
				return lgc.ListInstanceByPolicy(kit, resourceType, filter, page, extraCond)
			},
		}, nil
	}
}

func getSysAssociationTypeMethod(kit *rest.Kit, lgc *logics.Logics) (types.ResourcePullMethod, error) {
	// association types should not include preset ones, since they can't be operated
	extraCond := map[string]interface{}{
		common.BKIsPre: map[string]interface{}{
			common.BKDBNE: true,
		},
	}

	return types.ResourcePullMethod{
		ListInstance: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.ListInstanceFilter,
			page types.Page) (*types.ListInstanceResult, error) {
			return lgc.ListSystemInstance(kit, resourceType, filter, page, extraCond)
		},
		FetchInstanceInfo: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.FetchInstanceInfoFilter) (
			[]map[string]interface{}, error) {
			return lgc.FetchInstanceInfo(kit, resourceType, filter, extraCond)
		},
		ListInstanceByPolicy: func(kit *rest.Kit, resourceType iamtypes.TypeID,
			filter *types.ListInstanceByPolicyFilter,
			page types.Page) (result *types.ListInstanceResult, e error) {
			return lgc.ListInstanceByPolicy(kit, resourceType, filter, page, extraCond)
		},
	}, nil
}

func getResourcePoolDirectoryMethod(kit *rest.Kit, lgc *logics.Logics) (types.ResourcePullMethod, error) {

	resourcePoolBizID, err := lgc.GetResourcePoolBizID(kit)
	if err != nil {
		return types.ResourcePullMethod{}, err
	}
	// resource pool directory must be in the resource pool business
	extraCond := map[string]interface{}{
		common.BKAppIDField: resourcePoolBizID,
	}

	return types.ResourcePullMethod{
		ListInstance: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.ListInstanceFilter,
			page types.Page) (*types.ListInstanceResult, error) {
			return lgc.ListSystemInstance(kit, resourceType, filter, page, extraCond)
		},
		FetchInstanceInfo: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.FetchInstanceInfoFilter) (
			[]map[string]interface{}, error) {
			return lgc.FetchInstanceInfo(kit, resourceType, filter, extraCond)
		},
		ListInstanceByPolicy: func(kit *rest.Kit, resourceType iamtypes.TypeID,
			filter *types.ListInstanceByPolicyFilter,
			page types.Page) (result *types.ListInstanceResult, e error) {
			return lgc.ListInstanceByPolicy(kit, resourceType, filter, page, extraCond)
		},
	}, nil
}

func getNoResourceMethod(kit *rest.Kit, lgc *logics.Logics) (types.ResourcePullMethod, error) {
	return types.ResourcePullMethod{}, nil
}

// kubeWorkloadKinds kube workload kinds
// TODO define this in kube types folder, and replace the kinds with actual ones, this is only an example
var kubeWorkloadKinds = []string{"deployment", "statefulSet", "daemonSet"}

// genKubeWorkloadEventMethod generate iam callback methods for iamtypes.KubeWorkloadEvent resource type
func genKubeWorkloadEventMethod(kit *rest.Kit, lgc *logics.Logics) (types.ResourcePullMethod, error) {
	return types.ResourcePullMethod{
		ListInstance: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.ListInstanceFilter,
			page types.Page) (*types.ListInstanceResult, error) {
			limit := page.Limit
			if limit > common.BKMaxPageSize && limit != common.BKNoLimit {
				return nil, kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded)
			}
			if limit == 0 {
				return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "page.limit")
			}
			// get kube workload kinds that matches the filter
			kinds := kubeWorkloadKinds
			if filter != nil {
				if filter.Parent != nil {
					return &types.ListInstanceResult{Count: 0, Results: make([]types.InstanceResource, 0)}, nil
				}

				if len(filter.Keyword) != 0 {
					kinds = make([]string, 0)
					for _, kind := range kubeWorkloadKinds {
						if strings.Contains(strings.ToLower(kind), strings.ToLower(filter.Keyword)) {
							kinds = append(kinds, kind)
						}
					}
				}
			}
			// generate iam instance resource by kube workload kinds, do pagination
			kindsLen := int64(len(kinds))
			if page.Offset >= kindsLen {
				return &types.ListInstanceResult{Count: 0, Results: make([]types.InstanceResource, 0)}, nil
			}
			end := page.Offset + limit
			if end > kindsLen {
				end = kindsLen
			}
			res := make([]types.InstanceResource, 0)
			for _, kind := range kinds[page.Offset:end] {
				res = append(res, types.InstanceResource{
					ID:          kind,
					DisplayName: kind,
				})
			}

			return &types.ListInstanceResult{Count: kindsLen, Results: res}, nil
		},
		FetchInstanceInfo: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.FetchInstanceInfoFilter) (
			[]map[string]interface{}, error) {
			// only support query name field, name field is the same with the id field
			hasNameField := false
			for _, attr := range filter.Attrs {
				if attr == types.NameField {
					hasNameField = true
				}
			}
			if !hasNameField {
				return make([]map[string]interface{}, 0), nil
			}
			res := make([]map[string]interface{}, 0)
			for _, id := range filter.IDs {
				if util.InStrArr(kubeWorkloadKinds, id) {
					res = append(res, map[string]interface{}{types.NameField: id})
				}
			}
			return res, nil
		},
		ListInstanceByPolicy: func(kit *rest.Kit, resourceType iamtypes.TypeID,
			filter *types.ListInstanceByPolicyFilter, page types.Page) (*types.ListInstanceResult, error) {
			return nil, fmt.Errorf("%s do not support %s", iamtypes.KubeWorkloadEvent, types.ListInstanceByPolicyMethod)
		},
	}, nil
}

func genGeneralCacheMethod(kit *rest.Kit, lgc *logics.Logics) (types.ResourcePullMethod, error) {
	return types.ResourcePullMethod{
		ListInstance: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.ListInstanceFilter,
			page types.Page) (*types.ListInstanceResult, error) {

			// get supported general cache resource types that matches the filter
			resTypes := make([]string, 0)
			for resType := range general.SupportedResTypeMap {
				resTypes = append(resTypes, string(resType))
			}

			if filter != nil {
				if filter.Parent != nil {
					return &types.ListInstanceResult{Count: 0, Results: make([]types.InstanceResource, 0)}, nil
				}

				if len(filter.Keyword) != 0 {
					filteredResTypes := make([]string, 0)
					for _, resType := range resTypes {
						if strings.Contains(strings.ToLower(resType), strings.ToLower(filter.Keyword)) {
							filteredResTypes = append(filteredResTypes, resType)
						}
					}
					resTypes = filteredResTypes
				}
			}

			// generate iam instance resource by supported general cache resource types, do pagination
			resLen := int64(len(resTypes))
			if page.Offset >= resLen {
				return &types.ListInstanceResult{Count: 0, Results: make([]types.InstanceResource, 0)}, nil
			}

			end := page.Offset + page.Limit
			if end > resLen {
				end = resLen
			}

			res := make([]types.InstanceResource, 0)
			for _, resType := range resTypes[page.Offset:end] {
				res = append(res, types.InstanceResource{
					ID:          resType,
					DisplayName: resType,
				})
			}

			return &types.ListInstanceResult{Count: resLen, Results: res}, nil
		},
		FetchInstanceInfo: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.FetchInstanceInfoFilter) (
			[]map[string]interface{}, error) {

			// only support query name field, name field is the same with the id field
			hasNameField := false
			for _, attr := range filter.Attrs {
				if attr == types.NameField {
					hasNameField = true
					break
				}
			}
			if !hasNameField {
				return make([]map[string]interface{}, 0), nil
			}

			res := make([]map[string]interface{}, 0)
			for _, id := range filter.IDs {
				_, exists := general.SupportedResTypeMap[general.ResType(id)]
				if exists {
					res = append(res, map[string]interface{}{types.NameField: id})
				}
			}

			return res, nil
		},
		ListInstanceByPolicy: func(kit *rest.Kit, resourceType iamtypes.TypeID,
			filter *types.ListInstanceByPolicyFilter,
			page types.Page) (*types.ListInstanceResult, error) {
			return nil, fmt.Errorf("%s do not support %s", iamtypes.GeneralCache, types.ListInstanceByPolicyMethod)
		},
	}, nil
}

func genTenantSetMethod(kit *rest.Kit, lgc *logics.Logics) (types.ResourcePullMethod, error) {
	return types.ResourcePullMethod{
		ListInstance: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.ListInstanceFilter,
			page types.Page) (*types.ListInstanceResult, error) {

			if page.Offset > 0 {
				return &types.ListInstanceResult{Count: 0, Results: make([]types.InstanceResource, 0)}, nil
			}

			if filter != nil {
				if filter.Parent != nil {
					return &types.ListInstanceResult{Count: 0, Results: make([]types.InstanceResource, 0)}, nil
				}

				if len(filter.Keyword) != 0 {
					if !strings.Contains(strings.ToLower(tenantset.DefaultTenantSetName),
						strings.ToLower(filter.Keyword)) {
						return &types.ListInstanceResult{Count: 0, Results: make([]types.InstanceResource, 0)}, nil
					}
				}
			}

			return &types.ListInstanceResult{Count: 1, Results: []types.InstanceResource{{
				ID:          strconv.FormatInt(tenantset.DefaultTenantSetID, 10),
				DisplayName: tenantset.DefaultTenantSetName,
			}}}, nil
		},
		FetchInstanceInfo: func(kit *rest.Kit, resourceType iamtypes.TypeID, filter *types.FetchInstanceInfoFilter) (
			[]map[string]interface{}, error) {

			res := make([]map[string]interface{}, 0)
			for _, id := range filter.IDs {
				if id == strconv.FormatInt(tenantset.DefaultTenantSetID, 10) {
					res = append(res, map[string]interface{}{
						types.IDField:   id,
						types.NameField: tenantset.DefaultTenantSetName,
					})
					break
				}
			}
			return res, nil
		},
		ListInstanceByPolicy: func(kit *rest.Kit, resourceType iamtypes.TypeID,
			filter *types.ListInstanceByPolicyFilter,
			page types.Page) (*types.ListInstanceResult, error) {
			return nil, fmt.Errorf("%s do not support %s", iamtypes.TenantSet, types.ListInstanceByPolicyMethod)
		},
	}, nil
}
