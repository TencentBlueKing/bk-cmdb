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

package phpapi

import (
	"encoding/json"
	"errors"
	"fmt"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	errorIfs "configcenter/src/common/errors"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/logics"
	"configcenter/src/scene_server/validator"
	"configcenter/src/source_controller/api/auditlog"
	sourceAPI "configcenter/src/source_controller/api/object"
)

var host *hostAction = &hostAction{}

type hostAction struct {
	base.BaseAction
}

func init() {
	host.CreateAction()

}

// helpers
func UpdateHostMain(req *restful.Request, hostCondition, data map[string]interface{}, appID int, hostCtrl, objCtrl, auditCtrl string, errIf errorIfs.CCErrorIf) (string, error) {
	blog.Debug("updateHostMain start")
	blog.Debug("hostCondition:%v", hostCondition)
	_, hostIDArr, err := GetHostMapByCond(req, hostCondition)

	blog.Debug("hostIDArr:%v", hostIDArr)
	if nil != err {
		return "", errors.New(fmt.Sprintf("GetHostIDByCond error:%v", err))
	}

	lenOfHostIDArr := len(hostIDArr)

	if lenOfHostIDArr != 1 {
		blog.Debug("GetHostMapByCond condition: %v", hostCondition)
		return "", errors.New(fmt.Sprintf("not find host info "))
	}

	language := util.GetActionLanguage(req)
	forward := &sourceAPI.ForwardParam{Header: req.Request.Header}
	valid := validator.NewValidMapWithKeyFields(common.BKDefaultOwnerID, common.BKInnerObjIDHost, objCtrl, []string{common.CreateTimeField, common.LastTimeField, common.BKChildStr, common.BKOwnerIDField}, forward, errIf.CreateDefaultCCErrorIf(language))
	ok, validErr := valid.ValidMap(data, common.ValidUpdate, hostIDArr[0])
	if false == ok && nil != validErr {
		blog.Error("updateHostMain error: %v", validErr)
		return "", validErr
	}

	configData, err := logics.GetConfigByCond(req, hostCtrl, map[string]interface{}{
		common.BKAppIDField:  []int{appID},
		common.BKHostIDField: []int{hostIDArr[0]},
	})

	if nil != err {
		return "", errors.New(fmt.Sprintf("GetConfigByCond error:%v", err))
	}

	lenOfConfigData := len(configData)

	if lenOfConfigData == 0 {
		blog.Errorf("not expected config lenth: appid:%d, hostid:%d", appID, hostIDArr[0])
		return "", errors.New(fmt.Sprintf("not expected config length: %d", lenOfConfigData))
	}

	hostID := configData[0][common.BKHostIDField]

	condition := make(map[string]interface{})
	condition[common.BKHostIDField] = hostID

	param := make(map[string]interface{})
	param["condition"] = condition
	param["data"] = data

	uURL := objCtrl + "/object/v1/insts/host"
	paramJson, err := json.Marshal(param)
	if nil != err {
		return "", errors.New(fmt.Sprintf("Marshal json error:%v", err))
	}
	strHostID := fmt.Sprintf("%d", hostID)
	logContent := logics.NewHostLog(req, common.BKDefaultOwnerID, strHostID, hostCtrl, objCtrl, nil)
	res, err := httpcli.ReqHttp(req, uURL, common.HTTPUpdate, []byte(paramJson))
	if nil == err {
		//操作成功，新加操作日志日志
		resJs, err := simplejson.NewJson([]byte(res))
		if err == nil {
			bl, _ := resJs.Get("result").Bool()
			if bl {
				user := util.GetActionUser(req)
				opClient := auditlog.NewClient(auditCtrl)
				content, _ := logContent.GetHostLog(strHostID, false)
				//(id interface{}, Content interface{}, OpDesc string, InnerIP, ownerID, appID, user string, OpType auditoplog.AuditOpType)
				opClient.AuditHostLog(hostID, content, "修改主机", logContent.GetInnerIP(), common.BKDefaultOwnerID, fmt.Sprintf("%d", appID), user, auditoplog.AuditOpTypeModify)

			}
		}
	}
	return res, err
}

func AddHost(req *restful.Request, data map[string]interface{}, objURL string) (int64, error) {
	return addObj(req, data, common.BKInnerObjIDHost, objURL)
}

