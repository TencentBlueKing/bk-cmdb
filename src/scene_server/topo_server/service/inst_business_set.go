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
			if bizSetID, err = bizSet.Int64(common.BKAppSetIDField); err != nil {
				blog.Errorf("get biz set id failed, err: %v, biz: %#v, rid: %s", err, bizSet, ctx.Kit.Rid)
				return err
			}

			var bizSetName string
			if bizSetName, err = bizSet.String(common.BKAppSetNameField); err != nil {
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

	bizSetID, err := bizSet.Int64(common.BKAppSetIDField)
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
	if err := ctx.DecodeInto(&searchCond); err != nil {
		blog.Errorf("failed to parse the params, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}
	defErr := ctx.Kit.CCError

	if err := searchCond.Validate(); err != nil {
		blog.Errorf("bizPropertyFilter is illegal, err: %v, rid:%s", err, ctx.Kit.Rid)
		ccErr := defErr.CCErrorf(common.CCErrCommParamsInvalid, err.Error())
		ctx.RespAutoError(ccErr)
		return
	}
	// Only one of biz_property_filter and condition parameters can take effect, and condition is not recommended to
	// continue to use it.

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
		page.Sort = common.BKAppSetIDField
	}

	// 此场景下获取全部的业务列表，只需要返回业务集的id和name
	if errKey, err := page.Validate(true); err != nil {
		blog.Errorf("page parameter invalid, errKey: %v, err: %s, rid: %s", errKey, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errKey))
		return
	}

	bizSetList := make([]int64, 0)

	if s.AuthManager.Enabled() {
		authInput := meta.ListAuthorizedResourcesParam{
			UserName:     ctx.Kit.User,
			ResourceType: meta.BizSet,

			// TODO: 后续修改成业务集的访问权限
			Action: meta.ViewBusinessResource,
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
				ctx.RespEntityWithCount(0, make([]mapstr.MapStr, 0))
				return
			}
			// sort for prepare to find business with page.
			sort.Sort(util.Int64Slice(bizSetList))
		}
	}

	query := &metadata.CommonSearchFilter{
		Conditions: &querybuilder.QueryFilter{
			Rule: &querybuilder.AtomRule{
				Field:    common.BKAppSetIDField,
				Operator: querybuilder.OperatorIn,
				Value:    bizSetList,
			},
		},
		Fields: []string{common.BKAppSetIDField, common.BKAppSetNameField},
		Page:   page,
	}
	result, err := s.Logics.BusinessSetOperation().FindBizSet(ctx.Kit, query)
	if nil != err {
		blog.Errorf("failed to find the biz set list, error is %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespEntity(nil)
	}

	ctx.RespEntity(result)
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

func (s *Service) getAuthBizSetIDList(ctx *rest.Contexts, mgoFilter map[string]interface{}, searchCond *metadata.QueryBusinessSetRequest, bizSetIDs []int64) ([]int64, error) {

	// 用户传来的biz set id list
	//bizSetIDs := make([]int64, 0)

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
	bizAuthSetList := make([]int64, 0)

	// if isAny is true means we have all bizIds authority, else we should parse ids list that we have authority.
	if authorizedRes.IsAny {
		// if user assign the ids,add the ids to the condition.
		if len(bizSetIDs) > 0 {
			authBizSetIDs = bizSetIDs
		}

	} else {
		for _, resourceID := range authorizedRes.Ids {
			bizSetID, err := strconv.ParseInt(resourceID, 10, 64)
			if err != nil {
				blog.Errorf("parse bizID: %s, failed, err: %v, rid: %s", bizSetID, err, ctx.Kit.Rid)
				ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
				return []int64{}, err
			}
			bizAuthSetList = append(bizAuthSetList, bizSetID)
		}
		if len(bizSetIDs) > 0 {
			// this means that user want to find a specific business.now we check if he has this authority.
			for _, bizSetID := range bizAuthSetList {
				if util.InArray(bizSetID, bizAuthSetList) {
					// authBizIDs store the authorized bizIDs
					authBizSetIDs = append(authBizSetIDs, bizSetID)
				}
			}
			if len(authBizSetIDs) > 0 {
				mgoFilter[common.BKAppIDField] = mapstr.MapStr{common.BKDBIN: authBizSetIDs}
			} else {
				// if there are no qualified bizIDs, return null

				if searchCond.Page.EnableCount {

					ctx.RespEntity(&metadata.CommonCountResult{})
				} else {
					ctx.RespEntity(&metadata.CommonSearchResult{})
				}

				return []int64{}, err
			}
			// now you have the authority.
		} else {
			if len(bizAuthSetList) == 0 {
				if searchCond.Page.EnableCount {
					ctx.RespEntity(&metadata.CommonCountResult{})
				} else {
					ctx.RespEntity(&metadata.CommonSearchResult{})
				}
				return []int64{}, err
			}
			authBizSetIDs = bizAuthSetList
			// sort for prepare to find business with page.
			sort.Sort(util.Int64Slice(authBizSetIDs))
			// user can only find business that is already authorized.
		}
	}
	return authBizSetIDs, nil
}

// SearchBusiness search the business by condition
func (s *Service) SearchBusinessSet(ctx *rest.Contexts) {
	searchCond := new(metadata.QueryBusinessSetRequest)
	if err := ctx.DecodeInto(&searchCond); err != nil {
		blog.Errorf("failed to parse the params, error %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	if err := searchCond.Validate(); err != nil {
		blog.Errorf("the params is illegal, error is %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "")
		return
	}

	// 用户传来的biz set id list
	bizSetIDs := make([]int64, 0)

	// 最终有权限的biz set list
	authBizSetIDs := make([]int64, 0)

	mgoFilter, key, err := searchCond.BizSetPropertyFilter.ToMgo()
	if err != nil {
		blog.Errorf("BizPropertyFilter ToMgo failed: %s, err: %v, rid:%s", searchCond.BizSetPropertyFilter,
			err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid,
			err.Error()+fmt.Sprintf(", biz property filter.%s", key)))
		return
	}

	biz, exist := mgoFilter[common.BKAppSetIDField]
	if exist {
		// constrict that bk_biz_id field can only be a numeric value,
		// operators like or/in/and is not allowed.
		if bizSetCond, ok := biz.(map[string]interface{}); ok {
			if cond, ok := bizSetCond["$eq"]; ok {
				bizSetID, err := util.GetInt64ByInterface(cond)
				if err != nil {
					ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "", common.BKAppSetIDField)
					return
				}
				bizSetIDs = []int64{bizSetID}
			}
			if cond, ok := bizSetCond["$in"]; ok {
				if conds, ok := cond.([]interface{}); ok {
					for _, c := range conds {
						bizSetID, err := util.GetInt64ByInterface(c)
						if err != nil {
							ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "", common.BKAppSetIDField)
							return
						}
						bizSetIDs = append(bizSetIDs, bizSetID)
					}
				}
			}
		} else {
			bizSetID, err := util.GetInt64ByInterface(mgoFilter[common.BKAppSetIDField])
			if err != nil {
				ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppSetIDField))
				return
			}
			bizSetIDs = []int64{bizSetID}
		}
	}

	if s.AuthManager.Enabled() {
		authBizSetIDs, err = s.getAuthBizSetIDList(ctx, mgoFilter, searchCond, bizSetIDs)
		if err != nil {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppSetIDField))
			return
		}
	}

	if searchCond.Page.EnableCount {

		query := &metadata.CommonCountFilter{
			Conditions: searchCond.BizSetPropertyFilter,
		}
		result, err := s.Logics.BusinessSetOperation().CountBizSet(ctx.Kit, query)
		if nil != err {
			blog.Errorf("failed to find the objects(%s), error info is %s, rid: %s",
				ctx.Request.PathParameter("obj_id"), err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		ctx.RespEntity(result)
	} else {
		query := &metadata.CommonSearchFilter{
			Conditions: &querybuilder.QueryFilter{
				Rule: &querybuilder.AtomRule{
					Field:    common.BKAppSetIDField,
					Operator: querybuilder.OperatorIn,
					Value:    authBizSetIDs,
				},
			},
			Fields: searchCond.Fields,
			Page:   searchCond.Page,
		}
		result, err := s.Logics.BusinessSetOperation().FindBizSet(ctx.Kit, query)
		if nil != err {
			blog.Errorf("failed to find the objects(%s), error info is %s, rid: %s",
				ctx.Request.PathParameter("obj_id"), err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		ctx.RespEntity(result)
	}

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
