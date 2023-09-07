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
)

type objectUniqueAuditLog struct {
	audit
}

// NewObjectUniqueAuditLog new object unique auditLog
func NewObjectUniqueAuditLog(clientSet coreservice.CoreServiceClientInterface) *objectUniqueAuditLog {
	return &objectUniqueAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}

// GenerateAuditLog generate audit of model unique, if data is nil, will auto get current model unique data by id.
func (h *objectUniqueAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter, id int64,
	data *metadata.ObjectUnique) (*metadata.AuditLog, error) {
	kit := parameter.kit

	if data == nil {
		// get current model unique data by id.
		rsp, err := h.clientSet.Model().ReadModelAttrUnique(kit.Ctx, kit.Header,
			metadata.QueryCondition{Condition: mapstr.MapStr{common.BKFieldID: id}})
		if err != nil {
			blog.Errorf("generate audit log of model unique failed, failed to read model unique, err: %v, rid: %s", err,
				kit.Rid)
			return nil, err
		}

		if len(rsp.Info) <= 0 {
			blog.Errorf("generate audit log of model unique failed, failed to read model unique, err: %v, rid: %s",
				kit.CCError.CCError(common.CCErrorModelNotFound), kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrorModelNotFound)
		}

		data = &rsp.Info[0]
	}

	return &metadata.AuditLog{
		AuditType:       metadata.ModelType,
		ResourceType:    metadata.ModelUniqueRes,
		Action:          parameter.action,
		ResourceID:      id,
		OperateFrom:     parameter.operateFrom,
		OperationDetail: &metadata.GenericOpDetail{Data: data},
	}, nil
}
