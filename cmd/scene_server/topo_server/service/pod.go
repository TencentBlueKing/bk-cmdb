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
	types2 "configcenter/pkg/kube/types"
	"fmt"
	"strconv"

	"configcenter/pkg/auditlog"
	"configcenter/pkg/blog"
	"configcenter/pkg/common"
	"configcenter/pkg/http/rest"
	"configcenter/pkg/mapstr"
	"configcenter/pkg/metadata"
)

// FindPodPath find pod path
func (s *Service) FindPodPath(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	req := types2.PodPathReq{}
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
	fields := []string{common.BKFieldID, types2.BKClusterIDFiled, types2.BKNamespaceIDField, types2.NamespaceField,
		types2.RefField}
	query := &metadata.QueryCondition{
		Condition: cond,
		Fields:    fields,
	}
	resp, err := s.Engine.CoreAPI.CoreService().Kube().ListPod(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("find pod failed, cond: %v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	paths, err := s.buildPodPaths(ctx.Kit, bizName, resp.Info)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(types2.PodPathData{
		Info: paths,
	})
}

func (s *Service) buildPodPaths(kit *rest.Kit, bizName string, pods []types2.Pod) ([]types2.PodPath, error) {
	paths := make([]types2.PodPath, 0)
	clusterIDs := make([]int64, 0)
	for _, pod := range pods {
		id := pod.ID

		if pod.ClusterID == 0 {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types2.BKClusterIDFiled,
				kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types2.BKClusterIDFiled)
		}
		clusterID := pod.ClusterID
		clusterIDs = append(clusterIDs, clusterID)

		if pod.NameSpaceID == 0 {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types2.BKNamespaceIDField,
				kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types2.BKNamespaceIDField)
		}
		namespaceID := pod.NameSpaceID

		if pod.NameSpace == "" {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types2.NamespaceField, kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types2.NamespaceField)
		}
		namespace := pod.NameSpace

		//
		if pod.Workload.Kind == "" {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types2.RefKindField, kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types2.RefKindField)
		}
		if pod.Workload.Name == "" {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types2.RefNameField, kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types2.RefNameField)
		}
		if pod.Workload.ID == 0 {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types2.RefIDField, kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types2.RefIDField)
		}
		ref := pod.Workload

		path := types2.PodPath{
			BizName:      bizName,
			ClusterID:    clusterID,
			NamespaceID:  namespaceID,
			Namespace:    namespace,
			Kind:         types2.WorkloadType(ref.Kind),
			WorkloadID:   ref.ID,
			WorkloadName: ref.Name,
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

	req := types2.PodQueryReq{}
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
			types2.BKTableNameBasePod, []map[string]interface{}{cond})
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
		Condition: cond,
		Page:      req.Page,
		Fields:    req.Fields,
	}
	resp, err := s.Engine.CoreAPI.CoreService().Kube().ListPod(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("find pod failed, cond: %v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithCount(0, resp.Info)
}

// BatchCreatePod batch create pods.
func (s *Service) BatchCreatePod(ctx *rest.Contexts) {

	data := new(types2.CreatePodsOption)
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
			blog.Errorf("create pod failed, err: %v, rid: %s", err, ctx.Kit.Rid)
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

// DeletePods delete pods and their containers
func (s *Service) DeletePods(ctx *rest.Contexts) {
	opt := new(types2.DeletePodsOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// get to delete pods and containers in them
	orCond := make([]mapstr.MapStr, len(opt.Data))

	for index, delData := range opt.Data {
		orCond[index] = mapstr.MapStr{
			common.BKAppIDField: delData.BizID,
			types2.BKIDField:    mapstr.MapStr{common.BKDBIN: delData.PodIDs},
		}
	}

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKDBOR: orCond},
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}

	podResp, err := s.Engine.CoreAPI.CoreService().Kube().ListPod(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("find pod failed, cond: %v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// if all pods are already deleted, return
	if len(podResp.Info) == 0 {
		ctx.RespEntity(nil)
		return
	}

	// generate audit logs
	audit := auditlog.NewKubeAudit(s.Engine.CoreAPI.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditDelete)
	auditLogs, err := audit.GeneratePodAuditLog(generateAuditParameter, podResp.Info)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// delete pods
	podIDs := make([]int64, len(podResp.Info))
	for index, pod := range podResp.Info {
		podIDs[index] = pod.ID
	}

	delOpt := &types2.DeletePodsByIDsOption{
		PodIDs: podIDs,
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Engine.CoreAPI.CoreService().Kube().DeletePods(ctx.Kit.Ctx, ctx.Kit.Header, delOpt)
		if err != nil {
			blog.Errorf("delete pods failed, opt: %+v, del opt: %+v, err: %v, rid: %s", opt, delOpt, err, ctx.Kit.Rid)
			return err
		}

		// save audit logs
		if err = audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
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
