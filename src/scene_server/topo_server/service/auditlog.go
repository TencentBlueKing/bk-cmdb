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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
)

// SearchAuditDict returns all audit types with their name and actions for front-end display
func (s *Service) SearchAuditDict(ctx *rest.Contexts) {
	languageType := common.LanguageType(util.GetLanguage(ctx.Kit.Header))
	if languageType != common.Chinese && languageType != common.English {
		blog.Errorf("can not find language, transform to Chinese, language: %v, rid: %s", languageType, ctx.Kit.Rid)
		languageType = common.Chinese
	}

	dict := metadata.GetAuditDict(languageType)
	if dict == nil {
		blog.Errorf("can not find audit dict, language type: %v, rid: %s", languageType, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommNotFound))
		return
	}

	ctx.RespEntity(dict)
}

// SearchAuditList search audit log list, only contains information for front-end table display
func (s *Service) SearchAuditList(ctx *rest.Contexts) {
	query := metadata.AuditQueryInput{}
	if err := ctx.DecodeInto(&query); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := query.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// the front-end table display fields
	fields := []string{common.BKFieldID, common.BKUser, common.BKResourceTypeField, common.BKActionField,
		common.BKOperationTimeField, common.BKAppIDField, common.BKResourceIDField, common.BKResourceNameField}

	cond := mapstr.MapStr{}
	condition := query.Condition
	if err := condition.Validate(); err != nil {
		blog.Errorf("condition, user and resource_name cannot exist at the same time")
		ctx.RespAutoError(err)
		return
	}

	// parse front-end condition to db search cond
	for _, item := range condition.Condition {
		if item.Operator != querybuilder.OperatorIn && item.Operator != querybuilder.OperatorNotIn {
			if !(item.Field == common.BKResourceNameField && item.Operator == querybuilder.OperatorContains) {
				blog.Errorf("operator invalid, %s wrong, only can be in or not_in (resource_name can use contains)",
					item.Operator)
				ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid,
					"operator only can be in or not_in (resource_name can use contains)"))
				return
			}
		}
		condField, key, err := item.ToMgo()
		if err != nil {
			blog.Errorf("condition invalid, %s wrong, err: %s", key, err.Error())
			ctx.RespAutoError(err)
			return
		}
		cond.Merge(condField)
	}

	if condition.User != "" {
		cond[common.BKUser] = condition.User
	}

	if condition.ResourceName != "" {
		if condition.FuzzyQuery {
			cond[common.BKResourceNameField] = map[string]interface{}{
				common.BKDBLIKE:    condition.ResourceName,
				common.BKDBOPTIONS: "i",
			}
		} else {
			cond[common.BKResourceNameField] = condition.ResourceName
		}
	}

	if condition.ResourceType != "" {
		cond[common.BKResourceTypeField] = condition.ResourceType
	}

	if condition.OperateFrom != "" {
		cond[common.BKOperateFromField] = condition.OperateFrom
	}

	if len(condition.Action) > 0 {
		cond[common.BKActionField] = map[string]interface{}{
			common.BKDBIN: condition.Action,
		}
	}

	if condition.BizID != 0 {
		cond[common.BKAppIDField] = condition.BizID
	}

	if condition.ResourceID != nil {
		cond[common.BKResourceIDField] = condition.ResourceID
	}

	if condition.ObjID != "" {
		switch condition.ResourceType {
		case metadata.ModelInstanceRes:
			cond[common.BKOperationDetailField+"."+common.BKObjIDField] = condition.ObjID
		case metadata.InstanceAssociationRes:
			cond[common.BKOperationDetailField+"."+"src_obj_id"] = condition.ObjID
		default:
			blog.Errorf("unsupported resource type %s when query with object id", condition.ResourceType)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKResourceTypeField))
			return
		}
	}

	// parse operation start time and end time from string to time condition
	timeCond, err := parseOperationTimeCondition(ctx.Kit, condition.OperationTime)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(timeCond) != 0 {
		cond[common.BKOperationTimeField] = timeCond
	}

	// parse audit type condition by category and audit type condition
	auditTypeCond, notMatch := parseAuditTypeCondition(ctx.Kit, condition)
	if notMatch {
		ctx.RespEntity(map[string]interface{}{"count": 0, "info": []interface{}{}})
		return
	}

	if auditTypeCond != nil {
		cond[common.BKAuditTypeField] = auditTypeCond
	}

	auditQuery := metadata.QueryCondition{
		Condition: cond,
		Fields:    fields,
		Page:      query.Page,
	}
	blog.V(5).Infof("AuditQuery, AuditOperation auditQuery: %+v, rid: %s", auditQuery, ctx.Kit.Rid)

	ctx.SetReadPreference(common.SecondaryPreferredMode)
	rsp, err := s.Engine.CoreAPI.CoreService().Audit().SearchAuditLog(ctx.Kit.Ctx, ctx.Kit.Header, auditQuery)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithCount(rsp.Count, rsp.Info)
}

