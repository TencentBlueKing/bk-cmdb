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
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// DeleteBizSet delete business set
func (s *Service) DeleteBizSet(ctx *rest.Contexts) {
	opt := new(metadata.DeleteBizSetOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(opt.BizSetIDs) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "bk_biz_set_ids"))
		return
	}

	if len(opt.BizSetIDs) > 100 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "bk_biz_set_ids", 100))
		return
	}

	// delete bizSet instances and related resources
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logics.InstOperation().DeleteInstByInstID(ctx.Kit, common.BKInnerObjIDBizSet, opt.BizSetIDs, false)
		if err != nil {
			blog.Errorf("delete biz set failed, ids: %v, err: %v, rid: %s", opt.BizSetIDs, err, ctx.Kit.Rid)
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

// FindBizInBizSet find all biz id and name in biz set
func (s *Service) FindBizInBizSet(ctx *rest.Contexts) {
	opt := new(metadata.FindBizInBizSetOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if opt.BizSetID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKBizSetIDField))
		return
	}

	if rawErr := opt.Page.ValidateWithEnableCount(false, common.BKMaxInstanceLimit); rawErr.ErrCode != 0 {
		blog.Errorf("page is invalid, err: %v, option: %#v, rid: %s", rawErr, opt, ctx.Kit.Rid)
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// get biz mongo condition by biz scope in biz set
	bizSetBizCond, err := s.getBizSetBizCond(ctx.Kit, opt.BizSetID)
	if err != nil {
		blog.Errorf("get biz cond by biz set id %d failed, err: %v, rid: %s", opt.BizSetID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// count biz in biz set is enable count is set
	if opt.Page.EnableCount {
		filter := []map[string]interface{}{bizSetBizCond}

		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKTableNameBaseApp, filter)
		if err != nil {
			blog.Errorf("count biz failed, err: %v, cond: %#v, rid: %s", err, bizSetBizCond, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntity(mapstr.MapStr{"count": counts[0]})
		return
	}

	// get biz in biz set if enable count is set
	bizOpt := &metadata.QueryCondition{
		Condition:      bizSetBizCond,
		Fields:         opt.Fields,
		Page:           opt.Page,
		DisableCounter: true,
	}

	_, biz, err := s.Logics.BusinessOperation().FindBiz(ctx.Kit, bizOpt)
	if err != nil {
		blog.Errorf("find biz failed, err: %v, cond: %#v, rid: %s", err, bizOpt, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(mapstr.MapStr{"info": biz})
}

// getBizSetBizCond get biz mongo condition from the biz set scope
func (s *Service) getBizSetBizCond(kit *rest.Kit, bizSetID int64) (mapstr.MapStr, error) {
	bizSetCond := &metadata.QueryCondition{
		Fields:         []string{common.BKScopeField},
		Page:           metadata.BasePage{Limit: 1},
		Condition:      map[string]interface{}{common.BKBizSetIDField: bizSetID},
		DisableCounter: true,
	}

	bizSetRes := new(metadata.BizSetInstanceResponse)
	err := s.Engine.CoreAPI.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDBizSet,
		bizSetCond, &bizSetRes)
	if err != nil {
		blog.Errorf("get biz set failed, cond: %#v, err: %v, rid: %s", bizSetCond, err, kit.Rid)
		return nil, err
	}

	if err := bizSetRes.CCError(); err != nil {
		blog.Errorf("get biz set failed, cond: %#v, err: %v, rid: %s", bizSetCond, err, kit.Rid)
		return nil, err
	}

	if len(bizSetRes.Data.Info) == 0 {
		blog.Errorf("get no biz set by cond: %#v, rid: %s", bizSetCond, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKBizSetIDField)
	}

	bizSetBizCond, errKey, rawErr := bizSetRes.Data.Info[0].Scope.ToMgo()
	if rawErr != nil {
		blog.Errorf("parse biz set scope(%#v) failed, err: %v, rid: %s", bizSetRes.Data.Info[0].Scope, rawErr, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errKey)
	}

	// do not include resource pool biz in biz set by default
	if _, exists := bizSetBizCond[common.BKDefaultField]; !exists {
		bizSetBizCond[common.BKDefaultField] = mapstr.MapStr{common.BKDBNE: common.DefaultAppFlag}
	}

	return bizSetBizCond, nil
}
