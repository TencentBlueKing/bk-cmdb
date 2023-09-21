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

package kube

import (
	"errors"
	"sync"

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	acmeta "configcenter/src/ac/meta"
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
func (s *service) findKubeTopoPathInfo(kit *rest.Kit, option *types.KubeTopoPathOption, filter mapstr.MapStr,
	tableNames []string, isShared bool) (*types.KubeTopoPathRsp, error) {

	result := &types.KubeTopoPathRsp{Info: make([]types.KubeObjectInfo, 0)}

	query := &metadata.QueryCondition{
		Condition:      filter,
		Page:           option.Page,
		Fields:         []string{types.BKIDField, types.KubeNameField},
		DisableCounter: true}

	// according to the topology display, put the folder to the front
	if tableNames[0] == types.BKTableNameBaseNamespace && !isShared {
		folderName := types.KubeFolderName
		if util.GetLanguage(kit.Header) == string(common.English) {
			folderName = types.KubeFolderNameEn
		}
		result.Info = append(result.Info, types.KubeObjectInfo{
			ID: types.KubeFolderID, Name: folderName, Kind: types.KubeFolder,
		})
	}

	for _, tableName := range tableNames {
		switch tableName {
		case types.BKTableNameBaseCluster:
			clusters, err := s.ClientSet.CoreService().Kube().SearchCluster(kit.Ctx, kit.Header, query)
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

			namespaces, err := s.ClientSet.CoreService().Kube().ListNamespace(kit.Ctx, kit.Header, query)
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
			workloads, cErr := s.ClientSet.CoreService().Kube().ListWorkload(kit.Ctx, kit.Header, query,
				types.WorkloadType(kind[tableName]))
			if cErr != nil {
				blog.Errorf("find %s failed, cond: %v, err: %v, rid: %s", kind[tableName], query, cErr, kit.Rid)
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

func (s *service) combinationConditions(kit *rest.Kit, infos []types.KubeResourceInfo, bizID int64) (
	[]map[string]interface{}, error) {

	filters := make([]map[string]interface{}, 0)
	for _, info := range infos {
		var resFilter *filter.Expression

		switch info.Kind {
		case types.KubeFolder:
			resFilter = &filter.Expression{
				RuleFactory: &filter.CombinedRule{
					Condition: filter.And,
					Rules: []filter.RuleFactory{
						filtertools.GenAtomFilter(types.BKClusterIDFiled, filter.Equal, info.ID),
						filtertools.GenAtomFilter(types.HasPodField, filter.Equal, false),
					},
				},
			}
		case types.KubeCluster:
			resFilter = filtertools.GenAtomFilter(types.BKClusterIDFiled, filter.Equal, info.ID)
		case types.KubeNamespace:
			resFilter = filtertools.GenAtomFilter(types.BKNamespaceIDField, filter.Equal, info.ID)
		default:
			resFilter = &filter.Expression{
				RuleFactory: &filter.CombinedRule{
					Condition: filter.And,
					Rules: []filter.RuleFactory{
						filtertools.GenAtomFilter(types.RefIDField, filter.Equal, info.ID),
						filtertools.GenAtomFilter(types.RefKindField, filter.Equal, info.Kind),
					},
				},
			}
		}

		// compatible for shared cluster scenario
		cond, err := s.Logics.KubeOperation().GenSharedNsListCond(kit, types.KubePod, bizID, resFilter)
		if err != nil {
			return nil, err
		}

		filters = append(filters, cond)
	}
	return filters, nil
}

func (s *service) countKubeHostOrPodsByCond(kit *rest.Kit, option *types.KubeTopoCountOption, kind string) (
	[]types.KubeTopoCountRsp, error) {

	filters, err := s.combinationConditions(kit, option.ResourceInfos, option.BizID)
	if err != nil {
		return nil, err
	}

	switch kind {
	case types.KubeHostKind:
		result, err := s.getTopoHostNumber(kit, option.ResourceInfos, filters, option.BizID)
		if err != nil {
			blog.Errorf("get host number failed, option: %+v, err: %v, rid: %s", option, err)
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

		counts, err := s.ClientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
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
func (s *service) CountKubeTopoHostsOrPods(ctx *rest.Contexts) {
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

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.Business, Action: acmeta.ViewBusinessResource},
		BusinessID: option.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	kind := ctx.Request.PathParameter("type")

	result, err := s.countKubeHostOrPodsByCond(ctx.Kit, option, kind)
	if err != nil {
		blog.Errorf("failed to get(%s) number, option: %+v, err: %v, rid: %s", kind, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *service) getTopoHostNumber(kit *rest.Kit, resourceInfos []types.KubeResourceInfo,
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

			res, err := s.getTopoHostNumByFilter(kit, bizID, id, filter, resourceInfos)
			if err != nil {
				firstErr = err
				return
			}

			lock.Lock()
			result = append(result, res)
			lock.Unlock()
		}(id, filter)
	}

	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}
	return result, nil
}

func (s *service) getTopoHostNumByFilter(kit *rest.Kit, bizID int64, id int, filter map[string]interface{},
	resourceInfos []types.KubeResourceInfo) (types.KubeTopoCountRsp, error) {

	// determine whether this node is a folder If it is a folder, then you need to check the node table.
	if resourceInfos[id].Kind == types.KubeFolder {
		hostMap, err := s.getHostIDsInNodeByCond(kit, bizID, filter)
		if err != nil {
			return types.KubeTopoCountRsp{}, err
		}

		return types.KubeTopoCountRsp{
			Kind:  resourceInfos[id].Kind,
			ID:    resourceInfos[id].ID,
			Count: int64(len(hostMap)),
		}, nil
	}

	// what counts here is the number of hosts in the pod table excluding folders.
	hostMap, err := s.getHostIDsInPodsByCond(kit, bizID, filter)
	if err != nil {
		return types.KubeTopoCountRsp{}, err
	}

	workloadType := types.WorkloadType(util.GetStrByInterface(filter[types.RefKindField]))
	// the scenario dealt with here is the workload type calculation number.
	if err := workloadType.Validate(); err == nil {
		refID, err := util.GetInt64ByInterface(filter[types.RefIDField])
		if err != nil {
			return types.KubeTopoCountRsp{}, err
		}

		return types.KubeTopoCountRsp{
			Kind:  util.GetStrByInterface(filter[types.RefKindField]),
			ID:    refID,
			Count: int64(len(hostMap)),
		}, nil
	}

	resultHostMap, err := s.getClusterNumFromFolder(kit, bizID, filter)
	if err != nil {
		return types.KubeTopoCountRsp{}, err
	}

	for hostID := range hostMap {
		resultHostMap[hostID] = struct{}{}
	}

	return types.KubeTopoCountRsp{
		Kind:  resourceInfos[id].Kind,
		ID:    resourceInfos[id].ID,
		Count: int64(len(resultHostMap)),
	}, nil
}

// getClusterNumFromFolder for the calculation of the number of hosts under the cluster,
// it is necessary to add the number of hosts under the folder node under the cluster.
func (s *service) getClusterNumFromFolder(kit *rest.Kit, bizID int64, filter map[string]interface{}) (
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

func (s *service) getDistinctHostMap(kit *rest.Kit, bizID int64, hostIDMap map[int64]struct{}) (
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

	hostRelations, err := s.ClientSet.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, relationReq)
	if err != nil {
		return nil, err
	}
	for _, info := range hostRelations.Info {
		result[info.HostID] = struct{}{}
	}
	return result, nil
}

func (s *service) getHostIDsInNodeByCond(kit *rest.Kit, bizID int64, cond mapstr.MapStr) (map[int64]struct{},
	error) {

	query := &metadata.QueryCondition{
		Condition:      cond,
		Fields:         []string{common.BKHostIDField},
		DisableCounter: true,
	}
	nodes, cErr := s.ClientSet.CoreService().Kube().SearchNode(kit.Ctx, kit.Header, query)
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

func (s *service) getHostIDsInPodsByCond(kit *rest.Kit, bizID int64, cond mapstr.MapStr) (map[int64]struct{},
	error) {

	query := &metadata.QueryCondition{
		Condition:      cond,
		Fields:         []string{common.BKHostIDField},
		DisableCounter: true,
	}

	pods, cErr := s.ClientSet.CoreService().Kube().ListPod(kit.Ctx, kit.Header, query)
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
func (s *service) SearchKubeTopoPath(ctx *rest.Contexts) {
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

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.Business, Action: acmeta.ViewBusinessResource},
		BusinessID: option.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	// get the next level resource object.
	subObject, cond, isShared, err := s.GetKubeSubTopoObject(ctx.Kit, option.ReferenceObjID, option.ReferenceID,
		option.BizID)
	if err != nil {
		blog.Errorf("failed to get subObject, option: %+v, err: %v, rid: %s", option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
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
			counts, err := s.ClientSet.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header, tableName,
				[]map[string]interface{}{cond})
			if err != nil {
				blog.Errorf("count %s failed, err: %v, cond: %#v, rid: %s", err, tableName, cond, ctx.Kit.Rid)
				ctx.RespAutoError(err)
				return
			}
			// for the next-level topology of the cluster, a folder needs to be added in addition to the namespace.
			if tableName == types.BKTableNameBaseNamespace && !isShared {
				counts[0] += 1
			}
			count += counts[0]
		}

		ctx.RespEntityWithCount(count, make([]mapstr.MapStr, 0))
		return
	}

	result, err := s.findKubeTopoPathInfo(ctx.Kit, option, cond, tableNames, isShared)
	if err != nil {
		blog.Errorf("failed to get topo path, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// GetKubeSubTopoObject get the next-level topology resource object of the specified resource
func (s *service) GetKubeSubTopoObject(kit *rest.Kit, object string, id int64, bizID int64) (string,
	map[string]interface{}, bool, error) {

	switch object {
	case types.KubeBusiness:
		// compatible for shared cluster scenario
		cond, err := s.Logics.KubeOperation().GenSharedClusterListCond(kit, bizID, nil)
		if err != nil {
			return "", nil, false, err
		}
		return types.KubeCluster, cond, false, nil
	case types.KubeCluster:
		// compatible for shared cluster scenario
		clusterCond := &metadata.QueryCondition{
			Fields:    []string{common.BKAppIDField, types.TypeField},
			Page:      metadata.BasePage{Limit: 1},
			Condition: mapstr.MapStr{types.BKIDField: id},
		}
		clusterRes, err := s.ClientSet.CoreService().Kube().SearchCluster(kit.Ctx, kit.Header, clusterCond)
		if err != nil {
			blog.Errorf("get cluster by id %d failed, err: %v, rid: %s", id, err, kit.Rid)
			return "", nil, false, err
		}

		if len(clusterRes.Data) != 1 {
			blog.Errorf("get cluster by id %d, but count is wrong, rid: %s", id, kit.Rid)
			return "", nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, types.BKClusterIDFiled)
		}

		if clusterRes.Data[0].BizID == bizID {
			return types.KubeNamespace, mapstr.MapStr{types.BKClusterIDFiled: id}, false, nil
		}

		if clusterRes.Data[0].Type == nil || *clusterRes.Data[0].Type != types.SharedClusterType {
			blog.Errorf("cluster %d is not shared cluster, rid: %s", id, kit.Rid)
			return "", nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, types.BKClusterIDFiled)
		}
		return types.KubeNamespace, mapstr.MapStr{types.BKClusterIDFiled: id, types.BKBizIDField: bizID}, true, nil
	case types.KubeNamespace:
		return types.KubeWorkload, map[string]interface{}{
			types.BKNamespaceIDField: id,
		}, false, nil
	case types.KubeFolder, types.KubePod:
		return "", nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_reference_obj_id")
	default:
		return types.KubePod, make(map[string]interface{}), false, nil
	}
}

// FindResourceAttrs get the attribute information of the kube object
// for front-end use only
func (s *service) FindResourceAttrs(ctx *rest.Contexts) {

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

func (s *service) hasNextLevelResource(kit *rest.Kit, kind string, ids []int64) (bool, error) {
	var tables []string
	filter := make(map[string]interface{})

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

			counts, err := s.ClientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, table, filter)
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
