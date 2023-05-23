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

	filter = util.SetQueryOwner(filter, cts.Kit.SupplierAccount)

	if opt.Page.EnableCount {
		count, err := mongodb.Client().Table(common.BKTableNameFieldTemplate).Find(filter).Count(cts.Kit.Ctx)
		if err != nil {
			blog.Errorf("count field templates failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
			cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}

		cts.RespEntity(metadata.FieldTemplateInfo{Count: count})
		return
	}

	fieldTemplates := make([]metadata.FieldTemplate, 0)
	err = mongodb.Client().Table(common.BKTableNameFieldTemplate).Find(filter).Start(uint64(opt.Page.Start)).
		Limit(uint64(opt.Page.Limit)).Fields(opt.Fields...).All(cts.Kit.Ctx, &fieldTemplates)
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
	cond = util.SetQueryOwner(cond, kit.SupplierAccount)

	count, err := mongodb.Client().Table(common.BKTableNameObjAsst).Find(cond).Count(kit.Ctx)
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
	cnt, err := mongodb.Client().Table(common.BKTableNameFieldTemplate).Find(cond).Count(kit.Ctx)
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

	if err := mongodb.Client().Table(common.BKTableNameObjDes).Find(filter).Fields(common.BKObjIDField).
		All(kit.Ctx, &objs); err != nil {
		blog.Errorf("mongodb count failed, table: %s, err: %v, rid: %s", common.BKTableNameObjDes, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
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
			OwnerID:    kit.SupplierAccount,
		})
	}

	if err := mongodb.Client().Table(common.BKTableNameObjFieldTemplateRelation).Insert(kit.Ctx, rows); err != nil {
		blog.Errorf("create  failed, db insert failed, doc: %+v, err: %+v, rid: %s", rows, err, kit.Rid)
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
	cond = util.SetQueryOwner(cond, kit.SupplierAccount)

	result := make([]metadata.APITaskSyncStatus, 0)
	if err := mongodb.Client().Table(common.BKTableNameAPITaskSyncHistory).Find(cond).
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
	if len(result) > 2 {
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
	err := mongodb.Client().Table(common.BKTableNameAPITask).Delete(kit.Ctx, delCond)
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
		blog.Errorf("multi field template founded, cond: %+v, rid: %s", err, kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// 2、process tasks in task
	if err := s.dealProcessRunningTasks(kit, opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// 3、delete binding relationship
	cond := mapstr.MapStr{
		common.BKTemplateID:  opt.ID,
		common.ObjectIDField: opt.ObjectID,
	}
	cond = util.SetModOwner(cond, kit.SupplierAccount)
	if err := mongodb.Client().Table(common.BKTableNameObjFieldTemplateRelation).Delete(kit.Ctx, cond); err != nil {
		blog.Errorf("delete obj field template failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}
