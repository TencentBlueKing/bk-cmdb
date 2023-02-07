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
	"errors"
	"strconv"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/kube/types"
)

// findKubeTopoPathInfo generate different query conditions based on different resources.
func (s *Service) findKubeTopoPathInfo(kit *rest.Kit, option *types.KubeTopoPathOption, filter mapstr.MapStr,
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
			clusters, err := s.Engine.CoreAPI.CoreService().Kube().SearchCluster(kit.Ctx, kit.Header, query)
			if err != nil {
				blog.Errorf("search cluster failed, err: %v, rid: %s", err, kit.Rid)
				return result, err
			}
			for _, cluster := range clusters.Data {
				if cluster.Name == nil {
					blog.Errorf("cluster name is nil, cluster: %v, rid: %s", cluster, err, kit.Rid)
					return nil, errors.New("cluster name is nil")
				}
				result.Info = append(result.Info, types.KubeObjectInfo{
					ID: cluster.ID, Name: *cluster.Name, Kind: types.KubeCluster,
				})
			}
		case types.BKTableNameBaseNamespace:

			namespaces, err := s.Engine.CoreAPI.CoreService().Kube().ListNamespace(kit.Ctx, kit.Header, query)
			if err != nil {
				blog.Errorf("find namespace failed, cond: %v, err: %v, rid: %s", query, err, kit.Rid)
				return result, err
			}
			for _, namespace := range namespaces.Data {
				result.Info = append(result.Info, types.KubeObjectInfo{
					ID:   namespace.ID,
					Name: namespace.Name,
					Kind: types.KubeNamespace,
				})
			}
		default:

			kind, err := types.GetKindByWorkLoadTableNameMap(tableName)
			if err != nil {
				return result, err
			}
			workloads, cErr := s.Engine.CoreAPI.CoreService().Kube().ListWorkload(kit.Ctx, kit.Header, query,
				types.WorkloadType(kind[tableName]))
			if cErr != nil {
				blog.Errorf("find namespace failed, cond: %v, err: %v, rid: %s", query, cErr, kit.Rid)
				return result, cErr
			}
			for _, workload := range workloads.Info {
				result.Info = append(result.Info, types.KubeObjectInfo{
					ID:   workload.GetWorkloadBase().ID,
					Name: workload.GetWorkloadBase().Name,
					Kind: kind[tableName],
				})
			}
		}
	}

	return result, nil
}

func combinationConditions(infos []types.KubeResourceInfo, bizID int64) []map[string]interface{} {

	filters := make([]map[string]interface{}, 0)
	for _, info := range infos {
		switch info.Kind {
		case types.KubeFolder:
			filters = append(filters, map[string]interface{}{
				types.BKClusterIDFiled: info.ID,
				types.HasPodField:      false,
				types.BKBizIDField:     bizID,
			})

		case types.KubeCluster:
			filters = append(filters, map[string]interface{}{
				types.BKClusterIDFiled: info.ID,
				types.BKBizIDField:     bizID,
			})

		case types.KubeNamespace:
			filters = append(filters, map[string]interface{}{
				types.BKNamespaceIDField: info.ID,
				types.BKBizIDField:       bizID,
			})
		default:
			filters = append(filters, map[string]interface{}{
				types.RefIDField:   info.ID,
				types.RefKindField: info.Kind,
				types.BKBizIDField: bizID,
			})
		}
	}
	return filters
}

func (s *Service) countKubeHostOrPodsByCond(kit *rest.Kit, option *types.KubeTopoCountOption, bizID int64,
	kind string) ([]types.KubeTopoCountRsp, error) {

	filters := combinationConditions(option.ResourceInfos, bizID)

	switch kind {
	case types.KubeHostKind:
		result, err := s.getTopoHostNumber(kit, option.ResourceInfos, filters, bizID)
		if err != nil {
			blog.Errorf("get host number failed, option: %+v, bizID: %d, err: %v, rid: %s", option, bizID, err)
			return nil, err
		}
		return result, nil
	case types.KubePodKind:

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
		result := make([]types.KubeTopoCountRsp, 0)

		// here it is explained that the node to be queried is a folder,
		// then the number of pods must be 0.
		if len(podFilters) == 0 {
			return result, nil
		}

		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
			types.BKTableNameBasePod, podFilters)
		if err != nil {
			blog.Errorf("count pod failed, cond: %#v, err: %v, rid: %s", podFilters, err, kit.Rid)
			return nil, err
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
		return result, nil

	default:
		return nil, errors.New("count type error")
	}
}

