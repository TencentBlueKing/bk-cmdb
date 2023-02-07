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

// CreateNamespace create namespace
func (s *Service) CreateNamespace(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	req := new(types.NsCreateOption)
	if err := ctx.DecodeInto(req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	var data *metadata.RspIDs
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		data, err = s.Engine.CoreAPI.CoreService().Kube().CreateNamespace(ctx.Kit.Ctx, ctx.Kit.Header, bizID, req)
		if err != nil {
			blog.Errorf("create namespace failed, data: %v, err: %v, rid: %s", req, err, ctx.Kit.Rid)
			return err
		}

		// audit log.
		audit := auditlog.NewKubeAudit(s.Engine.CoreAPI.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditCreate)
		for idx := range req.Data {
			req.Data[idx].BizID = bizID
			req.Data[idx].ID = data.IDs[idx]
			req.Data[idx].SupplierAccount = ctx.Kit.SupplierAccount
		}
		auditLogs, err := audit.GenerateNamespaceAuditLog(auditParam, req.Data)
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

// UpdateNamespace update namespace
func (s *Service) UpdateNamespace(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	req := new(types.NsUpdateOption)
	if err := ctx.DecodeInto(req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKFieldID: mapstr.MapStr{common.BKDBIN: req.IDs},
		},
		DisableCounter: true,
	}
	resp, err := s.Engine.CoreAPI.CoreService().Kube().ListNamespace(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("list namespace failed, bizID: %d, data: %v, err: %v, rid: %s", bizID, req, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(resp.Data) == 0 {
		blog.Errorf("no namespace founded, bizID: %d, query: %+v, rid: %s", bizID, query, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	for _, namespace := range resp.Data {
		ids := make([]int64, 0)
		if namespace.BizID != bizID {
			ids = append(ids, namespace.ID)
		}

		if len(ids) != 0 {
			blog.Errorf("namespace does not belong to this business, ids: %v, bizID: %s, rid: %s", ids, bizID,
				ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, ids))
			return
		}
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Engine.CoreAPI.CoreService().Kube().UpdateNamespace(ctx.Kit.Ctx, ctx.Kit.Header, bizID, req)
		if err != nil {
			blog.Errorf("update namespace failed, data: %v, err: %v, rid: %s", req, err, ctx.Kit.Rid)
			return err
		}

		audit := auditlog.NewKubeAudit(s.Engine.CoreAPI.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
		updateFields, goErr := mapstr.Struct2Map(req.Data)
		if goErr != nil {
			blog.Errorf("update fields convert failed, err: %v, rid: %s", goErr, ctx.Kit.Rid)
			return goErr
		}
		auditParam.WithUpdateFields(updateFields)
		auditLogs, err := audit.GenerateNamespaceAuditLog(auditParam, resp.Data)
		if err != nil {
			blog.Errorf("generate audit log failed, data: %v, err: %v, rid: %s", resp.Data, err, ctx.Kit.Rid)
			return err
		}
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed, data: %v, err: %v, rid: %s", resp.Data, err, ctx.Kit.Rid)
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

// DeleteNamespace delete namespace
func (s *Service) DeleteNamespace(ctx *rest.Contexts) {

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	req := new(types.NsDeleteOption)
	if err := ctx.DecodeInto(req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: req.IDs}},
	}

	resp, err := s.Engine.CoreAPI.CoreService().Kube().ListNamespace(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("list namespace failed, bizID: %d, data: %v, err: %v, rid: %s", bizID, req, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(resp.Data) == 0 {
		ctx.RespEntity(nil)
		return
	}

	for _, namespace := range resp.Data {
		ids := make([]int64, 0)
		if namespace.BizID != bizID {
			ids = append(ids, namespace.ID)
		}

		if len(ids) != 0 {
			blog.Errorf("namespace does not belong to this business, ids: %v, bizID: %d, rid: %s", ids, bizID,
				ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, ids))
			return
		}
	}

	hasRes, err := s.hasNextLevelResource(ctx.Kit, types.KubeNamespace, bizID, req.IDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	if hasRes {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Engine.CoreAPI.CoreService().Kube().DeleteNamespace(ctx.Kit.Ctx, ctx.Kit.Header, bizID, req)
		if err != nil {
			blog.Errorf("delete namespace failed, data: %v, err: %v, rid: %s", req, err, ctx.Kit.Rid)
			return err
		}

		// audit log.
		audit := auditlog.NewKubeAudit(s.Engine.CoreAPI.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditDelete)
		auditLogs, err := audit.GenerateNamespaceAuditLog(auditParam, resp.Data)
		if err != nil {
			blog.Errorf("generate audit log failed, data: %v, err: %v, rid: %s", resp.Data, err, ctx.Kit.Rid)
			return err
		}
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed, data: %v, err: %v, rid: %s", resp.Data, err, ctx.Kit.Rid)
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

// ListNamespace list namespace
func (s *Service) ListNamespace(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	req := new(types.NsQueryOption)
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
		blog.Errorf("build query namespace condition failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if req.Page.EnableCount {
		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			types.BKTableNameBaseNamespace, []map[string]interface{}{cond})
		if err != nil {
			blog.Errorf("count namespace failed, cond: %v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
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
	resp, err := s.Engine.CoreAPI.CoreService().Kube().ListNamespace(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("list namespace failed, bizID: %s, data: %v, err: %v, rid: %s", bizID, req, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithCount(0, resp.Data)
}
