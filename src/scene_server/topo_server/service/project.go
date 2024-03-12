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
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// CreateProject create project
func (s *Service) CreateProject(ctx *rest.Contexts) {
	opt := new(metadata.CreateProjectOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	var ids []int64
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		ids, err = s.Logics.ProjectOperation().CreateProject(ctx.Kit, opt.Data)
		if err != nil {
			blog.Errorf("create project failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(metadata.ProjectDataResp{
		IDs: ids,
	})
}

// UpdateProject update project
func (s *Service) UpdateProject(ctx *rest.Contexts) {
	opt := new(metadata.UpdateProjectOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	projectFilter := mapstr.MapStr{
		common.BKFieldID: mapstr.MapStr{common.BKDBIN: opt.IDs},
	}
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logics.InstOperation().UpdateInst(ctx.Kit, projectFilter, opt.Data, common.BKInnerObjIDProject)
		if err != nil {
			blog.Errorf("update project failed, err: %v, filter: %v, data: %v, rid: %s", err, projectFilter, opt,
				ctx.Kit.Rid)
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

// SearchProject search project
func (s *Service) SearchProject(ctx *rest.Contexts) {
	opt := new(metadata.SearchProjectOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond := mapstr.MapStr{}
	if opt.Filter != nil {
		filterCond, key, err := opt.Filter.ToMgo()
		if err != nil {
			blog.Errorf("filter to mongo failed, err: %v, filter: %v, rid: %s", err, opt.Filter, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, fmt.Sprintf("filter.%s", key)))
			return
		}
		cond = filterCond
	}

	if opt.Page.Sort == "" {
		opt.Page.Sort = common.BKFieldID
	}

	query := &metadata.QueryCondition{
		Condition:     cond,
		TimeCondition: opt.TimeCondition,
		Page:          opt.Page,
		Fields:        opt.Fields,
	}
	if !opt.Page.EnableCount {
		query.DisableCounter = true
	}

	res, err := s.Logics.InstOperation().FindInst(ctx.Kit, common.BKInnerObjIDProject, query)
	if err != nil {
		blog.Errorf("failed to find the project, err: %v, query: %v, rid: %s", err, query, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if opt.Page.EnableCount {
		ctx.RespEntityWithCount(int64(res.Count), make([]mapstr.MapStr, 0))
		return
	}

	if len(res.Info) == 0 {
		ctx.RespEntityWithCount(0, []mapstr.MapStr{})
		return
	}
	ctx.RespEntity(res)
}

// DeleteProject delete project
func (s *Service) DeleteProject(ctx *rest.Contexts) {
	opt := new(metadata.DeleteProjectOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logics.InstOperation().DeleteInstByInstID(ctx.Kit, common.BKInnerObjIDProject, opt.IDs, false)
		if err != nil {
			blog.Errorf("delete project failed, ids: %v, err: %v, rid: %s", opt.IDs, err, ctx.Kit.Rid)
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

// UpdateProjectID 更新bk_project_id, 此接口为BCS进行项目数据迁移时的专用接口，其他平台不可使用
func (s *Service) UpdateProjectID(ctx *rest.Contexts) {
	opt := new(metadata.UpdateProjectIDOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKFieldID: opt.ID,
		},
		DisableCounter: true,
	}
	resp, err := s.Logics.InstOperation().FindInst(ctx.Kit, common.BKInnerObjIDProject, query)
	if err != nil {
		blog.Errorf("failed to find the project, err: %v, query: %v, rid: %s", err, query, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(resp.Info) == 0 {
		ctx.RespEntity(nil)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		ccErr := s.Engine.CoreAPI.CoreService().Project().UpdateProjectID(ctx.Kit.Ctx, ctx.Kit.Header, opt)
		if ccErr != nil {
			blog.Errorf("update project bk_project_id failed, err: %v, opt: %v, rid: %s", ccErr, opt, ctx.Kit.Rid)
			return ccErr
		}

		audit := auditlog.NewInstanceAudit(s.Engine.CoreAPI.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
		auditParam.WithUpdateFields(map[string]interface{}{common.BKProjectIDField: opt.ProjectID})
		auditLogs, err := audit.GenerateAuditLog(auditParam, common.BKInnerObjIDProject, resp.Info)
		if err != nil {
			blog.Errorf("generate audit log failed, err: %v, data: %v, rid: %s", err, resp.Info, ctx.Kit.Rid)
			return err
		}
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed, err: %v, data: %v, rid: %s", err, resp.Info, ctx.Kit.Rid)
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
