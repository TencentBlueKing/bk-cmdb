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
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// AddCloudHostToBiz add cloud host to biz idle module
func (s *Service) AddCloudHostToBiz(ctx *rest.Contexts) {
	input := new(metadata.AddCloudHostToBizParam)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	moduleID, err := s.preprocessAddCloudHostParam(ctx.Kit, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	var createdHostIDs []int64
	txnErr := s.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		hostRes, err := s.CoreAPI.CoreService().Instance().CreateManyInstance(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKInnerObjIDHost, &metadata.CreateManyModelInstance{Datas: input.HostInfo})
		if err != nil {
			blog.Errorf("create hosts failed, input: %#v, err: %v, rid: %s", input, err, ctx.Kit.Rid)
			return err
		}

		if len(hostRes.Repeated) > 0 || len(hostRes.Exceptions) > 0 {
			blog.Errorf("create hosts failed, res: %#v, input: %#v, err: %v, rid: %s", hostRes, input, err, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "host_info")
		}

		for _, created := range hostRes.Created {
			createdHostIDs = append(createdHostIDs, int64(created.ID))
		}

		transReq := &metadata.TransferHostToInnerModule{
			ApplicationID: input.BizID,
			ModuleID:      moduleID,
			HostID:        createdHostIDs,
		}

		res, err := s.CoreAPI.CoreService().Host().TransferToInnerModule(ctx.Kit.Ctx, ctx.Kit.Header, transReq)
		if err != nil {
			blog.Errorf("transfer hosts failed, req: %+v, res: %+v, err: %v, rid: %s", res, transReq, err, ctx.Kit.Rid)
			return err
		}

		// generate and save audit logs
		audit := auditlog.NewHostAudit(s.CoreAPI.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditCreate)
		auditCond := mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: createdHostIDs}}

		auditLogs, rawErr := audit.GenerateAuditLogByCond(auditParam, input.BizID, auditCond)
		if rawErr != nil {
			return rawErr
		}

		if err = audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			return err
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(metadata.RspIDs{IDs: createdHostIDs})
}

// preprocessAddCloudHostParam preprocess AddCloudHostToBizParam, returns biz idle module id to add to
func (s *Service) preprocessAddCloudHostParam(kit *rest.Kit, input *metadata.AddCloudHostToBizParam) (int64, error) {
	if rawErr := input.Validate(); rawErr.ErrCode != 0 {
		return 0, rawErr.ToCCError(kit.CCError)
	}

	// set cloud host identifier for the hosts and get the cloud area ids to validate if all host cloud areas are valid
	cloudAreaIDs := make([]int64, 0)
	for _, host := range input.HostInfo {
		host[common.BKCloudHostIdentifierField] = true

		cloudAreaID, err := util.GetInt64ByInterface(host[common.BKCloudIDField])
		if err != nil {
			blog.Error("host(%#v) cloud id field is invalid, err: %v, rid: %s", host, err, kit.Rid)
			return 0, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "host_info")
		}

		if cloudAreaID != common.BKDefaultDirSubArea {
			cloudAreaIDs = append(cloudAreaIDs, cloudAreaID)
		}
	}

	if len(cloudAreaIDs) > 0 {
		cloudAreaIDs = util.IntArrayUnique(cloudAreaIDs)
		cloudAreaCond := mapstr.MapStr{common.BKCloudIDField: mapstr.MapStr{common.BKDBIN: cloudAreaIDs}}

		counts, err := s.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
			common.BKTableNameBasePlat, []map[string]interface{}{cloudAreaCond})
		if err != nil {
			blog.Error("get cloud area count failed, cond: %+v, err: %v, rid: %s", cloudAreaCond, err, kit.Rid)
			return 0, err
		}

		if len(counts) != 1 || int(counts[0]) != len(cloudAreaIDs) {
			blog.Error("host cloud areas are invalid, input: %+v, rid: %s", input, kit.Rid)
			return 0, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "host_info")
		}
	}

	// get the idle module ID
	moduleCond := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKAppIDField:   input.BizID,
			common.BKDefaultField: common.DefaultResModuleFlag,
		},
		Fields:         []string{common.BKModuleIDField},
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		DisableCounter: true,
	}

	moduleRes, err := s.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		moduleCond)
	if err != nil {
		blog.Errorf("get idle module ID failed, cond: %+v, err: %v, rid: %s", moduleCond, err, kit.Rid)
		return 0, err
	}

	if len(moduleRes.Info) != 1 {
		blog.Errorf("biz idle module count is not one, cond: %+v, err: %v, rid: %s", moduleCond, err, kit.Rid)
		return 0, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField)
	}

	idleModuleID, err := util.GetInt64ByInterface(moduleRes.Info[0][common.BKModuleIDField])
	if err != nil {
		blog.Errorf("parse module id failed, err: %v, module: %+v, rid: %s", moduleRes.Info[0], err, kit.Rid)
		return 0, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKModuleIDField)
	}

	return idleModuleID, nil
}

