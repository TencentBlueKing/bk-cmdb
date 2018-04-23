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

package process

import (
	"configcenter/src/api_server/ccapi/actions/v2"
	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"
)

var proc *procAction = &procAction{}

type procAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "process/getProcessPortByApplicationID", Params: nil, Handler: proc.GetProcessPortByApplicationID, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "process/getProcessPortByIP", Params: nil, Handler: proc.GetProcessPortByIP, FilterHandler: nil, Version: v2.APIVersion})

	// set cc api interface
	proc.CreateAction()
}

// GetProcessPortByApplicationID  get process port and application id
func (cli *procAction) GetProcessPortByApplicationID(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetProcessPortByApplicationID start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("GetProcessPortByApplicationID error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("GetProcessPortByApplicationID http body data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID"})
	if !res {
		blog.Error("GetProcessPortByApplicationID error: %s", msg)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}

	appId, err := strconv.Atoi(formData.Get("ApplicationID"))
	appIdStr := formData.Get("ApplicationID")
	if nil != err {

		blog.Error("GetProcessPortByApplicationID error: ApplicationID is not number")
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	// 根据APPID获取所有模块
	modules, err := getModulesByAppId(appIdStr, req, resp)
	blog.Debug("modules:%v", modules)
	if nil != err {
		blog.Error("getModulesByAppId error:%v", err)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, err.Error()).Error(), resp)
		return
	}
	modulesMap := modules.([]map[string]interface{})

	paramJson, _ := json.Marshal(modulesMap)
	blog.Debug("paramJson:%v", string(paramJson))
	url := fmt.Sprintf("%s/process/v1/openapi/GetProcessPortByApplicationID/%d", cli.CC.ProcAPI(), appId)
	blog.Debug("url%v", url)
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, paramJson)
	blog.Debug("rspV3:%v", rspV3)
	if err != nil {
		blog.Error("GetProcessPortByApplicationID url:%s, data:%s, error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	rspMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(rspV3), &rspMap)
	if nil != err {
		blog.Error("GetProcessPortByApplicationID error  data:%s", string(rspV3))
		converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Errorf(common.CCErrCommJSONUnmarshalFailed, msg).Error(), resp)
		return
	}
	if !rspMap["result"].(bool) {
		code, _ := util.GetIntByInterface(rspMap[common.HTTPBKAPIErrorCode])
		blog.Error("GetProcessPortByApplicationID error:%s", rspMap[common.HTTPBKAPIErrorMessage])
		converter.RespFailV2(code, rspMap[common.HTTPBKAPIErrorMessage].(string), resp)
		return
	}
	if nil == rspMap["data"] {
		emptyData := make([]interface{}, 0)
		converter.RespSuccessV2(converter.ResV2ToForProcList(emptyData), resp)
		return
	}
	converter.RespSuccessV2(converter.GeneralV2Data(rspMap["data"]), resp)

}

// GetProcessPortByIP get process port by ip
func (cli *procAction) GetProcessPortByIP(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetProcessPortByIP start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("GetProcessPortByIP error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Debug("GetProcessPortByIP http body data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ips"})
	if !res {
		blog.Error("GetProcessPortByIP error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	ips := formData.Get("ips")
	ipArr := strings.Split(ips, ",")
	if len(ipArr) == 0 {
		blog.Error("GetProcessPortByIP error: ips is required")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ips").Error(), resp)
		return
	}
	param := make(common.KvMap)
	param["ipArr"] = ipArr
	paramJson, _ := json.Marshal(param)
	blog.Debug("paramJson:%v", paramJson)
	url := fmt.Sprintf("%s/process/v1/openapi/GetProcessPortByIP", cli.CC.ProcAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, paramJson)
	blog.Debug("rspV3:%v", rspV3)
	if err != nil {
		blog.Error("GetProcessPortByIP url;%s, params:%s, error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	rspMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(rspV3), &rspMap)
	if nil != err {
		blog.Error("GetProcessPortByIP replay:%s error:%s", rspV3, err.Error())
		converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}

	if !rspMap["result"].(bool) {
		blog.Error("GetProcessPortByIP error:%s", rspMap[common.HTTPBKAPIErrorMessage])
		errorIDInt, _ := util.GetIntByInterface(rspMap[common.HTTPBKAPIErrorCode])
		converter.RespFailV2(errorIDInt, rspMap[common.HTTPBKAPIErrorMessage].(string), resp)
		return
	}
	if nil == rspMap["data"] {
		emptyData := make([]interface{}, 0)
		converter.RespSuccessV2(converter.ResV2ToForProcList(emptyData), resp)
		return
	}

	converter.RespSuccessV2(converter.ResV2ToForProcList(rspMap["data"]), resp)
}

func getModulesByAppId(appId string, req *restful.Request, resp *restful.Response) (interface{}, error) {
	blog.Debug("appId:%v", appId)
	//set empty to get all fields
	param := map[string]interface{}{
		"condition": make(map[string]interface{}),
		"fields":    []string{common.BKModuleIDField, common.BKModuleNameField},
		"page": map[string]interface{}{
			"start": 0,
			"limit": 0,
			"sort":  "",
		},
	}
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appId

	searchParams := make(map[string]interface{})
	searchParams["condition"] = condition
	searchParams["fields"] = fmt.Sprintf("%s,%s", common.BKModuleIDField, common.BKModuleNameField)
	searchParams["start"] = 0
	searchParams["limit"] = 0
	searchParams["sort"] = ""
	paramJson, _ := json.Marshal(param)

	url := fmt.Sprintf("%s/topo/v1/openapi/module/searchByApp/%s", proc.CC.TopoAPI(), appId)
	blog.Debug("url:%v", url)
	moduleRes, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	if nil != err {
		blog.Errorf("getModulesByAppId url:%s, params:%s, error:%s", url, string(paramJson), err.Error())
		return nil, err
	}
	blog.Debug("moduleRes:%v", moduleRes)
	resData := make([]map[string]interface{}, 0)
	var rst map[string]interface{}
	if jsErr := json.Unmarshal([]byte(moduleRes), &rst); nil == jsErr {
		if rst["result"].(bool) {
			modules := (rst["data"].(map[string]interface{}))["info"]
			for _, module := range modules.([]interface{}) {
				resData = append(resData, module.(map[string]interface{}))
			}
			return resData, nil
		} else {
			return nil, errors.New(rst[common.HTTPBKAPIErrorMessage].(string))
		}
	} else {

		return nil, jsErr
	}
}
