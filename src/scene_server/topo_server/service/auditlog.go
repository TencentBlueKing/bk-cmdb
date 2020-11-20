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
	"configcenter/src/common/metadata"
)

// SearchAuditDict returns all audit types with their name and actions for front-end display
func (s *Service) SearchAuditDict(ctx *rest.Contexts) {
	ctx.RespEntity(metadata.GetAuditDict())
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

	cond := make(map[string]interface{})
	condition := query.Condition

	// parse front-end condition to db search cond
	if condition.ResourceType != "" {
		cond[common.BKResourceTypeField] = condition.ResourceType
	}

	if condition.User != "" {
		cond[common.BKUser] = condition.User
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

	if condition.ResourceName != "" {
		cond[common.BKResourceNameField] = map[string]interface{}{
			common.BKDBLIKE: condition.ResourceName,
		}
	}

	if condition.ObjID != "" {
		cond[common.BKOperationDetailField+"."+common.BKObjIDField] = condition.ObjID
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
	count, list, err := s.Core.AuditOperation().SearchAuditList(ctx.Kit, auditQuery)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithCount(count, list)
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
	blog.V(5).Infof("AuditDetailQuery, AuditOperation auditDetailQuery: %+v, rid: %s", auditDetailQuery, ctx.Kit.Rid)

	ctx.SetReadPreference(common.SecondaryPreferredMode)
	list, err := s.Core.AuditOperation().SearchAuditDetail(ctx.Kit, auditDetailQuery)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(list)
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
