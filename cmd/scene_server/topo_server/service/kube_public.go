/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"strconv"

	"configcenter/pkg/blog"
	"configcenter/pkg/common"
	"configcenter/pkg/http/rest"
	types2 "configcenter/pkg/kube/types"
	"configcenter/pkg/mapstr"
	"configcenter/pkg/metadata"
	"configcenter/pkg/util"
)

// convertKubeCondition generate different query conditions based on different resources.
func (s *Service) findKubeTopoPathIfo(kit *rest.Kit, option *types2.KubeTopoPathOption, filter mapstr.MapStr,
	tableNames []string) (*types2.KubeTopoPathRsp, error) {

	result := &types2.KubeTopoPathRsp{Info: make([]types2.KubeObjectInfo, 0)}

	query := &metadata.QueryCondition{
		Condition:      filter,
		Page:           option.Page,
		Fields:         []string{types2.BKIDField, types2.KubeNameField},
		DisableCounter: true}

	// according to the topology display, put the folder to the front
	if tableNames[0] == types2.BKTableNameBaseNamespace {
		result.Info = append(result.Info, types2.KubeObjectInfo{
			ID: types2.KubeFolderID, Name: types2.KubeFolderName, Kind: types2.KubeFolder,
		})
	}

	for _, tableName := range tableNames {
		switch tableName {
		case types2.BKTableNameBaseCluster:
			clusters, err := s.Engine.CoreAPI.CoreService().Container().SearchCluster(kit.Ctx, kit.Header, query)
			if err != nil {
				blog.Errorf("search cluster failed, err: %v, rid: %s", err, kit.Rid)
				return result, err
			}
			for _, cluster := range clusters.Data {
				result.Info = append(result.Info, types2.KubeObjectInfo{
					ID: cluster.ID, Name: *cluster.Name, Kind: types2.KubeCluster,
				})
			}
		case types2.BKTableNameBaseNamespace:

			option := &types2.QueryReq{
				Table:     types2.BKTableNameBaseNamespace,
				Condition: query,
			}
			namespaces, err := s.Engine.CoreAPI.CoreService().Kube().FindInst(kit.Ctx, kit.Header, option)
			if err != nil {
				blog.Errorf("find namespace failed, cond: %v, err: %v, rid: %s", query, err, kit.Rid)
				return result, err
			}
			for _, namespace := range namespaces.Info {
				id, err := util.GetInt64ByInterface(namespace[types2.BKIDField])
				if err != nil {
					blog.Errorf("find namespace failed, cond: %v, err: %v, rid: %s", query, err, kit.Rid)
					return result, err
				}
				result.Info = append(result.Info, types2.KubeObjectInfo{
					ID: id, Name: util.GetStrByInterface(namespace[types2.KubeNameField]), Kind: types2.KubeNamespace,
				})
			}
		default:

			option := &types2.QueryReq{
				Table:     tableName,
				Condition: query,
			}
			workloads, cErr := s.Engine.CoreAPI.CoreService().Kube().FindInst(kit.Ctx, kit.Header, option)
			if cErr != nil {
				blog.Errorf("find namespace failed, cond: %v, err: %v, rid: %s", query, cErr, kit.Rid)
				return result, cErr
			}

			kind, err := types2.GetKindByWorkLoadTableNameMap(tableName)
			if err != nil {
				return result, err
			}

			for _, workload := range workloads.Info {
				id, err := util.GetInt64ByInterface(workload[types2.BKIDField])
				if err != nil {
					blog.Errorf("find namespace failed, cond: %v, err: %v, rid: %s", query, err, kit.Rid)
					return result, err
				}
				result.Info = append(result.Info, types2.KubeObjectInfo{
					ID: id, Name: util.GetStrByInterface(workload[types2.KubeNameField]), Kind: kind[tableName],
				})
			}
		}
	}

	return result, nil
}

func combinationConditions(infos []types2.KubeResourceInfo, bizID int64,
	supplierAccount string) []map[string]interface{} {

	filters := make([]map[string]interface{}, 0)
	for _, info := range infos {
		switch info.Kind {
		case types2.KubeFolder:
			filters = append(filters, map[string]interface{}{
				types2.BKClusterIDFiled:       info.ID,
				types2.HasPodField:            false,
				types2.BKBizIDField:           bizID,
				types2.BKSupplierAccountField: supplierAccount,
			})

		case types2.KubeCluster:
			filters = append(filters, map[string]interface{}{
				types2.BKClusterIDFiled:       info.ID,
				types2.BKBizIDField:           bizID,
				types2.BKSupplierAccountField: supplierAccount,
			})

		case types2.KubeNamespace:
			filters = append(filters, map[string]interface{}{
				types2.BKNamespaceIDField:     info.ID,
				types2.BKBizIDField:           bizID,
				types2.BKSupplierAccountField: supplierAccount,
			})

		default:
			filters = append(filters, map[string]interface{}{
				types2.RefIDField:             info.ID,
				types2.RefKindField:           info.Kind,
				types2.BKBizIDField:           bizID,
				types2.BKSupplierAccountField: supplierAccount,
			})
		}
	}
	return filters
}

