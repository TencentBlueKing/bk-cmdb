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
	"reflect"
	"sort"
	"strconv"
	"strings"

	"configcenter/src/auth"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	gparams "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateBusiness create a new business
func (s *Service) CreateBusiness(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	var txnErr error
	// 判断是否使用事务
	if s.EnableTxn {
		sess, err := s.DB.StartSession()
		if err != nil {
			txnErr = err
			blog.Errorf("StartSession err: %s, rid: %s", err.Error(), params.ReqID)
			return nil, err
		}
		// 获取事务信息，将其存入context中
		txnInfo, err := sess.TxnInfo()
		if err != nil {
			txnErr = err
			blog.Errorf("TxnInfo err: %+v", err)
			return nil, params.Err.Error(common.CCErrObjectDBOpErrno)
		}
		params.Header = txnInfo.IntoHeader(params.Header)
		params.Context = util.TnxIntoContext(params.Context, txnInfo)
		err = sess.StartTransaction(params.Context)
		if err != nil {
			txnErr = err
			blog.Errorf("StartTransaction err: %+v", err)
			return nil, params.Err.Error(common.CCErrObjectDBOpErrno)
		}
		defer func() {
			if txnErr == nil {
				err = sess.CommitTransaction(params.Context)
				if err != nil {
					blog.Errorf("CommitTransaction err: %+v", err)
				}
			} else {
				blog.Errorf("Occur err:%v, begin AbortTransaction", txnErr)
				err = sess.AbortTransaction(params.Context)
				if err != nil {
					blog.Errorf("AbortTransaction err: %+v", err)
				}
			}
			sess.EndSession(params.Context)
		}()
	}

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		txnErr = err
		blog.Errorf("failed to search the business, %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	data.Set(common.BKDefaultField, common.DefaultFlagDefaultValue)
	business, err := s.Core.BusinessOperation().CreateBusiness(params, obj, data)
	if err != nil {
		txnErr = err
		blog.Errorf("create business failed, err: %v, rid: %s", err, params.ReqID)
		return nil, err
	}

	businessID, err := business.GetInstID()
	if err != nil {
		txnErr = err
		blog.Errorf("unexpected error, create business success, but get id failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommParamsInvalid)
	}

	// auth: register business to iam
	if err := s.AuthManager.RegisterBusinessesByID(params.Context, params.Header, businessID); err != nil {
		txnErr = err
		blog.Errorf("create business success, but register to iam failed, err: %v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}
	return business, nil
}

// DeleteBusiness delete the business
func (s *Service) DeleteBusiness(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-business]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	// auth: deregister business to iam
	if err := s.AuthManager.DeregisterBusinessesByID(params.Context, params.Header, bizID); err != nil {
		blog.Errorf("delete business failed, deregister business failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommUnRegistResourceToIAMFailed)
	}

	return nil, s.Core.BusinessOperation().DeleteBusiness(params, obj, bizID)
}

// UpdateBusiness update the business
func (s *Service) UpdateBusiness(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-business]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}

	err = s.Core.BusinessOperation().UpdateBusiness(params, data, obj, bizID)
	if err != nil {
		return nil, err
	}

	// auth: update registered business to iam
	if err := s.AuthManager.UpdateRegisteredBusinessByID(params.Context, params.Header, bizID); err != nil {
		blog.Errorf("update business success, but update registered business failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommRegistResourceToIAMFailed)
	}

	return nil, nil
}

// UpdateBusinessStatus update the business status
func (s *Service) UpdateBusinessStatus(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-business]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}
	data = mapstr.New()
	query := &metadata.QueryBusinessRequest{
		Condition: mapstr.MapStr{common.BKAppIDField: bizID},
	}
	_, bizs, err := s.Core.BusinessOperation().FindBusiness(params, query)
	if len(bizs) <= 0 {
		return nil, params.Err.Error(common.CCErrCommNotFound)
	}
	data = mapstr.New()
	switch common.DataStatusFlag(pathParams("flag")) {
	case common.DataStatusDisabled:
		innerCond := condition.CreateCondition()
		innerCond.Field(common.BKAsstObjIDField).Eq(obj.Object().ObjectID)
		innerCond.Field(common.BKAsstInstIDField).Eq(bizID)
		if err := s.Core.AssociationOperation().CheckBeAssociation(params, obj, innerCond); nil != err {
			return nil, err
		}

		// check if this business still has hosts.
		has, err := s.Core.BusinessOperation().HasHosts(params, bizID)
		if err != nil {
			return nil, err
		}
		if has {
			return nil, params.Err.Error(common.CCErrTopoArchiveBusinessHasHost)
		}

		data.Set(common.BKDataStatusField, pathParams("flag"))
	case common.DataStatusEnable:
		name, err := bizs[0].String(common.BKAppNameField)
		if nil != err {
			return nil, params.Err.Error(common.CCErrCommNotFound)
		}
		name = name + common.BKDataRecoverSuffix
		if len(name) >= common.FieldTypeSingleLenChar {
			name = name[:common.FieldTypeSingleLenChar]
		}
		data.Set(common.BKAppNameField, name)
		data.Set(common.BKDataStatusField, pathParams("flag"))
	default:
		return nil, params.Err.Errorf(common.CCErrCommParamsIsInvalid, pathParams("flag"))
	}

	err = s.Core.BusinessOperation().UpdateBusiness(params, data, obj, bizID)
	if err != nil {
		blog.Errorf("UpdateBusinessStatus failed, run update failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	if err := s.AuthManager.UpdateRegisteredBusinessByID(params.Context, params.Header, bizID); err != nil {
		blog.Errorf("UpdateBusinessStatus failed, update register business info failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}
	return nil, nil
}

// find business list with these info：
// 1. have any authorized resources in a business.
// 2. only returned with a few field for this business info.
func (s *Service) SearchReducedBusinessList(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	query := &metadata.QueryBusinessRequest{
		Fields: []string{common.BKAppIDField, common.BKAppNameField, "business_dept_id", "business_dept_name"},
		Page:   metadata.BasePage{},
		Condition: mapstr.MapStr{
			common.BKDataStatusField: mapstr.MapStr{common.BKDBNE: common.DataStatusDisabled},
			common.BKDefaultField:    0,
		},
	}

	if s.AuthManager.Enabled() {
		user := authmeta.UserInfo{UserName: params.User, SupplierAccount: params.SupplierAccount}
		appList, err := s.AuthManager.Authorize.GetAnyAuthorizedBusinessList(params.Context, user)
		if err != nil {
			blog.Errorf("[api-biz] SearchReducedBusinessList failed, GetExactAuthorizedBusinessList failed, user: %s, err: %s, rid: %s", user, err.Error(), params.ReqID)
			return nil, params.Err.Error(common.CCErrorTopoGetAuthorizedBusinessListFailed)
		}

		// sort for prepare to find business with page.
		sort.Sort(util.Int64Slice(appList))
		// user can only find business that is already authorized.
		query.Condition[common.BKAppIDField] = mapstr.MapStr{common.BKDBIN: appList}

	}

	cnt, instItems, err := s.Core.BusinessOperation().FindBusiness(params, query)
	if nil != err {
		blog.Errorf("[api-business] failed to find the objects(%s), error info is %s, rid: %s", pathParams("obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	datas := make([]mapstr.MapStr, 0)
	for _, item := range instItems {
		inst := mapstr.New()
		inst[common.BKAppIDField] = item[common.BKAppIDField]
		inst[common.BKAppNameField] = item[common.BKAppNameField]
		inst["business_dept_id"] = item["business_dept_id"]
		inst["business_dept_name"] = item["business_dept_name"]

		if val, exist := item["business_dept_id"]; exist {
			inst["business_dept_id"] = val
		} else {
			inst["business_dept_id"] = ""
		}
		if val, exist := item["business_dept_name"]; exist {
			inst["business_dept_name"] = val
		} else {
			inst["business_dept_name"] = ""
		}
		datas = append(datas, inst)
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", datas)
	return result, nil
}

func (s *Service) GetBusinessBasicInfo(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-business]failed to parse the biz id, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "business id")
	}
	query := &metadata.QueryCondition{
		Fields: []string{common.BKAppNameField, common.BKAppIDField},
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
		},
	}
	result, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(params.Context, params.Header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.Errorf("failed to get business by id, bizID: %s, err: %s, rid: %s", bizID, err.Error(), params.ReqID)
		return nil, err
	}
	if len(result.Data.Info) == 0 {
		blog.Errorf("GetBusinessBasicInfo failed, get business by id not found, bizID: %d, rid: %s", bizID, params.ReqID)
		err := params.Err.CCError(common.CCErrCommNotFound)
		return nil, err
	}
	bizData := result.Data.Info[0]
	return bizData, nil
}

// 4 scenarios, such as user's name user1, scenarios as follows:
// user1
// user1,user3
// user2,user1
// user2,user1,user4
const exactUserRegexp = `(^USER_PLACEHOLDER$)|(^USER_PLACEHOLDER[,]{1})|([,]{1}USER_PLACEHOLDER[,]{1})|([,]{1}USER_PLACEHOLDER$)`

func handleSpecialBusinessFieldSearchCond(input map[string]interface{}, userFieldArr []string) map[string]interface{} {
	output := make(map[string]interface{})
	for i, j := range input {
		objType := reflect.TypeOf(j)
		switch objType.Kind() {
		case reflect.String:
			targetStr := j.(string)
			if util.InStrArr(userFieldArr, i) {
				exactOr := make([]map[string]interface{}, 0)
				for _, user := range strings.Split(strings.Trim(targetStr, ","), ",") {
					// search with exactly the user's name with regexp
					like := strings.Replace(exactUserRegexp, "USER_PLACEHOLDER", gparams.SpecialCharChange(user), -1)
					exactOr = append(exactOr, mapstr.MapStr{i: mapstr.MapStr{common.BKDBLIKE: like}})
				}
				output[common.BKDBOR] = exactOr
			} else {
				output[i] = targetStr
			}
		default:
			output[i] = j
		}
	}
	return output
}

// SearchBusiness search the business by condition
func (s *Service) SearchBusiness(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	searchCond := new(metadata.QueryBusinessRequest)
	if err := data.MarshalJSONInto(&searchCond); nil != err {
		blog.Errorf("failed to parse the params, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}

	attrCond := condition.CreateCondition()
	attrCond.Field(metadata.AttributeFieldSupplierAccount).Eq(params.SupplierAccount)
	attrCond.Field(metadata.AttributeFieldObjectID).Eq(common.BKInnerObjIDApp)
	attrCond.Field(metadata.AttributeFieldPropertyType).Eq(common.FieldTypeUser)
	attrArr, err := s.Core.AttributeOperation().FindObjectAttribute(params, attrCond)
	if nil != err {
		blog.Errorf("failed get the business attribute, %s, rid:%s", err.Error(), util.GetHTTPCCRequestID(params.Header))
		return nil, err
	}
	// userFieldArr Fields in the business are user-type fields
	var userFields []string
	for _, attrInterface := range attrArr {
		userFields = append(userFields, attrInterface.Attribute().PropertyID)
	}

	searchCond.Condition = handleSpecialBusinessFieldSearchCond(searchCond.Condition, userFields)

	// parse business id from user's condition for testing.
	var bizIDs []int64
	biz, exist := searchCond.Condition[common.BKAppIDField]
	if exist {
		// constrict that bk_biz_id field can only be a numeric value,
		// operators like or/in/and is not allowed.
		if bizcond, ok := biz.(map[string]interface{}); ok {
			if cond, ok := bizcond["$eq"]; ok {
				if reflect.TypeOf(cond).ConvertibleTo(reflect.TypeOf(int64(1))) == false {
					return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
				}
				bizIDs = []int64{int64(cond.(float64))}
			}
			if cond, ok := bizcond["$in"]; ok {
				if conds, ok := cond.([]interface{}); ok {
					for _, c := range conds {
						if reflect.TypeOf(c).ConvertibleTo(reflect.TypeOf(int64(1))) == false {
							return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
						}
						bizIDs = append(bizIDs, int64(c.(float64)))
					}
				}
			}
		} else if reflect.TypeOf(biz).ConvertibleTo(reflect.TypeOf(int64(1))) {
			bizIDs = []int64{int64(searchCond.Condition[common.BKAppIDField].(float64))}
		} else {
			return nil, params.Err.New(common.CCErrCommParamsInvalid, common.BKAppIDField)
		}
	}

	if s.AuthManager.Enabled() {
		user := authmeta.UserInfo{UserName: params.User, SupplierAccount: params.SupplierAccount}
		appList, err := s.AuthManager.Authorize.GetExactAuthorizedBusinessList(params.Context, user)
		if err != nil {
			blog.Errorf("[api-biz] SearchBusiness failed, GetExactAuthorizedBusinessList failed, user: %s, err: %s, rid: %s", user, err.Error(), params.ReqID)
			return nil, params.Err.Error(common.CCErrorTopoGetAuthorizedBusinessListFailed)
		}

		if len(bizIDs) > 0 {
			// this means that user want to find a specific business.
			// now we check if he has this authority.
			for _, bizID := range bizIDs {
				if !util.InArray(bizID, appList) {
					noAuthResp, err := s.AuthManager.GenBusinessAuditNoPermissionResp(params.Context, params.Header, bizID)
					if err != nil {
						return nil, params.Err.Error(common.CCErrTopoAppSearchFailed)
					}
					return noAuthResp, auth.NoAuthorizeError
				}
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

	cnt, instItems, err := s.Core.BusinessOperation().FindBusiness(params, searchCond)
	if nil != err {
		blog.Errorf("[api-business] failed to find the objects(%s), error info is %s, rid: %s", pathParams("obj_id"), err.Error(), params.ReqID)
		return nil, err
	}

	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)

	return result, nil
}

// search archived business by condition
func (s *Service) SearchArchivedBusiness(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	supplierAccount := pathParams("owner_id")
	query := metadata.QueryBusinessRequest{
		Condition: mapstr.MapStr{
			common.BKDefaultField:    common.DefaultAppFlag,
			common.BkSupplierAccount: supplierAccount,
		},
	}

	if s.AuthManager.Enabled() {
		user := authmeta.UserInfo{UserName: params.User, SupplierAccount: params.SupplierAccount}
		appList, err := s.AuthManager.Authorize.GetExactAuthorizedBusinessList(params.Context, user)
		if err != nil {
			blog.Errorf("[api-biz] SearchArchivedBusiness failed, GetExactAuthorizedBusinessList failed, user: %s, err: %s, rid: %s", user, err.Error(), params.ReqID)
			return nil, params.Err.Error(common.CCErrorTopoGetAuthorizedBusinessListFailed)
		}
		// sort for prepare to find business with page.
		sort.Sort(util.Int64Slice(appList))
		// user can only find business that is already authorized.
		query.Condition[common.BKAppIDField] = mapstr.MapStr{common.BKDBIN: appList}
	}

	cnt, instItems, err := s.Core.BusinessOperation().FindBusiness(params, &query)
	if nil != err {
		blog.Errorf("[api-business] failed to find the objects(%s), error info is %s, rid: %s", pathParams("obj_id"), err.Error(), params.ReqID)
		return nil, err
	}
	result := mapstr.MapStr{}
	result.Set("count", cnt)
	result.Set("info", instItems)
	return result, nil
}

// CreateDefaultBusiness create the default business
func (s *Service) CreateDefaultBusiness(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	data.Set(common.BKDefaultField, common.DefaultAppFlag)
	business, err := s.Core.BusinessOperation().CreateBusiness(params, obj, data)
	if err != nil {
		return nil, fmt.Errorf("create business failed, err: %+v", err)
	}

	businessID, err := business.GetInstID()
	if err != nil {
		return nil, fmt.Errorf("unexpected error, create default business success, but get id failed, err: %+v", err)
	}

	// auth: register business to iam
	if err := s.AuthManager.RegisterBusinessesByID(params.Context, params.Header, businessID); err != nil {
		blog.Errorf("create default business failed, register business failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}

	return business, nil
}

func (s *Service) GetInternalModule(params types.ContextParams, pathParams, queryparams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	obj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("failed to search the business, %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}
	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		return nil, params.Err.New(common.CCErrTopoAppSearchFailed, err.Error())
	}

	_, result, err := s.Core.BusinessOperation().GetInternalModule(params, obj, bizID)
	if nil != err {
		return nil, err
	}

	return result, nil
}

func (s *Service) GetInternalModuleWithStatistics(params types.ContextParams, pathParams, queryparams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		return nil, params.Err.New(common.CCErrTopoAppSearchFailed, err.Error())
	}

	result, err := s.GetInternalModule(params, pathParams, queryparams, data)
	if err != nil {
		return result, err
	}
	innerAppTopo, ok := result.(*metadata.InnterAppTopo)
	if ok == false || innerAppTopo == nil {
		blog.ErrorJSON("GetInternalModuleWithStatistics failed, GetInternalModule return unexpected type: %s, rid: %s", innerAppTopo, params.ReqID)
		return result, err
	}
	moduleIDArr := make([]int64, 0)
	for _, item := range innerAppTopo.Module {
		moduleIDArr = append(moduleIDArr, item.ModuleID)
	}
	listHostOption := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		SetIDArr:      []int64{innerAppTopo.SetID},
		ModuleIDArr:   moduleIDArr,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	hostModuleRelations, e := s.Engine.CoreAPI.CoreService().Host().GetHostModuleRelation(params.Context, params.Header, listHostOption)
	if e != nil {
		blog.Errorf("GetInternalModuleWithStatistics failed, list host modules failed, option: %+v, err: %s, rid: %s", listHostOption, e.Error(), params.ReqID)
		return nil, e
	}
	setHostIDs := make([]int64, 0)
	moduleHostIDs := make(map[int64][]int64, 0)
	for _, relation := range hostModuleRelations.Data.Info {
		setHostIDs = append(setHostIDs, relation.HostID)
		if _, ok := moduleHostIDs[relation.ModuleID]; ok == false {
			moduleHostIDs[relation.ModuleID] = make([]int64, 0)
		}
		moduleHostIDs[relation.ModuleID] = append(moduleHostIDs[relation.ModuleID], relation.HostID)
	}
	set := mapstr.NewFromStruct(innerAppTopo, "field")
	if err != nil {
		blog.Errorf("GetInternalModuleWithStatistics failed, convert innerAppTopo to map failed, innerAppTopo: %+v, err: %s, rid: %s", innerAppTopo, e.Error(), params.ReqID)
		return nil, e
	}
	set["host_count"] = len(util.IntArrayUnique(setHostIDs))
	modules := make([]mapstr.MapStr, 0)
	for _, module := range innerAppTopo.Module {
		moduleItem := mapstr.NewFromStruct(module, "field")
		moduleItem["host_count"] = 0
		if hostIDs, ok := moduleHostIDs[module.ModuleID]; ok == true {
			moduleItem["host_count"] = len(util.IntArrayUnique(hostIDs))
		}
		modules = append(modules, moduleItem)
	}
	set["module"] = modules
	return set, nil
}

// ListAllBusinessSimplify list all businesses with return only several fields
func (s *Service) ListAllBusinessSimplify(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	fields := []string{
		common.BKAppIDField,
		common.BKAppNameField,
	}

	query := &metadata.QueryBusinessRequest{
		Fields: fields,
		Page:   metadata.BasePage{},
		Condition: mapstr.MapStr{
			common.BKDataStatusField: mapstr.MapStr{common.BKDBNE: common.DataStatusDisabled},
		},
	}
	cnt, instItems, err := s.Core.BusinessOperation().FindBusiness(params, query)
	if nil != err {
		blog.Errorf("ListAllBusinessSimplify failed, FindBusiness failed, err: %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}
	businesses := make([]metadata.BizBasicInfo, 0)
	for _, item := range instItems {
		business := metadata.BizBasicInfo{}
		if err := mapstruct.Decode2Struct(item, &business); err != nil {
			blog.Errorf("ListAllBusinessSimplify failed, decode biz from db failed, err: %s, rid: %s", err.Error(), params.ReqID)
			return nil, params.Err.CCError(common.CCErrCommParseDBFailed)
		}
		businesses = append(businesses, business)
	}

	result := map[string]interface{}{
		"count": cnt,
		"info":  businesses,
	}

	return result, nil
}
