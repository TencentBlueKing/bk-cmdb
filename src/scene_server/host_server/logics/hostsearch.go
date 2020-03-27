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
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	hostParse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
)

func (lgc *Logics) SearchHost(ctx context.Context, data *metadata.HostCommonSearch, isDetail bool) (*metadata.SearchHost, error) {
	searchHostInst := NewSearchHost(ctx, lgc, data)
	searchHostInst.ParseCondition()
	retHostInfo := &metadata.SearchHost{
		Info: make([]mapstr.MapStr, 0),
	}
	err := searchHostInst.SearchHostByConds()
	if err != nil {
		return retHostInfo, err
	}
	hostInfoArr, cnt, err := searchHostInst.FillTopologyData()
	if err != nil {
		return retHostInfo, err
	}

	retHostInfo.Count = cnt
	if cnt > 0 {
		retHostInfo.Info = hostInfoArr
	}
	return retHostInfo, nil
}

type searchHostConds struct {
	hostCond      metadata.SearchCondition
	appCond       metadata.SearchCondition
	setCond       metadata.SearchCondition
	moduleCond    metadata.SearchCondition
	mainlineCond  metadata.SearchCondition
	platCond      metadata.SearchCondition
	objectCondMap map[string][]metadata.ConditionItem
}

type searchHostTopologyShowSection struct {
	app    bool
	set    bool
	module bool
}

type searchHostModuleHostConfig struct {
	appIDArr      []int64
	moduleIDArr   []int64
	setIDArr      []int64
	asstHostIDArr []int64
}

type searchHostIDArr struct {
	moduleHostConfig searchHostModuleHostConfig
}

type searchHostInfoMapCache struct {
	appInfoMap           map[int64]mapstr.MapStr
	setInfoMap           map[int64]mapstr.MapStr
	moduleInfoMap        map[int64]mapstr.MapStr
	cloudAsstNameInfoMap map[int64]*InstNameAsst
}

type hostInfoStruct struct {
	hostID   int64
	hostInfo mapstr.MapStr
}

type searchHost struct {
	lgc             *Logics
	pheader         http.Header
	hostSearchParam *metadata.HostCommonSearch
	//  this part need to be displayed?
	topoShowSection searchHostTopologyShowSection

	conds searchHostConds
	//search end, condition not dsetAppConfigata
	noData       bool
	idArr        searchHostIDArr
	hostInfoArr  []hostInfoStruct // int64 is hostID
	cacheInfoMap searchHostInfoMapCache
	totalHostCnt int

	ccErr errors.DefaultCCErrorIf
	ccRid string
	ctx   context.Context
}

type setLevelInfo struct {
	setName       string
	setID         int64
	moduleInfoMap map[int64]int64
}

type appLevelInfo struct {
	appName    string
	appID      int64
	setInfoMap map[int64]*setLevelInfo
}

// searchHostInterface Too many methods, hiding private methods
type searchHostInterface interface {
	ParseCondition()
	SearchHostByConds() errors.CCError
	FillTopologyData() ([]mapstr.MapStr, int, errors.CCError)
}

func NewSearchHost(ctx context.Context, lgc *Logics, hostSearchParam *metadata.HostCommonSearch) searchHostInterface {
	sh := &searchHost{
		lgc:             lgc,
		pheader:         lgc.header,
		hostSearchParam: hostSearchParam,
		idArr:           searchHostIDArr{},
		ccRid:           lgc.rid,
		ccErr:           lgc.ccErr,
		ctx:             ctx,
	}

	sh.conds.objectCondMap = make(map[string][]metadata.ConditionItem)

	return sh
}

func (sh *searchHost) ParseCondition() {

	for _, object := range sh.hostSearchParam.Condition {
		if object.ObjectID == common.BKInnerObjIDHost {
			sh.conds.hostCond = object
		} else if object.ObjectID == common.BKInnerObjIDSet {
			sh.conds.setCond = object
			sh.topoShowSection.set = true
		} else if object.ObjectID == common.BKInnerObjIDModule {
			sh.conds.moduleCond = object
			sh.topoShowSection.module = true
		} else if object.ObjectID == common.BKInnerObjIDApp {
			sh.conds.appCond = object
			sh.topoShowSection.app = true
		} else if object.ObjectID == common.BKInnerObjIDObject {
			sh.conds.mainlineCond = object
		} else if object.ObjectID == common.BKInnerObjIDPlat {
			sh.conds.platCond = object
		} else {
			sh.conds.objectCondMap[object.ObjectID] = object.Condition
		}
	}

	sh.tryParseAppID()

}

