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
	"strconv"

	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/metadata"
)

// DynamicGroupAuditLog is audit log handler for dynamic grouping.
type DynamicGroupAuditLog struct {
	// audit base audit handler.
	audit
}

// NewDynamicGroupAuditLog creates a new DynamicGroupAuditLog object.
func NewDynamicGroupAuditLog(clientSet coreservice.CoreServiceClientInterface) *DynamicGroupAuditLog {
	return &DynamicGroupAuditLog{audit: audit{clientSet: clientSet}}
}

// GenerateAuditLog generates an audit log for dynamic grouping operations.
func (l *DynamicGroupAuditLog) GenerateAuditLog(param *generateAuditCommonParameter, dynamicGroup *metadata.DynamicGroup) ([]metadata.AuditLog, error) {
	if dynamicGroup == nil {
		return make([]metadata.AuditLog, 0), nil
	}

	kit := param.kit
	content := make(map[string]interface{})

	if param.action == metadata.AuditUpdate {
		bizID := strconv.FormatInt(dynamicGroup.AppID, 10)
		resp, err := l.clientSet.Host().GetDynamicGroup(kit.Ctx, bizID, dynamicGroup.ID, kit.Header)
		if err != nil {
			return nil, err
		}
		if !resp.Result {
			return nil, resp.CCError()
		}
		dynamicGroup = &resp.Data
	}

	content[common.BKFieldID] = dynamicGroup.ID
	content[common.BKFieldName] = dynamicGroup.Name
	content[common.BKAppIDField] = dynamicGroup.AppID
	content[common.BKObjIDField] = dynamicGroup.ObjID
	content["info"] = dynamicGroup.Info
	content["create_user"] = dynamicGroup.CreateUser
	content["modify_user"] = dynamicGroup.ModifyUser
	content[common.CreateTimeField] = dynamicGroup.CreateTime
	content[common.LastTimeField] = dynamicGroup.UpdateTime

	logs := []metadata.AuditLog{metadata.AuditLog{
		AuditType:       metadata.DynamicGroupType,
		ResourceType:    metadata.DynamicGroupRes,
		Action:          param.action,
		ResourceID:      dynamicGroup.ID,
		ResourceName:    dynamicGroup.Name,
		BusinessID:      dynamicGroup.AppID,
		OperationDetail: &metadata.BasicOpDetail{Details: param.NewBasicContent(content)},
	}}

	return logs, nil
}