// SearchAuditDetail search audit log detail by id
func (s *Service) SearchAuditDetail(ctx *rest.Contexts) {
	query := metadata.AuditDetailQueryInput{}
	if err := ctx.DecodeInto(&query); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := query.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond := make(map[string]interface{})
	cond[common.BKFieldID] = map[string]interface{}{
		common.BKDBIN: query.IDs,
	}

	auditDetailQuery := metadata.QueryCondition{
		Condition: cond,
	}
	blog.V(5).Infof("AuditDetailQuery, AuditOperation auditDetailQuery: %+v, rid: %s",
		auditDetailQuery, ctx.Kit.Rid)

	ctx.SetReadPreference(common.SecondaryPreferredMode)
	rsp, err := s.Engine.CoreAPI.CoreService().Audit().SearchAuditLog(ctx.Kit.Ctx, ctx.Kit.Header, auditDetailQuery)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(rsp.Info) == 0 {
		blog.Errorf("get no audit log detail, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKFieldID))
	}
	ctx.RespEntity(rsp.Info)
}

func parseOperationTimeCondition(kit *rest.Kit, operationTime metadata.OperationTimeCondition) (map[string]interface{}, error) {
	timeCond := make(map[string]interface{})

	if len(operationTime.Start) != 0 {
		timeCond[common.BKDBGTE] = operationTime.Start
	}

	if len(operationTime.End) != 0 {
		timeCond[common.BKDBLTE] = operationTime.End
	}

	return timeCond, nil
}

func parseAuditTypeCondition(kit *rest.Kit, condition metadata.AuditQueryCondition) (interface{}, bool) {
	if condition.Category != "" {
		auditTypes := metadata.GetAuditTypesByCategory(condition.Category)
		if condition.AuditType != "" {
			flag := false
			if condition.AuditType != metadata.HostType || condition.Category != "resource" {
				for _, audit := range auditTypes {
					if condition.AuditType == audit {
						flag = true
						break
					}
				}

				// audit type and category not match, result is empty
				if !flag {
					return nil, true
				}
			}
		} else {
			return map[string]interface{}{
				common.BKDBIN: auditTypes,
			}, false
		}
	}

	if condition.AuditType != "" {
		return condition.AuditType, false
	}

	return nil, false
}

// SearchInstAudit search instance audit, online allow front-end to use
// 前端在资源池内查看实例的变更记录需要针对当前用户的实例查询权限进行鉴权，原有审计查询接口鉴权条件为操作审计权限，需要进行区别
func (s *Service) SearchInstAudit(ctx *rest.Contexts) {
	query := new(metadata.InstAuditQueryInput)
	if err := ctx.DecodeInto(query); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := query.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	isMainline, err := s.Logics.AssociationOperation().IsMainlineObject(ctx.Kit, query.Condition.ObjID)
	if err != nil {
		blog.Errorf("check if object is mainline object failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if isMainline && query.Condition.BizID <= 0 {
		blog.Errorf("search mainline object audit must provide bizID, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	cond, err := buildInstAuditCondition(ctx, query.Condition)
	if err != nil {
		blog.Errorf("build audit condition failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	fields := make([]string, 0)
	if !query.WithDetail {
		fields = []string{common.BKFieldID, common.BKUser, common.BKResourceTypeField, common.BKActionField,
			common.BKOperationTimeField, common.BKAppIDField, common.BKResourceIDField, common.BKResourceNameField}
	}

	auditQuery := metadata.QueryCondition{Condition: cond, Fields: fields, Page: query.Page}

	ctx.SetReadPreference(common.SecondaryPreferredMode)
	rsp, err := s.Engine.CoreAPI.CoreService().Audit().SearchAuditLog(ctx.Kit.Ctx, ctx.Kit.Header, auditQuery)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithCount(rsp.Count, rsp.Info)
}

func buildInstAuditCondition(ctx *rest.Contexts, query metadata.InstAuditCondition) (mapstr.MapStr, error) {

	cond := mapstr.New()
	// BizID用于校验当前主线实例的权限，查询条件不应设置业务id，这样会导致主线实例的审计信息返回不全
	if query.User != "" {
		cond[common.BKUser] = query.User
	}

	if query.ResourceID != nil {
		cond[common.BKResourceIDField] = query.ResourceID
	}

	if query.ResourceName != "" {
		cond[common.BKResourceNameField] = query.ResourceName
	}

	if query.ResourceType != "" {
		cond[common.BKResourceTypeField] = query.ResourceType
	}

	if len(query.Action) > 0 {
		cond[common.BKActionField] = map[string]interface{}{common.BKDBIN: query.Action}
	}

	switch query.ResourceType {
	case metadata.InstanceAssociationRes:
		cond[common.BKOperationDetailField+"."+"src_obj_id"] = query.ObjID
	case metadata.ModelInstanceRes:
		cond[common.BKOperationDetailField+"."+common.BKObjIDField] = query.ObjID
	case metadata.BusinessRes, metadata.BizSetRes, metadata.HostRes:
		// host, biz and biz set auditlog not need bk_obj_id or operation_detail to select
		break
	default:
		blog.Errorf("unsupported resource type %s when query with object id", query.ResourceType)
		return mapstr.MapStr{}, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKResourceTypeField)
	}

	if len(query.ID) != 0 {
		cond[common.BKFieldID] = mapstr.MapStr{common.BKDBIN: query.ID}
	}

	timeCond, err := parseOperationTimeCondition(ctx.Kit, query.OperationTime)
	if err != nil {
		blog.Errorf("parse operation time condition failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		return mapstr.MapStr{}, err
	}

	if len(timeCond) != 0 {
		cond[common.BKOperationTimeField] = timeCond
	}

	return cond, nil
}
