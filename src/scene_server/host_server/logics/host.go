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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	hostParse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
	"configcenter/src/scene_server/validator"
)

func (lgc *Logics) GetHostAttributes(ownerID string, header http.Header) ([]metadata.Header, error) {
	searchOp := hutil.NewOperation().WithObjID(common.BKInnerObjIDHost).WithOwnerID(ownerID).Data()
	result, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), header, searchOp)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("search host obj log failed, err: %v, result err: %s", err, result.ErrMsg)
	}

	headers := make([]metadata.Header, 0)
	for _, p := range result.Data {
		if p.PropertyID == common.BKChildStr {
			continue
		}
		headers = append(headers, metadata.Header{
			PropertyID:   p.PropertyID,
			PropertyName: p.PropertyName,
		})
	}

	return headers, nil
}

func (lgc *Logics) GetHostInstanceDetails(pheader http.Header, ownerID, hostID string) (map[string]interface{}, string, error) {
	// get host details, pre data
	result, err := lgc.CoreAPI.HostController().Host().GetHostByID(context.Background(), hostID, pheader)
	if err != nil || (err == nil && !result.Result) {
		return nil, "", fmt.Errorf("get host  data failed, err, %v, %v", err, result.ErrMsg)
	}

	hostInfo := result.Data
	attributes, err := lgc.GetObjectAsst(ownerID, pheader)
	if err != nil {
		return nil, "", err
	}

	for key, val := range attributes {
		if item, ok := hostInfo[key]; ok {
			if item == nil {
				continue
			}

			strItem := util.GetStrByInterface(item)
			ids := make([]int64, 0)
			for _, strID := range strings.Split(strItem, ",") {
				id, err := strconv.ParseInt(strID, 10, 64)
				if err != nil {
					return nil, "", err
				}
				ids = append(ids, id)
			}

			//cond := make(map[string]interface{})
			//cond[common.BKHostIDField] = map[string]interface{}{"$in": ids}
			q := &metadata.QueryInput{
				Condition: nil, //cond,
				Fields:    "",
				Start:     0,
				Limit:     common.BKNoLimit,
				Sort:      "",
			}

			asst, _, err := lgc.getInstAsst(ownerID, val, strings.Split(strItem, ","), pheader, q)
			if err != nil {
				return nil, "", fmt.Errorf("get instance asst failed, err: %v", err)
			}
			hostInfo[key] = asst
		}
	}

	ip := hostInfo[common.BKHostInnerIPField].(string)
	return hostInfo, ip, nil
}

func (lgc *Logics) GetConfigByCond(pheader http.Header, cond map[string][]int64) ([]map[string]int64, error) {
	configArr := make([]map[string]int64, 0)

	if 0 == len(cond) {
		return configArr, nil
	}

	result, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(context.Background(), pheader, cond)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get module host config failed, err: %v, %v", err, result.ErrMsg)
	}

	for _, info := range result.Data {
		data := make(map[string]int64)
		data[common.BKAppIDField] = info.AppID
		data[common.BKSetIDField] = info.SetID
		data[common.BKModuleIDField] = info.ModuleID
		data[common.BKHostIDField] = info.HostID
		configArr = append(configArr, data)
	}
	return configArr, nil
}

