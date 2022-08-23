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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

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

// BatchCreatePod batch create pods.
func (s *Service) BatchCreatePod(ctx *rest.Contexts) {

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
			// todo: 补充namespace查询
			return result, nil

		default:
			// todo: 补充workload的查询
			return result, nil
		}
	}

	return result, nil
}

func combinationConditions(option []types.KubeResourceInfo) []map[string]interface{} {
	filter := make([]map[string]interface{}, 0)

	// 判断拓扑的资源的类别，分为cluster，namespace、和workload三大类。
	for _, res := range option {
		switch res.Kind {
		case types.KubeCluster:
			filter = append(filter, map[string]interface{}{
				types.BKClusterIDFiled: res.ID,
			})
		case types.NamespaceField:
			filter = append(filter, map[string]interface{}{
				types.BKNamespaceIDField: res.ID,
			})
		default:
			filter = append(filter, map[string]interface{}{
				types.ReferenceID:   res.ID,
				types.ReferenceKind: res.Kind,
			})
		}
	}
	return filter
}

func (s *Service) CountKubeTopoHostsOrPods(ctx *rest.Contexts) {

	option := new(types.KubeTopoCountReq)
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

	filter := combinationConditions(option.ResourceInfos)

	// 这里需要限制一下，分10批次，每次10个进行查询，之后进行组合
	counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
		types.BKTableNameBasePod, filter)
	if err != nil {
		blog.Errorf("count node failed, err: %v, cond: %#v, rid: %s", err, filter, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	result := make([]types.KubeTopoCountRsp, 0)
	for id, count := range counts {
		result = append(result, types.KubeTopoCountRsp{
			Kind:  option.ResourceInfos[id].Kind,
			ID:    option.ResourceInfos[id].ID,
			Count: count,
		})
	}
	ctx.RespEntity(result)

}

// SearchKubeTopoPath 查询容器拓扑路径
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

	filter := mapstr.MapStr{
		types.BKSupplierAccountField: ctx.Kit.SupplierAccount,
		types.BKBizIDField:           bizID,
	}

	// 获取下一级资源对象
	subObject := types.GetKubeSubTopoObject(option.ReferenceObjID)

	tableNames, err := types.GetCollectionWithObject(subObject)
	if err != nil {
		blog.Errorf("failed get , err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(tableNames) > 1 || tableNames[0] != types.BKTableNameBaseCluster {
		filter[types.BKIDField] = option.ReferenceID
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
			blog.Errorf("parse biz filter(%#v) failed, err: %v, rid: %s", searchCond.Filter, rawErr, ctx.Kit.Rid)
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
			blog.Errorf("count node failed, err: %v, cond: %#v, rid: %s", err, filter, ctx.Kit.Rid)
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
		blog.Errorf("search cluster failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		return
	}
	ctx.RespEntity(result.Data)

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
	ctx.RespEntity(result.Data)
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

// FindContainerAttrs 获取容器对象的属性信息
func (s *Service) FindContainerAttrs(ctx *rest.Contexts) {

	object := ctx.Request.PathParameter("object")
	if !types.IsContainerTopoResource(object) {
		blog.Errorf("the parameter is invalid and does not belong to the container object(%s)", object)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "object"))
		return
	}

	result := make([]types.ContainerAttrsRsp, 0)
	switch object {
	case types.KubeCluster:
		for _, descriptor := range types.ClusterSpecFieldsDescriptor {
			result = append(result, types.ContainerAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types.KubeNamespace:
		for _, descriptor := range types.NamespaceSpecFieldsDescriptor {
			result = append(result, types.ContainerAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types.KubeNode:
		for _, descriptor := range types.NodeSpecFieldsDescriptor {
			result = append(result, types.ContainerAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types.KubeWorkload:
		for _, descriptor := range types.WorkLoadSpecFieldsDescriptor {
			result = append(result, types.ContainerAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types.KubePod:
		for _, descriptor := range types.PodSpecFieldsDescriptor {
			result = append(result, types.ContainerAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	case types.KubeContainer:
		for _, descriptor := range types.ContainerSpecFieldsDescriptor {
			result = append(result, types.ContainerAttrsRsp{
				Field:    descriptor.Field,
				Type:     string(descriptor.Type),
				Required: descriptor.IsRequired,
			})
		}
	}
	ctx.RespEntity(result)
}
