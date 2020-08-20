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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/coccyx/timeparser"
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
	opDetailPrefix := common.BKOperationDetailField + "."
	fields := []string{common.BKFieldID, common.BKUser, common.BKResourceTypeField, common.BKActionField,
		common.BKOperationTimeField, opDetailPrefix + common.BKAppIDField}

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
		cond[opDetailPrefix+common.BKAppIDField] = condition.BizID
	}

	// parse operation time from string array to start time and end time
	if condition.OperationTime != nil && len(condition.OperationTime) > 0 {
		times := condition.OperationTime
		if 2 != len(times) {
			blog.Errorf("search operation log input params times error, info: %v, rid: %s", times, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKOperationTimeField))
			return
		}

		startTime, err := timeparser.TimeParserInLocation(times[0], time.Local)
		if nil != err {
			blog.Errorf("parse start time failed, error: %s, time: %s, rid: %s", err, times[0], ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKOperationTimeField))
			return
		}

		endTime, err := timeparser.TimeParserInLocation(times[1], time.Local)
		if nil != err {
			blog.Errorf("parse end time failed, error: %s, time: %s, rid: %s", err, times[1], ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKOperationTimeField))
			return
		}

		cond[common.BKOperationTimeField] = map[string]interface{}{
			common.BKDBGTE: startTime.Local(),
			common.BKDBLTE: endTime.Local(),
		}
	}

	andCond := make([]map[string]interface{}, 0)

	// parse resource id and name condition by resource type, use different fields for different resource type
	resourceID := condition.ResourceID
	resourceNameCond := map[string]interface{}{
		common.BKDBLIKE: condition.ResourceName,
	}

	resourceIDField := opDetailPrefix + common.BKResourceIDField
	srcInstIDField := opDetailPrefix + "src_instance_id"
	srcModelIDField := opDetailPrefix + "src_model_id"
	targetInstIDField := opDetailPrefix + "target_instance_id"
	targetModelIDField := opDetailPrefix + "target_model_id"

	resourceNameField := opDetailPrefix + common.BKResourceNameField
	srcInstNameField := opDetailPrefix + "src_instance_name"
	srcModelNameField := opDetailPrefix + "src_model_name"
	targetInstNameField := opDetailPrefix + "target_instance_name"
	targetModelNameField := opDetailPrefix + "target_model_name"

	switch condition.ResourceType {
	case "":
		fields = append(fields, resourceIDField, srcInstIDField, srcModelIDField, resourceNameField, srcInstNameField, srcModelNameField)

		if condition.ResourceID != 0 {
			andCond = append(andCond, map[string]interface{}{common.BKDBOR: []map[string]interface{}{
				{resourceIDField: resourceID},
				{srcInstIDField: resourceID},
				{targetInstIDField: resourceID},
				{srcModelIDField: resourceID},
				{targetModelIDField: resourceID},
			}})
		}

		if condition.ResourceName != "" {
			andCond = append(andCond, map[string]interface{}{common.BKDBOR: []map[string]interface{}{
				{resourceNameField: resourceNameCond},
				{srcInstNameField: resourceNameCond},
				{targetInstNameField: resourceNameCond},
				{srcModelNameField: resourceNameCond},
				{targetModelNameField: resourceNameCond},
			}})
		}

	case metadata.InstanceAssociationRes:
		fields = append(fields, srcInstIDField, srcInstNameField)

		if condition.ResourceID != 0 {
			andCond = append(andCond, map[string]interface{}{common.BKDBOR: []map[string]interface{}{
				{srcInstIDField: resourceID},
				{targetInstIDField: resourceID},
			}})
		}

		if condition.ResourceName != "" {
			andCond = append(andCond, map[string]interface{}{common.BKDBOR: []map[string]interface{}{
				{srcInstNameField: resourceNameCond},
				{targetInstNameField: resourceNameCond},
			}})
		}

	case metadata.ModelAssociationRes:
		fields = append(fields, srcModelIDField, srcModelNameField)

		if condition.ResourceID != 0 {
			andCond = append(andCond, map[string]interface{}{common.BKDBOR: []map[string]interface{}{
				{srcModelIDField: resourceID},
				{targetModelIDField: resourceID},
			}})
		}

		if condition.ResourceName != "" {
			andCond = append(andCond, map[string]interface{}{common.BKDBOR: []map[string]interface{}{
				{srcModelNameField: resourceNameCond},
				{targetModelNameField: resourceNameCond},
			}})
		}

	default:
		fields = append(fields, resourceIDField, resourceNameField)
		if condition.ResourceID != 0 {
			cond[resourceIDField] = resourceID
		}
		if condition.ResourceName != "" {
			cond[resourceNameField] = resourceNameCond
		}
	}

	// parse audit type condition by category and audit type condition
	auditTypeCond := make(map[string]interface{})
	if condition.AuditType != "" {
		auditTypeCond[common.BKAuditTypeField] = condition.AuditType
	}
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

				// audit type and category not match
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

		bizID, err := s.getDefaultBizID(ctx.Kit)
		if err != nil {
			blog.Errorf("get default biz id failed, err: %s", err.Error())
			ctx.RespAutoError(err)
			return
		}
		switch condition.Category {
		case "business":
			cond[common.BKAuditTypeField] = auditTypeCond[common.BKAuditTypeField]
			andCond = append(andCond, map[string]interface{}{
				common.BKActionField:                 map[string]interface{}{common.BKDBNE: metadata.AuditAssignHost},
				opDetailPrefix + common.BKAppIDField: map[string]interface{}{common.BKDBNIN: []int64{0, bizID}},
			})

		case "resource":
			andCond = append(andCond, map[string]interface{}{common.BKDBOR: []map[string]interface{}{
				auditTypeCond,
				{
					common.BKAuditTypeField: metadata.HostType,
					common.BKActionField:    map[string]interface{}{common.BKDBEQ: metadata.AuditAssignHost},
				},
				{
					common.BKAuditTypeField:              metadata.HostType,
					opDetailPrefix + common.BKAppIDField: map[string]interface{}{common.BKDBIN: []int64{0, bizID}},
				},
			}})
		default:
			cond[common.BKAuditTypeField] = auditTypeCond[common.BKAuditTypeField]
		}
	}

	if len(andCond) > 0 {
		cond[common.BKDBAND] = andCond
	}

	auditQuery := metadata.QueryCondition{
		Condition: cond,
		Fields:    fields,
		Page:      query.Page,
	}
	blog.V(5).Infof("AuditQuery, AuditOperation auditQuery: %+v, rid: %s", auditQuery, ctx.Kit.Rid)
	count, list, err := s.Core.AuditOperation().SearchAuditList(ctx.Kit, auditQuery)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithCount(count, list)
}

var defaultBizID int64

func (s *Service) getDefaultBizID(kit *rest.Kit) (int64, error) {
	if defaultBizID != 0 {
		return defaultBizID, nil
	}

	biz, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, &metadata.QueryCondition{
		Fields: []string{common.BKAppIDField},
		Page:   metadata.BasePage{Limit: 1},
		Condition: map[string]interface{}{
			common.BKDefaultField: common.DefaultAppFlag,
		},
	})

	if err != nil {
		return 0, err
	}

	if len(biz.Data.Info) == 0 {
		return 0, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	defaultBizID, err = util.GetInt64ByInterface(biz.Data.Info[0][common.BKAppIDField])
	if err != nil {
		return 0, err
	}

	return defaultBizID, nil
}
