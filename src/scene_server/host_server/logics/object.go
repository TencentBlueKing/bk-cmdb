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

	meta "configcenter/src/common/metadata"
	"configcenter/src/scene_server/host_server/service"
    "configcenter/src/common"
    "configcenter/src/common/blog"
    "configcenter/src/common/util"
    "encoding/json"
    "github.com/emicklei/go-restful"
)

// get the object attributes
func (lgc *Logics) GetObjectAttributes(ownerID, objID string, pheader http.Header, page meta.BasePage) ([]meta.Attribute, error) {
	opt := service.NewOperation().WithOwnerID(ownerID).WithPage(page).WithObjID(objID).Data()
	result, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), pheader, opt)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get object attribute failed, err: %v, %v", err, result.ErrMsg)
	}

	return result.Data, nil
}

func (lgc *Logics) GetTopoIDByName(pheader http.Header, c *meta.HostToAppModule) (int64, int64, int64, error) {
    if "" == c.AppName || "" == c.SetName || "" == c.ModuleName {
        return 0, 0, 0, nil
    }
    
    appInfo, appErr := lgc.GetSingleApp(pheader, common.KvMap{common.BKAppNameField: c.AppName, common.BKOwnerIDField: c.OwnerID})
    if nil != appErr {
        blog.Errorf("getTopoIDByName get app info error; %s", appErr.Error())
        return 0, 0, 0, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader)).Error(common.CCErrCommHTTPDoRequestFailed)
    }
    
    appID, err := appInfo.Int64(common.BKAppIDField)
    if err != nil {
        return 0, 0, 0, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader)).Error(common.CCErrCommParamsInvalid)
    }

    appIDConf := map[string]interface{}{
        "field":    common.BKAppIDField,
        "operator": common.BKDBEQ,
        "value":    appID,
    }
    setNameCond := map[string]interface{}{
        "field":    common.BKSetNameField,
        "operator": common.BKDBEQ,
        "value":    c.SetName,
    }
    var setCond []interface{}
    setCond = append(setCond, appIDConf, setNameCond)

    setIDs, setErr := lgc.GetSetIDByCond(pheader, setCond)
    if nil != setErr {
        blog.Errorf("getTopoIDByName get app info error; %s", setErr.Error())
        return 0, 0, 0, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader)).Error(common.CCErrCommHTTPDoRequestFailed)
    }
    if 0 == len(setIDs) || 0 >= setIDs[0] {
        blog.Info("getTopoIDByName get set info not found; applicationName: %s, setName: %s", c.AppName, c.SetName)
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
        "value":    c.ModuleName,
    }
    var moduleCond []interface{}
    moduleCond = append(moduleCond, appIDConf, setIDConds, moduleNameCond)
    moduleIDs, moduleErr := lgc.GetModuleIDByCond(pheader, moduleCond)
    if nil != moduleErr {
        blog.Errorf("getTopoIDByName get app info error; %s", setErr.Error())
        return 0, 0, 0, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader)).Error(common.CCErrCommHTTPDoRequestFailed)
    }
    if 0 == len(moduleIDs) || 0 >= moduleIDs[0] {
        blog.Info("getTopoIDByName get module info not found; applicationName: %s, setName: %s, moduleName: %s", c.AppName, c.SetName, c.ModuleName)
        return 0, 0, 0, nil
    }
    moduleID := moduleIDs[0]

    return appID, setID, moduleID, nil
}

func (lgc *Logics) GetSetIDByObjectCond(pheader http.Header, appID int64, objectCond []interface{}) ([]int64, error) {
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
        sSetIDArr, err := lgc.GetSetIDByCond(pheader, condition)
        if err != nil {
            return nil, err
        }

        if 0 != len(sSetIDArr) {
            return sSetIDArr, nil
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

func (lgc *Logics) func getObjectByParentID(pheader http.Header, valArr []int64) []int {
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