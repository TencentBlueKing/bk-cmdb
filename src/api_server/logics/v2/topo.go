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
	"fmt"
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
)

func (lgc *Logics) GetAppTopo(user string, pheader http.Header, appID int64, conds mapstr.MapStr) (mapstr.MapStr, errors.CCError) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	rid := util.GetHTTPCCRequestID(pheader)

	apps, err := lgc.getAppInfo(user, pheader, appID)
	if err != nil {
		return nil, err
	}
	if 0 == len(apps) {
		blog.Errorf("GetAppTopo not find app by id:%d", appID)
		return nil, defErr.Error(common.CCErrCommNotFound)
	}
	appInfo := apps[0]
	appName, err := appInfo.String(common.BKAppNameField)
	if err != nil {
		appName = ""
	}

	ret := mapstr.MapStr{
		"Level":           3,
		"ApplicationName": appName,
		"ApplicationID":   strconv.FormatInt(appID, 10),
		"Children":        make([]interface{}, 0),
	}

	moduleFields := []string{common.BKModuleIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField, common.BKDefaultField}
	modules, err := lgc.getModulesByConds(user, pheader, strconv.FormatInt(appID, 10), conds, moduleFields, "")
	if err != nil {
		return nil, err
	}
	modulesMap, err := getModuleMap(modules, appID, rid, defErr)
	if err != nil {
		return nil, err
	}

	sets, err := lgc.getSets(user, pheader, appID)
	if err != nil {
		return nil, err
	}
	if 0 == len(sets) {
		blog.Errorf("GetAppTopo not find set by app id:%d, rid:%s", appID, rid)
		return nil, defErr.Error(common.CCErrCommNotFound)
	}
	retSets := make([]mapstr.MapStr, 0)
	for _, setInfo := range sets {

		setID, err := util.GetInt64ByInterface(setInfo[common.BKSetIDField])
		if nil != err {
			blog.Errorf("GetAppTopo get set id by getsets  return  info:%v  error:%v", setInfo, err)
			return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		setName, _ := setInfo[common.BKSetNameField]
		if nil == setName {
			setName = ""
		}
		moduleArr, ok := modulesMap[setID]
		if false == ok {
			moduleArr = make([]mapstr.MapStr, 0)
		}
		retSets = append(retSets, mapstr.MapStr{
			"SetName":  setName,
			"SetID":    strconv.FormatInt(setID, 10),
			"Children": moduleArr,
		})

	}
	ret["Children"] = retSets

	return ret, nil

}

func getModuleMap(modules []mapstr.MapStr, appID int64, rid string, defErr errors.DefaultCCErrorIf) (map[int64][]mapstr.MapStr, errors.CCError) {

	// key is setID
	modulesMap := make(map[int64][]mapstr.MapStr)
	if 0 == len(modules) {
		blog.Errorf("GetAppTopo not find module by app id:%d,rid:%s", appID, rid)
		return nil, defErr.Error(common.CCErrCommNotFound)
	}
	strAppID := strconv.FormatInt(appID, 10)
	for _, module := range modules {

		setID, err := module.Int64(common.BKSetIDField)
		if nil != err {
			blog.Errorf("GetAppTopo get set id by getmodues  return info  error, module info:%v  error, rid:%s", module, rid)
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "module", "SetID", "int", err.Error())
		}
		moduleName, err := module.String(common.BKModuleNameField)
		if err != nil {
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "module", "ModuleName", "string", err.Error())

		}
		moduleID, err := util.GetInt64ByInterface(module[common.BKModuleIDField])
		if nil != err {
			blog.Errorf("GetAppTopo get module id by getmodues  return info  error, module info:%v  error, rid:%s", module, rid)
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "module", "ModuleID", "int", err.Error())
		}
		_, ok := modulesMap[setID]
		if false == ok {
			modulesMap[setID] = make([]mapstr.MapStr, 0)
		}
		modulesMap[setID] = append(modulesMap[setID], mapstr.MapStr{
			"SetID":         strconv.FormatInt(setID, 10),
			"ModuleID":      strconv.FormatInt(moduleID, 10),
			"ModuleName":    moduleName,
			"HostNum":       "0",
			"ApplicationID": strAppID,
			"ObjID":         common.BKInnerObjIDModule,
		})

	}

	return modulesMap, nil
}

