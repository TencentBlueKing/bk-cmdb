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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/validator"

	"github.com/emicklei/go-restful"
)

// updateHostPlat 根据条件更新主机信息
func (s *Service) UpdateHost(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	blog.V(5).Infof("updateHost start!,rid:%s", srvData.rid)

	appID, err := util.GetInt64ByInterface(req.PathParameter(common.BKAppIDField))
	if nil != err {
		blog.Errorf("convert appid %s to int error:%v,rid:%s", req.PathParameter(common.BKAppIDField), err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("updateHost , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	data, httpCode, errMsg := srvData.lgc.UpdateHost(srvData.ctx, input, appID)

	if nil != errMsg {
		blog.Errorf("UpdateHost update host, appID:%d, input:%+v, error:%s,rid:%s", appID, input, err, srvData.rid)
		resp.WriteError(httpCode, &meta.RespError{Msg: errMsg})
		return
	}
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     data,
	})

}

func (s *Service) UpdateHostByAppID(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	blog.V(5).Infof("updateHostByAppID start!,rid:%s", srvData.rid)
	appID, err := util.GetInt64ByInterface(req.PathParameter("appid"))
	if nil != err {
		blog.Errorf("convert appid to int error:%v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsLostField, "ApplicationID")})
		return
	}

	input := new(meta.UpdateHostParams)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("updateHostByAppID , but decode body failed, err: %v,input:%+v,rid:%s", err, input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.V(5).Infof("updateHostByAppID http body data: %v,srvData.rid", input, srvData.rid)
	result, httpCode, errMsg := srvData.lgc.UpdateHostByAppID(srvData.ctx, input, appID)
	if nil != errMsg {
		blog.Errorf("updateHostByAppID update host, appID:%d, input:%+v, error:%s,rid:%s", appID, input, err, srvData.rid)
		resp.WriteError(httpCode, &meta.RespError{Msg: errMsg})
		return
	}
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	})

}

func (s *Service) HostSearchByIP(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	input := new(meta.HostSearchByIPParams)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("HostSearchByIP , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if 0 == len(input.IpList) {
		blog.Error("input does not contains key IP,input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsLostField, common.BKIPListField)})
	}

	orCondition := []map[string]interface{}{
		map[string]interface{}{common.BKHostInnerIPField: map[string]interface{}{common.BKDBIN: input.IpList}},
		map[string]interface{}{common.BKHostOuterIPField: map[string]interface{}{common.BKDBIN: input.IpList}},
	}
	hostMapCondition := map[string]interface{}{common.BKDBOR: orCondition}

	if nil != input.CloudID {
		hostMapCondition[common.BKCloudIDField] = input.CloudID
	}

	phpapi := srvData.lgc.NewPHPAPI()
	hostMap, hostIDArr, err := phpapi.GetHostMapByCond(srvData.ctx, hostMapCondition)
	if err != nil {
		blog.Errorf("HostSearchByIP error : %s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	configCond := map[string][]int64{
		common.BKHostIDField: hostIDArr,
	}
	if 0 < len(input.AppID) {
		configCond[common.BKAppIDField] = input.AppID
	}

	configData, err := srvData.lgc.GetConfigByCond(srvData.ctx, configCond)
	if nil != err {
		blog.Errorf("HostSearchByIP error : %s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}
	hostData, err := phpapi.SetHostData(srvData.ctx, configData, hostMap)
	if nil != err {
		blog.Errorf("HostSearchByIP error : %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostData,
	})
}

func (s *Service) HostSearchByConds(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("HostSearchByConds , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	phpapi := srvData.lgc.NewPHPAPI()
	hostMap, hostIDArr, err := phpapi.GetHostMapByCond(srvData.ctx, input)
	if err != nil {
		blog.Errorf("HostSearchByConds error : %v, input:%+v,rid:%s", err, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	configCond := map[string][]int64{
		common.BKHostIDField: hostIDArr,
	}
	configData, err := srvData.lgc.GetConfigByCond(srvData.ctx, configCond)
	if nil != err {
		blog.Errorf("HostSearchByConds error : %v, input:%+v", err, input)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}

	hostData, err := phpapi.SetHostData(srvData.ctx, configData, hostMap)
	if nil != err {
		blog.Errorf("HostSearchByConds error : %v, input:%+v,rid:%s", err, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostData,
	})
}

func (s *Service) HostSearchByModuleID(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	input := new(meta.HostSearchByModuleIDParams)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("HostSearchByModuleID , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if nil == input.ApplicationID {
		blog.Error("HostSearchByModuleID input does not contains key ApplicationID.input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsLostField, "ApplicationID")})
		return
	}

	if nil == input.ModuleID {
		blog.Error("HostSearchByModuleID input does not contains key ModuleID.input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsLostField, "ModuleID")})
		return
	}

	configData, err := srvData.lgc.GetConfigByCond(srvData.ctx, map[string][]int64{
		common.BKModuleIDField: input.ModuleID,
		common.BKAppIDField:    []int64{*input.ApplicationID},
	})
	if nil != err {
		blog.Errorf("HostSearchByModuleID get host module config error:%s, input:%+v", err.Error(), input)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	phpapi := srvData.lgc.NewPHPAPI()
	hostData, err := phpapi.GetHostDataByConfig(srvData.ctx, configData)
	if nil != err {
		blog.Errorf("HostSearchByModuleID get host module config error:%s, input:%+v", err.Error(), input)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostData,
	})
}

func (s *Service) HostSearchBySetID(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	input := new(meta.HostSearchBySetIDParams)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("HostSearchBySetID , but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	if nil == input.ApplicationID {
		blog.Error("HostSearchBySetID input does not contains key ApplicationID")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsLostField, "ApplicationID")})
		return
	}

	conds := make(map[string][]int64)
	conds[common.BKAppIDField] = []int64{*input.ApplicationID}

	if len(input.SetID) > 0 {
		conds[common.BKSetIDField] = input.SetID
	}

	configData, err := srvData.lgc.GetConfigByCond(srvData.ctx, conds)
	if nil != err {
		blog.Errorf("HostSearchBySetID get host module config error:%s, input:%+v", err.Error(), input)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}

	phpapi := srvData.lgc.NewPHPAPI()
	hostData, err := phpapi.GetHostDataByConfig(srvData.ctx, configData)
	if nil != err {
		blog.Errorf("HostSearchByModuleID get host module config error:%s, input:%+v", err.Error(), input)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostData,
	})
}

func (s *Service) HostSearchByAppID(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	input := new(meta.HostSearchByAppIDParams)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("HostSearchByAppID , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if nil == input.ApplicationID {
		blog.Error("HostSearchByAppID input does not contains key ApplicationID,input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsLostField, "ApplicationID")})
		return
	}

	configData, err := srvData.lgc.GetConfigByCond(srvData.ctx, map[string][]int64{
		common.BKAppIDField: []int64{*input.ApplicationID},
	})

	if nil != err {
		blog.Errorf("HostSearchByModuleID get host module config error:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}

	phpapi := srvData.lgc.NewPHPAPI()
	hostData, err := phpapi.GetHostDataByConfig(srvData.ctx, configData)
	if nil != err {
		blog.Errorf("HostSearchByModuleID get host module config error:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostData,
	})

}

func (s *Service) HostSearchByProperty(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("HostSearchByProperty , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	appID, err := util.GetInt64ByInterface(input[common.BKAppIDField])
	if nil != err {
		blog.Error("HostSearchByProperty input does not contains key ApplicationID.input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsLostField, "ApplicationID")})
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

	blog.V(5).Infof("HostSearchByProperty setCondition: %+v,rid:%s", setCondition, srvData.rid)
	setIDArr, err := srvData.lgc.GetSetIDByCond(srvData.ctx, setCondition)
	if nil != err {
		blog.Errorf("HostSearchByProperty get host module config error:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetSetFaild, err.Error())})
		return
	}
	blog.V(5).Infof("HostSearchByProperty ApplicationID: %s, SetID: %v,input:%+v,rid:%s", appID, setIDArr, input, srvData.rid)

	condition := map[string][]int64{
		common.BKAppIDField: []int64{appID},
	}

	condition[common.BKSetIDField] = setIDArr
	configData, err := srvData.lgc.GetConfigByCond(srvData.ctx, condition)
	if nil != err {
		blog.Errorf("HostSearchByProperty get host module config error:%s, input:%+v,param:%+v,rid:%s", err.Error(), input, condition, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}

	phpapi := srvData.lgc.NewPHPAPI()
	hostData, err := phpapi.GetHostDataByConfig(srvData.ctx, configData)
	if nil != err {
		blog.Errorf("HostSearchByProperty get host module config error:%s, input:%+v,params:%+v,rid:%s", err.Error(), input, configData, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostData,
	})
}

func (s *Service) GetIPAndProxyByCompany(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	input := new(meta.GetIPAndProxyByCompanyParams)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("GetIPAndProxyByCompany , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if 0 == len(input.Ips) {
		blog.Error("GetIPAndProxyByCompany input does not contains key IP.input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsLostField, "IP")})
		return
	}

	appIDInt, err := util.GetInt64ByInterface(*input.AppIDStr)
	if nil != err {
		blog.Errorf("GetIPAndProxyByCompany input application id not integer, input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "ApplicationID")})
		return
	}

	platIDInt, err := util.GetInt64ByInterface(*input.CloudIDStr)
	if nil != err {
		blog.Errorf("GetIPAndProxyByCompany cloud id not integer, input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "PlatID")})
		return
	}
	resData, err := srvData.lgc.GetIPAndProxyByCompany(srvData.ctx, input.Ips, platIDInt, appIDInt)
	if nil != err {
		blog.Errorf("GetIPAndProxyByCompany error:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     resData,
	})

}

func (s *Service) UpdateCustomProperty(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("UpdateCustomProperty , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	blog.V(5).Infof("UpdateCustomProperty  input:%+v,rid:%s", input, srvData.rid)
	appID, err := util.GetInt64ByInterface(input[common.BKAppIDField])
	if nil != err {
		blog.Errorf("UpdateCustomProperty input not found appID, input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID")})
		return
	}
	hostID, err := util.GetInt64ByInterface(input[common.BKHostIDField])
	if nil != err {
		blog.Errorf("UpdateCustomProperty input not found hostID, input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, "HostID")})
		return
	}
	propertyJson, ok := input["property"].(string)
	if false == ok && "" == propertyJson {
		blog.Errorf("UpdateCustomPropertyinput not found property, input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, "property")})
		return
	}

	propertyMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(propertyJson), &propertyMap)
	if nil != err {
		blog.Errorf("UpdateCustomPropertyinput not found property, input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	res, err := srvData.lgc.UpdateCustomProperty(srvData.ctx, hostID, appID, propertyMap)
	if nil != err {
		blog.Errorf("UpdateCustomPropertyinput not found property, input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     res,
	})

}

func (s *Service) GetHostAppByCompanyId(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	input := &meta.GetHostAppByCompanyIDParams{}
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("GetHostAppByCompanyId , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.V(5).Infof("GetHostAppByCompanyId input:%+v,rid:%s", input, srvData.rid)
	platId, err := util.GetInt64ByInterface(input.CloudIDStr)
	if nil != err {
		blog.Errorf("GetHostAppByCompanyId cloud id not integer, input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "PlatID")})
		return
	}
	ipArr := strings.Split(input.IPs, ",")
	hostCon := map[string]interface{}{
		common.BKHostInnerIPField: map[string]interface{}{
			common.BKDBIN: ipArr,
		},
		common.BKCloudIDField: platId,
	}

	phpapi := srvData.lgc.NewPHPAPI()
	//根据i,platId获取主机
	hostArr, hostIdArr, err := phpapi.GetHostMapByCond(srvData.ctx, hostCon) // phpapilogic.GetHostMapByCond(req, hostCon)
	if nil != err {
		blog.Errorf("GetHostAppByCompanyId getHostMapByCond:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetFail)})
		return
	}
	blog.V(5).Infof("GetHostAppByCompanyId hostArr:%v, input:%+v,rid:%s", hostArr, input, srvData.rid)
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
	configArr, err := srvData.lgc.GetConfigByCond(srvData.ctx, configCon)
	if nil != err {
		blog.Errorf("GetHostAppByCompanyId getConfigByCond err:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetFail)})
		return
	}
	blog.V(5).Infof("GetHostAppByCompanyId configArr:%v, input:%+v,rid:%s", configArr, input, srvData.rid)
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
	hostMapArr, err := phpapi.SetHostData(srvData.ctx, configArr, hostArr) //phpapilogic.SetHostData(req, configArr, hostArr)
	if nil != err {
		blog.Errorf("GetHostAppByCompanyId setHostData err:%s, input:%+v", err.Error(), input)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
		return
	}
	blog.V(5).Infof("GetHostAppByCompanyId hostMap:%v, input:%+v,rid:%s", hostMapArr, input, srvData.rid)
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostMapArr,
	})
}

func (s *Service) DelHostInApp(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	input := &meta.DelHostInAppParams
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("DelHostInApp , but decode body failed, err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	blog.V(5).Infof("DelHostInApp input:%+v,rid:%s", input, srvData.rid)
	appID, err := util.GetInt64ByInterface(input.AppID)
	if nil != err {
		blog.Errorf("GetHostAppByCompanyId cloud id not integer, input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "ApplicationID")})
		return
	}
	hostID, err := util.GetInt64ByInterface(input.HostID)
	if nil != err {
		blog.Errorf("GetHostAppByCompanyId host id not integer, input:%+v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "HostID")})
		return
	}
	configCon := map[string][]int64{
		common.BKAppIDField:  []int64{appID},
		common.BKHostIDField: []int64{hostID},
	}

	configArr, err := srvData.lgc.GetConfigByCond(srvData.ctx, configCon)
	if err != nil {
		blog.Errorf("DelHostInApp GetConfigByCond err msg:%v, input:%+v,rid:%s", err, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}
	if len(configArr) == 0 {
		blog.Errorf("DelHostInApp not fint hostId:%v in appId:%v, input:%+v,rid:%s", hostID, appID, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrHostNotFound)})
		return
	}
	moduleIdArr := make([]int64, 0)
	for _, item := range configArr {
		moduleIdArr = append(moduleIdArr, item[common.BKModuleIDField])
	}
	moduleCon := mapstr.MapStr{
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: moduleIdArr,
		},
		common.BKDefaultField: common.DefaultResModuleFlag,
	}
	moduleArr, err := srvData.lgc.GetModuleMapByCond(srvData.ctx, []string{common.BKModuleIDField}, moduleCon)
	if err != nil {
		blog.Errorf("DelHostInApp GetConfigByCond err msg, error:%s, input:%s", err.Error(), input)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetModuleFail, err.Error())})
		return
	}
	blog.V(5).Infof("DelHostInApp moduleArr:%v, input:%+v,rid:%s", moduleArr, input, srvData.rid)
	if len(moduleArr) == 0 {
		blog.Errorf("DelHostInApp GetModuleMapByCond   not find host in idle module input: %v,rid:%s", input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostNotFound)})
		return
	}
	param := make(common.KvMap)
	param[common.BKAppIDField] = appID
	param[common.BKHostIDField] = hostID
	res, err := s.CoreAPI.ObjectController().OpenAPI().DeleteSetHost(srvData.ctx, req.Request.Header, param)
	if nil != err {
		blog.Errorf("DelHostInApp DeleteSetHost   error:%s,  input:%+v,param:%+v,,rid:%s", err.Error(), input, param, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostEditRelationPoolFail)})
		return
	}
	if false == res.Result {
		blog.Errorf("DelHostInApp DeleteSetHost   error:%s,  input:%+v,param:%+v,rid:%s", res.ErrMsg, input, param, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(res.Code, res.ErrMsg)})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     res.Data,
	})

}

func (s *Service) GetGitServerIp(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	input := new(meta.GitServerIpParams)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("GetGitServerIp , but decode body failed, err: %s,rid:%s", err.Error(), srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	var appID, setID, moduleID int64

	// 根据appName获取app
	appCondition := mapstr.MapStr{
		common.BKAppNameField: input.AppName,
	}
	appMap, err := srvData.lgc.GetAppMapByCond(srvData.ctx, nil, appCondition) //  logics.GetAppMapByCond(req, "", cli.CC.ObjCtrl(), appCondition)
	if nil != err {
		blog.Errorf("GetGitServerIp GetAppMapByCond error:%s, input:%s,param:%+v,rid:%s", err.Error(), input, appCondition, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetAPPFail, err.Error())})
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
	setMap, err := srvData.lgc.GetSetMapByCond(srvData.ctx, nil, setCondition)
	if nil != err {
		blog.Errorf("GetGitServerIp GetSetMapByCond error:%s, input:%s,param:%+v,rid:%s", err.Error(), input, setCondition, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetSetFaild, err.Error())})
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
	moduleCondition := mapstr.MapStr{
		common.BKModuleNameField: input.ModuleName,
		common.BKAppIDField:      appID,
	}
	moduleMap, err := srvData.lgc.GetModuleMapByCond(srvData.ctx, nil, moduleCondition)
	if nil != err {
		blog.Errorf("GetGitServerIp GetModuleMapByCond error:%s, input:%s,param:%s,rid:%s", err.Error(), input, moduleCondition, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetSetFaild, err.Error())})
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
	configData, err := srvData.lgc.GetConfigByCond(srvData.ctx, confMap)
	if nil != err {
		blog.Errorf("GetGitServerIp GetModuleMapByCond error:%s, input:%+v,param:%+v,rid:%s", err.Error(), input, confMap, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
		return
	}
	blog.V(5).Infof("GetGitServerIp configData:%v", configData)
	phpapi := srvData.lgc.NewPHPAPI()
	hostArr, err := phpapi.GetHostDataByConfig(srvData.ctx, configData)
	if nil != err {
		blog.Errorf("GetGitServerIp getHostDataByConfig error:%s, input:%+v,param:%+v,rid:%s", err.Error(), input, configData, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
		return
	}
	blog.V(5).Infof("GetGitServerIp hostArr:%v, input:%+v,rid:%s", hostArr, input, srvData.rid)

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostArr,
	})
}

func (s *Service) GetPlat(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	params := new(meta.QueryCondition)
	res, err := s.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDPlat, params)
	if nil != err {
		blog.Errorf("GetPlat htt do error: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrTopoGetCloudErrStrFaild, err.Error())})
		return
	}
	if false == res.Result {
		blog.Errorf("GetPlat http reply error. err code:%d, err msg:%s,rid:%s", res.Code, res.ErrMsg, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(res.Code, res.ErrMsg)})

	} else {
		resp.WriteEntity(meta.Response{
			BaseResp: meta.SuccessBaseResp,
			Data:     res.Data,
		})
	}

}

func (s *Service) CreatePlat(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); nil != err {
		blog.Errorf("CreatePlat , but decode body failed, err: %s,rid:%s", err.Error(), srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	ownerId := util.GetOwnerID(req.Request.Header)
	input[common.BKOwnerIDField] = ownerId

	valid := validator.NewValidMap(util.GetOwnerID(req.Request.Header), common.BKInnerObjIDPlat, srvData.header, s.Engine)
	validErr := valid.ValidMap(input, common.ValidCreate, 0)

	if nil != validErr {
		blog.Errorf("CreatePlat error: %v, input:%+v,rid:%s", validErr, input, srvData.rid)
		if se, ok := validErr.(errors.CCErrorCoder); ok {
			if se.GetCode() == common.CCErrCommDuplicateItem {
				resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommDuplicateItem, "")})
			}
		}
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrTopoInstCreateFailed)})
		return
	}

	instInfo := &meta.CreateModelInstance{
		Data: mapstr.NewFromMap(input),
	}

	res, err := s.CoreAPI.CoreService().Instance().CreateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDPlat, instInfo)
	if nil != err {
		blog.Errorf("CreatePlat error: %s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrTopoInstCreateFailed)})
		return
	}

	if false == res.Result {
		blog.Errorf("GetPlat error.err code:%d,err msg:%s,input:%+v,rid:%s", res.Code, res.ErrMsg, input, srvData.rid)
		resp.WriteHeaderAndJson(http.StatusInternalServerError, res, common.BKHTTPMIMEJSON)

	}
	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     res.Data,
	})

}

func (s *Service) DelPlat(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	platID, convErr := util.GetInt64ByInterface(req.PathParameter(common.BKCloudIDField))
	if nil != convErr || 0 == platID {
		blog.Errorf("the platID is invalid, error info is %s, input:%s.rid:%s", convErr.Error(), platID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, convErr.Error())})
		return
	}

	params := new(meta.QueryInput)
	params.Fields = common.BKHostIDField
	params.Condition = map[string]interface{}{
		common.BKCloudIDField: platID,
	}

	hostRes, err := s.CoreAPI.HostController().Host().GetHosts(srvData.ctx, srvData.header, params)
	if nil != err {
		blog.Errorf("DelPlat search host error: %s, input:%+v,rid:%s", err.Error(), platID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetFail)})
		return
	}
	if !hostRes.Result {
		blog.Errorf("DelPlat search host http response error.err code:%d,err msg:%s, input:%+v,rid:%s", hostRes.Code, hostRes.ErrMsg, platID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetFail)})
		return
	}

	if 0 < hostRes.Data.Count {
		blog.Errorf("DelPlat plat [%d] has host data, can not delete,rid:%s", platID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrTopoHasHostCheckFailed)})
		return
	}

	delCond := &meta.DeleteOption{
		Condition: mapstr.MapStr{common.BKCloudIDField: platID},
	}

	res, err := s.CoreAPI.CoreService().Instance().DeleteInstance(srvData.ctx, srvData.header, common.BKInnerObjIDPlat, delCond)
	if nil != err {
		blog.Errorf("DelPlat do error: %v, input:%d,rid:%s", err, platID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrTopoInstDeleteFailed)})
		return
	}

	if false == res.Result {
		blog.Errorf("DelPlat http reponse error. err code:%d,err msg:%s,input:%s,rid:%s", res.Code, res.ErrMsg, platID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(res.Code, res.ErrMsg)})

	} else {
		resp.WriteEntity(meta.Response{
			BaseResp: meta.SuccessBaseResp,
			Data:     "",
		})
	}
}

