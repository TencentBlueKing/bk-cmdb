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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/kube/types"
)

// convertKubeCondition generate different query conditions based on different resources.
func (s *Service) findKubeTopoPathIfo(kit *rest.Kit, option *types.KubeTopoPathReq, filter mapstr.MapStr,
	tableNames []string) (*types.KubeTopoPathRsp, error) {

	result := &types.KubeTopoPathRsp{Info: make([]types.KubeObjectInfo, 0)}

	query := &metadata.QueryCondition{
		Condition:      filter,
		Page:           option.Page,
		Fields:         []string{types.BKIDField, types.KubeNameField},
		DisableCounter: true}

	// according to the topology display, put the folder to the front
	if tableNames[0] == types.BKTableNameBaseNamespace {
		result.Info = append(result.Info, types.KubeObjectInfo{
			ID: types.KubeFolderID, Name: types.KubeFolderName, Kind: types.KubeFolder,
		})
	}

	for _, tableName := range tableNames {
		switch tableName {
		case types.BKTableNameBaseCluster:
			clusters, err := s.Logics.ContainerOperation().SearchCluster(kit, query)
			if err != nil {
				blog.Errorf("search cluster failed, err: %v, rid: %s", err, kit.Rid)
				return result, err
			}
			for _, cluster := range clusters.Data {
				result.Info = append(result.Info, types.KubeObjectInfo{
					ID: cluster.ID, Name: *cluster.Name, Kind: types.KubeCluster,
				})
			}
		case types.BKTableNameBaseNamespace:

			option := &types.QueryReq{
				Table:     types.BKTableNameBaseNamespace,
				Condition: query,
			}
			namespaces, err := s.Engine.CoreAPI.CoreService().Kube().FindInst(kit.Ctx, kit.Header, option)
			if err != nil {
				blog.Errorf("find namespace failed, cond: %v, err: %v, rid: %s", query, err, kit.Rid)
				return result, err
			}
			for _, namespace := range namespaces.Info {
				id, err := util.GetInt64ByInterface(namespace[types.BKIDField])
				if err != nil {
					blog.Errorf("find namespace failed, cond: %v, err: %v, rid: %s", query, err, kit.Rid)
					return result, err
				}
				result.Info = append(result.Info, types.KubeObjectInfo{
					ID: id, Name: util.GetStrByInterface(namespace[types.KubeNameField]), Kind: types.KubeNamespace,
				})
			}
		default:

			option := &types.QueryReq{
				Table:     tableName,
				Condition: query,
			}
			workloads, cErr := s.Engine.CoreAPI.CoreService().Kube().FindInst(kit.Ctx, kit.Header, option)
			if cErr != nil {
				blog.Errorf("find namespace failed, cond: %v, err: %v, rid: %s", query, cErr, kit.Rid)
				return result, cErr
			}

			kind, err := types.GetKindByWorkLoadTableNameMap(tableName)
			if err != nil {
				return result, err
			}

			for _, workload := range workloads.Info {
				id, err := util.GetInt64ByInterface(workload[types.BKIDField])
				if err != nil {
					blog.Errorf("find namespace failed, cond: %v, err: %v, rid: %s", query, err, kit.Rid)
					return result, err
				}
				result.Info = append(result.Info, types.KubeObjectInfo{
					ID: id, Name: util.GetStrByInterface(workload[types.KubeNameField]), Kind: kind[tableName],
				})
			}
		}
	}

	return result, nil
}

func combinationConditions(infos []types.KubeResourceInfo, bizID int64,
	supplierAccount string) []map[string]interface{} {

	filters := make([]map[string]interface{}, 0)
	for _, info := range infos {
		switch info.Kind {
		case types.KubeFolder:
			filters = append(filters, map[string]interface{}{
				types.BKClusterIDFiled:       info.ID,
				types.HasPodField:            false,
				types.BKBizIDField:           bizID,
				types.BKSupplierAccountField: supplierAccount,
			})

		case types.KubeCluster:
			filters = append(filters, map[string]interface{}{
				types.BKClusterIDFiled:       info.ID,
				types.BKBizIDField:           bizID,
				types.BKSupplierAccountField: supplierAccount,
			})

		case types.KubeNamespace:
			filters = append(filters, map[string]interface{}{
				types.BKNamespaceIDField:     info.ID,
				types.BKBizIDField:           bizID,
				types.BKSupplierAccountField: supplierAccount,
			})

		default:
			filters = append(filters, map[string]interface{}{
				types.RefIDField:             info.ID,
				types.RefKindField:           info.Kind,
				types.BKBizIDField:           bizID,
				types.BKSupplierAccountField: supplierAccount,
			})
		}
	}
	return filters
}

