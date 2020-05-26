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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type ObjAuditLog interface {
	buildSnapshotForPre() ObjAuditLog
	buildSnapshotForCur() ObjAuditLog
	SaveAuditLog(metadata.ActionType)
}

type ObjectAudit struct {
	kit          *rest.Kit
	clientSet    apimachinery.ClientSetInterface
	auditType    metadata.AuditType
	resourceType metadata.ResourceType
	id           int64
	bkObjID      string
	bkObjName    string
	preData      metadata.Object
	curData      metadata.Object
}

func NewObjectAudit(kit *rest.Kit, clientSet apimachinery.ClientSetInterface, ID int64) ObjAuditLog {
	return &ObjectAudit{
		kit:          kit,
		clientSet:    clientSet,
		auditType:    metadata.ModelType,
		resourceType: metadata.ModelRes,
		id:           ID,
	}
}

func (log *ObjectAudit) SaveAuditLog(auditAction metadata.ActionType) {
	preData := log.preData.ToMapStr()
	curData := log.curData.ToMapStr()
	switch auditAction {
	case metadata.AuditDelete:
		curData = nil
	case metadata.AuditCreate:
		preData = nil
	case metadata.AuditUpdate:
		//do nothing
	}
	auditLog := metadata.AuditLog{
		AuditType:    log.auditType,
		ResourceType: log.resourceType,
		Action:       auditAction,
		OperationDetail: &metadata.BasicOpDetail{
			ResourceID:   log.id,
			ResourceName: log.bkObjName,
			Details: &metadata.BasicContent{
				PreData: preData,
				CurData: curData,
			},
		},
	}
	auditResult, err := log.clientSet.CoreService().Audit().SaveAuditLog(log.kit.Ctx, log.kit.Header, auditLog)
	if err != nil {
		blog.ErrorJSON("SaveAuditLog %s %s audit log failed, err: %s, result: %+v,rid:%s", auditAction, log.resourceType, err, auditResult, log.kit.Rid)
	}
	if auditResult.Result != true {
		blog.ErrorJSON("SaveAuditLog %s %s audit log failed, err: %s, result: %s,rid:%s", auditAction, log.resourceType, err, auditResult, log.kit.Rid)
	}
	return
}

func (log *ObjectAudit) buildSnapshotForPre() ObjAuditLog {
	query := mapstr.MapStr{"id": log.id}
	rsp, err := log.clientSet.CoreService().Model().ReadModel(log.kit.Ctx, log.kit.Header, &metadata.QueryCondition{Condition: query})
	if err != nil {
		blog.Errorf("[audit] failed to build the objData, error info is %s, rid: %s", err.Error(), log.kit.Rid)
		return log
	}
	if rsp.Result != true {
		blog.Errorf("[audit] failed to build the objData,rsp code is %v, err: %s", rsp.Code, rsp.ErrMsg)
		return log
	}
	if len(rsp.Data.Info) <= 0 {
		blog.Errorf("[audit] failed to build the objData,err: %s", log.kit.CCError.CCError(common.CCErrorModelNotFound))
		return log
	}
	log.preData = rsp.Data.Info[0].Spec
	log.bkObjID = log.preData.ObjectID
	log.bkObjName = log.preData.ObjectName
	return log
}

func (log *ObjectAudit) buildSnapshotForCur() ObjAuditLog {
	query := mapstr.MapStr{"id": log.id}
	rsp, err := log.clientSet.CoreService().Model().ReadModel(log.kit.Ctx, log.kit.Header, &metadata.QueryCondition{Condition: query})
	if err != nil {
		blog.Errorf("[audit] failed to build the objData, error info is %s, rid: %s", err.Error(), log.kit.Rid)
		return log
	}
	if rsp.Result != true {
		blog.Errorf("[audit] failed to build the objData,rsp code is %v, err: %s", rsp.Code, rsp.ErrMsg)
		return log
	}
	if len(rsp.Data.Info) <= 0 {
		blog.Errorf("[audit] failed to build the objData,err: %s", log.kit.CCError.CCError(common.CCErrorModelNotFound))
		return log
	}
	log.curData = rsp.Data.Info[0].Spec
	log.bkObjID = log.curData.ObjectID
	log.bkObjName = log.curData.ObjectName
	return log
}
