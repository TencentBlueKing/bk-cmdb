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
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	gparams "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/thirdparty/hooks"
)

// CreateBusiness create a new business
func (s *Service) CreateBusiness(ctx *rest.Contexts) {
	data := mapstr.MapStr{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := hooks.ValidateCreateBusinessHook(ctx.Kit, s.Engine.CoreAPI, data); err != nil {
		blog.Errorf("validate create business hook failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	data.Set(common.BKDefaultField, common.DefaultFlagDefaultValue)
	// do with transaction
	var business mapstr.MapStr
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		business, err = s.Logics.BusinessOperation().CreateBusiness(ctx.Kit, data)
		if err != nil {
			blog.Errorf("create business failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}

		// register business resource creator action to iam
		if auth.EnableAuthorize() {
			var bizID int64
			if bizID, err = business.Int64(common.BKAppNameField); err != nil {
				blog.ErrorJSON("get biz name failed, err: %s, biz: %s, rid: %s", err, business, ctx.Kit.Rid)
				return err
			}
			var bizName string
			if bizName, err = business.String(common.BKAppNameField); err != nil {
				blog.ErrorJSON("get biz name failed, err: %s, biz: %s, rid: %s", err, business, ctx.Kit.Rid)
				return err
			}
			iamInstance := metadata.IamInstanceWithCreator{
				Type:    string(iam.Business),
				ID:      strconv.FormatInt(bizID, 10),
				Name:    bizName,
				Creator: ctx.Kit.User,
			}
			_, err = s.AuthManager.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance)
			if err != nil {
				blog.Errorf("register created business to iam failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(business)
}

// UpdateBusiness update the business
func (s *Service) UpdateBusiness(ctx *rest.Contexts) {
	data := mapstr.MapStr{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDApp)
	if err != nil {
		blog.Errorf("failed to search the business, %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Logics.BusinessOperation().UpdateBusiness(ctx.Kit, data, obj.Object(), bizID)
		if err != nil {
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

// UpdateBusinessStatus update the business status
func (s *Service) UpdateBusinessStatus(ctx *rest.Contexts) {
	data := struct {
		metadata.UpdateBusinessStatusOption `json:",inline"`
	}{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDApp)
	if err != nil {
		blog.Errorf("search business failed, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	query := &metadata.QueryBusinessRequest{
		Condition: mapstr.MapStr{common.BKAppIDField: bizID},
	}
	_, bizs, err := s.Logics.BusinessOperation().FindBiz(ctx.Kit, query, false)
	if len(bizs) <= 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommNotFound))
		return
	}
	biz := metadata.BizBasicInfo{}
	if err := mapstruct.Decode2Struct(bizs[0], &biz); err != nil {
		blog.Errorf("parse biz failed, biz: %+v, err: %v, rid: %s", bizs[0], err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	updateData := mapstr.MapStr{}
	switch common.DataStatusFlag(ctx.Request.PathParameter("flag")) {
	case common.DataStatusDisabled:
		if err := s.Core.AssociationOperation().CheckAssociation(ctx.Kit, obj.Object().ObjectID, bizID); err != nil {
			ctx.RespAutoError(err)
			return
		}

		// check if this business still has hosts.
		has, err := s.Logics.BusinessOperation().HasHosts(ctx.Kit, bizID)
		if err != nil {
			ctx.RespAutoError(err)
			return
		}
		if has {
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrTopoArchiveBusinessHasHost))
			return
		}
		achieveBizName, err := s.Logics.BusinessOperation().GenerateAchieveBusinessName(ctx.Kit, biz.BizName)
		if err != nil {
			ctx.RespAutoError(err)
			return
		}
		updateData.Set(common.BKAppNameField, achieveBizName)
		updateData.Set(common.BKDataStatusField, ctx.Request.PathParameter("flag"))
	case common.DataStatusEnable:
		if len(data.UpdateBusinessStatusOption.BizName) > 0 {
			updateData.Set(common.BKAppNameField, data.UpdateBusinessStatusOption.BizName)
		}
		updateData.Set(common.BKDataStatusField, ctx.Request.PathParameter("flag"))
	default:
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, ctx.Request.PathParameter))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Logics.BusinessOperation().UpdateBusiness(ctx.Kit, updateData, obj.Object(), bizID)
		if err != nil {
			blog.Errorf("UpdateBusinessStatus failed, run update failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
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

//SearchReducedBusinessList find business list with these info：
// 1. have any authorized resources in a business.
// 2. only returned with a few field for this business info.
func (s *Service) SearchReducedBusinessList(ctx *rest.Contexts) {
	page := metadata.BasePage{
		Limit: common.BKNoLimit,
	}
	sortParam := ctx.Request.QueryParameter("sort")
	if len(sortParam) > 0 {
		page.Sort = sortParam
	}
	query := &metadata.QueryBusinessRequest{
		Fields: []string{common.BKAppIDField, common.BKAppNameField},
		Page:   page,
		Condition: mapstr.MapStr{
			common.BKDataStatusField: mapstr.MapStr{common.BKDBNE: common.DataStatusDisabled},
			common.BKDefaultField:    0,
		},
	}

	if s.AuthManager.Enabled() {
		authInput := meta.ListAuthorizedResourcesParam{
			UserName:     ctx.Kit.User,
			ResourceType: meta.Business,
			Action:       meta.ViewBusinessResource,
		}
		authorizedResources, err := s.AuthManager.Authorizer.ListAuthorizedResources(ctx.Kit.Ctx, ctx.Kit.Header,
			authInput)
		if err != nil {
			blog.Errorf("ListAuthorizedResources failed, user: %s, err: %v, rid: %s", ctx.Kit.User, err,
				ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		appList := make([]int64, 0)
		for _, resourceID := range authorizedResources {
			bizID, err := strconv.ParseInt(resourceID, 10, 64)
			if err != nil {
				blog.Errorf("parse bizID(%s) failed, err: %v, rid: %s", bizID, err, ctx.Kit.Rid)
				ctx.RespAutoError(err)
				return
			}
			appList = append(appList, bizID)
		}

		// sort for prepare to find business with page.
		sort.Sort(util.Int64Slice(appList))
		// user can only find business that is already authorized.
		query.Condition[common.BKAppIDField] = mapstr.MapStr{common.BKDBIN: appList}
	}

	cnt, instItems, err := s.Logics.BusinessOperation().FindBiz(ctx.Kit, query, false)
	if err != nil {
		blog.Errorf(" find objects failed, err: %v, rid: %s", ctx.Request.PathParameter("obj_id"), err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	datas := make([]mapstr.MapStr, 0)
	for _, item := range instItems {
		inst := mapstr.New()
		inst[common.BKAppIDField] = item[common.BKAppIDField]
		inst[common.BKAppNameField] = item[common.BKAppNameField]
		datas = append(datas, inst)
	}

	result := mapstr.MapStr{
		"count": cnt,
		"info":  datas,
	}
	ctx.RespEntity(result)
}

func (s *Service) GetBusinessBasicInfo(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	query := &metadata.QueryCondition{
		Fields: []string{common.BKAppNameField, common.BKAppIDField},
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
		},
	}
	result, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDApp, query)
	if err != nil {
		blog.Errorf("get business failed, bizID: %s, err: %v, rid: %s", bizID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(result.Info) == 0 {
		blog.Errorf("get business by id not found, bizID: %d, rid: %s", bizID, ctx.Kit.Rid)
		err := ctx.Kit.CCError.CCError(common.CCErrCommNotFound)
		ctx.RespAutoError(err)
		return
	}
	bizData := result.Info[0]
	ctx.RespEntity(bizData)
}

// 4 scenarios, such as user's name user1, scenarios as follows:
// user1
// user1,user3
// user2,user1
// user2,user1,user4
const exactUserRegexp = `(^USER_PLACEHOLDER$)|(^USER_PLACEHOLDER[,]{1})|([,]{1}USER_PLACEHOLDER[,]{1})|([,
]{1}USER_PLACEHOLDER$)`

func handleSpecialBusinessFieldSearchCond(input map[string]interface{}, userFieldArr []string) map[string]interface{} {
	output := make(map[string]interface{})
	exactAnd := make([]map[string]interface{}, 0)
	for i, j := range input {
		if j == nil {
			output[i] = j
			continue
		}

		objType := reflect.TypeOf(j)
		switch objType.Kind() {
		case reflect.String:
			if _, ok := j.(json.Number); ok {
				output[i] = j
				continue
			}
			targetStr := j.(string)
			if util.InStrArr(userFieldArr, i) {
				for _, user := range strings.Split(strings.Trim(targetStr, ","), ",") {
					// search with exactly the user's name with regexpF
					like := strings.Replace(exactUserRegexp, "USER_PLACEHOLDER", gparams.SpecialCharChange(user),
						-1)
					exactAnd = append(exactAnd, mapstr.MapStr{i: mapstr.MapStr{common.BKDBLIKE: like}})
				}
			} else {
				attrVal := gparams.SpecialCharChange(targetStr)
				output[i] = map[string]interface{}{common.BKDBLIKE: attrVal, common.BKDBOPTIONS: "i"}
			}
		default:
			output[i] = j
		}
	}

	if len(exactAnd) > 0 {
		output[common.BKDBAND] = exactAnd
	}

	return output
}

// SearchBusiness search the business by condition
// func (s *Service) SearchBusiness(ctx *rest.Contexts) {
func (s *Service) SearchBusiness(ctx *rest.Contexts) {
	searchCond := new(metadata.QueryBusinessRequest)
	if err := ctx.DecodeInto(&searchCond); err != nil {
		blog.Errorf("failed to parse the params, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	attrCond := condition.CreateCondition()
	attrCond.Field(metadata.AttributeFieldObjectID).Eq(common.BKInnerObjIDApp)
	attrCond.Field(metadata.AttributeFieldPropertyType).Eq(common.FieldTypeUser)
	attrArr, err := s.Core.AttributeOperation().FindBusinessAttribute(ctx.Kit, attrCond.ToMapStr())
	if err != nil {
		blog.Errorf("failed get the business attribute, %s, rid:%s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	// userFieldArr Fields in the business are user-type fields
	var userFields []string
	for _, attribute := range attrArr {
		userFields = append(userFields, attribute.PropertyID)
	}

	searchCond.Condition = handleSpecialBusinessFieldSearchCond(searchCond.Condition, userFields)

	// parse business id from user's condition for testing.
	bizIDs := make([]int64, 0)
	authBizIDs := make([]int64, 0)
	biz, exist := searchCond.Condition[common.BKAppIDField]
	if exist {
		// constrict that bk_biz_id field can only be a numeric value,
		// operators like or/in/and is not allowed.
		if bizcond, ok := biz.(map[string]interface{}); ok {
			if cond, ok := bizcond["$eq"]; ok {
				bizID, err := util.GetInt64ByInterface(cond)
				if err != nil {
					ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "", common.BKAppIDField)
					return
				}
				bizIDs = []int64{bizID}
			}
			if cond, ok := bizcond["$in"]; ok {
				if conds, ok := cond.([]interface{}); ok {
					for _, c := range conds {
						bizID, err := util.GetInt64ByInterface(c)
						if err != nil {
							ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "", common.BKAppIDField)
							return
						}
						bizIDs = append(bizIDs, bizID)
					}
				}
			}
		} else {
			bizID, err := util.GetInt64ByInterface(searchCond.Condition[common.BKAppIDField])
			if err != nil {
				ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
				return
			}
			bizIDs = []int64{bizID}
		}
	}

	if s.AuthManager.Enabled() {
		authInput := meta.ListAuthorizedResourcesParam{
			UserName:     ctx.Kit.User,
			ResourceType: meta.Business,
			Action:       meta.Find,
		}
		authorizedResources, err := s.AuthManager.Authorizer.ListAuthorizedResources(ctx.Kit.Ctx, ctx.Kit.Header,
			authInput)
		if err != nil {
			blog.Errorf("SearchBusiness failed, ListAuthorizedResources failed, user: %s, err: %v, rid: %s",
				ctx.Kit.User, err, ctx.Kit.Rid)
			ctx.RespErrorCodeOnly(common.CCErrorTopoGetAuthorizedBusinessListFailed, "")
			return
		}
		appList := make([]int64, 0)
		for _, resourceID := range authorizedResources {
			bizID, err := strconv.ParseInt(resourceID, 10, 64)
			if err != nil {
				blog.Errorf("parse bizID(%s) failed, err: %v, rid: %s", bizID, err, ctx.Kit.Rid)
				ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
				return
			}
			appList = append(appList, bizID)
		}
		if len(bizIDs) > 0 {
			// this means that user want to find a specific business.
			// now we check if he has this authority.
			for _, bizID := range bizIDs {
				if util.InArray(bizID, appList) {
					// authBizIDs store the authorized bizIDs
					authBizIDs = append(authBizIDs, bizID)
				}
			}
			if len(authBizIDs) > 0 {
				searchCond.Condition[common.BKAppIDField] = mapstr.MapStr{common.BKDBIN: authBizIDs}
			} else {
				// if there are no qualified bizIDs, return null
				result := mapstr.MapStr{}
				result.Set("count", 0)
				result.Set("info", []mapstr.MapStr{})
				ctx.RespEntity(result)
				return
			}
			// now you have the authority.
		} else {
			// sort for prepare to find business with page.
			sort.Sort(util.Int64Slice(appList))
			// user can only find business that is already authorized.
			searchCond.Condition[common.BKAppIDField] = mapstr.MapStr{common.BKDBIN: appList}
		}
	}

	if _, ok := searchCond.Condition[common.BKDataStatusField]; !ok {
		searchCond.Condition[common.BKDataStatusField] = mapstr.MapStr{common.BKDBNE: common.DataStatusDisabled}
	}

	// can only find normal business, but not resource pool business
	searchCond.Condition[common.BKDefaultField] = 0

	cnt, instItems, err := s.Logics.BusinessOperation().FindBiz(ctx.Kit, searchCond, false)
	if err != nil {
		blog.Errorf("find business failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)

	ctx.RespEntity(result)
}

// SearchOwnerResourcePoolBusiness search archived business by condition
func (s *Service) SearchOwnerResourcePoolBusiness(ctx *rest.Contexts) {

	supplierAccount := ctx.Request.PathParameter("owner_id")
	query := metadata.QueryBusinessRequest{
		Condition: mapstr.MapStr{
			common.BKDefaultField:    common.DefaultAppFlag,
			common.BkSupplierAccount: supplierAccount,
		},
	}

	cnt, instItems, err := s.Logics.BusinessOperation().FindBiz(ctx.Kit, &query, false)
	if err != nil {
		blog.Errorf("find objects(%s) failed, err: %v, rid: %s", ctx.Request.PathParameter("obj_id"), err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if cnt == 0 {
		blog.InfoJSON("cond:%s, header:%s, rid:%s", query, ctx.Kit.Header, ctx.Kit.Rid)
	}
	result := mapstr.MapStr{
		"count": cnt,
		"info":  instItems,
	}
	ctx.RespEntity(result)
	return
}

// CreateDefaultBusiness create the default business
func (s *Service) CreateDefaultBusiness(ctx *rest.Contexts) {
	data := mapstr.MapStr{}
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	data.Set(common.BKDefaultField, common.DefaultAppFlag)

	var business mapstr.MapStr
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		business, err = s.Logics.BusinessOperation().CreateBusiness(ctx.Kit, data)
		if err != nil {
			blog.Errorf("create business failed, err: %+v", err)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(business)
}

// ListAllBusinessSimplify list all businesses with return only several fields
func (s *Service) ListAllBusinessSimplify(ctx *rest.Contexts) {
	page := metadata.BasePage{
		Limit: common.BKNoLimit,
	}
	sortParam := ctx.Request.QueryParameter("sort")
	if len(sortParam) > 0 {
		page.Sort = sortParam
	}

	fields := []string{
		common.BKAppIDField,
		common.BKAppNameField,
	}

	query := &metadata.QueryBusinessRequest{
		Fields: fields,
		Page:   page,
		Condition: mapstr.MapStr{
			common.BKDataStatusField: mapstr.MapStr{common.BKDBNE: common.DataStatusDisabled},
		},
	}
	cnt, instItems, err := s.Logics.BusinessOperation().FindBiz(ctx.Kit, query, false)
	if err != nil {
		blog.Errorf("find business failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	businesses := make([]metadata.BizBasicInfo, 0)
	for _, item := range instItems {
		business := metadata.BizBasicInfo{}
		if err := mapstruct.Decode2Struct(item, &business); err != nil {
			blog.Errorf("decode biz from db failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		businesses = append(businesses, business)
	}

	result := map[string]interface{}{
		"count": cnt,
		"info":  businesses,
	}
	ctx.RespEntity(result)
}

// GetBriefTopologyNodeRelation is used to get directly related business topology node information.
// As is, you can find modules belongs to a set; or you can find the set a module belongs to.
func (s *Service) GetBriefTopologyNodeRelation(ctx *rest.Contexts) {
	options := new(metadata.GetBriefBizRelationOptions)
	if err := ctx.DecodeInto(options); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := options.Validate()
	if rawErr.ErrCode != 0 {
		blog.Errorf("validate failed, err: %v, rid: %s", rawErr.Args, ctx.Kit.Rid)
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	relations, err := s.Logics.BusinessOperation().GetBriefTopologyNodeRelation(ctx.Kit, options)
	if err != nil {
		blog.Errorf("get brief topology node relation failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(&relations)
}
