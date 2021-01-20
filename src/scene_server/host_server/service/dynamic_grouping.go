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
	"time"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	meta "configcenter/src/common/metadata"
	parser "configcenter/src/common/paraparse"
	"configcenter/src/scene_server/host_server/logics"
)

// CreateDynamicGroup creates a new dynamic group object.
func (s *Service) CreateDynamicGroup(ctx *rest.Contexts) {
	newDynamicGroup := meta.DynamicGroup{}
	if err := ctx.DecodeInto(&newDynamicGroup); err != nil {
		blog.Errorf("create dynamic group failed, decode request body err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	//  validate dynamic group func.
	validatefunc := func(objectID string) ([]meta.Attribute, error) {
		return logics.NewLogics(s.Engine, s.CacheDB, s.AuthManager).
			SearchObjectAttributes(ctx.Kit, newDynamicGroup.AppID, objectID)
	}

	if err := newDynamicGroup.Validate(validatefunc); err != nil {
		blog.Errorf("create dynamic group failed, invalid param: %+v, input: %+v, rid: %s", err, newDynamicGroup, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, err.Error()))
		return
	}
	newDynamicGroup.CreateUser = ctx.Kit.User
	newDynamicGroup.CreateTime = time.Now().UTC()
	response := &meta.IDResult{}

	// create base on auto run txn with func.
	autoRunTxnFunc := func() error {
		var err error
		response, err = s.CoreAPI.CoreService().Host().CreateDynamicGroup(ctx.Kit.Ctx, ctx.Kit.Header, &newDynamicGroup)
		if err != nil {
			blog.Errorf("create dynamic group failed, err: %+v, input: %+v, rid: %s", err, newDynamicGroup, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !response.Result {
			blog.Errorf("create dynamic group failed, errcode: %d, errmsg: %s, input: %+v, rid: %s",
				response.Code, response.ErrMsg, newDynamicGroup, ctx.Kit.Rid)
			return response.CCError()
		}
		newDynamicGroup.ID = response.Data.ID

		// audit log.
		audit := auditlog.NewDynamicGroupAuditLog(s.CoreAPI.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, meta.AuditCreate)
		auditLogs, err := audit.GenerateAuditLog(auditParam, &newDynamicGroup)
		if err != nil {
			blog.Errorf("generate audit log failed after create new dynamic group[%s], err: %+v, rid: %s", newDynamicGroup.ID, err, ctx.Kit.Rid)
			return err
		}
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed after create new dynamic group[%s], err: %+v, rid: %s", newDynamicGroup.ID, err, ctx.Kit.Rid)
			return err
		}

		// register dynamic group resource create action to iam.
		if auth.EnableAuthorize() {
			bizID := strconv.FormatInt(newDynamicGroup.AppID, 10)

			resp, err := s.CoreAPI.CoreService().Host().GetDynamicGroup(ctx.Kit.Ctx, bizID, newDynamicGroup.ID, ctx.Kit.Header)
			if err != nil {
				blog.Errorf("get created new dynamic group failed, err: %+v, biz: %s, ID: %s, rid: %s", err, bizID, newDynamicGroup.ID, ctx.Kit.Rid)
				return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
			}
			if !resp.Result {
				blog.Errorf("get created new dynamic group failed, err: %+v, biz: %s, ID: %s, rid: %s", resp.ErrMsg, bizID, newDynamicGroup.ID, ctx.Kit.Rid)
				return resp.CCError()
			}

			iamInstance := meta.IamInstanceWithCreator{
				Type:    string(iam.BizCustomQuery),
				ID:      resp.Data.ID,
				Name:    resp.Data.Name,
				Creator: ctx.Kit.User,
			}

			if _, err = s.AuthManager.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance); err != nil {
				blog.Errorf("register created new dynamic group to iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
				return err
			}
		}
		return nil
	}

	// do create action now.
	if err := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, autoRunTxnFunc); err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(response.Data)
}

// UpdateDynamicGroup updates target dynamic group.
func (s *Service) UpdateDynamicGroup(ctx *rest.Contexts) {
	req := ctx.Request

	// application ID.
	bizID := req.PathParameter("bk_biz_id")

	// target dynamic group ID.
	targetID := req.PathParameter("id")

	bizIDInt64, err := strconv.ParseInt(bizID, 10, 64)
	if err != nil {
		blog.Errorf("update dynamic group failed, invalid bizID from path, bizID: %s, rid: %s", bizID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id"))
		return
	}

	// decode update request data(DynamicGroup struct to interfaces).
	params := make(map[string]interface{})
	if err := ctx.DecodeInto(&params); err != nil {
		blog.Errorf("update dynamic group failed, decode request body err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	// final updates.
	updates := make(map[string]interface{})

	if info, isExist := params["info"]; isExist {
		// update dynamic group info.
		row, err := json.Marshal(info)
		if err != nil {
			blog.Errorf("update dynamic group failed, invalid info, info: %+v, rid: %s", info, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "info"))
			return
		}

		dynamicGroupInfo := &meta.DynamicGroupInfo{}
		if err = json.Unmarshal(row, dynamicGroupInfo); err != nil {
			blog.Errorf("update dynamic group failed, invalid info, info: %+v, rid: %s", info, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "info"))
			return
		}

		objectIDParam, isObjIDExist := params["bk_obj_id"]
		if !isObjIDExist {
			blog.Errorf("update dynamic group failed, bk_obj_id is required in update info condition action, input: %+v, rid: %s",
				params, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "bk_obj_id"))
			return
		}
		objectID, ok := objectIDParam.(string)
		if !ok {
			blog.Errorf("update dynamic group failed, invalid bk_obj_id type, objectID: %+v, rid: %s", objectIDParam, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "bk_obj_id"))
			return
		}

		//  validate dynamic group func.
		validatefunc := func(objectID string) ([]meta.Attribute, error) {
			return logics.NewLogics(s.Engine, s.CacheDB, s.AuthManager).
				SearchObjectAttributes(ctx.Kit, bizIDInt64, objectID)
		}

		if err := dynamicGroupInfo.Validate(objectID, validatefunc); err != nil {
			blog.Errorf("update dynamic group failed, invalid param: %+v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, err.Error()))
			return
		}
		updates[common.BKObjIDField] = objectID
		updates["info"] = dynamicGroupInfo

	} else {
		_, isExist := params[common.BKObjIDField]
		if isExist {
			blog.Errorf("update dynamic group failed, info.condition is required in update bk_obj_id action, input: %+v, rid: %s",
				params, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "info.condition"))
			return
		}

		_, isExist = params[common.BKFieldName]
		if !isExist {
			blog.Errorf("update dynamic group failed, empty update content, bk_biz_id/info/name, input: %+v, rid: %s",
				params, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "info.condition"))
			return
		}
	}

	// update name.
	if name, isExist := params[common.BKFieldName]; isExist {
		updates[common.BKFieldName] = name
	}

	// update base on auto run txn with func.
	autoRunTxnFunc := func() error {
		// audit log.
		audit := auditlog.NewDynamicGroupAuditLog(s.CoreAPI.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, meta.AuditUpdate).WithUpdateFields(updates)
		auditLogs, err := audit.GenerateAuditLog(auditParam, &meta.DynamicGroup{ID: targetID, AppID: bizIDInt64})
		if err != nil {
			blog.Errorf("generate audit log failed after update dynamic group[%s], err: %+v, rid: %s", targetID, err, ctx.Kit.Rid)
			return err
		}
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed after update dynamic group[%s], err: %+v, rid: %s", targetID, err, ctx.Kit.Rid)
			return err
		}

		response, err := s.CoreAPI.CoreService().Host().UpdateDynamicGroup(ctx.Kit.Ctx, bizID, targetID, ctx.Kit.Header, updates)
		if err != nil {
			blog.Errorf("update dynamic group failed, err: %+v, biz: %s, input: %+v, rid: %s",
				err, bizID, params, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !response.Result {
			blog.Errorf("update dynamic group failed, errcode: %d, errmsg: %s, bizID: %s, input: %+v, rid: %s",
				response.Code, response.ErrMsg, bizID, params, ctx.Kit.Rid)
			return response.CCError()
		}
		return nil
	}

	// do update action now.
	if err := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, autoRunTxnFunc); err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// DeleteDynamicGroup deletes target dynamic group.
