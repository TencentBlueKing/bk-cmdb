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
		blog.Errorf("create dynamic group failed, decode request body err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	if err := s.createGroupParamCheck(ctx.Kit, newDynamicGroup); err != nil {
		blog.Errorf("create request param check failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
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
			blog.Errorf("create dynamic group failed, err: %v, input: %+v, rid: %s", err, newDynamicGroup,
				ctx.Kit.Rid)
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
			blog.Errorf("generate audit log failed, err: %v, id: %s, rid: %s", err, newDynamicGroup.ID, ctx.Kit.Rid)
			return err
		}
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed, err: %v, id: %s, rid: %s", err, newDynamicGroup.ID, ctx.Kit.Rid)
			return err
		}

		// register dynamic group resource create action to iam.
		if auth.EnableAuthorize() {
			if err := s.registerActionToIAM(ctx.Kit, newDynamicGroup); err != nil {
				blog.Errorf("register created new dynamic group to iam failed, err: %v, rid: %s", err, ctx.Kit.Rid)
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

// registerActionToIam register dynamic group resource create action to iam.
func (s *Service) registerActionToIAM(kit *rest.Kit, dynamicGroup meta.DynamicGroup) error {
	bizID := strconv.FormatInt(dynamicGroup.AppID, 10)
	resp, err := s.CoreAPI.CoreService().Host().GetDynamicGroup(kit.Ctx, bizID, dynamicGroup.ID, kit.Header)
	if err != nil {
		blog.Errorf("get created new dynamic group failed, err: %v, biz: %s, ID: %s, rid: %s", err, bizID,
			dynamicGroup.ID, kit.Rid)
		return err
	}
	if !resp.Result {
		blog.Errorf("get created new dynamic group failed, err: response result is false, "+
			"biz: %s, ID: %s, rid: %s", bizID, dynamicGroup.ID, kit.Rid)
		return resp.CCError()
	}

	iamInstance := meta.IamInstanceWithCreator{
		Type:    string(iam.BizCustomQuery),
		ID:      resp.Data.ID,
		Name:    resp.Data.Name,
		Creator: kit.User,
	}
	if _, err = s.AuthManager.Authorizer.RegisterResourceCreatorAction(kit.Ctx, kit.Header, iamInstance); err != nil {
		blog.Errorf("register created new dynamic group to iam failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

// createGroupParamCheck 新建动态分组接口请求参数检查
func (s *Service) createGroupParamCheck(kit *rest.Kit, dynamicGroup meta.DynamicGroup) error {
	//  validate dynamic group func.
	validateFunc := func(objectID string) ([]meta.Attribute, error) {
		return logics.NewLogics(s.Engine, s.CacheDB, s.AuthManager).
			SearchObjectAttributes(kit, dynamicGroup.AppID, objectID)
	}

	if err := dynamicGroup.Validate(validateFunc); err != nil {
		blog.Errorf("create dynamic group failed, invalid param, err: %v, input: %+v, rid: %s",
			err, dynamicGroup, kit.Rid)
		return err
	}

	return nil
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
		blog.Errorf("update dynamic group failed, invalid bizID from path, err: %v, bizID: %s, rid: %s",
			err, bizID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id"))
		return
	}

	// decode update request data(DynamicGroup struct to interfaces).
	params := make(map[string]interface{})
	if err := ctx.DecodeInto(&params); err != nil {
		blog.Errorf("update dynamic group failed, decode request body, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	// final updates.
	updates := make(map[string]interface{})
	if err := s.updateGroupParamCheck(ctx.Kit, bizIDInt64, params, updates); err != nil {
		blog.Errorf("update request param check failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// update base on auto run txn with func.
	autoRunTxnFunc := func() error {
		// audit log.
		audit := auditlog.NewDynamicGroupAuditLog(s.CoreAPI.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, meta.AuditUpdate).WithUpdateFields(updates)
		auditLogs, err := audit.GenerateAuditLog(auditParam, &meta.DynamicGroup{ID: targetID, AppID: bizIDInt64})
		if err != nil {
			blog.Errorf("generate audit log failed after update dynamic group[%s], err: %v, rid: %s",
				targetID, err, ctx.Kit.Rid)
			return err
		}
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed after update dynamic group[%s], err: %v, rid: %s", targetID, err,
				ctx.Kit.Rid)
			return err
		}

		response, err := s.CoreAPI.CoreService().Host().UpdateDynamicGroup(ctx.Kit.Ctx, bizID, targetID,
			ctx.Kit.Header, updates)
		if err != nil {
			blog.Errorf("update dynamic group failed, err: %v, biz: %s, input: %+v, rid: %s",
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

// updateGroupParamCheck 更新动态分组接口请求参数检查
func (s *Service) updateGroupParamCheck(kit *rest.Kit, bizID int64, params, updates map[string]interface{}) error {

	if info, isExist := params["info"]; isExist {
		// update dynamic group info.
		row, err := json.Marshal(info)
		if err != nil {
			blog.Errorf("update dynamic group failed, invalid info, err: %v, info: %+v, rid: %s", err, info, kit.Rid)
			return err
		}

		dynamicGroupInfo := &meta.DynamicGroupInfo{}
		if err = json.Unmarshal(row, dynamicGroupInfo); err != nil {
			blog.Errorf("update dynamic group failed, invalid info, err: %v, row: %+v, rid: %s", err, row, kit.Rid)
			return err
		}

		objectIDParam, isObjIDExist := params["bk_obj_id"]
		if !isObjIDExist {
			blog.Errorf("update dynamic group failed, err: bk_obj_id is required in update info condition action,"+
				" input: %+v, rid: %s", params, kit.Rid)
			return err
		}
		objectID, ok := objectIDParam.(string)
		if !ok {
			blog.Errorf("update dynamic group failed, err: invalid bk_obj_id type, objectID: %+v, rid: %s",
				objectIDParam, kit.Rid)
			return err
		}

		//  validate dynamic group func.
		validatefunc := func(objectID string) ([]meta.Attribute, error) {
			return logics.NewLogics(s.Engine, s.CacheDB, s.AuthManager).
				SearchObjectAttributes(kit, bizID, objectID)
		}
		if err := dynamicGroupInfo.Validate(objectID, validatefunc); err != nil {
			blog.Errorf("update dynamic group failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		updates[common.BKObjIDField] = objectID
		updates["info"] = dynamicGroupInfo

	} else {
		_, isExist := params[common.BKObjIDField]
		if isExist {
			blog.Errorf("update dynamic group failed, err: info.condition is required in update bk_obj_id action, "+
				"input: %+v, rid: %s", params, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid)
		}

		_, isExist = params[common.BKFieldName]
		if !isExist {
			blog.Errorf("update dynamic group failed, err: empty update content, bk_biz_id/info/name, input: %+v, "+
				"rid: %s", params, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid)
		}
	}
	// update name.
	if name, isExist := params[common.BKFieldName]; isExist {
		updates[common.BKFieldName] = name
	}
	return nil
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
		blog.Errorf("delete dynamic group failed, err: %v, biz: %s, rid: %s", err, bizID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !result.Result {
		blog.Errorf("delete dynamic group failed, errcode: %d, errmsg: %s, bizID: %s, rid: %s",
			result.Code, result.ErrMsg, bizID, ctx.Kit.Rid)
		ctx.RespAutoError(result.CCError())
		return
	}
	dynamicGroup := result.Data

	// delete base on auto run txn with func.
	autoRunTxnFunc := func() error {
		result, err := s.CoreAPI.CoreService().Host().DeleteDynamicGroup(ctx.Kit.Ctx, bizID, targetID, ctx.Kit.Header)
		if err != nil {
			blog.Errorf("delete dynamic group failed, err: %v, biz: %s, rid: %s", err, bizID, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.Errorf("delete dynamic group failed, errcode: %d, errmsg: %s, bizID: %s, rid: %s",
				result.Code, result.ErrMsg, bizID, ctx.Kit.Rid)
			return result.CCError()
		}

		// audit log.
		audit := auditlog.NewDynamicGroupAuditLog(s.CoreAPI.CoreService())
		auditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, meta.AuditDelete)
		auditLogs, err := audit.GenerateAuditLog(auditParam, &dynamicGroup)
		if err != nil {
			blog.Errorf("generate audit log failed after delete dynamic group[%s], err: %v, rid: %s", targetID, err,
				ctx.Kit.Rid)
			return err
		}
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			blog.Errorf("save audit log failed after delete dynamic group[%s], err: %v, rid: %s", targetID, err,
				ctx.Kit.Rid)
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
		blog.Errorf("get dynamic group failed, err: %v, biz: %s, ID: %s, rid: %s",
			err, bizID, targetID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !result.Result {
		blog.Errorf("get dynamic group failed, errcode: %d, errmsg: %s, bizID: %s, ID: %s, rid: %s",
			result.Code, result.ErrMsg, bizID, targetID, ctx.Kit.Rid)
		ctx.RespAutoError(result.CCError())
		return
	}
	changeTimeToMatchLocalZone(result.Data.Info.Condition)
	changeTimeToMatchLocalZone(result.Data.Info.VariableCondition)
	ctx.RespEntity(result.Data)
}

// SearchDynamicGroup returns dynamic group list with target conditions.
func (s *Service) SearchDynamicGroup(ctx *rest.Contexts) {
	req := ctx.Request

	// application ID.
	bizID := req.PathParameter("bk_biz_id")

	bizIDInt64, err := strconv.ParseInt(bizID, 10, 64)
	if err != nil {
		blog.Errorf("search dynamic groups failed, invalid bizID from path, bizID: %s, rid: %s",
			bizID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "bk_biz_id"))
		return
	}

	input := new(meta.QueryCondition)
	if err := ctx.DecodeInto(input); err != nil {
		blog.Errorf("search dynamic groups failed, decode request body err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	if input.Page.Start < 0 {
		blog.Errorf("search dynamic groups failed, invalid page start param, input: %+v, rid: %s",
			input, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "start"))
		return
	}
	if input.Page.IsIllegal() {
		blog.Errorf("search dynamic groups failed, invalid page limit param, input: %+v, rid: %s",
			input, ctx.Kit.Rid)
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
		blog.Errorf("search dynamic groups failed, err: %v, biz: %s, input: %+v, rid: %s", err, bizID, input,
			ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !result.Result {
		blog.Errorf("search dynamic groups failed, errcode: %d, errmsg: %s, bizID: %s, input: %+v, rid: %s",
			result.Code, result.ErrMsg, bizID, input, ctx.Kit.Rid)
		ctx.RespAutoError(result.CCError())
		return
	}
	for _, dynamicGroup := range result.Data.Info {
		changeTimeToMatchLocalZone(dynamicGroup.Info.Condition)
		changeTimeToMatchLocalZone(dynamicGroup.Info.VariableCondition)
	}
	ctx.RespEntity(result.Data)
}

// ExecuteDynamicGroup executes target dynamic group and returns search datas.
func (s *Service) ExecuteDynamicGroup(ctx *rest.Contexts) {
	req := ctx.Request

	// target dynamic group ID.
	targetID := req.PathParameter(common.BKFieldID)
	bizID := req.PathParameter(common.BKAppIDField)
	bizIDInt64, err := strconv.ParseInt(bizID, 10, 64)
	if err != nil {
		blog.Errorf("execute dynamic group failed, invalid bizID from path, err: %v, bizID: %s, rid: %s",
			err, bizID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField))
		return
	}

	input := new(meta.ExecuteOption)
	if err := ctx.DecodeInto(input); err != nil {
		blog.Errorf("execute dynamic group failed, decode request body err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	result, searchConditions, err := s.checkAndBuildParam(ctx.Kit, input, bizIDInt64, targetID)
	if err != nil {
		blog.Errorf("check and build request param failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// target dynamic group.
	targetDynamicGroup := result.Data
	// execute dynamic group with target object type.
	searchPage := input.Page
	switch targetDynamicGroup.ObjID {
	case common.BKInnerObjIDHost:
		// build host search conditions.
		searchHostCondition := meta.HostCommonSearch{AppID: bizIDInt64, Condition: searchConditions, Page: searchPage}
		// execute host object dynamic group.
		data, err := logics.NewLogics(s.Engine, s.CacheDB, s.AuthManager).
			ExecuteHostDynamicGroup(ctx.Kit, &searchHostCondition, input.Fields, input.DisableCounter)
		if err != nil {
			blog.Errorf("execute dynamic group failed, search hosts, err: %v, bizID: %s, ID: %s, rid: %s",
				err, bizID, targetID, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrGetUserCustomQueryDetailFailed, err.Error()))
			return
		}

		ctx.RespEntity(meta.InstDataInfo{
			Count: data.Count,
			Info:  data.Info,
		})
		return

	case common.BKInnerObjIDSet:
		// build set search conditions.
		searchSetCondition := meta.SetCommonSearch{AppID: bizIDInt64, Condition: searchConditions, Page: searchPage}
		// execute set object dynamic group.
		data, err := logics.NewLogics(s.Engine, s.CacheDB, s.AuthManager).
			ExecuteSetDynamicGroup(ctx.Kit, &searchSetCondition, input.Fields, input.DisableCounter)
		if err != nil {
			blog.Errorf("execute dynamic group failed, search sets, err: %v, bizID: %s, ID: %s, rid: %s",
				err, bizID, targetID, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrGetUserCustomQueryDetailFailed, err.Error()))
			return
		}

		ctx.RespEntity(meta.InstDataInfo{
			Count: data.Count,
			Info:  data.Info,
		})
		return

	default:
		// unknown group object type, should no-reach here.
		blog.Errorf("execute dynamic group failed, err: unknown group object type: [%s], data: %v, "+
			"bizID: %s, ID: %s, rid: %s", targetDynamicGroup.ObjID, result.Data, bizID, targetID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCSystemUnknownError))
	}
}

// checkAndBuildParam 执行动态分组接口请求参数检查和返回查询参数
func (s *Service) checkAndBuildParam(kit *rest.Kit, input *meta.ExecuteOption, bizID int64, targetID string) (
	*meta.GetDynamicGroupResult, []meta.SearchCondition, error) {

	if len(input.Fields) == 0 {
		blog.Errorf("execute dynamic group failed, err: fields is empty, input: %+v, rid: %s", input, kit.Rid)
		return &meta.GetDynamicGroupResult{}, nil, kit.CCError.Error(common.CCErrCommParamsIsInvalid)
	}

	if input.Page.Start < 0 {
		blog.Errorf("execute dynamic group failed, err: invalid page start param, input: %+v, rid: %s",
			input, kit.Rid)
		return &meta.GetDynamicGroupResult{}, nil, kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed)
	}
	if input.Page.IsIllegal() {
		blog.Errorf("execute dynamic group failed, err: invalid page limit param, input: %+v, rid: %s",
			input, kit.Rid)
		return &meta.GetDynamicGroupResult{}, nil, kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	// query target dynamic group.
	result, err := s.CoreAPI.CoreService().Host().GetDynamicGroup(kit.Ctx, strconv.FormatInt(bizID, 10), targetID,
		kit.Header)
	if err != nil {
		blog.Errorf("execute dynamic group failed, err: %v, bizID: %d, ID: %s, rid: %s", err, bizID, targetID, kit.Rid)
		return nil, nil, err
	}
	if !result.Result {
		blog.Errorf("execute dynamic group failed, errcode: %d, errmsg: %s, bizID: %d, ID: %s, rid: %s",
			result.Code, result.ErrMsg, bizID, targetID, kit.Rid)
		return nil, nil, result.CCError()
	}
	if len(result.Data.Name) == 0 {
		blog.Errorf("execute dynamic group failed, err: group not found, bizID: %d, ID: %s, rid: %s",
			bizID, targetID, kit.Rid)
		return nil, nil, kit.CCError.Error(common.CCErrCommNotFound)
	}

	validateFunc := func(objectID string) ([]meta.Attribute, error) {
		return logics.NewLogics(s.Engine, s.CacheDB, s.AuthManager).SearchObjectAttributes(kit, bizID, objectID)
	}
	if _, err = meta.ValidDynamicGroupCond(input.VariableCondition, result.Data.ObjID, validateFunc,
		map[string]map[string]struct{}{}); err != nil {
		blog.Errorf("dynamic group info condition is invalid, input: %+v, err: %v, rid: %s", input, err, kit.Rid)
		return nil, nil, err
	}

	info := result.Data.Info
	info.VariableCondition, err = buildFinalCond(kit, input.VariableCondition, info.VariableCondition)
	if err != nil {
		blog.Errorf("build final variable condition failed, info: %+v, input: %+v, err: %v, rid: %s", input, err,
			kit.Rid)
		return nil, nil, err
	}

	conditionMap := make(map[string]*meta.SearchCondition)
	conditionMap = parseCond(conditionMap, info.Condition)
	conditionMap = parseCond(conditionMap, info.VariableCondition)

	searchConditions := make([]meta.SearchCondition, 0)
	for _, condition := range conditionMap {
		searchConditions = append(searchConditions, *condition)
	}

	return result, searchConditions, nil
}

func buildFinalCond(kit *rest.Kit, reqCondArr []meta.DynamicGroupInfoCondition,
	originCondArr []meta.DynamicGroupInfoCondition) ([]meta.DynamicGroupInfoCondition, error) {

	if len(reqCondArr) == 0 {
		return originCondArr, nil
	}

	if len(reqCondArr) != 0 && len(originCondArr) == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "variable_condition")
	}

	reqCondMap := make(map[string]meta.DynamicGroupInfoCondition)
	for _, cond := range reqCondArr {
		reqCondMap[cond.ObjID] = cond
	}

	for idx, cond := range originCondArr {
		reqCond, ok := reqCondMap[cond.ObjID]
		if !ok {
			continue
		}

		updateCondMap := make(map[string]meta.DynamicGroupCondition)
		for _, fieldCond := range reqCond.Condition {
			updateCondMap[fieldCond.Field] = fieldCond
		}

		updateTimeRuleMap := make(map[string]meta.TimeConditionItem)
		if reqCond.TimeCondition != nil {
			for _, rule := range reqCond.TimeCondition.Rules {
				updateTimeRuleMap[rule.Field] = rule
			}
		}

		for subIdx, subCond := range cond.Condition {
			updateCond, exist := updateCondMap[subCond.Field]
			if !exist {
				continue
			}

			originCondArr[idx].Condition[subIdx] = updateCond
		}

		if cond.TimeCondition == nil {
			continue
		}

		for subIdx, rule := range cond.TimeCondition.Rules {
			updateTimeRule, exist := updateTimeRuleMap[rule.Field]
			if !exist {
				continue
			}

			originCondArr[idx].TimeCondition.Rules[subIdx] = updateTimeRule
		}
	}

	return originCondArr, nil
}

func parseCond(conditionMap map[string]*meta.SearchCondition,
	conditions []meta.DynamicGroupInfoCondition) map[string]*meta.SearchCondition {

	for _, cond := range conditions {
		if conditionMap[cond.ObjID] == nil {
			conditionMap[cond.ObjID] = &meta.SearchCondition{ObjectID: cond.ObjID}
		}

		for _, item := range cond.Condition {
			condItem := meta.ConditionItem{Field: item.Field, Operator: item.Operator, Value: item.Value}
			conditionMap[cond.ObjID].Condition = append(conditionMap[cond.ObjID].Condition, condItem)
		}

		if conditionMap[cond.ObjID].TimeCondition == nil {
			conditionMap[cond.ObjID].TimeCondition = cond.TimeCondition
			continue
		}

		conditionMap[cond.ObjID].TimeCondition.Rules = append(conditionMap[cond.ObjID].TimeCondition.Rules,
			cond.TimeCondition.Rules...)
	}

	return conditionMap
}

// changeTimeToMatchLocalZone TODO
// change the time in UTC format to the time in the local time zone
func changeTimeToMatchLocalZone(conditions []meta.DynamicGroupInfoCondition) {
	for _, condition := range conditions {
		if condition.TimeCondition == nil || condition.TimeCondition.Rules == nil {
			continue
		}
		for _, rule := range condition.TimeCondition.Rules {
			if rule.Start != nil {
				rule.Start.Time = rule.Start.Local()
			}
			if rule.End != nil {
				rule.End.Time = rule.End.Local()
			}
		}
	}
}
