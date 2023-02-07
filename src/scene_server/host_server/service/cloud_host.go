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
	"sort"
	"strings"

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

	err := s.preprocessAddCloudHostParam(ctx.Kit, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	needCreateHost, createIndexMap, updateParamMap, err := s.classifyHosts(ctx.Kit, input.HostInfo)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	var cloudHostIDs []int64
	txnErr := s.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.updateCloudHost(ctx.Kit, input.BizID, updateParamMap)
		if err != nil {
			return err
		}

		createdHostIDs, err := s.addCloudHostToBiz(ctx.Kit, input.BizID, needCreateHost)
		if err != nil {
			return err
		}

		// rearrange id index
		cloudHostIDs = make([]int64, len(input.HostInfo))
		for originIndex, index := range createIndexMap {
			cloudHostIDs[originIndex] = createdHostIDs[index]
		}
		for hostID, updateParam := range updateParamMap {
			cloudHostIDs[updateParam.index] = hostID
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(metadata.RspIDs{IDs: cloudHostIDs})
}

type updateCloudHostParams struct {
	currentHost mapstr.MapStr
	updateHost  mapstr.MapStr
	index       int
}

// preprocessAddCloudHostParam preprocess AddCloudHostToBizParam, validate the input and set cloud host identifier
func (s *Service) preprocessAddCloudHostParam(kit *rest.Kit, input *metadata.AddCloudHostToBizParam) error {
	if rawErr := input.Validate(); rawErr.ErrCode != 0 {
		return rawErr.ToCCError(kit.CCError)
	}

	// set cloud host identifier for the hosts and get the cloud area ids to validate if all host cloud areas are valid
	cloudAreaIDs := make([]int64, 0)
	for _, host := range input.HostInfo {
		host[common.BKCloudHostIdentifierField] = true

		cloudAreaID, err := util.GetInt64ByInterface(host[common.BKCloudIDField])
		if err != nil {
			blog.Error("host(%#v) cloud id field is invalid, err: %v, rid: %s", host, err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "host_info")
		}

		cloudAreaIDs = append(cloudAreaIDs, cloudAreaID)
	}

	if len(cloudAreaIDs) > 0 {
		cloudAreaIDs = util.IntArrayUnique(cloudAreaIDs)
		cloudAreaCond := mapstr.MapStr{common.BKCloudIDField: mapstr.MapStr{common.BKDBIN: cloudAreaIDs}}

		counts, err := s.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
			common.BKTableNameBasePlat, []map[string]interface{}{cloudAreaCond})
		if err != nil {
			blog.Error("get cloud area count failed, cond: %+v, err: %v, rid: %s", cloudAreaCond, err, kit.Rid)
			return err
		}

		if len(counts) != 1 || int(counts[0]) != len(cloudAreaIDs) {
			blog.Error("host cloud areas are invalid, input: %+v, rid: %s", input, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "host_info")
		}
	}

	return nil
}

// classifyHosts classify hosts to create and update ones
func (s *Service) classifyHosts(kit *rest.Kit, hosts []mapstr.MapStr) ([]mapstr.MapStr, map[int]int,
	map[int64]updateCloudHostParams, error) {

	// get the host ids that are already added, generate host innerIP+cloudID to host info map
	hostCond := make([]mapstr.MapStr, len(hosts))
	hostMap := make(map[string]classifyHostParams)
	for idx, host := range hosts {
		cloudAreaID, err := util.GetInt64ByInterface(host[common.BKCloudIDField])
		if err != nil {
			blog.Error("host(%#v) cloud id field is invalid, err: %v, rid: %s", host, err, kit.Rid)
			return nil, nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKCloudIDField)
		}

		innerIP := util.GetStrByInterface(host[common.BKHostInnerIPField])
		if len(innerIP) == 0 {
			blog.Error("host(%#v) inner ip field is empty, rid: %s", host, kit.Rid)
			return nil, nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKHostInnerIPField)
		}
		innerIPArr := strings.Split(innerIP, ",")

		hostCond[idx] = mapstr.MapStr{
			common.BKCloudIDField:     cloudAreaID,
			common.BKHostInnerIPField: mapstr.MapStr{common.BKDBIN: innerIPArr},
		}

		hostMap[uniqueHostKey(cloudAreaID, innerIPArr)] = classifyHostParams{
			host:  host,
			index: idx,
		}
	}

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKDBOR: hostCond},
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}
	res, err := s.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDHost, query)
	if err != nil {
		blog.Errorf("get exist hosts failed, err: %v, query: %#v, rid: %s", err, query, kit.Rid)
		return nil, nil, nil, err
	}

	// cross compare new hosts and exist hosts to classify hosts into need create and need update
	updateParamMap := make(map[int64]updateCloudHostParams)
	for _, existHost := range res.Info {
		cloudAreaID, err := util.GetInt64ByInterface(existHost[common.BKCloudIDField])
		if err != nil {
			blog.Error("exist host(%#v) cloud id field is invalid, err: %v, rid: %s", existHost, err, kit.Rid)
			return nil, nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKCloudIDField)
		}

		hostID, err := util.GetInt64ByInterface(existHost[common.BKHostIDField])
		if err != nil {
			blog.Error("exist host(%#v) id field is invalid, err: %v, rid: %s", existHost, err, kit.Rid)
			return nil, nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKHostIDField)
		}

		innerIP := util.GetStrByInterface(existHost[common.BKHostInnerIPField])
		innerIPArr := strings.Split(innerIP, ",")

		hostKey := uniqueHostKey(cloudAreaID, innerIPArr)

		updateData, exists := hostMap[hostKey]
		if !exists {
			blog.Error("exist host(%#v) has no matching update data, rid: %s", existHost, kit.Rid)
			return nil, nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "host_info")
		}

		updateParamMap[hostID] = updateCloudHostParams{
			currentHost: existHost,
			updateHost:  updateData.host,
			index:       updateData.index,
		}
		delete(hostMap, hostKey)
	}

	needCreateHosts := make([]mapstr.MapStr, 0)
	createIndexMap := make(map[int]int)
	idx := 0
	for _, hostInfo := range hostMap {
		needCreateHosts = append(needCreateHosts, hostInfo.host)
		createIndexMap[hostInfo.index] = idx
		idx++
	}

	return needCreateHosts, createIndexMap, updateParamMap, nil
}

