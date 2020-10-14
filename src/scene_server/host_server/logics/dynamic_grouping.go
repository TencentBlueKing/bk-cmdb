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
	"net/http"
	"sort"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	hostParse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
)

/* reuse old part codes of hostsearch.go */

// ExecuteHostDynamicGroup searches hosts base on conditions without filling topology informations.
func (lgc *Logics) ExecuteHostDynamicGroup(kit *rest.Kit, data *metadata.HostCommonSearch,
	fields []string, disableCounter bool) (*metadata.SearchHost, error) {

	// create search host action instance.
	executor := NewHostDynamicGroupExecutor(kit, lgc, data, fields, disableCounter)

	hostInfos, count, err := executor.Execute()
	if err != nil {
		return nil, err
	}
	if disableCounter {
		count = 0
	}

	return &metadata.SearchHost{Count: count, Info: hostInfos}, nil
}

// HostDynamicGroupExecutor handle host dynamic group action.
type HostDynamicGroupExecutor struct {
	kit *rest.Kit
	lgc *Logics

	ctx    context.Context
	ccErr  errors.DefaultCCErrorIf
	ccRid  string
	header http.Header

	// host search params and conditions.
	params *metadata.HostCommonSearch
	conds  searchHostConds

	idArr        searchHostIDArr
	cacheInfoMap searchHostInfoMapCache

	// final search results.
	total          int
	hosts          []hostInfoStruct
	fields         []string
	disableCounter bool

	isNotFound bool
	needPaged  bool
}

// NewHostDynamicGroupExecutor creates a new HostDynamicGroupExecutor object.
func NewHostDynamicGroupExecutor(kit *rest.Kit, lgc *Logics, params *metadata.HostCommonSearch,
	fileds []string, disableCounter bool) *HostDynamicGroupExecutor {

	executor := &HostDynamicGroupExecutor{
		kit:            kit,
		lgc:            lgc,
		ctx:            kit.Ctx,
		ccErr:          kit.CCError,
		ccRid:          kit.Rid,
		header:         kit.Header,
		params:         params,
		idArr:          searchHostIDArr{},
		fields:         fileds,
		disableCounter: disableCounter,
	}
	executor.conds.objectCondMap = make(map[string][]metadata.ConditionItem)

	return executor
}

// Execute executes host dynamic group.
func (e *HostDynamicGroupExecutor) Execute() ([]mapstr.MapStr, int, error) {
	// parse conditions.
	e.parseCondition()

	// search host with conditions.
	if err := e.searchHostByConds(); err != nil {
		return nil, 0, err
	}
	result, count := e.buildSearchResult()

	return result, count, nil
}

func (e *HostDynamicGroupExecutor) parseCondition() {
	for _, cond := range e.params.Condition {
		switch cond.ObjectID {
		case common.BKInnerObjIDHost:
			e.conds.hostCond = cond

		case common.BKInnerObjIDSet:
			e.conds.setCond = cond

		case common.BKInnerObjIDModule:
			e.conds.moduleCond = cond

		case common.BKInnerObjIDApp:
			e.conds.appCond = cond

		case common.BKInnerObjIDObject:
			e.conds.mainlineCond = cond

		case common.BKInnerObjIDPlat:
			e.conds.platCond = cond

		default:
			e.conds.objectCondMap[cond.ObjectID] = cond.Condition
		}
	}

	// parse and split conditions done, and clear orgin conditions.
	e.params.Condition = nil

	// add application id to app level.
	if e.params.AppID != -1 && e.params.AppID != 0 {
		condItem := metadata.ConditionItem{Field: common.BKAppIDField, Operator: common.BKDBEQ, Value: e.params.AppID}
		e.conds.appCond.Condition = append(e.conds.appCond.Condition, condItem)
	}
}

func (e *HostDynamicGroupExecutor) searchHostByConds() error {
	// search base on topology.
	err := e.searchByTopo()
	if err != nil {
		return err
	}
	if e.isNotFound {
		return nil
	}

	// search base on host conditions.
	err = e.searchByHostConds()
	if err != nil {
		return err
	}
	return nil
}

