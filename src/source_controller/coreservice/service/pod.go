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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/kube/types"
	"configcenter/src/storage/dal/table"
	"configcenter/src/storage/driver/mongodb"
)

func (s *coreService) combinePodData(kit *rest.Kit, inputData *types.CreatePodsOption, ids []uint64) ([]types.Pod,
	[]types.Container, map[int64]struct{}, error) {
	pods, containers := make([]types.Pod, 0), make([]types.Container, 0)
	now := time.Now().Unix()
	nodeIDMap := make(map[int64]struct{})

	var i int

	for _, info := range inputData.Data {
		for _, pod := range info.Pods {
			podTmp, nodeID, err := s.combinationPodsInfo(kit, pod, info.BizID, now, int64(ids[i]))
			if err != nil {
				return nil, nil, nil, err
			}
			// need to be compatible with scenarios where
			// there is no container in the pod, so move it next time.
			i++
			pods = append(pods, podTmp)
			if nodeID != 0 {
				nodeIDMap[nodeID] = struct{}{}
			}

			// skip if there is no container information in the pod
			if len(pod.Containers) == 0 {
				continue
			}
			// generate pod ids field
			cIDs, err := mongodb.Client().NextSequences(kit.Ctx, types.BKTableNameBaseContainer,
				len(pod.Containers))
			if err != nil {
				blog.Errorf("create container failed, generate ids failed, err: %+v, rid: %s", err, kit.Rid)
				return nil, nil, nil, err
			}

			for id, container := range pod.Containers {
				// due to the need to be compatible with the scenario where there is no container in the pod,
				// the left and right bits of the array "ids" need to be obtained to obtain the podID that really
				// needs redundancy.
				data := s.combinationContainerInfo(kit, int64(cIDs[id]), int64(ids[i-1]), now, container)
				containers = append(containers, data)
			}
		}
	}
	return pods, containers, nodeIDMap, nil
}

