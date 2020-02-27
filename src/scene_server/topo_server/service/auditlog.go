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

	"configcenter/src/auth"
	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

const CCTimeTypeParseFlag = "cc_time_type"

// AuditQuery search audit logs
func (s *Service) AuditQuery(ctx *rest.Contexts) {
	query := metadata.QueryInput{}
	if err := ctx.DecodeInto(&query); nil != err {
		ctx.RespAutoError(err)
		return
	}

	queryCondition := query.Condition
	if nil == queryCondition {
		query.Condition = common.KvMap{}
	} else {
		cond := queryCondition
		times, ok := cond[common.BKOpTimeField].([]interface{})
		if ok {
			if 2 != len(times) {
				blog.Errorf("search operation log input params times error, info: %v, rid: %s", times, ctx.Kit.Rid)
				ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsInvalid))
				return
			}

			cond[common.BKOpTimeField] = common.KvMap{
				common.BKDBGTE:      times[0],
				common.BKDBLTE:      times[1],
				CCTimeTypeParseFlag: "1",
			}
		}
		query.Condition = cond
	}
	if 0 == query.Limit {
		query.Limit = common.BKDefaultLimit
	}

	// add auth filter condition
	var businessID int64
	bizID, exist := query.Condition[common.BKAppIDField]
	if exist == true {
		id, err := util.GetInt64ByInterface(bizID)
		if err != nil {
			blog.Errorf("%s field in query condition but parse int failed, err: %+v, rid: %s", common.BKAppIDField, err, ctx.Kit.Rid)
		}
		businessID = id
	}

	// switch between two different control mechanism
	// TODO use global authorization for now, need more specific auth control
	if s.AuthManager.RegisterAuditCategoryEnabled == false {
		if err := s.AuthManager.AuthorizeAuditRead(ctx.Kit.Ctx, ctx.Kit.Header, 0); err != nil {
			blog.Errorf("AuditQuery failed, authorize failed, AuthorizeAuditRead failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			resp, err := s.AuthManager.GenAuthorizeAuditReadNoPermissionsResponse(ctx.Kit.Ctx, ctx.Kit.Header, 0)
			if err != nil {
				ctx.RespAutoError(fmt.Errorf("try authorize failed, err: %v", err))
				return
			}
			ctx.RespEntityWithError(resp, auth.NoAuthorizeError)
			return
		}
	} else {
		var hasAuthorize bool
		for _, bizID := range []int64{businessID, 0} {
			authCondition, hasAuthorization, err := s.AuthManager.MakeAuthorizedAuditListCondition(ctx.Kit.Ctx, ctx.Kit.Header, bizID)
			if err != nil {
				blog.Errorf("AuditQuery failed, make audit query condition from auth failed, %+v, rid: %s", err, ctx.Kit.Rid)
				ctx.RespAutoError(fmt.Errorf("make audit query condition from auth failed, %+v", err))
				return
			}

			if hasAuthorization == true {
				query.Condition[common.BKDBOR] = authCondition
				blog.V(5).Infof("AuditQuery, auth condition is: %+v, rid: %s", authCondition, ctx.Kit.Rid)
				hasAuthorize = hasAuthorization
				break
			}
		}
		if hasAuthorize == false {
			blog.Errorf("AuditQuery failed, user %+v has no authorization on audit, rid: %s", ctx.Kit.User, ctx.Kit.Rid)
			ctx.RespAutoError(auth.NoAuthorizeError)
			return
		}
	}

	blog.V(5).Infof("AuditQuery, AuditOperation parameter: %+v, rid: %s", query, ctx.Kit.Rid)
	resp, err := s.Core.AuditOperation().Query(ctx.Kit, query)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// InstanceAuditQuery search instance audit logs
// current use case: get host and process related audit log in cmdb web
func (s *Service) InstanceAuditQuery(ctx *rest.Contexts) {
	query := metadata.QueryInput{}
	if err := ctx.DecodeInto(&query); nil != err {
		ctx.RespAutoError(err)
		return
	}

	objectID := ctx.Request.PathParameter("bk_obj_id")
	if len(objectID) == 0 {
		blog.Errorf("InstanceAuditQuery failed, host audit query condition can't be empty, query: %+v, rid: %s", query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_obj_id"))
		return
	}

	queryCondition := query.Condition
	if nil == queryCondition {
		blog.Errorf("InstanceAuditQuery failed, host audit query condition can't be empty, query: %+v, rid: %s", query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "condition"))
		return
	}

	cond := queryCondition
	times, ok := cond[common.BKOpTimeField].([]interface{})
	if ok {
		if 2 != len(times) {
			blog.Errorf("InstanceAuditQuery failed, search operation log input params times error, info: %v, rid: %s", times, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "op_time"))
			return
		}

		cond[common.BKOpTimeField] = common.KvMap{
			"$gte":              times[0],
			"$lte":              times[1],
			CCTimeTypeParseFlag: "1",
		}
	}
	cond[common.BKOpTargetField] = objectID
	query.Condition = cond
	if 0 == query.Limit {
		query.Limit = common.BKDefaultLimit
	}

	// auth: check authorization on instance
	var businessID int64
	bizID, exist := query.Condition[common.BKAppIDField]
	if exist == true {
		id, err := util.GetInt64ByInterface(bizID)
		if err != nil {
			blog.Errorf("InstanceAuditQuery failed, %s field in query condition but parse int failed, err: %+v, rid: %s", common.BKAppIDField, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
			return
		}
		businessID = id
	}

	instID, exist := queryCondition["inst_id"]
	if exist == false {
		blog.Errorf("InstanceAuditQuery failed, instance audit query condition condition.ext_key not exist, query: %+v, rid: %s", query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "inst_id"))
		return
	}
	instanceID, err := util.GetInt64ByInterface(instID)
	if err != nil {
		blog.Errorf("InstanceAuditQuery failed, instance audit query condition instanceID in condition.ext_key.$in invalid, instanceID: %+v, query: %+v, rid: %s", instID, query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "inst_id"))
		return
	}

	opTarget, exist := queryCondition["op_target"]
	if exist {
		target, ok := opTarget.(string)
		if !ok {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "op_target"))
			return
		}
		if target == "biz" {
			businessID = instanceID
		}
	}

	action := meta.Find
	switch objectID {
	case common.BKInnerObjIDHost:
		err = s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, action, instanceID)
	case common.BKInnerObjIDProc:
		err = s.AuthManager.AuthorizeByProcessID(ctx.Kit.Ctx, ctx.Kit.Header, action, instanceID)
		if err != nil && err == auth.NoAuthorizeError {
			resp, err := s.AuthManager.GenProcessNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, businessID)
			if err != nil {
				ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrTopoGetAppFailed, bizID))
				return
			}
			ctx.RespEntityWithError(resp, auth.NoAuthorizeError)
			return
		}
	case common.BKInnerObjIDModule:
		err = s.AuthManager.AuthorizeByModuleID(ctx.Kit.Ctx, ctx.Kit.Header, action, instanceID)
		if err != nil && err == auth.NoAuthorizeError {
			ctx.RespEntityWithError(s.AuthManager.GenModuleSetNoPermissionResp(), auth.NoAuthorizeError)
			return
		}
	case common.BKInnerObjIDSet:
		err = s.AuthManager.AuthorizeBySetID(ctx.Kit.Ctx, ctx.Kit.Header, action, instanceID)
		if err != nil && err == auth.NoAuthorizeError {
			ctx.RespEntityWithError(s.AuthManager.GenModuleSetNoPermissionResp(), auth.NoAuthorizeError)
			return
		}
	case common.BKInnerObjIDApp:
		err = s.AuthManager.AuthorizeByBusinessID(ctx.Kit.Ctx, ctx.Kit.Header, action, instanceID)
		if err != nil && err == auth.NoAuthorizeError {
			resp, err := s.AuthManager.GenBusinessAuditNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, businessID)
			if err != nil {
				ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrTopoGetAppFailed))
				return
			}
			ctx.RespEntityWithError(resp, auth.NoAuthorizeError)
			return
		}
	default:
		err = s.AuthManager.AuthorizeByInstanceID(ctx.Kit.Ctx, ctx.Kit.Header, action, objectID, instanceID)
	}
	if err != nil {
		blog.Errorf("InstanceAuditQuery failed, query instance audit log failed, authorization on instance of model %s failed, err: %+v, rid: %s", objectID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommAuthorizeFailed))
		return
	}

	blog.V(4).Infof("InstanceAuditQuery failed, AuditOperation parameter: %+v, rid: %s", query, ctx.Kit.Rid)
	resp, err := s.Core.AuditOperation().Query(ctx.Kit, query)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}
