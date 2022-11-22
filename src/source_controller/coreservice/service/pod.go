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
				data, err := s.combinationContainerInfo(kit, int64(cIDs[id]), int64(ids[i-1]), now, container)
				if err != nil {
					return nil, nil, nil, err
				}
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

func (s *coreService) combinationPodsInfo(kit *rest.Kit, pod types.PodsInfo, bizID int64, now, id int64) (
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

	// this scenario shows that the hasPod flag has been set to true and does not need to be reset
	var nodeID int64
	if !hasPod {
		nodeID = sysSpec.NodeID
	}
	return podInfo, nodeID, nil
}

func (s *coreService) combinationContainerInfo(kit *rest.Kit, containerID, podID, now int64, info types.Container) (
	types.Container, error) {

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

	return container, nil
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
