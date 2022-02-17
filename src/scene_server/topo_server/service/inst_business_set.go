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
	"fmt"
	"strconv"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
)

// validateScopeFields validate if scope fields are all enum/organization type
func (s *Service) validateScopeFields(kit *rest.Kit, fields []string) error {
	// biz id field is allowed to use in biz set scope, exclude it in validation
	validFields := make([]string, 0)
	for _, field := range fields {
		if field != common.BKAppIDField {
			validFields = append(validFields, field)
		}
	}

	cond := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKPropertyIDField: map[string]interface{}{
				common.BKDBIN: validFields,
			},
		},
		Fields: []string{common.BKPropertyTypeField},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
	}
	res, err := s.Engine.CoreAPI.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header,
		common.BKInnerObjIDApp, cond)
	if err != nil {
		blog.Errorf("read model attribute failed, error: %v, cond: %+v, rid: %s", err, cond, kit.Rid)
		return err
	}
	// 数量上必须与查询一致
	if res.Count != int64(len(validFields)) {
		blog.Errorf("read model attribute failed, error: %v, cond: %+v, rid: %s", err, cond, kit.Rid)
		return err
	}

	// 严格校验每个字段类型是否是enum或 organization
	for _, info := range res.Info {
		if info.PropertyType != common.FieldTypeEnum && info.PropertyType != common.FieldTypeOrganization {
			blog.Errorf("model attribute type must be enum or organization, cond: %+v, rid: %s", cond, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "model attribute must enum or organization")
		}
	}
	return nil
}