// BatchCreatePod batch create pods
func (s *coreService) BatchCreatePod(ctx *rest.Contexts) {

	inputData := new(types.CreatePodsOption)
	if err := ctx.DecodeInto(inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}

	var podsLen int
	for _, info := range inputData.Data {
		podsLen += len(info.Pods)
	}
	if podsLen == 0 {
		ctx.RespAutoError(errors.New("no pods need created"))
		return
	}
	// generate pod ids field
	ids, err := mongodb.Client().NextSequences(ctx.Kit.Ctx, types.BKTableNameBasePod, podsLen)
	if err != nil {
		blog.Errorf("create pods failed, generate ids failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	pods, containers, nodeIDMap, err := s.combinePodData(ctx.Kit, inputData, ids)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := mongodb.Client().Table(types.BKTableNameBasePod).Insert(ctx.Kit.Ctx, pods); err != nil {
		blog.Errorf("create pod failed, db insert failed, pods: %+v, err: %+v, rid: %s", pods, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(containers) == 0 {
		ctx.RespEntity(pods)
		return
	}
	err = mongodb.Client().Table(types.BKTableNameBaseContainer).Insert(ctx.Kit.Ctx, containers)
	if err != nil {
		blog.Errorf("create container failed, db insert failed, containers: %+v, err: %+v, rid: %s",
			containers, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if err := s.updateNodeField(ctx.Kit, nodeIDMap); err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(pods)
}

func validateNsCluster(kit *rest.Kit, workload types.WorkloadBase, bizID int64) ccErr.CCErrorCoder {
	if workload.BizID != bizID && workload.ClusterType != types.ClusterShareTypeField {
		blog.Errorf("bizID(%d) in the request is inconsistent with the bizID(%d) in the cluster, "+
			"and the cluster type must be a shared cluster, type is %s, rid: %s", bizID,
			workload.BizID, workload.ClusterType, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed,
			errors.New("cluster must be share type"))
	}
	// in the scenario where the bizID is inconsistent, the ns relationship table needs to be verified.
	if workload.BizID != bizID && workload.ClusterType == types.ClusterShareTypeField {
		countFilter := map[string]interface{}{
			types.BKNamespaceIDField: workload.NamespaceID,
			types.ClusterUIDField:    workload.ClusterUID,
			common.BKAppIDField:      bizID,
			types.BKAsstBizIDField:   workload.BizAsstID,
		}
		count, err := mongodb.Client().Table(types.BKTableNsClusterRelation).Find(countFilter).Count(kit.Ctx)
		if err != nil {
			blog.Errorf("query ns relation failed, filter: %+v, err: %+v, rid: %s", countFilter, err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}

		if count == 0 {
			blog.Errorf("no ns relation founded, filter: %+v, rid: %s", count, countFilter, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed,
				errors.New("no ns relation founded"))
		}
		if count > 1 {
			blog.Errorf("query ns relation num(%d) error, filter: %+v, rid: %s", count, countFilter, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed,
				errors.New("query to multiple relation data"))
		}
	}
	return nil
}

// getSysSpecInfoByCond get the spec redundancy information required by the pod.
func getSysSpecInfoByCond(kit *rest.Kit, spec types.SpecSimpleInfo, bizID int64,
	hostID int64) (*types.SysSpec, bool, ccErr.CCErrorCoder) {
	// 通过workload kind 获取表名
	tableName, err := spec.Ref.Kind.Table()
	if err != nil {
		blog.Errorf("get collection failed, kind: %s, err: %v, rid: %s", spec.Ref.Kind, err, kit.Rid)
		return nil, false, kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	filter := map[string]interface{}{
		common.BKAppIDField:      bizID,
		types.BKClusterIDField:   spec.ClusterID,
		types.BKNamespaceIDField: spec.NamespaceID,
		types.BKIDField:          spec.Ref.ID,
	}
	util.SetModOwner(filter, kit.SupplierAccount)

	fields := []string{types.ClusterUIDField, types.NamespaceField, types.KubeNameField, types.TypeField,
		types.BKAsstBizIDField, common.BKAppIDField}

	workload := make([]types.WorkloadBase, 0)

	err = mongodb.Client().Table(tableName).Find(filter).Fields(fields...).All(kit.Ctx, &workload)
	if err != nil {
		blog.Errorf("query host module config failed, err: %s, rid:%s", err, kit.Rid)
		return nil, false, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(workload) != 1 {
		blog.Errorf("workload gets the wrong amount, filter: %+v, num: %d, rid: %s", filter, len(workload), kit.Rid)
		return nil, false, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed,
			errors.New("workload gets the wrong num"))
	}

	if err := validateNsCluster(kit, workload[0], bizID); err != nil {
		return nil, false, err
	}
	nodeName, hasPod, err := getNodeInfo(kit, spec, bizID, workload[0].BizID)
	if err != nil {
		return nil, false, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return &types.SysSpec{
		WorkloadSpec: types.WorkloadSpec{
			NamespaceSpec: types.NamespaceSpec{
				ClusterSpec: types.ClusterSpec{
					BizID:      bizID,
					ClusterUID: workload[0].ClusterUID,
					ClusterID:  spec.ClusterID,
					BizAsstID:  workload[0].BizAsstID,
				},
				Namespace:   workload[0].Namespace,
				NamespaceID: spec.NamespaceID,
			},
			Ref: types.Reference{Kind: spec.Ref.Kind, Name: workload[0].Name, ID: spec.Ref.ID},
		},
		SupplierAccount: kit.SupplierAccount,
		HostID:          hostID,
		NodeID:          spec.NodeID,
		Node:            nodeName,
	}, hasPod, nil
}

func getNodeInfo(kit *rest.Kit, spec types.SpecSimpleInfo, bizID, clusterBizID int64) (string, bool, error) {

	filter := map[string]interface{}{
		common.BKAppIDField:    bizID,
		common.BKOwnerIDField:  kit.SupplierAccount,
		types.BKClusterIDField: spec.ClusterID,
		types.BKIDField:        spec.NodeID,
	}
	util.SetModOwner(filter, kit.SupplierAccount)

	nodes := make([]types.Node, 0)
	fields := []string{types.KubeNameField, types.HasPodField, types.BKBizIDField, types.ClusterTypeField}
	err := mongodb.Client().Table(types.BKTableNameBaseNode).Find(filter).
		Fields(fields...).All(kit.Ctx, &nodes)
	if err != nil {
		blog.Errorf("query node failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return "", false, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(nodes) != 1 {
		blog.Errorf("node gets the wrong amount, filter: %+v, num: %d, rid: %s", filter, len(nodes), kit.Rid)
		return "", false, kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}

	if nodes[0].Name == nil || nodes[0].HasPod == nil {
		blog.Errorf("query node failed, name or has pod is nil, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return "", false, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if nodes[0].BizID != bizID && nodes[0].ClusterType != types.ClusterShareTypeField {
		countFilter := map[string]interface{}{
			types.BKNodeIDField:    nodes[0].ID,
			types.BKClusterIDField: nodes[0].ClusterID,
		}
		count, err := mongodb.Client().Table(types.BKTableNodeClusterRelation).Find(countFilter).Count(kit.Ctx)
		if err != nil {
			blog.Errorf("query ns relation failed, filter: %+v, err: %+v, rid: %s", countFilter, err, kit.Rid)
			return "", false, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}

		if count == 0 {
			blog.Errorf("no node relation founded, filter: %+v, rid: %s", countFilter, kit.Rid)
			return "", false, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed,
				errors.New("no ns relation founded"))
		}
		if count > 1 {
			blog.Errorf("query node relation num(%d) error, filter: %+v, rid: %s", count, countFilter, kit.Rid)
			return "", false, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed,
				errors.New("query to multiple relation data"))
		}
	}
	return *nodes[0].Name, *nodes[0].HasPod, nil
}

func (s *coreService) combinationPodsInfo(kit *rest.Kit, pod types.PodsInfo, bizID int64, now, id int64) (
	types.Pod, int64, error) {

	sysSpec, hasPod, err := getSysSpecInfoByCond(kit, pod.Spec, bizID, pod.HostID)

	if err != nil {
		return types.Pod{}, 0, err
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

	// this scenario shows that the hasPod flag has been set to true and does not need to be reset
	var nodeID int64
	if !hasPod {
		nodeID = sysSpec.NodeID
	}
	return podInfo, nodeID, nil
}

func (s *coreService) combinationContainerInfo(kit *rest.Kit, containerID, podID, now int64,
	info types.Container) types.Container {

	container := types.Container{
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

	return container
}

// ListPod list pod
func (s *coreService) ListPod(ctx *rest.Contexts) {
	input := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}
	util.SetQueryOwner(input.Condition, ctx.Kit.SupplierAccount)
	pods := make([]types.Pod, 0)
	err := mongodb.Client().Table(types.BKTableNameBasePod).Find(input.Condition).Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).
		Sort(input.Page.Sort).
		Fields(input.Fields...).All(ctx.Kit.Ctx, &pods)
	if err != nil {
		blog.Errorf("search pod failed, cond: %v, err: %v, rid: %s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := &types.PodDataResp{Info: pods}
	ctx.RespEntity(result)
}

func (s *coreService) DeletePods(ctx *rest.Contexts) {
	opt := new(types.DeletePodsByIDsOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// delete the containers in the pods
	delContainerCond := mapstr.MapStr{
		types.BKPodIDField: mapstr.MapStr{common.BKDBIN: opt.PodIDs},
	}

	err := mongodb.Client().Table(types.BKTableNameBaseContainer).Delete(ctx.Kit.Ctx, delContainerCond)
	if err != nil {
		blog.Errorf("delete containers failed, cond: %+v, err: %v, rid: %s", delContainerCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// delete the pods
	delPodCond := mapstr.MapStr{
		types.BKIDField: mapstr.MapStr{common.BKDBIN: opt.PodIDs},
	}

	err = mongodb.Client().Table(types.BKTableNameBasePod).Delete(ctx.Kit.Ctx, delPodCond)
	if err != nil {
		blog.Errorf("delete pods failed, cond: %+v, err: %v, rid: %s", delPodCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// ListContainer list container
func (s *coreService) ListContainer(ctx *rest.Contexts) {
	input := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	containers := make([]types.Container, 0)
	util.SetQueryOwner(input.Condition, ctx.Kit.SupplierAccount)
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
