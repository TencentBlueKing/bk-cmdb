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
	"configcenter/src/auth"
	"fmt"

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/types"
)

const CCTimeTypeParseFlag = "cc_time_type"

// AuditQuery search audit logs
func (s *Service) AuditQuery(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	query := metadata.QueryInput{}
	if err := data.MarshalJSONInto(&query); nil != err {
		blog.Errorf("[audit] failed to parse the input (%#v), error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommJSONUnmarshalFailed, err.Error())
	}

	queryCondition := query.Condition
	if nil == queryCondition {
		query.Condition = common.KvMap{common.BKOwnerIDField: params.SupplierAccount}
	} else {
		cond := queryCondition.(map[string]interface{})
		times, ok := cond[common.BKOpTimeField].([]interface{})
		if ok {
			if 2 != len(times) {
				blog.Errorf("search operation log input params times error, info: %v, rid: %s", times, params.ReqID)
				return nil, params.Err.Error(common.CCErrCommParamsInvalid)
			}

			cond[common.BKOpTimeField] = common.KvMap{
				common.BKDBGTE:      times[0],
				common.BKDBLTE:      times[1],
				CCTimeTypeParseFlag: "1",
			}
		}
		cond[common.BKOwnerIDField] = params.SupplierAccount
		query.Condition = cond
	}
	if 0 == query.Limit {
		query.Limit = common.BKDefaultLimit
	}

	// add auth filter condition
	var businessID int64
	bizID, exist := query.Condition.(map[string]interface{})[common.BKAppIDField]
	if exist == true {
		id, err := util.GetInt64ByInterface(bizID)
		if err != nil {
			blog.Errorf("%s field in query condition but parse int failed, err: %+v, rid: %s", common.BKAppIDField, err, params.ReqID)
		}
		businessID = id
	}

	// switch between tow different control mechanism
	if s.AuthManager.RegisterAuditCategoryEnabled == false {
		if err := s.AuthManager.AuthorizeAuditRead(params.Context, params.Header, 0); err != nil {
			blog.Infof("AuditQuery check authorize on global failed, we'll try to check authorize on business, err: %+v, rid: %s", err, params.ReqID)
			if err := s.AuthManager.AuthorizeAuditRead(params.Context, params.Header, businessID); err != nil {
				blog.Errorf("AuditQuery failed, authorize failed, AuthorizeAuditRead failed, err: %+v, rid: %s", err, params.ReqID)
				resp, err := s.AuthManager.GenAuthorizeAuditReadNoPermissionsResponse(params.Context, params.Header, businessID)
				if err != nil {
					return nil, fmt.Errorf("try authorize failed, err: %v", err)
				}
				return resp, auth.NoAuthorizeError
			}
		}
	} else {
		var hasAuthorize bool
		for _, bizID := range []int64{businessID, 0} {
			authCondition, hasAuthorization, err := s.AuthManager.MakeAuthorizedAuditListCondition(params.Context, params.Header, bizID)
			if err != nil {
				blog.Errorf("AuditQuery failed, make audit query condition from auth failed, %+v, rid: %s", err, params.ReqID)
				return nil, fmt.Errorf("make audit query condition from auth failed, %+v", err)
			}

			if hasAuthorization == true {
				query.Condition.(map[string]interface{})[common.BKDBOR] = authCondition
				blog.V(5).Infof("AuditQuery, auth condition is: %+v, rid: %s", authCondition, params.ReqID)
				hasAuthorize = hasAuthorization
				break
			}
		}
		if hasAuthorize == false {
			blog.Errorf("AuditQuery failed, user %+v has no authorization on audit, rid: %s", params.User, params.ReqID)
			return nil, auth.NoAuthorizeError
		}
	}

	blog.V(5).Infof("AuditQuery, AuditOperation parameter: %+v, rid: %s", query, params.ReqID)
	return s.Core.AuditOperation().Query(params, query)
}

// InstanceAuditQuery search instance audit logs
// current use case: get host and process related audit log in cmdb web
func (s *Service) InstanceAuditQuery(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	query := metadata.QueryInput{}
	if err := data.MarshalJSONInto(&query); nil != err {
		blog.Errorf("InstanceAuditQuery failed, failed to parse the input (%#v), error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, params.Err.New(common.CCErrCommJSONUnmarshalFailed, err.Error())
	}

	objectID := pathParams("bk_obj_id")
	if len(objectID) == 0 {
		blog.Errorf("InstanceAuditQuery failed, host audit query condition can't be empty, query: %+v, rid: %s", query, params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, "bk_obj_id")
	}

	queryCondition := query.Condition
	if nil == queryCondition {
		blog.Errorf("InstanceAuditQuery failed, host audit query condition can't be empty, query: %+v, rid: %s", query, params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, "condition")
	}

	cond := queryCondition.(map[string]interface{})
	times, ok := cond[common.BKOpTimeField].([]interface{})
	if ok {
		if 2 != len(times) {
			blog.Errorf("InstanceAuditQuery failed, search operation log input params times error, info: %v, rid: %s", times, params.ReqID)
			return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, "op_time")
		}

		cond[common.BKOpTimeField] = common.KvMap{
			"$gte":              times[0],
			"$lte":              times[1],
			CCTimeTypeParseFlag: "1",
		}
	}
	cond[common.BKOwnerIDField] = params.SupplierAccount
	cond[common.BKOpTargetField] = objectID
	query.Condition = cond
	if 0 == query.Limit {
		query.Limit = common.BKDefaultLimit
	}

	// auth: check authorization on instance
	var businessID int64
	bizID, exist := query.Condition.(map[string]interface{})[common.BKAppIDField]
	if exist == true {
		id, err := util.GetInt64ByInterface(bizID)
		if err != nil {
			blog.Errorf("InstanceAuditQuery failed, %s field in query condition but parse int failed, err: %+v, rid: %s", common.BKAppIDField, err, params.ReqID)
			return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
		}
		businessID = id
	}

	instID, exist := queryCondition.(map[string]interface{})["inst_id"]
	if exist == false {
		blog.Errorf("InstanceAuditQuery failed, instance audit query condition condition.ext_key not exist, query: %+v, rid: %s", query, params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, "inst_id")
	}
	instanceID, err := util.GetInt64ByInterface(instID)
	if err != nil {
		blog.Errorf("InstanceAuditQuery failed, instance audit query condition instanceID in condition.ext_key.$in invalid, instanceID: %+v, query: %+v, rid: %s", instID, query, params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, "inst_id")
	}

	opTarget, exist := queryCondition.(map[string]interface{})["op_target"]
	if exist {
		target, ok := opTarget.(string)
		if !ok {
			return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, "op_target")
		}
		if target == "biz" {
			businessID = instanceID
		}
	}

	action := meta.Find
	switch objectID {
	case common.BKInnerObjIDHost:
		err = s.AuthManager.AuthorizeByHostsIDs(params.Context, params.Header, action, instanceID)
	case common.BKInnerObjIDProc:
		err = s.AuthManager.AuthorizeByProcessID(params.Context, params.Header, action, instanceID)
		if err != nil && err == auth.NoAuthorizeError {
			resp, err := s.AuthManager.GenProcessNoPermissionResp(params.Context, params.Header, businessID)
			if err != nil {
				return nil, params.Err.Errorf(common.CCErrTopoGetAppFailed, bizID)
			}
			return resp, auth.NoAuthorizeError
		}
	case common.BKInnerObjIDModule:
		err = s.AuthManager.AuthorizeByModuleID(params.Context, params.Header, action, instanceID)
		if err != nil && err == auth.NoAuthorizeError {
			return s.AuthManager.GenModuleSetNoPermissionResp(), auth.NoAuthorizeError
		}
	case common.BKInnerObjIDSet:
		err = s.AuthManager.AuthorizeBySetID(params.Context, params.Header, action, instanceID)
		if err != nil && err == auth.NoAuthorizeError {
			return s.AuthManager.GenModuleSetNoPermissionResp(), auth.NoAuthorizeError
		}
	case common.BKInnerObjIDApp:
		err = s.AuthManager.AuthorizeByBusinessID(params.Context, params.Header, action, instanceID)
		if err != nil && err == auth.NoAuthorizeError {
			resp, err := s.AuthManager.GenBusinessAuditNoPermissionResp(params.Context, params.Header, businessID)
			if err != nil {
				return nil, params.Err.Error(common.CCErrTopoGetAppFailed)
			}
			return resp, auth.NoAuthorizeError
		}
	default:
		err = s.AuthManager.AuthorizeByInstanceID(params.Context, params.Header, action, objectID, instanceID)
	}
	if err != nil {
		blog.Errorf("InstanceAuditQuery failed, query instance audit log failed, authorization on instance of model %s failed, err: %+v, rid: %s", objectID, err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommAuthorizeFailed)
	}

	blog.V(4).Infof("InstanceAuditQuery failed, AuditOperation parameter: %+v, rid: %s", query, params.ReqID)
	return s.Core.AuditOperation().Query(params, query)
}
