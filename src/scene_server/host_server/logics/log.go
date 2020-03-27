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

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
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
	Content *metadata.Content
}

func (lgc *Logics) NewHostLog(ctx context.Context, ownerID string) *HostLog {
	return &HostLog{
		logic:   lgc,
		header:  lgc.header,
		ownerID: ownerID,
		Content: new(metadata.Content),
	}
}

func (h *HostLog) WithPrevious(ctx context.Context, hostID string, headers []metadata.Header) errors.CCError {
	var err error
	if headers != nil && len(headers) != 0 {
		h.Content.Headers = headers
	} else {
		h.Content.Headers, err = h.logic.GetHostAttributes(ctx, h.ownerID, nil)
		if err != nil {
			return err
		}
	}

	h.Content.PreData, h.ip, err = h.logic.GetHostInstanceDetails(ctx, h.ownerID, hostID)
	if err != nil {
		return err
	}

	return nil
}

func (h *HostLog) WithCurrent(ctx context.Context, hostID string) errors.CCError {
	var err error
	h.Content.CurData, h.ip, err = h.logic.GetHostInstanceDetails(ctx, h.ownerID, hostID)
	if err != nil {
		return err
	}

	return nil
}

func (h *HostLog) AuditLog(ctx context.Context, hostID int64) metadata.SaveAuditLogParams {
	return metadata.SaveAuditLogParams{
		ID:      hostID,
		Content: h.Content,
		ExtKey:  h.ip,
	}
}

func (h *HostLog) GetContent(hostID int64) *metadata.Content {
	return h.Content
}

type HostModuleLog struct {
	logic     *Logics
	header    http.Header
	instIDArr []int64
	pre       []metadata.ModuleHost
	cur       []metadata.ModuleHost
	hostInfos []mapstr.MapStr
	desc      string
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

func (h *HostModuleLog) SaveAudit(ctx context.Context, appID int64, user, desc string) errors.CCError {
	if err := h.WithCurrent(ctx); err != nil {
		return err
	}

	var setIDs, moduleIDs, appIDs []int64
	preMap := make(map[int64]map[int64]interface{})
	curMap := make(map[int64]map[int64]interface{})

	for _, val := range h.pre {
		if _, ok := preMap[val.HostID]; false == ok {
			preMap[val.HostID] = make(map[int64]interface{})
		}
		preMap[val.HostID][val.ModuleID] = val
		setIDs = append(setIDs, val.SetID)
		moduleIDs = append(moduleIDs, val.ModuleID)
		appIDs = append(appIDs, val.AppID)
	}

	for _, val := range h.cur {
		if _, ok := curMap[val.HostID]; false == ok {
			curMap[val.HostID] = make(map[int64]interface{})
		}
		curMap[val.HostID][val.ModuleID] = val
		setIDs = append(setIDs, val.SetID)
		moduleIDs = append(moduleIDs, val.ModuleID)
	}

	modules, err := h.getModules(ctx, moduleIDs)
	if err != nil {
		return err
	}

	sets, err := h.getSets(ctx, setIDs)
	if err != nil {
		return err
	}

	setMap := make(map[int64]metadata.Ref, 0)
	for _, setInfo := range sets {
		instID, err := util.GetInt64ByInterface(setInfo[common.BKSetIDField])
		if err != nil {
			return err
		}
		setMap[instID] = metadata.Ref{
			RefID:   instID,
			RefName: setInfo[common.BKSetNameField].(string),
		}
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

	type ModuleRef struct {
		metadata.Ref
		Set     []interface{} `json:"set"`
		appID   int64
		ownerID string
	}
	moduleMap := make(map[int64]ModuleRef, 0)
	for _, moduleInfo := range modules {
		mID, err := util.GetInt64ByInterface(moduleInfo[common.BKModuleIDField])
		if err != nil {
			return err
		}
		sID, err := util.GetInt64ByInterface(moduleInfo[common.BKSetIDField])
		if err != nil {
			return err
		}

		appID, err := moduleInfo.Int64(common.BKAppIDField)
		if err != nil {
			return h.logic.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
		}
		moduleRef := ModuleRef{}
		moduleRef.Set = append(moduleRef.Set, setMap[sID])
		moduleRef.RefID = mID
		moduleRef.RefName = moduleInfo[common.BKModuleNameField].(string)
		moduleRef.appID = appID
		moduleRef.ownerID = moduleInfo[common.BKOwnerIDField].(string)
		moduleMap[mID] = moduleRef
	}
	moduleReName := "module"
	setRefName := "set"
	headers := []metadata.Header{
		{PropertyID: moduleReName, PropertyName: "module"},
		{PropertyID: setRefName, PropertyName: "set"},
		{PropertyID: common.BKAppIDField, PropertyName: "business ID"},
		{PropertyID: common.BKAppNameField, PropertyName: "business Name"},
	}

	if len(desc) != 0 {
		h.desc = desc
	} else {
		h.desc = "host module change"
	}

	logs := make([]metadata.SaveAuditLogParams, 0)
	for _, host := range h.hostInfos {
		instID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			return err
		}

		preModule := make([]interface{}, 0)
		var preApp int64
		for moduleID := range preMap[instID] {
			preModule = append(preModule, moduleMap[moduleID])
			preApp = moduleMap[moduleID].appID
		}

		curModule := make([]interface{}, 0)
		var curApp int64

		for moduleID := range curMap[instID] {
			curModule = append(curModule, moduleMap[moduleID])
			curApp = moduleMap[moduleID].appID
		}

		logs = append(logs, metadata.SaveAuditLogParams{
			ID:    instID,
			Model: common.BKInnerObjIDHost,
			Content: metadata.Content{
				PreData: common.KvMap{moduleReName: preModule, common.BKAppIDField: preApp, common.BKAppNameField: appIDNameMap[preApp]},
				CurData: common.KvMap{moduleReName: curModule, common.BKAppIDField: curApp, common.BKAppNameField: appIDNameMap[curApp]},
				Headers: headers,
			},
			OpDesc: h.desc,
			OpType: auditoplog.AuditOpTypeHostModule,
			ExtKey: host[common.BKHostInnerIPField].(string),
			BizID:  appID,
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
		Limit:     metadata.SearchLimit{Offset: 0, Limit: common.BKNoLimit},
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
		Limit:     metadata.SearchLimit{Offset: 0, Limit: common.BKNoLimit},
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
		Limit:     metadata.SearchLimit{Offset: 0, Limit: common.BKNoLimit},
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
