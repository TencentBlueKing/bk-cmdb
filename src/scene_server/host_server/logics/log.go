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

package logics

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type HostLog struct {
	logic   *Logics
	header  http.Header
	ownerID string
	ip      string
	Content *metadata.BasicContent
}

func (lgc *Logics) NewHostLog(ctx context.Context, ownerID string) *HostLog {
	return &HostLog{
		logic:   lgc,
		header:  lgc.header,
		ownerID: ownerID,
		Content: new(metadata.BasicContent),
	}
}

func (h *HostLog) WithPrevious(ctx context.Context, hostID int64, properties []metadata.Property) errors.CCError {
	var err error
	if len(h.Content.Properties) == 0 {
		if properties != nil && len(properties) != 0 {
			h.Content.Properties = properties
		} else {
			h.Content.Properties, err = h.logic.GetHostAttributes(ctx, h.ownerID, nil)
			if err != nil {
				return err
			}
		}
	}

	id := strconv.FormatInt(hostID, 10)
	h.Content.PreData, h.ip, err = h.logic.GetHostInstanceDetails(ctx, h.ownerID, id)
	if err != nil {
		return err
	}

	return nil
}

func (h *HostLog) WithCurrent(ctx context.Context, hostID int64, properties []metadata.Property) errors.CCError {
	var err error
	if len(h.Content.Properties) == 0 {
		if properties != nil && len(properties) != 0 {
			h.Content.Properties = properties
		} else {
			h.Content.Properties, err = h.logic.GetHostAttributes(ctx, h.ownerID, nil)
			if err != nil {
				return err
			}
		}
	}

	id := strconv.FormatInt(hostID, 10)
	h.Content.CurData, h.ip, err = h.logic.GetHostInstanceDetails(ctx, h.ownerID, id)
	if err != nil {
		return err
	}

	return nil
}

func (h *HostLog) AuditLog(ctx context.Context, hostID int64, bizID int64, action metadata.ActionType) metadata.AuditLog {
	return metadata.AuditLog{
		AuditType:    metadata.HostType,
		ResourceType: metadata.HostRes,
		Action:       action,
		OperationDetail: &metadata.InstanceOpDetail{
			BasicOpDetail: metadata.BasicOpDetail{
				BusinessID:   bizID,
				ResourceID:   hostID,
				ResourceName: h.ip,
				Details:      h.Content,
			},
			ModelID: common.BKInnerObjIDHost,
		},
	}
}

type HostModuleLog struct {
	logic     *Logics
	header    http.Header
	instIDArr []int64
	pre       []metadata.ModuleHost
	cur       []metadata.ModuleHost
	hostInfos []mapstr.MapStr
}

func (lgc *Logics) NewHostModuleLog(instID []int64) *HostModuleLog {
	return &HostModuleLog{
		logic:     lgc,
		instIDArr: instID,
		pre:       make([]metadata.ModuleHost, 0),
		cur:       make([]metadata.ModuleHost, 0),
		header:    lgc.header,
	}
}

