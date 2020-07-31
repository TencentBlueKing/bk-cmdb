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

type ObjectUniqueAudit struct {
	kit          *rest.Kit
	clientSet    apimachinery.ClientSetInterface
	auditType    metadata.AuditType
	resourceType metadata.ResourceType
	id           int64
	bkObjID      string
	bkObjName    string
	preData      metadata.ObjectUnique
	curData      metadata.ObjectUnique
}

func NewObjectUniqueAudit(kit *rest.Kit, clientSet apimachinery.ClientSetInterface, ID int64) ObjAuditLog {
	return &ObjectUniqueAudit{
		kit:          kit,
		clientSet:    clientSet,
		auditType:    metadata.ModelType,
		resourceType: metadata.ModelUniqueRes,
		id:           ID,
	}
}

func (log *ObjectUniqueAudit) SaveAuditLog(auditAction metadata.ActionType) errors.CCError {
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
	//get objectName
	err := log.getObjectInfo(log.kit, log.bkObjID)
	if err != nil {
		blog.Errorf("[audit] failed to get object name, err: %s", err)
	}
	//make auditLog
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
		blog.ErrorJSON("%s %s audit log failed, err: %s, result: %+v, rid: %s", auditAction, log.resourceType, err, auditResult, log.kit.Rid)
		return err
	}
	if auditResult.Result != true {
		blog.ErrorJSON("%s %s audit log failed, err: %s, result: %+v, rid: %s", auditAction, log.resourceType, err, auditResult, log.kit.Rid)
		return errors.New(common.CCErrAuditSaveLogFailed, auditResult.ErrMsg)
	}
	return nil
}

func (log *ObjectUniqueAudit) buildSnapshotForPre() ObjAuditLog {
	query := mapstr.MapStr{"id": log.id}
	//get repData
	rsp, err := log.clientSet.CoreService().Model().ReadModelAttrUnique(log.kit.Ctx, log.kit.Header, metadata.QueryCondition{Condition: query})
	if err != nil {
		blog.Errorf("[audit] failed to get unique info, error info is %s, rid: %s", err.Error(), log.kit.Rid)
		return log
	}
	if rsp.Result != true {
		blog.Errorf("[audit] failed to get unique info, rsp code is %v, err: %s", rsp.Code, rsp.ErrMsg)
		return log
	}
	if len(rsp.Data.Info) <= 0 {
		blog.Errorf("[audit] failed to get unique info, err: %s", log.kit.CCError.CCError(common.CCErrorModelNotFound))
		return log
	}
	log.preData = rsp.Data.Info[0]
	log.bkObjID = log.preData.ObjID
	return log
}

func (log *ObjectUniqueAudit) buildSnapshotForCur() ObjAuditLog {
	query := mapstr.MapStr{"id": log.id}
	//get repData
	rsp, err := log.clientSet.CoreService().Model().ReadModelAttrUnique(log.kit.Ctx, log.kit.Header, metadata.QueryCondition{Condition: query})
	if err != nil {
		blog.Errorf("[audit] failed to get unique info, error info is %s, rid: %s", err.Error(), log.kit.Rid)
		return log
	}
	if rsp.Result != true {
		blog.Errorf("[audit] failed to get unique info, rsp code is %v, err: %s", rsp.Code, rsp.ErrMsg)
		return log
	}
	if len(rsp.Data.Info) <= 0 {
		blog.Errorf("[audit] failed to get unique info, err: %s", log.kit.CCError.CCError(common.CCErrorModelNotFound))
		return log
	}
	log.curData = rsp.Data.Info[0]
	log.bkObjID = log.curData.ObjID
	return log
}

//search DB to get bkObjName by bkObjID.
func (log *ObjectUniqueAudit) getObjectInfo(kit *rest.Kit, bkObjectID string) errors.CCError {
	query := make(map[string]interface{})
	if bkObjectID == "" {
		query = mapstr.MapStr{"bk_obj_id": log.bkObjID}
	} else {
		query = mapstr.MapStr{"bk_obj_id": bkObjectID}
	}
	//get objectName
	resp, err := log.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: query})
	if err != nil {
		return err
	}
	if resp.Result != true {
		return kit.CCError.New(resp.Code, resp.ErrMsg)
	}
	if len(resp.Data.Info) <= 0 {
		return kit.CCError.CCError(common.CCErrorModelNotFound)
	}
	log.bkObjName = resp.Data.Info[0].Spec.ObjectName
	return nil
}
