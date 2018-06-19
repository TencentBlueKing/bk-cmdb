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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/instapi"
	"configcenter/src/scene_server/validator"
	"github.com/bitly/go-simplejson"
)

func (lgc *Logics) GetConfigByCond(pheader http.Header, cond map[string][]int64) ([]map[string]int64, error) {
	configArr := make([]map[string]int64, 0)

	if 0 == len(cond) {
		return configArr, nil
	}

	result, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(context.Background(), pheader, cond)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get module host config failed, err: %v, %v", err, result.ErrMsg)
	}

	for _, infos := range result.Data {
		info := infos.(map[string]interface{})
		hostID, err := info[common.BKHostIDField].(json.Number).Int64()
		if err != nil {
			return nil, err
		}

		setID, err := info[common.BKSetIDField].(json.Number).Int64()
		if err != nil {
			return nil, err
		}

		moduleID, err := info[common.BKModuleIDField].(json.Number).Int64()
		if err != nil {
			return nil, err
		}

		appID, err := info[common.BKAppIDField].(json.Number).Int64()
		if err != nil {
			return nil, err
		}

		data := make(map[string]int64)
		data[common.BKAppIDField] = appID
		data[common.BKSetIDField] = setID
		data[common.BKModuleIDField] = moduleID
		data[common.BKHostIDField] = hostID
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
		valid := validator.NewValidMap(common.BKDefaultOwnerID, common.BKInnerObjIDHost, ObjAddr, forward, errHandle)
		_, hasErr = valid.ValidMap(host, "create", 0)

		if nil != hasErr {
			return hasErr
		}

		result, err := lgc.CoreAPI.HostController().Host().AddHost(context.Background(), pheader, host)
		if err != nil || (err == nil && !result.Result) {
			return errors.New(lang.Languagef("host_agent_add_host_fail", err.Error()))
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
	if err != nil || (err == nil && !result.Result) {
		blog.Error("enter ip, add module host config failed, err: %v", err)
		return errors.New(lang.Languagef("host_agent_add_host_module_fail", err.Error()))
	}

	//prepare the log
	hostLogFields, _ := GetHostLogFields(req, ownerID, ObjAddr)
	logObj := NewHostLog(req, common.BKDefaultOwnerID, "", hostAddr, ObjAddr, hostLogFields)
	content, _ := logObj.GetHostLog(fmt.Sprintf("%d", hostID), false)
	logAPIClient := sourceAuditAPI.NewClient(auditAddr)
	logAPIClient.AuditHostLog(hostID, content, "enter IP HOST", ip, ownerID, fmt.Sprintf("%d", appID), user, auditoplog.AuditOpTypeAdd)
	logClient, err := NewHostModuleConfigLog(req, nil, hostAddr, ObjAddr, auditAddr)
	logClient.SetHostID([]int{hostID})
	logClient.SetDesc("host module change")
	logClient.SaveLog(fmt.Sprintf("%d", appID), user)
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
func (lgc *Logics) SearchHost(pheader http.Header, data *metadata.HostCommonSearch, isDetail bool) (interface{}, error) {
	var hostCond, appCond, setCond, moduleCond, mainlineCond metadata.SearchCondition
	objectCondMap := make(map[string][]interface{}, 0)
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

	hostModuleMap := make(map[int64]interface{})
	hostSetMap := make(map[int64]interface{})
	hostAppMap := make(map[int64]interface{})

	result := make(map[string]interface{})
	totalInfo := make([]interface{}, 0)
	moduleHostConfig := make(map[string][]int, 0)

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
		cond := make(map[string]interface{})
		cond["field"] = common.BKAppIDField
		cond["operator"] = common.BKDBEQ
		cond["value"] = data.AppID
		appCond.Condition = append(appCond.Condition, cond)
	}
	if len(appCond.Condition) > 0 {
		var err error
		appIDArr, err = lgc.GetAppIDByCond(pheader, appCond.Condition)
		if err != nil {
			return nil, err
		}
	}
	//search mainline object by cond
	if len(mainlineCond.Condition) > 0 {
		objSetIDArr = GetSetIDByObjectCond(req, objCtrl, data.AppID, mainlineCond.Condition)
	}
	//search set by appcond
	if len(setCond.Condition) > 0 || len(mainlineCond.Condition) > 0 {
		if len(appCond.Condition) > 0 {
			cond := make(map[string]interface{})
			cond["field"] = common.BKAppIDField
			cond["operator"] = common.BKDBIN
			cond["value"] = appIDArr
			setCond.Condition = append(setCond.Condition, cond)
		}
		if len(mainlineCond.Condition) > 0 {
			cond := make(map[string]interface{})
			cond["field"] = common.BKSetIDField
			cond["operator"] = common.BKDBIN
			cond["value"] = objSetIDArr
			setCond.Condition = append(setCond.Condition, cond)
		}
		setIDArr, _ = GetSetIDByCond(req, objCtrl, setCond.Condition)
	}

	//search host id by object
	firstCond := true
	if len(objectCondMap) > 0 {
		for objID, objCond := range objectCondMap {
			instIDArr := GetObjectInstByCond(req, objID, objCtrl, objCond)
			instHostIDArr := GetHostIDByInstID(req, objID, objCtrl, instIDArr)
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
			cond := make(map[string]interface{})
			cond["field"] = common.BKSetIDField
			cond["operator"] = common.BKDBIN
			cond["value"] = setIDArr
			moduleCond.Condition = append(moduleCond.Condition, cond)
		}
		if len(appCond.Condition) > 0 {
			cond := make(map[string]interface{})
			cond["field"] = common.BKAppIDField
			cond["operator"] = common.BKDBIN
			cond["value"] = appIDArr
			moduleCond.Condition = append(moduleCond.Condition, cond)
		}
		//search module by cond
		moduleIDArr, _ = GetModuleIDByCond(req, objCtrl, moduleCond.Condition)
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
	hostIDArr, _ = GetHostIDByCond(req, hostCtrl, moduleHostConfig)

	if len(appCond.Condition) > 0 || len(setCond.Condition) > 0 || len(moduleCond.Condition) > 0 || -1 != data.AppID {
		cond := make(map[string]interface{})
		cond["field"] = common.BKHostIDField
		cond["operator"] = common.BKDBIN
		cond["value"] = hostIDArr
		hostCond.Condition = append(hostCond.Condition, cond)
	}
	if 0 != len(hostCond.Fields) {
		hostCond.Fields = append(hostCond.Fields, common.BKHostIDField)
	}

	url := hostCtrl + "/host/v1/hosts/search"
	start := data.Page.Start
	limit := data.Page.Limit
	sort := data.Page.Sort
	body := make(map[string]interface{})
	body["start"] = start
	body["limit"] = limit
	body["sort"] = sort
	body["fields"] = strings.Join(hostCond.Fields, ",")

	bodyContent, _ := json.Marshal(body)
	blog.Info("Get Host By Cond url :%s", url)
	blog.Info("Get Host By Cond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("Get Host By Cond return :%s", string(reply))
	if err != nil {
		//cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return nil, errors.New(common.CC_Err_Comm_Host_Get_FAIL_STR)
	}
	condition := make(map[string]interface{})
	hostParse.ParseHostParams(hostCond.Condition, condition)
	hostParse.ParseHostIPParams(data.Ip, condition)
	body["condition"] = condition
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
	instapi.Inst.InitInstHelper(hostCtrl, objCtrl)
	var retStrErr int
	if true == isDetail {
		hostResult, retStrErr = instapi.Inst.GetInstAsstDetailsSub(req, common.BKInnerObjIDHost, common.BKDefaultOwnerID, hostResult, map[string]interface{}{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  "",
		})
	} else {
		hostResult, retStrErr = instapi.Inst.GetInstDetailsSub(req, common.BKInnerObjIDHost, common.BKDefaultOwnerID, hostResult, map[string]interface{}{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  "",
		})
	}

	if common.CCSuccess != retStrErr {
		blog.Error("failed to replace association object, error code is %d", retStrErr)
	}

	cnt := hostResult["count"]
	hostInfo := hostResult["info"].([]interface{})
	result["count"] = cnt
	resHostIDArr := make([]int, 0)
	queryCond := make(map[string]interface{})
	for _, j := range hostInfo {
		host := j.(map[string]interface{})
		hostID, _ := host[common.BKHostIDField].(json.Number).Int64()
		resHostIDArr = append(resHostIDArr, int(hostID))

		queryCond[common.BKHostIDField] = resHostIDArr
	}
	mhconfig, _ := GetConfigByCond(req, hostCtrl, queryCond)
	blog.Info("get modulehostconfig map:%v", mhconfig)
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
		//get app fields

		exist := util.InArray(common.BKAppIDField, appCond.Fields)
		if 0 != len(appCond.Fields) && !exist {
			appCond.Fields = append(appCond.Fields, common.BKAppIDField)
		}
		cond := make(map[string]interface{})
		celld := make(map[string]interface{})
		celld[common.BKDBIN] = disAppIDArr
		cond[common.BKAppIDField] = celld
		fields := strings.Join(appCond.Fields, ",")
		hostAppMap, _ = GetAppMapByCond(req, fields, objCtrl, cond)
	}
	if nil != setCond.Fields {
		//get set fields

		exist := util.InArray(common.BKSetIDField, setCond.Fields)
		if !exist && 0 != len(setCond.Fields) {
			setCond.Fields = append(setCond.Fields, common.BKSetIDField)
		}
		cond := make(map[string]interface{})
		celld := make(map[string]interface{})
		celld[common.BKDBIN] = disSetIDArr
		cond[common.BKSetIDField] = celld
		fields := strings.Join(setCond.Fields, ",")
		hostSetMap, _ = GetSetMapByCond(req, fields, objCtrl, cond)
	}
	if nil != moduleCond.Fields {
		//get module fields

		exist := util.InArray(common.BKModuleIDField, moduleCond.Fields)
		if !exist && 0 != len(moduleCond.Fields) {
			moduleCond.Fields = append(moduleCond.Fields, common.BKModuleIDField)
		}
		cond := make(map[string]interface{})
		celld := make(map[string]interface{})
		celld[common.BKDBIN] = disModuleIDArr
		cond[common.BKModuleIDField] = celld
		fields := strings.Join(moduleCond.Fields, ",")
		hostModuleMap, _ = GetModuleMapByCond(req, fields, objCtrl, cond)
	}

	//com host info
	for _, j := range hostInfo {
		host := j.(map[string]interface{})
		hostID, _ := host[common.BKHostIDField].(json.Number).Int64()
		hostID32 := int(hostID)
		hostData := make(map[string]interface{})

		//appdata
		hostAppIDArr, ok := hostAppConfig[hostID32]
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
		hostData[common.BKInnerObjIDApp] = hostAppData

		//setdata
		hostSetIDArr, ok := hostSetConfig[hostID32]
		hostSetData := make([]interface{}, 0)
		for _, setID := range hostSetIDArr {
			setInfo, isOk := hostSetMap[setID]
			if false == isOk {
				continue
			}
			appID := setAppConfig[setID]
			if false == isOk {
				continue
			}
			appInfoI, isOk := hostAppMap[appID]
			if false == isOk {
				continue
			}
			appInfo, isOk := appInfoI.(map[string]interface{})
			if false == isOk {
				continue
			}
			appName, isOk := appInfo[common.BKAppNameField].(string)
			if false == isOk {
				continue
			}
			data, isOk := setInfo.(map[string]interface{})
			if false == isOk {
				continue
			}

			setName, isOk := data[common.BKSetNameField].(string)
			if false == isOk {
				continue
			}
			datacp := make(map[string]interface{})
			for key, val := range data {
				datacp[key] = val
			}
			datacp[TopoSetName] = appName + SplitFlag + setName
			hostSetData = append(hostSetData, datacp)
			setIDNameMap[setID] = setName
		}
		hostData[common.BKInnerObjIDSet] = hostSetData

		//moduledata
		hostModuleIDArr, ok := hostModuleConfig[hostID32]
		hostModuleData := make([]interface{}, 0)
		for _, ModuleID := range hostModuleIDArr {
			moduleInfo, ok := hostModuleMap[ModuleID]
			if false == ok {
				continue
			}
			setID := moduleSetConfig[ModuleID]
			if false == ok {
				continue
			}
			appID := setAppConfig[setID]
			if false == ok {
				continue
			}
			appInfoI, ok := hostAppMap[appID]
			if false == ok {
				continue
			}
			appInfo, ok := appInfoI.(map[string]interface{})
			if false == ok {
				continue
			}
			appName, ok := appInfo[common.BKAppNameField].(string)
			if false == ok {
				continue
			}
			data, ok := moduleInfo.(map[string]interface{})
			if false == ok {
				continue
			}

			moduleName, ok := data[common.BKModuleNameField].(string)
			if false == ok {
				continue
			}
			datacp := make(map[string]interface{})
			for key, val := range data {
				datacp[key] = val
			}
			setName := setIDNameMap[setID]
			datacp[TopoModuleName] = appName + SplitFlag + setName + SplitFlag + moduleName
			hostModuleData = append(hostModuleData, datacp)
		}
		hostData[common.BKInnerObjIDModule] = hostModuleData

		hostData[common.BKInnerObjIDHost] = j
		totalInfo = append(totalInfo, hostData)
	}

	result["info"] = totalInfo
	result["count"] = cnt

	return result, err
}
