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

// ListContainer list container
func (s *Service) ListContainer(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	req := types.ContainerQueryReq{}
	if err := ctx.DecodeInto(&req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond := mapstr.MapStr{
		common.BKAppIDField: bizID,
		common.BKFieldID:    req.PodID,
	}
	counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
		types.BKTableNameBasePod, []map[string]interface{}{cond})
	if err != nil {
		blog.Errorf("get pod failed, cond: %v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if counts[0] != 1 {
		blog.Errorf("get pod failed, count: %d, cond: %v, err: %v, rid: %s", counts[0], cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.BKPodIDField))
		return
	}

	cond, err = req.BuildCond()
	if err != nil {
		blog.Errorf("build query container condition failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if req.Page.EnableCount {
		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			types.BKTableNameBaseContainer, []map[string]interface{}{cond})
		if err != nil {
			blog.Errorf("count container failed, cond: %v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntityWithCount(counts[0], make([]mapstr.MapStr, 0))
		return
	}

	if req.Page.Sort == "" {
		req.Page.Sort = common.BKFieldID
	}

	query := &metadata.QueryCondition{
		Condition:      cond,
		Page:           req.Page,
		Fields:         req.Fields,
		DisableCounter: true,
	}

	option := &types.QueryReq{
		Table:     types.BKTableNameBaseContainer,
		Condition: query,
	}
	res, err := s.Engine.CoreAPI.CoreService().Kube().FindInst(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		blog.Errorf("find container failed, cond: %v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(res)
}

// BatchDeleteNode delete nodes.
func (s *Service) BatchDeleteNode(ctx *rest.Contexts) {
	option := new(types.BatchDeleteNodeOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	deleteOption := reorganizeDeleteOption(option)

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		err = s.Logics.ContainerOperation().BatchDeleteNode(ctx.Kit, bizID, deleteOption, ctx.Kit.SupplierAccount)
		if err != nil {
			blog.Errorf("delete cluster failed, biz: %d, option: %+v, err: %v, rid: %s", bizID, option, err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
}

// DeleteCluster delete cluster.
func (s *Service) DeleteCluster(ctx *rest.Contexts) {
	option := new(types.DeleteClusterOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err != nil {
		ctx.RespAutoError(err)
		return
	}
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	//1、重新整理一下请求

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		err = s.Logics.ContainerOperation().DeleteCluster(ctx.Kit, bizID, option, ctx.Kit.SupplierAccount)
		if err != nil {
			blog.Errorf("delete cluster failed, biz: %d, option: %+v, err: %v, rid: %s", bizID, option, err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
}

func reorganizeDeleteOption(option *types.BatchDeleteNodeOption) *types.ArrangeDeleteNodeOption {
	deleteOption := make(map[interface{}][]interface{})
	flag := false
	if len(option.NodeIDs) > 0 {
		nameMap := make(map[string]struct{})
		for _, nodeID := range option.NodeIDs {
			if _, ok := nameMap[nodeID.Name]; !ok {
				deleteOption[nodeID.ClusterUID] = append(deleteOption[nodeID.ClusterUID], nodeID.Name)
			}
		}
		flag = true
	}

	if len(option.NodeCmdbIDs) > 0 {
		idMap := make(map[int64]struct{})
		for _, nodeID := range option.NodeCmdbIDs {
			if _, ok := idMap[nodeID.ID]; !ok {
				deleteOption[nodeID.ClusterID] = append(deleteOption[nodeID.ClusterID], nodeID.ID)
			}
		}
	}

	return &types.ArrangeDeleteNodeOption{
		Option: deleteOption,
		Flag:   flag,
	}
}

// BatchCreatePod batch create pods.
func (s *Service) BatchCreatePod(ctx *rest.Contexts) {
	data := new(types.CreatePodsReq)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := data.Validate(); err != nil {
		blog.Errorf("batch create nodes param verification failed, data: %+v, err: %v, rid: %s", data, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	var ids []int64

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		ids, err = s.Logics.ContainerOperation().BatchCreatePod(ctx.Kit, data, bizID, ctx.Kit.SupplierAccount)
		if err != nil {
			blog.Errorf("create business cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(ids)

}

// BatchCreateNode batch create nodes.
func (s *Service) BatchCreateNode(ctx *rest.Contexts) {
	data := new(types.CreateNodesReq)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := data.ValidateCreate(); err != nil {
		blog.Errorf("batch create nodes param verification failed, data: %+v, err: %v, rid: %s", data, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	var ids []int64

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		ids, err = s.Logics.ContainerOperation().BatchCreateNode(ctx.Kit, data, bizID, ctx.Kit.SupplierAccount)
		if err != nil {
			blog.Errorf("create business cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(ids)

}

// convertKubeCondition 根据不同的资源生成不同的查询条件
func (s *Service) findKubeTopoPathIfo(kit *rest.Kit, option *types.KubeTopoPathReq, filter mapstr.MapStr,
	tableNames []string) (*types.KubeTopoPathRsp, error) {

	result := &types.KubeTopoPathRsp{
		Info: make([]types.KubeObjectInfo, 0),
	}

	query := &metadata.QueryCondition{
		Condition:      filter,
		Page:           option.Page,
		Fields:         []string{types.BKIDField, types.KubeNameField},
		DisableCounter: true,
	}

	// 按照拓扑的展示，将folder放到最前面
	if tableNames[0] == types.BKTableNameBaseNamespace {
		result.Info = append(result.Info, types.KubeObjectInfo{
			ID:   types.KubeFolderID,
			Name: types.KubeFolderName,
			Kind: types.KubeFolder,
		})
	}

	for _, tableName := range tableNames {
		// 根据转化的对象找到对应的表，然后根据id 应该是根据不同的
		switch tableName {
		case types.BKTableNameBaseCluster:
			clusters, err := s.Logics.ContainerOperation().SearchCluster(kit, query)
			if err != nil {
				blog.Errorf("search cluster failed, err: %v, rid: %s", err, kit.Rid)
				return result, err
			}
			for _, cluster := range clusters.Data {
				result.Info = append(result.Info, types.KubeObjectInfo{
					ID:   *cluster.ID,
					Name: *cluster.Name,
					Kind: types.KubeCluster,
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
					ID:   id,
					Name: util.GetStrByInterface(namespace[types.KubeNameField]),
					Kind: types.KubeNamespace,
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
					ID:   id,
					Name: util.GetStrByInterface(workload[types.KubeNameField]),
					Kind: kind[tableName],
				})
			}
		}
	}

	return result, nil
}

// 这里整理出来条件如果是folder类型的话，实际上pod数应该是0
// folder的主机数量实际上是需要看clusterID的对应的node表中的hasPod为false的hostID进行去重之后的本业务主机数量。
func combinationConditions(infos []types.KubeResourceInfo, bizID int64,
	supplierAccount string) []map[string]interface{} {

	filters := make([]map[string]interface{}, 0)
	blog.Errorf("0000000000 infos: %+v", infos)
	// 判断拓扑的资源的类别，分为cluster，namespace、和workload三大类。
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

// CountKubeTopoHostsOrPods 计算节点的数量
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
		// 这里需要限制一下，分10批次，每次10个进行查询，之后进行组合
		// filters 这里包含了folder的条件，在查询pod的时候需要将folder的条件过滤出来 并且赋予0值剩余的再进行查询pod数量。
		podFilters := make([]map[string]interface{}, 0)

		resIDMap := make(map[int]struct{})
		for id, filter := range filters {
			// 如果filter中含有 "has_pod"字段 那么说明folder节点，
			if _, ok := filter[types.HasPodField]; ok {
				resIDMap[id] = struct{}{}
				continue
			}
			// 重新整理需要查询pod的条件
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
	filters []map[string]interface{}, bizID int64) (
	[]types.KubeTopoCountRsp, error) {

	// 如果是要获取host的话，1、这块是需要返回所有的hostID。2、对这些hostID 进行去重。3、将这些hostID 和业务ID 组合起来查一下
	// modulehostconfig 表，最终得到的数量才是真正的主机数
	result := make([]types.KubeTopoCountRsp, 0)

	for id, filter := range filters {

		// 得判断一下是否是folder 如果是folder那么需要查的是node表
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

		// 此处计算的是 除了folder之外的pod表中的主机数量
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

		if clusterID, ok := filter[types.BKClusterIDFiled]; ok {
			// 计算集群在folder下的主机数量
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

	// 获取下一级资源对象
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
			// 对于cluster下一级拓扑，除了namespace 还需要增加一个folder
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

// SearchNodes 根据用户指定的条件查询 nodes
func (s *Service) SearchNodes(ctx *rest.Contexts) {

	searchCond := new(types.QueryNodeReq)
	if err := ctx.DecodeInto(searchCond); err != nil {
		blog.Errorf("failed to parse the params, error: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	if cErr := searchCond.Validate(); cErr.ErrCode != 0 {
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

	filter := mapstr.New()
	if searchCond.Filter != nil {
		cond, errKey, rawErr := searchCond.Filter.ToMgo()
		if rawErr != nil {
			blog.Errorf("parse biz failed, filter: %+v, err: %v, rid: %s", searchCond.Filter, rawErr, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errKey))
			return
		}
		filter = cond
	}

	// 无论条件中是否有bk_biz_id、supplier_account,这里统一替换成url中的bk_biz_id 和kit中的supplier_account
	filter[types.BKBizIDField] = bizID
	filter[types.BKSupplierAccountField] = ctx.Kit.SupplierAccount
	if searchCond.HostID != 0 {
		filter[common.BKHostIDField] = searchCond.HostID
	}
	if searchCond.ClusterUID != 0 {
		filter[types.ClusterUIDField] = searchCond.ClusterUID
	}
	if searchCond.ClusterID != 0 {
		filter[types.BKClusterIDFiled] = searchCond.ClusterID
	}

	// count biz in cluster enable count is set
	if searchCond.Page.EnableCount {
		filter := []map[string]interface{}{filter}
		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			types.BKTableNameBaseNode, filter)
		if err != nil {
			blog.Errorf("count node failed, cond: %+v, err: %v, rid: %s", filter, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntityWithCount(counts[0], make([]mapstr.MapStr, 0))
		return
	}

	query := &metadata.QueryCondition{
		Condition: filter,
		Page:      searchCond.Page,
		Fields:    searchCond.Fields,
	}
	result, err := s.Logics.ContainerOperation().SearchNode(ctx.Kit, query)
	if err != nil {
		blog.Errorf("search node failed, filter: %+v, err: %v, rid: %s", filter, err, ctx.Kit.Rid)
		return
	}
	ctx.RespEntityWithCount(0, result.Data)
}

// SearchClusters 根据用户指定的条件查询cluster
func (s *Service) SearchClusters(ctx *rest.Contexts) {

	searchCond := new(types.QueryClusterReq)
	if err := ctx.DecodeInto(searchCond); err != nil {
		blog.Errorf("failed to parse the params, error: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	if cErr := searchCond.Validate(); cErr.ErrCode != 0 {
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

	filter := mapstr.New()
	if searchCond.Filter != nil {
		cond, errKey, rawErr := searchCond.Filter.ToMgo()
		if rawErr != nil {
			blog.Errorf("parse biz filter(%#v) failed, err: %v, rid: %s", searchCond.Filter, rawErr, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errKey))
			return
		}
		filter = cond
	}
	// 无论条件中是否有bk_biz_id、supplier_account,这里统一替换成url中的bk_biz_id 和kit中的supplier_account
	filter[types.BKBizIDField] = bizID
	filter[types.BKSupplierAccountField] = ctx.Kit.SupplierAccount

	// count biz in cluster enable count is set
	if searchCond.Page.EnableCount {
		filter := []map[string]interface{}{filter}

		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			types.BKTableNameBaseCluster, filter)
		if err != nil {
			blog.Errorf("count biz failed, err: %v, cond: %#v, rid: %s", err, filter, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntityWithCount(counts[0], make([]mapstr.MapStr, 0))
		return
	}

	query := &metadata.QueryCondition{
		Condition: filter,
		Page:      searchCond.Page,
		Fields:    searchCond.Fields,
	}
	result, err := s.Logics.ContainerOperation().SearchCluster(ctx.Kit, query)
	if err != nil {
		blog.Errorf("search cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		return
	}
	ctx.RespEntityWithCount(0, result.Data)
}

// UpdateClusterFields 更新集群字段
func (s *Service) UpdateClusterFields(ctx *rest.Contexts) {
	// 还是用之前的结构体作为参数，json解析后看下不空的字段，根据tag获取是否是可编辑的字段，如果是可编辑的字段那么可以更新，如果是不可编辑字
	// 段那么直接报错
	data := new(types.UpdateClusterOption)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := data.Validate(); err != nil {
		ctx.RespAutoError(err)
		return
	}
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		err = s.Logics.ContainerOperation().UpdateClusterFields(ctx.Kit, data, bizID, ctx.Kit.SupplierAccount)
		if err != nil {
			blog.Errorf("create cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
}

// CreateCluster create a container cluster
func (s *Service) CreateCluster(ctx *rest.Contexts) {
	data := new(types.ClusterBaseFields)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := data.ValidateCreate(); err != nil {
		blog.Errorf("validate create container cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	var id int64
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		id, err = s.Logics.ContainerOperation().CreateCluster(ctx.Kit, data, bizID, ctx.Kit.SupplierAccount)
		if err != nil {
			blog.Errorf("create cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(id)
}

// FindResourceAttrs 获取容器对象的属性信息
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