// CreateBusinessSet create a new business set
func (s *Service) CreateBusinessSet(ctx *rest.Contexts) {
	data := new(metadata.CreateBizSetRequest)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	fields, errRaw := data.Validate()
	if errRaw.ErrCode != 0 {
		blog.Errorf("validate create business set failed, err: %v, rid: %s", errRaw, ctx.Kit.Rid)
		ctx.RespAutoError(errRaw.ToCCError(ctx.Kit.CCError))
		return
	}

	if err := s.validateScopeFields(ctx.Kit, fields); err != nil {
		ctx.RespAutoError(err)
		return
	}

	var bizSet mapstr.MapStr
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		bizSet, err = s.Logics.BusinessSetOperation().CreateBusinessSet(ctx.Kit, data)
		if err != nil {
			blog.Errorf("create business set failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	bizSetID, err := bizSet.Int64(common.BKBizSetIDField)
	if err != nil {
		blog.Errorf("get biz set id failed, err: %v, biz: %#v, rid: %s", err, bizSet, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(bizSetID)
}

// UpdateBizSet update business set
func (s *Service) UpdateBizSet(ctx *rest.Contexts) {
	opt := new(metadata.UpdateBizSetOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(opt.BizSetIDs) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "bk_biz_set_ids"))
		return
	}

	if opt.Data == nil || (opt.Data.BizSetAttr == nil && opt.Data.Scope == nil) {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "data"))
		return
	}

	updateData := make(mapstr.MapStr)
	if opt.Data.BizSetAttr != nil {
		updateData = opt.Data.BizSetAttr
	}

	// do not allow batch update biz set name and scope
	if len(opt.BizSetIDs) > 1 {
		if _, exists := updateData[common.BKBizSetNameField]; exists {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKBizSetNameField))
			return
		}

		if opt.Data.Scope != nil {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKBizSetScopeField))
			return
		}
	}

	// validate scope field
	if opt.Data.Scope != nil {
		fields, err := opt.Data.Scope.Validate()
		if err != nil {
			blog.Errorf("validate business set scope failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKBizSetScopeField))
			return
		}

		if err := s.validateScopeFields(ctx.Kit, fields); err != nil {
			ctx.RespAutoError(err)
			return
		}
	}

	bizSetFilter := mapstr.MapStr{
		common.BKBizSetIDField: mapstr.MapStr{common.BKDBIN: opt.BizSetIDs},
	}

	if opt.Data.Scope != nil {
		updateData[common.BKBizSetScopeField] = opt.Data.Scope
	}

	// update biz set instances
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logics.InstOperation().UpdateInst(ctx.Kit, bizSetFilter, updateData, common.BKInnerObjIDBizSet)
		if err != nil {
			blog.Errorf("update biz set failed, err: %v, opt: %#v, rid: %s", err, opt, ctx.Kit.Rid)
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

// PreviewBusinessSet  此预览接口用于创建业务集过程中的预览，支持进行条件匹配
func (s *Service) PreviewBusinessSet(ctx *rest.Contexts) {
	searchCond := new(metadata.PreviewBusinessSetRequest)
	if err := ctx.DecodeInto(searchCond); err != nil {
		blog.Errorf("failed to parse the params, error: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	errRaw := searchCond.Validate(true)
	if errRaw.ErrCode != 0 {
		blog.Errorf("biz property filter is illegal, err: %v, rid: %s", errRaw, ctx.Kit.Rid)
		ctx.RespAutoError(errRaw.ToCCError(ctx.Kit.CCError))
		return
	}

	mgoFilter := make(map[string]interface{})

	if searchCond.BizSetPropertyFilter != nil {
		filter, key, err := searchCond.BizSetPropertyFilter.ToMgo()
		if err != nil {
			blog.Errorf("BizPropertyFilter ToMgo failed: %s, err: %v, rid: %s", searchCond.BizSetPropertyFilter,
				err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid,
				fmt.Sprintf(", biz_property_filter.%s", key)))
			return
		}
		mgoFilter = filter
	}
	if searchCond.Filter != nil {
		filter, key, err := searchCond.Filter.ToMgo()
		if err != nil {
			blog.Errorf("BizPropertyFilter ToMgo failed: %s, err: %v, rid: %s", searchCond.Filter, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, key))
			return
		}
		mgoFilter = filter
	}

	mgoFilter[common.BKDataStatusField] = mapstr.MapStr{common.BKDBNE: common.DataStatusDisabled}
	mgoFilter[common.BKDefaultField] = 0
	query := new(metadata.QueryCondition)
	bizSetResult := new(metadata.QueryBusinessSetResponse)
	if searchCond.Page.EnableCount {
		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKTableNameBaseApp, []map[string]interface{}{mgoFilter})
		if err != nil {
			blog.Errorf("count biz failed, err: %v, cond: %#v, rid: %s", err, mgoFilter, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		bizSetResult.Count = int(counts[0])
		ctx.RespEntity(bizSetResult)
		return
	}
	query = &metadata.QueryCondition{
		Condition: mgoFilter,
		Fields:    []string{common.BKAppIDField, common.BKAppNameField},
		Page: metadata.BasePage{
			Start: searchCond.Page.Start,
			Limit: searchCond.Page.Limit,
			Sort:  common.BKAppIDField,
		},
	}

	_, instItems, err := s.Logics.BusinessOperation().FindBiz(ctx.Kit, query)
	if err != nil {
		blog.Errorf("find business failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// 底层默认统一加了default返回值，由于与前端约定只返回id和name，所以需要将default去掉，
	for _, item := range instItems {
		inst := mapstr.New()
		inst[common.BKAppIDField] = item[common.BKAppIDField]
		inst[common.BKAppNameField] = item[common.BKAppNameField]
		bizSetResult.Info = append(bizSetResult.Info, inst)
	}
	ctx.RespEntity(bizSetResult)
}

// SearchReducedBusinessSetList 此接口只用于前端左侧下拉栏的查询，只需要返回id和name即可
func (s *Service) SearchReducedBusinessSetList(ctx *rest.Contexts) {

	page := metadata.BasePage{
		Limit: common.BKNoLimit,
	}
	sortParam := ctx.Request.QueryParameter("sort")
	if len(sortParam) > 0 {
		page.Sort = sortParam
	} else {
		page.Sort = common.BKBizSetIDField
	}

	bizSetList := make([]int64, 0)
	bizSetResult := new(metadata.QueryBusinessSetResponse)
	if s.AuthManager.Enabled() {

		flag, authBizSetList, err := s.getAuthBizSetIDList(ctx.Kit, meta.AccessBizSet)
		if err != nil {
			blog.Errorf("get authorized biz set id failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespEntity(bizSetResult)
			return
		}
		if !flag && len(authBizSetList) == 0 {
			ctx.RespEntity(bizSetResult)
			return
		}
		bizSetList = authBizSetList
	}

	query := &metadata.CommonSearchFilter{
		ObjectID: common.BKInnerObjIDBizSet,
		Fields:   []string{common.BKBizSetIDField, common.BKBizSetNameField},
		Page:     page,
	}

	if len(bizSetList) > 0 {
		query.Conditions = &querybuilder.QueryFilter{
			Rule: &querybuilder.AtomRule{
				Field:    common.BKBizSetIDField,
				Operator: querybuilder.OperatorIn,
				Value:    bizSetList,
			},
		}
	}
	result, err := s.Logics.InstOperation().SearchObjectInstances(ctx.Kit, common.BKInnerObjIDBizSet, query)
	if err != nil {
		blog.Errorf("failed to find the biz set list, error is %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrorTopoGetAuthorizedBusinessSetListFailed, err.Error()))
		return
	}
	bizSetResult.Info = result.Info
	ctx.RespEntity(bizSetResult)
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

	var bizFilter mapstr.MapStr
	if opt.Filter != nil {
		cond, errKey, rawErr := opt.Filter.ToMgo()
		if rawErr != nil {
			blog.Errorf("parse biz filter(%#v) failed, err: %v, rid: %s", opt.Filter, rawErr, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errKey))
			return
		}
		bizFilter = cond
	}

	// get biz mongo condition by biz scope in biz set
	bizSetBizCond, err := s.getBizSetBizCond(ctx.Kit, opt.BizSetID)
	if err != nil {
		blog.Errorf("get biz cond by biz set id %d failed, err: %v, rid: %s", opt.BizSetID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// merge biz set scope mongo condition with extra biz condition to search specific biz in all biz in biz set
	if len(bizFilter) > 0 {
		bizSetBizCond = mapstr.MapStr{common.BKDBAND: []mapstr.MapStr{bizSetBizCond, bizFilter}}
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
		ctx.RespEntityWithCount(counts[0], make([]mapstr.MapStr, 0))
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
	ctx.RespEntityWithCount(0, biz)
}

// getBizSetBizCond get biz mongo condition from the biz set scope
func (s *Service) getBizSetBizCond(kit *rest.Kit, bizSetID int64) (mapstr.MapStr, error) {
	bizSetCond := &metadata.QueryCondition{
		Fields:         []string{common.BKBizSetScopeField},
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

	if bizSetRes.Data.Info[0].Scope.MatchAll {
		// do not include resource pool biz in biz set by default
		return mapstr.MapStr{
			common.BKDefaultField:    mapstr.MapStr{common.BKDBNE: common.DefaultAppFlag},
			common.BKDataStatusField: map[string]interface{}{common.BKDBNE: common.DataStatusDisabled},
		}, nil
	}

	if bizSetRes.Data.Info[0].Scope.Filter == nil {
		blog.Errorf("biz set(%#v) has no filter and is not match all, rid: %s", bizSetRes.Data.Info[0].Scope, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKBizSetIDField)
	}

	bizSetBizCond, errKey, rawErr := bizSetRes.Data.Info[0].Scope.Filter.ToMgo()
	if rawErr != nil {
		blog.Errorf("parse biz set scope(%#v) failed, err: %v, rid: %s", bizSetRes.Data.Info[0].Scope, rawErr, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errKey)
	}

	// do not include resource pool biz in biz set by default
	if _, exists := bizSetBizCond[common.BKDefaultField]; !exists {
		bizSetBizCond[common.BKDefaultField] = mapstr.MapStr{common.BKDBNE: common.DefaultAppFlag}
	}

	// do not include disabled biz in biz set by default
	if _, exists := bizSetBizCond[common.BKDataStatusField]; !exists {
		bizSetBizCond[common.BKDataStatusField] = map[string]interface{}{common.BKDBNE: common.DataStatusDisabled}
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

// getBizSetIDList 获取有权限的biz set ids
func (s *Service) getAuthBizSetIDList(kit *rest.Kit, action meta.Action) (bool, []int64, error) {

	// 最终有权限的biz set list
	authBizSetIDs := make([]int64, 0)

	authInput := meta.ListAuthorizedResourcesParam{
		UserName:     kit.User,
		ResourceType: meta.BizSet,
		Action:       action,
	}
	authorizedRes, err := s.AuthManager.Authorizer.ListAuthorizedResources(kit.Ctx, kit.Header, authInput)
	if err != nil {
		blog.Errorf("search business failed, list authorized resources failed, user: %s, err: %v, rid: %s",
			kit.User, err, kit.Rid)
		return false, []int64{}, err
	}

	// if isAny is true means we have all bizSetIds authority, else we should parse ids list that we have authority.
	if authorizedRes.IsAny {
		// if user assign the ids,add the ids to the condition.
		return true, []int64{}, nil

	} else {
		for _, resourceID := range authorizedRes.Ids {
			bizSetID, err := strconv.ParseInt(resourceID, 10, 64)
			if err != nil {
				blog.Errorf("parse bizID: %s, failed, err: %v, rid: %s", bizSetID, err, kit.Rid)
				return false, []int64{}, err
			}
			authBizSetIDs = append(authBizSetIDs, bizSetID)
		}

	}
	return false, authBizSetIDs, nil
}

// SearchBusinessSet search business set by condition
func (s *Service) SearchBusinessSet(ctx *rest.Contexts) {

	searchCond := new(metadata.QueryBusinessSetRequest)
	if err := ctx.DecodeInto(searchCond); err != nil {
		blog.Errorf("failed to parse the params, error %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	errRaw := searchCond.Validate(false)
	if errRaw.ErrCode != 0 {
		blog.Errorf("biz property filter is illegal, err: %v, rid: %s", errRaw, ctx.Kit.Rid)
		ctx.RespAutoError(errRaw.ToCCError(ctx.Kit.CCError))
		return
	}

	cond := new(metadata.CommonSearchFilter)
	bizSetResult := new(metadata.QueryBusinessSetResponse)
	if s.AuthManager.Enabled() {
		flag, authbizSetIDs, err := s.getAuthBizSetIDList(ctx.Kit, meta.Find)
		if err != nil {
			blog.Errorf("get authorized biz set id failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespEntity(bizSetResult)
			return
		}
		if !flag && len(authbizSetIDs) == 0 {
			ctx.RespEntity(bizSetResult)
			return
		}
		if flag {
			cond.Conditions = searchCond.BizSetPropertyFilter
		} else {
			cond.Conditions = &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionAnd, Rules: []querybuilder.Rule{&querybuilder.AtomRule{
						Field: common.BKBizSetIDField, Operator: querybuilder.OperatorIn, Value: authbizSetIDs},
						searchCond.BizSetPropertyFilter}}}
		}
	} else {
		cond.Conditions = searchCond.BizSetPropertyFilter
	}
	if searchCond.Page.EnableCount {
		countCond := &metadata.CommonCountFilter{
			ObjectID:      common.BKInnerObjIDBizSet,
			Conditions:    searchCond.BizSetPropertyFilter,
			TimeCondition: searchCond.TimeCondition}

		result, err := s.Logics.InstOperation().CountObjectInstances(ctx.Kit, common.BKInnerObjIDBizSet, countCond)
		if err != nil {
			blog.Errorf("count object %s instances failed, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntity(result)
		return
	}

	if searchCond.Page.Sort == "" {
		searchCond.Page.Sort = common.BKBizSetIDField
	}
	cond.TimeCondition = searchCond.TimeCondition
	cond.Fields = searchCond.Fields
	cond.Page = searchCond.Page
	cond.ObjectID = common.BKInnerObjIDBizSet
	result, err := s.Logics.InstOperation().SearchObjectInstances(ctx.Kit, common.BKInnerObjIDBizSet, cond)
	if err != nil {
		blog.Errorf("failed to find the biz set, error info is %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(result.Info) == 0 {
		ctx.RespEntity(bizSetResult)
		return
	}

	bizSetResult.Info = result.Info
	ctx.RespEntity(bizSetResult)
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

	// get parent object id to check if the parent node is a valid mainline instance that belongs to the biz set
	var childObjID string
	switch opt.ParentObjID {
	case common.BKInnerObjIDBizSet:
		if opt.ParentID != opt.BizSetID {
			blog.Errorf("biz parent id %d is not equal to biz set id %d, rid: %s", opt.ParentID, opt.BizSetID, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKParentIDField)
		}

		// find biz nodes by the condition in biz sets
		bizArr, err := s.getTopoBriefInfo(kit, common.BKInnerObjIDApp, bizSetBizCond)
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

	// check if the parent node belongs to a biz that is in the biz set
	if err := s.checkTopoNodeInBizSet(kit, opt.ParentObjID, opt.ParentID, bizSetBizCond); err != nil {
		blog.Errorf("check if parent %s node %d in biz failed, err: %v, biz cond: %#v, rid: %s", opt.ParentObjID,
			opt.ParentID, err, bizSetBizCond, kit.Rid)
		return nil, err
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
		return append(setArr, instArr...), nil
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
		Fields:         []string{instIDField, instNameField, common.BKDefaultField},
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
			common.BKDefaultField:  inst[common.BKDefaultField],
		}
	}

	return topoArr, nil
}

// CountBizSetTopoHostAndSrvInst count hosts and service instances in topo node under the biz set. **only for ui**
func (s *Service) CountBizSetTopoHostAndSrvInst(ctx *rest.Contexts) {
	urlBizSetID := ctx.Request.PathParameter(common.BKBizSetIDField)
	bizSetID, err := strconv.ParseInt(urlBizSetID, 10, 64)
	if err != nil {
		blog.Errorf("parse biz set id: %s from url failed, err: %v , rid: %s", urlBizSetID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	input := new(metadata.HostAndSerInstCountOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(input.Condition) > 20 {
		err := ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "condition length exceed 20")
		ctx.RespAutoError(err)
		return
	}

	// validate if the topo nodes are all in the biz set
	objInstMap := make(map[string][]int64)
	for _, node := range input.Condition {
		objInstMap[node.ObjID] = append(objInstMap[node.ObjID], node.InstID)
	}

	bizIDs := make([]interface{}, 0)
	bizIDMap := make(map[interface{}]struct{})
	for objID, instIDs := range objInstMap {
		distinctOpt := &metadata.DistinctFieldOption{
			TableName: common.GetInstTableName(objID, ctx.Kit.SupplierAccount),
			Field:     common.BKAppIDField,
			Filter:    map[string]interface{}{common.GetInstIDField(objID): mapstr.MapStr{common.BKDBIN: instIDs}},
		}

		distinctIDs, err := s.Engine.CoreAPI.CoreService().Common().GetDistinctField(ctx.Kit.Ctx, ctx.Kit.Header, distinctOpt)
		if err != nil {
			blog.Errorf("get biz ids failed, distinct opt: %+v, err: %v, rid: %s", distinctOpt, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		for _, id := range distinctIDs {
			if _, exists := bizIDMap[id]; !exists {
				bizIDMap[id] = struct{}{}
				bizIDs = append(bizIDs, id)
			}
		}
	}

	bizSetBizCond, err := s.getBizSetBizCond(ctx.Kit, bizSetID)
	if err != nil {
		blog.Errorf("get biz cond by biz set id %d failed, err: %v, rid: %s", bizSetID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	filter := []map[string]interface{}{{
		common.BKDBAND: []map[string]interface{}{
			bizSetBizCond, {common.BKAppIDField: mapstr.MapStr{common.BKDBIN: bizIDs}},
		}}}
	counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKTableNameBaseApp, filter)
	if err != nil {
		blog.Errorf("count topo nodes failed, err: %v, filter: %#v, rid: %s", err, filter, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(counts) != 1 || (int(counts[0]) != len(bizIDs)) {
		blog.Errorf("topo nodes are not all in biz set, biz ids: %v, rid: %s", bizIDs, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "condition"))
		return
	}

	result, err := s.Logics.InstAssociationOperation().TopoNodeHostAndSerInstCount(ctx.Kit, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}