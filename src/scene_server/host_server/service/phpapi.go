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
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/validator"

	"github.com/emicklei/go-restful"
)

// updateHostPlat 根据条件更新主机信息
func (s *Service) UpdateHost(req *restful.Request, resp *restful.Response) {
	blog.V(3).Infof("updateHost start!")
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))

	appID, err := util.GetInt64ByInterface(req.PathParameter(common.BKAppIDField))
	if nil != err {
		blog.Errorf("convert appid %s to int error:%v", req.PathParameter(common.BKAppIDField), err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("updateHost , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	data, httpCode, errMsg := s.Logics.UpdateHost(input, appID, req.Request.Header)

	if nil != errMsg {
		blog.Errorf("UpdateHost update host, appID:%d, input:%v, error:%s", appID, input, err)
		resp.WriteError(httpCode, &meta.RespError{Msg: errMsg})
		return
	}
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     data,
	})

}

func (s *Service) UpdateHostByAppID(req *restful.Request, resp *restful.Response) {
	blog.V(3).Infof("updateHostByAppID start!")
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))
	appID, err := util.GetInt64ByInterface(req.PathParameter("appid"))
	if nil != err {
		blog.Error("convert appid to int error:%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsLostField, "ApplicationID")})
		return
	}

	input := new(meta.UpdateHostParams)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("updateHostByAppID , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.V(3).Infof("updateHostByAppID http body data: %v", input)
	result, httpCode, errMsg := s.Logics.UpdateHostByAppID(input, appID, req.Request.Header)
	if nil != errMsg {
		blog.Errorf("updateHostByAppID update host, appID:%d, input:%v, error:%s", appID, input, err)
		resp.WriteError(httpCode, &meta.RespError{Msg: errMsg})
		return
	}
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	})

}

