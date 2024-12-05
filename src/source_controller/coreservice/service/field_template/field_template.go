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

package fieldtmpl

import (
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// ListFieldTemplate list field templates.
func (s *service) ListFieldTemplate(cts *rest.Contexts) {
	opt := new(metadata.CommonQueryOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	filter, err := opt.ToMgo()
	if err != nil {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	if opt.Page.EnableCount {
		count, err := mongodb.Shard(cts.Kit.ShardOpts()).Table(common.BKTableNameFieldTemplate).Find(filter).
			Count(cts.Kit.Ctx)
		if err != nil {
			blog.Errorf("count field templates failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
			cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}

		cts.RespEntity(metadata.FieldTemplateInfo{Count: count})
		return
	}

	fieldTemplates := make([]metadata.FieldTemplate, 0)
	err = mongodb.Shard(cts.Kit.ShardOpts()).Table(common.BKTableNameFieldTemplate).Find(filter).
		Start(uint64(opt.Page.Start)).Limit(uint64(opt.Page.Limit)).Sort(opt.Page.Sort).Fields(opt.Fields...).
		All(cts.Kit.Ctx, &fieldTemplates)
	if err != nil {
		blog.Errorf("list field templates failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
		cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	cts.RespEntity(metadata.FieldTemplateInfo{Info: fieldTemplates})
}

func canObjBindingFieldTemplate(kit *rest.Kit, objIDs []string) error {

	cond := mapstr.MapStr{
		common.AssociationKindIDField: common.AssociationKindMainline,
		common.BKDBAND: []mapstr.MapStr{
			{
				common.BKObjIDField: mapstr.MapStr{
					common.BKDBNE: common.BKInnerObjIDHost,
				},
			},
			{
				common.BKObjIDField: mapstr.MapStr{
					common.BKDBIN: objIDs,
				},
			},
		},
	}

	count, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAsst).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("search mainline failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}

	if count > 0 {
		return kit.CCError.CCError(common.CCErrCommParamsIsInvalid)
	}
	return nil
}

func (s *service) validateTemplateID(kit *rest.Kit, id int64) error {

	cond := mapstr.MapStr{
		common.BKFieldID: id,
	}
	cnt, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameFieldTemplate).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count field template failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if cnt == 0 {
		blog.Errorf("no field template founded, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommNotFound, "field_template")
	}
	if cnt > 1 {
		blog.Errorf("multi field template founded, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommGetMultipleObject, "field_template")
	}
	return nil
}

func (s *service) getObjectByIDs(kit *rest.Kit, ids []int64) ([]string, error) {

	for _, id := range ids {
		if id == 0 {
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKFieldID)
		}
	}

	filter := mapstr.MapStr{
		common.BKFieldID: mapstr.MapStr{
			common.BKDBIN: ids,
		},
	}
	objs := make([]metadata.Object, 0)

	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjDes).Find(filter).Fields(common.BKObjIDField).
		All(kit.Ctx, &objs); err != nil {
		blog.Errorf("mongodb count failed, table: %s, err: %v, rid: %s", common.BKTableNameObjDes, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	if len(ids) != len(objs) {
		blog.Errorf("mongodb count num failed, ids len: %d, objects len: %d, rid: %s", len(ids), len(objs), kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "ids")
	}

	if len(objs) == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommNotFound, "object_ids")
	}

	objIDs := make([]string, 0)
	for _, obj := range objs {
		objIDs = append(objIDs, obj.ObjectID)
	}

	return objIDs, nil
}

// FieldTemplateBindObject field template bind model.
func (s *service) FieldTemplateBindObject(ctx *rest.Contexts) {

	opt := new(metadata.FieldTemplateBindObjOpt)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}
	kit := ctx.Kit
	// determine whether the templateID is legal
	if err := s.validateTemplateID(kit, opt.ID); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ids := util.IntArrayUnique(opt.ObjectIDs)
	objIDs, err := s.getObjectByIDs(kit, ids)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := canObjBindingFieldTemplate(kit, objIDs); err != nil {
		blog.Errorf("validate objID failed, ids: %+v, err: %v, rid: %s", ids, err, kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	rows := make([]metadata.ObjFieldTemplateRelation, 0)
	for _, id := range ids {
		rows = append(rows, metadata.ObjFieldTemplateRelation{
			ObjectID:   id,
			TemplateID: opt.ID,
			TenantID:   kit.TenantID,
		})
	}

	if err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameObjFieldTemplateRelation).
		Insert(kit.Ctx, rows); err != nil {
		blog.Errorf("create failed, db insert failed, doc: %+v, err: %+v, rid: %s", rows, err, kit.Rid)
		ctx.RespAutoError(kit.CCError.CCError(common.CCErrCommDBInsertFailed))
		return
	}
	ctx.RespEntity(nil)
}

func (s *service) dealProcessRunningTasks(kit *rest.Kit, option *metadata.FieldTemplateUnbindObjOpt) error {

	// 1、get the status of the task
	cond := mapstr.MapStr{
		common.BKInstIDField:       option.ID,
		metadata.APITaskExtraField: option.ObjectID,
		common.BKStatusField: mapstr.MapStr{
			common.BKDBIN: []metadata.APITaskStatus{
				metadata.APITaskStatusExecute, metadata.APITaskStatusWaitExecute,
				metadata.APITaskStatusNew},
		},
	}

	result := make([]metadata.APITaskSyncStatus, 0)
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameAPITaskSyncHistory).Find(cond).
		Fields(common.BKStatusField, common.BKCloudSyncTaskID).
		All(kit.Ctx, &result); err != nil {
		blog.Errorf("search mainline failed cond: %+v, err: %s, rid: %s", cond, err, kit.Rid)
		return err
	}

	if len(result) == 0 {
		return nil
	}

	// 2、the possible task status scenarios are: one is executing,
	// one is waiting or new, but there will be no more than two tasks.
	if len(result) > metadata.APITaskFieldTemplateMaxNum {
		blog.Errorf("task num incorrect, template ID: %d, objID: %s, rid: %s", option.ID, option.ObjectID, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommGetMultipleObject,
			fmt.Sprintf("template ID: %d, objID: %s", option.ID, option.ObjectID))
	}

	// 3、if there is a running task, return an error directly.
	var taskID string
	for _, info := range result {
		if info.Status == metadata.APITaskStatusExecute {
			blog.Errorf("unbinding failed, sync task(%s) is running, template ID: %d, objID: %s, rid: %d")
			return kit.CCError.Errorf(common.CCErrTaskCreateConflict,
				fmt.Sprintf("template ID: %d, objID: %s", option.ID, option.ObjectID))
		}
		taskID = info.TaskID
	}

	// 4、if there is a queued task that needs to be cleared.
	delCond := mapstr.MapStr{
		common.BKTaskIDField: taskID,
		common.BKStatusField: mapstr.MapStr{
			common.BKDBIN: []metadata.APITaskStatus{metadata.APITaskStatusWaitExecute, metadata.APITaskStatusNew},
		},
	}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameAPITask).Delete(kit.Ctx, delCond)
	if err != nil {
		blog.Errorf("delete task failed, cond: %#v, err: %v, rid: %s", delCond, err, kit.Rid)
		return err
	}
	return nil
}

// FieldTemplateUnbindObject field template unbind model.
func (s *service) FieldTemplateUnbindObject(ctx *rest.Contexts) {

	opt := new(metadata.FieldTemplateUnbindObjOpt)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	kit := ctx.Kit
	// 1、judging the legitimacy of parameters
	if err := s.validateTemplateID(kit, opt.ID); err != nil {
		ctx.RespAutoError(err)
		return
	}

	objIDs, err := s.getObjectByIDs(kit, []int64{opt.ObjectID})
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := canObjBindingFieldTemplate(kit, objIDs); err != nil {
		blog.Errorf("validate failed, objIDs: %+v, err: %v, rid: %s", objIDs, err, kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// 2、process tasks in task
	if err := dealProcessRunningTasks(kit, []int64{opt.ID}, opt.ObjectID); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// 3、delete binding relationship
	if err := s.deleteFieldTmplRelation(kit, opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// 4、set the templateID in the involved model attribute field and model unique check field to 0
	if err := s.fieldTemplateUnbindAttrAndUnique(ctx.Kit, opt.ID, objIDs); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func dealProcessRunningTasks(kit *rest.Kit, ids []int64, objectID int64) error {

	// 1、get the status of the task
	cond := mapstr.MapStr{
		common.BKInstIDField: mapstr.MapStr{
			common.BKDBIN: ids,
		},
		metadata.APITaskExtraField: objectID,
		common.BKStatusField: mapstr.MapStr{
			common.BKDBIN: []metadata.APITaskStatus{
				metadata.APITaskStatusExecute, metadata.APITaskStatusWaitExecute,
				metadata.APITaskStatusNew},
		},
	}

	result := make([]metadata.APITaskSyncStatus, 0)
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameAPITaskSyncHistory).Find(cond).
		Fields(common.BKStatusField, common.BKTaskIDField).
		All(kit.Ctx, &result); err != nil {
		blog.Errorf("search task failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}

	if len(result) == 0 {
		return nil
	}

	// 2、the possible task status scenarios are: one is executing,
	// one is waiting or new, but there will be no more than two tasks.
	if len(result) > metadata.MaxFieldTemplateTaskNum {
		blog.Errorf("task num incorrect, template IDs: %v, objID: %d, rid: %s", ids, objectID, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommGetMultipleObject,
			fmt.Sprintf("template IDs: %v, objID: %d", ids, objectID))
	}

	// 3、if there is a running task, return an error directly.
	var taskID string
	for _, info := range result {
		if info.Status == metadata.APITaskStatusExecute {
			blog.Errorf("unbinding failed, sync task(%s) is running, template ID: %v, objID: %d, rid: %d", info.TaskID,
				ids, objectID, kit.Rid)
			return kit.CCError.Errorf(common.CCErrTaskDeleteConflict,
				fmt.Sprintf("template IDs: %v, objID: %d", ids, objectID))
		}
		taskID = info.TaskID
	}

	// 4、if there is a queued task that needs to be cleared.
	delCond := mapstr.MapStr{
		common.BKTaskIDField: taskID,
		common.BKStatusField: mapstr.MapStr{
			common.BKDBIN: []metadata.APITaskStatus{metadata.APITaskStatusWaitExecute, metadata.APITaskStatusNew},
		},
	}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameAPITask).Delete(kit.Ctx, delCond)
	if err != nil {
		blog.Errorf("delete task failed, cond: %#v, err: %v, rid: %s", delCond, err, kit.Rid)
		return err
	}
	return nil
}

func (s *service) deleteFieldTmplRelation(kit *rest.Kit, option *metadata.FieldTemplateUnbindObjOpt) error {

	cond := mapstr.MapStr{
		common.BKTemplateID:  option.ID,
		common.ObjectIDField: option.ObjectID,
	}

	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjFieldTemplateRelation).Delete(kit.Ctx,
		cond); err != nil {
		blog.Errorf("delete obj field template failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}
	return nil
}

func (s *service) fieldTemplateUnbindAttrAndUnique(kit *rest.Kit, id int64, objIDs []string) error {
	if len(objIDs) == 0 {
		return nil
	}

	tmplCond := mapstr.MapStr{
		common.BKTemplateID: id,
	}

	dbTmplAttrs := make([]metadata.FieldTemplateAttr, 0)
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAttDesTemplate).Find(tmplCond).
		Fields(common.BKFieldID).All(kit.Ctx, &dbTmplAttrs); err != nil {
		blog.Errorf("list field template attrs failed, filter: %+v, err: %v, rid: %v", tmplCond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	// there must be at least one attribute on the field template
	if len(dbTmplAttrs) == 0 {
		blog.Errorf("no attribute founded, id: %d, rid: %s", id, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKTemplateID)
	}

	attrTmplIDs := make([]int64, len(dbTmplAttrs))
	for idx, attr := range dbTmplAttrs {
		attrTmplIDs[idx] = attr.ID
	}

	updateCond := mapstr.MapStr{
		common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objIDs},
		common.BKTemplateID: mapstr.MapStr{common.BKDBIN: attrTmplIDs},
	}
	data := mapstr.MapStr{common.BKTemplateID: 0}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAttDes).Update(kit.Ctx, updateCond,
		data); err != nil {
		blog.Errorf("update object attrs failed, filter: %+v, err: %v, rid: %v", updateCond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}

	dbTmplUniques := make([]metadata.FieldTemplateUnique, 0)
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjectUniqueTemplate).Find(tmplCond).Fields(
		common.BKFieldID).All(kit.Ctx, &dbTmplUniques); err != nil {
		blog.Errorf("list field template uniques failed, filter: %+v, err: %v, rid: %v", tmplCond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(dbTmplUniques) == 0 {
		return nil
	}

	uniqueTmplIDs := make([]int64, len(dbTmplUniques))
	for idx, unique := range dbTmplUniques {
		uniqueTmplIDs[idx] = unique.ID
	}

	updateCond[common.BKTemplateID] = mapstr.MapStr{common.BKDBIN: uniqueTmplIDs}

	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjUnique).Update(kit.Ctx, updateCond,
		data); err != nil {
		blog.Errorf("update object uniques failed, filter: %+v, err: %v, rid: %v", updateCond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}

	return nil
}

// CreateFieldTemplate create field template.
func (s *service) CreateFieldTemplate(ctx *rest.Contexts) {
	template := new(metadata.FieldTemplate)
	if err := ctx.DecodeInto(template); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := template.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	filter := make(map[string]interface{})
	count, err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameFieldTemplate).Find(filter).
		Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("count field templates failed, err: %v, filter: %+v, rid: %s", err, filter, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	if count >= metadata.FieldTmplNumLimit {
		blog.Errorf("field template exceeds the maximum number limit, count: %d, limit: %d, rid: %s", count,
			metadata.FieldTmplNumLimit, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "field template", count))
		return
	}

	id, err := mongodb.Shard(ctx.Kit.SysShardOpts()).NextSequence(ctx.Kit.Ctx, common.BKTableNameFieldTemplate)
	if err != nil {
		blog.Errorf("get sequence id on the table (%s) failed, err: %v, rid: %s", common.BKTableNameFieldTemplate, err,
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error()))
		return
	}

	template.ID = int64(id)
	template.TenantID = ctx.Kit.TenantID
	template.Creator = ctx.Kit.User
	template.Modifier = ctx.Kit.User
	now := time.Now()
	template.CreateTime = &metadata.Time{Time: now}
	template.LastTime = &metadata.Time{Time: now}

	if err = mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameFieldTemplate).Insert(ctx.Kit.Ctx,
		template); err != nil {
		blog.Errorf("save field template failed, data: %v, err: %v, rid: %s", template, err, ctx.Kit.Rid)
		if mongodb.IsDuplicatedError(err) {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(err)))
			return
		}
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBInsertFailed))
		return
	}

	ctx.RespEntity(metadata.RspID{ID: int64(id)})
}

// DeleteFieldTemplate delete field template
func (s *service) DeleteFieldTemplate(ctx *rest.Contexts) {
	opt := new(metadata.DeleteOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}
	tmplCond := opt.Condition

	templates := make([]metadata.FieldTemplate, 0)
	err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameFieldTemplate).Find(tmplCond).Fields(
		common.BKFieldID).
		All(ctx.Kit.Ctx, &templates)
	if err != nil {
		blog.Errorf("find field template failed, cond: %v, err: %v, rid: %s", tmplCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(templates) == 0 {
		ctx.RespEntity(nil)
		return
	}

	tmplIDs := make([]int64, 0)
	for _, template := range templates {
		tmplIDs = append(tmplIDs, template.ID)
	}
	countCond := mapstr.MapStr{common.BKTemplateID: mapstr.MapStr{common.BKDBIN: tmplIDs}}

	relationCount, err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameObjFieldTemplateRelation).Find(
		countCond).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("count field template relation failed, filter: %+v, err: %v, rid: %v", countCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	if relationCount != 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCoreServiceFieldTemplateHasRelation))
		return
	}

	attrCount, err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameObjAttDesTemplate).Find(countCond).
		Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("count field template attribute failed, filter: %+v, err: %v, rid: %v", countCond, err,
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	if attrCount != 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCoreServiceFieldTemplateHasAttr))
		return
	}

	uniqueCount, err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameObjectUniqueTemplate).
		Find(countCond).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("count field template unique failed, filter: %+v, err: %v, rid: %v", countCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	if uniqueCount != 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCoreServiceFieldTemplateHasUnique))
		return
	}

	if err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameFieldTemplate).Delete(ctx.Kit.Ctx,
		tmplCond); err != nil {
		blog.Errorf("delete field template failed, cond: %v, err: %v, rid: %s", tmplCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// UpdateFieldTemplate update field template
func (s *service) UpdateFieldTemplate(ctx *rest.Contexts) {
	opt := new(metadata.FieldTemplate)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	if opt.ID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKFieldID))
		return
	}

	cond := map[string]interface{}{common.BKFieldID: opt.ID}
	dbTmpl := new(metadata.FieldTemplate)
	err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameFieldTemplate).Find(cond).
		One(ctx.Kit.Ctx, dbTmpl)
	if err != nil {
		blog.Errorf("find template failed, cond: %v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if dbTmpl == nil {
		ctx.RespEntity(nil)
		return
	}

	opt.TenantID = dbTmpl.TenantID
	opt.Creator = dbTmpl.Creator
	opt.CreateTime = dbTmpl.CreateTime
	opt.Modifier = ctx.Kit.User
	now := time.Now()
	opt.LastTime = &metadata.Time{Time: now}

	if err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameFieldTemplate).Update(ctx.Kit.Ctx, cond,
		opt); err != nil {
		blog.Errorf("update field template failed, cond: %v, data: %v, err: %v, rid: %s", cond, opt, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// FindFieldTmplSimplifyByAttr according to the template ID of the model attribute,
// obtain the brief information of the corresponding field template.
func (s *service) FindFieldTmplSimplifyByAttr(ctx *rest.Contexts) {

	opt := new(metadata.ListTmplSimpleByAttrOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// 1、whether data can be found according to attrID and model templateID in parameters
	countCond := mapstr.MapStr{
		common.BKFieldID:    opt.AttrID,
		common.BKTemplateID: opt.TemplateID,
	}

	count, err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameObjAttDes).Find(countCond).
		Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("count object attr failed, filter: %+v, err: %v, rid: %v", countCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	if count != 1 {
		blog.Errorf("count object attr num error, count: %d, cond: %v, rid: %s", count, countCond, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid,
			fmt.Sprintf("%s,%s", common.BKAttributeIDField, common.BKTemplateID)))
		return
	}

	// 2、according to the templateID on the attribute, go to the template attribute to find the
	// corresponding field template templateID
	attrCond := mapstr.MapStr{
		common.BKFieldID: opt.TemplateID,
	}

	dbTmplAttrs := make([]metadata.FieldTemplateAttr, 0)
	if err = mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameObjAttDesTemplate).Find(attrCond).Fields(
		common.BKTemplateID).
		All(ctx.Kit.Ctx, &dbTmplAttrs); err != nil {
		blog.Errorf("find field template attr failed, cond: %v, err: %v, rid: %s", attrCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(dbTmplAttrs) != 1 {
		blog.Errorf("count field template attr num error, count: %d, cond: %v, rid: %s", len(dbTmplAttrs), attrCond,
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAttributeIDField))
		return
	}

	// 3、find brief information based on the field template templateID
	tmplCond := mapstr.MapStr{
		common.BKFieldID: dbTmplAttrs[0].TemplateID,
	}

	templates := make([]metadata.FieldTemplate, 0)
	if err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameFieldTemplate).Find(tmplCond).
		Fields(common.BKFieldID, common.BKFieldName).All(ctx.Kit.Ctx, &templates); err != nil {
		blog.Errorf("find field template failed, cond: %v, err: %v, rid: %s", tmplCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(templates) != 1 {
		blog.Errorf("count template num error, count: %d, cond: %v, rid: %s", len(templates), tmplCond, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAttributeIDField))
		return
	}

	result := metadata.ListTmplSimpleResult{
		Name: templates[0].Name,
		ID:   templates[0].ID,
	}
	ctx.RespEntity(result)
}

// FindFieldTmplSimplifyByUnique get the brief information of the corresponding
// field template according to the uniquely verified template ID of the model.
func (s *service) FindFieldTmplSimplifyByUnique(ctx *rest.Contexts) {

	opt := new(metadata.ListTmplSimpleByUniqueOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// 1、whether the data can be found according to unique and the model templateID in the parameter
	countCond := mapstr.MapStr{
		common.BKFieldID:    opt.UniqueID,
		common.BKTemplateID: opt.TemplateID,
	}

	count, err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameObjUnique).Find(countCond).Count(
		ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("count object unique failed, filter: %+v, err: %v, rid: %v", countCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	if count != 1 {
		blog.Errorf("count object unique num error, count: %d, cond: %v, rid: %s", count, countCond, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid,
			fmt.Sprintf("%s,%s", "bk_unique_id", common.BKTemplateID)))
		return
	}

	// 2、find the corresponding field template templateID on the template
	// attribute according to the templateID on the unique check
	attrCond := mapstr.MapStr{
		common.BKFieldID: opt.TemplateID,
	}

	dbTmplUniques := make([]metadata.FieldTemplateUnique, 0)
	if err = mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameObjectUniqueTemplate).Find(attrCond).
		Fields(common.BKTemplateID).All(ctx.Kit.Ctx, &dbTmplUniques); err != nil {
		blog.Errorf("find field template unique failed, cond: %v, err: %v, rid: %s", attrCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(dbTmplUniques) != 1 {
		blog.Errorf("count template unique num error, count: %d, cond: %v, rid: %s", len(dbTmplUniques), attrCond,
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAttributeIDField))
		return
	}

	// 3、find brief information based on the field template templateID.
	tmplCond := mapstr.MapStr{
		common.BKFieldID: dbTmplUniques[0].TemplateID,
	}

	templates := make([]metadata.FieldTemplate, 0)
	if err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameFieldTemplate).Find(tmplCond).
		Fields(common.BKFieldID, common.BKFieldName).All(ctx.Kit.Ctx, &templates); err != nil {
		blog.Errorf("find field template failed, cond: %v, err: %v, rid: %s", tmplCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(templates) != 1 {
		blog.Errorf("count template num error, count: %d, cond: %v, rid: %s", len(templates), tmplCond, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAttributeIDField))
		return
	}

	result := metadata.ListTmplSimpleResult{
		Name: templates[0].Name,
		ID:   templates[0].ID,
	}
	ctx.RespEntity(result)

}
