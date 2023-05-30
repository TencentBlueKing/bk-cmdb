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
	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// ListObjFieldTmplRel list field template and object relations.
func (s *service) ListObjFieldTmplRel(cts *rest.Contexts) {
	opt := new(metadata.ListObjFieldTmplRelOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// TODO add find object or template auth check after find object operation authorization is supported

	// list field templates and object relations
	var relFilter *filter.Expression
	if len(opt.TemplateIDs) > 0 {
		relFilter = filtertools.GenAtomFilter(common.BKTemplateID, filter.In, opt.TemplateIDs)
	}

	if len(opt.ObjectIDs) > 0 {
		var err error
		relFilter, err = filtertools.And(relFilter, filtertools.GenAtomFilter(common.ObjectIDField, filter.In,
			opt.ObjectIDs))
		if err != nil {
			cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "filter"))
			return
		}
	}

	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: relFilter},
		Page:               metadata.BasePage{Limit: common.BKNoLimit},
	}
	res, err := s.clientSet.CoreService().FieldTemplate().ListObjFieldTmplRel(cts.Kit.Ctx, cts.Kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field templates failed, err: %v, opt: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	cts.RespEntity(res)
}

// ListFieldTmplByObj list field template by related object id.
func (s *service) ListFieldTmplByObj(cts *rest.Contexts) {
	opt := new(metadata.ListFieldTmplByObjOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// TODO add find object or template auth check after find object operation authorization is supported

	// get field templates ids by object id
	relOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: filtertools.GenAtomFilter(common.ObjectIDField,
			filter.Equal, opt.ObjectID)},
		Page: metadata.BasePage{Limit: common.BKNoLimit},
	}
	relRes, err := s.clientSet.CoreService().FieldTemplate().ListObjFieldTmplRel(cts.Kit.Ctx, cts.Kit.Header, relOpt)
	if err != nil {
		blog.Errorf("list field templates failed, err: %v, opt: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	templateIDs := make([]int64, len(relRes.Info))
	for i, relation := range relRes.Info {
		templateIDs[i] = relation.TemplateID
	}

	if len(templateIDs) == 0 {
		cts.RespEntity(&metadata.FieldTemplateInfo{Info: make([]metadata.FieldTemplate, 0)})
		return
	}

	// list filed template by ids
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: filtertools.GenAtomFilter(common.BKFieldID,
			filter.In, templateIDs)},
		Page: metadata.BasePage{Limit: common.BKNoLimit},
	}
	res, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplate(cts.Kit.Ctx, cts.Kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field templates failed, err: %v, opt: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	cts.RespEntity(res)
}

// ListObjByFieldTmpl list object by field template.
func (s *service) ListObjByFieldTmpl(cts *rest.Contexts) {
	opt := new(metadata.ListObjByFieldTmplOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// TODO add find object or template auth check after find object operation authorization is supported
	// if object detail is needed later, add object auth check

	// get object ids by field template id
	relOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: filtertools.GenAtomFilter(common.BKTemplateID,
			filter.Equal, opt.TemplateID)},
		Page: metadata.BasePage{Limit: common.BKNoLimit},
	}
	relRes, err := s.clientSet.CoreService().FieldTemplate().ListObjFieldTmplRel(cts.Kit.Ctx, cts.Kit.Header, relOpt)
	if err != nil {
		blog.Errorf("list field templates failed, err: %v, opt: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	objectIDs := make([]int64, len(relRes.Info))
	for i, relation := range relRes.Info {
		objectIDs[i] = relation.ObjectID
	}

	if len(objectIDs) == 0 {
		cts.RespEntity(&metadata.FieldTemplateInfo{Info: make([]metadata.FieldTemplate, 0)})
		return
	}

	// list object by ids
	listOpt := &metadata.QueryCondition{
		Fields:    []string{common.BKFieldID, common.BKFieldName},
		Page:      opt.Page,
		Condition: mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: objectIDs}},
	}
	res, objErr := s.clientSet.CoreService().Model().ReadModel(cts.Kit.Ctx, cts.Kit.Header, listOpt)
	if objErr != nil {
		blog.Errorf("list objects failed, err: %v, opt: %+v, rid: %s", objErr, opt, cts.Kit.Rid)
		cts.RespAutoError(objErr)
		return
	}

	cts.RespEntity(res)
}

// SyncFieldTemplateToObjectTask synchronize field template information to model tasks.
func (s *service) SyncFieldTemplateToObjectTask(ctx *rest.Contexts) {
	syncOption := new(metadata.SyncObjectTask)
	if err := ctx.DecodeInto(syncOption); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if rawErr := syncOption.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}
	// templateID 是否存在。

	// list filed template by ids
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: filtertools.GenAtomFilter(common.BKFieldID,
			filter.Equal, syncOption.TemplateID)},
		Page: metadata.BasePage{Limit: common.BKNoLimit},
	}
	res, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplate(ctx.Kit.Ctx, ctx.Kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field templates failed, err: %v, opt: %+v, rid: %s", err, syncOption, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(res.Info) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKTemplateID))
		return
	}
	// object 是否合法 主要是看一下通过这个id和templateID 是否能查到，如果查不到直接报错就好了

	// todo: 如果能查到还得看一下这个模型是否停用状态如果停用？这个涉及到停用模型的时候是否解绑关系，后续再补充这块的逻辑

	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.doSyncFieldTemplateTask(ctx.Kit, syncOption); err != nil {
			blog.Errorf("do sync service instance task(%#v) failed, err: %v, rid: %s", syncOption, err, ctx.Kit.Rid)
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

func (s *service) doSyncFieldTemplateTask(kit *rest.Kit, option *metadata.SyncObjectTask) error {

	return nil
}
