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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type attributeGroupAuditLog struct {
	audit
}

// GenerateTableAuditLog generate audit of table model attribute,
// if data is nil, will auto get current model attribute data by id.
func (h *objectAttributeAuditLog) GenerateTableAuditLog(parameter *generateAuditCommonParameter, objID string, bizID,
	id int64, data *metadata.Attribute) (*metadata.AuditLog, error) {
	kit := parameter.kit

	if data == nil {
		// get current model attribute data by id.
		query := mapstr.MapStr{metadata.AttributeFieldID: id}
		rsp, err := h.clientSet.Model().ReadModelAttrsWithTableByCondition(kit.Ctx, kit.Header, bizID,
			&metadata.QueryCondition{Condition: query})
		if err != nil {
			blog.Errorf("generate audit log of model attribute failed, failed to read model attribute, err: %v, "+
				"rid: %s", err.Error(), kit.Rid)
			return nil, err
		}

		if len(rsp.Info) <= 0 {
			blog.Errorf("generate audit log of model attribute failed, failed to read model attribute, query: %+v,"+
				" err: %s, rid: %s", query, kit.CCError.CCError(common.CCErrorModelNotFound), kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrorModelNotFound)
		}

		data = &rsp.Info[0]
	}

	objName, err := h.getObjNameByObjID(kit, objID)
	if err != nil {
		return nil, err
	}

	return &metadata.AuditLog{
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModelAttributeRes,
		Action:       parameter.action,
		BusinessID:   data.BizID,
		ResourceID:   id,
		ResourceName: data.PropertyName,
		OperateFrom:  parameter.operateFrom,
		OperationDetail: &metadata.ModelAttrOpDetail{
			BkObjID:   data.ObjectID,
			BkObjName: objName,
			BasicOpDetail: metadata.BasicOpDetail{
				Details: parameter.NewBasicContent(data.ToMapStr()),
			},
		},
	}, nil
}

// GenerateAuditLog generate audit of object attribute group,
// if data is nil, will auto get current attribute group data by id.
func (h *attributeGroupAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter, id int64,
	data *metadata.Group) (*metadata.AuditLog, error) {
	kit := parameter.kit

	if data == nil {
		// get current object attribute group data by id.
		query := mapstr.MapStr{metadata.GroupFieldID: id}
		rsp, err := h.clientSet.Model().ReadAttributeGroupByCondition(kit.Ctx, kit.Header,
			metadata.QueryCondition{Condition: query})
		if err != nil {
			blog.Errorf("generate audit log of attribute group failed, failed to read attribute group, err: %v, "+
				"rid: %s", err.Error(), kit.Rid)
			return nil, err
		}

		if len(rsp.Info) <= 0 {
			blog.Errorf("generate audit log of model attribute failed, failed to read attribute group, err: %s, "+
				"rid: %s", kit.CCError.CCError(common.CCErrorModelNotFound), kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrorModelNotFound)
		}

		data = &rsp.Info[0]
	}

	objName, err := h.getObjNameByObjID(kit, data.ObjectID)
	if err != nil {
		return nil, err
	}

	return &metadata.AuditLog{
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModelAttributeGroupRes,
		Action:       parameter.action,
		BusinessID:   data.BizID,
		ResourceID:   id,
		ResourceName: data.GroupName,
		OperateFrom:  parameter.operateFrom,
		OperationDetail: &metadata.ModelAttrOpDetail{
			BkObjID:   data.ObjectID,
			BkObjName: objName,
			BasicOpDetail: metadata.BasicOpDetail{
				Details: parameter.NewBasicContent(data.ToMapStr()),
			},
		},
	}, nil
}

// NewAttributeGroupAuditLog TODO
func NewAttributeGroupAuditLog(clientSet coreservice.CoreServiceClientInterface) *attributeGroupAuditLog {
	return &attributeGroupAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