func (s *Service) HostSearchByIP(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))

	input := new(meta.HostSearchByIPParams)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("HostSearchByIP , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if 0 == len(input.IpList) {
		blog.Error("input does not contains key IP")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsLostField, common.BKIPListField)})
	}

	orCondition := []map[string]interface{}{
		map[string]interface{}{common.BKHostInnerIPField: map[string]interface{}{common.BKDBIN: input.IpList}},
		map[string]interface{}{common.BKHostOuterIPField: map[string]interface{}{common.BKDBIN: input.IpList}},
	}
	hostMapCondition := map[string]interface{}{common.BKDBOR: orCondition}

	if nil != input.CloudID {
		hostMapCondition[common.BKCloudIDField] = input.CloudID
	}

	phpapi := s.Logics.NewPHPAPI(req.Request.Header)
	hostMap, hostIDArr, err := phpapi.GetHostMapByCond(hostMapCondition)
	if err != nil {
		blog.Errorf("HostSearchByIP error : %s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}

	configCond := map[string][]int64{
		common.BKHostIDField: hostIDArr,
	}
	if 0 < len(input.AppID) {
		configCond[common.BKAppIDField] = input.AppID
	}

	configData, err := s.Logics.GetConfigByCond(req.Request.Header, configCond)
	if nil != err {
		blog.Errorf("HostSearchByIP error : %s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}
	hostData, err := phpapi.SetHostData(configData, hostMap)
	if nil != err {
		blog.Error("HostSearchByIP error : %v", err)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostData,
	})
}

func (s *Service) HostSearchByConds(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("HostSearchByConds , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	phpapi := s.Logics.NewPHPAPI(req.Request.Header)
	hostMap, hostIDArr, err := phpapi.GetHostMapByCond(input)
	if err != nil {
		blog.Error("HostSearchByConds error : %v, input:%s", err, input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}

	configCond := map[string][]int64{
		common.BKHostIDField: hostIDArr,
	}
	configData, err := s.Logics.GetConfigByCond(req.Request.Header, configCond)
	if nil != err {
		blog.Error("HostSearchByConds error : %v, input:%v", err, input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}

	hostData, err := phpapi.SetHostData(configData, hostMap)
	if nil != err {
		blog.Error("HostSearchByConds error : %v, input:%v", err, input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostData,
	})
}

func (s *Service) HostSearchByModuleID(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))

	input := new(meta.HostSearchByModuleIDParams)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("HostSearchByModuleID , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if nil == input.ApplicationID {
		blog.Error("HostSearchByModuleID input does not contains key ApplicationID")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsLostField, "ApplicationID")})
		return
	}

	if nil == input.ModuleID {
		blog.Error("HostSearchByModuleID input does not contains key ModuleID")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsLostField, "ModuleID")})
		return
	}

	configData, err := s.Logics.GetConfigByCond(req.Request.Header, map[string][]int64{
		common.BKModuleIDField: input.ModuleID,
		common.BKAppIDField:    []int64{*input.ApplicationID},
	})
	if nil != err {
		blog.Errorf("HostSearchByModuleID get host module config error:%s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}

	phpapi := s.Logics.NewPHPAPI(req.Request.Header)
	hostData, err := phpapi.GetHostDataByConfig(configData)
	if nil != err {
		blog.Errorf("HostSearchByModuleID get host module config error:%s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostData,
	})
}

func (s *Service) HostSearchBySetID(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))

	input := new(meta.HostSearchBySetIDParams)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("HostSearchBySetID , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	if nil == input.ApplicationID {
		blog.Error("HostSearchBySetID input does not contains key ApplicationID")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsLostField, "ApplicationID")})
		return
	}

	conds := make(map[string][]int64)
	conds[common.BKAppIDField] = []int64{*input.ApplicationID}

	if len(input.SetID) > 0 {
		conds[common.BKSetIDField] = input.SetID
	}

	configData, err := s.Logics.GetConfigByCond(req.Request.Header, conds)
	if nil != err {
		blog.Errorf("HostSearchBySetID get host module config error:%s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}

	phpapi := s.Logics.NewPHPAPI(req.Request.Header)
	hostData, err := phpapi.GetHostDataByConfig(configData)
	if nil != err {
		blog.Errorf("HostSearchByModuleID get host module config error:%s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostData,
	})
}

func (s *Service) HostSearchByAppID(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))

	input := new(meta.HostSearchByAppIDParams)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("HostSearchByAppID , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if nil == input.ApplicationID {
		blog.Error("HostSearchByAppID input does not contains key ApplicationID")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsLostField, "ApplicationID")})
		return
	}

	configData, err := s.Logics.GetConfigByCond(req.Request.Header, map[string][]int64{
		common.BKAppIDField: []int64{*input.ApplicationID},
	})

	if nil != err {
		blog.Errorf("HostSearchByModuleID get host module config error:%s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}

	phpapi := s.Logics.NewPHPAPI(req.Request.Header)
	hostData, err := phpapi.GetHostDataByConfig(configData)
	if nil != err {
		blog.Errorf("HostSearchByModuleID get host module config error:%s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostData,
	})

}

func (s *Service) HostSearchByProperty(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("HostSearchByProperty , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	appID, err := util.GetInt64ByInterface(input[common.BKAppIDField])
	if nil != err {
		blog.Error("HostSearchByProperty input does not contains key ApplicationID")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsLostField, "ApplicationID")})
		return
	}

	setCondition := make([]meta.ConditionItem, 0)

	setIDArrI, hasSetID := input[common.BKSetIDField]
	if hasSetID {
		cond := meta.ConditionItem{}
		cond.Field = common.BKSetIDField
		cond.Operator = common.BKDBIN
		cond.Value = setIDArrI
		setCondition = append(setCondition, cond)
	}

	setEnvTypeArr, hasSetEnvType := input[common.BKSetEnvField]
	if hasSetEnvType {
		cond := meta.ConditionItem{}
		cond.Field = common.BKSetEnvField
		cond.Operator = common.BKDBIN
		cond.Value = setEnvTypeArr
		setCondition = append(setCondition, cond)
	}

	setSrvStatusArr, hasSetSrvStatus := input[common.BKSetStatusField]
	if hasSetSrvStatus {
		cond := meta.ConditionItem{}
		cond.Field = common.BKSetStatusField
		cond.Operator = common.BKDBIN
		cond.Value = setSrvStatusArr
		setCondition = append(setCondition, cond)
	}

	blog.V(3).Infof("HostSearchByProperty setCondition: %v\n", setCondition)
	setIDArr, err := s.Logics.GetSetIDByCond(req.Request.Header, setCondition)
	if nil != err {
		blog.Errorf("HostSearchByProperty get host module config error:%s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostGetSetFaild, err.Error())})
		return
	}
	blog.V(3).Infof("HostSearchByProperty ApplicationID: %s, SetID: %v\n", appID, setIDArr)

	condition := map[string][]int64{
		common.BKAppIDField: []int64{appID},
	}

	condition[common.BKSetIDField] = setIDArr
	configData, err := s.Logics.GetConfigByCond(req.Request.Header, condition)
	if nil != err {
		blog.Errorf("HostSearchByProperty get host module config error:%s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}

	phpapi := s.Logics.NewPHPAPI(req.Request.Header)
	hostData, err := phpapi.GetHostDataByConfig(configData)
	if nil != err {
		blog.Errorf("HostSearchByProperty get host module config error:%s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostData,
	})
}

func (s *Service) GetIPAndProxyByCompany(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))

	input := new(meta.GetIPAndProxyByCompanyParams)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("GetIPAndProxyByCompany , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if 0 == len(input.Ips) {
		blog.Error("GetIPAndProxyByCompany input does not contains key IP")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsLostField, "IP")})
		return
	}

	appIDInt, err := util.GetInt64ByInterface(*input.AppIDStr)
	if nil != err {
		blog.Error("GetIPAndProxyByCompany input application id not integer, input:%v", input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "ApplicationID")})
		return
	}

	platIDInt, err := util.GetInt64ByInterface(*input.CloudIDStr)
	if nil != err {
		blog.Error("GetIPAndProxyByCompany cloud id not integer, input:%v", input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "PlatID")})
		return
	}
	resData, err := s.Logics.GetIPAndProxyByCompany(input.Ips, platIDInt, appIDInt, req.Request.Header)
	if nil != err {
		blog.Errorf("GetIPAndProxyByCompany error:%s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: err})
		return
	}
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     resData,
	})

}

