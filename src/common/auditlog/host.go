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
	"configcenter/src/common/util"
)

type hostAuditLog struct {
	audit
}

// GenerateAuditLog generate audit log of hosts
func (h *hostAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter, bizID int64, data []mapstr.MapStr) (
	[]metadata.AuditLog, error) {

	return h.generateAuditLog(parameter, bizID, data)
}

<<<<<<< HEAD
// GenerateAuditLogByHostIDGetBizID generate audit log of host, auto get bizID by hostID and host topology relate,
// if data is nil, will auto get host current data by hostID, meanwhile update hostIP if hostIP is "".
func (h *hostAuditLog) GenerateAuditLogByHostIDGetBizID(parameter *generateAuditCommonParameter, hostID int64,
	innerIP string, data map[string]interface{}) (*metadata.AuditLog, error) {
	// get bizID by hostID and host topology related.
	bizID, err := h.getBizIDByHostID(parameter.kit, hostID)
	if err != nil {
		blog.Errorf("generate host audit log failed, failed to get bizID by hostID, hostID: %d, innerIP: %s, "+
			"err: %v, rid: %s", hostID, innerIP, err, parameter.kit.Rid)
=======
// GenerateAuditLogByCond generate audit log of hosts, auto get current hosts by condition.
func (h *hostAuditLog) GenerateAuditLogByCond(parameter *generateAuditCommonParameter, bizID int64,
	condition map[string]interface{}) ([]metadata.AuditLog, error) {
	data, err := h.getInstByCond(parameter.kit, common.BKInnerObjIDHost, condition, nil)
	if err != nil {
		blog.Errorf("get hosts failed, err: %v, condition: %#v, rid: %s", err, condition, parameter.kit.Rid)
>>>>>>> v3.9.x
		return nil, err
	}
	return h.generateAuditLog(parameter, bizID, data)
}

func (h *hostAuditLog) generateAuditLog(parameter *generateAuditCommonParameter, bizID int64,
	data []mapstr.MapStr) ([]metadata.AuditLog, error) {

	kit := parameter.kit
	if len(data) == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "host audit log data")
	}

	auditLogs := make([]metadata.AuditLog, len(data))
	hostIDs := make([]int64, len(data))
	for index, host := range data {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
<<<<<<< HEAD
			blog.Errorf("generate host audit log failed, failed to get host instance by bizID, hostID: %d, "+
				"err: %v, rid: %s", hostID, err, kit.Rid)
			return nil, err
=======
			blog.Errorf("parse host id failed, err: %v, host: %#v, rid: %s", err, host, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
>>>>>>> v3.9.x
		}
		hostIDs[index] = hostID

		auditLog := metadata.AuditLog{
			AuditType:    metadata.HostType,
			ResourceType: metadata.HostRes,
			Action:       parameter.action,
			BusinessID:   bizID,
			ResourceID:   hostIDs[index],
			ResourceName: util.GetStrByInterface(host[common.BKHostInnerIPField]),
			OperateFrom:  parameter.operateFrom,
			OperationDetail: &metadata.InstanceOpDetail{
				BasicOpDetail: metadata.BasicOpDetail{
					Details: parameter.NewBasicContent(host),
				},
				ModelID: common.BKInnerObjIDHost,
			},
		}
		auditLogs[index] = auditLog
	}

<<<<<<< HEAD
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

	var bizID int64
	if len(moduleHost.Info) > 0 {
		bizID = moduleHost.Info[0].AppID
=======
	if bizID == 0 {
		hostBizMap, err := h.getBizIDByHostID(parameter.kit, hostIDs)
		if err != nil {
			blog.Errorf("get biz id for hosts failed, err: %v, host ids: %+v, rid: %s", err, hostIDs, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
		}

		for index, auditLog := range auditLogs {
			auditLogs[index].BusinessID = hostBizMap[auditLog.ID]
		}
>>>>>>> v3.9.x
	}

	return auditLogs, nil
}

<<<<<<< HEAD
// getHostInstanceDetailByHostID get host data and hostIP by hostID.
func (h *hostAuditLog) getHostInstanceDetailByHostID(kit *rest.Kit, hostID int64) (map[string]interface{}, string,
	error) {
	// get host details.
	result, err := h.clientSet.Host().GetHostByID(kit.Ctx, kit.Header, hostID)
	if err != nil {
		blog.Errorf("get host instance detail failed, err: %v, hostID: %d, rid: %s", err, hostID, kit.Rid)
		return nil, "", kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get host instance detail failed, http response error, err code: %d, err msg: %s, "+
			"hostID: %d, rid: %s", result.Code, result.ErrMsg, hostID, kit.Rid)
		return nil, "", kit.CCError.New(result.Code, result.ErrMsg)
=======
// getBizIDByHostID get mapping of host id to biz id
func (h *hostAuditLog) getBizIDByHostID(kit *rest.Kit, hostIDs []int64) (map[int64]int64, error) {
	input := &metadata.HostModuleRelationRequest{HostIDArr: hostIDs, Fields: []string{common.BKHostIDField,
		common.BKAppIDField}}
	moduleHost, err := h.clientSet.Host().GetHostModuleRelation(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("failed to get host module relation, hostIDs: %+v, err: %v, rid: %s", hostIDs, err, kit.Rid)
		return nil, err
>>>>>>> v3.9.x
	}
	if err := moduleHost.CCError(); err != nil {
		blog.Errorf("failed to get host module relation, hostIDs: %+v, err: %v, rid: %s", hostIDs, err, kit.Rid)
		return nil, err
	}

<<<<<<< HEAD
	ip, ok := hostInfo[common.BKHostInnerIPField].(string)
	if !ok {
		blog.Errorf("get host instance detail failed, convert bk_host_innerip to string error, hostID: %d, "+
			"hostInfo: %+v, rid: %d", hostID, hostInfo, kit.Rid)
		return nil, "", kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost,
			common.BKHostInnerIPField, "string", "not string")

=======
	hostBizMap := make(map[int64]int64)
	for _, relation := range moduleHost.Data.Info {
		hostBizMap[relation.HostID] = relation.AppID
>>>>>>> v3.9.x
	}

	return hostBizMap, nil
}

func NewHostAudit(clientSet coreservice.CoreServiceClientInterface) *hostAuditLog {
	return &hostAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
