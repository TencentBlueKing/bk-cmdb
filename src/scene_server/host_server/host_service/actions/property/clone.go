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

package property

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/logics"
	apilogic "configcenter/src/scene_server/host_server/host_service/logics/phpapi"

	"github.com/emicklei/go-restful"
)

var host *hostAction = &hostAction{}

type hostAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/propery/clone", Params: nil, Handler: host.CloneHostProperty})

	// create CC object
	host.CreateAction()
}

// cloneHostProperty  clone host property
func (cli *hostAction) CloneHostProperty(req *restful.Request, resp *restful.Response) {
	defError := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CCErrCommHTTPReadBodyFailed, defError.Error(common.CCErrCommHTTPReadBodyFailed).Error(), resp)
		return
	}

	var input metadata.HostCloneInputParams
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Error("Unmarshal json failed, error:%v", err)
		cli.ResponseFailed(common.CCErrCommJSONUnmarshalFailed, defError.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}

	blog.Debug("CloneHostProperty input:%v", input)

	condition := common.KvMap{
		common.BKHostInnerIPField: input.OrgIP,
		common.BKCloudIDField:     input.PlatID,
	}

	// deal with origin IP
	hostMap, hostIdArr, err := apilogic.GetHostMapByCond(req, condition)

	blog.Debug("hostMapData:%v", hostMap)
	if err != nil {
		blog.Error("getHostMapByCond error : %v", err)
		cli.ResponseFailed(common.CCErrHostDetailFail, err.Error(), resp)
		return
	}

	if len(hostIdArr) == 0 {
		blog.Error("clone host getHostMapByCond error, ip:%s, platid:%s", input.OrgIP, input.PlatID)
		cli.ResponseFailed(common.CCErrHostDetailFail, "not found host ", resp)
		return
	}

	hostMapData, ok := hostMap[hostIdArr[0]].(map[string]interface{})
	if false == ok {
		blog.Error("getHostMapByCond not source ip , raw data format error: %v", hostMap)
		cli.ResponseFailed(common.CCErrHostDetailFail, "source ip not found", resp)
		return
	}

	hostIDI, ok := hostMapData[common.BKHostIDField].(int64)
	if false == ok {
		blog.Error("host id not int : %v", hostMapData[common.BKHostIDField])
		cli.ResponseFailed(common.CCErrHostDetailFail, "source ip not found", resp)
		return
	}

	configCond := map[string]interface{}{
		common.BKHostIDField: []int64{hostIDI},
		common.BKAppIDField:  []int64{input.AppID},
	}
	// check is ip exist
	configData, err := logics.GetConfigByCond(req, host.CC.HostCtrl(), configCond)
	blog.Debug("configData:%v", configData)
	if nil != err {
		blog.Error("clone host property error : %v", err)
		cli.ResponseFailed(common.CCErrHostDetailFail, err.Error(), resp)
		return
	}
	if len(configData) == 0 {
		msg := "no find host module relation "
		blog.Error("clone host property error : %v", msg)
		cli.ResponseFailed(common.CCErrHostDetailFail, msg, resp)
		return
	}
	// deal with destination IP
	dstIpArr := strings.Split(input.DstIP, ",")

	// get exist ip
	dstCondition := map[string]interface{}{
		common.BKHostInnerIPField: map[string]interface{}{
			common.BKDBIN: dstIpArr,
		},
		common.BKCloudIDField: input.PlatID,
	}
	dstHostMap, dstHostIdArr, err := apilogic.GetHostMapByCond(req, dstCondition)
	blog.Debug("dstHostMap:%v", dstHostMap)

	dstConfigCond := map[string]interface{}{
		common.BKAppIDField:  []int64{input.AppID},
		common.BKHostIDField: dstHostIdArr,
	}
	dstHostIdArrV, err := logics.GetHostIDByCond(req, host.CC.HostCtrl(), dstConfigCond)
	existIPMap := make(map[string]int64, 0)
	for _, id := range dstHostIdArrV {
		if dstHostMapData, ok := dstHostMap[id].(map[string]interface{}); ok {
			ip, ok := dstHostMapData[common.BKHostInnerIPField].(string)
			if false == ok {
				cli.ResponseFailed(common.CCErrHostDetailFail, "data format error, not found innerip", resp)
				return
			}

			hostID, err := util.GetInt64ByInterface(dstHostMapData[common.BKHostIDField])
			if nil != err {
				cli.ResponseFailed(common.CCErrHostDetailFail, "data format error, not found host id", resp)
				return
			}
			existIPMap[ip] = hostID
		} else {
			cli.ResponseFailed(common.CCErrHostDetailFail, "data format error", resp)
			return
		}
	}

	//do not update nil data
	updateHostData := make(map[string]interface{})
	for key, val := range hostMapData {
		if nil != val {
			updateHostData[key] = val
		}
	}
	// remote duplication ip
	dstIPMap := make(map[string]bool, len(dstIpArr))
	for _, ip := range dstIpArr {
		dstIPMap[ip] = true
	}

	blog.Debug("configData[0]:%v", configData[0])
	moduleIDs := make([]int64, 0)
	for _, configData := range configData {

		moduleID, err := util.GetInt64ByInterface(configData[common.BKModuleIDField])
		if nil != err {
			cli.ResponseFailed(common.CCErrGetOriginHostModuelRelationship, fmt.Sprintf("get source ip module error, error data:%v", configData), resp)
			return
		}
		moduleIDs = append(moduleIDs, moduleID)
	}

	// clone host, existing modification, new addition;
	// dstIpArr: all hosts to be cloned, existIpArr: existing host to be cloned
	blog.Debug("existIpArr:%v", existIPMap)
	for dstIpV, _ := range dstIPMap {
		if dstIpV == input.OrgIP {
			blog.Debug("clone host updateHostMain err:dstIp and orgIp cannot be the same")
			msg := "dstIp and orgIp cannot be the same"
			cli.ResponseFailed(common.CCErrHostCreateFail, msg, resp)
			return
		}
		blog.Debug("hostMapData:%v", hostMapData)
		hostID, oK := existIPMap[dstIpV]
		if true == oK {
			blog.Debug("clone update")
			hostCondition := map[string]interface{}{
				common.BKHostInnerIPField: dstIpV,
			}

			updateHostData[common.BKHostInnerIPField] = dstIpV
			delete(updateHostData, common.BKHostIDField)
			res, err := apilogic.UpdateHostMain(req, hostCondition, updateHostData, int(input.AppID), host.CC.HostCtrl(), host.CC.ObjCtrl(), host.CC.AuditCtrl(), cli.CC.Error)
			if nil != err {
				blog.Debug("clone host updateHostMain err: %v", err)
				msg := fmt.Sprintf("clone host error:%s", dstIpV)
				cli.ResponseFailed(common.CC_Err_Comm_Host_Update_FAIL_ERR, fmt.Sprintf("%s%s", common.CC_Err_Comm_Host_Update_FAIL_ERR_STR, msg), resp)
				return
			}
			blog.Debug("clone host updateHostMain res:%v", res)
			params := make(map[string]interface{})
			params[common.BKAppIDField] = input.AppID
			params[common.BKHostIDField] = hostID

			delModulesURL := cli.CC.HostCtrl() + "/host/v1/meta/hosts/modules"
			isSuccess, errMsg, _ := logics.GetHttpResult(req, delModulesURL, common.HTTPDelete, params)
			if !isSuccess {
				blog.Error("remove hosthostconfig error, params:%v, error:%s", params, errMsg)
				cli.ResponseFailed(common.CCErrHostTransferModule, fmt.Sprintf("remote host module config error, error:%s", errMsg), resp)
				return
			}
		} else {
			hostMapData[common.BKHostInnerIPField] = dstIpV
			blog.Debug("clone add")
			addHostMapData := hostMapData
			delete(addHostMapData, common.BKHostIDField)
			cloneHostId, err := apilogic.AddHost(req, addHostMapData, host.CC.ObjCtrl())
			if nil != err {
				blog.Debug("clone host addHost err:%v", err)
				msg := fmt.Sprintf("clone host error:%s", dstIpV)
				cli.ResponseFailed(common.CC_Err_Comm_HOST_CREATE_FAIL, fmt.Sprintf("%s%s", common.CC_Err_Comm_HOST_CREATE_FAIL_STR, msg), resp)
				return
			}

			blog.Debug("cloneHostId:%v", cloneHostId)
			hostID = cloneHostId

		}
		err = apilogic.AddModuleHostConfig(req, hostID, input.AppID, moduleIDs, host.CC.HostCtrl())
		if nil != err {
			blog.Debug("clone host addModuleHostConfig err:%v", err)
			msg := fmt.Sprintf("clone host get source ip host module relation failure, error:%s", err.Error())
			cli.ResponseFailed(common.CCErrHostModuleRelationAddFailed, msg, resp)
			return
		}
	}

	cli.ResponseSuccess(nil, resp)
}
