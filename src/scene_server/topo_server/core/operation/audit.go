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

package operation

import (
	"configcenter/src/apimachinery"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

type AuditOperationInterface interface {
	SearchAuditList(kit *rest.Kit, query metadata.QueryCondition) (int64, []metadata.AuditLogBasicInfo, error)
}

// NewAuditOperation create a new inst operation instance
func NewAuditOperation(client apimachinery.ClientSetInterface) AuditOperationInterface {
	return &audit{
		clientSet: client,
	}
}

type audit struct {
	clientSet apimachinery.ClientSetInterface
}

func (a *audit) SearchAuditList(kit *rest.Kit, query metadata.QueryCondition) (int64, []metadata.AuditLogBasicInfo, error) {
	rsp, err := a.clientSet.CoreService().Audit().SearchAuditLog(kit.Ctx, kit.Header, query)
	if nil != err {
		blog.ErrorJSON("search audit log list failed, error: %s, query: %s, rid: %s", err.Error(), query, kit.Rid)
		return 0, nil, err
	}

	auditList := make([]metadata.AuditLogBasicInfo, len(rsp.Data.Info))
	for index, auditLog := range rsp.Data.Info {
		auditInfo := metadata.AuditLogBasicInfo{
			ID:            auditLog.ID,
			User:          auditLog.User,
			ResourceType:  auditLog.ResourceType,
			Action:        auditLog.Action,
			OperationTime: auditLog.OperationTime,
		}

		switch auditLog.OperationDetail.WithName() {
		case "BasicDetail":
			operationDetail := auditLog.OperationDetail.(*metadata.BasicOpDetail)
			auditInfo.ResourceID = operationDetail.ResourceID
			auditInfo.ResourceName = operationDetail.ResourceName
			auditInfo.BusinessID = operationDetail.BusinessID

		case "InstanceOpDetail":
			operationDetail := auditLog.OperationDetail.(*metadata.InstanceOpDetail)
			auditInfo.ResourceID = operationDetail.ResourceID
			auditInfo.ResourceName = operationDetail.ResourceName
			auditInfo.BusinessID = operationDetail.BusinessID

		case "InstanceAssociationOpDetail":
			operationDetail := auditLog.OperationDetail.(*metadata.InstanceAssociationOpDetail)
			auditInfo.ResourceID = operationDetail.SourceInstanceID
			auditInfo.ResourceName = operationDetail.SourceInstanceName

		case "ModelAssociationOpDetail":
			operationDetail := auditLog.OperationDetail.(*metadata.ModelAssociationOpDetail)
			auditInfo.ResourceID = operationDetail.AssociationOpDetail.SourceModelID
			auditInfo.ResourceName = operationDetail.SourceModelName

		case "HostTransferOpDetail":
			operationDetail := auditLog.OperationDetail.(*metadata.HostTransferOpDetail)
			auditInfo.ResourceID = operationDetail.HostID
			auditInfo.ResourceName = operationDetail.HostInnerIP
			auditInfo.BusinessID = operationDetail.BusinessID

		case "ModelAttrDetail":
			operationDetail := auditLog.OperationDetail.(*metadata.ModelAttrOpDetail)
			auditInfo.ResourceID = operationDetail.ResourceID
			auditInfo.ResourceName = operationDetail.ResourceName
			auditInfo.BusinessID = operationDetail.BusinessID
		}
		auditList[index] = auditInfo
	}

	return rsp.Data.Count, auditList, nil
}
