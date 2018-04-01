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
	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"encoding/json"
	"fmt"
	"strconv"

	"configcenter/src/api_server/ccapi/actions/v2"

	"configcenter/src/common/util"
	"github.com/emicklei/go-restful"
)

var gse *gseAction = &gseAction{}

type gseAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/getAgentStatus", Params: nil, Handler: gse.GetAgentStatus, FilterHandler: nil, Version: v2.APIVersion})

	// set cc api interface
	gse.CreateAction()
}

//  GetAgentStatus get host agent host
func (cli *gseAction) GetAgentStatus(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetAgentStatus start!")

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("GetAgentStatus error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form
	blog.Debug("GetAgentStatus http body data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"appId"})
	if !res {
		blog.Error("GetAgentStatus error: %s", msg)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "appId").Error(), resp)
		return
	}

	appID, err := strconv.Atoi(formData["appId"][0])
	if nil != err {
		blog.Error("GetAgentStatus error: ApplicationID is not number")
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, "appId").Error(), resp)
		return
	}

	url := fmt.Sprintf("%s/host/v1/getAgentStatus/%d", cli.CC.HostAPI(), appID)
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectGet, nil)

	if err != nil {
		blog.Error("GetAgentStatus error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	rspMap := make(map[string]interface{})
	json.Unmarshal([]byte(rspV3), &rspMap)
	if nil != rspMap["result"] {
		if !rspMap["result"].(bool) {
			blog.Error("GetAgentStatus error:%s", rspMap[common.HTTPBKAPIErrorMessage])
			converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, rspMap[common.HTTPBKAPIErrorMessage]).Error(), resp)
			return
		}
	} else {
		msg = "获取AgentStatus失败"
		blog.Error("GetAgentStatus error rspMap:%v", rspMap)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	converter.RespSuccessV2(rspMap["data"], resp)

}
