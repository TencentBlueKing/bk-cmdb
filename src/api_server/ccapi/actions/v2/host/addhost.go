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

package host

import (
	"configcenter/src/api_server/ccapi/actions/v2"
	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"strings"
)

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/addhost", Params: nil, Handler: host.AddHost, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/enterIp", Params: nil, Handler: host.EnterIP, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/enterip", Params: nil, Handler: host.EnterIP, FilterHandler: nil, Version: v2.APIVersion})

}

func (cli *hostAction) AddHost(req *restful.Request, resp *restful.Response) {
	blog.Debug("EnterIP start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("EnterIP error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form
	ips := formData.Get("ip")
	moduleName := formData.Get("moduleName")
	appName := formData.Get("appName")
	setName := formData.Get("setName")
	platID := formData.Get("platId")

	if "" == ips {
		blog.Errorf("AddHost error ip empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ip").Error(), resp)
		return
	}
	ipArr := strings.Split(ips, ",")

	intPlatID, err := util.GetIntByInterface(platID)
	if nil != err {
		blog.Errorf("AddHost error ip empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "platId").Error(), resp)
		return
	}

	input := make(common.KvMap)
	input["ips"] = ipArr
	input[common.BKModuleNameField] = moduleName
	input[common.BKSetNameField] = setName
	input[common.BKAppNameField] = appName
	input[common.BKCloudIDField] = intPlatID
	input[common.BKOwnerIDField] = common.BKDefaultOwnerID
	input["is_increment"] = true

	paramJson, _ := json.Marshal(input)
	url := fmt.Sprintf("%s/host/v1/host/add/module", cli.CC.HostAPI())
	blog.Infof("http request for add module url:%s, params:%s", url, string(paramJson))
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	blog.Infof("http request for add module url:%s, reply:%s", url, rspV3)
	if err != nil {
		blog.Error("addhost url:%s, params:%s, error:%s", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	err = converter.ResToV2ForEnterIP(rspV3)
	if err != nil {
		blog.Error("convert addhost result to v2 error:%s, reply:%s", err.Error(), rspV3)
		converter.RespFailV2(common.CCErrAddHostToModuleFailStr, err.Error(), resp)
		return
	}
	converter.RespSuccessV2("", resp)
}

func (cli *hostAction) EnterIP(req *restful.Request, resp *restful.Response) {
	blog.Debug("EnterIP start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("EnterIP error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form
	ips := formData.Get("ip")
	hostName := formData.Get("hostname")
	moduleName := formData.Get("moduleName")
	appName := formData.Get("appName")
	setName := formData.Get("setName")
	osType := formData.Get("osType")

	if "" == ips {
		blog.Errorf("EnterIP error ips empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ip").Error(), resp)
		return
	}
	ipArr := strings.Split(ips, ",")
	var hostNameArr []string
	if "" != hostName {
		hostNameArr = strings.Split(hostName, ",")
	}
	if osType == "window" {
		osType = "windows"
	}
	if "" != osType && osType != "windows" && osType != "linux" {
		blog.Errorf("osType mast be windows or linux; not %s", osType)
		converter.RespFailV2(common.CCErrAPIServerV2OSTypeErr, defErr.Error(common.CCErrAPIServerV2OSTypeErr).Error(), resp)
		return
	}
	osType = "1"
	if "windows" == osType {
		osType = "2"
	}
	input := make(common.KvMap)
	input["ips"] = ipArr
	input[common.BKHostNameField] = hostNameArr
	input[common.BKModuleNameField] = moduleName
	input[common.BKSetNameField] = setName
	input[common.BKAppNameField] = appName
	input[common.BKOSTypeField] = osType
	input[common.BKOwnerIDField] = common.BKDefaultOwnerID

	paramJson, _ := json.Marshal(input)
	url := fmt.Sprintf("%s/host/v1/host/add/module", cli.CC.HostAPI())
	blog.Infof("http request for add module url:%s, params:%s", url, string(paramJson))
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	blog.Infof("http request for add module url:%s, reply:%s", url, rspV3)
	if err != nil {
		blog.Error("EnterIP url:%s, params:%s, error:%s", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	err = converter.ResToV2ForEnterIP(rspV3)
	if err != nil {
		blog.Error("convert EnterIP result to v2 error:%s, reply:%s", err.Error(), rspV3)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}
	converter.RespSuccessV2("", resp)
}