// GetDefaultTopo get resource topo
func GetDefaultTopo(req *restful.Request, appID string, topoApi string) (map[string]interface{}, error) {
	defaultTopo := make(map[string]interface{})
	url := fmt.Sprintf("%s/topo/v1/topo/internal/%s/%s", topoApi, common.BKDefaultOwnerID, appID)
	res, err := httpcli.ReqHttp(req, url, common.HTTPSelectGet, nil)
	if err != nil {
		blog.Errorf("getDefaultTopo error:%v", err)
		return nil, err
	}

	//appIDInt ,_:= strconv.Atoi(appID)
	resMap := make(map[string]interface{})

	err = json.Unmarshal([]byte(res), &resMap)
	if err != nil {
		blog.Errorf("json Unmarshal data:%s, err:%s", res, err.Error())
		return nil, err
	}
	if resMap["result"].(bool) {

		resMapData, ok := resMap["data"].(map[string]interface{})
		if false == ok {
			blog.Errorf("assign error resMap:%v", resMap)
			return defaultTopo, nil
		}
		defaultTopo["Children"] = make([]map[string]interface{}, 0)
		resModule, ok := resMapData["module"].([]interface{})
		if false == ok {
			blog.Errorf("assign error resMapData:%v", resMapData)
			return defaultTopo, nil
		}
		for _, module := range resModule {
			Module, ok := module.(map[string]interface{})
			if false == ok {
				blog.Errorf("assign error module:%v", module)
				continue
			}
			moduleMap := map[string]interface{}{
				"ModuleID":      Module[common.BKModuleIDField],
				"ModuleName":    Module[common.BKModuleNameField],
				"ApplicationID": appID,
				"ObjID":         "module",
			}
			defaultTopo["Children"] = append(defaultTopo["Children"].([]map[string]interface{}), moduleMap)
		}
		defaultTopo["SetName"] = common.DefaultResSetName
		setIdInt, _ := util.GetIntByInterface(resMap["data"].(map[string]interface{})[common.BKSetIDField])
		setIdStr := strconv.Itoa(setIdInt)
		defaultTopo["SetID"] = setIdStr
		defaultTopo["ObjID"] = "set"
	}

	return defaultTopo, nil
}

// AppendDefaultTopo combin  idle pool set
func AppendDefaultTopo(topo map[string]interface{}, defaultTopo map[string]interface{}) map[string]interface{} {
	topoChildren, ok := topo["Children"].([]map[string]interface{})
	if false == ok {
		err := fmt.Sprintf("assign error topo.Children is not []map[string]interface{},topo:%v", topo)
		blog.Error(err)
		return nil
	}

	children := make([]map[string]interface{}, 0)
	children = append(children, defaultTopo)
	for _, child := range topoChildren {
		children = append(children, child)
	}
	topo["Children"] = children
	return topo
}

// SetModuleHostCount get set host count
func (lgc *Logics) SetModuleHostCount(data []mapstr.MapStr, user string, pheader http.Header) error {
	for _, itemMap := range data {

		switch itemMap["ObjID"] {
		case common.BKInnerObjIDModule:

			mouduleId, getErr := util.GetIntByInterface(itemMap["ModuleID"])
			if nil != getErr {
				blog.Errorf("%v, %v", getErr, itemMap)
				return getErr
			}
			appId, getErr := util.GetIntByInterface(itemMap["ApplicationID"])
			if nil != getErr {
				blog.Errorf("%v, %v", getErr, itemMap)
				return getErr
			}
			hostNum, getErr := lgc.GetModuleHostCount(appId, mouduleId, user, pheader)
			if nil != getErr {
				return getErr
			}
			itemMap["HostNum"] = hostNum
		}

		if nil != itemMap["Children"] {
			children, ok := itemMap["Children"].([]mapstr.MapStr)
			if false == ok {
				children = make([]mapstr.MapStr, 0)
			}
			lgc.SetModuleHostCount(children, user, pheader)
		} else {
			children := make([]interface{}, 0)
			itemMap["Children"] = children
		}
	}
	return nil
}

