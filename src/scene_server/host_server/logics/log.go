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

func (lgc *Logics) NewHostLog(pheader http.Header, ownerID string) *HostLog {
	return &HostLog{
		logic:   lgc,
		header:  pheader,
		ownerID: ownerID,
		Content: new(metadata.Content),
	}
}

func (h *HostLog) WithPrevious(hostID string, headers []metadata.Header) error {
	var err error
	if headers != nil || len(headers) != 0 {
		h.Content.Headers = headers
	} else {
		h.Content.Headers, err = h.logic.GetHostAttributes(h.ownerID, h.header)
		if err != nil {
			return err
		}
	}

	h.Content.PreData, h.ip, err = h.logic.GetHostInstanceDetails(h.header, h.ownerID, hostID)
	if err != nil {
		return err
	}

	return nil
}

func (h *HostLog) WithCurrent(hostID string) error {
	var err error
	h.Content.CurData, h.ip, err = h.logic.GetHostInstanceDetails(h.header, h.ownerID, hostID)
	if err != nil {
		return err
	}

	return nil
}

func (h *HostLog) AuditLog(hostID int64) *auditoplog.AuditLogExt {
	return &auditoplog.AuditLogExt{
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
	instID    []int64
	pre       []metadata.ModuleHost
	cur       []metadata.ModuleHost
	hostInfos []map[string]interface{}
	desc      string
}

func (lgc *Logics) NewHostModuleLog(header http.Header, instID []int64) *HostModuleLog {
	return &HostModuleLog{
		logic:  lgc,
		instID: instID,
		pre:    make([]metadata.ModuleHost, 0),
		cur:    make([]metadata.ModuleHost, 0),
		header: header,
	}
}

func (h *HostModuleLog) WithPrevious() error {
	var err error
	h.pre, err = h.getHostModuleConfig()
	if err != nil {
		return err
	}

	h.hostInfos, err = h.getInnerIP()
	if err != nil {
		return err
	}

	return nil
}

func (h *HostModuleLog) WithCurrent() error {
	var err error
	h.cur, err = h.getHostModuleConfig()
	if err != nil {
		return err
	}
	return nil
}

func (h *HostModuleLog) SaveAudit(appID, user, desc string) error {
	if err := h.WithCurrent(); err != nil {
		return err
	}

	var setIDs, moduleIDs []int64
	preMap := make(map[int64]map[int64]interface{})
	curMap := make(map[int64]map[int64]interface{})

	for _, val := range h.pre {
		if _, ok := preMap[val.HostID]; false == ok {
			preMap[val.HostID] = make(map[int64]interface{})
		}
		preMap[val.HostID][val.ModuleID] = val
		setIDs = append(setIDs, val.SetID)
		moduleIDs = append(moduleIDs, val.ModuleID)
	}

	for _, val := range h.cur {
		if _, ok := curMap[val.HostID]; false == ok {
			curMap[val.HostID] = make(map[int64]interface{})
		}
		curMap[val.HostID][val.ModuleID] = val
		setIDs = append(setIDs, val.SetID)
		moduleIDs = append(moduleIDs, val.ModuleID)
	}

	modules, err := h.getModules(moduleIDs)
	if err != nil {
		return err
	}

	sets, err := h.getSets(setIDs)
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
	type ModuleRef struct {
		metadata.Ref
		Set     []interface{} `json:"set"`
		appID   interface{}
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
		moduleRef := ModuleRef{}
		moduleRef.Set = append(moduleRef.Set, setMap[sID])
		moduleRef.RefID = mID
		moduleRef.RefName = moduleInfo[common.BKModuleNameField].(string)
		moduleRef.appID = moduleInfo[common.BKAppIDField]
		moduleRef.ownerID = moduleInfo[common.BKOwnerIDField].(string)
		moduleMap[mID] = moduleRef
	}
	moduleReName := "module"
	setRefName := "set"
	headers := []metadata.Header{
		{PropertyID: moduleReName, PropertyName: "module"},
		{PropertyID: setRefName, PropertyName: "app"},
		{PropertyID: common.BKAppIDField, PropertyName: "business ID"},
	}

	logs := make([]auditoplog.AuditLogExt, 0)

	var ownerID string
	for _, host := range h.hostInfos {
		instID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			return err
		}
		log := auditoplog.AuditLogExt{ID: instID}
		log.ExtKey = host[common.BKHostInnerIPField].(string)

		preModule := make([]interface{}, 0)
		var preApp interface{}
		for moduleID, _ := range preMap[instID] {
			preModule = append(preModule, moduleMap[moduleID])
			preApp = moduleMap[moduleID].appID
			ownerID = moduleMap[moduleID].ownerID
		}

		curModule := make([]interface{}, 0)
		var curApp interface{}

		for moduleID, _ := range curMap[instID] {
			curModule = append(curModule, moduleMap[moduleID])
			curApp = moduleMap[moduleID].appID
			ownerID = moduleMap[moduleID].ownerID
		}

		log.Content = metadata.Content{
			PreData: common.KvMap{moduleReName: preModule, common.BKAppIDField: preApp},
			CurData: common.KvMap{moduleReName: curModule, common.BKAppIDField: curApp},
			Headers: headers,
		}
		logs = append(logs, log)

	}

	if len(desc) != 0 {
		h.desc = desc
	} else {
		h.desc = "host module change"
	}
	data := common.KvMap{common.BKContentField: logs, common.BKOpDescField: h.desc, common.BKOpTypeField: auditoplog.AuditOpTypeHostModule}
	result, err := h.logic.CoreAPI.AuditController().AddHostLogs(context.Background(), ownerID, appID, user, h.header, data)
	if err != nil || (err == nil && !result.Result) {
		return fmt.Errorf("%v, %v", err, result.ErrMsg)
	}
	return nil
}

