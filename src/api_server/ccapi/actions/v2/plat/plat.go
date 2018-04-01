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
 
package plat

import (
	logics "configcenter/src/api_server/ccapi/logics/v2"
	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"encoding/json"
	"fmt"
	"strings"

	"configcenter/src/api_server/ccapi/actions/v2"

	"github.com/emicklei/go-restful"
)

var plat *platAction = &platAction{}

type platAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "Plat/updateHost", Params: nil, Handler: plat.updateHost, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Plat/get", Params: nil, Handler: plat.GetPlats, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "Plat/get", Params: nil, Handler: plat.GetPlats, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Plat/delete", Params: nil, Handler: plat.DeletePlats, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "Plat/add", Params: nil, Handler: plat.CreatePlats, FilterHandler: nil, Version: v2.APIVersion})

	plat.CreateAction()
}

// updateHost update host cloud id by ApplicationID, ip,orgPlatID
func (cli *platAction) updateHost(req *restful.Request, resp *restful.Response) {
	blog.Debug("updateHost start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("updateHost error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("updateHost http body data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"appId", "orgPlatId", "ip", "dstPlatId"})
	if !res {
		blog.Error("ValidateFormData error:%s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	appID := formData["appId"][0]
	orgPlatID, _ := util.GetIntByInterface(formData["orgPlatId"][0])
	ip := strings.Trim(formData["ip"][0], " ")
	dstPlatID, _ := util.GetIntByInterface(formData["dstPlatId"][0])

	param := map[string]interface{}{
		"condition": map[string]interface{}{
			common.BKHostInnerIPField: ip,
			common.BKSubAreaField:     orgPlatID,
		},
		"data": map[string]interface{}{
			common.BKSubAreaField: dstPlatID,
		},
	}

	paramJson, _ := json.Marshal(param)

	url := fmt.Sprintf("%s/host/v1/openapi/host/%s", cli.CC.HostAPI(), appID)
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, []byte(paramJson))
	if err != nil {
		blog.Error("updateHost url:%s, params:%s, error:%v", url, paramJson, err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2([]byte(rspV3), resp)

}

// GetPlats: 获取所有云平台
func (cli *platAction) GetPlats(req *restful.Request, resp *restful.Response) {
	blog.Debug("GetPlats start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	url := fmt.Sprintf("%s/host/v1/plat", cli.CC.HostAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectGet, nil)
	if err != nil {
		blog.Error("GetPlats url:%s, error:%v", url, err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForPlatList(rspV3)
	if err != nil {
		blog.Error("convert plat res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

// todo: 确认删除逻辑
// DeletePlats: 根据PlatID删除云平台
func (cli *platAction) DeletePlats(req *restful.Request, resp *restful.Response) {
	blog.Debug("DeletePlats start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("DeletePlats error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("DeletePlats http body data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"platId"})
	if !res {
		blog.Error("ValidateFormData error:%s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	platID := formData["platId"][0]

	url := fmt.Sprintf("%s/host/v1/plat/%s", cli.CC.HostAPI(), platID)
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPDelete, nil)
	if err != nil {
		blog.Error("DeletePlats url:%s, error:%v", url, err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2([]byte(rspV3), resp)
}

// CreatePlats: 添加云平台
func (cli *platAction) CreatePlats(req *restful.Request, resp *restful.Response) {
	blog.Debug("CreatePlats start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("CreatePlats error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Debug("CreatePlats http body data:%v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"platName"})
	if !res {
		blog.Error("ValidateFormData error:%s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	platName := formData["platName"][0]
	param := map[string]interface{}{
		common.BKCloudNameField: platName,
	}

	param, _ = logics.AutoInputV3Filed(param, common.BKInnerObjIDPlat, cli.CC.TopoAPI(), req.Request.Header)
	paramJson, _ := json.Marshal(param)

	url := fmt.Sprintf("%s/host/v1/plat", cli.CC.HostAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPCreate, []byte(paramJson))
	if err != nil {
		blog.Error("CreatePlats url:%s, params:%s, error:%v", url, string(paramJson), err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	rspV3Map := make(map[string]interface{})
	err = json.Unmarshal([]byte(rspV3), &rspV3Map)
	if err != nil {
		blog.Error("Unmarshal v3 result   to map error:, data:%s, error:%v", rspV3, err)
		converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}

	if !rspV3Map["result"].(bool) {
		blog.Error("CreatePlats error:%s", rspV3Map[common.HTTPBKAPIErrorMessage])
		converter.RespFailV2(rspV3Map[common.HTTPBKAPIErrorCode].(int), rspV3Map[common.HTTPBKAPIErrorMessage].(string), resp)
		return
	}

	rspDataV3Map, ok := rspV3Map["data"].(map[string]interface{})
	if false == ok {
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(rspDataV3Map[common.BKCloudIDField], resp)
}
