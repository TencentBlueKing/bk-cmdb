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
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"encoding/json"
	"errors"
	"fmt"

	"configcenter/src/common/auditoplog"
	errorHandle "configcenter/src/common/errors"
	sencecommon "configcenter/src/scene_server/common"
	"configcenter/src/scene_server/validator"
	sourceAuditAPI "configcenter/src/source_controller/api/auditlog"
	sourceAPI "configcenter/src/source_controller/api/object"

	"time"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

//GetModuleByModuleID  get module by module id
func GetModuleByModuleID(req *restful.Request, appID int, moduleID int, hostAddr string) ([]interface{}, error) {
	//moduleURL := "http://" + cc.ObjCtrl + "/object/v1/insts/module/search"
	URL := hostAddr + "/object/v1/insts/module/search"
	params := make(map[string]interface{})

	conditon := make(map[string]interface{})
	conditon[common.BKAppIDField] = appID
	conditon[common.BKModuleIDField] = moduleID
	params["condition"] = conditon
	params["sort"] = common.BKModuleIDField
	params["start"] = 0
	params["limit"] = 1
	params["fields"] = common.BKModuleIDField
	isSuccess, errMsg, data := GetHttpResult(req, URL, common.HTTPSelectPost, params)
	if !isSuccess {
		blog.Error("get idle module error, params:%v, error:%s", params, errMsg)
		return nil, errors.New(errMsg)
	}
	dataStrArry := data.(map[string]interface{})
	dataInfo, ok := dataStrArry["info"].([]interface{})
	if !ok {
		blog.Error("get idle module error, params:%v, error:%s", params, errMsg)
		return nil, errors.New(errMsg)
	}

	return dataInfo, nil
}

//GetSingleModuleID  get single module id
func GetSingleModuleID(req *restful.Request, conds interface{}, hostAddr string) (int, error) {
	//moduleURL := "http://" + cc.ObjCtrl + "/object/v1/insts/module/search"
	url := hostAddr + "/object/v1/insts/module/search"
	params := make(map[string]interface{})

	params["condition"] = conds
	params["sort"] = common.BKModuleIDField
	params["start"] = 0
	params["limit"] = 1
	params["fields"] = common.BKModuleIDField
	isSuccess, errMsg, data := GetHttpResult(req, url, common.HTTPSelectPost, params)
	if !isSuccess {
		blog.Error("get idle module error, params:%v, error:%s", params, errMsg)
		return 0, errors.New(errMsg)
	}
	dataInterface := data.(map[string]interface{})
	info := dataInterface["info"].([]interface{})
	if 1 != len(info) {
		blog.Error("not find module error, params:%v, error:%s", params, errMsg)
		return 0, errors.New("获取集群，返回数据格式错误")
	}
	row := info[0].(map[string]interface{})
	moduleID, _ := util.GetIntByInterface(row[common.BKModuleIDField])

	if 0 == moduleID {
		blog.Error("not find module error, params:%v, error:%s", params, errMsg)
		return 0, errors.New("获取集群信息失败")
	}

	return moduleID, nil
}

