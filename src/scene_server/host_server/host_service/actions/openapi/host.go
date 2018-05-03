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
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/logics"
	sourceAPI "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	simplejson "github.com/bitly/go-simplejson"

	errorIfs "configcenter/src/common/errors"
	"configcenter/src/scene_server/validator"
	"errors"
	"strings"

	"configcenter/src/common/auditoplog"
	"configcenter/src/source_controller/api/auditlog"

	"github.com/emicklei/go-restful"
)

var host *hostAction = &hostAction{}

type hostAction struct {
	base.BaseAction
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/openapi/host/{" + common.BKAppIDField + "}", Params: nil, Handler: host.UpdateHost})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/host/updateHostByAppID/{appid}", Params: nil, Handler: host.UpdateHostByAppID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/host/getHostListByAppidAndField/{" + common.BKAppIDField + "}/{field}", Params: nil, Handler: host.GetHostListByAppidAndField})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/gethostlistbyip", Params: nil, Handler: host.HostSearchByIP})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/getmodulehostlist", Params: nil, Handler: host.HostSearchByModuleID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/getsethostlist", Params: nil, Handler: host.HostSearchBySetID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/getapphostlist", Params: nil, Handler: host.HostSearchByAppID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/gethostsbyproperty", Params: nil, Handler: host.HostSearchByProperty})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/getIPAndProxyByCompany", Params: nil, Handler: host.GetIPAndProxyByCompany})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/openapi/updatecustomproperty", Params: nil, Handler: host.UpdateCustomProperty})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "openapi/host/clonehostproperty", Params: nil, Handler: host.CloneHostProperty})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/openapi/host/getHostAppByCompanyId", Params: nil, Handler: host.GetHostAppByCompanyId})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/openapi/host/delhostinapp", Params: nil, Handler: host.DelHostInApp})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/openapi/host/getGitServerIp", Params: nil, Handler: host.GetGitServerIp})

	// create CC object
	host.CreateAction()
}

