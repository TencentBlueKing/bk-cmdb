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

type ObjectAsstAudit struct {
	kit           *rest.Kit
	clientSet     apimachinery.ClientSetInterface
	auditType     metadata.AuditType
	resourceType  metadata.ResourceType
	id            int64
	bkObjID       string
	bkObjName     string
	bkAsstObjID   string
	bkAsstObjName string
	preData       metadata.Association
	curData       metadata.Association
}

func NewObjectAsstAudit(kit *rest.Kit, clientSet apimachinery.ClientSetInterface, ID int64) *ObjectAsstAudit {
	return &ObjectAsstAudit{
		kit:          kit,
		clientSet:    clientSet,
		auditType:    metadata.ModelType,
		resourceType: metadata.ModelAssociationRes,
		id:           ID,
	}
}

func (log *ObjectAsstAudit) SaveAuditLog(auditAction metadata.ActionType) errors.CCError {
	var err errors.CCError
	preData := mapstr.MapStr{}
	curData := mapstr.MapStr{}
	switch auditAction {
	case metadata.AuditDelete:
		preData = log.preData.ToMapStr()
	case metadata.AuditCreate:
		curData = log.curData.ToMapStr()
	case metadata.AuditUpdate:
		preData = log.preData.ToMapStr()
		curData = log.curData.ToMapStr()
	}
	//get objectName
	log.bkObjName, err = log.getObjectInfo(log.kit, log.bkObjID)
	if err != nil {
		blog.Errorf("[audit] failed to get object info, err: %s", err)
	}
	//get objectName
	log.bkAsstObjName, err = log.getObjectInfo(log.kit, log.bkAsstObjID)
	if err != nil {
		blog.Errorf("[audit] failed to get object info, err: %s", err)
	}
	auditLog := metadata.AuditLog{
		AuditType:    log.auditType,
		ResourceType: log.resourceType,
		Action:       auditAction,
		OperationDetail: &metadata.ModelAsstOpDetail{
			BkObjID:       log.bkObjID,
			BkObjName:     log.bkObjName,
			BkAsstObjID:   log.bkAsstObjID,
			BkAsstObjName: log.bkAsstObjName,
			BasicOpDetail: metadata.BasicOpDetail{
				ResourceID:   log.id,
				ResourceName: log.bkObjID,
				Details: &metadata.BasicContent{
					PreData: preData,
					CurData: curData,
				},
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

func (log *ObjectAsstAudit) buildSnapshotForPre() ObjAuditLog {
	query := mapstr.MapStr{"id": log.id}
	rsp, err := log.clientSet.CoreService().Association().ReadModelAssociation(log.kit.Ctx, log.kit.Header, &metadata.QueryCondition{Condition: query})
	if err != nil {
		blog.Errorf("[audit] failed to get object association info, error info is %s, rid: %s", err.Error(), log.kit.Rid)
		return nil
	}
	if rsp.Result != true {
		blog.Errorf("[audit] failed to get object association info, rsp code is %v, err: %s", rsp.Code, rsp.ErrMsg)
		return nil
	}
	if len(rsp.Data.Info) <= 0 {
		blog.Errorf("[audit] failed to get object association info, err: %s", log.kit.CCError.CCError(common.CCErrorModelNotFound))
		return nil
	}
	log.preData = rsp.Data.Info[0]
	log.bkObjID = log.preData.ObjectID
	log.bkAsstObjID = log.preData.AsstObjID
	return log
}

func (log *ObjectAsstAudit) buildSnapshotForCur() ObjAuditLog {
	query := mapstr.MapStr{"id": log.id}
	rsp, err := log.clientSet.CoreService().Association().ReadModelAssociation(log.kit.Ctx, log.kit.Header, &metadata.QueryCondition{Condition: query})
	if err != nil {
		blog.Errorf("[audit] failed to get object association info, error info is %s, rid: %s", err.Error(), log.kit.Rid)
		return nil
	}
	if rsp.Result != true {
		blog.Errorf("[audit] failed to get object association info, rsp code is %v, err: %s", rsp.Code, rsp.ErrMsg)
		return nil
	}
	if len(rsp.Data.Info) <= 0 {
		blog.Errorf("[audit] failed to get object association info, err: %s", log.kit.CCError.CCError(common.CCErrorModelNotFound))
		return nil
	}
	log.curData = rsp.Data.Info[0]
	log.bkObjID = log.curData.ObjectID
	log.bkAsstObjID = log.curData.AsstObjID
	return log
}

func (log *ObjectAsstAudit) getObjectInfo(kit *rest.Kit, bkObjectID string) (string, errors.CCError) {
	query := make(map[string]interface{})
	if bkObjectID == "" {
		query = mapstr.MapStr{"bk_obj_id": log.bkObjID}
	} else {
		query = mapstr.MapStr{"bk_obj_id": bkObjectID}
	}
	//get objectName
	resp, err := log.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: query})
	if err != nil {
		return "", err
	}
	if resp.Result != true {
		return "", kit.CCError.New(resp.Code, resp.ErrMsg)
	}
	if len(resp.Data.Info) <= 0 {
		return "", kit.CCError.CCError(common.CCErrorModelNotFound)
	}
	return resp.Data.Info[0].Spec.ObjectName, nil
}

func (log *ObjectAsstAudit) transInfoToPre(data metadata.Association) {
	log.preData = data
	log.bkObjID = data.ObjectID
	log.bkAsstObjID = data.AsstObjID
}