// CountKubeTopoHostsOrPods count the number of node pods or hosts
func (s *Service) CountKubeTopoHostsOrPods(ctx *rest.Contexts) {

	option := new(types.KubeTopoCountOption)
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

	kind := ctx.Request.PathParameter("type")
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result, err := s.countKubeHostOrPodsByCond(ctx.Kit, option, bizID, kind)
	if err != nil {
		blog.Errorf("failed to get(%s) number, bizID: %d, option: %+v, err: %v, rid: %s",
			kind, bizID, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *Service) getTopoHostNumber(kit *rest.Kit, resourceInfos []types.KubeResourceInfo,
	filters []map[string]interface{}, bizID int64) ([]types.KubeTopoCountRsp, error) {

	// obtaining a host requires the following steps: 1、get all hostIDs of the node. 2、deduplicate hostID. 3、combine
	// the hostID and business ID to check the modulehostconfig table, and the final number is the real number of hosts.
	var wg sync.WaitGroup
	var lock sync.RWMutex
	var firstErr error
	pipeline, result := make(chan bool, 20), make([]types.KubeTopoCountRsp, 0)

	for id, filter := range filters {
		pipeline <- true
		wg.Add(1)
		go func(id int, filter map[string]interface{}) {
			defer func() {
				wg.Done()
				<-pipeline
			}()
			// determine whether this node is a folder If it is a folder, then you need to check the node table.
			if resourceInfos[id].Kind == types.KubeFolder {
				hostMap, err := s.getHostIDsInNodeByCond(kit, bizID, filter)
				if err != nil {
					firstErr = err
					return
				}
				lock.Lock()
				result = append(result, types.KubeTopoCountRsp{
					Kind:  resourceInfos[id].Kind,
					ID:    resourceInfos[id].ID,
					Count: int64(len(hostMap)),
				})
				lock.Unlock()
				return
			}

			// what counts here is the number of hosts in the pod table excluding folders.
			hostMap, err := s.getHostIDsInPodsByCond(kit, bizID, filter)
			if err != nil {
				firstErr = err
				return
			}

			workloadType := types.WorkloadType(util.GetStrByInterface(filter[types.RefKindField]))
			// the scenario dealt with here is the workload type calculation number.
			if err := workloadType.Validate(); err == nil {
				id, err := util.GetInt64ByInterface(filter[types.RefIDField])
				if err != nil {
					firstErr = err
					return
				}

				lock.Lock()
				result = append(result, types.KubeTopoCountRsp{
					Kind:  util.GetStrByInterface(filter[types.RefKindField]),
					ID:    id,
					Count: int64(len(hostMap)),
				})
				lock.Unlock()
				return
			}

			resultHostMap, err := s.getClusterNumFromFolder(kit, bizID, filter)
			if err != nil {
				firstErr = err
				return
			}

			for id := range hostMap {
				resultHostMap[id] = struct{}{}
			}

			lock.Lock()
			result = append(result, types.KubeTopoCountRsp{
				Kind:  resourceInfos[id].Kind,
				ID:    resourceInfos[id].ID,
				Count: int64(len(resultHostMap)),
			})
			lock.Unlock()
		}(id, filter)
	}

	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}
	return result, nil
}

// getClusterNumFromFolder for the calculation of the number of hosts under the cluster,
// it is necessary to add the number of hosts under the folder node under the cluster.
func (s *Service) getClusterNumFromFolder(kit *rest.Kit, bizID int64, filter map[string]interface{}) (
	map[int64]struct{}, error) {

	result := make(map[int64]struct{})
	clusterID, ok := filter[types.BKClusterIDFiled]
	if !ok {
		return result, nil
	}

	nodeFilter := mapstr.MapStr{
		types.BKClusterIDFiled: clusterID,
		types.HasPodField:      false,
		types.BKBizIDField:     bizID,
	}
	result, err := s.getHostIDsInNodeByCond(kit, bizID, nodeFilter)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *Service) getDistinctHostMap(kit *rest.Kit, bizID int64, hostIDMap map[int64]struct{}) (
	map[int64]struct{}, error) {

	result := make(map[int64]struct{}, 0)
	if len(hostIDMap) == 0 {
		return result, nil
	}

	hostIDs := make([]int64, 0)
	for id := range hostIDMap {
		hostIDs = append(hostIDs, id)
	}

	relationReq := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		HostIDArr:     hostIDs,
		Page:          metadata.BasePage{Limit: common.BKNoLimit},
		Fields:        []string{common.BKHostIDField},
	}

	hostRelations, err := s.Engine.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, relationReq)
	if err != nil {
		return nil, err
	}
	for _, info := range hostRelations.Info {
		result[info.HostID] = struct{}{}
	}
	return result, nil
}

