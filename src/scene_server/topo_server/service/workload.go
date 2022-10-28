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
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// CreateWorkload create workload
func (s *Service) CreateWorkload(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	if err := kind.Validate(); err != nil {
		blog.Errorf("workload kind is invalid, kind: %v, err: %v, rid: %s", kind, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	req := types.WlCreateReq{Kind: kind}
	if err := ctx.DecodeInto(&req); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	var data *types.WlCreateRespData
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		data, err = s.Engine.CoreAPI.CoreService().Kube().CreateWorkload(ctx.Kit.Ctx, ctx.Kit.Header, bizID, kind, &req)
		if err != nil {
			blog.Errorf("create workload failed, data: %v, err: %v, rid: %s", req, err, ctx.Kit.Rid)
			return err
		}

		// audit log.
		audit := auditlog.NewKubeAudit(s.Engine.CoreAPI.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditCreate)
		for idx := range req.Data {
			wlBase := req.Data[idx].GetWorkloadBase()
			wlBase.BizID = bizID
			wlBase.ID = data.IDs[idx]
			wlBase.SupplierAccount = ctx.Kit.SupplierAccount
			req.Data[idx].SetWorkloadBase(wlBase)
		}
		auditLogs, err := audit.GenerateWorkloadAuditLog(auditParam, req.Data, kind)
		if err != nil {
			blog.Errorf("generate audit log failed, ids: %v, err: %v, rid: %s", data.IDs, err, ctx.Kit.Rid)
			return err
		}
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed, ids: %v, err: %v, rid: %s", data.IDs, err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(data)
}

// UpdateWorkload update workload
func (s *Service) UpdateWorkload(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	if err := kind.Validate(); err != nil {
		ctx.RespAutoError(err)
		return
	}

	req := types.WlUpdateReq{Kind: kind}
	if err := ctx.DecodeInto(&req); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond, err := req.BuildCond(0, false, ctx.Kit.SupplierAccount)
	if err != nil {
		blog.Errorf("build workload condition failed, bizID: %s, data: %v, err: %v, rid: %s", bizID, req, err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	query := &metadata.QueryCondition{
		Condition: cond,
	}
	workloads, err := s.listWorkload(ctx.Kit, query, kind)
	if err != nil {
		blog.Errorf("list workload failed, bizID: %s, data: %v, err: %v, rid: %s", bizID, req, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	for _, workload := range workloads {
		ids := make([]int64, 0)
		if workload.GetWorkloadBase().BizID != bizID {
			ids = append(ids, workload.GetWorkloadBase().ID)
		}

		if len(ids) != 0 {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, ids))
			return
		}
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Engine.CoreAPI.CoreService().Kube().UpdateWorkload(ctx.Kit.Ctx, ctx.Kit.Header, bizID, kind, &req)
		if err != nil {
			blog.Errorf("update workload failed, data: %v, err: %v, rid: %s", req, err, ctx.Kit.Rid)
			return err
		}

		audit := auditlog.NewKubeAudit(s.Engine.CoreAPI.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
		updateFields, goErr := mapstr.Struct2Map(req.Info)
		if goErr != nil {
			blog.Errorf("update fields convert failed, err: %v, rid: %s", goErr, ctx.Kit.Rid)
			return goErr
		}
		auditParam.WithUpdateFields(updateFields)
		auditLogs, err := audit.GenerateWorkloadAuditLog(auditParam, workloads, kind)
		if err != nil {
			blog.Errorf("generate audit log failed, data: %v, err: %v, rid: %s", workloads, err, ctx.Kit.Rid)
			return err
		}
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed, data: %v, err: %v, rid: %s", workloads, err, ctx.Kit.Rid)
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

// DeleteWorkload delete workload
func (s *Service) DeleteWorkload(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	if err := kind.Validate(); err != nil {
		ctx.RespAutoError(err)
		return
	}

	req := types.WlDeleteReq{}
	if err := ctx.DecodeInto(&req); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond, err := req.BuildCond(0, false, ctx.Kit.SupplierAccount)
	if err != nil {
		blog.Errorf("delete workload failed, bizID: %s, data: %v, err: %v, rid: %s", bizID, req, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	query := &metadata.QueryCondition{
		Condition: cond,
	}
	workloads, err := s.listWorkload(ctx.Kit, query, kind)
	if err != nil {
		blog.Errorf("list workload failed, bizID: %s, req: %v, err: %v, rid: %s", bizID, req, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	for _, workload := range workloads {
		ids := make([]int64, 0)
		if workload.GetWorkloadBase().BizID != bizID {
			ids = append(ids, workload.GetWorkloadBase().ID)
		}

		if len(ids) != 0 {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, ids))
			return
		}
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Engine.CoreAPI.CoreService().Kube().DeleteWorkload(ctx.Kit.Ctx, ctx.Kit.Header, bizID, kind, &req)
		if err != nil {
			blog.Errorf("delete workload failed, data: %v, err: %v, rid: %s", req, err, ctx.Kit.Rid)
			return err
		}

		// audit log.
		audit := auditlog.NewKubeAudit(s.Engine.CoreAPI.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditDelete)
		auditLogs, err := audit.GenerateWorkloadAuditLog(auditParam, workloads, kind)
		if err != nil {
			blog.Errorf("generate audit log failed, data: %v, err: %v, rid: %s", workloads, err, ctx.Kit.Rid)
			return err
		}
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed, data: %v, err: %v, rid: %s", workloads, err, ctx.Kit.Rid)
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

func (s *Service) listWorkload(kit *rest.Kit, query *metadata.QueryCondition, kind types.WorkloadType) (
	[]types.WorkloadI, error) {

	resp, err := s.Engine.CoreAPI.CoreService().Kube().ListWorkload(kit.Ctx, kit.Header, query, kind)
	if err != nil {
		blog.Errorf("find workload failed, kind: %s, cond: %v, err: %v, rid: %s", kind, query, err, kit.Rid)
		return nil, err
	}

	return resp.Info, nil
}

// ListWorkload list workload
func (s *Service) ListWorkload(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	table, err := kind.Table()
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KindField))
		return
	}

	req := types.WlQueryReq{}
	if err := ctx.DecodeInto(&req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(kind); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond, err := req.BuildCond(bizID, ctx.Kit.SupplierAccount)
	if err != nil {
		blog.Errorf("build query workload condition failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if req.Page.EnableCount {
		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header, table,
			[]map[string]interface{}{cond})
		if err != nil {
			blog.Errorf("count workload failed, table: %s, cond: %v, err: %v, rid: %s", table, cond, err, ctx.Kit.Rid)
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

	workloads, err := s.listWorkload(ctx.Kit, query, kind)
	if err != nil {
		blog.Errorf("list workload failed, bizID: %s, cond: %v, err: %v, rid: %s", bizID, query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(workloads) == 0 {
		ctx.RespEntityWithCount(0, []mapstr.MapStr{})
		return
	}

	ctx.RespEntityWithCount(0, workloads)

}
