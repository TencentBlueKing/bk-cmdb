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

// Package idrule package
package idrule

import (
	"fmt"

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// UpdateIDRuleIncrID update id  rule self-increasing id
func (s *service) UpdateIDRuleIncrID(ctx *rest.Contexts) {
	opt := new(metadata.UpdateIDGenOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	if authResp, authorized := s.AuthManager.Authorize(ctx.Kit, meta.ResourceAttribute{Basic: meta.Basic{
		Type: meta.IDRuleIncrID, Action: meta.Update}}); !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	opt.Type = util.GetIDRule(opt.Type)
	err := s.ClientSet.CoreService().Model().UpdateIDGenerator(ctx.Kit.Ctx, ctx.Kit.Header, opt)
	if err != nil {
		blog.Errorf("update id generator failed, err: %v, opt: %+v, rid: %s", err, opt, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// SyncInstIDRule sync instance id rule field
func (s *service) SyncInstIDRule(ctx *rest.Contexts) {
	opt := new(metadata.SyncIDRuleOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// authorize
	authResp, authorized, err := s.AuthManager.HasUpdateModelInstAuth(ctx.Kit, []string{opt.ObjID})
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	if !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	cond := &metadata.QueryCondition{
		Fields: []string{common.BKPropertyIDField},
		Condition: map[string]interface{}{
			common.BKObjIDField:        opt.ObjID,
			common.BKPropertyTypeField: common.FieldTypeIDRule,
		},
	}
	attr, err := s.ClientSet.CoreService().Model().ReadModelAttr(ctx.Kit.Ctx, ctx.Kit.Header, opt.ObjID, cond)
	if err != nil {
		blog.Errorf("find attribute failed, cond: %+v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(attr.Info) != 1 {
		blog.Errorf("the num of id rule fields is not %d, cond: %+v, err: %v, rid: %s", metadata.IDRuleFieldLimit,
			cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid,
			fmt.Sprintf("the num of assetid fields is not %d", metadata.IDRuleFieldLimit)))
		return
	}

	propertyID := attr.Info[0].PropertyID
	taskID, err := s.buildTask(ctx.Kit, opt.ObjID, propertyID)
	if err != nil {
		blog.Errorf("build task failed, opt: %+v, propertyID: %s, err: %v, rid: %s", opt, propertyID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(&metadata.SyncIDRuleRes{TaskID: taskID})
}

func (s *service) buildTask(kit *rest.Kit, objID, propertyID string) (string, error) {
	idField := common.GetInstIDField(objID)
	cond := &metadata.QueryCondition{
		Fields: []string{idField},
		Condition: mapstr.MapStr{
			common.BKDBOR: []mapstr.MapStr{{propertyID: ""}, {propertyID: mapstr.MapStr{common.BKDBExists: false}}},
			idField:       mapstr.MapStr{common.BKDBGT: 0},
		},
		Page: metadata.BasePage{Limit: common.BKMaxPageSize, Sort: idField},
	}

	taskData := make([]interface{}, 0)
	data := metadata.UpdateInstIDRuleOption{ObjID: objID, PropertyID: propertyID}
	for {
		insts, err := s.ClientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, cond)
		if err != nil {
			blog.Errorf("find instance failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
			return "", err
		}
		if len(insts.Info) == 0 {
			break
		}

		var lastID int64
		for _, inst := range insts.Info {
			id, err := util.GetInt64ByInterface(inst[idField])
			if err != nil {
				blog.Errorf("get instance id failed, inst: %+v, err: %v, rid: %s", inst, err, kit.Rid)
				return "", kit.CCError.Errorf(common.CCErrCommDBSelectFailed)
			}
			data.IDs = append(data.IDs, id)
			lastID = id
		}

		taskData = append(taskData, data)
		data.IDs = make([]int64, 0)

		if len(insts.Info) < common.BKMaxPageSize {
			break
		}

		cond.Condition[idField] = mapstr.MapStr{common.BKDBGT: lastID}
	}

	modelCond := &metadata.CommonQueryOption{
		Fields: []string{common.BKFieldID},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		CommonFilterOption: metadata.CommonFilterOption{
			Filter: filtertools.GenAtomFilter(common.BKObjIDField, filter.Equal, objID),
		},
	}
	res, err := s.ClientSet.CoreService().Model().ListModel(kit.Ctx, kit.Header, modelCond)
	if err != nil {
		blog.Errorf("list object failed, cond: %+v, err: %v, opt: %+v, rid: %s", modelCond, err, kit.Rid)
		return "", err
	}
	if len(res.Info) != 1 {
		blog.Errorf("objID is invalid, objID: %s, rid: %s", objID, kit.Rid)
		return "", kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKObjIDField)
	}

	var taskID string
	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
		task, err := s.ClientSet.TaskServer().Task().Create(kit.Ctx, kit.Header, common.SyncInstIDRuleTaskFlag,
			res.Info[0].ID, taskData)
		if err != nil {
			blog.Errorf("create instance assetid sync failed, data: %+v, err: %v, rid: %s", data, err, kit.Rid)
			return err
		}

		taskID = task.TaskID
		blog.V(4).Infof("successfully create instance assetid sync task: %#v, rid: %s", task, kit.Rid)
		return nil
	})

	if txnErr != nil {
		return "", txnErr
	}

	return taskID, nil
}

// SyncInstIDRuleTask sync instance id rule field task
func (s *service) SyncInstIDRuleTask(ctx *rest.Contexts) {
	opt := new(metadata.UpdateInstIDRuleOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	txnErr := s.ClientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.ClientSet.CoreService().IDRule().UpdateInstIDRule(ctx.Kit.Ctx, ctx.Kit.Header, opt); err != nil {
			blog.Errorf("update instance id rule failed, opt: %+v, err: %v, rid: %s", opt, err, ctx.Kit.Rid)
			return err
		}

		cond := &metadata.QueryCondition{
			Page:      metadata.BasePage{Limit: common.BKNoLimit},
			Condition: mapstr.MapStr{metadata.GetInstIDFieldByObjID(opt.ObjID): mapstr.MapStr{common.BKDBIN: opt.IDs}},
		}
		insts, err := s.ClientSet.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, opt.ObjID, cond)
		if err != nil {
			blog.Errorf("find instance failed, cond: %+v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
			return err
		}

		audit := auditlog.NewInstanceAudit(s.ClientSet.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
		auditLogs := make([]metadata.AuditLog, 0)
		for _, inst := range insts.Info {
			val, exist := inst.Get(opt.PropertyID)
			if !exist {
				blog.Errorf("inst can not find property: %s, inst: %+v, rid: %s", opt.PropertyID, inst, ctx.Kit.Rid)
				return err
			}

			auditParam.WithUpdateFields(mapstr.MapStr{opt.PropertyID: val})
			delete(inst, opt.PropertyID)
			auditLog, err := audit.GenerateAuditLog(auditParam, opt.ObjID, []mapstr.MapStr{inst})
			if err != nil {
				blog.Errorf("generate audit log failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}

			auditLogs = append(auditLogs, auditLog...)
		}

		if err = audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit failed, err: %v, rid: %s", err, ctx.Kit.Rid)
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

// FindInstIDRuleTaskStatus find instance id rule task status
func (s *service) FindInstIDRuleTaskStatus(ctx *rest.Contexts) {
	opt := new(metadata.IDRuleTaskOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	resp, err := s.ClientSet.TaskServer().Task().TaskDetail(ctx.Kit.Ctx, ctx.Kit.Header, opt.TaskID)
	if err != nil {
		blog.Errorf("find id rule task status failed, task id: %s, err: %v, rid: %s", opt.TaskID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(map[string]interface{}{common.BKStatusField: resp.Info.Status})
}