//AddHost, return error info
func AddHost(req *restful.Request, ownerID string, appID int, hostInfos map[int]map[string]interface{}, moduleID int, hostAddr, ObjAddr, auditAddr string, errHandle errorHandle.DefaultCCErrorIf) (error, []string, []string, []string) {

	user := sencecommon.GetUserFromHeader(req)

	addHostURL := hostAddr + "/host/v1/insts/"
	uHostURL := ObjAddr + "/object/v1/insts/host"

	addParams := make(map[string]interface{})
	addParams[common.BKAppIDField] = appID
	addParams[common.BKModuleIDField] = []int{moduleID}
	addModulesURL := hostAddr + "/host/v1/meta/hosts/modules/"

	allHostList, err := GetHostInfoByConds(req, hostAddr, nil)
	if nil != err {
		return errors.New("查询主机信息失败"), nil, nil, nil
	}

	hostMap := convertHostInfo(allHostList)
	input := make(map[string]interface{}, 2)     //更新主机数据
	condInput := make(map[string]interface{}, 1) //更新主机条件
	var errMsg, succMsg, updateErrMsg []string   //新加错误， 成功，  更新失败
	iSubArea := common.BKDefaultDirSubArea

	defaultFields := getHostFields(ownerID, ObjAddr)
	ts := time.Now().UTC()
	//operator log
	var logConents []auditoplog.AuditLogExt
	hostLogFields, _ := GetHostLogFields(req, ownerID, ObjAddr)
	for index, host := range hostInfos {
		if nil == host {
			continue
		}

		innerIP, ok := host[common.BKHostInnerIPField].(string)
		if ok == false || "" == innerIP {
			errMsg = append(errMsg, fmt.Sprintf("%d行内网ip为空", index))
			continue
		}
		notExistFields := []string{} //没有赋值的key，不需要校验
		for key, value := range defaultFields {
			_, ok := host[key]
			if ok {
				//已经存在，
				continue
			}
			require, _ := util.GetIntByInterface(value["require"])
			if require == common.BKTrue {
				errMsg = append(errMsg, fmt.Sprintf("%d行内网%s必填", index, key))
				continue
			}
			notExistFields = append(notExistFields, key)
		}
		blog.Infof("no validate fields %v", notExistFields)

		valid := validator.NewValidMapWithKeyFileds(common.BKDefaultOwnerID, common.BKInnerObjIDHost, ObjAddr, notExistFields, errHandle)

		key := fmt.Sprintf("%s-%v", innerIP, iSubArea)
		iHost, ok := hostMap[key]
		//生产日志
		if ok {
			delete(host, common.BKCloudIDField)
			delete(host, "import_from")
			delete(host, common.CreateTimeField)
			hostInfo := iHost.(map[string]interface{})

			hostID, _ := util.GetIntByInterface(hostInfo[common.BKHostIDField])
			_, err = valid.ValidMap(host, common.ValidUpdate, hostID)
			if nil != err {
				updateErrMsg = append(updateErrMsg, fmt.Sprintf("%d行%v", index, err))
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
				ret := fmt.Sprintf("%s更新失败%v;", innerIP, message)
				updateErrMsg = append(updateErrMsg, fmt.Sprintf("%d行%v", index, ret))
				continue
			}
			logContent, _ := logObj.GetHostLog(strHostID, false)
			logConents = append(logConents, auditoplog.AuditLogExt{ID: hostID, Content: logContent, ExtKey: innerIP})

		} else {
			host[common.BKCloudIDField] = iSubArea
			host[common.CreateTimeField] = ts
			//补充未填写字段的默认值
			for key, val := range defaultFields {
				_, ok := host[key]
				if !ok {
					host[key] = val["default"]
				}
			}
			_, err := valid.ValidMap(host, common.ValidCreate, 0)

			if nil != err {
				errMsg = append(errMsg, fmt.Sprintf("%d行%v", index, err))
				continue
			}

			//prepare the log
			logObj := NewHostLog(req, common.BKDefaultOwnerID, "", hostAddr, ObjAddr, hostLogFields)

			isSuccess, message, retData := GetHttpResult(req, addHostURL, common.HTTPCreate, host)
			if !isSuccess {
				ret := fmt.Sprintf("%s新加失败%s;", host["InnerIP"].(string), message)
				errMsg = append(errMsg, fmt.Sprintf("%d行%v", index, ret))
				continue
			}

			retHost := retData.(map[string]interface{})
			hostID, _ := util.GetIntByInterface(retHost[common.BKHostIDField])
			addParams[common.BKHostIDField] = hostID
			innerIP := host[common.BKHostInnerIPField].(string)

			isSuccess, message, _ = GetHttpResult(req, addModulesURL, common.HTTPCreate, addParams)
			if !isSuccess {
				blog.Error("add hosthostconfig error, params:%v, error:%s", addParams, message)
				errMsg = append(errMsg, fmt.Sprintf("%d行%v", index, innerIP))
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
		_, err := logAPIClient.AuditHostsLog(logConents, "导入主机", ownerID, fmt.Sprintf("%d", appID), user, auditoplog.AuditOpTypeAdd)
		//addAuditLogs(req, logAdd, "新加主机", ownerID, appID, user, auditAddr)
		if nil != err {
			blog.Errorf("add audit log error %s", err.Error())
		}
	}

	if 0 < len(errMsg) || 0 < len(updateErrMsg) {
		return errors.New("导入主机出现错误"), succMsg, updateErrMsg, errMsg
	}

	return nil, succMsg, updateErrMsg, errMsg
}

//EnterIP 将机器导入到制定模块或者空闲机器， 已经存在机器，不操作
func EnterIP(req *restful.Request, ownerID string, appID, moduleID int, IP, osType, hostname, appName, setName, moduleName, hostAddr, ObjAddr, auditAddr string, errHandle errorHandle.DefaultCCErrorIf) error {

	user := sencecommon.GetUserFromHeader(req)

	addHostURL := hostAddr + "/host/v1/insts/"

	addParams := make(map[string]interface{})
	addParams[common.BKAppIDField] = appID
	addParams[common.BKModuleIDField] = []int{moduleID}
	addModulesURL := hostAddr + "/host/v1/meta/hosts/modules/"

	conds := map[string]interface{}{
		common.BKHostInnerIPField: IP,
		common.BKCloudIDField:     common.BKDefaultDirSubArea,
	}
	hostList, err := GetHostInfoByConds(req, hostAddr, conds)
	if nil != err {
		return errors.New("查询主机信息失败")
	}
	if len(hostList) > 0 {
		return nil
	}

	host := make(map[string]interface{})
	host[common.BKHostInnerIPField] = IP
	host[common.BKOSTypeField] = osType

	host["import_from"] = common.HostAddMethodAgent
	host[common.BKCloudIDField] = common.BKDefaultDirSubArea
	defaultFields := getHostFields(ownerID, ObjAddr)
	//补充未填写字段的默认值
	for key, val := range defaultFields {
		_, ok := host[key]
		if !ok {

			host[key] = val[common.BKDefaultField]
		}
	}

	isSuccess, message, retData := GetHttpResult(req, addHostURL, common.HTTPCreate, host)
	if !isSuccess {
		return errors.New(fmt.Sprintf("add host to cmdb error,error:%s", message))
	}

	retHost := retData.(map[string]interface{})
	hostID, _ := util.GetIntByInterface(retHost[common.BKHostIDField])
	addParams[common.BKHostIDField] = hostID

	isSuccess, message, _ = GetHttpResult(req, addModulesURL, common.HTTPCreate, addParams)
	if !isSuccess {
		blog.Error("enterip add hosthostconfig error, params:%v, error:%s", addParams, message)
		return errors.New(fmt.Sprintf("add hosthostconfig error,error:%s", message))
	}

	//prepare the log
	hostLogFields, _ := GetHostLogFields(req, ownerID, ObjAddr)
	logObj := NewHostLog(req, common.BKDefaultOwnerID, "", hostAddr, ObjAddr, hostLogFields)
	content, _ := logObj.GetHostLog(fmt.Sprintf("%d", hostID), false)
	logAPIClient := sourceAuditAPI.NewClient(auditAddr)
	logAPIClient.AuditHostLog(hostID, content, "enter IP HOST", IP, ownerID, fmt.Sprintf("%d", appID), user, auditoplog.AuditOpTypeAdd)
	logClient, err := NewHostModuleConfigLog(req, nil, hostAddr, ObjAddr, auditAddr)
	logClient.SetHostID([]int{hostID})
	logClient.SetDescPrefix("enter IP ")
	logClient.SaveLog(fmt.Sprintf("%d", appID), user)
	return nil

}

//getHostFields 获取主所有字段和默认值
func getHostFields(ownerID, ObjAddr string) map[string]map[string]interface{} {
	return GetObjectFields(ownerID, common.BKInnerObjIDHost, ObjAddr)
}

//GetObjectFields get object fields
func GetObjectFields(ownerID, objID, ObjAddr string) map[string]map[string]interface{} {
	data := make(map[string]interface{})
	data[common.BKOwnerIDField] = ownerID
	data[common.BKObjIDField] = objID
	info, _ := json.Marshal(data)
	client := sourceAPI.NewClient(ObjAddr)
	result, _ := client.SearchMetaObjectAtt([]byte(info))
	fields := make(map[string]map[string]interface{})
	for _, j := range result {
		propertyID := j.PropertyID
		fieldType := j.PropertyType
		switch fieldType {
		case common.FiledTypeSingleChar:
			fields[propertyID] = common.KvMap{"default": "", "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired}
		case common.FiledTypeLongChar:
			fields[propertyID] = common.KvMap{"default": "", "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired} //""
		case common.FiledTypeInt:
			fields[propertyID] = common.KvMap{"default": nil, "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired} //0
		case common.FiledTypeEnum:
			fields[propertyID] = common.KvMap{"default": nil, "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired}
		case common.FiledTypeDate:
			fields[propertyID] = common.KvMap{"default": nil, "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired}
		case common.FiledTypeTime:
			fields[propertyID] = common.KvMap{"default": nil, "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired}
		case common.FiledTypeUser:
			fields[propertyID] = common.KvMap{"default": nil, "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired}
		default:
			fields[propertyID] = common.KvMap{"default": nil, "name": j.PropertyName, "type": j.PropertyType, "require": j.IsRequired}
			continue
		}

	}
	return fields
}

//convertHostInfo convert host info，InnerIP+SubArea key map[string]interface
func convertHostInfo(hosts []interface{}) map[string]interface{} {
	var hostMap map[string]interface{} = make(map[string]interface{})
	for _, host := range hosts {
		h := host.(map[string]interface{})

		key := fmt.Sprintf("%v-%v", h[common.BKHostInnerIPField], h[common.BKCloudIDField])
		hostMap[key] = h
	}
	return hostMap
}

func GetHostInfoByConds(req *restful.Request, hostURL string, conds map[string]interface{}) ([]interface{}, error) {
	hostURL = hostURL + "/host/v1/hosts/search"
	getParams := make(map[string]interface{})
	getParams["fields"] = nil
	getParams["condition"] = conds
	getParams["start"] = 0
	getParams["limit"] = common.BKNoLimit
	getParams["sort"] = common.BKHostIDField
	blog.Info("get host info by conds url:%s", hostURL)
	blog.Info("get host info by conds params:%v", getParams)
	isSucess, message, iRetData := GetHttpResult(req, hostURL, common.HTTPSelectPost, getParams)
	blog.Info("get host info by conds return:%v", iRetData)
	if !isSucess {
		return nil, errors.New("获取主机信息失败;" + message)
	}
	if nil == iRetData {
		return nil, nil
	}
	retData := iRetData.(map[string]interface{})
	data, _ := retData["info"]
	if nil == data {
		return nil, nil
	}
	return data.([]interface{}), nil
}

//GetHttpResult get http result
func GetHttpResult(req *restful.Request, url, method string, params interface{}) (bool, string, interface{}) {
	var strParams []byte
	switch params.(type) {
	case string:
		strParams = []byte(params.(string))
	default:
		strParams, _ = json.Marshal(params)

	}
	blog.Info("get request url:%s", url)
	blog.Info("get request info  params:%v", string(strParams))
	reply, err := httpcli.ReqHttp(req, url, method, []byte(strParams))
	blog.Info("get request result:%v", string(reply))
	if err != nil {
		blog.Error("http do error, params:%s, error:%s", strParams, err.Error())
		return false, err.Error(), nil
	}

	addReply, err := simplejson.NewJson([]byte(reply))
	if err != nil {
		blog.Error("http do error, params:%s, reply:%s, error:%s", strParams, reply, err.Error())
		return false, err.Error(), nil
	}
	isSuccess, err := addReply.Get("result").Bool()
	if nil != err || !isSuccess {
		errMsg, _ := addReply.Get("message").String()
		blog.Error("http do error, url:%s, params:%s, error:%s", url, strParams, errMsg)
		return false, errMsg, addReply.Get("data").Interface()
	}
	return true, "", addReply.Get("data").Interface()
}

//MoveHostToResourcePool move host to resource pool
func MoveHost2ResourcePool(CC *api.APIResource, req *restful.Request, appID int, hostID []int) (interface{}, error) {
	user := sencecommon.GetUserFromHeader(req)

	conds := make(map[string]interface{})
	conds[common.BKAppIDField] = appID
	appinfo, err := GetAppInfo(req, common.BKOwnerIDField, conds, CC.ObjCtrl())
	if err != nil {
		return nil, err
	}
	ownerID := appinfo[common.BKOwnerIDField].(string)
	if "" == ownerID {
		return nil, errors.New("未找到资源池")
	}
	//get default biz
	ownerAppID, err := GetDefaultAppID(req, ownerID, common.BKAppIDField, CC.ObjCtrl())
	if err != nil {
		return nil, errors.New("获取资源池信息失败，" + err.Error())

	}
	if 0 == appID {
		return nil, errors.New("资源池不存在")
	}
	if ownerAppID == appID {
		return nil, errors.New("当前主机已经属于资源池，不需要转移")
	}

	//get resource set
	mconds := make(map[string]interface{})
	mconds[common.BKDefaultField] = common.DefaultResModuleFlag
	mconds[common.BKModuleNameField] = common.DefaultResModuleName
	mconds[common.BKAppIDField] = ownerAppID
	moduleID, err := GetSingleModuleID(req, mconds, CC.ObjCtrl())
	if nil != err {
		return nil, errors.New("获取资源池模块信息失败" + err.Error())
	}

	logClient, err := NewHostModuleConfigLog(req, nil, CC.HostCtrl(), CC.ObjCtrl(), CC.AuditCtrl())

	conds[common.BKHostIDField] = hostID
	conds["bk_owner_module_id"] = moduleID
	conds["bk_owner_biz_id"] = ownerAppID
	url := CC.HostCtrl() + "/host/v1/meta/hosts/resource"
	isSucess, errmsg, data := GetHttpResult(req, url, common.HTTPUpdate, conds)
	if !isSucess {
		return data, errors.New("更新主机关系失败;" + errmsg)
	}
	logClient.SetDescSuffix("; 转移主机到资源池")
	logClient.SaveLog(fmt.Sprintf("%d", appID), user)

	return data, err
}

//GetAppInfo get app info
func GetAppInfo(req *restful.Request, fields string, conditon map[string]interface{}, hostAddr string) (map[string]interface{}, error) {
	//moduleURL := "http://" + cc.ObjCtrl + "/object/v1/insts/module/search"
	URL := hostAddr + "/object/v1/insts/" + common.BKInnerObjIDApp + "/search"
	params := make(map[string]interface{})
	params["condition"] = conditon
	params["sort"] = common.BKAppIDField
	params["start"] = 0
	params["limit"] = 1
	params["fields"] = fields

	blog.Info("get application info  url:%s", URL)
	blog.Info("get application info  url:%v", params)
	isSuccess, errMsg, data := GetHttpResult(req, URL, common.HTTPSelectPost, params)
	if !isSuccess {
		blog.Error("get application info  error, params:%v, error:%s", params, errMsg)
		return nil, errors.New(errMsg)
	}
	dataInterface := data.(map[string]interface{})
	info := dataInterface["info"].([]interface{})
	if 1 != len(info) {
		blog.Error("not application info error, params:%v, error:%s", params, errMsg)
		return nil, errors.New("业务不存在")
	}
	row := info[0].(map[string]interface{})

	if 0 == len(row) {
		blog.Error("not application info error, params:%v, error:%s", params, errMsg)
		return nil, errors.New("业务存在")
	}

	return row, nil
}

//GetDefaultAppID get default biz id
func GetDefaultAppID(req *restful.Request, ownerID, fields, hostAddr string) (int, error) {
	conds := make(map[string]interface{})
	conds[common.BKOwnerIDField] = ownerID
	conds[common.BKDefaultField] = common.DefaultAppFlag
	appinfo, err := GetAppInfo(req, fields, conds, hostAddr)
	if nil != err {
		blog.Errorf("get default app info error:%v", err.Error())
		return 0, err
	}
	return util.GetIntByInterface(appinfo[common.BKAppIDField])
}

//GetDefaultAppID get supplier ID
func GetDefaultAppIDBySupplierID(req *restful.Request, supplierID int, fields, hostAddr string) (int, error) {
	conds := make(map[string]interface{})
	conds[common.BKSupplierIDField] = supplierID
	conds[common.BKDefaultField] = common.DefaultAppFlag
	appinfo, err := GetAppInfo(req, fields, conds, hostAddr)
	if nil != err {
		blog.Errorf("get default app info error:%v", err.Error())
		return 0, err
	}
	return util.GetIntByInterface(appinfo[common.BKAppIDField])
}

//IsExistHostIDInApp  is host exsit in app
func IsExistHostIDInApp(CC *api.APIResource, req *restful.Request, appID int, hostID int) (bool, error) {
	conds := common.KvMap{common.BKAppIDField: appID, common.BKHostIDField: hostID}
	url := CC.HostCtrl() + "/host/v1/meta/hosts/modules/search"
	isSucess, errmsg, data := GetHttpResult(req, url, common.HTTPSelectPost, conds)
	blog.Info("IsExistHostIDInApp request url:%s, params:{appid:%d, hostid:%d}", url, appID, hostID)
	blog.Info("IsExistHostIDInApp res:%v,%s, %v", isSucess, errmsg, data)
	if !isSucess {
		return false, errors.New("获取主机关系失败;" + errmsg)
	}
	//数据为空
	if nil == data {
		return false, nil
	}
	ids, ok := data.([]interface{})
	if !ok {
		return false, errors.New(fmt.Sprintf("获取主机关系返回值格式错误;%v", data))
	}

	if len(ids) > 0 {
		return true, nil
	}
	return false, nil

}
