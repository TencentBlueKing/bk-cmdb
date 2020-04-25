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

type ObjAuditLog interface {
	WithPrevious(*rest.Kit, int64) errors.CCError
	WithCurrent(*rest.Kit, int64) errors.CCError
	SaveAuditLog(*rest.Kit, metadata.ActionType) errors.CCError
}

type ObjectAudit struct {
	//auditType 定义这是模型相关审计（前端对应“其它”类）
	auditType metadata.AuditType
	//resourceType 模型审计中按资源种类有多个审计方式
	resourceType     metadata.ResourceType
	clientSet        apimachinery.ClientSetInterface
	preData          metadata.ModelData
	curData          metadata.ModelData
	objClsID         string
	objClsName       string
	objectID         string
	objectName       string
	objAttrID        string
	objAttrName      string
	objAttrGroupID   string
	objAttrGroupName string
}

func NewObjectAudit(clientSet apimachinery.ClientSetInterface, resourceType metadata.ResourceType) *ObjectAudit {
	return &ObjectAudit{
		auditType:    metadata.ModelType,
		resourceType: resourceType,
		clientSet:    clientSet,
	}
}

func (log *ObjectAudit) WithPrevious() *ObjectAudit {
	log.preData.ClassificationID = log.objClsID
	log.preData.ClassificationName = log.objClsName
	log.preData.ObjectID = log.objectID
	log.preData.ObjectName = log.objClsName
	log.preData.AttributeID = log.objAttrID
	log.preData.AttributeName = log.objAttrName
	log.preData.AttrGroupId = log.objAttrGroupID
	log.preData.AttrGroupName = log.objAttrGroupName
	return log
}

func (log *ObjectAudit) WithCurrent() *ObjectAudit {
	log.curData.ClassificationID = log.objClsID
	log.curData.ClassificationName = log.objClsName
	log.curData.ObjectID = log.objectID
	log.curData.ObjectName = log.objClsName
	log.curData.AttributeID = log.objAttrID
	log.curData.AttributeName = log.objAttrName
	log.curData.AttrGroupId = log.objAttrGroupID
	log.curData.AttrGroupName = log.objAttrGroupName
	return log
}

func (log *ObjectAudit) SaveAuditLog(kit *rest.Kit, auditAction metadata.ActionType) errors.CCError {
	auditLog := metadata.AuditLog{
		AuditType:    log.auditType,
		ResourceType: log.resourceType,
		Action:       auditAction,
		OperationDetail: &metadata.ModelOpDetail{
			ClassificationID:   log.objClsID,
			ClassificationName: log.objClsName,
			ObjectID:           log.objectID,
			ObjectName:         log.objectName,
			AttrGroupId:        log.objAttrGroupID,
			AttrGroupName:      log.objAttrGroupName,
			AttributeID:        log.objAttrID,
			AttributeName:      log.objAttrName,
			PreData:            log.preData,
			CurData:            log.curData,
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

func (log *ObjectAudit) buildObjData(kit *rest.Kit, ID int64) *ObjectAudit {
	query := mapstr.MapStr{"id": ID}
	rsp, err := log.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: query})
	if err != nil {
		blog.Errorf("[audit] failed to build the objData, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil
	}
	if rsp.Result != true {
		blog.Errorf("[audit] failed to build the objData,rsp code is %v, err: %s", rsp.Code, rsp.ErrMsg)
		return nil
	}
	if len(rsp.Data.Info) <= 0 {
		blog.Errorf("[audit] failed to build the objData,err: %s", kit.CCError.CCError(common.CCErrorModelNotFound))
		return nil
	}
	log.objectID = rsp.Data.Info[0].Spec.ObjectID
	log.objectName = rsp.Data.Info[0].Spec.ObjectName
	return log
}

func (log *ObjectAudit) buildObjClsData(kit *rest.Kit, ID int64) *ObjectAudit {
	query := mapstr.MapStr{"id": ID}
	rsp, err := log.clientSet.CoreService().Model().ReadModelClassification(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: query})
	if err != nil {
		blog.Errorf("[audit] failed to build the objClsData, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil
	}
	if rsp.Result != true {
		blog.Errorf("[audit] failed to build the objData,rsp code is %v, err: %s", rsp.Code, rsp.ErrMsg)
		return nil
	}
	if len(rsp.Data.Info) <= 0 {
		blog.Errorf("[audit] failed to build the objData,err: %s", kit.CCError.CCError(common.CCErrorModelNotFound))
		return nil
	}
	log.objClsID = rsp.Data.Info[0].ClassificationID
	log.objClsName = rsp.Data.Info[0].ClassificationName
	return log
}

func (log *ObjectAudit) buildObjAttrData(kit *rest.Kit, ID int64) *ObjectAudit {
	query := mapstr.MapStr{"id": ID}
	//get repData
	rsp, err := log.clientSet.CoreService().Model().ReadModelAttrByCondition(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: query})
	if err != nil {
		blog.Errorf("[audit] failed to build the objAttrData, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil
	}
	if rsp.Result != true {
		blog.Errorf("[audit] failed to build the objAttrData,rsp code is %v, err: %s", rsp.Code, rsp.ErrMsg)
		return nil
	}
	if len(rsp.Data.Info) <= 0 {
		blog.Errorf("[audit] failed to build the objAttrData,err: %s", kit.CCError.CCError(common.CCErrorModelNotFound))
		return nil
	}
	log.objectID = rsp.Data.Info[0].ObjectID
	log.objAttrID = rsp.Data.Info[0].PropertyID
	log.objAttrName = rsp.Data.Info[0].PropertyName
	log.objAttrGroupID = rsp.Data.Info[0].PropertyGroup
	log.objAttrGroupName = rsp.Data.Info[0].PropertyName
	err = log.getObjectInfo(kit, log.objectID)
	if err != nil {
		blog.Errorf("[audit] failed to get the objInfo,err: %s", err)
	}
	return log
}

func (log *ObjectAudit) getObjectInfo(kit *rest.Kit, objectID string) errors.CCError {
	query := make(map[string]interface{})
	if objectID == "" {
		query = mapstr.MapStr{"bk_obj_id": log.objectID}
	} else {
		query = mapstr.MapStr{"bk_obj_id": objectID}
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
	log.objectName = resp.Data.Info[0].Spec.ObjectName
	return nil
}

func (log *ObjectAudit) buildObjAttrGroupData(kit *rest.Kit, ID int64) *ObjectAudit {
	query := mapstr.MapStr{"id": ID}
	rsp, err := log.clientSet.CoreService().Model().ReadAttributeGroupByCondition(kit.Ctx, kit.Header, metadata.QueryCondition{Condition: query})
	if err != nil {
		blog.Errorf("[audit] failed to build the objAttrGroupData, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil
	}
	if rsp.Result != true {
		blog.Errorf("[audit] failed to build the objAttrGroupData,rsp code is %v, err: %s", rsp.Code, rsp.ErrMsg)
		return nil
	}
	if len(rsp.Data.Info) <= 0 {
		blog.Errorf("[audit] failed to build the objAttrGroupData,err: %s", kit.CCError.CCError(common.CCErrorModelNotFound))
		return nil
	}
	log.objectID = rsp.Data.Info[0].ObjectID
	log.objAttrGroupID = rsp.Data.Info[0].GroupID
	log.objAttrGroupName = rsp.Data.Info[0].GroupName
	err = log.getObjectInfo(kit, log.objectID)
	if err != nil {
		blog.Errorf("[audit] failed to get the objInfo,err: %s", err)
	}
	return log
}

//Todo 后续还应补充从excel导入的审计功能