// GetModuleHostCount get module host count
func (lgc *Logics) GetModuleHostCount(appID, mouduleID interface{}, user string, pheader http.Header) (int, error) {
	rid := util.GetHTTPCCRequestID(pheader)

	param := mapstr.MapStr{
		common.BKAppIDField:    appID,
		common.BKModuleIDField: []interface{}{mouduleID},
	}

	result, err := lgc.CoreAPI.HostServer().HostSearchByModuleID(context.Background(), pheader, param)
	if nil != err {
		blog.Errorf("getModuleHostCount , error:%v", err)
		return 0, err
	}
	if !result.Result {
		blog.Errorf("getModuleHostCount http response error, err code:%d, err msg:%d, cond:%+v, rid:%s", result.Code, result.ErrMsg, param, rid)
		return 0, err
	}

	rspV3MapData, ok := result.Data.([]interface{})
	if false == ok {
		blog.Errorf("assign error rspV3Map.data is not []interface{}, rspV3Map:%v", result.Data)
		return 0, nil
	}
	return len(rspV3MapData), nil
}

func (lgc *Logics) getSets(user string, pheader http.Header, appID int64) ([]mapstr.MapStr, errors.CCError) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	rid := util.GetHTTPCCRequestID(pheader)

	page := mapstr.MapStr{"start": 0, "limit": common.BKNoLimit}
	param := &params.SearchParams{Page: page, Condition: mapstr.MapStr{}}

	appIDStr := strconv.FormatInt(appID, 10)
	result, err := lgc.CoreAPI.TopoServer().Instance().SearchSet(context.Background(), user, appIDStr, pheader, param)
	if err != nil {
		blog.Errorf("get sets   error:%v, cond:%+v,rid:%s", err, param, rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get sets info http response error, err code:%d, err msg:%s, cond:%+v,rid:%s", result.Code, result.ErrMsg, param, rid)
		return nil, defErr.New(result.Code, result.ErrMsg)
	}
	return result.Data.Info, nil
}

func (lgc *Logics) getModules(user string, pheader http.Header, appID int64) ([]mapstr.MapStr, errors.CCError) {

	fields := []string{common.BKModuleIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField, common.BKDefaultField}

	return lgc.getModulesByConds(user, pheader, strconv.FormatInt(appID, 10), nil, fields, "")
}

func (lgc *Logics) getModulesByConds(user string, pheader http.Header, appIDStr string, conds map[string]interface{}, fields []string, sort string) ([]mapstr.MapStr, errors.CCError) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	rid := util.GetHTTPCCRequestID(pheader)

	searchParams := mapstr.MapStr{}
	searchParams["fields"] = fields
	searchParams["condition"] = conds
	searchParams["page"] = mapstr.MapStr{"start": 0, "limit": common.BKNoLimit, "sort": sort}
	result, err := lgc.CoreAPI.TopoServer().OpenAPI().SearchModuleByApp(context.Background(), appIDStr, pheader, searchParams)
	if err != nil {
		blog.Errorf("getModules   error:%v, cond:%+v,rid:%s", err, conds, rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get app info http response error, err code:%d, err msg:%s, cond:%+v,rid:%s", result.Code, result.ErrMsg, searchParams, rid)
		return nil, defErr.New(result.Code, result.ErrMsg)
	}
	return result.Data.Info, nil
}

