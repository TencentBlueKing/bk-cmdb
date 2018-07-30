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

package service

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func (s *Service) getCustomerGroupList(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("getCustomerGroupList error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form
	strAppIDs := formData.Get("ApplicationIDs")

	if "" == strAppIDs {
		blog.Error("getCustomerGroupList error: param ApplicationIDs is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationIDs").Error(), resp)
		return
	}

	appIDs := strings.Split(strAppIDs, ",")

	var postInput metadata.QueryInput
	postInput.Start = 0
	postInput.Limit = common.BKNoLimit

	var resDataV2 []mapstr.MapStr

	// all application ids
	for _, appID := range appIDs {

		result, err := s.CoreAPI.HostServer().GetUserCustomQuery(context.Background(), appID, pheader, &postInput)
		if err != nil {
			blog.Error("getCustomerGroupList error:%v", err)
			converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
			return
		}

		//translate cmdb v3 to v2 api result
		retItem, err := converter.ResToV2ForCustomerGroup(result.Result, result.ErrMsg, result.Data, appID)

		//translate cmdb v3 to v2 api result error,
		if err != nil {
			blog.Error("getCustomerGroupList error:%s, reply:%v", err.Error(), result.Data)
			converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
			return
		}
		if 0 == len(retItem) {
			continue
		}
		resDataV2 = append(resDataV2, mapstr.MapStr{"ApplicationID": appID, "CustomerGroup": retItem})

	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) getContentByCustomerGroupID(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

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

	name, err, errCode := s.GetNameByID(appID, id, pheader)
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
		skip = strconv.Itoa(intPage * intPageSize)

	} else {
		pageSize = strconv.Itoa(common.BKNoLimit)
	}

	result, err := s.CoreAPI.HostServer().GetUserCustomQueryResult(context.Background(), appID, id, skip, pageSize, pheader)
	if nil != err {
		blog.Errorf("http request  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	//translate cmdb v3 to v2 api result
	list, total, err := converter.ResToV2ForCustomerGroupResult(result.Result, result.ErrMsg, result.Data)

	//translate cmdb v3 to v2 api result error,
	if err != nil {
		blog.Error("getContentByCustomerGroupID  v%", result)
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

func (s *Service) GetNameByID(appID, id string, pheader http.Header) (string, error, int) {

	result, err := s.CoreAPI.HostServer().GetUserCustomQueryDetail(context.Background(), appID, id, pheader)
	//http request error
	if err != nil {
		blog.Errorf("GetNameByID error:%v", err)
		return "", nil, common.CCErrCommHTTPDoRequestFailed
	}

	if result.Result {
		name, _ := result.Data["name"].(string)
		return name, nil, 0
	} else {
		return "", errors.New(result.ErrMsg), common.CCErrCommHTTPDoRequestFailed
	}

}
