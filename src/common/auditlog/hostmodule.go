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

func NewHostModuleLog(clientSet coreservice.CoreServiceClientInterface, hostID []int64) *hostModuleLog {
	return &hostModuleLog{
		audit: audit{
			clientSet: clientSet,
		},
		hostIDArr: hostID,
	}
}

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

func (h *hostModuleLog) SaveAudit(kit *rest.Kit) errors.CCError {
	hostInfos, err := h.getInnerIP(kit)
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
	for _, val := range h.pre {
		setIDs = append(setIDs, val.SetID)
		moduleIDs = append(moduleIDs, val.ModuleID)
		appIDs = append(appIDs, val.AppID)
	}
	for _, val := range h.cur {
		setIDs = append(setIDs, val.SetID)
		moduleIDs = append(moduleIDs, val.ModuleID)
		appIDs = append(appIDs, val.AppID)
	}

	modules, err := h.getModules(kit, moduleIDs)
	if err != nil {
		return err
	}

	moduleNameMap := make(map[int64]string)
	for _, module := range modules {
		moduleID, err := util.GetInt64ByInterface(module[common.BKModuleIDField])
		if err != nil {
			return err
		}
		moduleName, err := module.String(common.BKModuleNameField)
		if err != nil {
			return err
		}
		moduleNameMap[moduleID] = moduleName
	}

	sets, err := h.getSets(kit, setIDs)
	if err != nil {
		return err
	}

	setNameMap := make(map[int64]string)
	for _, setInfo := range sets {
		setID, err := util.GetInt64ByInterface(setInfo[common.BKSetIDField])
		if err != nil {
			return err
		}
		setNameMap[setID], err = setInfo.String(common.BKSetNameField)
		if err != nil {
			return err
		}
	}

	preHostRelationMap := make(map[int64]map[int64][]metadata.Module)
	preHostAppMap := make(map[int64]int64)
	for _, val := range h.pre {
		if _, ok := preHostRelationMap[val.HostID]; false == ok {
			preHostRelationMap[val.HostID] = make(map[int64][]metadata.Module)
		}
		preHostAppMap[val.HostID] = val.AppID
		preHostRelationMap[val.HostID][val.SetID] = append(preHostRelationMap[val.HostID][val.SetID], metadata.Module{ModuleID: val.ModuleID, ModuleName: moduleNameMap[val.ModuleID]})
	}

	curHostRelationMap := make(map[int64]map[int64][]metadata.Module)
	curHostAppMap := make(map[int64]int64)
	for _, val := range h.cur {
		if _, ok := curHostRelationMap[val.HostID]; false == ok {
			curHostRelationMap[val.HostID] = make(map[int64][]metadata.Module)
		}
		curHostAppMap[val.HostID] = val.AppID
		curHostRelationMap[val.HostID][val.SetID] = append(curHostRelationMap[val.HostID][val.SetID], metadata.Module{ModuleID: val.ModuleID, ModuleName: moduleNameMap[val.ModuleID]})
	}

	appInfoArr, err := h.getApps(kit, appIDs)
	if err != nil {
		return err
	}

	appIDNameMap := make(map[int64]string, 0)
	for _, appInfo := range appInfoArr {
		bizID, err := appInfo.Int64(common.BKAppIDField)
		if err != nil {
			blog.ErrorJSON("appInfo get biz id err:%s, appInfo: %s, rid:%s", err.Error(), appInfo, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
		}
		name, err := appInfo.String(common.BKAppNameField)
		if err != nil {
			blog.ErrorJSON("appInfo get biz name err:%s, appInfo: %s, rid:%s", err.Error(), appInfo, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
		}

		appIDNameMap[bizID] = name
	}

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

		sets := make([]metadata.Topo, 0)
		for setID, modules := range preHostRelationMap[hostID] {
			sets = append(sets, metadata.Topo{
				SetID:   setID,
				SetName: setNameMap[setID],
				Module:  modules,
			})
		}

		preBizID := preHostAppMap[hostID]
		preData := metadata.HostBizTopo{
			BizID:   preBizID,
			BizName: appIDNameMap[preBizID],
			Set:     sets,
		}

		sets = make([]metadata.Topo, 0)
		for setID, modules := range curHostRelationMap[hostID] {
			sets = append(sets, metadata.Topo{
				SetID:   setID,
				SetName: setNameMap[setID],
				Module:  modules,
			})
		}
		curBizID := curHostAppMap[hostID]
		curData := metadata.HostBizTopo{
			BizID:   curBizID,
			BizName: appIDNameMap[curBizID],
			Set:     sets,
		}

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
			AuditType:    metadata.HostType,
			ResourceType: metadata.HostRes,
			Action:       action,
			BusinessID:   bizID,
			ResourceID:   hostID,
			ResourceName: hostIP,
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
		blog.Errorf("get host module config failed, http do error, err: %s, input: %+v, rid: %s", err.Error(), conds, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get host module config failed, http respond error, err code: %d, err msg: %s, input: %+v, rid: %s",
			result.Code, result.ErrMsg, conds, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}
	return result.Data.Info, nil
}

func (h *hostModuleLog) getInnerIP(kit *rest.Kit) ([]mapstr.MapStr, errors.CCError) {
	query := &metadata.QueryInput{
		Start:     0,
		Limit:     len(h.hostIDArr),
		Sort:      common.BKAppIDField,
		Condition: common.KvMap{common.BKHostIDField: common.KvMap{common.BKDBIN: h.hostIDArr}},
		Fields:    fmt.Sprintf("%s,%s", common.BKHostIDField, common.BKHostInnerIPField),
	}

	result, err := h.audit.clientSet.Host().GetHosts(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("get hosts failed, http do error, err: %v, input: %+v, rid: %s", err, query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get hosts failed, http respond error, err code: %d, err msg: %s, input: %+v, rid: %s",
			result.Code, result.ErrMsg, query, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	return result.Data.Info, nil
}

func (h *hostModuleLog) getModules(kit *rest.Kit, moduleIds []int64) ([]mapstr.MapStr, errors.CCError) {
	if moduleIds == nil {
		return make([]mapstr.MapStr, 0), nil
	}
	query := &metadata.QueryCondition{
		Page:      metadata.BasePage{Start: 0, Limit: common.BKNoLimit},
		Condition: mapstr.MapStr{common.BKModuleIDField: common.KvMap{common.BKDBIN: moduleIds}},
		Fields:    []string{common.BKModuleIDField, common.BKSetIDField, common.BKModuleNameField, common.BKAppIDField, common.BKOwnerIDField},
	}
	result, err := h.audit.clientSet.Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.Errorf("get modules failed, http do error, err: %v, input: %+v, rid: %s", err, query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get modules failed, http respond error, err code: %d, err msg: %s, input: %+v, rid: %s",
			result.Code, result.ErrMsg, query, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	return result.Data.Info, nil
}

func (h *hostModuleLog) getSets(kit *rest.Kit, setIDs []int64) ([]mapstr.MapStr, errors.CCError) {
	if setIDs == nil {
		return make([]mapstr.MapStr, 0), nil
	}
	query := &metadata.QueryCondition{
		Page:      metadata.BasePage{Start: 0, Limit: common.BKNoLimit},
		Condition: mapstr.MapStr{common.BKSetIDField: mapstr.MapStr{common.BKDBIN: setIDs}},
		Fields:    []string{common.BKSetNameField, common.BKSetIDField, common.BKOwnerIDField},
	}
	result, err := h.audit.clientSet.Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, query)
	if err != nil {
		blog.Errorf("get sets failed, err: %v, input: %+v, rid: %s", err, query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get sets failed, http response error, err code: %d, err msg: %s, input: %+v, rid: %s",
			result.Code, result.ErrMsg, query, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}
	return result.Data.Info, nil
}

func (h *hostModuleLog) getApps(kit *rest.Kit, appIDs []int64) ([]mapstr.MapStr, errors.CCError) {
	if appIDs == nil {
		return make([]mapstr.MapStr, 0), nil
	}
	query := &metadata.QueryCondition{
		Page:      metadata.BasePage{Start: 0, Limit: common.BKNoLimit},
		Condition: mapstr.MapStr{common.BKAppIDField: mapstr.MapStr{common.BKDBIN: appIDs}},
		Fields:    []string{common.BKAppIDField, common.BKAppNameField, common.BKOwnerIDField},
	}
	result, err := h.audit.clientSet.Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.Errorf("get business failed, http do error, err: %v, input: %+v, rid: %s", err, query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get business failed, http response error, err code: %d, err msg: %s, input: %+v, rid: %s",
			result.Code, result.ErrMsg, query, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}
	return result.Data.Info, nil
}

// GenerateAuditLog generate audit log of host module relate.
func (h *hostModuleLog) generateAuditLog(action metadata.ActionType, hostID, bizID int64, hostIP string,
	preData, curData metadata.HostBizTopo) *metadata.AuditLog {
	return &metadata.AuditLog{
		AuditType:    metadata.HostType,
		ResourceType: metadata.HostRes,
		Action:       action,
		BusinessID:   bizID,
		ResourceID:   hostID,
		ResourceName: hostIP,
		OperationDetail: &metadata.HostTransferOpDetail{
			PreData: preData,
			CurData: curData,
		},
	}
}
