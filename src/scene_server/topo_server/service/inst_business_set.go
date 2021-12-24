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
	"configcenter/src/common/mapstr"
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

// FindBizSetTopo find topo nodes id and name info by parent node in biz set
func (s *Service) FindBizSetTopo(ctx *rest.Contexts) {
	opt := new(metadata.FindBizSetTopoOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	topo, err := s.findBizSetTopo(ctx.Kit, opt)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(topo)
}

func (s *Service) findBizSetTopo(kit *rest.Kit, opt *metadata.FindBizSetTopoOption) ([]mapstr.MapStr, error) {
	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		blog.Errorf("option(%#v) is invalid, err: %v, rid: %s", opt, rawErr, kit.Rid)
		return nil, rawErr.ToCCError(kit.CCError)
	}

	// get biz mongo condition by biz scope in biz set
	bizSetBizCond, err := s.getBizSetBizCond(kit, opt.BizSetID)
	if err != nil {
		blog.Errorf("get biz cond by biz set id %d failed, err: %v, rid: %s", opt.BizSetID, err, kit.Rid)
		return nil, err
	}

	// check if the parent node belongs to a biz that is in the biz set
	if err := s.checkTopoNodeInBizSet(kit, opt.ParentObjID, opt.ParentID, bizSetBizCond); err != nil {
		blog.Errorf("check if parent %s node %d in biz failed, err: %v, biz cond: %#v, rid: %s", opt.ParentObjID,
			opt.ParentID, err, bizSetBizCond, kit.Rid)
		return nil, err
	}

	// get parent object id to check if the parent node is a valid mainline instance that belongs to the biz set
	var childObjID string
	switch opt.ParentObjID {
	case common.BKInnerObjIDBizSet:
		if opt.ParentID != opt.BizSetID {
			blog.Errorf("biz parent id %s is not equal to biz set id %s, rid: %s", opt.ParentID, opt.BizSetID, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKParentIDField)
		}

		// find biz nodes by the condition in biz sets
		bizArr, err := s.getTopoBriefInfo(kit, common.BKInnerObjIDSet, bizSetBizCond)
		if err != nil {
			return nil, err
		}
		return bizArr, nil
	case common.BKInnerObjIDSet:
		childObjID = common.BKInnerObjIDModule
	case common.BKInnerObjIDModule:
		blog.Errorf("module's child(host) is not a mainline object, **forbidden to search**, rid: %s", kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjIDField)
	default:
		asstOpt := &metadata.QueryCondition{
			Condition: mapstr.MapStr{
				common.AssociationKindIDField: common.AssociationKindMainline,
				common.BKAsstObjIDField:       opt.ParentObjID,
			},
		}

		asst, err := s.Engine.CoreAPI.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, asstOpt)
		if err != nil {
			blog.Errorf("search mainline association failed, err: %v, cond: %#v, rid: %s", err, asstOpt, kit.Rid)
			return nil, err
		}

		if len(asst.Info) == 0 {
			blog.Errorf("parent object %s is not mainline, **forbidden to search**, rid: %s", opt.ParentObjID, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAsstObjIDField)
		}

		childObjID = asst.Info[0].ObjectID
	}

	// find topo nodes' id and name by parent id
	instArr, err := s.getTopoBriefInfo(kit, childObjID, mapstr.MapStr{common.BKParentIDField: opt.ParentID})
	if err != nil {
		return nil, err
	}

	// if there exists custom level, biz can have both default set as child and its custom level children
	if opt.ParentObjID == common.BKInnerObjIDApp && childObjID != common.BKInnerObjIDSet {
		setCond := mapstr.MapStr{
			common.BKParentIDField: opt.ParentID,
			common.BKDefaultField:  common.DefaultResSetFlag,
		}

		setArr, err := s.getTopoBriefInfo(kit, common.BKInnerObjIDSet, setCond)
		if err != nil {
			return nil, err
		}
		return append(instArr, setArr...), nil
	}
	return instArr, nil
}

// checkTopoNodeInBizSet check if topo node belongs to biz that is in the biz set, input contains the biz set scope cond
func (s *Service) checkTopoNodeInBizSet(kit *rest.Kit, objID string, instID int64, bizSetBizCond mapstr.MapStr) error {
	instOpt := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.GetInstIDField(objID): instID},
		Fields:         []string{common.BKAppIDField},
		Page:           metadata.BasePage{Limit: 1},
		DisableCounter: true,
	}
	instRes, err := s.Logics.InstOperation().FindInst(kit, objID, instOpt)
	if err != nil {
		blog.Errorf("find %s inst failed, err: %v, cond: %+v, rid: %s", objID, err, instOpt, kit.Rid)
		return err
	}

	if len(instRes.Info) == 0 {
		blog.Errorf("inst %s/%d is not exist, rid: %s", objID, instID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, objID)
	}

	bizCond := &metadata.Condition{
		Condition: map[string]interface{}{
			common.BKDBAND: []mapstr.MapStr{bizSetBizCond, {common.BKAppIDField: instRes.Info[0][common.BKAppIDField]}},
		},
	}
	resp, err := s.Engine.CoreAPI.CoreService().Instance().CountInstances(kit.Ctx, kit.Header,
		common.BKInnerObjIDApp, bizCond)
	if err != nil {
		blog.Errorf("count biz failed, err: %v, cond: %#v, rid: %s", err, bizCond, kit.Rid)
		return err
	}

	if resp.Count == 0 {
		blog.Errorf("instance biz does not belong to the biz set, biz cond: %#v, rid: %s", bizCond, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, objID)
	}

	return nil
}

// getTopoBriefInfo get topo id and name by condition and parse to the form of topo node, sort in the order of inst id
func (s *Service) getTopoBriefInfo(kit *rest.Kit, objID string, condition mapstr.MapStr) ([]mapstr.MapStr, error) {
	instIDField := metadata.GetInstIDFieldByObjID(objID)
	instNameField := metadata.GetInstNameFieldName(objID)

	instOpt := &metadata.QueryCondition{
		Fields:         []string{instIDField, instNameField},
		Page:           metadata.BasePage{Limit: common.BKNoLimit, Sort: instIDField},
		DisableCounter: true,
		Condition:      condition,
	}

	instRes, err := s.Logics.InstOperation().FindInst(kit, objID, instOpt)
	if err != nil {
		blog.Errorf("find %s inst failed, err: %v, cond: %#v, rid: %s", objID, err, instOpt, kit.Rid)
		return nil, err
	}

	topoArr := make([]mapstr.MapStr, len(instRes.Info))
	for index, inst := range instRes.Info {
		topoArr[index] = mapstr.MapStr{
			common.BKObjIDField:    objID,
			common.BKInstIDField:   inst[instIDField],
			common.BKInstNameField: inst[instNameField],
		}
	}

	return topoArr, nil
}
