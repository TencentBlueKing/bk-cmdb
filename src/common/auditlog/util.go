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
	"context"
	"net/http"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// AuditInterface audit log methods
type AuditInterface interface {
	GetInstNameByID(ctx context.Context, objID string, instID int64) (string, error)
	GetAuditLogProperty(ctx context.Context, objID string) ([]metadata.Property, error)
}

type auditLog struct {
	clientSet apimachinery.ClientSetInterface
	header    http.Header
	rid       string
	ccErr     errors.DefaultCCErrorIf
}

func NewAudit(clientSet apimachinery.ClientSetInterface, header http.Header) AuditInterface {
	return &auditLog{
		clientSet: clientSet,
		header:    header,
		rid:       util.GetHTTPCCRequestID(header),
		ccErr:     util.GetDefaultCCError(header),
	}
}

func (a *auditLog) GetInstNameByID(ctx context.Context, objID string, instID int64) (string, error) {
	switch objID {
	case common.BKInnerObjIDHost:
		rsp, err := a.clientSet.CoreService().Host().GetHosts(ctx, a.header, &metadata.QueryInput{
			Fields: common.BKHostInnerIPField,
			Condition: map[string]interface{}{
				common.BKHostIDField: instID,
			},
		})
		if nil != err || !rsp.Result {
			blog.ErrorfDepth(1, "GetInstNameByID %d GetHosts failed, err: %v, rsp: %+v, rid: %s", instID, err, rsp, a.rid)
			return "", a.ccErr.CCError(common.CCErrAuditTakeSnapshotFailed)
		}
		if rsp.Data.Count != 1 {
			blog.ErrorfDepth(1, "GetInstNameByID %d GetHosts find %d insts, rid: %s", instID, rsp.Data.Count, a.rid)
			return "", a.ccErr.CCError(common.CCErrAuditTakeSnapshotFailed)
		}
		ip, err := rsp.Data.Info[0].String(common.BKHostInnerIPField)
		if err != nil {
			blog.ErrorfDepth(1, "GetInstNameByID %d GetHosts parse ip failed, data: %+v, rid: %s", instID, rsp.Data.Info[0], a.rid)
			return "", a.ccErr.CCError(common.CCErrAuditTakeSnapshotFailed)
		}
		return ip, nil
	default:
		rsp, err := a.clientSet.CoreService().Instance().ReadInstance(ctx, a.header, objID, &metadata.QueryCondition{
			Fields: []string{common.GetInstNameField(objID)},
			Condition: map[string]interface{}{
				common.GetInstIDField(objID): instID,
			},
		})
		if nil != err || !rsp.Result {
			blog.ErrorfDepth(1, "GetInstNameByID %s %d ReadInstance failed, err: %v, rsp: %+v, rid: %s", objID, instID, err, rsp, a.rid)
			return "", a.ccErr.CCError(common.CCErrAuditTakeSnapshotFailed)
		}
		if rsp.Data.Count != 1 {
			blog.ErrorfDepth(1, "GetInstNameByID %d ReadInstance find %d insts, rid: %s", instID, rsp.Data.Count, a.rid)
			return "", a.ccErr.CCError(common.CCErrAuditTakeSnapshotFailed)
		}
		instName, err := rsp.Data.Info[0].String(common.GetInstNameField(objID))
		if err != nil {
			blog.ErrorfDepth(1, "GetInstNameByID %d ReadInstance parse inst name failed, data: %+v, rid: %s", instID, rsp.Data.Info[0], a.rid)
			return "", a.ccErr.CCError(common.CCErrAuditTakeSnapshotFailed)
		}
		return instName, nil
	}
}

func (a *auditLog) GetAuditLogProperty(ctx context.Context, objID string) ([]metadata.Property, error) {
	supplierAccount := util.GetOwnerID(a.header)

	cond := map[string]interface{}{
		metadata.AttributeFieldObjectID:        objID,
		metadata.AttributeFieldSupplierAccount: supplierAccount,
		metadata.AttributeFieldIsSystem: map[string]interface{}{
			common.BKDBNE: true,
		},
		metadata.AttributeFieldIsAPI: map[string]interface{}{
			common.BKDBNE: true,
		},
	}
	rsp, err := a.clientSet.CoreService().Model().ReadModelAttr(ctx, a.header, objID, &metadata.QueryCondition{Condition: cond})
	if nil != err || !rsp.Result {
		blog.ErrorfDepth(1, "GetAuditLogProperty failed to get the object(%s)' attribute, error: %v, rsp: %+v, rid: %s", objID, err, rsp, a.rid)
		return nil, a.ccErr.CCError(common.CCErrAuditTakeSnapshotFailed)
	}

	properties := make([]metadata.Property, 0)
	for _, attr := range rsp.Data.Info {
		properties = append(properties, metadata.Property{
			PropertyID:   attr.PropertyID,
			PropertyName: attr.PropertyName,
		})
	}
	return properties, nil
}