func (sh *searchHost) SearchHostByConds() errors.CCError {

	err := sh.searchByTopo()
	if err != nil {
		return err
	}
	if sh.noData {
		return nil
	}
	err = sh.searchByHostConds()
	if err != nil {
		return err
	}
	if sh.noData {
		return nil
	}

	return nil

}

func (sh *searchHost) FillTopologyData() ([]mapstr.MapStr, int, errors.CCError) {

	if sh.noData {
		return nil, 0, nil
	}

	hostIDArr := make([]int64, 0)
	for _, hostInfoItem := range sh.hostInfoArr {
		hostIDArr = append(hostIDArr, hostInfoItem.hostID)
	}
	queryCond := metadata.HostModuleRelationRequest{
		HostIDArr: hostIDArr,
	}
	mhconfig, err := sh.lgc.GetConfigByCond(sh.ctx, queryCond)
	if err != nil {
		return nil, 0, err
	}
	appIDArr := make([]int64, 0)
	setIDArr := make([]int64, 0)
	moduleIDArr := make([]int64, 0)

	hostAppSetModuleConfig := make(map[int64]map[int64]*appLevelInfo, 0)

	blog.V(5).Infof("get modulehostconfig map:%v, rid:%s", mhconfig, sh.ccRid)
	for _, mh := range mhconfig {
		hostAppInfoLevelInst, ok := hostAppSetModuleConfig[mh.HostID]
		if !ok {
			hostAppInfoLevelInst = make(map[int64]*appLevelInfo, 0)
			hostAppSetModuleConfig[mh.HostID] = hostAppInfoLevelInst
		}

		appInfoLevelInst, ok := hostAppInfoLevelInst[mh.AppID]
		if !ok {
			appInfoLevelInst = &appLevelInfo{
				setInfoMap: make(map[int64]*setLevelInfo, 0),
			}
			hostAppInfoLevelInst[mh.AppID] = appInfoLevelInst
		}
		setInfoLevleInst, ok := appInfoLevelInst.setInfoMap[mh.SetID]
		if !ok {
			setInfoLevleInst = &setLevelInfo{
				moduleInfoMap: make(map[int64]int64, 0),
			}
			appInfoLevelInst.setInfoMap[mh.SetID] = setInfoLevleInst
		}
		setInfoLevleInst.moduleInfoMap[mh.ModuleID] = mh.ModuleID

		appIDArr = append(appIDArr, mh.AppID)
		setIDArr = append(setIDArr, mh.SetID)
		moduleIDArr = append(moduleIDArr, mh.ModuleID)
	}
	appInfoMap, err := sh.fetchTopoAppCacheInfo(appIDArr)
	if err != nil {
		return nil, 0, err
	}
	sh.cacheInfoMap.appInfoMap = appInfoMap

	setInfoMap, err := sh.fetchTopoSetCacheInfo(setIDArr)
	if err != nil {
		return nil, 0, err
	}
	sh.cacheInfoMap.setInfoMap = setInfoMap

	moduleInfoMap, err := sh.fetchTopoModuleCacheInfo(moduleIDArr)
	if err != nil {
		return nil, 0, err
	}
	sh.cacheInfoMap.moduleInfoMap = moduleInfoMap

	cloudAsstNameInfoMap, err := sh.fetchHostCloudCacheInfo()
	if err != nil {
		return nil, 0, err
	}
	sh.cacheInfoMap.cloudAsstNameInfoMap = cloudAsstNameInfoMap

	result := make([]mapstr.MapStr, 0)
	for _, hostInfoItem := range sh.hostInfoArr {
		searchHostItem := mapstr.New()
		levelInfo, ok := hostAppSetModuleConfig[hostInfoItem.hostID]
		if !ok {
			continue
		}
		searchHostItem = sh.fillHostAppInfo(levelInfo, searchHostItem)
		searchHostItem = sh.fillHostSetInfo(levelInfo, searchHostItem)
		searchHostItem = sh.fillHostModuleInfo(levelInfo, searchHostItem)
		hostInfo := sh.fillHostCloudInfo(hostInfoItem.hostInfo, searchHostItem)
		searchHostItem.Set(common.BKInnerObjIDHost, hostInfo)

		result = append(result, searchHostItem)
	}
	return result, sh.totalHostCnt, nil

}

