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

func (lgc *Logics) GetAppTopo(ctx context.Context, appID int64, conds mapstr.MapStr) (mapstr.MapStr, errors.CCError) {
	defErr := lgc.ccErr

	apps, err := lgc.getAppInfo(ctx, appID)
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
	modules, err := lgc.getModulesByConds(ctx, strconv.FormatInt(appID, 10), conds, moduleFields, "")
	if err != nil {
		return nil, err
	}
	modulesMap, err := getModuleMap(modules, appID, lgc.rid, defErr)
	if err != nil {
		return nil, err
	}

	sets, err := lgc.getSets(ctx, appID)
	if err != nil {
		return nil, err
	}
	if 0 == len(sets) {
		blog.Errorf("GetAppTopo not find set by app id:%d, rid:%s", appID, lgc.rid)
		return nil, defErr.Error(common.CCErrCommNotFound)
	}
	retSets := make([]mapstr.MapStr, 0)
	for _, setInfo := range sets {

		setID, err := setInfo.Int64(common.BKSetIDField)
		if nil != err {
			blog.Errorf("GetAppTopo get set id by getsets  return  info:%v  error:%v.rid:%s", setInfo, err, lgc.rid)
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDSet, common.BKSetIDField, "int", err.Error())
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

// TODO delete
// GetDefaultTopo get resource topo
func GetDefaultTopo_delete(req *restful.Request, appID string, topoApi string) (map[string]interface{}, error) {
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
func AppendDefaultTopo(topo map[string]interface{}, defaultTopo map[string]interface{}, rid string) map[string]interface{} {
	topoChildren, ok := topo["Children"].([]map[string]interface{})
	if false == ok {
		blog.Warnf("assign error topo.Children is not []map[string]interface{},topo:%#v,rid:%s", topo, rid)
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
func (lgc *Logics) SetModuleHostCount(ctx context.Context, data []mapstr.MapStr) error {
	for _, itemMap := range data {

		switch itemMap["ObjID"] {
		case common.BKInnerObjIDModule:

			mouduleID, getErr := itemMap.Int64("ModuleID")
			if nil != getErr {
				blog.Errorf("SetModuleHostCount error. err:%v, info:%#v, rid:%s", getErr, itemMap, lgc.rid)
				return lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, "module", "ModuleID", "int", getErr.Error())
			}
			appID, getErr := itemMap.Int64("ApplicationID")
			if nil != getErr {
				blog.Errorf("SetModuleHostCount error. err:%v, info:%#v, rid:%s", getErr, itemMap)
				return lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, "App", "ApplicationID", "int", getErr.Error())
			}
			hostNum, getErr := lgc.GetModuleHostCount(ctx, appID, mouduleID)
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
			lgc.SetModuleHostCount(ctx, children)
		} else {
			children := make([]interface{}, 0)
			itemMap["Children"] = children
		}
	}
	return nil
}

// GetModuleHostCount get module host count
func (lgc *Logics) GetModuleHostCount(ctx context.Context, appID, mouduleID interface{}) (int, error) {
	param := mapstr.MapStr{
		common.BKAppIDField:    appID,
		common.BKModuleIDField: []interface{}{mouduleID},
	}

	result, err := lgc.CoreAPI.HostServer().HostSearchByModuleID(ctx, lgc.header, param)
	if nil != err {
		blog.Errorf("getModuleHostCount http do error, error:%v,input:%#v,rid:%s", err, param, lgc.rid)
		return 0, err
	}
	if !result.Result {
		blog.Errorf("getModuleHostCount http response error, err code:%d, err msg:%d, input:%+v, rid:%s", result.Code, result.ErrMsg, param, lgc.rid)
		return 0, err
	}

	rspV3MapData, ok := result.Data.([]interface{})
	if false == ok {
		blog.Errorf("assign error rspV3Map.data is not []interface{}, rspV3Map:%#v,input:%#v,rid:%s", result.Data, param, lgc.rid)
		return 0, nil
	}
	return len(rspV3MapData), nil
}

func (lgc *Logics) getSets(ctx context.Context, appID int64) ([]mapstr.MapStr, errors.CCError) {
	defErr := lgc.ccErr

	page := mapstr.MapStr{"start": 0, "limit": common.BKNoLimit}
	param := &params.SearchParams{Page: page, Condition: mapstr.MapStr{}}

	appIDStr := strconv.FormatInt(appID, 10)
	result, err := lgc.CoreAPI.TopoServer().Instance().SearchSet(ctx, lgc.ownerID, appIDStr, lgc.header, param)
	if err != nil {
		blog.Errorf("get sets   error:%v, cond:%+v,rid:%s", err, param, lgc.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get sets info http response error, err code:%d, err msg:%s, cond:%+v,rid:%s", result.Code, result.ErrMsg, param, lgc.rid)
		return nil, defErr.New(result.Code, result.ErrMsg)
	}
	return result.Data.Info, nil
}

func (lgc *Logics) getModules(ctx context.Context, appID int64) ([]mapstr.MapStr, errors.CCError) {

	fields := []string{common.BKModuleIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField, common.BKDefaultField}

	return lgc.getModulesByConds(ctx, strconv.FormatInt(appID, 10), nil, fields, "")
}

func (lgc *Logics) getModulesByConds(ctx context.Context, appIDStr string, conds map[string]interface{}, fields []string, sort string) ([]mapstr.MapStr, errors.CCError) {
	defErr := lgc.ccErr

	searchParams := mapstr.MapStr{}
	searchParams["fields"] = fields
	searchParams["condition"] = conds
	searchParams["page"] = mapstr.MapStr{"start": 0, "limit": common.BKNoLimit, "sort": sort}
	result, err := lgc.CoreAPI.TopoServer().OpenAPI().SearchModuleByApp(ctx, appIDStr, lgc.header, searchParams)
	if err != nil {
		blog.Errorf("getModules http do error,err:%v, cond:%+v,rid:%s", err, conds, lgc.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get app info http response error, err code:%d, err msg:%s, cond:%+v,rid:%s", result.Code, result.ErrMsg, searchParams, lgc.rid)
		return nil, defErr.New(result.Code, result.ErrMsg)
	}
	return result.Data.Info, nil
}

func (lgc *Logics) getAppInfo(ctx context.Context, appID int64) ([]mapstr.MapStr, errors.CCError) {
	defErr := lgc.ccErr

	page := mapstr.MapStr{"start": 0, "limit": 2}
	condition := mapstr.MapStr{common.BKAppIDField: appID}
	param := &params.SearchParams{Condition: condition, Page: page}
	result, err := lgc.CoreAPI.TopoServer().Instance().SearchApp(ctx, lgc.ownerID, lgc.header, param)
	if nil != err {
		blog.Errorf("get app info error:%v, cond:%+v,rid:%s", err, condition, lgc.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("get app info http response error, err code:%d, err msg:%s, cond:%+v,rid:%s", result.Code, result.ErrMsg, condition, lgc.rid)
		return nil, defErr.New(result.Code, result.ErrMsg)
	}

	return result.Data.Info, nil
}

func getRspV3DataInfo(logPrex string, result bool, code int, data interface{}, rid string) ([]interface{}, int) {

	if false == result {
		return nil, code
	}
	dataMap, ok := data.(map[string]interface{})
	if false == ok {
		blog.Errorf("%s rspV3 json get data.info error, body:%s  error:%#v,rid:%s", logPrex, dataMap, rid)
		return nil, common.CCErrCommJSONUnmarshalFailed
	}
	dataInfo, ok := dataMap["info"].([]interface{})
	if false == ok {
		blog.Errorf("%s rspV3 json get data.info error, body:%s  error:%#v", logPrex, dataMap, rid)
		return nil, common.CCErrCommJSONUnmarshalFailed
	}
	return dataInfo, 0
}

func (lgc *Logics) CheckAppTopoIsThreeLevel(ctx context.Context) (bool, error) {

	result, err := lgc.CoreAPI.TopoServer().Object().SelectModel(ctx, lgc.ownerID, lgc.header)
	if nil != err {
		blog.Errorf("CheckAppTopoIsThreeLevel http do error.err:%s,input:%#v,rid:%s", err.Error(), lgc.ownerID, lgc.rid)
		return false, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
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

func (lgc *Logics) GetAppListByOwnerIDAndUin(ctx context.Context, uin string) (appInfoArr []mapstr.MapStr, errCode int, err errors.CCError) {
	defErr := lgc.ccErr

	userLike := mapstr.MapStr{
		common.BKDBLIKE: fmt.Sprintf("^%s,|,%s,|,%s$|^%s$", uin, uin, uin, uin),
	}
	var condItemArr []mapstr.MapStr
	var conds mapstr.MapStr
	if lgc.ownerID != uin {
		condItemArr = append(condItemArr, mapstr.MapStr{common.CreatorField: userLike})
		condItemArr = append(condItemArr, mapstr.MapStr{common.BKMaintainersField: userLike})
		conds = mapstr.MapStr{common.BKDBOR: condItemArr}
	}
	s := &metadata.SearchParams{}
	if nil != conds {
		s.Condition = conds
	}

	result, err := lgc.CoreAPI.TopoServer().Instance().InstSearch(ctx, lgc.ownerID, common.BKInnerObjIDApp, lgc.header, s)
	if nil != err {
		blog.Errorf("GetAppListByOwnerIDAndUin http do error, error:%s,input:%#v,rid:%s", err.Error(), s, lgc.rid)
		return nil, common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !result.Result {
		blog.Errorf("GetAppListByOwnerIDAndUin http reply error, error code:%d, error message:%s,input:%#v,rid:%s", result.Code, result.ErrMsg, s, lgc.rid)
		return nil, common.CCErrCommHTTPDoRequestFailed, defErr.New(result.Code, result.ErrMsg)
	}

	if 1 <= len(result.Data.Info) {
		return nil, common.CCErrCommNotFound, defErr.Error(common.CCErrCommNotFound)
	}

	return result.Data.Info, 0, nil

}
