/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// BindAgent bind gse agent to host, if the host has already bound another agent, change to this one.
func (s *Service) BindAgent(ctx *rest.Contexts) {
	input := new(metadata.BindAgentParam)
	if err := ctx.DecodeInto(&input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	rawErr := input.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		// check if host exists, get host data for audit log
		hostCond := mapstr.MapStr{
			common.BKHostIDField: input.HostID,
		}

		hosts, ccErr := s.Logic.GetHostInfoByConds(ctx.Kit, hostCond)
		if ccErr != nil {
			return ccErr
		}

		if len(hosts) != 1 {
			blog.Errorf("host to bind agent is not exist, input: %+v, rid: %s", input, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
		}

		if util.GetStrByInterface(hosts[0][common.BKAgentIDField]) == input.AgentID {
			return nil
		}

		// check if agent id has not already been bind to another host
		agentIDCond := []map[string]interface{}{{
			common.BKAgentIDField: input.AgentID,
		}}
		counts, ccErr := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKTableNameBaseHost, agentIDCond)
		if ccErr != nil {
			return ccErr
		}

		if len(counts) != 1 || counts[0] > 0 {
			blog.Errorf("agent is bind to another host, cannot bind again, input: %+v, rid: %s", input, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKAgentIDField)
		}

		// generate audit log.
		updateData := mapstr.MapStr{
			common.BKAgentIDField: input.AgentID,
		}

		genAuditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
		genAuditParam.WithUpdateFields(updateData)
		audit := auditlog.NewHostAudit(s.CoreAPI.CoreService())

		auditLog, err := audit.GenerateAuditLog(genAuditParam, 0, hosts)
		if err != nil {
			return err
		}

		// bind gse agent id to host
		opt := &metadata.UpdateOption{
			Condition:  hostCond,
			Data:       updateData,
			CanEditAll: true,
		}
		_, err = s.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDHost,
			opt)
		if err != nil {
			blog.Errorf("bind host agent failed, input: %+v, opt: %+v, err: %v, rid: %s", input, opt, err, ctx.Kit.Rid)
			return err
		}

		// save audit log.
		if err := audit.SaveAuditLog(ctx.Kit, auditLog...); err != nil {
			blog.Errorf("add hosts audit(%v) failed, err: %v, rid: %s", auditLog, err, ctx.Kit.Rid)
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

// UnbindAgent unbind gse agent to host, if the host is not bound to the agent, returns error.
func (s *Service) UnbindAgent(ctx *rest.Contexts) {
	input := new(metadata.UnbindAgentParam)
	if err := ctx.DecodeInto(&input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	rawErr := input.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		// check if host is bound to the agent, get host data for audit log
		hostCond := mapstr.MapStr{
			common.BKHostIDField: input.HostID,
			common.BKDBOR: []map[string]interface{}{{
				common.BKAgentIDField: map[string]interface{}{common.BKDBExists: false},
			}, {
				common.BKAgentIDField: map[string]interface{}{common.BKDBIN: []string{"", input.AgentID}},
			}},
		}

		hosts, ccErr := s.Logic.GetHostInfoByConds(ctx.Kit, hostCond)
		if ccErr != nil {
			return ccErr
		}

		if len(hosts) != 1 {
			blog.Errorf("host is not bound to the agent, input: %+v, rid: %s", input, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
		}

		if util.GetStrByInterface(hosts[0][common.BKAgentIDField]) == "" {
			return nil
		}

		// generate audit log.
		updateData := mapstr.MapStr{
			common.BKAgentIDField: "",
		}

		genAuditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
		genAuditParam.WithUpdateFields(updateData)
		audit := auditlog.NewHostAudit(s.CoreAPI.CoreService())

		auditLog, err := audit.GenerateAuditLog(genAuditParam, 0, hosts)
		if err != nil {
			return err
		}

		// unbind gse agent id to host
		opt := &metadata.UpdateOption{
			Condition:  hostCond,
			Data:       updateData,
			CanEditAll: true,
		}
		_, err = s.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDHost,
			opt)
		if err != nil {
			blog.Errorf("bind host agent failed, input: %+v, opt: %+v, err: %v, rid: %s", input, opt, err, ctx.Kit.Rid)
			return err
		}

		// save audit log.
		if err := audit.SaveAuditLog(ctx.Kit, auditLog...); err != nil {
			blog.Errorf("add hosts audit(%v) failed, err: %v, rid: %s", auditLog, err, ctx.Kit.Rid)
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
