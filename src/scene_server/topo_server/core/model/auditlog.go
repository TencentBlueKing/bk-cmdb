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

package model

import (
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type ModAuditLog interface {
	WithPrevious(*rest.Kit, int64) errors.CCError
	WithCurrent(*rest.Kit, int64) errors.CCError
	SaveAuditLog(*rest.Kit, metadata.ActionType) errors.CCError
}

func (log *ObjectAuditLog) WithPrevious(kit *rest.Kit, id int64) errors.CCError {
	preData, err := log.buildLogData(kit, id)
	if err != nil {
		return err
	}
	log.preData = preData
	return nil
}

func (log *ObjectAuditLog) WithCurrent(kit *rest.Kit, id int64) errors.CCError {
	curData, err := log.buildLogData(kit, id)
	if err != nil {
		return err
	}
	log.curData = curData
	return nil
}

func (log *ObjectAuditLog) SaveAuditLog(kit *rest.Kit, auditAction metadata.ActionType) errors.CCError {
	auditLog := metadata.AuditLog{
		AuditType:    log.auditType,
		ResourceType: log.resourceType,
		Action:       auditAction,
		OperationDetail: &metadata.InstanceOpDetail{
			ModelID: log.objectID,
			BasicOpDetail: metadata.BasicOpDetail{
				ResourceID:   log.resourceID,
				ResourceName: log.resourceName,
				Details: &metadata.BasicContent{
					PreData:    log.preData,
					CurData:    log.curData,
					Properties: objectProperty,
				},
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

func (log *ObjectAuditLog) buildLogData(kit *rest.Kit, ID int64) (map[string]interface{}, errors.CCError) {
	query := mapstr.MapStr{"id": ID}
	switch log.auditType {
	case metadata.ModelType:
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
		log.resourceID = rsp.Data.Info[0].Spec.ID
		log.resourceName = rsp.Data.Info[0].Spec.ObjectName
		rspData := rsp.Data.Info[0].Spec.ToMapStr()
		return rspData, nil
	case metadata.ModelClassificationType:
		rsp, err := log.clientSet.CoreService().Model().ReadModelClassification(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: query})
		if err != nil {
			return nil, err
		}
		if rsp.Result != true {
			return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
		}
		if len(rsp.Data.Info) <= 0 {
			return nil, kit.CCError.CCError(common.CCErrorModelNotFound)
		}
		log.objectID = rsp.Data.Info[0].ClassificationID
		log.resourceID = rsp.Data.Info[0].ID
		log.resourceName = rsp.Data.Info[0].ClassificationName
		rspData := rsp.Data.Info[0].ToMapStr()
		return rspData, nil
	case metadata.ModelAttributeType:
		rsp, err := log.clientSet.CoreService().Model().ReadModelAttrByCondition(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: query})
		if err != nil {
			return nil, err
		}
		if rsp.Result != true {
			return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
		}
		if len(rsp.Data.Info) <= 0 {
			return nil, kit.CCError.CCError(common.CCErrorModelNotFound)
		}
		log.objectID = rsp.Data.Info[0].ObjectID
		log.resourceID = rsp.Data.Info[0].ID
		log.resourceName = rsp.Data.Info[0].PropertyName
		rspData := rsp.Data.Info[0].ToMapStr()
		return rspData, nil
	case metadata.ModelGroupType:
		rsp, err := log.clientSet.CoreService().Model().ReadAttributeGroupByCondition(kit.Ctx, kit.Header, metadata.QueryCondition{Condition: query})
		if err != nil {
			return nil, err
		}
		if rsp.Result != true {
			return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
		}
		if len(rsp.Data.Info) <= 0 {
			return nil, kit.CCError.CCError(common.CCErrorModelNotFound)
		}
		log.objectID = rsp.Data.Info[0].ObjectID
		log.resourceID = rsp.Data.Info[0].ID
		log.resourceName = rsp.Data.Info[0].GroupName
		rspData := rsp.Data.Info[0].ToMapStr()
		return rspData, nil
	}
	return nil, kit.CCError.Errorf(common.CCErrAuditSelectFailed)
}

type ObjectAuditLog struct {
	auditType    metadata.AuditType
	resourceType metadata.ResourceType
	resourceID   int64
	resourceName string
	clientSet    apimachinery.ClientSetInterface
	preData      map[string]interface{}
	curData      map[string]interface{}
	objectID     string
}

var objectProperty = []metadata.Property{
	{PropertyID: "bk_obj_id", PropertyName: "模型id"},
	{PropertyID: "bk_obj_name", PropertyName: "模型名称"},
	{PropertyID: "bk_classification_id", PropertyName: "模型分组id"},
	{PropertyID: "bk_classification_name", PropertyName: "模型分组名"},
	{PropertyID: "bk_group_id", PropertyName: "模型字段分组id"},
	{PropertyID: "bk_group_name", PropertyName: "模型字段分组名"},
	//{PropertyID:"creator", PropertyName:"创建人"},
}

func NewObjectAuditLog(clientSet apimachinery.ClientSetInterface, auditType metadata.AuditType, resourceType metadata.ResourceType) *ObjectAuditLog {
	return &ObjectAuditLog{
		auditType:    auditType,
		resourceType: resourceType,
		resourceName: "",
		clientSet:    clientSet,
	}
}
