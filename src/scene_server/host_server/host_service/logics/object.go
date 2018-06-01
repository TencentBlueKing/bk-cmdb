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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpcli "configcenter/src/common/http/httpclient"
	parse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"encoding/json"

	"github.com/emicklei/go-restful"
	"github.com/tidwall/gjson"
)

type ObjectSResult struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    ObjectData  `json:"data"`
}
type ObjectData struct {
	Count int                      `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}

//GetSetIDByObjectCond get set id by object condition
func GetSetIDByObjectCond(req *restful.Request, objURL string, appID int, objectCond []interface{}) []int {
	objectIDArr := make([]int, 0)
	conc := make(map[string]interface{})
	appIDCond := make(map[string]interface{})
	condition := make([]interface{}, 0)

	appIDCond["field"] = common.BKAppIDField
	appIDCond["operator"] = common.BKDBEQ
	appIDCond["value"] = appID

	for _, i := range objectCond {
		condi := i.(map[string]interface{})
		field, ok := condi["field"].(string)
		if false == ok || field != common.BKInstIDField {
			continue
		}
		value, err := util.GetIntByInterface(condi["value"])
		if nil != err {
			continue
		}
		conc["field"] = common.BKInstParentStr
		conc["operator"] = condi["operator"]
		conc["value"] = condi["value"]
		objectIDArr = append(objectIDArr, value)
	}
	condition = append(condition, conc)
	condition = append(condition, appIDCond)
	for {
		sSetIDArr, _ := GetSetIDByCond(req, objURL, condition)

		if 0 != len(sSetIDArr) {
			return sSetIDArr
		}

		sObjectIDArr := getObjectByParentID(req, objectIDArr, objURL)
		objectIDArr = sObjectIDArr
		if 0 == len(sObjectIDArr) {
			return []int{}
		}
		conc = make(map[string]interface{})
		condition = make([]interface{}, 0)
		conc["field"] = common.BKInstParentStr
		conc["operator"] = common.BKDBIN
		conc["value"] = sObjectIDArr
		condition = append(condition, conc)
		condition = append(condition, appIDCond)
	}

}

//getHostIDByInstID get host Id by inst Id
func GetHostIDByInstID(req *restful.Request, asstObjId string, objURL string, instIDArr []int) []int {

	//search the association insts
	hostIDArr := make([]int, 0)
	sURL := objURL + "/object/v1/insts/" + common.BKTableNameInstAsst + "/search"
	input := make(map[string]interface{})
	conditon := make(map[string]interface{})
	bkAsstInstId := make(map[string]interface{})
	bkAsstInstId[common.BKDBIN] = instIDArr
	conditon[common.BKAsstInstIDField] = bkAsstInstId
	conditon[common.BKAsstObjIDField] = asstObjId
	conditon[common.BKObjIDField] = common.BKInnerObjIDHost
	input["condition"] = conditon

	inputJSON, _ := json.Marshal(input)
	blog.Info("get host id by inst id url: %s", sURL)
	blog.Info("get host id by inst id content: %s", inputJSON)
	instRes, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, []byte(inputJSON))
	if nil != err {
		blog.Errorf("failed to search the inst association, condition is %s ,error is %s", string(inputJSON), err.Error())
		return hostIDArr
	}
	blog.Info("get host id by inst id return: %s", instRes)
	gjson.Get(instRes, "data.info.#."+common.BKInstIDField).ForEach(func(key, value gjson.Result) bool {

		hostIDArr = append(hostIDArr, int(value.Int()))
		return true
	})

	return util.IntArrayUnique(hostIDArr)

}

//getObjectInstByCond get object inst id by condtion
func GetObjectInstByCond(req *restful.Request, objID string, objURL string, cond []interface{}) []int {
	var url string
	var outField string
	instIDArr := make([]int, 0)
	condition := make(map[string]interface{})
	condc := make(map[string]interface{})
	parse.ParseCommonParams(cond, condc)
	if objID == common.BKInnerObjIDPlat {
		url = objURL + "/object/v1/insts/plat/search"
		outField = common.BKCloudIDField
	} else {
		condc[common.BKObjIDField] = objID
		url = objURL + "/object/v1/insts/object/search"
		outField = common.BKInstIDField
	}

	condition["condition"] = condc
	bodyContent, _ := json.Marshal(condition)
	blog.Info("GetObjectInstByCond url :%s", url)
	blog.Info("GetObjectInstByCond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetObjectInstByCond return :%s", string(reply))
	if err != nil {
		blog.Error("GetObjectInstByCond err :%v", err)
		return instIDArr
	}
	var sResult ObjectSResult
	err = json.Unmarshal([]byte(reply), &sResult)
	if err != nil {
		blog.Error("GetObjectInstByCond err :%v", err)
		return instIDArr
	}
	if 0 == sResult.Data.Count {
		return instIDArr
	}
	for _, j := range sResult.Data.Info {
		cell, err := util.GetIntByInterface(j[outField])
		if nil != err {
			continue
		}
		instIDArr = append(instIDArr, cell)
	}
	return instIDArr
}

//getObjectByParentID get object by parent id
func getObjectByParentID(req *restful.Request, valArr []int, objURL string) []int {
	instIDArr := make([]int, 0)
	condCell := make(map[string]interface{})
	condition := make(map[string]interface{})
	sCond := make(map[string]interface{})
	condCell[common.BKDBIN] = valArr
	sCond[common.BKInstParentStr] = condCell
	condition["condition"] = sCond
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/object/search"
	blog.Info("GetObjectIDByCond url :%s", url)
	blog.Info("GetObjectIDByCond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetObjectIDByCond return :%s", string(reply))
	if err != nil {
		blog.Error("GetsetIDByCond err :%v", err)
		return instIDArr
	}
	var sResult ObjectSResult
	err = json.Unmarshal([]byte(reply), &sResult)
	if err != nil {
		blog.Error("GetObjectIDByCond err :%v", err)
		return instIDArr
	}
	if 0 == sResult.Data.Count {
		return instIDArr
	}
	for _, j := range sResult.Data.Info {
		cell, err := util.GetIntByInterface(j[common.BKInstIDField])
		if nil != err {
			continue
		}
		instIDArr = append(instIDArr, cell)
	}
	return instIDArr
}

//GetTopoIDByName  get topo id by name
func GetTopoIDByName(req *restful.Request, ownerID, appName, setName, moduleName, objURL string, defErr errors.DefaultCCErrorIf) (int, int, int, error) {
	if "" == appName || "" == setName || "" == moduleName {
		return 0, 0, 0, nil
	}
	appInfo, appErr := GetSingleApp(req, objURL, common.KvMap{common.BKAppNameField: appName, common.BKOwnerIDField: ownerID})
	if nil != appErr {
		blog.Errorf("getTopoIDByName get app info error; %s", appErr.Error())
		return 0, 0, 0, defErr.Errorf(common.CCErrCommHTTPDoRequestFailed)
	}

	appID, _ := util.GetIntByInterface(appInfo[common.BKAppIDField])
	if 0 >= appID {
		blog.Info("getTopoIDByName get app info not found; applicationName: %s", appName)
		return 0, 0, 0, nil
	}
	appIDConf := map[string]interface{}{
		"field":    common.BKAppIDField,
		"operator": common.BKDBEQ,
		"value":    appID,
	}
	setNameCond := map[string]interface{}{
		"field":    common.BKSetNameField,
		"operator": common.BKDBEQ,
		"value":    setName,
	}
	var setCond []interface{}
	setCond = append(setCond, appIDConf, setNameCond)

	setIDs, setErr := GetSetIDByCond(req, objURL, setCond)
	if nil != setErr {
		blog.Errorf("getTopoIDByName get app info error; %s", setErr.Error())
		return 0, 0, 0, defErr.Errorf(common.CCErrCommHTTPDoRequestFailed)
	}
	if 0 == len(setIDs) || 0 >= setIDs[0] {
		blog.Info("getTopoIDByName get set info not found; applicationName: %s, setName: %s", appName, setName)
		return 0, 0, 0, nil
	}
	setID := setIDs[0]
	setIDConds := map[string]interface{}{
		"field":    common.BKSetIDField,
		"operator": common.BKDBEQ,
		"value":    setID,
	}
	moduleNameCond := map[string]interface{}{
		"field":    common.BKModuleNameField,
		"operator": common.BKDBEQ,
		"value":    moduleName,
	}
	var moduleCond []interface{}
	moduleCond = append(moduleCond, appIDConf, setIDConds, moduleNameCond)
	moduleIDs, moduleErr := GetModuleIDByCond(req, objURL, moduleCond)

	if nil != moduleErr {
		blog.Errorf("getTopoIDByName get app info error; %s", setErr.Error())
		return 0, 0, 0, defErr.Errorf(common.CCErrCommHTTPDoRequestFailed)
	}
	if 0 == len(moduleIDs) || 0 >= moduleIDs[0] {
		blog.Info("getTopoIDByName get module info not found; applicationName: %s, setName: %s, moduleName: %s", appName, setName, moduleName)
		return 0, 0, 0, nil
	}
	moduleID := moduleIDs[0]

	return appID, setID, moduleID, nil
}