func AddModuleHostConfig(req *restful.Request, hostID, appID int64, moduleIDs []int64, hostCtrl string) error {
	data := common.KvMap{
		common.BKAppIDField:    appID,
		common.BKHostIDField:   hostID,
		common.BKModuleIDField: moduleIDs,
	}
	blog.Debug("addModuleHostConfig start, data: %v", data)

	resMap := make(map[string]interface{})
	inputJson, _ := json.Marshal(data)
	addModulesURL := hostCtrl + "/host/v1/meta/hosts/modules"
	blog.Info("addModuleHostconfig params", string(inputJson))
	res, err := httpcli.ReqHttp(req, addModulesURL, common.HTTPCreate, []byte(inputJson))
	if nil != err {
		blog.Errorf("http do error:url:%s, error:%s", addModulesURL, err.Error())
		return err
	}

	err = json.Unmarshal([]byte(res), &resMap)
	blog.Info("addModuleHostConfig reply:%s", res)
	if nil != err {
		blog.Errorf("http do error:url:%s, respone not json params:%s, reply:%s", addModulesURL, string(inputJson), res)
		return err
	}
	if !resMap["result"].(bool) {
		return errors.New(resMap[common.HTTPBKAPIErrorMessage].(string))
	}
	blog.Debug("addModuleHostConfig success, res: %v", resMap)
	return nil
}

func addObj(req *restful.Request, data map[string]interface{}, objType, objURL string) (int64, error) {
	resMap := make(map[string]interface{})

	url := objURL + "/object/v1/insts/" + objType
	inputJson, _ := json.Marshal(data)
	res, err := httpcli.ReqHttp(req, url, common.HTTPCreate, []byte(inputJson))
	if nil != err {
		return 0, err
	}

	err = json.Unmarshal([]byte(res), &resMap)
	if nil != err {
		return 0, err
	}

	if !resMap["result"].(bool) {
		return 0, errors.New(resMap[common.HTTPBKAPIErrorMessage].(string))
	}

	blog.Debug("add object result : %v", resMap)
	dataMap, ok := resMap["data"].(map[string]interface{})
	if false == ok {
		return 0, fmt.Errorf("add host reply error; reply:%s", res)
	}

	objID, err := util.GetInt64ByInterface(dataMap[common.BKHostIDField])
	if nil != err {
		return 0, fmt.Errorf("add host reply error, not found host id")
	}
	return objID, nil
}

//search host helpers

func SetHostData(req *restful.Request, moduleHostConfig []map[string]int, hostMap map[int]interface{}) ([]interface{}, error) {

	//total data
	hostData := make([]interface{}, 0)

	appIDArr := make([]int, 0)
	setIDArr := make([]int, 0)
	moduleIDArr := make([]int, 0)

	for _, config := range moduleHostConfig {
		setIDArr = append(setIDArr, config[common.BKSetIDField])
		moduleIDArr = append(moduleIDArr, config[common.BKModuleIDField])
		appIDArr = append(appIDArr, config[common.BKAppIDField])
	}

	moduleMap, err := logics.GetModuleMapByCond(req, "", host.CC.ObjCtrl(), map[string]interface{}{
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: moduleIDArr,
		},
	})
	if err != nil {
		return hostData, err
	}

	setMap, err := logics.GetSetMapByCond(req, "", host.CC.ObjCtrl(), map[string]interface{}{
		common.BKSetIDField: map[string]interface{}{
			common.BKDBIN: setIDArr,
		},
	})
	if err != nil {
		return hostData, err
	}

	blog.Debug("GetAppMapByCond , appIDArr:%v", appIDArr)
	appMap, err := logics.GetAppMapByCond(req, "", host.CC.ObjCtrl(), map[string]interface{}{
		common.BKAppIDField: map[string]interface{}{
			common.BKDBIN: appIDArr,
		},
	})

	if err != nil {
		return hostData, err
	}
	for _, config := range moduleHostConfig {
		host, hasHost := hostMap[config[common.BKHostIDField]].(map[string]interface{})
		if !hasHost {
			blog.Errorf("hostMap has not hostID: %d", config[common.BKHostIDField])
			continue
		}

		module := moduleMap[config[common.BKModuleIDField]].(map[string]interface{})
		set := setMap[config[common.BKSetIDField]].(map[string]interface{})
		app := appMap[config[common.BKAppIDField]].(map[string]interface{})

		hostStr, _ := json.Marshal(host)
		hostNew := make(map[string]interface{})
		json.Unmarshal(hostStr, &hostNew)

		hostNew[common.BKModuleIDField] = module[common.BKModuleIDField]
		hostNew[common.BKModuleNameField] = module[common.BKModuleNameField]
		hostNew[common.BKModuleTypeField] = module[common.BKModuleTypeField]
		hostNew[common.BKSetIDField] = set[common.BKSetIDField]
		hostNew[common.BKSetNameField] = set[common.BKSetNameField]
		hostNew[common.BKAppIDField] = app[common.BKAppIDField]
		hostNew[common.BKAppNameField] = app[common.BKAppNameField]
		hostNew[common.BKOwnerIDField] = app[common.BKOwnerIDField]
		hostNew[common.BKOperatorField] = module[common.BKOperatorField]
		hostNew[common.BKBakOperatorField] = module[common.BKBakOperatorField]

		hostData = append(hostData, hostNew)
	}
	return hostData, nil
}