// HostSearchByIP: 根据IP查询主机, 多个IP以英文逗号分隔
func (cli *hostAction) HostSearchByIP(req *restful.Request, resp *restful.Response) {
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		input := make(map[string]interface{})
		err = json.Unmarshal(value, &input)
		if nil != err {
			blog.Error("Unmarshal json failed, error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)

		}

		ipArr, hasIp := input[common.BKIPListField]
		if !hasIp {
			blog.Error("input does not contains key IP")
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
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

		hostMap, hostIDArr, err := getHostMapByCond(req, hostMapCondition)
		if err != nil {
			blog.Error("getHostMapByCond error : %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}

		configCond := map[string]interface{}{
			common.BKHostIDField: hostIDArr,
		}
		if hasAppID {
			configCond[common.BKAppIDField] = appIDArrInput
		}

		configData, err := logics.GetConfigByCond(req, host.CC.HostCtrl(), configCond)

		hostData, err := setHostData(req, configData, hostMap)
		if nil != err {
			blog.Error("HostSearchByIP error : %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}

		return http.StatusOK, hostData, nil
	}, resp)
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
		cli.respGetHostFailed(req, resp, err)
		return
	}

	hostData, err := getHostDataByConfig(req, configData)
	if nil != err {
		cli.respGetHostFailed(req, resp, err)
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
		cli.respGetHostFailed(req, resp, err)
		return
	}

	hostData, err := getHostDataByConfig(req, configData)
	if nil != err {
		cli.respGetHostFailed(req, resp, err)
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
		cli.respGetHostFailed(req, resp, err)
		return
	}

	hostData, err := getHostDataByConfig(req, configData)
	if nil != err {
		cli.respGetHostFailed(req, resp, err)
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
		cli.respGetHostFailed(req, resp, err)
		return
	}

	hostData, err := getHostDataByConfig(req, configData)
	if nil != err {
		cli.respGetHostFailed(req, resp, err)
		return
	}

	cli.ResponseSuccess(hostData, resp)
}

func (cli *hostAction) respGetHostFailed(req *restful.Request, resp *restful.Response, err error) {
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	cli.CallResponseEx(func() (int, interface{}, error) {
		blog.Error("get host error : %v", err)
		return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
	}, resp)
}

// GetHostListByAppidAndField 根据主机属性的值group主机列表
func (cli *hostAction) GetHostListByAppidAndField(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetHostListByAppidAndField start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	cli.CallResponseEx(func() (int, interface{}, error) {
		appID, err := strconv.Atoi(req.PathParameter(common.BKAppIDField))
		if nil != err {
			blog.Error("convert appid to int error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}

		field := req.PathParameter("field")

		configData, err := logics.GetConfigByCond(req, cli.CC.HostCtrl(), map[string]interface{}{
			common.BKAppIDField: []int{appID},
		})

		if nil != err {
			blog.Error("GetConfigByCond error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
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
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
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

		// cli.ResponseSuccess(retData, resp)
		return http.StatusOK, resData, nil
	}, resp)
}

// updateHostPlat 根据条件更新主机信息
func (cli *hostAction) UpdateHost(req *restful.Request, resp *restful.Response) {
	blog.Debug("updateHost start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	cli.CallResponseEx(func() (int, interface{}, error) {
		appID, err := strconv.Atoi(req.PathParameter(common.BKAppIDField))
		if nil != err {
			blog.Error("convert appid to int error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}

		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		blog.Debug("updateHost http body data: %s", value)

		input := make(map[string]interface{})
		err = json.Unmarshal(value, &input)
		if nil != err {
			blog.Error("unmarshal json error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}
		blog.Debug("input:%s", input, string(value))

		updateData, ok := input["data"]
		if !ok {
			blog.Error("params data is required:%s", string(value))
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}
		mapData, ok := updateData.(map[string]interface{})
		if !ok {
			blog.Error("params data must be object:%s", string(value))
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}

		dstPlat, ok := mapData[common.BKSubAreaField]
		if !ok {
			blog.Error("params data.bk_cloud_id is require:%s", string(value))
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}

		// dst host exist return souccess, hongsong tiyi
		dstHostCondition := map[string]interface{}{
			common.BKHostInnerIPField: input["condition"].(map[string]interface{})[common.BKHostInnerIPField],
			common.BKCloudIDField:     dstPlat,
		}
		_, hostIDArr, err := getHostMapByCond(req, dstHostCondition)
		blog.Debug("hostIDArr:%v", hostIDArr)
		if nil != err {
			blog.Error("updateHostMain error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostModifyFail)
		}

		if len(hostIDArr) != 0 {
			// cli.ResponseSuccess(nil, resp)
			return http.StatusOK, nil, nil
		}

		blog.Debug(input["condition"].(map[string]interface{})[common.BKCloudIDField])
		hostCondition := map[string]interface{}{
			common.BKHostInnerIPField: input["condition"].(map[string]interface{})[common.BKHostInnerIPField],
			common.BKCloudIDField:     input["condition"].(map[string]interface{})[common.BKCloudIDField],
		}
		data := input["data"].(map[string]interface{})
		data[common.BKHostInnerIPField] = input["condition"].(map[string]interface{})[common.BKHostInnerIPField]
		res, err := updateHostMain(req, hostCondition, data, appID, cli.CC.HostCtrl(), cli.CC.ObjCtrl(), cli.CC.AuditCtrl(), cli.CC.Error)

		if nil != err {
			blog.Error("updateHostMain error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostModifyFail)
		}

		cli.ResponseSuccess(res, resp)
		return http.StatusOK, nil, nil
	}, resp)
}

// updateHostByAppID 根据IP更新主机Proxy状态，如果不存在主机则添加到对应业务及默认模块
func (cli *hostAction) UpdateHostByAppID(req *restful.Request, resp *restful.Response) {
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	cli.CallResponseEx(func() (int, interface{}, error) {
		blog.Debug("updateHostByAppID start!")
		appID, err := strconv.Atoi(req.PathParameter("appid"))
		if nil != err {
			blog.Error("convert appid to int error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}

		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		blog.Debug("updateHostByAppID http body data: %s", value)

		input := make(map[string]interface{})
		err = json.Unmarshal(value, &input)
		if nil != err {
			blog.Error("unmarshal json error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		proxyArr := input[common.BKProxyListField].([]interface{})
		platID, _ := util.GetIntByInterface(input[common.BKCloudIDField])

		blog.Debug("proxyArr:%v", proxyArr)
		defaultModule, err := getDefaultModules(req, appID, cli.CC.ObjCtrl())

		if nil != err {
			blog.Error("getDefaultModules error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}

		//defaultSetID := defaultModule["SetID"]
		defaultModuleID := (defaultModule.(map[string]interface{}))[common.BKModuleIDField]

		for _, pro := range proxyArr {
			proMap := pro.(map[string]interface{})
			var hostID int
			innerIP := proMap[common.BKHostInnerIPField]
			outerIP, ok := proMap[common.BKHostOuterIPField]
			if !ok {
				outerIP = ""
			}

			hostData, err := getHostByIPAndSource(req, innerIP.(string), platID, cli.CC.ObjCtrl())
			blog.Error("hostData:%v", hostData)
			if nil != err {
				blog.Error("getHostByIPAndSource error:%v", err)
				return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
			}

			hostDataArr := hostData.([]interface{})

			if len(hostDataArr) == 0 {
				blog.Debug("procMap:%v", proMap)
				hostIDNew, err := addHost(req, proMap, cli.CC.ObjCtrl())

				if nil != err {
					blog.Error("addHost error:%v", err)
					return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostUpdateFail)
				}

				hostID = hostIDNew

				blog.Debug("addHost success, hostID: %d", hostID)

				err = addModuleHostConfig(req, map[string]interface{}{
					common.BKAppIDField:    appID,
					common.BKModuleIDField: []float64{defaultModuleID.(float64)},
					common.BKHostIDField:   hostID,
				}, cli.CC.HostCtrl())

				if nil != err {
					blog.Error("addModuleHostConfig error:%v", err)
					return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostUpdateFail)
				}

			} else {
				hostMap := hostDataArr[0].(map[string]interface{})
				hostIDTemp := hostMap[common.BKHostIDField].(float64)
				hostID = int(hostIDTemp)
			}

			if outerIP != "" {
				hostCondition := map[string]interface{}{
					common.BKHostIDField: hostID,
				}
				data := map[string]interface{}{
					// TODO 没有gse_proxy字段，暂时不修改;2018/03/09
					//common.BKGseProxyField: 1,
				}

				_, err := updateHostMain(req, hostCondition, data, appID, cli.CC.HostCtrl(), cli.CC.ObjCtrl(), cli.CC.AuditCtrl(), cli.CC.Error)
				if nil != err {
					blog.Error("updateHostMain error:%v", err)
					return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostUpdateFail)
				}
			}

		}

		// cli.ResponseSuccess(nil, resp)
		return http.StatusOK, nil, nil
	}, resp)
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

	hosts, err := getHostByCond(req, param, host.CC.ObjCtrl())

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

	hostProxys, err := getHostByCond(req, paramProxy, host.CC.ObjCtrl())

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

// updateCustomProperty 修改主机自定义属性
func (cli *hostAction) UpdateCustomProperty(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)

	// 获取该语系下的错误码
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		input := make(map[string]interface{})
		err = json.Unmarshal(value, &input)
		if nil != err {
			blog.Error("Unmarshal json failed, error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)

		}
		blog.Error("UpdateCustomProperty :%v", input)
		appID, _ := strconv.Atoi(input[common.BKAppIDField].(string))
		hostID, _ := strconv.Atoi(input[common.BKHostIDField].(string))
		propertyJSON := input["property"]

		propertyMap := make(map[string]interface{})
		if nil != propertyJSON {
			err = json.Unmarshal([]byte(propertyJSON.(string)), &propertyMap)
		}
		if nil != err {
			blog.Error("Unmarshal json failed, error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppUpdateFailed)
		}
		condition := make(common.KvMap)
		condition[common.BKAppIDField] = appID
		fileds := fmt.Sprintf("%s,%s", common.BKAppIDField, common.BKOwnerIDField)
		apps, err := logics.GetAppMapByCond(req, fileds, cli.CC.ObjCtrl(), condition)
		if nil != err {
			blog.Error("UpdateCustomProperty GetAppMapByCond, error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppUpdateFailed)

		}
		blog.Debug("UpdateCustomProperty apps:%v", apps)
		if _, ok := apps[appID]; !ok {
			msg := "业务不存在"
			blog.Debug("UpdateCustomProperty error:%v", msg)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppUpdateFailed)
		}

		appMap := apps[appID]
		ownerID := appMap.(map[string]interface{})[common.BKOwnerIDField]
		propertys, _ := getCustomerPropertyByOwner(req, ownerID, cli.CC.ObjCtrl())
		params := make(common.KvMap)
		for _, attrMap := range propertys {
			PropertyID, ok := attrMap[common.BKPropertyIDField].(string)
			if !ok {
				continue
			}
			blog.Debug("input[PropertyId]:%v", input[PropertyID])
			if _, ok := propertyMap[PropertyID]; ok {
				params[PropertyID] = propertyMap[PropertyID]
			}
		}
		blog.Debug("params:%v", params)
		hostCondition := map[string]interface{}{
			common.BKHostIDField: hostID,
		}
		_, err = updateHostMain(req, hostCondition, params, appID, cli.CC.HostCtrl(), cli.CC.ObjCtrl(), cli.CC.AuditCtrl(), cli.CC.Error)
		if nil != err {
			blog.Error("UpdateCustomProperty updateHostMain error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppUpdateFailed)
		}
		return http.StatusOK, nil, nil
	}, resp)
}

// cloneHostProperty 克隆主机
func (cli *hostAction) CloneHostProperty(req *restful.Request, resp *restful.Response) {
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		input := make(map[string]interface{})
		err = json.Unmarshal(value, &input)
		if nil != err {
			blog.Error("Unmarshal json failed, error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)

		}
		blog.Debug("CloneHostProperty input:%v", input)
		appID, _ := strconv.Atoi(input[common.BKAppIDField].(string))
		orgIP := input[common.BKOrgIPField]
		dstIP := input[common.BKDstIPField]

		platID, hasPlatID := input[common.BKCloudIDField]
		platIDInt, _ := strconv.Atoi(input[common.BKCloudIDField].(string))
		condition := common.KvMap{
			common.BKHostInnerIPField: orgIP,
		}

		if hasPlatID && platID != nil && platID != "" {
			condition[common.BKCloudIDField] = platIDInt
		}
		// 处理源IP
		hostMap, hostIDArr, err := getHostMapByCond(req, condition)

		blog.Debug("hostMapData:%v", hostMap)
		if err != nil {
			blog.Error("getHostMapByCond error : %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)

		}

		if len(hostIDArr) == 0 {
			blog.Error("clone host getHostMapByCond error, ip:%s, platid:%s", orgIP, platID)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}

		hostMapData, ok := hostMap[hostIDArr[0]].(map[string]interface{})
		if false == ok {
			blog.Error("getHostMapByCond not source ip : %s", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}

		configCond := map[string]interface{}{
			common.BKHostIDField: []interface{}{hostMapData[common.BKHostIDField]},
			common.BKAppIDField:  []int{appID},
		}
		// 判断源IP是否存在
		configData, err := logics.GetConfigByCond(req, host.CC.HostCtrl(), configCond)
		blog.Debug("configData:%v", configData)
		if nil != err {
			blog.Error("clone host property error : %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}
		if len(configData) == 0 {
			msg := "no find host"
			blog.Error("clone host property error : %v", msg)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}
		// 处理目标IP
		dstIPArr := strings.Split(dstIP.(string), ",")
		// 获得已存在的主机
		dstCondition := map[string]interface{}{
			common.BKHostInnerIPField: map[string]interface{}{
				common.BKDBIN: dstIPArr,
			},
			common.BKCloudIDField: platIDInt,
		}
		dstHostMap, dstHostIDArr, err := getHostMapByCond(req, dstCondition)
		blog.Debug("dstHostMap:%v", dstHostMap)

		dstConfigCond := map[string]interface{}{
			common.BKAppIDField:  []int{appID},
			common.BKHostIDField: dstHostIDArr,
		}
		dstHostIDArrV, err := logics.GetHostIDByCond(req, host.CC.HostCtrl(), dstConfigCond)
		existIPArr := make([]string, 0)
		for _, id := range dstHostIDArrV {
			if dstHostMapData, ok := dstHostMap[id].(map[string]interface{}); ok {
				existIPArr = append(existIPArr, dstHostMapData[common.BKHostInnerIPField].(string))
			}
		}

		//更新的时候，不修改为nil的数据
		updateHostData := make(map[string]interface{})
		for key, val := range hostMapData {
			if nil != val {
				updateHostData[key] = val
			}
		}
		// 克隆主机, 已存在的修改，不存在的新增；dstIpArr: 全部要克隆的主机，existIpArr：已存在的要克隆的主机
		blog.Debug("existIpArr:%v", existIPArr)
		for _, dstIPV := range dstIPArr {
			if dstIPV == orgIP {
				blog.Debug("clone host updateHostMain err:%v", err)
				// msg := "dstIp 和 orgIp不能相同"
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostUpdateFail)

			}
			blog.Debug("hostMapData:%v", hostMapData)
			if in_existIpArr(existIPArr, dstIPV) {
				blog.Debug("clone update")
				hostCondition := map[string]interface{}{
					common.BKHostInnerIPField: dstIPV,
				}

				updateHostData[common.BKHostInnerIPField] = dstIPV
				delete(updateHostData, common.BKHostIDField)
				res, err := updateHostMain(req, hostCondition, updateHostData, appID, host.CC.HostCtrl(), host.CC.ObjCtrl(), host.CC.AuditCtrl(), cli.CC.Error)
				if nil != err {
					blog.Debug("clone host updateHostMain err:%v", err)
					// msg := fmt.Sprintf("clone host error:%s", dstIpV)
					return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostUpdateFail)
				}
				blog.Debug("clone host updateHostMain res:%v", res)
			} else {
				hostMapData[common.BKHostInnerIPField] = dstIPV
				blog.Debug("clone add")
				addHostMapData := hostMapData
				delete(addHostMapData, common.BKHostIDField)
				cloneHostID, err := addHost(req, addHostMapData, host.CC.ObjCtrl())
				if nil != err {
					blog.Debug("clone host addHost err:%v", err)
					// msg := fmt.Sprintf("clone host error:%s", dstIpV)
					return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostCreateFail)
				}

				blog.Debug("cloneHostId:%v", cloneHostID)
				blog.Debug("configData[0]:%v", configData[0])
				configDataMap := make(map[string]interface{}, 0)
				configDataMap[common.BKHostIDField] = cloneHostID
				configDataMap[common.BKModuleIDField] = []int{configData[0][common.BKModuleIDField]}
				configDataMap[common.BKAppIDField] = configData[0][common.BKAppIDField]
				configDataMap[common.BKSetIDField] = configData[0][common.BKSetIDField]
				err = addModuleHostConfig(req, configDataMap, host.CC.HostCtrl())
				if nil != err {
					blog.Debug("clone host addModuleHostConfig err:%v", err)
					// msg := fmt.Sprintf("clone host error:%s", dstIpV)
					return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostCreateFail)

				}
			}
		}
		//成功
		return http.StatusOK, nil, nil
	}, resp)
}

func in_existIpArr(arr []string, ip string) bool {
	for _, v := range arr {
		if ip == v {
			return true
		}
	}
	return false
}

//GetHostAppByCompanyId: 根据开发商ID和平台ID获取业务
func (cli *hostAction) GetHostAppByCompanyId(req *restful.Request, resp *restful.Response) {
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	cli.CallResponseEx(func() (int, interface{}, error) {
		value, _ := ioutil.ReadAll(req.Request.Body)
		js, err := simplejson.NewJson([]byte(value))
		if err != nil {
			blog.Error("GetHostAppByCompanyId failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}
		input, err := js.Map()
		if err != nil {
			blog.Error("GetHostAppByCompanyId failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
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
		hostArr, hostIdArr, err := getHostMapByCond(req, hostCon)
		if nil != err {
			blog.Error("GetHostAppByCompanyId getHostMapByCond:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}
		blog.Debug("GetHostAppByCompanyId hostArr:%v", hostArr)
		if len(hostIdArr) == 0 {
			// cli.ResponseSuccess("", resp)
			return http.StatusOK, nil, nil
		}
		// 根据主机hostId获取app_id,module_id,set_id
		configCon := map[string]interface{}{
			common.BKHostIDField: hostIdArr,
		}
		configArr, err := getConfigByCond(req, cli.CC.HostCtrl(), configCon)
		if nil != err {
			blog.Error("GetHostAppByCompanyId getConfigByCond:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}
		blog.Debug("GetHostAppByCompanyId configArr:%v", configArr)
		if len(configArr) == 0 {
			// cli.ResponseSuccess("", resp)
			return http.StatusOK, nil, nil

		}
		appIdArr := make([]int, 0)
		setIdArr := make([]int, 0)
		moduleIdArr := make([]int, 0)
		for _, item := range configArr {
			appIdArr = append(appIdArr, item[common.BKAppIDField])
			setIdArr = append(setIdArr, item[common.BKSetIDField])
			moduleIdArr = append(moduleIdArr, item[common.BKModuleIDField])
		}
		hostMapArr, err := setHostData(req, configArr, hostArr)
		if nil != err {
			blog.Error("GetHostAppByCompanyId setHostData:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}
		blog.Debug("GetHostAppByCompanyId hostMap:%v", hostMapArr)
		hostDataArr := make([]interface{}, 0)
		for _, h := range hostMapArr {
			hostMap := h.(map[string]interface{})
			if hostMap[common.BKOwnerIDField] == ownerId {
				hostDataArr = append(hostDataArr, hostMap)
			}
		}
		// cli.ResponseSuccess(hostDataArr, resp)
		return http.StatusOK, hostDataArr, nil
	}, resp)
}

//DelHostInApp: 从业务空闲机集群中删除主机
func (cli *hostAction) DelHostInApp(req *restful.Request, resp *restful.Response) {
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	cli.CallResponseEx(func() (int, interface{}, error) {
		value, _ := ioutil.ReadAll(req.Request.Body)
		js, err := simplejson.NewJson([]byte(value))
		if err != nil {
			blog.Error("DelHostInApp failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}
		input, err := js.Map()
		if err != nil {
			blog.Error("DelHostInApp failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}
		blog.Debug(" input:%v", input)
		appId, _ := strconv.Atoi(input["appId"].(string))
		hostId, _ := strconv.Atoi(input["hostId"].(string))
		configCon := map[string]interface{}{
			"ApplicationID": []int{appId},
			"HostID":        []int{hostId},
		}

		configArr, err := logics.GetConfigByCond(req, cli.CC.HostCtrl(), configCon)
		if err != nil {
			blog.Error("DelHostInApp GetConfigByCond err msg : %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}
		if len(configArr) == 0 {
			msg := fmt.Sprintf("not fint hostId:%v in appId:%v", hostId, appId)
			blog.Info("DelHostInApp GetConfigByCond  msg : %v", msg)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)

		}
		moduleIdArr := make([]int, 0)
		for _, item := range configArr {
			moduleIdArr = append(moduleIdArr, item["ModuleID"])
		}
		moduleCon := map[string]interface{}{
			"ModuleID": map[string]interface{}{
				common.BKDBIN: moduleIdArr,
			},
			"Default": common.DefaultResModuleFlag,
		}
		fields := "ModuleID"
		moduleArr, err := logics.GetModuleMapByCond(req, fields, cli.CC.ObjCtrl(), moduleCon)
		if err != nil {
			blog.Error("DelHostInApp GetConfigByCond err msg : %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}
		blog.Debug("moduleArr:%v", moduleArr)
		if len(moduleArr) == 0 {
			msg := fmt.Sprintf("非空闲主机不能删除")
			blog.Debug("DelHostInApp GetModuleMapByCond  msg : %v", msg)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}
		param := make(common.KvMap)
		param["ApplicationID"] = appId
		param["HostID"] = hostId
		uUrl := cli.CC.ObjCtrl() + "/object/v1/openapi/set/delhost"
		blog.Debug("uUrl%v", uUrl)
		inputJson, err := json.Marshal(param)
		blog.Debug("inputJson%v", string(inputJson))

		if nil != err {
			blog.Error("Marshal json error:%v", err)
			return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONMarshalFailed)

		}

		res, err := httpcli.ReqHttp(req, uUrl, common.HTTPDelete, []byte(inputJson))
		blog.Debug("del res:%v", res)
		if nil != err {
			blog.Error("request ctrl error:%v", err)
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommHTTPDoRequestFailed)

		}
		blog.Debug("res:%v", res)
		//err = delSetConfigHost(param)
		var rst api.BKAPIRsp
		if "not found" == fmt.Sprintf("%v", err) {
			return http.StatusOK, &rst, nil
		}
		if nil != err {
			blog.Error("delSetConfigHost error:%v", err)
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoSetDeleteFailed)

		}

		// deal result
		return http.StatusOK, &rst, nil
	}, resp)
}

func (cli *hostAction) GetGitServerIp(req *restful.Request, resp *restful.Response) {
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	cli.CallResponseEx(func() (int, interface{}, error) {
		blog.Debug("GetGitServerIp start")
		value, _ := ioutil.ReadAll(req.Request.Body)
		js, err := simplejson.NewJson([]byte(value))
		if err != nil {
			blog.Error("GetGitServerIp failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}
		input, err := js.Map()
		if err != nil {
			blog.Error("GetGitServerIp failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
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
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}
		if len(appMap) == 0 {
			// cli.ResponseSuccess("", resp)
			return http.StatusOK, nil, nil
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
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}
		if len(setMap) == 0 {
			// cli.ResponseSuccess("", resp)
			return http.StatusOK, nil, nil
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
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}
		if len(moduleMap) == 0 {
			cli.ResponseSuccess(nil, resp)
			return http.StatusOK, nil, nil
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
		hostArr, err := getHostDataByConfig(req, configData)
		if nil != err {
			blog.Error("GetGitServerIp getHostDataByConfig err msg : %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}
		blog.Debug("hostArr:%v", hostArr)

		// cli.ResponseSuccess(hostArr, resp)
		return http.StatusOK, hostArr, nil

	}, resp)
}

// helpers
func updateHostMain(req *restful.Request, hostCondition, data map[string]interface{}, appID int, hostCtrl, objCtrl, auditCtrl string, errIf errorIfs.CCErrorIf) (string, error) {
	blog.Debug("updateHostMain start")
	blog.Debug("hostCondition:%v", hostCondition)
	_, hostIDArr, err := getHostMapByCond(req, hostCondition)

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
	valid := validator.NewValidMapWithKeyFields(common.BKDefaultOwnerID, common.BKInnerObjIDHost, objCtrl, []string{common.CreateTimeField, common.LastTimeField, common.BKChildStr}, forward, errIf.CreateDefaultCCErrorIf(language))
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

func getDefaultModules(req *restful.Request, appID int, objURL string) (interface{}, error) {

	param := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKAppIDField:   appID,
			common.BKDefaultField: 1,
		},
		"fields": fmt.Sprintf("%s,%s", common.BKSetIDField, common.BKModuleIDField),
	}

	resMap, err := getObjByCondition(req, param, common.BKInnerObjIDModule, objURL)

	if nil != err {
		return nil, err
	}

	blog.Debug("getDefaultModules complete, res: %v", resMap)

	if !resMap["result"].(bool) {
		return nil, errors.New(resMap["message"].(string))
	}

	resDataMap := resMap["data"].(map[string]interface{})

	if resDataMap["count"] == 0 {
		return nil, errors.New(fmt.Sprintf("can not found default module, appid: %d", appID))
	}

	return (resDataMap["info"].([]interface{}))[0], nil

}

func getHostByIPAndSource(req *restful.Request, innerIP string, platID int, objURL string) (interface{}, error) {

	param := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKHostInnerIPField: innerIP,
			common.BKCloudIDField:     platID,
		},
		"fields": common.BKHostIDField,
	}

	resMap, err := getObjByCondition(req, param, common.BKInnerObjIDHost, objURL)

	if nil != err {
		return nil, err
	}

	if !resMap["result"].(bool) {
		return nil, errors.New(resMap["message"].(string))
	}

	resDataMap := resMap["data"].(map[string]interface{})

	blog.Debug("getHostByIPAndSource res: %v", resDataMap)

	return resDataMap["info"], nil
}

func getHostByCond(req *restful.Request, param map[string]interface{}, objURL string) (interface{}, error) {

	blog.Debug("param:%v", param)
	resMap, err := getObjByCondition(req, param, common.BKInnerObjIDHost, objURL)

	if nil != err {
		return nil, err
	}

	if !resMap["result"].(bool) {
		return nil, errors.New(resMap["message"].(string))
	}

	resDataMap := resMap["data"].(map[string]interface{})

	blog.Debug("getHostByIPArrAndSource res: %v", resDataMap)

	return resDataMap["info"], nil
}

func getObjByCondition(req *restful.Request, param map[string]interface{}, objType, objURL string) (map[string]interface{}, error) {
	resMap := make(map[string]interface{})

	url := objURL + "/object/v1/insts/" + objType + "/search"
	inputJson, _ := json.Marshal(param)
	res, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(inputJson))
	if nil != err {
		return nil, err
	}

	err = json.Unmarshal([]byte(res), &resMap)
	if nil != err {
		return nil, err
	}

	return resMap, nil
}

func addHost(req *restful.Request, data map[string]interface{}, objURL string) (int, error) {
	return addObj(req, data, common.BKInnerObjIDHost, objURL)
}

func addModuleHostConfig(req *restful.Request, data map[string]interface{}, hostCtrl string) error {
	blog.Debug("addModuleHostConfig start, data: %v", data)

	resMap := make(map[string]interface{})
	inputJson, _ := json.Marshal(data)
	addModulesURL := hostCtrl + "/host/v1/meta/hosts/modules"
	res, err := httpcli.ReqHttp(req, addModulesURL, common.HTTPCreate, []byte(inputJson))
	if nil != err {
		return err
	}

	json.Unmarshal([]byte(res), &resMap)

	if !resMap["result"].(bool) {
		return errors.New(resMap["message"].(string))
	}
	blog.Debug("addModuleHostConfig success, res: %v", resMap)
	return nil
}

func addObj(req *restful.Request, data map[string]interface{}, objType, objURL string) (int, error) {
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
		return 0, errors.New(resMap["message"].(string))
	}

	blog.Debug("add object result : %v", resMap)

	objID := (resMap["data"].(map[string]interface{}))[common.BKHostIDField].(float64)
	return int(objID), nil
}

//search host helpers

func setHostData(req *restful.Request, moduleHostConfig []map[string]int, hostMap map[int]interface{}) ([]interface{}, error) {

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

func getHostMapByCond(req *restful.Request, condition map[string]interface{}) (map[int]interface{}, []int, error) {
	hostMap := make(map[int]interface{})
	hostIDArr := make([]int, 0)

	// build host controller url
	url := host.CC.HostCtrl() + "/host/v1/hosts/search"
	searchParams := map[string]interface{}{
		"fields":    "",
		"condition": condition,
	}
	inputJson, err := json.Marshal(searchParams)
	if nil != err {
		return nil, nil, err
	}
	hostInfo, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(inputJson))
	blog.Debug("appInfo:%v", hostInfo)
	if nil != err {
		blog.Errorf("getHostMapByCond error:%s, params:%s, error:%s", url, string(inputJson), err.Error())
		return hostMap, hostIDArr, err
	}

	js, err := simplejson.NewJson([]byte(hostInfo))
	if nil != err {
		return nil, nil, err
	}

	resDataInfo, err := js.Get("data").Get("info").Array() //res["data"].(map[string]interface{})
	if nil != err {
		return nil, nil, err
	}

	for _, item := range resDataInfo {
		host := item.(map[string]interface{})
		host_id, err := util.GetIntByInterface(host[common.BKHostIDField])
		if nil != err {
			return nil, nil, err
		}

		hostMap[host_id] = host
		hostIDArr = append(hostIDArr, host_id)
	}
	return hostMap, hostIDArr, nil
}

func getHostDataByConfig(req *restful.Request, configData []map[string]int) ([]interface{}, error) {

	hostIDArr := make([]int, 0)

	for _, config := range configData {
		hostIDArr = append(hostIDArr, config[common.BKHostIDField])
	}

	hostMapCondition := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDArr,
		},
	}

	hostMap, _, err := getHostMapByCond(req, hostMapCondition)
	if nil != err {
		return nil, err
	}

	hostData, err := setHostData(req, configData, hostMap)
	if nil != err {
		return hostData, err
	}

	return hostData, nil
}

func getCustomerPropertyByOwner(req *restful.Request, OwnerId interface{}, ObjCtrl string) ([]map[string]interface{}, error) {
	blog.Debug("getCustomerPropertyByOwner start")
	gHostAttrUrl := ObjCtrl + "/object/v1/meta/objectatts"
	searchBody := make(map[string]interface{})
	searchBody[common.BKObjIDField] = common.BKInnerObjIDHost
	searchBody[common.BKOwnerIDField] = OwnerId
	searchJson, _ := json.Marshal(searchBody)
	gHostAttrRe, err := httpcli.ReqHttp(req, gHostAttrUrl, common.HTTPSelectPost, []byte(searchJson))
	if nil != err {
		blog.Error("GetHostDetailById  attr error :%v", err)
		return nil, err
	}
	js, err := simplejson.NewJson([]byte(gHostAttrRe))
	gHostAttr, _ := js.Map()

	gAttrResult := gHostAttr["result"].(bool)
	if false == gAttrResult {
		blog.Error("GetHostDetailById  attr error :%v", err)
		return nil, err
	}
	hostAttrArr := gHostAttr["data"].([]interface{})
	customAttrArr := make([]map[string]interface{}, 0)
	for _, attr := range hostAttrArr {
		if !attr.(map[string]interface{})[common.BKIsPre].(bool) {
			customAttrArr = append(customAttrArr, attr.(map[string]interface{}))
		}
	}
	return customAttrArr, nil
}
