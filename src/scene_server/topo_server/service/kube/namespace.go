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
	acmeta "configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// CreateNamespace create namespace
func (s *service) CreateNamespace(ctx *rest.Contexts) {
	req := new(types.NsCreateOption)
	if err := ctx.DecodeInto(req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeNamespace, Action: acmeta.Create},
		BusinessID: req.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	var data *metadata.RspIDs
	var err error
	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		data, err = s.ClientSet.CoreService().Kube().CreateNamespace(ctx.Kit.Ctx, ctx.Kit.Header, req.Data)
		if err != nil {
			blog.Errorf("create namespace failed, data: %v, err: %v, rid: %s", req, err, ctx.Kit.Rid)
			return err
		}

		// audit log.
		audit := auditlog.NewKubeAudit(s.ClientSet.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditCreate)
		for idx := range req.Data {
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
func (s *service) UpdateNamespace(ctx *rest.Contexts) {
	req := new(types.NsUpdateOption)
	if err := ctx.DecodeInto(req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeNamespace, Action: acmeta.Update},
		BusinessID: req.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKFieldID: mapstr.MapStr{common.BKDBIN: req.IDs},
		},
		DisableCounter: true,
	}
	resp, err := s.ClientSet.CoreService().Kube().ListNamespace(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("list namespace failed, bizID: %d, data: %v, err: %v, rid: %s", req.BizID, req, err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(resp.Data) == 0 {
		blog.Errorf("no namespace founded, bizID: %d, query: %+v, rid: %s", req.BizID, query, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// checks if namespace is a shared namespace and if its biz id is not the same with the input biz id
	mismatchIDs := make([]int64, 0)
	for _, namespace := range resp.Data {
		if namespace.BizID != req.BizID {
			mismatchIDs = append(mismatchIDs, namespace.ID)
		}
	}

	if len(mismatchIDs) > 0 {
		mismatchNsMap := map[int64][]int64{req.BizID: mismatchIDs}
		if err := s.Logics.KubeOperation().CheckPlatBizSharedNs(ctx.Kit, mismatchNsMap); err != nil {
			ctx.RespAutoError(err)
			return
		}
	}

	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.ClientSet.CoreService().Kube().UpdateNamespace(ctx.Kit.Ctx, ctx.Kit.Header, &req.NsUpdateByIDsOption)
		if err != nil {
			blog.Errorf("update namespace failed, data: %v, err: %v, rid: %s", req, err, ctx.Kit.Rid)
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
func (s *service) DeleteNamespace(ctx *rest.Contexts) {
	req := new(types.NsDeleteOption)
	if err := ctx.DecodeInto(req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeNamespace, Action: acmeta.Delete},
		BusinessID: req.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: req.IDs}},
	}

	resp, err := s.ClientSet.CoreService().Kube().ListNamespace(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("list namespace failed, bizID: %d, data: %v, err: %v, rid: %s", req.BizID, req, err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(resp.Data) == 0 {
		ctx.RespEntity(nil)
		return
	}

	// checks if namespace is a shared namespace and if its biz id is not the same with the input biz id
	mismatchIDs := make([]int64, 0)
	for _, namespace := range resp.Data {
		if namespace.BizID != req.BizID {
			mismatchIDs = append(mismatchIDs, namespace.ID)
		}
	}

	if len(mismatchIDs) > 0 {
		mismatchNsMap := map[int64][]int64{req.BizID: mismatchIDs}
		if err := s.Logics.KubeOperation().CheckPlatBizSharedNs(ctx.Kit, mismatchNsMap); err != nil {
			ctx.RespAutoError(err)
			return
		}
	}

	hasRes, rawErr := s.hasNextLevelResource(ctx.Kit, types.KubeNamespace, req.IDs)
	if rawErr != nil {
		ctx.RespAutoError(rawErr)
		return
	}
	if hasRes {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID))
		return
	}

	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.ClientSet.CoreService().Kube().DeleteNamespace(ctx.Kit.Ctx, ctx.Kit.Header, &req.NsDeleteByIDsOption)
		if err != nil {
			blog.Errorf("delete namespace failed, data: %v, err: %v, rid: %s", req, err, ctx.Kit.Rid)
			return err
		}

		// audit log.
		audit := auditlog.NewKubeAudit(s.ClientSet.CoreService())
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
func (s *service) ListNamespace(ctx *rest.Contexts) {
	req := new(types.NsQueryOption)
	if err := ctx.DecodeInto(req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubeNamespace, Action: acmeta.Find},
		BusinessID: req.BizID}
	if resp, authorized := s.AuthManager.Authorize(ctx.Kit, authRes); !authorized {
		ctx.RespNoAuth(resp)
		return
	}

	// compatible for shared cluster scenario
	cond, err := s.Logics.KubeOperation().GenSharedNsListCond(ctx.Kit, types.KubeNamespace, req.BizID, req.Filter)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if req.Page.EnableCount {
		counts, err := s.ClientSet.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
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
	resp, err := s.ClientSet.CoreService().Kube().ListNamespace(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("list namespace failed, bizID: %s, data: %v, err: %v, rid: %s", req.BizID, req, err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithCount(0, resp.Data)
}
