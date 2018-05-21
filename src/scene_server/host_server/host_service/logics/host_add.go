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
	"encoding/json"
	"errors"
	"fmt"
	"time"

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/util"
	scenecommon "configcenter/src/scene_server/common"
	"configcenter/src/scene_server/validator"
	sourceAuditAPI "configcenter/src/source_controller/api/auditlog"
	sourceAPI "configcenter/src/source_controller/api/object"
)

//AddHost, return error info
func AddHost(req *restful.Request, ownerID string, appID int, hostInfos map[int]map[string]interface{}, moduleID int, cc *api.APIResource) (error, []string, []string, []string) {
	forward := &sourceAPI.ForwardParam{Header: req.Request.Header}
	user := scenecommon.GetUserFromHeader(req)

	hostAddr := cc.HostCtrl()
	ObjAddr := cc.ObjCtrl()
	auditAddr := cc.AuditCtrl()
	addHostURL := hostAddr + "/host/v1/insts/"
	uHostURL := ObjAddr + "/object/v1/insts/host"

	language := util.GetActionLanguage(req)
	errHandle := cc.Error.CreateDefaultCCErrorIf(language)
	langHandle := cc.Lang.CreateDefaultCCLanguageIf(language)

	addParams := make(map[string]interface{})
	addParams[common.BKAppIDField] = appID
	addParams[common.BKModuleIDField] = []int{moduleID}
	addModulesURL := hostAddr + "/host/v1/meta/hosts/modules/"

	allHostList, err := GetHostInfoByConds(req, hostAddr, nil, langHandle)
	if nil != err {
		return errors.New(langHandle.Language("host_search_fail")), nil, nil, nil
	}

	//get asst field
	objCli := sourceAPI.NewClient("")
	objCli.SetAddress(ObjAddr)
	asst := map[string]interface{}{}
	asst[common.BKOwnerIDField] = ownerID
	asst[common.BKObjIDField] = common.BKInnerObjIDHost
	searchData, _ := json.Marshal(asst)
	objCli.SetAddress(ObjAddr)
	asstDes, err := objCli.SearchMetaObjectAsst(&sourceAPI.ForwardParam{Header: req.Request.Header}, searchData)
	if nil != err {
		return errors.New(langHandle.Language("host_search_fail")), nil, nil, nil
	}

	hostMap := convertHostInfo(allHostList)
	input := make(map[string]interface{}, 2)     //更新主机数据
	condInput := make(map[string]interface{}, 1) //更新主机条件
	var errMsg, succMsg, updateErrMsg []string   //新加错误， 成功，  更新失败
	iSubArea := common.BKDefaultDirSubArea

	defaultFields, err := getHostFields(forward, ownerID, ObjAddr)
	if nil != err {
		return errors.New("get host property failure"), nil, nil, nil
	}

	assObjectInt := scenecommon.NewAsstObjectInst(req, ownerID, ObjAddr, defaultFields, langHandle)
	err = assObjectInt.GetObjAsstObjectPrimaryKey()
	if nil != err {
		return fmt.Errorf("get host assocate object  property failure, error:%s", err.Error()), nil, nil, nil
	}
	rowErr, err := assObjectInt.InitInstFromData(hostInfos)
	if nil != err {
		return fmt.Errorf("get host assocate object instance data failure, error:%s", err.Error()), nil, nil, nil
	}

	ts := time.Now().UTC()
	//operator log
	var logConents []auditoplog.AuditLogExt
	hostLogFields, _ := GetHostLogFields(req, ownerID, ObjAddr)
	filterFields := []string{common.CreateTimeField}
	for index, host := range hostInfos {
		if nil == host {
			continue
		}

		innerIP, isOk := host[common.BKHostInnerIPField].(string)
		if isOk == false || "" == innerIP {
			errMsg = append(errMsg, langHandle.Languagef("host_import_innerip_empty", index))
			continue
		}

		valid := validator.NewValidMapWithKeyFields(common.BKDefaultOwnerID, common.BKInnerObjIDHost, ObjAddr, filterFields, forward, errHandle)
		key := fmt.Sprintf("%s-%v", innerIP, iSubArea)
		iHost, isOk := hostMap[key]

		//生产日志
		if isOk {
			if err, ok := rowErr[index]; true == ok {
				updateErrMsg = append(updateErrMsg, langHandle.Languagef("import_row_int_error_str", index, err.Error())) //fmt.Sprintf("%d行%s", index, err.Error()))
				continue
			}
			//delete(host, common.BKCloudIDField)
			delete(host, "import_from")
			delete(host, common.CreateTimeField)
			err := assObjectInt.SetObjAsstPropertyVal(host)
			if nil != err {
				blog.Error("host assocate property error %d %s", index, err)
				updateErrMsg = append(updateErrMsg, langHandle.Languagef("import_row_int_error_str", index, err.Error()))

				continue
			}

			hostID, _ := util.GetIntByInterface(iHost[common.BKHostIDField])
			fmt.Println("assObjectInt", hostID)

			_, err = valid.ValidMap(host, common.ValidUpdate, hostID)
			if nil != err {
				blog.Error("host valid error %v %v", index, err)
				updateErrMsg = append(updateErrMsg, langHandle.Languagef("import_row_int_error_str", index, err.Error()))
				fmt.Println("assObjectInt valid", langHandle.Languagef("import_row_int_error_str", index, err.Error()))

				continue
			}
			//update host asst attr
			err = scenecommon.UpdateInstAssociation(ObjAddr, req, hostID, ownerID, common.BKInnerObjIDHost, host) //hostAsstData, ownerID, host)
			if nil != err {
				blog.Error("update host asst attr error : %v", err)
				updateErrMsg = append(updateErrMsg, langHandle.Languagef("import_row_int_error_str", index, err.Error()))
				continue
			}
			//prepare the log
			strHostID := fmt.Sprintf("%d", hostID)
			logObj := NewHostLog(req, common.BKDefaultOwnerID, strHostID, hostAddr, ObjAddr, hostLogFields)

			condInput[common.BKHostIDField] = hostID
			input["condition"] = condInput
			input["data"] = host
			isSuccess, message, _ := GetHttpResult(req, uHostURL, common.HTTPUpdate, input)
			innerIP := host[common.BKHostInnerIPField].(string)
			if !isSuccess {
				blog.Error("host update error %v %v", index, message)
				updateErrMsg = append(updateErrMsg, langHandle.Languagef("host_import_update_fail", index, innerIP, message))
				continue
			}
			logContent, _ := logObj.GetHostLog(strHostID, false)
			logConents = append(logConents, auditoplog.AuditLogExt{ID: hostID, Content: logContent, ExtKey: innerIP})

		} else {
			if err, ok := rowErr[index]; true == ok {
				errMsg = append(errMsg, langHandle.Languagef("import_row_int_error_str", index, err.Error()))
				continue
			}
			err := assObjectInt.SetObjAsstPropertyVal(host)
			if nil != err {
				blog.Error("host assocate property error %v %v", index, err)
				updateErrMsg = append(updateErrMsg, langHandle.Languagef("import_row_int_error_str", index, err.Error()))
				continue
			}
			_, ok := host[common.BKCloudIDField]
			if false == ok {
				host[common.BKCloudIDField] = iSubArea
			}

			//补充未填写字段的默认值
			for _, field := range defaultFields {
				_, ok := host[field.PropertyID]
				if !ok {

					host[field.PropertyID] = ""
				}
			}

			_, err = valid.ValidMap(host, common.ValidCreate, 0)
			host[common.CreateTimeField] = ts

			if nil != err {
				errMsg = append(errMsg, langHandle.Languagef("import_row_int_error_str", index, err.Error()))
				continue
			}

			//prepare the log
			logObj := NewHostLog(req, common.BKDefaultOwnerID, "", hostAddr, ObjAddr, hostLogFields)

			isSuccess, message, retData := GetHttpResult(req, addHostURL, common.HTTPCreate, host)
			if !isSuccess {
				ip, _ := host["InnerIP"].(string)
				errMsg = append(errMsg, langHandle.Languagef("host_import_add_fail", index, ip, message))
				continue
			}

			retHost := retData.(map[string]interface{})
			hostID, _ := util.GetIntByInterface(retHost[common.BKHostIDField])

			//add host asst attr
			hostAsstData := scenecommon.ExtractDataFromAssociationField(int64(hostID), host, asstDes)
			err = scenecommon.CreateInstAssociation(ObjAddr, req, hostAsstData)
			if nil != err {
				blog.Error("add host asst attr error : %v", err)
				errMsg = append(errMsg, langHandle.Languagef("import_row_int_error_str", index, err.Error()))
				continue
			}

			addParams[common.BKHostIDField] = hostID
			innerIP := host[common.BKHostInnerIPField].(string)

			isSuccess, message, _ = GetHttpResult(req, addModulesURL, common.HTTPCreate, addParams)
			if !isSuccess {
				blog.Error("add hosthostconfig error, params:%v, error:%s", addParams, message)
				errMsg = append(errMsg, langHandle.Languagef("host_import_add_host_module", index, innerIP))
				continue
			}
			strHostID := fmt.Sprintf("%d", hostID)
			logContent, _ := logObj.GetHostLog(strHostID, false)

			logConents = append(logConents, auditoplog.AuditLogExt{ID: hostID, Content: logContent, ExtKey: innerIP})

		}

		succMsg = append(succMsg, fmt.Sprintf("%d", index))
	}

	if 0 < len(logConents) {
		logAPIClient := sourceAuditAPI.NewClient(auditAddr)
		_, err := logAPIClient.AuditHostsLog(logConents, "import host", ownerID, fmt.Sprintf("%d", appID), user, auditoplog.AuditOpTypeAdd)
		//addAuditLogs(req, logAdd, "新加主机", ownerID, appID, user, auditAddr)
		if nil != err {
			blog.Errorf("add audit log error %s", err.Error())
		}
	}

	if 0 < len(errMsg) || 0 < len(updateErrMsg) {
		return errors.New(langHandle.Language("host_import_err")), succMsg, updateErrMsg, errMsg
	}

	return nil, succMsg, updateErrMsg, errMsg
}

