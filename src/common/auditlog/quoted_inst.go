/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package auditlog

import (
	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type quotedInstAuditLog struct {
	audit
}

// GenerateAuditLog generate audit log of quoted instance.
func (h *quotedInstAuditLog) GenerateAuditLog(params *generateAuditCommonParameter, objID, srcObjID, attrID string,
	data []mapstr.MapStr) ([]metadata.AuditLog, error) {

	auditLogs := make([]metadata.AuditLog, len(data))
	kit := params.kit

	for index, inst := range data {
		id, err := util.GetInt64ByInterface(inst[common.BKFieldID])
		if err != nil {
			blog.Errorf("parse quoted inst id failed, err: %v, id: %+v, rid: %s", err, inst[common.BKFieldID], kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID)
		}

		auditLog := metadata.AuditLog{
			AuditType:    metadata.QuotedInstType,
			ResourceType: metadata.QuotedInst,
			Action:       params.action,
			ResourceID:   id,
			OperateFrom:  params.operateFrom,
			ResourceName: util.GetStrByInterface(inst[metadata.GetInstNameFieldName(objID)]),
			OperationDetail: &metadata.QuotedInstOpDetail{
				BasicOpDetail: metadata.BasicOpDetail{
					Details: params.NewBasicContent(inst),
				},
				ObjID:    objID,
				SrcObjID: srcObjID,
				AttrID:   attrID,
			},
		}
		auditLogs[index] = auditLog
	}

	return auditLogs, nil
}

// NewQuotedInstAuditLog new audit log utility for quoted inst
func NewQuotedInstAuditLog(clientSet coreservice.CoreServiceClientInterface) *quotedInstAuditLog {
	return &quotedInstAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