func (s *Service) UpdateCustomProperty(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("UpdateCustomProperty , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	blog.Error("UpdateCustomProperty  input:%v", input)
	appID, err := util.GetInt64ByInterface(input[common.BKAppIDField])
	if nil != err {
		blog.Error("UpdateCustomProperty input not found appID, input:%v", input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID")})
		return
	}
	hostID, err := util.GetInt64ByInterface(input[common.BKHostIDField])
	if nil != err {
		blog.Error("UpdateCustomProperty input not found hostID, input:%v", input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, "HostID")})
		return
	}
	propertyJson, ok := input["property"].(string)
	if false == ok && "" == propertyJson {
		blog.Error("UpdateCustomPropertyinput not found property, input:%v", input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "property")})
		return
	}

	propertyMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(propertyJson), &propertyMap)
	if nil != err {
		blog.Error("UpdateCustomPropertyinput not found property, input:%v", input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	res, err := s.Logics.UpdateCustomProperty(hostID, appID, propertyMap, req.Request.Header)
	if nil != err {
		blog.Error("UpdateCustomPropertyinput not found property, input:%v", input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     res,
	})

}

func (s *Service) CloneHostProperty(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	input := &meta.CloneHostPropertyParams{}
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("CloneHostProperty , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.V(0).Infof("CloneHostProperty input:%v", input)
	appID, err := util.GetInt64ByInterface(input.AppIDStr) // util.GetInt64ByInterface(input.[common.BKAppIDField])
	if nil != err {
		blog.Errorf("CloneHostProperty ,appliation not int , err: %v, input:%v", err, input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID")})
		return
	}
	if "" == input.CloudIDStr {
		blog.Errorf("CloneHostProperty ,set not found , err: %v, input:%v", err, input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "ApplicationID")})
		return
	}

	CloudID, err := util.GetInt64ByInterface(input.CloudIDStr)
	if nil != err {
		blog.Errorf("CloneHostProperty ,appliation not int , err: %v, input:%v", err, input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID")})
		return
	}
	res, err := s.Logics.CloneHostProperty(input, appID, CloudID, req.Request.Header)
	if nil != err {
		blog.Errorf("CloneHostProperty ,appliation not int , err: %v, input:%v", err, input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     res,
	})
}

func (s *Service) GetHostAppByCompanyId(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	input := &meta.GetHostAppByCompanyIDParams{}
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("GetHostAppByCompanyId , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.V(3).Infof("GetHostAppByCompanyId input:%v", input)
	platId, err := util.GetInt64ByInterface(input.CloudIDStr)
	if nil != err {
		blog.Errorf("GetHostAppByCompanyId cloud id not integer, input:%v", input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "PlatID")})
		return
	}
	ipArr := strings.Split(input.IPs, ",")
	hostCon := map[string]interface{}{
		common.BKHostInnerIPField: map[string]interface{}{
			common.BKDBIN: ipArr,
		},
		common.BKCloudIDField: platId,
	}

	phpapi := s.Logics.NewPHPAPI(req.Request.Header)
	//根据i,platId获取主机
	hostArr, hostIdArr, err := phpapi.GetHostMapByCond(hostCon) // phpapilogic.GetHostMapByCond(req, hostCon)
	if nil != err {
		blog.Errorf("GetHostAppByCompanyId getHostMapByCond:%s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostGetFail)})
		return
	}
	blog.V(3).Infof("GetHostAppByCompanyId hostArr:%v, input:%v", hostArr, input)
	if len(hostIdArr) == 0 {
		resp.WriteEntity(meta.Response{
			BaseResp: meta.SuccessBaseResp,
			Data:     make([]interface{}, 0),
		})
		return
	}
	// 根据主机hostId获取app_id,module_id,set_id
	configCon := map[string][]int64{
		common.BKHostIDField: hostIdArr,
	}
	configArr, err := s.Logics.GetConfigByCond(req.Request.Header, configCon)
	if nil != err {
		blog.Errorf("GetHostAppByCompanyId getConfigByCond err:%s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostGetFail)})
		return
	}
	blog.V(3).Infof("GetHostAppByCompanyId configArr:%v, input:%v", configArr, input)
	if len(configArr) == 0 {
		resp.WriteEntity(meta.Response{
			BaseResp: meta.SuccessBaseResp,
			Data:     make([]interface{}, 0),
		})
		return
	}
	appIdArr := make([]int64, 0)
	setIdArr := make([]int64, 0)
	moduleIdArr := make([]int64, 0)
	for _, item := range configArr {
		appIdArr = append(appIdArr, item[common.BKAppIDField])
		setIdArr = append(setIdArr, item[common.BKSetIDField])
		moduleIdArr = append(moduleIdArr, item[common.BKModuleIDField])
	}
	hostMapArr, err := phpapi.SetHostData(configArr, hostArr) //phpapilogic.SetHostData(req, configArr, hostArr)
	if nil != err {
		blog.Errorf("GetHostAppByCompanyId setHostData err:%s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostGetFail)})
		return
	}
	blog.V(3).Infof("GetHostAppByCompanyId hostMap:%v, input:%v", hostMapArr, input)
	hostDataArr := make([]interface{}, 0)
	for _, h := range hostMapArr {
		hostMap := h.(map[string]interface{})
		hostDataArr = append(hostDataArr, hostMap)
	}
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostDataArr,
	})
}

func (s *Service) DelHostInApp(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	input := &meta.DelHostInAppParams
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("DelHostInApp , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	blog.V(3).Infof("DelHostInApp input:%v", input)
	appID, err := util.GetInt64ByInterface(input.AppID)
	if nil != err {
		blog.Errorf("GetHostAppByCompanyId cloud id not integer, input:%v", input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "ApplicationID")})
		return
	}
	hostID, err := util.GetInt64ByInterface(input.HostID)
	if nil != err {
		blog.Errorf("GetHostAppByCompanyId host id not integer, input:%v", input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, "HostID")})
		return
	}
	configCon := map[string][]int64{
		common.BKAppIDField:  []int64{appID},
		common.BKHostIDField: []int64{hostID},
	}

	configArr, err := s.Logics.GetConfigByCond(req.Request.Header, configCon)
	if err != nil {
		blog.Errorf("DelHostInApp GetConfigByCond err msg:%v, input:%v", err, input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}
	if len(configArr) == 0 {
		blog.Errorf("DelHostInApp not fint hostId:%v in appId:%v, input:%v", hostID, appID, input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrHostNotFound)})
		return
	}
	moduleIdArr := make([]int64, 0)
	for _, item := range configArr {
		moduleIdArr = append(moduleIdArr, item[common.BKModuleIDField])
	}
	moduleCon := map[string]interface{}{
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: moduleIdArr,
		},
		common.BKDefaultField: common.DefaultResModuleFlag,
	}
	moduleArr, err := s.Logics.GetModuleMapByCond(req.Request.Header, common.BKModuleIDField, moduleCon)
	if err != nil {
		blog.Errorf("DelHostInApp GetConfigByCond err msg, error:%s, input:%s", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostGetModuleFail, err.Error())})
		return
	}
	blog.V(3).Infof("DelHostInApp moduleArr:%v, input:%v", moduleArr, input)
	if len(moduleArr) == 0 {
		blog.Errorf("DelHostInApp GetModuleMapByCond   not find host in idle module input: %v", input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostNotFound)})
		return
	}
	param := make(common.KvMap)
	param[common.BKAppIDField] = appID
	param[common.BKHostIDField] = hostID
	res, err := s.CoreAPI.ObjectController().OpenAPI().DeleteSetHost(context.Background(), req.Request.Header, param)
	if nil != err {
		blog.Errorf("DelHostInApp DeleteSetHost   error:%s,  input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostEditRelationPoolFail)})
		return
	}
	if false == res.Result {
		blog.Errorf("DelHostInApp DeleteSetHost   error:%s,  input:%v", res.ErrMsg, input)
		resp.WriteHeaderAndJson(http.StatusBadGateway, res, common.BKHTTPMIMEJSON)
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     res.Data,
	})

}