func uniqueHostKey(cloudID int64, innerIP []string) string {
	sort.Strings(innerIP)
	return fmt.Sprintf("%d-%v", cloudID, innerIP)
}

type classifyHostParams struct {
	host  mapstr.MapStr
	index int
}

// addCloudHostToBiz add cloud host to biz idle module
func (s *Service) addCloudHostToBiz(kit *rest.Kit, bizID int64, hosts []mapstr.MapStr) ([]int64, error) {
	if len(hosts) == 0 {
		return make([]int64, 0), nil
	}

	// get the idle module ID of the biz
	moduleCond := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKAppIDField:   bizID,
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
		return nil, err
	}

	if len(moduleRes.Info) != 1 {
		blog.Errorf("biz idle module count is not one, cond: %+v, err: %v, rid: %s", moduleCond, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField)
	}

	moduleID, err := util.GetInt64ByInterface(moduleRes.Info[0][common.BKModuleIDField])
	if err != nil {
		blog.Errorf("parse module id failed, err: %v, module: %+v, rid: %s", moduleRes.Info[0], err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKModuleIDField)
	}

	hostRes, err := s.CoreAPI.CoreService().Instance().CreateManyInstance(kit.Ctx, kit.Header, common.BKInnerObjIDHost,
		&metadata.CreateManyModelInstance{Datas: hosts})
	if err != nil {
		blog.Errorf("create hosts failed, input: %#v, err: %v, rid: %s", hosts, err, kit.Rid)
		return nil, err
	}

	if len(hostRes.Repeated) > 0 || len(hostRes.Exceptions) > 0 || len(hostRes.Created) == 0 {
		blog.Errorf("create hosts failed, res: %#v, input: %#v, err: %v, rid: %s", hostRes, hosts, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "host_info")
	}

	createdHostIDs := make([]int64, len(hosts))
	for _, created := range hostRes.Created {
		createdHostIDs[created.OriginIndex] = int64(created.ID)
	}

	transReq := &metadata.TransferHostToInnerModule{
		ApplicationID: bizID,
		ModuleID:      moduleID,
		HostID:        createdHostIDs,
	}

	res, err := s.CoreAPI.CoreService().Host().TransferToInnerModule(kit.Ctx, kit.Header, transReq)
	if err != nil {
		blog.Errorf("transfer hosts failed, req: %+v, res: %+v, err: %v, rid: %s", res, transReq, err, kit.Rid)
		return nil, err
	}

	// generate and save audit logs
	audit := auditlog.NewHostAudit(s.CoreAPI.CoreService())
	auditParam := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditCond := mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: createdHostIDs}}

	auditLogs, rawErr := audit.GenerateAuditLogByCond(auditParam, bizID, auditCond)
	if rawErr != nil {
		return nil, rawErr
	}

	if err = audit.SaveAuditLog(kit, auditLogs...); err != nil {
		return nil, err
	}

	return createdHostIDs, nil
}

// updateCloudHost update cloud host info
func (s *Service) updateCloudHost(kit *rest.Kit, bizID int64, updateHostMap map[int64]updateCloudHostParams) error {

	if len(updateHostMap) == 0 {
		return nil
	}

	updateHostIDs := make([]int64, 0)
	for hostID := range updateHostMap {
		updateHostIDs = append(updateHostIDs, hostID)
	}

	// validate if all exist hosts are in the correct biz
	validateOpt := &metadata.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: []int64{bizID},
		HostIDArr:        updateHostIDs,
	}
	hostIDs, err := s.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(kit.Ctx, kit.Header, validateOpt)
	if err != nil {
		blog.Errorf("validate host in biz failed, err: %v, opt: %+v, rid: %s", err, validateOpt, kit.Rid)
		return err
	}

	if len(hostIDs) != len(updateHostMap) {
		blog.Errorf("not all hosts are in biz %d, valid ids: %+v, all ids: %+v, rid: %s", bizID, hostIDs,
			updateHostIDs, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "host_info")
	}

	audit := auditlog.NewHostAudit(s.CoreAPI.CoreService())
	auditLogs := make([]metadata.AuditLog, 0)

	for hostID, param := range updateHostMap {
		// generator audit log
		auditParam := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(
			param.updateHost)
		logs, rawErr := audit.GenerateAuditLog(auditParam, bizID, []mapstr.MapStr{param.currentHost})
		if rawErr != nil {
			return rawErr
		}
		auditLogs = append(auditLogs, logs...)

		// update host
		opt := &metadata.UpdateOption{
			Condition:  mapstr.MapStr{common.BKHostIDField: hostID},
			Data:       param.updateHost,
			CanEditAll: true,
		}
		_, err := s.CoreAPI.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDHost, opt)
		if err != nil {
			blog.Errorf("update host failed, err: %v, opt: %+v, rid: %s", err, opt, kit.Rid)
			return err
		}
	}

	// save audit log
	if err := audit.SaveAuditLog(kit, auditLogs...); err != nil {
		return err
	}

	return nil
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