/* ** fill host cloud info  ** */

func (sh *searchHost) fillHostCloudInfo(hostInfo, searchHostItem mapstr.MapStr) mapstr.MapStr {
	clouldID, err := hostInfo.Int64(common.BKCloudIDField)
	if err != nil {
		blog.Warnf("search host fillHostCloudInfo host get cloud id error, hostinfo:%+v, error:%s, rid:%s", hostInfo, err.Error(), sh.ccRid)
		hostInfo.Set(common.BKCloudIDField, make([]InstNameAsst, 0))
		return hostInfo
	}
	instAsst, ok := sh.cacheInfoMap.cloudAsstNameInfoMap[clouldID]
	if !ok {
		blog.Warnf("search host fillHostCloudInfo host  cloud id not found, cloud id:%d, hostinfo:%+v, rid:%s", clouldID, hostInfo, sh.ccRid)
		hostInfo.Set(common.BKCloudIDField, make([]InstNameAsst, 0))
		return hostInfo
	}
	hostInfo.Set(common.BKCloudIDField, []*InstNameAsst{instAsst})
	return hostInfo
}

func (sh *searchHost) fetchHostCloudCacheInfo() (map[int64]*InstNameAsst, errors.CCError) {
	cloudIDMap := make(map[int64]bool, 0)
	for _, hostInfoItem := range sh.hostInfoArr {

		cloudID, err := hostInfoItem.hostInfo.Int64(common.BKCloudIDField)
		if err != nil {
			blog.Warnf("hostSearch not found  cloud id in hsot, hostInfo:%d, rid:%s", hostInfoItem.hostInfo, sh.ccRid)
			continue
		}
		cloudIDMap[cloudID] = true
	}
	var cloudIDArr []int64
	for cloudID := range cloudIDMap {
		cloudIDArr = append(cloudIDArr, cloudID)
	}
	queryInput := &metadata.QueryCondition{}
	queryInput.Condition = mapstr.MapStr{
		common.BKCloudIDField: mapstr.MapStr{
			common.BKDBIN: cloudIDArr,
		},
	}
	result, err := sh.lgc.CoreAPI.CoreService().Instance().ReadInstance(sh.ctx, sh.pheader, common.BKInnerObjIDPlat, queryInput)
	if err != nil {
		blog.Errorf("fetchHostCloudCacheInfo SearchObjects http do error, err:%s,input:%+v,rid:%s", err.Error(), queryInput, sh.ccRid)
		return nil, sh.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("fetchHostCloudCacheInfo SearchObjects http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, queryInput, sh.ccRid)
		return nil, sh.ccErr.New(result.Code, result.ErrMsg)
	}

	cloudInfoMap := make(map[int64]*InstNameAsst)
	for _, info := range result.Data.Info {
		asstInst, err := sh.convInstInfoToAssociateInfo(common.BKCloudIDField, common.BKCloudNameField, common.BKInnerObjIDPlat, info)
		if err != nil {
			return nil, err
		}
		cloudInfoMap[asstInst.ObjectID] = asstInst
	}
	return cloudInfoMap, nil

}

func (sh *searchHost) convInstInfoToAssociateInfo(instIDKey, instNameKey, objID string, instInfo mapstr.MapStr) (*InstNameAsst, errors.CCError) {
	if val, exist := instInfo[instNameKey]; exist {
		asstInst := &InstNameAsst{}
		if name, can := val.(string); can {
			asstInst.Name = name
			asstInst.ObjID = objID
		}
		instID, err := instInfo.Int64(instIDKey)
		if err != nil {
			return nil, sh.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, objID, instIDKey, "int", err.Error())
		}
		asstInst.ID = strconv.FormatInt(instID, 10)
		asstInst.ObjectID = instID
		return asstInst, nil
	}

	return nil, nil
}

/* ** fill host topology data ** */

