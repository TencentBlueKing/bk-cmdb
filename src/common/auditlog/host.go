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
	"configcenter/src/common/metadata"
)

type hostAuditLog struct {
	audit
}

// GenerateAuditLog generate audit log of host, if data is nil, will auto get host current data by hostID,
// meanwhile update hostIP if hostIP is "".
func (h *hostAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter, hostID, bizID int64, innerIP string,
	data map[string]interface{}) (*metadata.AuditLog, error) {
	return h.generateAuditLog(parameter, hostID, bizID, innerIP, data)
}

// GenerateAuditLogByHostIDGetBizID generate audit log of host, auto get bizID by hostID and host topology relate,
// if data is nil, will auto get host current data by hostID, meanwhile update hostIP if hostIP is "".
func (h *hostAuditLog) GenerateAuditLogByHostIDGetBizID(parameter *generateAuditCommonParameter, hostID int64, innerIP string,
	data map[string]interface{}) (*metadata.AuditLog, error) {
	// get bizID by hostID and host topology related.
	bizID, err := h.getBizIDByHostID(parameter.kit, hostID)
	if err != nil {
		blog.Errorf("generate host audit log failed, failed to get bizID by hostID, hostID: %d, innerIP: %s, err: %v, rid: %s",
			hostID, innerIP, err, parameter.kit.Rid)
		return nil, err
	}

	return h.generateAuditLog(parameter, hostID, bizID, innerIP, data)
}

func (h *hostAuditLog) generateAuditLog(parameter *generateAuditCommonParameter, hostID, bizID int64, hostIP string,
	data map[string]interface{}) (*metadata.AuditLog, error) {
	kit := parameter.kit

	// get current data and inner ip by hostID
	if data == nil {
		var innerIP string
		var err error
		data, innerIP, err = h.getHostInstanceDetailByHostID(kit, hostID)
		if err != nil {
			blog.Errorf("generate host audit log failed, failed to get host instance by bizID, hostID: %d, err: %v, rid: %s",
				hostID, err, kit.Rid)
			return nil, err
		}

		if innerIP != "" {
			hostIP = innerIP
		}
	}

	return &metadata.AuditLog{
		AuditType:    metadata.HostType,
		ResourceType: metadata.HostRes,
		Action:       parameter.action,
		BusinessID:   bizID,
		ResourceID:   hostID,
		ResourceName: hostIP,
		OperateFrom:  parameter.operateFrom,
		OperationDetail: &metadata.InstanceOpDetail{
			BasicOpDetail: metadata.BasicOpDetail{
				Details: parameter.NewBasicContent(data),
			},
			ModelID: common.BKInnerObjIDHost,
		},
	}, nil
}

// getBizIDByHostID get the bizID for host belong business.
func (h *hostAuditLog) getBizIDByHostID(kit *rest.Kit, hostID int64) (int64, error) {
	input := &metadata.HostModuleRelationRequest{HostIDArr: []int64{hostID}, Fields: []string{common.BKAppIDField}}
	moduleHost, err := h.clientSet.Host().GetHostModuleRelation(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("failed to get host module relation, hostID: %d, err: %v, rid: %s", hostID, err, kit.Rid)
		return 0, err
	}
	if !moduleHost.Result {
		blog.Errorf("failed to get host module relation, http response error, hostID: %d, errCode: %d, errMsg: %s, rid: %s", hostID,
			moduleHost.Code, moduleHost.ErrMsg, kit.Rid)
		return 0, kit.CCError.New(moduleHost.Code, moduleHost.ErrMsg)
	}

	var bizID int64
	if len(moduleHost.Data.Info) > 0 {
		bizID = moduleHost.Data.Info[0].AppID
	}

	return bizID, nil
}

// getHostInstanceDetailByHostID get host data and hostIP by hostID.
func (h *hostAuditLog) getHostInstanceDetailByHostID(kit *rest.Kit, hostID int64) (map[string]interface{}, string, error) {
	// get host details.
	result, err := h.clientSet.Host().GetHostByID(kit.Ctx, kit.Header, hostID)
	if err != nil {
		blog.Errorf("get host instance detail failed, err: %v, hostID: %d, rid: %s", err, hostID, kit.Rid)
		return nil, "", kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get host instance detail failed, http response error, err code: %d, err msg: %s, hostID: %d, rid: %s",
			result.Code, result.ErrMsg, hostID, kit.Rid)
		return nil, "", kit.CCError.New(result.Code, result.ErrMsg)
	}

	hostInfo := result.Data
	if len(hostInfo) == 0 {
		return nil, "", nil
	}

	ip, ok := hostInfo[common.BKHostInnerIPField].(string)
	if !ok {
		blog.Errorf("get host instance detail failed, convert bk_host_innerip to string error, hostID: %d, hostInfo: %+v, rid: %d",
			hostID, hostInfo, kit.Rid)
		return nil, "", kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostInnerIPField, "string", "not string")

	}

	return hostInfo, ip, nil
}

func NewHostAudit(clientSet coreservice.CoreServiceClientInterface) *hostAuditLog {
	return &hostAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
