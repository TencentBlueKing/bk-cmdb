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

type ObjectClsAudit struct {
	kit          *rest.Kit
	clientSet    apimachinery.ClientSetInterface
	auditType    metadata.AuditType
	resourceType metadata.ResourceType
	id           int64
	objClsID     string
	objClsName   string
	preData      metadata.Classification
	curData      metadata.Classification
}

func NewObjectClsAudit(kit *rest.Kit, clientSet apimachinery.ClientSetInterface, ID int64) ObjAuditLog {
	return &ObjectClsAudit{
		kit:          kit,
		clientSet:    clientSet,
		auditType:    metadata.ModelType,
		resourceType: metadata.ModelClassificationRes,
		id:           ID,
	}
}

func (log *ObjectClsAudit) SaveAuditLog(auditAction metadata.ActionType) errors.CCError {
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
		ResourceID:   log.id,
		ResourceName: log.objClsName,
		OperationDetail: &metadata.BasicOpDetail{
			Details: &metadata.BasicContent{
				PreData: preData,
				CurData: curData,
			},
		},
	}
	auditResult, err := log.clientSet.CoreService().Audit().SaveAuditLog(log.kit.Ctx, log.kit.Header, auditLog)
	if err != nil {
		blog.ErrorJSON("SaveAuditLog %s %s audit log failed, err: %s, result: %+v,rid:%s", auditAction, log.resourceType, err, auditResult, log.kit.Rid)
		return log.kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}
	if auditResult.Result != true {
		blog.ErrorJSON("SaveAuditLog %s %s audit log failed, err: %s, result: %s,rid:%s", auditAction, log.resourceType, err, auditResult, log.kit.Rid)
		return log.kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}
	return nil
}

func (log *ObjectClsAudit) buildSnapshotForPre() ObjAuditLog {
	query := mapstr.MapStr{"id": log.id}
	rsp, err := log.clientSet.CoreService().Model().ReadModelClassification(log.kit.Ctx, log.kit.Header, &metadata.QueryCondition{Condition: query})
	if err != nil {
		blog.Errorf("[audit] failed to build the objClsData, error info is %s, rid: %s", err.Error(), log.kit.Rid)
		return nil
	}
	if rsp.Result != true {
		blog.Errorf("[audit] failed to build the objClsData,rsp code is %v, err: %s", rsp.Code, rsp.ErrMsg)
		return nil
	}
	if len(rsp.Data.Info) <= 0 {
		blog.Errorf("[audit] failed to build the objClsData,err: %s", log.kit.CCError.CCError(common.CCErrorModelNotFound))
		return nil
	}
	log.preData = rsp.Data.Info[0]
	log.objClsID = log.preData.ClassificationID
	log.objClsName = log.preData.ClassificationName
	return log
}

func (log *ObjectClsAudit) buildSnapshotForCur() ObjAuditLog {
	query := mapstr.MapStr{"id": log.id}
	rsp, err := log.clientSet.CoreService().Model().ReadModelClassification(log.kit.Ctx, log.kit.Header, &metadata.QueryCondition{Condition: query})
	if err != nil {
		blog.Errorf("[audit] failed to build the objClsData, error info is %s, rid: %s", err.Error(), log.kit.Rid)
		return nil
	}
	if rsp.Result != true {
		blog.Errorf("[audit] failed to build the objClsData,rsp code is %v, err: %s", rsp.Code, rsp.ErrMsg)
		return nil
	}
	if len(rsp.Data.Info) <= 0 {
		blog.Errorf("[audit] failed to build the objClsData,err: %s", log.kit.CCError.CCError(common.CCErrorModelNotFound))
		return nil
	}
	log.curData = rsp.Data.Info[0]
	log.objClsID = log.curData.ClassificationID
	log.objClsName = log.curData.ClassificationName
	return log
}
