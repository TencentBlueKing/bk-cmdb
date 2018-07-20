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
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"configcenter/src/api_server/ccapi/logics/v2/common/converter"
	"configcenter/src/api_server/ccapi/logics/v2/common/utils"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func (s *Service) getSets(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	user := util.GetUser(pheader)

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("getSets error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("getSets data:%v", formData)

	if len(formData["ApplicationID"]) == 0 {
		blog.Error("getSets error: ApplicationID is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID").Error(), resp)
		return
	}

	condition := mapstr.MapStr{}

	appID := formData.Get("ApplicationID")
	condition[common.BKAppIDField] = appID

	blog.Errorf("len:%v", len(formData["SetServiceStatus"]))
	if "" != formData.Get("SetServiceStatus") {
		formStatus := formData.Get("SetServiceStatus")

		//service status，include 0：close，1：open，default 1
		strStatus := "1"
		if "0" == formStatus {
			strStatus = "0"
		}
		condition[common.BKSetStatusField] = strStatus
	}

	if "" != formData.Get("SetEnviType") {
		formEnv := formData.Get("SetEnviType")

		//env 1 test 2：tiyan 3：formal，default 3
		strEnv := "3"
		switch formEnv {
		case "1":
			strEnv = "1"
		case "2":
			strEnv = "2"
		}
		condition[common.BKSetEnvField] = strEnv
	}

	param := &params.SearchParams{Condition: condition,
		Fields: []string{},
		Page: map[string]interface{}{
			"start": 0,
			"limit": 0,
			"sort":  "",
		}}
	result, err := s.CoreAPI.TopoServer().Instance().SearchSet(context.Background(), user, appID, pheader, param)

	if err != nil {
		blog.Errorf("getSets  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	resDataV2, err := converter.ResToV2ForSetList(result.Result, result.ErrMsg, result.Data)
	if err != nil {
		blog.Errorf("convert set res to v2 error:%v, rspV3:%v", err, result.Data)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) getsetproperty(req *restful.Request, resp *restful.Response) {
	resDataV2 := common.KvMap{
		/*"option" : [
		    {
		        "id" : "2",
		        "name" : "体验",
		        "type" : "text",
		        "is_default" : false
		    },
		    {
		        "id" : "3",
		        "name" : "正式",
		        "type" : "text",
		        "is_default" : true
		    },
		    {
		        "id" : "1",
		        "name" : "测试",
		        "type" : "text",
		        "is_default" : false
		    }
		]*/
		"SetEnviType": []common.KvMap{
			common.KvMap{"Property": "1", "value": "1"},
			common.KvMap{"Property": "2", "value": "2"},
			common.KvMap{"Property": "3", "value": "3"},
		},
		/*"option" : [
		    {
		        "id" : "2",
		        "name" : "关闭",
		        "type" : "text",
		        "is_default" : false
		    },
		    {
		        "id" : "1",
		        "name" : "开放",
		        "type" : "text",
		        "is_default" : true
		    }
		]*/
		"SetServiceStatus": []common.KvMap{
			common.KvMap{"Property": "0", "value": "2"},
			common.KvMap{"Property": "1", "value": "1"},
		},
	}

	converter.RespSuccessV2(resDataV2, resp)
}

func (s *Service) getModulesByProperty(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Error("getModulesByProperty error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("getModulesByProperty data:%v", formData)

	if len(formData["ApplicationID"]) == 0 {
		blog.Error("getModulesByProperty error: ApplicationID is empty!")
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	appID := formData["ApplicationID"][0]

	param := make(map[string]interface{})

	if len(formData["SetID"]) > 0 && formData["SetID"][0] != "" {
		setIDArr, sliceErr := utils.SliceStrToInt(strings.Split(formData["SetID"][0], ","))
		if sliceErr != nil {
			blog.Errorf("SliceStrToInt error:%v", sliceErr)
			converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, sliceErr.Error()).Error(), resp)
			return
		}

		param[common.BKSetIDField] = setIDArr
	}

	if len(formData["SetEnviType"]) > 0 && "" != formData["SetEnviType"][0] && formData["SetID"][0] != "" {
		param[common.BKSetEnvField] = strings.Split(formData["SetEnviType"][0], ",")
	}

	if len(formData["SetServiceStatus"]) > 0 && "" != formData["SetServiceStatus"][0] {
		param[common.BKSetStatusField] = strings.Split(formData["SetServiceStatus"][0], ",")
	}

	result, err := s.CoreAPI.TopoServer().OpenAPI().SearchModuleByProperty(context.Background(), appID, pheader, param)
	if err != nil {
		blog.Errorf("getModulesByProperty  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	resDataV2, err := converter.ResToV2ForModuleMapList(result.Result, result.ErrMsg, result.Data)
	if err != nil {
		blog.Errorf("convert module res to v2 error:%v", err)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}

	converter.RespSuccessV2(resDataV2, resp)

}

func (s *Service) addSet(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	user := util.GetUser(pheader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	defLang := s.Language.CreateDefaultCCLanguageIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("addSet error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("addSet  data: %v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "SetName"})
	if !res {
		blog.Errorf("addSet error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	reqParam := make(map[string]interface{})
	if len(formData["SetName"][0]) > 24 {
		msg = defLang.Language("apiv2_set_name_lt_24")
		blog.Errorf("addSet error:%v", msg)
		converter.RespFailV2(common.CCErrAPIServerV2SetNameLenErr, defErr.Errorf(common.CCErrAPIServerV2SetNameLenErr, msg).Error(), resp)
		return
	}

	reqParam[common.BKSetNameField] = formData["SetName"][0]

	appID := formData["ApplicationID"][0]
	reqParam[common.BKInstParentStr], err = strconv.Atoi(appID)
	if nil != err {
		blog.Errorf("convert appid to int error:%v", err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID").Error(), resp)
		return
	}

	if len(formData["properties"]) > 0 {
		propertiesJson := formData["properties"][0]
		propertiesMap := make(map[string]interface{})
		err := json.Unmarshal([]byte(propertiesJson), &propertiesMap)
		if nil != err {
			blog.Errorf("addSet error:%v", err)
			converter.RespFailV2(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
			return
		}
		for k, v := range propertiesMap {
			if k == "EnviType" {
				vInt, err := util.GetIntByInterface(v)
				if nil != err {
					blog.Errorf("addSet error:%v", err)
					converter.RespFailV2(common.CCErrCommParamsIsInvalid, defErr.Errorf(common.CCErrCommParamsIsInvalid, "EnviType").Error(), resp)
					return
				}
				vStr := strconv.Itoa(vInt)

				env := getAddSetEnv(vStr)
				reqParam[common.BKSetEnvField] = env
				continue
			}
			if k == "ServiceStatus" {
				vInt, err := util.GetIntByInterface(v)
				if nil != err {
					blog.Errorf("addSet error:%v", err)
					converter.RespFailV2(common.CCErrCommParamsIsInvalid, defErr.Errorf(common.CCErrCommParamsIsInvalid, "ServiceStatus").Error(), resp)
					return
				}
				vStr := strconv.Itoa(vInt)
				status := getAddSetService(vStr)
				reqParam[common.BKSetStatusField] = status
				continue
			}
			if k == "Description" {
				reqParam[common.BKSetDescField] = v
				continue
			}
			if k == "Capacity" {
				reqParam[common.BKSetCapacityField], err = util.GetIntByInterface(v)
				if nil != err {
					blog.Errorf("addSet GetIntByInterface error:%v", err)
					converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsIsInvalid, "Capacity").Error(), resp)
					return
				}
				continue
			}
			reqParam[k] = v
		}
	}
	topoLevel, err := s.Logics.CheckAppTopoIsThreeLevel(user, pheader)
	if err != nil {
		blog.Error("AddSet CheckAppTopoIsThreeLevel error:%v", err)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, err.Error()).Error(), resp)
		return
	}
	if false == topoLevel {
		blog.Error("AddSet CheckAppTopoIsThreeLevel  mainline topology level not three")
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, "business topology level is more than three, please use the v3 api instead").Error(), resp)
		return
	}
	delete(reqParam, "ChnName")

	reqParam, err = s.Logics.AutoInputV3Field(reqParam, common.BKInnerObjIDSet, user, pheader)
	blog.Infof("addSet reqParam:%v", reqParam)

	result, err := s.CoreAPI.TopoServer().Instance().CreateSet(context.Background(), appID, pheader, reqParam)
	if nil != err {
		blog.Errorf("addSet error:%v", err)
		converter.RespFailV2(common.CCErrTopoSetCreateFailed, defErr.Error(common.CCErrTopoSetCreateFailed).Error(), resp)
		return
	}
	if !result.Result {
		msg = fmt.Sprintf("%s", result.ErrMsg)
		blog.Errorf("addSet error:%s", msg)
		converter.RespFailV2(common.CCErrTopoSetCreateFailed, defErr.Error(common.CCErrTopoSetCreateFailed).Error(), resp)
		return
	}
	rspDataV3Map := result.Data
	blog.Infof("rsp_v3:%v", result.Data)
	converter.RespSuccessV2(rspDataV3Map, resp)
}

func getSetService(formData url.Values) string {
	formStatus := ""
	if len(formData["ServiceStatus"]) > 0 {
		formStatus = formData["ServiceStatus"][0]
	}

	strStatus := "1"
	if "0" == formStatus {
		strStatus = "2"
	}
	return strStatus
}

func getAddSetService(v string) string {
	formStatus := v

	strStatus := "1"
	if "0" == formStatus {
		strStatus = "2"
	}
	return strStatus
}

func getAddSetEnv(v string) string {

	formEnv := v

	//env 1：测试 2：体验 3：正式，默认为3
	strEnv := "3"
	switch formEnv {
	case "1":
		strEnv = "1"
	case "2":
		strEnv = "2"
	}
	return strEnv
}

func getSetServiceAndEnv(formData url.Values) (string, string) {

	formEnv := ""

	if len(formData["SetEnviType"]) > 0 {
		formEnv = formData["SetEnviType"][0]
	}
	//env 1：测试 2：体验 3：正式，默认为3
	strEnv := "3"
	switch formEnv {
	case "1":
		strEnv = "1"
	case "2":
		strEnv = "1"
	}
	strStatus := getSetService(formData)
	return strStatus, strEnv
}

func (s *Service) updateSet(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("updateSet error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("updateSet error data: %v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "SetID"})
	if !res {
		blog.Errorf("updateSet error error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	reqParam := make(map[string]interface{})
	reqData := make(map[string]interface{})

	setIDStrArr := strings.Split(formData["SetID"][0], ",")
	setIDArr, err := utils.SliceStrToInt(setIDStrArr)
	if nil != err {
		blog.Errorf("updateSet error:%v", err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "SetID").Error(), resp)
		return
	}

	appID := formData["ApplicationID"][0]

	reqParam[common.BKSetIDField] = setIDArr
	setName := formData["SetName"][0]

	if len(setIDArr) == 1 && len(formData["SetName"]) > 0 {

		reqData[common.BKSetNameField] = setName
	} else {
		msg := "once can only update one set"
		blog.Errorf("update set error:%v", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}
	if len(formData["SetName"][0]) > 24 {
		msg = "set name over 24"
		blog.Errorf("add set error:%v", msg)
		converter.RespFailV2(common.CCErrAPIServerV2SetNameLenErr, defErr.Error(common.CCErrAPIServerV2SetNameLenErr).Error(), resp)
		return
	}
	status, env := getSetServiceAndEnv(formData)
	if len(formData["SetEnviType"]) > 0 {
		reqData[common.BKSetEnvField] = env
	}
	if len(formData["Des"]) > 0 {
		reqData[common.BKSetDescField] = formData["Des"][0]
	}
	if len(formData["Capacity"]) > 0 {
		reqData[common.BKSetCapacityField], err = util.GetIntByInterface(formData["Capacity"][0])
		if nil != err {
			blog.Error("updateSet error:%v", err)
			converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "Capacity").Error(), resp)
			return
		}
	}

	if len(formData["ServiceStatus"]) > 0 {
		reqData[common.BKSetStatusField] = status
	}

	delete(reqData, "ChnName")

	reqParam["data"] = reqData

	result, err := s.CoreAPI.TopoServer().OpenAPI().UpdateMultiSet(context.Background(), appID, pheader, reqParam)

	if nil != err {
		blog.Errorf("updateSet ,error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)
}

func (s *Service) updateSetServiceStatus(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("updateSetServiceStatus error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("updateSetServiceStatus  data: %v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "SetID", "Status"})
	if !res {
		blog.Errorf("updateSetServiceStatus error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	reqParam := make(map[string]interface{})
	reqData := make(map[string]interface{})

	setIDStrArr := strings.Split(formData["SetID"][0], ",")
	setIDArr, err := utils.SliceStrToInt(setIDStrArr)
	if nil != err {
		blog.Errorf("updateSetServiceStatus error:%v", err)
		converter.RespFailV2(common.CCErrCommParamsNeedInt, defErr.Errorf(common.CCErrCommParamsNeedInt, "SetID").Error(), resp)
		return
	}

	appID := formData["ApplicationID"][0]

	reqData[common.BKSetStatusField] = formData["Status"][0]

	// service status  combin 0：关闭，1：开启，默认为1
	switch reqData[common.BKSetStatusField] {
	case "0":
		reqData[common.BKSetStatusField] = "2"
	case "1":
		reqData[common.BKSetStatusField] = "1"

	}

	reqParam[common.BKSetIDField] = setIDArr
	reqParam["data"] = reqData

	result, err := s.CoreAPI.TopoServer().OpenAPI().UpdateMultiSet(context.Background(), appID, pheader, reqParam)
	if nil != err {
		blog.Errorf("updateSetServiceStatus  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)
}

func (s *Service) delSet(req *restful.Request, resp *restful.Response) {

	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("delSet error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.Infof("delSet data: %v", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "SetID"})
	if !res {
		blog.Errorf("delSet error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	appID := formData["ApplicationID"][0]

	reqParam := make(map[string]interface{})
	reqParam[common.BKSetIDField] = formData["SetID"][0]

	result, err := s.CoreAPI.TopoServer().OpenAPI().DeleteMultiSet(context.Background(), appID, pheader, reqParam)
	if nil != err {
		blog.Errorf("delSet  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)
}

func (s *Service) delSetHost(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("delSetHost error:%v", err)
		converter.RespFailV2(common.CCErrCommPostInputParseError, defErr.Error(common.CCErrCommPostInputParseError).Error(), resp)
		return
	}

	formData := req.Request.Form

	blog.V(3).Infof("delSetHost data: %s", formData)

	res, msg := utils.ValidateFormData(formData, []string{"ApplicationID", "SetID"})
	if !res {
		blog.Errorf("delSetHost error: %s", msg)
		converter.RespFailV2(common.CCErrAPIServerV2DirectErr, defErr.Errorf(common.CCErrAPIServerV2DirectErr, msg).Error(), resp)
		return
	}

	appID := formData["ApplicationID"][0]
	setIDStrArr := strings.Split(formData["SetID"][0], ",")
	setIDArr, err := utils.SliceStrToInt(setIDStrArr)
	reqParam := make(map[string]interface{})
	reqParam[common.BKSetIDField] = setIDArr

	result, err := s.CoreAPI.TopoServer().OpenAPI().DeleteSetHost(context.Background(), appID, pheader, reqParam)
	if err != nil {
		blog.Errorf("delSetHost  error:%v", err)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}

	converter.RespCommonResV2(result.Result, result.Code, result.ErrMsg, resp)

}
