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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"

	v2 "configcenter/src/api_server/ccapi/actions/v2/phpapi"
	logics "configcenter/src/api_server/ccapi/logics/v2"
	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

var host *hostAction = &hostAction{}

type hostAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Host/gethostlistbyip", Params: nil, Handler: host.GetHostListByIP, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/gethostlistbyip", Params: nil, Handler: host.GetHostListByIP, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Host/getsethostlist", Params: nil, Handler: host.GetSetHostList, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/getmodulehostlist", Params: nil, Handler: host.GetModuleHostList, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Host/getmodulehostlist", Params: nil, Handler: host.GetModuleHostList, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/getapphostlist", Params: nil, Handler: host.GetAppHostList, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Host/getapphostlist", Params: nil, Handler: host.GetAppHostList, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Set/gethostsbyproperty", Params: nil, Handler: host.GetHostsByProperty, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "set/gethostsbyproperty", Params: nil, Handler: host.GetHostsByProperty, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "Host/updateHostStatus", Params: nil, Handler: host.UpdateHostStatus, FilterHandler: nil, Version: v2.APIVersion})

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "Host/updateHostByAppId", Params: nil, Handler: host.UpdateHostByAppID, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Host/getCompanyIdByIps", Params: nil, Handler: host.GetCompanyIDByIps, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/getCompanyIdByIps", Params: nil, Handler: host.GetCompanyIDByIps, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Host/getHostListByAppidAndField", Params: nil, Handler: host.GetHostListByAppIDAndField, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/getHostListByAppidAndField", Params: nil, Handler: host.GetHostListByAppIDAndField, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Host/getIPAndProxyByCompany", Params: nil, Handler: host.GetIPAndProxyByCompany, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "Host/updatehostmodule", Params: nil, Handler: host.UpdateHostModule, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "host/updatehostmodule", Params: nil, Handler: host.UpdateHostModule, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "host/updateCustomProperty", Params: nil, Handler: host.UpdateCustomProperty, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "host/cloneHostProperty", Params: nil, Handler: host.CloneHostProperty, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "host/delHostInApp", Params: nil, Handler: host.DelHostInApp, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "host/getgitServerIp", Params: nil, Handler: host.GetGitServerIp, FilterHandler: nil, Version: v2.APIVersion})

	// set cc api interface
	host.CreateAction()
}

