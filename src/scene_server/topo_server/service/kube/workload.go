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

	acmeta "configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// CreateWorkload create workload
func (s *service) CreateWorkload(ctx *rest.Contexts) {
	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	if err := kind.Validate(); err != nil {
		blog.Errorf("workload kind is invalid, kind: %v, err: %v, rid: %s", kind, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	req := types.WlCreateOption{Kind: kind}
	if err := ctx.DecodeInto(&req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeWorkload, Action: acmeta.Create},
		BusinessID: req.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	var data *metadata.RspIDs
	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		res, err := s.createWorkload(ctx.Kit, kind, req)
		if err != nil {
			return err
		}
		data = res
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(data)
}

func (s *service) createWorkload(kit *rest.Kit, kind types.WorkloadType, req types.WlCreateOption) (*metadata.RspIDs,
	error) {

	data, err := s.ClientSet.CoreService().Kube().CreateWorkload(kit.Ctx, kit.Header, kind, req.Data)
	if err != nil {
		blog.Errorf("create workload failed, data: %v, err: %v, rid: %s", req, err, kit.Rid)
		return nil, err
	}

	// audit log.
	audit := auditlog.NewKubeAudit(s.ClientSet.CoreService())
	auditParam := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	for idx := range req.Data {
		wlBase := req.Data[idx].GetWorkloadBase()
		wlBase.BizID = req.BizID
		wlBase.ID = data.IDs[idx]
		wlBase.TenantID = kit.TenantID
		req.Data[idx].SetWorkloadBase(wlBase)
	}

	auditLogs, err := audit.GenerateWorkloadAuditLog(auditParam, req.Data, kind)
	if err != nil {
		blog.Errorf("generate audit log failed, ids: %v, err: %v, rid: %s", data.IDs, err, kit.Rid)
		return nil, err
	}

	if err = audit.SaveAuditLog(kit, auditLogs...); err != nil {
		blog.Errorf("save audit log failed, ids: %v, err: %v, rid: %s", data.IDs, err, kit.Rid)
		return nil, err
	}

	return data, nil
}

// UpdateWorkload update workload
func (s *service) UpdateWorkload(ctx *rest.Contexts) {
	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	if err := kind.Validate(); err != nil {
		ctx.RespAutoError(err)
		return
	}

	req := new(types.WlUpdateOption)
	req.Kind = kind
	if err := ctx.DecodeInto(req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeWorkload, Action: acmeta.Update},
		BusinessID: req.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	wlData, err := s.checkWorkloadData(ctx.Kit, req.BizID, req.IDs, kind)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(wlData) == 0 {
		blog.Errorf("no workload founded, bizID: %s, data: %v, rid: %s", req.BizID, req, ctx.Kit.Rid)
		ctx.RespAutoError(errors.New("no workload founded"))
		return
	}

	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.ClientSet.CoreService().Kube().UpdateWorkload(ctx.Kit.Ctx, ctx.Kit.Header, kind,
			&req.WlUpdateByIDsOption)
		if err != nil {
			blog.Errorf("update workload failed, data: %v, err: %v, rid: %s", req, err, ctx.Kit.Rid)
			return err
		}

		audit := auditlog.NewKubeAudit(s.ClientSet.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
		updateFields, goErr := mapstr.Struct2Map(req.Data)
		if goErr != nil {
			blog.Errorf("update fields convert failed, err: %v, rid: %s", goErr, ctx.Kit.Rid)
			return goErr
		}
		auditParam.WithUpdateFields(updateFields)
		auditLogs, err := audit.GenerateWorkloadAuditLog(auditParam, wlData, kind)
		if err != nil {
			blog.Errorf("generate audit log failed, data: %v, err: %v, rid: %s", wlData, err, ctx.Kit.Rid)
			return err
		}
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed, data: %v, err: %v, rid: %s", wlData, err, ctx.Kit.Rid)
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

// checkWlSharedNs checks if workload's ns is a shared ns and if its biz id is not the same with the input biz id
func (s *service) checkWlSharedNs(kit *rest.Kit, workloads []types.WorkloadInterface, bizID int64) error {
	mismatchNsIDs := make([]int64, 0)
	for _, workload := range workloads {
		wl := workload.GetWorkloadBase()
		if wl.BizID != bizID {
			mismatchNsIDs = append(mismatchNsIDs, wl.NamespaceID)
		}
	}

	if len(mismatchNsIDs) > 0 {
		mismatchNsMap := map[int64][]int64{bizID: mismatchNsIDs}
		if err := s.Logics.KubeOperation().CheckPlatBizSharedNs(kit, mismatchNsMap); err != nil {
			return err
		}
	}
	return nil
}

// DeleteWorkload delete workload
func (s *service) DeleteWorkload(ctx *rest.Contexts) {
	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	if err := kind.Validate(); err != nil {
		ctx.RespAutoError(err)
		return
	}

	req := new(types.WlDeleteOption)
	if err := ctx.DecodeInto(req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeWorkload, Action: acmeta.Delete},
		BusinessID: req.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	wlData, err := s.checkWorkloadData(ctx.Kit, req.BizID, req.IDs, kind)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// if all workloads are already deleted, return
	if len(wlData) == 0 {
		ctx.RespEntity(nil)
		return
	}

	hasRes, rawErr := s.hasNextLevelResource(ctx.Kit, string(kind), req.IDs)
	if rawErr != nil {
		ctx.RespAutoError(rawErr)
		return
	}
	if hasRes {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID))
		return
	}

	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.ClientSet.CoreService().Kube().DeleteWorkload(ctx.Kit.Ctx, ctx.Kit.Header, kind,
			&req.WlDeleteByIDsOption)
		if err != nil {
			blog.Errorf("delete workload failed, data: %v, err: %v, rid: %s", req, err, ctx.Kit.Rid)
			return err
		}

		audit := auditlog.NewKubeAudit(s.ClientSet.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditDelete)
		auditLogs, err := audit.GenerateWorkloadAuditLog(auditParam, wlData, kind)
		if err != nil {
			blog.Errorf("generate audit log failed, data: %v, err: %v, rid: %s", wlData, err, ctx.Kit.Rid)
			return err
		}
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed, data: %v, err: %v, rid: %s", wlData, err, ctx.Kit.Rid)
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

func (s *service) checkWorkloadData(kit *rest.Kit, bizID int64, ids []int64, kind types.WorkloadType) (
	[]types.WorkloadInterface, error) {

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: ids}},
	}
	resp, err := s.ClientSet.CoreService().Kube().ListWorkload(kit.Ctx, kit.Header, query, kind)
	if err != nil {
		blog.Errorf("list workload failed, bizID: %s, ids: %+v, err: %v, rid: %s", bizID, ids, err, kit.Rid)
		return nil, err
	}

	if len(resp.Info) == 0 {
		return nil, nil
	}

	if err := s.checkWlSharedNs(kit, resp.Info, bizID); err != nil {
		return nil, err
	}

	return resp.Info, nil
}

// ListWorkload list workload
func (s *service) ListWorkload(ctx *rest.Contexts) {
	kind := types.WorkloadType(ctx.Request.PathParameter(types.KindField))
	table, err := kind.Table()
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types.KindField))
		return
	}

	req := new(types.WlQueryOption)
	if err := ctx.DecodeInto(req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(kind); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeWorkload, Action: acmeta.Find},
		BusinessID: req.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	// compatible for shared cluster scenario
	cond, err := s.Logics.KubeOperation().GenSharedNsListCond(ctx.Kit, types.KubeWorkload, req.BizID, req.Filter)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if req.Page.EnableCount {
		counts, err := s.ClientSet.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header, table,
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

	resp, err := s.ClientSet.CoreService().Kube().ListWorkload(ctx.Kit.Ctx, ctx.Kit.Header, query, kind)
	if err != nil {
		blog.Errorf("list workload failed, bizID: %s, cond: %v, err: %v, rid: %s", req.BizID, query, err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(resp.Info) == 0 {
		ctx.RespEntityWithCount(0, []mapstr.MapStr{})
		return
	}

	ctx.RespEntityWithCount(0, resp.Info)

}
