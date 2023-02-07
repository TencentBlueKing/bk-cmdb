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
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/auditlog"
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

	req := new(types.PodPathOption)
	if err := ctx.DecodeInto(req); err != nil {
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
	resp, err := s.Engine.CoreAPI.CoreService().Kube().ListPod(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("find pod failed, cond: %v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(resp.Info) == 0 {
		ctx.RespEntity(types.PodPathData{Info: []types.PodPath{}})
		return
	}

	paths, err := s.buildPodPaths(ctx.Kit, bizName, resp.Info)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(types.PodPathData{
		Info: paths,
	})
}

func (s *Service) buildPodPaths(kit *rest.Kit, bizName string, pods []types.Pod) ([]types.PodPath, error) {
	paths := make([]types.PodPath, 0)
	clusterIDs := make([]int64, 0)
	for _, pod := range pods {
		id := pod.ID

		if pod.ClusterID == 0 {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.BKClusterIDFiled,
				kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types.BKClusterIDFiled)
		}
		clusterID := pod.ClusterID
		clusterIDs = append(clusterIDs, clusterID)

		if pod.NamespaceID == 0 {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.BKNamespaceIDField,
				kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types.BKNamespaceIDField)
		}
		namespaceID := pod.NamespaceID

		if pod.Namespace == "" {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.NamespaceField, kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types.NamespaceField)
		}
		namespace := pod.Namespace

		if pod.Ref.Kind == "" {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.RefKindField, kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types.RefKindField)
		}
		if pod.Ref.Name == "" {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.RefNameField, kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types.RefNameField)
		}
		if pod.Ref.ID == 0 {
			blog.Errorf("get pod attribute failed, attr: %s, pod: %v, err: %v, rid: %s", types.RefIDField, kit.Rid)
			return nil, fmt.Errorf("get pod attribute failed, attr: %s", types.RefIDField)
		}
		ref := pod.Ref

		path := types.PodPath{
			BizName:      bizName,
			ClusterID:    clusterID,
			NamespaceID:  namespaceID,
			Namespace:    namespace,
			Kind:         ref.Kind,
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

	req := new(types.PodQueryOption)
	if err := ctx.DecodeInto(req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond, err := req.BuildCond(bizID)
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

	data := new(types.CreatePodsOption)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := data.Validate(); err.ErrCode != 0 {
		blog.Errorf("batch create pods param verification failed, data: %+v, err: %v, rid: %s", data, err, ctx.Kit.Rid)
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	filters := make([]map[string]interface{}, 0)
	for _, info := range data.Data {
		for _, pod := range info.Pods {
			filter := map[string]interface{}{
				types.BKBizIDField:       info.BizID,
				types.BKClusterIDFiled:   pod.ClusterID,
				types.BKNamespaceIDField: pod.NamespaceID,
				types.BKNodeIDField:      pod.NodeID,
				types.KubeNameField:      *pod.Name,
				types.RefKindField:       pod.Ref.Kind,
				types.RefIDField:         pod.Ref.ID,
			}
			filters = append(filters, filter)
		}
	}

	counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
		types.BKTableNameBasePod, filters)
	if err != nil {
		blog.Errorf("count pods failed, filter: %#v, err: %v, rid: %s", filters, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	var podNum int64
	for _, count := range counts {
		podNum += count
	}
	if podNum > 0 {
		blog.Errorf("some pods already exists and the creation fails, filter: %#v, rid: %s", filters, ctx.Kit.Rid)
		ctx.RespAutoError(errors.New("some pod already exists and the creation fails"))
		return
	}

	var ids []int64
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		ids, err = s.Logics.KubeOperation().BatchCreatePod(ctx.Kit, data)
		if err != nil {
			blog.Errorf("create pods failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(metadata.RspIDs{IDs: ids})
}

// DeletePods delete pods and their containers
func (s *Service) DeletePods(ctx *rest.Contexts) {
	opt := new(types.DeletePodsOption)
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
			types.BKIDField:     mapstr.MapStr{common.BKDBIN: delData.PodIDs},
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

	delOpt := &types.DeletePodsByIDsOption{
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

// ListContainer list container
func (s *Service) ListContainer(ctx *rest.Contexts) {
	req := new(types.ContainerQueryOption)
	if err := ctx.DecodeInto(req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond, err := req.BuildCond()
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
		Condition: cond,
		Page:      req.Page,
		Fields:    req.Fields,
	}

	resp, err := s.Engine.CoreAPI.CoreService().Kube().ListContainer(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("find container failed, cond: %v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithCount(0, resp.Info)
}
