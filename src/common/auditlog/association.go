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
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// instanceAuditLog provides methods to generate and save instance association audit log
type instanceAssociationAuditLog struct {
	audit
}

func (a *instanceAssociationAuditLog) GenerateAuditLog(kit *rest.Kit, data *metadata.InstAsst, id int64, action metadata.ActionType) (
	metadata.AuditLog, errors.CCErrorCoder) {

	if data == nil {
		cond := metadata.QueryCondition{Condition: map[string]interface{}{common.BKFieldID: id}}
		result, err := a.clientSet.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header, &cond)
		if err != nil {
			blog.Errorf("generate inst asst audit log failed, get instance association failed, err: %v, rid: %s", err, kit.Rid)
			return metadata.AuditLog{}, kit.CCError.CCError(common.CCErrAuditTakeSnapshotFailed)
		}

		if len(result.Data.Info) == 0 || len(result.Data.Info) > 1 {
			blog.Errorf("generate inst asst audit log failed, get instance association by id(%d) get none or multiple result, rid: %s", id, kit.Rid)
			return metadata.AuditLog{}, kit.CCError.CCError(common.CCErrAuditTakeSnapshotFailed)
		}

		data = &result.Data.Info[0]
	}

	srcInstName, err := a.getInstNameByID(kit, data.ObjectID, data.InstID)
	if err != nil {
		return metadata.AuditLog{}, kit.CCError.CCError(common.CCErrAuditTakeSnapshotFailed)
	}

	targetInstName, err := a.getInstNameByID(kit, data.AsstObjectID, data.AsstInstID)
	if err != nil {
		return metadata.AuditLog{}, kit.CCError.CCError(common.CCErrAuditTakeSnapshotFailed)
	}

	return metadata.AuditLog{
		AuditType:    metadata.ModelInstanceType,
		ResourceType: metadata.InstanceAssociationRes,
		Action:       action,
		OperationDetail: &metadata.InstanceAssociationOpDetail{
			AssociationOpDetail: metadata.AssociationOpDetail{
				AssociationID:   data.ObjectAsstID,
				AssociationKind: data.AssociationKindID,
				SourceModelID:   data.ObjectID,
				TargetModelID:   data.AsstObjectID,
			},
			SourceInstanceID:   data.InstID,
			SourceInstanceName: srcInstName,
			TargetInstanceID:   data.AsstInstID,
			TargetInstanceName: targetInstName,
		},
	}, nil
}

func NewInstanceAssociationAudit(clientSet apimachinery.ClientSetInterface) *instanceAssociationAuditLog {
	return &instanceAssociationAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
