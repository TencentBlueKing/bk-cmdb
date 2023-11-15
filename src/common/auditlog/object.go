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

package auditlog

import (
	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
)

type objectAuditLog struct {
	audit
}

// GenerateAuditLog generate audit of model, if data is nil, will auto get current model data by id.
func (h *objectAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter, id int64, data *metadata.Object) (
	*metadata.AuditLog, error) {
	kit := parameter.kit
	if data == nil {
		// get current model data by id.
		query := mapstr.MapStr{metadata.ModelFieldID: id}
		rsp, err := h.clientSet.Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: query})
		if err != nil {
			blog.Errorf("generate audit log of model failed, failed to read model, err: %v, rid: %s",
				err.Error(), kit.Rid)
			return nil, err
		}

		if len(rsp.Info) <= 0 {
			blog.Errorf("generate audit log of model failed, failed to read model, err: %s, rid: %s",
				kit.CCError.CCError(common.CCErrorModelNotFound), kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrorModelNotFound)
		}

		data = &rsp.Info[0]
	}

	dataMap, err := mapstruct.Struct2Map(data)
	if err != nil {
		blog.Errorf("convert model(%+v) to map failed, err: %v, rid: %s", data, err, kit.Rid)
		return nil, err
	}

	return &metadata.AuditLog{
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModelRes,
		Action:       parameter.action,
		ResourceID:   id,
		ResourceName: data.ObjectName,
		OperateFrom:  parameter.operateFrom,
		OperationDetail: &metadata.BasicOpDetail{
			Details: parameter.NewBasicContent(dataMap),
		},
	}, nil
}

// GenerateAuditLogForBindingFieldTemplate specific generate audit log function for model binding template scenarios.
func (h *objectAuditLog) GenerateAuditLogForBindingFieldTemplate(parameter *generateAuditCommonParameter,
	objIDs []int64, templateID int64) ([]metadata.AuditLog, error) {

	kit := parameter.kit

	objectLen := len(objIDs)
	objectTmplIDMap := make(map[int64][]int64)

	auditLogs := make([]metadata.AuditLog, 0)
	for start := 0; start < objectLen; start += common.BKMaxPageSize {
		limit := start + common.BKMaxPageSize
		if limit > objectLen {
			limit = objectLen
		}
		query := &metadata.QueryCondition{
			Condition: mapstr.MapStr{
				metadata.ModelFieldID: mapstr.MapStr{
					common.BKDBIN: objIDs[start:limit],
				},
			},
			DisableCounter: true,
		}
		rsp, err := h.clientSet.Model().ReadModel(kit.Ctx, kit.Header, query)
		if err != nil {
			blog.Errorf("failed to read model, cond: %+v, err: %v, rid: %s", query, err, kit.Rid)
			return nil, err
		}

		if len(rsp.Info) <= 0 {
			blog.Errorf("no model founded, cond: %+v, rid: %s", query, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrorModelNotFound)
		}
		if len(rsp.Info) != limit-start {
			blog.Errorf("fetching model data does not meet expectations, cond: %+v, rid: %s", query, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommParamsInvalid)
		}
		// todo: 获取关联关系接口

		for _, data := range rsp.Info {
			obj, err := mapstruct.Struct2Map(data)
			if err != nil {
				blog.Errorf("convert model(%+v) to map failed, err: %v, rid: %s", data, err, kit.Rid)
				return nil, err
			}
			objectTmplIDMap[data.ID] = append(objectTmplIDMap[data.ID], templateID)
			obj[string(metadata.ObjTemplateIDs)] = objectTmplIDMap[data.ID]
			parameter.updateFields = obj
			auditLog := metadata.AuditLog{
				AuditType:    metadata.ModelType,
				ResourceType: metadata.ModelRes,
				Action:       parameter.action,
				ResourceID:   data.ID,
				ResourceName: data.ObjectName,
				OperateFrom:  parameter.operateFrom,
				OperationDetail: &metadata.BasicOpDetail{
					Details: parameter.NewBasicContent(obj),
				},
			}
			auditLogs = append(auditLogs, auditLog)
		}
	}
	return auditLogs, nil
}

// NewObjectAuditLog TODO
func NewObjectAuditLog(clientSet coreservice.CoreServiceClientInterface) *objectAuditLog {
	return &objectAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
