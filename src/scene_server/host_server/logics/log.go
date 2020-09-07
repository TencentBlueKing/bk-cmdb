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
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type HostLog struct {
	kit     *rest.Kit
	logic   *Logics
	header  http.Header
	ownerID string
	ip      string
	Content *metadata.BasicContent
}

func (lgc *Logics) NewHostLog(kit *rest.Kit) *HostLog {
	return &HostLog{
		kit:     kit,
		logic:   lgc,
		header:  kit.Header,
		ownerID: kit.SupplierAccount,
		Content: new(metadata.BasicContent),
	}
}

func (h *HostLog) WithPrevious(hostID int64) errors.CCError {
	var err error

	h.Content.PreData, h.ip, err = h.logic.GetHostInstanceDetails(h.kit, hostID)
	if err != nil {
		return err
	}

	return nil
}

func (h *HostLog) WithCurrent(hostID int64) errors.CCError {
	var err error

	h.Content.CurData, h.ip, err = h.logic.GetHostInstanceDetails(h.kit, hostID)
	if err != nil {
		return err
	}

	return nil
}

func (h *HostLog) AuditLog(hostID int64, bizID int64, action metadata.ActionType) (metadata.AuditLog, error) {
	return metadata.AuditLog{
		AuditType:    metadata.HostType,
		ResourceType: metadata.HostRes,
		Action:       action,
		BusinessID:   bizID,
		ResourceID:   hostID,
		ResourceName: h.ip,
		OperationDetail: &metadata.InstanceOpDetail{
			BasicOpDetail: metadata.BasicOpDetail{
				Details: h.Content,
			},
			ModelID: common.BKInnerObjIDHost,
		},
	}, nil
}

type HostModuleLog struct {
	kit       *rest.Kit
	logic     *Logics
	header    http.Header
	hostIDArr []int64
	pre       []metadata.ModuleHost
	cur       []metadata.ModuleHost
}

func (lgc *Logics) NewHostModuleLog(kit *rest.Kit, hostID []int64) *HostModuleLog {
	return &HostModuleLog{
		kit:	   kit,
		logic:     lgc,
		hostIDArr: hostID,
		header:    kit.Header,
	}
}

