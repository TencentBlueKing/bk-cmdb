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

type objectAuditLog struct {
	audit
}

// GenerateAuditLog generate audit of model, if data is nil, will auto get current model data by id.
func (h *objectAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter, id int64, data *metadata.Object) (
	*metadata.AuditLog, error) {
	kit := parameter.kit
	if data == nil {
		// get current model data by id.
		query := mapstr.MapStr{metadata.ModelFieldID: id}
		rsp, err := h.clientSet.Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: query})
		if err != nil {
			blog.Errorf("generate audit log of model failed, failed to read model, err: %v, rid: %s",
				err.Error(), kit.Rid)
			return nil, err
		}
		if rsp.Result != true {
			blog.Errorf("generate audit log of model failed, failed to read model, rsp code is %v, err: %s, rid: %s",
				rsp.Code, rsp.ErrMsg, kit.Rid)
			return nil, parameter.kit.CCError.New(rsp.Code, rsp.ErrMsg)
		}
		if len(rsp.Data.Info) <= 0 {
			blog.Errorf("generate audit log of model failed, failed to read model, err: %s, rid: %s",
				kit.CCError.CCError(common.CCErrorModelNotFound), kit.Rid)
			return nil, err
		}

		data = &rsp.Data.Info[0].Spec
	}

	return &metadata.AuditLog{
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModelRes,
		Action:       parameter.action,
		ResourceID:   id,
		ResourceName: data.ObjectName,
		OperateFrom:  parameter.operateFrom,
		OperationDetail: &metadata.BasicOpDetail{
			Details: parameter.NewBasicContent(data.ToMapStr()),
		},
	}, nil
}

func NewObjectAuditLog(clientSet coreservice.CoreServiceClientInterface) *objectAuditLog {
	return &objectAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