func (sh *searchHost) fillHostAppInfo(appInfoLevelInst map[int64]*appLevelInfo, searchHostItem mapstr.MapStr) mapstr.MapStr {

	appInfoArr := make([]mapstr.MapStr, 0)
	var err error
	//appdata
	for appID, appLevelInfo := range appInfoLevelInst {

		appInfo, mapOk := sh.cacheInfoMap.appInfoMap[appID]
		if mapOk {
			appInfoArr = append(appInfoArr, appInfo)
		}
		appLevelInfo.appID = appID
		appLevelInfo.appName, err = appInfo.String(common.BKAppNameField)
		if err != nil {
			blog.Warnf("hostSearch not found app name, appInfo:%d, rid:%s", appInfo, sh.ccRid)
			continue
		}

	}
	searchHostItem.Set(common.BKInnerObjIDApp, appInfoArr)
	return searchHostItem

}

func (sh *searchHost) fillHostSetInfo(appInfoLevelInst map[int64]*appLevelInfo, searchHostItem mapstr.MapStr) mapstr.MapStr {

	setInfoArr := make([]mapstr.MapStr, 0)
	for _, appLevelInfo := range appInfoLevelInst {
		for setID, setLevelInfo := range appLevelInfo.setInfoMap {
			setInfo, isOk := sh.cacheInfoMap.setInfoMap[setID]
			if false == isOk {
				continue
			}

			setName, err := setInfo.String(common.BKSetNameField)
			if nil != err {
				blog.Warnf("hostSearch not found set name, setInfo:%d, rid:%s", setInfo, sh.ccRid)
				continue
			}
			setLevelInfo.setID = setID
			setLevelInfo.setName = setName
			setInfo.Set(TopoModuleName, appLevelInfo.appName+SplitFlag+setName)

			setInfoArr = append(setInfoArr, setInfo)
		}

	}

	searchHostItem[common.BKInnerObjIDSet] = setInfoArr
	return searchHostItem
}

func (sh *searchHost) fillHostModuleInfo(appInfoLevelInst map[int64]*appLevelInfo, searchHostItem mapstr.MapStr) mapstr.MapStr {

	moduleInfoArr := make([]mapstr.MapStr, 0)
	for _, appLevelInfo := range appInfoLevelInst {
		for _, setLevelInfo := range appLevelInfo.setInfoMap {
			for mdouleID := range setLevelInfo.moduleInfoMap {
				moduleInfo, ok := sh.cacheInfoMap.moduleInfoMap[mdouleID]
				if false == ok {
					blog.V(5).Infof("hostSearch not found module id, moduleID:%d, hostModuleMap:%v, rid:%s", mdouleID, sh.cacheInfoMap.moduleInfoMap, sh.ccRid)
					continue
				}

				moduleName, err := moduleInfo.String(common.BKModuleNameField)
				if nil != err {
					blog.V(5).Infof("hostSearch not found module name, moduleID:%d, hostModuleMap:%v, rid:%s", mdouleID, sh.cacheInfoMap.moduleInfoMap, sh.ccRid)
					continue
				}
				datacp := make(map[string]interface{})
				for key, val := range moduleInfo {
					datacp[key] = val
				}
				datacp[TopoModuleName] = appLevelInfo.appName + SplitFlag + setLevelInfo.setName + SplitFlag + moduleName
				moduleInfoArr = append(moduleInfoArr, datacp)
			}
		}

	}

	searchHostItem[common.BKInnerObjIDModule] = moduleInfoArr
	return searchHostItem
}

func (sh *searchHost) fetchTopoAppCacheInfo(appIDArr []int64) (map[int64]mapstr.MapStr, errors.CCError) {

	if nil != sh.conds.appCond.Fields {
		if len(sh.conds.appCond.Fields) != 0 {
			sh.conds.appCond.Fields = append(sh.conds.appCond.Fields, common.BKAppIDField)
			sh.conds.appCond.Fields = append(sh.conds.appCond.Fields, common.BKAppNameField)
		}
		cond := mapstr.New()
		celld := mapstr.New()
		celld.Set(common.BKDBIN, appIDArr)
		cond.Set(common.BKAppIDField, celld)
		return sh.lgc.GetAppMapByCond(sh.ctx, sh.conds.appCond.Fields, cond)

	}
	return nil, nil
}

func (sh *searchHost) fetchTopoSetCacheInfo(setIDArr []int64) (map[int64]mapstr.MapStr, errors.CCError) {

	if nil != sh.conds.setCond.Fields {
		exist := util.InArray(common.BKSetIDField, sh.conds.setCond.Fields)
		if !exist && 0 != len(sh.conds.setCond.Fields) {
			sh.conds.setCond.Fields = append(sh.conds.setCond.Fields, common.BKSetIDField)
		}
		cond := mapstr.New()
		celld := mapstr.New()
		celld.Set(common.BKDBIN, setIDArr)
		cond.Set(common.BKSetIDField, celld)
		return sh.lgc.GetSetMapByCond(sh.ctx, sh.conds.setCond.Fields, cond)
	}

	return nil, nil
}