func (s *Service) getHostIDsInNodeByCond(kit *rest.Kit, bizID int64, cond mapstr.MapStr) (map[int64]struct{}, error) {

	query := &metadata.QueryCondition{
		Condition:      cond,
		Fields:         []string{common.BKHostIDField},
		DisableCounter: true,
	}
	nodes, cErr := s.Engine.CoreAPI.CoreService().Kube().SearchNode(kit.Ctx, kit.Header, query)
	if cErr != nil {
		blog.Errorf("find nodes failed, cond: %v, err: %v, rid: %s", query, cErr, kit.Rid)
		return nil, cErr
	}

	hostIDMap := make(map[int64]struct{})
	for _, node := range nodes.Data {
		hostIDMap[node.HostID] = struct{}{}
	}

	result, err := s.getDistinctHostMap(kit, bizID, hostIDMap)
	if err != nil {
		blog.Errorf("get host module relation failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	return result, nil
}

func (s *Service) getHostIDsInPodsByCond(kit *rest.Kit, bizID int64, cond mapstr.MapStr) (map[int64]struct{}, error) {

	query := &metadata.QueryCondition{
		Condition:      cond,
		Fields:         []string{common.BKHostIDField},
		DisableCounter: true,
	}

	pods, cErr := s.Engine.CoreAPI.CoreService().Kube().ListPod(kit.Ctx, kit.Header, query)
	if cErr != nil {
		blog.Errorf("find pods failed, cond: %v, err: %v, rid: %s", query, cErr, kit.Rid)
		return nil, cErr
	}
	hostIDMap := make(map[int64]struct{})
	for _, pod := range pods.Info {
		hostIDMap[pod.HostID] = struct{}{}
	}

	result, err := s.getDistinctHostMap(kit, bizID, hostIDMap)
	if err != nil {
		blog.Errorf("get host module relation failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	return result, nil
}

// SearchKubeTopoPath querying container topology paths.
func (s *Service) SearchKubeTopoPath(ctx *rest.Contexts) {

	option := new(types.KubeTopoPathOption)
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
	if filter == nil {
		blog.Errorf("failed to get subObject, option: %+v, err: %v, rid: %s", option, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsInvalid))
		return
	}
	tableNames, err := types.GetCollectionWithObject(subObject)
	if err != nil {
		blog.Errorf("failed get tableName, subObject: %s, err: %v, rid: %s", subObject, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

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

	result, err := s.findKubeTopoPathInfo(ctx.Kit, option, filter, tableNames)
	if err != nil {
		blog.Errorf("failed to get topo path, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// FindResourceAttrs get the attribute information of the kube object
// for front-end use only
func (s *Service) FindResourceAttrs(ctx *rest.Contexts) {

	object := ctx.Request.PathParameter("object")
	if !types.IsKubeTopoResource(object) {
		blog.Errorf("the param is invalid and does not belong to the kube object(%s)", object)
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
		// TODO compatible for different types of workload
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

func (s *Service) hasNextLevelResource(kit *rest.Kit, kind string, bizID int64, ids []int64) (bool, error) {
	var tables []string
	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
	}

	switch kind {
	case types.KubeCluster:
		tables = []string{types.BKTableNameBaseNamespace, types.BKTableNameBaseNode, types.BKTableNameBasePod}
		workLoads := types.GetWorkLoadTables()
		tables = append(tables, workLoads...)
		filter[types.BKClusterIDFiled] = map[string]interface{}{common.BKDBIN: ids}

	case types.KubeNamespace:
		tables = []string{types.BKTableNameBasePod}
		workLoads := types.GetWorkLoadTables()
		tables = append(tables, workLoads...)
		filter[types.BKNamespaceIDField] = map[string]interface{}{common.BKDBIN: ids}

	default:
		tables = []string{types.BKTableNameBasePod}
		if err := types.WorkloadType(kind).Validate(); err != nil {
			return false, err
		}

		filter[types.RefKindField] = kind
		filter[types.RefIDField] = map[string]interface{}{common.BKDBIN: ids}
	}

	var (
		wg       sync.WaitGroup
		firstErr ccErr.CCErrorCoder
		hasRes   bool
	)
	for _, table := range tables {
		wg.Add(1)
		go func(table string, filter []map[string]interface{}) {
			defer func() {
				wg.Done()
			}()

			counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, table, filter)
			if err != nil {
				blog.Errorf("count resource failed, cond: %v, table: %s, err: %v, rid: %s", filter, table, err, kit.Rid)
				firstErr = err
				return
			}

			if counts[0] > 0 {
				blog.Errorf("there are resources under the target that cannot be deleted, filter: %v, table: %s, "+
					"rid: %s", filter, table, kit.Rid)
				hasRes = true
				return
			}

		}(table, []map[string]interface{}{filter})
	}
	wg.Wait()

	if firstErr != nil {
		return hasRes, firstErr
	}

	return hasRes, nil
}