func (h *HostModuleLog) getHostModuleConfig() ([]metadata.ModuleHost, error) {
	conds := map[string][]int64{common.BKHostIDField: h.instID}
	result, err := h.logic.CoreAPI.HostController().Module().GetModulesHostConfig(context.Background(), h.header, conds)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("%v, %v", err, result.ErrMsg)
	}
	return result.Data, nil
}

func (h *HostModuleLog) getInnerIP() ([]map[string]interface{}, error) {
	query := &metadata.QueryInput{
		Start:     0,
		Limit:     1,
		Sort:      common.BKAppIDField,
		Condition: common.KvMap{common.BKHostIDField: common.KvMap{common.BKDBIN: h.instID}},
		Fields:    fmt.Sprintf("%s,%s", common.BKHostIDField, common.BKHostInnerIPField),
	}

	result, err := h.logic.CoreAPI.HostController().Host().GetHosts(context.Background(), h.header, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("%v, %v", err, result.ErrMsg)
	}

	return result.Data.Info, nil
}

func (h *HostModuleLog) getModules(moduleIds []int64) ([]mapstr.MapStr, error) {
	if moduleIds == nil {
		return make([]mapstr.MapStr, 0), nil
	}
	query := &metadata.QueryInput{
		Start:     0,
		Limit:     common.BKNoLimit,
		Condition: common.KvMap{common.BKModuleIDField: common.KvMap{common.BKDBIN: moduleIds}},
		Fields:    fmt.Sprintf("%s,%s,%s,%s,%s", common.BKModuleIDField, common.BKSetIDField, common.BKModuleNameField, common.BKAppIDField, common.BKOwnerIDField),
	}

	result, err := h.logic.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, h.header, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get modules with id failed, err: %v, result err: %s", err, result.ErrMsg)
	}

	return result.Data.Info, nil
}

func (h *HostModuleLog) getSets(setIDs []int64) ([]mapstr.MapStr, error) {
	if setIDs == nil {
		return make([]mapstr.MapStr, 0), nil
	}
	query := &metadata.QueryInput{
		Start:     0,
		Limit:     common.BKNoLimit,
		Condition: common.KvMap{common.BKSetIDField: common.KvMap{common.BKDBIN: setIDs}},
		Fields:    fmt.Sprintf("%s,%s,%s", common.BKSetNameField, common.BKSetIDField, common.BKOwnerIDField),
	}

	result, err := h.logic.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDSet, h.header, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get modules with id failed, err: %v, result err: %s", err, result.ErrMsg)
	}

	return result.Data.Info, nil
}
