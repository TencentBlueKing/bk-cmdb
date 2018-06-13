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

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/scene_server/host_server/host_service/logics"
	apilogic "configcenter/src/scene_server/host_server/host_service/logics/phpapi"
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

type InputParams struct {
	OrgIP  string `json:"bk_org_ip"`
	DstIP  string `json:"bk_dst_ip"`
	AppID  int    `json:"bk_biz_id"`
	PlatID int    `json:"bk_cloud_id"`
}

// CloneHostProperty  clone host property and host module config
func (cli *hostAction) CloneHostProperty(req *restful.Request, resp *restful.Response) {
	value, err := ioutil.ReadAll(req.Request.Body)
	if nil != err {
		blog.Error("read request body failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_DO_STR, resp)
		return
	}
	var input InputParams
	err = json.Unmarshal(value, &input)
	if nil != err {
		blog.Error("Unmarshal json failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	blog.Debug("CloneHostProperty input:%v", input)
	appID := input.AppID
	orgIP := input.OrgIP
	dstIP := input.DstIP
	platID := input.PlatID

	condition := common.KvMap{
		common.BKHostInnerIPField: orgIP,
		common.BKCloudIDField:     platID,
	}

	// deal with origin IP
	hostMap, hostIdArr, err := apilogic.GetHostMapByCond(req, condition)

	blog.Debug("hostMapData:%v", hostMap)
	if err != nil {
		blog.Error("getHostMapByCond error : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return
	}

	if len(hostIdArr) == 0 {
		blog.Error("clone host getHostMapByCond error, ip:%s, platid:%s", orgIP, platID)
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
		common.BKAppIDField:  []int{appID},
	}

	// check origin ip is exist
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

	// deal with destination IP
	dstIpArr := strings.Split(dstIP, ",")

	// get ip that in db
	dstCondition := map[string]interface{}{
		common.BKHostInnerIPField: map[string]interface{}{
			common.BKDBIN: dstIpArr,
		},
		common.BKCloudIDField: platID,
	}
	dstHostMap, dstHostIdArr, err := apilogic.GetHostMapByCond(req, dstCondition)
	blog.Debug("dstHostMap:%v", dstHostMap)

	dstConfigCond := map[string]interface{}{
		common.BKAppIDField:  []int{appID},
		common.BKHostIDField: dstHostIdArr,
	}
	dstHostIdArrV, err := logics.GetHostIDByCond(req, host.CC.HostCtrl(), dstConfigCond)
	existIpArr := make([]string, 0)
	for _, id := range dstHostIdArrV {
		if dstHostMapData, ok := dstHostMap[id].(map[string]interface{}); ok {
			existIpArr = append(existIpArr, dstHostMapData[common.BKHostInnerIPField].(string))
		}
	}

	// update host data not eq nil
	updateHostData := make(map[string]interface{})
	for key, val := range hostMapData {
		if nil != val {
			updateHostData[key] = val
		}
	}

	// clone host , the exist host should update and none exist host should add；
	// dstIpArr: all host for clone，existIpArr：the exist host
	blog.Debug("existIpArr:%v", existIpArr)
	for _, dstIpV := range dstIpArr {
		if dstIpV == orgIP {
			blog.Debug("clone host updateHostMain err:%v", err)
			msg := "dstIp and orgIp should not be the same"
			cli.ResponseFailed(common.CC_Err_Comm_Host_Update_FAIL_ERR, fmt.Sprintf("%s%s", common.CC_Err_Comm_Host_Update_FAIL_ERR_STR, msg), resp)
			return
		}
		blog.Debug("hostMapData:%v", hostMapData)

		if apilogic.In_existIpArr(existIpArr, dstIpV) {
			blog.Debug("clone update")
			hostCondition := map[string]interface{}{
				common.BKHostInnerIPField: dstIpV,
			}

			updateHostData[common.BKHostInnerIPField] = dstIpV
			delete(updateHostData, common.BKHostIDField)
			res, err := apilogic.UpdateHostMain(req, hostCondition, updateHostData, appID, host.CC.HostCtrl(), host.CC.ObjCtrl(), host.CC.AuditCtrl(), cli.CC.Error)
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
			cloneHostId, err := apilogic.AddHost(req, addHostMapData, host.CC.ObjCtrl())
			if nil != err {
				blog.Debug("clone host addHost err:%v", err)
				msg := fmt.Sprintf("clone host error:%s", dstIpV)
				cli.ResponseFailed(common.CC_Err_Comm_HOST_CREATE_FAIL, fmt.Sprintf("%s%s", common.CC_Err_Comm_HOST_CREATE_FAIL_STR, msg), resp)
				return
			}

			blog.Debug("cloneHostId:%v configData[0]:%v", cloneHostId, configData[0])

			configDataMap := make(map[string]interface{}, 0)
			configDataMap[common.BKHostIDField] = cloneHostId
			configDataMap[common.BKModuleIDField] = []int{configData[0][common.BKModuleIDField]}
			configDataMap[common.BKAppIDField] = configData[0][common.BKAppIDField]
			configDataMap[common.BKSetIDField] = configData[0][common.BKSetIDField]
			err = apilogic.AddModuleHostConfig(req, configDataMap, host.CC.HostCtrl())
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
