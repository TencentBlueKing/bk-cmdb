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
	"fmt"

	"configcenter/src/auth"
	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// AuditQuery search audit logs
func (s *Service) AuditQuery(ctx *rest.Contexts) {
	query := metadata.QueryInput{}
	if err := ctx.DecodeInto(&query); nil != err {
		ctx.RespAutoError(err)
		return
	}

	rawErr := query.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	var businessID int64

	queryCondition := query.Condition
	if nil == queryCondition {
		query.Condition = make(map[string]interface{})
	} else {
		condition := metadata.AuditQueryCondition{}
		js, err := json.Marshal(queryCondition)
		if nil != err {
			ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommJSONMarshalFailed, err.Error()))
			return
		}
		err = json.Unmarshal(js, &condition)
		if err != nil {
			ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommJSONUnmarshalFailed, err.Error()))
			return
		}

		cond := make(map[string]interface{})
		auditTypeCond := make(map[string]interface{})
		if condition.AuditType != "" {
			auditTypeCond[common.BKAuditTypeField] = condition.AuditType
		}
		if condition.User != "" {
			cond[common.BKUser] = condition.User
		}
		if condition.OperateFrom != "" {
			cond[common.BKOperateFromField] = condition.OperateFrom
		}
		if condition.Action != nil && len(condition.Action) > 0 {
			cond[common.BKActionField] = map[string]interface{}{
				common.BKDBIN: condition.Action,
			}
		}

		if condition.OperationTime != nil && len(condition.OperationTime) > 0 {
			times := condition.OperationTime
			if 2 != len(times) {
				blog.Errorf("search operation log input params times error, info: %v, rid: %s", times, ctx.Kit.Rid)
				ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKOperationTimeField))
				return
			}
			cond[common.BKOperationTimeField] = map[string]interface{}{
				common.BKDBGTE:             times[0],
				common.BKDBLTE:             times[1],
				common.BKTimeTypeParseFlag: "1",
			}
		}

		andCond := make([]map[string]interface{}, 0)

		// add auth filter condition
		if condition.BizID != 0 {
			businessID = condition.BizID
			andCond = append(andCond, map[string]interface{}{common.BKOperationDetailField + "." + common.BKAppIDField: businessID})
		}

		if condition.ResourceID != 0 {
			resourceID := condition.ResourceID
			andCond = append(andCond, map[string]interface{}{common.BKDBOR: []map[string]interface{}{
				{common.BKOperationDetailField + "." + common.BKResourceIDField: resourceID},
				{common.BKOperationDetailField + "." + common.BKHostIDField: resourceID},
				{common.BKOperationDetailField + ".src_instance_id": resourceID},
				{common.BKOperationDetailField + ".target_instance_id": resourceID},
			}})
		}

		if condition.ResourceName != "" {
			resourceNameCond := map[string]interface{}{
				common.BKDBLIKE: condition.ResourceName,
			}
			andCond = append(andCond, map[string]interface{}{common.BKDBOR: []map[string]interface{}{
				{common.BKOperationDetailField + "." + common.BKResourceNameField: resourceNameCond},
				{common.BKOperationDetailField + "." + common.BKHostInnerIPField: resourceNameCond},
				{common.BKOperationDetailField + ".src_instance_name": resourceNameCond},
				{common.BKOperationDetailField + ".target_instance_name": resourceNameCond},
			}})
		}

		if condition.Category != "" {
			auditTypes := metadata.GetAuditTypesByCategory(condition.Category)
			if condition.AuditType != "" {
				flag := false
				if condition.AuditType != metadata.HostType || (condition.Category != "business" && condition.Category != "resource") {
					for _, audit := range auditTypes {
						if condition.AuditType == audit {
							flag = true
						}
					}
					if !flag {
						ctx.RespEntity(map[string]interface{}{"count": 0, "info": []interface{}{}})
						return
					}
				}
			} else {
				auditTypeCond[common.BKAuditTypeField] = map[string]interface{}{
					common.BKDBIN: auditTypes,
				}
			}
			biz, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDApp, &metadata.QueryCondition{
				Fields: []string{common.BKAppIDField},
				Page: metadata.BasePage{
					Limit: 1,
					Start: 0,
				},
				Condition: map[string]interface{}{
					common.BKDefaultField: common.DefaultAppFlag,
				},
			})
			if nil != err {
				blog.Errorf("find default biz failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				ctx.RespAutoError(err)
				return
			}
			if len(biz.Data.Info) == 0 {
				blog.Errorf("find default biz get no result, rid: %s", ctx.Kit.Rid)
				ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommBizNotFoundError, "default"))
				return
			}
			defaultBizID := biz.Data.Info[0][common.BKAppIDField]
			switch condition.Category {
			case "business":
				andCond = append(andCond, map[string]interface{}{common.BKDBOR: []map[string]interface{}{
					auditTypeCond,
					{
						common.BKAuditTypeField: metadata.HostType,
						common.BKActionField:    map[string]interface{}{common.BKDBNE: metadata.AuditAssignHost},
						common.BKOperationDetailField + "." + common.BKAppIDField: map[string]interface{}{common.BKDBNE: defaultBizID},
					},
				}})
			case "resource":
				andCond = append(andCond, map[string]interface{}{common.BKDBOR: []map[string]interface{}{
					auditTypeCond,
					{
						common.BKAuditTypeField: metadata.HostType,
						common.BKActionField:    map[string]interface{}{common.BKDBEQ: metadata.AuditAssignHost},
					},
					{
						common.BKAuditTypeField:                                   metadata.HostType,
						common.BKOperationDetailField + "." + common.BKAppIDField: map[string]interface{}{common.BKDBEQ: defaultBizID},
					},
				}})
			default:
				cond[common.BKAuditTypeField] = auditTypeCond[common.BKAuditTypeField]
			}
		}

		labelCond := make([]map[string]interface{}, 0)
		if condition.Label != nil {
			for _, label := range condition.Label {
				labelCond = append(labelCond, map[string]interface{}{
					common.BKLabelField + "." + label: map[string]interface{}{
						common.BKDBExists: true,
						common.BKDBNE:     nil,
					},
				})
			}
		}
		if condition.ResourceType != nil && len(condition.ResourceType) > 0 {
			if len(labelCond) > 0 {
				labelCond = append(labelCond, map[string]interface{}{
					common.BKResourceTypeField: map[string]interface{}{
						common.BKDBIN: condition.ResourceType,
					},
				})
				andCond = append(andCond, map[string]interface{}{
					common.BKDBOR: labelCond,
				})
			} else {
				cond[common.BKResourceTypeField] = map[string]interface{}{
					common.BKDBIN: condition.ResourceType,
				}
			}
		}

		if len(andCond) > 0 {
			cond[common.BKDBAND] = andCond
		}
		query.Condition = cond
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

	rawErr := query.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	objectID := ctx.Request.PathParameter(common.BKObjIDField)
	if len(objectID) == 0 {
		blog.Errorf("InstanceAuditQuery failed, object ID can't be empty, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjIDField))
		return
	}

	queryCondition := query.Condition
	if nil == queryCondition {
		blog.Errorf("InstanceAuditQuery failed, host audit query condition can't be empty, query: %+v, rid: %s", query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "condition"))
		return
	}
	condition := metadata.AuditQueryCondition{}
	js, err := json.Marshal(queryCondition)
	if nil != err {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommJSONMarshalFailed, err.Error()))
		return
	}
	err = json.Unmarshal(js, &condition)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommJSONUnmarshalFailed, err.Error()))
		return
	}

	cond := make(map[string]interface{})
	if condition.User != "" {
		cond[common.BKUser] = condition.User
	}
	if condition.OperateFrom != "" {
		cond[common.BKOperateFromField] = condition.OperateFrom
	}
	if condition.Action != nil && len(condition.Action) > 0 {
		cond[common.BKActionField] = map[string]interface{}{
			common.BKDBIN: condition.Action,
		}
	}

	if condition.OperationTime != nil && len(condition.OperationTime) > 0 {
		times := condition.OperationTime
		if 2 != len(times) {
			blog.Errorf("search operation log input params times error, info: %v, rid: %s", times, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKOperationTimeField))
			return
		}
		cond[common.BKOperationTimeField] = map[string]interface{}{
			common.BKDBGTE:             times[0],
			common.BKDBLTE:             times[1],
			common.BKTimeTypeParseFlag: "1",
		}
	}

	andCond := make([]map[string]interface{}, 0)

	// auth: check authorization on instance
	var businessID int64
	if condition.BizID != 0 {
		businessID = condition.BizID
		andCond = append(andCond, map[string]interface{}{common.BKOperationDetailField + "." + common.BKAppIDField: businessID})
	}

	if condition.ResourceID == 0 {
		blog.Errorf("InstanceAuditQuery failed, instance audit query condition condition.resource_id not exist, query: %+v, rid: %s", query, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKResourceIDField))
		return
	}

	instanceID := condition.ResourceID
	if objectID == common.BKInnerObjIDApp {
		businessID = instanceID
	}
	isMainline, err := s.Core.AssociationOperation().IsMainlineObject(ctx.Kit, objectID)
	if err != nil {
		blog.Errorf("InstanceAuditQuery failed, check if object(%s) is mainline object failed, err: %s, rid: %s", err.Error(), objectID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommHTTPDoRequestFailed, err.Error()))
		return
	}
	orCond := []map[string]interface{}{
		{
			common.BKOperationDetailField + "." + common.BKResourceIDField: instanceID,
			common.BKResourceTypeField:                                     metadata.GetResourceTypeByObjID(objectID, isMainline),
		},
		{
			common.BKOperationDetailField + ".src_instance_id":                instanceID,
			common.BKOperationDetailField + ".basic_asst_detail.src_model_id": objectID,
			common.BKResourceTypeField:                                        metadata.InstanceAssociationRes,
		},
		{
			common.BKOperationDetailField + ".target_instance_id":                instanceID,
			common.BKOperationDetailField + ".basic_asst_detail.target_model_id": objectID,
			common.BKResourceTypeField:                                           metadata.InstanceAssociationRes,
		},
	}
	if objectID == common.BKInnerObjIDHost {
		orCond = append(orCond, map[string]interface{}{
			common.BKOperationDetailField + "." + common.BKHostIDField: instanceID,
			common.BKResourceTypeField:                                 metadata.HostRes,
		})
	}
	andCond = append(andCond, map[string]interface{}{common.BKDBOR: orCond})

	cond[common.BKDBAND] = andCond
	query.Condition = cond

	action := meta.Find
	switch objectID {
	case common.BKInnerObjIDHost:
		err = s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, action, instanceID)
	case common.BKInnerObjIDProc:
		err = s.AuthManager.AuthorizeByProcessID(ctx.Kit.Ctx, ctx.Kit.Header, action, instanceID)
		if err != nil && err == auth.NoAuthorizeError {
			resp, err := s.AuthManager.GenProcessNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, businessID)
			if err != nil {
				ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrTopoGetAppFailed, businessID))
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