// CountKubeTopoHostsOrPods count the number of node pods or hosts
func (s *Service) CountKubeTopoHostsOrPods(ctx *rest.Contexts) {

	option := new(types2.KubeTopoCountOption)
	if err := ctx.DecodeInto(option); err != nil {
		blog.Errorf("failed to parse the params, error: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	kind := ctx.Request.PathParameter("type")
	if kind != types2.KubeHostKind && kind != types2.KubePodKind {
		blog.Errorf("failed to parse the params, error: %v, rid: %s", ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	if cErr := option.Validate(); cErr.ErrCode != 0 {
		blog.Errorf("validate request failed, err: %v, rid: %s", cErr, ctx.Kit.Rid)
		ctx.RespAutoError(cErr.ToCCError(ctx.Kit.CCError))
		return
	}
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	filters := combinationConditions(option.ResourceInfos, bizID, ctx.Kit.SupplierAccount)
	result := make([]types2.KubeTopoCountRsp, 0)
	if kind == types2.KubePodKind {
		podFilters := make([]map[string]interface{}, 0)

		resIDMap := make(map[int]struct{})
		for id, filter := range filters {
			// if the filter contains the "has_pod" field, it indicates the folder node
			if _, ok := filter[types2.HasPodField]; ok {
				resIDMap[id] = struct{}{}
				continue
			}
			podFilters = append(podFilters, filter)
		}
		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			types2.BKTableNameBasePod, podFilters)
		if err != nil {
			blog.Errorf("count pod failed, cond: %#v, err: %v, rid: %s", podFilters, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		var idx int
		for id := range option.ResourceInfos {
			if _, ok := resIDMap[id]; ok {
				result = append(result, types2.KubeTopoCountRsp{
					Kind:  option.ResourceInfos[id].Kind,
					ID:    option.ResourceInfos[id].ID,
					Count: 0,
				})
				continue
			}

			result = append(result, types2.KubeTopoCountRsp{
				Kind:  option.ResourceInfos[id].Kind,
				ID:    option.ResourceInfos[id].ID,
				Count: counts[idx],
			})
			idx++
		}

		ctx.RespEntity(result)
		return
	}

	result, err = s.getTopoHostNumber(ctx, option.ResourceInfos, filters, bizID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *Service) getTopoHostNumber(ctx *rest.Contexts, resourceInfos []types2.KubeResourceInfo,
	filters []map[string]interface{}, bizID int64) ([]types2.KubeTopoCountRsp, error) {

	// obtaining a host requires the following steps:
	// 1、get all hostIDs of the node.
	// 2、deduplicate hostID.
	// 3、combine the hostID and business ID to check the modulehostconfig table,
	// and the final number is the real number of hosts.
	result := make([]types2.KubeTopoCountRsp, 0)

	for id, filter := range filters {

		// determine whether this node is a folder If it is a folder, then you need to check the node table.
		if resourceInfos[id].Kind == types2.KubeFolder {
			count, err := s.getHostIDsByCond(ctx.Kit, filter, types2.BKTableNameBaseNode, bizID)
			if err != nil {
				return nil, err
			}
			result = append(result, types2.KubeTopoCountRsp{
				Kind:  resourceInfos[id].Kind,
				ID:    resourceInfos[id].ID,
				Count: count,
			})
			continue
		}

		// what counts here is the number of hosts in the pod table excluding folders.
		count, err := s.getHostIDsByCond(ctx.Kit, filter, types2.BKTableNameBasePod, bizID)
		if err != nil {
			return nil, err
		}

		if types2.IsWorkLoadKind(util.GetStrByInterface(filter[types2.RefKindField])) {
			id, _ := util.GetInt64ByInterface(filter[types2.RefIDField])
			result = append(result, types2.KubeTopoCountRsp{
				Kind:  util.GetStrByInterface(filter[types2.RefKindField]),
				ID:    id,
				Count: count,
			})
			continue
		}

		var (
			folderHostCount int64
		)
		// for the calculation of the number of hosts under the cluster,
		// it is necessary to add the number of hosts under the folder node under the cluster.
		if clusterID, ok := filter[types2.BKClusterIDFiled]; ok {
			nodeFilter := mapstr.MapStr{
				types2.BKClusterIDFiled:       clusterID,
				types2.HasPodField:            false,
				types2.BKBizIDField:           bizID,
				types2.BKSupplierAccountField: ctx.Kit.SupplierAccount,
			}
			folderHostCount, err = s.getHostIDsByCond(ctx.Kit, nodeFilter, types2.BKTableNameBaseNode, bizID)
			if err != nil {
				return nil, err
			}
		}

		result = append(result, types2.KubeTopoCountRsp{
			Kind:  resourceInfos[id].Kind,
			ID:    resourceInfos[id].ID,
			Count: count + folderHostCount,
		})
	}
	return result, nil
}

func (s *Service) getHostIDsByCond(kit *rest.Kit, cond mapstr.MapStr, table string, bizID int64) (int64, error) {

	query := &metadata.QueryCondition{
		Condition: cond,
		Fields:    []string{common.BKHostIDField},
	}
	option := &types2.QueryReq{
		Table:     table,
		Condition: query,
	}
	var err error
	insts, err := s.Engine.CoreAPI.CoreService().Kube().FindInst(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("find inst failed, cond: %v, err: %v, rid: %s", query, err, kit.Rid)
		return 0, err
	}

	hostIDMap := make(map[int64]struct{})
	for _, inst := range insts.Info {
		hostID, err := util.GetInt64ByInterface(inst[common.BKHostIDField])
		if err != nil {
			blog.Errorf("get inst attribute failed, attr: %s, node: %v, err: %v, rid: %s", common.BKHostIDField, inst,
				err, kit.Rid)
			return 0, err
		}
		hostIDMap[hostID] = struct{}{}
	}

	hostIDs := make([]int64, 0)

	for id := range hostIDMap {
		hostIDs = append(hostIDs, id)
	}

	countOp := []map[string]interface{}{{
		common.BKAppIDField: bizID,
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDs,
		},
	}}

	rsp, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameModuleHostConfig, countOp)
	if err != nil {
		blog.Errorf("get host module relation failed, err: %v, rid: %s", err, kit.Rid)
		return 0, err
	}

	return rsp[0], nil

}

// SearchKubeTopoPath querying container topology paths.
func (s *Service) SearchKubeTopoPath(ctx *rest.Contexts) {

	option := new(types2.KubeTopoPathOption)
	if err := ctx.DecodeInto(option); err != nil {
		blog.Errorf("failed to parse the params, error: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	if cErr := option.Validate(); cErr.ErrCode != 0 {
		blog.Errorf("validate request failed, err: %v, rid: %s", cErr, ctx.Kit.Rid)
		ctx.RespAutoError(cErr.ToCCError(ctx.Kit.CCError))
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// get the next level resource object.
	subObject, filter := types2.GetKubeSubTopoObject(option.ReferenceObjID, option.ReferenceID, bizID)
	tableNames, err := types2.GetCollectionWithObject(subObject)
	if err != nil {
		blog.Errorf("failed get , err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	filter[types2.BKSupplierAccountField] = ctx.Kit.SupplierAccount

	if option.Page.EnableCount {
		var count int64
		for _, tableName := range tableNames {
			filter := []map[string]interface{}{filter}
			counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
				tableName, filter)
			if err != nil {
				blog.Errorf("count node failed, err: %v, cond: %#v, rid: %s", err, filter, ctx.Kit.Rid)
				ctx.RespAutoError(err)
				return
			}
			// for the next-level topology of the cluster, a folder needs to be added in addition to the namespace.
			if tableName == types2.BKTableNameBaseNamespace {
				counts[0] += 1
			}
			count += counts[0]
		}

		ctx.RespEntityWithCount(count, make([]mapstr.MapStr, 0))
		return
	}

	result, err := s.findKubeTopoPathIfo(ctx.Kit, option, filter, tableNames)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// FindResourceAttrs get the attribute information of the container object
func (s *Service) FindResourceAttrs(ctx *rest.Contexts) {

	object := ctx.Request.PathParameter("object")
	if !types2.IsContainerTopoResource(object) {
		blog.Errorf("the parameter is invalid and does not belong to the container object(%s)", object)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "object"))
		return
	}

	result := make([]types2.KubeAttrsRsp, 0)
	switch object {
	case types2.KubeCluster:
		for _, descriptor := range types2.ClusterSpecFieldsDescriptor {
			result = append(result, types2.KubeAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types2.KubeNamespace:
		for _, descriptor := range types2.NamespaceSpecFieldsDescriptor {
			result = append(result, types2.KubeAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types2.KubeNode:
		for _, descriptor := range types2.NodeSpecFieldsDescriptor {
			result = append(result, types2.KubeAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types2.KubeWorkload:
		for _, descriptor := range types2.WorkLoadSpecFieldsDescriptor {
			result = append(result, types2.KubeAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types2.KubePod:
		for _, descriptor := range types2.PodSpecFieldsDescriptor {
			result = append(result, types2.KubeAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types2.KubeContainer:
		for _, descriptor := range types2.ContainerSpecFieldsDescriptor {
			result = append(result, types2.KubeAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	}
	ctx.RespEntity(result)
}
