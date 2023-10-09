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
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type hostModuleLog struct {
	audit     audit
	hostIDArr []int64
	pre       []metadata.ModuleHost
	cur       []metadata.ModuleHost
}

// NewHostModuleLog TODO
func NewHostModuleLog(clientSet coreservice.CoreServiceClientInterface, hostID []int64) *hostModuleLog {
	return &hostModuleLog{
		audit: audit{
			clientSet: clientSet,
		},
		hostIDArr: hostID,
	}
}

// WithPrevious TODO
func (h *hostModuleLog) WithPrevious(kit *rest.Kit) errors.CCError {
	if h.pre != nil {
		return nil
	}
	var err error
	h.pre, err = h.getHostModuleConfig(kit)
	if err != nil {
		return err
	}
	return nil
}

// WithCurrent TODO
func (h *hostModuleLog) WithCurrent(kit *rest.Kit) errors.CCError {
	if h.cur != nil {
		return nil
	}
	var err error
	h.cur, err = h.getHostModuleConfig(kit)
	if err != nil {
		return err
	}
	return nil
}

// SaveAudit save the audit msg
func (h *hostModuleLog) SaveAudit(kit *rest.Kit) errors.CCError {
	hostInfos, err := h.getInnerIPAndInnerIPv6(kit)
	if err != nil {
		return err
	}

	if err := h.WithCurrent(kit); err != nil {
		return err
	}

	defaultBizID, err := h.audit.getDefaultAppID(kit)
	if err != nil {
		blog.ErrorJSON("save audit failed, failed to get default appID, err: %s, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	var setIDs, moduleIDs, appIDs []int64
	for _, val := range append(h.pre, h.cur...) {
		setIDs = append(setIDs, val.SetID)
		moduleIDs = append(moduleIDs, val.ModuleID)
		appIDs = append(appIDs, val.AppID)
	}

	moduleMap, err := h.getInstIDNameMap(kit, common.BKInnerObjIDModule, moduleIDs)
	if err != nil {
		return err
	}

	setMap, err := h.getInstIDNameMap(kit, common.BKInnerObjIDSet, setIDs)
	if err != nil {
		return err
	}

	bizMap, err := h.getInstIDNameMap(kit, common.BKInnerObjIDApp, appIDs)
	if err != nil {
		return err
	}

	preDataMap := h.getHostTransferDataMap(h.pre, bizMap, setMap, moduleMap)
	curDataMap := h.getHostTransferDataMap(h.cur, bizMap, setMap, moduleMap)

	var logs = make([]metadata.AuditLog, 0)
	for _, host := range hostInfos {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			return err
		}

		hostIP, err := host.String(common.BKHostInnerIPField)
		if err != nil {
			return err
		}

		hostIPv6, err := host.String(common.BKHostInnerIPv6Field)
		if err != nil {
			return err
		}

		preData, curData := preDataMap[hostID], curDataMap[hostID]
		preBizID, curBizID := preData.BizID, curData.BizID

		var action metadata.ActionType
		var bizID int64
		if preBizID != curBizID && preBizID == defaultBizID {
			action = metadata.AuditAssignHost
			bizID = curBizID
		} else if preBizID != curBizID && curBizID == defaultBizID {
			action = metadata.AuditUnassignHost
			bizID = preBizID
		} else {
			action = metadata.AuditTransferHostModule
			bizID = curBizID
		}

		// generate audit log.
		logs = append(logs, metadata.AuditLog{
			AuditType:          metadata.HostType,
			ResourceType:       metadata.HostRes,
			Action:             action,
			BusinessID:         bizID,
			ResourceID:         hostID,
			ResourceName:       hostIP,
			ExtendResourceName: hostIPv6,
			OperationDetail: &metadata.HostTransferOpDetail{
				PreData: preData,
				CurData: curData,
			},
		})
	}

	// save audit log.
	if err := h.audit.SaveAuditLog(kit, logs...); err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func (h *hostModuleLog) getHostModuleConfig(kit *rest.Kit) ([]metadata.ModuleHost, errors.CCError) {
	conds := &metadata.HostModuleRelationRequest{
		HostIDArr: h.hostIDArr,
		Fields:    []string{common.BKAppIDField, common.BKSetIDField, common.BKModuleIDField, common.BKHostIDField},
	}
	result, err := h.audit.clientSet.Host().GetHostModuleRelation(kit.Ctx, kit.Header, conds)
	if err != nil {
		blog.Errorf("get host module config failed, http do error, err: %s, input: %+v, rid: %s", err.Error(),
			conds, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	return result.Info, nil
}

func (h *hostModuleLog) getInnerIPAndInnerIPv6(kit *rest.Kit) ([]mapstr.MapStr, errors.CCError) {
	query := &metadata.QueryInput{
		Start:     0,
		Limit:     len(h.hostIDArr),
		Sort:      common.BKAppIDField,
		Condition: common.KvMap{common.BKHostIDField: common.KvMap{common.BKDBIN: h.hostIDArr}},
		Fields: fmt.Sprintf("%s,%s,%s", common.BKHostIDField, common.BKHostInnerIPField,
			common.BKHostInnerIPv6Field),
	}

	result, err := h.audit.clientSet.Host().GetHosts(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("get hosts failed, http do error, err: %v, input: %+v, rid: %s", err, query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	return result.Info, nil
}

func (h *hostModuleLog) getInstIDNameMap(kit *rest.Kit, objID string, ids []int64) (map[int64]string, error) {
	if ids == nil {
		return make(map[int64]string), nil
	}

	idField := metadata.GetInstIDFieldByObjID(objID)
	nameField := metadata.GetInstNameFieldName(objID)

	query := &metadata.QueryCondition{
		Page:      metadata.BasePage{Start: 0, Limit: common.BKNoLimit},
		Condition: mapstr.MapStr{idField: common.KvMap{common.BKDBIN: ids}},
		Fields:    []string{idField, nameField},
	}

	result, err := h.audit.clientSet.Instance().ReadInstance(kit.Ctx, kit.Header, objID, query)
	if err != nil {
		blog.Errorf("get %s id to name map failed, err: %v, input: %+v, rid: %s", objID, err, query, kit.Rid)
		return nil, err
	}

	idNameMap := make(map[int64]string)
	for _, info := range result.Info {
		id, err := util.GetInt64ByInterface(info[idField])
		if err != nil {
			blog.Errorf("parse %s id %+v to int failed, err: %v, rid: %s", objID, info[idField], err, kit.Rid)
			return nil, err
		}

		idNameMap[id] = util.GetStrByInterface(info[nameField])
	}

	return idNameMap, nil
}

func (h *hostModuleLog) getHostTransferDataMap(relations []metadata.ModuleHost, bizMap, setMap,
	moduleMap map[int64]string) map[int64]metadata.HostBizTopo {

	hostRelationMap := make(map[int64]map[int64][]metadata.Module)
	hostAppMap := make(map[int64]int64)

	for _, val := range relations {
		if _, ok := hostRelationMap[val.HostID]; !ok {
			hostRelationMap[val.HostID] = make(map[int64][]metadata.Module)
		}

		hostAppMap[val.HostID] = val.AppID
		hostRelationMap[val.HostID][val.SetID] = append(hostRelationMap[val.HostID][val.SetID],
			metadata.Module{ModuleID: val.ModuleID, ModuleName: moduleMap[val.ModuleID]})
	}

	hostDataMap := make(map[int64]metadata.HostBizTopo)
	for hostID, relationMap := range hostRelationMap {
		sets := make([]metadata.Topo, 0)
		for setID, modules := range relationMap {
			sets = append(sets, metadata.Topo{
				SetID:   setID,
				SetName: setMap[setID],
				Module:  modules,
			})
		}

		bizID := hostAppMap[hostID]
		hostDataMap[hostID] = metadata.HostBizTopo{
			BizID:   bizID,
			BizName: bizMap[bizID],
			Set:     sets,
		}
	}

	return hostDataMap
}
