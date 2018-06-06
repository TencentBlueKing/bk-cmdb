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
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/logics"
	phpapilogic "configcenter/src/scene_server/host_server/host_service/logics/phpapi"
)

var host *hostAction = &hostAction{}

type hostAction struct {
	base.BaseAction
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/openapi/host/{" + common.BKAppIDField + "}", Params: nil, Handler: host.UpdateHost})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/host/updateHostByAppID/{appid}", Params: nil, Handler: host.UpdateHostByAppID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/host/getHostListByAppidAndField/{" + common.BKAppIDField + "}/{field}", Params: nil, Handler: host.GetHostListByAppidAndField})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/openapi/gethostlistbyip", Params: nil, Handler: host.HostSearchByConds})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/gethostlistbyconds", Params: nil, Handler: host.HostSearchByConds})
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

// updateHostPlat 根据条件更新主机信息
func (cli *hostAction) UpdateHost(req *restful.Request, resp *restful.Response) {
	blog.Debug("updateHost start!")

	appID, err := strconv.Atoi(req.PathParameter(common.BKAppIDField))
	if nil != err {
		blog.Error("convert appid to int error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}

	blog.Debug("updateHost http body data: %s", value)

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Error("unmarshal json error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}
	blog.Debug("input:%s", input, string(value))

	updateData, ok := input["data"]
	if !ok {
		blog.Error("params data is required:%s", string(value))
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}
	mapData, ok := updateData.(map[string]interface{})
	if !ok {
		blog.Error("params data must be object:%s", string(value))
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	dstPlat, ok := mapData[common.BKSubAreaField]
	if !ok {
		blog.Error("params data.bk_cloud_id is require:%s", string(value))
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	// dst host exist return souccess, hongsong tiyi
	dstHostCondition := map[string]interface{}{
		common.BKHostInnerIPField: input["condition"].(map[string]interface{})[common.BKHostInnerIPField],
		common.BKCloudIDField:     dstPlat,
	}
	_, hostIDArr, err := phpapilogic.GetHostMapByCond(req, dstHostCondition)
	blog.Debug("hostIDArr:%v", hostIDArr)
	if nil != err {
		blog.Error("updateHostMain error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_HOST_MODIFY_FAIL, err.Error(), resp)
		return
	}

	if len(hostIDArr) != 0 {
		cli.ResponseSuccess(nil, resp)
		return
	}

	blog.Debug(input["condition"].(map[string]interface{})[common.BKCloudIDField])
	hostCondition := map[string]interface{}{
		common.BKHostInnerIPField: input["condition"].(map[string]interface{})[common.BKHostInnerIPField],
		common.BKCloudIDField:     input["condition"].(map[string]interface{})[common.BKCloudIDField],
	}
	data := input["data"].(map[string]interface{})
	data[common.BKHostInnerIPField] = input["condition"].(map[string]interface{})[common.BKHostInnerIPField]
	res, err := phpapilogic.UpdateHostMain(req, hostCondition, data, appID, cli.CC.HostCtrl(), cli.CC.ObjCtrl(), cli.CC.AuditCtrl(), cli.CC.Error)

	if nil != err {
		blog.Error("updateHostMain error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_HOST_MODIFY_FAIL, err.Error(), resp)
		return
	}

	cli.ResponseSuccess(res, resp)
	return
}

// updateHostByAppID 根据IP更新主机Proxy状态，如果不存在主机则添加到对应业务及默认模块
func (cli *hostAction) UpdateHostByAppID(req *restful.Request, resp *restful.Response) {
	blog.Debug("updateHostByAppID start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	appID, err := strconv.Atoi(req.PathParameter("appid"))
	if nil != err {
		blog.Error("convert appid to int error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}

	blog.Debug("updateHostByAppID http body data: %s", value)

	input := make(map[string]interface{})
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Error("unmarshal json error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	proxyArr := input[common.BKProxyListField].([]interface{})
	platID, _ := util.GetIntByInterface(input[common.BKCloudIDField])

	blog.Debug("proxyArr:%v", proxyArr)
	defaultModule, err := phpapilogic.GetDefaultModules(req, appID, cli.CC.ObjCtrl())

	if nil != err {
		blog.Error("getDefaultModules error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	//defaultSetID := defaultModule["SetID"]
	defaultModuleMap, _ := defaultModule.(map[string]interface{})
	if nil != err {
		blog.Error("getDefaultModules error:%v", err)
		ccErr := defErr.Error(common.CCErrGetModule)
		cli.ResponseFailed(common.CCErrGetModule, ccErr.Error(), resp)
		return
	}
	defaultModuleID, err := util.GetInt64ByInterface(defaultModuleMap[common.BKModuleIDField])
	if nil != err {
		blog.Error("getDefaultModules error:%v", err)
		ccErr := defErr.Error(common.CCErrGetModule)
		cli.ResponseFailed(common.CCErrGetModule, ccErr.Error(), resp)
		return
	}
	for _, pro := range proxyArr {
		proMap := pro.(map[string]interface{})
		var hostID int
		innerIP := proMap[common.BKHostInnerIPField]
		outerIP, ok := proMap[common.BKHostOuterIPField]
		if !ok {
			outerIP = ""
		}

		hostData, err := phpapilogic.GetHostByIPAndSource(req, innerIP.(string), platID, cli.CC.ObjCtrl())
		blog.Error("hostData:%v", hostData)
		if nil != err {
			blog.Error("getHostByIPAndSource error:%v", err)
			cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
			return
		}

		hostDataArr := hostData.([]interface{})

		if len(hostDataArr) == 0 {
			platID, ok := proMap[common.BKCloudIDField]
			if ok {
				platConds := common.KvMap{
					common.BKCloudIDField: platID,
				}
				bl, err := logics.IsExistPlat(req, cli.CC.ObjCtrl(), platConds)
				if nil != err {
					blog.Errorf("is exist plat  error:%s", err.Error())
					cli.ResponseFailed(common.CCErrTopoGetCloudErrStrFaild, defErr.Errorf(common.CCErrTopoGetCloudErrStrFaild, err.Error()).Error(), resp)
					return
				}
				if !bl {
					blog.Errorf("is exist plat  not foud platid :%v", platID)
					cli.ResponseFailed(common.CCErrTopoCloudNotFound, defErr.Error(common.CCErrTopoCloudNotFound).Error(), resp)
					return
				}
			}
			blog.Debug("procMap:%v", proMap)
			proMap["import_from"] = common.HostAddMethodAgent
			hostIDNew, err := phpapilogic.AddHost(req, proMap, cli.CC.ObjCtrl())

			if nil != err {
				blog.Error("addHost error:%v", err)
				cli.ResponseFailed(common.CC_Err_Comm_Host_Update_FAIL_ERR, common.CC_Err_Comm_Host_Update_FAIL_ERR_STR, resp)
				return
			}

			hostID = hostIDNew

			blog.Debug("addHost success, hostID: %d", hostID)

			err = phpapilogic.AddModuleHostConfig(req, map[string]interface{}{
				common.BKAppIDField:    appID,
				common.BKModuleIDField: []int64{defaultModuleID},
				common.BKHostIDField:   hostID,
			}, cli.CC.HostCtrl())

			if nil != err {
				blog.Error("addModuleHostConfig error:%v", err)
				cli.ResponseFailed(common.CC_Err_Comm_Host_Update_FAIL_ERR, common.CC_Err_Comm_Host_Update_FAIL_ERR_STR, resp)
				return
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

			_, err := phpapilogic.UpdateHostMain(req, hostCondition, data, appID, cli.CC.HostCtrl(), cli.CC.ObjCtrl(), cli.CC.AuditCtrl(), cli.CC.Error)
			if nil != err {
				blog.Error("updateHostMain error:%v", err)
				cli.ResponseFailed(common.CC_Err_Comm_Host_Update_FAIL_ERR, err.Error(), resp)
				return
			}
		}

	}

	cli.ResponseSuccess(nil, resp)
}

// updateCustomProperty 修改主机自定义属性
func (cli *hostAction) UpdateCustomProperty(req *restful.Request, resp *restful.Response) {
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
	blog.Error("UpdateCustomProperty :%v", input)
	appId, _ := strconv.Atoi(input[common.BKAppIDField].(string))
	hostId, _ := strconv.Atoi(input[common.BKHostIDField].(string))
	propertyJson := input["property"]

	propertyMap := make(map[string]interface{})
	if nil != propertyJson {
		err = json.Unmarshal([]byte(propertyJson.(string)), &propertyMap)
	}
	if nil != err {
		blog.Error("Unmarshal json failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_APP_Update_FAIL, common.CC_Err_Comm_APP_Update_FAIL_STR, resp)
		return
	}
	condition := make(common.KvMap)
	condition[common.BKAppIDField] = appId
	fileds := fmt.Sprintf("%s,%s", common.BKAppIDField, common.BKOwnerIDField)
	apps, err := logics.GetAppMapByCond(req, fileds, cli.CC.ObjCtrl(), condition)
	if nil != err {
		blog.Error("UpdateCustomProperty GetAppMapByCond, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_APP_Update_FAIL, common.CC_Err_Comm_APP_Update_FAIL_STR, resp)
		return
	}
	blog.Debug("UpdateCustomProperty apps:%v", apps)
	if _, ok := apps[appId]; !ok {
		msg := "业务不存在"
		blog.Debug("UpdateCustomProperty error:%v", msg)
		cli.ResponseFailed(common.CC_Err_Comm_APP_Update_FAIL, msg, resp)
		return
	}

	appMap := apps[appId]
	ownerId := appMap.(map[string]interface{})[common.BKOwnerIDField]
	propertys, _ := phpapilogic.GetCustomerPropertyByOwner(req, ownerId, cli.CC.ObjCtrl())
	params := make(common.KvMap)
	for _, attrMap := range propertys {
		PropertyId, ok := attrMap[common.BKPropertyIDField].(string)
		if !ok {
			continue
		}
		blog.Debug("input[PropertyId]:%v", input[PropertyId])
		if _, ok := propertyMap[PropertyId]; ok {
			params[PropertyId] = propertyMap[PropertyId]
		}
	}
	blog.Debug("params:%v", params)
	hostCondition := map[string]interface{}{
		common.BKHostIDField: hostId,
	}
	res, err := phpapilogic.UpdateHostMain(req, hostCondition, params, appId, cli.CC.HostCtrl(), cli.CC.ObjCtrl(), cli.CC.AuditCtrl(), cli.CC.Error)
	if nil != err {
		msg := fmt.Sprintf("%v", err)
		blog.Error("UpdateCustomProperty updateHostMain error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_APP_Update_FAIL, msg, resp)
		return
	}
	cli.ResponseSuccess(res, resp)
}

// cloneHostProperty 克隆主机
func (cli *hostAction) CloneHostProperty(req *restful.Request, resp *restful.Response) {
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
	blog.Debug("CloneHostProperty input:%v", input)
	appId, _ := strconv.Atoi(input[common.BKAppIDField].(string))
	orgIp := input[common.BKOrgIPField]
	dstIp := input[common.BKDstIPField]

	platId, hasPlatId := input[common.BKCloudIDField]
	platIdInt, _ := strconv.Atoi(input[common.BKCloudIDField].(string))
	condition := common.KvMap{
		common.BKHostInnerIPField: orgIp,
	}

	if hasPlatId && platId != nil && platId != "" {
		condition[common.BKCloudIDField] = platIdInt
	}
	// 处理源IP
	hostMap, hostIdArr, err := phpapilogic.GetHostMapByCond(req, condition)

	blog.Debug("hostMapData:%v", hostMap)
	if err != nil {
		blog.Error("getHostMapByCond error : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return
	}

	if len(hostIdArr) == 0 {
		blog.Error("clone host getHostMapByCond error, ip:%s, platid:%s", orgIp, platId)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return
	}

	hostMapData, ok := hostMap[hostIdArr[0]].(map[string]interface{})
	if false == ok {
		blog.Error("getHostMapByCond not source ip : %s", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return
	}

	configCond := map[string]interface{}{
		common.BKHostIDField: []interface{}{hostMapData[common.BKHostIDField]},
		common.BKAppIDField:  []int{appId},
	}
	// 判断源IP是否存在
	configData, err := logics.GetConfigByCond(req, host.CC.HostCtrl(), configCond)
	blog.Debug("configData:%v", configData)
	if nil != err {
		blog.Error("clone host property error : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return
	}
	if len(configData) == 0 {
		msg := "no find host"
		blog.Error("clone host property error : %v", msg)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, fmt.Sprintf("%s %s", common.CC_Err_Comm_Host_Get_FAIL_STR, msg), resp)
		return
	}
	// 处理目标IP
	dstIpArr := strings.Split(dstIp.(string), ",")
	// 获得已存在的主机
	dstCondition := map[string]interface{}{
		common.BKHostInnerIPField: map[string]interface{}{
			common.BKDBIN: dstIpArr,
		},
		common.BKCloudIDField: platIdInt,
	}
	dstHostMap, dstHostIdArr, err := phpapilogic.GetHostMapByCond(req, dstCondition)
	blog.Debug("dstHostMap:%v", dstHostMap)

	dstConfigCond := map[string]interface{}{
		common.BKAppIDField:  []int{appId},
		common.BKHostIDField: dstHostIdArr,
	}
	dstHostIdArrV, err := logics.GetHostIDByCond(req, host.CC.HostCtrl(), dstConfigCond)
	existIpArr := make([]string, 0)
	for _, id := range dstHostIdArrV {
		if dstHostMapData, ok := dstHostMap[id].(map[string]interface{}); ok {
			existIpArr = append(existIpArr, dstHostMapData[common.BKHostInnerIPField].(string))
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
	blog.Debug("existIpArr:%v", existIpArr)
	for _, dstIpV := range dstIpArr {
		if dstIpV == orgIp {
			blog.Debug("clone host updateHostMain err:%v", err)
			msg := "dstIp 和 orgIp不能相同"
			cli.ResponseFailed(common.CC_Err_Comm_Host_Update_FAIL_ERR, fmt.Sprintf("%s%s", common.CC_Err_Comm_Host_Update_FAIL_ERR_STR, msg), resp)
			return
		}
		blog.Debug("hostMapData:%v", hostMapData)
		if phpapilogic.In_existIpArr(existIpArr, dstIpV) {
			blog.Debug("clone update")
			hostCondition := map[string]interface{}{
				common.BKHostInnerIPField: dstIpV,
			}

			updateHostData[common.BKHostInnerIPField] = dstIpV
			delete(updateHostData, common.BKHostIDField)
			res, err := phpapilogic.UpdateHostMain(req, hostCondition, updateHostData, appId, host.CC.HostCtrl(), host.CC.ObjCtrl(), host.CC.AuditCtrl(), cli.CC.Error)
			if nil != err {
				blog.Debug("clone host updateHostMain err:%v", err)
				msg := fmt.Sprintf("clone host error:%s", dstIpV)
				cli.ResponseFailed(common.CC_Err_Comm_Host_Update_FAIL_ERR, fmt.Sprintf("%s%s", common.CC_Err_Comm_Host_Update_FAIL_ERR_STR, msg), resp)
				return
			}
			blog.Debug("clone host updateHostMain res:%v", res)
		} else {
			hostMapData[common.BKHostInnerIPField] = dstIpV
			blog.Debug("clone add")
			addHostMapData := hostMapData
			delete(addHostMapData, common.BKHostIDField)
			cloneHostId, err := phpapilogic.AddHost(req, addHostMapData, host.CC.ObjCtrl())
			if nil != err {
				blog.Debug("clone host addHost err:%v", err)
				msg := fmt.Sprintf("clone host error:%s", dstIpV)
				cli.ResponseFailed(common.CC_Err_Comm_HOST_CREATE_FAIL, fmt.Sprintf("%s%s", common.CC_Err_Comm_HOST_CREATE_FAIL_STR, msg), resp)
				return
			}

			blog.Debug("cloneHostId:%v", cloneHostId)
			blog.Debug("configData[0]:%v", configData[0])
			configDataMap := make(map[string]interface{}, 0)
			configDataMap[common.BKHostIDField] = cloneHostId
			configDataMap[common.BKModuleIDField] = []int{configData[0][common.BKModuleIDField]}
			configDataMap[common.BKAppIDField] = configData[0][common.BKAppIDField]
			configDataMap[common.BKSetIDField] = configData[0][common.BKSetIDField]
			err = phpapilogic.AddModuleHostConfig(req, configDataMap, host.CC.HostCtrl())
			if nil != err {
				blog.Debug("clone host addModuleHostConfig err:%v", err)
				msg := fmt.Sprintf("clone host error:%s", dstIpV)
				cli.ResponseFailed(common.CC_Err_Comm_HOST_CREATE_FAIL, fmt.Sprintf("%s%s", common.CC_Err_Comm_HOST_CREATE_FAIL_STR, msg), resp)
				return
			}
		}
	}

	cli.ResponseSuccess(nil, resp)
}

//DelHostInApp: 从业务空闲机集群中删除主机
func (cli *hostAction) DelHostInApp(req *restful.Request, resp *restful.Response) {
	value, _ := ioutil.ReadAll(req.Request.Body)
	js, err := simplejson.NewJson([]byte(value))
	if err != nil {
		blog.Error("DelHostInApp failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}
	input, err := js.Map()
	if err != nil {
		blog.Error("DelHostInApp failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
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
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL, resp)
		return
	}
	if len(configArr) == 0 {
		msg := fmt.Sprintf("not fint hostId:%v in appId:%v", hostId, appId)
		blog.Info("DelHostInApp GetConfigByCond  msg : %v", msg)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, fmt.Sprintf("%s:%s", common.CC_Err_Comm_Host_Get_FAIL_STR, msg), resp)
		return
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
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL, resp)
		return
	}
	blog.Debug("moduleArr:%v", moduleArr)
	if len(moduleArr) == 0 {
		msg := fmt.Sprintf("非空闲主机不能删除")
		blog.Debug("DelHostInApp GetModuleMapByCond  msg : %v", msg)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, msg, resp)
		return
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
		return
	}

	res, err := httpcli.ReqHttp(req, uUrl, common.HTTPDelete, []byte(inputJson))
	blog.Debug("del res:%v", res)
	if nil != err {
		blog.Error("request ctrl error:%v", err)
		return
	}
	blog.Debug("res:%v", res)
	//err = delSetConfigHost(param)
	var rst api.BKAPIRsp
	if "not found" == fmt.Sprintf("%v", err) {
		cli.Response(&rst, resp)
		return
	}
	if nil != err {
		blog.Error("delSetConfigHost error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Set_Delete_FAIL, common.CC_Err_Comm_Set_Delete_FAIL, resp)
		return
	}

	// deal result

	cli.Response(&rst, resp)
}
