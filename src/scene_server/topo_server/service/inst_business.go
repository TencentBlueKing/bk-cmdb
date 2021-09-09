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
	"configcenter/src/ac"
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
	"configcenter/src/scene_server/topo_server/core/inst"
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

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	data.Set(common.BKDefaultField, common.DefaultFlagDefaultValue)
	// do with transaction
	var business inst.Inst
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		business, err = s.Core.BusinessOperation().CreateBusiness(ctx.Kit, obj, data)
		if err != nil {
			blog.Errorf("create business failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return err
		}

		// register business resource creator action to iam
		if auth.EnableAuthorize() {
			var bizID int64
			if bizID, err = business.GetBizID(); err != nil {
				blog.ErrorJSON("get biz id failed, err: %s, biz: %s, rid: %s", err, business, ctx.Kit.Rid)
				return err
			}
			var bizName string
			if bizName, err = business.GetInstName(); err != nil {
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
	if nil != err {
		blog.Errorf("failed to search the business, %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-business]failed to parse the biz id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Core.BusinessOperation().UpdateBusiness(ctx.Kit, data, obj, bizID)
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
	if nil != err {
		blog.Errorf("failed to search the business, %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-business]failed to parse the biz id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}
	query := &metadata.QueryBusinessRequest{
		Condition: mapstr.MapStr{common.BKAppIDField: bizID},
	}
	_, bizs, err := s.Core.BusinessOperation().FindBiz(ctx.Kit, query)
	if len(bizs) <= 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommNotFound))
		return
	}
	biz := metadata.BizBasicInfo{}
	if err := mapstruct.Decode2Struct(bizs[0], &biz); err != nil {
		blog.Errorf("[api-business]failed, parse biz failed, biz: %+v, err: %s, rid: %s", bizs[0], err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParseDBFailed))
		return
	}

	updateData := mapstr.New()
	switch common.DataStatusFlag(ctx.Request.PathParameter("flag")) {
	case common.DataStatusDisabled:
		if err := s.Core.AssociationOperation().CheckAssociation(ctx.Kit, obj.Object().ObjectID, bizID); nil != err {
			ctx.RespAutoError(err)
			return
		}

		// check if this business still has hosts.
		has, err := s.Core.BusinessOperation().HasHosts(ctx.Kit, bizID)
		if err != nil {
			ctx.RespAutoError(err)
			return
		}
		if has {
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrTopoArchiveBusinessHasHost))
			return
		}
		achieveBizName, err := s.Core.BusinessOperation().GenerateAchieveBusinessName(ctx.Kit, biz.BizName)
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
		err = s.Core.BusinessOperation().UpdateBusiness(ctx.Kit, updateData, obj, bizID)
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

// UpdateBizPropertyBatch batch update business properties
func (s *Service) UpdateBizPropertyBatch(ctx *rest.Contexts) {
	param := new(metadata.UpdateBizPropertyBatchParameter)
	if err := ctx.DecodeInto(&param); nil != err {
		ctx.RespAutoError(err)
		return
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	updateCond := condition.CreateCondition()
	if param.EditAll == true {
		if s.AuthManager.Enabled() {
			bizList, err := s.getAllEditableBiz(ctx)
			if err != nil {
				ctx.RespErrorCodeOnly(common.CCErrorTopoGetAuthorizedBusinessListFailed, "")
				return
			}
			updateCond.Field(common.BKAppIDField).In(bizList)
		}
	} else {
		bizList := param.BizID
		if len(bizList) > common.BKMaxPageSize {
			blog.Errorf("bk_biz_id len %d exceed max page size %d, rid:%s", len(bizList),
				common.BKMaxPageSize, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "update",
				common.BKMaxPageSize))
			return
		}
		// check authorization
		if err := s.checkBizEditable(ctx, bizList); err != nil {
			// resp in checkBizEditable, need not resp again
			return
		}
		updateCond.Field(common.BKAppIDField).In(bizList)
	}

	// exclude archived biz
	updateCond.Field(common.BKDataStatusField).NotEq(common.DataStatusDisabled)
	// exclude default biz
	updateCond.Field(common.BKDefaultField).NotEq(common.DefaultAppFlag)

	data, err := mapstr.NewFromInterface(param.Properties)
	if err != nil {
		blog.Errorf("update biz property batch convert properties[%v] to mapstr failed, err: %v, rid: %s",
			param.Properties, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	// cannot update biz bk_data_status
	data.Remove(common.BKDataStatusField)

	// update biz instances
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.Core.BusinessOperation().UpdateBusinessByCond(ctx.Kit, data, obj, updateCond); err != nil {
			blog.Errorf("update many biz err: %+v, rid: %s", err, ctx.Kit.Rid)
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

func (s *Service) getAllEditableBiz(ctx *rest.Contexts) ([]int64, error) {
	bizList := make([]int64, 0)
	if s.AuthManager.Enabled() == false {
		return bizList, nil
	}

	authInput := meta.ListAuthorizedResourcesParam{
		UserName:     ctx.Kit.User,
		ResourceType: meta.Business,
		Action:       meta.Update,
	}
	authorizedResources, err := s.AuthManager.Authorizer.ListAuthorizedResources(ctx.Kit.Ctx, ctx.Kit.Header, authInput)
	if err != nil {
		blog.Errorf("get editable biz err: %s, user: %s, rid: %s", err.Error(), ctx.Kit.User, ctx.Kit.Rid)
		return bizList, err
	}

	for _, resourceID := range authorizedResources {
		bizID, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			blog.Errorf("parse bizID(%s) failed, err: %v, rid: %s", resourceID, err, ctx.Kit.Rid)
			return bizList, err
		}
		bizList = append(bizList, bizID)
	}

	return bizList, nil
}

func (s *Service) checkBizEditable(ctx *rest.Contexts, bizList []int64) error {
	if s.AuthManager.Enabled() == false {
		return nil
	}

	err := s.AuthManager.AuthorizeByBusinessID(ctx.Kit.Ctx, ctx.Kit.Header, meta.Update, bizList...)
	if err == nil {
		return nil
	}
	if err != ac.NoAuthorizeError {
		blog.ErrorJSON("check biz authorization failed, biz: %+v, err: %s, rid: %s", bizList, err.Error(),
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
		return err
	}

	perm, err := s.AuthManager.GenBizBatchNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, meta.Update, bizList)
	if err != nil && err != ac.NoAuthorizeError {
		blog.ErrorJSON("biz authorization get permission failed, biz: %+v, err: %s, rid: %s", bizList,
			err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
		return err
	}
	ctx.RespEntityWithError(perm, ac.NoAuthorizeError)
	return ac.NoAuthorizeError
}

// find business list with these info：
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
	if errKey, err := page.Validate(true); err != nil {
		blog.Errorf("[api-biz] SearchReducedBusinessList failed, page parameter invalid, errKey: %s, err: %s, rid: %s", errKey, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errKey))
		return
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
		authorizedResources, err := s.AuthManager.Authorizer.ListAuthorizedResources(ctx.Kit.Ctx, ctx.Kit.Header, authInput)
		if err != nil {
			blog.Errorf("[api-biz] SearchReducedBusinessList failed, ListAuthorizedResources failed, user: %s, err: %s, rid: %s", ctx.Kit.User, err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoGetAuthorizedBusinessListFailed))
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

		// sort for prepare to find business with page.
		sort.Sort(util.Int64Slice(appList))
		// user can only find business that is already authorized.
		query.Condition[common.BKAppIDField] = mapstr.MapStr{common.BKDBIN: appList}

	}

	cnt, instItems, err := s.Core.BusinessOperation().FindBiz(ctx.Kit, query)
	if nil != err {
		blog.Errorf("[api-business] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("obj_id"), err.Error(), ctx.Kit.Rid)
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

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", datas)
	ctx.RespEntity(result)
}

func (s *Service) GetBusinessBasicInfo(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-business]failed to parse the biz id, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "business id"))
		return
	}
	query := &metadata.QueryCondition{
		Fields: []string{common.BKAppNameField, common.BKAppIDField},
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
		},
	}
	result, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.Errorf("failed to get business by id, bizID: %s, err: %s, rid: %s", bizID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(result.Data.Info) == 0 {
		blog.Errorf("GetBusinessBasicInfo failed, get business by id not found, bizID: %d, rid: %s", bizID, ctx.Kit.Rid)
		err := ctx.Kit.CCError.CCError(common.CCErrCommNotFound)
		ctx.RespAutoError(err)
		return
	}
	bizData := result.Data.Info[0]
	ctx.RespEntity(bizData)
}

// 4 scenarios, such as user's name user1, scenarios as follows:
// user1
// user1,user3
// user2,user1
// user2,user1,user4
const exactUserRegexp = `(^USER_PLACEHOLDER$)|(^USER_PLACEHOLDER[,]{1})|([,]{1}USER_PLACEHOLDER[,]{1})|([,]{1}USER_PLACEHOLDER$)`

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
					like := strings.Replace(exactUserRegexp, "USER_PLACEHOLDER", gparams.SpecialCharChange(user), -1)
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
	if err := ctx.DecodeInto(&searchCond); nil != err {
		blog.Errorf("failed to parse the params, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	attrCond := condition.CreateCondition()
	attrCond.Field(metadata.AttributeFieldObjectID).Eq(common.BKInnerObjIDApp)
	attrCond.Field(metadata.AttributeFieldPropertyType).Eq(common.FieldTypeUser)
	attrArr, err := s.Core.AttributeOperation().FindBusinessAttribute(ctx.Kit, attrCond.ToMapStr())
	if nil != err {
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
	var bizIDs, authBizIDs []int64
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
		authorizedResources, err := s.AuthManager.Authorizer.ListAuthorizedResources(ctx.Kit.Ctx, ctx.Kit.Header, authInput)
		if err != nil {
			blog.Errorf("[api-biz] SearchBusiness failed, ListAuthorizedResources failed, user: %s, err: %s, rid: %s", ctx.Kit.User, err.Error(), ctx.Kit.Rid)
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

	cnt, instItems, err := s.Core.BusinessOperation().FindBiz(ctx.Kit, searchCond)
	if nil != err {
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

	cnt, instItems, err := s.Core.BusinessOperation().FindBiz(ctx.Kit, &query)
	if nil != err {
		blog.Errorf("[api-business] failed to find the objects(%s), error info is %s, rid: %s", ctx.Request.PathParameter("obj_id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if cnt == 0 {
		blog.InfoJSON("cond:%s, header:%s, rid:%s", query, ctx.Kit.Header, ctx.Kit.Rid)
	}
	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
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
	obj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	data.Set(common.BKDefaultField, common.DefaultAppFlag)

	var business inst.Inst
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		business, err = s.Core.BusinessOperation().CreateBusiness(ctx.Kit, obj, data)
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

func (s *Service) GetInternalModule(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrTopoAppSearchFailed, err.Error()))
		return
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)

	_, result, err := s.Core.BusinessOperation().GetInternalModule(ctx.Kit, bizID)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *Service) GetInternalModuleWithStatistics(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrTopoAppSearchFailed, err.Error()))
		return
	}

	_, innerAppTopo, err := s.Core.BusinessOperation().GetInternalModule(ctx.Kit, bizID)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	if innerAppTopo == nil {
		blog.ErrorJSON("get internal module with statistics failed, type: %s, rid: %s", innerAppTopo, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	moduleIDArr := make([]int64, 0)
	for _, item := range innerAppTopo.Module {
		moduleIDArr = append(moduleIDArr, item.ModuleID)
	}

	// count host apply rules
	listApplyRuleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: moduleIDArr,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	hostApplyRules, err := s.Engine.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx,
		ctx.Kit.Header, bizID, listApplyRuleOption)
	if err != nil {
		blog.Errorf("fillStatistics failed, ListHostApplyRule failed, bizID: %d, option: %+v, err: %+v, rid: %s",
			bizID, listApplyRuleOption, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	moduleRuleCount := make(map[int64]int64)
	for _, item := range hostApplyRules.Info {
		if _, exist := moduleRuleCount[item.ModuleID]; exist == false {
			moduleRuleCount[item.ModuleID] = 0
		}
		moduleRuleCount[item.ModuleID] += 1
	}

	set := mapstr.NewFromStruct(innerAppTopo, "field")
	modules := make([]mapstr.MapStr, 0)
	for _, module := range innerAppTopo.Module {
		moduleItem := mapstr.NewFromStruct(module, "field")
		moduleItem["host_apply_rule_count"] = 0
		if ruleCount, ok := moduleRuleCount[module.ModuleID]; ok == true {
			moduleItem["host_apply_rule_count"] = ruleCount
		}
		modules = append(modules, moduleItem)
	}
	set["module"] = modules
	ctx.RespEntity(set)
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
	if errKey, err := page.Validate(true); err != nil {
		blog.Errorf("[api-biz] ListAllBusinessSimplify failed, page parameter invalid, errKey: %s, err: %s, rid: %s", errKey, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errKey))
		return
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
	cnt, instItems, err := s.Core.BusinessOperation().FindBiz(ctx.Kit, query)
	if nil != err {
		blog.Errorf("ListAllBusinessSimplify failed, FindBusiness failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	businesses := make([]metadata.BizBasicInfo, 0)
	for _, item := range instItems {
		business := metadata.BizBasicInfo{}
		if err := mapstruct.Decode2Struct(item, &business); err != nil {
			blog.Errorf("ListAllBusinessSimplify failed, decode biz from db failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParseDBFailed))
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
