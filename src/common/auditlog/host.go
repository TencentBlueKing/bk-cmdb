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

// GenerateAuditLogByCond generate audit log of hosts, auto get current hosts by condition.
func (h *hostAuditLog) GenerateAuditLogByCond(parameter *generateAuditCommonParameter, bizID int64,
	condition map[string]interface{}) ([]metadata.AuditLog, error) {
	data, err := h.getInstByCond(parameter.kit, common.BKInnerObjIDHost, condition, nil)
	if err != nil {
		blog.Errorf("get hosts failed, err: %v, condition: %#v, rid: %s", err, condition, parameter.kit.Rid)
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
			blog.Errorf("parse host id failed, err: %v, host: %#v, rid: %s", err, host, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
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

	if bizID == 0 {
		hostBizMap, err := h.getBizIDByHostID(parameter.kit, hostIDs)
		if err != nil {
			blog.Errorf("get biz id for hosts failed, err: %v, host ids: %+v, rid: %s", err, hostIDs, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
		}

		for index := range auditLogs {
			auditLogs[index].BusinessID = hostBizMap[hostIDs[index]]
		}
	}

	return auditLogs, nil
}

// getBizIDByHostID get mapping of host id to biz id
func (h *hostAuditLog) getBizIDByHostID(kit *rest.Kit, hostIDs []int64) (map[int64]int64, error) {
	input := &metadata.HostModuleRelationRequest{HostIDArr: hostIDs, Fields: []string{common.BKHostIDField,
		common.BKAppIDField}}
	moduleHost, err := h.clientSet.Host().GetHostModuleRelation(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("failed to get host module relation, hostIDs: %+v, err: %v, rid: %s", hostIDs, err, kit.Rid)
		return nil, err
	}

	hostBizMap := make(map[int64]int64)
	for _, relation := range moduleHost.Info {
		hostBizMap[relation.HostID] = relation.AppID
	}

	return hostBizMap, nil
}

// NewHostAudit TODO
func NewHostAudit(clientSet coreservice.CoreServiceClientInterface) *hostAuditLog {
	return &hostAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