func (s *Service) GetGitServerIp(req *restful.Request, resp *restful.Response) {

	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	input := new(meta.GitServerIpParams)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("GetGitServerIp , but decode body failed, err: %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	var appID, setID, moduleID int64

	// 根据appName获取app
	appCondition := map[string]interface{}{
		common.BKAppNameField: input.AppName,
	}
	appMap, err := s.Logics.GetAppMapByCond(req.Request.Header, "", appCondition) //  logics.GetAppMapByCond(req, "", cli.CC.ObjCtrl(), appCondition)
	if nil != err {
		blog.Errorf("GetGitServerIp GetAppMapByCond error:%s, input:%s", err.Error(), input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostGetAPPFail, err.Error())})
		return
	}
	if 0 == len(appMap) {
		resp.WriteEntity(meta.Response{
			BaseResp: meta.SuccessBaseResp,
			Data:     make([]interface{}, 0),
		})
		return
	}

	for key, _ := range appMap {
		appID = key
	}

	// 根据setName获取set信息
	setCondition := map[string]interface{}{
		common.BKSetNameField: input.AppName,
		common.BKAppIDField:   appID,
	}
	setMap, err := s.Logics.GetSetMapByCond(req.Request.Header, "", setCondition)
	if nil != err {
		blog.Errorf("GetGitServerIp GetSetMapByCond error:%s, input:%s", err.Error(), input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostGetSetFaild, err.Error())})
		return
	}
	if 0 == len(setMap) {
		resp.WriteEntity(meta.Response{
			BaseResp: meta.SuccessBaseResp,
			Data:     make([]interface{}, 0),
		})
		return
	}
	for key, _ := range setMap {
		setID = key
	}

	// 根据moduleName获取module信息
	moduleCondition := map[string]interface{}{
		common.BKModuleNameField: input.ModuleName,
		common.BKAppIDField:      appID,
	}
	moduleMap, err := s.Logics.GetModuleMapByCond(req.Request.Header, "", moduleCondition)
	if nil != err {
		blog.Errorf("GetGitServerIp GetModuleMapByCond error:%s, input:%s", err.Error(), input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostGetSetFaild, err.Error())})
		return
	}
	for key, _ := range moduleMap {
		moduleID = key
	}

	if len(moduleMap) == 0 {
		resp.WriteEntity(meta.Response{
			BaseResp: meta.SuccessBaseResp,
			Data:     make([]interface{}, 0),
		})
		return
	}
	// 根据 appId,setId,moduleId 获取主机信息
	//configData := make([]map[string]int,0)
	confMap := map[string][]int64{
		common.BKAppIDField:    []int64{appID},
		common.BKSetIDField:    []int64{setID},
		common.BKModuleIDField: []int64{moduleID},
	}
	configData, err := s.Logics.GetConfigByCond(req.Request.Header, confMap)
	if nil != err {
		blog.Errorf("GetGitServerIp GetModuleMapByCond error:%s, input:%s", err.Error(), input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}
	blog.V(3).Infof("GetGitServerIp configData:%v", configData)
	phpapi := s.Logics.NewPHPAPI(req.Request.Header)
	hostArr, err := phpapi.GetHostDataByConfig(configData)
	if nil != err {
		blog.Error("GetGitServerIp getHostDataByConfig error:%s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}
	blog.V(3).Infof("GetGitServerIp hostArr:%v, input:%v", hostArr, input)

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostArr,
	})
}