func (sh *searchHost) fetchTopoModuleCacheInfo(moduleIDArr []int64) (map[int64]mapstr.MapStr, errors.CCError) {
	if nil != sh.conds.moduleCond.Fields {
		exist := util.InArray(common.BKModuleIDField, sh.conds.moduleCond.Fields)
		if !exist && 0 != len(sh.conds.moduleCond.Fields) {
			sh.conds.moduleCond.Fields = append(sh.conds.moduleCond.Fields, common.BKModuleIDField)
		}
		cond := mapstr.New()
		celld := mapstr.New()
		celld.Set(common.BKDBIN, moduleIDArr)
		cond.Set(common.BKModuleIDField, celld)
		return sh.lgc.GetModuleMapByCond(sh.ctx, sh.conds.moduleCond.Fields, cond)

	}

	return nil, nil

}

/* ** The following is the processing of querying data according to conditions. ** */

func (sh *searchHost) searchByTopo() errors.CCError {
	err := sh.searchByApp()
	if err != nil {
		return err
	}
	err = sh.searchByMainline()
	if err != nil {
		return err
	}
	err = sh.searchByModule()
	if err != nil {
		return err
	}
	//Query host information based on associated objects, alternate code
	//sh.searchByAssocation()
	err = sh.searchByPlatCondition()
	if err != nil {
		return err
	}
	if sh.noData {
		return nil
	}
	return nil
}

func (sh *searchHost) searchByPlatCondition() errors.CCError {
	if sh.noData {
		return nil
	}
	if len(sh.conds.platCond.Condition) > 0 {
		instIDArr, err := sh.lgc.GetObjectInstByCond(sh.ctx, common.BKInnerObjIDPlat, sh.conds.platCond.Condition)
		if err != nil {
			return err
		}
		if len(instIDArr) == 0 {
			sh.noData = true
			return nil
		}
		sh.conds.hostCond.Condition = append(sh.conds.hostCond.Condition, metadata.ConditionItem{
			Field:    common.BKCloudIDField,
			Operator: common.BKDBIN,
			Value:    instIDArr,
		})

	}

	return nil
}

func (sh *searchHost) searchByApp() errors.CCError {
	if sh.noData {
		return nil
	}
	if len(sh.conds.appCond.Condition) > 0 {
		appIDArr, err := sh.lgc.GetAppIDByCond(sh.ctx, sh.conds.appCond.Condition)
		if err != nil {
			return err
		}
		if len(appIDArr) == 0 {
			sh.noData = true
			return nil
		}
		sh.idArr.moduleHostConfig.appIDArr = appIDArr
	}
	return nil
}

func (sh *searchHost) searchByMainline() errors.CCError {

	if sh.noData {
		return nil
	}

	var err error
	setIDArr := make([]int64, 0)
	objSetIDArr := make([]int64, 0)

	// search mainline object by cond
	if len(sh.conds.mainlineCond.Condition) > 0 {
		objSetIDArr, err = sh.lgc.GetSetIDByObjectCond(sh.ctx, sh.hostSearchParam.AppID, sh.conds.mainlineCond.Condition)
		if err != nil {
			return err
		}
		if len(objSetIDArr) == 0 {
			sh.noData = true
			return nil
		}
	}
	// search set by appcond
	if len(sh.conds.setCond.Condition) > 0 || len(sh.conds.mainlineCond.Condition) > 0 {
		if len(sh.conds.appCond.Condition) > 0 {
			sh.conds.setCond.Condition = append(sh.conds.setCond.Condition, metadata.ConditionItem{
				Field:    common.BKAppIDField,
				Operator: common.BKDBIN,
				Value:    sh.idArr.moduleHostConfig.appIDArr,
			})
		}
		if len(sh.conds.mainlineCond.Condition) > 0 {
			sh.conds.setCond.Condition = append(sh.conds.setCond.Condition, metadata.ConditionItem{
				Field:    common.BKSetIDField,
				Operator: common.BKDBIN,
				Value:    objSetIDArr,
			})
		}
		setIDArr, err = sh.lgc.GetSetIDByCond(sh.ctx, sh.conds.setCond.Condition)
		if err != nil {
			return err
		}
		if len(setIDArr) == 0 {
			sh.noData = true
			return nil
		}
	}

	if len(sh.conds.setCond.Condition) > 0 {
		sh.idArr.moduleHostConfig.setIDArr = setIDArr
	}

	return nil
}

