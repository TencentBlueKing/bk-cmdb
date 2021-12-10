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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// BatchCreateSet batch create set
func (s *Service) BatchCreateSet(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if err != nil {
		blog.Errorf("batch create set failed, parse biz id from url failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	batchBody := metadata.BatchCreateSetRequest{}
	if err := ctx.DecodeInto(&batchBody); err != nil {
		ctx.RespAutoError(err)
		return
	}

	batchCreateResult := make([]metadata.OneSetCreateResult, 0)
	var firstErr error
	for idx, set := range batchBody.Sets {
		if _, ok := set[common.BkSupplierAccount]; !ok {
			set[common.BkSupplierAccount] = ctx.Kit.SupplierAccount
		}
		set[common.BKAppIDField] = bizID

		result := make(mapstr.MapStr)
		// to avoid judging to be nested transaction, need a new header
		ctx.Kit.Header = ctx.Kit.NewHeader()
		txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
			var err error
			result, err = s.Logics.SetOperation().CreateSet(ctx.Kit, bizID, set)
			if err != nil && firstErr == nil {
				firstErr = err
			}
			if err != nil && blog.V(3) {
				blog.Errorf("batch create set failed, idx: %d, data: %#v, err: %v, rid: %s", idx, set, err,
					ctx.Kit.Rid)
			}
			return err
		})

		errMsg := ""
		if txnErr != nil {
			errMsg = txnErr.Error()
		}
		batchCreateResult = append(batchCreateResult, metadata.OneSetCreateResult{
			Index:    idx,
			Data:     result,
			ErrorMsg: errMsg,
		})
	}

	ctx.RespEntityWithError(batchCreateResult, firstErr)
}

// CreateSet create a new set
func (s *Service) CreateSet(ctx *rest.Contexts) {
	data := mapstr.MapStr{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id from url, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	resp := make(mapstr.MapStr)
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		resp, err = s.Logics.SetOperation().CreateSet(ctx.Kit, bizID, data)
		if err != nil {
			blog.Errorf("create set failed, bizID: %d, data: %#v, err: %v, rid: %s", bizID, data, err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(resp)
}

// checkIsBuiltInSet check if set is built-in set
func (s *Service) checkIsBuiltInSet(kit *rest.Kit, setIDs ...int64) error {
	// 检查是否是内置集群
	cond := &metadata.Condition{
		Condition: map[string]interface{}{
			common.BKSetIDField: map[string]interface{}{
				common.BKDBIN: setIDs,
			},
			// 当default值不等于0时为内置集群
			common.BKDefaultField: map[string]interface{}{
				common.BKDBNE: common.DefaultFlagDefaultValue,
			},
		},
	}

	rsp, e := s.Engine.CoreAPI.CoreService().Instance().CountInstances(kit.Ctx, kit.Header, common.BKInnerObjIDSet, cond)
	if e != nil {
		blog.Errorf("check is built in set failed, option: %s, err: %v, rid: %s", cond, e, kit.Rid)
		return e
	}

	if rsp.Count > 0 {
		return kit.CCError.CCError(common.CCErrorTopoForbiddenDeleteOrUpdateBuiltInSetModule)
	}

	return nil
}

// DeleteSets batch delete the set
func (s *Service) DeleteSets(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id from url, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	data := new(metadata.OpCondition)
	if err = ctx.DecodeInto(data); err != nil {
		blog.Errorf("failed to parse to the operation condition, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	setIDs := data.Delete.InstID
	// 检查是否是内置集群
	if err := s.checkIsBuiltInSet(ctx.Kit, setIDs...); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Logics.SetOperation().DeleteSet(ctx.Kit, bizID, setIDs)
		if err != nil {
			blog.Errorf("delete set failed, intIDs: %d, err: %v, rid: %s", setIDs, err, ctx.Kit.Rid)
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

// DeleteSet delete the set
func (s *Service) DeleteSet(ctx *rest.Contexts) {
	if "batch" == ctx.Request.PathParameter("set_id") {
		s.DeleteSets(ctx)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id from url, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	setID, err := strconv.ParseInt(ctx.Request.PathParameter("set_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the set id from url, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "set id"))
		return
	}

	// 检查是否时内置集群
	if err := s.checkIsBuiltInSet(ctx.Kit, setID); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Logics.SetOperation().DeleteSet(ctx.Kit, bizID, []int64{setID})
		if err != nil {
			blog.Errorf("delete set failed, bizID: %d, setID: %d, err: %v rid: %s", bizID, setID, err, ctx.Kit.Rid)
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

// UpdateSet update the set
func (s *Service) UpdateSet(ctx *rest.Contexts) {
	data := make(mapstr.MapStr)
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id from url, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	setID, err := strconv.ParseInt(ctx.Request.PathParameter("set_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the set id from url, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "set id"))
		return
	}

	// 不允许修改内置集群
	if err := s.checkIsBuiltInSet(ctx.Kit, setID); err != nil {
		blog.Errorf("check is builtIn set failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Logics.SetOperation().UpdateSet(ctx.Kit, data, bizID, setID)
		if err != nil {
			blog.Errorf("update set failed, data: %#v, setID: %d, err: %v, rid: %s", data, setID, err, ctx.Kit.Rid)
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

// SearchSet search the set
func (s *Service) SearchSet(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id from url, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	queryCond := new(metadata.QueryCondition)
	if err = ctx.DecodeInto(queryCond); err != nil {
		blog.Errorf("search set failed, decode parameter condition failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if queryCond.Condition == nil {
		queryCond.Condition = mapstr.New()
	}

	queryCond.Condition[common.BKAppIDField] = bizID

	instItems, err := s.Logics.InstOperation().FindInst(ctx.Kit, common.BKInnerObjIDSet, queryCond)
	if err != nil {
		blog.Errorf("failed to find inst, err: %v, rid: %s", ctx.Request.PathParameter("obj_id"), err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(instItems)

	return
}

// SearchSetBatch search the sets in one biz
func (s *Service) SearchSetBatch(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	option := metadata.SearchInstBatchOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := option.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	setIDs := util.IntArrayUnique(option.IDs)
	cond := mapstr.MapStr{
		common.BKAppIDField: bizID,
		common.BKSetIDField: mapstr.MapStr{
			common.BKDBIN: setIDs,
		},
	}

	qc := &metadata.QueryCondition{
		Fields: option.Fields,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: cond,
	}
	instanceResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDSet, qc)
	if err != nil {
		blog.Errorf("search module batch failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(instanceResult.Info)
}