func (s *Service) GetPlat(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	params := new(meta.QueryInput)
	params.Limit = 0
	res, err := s.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDPlat, req.Request.Header, params)
	if nil != err {
		blog.Error("GetPlat error: %v", err)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrTopoGetCloudErrStrFaild, err.Error())})
		return
	}
	if false == res.Result {
		blog.Error("GetPlat error: %s", res.ErrMsg)
		resp.WriteHeaderAndJson(http.StatusBadGateway, res, common.BKHTTPMIMEJSON)

	} else {
		resp.WriteEntity(meta.Response{
			BaseResp: meta.SuccessBaseResp,
			Data:     res.Data,
		})
	}

}

func (s *Service) CreatePlat(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); nil != err {
		blog.Errorf("CreatePlat , but decode body failed, err: %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	ownerId := util.GetOwnerID(req.Request.Header)
	input[common.BKOwnerIDField] = ownerId

	valid := validator.NewValidMap(util.GetOwnerID(req.Request.Header), common.BKInnerObjIDPlat, req.Request.Header, s.Engine)
	validErr := valid.ValidMap(input, common.ValidCreate, 0)
	if nil != validErr {
		blog.Errorf("CreatePlat error: %v, input:%v", validErr, input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrTopoInstCreateFailed)})
		return
	}

	res, err := s.CoreAPI.ObjectController().Instance().CreateObject(context.Background(), common.BKInnerObjIDPlat, req.Request.Header, input)
	if nil != err {
		blog.Errorf("CreatePlat error: %s, input:%v", err.Error(), input)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrTopoInstCreateFailed)})
		return
	}

	if false == res.Result {
		blog.Errorf("GetPlat error: %s", res.ErrMsg)
		resp.WriteHeaderAndJson(http.StatusBadGateway, res, common.BKHTTPMIMEJSON)

	}
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     res.Data,
	})

}

