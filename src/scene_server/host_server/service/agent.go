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
	if err := ctx.DecodeInto(&input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// validate bind agent input
	rawErr := input.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	hostMap, err := s.validateHostByAgentRelations(ctx.Kit, input.List, nil)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		audit := auditlog.NewHostAudit(s.CoreAPI.CoreService())
		auditLogs := make([]metadata.AuditLog, 0)

		for _, param := range input.List {
			host := hostMap[param.HostID]

			if util.GetStrByInterface(host[common.BKAgentIDField]) == param.AgentID {
				continue
			}

			// generate audit log.
			updateData := mapstr.MapStr{
				common.BKAgentIDField: param.AgentID,
			}

			genAuditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
			genAuditParam.WithUpdateFields(updateData)
			auditLog, err := audit.GenerateAuditLog(genAuditParam, 0, []mapstr.MapStr{host})
			if err != nil {
				return err
			}

			auditLogs = append(auditLogs, auditLog...)

			// bind gse agent id to host
			opt := &metadata.UpdateOption{
				Condition:  mapstr.MapStr{common.BKHostIDField: param.HostID},
				Data:       updateData,
				CanEditAll: true,
			}
			_, err = s.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header,
				common.BKInnerObjIDHost, opt)
			if err != nil {
				blog.Errorf("bind host agent failed, opt: %+v, err: %v, rid: %s", input, opt, err, ctx.Kit.Rid)
				return err
			}
		}

		// save audit log.
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("add hosts audit(%v) failed, err: %v, rid: %s", auditLogs, err, ctx.Kit.Rid)
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
	if err := ctx.DecodeInto(&input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// validate unbind agent input
	rawErr := input.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// host to unbind agent must have the same agent id or have no agent id at all
	agentIDs := make([]string, len(input.List))
	for index, param := range input.List {
		agentIDs[index] = param.AgentID
	}
	agentIDs = append(agentIDs, "")

	hostCond := mapstr.MapStr{
		common.BKDBOR: []map[string]interface{}{{
			common.BKAgentIDField: map[string]interface{}{common.BKDBExists: false},
		}, {
			common.BKAgentIDField: map[string]interface{}{common.BKDBIN: agentIDs},
		}},
	}
	hostMap, err := s.validateHostByAgentRelations(ctx.Kit, input.List, hostCond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		audit := auditlog.NewHostAudit(s.CoreAPI.CoreService())
		auditLogs := make([]metadata.AuditLog, 0)

		for _, param := range input.List {
			host := hostMap[param.HostID]

			// skip the host that is not bound to any agent, and returns error if the host is bound to another agent
			agentID := util.GetStrByInterface(host[common.BKAgentIDField])
			if agentID == "" {
				continue
			}

			if agentID != param.AgentID {
				return ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAgentIDField)
			}

			// generate audit log.
			updateData := mapstr.MapStr{
				common.BKAgentIDField: "",
			}

			genAuditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
			genAuditParam.WithUpdateFields(updateData)

			auditLog, err := audit.GenerateAuditLog(genAuditParam, 0, []mapstr.MapStr{host})
			if err != nil {
				return err
			}
			auditLogs = append(auditLogs, auditLog...)

			// unbind gse agent id to host
			opt := &metadata.UpdateOption{
				Condition:  mapstr.MapStr{common.BKHostIDField: param.HostID},
				Data:       updateData,
				CanEditAll: true,
			}
			_, err = s.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header,
				common.BKInnerObjIDHost, opt)
			if err != nil {
				blog.Errorf("unbind host agent failed, opt: %+v, err: %v, rid: %s", opt, err, ctx.Kit.Rid)
				return err
			}
		}

		// save audit log.
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("add hosts audit(%v) failed, err: %v, rid: %s", auditLogs, err, ctx.Kit.Rid)
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

// validateHostByAgentRelations check if hosts that needs to bind/unbind agent exist, get host data for audit log
func (s *Service) validateHostByAgentRelations(kit *rest.Kit, relations []metadata.HostAgentRelation,
	hostCond mapstr.MapStr) (map[int64]mapstr.MapStr, error) {

	hostIDs := make([]int64, len(relations))
	for index, relation := range relations {
		hostIDs[index] = relation.HostID
	}

	if hostCond == nil {
		hostCond = make(mapstr.MapStr)
	}
	hostCond[common.BKHostIDField] = mapstr.MapStr{common.BKDBIN: hostIDs}

	hosts, ccErr := s.Logic.GetHostInfoByConds(kit, hostCond)
	if ccErr != nil {
		return nil, ccErr
	}

	if len(hosts) != len(hostIDs) {
		blog.Errorf("not all host is exist, relations: %+v, rid: %s", relations, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrHostNotFound)
	}

	hostMap := make(map[int64]mapstr.MapStr)
	for _, host := range hosts {
		hostID, err := host.Int64(common.BKHostIDField)
		if err != nil {
			blog.Errorf("parse host id %v failed, err: %v, rid: %s", host[common.BKHostIDField], err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
		}
		hostMap[hostID] = host
	}

	return hostMap, nil
}