func (s *Service) DeleteDynamicGroup(ctx *rest.Contexts) {
	req := ctx.Request

	// application ID.
	bizID := req.PathParameter("bk_biz_id")

	// target dynamic group ID.
	targetID := req.PathParameter("id")

	// pre-query target dynamic group object.
	result, err := s.CoreAPI.CoreService().Host().GetDynamicGroup(ctx.Kit.Ctx, bizID, targetID, ctx.Kit.Header)
	if err != nil {
		blog.Errorf("delete dynamic group failed, err: %+v, biz: %s, rid: %s", err, bizID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !result.Result {
		blog.Errorf("delete dynamic group failed, errcode: %d, errmsg: %s, bizID: %s, rid: %s", result.Code, result.ErrMsg, bizID, ctx.Kit.Rid)
		ctx.RespAutoError(result.CCError())
		return
	}
	dynamicGroup := result.Data

	// delete base on auto run txn with func.
	autoRunTxnFunc := func() error {
		result, err := s.CoreAPI.CoreService().Host().DeleteDynamicGroup(ctx.Kit.Ctx, bizID, targetID, ctx.Kit.Header)
		if err != nil {
			blog.Errorf("delete dynamic group failed, err: %+v, biz: %s, rid: %s", err, bizID, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.Errorf("delete dynamic group failed, errcode: %d, errmsg: %s, bizID: %s, rid: %s", result.Code, result.ErrMsg, bizID, ctx.Kit.Rid)
			return result.CCError()
		}

		// audit log.
		audit := auditlog.NewDynamicGroupAuditLog(s.CoreAPI.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, meta.AuditDelete)
		auditLogs, err := audit.GenerateAuditLog(auditParam, &dynamicGroup)
		if err != nil {
			blog.Errorf("generate audit log failed after delete dynamic group[%s], err: %+v, rid: %s", targetID, err, ctx.Kit.Rid)
			return err
		}
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed after delete dynamic group[%s], err: %+v, rid: %s", targetID, err, ctx.Kit.Rid)
			return err
		}
		return nil
	}

	// do delete action now.
	if err := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, autoRunTxnFunc); err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// GetDynamicGroup returns target dynamic group detail.
func (s *Service) GetDynamicGroup(ctx *rest.Contexts) {
	req := ctx.Request

	// application ID.
	bizID := req.PathParameter("bk_biz_id")

	// target dynamic group ID.
	targetID := req.PathParameter("id")

	// do query action now.
	result, err := s.CoreAPI.CoreService().Host().GetDynamicGroup(ctx.Kit.Ctx, bizID, targetID, ctx.Kit.Header)
	if err != nil {
		blog.Errorf("get dynamic group failed, err: %+v, biz: %s, ID: %s, rid: %s", err, bizID, targetID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !result.Result {
		blog.Errorf("get dynamic group failed, errcode: %d, errmsg: %s, bizID: %s, ID: %s, rid: %s",
			result.Code, result.ErrMsg, bizID, targetID, ctx.Kit.Rid)
		ctx.RespAutoError(result.CCError())
		return
	}
	ctx.RespEntity(result.Data)
}

// SearchDynamicGroup returns dynamic group list with target conditions.
func (s *Service) SearchDynamicGroup(ctx *rest.Contexts) {
	req := ctx.Request

	// application ID.
	bizID := req.PathParameter("bk_biz_id")

	bizIDInt64, err := strconv.ParseInt(bizID, 10, 64)
	if err != nil {
		blog.Errorf("search dynamic groups failed, invalid bizID from path, bizID: %s, rid: %s", bizID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id"))
		return
	}

	input := new(meta.QueryCondition)
	if err := ctx.DecodeInto(input); err != nil {
		blog.Errorf("search dynamic groups failed, decode request body err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	if input.Page.Start < 0 {
		blog.Errorf("search dynamic groups failed, invalid page start param, input: %+v, rid: %s", input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "start"))
		return
	}
	if input.Page.IsIllegal() {
		blog.Errorf("search dynamic groups failed, invalid page limit param, input: %+v, rid: %s", input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "limit"))
		return
	}

	// handle search conditions.
	var condition map[string]interface{}
	if input.Condition != nil {
		condition = input.Condition
	} else {
		condition = make(map[string]interface{})
	}
	condition[common.BKAppIDField] = bizIDInt64

	// add like search when there is name field in condition.
	name, ok := condition["name"].(string)
	if ok && len(name) != 0 {
		condition["name"] = common.KvMap{common.BKDBLIKE: parser.SpecialCharChange(name)}
	}

	// reset final conditions.
	input.Condition = condition

	// do search action now.
	result, err := s.CoreAPI.CoreService().Host().SearchDynamicGroup(ctx.Kit.Ctx, ctx.Kit.Header, input)
	if err != nil {
		blog.Errorf("search dynamic groups failed, err: %+v, biz: %s, input: %+v, rid: %s", err, bizID, input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !result.Result {
		blog.Errorf("search dynamic groups failed, errcode: %d, errmsg: %s, bizID: %s, input: %+v, rid: %s",
			result.Code, result.ErrMsg, bizID, input, ctx.Kit.Rid)
		ctx.RespAutoError(result.CCError())
		return
	}
	ctx.RespEntity(result.Data)
}

// ExecuteDynamicGroup executes target dynamic group and returns search datas.
func (s *Service) ExecuteDynamicGroup(ctx *rest.Contexts) {
	req := ctx.Request

	// application ID.
	bizID := req.PathParameter("bk_biz_id")

	// target dynamic group ID.
	targetID := req.PathParameter("id")

	bizIDInt64, err := strconv.ParseInt(bizID, 10, 64)
	if err != nil {
		blog.Errorf("execute dynamic group failed, invalid bizID from path, bizID: %s, rid: %s", bizID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id"))
		return
	}

	input := new(meta.QueryCondition)
	if err := ctx.DecodeInto(input); err != nil {
		blog.Errorf("execute dynamic group failed, decode request body err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	if len(input.Fields) == 0 {
		blog.Errorf("execute dynamic group failed, empty fields input: %+v, rid: %s", input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "fields"))
		return
	}

	if input.Page.Start < 0 {
		blog.Errorf("execute dynamic group failed, invalid page start param, input: %+v, rid: %s", input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	if input.Page.IsIllegal() {
		blog.Errorf("execute dynamic group failed, invalid page limit param, input: %+v, rid: %s", input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	searchPage := input.Page

	// query target dynamic group.
	result, err := s.CoreAPI.CoreService().Host().GetDynamicGroup(ctx.Kit.Ctx, bizID, targetID, ctx.Kit.Header)
	if err != nil {
		blog.Errorf("execute dynamic group failed, err: %+v, bizID: %s, ID: %s, rid: %s", err, bizID, targetID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrGetUserCustomQueryDetailFailed, err.Error()))
		return
	}
	if !result.Result {
		blog.Errorf("execute dynamic group failed, errcode: %d, errmsg: %s, bizID: %s, ID: %s, rid: %s",
			result.Code, result.ErrMsg, bizID, targetID, ctx.Kit.Rid)
		ctx.RespAutoError(result.CCError())
		return
	}
	if len(result.Data.Name) == 0 {
		blog.Errorf("execute dynamic group failed, group not found, bizID: %s, ID: %s, rid: %s", bizID, targetID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommNotFound))
		return
	}

	// target dynamic group.
	targetDynamicGroup := result.Data

	// build search conditions.
	searchConditions := []meta.SearchCondition{}

	// parse all dynamic group conditions to search condition.
	for _, cond := range targetDynamicGroup.Info.Condition {
		searchCondition := meta.SearchCondition{ObjectID: cond.ObjID, Condition: []meta.ConditionItem{}}

		// build condition items.
		for _, item := range cond.Condition {
			condItem := meta.ConditionItem{Field: item.Field, Operator: item.Operator, Value: item.Value}
			searchCondition.Condition = append(searchCondition.Condition, condItem)
		}
		searchConditions = append(searchConditions, searchCondition)
	}

	// execute dynamic group with target object type.
	if targetDynamicGroup.ObjID == common.BKInnerObjIDHost {
		// build host search conditions.
		searchHostCondition := meta.HostCommonSearch{AppID: bizIDInt64, Condition: searchConditions, Page: searchPage}

		// execute host object dynamic group.
		data, err := logics.NewLogics(s.Engine, s.CacheDB, s.AuthManager).
			ExecuteHostDynamicGroup(ctx.Kit, &searchHostCondition, input.Fields, input.DisableCounter)
		if err != nil {
			blog.Errorf("execute dynamic group failed, search hosts, err: %+v, bizID: %s, ID: %s, rid: %s",
				err, bizID, targetID, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrGetUserCustomQueryDetailFailed, err.Error()))
			return
		}

		ctx.RespEntity(meta.InstDataInfo{
			Count: data.Count,
			Info:  data.Info,
		})
		return

	} else if targetDynamicGroup.ObjID == common.BKInnerObjIDSet {
		// build host search conditions.
		searchSetCondition := meta.SetCommonSearch{AppID: bizIDInt64, Condition: searchConditions, Page: searchPage}

		// execute set object dynamic group.
		data, err := logics.NewLogics(s.Engine, s.CacheDB, s.AuthManager).
			ExecuteSetDynamicGroup(ctx.Kit, &searchSetCondition, input.Fields, input.DisableCounter)
		if err != nil {
			blog.Errorf("execute dynamic group failed, search sets, err: %+v, bizID: %s, ID: %s, rid: %s",
				err, bizID, targetID, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrGetUserCustomQueryDetailFailed, err.Error()))
			return
		}

		ctx.RespEntity(meta.InstDataInfo{
			Count: data.Count,
			Info:  data.Info,
		})
		return
	}

	// unknown group object type, should no-reach here.
	blog.Errorf("execute dynamic group failed, unknown group object type[%s], bizID: %s, ID: %s, rid: %s",
		targetDynamicGroup.ObjID, bizID, targetID, ctx.Kit.Rid)
	ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCSystemUnknownError))
}
