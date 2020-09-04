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
	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type attributeGroupAuditLog struct {
	audit
}

// GenerateAuditLog generate audit of model attribute, if data is nil, then auto get current model attribute data by id.
func (h *attributeGroupAuditLog) GenerateAuditLog(kit *rest.Kit, action metadata.ActionType, id int64, OperateFrom metadata.OperateFromType,
	data *metadata.Group, updateFields map[string]interface{}) (*metadata.AuditLog, error) {
	// get data by attribute instance id.
	if data == nil {
		query := mapstr.MapStr{"id": id}
		// get current model attribute data by id.
		rsp, err := h.clientSet.Model().ReadAttributeGroupByCondition(kit.Ctx, kit.Header, metadata.QueryCondition{Condition: query})
		if err != nil {
			blog.Errorf("generate audit log of attribute group failed, failed to read attribute group, err: %v, rid: %s",
				err.Error(), kit.Rid)
			return nil, err
		}
		if rsp.Result != true {
			blog.Errorf("generate audit log of attribute group failed, failed to read attribute group, rsp code is %v, err: %s",
				rsp.Code, rsp.ErrMsg)
			return nil, err
		}
		if len(rsp.Data.Info) <= 0 {
			blog.Errorf("generate audit log of model attribute failed, failed to read attribute group, err: %s",
				kit.CCError.CCError(common.CCErrorModelNotFound))
			return nil, err
		}

		data = &rsp.Data.Info[0]
	}

	bizID, _ := data.Metadata.ParseBizID()
	groupName := data.GroupName
	objID := data.ObjectID
	objName, err := h.getObjNameByObjID(kit, objID)
	if err != nil {
		return nil, err
	}

	var basicDetail *metadata.BasicContent
	switch action {
	case metadata.AuditCreate:
		basicDetail = &metadata.BasicContent{
			CurData: data.ToMapStr(),
		}
	case metadata.AuditDelete:
		basicDetail = &metadata.BasicContent{
			PreData: data.ToMapStr(),
		}
	case metadata.AuditUpdate:
		basicDetail = &metadata.BasicContent{
			PreData:      data.ToMapStr(),
			UpdateFields: updateFields,
		}
	}

	var auditLog = &metadata.AuditLog{
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModelGroupRes,
		Action:       action,
		BusinessID:   bizID,
		ResourceID:   id,
		ResourceName: groupName,
		OperateFrom:  OperateFrom,
		OperationDetail: &metadata.ModelAttrOpDetail{
			BkObjID:   objID,
			BkObjName: objName,
			BasicOpDetail: metadata.BasicOpDetail{
				Details: basicDetail,
			},
		},
	}

	return auditLog, nil
}

func NewAttributeGroupAuditLog(clientSet coreservice.CoreServiceClientInterface) *attributeGroupAuditLog {
	return &attributeGroupAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