func (e *HostDynamicGroupExecutor) searchByTopo() error {
	// search base on application.
	err := e.searchByApp()
	if err != nil {
		return err
	}

	// search base on set.
	err = e.searchByMainline()
	if err != nil {
		return err
	}

	// search base on module.
	err = e.searchByModule()
	if err != nil {
		return err
	}

	// search base on plat.
	err = e.searchByPlatCondition()
	if err != nil {
		return err
	}

	return nil
}

func (e *HostDynamicGroupExecutor) searchByApp() error {
	if e.isNotFound {
		return nil
	}

	if len(e.conds.appCond.Condition) == 0 {
		return nil
	}

	appIDs, err := e.lgc.GetAppIDByCond(e.kit, e.conds.appCond.Condition)
	if err != nil {
		return err
	}

	if len(appIDs) == 0 {
		e.isNotFound = true
		return nil
	}

	e.conds.appCond.Condition = nil
	e.idArr.moduleHostConfig.appIDArr = appIDs

	return nil
}

func (e *HostDynamicGroupExecutor) searchByMainline() error {
	if e.isNotFound {
		return nil
	}

	var err error
	setIDs := []int64{}
	objSetIDs := []int64{}

	// search mainline object.
	if len(e.conds.mainlineCond.Condition) > 0 {
		objSetIDs, err = e.lgc.GetSetIDByObjectCond(e.kit, e.params.AppID, e.conds.mainlineCond.Condition)
		if err != nil {
			return err
		}

		if len(objSetIDs) == 0 {
			e.isNotFound = true
			return nil
		}
		e.conds.mainlineCond.Condition = nil

		e.conds.setCond.Condition = append(e.conds.setCond.Condition, metadata.ConditionItem{
			Field:    common.BKSetIDField,
			Operator: common.BKDBIN,
			Value:    objSetIDs,
		})
	}

	// search set.
	if len(e.conds.setCond.Condition) > 0 {
		if len(e.idArr.moduleHostConfig.appIDArr) > 0 {
			e.conds.setCond.Condition = append(e.conds.setCond.Condition, metadata.ConditionItem{
				Field:    common.BKAppIDField,
				Operator: common.BKDBIN,
				Value:    e.idArr.moduleHostConfig.appIDArr,
			})
		}

		setIDs, err = e.lgc.GetSetIDByCond(e.kit, e.conds.setCond.Condition)
		if err != nil {
			return err
		}

		if len(setIDs) == 0 {
			e.isNotFound = true
			return nil
		}
		e.conds.setCond.Condition = nil
		e.idArr.moduleHostConfig.setIDArr = setIDs
	}

	return nil
}

func (e *HostDynamicGroupExecutor) searchByModule() error {
	if e.isNotFound {
		return nil
	}

	if len(e.conds.moduleCond.Condition) == 0 {
		return nil
	}

	if len(e.idArr.moduleHostConfig.setIDArr) > 0 {
		e.conds.moduleCond.Condition = append(e.conds.moduleCond.Condition, metadata.ConditionItem{
			Field:    common.BKSetIDField,
			Operator: common.BKDBIN,
			Value:    e.idArr.moduleHostConfig.setIDArr,
		})
	}

	if len(e.idArr.moduleHostConfig.appIDArr) > 0 {
		e.conds.moduleCond.Condition = append(e.conds.moduleCond.Condition, metadata.ConditionItem{
			Field:    common.BKAppIDField,
			Operator: common.BKDBIN,
			Value:    e.idArr.moduleHostConfig.appIDArr,
		})
	}

	// search module.
	moduleIDs, err := e.lgc.GetModuleIDByCond(e.kit, e.conds.moduleCond.Condition)
	if err != nil {
		return err
	}

	if len(moduleIDs) == 0 {
		e.isNotFound = true
		return nil
	}

	e.conds.moduleCond.Condition = nil
	e.idArr.moduleHostConfig.moduleIDArr = moduleIDs

	return nil
}

