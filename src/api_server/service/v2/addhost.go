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

package v2

import (
	"strings"

	"configcenter/src/api_server/logics/v2/common/converter"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func (s *Service) addHost(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr
	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("addHost error:%v,rid:%s", err, srvData.rid)
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
		blog.Errorf("addHost error ip empty,input:%#v,rid:%s", formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "ip").Error(), resp)
		return
	}
	ipArr := strings.Split(ips, ",")

	intPlatID, err := util.GetInt64ByInterface(platID)
	if nil != err {
		blog.Errorf("addHost error platID not integer, input:%#v,rid:%s", formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "platId").Error(), resp)
		return
	}

	param := &metadata.HostToAppModule{Ips: ipArr,
		ModuleName:  moduleName,
		SetName:     setName,
		AppName:     appName,
		OwnerID:     srvData.ownerID,
		PlatID:      intPlatID,
		IsIncrement: true}
	result, err := s.CoreAPI.HostServer().AssignHostToAppModule(srvData.ctx, srvData.header, param)
	if err != nil {
		blog.Errorf("addHost  error:%v,input:%#v,codition:%#v,rid:%s", err, formData, param, srvData.rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("addHost  error, err code:%d,err msg:%s,input:%#v,codition:%#v,rid:%s", result.Code, result.ErrMsg, formData, param, srvData.rid)
		converter.RespFailV2Error(defErr.New(result.Code, result.ErrMsg), resp)
		return
	}

	err = converter.ResToV2ForEnterIP(result.Result, result.ErrMsg, result.Data)

	if err != nil {
		blog.Errorf("convert addHost result to v2 error:%s, reply:%v", err.Error(), result.Data)
		converter.RespFailV2(common.CCErrAddHostToModuleFailStr, err.Error(), resp)
		return
	}
	converter.RespSuccessV2("", resp)
}

func (s *Service) enterIP(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	err := req.Request.ParseForm()
	if err != nil {
		blog.Errorf("enterIP error:%v,rid:%s", err, srvData.rid)
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
		blog.Errorf("enterIP error ips empty,input:%+v,rid:%s", formData, srvData.rid)
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
		blog.Errorf("osType mast be windows or linux; not %s,input:%+v,rid:%s", osType, formData, srvData.rid)
		converter.RespFailV2(common.CCErrAPIServerV2OSTypeErr, defErr.Error(common.CCErrAPIServerV2OSTypeErr).Error(), resp)
		return
	}
	osTypeEnumKey := "1"
	if "windows" == osType {
		osTypeEnumKey = "2"
	}

	param := &metadata.HostToAppModule{Ips: ipArr,
		HostName:    hostNameArr,
		ModuleName:  moduleName,
		SetName:     setName,
		AppName:     appName,
		OwnerID:     srvData.ownerID,
		OsType:      osTypeEnumKey,
		IsIncrement: true}
	result, err := s.CoreAPI.HostServer().AssignHostToAppModule(srvData.ctx, srvData.header, param)
	if err != nil {
		blog.Errorf("enterIP  error:%v,input:%+v.condition:%#v,rid:%s", err, formData, param, srvData.rid)
		converter.RespFailV2(common.CCErrCommHTTPDoRequestFailed, defErr.Error(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	if !result.Result {
		blog.Errorf("enterIP  error, err code:%d,err msg:%s,input:%+v.condition:%#v,rid:%s", result.Code, result.ErrMsg, formData, param, srvData.rid)
		converter.RespFailV2Error(defErr.New(result.Code, result.ErrMsg), resp)
		return
	}

	err = converter.ResToV2ForEnterIP(result.Result, result.ErrMsg, result.Data)
	if err != nil {
		blog.Errorf("convert enterip result to v2 error:%v, reply:%v,input:%+v,rid:%s", err.Error(), result.Data, formData, srvData.rid)
		converter.RespFailV2(common.CCErrCommReplyDataFormatError, defErr.Error(common.CCErrCommReplyDataFormatError).Error(), resp)
		return
	}
	converter.RespSuccessV2("", resp)
}