// EnterIP 将机器导入到制定模块或者空闲机器， 已经存在机器，不操作
func (lgc *Logics) EnterIP(pheader http.Header, ownerID string, appID, moduleID int64, ip string, cloudID int64, host map[string]interface{}, isIncrement bool) error {

	user := util.GetUser(pheader)
	lang := lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(pheader))

	isExist, err := lgc.IsPlatExist(pheader, common.KvMap{common.BKCloudIDField: cloudID})
	if nil != err {
		return errors.New(lang.Languagef("plat_get_str_err", err.Error())) // "查询主机信息失败")
	}
	if !isExist {
		return errors.New(lang.Language("plat_id_not_exist"))
	}
	conds := map[string]interface{}{
		common.BKHostInnerIPField: ip,
		common.BKCloudIDField:     cloudID,
	}
	hostList, err := lgc.GetHostInfoByConds(pheader, conds)
	if nil != err {
		return errors.New(lang.Languagef("host_search_fail", err.Error())) // "查询主机信息失败")
	}

	hostID := int64(0)
	if len(hostList) == 0 {
		//host not exist, add host
		host[common.BKHostInnerIPField] = ip
		host[common.BKCloudIDField] = cloudID
		host["import_from"] = common.HostAddMethodAgent
		defaultFields, hasErr := lgc.getHostFields(ownerID, pheader)
		if nil != hasErr {
			blog.Errorf("get host property error; error:%s", hasErr.Error())
			return errors.New("get host property error")
		}
		//补充未填写字段的默认值
		for _, field := range defaultFields {
			_, ok := host[field.PropertyID]
			if !ok {
				if true == util.IsStrProperty(field.PropertyType) {
					host[field.PropertyID] = ""
				} else {
					host[field.PropertyID] = nil
				}
			}
		}
		valid := validator.NewValidMap(util.GetOwnerID(pheader), common.BKInnerObjIDHost, pheader, lgc.Engine)
		hasErr = valid.ValidMap(host, "create", 0)

		if nil != hasErr {
			return hasErr
		}

		result, err := lgc.CoreAPI.HostController().Host().AddHost(context.Background(), pheader, host)
		if err != nil {
			return errors.New(lang.Languagef("host_agent_add_host_fail", err.Error()))
		} else if err == nil && !result.Result {
			return errors.New(lang.Languagef("host_agent_add_host_fail", result.ErrMsg))
		}

		retHost := result.Data.(map[string]interface{})
		hostID, err = util.GetInt64ByInterface(retHost[common.BKHostIDField])
		if err != nil {
			return errors.New(lang.Languagef("host_agent_add_host_fail", err.Error()))
		}

	} else if false == isIncrement {
		//Not an additional relationship model
		return nil
	} else {

		hostID, err = util.GetInt64ByInterface(hostList[0][common.BKHostIDField])
		if err != nil {
			return errors.New(lang.Languagef("host_search_fail", err.Error())) // "查询主机信息失败"
		}
		if 0 == hostID {
			return errors.New(lang.Languagef("host_search_fail", err.Error()))
		}
		bl, hasErr := lgc.IsHostExistInApp(appID, hostID, pheader)
		if nil != hasErr {
			blog.Errorf("check host is exist in app error, params:{appid:%d, hostid:%s}, error:%s", appID, hostID, hasErr.Error())
			return lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader)).Errorf(common.CCErrHostNotINAPPFail, hostID)

		}
		if false == bl {
			blog.Errorf("Host does not belong to the current application; error, params:{appid:%d, hostid:%s}", appID, hostID)
			return lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader)).Errorf(common.CCErrHostNotINAPP, hostID)
		}

	}

	//del host relation from default  module
	conf := &metadata.ModuleHostConfigParams{
		ApplicationID: appID,
		HostID:        hostID,
	}
	result, err := lgc.CoreAPI.HostController().Module().DelDefaultModuleHostConfig(context.Background(), pheader, conf)
	if err != nil || (err == nil && !result.Result) {
		return lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader)).Errorf(common.CCErrHostDELResourcePool, hostID)
	}

	cfg := &metadata.ModuleHostConfigParams{
		ApplicationID: appID,
		ModuleID:      []int64{moduleID},
		HostID:        hostID,
	}
	result, err = lgc.CoreAPI.HostController().Module().AddModuleHostConfig(context.Background(), pheader, cfg)
	if err != nil {
		blog.Error("enter ip, add module host config failed, err: %v", err)
		return errors.New(lang.Languagef("host_agent_add_host_module_fail", err.Error()))
	} else if err == nil && !result.Result {
		blog.Errorf("enter ip, add module host config failed, err: %v", result.ErrMsg)
		return errors.New(lang.Languagef("host_agent_add_host_module_fail", result.ErrMsg))
	}

	audit := lgc.NewHostLog(pheader, ownerID)
	if err := audit.WithPrevious(strconv.FormatInt(hostID, 10), nil); err != nil {
		return fmt.Errorf("audit host log, but get pre data failed, err: %v", err)
	}
	content := audit.GetContent(hostID)
	log := common.KvMap{common.BKContentField: content, common.BKOpDescField: "enter ip host", common.BKHostInnerIPField: audit.ip, common.BKOpTypeField: auditoplog.AuditOpTypeAdd, "inst_id": hostID}
	aResult, err := lgc.CoreAPI.AuditController().AddHostLog(context.Background(), ownerID, strconv.FormatInt(appID, 10), user, pheader, log)
	if err != nil || (err == nil && !aResult.Result) {
		return fmt.Errorf("audit host module log failed, err: %v, %v", err, aResult.ErrMsg)
	}

	hmAudit := lgc.NewHostModuleLog(pheader, []int64{hostID})
	if err := hmAudit.WithPrevious(); err != nil {
		return fmt.Errorf("audit host module log, but get pre data failed, err: %v", err)
	}
	if err := hmAudit.SaveAudit(strconv.FormatInt(appID, 10), user, "host module change"); err != nil {
		return fmt.Errorf("audit host module log, but get pre data failed, err: %v", err)
	}
	return nil
}

