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
	"configcenter/src/ac/meta"
	"strconv"

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// ListFieldTemplate list field templates.
// NOTICE: this api only returns basic info of field template, do not need to authorize
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

	// list field templates
	res, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplate(cts.Kit.Ctx, cts.Kit.Header, opt)
	if err != nil {
		blog.Errorf("list field templates failed, err: %v, req: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	cts.RespEntity(res)
}

// FindFieldTemplateByID find the field combination template involved by ID.
func (s *service) FindFieldTemplateByID(ctx *rest.Contexts) {

	idStr := ctx.Request.PathParameter(common.BKFieldID)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the path params id(%s), err: %v, rid: %s", idStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKFieldID))
		return
	}

	if id == 0 {
		blog.Errorf("template id cannot be zero, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKFieldID))
		return
	}

	query := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{
			Filter: filtertools.GenAtomFilter(common.BKFieldID, filter.Equal, id),
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}

	// list field template by id
	res, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplate(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("get field template failed, req: %+v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(res.Info) > 1 {
		blog.Errorf("multiple field templates found, req: %+v, rid: %s", query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommGetMultipleObject))
		return
	}

	if len(res.Info) == 0 {
		ctx.RespEntity(nil)
		return
	}
	ctx.RespEntity(res.Info[0])
}

// canObjsBindFieldTemplate currently only supports "host" in the mainline model
func (s *service) canObjsBindFieldTemplate(kit *rest.Kit, ids []int64) ccErr.CCErrorCoder {

	query := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKFieldID: mapstr.MapStr{
				common.BKDBIN: ids,
			},
		},
		Fields:         []string{common.BKObjIDField},
		DisableCounter: true,
	}

	obj, err := s.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("failed to find objects by query(%#v), err: %v, rid: %s", query, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommNotFound, err.Error())
	}

	if len(obj.Info) == 0 {
		blog.Errorf("failed to find objIDs by query(%#v),  rid: %s", query, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommNotFound)
	}
	if len(obj.Info) != len(ids) {
		blog.Errorf("obj ids are invalid, input: %+v, rid: %s", query, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "object_ids")
	}

	objIDs := make([]string, 0)
	for _, info := range obj.Info {
		objIDs = append(objIDs, info.ObjectID)
	}

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

	counts, ccErr := s.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameObjAsst, []map[string]interface{}{cond})
	if ccErr != nil {
		blog.Errorf("get mainline obj count failed, cond: %+v, err: %v, rid: %s", cond, ccErr, kit.Rid)
		return ccErr
	}

	if len(counts) != 1 || int(counts[0]) > 0 {
		blog.Errorf("obj ids cannot be a non-host mainline model, ids: %+v, rid: %s", ids, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "object_ids")
	}

	return nil
}

// FieldTemplateBindObject field template binding model
func (s *service) FieldTemplateBindObject(ctx *rest.Contexts) {

	opt := new(metadata.FieldTemplateBindObjOpt)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	objIDs := make([]int64, 0)
	objIDs = util.IntArrayUnique(opt.ObjectIDs)
	if err := s.canObjsBindFieldTemplate(ctx.Kit, objIDs); err != nil {
		ctx.RespAutoError(err)
		return
	}

	option := &metadata.FieldTemplateBindObjOpt{
		ID:        opt.ID,
		ObjectIDs: objIDs,
	}

	// todo:待补充鉴权日志
	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.clientSet.CoreService().FieldTemplate().FieldTemplateBindObject(ctx.Kit.Ctx, ctx.Kit.Header, option)
		if err != nil {
			blog.Errorf("field template bind model failed, err: %v , rid: %s", err, ctx.Kit.Rid)
			return err
		}
		// todo:开启事务记录审计日志
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// FieldTemplateUnbindObject field template binding model
func (s *service) FieldTemplateUnbindObject(ctx *rest.Contexts) {
	opt := new(metadata.FieldTemplateUnbindObjOpt)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	if err := s.canObjsBindFieldTemplate(ctx.Kit, []int64{opt.ObjectID}); err != nil {
		ctx.RespAutoError(err)
		return
	}
	// todo:待补充鉴权日志

	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.clientSet.CoreService().FieldTemplate().FieldTemplateUnbindObject(ctx.Kit.Ctx, ctx.Kit.Header, opt)
		if err != nil {
			blog.Errorf("field template unbind model failed, err: %v , rid: %s", err, ctx.Kit.Rid)
			return err
		}
		// todo:开启事务记录审计日志
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// CreateFieldTemplate create field template(contains field template brief information, attributes and uniques)
func (s *service) CreateFieldTemplate(ctx *rest.Contexts) {
	opt := new(metadata.CreateFieldTmplOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := opt.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	if authResp, authorized := s.auth.Authorize(ctx.Kit, meta.ResourceAttribute{Basic: meta.Basic{
		Type: meta.FieldTemplate, Action: meta.Create}}); !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	var res *metadata.RspID
	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		res, err = s.logics.FieldTemplateOperation().CreateFieldTemplate(ctx.Kit, opt)
		if err != nil {
			blog.Errorf("create field template failed, opt: %v, err: %v, rid: %s", opt, err, ctx.Kit.Rid)
			return err
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(res)
}
