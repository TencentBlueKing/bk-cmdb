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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

type hostAuditLog struct {
	audit
}

// GenerateAuditLog generate audit log of host, if data is nil, will auto get host current data by hostID,
// meanwhile update hostIP if hostIP is "".
func (h *hostAuditLog) GenerateAuditLog(kit *rest.Kit, action metadata.ActionType, hostID, bizID int64, innerIP string,
	OperateFrom metadata.OperateFromType, data, updateFields map[string]interface{}) (*metadata.AuditLog, error) {
	return h.generateAuditLog(kit, action, hostID, bizID, innerIP, OperateFrom, data, updateFields)
}

// GenerateAuditLogByHostIDGetBizID generate audit log of host, auto get bizID by hostID and host topology relate,
// if data is nil, will auto get host current data by hostID, meanwhile update hostIP if hostIP is "".
func (h *hostAuditLog) GenerateAuditLogByHostIDGetBizID(kit *rest.Kit, action metadata.ActionType, hostID int64,
	innerIP string, OperateFrom metadata.OperateFromType, data, updateFields map[string]interface{}) (*metadata.AuditLog, error) {
	// get bizID by hostID and host topology related.
	bizID, err := h.getBizIDByHostID(kit, hostID)
	if err != nil {
		blog.Errorf("generate host audit log failed, failed to get bizID by hostID, hostID: %d, innerIP: %s, err: %v, rid: %s",
			hostID, innerIP, err, kit.Rid)
		return nil, err
	}

	return h.generateAuditLog(kit, action, hostID, bizID, innerIP, OperateFrom, data, updateFields)
}

func (h *hostAuditLog) generateAuditLog(kit *rest.Kit, action metadata.ActionType, hostID, bizID int64, hostIP string,
	OperateFrom metadata.OperateFromType, data, updateFields map[string]interface{}) (*metadata.AuditLog, error) {
	// get current data and inner ip by hostID
	if data == nil {
		var innerIP string
		var err error

		data, innerIP, err = h.getHostInstanceDetailByHostID(kit, hostID)
		if err != nil {
			blog.Errorf("generate host audit log failed, failed to get host instance by bizID, hostID: %d, err: %v, rid: %s",
				hostID, err, kit.Rid)
			return nil, err
		}

		if innerIP != "" {
			hostIP = innerIP
		}
	}

	var basicDetail *metadata.BasicContent
	switch action {
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
			UpdateFields: updateFields,
		}
	}

	var auditLog = &metadata.AuditLog{
		AuditType:    metadata.HostType,
		ResourceType: metadata.HostRes,
		Action:       action,
		BusinessID:   bizID,
		ResourceID:   hostID,
		ResourceName: hostIP,
		OperateFrom:  OperateFrom,
		OperationDetail: &metadata.InstanceOpDetail{
			BasicOpDetail: metadata.BasicOpDetail{
				Details: basicDetail,
			},
			ModelID: common.BKInnerObjIDHost,
		},
	}
	return auditLog, nil
}

func NewHostAudit(clientSet coreservice.CoreServiceClientInterface) *hostAuditLog {
	return &hostAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