// updateHostStatus update host proxy status , set GseProxy field equal 1
func (cli *hostAction) UpdateHostStatus(req *restful.Request, resp *restful.Response) {
	blog.Debug("updateHostStatus start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("updateHostStatus error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("updateHostStatus http body data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"appId", "platId", "ip"})
	if !res {
		blog.Error("updateHostStatus error: %s", msg)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	appID := formData["appId"][0]

	ip := formData["ip"][0]
	platID, _ := strconv.Atoi(formData["platId"][0])

	param := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKHostInnerIPField: ip,
			common.BKSubAreaField:     platID,
		},
		"data": map[string]interface{}{
			common.BKGseProxyField: "1",
			common.BKSubAreaField:  platID,
		},
	}
	paramJson, err := json.Marshal(param)
	if err != nil {
		blog.Error("updateHostStatus error:%v", err)
		converter.RespFailV2(common.CCErrCommJSONMarshalFailed, defErr.Error(common.CCErrCommJSONMarshalFailed).Error(), resp)
		return
	}

	url := fmt.Sprintf("%s/host/v1/openapi/host/%s", cli.CC.HostAPI(), appID)
	blog.Debug("url:%v", url)

	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, []byte(paramJson))
	blog.Debug("rspV3:%v", rspV3)
	if err != nil {
		blog.Error("updateHostStatus url:%s,  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2([]byte(rspV3), resp)
}

// updateHostStatus update proxy satus by ip
func (cli *hostAction) UpdateHostByAppID(req *restful.Request, resp *restful.Response) {
	blog.Debug("updateHostByAppID start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("updateHostByAppID error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("updateHostByAppID http body data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"appId", "platId", "proxyList"})
	if !res {
		blog.Error("updateHostByAppID error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	appID := formData["appId"][0]

	proxyList := formData["proxyList"][0]

	platID := formData.Get("platId")
	platIdInt, _ := strconv.Atoi(platID)
	proxyListArr := make([]map[string]interface{}, 0)
	json.Unmarshal([]byte(proxyList), &proxyListArr)
	proxyListArrV3 := make([]common.KvMap, 0)
	for _, proxy := range proxyListArr {
		proxyNew := make(map[string]interface{})
		proxyNew[common.BKCloudIDField] = platIdInt
		proxyNew[common.BKHostInnerIPField] = proxy["InnerIP"]
		proxyNew[common.BKHostOuterIPField] = proxy["OuterIP"]
		proxyNew, inputErr := logics.AutoInputV3Field(proxyNew, common.BKInnerObjIDHost, host.CC.TopoAPI(), req.Request.Header)

		if inputErr != nil {
			blog.Error("AutoInputV3Field error:%v", inputErr)
			converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, inputErr.Error()).Error(), resp)
			return
		}
		proxyListArrV3 = append(proxyListArrV3, proxyNew)
	}

	blog.Debug("proxyListArrV3:%v", proxyListArrV3)
	param := map[string]interface{}{

		common.BKCloudIDField:   platID,
		common.BKProxyListField: proxyListArrV3,
	}
	paramJson, _ := json.Marshal(param)

	url := fmt.Sprintf("%s/host/v1/host/updateHostByAppID/%s", cli.CC.HostAPI(), appID)
	blog.Debug(url)
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, []byte(paramJson))
	blog.Debug("rspV3:%v", rspV3)
	if err != nil {
		blog.Error("updateHostByAppID  url:%s, params:%s, error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2([]byte(rspV3), resp)
}

//updateHostModule update host module relation
func (cli *hostAction) UpdateHostModule(req *restful.Request, resp *restful.Response) {
	blog.Debug("UpdateHostModule start")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("UpdateHostModule Error %v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form

	blog.Infof("UpdateHostModule http body data:%v", formData)

	appID := formData.Get("ApplicationID")
	platID := formData.Get("platId")
	moduleID := formData.Get("dstModuleID")
	ips := formData.Get("ip")

	if "" == appID {
		blog.Errorf("UpdateHostModule error ApplicationID empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}
	if "" == platID {
		blog.Errorf("UpdateHostModule error platID empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "platID").Error(), resp)
		return
	}
	if "" == moduleID {
		blog.Errorf("UpdateHostModule error moduleID empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "dstModuleID").Error(), resp)
		return
	}
	if "" == ips {
		blog.Errorf("UpdateHostModule error ips empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ips").Error(), resp)
		return
	}
	ipArr := strings.Split(ips, ",")
	moduleIDArr, err := utils.SliceStrToInt(strings.Split(moduleID, ","))
	if nil != err {
		blog.Error("getHostListByAppidAndField error: %v", err)
		converter.RespFailV2(common.CCErrAPIServerV2MultiModuleIDErr, defErr.Error(common.CCErrAPIServerV2MultiModuleIDErr).Error(), resp)
		return
	}
	appIDArr := make([]int, 0)
	appIDInt, _ := strconv.Atoi(appID)
	appIDArr = append(appIDArr, appIDInt)
	platIDInt, _ := strconv.Atoi(platID)
	param := map[string]interface{}{
		common.BKIPListField:  ipArr,
		common.BKAppIDField:   appIDArr,
		common.BKSubAreaField: platIDInt,
	}
	paramJson, _ := json.Marshal(param)

	url := fmt.Sprintf("%s/host/v1/gethostlistbyip", cli.CC.HostAPI())
	hosts, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	if nil != err {
		blog.Error("UpdateHostModule url:%s, params:%s, error:%s ", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	js, err := simplejson.NewJson([]byte(hosts))
	if nil != err {
		converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		blog.Errorf("UpdateHostModule error simplejson.NewJson error, data:%s", string(hosts))
		return
	}
	hostsMap, err := js.Map()
	if nil != err {
		blog.Errorf("UpdateHostModule error js.Map error, data:%s", string(hosts))
		converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}
	HostIDArr := make([]int, 0)
	for _, item := range hostsMap["data"].([]interface{}) {
		blog.Debug("item:%v", item)
		hostID := item.(map[string]interface{})[common.BKHostIDField]
		blog.Debug("hostIDInt:%d", hostID)
		hostIDInt, _ := util.GetIntByInterface(hostID)
		HostIDArr = append(HostIDArr, int(hostIDInt))

	}
	blog.Debug("HostIDArr:%v", HostIDArr)
	urlInput := ""
	// host translate module
	moduleMap, err := GetModulesByApp(req, map[string]interface{}{
		common.BKAppIDField: appID}, cli.CC.TopoAPI())
	blog.Debug("moduleMap:%v", moduleMap)
	if nil != err {
		blog.Error("GetModulesByApp error: %v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Errorf(common.CCErrCommReplyDataFormatError, err.Error()).Error(), resp)
		return
	}

	input := make(common.KvMap)
	input[common.BKAppIDField] = appIDInt
	input[common.BKHostIDField] = HostIDArr
	if len(moduleIDArr) > 1 {
		for moduleId := range moduleIDArr {
			moduleName := moduleMap[moduleId].(map[string]interface{})[common.BKModuleNameField].(string)
			if moduleName == common.DefaultFaultModuleName || moduleName == common.DefaultResModuleName {
				msg := defErr.Error(common.CCErrAPIServerV2HostModuleContainDefaultModuleErr).Error()
				blog.Error("update host module error: %v", msg)
				converter.RespFailV2(common.CCErrAPIServerV2HostModuleContainDefaultModuleErr, msg, resp)
				return
			}
		}
		urlInput = fmt.Sprintf("%s/host/v1/hosts/modules", cli.CC.HostAPI())
		input[common.BKModuleIDField] = moduleIDArr
		input[common.BKIsIncrementField] = false
	} else {
		if moduleMap[moduleIDArr[0]].(map[string]interface{})[common.BKModuleNameField].(string) == common.DefaultFaultModuleName {
			urlInput = fmt.Sprintf("%s/host/v1/hosts/faultmodule", cli.CC.HostAPI())
		}
		if moduleMap[moduleIDArr[0]].(map[string]interface{})[common.BKModuleNameField].(string) == common.DefaultResModuleName {
			urlInput = fmt.Sprintf("%s/host/v1/hosts/emptymodule", cli.CC.HostAPI())
		} else {
			urlInput = fmt.Sprintf("%s/host/v1/hosts/modules", cli.CC.HostAPI())
			input[common.BKModuleIDField] = moduleIDArr
			input[common.BKIsIncrementField] = false
		}
	}

	inputJson, _ := json.Marshal(input)
	blog.Infof("update host module url:%v", url)
	blog.Infof("http request for hosts search url:%s, params:%s", urlInput, string(inputJson))
	rspV3, err := httpcli.ReqHttp(req, urlInput, common.HTTPCreate, []byte(inputJson))
	blog.Debug("http request for  hosts search url:%s, reply:%s", urlInput, rspV3)
	if err != nil {
		blog.Error("UpdateHostModule url:%s, params:%s, error:%s , reply:%s", url, string(inputJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	converter.RespCommonResV2([]byte(rspV3), resp)
}

// UpdateCustomProperty update host property
func (cli *hostAction) UpdateCustomProperty(req *restful.Request, resp *restful.Response) {
	blog.Debug("UpdateCustomProperty start")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if nil != err {
		blog.Error("UpdateHostModule Error %v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Error("formData:%v", formData)
	appId := formData.Get("ApplicationID")
	hostId := formData.Get("HostID")
	property := formData.Get("Property")
	if "" == appId {
		blog.Errorf("UpdateCustomProperty error platId empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}
	if "" == hostId {
		blog.Errorf("UpdateCustomProperty error host empty")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "HostID").Error(), resp)
		return
	}
	param := make(common.KvMap)
	param[common.BKAppIDField] = appId
	param[common.BKHostIDField] = hostId
	param["property"] = property
	paramJson, _ := json.Marshal(param)
	url := fmt.Sprintf("%s/host/v1/openapi/updatecustomproperty", cli.CC.HostAPI())
	blog.Infof("http request for UpdateCustomProperty url:%s, params:%s", url, string(paramJson))
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, []byte(paramJson))
	blog.Infof("http request for UpdateCustomProperty url:%s, reply:%s", url, rspV3)
	if err != nil {
		blog.Error("UpdateCustomProperty url:%s, params:%s, error:%s , reply:%s", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	rspV3Map := make(map[string]interface{})
	err = json.Unmarshal([]byte(rspV3), &rspV3Map)
	if nil != err {
		converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}
	converter.RespCommonResV2([]byte(rspV3), resp)
}

// CloneHostProperty clone host
func (cli *hostAction) CloneHostProperty(req *restful.Request, resp *restful.Response) {
	blog.Debug("CloneHostProperty start")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if nil != err {
		blog.Error("CloneHostProperty Error %v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Error("formData:%v", formData)
	res, msg := utils.ValidateFormData(formData, []string{
		"ApplicationID",
		"orgIp",
		"dstIp",
	})
	if !res {
		blog.Error("CloneHostProperty error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	appId := formData.Get("ApplicationID")
	orgIp := formData.Get("orgIp")
	dstIp := formData.Get("dstIp")
	platId := formData.Get("platId")

	appIdI, err := strconv.ParseInt(appId, 10, 64)
	if nil != err {
		blog.Error("CloneHostProperty error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	platIdI, err := strconv.ParseInt(platId, 10, 64)

	if nil != err {
		blog.Error("CloneHostProperty error: %s", msg)
		platIdI = 0
	}

	var param metadata.HostCloneInputParams
	param.AppID = appIdI
	param.DstIP = dstIp
	param.OrgIP = orgIp
	param.PlatID = platIdI
	paramJson, _ := json.Marshal(param)

	url := fmt.Sprintf("%s/host/v1/propery/clone", cli.CC.HostAPI())
	blog.Infof("http request for CloneHostProperty url:%s, params:%s", url, string(paramJson))
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, []byte(paramJson))
	blog.Infof("http request for CloneHostProperty url:%s, reply:%s", url, rspV3)
	if err != nil {
		blog.Error("CloneHostProperty url:%s, params:%s, error:%s", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	rspV3Map := make(map[string]interface{})
	err = json.Unmarshal([]byte(rspV3), &rspV3Map)
	if nil != err {
		blog.Error("CloneHostProperty Unmarshal json error:%v, rspV3:%s", err, rspV3)
		return
	}
	converter.RespCommonResV2([]byte(rspV3), resp)
}

// DelHostInApp del host from application idle set
func (cli *hostAction) DelHostInApp(req *restful.Request, resp *restful.Response) {
	blog.Debug("DelHostInApp start")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if nil != err {
		blog.Error("DelHostInApp Error %v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}
	formData := req.Request.Form
	blog.Debug("formData:%v", formData)
	res, msg := utils.ValidateFormData(formData, []string{
		"ApplicationID",
		"HostID",
	})
	if !res {
		blog.Error("DelHostInApp error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	appId := formData.Get("ApplicationID")
	hostId, _ := util.GetInt64ByInterface(formData.Get("HostID"))
	param := make(common.KvMap)
	param[common.BKAppIDField], _ = util.GetInt64ByInterface(appId)
	param[common.BKHostIDField] = []int64{hostId}
	paramJson, _ := json.Marshal(param)
	//url := fmt.Sprintf("%s/host/v1/openapi/host/delhostinapp", cli.CC.HostAPI())
	url := fmt.Sprintf("%s/host/v1/hosts/resource", cli.CC.HostAPI())
	blog.Infof("http request for DelHostInApp url:%s, params:%s", url, string(paramJson))
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPCreate, []byte(paramJson))
	blog.Infof("http request for DelHostInApp url:%s, reply:%s", url, rspV3)
	if err != nil {
		blog.Error("DelHostInApp url:%s, params:%s, error:%s", url, string(paramJson), err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	rspV3Map := make(map[string]interface{})
	err = json.Unmarshal([]byte(rspV3), &rspV3Map)
	if nil != err {
		blog.Error("DelHostInApp Unmarshal json error:%v, rspV3:%s", err, rspV3)
		return
	}
	converter.RespCommonResV2([]byte(rspV3), resp)
}