// CountKubeTopoHostsOrPods count the number of node pods or hosts
func (s *Service) CountKubeTopoHostsOrPods(ctx *rest.Contexts) {

	option := new(types.KubeTopoCountOption)
	if err := ctx.DecodeInto(option); err != nil {
		blog.Errorf("failed to parse the params, error: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	kind := ctx.Request.PathParameter("type")
	if kind != types.KubeHostKind && kind != types.KubePodKind {
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
	result := make([]types.KubeTopoCountRsp, 0)
	if kind == types.KubePodKind {
		podFilters := make([]map[string]interface{}, 0)

		resIDMap := make(map[int]struct{})
		for id, filter := range filters {
			// if the filter contains the "has_pod" field, it indicates the folder node
			if _, ok := filter[types.HasPodField]; ok {
				resIDMap[id] = struct{}{}
				continue
			}
			podFilters = append(podFilters, filter)
		}
		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			types.BKTableNameBasePod, podFilters)
		if err != nil {
			blog.Errorf("count pod failed, cond: %#v, err: %v, rid: %s", podFilters, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		var idx int
		for id := range option.ResourceInfos {
			if _, ok := resIDMap[id]; ok {
				result = append(result, types.KubeTopoCountRsp{
					Kind:  option.ResourceInfos[id].Kind,
					ID:    option.ResourceInfos[id].ID,
					Count: 0,
				})
				continue
			}

			result = append(result, types.KubeTopoCountRsp{
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

func (s *Service) getTopoHostNumber(ctx *rest.Contexts, resourceInfos []types.KubeResourceInfo,
	filters []map[string]interface{}, bizID int64) ([]types.KubeTopoCountRsp, error) {

	// obtaining a host requires the following steps:
	// 1、get all hostIDs of the node.
	// 2、deduplicate hostID.
	// 3、combine the hostID and business ID to check the modulehostconfig table,
	// and the final number is the real number of hosts.
	result := make([]types.KubeTopoCountRsp, 0)

	for id, filter := range filters {

		// determine whether this node is a folder If it is a folder, then you need to check the node table.
		if resourceInfos[id].Kind == types.KubeFolder {
			count, err := s.getHostIDsByCond(ctx.Kit, filter, types.BKTableNameBaseNode, bizID)
			if err != nil {
				return nil, err
			}
			result = append(result, types.KubeTopoCountRsp{
				Kind:  resourceInfos[id].Kind,
				ID:    resourceInfos[id].ID,
				Count: count,
			})
			continue
		}

		// what counts here is the number of hosts in the pod table excluding folders.
		count, err := s.getHostIDsByCond(ctx.Kit, filter, types.BKTableNameBasePod, bizID)
		if err != nil {
			return nil, err
		}

		if types.IsWorkLoadKind(util.GetStrByInterface(filter[types.RefKindField])) {
			id, _ := util.GetInt64ByInterface(filter[types.RefIDField])
			result = append(result, types.KubeTopoCountRsp{
				Kind:  util.GetStrByInterface(filter[types.RefKindField]),
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
		if clusterID, ok := filter[types.BKClusterIDFiled]; ok {
			nodeFilter := mapstr.MapStr{
				types.BKClusterIDFiled:       clusterID,
				types.HasPodField:            false,
				types.BKBizIDField:           bizID,
				types.BKSupplierAccountField: ctx.Kit.SupplierAccount,
			}
			folderHostCount, err = s.getHostIDsByCond(ctx.Kit, nodeFilter, types.BKTableNameBaseNode, bizID)
			if err != nil {
				return nil, err
			}
		}

		result = append(result, types.KubeTopoCountRsp{
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
	option := &types.QueryReq{
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

	option := new(types.KubeTopoPathReq)
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
	subObject, filter := types.GetKubeSubTopoObject(option.ReferenceObjID, option.ReferenceID, bizID)
	tableNames, err := types.GetCollectionWithObject(subObject)
	if err != nil {
		blog.Errorf("failed get , err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	filter[types.BKSupplierAccountField] = ctx.Kit.SupplierAccount

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
			if tableName == types.BKTableNameBaseNamespace {
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
	if !types.IsContainerTopoResource(object) {
		blog.Errorf("the parameter is invalid and does not belong to the container object(%s)", object)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "object"))
		return
	}

	result := make([]types.KubeAttrsRsp, 0)
	switch object {
	case types.KubeCluster:
		for _, descriptor := range types.ClusterSpecFieldsDescriptor {
			result = append(result, types.KubeAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types.KubeNamespace:
		for _, descriptor := range types.NamespaceSpecFieldsDescriptor {
			result = append(result, types.KubeAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types.KubeNode:
		for _, descriptor := range types.NodeSpecFieldsDescriptor {
			result = append(result, types.KubeAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types.KubeWorkload:
		for _, descriptor := range types.WorkLoadSpecFieldsDescriptor {
			result = append(result, types.KubeAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types.KubePod:
		for _, descriptor := range types.PodSpecFieldsDescriptor {
			result = append(result, types.KubeAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types.KubeContainer:
		for _, descriptor := range types.ContainerSpecFieldsDescriptor {
			result = append(result, types.KubeAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	}
	ctx.RespEntity(result)
}
