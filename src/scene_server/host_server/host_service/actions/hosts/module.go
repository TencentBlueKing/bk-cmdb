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

package hosts

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/logics"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"
)

type CloudHostModuleParams struct {
	ApplicationID int          `json:"bk_biz_id"`
	HostInfoArr   []BkHostInfo `json:"host_info"`
	ModuleID      int          `json:"bk_module_id"`
}

type BkHostInfo struct {
	IP      string `json:"bk_host_innerip"`
	CloudID int    `json:"bk_cloud_id"`
}

func init() {

	hostModuleConfig.CreateAction()

	//this api only exsit when host allow in mutile biz
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/modules/biz/mutiple", Params: nil, Handler: hostModuleConfig.AddHostMutiltAppModuleRelation})
}

// HostModuleRelation add host module relation
func (m *hostModuleConfigAction) AddHostMutiltAppModuleRelation(req *restful.Request, resp *restful.Response) {
	defErr := m.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	m.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read input body error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)

		}
		var data CloudHostModuleParams
		//get data
		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		//check if this module exist to do
		module, err := logics.GetModuleByModuleID(req, data.ApplicationID, data.ModuleID, m.CC.ObjCtrl())
		if nil != err {
			blog.Error("get destination module info error, params:%v, error:%v", data.ModuleID, err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoModuleSelectFailed)
		}
		if 0 == len(module) {
			blog.Error("destination module  not found , params:%v, error:%v", data.ModuleID, err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoMulueIDNotfoundFailed)
		}
		var errMsg []string
		var succ []string
		var hostIdArr []int
		for index, hostInfo := range data.HostInfoArr {

			//check if host exist
			hostCond := common.KvMap{
				common.BKHostInnerIPField: hostInfo.IP,
				common.BKCloudIDField:     hostInfo.CloudID,
			}
			hostList, err := logics.GetHostInfoByConds(req, m.CC.HostCtrl(), hostCond)
			if nil != err || 0 == len(hostList) {
				blog.Error("get host info error, params:%v, error:%v", hostCond, err.Error())
				errMsg = append(errMsg, fmt.Sprintf("%s 主机IP在系统中不存在", hostInfo.IP))
				continue

			}

			//check if host in this module
			hostData := hostList[0]
			hostMap, ok := hostData.(map[string]interface{})
			if false == ok {
				blog.Error("host not exsit, params:%v, error:%v", hostCond, err.Error())
				errMsg = append(errMsg, fmt.Sprintf("%s 主机IP在系统中不存在", hostInfo.IP))
				continue
			}
			hostId, err := util.GetIntByInterface(hostMap[common.BKHostIDField])
			if nil != err {
				blog.Error("host not exsit, params:%v, error:%v", hostCond, err.Error())
				errMsg = append(errMsg, fmt.Sprintf("%s 主机IP在系统中不存在", hostInfo.IP))
				continue
			}
			moduleHostCond := common.KvMap{
				common.BKHostIDField:   []int{hostId},
				common.BKModuleIDField: []int{data.ModuleID},
			}
			moduleHostConfig, err := logics.GetConfigByCond(req, m.CC.HostCtrl(), moduleHostCond)
			if nil != err {
				blog.Error("get module host config error, params:%v, error:%v", moduleHostCond, err.Error())
				errMsg = append(errMsg, fmt.Sprintf("%s 获取主机模块关系失败", hostInfo.IP))
				continue
			}
			if 0 != len(moduleHostConfig) {
				blog.Error("host exist in module, params:%v, error:%v", moduleHostCond, err)
				errMsg = append(errMsg, fmt.Sprintf("%s 主机已经存在于当前模块中", hostInfo.IP))
				continue
			}

			//add host to this module
			params := make(map[string]interface{})
			addModulesURL := m.CC.HostCtrl() + "/host/v1/meta/hosts/modules"
			params[common.BKAppIDField] = data.ApplicationID
			params[common.BKModuleIDField] = []int{data.ModuleID}
			params[common.BKHostIDField] = hostId
			isSuccess, errMsgStr, _ := logics.GetHttpResult(req, addModulesURL, common.HTTPCreate, params)
			if !isSuccess {
				blog.Error("add modulehostconfig error, params:%v, error:%s", params, errMsgStr)
				errMsg = append(errMsg, fmt.Sprintf("%s 主机添加到模块失败", hostInfo.IP))

			}
			hostIdArr = append(hostIdArr, hostId)
			if nil != err {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommResourceInitFailed)
			}
			succ = append(succ, fmt.Sprintf("%d", index))
		}

		if 0 != len(errMsg) {
			retData := make(map[string]interface{})
			retData["success"] = succ
			retData["error"] = errMsg
			return http.StatusInternalServerError, retData, defErr.Error(common.CCErrAddHostToModule)
		}

		logClient, err := logics.NewHostModuleConfigLog(req, hostIdArr, m.CC.HostCtrl(), m.CC.ObjCtrl(), m.CC.AuditCtrl())
		user := util.GetActionUser(req)
		logClient.SaveLog(fmt.Sprintf("%d", data.ApplicationID), user)
		return http.StatusOK, nil, nil
	}, resp)

}