func (h *HostModuleLog) WithPrevious(ctx context.Context) errors.CCError {
	var err error
	h.pre, err = h.getHostModuleConfig(ctx)
	if err != nil {
		return err
	}

	h.hostInfos, err = h.getInnerIP(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (h *HostModuleLog) WithCurrent(ctx context.Context) errors.CCError {
	var err error
	h.cur, err = h.getHostModuleConfig(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (h *HostModuleLog) SaveAudit(ctx context.Context) errors.CCError {
	if err := h.WithCurrent(ctx); err != nil {
		return err
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

	modules, err := h.getModules(ctx, moduleIDs)
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

	sets, err := h.getSets(ctx, setIDs)
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

	preHostRelationMap := make(map[int64]map[int64]map[int64][]metadata.Module)
	for _, val := range h.pre {
		if _, ok := preHostRelationMap[val.HostID]; false == ok {
			preHostRelationMap[val.HostID] = make(map[int64]map[int64][]metadata.Module)
		}
		if _, ok := preHostRelationMap[val.HostID][val.AppID]; false == ok {
			preHostRelationMap[val.HostID][val.AppID] = make(map[int64][]metadata.Module)
		}
		preHostRelationMap[val.HostID][val.AppID][val.SetID] = append(preHostRelationMap[val.HostID][val.AppID][val.SetID], metadata.Module{ModuleID: val.ModuleID, ModuleName: moduleNameMap[val.ModuleID]})
	}
	curHostRelationMap := make(map[int64]map[int64]map[int64][]metadata.Module)
	for _, val := range h.cur {
		if _, ok := curHostRelationMap[val.HostID]; false == ok {
			curHostRelationMap[val.HostID] = make(map[int64]map[int64][]metadata.Module)
		}
		if _, ok := curHostRelationMap[val.HostID][val.AppID]; false == ok {
			curHostRelationMap[val.HostID][val.AppID] = make(map[int64][]metadata.Module)
		}
		curHostRelationMap[val.HostID][val.AppID][val.SetID] = append(curHostRelationMap[val.HostID][val.AppID][val.SetID], metadata.Module{ModuleID: val.ModuleID, ModuleName: moduleNameMap[val.ModuleID]})
	}

	appInfoArr, err := h.getApps(ctx, appIDs)
	if err != nil {
		return err
	}
	appIDNameMap := make(map[int64]string, 0)
	for _, appInfo := range appInfoArr {
		instID, err := appInfo.Int64(common.BKAppIDField)
		if err != nil {
			blog.ErrorJSON("appInfo get biz id err:%s, appInfo: %s, rid:%s", err.Error(), appInfo, h.logic.rid)
			return h.logic.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
		}
		name, err := appInfo.String(common.BKAppNameField)
		if err != nil {
			blog.ErrorJSON("appInfo get biz name err:%s, appInfo: %s, rid:%s", err.Error(), appInfo, h.logic.rid)
			return h.logic.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
		}

		appIDNameMap[instID] = name
	}

	logs := make([]metadata.AuditLog, 0)
	for _, host := range h.hostInfos {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			return err
		}

		hostIP, err := host.String(common.BKHostInnerIPField)
		if err != nil {
			return err
		}

		preData := make([]metadata.HostBizTopo, 0)
		for bizID, setModuleMap := range preHostRelationMap[hostID] {
			sets := make([]metadata.Topo, 0)
			for setID, modules := range setModuleMap {
				sets = append(sets, metadata.Topo{
					SetID:   setID,
					SetName: setNameMap[setID],
					Module:  modules,
				})
			}
			preData = append(preData, metadata.HostBizTopo{
				BizID:   bizID,
				BizName: appIDNameMap[bizID],
				Set:     sets,
			})
		}

		curData := make([]metadata.HostBizTopo, 0)
		for bizID, setModuleMap := range curHostRelationMap[hostID] {
			sets := make([]metadata.Topo, 0)
			for setID, modules := range setModuleMap {
				sets = append(sets, metadata.Topo{
					SetID:   setID,
					SetName: setNameMap[setID],
					Module:  modules,
				})
			}
			curData = append(curData, metadata.HostBizTopo{
				BizID:   bizID,
				BizName: appIDNameMap[bizID],
				Set:     sets,
			})
		}

		logs = append(logs, metadata.AuditLog{
			AuditType:    metadata.HostType,
			ResourceType: metadata.HostRes,
			Action:       metadata.AuditTransferHost,
			OperationDetail: &metadata.HostTransferOpDetail{
				HostID:      hostID,
				HostInnerIP: hostIP,
				PreData:     preData,
				CurData:     curData,
			},
		})
	}

	result, err := h.logic.CoreAPI.CoreService().Audit().SaveAuditLog(context.Background(), h.header, logs...)
	if err != nil {
		blog.Errorf("AddHostLogs http do error, err:%s,input:%+v,rid:%s", err.Error(), logs, h.logic.rid)
		return h.logic.ccErr.Error(common.CCErrAuditSaveLogFailed)
	}
	if !result.Result {
		blog.Errorf("AddHostLogs  http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, logs, h.logic.rid)
		return h.logic.ccErr.New(result.Code, result.ErrMsg)
	}

	return nil
}

func (h *HostModuleLog) getHostModuleConfig(ctx context.Context) ([]metadata.ModuleHost, errors.CCError) {
	conds := &metadata.HostModuleRelationRequest{
		HostIDArr: h.instIDArr,
	}
	result, err := h.logic.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx, h.header, conds)
	if err != nil {
		blog.Errorf("getHostModuleConfig http do error, err:%s,input:%+v,rid:%s", err.Error(), conds, h.logic.rid)
		return nil, h.logic.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("getHostModuleConfig http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, conds, h.logic.rid)
		return nil, h.logic.ccErr.New(result.Code, result.ErrMsg)
	}
	return result.Data.Info, nil
}

func (h *HostModuleLog) getInnerIP(ctx context.Context) ([]mapstr.MapStr, errors.CCError) {
	query := &metadata.QueryInput{
		Start:     0,
		Limit:     len(h.instIDArr),
		Sort:      common.BKAppIDField,
		Condition: common.KvMap{common.BKHostIDField: common.KvMap{common.BKDBIN: h.instIDArr}},
		Fields:    fmt.Sprintf("%s,%s", common.BKHostIDField, common.BKHostInnerIPField),
	}

	result, err := h.logic.CoreAPI.CoreService().Host().GetHosts(ctx, h.header, query)
	if err != nil {
		blog.Errorf("GetHosts http do error, err:%s,input:%+v,rid:%s", err.Error(), query, h.logic.rid)
		return nil, h.logic.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetHosts http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, h.logic.rid)
		return nil, h.logic.ccErr.New(result.Code, result.ErrMsg)
	}

	return result.Data.Info, nil
}

func (h *HostModuleLog) getModules(ctx context.Context, moduleIds []int64) ([]mapstr.MapStr, errors.CCError) {
	if moduleIds == nil {
		return make([]mapstr.MapStr, 0), nil
	}
	query := &metadata.QueryCondition{
		Page:      metadata.BasePage{Start: 0, Limit: common.BKNoLimit},
		Condition: mapstr.MapStr{common.BKModuleIDField: common.KvMap{common.BKDBIN: moduleIds}},
		Fields:    []string{common.BKModuleIDField, common.BKSetIDField, common.BKModuleNameField, common.BKAppIDField, common.BKOwnerIDField},
	}
	result, err := h.logic.CoreAPI.CoreService().Instance().ReadInstance(ctx, h.header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.Errorf("getModules http do error, err:%s,input:%+v,rid:%s", err.Error(), query, h.logic.rid)
		return nil, h.logic.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("getModules http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, h.logic.rid)
		return nil, h.logic.ccErr.New(result.Code, result.ErrMsg)
	}

	return result.Data.Info, nil
}

func (h *HostModuleLog) getSets(ctx context.Context, setIDs []int64) ([]mapstr.MapStr, errors.CCError) {
	if setIDs == nil {
		return make([]mapstr.MapStr, 0), nil
	}
	query := &metadata.QueryCondition{
		Page:      metadata.BasePage{Start: 0, Limit: common.BKNoLimit},
		Condition: mapstr.MapStr{common.BKSetIDField: mapstr.MapStr{common.BKDBIN: setIDs}},
		Fields:    []string{common.BKSetNameField, common.BKSetIDField, common.BKOwnerIDField},
	}
	result, err := h.logic.CoreAPI.CoreService().Instance().ReadInstance(ctx, h.header, common.BKInnerObjIDSet, query)
	if err != nil {
		blog.Errorf("getSets http do error, err:%s,input:%+v,rid:%s", err.Error(), query, h.logic.rid)
		return nil, h.logic.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("getSets http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, h.logic.rid)
		return nil, h.logic.ccErr.New(result.Code, result.ErrMsg)
	}
	return result.Data.Info, nil
}

func (h *HostModuleLog) getApps(ctx context.Context, appIDs []int64) ([]mapstr.MapStr, errors.CCError) {
	if appIDs == nil {
		return make([]mapstr.MapStr, 0), nil
	}
	query := &metadata.QueryCondition{
		Page:      metadata.BasePage{Start: 0, Limit: common.BKNoLimit},
		Condition: mapstr.MapStr{common.BKAppIDField: mapstr.MapStr{common.BKDBIN: appIDs}},
		Fields:    []string{common.BKAppIDField, common.BKAppNameField, common.BKOwnerIDField},
	}
	result, err := h.logic.CoreAPI.CoreService().Instance().ReadInstance(ctx, h.header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.Errorf("getApps http do error, err:%s,input:%+v,rid:%s", err.Error(), query, h.logic.rid)
		return nil, h.logic.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("getApps http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, h.logic.rid)
		return nil, h.logic.ccErr.New(result.Code, result.ErrMsg)
	}
	return result.Data.Info, nil
}