// DeleteCloudHostFromBiz delete cloud hosts from biz
func (s *Service) DeleteCloudHostFromBiz(ctx *rest.Contexts) {
	input := new(metadata.DeleteCloudHostFromBizParam)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := s.validateDeleteCloudHostParam(ctx.Kit, input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		// delete all instance associations of the hosts
		asstCond := &metadata.InstAsstDeleteOption{
			Opt: metadata.DeleteOption{
				Condition: mapstr.MapStr{
					common.BKDBOR: []mapstr.MapStr{{
						common.BKObjIDField:  common.BKInnerObjIDHost,
						common.BKInstIDField: mapstr.MapStr{common.BKDBIN: input.HostIDs},
					}, {
						common.BKAsstObjIDField:  common.BKInnerObjIDHost,
						common.BKAsstInstIDField: mapstr.MapStr{common.BKDBIN: input.HostIDs},
					}},
				}},
			ObjID: common.BKInnerObjIDHost,
		}

		_, err := s.CoreAPI.CoreService().Association().DeleteInstAssociation(ctx.Kit.Ctx, ctx.Kit.Header, asstCond)
		if err != nil {
			blog.Errorf("delete host association by cond(%#v) failed, err: %v, rid: %s", asstCond, err, ctx.Kit.Rid)
			return err
		}

		// generate host audit log
		audit := auditlog.NewHostAudit(s.CoreAPI.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditDelete)
		auditCond := mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: input.HostIDs}}

		auditLogs, err := audit.GenerateAuditLogByCond(auditParam, input.BizID, auditCond)
		if err != nil {
			return err
		}

		// delete host
		delReq := &metadata.DeleteHostRequest{
			ApplicationID: input.BizID,
			HostIDArr:     input.HostIDs,
		}

		err = s.CoreAPI.CoreService().Host().DeleteHostFromSystem(ctx.Kit.Ctx, ctx.Kit.Header, delReq)
		if err != nil {
			blog.Error("delete host failed, request: %#v, err: %v, rid: %s", delReq, err, ctx.Kit.Rid)
			return err
		}

		// save audit logs
		if len(auditLogs) > 0 {
			if err = audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
				return err
			}
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// validateDeleteCloudHostParam validate DeleteCloudHostFromBizParam
func (s *Service) validateDeleteCloudHostParam(kit *rest.Kit, input *metadata.DeleteCloudHostFromBizParam) error {
	if rawErr := input.Validate(); rawErr.ErrCode != 0 {
		return rawErr.ToCCError(kit.CCError)
	}

	// check if host belongs to the idle set
	moduleCond := mapstr.MapStr{
		common.BKAppIDField:   input.BizID,
		common.BKDefaultField: mapstr.MapStr{common.BKDBNE: common.DefaultFlagDefaultValue},
	}
	if err := s.Logic.ValidateHostInModule(kit, input.HostIDs, moduleCond); err != nil {
		return err
	}

	// check if hosts are all cloud hosts
	hostCond := mapstr.MapStr{
		common.BKHostIDField:              mapstr.MapStr{common.BKDBIN: input.HostIDs},
		common.BKCloudHostIdentifierField: mapstr.MapStr{common.BKDBNE: true},
	}

	counts, err := s.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameBaseHost, []map[string]interface{}{hostCond})
	if err != nil {
		blog.Error("get normal host count failed, cond: %+v, err: %v, rid: %s", hostCond, err, kit.Rid)
		return err
	}

	if len(counts) != 1 || int(counts[0]) > 0 {
		blog.Error("host are not all cloud hosts, input: %+v, rid: %s", input, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "bk_host_ids")
	}

	return nil
}