func (sh *searchHost) searchByModule() errors.CCError {
	if sh.noData {
		return nil
	}
	if len(sh.conds.moduleCond.Condition) > 0 {
		if len(sh.conds.setCond.Condition) > 0 {
			sh.conds.moduleCond.Condition = append(sh.conds.moduleCond.Condition, metadata.ConditionItem{
				Field:    common.BKSetIDField,
				Operator: common.BKDBIN,
				Value:    sh.idArr.moduleHostConfig.setIDArr,
			})
		}
		if len(sh.conds.appCond.Condition) > 0 {
			sh.conds.moduleCond.Condition = append(sh.conds.moduleCond.Condition, metadata.ConditionItem{
				Field:    common.BKAppIDField,
				Operator: common.BKDBIN,
				Value:    sh.idArr.moduleHostConfig.appIDArr,
			})
		}
		//search module by cond
		moduleIDArr, err := sh.lgc.GetModuleIDByCond(sh.ctx, sh.conds.moduleCond.Condition)
		if err != nil {
			return err
		}
		if len(moduleIDArr) == 0 {
			sh.noData = true
			return nil
		}
		if len(sh.conds.moduleCond.Condition) > 0 {
			sh.idArr.moduleHostConfig.moduleIDArr = moduleIDArr
		}
	}

	return nil
}

func (sh *searchHost) searchByHostConds() errors.CCError {
	if sh.noData {
		return nil
	}

	err := sh.appendHostTopoConds()
	if err != nil {
		return err
	}

	if 0 != len(sh.conds.hostCond.Fields) {
		sh.conds.hostCond.Fields = append(sh.conds.hostCond.Fields, common.BKHostIDField, common.BKCloudIDField)
	}

	condition := make(map[string]interface{})
	err = hostParse.ParseHostParams(sh.conds.hostCond.Condition, condition)
	if err != nil {
		return err
	}

	err = hostParse.ParseHostIPParams(sh.hostSearchParam.Ip, condition)
	if err != nil {
		return err
	}

	query := &metadata.QueryInput{
		Condition: condition,
		Start:     sh.hostSearchParam.Page.Start,
		Limit:     sh.hostSearchParam.Page.Limit,
		Sort:      sh.hostSearchParam.Page.Sort,
		Fields:    strings.Join(sh.conds.hostCond.Fields, ","),
	}

	gResult, err := sh.lgc.CoreAPI.CoreService().Host().GetHosts(sh.ctx, sh.pheader, query)
	if err != nil {
		blog.Errorf("get hosts failed, err: %v, rid: %s", err, sh.ccRid)
		return err
	}
	if !gResult.Result {
		blog.Errorf("get host failed, error code:%d, error message:%s, rid: %s", gResult.Code, gResult.ErrMsg, sh.ccRid)
		return sh.ccErr.New(gResult.Code, gResult.ErrMsg)
	}

	if len(gResult.Data.Info) == 0 {
		sh.noData = true
	}

	sh.totalHostCnt = gResult.Data.Count
	for _, host := range gResult.Data.Info {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			return err
		}
		sh.hostInfoArr = append(sh.hostInfoArr, hostInfoStruct{
			hostID:   hostID,
			hostInfo: host,
		})
	}
	return nil
}

