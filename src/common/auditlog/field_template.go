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
	"configcenter/src/common/metadata"
)

// FieldTemplateSettingAuditLog is audit log handler for field template.
type FieldTemplateSettingAuditLog struct {
	// audit base audit handler.
	audit
}

// NewFieldTemplateSettingAuditLog creates a new field template object.
func NewFieldTemplateSettingAuditLog(clientSet coreservice.CoreServiceClientInterface) *PlatFormSettingAuditLog {
	return &PlatFormSettingAuditLog{audit: audit{clientSet: clientSet}}
}

// GenerateAuditLog generate audit of model, if data is nil, will auto get current model data by id.
// todo：待产品确定审计格式
func (h *FieldTemplateSettingAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter, objID string,
	data *metadata.Object) (*metadata.AuditLog, error) {
	kit := parameter.kit
	if data == nil {
		// get current model data by id.
		query := mapstr.MapStr{metadata.ModelFieldObjectID: objID}
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

	return &metadata.AuditLog{
		AuditType:    metadata.FieldTemplateType,
		ResourceType: metadata.FieldTemplateRes,
		Action:       parameter.action,
		ResourceID:   objID,
		ResourceName: data.ObjectName,
		OperateFrom:  parameter.operateFrom,
		OperationDetail: &metadata.BasicOpDetail{
			Details: parameter.NewBasicContent(data.ToMapStr()),
		},
	}, nil
}
