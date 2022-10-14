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
	"configcenter/src/storage/dal/table"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/orm"
	"configcenter/src/kube/types"
	"configcenter/src/storage/driver/mongodb"
)

// BatchCreatePod batch create nodes
func (s *coreService) BatchCreatePod(ctx *rest.Contexts) {

	inputData := new(types.CreatePodsOption)
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	// generate pod ids field
	ids, err := mongodb.Client().NextSequences(ctx.Kit.Ctx, types.BKTableNameBasePod, len(inputData.Data))
	if err != nil {
		blog.Errorf("create pods failed, generate ids failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	pods := make([]types.Pod, 0)
	now := time.Now().Unix()
	nodeIDMap := make(map[int64]struct{})
	for _, info := range inputData.Data {
		for idx, pod := range info.Pods {
			podTmp, nodeID, err := s.insertPodTable(ctx.Kit, pod, info.BizID, now, int64(ids[idx]))
			if err != nil {
				ctx.RespAutoError(err)
				return
			}
			pods = append(pods, podTmp)
			if nodeID != 0 {
				nodeIDMap[nodeID] = struct{}{}
			}

			// skip if there is no container information in the pod
			if len(pod.Containers) == 0 {
				continue
			}
			// generate pod ids field
			containerIDs, err := mongodb.Client().NextSequences(ctx.Kit.Ctx, types.BKTableNameBaseContainer,
				len(pod.Containers))
			if err != nil {
				blog.Errorf("create container failed, generate ids failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
				ctx.RespAutoError(err)
				return
			}

			for id, container := range pod.Containers {
				err := s.insertContainerTable(ctx.Kit, int64(containerIDs[id]), int64(ids[idx]), container, now)
				if err != nil {
					ctx.RespAutoError(err)
					return
				}
			}
		}
	}
	nodeIDs := make([]int64, 0)
	for id := range nodeIDMap {
		nodeIDs = append(nodeIDs, id)
	}

	if len(nodeIDs) == 0 {
		ctx.RespEntityWithError(pods, nil)
		return
	}

	if err := s.updateNodeField(ctx.Kit, nodeIDs); err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(pods, nil)
}

func (s *coreService) insertPodTable(kit *rest.Kit, pod types.PodsInfo, bizID int64, now, id int64) (
	types.Pod, int64, error) {

	sysSpec, hasPod, ccErr := s.core.KubeOperation().GetSysSpecInfoByCond(kit, pod.Spec, bizID, pod.HostID)
	if ccErr != nil {
		return types.Pod{}, 0, ccErr
	}

	podInfo := types.Pod{
		ID:            id,
		SysSpec:       *sysSpec,
		Name:          pod.Name,
		Priority:      pod.Priority,
		Labels:        pod.Labels,
		IP:            pod.IP,
		IPs:           pod.IPs,
		Volumes:       pod.Volumes,
		QOSClass:      pod.QOSClass,
		NodeSelectors: pod.NodeSelectors,
		Tolerations:   pod.Tolerations,
		Revision: table.Revision{
			CreateTime: now,
			LastTime:   now,
			Creator:    kit.User,
			Modifier:   kit.User,
		},
	}
	if err := mongodb.Client().Table(types.BKTableNameBasePod).Insert(kit.Ctx, &podInfo); err != nil {
		blog.Errorf("create pod failed, db insert failed, pod: %+v, err: %+v, rid: %s", podInfo, err, kit.Rid)
		return types.Pod{}, 0, err
	}
	// this scenario shows that the hasPod flag has been set to true and does not need to be reset
	var nodeID int64
	if !hasPod {
		nodeID = sysSpec.NodeID
	}
	return podInfo, nodeID, nil
}

func (s *coreService) insertContainerTable(kit *rest.Kit, containerID, podID int64, info types.Container,
	now int64) error {

	container := &types.Container{
		ID:              containerID,
		SupplierAccount: kit.SupplierAccount,
		PodID:           podID,
		Name:            info.Name,
		ContainerID:     info.ContainerID,
		Image:           info.Image,
		Ports:           info.Ports,
		HostPorts:       info.HostPorts,
		Args:            info.Args,
		Started:         info.Started,
		Limits:          info.Limits,
		ReqSysSpecuests: info.ReqSysSpecuests,
		Liveness:        info.Liveness,
		Environment:     info.Environment,
		Mounts:          info.Mounts,
		Revision: table.Revision{
			CreateTime: now,
			LastTime:   now,
			Creator:    kit.User,
			Modifier:   kit.User,
		},
	}
	err := mongodb.Client().Table(types.BKTableNameBaseContainer).Insert(kit.Ctx, container)
	if err != nil {
		blog.Errorf("create container failed, db insert failed, container: %+v, err: %+v, rid: %s",
			container, err, kit.Rid)
		return err
	}
	return nil
}

// updateNodeField here you need to update the has_pod in the node uniformly
func (s *coreService) updateNodeField(kit *rest.Kit, nodeIDs []int64) error {

	filter := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: nodeIDs,
		},
	}

	updateData := map[string]interface{}{
		types.HasPodField: true,
	}
	if err := mongodb.Client().Table(types.BKTableNameBaseNode).Update(kit.Ctx, filter, updateData); err != nil {
		blog.Errorf("update node has_pod field failed, filter: %v, err: %+v, rid: %s", filter, err, kit.Rid)
		return err
	}
	return nil
}

// BatchCreateNode batch create nodes
func (s *coreService) BatchCreateNode(ctx *rest.Contexts) {

	inputData := new(types.CreateNodesOption)
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, rid: %s", ctx.Request.PathParameter("bk_biz_id"),
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	ctx.RespEntityWithError(s.core.KubeOperation().BatchCreateNode(ctx.Kit, bizID, inputData.Nodes))
}

// SearchClusters search clusters
func (s *coreService) SearchClusters(ctx *rest.Contexts) {

	input := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	clusters := make([]types.Cluster, 0)
	if input.Condition == nil {
		input.Condition = mapstr.New()
	}
	err := mongodb.Client().Table(types.BKTableNameBaseCluster).Find(input.Condition).Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(ctx.Kit.Ctx, &clusters)
	if err != nil {
		blog.Errorf("search cluster failed, cond: %+v, err: %v, rid: %s", input.Condition, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	result := &types.ResponseCluster{Data: clusters}

	ctx.RespEntityWithError(result, nil)

}

// SearchNodes search nodes
func (s *coreService) SearchNodes(ctx *rest.Contexts) {
	input := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if input.Condition == nil {
		input.Condition = mapstr.New()
	}
	nodes := make([]types.Node, 0)
	err := mongodb.Client().Table(types.BKTableNameBaseNode).Find(input.Condition).
		Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(ctx.Kit.Ctx, &nodes)
	if err != nil {
		blog.Errorf("search nodes failed, input %+v, err: %v, rid: %s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := &types.SearchNodeRsp{Data: nodes}
	ctx.RespEntityWithError(result, nil)
}

// BatchUpdateNode batch update node.
func (s *coreService) BatchUpdateNode(ctx *rest.Contexts) {

	input := new(types.UpdateNodeOption)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, rid: %s", ctx.Request.PathParameter("bk_biz_id"),
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	supplierAccount := ctx.Request.PathParameter("supplierAccount")
	if supplierAccount == "" {
		blog.Error("url parameter supplierAccount is not set, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, common.BKOwnerIDField))
		return
	}

	filter := map[string]interface{}{
		types.BKBizIDField:    bizID,
		common.BKOwnerIDField: supplierAccount,
	}
	for _, node := range input.Nodes {

		filter[types.BKIDField] = map[string]interface{}{
			common.BKDBIN: node.NodeIDs,
		}

		opts := orm.NewFieldOptions().AddIgnoredFields(types.IgnoredUpdateNodeFields...)
		updateData, err := orm.GetUpdateFieldsWithOption(node.Data, opts)
		if err != nil {
			blog.Errorf("get update data failed, data: %v, err: %v, rid: %s", node.Data, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		err = mongodb.Client().Table(types.BKTableNameBaseNode).Update(ctx.Kit.Ctx, filter, updateData)
		if err != nil {
			blog.Errorf("update node failed, filter: %v, updateData: %v, err: %v, rid: %s", filter, updateData,
				err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
	}

	ctx.RespEntityWithError(&metadata.UpdatedCount{Count: uint64(len(input.Nodes))}, nil)
}

// BatchUpdateCluster update cluster.
func (s *coreService) BatchUpdateCluster(ctx *rest.Contexts) {

	input := new(types.UpdateClusterOption)

	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, err: %v, rid: %s",
			ctx.Request.PathParameter("bk_biz_id"), err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	supplierAccount := ctx.Request.PathParameter("supplierAccount")
	if supplierAccount == "" {
		blog.Error("url parameter supplierAccount is not set, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, common.BKOwnerIDField))
		return
	}

	for _, one := range input.Clusters {
		filter := map[string]interface{}{
			types.BKBizIDField:    bizID,
			common.BKOwnerIDField: supplierAccount,
		}

		if one.ID != 0 {
			filter[types.BKIDField] = one.ID
		}

		opts := orm.NewFieldOptions().AddIgnoredFields(types.IgnoredUpdateClusterFields...)
		updateData, err := orm.GetUpdateFieldsWithOption(one.Data, opts)
		if err != nil {
			blog.Errorf("get update data failed, data: %v, err: %v, rid: %s", one, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommDBUpdateFailed))
			return
		}

		err = mongodb.Client().Table(types.BKTableNameBaseCluster).Update(ctx.Kit.Ctx, filter, updateData)
		if err != nil {
			blog.Errorf("update cluster failed, filter: %v, updateData: %v, err: %v, rid: %s", filter, updateData,
				err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommDBUpdateFailed))
			return
		}
	}

	ctx.RespEntityWithError(&metadata.UpdatedCount{Count: uint64(len(input.Clusters))}, nil)
}

// CreateCluster create cluster instance.
func (s *coreService) CreateCluster(ctx *rest.Contexts) {

	inputData := new(types.Cluster)

	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, rid: %s", ctx.Request.PathParameter("bk_biz_id"),
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	ctx.RespEntityWithError(s.core.KubeOperation().CreateCluster(ctx.Kit, bizID, inputData))
}

// BatchDeleteCluster delete clusters.
func (s *coreService) BatchDeleteCluster(ctx *rest.Contexts) {

	option := new(types.DeleteClusterOption)
	if err := ctx.DecodeInto(option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, rid: %s", ctx.Request.PathParameter("bk_biz_id"),
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	filter := make(map[string]interface{}, 0)
	num := 0
	if len(option.IDs) > 0 {
		num = len(option.IDs)
		filter = map[string]interface{}{
			common.BKAppIDField:   bizID,
			common.BKOwnerIDField: ctx.Kit.SupplierAccount,
			types.BKIDField: map[string]interface{}{
				common.BKDBIN: option.IDs,
			},
		}
	}
	if err := mongodb.Client().Table(types.BKTableNameBaseCluster).Delete(ctx.Kit.Ctx, filter); err != nil {
		blog.Errorf("delete cluster failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(uint64(num), nil)
}

// BatchDeleteNode delete clusters.
func (s *coreService) BatchDeleteNode(ctx *rest.Contexts) {
	option := new(types.BatchDeleteNodeOption)
	if err := ctx.DecodeInto(option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err != nil {
		ctx.RespAutoError(err)
		return
	}
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if err != nil {
		blog.Error("url parameter bk_biz_id not integer, bizID: %s, rid: %s", ctx.Request.PathParameter("bk_biz_id"),
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	filter := map[string]interface{}{
		common.BKAppIDField:   bizID,
		common.BKOwnerIDField: ctx.Kit.SupplierAccount,
		types.BKIDField: map[string]interface{}{
			common.BKDBIN: option.IDs,
		},
	}
	if err := mongodb.Client().Table(types.BKTableNameBaseNode).Delete(ctx.Kit.Ctx, filter); err != nil {
		blog.Errorf("delete cluster failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	num := len(option.IDs)

	ctx.RespEntityWithError(&metadata.DeletedCount{Count: uint64(num)}, nil)

}

// ListContainer list container
func (s *coreService) ListContainer(ctx *rest.Contexts) {
	input := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	containers := make([]types.Container, 0)
	err := mongodb.Client().Table(types.BKTableNameBaseContainer).Find(input.Condition).Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(ctx.Kit.Ctx, &containers)
	if err != nil {
		blog.Errorf("search container failed, cond: %v, err: %v, rid: %s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := &types.ContainerDataResp{Info: containers}
	ctx.RespEntity(result)
}
