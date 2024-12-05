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
	"fmt"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/kube/types"
	"configcenter/src/storage/driver/mongodb"
)

// ListContainerByPod list container by pod condition
func (k *kubeOperation) ListContainerByPod(kit *rest.Kit, input *types.GetContainerByPodOption) (
	*types.GetContainerByPodResp, errors.CCErrorCoder) {

	// do not have pod condition, use container condition to list container
	if len(input.PodCond) == 0 {
		return k.listContainerByContainerCond(kit, input)
	}

	// do not have container condition, use pod condition to list container
	if len(input.ContainerCond) == 0 {
		return k.listContainerByPodCond(kit, input)
	}

	// use both container and pod condition to list container
	return k.listContainerByBothCond(kit, input)
}

// listContainerByContainerCond list container by container condition
func (k *kubeOperation) listContainerByContainerCond(kit *rest.Kit, input *types.GetContainerByPodOption) (
	*types.GetContainerByPodResp, errors.CCErrorCoder) {

	cond := input.ContainerCond

	if input.Page.EnableCount {
		cnt, err := mongodb.Shard(kit.ShardOpts()).Table(types.BKTableNameBaseContainer).Find(cond).Count(kit.Ctx)
		if err != nil {
			blog.Errorf("count container by cond: %+v failed, err: %v, rid: %s", cond, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}

		return &types.GetContainerByPodResp{Count: int64(cnt)}, nil
	}

	containers := make([]mapstr.MapStr, 0)
	err := mongodb.Shard(kit.ShardOpts()).Table(types.BKTableNameBaseContainer).Find(cond).
		Start(uint64(input.Page.Start)).Limit(uint64(input.Page.Limit)).Sort(input.Page.Sort).Fields(input.Fields...).
		All(kit.Ctx, &containers)
	if err != nil {
		blog.Errorf("list container by cond: %+v failed, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return &types.GetContainerByPodResp{Info: containers}, nil
}

// listContainerByPodCond list container by pod condition
func (k *kubeOperation) listContainerByPodCond(kit *rest.Kit, input *types.GetContainerByPodOption) (
	*types.GetContainerByPodResp, errors.CCErrorCoder) {

	fieldPrefix := "container"
	fieldMap := make(map[string]string)
	for _, name := range types.ContainerFields.GetFieldNames() {
		// do not need to specify inner fields
		if strings.Contains(name, ".") {
			continue
		}
		fieldMap[name] = fmt.Sprintf("$%s.%s", fieldPrefix, name)
	}

	filter := []map[string]interface{}{
		{common.BKDBMatch: input.PodCond},
		{
			common.BKDBLookup: map[string]string{
				common.BKDBFrom:         types.BKTableNameBaseContainer,
				common.BKDBLocalField:   common.BKFieldID,
				common.BKDBForeignField: types.BKPodIDField,
				common.BKDBAs:           fieldPrefix,
			},
		},
		{common.BKDBUnwind: fmt.Sprintf("$%s", fieldPrefix)},
		{common.BKDBProject: fieldMap},
	}

	return k.aggregateContainer(kit, types.BKTableNameBasePod, input, filter)
}

// listContainerByBothCond list container by container and pod condition
func (k *kubeOperation) listContainerByBothCond(kit *rest.Kit, input *types.GetContainerByPodOption) (
	*types.GetContainerByPodResp, errors.CCErrorCoder) {

	filter := []map[string]interface{}{
		{common.BKDBMatch: input.ContainerCond},
		{
			common.BKDBLookup: map[string]string{
				common.BKDBFrom:         types.BKTableNameBasePod,
				common.BKDBLocalField:   types.BKPodIDField,
				common.BKDBForeignField: common.BKFieldID,
				common.BKDBAs:           "pod",
			},
		},
		{common.BKDBMatch: input.PodCond},
	}

	return k.aggregateContainer(kit, types.BKTableNameBaseContainer, input, filter)
}

func (k *kubeOperation) aggregateContainer(kit *rest.Kit, table string, input *types.GetContainerByPodOption,
	filter []map[string]interface{}) (*types.GetContainerByPodResp, errors.CCErrorCoder) {

	if input.Page.EnableCount {
		total := "total"
		result := map[string]int64{}
		filter = append(filter, map[string]interface{}{common.BKDBCount: total})
		err := mongodb.Shard(kit.ShardOpts()).Table(table).AggregateOne(kit.Ctx, filter, &result)
		if err != nil && !mongodb.IsNotFoundError(err) {
			blog.Errorf("get container count failed, cond: %+v, err: %v, rid: %s", filter, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}

		return &types.GetContainerByPodResp{Count: result[total]}, nil
	}

	containerFieldMap := make(map[string]int64)
	if len(input.Fields) > 0 {
		for _, field := range input.Fields {
			containerFieldMap[field] = 1
		}
	} else {
		for _, field := range types.PodFields.GetFieldNames() {
			// do not need to specify inner fields
			if strings.Contains(field, ".") {
				continue
			}
			containerFieldMap[field] = 1
		}
	}
	filter = append(filter, map[string]interface{}{common.BKDBProject: containerFieldMap})
	filter = append(filter, []map[string]interface{}{{common.BKDBSort: map[string]int64{input.Page.Sort: 1}},
		{common.BKDBSkip: input.Page.Start}, {common.BKLimit: input.Page.Limit}}...)

	containers := make([]mapstr.MapStr, 0)
	err := mongodb.Shard(kit.ShardOpts()).Table(table).AggregateAll(kit.Ctx, filter, &containers)
	if err != nil {
		blog.Errorf("get all object models failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return &types.GetContainerByPodResp{Info: containers}, nil
}
