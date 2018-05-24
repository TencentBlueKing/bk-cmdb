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

package openapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/scene_server/host_server/host_service/logics"
	phpapilogic "configcenter/src/scene_server/host_server/host_service/logics/phpapi"
)

// HostSearchByIP: 根据IP查询主机, 多个IP以英文逗号分隔
func (cli *hostAction) HostSearchByIP(req *restful.Request, resp *restful.Response) {

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Error("Unmarshal json failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	ipArr, hasIp := input[common.BKIPListField]
	if !hasIp {
		blog.Error("input does not contains key IP")
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	appIDArrInput, hasAppID := input[common.BKAppIDField]
	subArea, hasSubArea := input[common.BKCloudIDField]

	orCondition := []map[string]interface{}{
		map[string]interface{}{common.BKHostInnerIPField: map[string]interface{}{common.BKDBIN: ipArr}},
		map[string]interface{}{common.BKHostOuterIPField: map[string]interface{}{common.BKDBIN: ipArr}},
	}
	hostMapCondition := map[string]interface{}{common.BKDBOR: orCondition}

	if hasSubArea && subArea != nil && subArea != "" {
		hostMapCondition[common.BKCloudIDField] = subArea
	}

	hostMap, hostIDArr, err := phpapilogic.GetHostMapByCond(req, hostMapCondition)
	if err != nil {
		blog.Error("getHostMapByCond error : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return
	}

	configCond := map[string]interface{}{
		common.BKHostIDField: hostIDArr,
	}
	if hasAppID {
		configCond[common.BKAppIDField] = appIDArrInput
	}

	configData, err := logics.GetConfigByCond(req, host.CC.HostCtrl(), configCond)

	hostData, err := phpapilogic.SetHostData(req, configData, hostMap)
	if nil != err {
		blog.Error("HostSearchByIP error : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return
	}

	cli.ResponseSuccess(hostData, resp)
}

// HostSearchByModuleID: 根据ModuleID查询主机
func (cli *hostAction) HostSearchByModuleID(req *restful.Request, resp *restful.Response) {

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Error("Unmarshal json failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	appID, hasAppID := input[common.BKAppIDField]
	if !hasAppID {
		blog.Error("input does not contains key ApplicationID")
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	moduleIDArr, hasModuleID := input[common.BKModuleIDField]
	if !hasModuleID {
		blog.Error("input does not contains key ModuleID")
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	configData, err := logics.GetConfigByCond(req, host.CC.HostCtrl(), map[string]interface{}{
		common.BKModuleIDField: moduleIDArr.([]interface{}),
		common.BKAppIDField:    []interface{}{appID},
	})
	if nil != err {
		cli.respGetHostFailed(resp, err)
		return
	}

	hostData, err := phpapilogic.GetHostDataByConfig(req, configData)
	if nil != err {
		cli.respGetHostFailed(resp, err)
		return
	}

	cli.ResponseSuccess(hostData, resp)
}

// HostSearchBySetID: 根据SetID查询主机
func (cli *hostAction) HostSearchBySetID(req *restful.Request, resp *restful.Response) {

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)

	if nil != err {
		blog.Error("Unmarshal json failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	appID, hasAppID := input[common.BKAppIDField]
	if !hasAppID {
		blog.Error("input does not contains key ApplicationID")
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	setIDArr, _ := input[common.BKSetIDField].([]interface{})
	conds := common.KvMap{common.BKAppIDField: []interface{}{appID}}
	if len(setIDArr) > 0 {
		conds[common.BKSetIDField] = setIDArr
	}

	configData, err := logics.GetConfigByCond(req, host.CC.HostCtrl(), conds)
	if nil != err {
		cli.respGetHostFailed(resp, err)
		return
	}

	hostData, err := phpapilogic.GetHostDataByConfig(req, configData)
	if nil != err {
		cli.respGetHostFailed(resp, err)
		return
	}

	cli.ResponseSuccess(hostData, resp)
}

// HostSearchByAppID: 根据业务ID查询主机
func (cli *hostAction) HostSearchByAppID(req *restful.Request, resp *restful.Response) {

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Error("Unmarshal json failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	appID, hasAppID := input[common.BKAppIDField]
	if !hasAppID {
		blog.Error("input does not contains key ApplicationID")
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	configData, err := logics.GetConfigByCond(req, host.CC.HostCtrl(), map[string]interface{}{
		common.BKAppIDField: []interface{}{appID},
	})
	if nil != err {
		cli.respGetHostFailed(resp, err)
		return
	}

	hostData, err := phpapilogic.GetHostDataByConfig(req, configData)
	if nil != err {
		cli.respGetHostFailed(resp, err)
		return
	}

	cli.ResponseSuccess(hostData, resp)
}

// HostSearchByProperty: 查根据set属性查询主机
func (cli *hostAction) HostSearchByProperty(req *restful.Request, resp *restful.Response) {

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Error("Unmarshal json failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	appID, hasAppID := input[common.BKAppIDField]
	if !hasAppID {
		blog.Error("input does not contains key ApplicationID")
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	setCondition := make([]interface{}, 0)

	setIDArrI, hasSetID := input[common.BKSetIDField]
	if hasSetID {
		cond := make(map[string]interface{})
		cond["field"] = common.BKSetIDField
		cond["operator"] = common.BKDBIN
		cond["value"] = setIDArrI
		setCondition = append(setCondition, cond)
	}

	setEnvTypeArr, hasSetEnvType := input[common.BKSetEnvField]
	if hasSetEnvType {
		cond := make(map[string]interface{})
		cond["field"] = common.BKSetEnvField
		cond["operator"] = common.BKDBIN
		cond["value"] = setEnvTypeArr
		setCondition = append(setCondition, cond)
	}

	setSrvStatusArr, hasSetSrvStatus := input[common.BKSetStatusField]
	if hasSetSrvStatus {
		cond := make(map[string]interface{})
		cond["field"] = common.BKSetStatusField
		cond["operator"] = common.BKDBIN
		cond["value"] = setSrvStatusArr
		setCondition = append(setCondition, cond)
	}

	blog.Debug("setCondition: %v\n", setCondition)
	setIDArr, err := logics.GetSetIDByCond(req, host.CC.ObjCtrl(), setCondition)
	blog.Debug("ApplicationID: %s, SetID: %v\n", appID, setIDArr)

	condition := map[string]interface{}{
		common.BKAppIDField: []interface{}{appID},
	}

	condition[common.BKSetIDField] = setIDArr

	configData, err := logics.GetConfigByCond(req, host.CC.HostCtrl(), condition)

	if nil != err {
		cli.respGetHostFailed(resp, err)
		return
	}

	hostData, err := phpapilogic.GetHostDataByConfig(req, configData)
	if nil != err {
		cli.respGetHostFailed(resp, err)
		return
	}

	cli.ResponseSuccess(hostData, resp)
}

func (cli *hostAction) respGetHostFailed(resp *restful.Response, err error) {
	blog.Error("get host error : %v", err)
	cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
	return
}

// GetHostListByAppidAndField 根据主机属性的值group主机列表
func (cli *hostAction) GetHostListByAppidAndField(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetHostListByAppidAndField start!")

	appID, err := strconv.Atoi(req.PathParameter(common.BKAppIDField))
	if nil != err {
		blog.Error("convert appid to int error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	field := req.PathParameter("field")

	configData, err := logics.GetConfigByCond(req, cli.CC.HostCtrl(), map[string]interface{}{
		common.BKAppIDField: []int{appID},
	})

	if nil != err {
		blog.Error("GetConfigByCond error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return
	}

	hostIDArr := make([]int, 0)
	for _, config := range configData {
		hostIDArr = append(hostIDArr, config[common.BKHostIDField])
	}

	//get host

	url := cli.CC.HostCtrl() + "/host/v1/hosts/search"
	searchParams := map[string]interface{}{
		"fields": fmt.Sprintf("%s,%s,%s,%s,", common.BKHostInnerIPField, common.BKCloudIDField, common.BKAppIDField, common.BKHostIDField) + field,
		"condition": map[string]interface{}{
			common.BKHostIDField: map[string]interface{}{
				common.BKDBIN: hostIDArr,
			},
		},
	}
	inputJson, _ := json.Marshal(searchParams)
	dataInfo, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(inputJson))
	if nil != err {
		blog.Error("get host by condition error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return
	}
	js, err := simplejson.NewJson([]byte(dataInfo))

	res, _ := js.Map()
	resData := res["data"].(map[string]interface{})
	hostData := resData["info"].([]interface{})

	retData := make(map[string][]interface{})
	blog.Debug("host data: %v", hostData)
	for _, item := range hostData {
		itemMap := item.(map[string]interface{})
		fieldValue, ok := itemMap[field]
		if !ok {
			continue
		}

		fieldValueStr := fmt.Sprintf("%v", fieldValue)
		groupData, ok := retData[fieldValueStr]
		if ok {
			retData[fieldValueStr] = append(groupData, itemMap)
		} else {
			retData[fieldValueStr] = []interface{}{
				itemMap,
			}
		}

	}

	cli.ResponseSuccess(retData, resp)
}

// getIPAndProxyByCompany 获取Company下proxy列表
func (cli *hostAction) GetIPAndProxyByCompany(req *restful.Request, resp *restful.Response) {

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Error("Unmarshal json failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	ipArr, hasIp := input["ips"]
	if !hasIp {
		blog.Error("input does not contains key IP")
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	appID, _ := input[common.BKAppIDField]
	appIDInt, _ := strconv.Atoi(appID.(string))
	platID, _ := input[common.BKCloudIDField]
	platIDInt, _ := strconv.Atoi(platID.(string))
	// 获取不合法的IP列表
	param := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKHostInnerIPField: map[string]interface{}{common.BKDBIN: ipArr},
			common.BKCloudIDField:     platIDInt,
		},
		"fields": fmt.Sprintf("%s,%s", common.BKHostIDField, common.BKHostInnerIPField),
	}

	hosts, err := phpapilogic.GetHostByCond(req, param, host.CC.ObjCtrl())

	hostIDArr := make([]interface{}, 0)
	hostMap := make(map[string]interface{})

	for _, host := range hosts.([]interface{}) {
		hostID := host.(map[string]interface{})[common.BKHostIDField]
		hostIDArr = append(hostIDArr, hostID)
		hostMap[fmt.Sprintf("%v", hostID)] = host
	}

	if nil != err {
		blog.Error("getHostByIPArrAndSource failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}
	blog.Debug("hostIDArr:%v", hostIDArr)
	muduleHostConfigs, err := logics.GetConfigByCond(req, cli.CC.HostCtrl(), map[string]interface{}{
		common.BKHostIDField: hostIDArr,
	})
	blog.Debug("vaildIPArr:%v", muduleHostConfigs)

	validIpArr := make([]interface{}, 0)

	appMap, err := logics.GetAppMapByCond(req, "", host.CC.ObjCtrl(), map[string]interface{}{})

	invalidIpMap := make(map[string]map[string]interface{})

	for _, config := range muduleHostConfigs {
		appIDTemp := fmt.Sprintf("%v", config[common.BKAppIDField])
		appIDIntTemp := config[common.BKAppIDField]
		hostID := config[common.BKHostIDField]
		ip := hostMap[fmt.Sprintf("%v", hostID)].(map[string]interface{})[common.BKHostInnerIPField]

		appName := appMap[appIDIntTemp].(map[string]interface{})[common.BKAppNameField]

		if appIDIntTemp != appIDInt {

			_, ok := invalidIpMap[appIDTemp]
			if !ok {
				invalidIpMap[appIDTemp] = make(map[string]interface{})
				invalidIpMap[appIDTemp][common.BKAppNameField] = appName
				invalidIpMap[appIDTemp]["ips"] = make([]string, 0)
			}

			invalidIpMap[appIDTemp]["ips"] = append(invalidIpMap[appIDTemp]["ips"].([]string), ip.(string))

		} else {
			validIpArr = append(validIpArr, ip)
		}
	}

	// 获取所有的proxy ip列表
	paramProxy := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKGseProxyField: 1,
			common.BKCloudIDField:  platIDInt,
		},
		"fields": fmt.Sprintf("%s,%s", common.BKHostIDField, common.BKHostInnerIPField),
	}

	hostProxys, err := phpapilogic.GetHostByCond(req, paramProxy, host.CC.ObjCtrl())

	proxyIpArr := make([]interface{}, 0)

	for _, host := range hostProxys.([]interface{}) {
		h := make(map[string]interface{})
		h[common.BKHostInnerIPField] = host.(map[string]interface{})[common.BKHostInnerIPField]
		h[common.BKHostOuterIPField] = ""
		proxyIpArr = append(proxyIpArr, h)
	}
	blog.Debug("proxyIpArr:%v", proxyIpArr)

	resData := make(map[string]interface{})
	resData[common.BKIPListField] = validIpArr
	resData[common.BKProxyListField] = proxyIpArr
	resData[common.BKInvalidIPSField] = invalidIpMap
	cli.ResponseSuccess(resData, resp)
}

//GetHostAppByCompanyId: 根据开发商ID和平台ID获取业务
func (cli *hostAction) GetHostAppByCompanyId(req *restful.Request, resp *restful.Response) {
	value, _ := ioutil.ReadAll(req.Request.Body)
	js, err := simplejson.NewJson([]byte(value))
	if err != nil {
		blog.Error("GetHostAppByCompanyId failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}
	input, err := js.Map()
	if err != nil {
		blog.Error("GetHostAppByCompanyId failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}
	blog.Debug(" input:%v", input)
	ownerId := input[common.BKOwnerIDField]
	platId, _ := strconv.Atoi(input[common.BKCloudIDField].(string))
	ip := input["ip"].(string)
	ipArr := strings.Split(ip, ",")
	hostCon := map[string]interface{}{
		common.BKHostInnerIPField: map[string]interface{}{
			common.BKDBIN: ipArr,
		},
		common.BKCloudIDField: platId,
	}
	//根据i,platId获取主机
	hostArr, hostIdArr, err := phpapilogic.GetHostMapByCond(req, hostCon)
	if nil != err {
		blog.Error("GetHostAppByCompanyId getHostMapByCond:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return
	}
	blog.Debug("GetHostAppByCompanyId hostArr:%v", hostArr)
	if len(hostIdArr) == 0 {
		cli.ResponseSuccess("", resp)
		return
	}
	// 根据主机hostId获取app_id,module_id,set_id
	configCon := map[string]interface{}{
		common.BKHostIDField: hostIdArr,
	}
	configArr, err := getConfigByCond(req, cli.CC.HostCtrl(), configCon)
	if nil != err {
		blog.Error("GetHostAppByCompanyId getConfigByCond:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return
	}
	blog.Debug("GetHostAppByCompanyId configArr:%v", configArr)
	if len(configArr) == 0 {
		cli.ResponseSuccess("", resp)
		return
	}
	appIdArr := make([]int, 0)
	setIdArr := make([]int, 0)
	moduleIdArr := make([]int, 0)
	for _, item := range configArr {
		appIdArr = append(appIdArr, item[common.BKAppIDField])
		setIdArr = append(setIdArr, item[common.BKSetIDField])
		moduleIdArr = append(moduleIdArr, item[common.BKModuleIDField])
	}
	hostMapArr, err := phpapilogic.SetHostData(req, configArr, hostArr)
	if nil != err {
		blog.Error("GetHostAppByCompanyId setHostData:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return
	}
	blog.Debug("GetHostAppByCompanyId hostMap:%v", hostMapArr)
	hostDataArr := make([]interface{}, 0)
	for _, h := range hostMapArr {
		hostMap := h.(map[string]interface{})
		if hostMap[common.BKOwnerIDField] == ownerId {
			hostDataArr = append(hostDataArr, hostMap)
		}
	}
	cli.ResponseSuccess(hostDataArr, resp)
}

func (cli *hostAction) GetGitServerIp(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetGitServerIp start")
	value, _ := ioutil.ReadAll(req.Request.Body)
	js, err := simplejson.NewJson([]byte(value))
	if err != nil {
		blog.Error("GetGitServerIp failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}
	input, err := js.Map()
	if err != nil {
		blog.Error("GetGitServerIp failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}
	blog.Debug(" input:%v", input)
	appName := input[common.BKAppNameField].(string)
	setName := input[common.BKSetNameField].(string)
	moduleName := input[common.BKModuleNameField].(string)
	var appId, setId, moduleId int

	// 根据appName获取app
	appCondition := map[string]interface{}{
		common.BKAppNameField: appName,
	}
	appMap, err := logics.GetAppMapByCond(req, "", cli.CC.ObjCtrl(), appCondition)
	for key, _ := range appMap {
		appId = key
	}
	if err != nil {
		blog.Error("GetGitServerIp GetAppMapByCond err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL, resp)
		return
	}
	if len(appMap) == 0 {
		cli.ResponseSuccess("", resp)
		return
	}
	// 根据setName获取set信息
	setCondition := map[string]interface{}{
		common.BKSetNameField: setName,
		common.BKAppIDField:   appId,
	}
	setMap, err := logics.GetSetMapByCond(req, "", cli.CC.ObjCtrl(), setCondition)
	for key, _ := range setMap {
		setId = key
	}
	if err != nil {
		blog.Error("GetGitServerIp GetSetMapByCond err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL, resp)
		return
	}
	if len(setMap) == 0 {
		cli.ResponseSuccess("", resp)
		return
	}

	// 根据moduleName获取module信息
	moduleCondition := map[string]interface{}{
		common.BKModuleNameField: moduleName,
		common.BKAppIDField:      appId,
	}
	moduleMap, err := logics.GetModuleMapByCond(req, "", cli.CC.ObjCtrl(), moduleCondition)
	for key, _ := range moduleMap {
		moduleId = key
	}
	if err != nil {
		blog.Error("GetGitServerIp GetModuleMapByCond err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL, resp)
		return
	}
	if len(moduleMap) == 0 {
		cli.ResponseSuccess(nil, resp)
		return
	}
	// 根据 appId,setId,moduleId 获取主机信息
	//configData := make([]map[string]int,0)
	confMap := map[string]interface{}{
		common.BKAppIDField:    []int{appId},
		common.BKSetIDField:    []int{setId},
		common.BKModuleIDField: []int{moduleId},
	}
	configData, err := logics.GetConfigByCond(req, cli.CC.HostCtrl(), confMap)
	blog.Debug("configData:%v", configData)
	hostArr, err := phpapilogic.GetHostDataByConfig(req, configData)
	if nil != err {
		blog.Error("GetGitServerIp getHostDataByConfig err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL, resp)
		return
	}
	blog.Debug("hostArr:%v", hostArr)

	cli.ResponseSuccess(hostArr, resp)
}
