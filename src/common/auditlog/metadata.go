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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// generateAuditCommonParameter include common parameter for generate audit log.
type generateAuditCommonParameter struct {
	kit          *rest.Kit
	action       metadata.ActionType
	operateFrom  metadata.OperateFromType
	updateFields map[string]interface{}
}

func NewGenerateAuditCommonParameter(kit *rest.Kit, action metadata.ActionType) *generateAuditCommonParameter {
	return &generateAuditCommonParameter{
		kit:    kit,
		action: action,
	}
}

func (a *generateAuditCommonParameter) WithOperateFrom(operateFrom metadata.OperateFromType) *generateAuditCommonParameter {
	a.operateFrom = operateFrom
	return a
}

func (a *generateAuditCommonParameter) WithUpdateFields(updateFields map[string]interface{}) *generateAuditCommonParameter {
	a.updateFields = updateFields
	return a
}

// NewBasicContent get basicContent by data and self.
func (a *generateAuditCommonParameter) NewBasicContent(data map[string]interface{}) *metadata.BasicContent {
	var basicDetail *metadata.BasicContent
	switch a.action {
	case metadata.AuditCreate:
		basicDetail = &metadata.BasicContent{
			CurData: data,
		}
	case metadata.AuditDelete:
		basicDetail = &metadata.BasicContent{
			PreData: data,
		}
	case metadata.AuditUpdate:
		basicDetail = &metadata.BasicContent{
			PreData:      data,
			UpdateFields: a.updateFields,
		}
	}

	return basicDetail
}