func (s *Service) DelPlat(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	platID, convErr := util.GetInt64ByInterface(req.PathParameter(common.BKCloudIDField))
	if nil != convErr || 0 == platID {
		blog.Error("the platID is invalid, error info is %s, input:%s", convErr.Error(), platID)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, convErr.Error())})
		return
	}

	params := new(meta.QueryInput)
	params.Fields = common.BKHostIDField
	params.Condition = map[string]interface{}{
		common.BKCloudIDField: platID,
	}

	hostRes, err := s.CoreAPI.HostController().Host().GetHosts(context.Background(), req.Request.Header, params)
	if nil != err {
		blog.Error("DelPlat search host error: %s, input:%v", err.Error(), platID)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostGetFail)})
		return
	}

	if 0 < hostRes.Data.Count {
		blog.Error("DelPlat plat [%d] has host data, can not delete", platID)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrTopoHasHostCheckFailed)})
		return
	}

	param := make(map[string]interface{})
	param[common.BKCloudIDField] = platID
	res, err := s.CoreAPI.ObjectController().Instance().DelObject(context.Background(), common.BKInnerObjIDPlat, req.Request.Header, param)
	if nil != err {
		blog.Error("DelPlat error: %v, input:%d", err, platID)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrTopoInstDeleteFailed)})
		return
	}

	if false == res.Result {
		blog.Errorf("GetPlat error: %s", res.ErrMsg)
		resp.WriteHeaderAndJson(http.StatusBadGateway, res, common.BKHTTPMIMEJSON)

	} else {
		resp.WriteEntity(meta.Response{
			BaseResp: meta.SuccessBaseResp,
			Data:     "",
		})
	}
}