func (lgc *Logics) GetHostInfoByConds(pheader http.Header, cond map[string]interface{}) ([]map[string]interface{}, error) {
	query := &metadata.QueryInput{
		Condition: cond,
		Start:     0,
		Limit:     common.BKNoLimit,
		Sort:      common.BKHostIDField,
	}

	result, err := lgc.CoreAPI.HostController().Host().GetHosts(context.Background(), pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get hosts info failed, err: %v, %v", err, result.ErrMsg)
	}

	return result.Data.Info, nil
}

// HostSearch search host by mutiple condition
const (
	SplitFlag      = "##"
	TopoSetName    = "TopSetName"
	TopoModuleName = "TopModuleName"
)

func (lgc *Logics) SearchHost(pheader http.Header, data *metadata.HostCommonSearch, isDetail bool) (*metadata.SearchHost, error) {
	var hostCond, appCond, setCond, moduleCond, mainlineCond metadata.SearchCondition
	objectCondMap := make(map[string][]metadata.ConditionItem, 0)
	appIDArr := make([]int64, 0)
	setIDArr := make([]int64, 0)
	moduleIDArr := make([]int64, 0)
	hostIDArr := make([]int64, 0)
	instAsstHostIDArr := make([]int64, 0)
	objSetIDArr := make([]int64, 0)
	disAppIDArr := make([]int64, 0)
	disSetIDArr := make([]int64, 0)
	disModuleIDArr := make([]int64, 0)
	hostAppConfig := make(map[int64][]int64)
	hostSetConfig := make(map[int64][]int64)
	hostModuleConfig := make(map[int64][]int64)
	moduleSetConfig := make(map[int64]int64)
	setAppConfig := make(map[int64]int64)
	setIDNameMap := make(map[int64]string)

	hostModuleMap := make(map[int64]mapstr.MapStr)
	hostSetMap := make(map[int64]mapstr.MapStr)
	hostAppMap := make(map[int64]mapstr.MapStr)

	totalInfo := make([]mapstr.MapStr, 0)
	moduleHostConfig := make(map[string][]int64)

	for _, object := range data.Condition {
		if object.ObjectID == common.BKInnerObjIDHost {
			hostCond = object
		} else if object.ObjectID == common.BKInnerObjIDSet {
			setCond = object
		} else if object.ObjectID == common.BKInnerObjIDModule {
			moduleCond = object
		} else if object.ObjectID == common.BKInnerObjIDApp {
			appCond = object
		} else if object.ObjectID == common.BKINnerObjIDObject {
			mainlineCond = object
		} else {
			objectCondMap[object.ObjectID] = object.Condition
		}
	}

	//search appID by cond
	if -1 != data.AppID && 0 != data.AppID {
		appCond.Condition = append(appCond.Condition, metadata.ConditionItem{
			Field:    common.BKAppIDField,
			Operator: common.BKDBEQ,
			Value:    data.AppID,
		})
	}

	var err error
	if len(appCond.Condition) > 0 {
		appIDArr, err = lgc.GetAppIDByCond(pheader, appCond.Condition)
		if err != nil {
			return nil, err
		}
	}
	//search mainline object by cond
	if len(mainlineCond.Condition) > 0 {
		objSetIDArr, err = lgc.GetSetIDByObjectCond(pheader, data.AppID, mainlineCond.Condition)
		if err != nil {
			return nil, err
		}
	}
	//search set by appcond
	if len(setCond.Condition) > 0 || len(mainlineCond.Condition) > 0 {
		if len(appCond.Condition) > 0 {
			setCond.Condition = append(setCond.Condition, metadata.ConditionItem{
				Field:    common.BKAppIDField,
				Operator: common.BKDBIN,
				Value:    appIDArr,
			})
		}
		if len(mainlineCond.Condition) > 0 {
			setCond.Condition = append(setCond.Condition, metadata.ConditionItem{
				Field:    common.BKSetIDField,
				Operator: common.BKDBIN,
				Value:    objSetIDArr,
			})
		}
		setIDArr, err = lgc.GetSetIDByCond(pheader, setCond.Condition)
		if err != nil {
			return nil, err
		}
	}

	//search host id by object
	firstCond := true
	if len(objectCondMap) > 0 {
		for objID, objCond := range objectCondMap {
			instIDArr, err := lgc.GetObjectInstByCond(pheader, objID, objCond)
			if err != nil {
				return nil, err
			}
			instHostIDArr, err := lgc.GetHostIDByInstID(pheader, objID, instIDArr)
			if err != nil {
				return nil, err
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
	if len(moduleCond.Condition) > 0 {
		if len(setCond.Condition) > 0 {
			moduleCond.Condition = append(moduleCond.Condition, metadata.ConditionItem{
				Field:    common.BKSetIDField,
				Operator: common.BKDBIN,
				Value:    setIDArr,
			})
		}
		if len(appCond.Condition) > 0 {
			moduleCond.Condition = append(moduleCond.Condition, metadata.ConditionItem{
				Field:    common.BKAppIDField,
				Operator: common.BKDBIN,
				Value:    appIDArr,
			})
		}
		//search module by cond
		moduleIDArr, err = lgc.GetModuleIDByCond(pheader, moduleCond.Condition)
		if err != nil {
			return nil, err
		}
	}

	if len(appCond.Condition) > 0 {
		moduleHostConfig[common.BKAppIDField] = appIDArr
	}
	if len(setCond.Condition) > 0 {
		moduleHostConfig[common.BKSetIDField] = setIDArr
	}
	if len(moduleCond.Condition) > 0 {
		moduleHostConfig[common.BKModuleIDField] = moduleIDArr
	}
	if len(objectCondMap) > 0 {
		moduleHostConfig[common.BKHostIDField] = instAsstHostIDArr
	}
	hostIDArr, err = lgc.GetHostIDByCond(pheader, moduleHostConfig)
	if err != nil {
		return nil, err
	}

	if len(appCond.Condition) > 0 || len(setCond.Condition) > 0 || len(moduleCond.Condition) > 0 || -1 != data.AppID {
		hostCond.Condition = append(hostCond.Condition, metadata.ConditionItem{
			Field:    common.BKHostIDField,
			Operator: common.BKDBIN,
			Value:    hostIDArr,
		})
	}
	if 0 != len(hostCond.Fields) {
		hostCond.Fields = append(hostCond.Fields, common.BKHostIDField)
	}

	condition := make(map[string]interface{})
	hostParse.ParseHostParams(hostCond.Condition, condition)
	hostParse.ParseHostIPParams(data.Ip, condition)

	query := &metadata.QueryInput{
		Condition: condition,
		Start:     data.Page.Start,
		Limit:     data.Page.Limit,
		Sort:      data.Page.Sort,
	}
	gResult, err := lgc.CoreAPI.HostController().Host().GetHosts(context.Background(), pheader, query)
	if err != nil || (err == nil && !gResult.Result) {
		blog.Errorf("get hosts failed, err: %v", err)
		return nil, err
	}

	hostResult := gResult.Data.Info

	// deal the host
	var retStrErr error
	page := metadata.BasePage{Start: 0, Limit: common.BKNoLimit}
	if true == isDetail {
		hostResult, retStrErr = lgc.GetInstAsstDetailsSub(pheader, common.BKInnerObjIDHost, common.BKDefaultOwnerID, hostResult, page)
	} else {
		hostResult, retStrErr = lgc.GetInstDetailsSub(pheader, common.BKInnerObjIDHost, common.BKDefaultOwnerID, hostResult, page)
	}
	if nil != retStrErr {
		blog.Errorf("failed to replace association object, error code is %s, input:%v", retStrErr.Error(), data)
		return nil, retStrErr
	}

	resHostIDArr := make([]int64, 0)
	queryCond := make(map[string][]int64)
	for _, j := range hostResult {
		hostID, err := util.GetInt64ByInterface(j[common.BKHostIDField])
		if err != nil {
			return nil, err
		}
		resHostIDArr = append(resHostIDArr, hostID)
	}
	queryCond[common.BKHostIDField] = resHostIDArr

	mhconfig, err := lgc.GetConfigByCond(pheader, queryCond)
	if err != nil {
		return nil, err
	}
	blog.V(3).Infof("get modulehostconfig map:%v", mhconfig)
	for _, mh := range mhconfig {
		hostID := mh[common.BKHostIDField]
		hostAppConfig[hostID] = append(hostAppConfig[hostID], mh[common.BKAppIDField])
		hostSetConfig[hostID] = append(hostSetConfig[hostID], mh[common.BKSetIDField])
		hostModuleConfig[hostID] = append(hostModuleConfig[hostID], mh[common.BKModuleIDField])

		moduleSetConfig[mh[common.BKModuleIDField]] = mh[common.BKSetIDField]
		setAppConfig[mh[common.BKSetIDField]] = mh[common.BKAppIDField]

		disAppIDArr = append(disAppIDArr, mh[common.BKAppIDField])
		disSetIDArr = append(disSetIDArr, mh[common.BKSetIDField])
		disModuleIDArr = append(disModuleIDArr, mh[common.BKModuleIDField])
	}
	if nil != appCond.Fields {
		exist := util.InArray(common.BKAppIDField, appCond.Fields)
		if 0 != len(appCond.Fields) && !exist {
			appCond.Fields = append(appCond.Fields, common.BKAppIDField)
		}
		cond := make(map[string]interface{})
		celld := make(map[string]interface{})
		celld[common.BKDBIN] = disAppIDArr
		cond[common.BKAppIDField] = celld
		fields := strings.Join(appCond.Fields, ",")
		hostAppMap, err = lgc.GetAppMapByCond(pheader, fields, cond)
		if err != nil {
			return nil, err
		}
	}
	if nil != setCond.Fields {
		exist := util.InArray(common.BKSetIDField, setCond.Fields)
		if !exist && 0 != len(setCond.Fields) {
			setCond.Fields = append(setCond.Fields, common.BKSetIDField)
		}
		cond := make(map[string]interface{})
		celld := make(map[string]interface{})
		celld[common.BKDBIN] = disSetIDArr
		cond[common.BKSetIDField] = celld
		fields := strings.Join(setCond.Fields, ",")
		hostSetMap, err = lgc.GetSetMapByCond(pheader, fields, cond)
		if err != nil {
			return nil, err
		}
	}
	if nil != moduleCond.Fields {
		exist := util.InArray(common.BKModuleIDField, moduleCond.Fields)
		if !exist && 0 != len(moduleCond.Fields) {
			moduleCond.Fields = append(moduleCond.Fields, common.BKModuleIDField)
		}
		cond := make(map[string]interface{})
		celld := make(map[string]interface{})
		celld[common.BKDBIN] = disModuleIDArr
		cond[common.BKModuleIDField] = celld
		fields := strings.Join(moduleCond.Fields, ",")
		hostModuleMap, err = lgc.GetModuleMapByCond(pheader, fields, cond)
		if err != nil {
			return nil, err
		}
	}

	//com host info
	for _, host := range hostResult {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			return nil, fmt.Errorf("invalid hostid: %v", err)
		}

		//appdata
		hostAppIDArr, ok := hostAppConfig[hostID]
		if false == ok {
			continue
		}
		hostAppData := make([]interface{}, 0)
		for _, appID := range hostAppIDArr {
			appInfo, mapOk := hostAppMap[appID]
			if mapOk {
				hostAppData = append(hostAppData, appInfo)
			}
		}
		hostData := make(map[string]interface{})
		hostData[common.BKInnerObjIDApp] = hostAppData

		//setdata
		hostSetIDArr, ok := hostSetConfig[hostID]
		hostSetData := make([]interface{}, 0)
		for _, setID := range hostSetIDArr {
			setInfo, isOk := hostSetMap[setID]
			if false == isOk {
				continue
			}
			appID := setAppConfig[setID]
			if false == isOk {
				blog.Warnf("hostSearch not found set id, setID:%d, setAppConfig:%v, input:%v", setID, setAppConfig, data)

				continue
			}
			appInfo, isOk := hostAppMap[appID]
			if false == isOk {
				blog.Warnf("hostSearch not found application id, appID:%d, hostAppMap:%v, input:%v", appID, hostAppMap, data)
				continue
			}

			appName, err := appInfo.String(common.BKAppNameField)
			if nil != err {
				blog.Warnf("hostSearch not found application name,  appID:%d, appInfo:%v, input:%v", appID, appInfo, data)
				continue
			}

			setName, err := setInfo.String(common.BKSetNameField)
			if nil != err {
				blog.Warnf("hostSearch not found set name, setInfo:%d, input:%v", setInfo, data)
				continue
			}
			datacp := make(map[string]interface{})
			for key, val := range setInfo {
				datacp[key] = val
			}
			datacp[TopoSetName] = appName + SplitFlag + setName
			hostSetData = append(hostSetData, datacp)
			setIDNameMap[setID] = setName
		}
		hostData[common.BKInnerObjIDSet] = hostSetData

		//moduledata
		hostModuleIDArr, ok := hostModuleConfig[hostID]
		hostModuleData := make([]interface{}, 0)
		for _, ModuleID := range hostModuleIDArr {
			moduleInfo, ok := hostModuleMap[ModuleID]
			if false == ok {
				blog.Warnf("hostSearch not found module id, moduleID:%d, hostModuleMap:%v, input:%v", ModuleID, hostModuleMap, data)
				continue
			}
			setID := moduleSetConfig[ModuleID]
			if false == ok {
				blog.Warnf("hostSearch not found application id, moduleID:%d, moduleSetConfig:%v, input:%v", ModuleID, moduleSetConfig, data)
				continue
			}
			appID := setAppConfig[setID]
			if false == ok {
				blog.Warnf("hostSearch not found application id, moduleID:%d, moduleSetConfig:%v, input:%v", ModuleID, setAppConfig, data)
				continue
			}
			appInfo, ok := hostAppMap[appID]
			if false == ok {
				blog.Warnf("hostSearch not found application info, moduleID:%d, moduleSetConfig:%v, input:%v", ModuleID, hostAppMap, data)
				continue
			}

			appName, err := appInfo.String(common.BKAppNameField)
			if nil != err {
				blog.Warnf("hostSearch not found application name, moduleID:%d, moduleSetConfig:%v, input:%v", ModuleID, appInfo, data)
				continue
			}

			moduleName, err := moduleInfo.String(common.BKModuleNameField)
			if nil != err {
				blog.Warnf("hostSearch not found module name,input:%v, item:%v, input:%v", mainlineCond, moduleInfo, data)
				continue
			}
			datacp := make(map[string]interface{})
			for key, val := range moduleInfo {
				datacp[key] = val
			}
			setName := setIDNameMap[setID]
			datacp[TopoModuleName] = appName + SplitFlag + setName + SplitFlag + moduleName
			hostModuleData = append(hostModuleData, datacp)
		}
		hostData[common.BKInnerObjIDModule] = hostModuleData

		hostData[common.BKInnerObjIDHost] = host
		totalInfo = append(totalInfo, hostData)
	}

	return &metadata.SearchHost{
		Info:  totalInfo,
		Count: gResult.Data.Count,
	}, nil
}

func (lgc *Logics) GetHostIDByCond(pheader http.Header, cond map[string][]int64) ([]int64, error) {
	result, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(context.Background(), pheader, cond)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("%v, %v", err, result.ErrMsg)
	}

	hostIDs := make([]int64, 0)
	for _, val := range result.Data {
		hostIDs = append(hostIDs, val.HostID)
	}

	return hostIDs, nil
}
