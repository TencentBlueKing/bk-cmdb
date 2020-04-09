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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type ObjectAudit struct {
	auditType    metadata.AuditType
	resourceType metadata.ResourceType
	clientSet    apimachinery.ClientSetInterface
	preData      map[string]interface{}
	curData      map[string]interface{}
	objectID     string
	objectName   string
}

func NewObjectAudit(clientSet apimachinery.ClientSetInterface, resourceType metadata.ResourceType) *ObjectAudit {
	return &ObjectAudit{
		auditType:    metadata.ModelType,
		resourceType: metadata.ModelRes,
		clientSet:    clientSet,
	}
}

func (log *ObjectAudit) WithPrevious(kit *rest.Kit, id int64) errors.CCError {
	preData, err := log.buildLogData(kit, id)
	if err != nil {
		return err
	}
	log.preData = preData
	return nil
}

func (log *ObjectAudit) MakeCurrent(curData map[string]interface{}) errors.CCError {
	log.curData = curData
	return nil
}

func (log *ObjectAudit) WithCurrent(kit *rest.Kit, id int64) errors.CCError {
	curData, err := log.buildLogData(kit, id)
	if err != nil {
		return err
	}
	log.curData = curData
	return nil
}

func (log *ObjectAudit) SaveAuditLog(kit *rest.Kit, auditAction metadata.ActionType) errors.CCError {
	auditLog := metadata.AuditLog{
		AuditType:    log.auditType,
		ResourceType: log.resourceType,
		Action:       auditAction,
		OperationDetail: &metadata.ModelOpDetail{
			ObjectID:   log.objectID,
			ObjectName: log.objectName,
			Details: &metadata.BasicContent{
				PreData: log.preData,
				CurData: log.curData,
			},
		},
	}
	auditResult, err := log.clientSet.CoreService().Audit().SaveAuditLog(kit.Ctx, kit.Header, auditLog)
	if err != nil {
		blog.ErrorJSON("SaveAuditLog %s %s audit log failed, err: %s, result: %+v,rid:%s", auditAction, auditAction, err, auditResult, kit.Rid)
		return kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}
	if auditResult.Result != true {
		blog.ErrorJSON("SaveAuditLog %s %s audit log failed, err: %s, result: %s,rid:%s", auditAction, auditAction, err, auditResult, kit.Rid)
		return kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

func (log *ObjectAudit) buildLogData(kit *rest.Kit, ID int64) (map[string]interface{}, errors.CCError) {
	query := mapstr.MapStr{"id": ID}
	rsp, err := log.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: query})
	if err != nil {
		return nil, err
	}
	if rsp.Result != true {
		return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}
	if len(rsp.Data.Info) <= 0 {
		return nil, kit.CCError.CCError(common.CCErrorModelNotFound)
	}
	log.objectID = rsp.Data.Info[0].Spec.ObjectID
	log.objectName = rsp.Data.Info[0].Spec.ObjectName
	rspData := rsp.Data.Info[0].Spec.ToMapStr()
	return rspData, nil
}
