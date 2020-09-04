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
	"configcenter/src/common/metadata"
)

type hostModuleAuditLog struct {
	audit
}

// GenerateAuditLog generate audit log of host module relate.
func (h *hostModuleAuditLog) GenerateAuditLog(action metadata.ActionType, hostID, bizID int64, hostIP string,
	preData, curData metadata.HostBizTopo) *metadata.AuditLog {
	var auditLog = metadata.AuditLog{
		AuditType:    metadata.HostType,
		ResourceType: metadata.HostRes,
		Action:       action,
		BusinessID:   bizID,
		ResourceID:   hostID,
		ResourceName: hostIP,
		OperationDetail: &metadata.HostTransferOpDetail{
			PreData: preData,
			CurData: curData,
		},
	}

	return &auditLog
}

func NewHostModuleAudit(clientSet coreservice.CoreServiceClientInterface) *hostModuleAuditLog {
	return &hostModuleAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