func (s *Service) getHostListByAppidAndField(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	// 获取AppID
	pathParams := req.PathParameters()
	appID, err := util.GetInt64ByInterface(pathParams[common.BKAppIDField])
	if nil != err {
		blog.Errorf("getHostListByAppidAndField error :%s,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, err.Error())})
		return
	}

	field := req.PathParameter("field")

	configData, err := srvData.lgc.GetConfigByCond(srvData.ctx, map[string][]int64{
		common.BKAppIDField: []int64{appID},
	})

	if nil != err {
		blog.Errorf("getHostListByAppidAndField error : %s, input:%+v,rid:%s", err.Error(), common.KvMap{"appid": appID, "field": field}, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostModuleConfigFaild, err.Error())})
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
	ret, err := srvData.lgc.CoreAPI.HostController().Host().GetHosts(srvData.ctx, req.Request.Header, query)
	if nil != err {
		blog.Errorf("getHostListByAppidAndField search host error: %s, input:%+v,rid:%s", err.Error(), common.KvMap{"appid": appID, "field": field}, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetFail)})
		return
	}
	if !ret.Result {
		blog.Errorf("getHostListByAppidAndField search host error. err code:%d,err msg:%s, input:%+v,rid:%s", ret.Code, ret.ErrMsg, common.KvMap{"appid": appID, "field": field}, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(ret.Code, ret.ErrMsg)})
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