func (sh *searchHost) appendHostTopoConds() errors.CCError {
	var moduleHostConfig metadata.HostModuleRelationRequest
	isAddHostID := false

	if len(sh.conds.setCond.Condition) > 0 {
		moduleHostConfig.SetIDArr = sh.idArr.moduleHostConfig.setIDArr
		isAddHostID = true
	}
	if len(sh.conds.moduleCond.Condition) > 0 {
		moduleHostConfig.ModuleIDArr = sh.idArr.moduleHostConfig.moduleIDArr
		isAddHostID = true
	}
	if len(sh.conds.objectCondMap) > 0 {
		moduleHostConfig.HostIDArr = sh.idArr.moduleHostConfig.asstHostIDArr
		isAddHostID = true
	}

	var appIDArr []int64
	if len(sh.conds.appCond.Condition) > 0 {
		appIDArr = sh.idArr.moduleHostConfig.appIDArr
		isAddHostID = true
	}
	if !isAddHostID {
		return nil
	}
	var moduleHostConfigArr []metadata.HostModuleRelationRequest
	if len(appIDArr) > 0 {
		//
		for _, appID := range appIDArr {
			newModuleHostConfig := *(&moduleHostConfig)
			newModuleHostConfig.ApplicationID = appID
			moduleHostConfigArr = append(moduleHostConfigArr, newModuleHostConfig)
		}
	} else {
		moduleHostConfigArr = append(moduleHostConfigArr, moduleHostConfig)
	}
	hostIDArr := make([]int64, 0)
	for _, moduleHostConfig := range moduleHostConfigArr {
		hostIDArrItem, err := sh.lgc.GetHostIDByCond(sh.ctx, moduleHostConfig)
		if err != nil {
			blog.Errorf("GetHostIDByCond get hosts failed, err: %v, rid: %s", err, sh.ccRid)
			return err
		}
		hostIDArr = append(hostIDArr, hostIDArrItem...)
	}

	// 合并两种涞源的根据 host_id 查询的 condition
	// 详情见issue: https://github.com/Tencent/bk-cmdb/issues/2461
	hostIDConditionExist := false
	for idx, cond := range sh.conds.hostCond.Condition {
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
				Value: map[string]interface{}{
					cond.Operator: cond.Value,
					common.BKDBIN: hostIDArr,
				},
			}
			sh.conds.hostCond.Condition[idx] = cond
		} else {
			// intersection of two array
			value, ok := cond.Value.([]interface{})
			if ok == false {
				blog.Errorf("invalid query condition with $in operator, value must be []int64, but got: %+v, rid: %s", cond.Value, sh.ccRid)
				return sh.ccErr.New(common.CCErrCommParamsIsInvalid, common.BKHostIDField)
			}
			shareIDs := make([]int64, 0)
			for _, hostID := range value {
				id, err := util.GetInt64ByInterface(hostID)
				if err != nil {
					blog.Errorf("invalid query condition with $in operator, value must be []int64, but got: %+v, rid: %s", cond.Value, sh.ccRid)
					return sh.ccErr.New(common.CCErrCommParamsIsInvalid, common.BKHostIDField)
				}
				if in := util.InArray(id, hostIDArr); in == true {
					shareIDs = append(shareIDs, id)
				}
			}
			sh.conds.hostCond.Condition[idx].Value = shareIDs
		}
	}
	if hostIDConditionExist == false {
		sh.conds.hostCond.Condition = append(sh.conds.hostCond.Condition, metadata.ConditionItem{
			Field:    common.BKHostIDField,
			Operator: common.BKDBIN,
			Value:    hostIDArr,
		})
	}

	return nil
}

// Query host information based on associated objects, alternate code
func (sh *searchHost) searchByAssociation() errors.CCError {
	instAsstHostIDArr := make([]int64, 0)
	//search host id by object
	firstCond := true
	if len(sh.conds.objectCondMap) > 0 {
		for objID, objCond := range sh.conds.objectCondMap {
			instIDArr, err := sh.lgc.GetObjectInstByCond(sh.ctx, objID, objCond)
			if err != nil {
				return err
			}
			instHostIDArr, err := sh.lgc.GetHostIDByInstID(sh.ctx, objID, instIDArr)
			if err != nil {
				return err
			}
			if firstCond {
				instAsstHostIDArr = instHostIDArr
			} else {
				instAsstHostIDArr = util.IntArrIntersection(instAsstHostIDArr, instHostIDArr)
			}
			firstCond = false
		}

	}
	instAsstHostIDArr = util.IntArrayUnique(instAsstHostIDArr)
	if len(sh.conds.objectCondMap) > 0 {
		sh.idArr.moduleHostConfig.asstHostIDArr = instAsstHostIDArr
	}

	return nil

}

func (sh *searchHost) tryParseAppID() {
	//search appID by cond
	if -1 != sh.hostSearchParam.AppID && 0 != sh.hostSearchParam.AppID {
		sh.conds.appCond.Condition = append(sh.conds.appCond.Condition, metadata.ConditionItem{
			Field:    common.BKAppIDField,
			Operator: common.BKDBEQ,
			Value:    sh.hostSearchParam.AppID,
		})
	}
}