func (e *HostDynamicGroupExecutor) searchByPlatCondition() error {
	if e.isNotFound {
		return nil
	}

	if len(e.conds.platCond.Condition) == 0 {
		return nil
	}

	instIDs, err := e.lgc.GetObjectInstByCond(e.kit, common.BKInnerObjIDPlat, e.conds.platCond.Condition)
	if err != nil {
		return err
	}

	if len(instIDs) == 0 {
		e.isNotFound = true
		return nil
	}

	e.conds.platCond.Condition = nil
	e.conds.hostCond.Condition = append(e.conds.hostCond.Condition, metadata.ConditionItem{
		Field:    common.BKCloudIDField,
		Operator: common.BKDBIN,
		Value:    instIDs,
	})

	return nil
}

func (e *HostDynamicGroupExecutor) searchByHostConds() error {
	if e.isNotFound {
		return nil
	}

	// add topology conditions.
	err := e.appendHostTopoConds()
	if err != nil {
		return err
	}
	if e.isNotFound {
		return nil
	}
	e.conds.hostCond.Fields = append(e.conds.hostCond.Fields, e.fields...)

	// empty means all fileds.
	if len(e.conds.hostCond.Fields) != 0 {
		// add more fields.
		e.conds.hostCond.Fields = append(e.conds.hostCond.Fields, common.BKHostIDField, common.BKCloudIDField)
	}

	condition := make(map[string]interface{})
	err = hostParse.ParseHostParams(e.conds.hostCond.Condition, condition)
	if err != nil {
		return err
	}
	err = hostParse.ParseHostIPParams(e.params.Ip, condition)
	if err != nil {
		return err
	}

	query := &metadata.QueryInput{
		Fields:         strings.Join(e.conds.hostCond.Fields, ","),
		Condition:      condition,
		Start:          e.params.Page.Start,
		Limit:          e.params.Page.Limit,
		Sort:           e.params.Page.Sort,
		DisableCounter: e.disableCounter,
	}

	e.conds.hostCond.Fields = nil
	e.params = nil

	if e.needPaged {
		query.Start = 0
	}

	result, err := e.lgc.CoreAPI.CoreService().Host().GetHosts(e.ctx, e.header, query)
	if err != nil {
		blog.Errorf("get hosts failed, err: %v, rid: %s", err, e.ccRid)
		return err
	}
	if !result.Result {
		blog.Errorf("get host failed, error code:%d, error message:%s, rid: %s", result.Code, result.ErrMsg, e.ccRid)
		return e.ccErr.New(result.Code, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		e.isNotFound = true
	}

	if !e.needPaged {
		e.total = result.Data.Count
	}

	for _, host := range result.Data.Info {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			return err
		}
		e.hosts = append(e.hosts, hostInfoStruct{hostID: hostID, hostInfo: host})
	}

	return nil
}

