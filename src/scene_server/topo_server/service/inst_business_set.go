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
	"sort"
	"strconv"

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
)

// CreateBusinessSet create a new business set
func (s *Service) CreateBusinessSet(ctx *rest.Contexts) {
	data := new(metadata.CreateBizSetRequest)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := data.Validate(); err != nil {
		blog.Errorf("validate create business set failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// do with transaction
	var bizSet mapstr.MapStr
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		bizSet, err = s.Logics.BusinessSetOperation().CreateBusinessSet(ctx.Kit, data)
		if err != nil {
			blog.Errorf("create business set failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}

		// register business set resource creator action to iam
		if auth.EnableAuthorize() {
			var bizSetID int64
			if bizSetID, err = bizSet.Int64(common.BKBizSetIDField); err != nil {
				blog.Errorf("get biz set id failed, err: %v, biz: %#v, rid: %s", err, bizSet, ctx.Kit.Rid)
				return err
			}

			var bizSetName string
			if bizSetName, err = bizSet.String(common.BKBizSetNameField); err != nil {
				blog.Errorf("get biz set name failed, err: %v, biz: %#v, rid: %s", err, bizSet, ctx.Kit.Rid)
				return err
			}

			iamInstance := metadata.IamInstanceWithCreator{
				Type:    string(iam.BizSet),
				ID:      strconv.FormatInt(bizSetID, 10),
				Name:    bizSetName,
				Creator: ctx.Kit.User,
			}
			_, err = s.AuthManager.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance)
			if err != nil {
				blog.Errorf("register created business set to iam failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}
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

// find business set list with these info：
// 1. have any authorized resources in a business.
// 2. only returned with a few field for this business info.
func (s *Service) PreviewBusinessSet(ctx *rest.Contexts) {
	searchCond := new(metadata.PreviewBusinessSetRequest)
	if err := ctx.DecodeInto(searchCond); err != nil {
		blog.Errorf("failed to parse the params, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if err := searchCond.Validate(); err != nil {
		blog.Errorf("bizPropertyFilter is illegal, err: %v, rid:%s", err, ctx.Kit.Rid)
		ccErr := ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, err.Error())
		ctx.RespAutoError(ccErr)
		return
	}

	opt := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			metadata.AttributeFieldObjectID:     common.BKInnerObjIDApp,
			metadata.AttributeFieldPropertyType: common.FieldTypeUser,
		},
		DisableCounter: true,
	}
	attrArr, err := s.Engine.CoreAPI.CoreService().Model().ReadModelAttr(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDApp, opt)
	if err != nil {
		blog.Errorf("failed get the business attribute, %s, rid:%s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	// userFieldArr Fields in the business are user-type fields
	var userFields []string
	for _, attribute := range attrArr.Info {
		userFields = append(userFields, attribute.PropertyID)
	}

	mgoFilter, key, err := searchCond.BizSetPropertyFilter.ToMgo()
	if err != nil {
		blog.Errorf("BizPropertyFilter ToMgo failed: %s, err: %v,  rid:%s", searchCond.BizSetPropertyFilter,
			err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid,
			err.Error()+fmt.Sprintf(", biz_property_filter.%s", key)))
		return
	}

	query := &metadata.QueryCondition{
		Condition: mgoFilter,
		Fields:    []string{common.BKAppIDField, common.BKAppNameField},
		Page: metadata.BasePage{
			Limit: common.BKMaxInstanceLimit,
			Sort:  common.BKAppIDField,
		},
	}

	cnt, instItems, err := s.Logics.BusinessOperation().FindBiz(ctx.Kit, query)
	if err != nil {
		blog.Errorf("find business failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := make(mapstr.MapStr)
	result.Set("count", cnt)
	result.Set("info", instItems)

	ctx.RespEntity(result)
}

// find business set list with these info：
// 1. have any authorized resources in a business set.
// 2. only returned with a few field for this business set info.
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

	if s.AuthManager.Enabled() {
		authInput := meta.ListAuthorizedResourcesParam{
			UserName:     ctx.Kit.User,
			ResourceType: meta.BizSet,
			Action:       meta.AccessBizSet,
		}
		authorizedRes, err := s.AuthManager.Authorizer.ListAuthorizedResources(ctx.Kit.Ctx, ctx.Kit.Header, authInput)
		if err != nil {
			blog.Errorf("get authorized resource failed, user: %s, err: %v, rid: %s", ctx.Kit.User, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoGetAuthorizedBusinessSetListFailed))
			return
		}

		// if isAny is false,we should add bizIds condition.
		if !authorizedRes.IsAny {
			for _, resourceID := range authorizedRes.Ids {
				bizSetID, err := strconv.ParseInt(resourceID, 10, 64)
				if err != nil {
					blog.Errorf("parse bizSetID %s failed, err: %v, rid: %s", bizSetID, err, ctx.Kit.Rid)
					ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
					return
				}
				bizSetList = append(bizSetList, bizSetID)
			}
			if len(bizSetList) == 0 {
				ctx.RespEntity(make([]interface{}, 0))
				return
			}
			// sort for prepare to find business with page.
			sort.Sort(util.Int64Slice(bizSetList))
		}
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
	result, err := s.Logics.BusinessSetOperation().FindBizSet(ctx.Kit, query)
	if nil != err {
		blog.Errorf("failed to find the biz set list, error is %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespEntity(nil)
		return
	}

	ctx.RespEntity(result.Info)
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
			common.BKDefaultField: mapstr.MapStr{common.BKDBNE: common.DefaultAppFlag},
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

// getBizSetIDList 通过用户所传条件算出的的bizSetIDs 与此用户拥有的bizSetIDs 做交集，得到用户最终能够获取到的业务集
func (s *Service) getBizSetIDList(ctx *rest.Contexts, searchCond *metadata.QueryBusinessSetRequest, bizSetIDs []int64) (
	[]int64, error) {

	// 最终有权限的biz set list
	authBizSetIDs := make([]int64, 0)

	authInput := meta.ListAuthorizedResourcesParam{
		UserName:     ctx.Kit.User,
		ResourceType: meta.BizSet,
		Action:       meta.Find,
	}
	authorizedRes, err := s.AuthManager.Authorizer.ListAuthorizedResources(ctx.Kit.Ctx, ctx.Kit.Header, authInput)
	if err != nil {
		blog.Errorf("search business failed, list authorized resources failed, user: %s, err: %v, rid: %s",
			ctx.Kit.User, err, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrorTopoGetAuthorizedBusinessListFailed, "")
		return []int64{}, err
	}

	// if isAny is true means we have all bizIds authority, else we should parse ids list that we have authority.
	if authorizedRes.IsAny {
		// if user assign the ids,add the ids to the condition.
		authBizSetIDs = bizSetIDs

	} else {
		bizAuthSetList := make([]int64, 0)
		for _, resourceID := range authorizedRes.Ids {
			bizSetID, err := strconv.ParseInt(resourceID, 10, 64)
			if err != nil {
				blog.Errorf("parse bizID: %s, failed, err: %v, rid: %s", bizSetID, err, ctx.Kit.Rid)
				ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
				return []int64{}, err
			}
			bizAuthSetList = append(bizAuthSetList, bizSetID)
		}
		// this means that user want to find a specific business.now we check if he has this authority.
		for _, bizSetID := range bizSetIDs {
			if util.InArray(bizSetID, bizAuthSetList) {
				// authBizIDs store the authorized bizIDs
				authBizSetIDs = append(authBizSetIDs, bizSetID)
			}
		}

	}
	return authBizSetIDs, nil
}

// searchBizSetByUserCondition  获取包含bk_biz_set_id 的查询结果。
func (s *Service) searchBizSetByUserCondition(ctx *rest.Contexts, filter *metadata.QueryBusinessSetRequest) (
	*metadata.CommonSearchResult, error) {
	query := &metadata.CommonSearchFilter{
		ObjectID:      common.BKInnerObjIDBizSet,
		Fields:        filter.Fields,
		Page:          filter.Page,
		TimeCondition: filter.TimeCondition,
		Conditions:    filter.BizSetPropertyFilter,
	}

	result, err := s.Logics.BusinessSetOperation().FindBizSet(ctx.Kit, query)
	if nil != err {
		blog.Errorf("failed to find the objects(%s), error info is %s, rid: %s",
			ctx.Request.PathParameter("obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return &metadata.CommonSearchResult{}, err
	}
	return result, nil
}

func getBusinessSetResult(bizSetIDs []int64, bizSetFlag bool,
	bizSetResult *metadata.CommonSearchResult) (*metadata.CommonSearchResult, error) {
	result := new(metadata.CommonSearchResult)

	for _, info := range bizSetResult.Info {
		if _, ok := info.(*mapstr.MapStr); !ok {
			blog.Errorf("biz set result type error,info: %+v", info)
			return nil, fmt.Errorf("biz set result type error")
		}
		value := *info.(*mapstr.MapStr)
		bizSetID, err := util.GetInt64ByInterface(value[common.BKBizSetIDField])
		if err != nil {
			return nil, err
		}
		if util.InArray(bizSetID, bizSetIDs) {
			if !bizSetFlag {
				delete(info.(map[string]interface{}), common.BKBizSetIDField)
			}
			result.Info = append(result.Info, info)
		}
	}
	return result, nil
}

// SearchBusiness search the business by condition
func (s *Service) SearchBusinessSet(ctx *rest.Contexts) {

	searchCond := new(metadata.QueryBusinessSetRequest)
	if err := ctx.DecodeInto(searchCond); err != nil {
		blog.Errorf("failed to parse the params, error %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	// BizSetPropertyFilter
	if searchCond.BizSetPropertyFilter != nil {
		if err := searchCond.Validate(); err != nil {
			blog.Errorf("the params is illegal, error is %v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "")
			return
		}

	}
	if searchCond.Page.Sort == "" {
		searchCond.Page.Sort = common.BKBizSetIDField
	}
	// 这个标记标明用户是否需要返回field
	bizSetFlag := false
	fields := make([]string, 0)

	if len(searchCond.Fields) != 0 {
		for _, field := range searchCond.Fields {
			if field == common.BKBizSetIDField {
				bizSetFlag = true
			}
			fields = append(fields, field)
		}
		if !bizSetFlag {
			fields = append(fields, common.BKBizSetIDField)
		}
	} else {
		bizSetFlag = true
	}

	bizSetResult, err := s.searchBizSetByUserCondition(ctx, searchCond)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "")
		return
	}
	if len(bizSetResult.Info) == 0 {
		ctx.RespEntity(&metadata.CommonSearchResult{})
		return
	}

	// 需要鉴权的 biz set id list
	userBizSetIDs := make([]int64, 0)
	for _, info := range bizSetResult.Info {
		if _, ok := info.(*mapstr.MapStr); !ok {
			ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "")
			return
		}
		value := *info.(*mapstr.MapStr)

		bizSetID, err := util.GetInt64ByInterface(value[common.BKBizSetIDField])
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "")
			return
		}
		userBizSetIDs = append(userBizSetIDs, bizSetID)
	}
	// 初始化将用户的业务ID
	bizSetIDs := userBizSetIDs

	if s.AuthManager.Enabled() {
		bizSetIDs, err = s.getBizSetIDList(ctx, searchCond, userBizSetIDs)
		if err != nil {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKBizSetIDField))
			return
		}
	}

	result := new(metadata.CommonSearchResult)

	// 1、只返回用户需要查询的并且有权限的业务集.2、用户没有指定获取biz_set_id场景下需要把biz_set_id删掉
	result, err = getBusinessSetResult(bizSetIDs, bizSetFlag, bizSetResult)
	if err != nil {
		blog.Errorf("get business result fail err: %v", err)
		ctx.RespErrorCodeOnly(common.CCErrCommParseDataFailed, "")

	}
	ctx.RespEntity(result)
}

// CountBusinessSet count the business by condition
func (s *Service) CountBusinessSet(ctx *rest.Contexts) {

	searchCond := new(metadata.QueryBusinessSetRequest)
	if err := ctx.DecodeInto(&searchCond); err != nil {
		blog.Errorf("failed to parse the params, error %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}
	if searchCond.BizSetPropertyFilter != nil {
		if err := searchCond.Validate(); err != nil {
			blog.Errorf("the params is illegal, error is %v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "")
			return
		}

	}

	searchCond.Page = metadata.BasePage{
		Limit: common.BKNoLimit,
		Sort:  common.BKBizSetIDField,
	}
	searchCond.Fields = []string{common.BKBizSetIDField}
	bizSetResult, err := s.searchBizSetByUserCondition(ctx, searchCond)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "")
		return
	}

	if len(bizSetResult.Info) == 0 {
		ctx.RespEntity(0)
		return
	}

	// 需要鉴权的 biz set id list.
	userBizSetIDs := make([]int64, 0)
	for _, info := range bizSetResult.Info {
		if _, ok := info.(*mapstr.MapStr); !ok {
			ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "")
			return
		}
		value := *info.(*mapstr.MapStr)
		bizSetID, err := util.GetInt64ByInterface(value[common.BKBizSetIDField])
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "")
			return
		}
		userBizSetIDs = append(userBizSetIDs, bizSetID)
	}

	// 初始化将用户的业务ID.
	bizSetIDs := userBizSetIDs

	if s.AuthManager.Enabled() {
		bizSetIDs, err = s.getBizSetIDList(ctx, searchCond, userBizSetIDs)
		if err != nil {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKBizSetIDField))
			return
		}
	}
	setIds := make([]int64, 0)
	for _, info := range bizSetResult.Info {
		if _, ok := info.(*mapstr.MapStr); !ok {
			blog.Errorf("biz set result type error,info: %+v", info)
			ctx.RespEntity(fmt.Errorf("biz set result type error"))
			return
		}
		value := *info.(*mapstr.MapStr)
		bizSetID, err := util.GetInt64ByInterface(value[common.BKBizSetIDField])
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "")
			return
		}
		if util.InArray(bizSetID, bizSetIDs) {
			setIds = append(setIds, bizSetID)
		}
	}
	ctx.RespEntity(len(setIds))
}

func (s *Service) findBizSetTopo(kit *rest.Kit,
	opt *metadata.FindBizSetTopoOption) ([]mapstr.MapStr, error) {
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