// EnterIP 将机器导入到制定模块或者空闲机器， 已经存在机器，不操作
func EnterIP(req *restful.Request, ownerID string, appID, moduleID int, ip string, cloudID int64, host map[string]interface{}, isIncrement bool, cc *api.APIResource) error {

	user := scenecommon.GetUserFromHeader(req)

	hostAddr := cc.HostCtrl()
	ObjAddr := cc.ObjCtrl()
	auditAddr := cc.AuditCtrl()

	language := util.GetActionLanguage(req)
	errHandle := cc.Error.CreateDefaultCCErrorIf(language)
	langHandle := cc.Lang.CreateDefaultCCLanguageIf(language)

	addHostURL := hostAddr + "/host/v1/insts/"

	addParams := make(map[string]interface{})
	addParams[common.BKAppIDField] = appID
	addParams[common.BKModuleIDField] = []int{moduleID}
	addModulesURL := hostAddr + "/host/v1/meta/hosts/modules/"

	isExist, err := IsExistPlat(req, ObjAddr, common.KvMap{common.BKCloudIDField: cloudID})
	if nil != err {
		return errors.New(langHandle.Language("host_search_fail")) // "查询主机信息失败")
	}
	if !isExist {
		return errors.New(langHandle.Language("plat_id_not_exist"))
	}
	conds := map[string]interface{}{
		common.BKHostInnerIPField: ip,
		common.BKCloudIDField:     cloudID,
	}
	hostList, err := GetHostInfoByConds(req, hostAddr, conds, langHandle)
	if nil != err {
		return errors.New(langHandle.Language("host_search_fail")) // "查询主机信息失败")
	}

	hostID := 0
	if len(hostList) == 0 {
		//host not exist, add host
		host[common.BKHostInnerIPField] = ip
		host[common.BKCloudIDField] = cloudID
		host["import_from"] = common.HostAddMethodAgent
		forward := &sourceAPI.ForwardParam{Header: req.Request.Header}
		defaultFields, hasErr := getHostFields(forward, ownerID, ObjAddr)
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
		_, hasErr = valid.ValidMap(host, "update", 0)

		if nil != hasErr {
			return hasErr
		}

		isSuccess, message, retData := GetHttpResult(req, addHostURL, common.HTTPCreate, host)
		if !isSuccess {
			return errors.New(langHandle.Languagef("host_agent_add_host_fail", message))
		}

		retHost := retData.(map[string]interface{})
		hostID, _ = util.GetIntByInterface(retHost[common.BKHostIDField])
	} else if false == isIncrement {
		//Not an additional relationship model
		return nil
	} else {
		hostMap, ok := hostList[0].(map[string]interface{})
		if false == ok {
			return errors.New(langHandle.Language("host_search_fail")) // "查询主机信息失败")
		}
		hostID, _ = util.GetIntByInterface(hostMap[common.BKHostIDField])
		if 0 == hostID {
			return errors.New(langHandle.Language("host_search_fail")) // "查询主机信息失败")
		}
		//func IsExistHostIDInApp(CC *api.APIResource, req *restful.Request, appID int, hostID int, defLang language.DefaultCCLanguageIf) (bool, error) {
		bl, hasErr := IsExistHostIDInApp(cc, req, appID, hostID, langHandle)
		if nil != hasErr {
			blog.Error("check host is exist in app error, params:{appid:%d, hostid:%s}, error:%s", appID, hostID, hasErr.Error())
			return errHandle.Errorf(common.CCErrHostNotINAPPFail, hostID)

		}
		if false == bl {
			blog.Error("Host does not belong to the current application; error, params:{appid:%d, hostid:%s}", appID, hostID)
			return errHandle.Errorf(common.CCErrHostNotINAPP, fmt.Sprintf("%d", hostID))
		}

	}

	//del host relation from default  module
	params := make(map[string]interface{})
	params[common.BKAppIDField] = appID
	params[common.BKHostIDField] = hostID
	delModulesURL := cc.HostCtrl() + "/host/v1/meta/hosts/defaultmodules"
	isSuccess, errMsg, _ := GetHttpResult(req, delModulesURL, common.HTTPDelete, params)
	if !isSuccess {
		blog.Error("remove hosthostconfig error, params:%v, error:%s", params, errMsg)
		return errHandle.Errorf(common.CCErrHostDELResourcePool, hostID)
	}

	addParams[common.BKHostIDField] = hostID

	isSuccess, message, _ := GetHttpResult(req, addModulesURL, common.HTTPCreate, addParams)
	if !isSuccess {
		blog.Error("enterip add hosthostconfig error, params:%v, error:%s", addParams, message)
		return errors.New(langHandle.Languagef("host_agent_add_host_module_fail", message))
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