func (e *HostDynamicGroupExecutor) appendHostTopoConds() error {
	var moduleHostConfig metadata.DistinctHostIDByTopoRelationRequest
	isAddHostID := false

	if len(e.idArr.moduleHostConfig.setIDArr) > 0 {
		moduleHostConfig.SetIDArr = e.idArr.moduleHostConfig.setIDArr
		isAddHostID = true
	}
	if len(e.idArr.moduleHostConfig.moduleIDArr) > 0 {
		moduleHostConfig.ModuleIDArr = e.idArr.moduleHostConfig.moduleIDArr
		isAddHostID = true
	}
	if len(e.conds.objectCondMap) > 0 {
		moduleHostConfig.HostIDArr = e.idArr.moduleHostConfig.asstHostIDArr
		isAddHostID = true
	}
	if len(e.idArr.moduleHostConfig.appIDArr) > 0 {
		moduleHostConfig.ApplicationIDArr = e.idArr.moduleHostConfig.appIDArr
		isAddHostID = true
	}

	if !isAddHostID {
		// no module host config condition level.
		return nil
	}

	var hostIDs []int64

	respHostIDInfo, err := e.lgc.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(e.ctx, e.kit.Header, &moduleHostConfig)
	if err != nil {
		blog.Errorf("get hosts failed, err: %v, rid: %s", err, e.ccRid)
		return e.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if err := respHostIDInfo.CCError(); err != nil {
		blog.Errorf("get host id by topology relation failed, error code:%d, error message:%s, cond: %s, rid: %s",
			respHostIDInfo.Code, respHostIDInfo.ErrMsg, moduleHostConfig, e.ccRid)
		return err
	}
	e.total = len(respHostIDInfo.Data.IDArr)

	// 当有根据主机实例内容查询的时候的时候，无法在程序中完成分页
	hasHostCond := false
	if len(e.params.Ip.Data) > 0 || len(e.conds.hostCond.Condition) > 0 {
		hasHostCond = true
	}

	if !hasHostCond && e.params.Page.Limit > 0 {
		start := e.params.Page.Start
		limit := start + e.params.Page.Limit

		uniqHostIDCnt := len(respHostIDInfo.Data.IDArr)
		if start < 0 {
			start = 0
		}
		if start >= uniqHostIDCnt {
			e.isNotFound = true
			return nil
		}

		allHostIDs := respHostIDInfo.Data.IDArr
		sort.Slice(allHostIDs, func(i, j int) bool { return allHostIDs[i] < allHostIDs[j] })

		if uniqHostIDCnt <= limit {
			hostIDs = allHostIDs[start:]
		} else {
			hostIDs = allHostIDs[start:limit]
		}

		e.needPaged = true
	} else {
		if len(respHostIDInfo.Data.IDArr) == 0 {
			e.isNotFound = true
			return nil
		}
		hostIDs = respHostIDInfo.Data.IDArr
	}

	// 合并两种根据host_id查询的condition
	// 详情见issue: https://github.com/Tencent/bk-cmdb/issues/2461
	hostIDConditionExist := false
	for idx, cond := range e.conds.hostCond.Condition {
		if cond.Field != common.BKHostIDField {
			continue
		}

		// merge two condition
		// {"field": "bk_host_id", "operator": "$eq", "value": 1}
		// {"field": "bk_host_id", "operator": "$eq", "value": [1, 2]}
		// ==> {"field": "bk_host_id", "operator": "", "value": {"$in": [1,2], "$eq": 1}}
		hostIDConditionExist = true
		if cond.Operator != common.BKDBIN {
			// it's somewhat trick here to use common.BKDBEQ as merge operator
			cond = metadata.ConditionItem{
				Field:    common.BKHostIDField,
				Operator: common.BKDBEQ,
				Value:    map[string]interface{}{cond.Operator: cond.Value, common.BKDBIN: hostIDs},
			}
			e.conds.hostCond.Condition[idx] = cond

		} else {
			// intersection of two array
			value, ok := cond.Value.([]interface{})
			if !ok {
				blog.Errorf("invalid query condition with $in operator, value must be []int64, but got: %+v, rid: %s", cond.Value, e.ccRid)
				return e.ccErr.New(common.CCErrCommParamsIsInvalid, common.BKHostIDField)
			}

			hostIDMap := make(map[int64]bool)
			for _, hostID := range hostIDs {
				hostIDMap[hostID] = true
			}

			shareIDs := make([]int64, 0)
			for _, hostID := range value {
				id, err := util.GetInt64ByInterface(hostID)
				if err != nil {
					blog.Errorf("invalid query condition with $in operator, value must be []int64, but got: %+v, rid: %s", cond.Value, e.ccRid)
					return e.ccErr.New(common.CCErrCommParamsIsInvalid, common.BKHostIDField)
				}

				if hostIDMap[id] {
					shareIDs = append(shareIDs, id)
				}
			}
			e.conds.hostCond.Condition[idx].Value = shareIDs
		}
	}

	if !hostIDConditionExist {
		e.conds.hostCond.Condition = append(e.conds.hostCond.Condition, metadata.ConditionItem{
			Field:    common.BKHostIDField,
			Operator: common.BKDBIN,
			Value:    hostIDs,
		})
	}

	return nil
}

func (e *HostDynamicGroupExecutor) buildSearchResult() ([]mapstr.MapStr, int) {
	result := make([]mapstr.MapStr, 0)

	if e.isNotFound {
		return result, 0
	}

	for _, host := range e.hosts {
		result = append(result, host.hostInfo)
	}

	// return search result and total num of row matched with the conditions.
	return result, e.total
}