func (lgc *Logics) getAppInfo(user string, pheader http.Header, appID int64) ([]mapstr.MapStr, errors.CCError) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	rid := util.GetHTTPCCRequestID(pheader)

	page := mapstr.MapStr{"start": 0, "limit": 2}
	condition := mapstr.MapStr{common.BKAppIDField: appID}
	param := &params.SearchParams{Condition: condition, Page: page}
	result, err := lgc.CoreAPI.TopoServer().Instance().SearchApp(context.Background(), user, pheader, param)
	if nil != err {
		blog.Errorf("get app info error:%v, cond:%+v,rid:%s", err, condition, rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get app info http response error, err code:%d, err msg:%s, cond:%+v,rid:%s", result.Code, result.ErrMsg, condition, rid)
		return nil, defErr.New(result.Code, result.ErrMsg)
	}

	return result.Data.Info, nil
}

func getRspV3DataInfo(logPrex string, result bool, code int, data interface{}) ([]interface{}, int) {

	if false == result {
		return nil, code
	}
	dataMap, ok := data.(map[string]interface{})
	if false == ok {
		blog.Errorf("%s rspV3 json get data.info error, body:%#v", logPrex, dataMap)
		return nil, common.CCErrCommJSONUnmarshalFailed
	}
	dataInfo, ok := dataMap["info"].([]interface{})
	if false == ok {
		blog.Errorf("%s rspV3 json get data.info error, body:%#v", logPrex, dataMap)
		return nil, common.CCErrCommJSONUnmarshalFailed
	}
	return dataInfo, 0
}

func (lgc *Logics) CheckAppTopoIsThreeLevel(user string, pheader http.Header) (bool, error) {

	result, err := lgc.CoreAPI.TopoServer().Object().SelectModel(context.Background(), util.GetOwnerID(pheader), pheader)
	if nil != err {
		blog.Errorf("CheckAppTopoIsThreeLevel  error:%s ", err.Error())
		return false, err
	}
	if false == result.Result {
		blog.Errorf("CheckAppTopoIsThreeLevel reply:%s ", result.ErrMsg)
		return false, fmt.Errorf(result.ErrMsg)
	}

	for _, item := range result.Data {
		if common.BKInnerObjIDApp == item.ObjID {
			if common.BKInnerObjIDSet != item.NextObj {
				return false, nil
			}
		}

		if common.BKInnerObjIDSet == item.ObjID {
			if common.BKInnerObjIDModule != item.NextObj {
				return false, nil
			}
		}
	}

	return true, nil
}

func (lgc *Logics) GetAppListByOwnerIDAndUin(ctx context.Context, pheader http.Header, ownerID, uin string) (appInfoArr []mapstr.MapStr, errCode int, err errors.CCError) {

	pheader.Set(common.BKHTTPOwnerID, ownerID)
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	userLike := mapstr.MapStr{
		common.BKDBLIKE: fmt.Sprintf("^%s,|,%s,|,%s$|^%s$", uin, uin, uin, uin),
	}
	var condItemArr []mapstr.MapStr
	var conds mapstr.MapStr
	if ownerID != uin {
		condItemArr = append(condItemArr, mapstr.MapStr{common.CreatorField: userLike})
		condItemArr = append(condItemArr, mapstr.MapStr{common.BKMaintainersField: userLike})
		conds = mapstr.MapStr{common.BKDBOR: condItemArr}
	}
	s := &metadata.SearchParams{}
	if nil != conds {
		s.Condition = conds
	}

	result, err := lgc.CoreAPI.TopoServer().Instance().InstSearch(ctx, ownerID, common.BKInnerObjIDApp, pheader, s)
	if nil != err {
		blog.Errorf("GetAppListByOwnerIDAndUin http do error, error:%s, request-id:%s", err.Error(), util.GetHTTPCCRequestID(pheader))
		return nil, common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("GetAppListByOwnerIDAndUin http reply error, error code:%d, error message:%s, request-id:%s", result.Code, result.ErrMsg, util.GetHTTPCCRequestID(pheader))
		return nil, common.CCErrCommHTTPDoRequestFailed, defErr.New(result.Code, result.ErrMsg)
	}

	if 1 <= len(result.Data.Info) {
		return nil, common.CCErrCommNotFound, defErr.Error(common.CCErrCommNotFound)
	}

	return result.Data.Info, 0, nil

}
