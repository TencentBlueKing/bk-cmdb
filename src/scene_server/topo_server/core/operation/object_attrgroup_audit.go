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

type ObjectAttrGroupAudit struct {
	kit          *rest.Kit
	clientSet    apimachinery.ClientSetInterface
	auditType    metadata.AuditType
	resourceType metadata.ResourceType
	id           int64
	bizID        int64
	bkGroupID    string
	bkGroupName  string
	bkObjectID   string
	bkObjectName string
	preData      metadata.Group
	curData      metadata.Group
}

func NewObjectAttrGroupAudit(kit *rest.Kit, clientSet apimachinery.ClientSetInterface, ID int64) ObjAuditLog {
	return &ObjectAttrGroupAudit{
		kit:          kit,
		clientSet:    clientSet,
		auditType:    metadata.ModelType,
		resourceType: metadata.ModelGroupRes,
		id:           ID,
	}
}

func (log *ObjectAttrGroupAudit) SaveAuditLog(auditAction metadata.ActionType) errors.CCError {
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
	err := log.getObjectInfo(log.kit, log.bkObjectID)
	if err != nil {
		blog.Errorf("[audit] failed to get the objInfo,err: %s", err)
	}
	//make auditLog
	auditLog := metadata.AuditLog{
		AuditType:    log.auditType,
		ResourceType: log.resourceType,
		Action:       auditAction,
		OperationDetail: &metadata.ModelAttrOpDetail{
			BkObjID:   log.bkObjectID,
			BkObjName: log.bkObjectName,
			BasicOpDetail: metadata.BasicOpDetail{
				BusinessID:   log.bizID,
				ResourceID:   log.id,
				ResourceName: log.bkGroupName,
				Details: &metadata.BasicContent{
					PreData: preData,
					CurData: curData,
				},
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

func (log *ObjectAttrGroupAudit) buildSnapshotForPre() ObjAuditLog {
	query := mapstr.MapStr{"id": log.id}
	rsp, err := log.clientSet.CoreService().Model().ReadAttributeGroupByCondition(log.kit.Ctx, log.kit.Header, metadata.QueryCondition{Condition: query})
	if err != nil {
		blog.Errorf("[audit] failed to build the objAttrGroupData, error info is %s, rid: %s", err.Error(), log.kit.Rid)
		return nil
	}
	if rsp.Result != true {
		blog.Errorf("[audit] failed to build the objAttrGroupData,rsp code is %v, err: %s", rsp.Code, rsp.ErrMsg)
		return nil
	}
	if len(rsp.Data.Info) <= 0 {
		blog.Errorf("[audit] failed to build the objAttrGroupData,err: %s", log.kit.CCError.CCError(common.CCErrorModelNotFound))
		return nil
	}
	log.preData = rsp.Data.Info[0]
	log.bkObjectID = log.preData.ObjectID
	log.bkGroupID = log.preData.GroupID
	log.bkGroupName = log.preData.GroupName
	log.bizID, _ = log.preData.Metadata.ParseBizID()
	return log
}

func (log *ObjectAttrGroupAudit) buildSnapshotForCur() ObjAuditLog {
	query := mapstr.MapStr{"id": log.id}
	rsp, err := log.clientSet.CoreService().Model().ReadAttributeGroupByCondition(log.kit.Ctx, log.kit.Header, metadata.QueryCondition{Condition: query})
	if err != nil {
		blog.Errorf("[audit] failed to build the objAttrGroupData, error info is %s, rid: %s", err.Error(), log.kit.Rid)
		return nil
	}
	if rsp.Result != true {
		blog.Errorf("[audit] failed to build the objAttrGroupData,rsp code is %v, err: %s", rsp.Code, rsp.ErrMsg)
		return nil
	}
	if len(rsp.Data.Info) <= 0 {
		blog.Errorf("[audit] failed to build the objAttrGroupData,err: %s", log.kit.CCError.CCError(common.CCErrorModelNotFound))
		return nil
	}
	log.curData = rsp.Data.Info[0]
	log.bkObjectID = log.curData.ObjectID
	log.bkGroupID = log.curData.GroupID
	log.bkGroupName = log.curData.GroupName
	log.bizID, _ = log.preData.Metadata.ParseBizID()
	return log
}

//search DB to get bkObjName by bkObjID.
func (log *ObjectAttrGroupAudit) getObjectInfo(kit *rest.Kit, bkObjectID string) errors.CCError {
	query := make(map[string]interface{})
	if bkObjectID == "" {
		query = mapstr.MapStr{"bk_obj_id": log.bkObjectID}
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
	log.bkObjectName = resp.Data.Info[0].Spec.ObjectName
	return nil
}
