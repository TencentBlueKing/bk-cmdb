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
	"configcenter/src/common/util"
	"encoding/json"

	"github.com/emicklei/go-restful"
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
func GetSetIDByObjectCond(req *restful.Request, objURL string, objectCond []interface{}) []int {
	objectIDArr := make([]int, 0)
	conc := make(map[string]interface{})
	condition := make([]interface{}, 0)
	for _, i := range objectCond {
		condi := i.(map[string]interface{})
		field, ok := condi["field"].(string)
		if false == ok || field != common.BKInstIDField {
			continue
		}
		value, ok := condi["value"].(float64)
		if false == ok {
			continue
		}
		conc["field"] = common.BKInstParentStr
		conc["operator"] = condi["operator"]
		conc["value"] = condi["value"]
		objectIDArr = append(objectIDArr, int(value))
	}
	condition = append(condition, conc)
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
	}

}

//getObjectByParentID get object by parent id
func getObjectByParentID(req *restful.Request, valArr []int, objURL string) []int {
	objectIDArr := make([]int, 0)
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
		return objectIDArr
	}
	var sResult ObjectSResult
	err = json.Unmarshal([]byte(reply), &sResult)
	if err != nil {
		blog.Error("GetObjectIDByCond err :%v", err)
		return objectIDArr
	}
	if 0 == sResult.Data.Count {
		return objectIDArr
	}
	for _, j := range sResult.Data.Info {
		var cell64 float64
		var cell int
		cell, ok := j[common.BKInstIDField].(int)
		if false == ok {
			cell64 = j[common.BKInstIDField].(float64)
			cell = int(cell64)
		}
		objectIDArr = append(objectIDArr, cell)
	}
	return objectIDArr
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
