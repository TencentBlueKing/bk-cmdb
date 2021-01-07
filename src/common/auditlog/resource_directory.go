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
	"fmt"

	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type resourceDirAuditLog struct {
	audit
}

// GenerateAuditLog generate audit log of resource directory, if data is nil, will auto get data by instModuleID.
func (h *resourceDirAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter, instModuleID, bizID int64, data mapstr.MapStr) (*metadata.AuditLog, error) {
	if data == nil {
		query := &metadata.QueryCondition{Condition: mapstr.MapStr{common.BKModuleIDField: instModuleID}}
		rsp, err := h.clientSet.Instance().ReadInstance(parameter.kit.Ctx, parameter.kit.Header, common.BKInnerObjIDModule, query)
		if err != nil {
			blog.Errorf("generate audit log of resource directory failed, failed to read resource directory, err: %v, rid: %s",
				err.Error(), parameter.kit.Rid)
			return nil, err
		}
		if rsp.Result != true {
			blog.Errorf("generate audit log of resource directory failed, failed to read resource directory, rsp code is %v, err: %s, rid: %s",
				rsp.Code, rsp.ErrMsg, parameter.kit.Rid)
			return nil, parameter.kit.CCError.New(rsp.Code, rsp.ErrMsg)
		}
		if len(rsp.Data.Info) <= 0 {
			blog.Errorf("generate audit log of resource directory failed, not find resource directory, instModuleID: %d, rid: %s",
				instModuleID, parameter.kit.Rid)
			return nil, fmt.Errorf("generate audit log of resource directory failed, not find resource directory")
		}

		data = rsp.Data.Info[0]
	}

	// get resource directory name.
	moduleName, err := data.String(common.BKModuleNameField)
	if err != nil {
		return nil, parameter.kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDModule,
			common.BKModuleNameField, "string", err.Error())
	}

	return &metadata.AuditLog{
		AuditType:    metadata.ModelInstanceType,
		ResourceType: metadata.ResourceDirRes,
		Action:       parameter.action,
		BusinessID:   bizID,
		ResourceID:   instModuleID,
		ResourceName: moduleName,
		OperateFrom:  parameter.operateFrom,
		OperationDetail: &metadata.InstanceOpDetail{
			BasicOpDetail: metadata.BasicOpDetail{
				Details: parameter.NewBasicContent(data),
			},
			ModelID: common.BKInnerObjIDModule,
		},
	}, nil
}

func NewResourceDirAuditLog(clientSet coreservice.CoreServiceClientInterface) *resourceDirAuditLog {
	return &resourceDirAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
