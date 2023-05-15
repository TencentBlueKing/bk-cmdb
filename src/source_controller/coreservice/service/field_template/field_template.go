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

func getFieldTemplatesUnSupportedModels(kit *rest.Kit) (map[string]struct{}, error) {

	// 1、query the mainline model
	result := make([]metadata.Association, 0)
	cond := mapstr.MapStr{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}
	cond = util.SetQueryOwner(cond, kit.SupplierAccount)
	if err := mongodb.Client().Table(common.BKTableNameObjAsst).Find(cond).Fields(common.BKObjIDField).
		All(kit.Ctx, &result); err != nil {
		blog.Errorf("search mainline failed cond: %+v, err: %s, rid: %s", cond, err, kit.Rid)
		return nil, err
	}
	// 2、get a list of models that do not support binding field composition templates
	objID := make(map[string]struct{})
	for _, data := range result {
		if data.ObjectID == common.BKInnerObjIDHost {
			continue
		}
		objID[data.ObjectID] = struct{}{}
	}
	return objID, nil
}

func (s *service) validateTemplateID(kit *rest.Kit, id int64) error {
	// 判断templateID 是否合法
	cond := mapstr.MapStr{common.BKFieldID: id}
	cnt, err := mongodb.Client().Table(common.BKTableNameFieldTemplate).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count field template failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if cnt == 0 {
		blog.Errorf("no field template founded, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommNotFound)
	}
	if cnt > 1 {
		blog.Errorf("multi field template founded, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}
	return nil
}

// FieldTemplateBindObject field template bind model.
func (s *service) FieldTemplateBindObject(ctx *rest.Contexts) {

	opt := new(metadata.FieldTemplateBindObjOpt)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}
	kit := ctx.Kit

	// 判断templateID 是否合法
	if err := s.validateTemplateID(kit, opt.ID); err != nil {
		ctx.RespAutoError(err)
		return
	}

	objIDMap, err := getFieldTemplatesUnSupportedModels(ctx.Kit)
	if err != nil {
		blog.Errorf("multi field template founded, cond: %+v, rid: %s", err, kit.Rid)
		ctx.RespAutoError(kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	rows := make([]metadata.ObjFieldTemplateRelation, 0)
	for _, objID := range opt.ObjectIDs {
		// 判断objectID是否合法
		if _, ok := objIDMap[objID]; ok || objID == "" {
			blog.Errorf("object(%s) is not allowed to bind field template, rid: %s", objID, kit.Rid)
			ctx.RespAutoError(kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
			return
		}
		rows = append(rows, metadata.ObjFieldTemplateRelation{
			ObjectID:   objID,
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

func (s *service) dealProcessRunningTasks() error {
	// 1、获取任务
	// 2、如果有运行中的任务直接返回报错
	// 3、如果有排队中的任务清楚掉任务

	return nil
}

// FieldTemplateUnBindObject field template unbind model.
func (s *service) FieldTemplateUnBindObject(ctx *rest.Contexts) {

	opt := new(metadata.FieldTemplateUnBindObjOpt)
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

	objIDMap, err := getFieldTemplatesUnSupportedModels(ctx.Kit)
	if err != nil {
		blog.Errorf("multi field template founded, cond: %+v, rid: %s", err, kit.Rid)
		ctx.RespAutoError(kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	if _, ok := objIDMap[opt.ObjectID]; ok || opt.ObjectID == "" {
		blog.Errorf("object(%s) is not allowed to bind field template, rid: %s", opt.ObjectID, kit.Rid)
		ctx.RespAutoError(kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
		return
	}

	// 2、process tasks in task
	if err := s.dealProcessRunningTasks(); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// 3、delete binding relationship
	cond := mapstr.MapStr{
		common.BKTemplateID: opt.ID,
		common.BKObjIDField: opt.ObjectID,
	}
	cond = util.SetModOwner(cond, kit.SupplierAccount)
	if err := mongodb.Client().Table(common.BKTableNameObjFieldTemplateRelation).Delete(kit.Ctx, cond); err != nil {
		blog.Errorf("delete obj field template failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// FindFieldTemplateTasksStatus 查找指定字段模版下的各个模型同步任务状态.
func (s *service) FindFieldTemplateTasksStatus(ctx *rest.Contexts) {

}
