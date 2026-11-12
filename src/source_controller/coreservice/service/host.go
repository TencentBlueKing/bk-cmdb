/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core/instances"
	"configcenter/src/storage/driver/mongodb"
)

// TransferHostToInnerModule TODO
func (s *coreService) TransferHostToInnerModule(ctx *rest.Contexts) {
	inputData := &metadata.TransferHostToInnerModule{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	err := s.core.HostOperation().TransferToInnerModule(ctx.Kit, inputData)
	if err != nil {
		blog.Errorf("transfer host to default module failed. err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// TransferHostToNormalModule TODO
func (s *coreService) TransferHostToNormalModule(ctx *rest.Contexts) {
	inputData := &metadata.HostsModuleRelation{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	err := s.core.HostOperation().TransferToNormalModule(ctx.Kit, inputData)
	if err != nil {
		blog.Errorf("transfer host to normal module failed. err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// TransferHostToAnotherBusiness TODO
func (s *coreService) TransferHostToAnotherBusiness(ctx *rest.Contexts) {
	inputData := &metadata.TransferHostsCrossBusinessRequest{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	err := s.core.HostOperation().TransferToAnotherBusiness(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("transfer host across business failed. err: %s, input: %s, rid: %s", err.Error(), inputData,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// RemoveFromModule TODO
func (s *coreService) RemoveFromModule(ctx *rest.Contexts) {
	inputData := &metadata.RemoveHostsFromModuleOption{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	err := s.core.HostOperation().RemoveFromModule(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("remove host from module failed. err: %s, input: %s, rid: %s", err.Error(), inputData,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// GetHostModuleRelation TODO
func (s *coreService) GetHostModuleRelation(ctx *rest.Contexts) {
	inputData := &metadata.HostModuleRelationRequest{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	relationArr, err := s.core.HostOperation().GetHostModuleRelation(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("GetHostModuleRelation  error. err:%s, rid:%s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(relationArr)
}

// DeleteHostFromSystem TODO
func (s *coreService) DeleteHostFromSystem(ctx *rest.Contexts) {
	inputData := &metadata.DeleteHostRequest{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	err := s.core.HostOperation().DeleteFromSystem(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("delete host error. err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// HostIdentifier TODO
func (s *coreService) HostIdentifier(ctx *rest.Contexts) {
	inputData := &metadata.SearchHostIdentifierParam{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	hostIdentifierArr, err := s.core.HostOperation().Identifier(ctx.Kit, inputData)

	if err != nil {
		blog.InfoJSON("Identifier host identifier handle error. err: %s, input: %s, rid:%s", err.Error(), inputData,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(metadata.SearchHostIdentifierData{Info: hostIdentifierArr, Count: len(hostIdentifierArr)})
}

// GetHostByID TODO
func (s *coreService) GetHostByID(ctx *rest.Contexts) {
	hostID, err := strconv.Atoi(ctx.Request.PathParameter("bk_host_id"))
	if err != nil {
		blog.Errorf("GetHostByID failed, get host by id, but got invalid host id, hostID: %s, err: %+v, rid: %s",
			hostID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKHostIDField))
		return
	}

	result := make(metadata.HostMapStr, 0)
	condition := common.KvMap{common.BKHostIDField: hostID}
	condition = util.SetModOwner(condition, ctx.Kit.SupplierAccount)
	err = mongodb.Client().Table(common.BKTableNameBaseHost).Find(condition).One(ctx.Kit.Ctx, &result)
	// TODO: return error for not found and deal error with all callers
	if err != nil && !mongodb.Client().IsNotFoundError(err) {
		blog.Errorf("GetHostByID failed, get host by id[%d] failed, err: %+v, rid: %s", hostID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	if result[common.BKCreatedAt] == nil && result[common.CreateTimeField] != nil {
		result[common.BKCreatedAt] = result[common.CreateTimeField]
	}

	if result[common.BKUpdatedAt] == nil && result[common.LastTimeField] != nil {
		result[common.BKUpdatedAt] = result[common.LastTimeField]
	}

	ctx.RespEntity(result)
}

// GetHosts TODO
func (s *coreService) GetHosts(ctx *rest.Contexts) {
	var dat metadata.QueryInput
	if err := ctx.DecodeInto(&dat); err != nil {
		ctx.RespAutoError(err)
		return
	}

	condition := dat.Condition
	if dat.TimeCondition != nil {
		var err error
		condition, err = dat.TimeCondition.MergeTimeCondition(condition)
		if err != nil {
			blog.ErrorJSON("merge time condition failed, error: %s, input: %s, rid: %s", err, dat, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
	}

	condition = util.SetModOwner(condition, ctx.Kit.SupplierAccount)
	fieldArr := util.SplitStrField(dat.Fields, ",")

	query := mongodb.Client().Table(common.BKTableNameBaseHost).Find(condition).Sort(dat.Sort).
		Start(uint64(dat.Start)).Limit(uint64(dat.Limit))
	info, err := instances.FindInst(ctx.Kit, fieldArr, query, common.BKInnerObjIDHost)
	if err != nil {
		blog.ErrorJSON("failed to query the host, cond: %v, err: %s, rid: %s", condition, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	var finalCount uint64

	if !dat.DisableCounter {
		count, err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(condition).Count(ctx.Kit.Ctx)
		if err != nil {
			blog.Errorf("get object failed type:%s ,input: %v error: %v, rid: %s", common.BKInnerObjIDHost, dat, err,
				ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostSelectInst))
			return
		}
		finalCount = count
	}

	ctx.RespEntity(metadata.HostInfo{
		Count: int(finalCount),
		Info:  info,
	})
}

// GetDistinctHostIDsByTopoRelation get all  host ids by topology relation condition
func (s *coreService) GetDistinctHostIDsByTopoRelation(ctx *rest.Contexts) {
	inputData := &metadata.DistinctHostIDByTopoRelationRequest{}
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	hostIDArr, err := s.core.HostOperation().GetDistinctHostIDsByTopoRelation(ctx.Kit, inputData)
	if err != nil {
		blog.ErrorJSON("GetDistinctHostIDsByTopoRelation  error. err:%s, rid:%s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(metadata.DistinctID{IDArr: hostIDArr})
}

// TransferHostResourceDirectory TODO
func (s *coreService) TransferHostResourceDirectory(ctx *rest.Contexts) {
	input := &metadata.TransferHostResourceDirectory{}
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	err := s.core.HostOperation().TransferResourceDirectory(ctx.Kit, input)
	if err != nil {
		blog.ErrorJSON("TransferHostResourceDirectory  error. err:%s, input:%s, rid:%s", err.Error(), input,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// modifyHostInnerIPConfigKey the config key of the switch of the host inner ip modification api.
//
// NOTICE: it must be the same as the key used by host_server, the switch is checked here as
// well so that the api can not be used by calling coreservice directly when it is disabled.
const modifyHostInnerIPConfigKey = "hostServer.enableModifyHostInnerIP"

// UpdateHostInnerIP updates the inner ip fields of a static host.
func (s *coreService) UpdateHostInnerIP(ctx *rest.Contexts) {
	// the feature is disabled by default, the config is read on each request so that the api
	// can be closed as soon as the config is updated, no service restart is needed.
	enabled, err := cc.Bool(modifyHostInnerIPConfigKey)
	if err != nil {
		blog.Errorf("get config %s failed, err: %v, rid: %s", modifyHostInnerIPConfigKey, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostModifyInnerIPDisabled))
		return
	}
	if !enabled {
		blog.Errorf("config %s is not enabled, host inner ip modification is forbidden, rid: %s",
			modifyHostInnerIPConfigKey, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrHostModifyInnerIPDisabled))
		return
	}

	input := new(metadata.ModifyHostInnerIPOpt)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// validate and normalize the input here, so that the audit log records the same value as
	// the one written to the db, the core operation validates it again authoritatively.
	if rawErr := input.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	result, err := s.core.HostOperation().UpdateHostInnerIP(ctx.Kit, *input)
	if err != nil {
		blog.Errorf("update host %d inner ip failed, err: %v, rid: %s", input.HostID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// the host is changed, generate and save its audit log in the same transaction as the update
	if result.Changed && len(result.ChangedFields) > 0 {
		updateFields := make(mapstr.MapStr, len(result.ChangedFields))
		for _, field := range result.ChangedFields {
			if field == common.BKHostInnerIPv6Field {
				updateFields[field] = metadata.ConvertHostIPListToString(input.DestIPv6)
				continue
			}
			updateFields[field] = metadata.ConvertHostIPListToString(input.DestIPv4)
		}

		audit := auditlog.NewHostAudit(s.engine.CoreAPI.CoreService())
		genAuditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate).
			WithUpdateFields(updateFields)

		// bizID is passed as 0, the audit log generator looks up the biz id of the host itself
		// (auditlog.getBizIDByHostID), which is the same as the other host audit logs: a host
		// in the resource pool is attributed to the resource pool business, and only a host
		// without any module relation(abnormal data) is recorded with biz id 0.
		auditLogs, err := audit.GenerateAuditLog(genAuditParam, 0, []mapstr.MapStr{result.PreHost})
		if err != nil {
			blog.Errorf("generate the audit log of host %d inner ip failed, err: %v, rid: %s", input.HostID, err,
				ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrAuditSaveLogFailed))
			return
		}

		if err := s.core.AuditOperation().CreateAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save the audit log of host %d inner ip failed, err: %v, rid: %s", input.HostID, err,
				ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrAuditSaveLogFailed))
			return
		}
	}

	ctx.RespEntity(result)
}