func (s *Service) GetAgentStatus(req *restful.Request, resp *restful.Response) {
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	// 获取AppID
	pathParams := req.PathParameters()
	appID, err := util.GetInt64ByInterface(pathParams["appid"])
	if nil != err {
		blog.Errorf("GetAgentStatus error :%s", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	res, err := s.Logics.GetAgentStatus(appID, &s.Config.Gse, req.Request.Header)
	if nil != err {
		blog.Error("GetAgentStatus error :%v", err)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, err.Error())})
		return
	}
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     res,
	})

}

func (s *Service) getHostListByAppidAndField(req *restful.Request, resp *restful.Response) {

	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	// 获取AppID
	pathParams := req.PathParameters()
	appID, err := util.GetInt64ByInterface(pathParams[common.BKAppIDField])
	if nil != err {
		blog.Errorf("getHostListByAppidAndField error :%s", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	field := req.PathParameter("field")

	configData, err := s.Logics.GetConfigByCond(req.Request.Header, map[string][]int64{
		common.BKAppIDField: []int64{appID},
	})

	if nil != err {
		blog.Errorf("getHostListByAppidAndField error : %s, input:%v", err.Error(), common.KvMap{"appid": appID, "field": field})
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}

	hostIDArr := make([]int64, 0)
	for _, config := range configData {
		hostIDArr = append(hostIDArr, config[common.BKHostIDField])
	}

	query := new(meta.QueryInput)
	query.Fields = fmt.Sprintf("%s,%s,%s,%s,", common.BKHostInnerIPField, common.BKCloudIDField, common.BKAppIDField, common.BKHostIDField) + field
	query.Condition = map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDArr,
		},
	}
	ret, err := s.Logics.CoreAPI.HostController().Host().GetHosts(context.Background(), req.Request.Header, query)
	if nil != err {
		blog.Error("getHostListByAppidAndField search host error: %s, input:%v", err.Error(), common.KvMap{"appid": appID, "field": field})
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostGetFail)})
		return
	}
	if !ret.Result {
		blog.Error("getHostListByAppidAndField search host error: %s, input:%v", ret.ErrMsg, common.KvMap{"appid": appID, "field": field})
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrHostGetFail)})
		return
	}
	retData := make(map[string][]interface{})
	for _, itemMap := range ret.Data.Info {
		fieldValue, ok := itemMap[field]
		if !ok {
			continue
		}

		fieldValueStr := fmt.Sprintf("%v", fieldValue)
		groupData, ok := retData[fieldValueStr]
		if ok {
			retData[fieldValueStr] = append(groupData, itemMap)
		} else {
			retData[fieldValueStr] = []interface{}{
				itemMap,
			}
		}
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     retData,
	})
}
