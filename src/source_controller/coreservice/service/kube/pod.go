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

func (s *service) combinePodData(kit *rest.Kit, inputData *types.CreatePodsOption, podIDs, containerIDs []uint64) (
	[]types.Pod, []types.Container, []int64, error) {

	pods, containers := make([]types.Pod, 0), make([]types.Container, 0)
	now := time.Now().Unix()

	sysSpecArr, nodeIDs, ccErr := s.core.KubeOperation().GetSysSpecInfoByCond(kit, inputData.Data)
	if ccErr != nil {
		return nil, nil, nil, ccErr
	}

	var podIdx, containerIdx int

	for _, info := range inputData.Data {
		for _, pod := range info.Pods {
			id := int64(podIDs[podIdx])
			sysSpec := sysSpecArr[podIdx]
			podIdx++

			podTmp, err := s.combinationPodsInfo(kit, pod, sysSpec, now, id)
			if err != nil {
				return nil, nil, nil, err
			}

			// need to be compatible with scenarios where
			// there is no container in the pod, so move it next time.
			pods = append(pods, podTmp)

			// skip if there is no container information in the pod
			if len(pod.Containers) == 0 {
				continue
			}

			for _, container := range pod.Containers {
				containerID := int64(containerIDs[containerIdx])
				containerIdx++
				// due to the need to be compatible with the scenario where there is no container in the pod,
				// the left and right bits of the array "ids" need to be obtained to obtain the podID that really
				// needs redundancy.
				data, err := s.combinationContainerInfo(kit, containerID, id, now, container)
				if err != nil {
					return nil, nil, nil, err
				}
				containers = append(containers, data)
			}
		}
	}

	return pods, containers, nodeIDs, nil
}

// BatchCreatePod batch create pods
func (s *service) BatchCreatePod(ctx *rest.Contexts) {
	inputData := new(types.CreatePodsOption)
	if err := ctx.DecodeInto(inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}

	var podsLen, containerLen int
	for _, info := range inputData.Data {
		podsLen += len(info.Pods)
		for _, pod := range info.Pods {
			containerLen += len(pod.Containers)
		}
	}

	if podsLen == 0 {
		ctx.RespAutoError(errors.New("no pods need created"))
		return
	}

	// generate pod ids field
	podIDs, err := mongodb.Client().NextSequences(ctx.Kit.Ctx, types.BKTableNameBasePod, podsLen)
	if err != nil {
		blog.Errorf("generate %d pod ids failed, err: %v, rid: %s", podsLen, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	// generate container ids field
	containerIDs, err := mongodb.Client().NextSequences(ctx.Kit.Ctx, types.BKTableNameBaseContainer, containerLen)
	if err != nil {
		blog.Errorf("generate %d container ids failed, err: %v, rid: %s", containerLen, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	pods, containers, nodeIDs, err := s.combinePodData(ctx.Kit, inputData, podIDs, containerIDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err = mongodb.Client().Table(types.BKTableNameBasePod).Insert(ctx.Kit.Ctx, pods); err != nil {
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

	if err = s.updateNodeHasPodField(ctx.Kit, nodeIDs); err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(pods)
}

func (s *service) combinationPodsInfo(kit *rest.Kit, pod types.PodsInfo, sysSpec types.SysSpec, now, id int64) (
	types.Pod, error) {

	podInfo := types.Pod{
		ID:            id,
		SysSpec:       sysSpec,
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

	return podInfo, nil
}

func (s *service) combinationContainerInfo(kit *rest.Kit, containerID, podID, now int64, info types.Container) (
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
func (s *service) ListPod(ctx *rest.Contexts) {
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

func (s *service) DeletePods(ctx *rest.Contexts) {
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
func (s *service) ListContainer(ctx *rest.Contexts) {
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
