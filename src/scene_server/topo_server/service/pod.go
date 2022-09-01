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
	"configcenter/src/kube/types"
)

// FindPodPath find pod path
func (s *Service) FindPodPath(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	req := types.PodPathReq{}
	if err := ctx.DecodeInto(&req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	bizIDWithName, err := s.getBizIDWithName(ctx.Kit, []int64{bizID})
	if err != nil {
		blog.Errorf("get bizID with name failed, bizID: %s, err: %v, rid: %s", bizID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	bizName := bizIDWithName[bizID]

	cond := mapstr.MapStr{
		common.BKAppIDField: mapstr.MapStr{common.BKDBEQ: bizID},
		common.BKFieldID:    mapstr.MapStr{common.BKDBIN: req.PodIDs},
	}
	fields := []string{common.BKFieldID, types.BKClusterIDFiled, types.BKNamespaceIDField, types.NamespaceField,
		types.RefField}
	query := &metadata.QueryCondition{
		Condition: cond,
		Fields:    fields,
	}
	option := &types.QueryReq{
		Table:     types.BKTableNameBasePod,
		Condition: query,
	}

	result, err := s.Engine.CoreAPI.CoreService().Kube().FindInst(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		blog.Errorf("find node failed, cond: %v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	paths, err := s.buildPodPaths(ctx.Kit, bizName, result.Info)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(types.PodPathData{
		Info: paths,
	})
}

func (s *Service) buildPodPaths(kit *rest.Kit, bizName string, pods []mapstr.MapStr) ([]types.PodPath, error) {
	paths := make([]types.PodPath, 0)
	clusterIDs := make([]int64, 0)
	for _, pod := range pods {
		id, err := pod.Int64(common.BKFieldID)
		if err != nil {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", common.BKFieldID, pod, err,
				kit.Rid)
			return nil, err
		}

		clusterID, err := pod.Int64(types.BKClusterIDFiled)
		if err != nil {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.BKClusterIDFiled, pod,
				err, kit.Rid)
			return nil, err
		}
		clusterIDs = append(clusterIDs, clusterID)

		namespaceID, err := pod.Int64(types.BKNamespaceIDField)
		if err != nil {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.BKNamespaceIDField, pod,
				err, kit.Rid)
			return nil, err
		}

		namespace, err := pod.String(types.NamespaceField)
		if err != nil {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.NamespaceField, pod,
				err, kit.Rid)
			return nil, err
		}

		ref, err := pod.MapStr(types.RefField)
		if err != nil {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.RefField, pod, err,
				kit.Rid)
			return nil, err
		}

		workloadID, err := ref.Int64(common.BKFieldID)
		if err != nil {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.RefIDField, pod, err,
				kit.Rid)
			return nil, err
		}

		workloadName, err := ref.String(common.BKFieldName)
		if err != nil {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.RefNameField, pod, err,
				kit.Rid)
			return nil, err
		}

		workloadKind, err := ref.String(types.KindField)
		if err != nil {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.RefKindField, pod, err,
				kit.Rid)
			return nil, err
		}

		path := types.PodPath{
			BizName:      bizName,
			ClusterID:    clusterID,
			NamespaceID:  namespaceID,
			Namespace:    namespace,
			Kind:         types.WorkloadType(workloadKind),
			WorkloadID:   workloadID,
			WorkloadName: workloadName,
			PodID:        id,
		}
		paths = append(paths, path)
	}

	if len(clusterIDs) != 0 {
		clusterIDWithName, err := s.getClusterIDWithName(kit, clusterIDs)
		if err != nil {
			blog.Errorf("get cluster id with name failed, clusterIDs: %v, err: %v, rid: %s", clusterIDs, err, kit.Rid)
			return nil, err
		}

		for idx, path := range paths {
			paths[idx].ClusterName = clusterIDWithName[path.ClusterID]
		}
	}

	return paths, nil
}

// ListPod list pod
func (s *Service) ListPod(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	req := types.PodQueryReq{}
	if err := ctx.DecodeInto(&req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond, err := req.BuildCond(bizID, ctx.Kit.SupplierAccount)
	if err != nil {
		blog.Errorf("build query pod condition failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if req.Page.EnableCount {
		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			types.BKTableNameBasePod, []map[string]interface{}{cond})
		if err != nil {
			blog.Errorf("count pod failed, cond: %v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
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
		Table:     types.BKTableNameBasePod,
		Condition: query,
	}
	res, err := s.Engine.CoreAPI.CoreService().Kube().FindInst(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		blog.Errorf("find pod failed, cond: %v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(res)
}

// BatchCreatePod batch create pods.
func (s *Service) BatchCreatePod(ctx *rest.Contexts) {
	data := new(types.CreatePodsOption)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := data.Validate(); err != nil {
		blog.Errorf("batch create pods param verification failed, data: %+v, err: %v, rid: %s", data, err, ctx.Kit.Rid)
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
		ids, err = s.Logics.ContainerOperation().BatchCreatePod(ctx.Kit, data, bizID)
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
