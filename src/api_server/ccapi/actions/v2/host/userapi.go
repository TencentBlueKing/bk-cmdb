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
	logics "configcenter/src/api_server/ccapi/logics/v2"
	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/common/commondata"
	"encoding/json"
	"fmt"
	"strings"

	"configcenter/src/api_server/ccapi/actions/v2"

	restful "github.com/emicklei/go-restful"
)

var userAPI *userAPIAction = &userAPIAction{}

type userAPIAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/CustomerGroup/getContentByCustomerGroupID", Params: nil, Handler: userAPI.getContentByCustomerGroupID, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/CustomerGroup/getContentByCustomerGroupId", Params: nil, Handler: userAPI.getContentByCustomerGroupID, FilterHandler: nil, Version: v2.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/CustomerGroup/getCustomerGroupList", Params: nil, Handler: userAPI.getCustomerGroupList, FilterHandler: nil, Version: v2.APIVersion})

	// set cc api interface
	userAPI.CreateAction()
}

// getCustomerGroupList get user api list
func (cli *userAPIAction) getCustomerGroupList(req *restful.Request, resp *restful.Response) {
	blog.Debug("getCustomerGroupList start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("getCustomerGroupList error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form
	strAppIDs := formData.Get("ApplicationIDs")
	//userAPIType := formData.Get("Type")

	if "" == strAppIDs {
		blog.Error("getCustomerGroupList error: param ApplicationIDs is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationIDs").Error(), resp)
		return
	}

	appIDs := strings.Split(strAppIDs, ",")

	urlPrefix := fmt.Sprintf("%s/host/v1/userapi/search/", cli.CC.HostAPI())

	var postInput commondata.ObjQueryInput
	postInput.Start = 0
	postInput.Limit = common.BKNoLimit
	paramJson, _ := json.Marshal(postInput)

	var resDataV2 []common.KvMap
	//traverse all application ids
	for _, appID := range appIDs {
		url := fmt.Sprintf("%s%s", urlPrefix, appID)
		rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
		//http request error
		if err != nil {
			blog.Error("getCustomerGroupList url:%s, params:%s error:%v", url, string(paramJson), err)
			converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
			return
		}

		//translate cmdb v3 to v2 api result
		retItem, err := converter.ResToV2ForCustomerGroup(rspV3, appID)

		//translate cmdb v3 to v2 api result error,
		if err != nil {
			blog.Error("GetHostListByIP error:%s, reply:%s", err.Error(), rspV3)
			converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
			return
		}
		if 0 == len(retItem) {
			continue
		}
		resDataV2 = append(resDataV2, common.KvMap{"ApplicationID": appID, "CustomerGroup": retItem})

	}

	converter.RespSuccessV2(resDataV2, resp)
}

// getContentByCustomerGroupID  get user api result by user api id
func (cli *userAPIAction) getContentByCustomerGroupID(req *restful.Request, resp *restful.Response) {
	blog.Debug("getContentByCustomerGroupID start!")
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("getContentByCustomerGroupID error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	appID := formData.Get("ApplicationID")
	id := formData.Get("CustomerGroupID")

	version := formData.Get("version")
	page := formData.Get("page")
	pageSize := formData.Get("pageSize")

	if "" == appID {
		blog.Error("getContentByCustomerGroupID error: param ApplicationID is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}

	if "" == id {
		blog.Error("getContentByCustomerGroupID error: param CustomerGroupID is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "CustomerGroupID").Error(), resp)
		return
	}

	detailURL := fmt.Sprintf("%s/host/v1/userapi/detail/%s/%s", cli.CC.HostAPI(), appID, id)
	userAPI := logics.NewUserAPI()
	name, err, errCode := userAPI.GetNameByID(req, detailURL)
	if nil != err {
		blog.Errorf("getContentByCustomerGroupID error: get CustomerGroup name is error! %s", err.Error())
		converter.RespFailV2(errCode, err.Error(), resp)
		return
	}

	skip := "0"
	if "1" == version {
		intPage, _ := util.GetIntByInterface(page)
		intPageSize, _ := util.GetIntByInterface(pageSize)
		if intPage > 0 {
			intPage -= 1
		}
		skip = fmt.Sprintf("%d", intPage*intPageSize)

	} else {
		pageSize = fmt.Sprintf("%d", common.BKNoLimit)
	}
	dataURL := fmt.Sprintf("%s/host/v1/userapi/data/%s/%s/%s/%s", cli.CC.HostAPI(), appID, id, skip, pageSize)
	rspDetailV3, err := httpcli.ReqHttp(req, dataURL, common.HTTPSelectGet, nil)
	if nil != err {
		blog.Errorf("http request url:%s, error:%s", dataURL, err.Error())
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	//translate cmdb v3 to v2 api result
	list, total, err := converter.ResToV2ForCustomerGroupResult(rspDetailV3)

	//translate cmdb v3 to v2 api result error,
	if err != nil {
		blog.Error("GetHostListByIP error:%s, reply:%s", err.Error(), rspDetailV3)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	if "1" == version {
		ret := make(common.KvMap)
		ret["list"] = list
		ret["total"] = total
		ret["page"] = page
		ret["pageSize"] = pageSize
		ret["GroupName"] = name

		converter.RespSuccessV2(ret, resp)
	} else {
		converter.RespSuccessV2(list, resp)

	}

}