func (h *HostModuleLog) WithPrevious(ctx context.Context) errors.CCError {
	if h.pre != nil {
		return nil
	}
	var err error
	h.pre, err = h.getHostModuleConfig(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (h *HostModuleLog) WithCurrent(ctx context.Context) errors.CCError {
	if h.cur != nil {
		return nil
	}
	var err error
	h.cur, err = h.getHostModuleConfig(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (h *HostModuleLog) SaveAudit(ctx context.Context) errors.CCError {
	hostInfos, err := h.getInnerIP(ctx)
	if err != nil {
		return err
	}

	if err := h.WithCurrent(ctx); err != nil {
		return err
	}

	defaultBizID, err := h.logic.GetDefaultAppID(h.kit)
	if err != nil {
		blog.ErrorJSON("save audit GetDefaultAppID failed, err: %s, rid: %s", err, h.kit.Rid)
		return h.kit.CCError.Error(common.CCErrAuditSaveLogFailed)
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

	appInfoArr, err := h.getApps(ctx, appIDs)
	if err != nil {
		return err
	}
	appIDNameMap := make(map[int64]string, 0)
	for _, appInfo := range appInfoArr {
		bizID, err := appInfo.Int64(common.BKAppIDField)
		if err != nil {
			blog.ErrorJSON("appInfo get biz id err:%s, appInfo: %s, rid:%s", err.Error(), appInfo, h.kit.Rid)
			return h.kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
		}
		name, err := appInfo.String(common.BKAppNameField)
		if err != nil {
			blog.ErrorJSON("appInfo get biz name err:%s, appInfo: %s, rid:%s", err.Error(), appInfo, h.kit.Rid)
			return h.kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
		}

		appIDNameMap[bizID] = name
	}

	logs := make([]metadata.AuditLog, 0)
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

	result, err := h.logic.CoreAPI.CoreService().Audit().SaveAuditLog(context.Background(), h.header, logs...)
	if err != nil {
		blog.Errorf("AddHostLogs http do error, err:%s,input:%+v,rid:%s", err.Error(), logs, h.kit.Rid)
		return h.kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}
	if !result.Result {
		blog.Errorf("AddHostLogs  http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, logs, h.kit.Rid)
		return h.kit.CCError.New(result.Code, result.ErrMsg)
	}

	return nil
}

func (h *HostModuleLog) getHostModuleConfig(ctx context.Context) ([]metadata.ModuleHost, errors.CCError) {
	conds := &metadata.HostModuleRelationRequest{
		HostIDArr: h.hostIDArr,
		Fields:    []string{common.BKAppIDField, common.BKSetIDField, common.BKModuleIDField, common.BKHostIDField},
	}
	result, err := h.logic.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx, h.header, conds)
	if err != nil {
		blog.Errorf("getHostModuleConfig http do error, err:%s,input:%+v,rid:%s", err.Error(), conds, h.kit.Rid)
		return nil, h.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("getHostModuleConfig http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, conds, h.kit.Rid)
		return nil, h.kit.CCError.New(result.Code, result.ErrMsg)
	}
	return result.Data.Info, nil
}

func (h *HostModuleLog) getInnerIP(ctx context.Context) ([]mapstr.MapStr, errors.CCError) {
	query := &metadata.QueryInput{
		Start:     0,
		Limit:     len(h.hostIDArr),
		Sort:      common.BKAppIDField,
		Condition: common.KvMap{common.BKHostIDField: common.KvMap{common.BKDBIN: h.hostIDArr}},
		Fields:    fmt.Sprintf("%s,%s", common.BKHostIDField, common.BKHostInnerIPField),
	}

	result, err := h.logic.CoreAPI.CoreService().Host().GetHosts(ctx, h.header, query)
	if err != nil {
		blog.Errorf("GetHosts http do error, err:%s,input:%+v,rid:%s", err.Error(), query, h.kit.Rid)
		return nil, h.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetHosts http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, h.kit.Rid)
		return nil, h.kit.CCError.New(result.Code, result.ErrMsg)
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
		blog.Errorf("getModules http do error, err:%s,input:%+v,rid:%s", err.Error(), query, h.kit.Rid)
		return nil, h.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("getModules http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, h.kit.Rid)
		return nil, h.kit.CCError.New(result.Code, result.ErrMsg)
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
		blog.Errorf("getSets http do error, err:%s,input:%+v,rid:%s", err.Error(), query, h.kit.Rid)
		return nil, h.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("getSets http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, h.kit.Rid)
		return nil, h.kit.CCError.New(result.Code, result.ErrMsg)
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
		blog.Errorf("getApps http do error, err:%s,input:%+v,rid:%s", err.Error(), query, h.kit.Rid)
		return nil, h.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("getApps http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, h.kit.Rid)
		return nil, h.kit.CCError.New(result.Code, result.ErrMsg)
	}
	return result.Data.Info, nil
}

type CloudAreaAuditLog struct {
	kit 	       *rest.Kit
	logic          *Logics
	header         http.Header
	ownerID        string
	MultiCloudArea map[int64]*SingleCloudArea
}

type SingleCloudArea struct {
	CloudName string
	PreData   map[string]interface{}
	CurData   map[string]interface{}
}

func (lgc *Logics) NewCloudAreaLog(kit *rest.Kit) *CloudAreaAuditLog {
	return &CloudAreaAuditLog{
		kit:			kit,
		logic:          lgc,
		header:         kit.Header,
		ownerID:        kit.SupplierAccount,
		MultiCloudArea: make(map[int64]*SingleCloudArea),
	}
}

func (c *CloudAreaAuditLog) WithPrevious(ctx context.Context, platIDs ...int64) errors.CCError {
	return c.buildAuditLogData(ctx, true, false, platIDs...)
}

func (c *CloudAreaAuditLog) WithCurrent(ctx context.Context, platIDs ...int64) errors.CCError {
	return c.buildAuditLogData(ctx, false, true, platIDs...)
}

func (c *CloudAreaAuditLog) buildAuditLogData(ctx context.Context, withPrevious, withCurrent bool, platIDs ...int64) errors.CCError {
	var err error
	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKCloudIDField: mapstr.MapStr{common.BKDBIN: platIDs}},
	}
	
	res, err := c.logic.CoreAPI.CoreService().Instance().ReadInstance(ctx, c.header, common.BKInnerObjIDPlat, query)
	if nil != err {
		return err
	}
	
	if len(res.Data.Info) <= 0 {
		return errors.New(common.CCErrTopoCloudNotFound, "")
	}

	for _, data := range res.Data.Info {
		cloudID, err := data.Int64(common.BKCloudIDField)
		if err != nil {
			return err
		}

		cloudName, err := data.String(common.BKCloudNameField)
		if err != nil {
			return err
		}

		if c.MultiCloudArea[cloudID] == nil {
			c.MultiCloudArea[cloudID] = new(SingleCloudArea)
		}

		c.MultiCloudArea[cloudID].CloudName = cloudName

		if withPrevious {
			c.MultiCloudArea[cloudID].PreData = data
		}

		if withCurrent {
			c.MultiCloudArea[cloudID].CurData = data
		}
	}

	return nil
}

func (c *CloudAreaAuditLog) SaveAuditLog(ctx context.Context, action metadata.ActionType) errors.CCError {
	logs := make([]metadata.AuditLog, 0)
	for cloudID, cloudarea := range c.MultiCloudArea {
		logs = append(logs, metadata.AuditLog{
			AuditType:    metadata.CloudResourceType,
			ResourceType: metadata.CloudAreaRes,
			Action:       action,
			ResourceID:   cloudID,
			ResourceName: cloudarea.CloudName,
			OperationDetail: &metadata.InstanceOpDetail{
				BasicOpDetail: metadata.BasicOpDetail{

					Details: &metadata.BasicContent{
						PreData:    cloudarea.PreData,
						CurData:    cloudarea.CurData,
					},
				},
				ModelID: common.BKInnerObjIDPlat,
			},
		})
	}

	auditResult, err := c.logic.CoreAPI.CoreService().Audit().SaveAuditLog(ctx, c.header, logs...)
	if err != nil {
		blog.ErrorJSON("SaveAuditLog add cloud area audit log failed, err: %s, result: %+v,rid:%s", err, auditResult, c.kit.Rid)
		return c.kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}
	if auditResult.Result != true {
		blog.ErrorJSON("SaveAuditLog add cloud area audit log failed, err: %s, result: %s,rid:%s", err, auditResult, c.kit.Rid)
		return c.kit.CCError.Errorf(common.CCErrAuditSaveLogFailed)
	}

	return nil
}
